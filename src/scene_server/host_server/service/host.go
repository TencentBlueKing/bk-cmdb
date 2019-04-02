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

package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"
	hutil "configcenter/src/scene_server/host_server/util"
)

type AppResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    DataInfo    `json:"data"`
}

type DataInfo struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

func (s *Service) DeleteHostBatch(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	opt := new(meta.DeleteHostBatchOpt)
	if err := json.NewDecoder(req.Request.Body).Decode(opt); err != nil {
		blog.Errorf("delete host batch , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	hostIDArr := strings.Split(opt.HostID, ",")
	var iHostIDArr []int64
	for _, i := range hostIDArr {
		iHostID, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			blog.Errorf("delete host batch, but got invalid host id, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
			return
		}
		iHostIDArr = append(iHostIDArr, iHostID)
	}

	condition := make(map[string]interface{})
	condition = hutil.NewOperation().WithDefaultField(int64(common.DefaultAppFlag)).WithOwnerID(srvData.ownerID).MapStr()
	query := meta.QueryCondition{Condition: condition}
	query.Limit.Limit = 1
	query.Limit.Offset = 0
	result, err := s.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDApp, &query)
	if err != nil {
		blog.Errorf("delete host batch  SearchObjects http do error, err: %v,input:+v,rid:%s", err, opt, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("delete host in batch SearchObjects http reponse erro, err code:%d,err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, opt, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	if len(result.Data.Info) == 0 {
		blog.Errorf("delete host batch, but can not found it's instance. input:%+v,rid:%s", opt, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostNotFound)})
		return
	}

	appID, err := result.Data.Info[0].Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("delete host batch, but got invalid app id, err: %v, appinfo:%+v,input:%s,rid:%s", err, result.Data.Info[0], opt, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	hostFields, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, nil)
	if err != nil {
		blog.Errorf("delete host batch failed, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	var logConents []auditoplog.AuditLogExt
	for _, hostID := range iHostIDArr {
		logger := srvData.lgc.NewHostLog(srvData.ctx, srvData.ownerID)
		if err := logger.WithPrevious(srvData.ctx, strconv.FormatInt(hostID, 10), hostFields); err != nil {
			blog.Errorf("delete host batch, but get pre host data failed, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}

		delOptConfig := meta.ModuleHostConfigParams{
			HostID:        hostID,
			ApplicationID: appID,
		}
		result, err := s.CoreAPI.HostController().Module().DelModuleHostConfig(srvData.ctx, srvData.header, &delOptConfig)
		if err != nil {
			blog.Errorf("delete host batch  DelModuleHostConfig http do error, err: %v,input:+v,params:%s,rid:%s", err, opt, delOptConfig, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
			return
		}
		if !result.Result {
			blog.Errorf("delete host batch  DelModuleHostConfig http response error, err code:%s,err msg:%s,input:+v,params:%s,rid:%s", result.Code, result.ErrMsg, opt, delOptConfig, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
			return
		}

		logConents = append(logConents, *logger.AuditLog(srvData.ctx, hostID))
	}

	hostCond := mapstr.MapStr{
		common.BKDBIN: iHostIDArr,
	}
	condInput := &meta.DeleteOption{
		Condition: mapstr.MapStr{
			common.BKHostIDField: hostCond,
		},
	}
	delResult, err := s.CoreAPI.CoreService().Instance().DeleteInstanceCascade(srvData.ctx, srvData.header, common.BKInnerObjIDHost, condInput)
	if err != nil || (err == nil && !delResult.Result) {
		blog.Errorf("delete host in batch, but delete host failed, err: %v, result err: %v,rid:%s", err, delResult.ErrMsg, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail)})
		return
	}

	addHostLogs := common.KvMap{common.BKContentField: logConents, common.BKOpDescField: "delete host", common.BKOpTypeField: auditoplog.AuditOpTypeDel}
	auditResult, err := s.CoreAPI.AuditController().AddHostLogs(srvData.ctx, srvData.ownerID, strconv.FormatInt(appID, 10), srvData.user, srvData.header, addHostLogs)
	if err != nil || (err == nil && !auditResult.Result) {
		blog.Errorf("delete host in batch, but add host audit log failed, err: %v, result err: %v,rid:%s", err, auditResult.ErrMsg, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrAuditSaveLogFaile)})
		return
	}
	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) GetHostInstanceProperties(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	hostID := req.PathParameter("bk_host_id")

	details, _, err := srvData.lgc.GetHostInstanceDetails(srvData.ctx, srvData.ownerID, hostID)
	if err != nil {
		blog.Errorf("get host defails failed, err: %v,host:%s,rid:%s", err, hostID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	attribute, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, nil)
	if err != nil {
		blog.Errorf("get host attribute fields failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	result := make([]meta.HostInstanceProperties, 0)
	for _, attr := range attribute {
		if attr.PropertyID == common.BKChildStr {
			continue
		}
		result = append(result, meta.HostInstanceProperties{
			PropertyID:    attr.PropertyID,
			PropertyName:  attr.PropertyName,
			PropertyValue: details[attr.PropertyID],
		})
	}

	resp.WriteEntity(meta.HostInstancePropertiesResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}

func (s *Service) HostSnapInfo(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	hostID := req.PathParameter(common.BKHostIDField)
	result, err := s.CoreAPI.HostController().Host().GetHostSnap(srvData.ctx, hostID, srvData.header)

	if err != nil {
		blog.Errorf("HostSnapInfohttp do error, err: %v,input:+v,rid:%s", err, hostID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPReadBodyFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("HostSnapInfohttp reponse erro, err code:%d,err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, hostID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	snap, err := logics.ParseHostSnap(result.Data.Data)
	if err != nil {
		blog.Errorf("get host snap info, but parse snap info failed, err: %v, hostID:%v,rid:%s", err, hostID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.HostSnapResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     snap,
	})
}

func (s *Service) AddHost(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	hostList := new(meta.HostList)
	if err := json.NewDecoder(req.Request.Body).Decode(hostList); err != nil {
		blog.Errorf("add host failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if hostList.HostInfo == nil {
		blog.Errorf("add host, but host info is nil.input:%+v,rid:%s", hostList, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	appID := hostList.ApplicationID
	if appID == 0 {
		// get default app id
		var err error
		appID, err = srvData.lgc.GetDefaultAppIDWithSupplier(srvData.ctx)
		if err != nil {
			blog.Errorf("add host, but get default appid failed, err: %v,input:%+v,rid:%s", err, hostList, srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
			return
		}
	}

	cond := hutil.NewOperation().WithModuleName(common.DefaultResModuleName).WithAppID(appID).MapStr()
	cond.Set(common.BKDefaultField, common.DefaultResModuleFlag)
	moduleID, err := srvData.lgc.GetResoulePoolModuleID(srvData.ctx, cond)
	if err != nil {
		blog.Errorf("add host, but get module id failed, err: %s,input:%+v,rid:%s", err.Error(), hostList, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	succ, updateErrRow, errRow, err := srvData.lgc.AddHost(srvData.ctx, appID, []int64{moduleID}, srvData.ownerID, hostList.HostInfo, hostList.InputType)
	retData := make(map[string]interface{})
	if err != nil {
		blog.Errorf("add host failed, succ: %v, update: %v, err: %v, %v,input:%+v,rid:%s", succ, updateErrRow, err, errRow, hostList, srvData.rid)
		retData["error"] = errRow
		retData["update_error"] = updateErrRow
		resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{false, common.CCErrHostCreateFail, srvData.ccErr.Error(common.CCErrHostCreateFail).Error()},
			Data:     retData,
		})
		return
	}
	retData["success"] = succ
	resp.WriteEntity(meta.NewSuccessResp(retData))
}

func (s *Service) AddHostFromAgent(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	agents := new(meta.AddHostFromAgentHostList)
	if err := json.NewDecoder(req.Request.Body).Decode(&agents); err != nil {
		blog.Errorf("add host from agent failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(agents.HostInfo) == 0 {
		blog.Errorf("add host from agent, but got 0 agents from body.input:%+v,rid:%s", agents, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "HostInfo")})
		return
	}

	appID, err := srvData.lgc.GetDefaultAppID(srvData.ctx, srvData.ownerID)
	if err != nil {
		blog.Errorf("AddHostFromAgent GetDefaultAppID error.input:%#v,rid:%s", agents, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}
	if 0 == appID {
		blog.Errorf("add host from agent, but got invalid default appid, err: %v,ownerID:%s,input:%#v,rid:%s", err, srvData.ownerID, agents, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrAddHostToModule, "bussiness not found")})
		return
	}

	opt := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(appID)
	moduleID, err := srvData.lgc.GetResoulePoolModuleID(srvData.ctx, opt.MapStr())
	if err != nil {
		blog.Errorf("add host from agent , but get module id failed, err: %v,ownerID:%s,input:%+v,rid:%s", err, srvData.ownerID, agents, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	agents.HostInfo["import_from"] = common.HostAddMethodAgent
	addHost := make(map[int64]map[string]interface{})
	addHost[1] = agents.HostInfo

	succ, updateErrRow, errRow, err := srvData.lgc.AddHost(srvData.ctx, appID, []int64{moduleID}, common.BKDefaultOwnerID, addHost, "")
	if err != nil {
		blog.Errorf("add host failed, succ: %v, update: %v, err: %v, %v,input:%+v,rid:%s", succ, updateErrRow, err, errRow, agents, srvData.rid)

		retData := make(map[string]interface{})
		retData["success"] = succ
		retData["error"] = errRow
		retData["update_error"] = updateErrRow
		resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{Result: false, Code: common.CCErrHostCreateFail, ErrMsg: srvData.ccErr.Error(common.CCErrHostCreateFail).Error()},
			Data:     retData,
		})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(succ))
}

func (s *Service) SearchHost(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	body := new(meta.HostCommonSearch)
	if err := json.NewDecoder(req.Request.Body).Decode(body); err != nil {
		blog.Errorf("search host failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	host, err := srvData.lgc.SearchHost(srvData.ctx, body, false)
	if err != nil {
		blog.Errorf("search host failed, err: %v,input:%+v,rid:%s", err, body, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostGetFail)})
		return
	}

	resp.WriteEntity(meta.SearchHostResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *host,
	})
}

func (s *Service) SearchHostWithAsstDetail(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	body := new(meta.HostCommonSearch)
	if err := json.NewDecoder(req.Request.Body).Decode(body); err != nil {
		blog.Errorf("search host failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	host, err := srvData.lgc.SearchHost(srvData.ctx, body, true)
	if err != nil {
		blog.Errorf("search host failed, err: %v,input:%+v,rid:%s", err, body, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.SearchHostResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *host,
	})
}

func (s *Service) UpdateHostBatch(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	data := mapstr.New()
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update host batch failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	hostIDStr, err := data.String(common.BKHostIDField)
	if err != nil {
		blog.Errorf("update host batch failed, but without host id field, err:%s.input:%+v,rid:%s", err.Error(), data, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedString, common.BKHostIDField)})
		return

	}

	businessMedata := data.Remove(common.MetadataField)
	data.Remove(common.BKHostIDField)
	hostFields, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, nil)
	if err != nil {
		blog.Errorf("update host batch, but get host attribute for audit failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDetailFail)})
		return
	}


	logPreConents := make(map[int64]auditoplog.AuditLogExt, 0)
	hostIDs := make([]int64, 0)
	for _, id := range strings.Split(hostIDStr, ",") {
		hostID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			blog.Errorf("update host batch, but got invalid host id[%s], err: %v,rid:%s", id, err, srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
			return
		}
		hostIDs = append(hostIDs, hostID)
		conds := mapstr.New()
		if businessMedata != nil {
			//conds.Set(common.MetadataField, businessMedata)
			// TODO use metadata
		}
		conds.Set(common.BKHostIDField, hostID)
		opt := &meta.UpdateOption{
			Condition: conds,
			Data:      mapstr.NewFromMap(data),
		}
        audit := srvData.lgc.NewHostLog(srvData.ctx, srvData.ownerID)
        if err := audit.WithPrevious(srvData.ctx, id, hostFields); err != nil {
            blog.Errorf("update host batch, but get host[%s] pre data for audit failed, err: %v", id, err)
            resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDetailFail)})
            return
        }
        result, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, opt)
		if err != nil {
			blog.Errorf("UpdateHostBatch UpdateObject http do error, err: %v,input:%+v,param:%+v,rid:%s", err, data, opt, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
			return
		}
		if !result.Result {
            blog.Errorf("UpdateHostBatch UpdateObject http response error, err code:%s,err msg:%s,input:%+v,param:%+v,rid:%s", result.Code, data, opt, srvData.rid)
            resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
            return 
        }

		logPreConents[hostID] = *audit.AuditLog(srvData.ctx, hostID)
	}

	hostModuleConfig, err := srvData.lgc.GetConfigByCond(srvData.ctx, map[string][]int64{common.BKHostIDField: hostIDs})
	if err != nil {
		blog.Errorf("update host batch failed, ids[%v], err: %v,input:%+v,rid:%s", hostIDs, err, data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	appID := "0"
	if len(hostModuleConfig) != 0 {
		appID = strconv.FormatInt(hostModuleConfig[0][common.BKAppIDField], 10)
	}

	logLastConents := make([]auditoplog.AuditLogExt, 0)
	for _, hostID := range hostIDs {

		audit := srvData.lgc.NewHostLog(srvData.ctx, common.BKDefaultOwnerID)
		if err := audit.WithPrevious(srvData.ctx, strconv.FormatInt(hostID, 10), hostFields); err != nil {
			blog.Errorf("update host batch, but get host[%v] pre data for audit failed, err: %v", hostID, err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDetailFail)})
			return
		}
		logContent := audit.Content
		logContent.CurData = logContent.PreData
		preLogContent, ok := logPreConents[hostID]
		if ok {
			content, ok := preLogContent.Content.(*meta.Content)
			if ok {
				logContent.PreData = content.PreData
			}
		}

		logLastConents = append(logLastConents, auditoplog.AuditLogExt{ID: hostID, Content: logContent, ExtKey: preLogContent.ExtKey})
	}
	log := common.KvMap{common.BKContentField: logLastConents, common.BKOpDescField: "update host", common.BKOpTypeField: auditoplog.AuditOpTypeModify}
	aResult, err := s.CoreAPI.AuditController().AddHostLogs(srvData.ctx, srvData.ownerID, appID, srvData.user, srvData.header, log)
	if err != nil || (err == nil && !aResult.Result) {
		blog.Errorf("update host batch, but add host[%v] audit failed, err: %v, %v,rid:%s", hostIDs, err, aResult.ErrMsg, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDetailFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) NewHostSyncAppTopo(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	hostList := new(meta.HostSyncList)
	if err := json.NewDecoder(req.Request.Body).Decode(hostList); err != nil {
		blog.Errorf("add host failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if hostList.HostInfo == nil {
		blog.Errorf("add host, but host info is nil.input:%+v,rid:%s", hostList, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}
	if 0 == len(hostList.ModuleID) {
		blog.Errorf("host sync app  parameters required moduleID,input:%+v,rid:%s", hostList, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, common.BKModuleIDField)})
		return
	}

	if common.BatchHostAddMaxRow < len(hostList.HostInfo) {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommXXExceedLimit, "host_info ", common.BatchHostAddMaxRow)})
		return
	}

	appConds := map[string]interface{}{
		common.BKAppIDField: hostList.ApplicationID,
	}
	appInfo, err := srvData.lgc.GetAppDetails(srvData.ctx, "", appConds)
	if nil != err {
		blog.Errorf("host sync app %d error:%s,input:%+v,rid:%s", hostList.ApplicationID, err.Error(), hostList, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if 0 == len(appInfo) {
		blog.Errorf("host sync app %d not found, reply:%+v,input:%+v,rid:%s", hostList.ApplicationID, appInfo, hostList, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrTopoGetAppFaild, "not foud ")})
		return
	}

	moduleCond := []meta.ConditionItem{
		meta.ConditionItem{
			Field:    common.BKModuleIDField,
			Operator: common.BKDBIN,
			Value:    hostList.ModuleID,
		},
	}
	moduleIDS, err := srvData.lgc.GetModuleIDByCond(srvData.ctx, moduleCond) //srvData.lgc..NewHostSyncValidModule(req, data.ApplicationID, data.ModuleID, m.CC.ObjCtrl())
	if nil != err {
		blog.Errorf("NewHostSyncAppTop GetModuleIDByCond error. err:%s,input:%+v,rid:%s", err.Error(), hostList, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if len(moduleIDS) != len(hostList.ModuleID) {
		blog.Errorf("not found part module: source:%v, db:%v", hostList.ModuleID, moduleIDS)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrTopoGetModuleFailed, " not found part moudle id")})
		return

	}
	succ, updateErrRow, errRow, err := srvData.lgc.AddHost(srvData.ctx, hostList.ApplicationID, hostList.ModuleID, srvData.ownerID, hostList.HostInfo, common.InputTypeApiNewHostSync)
	if err != nil {
		blog.Errorf("add host failed, succ: %v, update: %v, err: %v, %v", succ, updateErrRow, err, errRow)

		retData := make(map[string]interface{})
		retData["success"] = succ
		retData["error"] = errRow
		retData["update_error"] = updateErrRow
		resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{Result: false, Code: common.CCErrHostCreateFail, ErrMsg: srvData.ccErr.Error(common.CCErrHostCreateFail).Error()},
			Data:     retData,
		})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(succ))
}

// MoveSetHost2IdleModule bk_set_id and bk_module_id cannot be empty at the same time
// Remove the host from the module or set.
// The host belongs to the current module or host only, and puts the host into the idle machine of the current service.
// When the host data is in multiple modules or sets. Disconnect the host from the module or set only
func (s *Service) MoveSetHost2IdleModule(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	srvData := s.newSrvComm(req.Request.Header)

	var data meta.SetHostConfigParams
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("MoveSetHost2IdleModule failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if 0 == data.ApplicationID {
		blog.Errorf("MoveSetHost2IdleModule bk_biz_id cannot be empty at the same time,input:%#v,rid:%s", data, util.GetHTTPCCRequestID(pheader))
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	//get host in set
	condition := make(map[string][]int64)
	hostIDArr := make([]int64, 0)
	sModuleIDArr := make([]int64, 0)

	if 0 == data.SetID && 0 == data.ModuleID {
		blog.Errorf("MoveSetHost2IdleModule bk_set_id and bk_module_id cannot be empty at the same time,input:%#v,rid:%s", data, util.GetHTTPCCRequestID(pheader))
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	if 0 != data.SetID {
		condition[common.BKSetIDField] = []int64{data.SetID}
	}
	if 0 != data.ModuleID {
		condition[common.BKModuleIDField] = []int64{data.ModuleID}
	}

	condition[common.BKAppIDField] = []int64{data.ApplicationID}
	hostResult, err := srvData.lgc.GetConfigByCond(srvData.ctx, condition) //logics.GetConfigByCond(req, m.CC.HostCtrl(), condition)
	if nil != err {
		blog.Errorf("read host from application  error:%v,input:%+v,rid:%s", err, data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	if 0 == len(hostResult) {
		blog.Warnf("no host in set,rid:%s", srvData.rid)
		resp.WriteEntity(meta.NewSuccessResp(nil))
		return
	}
	for _, cell := range hostResult {
		hostIDArr = append(hostIDArr, cell[common.BKHostIDField])
		sModuleIDArr = append(sModuleIDArr, cell[common.BKModuleIDField])
	}

	getModuleCond := make([]meta.ConditionItem, 0)
	getModuleCond = append(getModuleCond, meta.ConditionItem{Field: common.BKDefaultField, Operator: common.BKDBEQ, Value: common.DefaultResModuleFlag})
	getModuleCond = append(getModuleCond, meta.ConditionItem{Field: common.BKModuleNameField, Operator: common.BKDBEQ, Value: common.DefaultResModuleName})
	getModuleCond = append(getModuleCond, meta.ConditionItem{Field: common.BKAppIDField, Operator: common.BKDBEQ, Value: data.ApplicationID})

	moduleIDArr, err := srvData.lgc.GetModuleIDByCond(srvData.ctx, getModuleCond) //GetSingleModuleID(req, conds, m.CC.ObjCtrl())
	if nil != err {
		blog.Errorf("module params   error:%v,input:%+v,rid:%s", err, data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if 0 == len(moduleIDArr) {
		blog.Errorf("module params   error:%v,input:%+v,rid:%s", err, data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	moduleID := moduleIDArr[0]

	// idle modle not change
	if moduleID == data.ModuleID {
		resp.WriteEntity(meta.NewSuccessResp(nil))
		return
	}

	moduleHostConfigParams := make(map[string]interface{})
	moduleHostConfigParams[common.BKAppIDField] = data.ApplicationID
	audit := srvData.lgc.NewHostModuleLog(hostIDArr)

	for _, hostID := range hostIDArr {

		bl, err := srvData.lgc.IsHostExistInApp(srvData.ctx, data.ApplicationID, hostID)
		if nil != err {
			blog.Errorf("check host is exist in app error, params:{appid:%d, hostid:%s}, error:%s,rid:%s", data.ApplicationID, hostID, err.Error(), srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		if false == bl {
			blog.Errorf("host do not belong to the current application; error, params:{appid:%d, hostid:%s},rid:%s", data.ApplicationID, hostID, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostNotINAPP)})
			return
		}

		var toEmptyModule = true

		sCond := make(map[string][]int64)
		sCond[common.BKAppIDField] = []int64{data.ApplicationID}
		sCond[common.BKHostIDField] = []int64{hostID}
		configResult, err := srvData.lgc.GetConfigByCond(srvData.ctx, sCond)
		if nil != err {
			blog.Errorf("remove hosthostconfig error, params:%v, error:%v,input:%+v,rid:%s", sCond, err, data, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		for _, config := range configResult {
			if 0 != data.SetID && config[common.BKSetIDField] != data.SetID {
				toEmptyModule = false
			}
			if 0 != data.ModuleID && config[common.BKModuleIDField] != data.ModuleID {
				toEmptyModule = false
			}
		}

		moduleHostConfigParams := meta.ModuleHostConfigParams{HostID: hostID, ApplicationID: data.ApplicationID, ModuleID: sModuleIDArr}

		result, err := s.CoreAPI.HostController().Module().DelModuleHostConfig(srvData.ctx, srvData.header, &moduleHostConfigParams)
		if err != nil {
			blog.Errorf("remove hosthostconfig http do error, err: %v,input:%+v,param:%+v,rid:%s", err, data, moduleHostConfigParams, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
			return
		}
		if !result.Result {
			blog.Errorf("remove hosthostconfig  http response error, err code:%s,err msg:%s,input:%+v,param:%+v,rid:%s", result.Code, data, moduleHostConfigParams, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
			return
		}

		if toEmptyModule {
			moduleHostConfigParams = meta.ModuleHostConfigParams{HostID: hostID, ModuleID: []int64{moduleID}, ApplicationID: data.ApplicationID}
			result, err = s.CoreAPI.HostController().Module().AddModuleHostConfig(srvData.ctx, srvData.header, &moduleHostConfigParams)
			if err != nil {
				blog.Errorf("add hosthostconfig http do error, err: %v,input:%+v,param:%+v,rid:%s", err, data, moduleHostConfigParams, srvData.rid)
				resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
				return
			}
			if !result.Result {
				blog.Errorf("add hosthostconfig  http response error, err code:%s,err msg:%s,input:%+v,param:%+v,rid:%s", result.Code, data, moduleHostConfigParams, srvData.rid)
				resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
				return
			}
		}

	}

	audit.SaveAudit(srvData.ctx, strconv.FormatInt(data.ApplicationID, 10), srvData.user, "host to empty module")

	resp.WriteEntity(meta.NewSuccessResp(nil))
	return
}

func (s *Service) CloneHostProperty(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := &meta.CloneHostPropertyParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("CloneHostProperty , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if 0 == input.AppID {
		blog.Errorf("CloneHostProperty ,appliation not foud input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID")})
		return
	}

	res, err := srvData.lgc.CloneHostProperty(srvData.ctx, input, input.AppID, input.CloudID)
	if nil != err {
		blog.Errorf("CloneHostProperty ,appliation not int , err: %v, input:%v", err, input)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res,
	})
}
