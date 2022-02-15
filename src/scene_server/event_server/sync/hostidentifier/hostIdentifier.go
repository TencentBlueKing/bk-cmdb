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
	"net/http"
	"strconv"
	"time"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/gse/client"
	getstatus "configcenter/src/thirdparty/gse/get_agent_state_forsyncdata"
	pushfile "configcenter/src/thirdparty/gse/push_file_forsyncdata"

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
func NewHostIdentifier(ctx context.Context, redisCli redis.Client, engine *backbone.Engine, conf *HostIdentifierConf,
	taskClient *client.GseTaskServerClient, apiClient *client.GseApiServerClient) (*HostIdentifier, error) {
	h := &HostIdentifier{
		redisCli:            redisCli,
		ctx:                 ctx,
		engine:              engine,
		gseTaskServerClient: taskClient,
		gseApiServerClient:  apiClient,
		winFileConfig:       conf.WinFileConf,
		linuxFileConfig:     conf.LinuxFileConf,
		watchLimiter:        flowctrl.NewRateLimiter(conf.RateLimiter.Qps, conf.RateLimiter.Burst),
		fullLimiter:         flowctrl.NewRateLimiter(conf.RateLimiter.Qps, conf.RateLimiter.Burst),
	}
	return h, nil
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

		h.watchToSyncHostIdentifier(events, rid)

		eventOp.setCursor(lastCursor, rid)
	}
}

func (h *HostIdentifier) watchToSyncHostIdentifier(events []*IdentifierEvent, rid string) {

	// 1、查询主机状态
	statusReq := new(getstatus.AgentStatusRequest)
	for _, event := range events {
		info := buildForStatus(strconv.FormatInt(event.CloudID, 10), event.InnerIP)
		statusReq.Hosts = append(statusReq.Hosts, info...)
	}

	resp, err := h.getAgentStatus(statusReq, true, rid)
	if err != nil {
		blog.Errorf("get agent status error, host: %v, err: %v, rid: %s", events, err, rid)
		return
	}

	// 2、将处于on状态的主机拿出来构造推送信息
	fList := make([]*pushfile.API_FileInfoV2, 0)
	hostInfos := make([]*HostInfo, 0)
	for _, event := range events {
		isOn, hostIP := getStatusOnAgentIP(strconv.FormatInt(event.CloudID, 10), event.InnerIP, resp.Result_)
		if !isOn {
			blog.Infof("host %v agent status is off, rid: %s", event, rid)
			continue
		}

		blog.Infof("host %v agent status is on, ip: %s, rid: %s", event, hostIP, rid)

		file := h.buildPushFile(event.RawEvent, hostIP, event.CloudID)
		fList = append(fList, file)
		hostInfos = append(hostInfos, &HostInfo{
			HostID:      event.HostID,
			HostInnerIP: hostIP,
			CloudID:     event.CloudID,
		})
	}

	if len(fList) == 0 {
		return
	}

	// 3、推送主机身份信息
	h.watchLimiter.AcceptMany(int64(len(fList)))
	if _, err := h.pushFile(true, hostInfos, fList, rid); err != nil {
		blog.Errorf("push host identifier to gse error, err: %v, rid: %s", err, rid)
	}
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
			Fields: []string{common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField},
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
	statusReq := &getstatus.AgentStatusRequest{}
	for _, hostInfo := range hosts {
		info := buildForStatus(util.GetStrByInterface(hostInfo[common.BKCloudIDField]),
			util.GetStrByInterface(hostInfo[common.BKHostInnerIPField]))

		statusReq.Hosts = append(statusReq.Hosts, info...)
	}

	resp, err := h.getAgentStatus(statusReq, false, rid)
	if err != nil {
		blog.Errorf("get agent status error,  hostInfo: %v, err: %v, rid: %s", hosts, err, rid)
		return nil, err
	}

	// 2、将处于on状态的主机拿出来构造主机身份推送信息
	hostIDs := make([]int64, 0)
	hostInfos := make([]*HostInfo, 0)
	// 此map保存hostID和该host处于on的agent的ip的对应关系
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
		isOn, hostIP := getStatusOnAgentIP(strconv.FormatInt(cloudID, 10), innerIP, resp.Result_)
		if !isOn {
			blog.Infof("host %v agent status is off, rid: %s", hostInfo, rid)
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
		return nil, errors.New("the host agent status is off")
	}

	// 3、查询主机身份并推送
	return h.getHostIdentifierAndPush(hostIDs, hostMap, hostInfos, rid, header)
}

func (h *HostIdentifier) getAgentStatus(status *getstatus.AgentStatusRequest,
	always bool, rid string) (*getstatus.AgentStatusResponse, error) {

	var err error
	failCount := 0
	resp := new(getstatus.AgentStatusResponse)

	// 调用gse api server 查询agent状态
	for always || failCount < retryTimes {
		resp, err = h.gseApiServerClient.GetAgentStatus(context.Background(), status)
		if err != nil {
			blog.Errorf("get host agent status error, err: %v, rid: %s", err, rid)
			failCount++
			sleepForFail(failCount)
			continue
		}

		if resp.BkErrorCode != common.CCSuccess {
			blog.Errorf("get agent status fail, code: %d, msg: %s, rid: %s", resp.BkErrorCode, resp.BkErrorMsg, rid)
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}

	if !always && failCount >= retryTimes {
		return nil, errors.New("find agent status from apiServer error")
	}

	return resp, nil
}

func (h *HostIdentifier) buildPushFile(hostIdentifier, hostIP string, cloudID int64) *pushfile.API_FileInfoV2 {
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
	fileInfo.MFile.MRight = conf.FilePrivilege

	return fileInfo
}

func (h *HostIdentifier) getHostIdentifierAndPush(hostIDs []int64, hostMap map[int64]string,
	hostInfos []*HostInfo, rid string, header http.Header) (*Task, error) {
	if len(hostIDs) == 0 {
		blog.Errorf("hostIDs count is 0, rid: %s", rid)
		return nil, errors.New("hostIDs count is 0")
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
		return nil, errors.New("get host identifier success, but can not build file to push")
	}

	// 3、推送主机身份信息
	h.fullLimiter.AcceptMany(int64(len(fileList)))
	return h.pushFile(false, hostInfos, fileList, rid)
}
