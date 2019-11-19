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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
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

func (s *coreService) SetOneModelClassification(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.SetOneModelClassification{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().SetOneModelClassification(params, inputData)
}

func (s *coreService) SetManyModelClassification(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

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
	return s.core.ModelOperation().DeleteModelClassification(params, inputData)
}

func (s *coreService) SearchModelClassification(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	dataResult, err := s.core.ModelOperation().SearchModelClassification(params, inputData)
	if nil != err {
		return dataResult, err
	}

	// translate language
	for index := range dataResult.Info {
		dataResult.Info[index].ClassificationName = s.TranslateClassificationName(params.Lang, &dataResult.Info[index])
	}

	return dataResult, err
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

	idStr := pathParams(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID)
	}
	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().CascadeDeleteModel(params, id)
}

func (s *coreService) SearchModel(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	dataResult, err := s.core.ModelOperation().SearchModelWithAttribute(params, inputData)
	if nil != err {
		return dataResult, err
	}

	// translate
	for modelIdx := range dataResult.Info {
		dataResult.Info[modelIdx].Spec.ObjectName = s.TranslateObjectName(params.Lang, &dataResult.Info[modelIdx].Spec)

		for attributeIdx := range dataResult.Info[modelIdx].Attributes {
			dataResult.Info[modelIdx].Attributes[attributeIdx].PropertyName = s.TranslatePropertyName(params.Lang, &dataResult.Info[modelIdx].Attributes[attributeIdx])
			dataResult.Info[modelIdx].Attributes[attributeIdx].Placeholder = s.TranslatePlaceholder(params.Lang, &dataResult.Info[modelIdx].Attributes[attributeIdx])
			if dataResult.Info[modelIdx].Attributes[attributeIdx].PropertyType == common.FieldTypeEnum {
				dataResult.Info[modelIdx].Attributes[attributeIdx].Option = s.TranslateEnumName(params.Context, params.Lang, &dataResult.Info[modelIdx].Attributes[attributeIdx], dataResult.Info[modelIdx].Attributes[attributeIdx].Option)
			}
		}
	}

	return dataResult, err
}

// GetModelStatistics 用于统计各个模型的实例数(Web页面展示需要)
func (s *coreService) GetModelStatistics(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	filter := map[string]interface{}{}
	setCount, err := s.db.Table(common.BKTableNameBaseSet).Find(filter).Count(params.Context)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count set model instances failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	moduleCount, err := s.db.Table(common.BKTableNameBaseModule).Find(filter).Count(params.Context)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count module model instances failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	hostCount, err := s.db.Table(common.BKTableNameBaseHost).Find(filter).Count(params.Context)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count host model instances failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	appFilter := map[string]interface{}{
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultAppFlag,
		},
		common.BKDataStatusField: map[string]interface{}{
			common.BKDBNE: common.DataStatusDisabled,
		},
	}
	bizCount, err := s.db.Table(common.BKTableNameBaseApp).Find(appFilter).Count(params.Context)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count application model instances failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	// db.getCollection('cc_ObjectBase').aggregate([{$group: {_id: "$bk_obj_id", count: {$sum : 1}}}])
	pipeline := []map[string]interface{}{
		{
			common.BKDBGroup: map[string]interface{}{
				"_id": "$bk_obj_id",
				"count": map[string]interface{}{
					common.BKDBSum: 1,
				},
			},
		},
	}
	type AggregationItem struct {
		ObjID string `bson:"_id" json:"bk_obj_id"`
		Count int64  `bson:"count" json:"instance_count"`
	}
	aggregationItems := make([]AggregationItem, 0)
	if err := s.db.Table(common.BKTableNameBaseInst).AggregateAll(params.Context, pipeline, &aggregationItems); err != nil {
		return nil, err
	}
	aggregationItems = append(aggregationItems, AggregationItem{
		ObjID: common.BKInnerObjIDHost,
		Count: int64(hostCount),
	}, AggregationItem{
		ObjID: common.BKInnerObjIDSet,
		Count: int64(setCount),
	}, AggregationItem{
		ObjID: common.BKInnerObjIDModule,
		Count: int64(moduleCount),
	}, AggregationItem{
		ObjID: common.BKInnerObjIDApp,
		Count: int64(bizCount),
	})

	return aggregationItems, nil
}

func (s *coreService) CreateModelAttributeGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.CreateModelAttributeGroup{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().CreateModelAttributeGroup(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) SetModelAttributeGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.SetModelAttributeGroup{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
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

func (s *coreService) UpdateModelAttributeGroupByCondition(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.UpdateOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().UpdateModelAttributeGroupByCondition(params, inputData)
}

func (s *coreService) SearchModelAttributeGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributeGroup(params, pathParams("bk_obj_id"), inputData)
	if nil != err {
		return dataResult, err
	}
	for index := range dataResult.Info {
		dataResult.Info[index].GroupName = s.TranslatePropertyGroupName(params.Lang, &dataResult.Info[index])
	}
	return dataResult, err
}

func (s *coreService) SearchModelAttributeGroupByCondition(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributeGroupByCondition(params, inputData)
	if nil != err {
		return dataResult, err
	}
	for index := range dataResult.Info {
		dataResult.Info[index].GroupName = s.TranslatePropertyGroupName(params.Lang, &dataResult.Info[index])
	}
	return dataResult, err
}

func (s *coreService) DeleteModelAttributeGroup(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().DeleteModelAttributeGroup(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) DeleteModelAttributeGroupByCondition(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().DeleteModelAttributeGroupByCondition(params, inputData)
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

func (s *coreService) UpdateModelAttributesByCondition(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.UpdateOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().UpdateModelAttributesByCondition(params, inputData)
}

func (s *coreService) DeleteModelAttribute(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.DeleteOption{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().DeleteModelAttributes(params, pathParams("bk_obj_id"), inputData)
}

func (s *coreService) SearchModelAttributesByCondition(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributesByCondition(params, inputData)
	if nil != err {
		return dataResult, err
	}

	// translate
	for index := range dataResult.Info {
		dataResult.Info[index].PropertyName = s.TranslatePropertyName(params.Lang, &dataResult.Info[index])
		dataResult.Info[index].Placeholder = s.TranslatePlaceholder(params.Lang, &dataResult.Info[index])
		if dataResult.Info[index].PropertyType == common.FieldTypeEnum {
			dataResult.Info[index].Option = s.TranslateEnumName(params.Context, params.Lang, &dataResult.Info[index], dataResult.Info[index].Option)
		}
	}

	return dataResult, err
}

func (s *coreService) SearchModelAttributes(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributes(params, pathParams("bk_obj_id"), inputData)
	if nil != err {
		return dataResult, err
	}

	// translate 主机内置字段bk_state不做翻译
	for index := range dataResult.Info {
		dataResult.Info[index].PropertyName = s.TranslatePropertyName(params.Lang, &dataResult.Info[index])
		dataResult.Info[index].Placeholder = s.TranslatePlaceholder(params.Lang, &dataResult.Info[index])
		if dataResult.Info[index].PropertyType == common.FieldTypeEnum {
			dataResult.Info[index].Option = s.TranslateEnumName(params.Context, params.Lang, &dataResult.Info[index], dataResult.Info[index].Option)
		}
	}

	return dataResult, err
}

func (s *coreService) SearchModelAttrUnique(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	inputData := metadata.QueryCondition{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return s.core.ModelOperation().SearchModelAttrUnique(params, inputData)
}

func (s *coreService) CreateModelAttrUnique(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputDatas := metadata.CreateModelAttrUnique{}
	if err := data.MarshalJSONInto(&inputDatas); nil != err {
		return nil, err
	}

	return s.core.ModelOperation().CreateModelAttrUnique(params, pathParams("bk_obj_id"), inputDatas)
}

func (s *coreService) UpdateModelAttrUnique(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputDatas := metadata.UpdateModelAttrUnique{}
	if err := data.MarshalJSONInto(&inputDatas); nil != err {
		return nil, err
	}
	id, err := strconv.ParseUint(pathParams("id"), 10, 64)
	if err != nil {
		return nil, params.Error.Errorf(common.CCErrCommParamsNeedInt, "id")
	}
	return s.core.ModelOperation().UpdateModelAttrUnique(params, pathParams("bk_obj_id"), id, inputDatas)
}

func (s *coreService) DeleteModelAttrUnique(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputDatas := metadata.DeleteModelAttrUnique{}
	if err := data.MarshalJSONInto(&inputDatas); nil != err {
		return nil, err
	}

	id, err := strconv.ParseUint(pathParams("id"), 10, 64)
	if err != nil {
		return nil, params.Error.Errorf(common.CCErrCommParamsNeedInt, "id")
	}

	return s.core.ModelOperation().DeleteModelAttrUnique(params, pathParams("bk_obj_id"), id, metadata.DeleteModelAttrUnique{Metadata: inputDatas.Metadata})
}
