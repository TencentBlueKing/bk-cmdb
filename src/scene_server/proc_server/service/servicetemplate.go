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

package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

// CreateServiceTemplate TODO
func (ps *ProcServer) CreateServiceTemplate(ctx *rest.Contexts) {
	option := new(metadata.CreateServiceTemplateOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	newTemplate := &metadata.ServiceTemplate{
		BizID:             option.BizID,
		Name:              option.Name,
		ServiceCategoryID: option.ServiceCategoryID,
		SupplierAccount:   ctx.Kit.SupplierAccount,
		HostApplyEnabled:  option.HostApplyEnabled,
	}

	var tpl *metadata.ServiceTemplate
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		tpl, err = ps.CoreAPI.CoreService().Process().CreateServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, newTemplate)
		if err != nil {
			blog.Errorf("create service template failed, err: %v", err)
			return err
		}

		// register service template resource creator action to iam
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.BizProcessServiceTemplate),
				ID:      strconv.FormatInt(tpl.ID, 10),
				Name:    tpl.Name,
				Creator: ctx.Kit.User,
			}
			_, err = ps.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created service template to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(tpl)
}

// CreateServiceTemplateAllInfo create service template all info, including attributes and process templates
func (ps *ProcServer) CreateServiceTemplateAllInfo(ctx *rest.Contexts) {
	option := new(metadata.CreateSvcTempAllInfoOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var templateID int64
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// create service template
		template := &metadata.ServiceTemplate{
			BizID:             option.BizID,
			Name:              option.Name,
			ServiceCategoryID: option.ServiceCategoryID,
		}

		tpl, err := ps.CoreAPI.CoreService().Process().CreateServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, template)
		if err != nil {
			blog.Errorf("create service template(%+v) failed, err: %v, rid: %s", template, err, ctx.Kit.Rid)
			return err
		}

		templateID = tpl.ID

		// create service template attributes
		if len(option.Attributes) > 0 {
			attrOpt := &metadata.CreateSvcTempAttrsOption{
				BizID:             option.BizID,
				ServiceTemplateID: templateID,
				Attributes:        option.Attributes,
			}

			_, err = ps.CoreAPI.CoreService().Process().CreateServiceTemplateAttrs(ctx.Kit.Ctx, ctx.Kit.Header, attrOpt)
			if err != nil {
				blog.Errorf("create service template attrs(%+v) failed, err: %v, rid: %s", attrOpt, err, ctx.Kit.Rid)
				return err
			}
		}

		// create process templates
		for _, process := range option.Processes {
			t := &metadata.ProcessTemplate{
				BizID:             option.BizID,
				ServiceTemplateID: templateID,
				Property:          process.Property,
			}

			_, err := ps.CoreAPI.CoreService().Process().CreateProcessTemplate(ctx.Kit.Ctx, ctx.Kit.Header, t)
			if err != nil {
				blog.Errorf("create process template(%+v) failed, err: %v, rid: %s", t, err, ctx.Kit.Rid)
				return err
			}
		}

		// register service template resource creator action to iam
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.BizProcessServiceTemplate),
				ID:      strconv.FormatInt(tpl.ID, 10),
				Name:    tpl.Name,
				Creator: ctx.Kit.User,
			}

			_, err := ps.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created service template to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(metadata.RspID{ID: templateID})
}

// GetServiceTemplate TODO
func (ps *ProcServer) GetServiceTemplate(ctx *rest.Contexts) {
	templateIDStr := ctx.Request.PathParameter(common.BKServiceTemplateIDField)
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if err != nil {
		ctx.RespErrorCodeF(common.CCErrCommParamsInvalid, "create service template failed, err: %v", common.BKServiceTemplateIDField, err)
		return
	}
	template, err := ps.CoreAPI.CoreService().Process().GetServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, templateID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "get service template failed, err: %v", err)
		return
	}

	ctx.RespEntity(template)
}

// GetServiceTemplateDetail return more info than GetServiceTemplate
func (ps *ProcServer) GetServiceTemplateDetail(ctx *rest.Contexts) {
	templateIDStr := ctx.Request.PathParameter(common.BKServiceTemplateIDField)
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if err != nil {
		ctx.RespErrorCodeF(common.CCErrCommParamsInvalid, "create service template failed, err: %v", common.BKServiceTemplateIDField, err)
		return
	}
	templateDetail, err := ps.CoreAPI.CoreService().Process().GetServiceTemplateWithStatistics(ctx.Kit.Ctx, ctx.Kit.Header, templateID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "get service template failed, err: %v", err)
		return
	}

	ctx.RespEntity(templateDetail)
}

// GetServiceTemplateAllInfo get service template all info, including attributes and process templates
func (ps *ProcServer) GetServiceTemplateAllInfo(ctx *rest.Contexts) {
	option := new(metadata.GetSvcTempAllInfoOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if option.BizID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAppIDField))
		return
	}

	if option.ID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKFieldID))
		return
	}

	allInfo, err := ps.getServiceTemplateAllInfo(ctx.Kit, option.ID, option.BizID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(allInfo)
}

func (ps *ProcServer) getServiceTemplateAllInfo(kit *rest.Kit, id, bizID int64) (*metadata.SvcTempAllInfo,
	errors.CCErrorCoder) {

	// get service template
	svcTemp, err := ps.CoreAPI.CoreService().Process().GetServiceTemplate(kit.Ctx, kit.Header, id)
	if err != nil {
		blog.Errorf("get service template %d failed, err: %v, rid: %s", id, err, kit.Rid)
		return nil, err
	}

	if svcTemp.BizID != bizID {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	// get service template attributes
	attrOpt := &metadata.ListServTempAttrOption{
		BizID: bizID,
		ID:    id,
	}
	attrs, err := ps.CoreAPI.CoreService().Process().ListServiceTemplateAttribute(kit.Ctx, kit.Header, attrOpt)
	if err != nil {
		blog.Errorf("get service template %d attributes failed, err: %v, rid: %s", id, err, kit.Rid)
		return nil, err
	}

	// get process templates
	procTempOpt := &metadata.ListProcessTemplatesOption{
		BusinessID:         bizID,
		ServiceTemplateIDs: []int64{id},
	}

	procTemps, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(kit.Ctx, kit.Header, procTempOpt)
	if err != nil {
		blog.Errorf("list process templates failed, input: %+v, err: %v, rid: %s", procTempOpt, err, kit.Rid)
		return nil, err
	}

	return &metadata.SvcTempAllInfo{
		ID:                svcTemp.ID,
		BizID:             svcTemp.BizID,
		Name:              svcTemp.Name,
		ServiceCategoryID: svcTemp.ServiceCategoryID,
		Attributes:        attrs.Attributes,
		Processes:         procTemps.Info,
	}, nil
}

func (ps *ProcServer) getHostIDByCondition(kit *rest.Kit, bizID int64, serviceTemplateIDs []int64,
	hostIDs []int64) ([]int64, errors.CCErrorCoder) {

	// 1、get module ids by template ids.
	moduleCond := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}

	if len(serviceTemplateIDs) > 0 {
		moduleCond[common.BKServiceTemplateIDField] = mapstr.MapStr{common.BKDBIN: serviceTemplateIDs}
	}

	moduleFilter := &metadata.QueryCondition{
		Condition:      moduleCond,
		Fields:         []string{common.BKModuleIDField},
		DisableCounter: true,
	}

	moduleRes := new(metadata.ResponseModuleInstance)
	err := ps.CoreAPI.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		moduleFilter, moduleRes)
	if err != nil {
		blog.Errorf("get module failed, filter: %#v, err: %v, rid: %s", moduleFilter, err, kit.Rid)
		return nil, err
	}
	if err := moduleRes.CCError(); err != nil {
		blog.Errorf("get module failed, filter: %#v, err: %v, rid: %s", moduleFilter, err, kit.Rid)
		return nil, err
	}
	modIDs := make([]int64, 0)

	// need to be compatible with scenarios without modules under the service template. In this scenario, only attribute
	// rules need to be applied to service templates, not host attribute rules (without modules, there must be no hosts)
	if len(moduleRes.Data.Info) == 0 {
		return modIDs, nil
	}

	for _, modID := range moduleRes.Data.Info {
		modIDs = append(modIDs, modID.ModuleID)
	}

	// 2、get the corresponding hostIDs list through the module ids.
	relReq := &metadata.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{bizID},
	}
	if hostIDs != nil {
		relReq.HostIDArr = hostIDs
	}
	if len(modIDs) > 0 {
		relReq.ModuleIDArr = modIDs
	}

	relRsp, relErr := ps.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(kit.Ctx, kit.Header, relReq)
	if relErr != nil {
		blog.Errorf("get host ids failed, req: %s, err: %s, rid: %s", relReq, relErr, kit.Rid)
		return relRsp, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	return relRsp, nil

}

// ExecServiceTemplateHostApplyRule execute the host automatic application task in the template scenario.
func (ps *ProcServer) ExecServiceTemplateHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	planReq := new(metadata.HostApplyServiceTemplateOption)
	if err := ctx.DecodeInto(planReq); err != nil {
		ctx.RespAutoError(err)
		return
	}
	hostIDs, err := ps.getHostIDByCondition(ctx.Kit, planReq.BizID, planReq.ServiceTemplateIDs, planReq.HostIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// enable host apply on service template
		updateOption := &metadata.UpdateOption{
			Condition: map[string]interface{}{
				common.BKFieldID:    map[string]interface{}{common.BKDBIN: planReq.ServiceTemplateIDs},
				common.BKAppIDField: planReq.BizID,
			},
			Data: map[string]interface{}{common.HostApplyEnabledField: true},
		}

		err := ps.CoreAPI.CoreService().Process().UpdateBatchServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, updateOption)
		if err != nil {
			blog.Errorf("update service template failed, err: %v", err)
			return err
		}

		// 1、update or add rules.
		rulesOption := make([]metadata.CreateOrUpdateApplyRuleOption, 0)
		for _, rule := range planReq.AdditionalRules {
			rulesOption = append(rulesOption, metadata.CreateOrUpdateApplyRuleOption{
				AttributeID:       rule.AttributeID,
				ServiceTemplateID: rule.ServiceTemplateID,
				PropertyValue:     rule.PropertyValue,
			})
		}
		saveRuleOp := metadata.BatchCreateOrUpdateApplyRuleOption{Rules: rulesOption}
		if _, ccErr := ps.CoreAPI.CoreService().HostApplyRule().BatchUpdateHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header,
			planReq.BizID, saveRuleOp); ccErr != nil {
			blog.Errorf("update host rule failed, bizID: %s, req: %s, err: %v, rid: %s", planReq.BizID, saveRuleOp,
				ccErr, rid)
			return ccErr
		}

		// 2、delete rules.
		if len(planReq.RemoveRuleIDs) > 0 {
			removeOp := metadata.DeleteHostApplyRuleOption{
				RuleIDs:            planReq.RemoveRuleIDs,
				ServiceTemplateIDs: planReq.ServiceTemplateIDs,
			}
			if ccErr := ps.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header,
				planReq.BizID, removeOp); ccErr != nil {
				blog.Errorf("delete apply rule failed, bizID: %d, req: %s, err: %v, rid: %s", planReq.BizID, removeOp,
					ccErr, rid)
				return ccErr
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(&metadata.RespError{Msg: txnErr})
		return
	}

	// The following three scenarios do not require the update of the host properties to be automatically applied:
	// 1. The changed flag is false.
	// 2. This request only deletes the rule scenario.
	// 3. No module is created under the service template or there is no eligible host under the module.
	if !planReq.Changed || len(planReq.AdditionalRules) == 0 || len(hostIDs) == 0 {
		ctx.RespEntity(nil)
		return
	}

	// update host operation is not done in a transaction, since the successfully updated hosts need not roll back
	ctx.Kit.Header.Del(common.TransactionIdHeader)

	// host apply attribute rules to the host.
	err = ps.updateHostAttributes(ctx.Kit, planReq, hostIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *ProcServer) getUpdateDataStrByApplyRule(kit *rest.Kit, rules []metadata.CreateHostApplyRuleOption) (
	string, errors.CCErrorCoder) {
	attributeIDs := make([]int64, 0)
	attrIDMap := make(map[int64]struct{})
	for _, rule := range rules {
		if _, ok := attrIDMap[rule.AttributeID]; ok {
			continue
		}
		attrIDMap[rule.AttributeID] = struct{}{}
		attributeIDs = append(attributeIDs, rule.AttributeID)
	}

	attCond := &metadata.QueryCondition{
		Fields: []string{common.BKFieldID, common.BKPropertyIDField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKFieldID: map[string]interface{}{
				common.BKDBIN: attributeIDs,
			},
		},
	}

	attrRes, err := s.CoreAPI.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDHost, attCond)
	if err != nil {
		blog.Errorf("read model attr failed, err: %v, attrCond: %#v, rid: %s", err, attCond, kit.Rid)
		return "", kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	attrMap := make(map[int64]string)
	for _, attr := range attrRes.Info {
		attrMap[attr.ID] = attr.PropertyID
	}

	fields := make([]string, len(rules))

	for index, field := range rules {
		value, _ := json.Marshal(field.PropertyValue)
		fields[index] = fmt.Sprintf(`"%s":%s`, attrMap[field.AttributeID], string(value))
	}

	sort.Strings(fields)
	return "{" + strings.Join(fields, ",") + "}", nil
}

func generateCondition(dataStr string, hostIDs []int64) (map[string]interface{}, map[string]interface{}) {
	data := make(map[string]interface{})
	_ = json.Unmarshal([]byte(dataStr), &data)

	cond := make([]map[string]interface{}, 0)

	for key, value := range data {
		cond = append(cond, map[string]interface{}{
			key: map[string]interface{}{common.BKDBNE: value},
		})
	}
	mergeCond := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs},
		common.BKDBOR:        cond,
	}
	return mergeCond, data
}

func (s *ProcServer) updateHostAttributes(kit *rest.Kit, planResult *metadata.HostApplyServiceTemplateOption,
	hostIDs []int64) errors.CCErrorCoder {

	dataStr, err := s.getUpdateDataStrByApplyRule(kit, planResult.AdditionalRules)
	if err != nil {
		return err
	}
	mergeCond, data := generateCondition(dataStr, hostIDs)
	counts, cErr := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameBaseHost, []map[string]interface{}{mergeCond})
	if cErr != nil {
		blog.Errorf("get hosts count failed, filter: %+v, err: %v, rid: %s", mergeCond, cErr, kit.Rid)
		return cErr
	}
	if counts[0] == 0 {
		blog.V(5).Infof("no hosts founded, filter: %+v, rid: %s", mergeCond, kit.Rid)
		return nil
	}

	// If there is no eligible host, then return directly.
	updateOp := &metadata.UpdateOption{Data: data, Condition: mergeCond}

	_, e := s.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDHost, updateOp)
	if e != nil {
		blog.Errorf("update host failed, option: %s, err: %v, rid: %s", updateOp, e, kit.Rid)
		return errors.New(common.CCErrCommHTTPDoRequestFailed, e.Error())
	}
	return nil
}

// UpdateServiceTemplateHostApplyRule update host auto-apply rules in service template dimension.
func (ps *ProcServer) UpdateServiceTemplateHostApplyRule(ctx *rest.Contexts) {

	syncOpt := new(metadata.HostApplyServiceTemplateOption)
	if err := ctx.DecodeInto(syncOpt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := syncOpt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	taskInfo := metadata.APITaskDetail{}

	// The host is automatically updated asynchronously in the application scenario. The instID corresponds to the
	// BizID, but if the task is created according to the business level, a large number of task conflict scenarios will
	// appear. This scenario allows repeated execution of the same task, and only the execution result of the last task
	// is retained. When querying the task result, the history api can be used without passing the instID. Therefore,
	// the instID here can be assigned a random number. Random instID from 10000 to 20000 in template scene.
	randInstNum := util.RandInt64WithRange(int64(10000), int64(20000))

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		taskRes, err := ps.CoreAPI.TaskServer().Task().Create(ctx.Kit.Ctx, ctx.Kit.Header,
			common.SyncServiceTemplateHostApplyTaskFlag, randInstNum, []interface{}{syncOpt})
		if err != nil {
			blog.Errorf("create service template host apply sync rule task failed, opt: %+v, err: %v, rid: %s",
				syncOpt, err, ctx.Kit.Rid)
			return err
		}
		taskInfo = taskRes
		blog.V(4).Infof("successfully create service template host apply sync task: %#v, rid: %s", taskRes, ctx.Kit.Rid)
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(metadata.HostApplyTaskResult{BizID: taskInfo.InstID, TaskID: taskInfo.TaskID})
}

// UpdateServiceTemplateHostApplyEnableStatus update object host if apply's status is enabled
func (ps *ProcServer) UpdateServiceTemplateHostApplyEnableStatus(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("parse bk_biz_id failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	requestBody := metadata.UpdateHostApplyEnableStatusOption{}
	if err := ctx.DecodeInto(&requestBody); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := requestBody.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	updateOption := &metadata.UpdateOption{
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKFieldID:    mapstr.MapStr{common.BKDBIN: requestBody.IDs},
		},
		Data: map[string]interface{}{
			common.HostApplyEnabledField: requestBody.Enable,
		},
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := ps.CoreAPI.CoreService().Process().UpdateBatchServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, updateOption)
		if err != nil {
			blog.Errorf("update service template failed, err: %v", err)
			return err
		}

		// in the scenario of turning on the host's automatic application state, there is no clear rule action, and
		// return directly
		if requestBody.Enable {
			return nil
		}

		if requestBody.ClearRules {
			listRuleOption := metadata.ListHostApplyRuleOption{
				ServiceTemplateIDs: requestBody.IDs,
				Page: metadata.BasePage{
					Limit: common.BKNoLimit,
				},
			}
			listRuleResult, ccErr := ps.Engine.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx,
				ctx.Kit.Header, bizID, listRuleOption)
			if ccErr != nil {
				blog.Errorf("get list host apply rule failed, bizID: %d,listRuleOption: %#v, rid: %s", bizID,
					listRuleOption, ctx.Kit.Rid)
				return ccErr
			}
			ruleIDs := make([]int64, 0)
			for _, item := range listRuleResult.Info {
				ruleIDs = append(ruleIDs, item.ID)
			}
			if len(ruleIDs) > 0 {
				deleteRuleOption := metadata.DeleteHostApplyRuleOption{
					RuleIDs:            ruleIDs,
					ServiceTemplateIDs: requestBody.IDs,
				}
				if ccErr := ps.Engine.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx,
					ctx.Kit.Header, bizID, deleteRuleOption); ccErr != nil {
					blog.Errorf("delete list host apply rule failed, bizID: %d, listRuleOption: %#v, rid: %s",
						bizID, listRuleOption, ctx.Kit.Rid)
					return ccErr
				}
			}
		}
		return nil
	})
	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// GetHostApplyTaskStatus get host auto-apply asynchronous task status.
func (s *ProcServer) GetHostApplyTaskStatus(ctx *rest.Contexts) {

	syncStatusOpt := new(metadata.HostApplyTaskStatusOption)
	if err := ctx.DecodeInto(syncStatusOpt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if rawErr := syncStatusOpt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// get host auto-apply task status by task ids. query the automatic application status of the host. Since the instID
	// when creating a task is a random number, the instID input condition is not required when querying.
	statusOpt := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKTaskTypeField: common.SyncServiceTemplateHostApplyTaskFlag,
			common.BKTaskIDField:   map[string]interface{}{common.BKDBIN: syncStatusOpt.TaskIDs},
		},
		Fields:         []string{common.BKStatusField, common.BKTaskIDField},
		DisableCounter: true,
	}

	tasksStatus, err := s.CoreAPI.TaskServer().Task().ListSyncStatusHistory(ctx.Kit.Ctx, ctx.Kit.Header, statusOpt)
	if err != nil {
		blog.Errorf("list sync status history failed, option: %#v, err: %v, rid: %s", statusOpt, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	result := metadata.HostApplyTaskStatusRsp{
		BizID: syncStatusOpt.BizID,
	}
	for _, task := range tasksStatus.Info {
		result.TaskInfo = append(result.TaskInfo, metadata.HostAppyTaskInfo{
			TaskID: task.TaskID,
			Status: string(task.Status),
		})
	}
	ctx.RespEntity(result)
	return
}

// DeleteHostApplyRule delete the host automatic application rule in the service template scenario.
func (ps *ProcServer) DeleteHostApplyRule(ctx *rest.Contexts) {

	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}
	option := metadata.DeleteHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := option.ValidateServiceTemplateOption(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := ps.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID,
			option); err != nil {
			blog.Errorf("DeleteHostApplyRule failed, core service DeleteHostApplyRule failed, bizID: %s, option: %s,"+
				" err: %v, rid: %s", bizID, option, err, rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)

}

// UpdateServiceTemplate TODO
func (ps *ProcServer) UpdateServiceTemplate(ctx *rest.Contexts) {
	option := new(metadata.UpdateServiceTemplateOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	updateParam := &metadata.ServiceTemplate{
		ID:                option.ID,
		Name:              option.Name,
		ServiceCategoryID: option.ServiceCategoryID,
	}

	var tpl *metadata.ServiceTemplate
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		tpl, err = ps.CoreAPI.CoreService().Process().UpdateServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, option.ID, updateParam)
		if err != nil {
			blog.Errorf("update service template failed, err: %v", err)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(tpl)
}

// UpdateServiceTemplateAllInfo update service template all info, including attributes and process templates
func (ps *ProcServer) UpdateServiceTemplateAllInfo(ctx *rest.Contexts) {
	option := new(metadata.UpdateSvcTempAllInfoOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	allInfo, err := ps.getServiceTemplateAllInfo(ctx.Kit, option.ID, option.BizID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// update service template
		if option.Name != allInfo.Name || option.ServiceCategoryID != allInfo.ServiceCategoryID {
			opt := &metadata.ServiceTemplate{
				Name:              option.Name,
				ServiceCategoryID: option.ServiceCategoryID,
			}

			if _, err := ps.CoreAPI.CoreService().Process().UpdateServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header,
				option.ID, opt); err != nil {
				blog.Errorf("update svc temp %d failed, opt: %+v, err: %v, rid: %s", option.ID, opt, err, ctx.Kit.Rid)
				return err
			}
		}

		// update service template attributes
		err = ps.updateSvcTempAllAttrs(ctx.Kit, allInfo.ID, allInfo.BizID, allInfo.Attributes, option.Attributes)
		if err != nil {
			return err
		}

		// update process templates
		err = ps.updateSvcTempAllProcTemps(ctx.Kit, allInfo.ID, allInfo.BizID, allInfo.Processes, option.Processes)
		if err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// updateSvcTempAllAttrs update service template attributes, add new attributes and delete redundant attributes
func (ps *ProcServer) updateSvcTempAllAttrs(kit *rest.Kit, id, bizID int64, prevAttrs []metadata.ServiceTemplateAttr,
	updateAttrs []metadata.SvcTempAttr) errors.CCErrorCoder {

	attrMap := make(map[int64]interface{})
	for _, attribute := range prevAttrs {
		attrMap[attribute.AttributeID] = attribute
	}

	// cross compare previous attributes and update attributes to find need add/update/delete attributes
	var addedAttrs, updatedAttrs []metadata.SvcTempAttr
	for _, attribute := range updateAttrs {
		value, exists := attrMap[attribute.AttributeID]
		if !exists {
			addedAttrs = append(addedAttrs, metadata.SvcTempAttr{
				AttributeID:   attribute.AttributeID,
				PropertyValue: attribute.PropertyValue,
			})
			continue
		}

		delete(attrMap, attribute.AttributeID)
		if !reflect.DeepEqual(value, attribute.PropertyValue) {
			updatedAttrs = append(updatedAttrs, metadata.SvcTempAttr{
				AttributeID:   attribute.AttributeID,
				PropertyValue: attribute.PropertyValue,
			})
		}
	}

	// delete service template attributes
	if len(attrMap) > 0 {
		deletedAttrIDs := make([]int64, 0)
		for attrID := range attrMap {
			deletedAttrIDs = append(deletedAttrIDs, attrID)
		}

		deleteOpt := &metadata.DeleteServTempAttrOption{
			BizID:        bizID,
			ID:           id,
			AttributeIDs: deletedAttrIDs,
		}
		err := ps.CoreAPI.CoreService().Process().DeleteServiceTemplateAttribute(kit.Ctx, kit.Header, deleteOpt)
		if err != nil {
			blog.Errorf("delete service template attrs failed, opt: %+v, err: %v, rid: %s", deleteOpt, err, kit.Rid)
			return err
		}
	}

	// add service template attributes
	if len(addedAttrs) > 0 {
		addOpt := &metadata.CreateSvcTempAttrsOption{
			BizID:             bizID,
			ServiceTemplateID: id,
			Attributes:        addedAttrs,
		}

		_, err := ps.CoreAPI.CoreService().Process().CreateServiceTemplateAttrs(kit.Ctx, kit.Header, addOpt)
		if err != nil {
			blog.Errorf("add service template attrs failed, opt: %+v, err: %v, rid: %s", addOpt, err, kit.Rid)
			return err
		}
	}

	// update service template attributes
	if len(updatedAttrs) > 0 {
		updateOpt := &metadata.UpdateServTempAttrOption{
			BizID:      bizID,
			ID:         id,
			Attributes: updatedAttrs,
		}

		err := ps.CoreAPI.CoreService().Process().UpdateServiceTemplateAttribute(kit.Ctx, kit.Header, updateOpt)
		if err != nil {
			blog.Errorf("update service template attrs failed, opt: %+v, err: %v, rid: %s", updateOpt, err, kit.Rid)
			return err
		}
	}

	return nil
}

// updateSvcTempAllProcTemps update service template procTemps, add new procTemps and delete redundant procTemps
func (ps *ProcServer) updateSvcTempAllProcTemps(kit *rest.Kit, id, bizID int64, prevProcTemps,
	updateProcTemps []metadata.ProcessTemplate) errors.CCErrorCoder {

	procTempMap := make(map[int64]*metadata.ProcessProperty)
	for _, procTemp := range prevProcTemps {
		procTempMap[procTemp.ID] = procTemp.Property
	}

	// cross compare previous procTemps and update procTemps to find need add/update/delete procTemps
	var addedProcTemps, updatedProcTemps []metadata.ProcessTemplate
	for _, procTemp := range updateProcTemps {
		value, exists := procTempMap[procTemp.ID]
		if !exists {
			procTemp.BizID = bizID
			procTemp.ServiceTemplateID = id
			addedProcTemps = append(addedProcTemps, procTemp)
			continue
		}

		delete(procTempMap, procTemp.ID)
		if !reflect.DeepEqual(value, procTemp.Property) {
			updatedProcTemps = append(updatedProcTemps, procTemp)
		}
	}

	// delete service template procTemps
	for procTempID := range procTempMap {
		err := ps.CoreAPI.CoreService().Process().DeleteProcessTemplate(kit.Ctx, kit.Header, procTempID)
		if err != nil {
			blog.Errorf("delete process template %d failed, err: %v, rid: %s", procTempID, err, kit.Rid)
			return err
		}
	}

	// add service template procTemps
	for _, procTemp := range addedProcTemps {
		_, err := ps.CoreAPI.CoreService().Process().CreateProcessTemplate(kit.Ctx, kit.Header, &procTemp)
		if err != nil {
			blog.Errorf("add process template(%+v) failed, err: %v, rid: %s", procTemp, err, kit.Rid)
			return err
		}
	}

	// update service template procTemps
	for _, procTemp := range updatedProcTemps {
		property, rawErr := mapstr.Struct2Map(procTemp.Property)
		if rawErr != nil {
			blog.Errorf("convert proc temp property(%+v) failed, err: %v, rid: %s", procTemp.Property, rawErr, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommJSONMarshalFailed)
		}

		_, err := ps.CoreAPI.CoreService().Process().UpdateProcessTemplate(kit.Ctx, kit.Header, procTemp.ID, property)
		if err != nil {
			blog.Errorf("update process template(%+v) failed, err: %v, rid: %s", procTemp, err, kit.Rid)
			return err
		}
	}

	return nil
}

// ListServiceTemplates TODO
func (ps *ProcServer) ListServiceTemplates(ctx *rest.Contexts) {
	input := new(metadata.ListServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if input.Page.IsIllegal() {
		ctx.RespErrorCodeOnly(common.CCErrCommPageLimitIsExceeded, "list service template, but page limit:%d is over "+
			"limited.", input.Page.Limit)
		return
	}

	option := metadata.ListServiceTemplateOption{
		BusinessID:         input.BizID,
		Page:               input.Page,
		ServiceCategoryID:  &input.ServiceCategoryID,
		Search:             input.Search,
		IsExact:            input.IsExact,
		ServiceTemplateIDs: input.ServiceTemplateIDs,
	}
	temp, err := ps.CoreAPI.CoreService().Process().ListServiceTemplates(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "list service template failed, input: %+v", input)
		return
	}

	ctx.RespEntity(temp)
}

// FindServiceTemplateCountInfo find count info of service templates
func (ps *ProcServer) FindServiceTemplateCountInfo(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("FindServiceTemplateCountInfo failed, parse bk_biz_id error, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	input := new(metadata.FindServiceTemplateCountInfoOption)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// generate count conditions
	filters := make([]map[string]interface{}, len(input.ServiceTemplateIDs))
	for idx, serviceTemplateID := range input.ServiceTemplateIDs {
		filters[idx] = map[string]interface{}{
			common.BKAppIDField:             bizID,
			common.BKServiceTemplateIDField: serviceTemplateID,
		}
	}

	// process templates reference count
	processTemplateCounts, err := ps.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header, common.BKTableNameProcessTemplate, filters)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed, "count process template by filters: %+v failed.", filters)
		return
	}
	if len(processTemplateCounts) != len(input.ServiceTemplateIDs) {
		ctx.RespWithError(ctx.Kit.CCError.CCError(common.CCErrProcGetProcessTemplatesFailed), common.CCErrProcGetProcessTemplatesFailed,
			"the count of process must be equal with the count of service templates, filters:%#v", filters)
		return
	}

	// module reference count
	moduleCounts, err := ps.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header, common.BKTableNameBaseModule, filters)
	if err != nil {
		ctx.RespWithError(err, common.CCErrTopoModuleSelectFailed, "count process template by filters: %+v failed.", filters)
		return
	}
	if len(moduleCounts) != len(input.ServiceTemplateIDs) {
		ctx.RespWithError(ctx.Kit.CCError.CCError(common.CCErrTopoModuleSelectFailed), common.CCErrTopoModuleSelectFailed,
			"the count of modules must be equal with the count of service templates, filters:%#v", filters)
		return
	}

	// service instance reference count
	serviceInstanceCounts, err := ps.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header, common.BKTableNameServiceInstance, filters)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "count process template by filters: %+v failed.", filters)
		return
	}
	if len(serviceInstanceCounts) != len(input.ServiceTemplateIDs) {
		ctx.RespWithError(ctx.Kit.CCError.CCError(common.CCErrProcGetServiceInstancesFailed), common.CCErrProcGetServiceInstancesFailed,
			"the count of service instance must be equal with the count of service templates, filters:%#v", filters)
		return
	}

	result := make([]metadata.FindServiceTemplateCountInfoResult, 0)
	for idx, serviceTemplateID := range input.ServiceTemplateIDs {
		result = append(result, metadata.FindServiceTemplateCountInfoResult{
			ServiceTemplateID:    serviceTemplateID,
			ProcessTemplateCount: processTemplateCounts[idx],
			ServiceInstanceCount: serviceInstanceCounts[idx],
			ModuleCount:          moduleCounts[idx],
		})
	}

	ctx.RespEntity(result)
}

// DeleteServiceTemplate TODO
// a service template can be delete only when it is not be used any more,
// which means that no process instance belongs to it.
func (ps *ProcServer) DeleteServiceTemplate(ctx *rest.Contexts) {
	input := new(metadata.DeleteServiceTemplatesInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := ps.CoreAPI.CoreService().Process().DeleteServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, input.ServiceTemplateID)
		if err != nil {
			blog.Errorf("delete service template: %d failed", input.ServiceTemplateID)
			return ctx.Kit.CCError.CCError(common.CCErrProcDeleteServiceTemplateFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// GetServiceTemplateSyncStatus check if service templates or modules with template need sync, return the status
func (ps *ProcServer) GetServiceTemplateSyncStatus(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil || bizID <= 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField))
		return
	}

	opt := new(metadata.GetServiceTemplateSyncStatusOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	const maxIDLen = 100
	if opt.IsPartial {
		if len(opt.ServiceTemplateIDs) == 0 {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "service_template_ids"))
			return
		}

		if len(opt.ServiceTemplateIDs) > maxIDLen {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "service_template_ids", maxIDLen))
			return
		}

		moduleCond := map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKServiceTemplateIDField: map[string]interface{}{
				common.BKDBIN: opt.ServiceTemplateIDs,
			},
		}

		statuses, _, err := ps.Logic.GetSvcTempSyncStatus(ctx.Kit, bizID, moduleCond, true)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}

		ctx.RespEntity(metadata.ServiceTemplateSyncStatus{ServiceTemplates: statuses})
		return
	} else {
		if len(opt.ModuleIDs) == 0 {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "bk_module_ids"))
			return
		}

		if len(opt.ModuleIDs) > maxIDLen {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "bk_module_ids", maxIDLen))
			return
		}

		moduleCond := map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: opt.ModuleIDs,
			},
			common.BKAppIDField: bizID,
			common.BKServiceTemplateIDField: map[string]interface{}{
				common.BKDBNE: common.ServiceTemplateIDNotSet,
			},
		}

		_, statuses, err := ps.Logic.GetSvcTempSyncStatus(ctx.Kit, bizID, moduleCond, false)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}

		ctx.RespEntity(metadata.ServiceTemplateSyncStatus{Modules: statuses})
		return
	}
}

// SearchRuleRelatedServiceTemplates search rule related service templates
func (ps *ProcServer) SearchRuleRelatedServiceTemplates(ctx *rest.Contexts) {
	requestBody := new(metadata.RuleRelatedServiceTemplateOption)
	if err := ctx.DecodeInto(requestBody); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if requestBody.ApplicationID == 0 {
		blog.Errorf("bk_biz_id should not be empty, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_biz_id"))
		return
	}

	if requestBody.QueryFilter == nil {
		blog.Errorf("search query_filter should not be empty, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "query_filter"))
		return
	}

	key, err := requestBody.QueryFilter.Validate(&querybuilder.RuleOption{NeedSameSliceElementType: true})
	if err != nil {
		blog.Errorf("search query_filter.%s validate failed, err: %v, rid: %s", key, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "query_filter."+key))
		return
	}

	templates, err := ps.Engine.CoreAPI.CoreService().HostApplyRule().SearchRuleRelatedServiceTemplates(ctx.Kit.Ctx,
		ctx.Kit.Header, requestBody)
	if err != nil {
		blog.Errorf("search rule related service templates failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(templates)
}

// UpdateServiceTemplateAttribute update service template attribute
func (ps *ProcServer) UpdateServiceTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.UpdateServTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := ps.Engine.CoreAPI.CoreService().Process().UpdateServiceTemplateAttribute(ctx.Kit.Ctx, ctx.Kit.Header,
			option); err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

// DeleteServiceTemplateAttribute delete service template attribute
func (ps *ProcServer) DeleteServiceTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.DeleteServTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := ps.Engine.CoreAPI.CoreService().Process().DeleteServiceTemplateAttribute(ctx.Kit.Ctx, ctx.Kit.Header,
			option); err != nil {
			return err
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

// ListServiceTemplateAttribute list service template attribute
func (ps *ProcServer) ListServiceTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.ListServTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	data, err := ps.Engine.CoreAPI.CoreService().Process().ListServiceTemplateAttribute(ctx.Kit.Ctx, ctx.Kit.Header,
		option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(data)
}
