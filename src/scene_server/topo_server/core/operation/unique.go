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
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// UniqueOperationInterface group operation methods
type UniqueOperationInterface interface {
	Create(params types.ContextParams, objectID string, request *metadata.CreateUniqueRequest) (uniqueID *metadata.RspID, err error)
	Update(params types.ContextParams, objectID string, id uint64, request *metadata.UpdateUniqueRequest) (err error)
	Delete(params types.ContextParams, objectID string, id uint64) (err error)
	Search(params types.ContextParams, objectID string) (objectUniques []metadata.ObjectUnique, err error)
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

func (a *unique) Create(params types.ContextParams, objectID string, request *metadata.CreateUniqueRequest) (uniqueID *metadata.RspID, err error) {
	resp, err := a.clientSet.ObjectController().Unique().Create(context.Background(), params.Header, objectID, request)
	if err != nil {
		blog.Errorf("[UniqueOperation] create for %s, %#v failed %v", objectID, request, err)
		return nil, params.Err.Error(common.CCErrTopoObjectUniqueCreateFailed)
	}
	if !resp.Result {
		return nil, params.Err.New(resp.Code, resp.ErrMsg)
	}
	return &resp.Data, nil

}
func (a *unique) Update(params types.ContextParams, objectID string, id uint64, request *metadata.UpdateUniqueRequest) (err error) {
	resp, err := a.clientSet.ObjectController().Unique().Update(context.Background(), params.Header, objectID, id, request)
	if err != nil {
		blog.Errorf("[UniqueOperation] update for %s, %d, %#v failed %v", objectID, id, request, err)
		return params.Err.Error(common.CCErrTopoObjectUniqueUpdateFailed)
	}
	if !resp.Result {
		return params.Err.New(resp.Code, resp.ErrMsg)
	}
	return nil
}
func (a *unique) Delete(params types.ContextParams, objectID string, id uint64) (err error) {
	resp, err := a.clientSet.ObjectController().Unique().Delete(context.Background(), params.Header, objectID, id)
	if err != nil {
		blog.Errorf("[UniqueOperation] delete for %s, %d failed %v", objectID, id, err)
		return params.Err.Error(common.CCErrTopoObjectUniqueDeleteFailed)
	}
	if !resp.Result {
		return params.Err.New(resp.Code, resp.ErrMsg)
	}
	return nil
}
func (a *unique) Search(params types.ContextParams, objectID string) (objectUniques []metadata.ObjectUnique, err error) {
	resp, err := a.clientSet.ObjectController().Unique().Search(context.Background(), params.Header, objectID)
	if err != nil {
		blog.Errorf("[UniqueOperation] search for %s, failed %v", objectID, err)
		return nil, params.Err.Error(common.CCErrTopoObjectUniqueSearchFailed)
	}
	if !resp.Result {
		return nil, params.Err.New(resp.Code, resp.ErrMsg)
	}
	return resp.Data, nil
}
