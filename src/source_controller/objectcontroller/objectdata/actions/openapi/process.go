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
		if err != nil {
			blog.Error("fail to get proc %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}
		/*resData := make([]map[string]interface{}, 0)
		for _, item := range resultProc {
			itemMap, ok := item.(bson.M)
			if false == ok {
				blog.Error("assign error item is not bson.M,item:%v",item)
				continue
			}
			lastTime ,ok := itemMap[common.LastTimeField].(time.Time)
			LastTime := ""
			if false == ok {
				blog.Error("assign error item is itemMap['last_time'] not (time.Time),itemMap:%v",itemMap)
			}else{
				LastTime = lastTime.Format("2006-01-02 15:04:05")
			}
			createTime ,ok := itemMap[common.CreateTimeField].(time.Time)
			CreateTime := ""
			if false == ok {
				blog.Error("assign error item is itemMap['create_time'] not (time.Time),itemMap:%v",itemMap)
			}else{
				CreateTime = createTime.Format("2006-01-02 15:04:05")
			}

			resData = append(resData, map[string]interface{}{
				"WorkPath":    itemMap[common.BKWorkPath],
				"AutoTimeGap": "0",
				"LastTime":    LastTime,
				"StartCmd":    "",
				"FuncID":      "0",
				"BindIP":      itemMap[common.BKBindIP],
				"FuncName":    itemMap[common.BKFuncName],
				//"Flag":        itemMap["Flag"],
				"Flag":        "",
				"User":        itemMap[	common.BKUser],
				"StopCmd":     "",
				"ProNum":      "",
				"ReloadCmd":   "",
				"ProcessName": itemMap[common.BKProcNameField],
				"OpTimeout":   "0",
				"KillCmd":     "",
				"Protocol":    itemMap[common.BKProtocol],
				"Seq":         "0",
				"ProcGrp":     "",
				"Port":        itemMap[common.BKPort],
				"ReStartCmd":  "",
				"AutoStart":   "0",
				"CreateTime":  CreateTime,
				"PidFile":     "",
			})
		}*/
		return http.StatusOK, resultProc, nil
	}, resp)
}
