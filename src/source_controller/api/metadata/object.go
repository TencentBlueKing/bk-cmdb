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
)

// ObjectAsst define object association struct
type ObjectAsst struct {
	ID          int       `bson:"id"                    json:"id"`
	ObjectID    string    `bson:"bk_obj_id"                json:"bk_obj_id"`
	ObjectAttID string    `bson:"bk_object_att_id"         json:"bk_object_att_id"`
	OwnerID     string    `bson:"bk_supplier_account"      json:"bk_supplier_account"`
	AsstForward string    `bson:"bk_asst_forward"          json:"bk_asst_forward"`
	AsstObjID   string    `bson:"bk_asst_obj_id"           json:"bk_asst_obj_id"`
	AsstName    string    `bson:"bk_asst_name"             json:"bk_asst_name"`
	Page        *BasePage `bson:"-"                        json:"page,omitempty"`
}

// TableName return the table name
func (ObjectAsst) TableName() string {
	return "cc_ObjAsst"
}

// ObjectAttDes define the object attribute struct
type ObjectAttDes struct {
	ID            int         `bson:"id"                     json:"id"`
	OwnerID       string      `bson:"bk_supplier_account"    json:"bk_supplier_account"`
	ObjectID      string      `bson:"bk_obj_id"              json:"bk_obj_id"`
	PropertyID    string      `bson:"bk_property_id"         json:"bk_property_id"`
	PropertyName  string      `bson:"bk_property_name"       json:"bk_property_name"`
	PropertyGroup string      `bson:"bk_property_group"      json:"bk_property_group"`
	PropertyIndex int         `bson:"bk_property_index"      json:"bk_property_index"`
	Unit          string      `bson:"unit"                   json:"unit"`
	Placeholder   string      `bson:"placeholder"            json:"placeholder"`
	Editable      bool        `bson:"editable"               json:"editable"`
	IsPre         bool        `bson:"ispre"                  json:"ispre"`
	IsRequired    bool        `bson:"isrequired"             json:"isrequired"`
	IsReadOnly    bool        `bson:"isreadonly"             json:"isreadonly"`
	IsOnly        bool        `bson:"isonly"                 json:"isonly"`
	IsSystem      bool        `bson:"bk_issystem"            json:"bk_issystem"`
	IsAPI         bool        `bson:"bk_isapi"               json:"bk_isapi"`
	PropertyType  string      `bson:"bk_property_type"       json:"bk_property_type"`
	Option        interface{} `bson:"option"                 json:"option"`
	Description   string      `bson:"description"            json:"description"`
	Creator       string      `bson:"creator"                json:"creator"`
	CreateTime    *time.Time  `bson:"create_time"            json:"create_time"`
	LastTime      *time.Time  `bson:"last_time"              json:"last_time"`
	Page          *BasePage   `bson:"-"                      json:"page,omitempty"`
}

// TableName return the table name
func (ObjectAttDes) TableName() string {
	return "cc_ObjAttDes"
}

// ObjectDes define Object struct
type ObjectDes struct {
	ID          int        `bson:"id"                json:"id"`
	ObjCls      string     `bson:"bk_classification_id" json:"bk_classification_id"`
	ObjIcon     string     `bson:"bk_obj_icon"          json:"bk_obj_icon"`
	ObjectID    string     `bson:"bk_obj_id"            json:"bk_obj_id"`
	ObjectName  string     `bson:"bk_obj_name"          json:"bk_obj_name"`
	IsPre       bool       `bson:"ispre"             json:"ispre"`
	IsPaused    bool       `bson:"bk_ispaused"          json:"bk_ispaused"`
	Position    string     `bson:"position"          json:"position"`
	OwnerID     string     `bson:"bk_supplier_account"  json:"bk_supplier_account"`
	Description string     `bson:"description"       json:"description"`
	Creator     string     `bson:"creator"           json:"creator"`
	Modifier    string     `bson:"modifier"          json:"modifier"`
	CreateTime  *time.Time `bson:"create_time"       json:"create_time"`
	LastTime    *time.Time `bson:"last_time"         json:"last_time"`
	Page        *BasePage  `bson:"-"                    json:"page,omitempty"`
}

// TableName return the table name
func (ObjectDes) TableName() string {
	return "cc_ObjDes"
}

// ObjClassification 模型分类
type ObjClassification struct {
	ID                 int       `bson:"id"                      json:"id"`
	ClassificationID   string    `bson:"bk_classification_id"    json:"bk_classification_id"`
	ClassificationName string    `bson:"bk_classification_name"  json:"bk_classification_name"`
	ClassificationType string    `bson:"bk_classification_type"  json:"bk_classification_type"`
	ClassificationIcon string    `bson:"bk_classification_icon"  json:"bk_classification_icon"`
	OwnerID            string    `bson:"bk_supplier_account"     json:"bk_supplier_account"`
	Page               *BasePage `bson:"-"                       json:"page,omitempty"`
}

// ObjClassificationObject define the class object class
type ObjClassificationObject struct {
	ObjClassification `bson:",inline"`
	Objects           []ObjectDes              `bson:"-" json:"bk_objects"`
	AsstObjects       map[string][]interface{} `bson:"-" json:"bk_asst_objects"`
}

// TableName return the table name
func (ObjClassification) TableName() string {
	return "cc_ObjClassification"
}

/** TODO: delete
// ResourceCls define the resource class
type ResourceCls struct {
	ID          int       `bson:"ID,omitempty"`
	ClsName     string    `bson:"ClsName,omitempty"`
	Description string    `bson:"Description,omitempty"`
	Page        *BasePage `bson:"-" json:"Page,omitempty"`
}

// TableName return the table name
func (ResourceCls) TableName() string {
	return "cc_ResCls"
}
*/

// PropertyGroup 属性分组结构定义
type PropertyGroup struct {
	ID         int       `bson:"id"                  json:"id"`
	GroupID    string    `bson:"bk_group_id"            json:"bk_group_id"`
	GroupName  string    `bson:"bk_group_name"          json:"bk_group_name"`
	GroupIndex int       `bson:"bk_group_index"         json:"bk_group_index"`
	ObjectID   string    `bson:"bk_obj_id"              json:"bk_obj_id"`
	OwnerID    string    `bson:"bk_supplier_account"    json:"bk_supplier_account"`
	IsDefault  bool      `bson:"bk_isdefault"           json:"bk_isdefault"`
	IsPre      bool      `bson:"ispre"               json:"ispre"`
	Page       *BasePage `bson:"-"                      json:"page,omitempty"`
}

// TableName return the table name
func (PropertyGroup) TableName() string {
	return "cc_PropertyGroup"
}

// InstAsst an association definition between instances.
type InstAsst struct {
	ID           int64  `bson:"id" json:"-"`
	InstID       int64  `bson:"bk_inst_id" json:"bk_inst_id"`
	ObjectID     string `bson:"bk_obj_id" json:"bk_obj_id"`
	AsstInstID   int64  `bson:"bk_asst_inst_id" json:"bk_asst_inst_id"`
	AsstObjectID string `bson:"bk_asst_obj_id" json:"bk_asst_obj_id"`
}

// TableName return the table name
func (InstAsst) TableName() string {
	return "cc_InstAsst"
}
