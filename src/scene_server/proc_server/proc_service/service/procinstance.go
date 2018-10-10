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

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (ps *ProcServer) OperateProcessInstance(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	procOpParam := new(meta.ProcessOperate)
	if err := json.NewDecoder(req.Request.Body).Decode(procOpParam); err != nil {
		blog.Errorf("fail to decode process operation parameter. err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	header := req.Request.Header
	matchProcInstParam := new(meta.MatchProcInstParam)
	matchProcInstParam.ApplicationID = procOpParam.ApplicationID
	matchProcInstParam.FuncID = procOpParam.FuncID
	matchProcInstParam.HostInstanceID = procOpParam.HostInstanceID
	matchProcInstParam.ModuleName = procOpParam.ModuleName
	matchProcInstParam.SetName = procOpParam.SetName
	procInstModel, err := ps.Logics.MatchProcessInstance(context.Background(), matchProcInstParam, header)
	if err != nil {
		blog.Errorf("match process instance failed in OperateProcessInstance. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcOperateFaile)})
		return
	}

	result, err := ps.Logics.OperateProcInstanceByGse(context.Background(), procOpParam, procInstModel, header)
	if err != nil {
		blog.Errorf("operate process failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcOperateFaile)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}

func (ps *ProcServer) QueryProcessOperateResult(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	taskID := req.PathParameter("taskID")
	succ, waitExec, mapExceErr, err := ps.Logics.QueryProcessOperateResult(context.Background(), taskID, req.Request.Header)
	if nil != err {
		data := common.KvMap{"error": mapExceErr}
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcQueryTaskInfoFail), Data: data})
		return
	}
	if 0 != len(mapExceErr) {
		data := common.KvMap{"error": mapExceErr}
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcQueryTaskOPErrFail), Data: data})
		return
	}
	if 0 != len(waitExec) {
		data := common.KvMap{"wait": waitExec}
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcQueryTaskWaitOPFail), Data: data})
		return
	}
	resp.WriteEntity(meta.NewSuccessResp(succ))
}

func (ps *ProcServer) RefreshProcHostInstByEvent(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	input := new(meta.EventInst)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("fail to decode RefreshProcHostInstByEvent request body. err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	ps.Logics.HandleHostProcDataChange(context.Background(), input)
	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) deleteProcInstanceModel(appId, procId, moduleName string, forward http.Header) error {
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId
	condition[common.BKProcessIDField] = procId
	condition[common.BKModuleNameField] = moduleName

	ret, err := ps.CoreAPI.ProcController().DeleteProcInstanceModel(context.Background(), forward, condition)
	if err != nil || (err == nil && !ret.Result) {
		return fmt.Errorf("fail to delete process instance model. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
	}

	return nil
}
