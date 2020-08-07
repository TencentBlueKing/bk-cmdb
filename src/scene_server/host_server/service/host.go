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
	"configcenter/src/auth/extensions"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
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
	delCondsArr := make([][]map[string]interface{}, 0)
	for _, i := range hostIDArr {
		iHostID, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			blog.Errorf("delete host batch, but got invalid host id, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
			_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, iHostID)})
			return
		}
		iHostIDArr = append(iHostIDArr, iHostID)
	}
	iHostIDArr = util.IntArrayUnique(iHostIDArr)

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Delete, iHostIDArr...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", iHostIDArr, err, srvData.rid)
		if err != auth.NoAuthorizeError {
			_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail)})
			return
		}
		perm, err := s.AuthManager.GenEditHostBatchNoPermissionResp(srvData.ctx, srvData.header, authcenter.Delete, iHostIDArr)
		if err != nil {
			_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDeleteFail)})
			return
		}
		_ = resp.WriteEntity(perm)
		return
	}

	for _, iHostID := range iHostIDArr {
		asstCond := map[string]interface{}{
			common.BKDBOR: []map[string]interface{}{
				{
					common.BKObjIDField:  common.BKInnerObjIDHost,
					common.BKHostIDField: iHostID,
				},
				{
					common.BKAsstObjIDField:  common.BKInnerObjIDHost,
					common.BKAsstInstIDField: iHostID,
				},
			},
		}
		rsp, err := s.CoreAPI.CoreService().Association().ReadInstAssociation(srvData.ctx, srvData.header, &meta.QueryCondition{Condition: asstCond})
		if nil != err {
			blog.ErrorJSON("DeleteHostBatch read host association do request failed , err: %s, rid: %s", err.Error(), srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)})
			return
		}
		if !rsp.Result {
			blog.ErrorJSON("DeleteHostBatch read host association failed , err message: %s, rid: %s", rsp.ErrMsg, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: rsp.CCError()})
			return
		}
		if rsp.Data.Count <= 0 {
			continue
		}
		objIDs := make([]string, 0)
		asstInstMap := make(map[string][]int64, 0)
		for _, asst := range rsp.Data.Info {
			if asst.ObjectID == common.BKInnerObjIDHost && iHostID == asst.InstID {
				objIDs = append(objIDs, asst.AsstObjectID)
				asstInstMap[asst.AsstObjectID] = append(asstInstMap[asst.AsstObjectID], asst.AsstInstID)
			} else if asst.AsstObjectID == common.BKInnerObjIDHost && iHostID == asst.AsstInstID {
				objIDs = append(objIDs, asst.ObjectID)
				asstInstMap[asst.ObjectID] = append(asstInstMap[asst.ObjectID], asst.InstID)
			} else {
				_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(common.CCErrCommDBSelectFailed, "host is not associated in selected association")})
				return
			}
		}
		delConds := make([]map[string]interface{}, 0)
		for objID, instIDs := range asstInstMap {
			if len(instIDs) < 0 {
				continue
			}
			instIDField := common.GetInstIDField(objID)
			instCond := map[string]interface{}{
				instIDField: map[string]interface{}{
					common.BKDBIN: instIDs,
				},
			}
			instRsp, err := s.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, objID, &meta.QueryCondition{Condition: instCond})
			if err != nil {
				blog.ErrorJSON("DeleteHostBatch read associated instances do request failed , err: %s, rid: %s", err.Error(), srvData.rid)
				_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)})
				return
			}
			if !instRsp.Result {
				blog.ErrorJSON("DeleteHostBatch read associated instances failed , err message: %s, rid: %s", instRsp.ErrMsg, srvData.rid)
				_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: instRsp.CCError()})
				return
			}
			if len(instRsp.Data.Info) > 0 {
				blog.ErrorJSON("DeleteHostBatch host %s has been associated, can't be deleted, rid: %s", iHostID, srvData.rid)
				_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.CCErrorf(common.CCErrTopoInstHasBeenAssociation, iHostID)})
				return
			}
			delConds = append(delConds, map[string]interface{}{
				common.BKObjIDField: objID,
				instIDField: map[string]interface{}{
					common.BKDBIN: instIDs,
				},
			}, map[string]interface{}{
				common.AssociatedObjectIDField: objID,
				instIDField: map[string]interface{}{
					common.BKDBIN: instIDs,
				},
			})
		}
		if len(delConds) > 0 {
			delCondsArr = append(delCondsArr, delConds)
		}
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		for _, delConds := range delCondsArr {
			delRsp, err := s.CoreAPI.CoreService().Association().DeleteInstAssociation(srvData.ctx, srvData.header, &meta.DeleteOption{Condition: map[string]interface{}{common.BKDBOR: delConds}})
			if err != nil {
				blog.ErrorJSON("DeleteHostBatch delete host redundant association do request failed , err: %s, rid: %s", err.Error(), srvData.rid)
				return srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
			}
			if !delRsp.Result {
				blog.ErrorJSON("DeleteHostBatch delete host redundant association failed , err message: %s, rid: %s", delRsp.ErrMsg, srvData.rid)
				return delRsp.CCError()
			}
		}
		appID, err := srvData.lgc.GetDefaultAppID(srvData.ctx)
		if err != nil {
			blog.Errorf("delete host batch, but got invalid app id, err: %v,input:%s,rid:%s", err, opt, srvData.rid)
			return srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
		}

		hostFields, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, meta.BizLabelNotExist)
		if err != nil {
			blog.Errorf("delete host batch failed, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
			return err
		}

		logContentMap := make(map[int64]meta.AuditLog, 0)
		hosts := make([]extensions.HostSimplify, 0)
		for _, hostID := range iHostIDArr {
			logger := srvData.lgc.NewHostLog(srvData.ctx, srvData.ownerID)
			if err := logger.WithPrevious(srvData.ctx, hostID, hostFields); err != nil {
				blog.Errorf("delete host batch, but get pre host data failed, err: %v,input:%+v,rid:%s", err, opt, srvData.rid)
				return err
			}

			logContentMap[hostID], err = logger.AuditLog(srvData.ctx, hostID, appID, meta.AuditDelete)
			if err != nil {
				blog.Errorf("delete host batch, but get host[%d] biz[%d] data failed, err: %v, rid:%s", hostID, appID, err, srvData.rid)
				return err
			}

			detail, ok := logContentMap[hostID].OperationDetail.(*meta.InstanceOpDetail)
			if !ok {
				blog.Errorf("delete host batch, but got invalid operation detail, rid:%s", srvData.rid)
				return errors.New(common.CCErrCommParamsValueInvalidError, "")
			}

			hosts = append(hosts, extensions.HostSimplify{
				BKAppIDField:       0,
				BKHostIDField:      hostID,
				BKHostInnerIPField: detail.ResourceName,
			})
		}

		input := &meta.DeleteHostRequest{
			ApplicationID: appID,
			HostIDArr:     iHostIDArr,
		}
		delResult, err := s.CoreAPI.CoreService().Host().DeleteHostFromSystem(srvData.ctx, srvData.header, input)
		if err != nil {
			blog.Error("DeleteHostBatch DeleteHost http do error. err:%s, input:%s, rid:%s", err.Error(), input, srvData.rid)
			return srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !delResult.Result {
			blog.Errorf("DeleteHostBatch DeleteHost http reply error. result: %#v, input:%#v, rid:%s", delResult, input, srvData.rid)
			return srvData.ccErr.Error(common.CCErrHostDeleteFail)
		}

		// auth: unregister hosts
		if err := s.AuthManager.DeregisterHosts(srvData.ctx, srvData.header, hosts...); err != nil {
			blog.ErrorJSON("deregister host from iam failed, hosts: %s, err: %s, rid: %s", hosts, err, srvData.rid)
			return srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)
		}

		// ensure delete host add log
		for _, ex := range delResult.Data {
			delete(logContentMap, ex.OriginIndex)
		}
		var logContents []meta.AuditLog
		for _, item := range logContentMap {
			logContents = append(logContents, item)
		}
		if len(logContents) > 0 {
			auditResult, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, logContents...)
			if err != nil || !auditResult.Result {
				blog.ErrorJSON("delete host in batch, but add host audit log failed, err: %s, result: %s,rid:%s", err, auditResult, srvData.rid)
				return srvData.ccErr.Error(common.CCErrAuditSaveLogFailed)
			}
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr})
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

	details, _, err := srvData.lgc.GetHostInstanceDetails(srvData.ctx, hostIDInt64)
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

// HostSnapInfoBatch get the host snapshot in batch
func (s *Service) HostSnapInfoBatch(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	option := meta.SearchInstBatchOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("HostSnapInfoBatch failed, decode body err: %v, rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	rawErr := option.Validate()
	if rawErr.ErrCode != 0 {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: rawErr.ToCCError(srvData.ccErr)})
		return
	}

	hostIDs := util.IntArrayUnique(option.IDs)

	// check authorization
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.Find, hostIDs...); err != nil {
		blog.Errorf("check host authorization failed, hostIDs: %#v, err: %v, rid: %s", hostIDs, err, srvData.rid)
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		return
	}

	input := meta.HostSnapBatchInput{HostIDs: hostIDs}
	// get snapshot
	result, err := s.CoreAPI.CoreService().Host().GetHostSnapBatch(srvData.ctx, srvData.header, input)
	if err != nil {
		blog.Errorf("HostSnapInfoBatch failed, http do error, err: %v ,input:%#v, rid:%s", err, input, srvData.rid)
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPReadBodyFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("HostSnapInfoBatch failed, http response error, err code:%d, err msg:%s, input:%#v, rid:%s", result.Code, result.ErrMsg, input, srvData.rid)
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	ret := make([]map[string]interface{}, 0)
	for hostID, snapData := range result.Data {
		if snapData == "" {
			blog.Infof("snapData is empty, hostID:%v, rid:%s", hostID, srvData.rid)
			ret = append(ret, map[string]interface{}{"bk_host_id": hostID})
			continue
		}
		snap, err := logics.ParseHostSnap(snapData)
		if err != nil {
			blog.Errorf("HostSnapInfoBatch failed, ParseHostSnap err: %v, hostID:%v, rid:%s", err, hostID, srvData.rid)
			_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: err})
			return
		}
		snapFields := make(map[string]interface{})
		for _, field := range option.Fields {
			if _, ok := snap[field]; ok {
				snapFields[field] = snap[field]
			}
		}
		snapFields["bk_host_id"] = hostID
		ret = append(ret, snapFields)
	}

	responseData := meta.HostSnapBatchResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     ret,
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

	retData := make(map[string]interface{})
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		hostIDs, success, updateErrRow, errRow, err := srvData.lgc.AddHost(srvData.ctx, appID, []int64{moduleID}, srvData.ownerID, hostList.HostInfo, hostList.InputType)
		if err != nil {
			blog.Errorf("add host failed, success: %v, update: %v, err: %v, %v,input:%+v,rid:%s", success, updateErrRow, err, errRow, hostList, srvData.rid)
			retData["error"] = errRow
			retData["update_error"] = updateErrRow
			return srvData.ccErr.Error(common.CCErrHostCreateFail)
		}
		retData["success"] = success

		// auth: register hosts
		if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDs...); err != nil {
			blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid: %s", hostIDs, err, srvData.rid)
			return srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr, Data: retData})
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
	var hostIDs []int64
	var success, updateErrRow, errRow []string
	retData := make(map[string]interface{})
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		var err error
		hostIDs, success, updateErrRow, errRow, err = srvData.lgc.AddHost(srvData.ctx, appID, []int64{moduleID}, common.BKDefaultOwnerID, addHost, "")
		if err != nil {
			blog.Errorf("add host failed, success: %v, update: %v, err: %v, %v,input:%+v,rid:%s", success, updateErrRow, err, errRow, agents, srvData.rid)

			retData["success"] = success
			retData["error"] = errRow
			retData["update_error"] = updateErrRow
			return srvData.ccErr.Error(common.CCErrHostCreateFail)
		}

		// register hosts
		// auth: register hosts
		if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDs...); err != nil {
			blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid: %s", hostIDs, err, srvData.rid)
			return srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr, Data: retData})
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

func (s *Service) getHostApplyRelatedFields(srvData *srvComm, hostIDArr []int64) (hostProperties map[int64][]string, hasRules bool, ccErr errors.CCErrorCoder) {
	// filter fields locked by host apply rule
	listRuleOption := meta.ListHostRelatedApplyRuleOption{
		HostIDs: hostIDArr,
		Page: meta.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostRules, ccErr := s.listHostRelatedApplyRule(srvData, 0, listRuleOption)
	if ccErr != nil {
		blog.Errorf("update host batch, listHostRelatedApplyRule failed, option: %+v, err: %v, rid: %s", listRuleOption, ccErr, srvData.rid)
		return nil, false, ccErr
	}
	attributeIDs := make([]int64, 0)
	for _, rules := range hostRules {
		for _, rule := range rules {
			attributeIDs = append(attributeIDs, rule.AttributeID)
		}
	}
	if len(attributeIDs) == 0 {
		return nil, false, nil
	}
	hostAttributesFilter := &meta.QueryCondition{
		Fields: []string{common.BKPropertyIDField, common.BKFieldID},
		Page: meta.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKFieldID: map[string]interface{}{
				common.BKDBIN: attributeIDs,
			},
		},
	}
	attributeResult, err := s.CoreAPI.CoreService().Model().ReadModelAttr(srvData.ctx, srvData.header, common.BKInnerObjIDHost, hostAttributesFilter)
	if err != nil {
		blog.Errorf("UpdateHostBatch failed, ReadModelAttr failed, param: %+v, err: %+v, rid:%s", hostAttributesFilter, err, srvData.rid)
		return nil, true, srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if ccErr := attributeResult.CCError(); ccErr != nil {
		blog.Errorf("UpdateHostBatch failed, ReadModelAttr failed, param: %+v, output: %+v, rid:%s", hostAttributesFilter, attributeResult, srvData.rid)
		return nil, true, ccErr
	}
	attributeMap := make(map[int64]meta.Attribute)
	for _, item := range attributeResult.Data.Info {
		attributeMap[item.ID] = item
	}
	hostProperties = make(map[int64][]string)
	for hostID, rules := range hostRules {
		if _, exist := hostProperties[hostID]; exist == false {
			hostProperties[hostID] = make([]string, 0)
		}
		for _, rule := range rules {
			attribute, ok := attributeMap[rule.AttributeID]
			if ok == false {
				continue
			}
			hostProperties[hostID] = append(hostProperties[hostID], attribute.PropertyID)
		}
	}
	return hostProperties, true, nil
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
	hostFields, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, meta.BizLabelNotExist)
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

	logPreContents := make(map[int64]*logics.HostLog, 0)
	for _, hostID := range hostIDArr {
		audit := srvData.lgc.NewHostLog(srvData.ctx, srvData.ownerID)
		if err := audit.WithPrevious(srvData.ctx, hostID, hostFields); err != nil {
			blog.Errorf("update host batch, but get host[%s] pre data for audit failed, err: %v, rid: %s", id, err, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostDetailFail)})
			return
		}

		logPreContents[hostID] = audit
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		hasHostUpdateWithoutHostApplyFiled := false
		// 功能开关：更新主机属性时是否剔除自动应用字段
		if meta.HostUpdateWithoutHostApplyFiled == true {
			hostProperties, hasRules, err := s.getHostApplyRelatedFields(srvData, hostIDArr)
			if err != nil {
				blog.Errorf("UpdateHostBatch failed, getHostApplyRelatedFields failed, hostIDArr: %+v, err: %v, rid:%s", hostIDArr, err, srvData.rid)
				return err
			}
			// get host attributes
			if hasRules == true {
				hasHostUpdateWithoutHostApplyFiled = true
				for _, hostID := range hostIDArr {
					updateData := make(map[string]interface{})
					for key, value := range data {
						properties, ok := hostProperties[hostID]
						if ok == true && util.InStrArr(properties, key) {
							continue
						}
						updateData[key] = value
					}
					opt := &meta.UpdateOption{
						Condition: mapstr.MapStr{common.BKHostIDField: hostID},
						Data:      mapstr.NewFromMap(updateData),
					}
					result, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, opt)
					if err != nil {
						blog.Errorf("UpdateHostBatch UpdateObject http do error, err: %v,input:%+v,param:%+v,rid:%s", err, data, opt, srvData.rid)
						return srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
					}
					if !result.Result {
						blog.ErrorJSON("UpdateHostBatch failed, UpdateObject failed, param:%s, response: %s, rid:%s", opt, result, srvData.rid)
						return srvData.ccErr.New(result.Code, result.ErrMsg)
					}
				}
			}
		}

		if hasHostUpdateWithoutHostApplyFiled == false {
			// 退化到批量编辑
			opt := &meta.UpdateOption{
				Condition: mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: hostIDArr}},
				Data:      mapstr.NewFromMap(data),
			}
			result, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, opt)
			if err != nil {
				blog.Errorf("UpdateHostBatch UpdateObject http do error, err: %v,input:%+v,param:%+v,rid:%s", err, data, opt, srvData.rid)
				return srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
			}
			if !result.Result {
				blog.ErrorJSON("UpdateHostBatch failed, UpdateObject failed, param:%s, response: %s, rid:%s", opt, result, srvData.rid)
				return srvData.ccErr.New(result.Code, result.ErrMsg)
			}
		}

		hostModuleConfig, err := srvData.lgc.GetConfigByCond(srvData.ctx, meta.HostModuleRelationRequest{HostIDArr: hostIDArr, Fields: []string{common.BKAppIDField, common.BKHostIDField}})
		if err != nil {
			blog.Errorf("update host batch GetConfigByCond failed, hostIDArr[%v], err: %v,input:%+v,rid:%s", hostIDArr, err, data, srvData.rid)
			return err
		}
		appIDMap := make(map[int64]int64)
		for _, hostModule := range hostModuleConfig {
			appIDMap[hostModule.HostID] = hostModule.AppID
		}

		logLastContents := make([]meta.AuditLog, 0)
		for _, hostID := range hostIDArr {
			audit, ok := logPreContents[hostID]
			if !ok {
				audit = srvData.lgc.NewHostLog(srvData.ctx, common.BKDefaultOwnerID)
			}
			if err := audit.WithCurrent(srvData.ctx, hostID, hostFields); err != nil {
				blog.Errorf("update host batch, but get host[%v] pre data for audit failed, err: %v, rid: %s", hostID, err, srvData.rid)
				return srvData.ccErr.Error(common.CCErrHostDetailFail)
			}
			auditLog, err := audit.AuditLog(srvData.ctx, hostID, appIDMap[hostID], meta.AuditUpdate)
			if err != nil {
				blog.Errorf("update host batch, but get host[%v] biz[%v] data for audit failed, err: %v, rid: %s", hostID, appIDMap[hostID], err, srvData.rid)
				return err
			}
			logLastContents = append(logLastContents, auditLog)
		}
		auditResp, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, logLastContents...)
		if err != nil {
			blog.Errorf("update host property batch, but add host[%v] audit failed, err: %v, rid:%s", hostIDArr, err, srvData.rid)
			return err
		}
		if !auditResp.Result {
			blog.Errorf("update host property batch, but add host[%v] audit failed, err: %v, rid:%s", hostIDArr, auditResp.ErrMsg, srvData.rid)
			return srvData.ccErr.New(auditResp.Code, auditResp.ErrMsg)
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr})
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

	hostFields, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, meta.BizLabelNotExist)
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

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		auditLogs := make([]meta.AuditLog, 0)
		for _, update := range parameter.Update {
			cond := mapstr.New()
			cond.Set(common.BKHostIDField, update.HostID)
			data, err := mapstr.NewFromInterface(update.Properties)
			if err != nil {
				blog.Errorf("update host property batch, but convert properties[%v] to mapstr failed, err: %v, rid: %s", update.Properties, err, srvData.rid)
				return err
			}
			// can't update host's cloud area using this api
			data.Remove(common.BKCloudIDField)
			data.Remove(common.BKHostIDField)
			opt := &meta.UpdateOption{
				Condition: cond,
				Data:      data,
			}
			hostLog := srvData.lgc.NewHostLog(srvData.ctx, srvData.ownerID)
			if err := hostLog.WithPrevious(srvData.ctx, update.HostID, hostFields); err != nil {
				blog.Errorf("update host property batch, but get host[%d] pre data for audit failed, err: %v, rid: %s", update.HostID, err, srvData.rid)
				return err
			}
			result, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, opt)
			if err != nil {
				blog.Errorf("UpdateHostPropertyBatch UpdateInstance http do error, err: %v,input:%+v,param:%+v,rid:%s", err, data, opt, srvData.rid)
				return err
			}
			if !result.Result {
				blog.Errorf("UpdateHostPropertyBatch UpdateObject http response error, err code:%d,err msg:%s,input:%+v,param:%+v,rid:%s", result.Code, data, opt, srvData.rid)
				return srvData.ccErr.New(result.Code, result.ErrMsg)
			}

			if err := hostLog.WithCurrent(srvData.ctx, update.HostID, nil); err != nil {
				blog.Errorf("update host property batch, but get host[%d] pre data for audit failed, err: %v, rid: %s", update.HostID, err, srvData.rid)
				return err
			}

			hostModuleConfig, err := srvData.lgc.GetConfigByCond(srvData.ctx, meta.HostModuleRelationRequest{HostIDArr: []int64{update.HostID}, Fields: []string{common.BKAppIDField}})
			if err != nil {
				blog.Errorf("update host property batch GetConfigByCond failed, hostID[%v], err: %v,rid:%s", update.HostID, err, srvData.rid)
				return err
			}
			var appID int64
			if len(hostModuleConfig) > 0 {
				appID = hostModuleConfig[0].AppID
			}
			auditLog, err := hostLog.AuditLog(srvData.ctx, update.HostID, appID, meta.AuditUpdate)
			if err != nil {
				blog.Errorf("update host property batch, but get host[%d] biz[%d] data for audit failed, err: %v, rid: %s", update.HostID, appID, err, srvData.rid)
				return err
			}
			auditLogs = append(auditLogs, auditLog)
		}

		auditResp, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, auditLogs...)
		if err != nil {
			blog.Errorf("update host property batch, but add host[%v] audit failed, err: %v, rid:%s", hostIDArr, err, srvData.rid)
			return srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !auditResp.Result {
			blog.Errorf("update host property batch, but add host[%v] audit failed, err: %v, rid:%s", hostIDArr, auditResp.ErrMsg, srvData.rid)
			return srvData.ccErr.New(auditResp.Code, auditResp.ErrMsg)
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr})
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

	retData := make(map[string]interface{})
	var hostIDs []int64
	var success, updateErrRow, errRow []string
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		var err error
		hostIDs, success, updateErrRow, errRow, err = srvData.lgc.AddHost(srvData.ctx, hostList.ApplicationID, hostList.ModuleID, srvData.ownerID, hostList.HostInfo, common.InputTypeApiNewHostSync)
		if err != nil {
			blog.Errorf("add host failed, success: %v, update: %v, err: %v, %v, rid: %s", success, updateErrRow, err, errRow, srvData.rid)

			retData["success"] = success
			retData["error"] = errRow
			retData["update_error"] = updateErrRow
			return srvData.ccErr.Error(common.CCErrHostCreateFail)
		}

		// register host to iam
		// auth: check authorization
		if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDs...); err != nil {
			blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid: %s", hostIDs, err, srvData.rid)
			return srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr, Data: retData})
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
	condition := &meta.DistinctHostIDByTopoRelationRequest{}

	if 0 != data.SetID {
		condition.SetIDArr = []int64{data.SetID}
	}
	if 0 != data.ModuleID {
		condition.ModuleIDArr = []int64{data.ModuleID}
	}

	condition.ApplicationIDArr = []int64{data.ApplicationID}
	hostResult, err := srvData.lgc.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(srvData.ctx, header, condition)
	if err != nil {
		blog.Errorf("get host ids failed, err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if err := hostResult.CCError(); err != nil {
		blog.ErrorJSON("get host id by topology relation failed, error code: %s, error message: %s, cond: %s, rid: %s", hostResult.Code, hostResult.ErrMsg, condition, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	hostIDArr := hostResult.Data.IDArr
	if 0 == len(hostIDArr) {
		blog.Warnf("no host in set,rid:%s", srvData.rid)
		_ = resp.WriteEntity(meta.NewSuccessResp(nil))
		return
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
		blog.Errorf("MoveSetHost2IdleModule GetModuleIDByCond idle module not exist, input:%#v, param:%#v, rid:%s", data, moduleCond, srvData.rid)
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

	var exceptionArr []meta.ExceptionResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		// step3. deregister host from iam
		if err := s.AuthManager.DeregisterHostsByID(srvData.ctx, srvData.header, hostIDArr...); err != nil {
			blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v, rid: %s", hostIDArr, err, srvData.rid)
			return srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)
		}
		hmInput := &meta.HostModuleRelationRequest{
			ApplicationID: data.ApplicationID,
			HostIDArr:     hostIDArr,
			Fields:        []string{common.BKSetIDField, common.BKModuleIDField, common.BKHostIDField},
		}
		configResult, err := srvData.lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(srvData.ctx, srvData.header, hmInput)
		if nil != err {
			blog.Errorf("remove hostModuleConfig, http do error, error:%v, params:%v, input:%+v, rid:%s", err, hmInput, data, srvData.rid)
			return err
		}
		if !configResult.Result {
			blog.Errorf("remove hostModuleConfig http reply error, result:%v, params:%v, input:%+v, rid:%s", configResult, hmInput, data, srvData.rid)
			return err
		}
		hostIDMHMap := make(map[int64][]meta.ModuleHost, 0)
		for _, item := range configResult.Data.Info {
			hostIDMHMap[item.HostID] = append(hostIDMHMap[item.HostID], item)
		}

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

		if err := audit.SaveAudit(srvData.ctx); err != nil {
			blog.Errorf("SaveAudit failed, err: %s, rid: %s", err.Error(), srvData.rid)
		}

		// register host to iam
		// auth: check authorization
		if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDArr...); err != nil {
			blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid:%s", hostIDArr, err, srvData.rid)
			return srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)
		}
		if len(exceptionArr) > 0 {
			blog.Errorf("MoveSetHost2IdleModule has exception. exception:%#v, rid:%s", exceptionArr, srvData.rid)
			return srvData.ccErr.Error(common.CCErrHostDeleteFail)
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr, Data: exceptionArr})
		return
	}
	_ = resp.WriteEntity(meta.NewSuccessResp(nil))
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

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		err = srvData.lgc.CloneHostProperty(srvData.ctx, input.AppID, srcHostID, dstHostID)
		if nil != err {
			blog.Errorf("CloneHostProperty  error , err: %v, input:%#v, rid:%s", err, input, srvData.rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr})
		return
	}
	_ = resp.WriteEntity(meta.NewSuccessResp(nil))
}

// UpdateImportHosts update excel import hosts
func (s *Service) UpdateImportHosts(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	hostList := new(meta.HostList)
	if err := json.NewDecoder(req.Request.Body).Decode(hostList); err != nil {
		blog.Errorf("UpdateImportHosts failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if hostList.HostInfo == nil {
		blog.Errorf("UpdateImportHosts, but host info is nil.input:%+v,rid:%s", hostList, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	hostFields, err := srvData.lgc.GetHostAttributes(srvData.ctx, srvData.ownerID, meta.BizLabelNotExist)
	if err != nil {
		blog.Errorf("UpdateImportHosts, but get host attribute for audit failed, err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}

	hostIDArr := make([]int64, 0)
	hosts := make(map[int64]map[string]interface{}, 0)
	indexHostIDMap := make(map[int64]int64, 0)
	var errMsg, successMsg []string
	for index, hostInfo := range hostList.HostInfo {
		if hostInfo == nil {
			continue
		}
		var intHostID int64
		hostID, ok := hostInfo[common.BKHostIDField]
		if !ok {
			blog.Errorf("UpdateImportHosts failed, because bk_host_id field not exits innerIp: %v, rid: %v", hostInfo[common.BKHostInnerIPField], srvData.rid)
			errMsg = append(errMsg, srvData.ccLang.Languagef("import_update_host_miss_hostID", index))
			continue
		}
		intHostID, err = util.GetInt64ByInterface(hostID)
		if err != nil {
			errMsg = append(errMsg, srvData.ccLang.Languagef("import_update_host_hostID_not_int", index))
			continue
		}
		// bk_host_innerip should not update
		delete(hostInfo, common.BKHostInnerIPField)
		hostIDArr = append(hostIDArr, intHostID)
		hosts[index] = hostInfo
		indexHostIDMap[index] = intHostID
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
			_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		_ = resp.WriteEntity(perm)
		return
	}

	logPreContents := make(map[int64]*logics.HostLog, 0)
	for _, hostID := range hostIDArr {
		audit := srvData.lgc.NewHostLog(srvData.ctx, srvData.ownerID)
		logPreContents[hostID] = audit
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(srvData.ctx, s.EnableTxn, srvData.header, func() error {
		hasHostUpdateWithoutHostApplyFiled := false
		// 功能开关：更新主机属性时是否剔除自动应用字段
		if meta.HostUpdateWithoutHostApplyFiled == true {
			hostProperties, hasRules, err := s.getHostApplyRelatedFields(srvData, hostIDArr)
			if err != nil {
				blog.Errorf("UpdateImportHosts failed, getHostApplyRelatedFields failed, hostIDArr: %+v, err: %v, rid:%s", hostIDArr, err, srvData.rid)
				return err
			}
			// get host attributes
			if hasRules == true {
				hasHostUpdateWithoutHostApplyFiled = true
				for index, hostInfo := range hosts {
					delete(hostInfo, common.BKHostIDField)
					intHostID := indexHostIDMap[index]
					updateData := make(map[string]interface{})
					for key, value := range hostInfo {
						properties, ok := hostProperties[intHostID]
						if ok == true && util.InStrArr(properties, key) {
							continue
						}
						updateData[key] = value
					}
					opt := &meta.UpdateOption{
						Condition: mapstr.MapStr{common.BKHostIDField: intHostID},
						Data:      mapstr.NewFromMap(updateData),
					}
					result, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, opt)
					if err != nil {
						blog.Errorf("UpdateImportHosts UpdateObject http do error, err: %v,input:%+v,param:%+v,rid:%s", err, hostList.HostInfo, opt, srvData.rid)
						errMsg = append(errMsg, srvData.ccLang.Languagef("import_host_update_fail", index, err.Error()))
						continue
					}
					if !result.Result {
						blog.ErrorJSON("UpdateImportHosts failed, UpdateObject failed, param:%s, response: %s, rid:%s", opt, result, srvData.rid)
						errMsg = append(errMsg, srvData.ccLang.Languagef("import_host_update_fail", index, result.ErrMsg))
						continue
					}
					successMsg = append(successMsg, strconv.FormatInt(index, 10))
				}
			}
		}

		if hasHostUpdateWithoutHostApplyFiled == false {
			for index, hostInfo := range hosts {
				delete(hostInfo, common.BKHostIDField)
				intHostID := indexHostIDMap[index]
				opt := &meta.UpdateOption{
					Condition: mapstr.MapStr{common.BKHostIDField: intHostID},
					Data:      mapstr.NewFromMap(hostInfo),
				}
				result, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, opt)
				if err != nil {
					blog.ErrorJSON("UpdateImportHosts UpdateInstance http do error, err: %v,input:%+v,param:%+v,rid:%s", err, hostList.HostInfo, opt, srvData.rid)
					errMsg = append(errMsg, srvData.ccLang.Languagef("import_host_update_fail", index, err.Error()))
					continue
				}
				if !result.Result {
					blog.ErrorJSON("UpdateImportHosts failed, UpdateInstance failed, param:%s, response: %s, rid:%s", opt, result, srvData.rid)
					errMsg = append(errMsg, srvData.ccLang.Languagef("import_host_update_fail", index, result.ErrMsg))
					continue
				}
				successMsg = append(successMsg, strconv.FormatInt(index, 10))
			}
		}

		logLastContents := make([]meta.AuditLog, 0)
		for _, hostID := range hostIDArr {
			audit := logPreContents[hostID]
			if err := audit.WithCurrent(srvData.ctx, hostID, hostFields); err != nil {
				blog.Errorf("UpdateImportHosts, but get host[%d] pre data for audit failed, err: %v, rid: %s", hostID, err, srvData.rid)
				return srvData.ccErr.Error(common.CCErrHostDetailFail)
			}
			hostModuleConfig, err := srvData.lgc.GetConfigByCond(srvData.ctx, meta.HostModuleRelationRequest{HostIDArr: []int64{hostID}, Fields: []string{common.BKAppIDField}})
			if err != nil {
				blog.Errorf("UpdateImportHosts GetConfigByCond failed, id[%v], err: %v,input:%+v,rid:%s", hostID, err, hostList.HostInfo, srvData.rid)
				return err
			}
			var appID int64
			if len(hostModuleConfig) > 0 {
				appID = hostModuleConfig[0].AppID
			}
			auditLog, err := audit.AuditLog(srvData.ctx, hostID, appID, meta.AuditUpdate)
			if err != nil {
				blog.Errorf("UpdateImportHosts create audit log failed, id[%v], err: %v,rid:%s", hostID, err, srvData.rid)
				return err
			}
			logLastContents = append(logLastContents, auditLog)
		}

		auditResp, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, logLastContents...)
		if err != nil {
			blog.Errorf("UpdateImportHosts, but add host[%v] audit failed, err: %v, rid:%s", hostIDArr, err, srvData.rid)
			return err
		}
		if !auditResp.Result {
			blog.Errorf("UpdateImportHosts, but add host[%v] audit failed, err: %v, rid:%s", hostIDArr, auditResp.ErrMsg, srvData.rid)
			return srvData.ccErr.New(auditResp.Code, auditResp.ErrMsg)
		}
		return nil
	})

	if txnErr != nil {
		_ = resp.WriteError(http.StatusOK, &meta.RespError{Msg: txnErr})
		return
	}
	retData := map[string]interface{}{
		"error":   errMsg,
		"success": successMsg,
	}
	_ = resp.WriteEntity(meta.NewSuccessResp(retData))
}
