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
	"configcenter/src/apimachinery"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

type modelAttribute struct {
	clientSet apimachinery.ClientSetInterface
}

func (m *modelAttribute) CreateModelAttributes(ctx core.ContextParams, objID string, inputParam metadata.CreateModelAttributes) (*metadata.CreateManyDataResult, error) {
	return nil, nil
}

func (m *modelAttribute) SetModelAttributes(ctx core.ContextParams, objID string, inputParam metadata.SetModelAttributes) (*metadata.SetManyDataResult, error) {
	return nil, nil
}
func (m *modelAttribute) UpdateModelAttributes(ctx core.ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error) {
	return nil, nil
}
func (m *modelAttribute) DeleteModelAttributes(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error) {
	return nil, nil
}
func (m *modelAttribute) SearchModelAttributes(ctx core.ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	return nil, nil
}
