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
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
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
		blog.Errorf("search topo graphics, but read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	selector := meta.TopoGraphics{}
	if jsErr := json.Unmarshal(value, &selector); nil != jsErr {
		blog.Errorf("search topo graphics, but failed to unmarshal the data, data is %s, error info is %s ", value, jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	cond := mapstr.MapStr{
		"scope_type":          selector.ScopeType,
		"scope_id":            selector.ScopeID,
		"bk_supplier_account": ownerID,
	}
	_, err = selector.Metadata.Label.GetBusinessID()
	if nil == err {
		cond.Merge(meta.PublicAndBizCondition(selector.Metadata))
	} else {
		cond.Merge(meta.BizLabelNotExist)
	}

	results := []meta.TopoGraphics{}
	if selErr := db.Table(common.BKTableNameTopoGraphics).Find(cond).All(ctx, &results); nil != selErr {
		blog.Errorf("search topo graphics, but select data failed, error information is %s", selErr.Error())
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
		blog.Errorf("update topo graphics, but read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	datas := []meta.TopoGraphics{}
	if jsErr := json.Unmarshal(value, &datas); nil != jsErr {
		blog.Errorf("update topo graphics, but failed to unmarshal the data, data is %s, error info is %s ", value, jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	for index := range datas {
		datas[index].SetSupplierAccount(ownerID)
		cond := mapstr.MapStr{
			"scope_type":          datas[index].ScopeType,
			"scope_id":            datas[index].ScopeID,
			"node_type":           datas[index].NodeType,
			"bk_obj_id":           datas[index].ObjID,
			"bk_inst_id":          datas[index].InstID,
			"bk_supplier_account": ownerID,
		}
		_, err := datas[index].Metadata.Label.GetBusinessID()
		if nil != err {
			cond.Merge(meta.BizLabelNotExist)
		} else {
			cond.Set("metadata", datas[index].Metadata)
		}
		cnt, err := db.Table(common.BKTableNameTopoGraphics).Find(cond).Count(ctx)
		if nil != err {
			blog.Errorf("update topo graphics, search data error: %s", value, err.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
			return
		}
		if 0 == cnt {
			err = db.Table(common.BKTableNameTopoGraphics).Insert(ctx, datas[index])
			if nil != err {
				blog.Errorf("update topo graphics, but insert data failed, err:%s", err.Error())
				resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
				return
			}
		} else {
			if err = cli.Instance.Table(common.BKTableNameTopoGraphics).Update(context.Background(), cond, datas[index]); err != nil {
				blog.Errorf("update topo graphics, but update failed, err: %s", err.Error())
				resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBUpdateFailed, err.Error())})
				return
			}
		}

	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}
