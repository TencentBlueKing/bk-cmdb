/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parser

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// 注意: 最后返回的模型不一定属于 possibleBizID 对应的业务， 可能是个公有模型, 也可能是私有模型(业务下模型)
// 背景: possibleBizID 是前端传递过来的请求参数，可能只是用于给公有模型创建一个业务私有字段
// 取名 possibleBizID 是为了显式告诉调用方获取到的模型不一定是私有模型
func (ps *parseStream) getPublicOrBizModelByObjectID(possibleBizID int64, objectID string) (metadata.Object, error) {
	filter := map[string]interface{}{
		common.BKObjIDField: objectID,
	}
	if possibleBizID == 0 {
		filter[metadata.MetadataBizField] = mapstr.MapStr{common.BKDBExists: false}
	} else {
		filter[common.BKDBOR] = []map[string]interface{}{
			{
				metadata.MetadataBizField: mapstr.MapStr{common.BKDBExists: false},
			}, {
				metadata.MetadataBizField: possibleBizID,
			},
		}
	}
	return ps.getOneModel(filter)
}

func (ps *parseStream) getOneModel(cond mapstr.MapStr) (metadata.Object, error) {
	model := metadata.Object{}
	models, err := ps.searchModels(cond)
	if err != nil {
		return model, err
	}
	if len(models) == 0 {
		return model, fmt.Errorf("model [%+v] not found", cond)
	}
	if len(models) > 1 {
		return model, fmt.Errorf("model [%+v] not found", cond)
	}
	return models[0], nil
}

func (ps *parseStream) searchModels(cond mapstr.MapStr) ([]metadata.Object, error) {
	model, err := ps.engine.CoreAPI.CoreService().Model().ReadModel(context.Background(), ps.RequestCtx.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return nil, err
	}

	if !model.Result {
		return nil, errors.New(model.Code, model.ErrMsg)
	}
	if len(model.Data.Info) <= 0 {
		return nil, fmt.Errorf("model [%+v] not found", cond)
	}

	models := []metadata.Object{}
	for _, info := range model.Data.Info {
		models = append(models, info.Spec)
	}
	return models, nil
}

func (ps *parseStream) getModelAttribute(cond mapstr.MapStr) ([]metadata.Attribute, error) {
	attr, err := ps.engine.CoreAPI.CoreService().Model().ReadModelAttrByCondition(context.Background(), ps.RequestCtx.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return nil, err
	}

	if !attr.Result {
		return nil, errors.New(attr.Code, attr.ErrMsg)
	}
	if len(attr.Data.Info) <= 0 {
		return nil, fmt.Errorf("attribute [%+v] not found", cond)
	}

	attrs := []metadata.Attribute{}
	for _, info := range attr.Data.Info {
		attrs = append(attrs, info)
	}
	return attrs, nil
}

func (ps *parseStream) getCls(clsID string) (metadata.Classification, error) {
	model, err := ps.engine.CoreAPI.CoreService().Model().ReadModelClassification(context.Background(), ps.RequestCtx.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{common.BKClassificationIDField: clsID}})
	if err != nil {
		return metadata.Classification{}, err
	}
	if len(model.Data.Info) <= 0 {
		return metadata.Classification{}, fmt.Errorf("model classification [%s] not found", clsID)
	}
	return model.Data.Info[0], nil
}

func (ps *parseStream) getAttributeGroup(cond interface{}) ([]metadata.Group, error) {
	mspstrCond, err := mapstr.NewFromInterface(cond)
	if err != nil {
		return nil, err
	}
	groups, err := ps.engine.CoreAPI.CoreService().Model().ReadAttributeGroupByCondition(context.Background(), ps.RequestCtx.Header,
		metadata.QueryCondition{Condition: mspstrCond})
	if err != nil {
		return nil, err
	}

	if !groups.Result {
		return nil, errors.New(groups.Code, groups.ErrMsg)
	}

	return groups.Data.Info, nil
}

func (ps *parseStream) isMainlineModel(modelID string) (bool, error) {
	asst, err := ps.getModelAssociation(mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline})
	if err != nil {
		return false, err
	}
	for _, mainline := range asst {
		if mainline.ObjectID == modelID {
			return true, nil
		}
	}

	return false, nil
}

func (ps *parseStream) getModelAssociation(cond mapstr.MapStr) ([]metadata.Association, error) {
	asst, err := ps.engine.CoreAPI.CoreService().Association().ReadModelAssociation(context.Background(), ps.RequestCtx.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return nil, err
	}

	if !asst.Result {
		return nil, errors.New(asst.Code, asst.ErrMsg)
	}

	if len(asst.Data.Info) <= 0 {
		return nil, fmt.Errorf("model association [%+v] not found", cond)
	}

	return asst.Data.Info, nil
}

func (ps *parseStream) getInstAssociation(cond mapstr.MapStr) (metadata.InstAsst, error) {
	asst, err := ps.engine.CoreAPI.CoreService().Association().ReadInstAssociation(context.Background(), ps.RequestCtx.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return metadata.InstAsst{}, err
	}

	if !asst.Result {
		return metadata.InstAsst{}, errors.New(asst.Code, asst.ErrMsg)
	}

	if len(asst.Data.Info) <= 0 {
		return metadata.InstAsst{}, fmt.Errorf("model association [%+v] not found", cond)
	}

	return asst.Data.Info[0], nil
}

func (ps *parseStream) getOneClassification(cond mapstr.MapStr) (metadata.Classification, error) {
	classification := metadata.Classification{}
	classifications, err := ps.getClassification(cond)
	if err != nil {
		return classification, err
	}
	if len(classifications) <= 0 {
		return classification, fmt.Errorf("classification [%+v] not found", cond)
	}
	if len(classifications) > 1 {
		return classification, fmt.Errorf("get multiple classification with [%+v]", cond)
	}
	return classifications[0], nil
}

func (ps *parseStream) getClassification(cond mapstr.MapStr) ([]metadata.Classification, error) {
	classificationResult, err := ps.engine.CoreAPI.CoreService().Model().ReadModelClassification(context.Background(), ps.RequestCtx.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return nil, err
	}

	if !classificationResult.Result {
		return nil, errors.New(classificationResult.Code, classificationResult.ErrMsg)
	}
	if len(classificationResult.Data.Info) <= 0 {
		return nil, fmt.Errorf("classification [%+v] not found", cond)
	}

	return classificationResult.Data.Info, nil
}

func (ps *parseStream) getModelUnique(cond mapstr.MapStr) (metadata.ObjectUnique, error) {
	unique := metadata.ObjectUnique{}
	filter := metadata.QueryCondition{Condition: cond}
	modelUniqueResult, err := ps.engine.CoreAPI.CoreService().Model().ReadModelAttrUnique(context.Background(), ps.RequestCtx.Header, filter)
	if err != nil {
		return unique, err
	}

	if !modelUniqueResult.Result {
		return unique, errors.New(modelUniqueResult.Code, modelUniqueResult.ErrMsg)
	}
	if len(modelUniqueResult.Data.Info) <= 0 {
		return unique, fmt.Errorf("model unique [%+v] not found", cond)
	}
	if len(modelUniqueResult.Data.Info) > 1 {
		return unique, fmt.Errorf("get multiple model unique with [%+v]", cond)
	}
	return modelUniqueResult.Data.Info[0], nil
}
