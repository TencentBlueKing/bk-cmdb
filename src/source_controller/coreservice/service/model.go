/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) CreateManyModelClassification(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputDatas := metadata.CreateManyModelClassifiaction{}
	if err := data.MarshalJSONInto(&inputDatas); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().CreateManyModelClassification(params, inputDatas)
}

func (s *coreService) CreateOneModelClassification(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.CreateOneModelClassification{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().CreateOneModelClassification(params, inputData)
}

func (s *coreService) SetOneModelClassificaition(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.SetOneModelClassification{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().SetOneModelClassification(params, inputData)
}

func (s *coreService) SetManyModelClassificaiton(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputDatas := metadata.SetManyModelClassification{}
	if err := data.MarshalJSONInto(&inputDatas); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().SetManyModelClassification(params, inputDatas)
}

func (s *coreService) UpdateModelClassification(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.UpdateOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().UpdateModelClassification(params, inputData)
}

func (s *coreService) DeleteModelClassification(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().DeleteModelClassificaiton(params, inputData)
}

func (s *coreService) CascadeDeleteModelClassification(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().CascadeDeleteModeClassification(params, inputData)
}

func (s *coreService) SearchModelClassification(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().SearchModelClassification(params, inputData)
}

func (s *coreService) CreateModel(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.CreateModel{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().CreateModel(params, inputData)
}

func (s *coreService) SetModel(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.SetModel{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().SetModel(params, inputData)
}

func (s *coreService) UpdateModel(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.UpdateOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().UpdateModel(params, inputData)
}

func (s *coreService) DeleteModel(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().DeleteModel(params, inputData)
}

func (s *coreService) CascadeDeleteModel(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().CascadeDeleteModel(params, inputData)
}

func (s *coreService) SearchModel(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().SearchModel(params, inputData)
}

func (s *coreService) CreateModelAttributeGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.CreateModelAttributeGroup{}
	if err := data.ToStructByTag(inputData.Data, "field"); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().CreateModelAttributeGroup(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) SetModelAttributeGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.SetModelAttributeGroup{}
	if err := data.ToStructByTag(inputData.Data, "field"); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().SetModelAttributeGroup(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) UpdateModelAttributeGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.UpdateOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().UpdateModelAttributeGroup(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) SearchModelAttributeGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().SearchModelAttributeGroup(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) DeleteModelAttributeGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().DeleteModelAttributeGroup(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) CreateModelAttributes(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.CreateModelAttributes{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().CreateModelAttributes(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) SetModelAttributes(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.SetModelAttributes{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().SetModelAttributes(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) UpdateModelAttributes(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.UpdateOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().UpdateModelAttributes(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) DeleteModelAttribute(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().DeleteModelAttributes(params, pathParams("bk_obj_id"), inputData)
}
func (s *coreService) SearchModelAttributes(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().SearchModelAttributes(params, pathParams("bk_obj_id"), inputData)
}
