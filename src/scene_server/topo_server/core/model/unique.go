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
	"context"
	"encoding/json"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	frtypes "configcenter/src/common/mapstr"
	metadata "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// Unique group opeartion interface declaration
type Unique interface {
	Operation

	Origin() metadata.ObjectUnique

	SetKeys([]metadata.UinqueKey)
	GetKeys() []metadata.UinqueKey

	SetMustCheck(bool)
	GetMustCheck() bool

	SetObjectID(objID string)
	GetObjectID() string

	SetSupplierAccount(supplierAccount string)
	GetSupplierAccount() string

	SetIsPre(isPre bool)
	GetIsPre() bool

	SetRecordID(uint64)
	GetRecordID() uint64
}

var _ Unique = (*unique)(nil)

type unique struct {
	FieldValid
	data      metadata.ObjectUnique
	isNew     bool
	params    types.ContextParams
	clientSet apimachinery.ClientSetInterface
}

func (g *unique) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.data)
}

func (g *unique) Origin() metadata.ObjectUnique {
	return g.data
}

func (g *unique) SetObjectID(objID string) {
	g.data.ObjID = objID
}
func (g *unique) GetObjectID() string {
	return g.data.ObjID
}

func (g *unique) Create() error {
	data := metadata.CreateUniqueRequest{
		ObjID:     g.data.ObjID,
		MustCheck: g.data.MustCheck,
		Keys:      g.data.Keys,
	}
	rsp, err := g.clientSet.ObjectController().Unique().Create(context.Background(), g.params.Header, g.data.ObjID, &data)

	if nil != err {
		blog.Errorf("[model-unique] failed to request object controller, err: %s", err.Error())
		return g.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-unique] failed to create the unique(%#v), err: %s", g.data, rsp.ErrMsg)
		return g.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	g.data.ID = uint64(rsp.Data.ID)
	return nil
}

func (g *unique) Update(data frtypes.MapStr) error {
	updateReq := metadata.UpdateUniqueRequest{
		MustCheck: g.data.MustCheck,
		Keys:      g.data.Keys,
	}
	rsp, err := g.clientSet.ObjectController().Unique().Update(context.Background(), g.params.Header, g.data.ObjID, g.data.ID, &updateReq)

	if nil != err {
		blog.Errorf("[model-unique]failed to request object controller, err: %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-unique]failed to update the object %s(%d) to  %v, err: %s", g.data.ObjID, g.data.ID, updateReq, rsp.ErrMsg)
		return g.params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	return nil
}

func (g *unique) Save(data frtypes.MapStr) error {
	searchResp, err := g.clientSet.ObjectController().Unique().Search(context.Background(), g.params.Header, g.data.ObjID)
	if nil != err {
		blog.Errorf("[model-unique]failed to request object controller, err: %s", err.Error())
		return err
	}

	if !searchResp.Result {
		blog.Errorf("[model-unique]failed to search the object unique by %s, err: %s", g.data.ObjID, searchResp.ErrMsg)
		return g.params.Err.New(searchResp.Code, searchResp.ErrMsg)
	}

	keyhash := g.data.KeysHash()
	var exists *metadata.ObjectUnique
	for _, uni := range searchResp.Data {
		if uni.KeysHash() == keyhash {
			exists = &uni
			break
		}
	}

	if exists != nil {
		g.data.ID = exists.ID
		return g.Update(data)
	}
	return g.Create()
}

func (g *unique) IsExists() (bool, error) {
	searchResp, err := g.clientSet.ObjectController().Unique().Search(context.Background(), g.params.Header, g.data.ObjID)
	if nil != err {
		blog.Errorf("[model-unique]failed to request object controller, err: %s", err.Error())
		return false, err
	}

	if !searchResp.Result {
		blog.Errorf("[model-unique]failed to search the object unique by %s, err: %s", g.data.ObjID, searchResp.ErrMsg)
		return false, g.params.Err.New(searchResp.Code, searchResp.ErrMsg)
	}

	keyhash := g.data.KeysHash()
	for _, uni := range searchResp.Data {
		if uni.KeysHash() == keyhash {
			return true, nil
		}
	}

	return false, nil
}

func (g *unique) SetRecordID(recordID uint64) {
	g.data.ID = recordID
}

func (g *unique) GetRecordID() uint64 {
	return g.data.ID
}

func (g *unique) SetSupplierAccount(supplierAccount string) {
	g.data.OwnerID = supplierAccount
}

func (g *unique) GetSupplierAccount() string {
	return g.data.OwnerID
}

func (g *unique) SetIsPre(isPre bool) {
	g.data.Ispre = isPre
}

func (g *unique) GetIsPre() bool {
	return g.data.Ispre
}

func (g *unique) SetKeys(keys []metadata.UinqueKey) {
	g.data.Keys = keys
}
func (g *unique) GetKeys() []metadata.UinqueKey {
	return g.data.Keys
}
func (g *unique) SetMustCheck(mustcheck bool) {
	g.data.MustCheck = mustcheck
}
func (g *unique) GetMustCheck() bool {
	return g.data.MustCheck
}
