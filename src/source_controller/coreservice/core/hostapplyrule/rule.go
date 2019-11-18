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
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
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

func (p *hostApplyRule) getHostAttribute(ctx core.ContextParams, bizID int64, hostAttributeID int64) (metadata.Attribute, errors.CCErrorCoder) {
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
			},
		},
		common.BKFieldID: hostAttributeID,
	}
	attribute := metadata.Attribute{}
	err := p.dbProxy.Table(common.BKTableNameObjAttDes).Find(filter).One(ctx.Context, &attribute)
	if err != nil {
		if p.dbProxy.IsNotFoundError(err) {
			blog.Errorf("get host attribute failed, not found, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
			return attribute, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		blog.Errorf("get host attribute failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
		return attribute, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	return attribute, nil
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
