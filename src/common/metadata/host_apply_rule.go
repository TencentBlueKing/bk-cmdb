/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metadata

import (
	"time"
)

// HostApplyRule represent one rule of host property auto apply
type HostApplyRule struct {
	ID       int64 `field:"id" json:"id" bson:"id" mapstructure:"id"`
	BizID    int64 `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id" mapstructure:"bk_biz_id"`
	ModuleID int64 `field:"bk_module_id" json:"bk_module_id" bson:"bk_module_id" mapstructure:"bk_module_id"`
	// `id` field of table: `cc_AsstDes`, not the same with bk_property_id
	AttributeID   int64       `field:"bk_attribute_id" json:"bk_attribute_id" bson:"bk_attribute_id" mapstructure:"bk_attribute_id"`
	PropertyValue interface{} `field:"bk_property_value" json:"bk_property_value" bson:"bk_property_value" mapstructure:"bk_property_value"`

	// 通用字段
	Creator         string    `field:"creator" json:"creator" bson:"creator" mapstructure:"creator"`
	Modifier        string    `field:"modifier" json:"modifier" bson:"modifier" mapstructure:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time" bson:"create_time" mapstructure:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time" bson:"last_time" mapstructure:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account" mapstructure:"bk_supplier_account"`
}

func (h *HostApplyRule) Validate() (string, error) {
	return "", nil
}

type CreateHostApplyRuleOption struct {
	AttributeID   int64       `field:"bk_attribute_id" json:"bk_attribute_id" bson:"bk_attribute_id" mapstructure:"bk_attribute_id"`
	ModuleID      int64       `field:"bk_module_id" json:"bk_module_id" bson:"bk_module_id" mapstructure:"bk_module_id"`
	PropertyValue interface{} `field:"bk_property_value" json:"bk_property_value" bson:"bk_property_value" mapstructure:"bk_property_value"`
}

type UpdateHostApplyRuleOption struct {
	PropertyValue interface{} `field:"bk_property_value" json:"bk_property_value" bson:"bk_property_value" mapstructure:"bk_property_value"`
}

type MultipleHostApplyRuleResult struct {
	Count int64           `json:"count"`
	Info  []HostApplyRule `json:"info"`
}

type ListHostApplyRuleOption struct {
	ModuleIDs []int64  `field:"bk_module_ids" json:"bk_module_ids" bson:"bk_module_ids"`
	Page      BasePage `json:"page"`
}

type DeleteHostApplyRuleOption struct {
	RuleIDs []int64 `field:"host_apply_rule_ids" json:"host_apply_rule_ids" bson:"host_apply_rule_ids"`
}
