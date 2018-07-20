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

package actions

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/api/auditlog"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	//	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}", Params: nil, Handler: process.GetProcessBindModule})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{bk_module_name}", Params: nil, Handler: process.BindModuleProcess})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{bk_module_name}", Params: nil, Handler: process.DeleteModuleProcessBind})

	process.CreateAction()
}

//BindModuleProcess bind proce 2 module
func (cli *procAction) BindModuleProcess(req *restful.Request, resp *restful.Response) {
	user := util.GetActionUser(req)
	ownerID := util.GetActionOnwerID(req)
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

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
		params = append(params, cell)
		pJson, _ := json.Marshal(params)

		bindURL := cli.CC.ProcCtrl() + "/process/v1/module"
		blog.Info("bind proc module config url: %v", bindURL)
		blog.Info("bind proc module config params: %v", string(pJson))
		bindData, err := httpcli.ReqHttp(req, bindURL, common.HTTPCreate, []byte(pJson))
		blog.Info("bind proc module config return: %v", bindData)
		if nil != err {
			blog.Error("BindModuleProcess   error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcBindToMoudleFaile)
		}
		var bindre ProcessResult
		err = json.Unmarshal([]byte(bindData), &bindre)
		if nil != err || false == bindre.Result {
			blog.Error("BindModuleProcess   error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcBindToMoudleFaile)
		}

		auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header).AuditProcLog(procID, "", fmt.Sprintf("bind module [%s]", moduleName), ownerID, appIDStr, user, auditoplog.AuditOpTypeModify)

		return http.StatusOK, nil, nil
	}, resp)
}

//DeleteModuleProcessBind delete process module bind
func (cli *procAction) DeleteModuleProcessBind(req *restful.Request, resp *restful.Response) {
	user := util.GetActionUser(req)
	ownerID := util.GetActionOnwerID(req)
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

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
		pJson, _ := json.Marshal(cell)

		dURL := cli.CC.ProcCtrl() + "/process/v1/module"
		blog.Info("delete proc module config bind url: %v", dURL)
		blog.Info("delete proc module config bind params: %v", string(pJson))
		bindData, err := httpcli.ReqHttp(req, dURL, common.HTTPDelete, []byte(pJson))
		blog.Info("delete bind proc module config bind return: %v", bindData)
		if nil != err {
			blog.Error("delete module process bind  error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcUnBindToMoudleFaile)
		}
		var bindre ProcessResult
		err = json.Unmarshal([]byte(bindData), &bindre)
		if nil != err || false == bindre.Result {
			blog.Error("delete module process bind  error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcUnBindToMoudleFaile)
		}
		auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header).AuditProcLog(procID, "", fmt.Sprintf("unbind module [%s]", moduleName), ownerID, appIDStr, user, auditoplog.AuditOpTypeModify)
		return http.StatusOK, nil, nil
	}, resp)
}

//GetProcessBindModule get process bind module
func (cli *procAction) GetProcessBindModule(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		appIDStr := pathParams[common.BKAppIDField]
		appID, _ := strconv.Atoi(appIDStr)
		procIDStr := pathParams[common.BKProcIDField]
		procID, _ := strconv.Atoi(procIDStr)
		condition := make(map[string]interface{})
		condition[common.BKAppIDField] = appID
		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		sCondJson, _ := json.Marshal(searchParams)

		gModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/module/search"
		blog.Info("get module query url: %v", gModuleURL)
		blog.Info("get module query params: %v", string(sCondJson))
		gModuleRe, err := httpcli.ReqHttp(req, gModuleURL, common.HTTPSelectPost, []byte(sCondJson))
		blog.Info("get module query return: %v", gModuleRe)
		if nil != err {
			blog.Error("GetProcessBindModule Module  error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleSelectFailed)
		}
		var modules ModuleSResult
		err = json.Unmarshal([]byte(gModuleRe), &modules)
		if nil != err {
			blog.Error("GetProcessBindModule Module  error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleSelectFailed)
		}
		moduleArr := modules.Data.Info
		gProc2ModuleURL := cli.CC.ProcCtrl() + "/process/v1/module/search"
		condition[common.BKProcIDField] = procID
		sCondJson, _ = json.Marshal(condition)
		blog.Info("get module config query url: %v", gModuleURL)
		blog.Info("get module config params: %v", string(sCondJson))
		gPorc2ModuleRe, err := httpcli.ReqHttp(req, gProc2ModuleURL, common.HTTPSelectPost, []byte(sCondJson))
		blog.Info("get module config return: %v", gPorc2ModuleRe)
		if nil != err {
			blog.Error("get module config params  error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcSelectBindToMoudleFaile)
		}
		var pro2Module ProcModuleResult
		err = json.Unmarshal([]byte(gPorc2ModuleRe), &pro2Module)
		if nil != err {
			blog.Error("get module config params  error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcSelectBindToMoudleFaile)
		}
		procModuleData := pro2Module.Data
		disModuleNameArr := make([]string, 0)
		for _, i := range moduleArr {
			if !util.InArray(i[common.BKModuleNameField], disModuleNameArr) {
				moduleName, ok := i[common.BKModuleNameField].(string)
				if false == ok {
					continue
				}
				isDefault64, ok := i[common.BKDefaultField].(float64)
				if false == ok {
					isDefault, ok := i[common.BKDefaultField].(int)
					if false != ok || 0 != isDefault {
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
		for _, j := range disModuleNameArr {
			num := 0
			isBind := 0
			for _, k := range moduleArr {
				moduleName, ok := k[common.BKModuleNameField].(string)
				if false == ok {
					continue
				}
				if j == moduleName {
					num++
				}
			}
			for _, m := range procModuleData {
				if j == m.ModuleName {
					isBind = 1
					break
				}
			}
			data := make(map[string]interface{})
			data[common.BKModuleNameField] = j
			data["set_num"] = num
			data["is_bind"] = isBind
			result = append(result, data)
		}
		return http.StatusOK, result, nil
	}, resp)
}
