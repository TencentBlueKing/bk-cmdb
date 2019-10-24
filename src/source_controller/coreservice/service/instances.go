/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) CreateOneModelInstance(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := metadata.CreateModelInstance{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.InstanceOperation().CreateModelInstance(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) CreateManyModelInstances(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := metadata.CreateManyModelInstance{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.InstanceOperation().CreateManyModelInstance(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) UpdateModelInstances(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := metadata.UpdateOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	// TODO: remove this logic when biz model is changed.
	cond := metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.AssociationKindIDField:  common.AssociationKindMainline,
			common.AssociatedObjectIDField: pathParams("bk_obj_id"),
		},
	}
	result, err := s.core.AssociationOperation().SearchModelAssociation(params, cond)
	if err != nil {
		return nil, err
	}

	if len(result.Info) != 0 {
		// this is a mainline object, need to delete metadata field.
		// otherwise, can not find this object, then update failed.
		inputData.Condition.Remove("metadata")
	}

	return s.core.InstanceOperation().UpdateModelInstance(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) SearchModelInstances(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	// 判断是否有要根据default字段，需要国际化的内容
	if _, ok := defaultNameLanguagePkg[pathParams("bk_obj_id")]; ok {
		// 大于两个字段
		if len(inputData.Fields) > 1 {
			inputData.Fields = append(inputData.Fields, common.BKDefaultField)
		} else if len(inputData.Fields) == 1 && inputData.Fields[0] != "" {
			// 只有一个字段，如果字段为空白字符，则不处理
			inputData.Fields = append(inputData.Fields, common.BKDefaultField)
		}
	}

	dataResult, err := s.core.InstanceOperation().SearchModelInstance(params, pathParams("bk_obj_id"), inputData)
	if nil != err {
		return dataResult, err
	}

	// translate language for default name
	if m, ok := defaultNameLanguagePkg[pathParams("bk_obj_id")]; ok {
		for idx := range dataResult.Info {
			subResult := m[fmt.Sprint(dataResult.Info[idx][common.BKDefaultField])]
			if len(subResult) >= 3 {
				dataResult.Info[idx][subResult[1]] = util.FirstNotEmptyString(params.Lang.Language(subResult[0]), fmt.Sprint(dataResult.Info[idx][subResult[1]]), fmt.Sprint(dataResult.Info[idx][subResult[2]]))
			}
		}

	}
	return dataResult, err
}

func (s *coreService) DeleteModelInstances(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.InstanceOperation().DeleteModelInstance(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) CascadeDeleteModelInstances(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.InstanceOperation().CascadeDeleteModelInstance(params, pathParams("bk_obj_id"), inputData)
}
