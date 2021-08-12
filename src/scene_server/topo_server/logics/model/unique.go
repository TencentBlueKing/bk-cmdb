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

package model

import (
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// UniqueOperationInterface unique operation methods
type UniqueOperationInterface interface {
	// CreateUnique create a unique by objectID and the unique keys
	CreateUnique(kit *rest.Kit, objectID string, request *metadata.CreateUniqueRequest) (*metadata.RspID, error)
	// UpdateUnique update the unique specified by the objectID
	UpdateUnique(kit *rest.Kit, objectID string, id uint64, request *metadata.UpdateUniqueRequest) error
	// DeleteUnique delete the unique specified by the objectID and unique id
	DeleteUnique(kit *rest.Kit, objectID string, id uint64) error
	// SearchUnique search all unique specified by the objectID
	SearchUnique(kit *rest.Kit, objectID string) ([]metadata.ObjectUnique, error)
}

// NewUniqueOperation create a new unique operation instance
func NewUniqueOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) UniqueOperationInterface {

	return &unique{
		clientSet:   client,
		authManager: authManager,
	}
}

type unique struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

// CreateUnique create a unique by objectID and the unique keys
func (a *unique) CreateUnique(kit *rest.Kit, objectID string, request *metadata.CreateUniqueRequest) (*metadata.RspID,
	error) {

	unique := metadata.ObjectUnique{
		ObjID: request.ObjID,
		Keys:  request.Keys,
	}

	resp, err := a.clientSet.CoreService().Model().CreateModelAttrUnique(kit.Ctx, kit.Header,
		objectID, metadata.CreateModelAttrUnique{Data: unique})
	if err != nil {
		blog.Errorf("create for %s, %#v failed %v, rid: %s", objectID, request, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoObjectUniqueCreateFailed)
	}

	if err = resp.CCError(); err != nil {
		blog.Errorf("create for %s, %#v failed %v, rid: %s", objectID, request, err, kit.Rid)
		return nil, err
	}

	return &metadata.RspID{ID: int64(resp.Data.Created.ID)}, nil
}

// UpdateUnique update the unique specified by the objectID
func (a *unique) UpdateUnique(kit *rest.Kit, objectID string, id uint64, request *metadata.UpdateUniqueRequest) error {

	update := metadata.UpdateModelAttrUnique{
		Data: *request,
	}
	resp, err := a.clientSet.CoreService().Model().UpdateModelAttrUnique(kit.Ctx, kit.Header, objectID, id, update)
	if err != nil {
		blog.Errorf("update for %s, %d, %#v failed %v, rid: %s", objectID, id, request, err, kit.Rid)
		return kit.CCError.Error(common.CCErrTopoObjectUniqueUpdateFailed)
	}

	if err = resp.CCError(); err != nil {
		blog.Errorf("update for %s, %d, %#v failed %v, rid: %s", objectID, id, request, err, kit.Rid)
		return err
	}

	return nil
}

// DeleteUnique delete the unique specified by the objectID and unique id
func (a *unique) DeleteUnique(kit *rest.Kit, objectID string, id uint64) error {

	resp, err := a.clientSet.CoreService().Model().DeleteModelAttrUnique(kit.Ctx, kit.Header, objectID, id)
	if err != nil {
		blog.Errorf("delete for %s, %d failed %v, rid: %s", objectID, id, err, kit.Rid)
		return kit.CCError.Error(common.CCErrTopoObjectUniqueDeleteFailed)
	}

	if err = resp.CCError(); err != nil {
		blog.Errorf("delete for %s, %d failed %v, rid: %s", objectID, id, err, kit.Rid)
		return err
	}

	return nil
}

// SearchUnique search all unique specified by the objectID
func (a *unique) SearchUnique(kit *rest.Kit, objectID string) ([]metadata.ObjectUnique, error) {

	cond := metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKObjIDField: objectID,
		},
	}
	resp, err := a.clientSet.CoreService().Model().ReadModelAttrUnique(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("search for %s, failed %v, rid: %s", objectID, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoObjectUniqueSearchFailed)
	}

	if err = resp.CCError(); err != nil {
		blog.Errorf("search for %s, failed %v, rid: %s", objectID, err, kit.Rid)
		return nil, err
	}

	return resp.Data.Info, nil
}
