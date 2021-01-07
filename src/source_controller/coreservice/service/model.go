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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

func (s *coreService) CreateManyModelClassification(ctx *rest.Contexts) {
	inputDatas := metadata.CreateManyModelClassifiaction{}
	if err := ctx.DecodeInto(&inputDatas); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CreateManyModelClassification(ctx.Kit, inputDatas))
}

func (s *coreService) CreateOneModelClassification(ctx *rest.Contexts) {
	inputData := metadata.CreateOneModelClassification{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CreateOneModelClassification(ctx.Kit, inputData))
}

func (s *coreService) SetOneModelClassification(ctx *rest.Contexts) {
	inputData := metadata.SetOneModelClassification{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().SetOneModelClassification(ctx.Kit, inputData))
}

func (s *coreService) SetManyModelClassification(ctx *rest.Contexts) {
	inputDatas := metadata.SetManyModelClassification{}
	if err := ctx.DecodeInto(&inputDatas); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().SetManyModelClassification(ctx.Kit, inputDatas))
}

func (s *coreService) UpdateModelClassification(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelClassification(ctx.Kit, inputData))
}

func (s *coreService) DeleteModelClassification(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModelClassification(ctx.Kit, inputData))
}

func (s *coreService) SearchModelClassification(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelClassification(ctx.Kit, inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	// translate language
	lang := s.Language(ctx.Kit.Header)
	for index := range dataResult.Info {
		defaultClassificationMap := map[string]bool{
			"bk_host_manage":  true,
			"bk_biz_topo":     true,
			"bk_organization": true,
			"bk_network":      true,
		}
		if defaultClassificationMap[dataResult.Info[index].ClassificationID] {
			dataResult.Info[index].ClassificationName = s.TranslateClassificationName(lang, &dataResult.Info[index])
		}
	}
	ctx.RespEntity(dataResult)
}

func (s *coreService) CreateModel(ctx *rest.Contexts) {
	inputData := metadata.CreateModel{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CreateModel(ctx.Kit, inputData))
}

func (s *coreService) SetModel(ctx *rest.Contexts) {
	inputData := metadata.SetModel{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().SetModel(ctx.Kit, inputData))
}

func (s *coreService) UpdateModel(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModel(ctx.Kit, inputData))
}

func (s *coreService) DeleteModel(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModel(ctx.Kit, inputData))
}

func (s *coreService) CascadeDeleteModel(ctx *rest.Contexts) {
	idStr := ctx.Request.PathParameter(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID))
		return
	}
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CascadeDeleteModel(ctx.Kit, id))
}

func (s *coreService) SearchModel(ctx *rest.Contexts) {

	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelWithAttribute(ctx.Kit, inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	// translate
	lang := s.Language(ctx.Kit.Header)
	for modelIdx := range dataResult.Info {
		if needTranslateObjMap[dataResult.Info[modelIdx].Spec.ObjectID] {
			dataResult.Info[modelIdx].Spec.ObjectName = s.TranslateObjectName(lang, &dataResult.Info[modelIdx].Spec)
		}
		for attributeIdx := range dataResult.Info[modelIdx].Attributes {
			if dataResult.Info[modelIdx].Attributes[attributeIdx].IsPre || dataResult.Info[modelIdx].Spec.IsPre || needTranslateObjMap[dataResult.Info[modelIdx].Spec.ObjectID] {
				dataResult.Info[modelIdx].Attributes[attributeIdx].PropertyName = s.TranslatePropertyName(lang, &dataResult.Info[modelIdx].Attributes[attributeIdx])
				dataResult.Info[modelIdx].Attributes[attributeIdx].Placeholder = s.TranslatePlaceholder(lang, &dataResult.Info[modelIdx].Attributes[attributeIdx])
				if dataResult.Info[modelIdx].Attributes[attributeIdx].PropertyType == common.FieldTypeEnum {
					dataResult.Info[modelIdx].Attributes[attributeIdx].Option = s.TranslateEnumName(ctx.Kit.Ctx, lang, &dataResult.Info[modelIdx].Attributes[attributeIdx], dataResult.Info[modelIdx].Attributes[attributeIdx].Option)
				}
			}
		}
	}

	ctx.RespEntity(dataResult)
}

// GetModelStatistics 用于统计各个模型的实例数(Web页面展示需要)
func (s *coreService) GetModelStatistics(ctx *rest.Contexts) {
	filter := map[string]interface{}{}
	setCount, err := mongodb.Client().Table(common.BKTableNameBaseSet).Find(filter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count set model instances failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	moduleCount, err := mongodb.Client().Table(common.BKTableNameBaseModule).Find(filter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count module model instances failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	hostCount, err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(filter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count host model instances failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	appFilter := map[string]interface{}{
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultAppFlag,
		},
		common.BKDataStatusField: map[string]interface{}{
			common.BKDBNE: common.DataStatusDisabled,
		},
	}
	bizCount, err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(appFilter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count application model instances failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
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
	if err := mongodb.Client().Table(common.BKTableNameBaseInst).AggregateAll(ctx.Kit.Ctx, pipeline, &aggregationItems); err != nil {
		ctx.RespAutoError(err)
		return
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

	ctx.RespEntity(aggregationItems)
}

func (s *coreService) CreateModelAttributeGroup(ctx *rest.Contexts) {
	inputData := metadata.CreateModelAttributeGroup{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().CreateModelAttributeGroup(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

func (s *coreService) SetModelAttributeGroup(ctx *rest.Contexts) {
	inputData := metadata.SetModelAttributeGroup{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().SetModelAttributeGroup(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

func (s *coreService) UpdateModelAttributeGroup(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelAttributeGroup(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

func (s *coreService) UpdateModelAttributeGroupByCondition(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelAttributeGroupByCondition(ctx.Kit, inputData))
}

func (s *coreService) SearchModelAttributeGroup(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributeGroup(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	lang := s.Language(ctx.Kit.Header)
	for index := range dataResult.Info {
		if dataResult.Info[index].IsDefault {
			dataResult.Info[index].GroupName = s.TranslatePropertyGroupName(lang, &dataResult.Info[index])
		}
	}
	ctx.RespEntity(dataResult)
}

func (s *coreService) SearchModelAttributeGroupByCondition(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributeGroupByCondition(ctx.Kit, inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}
	lang := s.Language(ctx.Kit.Header)
	for index := range dataResult.Info {
		if dataResult.Info[index].IsDefault {
			dataResult.Info[index].GroupName = s.TranslatePropertyGroupName(lang, &dataResult.Info[index])
		}
	}
	ctx.RespEntity(dataResult)
}

func (s *coreService) DeleteModelAttributeGroup(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModelAttributeGroup(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

func (s *coreService) DeleteModelAttributeGroupByCondition(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModelAttributeGroupByCondition(ctx.Kit, inputData))
}

func (s *coreService) CreateModelAttributes(ctx *rest.Contexts) {

	inputData := metadata.CreateModelAttributes{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CreateModelAttributes(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

func (s *coreService) SetModelAttributes(ctx *rest.Contexts) {

	inputData := metadata.SetModelAttributes{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().SetModelAttributes(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

func (s *coreService) UpdateModelAttributes(ctx *rest.Contexts) {

	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelAttributes(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

func (s *coreService) UpdateModelAttributesIndex(ctx *rest.Contexts) {

	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelAttributesIndex(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

func (s *coreService) UpdateModelAttributesByCondition(ctx *rest.Contexts) {

	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelAttributesByCondition(ctx.Kit, inputData))
}

func (s *coreService) DeleteModelAttribute(ctx *rest.Contexts) {

	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModelAttributes(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData))
}

func (s *coreService) SearchModelAttributesByCondition(ctx *rest.Contexts) {

	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributesByCondition(ctx.Kit, inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	// translate
	lang := s.Language(ctx.Kit.Header)
	for index := range dataResult.Info {
		if dataResult.Info[index].IsPre || needTranslateObjMap[dataResult.Info[index].ObjectID] {
			dataResult.Info[index].PropertyName = s.TranslatePropertyName(lang, &dataResult.Info[index])
			dataResult.Info[index].Placeholder = s.TranslatePlaceholder(lang, &dataResult.Info[index])
			if dataResult.Info[index].PropertyType == common.FieldTypeEnum {
				dataResult.Info[index].Option = s.TranslateEnumName(ctx.Kit.Ctx, lang, &dataResult.Info[index], dataResult.Info[index].Option)
			}
		}
	}

	ctx.RespEntity(dataResult)
}

func (s *coreService) SearchModelAttributes(ctx *rest.Contexts) {

	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributes(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	// translate 主机内置字段bk_state不做翻译
	lang := s.Language(ctx.Kit.Header)
	for index := range dataResult.Info {
		if dataResult.Info[index].IsPre || needTranslateObjMap[dataResult.Info[index].ObjectID] {
			dataResult.Info[index].PropertyName = s.TranslatePropertyName(lang, &dataResult.Info[index])
			dataResult.Info[index].Placeholder = s.TranslatePlaceholder(lang, &dataResult.Info[index])
			if dataResult.Info[index].PropertyType == common.FieldTypeEnum {
				dataResult.Info[index].Option = s.TranslateEnumName(ctx.Kit.Ctx, lang, &dataResult.Info[index], dataResult.Info[index].Option)
			}
		}
	}

	ctx.RespEntity(dataResult)
}

func (s *coreService) SearchModelAttrUnique(ctx *rest.Contexts) {

	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().SearchModelAttrUnique(ctx.Kit, inputData))
}

func (s *coreService) CreateModelAttrUnique(ctx *rest.Contexts) {
	inputDatas := metadata.CreateModelAttrUnique{}
	if err := ctx.DecodeInto(&inputDatas); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().CreateModelAttrUnique(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), inputDatas))
}

func (s *coreService) UpdateModelAttrUnique(ctx *rest.Contexts) {
	inputDatas := metadata.UpdateModelAttrUnique{}
	if err := ctx.DecodeInto(&inputDatas); nil != err {
		ctx.RespAutoError(err)
		return
	}
	id, err := strconv.ParseUint(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "id"))
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelAttrUnique(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), id, inputDatas))
}

func (s *coreService) DeleteModelAttrUnique(ctx *rest.Contexts) {
	id, err := strconv.ParseUint(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "id"))
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModelAttrUnique(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), id))
}
