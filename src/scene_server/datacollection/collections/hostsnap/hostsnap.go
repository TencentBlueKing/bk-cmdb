/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hostsnap

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"

	"github.com/tidwall/gjson"
)

const (
	// defaultChangeRangePercent is the value of the default percentage of data fluctuation
	defaultChangeRangePercent = 10
	// minChangeRangePercent is the value of the minimum percentage of data fluctuation
	minChangeRangePercent = 1
	// defaultRateLimiterQPS is the default value of rateLimiter qps
	defaultRateLimiterQPS = 40
	// defaultRateLimiterBurst is the default value of rateLimiter burst
	defaultRateLimiterBurst = 100
)

var (
	// 需要参与变化对比的字段
	compareFields = []string{"bk_cpu", "bk_cpu_module", "bk_cpu_mhz", "bk_disk", "bk_mem", "bk_os_type", "bk_os_name",
		"bk_os_version", "bk_host_name", "bk_outer_mac", "bk_mac", "bk_os_bit"}
	reqireFields = append(compareFields, "bk_host_id", "bk_host_innerip", "bk_host_outerip")

	// notice: 为了对应不同版本和环境差异，再当前版本中设置compareFields中不参加对比的字段
	ignoreCompareField = make(map[string]struct{}, 0)
)

type HostSnap struct {
	redisCli    redis.Client
	authManager *extensions.AuthManager
	*backbone.Engine
	rateLimit flowctrl.RateLimiter
	filter    *filter
	ctx       context.Context
	db        dal.RDB
	window    *Window
}

func NewHostSnap(ctx context.Context, redisCli redis.Client, db dal.RDB, engine *backbone.Engine, authManager *extensions.AuthManager) *HostSnap {
	qps, burst := getRateLimiterConfig()
	h := &HostSnap{
		redisCli:    redisCli,
		ctx:         ctx,
		db:          db,
		rateLimit:   flowctrl.NewRateLimiter(int64(qps), int64(burst)),
		authManager: authManager,
		Engine:      engine,
		filter:      newFilter(),
		window:      newWindow(),
	}
	return h
}

func getRateLimiterConfig() (int, int) {
	qps, err := cc.Int("datacollection.hostsnap.rateLimiter.qps")
	if err != nil {
		blog.Errorf("can't find the value of datacollection.hostsnap.rateLimiter.qps settings, set the default value: %s", defaultRateLimiterQPS)
		qps = defaultRateLimiterQPS
	}
	burst, err := cc.Int("datacollection.hostsnap.rateLimiter.burst")
	if err != nil {
		blog.Errorf("can't find the value of datacollection.hostsnap.rateLimiter.burst setting,set the default value: %s", defaultRateLimiterBurst)
		burst = defaultRateLimiterBurst
	}
	return qps, burst
}

// Hash returns hash value base on message.
func (h *HostSnap) Hash(cloudid, ip string) (string, error) {
	if len(cloudid) == 0 {
		return "", fmt.Errorf("can't make hash from invalid message format, cloudid empty")
	}
	if len(ip) == 0 {
		return "", fmt.Errorf("can't make hash from invalid message format, ip empty")
	}

	hash := fmt.Sprintf("%s:%s", cloudid, ip)

	return hash, nil
}

// Mock returns local mock message for testing.
func (h *HostSnap) Mock() string {
	return MockMessage
}

// getLimitConfig return the configured limit value
func getLimitConfig(config string, defaultValue, minValue int) int {
	var err error
	limit := defaultValue
	if cc.IsExist(config) {
		limit, err = cc.Int(config)
		if err != nil {
			blog.Errorf("get %s value error, err: %v ", config, err)
			limit = defaultValue
		}
		if limit < minValue {
			limit = minValue
		}
	}
	return limit
}

func (h *HostSnap) Analyze(msg *string) error {
	if msg == nil {
		return fmt.Errorf("message nil")
	}

	var data string

	if !gjson.Get(*msg, "cloudid").Exists() {
		data = gjson.Get(*msg, "data").String()
	} else {
		data = *msg
	}

	header, rid := newHeaderWithRid()

	val := gjson.Parse(data)
	cloudID := val.Get("cloudid").Int()
	ips := getIPS(&val)
	host, err := h.getHostByVal(header, cloudID, ips, &val)
	if err != nil {
		blog.Errorf("get host detail with ips: %v failed, err: %v, rid: %s", ips, err, rid)
		return err
	}
	elements := gjson.GetMany(host, common.BKHostIDField, common.BKHostInnerIPField, common.BKHostOuterIPField)
	// check host id field
	if !elements[0].Exists() {
		blog.Errorf("snapshot analyze, but host id not exist, host: %s, ips: %v, rid: %s", host, ips, rid)
		return errors.New("host id not exist")
	}
	hostID := elements[0].Int()
	if hostID == 0 {
		blog.Errorf("snapshot analyze, but host id is 0, host: %s, ips: %v, rid: %s", host, ips, rid)
		return errors.New("host id can not be 0")
	}

	// check inner ip
	if !elements[1].Exists() {
		blog.Errorf("snapshot analyze, but host inner ip not exist, host: %s, ips: %v, rid: %s", host, ips, rid)
		return errors.New("host inner ip not exist")
	}

	innerIP := elements[1].String()
	outerIP := elements[2].String()

	// save host snapshot in redis

	if !val.Get("data.apiVer").Exists() {
		h.saveHostsnap(header, &val, hostID)
	}

	// window restriction on request
	if !h.window.canPassWindow() {
		if blog.V(4) {
			blog.Infof("not within the time window that can pass, skip host snapshot data update, host id: %d, ip: %s, cloud id: %d, rid: %s",
				hostID, innerIP, cloudID, rid)
		}
		return nil
	}
	setter, raw := parseSetter(&val, innerIP, outerIP)
	// no need to update
	if !needToUpdate(raw, host) {
		return nil
	}

	// limit the number of requests
	if !h.rateLimit.TryAccept() {
		blog.Warnf("skip host snapshot data update due to request limit, host id: %d, ip: %s, cloud id: %d, rid: %s",
			hostID, innerIP, cloudID, rid)
		return nil
	}

	blog.V(5).Infof("snapshot for host changed, need update, host id: %d, ip: %s, cloud id: %d, from %s to %s, rid: %s",
		hostID, innerIP, cloudID, host, raw, rid)

	// get audit interface of host.
	audit := auditlog.NewHostAudit(h.CoreAPI.CoreService())
	kit := &rest.Kit{
		Rid:             rid,
		Header:          header,
		Ctx:             h.ctx,
		CCError:         h.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)),
		User:            common.CCSystemCollectorUserName,
		SupplierAccount: common.BKDefaultOwnerID,
	}

	// generate audit log for update host.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).
		WithOperateFrom(metadata.FromDataCollection).WithUpdateFields(setter)
	auditLog, err := audit.GenerateAuditLogByHostIDGetBizID(generateAuditParameter, hostID, innerIP, nil)
	if err != nil {
		blog.Errorf("generate host snap audit log failed before update host, host: %d/%s, err: %v, rid: %s", hostID, innerIP, err, rid)
		return err
	}

	// notice: needToUpdate 需要顺序，只能在更新数据库之前，删除需要忽略更新的字段
	for field := range ignoreCompareField {
		delete(setter, field)
	}

	opt := &metadata.UpdateOption{
		Condition: map[string]interface{}{
			common.BKHostIDField: hostID,
		},
		Data:       setter,
		CanEditAll: true,
	}

	res, err := h.CoreAPI.CoreService().Instance().UpdateInstance(h.ctx, header, common.BKInnerObjIDHost, opt)
	if err != nil {
		blog.Errorf("snapshot changed, update host %d/%s snapshot failed, err: %v, rid: %s", hostID, innerIP, err, rid)
		return err
	}
	if !res.Result {
		blog.Errorf("snapshot changed, update host %d/%s snapshot failed, err: %s, rid: %s", hostID, innerIP, res.ErrMsg, rid)
		return fmt.Errorf("update snapshot failed, err: %s", res.ErrMsg)
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("save host snap audit log failed after update host, host %d/%s, err: %v, rid: %s", hostID, innerIP, err, rid)
		return err
	}

	blog.V(5).Infof("snapshot for host changed, update success, host id: %d, ip: %s, cloud id: %d, rid: %s",
		hostID, innerIP, cloudID, rid)

	return nil
}

func needToUpdate(src, toCompare string) bool {
	// get data fluctuation limit
	changeRangePercent := getLimitConfig("datacollection.hostsnap.changeRangePercent", defaultChangeRangePercent, minChangeRangePercent)
	srcElements := gjson.GetMany(src, compareFields...)
	compareElements := gjson.GetMany(toCompare, compareFields...)
	for idx, field := range compareFields {
		if _, ok := ignoreCompareField[field]; ok {
			// 忽略变更对比的字段直接过滤掉
			continue
		}
		// compare these value with string directly to avoid empty value or null value.
		if srcElements[idx].String() != compareElements[idx].String() {
			compareField := compareFields[idx]
			// tolerate bk_cpu, bk_cpu_mhz, bk_disk, bk_mem changes less than the set value
			if compareField == "bk_cpu" || compareField == "bk_cpu_mhz" || compareField == "bk_disk" || compareField == "bk_mem" {
				val := compareElements[idx].Float() * (float64(changeRangePercent) / 100.0)
				diff := srcElements[idx].Float() - compareElements[idx].Float()
				if -val < diff && diff < val {
					continue
				}
			}
			return true
		}
	}
	return false
}

func parseSetter(val *gjson.Result, innerIP, outerIP string) (map[string]interface{}, string) {
	if val.Get("data.apiVer").String() == "v1.0" {
		return parseV10Setter(val, innerIP, outerIP)
	}

	var cpumodule = val.Get("data.cpu.cpuinfo.0.modelName").String()
	cpumodule = strings.TrimSpace(cpumodule)
	var cpunum int64
	for _, core := range val.Get("data.cpu.cpuinfo.#.cores").Array() {
		cpunum += core.Int()
	}
	var CPUMhz = val.Get("data.cpu.cpuinfo.0.mhz").Int()
	var disk uint64
	for _, disktotal := range val.Get("data.disk.usage.#.total").Array() {
		disk += disktotal.Uint() >> 10 >> 10 >> 10
	}
	var mem = val.Get("data.mem.meminfo.total").Uint()
	var hostname = val.Get("data.system.info.hostname").String()
	hostname = strings.TrimSpace(hostname)
	var ostype = val.Get("data.system.info.os").String()
	ostype = strings.TrimSpace(ostype)
	var osname string
	platform := val.Get("data.system.info.platform").String()
	platform = strings.TrimSpace(platform)
	version := val.Get("data.system.info.platformVersion").String()
	switch strings.ToLower(ostype) {
	case "linux":
		version = strings.Replace(version, ".x86_64", "", 1)
		version = strings.Replace(version, ".i386", "", 1)
		osname = fmt.Sprintf("%s %s", ostype, platform)
		ostype = common.HostOSTypeEnumLinux
	case "windows":
		version = strings.Replace(version, "Microsoft ", "", 1)
		platform = strings.Replace(platform, "Microsoft ", "", 1)
		osname = fmt.Sprintf("%s", platform)
		ostype = common.HostOSTypeEnumWindows
	case "aix":
		osname = platform
		ostype = common.HostOSTypeEnumAIX
	default:
		osname = fmt.Sprintf("%s", platform)
	}
	version = strings.TrimSpace(version)
	osname = strings.TrimSpace(osname)
	var OuterMAC, InnerMAC string
	for _, inter := range val.Get("data.net.interface").Array() {
		for _, addr := range inter.Get("addrs.#.addr").Array() {
			splitAddr := strings.Split(addr.String(), "/")
			if len(splitAddr) == 0 {
				continue
			}
			ip := splitAddr[0]
			if ip == innerIP {
				InnerMAC = inter.Get("hardwareaddr").String()
				InnerMAC = strings.TrimSpace(InnerMAC)
			} else if ip == outerIP {
				OuterMAC = inter.Get("hardwareaddr").String()
				OuterMAC = strings.TrimSpace(OuterMAC)
			}
		}
	}

	osbit := val.Get("data.system.info.systemtype").String()
	osbit = strings.TrimSpace(osbit)
	mem = mem >> 10 >> 10

	setter := make(map[string]interface{})
	raw := strings.Builder{}
	raw.WriteByte('{')

	if cpunum <= 0 {
		blog.V(4).Infof("bk_cpu not found in message for %s", innerIP)
	} else {
		setter["bk_cpu"] = cpunum
		raw.WriteString("\"bk_cpu\":")
		raw.WriteString(strconv.FormatInt(cpunum, 10))

	}

	if cpumodule == "" {
		blog.V(4).Infof("bk_cpu_module not found in message for %s", innerIP)
	} else {
		setter["bk_cpu_module"] = cpumodule
		raw.WriteString(",")
		raw.WriteString("\"bk_cpu_module\":")
		raw.Write([]byte("\"" + cpumodule + "\""))
	}

	if CPUMhz <= 0 {
		blog.V(4).Infof("bk_cpu_mhz not found in message for %s", innerIP)
	} else {
		setter["bk_cpu_mhz"] = CPUMhz
		raw.WriteString(",")
		raw.WriteString("\"bk_cpu_mhz\":")
		raw.WriteString(strconv.FormatInt(CPUMhz, 10))
	}

	if disk <= 0 {
		blog.V(4).Infof("bk_disk not found in message for %s", innerIP)
	} else {
		setter["bk_disk"] = disk
		raw.WriteString(",")
		raw.WriteString("\"bk_disk\":")
		raw.WriteString(strconv.FormatUint(disk, 10))
	}

	if mem <= 0 {
		blog.V(4).Infof("bk_mem not found in message for %s", innerIP)
	} else {
		setter["bk_mem"] = mem
		raw.WriteString(",")
		raw.WriteString("\"bk_mem\":")
		raw.WriteString(strconv.FormatUint(mem, 10))
	}

	if ostype == "" {
		blog.V(4).Infof("bk_os_type not found in message for %s", innerIP)
	} else {
		setter["bk_os_type"] = ostype
		raw.WriteString(",")
		raw.WriteString("\"bk_os_type\":")
		raw.Write([]byte("\"" + ostype + "\""))
	}

	if osname == "" {
		blog.V(4).Infof("bk_os_name not found in message for %s", innerIP)
	} else {
		setter["bk_os_name"] = osname
		raw.WriteString(",")
		raw.WriteString("\"bk_os_name\":")
		raw.Write([]byte("\"" + osname + "\""))
	}

	if version == "" {
		blog.V(4).Infof("bk_os_version not found in message for %s", innerIP)
	} else {
		setter["bk_os_version"] = version
		raw.WriteString(",")
		raw.WriteString("\"bk_os_version\":")
		raw.Write([]byte("\"" + version + "\""))
	}

	if hostname == "" {
		blog.V(4).Infof("bk_host_name not found in message for %s", innerIP)
	} else {
		setter["bk_host_name"] = hostname
		raw.WriteString(",")
		raw.WriteString("\"bk_host_name\":")
		raw.Write([]byte("\"" + hostname + "\""))
	}

	if outerIP != "" && OuterMAC == "" {
		blog.V(4).Infof("bk_outer_mac not found in message for %s", innerIP)
	} else {
		setter["bk_outer_mac"] = OuterMAC
		raw.WriteString(",")
		raw.WriteString("\"bk_outer_mac\":")
		raw.Write([]byte("\"" + OuterMAC + "\""))
	}

	if InnerMAC == "" {
		blog.V(4).Infof("bk_mac not found in message for %s", innerIP)
	} else {
		setter["bk_mac"] = InnerMAC
		raw.WriteString(",")
		raw.WriteString("\"bk_mac\":")
		raw.Write([]byte("\"" + InnerMAC + "\""))
	}

	if osbit == "" {
		blog.V(4).Infof("bk_os_bit not found in message for %s", innerIP)
	} else {
		setter["bk_os_bit"] = osbit
		raw.WriteString(",")
		raw.WriteString("\"bk_os_bit\":")
		raw.Write([]byte("\"" + osbit + "\""))
	}

	raw.WriteByte('}')

	return setter, raw.String()
}

func parseV10Setter(val *gjson.Result, innerIP, outerIP string) (map[string]interface{}, string) {
	var cpumodule = val.Get("data.cpu.model").String()
	cpumodule = strings.TrimSpace(cpumodule)
	var cpunum = val.Get("data.cpu.total").Int()
	var disk = val.Get("data.disk.total").Uint() >> 10 >> 10 >> 10
	var mem = val.Get("data.mem.total").Uint() >> 10 >> 10
	var hostname = val.Get("data.system.hostname").String()
	hostname = strings.TrimSpace(hostname)
	var ostype = val.Get("data.system.os").String()
	ostype = strings.TrimSpace(ostype)
	var osname string
	platform := val.Get("data.system.platform").String()
	platform = strings.TrimSpace(platform)
	version := val.Get("data.system.platVer").String()
	switch strings.ToLower(ostype) {
	case "linux":
		version = strings.Replace(version, ".x86_64", "", 1)
		version = strings.Replace(version, ".i386", "", 1)
		osname = fmt.Sprintf("%s %s", ostype, platform)
		ostype = common.HostOSTypeEnumLinux
	case "windows":
		version = strings.Replace(version, "Microsoft ", "", 1)
		platform = strings.Replace(platform, "Microsoft ", "", 1)
		osname = fmt.Sprintf("%s", platform)
		ostype = common.HostOSTypeEnumWindows
	case "aix":
		osname = platform
		ostype = common.HostOSTypeEnumAIX
	default:
		osname = fmt.Sprintf("%s", platform)
	}
	version = strings.TrimSpace(version)
	osname = strings.TrimSpace(osname)
	var OuterMAC, InnerMAC string
	for _, inter := range val.Get("data.net.interface").Array() {
		for _, addr := range inter.Get("addrs").Array() {
			splitAddr := strings.Split(addr.String(), "/")
			if len(splitAddr) == 0 {
				continue
			}
			ip := splitAddr[0]
			if ip == innerIP {
				InnerMAC = inter.Get("mac").String()
				InnerMAC = strings.TrimSpace(InnerMAC)
			} else if ip == outerIP {
				OuterMAC = inter.Get("mac").String()
				OuterMAC = strings.TrimSpace(OuterMAC)
			}
		}
	}

	osbit := val.Get("data.system.sysType").String()
	osbit = strings.TrimSpace(osbit)

	setter := make(map[string]interface{})
	raw := strings.Builder{}
	raw.WriteByte('{')

	if cpunum <= 0 {
		blog.V(4).Infof("bk_cpu not found in message for %s", innerIP)
	} else {
		setter["bk_cpu"] = cpunum
		raw.WriteString("\"bk_cpu\":")
		raw.WriteString(strconv.FormatInt(cpunum, 10))

	}

	if cpumodule == "" {
		blog.V(4).Infof("bk_cpu_module not found in message for %s", innerIP)
	} else {
		setter["bk_cpu_module"] = cpumodule
		raw.WriteString(",")
		raw.WriteString("\"bk_cpu_module\":")
		raw.Write([]byte("\"" + cpumodule + "\""))
	}

	if disk <= 0 {
		blog.V(4).Infof("bk_disk not found in message for %s", innerIP)
	} else {
		setter["bk_disk"] = disk
		raw.WriteString(",")
		raw.WriteString("\"bk_disk\":")
		raw.WriteString(strconv.FormatUint(disk, 10))
	}

	if mem <= 0 {
		blog.V(4).Infof("bk_mem not found in message for %s", innerIP)
	} else {
		setter["bk_mem"] = mem
		raw.WriteString(",")
		raw.WriteString("\"bk_mem\":")
		raw.WriteString(strconv.FormatUint(mem, 10))
	}

	if ostype == "" {
		blog.V(4).Infof("bk_os_type not found in message for %s", innerIP)
	} else {
		setter["bk_os_type"] = ostype
		raw.WriteString(",")
		raw.WriteString("\"bk_os_type\":")
		raw.Write([]byte("\"" + ostype + "\""))
	}

	if osname == "" {
		blog.V(4).Infof("bk_os_name not found in message for %s", innerIP)
	} else {
		setter["bk_os_name"] = osname
		raw.WriteString(",")
		raw.WriteString("\"bk_os_name\":")
		raw.Write([]byte("\"" + osname + "\""))
	}

	if version == "" {
		blog.V(4).Infof("bk_os_version not found in message for %s", innerIP)
	} else {
		setter["bk_os_version"] = version
		raw.WriteString(",")
		raw.WriteString("\"bk_os_version\":")
		raw.Write([]byte("\"" + version + "\""))
	}

	if hostname == "" {
		blog.V(4).Infof("bk_host_name not found in message for %s", innerIP)
	} else {
		setter["bk_host_name"] = hostname
		raw.WriteString(",")
		raw.WriteString("\"bk_host_name\":")
		raw.Write([]byte("\"" + hostname + "\""))
	}

	if outerIP != "" && OuterMAC == "" {
		blog.V(4).Infof("bk_outer_mac not found in message for %s", innerIP)
	} else {
		setter["bk_outer_mac"] = OuterMAC
		raw.WriteString(",")
		raw.WriteString("\"bk_outer_mac\":")
		raw.Write([]byte("\"" + OuterMAC + "\""))
	}

	if InnerMAC == "" {
		blog.V(4).Infof("bk_mac not found in message for %s", innerIP)
	} else {
		setter["bk_mac"] = InnerMAC
		raw.WriteString(",")
		raw.WriteString("\"bk_mac\":")
		raw.Write([]byte("\"" + InnerMAC + "\""))
	}

	if osbit == "" {
		blog.V(4).Infof("bk_os_bit not found in message for %s", innerIP)
	} else {
		setter["bk_os_bit"] = osbit
		raw.WriteString(",")
		raw.WriteString("\"bk_os_bit\":")
		raw.Write([]byte("\"" + osbit + "\""))
	}

	raw.WriteByte('}')

	return setter, raw.String()
}

func (h *HostSnap) getHostByVal(header http.Header, cloudID int64, ips []string, val *gjson.Result) (string, error) {
	rid := util.GetHTTPCCRequestID(header)

	if len(ips) == 0 {
		blog.Warnf("snapshot message has no ip, message:%s, rid: %s", val.String(), rid)
		return "", errors.New("snapshot has no ip fields")
	}

	for _, ip := range ips {
		if h.filter.Exist(ip, cloudID) {
			// skip the cached invalid ip which may not exist.
			continue
		}

		opt := &metadata.SearchHostWithInnerIPOption{
			InnerIP: ip,
			CloudID: cloudID,
			Fields:  reqireFields,
		}

		host, err := h.Engine.CoreAPI.CacheService().Cache().Host().SearchHostWithInnerIP(context.Background(), header, opt)
		if err != nil {
			blog.Errorf("get host info with ip: %s, cloud id: %d failed, err: %v, rid: %s", ip, cloudID, err, rid)
			if ccErr, ok := err.(ccErr.CCErrorCoder); ok {
				if ccErr.GetCode() == common.CCErrCommDBSelectFailed {
					h.filter.Set(ip, cloudID)
				}
			}
			// do not return, continue search with next ip
		}

		if len(host) == 0 {
			// not find host
			continue
		}

		return host, nil

	}

	return "", errors.New("can not find ip detail from cache")
}

func getIPS(val *gjson.Result) []string {
	ipv4 := make([]string, 0)
	ipv6 := make([]string, 0)

	rootIP := val.Get("ip").String()
	if !strings.HasPrefix(rootIP, "127.0.0.") && net.ParseIP(rootIP) != nil {
		if strings.Contains(rootIP, ":") {
			// not support ipv6 for now.
			// ipv6 = append(ipv6, rootIP)
		} else {
			ipv4 = append(ipv4, rootIP)
		}
	}

	interfaces := val.Get("data.net.interface.#.addrs.#.addr").Array()
	for _, addrs := range interfaces {
		for _, addr := range addrs.Array() {
			ip := strings.Split(addr.String(), "/")[0]
			if strings.HasPrefix(ip, "127.0.0.") {
				continue
			}

			if net.ParseIP(ip) == nil {
				// invalid ip address
				continue
			}

			if strings.Contains(ip, ":") {
				// not support ipv6 for now.
				// ipv6 = append(ipv6, ip)
			} else {
				ipv4 = append(ipv4, ip)
			}
		}
	}
	return append(ipv4, ipv6...)
}

// saveHostsnap save host snapshot in redis
func (h *HostSnap) saveHostsnap(header http.Header, hostData *gjson.Result, hostID int64) error {
	rid := util.GetHTTPCCRequestID(header)

	snapshot, err := ParseHostSnap(hostData)
	if err != nil {
		blog.Errorf("saveHostsnap failed, ParseHostSnap err: %v, hostID:%v, rid:%s", err, hostID, rid)
		return err
	}

	key := common.RedisSnapKeyPrefix + strconv.FormatInt(hostID, 10)
	if err := h.redisCli.Set(context.Background(), key, *snapshot, time.Minute*10).Err(); err != nil {
		blog.Errorf("saveHostsnap failed, set key: %s to redis err: %v, rid: %s", key, err, rid)
		return err
	}

	return nil
}

func newHeaderWithRid() (http.Header, string) {
	header := http.Header{}
	header.Add(common.BKHTTPOwnerID, common.BKDefaultOwnerID)
	header.Add(common.BKHTTPHeaderUser, common.CCSystemCollectorUserName)
	rid := util.GenerateRID()
	header.Add(common.BKHTTPCCRequestID, rid)
	return header, rid
}

const MockMessage = "{\"localTime\": \"2017-09-19 16:57:00\", \"data\": \"{\\\"ip\\\":\\\"192.168.1.7\\\",\\\"bizid\\\":0,\\\"cloudid\\\":0,\\\"data\\\":{\\\"timezone\\\":8,\\\"datetime\\\":\\\"2017-09-19 16:57:07\\\",\\\"utctime\\\":\\\"2017-09-19 08:57:07\\\",\\\"country\\\":\\\"Asia\\\",\\\"city\\\":\\\"Shanghai\\\",\\\"cpu\\\":{\\\"cpuinfo\\\":[{\\\"cpu\\\":0,\\\"vendorID\\\":\\\"GenuineIntel\\\",\\\"family\\\":\\\"6\\\",\\\"model\\\":\\\"63\\\",\\\"stepping\\\":2,\\\"physicalID\\\":\\\"0\\\",\\\"coreID\\\":\\\"0\\\",\\\"cores\\\":1,\\\"modelName\\\":\\\"Intel(R) Xeon(R) CPU E5-26xx v3\\\",\\\"mhz\\\":2294.01,\\\"cacheSize\\\":4096,\\\"flags\\\":[\\\"fpu\\\",\\\"vme\\\",\\\"de\\\",\\\"pse\\\",\\\"tsc\\\",\\\"msr\\\",\\\"pae\\\",\\\"mce\\\",\\\"cx8\\\",\\\"apic\\\",\\\"sep\\\",\\\"mtrr\\\",\\\"pge\\\",\\\"mca\\\",\\\"cmov\\\",\\\"pat\\\",\\\"pse36\\\",\\\"clflush\\\",\\\"mmx\\\",\\\"fxsr\\\",\\\"sse\\\",\\\"sse2\\\",\\\"ss\\\",\\\"ht\\\",\\\"syscall\\\",\\\"nx\\\",\\\"lm\\\",\\\"constant_tsc\\\",\\\"up\\\",\\\"rep_good\\\",\\\"unfair_spinlock\\\",\\\"pni\\\",\\\"pclmulqdq\\\",\\\"ssse3\\\",\\\"fma\\\",\\\"cx16\\\",\\\"pcid\\\",\\\"sse4_1\\\",\\\"sse4_2\\\",\\\"x2apic\\\",\\\"movbe\\\",\\\"popcnt\\\",\\\"tsc_deadline_timer\\\",\\\"aes\\\",\\\"xsave\\\",\\\"avx\\\",\\\"f16c\\\",\\\"rdrand\\\",\\\"hypervisor\\\",\\\"lahf_lm\\\",\\\"abm\\\",\\\"xsaveopt\\\",\\\"bmi1\\\",\\\"avx2\\\",\\\"bmi2\\\"],\\\"microcode\\\":\\\"1\\\"}],\\\"per_usage\\\":[3.0232169701043103],\\\"total_usage\\\":3.0232169701043103,\\\"per_stat\\\":[{\\\"cpu\\\":\\\"cpu0\\\",\\\"user\\\":5206.09,\\\"system\\\":6107.04,\\\"idle\\\":337100.84,\\\"nice\\\":6.68,\\\"iowait\\\":528.24,\\\"irq\\\":0.02,\\\"softirq\\\":13.48,\\\"steal\\\":0,\\\"guest\\\":0,\\\"guestNice\\\":0,\\\"stolen\\\":0}],\\\"total_stat\\\":{\\\"cpu\\\":\\\"cpu-total\\\",\\\"user\\\":5206.09,\\\"system\\\":6107.04,\\\"idle\\\":337100.84,\\\"nice\\\":6.68,\\\"iowait\\\":528.24,\\\"irq\\\":0.02,\\\"softirq\\\":13.48,\\\"steal\\\":0,\\\"guest\\\":0,\\\"guestNice\\\":0,\\\"stolen\\\":0}},\\\"env\\\":{\\\"crontab\\\":[{\\\"user\\\":\\\"root\\\",\\\"content\\\":\\\"#secu-tcs-agent monitor, install at Fri Sep 15 16:12:02 CST 2017\\\\n* * * * * /usr/local/sa/agent/secu-tcs-agent-mon-safe.sh /usr/local/sa/agent \\\\u003e /dev/null 2\\\\u003e\\\\u00261\\\\n*/1 * * * * /usr/local/qcloud/stargate/admin/start.sh \\\\u003e /dev/null 2\\\\u003e\\\\u00261 \\\\u0026\\\\n*/20 * * * * /usr/sbin/ntpdate ntpupdate.tencentyun.com \\\\u003e/dev/null \\\\u0026\\\\n*/1 * * * * cd /usr/local/gse/gseagent; ./cron_agent.sh 1\\\\u003e/dev/null 2\\\\u003e\\\\u00261\\\\n\\\"}],\\\"host\\\":\\\"127.0.0.1  localhost  localhost.localdomain  VM_0_31_centos\\\\n::1         localhost localhost.localdomain localhost6 localhost6.localdomain6\\\\n\\\",\\\"route\\\":\\\"Kernel IP routing table\\\\nDestination     Gateway         Genmask         Flags Metric Ref    Use Iface\\\\n10.0.0.0        0.0.0.0         255.255.255.0   U     0      0        0 eth0\\\\n169.254.0.0     0.0.0.0         255.255.0.0     U     1002   0        0 eth0\\\\n0.0.0.0         10.0.0.1        0.0.0.0         UG    0      0        0 eth0\\\\n\\\"},\\\"disk\\\":{\\\"diskstat\\\":{\\\"vda1\\\":{\\\"major\\\":252,\\\"minor\\\":1,\\\"readCount\\\":24347,\\\"mergedReadCount\\\":570,\\\"writeCount\\\":696357,\\\"mergedWriteCount\\\":4684783,\\\"readBytes\\\":783955968,\\\"writeBytes\\\":22041231360,\\\"readSectors\\\":1531164,\\\"writeSectors\\\":43049280,\\\"readTime\\\":80626,\\\"writeTime\\\":12704736,\\\"iopsInProgress\\\":0,\\\"ioTime\\\":822057,\\\"weightedIoTime\\\":12785026,\\\"name\\\":\\\"vda1\\\",\\\"serialNumber\\\":\\\"\\\",\\\"speedIORead\\\":0,\\\"speedByteRead\\\":0,\\\"speedIOWrite\\\":2.9,\\\"speedByteWrite\\\":171144.53333333333,\\\"util\\\":0.0025666666666666667,\\\"avgrq_sz\\\":115.26436781609195,\\\"avgqu_sz\\\":0.06568333333333334,\\\"await\\\":22.649425287356323,\\\"svctm\\\":0.8850574712643678}},\\\"partition\\\":[{\\\"device\\\":\\\"/dev/vda1\\\",\\\"mountpoint\\\":\\\"/\\\",\\\"fstype\\\":\\\"ext3\\\",\\\"opts\\\":\\\"rw,noatime,acl,user_xattr\\\"}],\\\"usage\\\":[{\\\"path\\\":\\\"/\\\",\\\"fstype\\\":\\\"ext2/ext3\\\",\\\"total\\\":52843638784,\\\"free\\\":47807447040,\\\"used\\\":2351915008,\\\"usedPercent\\\":4.4507060113962345,\\\"inodesTotal\\\":3276800,\\\"inodesUsed\\\":29554,\\\"inodesFree\\\":3247246,\\\"inodesUsedPercent\\\":0.9019165039062501}]},\\\"load\\\":{\\\"load_avg\\\":{\\\"load1\\\":0,\\\"load5\\\":0,\\\"load15\\\":0}},\\\"mem\\\":{\\\"meminfo\\\":{\\\"total\\\":1044832256,\\\"available\\\":805912576,\\\"used\\\":238919680,\\\"usedPercent\\\":22.866797864249705,\\\"free\\\":92041216,\\\"active\\\":521183232,\\\"inactive\\\":352964608,\\\"wired\\\":0,\\\"buffers\\\":110895104,\\\"cached\\\":602976256,\\\"writeback\\\":0,\\\"dirty\\\":151552,\\\"writebacktmp\\\":0},\\\"vmstat\\\":{\\\"total\\\":0,\\\"used\\\":0,\\\"free\\\":0,\\\"usedPercent\\\":0,\\\"sin\\\":0,\\\"sout\\\":0}},\\\"net\\\":{\\\"interface\\\":[{\\\"mtu\\\":65536,\\\"name\\\":\\\"lo\\\",\\\"hardwareaddr\\\":\\\"28:31:52:1d:c6:0a\\\",\\\"flags\\\":[\\\"up\\\",\\\"loopback\\\"],\\\"addrs\\\":[{\\\"addr\\\":\\\"127.0.0.1/8\\\"}]},{\\\"mtu\\\":1500,\\\"name\\\":\\\"eth0\\\",\\\"hardwareaddr\\\":\\\"52:54:00:19:2e:e8\\\",\\\"flags\\\":[\\\"up\\\",\\\"broadcast\\\",\\\"multicast\\\"],\\\"addrs\\\":[{\\\"addr\\\":\\\"127.0.0.1/24\\\"}]}],\\\"dev\\\":[{\\\"name\\\":\\\"lo\\\",\\\"speedSent\\\":0,\\\"speedRecv\\\":0,\\\"speedPacketsSent\\\":0,\\\"speedPacketsRecv\\\":0,\\\"bytesSent\\\":604,\\\"bytesRecv\\\":604,\\\"packetsSent\\\":2,\\\"packetsRecv\\\":2,\\\"errin\\\":0,\\\"errout\\\":0,\\\"dropin\\\":0,\\\"dropout\\\":0,\\\"fifoin\\\":0,\\\"fifoout\\\":0},{\\\"name\\\":\\\"eth0\\\",\\\"speedSent\\\":574,\\\"speedRecv\\\":214,\\\"speedPacketsSent\\\":3,\\\"speedPacketsRecv\\\":2,\\\"bytesSent\\\":161709123,\\\"bytesRecv\\\":285910298,\\\"packetsSent\\\":1116625,\\\"packetsRecv\\\":1167796,\\\"errin\\\":0,\\\"errout\\\":0,\\\"dropin\\\":0,\\\"dropout\\\":0,\\\"fifoin\\\":0,\\\"fifoout\\\":0}],\\\"netstat\\\":{\\\"established\\\":2,\\\"syncSent\\\":1,\\\"synRecv\\\":0,\\\"finWait1\\\":0,\\\"finWait2\\\":0,\\\"timeWait\\\":0,\\\"close\\\":0,\\\"closeWait\\\":0,\\\"lastAck\\\":0,\\\"listen\\\":2,\\\"closing\\\":0},\\\"protocolstat\\\":[{\\\"protocol\\\":\\\"udp\\\",\\\"stats\\\":{\\\"inDatagrams\\\":176253,\\\"inErrors\\\":0,\\\"noPorts\\\":1,\\\"outDatagrams\\\":199569,\\\"rcvbufErrors\\\":0,\\\"sndbufErrors\\\":0}}]},\\\"system\\\":{\\\"info\\\":{\\\"hostname\\\":\\\"VM_0_31_centos\\\",\\\"uptime\\\":348315,\\\"bootTime\\\":1505463112,\\\"procs\\\":142,\\\"os\\\":\\\"linux\\\",\\\"platform\\\":\\\"centos\\\",\\\"platformFamily\\\":\\\"rhel\\\",\\\"platformVersion\\\":\\\"6.2\\\",\\\"kernelVersion\\\":\\\"2.6.32-504.30.3.el6.x86_64\\\",\\\"virtualizationSystem\\\":\\\"\\\",\\\"virtualizationRole\\\":\\\"\\\",\\\"hostid\\\":\\\"96D0F4CA-2157-40E6-BF22-6A7CD9B6EB8C\\\",\\\"systemtype\\\":\\\"64-bit\\\"}}}}\", \"timestamp\": 1505811427, \"dtEventTime\": \"2017-09-19 16:57:07\", \"dtEventTimeStamp\": 1505811427000}"
