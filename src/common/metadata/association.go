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

import types "configcenter/src/common/mapstr"

const (
	// AssociationFieldObjectID the association data field definition
	AssociationFieldObjectID = "bk_obj_id"
	// AssociationFieldObjectAttributeID the association data field definition
	AssociationFieldObjectAttributeID = "bk_object_att_id"
	// AssociationFieldSupplierAccount the association data field definition
	AssociationFieldSupplierAccount = "bk_supplier_account"
	// AssociationFieldAssociationForward the association data field definition
	AssociationFieldAssociationForward = "bk_asst_forward"
	// AssociationFieldAssociationObjectID the association data field definition
	AssociationFieldAssociationObjectID = "bk_asst_obj_id"
	// AssociationFieldAssociationName the association data field definition
	AssociationFieldAssociationName = "bk_asst_name"
)

// Association define object association struct
type Association struct {
	ID               int64  `field:"id" json:"id" bson:"id"`
	ObjectID         string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	OwnerID          string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	AsstForward      string `field:"bk_asst_forward" json:"bk_asst_forward" bson:"bk_asst_forward"`
	AsstObjID        string `field:"bk_asst_obj_id" json:"bk_asst_obj_id" bson:"bk_asst_obj_id"`
	AsstName         string `field:"bk_asst_name" json:"bk_asst_name" bson:"bk_asst_name"`
	ObjectAttID      string `field:"bk_object_att_id" json:"bk_object_att_id" bson:"bk_object_att_id"`
	ClassificationID string `field:"bk_classification_id" bson:"-"`
	ObjectIcon       string `field:"bk_obj_icon" bson:"-"`
	ObjectName       string `field:"bk_obj_name" bson:"-"`
}

// Parse load the data from mapstr attribute into attribute instance
func (cli *Association) Parse(data types.MapStr) (*Association, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Association) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}

// InstAsst an association definition between instances.
type InstAsst struct {
	ID           int64  `field:"id" json:"-"`
	InstID       int64  `field:"bk_inst_id" json:"bk_inst_id"`
	ObjectID     string `field:"bk_obj_id" json:"bk_obj_id"`
	AsstInstID   int64  `field:"bk_asst_inst_id" json:"bk_asst_inst_id"`
	AsstObjectID string `field:"bk_asst_obj_id" json:"bk_asst_obj_id"`
}
type InstNameAsst struct {
	ID         string                 `json:"id"`
	ObjID      string                 `json:"bk_obj_id"`
	ObjIcon    string                 `json:"bk_obj_icon"`
	InstID     int64                  `json:"bk_inst_id"`
	ObjectName string                 `json:"bk_obj_name"`
	InstName   string                 `json:"bk_inst_name"`
	InstInfo   map[string]interface{} `json:"inst_info,omitempty"`
}

// Parse load the data from mapstr attribute into attribute instance
func (cli *InstAsst) Parse(data types.MapStr) (*InstAsst, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *InstAsst) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}

// MainlineObjectTopo the mainline object topo
type MainlineObjectTopo struct {
	ObjID      string `field:"bk_obj_id" json:"bk_obj_id"`
	ObjName    string `field:"bk_obj_name" json:"bk_obj_name"`
	OwnerID    string `field:"bk_supplier_account" json:"bk_supplier_account"`
	NextObj    string `field:"bk_next_obj" json:"bk_next_obj"`
	NextName   string `field:"bk_next_name" json:"bk_next_name"`
	PreObjID   string `field:"bk_pre_obj_id" json:"bk_pre_obj_id"`
	PreObjName string `field:"bk_pre_obj_name" json:"bk_pre_obj_name"`
}

// Parse load the data from mapstr attribute into attribute instance
func (cli *MainlineObjectTopo) Parse(data types.MapStr) (*MainlineObjectTopo, error) {

	err := SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *MainlineObjectTopo) ToMapStr() types.MapStr {
	return SetValueToMapStrByTags(cli)
}

// TopoInst 实例拓扑结构
type TopoInst struct {
	InstID   int64  `json:"bk_inst_id"`
	InstName string `json:"bk_inst_name"`
	ObjID    string `json:"bk_obj_id"`
	ObjName  string `json:"bk_obj_name"`
	Default  int    `json:"default"`
}

// TopoInstRst 拓扑实例
type TopoInstRst struct {
	TopoInst `json:",inline"`
	Child    []TopoInstRst `json:"child"`
}

// ConditionItem subcondition
type ConditionItem struct {
	Field    string      `json:"field,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

// AssociationParams  association params
type AssociationParams struct {
	Page      BasePage                   `json:"page,omitempty"`
	Fields    map[string][]string        `json:"fields,omitempty"`
	Condition map[string][]ConditionItem `json:"condition,omitempty"`
}
