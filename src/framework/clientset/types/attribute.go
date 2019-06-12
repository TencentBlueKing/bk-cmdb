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

package types

import (
	"time"

	"configcenter/src/framework/core/types"
)

type AttributeType string

const (
	SingleChar AttributeType = "singlechar"
	LongChar   AttributeType = "longchar"
	Int        AttributeType = "int"
	Enum       AttributeType = "enum"
	Date       AttributeType = "date"
	Time       AttributeType = "time"
	ObjectUser AttributeType = "objuser"
	TimeZone   AttributeType = "timezone"
	Bool       AttributeType = "bool"
)

type Attribute struct {
	Description   string        `json:"description,omitempty"`
	UpdateAt      time.Time     `json:"last_time,omitempty"`
	Tenancy       string        `json:"bk_supplier_account,omitempty"`
	Name          string        `json:"bk_property_name,omitempty"`
	IsRequired    bool          `json:"isrequired,omitempty"`
	Type          AttributeType `json:"bk_property_type,omitempty"`
	Option        string        `json:"option,omitempty"`
	Creator       string        `json:"creator,omitempty"`
	ID            int64         `json:"id,omitempty"`
	ModelID       string        `json:"bk_object_id,omitempty"`
	Placeholder   string        `json:"placeholder,omitempty"`
	Editable      bool          `json:"editable,omitempty"`
	Unit          string        `json:"unit,omitempty"`
	IsSystem      bool          `json:"is_system,omitempty"`
	CreatedAt     time.Time     `json:"create_time,omitempty"`
	IsAPI         bool          `json:"bk_isapi,omitempty"`
	PropertyID    string        `json:"bk_property_id,omitempty"`
	PropertyGroup string        `json:"bk_property_group,omitempty"`
	PropertyIndex int64         `json:"bk_property_index,omitempty"`
	IsPre         bool          `json:"ispre,omitempty"`
}

type CreateAttributeCtx struct {
	BaseCtx
	Attribute Attribute
}

type CreateAttributeResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		ID int64 `json:"id"`
	} `json:"data"`
}

type DeleteAttributeCtx struct {
	BaseCtx
	AttributeID int64
}

type UpdateAttributeCtx struct {
	BaseCtx
	AttributeID int64
	Attribute   Attribute
}

type GetAttributeCtx struct {
	BaseCtx
	Filter types.MapStr
}

type GetAttributeResult struct {
	BaseResp `json:",inline"`
	Data     []Attribute `json:"data"`
}
