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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

// CreateClassification create object's classification
func (s *coreService) SearchTopoGraphics(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	selector := meta.TopoGraphics{}
	if jsErr := data.MarshalJSONInto(&selector); nil != jsErr {
		blog.Errorf("search topo graphics, but failed to unmarshal the data, data: %+v, err: %s, rid: %s", data, jsErr.Error(), params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	cond := mapstr.MapStr{
		"scope_type":          selector.ScopeType,
		"scope_id":            selector.ScopeID,
		"bk_supplier_account": params.SupplierAccount,
	}
	_, err := selector.Metadata.Label.GetBusinessID()
	if nil == err {
		cond.Merge(meta.PublicAndBizCondition(selector.Metadata))
	} else {
		cond.Merge(meta.BizLabelNotExist)
	}

	results := make([]meta.TopoGraphics, 0)
	if selErr := s.db.Table(common.BKTableNameTopoGraphics).Find(cond).All(params.Context, &results); nil != selErr {
		blog.Errorf("search topo graphics, but select data failed, error information is %s, rid: %s", selErr.Error(), params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	return results, nil
}

func (s *coreService) UpdateTopoGraphics(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputBody := struct {
		Data []meta.TopoGraphics `json:"data" field:"data" bson:"data"`
	}{}
	if jsErr := data.MarshalJSONInto(&inputBody); nil != jsErr {
		blog.Errorf("update topo graphics, but failed to unmarshal the data, data: %+v, err: %s, rid: %s", data, jsErr.Error(), params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	for index := range inputBody.Data {
		inputBody.Data[index].SetSupplierAccount(params.SupplierAccount)
		cond := mapstr.MapStr{
			"scope_type":          inputBody.Data[index].ScopeType,
			"scope_id":            inputBody.Data[index].ScopeID,
			"node_type":           inputBody.Data[index].NodeType,
			"bk_obj_id":           inputBody.Data[index].ObjID,
			"bk_inst_id":          inputBody.Data[index].InstID,
			"bk_supplier_account": params.SupplierAccount,
		}

		cnt, err := s.db.Table(common.BKTableNameTopoGraphics).Find(cond).Count(params.Context)
		if nil != err {
			blog.Errorf("update topo graphics, search data failed, data: %+v, err: %s, rid: %s", data, err.Error(), params.ReqID)
			return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
		}
		if 0 == cnt {
			err = s.db.Table(common.BKTableNameTopoGraphics).Insert(params.Context, inputBody.Data[index])
			if nil != err {
				blog.Errorf("update topo graphics, but insert data failed, err:%s, rid: %s", err.Error(), params.ReqID)
				return nil, params.Error.CCError(common.CCErrCommDBInsertFailed)
			}
		} else {
			if err = s.db.Table(common.BKTableNameTopoGraphics).Update(context.Background(), cond, inputBody.Data[index]); err != nil {
				blog.Errorf("update topo graphics, but update failed, err: %s, rid: %s", err.Error(), params.ReqID)
				return nil, params.Error.CCError(common.CCErrCommDBUpdateFailed)
			}
		}
	}

	return nil, nil
}
