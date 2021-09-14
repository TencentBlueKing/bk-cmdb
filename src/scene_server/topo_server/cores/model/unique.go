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
	"encoding/json"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	metadata "configcenter/src/common/metadata"
)

// Unique group opeartion interface declaration
type Unique interface {
	Operation

	Origin() metadata.ObjectUnique

	SetKeys([]metadata.UniqueKey)
	GetKeys() []metadata.UniqueKey

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
	kit       *rest.Kit
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

// Create unique
func (g *unique) Create() error {
	data := metadata.ObjectUnique{
		ObjID: g.data.ObjID,
		Keys:  g.data.Keys,
	}

	rsp, err := g.clientSet.CoreService().Model().CreateModelAttrUnique(g.kit.Ctx, g.kit.Header, g.data.ObjID,
		metadata.CreateModelAttrUnique{Data: data})
	if nil != err {
		blog.Errorf("[model-unique] failed to request object controller, err: %s, rid: %s", err.Error(), g.kit.Rid)
		return g.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	g.data.ID = uint64(rsp.Created.ID)
	return nil
}

// Update unique
func (g *unique) Update(data mapstr.MapStr) error {
	updateReq := metadata.UpdateUniqueRequest{
		Keys: g.data.Keys,
	}

	_, err := g.clientSet.CoreService().Model().UpdateModelAttrUnique(g.kit.Ctx, g.kit.Header, g.data.ObjID,
		g.data.ID, metadata.UpdateModelAttrUnique{Data: updateReq})
	if nil != err {
		blog.Errorf("[model-unique]failed to request object controller, err: %s, rid: %s", err.Error(), g.kit.Rid)
		return err
	}

	return nil
}

// Save create or update unique
func (g *unique) Save(data mapstr.MapStr) error {
	cond := condition.CreateCondition().Field(common.BKObjIDField).Eq(g.data.ObjID)
	searchResp, err := g.clientSet.CoreService().Model().ReadModelAttrUnique(g.kit.Ctx, g.kit.Header,
		metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[model-unique]failed to request object controller, err: %s, rid: %s", err.Error(), g.kit.Rid)
		return err
	}

	keyhash := g.data.KeysHash()
	var exists *metadata.ObjectUnique
	for _, uni := range searchResp.Info {
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

// IsExists check unique if exists
func (g *unique) IsExists() (bool, error) {
	cond := condition.CreateCondition().Field(common.BKObjIDField).Eq(g.data.ObjID)
	searchResp, err := g.clientSet.CoreService().Model().ReadModelAttrUnique(g.kit.Ctx, g.kit.Header,
		metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[model-unique]failed to request object controller, err: %s, rid: %s", err.Error(), g.kit.Rid)
		return false, err
	}

	keyhash := g.data.KeysHash()
	for _, uni := range searchResp.Info {
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

func (g *unique) SetKeys(keys []metadata.UniqueKey) {
	g.data.Keys = keys
}
func (g *unique) GetKeys() []metadata.UniqueKey {
	return g.data.Keys
}
