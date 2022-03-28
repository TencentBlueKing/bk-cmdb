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
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *Service) CreateHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("CreateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.CreateHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	var rule metadata.HostApplyRule
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		rule, err = s.CoreAPI.CoreService().HostApplyRule().CreateHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID,
			option)
		if err != nil {
			blog.ErrorJSON("CreateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, "+
				"option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(rule)
}

func (s *Service) UpdateHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	ruleIDStr := ctx.Request.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, parse biz id failed, ruleIDStr: %s, err: %v,rid:%s", ruleIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField))
		return
	}

	option := metadata.UpdateHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	var rule metadata.HostApplyRule
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		rule, err = s.CoreAPI.CoreService().HostApplyRule().UpdateHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, ruleID, option)
		if err != nil {
			blog.ErrorJSON("UpdateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(rule)
}

// DeleteHostApplyRule delete the host automatic application rule in the module scenario.
func (s *Service) DeleteHostApplyRule(ctx *rest.Contexts) {
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

	if rawErr := option.ValidateModuleOption(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option); err != nil {
			blog.ErrorJSON("DeleteHostApplyRule failed, core service DeleteHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(make(map[string]interface{}))
}

func (s *Service) GetHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	ruleIDStr := ctx.Request.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, parse biz id failed, ruleIDStr: %s, err: %v,rid:%s", ruleIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField))
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().GetHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, ruleID)
	if err != nil {
		blog.ErrorJSON("GetHostApplyRule failed, core service GetHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, err.Error(), rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(rule)
}

func (s *Service) ListHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ListHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.ListHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if len(option.ModuleIDs) == 0 {
		blog.Errorf("ListHostApplyRule failed, parameter bk_module_ids empty, rid:%s", err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_module_ids"))
		return
	}

	ruleResult, err := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option)
	if err != nil {
		blog.ErrorJSON("ListHostApplyRule failed, core service ListHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(ruleResult)
}

func (s *Service) BatchCreateOrUpdateHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("BatchCreateOrUpdateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.BatchCreateOrUpdateApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	var batchResult metadata.BatchCreateOrUpdateHostApplyRuleResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		batchResult, err = s.CoreAPI.CoreService().HostApplyRule().BatchUpdateHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option)
		if err != nil {
			blog.ErrorJSON("BatchCreateOrUpdateHostApplyRule failed, coreservice BatchUpdateHostApplyRule failed, option: %s, result: %s, err: %s, rid:%s", option, batchResult, err, rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	var firstErr errors.CCErrorCoder
	for _, item := range batchResult.Items {
		if err := item.GetError(); err != nil {
			firstErr = err
			break
		}
	}
	if firstErr != nil {
		ctx.RespEntityWithError(batchResult, firstErr)
		return
	}
	ctx.RespEntity(batchResult)
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

func (s *Service) updateHostPlan(planResult metadata.HostApplyPlanResult, kit *rest.Kit) (
	[]metadata.HostApplyResult, errors.CCErrorCoder) {
	var (
		wg       sync.WaitGroup
		firstErr errors.CCErrorCoder
	)

	// update host instances, allow partial success
	updateMap := make(map[string][]int64, 0)
	for _, plan := range planResult.Plans {
		if len(plan.UpdateFields) == 0 {
			continue
		}
		dataStr := plan.GetUpdateDataStr()
		updateMap[dataStr] = append(updateMap[dataStr], plan.HostID)
	}

	hostApplyResults := make([]metadata.HostApplyResult, 0)
	pipeline := make(chan bool, 5)
	for dataStr, hostIDs := range updateMap {

		pipeline <- true
		wg.Add(1)

		go func(dataStr string, hostIDs []int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			mergeCond, data := generateCondition(dataStr, hostIDs)
			counts, cErr := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
				common.BKTableNameBaseHost, []map[string]interface{}{mergeCond})
			if cErr != nil {
				if firstErr == nil {
					firstErr = cErr
				}
				blog.Errorf("get hosts count failed, filter: %+v, err: %v, rid: %s", mergeCond, cErr, kit.Rid)
				return
			}
			if counts[0] == 0 {
				blog.V(5).Infof("no hosts founded, filter: %+v, rid: %s", mergeCond, kit.Rid)
				return
			}

			// If there is no eligible host, then return directly.
			updateOp := &metadata.UpdateOption{Data: data, Condition: mergeCond}

			_, err := s.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header,
				common.BKInnerObjIDHost, updateOp)
			if err != nil {
				blog.Errorf("update host failed, option: %s, err: %v, rid: %s", updateOp, err, kit.Rid)
				for _, hostID := range hostIDs {
					hostApplyResult := metadata.HostApplyResult{HostID: hostID}
					hostApplyResult.SetError(kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
					hostApplyResults = append(hostApplyResults, hostApplyResult)
				}
				if firstErr == nil {
					firstErr = errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
				}
				return
			}

			for _, hostID := range hostIDs {
				hostApplyResult := metadata.HostApplyResult{HostID: hostID}
				hostApplyResults = append(hostApplyResults, hostApplyResult)
			}

		}(dataStr, hostIDs)
	}

	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}

	return hostApplyResults, nil
}

// getHostIDByCondition get the final list of hostIDs.
func (s *Service) getHostIDByCondition(kit *rest.Kit, bizID int64, modIDs []int64, hostIDs []int64) ([]int64, error) {
	relationReq := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   modIDs,
		Page:          metadata.BasePage{Limit: common.BKNoLimit},
		Fields:        []string{common.BKModuleIDField, common.BKHostIDField},
	}

	// hostIDs are not empty in the invalid host scenario.
	if hostIDs != nil {
		relationReq.HostIDArr = hostIDs
	}
	hostRelations, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, relationReq)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, kit.Rid)
		return []int64{}, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if len(hostRelations.Info) == 0 {
		return []int64{}, nil
	}
	hostModuleMap := make(map[int64]struct{})
	for _, item := range hostRelations.Info {
		hostModuleMap[item.HostID] = struct{}{}
	}
	result := make([]int64, 0)
	for hostID := range hostModuleMap {
		result = append(result, hostID)
	}
	return result, nil
}

func (s *Service) generateApplyPlan(ctx *rest.Contexts, bizID int64, planRequest metadata.HostApplyPlanRequest) (
	metadata.HostApplyPlanResult, errors.CCErrorCoder) {

	rid := ctx.Kit.Rid

	relationReq := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   planRequest.ModuleIDs,
		Page:          metadata.BasePage{Limit: common.BKNoLimit},
		Fields:        []string{common.BKModuleIDField, common.BKHostIDField},
	}
	if planRequest.HostIDs != nil {
		relationReq.HostIDArr = planRequest.HostIDs
	}
	hostRelations, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header, relationReq)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, rid)
		return metadata.HostApplyPlanResult{}, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	hostModuleMap := make(map[int64][]int64)
	moduleIDs := make([]int64, 0)
	for _, item := range hostRelations.Info {
		hostModuleMap[item.HostID] = append(hostModuleMap[item.HostID], item.ModuleID)
		moduleIDs = append(moduleIDs, item.ModuleID)
	}
	hostModules := make([]metadata.Host2Modules, 0)
	for hostID, moduleIDs := range hostModuleMap {
		hostModules = append(hostModules, metadata.Host2Modules{
			HostID:    hostID,
			ModuleIDs: moduleIDs,
		})
	}

	ruleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: moduleIDs,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}
	rules, ccErr := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID,
		ruleOption)
	if ccErr != nil {
		blog.Errorf("list host apply rule failed, bizID: %d, opt: %#v, err: %v, rid: %s", bizID, ruleOption, ccErr, rid)
		return metadata.HostApplyPlanResult{}, ccErr
	}

	now := time.Now()
	if len(planRequest.AdditionalRules) > 0 {
	OuterLoop:
		for _, item := range planRequest.AdditionalRules {
			for index, rule := range rules.Info {
				if item.ModuleID == rule.ModuleID && item.AttributeID == rule.AttributeID {
					rules.Info[index].PropertyValue = item.PropertyValue
					continue OuterLoop
				}
			}
			rules.Info = append(rules.Info, metadata.HostApplyRule{BizID: bizID, ModuleID: item.ModuleID,
				AttributeID: item.AttributeID, PropertyValue: item.PropertyValue, Creator: ctx.Kit.User,
				Modifier: ctx.Kit.User, CreateTime: now, LastTime: now, SupplierAccount: ctx.Kit.SupplierAccount,
			})
		}
	}

	// filter out removed rules
	finalRules := make([]metadata.HostApplyRule, 0)
	for _, item := range rules.Info {
		if util.InArray(item.ID, planRequest.RemoveRuleIDs) || util.InArray(item.ID, planRequest.IgnoreRuleIDs) {
			continue
		}
		finalRules = append(finalRules, item)
	}

	planOption := metadata.HostApplyPlanOption{
		Rules:             finalRules,
		HostModules:       hostModules,
		ConflictResolvers: planRequest.ConflictResolvers,
	}

	planResult, ccErr := s.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(ctx.Kit.Ctx, ctx.Kit.Header, bizID,
		planOption)
	if ccErr != nil {
		blog.Errorf("generate apply plan failed, bizID: %d, opt: %#v, err: %v, rid: %s", bizID, planOption, ccErr, rid)
		return planResult, ccErr
	}
	planResult.Rules = rules.Info
	return planResult, nil
}

func (s *Service) getUpdateDataStr(kit *rest.Kit, rules []metadata.HostAttribute) (
	string, errors.CCErrorCoder) {
	attributeIDs := make([]int64, 0)
	for _, rule := range rules {
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

func (s *Service) updateHostAttributes(kit *rest.Kit, planResult []metadata.HostAttribute,
	hostIDs []int64) errors.CCErrorCoder {

	dataStr, err := s.getUpdateDataStr(kit, planResult)
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

// GetHostApplyTaskStatus get host auto-apply asynchronous task status.
func (s *Service) GetHostApplyTaskStatus(ctx *rest.Contexts) {

	syncStatusOpt := new(metadata.HostApplyTaskStatusOption)
	if err := ctx.DecodeInto(syncStatusOpt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if rawErr := syncStatusOpt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// get host auto-apply task status by task ids. Query the automatic application status of the host. Since the instID
	// when creating a task is a random number, the instID input condition is not required when querying.
	statusOpt := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKTaskTypeField: common.SyncModuleHostApplyTaskFlag,
			common.BKTaskIDField:   map[string]interface{}{common.BKDBIN: syncStatusOpt.TaskIDS},
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

// UpdateModuleHostApplyRule update host auto-apply rules in module dimension.
func (s *Service) UpdateModuleHostApplyRule(ctx *rest.Contexts) {

	syncOpt := new(metadata.HostApplyModulesOption)
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
	// the instID here can be assigned a random number. Random instID from 1 to 10000 in module scenario.
	randInstNum := util.RandInt64WithRange(int64(1), int64(10000))

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		taskRes, err := s.CoreAPI.TaskServer().Task().Create(ctx.Kit.Ctx, ctx.Kit.Header,
			common.SyncModuleHostApplyTaskFlag, randInstNum, []interface{}{syncOpt})
		if err != nil {
			blog.Errorf("create module host apply sync rule task failed, opt: %+v, err: %v, rid: %s",
				syncOpt, err, ctx.Kit.Rid)
			return err
		}
		taskInfo = taskRes
		blog.V(4).Infof("successfully created module host apply sync task: %#v, rid: %s", taskRes, ctx.Kit.Rid)
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(metadata.HostApplyTaskResult{BizID: taskInfo.InstID, TaskID: taskInfo.TaskID})
}

// ExecModuleHostApplyRule the host automatically applies rules in the asynchronous execution module scenario.
func (s *Service) ExecModuleHostApplyRule(ctx *rest.Contexts) {

	rid := ctx.Kit.Rid

	planReq := new(metadata.HostApplyModulesOption)
	if err := ctx.DecodeInto(planReq); err != nil {
		ctx.RespAutoError(err)
		return
	}
	hostIDs, err := s.getHostIDByCondition(ctx.Kit, planReq.BizID, planReq.ModuleIDs, planReq.HostIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// enable host apply on module
		op := &metadata.UpdateOption{
			Condition: map[string]interface{}{
				common.BKModuleIDField: map[string]interface{}{common.BKDBIN: planReq.ModuleIDs}},
			Data: map[string]interface{}{common.HostApplyEnabledField: true},
		}

		_, err := s.Engine.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKInnerObjIDModule, op)
		if err != nil {
			blog.Errorf("update instance of module failed, option: %s, err: %v, rid: %s", op, err, rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}

		rulesOption := make([]metadata.CreateOrUpdateApplyRuleOption, 0)
		for _, rule := range planReq.AdditionalRules {

			rulesOption = append(rulesOption, metadata.CreateOrUpdateApplyRuleOption{
				AttributeID:   rule.AttributeID,
				ModuleID:      rule.ModuleID,
				PropertyValue: rule.PropertyValue})
		}

		// 1、update or add rules.
		saveRuleOp := metadata.BatchCreateOrUpdateApplyRuleOption{Rules: rulesOption}
		if _, ccErr := s.CoreAPI.CoreService().HostApplyRule().BatchUpdateHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header,
			planReq.BizID, saveRuleOp); ccErr != nil {
			blog.Errorf("update host rule failed, bizID: %s, req: %s, err: %v, rid: %s", planReq.BizID, saveRuleOp,
				ccErr, rid)
			return ccErr
		}

		// 2、delete rules.
		if len(planReq.RemoveRuleIDs) > 0 {
			removeOp := metadata.DeleteHostApplyRuleOption{
				RuleIDs:   planReq.RemoveRuleIDs,
				ModuleIDs: planReq.ModuleIDs}
			if ccErr := s.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header,
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

	// the following three scenarios do not require the update of the host properties to be automatically applied:
	// 1. The changed flag is false. 2. This request only deletes the rule scenario. 3. No eligible host was found.
	if !planReq.Changed || len(planReq.AdditionalRules) == 0 || len(hostIDs) == 0 {
		ctx.RespEntity(nil)
		return
	}

	// update host operation is not done in a transaction, since the successfully updated hosts need not roll back
	ctx.Kit.Header.Del(common.TransactionIdHeader)

	attributes := make([]metadata.HostAttribute, 0)

	for _, rule := range planReq.AdditionalRules {
		attributes = append(attributes, metadata.HostAttribute{
			AttributeID:   rule.AttributeID,
			PropertyValue: rule.PropertyValue})
	}
	// apply module attribute rules to the host.
	err = s.updateHostAttributes(ctx.Kit, attributes, hostIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// ListHostRelatedApplyRule 返回主机关联的规则信息（仅返回启用模块的规则）
func (s *Service) ListHostRelatedApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ListHostRelatedApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.ListHostRelatedApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	result, ccErr := s.listHostRelatedApplyRule(ctx, bizID, option)
	if ccErr != nil {
		blog.Errorf("ListHostRelatedApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	ctx.RespEntity(result)
}

func (s *Service) listHostRelatedApplyRule(ctx *rest.Contexts, bizID int64,
	option metadata.ListHostRelatedApplyRuleOption) (map[int64][]metadata.HostApplyRule, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid

	relationOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     option.HostIDs,
		Page:          metadata.BasePage{Limit: common.BKNoLimit},
		Fields:        []string{common.BKModuleIDField, common.BKHostIDField},
	}
	relationResult, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header,
		relationOption)
	if err != nil {
		blog.Errorf("get host module relation failed, option: %+v, err: %v, rid: %s", relationOption, err, rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	hostModuleIDMap := make(map[int64][]int64)
	moduleIDs := make([]int64, 0)
	for _, item := range relationResult.Info {
		moduleIDs = append(moduleIDs, item.ModuleID)
		hostModuleIDMap[item.HostID] = append(hostModuleIDMap[item.HostID], item.ModuleID)
	}

	// filter enabled modules
	moduleFilter := &metadata.QueryCondition{
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Fields: []string{common.BKModuleIDField},
		Condition: map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: moduleIDs,
			},
			common.HostApplyEnabledField: true,
		},
	}
	moduleResult, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDModule, moduleFilter)
	if err != nil {
		blog.Errorf("get module failed, option: %#v, err: %v, rid: %s", moduleFilter, err, rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	validModuleIDs := make([]int64, 0)
	for _, item := range moduleResult.Info {
		moduleID, err := util.GetInt64ByInterface(item[common.BKModuleIDField])
		if err != nil {
			blog.Errorf("parse module id failed, err: %v, module: %#v, rid: %s", err, item, rid)
			return nil, ctx.Kit.CCError.CCError(common.CCErrCommParseDBFailed)
		}
		validModuleIDs = append(validModuleIDs, moduleID)
	}

	ruleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: validModuleIDs,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}
	ruleResult, ccErr := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header,
		bizID, ruleOption)
	if ccErr != nil {
		blog.Errorf("list host apply rule failed, bizID: %d, option: %#v, err: %v, rid: %s", bizID, option, ccErr, rid)
		return nil, ccErr
	}
	// moduleID -> []hostApplyRule
	moduleRules := make(map[int64][]metadata.HostApplyRule)
	for _, item := range ruleResult.Info {
		moduleRules[item.ModuleID] = append(moduleRules[item.ModuleID], item)
	}

	// hostID -> []moduleIDs
	result := make(map[int64][]metadata.HostApplyRule)
	for _, hostID := range option.HostIDs {
		moduleIDs, exist := hostModuleIDMap[hostID]
		if !exist {
			continue
		}
		for _, moduleID := range moduleIDs {
			rules, exist := moduleRules[moduleID]
			if exist {
				result[hostID] = append(result[hostID], rules...)
			}
		}
	}
	return result, nil
}
