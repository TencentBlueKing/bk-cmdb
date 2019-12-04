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

	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
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
func (s *Service) DeleteHostBatchFromResourcePool(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	opt := new(meta.DeleteHostBatchOpt)
	if err := json.NewDecoder(req.Request.Body).Decode(opt); err != nil {
		blog.Errorf("delete host batch , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	hostIDArr := strings.Split(opt.HostID, ",")
	var iHostIDArr []int64
	for _, i := range hostIDArr {
		iHostID, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			blog.Errorf("delete host batch, but got invalid host id, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
			_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
			return
		}
		iHostIDArr = append(iHostIDArr, iHostID)
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Delete, iHostIDArr...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", iHostIDArr, err, srvData.rid)
		if err != auth.NoAuthorizeError {
			_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail)})
			return
		}
		perm, err := s.AuthManager.GenEditHostBatchNoPermissionResp(srvData.ctx, srvData.header, authcenter.Delete, iHostIDArr)
		if err != nil {
			resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail)})
			return
		}
		resp.WriteEntity(perm)
		return
	}

	appID, err := srvData.lgc.GetDefaultAppID(srvData.ctx)
	if err != nil {
		blog.Errorf("delete host batch, but got invalid app id, err: %v,input:%s,rid:%s", err, opt, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	hostFields, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, nil)
	if err != nil {
		blog.Errorf("delete host batch failed, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	// auth: unregister hosts
	if err := s.AuthManager.DeregisterHostsByID(srvData.ctx, srvData.header, iHostIDArr...); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v, rid: %s", iHostIDArr, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)})
		return
	}

	logContentMap := make(map[int64]meta.SaveAuditLogParams, 0)
	for _, hostID := range iHostIDArr {
		logger := srvData.lgc.NewHostLog(srvData.ctx, srvData.ownerID)
		if err := logger.WithPrevious(srvData.ctx, strconv.FormatInt(hostID, 10), hostFields); err != nil {
			blog.Errorf("delete host batch, but get pre host data failed, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}

		logContentMap[hostID] = logger.AuditLog(srvData.ctx, hostID)
	}

	input := &meta.DeleteHostRequest{
		ApplicationID: appID,
		HostIDArr:     iHostIDArr,
	}
	delResult, err := s.CoreAPI.CoreService().Host().DeleteHostFromSystem(srvData.ctx, srvData.header, input)
	if err != nil {
		blog.Error("DeleteHostBatch DeleteHost http do error. err:%s, input:%s, rid:%s", err.Error(), input, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}

	// ensure delete host add log
	for _, ex := range delResult.Data {
		delete(logContentMap, ex.OriginIndex)
	}
	var logContents []meta.SaveAuditLogParams
	for _, item := range logContentMap {
		item.Model = common.BKInnerObjIDHost
		item.OpDesc = "delete host"
		item.OpType = auditoplog.AuditOpTypeDel
		logContents = append(logContents, item)
	}
	if len(logContents) > 0 {
		auditResult, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, logContents...)
		if err != nil || (err == nil && !auditResult.Result) {
			blog.Errorf("delete host in batch, but add host audit log failed, err: %v, result err: %v,rid:%s", err, auditResult.ErrMsg, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrAuditSaveLogFailed)})
			return
		}
	}

	if !delResult.Result {
		blog.Errorf("DeleteHostBatch DeleteHost http reply error. result: %#v, input:%#v, rid:%s", delResult, input, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail)})
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(nil))
}

// get host instance's properties as follows:
// host object property id: "bk_host_name"
// host object property name: "host"
// host object property value: "centos7"

func (s *Service) GetHostInstanceProperties(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	hostID := req.PathParameter("bk_host_id")
	hostIDInt64, err := util.GetInt64ByInterface(hostID)
	if err != nil {
		blog.Errorf("convert hostID to int64, err: %v,host:%s,rid:%s", err, hostID, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, common.BKHostIDField)})
		return
	}

	details, _, err := srvData.lgc.GetHostInstanceDetails(srvData.ctx, srvData.ownerID, hostID)
	if err != nil {
		blog.Errorf("get host details failed, err: %v,host:%s,rid:%s", err, hostID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if len(details) == 0 {
		blog.Errorf("host not found, hostID: %v,rid:%s", hostID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostNotFound)})
		return
	}

	// check authorization
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, hostIDInt64); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDInt64, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	attribute, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, nil)
	if err != nil {
		blog.Errorf("get host attribute fields failed, err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
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

	responseData := meta.HostInstancePropertiesResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	}
	_ = resp.WriteEntity(responseData)
}

// HostSnapInfo return host state
func (s *Service) HostSnapInfo(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	hostID := req.PathParameter(common.BKHostIDField)
	hostIDInt64, err := strconv.ParseInt(hostID, 10, 64)
	if err != nil {
		blog.Errorf("HostSnapInfo hostID convert to int64 failed, err:%v, input:%+v, rid:%s", err, hostID, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedInt)})
		return
	}

	// check authorization
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, hostIDInt64); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDInt64, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	// get snapshot
	result, err := s.CoreAPI.CoreService().Host().GetHostSnap(srvData.ctx, srvData.header, hostID)

	if err != nil {
		blog.Errorf("HostSnapInfo, http do error, err: %v ,input:%#v, rid:%s", err, hostID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPReadBodyFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("HostSnapInfo, http response error, err code:%d,err msg:%s, input:%#v, rid:%s", result.Code, result.ErrMsg, hostID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	snap, err := logics.ParseHostSnap(result.Data.Data)
	if err != nil {
		blog.Errorf("get host snap info, but parse snap info failed, err: %v, hostID:%v,rid:%s", err, hostID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	responseData := meta.HostSnapResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     snap,
	}
	_ = resp.WriteEntity(responseData)
}

// add host to host resource pool
func (s *Service) AddHost(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	hostList := new(meta.HostList)
	if err := json.NewDecoder(req.Request.Body).Decode(hostList); err != nil {
		blog.Errorf("add host failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if hostList.HostInfo == nil {
		blog.Errorf("add host, but host info is nil.input:%+v,rid:%s", hostList, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	appID := hostList.ApplicationID
	if appID == 0 {
		// get default app id
		var err error
		appID, err = srvData.lgc.GetDefaultAppIDWithSupplier(srvData.ctx)
		if err != nil {
			blog.Errorf("add host, but get default app id failed, err: %v,input:%+v,rid:%s", err, hostList, srvData.rid)
			_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
			return
		}
	}

	// 获取目标业务空先机模块ID
	cond := hutil.NewOperation().WithModuleName(common.DefaultResModuleName).WithAppID(appID).MapStr()
	cond.Set(common.BKDefaultField, common.DefaultResModuleFlag)
	moduleID, err := srvData.lgc.GetResourcePoolModuleID(srvData.ctx, cond)
	if err != nil {
		blog.Errorf("add host, but get module id failed, err: %s,input: %+v,rid: %s", err.Error(), hostList, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	hostIDs, success, updateErrRow, errRow, err := srvData.lgc.AddHost(srvData.ctx, appID, []int64{moduleID}, srvData.ownerID, hostList.HostInfo, hostList.InputType)
	retData := make(map[string]interface{})
	if err != nil {
		blog.Errorf("add host failed, success: %v, update: %v, err: %v, %v,input:%+v,rid:%s", success, updateErrRow, err, errRow, hostList, srvData.rid)
		retData["error"] = errRow
		retData["update_error"] = updateErrRow
		_ = resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{Result: false, Code: common.CCErrHostCreateFail, ErrMsg: srvData.ccErr.Error(common.CCErrHostCreateFail).Error()},
			Data:     retData,
		})
		return
	}
	retData["success"] = success

	// auth: register hosts
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDs...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid: %s", hostIDs, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(retData))
}

// Deprecated:
func (s *Service) AddHostFromAgent(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	agents := new(meta.AddHostFromAgentHostList)
	if err := json.NewDecoder(req.Request.Body).Decode(&agents); err != nil {
		blog.Errorf("add host from agent failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(agents.HostInfo) == 0 {
		blog.Errorf("add host from agent, but got 0 agents from body.input:%+v,rid:%s", agents, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "HostInfo")})
		return
	}

	appID, err := srvData.lgc.GetDefaultAppID(srvData.ctx)
	if err != nil {
		blog.Errorf("AddHostFromAgent GetDefaultAppID error.input:%#v,rid:%s", agents, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}
	if 0 == appID {
		blog.Errorf("add host from agent, but got invalid default appID, err: %v,ownerID:%s,input:%#v,rid:%s", err, srvData.ownerID, agents, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrAddHostToModule, "business not found")})
		return
	}

	// check authorization
	// is AddHostFromAgent's authentication the same with common api?
	// auth: check authorization
	// if err := s.AuthManager.AuthorizeCreateHost(srvData.ctx, srvData.header, appID); err != nil {
	// 	blog.Errorf("check add host authorization failed, business: %+v, err: %v, rid: %s", appID, err, srvData.rid)
	// 	resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
	// 	return
	// }

	opt := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(appID)
	moduleID, err := srvData.lgc.GetResourcePoolModuleID(srvData.ctx, opt.MapStr())
	if err != nil {
		blog.Errorf("add host from agent , but get module id failed, err: %v,ownerID:%s,input:%+v,rid:%s", err, srvData.ownerID, agents, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	agents.HostInfo["import_from"] = common.HostAddMethodAgent
	addHost := make(map[int64]map[string]interface{})
	addHost[1] = agents.HostInfo

	hostIDs, success, updateErrRow, errRow, err := srvData.lgc.AddHost(srvData.ctx, appID, []int64{moduleID}, common.BKDefaultOwnerID, addHost, "")
	if err != nil {
		blog.Errorf("add host failed, success: %v, update: %v, err: %v, %v,input:%+v,rid:%s", success, updateErrRow, err, errRow, agents, srvData.rid)

		retData := make(map[string]interface{})
		retData["success"] = success
		retData["error"] = errRow
		retData["update_error"] = updateErrRow
		_ = resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{Result: false, Code: common.CCErrHostCreateFail, ErrMsg: srvData.ccErr.Error(common.CCErrHostCreateFail).Error()},
			Data:     retData,
		})
		return
	}

	// register hosts
	// auth: register hosts
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDs...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid: %s", hostIDs, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(success))
}

func (s *Service) SearchHost(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	body := new(meta.HostCommonSearch)
	if err := json.NewDecoder(req.Request.Body).Decode(body); err != nil {
		blog.Errorf("search host failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	host, err := srvData.lgc.SearchHost(srvData.ctx, body, false)
	if err != nil {
		blog.Errorf("search host failed, err: %v,input:%+v,rid:%s", err, body, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostGetFail)})
		return
	}

	hostIDArray := host.ExtractHostIDs()
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, *hostIDArray...); err != nil {
		blog.Errorf("check host authorization failed, hostID: %+v, err: %+v, rid: %s", hostIDArray, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	_ = resp.WriteEntity(meta.SearchHostResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *host,
	})
}

func (s *Service) SearchHostWithAsstDetail(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	body := new(meta.HostCommonSearch)
	if err := json.NewDecoder(req.Request.Body).Decode(body); err != nil {
		blog.Errorf("search host failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	host, err := srvData.lgc.SearchHost(srvData.ctx, body, true)
	if err != nil {
		blog.Errorf("search host failed, err: %v,input:%+v,rid:%s", err, body, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	// auth: check authorization
	hostIDArray := host.ExtractHostIDs()
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, *hostIDArray...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDArray, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	_ = resp.WriteEntity(meta.SearchHostResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *host,
	})
}

func (s *Service) UpdateHostBatch(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	data := mapstr.New()
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update host batch failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	// TODO: this is a wrong usage, just for compatible the wrong usage before.
	// delete this, when the frontend use the right request field. not the number.
	id := data[common.BKHostIDField]
	hostIDStr := ""
	switch id.(type) {
	case float64:
		floatID := id.(float64)
		hostIDStr = strconv.FormatInt(int64(floatID), 10)
	case string:
		hostIDStr = id.(string)
	default:
		blog.Errorf("update host batch failed, got invalid host id(%v) data type,rid:%s", id, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsIsInvalid, "bk_host_id")})
		return
	}

	data.Remove(common.MetadataField)
	data.Remove(common.BKHostIDField)
	hostFields, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, nil)
	if err != nil {
		blog.Errorf("update host batch, but get host attribute for audit failed, err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	// check authorization
	hostIDArr := make([]int64, 0)
	for _, id := range strings.Split(hostIDStr, ",") {
		hostID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			blog.Errorf("update host batch, but got invalid host id[%s], err: %v,rid:%s", id, err, srvData.rid)
			_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
			return
		}
		hostIDArr = append(hostIDArr, hostID)
	}
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Update, hostIDArr...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDArr, err, srvData.rid)
		if err != nil && err != auth.NoAuthorizeError {
			blog.ErrorJSON("check host authorization failed, hosts: %s, err: %s, rid: %s", hostIDArr, err.Error(), srvData.rid)
			_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		perm, err := s.AuthManager.GenEditHostBatchNoPermissionResp(srvData.ctx, srvData.header, authcenter.Edit, hostIDArr)
		if err != nil && err != auth.NoAuthorizeError {
			blog.ErrorJSON("check host authorization get permission failed, hosts: %s, err: %s, rid: %s", hostIDArr, err.Error(), srvData.rid)
			resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		resp.WriteEntity(perm)
		return
	}

	logPreContents := make(map[int64]meta.SaveAuditLogParams, 0)
	for _, id := range strings.Split(hostIDStr, ",") {
		hostID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			blog.Errorf("update host batch, but got invalid host id[%s], err: %v,rid:%s", id, err, srvData.rid)
			_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsInvalid)})
			return
		}
		audit := srvData.lgc.NewHostLog(srvData.ctx, srvData.ownerID)
		if err := audit.WithPrevious(srvData.ctx, id, hostFields); err != nil {
			blog.Errorf("update host batch, but get host[%s] pre data for audit failed, err: %v, rid: %s", id, err, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDetailFail)})
			return
		}

		logPreContents[hostID] = audit.AuditLog(srvData.ctx, hostID)
	}

	opt := &meta.UpdateOption{
		Condition: mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: hostIDArr}},
		Data:      mapstr.NewFromMap(data),
	}
	result, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, opt)
	if err != nil {
		blog.Errorf("UpdateHostBatch UpdateObject http do error, err: %v,input:%+v,param:%+v,rid:%s", err, data, opt, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.ErrorJSON("UpdateHostBatch failed, UpdateObject failed, param:%s, response: %s, rid:%s", opt, result, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	logLastContents := make([]meta.SaveAuditLogParams, 0)
	for _, hostID := range hostIDArr {
		audit := srvData.lgc.NewHostLog(srvData.ctx, common.BKDefaultOwnerID)
		if err := audit.WithPrevious(srvData.ctx, strconv.FormatInt(hostID, 10), hostFields); err != nil {
			blog.Errorf("update host batch, but get host[%v] pre data for audit failed, err: %v, rid: %s", hostID, err, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDetailFail)})
			return
		}
		logContent := audit.Content
		logContent.CurData = logContent.PreData
		preLogContent, ok := logPreContents[hostID]
		if ok {
			content, ok := preLogContent.Content.(*meta.Content)
			if ok {
				logContent.PreData = content.PreData
			}
		}

		hostModuleConfig, err := srvData.lgc.GetConfigByCond(srvData.ctx, meta.HostModuleRelationRequest{HostIDArr: []int64{hostID}})
		if err != nil {
			blog.Errorf("update host batch GetConfigByCond failed, id[%v], err: %v,input:%+v,rid:%s", hostID, err, data, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		var appID int64
		if len(hostModuleConfig) > 0 {
			appID = hostModuleConfig[0].AppID
		}

		logLastContents = append(logLastContents,
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
	auditResp, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, logLastContents...)
	if err != nil {
		blog.Errorf("update host property batch, but add host[%v] audit failed, err: %v, rid:%s", hostIDArr, err, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if !auditResp.Result {
		blog.Errorf("update host property batch, but add host[%v] audit failed, err: %v, rid:%s", hostIDArr, auditResp.ErrMsg, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(auditResp.Code, auditResp.ErrMsg)})
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) UpdateHostPropertyBatch(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	parameter := new(meta.UpdateHostPropertyBatchParameter)
	if err := json.NewDecoder(req.Request.Body).Decode(&parameter); err != nil {
		blog.Errorf("update host property batch failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(parameter.Update) > common.BKMaxPageSize {
		blog.Errorf("UpdateHostPropertyBatch failed, data len %d exceed max pageSize %d, rid:%s", len(parameter.Update), common.BKMaxPageSize, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommXXExceedLimit, "update", common.BKMaxPageSize)})
		return
	}

	hostFields, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, nil)
	if err != nil {
		blog.Errorf("update host property batch, but get host attribute for audit failed, err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	// check authorization
	hostIDArr := make([]int64, 0)
	for _, update := range parameter.Update {
		hostIDArr = append(hostIDArr, update.HostID)
	}
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Update, hostIDArr...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDArr, err, srvData.rid)
		if err != nil && err != auth.NoAuthorizeError {
			blog.ErrorJSON("check host authorization failed, hosts: %s, err: %s, rid: %s", hostIDArr, err.Error(), srvData.rid)
			_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		perm, err := s.AuthManager.GenEditHostBatchNoPermissionResp(srvData.ctx, srvData.header, authcenter.Edit, hostIDArr)
		if err != nil && err != auth.NoAuthorizeError {
			blog.ErrorJSON("check host authorization get permission failed, hosts: %s, err: %s, rid: %s", hostIDArr, err.Error(), srvData.rid)
			resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		resp.WriteEntity(perm)
		return
	}

	auditLogs := make([]meta.SaveAuditLogParams, 0)
	for _, update := range parameter.Update {
		id := strconv.FormatInt(update.HostID, 10)
		cond := mapstr.New()
		cond.Set(common.BKHostIDField, update.HostID)
		data, err := mapstr.NewFromInterface(update.Properties)
		if err != nil {
			blog.Errorf("update host property batch, but convert properties[%v] to mapstr failed, err: %v, rid: %s", update.Properties, err, srvData.rid)
			_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
			return
		}
		// can't update host's cloud area using this api
		data.Remove(common.BKCloudIDField)
		data.Remove(common.BKHostIDField)
		opt := &meta.UpdateOption{
			Condition: cond,
			Data:      data,
		}
		hostLog := srvData.lgc.NewHostLog(srvData.ctx, srvData.ownerID)
		if err := hostLog.WithPrevious(srvData.ctx, id, hostFields); err != nil {
			blog.Errorf("update host property batch, but get host[%s] pre data for audit failed, err: %v, rid: %s", id, err, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		result, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, opt)
		if err != nil {
			blog.Errorf("UpdateHostPropertyBatch UpdateInstance http do error, err: %v,input:%+v,param:%+v,rid:%s", err, data, opt, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		if !result.Result {
			blog.Errorf("UpdateHostPropertyBatch UpdateObject http response error, err code:%d,err msg:%s,input:%+v,param:%+v,rid:%s", result.Code, data, opt, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
			return
		}

		if err := hostLog.WithCurrent(srvData.ctx, id); err != nil {
			blog.Errorf("update host property batch, but get host[%s] pre data for audit failed, err: %v, rid: %s", id, err, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}

		hostModuleConfig, err := srvData.lgc.GetConfigByCond(srvData.ctx, meta.HostModuleRelationRequest{HostIDArr: []int64{update.HostID}})
		if err != nil {
			blog.Errorf("update host property batch GetConfigByCond failed, hostID[%v], err: %v,rid:%s", update.HostID, err, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
			return
		}
		var appID int64
		if len(hostModuleConfig) > 0 {
			appID = hostModuleConfig[0].AppID
		}
		auditLog := hostLog.AuditLog(srvData.ctx, update.HostID)
		auditLog.Model = common.BKInnerObjIDHost
		auditLog.OpType = auditoplog.AuditOpTypeModify
		auditLog.BizID = appID
		auditLog.OpDesc = "update host"
		auditLogs = append(auditLogs, auditLog)
	}

	auditResp, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, auditLogs...)
	if err != nil {
		blog.Errorf("update host property batch, but add host[%v] audit failed, err: %v, rid:%s", hostIDArr, err, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if !auditResp.Result {
		blog.Errorf("update host property batch, but add host[%v] audit failed, err: %v, rid:%s", hostIDArr, auditResp.ErrMsg, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(auditResp.Code, auditResp.ErrMsg)})
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(nil))
}

// NewHostSyncAppTopo add new hosts to the business
// synchronize hosts directly to a module in a business if this host does not exist.
// otherwise, this operation will only change host's attribute.
// TODO: used by framework.
func (s *Service) NewHostSyncAppTopo(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	hostList := new(meta.HostSyncList)
	if err := json.NewDecoder(req.Request.Body).Decode(hostList); err != nil {
		blog.Errorf("add host failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if hostList.HostInfo == nil {
		blog.Errorf("add host, but host info is nil.input:%+v,rid:%s", hostList, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "host_info")})
		return
	}
	if 0 == len(hostList.ModuleID) {
		blog.Errorf("host sync app  parameters required moduleID,input:%+v,rid:%s", hostList, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, common.BKModuleIDField)})
		return
	}

	if common.BatchHostAddMaxRow < len(hostList.HostInfo) {
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommXXExceedLimit, "host_info ", common.BatchHostAddMaxRow)})
		return
	}

	appConds := map[string]interface{}{
		common.BKAppIDField: hostList.ApplicationID,
	}
	appInfo, err := srvData.lgc.GetAppDetails(srvData.ctx, "", appConds)
	if nil != err {
		blog.Errorf("host sync app %d error:%s,input:%+v,rid:%s", hostList.ApplicationID, err.Error(), hostList, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if 0 == len(appInfo) {
		blog.Errorf("host sync app %d not found, reply:%+v,input:%+v,rid:%s", hostList.ApplicationID, appInfo, hostList, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrTopoGetAppFailed)})
		return
	}

	moduleCond := []meta.ConditionItem{
		{
			Field:    common.BKModuleIDField,
			Operator: common.BKDBIN,
			Value:    hostList.ModuleID,
		},
	}
	if len(hostList.ModuleID) > 1 {
		moduleCond = append(moduleCond, meta.ConditionItem{
			Field:    common.BKDefaultField,
			Operator: common.BKDBEQ,
			Value:    common.DefaultFlagDefaultValue,
		})
	}
	// srvData.lgc..NewHostSyncValidModule(req, data.ApplicationID, data.ModuleID, m.CC.ObjCtrl())
	moduleIDS, err := srvData.lgc.GetModuleIDByCond(srvData.ctx, moduleCond)
	if nil != err {
		blog.Errorf("NewHostSyncAppTop GetModuleIDByCond error. err:%s,input:%+v,rid:%s", err.Error(), hostList, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if len(moduleIDS) != len(hostList.ModuleID) {
		blog.Errorf("not found part module: source:%v, db:%v, rid: %s", hostList.ModuleID, moduleIDS, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleIDNotFoundORHasMultipleInnerModuleIDFailed)})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeCreateHost(srvData.ctx, srvData.header, hostList.ApplicationID); err != nil {
		blog.Errorf("check add hosts authorization failed, business: %d, err: %v, rid: %s", hostList.ApplicationID, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	hostIDs, success, updateErrRow, errRow, err := srvData.lgc.AddHost(srvData.ctx, hostList.ApplicationID, hostList.ModuleID, srvData.ownerID, hostList.HostInfo, common.InputTypeApiNewHostSync)
	if err != nil {
		blog.Errorf("add host failed, success: %v, update: %v, err: %v, %v, rid: %s", success, updateErrRow, err, errRow, srvData.rid)

		retData := make(map[string]interface{})
		retData["success"] = success
		retData["error"] = errRow
		retData["update_error"] = updateErrRow
		_ = resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{Result: false, Code: common.CCErrHostCreateFail, ErrMsg: srvData.ccErr.Error(common.CCErrHostCreateFail).Error()},
			Data:     retData,
		})
		return
	}

	// register host to iam
	// auth: check authorization
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDs...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid: %s", hostIDs, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(success))
}

// MoveSetHost2IdleModule bk_set_id and bk_module_id cannot be empty at the same time
// Remove the host from the module or set.
// The host belongs to the current module or host only, and puts the host into the idle machine of the current service.
// When the host data is in multiple modules or sets. Disconnect the host from the module or set only
// TODO: used by v2 version, remove this api when v2 is offline.
func (s *Service) MoveSetHost2IdleModule(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	srvData := s.newSrvComm(req.Request.Header)

	var data meta.SetHostConfigParams
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("MoveSetHost2IdleModule failed with decode body err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if 0 == data.ApplicationID {
		blog.Errorf("MoveSetHost2IdleModule bk_biz_id cannot be empty at the same time,input:%#v,rid:%s", data, util.GetHTTPCCRequestID(header))
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	if 0 == data.SetID && 0 == data.ModuleID {
		blog.Errorf("MoveSetHost2IdleModule bk_set_id and bk_module_id cannot be empty at the same time,input:%#v,rid:%s", data, util.GetHTTPCCRequestID(header))
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	// get host in set
	var condition meta.HostModuleRelationRequest
	hostIDArr := make([]int64, 0)

	if 0 != data.SetID {
		condition.SetIDArr = []int64{data.SetID}
	}
	if 0 != data.ModuleID {
		condition.ModuleIDArr = []int64{data.ModuleID}
	}

	condition.ApplicationID = data.ApplicationID
	hostResult, err := srvData.lgc.GetConfigByCond(srvData.ctx, condition)
	if nil != err {
		blog.Errorf("read host from application  error:%v,input:%+v,rid:%s", err, data, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	if 0 == len(hostResult) {
		blog.Warnf("no host in set,rid:%s", srvData.rid)
		_ = resp.WriteEntity(meta.NewSuccessResp(nil))
		return
	}
	for _, cell := range hostResult {
		hostIDArr = append(hostIDArr, cell.HostID)
	}
	moduleCond := []meta.ConditionItem{
		{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    data.ApplicationID,
		},
		{
			Field:    common.BKDefaultField,
			Operator: common.BKDBEQ,
			Value:    common.DefaultResModuleFlag,
		},
	}

	moduleIDArr, err := srvData.lgc.GetModuleIDByCond(srvData.ctx, moduleCond)
	if err != nil {
		blog.Errorf("MoveSetHost2IdleModule GetModuleIDByCond error. err:%s, input:%#v, param:%#v, rid:%s", err.Error(), data, moduleCond, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if len(moduleIDArr) == 0 {
		blog.Errorf("MoveSetHost2IdleModule GetModuleIDByCond error. err:%s, input:%#v, param:%#v, rid:%s", err.Error(), data, moduleCond, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleNotExist, "idle module")})
		return
	}
	idleModuleID := moduleIDArr[0]
	moduleHostConfigParams := make(map[string]interface{})
	moduleHostConfigParams[common.BKAppIDField] = data.ApplicationID
	audit := srvData.lgc.NewHostModuleLog(hostIDArr)

	// auth: check authorization
	// if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.MoveHostToBizIdleModule, hostIDArr...); err != nil {
	// 	blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDArr, err, srvData.rid)
	// 	resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
	// 	return
	// }
	// // step2. check permission for target business
	// if err := s.AuthManager.AuthorizeCreateHost(srvData.ctx, srvData.header, data.ApplicationID); err != nil {
	// 	blog.Errorf("check add host authorization failed, business: %d, err: %v, rid: %s", data.ApplicationID, err, srvData.rid)
	// 	resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
	// 	return
	// }
	// step3. deregister host from iam
	if err := s.AuthManager.DeregisterHostsByID(srvData.ctx, srvData.header, hostIDArr...); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v, rid: %s", hostIDArr, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)})
		return
	}
	hmInput := &meta.HostModuleRelationRequest{
		ApplicationID: data.ApplicationID,
		HostIDArr:     hostIDArr,
	}
	configResult, err := srvData.lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(srvData.ctx, srvData.header, hmInput)
	if nil != err {
		blog.Errorf("remove hostModuleConfig, http do error, error:%v, params:%v, input:%+v, rid:%s", err, hmInput, data, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	if !configResult.Result {
		blog.Errorf("remove hostModuleConfig http reply error, result:%v, params:%v, input:%+v, rid:%s", configResult, hmInput, data, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	hostIDMHMap := make(map[int64][]meta.ModuleHost, 0)
	for _, item := range configResult.Data.Info {
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
			opResult, err = srvData.lgc.CoreAPI.CoreService().Host().TransferToInnerModule(srvData.ctx, srvData.header, input)
		} else {
			input := &meta.HostsModuleRelation{
				ApplicationID: data.ApplicationID,
				HostID:        []int64{hostID},
				ModuleID:      newModuleIDArr,
			}
			opResult, err = srvData.lgc.CoreAPI.CoreService().Host().TransferToNormalModule(srvData.ctx, srvData.header, input)
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

	if err := audit.SaveAudit(srvData.ctx, data.ApplicationID, srvData.user, "host to empty module"); err != nil {
		blog.Errorf("SaveAudit failed, err: %s, rid: %s", err.Error(), srvData.rid)
	}

	// register host to iam
	// auth: check authorization
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDArr...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid:%s", hostIDArr, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
		return
	}
	if len(exceptionArr) > 0 {
		blog.Errorf("MoveSetHost2IdleModule has exception. exception:%#v, rid:%s", exceptionArr, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail), Data: exceptionArr})
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(nil))
	return
}

func (s *Service) ip2hostID(srvData *srvComm, ip string, cloudID int64) (hostID int64, err error) {
	_, hostID, err = srvData.lgc.IPCloudToHost(srvData.ctx, ip, cloudID)
	return hostID, err
}

// CloneHostProperty clone host property from src host to dst host
func (s *Service) CloneHostProperty(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := &meta.CloneHostPropertyParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("CloneHostProperty , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if input.OrgIP == input.DstIP {
		result := meta.Response{
			BaseResp: meta.SuccessBaseResp,
			Data:     nil,
		}
		_ = resp.WriteEntity(result)
		return
	}

	if 0 == input.AppID {
		blog.Errorf("CloneHostProperty, application not found input:%+v,rid:%s", input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID")})
		return
	}
	if input.OrgIP == "" {
		blog.Errorf("CloneHostProperty, OrgIP not found input:%+v,rid:%s", input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "bk_org_ip")})
		return
	}
	if input.DstIP == "" {
		blog.Errorf("CloneHostProperty, OrgIP not found input:%+v,rid:%s", input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "bk_dst_ip")})
		return
	}

	// authorization check
	srcHostID, err := s.ip2hostID(srvData, input.OrgIP, input.CloudID)
	if err != nil {
		blog.Errorf("ip2hostID failed, ip:%s, input:%+v, rid:%s", input.OrgIP, input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, "OrgIP")})
		return
	}
	// check source host exist
	if srcHostID == 0 {
		blog.Errorf("host not found. params:%s,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.CCErrorf(common.CCErrHostNotFound)})
		return
	}
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, srcHostID); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid:%s", srcHostID, err, srvData.rid)
		_ = resp.WriteError(http.StatusForbidden, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}
	// step2. verify has permission to update dst host
	dstHostID, err := s.ip2hostID(srvData, input.DstIP, input.CloudID)
	if err != nil {
		blog.Errorf("ip2hostID failed, ip:%s, input:%+v, rid:%s", input.DstIP, input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, "DstIP")})
		return
	}
	// check whether destination host exist
	if dstHostID == 0 {
		blog.Errorf("host not found. params:%s,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.CCErrorf(common.CCErrHostNotFound)})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Update, dstHostID); err != nil {
		if err != auth.NoAuthorizeError {
			blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid:%s", dstHostID, err, srvData.rid)
			resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		perm, err := s.AuthManager.GenEditBizHostNoPermissionResp(srvData.ctx, srvData.header, []int64{dstHostID})
		if err != nil {
			resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		resp.WriteEntity(perm)
		return
	}

	err = srvData.lgc.CloneHostProperty(srvData.ctx, input.AppID, srcHostID, dstHostID)
	if nil != err {
		blog.Errorf("CloneHostProperty  error , err: %v, input:%#v, rid:%s", err, input, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	result := meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     nil,
	}
	_ = resp.WriteEntity(result)
}
