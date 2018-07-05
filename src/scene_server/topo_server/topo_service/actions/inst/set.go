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

var set = &setAction{}

type setAction struct {
	base.BaseAction
}

func init() {

	// init action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/set/{app_id}", Params: nil, Handler: set.CreateSet})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/set/{app_id}/{set_id}", Params: nil, Handler: set.DeleteSet})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/set/{app_id}/{set_id}", Params: nil, Handler: set.UpdateSet})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/set/search/{owner_id}/{app_id}", Params: nil, Handler: set.SearchSet})

	// set cc interface
	set.CreateAction()
}

// CreateModule
func (cli *setAction) CreateSet(req *restful.Request, resp *restful.Response) {

	blog.Debug("create set")

	// get language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	user := util.GetActionUser(req)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		forward := &api.ForwardParam{Header: req.Request.Header}
		//create default module
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		js, err := simplejson.NewJson(value)
		if nil != err {
			blog.Error("failed to unmarshal the data , error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input, jsonErr := js.Map()
		if nil != jsonErr {
			blog.Error("failed to unmarshal the data , error info is %s", jsonErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		_, setOK := input[common.BKSetNameField]
		if !setOK {
			blog.Errorf("not set '%s'", common.BKSetNameField)
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKSetNameField)
		}

		ownerID, ownerOK := input[common.BKOwnerIDField]
		if !ownerOK {
			blog.Error("'%s' field must be setted", common.BKOwnerIDField)
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKOwnerIDField)
		}

		_, parentOK := input[common.BKInstParentStr]
		if !parentOK {
			blog.Error("'%s' field must be setted", common.BKInstParentStr)
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsLostField, common.BKInstParentStr)
		}

		appID, convErr := strconv.Atoi(req.PathParameter("app_id"))
		if nil != convErr {
			blog.Error("failed to convert the appid to int, error info is %s", convErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "app_id")
		}

		tmpID, ok := ownerID.(string)
		if !ok {
			blog.Error("'OwnerID' must be a string value")
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKOwnerIDField)
		}

		input[common.BKAppIDField] = appID
		// check
		valid := validator.NewValidMapWithKeyFields(tmpID, common.BKInnerObjIDSet, cli.CC.ObjCtrl(), []string{common.BKInstParentStr, common.BKOwnerIDField}, forward, defErr)
		_, err = valid.ValidMap(input, common.ValidCreate, 0)
		if nil != err {
			blog.Error("failed to valid the input data, error info is %s", err.Error())
			return http.StatusBadRequest, "", err
		}

		// create
		input[common.BKDefaultField] = 0
		input[common.CreateTimeField] = util.GetCurrentTimeStr()

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		cModuleURL := cli.CC.ObjCtrl() + "/object/v1/insts/set"
		moduleRes, err := httpcli.ReqHttp(req, cModuleURL, common.HTTPCreate, inputJSON)
		if nil != err {
			blog.Error("failed to create the set, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetCreateFailed)
		}

		{
			// save change log
			instID := gjson.Get(moduleRes, "data."+common.BKSetIDField).Int()
			ownerID := fmt.Sprint(input[common.BKOwnerIDField])
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDSet)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetCreateFailed)
			}

			curData, retStrErr := inst.getInstDetail(req, int(instID), common.BKInnerObjIDSet, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetCreateFailed)
			}
			auditContent := metadata.Content{
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditSetLog(instID, auditContent, "create set", ownerID, fmt.Sprint(appID), user, auditoplog.AuditOpTypeAdd)
		}

		return http.StatusOK, moduleRes, nil
	}, resp)

}

// DeleteSet delete set by conditions
func (cli *setAction) DeleteSet(req *restful.Request, resp *restful.Response) {

	blog.Debug("delete set")

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

		operationInst := &operation{}
		operationInst.Delete.InstID = append(operationInst.Delete.InstID, setID)
		if setID < 0 { // if the inst less than zeor, it means to batch to delete the inst
			//create default module
			value, err := ioutil.ReadAll(req.Request.Body)
			if nil != err {
				blog.Error("read request body failed, error:%v", err)
				return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
			}
			if 0 == len(value) {

				blog.Error("read request body failed, it is empty")
				return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
			}
			if err = json.Unmarshal(value, operationInst); nil != err {
				blog.Errorf("failed to unmarshal the body params, error info is %s", err.Error())
				return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
			}
		}
		for _, operate := range operationInst.Delete.InstID {

			setID = operate

			// check wether it can be delete
			rstOk, rstErr := hasHost(req, cli.CC.HostCtrl(), map[string][]int{
				common.BKAppIDField: []int{appID},
				common.BKSetIDField: []int{setID},
			})
			if nil != rstErr {
				blog.Error("failed to check set wether it has hosts, error info is %s", rstErr.Error())
				return http.StatusBadRequest, "", defErr.Error(common.CCErrTopoHasHostCheckFailed)
			}

			if !rstOk {
				blog.Error("failed to delete set, because of it has some hosts")
				return http.StatusBadRequest, "", defErr.Error(common.CCErrTopoHasHost)
			}

			// take snapshot before operation
			ownerID := app.getOwnerIDByAppID(req, appID)
			if ownerID == "" {
				blog.Errorf("owner id not found")
			}
			preData, retStrErr := inst.getInstDetail(req, setID, common.BKInnerObjIDSet, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, "", defErr.Error(retStrErr)
			}

			//delete set
			input := make(map[string]interface{})
			input[common.BKAppIDField] = appID
			input[common.BKSetIDField] = setID

			uURL := cli.CC.ObjCtrl() + "/object/v1/insts/set"

			inputJSON, jsErr := json.Marshal(input)
			if nil != jsErr {
				blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
			}

			_, err := httpcli.ReqHttp(req, uURL, common.HTTPDelete, []byte(inputJSON))
			if nil != err {
				blog.Error("failed to delete the set, error info is %s", err.Error())
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetDeleteFailed)
			}

			//delete module
			input = make(map[string]interface{})
			input[common.BKAppIDField] = appID
			input[common.BKSetIDField] = setID

			uURL = cli.CC.ObjCtrl() + "/object/v1/insts/module"
			inputJSON, jsErr = json.Marshal(input)
			if nil != jsErr {
				blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
			}

			moduleRes, err := httpcli.ReqHttp(req, uURL, common.HTTPDelete, []byte(inputJSON))
			if nil != err {
				blog.Error("failed to delete the module, error info is %s", err.Error())
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoModuleDeleteFailed)
			}
			if rsp, ok := cli.IsSuccess([]byte(moduleRes)); !ok {
				blog.Error("failed to update the module, error info is %v", rsp.Message)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetUpdateFailed)
			}

			{
				// save change log
				instID := gjson.Get(moduleRes, "data.bk_set_id").Int()
				headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDSet)
				if common.CCSuccess != attErr {
					return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetDeleteFailed)
				}

				auditContent := metadata.Content{
					PreData: preData,
					Headers: headers,
				}
				auditlog.NewClient(cli.CC.AuditCtrl()).AuditSetLog(instID, auditContent, "delete set", ownerID, fmt.Sprint(appID), user, auditoplog.AuditOpTypeDel)
			}
		} // delete the set
		return http.StatusOK, nil, nil
	}, resp)
}

// UpdateSet update set by condition
func (cli *setAction) UpdateSet(req *restful.Request, resp *restful.Response) {
	blog.Debug("updte set")

	// get language
	language := util.GetActionLanguage(req)

	// get error by language
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

		//update set
		input := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition[common.BKAppIDField] = appID
		condition[common.BKSetIDField] = setID

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("read request body failed, error:%s", readErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		js, jsErr := simplejson.NewJson([]byte(value))
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		data, jsErr := js.Map()
		data[common.BKAppIDField] = appID
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}
		valid := validator.NewValidMapWithKeyFields(common.BKDefaultOwnerID, common.BKInnerObjIDSet, cli.CC.ObjCtrl(), []string{common.BKInstParentStr, common.BKOwnerIDField, common.BKSetNameField}, forward, defErr)
		_, err := valid.ValidMap(data, common.ValidUpdate, setID)
		if nil != err {
			blog.Error("failed to valid the input data, error info is %s", err.Error())
			return http.StatusBadRequest, "", err
		}

		// take snapshot before operation
		ownerID := app.getOwnerIDByAppID(req, appID)
		if ownerID == "" {
			blog.Errorf("owner id not found")
		}
		preData, retStrErr := inst.getInstDetail(req, setID, common.BKInnerObjIDSet, ownerID)
		if common.CCSuccess != retStrErr {
			blog.Errorf("get inst detail error: %v", retStrErr)
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		input["condition"] = condition
		input["data"] = data

		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/set"

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to marshal the data, error info is %s", jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		moduleRes, err := httpcli.ReqHttp(req, uURL, "PUT", []byte(inputJSON))
		if nil != err {
			blog.Error("failed to delete the set, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetUpdateFailed)
		}

		{
			// save change log
			instID := setID //gjson.Get(moduleRes, "data.bk_set_id").Int()
			//ownerID := fmt.Sprint(input[common.BKOwnerIDField])
			headers, attErr := inst.getHeader(forward, ownerID, common.BKInnerObjIDSet)
			if common.CCSuccess != attErr {
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetUpdateFailed)
			}

			curData, retStrErr := inst.getInstDetail(req, int(instID), common.BKInnerObjIDSet, ownerID)
			if common.CCSuccess != retStrErr {
				blog.Errorf("get inst detail error: %v", retStrErr)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetUpdateFailed)
			}
			auditContent := metadata.Content{
				PreData: preData,
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl()).AuditSetLog(instID, auditContent, "update set", ownerID, fmt.Sprint(appID), user, auditoplog.AuditOpTypeModify)
		}
		return http.StatusOK, moduleRes, nil
	}, resp)

}

// SearfhModule
func (cli *setAction) SearchSet(req *restful.Request, resp *restful.Response) {
	blog.Debug("search set")
	// get language
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

		value, readErr := ioutil.ReadAll(req.Request.Body)
		if nil != readErr {
			blog.Error("read request body failed, error:%s", readErr.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		var js params.SearchParams
		err := json.Unmarshal([]byte(value), &js)
		if nil != err {
			blog.Error("failed to unmarshal the data , error info is %s", err.Error())
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		condition := params.ParseAppSearchParams(js.Condition)

		condition[common.BKAppIDField] = appID

		page := js.Page

		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		searchParams["fields"] = strings.Join(js.Fields, ",")
		searchParams["start"] = page["start"]
		searchParams["limit"] = page["limit"]
		searchParams["sort"] = page["sort"]

		//search
		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/set/search"
		inputJSON, _ := json.Marshal(searchParams)
		moduleRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
		if nil != err {
			blog.Error("failed to select the set, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetSelectFailed)
		}

		// replace the association id into name
		retStr, retStrErr := inst.getInstDetails(req, common.BKInnerObjIDSet, ownerID, moduleRes, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})
		if common.CCSuccess != retStrErr {
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		return http.StatusOK, retStr["data"], nil

	}, resp)

}
