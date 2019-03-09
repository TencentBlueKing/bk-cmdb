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
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2/bson"
	//"time"
)

func (cli *Service) GetProcessesByModuleName(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Errorf("read request body failed, error:%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Errorf("unmarshal json error:%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	moduleName := input[common.BKModuleNameField]
	query := make(map[string]interface{})
	query[common.BKModuleNameField] = moduleName

	fields := []string{common.BKProcIDField, common.BKAppIDField, common.BKModuleNameField}
	var result []interface{}
	query = util.SetModOwner(query, ownerID)
	err = db.Table(common.BKTableNameProcModule).Find(query).Limit(common.BKNoLimit).Sort(common.BKHostIDField).Fields(fields...).All(ctx, &result)
	if err != nil {
		blog.Errorf("fail to get module proc config %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
		return
	}
	processIdArr := make([]int, 0)
	for _, item := range result {
		itemMap, ok := item.(bson.M)
		if false == ok {
			blog.Errorf("assign error item is not bson.M,item:%v", item)
			continue
		}
		processID, err := util.GetIntByInterface(itemMap[common.BKProcIDField])
		if nil != err {
			blog.Errorf("GetIntByInterface err:%v", err)
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
	err = db.Table(common.BKTableNameBaseProcess).Find(procQuery).Sort(common.BKProcIDField).Limit(common.BKNoLimit).All(ctx, &resultProc)

	if err != nil {
		blog.Errorf("fail to get proc %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: resultProc})
}
