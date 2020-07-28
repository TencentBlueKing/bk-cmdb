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
	"sync/atomic"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
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

func (ps *parseStream) getClassIDWithObject(obj string) (int64, error) {
	filter := mapstr.MapStr{
		common.BKObjIDField: obj,
	}

	model, err := ps.getOneModel(filter)
	if err != nil {
		return 0, err
	}

	filter = map[string]interface{}{
		common.BKClassificationIDField: model.ObjCls,
	}

	class, err := ps.getOneClassification(filter)
	if err != nil {
		return 0, err
	}

	return class.ID, nil
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

func (ps *parseStream) getInstanceTypeByObject(objID string) (meta.ResourceType, error) {
	switch objID {
	case common.BKInnerObjIDPlat:
		return meta.CloudAreaInstance, nil
	case common.BKInnerObjIDHost:
		return meta.HostInstance, nil
	case common.BKInnerObjIDModule:
		return meta.ModelModule, nil
	case common.BKInnerObjIDSet:
		return meta.ModelSet, nil
	case common.BKInnerObjIDApp:
		return meta.Business, nil
	case common.BKInnerObjIDProc:
		return meta.Process, nil
	}
	isMainline, err := ps.isMainlineModel(objID)
	if err != nil {
		return "", err
	}
	if isMainline {
		return meta.MainlineInstance, nil
	}
	return meta.ModelInstance, nil
}

func (ps *parseStream) getBizIDByHostID(hostID int64) (int64, error) {
	result, err := ps.engine.CoreAPI.CoreService().Host().GetHostModuleRelation(context.Background(), ps.RequestCtx.Header,
		&metadata.HostModuleRelationRequest{HostIDArr: []int64{hostID}, Fields: []string{common.BKAppIDField}})
	if err != nil {
		return 0, err
	}
	for _, relation := range result.Data.Info {
		return relation.AppID, nil
	}
	return 0, nil
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

// get hosts relation which these hosts must be in the resource pool
func (ps *parseStream) getRscPoolHostModuleRelation(hostIDs []int64) (map[int64]int64, error) {
	opt := &metadata.HostModuleRelationRequest{
		HostIDArr: hostIDs,
		Fields:    []string{common.BKAppIDField, common.BKHostIDField, common.BKModuleIDField},
	}

	result, err := ps.engine.CoreAPI.CoreService().Host().GetHostModuleRelation(context.Background(),
		ps.RequestCtx.Header, opt)
	if err != nil {
		return nil, err
	}

	if !result.Result {
		return nil, errors.New(result.Code, result.ErrMsg)
	}

	resourceBiz, err := ps.getResourcePoolBusinessID()
	if err != nil {
		return nil, err
	}

	relation := make(map[int64]int64)
	for _, rel := range result.Data.Info {
		if rel.AppID != resourceBiz {
			return nil, errors.New(common.CCErrCommParamsInvalid, "host does not belongs to host pool")
		}

		relation[rel.HostID] = rel.ModuleID
	}

	return relation, nil
}

// TODO: support multiple supplier account.
var resourcePoolBizID int64

func (ps *parseStream) getResourcePoolBusinessID() (biz int64, err error) {
	if id := atomic.LoadInt64(&resourcePoolBizID); id != 0 {
		return id, nil
	}

	opt := &metadata.QueryCondition{
		Fields: []string{common.BKAppIDField},
		Page:   metadata.BasePage{},
		Condition: mapstr.MapStr{
			common.BkSupplierAccount: ps.RequestCtx.Header.Get(common.BKHTTPOwnerID),
			"default":                1,
		},
	}

	result, err := ps.engine.CoreAPI.CoreService().Instance().ReadInstance(context.Background(), ps.RequestCtx.Header,
		common.BKInnerObjIDApp, opt)
	if err != nil {
		return 0, err
	}

	if !result.Result {
		return 0, errors.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		return 0, errors.New(common.CCErrCommParamsIsInvalid, "")
	}

	id, err := util.GetInt64ByInterface(result.Data.Info[0][common.BKAppIDField])
	if err != nil {
		return 0, errors.New(common.CCErrorUnknownOrUnrecognizedError, "invalid resource biz id")
	}

	atomic.StoreInt64(&resourcePoolBizID, id)
	return id, nil
}

type hostPool struct {
	Business int64   `json:"bk_biz_id"`
	HostID   []int64 `json:"bk_host_id"`
}

var resourcePoolDefaultDirID int64

func (ps *parseStream) getResourcePoolDefaultDirID() (dirID int64, err error) {
	if id := atomic.LoadInt64(&resourcePoolDefaultDirID); id != 0 {
		return id, nil
	}

	opt := &metadata.QueryCondition{
		Fields: []string{common.BKModuleIDField},
		Page:   metadata.BasePage{},
		Condition: mapstr.MapStr{
			common.BkSupplierAccount: ps.RequestCtx.Header.Get(common.BKHTTPOwnerID),
			common.BKDefaultField:                common.DefaultResModuleFlag,
		},
	}

	result, err := ps.engine.CoreAPI.CoreService().Instance().ReadInstance(context.Background(), ps.RequestCtx.Header,
		common.BKInnerObjIDModule, opt)
	if err != nil {
		return 0, err
	}

	if !result.Result {
		return 0, result.CCError()
	}

	if len(result.Data.Info) == 0 {
		return 0, errors.New(common.CCErrCommParamsIsInvalid, common.BKModuleIDField)
	}

	id, err := util.GetInt64ByInterface(result.Data.Info[0][common.BKModuleIDField])
	if err != nil {
		return 0, errors.New(common.CCErrorUnknownOrUnrecognizedError, "invalid resource pool default directory id")
	}

	atomic.StoreInt64(&resourcePoolDefaultDirID, id)
	return id, nil
}
