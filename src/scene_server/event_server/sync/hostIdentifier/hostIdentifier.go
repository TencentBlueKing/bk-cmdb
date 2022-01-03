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

package hostIdentifier

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
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
	// defaultHostIdentifierBatchSyncPerLimit default batch sync limit for host identifier
	defaultHostIdentifierBatchSyncPerLimit = 500
)

// HostIdentifier manipulate the structure of the host Identifier
type HostIdentifier struct {
	redisCli            redis.Client
	engine              *backbone.Engine
	ctx                 context.Context
	gseTaskServerClient *client.GseTaskServerClient
	gseApiServerClient  *client.GseApiServerClient
}

// NewHostIdentifier new HostIdentifier struct
func NewHostIdentifier(ctx context.Context, redisCli redis.Client, engine *backbone.Engine,
	gseTaskServerClient *client.GseTaskServerClient,
	gseApiServerClient *client.GseApiServerClient) *HostIdentifier {
	h := &HostIdentifier{
		redisCli:            redisCli,
		ctx:                 ctx,
		engine:              engine,
		gseTaskServerClient: gseTaskServerClient,
		gseApiServerClient:  gseApiServerClient,
	}
	return h
}

// WatchToSyncHostIdentifier watch to sync host identifier
func (h *HostIdentifier) WatchToSyncHostIdentifier() {
	var err error
	options := &watch.WatchEventOptions{
		EventTypes: []watch.EventType{watch.Create, watch.Update},
		Resource:   watch.HostIdentifier,
	}
	// 从redis里拿cursor，否则从当前时间watch
	cursor, err := h.redisCli.Get(h.ctx, hostIdentifierCursor).Result()
	if err != nil {
		blog.Errorf("get host identity cursor from redis error, err: %v", err)
	} else {
		options.Cursor = cursor
	}

	// start to watch and push host identifier
	for {
		if !h.engine.Discovery().IsMaster() {
			return
		}

		header, rid := newHeaderWithRid()
		watchEvents, watchErr := h.engine.CoreAPI.CacheService().Cache().Event().WatchEvent(h.ctx, header, options)
		if watchErr != nil && watchErr.GetCode() == common.CCErrEventChainNodeNotExist {
			// 设置从当前时间开始watch
			options.Cursor = ""
			blog.Errorf("watch host_identifier event error, err: %v, rid: %s", err, rid)
			continue
		}
		if watchErr != nil {
			blog.Errorf("watch host_identifier event error, err: %v, rid: %s", err, rid)
			time.Sleep(time.Second)
			continue
		}
		if !gjson.Get(*watchEvents, "bk_watched").Bool() {
			options.Cursor = gjson.Get(*watchEvents, "bk_events.0.bk_cursor").String()
			continue
		}
		events := gjson.Get(*watchEvents, "bk_events")

		h.watchToSyncHostIdentifier(events)

		// 保存新的cursor到内存和redis中
		redisFailCount := 0
		options.Cursor = gjson.Get(*watchEvents,
			"bk_events."+strconv.Itoa(len(events.Array())-1)+".bk_cursor").String()
		for redisFailCount < retryTimes {
			if err := h.redisCli.Set(h.ctx, hostIdentifierCursor, options.Cursor, 3*time.Hour).Err(); err != nil {
				blog.Errorf("set redis key: %s val: %s error, err: %v", hostIdentifierCursor,
					options.Cursor, err)
				redisFailCount++
				sleepForFail(redisFailCount)
				continue
			}
			break
		}
	}
}

func (h *HostIdentifier) watchToSyncHostIdentifier(events gjson.Result) {
	// 查询主机状态
	agentStatusRequest := new(get_agent_state_forsyncdata.AgentStatusRequest)
	for _, event := range events.Array() {
		eventDetail := event.Map()["bk_detail"].Map()
		agentStatusRequest.Hosts = append(agentStatusRequest.Hosts,
			buildAgentStatusRequestHostInfo(eventDetail[common.BKCloudIDField].String(),
				eventDetail[common.BKHostInnerIPField].String())...)
	}
	gseFailCount := 0
	agentStatus := new(get_agent_state_forsyncdata.AgentStatusResponse)
	var err error
	for {
		agentStatus, err = h.gseApiServerClient.GetAgentStatus(h.ctx, agentStatusRequest)
		if err != nil || agentStatus.BkErrorCode != common.CCSuccess {
			blog.Errorf("get host agent status error: err: %v, errCode: %d, errMessage: %s",
				err, agentStatus.BkErrorCode, agentStatus.BkErrorMsg)
			gseFailCount++
			sleepForFail(gseFailCount)
			continue
		}
		break
	}

	// 将处于on状态的主机拿出来构造推送信息
	var fileList []*push_file_forsyncdata.API_FileInfoV2
	var hostInfos []*HostInfo
	for _, event := range events.Array() {
		eventDetail := event.Map()["bk_detail"].Map()
		isOn, hostIP := getAgentIPStatusIsOn(eventDetail[common.BKCloudIDField].String(),
			eventDetail[common.BKHostInnerIPField].String(), agentStatus.Result_)
		if isOn {
			fileList = append(fileList, h.buildPushFile(event.Map()["bk_detail"].String(),
				hostIP, eventDetail[common.BKCloudIDField].Int()))
			hostInfos = append(hostInfos, &HostInfo{
				HostInnerIP: hostIP,
				CloudID:     eventDetail[common.BKCloudIDField].Int(),
				Times:       retryTimes,
			})
		}
	}

	if len(fileList) == 0 {
		return
	}

	// 推送主机身份信息
	gseFailCount = 0
	pushFileResp := new(push_file_forsyncdata.API_CommRsp)
	for {
		pushFileResp, err = h.gseTaskServerClient.PushFileV2(h.ctx, fileList)
		if err != nil || pushFileResp.MErrcode != common.CCSuccess {
			blog.Errorf("push host identifier to gse error: err: %v, errCode: %d, errMessage: %s",
				err, pushFileResp.MErrcode, pushFileResp.MErrmsg)
			gseFailCount++
			sleepForFail(gseFailCount)
			continue
		}
		break
	}
	blog.Infof("push host identifier success: file: %v", fileList)

	// 将任务放到redis维护的任务list中
	if err := h.createNewTaskToTaskList(pushFileResp.MContent, hostInfos); err != nil {
		blog.Errorf("add task to redis task list error, taskID: %s, hostInfos: %v, err: %v",
			pushFileResp.MContent, hostInfos, err)
	}
}

// BatchSyncHostIdentifier batch sync host identifier
func (h *HostIdentifier) BatchSyncHostIdentifier() {
	start := 0
	limit, err := cc.Int("eventServer.hostIdentifier.batchSyncPerLimit")
	if err != nil {
		blog.Errorf("get host identifier batchSyncPerLimit config error, err: %v", err)
		limit = defaultHostIdentifierBatchSyncPerLimit
	}
	for {
		if !h.engine.Discovery().IsMaster() {
			return
		}

		header, rid := newHeaderWithRid()
		params := metadata.ListHostsWithNoBizParameter{
			Fields: []string{common.BKHostInnerIPField, common.BKCloudIDField},
			Page: metadata.BasePage{
				Start: start,
				Limit: limit,
			},
		}
		hosts, err := h.engine.CoreAPI.ApiServer().ListHostWithoutApp(h.ctx, header, params)
		if err != nil || !hosts.Result {
			blog.Errorf("get host in batch error, resp: %v, err: %v, rid: %s", hosts, err, rid)
			continue
		}
		if len(hosts.Data.Info) == 0 {
			break
		}

		h.batchSyncHostIdentifier(hosts.Data.Info)

		start++
		if start*limit >= hosts.Data.Count {
			break
		}
	}
}

func (h *HostIdentifier) batchSyncHostIdentifier(hosts []map[string]interface{}) {
	header, rid := newHeaderWithRid()
	agentStatus := new(get_agent_state_forsyncdata.AgentStatusResponse)
	var err error
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
		if err != nil || agentStatus.BkErrorCode != common.CCSuccess {
			blog.Errorf("get host agent status error: err: %v, errCode: %d, errMessage: %s",
				err, agentStatus.BkErrorCode, agentStatus.BkErrorMsg)
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}
	if failCount >= retryTimes {
		return
	}

	// 将处于on状态的主机拿出来构造主机身份推送信息
	var fileList []*push_file_forsyncdata.API_FileInfoV2
	var hostInfos []*HostInfo
	for _, hostInfo := range hosts {
		cloudID := util.GetStrByInterface(hostInfo[common.BKCloudIDField])
		innerIP := util.GetStrByInterface(hostInfo[common.BKHostInnerIPField])
		cloudIDInt64Val, err := strconv.ParseInt(cloudID, 10, 64)
		if err != nil {
			blog.Errorf("convert cloudID string to int64 error, hostInfo: %v, error: %v", hostInfo, err)
			continue
		}
		isOn, hostIP := getAgentIPStatusIsOn(cloudID, innerIP, agentStatus.Result_)
		if isOn {
			identifier, err := h.findHostIdentifier(cloudIDInt64Val, innerIP, header)
			if err != nil {
				blog.Errorf("search identifier error, cloudID: %d, innerIP: %s, err: %v, rid: %s",
					cloudID, innerIP, err, header, rid)
				continue
			}
			fileList = append(fileList, h.buildPushFile(identifier, hostIP, cloudIDInt64Val))
			hostInfos = append(hostInfos, &HostInfo{
				HostInnerIP: hostIP,
				CloudID:     cloudIDInt64Val,
				Times:       retryTimes,
			})
		}
	}

	if len(fileList) == 0 {
		return
	}

	// 推送主机身份信息
	pushFileResp := new(push_file_forsyncdata.API_CommRsp)
	failCount = 0
	for failCount < retryTimes {
		pushFileResp, err = h.gseTaskServerClient.PushFileV2(h.ctx, fileList)
		if err != nil || pushFileResp.MErrcode != common.CCSuccess {
			blog.Errorf("push host identifier to gse error: err: %v, errCode: %d, errMessage: %s",
				err, pushFileResp.MErrcode, pushFileResp.MErrmsg)
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}
	if failCount >= retryTimes {
		return
	}
	blog.Infof("push host identifier success: file: %v", fileList)

	// 将任务放到redis维护的任务list中
	if err := h.createNewTaskToTaskList(pushFileResp.MContent, hostInfos); err != nil {
		blog.Errorf("add task to redis task list error, taskID: %s, hostInfos: %v, err: %v",
			pushFileResp.MContent, hostInfos, err)
	}
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
	conf := getHostIdentifierFileConf(osType)
	fileInfo.MFile.MName = conf.FileName
	fileInfo.MFile.MPath = conf.FilePath
	fileInfo.MFile.MOwner = conf.FileOwner
	fileInfo.MFile.MRight = conf.FileRight
	return fileInfo
}

func (h *HostIdentifier) findHostIdentifier(cloudID int64, innerIP string, header http.Header) (string, error) {
	input := &metadata.SearchIdentifierParam{
		IP: metadata.IPParam{
			Data:    strings.Split(innerIP, ","),
			CloudID: &cloudID,
		},
	}
	resp, err := h.engine.CoreAPI.TopoServer().Instance().SearchIdentifier(h.ctx, common.BKInnerObjIDHost,
		input, header)
	if err != nil || len(resp.Data.Info) == 0 {
		return "", fmt.Errorf("get host identifier error")
	}
	hostIdentifier, err := json.Marshal(resp.Data.Info[0])
	if err != nil {
		return "", err
	}
	return string(hostIdentifier), nil
}
