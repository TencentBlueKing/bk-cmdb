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
	AttributeFieldID              = "id"
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
	ID                int64       `field:"id" json:"id" bson:"id"`
	OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
	PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
	PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
	PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
	PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
	Unit              string      `field:"unit" json:"unit" bson:"unit"`
	Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
	IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
	IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
	IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
	IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
	IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
	IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
	IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
	PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
	Option            interface{} `field:"option" json:"option" bson:"option"`
	Description       string      `field:"description" json:"description" bson:"description"`
	Creator           string      `field:"creator" json:"creator" bson:"creator"`
	CreateTime        *time.Time  `json:"create_time" bson:"creaet_time"`
	LastTime          *time.Time  `json:"last_time" bson:"last_time"`
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

// ObjAttDes 对象模型属性
type ObjAttDes struct {
	Attribute         `json:",inline"`
	AssoType          int    `json:"bk_asst_type"`
	AsstForward       string `json:"bk_asst_forward"`
	AssociationID     string `json:"bk_asst_obj_id"`
	PropertyGroupName string `json:"bk_property_group_name"`
}
