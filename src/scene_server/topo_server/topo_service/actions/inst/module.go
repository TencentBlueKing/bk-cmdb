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

package inst

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/auditlog"
	"configcenter/src/source_controller/api/metadata"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

var module = &moduleAction{}

type moduleAction struct {
	base.BaseAction
}

func init() {

	// init action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/module/{app_id}/{set_id}", Params: nil, Handler: module.CreateModule})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/module/{app_id}/{set_id}/{module_id}", Params: nil, Handler: module.DeleteModule})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/module/{app_id}/{set_id}/{module_id}", Params: nil, Handler: module.UpdateModule})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/module/search/{owner_id}/{app_id}/{set_id}", Params: nil, Handler: module.SearchModule})

	// set cc interface
	module.CreateAction()
}

// CreateModule
func (cli *moduleAction) CreateModule(req *restful.Request, resp *restful.Response) {

	blog.Debug("create module")

	// get language
	language := util.GetActionLanguage(req)
	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}
	user := util.GetActionUser(req)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		//create default module
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		js, err := simplejson.NewJson(value)
		if nil != err {
			blog.Error("the input json is invalid, error info is %s", err.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		input, jsErr := js.Map()
		if nil != jsErr {
			blog.Error("the input json is invalid, error info is %s", jsErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		setID, convErr := strconv.Atoi(req.PathParameter("set_id"))
		if nil != convErr {
			blog.Error("the setid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "set_id")
		}

		appID, convErr := strconv.Atoi(req.PathParameter("app_id"))
		if nil != convErr {
			blog.Error("the appid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "app_id")
		}

		if _, ok := input[common.BKOwnerIDField]; !ok {
			blog.Error("not set %s", common.BKOwnerIDField)
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKOwnerIDField)
		}

		if _, ok := input[common.BKModuleNameField]; !ok {
			blog.Error("not set ModuleName")
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKModuleNameField)
		}

		if _, ok := input[common.BKInstParentStr]; !ok {
			blog.Error("not set %s", common.BKInstParentStr)
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKInstParentStr)
		}

		tmpID, ok := input[common.BKOwnerIDField].(string)
		if !ok {
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedString, common.BKOwnerIDField)
		}

		// create
		input[common.BKSetIDField] = setID
		input[common.BKAppIDField] = appID
		// check
		valid := validator.NewValidMapWithKeyFields(tmpID, common.BKInnerObjIDModule, cli.CC.ObjCtrl(), []string{common.BKOwnerIDField, common.BKInstParentStr}, forward, defErr)
		_, err = valid.ValidMap(input, common.ValidCreate, 0)
		if nil != err {
			blog.Error("failed to valide, error is %s", err.Error())
			return http.StatusBadRequest, "", err
		}

		// create
		input[common.BKDefaultField] = 0
		input[common.CreateTimeField] = util.GetCurrentTimeStr()

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the json, error is info %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		cModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/module"

		moduleRes, err := httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, inputJSON)
		if nil != err {
			blog.Error("failed to create the module, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleCreateFailed)
		}

		{
			// save change log
			instID := gjson.Get(moduleRes, "data."+common.BKModuleIDField).Int()
			if instID == 0 {
				blog.Errorf("inst id not found")
			}
			ownerID := app.getOwnerIDByAppID(req, appID)
			if ownerID == "" {
				blog.Errorf("owner id not found")
			}
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDModule)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(attErr)
			}

			curData, retStrErr := inst.getInstDetail(req, int(instID), common.BKInnerObjIDModule, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleCreateFailed)
			}
			auditContent := metadata.Content{
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditModuleLog(instID, auditContent, "create module", ownerID, fmt.Sprint(appID), user, auditoplog.AuditOpTypeAdd)
		}
		return http.StatusOK, moduleRes, nil
	}, resp)

}

// DeleteModule delete module by condition
func (cli *moduleAction) DeleteModule(req *restful.Request, resp *restful.Response) {

	blog.Debug("delete module")

	// get the language
	language := util.GetActionLanguage(req)
	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	user := util.GetActionUser(req)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		appID, convErr := strconv.Atoi(req.PathParameter("app_id"))
		if nil != convErr {
			blog.Error("the appid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "app_id")
		}

		setID, convErr := strconv.Atoi(req.PathParameter("set_id"))
		if nil != convErr {
			blog.Error("the setid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "set_id")
		}

		moduleID, convErr := strconv.Atoi(req.PathParameter("module_id"))
		if nil != convErr {
			blog.Error("the moduleid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "module_id")
		}

		// check wether it can be delete
		rstOk, rstErr := hasHost(req, cli.CC.HostCtrl(), map[string][]int{
			common.BKAppIDField:    []int{appID},
			common.BKModuleIDField: []int{moduleID},
			common.BKSetIDField:    []int{setID},
		})
		if nil != rstErr {
			blog.Error("failed to check module wether it has hosts, error info is %s", rstErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrTopoHasHostCheckFailed)
		}

		if !rstOk {
			blog.Error("failed to delete module, because of it has some hosts")
			return http.StatusBadRequest, "", defErr.Error(common.CCErrTopoHasHost)
		}

		// take snapshot before operation
		ownerID := app.getOwnerIDByAppID(req, appID)
		if ownerID == "" {
			blog.Errorf("owner id not found")
		}
		preData, retStrErr := inst.getInstDetail(req, moduleID, common.BKInnerObjIDModule, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		//delete module
		input := make(map[string]interface{})
		input[common.BKAppIDField] = appID
		input[common.BKSetIDField] = setID
		input[common.BKModuleIDField] = moduleID

		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/module"

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		moduleRes, err := httpcli.ReqHttp(req, uURL, "DELETE", []byte(inputJSON))
		if nil != err {
			blog.Error("failed to delete the module, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleDeleteFailed)
		}

		{
			// save change log
			instID := gjson.Get(moduleRes, "data.bk_module_id").Int()
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDModule)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleDeleteFailed)
			}
			auditContent := metadata.Content{
				PreData: preData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditModuleLog(instID, auditContent, "delete module", ownerID, fmt.Sprint(appID), user, auditoplog.AuditOpTypeDel)
		}
		return http.StatusOK, moduleRes, nil

	}, resp)

}

// UpdateModule
func (cli *moduleAction) UpdateModule(req *restful.Request, resp *restful.Response) {
	blog.Debug("update module")

	// get language
	language := util.GetActionLanguage(req)
	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	user := util.GetActionUser(req)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		appID, convErr := strconv.Atoi(req.PathParameter("app_id"))
		if nil != convErr {
			blog.Error("the appid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "app_id")
		}

		setID, convErr := strconv.Atoi(req.PathParameter("set_id"))
		if nil != convErr {
			blog.Error("the setid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "set_id")
		}

		moduleID, _ := strconv.Atoi(req.PathParameter("module_id"))
		if nil != convErr {
			blog.Error("the moduleid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "module_id")
		}

		//update module
		input := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition[common.BKAppIDField] = appID
		condition[common.BKSetIDField] = setID
		condition[common.BKModuleIDField] = moduleID

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("failed to read the http request, error info is %s", readErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			blog.Error("failed to create simplejson, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		data, jsErr := js.Map()
		if nil != jsErr {
			blog.Error("failed to unmarshal data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		data[common.BKAppIDField] = appID
		data[common.BKSetIDField] = setID
		valid := validator.NewValidMapWithKeyFields(common.BKDefaultOwnerID, common.BKInnerObjIDModule, cli.CC.ObjCtrl(), []string{common.BKOwnerIDField, common.BKInstParentStr, common.BKModuleNameField}, forward, defErr)
		_, err = valid.ValidMap(data, common.ValidUpdate, moduleID)

		if nil != err {
			blog.Error("failed to valid the input , error is %s", err.Error())
			return http.StatusBadRequest, "", err
		}

		// take snapshot before operation
		ownerID := app.getOwnerIDByAppID(req, appID)
		if ownerID == "" {
			blog.Errorf("owner id not found")
		}
		preData, retStrErr := inst.getInstDetail(req, moduleID, common.BKInnerObjIDModule, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		input["condition"] = condition
		input["data"] = data
		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/module"
		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		moduleRes, err := httpcli.ReqHttp(req, uURL, "PUT", []byte(inputJSON))

		if nil != err {
			blog.Error("failed to update the module, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleUpdateFailed)
		}

		if rsp, ok := cli.IsSuccess([]byte(moduleRes)); !ok {
			blog.Error("failed to update the module, error info is %v", rsp.Message)
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleUpdateFailed)
		}

		{
			// save change log
			instID := moduleID
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDModule)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleCreateFailed)
			}

			curData, retStrErr := inst.getInstDetail(req, instID, common.BKInnerObjIDModule, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleCreateFailed)
			}
			auditContent := metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditModuleLog(instID, auditContent, "update module", ownerID, fmt.Sprint(appID), user, auditoplog.AuditOpTypeModify)
		}

		return http.StatusOK, nil, nil

	}, resp)

	return

}

// SearfhModule search modules
func (cli *moduleAction) SearchModule(req *restful.Request, resp *restful.Response) {
	blog.Debug("search module")

	// get the language
	language := util.GetActionLanguage(req)
	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		ownerID := req.PathParameter("owner_id")
		appID, convErr := strconv.Atoi(req.PathParameter("app_id"))
		if nil != convErr {
			blog.Error("the appid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "app_id")
		}
		setID, convErr := strconv.Atoi(req.PathParameter("set_id"))
		if nil != convErr {
			blog.Error("the setid is invalid, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "set_id")
		}

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("failed to read the http request, error info is %s", readErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		var js params.SearchParams
		err := json.Unmarshal([]byte(value), &js)
		if nil != err {
			blog.Error("failed to unmarshal the input, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		condition := params.ParseAppSearchParams(js.Condition)

		condition[common.BKAppIDField] = appID
		condition[common.BKSetIDField] = setID

		page := js.Page

		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		searchParams["fields"] = strings.Join(js.Fields, ",")
		searchParams["start"] = page["start"]
		searchParams["limit"] = page["limit"]
		searchParams["sort"] = page["sort"]

		//search
		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/module/search"
		inputJSON, jsErr := json.Marshal(searchParams)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		moduleRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
		if nil != err {
			blog.Error("failed to update the module, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleSelectFailed)
		}

		// replace the association id to name
		retStr, retStrErr := inst.getInstDetails(req, common.BKInnerObjIDModule, ownerID, moduleRes, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})
		if common.CCSuccess != retStrErr {
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		return http.StatusOK, retStr["data"], nil
	}, resp)

	return
}
