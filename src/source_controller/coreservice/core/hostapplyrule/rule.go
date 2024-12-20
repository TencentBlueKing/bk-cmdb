/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hostapplyrule

import (
	"context"
	"strconv"
	"strings"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/thirdparty/hooks"
)

type ruleType string

const (
	module          ruleType = "module"
	serviceTemplate ruleType = "serviceTemplate"
)

type hostApplyRule struct {
	dependence HostApplyDependence
	cs         apimachinery.ClientSetInterface
}

// HostApplyDependence TODO
type HostApplyDependence interface {
	UpdateModelInstance(kit *rest.Kit, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
}

// New TODO
func New(dependence HostApplyDependence, cs apimachinery.ClientSetInterface) core.HostApplyRuleOperation {
	rule := &hostApplyRule{
		dependence: dependence,
		cs:         cs,
	}
	return rule
}

func (p *hostApplyRule) validateID(kit *rest.Kit, bizID int64, moduleID int64,
	serviceTemplateID int64) errors.CCErrorCoder {

	if moduleID != 0 && serviceTemplateID != 0 {
		blog.Errorf("bk_module_id and service_template_id can not exist together, rid: %s", kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_module_id and service_template_id")
	}

	if moduleID == 0 && serviceTemplateID == 0 {
		blog.Errorf("bk_module_id or service_template_id no exist, rid: %s", kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_module_id and service_template_id")
	}

	if moduleID != 0 {
		modFilter := map[string]interface{}{
			common.BKAppIDField:    bizID,
			common.BKModuleIDField: moduleID,
		}
		moduleCount, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Find(modFilter).
			Count(kit.Ctx)
		if err != nil {
			blog.Errorf("count module failed, filter: %+v, err: %v, rid: %s", modFilter, err, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		if moduleCount == 0 {
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
		}
		return nil
	}

	tempFilter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKFieldID:    serviceTemplateID,
	}
	templateCount, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceTemplate).Find(tempFilter).
		Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count service template failed, filter: %+v, err: %v, rid: %s", tempFilter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if templateCount == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}
	return nil
}

func (p *hostApplyRule) listHostAttributes(kit *rest.Kit, bizID int64, hostAttributeIDs ...int64) ([]metadata.Attribute,
	errors.CCErrorCoder) {

	filter := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{
				// business private attribute
				common.BKAppIDField: bizID,
			}, {
				// global attribute
				common.BKAppIDField: 0,
			},
		},
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: hostAttributeIDs,
		},
	}
	attributes := make([]metadata.Attribute, 0)
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAttDes).Find(filter).All(kit.Ctx, &attributes)
	if err != nil {
		if mongodb.IsNotFoundError(err) {
			blog.Errorf("get host attribute failed, not found, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
			return attributes, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		blog.Errorf("get host attribute failed, db select failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return attributes, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	return attributes, nil
}

func (p *hostApplyRule) getHostAttribute(kit *rest.Kit, bizID int64, hostAttributeID int64) (metadata.Attribute,
	errors.CCErrorCoder) {

	attribute := metadata.Attribute{}
	attributes, err := p.listHostAttributes(kit, bizID, hostAttributeID)
	if err != nil {
		blog.Errorf("listHostAttributes failed, bizID: %d, attribute: %d, err: %v, rid: %s", bizID, hostAttributeID,
			err, kit.Rid)
		return attribute, err
	}
	if len(attributes) == 0 {
		return attribute, kit.CCError.CCError(common.CCErrCommNotFound)
	}
	if len(attributes) > 1 {
		return attribute, kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}
	return attributes[0], nil
}

// CreateHostApplyRule create host apply rule
func (p *hostApplyRule) CreateHostApplyRule(kit *rest.Kit, bizID int64,
	option metadata.CreateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder) {

	now := time.Now()
	rule := metadata.HostApplyRule{
		ID:                0,
		BizID:             bizID,
		AttributeID:       option.AttributeID,
		ModuleID:          option.ModuleID,
		ServiceTemplateID: option.ServiceTemplateID,
		PropertyValue:     option.PropertyValue,
		Creator:           kit.User,
		Modifier:          kit.User,
		CreateTime:        now,
		LastTime:          now,
		TenantID:          kit.TenantID,
	}
	if key, err := rule.Validate(); err != nil {
		blog.Errorf("host apply rule parameter invalid, key: %s, err: %v, rid: %s", key, err, kit.Rid)
		return rule, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, key)
	}

	// validate relation id
	if err := p.validateID(kit, bizID, rule.ModuleID, rule.ServiceTemplateID); err != nil {
		blog.Errorf("validate relation id failed, bizID: %d, err: %v, rid: %s", bizID, err, kit.Rid)
		return rule, err
	}

	attribute, ccErr := p.getHostAttribute(kit, bizID, rule.AttributeID)
	if ccErr != nil {
		blog.Errorf("get host attribute failed, bizID: %d, attributeID: %d, err: %v, rid: %s", bizID, rule.AttributeID,
			ccErr, kit.Rid)
		return rule, ccErr
	}

	if value, ok := option.PropertyValue.(string); ok {
		option.PropertyValue = strings.TrimSpace(value)
	}
	rawError := attribute.Validate(kit.Ctx, option.PropertyValue, common.BKPropertyValueField)
	if rawError.ErrCode != 0 {
		ccErr := rawError.ToCCError(kit.CCError)
		blog.Errorf("validate host attribute value failed, attribute: %+v, value: %+v, err: %v, rid: %s",
			attribute, option.PropertyValue, ccErr, kit.Rid)
		return rule, ccErr
	}

	if err := hooks.ValidHostApplyStatusHook(kit, p.cs, attribute.PropertyID, option.PropertyValue); err != nil {
		return rule, err
	}

	// generate id field
	id, err := mongodb.Shard(kit.SysShardOpts()).NextSequence(kit.Ctx, common.BKTableNameHostApplyRule)
	if err != nil {
		blog.Errorf("%s generate id failed, err: %v, rid: %s", common.BKTableNameHostApplyRule, err, kit.Rid)
		return rule, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	rule.ID = int64(id)

	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostApplyRule).Insert(kit.Ctx, rule); err != nil {
		if mongodb.IsDuplicatedError(err) {
			blog.Errorf("table %s duplicated error, doc: %+v, err: %v, rid: %s", common.BKTableNameHostApplyRule, rule,
				err, kit.Rid)
			return rule, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKAttributeIDField)
		}
		blog.Errorf("db insert failed, table: %s, doc: %+v, err: %v, rid: %s", common.BKTableNameHostApplyRule, rule,
			err, kit.Rid)
		return rule, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
	}

	return rule, nil
}

// UpdateHostApplyRule update host apply rule
func (p *hostApplyRule) UpdateHostApplyRule(kit *rest.Kit, bizID int64, ruleID int64,
	option metadata.UpdateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder) {

	rule, ccErr := p.GetHostApplyRule(kit, bizID, ruleID)
	if ccErr != nil {
		blog.Errorf("GetHostApplyRule failed, bizID: %d, id: %d, err: %v, rid: %s", bizID, ruleID, ccErr, kit.Rid)
		return rule, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	attribute, ccErr := p.getHostAttribute(kit, bizID, rule.AttributeID)
	if ccErr != nil {
		blog.Errorf("getHostAttribute failed, bizID: %d, attributeID: %d, err: %v, rid: %s", bizID, rule.AttributeID,
			ccErr, kit.Rid)
		return rule, ccErr
	}
	if value, ok := option.PropertyValue.(string); ok {
		option.PropertyValue = strings.TrimSpace(value)
	}
	rawError := attribute.Validate(kit.Ctx, option.PropertyValue, common.BKPropertyValueField)
	if rawError.ErrCode != 0 {
		ccErr := rawError.ToCCError(kit.CCError)
		blog.Errorf("validate host attribute value failed, attribute: %+v, value: %+v, err: %v, rid: %s", attribute,
			option.PropertyValue, ccErr, kit.Rid)
		return rule, ccErr
	}

	if err := hooks.ValidHostApplyStatusHook(kit, p.cs, attribute.PropertyID, option.PropertyValue); err != nil {
		return rule, err
	}

	rule.LastTime = time.Now()
	rule.Modifier = kit.User
	rule.PropertyValue = option.PropertyValue

	filter := map[string]interface{}{
		common.BKFieldID: ruleID,
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostApplyRule).Update(kit.Ctx, filter,
		rule); err != nil {
		blog.Errorf("update host apply rule failed, filter: %+v, doc: %+v, err: %v, rid: %s", filter, rule, err,
			kit.Rid)
		return rule, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	return rule, nil
}

// DeleteHostApplyRule delete host apply rule by condition, bizID maybe 0
func (p *hostApplyRule) DeleteHostApplyRule(kit *rest.Kit, bizID int64,
	option metadata.DeleteHostApplyRuleOption) errors.CCErrorCoder {

	if len(option.RuleIDs) == 0 && len(option.ModuleIDs) == 0 && len(option.ServiceTemplateIDs) == 0 {
		blog.Errorf("HTTP request body data is not set failed, err: params is empty, option: %v, rid: %s",
			option, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommHTTPBodyEmpty)
	}

	filter := make(map[string]interface{})
	if bizID != 0 {
		filter[common.BKAppIDField] = bizID
	}
	if len(option.RuleIDs) > 0 {
		filter[common.BKFieldID] = map[string]interface{}{
			common.BKDBIN: option.RuleIDs,
		}
	}
	if len(option.ModuleIDs) > 0 {
		filter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: option.ModuleIDs,
		}
	}
	if len(option.ServiceTemplateIDs) > 0 {
		filter[common.BKServiceTemplateIDField] = map[string]interface{}{
			common.BKDBIN: option.ServiceTemplateIDs,
		}
	}

	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostApplyRule).Delete(kit.Ctx,
		filter); err != nil {
		blog.Errorf("delete host apply rules failed, err: %v, filter: %v, rid: %s", err, filter, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
	}
	return nil
}

// GetHostApplyRule get host apply rule by condition
func (p *hostApplyRule) GetHostApplyRule(kit *rest.Kit, bizID int64, ruleID int64) (metadata.HostApplyRule,
	errors.CCErrorCoder) {

	rule := metadata.HostApplyRule{}
	filter := map[string]interface{}{
		common.TenantID:     kit.TenantID,
		common.BKAppIDField: bizID,
		common.BKFieldID:    ruleID,
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostApplyRule).Find(filter).One(kit.Ctx,
		&rule); err != nil {
		if mongodb.IsNotFoundError(err) {
			blog.Errorf("find host apply rule failed, not found, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
			return rule, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		blog.Errorf("find host apply rule failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return rule, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	return rule, nil
}

// GetHostApplyRuleByAttributeID get host apply rule by attribute id
func (p *hostApplyRule) GetHostApplyRuleByAttributeID(kit *rest.Kit,
	bizID, moduleID, attributeID int64) (metadata.HostApplyRule, errors.CCErrorCoder) {

	rule := metadata.HostApplyRule{}
	filter := map[string]interface{}{
		common.TenantID:           kit.TenantID,
		common.BKAppIDField:       bizID,
		common.BKModuleIDField:    moduleID,
		common.BKAttributeIDField: attributeID,
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostApplyRule).Find(filter).One(kit.Ctx,
		&rule); err != nil {
		if mongodb.IsNotFoundError(err) {
			blog.Errorf("get host apply rule by attribute id failed, not found, filter: %+v, err: %v, rid: %s", filter,
				err, kit.Rid)
			return rule, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		blog.Errorf("get host apply rule by attribute id failed, filter: %+v, err: %v, rid: %s", filter,
			err, kit.Rid)
		return rule, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	return rule, nil
}

// ListHostApplyRule by condition, bizID maybe 0
func (p *hostApplyRule) ListHostApplyRule(kit *rest.Kit, bizID int64,
	option metadata.ListHostApplyRuleOption) (metadata.MultipleHostApplyRuleResult, errors.CCErrorCoder) {
	result := metadata.MultipleHostApplyRuleResult{}

	filter := make(map[string]interface{})
	if bizID != 0 {
		filter[common.BKAppIDField] = bizID
	}
	if len(option.ModuleIDs) != 0 {
		filter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: option.ModuleIDs,
		}
	}
	if len(option.ServiceTemplateIDs) != 0 {
		filter[common.BKServiceTemplateIDField] = map[string]interface{}{
			common.BKDBIN: option.ServiceTemplateIDs,
		}
	}

	if len(option.AttributeIDs) != 0 {
		filter[common.BKAttributeIDField] = map[string]interface{}{
			common.BKDBIN: option.AttributeIDs,
		}
	}
	query := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostApplyRule).Find(filter)
	total, err := query.Count(kit.Ctx)
	if err != nil {
		blog.Errorf("get host apply rule failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	result.Count = int64(total)

	if len(option.Page.Sort) > 0 {
		query = query.Sort(option.Page.Sort)
	}
	if option.Page.Limit > 0 {
		query = query.Limit(uint64(option.Page.Limit))
	}
	if option.Page.Start > 0 {
		query = query.Start(uint64(option.Page.Start))
	}

	rules := make([]metadata.HostApplyRule, 0)
	if err := query.All(kit.Ctx, &rules); err != nil {
		blog.Errorf("list host apply rule failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	result.Info = rules
	return result, nil
}

// SearchRuleRelatedModules 用于过滤主机应用规则相关的模块
func (p *hostApplyRule) SearchRuleRelatedModules(kit *rest.Kit, bizID int64,
	option metadata.SearchRuleRelatedModulesOption) ([]metadata.Module, errors.CCErrorCoder) {

	// 1.获取与查询条件中的属性关联的rule和attribute
	rules, attributeMap, ccErr := getRuleAndAttribute(kit, bizID, option.QueryFilter, module)
	if ccErr != nil {
		return nil, ccErr
	}

	// 如果没有rule匹配或者是小于attribute的数量，那么说明没有module是满足查询条件的
	if len(rules) == 0 || len(rules) < len(attributeMap) {
		return nil, nil
	}

	// 2. 将模块与rule进行关联
	moduleToRules, moduleIDs := getRuleRelationIDs(rules, module)

	moduleFilter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.TenantID:     kit.TenantID,
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIDs,
		},
	}
	modules := make([]metadata.Module, 0)
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Find(moduleFilter).All(kit.Ctx, &modules)
	if err != nil {
		blog.Errorf("find modules failed, filter: %+v, err: %v, rid: %s", moduleFilter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	// 3.根据匹配的规则过滤出模块
	resultModules := make([]metadata.Module, 0)
	for _, module := range modules {
		rules, exist := moduleToRules[module.ModuleID]
		if !exist {
			continue
		}

		if match(kit.Ctx, rules, attributeMap, option.QueryFilter) {
			resultModules = append(resultModules, module)
		}
	}

	return resultModules, nil
}

func getRuleAndAttribute(kit *rest.Kit, bizID int64, filter *querybuilder.QueryFilter, rType ruleType) (
	[]metadata.HostApplyRule, map[int64]metadata.Attribute, errors.CCErrorCoder) {

	attributeIDs, ccErr := getAttributeIDs(kit, filter)
	if ccErr != nil {
		return nil, nil, ccErr
	}

	if len(attributeIDs) == 0 {
		return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "query_filter")
	}

	ruleFilter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.TenantID:     kit.TenantID,
		common.BKAttributeIDField: map[string]interface{}{
			common.BKDBIN: attributeIDs,
		},
	}

	switch rType {
	case module:
		ruleFilter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBGT: 0,
		}
	case serviceTemplate:
		ruleFilter[common.BKServiceTemplateIDField] = map[string]interface{}{
			common.BKDBGT: 0,
		}
	}

	var err error
	rules := make([]metadata.HostApplyRule, 0)
	err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostApplyRule).Find(ruleFilter).All(kit.Ctx, &rules)
	if err != nil {
		blog.Errorf("find rules failed, filter: %+v, err: %v, rid: %s", ruleFilter, err, kit.Rid)
		return nil, nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	attributeFilter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: attributeIDs,
		},
	}
	attributes := make([]metadata.Attribute, 0)
	err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAttDes).Find(attributeFilter).All(kit.Ctx,
		&attributes)
	if err != nil {
		blog.Errorf("find attributes failed, filter: %+v, err: %v, rid: %s", attributeFilter, err, kit.Rid)
		return nil, nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	attributeMap := make(map[int64]metadata.Attribute)
	for _, attribute := range attributes {
		attributeMap[attribute.ID] = attribute
	}

	return rules, attributeMap, nil
}

func getAttributeIDs(kit *rest.Kit, filter *querybuilder.QueryFilter) ([]int64, errors.CCErrorCoder) {
	fields := filter.GetField()
	if fields == nil || len(fields) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "query_filter")
	}

	attributeIDs := make([]int64, len(fields))
	for index, val := range fields {
		attributeID, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "query_filter")
		}
		attributeIDs[index] = attributeID
	}

	return attributeIDs, nil
}

func getRuleRelationIDs(rules []metadata.HostApplyRule, rType ruleType) (map[int64]map[string]metadata.HostApplyRule,
	[]int64) {

	idToRules := make(map[int64]map[string]metadata.HostApplyRule)
	ids := make([]int64, 0)

	for _, rule := range rules {
		var id int64
		switch rType {
		case module:
			if rule.ServiceTemplateID != 0 {
				continue
			}
			id = rule.ModuleID
		case serviceTemplate:
			if rule.ModuleID != 0 {
				continue
			}
			id = rule.ServiceTemplateID
		}

		ruleMap, exist := idToRules[id]
		if !exist {
			ids = append(ids, id)
			ruleMap = make(map[string]metadata.HostApplyRule)
		}

		ruleMap[strconv.FormatInt(rule.AttributeID, 10)] = rule
		idToRules[id] = ruleMap
	}

	return idToRules, ids
}

func match(ctx context.Context, rules map[string]metadata.HostApplyRule, attributeMap map[int64]metadata.Attribute,
	filter *querybuilder.QueryFilter) bool {

	rid := util.ExtractRequestIDFromContext(ctx)
	return filter.Match(func(r querybuilder.AtomRule) bool {
		rule, exist := rules[r.Field]
		if !exist {
			return false
		}

		if r.Operator == querybuilder.OperatorExist {
			return true
		}

		prettyValue, err := attributeMap[rule.AttributeID].PrettyValue(ctx, rule.PropertyValue)
		if err != nil {
			blog.Errorf("prettyValue failed, err: %v, rid: %s", err, rid)
			return false
		}

		strValue, ok := r.Value.(string)
		if !ok {
			return false
		}
		if r.Operator == querybuilder.OperatorContains {
			if util.CaseInsensitiveContains(prettyValue, strValue) {
				return true
			}
		}
		return false
	})
}

// BatchUpdateHostApplyRule batch update host apply rule
func (p *hostApplyRule) BatchUpdateHostApplyRule(kit *rest.Kit, bizID int64,
	option metadata.BatchCreateOrUpdateApplyRuleOption) (metadata.BatchCreateOrUpdateHostApplyRuleResult,
	errors.CCErrorCoder) {
	rid := kit.Rid
	batchResult := metadata.BatchCreateOrUpdateHostApplyRuleResult{
		Items: make([]metadata.CreateOrUpdateHostApplyRuleResult, 0),
	}
	now := time.Now()
	for index, item := range option.Rules {
		itemResult := metadata.CreateOrUpdateHostApplyRuleResult{
			Index: index,
		}
		ruleFilter := map[string]interface{}{
			common.BKAppIDField:             bizID,
			common.TenantID:                 kit.TenantID,
			common.BKAttributeIDField:       item.AttributeID,
			common.BKModuleIDField:          item.ModuleID,
			common.BKServiceTemplateIDField: item.ServiceTemplateID,
		}
		count, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostApplyRule).Find(ruleFilter).Count(
			kit.Ctx)
		if err != nil {
			blog.Errorf("get host apply rule failed, filter: %+v, err: %v, rid: %s", ruleFilter, err, rid)
			ccErr := kit.CCError.CCError(common.CCErrCommDBSelectFailed)
			itemResult.SetError(ccErr)
			batchResult.Items = append(batchResult.Items, itemResult)
			continue
		}

		// valid host apply attribute
		attribute, ccErr := p.getHostAttribute(kit, bizID, item.AttributeID)
		if ccErr != nil {
			blog.Errorf("getHostAttribute failed, attribute: %d, err: %v, rid: %s", item.AttributeID, ccErr, rid)
			itemResult.SetError(ccErr)
			batchResult.Items = append(batchResult.Items, itemResult)
			continue
		}
		if value, ok := item.PropertyValue.(string); ok {
			item.PropertyValue = strings.TrimSpace(value)
		}
		rawError := attribute.Validate(kit.Ctx, item.PropertyValue, common.BKPropertyValueField)
		if rawError.ErrCode != 0 {
			ccErr := rawError.ToCCError(kit.CCError)
			blog.Errorf("validate host attribute value failed, attribute: %+v, value: %+v, err: %v, rid: %s",
				attribute, item.PropertyValue, ccErr, kit.Rid)
			itemResult.SetError(ccErr)
			batchResult.Items = append(batchResult.Items, itemResult)
			continue
		}

		if err := hooks.ValidHostApplyStatusHook(kit, p.cs, attribute.PropertyID, item.PropertyValue); err != nil {
			return batchResult, err
		}

		// update rule
		if count > 0 {
			updateData := map[string]interface{}{
				common.BKPropertyValueField: item.PropertyValue,
				common.LastTimeField:        now,
				common.ModifierField:        kit.User,
			}
			if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostApplyRule).Update(kit.Ctx, ruleFilter,
				updateData); err != nil {
				blog.Errorf("update host apply rule failed, filter: %+v, doc: %+v, err: %v, rid: %s", ruleFilter,
					updateData, err, rid)
				ccErr := kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
				itemResult.SetError(ccErr)
			}
			batchResult.Items = append(batchResult.Items, itemResult)
			continue
		}

		// create new rule
		newRuleID, err := mongodb.Shard(kit.SysShardOpts()).NextSequence(kit.Ctx, common.BKTableNameHostApplyRule)
		if err != nil {
			blog.Errorf("%s generate id field failed, err: %s, rid: %s", common.BKTableNameHostApplyRule, err, rid)
			ccErr := kit.CCError.CCError(common.CCErrCommGenerateRecordIDFailed)
			itemResult.SetError(ccErr)
			batchResult.Items = append(batchResult.Items, itemResult)
			continue
		}
		rule := metadata.HostApplyRule{
			ID:                int64(newRuleID),
			BizID:             bizID,
			ModuleID:          item.ModuleID,
			ServiceTemplateID: item.ServiceTemplateID,
			AttributeID:       item.AttributeID,
			PropertyValue:     item.PropertyValue,
			Creator:           kit.User,
			Modifier:          kit.User,
			CreateTime:        now,
			LastTime:          now,
			TenantID:          kit.TenantID,
		}
		if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameHostApplyRule).Insert(kit.Ctx,
			rule); err != nil {
			blog.Errorf("insert host apply rule failed, doc: %+v, err: %v, rid: %s", rule, err, rid)
			ccErr := kit.CCError.CCError(common.CCErrCommDBInsertFailed)
			itemResult.SetError(ccErr)
			batchResult.Items = append(batchResult.Items, itemResult)
			continue
		}
		batchResult.Items = append(batchResult.Items, itemResult)
	}

	for index, item := range option.Rules {
		rule, ccErr := p.GetHostApplyRuleByAttributeID(kit, bizID, item.ModuleID, item.AttributeID)
		if ccErr != nil {
			blog.Errorf("GetHostApplyRuleByAttributeID failed, bizID: %d, moduleID: %d, attribute: %d, err: %s,"+
				" rid: %s", bizID, item.ModuleID, item.AttributeID, ccErr, rid)
			if err := batchResult.Items[index].GetError(); err == nil {
				batchResult.Items[index].SetError(ccErr)
			}
		}
		batchResult.Items[index].Rule = rule
	}

	return batchResult, nil
}

// SearchRuleRelatedServiceTemplates 用于过滤主机应用规则相关的服务模版
func (p *hostApplyRule) SearchRuleRelatedServiceTemplates(kit *rest.Kit,
	option metadata.RuleRelatedServiceTemplateOption) ([]metadata.SrvTemplate, errors.CCErrorCoder) {

	// 1.获取与查询条件中的属性关联的rule和attribute
	rules, attributeMap, ccErr := getRuleAndAttribute(kit, option.ApplicationID, option.QueryFilter, serviceTemplate)
	if ccErr != nil {
		return nil, ccErr
	}

	// 如果没有rule匹配或者是小于attribute的数量，那么说明没有service template是满足查询条件的
	if len(rules) == 0 || len(rules) < len(attributeMap) {
		return nil, nil
	}

	// 2. 将模版与rule进行关联
	srvTemplateToRules, srvTemplateIDs := getRuleRelationIDs(rules, serviceTemplate)

	srvTemplateFilter := map[string]interface{}{
		common.BKAppIDField: option.ApplicationID,
		common.TenantID:     kit.TenantID,
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: srvTemplateIDs,
		},
	}
	srvTemplates := make([]metadata.SrvTemplate, 0)
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceTemplate).Find(srvTemplateFilter).
		All(kit.Ctx, &srvTemplates)
	if err != nil {
		blog.Errorf("find service templates failed, filter: %+v, err: %v, rid: %s", srvTemplateFilter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	// 3.根据匹配的规则过滤出模版
	resultSrvTemplates := make([]metadata.SrvTemplate, 0)
	for _, srvTemplate := range srvTemplates {
		rules, exist := srvTemplateToRules[srvTemplate.ID]
		if !exist {
			continue
		}

		if match(kit.Ctx, rules, attributeMap, option.QueryFilter) {
			resultSrvTemplates = append(resultSrvTemplates, srvTemplate)
		}
	}

	return resultSrvTemplates, nil
}
