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
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
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
func NewUniqueOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) UniqueOperationInterface {
	return &unique{
		clientSet:   client,
		authManager: authManager,
	}
}

type unique struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

func (a *unique) Create(params types.ContextParams, objectID string, request *metadata.CreateUniqueRequest) (uniqueID *metadata.RspID, err error) {

	unique := metadata.ObjectUnique{
		ObjID:     request.ObjID,
		Keys:      request.Keys,
		MustCheck: request.MustCheck,
	}

	if nil != params.MetaData {
		unique.Metadata = *params.MetaData
	}
	resp, err := a.clientSet.CoreService().Model().CreateModelAttrUnique(context.Background(), params.Header, objectID, metadata.CreateModelAttrUnique{Data: unique})
	if err != nil {
		blog.Errorf("[UniqueOperation] create for %s, %#v failed %v, rid: %s", objectID, request, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrTopoObjectUniqueCreateFailed)
	}
	if !resp.Result {
		return nil, params.Err.New(resp.Code, resp.ErrMsg)
	}

	return &metadata.RspID{ID: int64(resp.Data.Created.ID)}, nil
}

func (a *unique) Update(params types.ContextParams, objectID string, id uint64, request *metadata.UpdateUniqueRequest) (err error) {
	update := metadata.UpdateModelAttrUnique{
		Data: *request,
	}
	resp, err := a.clientSet.CoreService().Model().UpdateModelAttrUnique(context.Background(), params.Header, objectID, id, update)
	if err != nil {
		blog.Errorf("[UniqueOperation] update for %s, %d, %#v failed %v, rid: %s", objectID, id, request, err, params.ReqID)
		return params.Err.Error(common.CCErrTopoObjectUniqueUpdateFailed)
	}
	if !resp.Result {
		return params.Err.New(resp.Code, resp.ErrMsg)
	}

	// auth: update register to iam
	if err := a.authManager.UpdateRegisteredModelUniqueByID(params.Context, params.Header, int64(id)); err != nil {
		blog.V(2).Infof("update unique %d for model %s failed, authorization failed, err: %+v, rid: %s", id, objectID, err, params.ReqID)
		return err
	}
	return nil
}

func (a *unique) Delete(params types.ContextParams, objectID string, id uint64) (err error) {
	meta := metadata.Metadata{}
	if params.MetaData != nil {
		meta = *params.MetaData
	}
	resp, err := a.clientSet.CoreService().Model().DeleteModelAttrUnique(context.Background(), params.Header, objectID, id, metadata.DeleteModelAttrUnique{Metadata: meta})
	if err != nil {
		blog.Errorf("[UniqueOperation] delete for %s, %d failed %v, rid: %s", objectID, id, err, params.ReqID)
		return params.Err.Error(common.CCErrTopoObjectUniqueDeleteFailed)
	}
	if !resp.Result {
		return params.Err.New(resp.Code, resp.ErrMsg)
	}
	// auth: deregister to iam
	if err := a.authManager.DeregisterModelUniqueByID(params.Context, params.Header, int64(id)); err != nil {
		blog.V(2).Infof("deregister unique %d for model %s failed, authorization failed, err: %+v, rid: %s", id, objectID, err, params.ReqID)
		return err
	}
	return nil
}

func (a *unique) Search(params types.ContextParams, objectID string) (objectUniques []metadata.ObjectUnique, err error) {
	fCond := condition.CreateCondition().Field(common.BKObjIDField).Eq(objectID).ToMapStr()
	if nil != params.MetaData {
		fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	cond := metadata.QueryCondition{
		Condition: fCond,
	}
	resp, err := a.clientSet.CoreService().Model().ReadModelAttrUnique(context.Background(), params.Header, cond)
	if err != nil {
		blog.Errorf("[UniqueOperation] search for %s, failed %v, rid: %s", objectID, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrTopoObjectUniqueSearchFailed)
	}
	if !resp.Result {
		return nil, params.Err.New(resp.Code, resp.ErrMsg)
	}
	return resp.Data.Info, nil
}
