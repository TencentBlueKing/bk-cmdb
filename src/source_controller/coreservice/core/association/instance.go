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

type associationInstance struct {
	dbProxy dal.RDB
}

func (m *associationInstance) CreateOneInstanceAssociation(ctx core.ContextParams, inputParam metadata.CreateOneInstanceAssociation) (*metadata.CreateOneDataResult, error) {
	return nil, nil
}
func (m *associationInstance) SetOneInstanceAssociation(ctx core.ContextParams, inputParam metadata.SetOneInstanceAssociation) (*metadata.SetDataResult, error) {
	return nil, nil
}
func (m *associationInstance) CreateManyInstanceAssociation(ctx core.ContextParams, inputParam metadata.CreateManyInstanceAssociation) (*metadata.CreateManyDataResult, error) {
	return nil, nil
}
func (m *associationInstance) SetManyInstanceAssociation(ctx core.ContextParams, inputParam metadata.SetManyInstanceAssociation) (*metadata.SetDataResult, error) {
	return nil, nil
}
func (m *associationInstance) UpdateInstanceAssociation(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error) {
	return nil, nil
}
func (m *associationInstance) SearchInstanceAssociation(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	return nil, nil
}
func (m *associationInstance) DeleteInstanceAssociation(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error) {
	return nil, nil
}
