// Package metadata TODO
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
	"fmt"
	"sort"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/json"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

const (
	// TopoNodeKeyword TODO
	TopoNodeKeyword = "keyword"
)

// HostApplyRule represent one rule of host property auto apply
type HostApplyRule struct {
	ID       int64 `field:"id" json:"id" bson:"id" mapstructure:"id"`
	BizID    int64 `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id" mapstructure:"bk_biz_id"`
	ModuleID int64 `field:"bk_module_id" json:"bk_module_id" bson:"bk_module_id" mapstructure:"bk_module_id"`
	// NOCC:tosa/linelength(忽略长度)
	ServiceTemplateID int64 `field:"service_template_id" json:"service_template_id" bson:"service_template_id" mapstructure:"service_template_id"`

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

// Validate TODO
func (h *HostApplyRule) Validate() (string, error) {
	return "", nil
}

// CreateHostApplyRuleOption create host auto-apply rules.
type CreateHostApplyRuleOption struct {
	ModuleID          int64       `json:"bk_module_id,omitempty"`
	ServiceTemplateID int64       `json:"service_template_id,omitempty"`
	AttributeID       int64       `json:"bk_attribute_id"`
	PropertyValue     interface{} `json:"bk_property_value"`
}

// UpdateHostApplyRuleOption TODO
type UpdateHostApplyRuleOption struct {
	PropertyValue interface{} `field:"bk_property_value" json:"bk_property_value" bson:"bk_property_value" mapstructure:"bk_property_value"`
}

// MultipleHostApplyRuleResult TODO
type MultipleHostApplyRuleResult struct {
	Count int64           `json:"count" mapstructure:"count"`
	Info  []HostApplyRule `json:"info" mapstructure:"info"`
}

// ListHostApplyRuleOption TODO
type ListHostApplyRuleOption struct {
	ApplicationID      int64    `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ModuleIDs          []int64  `json:"bk_module_ids" bson:"bk_module_ids"`
	ServiceTemplateIDs []int64  `json:"service_template_ids" bson:"service_template_ids"`
	AttributeIDs       []int64  `json:"bk_attribute_ids" bson:"bk_attribute_ids"`
	Page               BasePage `json:"page" bson:"page"`
}

// ListHostRelatedApplyRuleOption TODO
type ListHostRelatedApplyRuleOption struct {
	HostIDs []int64  `field:"bk_host_ids" json:"bk_host_ids" bson:"bk_host_ids" mapstructure:"bk_host_ids"`
	Page    BasePage `json:"page" mapstructure:"page"`
}

// DeleteHostApplyRuleOption delete host auto-apply rule request.
type DeleteHostApplyRuleOption struct {
	ModuleIDs          []int64 `json:"bk_module_ids"`
	ServiceTemplateIDs []int64 `json:"service_template_ids"`
	RuleIDs            []int64 `json:"host_apply_rule_ids"`
}

// ValidateServiceTemplateOption judging the validity of parameters in the service template scenario.
func (option *DeleteHostApplyRuleOption) ValidateServiceTemplateOption() (rawError errors.RawErrorInfo) {

	if len(option.RuleIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"host_apply_rule_ids"},
		}
	}

	if len(option.ModuleIDs) > 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"bk_module_ids"},
		}
	}

	if len(option.ServiceTemplateIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_template_ids"},
		}
	}

	// serviceTemplate id does not allow id 0 to exist.
	if util.InArray(0, option.ServiceTemplateIDs) {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"bk_templates_ids"},
		}
	}
	return errors.RawErrorInfo{}
}

// ValidateModuleOption judging the validity of parameters in the module scenario.
func (option *DeleteHostApplyRuleOption) ValidateModuleOption() (rawError errors.RawErrorInfo) {

	if len(option.RuleIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"host_apply_rule_ids"},
		}
	}
	if len(option.ServiceTemplateIDs) > 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"service_template_ids"},
		}
	}

	if len(option.ModuleIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_module_ids"},
		}
	}

	// module id does not allow id 0 to exist.
	if len(option.ModuleIDs) > 0 {
		if util.InArray(0, option.ModuleIDs) {
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsIsInvalid,
				Args:    []interface{}{"bk_module_ids"},
			}
		}
	}
	return errors.RawErrorInfo{}
}

// BatchCreateOrUpdateApplyRuleOption TODO
type BatchCreateOrUpdateApplyRuleOption struct {
	Rules []CreateOrUpdateApplyRuleOption `field:"host_apply_rules" json:"host_apply_rules" bson:"host_apply_rules" mapstructure:"host_apply_rules"`
}

// UpdateHostByHostApplyRuleOption TODO
type UpdateHostByHostApplyRuleOption struct {
	HostIDs []int64 `field:"bk_host_ids" json:"bk_host_ids" bson:"bk_host_ids" mapstructure:"bk_host_ids"`
}

// CreateOrUpdateApplyRuleOption create or update host auto-apply rule parameters.
type CreateOrUpdateApplyRuleOption struct {
	ModuleID          int64       `json:"bk_module_id,omitempty" bson:"bk_module_id"`
	ServiceTemplateID int64       `json:"service_template_id,omitempty" bson:"service_template_id"`
	AttributeID       int64       `json:"bk_attribute_id" bson:"bk_attribute_id"`
	PropertyValue     interface{} `json:"bk_property_value" bson:"bk_property_value"`
}

// BatchCreateOrUpdateHostApplyRuleResult TODO
type BatchCreateOrUpdateHostApplyRuleResult struct {
	Items []CreateOrUpdateHostApplyRuleResult `json:"items" mapstructure:"items"`
}

// CreateOrUpdateHostApplyRuleResult TODO
type CreateOrUpdateHostApplyRuleResult struct {
	ErrorContainer `json:",inline"`
	Index          int           `json:"index" mapstructure:"index"`
	Rule           HostApplyRule `json:"rule" mapstructure:"rule"`
}

// HostAttribute host auto-apply properties.
type HostAttribute struct {
	AttributeID   int64       `json:"bk_attribute_id" bson:"bk_attribute_id"`
	PropertyValue interface{} `json:"bk_property_value" bson:"bk_property_value"`
}

// HostApplyConflictResolver define a resolution to a single conflict.
type HostApplyConflictResolver struct {
	HostID        int64 `json:"bk_host_id" bson:"bk_host_id"`
	HostAttribute `json:",inline"`
}

// HostApplyTransRules module attribute value setting in the host transfer scenario.
type HostApplyTransRules struct {
	Changed    bool            `json:"changed"`
	FinalRules []HostAttribute `json:"final_rules"`
}

// Host2Modules TODO
type Host2Modules struct {
	HostID    int64   `field:"bk_host_id" json:"bk_host_id" bson:"bk_host_id" mapstructure:"bk_host_id"`
	ModuleIDs []int64 `field:"bk_module_ids" json:"bk_module_ids" bson:"bk_module_ids" mapstructure:"bk_module_ids"`
}

// HostApplyPlanOption 主机属性自动应用执行计划生成逻辑核心数据结构
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

// HostApplyConflictField TODO
type HostApplyConflictField struct {
	AttributeID   int64       `field:"bk_attribute_id" json:"bk_attribute_id" bson:"bk_attribute_id" mapstructure:"bk_attribute_id"`
	PropertyID    string      `field:"bk_property_id" json:"bk_property_id" mapstructure:"bk_property_id"`
	PropertyValue interface{} `field:"bk_property_value" json:"bk_property_value" mapstructure:"bk_property_value"`

	Rules []HostApplyRule `field:"host_apply_rules" json:"host_apply_rules" bson:"host_apply_rules" mapstructure:"host_apply_rules"`
	// UnresolvedConflictExist show whether conflict still exist after use possible conflict resolver
	// if there is a conflict, but has a resolver for it, ConflictedStillExist will be false
	UnresolvedConflictExist bool `field:"unresolved_conflict_exist" json:"unresolved_conflict_exist" mapstructure:"unresolved_conflict_exist"`
}

// HostApplyUpdateField TODO
type HostApplyUpdateField struct {
	AttributeID   int64       `field:"bk_attribute_id" json:"bk_attribute_id" bson:"bk_attribute_id" mapstructure:"bk_attribute_id"`
	PropertyID    string      `field:"bk_property_id" json:"bk_property_id" mapstructure:"bk_property_id"`
	PropertyValue interface{} `field:"bk_property_value" json:"bk_property_value" mapstructure:"bk_property_value"`
}

// OneHostApplyPlan TODO
type OneHostApplyPlan struct {
	ErrorContainer `json:",inline"`
	HostID         int64     `field:"bk_host_id" json:"bk_host_id" bson:"bk_host_id" mapstructure:"bk_host_id"`
	CloudInfo      CloudInst `field:"cloud_area" json:"cloud_area" bson:"cloud_area" mapstructure:"cloud_area"`
	ModuleIDs      []int64   `field:"bk_module_ids" json:"bk_module_ids" bson:"bk_module_ids" mapstructure:"bk_module_ids"`
	// 预计执行后端主机信息
	ExpectHost     map[string]interface{}   `field:"expect_host" json:"expect_host" bson:"expect_host" mapstructure:"expect_host"`
	UpdateFields   []HostApplyUpdateField   `field:"update_fields" json:"update_fields" bson:"update_fields" mapstructure:"update_fields"`
	ConflictFields []HostApplyConflictField `field:"conflicts" json:"conflicts" bson:"conflicts" mapstructure:"conflicts"`
	// 未解决的冲突字段数
	UnresolvedConflictCount int64 `field:"unresolved_conflict_count" json:"unresolved_conflict_count" mapstructure:"unresolved_conflict_count"`
}

// GetUpdateData TODO
func (plan OneHostApplyPlan) GetUpdateData() map[string]interface{} {
	updateData := map[string]interface{}{}
	for _, field := range plan.UpdateFields {
		updateData[field.PropertyID] = field.PropertyValue
	}
	return updateData
}

// GetUpdateDataStr TODO
func (plan OneHostApplyPlan) GetUpdateDataStr() string {
	fields := make([]string, len(plan.UpdateFields))

	for index, field := range plan.UpdateFields {
		value, _ := json.Marshal(field.PropertyValue)
		fields[index] = fmt.Sprintf(`"%s":%s`, field.PropertyID, string(value))
	}

	sort.Strings(fields)
	return "{" + strings.Join(fields, ",") + "}"
}

// HostApplyPlanResult TODO
type HostApplyPlanResult struct {
	Plans []OneHostApplyPlan `field:"plans" json:"plans" bson:"plans" mapstructure:"plans"`
	// 未解决的冲突主机数
	UnresolvedConflictCount int64           `field:"unresolved_conflict_count" json:"unresolved_conflict_count" mapstructure:"unresolved_conflict_count"`
	HostAttributes          []Attribute     `field:"host_attributes" json:"host_attributes" bson:"host_attributes" mapstructure:"host_attributes"`
	Count                   int             `field:"count" json:"count" bson:"count" mapstructure:"count"`
	Rules                   []HostApplyRule `field:"final_rules" json:"final_rules" mapstructure:"final_rules"`
}

// HostApplyPlanBase  host auto-apply Infrastructure
type HostApplyPlanBase struct {
	BizID int64 `json:"bk_biz_id"`

	// Changed true: indicates that the host needs to apply the attribute rules in the request; false: indicates that
	// the host does not apply the attribute rule in the request
	Changed         bool                        `json:"changed"`
	RemoveRuleIDs   []int64                     `json:"remove_rule_ids"`
	IgnoreRuleIDs   []int64                     `json:"ignore_rule_ids"`
	AdditionalRules []CreateHostApplyRuleOption `json:"additional_rules"`
	// optional, if set, only hostID in HostIDs will be used
	HostIDs []int64 `json:"bk_host_ids" bson:"bk_host_ids"`
}

// HostApplyTaskStatusOption get task status.
type HostApplyTaskStatusOption struct {
	BizID   int64    `json:"bk_biz_id"`
	TaskIDs []string `json:"task_ids"`
}

// Validate validate HostApplyTaskStatusOption
func (op *HostApplyTaskStatusOption) Validate() (rawError errors.RawErrorInfo) {

	if op.BizID == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}
	if len(op.TaskIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"task_ids"},
		}
	}
	return errors.RawErrorInfo{}
}

// HostAppyTaskInfo the execution state corresponding to the task ID.
type HostAppyTaskInfo struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}

// HostApplyTaskStatusRsp host automatic application task status query information.
type HostApplyTaskStatusRsp struct {
	BizID    int64              `json:"bk_biz_id"`
	TaskInfo []HostAppyTaskInfo `json:"task_info"`
}

// HostApplyTaskResult task result
type HostApplyTaskResult struct {
	BizID  int64  `json:"bk_biz_id"`
	TaskID string `json:"task_id"`
}

// HostApplyModulesOption request parameters automatically applied by the host in the module scenario.
type HostApplyModulesOption struct {
	HostApplyPlanBase `json:",inline"`
	// ModuleIDs Module list
	ModuleIDs []int64 `json:"bk_module_ids" bson:"bk_module_ids"`
}

func hostApplyBaseValidate(bizID int64, addRules []CreateHostApplyRuleOption, removeIDs []int64) errors.RawErrorInfo {
	if bizID == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}

	if len(addRules) == 0 && len(removeIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"additional_rules and remove_rule_ids"},
		}
	}
	return errors.RawErrorInfo{}
}

// Validate validate HostApplyModulesOption
func (op *HostApplyModulesOption) Validate() (rawError errors.RawErrorInfo) {

	if len(op.ModuleIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_module_ids"},
		}
	}
	if err := hostApplyBaseValidate(op.BizID, op.AdditionalRules, op.RemoveRuleIDs); err.ErrCode != 0 {
		return err
	}
	return errors.RawErrorInfo{}
}

// HostApplyServiceTemplateOption Request parameters automatically applied by the host in the service template scenario.
type HostApplyServiceTemplateOption struct {
	HostApplyPlanBase  `json:",inline"`
	ServiceTemplateIDs []int64 ` json:"service_template_ids" bson:"service_template_ids" `
}

// Validate validate HostApplyServiceTemplateOption
func (op *HostApplyServiceTemplateOption) Validate() (rawError errors.RawErrorInfo) {

	if len(op.ServiceTemplateIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_service_template_ids"},
		}
	}
	if err := hostApplyBaseValidate(op.BizID, op.AdditionalRules, op.RemoveRuleIDs); err.ErrCode != 0 {
		return err
	}
	return errors.RawErrorInfo{}
}

// HostApplyRsp the host automatically applies the asynchronous task response.
type HostApplyRsp struct {
	BizId  int    `json:"bk_biz_id"`
	TaskId string `json:"task_id"`
}

// HostApplyPlanRequest TODO
type HostApplyPlanRequest struct {
	ApplicationID      int64                       `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	RemoveRuleIDs      []int64                     `field:"remove_rule_ids" json:"remove_rule_ids" bson:"remove_rule_ids" mapstructure:"remove_rule_ids"`
	IgnoreRuleIDs      []int64                     `field:"ignore_rule_ids" json:"ignore_rule_ids" bson:"ignore_rule_ids" mapstructure:"ignore_rule_ids"`
	AdditionalRules    []CreateHostApplyRuleOption `field:"additional_rules" json:"additional_rules" bson:"additional_rules" mapstructure:"additional_rules"`
	ConflictResolvers  []HostApplyConflictResolver `field:"conflict_resolvers" json:"conflict_resolvers" bson:"conflict_resolvers" mapstructure:"conflict_resolvers"`
	ModuleIDs          []int64                     `field:"bk_module_ids" json:"bk_module_ids" bson:"bk_module_ids" mapstructure:"bk_module_ids"`
	ServiceTemplateIDs []int64                     `field:"service_template_ids" json:"service_template_ids" bson:"service_template_ids"`
	// optional, if set, only hostID in HostIDs will be used
	HostIDs []int64 `field:"bk_host_ids" json:"bk_host_ids" bson:"bk_host_ids" mapstructure:"bk_host_ids"`
}

// HostApplyResult TODO
type HostApplyResult struct {
	ErrorContainer `json:",inline"`
	HostID         int64 `field:"bk_host_id" json:"bk_host_id" bson:"bk_host_id" mapstructure:"bk_host_id"`
}

// MultipleHostApplyResult TODO
type MultipleHostApplyResult struct {
	ErrorContainer `json:",inline"`
	HostResults    []HostApplyResult `field:"host_results" json:"host_results" bson:"host_results" mapstructure:"host_results"`
}

// ErrorContainer TODO
type ErrorContainer struct {
	ErrCode int    `json:"error_code" json:"error_code" bson:"error_code" mapstructure:"error_code"`
	ErrMsg  string `json:"error_message" json:"error_message" bson:"error_message" mapstructure:"error_message"`
}

// SetError TODO
func (container *ErrorContainer) SetError(err errors.CCErrorCoder) {
	if err == nil {
		return
	}
	container.ErrCode = err.GetCode()
	container.ErrMsg = err.Error()
}

// GetError TODO
func (container *ErrorContainer) GetError() errors.CCErrorCoder {
	if container.ErrCode == 0 {
		return nil
	}
	return errors.New(container.ErrCode, container.ErrMsg)
}

// SearchRuleRelatedModulesOption TODO
type SearchRuleRelatedModulesOption struct {
	QueryFilter *querybuilder.QueryFilter `json:"query_filter" field:"query_filter" mapstructure:"query_filter"`
}

// UpdateHostApplyEnableStatusOption update host auto-applied status request.
type UpdateHostApplyEnableStatusOption struct {

	// IDs the module scene represents the list of module IDs. In the template scene,this ID is the list of template IDs.
	IDs        []int64 `json:"ids" mapstructure:"ids"`
	Enable     bool    `json:"enabled" mapstructure:"enabled"`
	ClearRules bool    `json:"clear_rules" mapstructure:"clear_rules"`
}

// Validate validate UpdateHostApplyEnableStatusOption
func (op *UpdateHostApplyEnableStatusOption) Validate() (rawError errors.RawErrorInfo) {

	if len(op.IDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"ids"},
		}
	}

	if util.InArray(0, op.IDs) {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"ids"},
		}

	}
	return errors.RawErrorInfo{}
}

// GetHostApplyStatusParam get service template host apply status param
type GetHostApplyStatusParam struct {
	ApplicationID int64   `json:"bk_biz_id"`
	ModuleIDs     []int64 `json:"bk_module_ids"`
}

// HostApplyStatusResult host apply status result
type HostApplyStatusResult struct {
	ModuleID         int64 `json:"bk_module_id"`
	HostApplyEnabled bool  `json:"host_apply_enabled"`
}

// RuleRelatedServiceTemplateOption rule related service template option
type RuleRelatedServiceTemplateOption struct {
	ApplicationID int64                     `json:"bk_biz_id"`
	QueryFilter   *querybuilder.QueryFilter `json:"query_filter"`
}

// InvalidHostCountOption request parameters struct about invalid host count
type InvalidHostCountOption struct {
	ApplicationID int64 `json:"bk_biz_id"`
	ID            int64 `json:"id"`
}

// InvalidHostCountResult the result struct about invalid host count
type InvalidHostCountResult struct {
	Count int64 `json:"count"`
}

// HostApplyRuleCountOption service template host apply rule count option
type HostApplyRuleCountOption struct {
	ApplicationID      int64   `json:"bk_biz_id"`
	ServiceTemplateIDs []int64 `json:"service_template_ids"`
}

// HostApplyRuleCountResult service template host apply rule count
type HostApplyRuleCountResult struct {
	ServiceTemplateID int64 `json:"service_template_id"`
	Count             int64 `json:"count"`
}

// ModuleFinalRulesParam module final rules param
type ModuleFinalRulesParam struct {
	ApplicationID int64   `json:"bk_biz_id"`
	ModuleIDs     []int64 `json:"bk_module_ids"`
}

// ServiceTemplatesResponse service template response
type ServiceTemplatesResponse struct {
	BaseResp `json:",inline"`
	Data     []SrvTemplate `json:"data"`
}
