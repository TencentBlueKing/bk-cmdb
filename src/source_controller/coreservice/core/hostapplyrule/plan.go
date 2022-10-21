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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// get attributes
	attributeIDs := make([]int64, 0)
	for _, item := range option.Rules {
		attributeIDs = append(attributeIDs, item.AttributeID)
	}
	attributes, err := p.listHostAttributes(kit, bizID, attributeIDs...)
	if err != nil {
		blog.Errorf("GenerateApplyPlan failed, listHostAttributes failed, attributeIDs: %s, err: %s, rid: %s",
			attributeIDs, err.Error(), rid)
		return result, err
	}

	fields := []string{common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField, common.BKHostNameField}
	for _, attr := range attributes {
		fields = append(fields, attr.PropertyID)
	}

	hosts := make([]metadata.HostMapStr, 0)
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(hostFilter).Fields(fields...).All(kit.Ctx, &hosts); err != nil {
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
			hostApplyPlans = append(hostApplyPlans, hostApplyPlan)
		}
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
		Count:                   len(option.HostModules),
		UnresolvedConflictCount: unresolvedConflictCount,
		HostAttributes:          attributes,
	}
	return result, nil
}

// isRuleEqualOrNot : When the attribute type is "organization", the type obtained from the database is the database's
// native primitive.A type. When converted to []interface{}, the type is int64, the type of propertyValue is
// []interface{}, and the type of each element inside is json.Number, which needs to be unified before comparing.
// The rest of the attribute types can be compared directly in non-organization scenarios.
func isRuleEqualOrNot(pType string, expectValue interface{}, propertyValue interface{}) (bool, errors.CCErrorCoder) {

	// in the transfer host scenario, the rule may be empty.
	if expectValue == nil {
		return false, nil
	}
	switch pType {
	case common.FieldTypeOrganization:

		value, ok := expectValue.(primitive.A)
		if !ok {
			return false, errors.New(common.CCErrCommUnexpectedFieldType, "expect value type error")
		}
		if _, ok := propertyValue.([]interface{}); !ok {
			return false, errors.New(common.CCErrCommUnexpectedFieldType, "property value type error")
		}

		expectValueList := make([]int, 0)
		for _, eValue := range []interface{}(value) {
			value, err := util.GetIntByInterface(eValue)
			if err != nil {
				return false, errors.New(common.CCErrCommUnexpectedFieldType, err.Error())
			}
			expectValueList = append(expectValueList, value)
		}

		ruleValueList := make([]int, 0)
		for _, rValue := range propertyValue.([]interface{}) {
			value, err := util.GetIntByInterface(rValue)
			if err != nil {
				return false, errors.New(common.CCErrCommUnexpectedFieldType, err.Error())
			}
			ruleValueList = append(ruleValueList, value)
		}
		if cmp.Equal(expectValueList, ruleValueList) {
			return true, nil
		}

	// 当属性是int类型时，需要转为统一类型进行对比
	case common.FieldTypeInt:
		origin, err := util.GetIntByInterface(propertyValue)
		if err != nil {
			return false, errors.New(common.CCErrCommUnexpectedFieldType, err.Error())
		}
		expect, err := util.GetIntByInterface(expectValue)
		if err != nil {
			return false, errors.New(common.CCErrCommUnexpectedFieldType, err.Error())
		}
		if cmp.Equal(origin, expect) {
			return true, nil
		}

	case common.FieldTypeTime:
		expectVal, ok := expectValue.(primitive.DateTime)
		if !ok {
			return false, errors.New(common.CCErrCommUnexpectedFieldType, "expect value type error")
		}
		expectTimeVal := expectVal.Time()

		propertyTimeValue, err := metadata.ParseTime(propertyValue)
		if err != nil {
			return false, errors.New(common.CCErrCommUnexpectedFieldType, err.Error())
		}

		if cmp.Equal(expectTimeVal, propertyTimeValue) {
			return true, nil
		}

	default:
		if cmp.Equal(expectValue, propertyValue) {
			return true, nil
		}
	}
	return false, nil
}

func preCheckRules(targetRules []metadata.HostApplyRule, attributeID int64, attrMap map[int64]metadata.Attribute,
	rid string) (metadata.Attribute, bool) {
	if len(targetRules) == 0 {
		return metadata.Attribute{}, false
	}
	attribute, exist := attrMap[attributeID]
	if !exist {
		blog.Infof("attribute id field not exist, attributeID: %s, rid: %s", attributeID, rid)
		return metadata.Attribute{}, false
	}
	if !metadata.CheckAllowHostApplyOnField(&attribute) {
		return metadata.Attribute{}, false
	}
	return attribute, true
}

func getOneHostApplyPlan(kit *rest.Kit, attrRules map[int64][]metadata.HostApplyRule,
	attrMap map[int64]metadata.Attribute, hostID int64, host map[string]interface{}, moduleIDs []int64,
	resolverMap map[int64]interface{}) (metadata.OneHostApplyPlan, errors.CCErrorCoder) {

	rid := util.ExtractRequestUserFromContext(kit.Ctx)
	plan := metadata.OneHostApplyPlan{
		HostID:         hostID,
		ModuleIDs:      moduleIDs,
		ExpectHost:     host,
		ConflictFields: make([]metadata.HostApplyConflictField, 0),
		UpdateFields:   make([]metadata.HostApplyUpdateField, 0),
	}

	for attributeID, targetRules := range attrRules {

		attribute, need := preCheckRules(targetRules, attributeID, attrMap, rid)
		if !need {
			continue
		}
		propertyIDField := attribute.PropertyID
		originalValue, ok := host[propertyIDField]
		if !ok {
			originalValue = nil
		}

		expectValue := originalValue

		// check conflicts and if needChange
		conflictedStillExist, needChange := false, false
		// check if host needs to be changed by the host apply rules, if not, do not append the field to the update
		// fields
		for _, rule := range targetRules {
			isEqual, err := isRuleEqualOrNot(attribute.PropertyType, expectValue, rule.PropertyValue)
			if err != nil {
				blog.Errorf("compare rule value failed, err: %v, rid: %s", err, rid)
				return metadata.OneHostApplyPlan{}, err
			}

			if isEqual {
				continue
			}

			needChange = true
			expectValue = rule.PropertyValue
			conflictedStillExist = true
			if propertyValue, exist := resolverMap[attribute.ID]; exist {
				conflictedStillExist = false
				expectValue = propertyValue
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

		if !needChange {
			continue
		}

		if conflictedStillExist {
			plan.UnresolvedConflictCount += 1
		}

		// validate property value before update to host
		if value, ok := expectValue.(string); ok {
			expectValue = strings.TrimSpace(value)
			targetRules[0].PropertyValue = expectValue
		}
		rawErr := attribute.Validate(kit.Ctx, expectValue, propertyIDField)
		if rawErr.ErrCode != 0 {
			blog.Errorf("attribute validate failed, attribute: %s, firstValue: %s, propertyID: %s, err: %s, rid: %s",
				attribute, expectValue, propertyIDField, rawErr, rid)
			plan.ErrCode = rawErr.ToCCError(kit.CCError).GetCode()
			plan.ErrMsg = rawErr.ToCCError(kit.CCError).Error()
			break
		}

		plan.ExpectHost[propertyIDField] = expectValue
		plan.UpdateFields = append(plan.UpdateFields, metadata.HostApplyUpdateField{
			AttributeID:   attributeID,
			PropertyID:    propertyIDField,
			PropertyValue: expectValue,
		})
	}
	return plan, nil
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

	resolverMap := make(map[int64]interface{})
	for _, item := range resolvers {
		if item.HostID != hostID {
			continue
		}
		resolverMap[item.AttributeID] = item.PropertyValue
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

	plan, err := getOneHostApplyPlan(kit, attributeRules, attributeMap, hostID, host, moduleIDs, resolverMap)
	if err != nil {
		return metadata.OneHostApplyPlan{}, err
	}

	sort.SliceStable(plan.UpdateFields, func(i, j int) bool {
		return plan.UpdateFields[i].PropertyID < plan.UpdateFields[j].PropertyID
	})

	sort.SliceStable(plan.ConflictFields, func(i, j int) bool {
		return plan.ConflictFields[i].PropertyID < plan.ConflictFields[j].PropertyID
	})

	return plan, nil
}

func getModuleIDsAndSrvTempIDs(kit *rest.Kit, modules []metadata.ModuleInst) (map[int64]struct{}, []int64,
	map[int64][]int64, errors.CCErrorCoder) {

	srvTemplateIDs := make([]int64, 0)
	moduleIDHostApplyEnabledMap := make(map[int64]bool)

	for _, module := range modules {
		if module.ServiceTemplateID != 0 {
			srvTemplateIDs = append(srvTemplateIDs, module.ServiceTemplateID)
		}
		moduleIDHostApplyEnabledMap[module.ModuleID] = module.HostApplyEnabled
	}
	existSrvTempIDsMap := make(map[int64]struct{})

	// the list of modules that are automatically applied
	// to the host, including the modules that are automatically
	// applied to the corresponding template host.
	enableModuleMap := make(map[int64]struct{})

	// store the correspondence between templates and modules
	srvTempModulesMap := make(map[int64][]int64)

	if len(srvTemplateIDs) > 0 {
		filter := map[string]interface{}{
			common.BKFieldID:             map[string]interface{}{common.BKDBIN: srvTemplateIDs},
			common.HostApplyEnabledField: true,
		}
		serviceTemplates := make([]metadata.ServiceTemplate, 0)
		if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(filter).Fields(common.BKFieldID).
			All(kit.Ctx, &serviceTemplates); err != nil {
			blog.Errorf("get serviceTemplate failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
			return enableModuleMap, []int64{}, srvTempModulesMap, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		for _, serviceTemplate := range serviceTemplates {
			existSrvTempIDsMap[serviceTemplate.ID] = struct{}{}

		}
	}

	// list of valid host auto-apply module configurations
	haveHostApplyModuleIDs := make([]int64, 0)

	// srvTempModulesMap: at present, there is a module without a template, or the host
	// automatic application of the template corresponding to the module is turned off
	for _, module := range modules {
		if _, ok := existSrvTempIDsMap[module.ServiceTemplateID]; ok {
			srvTempModulesMap[module.ServiceTemplateID] = append(srvTempModulesMap[module.ServiceTemplateID],
				module.ModuleID)
			enableModuleMap[module.ModuleID] = struct{}{}

			continue
		}
		if moduleIDHostApplyEnabledMap[module.ModuleID] {
			enableModuleMap[module.ModuleID] = struct{}{}
			haveHostApplyModuleIDs = append(haveHostApplyModuleIDs, module.ModuleID)
		}
	}

	return enableModuleMap, haveHostApplyModuleIDs, srvTempModulesMap, nil
}

// getFinalRules If there are multiple target modules to be transferred, the configuration values in the template
// will be preferred. If a property exists in multiple templates, the value will be randomly selected
func (p *hostApplyRule) getFinalRules(kit *rest.Kit, bizID int64, haveHostApplyIDs, serviceTemplateIDs []int64,
	srvTemplateIDMap map[int64][]int64) ([]metadata.HostApplyRule, errors.CCErrorCoder) {

	finalRules := make([]metadata.HostApplyRule, 0)

	if len(serviceTemplateIDs) > 0 {
		srvTempRuleOp := metadata.ListHostApplyRuleOption{
			ServiceTemplateIDs: serviceTemplateIDs,
			Page:               metadata.BasePage{Limit: common.BKNoLimit},
		}
		srvTempRules, ccErr := p.ListHostApplyRule(kit, bizID, srvTempRuleOp)
		if ccErr != nil {
			blog.Errorf("list service template host apply rule failed, opt: %v, err: %v, rid: %s", srvTempRuleOp,
				ccErr, kit.Rid)
			return finalRules, ccErr
		}

		for _, rule := range srvTempRules.Info {
			for _, moduleID := range srvTemplateIDMap[rule.ServiceTemplateID] {
				rule.ModuleID = moduleID
				finalRules = append(finalRules, rule)
			}
		}
	}

	if len(haveHostApplyIDs) > 0 {
		moduleRuleOp := metadata.ListHostApplyRuleOption{
			ModuleIDs: haveHostApplyIDs,
			Page:      metadata.BasePage{Limit: common.BKNoLimit},
		}
		moduleRules, ccErr := p.ListHostApplyRule(kit, bizID, moduleRuleOp)
		if ccErr != nil {
			blog.Errorf("list module host apply rule failed, opt: %v, err: %v, rid: %s", moduleRuleOp, ccErr, kit.Rid)
			return finalRules, ccErr
		}
		finalRules = append(finalRules, moduleRules.Info...)
	}

	return finalRules, nil
}

// RunHostApplyOnHosts run host apply rule on specified host
func (p *hostApplyRule) RunHostApplyOnHosts(kit *rest.Kit, bizID int64, relations []metadata.ModuleHost) (
	metadata.MultipleHostApplyResult, errors.CCErrorCoder) {

	result := metadata.MultipleHostApplyResult{HostResults: make([]metadata.HostApplyResult, 0)}

	moduleIDs := make([]int64, 0)
	for _, item := range relations {
		moduleIDs = append(moduleIDs, item.ModuleID)
	}

	modules := make([]metadata.ModuleInst, 0)
	moduleFilter := map[string]interface{}{common.BKModuleIDField: map[string]interface{}{common.BKDBIN: moduleIDs}}

	fields := []string{common.BKModuleIDField, common.BKServiceTemplateIDField, common.HostApplyEnabledField}
	err := mongodb.Client().Table(common.BKTableNameBaseModule).Find(moduleFilter).Fields(fields...).
		All(kit.Ctx, &modules)
	if err != nil {
		blog.Errorf("search modules info failed, filter: %s, err: %v, rid: %s", moduleFilter, err, kit.Rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	enableModuleMap, haveHostApplyIDs, srvTemplateIDMap, cErr := getModuleIDsAndSrvTempIDs(kit, modules)
	if err != nil {
		return result, cErr
	}

	host2Modules := make(map[int64][]int64)
	for _, relation := range relations {
		if _, exist := enableModuleMap[relation.ModuleID]; !exist {
			continue
		}
		if _, exist := host2Modules[relation.HostID]; !exist {
			host2Modules[relation.HostID] = make([]int64, 0)
		}
		// checkout host apply enabled status on module
		host2Modules[relation.HostID] = append(host2Modules[relation.HostID], relation.ModuleID)
	}
	hostModules := make([]metadata.Host2Modules, 0)
	for hostID, moduleIDs := range host2Modules {
		hostModules = append(hostModules, metadata.Host2Modules{
			HostID:    hostID,
			ModuleIDs: moduleIDs})
	}

	serviceTemplateIDs := make([]int64, 0)
	for serviceTemplateID := range srvTemplateIDMap {
		serviceTemplateIDs = append(serviceTemplateIDs, serviceTemplateID)
	}

	finalRules, cErr := p.getFinalRules(kit, bizID, haveHostApplyIDs, serviceTemplateIDs, srvTemplateIDMap)
	if cErr != nil {
		return result, cErr
	}

	if len(finalRules) == 0 {
		return result, nil
	}

	planOption := metadata.HostApplyPlanOption{
		Rules:       finalRules,
		HostModules: hostModules,
	}

	planResult, ccErr := p.GenerateApplyPlan(kit, bizID, planOption)
	if ccErr != nil {
		blog.Errorf("generate apply plan failed, option: %v, err: %v, rid: %s", planOption, ccErr, kit.Rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	result.HostResults = p.carryOutPlan(kit, planResult.Plans)
	for _, hostResult := range result.HostResults {
		if ccErr := hostResult.GetError(); ccErr != nil {
			result.SetError(ccErr)
			break
		}
	}
	return result, result.GetError()
}

type updateHostOption struct {
	hostIDs []int64
	data    map[string]interface{}
}

func (p *hostApplyRule) carryOutPlan(kit *rest.Kit, plans []metadata.OneHostApplyPlan) []metadata.HostApplyResult {
	hostResults := make([]metadata.HostApplyResult, 0)

	// group hosts with the same host apply update data together, updateDataMap key is the json format of update data
	updateDataMap := make(map[string]*updateHostOption)
	for _, plan := range plans {
		if len(plan.UpdateFields) == 0 {
			hostResults = append(hostResults, metadata.HostApplyResult{HostID: plan.HostID})
			continue
		}

		dataStr := plan.GetUpdateDataStr()
		if _, exists := updateDataMap[dataStr]; !exists {
			updateDataMap[dataStr] = &updateHostOption{
				hostIDs: make([]int64, 0),
				data:    plan.GetUpdateData(),
			}
		}
		updateDataMap[dataStr].hostIDs = append(updateDataMap[dataStr].hostIDs, plan.HostID)
	}

	// batch update the hosts with the same update data together to improve performance
	for _, updateOpt := range updateDataMap {
		updateOption := metadata.UpdateOption{
			Data: updateOpt.data,
			Condition: map[string]interface{}{
				common.BKHostIDField: map[string]interface{}{common.BKDBIN: updateOpt.hostIDs},
			},
		}

		if _, err := p.dependence.UpdateModelInstance(kit, common.BKInnerObjIDHost, updateOption); err != nil {
			blog.Warnf("update host failed, updateOption: %+v, err: %v, rid: %s", updateOption, err, kit.Rid)
			ccErr, ok := err.(errors.CCErrorCoder)
			if !ok {
				ccErr = kit.CCError.CCError(common.CCErrHostUpdateFail)
			}
			applyResult := metadata.HostApplyResult{}
			applyResult.SetError(ccErr)
			for _, hostID := range updateOpt.hostIDs {
				applyResult.HostID = hostID
				hostResults = append(hostResults, applyResult)
			}
			continue
		}
		for _, hostID := range updateOpt.hostIDs {
			hostResults = append(hostResults, metadata.HostApplyResult{HostID: hostID})
		}
	}

	return hostResults
}
