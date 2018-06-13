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
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"
)

var topo = &topoAction{}

// topoAction
type topoAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/topo/internal/{owner_id}/{app_id}", Params: nil, Handler: topo.GetInternalModule})
	topo.CreateAction()
}

//get built in module
func (cli *topoAction) GetInternalModule(req *restful.Request, resp *restful.Response) {

	// get language
	language := util.GetActionLanguage(req)

	// get the error by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		appIDStr := req.PathParameter("app_id")
		appID, _ := strconv.Atoi(appIDStr)
		setCond := make(map[string]interface{})
		cond := make(map[string]interface{})
		cond[common.BKAppIDField] = appID
		cond[common.BKDefaultField] = common.DefaultResModuleFlag
		setCond["condition"] = cond

		//search set
		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/set/search"
		inputJSON, _ := json.Marshal(setCond)
		setRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
		blog.Info("search set params: %s", string(inputJSON))
		blog.Info("search set return: %s", string(setRes))
		if nil != err {
			blog.Error("search set error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoSetSelectFailed)
		}

		moduleCond := make(map[string]interface{})
		defaultCond := make(map[string]interface{})
		defaultCond[common.BKDBIN] = []int{common.DefaultResModuleFlag, common.DefaultFaultModuleFlag}
		cond[common.BKDefaultField] = defaultCond
		moduleCond["condition"] = cond

		//search module
		sURL = cli.CC.ObjCtrl() + "/object/v1/insts/module/search"
		inputJSON, _ = json.Marshal(moduleCond)
		moduleRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
		blog.Debug("search module params: %s", string(inputJSON))
		blog.Info("search module return: %s", string(moduleRes))
		if nil != err {
			blog.Error("search module error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleSelectFailed)
		}
		var setObj api.APIRsp
		var moduleObj api.APIRsp
		err = json.Unmarshal([]byte(setRes), &setObj)
		if nil != err || !setObj.Result {
			blog.Error("search set error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoSetSelectFailed)
		}
		err = json.Unmarshal([]byte(moduleRes), &moduleObj)
		if nil != err || !moduleObj.Result {
			blog.Error("search set error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleSelectFailed)
		}
		setData := setObj.Data.(map[string]interface{})
		setInfo := setData["info"].([]interface{})
		moduleData := moduleObj.Data.(map[string]interface{})
		moduleInfo := moduleData["info"].([]interface{})
		if 0 == len(setInfo) || 0 == len(moduleInfo) {
			blog.Error("search set error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleSelectFailed)
		}

		setResult := make(map[string]interface{})
		moduleResult := make([]map[string]interface{}, 0)
		for _, i := range setInfo {
			set := i.(map[string]interface{})
			setResult[common.BKSetIDField] = set[common.BKSetIDField]
			setResult[common.BKSetNameField] = set[common.BKSetNameField]
		}
		for _, j := range moduleInfo {
			moduleR := make(map[string]interface{})
			module := j.(map[string]interface{})
			moduleR[common.BKModuleIDField] = module[common.BKModuleIDField]
			moduleR[common.BKModuleNameField] = module[common.BKModuleNameField]
			moduleResult = append(moduleResult, moduleR)
		}
		setResult[common.BKInnerObjIDModule] = moduleResult

		return http.StatusOK, setResult, nil
	}, resp)

}
