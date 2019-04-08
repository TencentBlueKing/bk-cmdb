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

package model

import (
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type modelAttrUnique struct {
	dbProxy dal.RDB
}

func (m *modelAttrUnique) CreateModelAttrUnique(ctx core.ContextParams, objID string, data metadata.CreateModelAttrUnique) (*metadata.CreateOneDataResult, error) {
	id, err := m.createModelAttrUnique(ctx, objID, data)
	if err != nil {
		return nil, err
	}
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, nil
}

func (m *modelAttrUnique) UpdateModelAttrUnique(ctx core.ContextParams, objID string, id uint64, data metadata.UpdateModelAttrUnique) (*metadata.UpdatedCount, error) {
	err := m.updateModelAttrUnique(ctx, objID, id, data)
	if err != nil {
		return nil, err
	}
	return &metadata.UpdatedCount{Count: 1}, nil
}

func (m *modelAttrUnique) DeleteModelAttrUnique(ctx core.ContextParams, objID string, id uint64, meta metadata.DeleteModelAttrUnique) (*metadata.DeletedCount, error) {
	err := m.deleteModelAttrUnique(ctx, objID, id, meta)
	if err != nil {
		return nil, err
	}
	return &metadata.DeletedCount{Count: 1}, nil
}

func (m *modelAttrUnique) SearchModelAttrUnique(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryUniqueResult, error) {

	uniqueItems, err := m.searchModelAttrUnique(ctx, inputParam)
	if nil != err {
		return &metadata.QueryUniqueResult{Info: []metadata.ObjectUnique{}}, err
	}
	dataResult := &metadata.QueryUniqueResult{Info: []metadata.ObjectUnique{}}
	dataResult.Count, err = m.countModelAttrUnique(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.QueryUniqueResult{Info: []metadata.ObjectUnique{}}, err
	}
	if len(uniqueItems) > 0 {
		dataResult.Info = uniqueItems
	}

	return dataResult, nil
}
