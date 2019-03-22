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
	"configcenter/src/auth"
	"context"
	"fmt"

	"configcenter/src/apimachinery"
	"configcenter/src/auth/extensions"
	"configcenter/src/auth/meta"
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
func NewUniqueOperation(client apimachinery.ClientSetInterface, authorize auth.Authorize) UniqueOperationInterface {
	return &unique{
		clientSet: client,
		authorize:      authorize,
	}
}

type unique struct {
	clientSet apimachinery.ClientSetInterface
	authorize         auth.Authorize
}

func (a *unique) Create(params types.ContextParams, objectID string, request *metadata.CreateUniqueRequest) (uniqueID *metadata.RspID, err error) {
	
	// auth: check authorization
	authManager := extensions.NewAuthManager(a.clientSet, a.authorize, params.Err)
	if err := authManager.AuthorizeByObjectID(params.Context, params.Header, meta.Update, objectID); err != nil {
		blog.V(2).Infof("create unique for model %s failed, authorization failed, err: %+v", objectID, err)
		return nil, err
	}
	
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
		blog.Errorf("[UniqueOperation] create for %s, %#v failed %v", objectID, request, err)
		return nil, params.Err.Error(common.CCErrTopoObjectUniqueCreateFailed)
	}
	if !resp.Result {
		return nil, params.Err.New(resp.Code, resp.ErrMsg)
	}
	
	// auth: register unique to iam
	uniqueid := int64(resp.Data.Created.ID)
	if err := authManager.UpdateRegisteredModelUniqueByID(params.Context, params.Header, uniqueid); err != nil {
		return nil, fmt.Errorf("register model attribute unique to iam failed, err: %+v", err)
	}
	return &metadata.RspID{ID: int64(resp.Data.Created.ID)}, nil
}

func (a *unique) Update(params types.ContextParams, objectID string, id uint64, request *metadata.UpdateUniqueRequest) (err error) {

	// auth: check authorization
	authManager := extensions.NewAuthManager(a.clientSet, a.authorize, params.Err)
	if err := authManager.AuthorizeByObjectID(params.Context, params.Header, meta.Update, objectID); err != nil {
		blog.V(2).Infof("update unique %d for model %s failed, authorization failed, err: %+v", id, objectID, err)
		return err
	}

	update := metadata.UpdateModelAttrUnique{
		Data: *request,
	}
	resp, err := a.clientSet.CoreService().Model().UpdateModelAttrUnique(context.Background(), params.Header, objectID, id, update)
	if err != nil {
		blog.Errorf("[UniqueOperation] update for %s, %d, %#v failed %v", objectID, id, request, err)
		return params.Err.Error(common.CCErrTopoObjectUniqueUpdateFailed)
	}
	if !resp.Result {
		return params.Err.New(resp.Code, resp.ErrMsg)
	}
	
	// auth: update register to iam
	if err := authManager.UpdateRegisteredModelUniqueByID(params.Context, params.Header, int64(id)); err != nil {
		blog.V(2).Infof("update unique %d for model %s failed, authorization failed, err: %+v", id, objectID, err)
		return err
	}
	return nil
}

func (a *unique) Delete(params types.ContextParams, objectID string, id uint64) (err error) {
	
	// auth: check authorization
	authManager := extensions.NewAuthManager(a.clientSet, a.authorize, params.Err)
	if err := authManager.AuthorizeByObjectID(params.Context, params.Header, meta.Update, objectID); err != nil {
		blog.V(2).Infof("delete unique %d for model %s failed, authorization failed, %+v", id, objectID, err)
		return err
	}
	
	resp, err := a.clientSet.CoreService().Model().DeleteModelAttrUnique(context.Background(), params.Header, objectID, id)
	if err != nil {
		blog.Errorf("[UniqueOperation] delete for %s, %d failed %v", objectID, id, err)
		return params.Err.Error(common.CCErrTopoObjectUniqueDeleteFailed)
	}
	if !resp.Result {
		return params.Err.New(resp.Code, resp.ErrMsg)
	}
	// auth: deregister to iam
	if err := authManager.DeregisterModelUniqueByID(params.Context, params.Header, int64(id)); err != nil {
		blog.V(2).Infof("deregister unique %d for model %s failed, authorization failed, err: %+v", id, objectID, err)
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
		blog.Errorf("[UniqueOperation] search for %s, %#v failed %v", objectID, err)
		return nil, params.Err.Error(common.CCErrTopoObjectUniqueSearchFailed)
	}
	if !resp.Result {
		return nil, params.Err.New(resp.Code, resp.ErrMsg)
	}
	return resp.Data.Info, nil
}
