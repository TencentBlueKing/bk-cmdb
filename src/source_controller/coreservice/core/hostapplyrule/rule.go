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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type hostApplyRule struct {
	dbProxy dal.RDB
}

func New(dbProxy dal.RDB) core.HostApplyRuleOperation {
	rule := &hostApplyRule{
		dbProxy: dbProxy,
	}
	return rule
}

func (p *hostApplyRule) validateModuleID(ctx core.ContextParams, bizID int64, moduleID int64) errors.CCErrorCoder {
	filter := map[string]interface{}{
		common.BKAppIDField:    bizID,
		common.BKModuleIDField: moduleID,
	}
	count, err := p.dbProxy.Table(common.BKTableNameBaseModule).Find(filter).Count(ctx.Context)
	if err != nil {
		blog.Errorf("ValidateModuleID failed, validate module id failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
		return ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	if count == 0 {
		return ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}
	return nil
}

func (p *hostApplyRule) listHostAttributes(ctx core.ContextParams, bizID int64, hostAttributeIDs ...int64) ([]metadata.Attribute, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{
				// business private attribute
				metadata.MetadataBizField: map[string]interface{}{
					common.BKDBEQ: strconv.FormatInt(bizID, 10),
				},
			}, {
				// global attribute
				metadata.MetadataBizField: map[string]interface{}{
					common.BKDBExists: false,
				},
			}, {
				// global attribute
				metadata.BKMetadata: map[string]interface{}{
					common.BKDBExists: false,
				},
			},
		},
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: hostAttributeIDs,
		},
	}
	attributes := make([]metadata.Attribute, 0)
	err := p.dbProxy.Table(common.BKTableNameObjAttDes).Find(filter).All(ctx.Context, &attributes)
	if err != nil {
		if p.dbProxy.IsNotFoundError(err) {
			blog.Errorf("get host attribute failed, not found, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
			return attributes, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		blog.Errorf("get host attribute failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
		return attributes, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	return attributes, nil
}

func (p *hostApplyRule) getHostAttribute(ctx core.ContextParams, bizID int64, hostAttributeID int64) (metadata.Attribute, errors.CCErrorCoder) {
	attribute := metadata.Attribute{}
	attributes, err := p.listHostAttributes(ctx, bizID, hostAttributeID)
	if err != nil {
		blog.Errorf("getHostAttribute failed, listHostAttributes failed, bizID: %d, attribute: %d, err: %s, rid: %s", bizID, hostAttributeID, err.Error(), ctx.ReqID)
		return attribute, err
	}
	if len(attributes) == 0 {
		return attribute, ctx.Error.CCError(common.CCErrCommNotFound)
	}
	if len(attributes) > 1 {
		return attribute, ctx.Error.CCError(common.CCErrCommGetMultipleObject)
	}
	return attributes[0], nil
}

func (p *hostApplyRule) CreateHostApplyRule(ctx core.ContextParams, bizID int64, option metadata.CreateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder) {
	now := time.Now()
	rule := metadata.HostApplyRule{
		ID:              0,
		BizID:           bizID,
		AttributeID:     option.AttributeID,
		ModuleID:        option.ModuleID,
		PropertyValue:   option.PropertyValue,
		Creator:         ctx.User,
		Modifier:        ctx.User,
		CreateTime:      now,
		LastTime:        now,
		SupplierAccount: ctx.SupplierAccount,
	}
	if key, err := rule.Validate(); err != nil {
		blog.Errorf("CreateHostApplyRule failed, parameter invalid, key: %s, err: %+v, rid: %s", key, err, ctx.ReqID)
		return rule, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, key)
	}

	// validate bk_module_id
	if err := p.validateModuleID(ctx, bizID, rule.ModuleID); err != nil {
		blog.Errorf("CreateHostApplyRule failed, validate bk_module_id failed, bizID: %d, moduleID: %d, err: %s, rid: %s", bizID, err.Error(), ctx.ReqID)
		return rule, err
	}

	attribute, ccErr := p.getHostAttribute(ctx, bizID, rule.AttributeID)
	if ccErr != nil {
		blog.Errorf("CreateHostApplyRule failed, get host attribute failed, bizID: %d, attributeID: %d, err: %+v, rid: %s", bizID, rule.AttributeID, ccErr, ctx.ReqID)
		return rule, ccErr
	}

	rawError := attribute.Validate(ctx.Context, option.PropertyValue, common.BKPropertyValueField)
	if rawError.ErrCode != 0 {
		ccErr := rawError.ToCCError(ctx.Error)
		blog.Errorf("CreateHostApplyRule failed, validate host attribute value failed,  attribute: %+v, value: %+v, err: %+v, rid: %s", attribute, option.PropertyValue, ccErr, ctx.ReqID)
		return rule, ccErr
	}

	// generate id field
	id, err := p.dbProxy.NextSequence(ctx, common.BKTableNameHostApplyRule)
	if nil != err {
		blog.Errorf("CreateHostApplyRule failed, generate id failed, err: %+v, rid: %s", err, ctx.ReqID)
		return rule, ctx.Error.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	rule.ID = int64(id)

	if err := p.dbProxy.Table(common.BKTableNameHostApplyRule).Insert(ctx.Context, rule); err != nil {
		if p.dbProxy.IsDuplicatedError(err) {
			blog.Errorf("CreateHostApplyRule failed, duplicated error, doc: %+v, err: %+v, rid: %s", rule, err, ctx.ReqID)
			return rule, ctx.Error.CCErrorf(common.CCErrCommDuplicateItem, common.BKAttributeIDField)
		}
		blog.Errorf("CreateHostApplyRule failed, db insert failed, doc: %+v, err: %+v, rid: %s", rule, err, ctx.ReqID)
		return rule, ctx.Error.CCError(common.CCErrCommDBInsertFailed)
	}

	return rule, nil
}

func (p *hostApplyRule) UpdateHostApplyRule(ctx core.ContextParams, bizID int64, ruleID int64, option metadata.UpdateHostApplyRuleOption) (metadata.HostApplyRule, errors.CCErrorCoder) {
	rule, err := p.GetHostApplyRule(ctx, bizID, ruleID)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, rule not found, bizID: %d, id: %d, rid: %s", bizID, ruleID, ctx.ReqID)
		return rule, ctx.Error.CCError(common.CCErrCommNotFound)
	}

	attribute, err := p.getHostAttribute(ctx, bizID, rule.AttributeID)
	rawError := attribute.Validate(ctx.Context, option.PropertyValue, common.BKPropertyValueField)
	if rawError.ErrCode != 0 {
		ccErr := rawError.ToCCError(ctx.Error)
		blog.Errorf("UpdateHostApplyRule failed, validate host attribute value failed,  attribute: %+v, value: %+v, err: %+v, rid: %s", attribute, option.PropertyValue, ccErr, ctx.ReqID)
		return rule, ccErr
	}

	rule.LastTime = time.Now()
	rule.Modifier = ctx.User
	rule.PropertyValue = option.PropertyValue

	filter := map[string]interface{}{
		common.BKFieldID: ruleID,
	}
	if err := p.dbProxy.Table(common.BKTableNameSetTemplate).Update(ctx.Context, filter, rule); err != nil {
		blog.ErrorJSON("UpdateHostApplyRule failed, db update failed, filter: %s, doc: %s, err: %s, rid: %s", filter, rule, err, ctx.ReqID)
		return rule, ctx.Error.CCError(common.CCErrCommDBUpdateFailed)
	}

	return rule, nil
}

func (p *hostApplyRule) DeleteHostApplyRule(ctx core.ContextParams, bizID int64, ruleIDs ...int64) errors.CCErrorCoder {
	if len(ruleIDs) == 0 {
		return ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "host_apply_rule_ids")
	}
	filter := map[string]interface{}{
		common.BKOwnerIDField: ctx.SupplierAccount,
		common.BKAppIDField:   bizID,
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: ruleIDs,
		},
	}
	if err := p.dbProxy.Table(common.BKTableNameHostApplyRule).Delete(ctx.Context, filter); err != nil {
		blog.Errorf("DeleteHostApplyRule failed, db remove failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
		return ctx.Error.CCError(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

func (p *hostApplyRule) GetHostApplyRule(ctx core.ContextParams, bizID int64, ruleID int64) (metadata.HostApplyRule, errors.CCErrorCoder) {
	rule := metadata.HostApplyRule{}
	filter := map[string]interface{}{
		common.BkSupplierAccount: ctx.SupplierAccount,
		common.BKAppIDField:      bizID,
		common.BKFieldID:         ruleID,
	}
	if err := p.dbProxy.Table(common.BKTableNameHostApplyRule).Find(filter).One(ctx.Context, &rule); err != nil {
		if p.dbProxy.IsNotFoundError(err) {
			blog.Errorf("GetHostApplyRule failed, db select failed, not found, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
			return rule, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		blog.Errorf("GetHostApplyRule failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
		return rule, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	return rule, nil
}

func (p *hostApplyRule) GetHostApplyRuleByAttributeID(ctx core.ContextParams, bizID, moduleID, attributeID int64) (metadata.HostApplyRule, errors.CCErrorCoder) {
	rule := metadata.HostApplyRule{}
	filter := map[string]interface{}{
		common.BkSupplierAccount:  ctx.SupplierAccount,
		common.BKAppIDField:       bizID,
		common.BKModuleIDField:    moduleID,
		common.BKAttributeIDField: attributeID,
	}
	if err := p.dbProxy.Table(common.BKTableNameHostApplyRule).Find(filter).One(ctx.Context, &rule); err != nil {
		if p.dbProxy.IsNotFoundError(err) {
			blog.Errorf("GetHostApplyRuleByAttributeID failed, db select failed, not found, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
			return rule, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		blog.Errorf("GetHostApplyRuleByAttributeID failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
		return rule, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	return rule, nil
}

func (p *hostApplyRule) ListHostApplyRule(ctx core.ContextParams, bizID int64, option metadata.ListHostApplyRuleOption) (metadata.MultipleHostApplyRuleResult, errors.CCErrorCoder) {
	result := metadata.MultipleHostApplyRuleResult{}
	if option.Page.Limit > common.BKMaxPageSize && option.Page.Limit != common.BKNoLimit {
		return result, ctx.Error.CCError(common.CCErrCommPageLimitIsExceeded)
	}

	filter := map[string]interface{}{
		common.BkSupplierAccount: ctx.SupplierAccount,
		common.BKAppIDField:      bizID,
	}
	if option.ModuleIDs != nil {
		filter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: option.ModuleIDs,
		}
	}
	query := p.dbProxy.Table(common.BKTableNameHostApplyRule).Find(filter)
	total, err := query.Count(ctx.Context)
	if err != nil {
		blog.ErrorJSON("ListHostApplyRule failed, db count failed, filter: %s, err: %s, rid: %s", filter, err.Error(), ctx.ReqID)
		return result, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
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
	if err := query.All(ctx.Context, &rules); err != nil {
		blog.ErrorJSON("ListHostApplyRule failed, db select failed, filter: %s, err: %s, rid: %s", filter, err.Error(), ctx.ReqID)
		return result, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	result.Info = rules
	return result, nil
}

// SearchRuleRelatedModules 用于过滤主机应用规则相关的模块
/*
支持场景：
_ 支持通过模块名过滤
_ 支持通过模块上设置的主机应用配置字段名过滤
_ 支持通过模块上设置的主机应用配置字段值过滤，字段值需要支持数值型和枚举字段的过滤，枚举类型翻译成对应的name域再过滤
*/
func (p *hostApplyRule) SearchRuleRelatedModules(ctx core.ContextParams, bizID int64, option metadata.SearchRuleRelatedModulesOption) ([]metadata.Module, errors.CCErrorCoder) {
	rid := ctx.ReqID

	// list modules
	moduleFilter := map[string]interface{}{
		common.BKAppIDField:      bizID,
		common.BkSupplierAccount: ctx.SupplierAccount,
	}
	modules := make([]metadata.Module, 0)
	if err := p.dbProxy.Table(common.BKTableNameBaseModule).Find(moduleFilter).All(ctx.Context, &modules); err != nil {
		blog.ErrorJSON("SearchRuleRelatedModules failed, find modules failed, filter: %s, err: %s, rid: %s", moduleFilter, err.Error(), rid)
		return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	moduleMap := make(map[int64]metadata.Module)
	for _, module := range modules {
		moduleMap[module.ModuleID] = module
	}

	// list rules
	ruleFilter := map[string]interface{}{
		common.BKAppIDField:      bizID,
		common.BkSupplierAccount: ctx.SupplierAccount,
	}
	rules := make([]metadata.HostApplyRule, 0)
	if err := p.dbProxy.Table(common.BKTableNameHostApplyRule).Find(ruleFilter).All(ctx.Context, &rules); err != nil {
		blog.ErrorJSON("SearchRuleRelatedModules failed, find rules failed, filter: %s, err: %s, rid: %s", ruleFilter, err.Error(), rid)
		return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	// list attributes
	attributeIDs := make([]int64, 0)
	for _, item := range rules {
		attributeIDs = append(attributeIDs, item.AttributeID)
	}
	attributeFilter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: attributeIDs,
		},
	}
	attributes := make([]metadata.Attribute, 0)
	if err := p.dbProxy.Table(common.BKTableNameObjAttDes).Find(attributeFilter).All(ctx.Context, &attributes); err != nil {
		blog.ErrorJSON("SearchRuleRelatedModules failed, find attributes failed, filter: %s, err: %s, rid: %s", attributeFilter, err.Error(), rid)
		return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	// attribute map
	attributeMap := make(map[int64]metadata.Attribute)
	for _, attribute := range attributes {
		attributeMap[attribute.ID] = attribute
	}

	resultModuleMap := make(map[int64]bool)
	resultModules := make([]metadata.Module, 0)
	for _, module := range modules {
		if matchModule(ctx, module, option) {
			resultModuleMap[module.ModuleID] = true
			resultModules = append(resultModules, module)
			continue
		}
	}

	for _, rule := range rules {
		attribute, exist := attributeMap[rule.AttributeID]
		if exist == false {
			continue
		}
		if matchRule(ctx, rule, attribute, option) {
			module, exist := moduleMap[rule.ModuleID]
			if exist == false {
				continue
			}
			// avoid repeat
			if _, exist := resultModuleMap[module.ModuleID]; exist == true {
				continue
			}
			resultModules = append(resultModules, module)
		}
	}
	return resultModules, nil
}

func matchModule(ctx context.Context, module metadata.Module, option metadata.SearchRuleRelatedModulesOption) bool {
	if option.QueryFilter == nil {
		return true
	}
	return option.QueryFilter.Match(func(r querybuilder.AtomRule) bool {
		if r.Field != metadata.TopoNodeKeyword {
			return false
		}
		strValue, ok := r.Value.(string)
		if ok == false {
			return false
		}
		if util.CaseInsensitiveContains(module.ModuleName, strValue) {
			return true
		}
		return false
	})
}

func matchRule(ctx context.Context, rule metadata.HostApplyRule, attribute metadata.Attribute, option metadata.SearchRuleRelatedModulesOption) bool {
	rid := util.ExtractRequestIDFromContext(ctx)

	prettyValue, err := attribute.PrettyValue(ctx, rule.PropertyValue)
	if err != nil {
		blog.Errorf("matchRule failed, PrettyValue failed, err: %s, rid: %s", err.Error(), rid)
		return false
	}

	return option.QueryFilter.Match(func(r querybuilder.AtomRule) bool {
		if r.Field != strconv.FormatInt(attribute.ID, 10) {
			return false
		}
		strValue, ok := r.Value.(string)
		if ok == false {
			return false
		}
		if util.CaseInsensitiveContains(prettyValue, strValue) {
			return true
		}
		return false
	})
}

func (p *hostApplyRule) BatchUpdateHostApplyRule(ctx core.ContextParams, bizID int64, option metadata.BatchCreateOrUpdateApplyRuleOption) (metadata.BatchCreateOrUpdateHostApplyRuleResult, errors.CCErrorCoder) {
	rid := ctx.ReqID
	batchResult := metadata.BatchCreateOrUpdateHostApplyRuleResult{
		Items: make([]metadata.CreateOrUpdateHostApplyRuleResult, 0),
	}
	now := time.Now()
	for index, item := range option.Rules {
		itemResult := metadata.CreateOrUpdateHostApplyRuleResult{
			Index: index,
		}
		ruleFilter := map[string]interface{}{
			common.BKAppIDField:       bizID,
			common.BkSupplierAccount:  ctx.SupplierAccount,
			common.BKAttributeIDField: item.AttributeID,
			common.BKModuleIDField:    item.ModuleID,
		}
		count, err := p.dbProxy.Table(common.BKTableNameHostApplyRule).Find(ruleFilter).Count(ctx.Context)
		if err != nil {
			blog.ErrorJSON("BatchUpdateHostApplyRule failed, find rule failed, filter: %s, err: %s, rid: %s", ruleFilter, err.Error(), rid)
			ccErr := ctx.Error.CCError(common.CCErrCommDBSelectFailed)
			itemResult.SetError(ccErr)
			batchResult.Items = append(batchResult.Items, itemResult)
			continue
		}

		// update rule
		if count > 0 {
			updateData := map[string]interface{}{
				common.BKPropertyValueField: item.PropertyValue,
				common.LastTimeField:        now,
				common.ModifierField:        ctx.User,
			}
			if err := p.dbProxy.Table(common.BKTableNameHostApplyRule).Update(ctx.Context, ruleFilter, updateData); err != nil {
				blog.ErrorJSON("BatchUpdateHostApplyRule failed, update rule failed, filter: %s, doc: %s, err: %s, rid: %s", ruleFilter, updateData, err.Error(), rid)
				ccErr := ctx.Error.CCError(common.CCErrCommDBUpdateFailed)
				itemResult.SetError(ccErr)
			}
			batchResult.Items = append(batchResult.Items, itemResult)
			continue
		}

		// create new rule
		newRuleID, err := p.dbProxy.NextSequence(ctx.Context, common.BKTableNameHostApplyRule)
		if err != nil {
			blog.ErrorJSON("BatchUpdateHostApplyRule failed, generate id field failed, err: %s, rid: %s", err.Error(), rid)
			ccErr := ctx.Error.CCError(common.CCErrCommGenerateRecordIDFailed)
			itemResult.SetError(ccErr)
			batchResult.Items = append(batchResult.Items, itemResult)
			continue
		}
		rule := metadata.HostApplyRule{
			ID:              int64(newRuleID),
			BizID:           bizID,
			ModuleID:        item.ModuleID,
			AttributeID:     item.AttributeID,
			PropertyValue:   item.PropertyValue,
			Creator:         ctx.User,
			Modifier:        ctx.User,
			CreateTime:      now,
			LastTime:        now,
			SupplierAccount: ctx.SupplierAccount,
		}
		if err := p.dbProxy.Table(common.BKTableNameHostApplyRule).Insert(ctx.Context, rule); err != nil {
			blog.ErrorJSON("BatchUpdateHostApplyRule failed, insert rule failed, doc: %s, err: %s, rid: %s", rule, err.Error(), rid)
			ccErr := ctx.Error.CCError(common.CCErrCommDBInsertFailed)
			itemResult.SetError(ccErr)
			batchResult.Items = append(batchResult.Items, itemResult)
			continue
		}
		batchResult.Items = append(batchResult.Items, itemResult)
	}

	for index, item := range option.Rules {
		rule, ccErr := p.GetHostApplyRuleByAttributeID(ctx, bizID, item.ModuleID, item.AttributeID)
		if ccErr != nil {
			blog.Errorf("GetHostApplyRuleByAttributeID failed, bizID: %d, moduleID: %d, attribute: %d, err: %s, rid: %s", bizID, item.ModuleID, item.AttributeID, ccErr.Error(), rid)
			if err := batchResult.Items[index].GetError(); err == nil {
				batchResult.Items[index].SetError(ccErr)
			}
		}
		batchResult.Items[index].Rule = rule
	}

	return batchResult, nil
}
