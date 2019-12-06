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

	"configcenter/src/common/errors"
	"configcenter/src/common/querybuilder"
)

const (
	TopoNodeKeyword = "keyword"
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
	Count int64           `json:"count" mapstructure:"count"`
	Info  []HostApplyRule `json:"info" mapstructure:"info"`
}

type ListHostApplyRuleOption struct {
	ModuleIDs []int64  `field:"bk_module_ids" json:"bk_module_ids" bson:"bk_module_ids" mapstructure:"bk_module_ids"`
	Page      BasePage `json:"page" mapstructure:"page"`
}

type ListHostRelatedApplyRuleOption struct {
	HostIDs []int64  `field:"bk_host_ids" json:"bk_host_ids" bson:"bk_host_ids" mapstructure:"bk_host_ids"`
	Page    BasePage `json:"page" mapstructure:"page"`
}

type DeleteHostApplyRuleOption struct {
	RuleIDs []int64 `field:"host_apply_rule_ids" json:"host_apply_rule_ids" bson:"host_apply_rule_ids" mapstructure:"host_apply_rule_ids"`
}

type BatchCreateOrUpdateApplyRuleOption struct {
	Rules []CreateOrUpdateApplyRuleOption `field:"host_apply_rules" json:"host_apply_rules" bson:"host_apply_rules" mapstructure:"host_apply_rules"`
}

type CreateOrUpdateApplyRuleOption struct {
	AttributeID   int64       `field:"bk_attribute_id" json:"bk_attribute_id" bson:"bk_attribute_id" mapstructure:"bk_attribute_id"`
	ModuleID      int64       `field:"bk_module_id" json:"bk_module_id" bson:"bk_module_id" mapstructure:"bk_module_id"`
	PropertyValue interface{} `field:"bk_property_value" json:"bk_property_value" bson:"bk_property_value" mapstructure:"bk_property_value"`
}

type BatchCreateOrUpdateHostApplyRuleResult struct {
	Items []CreateOrUpdateHostApplyRuleResult `json:"items" mapstructure:"items"`
}

type CreateOrUpdateHostApplyRuleResult struct {
	ErrorContainer `json:",inline"`
	Index          int           `json:"index" mapstructure:"index"`
	Rule           HostApplyRule `json:"rule" mapstructure:"rule"`
}

// ConflictResolver 定义单个冲突的解决办法
type HostApplyConflictResolver struct {
	HostID        int64       `field:"bk_host_id" json:"bk_host_id" bson:"bk_host_id" mapstructure:"bk_host_id"`
	AttributeID   int64       `field:"bk_attribute_id" json:"bk_attribute_id" bson:"bk_attribute_id" mapstructure:"bk_attribute_id"`
	PropertyValue interface{} `field:"bk_property_value" json:"bk_property_value" bson:"bk_property_value" mapstructure:"bk_property_value"`
}

type Host2Modules struct {
	HostID    int64   `field:"bk_host_id" json:"bk_host_id" bson:"bk_host_id" mapstructure:"bk_host_id"`
	ModuleIDs []int64 `field:"bk_module_ids" json:"bk_module_ids" bson:"bk_module_ids" mapstructure:"bk_module_ids"`
}

// 主机属性自动应用执行计划生成逻辑核心数据结构
// 设计背景：该数据结构需要支持如下三种场景
// 1. 应用模块配置到主机属性
// 2. 编辑模块配置(可能未保存), 预览应用效果(查看是否有冲突)
// 3. 将主机转移到模块下前预览应用效果（查看是否有冲突）
// 字段说明
// - Rules: 主机属性应用规则，由于上述case2的存在，其中 ID 可能为0
// - HostModules: 主机所有模块信息，case3的存在，导致不能直接从db中查询主机所属模块
// - ConflictResolvers: 可选参数，用于表示主机属性应用出现冲突时，如何设置应用值，如果未设置则冲突的字段不会被更新
type HostApplyPlanOption struct {
	Rules             []HostApplyRule             `field:"host_apply_rules" json:"host_apply_rules" bson:"host_apply_rules" mapstructure:"host_apply_rules"`
	HostModules       []Host2Modules              `field:"host_modules" json:"host_modules" bson:"host_modules" mapstructure:"host_modules"`
	ConflictResolvers []HostApplyConflictResolver `field:"conflict_resolvers" json:"conflict_resolvers" bson:"conflict_resolvers" mapstructure:"conflict_resolvers"`
}

type HostApplyConflictField struct {
	AttributeID   int64       `field:"bk_attribute_id" json:"bk_attribute_id" bson:"bk_attribute_id" mapstructure:"bk_attribute_id"`
	PropertyID    string      `field:"bk_property_id" json:"bk_property_id" mapstructure:"bk_property_id"`
	PropertyValue interface{} `field:"bk_property_value" json:"bk_property_value" mapstructure:"bk_property_value"`

	Rules []HostApplyRule `field:"host_apply_rules" json:"host_apply_rules" bson:"host_apply_rules" mapstructure:"host_apply_rules"`
	// UnresolvedConflictExist show whether conflict still exist after use possible conflict resolver
	// if there is a conflict, but has a resolver for it, ConflictedStillExist will be false
	UnresolvedConflictExist bool `field:"unresolved_conflict_exist" json:"unresolved_conflict_exist" mapstructure:"unresolved_conflict_exist"`
}

type HostApplyUpdateField struct {
	AttributeID   int64       `field:"bk_attribute_id" json:"bk_attribute_id" bson:"bk_attribute_id" mapstructure:"bk_attribute_id"`
	PropertyID    string      `field:"bk_property_id" json:"bk_property_id" mapstructure:"bk_property_id"`
	PropertyValue interface{} `field:"bk_property_value" json:"bk_property_value" mapstructure:"bk_property_value"`
}

type OneHostApplyPlan struct {
	ErrorContainer `json:",inline"`
	HostID         int64     `field:"bk_host_id" json:"bk_host_id" bson:"bk_host_id" mapstructure:"bk_host_id"`
	CloudInfo      CloudInst `field:"cloud_area" json:"cloud_area" bson:"cloud_area" mapstructure:"cloud_area"`
	ModuleIDs      []int64   `field:"bk_module_ids" json:"bk_module_ids" bson:"bk_module_ids" mapstructure:"bk_module_ids"`
	// 预计执行后端主机信息
	ExpiredHost    map[string]interface{}   `field:"expired_host" json:"expired_host" bson:"expired_host" mapstructure:"expired_host"`
	UpdateFields   []HostApplyUpdateField   `field:"update_fields" json:"update_fields" bson:"update_fields" mapstructure:"update_fields"`
	ConflictFields []HostApplyConflictField `field:"conflicts" json:"conflicts" bson:"conflicts" mapstructure:"conflicts"`
	// 未解决的冲突字段数
	UnresolvedConflictCount int64 `field:"unresolved_conflict_count" json:"unresolved_conflict_count" mapstructure:"unresolved_conflict_count"`
}

func (plan OneHostApplyPlan) GetUpdateData() map[string]interface{} {
	updateData := map[string]interface{}{}
	for _, field := range plan.UpdateFields {
		updateData[field.PropertyID] = field.PropertyValue
	}
	return updateData
}

type HostApplyPlanResult struct {
	Plans []OneHostApplyPlan `field:"plans" json:"plans" bson:"plans" mapstructure:"plans"`
	// 未解决的冲突主机数
	UnresolvedConflictCount int64           `field:"unresolved_conflict_count" json:"unresolved_conflict_count" mapstructure:"unresolved_conflict_count"`
	HostAttributes          []Attribute     `field:"host_attributes" json:"host_attributes" bson:"host_attributes" mapstructure:"host_attributes"`
	Count                   int             `field:"count" json:"count" bson:"count" mapstructure:"count"`
	Rules                   []HostApplyRule `field:"final_rules" json:"final_rules" mapstructure:"final_rules"`
}

type HostApplyPlanRequest struct {
	RemoveRuleIDs     []int64                     `field:"remove_rule_ids" json:"remove_rule_ids" bson:"remove_rule_ids" mapstructure:"remove_rule_ids"`
	IgnoreRuleIDs     []int64                     `field:"ignore_rule_ids" json:"ignore_rule_ids" bson:"ignore_rule_ids" mapstructure:"ignore_rule_ids"`
	AdditionalRules   []CreateHostApplyRuleOption `field:"additional_rules" json:"additional_rules" bson:"additional_rules" mapstructure:"additional_rules"`
	ConflictResolvers []HostApplyConflictResolver `field:"conflict_resolvers" json:"conflict_resolvers" bson:"conflict_resolvers" mapstructure:"conflict_resolvers"`
	ModuleIDs         []int64                     `field:"bk_module_ids" json:"bk_module_ids" bson:"bk_module_ids" mapstructure:"bk_module_ids"`
	// optional, if set, only hostID in HostIDs will be used
	HostIDs []int64 `field:"bk_host_ids" json:"bk_host_ids" bson:"bk_host_ids" mapstructure:"bk_host_ids"`
}

type HostApplyResult struct {
	ErrorContainer `json:",inline"`
	HostID         int64 `field:"bk_host_id" json:"bk_host_id" bson:"bk_host_id" mapstructure:"bk_host_id"`
}

type ErrorContainer struct {
	ErrCode int    `json:"error_code" json:"error_code" bson:"error_code" mapstructure:"error_code"`
	ErrMsg  string `json:"error_message" json:"error_message" bson:"error_message" mapstructure:"error_message"`
}

func (container *ErrorContainer) SetError(err errors.CCErrorCoder) {
	if err == nil {
		return
	}
	container.ErrCode = err.GetCode()
	container.ErrMsg = err.Error()
}

func (container *ErrorContainer) GetError() errors.CCErrorCoder {
	if container.ErrCode == 0 {
		return nil
	}
	return errors.New(container.ErrCode, container.ErrMsg)
}

type SearchRuleRelatedModulesOption struct {
	Keyword     string                    `json:"keyword" field:"keyword" mapstructure:"keyword"`
	QueryFilter *querybuilder.QueryFilter `json:"query_filter" field:"query_filter" mapstructure:"query_filter"`
}

type UpdateModuleHostApplyEnableStatusOption struct {
	Enable     bool `json:"enable" mapstructure:"enable"`
	ClearRules bool `json:"clear_rules" mapstructure:"clear_rules"`
}
