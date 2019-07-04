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
	"configcenter/src/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"
	hutil "configcenter/src/scene_server/host_server/util"

	"github.com/emicklei/go-restful"
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

// delete hosts from resource pool
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

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Delete, iHostIDArr...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v", iHostIDArr, err)
		if err != auth.NoAuthorizeError {
			resp.WriteEntity(&meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail)})
			return
		}
		resp.WriteEntity(s.AuthManager.GenDeleteHostBatchNoPermissionResp(iHostIDArr))
		return
	}

	appID, err := srvData.lgc.GetDefaultAppID(srvData.ctx)
	if err != nil {
		blog.Errorf("delete host batch, but got invalid app id, err: %v,input:%s,rid:%s", err, opt, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	hostFields, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, nil)
	if err != nil {
		blog.Errorf("delete host batch failed, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	// auth: unregister hosts
	if err := s.AuthManager.DeregisterHostsByID(srvData.ctx, srvData.header, iHostIDArr...); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v", iHostIDArr, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)})
		return
	}

	logConentsMap := make(map[int64]meta.SaveAuditLogParams, 0)
	for _, hostID := range iHostIDArr {
		logger := srvData.lgc.NewHostLog(srvData.ctx, srvData.ownerID)
		if err := logger.WithPrevious(srvData.ctx, strconv.FormatInt(hostID, 10), hostFields); err != nil {
			blog.Errorf("delete host batch, but get pre host data failed, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}

		logConentsMap[hostID] = logger.AuditLog(srvData.ctx, hostID)
	}

	input := &meta.DeleteHostRequest{
		ApplicationID: appID,
		HostIDArr:     iHostIDArr,
	}
	delResult, err := s.CoreAPI.CoreService().Host().DeleteHost(srvData.ctx, srvData.header, input)
	if err != nil {
		blog.Error("DeleteHostBatch DeleteHost http do error. err:%s, input:%s, rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}

	// ensure delete host add log
	for _, ex := range delResult.Data {
		delete(logConentsMap, ex.OriginIndex)
	}
	var logConents []meta.SaveAuditLogParams
	for _, item := range logConentsMap {
		item.Model = common.BKInnerObjIDHost
		item.OpDesc = "delete host"
		item.OpType = auditoplog.AuditOpTypeDel
		logConents = append(logConents, item)
	}
	if len(logConents) > 0 {
		auditResult, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, logConents...)
		if err != nil || (err == nil && !auditResult.Result) {
			blog.Errorf("delete host in batch, but add host audit log failed, err: %v, result err: %v,rid:%s", err, auditResult.ErrMsg, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrAuditSaveLogFaile)})
			return
		}
	}

	if !delResult.Result {
		blog.Errorf("DeleteHostBatch DeleteHost http reply error. result: %#v, input:%#v, rid:%s", delResult, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) GetHostInstanceProperties(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	hostID := req.PathParameter("bk_host_id")
	hostIDInt64, err := util.GetInt64ByInterface(hostID)
	if err != nil {
		blog.Errorf("convert hostID to int64, err: %v,host:%s,rid:%s", err, hostID, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, common.BKHostIDField)})
		return
	}

	details, _, err := srvData.lgc.GetHostInstanceDetails(srvData.ctx, srvData.ownerID, hostID)
	if err != nil {
		blog.Errorf("get host defails failed, err: %v,host:%s,rid:%s", err, hostID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if len(details) == 0 {
		blog.Errorf("host not found, hostID: %v,rid:%s", hostID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostNotFound)})
		return
	}

	// check authorization
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, hostIDInt64); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v", hostIDInt64, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
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

// HostSnapInfo return host state
func (s *Service) HostSnapInfo(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	hostID := req.PathParameter(common.BKHostIDField)
	hostIDInt64, err := strconv.ParseInt(hostID, 10, 64)
	if err != nil {
		blog.Errorf("HostSnapInfohttp hostID convert to int64 failed, err:%v, input:%+v, rid:%s", err, hostID, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedInt)})
		return
	}

	// check authorization
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, hostIDInt64); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v", hostIDInt64, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	// get snapshot
	result, err := s.CoreAPI.HostController().Host().GetHostSnap(srvData.ctx, hostID, srvData.header)

	if err != nil {
		blog.Errorf("HostSnapInfohttp do error, err: %v ,input:%#v, rid:%s", err, hostID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPReadBodyFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("HostSnapInfohttp reponse erro, err code:%d,err msg:%s, input:%#v, rid:%s", result.Code, result.ErrMsg, hostID, srvData.rid)
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

// add host to host resource pool
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
			blog.Errorf("add host, but get default app id failed, err: %v,input:%+v,rid:%s", err, hostList, srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
			return
		}
	}

	cond := hutil.NewOperation().WithModuleName(common.DefaultResModuleName).WithAppID(appID).MapStr()
	cond.Set(common.BKDefaultField, common.DefaultResModuleFlag)
	moduleID, err := srvData.lgc.GetResoulePoolModuleID(srvData.ctx, cond)
	if err != nil {
		blog.Errorf("add host, but get module id failed, err: %s,input: %+v,rid: %s", err.Error(), hostList, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	hostIDs, succ, updateErrRow, errRow, err := srvData.lgc.AddHost(srvData.ctx, appID, []int64{moduleID}, srvData.ownerID, hostList.HostInfo, hostList.InputType)
	retData := make(map[string]interface{})
	if err != nil {
		blog.Errorf("add host failed, succ: %v, update: %v, err: %v, %v,input:%+v,rid:%s", succ, updateErrRow, err, errRow, hostList, srvData.rid)
		retData["error"] = errRow
		retData["update_error"] = updateErrRow
		resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{Result: false, Code: common.CCErrHostCreateFail, ErrMsg: srvData.ccErr.Error(common.CCErrHostCreateFail).Error()},
			Data:     retData,
		})
		return
	}
	retData["success"] = succ

	// auth: register hosts
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDs...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v", hostIDs, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(retData))
}

// Deprecated:
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

	appID, err := srvData.lgc.GetDefaultAppID(srvData.ctx)
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

	// check authorization
	// is AddHostFromAgent's authentication the same with common api?
	// auth: check authorization
	// if err := s.AuthManager.AuthorizeCreateHost(srvData.ctx, srvData.header, appID); err != nil {
	// 	blog.Errorf("check add host authorization failed, business: %+v, err: %v", appID, err)
	// 	resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
	// 	return
	// }

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

	hostIDs, succ, updateErrRow, errRow, err := srvData.lgc.AddHost(srvData.ctx, appID, []int64{moduleID}, common.BKDefaultOwnerID, addHost, "")
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

	// register hosts
	// auth: register hosts
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDs...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v", hostIDs, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
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

	hostIDArray := host.ExtractHostIDs()
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, *hostIDArray...); err != nil {
		blog.Errorf("check host authorization failed, hostID: %+v, err: %+v", hostIDArray, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
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

	// auth: check authorization
	hostIDArray := host.ExtractHostIDs()
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, *hostIDArray...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v", hostIDArray, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
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

	// check authorization
	hostIDArr := make([]int64, 0)
	for _, id := range strings.Split(hostIDStr, ",") {
		hostID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			blog.Errorf("update host batch, but got invalid host id[%s], err: %v,rid:%s", id, err, srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
			return
		}
		hostIDArr = append(hostIDArr, hostID)
	}
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Update, hostIDArr...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v", hostIDArr, err)
		if err != auth.NoAuthorizeError {
			resp.WriteEntity(&meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail)})
			return
		}
		resp.WriteEntity(s.AuthManager.GenDeleteHostBatchNoPermissionResp(hostIDArr))
		return
	}

	logPreConents := make(map[int64]meta.SaveAuditLogParams, 0)
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
			// conds.Set(common.MetadataField, businessMedata)
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
			blog.Errorf("UpdateHostBatch UpdateObject http response error, err code:%d,err msg:%s,input:%+v,param:%+v,rid:%s", result.Code, data, opt, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
			return
		}

		logPreConents[hostID] = audit.AuditLog(srvData.ctx, hostID)
	}

	opt := &meta.UpdateOption{
		Condition: mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: hostIDs}},
		Data:      mapstr.NewFromMap(data),
	}
	result, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, opt)
	if err != nil {
		blog.Errorf("UpdateHostBatch UpdateObject http do error, err: %v,input:%+v,param:%+v,rid:%s", err, data, opt, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("UpdateHostBatch UpdateObject http response error, err code:%d, err msg:%s, input:%+v, param:%+v, rid:%s", result.Code, data, opt, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	hostModuleConfig, err := srvData.lgc.GetConfigByCond(srvData.ctx, meta.HostModuleRelationRequest{HostIDArr: hostIDs})
	if err != nil {
		blog.Errorf("update host batch failed, ids[%v], err: %v,input:%+v,rid:%s", hostIDs, err, data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	var appID int64
	if len(hostModuleConfig) > 0 {
		appID = hostModuleConfig[0].AppID
	}

	logLastConents := make([]meta.SaveAuditLogParams, 0)
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

		logLastConents = append(logLastConents,
			meta.SaveAuditLogParams{
				ID:      hostID,
				Model:   common.BKInnerObjIDHost,
				Content: logContent,
				OpDesc:  "update host",
				OpType:  auditoplog.AuditOpTypeModify,
				ExtKey:  preLogContent.ExtKey,
				BizID:   appID,
			},
		)

	}
	auditresp, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, logLastConents...)
	if err != nil || (err == nil && !auditresp.Result) {
		blog.Errorf("update host batch, but add host[%v] audit failed, err: %v, %v,rid:%s", hostIDs, err, auditresp.ErrMsg, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDetailFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

// NewHostSyncAppTopo add new hosts to the business
// synchronize hosts directly to a module in a business if this host does not exist.
// otherwise, this operation will only change host's attribute.
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
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "host_info")})
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
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrTopoGetAppFailed)})
		return
	}

	moduleCond := []meta.ConditionItem{
		meta.ConditionItem{
			Field:    common.BKModuleIDField,
			Operator: common.BKDBIN,
			Value:    hostList.ModuleID,
		},
	}
	if len(hostList.ModuleID) > 1 {
		moduleCond = append(moduleCond, meta.ConditionItem{
			Field:    common.BKDefaultField,
			Operator: common.BKDBEQ,
			Value:    0,
		})
	}
	moduleIDS, err := srvData.lgc.GetModuleIDByCond(srvData.ctx, moduleCond) //srvData.lgc..NewHostSyncValidModule(req, data.ApplicationID, data.ModuleID, m.CC.ObjCtrl())
	if nil != err {
		blog.Errorf("NewHostSyncAppTop GetModuleIDByCond error. err:%s,input:%+v,rid:%s", err.Error(), hostList, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if len(moduleIDS) != len(hostList.ModuleID) {
		blog.Errorf("not found part module: source:%v, db:%v", hostList.ModuleID, moduleIDS)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostMulueIDNotFoundORHasMutliInnerModuleIDFailed)})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeCreateHost(srvData.ctx, srvData.header, hostList.ApplicationID); err != nil {
		blog.Errorf("check add hosts authorization failed, business: %d, err: %v", hostList.ApplicationID, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	hostIDs, succ, updateErrRow, errRow, err := srvData.lgc.AddHost(srvData.ctx, hostList.ApplicationID, hostList.ModuleID, srvData.ownerID, hostList.HostInfo, common.InputTypeApiNewHostSync)
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

	// register host to iam
	// auth: check authorization
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDs...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v", hostIDs, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
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

	// get host in set
	var condition meta.HostModuleRelationRequest
	hostIDArr := make([]int64, 0)

	if 0 == data.SetID && 0 == data.ModuleID {
		blog.Errorf("MoveSetHost2IdleModule bk_set_id and bk_module_id cannot be empty at the same time,input:%#v,rid:%s", data, util.GetHTTPCCRequestID(pheader))
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	if 0 != data.SetID {
		condition.SetIDArr = []int64{data.SetID}
	}
	if 0 != data.ModuleID {
		condition.ModuleIDArr = []int64{data.ModuleID}
	}

	condition.ApplicationID = data.ApplicationID
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
		hostIDArr = append(hostIDArr, cell.HostID)
	}
	moduleCond := []meta.ConditionItem{
		meta.ConditionItem{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    data.ApplicationID,
		},
		meta.ConditionItem{
			Field:    common.BKDefaultField,
			Operator: common.BKDBEQ,
			Value:    common.DefaultResModuleFlag,
		},
	}

	moduleIDArr, err := srvData.lgc.GetModuleIDByCond(srvData.ctx, moduleCond)
	if err != nil {
		blog.Errorf("MoveSetHost2IdleModule GetModuleIDByCond error. err:%s, input:%#v, param:%#v, rid:%s", err.Error(), data, moduleCond, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
	}
	if len(moduleIDArr) == 0 {
		blog.Errorf("MoveSetHost2IdleModule GetModuleIDByCond error. err:%s, input:%#v, param:%#v, rid:%s", err.Error(), data, moduleCond, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleNotExist, "idle module")})
		return
	}
	idleModuleID := moduleIDArr[0]
	moduleHostConfigParams := make(map[string]interface{})
	moduleHostConfigParams[common.BKAppIDField] = data.ApplicationID
	audit := srvData.lgc.NewHostModuleLog(hostIDArr)

	// auth: check authorization
	// if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.MoveHostToBizIdleModule, hostIDArr...); err != nil {
	// 	blog.Errorf("check host authorization failed, hosts: %+v, err: %v", hostIDArr, err)
	// 	resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
	// 	return
	// }
	// // step2. check permission for target business
	// if err := s.AuthManager.AuthorizeCreateHost(srvData.ctx, srvData.header, data.ApplicationID); err != nil {
	// 	blog.Errorf("check add host authorization failed, business: %d, err: %v", data.ApplicationID, err)
	// 	resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
	// 	return
	// }
	// step3. deregister host from iam
	if err := s.AuthManager.DeregisterHostsByID(srvData.ctx, srvData.header, hostIDArr...); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v", hostIDArr, err)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)})
		return
	}
	hmInput := &meta.HostModuleRelationRequest{
		ApplicationID: data.ApplicationID,
		HostIDArr:     hostIDArr,
	}
	configResult, err := srvData.lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(srvData.ctx, srvData.header, hmInput)
	if nil != err {
		blog.Errorf("remove hosthostconfig http do error, error:%v, params:%v, input:%+v, rid:%s", err, hmInput, data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if !configResult.Result {
		blog.Errorf("remove hosthostconfig http reply error, result:%v, params:%v, input:%+v, rid:%s", configResult, hmInput, data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	hostIDMHMap := make(map[int64][]meta.ModuleHost, 0)
	for _, item := range configResult.Data {
		hostIDMHMap[item.HostID] = append(hostIDMHMap[item.HostID], item)
	}

	var exceptionArr []meta.ExceptionResult
	for _, hostID := range hostIDArr {
		hostMHArr, ok := hostIDMHMap[hostID]
		if !ok {
			// ignore  not exist the host under the current business,
			continue
		}
		toEmptyModule := true
		var newModuleIDArr []int64
		for _, item := range hostMHArr {
			if 0 != data.ModuleID && item.ModuleID == data.ModuleID {
				continue
			}
			if 0 != data.SetID && 0 == data.ModuleID && item.SetID == data.SetID {
				continue
			}

			toEmptyModule = false
			newModuleIDArr = append(newModuleIDArr, item.ModuleID)
		}

		var opResult *meta.OperaterException
		if toEmptyModule {
			input := &meta.TransferHostToInnerModule{
				ApplicationID: data.ApplicationID,
				ModuleID:      idleModuleID,
				HostID:        []int64{hostID},
			}
			opResult, err = srvData.lgc.CoreAPI.CoreService().Host().TransferHostToInnerModule(srvData.ctx, srvData.header, input)
		} else {
			input := &meta.HostsModuleRelation{
				ApplicationID: data.ApplicationID,
				HostID:        []int64{hostID},
				ModuleID:      newModuleIDArr,
			}
			opResult, err = srvData.lgc.CoreAPI.CoreService().Host().TransferHostModule(srvData.ctx, srvData.header, input)
		}
		if err != nil {
			blog.Errorf("MoveSetHost2IdleModule handle error. err:%s, to idle module:%v, input:%#v, hostID:%d, rid:%s", err.Error(), toEmptyModule, data, hostID, srvData.rid)
			ccErr := srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
			exceptionArr = append(exceptionArr, meta.ExceptionResult{Code: int64(ccErr.GetCode()), Message: ccErr.Error(), OriginIndex: hostID})
		}
		if !opResult.Result {
			if len(opResult.Data) > 0 {
				blog.Errorf("MoveSetHost2IdleModule handle reply error. result:%#v, to idle module:%v, input:%#v, hostID:%d, rid:%s", opResult, toEmptyModule, data, hostID, srvData.rid)
				exceptionArr = append(exceptionArr, opResult.Data...)
			} else {
				blog.Errorf("MoveSetHost2IdleModule handle reply error. result:%#v, to idle module:%v, input:%#v, hostID:%d, rid:%s", opResult, toEmptyModule, data, hostID, srvData.rid)
				exceptionArr = append(exceptionArr, meta.ExceptionResult{
					Code:        int64(opResult.Code),
					Message:     opResult.ErrMsg,
					OriginIndex: hostID,
				})
			}
		}

	}

	audit.SaveAudit(srvData.ctx, data.ApplicationID, srvData.user, "host to empty module")

	// register host to iam
	// auth: check authorization
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDArr...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid:%s", hostIDArr, err, srvData.rid)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
		return
	}
	if len(exceptionArr) > 0 {
		blog.Errorf("MoveSetHost2IdleModule has exception. exception:%#v, rid:%s", exceptionArr, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail), Data: exceptionArr})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
	return
}

func (s *Service) ip2hostID(srvData *srvComm, ip string, cloudID int64) (hostID int64, err error) {
	// FIXME there must be a better ip to hostID solution
	condition := common.KvMap{
		common.BKHostInnerIPField: ip,
		common.BKCloudIDField:     cloudID,
	}

	phpapi := srvData.lgc.NewPHPAPI()
	hostMap, hostIDArr, err := phpapi.GetHostMapByCond(srvData.ctx, condition)
	if err != nil {
		err := fmt.Errorf("GetHostMapByCond failed, %v", err)
		return 0, err
	}
	if len(hostIDArr) == 0 {
		return 0, nil
	}

	hostMapData, ok := hostMap[hostIDArr[0]]
	if false == ok {
		blog.Errorf("ip2hostID source ip invalid, raw data format hostMap:%+v, ip:%+v, cloudID:%+v, rid:%s", hostMap, ip, cloudID, srvData.rid)
		return 0, fmt.Errorf("ip %d:%s not found", cloudID, ip)
	}

	hostID, err = util.GetInt64ByInterface(hostMapData[common.BKHostIDField])
	if nil != err {
		blog.Errorf("ip2hostID bk_host_id field not found hostmap:%+v ip:%+v, cloudID:%+v,rid:%s", hostMapData, ip, cloudID, srvData.rid)
		return 0, fmt.Errorf("ip %+v:%+v not found", cloudID, ip)
	}

	return hostID, nil
}

// CloneHostProperty clone host property from src host to dst host
func (s *Service) CloneHostProperty(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := &meta.CloneHostPropertyParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("CloneHostProperty , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if 0 == input.AppID {
		blog.Errorf("CloneHostProperty, application not found input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID")})
		return
	}

	// authorization check
	srcHostID, err := s.ip2hostID(srvData, input.OrgIP, input.CloudID)
	if err != nil {
		blog.Errorf("ip2hostID failed, ip:%s, input:%+v, rid:%s", input.OrgIP, input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, "OrgIP")})
		return
	}
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, srcHostID); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid:%s", srcHostID, err, srvData.rid)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}
	// step2. verify has permission to update dst host
	dstHostID, err := s.ip2hostID(srvData, input.DstIP, input.CloudID)
	if err != nil {
		blog.Errorf("ip2hostID failed, ip:%s, input:%+v, rid:%s", input.DstIP, input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, "DstIP")})
		return
	}
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Update, dstHostID); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid:%s", dstHostID, err, srvData.rid)
		resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	res, err := srvData.lgc.CloneHostProperty(srvData.ctx, input, input.AppID, input.CloudID)
	if nil != err {
		blog.Errorf("CloneHostProperty ,application not int , err: %v, input:%#v, rid:%s", err, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res,
	})
}
