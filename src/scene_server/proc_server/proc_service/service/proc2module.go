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

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (ps *ProcServer) BindModuleProcess(req *restful.Request, resp *restful.Response) {
	user := util.GetUser(req.Request.Header)
	ownerID := util.GetOwnerID(req.Request.Header)
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	pathParams := req.PathParameters()
	appIDStr := pathParams[common.BKAppIDField]
	appID, _ := strconv.Atoi(appIDStr)
	procIDStr := pathParams[common.BKProcIDField]
	procID, _ := strconv.Atoi(procIDStr)
	moduleName := pathParams[common.BKModuleNameField]
	params := make([]interface{}, 0)
	cell := make(map[string]interface{})
	cell[common.BKAppIDField] = appID
	cell[common.BKProcIDField] = procID
	cell[common.BKModuleNameField] = moduleName
	cell[common.BKOwnerIDField] = util.GetOwnerID(req.Request.Header)
	params = append(params, cell)

	// TODO use change use chan, process model trigger point
	// if err := ps.createProcInstanceModel(appIDStr, procIDStr, moduleName, ownerID, &sourceAPI.ForwardParam{Header:req.Request.Header}); err != nil {
	//     blog.Errorf("fail to create process instance model. err: %v", err)
	//     resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg:defErr.Error(common.CCErrProcBindToMoudleFaile)})
	//     return
	// }

	ret, err := ps.CoreAPI.ProcController().CreateProc2Module(context.Background(), req.Request.Header, params)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("fail to BindModuleProcess. err: %v, errcode:%d, errmsg: %s", err.Error(), ret.Code, ret.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcBindToMoudleFaile)})
		return
	}

	// save operation log
	log := common.KvMap{common.BKOpDescField: fmt.Sprintf("bind module [%s]", moduleName), common.BKOpTypeField: auditoplog.AuditOpTypeAdd, "inst_id": procID, common.BKContentField: meta.Content{}}
	ps.CoreAPI.AuditController().AddProcLog(context.Background(), ownerID, appIDStr, user, req.Request.Header, log)

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) DeleteModuleProcessBind(req *restful.Request, resp *restful.Response) {
	user := util.GetUser(req.Request.Header)
	ownerID := util.GetOwnerID(req.Request.Header)
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	pathParams := req.PathParameters()
	appIDStr := pathParams[common.BKAppIDField]
	appID, _ := strconv.Atoi(appIDStr)
	procIDStr := pathParams[common.BKProcIDField]
	procID, _ := strconv.Atoi(procIDStr)
	moduleName := pathParams[common.BKModuleNameField]
	cell := make(map[string]interface{})
	cell[common.BKAppIDField] = appID
	cell[common.BKProcIDField] = procID
	cell[common.BKModuleNameField] = moduleName

	if err := ps.deleteProcInstanceModel(appIDStr, procIDStr, moduleName, req.Request.Header); err != nil {
		blog.Errorf("%v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcUnBindToMoudleFaile)})
		return
	}

	ret, err := ps.CoreAPI.ProcController().DeleteProc2Module(context.Background(), req.Request.Header, cell)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("fail to delete module process bind. err: %v, errcode:%s, errmsg: %s", err, ret.Code, ret.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcUnBindToMoudleFaile)})
		return
	}

	// save operation log
	log := common.KvMap{common.BKOpDescField: fmt.Sprintf("unbind module [%s]", moduleName), common.BKOpTypeField: auditoplog.AuditOpTypeAdd, "inst_id": procID, common.BKContentField: meta.Content{}}
	ps.CoreAPI.AuditController().AddProcLog(context.Background(), ownerID, appIDStr, user, req.Request.Header, log)

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) GetProcessBindModule(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	pathParams := req.PathParameters()
	appIDStr := pathParams[common.BKAppIDField]
	appID, _ := strconv.Atoi(appIDStr)
	procIDStr := pathParams[common.BKProcIDField]
	procID, _ := strconv.Atoi(procIDStr)

	// search object instance
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appID
	input := new(meta.QueryInput)
	input.Condition = condition

	objRet, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, req.Request.Header, input)
	if err != nil || (err == nil && !objRet.Result) {
		blog.Errorf("fail to GetProcessBindModule when do searchobject. err:%v, errcode:%d, errmsg:%s", err, objRet.Code, objRet.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrObjectSelectInstFailed)})
		return
	}

	condition[common.BKProcIDField] = procID
	// get process by module
	p2mRet, err := ps.CoreAPI.ProcController().GetProc2Module(context.Background(), req.Request.Header, condition)
	if err != nil || (err == nil && !p2mRet.Result) {
		blog.Errorf("fail to GetProcessBindModule when do GetProc2Module. err:%v, errcode:%d, errmsg:%s", err, p2mRet.Code, p2mRet.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcSelectBindToMoudleFaile)})
		return
	}

	disModuleNameArr := make([]string, 0)
	for _, moduleInfo := range objRet.Data.Info {
		if !util.InArray(moduleInfo[common.BKModuleNameField], disModuleNameArr) {
			moduleName, ok := moduleInfo[common.BKModuleNameField].(string)
			if !ok {
				continue
			}

			isDefault64, ok := moduleInfo[common.BKDefaultField].(float64)
			if !ok {
				isDefault, ok := moduleInfo[common.BKDefaultField].(int)
				if !ok || 0 != isDefault {
					continue
				}
			} else {
				if 0 != isDefault64 {
					continue
				}
			}
			disModuleNameArr = append(disModuleNameArr, moduleName)
		}
	}

	result := make([]interface{}, 0)
	for _, disModuleName := range disModuleNameArr {
		num := 0
		isBind := 0
		for _, moduleInfo := range objRet.Data.Info {
			moduleName, ok := moduleInfo[common.BKModuleNameField].(string)
			if !ok {
				continue
			}

			if disModuleName == moduleName {
				num++
			}
		}

		for _, procModule := range p2mRet.Data {
			if disModuleName == procModule.ModuleName {
				isBind = 1
				break
			}
		}

		data := make(map[string]interface{})
		data[common.BKModuleNameField] = disModuleName
		data["set_num"] = num
		data["is_bind"] = isBind
		result = append(result, data)
	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}
