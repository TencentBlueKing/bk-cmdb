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
	"fmt"
	"sort"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"

	"github.com/google/go-cmp/cmp"
)

// GenerateApplyPlan 生成主机属性自动应用执行计划
func (p *hostApplyRule) GenerateApplyPlan(kit *rest.Kit, bizID int64, option metadata.HostApplyPlanOption) (metadata.HostApplyPlanResult, errors.CCErrorCoder) {
	rid := kit.Rid

	result := metadata.HostApplyPlanResult{
		Plans:          make([]metadata.OneHostApplyPlan, 0),
		HostAttributes: make([]metadata.Attribute, 0),
	}

	// get hosts
	hostIDs := make([]int64, 0)
	for _, item := range option.HostModules {
		hostIDs = append(hostIDs, item.HostID)
	}
	hostFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}
	hosts := make([]metadata.HostMapStr, 0)
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(hostFilter).All(kit.Ctx, &hosts); err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, list hosts failed, filter: %s, err: %s, rid: %s", hostFilter, err.Error(), rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	// convert to map
	hostID2CloudID := make(map[int64]int64)
	hostMap := make(map[int64]map[string]interface{})
	for _, item := range hosts {
		host := struct {
			HostID  int64 `mapstructure:"bk_host_id"`
			CloudID int64 `mapstructure:"bk_cloud_id"`
		}{}
		if err := mapstruct.Decode2Struct(item, &host); err != nil {
			blog.ErrorJSON("GenerateApplyPlan failed, parse hostID failed, host: %s, err: %s, rid: %s", item, err.Error(), rid)
			return result, kit.CCError.CCError(common.CCErrCommParseDBFailed)
		}
		hostMap[host.HostID] = item
		hostID2CloudID[host.HostID] = host.CloudID
	}

	cloudIDs := make([]int64, 0)
	for _, cloudID := range hostID2CloudID {
		cloudIDs = append(cloudIDs, cloudID)
	}
	clouds := make([]metadata.CloudInst, 0)
	cloudFilter := map[string]interface{}{
		common.BKCloudIDField: map[string]interface{}{
			common.BKDBIN: util.IntArrayUnique(cloudIDs),
		},
	}
	if err := mongodb.Client().Table(common.BKTableNameBasePlat).Find(cloudFilter).All(kit.Ctx, &clouds); err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, read cloud failed, filter: %s, err: %s, rid: %s",
			cloudFilter, err.Error(), rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	cloudMap := make(map[int64]metadata.CloudInst)
	for _, item := range clouds {
		cloudMap[item.CloudID] = item
	}

	// get attributes
	attributeIDs := make([]int64, 0)
	for _, item := range option.Rules {
		attributeIDs = append(attributeIDs, item.AttributeID)
	}
	attributes, err := p.listHostAttributes(kit, bizID, attributeIDs...)
	if err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, listHostAttributes failed, attributeIDs: %s, err: %s, rid: %s",
			attributeIDs, err.Error(), rid)
		return result, err
	}

	// compute apply plan one by one
	hostApplyPlans := make([]metadata.OneHostApplyPlan, 0)
	var hostApplyPlan metadata.OneHostApplyPlan
	unresolvedConflictCount := int64(0)
	for _, hostModule := range option.HostModules {
		host, exist := hostMap[hostModule.HostID]
		if !exist {
			err := errors.New(common.CCErrCommNotFound, fmt.Sprintf("host[%d] not found", hostModule.HostID))
			hostApplyPlan = metadata.OneHostApplyPlan{
				HostID:         hostModule.HostID,
				ExpectHost:     host,
				ConflictFields: nil,
			}
			hostApplyPlan.SetError(err)
			hostApplyPlans = append(hostApplyPlans, hostApplyPlan)
			continue
		}
		hostApplyPlan, err = p.generateOneHostApplyPlan(kit, hostModule.HostID, host, hostModule.ModuleIDs, option.Rules, attributes, option.ConflictResolvers)
		if err != nil {
			blog.ErrorJSON("generateOneHostApplyPlan failed, host: %s, moduleIDs: %s, rules: %s, err: %s, rid: %s", host, hostModule.ModuleIDs, option.Rules, err.Error(), rid)
			return result, err
		}
		if hostApplyPlan.UnresolvedConflictCount > 0 {
			unresolvedConflictCount += 1
		}
		hostApplyPlans = append(hostApplyPlans, hostApplyPlan)
	}

	sort.SliceStable(hostApplyPlans, func(i, j int) bool {
		return hostApplyPlans[i].UnresolvedConflictCount > hostApplyPlans[j].UnresolvedConflictCount
	})

	// fill cloud area info
	for index, item := range hostApplyPlans {
		cloudID, ok := hostID2CloudID[item.HostID]
		if !ok {
			continue
		}
		cloudArea, ok := cloudMap[cloudID]
		if !ok {
			continue
		}
		hostApplyPlans[index].CloudInfo = cloudArea
	}

	result = metadata.HostApplyPlanResult{
		Plans:                   hostApplyPlans,
		Count:                   len(hostApplyPlans),
		UnresolvedConflictCount: unresolvedConflictCount,
		HostAttributes:          attributes,
	}
	return result, nil
}

func (p *hostApplyRule) generateOneHostApplyPlan(
	kit *rest.Kit,
	hostID int64,
	host map[string]interface{},
	moduleIDs []int64,
	rules []metadata.HostApplyRule,
	attributes []metadata.Attribute,
	resolvers []metadata.HostApplyConflictResolver,
) (metadata.OneHostApplyPlan, errors.CCErrorCoder) {
	rid := util.ExtractRequestUserFromContext(kit.Ctx)

	resolverMap := make(map[int64]interface{})
	for _, item := range resolvers {
		if item.HostID != hostID {
			continue
		}
		resolverMap[item.AttributeID] = item.PropertyValue
	}

	plan := metadata.OneHostApplyPlan{
		HostID:                  hostID,
		ModuleIDs:               moduleIDs,
		ExpectHost:              host,
		ConflictFields:          make([]metadata.HostApplyConflictField, 0),
		UpdateFields:            make([]metadata.HostApplyUpdateField, 0),
		UnresolvedConflictCount: 0,
	}

	moduleIDSet := make(map[int64]bool)
	for _, moduleID := range moduleIDs {
		moduleIDSet[moduleID] = true
	}
	attributeRules := make(map[int64][]metadata.HostApplyRule)
	for _, rule := range rules {
		if _, exist := moduleIDSet[rule.ModuleID]; !exist {
			continue
		}
		if _, exist := attributeRules[rule.AttributeID]; !exist {
			attributeRules[rule.AttributeID] = make([]metadata.HostApplyRule, 0)
		}
		attributeRules[rule.AttributeID] = append(attributeRules[rule.AttributeID], rule)
	}

	attributeMap := make(map[int64]metadata.Attribute)
	for _, attribute := range attributes {
		attributeMap[attribute.ID] = attribute
	}

	// update host if conflicts not exist
	for attributeID, targetRules := range attributeRules {
		if len(targetRules) == 0 {
			continue
		}
		attribute, exist := attributeMap[attributeID]
		if !exist {
			blog.Infof("generateOneHostApplyPlan attribute id filed not exist, attributeID: %s, rid: %s", attributeID, rid)
			continue
		}
		if metadata.CheckAllowHostApplyOnField(attribute.PropertyID) == false {
			continue
		}
		propertyIDField := attribute.PropertyID
		originalValue, ok := host[propertyIDField]
		if !ok {
			originalValue = nil
		}

		// check conflicts
		firstValue := targetRules[0].PropertyValue
		conflictedStillExist := false
		for _, rule := range targetRules {
			if cmp.Equal(firstValue, rule.PropertyValue) {
				continue
			}

			conflictedStillExist = true
			if propertyValue, exist := resolverMap[attribute.ID]; exist {
				conflictedStillExist = false
				firstValue = propertyValue
			}

			plan.ConflictFields = append(plan.ConflictFields, metadata.HostApplyConflictField{
				AttributeID:             attributeID,
				PropertyID:              propertyIDField,
				PropertyValue:           originalValue,
				Rules:                   targetRules,
				UnresolvedConflictExist: conflictedStillExist,
			})
			break
		}

		if conflictedStillExist {
			plan.UnresolvedConflictCount += 1
			continue
		}

		// validate property value before update to host
		if value, ok := firstValue.(string); ok {
			firstValue = strings.TrimSpace(value)
			targetRules[0].PropertyValue = firstValue
		}
		rawErr := attribute.Validate(kit.Ctx, firstValue, propertyIDField)
		if rawErr.ErrCode != 0 {
			err := rawErr.ToCCError(kit.CCError)
			blog.ErrorJSON("generateOneHostApplyPlan failed, Validate failed, "+
				"attribute: %s, firstValue: %s, propertyIDField: %s, rawErr: %s, rid: %s",
				attribute, firstValue, propertyIDField, rawErr, rid)
			plan.ErrCode = err.GetCode()
			plan.ErrMsg = err.Error()
			break
		}

		plan.ExpectHost[propertyIDField] = firstValue
		plan.UpdateFields = append(plan.UpdateFields, metadata.HostApplyUpdateField{
			AttributeID:   attributeID,
			PropertyID:    propertyIDField,
			PropertyValue: firstValue,
		})
	}

	sort.SliceStable(plan.UpdateFields, func(i, j int) bool {
		return plan.UpdateFields[i].PropertyID < plan.UpdateFields[j].PropertyID
	})

	sort.SliceStable(plan.ConflictFields, func(i, j int) bool {
		return plan.ConflictFields[i].PropertyID < plan.ConflictFields[j].PropertyID
	})

	return plan, nil
}

func (p *hostApplyRule) RunHostApplyOnHosts(kit *rest.Kit, bizID int64, option metadata.UpdateHostByHostApplyRuleOption) (metadata.MultipleHostApplyResult, errors.CCErrorCoder) {
	rid := kit.Rid
	result := metadata.MultipleHostApplyResult{
		HostResults: make([]metadata.HostApplyResult, 0),
	}
	relationFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: option.HostIDs,
		},
	}
	relations := make([]metadata.ModuleHost, 0)
	if err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(relationFilter).All(kit.Ctx, &relations); err != nil {
		blog.ErrorJSON("RunHostApplyOnHosts failed, find %s failed, filter: %s, err: %s, rid: %s", common.BKTableNameModuleHostConfig, relationFilter, err.Error(), rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	moduleIDs := make([]int64, 0)
	for _, item := range relations {
		moduleIDs = append(moduleIDs, item.ModuleID)
	}
	modules := make([]metadata.ModuleInst, 0)
	moduleFilter := map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIDs,
		},
		common.HostApplyEnabledField: true,
	}
	if err := mongodb.Client().Table(common.BKTableNameBaseModule).Find(moduleFilter).All(kit.Ctx, &modules); err != nil {
		blog.ErrorJSON("RunHostApplyOnHosts failed, find %s failed, filter: %s, err: %s, rid: %s", common.BKTableNameBaseModule, moduleFilter, err.Error(), rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	enableModuleMap := make(map[int64]bool)
	for _, module := range modules {
		enableModuleMap[module.ModuleID] = true
	}
	host2Modules := make(map[int64][]int64)
	for _, relation := range relations {
		if _, exist := host2Modules[relation.HostID]; !exist {
			host2Modules[relation.HostID] = make([]int64, 0)
		}
		if _, exist := enableModuleMap[relation.ModuleID]; !exist {
			continue
		}
		// checkout host apply enabled status on module
		host2Modules[relation.HostID] = append(host2Modules[relation.HostID], relation.ModuleID)
	}
	hostModules := make([]metadata.Host2Modules, 0)
	for hostID, moduleIDs := range host2Modules {
		hostModules = append(hostModules, metadata.Host2Modules{
			HostID:    hostID,
			ModuleIDs: moduleIDs,
		})
	}
	listHostApplyRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: moduleIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	rules, ccErr := p.ListHostApplyRule(kit, bizID, listHostApplyRuleOption)
	if ccErr != nil {
		blog.ErrorJSON("RunHostApplyOnHosts failed, ListHostApplyRule failed, option: %s, err: %s, rid: %s", common.BKTableNameModuleHostConfig, listHostApplyRuleOption, ccErr.Error(), rid)
		return result, ccErr
	}
	planOption := metadata.HostApplyPlanOption{
		Rules:       rules.Info,
		HostModules: hostModules,
	}
	planResult, ccErr := p.GenerateApplyPlan(kit, bizID, planOption)
	if ccErr != nil {
		blog.ErrorJSON("RunHostApplyOnHosts failed, find %s failed, filter: %s, err: %s, rid: %s", common.BKTableNameModuleHostConfig, relationFilter, ccErr.Error(), rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	for _, plan := range planResult.Plans {
		applyResult := metadata.HostApplyResult{
			ErrorContainer: metadata.ErrorContainer{},
			HostID:         0,
		}
		updateData := plan.GetUpdateData()
		if len(updateData) == 0 {
			result.HostResults = append(result.HostResults, applyResult)
			continue
		}

		updateOption := metadata.UpdateOption{
			Condition: map[string]interface{}{
				common.BKHostIDField: plan.HostID,
			},
			Data: updateData,
		}
		_, err := p.dependence.UpdateModelInstance(kit, common.BKInnerObjIDHost, updateOption)
		blog.Warnf("RunHostApplyOnHosts failed, UpdateModelInstance failed, hostID: %d, updateOption: %+v, err: %+v, rid: %s", plan.HostID, updateOption, err, rid)
		if err != nil {
			ccErr, ok := err.(errors.CCErrorCoder)
			if ok {
				applyResult.SetError(ccErr)
			} else {
				ccErr := kit.CCError.CCError(common.CCErrHostUpdateFail)
				applyResult.SetError(ccErr)
			}
		}
		result.HostResults = append(result.HostResults, applyResult)
	}

	for _, hostResult := range result.HostResults {
		if ccErr := hostResult.GetError(); ccErr != nil {
			result.SetError(ccErr)
		}
	}
	return result, result.GetError()
}
