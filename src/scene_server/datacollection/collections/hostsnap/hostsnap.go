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

	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"

	"github.com/tidwall/gjson"
	"gopkg.in/redis.v5"
)

type HostSnap struct {
	redisCli    *redis.Client
	authManager *extensions.AuthManager
	*backbone.Engine

	filter *filter
	ctx    context.Context
	db     dal.RDB
}

func NewHostSnap(ctx context.Context, redisCli *redis.Client, db dal.RDB, engine *backbone.Engine, authManager *extensions.AuthManager) *HostSnap {
	h := &HostSnap{
		redisCli:    redisCli,
		ctx:         ctx,
		db:          db,
		authManager: authManager,
		Engine:      engine,
		filter:      newFilter(),
	}
	return h
}

var compareFields = []string{"bk_cpu", "bk_cpu_module", "bk_cpu_mhz", "bk_disk", "bk_mem", "bk_os_type", "bk_os_name",
	"bk_os_version", "bk_host_name", "bk_outer_mac", "bk_mac", "bk_os_bit",
	common.HostFieldDockerClientVersion, common.HostFieldDockerServerVersion}

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

	key := common.RedisSnapKeyPrefix + strconv.FormatInt(hostID, 10)
	if err := h.redisCli.Set(key, data, time.Minute*10).Err(); err != nil {
		blog.Errorf("save snapshot key: %s to redis failed: %v, rid: %s", key, err, rid)
	}

	setter, raw := parseSetter(&val, innerIP, outerIP)
	// no need to update
	if !needToUpdate(raw, host) {
		return nil
	}

	blog.V(5).Infof("snapshot for host changed, need update, host id: %d, ip: %s, cloud id: %d, from %s to %s, rid: %s",
		hostID, innerIP, cloudID, host, raw, rid)

	// add auditLog
	preData, err := h.CoreAPI.CoreService().Host().GetHostByID(h.ctx, header, hostID)
	if err != nil {
		blog.Errorf("snapshot get host previous data failed, err: %s, hostID: %d, rid: %s", err, hostID, rid)
		return err
	}
	if !preData.Result {
		blog.Errorf("snapshot get host previous data failed, code: %d, err: %s, hostID: %d, rid: %s",
			preData.Code, preData.ErrMsg, hostID, rid)
		return err
	}

	opt := &metadata.UpdateOption{
		Condition: map[string]interface{}{
			common.BKHostIDField: hostID,
		},
		Data:       mapstr.NewFromMap(setter),
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

	curData := make(map[string]interface{}, 0)
	for k, v := range preData.Data {
		if value, exist := setter[k]; exist {
			// set with the new value
			curData[k] = value
			continue
		}
		curData[k] = v
	}

	input := &metadata.HostModuleRelationRequest{HostIDArr: []int64{hostID}, Fields: []string{common.BKAppIDField}}
	moduleHost, err := h.CoreAPI.CoreService().Host().GetHostModuleRelation(h.ctx, header, input)
	if err != nil {
		blog.Errorf("snapshot get host: %d/%s module relation failed, err:%v, rid: %s", hostID, innerIP, err, rid)
		return err
	}
	if !moduleHost.Result {
		blog.Errorf("snapshot get host: %d%s module relation failed, err: %v, rid: %s", hostID, innerIP, moduleHost.ErrMsg, rid)
		return fmt.Errorf("snapshot get moduleHostConfig failed, fail to create auditLog")
	}

	audit := auditlog.NewAudit(h.CoreAPI, header)
	properties, err := audit.GetAuditLogProperty(h.ctx, common.BKInnerObjIDHost)
	if err != nil {
		return err
	}
	var bizID int64
	if len(moduleHost.Data.Info) > 0 {
		bizID = moduleHost.Data.Info[0].AppID
	}
	bizName := ""
	if bizID > 0 {
		bizName, err = audit.GetInstNameByID(h.ctx, common.BKInnerObjIDApp, bizID)
		if err != nil {
			return err
		}
	}
	auditLog := metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       metadata.AuditUpdate,
		OperateFrom:  metadata.FromDataCollection,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				BusinessID:   bizID,
				BusinessName: bizName,
				ResourceID:   hostID,
				ResourceName: innerIP,
				Details: &metadata.BasicContent{
					PreData:    preData.Data,
					CurData:    curData,
					Properties: properties,
				},
			},
			ModelID: common.BKInnerObjIDHost,
		},
	}
	result, err := h.CoreAPI.CoreService().Audit().SaveAuditLog(h.ctx, header, auditLog)
	if err != nil {
		blog.Errorf("snapshot create host %d/%s audit log failed, err: %v, rid: %s", hostID, innerIP, err.Error())
		return err
	}
	if !result.Result {
		blog.Errorf("snapshot create host %d/%s audit log failed, err: %s, rid: %s", hostID, innerIP, result.ErrMsg, rid)
		return fmt.Errorf("create host audit log failed, err: %s", result.ErrMsg)
	}

	blog.V(5).Infof("snapshot for host changed, update success, host id: %d, ip: %s, cloud id: %d, rid: %s",
		hostID, innerIP, cloudID, rid)

	return nil
}

func needToUpdate(src, toCompare string) bool {
	srcElements := gjson.GetMany(src, compareFields...)
	compareElements := gjson.GetMany(toCompare, compareFields...)
	for idx := range compareFields {
		// compare these value with string directly to avoid empty value or null value.
		if srcElements[idx].String() != compareElements[idx].String() {
			// tolerate bk_cpu_mhz changes less than 100
			if compareFields[idx] == "bk_cpu_mhz" {
				diff := srcElements[idx].Int() - compareElements[idx].Int()
				if -100 < diff && diff < 100 {
					continue
				}
			}
			return true
		}
	}
	return false
}

func parseSetter(val *gjson.Result, innerIP, outerIP string) (map[string]interface{}, string) {
	var cpumodule = val.Get("data.cpu.cpuinfo.0.modelName").String()
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
	var ostype = val.Get("data.system.info.os").String()
	var osname string
	platform := val.Get("data.system.info.platform").String()
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
	var OuterMAC, InnerMAC string
	for _, inter := range val.Get("data.net.interface").Array() {
		for _, addr := range inter.Get("addrs.#.addr").Array() {
			ip := strings.Split(addr.String(), "/")[0]
			if ip == innerIP {
				InnerMAC = inter.Get("hardwareaddr").String()
			} else if ip == outerIP {
				OuterMAC = inter.Get("hardwareaddr").String()
			}
		}
	}

	osbit := val.Get("data.system.info.systemtype").String()

	dockerClientVersion := val.Get("data.system.docker.Client.Version").String()
	dockerServerVersion := val.Get("data.system.docker.Server.Version").String()

	mem = mem >> 10 >> 10
	setter := map[string]interface{}{
		"bk_cpu":                            cpunum,
		"bk_cpu_module":                     cpumodule,
		"bk_cpu_mhz":                        CPUMhz,
		"bk_disk":                           disk,
		"bk_mem":                            mem,
		"bk_os_type":                        ostype,
		"bk_os_name":                        osname,
		"bk_os_version":                     version,
		"bk_host_name":                      hostname,
		"bk_outer_mac":                      OuterMAC,
		"bk_mac":                            InnerMAC,
		"bk_os_bit":                         osbit,
		common.HostFieldDockerClientVersion: dockerClientVersion,
		common.HostFieldDockerServerVersion: dockerServerVersion,
	}

	raw := strings.Builder{}
	raw.WriteByte('{')
	raw.WriteString("\"bk_cpu\":")
	raw.WriteString(strconv.FormatInt(cpunum, 10))
	raw.WriteString(",")
	raw.WriteString("\"bk_cpu_module\":")
	raw.Write([]byte("\"" + cpumodule + "\""))
	raw.WriteString(",")
	raw.WriteString("\"bk_cpu_mhz\":")
	raw.WriteString(strconv.FormatInt(CPUMhz, 10))
	raw.WriteString(",")
	raw.WriteString("\"bk_disk\":")
	raw.WriteString(strconv.FormatUint(disk, 10))
	raw.WriteString(",")
	raw.WriteString("\"bk_mem\":")
	raw.WriteString(strconv.FormatUint(mem, 10))
	raw.WriteString(",")
	raw.WriteString("\"bk_os_type\":")
	raw.Write([]byte("\"" + ostype + "\""))
	raw.WriteString(",")
	raw.WriteString("\"bk_os_name\":")
	raw.Write([]byte("\"" + osname + "\""))
	raw.WriteString(",")
	raw.WriteString("\"bk_os_version\":")
	raw.Write([]byte("\"" + version + "\""))
	raw.WriteString(",")
	raw.WriteString("\"bk_host_name\":")
	raw.Write([]byte("\"" + hostname + "\""))
	raw.WriteString(",")
	raw.WriteString("\"bk_outer_mac\":")
	raw.Write([]byte("\"" + OuterMAC + "\""))
	raw.WriteString(",")
	raw.WriteString("\"bk_mac\":")
	raw.Write([]byte("\"" + InnerMAC + "\""))
	raw.WriteString(",")
	raw.WriteString("\"bk_os_bit\":")
	raw.Write([]byte("\"" + osbit + "\""))
	raw.WriteString(",")
	raw.WriteString("\"docker_client_version\":")
	raw.Write([]byte("\"" + dockerClientVersion + "\""))
	raw.WriteString(",")
	raw.WriteString("\"docker_server_version\":")
	raw.Write([]byte("\"" + dockerServerVersion + "\""))
	raw.WriteByte('}')

	if cpunum <= 0 {
		blog.V(4).Infof("bk_cpu not found in message for %s", innerIP)
	}
	if cpumodule == "" {
		blog.V(4).Infof("bk_cpu_module not found in message for %s", innerIP)
	}
	if CPUMhz <= 0 {
		blog.V(4).Infof("bk_cpu_mhz not found in message for %s", innerIP)
	}
	if disk <= 0 {
		blog.V(4).Infof("bk_disk not found in message for %s", innerIP)
	}
	if mem <= 0 {
		blog.V(4).Infof("bk_mem not found in message for %s", innerIP)
	}
	if ostype == "" {
		blog.V(4).Infof("bk_os_type not found in message for %s", innerIP)
	}
	if osname == "" {
		blog.V(4).Infof("bk_os_name not found in message for %s", innerIP)
	}
	if version == "" {
		blog.V(4).Infof("bk_os_version not found in message for %s", innerIP)
	}
	if hostname == "" {
		blog.V(4).Infof("bk_host_name not found in message for %s", innerIP)
	}
	if outerIP != "" && OuterMAC == "" {
		blog.V(4).Infof("bk_outer_mac not found in message for %s", innerIP)
	}
	if InnerMAC == "" {
		blog.V(4).Infof("bk_mac not found in message for %s", innerIP)
	}

	return setter, raw.String()
}

var reqireFields = append(compareFields, "bk_host_id", "bk_host_innerip", "bk_host_outerip")

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

		host, err := h.Engine.CoreAPI.CoreService().Cache().SearchHostWithInnerIP(context.Background(), header, opt)
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

func newHeaderWithRid() (http.Header, string) {
	header := http.Header{}
	header.Add(common.BKHTTPOwnerID, common.BKDefaultOwnerID)
	header.Add(common.BKHTTPHeaderUser, common.CCSystemCollectorUserName)
	rid := util.GenerateRID()
	header.Add(common.BKHTTPCCRequestID, rid)
	return header, rid
}

const MockMessage = "{\"localTime\": \"2017-09-19 16:57:00\", \"data\": \"{\\\"ip\\\":\\\"192.168.1.7\\\",\\\"bizid\\\":0,\\\"cloudid\\\":0,\\\"data\\\":{\\\"timezone\\\":8,\\\"datetime\\\":\\\"2017-09-19 16:57:07\\\",\\\"utctime\\\":\\\"2017-09-19 08:57:07\\\",\\\"country\\\":\\\"Asia\\\",\\\"city\\\":\\\"Shanghai\\\",\\\"cpu\\\":{\\\"cpuinfo\\\":[{\\\"cpu\\\":0,\\\"vendorID\\\":\\\"GenuineIntel\\\",\\\"family\\\":\\\"6\\\",\\\"model\\\":\\\"63\\\",\\\"stepping\\\":2,\\\"physicalID\\\":\\\"0\\\",\\\"coreID\\\":\\\"0\\\",\\\"cores\\\":1,\\\"modelName\\\":\\\"Intel(R) Xeon(R) CPU E5-26xx v3\\\",\\\"mhz\\\":2294.01,\\\"cacheSize\\\":4096,\\\"flags\\\":[\\\"fpu\\\",\\\"vme\\\",\\\"de\\\",\\\"pse\\\",\\\"tsc\\\",\\\"msr\\\",\\\"pae\\\",\\\"mce\\\",\\\"cx8\\\",\\\"apic\\\",\\\"sep\\\",\\\"mtrr\\\",\\\"pge\\\",\\\"mca\\\",\\\"cmov\\\",\\\"pat\\\",\\\"pse36\\\",\\\"clflush\\\",\\\"mmx\\\",\\\"fxsr\\\",\\\"sse\\\",\\\"sse2\\\",\\\"ss\\\",\\\"ht\\\",\\\"syscall\\\",\\\"nx\\\",\\\"lm\\\",\\\"constant_tsc\\\",\\\"up\\\",\\\"rep_good\\\",\\\"unfair_spinlock\\\",\\\"pni\\\",\\\"pclmulqdq\\\",\\\"ssse3\\\",\\\"fma\\\",\\\"cx16\\\",\\\"pcid\\\",\\\"sse4_1\\\",\\\"sse4_2\\\",\\\"x2apic\\\",\\\"movbe\\\",\\\"popcnt\\\",\\\"tsc_deadline_timer\\\",\\\"aes\\\",\\\"xsave\\\",\\\"avx\\\",\\\"f16c\\\",\\\"rdrand\\\",\\\"hypervisor\\\",\\\"lahf_lm\\\",\\\"abm\\\",\\\"xsaveopt\\\",\\\"bmi1\\\",\\\"avx2\\\",\\\"bmi2\\\"],\\\"microcode\\\":\\\"1\\\"}],\\\"per_usage\\\":[3.0232169701043103],\\\"total_usage\\\":3.0232169701043103,\\\"per_stat\\\":[{\\\"cpu\\\":\\\"cpu0\\\",\\\"user\\\":5206.09,\\\"system\\\":6107.04,\\\"idle\\\":337100.84,\\\"nice\\\":6.68,\\\"iowait\\\":528.24,\\\"irq\\\":0.02,\\\"softirq\\\":13.48,\\\"steal\\\":0,\\\"guest\\\":0,\\\"guestNice\\\":0,\\\"stolen\\\":0}],\\\"total_stat\\\":{\\\"cpu\\\":\\\"cpu-total\\\",\\\"user\\\":5206.09,\\\"system\\\":6107.04,\\\"idle\\\":337100.84,\\\"nice\\\":6.68,\\\"iowait\\\":528.24,\\\"irq\\\":0.02,\\\"softirq\\\":13.48,\\\"steal\\\":0,\\\"guest\\\":0,\\\"guestNice\\\":0,\\\"stolen\\\":0}},\\\"env\\\":{\\\"crontab\\\":[{\\\"user\\\":\\\"root\\\",\\\"content\\\":\\\"#secu-tcs-agent monitor, install at Fri Sep 15 16:12:02 CST 2017\\\\n* * * * * /usr/local/sa/agent/secu-tcs-agent-mon-safe.sh /usr/local/sa/agent \\\\u003e /dev/null 2\\\\u003e\\\\u00261\\\\n*/1 * * * * /usr/local/qcloud/stargate/admin/start.sh \\\\u003e /dev/null 2\\\\u003e\\\\u00261 \\\\u0026\\\\n*/20 * * * * /usr/sbin/ntpdate ntpupdate.tencentyun.com \\\\u003e/dev/null \\\\u0026\\\\n*/1 * * * * cd /usr/local/gse/gseagent; ./cron_agent.sh 1\\\\u003e/dev/null 2\\\\u003e\\\\u00261\\\\n\\\"}],\\\"host\\\":\\\"127.0.0.1  localhost  localhost.localdomain  VM_0_31_centos\\\\n::1         localhost localhost.localdomain localhost6 localhost6.localdomain6\\\\n\\\",\\\"route\\\":\\\"Kernel IP routing table\\\\nDestination     Gateway         Genmask         Flags Metric Ref    Use Iface\\\\n10.0.0.0        0.0.0.0         255.255.255.0   U     0      0        0 eth0\\\\n169.254.0.0     0.0.0.0         255.255.0.0     U     1002   0        0 eth0\\\\n0.0.0.0         10.0.0.1        0.0.0.0         UG    0      0        0 eth0\\\\n\\\"},\\\"disk\\\":{\\\"diskstat\\\":{\\\"vda1\\\":{\\\"major\\\":252,\\\"minor\\\":1,\\\"readCount\\\":24347,\\\"mergedReadCount\\\":570,\\\"writeCount\\\":696357,\\\"mergedWriteCount\\\":4684783,\\\"readBytes\\\":783955968,\\\"writeBytes\\\":22041231360,\\\"readSectors\\\":1531164,\\\"writeSectors\\\":43049280,\\\"readTime\\\":80626,\\\"writeTime\\\":12704736,\\\"iopsInProgress\\\":0,\\\"ioTime\\\":822057,\\\"weightedIoTime\\\":12785026,\\\"name\\\":\\\"vda1\\\",\\\"serialNumber\\\":\\\"\\\",\\\"speedIORead\\\":0,\\\"speedByteRead\\\":0,\\\"speedIOWrite\\\":2.9,\\\"speedByteWrite\\\":171144.53333333333,\\\"util\\\":0.0025666666666666667,\\\"avgrq_sz\\\":115.26436781609195,\\\"avgqu_sz\\\":0.06568333333333334,\\\"await\\\":22.649425287356323,\\\"svctm\\\":0.8850574712643678}},\\\"partition\\\":[{\\\"device\\\":\\\"/dev/vda1\\\",\\\"mountpoint\\\":\\\"/\\\",\\\"fstype\\\":\\\"ext3\\\",\\\"opts\\\":\\\"rw,noatime,acl,user_xattr\\\"}],\\\"usage\\\":[{\\\"path\\\":\\\"/\\\",\\\"fstype\\\":\\\"ext2/ext3\\\",\\\"total\\\":52843638784,\\\"free\\\":47807447040,\\\"used\\\":2351915008,\\\"usedPercent\\\":4.4507060113962345,\\\"inodesTotal\\\":3276800,\\\"inodesUsed\\\":29554,\\\"inodesFree\\\":3247246,\\\"inodesUsedPercent\\\":0.9019165039062501}]},\\\"load\\\":{\\\"load_avg\\\":{\\\"load1\\\":0,\\\"load5\\\":0,\\\"load15\\\":0}},\\\"mem\\\":{\\\"meminfo\\\":{\\\"total\\\":1044832256,\\\"available\\\":805912576,\\\"used\\\":238919680,\\\"usedPercent\\\":22.866797864249705,\\\"free\\\":92041216,\\\"active\\\":521183232,\\\"inactive\\\":352964608,\\\"wired\\\":0,\\\"buffers\\\":110895104,\\\"cached\\\":602976256,\\\"writeback\\\":0,\\\"dirty\\\":151552,\\\"writebacktmp\\\":0},\\\"vmstat\\\":{\\\"total\\\":0,\\\"used\\\":0,\\\"free\\\":0,\\\"usedPercent\\\":0,\\\"sin\\\":0,\\\"sout\\\":0}},\\\"net\\\":{\\\"interface\\\":[{\\\"mtu\\\":65536,\\\"name\\\":\\\"lo\\\",\\\"hardwareaddr\\\":\\\"28:31:52:1d:c6:0a\\\",\\\"flags\\\":[\\\"up\\\",\\\"loopback\\\"],\\\"addrs\\\":[{\\\"addr\\\":\\\"127.0.0.1/8\\\"}]},{\\\"mtu\\\":1500,\\\"name\\\":\\\"eth0\\\",\\\"hardwareaddr\\\":\\\"52:54:00:19:2e:e8\\\",\\\"flags\\\":[\\\"up\\\",\\\"broadcast\\\",\\\"multicast\\\"],\\\"addrs\\\":[{\\\"addr\\\":\\\"127.0.0.1/24\\\"}]}],\\\"dev\\\":[{\\\"name\\\":\\\"lo\\\",\\\"speedSent\\\":0,\\\"speedRecv\\\":0,\\\"speedPacketsSent\\\":0,\\\"speedPacketsRecv\\\":0,\\\"bytesSent\\\":604,\\\"bytesRecv\\\":604,\\\"packetsSent\\\":2,\\\"packetsRecv\\\":2,\\\"errin\\\":0,\\\"errout\\\":0,\\\"dropin\\\":0,\\\"dropout\\\":0,\\\"fifoin\\\":0,\\\"fifoout\\\":0},{\\\"name\\\":\\\"eth0\\\",\\\"speedSent\\\":574,\\\"speedRecv\\\":214,\\\"speedPacketsSent\\\":3,\\\"speedPacketsRecv\\\":2,\\\"bytesSent\\\":161709123,\\\"bytesRecv\\\":285910298,\\\"packetsSent\\\":1116625,\\\"packetsRecv\\\":1167796,\\\"errin\\\":0,\\\"errout\\\":0,\\\"dropin\\\":0,\\\"dropout\\\":0,\\\"fifoin\\\":0,\\\"fifoout\\\":0}],\\\"netstat\\\":{\\\"established\\\":2,\\\"syncSent\\\":1,\\\"synRecv\\\":0,\\\"finWait1\\\":0,\\\"finWait2\\\":0,\\\"timeWait\\\":0,\\\"close\\\":0,\\\"closeWait\\\":0,\\\"lastAck\\\":0,\\\"listen\\\":2,\\\"closing\\\":0},\\\"protocolstat\\\":[{\\\"protocol\\\":\\\"udp\\\",\\\"stats\\\":{\\\"inDatagrams\\\":176253,\\\"inErrors\\\":0,\\\"noPorts\\\":1,\\\"outDatagrams\\\":199569,\\\"rcvbufErrors\\\":0,\\\"sndbufErrors\\\":0}}]},\\\"system\\\":{\\\"info\\\":{\\\"hostname\\\":\\\"VM_0_31_centos\\\",\\\"uptime\\\":348315,\\\"bootTime\\\":1505463112,\\\"procs\\\":142,\\\"os\\\":\\\"linux\\\",\\\"platform\\\":\\\"centos\\\",\\\"platformFamily\\\":\\\"rhel\\\",\\\"platformVersion\\\":\\\"6.2\\\",\\\"kernelVersion\\\":\\\"2.6.32-504.30.3.el6.x86_64\\\",\\\"virtualizationSystem\\\":\\\"\\\",\\\"virtualizationRole\\\":\\\"\\\",\\\"hostid\\\":\\\"96D0F4CA-2157-40E6-BF22-6A7CD9B6EB8C\\\",\\\"systemtype\\\":\\\"64-bit\\\"}}}}\", \"timestamp\": 1505811427, \"dtEventTime\": \"2017-09-19 16:57:07\", \"dtEventTimeStamp\": 1505811427000}"
