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

package openapi

import (
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	//"time"
)

var proc *procAction = &procAction{}

type procAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/openapi/proc/getProcModule", Params: nil, Handler: proc.GetProcessesByModuleName, FilterHandler: nil})

	// create CC object
	proc.CreateAction()
}

func (cli *procAction) GetProcessesByModuleName(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		blog.Debug("GetProcessesByModuleName start !")
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		input := make(map[string]interface{})
		err = json.Unmarshal(value, &input)
		if nil != err {
			blog.Error("unmarshal json error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		moduleName := input[common.BKModuleNameField]
		query := make(map[string]interface{})
		query[common.BKModuleNameField] = moduleName
		query = util.SetModOwner(query, ownerID)
		blog.Debug("query;%v", query)
		fields := []string{common.BKProcIDField, common.BKAppIDField, common.BKModuleNameField}
		var result []interface{}
		err = proc.CC.InstCli.GetMutilByCondition("cc_Proc2Module", fields, query, &result, common.BKHostIDField, 0, 100000)
		if err != nil {
			blog.Error("fail to get module proc config %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}
		processIdArr := make([]int, 0)
		for _, item := range result {
			itemMap, ok := item.(bson.M)
			if false == ok {
				blog.Error("assign error item is not bson.M,item:%v", item)
				continue
			}
			processID, err := util.GetIntByInterface(itemMap[common.BKProcIDField])
			if nil != err {
				blog.Error("GetIntByInterface err:%v", err)
				continue
			}
			processIdArr = append(processIdArr, processID)
		}
		procQuery := make(map[string]interface{})
		procQuery[common.BKProcIDField] = map[string]interface{}{
			"$in": processIdArr,
		}
		var resultProc []interface{}
		err = proc.CC.InstCli.GetMutilByCondition("cc_Process", []string{}, procQuery, &resultProc, common.BKProcIDField, 0, 100000)
		blog.Infof("GetProcessesByModuleName params:%v, result:%v", procQuery, resultProc)
		if err != nil {
			blog.Error("fail to get proc %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}
		return http.StatusOK, resultProc, nil
	}, resp)
}
