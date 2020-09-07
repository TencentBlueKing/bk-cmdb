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
	"configcenter/src/common/mapstr"
)

const (
	GroupFieldID              = "id"
	GroupFieldGroupID         = "bk_group_id"
	GroupFieldGroupName       = "bk_group_name"
	GroupFieldGroupIndex      = "bk_group_index"
	GroupFieldObjectID        = "bk_obj_id"
	GroupFieldSupplierAccount = "bk_supplier_account"
	GroupFieldIsDefault       = "bk_isdefault"
	GroupFieldIsPre           = "ispre"
)

// PropertyGroupObjectAtt uset to update or delete the property group object attribute
type PropertyGroupObjectAtt struct {
	Condition struct {
		OwnerID    string `field:"bk_supplier_account" json:"bk_supplier_account"`
		ObjectID   string `field:"bk_obj_id" json:"bk_obj_id"`
		PropertyID string `field:"bk_property_id" json:"bk_property_id"`
	} `json:"condition"`
	Data struct {
		PropertyGroupID string `field:"bk_property_group" json:"bk_property_group"`
		PropertyIndex   int    `field:"bk_property_index" json:"bk_property_index"`
	} `json:"data"`
}

// Group group metadata definition
type Group struct {
	BizID      int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ID         int64  `field:"id" json:"id" bson:"id"`
	GroupID    string `field:"bk_group_id" json:"bk_group_id" bson:"bk_group_id"`
	GroupName  string `field:"bk_group_name" json:"bk_group_name" bson:"bk_group_name"`
	GroupIndex int64  `field:"bk_group_index" json:"bk_group_index" bson:"bk_group_index"`
	ObjectID   string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	OwnerID    string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	IsDefault  bool   `field:"bk_isdefault" json:"bk_isdefault" bson:"bk_isdefault"`
	IsPre      bool   `field:"ispre" json:"ispre" bson:"ispre"`
	IsCollapse bool   `field:"is_collapse" json:"is_collapse" bson:"is_collapse"`
}

// Parse load the data from mapstr group into group instance
func (cli *Group) Parse(data mapstr.MapStr) (*Group, error) {

	err := mapstr.SetValueToStructByTags(cli, data)
	if nil != err {
		return nil, err
	}

	return cli, err
}

// ToMapStr to mapstr
func (cli *Group) ToMapStr() mapstr.MapStr {
	return mapstr.SetValueToMapStrByTags(cli)
}
