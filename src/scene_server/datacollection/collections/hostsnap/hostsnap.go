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

// Package hostsnap TODO
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
	httpheader "configcenter/src/common/http/header"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/redis"

	goredis "github.com/go-redis/redis/v7"
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
	// redisConsumptionCheckPrefix redis consumption check prefix
	redisConsumptionCheckPrefix = "consumptionCheck:"
)

var (
	// todo: 这里有一个问题，目前在动态IP场景下是会将上报的ip地址都放到innerIP中，后续可能涉及到调整
	// Note:Among them, ipv4 and ipv6 addresses involve updating in dynamic scenarios, but are not allowed to be updated
	// in static ip scenarios, and require special processing
	compareFields = []string{"bk_cpu", "bk_cpu_module", "bk_disk", "bk_mem", "bk_os_type", "bk_os_name",
		"bk_os_version", "bk_host_name", "bk_outer_mac", "bk_mac", "bk_os_bit", "bk_cpu_architecture",
		common.BKHostInnerIPField, common.BKHostInnerIPv6Field}
	reqireFields = append(compareFields, common.BKHostIDField, common.BKAddressingField, common.BKHostOuterIPField,
		common.BKHostOuterIPv6Field)

	// notice: 为了对应不同版本和环境差异，再当前版本中设置compareFields中不参加对比的字段
	ignoreCompareField = make(map[string]struct{}, 0)
)

// HostSnap TODO
type HostSnap struct {
	redisCli    redis.Client
	authManager *extensions.AuthManager
	*backbone.Engine
	rateLimit flowctrl.RateLimiter
	filter    *filter
	ctx       context.Context
	window    *Window
}

// NewHostSnap new hostsnap
func NewHostSnap(ctx context.Context, redisCli redis.Client, engine *backbone.Engine,
	authManager *extensions.AuthManager) *HostSnap {
	qps, burst := getRateLimiterConfig()
	h := &HostSnap{
		redisCli:    redisCli,
		ctx:         ctx,
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
		blog.Errorf("can't find the value of datacollection.hostsnap.rateLimiter.qps settings, "+
			"set the default value: %d", defaultRateLimiterQPS)
		qps = defaultRateLimiterQPS
	}
	burst, err := cc.Int("datacollection.hostsnap.rateLimiter.burst")
	if err != nil {
		blog.Errorf("can't find the value of datacollection.hostsnap.rateLimiter.burst setting, "+
			"set the default value: %d", defaultRateLimiterBurst)
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

func (h *HostSnap) putDataIntoDelayQueue(rid, msg string) error {
	blog.V(5).Infof("put msg to delay queue, msg: %s, rid: %s", msg, rid)

	// There is a question here, is the member's setting agentID or msg:
	// 1. If msg is set as a member, there are two problems:
	// a. What if a new msg corresponding to the same agentID comes up at this time? Because the queue is old at this
	// time, the timestamp will be judged later in the processing. no problem.
	// b. If another msg comes up at this time, and the corresponding host is not found in the database for this msg, it
	// will still be added to the delay queue at this time. At this time, the timestamp will be compared when it is
	// retrieved later. That is, it will increase storage and subsequent redundant processing. According to the current
	// data reported at most one minute. msg won't be too many. From the point of view of data processing, it is
	// protected by a timestamp behind it. It won't be a problem either.
	// 2. What if agentID is used as a member?
	// a. At this time, two data structures are needed to store agentID in zset, and messages need to be stored in
	// another key-value structure. One message needs to involve two data structures, and there will also be some
	// redundant space occupation at this time. .
	// b. At this time, a delay queue involves two data structures. To increase the complexity of the operation, it is
	// necessary to ensure the transactional nature of the operation. When adding, you need to operate on two data
	// structures, and when deleting, you also need to operate on two data structures. At present, the benefits are not
	// large.
	body := &goredis.Z{
		Score:  float64(time.Now().Unix()),
		Member: msg,
	}

	if err := h.redisCli.ZAdd(context.Background(), common.RedisHostSnapMsgDelayQueue, body).Err(); err != nil {
		return err
	}

	return nil
}

// getBaseInfoFromCollectorsMsg obtain basic information from the reported msg: agentID, ipv4, ipv6, cloudID, and the
// entire parsed message body.
func getBaseInfoFromCollectorsMsg(msg *string, rid string) (string, []string, []string, int64, gjson.Result, error) {

	var data string
	if !gjson.Get(*msg, "cloudid").Exists() {
		data = gjson.Get(*msg, "data").String()
	} else {
		data = *msg
	}
	val := gjson.Parse(data)
	agentID := gjson.Get(*msg, "bk_agent_id").String()
	cloudID := val.Get("cloudid").Int()

	ipv4, ipv6 := getIPsFromMsg(&val, agentID, rid)
	if len(ipv4) == 0 && len(ipv6) == 0 {
		return "", nil, nil, 0, gjson.Result{}, errors.New("msg has no ipv4 and ipv6 address")
	}
	return agentID, ipv4, ipv6, cloudID, val, nil
}

func checkMsgInfoValid(rid, agentID, host string, elements []gjson.Result) error {
	// check host id field
	if !elements[0].Exists() {
		blog.Errorf("snapshot analyze, but host id not exist, host: %s, rid: %s", host, rid)
		return errors.New("host id not exist")
	}

	hostID := elements[0].Int()
	if hostID == 0 {
		blog.Errorf("snapshot analyze, but host id is 0, host: %s, rid: %s", host, rid)
		return errors.New("host id can not be 0")
	}

	// check inner ip
	if !elements[1].Exists() {
		blog.Errorf("snapshot analyze, but host inner ip not exist, host: %s, rid: %s", host, rid)
		return errors.New("host inner ip not exist")
	}

	// there must be an addressing field, the default is static.
	if !elements[3].Exists() {
		blog.Errorf("snapshot analyze, but host addressing not exist, host: %s, rid: %s", host, rid)
		return errors.New("host addressing not exist")
	}

	// If the data is obtained through ip+cloudID, it means that there is no agentID in the data reported at this time,
	// and if there is an agentID in the data stored in the database, there is a problem, so it should be discarded.
	if agentID == "" && elements[4].Exists() && elements[4].String() != "" {
		blog.Errorf("snapshot analyze, agentID is inconsistent, host: %s, agentID: %s, rid: %s",
			host, elements[4].String(), rid)
		return errors.New("the agentID field is inconsistent with the data")
	}

	// if agentID is not empty, the bk_addressing field must be present.
	if agentID != "" && !elements[3].Exists() {
		blog.Errorf("snapshot analyze, bk_addressing is null, host: %s, agentID: %s rid: %s", host, agentID, rid)
		return errors.New("bk_addressing is null")
	}

	return nil
}

func (h *HostSnap) getHostDetail(header http.Header, rid, agentID, msg, sourceType string, ips []string,
	cloudID int64) (host string, err error) {

	if agentID != "" {
		host, err = h.getHostByAgentID(header, rid, agentID)
		if err != nil {
			// todo 由于采集器会上报没有绑定agent id的主机信息，给db造成了压力，这里先把加入延迟队列的逻辑去掉
			// if err := h.putDataIntoDelayQueue(rid, msg); err != nil {
			//	blog.Errorf("put msg to delay queue failed, agentID: %, err: %v, rid: %s", agentID, err, rid)
			// }
			blog.Errorf("get host detail with agentID: %v failed, err: %v, rid: %s", agentID, err, rid)
			return "", errors.New("no host founded")
		}
		if sourceType == metadata.HostSnapDataSourcesDelayQueue {
			// If the data is obtained from the delay queue, then it means that the host information that could not be
			// found before can be found now, and the data needs to be deleted from the delay queue.
			err := h.redisCli.ZRem(context.Background(), common.RedisHostSnapMsgDelayQueue, msg).Err()
			if err != nil {
				blog.Errorf("remove member failed, msg: %v, err: %v, rid: %s", msg, err, rid)
				return "", nil
			}
		}
	} else {
		// when querying according to ip+cloudID, it is in a static IP scenario, so here is also a host details.
		host, err = h.getHostByVal(header, rid, cloudID, ips)
		if err != nil {
			blog.Errorf("get host detail with ips: %v failed, err: %v, rid: %s", ips, err, rid)
			return "", err
		}
	}
	return host, nil
}

// 1. Parse out the fields such as agentID, ip, cloudID, etc. from the reported message. The ipv4 field is a required
// field.
// 2. If the agentID is not empty, you need to use the agentID as the key to query the host. If the agentID is empty or
// there is no agentID field, you need to query the information in the static IP scenario through ip+cloudID. It should
// be noted here that if both gse agent and cmdb support dynamic IP scenarios, then bkmonitorBeat At this time, the old
// version will cause an update error in the dynamic ip scenario. It is necessary to ensure that bkMonitorBeat also
// needs to support agentID.
// 3. If the host information is not found in cmdb at this time, put this information in the delay queue, and the delay
// queue will re-report the query failure message every 10s. The information in the delay queue is kept for a maximum
// of ten minutes, and the messages before ten minutes are directly discarded.
// 4. After the obtained message, it will be judged whether an updated message has been updated to cmdb. If the
// timestamp of this message is earlier than the data in cmdb, it will be discarded directly.

// Analyze analyze host snap
func (h *HostSnap) Analyze(msg *string, sourceType string) (bool, error) {
	if msg == nil {
		return false, errors.New("message nil")
	}

	header, rid := newHeaderWithRid()

	agentID, ipv4, ipv6, cloudID, val, err := getBaseInfoFromCollectorsMsg(msg, rid)
	if err != nil {
		blog.Errorf("parse base info failed, msg: %s, err: %v, rid: %s", *msg, err, rid)
		return false, err
	}
	host, err := h.getHostDetail(header, rid, agentID, *msg, sourceType, ipv4, cloudID)
	if err != nil {
		blog.Errorf("get host detail failed, agentID: %s, ips: %v, err: %v, rid: %s", agentID, ipv4, err, rid)
		return false, err
	}

	if host == "" {
		blog.Errorf("get host detail failed, agentID: %s, ips: %v, err: %v, rid: %s", agentID, ipv4, err, rid)
		return false, errors.New("get host detail failed")
	}

	fields := []string{common.BKHostIDField, common.BKHostInnerIPField, common.BKHostOuterIPField,
		common.BKAddressingField, common.BKAgentIDField, common.BKHostInnerIPv6Field, common.BKHostOuterIPv6Field}
	elements := gjson.GetMany(host, fields...)

	if err := checkMsgInfoValid(rid, agentID, host, elements); err != nil {
		return false, err
	}
	hostID := elements[0].Int()
	innerIP := elements[1].String()

	// save host snapshot in redis
	if !val.Get("data.apiVer").Exists() {
		h.saveHostsnap(header, &val, hostID)
	}

	// window restriction on request when no apiVer information reported
	if !val.Get("data.apiVer").Exists() && !h.window.canPassWindow() {
		if blog.V(4) {
			blog.Infof("not within the time window that can pass, skip host snapshot data update, "+
				"host id: %d, ip: %s, cloud id: %d, rid: %s", hostID, innerIP, cloudID, rid)
		}
		return false, nil
	}

	if h.skipMsg(val, innerIP, rid, hostID, cloudID) {
		return false, nil
	}

	updateIPv4, updateIPv6, err := getIPv4AndIPv6UpdateData(elements[3].String(), ipv4, ipv6)
	if err != nil {
		blog.Errorf("get host ipv4 and ipv6 update data failed, agentID: %s, ipv4: %v, ipv6: %v, err: %v, rid: %s",
			agentID, ipv4, ipv6, err, rid)
		return false, err
	}

	outerIP := elements[2].String()

	setter, raw := make(map[string]interface{}), ""

	if val.Get("data.apiVer").String() >= "v1.0" {
		setter, raw = parseV10Setter(&val, &hostInfo{innerIPv4: innerIP, outerIPv4: outerIP,
			innerIPv6: elements[5].String(), outerIPv6: elements[6].String(), addressing: elements[3].String(),
			updateIPv4: updateIPv4, updateIPv6: updateIPv6})
	} else {
		setter, raw = parseSetter(&val, innerIP, outerIP)
	}

	// no need to update
	if !needToUpdate(raw, host, elements[3].String()) {
		return false, nil
	}

	// limit request from the old collection plug-in
	if !val.Get("data.apiVer").Exists() && !h.rateLimit.TryAccept() {
		blog.Warnf("skip host snapshot data update due to request limit, host id: %d, ip: %s, cloud id: %d, "+
			"rid: %s", hostID, innerIP, cloudID, rid)
		return false, nil
	}

	// limit the request from the new collection plug-in
	if val.Get("data.apiVer").Exists() {
		h.rateLimit.Accept()
	}

	blog.V(5).Infof("snapshot for host changed, need update, host id: %d, ip: %s, cloud id: %d, from %s "+
		"to %s, rid: %s", hostID, innerIP, cloudID, host, raw, rid)

	hostOption := updateHostOption{setter: setter, host: host, innerIP: innerIP, hostID: hostID, cloudID: cloudID}

	txnErr := h.updateHostWithColletorMsg(header, rid, hostOption)
	if txnErr != nil {
		return true, txnErr
	}
	return false, nil
}

func getIPv4AndIPv6UpdateData(addressing string, ipv4 []string, ipv6 []string) ([]string, []string, error) {
	var err error
	var updateIPv4, updateIPv6 []string
	// in the dynamic ip scenario, innerIP needs to be updated
	if addressing == common.BKAddressingDynamic {
		updateIPv4 = ipv4
		updateIPv6, err = common.ConvertHostIpv6Val(ipv6)
		if err != nil {
			return nil, nil, err
		}
	}

	return updateIPv4, updateIPv6, nil
}

type updateHostOption struct {
	hostID  int64
	host    string
	innerIP string
	cloudID int64
	setter  map[string]interface{}
}

func (h *HostSnap) updateHostWithColletorMsg(header http.Header, rid string, hostOption updateHostOption) error {

	txnErr := h.CoreAPI.CoreService().Txn().AutoRunTxn(h.ctx, header, func() error {
		// get audit interface of host.
		audit := auditlog.NewHostAudit(h.CoreAPI.CoreService())
		kit := &rest.Kit{
			Rid:      rid,
			Header:   header,
			Ctx:      h.ctx,
			CCError:  h.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(header)),
			User:     common.CCSystemCollectorUserName,
			TenantID: common.BKDefaultTenantID,
		}

		// generate audit log for update host.
		hostData := make(mapstr.MapStr)
		err := json.Unmarshal([]byte(hostOption.host), &hostData)
		if err != nil {
			blog.Errorf("unmarshal host %s failed, err: %v, rid: %s", hostOption.host, err, rid)
			return err
		}
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).
			WithOperateFrom(metadata.FromDataCollection).WithUpdateFields(hostOption.setter)
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, 0, []mapstr.MapStr{hostData})
		if err != nil {
			blog.Errorf("generate host snap audit log failed before update host, host: %d/%s, err: %v, rid: %s",
				hostOption.hostID, hostOption.innerIP, err, rid)
			return err
		}

		// notice: needToUpdate 需要顺序，只能在更新数据库之前，删除需要忽略更新的字段
		for field := range ignoreCompareField {
			delete(hostOption.setter, field)
		}

		opt := &metadata.UpdateOption{
			Condition: map[string]interface{}{
				common.BKHostIDField: hostOption.hostID,
			},
			Data:       hostOption.setter,
			CanEditAll: true,
		}

		_, err = h.CoreAPI.CoreService().Instance().UpdateInstance(h.ctx, header, common.BKInnerObjIDHost, opt)
		if err != nil {
			blog.Errorf("snapshot changed, update host %d/%s snapshot failed, err: %v, rid: %s",
				hostOption.hostID, hostOption.innerIP, err, rid)
			return err
		}
		// save audit log.
		if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
			blog.Errorf("save host snap audit log failed after update host, host %d/%s, err: %v, rid: %s",
				hostOption.hostID, hostOption.innerIP, err, rid)
			return err
		}
		blog.V(5).Infof("snapshot for host changed, update success, host id: %d, ip: %s, cloud id: %d, rid: %s",
			hostOption.hostID, hostOption.innerIP, hostOption.cloudID, rid)

		return nil
	})

	return txnErr
}

// skipMsg verify the timestamp to determine whether the host sequence is correct, if it is old message, skip.
func (h *HostSnap) skipMsg(val gjson.Result, innerIP, rid string, hostID, cloudID int64) bool {
	if !val.Get("data.apiVer").Exists() {
		return false
	}

	key := redisConsumptionCheckPrefix + strconv.FormatInt(hostID, 10)
	timestamp, err := h.redisCli.Get(context.Background(), key).Result()
	if err != nil && !redis.IsNilErr(err) {
		blog.Errorf("get key: %s from redis err: %v, rid: %s", key, err, rid)
		return false
	}

	var oldTimestamp int64
	if timestamp != "" {
		oldTimestamp, err = strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			blog.Errorf("parseInt timestamp %s error, key: %s to redis err: %v, rid: %s",
				timestamp, key, err, rid)
			return false
		}
	}

	newTimestamp := val.Get("data.timestamp").Int()
	if err == nil && oldTimestamp != 0 && oldTimestamp > newTimestamp {
		blog.Warnf("skip host snapshot data update due to it is old data, host id: %d, ip: %s, "+
			"cloud id: %d, timestamp: %d, rid: %s", hostID, innerIP, cloudID, newTimestamp, rid)
		return true
	}

	randTime := util.RandInt64WithRange(int64(5), int64(10))
	err = h.redisCli.Set(context.Background(), key, newTimestamp, time.Minute*time.Duration(randTime)).Err()
	if err != nil {
		blog.Errorf("set key: %s to redis err: %v, rid: %s", key, err, rid)
		return false
	}

	return false
}

func needToUpdate(src, toCompare, addressing string) bool {
	// get data fluctuation limit
	changeRangePercent := getLimitConfig("datacollection.hostsnap.changeRangePercent",
		defaultChangeRangePercent, minChangeRangePercent)
	srcElements := gjson.GetMany(src, compareFields...)
	compareElements := gjson.GetMany(toCompare, compareFields...)
	for idx, field := range compareFields {
		if _, ok := ignoreCompareField[field]; ok {
			// 忽略变更对比的字段直接过滤掉
			continue
		}

		// 当不存在该字段时，需要跳过，防止对比出现差异记录了审计
		if !srcElements[idx].Exists() {
			continue
		}

		// compare these value with string directly to avoid empty value or null value.
		if srcElements[idx].String() != compareElements[idx].String() {
			compareField := compareFields[idx]
			// in the static scenario, it is not necessary to compare whether the inner IP has changed, but in the
			// dynamic ip scenario, it is necessary to compare whether the ip has changed.
			if addressing == common.BKAddressingStatic &&
				(compareField == common.BKHostInnerIPField || compareField == common.BKHostInnerIPv6Field) {
				continue
			}
			// tolerate bk_cpu, bk_disk, bk_mem changes less than the set value
			if compareField == "bk_cpu" || compareField == "bk_disk" || compareField == "bk_mem" {
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

type hostDiscoverMsg struct {
	ostype      string
	osname      string
	platform    string
	version     string
	cpumodule   string
	hostname    string
	osbit       string
	arch        string
	cpunum      int64
	disk        uint64
	mem         uint64
	outerMACArr []string
	innerMACArr []string
	hasOuterMAC bool
	hasInnerMAC bool
}

type hostInfo struct {
	innerIPv4  string
	outerIPv4  string
	innerIPv6  string
	outerIPv6  string
	addressing string
	updateIPv4 []string
	updateIPv6 []string
}

func getHostInfoFromMsgV10(val *gjson.Result, host *hostInfo) *hostDiscoverMsg {

	hostMsg := new(hostDiscoverMsg)

	hostMsg.cpumodule = strings.TrimSpace(val.Get("data.cpu.model").String())
	hostMsg.cpunum = val.Get("data.cpu.total").Int()

	hostMsg.disk = val.Get("data.disk.total").Uint() >> 10 >> 10 >> 10
	hostMsg.mem = val.Get("data.mem.total").Uint() >> 10 >> 10

	hostMsg.hostname = strings.TrimSpace(val.Get("data.system.hostname").String())
	hostMsg.ostype = strings.TrimSpace(val.Get("data.system.os").String())
	hostMsg.platform = strings.TrimSpace(val.Get("data.system.platform").String())
	hostMsg.version = val.Get("data.system.platVer").String()
	hostMsg.arch = strings.TrimSpace(val.Get("data.system.arch").String())

	switch strings.ToLower(hostMsg.ostype) {
	case common.HostOSTypeName[common.HostOSTypeEnumLinux]:
		hostMsg.version = strings.Replace(hostMsg.version, ".x86_64", "", 1)
		hostMsg.version = strings.Replace(hostMsg.version, ".i386", "", 1)
		hostMsg.osname = fmt.Sprintf("%s %s", hostMsg.ostype, hostMsg.platform)
		hostMsg.ostype = common.HostOSTypeEnumLinux
	case common.HostOSTypeName[common.HostOSTypeEnumWindows]:
		hostMsg.version = strings.Replace(hostMsg.version, "Microsoft ", "", 1)
		hostMsg.platform = strings.Replace(hostMsg.platform, "Microsoft ", "", 1)
		hostMsg.osname = fmt.Sprintf("%s", hostMsg.platform)
		hostMsg.ostype = common.HostOSTypeEnumWindows
	case common.HostOSTypeName[common.HostOSTypeEnumAIX]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumAIX
	case common.HostOSTypeName[common.HostOSTypeEnumUNIX]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumUNIX
	case common.HostOSTypeName[common.HostOSTypeEnumSolaris]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumSolaris
	case common.HostOSTypeName[common.HostOSTypeEnumHpUX]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumHpUX
	case common.HostOSTypeName[common.HostOSTypeEnumFreeBSD]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumFreeBSD
	case common.HostOSTypeName[common.HostOSTypeEnumMacOS]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumMacOS
	default:
		hostMsg.osname = fmt.Sprintf("%s", hostMsg.platform)
	}

	hostMsg.version = strings.TrimSpace(hostMsg.version)
	hostMsg.osname = strings.TrimSpace(hostMsg.osname)

	// 静态寻址方式时：根据主机现在的内网ipv4和ipv6地址，如果为ipv6单栈，那么就通过主机当前的内外网ipv6地址，找到对应的内外网mac地址；否则使
	// 用主机当前的内外ipv4地址，去找到对应的内外网mac地址；
	// 动态寻址方式时：根据采集器上报上来的内网ipv4和ipv6地址，如果为ipv6单栈，那么就通过上报上来的内网ipv6地址，和当前主机的外网ipv6地址，找
	// 到对应的内外网mac地址；否则使用采集器上报上来的内网ipv4地址，和当前主机的外网ipv4地址，找到对应的内外网mac地址。
	var innerMACArr, outerMACArr []string
	switch host.addressing {
	case common.BKAddressingStatic:
		if host.innerIPv4 == "" && host.innerIPv6 != "" {
			innerIPv6Arr := strings.Split(host.innerIPv6, ",")
			outerIPv6Arr := strings.Split(host.outerIPv6, ",")
			innerMACArr, outerMACArr = getMacAddr(val, innerIPv6Arr, outerIPv6Arr)
		} else {
			innerIPv4Arr := strings.Split(host.innerIPv4, ",")
			outerIPv4Arr := strings.Split(host.outerIPv4, ",")
			innerMACArr, outerMACArr = getMacAddr(val, innerIPv4Arr, outerIPv4Arr)
		}

	case common.BKAddressingDynamic:
		if len(host.updateIPv4) == 0 && len(host.updateIPv6) != 0 {
			outerIPv6Arr := strings.Split(host.outerIPv6, ",")
			innerMACArr, outerMACArr = getMacAddr(val, host.updateIPv6, outerIPv6Arr)
		} else {
			outerIPv4Arr := strings.Split(host.outerIPv4, ",")
			innerMACArr, outerMACArr = getMacAddr(val, host.updateIPv4, outerIPv4Arr)
		}

	default:
		blog.Errorf("can not support this host addressing, host: %v", host)
	}

	if len(innerMACArr) != 0 {
		hostMsg.hasInnerMAC = true
		hostMsg.innerMACArr = innerMACArr
	}

	if len(outerMACArr) != 0 {
		hostMsg.hasOuterMAC = true
		hostMsg.outerMACArr = outerMACArr
	}

	hostMsg.osbit = strings.TrimSpace(val.Get("data.system.sysType").String())

	printSetterInfo(hostMsg, host.innerIPv4, host.outerIPv4)
	return hostMsg
}

func getMacAddr(val *gjson.Result, innerIP, outerIP []string) ([]string, []string) {
	innerIPMap := make(map[string]int)
	for index, ip := range innerIP {
		innerIPMap[ip] = index
	}

	outerIPMap := make(map[string]int)
	for index, ip := range outerIP {
		outerIPMap[ip] = index
	}

	innerMACArr, outerMACArr := make([]string, len(innerIP)), make([]string, len(outerIP))

	for _, inter := range val.Get("data.net.interface").Array() {
		for _, addr := range inter.Get("addrs").Array() {
			splitAddr := strings.Split(addr.String(), "/")
			if len(splitAddr) == 0 {
				continue
			}
			ip := splitAddr[0]
			var err error
			if strings.Contains(ip, ":") {
				ip, err = common.GetIPv4IfEmbeddedInIPv6(ip)
				if err != nil {
					blog.Warnf("get ip failed, addr: %s, err: %v", ip, err)
					continue
				}
			}

			if strings.Contains(ip, ":") {
				ip, err = common.ConvertIPv6ToStandardFormat(ip)
				if err != nil {
					blog.Warnf("convert ipv6 to standard format failed, addr: %s, err: %v", ip, err)
					continue
				}
			}

			if index, exists := innerIPMap[ip]; exists {
				innerMAC := strings.TrimSpace(inter.Get("mac").String())
				if len(innerMACArr[index]) == 0 {
					innerMACArr[index] = innerMAC
				} else {
					blog.Errorf("innerIP has different mac address, use first mac address by default, ip: %s, "+
						"mac: %s, other mac: %s", ip, innerMACArr[index], innerMAC)
				}
				continue
			}

			if index, exists := outerIPMap[ip]; exists {
				outerMAC := strings.TrimSpace(inter.Get("mac").String())
				if len(outerMACArr[index]) == 0 {
					outerMACArr[index] = outerMAC
				} else {
					blog.Errorf("outerIP has different mac address, use first mac address by default, ip: %s, "+
						"mac: %s, other mac: %s", ip, outerMACArr[index], outerMAC)
				}
			}
		}
	}

	return innerMACArr, outerMACArr
}

func getOsInfoFromMsg(val *gjson.Result, innerIP, outerIP string) *hostDiscoverMsg {

	hostMsg := new(hostDiscoverMsg)
	hostMsg.cpumodule = strings.TrimSpace(val.Get("data.cpu.cpuinfo.0.modelName").String())

	for _, core := range val.Get("data.cpu.cpuinfo.#.cores").Array() {
		hostMsg.cpunum += core.Int()
	}

	for _, disktotal := range val.Get("data.disk.usage.#.total").Array() {
		hostMsg.disk += disktotal.Uint()
	}
	hostMsg.disk = hostMsg.disk >> 10 >> 10 >> 10

	hostMsg.mem = val.Get("data.mem.meminfo.total").Uint()
	hostMsg.hostname = strings.TrimSpace(val.Get("data.system.info.hostname").String())
	hostMsg.ostype = strings.TrimSpace(val.Get("data.system.info.os").String())
	hostMsg.platform = strings.TrimSpace(val.Get("data.system.info.platform").String())
	hostMsg.version = val.Get("data.system.info.platformVersion").String()

	switch strings.ToLower(hostMsg.ostype) {
	case common.HostOSTypeName[common.HostOSTypeEnumLinux]:
		hostMsg.version = strings.Replace(hostMsg.version, ".x86_64", "", 1)
		hostMsg.version = strings.Replace(hostMsg.version, ".i386", "", 1)
		hostMsg.osname = fmt.Sprintf("%s %s", hostMsg.ostype, hostMsg.platform)
		hostMsg.ostype = common.HostOSTypeEnumLinux
	case common.HostOSTypeName[common.HostOSTypeEnumWindows]:
		hostMsg.version = strings.Replace(hostMsg.version, "Microsoft ", "", 1)
		hostMsg.platform = strings.Replace(hostMsg.platform, "Microsoft ", "", 1)
		hostMsg.osname = fmt.Sprintf("%s", hostMsg.platform)
		hostMsg.ostype = common.HostOSTypeEnumWindows
	case common.HostOSTypeName[common.HostOSTypeEnumAIX]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumAIX
	case common.HostOSTypeName[common.HostOSTypeEnumUNIX]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumUNIX
	case common.HostOSTypeName[common.HostOSTypeEnumSolaris]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumSolaris
	case common.HostOSTypeName[common.HostOSTypeEnumHpUX]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumHpUX
	case common.HostOSTypeName[common.HostOSTypeEnumFreeBSD]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumFreeBSD
	case common.HostOSTypeName[common.HostOSTypeEnumMacOS]:
		hostMsg.osname = hostMsg.platform
		hostMsg.ostype = common.HostOSTypeEnumMacOS
	default:
		hostMsg.osname = fmt.Sprintf("%s", hostMsg.platform)
	}

	hostMsg.version = strings.TrimSpace(hostMsg.version)
	hostMsg.osname = strings.TrimSpace(hostMsg.osname)

	innerIPMap, outerIPMap, innerLen, outerLen := getIpAddrs(innerIP, outerIP)

	hostMsg.outerMACArr, hostMsg.innerMACArr = make([]string, outerLen), make([]string, innerLen)

	for _, inter := range val.Get("data.net.interface").Array() {
		for _, addr := range inter.Get("addrs.#.addr").Array() {
			splitAddr := strings.Split(addr.String(), "/")
			if len(splitAddr) == 0 {
				continue
			}
			ip := splitAddr[0]
			if index, exists := innerIPMap[ip]; exists {
				hostMsg.hasInnerMAC = true
				innerMAC := strings.TrimSpace(inter.Get("hardwareaddr").String())
				if len(hostMsg.innerMACArr[index]) == 0 {
					hostMsg.innerMACArr[index] = innerMAC
				}
			} else if index, exists := outerIPMap[ip]; exists {
				hostMsg.hasOuterMAC = true
				outerMAC := strings.TrimSpace(inter.Get("hardwareaddr").String())
				if len(hostMsg.outerMACArr[index]) == 0 {
					hostMsg.outerMACArr[index] = outerMAC
				}
			}
		}
	}

	hostMsg.osbit = strings.TrimSpace(val.Get("data.system.info.systemtype").String())
	hostMsg.mem = hostMsg.mem >> 10 >> 10

	printSetterInfo(hostMsg, innerIP, outerIP)
	return hostMsg
}

func getIpAddrs(innerIP, outerIP string) (map[string]int, map[string]int, int, int) {
	innerIPArr := strings.Split(innerIP, ",")
	innerIPMap := make(map[string]int)
	for index, ip := range innerIPArr {
		innerIPMap[ip] = index
	}

	outerIPArr := strings.Split(outerIP, ",")
	outerIPMap := make(map[string]int)
	for index, ip := range outerIPArr {
		outerIPMap[ip] = index
	}
	return innerIPMap, outerIPMap, len(innerIPMap), len(outerIPMap)
}

func parseSetter(val *gjson.Result, innerIP, outerIP string) (map[string]interface{}, string) {

	hostMsg := getOsInfoFromMsg(val, innerIP, outerIP)

	setter := make(map[string]interface{})
	raw := strings.Builder{}
	raw.WriteByte('{')

	if hostMsg.cpunum > 0 {
		setter["bk_cpu"] = hostMsg.cpunum
		raw.WriteString("\"bk_cpu\":")
		raw.WriteString(strconv.FormatInt(hostMsg.cpunum, 10))
	}

	if hostMsg.cpumodule != "" {
		setter["bk_cpu_module"] = hostMsg.cpumodule
		raw.WriteString(",")
		raw.WriteString("\"bk_cpu_module\":")
		raw.Write([]byte("\"" + hostMsg.cpumodule + "\""))
	}

	if hostMsg.disk > 0 {
		setter["bk_disk"] = hostMsg.disk
		raw.WriteString(",")
		raw.WriteString("\"bk_disk\":")
		raw.WriteString(strconv.FormatUint(hostMsg.disk, 10))
	}

	if hostMsg.mem > 0 {
		setter["bk_mem"] = hostMsg.mem
		raw.WriteString(",")
		raw.WriteString("\"bk_mem\":")
		raw.WriteString(strconv.FormatUint(hostMsg.mem, 10))
	}

	if hostMsg.ostype != "" {
		setter["bk_os_type"] = hostMsg.ostype
		raw.WriteString(",")
		raw.WriteString("\"bk_os_type\":")
		raw.Write([]byte("\"" + hostMsg.ostype + "\""))
	}

	if hostMsg.osname != "" {
		setter["bk_os_name"] = hostMsg.osname
		raw.WriteString(",")
		raw.WriteString("\"bk_os_name\":")
		raw.Write([]byte("\"" + hostMsg.osname + "\""))
	}

	if hostMsg.version != "" {
		setter["bk_os_version"] = hostMsg.version
		raw.WriteString(",")
		raw.WriteString("\"bk_os_version\":")
		raw.Write([]byte("\"" + hostMsg.version + "\""))
	}

	if hostMsg.hostname != "" {
		setter["bk_host_name"] = hostMsg.hostname
		raw.WriteString(",")
		raw.WriteString("\"bk_host_name\":")
		raw.Write([]byte("\"" + hostMsg.hostname + "\""))
	}

	if outerIP == "" || hostMsg.hasOuterMAC {
		outerMAC := strings.Join(hostMsg.outerMACArr, ",")
		setter["bk_outer_mac"] = outerMAC
		raw.WriteString(",")
		raw.WriteString("\"bk_outer_mac\":")
		raw.Write([]byte("\"" + outerMAC + "\""))
	}

	if hostMsg.hasInnerMAC {
		innerMAC := strings.Join(hostMsg.innerMACArr, ",")
		setter["bk_mac"] = innerMAC
		raw.WriteString(",")
		raw.WriteString("\"bk_mac\":")
		raw.Write([]byte("\"" + innerMAC + "\""))
	}

	if hostMsg.osbit == "" {
		setter["bk_os_bit"] = hostMsg.osbit
		raw.WriteString(",")
		raw.WriteString("\"bk_os_bit\":")
		raw.Write([]byte("\"" + hostMsg.osbit + "\""))
	}

	raw.WriteByte('}')
	return setter, raw.String()
}

func printSetterInfo(hostMsg *hostDiscoverMsg, innerIP, outerIP string) {
	if hostMsg.cpunum <= 0 {
		blog.V(4).Infof("bk_cpu not found in message for %s", innerIP)
	}
	if hostMsg.cpumodule == "" {
		blog.V(4).Infof("bk_cpu_module not found in message for %s", innerIP)
	}
	if hostMsg.disk <= 0 {
		blog.V(4).Infof("bk_disk not found in message for %s", innerIP)
	}
	if hostMsg.mem <= 0 {
		blog.V(4).Infof("bk_mem not found in message for %s", innerIP)
	}
	if hostMsg.ostype == "" {
		blog.V(4).Infof("bk_os_type not found in message for %s", innerIP)
	}
	if hostMsg.osname == "" {
		blog.V(4).Infof("bk_os_name not found in message for %s", innerIP)
	}
	if hostMsg.version == "" {
		blog.V(4).Infof("bk_os_version not found in message for %s", innerIP)
	}
	if hostMsg.hostname == "" {
		blog.V(4).Infof("bk_host_name not found in message for %s", innerIP)
	}
	if outerIP != "" && !hostMsg.hasOuterMAC {
		blog.V(4).Infof("bk_outer_mac not found in message for %s", outerIP)
	}
	if !hostMsg.hasInnerMAC {
		blog.V(4).Infof("bk_mac not found in message for %s", innerIP)
	}
	if hostMsg.osbit == "" {
		blog.V(4).Infof("bk_os_bit not found in message for %s", innerIP)
	}
}

// NOCC:golint/fnsize(解析操作需要放到一个函数中)
func parseV10Setter(val *gjson.Result, host *hostInfo) (
	map[string]interface{}, string) {

	hostMsg := getHostInfoFromMsgV10(val, host)
	setter, raw := make(map[string]interface{}), strings.Builder{}
	raw.WriteByte('{')

	if hostMsg.cpunum > 0 {
		setter["bk_cpu"] = hostMsg.cpunum
		raw.WriteString("\"bk_cpu\":")
		raw.WriteString(strconv.FormatInt(hostMsg.cpunum, 10))
	}

	if hostMsg.cpumodule != "" {
		setter["bk_cpu_module"] = hostMsg.cpumodule
		raw.WriteString(",")
		raw.WriteString("\"bk_cpu_module\":")
		raw.Write([]byte("\"" + hostMsg.cpumodule + "\""))
	}

	if hostMsg.disk > 0 {
		setter["bk_disk"] = hostMsg.disk
		raw.WriteString(",")
		raw.WriteString("\"bk_disk\":")
		raw.WriteString(strconv.FormatUint(hostMsg.disk, 10))
	}

	if hostMsg.mem > 0 {
		setter["bk_mem"] = hostMsg.mem
		raw.WriteString(",")
		raw.WriteString("\"bk_mem\":")
		raw.WriteString(strconv.FormatUint(hostMsg.mem, 10))
	}

	if hostMsg.ostype != "" {
		setter["bk_os_type"] = hostMsg.ostype
		raw.WriteString(",")
		raw.WriteString("\"bk_os_type\":")
		raw.Write([]byte("\"" + hostMsg.ostype + "\""))
	}

	if hostMsg.osname != "" {
		setter["bk_os_name"] = hostMsg.osname
		raw.WriteString(",")
		raw.WriteString("\"bk_os_name\":")
		raw.Write([]byte("\"" + hostMsg.osname + "\""))
	}

	if hostMsg.version != "" {
		setter["bk_os_version"] = hostMsg.version
		raw.WriteString(",")
		raw.WriteString("\"bk_os_version\":")
		raw.Write([]byte("\"" + hostMsg.version + "\""))
	}

	if hostMsg.hostname != "" {
		setter["bk_host_name"] = hostMsg.hostname
		raw.WriteString(",")
		raw.WriteString("\"bk_host_name\":")
		raw.Write([]byte("\"" + hostMsg.hostname + "\""))
	}

	if len(host.updateIPv4) > 0 {
		setter[common.BKHostInnerIPField] = strings.Join(host.updateIPv4, ",")
		raw.WriteString(",")
		raw.WriteString("\"bk_host_innerip\":")
		raw.Write([]byte("\"" + strings.Join(host.updateIPv4, ",") + "\""))
	}

	if len(host.updateIPv6) > 0 {
		setter[common.BKHostInnerIPv6Field] = strings.Join(host.updateIPv6, ",")
		raw.WriteString(",")
		raw.WriteString("\"bk_host_innerip_v6\":")
		raw.Write([]byte("\"" + strings.Join(host.updateIPv6, ",") + "\""))
	}

	if host.outerIPv4 == "" || hostMsg.hasOuterMAC {
		outerMAC := strings.Join(hostMsg.outerMACArr, ",")
		setter["bk_outer_mac"] = outerMAC
		raw.WriteString(",")
		raw.WriteString("\"bk_outer_mac\":")
		raw.Write([]byte("\"" + outerMAC + "\""))
	}

	if hostMsg.hasInnerMAC {
		innerMAC := strings.Join(hostMsg.innerMACArr, ",")
		setter["bk_mac"] = innerMAC
		raw.WriteString(",")
		raw.WriteString("\"bk_mac\":")
		raw.Write([]byte("\"" + innerMAC + "\""))
	}

	if hostMsg.osbit != "" {
		setter["bk_os_bit"] = hostMsg.osbit
		raw.WriteString(",")
		raw.WriteString("\"bk_os_bit\":")
		raw.Write([]byte("\"" + hostMsg.osbit + "\""))
	}

	if hostMsg.arch != "" {
		setter["bk_cpu_architecture"] = hostMsg.arch
		raw.WriteString(",")
		raw.WriteString("\"bk_cpu_architecture\":")
		raw.Write([]byte("\"" + hostMsg.arch + "\""))
	}

	raw.WriteByte('}')
	return setter, raw.String()
}

func (h *HostSnap) getHostByAgentID(header http.Header, rid, agentID string) (string, error) {

	opt := &metadata.SearchHostWithAgentID{
		AgentID: agentID,
		Fields:  reqireFields,
	}

	host, err := h.Engine.CoreAPI.CacheService().Cache().Host().SearchHostWithAgentID(context.Background(), header, opt)
	if err != nil {
		blog.Errorf("get host info with agentID: %s failed, err: %v, rid: %s", agentID, err, rid)
		return "", err
	}
	return host, nil
}

func (h *HostSnap) getHostByVal(header http.Header, rid string, cloudID int64, ips []string) (string, error) {

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

		host, err := h.Engine.CoreAPI.CacheService().Cache().Host().SearchHostWithInnerIPForStatic(context.Background(),
			header, opt)
		if err != nil {
			blog.Errorf("get host info with ip: %s, cloud id: %d failed, err: %v, rid: %s", ip, cloudID, err, rid)
			if ccErr, ok := err.(ccErr.CCErrorCoder); ok {
				if ccErr.GetCode() == common.CCErrCommDBSelectFailed {
					h.filter.Set(ip, cloudID)
				}
			}
			// do not return, continue search with next ip
			continue
		}

		// not find host
		if len(host) == 0 {
			continue
		}
		return host, nil
	}

	return "", errors.New("can not find ip detail from cache")
}

func getIPsFromMsg(val *gjson.Result, agentID string, rid string) ([]string, []string) {
	ipv4Map := make(map[string]struct{})
	ipv6Map := make(map[string]struct{})

	rootIP := val.Get("ip").String()
	rootIP = strings.TrimSpace(rootIP)
	if rootIP != metadata.IPv4LoopBackIpPrefix && rootIP != metadata.IPv6LoopBackIp &&
		!strings.HasPrefix(rootIP, metadata.IPv6LinkLocalAddressPrefix) && net.ParseIP(rootIP) != nil {
		if strings.Contains(rootIP, ":") {
			ipv6Map[rootIP] = struct{}{}
		} else {
			ipv4Map[rootIP] = struct{}{}
		}
	}
	// need to be compatible with the old and new versions of the format
	// new format:
	//  "net":{
	//            "interface":[
	//                {
	//                    "addrs":["::1/64", "127.0.0.1/32"]
	//                }
	//            ]
	//        }

	//  old format:
	//   "interface":[
	//        {
	//            "addrs":[
	//                {
	//                    "addr":"::1/128"
	//                }
	//            ]
	//        }
	//    ]

	interfaces := make([]gjson.Result, 0)
	if val.Get("data.apiVer").Exists() && val.Get("data.apiVer").String() >= "v1.0" {
		interfaces = val.Get("data.net.interface.#.addrs").Array()
	} else {
		interfaces = val.Get("data.net.interface.#.addrs.#.addr").Array()
	}

	for _, addrs := range interfaces {
		for _, addr := range addrs.Array() {
			ip := strings.Split(addr.String(), "/")[0]
			ip = strings.TrimSpace(ip)
			if ip == metadata.IPv4LoopBackIpPrefix || ip == metadata.IPv6LoopBackIp ||
				strings.HasPrefix(ip, metadata.IPv6LinkLocalAddressPrefix) {
				continue
			}

			var err error
			if strings.Contains(ip, ":") {
				ip, err = common.GetIPv4IfEmbeddedInIPv6(ip)
				if err != nil {
					blog.Warnf("get ip failed, agentID: %v, addr: %v, err: %v", agentID, addr, err, rid)
					continue
				}
			}

			if net.ParseIP(ip) == nil {
				// invalid ip address
				continue
			}

			if strings.Contains(ip, ":") {
				ipv6Map[ip] = struct{}{}
			} else {
				ipv4Map[ip] = struct{}{}
			}
		}
	}

	ipv4List := make([]string, 0)
	for ipv4 := range ipv4Map {
		ipv4List = append(ipv4List, ipv4)
	}

	ipv6List := make([]string, 0)
	for ipv6 := range ipv6Map {
		ipv6List = append(ipv6List, ipv6)
	}
	return ipv4List, ipv6List
}

// saveHostsnap save host snapshot in redis
func (h *HostSnap) saveHostsnap(header http.Header, hostData *gjson.Result, hostID int64) error {
	rid := httpheader.GetRid(header)

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
	rid := util.GenerateRID()
	header := headerutil.GenCommonHeader(common.CCSystemCollectorUserName, common.BKDefaultTenantID, rid)
	return header, rid
}

// MockMessage TODO
const MockMessage = "{\"localTime\": \"2017-09-19 16:57:00\", \"data\": \"{\\\"ip\\\":\\\"127.0.0.1\\\"," +
	"\\\"bizid\\\":0,\\\"cloudid\\\":0,\\\"data\\\":{\\\"timezone\\\":8,\\\"datetime\\\":\\\"2017-09-19 16:57:07\\\"," +
	"\\\"utctime\\\":\\\"2017-09-19 08:57:07\\\",\\\"country\\\":\\\"Asia\\\",\\\"city\\\":\\\"Shanghai\\\"," +
	"\\\"cpu\\\":{\\\"cpuinfo\\\":[{\\\"cpu\\\":0,\\\"vendorID\\\":\\\"GenuineIntel\\\",\\\"family\\\":\\\"6\\\"," +
	"\\\"model\\\":\\\"63\\\",\\\"stepping\\\":2,\\\"physicalID\\\":\\\"0\\\",\\\"coreID\\\":\\\"0\\\"," +
	"\\\"cores\\\":1,\\\"modelName\\\":\\\"Intel(R) Xeon(R) CPU E5-26xx v3\\\",\\\"mhz\\\":2294.01," +
	"\\\"cacheSize\\\":4096,\\\"flags\\\":[\\\"fpu\\\",\\\"vme\\\",\\\"de\\\",\\\"pse\\\",\\\"tsc\\\"," +
	"\\\"msr\\\",\\\"pae\\\",\\\"mce\\\",\\\"cx8\\\",\\\"apic\\\",\\\"sep\\\",\\\"mtrr\\\",\\\"pge\\\"," +
	"\\\"mca\\\",\\\"cmov\\\",\\\"pat\\\",\\\"pse36\\\",\\\"clflush\\\",\\\"mmx\\\",\\\"fxsr\\\"," +
	"\\\"sse\\\",\\\"sse2\\\",\\\"ss\\\",\\\"ht\\\",\\\"syscall\\\",\\\"nx\\\",\\\"lm\\\",\\\"constant_tsc\\\"," +
	"\\\"up\\\",\\\"rep_good\\\",\\\"unfair_spinlock\\\",\\\"pni\\\",\\\"pclmulqdq\\\",\\\"ssse3\\\",\\\"fma\\\"," +
	"\\\"cx16\\\",\\\"pcid\\\",\\\"sse4_1\\\",\\\"sse4_2\\\",\\\"x2apic\\\",\\\"movbe\\\",\\\"popcnt\\\"," +
	"\\\"tsc_deadline_timer\\\",\\\"aes\\\",\\\"xsave\\\",\\\"avx\\\",\\\"f16c\\\",\\\"rdrand\\\",\\\"hypervisor\\\"," +
	"\\\"lahf_lm\\\",\\\"abm\\\",\\\"xsaveopt\\\",\\\"bmi1\\\",\\\"avx2\\\",\\\"bmi2\\\"],\\\"microcode\\\":" +
	"\\\"1\\\"}],\\\"per_usage\\\":[3.0232169701043103],\\\"total_usage\\\":3.0232169701043103,\\\"per_stat\\\"" +
	":[{\\\"cpu\\\":\\\"cpu0\\\",\\\"user\\\":5206.09,\\\"system\\\":6107.04,\\\"idle\\\":337100.84,\\\"nice\\\"" +
	":6.68,\\\"iowait\\\":528.24,\\\"irq\\\":0.02,\\\"softirq\\\":13.48,\\\"steal\\\":0,\\\"guest\\\":0," +
	"\\\"guestNice\\\":0,\\\"stolen\\\":0}],\\\"total_stat\\\":{\\\"cpu\\\":\\\"cpu-total\\\",\\\"user\\\"" +
	":5206.09,\\\"system\\\":6107.04,\\\"idle\\\":337100.84,\\\"nice\\\":6.68,\\\"iowait\\\":528.24,\\\"irq\\\"" +
	":0.02,\\\"softirq\\\":13.48,\\\"steal\\\":0,\\\"guest\\\":0,\\\"guestNice\\\":0,\\\"stolen\\\":0}}," +
	"\\\"env\\\":{\\\"crontab\\\":[{\\\"user\\\":\\\"root\\\",\\\"content\\\":\\\"#secu-tcs-agent monitor, " +
	"install at Fri Sep 15 16:12:02 CST 2017\\\\n* * * * * /usr/local/sa/agent/secu-tcs-agent-mon-safe.sh " +
	"/usr/local/sa/agent \\\\u003e /dev/null 2\\\\u003e\\\\u00261\\\\n*/1 * * * * /usr/local/qcloud/stargate" +
	"/admin/start.sh \\\\u003e /dev/null 2\\\\u003e\\\\u00261 \\\\u0026\\\\n*/20 * * * * /usr/sbin/ntpdate " +
	"ntpupdate.tencentyun.com \\\\u003e/dev/null \\\\u0026\\\\n*/1 * * * * cd /usr/local/gse/gseagent; " +
	"./cron_agent.sh 1\\\\u003e/dev/null 2\\\\u003e\\\\u00261\\\\n\\\"}],\\\"host\\\":\\\"127.0.0.1  localhost" +
	"  localhost.localdomain  VM_0_31_centos\\\\n::1         localhost localhost.localdomain localhost6 " +
	"localhost6.localdomain6\\\\n\\\",\\\"route\\\":\\\"Kernel IP routing table\\\\nDestination     " +
	"Gateway         Genmask         Flags Metric Ref    Use Iface\\\\n10.0.0.0        0.0.0.0     " +
	"    255.255.255.0   U     0      0        0 eth0\\\\n169.254.0.0     0.0.0.0         255.255.0.0  " +
	"   U     1002   0        0 eth0\\\\n0.0.0.0         127.0.0.1        0.0.0.0         UG    0      0    " +
	"    0 eth0\\\\n\\\"},\\\"disk\\\":{\\\"diskstat\\\":{\\\"vda1\\\":{\\\"major\\\":252,\\\"minor\\\"" +
	":1,\\\"readCount\\\":24347,\\\"mergedReadCount\\\":570,\\\"writeCount\\\":696357,\\\"mergedWriteCount\\\"" +
	":4684783,\\\"readBytes\\\":783955968,\\\"writeBytes\\\":22041231360,\\\"readSectors\\\":1531164,\\\"" +
	"writeSectors\\\":43049280,\\\"readTime\\\":80626,\\\"writeTime\\\":12704736,\\\"iopsInProgress\\\":0," +
	"\\\"ioTime\\\":822057,\\\"weightedIoTime\\\":12785026,\\\"name\\\":\\\"vda1\\\",\\\"serialNumber\\\"" +
	":\\\"\\\",\\\"speedIORead\\\":0,\\\"speedByteRead\\\":0,\\\"speedIOWrite\\\":2.9,\\\"speedByteWrite\\\"" +
	":171144.53333333333,\\\"util\\\":0.0025666666666666667,\\\"avgrq_sz\\\":115.26436781609195,\\\"avgqu_sz\\\"" +
	":0.06568333333333334,\\\"await\\\":22.649425287356323,\\\"svctm\\\":0.8850574712643678}},\\\"" +
	"partition\\\":[{\\\"device\\\":\\\"/dev/vda1\\\",\\\"mountpoint\\\":\\\"/\\\",\\\"fstype\\\":\\\"ext3\\\"," +
	"\\\"opts\\\":\\\"rw,noatime,acl,user_xattr\\\"}],\\\"usage\\\":[{\\\"path\\\":\\\"/\\\",\\\"fstype\\\":" +
	"\\\"ext2/ext3\\\",\\\"total\\\":52843638784,\\\"free\\\":47807447040,\\\"used\\\":2351915008," +
	"\\\"usedPercent\\\":4.4507060113962345,\\\"inodesTotal\\\":3276800,\\\"inodesUsed\\\":29554,\\\"inodesFree" +
	"\\\":3247246,\\\"inodesUsedPercent\\\":0.9019165039062501}]},\\\"load\\\":{\\\"load_avg\\\":{\\\"load1\\\"" +
	":0,\\\"load5\\\":0,\\\"load15\\\":0}},\\\"mem\\\":{\\\"meminfo\\\":{\\\"total\\\":1044832256,\\\"available" +
	"\\\":805912576,\\\"used\\\":238919680,\\\"usedPercent\\\":22.866797864249705,\\\"free\\\":92041216," +
	"\\\"active\\\":521183232,\\\"inactive\\\":352964608,\\\"wired\\\":0,\\\"buffers\\\":110895104,\\\"cached" +
	"\\\":602976256,\\\"writeback\\\":0,\\\"dirty\\\":151552,\\\"writebacktmp\\\":0},\\\"vmstat\\\":" +
	"{\\\"total\\\":0,\\\"used\\\":0,\\\"free\\\":0,\\\"usedPercent\\\":0,\\\"sin\\\":0,\\\"sout\\\":0}}," +
	"\\\"net\\\":{\\\"interface\\\":[{\\\"mtu\\\":65536,\\\"name\\\":\\\"lo\\\",\\\"hardwareaddr\\\":\\\"" +
	"28:31:52:1d:c6:0a\\\",\\\"flags\\\":[\\\"up\\\",\\\"loopback\\\"],\\\"addrs\\\":[{\\\"addr\\\":\\\"" +
	"127.0.0.1/8\\\"}]},{\\\"mtu\\\":1500,\\\"name\\\":\\\"eth0\\\",\\\"hardwareaddr\\\":\\\"52:54:00:19:" +
	"2e:e8\\\",\\\"flags\\\":[\\\"up\\\",\\\"broadcast\\\",\\\"multicast\\\"],\\\"addrs\\\":[{\\\"addr\\\":\\" +
	"\"127.0.0.1/24\\\"}]}],\\\"dev\\\":[{\\\"name\\\":\\\"lo\\\",\\\"speedSent\\\":0,\\\"speedRecv\\\":0,\\\"" +
	"speedPacketsSent\\\":0,\\\"speedPacketsRecv\\\":0,\\\"bytesSent\\\":604,\\\"bytesRecv\\\":604,\\\"packet" +
	"sSent\\\":2,\\\"packetsRecv\\\":2,\\\"errin\\\":0,\\\"errout\\\":0,\\\"dropin\\\":0,\\\"dropout\\\":0,\\\"" +
	"fifoin\\\":0,\\\"fifoout\\\":0},{\\\"name\\\":\\\"eth0\\\",\\\"speedSent\\\":574,\\\"speedRecv\\\":214,\\" +
	"\"speedPacketsSent\\\":3,\\\"speedPacketsRecv\\\":2,\\\"bytesSent\\\":161709123,\\\"bytesRecv\\\":285910" +
	"298,\\\"packetsSent\\\":1116625,\\\"packetsRecv\\\":1167796,\\\"errin\\\":0,\\\"errout\\\":0,\\\"dropin\\" +
	"\":0,\\\"dropout\\\":0,\\\"fifoin\\\":0,\\\"fifoout\\\":0}],\\\"netstat\\\":{\\\"established\\\":2,\\\"syn" +
	"cSent\\\":1,\\\"synRecv\\\":0,\\\"finWait1\\\":0,\\\"finWait2\\\":0,\\\"timeWait\\\":0,\\\"close\\\":0,\\\"" +
	"closeWait\\\":0,\\\"lastAck\\\":0,\\\"listen\\\":2,\\\"closing\\\":0},\\\"protocolstat\\\":[{\\\"protocol\\\"" +
	":\\\"udp\\\",\\\"stats\\\":{\\\"inDatagrams\\\":176253,\\\"inErrors\\\":0,\\\"noPorts\\\":1,\\\"outDatagrams" +
	"\\\":199569,\\\"rcvbufErrors\\\":0,\\\"sndbufErrors\\\":0}}]},\\\"system\\\":{\\\"info\\\":{\\\"hostname\\\"" +
	":\\\"VM_0_31_centos\\\",\\\"uptime\\\":348315,\\\"bootTime\\\":1505463112,\\\"procs\\\":142,\\\"os\\\":\\" +
	"\"linux\\\",\\\"platform\\\":\\\"centos\\\",\\\"platformFamily\\\":\\\"rhel\\\",\\\"platformVersion\\\":" +
	"\\\"6.2\\\",\\\"kernelVersion\\\":\\\"2.6.32-504.30.3.el6.x86_64\\\",\\\"virtualizationSystem\\\":\\\"\\\"" +
	",\\\"virtualizationRole\\\":\\\"\\\",\\\"hostid\\\":\\\"96D0F4CA-2157-40E6-BF22-6A7CD9B6EB8C\\\",\\\"syst" +
	"emtype\\\":\\\"64-bit\\\"}}}}\", \"timestamp\": 1505811427, \"dtEventTime\": \"2017-09-19 16:57:07\", \"" +
	"dtEventTimeStamp\": 1505811427000}"
