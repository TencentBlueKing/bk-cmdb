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
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/auditlog"
	"configcenter/src/source_controller/api/metadata"
	sourceAPI "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	api "configcenter/src/source_controller/api/object"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var process *procAction = &procAction{}

type procAction struct {
	base.BaseAction
	objcli *api.Client
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/{bk_supplier_account}/{bk_biz_id}", Params: nil, Handler: process.CreateProcess})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}", Params: nil, Handler: process.DeleteProcess})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/search/{bk_supplier_account}/{bk_biz_id}", Params: nil, Handler: process.SearchProcess})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}", Params: nil, Handler: process.UpdateProcess})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/{bk_supplier_account}/{bk_biz_id}", Params: nil, Handler: process.BatchUpdateProcess})
	process.CreateAction()
	process.objcli = api.NewClient("")
}

//UpdateProcess update process
func (cli *procAction) UpdateProcess(req *restful.Request, resp *restful.Response) {
	user := util.GetActionUser(req)
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &sourceAPI.ForwardParam{Header: req.Request.Header}
	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID := pathParams[common.BKOwnerIDField]
		appIDStr := pathParams[common.BKAppIDField]
		procIDStr := pathParams[common.BKProcIDField]
		appID, _ := strconv.Atoi(appIDStr)
		procID, _ := strconv.Atoi(procIDStr)

		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		data, err := simplejson.NewJson([]byte(value))
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		procData, err := data.Map()
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		procData[common.BKAppIDField] = appID
		valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDProc, cli.CC.ObjCtrl(), forward, defErr)
		_, err = valid.ValidMap(procData, common.ValidUpdate, procID)
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommFieldNotValid)
		}

		// take snapshot before operation
		preDetails, err := cli.getProcDetail(req, ownerID, appID, int(procID))
		if err != nil {
			blog.Errorf("get inst detail error: %v", err)
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrAuditSaveLogFaile)
		}

		input := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKAppIDField] = appID
		condition[common.BKProcIDField] = procID
		input["condition"] = condition
		input["data"] = procData
		procInfoJson, _ := json.Marshal(input)
		cProcURL := cli.CC.ObjCtrl() + "/object/v1/insts/process"
		blog.Info("update process url:%v", cProcURL)
		blog.Info("update process data:%s", string(procInfoJson))
		sProcRes, err := httpcli.ReqHttp(req, cProcURL, common.HTTPUpdate, []byte(procInfoJson))
		blog.Info("update process return:%s", string(sProcRes))
		if nil != err {
			blog.Error("update process error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcUpdateProcessFaile)
		}

		{
			// save change log
			headers := []metadata.Header{}
			// take snapshot before operation
			details, getErr := cli.getProcDetail(req, ownerID, appID, int(procID))
			if getErr != nil {
				blog.Errorf("get inst detail error: %v", getErr)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrAuditSaveLogFaile)
			}
			curData := map[string]interface{}{}
			for _, detail := range details {
				curData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
				headers = append(headers,
					metadata.Header{
						PropertyID:   fmt.Sprint(detail[common.BKPropertyIDField].(string)),
						PropertyName: fmt.Sprint(detail[common.BKPropertyNameField]),
					})
			}
			preData := map[string]interface{}{}
			for _, detail := range preDetails {
				preData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
			}
			auditContent := metadata.Content{
				CurData: curData,
				PreData: preData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header).AuditProcLog(procID, auditContent, "update process", ownerID, fmt.Sprint(appID), user, auditoplog.AuditOpTypeModify)
		}

		js, err := simplejson.NewJson([]byte(sProcRes))
		procResData, _ := js.Map()
		return http.StatusOK, procResData["data"], nil
	}, resp)
}

//SearchProcess search process
func (cli *procAction) SearchProcess(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID := pathParams[common.BKOwnerIDField]
		appIDStr := pathParams[common.BKAppIDField]
		appID, _ := strconv.Atoi(appIDStr)

		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		var js params.SearchParams
		err = json.Unmarshal([]byte(value), &js)
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		condition := js.Condition
		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKAppIDField] = appID
		if processName, ok := condition["bk_process_name"]; ok {
			processNameStr, ok := processName.(string)
			if ok {
				condition["bk_process_name"] = map[string]interface{}{common.BKDBLIKE: params.SpeceialCharChange(processNameStr)}

			} else {
				condition["bk_process_name"] = map[string]interface{}{common.BKDBLIKE: processName}

			}
		}
		page := js.Page
		searchParams := make(map[string]interface{})
		searchParams["condition"] = condition
		searchParams["fields"] = strings.Join(js.Fields, ",")
		searchParams["start"] = page["start"]
		searchParams["limit"] = page["limit"]
		searchParams["sort"] = page["sort"]

		// query process by module name
		if moduleName, ok := condition["bk_module_name"]; ok {

			procIdUrl := cli.CC.ProcCtrl() + "/process/v1/module/search/"

			qJ, _ := json.Marshal(map[string]interface{}{
				common.BKAppIDField:      appID,
				common.BKModuleNameField: moduleName,
			})

			blog.Info("get process arr by module: %s(%s)", procIdUrl, qJ)
			if res, err := httpcli.ReqHttp(req, procIdUrl, common.HTTPSelectPost, qJ); err == nil {

				var resMap ProcModuleResult

				err = json.Unmarshal([]byte(res), &resMap)
				if err != nil || !resMap.Result {
					blog.Error("get process arr failed: %v", res)
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcSearchProcessFaile)

				}

				// parse procid array
				procIdArr := make([]int, 0)
				for _, item := range resMap.Data {
					procIdArr = append(procIdArr, item.ProcessID)
				}

				// update condition
				condition[common.BKProcIDField] = map[string]interface{}{
					"$in": procIdArr,
				}
				delete(condition, common.BKModuleNameField)
				searchParams["condition"] = condition

			} else {
				blog.Error("get process arr err: %v, %v", res, err)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcSearchProcessFaile)
			}
		}

		// search process
		cProcURL := cli.CC.ObjCtrl() + "/object/v1/insts/process/search"
		procInfoJson, _ := json.Marshal(searchParams)
		blog.Info("search process url:%v, data: %s", cProcURL, string(procInfoJson))
		sProcRes, err := httpcli.ReqHttp(req, cProcURL, common.HTTPSelectPost, []byte(procInfoJson))
		blog.Info("search process return:%s", string(sProcRes))
		if nil != err {
			blog.Error("search process error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcSearchProcessFaile)
		}

		var sResult ProcessResult
		err = json.Unmarshal([]byte(sProcRes), &sResult)
		if nil != err {
			blog.Error("search process error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcSearchProcessFaile)
		}

		// query fail
		if !sResult.Result {
			return http.StatusOK, sResult.Data, defErr.Error(sResult.Code)
		}

		return http.StatusOK, sResult.Data, nil
	}, resp)
}

//DeleteProcess delete process
func (cli *procAction) DeleteProcess(req *restful.Request, resp *restful.Response) {
	user := util.GetActionUser(req)
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID := pathParams[common.BKOwnerIDField]
		appIDStr := pathParams[common.BKAppIDField]
		proIDStr := pathParams[common.BKProcIDField]
		appID, _ := strconv.Atoi(appIDStr)
		procID, _ := strconv.Atoi(proIDStr)

		// take snapshot before operation
		details, err := cli.getProcDetail(req, ownerID, appID, int(procID))
		if err != nil {
			blog.Errorf("get inst detail error: %v", err)
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrAuditSaveLogFaile)
		}

		conditon := make(map[string]interface{})
		conditon[common.BKAppIDField] = appID
		conditon[common.BKProcIDField] = procID
		conditon[common.BKOwnerIDField] = ownerID
		//delete process
		procInfoJson, _ := json.Marshal(conditon)
		dProcURL := cli.CC.ObjCtrl() + "/object/v1/insts/process"
		blog.Info("delete process url:%v", dProcURL)
		blog.Info("delete process data:%s", string(procInfoJson))
		cProcRes, err := httpcli.ReqHttp(req, dProcURL, common.HTTPDelete, []byte(procInfoJson))
		blog.Info("delete process return:%s", string(cProcRes))
		if nil != err {
			blog.Error("delete process error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcDeleteProcessFaile)
		}
		var info ProcessResult
		err = json.Unmarshal([]byte(cProcRes), &info)
		if nil != err {
			blog.Error("delete process error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcDeleteProcessFaile)
		}

		{
			// save change log
			headers := []metadata.Header{}

			preData := map[string]interface{}{}
			for _, detail := range details {
				preData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
				headers = append(headers,
					metadata.Header{
						PropertyID:   fmt.Sprint(detail[common.BKPropertyIDField].(string)),
						PropertyName: fmt.Sprint(detail[common.BKPropertyNameField]),
					})
			}
			auditContent := metadata.Content{
				PreData: preData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header).AuditProcLog(procID, auditContent, "delete process", ownerID, fmt.Sprint(appID), user, auditoplog.AuditOpTypeDel)
		}

		return http.StatusOK, nil, nil
	}, resp)
}

// CreateProcess create process
func (cli *procAction) CreateProcess(req *restful.Request, resp *restful.Response) {
	user := util.GetActionUser(req)
	language := util.GetActionLanguage(req)
	forward := &sourceAPI.ForwardParam{Header: req.Request.Header}
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		ownerID := pathParams[common.BKOwnerIDField]
		appIDStr := pathParams[common.BKAppIDField]
		appID, _ := strconv.Atoi(appIDStr)
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		input, err := js.Map()
		input[common.BKAppIDField] = appID

		valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDProc, cli.CC.ObjCtrl(), forward, defErr)
		_, err = valid.ValidMap(input, common.ValidCreate, 0)
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommFieldNotValid)
		}
		//create process
		input[common.BKOwnerIDField] = ownerID
		procInfoJson, _ := json.Marshal(input)
		cProcURL := cli.CC.ObjCtrl() + "/object/v1/insts/process"
		blog.Info("create process url:%v", cProcURL)
		blog.Info("create process data:%s", string(procInfoJson))
		cProcRes, err := httpcli.ReqHttp(req, cProcURL, common.HTTPCreate, []byte(procInfoJson))
		blog.Info("create process return:%s", string(cProcRes))
		if nil != err {
			blog.Error("create process error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcCreateProcessFaile)
		}
		var info ProcessCResult
		err = json.Unmarshal([]byte(cProcRes), &info)
		if nil != err {
			blog.Error("create process error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcCreateProcessFaile)
		}

		{
			// save change log
			instID := gjson.Get(cProcRes, "data."+common.BKProcIDField).Int()
			if instID == 0 {
				blog.Errorf("inst id not found")
			}
			headers := []metadata.Header{}

			curData := map[string]interface{}{}
			details, err := cli.getProcDetail(req, ownerID, appID, int(instID))
			if err != nil {
				blog.Errorf("get inst detail error: %v", err)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrAuditSaveLogFaile)
			}
			for _, detail := range details {
				curData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
				headers = append(headers,
					metadata.Header{
						PropertyID:   fmt.Sprint(detail[common.BKPropertyIDField].(string)),
						PropertyName: fmt.Sprint(detail[common.BKPropertyNameField]),
					})
			}
			auditContent := metadata.Content{
				CurData: curData,
				Headers: headers,
			}
			auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header).AuditProcLog(instID, auditContent, "create process", ownerID, fmt.Sprint(appID), user, auditoplog.AuditOpTypeAdd)
		}

		result := make(map[string]interface{})
		data := info.Data
		result[common.BKProcIDField] = data[common.BKProcIDField]

		return http.StatusOK, result, nil
	}, resp)
}

// BatchUpdateProcess batch update process
func (cli *procAction) BatchUpdateProcess(req *restful.Request, resp *restful.Response) {
	user := util.GetActionUser(req)
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &sourceAPI.ForwardParam{Header: req.Request.Header}

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID := pathParams[common.BKOwnerIDField]
		appIDStr := pathParams[common.BKAppIDField]
		appID, _ := strconv.Atoi(appIDStr)

		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		data, err := simplejson.NewJson([]byte(value))
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		procData, err := data.Map()
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		// split process str by comma
		var procIDArr []string
		if procIDStr, ok := procData[common.BKProcessIDField].(string); !ok {
			return http.StatusInternalServerError, nil, defErr.Error(common.CC_Err_Comm_PROC_Field_VALID_FAIL)
		} else {
			procIDArr = strings.Split(procIDStr, ",")
		}

		// update condition
		procData[common.BKAppIDField] = appID
		delete(procData, common.BKProcessIDField)

		// forbidden edit process name and func id
		delete(procData, common.BKProcessNameField)
		delete(procData, common.BKFuncIDField)

		// parse process id and valid
		var iProcIDArr []int
		auditContentArr := make([]metadata.Content, len(procIDArr))

		valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDProc, cli.CC.ObjCtrl(), forward, defErr)

		for index, pidStr := range procIDArr {

			procID, _ := strconv.Atoi(pidStr)

			_, err = valid.ValidMap(procData, common.ValidUpdate, procID)
			if nil != err {
				blog.Error("valid procID = %v err: %v", err)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommFieldNotValid)
			}

			iProcIDArr = append(iProcIDArr, procID)

			// take snapshot before operation
			details, err := cli.getProcDetail(req, ownerID, appID, int(procID))
			if err != nil {
				blog.Errorf("get inst detail error: %v", err)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrAuditSaveLogFaile)
			}

			// save change log
			headers := []metadata.Header{}

			curData := map[string]interface{}{}
			for _, detail := range details {
				curData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
				headers = append(headers,
					metadata.Header{
						PropertyID:   fmt.Sprint(detail[common.BKPropertyIDField].(string)),
						PropertyName: fmt.Sprint(detail[common.BKPropertyNameField]),
					})
			}

			curData["bk_process_id"] = procID

			// save proc info before modify
			auditContentArr[index] = metadata.Content{
				CurData: make(map[string]interface{}),
				PreData: curData,
				Headers: headers,
			}

		}

		// update processes
		input := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition[common.BKOwnerIDField] = ownerID
		condition[common.BKAppIDField] = appID
		condition[common.BKProcIDField] = map[string]interface{}{
			common.BKDBIN: iProcIDArr,
		}

		input["condition"] = condition
		input["data"] = procData

		procInfoJson, _ := json.Marshal(input)
		cProcURL := cli.CC.ObjCtrl() + "/object/v1/insts/process"

		blog.Info("update process url:%v", cProcURL)
		blog.Info("update process data:%s", string(procInfoJson))
		sProcRes, err := httpcli.ReqHttp(req, cProcURL, common.HTTPUpdate, []byte(procInfoJson))
		blog.Info("update process return:%s", string(sProcRes))
		if nil != err {
			blog.Error("update process error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcUpdateProcessFaile)
		}

		// fetch and return data
		js, err := simplejson.NewJson([]byte(sProcRes))
		procResData, _ := js.Map()

		// update audit log content after modify
		for index, procID := range iProcIDArr {

			// take snapshot before operation
			details, err := cli.getProcDetail(req, ownerID, appID, int(procID))

			if err != nil {
				blog.Errorf("get inst detail error: %v", err)
				return http.StatusInternalServerError, "", defErr.Error(common.CCErrAuditSaveLogFaile)
			}

			// save change log
			curData := map[string]interface{}{}
			for _, detail := range details {
				curData[detail[common.BKPropertyIDField].(string)] = detail[common.BKPropertyValueField]
			}

			curData["bk_process_id"] = procID

			// save proc info before modify
			auditContentArr[index].CurData = curData

			auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header).AuditProcLog(procID, auditContentArr[index], "update process", ownerID, fmt.Sprint(appID), user, auditoplog.AuditOpTypeModify)

			blog.Info("update process with: %s", auditContentArr[index])
		}

		return http.StatusOK, procResData["data"], nil

	}, resp)
}
