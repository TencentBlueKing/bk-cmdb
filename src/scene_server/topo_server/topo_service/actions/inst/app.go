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
	"configcenter/src/common/errors"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	sencecommon "configcenter/src/scene_server/common"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/auditlog"
	"configcenter/src/source_controller/api/metadata"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tidwall/gjson"

	"io/ioutil"
	"strconv"
	"strings"

	api "configcenter/src/source_controller/api/object"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var app = &appAction{}

type appAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/app/{owner_id}", Params: nil, Handler: app.CreateApp})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/app/{owner_id}/{app_id}", Params: nil, Handler: app.DeleteApp})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/app/{owner_id}/{app_id}", Params: nil, Handler: app.UpdateApp})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/app/status/{flag}/{owner_id}/{app_id}", Params: nil, Handler: app.UpdateAppDataStatus})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/app/search/{owner_id}", Params: nil, Handler: app.SearchApp})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/app/default/{owner_id}/search", Params: nil, Handler: app.GetDefaultApp})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/app/default/{owner_id}", Params: nil, Handler: app.CreateDefaultApp})

	// create CC object
	app.CreateAction()
}

//delete application
func (cli *appAction) DeleteApp(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get error code in language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		//new feature, app not allow deletion
		blog.Error("app not allow deletion")
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppDeleteFailed)

		forward := &api.ForwardParam{Header: req.Request.Header}
		pathParams := req.PathParameters()
		appID, _ := strconv.Atoi(pathParams["app_id"])
		ownerID, _ := pathParams["owner_id"]
		user := sencecommon.GetUserFromHeader(req)

		// check wether it can be delete
		rstOk, rstErr := hasHost(req, cli.CC.HostCtrl(), map[string][]int{common.BKAppIDField: []int{appID}})
		if nil != rstErr {
			blog.Error("failed to check app wether it has hosts, error info is %s", rstErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoHasHostCheckFailed)
		}

		if !rstOk {
			blog.Error("failed to delete app, because of it has some hosts")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoHasHostCheckFailed)
		}

		// take snapshot before operation
		preData, retStrErr := inst.getInstDetail(req, appID, common.BKInnerObjIDApp, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
		}
		appData, ok := preData.(map[string]interface{})
		if false == ok {
			blog.Error("failed to get app detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppDeleteFailed)
		}
		appNameI, ok := appData[common.BKAppNameField]
		if false == ok {
			blog.Error("failed to get app detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppDeleteFailed)
		}
		bkAppName, ok := appNameI.(string)
		if false == ok {
			blog.Error("failed to get app detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppDeleteFailed)
		}
		if common.BKAppName == bkAppName {
			blog.Error("failed to delete bk default app")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoBkAppNotAllowedDelete)

		}
		//delete app
		input := make(map[string]interface{})
		input[common.BKAppIDField] = appID
		dAppURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDApp
		inputJSON, _ := json.Marshal(input)
		blog.Info("delete app url:%s", dAppURL)
		_, err := httpcli.ReqHttp(req, dAppURL, common.HTTPDelete, []byte(inputJSON))
		if nil != err {
			blog.Error("delete app error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppDeleteFailed)
		}
		{
			// save change log
			instID, _ := strconv.Atoi(fmt.Sprint(appID))
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDApp)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditSaveLogFaile)
			}

			auditContent := metadata.Content{
				PreData: preData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header).AuditObjLog(instID, auditContent, "delete app", common.BKInnerObjIDApp, ownerID, "0", user, auditoplog.AuditOpTypeDel)
		}
		//delete set in app
		setInput := make(map[string]interface{})
		setInput[common.BKAppIDField] = appID
		inputSetJSON, _ := json.Marshal(setInput)
		dSetURL := cli.CC.ObjCtrl() + "/object/v1/insts/set"
		_, err = httpcli.ReqHttp(req, dSetURL, common.HTTPDelete, []byte(inputSetJSON))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppDeleteFailed)
		}
		//delete module in app
		moduleInput := make(map[string]interface{})
		moduleInput[common.BKAppIDField] = appID
		inputModuleJSON, _ := json.Marshal(moduleInput)
		dModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/module"
		_, err = httpcli.ReqHttp(req, dModuleURL, common.HTTPDelete, []byte(inputModuleJSON))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleDeleteFailed)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

// update  application data status
func (cli *appAction) UpdateAppDataStatus(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get error code in language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		pathParams := req.PathParameters()
		appID, _ := strconv.Atoi(pathParams["app_id"])
		ownerID, _ := pathParams["owner_id"]
		flag, _ := pathParams["flag"]
		data := make(map[string]interface{})
		var appName string
		if common.DataStatusFlag(flag) != common.DataStatusDisabled && common.DataStatusFlag(flag) != common.DataStatusEnable {
			blog.Error("input params error:")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// check wether it can be delete
		rstOk, rstErr := hasHost(req, cli.CC.HostCtrl(), map[string][]int{common.BKAppIDField: []int{appID}})
		if nil != rstErr {
			blog.Error("failed to check app wether it has hosts, error info is %s", rstErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoHasHostCheckFailed)
		}

		if !rstOk {
			blog.Error("failed to delete app, because of it has some hosts")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoHasHostCheckFailed)
		}

		if common.DataStatusFlag(flag) == common.DataStatusEnable {
			condition := make(map[string]interface{})
			searchParams := make(map[string]interface{})
			condition[common.BKAppIDField] = appID
			searchParams["condition"] = condition

			//get app by appid
			sAppURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
			inputJSON, _ := json.Marshal(searchParams)
			appInfo, err := httpcli.ReqHttp(req, sAppURL, common.HTTPSelectPost, []byte(inputJSON))
			blog.Infof("get app params: %s", string(inputJSON))
			if nil != err {
				blog.Errorf("get app error: %v", err)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppSearchFailed)
			}
			blog.Infof("get app return %s", string(appInfo))
			appName = gjson.Get(string(appInfo), "data.info.0.bk_biz_name").String()

			//valid update name
			data[common.BKAppNameField] = appName + "(" + common.BKBizRecovery + ")"
			valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDApp, cli.CC.ObjCtrl(), forward, defErr)
			_, err = valid.ValidMap(data, common.ValidUpdate, appID)
			if nil != err {
				blog.Errorf("update app vaild error:%s", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommFieldNotValid)
			}

		}
		user := sencecommon.GetUserFromHeader(req)

		//update app
		input := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition[common.BKAppIDField] = appID
		condition[common.BKOwnerIDField] = ownerID
		data[common.BKDataStatusField] = flag

		// take snapshot before operation
		preData, retStrErr := inst.getInstDetail(req, appID, common.BKInnerObjIDApp, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
		}

		input["condition"] = condition
		input["data"] = data
		uAppURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDApp
		inputJSON, _ := json.Marshal(input)
		_, err := httpcli.ReqHttp(req, uAppURL, common.HTTPUpdate, []byte(inputJSON))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppUpdateFailed)
		}

		{
			// save change log
			instID, _ := strconv.Atoi(fmt.Sprint(appID))
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDApp)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
			}

			curData, retStrErr := inst.getInstDetail(req, instID, common.BKInnerObjIDApp, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditSaveLogFaile)
			}

			auditContent := metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header).AuditObjLog(instID, auditContent, "update app", common.BKInnerObjIDApp, ownerID, "0", user, auditoplog.AuditOpTypeModify)
		}

		return http.StatusOK, nil, nil
	}, resp)

}

//update application
func (cli *appAction) UpdateApp(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get error code in language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		pathParams := req.PathParameters()
		appID, _ := strconv.Atoi(pathParams["app_id"])
		ownerID, _ := pathParams["owner_id"]
		user := sencecommon.GetUserFromHeader(req)
		//update app
		input := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition[common.BKAppIDField] = appID
		condition[common.BKOwnerIDField] = ownerID
		value, err := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		data, _ := js.Map()
		valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDApp, cli.CC.ObjCtrl(), forward, defErr)
		_, err = valid.ValidMap(data, common.ValidUpdate, appID)
		if nil != err {
			blog.Errorf("UpdateApp vaild error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommFieldNotValid)
		}

		// take snapshot before operation
		preData, retStrErr := inst.getInstDetail(req, appID, common.BKInnerObjIDApp, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
		}

		appData, ok := preData.(map[string]interface{})
		if false == ok {
			blog.Error("failed to get app detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppUpdateFailed)
		}
		appNameI, ok := appData[common.BKAppNameField]
		if false == ok {
			blog.Error("failed to get app detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppUpdateFailed)
		}
		bkAppName, ok := appNameI.(string)
		if false == ok {
			blog.Error("failed to get app detail")
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppUpdateFailed)
		}
		if common.BKAppName == bkAppName {
			_, ok := data[common.BKAppNameField]
			if ok {
				delete(data, common.BKAppNameField)
			}
		}

		input["condition"] = condition
		input["data"] = data
		uAppURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDApp
		inputJSON, _ := json.Marshal(input)
		_, err = httpcli.ReqHttp(req, uAppURL, common.HTTPUpdate, []byte(inputJSON))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppUpdateFailed)
		}

		{
			// save change log
			instID, _ := strconv.Atoi(fmt.Sprint(appID))
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDApp)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
			}

			curData, retStrErr := inst.getInstDetail(req, instID, common.BKInnerObjIDApp, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditSaveLogFaile)
			}

			auditContent := metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header).AuditObjLog(instID, auditContent, "update app", common.BKInnerObjIDApp, ownerID, "0", user, auditoplog.AuditOpTypeModify)
		}

		return http.StatusOK, nil, nil
	}, resp)

}

func (cli *appAction) getOwnerIDByAppID(req *restful.Request, appID int) (ownerID string) {
	condition := map[string]interface{}{}
	condition[common.BKAppIDField] = appID
	sAppURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
	inputJSON, _ := json.Marshal(map[string]interface{}{"condition": condition})
	appInfo, err := httpcli.ReqHttp(req, sAppURL, common.HTTPSelectPost, []byte(inputJSON))
	if nil != err {
		blog.Error("search app error: %v", err)
		return
	}
	return gjson.Get(appInfo, "data.info.0."+common.BKOwnerIDField).String()
}

//search application
func (cli *appAction) SearchApp(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID, _ := pathParams["owner_id"]
		value, _ := ioutil.ReadAll(req.Request.Body)
		var js params.SearchParams
		err := json.Unmarshal([]byte(value), &js)
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		var condition map[string]interface{}
		if 1 == js.Native {
			condition = js.Condition
		} else {
			condition = params.ParseAppSearchParams(js.Condition)
		}

		//search app in enable status default
		_, ok := condition[common.BKDataStatusField]
		if !ok {
			condition[common.BKDataStatusField] = map[string]interface{}{common.BKDBNE: common.DataStatusDisabled}
		}

		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKDefaultField] = 0
		page := js.Page
		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		searchParams["fields"] = strings.Join(js.Fields, ",")
		searchParams["start"] = page["start"]
		searchParams["limit"] = page["limit"]
		searchParams["sort"] = page["sort"]
		//search app
		sAppURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
		inputJSON, _ := json.Marshal(searchParams)
		appInfo, err := httpcli.ReqHttp(req, sAppURL, common.HTTPSelectPost, []byte(inputJSON))
		blog.Debug("search app params: %s", string(inputJSON))
		if nil != err {
			blog.Error("search app error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppSearchFailed)
		}
		blog.Debug("search app return %v", appInfo)
		// replace the association id to name
		retstr, retStrErr := inst.getInstDetails(req, common.BKInnerObjIDApp, ownerID, appInfo, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})
		if common.CCSuccess != retStrErr {
			blog.Error("search app error: %v", retStrErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppSearchFailed)
		}
		blog.Info("search app return %v", retstr)
		return http.StatusOK, retstr["data"], nil
	}, resp)
}

//create application
func (cli *appAction) CreateApp(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		pathParams := req.PathParameters()
		ownerID := pathParams["owner_id"]
		user := sencecommon.GetUserFromHeader(req)
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input, err := js.Map()
		valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDApp, cli.CC.ObjCtrl(), forward, defErr)
		_, err = valid.ValidMap(input, common.ValidCreate, 0)
		if nil != err {
			blog.Errorf("create app valid eror:%s, data:%v", err.Error(), string(value))
			if _, ok := err.(errors.CCErrorCoder); ok {
				return http.StatusInternalServerError, nil, err
			}
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommFieldNotValid)
		}
		input[common.BKOwnerIDField] = ownerID
		input[common.BKDefaultField] = 0
		input[common.BKSupplierIDField] = common.BKDefaultSupplierID
		appInfoJSON, _ := json.Marshal(input)
		cAppURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDApp
		cAppRes, err := httpcli.ReqHttp(req, cAppURL, common.HTTPCreate, []byte(appInfoJSON))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppCreateFailed)
		}
		js, err = simplejson.NewJson([]byte(cAppRes))
		appResData, _ := js.Map()
		appIDInfo := appResData["data"].(map[string]interface{})
		appID := appIDInfo[common.BKAppIDField]
		{
			// save change log
			instID, _ := strconv.Atoi(fmt.Sprint(appID))
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDApp)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditTakeSnapshotFaile)
			}

			curData, retStrErr := inst.getInstDetail(req, instID, common.BKInnerObjIDApp, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAuditSaveLogFaile)
			}
			auditContent := metadata.Content{
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header).AuditObjLog(instID, auditContent, "create app", common.BKInnerObjIDApp, ownerID, "0", user, auditoplog.AuditOpTypeAdd)
		}
		//create default set
		inputSetInfo := make(map[string]interface{})
		inputSetInfo[common.BKAppIDField] = appID
		inputSetInfo[common.BKInstParentStr] = appID
		inputSetInfo[common.BKSetNameField] = common.DefaultResSetName
		inputSetInfo[common.BKDefaultField] = common.DefaultResSetFlag
		inputSetInfo[common.BKOwnerIDField] = ownerID
		cSetURL := cli.CC.ObjCtrl() + "/object/v1/insts/set"
		setJSONData, _ := json.Marshal(inputSetInfo)
		cSetRes, err := httpcli.ReqHttp(req, cSetURL, common.HTTPCreate, []byte(setJSONData))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoSetCreateFailed)
		}
		//create default module
		js, err = simplejson.NewJson([]byte(cSetRes))
		setResData, _ := js.Map()
		setIDInfo := setResData["data"].(map[string]interface{})
		setID := setIDInfo[common.BKSetIDField]
		inputResModuleInfo := make(map[string]interface{})
		inputResModuleInfo[common.BKSetIDField] = setID
		inputResModuleInfo[common.BKInstParentStr] = setID
		inputResModuleInfo[common.BKAppIDField] = appID
		inputResModuleInfo[common.BKModuleNameField] = common.DefaultResModuleName
		inputResModuleInfo[common.BKDefaultField] = common.DefaultResModuleFlag
		inputResModuleInfo[common.BKOwnerIDField] = ownerID
		cModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/module"
		resModuleJSONData, _ := json.Marshal(inputResModuleInfo)
		_, err = httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, []byte(resModuleJSONData))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleCreateFailed)
		}
		inputFaultModuleInfo := make(map[string]interface{})
		inputFaultModuleInfo[common.BKSetIDField] = setID
		inputFaultModuleInfo[common.BKInstParentStr] = setID
		inputFaultModuleInfo[common.BKAppIDField] = appID
		inputFaultModuleInfo[common.BKModuleNameField] = common.DefaultFaultModuleName
		inputFaultModuleInfo[common.BKDefaultField] = common.DefaultFaultModuleFlag
		inputFaultModuleInfo[common.BKOwnerIDField] = ownerID
		resFaultModuleJSONData, _ := json.Marshal(inputFaultModuleInfo)
		_, err = httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, []byte(resFaultModuleJSONData))
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppCreateFailed)
		}
		result := make(map[string]interface{})
		result[common.BKAppIDField] = appID

		return http.StatusOK, result, nil
	}, resp)
}

//get default application
func (cli *appAction) GetDefaultApp(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID, _ := pathParams["owner_id"]
		value, _ := ioutil.ReadAll(req.Request.Body)
		var js params.SearchParams
		err := json.Unmarshal([]byte(value), &js)
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		condition := js.Condition
		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKDefaultField] = common.DefaultAppFlag
		page := js.Page
		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		searchParams["fields"] = strings.Join(js.Fields, ",")
		searchParams["start"] = page["start"]
		searchParams["limit"] = page["limit"]
		searchParams["sort"] = page["sort"]
		//search app
		sAppURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
		inputJSON, _ := json.Marshal(searchParams)
		appInfo, err := httpcli.ReqHttp(req, sAppURL, common.HTTPSelectPost, []byte(inputJSON))
		blog.Info("get default app params: %s", string(inputJSON))
		if nil != err {
			blog.Error("search app error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppSearchFailed)
		}
		blog.Info("get default a app return %v", appInfo)
		appJson, err := simplejson.NewJson([]byte(appInfo))
		appResData, _ := appJson.Map()
		return http.StatusOK, appResData["data"], nil
	}, resp)

}

//create default application
func (cli *appAction) CreateDefaultApp(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		pathParams := req.PathParameters()
		ownerID := pathParams["owner_id"]
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			blog.Errorf("create default app get params error %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input, err := js.Map()
		valid := validator.NewValidMap(ownerID, common.BKInnerObjIDApp, cli.CC.ObjCtrl(), forward, defErr)
		_, err = valid.ValidMap(input, common.ValidCreate, 0)
		if nil != err {
			blog.Errorf("create default app get params error %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommFieldNotValid)
		}

		//create application
		input[common.BKOwnerIDField] = ownerID
		input[common.BKSupplierIDField] = common.BKDefaultSupplierID
		input[common.BKDefaultField] = common.DefaultAppFlag
		appInfoJSON, _ := json.Marshal(input)
		cAppURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDApp
		cAppRes, err := httpcli.ReqHttp(req, cAppURL, common.HTTPCreate, []byte(appInfoJSON))
		if nil != err {
			blog.Errorf("add default application error, ownerID:%s, error:%v ", ownerID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppCreateFailed)
		}
		//create default set
		js, err = simplejson.NewJson([]byte(cAppRes))
		appResData, _ := js.Map()
		appIDInfo := appResData["data"].(map[string]interface{})
		appID := appIDInfo[common.BKAppIDField]
		inputSetInfo := make(map[string]interface{})
		inputSetInfo[common.BKAppIDField] = appID
		inputSetInfo[common.BKInstParentStr] = appID
		inputSetInfo[common.BKSetNameField] = common.DefaultResSetName
		inputSetInfo[common.BKDefaultField] = common.DefaultResSetFlag
		inputSetInfo[common.BKOwnerIDField] = ownerID
		cSetURL := cli.CC.ObjCtrl() + "/object/v1/insts/set"
		setJSONData, _ := json.Marshal(inputSetInfo)
		cSetRes, err := httpcli.ReqHttp(req, cSetURL, common.HTTPCreate, []byte(setJSONData))
		if nil != err {
			blog.Errorf("add default application Set error, ownerID:%s, error:%v ", ownerID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppCreateFailed)
		}
		//create default module
		js, err = simplejson.NewJson([]byte(cSetRes))
		setResData, _ := js.Map()
		setIDInfo := setResData["data"].(map[string]interface{})
		setID := setIDInfo[common.BKSetIDField]
		inputResModuleInfo := make(map[string]interface{})
		inputResModuleInfo[common.BKSetIDField] = setID
		inputResModuleInfo[common.BKInstParentStr] = setID
		inputResModuleInfo[common.BKAppIDField] = appID
		inputResModuleInfo[common.BKModuleNameField] = common.DefaultResModuleName
		inputResModuleInfo[common.BKDefaultField] = common.DefaultResModuleFlag
		inputResModuleInfo[common.BKOwnerIDField] = ownerID
		cModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/module"
		resModuleJSONData, _ := json.Marshal(inputResModuleInfo)
		_, err = httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, []byte(resModuleJSONData))
		if nil != err {
			blog.Errorf("add default application module error, ownerID:%s, error:%v ", ownerID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppCreateFailed)
		}
		inputFaultModuleInfo := make(map[string]interface{})
		inputFaultModuleInfo[common.BKSetIDField] = setID
		inputFaultModuleInfo[common.BKInstParentStr] = setID
		inputFaultModuleInfo[common.BKAppIDField] = appID
		inputFaultModuleInfo[common.BKModuleNameField] = common.DefaultFaultModuleName
		inputFaultModuleInfo[common.BKDefaultField] = common.DefaultFaultModuleFlag
		inputFaultModuleInfo[common.BKOwnerIDField] = ownerID
		resFaultModuleJSONData, _ := json.Marshal(inputFaultModuleInfo)
		_, err = httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, []byte(resFaultModuleJSONData))
		if nil != err {
			blog.Errorf("add default application module error, ownerID:%s, error info is %v ", ownerID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppCreateFailed)
		}

		if ownerID != common.BKDefaultOwnerID {
			headerOwner := util.GetActionOnwerID(req)
			blog.Infof("copy asst for %s, header owner: %s", ownerID, headerOwner)
			searchAsstURL := cli.CC.ObjCtrl() + "/object/v1/meta/objectassts"
			searchAsstCondition := map[string]interface{}{
				common.BKOwnerIDField: common.BKDefaultOwnerID,
			}
			searchAsstData, _ := json.Marshal(searchAsstCondition)
			searchAsstReply, err := httpcli.ReqHttp(req, searchAsstURL, common.HTTPSelectPost, []byte(searchAsstData))
			if nil != err {
				blog.Errorf("add default application module error, ownerID:%s, error:%v ", ownerID, err)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppCreateFailed)
			}

			searchAsstJSON := gjson.Parse(searchAsstReply)
			if !searchAsstJSON.Get("result").Bool() {
				return http.StatusInternalServerError, nil, defErr.Error(int(searchAsstJSON.Get(common.HTTPBKAPIErrorCode).Int()))
			}

			assts := []map[string]interface{}{}
			json.Unmarshal([]byte(searchAsstJSON.Get("data").String()), &assts)

			blog.Infof("copy asst for %s, %+v", ownerID, assts)

			for index := range assts {
				assts[index][common.BKOwnerIDField] = ownerID

				createAsstURL := cli.CC.ObjCtrl() + "/object/v1/meta/objectasst"
				createAsstData, _ := json.Marshal(assts[index])
				createAsstReply, err := httpcli.ReqHttp(req, createAsstURL, common.HTTPSelectPost, []byte(createAsstData))
				if nil != err {
					blog.Errorf("add default application module error, ownerID:%s, error:%v ", ownerID, err)
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppCreateFailed)
				}
				createAsstJSON := gjson.Parse(createAsstReply)
				if !createAsstJSON.Get("result").Bool() {
					return http.StatusInternalServerError, nil, defErr.Error(int(createAsstJSON.Get(common.HTTPBKAPIErrorCode).Int()))
				}
			}

		}

		result := make(map[string]interface{})
		result[common.BKAppIDField] = appID
		return http.StatusOK, result, nil
	}, resp)
}
