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

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateClassification create object's classification
func (cli *Service) SearchTopoGraphics(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	blog.Infof("search param %s", value)

	selector := meta.TopoGraphics{}
	if jsErr := json.Unmarshal(value, &selector); nil != jsErr {
		blog.Error("failed to unmarshal the data, data is %s, error info is %s ", value, jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	selector.SetSupplierAccount(ownerID)
	results := []meta.TopoGraphics{}
	if selErr := db.Table(common.BKTableNameTopoGraphics).Find(selector).All(ctx, &results); nil != selErr {
		blog.Error("select data failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})
}

func (cli *Service) UpdateTopoGraphics(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// execute
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	datas := []meta.TopoGraphics{}
	if jsErr := json.Unmarshal(value, &datas); nil != jsErr {
		blog.Error("failed to unmarshal the data, data is %s, error info is %s ", value, jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	for index := range datas {
		datas[index].SetSupplierAccount(ownerID)
		err = db.Table(common.BKTableNameTopoGraphics).Insert(ctx, datas[index].FillBlank())
		if cli.Instance.IsDuplicatedError(err) {
			condition := meta.TopoGraphics{}
			condition.SetScopeType(*datas[index].ScopeType)
			condition.SetScopeID(*datas[index].ScopeID)
			condition.SetNodeType(*datas[index].NodeType)
			condition.SetObjID(*datas[index].ObjID)
			condition.SetInstID(*datas[index].InstID)
			condition.SetSupplierAccount(ownerID)
			if err = cli.Instance.Table(common.BKTableNameTopoGraphics).Update(context.Background(), condition, datas[index]); err != nil {
				blog.Error("update data failed, error information is %s", err.Error())
				resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBUpdateFailed, err.Error())})
				return
			}
		} else if err != nil {
			blog.Error("insert data failed, error information is %s", err.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
			return
		}
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}
