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
	"fmt"
	"strconv"
	"time"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/gse/client"
	"configcenter/src/thirdparty/gse/get_agent_state_forsyncdata"
	"configcenter/src/thirdparty/gse/push_file_forsyncdata"

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
)

// HostIdentifier manipulate the structure of the host Identifier
type HostIdentifier struct {
	redisCli            redis.Client
	engine              *backbone.Engine
	ctx                 context.Context
	gseTaskServerClient *client.GseTaskServerClient
	gseApiServerClient  *client.GseApiServerClient
	winFileConfig       *FileConf
	linuxFileConfig     *FileConf
	watchLimiter        flowctrl.RateLimiter
	fullLimiter         flowctrl.RateLimiter
}

// NewHostIdentifier new HostIdentifier struct
func NewHostIdentifier(ctx context.Context, redisCli redis.Client, engine *backbone.Engine,
	taskClient *client.GseTaskServerClient, apiClient *client.GseApiServerClient) (*HostIdentifier, error) {
	winFileConfig, err := newHostIdentifierFileConf("windows")
	if err != nil {
		return nil, err
	}

	linuxFileConfig, err := newHostIdentifierFileConf("linux")
	if err != nil {
		return nil, err
	}

	qps, burst, err := getRateLimiterConfig()
	if err != nil {
		return nil, err
	}

	h := &HostIdentifier{
		redisCli:            redisCli,
		ctx:                 ctx,
		engine:              engine,
		gseTaskServerClient: taskClient,
		gseApiServerClient:  apiClient,
		winFileConfig:       winFileConfig,
		linuxFileConfig:     linuxFileConfig,
		watchLimiter:        flowctrl.NewRateLimiter(qps, burst),
		fullLimiter:         flowctrl.NewRateLimiter(qps, burst),
	}
	return h, nil
}

// WatchToSyncHostIdentifier watch to sync host identifier
func (h *HostIdentifier) WatchToSyncHostIdentifier() {
	errFrequency := util.NewErrFrequency(nil)
	options := &watch.WatchEventOptions{
		EventTypes: []watch.EventType{watch.Create, watch.Update},
		Resource:   watch.HostIdentifier,
	}
	// 用于判断之前的状态是从节点，此值当第一次从slave切换到master时，会去查询redis里的cursor
	isBeforeStatusSlave := true
	// start to watch and push host identifier
	for {
		if !h.engine.Discovery().IsMaster() {
			isBeforeStatusSlave = true
			time.Sleep(time.Minute)
			continue
		}

		if isBeforeStatusSlave {
			options.Cursor = ""
			// 从redis里拿cursor，否则从当前时间watch
			cursor, err := h.redisCli.Get(h.ctx, hostIdentifierCursor).Result()
			if err != nil {
				blog.Errorf("get host identity cursor from redis error, err: %v", err)
			} else {
				options.Cursor = cursor
			}
			isBeforeStatusSlave = false
		}

		header, rid := newHeaderWithRid()
		watchEvents, watchErr := h.engine.CoreAPI.CacheService().Cache().Event().WatchEvent(h.ctx, header, options)
		if watchErr != nil {
			if watchErr.GetCode() == common.CCErrEventChainNodeNotExist {
				// 设置从当前时间开始watch
				options.Cursor = ""
				if err := h.redisCli.Del(h.ctx, hostIdentifierCursor).Err(); err != nil {
					blog.Errorf("delete redis key failed, key: %s, err: %v, rid: %s", hostIdentifierCursor, err, rid)
				}
			}

			if errFrequency.IsErrAlwaysAppear(watchErr) {
				options.Cursor = ""
				if err := h.redisCli.Del(h.ctx, hostIdentifierCursor).Err(); err != nil {
					blog.Errorf("delete redis key failed, key: %s, err: %v, rid: %s", hostIdentifierCursor, err, rid)
				}
			}

			blog.Errorf("watch host_identifier event failed, reset watch it from current time, err: %v, rid: %s",
				watchErr, rid)
			time.Sleep(time.Second)
			continue
		}

		if !gjson.Get(*watchEvents, "bk_watched").Bool() {
			options.Cursor = gjson.Get(*watchEvents, "bk_events.0.bk_cursor").String()
			continue
		}
		events := gjson.Get(*watchEvents, "bk_events").Array()
		if len(events) == 0 {
			blog.Errorf("can not get events from bk_events, watchEvent: %v, rid: %s", *watchEvents, rid)
			continue
		}

		h.watchToSyncHostIdentifier(events)

		// 保存新的cursor到内存和redis中
		redisFailCount := 0
		options.Cursor = events[len(events)-1].Get("bk_cursor").String()
		for redisFailCount < retryTimes {
			if err := h.redisCli.Set(h.ctx, hostIdentifierCursor, options.Cursor, 3*time.Hour).Err(); err != nil {
				blog.Errorf("set redis error, key: %s, val: %s, err: %v", hostIdentifierCursor, options.Cursor, err)
				redisFailCount++
				sleepForFail(redisFailCount)
				continue
			}
			break
		}
	}
}

func (h *HostIdentifier) watchToSyncHostIdentifier(events []gjson.Result) {
	// 查询主机状态处于on还是off
	status := new(get_agent_state_forsyncdata.AgentStatusRequest)
	for _, event := range events {
		eventDetail := event.Map()["bk_detail"]
		if !eventDetail.Exists() {
			continue
		}

		eventDetailMap := eventDetail.Map()
		if !eventDetailMap[common.BKCloudIDField].Exists() || !eventDetailMap[common.BKHostInnerIPField].Exists() {
			continue
		}

		status.Hosts = append(status.Hosts,
			buildAgentStatusRequestHostInfo(eventDetailMap[common.BKCloudIDField].String(),
				eventDetailMap[common.BKHostInnerIPField].String())...)
	}

	gseFailCount := 0
	agentStatus := new(get_agent_state_forsyncdata.AgentStatusResponse)
	var err error
	for {
		agentStatus, err = h.gseApiServerClient.GetAgentStatus(h.ctx, status)
		if err != nil {
			blog.Errorf("get host agent status error, err: %v", err)
			gseFailCount++
			sleepForFail(gseFailCount)
			continue
		}

		if agentStatus.BkErrorCode != common.CCSuccess {
			blog.Errorf("get host agent status fail, errCode: %d, errMessage: %s",
				agentStatus.BkErrorCode, agentStatus.BkErrorMsg)
			gseFailCount++
			sleepForFail(gseFailCount)
			continue
		}
		break
	}

	// 将处于on状态的主机拿出来构造推送信息
	fileList := make([]*push_file_forsyncdata.API_FileInfoV2, 0)
	hostInfos := make([]*HostInfo, 0)
	for _, event := range events {
		eventDetail := event.Map()["bk_detail"]
		if !eventDetail.Exists() {
			continue
		}

		eventDetailMap := eventDetail.Map()
		if !eventDetailMap[common.BKCloudIDField].Exists() || !eventDetailMap[common.BKHostInnerIPField].Exists() {
			continue
		}

		isOn, hostIP := getStatusOnAgentIP(eventDetailMap[common.BKCloudIDField].String(),
			eventDetailMap[common.BKHostInnerIPField].String(), agentStatus.Result_)
		if !isOn {
			continue
		}
		fileList = append(fileList, h.buildPushFile(eventDetail.String(),
			hostIP, eventDetailMap[common.BKCloudIDField].Int()))
		hostInfos = append(hostInfos, &HostInfo{
			HostID:      eventDetailMap[common.BKHostIDField].Int(),
			HostInnerIP: hostIP,
			CloudID:     eventDetailMap[common.BKCloudIDField].Int(),
		})
	}

	if len(fileList) == 0 {
		return
	}

	h.watchLimiter.AcceptMany(int64(len(fileList)))
	// 推送主机身份信息
	if _, err := h.pushFile(true, hostInfos, fileList); err != nil {
		blog.Errorf("push host identifier to gse error, err: %v", err)
	}
}

// FullSyncHostIdentifier Fully synchronize host identity
func (h *HostIdentifier) FullSyncHostIdentifier() {
	start := 0
	for {
		if !h.engine.Discovery().IsMaster() {
			return
		}

		header, rid := newHeaderWithRid()
		util.SetReadPreference(h.ctx, header, common.SecondaryPreferredMode)
		option := &metadata.ListHosts{
			Fields: []string{common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField},
			Page: metadata.BasePage{
				Start: start,
				Limit: hostIdentifierBatchSyncPerLimit,
			},
		}
		hosts, err := h.engine.CoreAPI.CoreService().Host().ListHosts(h.ctx, header, option)
		if err != nil {
			blog.Errorf("get host in batch error, resp: %v, err: %v, rid: %s", hosts, err, rid)
			continue
		}

		if _, err := h.BatchSyncHostIdentifier(hosts.Info); err != nil {
			blog.Errorf("full sync host identifier error, hosts: %v, err: %v, rid: %s", hosts.Info, err, rid)
		}

		start += hostIdentifierBatchSyncPerLimit
		if start >= hosts.Count {
			break
		}
	}
}

// BatchSyncHostIdentifier batch sync host identifier
func (h *HostIdentifier) BatchSyncHostIdentifier(hosts []map[string]interface{}) (*Task, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("the hosts count is 0")
	}

	var err error
	agentStatus := new(get_agent_state_forsyncdata.AgentStatusResponse)
	failCount := 0
	for failCount < retryTimes {
		// 查询主机状态处于on还是off
		agentStatusRequest := &get_agent_state_forsyncdata.AgentStatusRequest{}
		for _, hostInfo := range hosts {
			agentStatusRequest.Hosts = append(agentStatusRequest.Hosts,
				buildAgentStatusRequestHostInfo(util.GetStrByInterface(hostInfo[common.BKCloudIDField]),
					util.GetStrByInterface(hostInfo[common.BKHostInnerIPField]))...)
		}

		agentStatus, err = h.gseApiServerClient.GetAgentStatus(h.ctx, agentStatusRequest)
		if err != nil {
			blog.Errorf("get host agent status error: err: %v", err)
			failCount++
			sleepForFail(failCount)
			continue
		}

		if agentStatus.BkErrorCode != common.CCSuccess {
			blog.Errorf("get host agent status fail, errCode: %d, errMessage: %s",
				agentStatus.BkErrorCode, agentStatus.BkErrorMsg)
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}

	if failCount >= retryTimes {
		return nil, err
	}

	// 将处于on状态的主机拿出来构造主机身份推送信息
	hostIDs := make([]int64, 0)
	hostInfos := make([]*HostInfo, 0)
	// 此map保存hostID和该host处于on的agent的ip的对应关系
	hostMap := make(map[int64]string)
	for _, hostInfo := range hosts {
		hostID, err := util.GetInt64ByInterface(hostInfo[common.BKHostIDField])
		if err != nil {
			blog.Errorf("get hostID error, hostInfo: %v, error: %v", hostInfo, err)
			continue
		}
		cloudID, err := util.GetInt64ByInterface(hostInfo[common.BKCloudIDField])
		if err != nil {
			blog.Errorf("get cloudID error, hostInfo: %v, error: %v", hostInfo, err)
			continue
		}
		innerIP := util.GetStrByInterface(hostInfo[common.BKHostInnerIPField])
		isOn, hostIP := getStatusOnAgentIP(strconv.FormatInt(cloudID, 10), innerIP, agentStatus.Result_)
		if !isOn {
			continue
		}
		hostIDs = append(hostIDs, hostID)
		hostMap[hostID] = hostIP
		hostInfos = append(hostInfos, &HostInfo{
			HostID:      hostID,
			HostInnerIP: hostIP,
			CloudID:     cloudID,
		})
	}

	if len(hostIDs) == 0 {
		return nil, fmt.Errorf("the host agent status is off")
	}

	return h.getHostIdentifierAndPush(hostIDs, hostMap, hostInfos)
}

func (h *HostIdentifier) buildPushFile(hostIdentifier, hostIP string,
	cloudID int64) *push_file_forsyncdata.API_FileInfoV2 {
	fileInfo := &push_file_forsyncdata.API_FileInfoV2{
		MFile: &push_file_forsyncdata.API_BaseFileInfo{
			MMd5: strMd5(hostIdentifier),
		},
		MHostlist: []*push_file_forsyncdata.API_Host{
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
	fileInfo.MFile.MRight = conf.FilePrivilege

	return fileInfo
}

func (h *HostIdentifier) getHostIdentifierAndPush(hostIDs []int64, hostMap map[int64]string,
	hostInfos []*HostInfo) (*Task, error) {
	// 查询主机身份
	header, rid := newHeaderWithRid()
	util.SetReadPreference(h.ctx, header, common.SecondaryPreferredMode)
	queryHostIdentifier := &metadata.SearchHostIdentifierParam{HostIDs: hostIDs}
	rsp, err := h.engine.CoreAPI.CoreService().Host().FindIdentifier(h.ctx, header, queryHostIdentifier)
	if err != nil {
		blog.Errorf("find host identifier error, hostIDs: %v, err: %v, rid: %s", hostIDs, err, rid)
		return nil, err
	}
	if rsp.Count == 0 {
		blog.Errorf("can not find host identifier, hostIDs: %v, rid: %s", hostIDs, rid)
		return nil, fmt.Errorf("can not find host identifier")
	}

	// 构造想要推送的主机身份文件信息
	fileList := make([]*push_file_forsyncdata.API_FileInfoV2, 0)
	for _, identifier := range rsp.Info {
		hostIdentifier, err := json.Marshal(identifier)
		if err != nil {
			blog.Errorf("marshal host identifier failed, val: %v, err: %v, rid: %s", identifier, err, rid)
			continue
		}
		fileList = append(fileList, h.buildPushFile(string(hostIdentifier), hostMap[identifier.HostID],
			identifier.CloudID))
	}

	if len(fileList) == 0 {
		blog.Errorf("get host identifier success, but can not build file to push, hostID: %v, rid: %s", hostIDs, rid)
		return nil, fmt.Errorf("get host identifier success, but can not build file to push")
	}

	h.fullLimiter.AcceptMany(int64(len(fileList)))
	// 推送主机身份信息
	return h.pushFile(false, hostInfos, fileList)
}
