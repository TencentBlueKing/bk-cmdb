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

package metadata

import (
	"fmt"
	"sort"
	"strings"

	"configcenter/src/common/mapstr"
)

// ObjectUnique TODO
type ObjectUnique struct {
	ID       uint64      `json:"id" bson:"id"`
	ObjID    string      `json:"bk_obj_id" bson:"bk_obj_id"`
	Keys     []UniqueKey `json:"keys" bson:"keys"`
	Ispre    bool        `json:"ispre" bson:"ispre"`
	OwnerID  string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime Time        `json:"last_time" bson:"last_time"`
}

// Parse load the data from mapstr attribute into ObjectUnique instance
func (cli *ObjectUnique) Parse(data mapstr.MapStr) (*ObjectUnique, error) {

	err := mapstr.SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// KeysHash TODO
func (u ObjectUnique) KeysHash() string {
	keys := []string{}
	for _, key := range u.Keys {
		keys = append(keys, fmt.Sprintf("%s:%d", key.Kind, key.ID))
	}
	sort.Strings(keys)
	return strings.Join(keys, "#")
}

// UniqueKey TODO
type UniqueKey struct {
	Kind string `json:"key_kind" bson:"key_kind"`
	ID   uint64 `json:"key_id" bson:"key_id"`
}

const (
	// UniqueKeyKindProperty TODO
	UniqueKeyKindProperty = "property"
	// UniqueKeyKindAssociation TODO
	UniqueKeyKindAssociation = "association"
)

// CreateUniqueRequest TODO
type CreateUniqueRequest struct {
	ObjID string      `json:"bk_obj_id" bson:"bk_obj_id"`
	Keys  []UniqueKey `json:"keys" bson:"keys"`
}

// CreateUniqueResult TODO
type CreateUniqueResult struct {
	BaseResp
	Data RspID `json:"data"`
}

// UpdateUniqueRequest TODO
type UpdateUniqueRequest struct {
	Keys     []UniqueKey `json:"keys" bson:"keys"`
	LastTime Time        `json:"last_time" bson:"last_time"`
}

// UpdateUniqueResult TODO
type UpdateUniqueResult struct {
	BaseResp
}

// DeleteUniqueRequest TODO
type DeleteUniqueRequest struct {
	ID    uint64 `json:"id"`
	ObjID string `json:"bk_obj_id"`
}

// DeleteUniqueResult TODO
type DeleteUniqueResult struct {
	BaseResp
}

// SearchUniqueRequest TODO
type SearchUniqueRequest struct {
	ObjID string `json:"bk_obj_id"`
}

// SearchUniqueResult TODO
type SearchUniqueResult struct {
	BaseResp
	Data []ObjectUnique `json:"data"`
}

// QueryUniqueResult TODO
type QueryUniqueResult struct {
	Count uint64         `json:"count"`
	Info  []ObjectUnique `json:"info"`
}
