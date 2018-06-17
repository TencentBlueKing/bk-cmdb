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
	ModelFieldObjCls      = "bk_classification_id"
	ModelFieldObjIcon     = "bk_obj_icon"
	ModelFieldObjectID    = "bk_obj_id"
	ModelFieldObjectName  = "bk_obj_name"
	ModelFieldIsPre       = "ispre"
	ModelFieldIsPaused    = "bk_ispaused"
	ModelFieldPosition    = "position"
	ModelFieldOwnerID     = "bk_supplier_account"
	ModelFieldDescription = "description"
	ModelFieldCreator     = "creator"
	ModelFieldModifier    = "modifier"
	ModelFieldCreateTime  = "create_time"
	ModelFieldLastTime    = "last_time"
)

// Object object metadata definition
type Object struct {
	ID          int64      `field:"id"`
	ObjCls      string     `field:"bk_classification_id"`
	ObjIcon     string     `field:"bk_obj_icon"`
	ObjectID    string     `field:"bk_obj_id"`
	ObjectName  string     `field:"bk_obj_name"`
	IsPre       bool       `field:"ispre"`
	IsPaused    bool       `field:"bk_ispaused"`
	Position    string     `field:"position"`
	OwnerID     string     `field:"bk_supplier_account"`
	Description string     `field:"description"`
	Creator     string     `field:"creator"`
	Modifier    string     `field:"modifier"`
	CreateTime  *time.Time `field:"create_time"`
	LastTime    *time.Time `field:"last_time"`
}

// Parse load the data from mapstr object into object instance
func (cli *Object) Parse(data types.MapStr) (*Object, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Object) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}

// MainLineObject main line object definition
type MainLineObject struct {
	ObjectDes     `json:",inline"`
	AssociationID string `json:"bk_asst_obj_id"`
}

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

// ObjectAsst define object association struct
type Association struct {
	ID          int       `bson:"id"                    json:"id"`
	ObjectID    string    `bson:"bk_obj_id"                json:"bk_obj_id"`
	ObjectAttID string    `bson:"bk_object_att_id"         json:"bk_object_att_id"`
	OwnerID     string    `bson:"bk_supplier_account"      json:"bk_supplier_account"`
	AsstForward string    `bson:"bk_asst_forward"          json:"bk_asst_forward"`
	AsstObjID   string    `bson:"bk_asst_obj_id"           json:"bk_asst_obj_id"`
	AsstName    string    `bson:"bk_asst_name"             json:"bk_asst_name"`
	Page        *BasePage `bson:"-"                        json:"page,omitempty"`
}

func (cli *Association) Parse(data types.MapStr) (*Association, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

type ObjectClsDes struct {
	ID      int    `json:"id"`
	ClsID   string `json:"bk_classification_id"`
	ClsName string `json:"bk_classification_name"`
	ClsType string `json:"bk_classification_type"`
	ClsIcon string `json:"bk_classification_icon"`
}

type ObjAttDes struct {
	Attribute         `json:",inline"`
	AssoType          int    `json:"bk_asst_type"`
	AsstForward       string `json:"bk_asst_forward"`
	AssociationID     string `json:"bk_asst_obj_id"`
	PropertyGroupName string `json:"bk_property_group_name"`
}

type InstAsst struct {
	ID           int64  `bson:"id" json:"-"`
	InstID       int64  `bson:"bk_inst_id" json:"bk_inst_id"`
	ObjectID     string `bson:"bk_obj_id" json:"bk_obj_id"`
	AsstInstID   int64  `bson:"bk_asst_inst_id" json:"bk_asst_inst_id"`
	AsstObjectID string `bson:"bk_asst_obj_id" json:"bk_asst_obj_id"`
}
