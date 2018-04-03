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
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"net/http"
	"strconv"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/{" + common.BKOwnerIDField + "}/{" + common.BKAppIDField + "}/{" + common.BKProcIDField + "}", Params: nil, Handler: process.GetProcessDetailByID})
	process.CreateAction()
}

//GetProcessDetailByID  get process detail info
func (cli *procAction) GetProcessDetailByID(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID := pathParams[common.BKOwnerIDField]
		appIDStr := pathParams[common.BKAppIDField]
		appID, _ := strconv.Atoi(appIDStr)
		procIDStr := pathParams[common.BKProcIDField]
		procID, _ := strconv.Atoi(procIDStr)

		reResult, err := cli.getProcDetail(req, ownerID, appID, procID)
		if err != nil {
			blog.Error("GetProcessDetailByID info error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrProcSearchDetailFaile)
		}
		return http.StatusOK, reResult, nil
	}, resp)
}

func (cli *procAction) getProcDetail(req *restful.Request, ownerID string, appID, procID int) ([]map[string]interface{}, error) {
	condition := make(map[string]interface{})
	condition[common.BKOwnerIDField] = ownerID
	condition[common.BKAppIDField] = appID
	condition[common.BKProcIDField] = procID
	searchParams := make(map[string]interface{})
	searchParams["condition"] = condition
	procInfoJson, _ := json.Marshal(searchParams)
	gProcURL := cli.CC.ObjCtrl() + "/object/v1/insts/process/search"
	blog.Info("get obj detail query url: %v", gProcURL)
	blog.Info("get obj detail query params: %v", string(procInfoJson))
	gProcRe, err := httpcli.ReqHttp(req, gProcURL, common.HTTPSelectPost, []byte(procInfoJson))
	blog.Info("get obj detail query return: %v", gProcRe)
	if nil != err {
		blog.Error("GetProcessDetailByID info error :%v", err)
		return nil, err
	}
	var getResult ProcessResult
	err = json.Unmarshal([]byte(gProcRe), &getResult)
	if nil != err {
		blog.Error("GetProcessDetailByID info error :%v", err)
		return nil, err
	}
	if false == getResult.Result {
		blog.Error("GetProcDetailByID  info error :%v", err)
		return nil, err
	}
	procData, ok := getResult.Data.(map[string]interface{})
	if false == ok {
		blog.Error("GetProcDetailByID  info error :%v", err)
		return nil, err
	}
	procInfoArr := procData["info"].([]interface{})
	proc := make(map[string]interface{}, 0)
	for _, j := range procInfoArr {
		procCell := j.(map[string]interface{})
		for k, v := range procCell {
			proc[k] = v
		}
	}
	gProcAttrURL := cli.CC.ObjCtrl() + "/object/v1/meta/objectatts"
	searchBody := make(map[string]interface{})
	searchBody[common.BKObjIDField] = common.BKInnerObjIDProc
	searchBody[common.BKOwnerIDField] = ownerID
	searchJson, _ := json.Marshal(searchBody)
	gProcAttrRe, err := httpcli.ReqHttp(req, gProcAttrURL, common.HTTPSelectPost, []byte(searchJson))
	if nil != err {
		blog.Error("GetProcDetailByID  attr error :%v", err)
		return nil, err
	}

	js, err := simplejson.NewJson([]byte(gProcAttrRe))
	gProcAttr, _ := js.Map()
	gAttrResult := gProcAttr["result"].(bool)
	if false == gAttrResult {
		blog.Error("GetProcDetailByID  attr error :%v", err)
		return nil, err
	}
	procAttrArr := gProcAttr["data"].([]interface{})
	reResult := make([]map[string]interface{}, 0)
	for _, i := range procAttrArr {
		attr := i.(map[string]interface{})
		data := make(map[string]interface{})
		propertyID := attr[common.BKPropertyIDField].(string)
		if propertyID == common.BKChildStr {
			continue
		}
		data[common.BKPropertyIDField] = propertyID
		data[common.BKPropertyNameField] = attr[common.BKPropertyNameField]

		data[common.BKPropertyValueField] = proc[propertyID]
		reResult = append(reResult, data)
	}
	return reResult, nil
}
