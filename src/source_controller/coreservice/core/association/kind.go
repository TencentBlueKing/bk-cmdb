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

type associationKind struct {
	dbProxy dal.RDB
}

func (m *associationKind) CreateAssociationKind(ctx core.ContextParams, inputParam metadata.CreateAssociationKind) (*metadata.CreateOneDataResult, error) {
	return nil, nil
}
func (m *associationKind) CreateManyAssociationKind(ctx core.ContextParams, inputParam metadata.CreateManyAssociationKind) (*metadata.CreateManyAssociationKind, error) {
	return nil, nil
}
func (m *associationKind) SetAssociationKind(ctx core.ContextParams, inputParam metadata.SetAssociationKind) (*metadata.SetOneDataResult, error) {
	return nil, nil
}
func (m *associationKind) SetManyAssociationKind(ctx core.ContextParams, inputParam metadata.SetManyAssociationKind) (*metadata.SetManyDataResult, error) {
	return nil, nil
}
func (m *associationKind) UpdateAssociationKind(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error) {
	return nil, nil
}
func (m *associationKind) DeleteAssociationKind(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error) {
	return nil, nil
}
func (m *associationKind) CascadeDeleteAssociationKind(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteOption, error) {
	return nil, nil
}
func (m *associationKind) SearchAssociationKind(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	return nil, nil
}
