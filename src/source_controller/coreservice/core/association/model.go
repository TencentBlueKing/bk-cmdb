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

package association

import (
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type associationModel struct {
	dbProxy dal.RDB
}

func (m *associationModel) CreateModelAssociation(ctx core.ContextParams, inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error) {
	return nil, nil
}
func (m *associationModel) SetModelAssociation(ctx core.ContextParams, inputParam metadata.SetModelAssociation) (*metadata.SetDataResult, error) {
	return nil, nil
}
func (m *associationModel) UpdateModelAssociation(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error) {
	return nil, nil
}
func (m *associationModel) SearchModelAssociation(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	return nil, nil
}
func (m *associationModel) DeleteModelAssociation(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error) {
	return nil, nil
}
func (m *associationModel) CascadeDeleteModelAssociation(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error) {
	return nil, nil
}
