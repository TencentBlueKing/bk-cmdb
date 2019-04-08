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
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"

	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
)

func (ps *ProcServer) CreateProcess(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	ownerID := srvData.ownerID
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		blog.Errorf("convert appid from string to int failed!, err: %s, input:%s,rid:%s", err.Error(), appID, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	input := mapstr.New()
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("create process failed! decode request body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	procName, err := input.String(common.BKProcNameField)
	if err != nil {
		blog.Errorf("create process failed! get process name error, err: %v,input:%#v,rid:%s", err, input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedString, common.BKAppNameField)})
		return
	}

	input[common.BKAppIDField] = appID
	input[common.BKOwnerIDField] = ownerID
	ret, err := ps.CoreAPI.CoreService().Instance().CreateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDProc, &meta.CreateModelInstance{Data: input})
	if err != nil {
		blog.Errorf("CreateProcess http do error. err:%s,input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !ret.Result {
		blog.Errorf("CreateProcess http reply error.err code:%d err msg:%s,input:%+v,rid:%s", ret.Code, ret.ErrMsg, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	// save change log
	instID := ret.Data.Created.ID
	curDetail, err := ps.getProcDetail(req, ownerID, appID, int(instID))
	if err != nil {
		blog.Errorf("get process instance detail failed. err:%s,input:%+v,rid:%s", err.Error(), input, srvData.rid)
	} else {
		ps.addProcLog(srvData.ctx, srvData.ownerID, appIDStr, srvData.user, nil, curDetail, auditoplog.AuditOpTypeAdd, int(instID), srvData.header)
	}

	// auth: register process to iam
	processSimplify := extensions.ProcessSimplify{
		ProcessID:    int64(appID),
		ProcessName:  procName,
		BKAppIDField: int64(instID),
	}
	if err := ps.AuthManager.RegisterProcesses(srvData.ctx, srvData.header, processSimplify); err != nil {
		blog.Warnf("create process sucess, but register to iam failed, err: %+v, process: %+v, rid: %s", err, processSimplify, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}

	// return success
	result := make(map[string]interface{})
	result[common.BKProcessIDField] = instID
	resp.WriteEntity(meta.NewSuccessResp(result))
}

func (ps *ProcServer) UpdateProcess(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	ownerID := req.PathParameter(common.BKOwnerIDField)
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		blog.Errorf("convert appid from string to int failed!, err: %s,appID:%v,rid:%s", err.Error(), appIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	procIDStr := req.PathParameter(common.BKProcessIDField)
	procID, err := strconv.Atoi(procIDStr)
	if err != nil {
		blog.Errorf("convert procid from string to int failed!, err: %s,procID:%+v,rid:%s", err.Error(), procIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	procData := mapstr.New()
	if err := json.NewDecoder(req.Request.Body).Decode(&procData); err != nil {
		blog.Errorf("create process failed! decode request body err: %v,appID:%+v,procID:%+v,rid:%s", err, appID, procID, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	procData[common.BKAppIDField] = appID

	// take snapshot before operation
	preProcDetail, err := ps.getProcDetail(req, ownerID, appID, procID)
	if err != nil {
		blog.Errorf("get process instance detail failed. err:%s", err.Error())
	}
	// change proc name set value
	procName := ""
	if procData.Exists(common.BKProcNameField) {
		procName, err = procData.String(common.BKProcNameField)
		if err != nil {
			blog.Errorf("create process failed! get process name error, err: %v,input:%#v, url para:%#v,rid:%s", err, procData, req.PathParameters(), srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedString, common.BKAppNameField)})
			return
		}
	}

	input := new(meta.UpdateOption)
	condition := make(map[string]interface{})
	condition[common.BKOwnerIDField] = ownerID
	condition[common.BKAppIDField] = appID
	condition[common.BKProcessIDField] = procID
	input.Condition = condition
	input.Data = procData
	ret, err := ps.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDProc, input)
	if err != nil {
		blog.Errorf("UpdateProcess http do error. err:%s,input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !ret.Result {
		blog.Errorf("UpdateProcess http reply error.err code:%d,err msg:%s,input:%+v,rid:%s", ret.Code, ret.ErrMsg, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	// save operation log
	// take snapshot before operation
	curDetail, err := ps.getProcDetail(req, ownerID, appID, procID)
	if err != nil {
		blog.Errorf("get process instance detail failed. err:%s,rid:%s", err.Error(), srvData.rid)
	}
	ps.addProcLog(srvData.ctx, ownerID, appIDStr, srvData.user, preProcDetail, curDetail, auditoplog.AuditOpTypeModify, procID, srvData.header)
	if procData.Exists(common.BKProcNameField) {
		// auth: register process to iam
		processSimplify := extensions.ProcessSimplify{
			ProcessID:    int64(appID),
			ProcessName:  procName,
			BKAppIDField: int64(procID),
		}
		if err := ps.AuthManager.UpdateRegisteredProcesses(srvData.ctx, srvData.header, processSimplify); err != nil {
			blog.Warnf("update process success, but update register to iam failed, err: %+v, process: %+v, rid: %s", err, processSimplify, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) BatchUpdateProcess(req *restful.Request, resp *restful.Response) {
	// user := util.GetUser(req.Request.Header)
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	ownerID := req.PathParameter(common.BKOwnerIDField)
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		blog.Errorf("convert appid from string to int failed!, err: %s,appID:%s,rid:%s", err.Error(), appIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	procData := mapstr.New()
	if err := json.NewDecoder(req.Request.Body).Decode(&procData); err != nil {
		blog.Errorf("create process failed! decode request body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	// split process str by comma
	var procIDArr []string
	if procIDStr, ok := procData[common.BKProcessIDField].(string); !ok {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrProcFieldValidFaile)})
		return
	} else {
		procIDArr = strings.Split(procIDStr, ",")
	}

	// update condition
	procData[common.BKAppIDField] = appID
	delete(procData, common.BKProcessIDField)

	// parse process id and valid
	var iProcIDArr []int
	auditContentArr := make([]meta.Content, len(procIDArr))

	// change process name set value
	procName := ""
	if procData.Exists(common.BKProcNameField) {
		procName, err = procData.String(common.BKProcNameField)
		if err != nil {
			blog.Errorf("create process failed! get process name error, err: %v,input:%#v, url para:%#v,rid:%s", err, procData, req.PathParameters(), srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedString, common.BKAppNameField)})
			return
		}
	}

	updatedProcesses := make([]extensions.ProcessSimplify, 0)
	for index, procIDStr := range procIDArr {
		procID, err := strconv.Atoi(procIDStr)
		if err != nil {
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
			return
		}
		details, err := ps.getProcDetail(req, ownerID, appID, procID)
		if err != nil {
			blog.Errorf("get inst detail error: %v, input:%+v,rid:%s", err, procData, srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrAuditSaveLogFaile)})
			return
		}
		// save change log
		headers := []meta.Header{}
		curData := map[string]interface{}{}
		for _, detail := range details {
			curData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
			headers = append(headers,
				meta.Header{
					PropertyID:   fmt.Sprint(detail[common.BKPropertyIDField].(string)),
					PropertyName: fmt.Sprint(detail[common.BKPropertyNameField]),
				})
		}

		curData[common.BKProcIDField] = procID

		// save proc info before modify
		auditContentArr[index] = meta.Content{
			CurData: make(map[string]interface{}),
			PreData: curData,
			Headers: headers,
		}
		iProcIDArr = append(iProcIDArr, procID)
		// update processes
		input := new(meta.UpdateOption)
		condition := make(map[string]interface{})
		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKAppIDField] = appID
		condition[common.BKProcessIDField] = procID

		input.Condition = condition
		input.Data = procData
		ret, err := ps.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDProc, input)
		if err != nil {
			blog.Errorf("BatchUpdateProcess http do error.err:%s,input:%+v,rid:%s", err.Error(), input, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
			return
		}
		if !ret.Result {
			blog.Errorf("BatchUpdateProcess http reply error.err code:%d,err msg:%s,input:%+v,rid:%s", ret.Code, ret.ErrMsg, input, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
			return
		}
		if procData.Exists(common.BKProcNameField) {
			processSimplify := extensions.ProcessSimplify{
				ProcessID:    int64(procID),
				ProcessName:  procName,
				BKAppIDField: int64(appID),
			}
			updatedProcesses = append(updatedProcesses, processSimplify)
		}
	} // end for procIDArr

	if len(updatedProcesses) > 0 {
		if err := ps.AuthManager.UpdateRegisteredProcesses(srvData.ctx, srvData.header, updatedProcesses...); err != nil {
			blog.Errorf("batch update processes success, but update register to iam failed, err: %+v, processes:%+v, rid:%s", err, updatedProcesses, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
			return
		}
	}

	logscontent := make([]auditoplog.AuditLogContext, 0)
	// update audit log content after modify
	for index, procID := range iProcIDArr {

		// take snapshot before operation
		details, err := ps.getProcDetail(req, ownerID, appID, int(procID))
		if err != nil {
			blog.Errorf("get inst detail error: %v, procID:%s,rid:%s", err, procID, srvData.rid)
			continue
		}

		// save change log
		curData := map[string]interface{}{}
		for _, detail := range details {
			curData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
		}

		curData["bk_process_id"] = procID

		// save proc info before modify
		auditContentArr[index].CurData = curData
		logscontent = append(logscontent, auditoplog.AuditLogContext{ID: int64(procID), Content: auditContentArr[index]})
	}
	if 0 < len(logscontent) {
		logs := meta.AuditProcsParams{
			Content: logscontent,
			OpDesc:  "update process",
			OpType:  auditoplog.AuditOpTypeModify,
		}

		ps.CoreAPI.AuditController().AddProcLogs(srvData.ctx, ownerID, appIDStr, srvData.user, srvData.header, logs)
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) DeleteProcess(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	ownerID := req.PathParameter(common.BKOwnerIDField)
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		blog.Errorf("convert appid from string to int failed!, err: %s,appID:%v,rid:%s", err.Error(), appIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	procIDStr := req.PathParameter(common.BKProcessIDField)
	procID, err := strconv.Atoi(procIDStr)
	if err != nil {
		blog.Errorf("convert procid from string to int failed!, err: %s,procID:%v,rid:%s", err.Error(), procIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	condition := mapstr.MapStr{common.BKProcessIDField: procID}
	// get process by module,restrict the process deletion of the associated module
	p2mRet, err := ps.CoreAPI.ProcController().GetProc2Module(srvData.ctx, srvData.header, condition)
	if err != nil {
		blog.Errorf("DeleteProcess GetProc2Module http do error. err:%s,input:(appID:%v,procID:%v),query:%+v,rid:%s", err.Error(), appID, procID, condition, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !p2mRet.Result {
		blog.Errorf("DeleteProcess GetProc2Module http reply error.err code:%d err msg:%s, input:(appID:%v,procID:%v),query:%+v,rid:%s", p2mRet.Code, p2mRet.ErrMsg, appID, procID, condition, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(p2mRet.Code, p2mRet.ErrMsg)})
		return
	}
	// take snapshot before operation
	preProcDetail, err := ps.getProcDetail(req, ownerID, appID, procID)
	if err != nil {
		blog.Errorf("get process instance detail failed. err:%s", err.Error())
	}

	conditon := make(map[string]interface{})
	conditon[common.BKAppIDField] = appID
	conditon[common.BKProcessIDField] = procID
	conditon[common.BKOwnerIDField] = ownerID
	ret, err := ps.CoreAPI.CoreService().Instance().DeleteInstance(srvData.ctx, srvData.header, common.BKInnerObjIDProc, &meta.DeleteOption{Condition: conditon})
	if err != nil {
		blog.Errorf("DeleteProcess DelObject http do error.err:%s,input:%+v,rid:%s", err.Error(), conditon, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !ret.Result {
		blog.Errorf("DeleteProcess DelObject http reply error.err code:%d,err msg:%s,input:%+v,rid:%s", ret.Code, ret.ErrMsg, conditon, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	// save operation log
	ps.addProcLog(srvData.ctx, srvData.ownerID, appIDStr, srvData.user, preProcDetail, nil, auditoplog.AuditOpTypeDel, procID, srvData.header)

	// auth: dereigster process from iam
	processSimplify := extensions.ProcessSimplify{
		ProcessID:    int64(procID),
		BKAppIDField: int64(appID),
	}
	if err := ps.AuthManager.DeregisterProcesses(srvData.ctx, srvData.header, processSimplify); err != nil {
		blog.Errorf("delete processes success, but deregister from iam failed, err: %+v, processes:%+v, rid:%s", err, processSimplify, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) SearchProcess(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	ownerID := req.PathParameter(common.BKOwnerIDField)
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		blog.Errorf("convert appid from string to int failed!, err: %s,appID:%v,rid:%s", err.Error(), appIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	var srchparam params.SearchParams
	srchparam.Condition = mapstr.New()
	if err := json.NewDecoder(req.Request.Body).Decode(&srchparam); err != nil {
		blog.Errorf("decode request body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	condition := srchparam.Condition
	condition[common.BKOwnerIDField] = ownerID
	condition[common.BKAppIDField] = appID
	if processName, ok := condition[common.BKProcessNameField]; ok {
		processNameStr, ok := processName.(string)
		if ok {
			condition[common.BKProcessNameField] = map[string]interface{}{common.BKDBLIKE: params.SpecialCharChange(processNameStr)}

		} else {
			condition[common.BKProcessNameField] = map[string]interface{}{common.BKDBLIKE: processName}

		}
	}
	page := srchparam.Page
	searchParams := new(meta.QueryCondition)
	searchParams.Condition = condition
	searchParams.Fields = srchparam.Fields
	searchParams.Limit.Offset, err = util.GetInt64ByInterface(page["start"])
	if nil != err {
		searchParams.Limit.Offset = 0
	}
	searchParams.Limit.Limit, err = util.GetInt64ByInterface(page["limit"])
	if nil != err {
		searchParams.Limit.Limit = common.BKNoLimit
	}
	if sort, ok := page["sort"].(string); !ok {
		searchParams.SortArr = nil
	} else {
		searchParams.SortArr = meta.NewSearchSortParse().String(sort).ToSearchSortArr()
	}

	// query process by module name
	if moduleName, ok := condition[common.BKModuleNameField]; ok {
		reqParam := make(map[string]interface{})
		reqParam[common.BKAppIDField] = appID
		reqParam[common.BKModuleNameField] = moduleName

		blog.V(5).Infof("get process arr by module(%s), reqParam: %+v,rid:%s", moduleName, reqParam, srvData.rid)
		ret, err := ps.CoreAPI.ProcController().GetProc2Module(srvData.ctx, srvData.header, reqParam)
		if err != nil {
			blog.Errorf("SearchProcess GetProc2Module http do error.err:%s,input:%+v,rid:%s", err.Error(), reqParam, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
			return
		}
		if !ret.Result {
			blog.Errorf("SearchProcess GetProc2Module http reply error.err code:%d,err msg:%s,input:%+v,rid:%s", ret.Code, ret.ErrMsg, reqParam, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
			return
		}

		// parse procid array
		procIdArr := make([]int64, 0)
		for _, item := range ret.Data {
			procIdArr = append(procIdArr, item.ProcessID)
		}

		// update condition
		condition[common.BKProcessIDField] = map[string]interface{}{
			"$in": procIdArr,
		}
		delete(condition, common.BKModuleNameField)
		searchParams.Condition = condition
	}

	// search process
	ret, err := ps.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDProc, searchParams)
	if err != nil {
		blog.Errorf("SearchProcess SearchObjects http do error.err:%s,input:%+v,rid:%s", err.Error(), searchParams, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !ret.Result {
		blog.Errorf("SearchProcess SearchObjects http reply error.err code:%d,err msg:%s,input:%+v,rid:%s", ret.Code, ret.ErrMsg, searchParams, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(ret.Data))
}

func (ps *ProcServer) addProcLog(ctx context.Context, ownerID, appID, user string, preProcDetails, curProcDetails []map[string]interface{}, auditType auditoplog.AuditOpType, instanceID int, header http.Header) error {

	headers := []meta.Header{}
	curData := map[string]interface{}{}
	preData := map[string]interface{}{}
	for _, detail := range preProcDetails {
		preData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
		if nil == curProcDetails {
			headers = append(headers,
				meta.Header{
					PropertyID:   fmt.Sprint(detail[common.BKPropertyIDField]),
					PropertyName: fmt.Sprint(detail[common.BKPropertyNameField]),
				})
		}
	}

	for _, detail := range curProcDetails {
		curData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
		headers = append(headers,
			meta.Header{
				PropertyID:   fmt.Sprint(detail[common.BKPropertyIDField]),
				PropertyName: fmt.Sprint(detail[common.BKPropertyNameField]),
			})
	}

	auditContent := meta.Content{
		CurData: curData,
		PreData: preData,
		Headers: headers,
	}
	desc := fmt.Sprintf("unknown type:%s", auditType)
	switch auditType {
	case auditoplog.AuditOpTypeAdd:
		desc = "create process"
	case auditoplog.AuditOpTypeDel:
		desc = "delete process"
	case auditoplog.AuditOpTypeModify:
		desc = "update process"
	}

	log := common.KvMap{common.BKContentField: auditContent, common.BKOpDescField: desc, common.BKOpTypeField: auditType, "inst_id": instanceID}
	_, err := ps.CoreAPI.AuditController().AddProcLog(ctx, ownerID, appID, user, header, log)
	return err
}
