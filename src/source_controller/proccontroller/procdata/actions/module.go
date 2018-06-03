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
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	eventtypes "configcenter/src/scene_server/event_server/types"
	"configcenter/src/source_controller/common/eventdata"
	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2/bson"
)

var proc *proc2moduleAction = &proc2moduleAction{}

// ProcAction
type proc2moduleAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/module", Params: nil, Handler: proc.CreateProc2Module})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/module/search", Params: nil, Handler: proc.GetProc2Module})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/module/process", Params: nil, Handler: proc.GetModule2Proc})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/module", Params: nil, Handler: proc.DeleteProc2Module})
	proc.CreateAction()
}

//DeleteProc2Module delete proc module config
func (cli *proc2moduleAction) DeleteProc2Module(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("DeleteProc2Module read json fail", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			blog.Error("DeleteProc2Module json decode fail", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input, err := js.Map()
		if err != nil {
			blog.Error("DeleteProc2Module json not array", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommParamsInvalid)
		}

		// retrieve original data
		var originals []interface{}
		err = proc.CC.InstCli.GetMutilByCondition(common.BKTableNameProcModule, []string{}, input, &originals, "", 0, 0)
		if err != nil {
			blog.Error("retrieve original error:%v", err)
		}

		blog.Info("delete proc module config %v", input)
		err = proc.CC.InstCli.DelByCondition(common.BKTableNameProcModule, input)
		if err != nil {
			blog.Error("delete proc module config error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcDeleteProc2Module)
		}

		// send event
		if len(originals) > 0 {
			ec := eventdata.NewEventContextByReq(req)
			for _, i := range originals {
				err := ec.InsertEvent(eventtypes.EventTypeRelation, "processmodule", eventtypes.EventActionDelete, nil, i)
				if err != nil {
					blog.Error("create event error:%s", err.Error())
				}
			}
		}
		return http.StatusOK, nil, nil
	}, resp)
}

//CreateProc2Module create proc module config
func (cli *proc2moduleAction) CreateProc2Module(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("CreateProc2Module read json fail", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			blog.Error("CreateProc2Module json decode fail", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input, err := js.Array()
		if err != nil {
			blog.Error("CreateProc2Module json not array", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommParamsInvalid)
		}

		blog.Info("create proc module config ", input)
		ec := eventdata.NewEventContextByReq(req)
		for _, i := range input {
			_, err = proc.CC.InstCli.Insert(common.BKTableNameProcModule, i)
			if err != nil {
				blog.Error("create proc module config error:%v", err)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcCreateProc2Module)
			}
			//  record events
			err := ec.InsertEvent(eventtypes.EventTypeRelation, "processmodule", eventtypes.EventActionCreate, i, nil)
			if err != nil {
				blog.Error("create event error:%v", err)
			}
		}

		return http.StatusOK, nil, nil
	}, resp)
}

//GetProc2Module get proc module config
func (cli *proc2moduleAction) GetProc2Module(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("GetProc2Module fail", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			blog.Error("GetProc2Module fail", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input, err := js.Map()
		if err != nil {
			blog.Error("GetProc2Module fail", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommParamsInvalid)
		}

		blog.Info("get proc module config condition ", input)
		var result []interface{}
		err = proc.CC.InstCli.GetMutilByCondition(common.BKTableNameProcModule, []string{}, input, &result, "", 0, 0)
		if err != nil {
			blog.Error("create proc module config error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcSelectProc2Module)
		}
		return http.StatusOK, result, nil
	}, resp)
}

// GetModule2Proc get process's id which bind to a module name
func (cli *proc2moduleAction) GetModule2Proc(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)

	// get the error factory by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		// read param
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("GetModule2Proc read param err: %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// parse json param
		js, err := simplejson.NewJson(value)
		if err != nil {
			blog.Error("GetModule2Proc parse json err: %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		// map json data
		input, err := js.Map()
		if err != nil {
			blog.Error("GetModule2Proc map json err: %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommParamsInvalid)
		}

		// get module name
		inputCond := input["condition"]
		condition, ok := inputCond.(map[string]interface{})
		if !ok {
			blog.Error("GetModule2Proc get module name failed")
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommParamsInvalid)
		}

		// get module bind procid list
		var records []interface{}

		query := make(map[string]interface{})
		query[common.BKAppIDField] = condition[common.BKAppIDField]
		query[common.BKModuleNameField] = condition[common.BKModuleNameField]
		fields := []string{common.BKProcIDField}
		blog.Info("GetModule2Proc: query=%v, fields=%v", query, fields)

		err = cli.CC.InstCli.GetMutilByCondition(common.BKTableNameProcModule, fields, query, &records, "", 0, 0)
		if err != nil {
			blog.Error("GetModule2Proc query err: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}

		// parse procid array
		procIdArr := make([]int, 0)
		for _, item := range records {
			itemMap, ok := item.(bson.M)
			if !ok {
				blog.Error("GetModule2Proc item %v is not bson.M", item)
				continue
			}
			procID, err := util.GetIntByInterface(itemMap[common.BKProcessIDField])
			if err != nil {
				blog.Error("GetModule2Proc convert id err： %v", err)
				continue
			}
			procIdArr = append(procIdArr, procID)
		}

		return http.StatusOK, procIdArr, nil
	}, resp)

}






