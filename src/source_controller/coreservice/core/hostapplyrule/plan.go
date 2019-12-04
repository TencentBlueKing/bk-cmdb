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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"

	"github.com/google/go-cmp/cmp"
)

// GenerateApplyPlan 生成主机属性自动应用执行计划
func (p *hostApplyRule) GenerateApplyPlan(ctx core.ContextParams, bizID int64, option metadata.HostApplyPlanOption) (metadata.HostApplyPlanResult, errors.CCErrorCoder) {
	rid := ctx.ReqID

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
	hosts := make([]map[string]interface{}, 0)
	if err := p.dbProxy.Table(common.BKTableNameBaseHost).Find(hostFilter).All(ctx.Context, &hosts); err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, list hosts failed, filter: %s, err: %s, rid: %s", hostFilter, err.Error(), rid)
		return result, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
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
			return result, ctx.Error.CCError(common.CCErrCommParseDBFailed)
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
	if err := p.dbProxy.Table(common.BKTableNameBasePlat).Find(cloudFilter).All(ctx.Context, &clouds); err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, read cloud failed, filter: %s, err: %s, rid: %s",
			cloudFilter, err.Error(), rid)
		return result, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
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
	attributes, err := p.listHostAttributes(ctx, bizID, attributeIDs...)
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
		if exist == false {
			err := errors.New(common.CCErrCommNotFound, fmt.Sprintf("host[%d] not found", hostModule.HostID))
			hostApplyPlan = metadata.OneHostApplyPlan{
				HostID:         hostModule.HostID,
				ExpiredHost:    host,
				ConflictFields: nil,
			}
			hostApplyPlan.SetError(err)
			hostApplyPlans = append(hostApplyPlans, hostApplyPlan)
			continue
		}
		hostApplyPlan, err = p.generateOneHostApplyPlan(ctx, hostModule.HostID, host, hostModule.ModuleIDs, option.Rules, attributes, option.ConflictResolvers)
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
		if ok == false {
			continue
		}
		cloudArea, ok := cloudMap[cloudID]
		if ok == false {
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
	ctx core.ContextParams,
	hostID int64,
	host map[string]interface{},
	moduleIDs []int64,
	rules []metadata.HostApplyRule,
	attributes []metadata.Attribute,
	resolvers []metadata.HostApplyConflictResolver,
) (metadata.OneHostApplyPlan, errors.CCErrorCoder) {
	rid := util.ExtractRequestUserFromContext(ctx)

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
		ExpiredHost:             host,
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
		if _, exist := moduleIDSet[rule.ModuleID]; exist == false {
			continue
		}
		if _, exist := attributeRules[rule.AttributeID]; exist == false {
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
		if exist == false {
			blog.ErrorJSON("generateOneHostApplyPlan failed, attribute id filed not exist, "+
				"attributeMap: %s, attributeID: %s, rid: %s", attributeMap, attributeID, rid)
			err := ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAttributeIDField)
			plan.ErrCode = err.GetCode()
			plan.ErrMsg = err.Error()
		}
		propertyIDField := attribute.PropertyID
		originalValue, ok := host[propertyIDField]
		if ok == false {
			originalValue = nil
		}

		// check conflicts
		firstValue := targetRules[0].PropertyValue
		conflictedStillExist := false
		for _, rule := range targetRules {
			if cmp.Equal(firstValue, rule.PropertyValue) == false {
				conflictedStillExist = true
				if propertyValue, exist := resolverMap[attribute.ID]; exist == true {
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
		}

		if conflictedStillExist == true {
			plan.UnresolvedConflictCount += 1
			continue
		}

		// validate property value before update to host
		rawErr := attribute.Validate(ctx.Context, firstValue, propertyIDField)
		if rawErr.ErrCode != 0 {
			err := rawErr.ToCCError(ctx.Error)
			blog.ErrorJSON("generateOneHostApplyPlan failed, Validate failed, "+
				"attribute: %s, firstValue: %s, propertyIDField: %s, rawErr: %s, rid: %s",
				attribute, firstValue, propertyIDField, rawErr, rid)
			plan.ErrCode = err.GetCode()
			plan.ErrMsg = err.Error()
			break
		}

		plan.ExpiredHost[propertyIDField] = firstValue
		plan.UpdateFields = append(plan.UpdateFields, metadata.HostApplyUpdateField{
			AttributeID:   attributeID,
			PropertyID:    propertyIDField,
			PropertyValue: firstValue,
		})
	}

	return plan, nil
}
