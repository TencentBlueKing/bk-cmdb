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
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"configcenter/src/ac"
	"configcenter/src/ac/extensions"
	authmeta "configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

// AppResult TODO
type AppResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    DataInfo    `json:"data"`
}

// DataInfo TODO
type DataInfo struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

// DeleteHostBatchFromResourcePool delete hosts from resource pool
func (s *Service) DeleteHostBatchFromResourcePool(ctx *rest.Contexts) {

	opt := new(meta.DeleteHostBatchOpt)
	if err := ctx.DecodeInto(&opt); nil != err {
		ctx.RespAutoError(err)
		return
	}

	hostIDArr := strings.Split(opt.HostID, ",")
	var iHostIDArr []int64
	delCondsArr := make([][]map[string]interface{}, 0)
	for _, i := range hostIDArr {
		iHostID, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			blog.Errorf("delete host batch, but got invalid host id, err: %v,input:%+v,rid:%s", err, opt, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, iHostID))
			return
		}
		iHostIDArr = append(iHostIDArr, iHostID)
	}
	iHostIDArr = util.IntArrayUnique(iHostIDArr)

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Delete,
		iHostIDArr...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", iHostIDArr, err, ctx.Kit.Rid)
		if err != ac.NoAuthorizeError {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostDeleteFail))
			return
		}
		perm, err := s.AuthManager.GenHostBatchNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Delete,
			iHostIDArr)
		if err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostDeleteFail))
			return
		}
		ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
		return
	}

	for _, iHostID := range iHostIDArr {
		asstCond := map[string]interface{}{
			common.BKDBOR: []map[string]interface{}{
				{
					common.BKObjIDField:  common.BKInnerObjIDHost,
					common.BKInstIDField: iHostID,
				},
				{
					common.BKAsstObjIDField:  common.BKInnerObjIDHost,
					common.BKAsstInstIDField: iHostID,
				},
			},
		}

		queryCond := &meta.InstAsstQueryCondition{
			Cond:  meta.QueryCondition{Condition: asstCond},
			ObjID: common.BKInnerObjIDHost,
		}

		rsp, err := s.CoreAPI.CoreService().Association().ReadInstAssociation(ctx.Kit.Ctx, ctx.Kit.Header, queryCond)
		if nil != err {
			blog.ErrorJSON("DeleteHostBatch read host association do request failed , err: %s, rid: %s", err.Error(),
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
			return
		}

		if rsp.Count <= 0 {
			continue
		}
		asstInstMap := make(map[string][]int64, 0)
		for _, asst := range rsp.Info {
			if asst.ObjectID == common.BKInnerObjIDHost && iHostID == asst.InstID {
				asstInstMap[asst.AsstObjectID] = append(asstInstMap[asst.AsstObjectID], asst.AsstInstID)
			} else if asst.AsstObjectID == common.BKInnerObjIDHost && iHostID == asst.AsstInstID {
				asstInstMap[asst.ObjectID] = append(asstInstMap[asst.ObjectID], asst.InstID)
			} else {
				ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommDBSelectFailed,
					"host is not associated in selected association"))
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
			instRsp, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, objID,
				&meta.QueryCondition{Condition: instCond})
			if err != nil {
				blog.ErrorJSON("DeleteHostBatch read associated instances do request failed , err: %s, rid: %s",
					err.Error(), ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
				return
			}

			if len(instRsp.Info) > 0 {
				blog.ErrorJSON("DeleteHostBatch host %s has been associated, can't be deleted, rid: %s", iHostID,
					ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrTopoInstHasBeenAssociation, iHostID))
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

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		for _, delConds := range delCondsArr {
			opt := &meta.InstAsstDeleteOption{
				Opt:   meta.DeleteOption{Condition: map[string]interface{}{common.BKDBOR: delConds}},
				ObjID: common.BKInnerObjIDHost,
			}
			_, err := s.CoreAPI.CoreService().Association().DeleteInstAssociation(ctx.Kit.Ctx, ctx.Kit.Header, opt)
			if err != nil {
				blog.ErrorJSON("DeleteHostBatch delete host redundant association do request failed , err: %s, "+
					"rid: %s", err.Error(), ctx.Kit.Rid)
				return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
			}
		}
		appID, err := s.Logic.GetDefaultAppID(ctx.Kit)
		if err != nil {
			blog.Errorf("delete host batch, but got invalid app id, err: %v,input:%s,rid:%s", err, opt, ctx.Kit.Rid)
			return ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
		}

		// for audit log.
		audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, meta.AuditDelete)

		// to generate audit log about deleting host, and host information.
		auditCond := map[string]interface{}{common.BKHostIDField: map[string]interface{}{common.BKDBIN: iHostIDArr}}
		logContents, err := audit.GenerateAuditLogByCond(generateAuditParameter, appID, auditCond)
		if err != nil {
			blog.Errorf("generate host audit log failed before delete host, hostIDs: %+v, bizID: %d, err: %v, "+
				"rid: %s", iHostIDArr, appID, err, ctx.Kit.Rid)
			return err
		}

		hosts := make([]extensions.HostSimplify, len(logContents))
		for index, logContent := range logContents {
			hosts[index] = extensions.HostSimplify{
				BKAppIDField:       0,
				BKHostIDField:      logContent.ID,
				BKHostInnerIPField: logContent.ResourceName,
			}
		}

		input := &meta.DeleteHostRequest{
			ApplicationID: appID,
			HostIDArr:     iHostIDArr,
		}
		delResult, err := s.CoreAPI.CoreService().Host().DeleteHostFromSystem(ctx.Kit.Ctx, ctx.Kit.Header, input)
		if err != nil {
			blog.Error("DeleteHostBatch DeleteHost http do error. err:%s, input:%s, rid:%s", err.Error(), input,
				ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !delResult.Result {
			blog.Errorf("DeleteHostBatch DeleteHost http reply error. result: %#v, input:%#v, rid:%s", delResult,
				input, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrHostDeleteFail)
		}

		// to save audit.
		if len(logContents) > 0 {
			if err := audit.SaveAuditLog(ctx.Kit, logContents...); err != nil {
				blog.Errorf("save host audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return ctx.Kit.CCError.CCError(common.CCErrAuditSaveLogFailed)
			}
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// GetHostInstanceProperties get host instance's properties as follows:
// host object property id: "bk_host_name"
// host object property name: "host"
// host object property value: "centos7"
func (s *Service) GetHostInstanceProperties(ctx *rest.Contexts) {

	hostID := ctx.Request.PathParameter("bk_host_id")
	hostIDInt64, err := strconv.ParseInt(hostID, 10, 64)
	if err != nil {
		blog.Errorf("convert hostID to int64, err: %v,host:%s,rid:%s", err, hostID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKHostIDField))
		return
	}
	details, _, err := s.Logic.GetHostInstanceDetails(ctx.Kit, hostIDInt64)
	if err != nil {
		blog.Errorf("get host details failed, err: %v,host:%s,rid:%s", err, hostID, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(details) == 0 {
		blog.Errorf("host not found, hostID: %v,rid:%s", hostID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostNotFound))
		return
	}

	// check authorization
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Find, hostIDInt64); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDInt64, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
		return
	}
	attribute, err := s.Logic.GetHostAttributes(ctx.Kit, nil)
	if err != nil {
		blog.Errorf("get host attribute fields failed, err: %v,rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := make([]meta.HostInstanceProperties, 0)
	for _, attr := range attribute {
		result = append(result, meta.HostInstanceProperties{
			PropertyID:    attr.PropertyID,
			PropertyName:  attr.PropertyName,
			PropertyValue: details[attr.PropertyID],
		})
	}

	ctx.RespEntity(result)

}

// AddHost TODO
// add host to host resource pool
func (s *Service) AddHost(ctx *rest.Contexts) {
	hostList := new(meta.HostList)
	if err := ctx.DecodeInto(&hostList); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if hostList.HostInfo == nil {
		blog.Errorf("add host, but host info is nil.input:%+v,rid:%s", hostList, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsNeedSet))
		return
	}

	appID := hostList.ApplicationID
	if appID == 0 {
		// get default app id
		var err error
		appID, err = s.Logic.GetDefaultAppIDWithSupplier(ctx.Kit)
		if err != nil {
			blog.Errorf("add host, but get default app id failed, err: %v,input:%+v,rid:%s", err, hostList, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}

	// get target biz's idle module ID
	cond := hutil.NewOperation().WithAppID(appID).MapStr()
	cond.Set(common.BKDefaultField, common.DefaultResModuleFlag)
	moduleID, _, err := s.Logic.GetResourcePoolModuleID(ctx.Kit, cond)
	if err != nil {
		blog.Errorf("add host, but get module id failed, err: %s,input: %+v,rid: %s", err.Error(), hostList, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	retData := make(map[string]interface{})
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		_, success, updateErrRow, errRow, err := s.Logic.AddHost(ctx.Kit, appID, []int64{moduleID},
			ctx.Kit.SupplierAccount, hostList.HostInfo, hostList.InputType)
		if err != nil {
			blog.Errorf("add host failed, success: %v, update: %v, err: %v, %v,input:%+v,rid:%s",
				success, updateErrRow, err, errRow, hostList, ctx.Kit.Rid)
			retData["error"] = errRow
			retData["update_error"] = updateErrRow
			return ctx.Kit.CCError.CCError(common.CCErrHostCreateFail)
		}
		retData["success"] = success
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(retData, txnErr)
		return
	}
	ctx.RespEntity(retData)
}

// AddHostByExcel TODO
// add host come from excel to host resource pool
func (s *Service) AddHostByExcel(ctx *rest.Contexts) {
	hostList := new(meta.HostList)
	if err := ctx.DecodeInto(&hostList); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if hostList.HostInfo == nil {
		blog.Errorf("add host, but host info is nil.input:%+v,rid:%s", hostList, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsNeedSet))
		return
	}

	appID := hostList.ApplicationID
	if appID == 0 {
		// get default app id
		var err error
		appID, err = s.Logic.GetDefaultAppIDWithSupplier(ctx.Kit)
		if err != nil {
			blog.Errorf("add host, but get default app id failed, err: %v,input:%+v,rid:%s", err, hostList, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}

	moduleID := hostList.ModuleID
	if moduleID == 0 {
		// get target biz's idle module ID
		cond := hutil.NewOperation().WithAppID(appID).MapStr()
		cond.Set(common.BKDefaultField, common.DefaultResModuleFlag)
		var err error
		moduleID, _, err = s.Logic.GetResourcePoolModuleID(ctx.Kit, cond)
		if err != nil {
			blog.Errorf("add host, but get module id failed, err: %s,input: %+v,rid: %s", err.Error(), hostList, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}

	retData := make(map[string]interface{})
	_, success, errRow, err := s.Logic.AddHostByExcel(ctx.Kit, appID, moduleID, ctx.Kit.SupplierAccount, hostList.HostInfo)
	retData["success"] = success
	retData["error"] = errRow
	if err != nil {
		blog.Errorf("add host failed, success: %v, errRow:%v, err: %v, hostList:%#v, rid:%s",
			success, errRow, err, hostList, ctx.Kit.Rid)
		ctx.RespEntityWithError(retData, ctx.Kit.CCError.CCError(common.CCErrHostCreateFail))
	}

	ctx.RespEntity(retData)
}

// AddHostToResourcePool TODO
// add host to resource pool, returns bk_host_id of the successfully added hosts
func (s *Service) AddHostToResourcePool(ctx *rest.Contexts) {

	hostList := new(meta.AddHostToResourcePoolHostList)
	body, err := ioutil.ReadAll(ctx.Request.Request.Body)
	if err != nil {
		blog.Errorf("read request body failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPReadBodyFailed))
		return
	}
	if err := json.Unmarshal(body, hostList); err != nil {
		blog.Errorf("add host failed with decode body err: %v, body: %s, rid:%s", err, string(body), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	if hostList.HostInfo == nil {
		blog.ErrorJSON("add host, but host info is nil. input:%s, rid:%s", hostList, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsNeedSet))
		return
	}
	_, retData, err := s.Logic.AddHostToResourcePool(ctx.Kit, *hostList)

	if err != nil {
		blog.ErrorJSON("add host failed, retData: %s, err: %s, input:%s, rid:%s", retData, err, hostList, ctx.Kit.Rid)
		ctx.RespEntityWithError(retData, err)
		return
	}
	ctx.RespEntity(retData)
}

// AddHostFromAgent TODO
// Deprecated:
func (s *Service) AddHostFromAgent(ctx *rest.Contexts) {

	agents := new(meta.AddHostFromAgentHostList)
	if err := ctx.DecodeInto(&agents); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if len(agents.HostInfo) == 0 {
		blog.Errorf("add host from agent, but got 0 agents from body.input:%+v,rid:%s", agents, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "HostInfo"))
		return
	}
	appID, err := s.Logic.GetDefaultAppID(ctx.Kit)
	if err != nil {
		blog.Errorf("AddHostFromAgent GetDefaultAppID error.input:%#v,rid:%s", agents, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if 0 == appID {
		blog.Errorf("add host from agent, but got invalid default appID, err: %v,ownerID:%s,input:%#v,rid:%s", err, ctx.Kit.SupplierAccount, agents, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrAddHostToModule, "business not found"))
		return
	}

	opt := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithAppID(appID)
	moduleID, _, err := s.Logic.GetResourcePoolModuleID(ctx.Kit, opt.MapStr())
	if err != nil {
		blog.Errorf("add host from agent , but get module id failed, err: %v,ownerID:%s,input:%+v,rid:%s", err, ctx.Kit.SupplierAccount, agents, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	agents.HostInfo["import_from"] = common.HostAddMethodAgent
	addHost := make(map[int64]map[string]interface{})
	addHost[1] = agents.HostInfo
	var success, updateErrRow, errRow []string
	retData := make(map[string]interface{})
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		_, success, updateErrRow, errRow, err = s.Logic.AddHost(ctx.Kit, appID, []int64{moduleID},
			common.BKDefaultOwnerID, addHost, "")
		if err != nil {
			blog.Errorf("add host failed, success: %v, update: %v, err: %v, %v,input:%+v,rid:%s",
				success, updateErrRow, err, errRow, agents, ctx.Kit.Rid)

			retData["success"] = success
			retData["error"] = errRow
			retData["update_error"] = updateErrRow
			return ctx.Kit.CCError.CCError(common.CCErrHostCreateFail)
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(retData, txnErr)
		return
	}
	ctx.RespEntity(success)
}

// SearchHost TODO
func (s *Service) SearchHost(ctx *rest.Contexts) {

	body := new(meta.HostCommonSearch)
	if err := ctx.DecodeInto(&body); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	host, err := s.Logic.SearchHost(ctx.Kit, body, false)
	if err != nil {
		blog.Errorf("search host failed, err: %v,input:%+v,rid:%s", err, body, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostGetFail))
		return
	}

	hostIDArray := host.ExtractHostIDs()
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Find, *hostIDArray...); err != nil {
		blog.Errorf("check host authorization failed, hostID: %+v, err: %+v, rid: %s", hostIDArray, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
		return
	}

	ctx.RespEntity(host)

}

// SearchHostWithAsstDetail TODO
func (s *Service) SearchHostWithAsstDetail(ctx *rest.Contexts) {

	body := new(meta.HostCommonSearch)
	if err := ctx.DecodeInto(&body); nil != err {
		ctx.RespAutoError(err)
		return
	}

	host, err := s.Logic.SearchHost(ctx.Kit, body, true)
	if err != nil {
		blog.Errorf("search host failed, err: %v,input:%+v,rid:%s", err, body, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(host)
}

// UpdateHostBatch update many hosts once
func (s *Service) UpdateHostBatch(ctx *rest.Contexts) {

	data := mapstr.New()
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	// TODO: this is a wrong usage, just for compatible the wrong usage before.
	// delete this, when the frontend use the rigListHostInstanceht request field. not the number.
	id := data[common.BKHostIDField]
	hostIDStr := ""
	switch id.(type) {
	case float64:
		floatID := id.(float64)
		hostIDStr = strconv.FormatInt(int64(floatID), 10)
	case string:
		hostIDStr = id.(string)
	default:
		blog.Errorf("update host batch failed, got invalid host id(%v) data type,rid:%s", id, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "bk_host_id"))
		return
	}

	data.Remove(common.MetadataField)
	data.Remove(common.BKHostIDField)
	data.Remove(common.BKCloudIDField)

	// check authorization
	hostIDArr := make([]int64, 0)
	for _, id := range strings.Split(hostIDStr, ",") {
		hostID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			blog.Errorf("update host batch, but got invalid host id[%s], err: %v,rid:%s", id, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsInvalid))
			return
		}
		hostIDArr = append(hostIDArr, hostID)
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Update,
		hostIDArr...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDArr, err, ctx.Kit.Rid)
		if err != nil && err != ac.NoAuthorizeError {
			blog.ErrorJSON("check host authorization failed, hosts: %s, err: %s, rid: %s", hostIDArr, err.Error(),
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
			return
		}
		perm, err := s.AuthManager.GenHostBatchNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Update, hostIDArr)
		if err != nil && err != ac.NoAuthorizeError {
			blog.ErrorJSON("check host authorization get permission failed, hosts: %s, err: %s, rid: %s", hostIDArr,
				err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
			return
		}
		ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
		return
	}

	// for audit log.
	audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// generator audit log.
		genAuditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, meta.AuditUpdate).WithUpdateFields(data)
		auditCond := map[string]interface{}{common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDArr}}
		auditLogs, err := audit.GenerateAuditLogByCond(genAuditParam, 0, auditCond)
		if err != nil {
			blog.Errorf("generate host audit log failed, hostIDs: %+v, err: %v, rid: %s", hostIDArr, err, ctx.Kit.Rid)
			return err
		}

		// to update host.
		opt := &meta.UpdateOption{
			Condition: mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: hostIDArr}},
			Data:      mapstr.NewFromMap(data),
		}
		_, err = s.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKInnerObjIDHost, opt)
		if err != nil {
			blog.Errorf("UpdateHostBatch UpdateObject http do error, err: %v, input: %+v, param: %+v, rid: %s",
				err, data, opt, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}

		// save audit log.
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save host audit log failed after update host, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (s *Service) getHostBizMapAndHostInfoMap(kit *rest.Kit, hostIDs []int64) (map[int64]int64,
	map[int64]mapstr.MapStr, error) {

	hostCond := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}
	hosts, err := s.Logic.GetHostInfoByConds(kit, hostCond)
	if err != nil {
		blog.Errorf("get hosts failed, condition: %#v, err: %v, rid: %s", hostCond, err, kit.Rid)
		return nil, nil, err
	}

	hostMap := make(map[int64]mapstr.MapStr)
	for _, host := range hosts {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.Errorf("parse host id failed, host: %#v, err: %v, rid: %s", host, err, kit.Rid)
			return nil, nil, err
		}
		hostMap[hostID] = host
	}

	input := meta.HostModuleRelationRequest{
		HostIDArr: hostIDs,
		Fields:    []string{common.BKAppIDField, common.BKHostIDField},
	}
	hostRelations, rawErr := s.Logic.GetHostRelations(kit, input)
	if rawErr != nil {
		blog.Errorf("get host relations failed, hostIDs: %+v, err: %v, rid: %s", hostIDs, err, kit.Rid)
		return nil, nil, err
	}

	hostBizMap := make(map[int64]int64)
	for _, relation := range hostRelations {
		hostBizMap[relation.HostID] = relation.AppID
	}

	return hostBizMap, hostMap, nil
}

// updateHostPropertyBatch concurrent update of host's property fields.
func (s *Service) updateHostPropertyBatch(kit *rest.Kit, hostIDArr []int64,
	parameter *meta.UpdateHostPropertyBatchParameter) error {

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
		hostBizMap, hostMap, err := s.getHostBizMapAndHostInfoMap(kit, hostIDArr)
		if err != nil {
			return err
		}

		auditContexts := make([]meta.AuditLog, 0)
		audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())

		var (
			wg       sync.WaitGroup
			firstErr error
		)
		pipeline := make(chan bool, 5)

		for _, update := range parameter.Update {
			if firstErr != nil {
				break
			}

			pipeline <- true
			wg.Add(1)

			genAuditParam := auditlog.NewGenerateAuditCommonParameter(kit, meta.AuditUpdate)

			go func(update meta.UpdateHostProperty) {
				defer func() {
					wg.Done()
					<-pipeline
				}()

				cond := mapstr.New()
				cond.Set(common.BKHostIDField, update.HostID)
				data, err := mapstr.NewFromInterface(update.Properties)
				if err != nil {
					blog.Errorf("convert properties: %v to mapStr failed, err: %v, rid: %s", update.Properties,
						err, kit.Rid)
					firstErr = err
					return
				}

				// can't update host's cloud area using this api
				data.Remove(common.BKCloudIDField)
				data.Remove(common.BKHostIDField)

				// generate audit log.
				genAuditParam.WithUpdateFields(data)
				hostInfo := []mapstr.MapStr{hostMap[update.HostID]}

				auditLog, err := audit.GenerateAuditLog(genAuditParam, hostBizMap[update.HostID], hostInfo)
				if err != nil {
					blog.Errorf("write audit log failed, hostID: %d, err: %v, rid: %s", update.HostID, err, kit.Rid)
					firstErr = err
					return
				}

				opt := &meta.UpdateOption{
					Condition: cond,
					Data:      data,
				}
				// to update data.
				_, err = s.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header,
					common.BKInnerObjIDHost, opt)
				if err != nil {
					blog.Errorf("update host property failed, input: %+v, opt: %+v, err: %v, rid: %s", data, opt, err,
						kit.Rid)
					firstErr = err
					return
				}
				// add audit log.
				auditContexts = append(auditContexts, auditLog...)
			}(update)
		}

		wg.Wait()
		if firstErr != nil {
			return firstErr
		}

		// save audit log.
		if err := audit.SaveAuditLog(kit, auditContexts...); err != nil {
			blog.Errorf("add hosts %+v audit failed, err: %v, rid: %s", hostIDArr, err, kit.Rid)
			return err
		}
		return nil
	})

	return txnErr
}

// UpdateHostPropertyBatch batch update host properties.
func (s *Service) UpdateHostPropertyBatch(ctx *rest.Contexts) {

	parameter := new(meta.UpdateHostPropertyBatchParameter)
	if err := ctx.DecodeInto(&parameter); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if len(parameter.Update) > common.BKMaxPageSize {
		blog.Errorf("update host property batch failed, data len %d exceed max pageSize %d, rid: %s",
			len(parameter.Update), common.BKMaxPageSize, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "update", common.BKMaxPageSize))
		return
	}

	hostIDs := make([]int64, 0)
	for _, update := range parameter.Update {
		hostIDs = append(hostIDs, update.HostID)
	}

	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Update, hostIDs...); err != nil {
		if err != nil && err != ac.NoAuthorizeError {
			blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDs, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
			return
		}
		perm, err := s.AuthManager.GenHostBatchNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Update, hostIDs)
		if err != nil && err != ac.NoAuthorizeError {
			blog.Errorf("check host authorization get permission failed, hosts: %+v, err: %v, rid: %s", hostIDs,
				err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
			return
		}
		blog.Errorf("hosts no authorized, hosts: %+v, rid: %s", hostIDs, ctx.Kit.Rid)
		ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
		return
	}

	if err := s.updateHostPropertyBatch(ctx.Kit, hostIDs, parameter); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// NewHostSyncAppTopo add new hosts to the business
// synchronize hosts directly to a module in a business if this host does not exist.
// otherwise, this operation will only change host's attribute.
// TODO: used by framework.
func (s *Service) NewHostSyncAppTopo(ctx *rest.Contexts) {

	hostList := new(meta.HostSyncList)
	if err := ctx.DecodeInto(&hostList); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if hostList.HostInfo == nil {
		blog.Errorf("add host, but host info is nil.input:%+v,rid:%s", hostList, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "host_info"))
		return
	}
	if 0 == len(hostList.ModuleID) {
		blog.Errorf("host sync app  parameters required moduleID,input:%+v,rid:%s", hostList, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKModuleIDField))
		return
	}

	if common.BatchHostAddMaxRow < len(hostList.HostInfo) {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "host_info ", common.BatchHostAddMaxRow))
		return
	}

	appConds := map[string]interface{}{
		common.BKAppIDField: hostList.ApplicationID,
	}

	appInfo, err := s.Logic.GetAppDetails(ctx.Kit, "", appConds)
	if nil != err {
		blog.Errorf("host sync app %d error:%s,input:%+v,rid:%s", hostList.ApplicationID, err.Error(), hostList, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if 0 == len(appInfo) {
		blog.Errorf("host sync app %d not found, reply:%+v,input:%+v,rid:%s", hostList.ApplicationID, appInfo, hostList, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoGetAppFailed))
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
	moduleIDS, err := s.Logic.GetModuleIDByCond(ctx.Kit, meta.ConditionWithTime{Condition: moduleCond})
	if nil != err {
		blog.Errorf("NewHostSyncAppTop GetModuleIDByCond error. err:%s,input:%+v,rid:%s", err.Error(), hostList, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(moduleIDS) != len(hostList.ModuleID) {
		blog.Errorf("not found part module: source:%v, db:%v, rid: %s", hostList.ModuleID, moduleIDS, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrHostModuleIDNotFoundORHasMultipleInnerModuleIDFailed))
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeCreateHost(ctx.Kit.Ctx, ctx.Kit.Header, hostList.ApplicationID); err != nil {
		blog.Errorf("check add hosts authorization failed, business: %d, err: %v, rid: %s", hostList.ApplicationID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
		return
	}

	retData := make(map[string]interface{})
	var success, updateErrRow, errRow []string
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		_, success, updateErrRow, errRow, err = s.Logic.AddHost(ctx.Kit, hostList.ApplicationID,
			hostList.ModuleID, ctx.Kit.SupplierAccount, hostList.HostInfo, common.InputTypeApiNewHostSync)
		if err != nil {
			blog.Errorf("add host failed, success: %v, update: %v, err: %v, %v, rid: %s",
				success, updateErrRow, err, errRow, ctx.Kit.Rid)

			retData["success"] = success
			retData["error"] = errRow
			retData["update_error"] = updateErrRow
			return ctx.Kit.CCError.CCError(common.CCErrHostCreateFail)
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(retData, txnErr)
		return
	}
	ctx.RespEntity(success)
}

// MoveSetHost2IdleModule bk_set_id and bk_module_id cannot be empty at the same time
// Remove the host from the module or set.
// The host belongs to the current module or host only, and puts the host into the idle machine of the current service.
// When the host data is in multiple modules or sets. Disconnect the host from the module or set only
// TODO: used by v2 version, remove this api when v2 is offline.
func (s *Service) MoveSetHost2IdleModule(ctx *rest.Contexts) {
	header := ctx.Kit.Header

	var data meta.SetHostConfigParams
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if 0 == data.ApplicationID {
		blog.Errorf("MoveSetHost2IdleModule bk_biz_id cannot be empty at the same time,input:%#v,rid:%s", data,
			util.GetHTTPCCRequestID(header))
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsNeedSet))
		return
	}

	if 0 == data.SetID && 0 == data.ModuleID {
		blog.Errorf("MoveSetHost2IdleModule bk_set_id and bk_module_id cannot be empty at the same time,input:%#v, "+
			"rid:%s", data, util.GetHTTPCCRequestID(header))
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsNeedSet))
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
	hostIDArr, ccErr := s.Logic.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(ctx.Kit.Ctx, header, condition)
	if ccErr != nil {
		blog.Errorf("get host ids failed, err: %v, rid: %s", ccErr, ctx.Kit.Rid)
		ctx.RespAutoError(ccErr)
		return
	}

	if 0 == len(hostIDArr) {
		blog.Warnf("no host in set,rid:%s", ctx.Kit.Rid)
		ctx.RespEntity(nil)
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

	moduleIDArr, err := s.Logic.GetModuleIDByCond(ctx.Kit, meta.ConditionWithTime{Condition: moduleCond})
	if err != nil {
		blog.Errorf("MoveSetHost2IdleModule GetModuleIDByCond error. err:%s, input:%#v, param:%#v, rid:%s",
			err.Error(), data, moduleCond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(moduleIDArr) == 0 {
		blog.Errorf("MoveSetHost2IdleModule GetModuleIDByCond idle module not exist, input:%#v, param:%#v, rid:%s",
			data, moduleCond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrHostModuleNotExist, "idle module"))
		return
	}
	idleModuleID := moduleIDArr[0]
	moduleHostConfigParams := make(map[string]interface{})
	moduleHostConfigParams[common.BKAppIDField] = data.ApplicationID
	audit := auditlog.NewHostModuleLog(s.CoreAPI.CoreService(), hostIDArr)

	var exceptionArr []meta.ExceptionResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {

		hmInput := &meta.HostModuleRelationRequest{
			ApplicationID: data.ApplicationID,
			HostIDArr:     hostIDArr,
			Fields:        []string{common.BKSetIDField, common.BKModuleIDField, common.BKHostIDField},
		}
		configResult, err := s.Logic.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header,
			hmInput)
		if nil != err {
			blog.Errorf("remove hostModuleConfig, http do error, error:%v, params:%v, input:%+v, rid:%s", err,
				hmInput, data, ctx.Kit.Rid)
			return err
		}

		hostIDMHMap := make(map[int64][]meta.ModuleHost, 0)
		for _, item := range configResult.Info {
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

			var opResult []meta.ExceptionResult
			var ccErr errors.CCErrorCoder
			if toEmptyModule {
				input := &meta.TransferHostToInnerModule{
					ApplicationID: data.ApplicationID,
					ModuleID:      idleModuleID,
					HostID:        []int64{hostID},
				}
				opResult, ccErr = s.Logic.CoreAPI.CoreService().Host().TransferToInnerModule(ctx.Kit.Ctx,
					ctx.Kit.Header, input)
			} else {
				input := &meta.HostsModuleRelation{
					ApplicationID: data.ApplicationID,
					HostID:        []int64{hostID},
					ModuleID:      newModuleIDArr,
				}
				opResult, ccErr = s.Logic.CoreAPI.CoreService().Host().TransferToNormalModule(ctx.Kit.Ctx,
					ctx.Kit.Header, input)
			}

			if ccErr != nil {
				blog.Errorf("transfer host failed, err: %v, result: %#v, to idle module:%v, input: %#v, rid: %s",
					err, opResult, toEmptyModule, data, ctx.Kit.Rid)
				if len(opResult) > 0 {
					exceptionArr = append(exceptionArr, opResult...)
				} else {
					exceptionArr = append(exceptionArr, meta.ExceptionResult{
						Code:        int64(ccErr.GetCode()),
						Message:     ccErr.Error(),
						OriginIndex: hostID,
					})
				}
			}
		}

		if err := audit.SaveAudit(ctx.Kit); err != nil {
			blog.Errorf("SaveAudit failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrHostDeleteFail)
		}

		if len(exceptionArr) > 0 {
			blog.Errorf("MoveSetHost2IdleModule has exception. exception:%#v, rid:%s", exceptionArr, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrHostDeleteFail)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(exceptionArr, txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (s *Service) ip2hostID(kit *rest.Kit, input *meta.CloneHostPropertyParams) (src int64, dst int64, err error) {

	cond := meta.QueryCondition{
		Fields: []string{common.BKHostIDField, common.BKHostInnerIPField},
		Page: meta.BasePage{
			Limit: common.BKNoLimit,
		},
	}

	useIP := false
	if len(input.OrgIP) != 0 {
		// use host inner ip
		cond.Condition = map[string]interface{}{
			common.BKHostInnerIPField: map[string]interface{}{common.BKDBIN: []string{input.OrgIP, input.DstIP}},
			common.BKCloudIDField:     input.CloudID,
		}
		useIP = true
	} else {
		// use host id
		cond.Condition = map[string]interface{}{
			common.BKHostIDField:  map[string]interface{}{common.BKDBIN: []int64{input.OrgID, input.DstID}},
			common.BKCloudIDField: input.CloudID,
		}
	}

	hosts, err := s.Logic.SearchHostInfo(kit, cond)
	if err != nil {
		blog.ErrorJSON("search hosts failed, err: %s, input: %s, rid: %s", err, cond, kit.Rid)
		return 0, 0, err
	}

	if !useIP {
		if len(hosts) != 2 {
			return 0, 0, errors.New(common.CCErrCommParamsInvalid, "src or dst id is not exists")
		}

		return input.OrgID, input.DstID, nil
	}

	// use ip
	orgID, dstID := int64(0), int64(0)
	for _, host := range hosts {
		hostID, err := host.Int64(common.BKHostIDField)
		if err != nil {
			blog.ErrorJSON("parse host id failed, err: %s, host: %s, rid: %s", err, host, kit.Rid)
			return 0, 0, err
		}

		hostIP, err := host.String(common.BKHostInnerIPField)
		if err != nil {
			blog.ErrorJSON("parse host ip failed, err: %s, host: %s, rid: %s", err, host, kit.Rid)
			return 0, 0, err
		}

		ipArr := strings.Split(hostIP, ",")
		for _, slicedIP := range ipArr {
			if slicedIP == input.OrgIP {
				orgID = hostID
			}

			if slicedIP == input.DstIP {
				dstID = hostID
			}
		}
	}

	if orgID == 0 || dstID == 0 {
		return 0, 0, errors.New(common.CCErrCommParamsInvalid, "invalid org or dst data")
	}

	return orgID, dstID, nil
}

// CloneHostProperty clone host property from src host to dst host
// can only clone editable fields that are not in host model unique rules.
// origin ip and dest ip can only be one ip.
func (s *Service) CloneHostProperty(ctx *rest.Contexts) {

	input := new(meta.CloneHostPropertyParams)
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if input.AppID <= 0 {
		blog.Errorf("invalid bk_biz_id: %d ,rid: %s", input.AppID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "bk_biz_id"))
		return
	}

	if input.CloudID < 0 {
		blog.Errorf("invalid bk_cloud_id: %d ,rid: %s", input.CloudID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "bk_cloud_id"))
		return
	}

	// can only use ip or id for one.
	if (len(input.OrgIP) != 0 || len(input.DstIP) != 0) && (input.OrgID > 0 || input.DstID > 0) {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid org/dst ip or id")
		return
	}

	if (len(input.OrgIP) == 0 && len(input.DstIP) == 0) && (input.OrgID <= 0 && input.DstID <= 0) {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid org/dst ip or id")
		return
	}

	if (len(input.OrgIP) != 0 || len(input.DstIP) != 0) && (len(input.OrgIP) == 0 || len(input.DstIP) == 0) {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "no parameter")
		return
	}

	if input.OrgID < 0 || input.DstID < 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid org/dst id")
		return
	}

	if (input.OrgID > 0 || input.DstID > 0) && (input.OrgID <= 0 || input.DstID <= 0) {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid org/dst id")
		return
	}

	if (len(input.OrgIP) != 0 && len(input.DstIP) != 0) && (input.OrgIP == input.DstIP) {
		ctx.RespEntity(nil)
		return
	}

	if (input.OrgID > 0 && input.DstID > 0) && (input.OrgID == input.DstID) {
		ctx.RespEntity(nil)
		return
	}

	// authorization check
	orgID, dstID, err := s.ip2hostID(ctx.Kit, input)
	if err != nil {
		blog.ErrorJSON("get host id from ip failed, input: %s, err: %s, rid:%s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// if both src ip and dst ip belongs to the same host, do not need to clone
	if orgID == dstID {
		ctx.RespEntity(nil)
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Find, orgID); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid:%s", orgID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
		return
	}

	// step2. verify has permission to update dst host
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Update, dstID); err != nil {
		if err != ac.NoAuthorizeError {
			blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid:%s", dstID, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
			return
		}
		perm, err := s.AuthManager.GenEditBizHostNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, []int64{dstID})
		if err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
			return
		}
		ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logic.CloneHostProperty(ctx.Kit, input.AppID, orgID, dstID)
		if nil != err {
			blog.Errorf("CloneHostProperty  error , err: %v, input:%#v, rid:%s", err, input, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// UpdateImportHosts update excel import hosts
func (s *Service) UpdateImportHosts(ctx *rest.Contexts) {
	hostList := new(meta.HostList)
	if err := ctx.DecodeInto(&hostList); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if hostList.HostInfo == nil {
		blog.Errorf("UpdateImportHosts, but host info is nil.input:%+v,rid:%s", hostList, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsNeedSet))
		return
	}

	hostIDArr := make([]int64, 0)
	hosts := make(map[int64]map[string]interface{}, 0)
	indexHostIDMap := make(map[int64]int64, 0)
	var errMsg, successMsg []string
	CCLang := s.Language.CreateDefaultCCLanguageIf(util.GetLanguage(ctx.Kit.Header))
	for _, index := range util.SortedMapInt64Keys(hostList.HostInfo) {
		hostInfo := hostList.HostInfo[index]
		if hostInfo == nil {
			continue
		}
		var intHostID int64
		hostID, ok := hostInfo[common.BKHostIDField]
		if !ok {
			blog.Errorf("UpdateImportHosts failed, because bk_host_id field not exits innerIp: %v, rid: %v",
				hostInfo[common.BKHostInnerIPField], ctx.Kit.Rid)

			errMsg = append(errMsg, CCLang.Languagef("import_update_host_miss_hostID", index))
			continue
		}
		intHostID, err := util.GetInt64ByInterface(hostID)
		if err != nil {
			errMsg = append(errMsg, CCLang.Languagef("import_update_host_hostID_not_int", index))
			continue
		}

		// remove unchangeable fields
		delete(hostInfo, common.BKHostInnerIPField)
		delete(hostInfo, common.BKCloudIDField)
		delete(hostInfo, common.BKImportFrom)
		delete(hostInfo, common.CreateTimeField)

		hostIDArr = append(hostIDArr, intHostID)
		hosts[index] = hostInfo
		indexHostIDMap[index] = intHostID
	}

	if len(hostIDArr) == 0 {
		ctx.RespEntity(map[string]interface{}{
			"error":   errMsg,
			"success": []string{},
		})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Update,
		hostIDArr...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDArr, err, ctx.Kit.Rid)
		if err != nil && err != ac.NoAuthorizeError {
			blog.ErrorJSON("check host authorization failed, hosts: %s, err: %s, rid: %s", hostIDArr, err.Error(),
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
			return
		}
		perm, err := s.AuthManager.GenHostBatchNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Update, hostIDArr)
		if err != nil && err != ac.NoAuthorizeError {
			blog.ErrorJSON("check host authorization get permission failed, hosts: %s, err: %s, rid: %s", hostIDArr,
				err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
			return
		}
		ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
		return
	}

	// audit interface of host audit log.
	audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())
	auditContexts := make([]meta.AuditLog, 0)

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		hostCond := map[string]interface{}{common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDArr}}
		hostInfoArr, err := s.Logic.GetHostInfoByConds(ctx.Kit, hostCond)
		if err != nil {
			blog.Errorf("get hosts failed, err: %v, condition: %#v, rid: %s", err, hostCond, ctx.Kit.Rid)
			return err
		}

		hostMap := make(map[int64]mapstr.MapStr)
		for _, host := range hostInfoArr {
			hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
			if err != nil {
				blog.Errorf("parse host id failed, err: %v, host: %#v, rid: %s", err, host, ctx.Kit.Rid)
				return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
			}
			hostMap[hostID] = host
		}

		hostRelations, rawErr := s.Logic.GetHostRelations(ctx.Kit, meta.HostModuleRelationRequest{HostIDArr: hostIDArr,
			Fields: []string{common.BKAppIDField, common.BKHostIDField}})
		if rawErr != nil {
			blog.Errorf("get host relations failed, err: %v, hostIDs: %+v, rid: %s", err, hostIDArr, ctx.Kit.Rid)
			return rawErr
		}

		hostBizMap := make(map[int64]int64)
		for _, relation := range hostRelations {
			hostBizMap[relation.HostID] = relation.AppID
		}

		ccLang := s.Language.CreateDefaultCCLanguageIf(util.GetLanguage(ctx.Kit.Header))
		for _, index := range util.SortedMapInt64Keys(hosts) {
			host := hosts[index]
			delete(host, common.BKHostIDField)
			intHostID := indexHostIDMap[index]

			// generate audit log.
			genAuditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, meta.AuditUpdate).WithUpdateFields(host)
			auditLog, err := audit.GenerateAuditLog(genAuditParam, hostBizMap[intHostID],
				[]mapstr.MapStr{hostMap[intHostID]})
			if err != nil {
				blog.Errorf("generate host audit log failed, hostID: %d, err: %v, rid: %s", intHostID, err, ctx.Kit.Rid)
				errMsg = append(errMsg, err.Error())
				continue
			}

			// to update data.
			opt := &meta.UpdateOption{
				Condition: mapstr.MapStr{common.BKHostIDField: intHostID},
				Data:      mapstr.NewFromMap(host),
			}
			_, err = s.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header,
				common.BKInnerObjIDHost, opt)
			if err != nil {
				blog.ErrorJSON("UpdateImportHosts UpdateInstance http do error, err: %v,input:%+v,param:%+v,rid:%s",
					err, hostList.HostInfo, opt, ctx.Kit.Rid)
				errMsg = append(errMsg, ccLang.Languagef("import_host_update_fail", index, err.Error()))
				continue
			}

			successMsg = append(successMsg, strconv.FormatInt(index, 10))
			auditContexts = append(auditContexts, auditLog...)
		}

		// save audit log.
		if err := audit.SaveAuditLog(ctx.Kit, auditContexts...); err != nil {
			blog.Errorf("success update host, but add host[%v] audit failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	retData := map[string]interface{}{
		"error":   errMsg,
		"success": successMsg,
	}
	ctx.RespEntity(retData)
}

// CountHostCPU 查询业务下的主机CPU数量的特殊接口，给成本管理使用
func (s *Service) CountHostCPU(ctx *rest.Contexts) {
	req := new(meta.CountHostCPUReq)
	if err := ctx.DecodeInto(&req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// get specified biz host cpu count
	if req.BizID != 0 {
		cnt, err := s.countBizHostCPU(ctx.Kit, req.BizID)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntity([]meta.BizHostCpuCount{cnt})
		return
	}

	// get paged biz ids(including resource pool & not archived biz) sort by id, then get host cpu count of each biz
	bizReq := &meta.QueryCondition{
		Condition: mapstr.MapStr{common.BKDataStatusField: mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled}},
		Page:      meta.BasePage{Start: req.Page.Start, Limit: req.Page.Limit, Sort: common.BKAppIDField},
		Fields:    []string{common.BKAppIDField},
	}

	bizRes, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDApp,
		bizReq)
	if err != nil {
		blog.Errorf("get biz ids failed, input: %+v, err: %v, rid: %s", bizReq, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := make([]meta.BizHostCpuCount, len(bizRes.Info))
	for idx, biz := range bizRes.Info {
		bizID, err := biz.Int64(common.BKAppIDField)
		if err != nil {
			blog.Errorf("parse biz id failed, biz: %+v, err: %v, rid: %s", biz, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
			return
		}

		cnt, err := s.countBizHostCPU(ctx.Kit, bizID)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		result[idx] = cnt
	}

	ctx.RespEntity(result)
}

// countBizHostCPU count host cpu num in one biz
func (s *Service) countBizHostCPU(kit *rest.Kit, bizID int64) (meta.BizHostCpuCount, error) {
	cnt := meta.BizHostCpuCount{BizID: bizID}

	pageSize := 500
	relReq := &meta.HostModuleRelationRequest{
		ApplicationID: bizID,
		Page:          meta.BasePage{Start: 0, Limit: pageSize, Sort: common.BKHostIDField},
		Fields:        []string{common.BKHostIDField},
	}

	// get host ids in biz and count cpu num, process 1000 relations at a time
	var prevHostID int64
	for ; ; relReq.Page.Start += pageSize {
		relRes, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, relReq)
		if err != nil {
			blog.Errorf("get biz host ids failed, req: %+v, err: %v, rid: %s", relReq, err, kit.Rid)
			return meta.BizHostCpuCount{}, err
		}

		relLen := len(relRes.Info)
		if relLen == 0 {
			break
		}

		hostIDs := make([]int64, 0)
		for _, rel := range relRes.Info {
			// since relations is sorted by host id, we use the previous host id to distinct the ids
			if rel.HostID == prevHostID {
				continue
			}
			hostIDs = append(hostIDs, rel.HostID)
			prevHostID = rel.HostID
		}

		if len(hostIDs) == 0 {
			continue
		}

		// get host cpu num, count total host num and cpu num and host with no cpu field num
		hostReq := &meta.QueryInput{
			Condition:      mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: hostIDs}},
			Fields:         "bk_cpu",
			DisableCounter: true,
		}

		hostRes, err := s.CoreAPI.CoreService().Host().GetHosts(kit.Ctx, kit.Header, hostReq)
		if err != nil {
			blog.Errorf("get hosts failed, req: %+v, err: %v, rid: %s", hostReq, err, kit.Rid)
			return meta.BizHostCpuCount{}, err
		}

		for _, host := range hostRes.Info {
			cnt.HostCount++
			cpuCnt, exists := host["bk_cpu"]
			if !exists || cpuCnt == nil {
				cnt.NoCpuHostCount++
				continue
			}

			cpuCount, err := util.GetInt64ByInterface(cpuCnt)
			if err != nil {
				blog.Errorf("parse host cpu count(%+v) failed, err: %v, rid: %s", cpuCnt, err, kit.Rid)
				return meta.BizHostCpuCount{}, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_cpu")
			}

			if cpuCnt == 0 {
				cnt.NoCpuHostCount++
				continue
			}

			cnt.CpuCount += cpuCount
		}

	}

	return cnt, nil
}
