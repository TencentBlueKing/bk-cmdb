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

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (ps *parseStream) getModel(cond mapstr.MapStr) ([]metadata.Object, error) {
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
		return meta.Plat, nil
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
