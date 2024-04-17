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

package hostidentifier

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/apigw/gse"
	"configcenter/src/thirdparty/gse/client"
	getstatus "configcenter/src/thirdparty/gse/get_agent_state_forsyncdata"
	pushfile "configcenter/src/thirdparty/gse/push_file_forsyncdata"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

const (
	// agentOnStatus express agent status is on
	agentOnStatus = 1
	// retryTimes indicates how many times will retry if fail
	retryTimes = 10
	// hostIdentifierCursor host identifier cursor in redis
	hostIdentifierCursor = "host_identifier:cursor"
	// callerName the name of the service that call the gse server
	callerName = "CMDB"
	// hostIdentifierBatchSyncPerLimit batch sync limit for host identifier
	hostIdentifierBatchSyncPerLimit = 200
	// metricsNamespacePrefix is prefix of metrics namespace.
	metricsNamespacePrefix = "cmdb_sync_data"
	// v2ApiAgentOnStatus express agent status is on about use gse api gateway api
	v2ApiAgentOnStatus = "2"
)

type hostIdentifierMetric struct {
	// getAgentStatusTotal call gse get agent status api total
	getAgentStatusTotal *prometheus.CounterVec
	// pushFileTotal call gse push file api total
	pushFileTotal *prometheus.CounterVec
	// getResultTotal call gse get task result api total
	getResultTotal *prometheus.CounterVec
	// agentStatusTotal host agent status total
	agentStatusTotal *prometheus.CounterVec
	// hostResultTotal host result total
	hostResultTotal *prometheus.CounterVec
}

// HostIdentifier manipulate the structure of the host Identifier
type HostIdentifier struct {
	redisCli            redis.Client
	engine              *backbone.Engine
	ctx                 context.Context
	gseTaskServerClient *client.GseTaskServerClient
	gseApiServerClient  *client.GseApiServerClient
	gseApiGWClient      gse.GseClientInterface
	apiVersion          types.ApiVersion
	winFileConfig       *FileConf
	linuxFileConfig     *FileConf
	watchLimiter        flowctrl.RateLimiter
	fullLimiter         flowctrl.RateLimiter
	metric              *hostIdentifierMetric
}

// NewHostIdentifier new HostIdentifier struct
func NewHostIdentifier(ctx context.Context, redisCli redis.Client, engine *backbone.Engine, conf *HostIdentifierConf,
	apiGWClient gse.GseClientInterface, taskClient *client.GseTaskServerClient,
	apiClient *client.GseApiServerClient, apiVersion types.ApiVersion) (*HostIdentifier, error) {
	if apiGWClient == nil && (apiClient == nil || taskClient == nil) {
		return nil, errors.New("connect to gse client is missing")
	}
	h := &HostIdentifier{
		redisCli:            redisCli,
		ctx:                 ctx,
		engine:              engine,
		gseTaskServerClient: taskClient,
		gseApiServerClient:  apiClient,
		gseApiGWClient:      apiGWClient,
		winFileConfig:       conf.WinFileConf,
		linuxFileConfig:     conf.LinuxFileConf,
		watchLimiter:        flowctrl.NewRateLimiter(conf.RateLimiter.Qps, conf.RateLimiter.Burst),
		fullLimiter:         flowctrl.NewRateLimiter(conf.RateLimiter.Qps, conf.RateLimiter.Burst),
		apiVersion:          apiVersion,
	}

	h.registerMetrics()

	return h, nil
}

// registerMetrics registers prometheus metrics.
func (h *HostIdentifier) registerMetrics() {
	getAgentStatusTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_agent_status_total", metricsNamespacePrefix),
			Help: "call gse get agent status api total.",
		},
		[]string{"status"},
	)
	h.engine.Metric().Registry().MustRegister(getAgentStatusTotal)

	pushFileTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_push_file_total", metricsNamespacePrefix),
			Help: "call gse push file api total.",
		},
		[]string{"status"},
	)
	h.engine.Metric().Registry().MustRegister(pushFileTotal)

	getResultTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_result_total", metricsNamespacePrefix),
			Help: "call gse get task result api total.",
		},
		[]string{"status"},
	)
	h.engine.Metric().Registry().MustRegister(getResultTotal)

	agentStatusTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_host_agent_status_total", metricsNamespacePrefix),
			Help: "host agent status total.",
		},
		[]string{"status"},
	)
	h.engine.Metric().Registry().MustRegister(agentStatusTotal)

	hostResultTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_host_result_total", metricsNamespacePrefix),
			Help: "host result total.",
		},
		[]string{"status"},
	)
	h.engine.Metric().Registry().MustRegister(hostResultTotal)

	h.engine.Metric().Registry().MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_fail_host_list_length", metricsNamespacePrefix),
			Help: "current length of redis host identifier fail host list.",
		},
		func() float64 {
			val, err := h.redisCli.LLen(context.Background(), RedisFailHostListName).Result()
			if err != nil {
				blog.Errorf("get redis host identifier fail host list length error, err: %v", err)
				return 0
			}
			return float64(val)
		},
	))

	h.engine.Metric().Registry().MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_task_list_length", metricsNamespacePrefix),
			Help: "current length of redis host identifier task list.",
		},
		func() float64 {
			val, err := h.redisCli.LLen(context.Background(), redisTaskListName).Result()
			if err != nil {
				blog.Errorf("get redis host identifier task list length error, err: %v", err)
				return 0
			}
			return float64(val)
		},
	))

	h.metric = &hostIdentifierMetric{
		getAgentStatusTotal: getAgentStatusTotal,
		pushFileTotal:       pushFileTotal,
		getResultTotal:      getResultTotal,
		agentStatusTotal:    agentStatusTotal,
		hostResultTotal:     hostResultTotal,
	}
}

// WatchToSyncHostIdentifier watch to sync host identifier
func (h *HostIdentifier) WatchToSyncHostIdentifier() {
	eventOp := &Event{
		engine:   h.engine,
		errFreq:  util.NewErrFrequency(nil),
		redisCli: h.redisCli,
	}
	observer := observer{
		isMaster: h.engine.Discovery(),
	}

	// start to watch and push host identifier
	for {
		preStatus, loop := observer.canLoop()
		if !loop {
			blog.V(4).Infof("loop watch host identifier, but not master, skip.")
			time.Sleep(time.Minute)
			continue
		}

		header, rid := newHeaderWithRid()
		events, lastCursor, ok := eventOp.getEvent(header, rid, preStatus)
		if !ok {
			time.Sleep(time.Second)
			continue
		}

		h.watchToSyncHostIdentifier(events, header, rid)

		eventOp.setCursor(lastCursor, rid)
	}
}

func (h *HostIdentifier) watchToSyncHostIdentifier(events []*IdentifierEvent, header http.Header, rid string) {
	// 1、查询主机状态
	statusReqList := make([]StatusReq, 0)
	for _, event := range events {
		statusReqList = append(statusReqList, StatusReq{
			CloudID:      strconv.FormatInt(event.CloudID, 10),
			InnerIP:      event.InnerIP,
			AgentID:      event.AgentID,
			BKAddressing: event.BKAddressing,
		})
	}

	resp, err := h.getAgentStatus(statusReqList, true, header, rid)
	if err != nil {
		blog.Errorf("get agent status error, host: %v, err: %v, rid: %s", events, err, rid)
		return
	}

	// 2、将处于on状态的主机拿出来构造推送信息
	task, hostInfos := h.getTaskFromEvent(events, resp, rid)
	if task == nil || (len(task.V1Task) == 0 && len(task.V2Task) == 0) {
		return
	}

	// 3、推送主机身份信息
	h.watchLimiter.AcceptMany(int64(len(hostInfos)))
	if _, err := h.pushFile(true, hostInfos, task, header, rid); err != nil {
		blog.Errorf("push host identifier to gse error, err: %v, rid: %s", err, rid)
	}
}

func (h *HostIdentifier) getTaskFromEvent(events []*IdentifierEvent, statusMap map[string]string, rid string) (
	*TaskInfo, []*HostInfo) {

	switch h.apiVersion {
	case types.V2:
		return h.getV2Task(events, statusMap, rid)

	case types.V1:
		return h.getV1Task(events, statusMap, rid)
	}

	return nil, nil
}

func (h *HostIdentifier) getV2Task(events []*IdentifierEvent, statusMap map[string]string, rid string) (*TaskInfo,
	[]*HostInfo) {

	fList := make([]*gse.Task, 0)
	hostInfos := make([]*HostInfo, 0)
	for _, event := range events {
		agentID := event.AgentID
		if agentID != "" {
			if statusMap[event.AgentID] != v2ApiAgentOnStatus {
				blog.Infof("agent status is off, agentID: %s, rid: %s", event.AgentID, rid)
				h.metric.agentStatusTotal.WithLabelValues("off").Inc()
				continue
			}
		} else {
			ids := buildV2ForStatus(strconv.FormatInt(event.CloudID, 10), event.InnerIP)
			isOn := false
			for _, id := range ids {
				if statusMap[id] == v2ApiAgentOnStatus {
					isOn = true
					agentID = id
					break
				}
			}

			if !isOn {
				blog.Infof("agent status is off, hostID: %d, ip: %s, cloudID: %d, rid: %s", event.HostID, event.InnerIP,
					event.CloudID, rid)
				h.metric.agentStatusTotal.WithLabelValues("off").Inc()
				continue
			}
		}

		if isFileExceedLimit(event.RawEvent) {
			blog.Errorf("file exceed limit: %d, unit:byte, hostID: %d, rid: %s", fileLimit, event.HostID, rid)
			continue
		}

		h.metric.agentStatusTotal.WithLabelValues("on").Inc()
		blog.Infof("agent status is on, agentID: %s, rid: %s", event.AgentID, rid)

		file := h.buildV2PushFile(event.RawEvent, agentID)
		fList = append(fList, file)
		hostInfos = append(hostInfos, &HostInfo{
			HostID:       event.HostID,
			CloudID:      event.CloudID,
			BKAddressing: event.BKAddressing,
			AgentID:      agentID,
			HostInnerIP:  event.InnerIP,
		})
	}

	taskInfo := &TaskInfo{
		V2Task: fList,
	}
	return taskInfo, hostInfos
}

func (h *HostIdentifier) getV1Task(events []*IdentifierEvent, statusMap map[string]string, rid string) (*TaskInfo,
	[]*HostInfo) {

	fList := make([]*pushfile.API_FileInfoV2, 0)
	hostInfos := make([]*HostInfo, 0)
	for _, event := range events {
		isOn, hostIP := getStatusOnAgentIP(strconv.FormatInt(event.CloudID, 10), event.InnerIP, statusMap)
		if !isOn {
			blog.Infof("agent status is off, hostID: %d, ip: %s, cloudID: %d, rid: %s",
				event.HostID, event.InnerIP, event.CloudID, rid)
			h.metric.agentStatusTotal.WithLabelValues("off").Inc()
			continue
		}

		h.metric.agentStatusTotal.WithLabelValues("on").Inc()
		blog.Infof("agent status is on, hostID: %d, ip: %s, cloudID: %d, rid: %s", event.HostID, hostIP, event.CloudID,
			rid)

		file := h.buildV1PushFile(event.RawEvent, hostIP, event.CloudID)
		fList = append(fList, file)
		hostInfos = append(hostInfos, &HostInfo{
			HostID:      event.HostID,
			HostInnerIP: hostIP,
			CloudID:     event.CloudID,
		})
	}

	taskInfo := &TaskInfo{
		V1Task: fList,
	}
	return taskInfo, hostInfos
}

// FullSyncHostIdentifier Fully synchronize host identity
func (h *HostIdentifier) FullSyncHostIdentifier() {
	start := 0
	for {
		if !h.engine.Discovery().IsMaster() {
			blog.V(4).Infof("loop full sync host identifier, but not master, skip.")
			return
		}

		header, rid := newHeaderWithRid()
		ctx, header := util.SetReadPreference(h.ctx, header, common.SecondaryPreferredMode)
		option := &metadata.ListHosts{
			Fields: []string{common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField,
				common.BKAddressingField, common.BKAgentIDField},
			Page: metadata.BasePage{
				Start: start,
				Limit: hostIdentifierBatchSyncPerLimit,
			},
		}

		hosts, err := h.engine.CoreAPI.CoreService().Host().ListHosts(ctx, header, option)
		if err != nil {
			blog.Errorf("get host in batch error, resp: %v, err: %v, rid: %s", hosts, err, rid)
			continue
		}

		if _, err := h.BatchSyncHostIdentifier(hosts.Info, header, rid); err != nil {
			blog.Errorf("full sync host identifier error, hosts: %v, err: %v, rid: %s", hosts.Info, err, rid)
		}

		start += hostIdentifierBatchSyncPerLimit
		if start >= hosts.Count {
			break
		}
	}
}

// BatchSyncHostIdentifier batch sync host identifier
func (h *HostIdentifier) BatchSyncHostIdentifier(hosts []map[string]interface{}, header http.Header,
	rid string) (*Task, error) {

	if len(hosts) == 0 {
		return nil, errors.New("the hosts count is 0")
	}

	// 1、查询主机状态
	statusReqList := make([]StatusReq, 0)
	for _, hostInfo := range hosts {
		statusReqList = append(statusReqList, StatusReq{
			CloudID:      util.GetStrByInterface(hostInfo[common.BKCloudIDField]),
			InnerIP:      util.GetStrByInterface(hostInfo[common.BKHostInnerIPField]),
			AgentID:      util.GetStrByInterface(hostInfo[common.BKAgentIDField]),
			BKAddressing: util.GetStrByInterface(hostInfo[common.BKAddressingField]),
		})
	}

	resp, err := h.getAgentStatus(statusReqList, false, header, rid)
	if err != nil {
		blog.Errorf("get agent status error,  hostInfo: %v, err: %v, rid: %s", hosts, err, rid)
		return nil, err
	}

	// 2、将处于on状态的主机拿出来构造主机身份推送信息
	hostIDs, hostInfos, hostMap := h.getOnStatusAgent(hosts, resp, rid)
	if len(hostIDs) == 0 {
		return nil, errors.New("the host agent status is off")
	}

	// 3、查询主机身份并推送
	return h.getHostIdentifierAndPush(hostIDs, hostMap, hostInfos, rid, header)
}

func (h *HostIdentifier) getOnStatusAgent(hosts []map[string]interface{}, statusMap map[string]string, rid string) (
	[]int64, []*HostInfo, map[int64]string) {

	switch h.apiVersion {
	case types.V2:
		return h.getV2OnStatusAgent(hosts, statusMap, rid)

	case types.V1:
		return h.getV1OnStatusAgent(hosts, statusMap, rid)
	}

	return nil, nil, nil
}

func (h *HostIdentifier) getV2OnStatusAgent(hosts []map[string]interface{}, statusMap map[string]string, rid string) (
	[]int64, []*HostInfo, map[int64]string) {

	hostIDs := make([]int64, 0)
	hostInfos := make([]*HostInfo, 0)
	// 此map保存hostID和该host处于on的agent的agentID的对应关系
	hostMap := make(map[int64]string)
	for _, hostInfo := range hosts {
		hostID, err := util.GetInt64ByInterface(hostInfo[common.BKHostIDField])
		if err != nil {
			blog.Errorf("get hostID error, hostInfo: %v, error: %v, rid: %s", hostInfo, err, rid)
			continue
		}

		cloudID, err := util.GetInt64ByInterface(hostInfo[common.BKCloudIDField])
		if err != nil {
			blog.Errorf("get cloudID error, hostInfo: %v, error: %v, rid: %s", hostInfo, err, rid)
			continue
		}
		innerIP := util.GetStrByInterface(hostInfo[common.BKHostInnerIPField])

		agentID := util.GetStrByInterface(hostInfo[common.BKAgentIDField])
		if agentID != "" {
			if statusMap[agentID] != v2ApiAgentOnStatus {
				blog.Infof("agent status is off, agentID: %s, hostID: %d, rid: %s", agentID, hostID, rid)
				h.metric.agentStatusTotal.WithLabelValues("off").Inc()
				continue
			}
		} else {
			ids := buildV2ForStatus(strconv.FormatInt(cloudID, 10), innerIP)
			isOn := false
			for _, id := range ids {
				if statusMap[id] == v2ApiAgentOnStatus {
					isOn = true
					agentID = id
				}
			}

			if !isOn {
				blog.Infof("agent status is off, hostID: %d, ip: %s, cloudID: %d, rid: %s", hostID, innerIP, cloudID,
					rid)
				h.metric.agentStatusTotal.WithLabelValues("off").Inc()
				continue
			}
		}

		h.metric.agentStatusTotal.WithLabelValues("on").Inc()
		hostIDs = append(hostIDs, hostID)
		hostMap[hostID] = agentID
		addressing := util.GetStrByInterface(hostInfo[common.BKAddressingField])
		hostInfos = append(hostInfos, &HostInfo{
			HostID:       hostID,
			AgentID:      agentID,
			BKAddressing: addressing,
			HostInnerIP:  innerIP,
			CloudID:      cloudID,
		})
	}

	return hostIDs, hostInfos, hostMap
}

func (h *HostIdentifier) getV1OnStatusAgent(hosts []map[string]interface{}, statusMap map[string]string, rid string) (
	[]int64, []*HostInfo, map[int64]string) {

	hostIDs := make([]int64, 0)
	// 此map保存hostID和该host处于on的agent的ip的对应关系
	hostMap := make(map[int64]string)
	hostInfos := make([]*HostInfo, 0)
	for _, hostInfo := range hosts {
		hostID, err := util.GetInt64ByInterface(hostInfo[common.BKHostIDField])
		if err != nil {
			blog.Errorf("get hostID error, hostInfo: %v, error: %v, rid: %s", hostInfo, err, rid)
			continue
		}

		cloudID, err := util.GetInt64ByInterface(hostInfo[common.BKCloudIDField])
		if err != nil {
			blog.Errorf("get cloudID error, hostInfo: %v, error: %v, rid: %s", hostInfo, err, rid)
			continue
		}

		innerIP := util.GetStrByInterface(hostInfo[common.BKHostInnerIPField])
		isOn, hostIP := getStatusOnAgentIP(strconv.FormatInt(cloudID, 10), innerIP, statusMap)
		if !isOn {
			blog.Infof("agent status is off, hostID: %d, ip: %s, cloudID: %d, rid: %s", hostID, innerIP, cloudID, rid)
			h.metric.agentStatusTotal.WithLabelValues("off").Inc()
			continue
		}

		h.metric.agentStatusTotal.WithLabelValues("on").Inc()
		hostIDs = append(hostIDs, hostID)
		hostMap[hostID] = hostIP
		hostInfos = append(hostInfos, &HostInfo{
			HostID:      hostID,
			HostInnerIP: hostIP,
			CloudID:     cloudID,
		})
	}

	return hostIDs, hostInfos, hostMap
}

func (h *HostIdentifier) getAgentStatus(statusReqList []StatusReq, always bool, header http.Header, rid string) (
	map[string]string, error) {

	if len(statusReqList) == 0 {
		return nil, errors.New("agent request array is empty")
	}

	switch h.apiVersion {
	case types.V2:
		result, err := h.getAgentStatusByV2Api(statusReqList, always, header)
		if err != nil {
			blog.Errorf("get host agent status error, err: %v, rid: %s", err, rid)
			return nil, err
		}
		h.metric.getAgentStatusTotal.WithLabelValues("success").Inc()
		return result, nil

	case types.V1:
		result, err := h.getAgentStatusByV1Api(statusReqList, always, rid)
		if err != nil {
			blog.Errorf("get host agent status error, err: %v, rid: %s", err, rid)
			return nil, err
		}
		h.metric.getAgentStatusTotal.WithLabelValues("success").Inc()
		return result, nil
	}

	return nil, errors.New("can not find api about get agent status")
}

func (h *HostIdentifier) getAgentStatusByV2Api(statusReqList []StatusReq, always bool, header http.Header) (
	map[string]string, error) {

	// build agentID request list
	agentIDList := make([]string, 0)
	for _, statusReq := range statusReqList {
		if statusReq.AgentID != "" {
			agentIDList = append(agentIDList, statusReq.AgentID)
			continue
		}

		if statusReq.BKAddressing == common.BKAddressingStatic {
			list := buildV2ForStatus(statusReq.CloudID, statusReq.InnerIP)
			agentIDList = append(agentIDList, list...)
		}
	}

	if len(agentIDList) == 0 {
		return nil, errors.New("can not build agentID info to find agent status")
	}

	// find agent status
	req := &gse.ListAgentStateRequest{
		AgentIDList: agentIDList,
	}
	statusMap := make(map[string]string)
	failCount := 0
	var err error
	var resp *gse.ListAgentStateResp
	for always || failCount < retryTimes {
		resp, err = h.gseApiGWClient.ListAgentState(h.ctx, header, req)
		if err != nil {
			h.metric.getAgentStatusTotal.WithLabelValues("failed").Inc()
			failCount++
			sleepForFail(failCount)
			continue
		}

		for _, agentStatus := range resp.Data {
			statusMap[agentStatus.BKAgentID] = strconv.Itoa(agentStatus.StatusCode)
		}
		break
	}

	if !always && failCount >= retryTimes {
		return nil, err
	}

	return statusMap, nil
}

func (h *HostIdentifier) getAgentStatusByV1Api(statusReqList []StatusReq, always bool, rid string) (map[string]string,
	error) {

	req := new(getstatus.AgentStatusRequest)
	for _, statusReq := range statusReqList {
		info := buildForStatus(statusReq.CloudID, statusReq.InnerIP)
		req.Hosts = append(req.Hosts, info...)
	}

	var err error
	failCount := 0
	resp := new(getstatus.AgentStatusResponse)
	// 调用gse api server 查询agent状态
	for always || failCount < retryTimes {
		resp, err = h.gseApiServerClient.GetAgentStatus(context.Background(), req)
		if err != nil {
			blog.Errorf("get host agent status error, err: %v, rid: %s", err, rid)
			h.metric.getAgentStatusTotal.WithLabelValues("failed").Inc()
			failCount++
			sleepForFail(failCount)
			continue
		}

		if resp.BkErrorCode != common.CCSuccess {
			blog.Errorf("get agent status fail, code: %d, msg: %s, rid: %s", resp.BkErrorCode, resp.BkErrorMsg, rid)
			h.metric.getAgentStatusTotal.WithLabelValues("failed").Inc()
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}

	if !always && failCount >= retryTimes {
		return nil, errors.New("find agent status from apiServer error")
	}

	return resp.Result_, nil
}

func (h *HostIdentifier) buildV2PushFile(hostIdentifier, agentID string) *gse.Task {
	osType := gjson.Get(hostIdentifier, common.BKOSTypeField).String()
	conf := h.getHostIdentifierFileConf(osType)
	fileOwner := conf.FileOwner
	// 如果是window系统，并且通过cloudID:ip的方式下发，代表是window的安装了1.0agent的机器，需要设置成特殊的owner
	if osType == common.HostOSTypeEnumWindows && strings.Contains(agentID, ":") {
		fileOwner = v1AgentFileOwner
	}

	return &gse.Task{
		FileName:    conf.FileName,
		StoreDir:    conf.FilePath,
		FileContent: hostIdentifier,
		Owner:       fileOwner,
		Right:       conf.FilePrivilege,
		AgentIDList: []string{agentID},
	}
}

// v1AgentFileOwner 由于安装agent1.0的window系统，下发主机身份时，文件的所有者只能是root，所以涉及到的地方需要进行特殊处理，配置成该值
const v1AgentFileOwner = "root"

func (h *HostIdentifier) buildV1PushFile(hostIdentifier, hostIP string, cloudID int64) *pushfile.API_FileInfoV2 {
	fileInfo := &pushfile.API_FileInfoV2{
		MFile: &pushfile.API_BaseFileInfo{
			MMd5: strMd5(hostIdentifier),
		},
		MHostlist: []*pushfile.API_Host{
			{
				MIP:         hostIP,
				MBusinessid: int32(cloudID),
			},
		},
		MContent: hostIdentifier,
		MCaller: map[string]string{
			"CALLER_NAME": callerName,
			"CALLER_IP":   h.engine.GetSrvInfo().IP,
		},
	}

	osType := gjson.Get(hostIdentifier, common.BKOSTypeField).String()
	conf := h.getHostIdentifierFileConf(osType)
	fileInfo.MFile.MName = conf.FileName
	fileInfo.MFile.MPath = conf.FilePath
	fileInfo.MFile.MOwner = conf.FileOwner
	// thrift的接口只会通过cloudID:innerIP的方式路由到1.0的agent，所以对于window操作系统可以直接设置
	if osType == common.HostOSTypeEnumWindows {
		fileInfo.MFile.MOwner = v1AgentFileOwner
	}
	fileInfo.MFile.MRight = conf.FilePrivilege
	fileInfo.MFile.MBackupName = conf.FileName + ".bak"

	return fileInfo
}

func (h *HostIdentifier) getHostIdentifierAndPush(hostIDs []int64, hostMap map[int64]string,
	hostInfos []*HostInfo, rid string, header http.Header) (*Task, error) {
	if len(hostIDs) == 0 {
		blog.Errorf("hostIDs count is 0, rid: %s", rid)
		return nil, errors.New("hostIDs count is 0")
	}
	hostInfoMap := make(map[int64]*HostInfo)
	for _, hostInfo := range hostInfos {
		if hostInfo == nil {
			continue
		}
		hostInfoMap[hostInfo.HostID] = hostInfo
	}

	// 1、查询主机身份
	ctx, header := util.SetReadPreference(h.ctx, header, common.SecondaryPreferredMode)
	queryHostIdentifier := &metadata.SearchHostIdentifierParam{HostIDs: hostIDs}
	rsp, err := h.engine.CoreAPI.CoreService().Host().FindIdentifier(ctx, header, queryHostIdentifier)
	if err != nil {
		blog.Errorf("find host identifier error, hostIDs: %v, err: %v, rid: %s", hostIDs, err, rid)
		return nil, err
	}
	if rsp.Count == 0 {
		blog.Errorf("can not find host identifier, hostIDs: %v, rid: %s", hostIDs, rid)
		return nil, errors.New("can not find host identifier")
	}

	// 2、构造想要推送的主机身份文件信息
	fileList := make([]*pushfile.API_FileInfoV2, 0)
	file2List := make([]*gse.Task, 0)
	hosts := make([]*HostInfo, 0)
	for _, identifier := range rsp.Info {
		hostIdentifier, err := json.Marshal(identifier)
		if err != nil {
			blog.Errorf("marshal host identifier failed, val: %v, err: %v, rid: %s", identifier, err, rid)
			continue
		}

		switch h.apiVersion {
		case types.V1:
			fileList = append(fileList, h.buildV1PushFile(string(hostIdentifier), hostMap[identifier.HostID],
				identifier.CloudID))
			hosts = append(hosts, hostInfoMap[identifier.HostID])
			continue

		case types.V2:
			if isFileExceedLimit(string(hostIdentifier)) {
				blog.Errorf("file exceed limit: %d, unit:byte, hostID: %d, rid: %s", fileLimit, identifier.HostID, rid)
				continue
			}
			file := h.buildV2PushFile(string(hostIdentifier), hostMap[identifier.HostID])
			file2List = append(file2List, file)
			hosts = append(hosts, hostInfoMap[identifier.HostID])
		}
	}
	switch h.apiVersion {
	case types.V1:
		if len(fileList) == 0 {
			blog.Errorf("get identifier success, but can not build file to push, hostID: %v, rid: %s", hostIDs, rid)
			return nil, errors.New("get host identifier success, but can not build file to push")
		}

	case types.V2:
		if len(file2List) == 0 {
			blog.Errorf("get identifier success, but can not build file to push, hostID: %v, rid: %s", hostIDs, rid)
			return nil, errors.New("get host identifier success, but can not build file to push")
		}
	}

	// 3、推送主机身份信息
	task := new(TaskInfo)
	switch h.apiVersion {
	case types.V1:
		h.fullLimiter.AcceptMany(int64(len(fileList)))
		task.V1Task = fileList

	case types.V2:
		h.fullLimiter.AcceptMany(int64(len(file2List)))
		task.V2Task = file2List
	}

	return h.pushFile(false, hosts, task, header, rid)
}
