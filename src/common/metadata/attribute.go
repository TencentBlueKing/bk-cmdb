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
	"time"

	types "configcenter/src/common/mapstr"
)

const (
	AttributeFieldSupplierAccount = "bk_supplier_account"
	AttributeFieldObjectID        = "bk_obj_id"
	AttributeFieldPropertyID      = "bk_property_id"
	AttributeFieldPropertyName    = "bk_property_name"
	AttributeFieldPropertyGroup   = "bk_property_group"
	AttributeFieldPropertyIndex   = "bk_property_index"
	AttributeFieldUnit            = "unit"
	AttributeFieldPlaceHoler      = "placeholder"
	AttributeFieldIsEditable      = "editable"
	AttributeFieldIsPre           = "ispre"
	AttributeFieldIsRequired      = "isrequired"
	AttributeFieldIsReadOnly      = "isreadonly"
	AttributeFieldIsOnly          = "isonly"
	AttributeFieldIsSystem        = "bk_issystem"
	AttributeFieldIsAPI           = "bk_isapi"
	AttributeFieldPropertyType    = "bk_property_type"
	AttributeFieldOption          = "option"
	AttributeFieldDescription     = "description"
	AttributeFieldCreator         = "creator"
	AttributeFieldCreateTime      = "create_time"
	AttributeFieldLastTime        = "last_time"
)

// Attribute attribute metadata definition
type Attribute struct {
	ID            int         `field:"id"`
	OwnerID       string      `field:"bk_supplier_account"`
	ObjectID      string      `field:"bk_obj_id"`
	PropertyID    string      `field:"bk_property_id"`
	PropertyName  string      `field:"bk_property_name"`
	PropertyGroup string      `field:"bk_property_group"`
	PropertyIndex int         `field:"bk_property_index"`
	Unit          string      `field:"unit"`
	Placeholder   string      `field:"placeholder"`
	IsEditable    bool        `field:"editable"`
	IsPre         bool        `field:"ispre"`
	IsRequired    bool        `field:"isrequired"`
	IsReadOnly    bool        `field:"isreadonly"`
	IsOnly        bool        `field:"isonly"`
	IsSystem      bool        `field:"bk_issystem"`
	IsAPI         bool        `field:"bk_isapi"`
	PropertyType  string      `field:"bk_property_type"`
	Option        interface{} `field:"option"`
	Description   string      `field:"description"`
	Creator       string      `field:"creator"`
	CreateTime    *time.Time  `field:"create_time"`
	LastTime      *time.Time  `field:"last_time"`
}

// Parse load the data from mapstr attribute into attribute instance
func (cli *Attribute) Parse(data types.MapStr) (*Attribute, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Attribute) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}
