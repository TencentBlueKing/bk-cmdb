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

package instances

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.InstanceOperation = (*instanceManager)(nil)

type instanceManager struct {
	dbProxy dal.RDB
}

// New create a new instance manager instance
func New(dbProxy dal.RDB) core.InstanceOperation {
	return &instanceManager{
		dbProxy: dbProxy,
	}
}

func (m *instanceManager) CreateModelInstance(ctx core.ContextParams, objID string, inputParam metadata.CreateModelInstance) (*metadata.CreateOneDataResult, error) {
	return nil, nil
}
func (m *instanceManager) CreateManyModelInstance(ctx core.ContextParams, objID string, inputParam metadata.CreateManyModelInstance) (*metadata.CreateManyDataResult, error) {
	return nil, nil
}
func (m *instanceManager) SetModelInstance(ctx core.ContextParams, objID string, inputParam metadata.SetModelInstance) (*metadata.SetDataResult, error) {
	return nil, nil
}
func (m *instanceManager) SetManyModelInstance(ctx core.ContextParams, objID string, inputParam metadata.SetManyModelInstance) (*metadata.SetDataResult, error) {
	return nil, nil
}
func (m *instanceManager) UpdateModelInstance(ctx core.ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {
	return nil, nil
}
func (m *instanceManager) SearchModelInstance(ctx core.ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	return nil, nil
}
func (m *instanceManager) DeleteModelInstance(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	return nil, nil
}
func (m *instanceManager) CascadeDeleteModelInstance(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	return nil, nil
}
