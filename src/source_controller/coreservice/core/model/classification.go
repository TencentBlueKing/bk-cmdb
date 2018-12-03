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

package model

import (
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

type modelClassification struct {
}

func (m *modelClassification) CreateOneModelClassification(ctx core.ContextParams, inputParam metadata.CreateOneModelClassification) (*metadata.CreateOneDataResult, error) {
	return nil, nil
}
func (m *modelClassification) CreateManyModelClassification(ctx core.ContextParams, inputParam metadata.CreateManyModelClassifiaction) (*metadata.CreateManyDataResult, error) {
	return nil, nil
}
func (m *modelClassification) SetManyModelClassification(ctx core.ContextParams, inputParam metadata.SetManyModelClassification) (*metadata.SetManyResult, error) {
	return nil, nil
}
func (m *modelClassification) SetOneModelClassification(ctx core.ContextParams, inputParam metadata.SetOneModelClassification) (*metadata.SetOneResult, error) {
	return nil, nil
}
func (m *modelClassification) UpdateModelClassification(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error) {
	return nil, nil
}
func (m *modelClassification) DeleteModelClassificaiton(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error) {
	return nil, nil
}
func (m *modelClassification) CascadeDeleteModeClassification(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error) {
	return nil, nil
}
func (m *modelClassification) SearchModelClassification(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	return nil, nil
}
