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

package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2/bson"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	//"time"
)

func (cli *Service) GetProcessesByModuleName(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	blog.Debug("GetProcessesByModuleName start !")
	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("unmarshal json error:%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	moduleName := input[common.BKModuleNameField]
	query := make(map[string]interface{})
	query[common.BKModuleNameField] = moduleName
	blog.Debug("query;%v", query)
	fields := []string{common.BKProcIDField, common.BKAppIDField, common.BKModuleNameField}
	var result []interface{}
	query = util.SetModOwner(query, ownerID)
	err = cli.Instance.GetMutilByCondition("cc_Proc2Module", fields, query, &result, common.BKHostIDField, 0, 100000)
	if err != nil && !cli.Instance.IsNotFoundErr(err) {
		blog.Error("fail to get module proc config %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
		return
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
	procQuery = util.SetModOwner(procQuery, ownerID)
	var resultProc []interface{}
	err = cli.Instance.GetMutilByCondition("cc_Process", []string{}, procQuery, &resultProc, common.BKProcIDField, 0, 100000)
	blog.Infof("GetProcessesByModuleName params:%v, result:%v", procQuery, resultProc)
	if err != nil && !cli.Instance.IsNotFoundErr(err) {
		blog.Error("fail to get proc %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: resultProc})
}
