/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package operation

import (
	"context"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// UniqueOperationInterface group operation methods
type UniqueOperationInterface interface {
	CreateObjectUnique(params types.ContextParams, data mapstr.MapStr) (model.Group, error)
	DeleteObjectUnique(params types.ContextParams, groupID int64) error
	FindObjectUnique(params types.ContextParams, cond condition.Condition) ([]model.Group, error)
	UpdateObjectUnique(params types.ContextParams, cond *metadata.UpdateGroupCondition) error
	SetProxy(modelFactory model.Factory, instFactory inst.Factory, obj ObjectOperationInterface)
}

// NewUniqueOperation create a new group operation instance
func NewUniqueOperation(client apimachinery.ClientSetInterface) UniqueOperationInterface {
	return &unique{
		clientSet: client,
	}
}

type unique struct {
	clientSet apimachinery.ClientSetInterface
}

func (a *unique) Search(params types.ContextParams, objectID string) (resp *metadata.SearchUniqueResult, err error) {
	resp, err := a.clientSet.ObjectController().Unique().Search(context.Background(), params.Header, objectID)
	if err != nil {
		return nil, params.Err.New(errorCode, msg)
	}
	if !resp.Result {
		return nil, params.Err.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())

	}
	return
}
func (a *unique) Create(params types.ContextParams, objectID string, request *metadata.UniqueKind) (resp *metadata.CreateUniqueTypeResult, err error) {
	return a.clientSet.ObjectController().Unique().Create(context.Background(), params.Header, request)
}
func (a *unique) Update(params types.ContextParams, objectID string, id int64, request *metadata.UpdateUniqueTypeRequest) (resp *metadata.UpdateUniqueTypeResult, err error) {
	return a.clientSet.ObjectController().Unique().Update(context.Background(), params.Header, asstTypeID, request)
}
func (a *unique) Delete(params types.ContextParams, objectID string, id int64) (resp *metadata.DeleteUniqueTypeResult, err error) {
	return a.clientSet.ObjectController().Unique().Delete(context.Background(), params.Header, asstTypeID)
}
