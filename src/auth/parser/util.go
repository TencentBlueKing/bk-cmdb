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

func (ps *parseStream) getModel(modelID interface{}) ([]metadata.Object, error) {
	model, err := ps.engine.CoreAPI.CoreService().Model().ReadModel(context.Background(), ps.RequestCtx.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{common.BKObjIDField: modelID}})
	if err != nil {
		return nil, err
	}

	if !model.Result {
		return nil, errors.New(model.Code, model.ErrMsg)
	}
	if len(model.Data.Info) <= 0 {
		return nil, fmt.Errorf("model [%s] not found", modelID)
	}

	models := []metadata.Object{}
	for _, info := range model.Data.Info {
		models = append(models, info.Spec)
	}
	return models, nil
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

func (ps *parseStream) getModelAssociation(id int64) (metadata.Association, error) {
	asst, err := ps.engine.CoreAPI.CoreService().Association().ReadModelAssociation(context.Background(), ps.RequestCtx.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{common.BKFieldID: id}})
	if err != nil {
		return metadata.Association{}, err
	}

	if !asst.Result {
		return metadata.Association{}, errors.New(asst.Code, asst.ErrMsg)
	}

	if len(asst.Data.Info) <= 0 {
		return metadata.Association{}, fmt.Errorf("model association [%d] not found", id)
	}

	return asst.Data.Info[0], nil
}
