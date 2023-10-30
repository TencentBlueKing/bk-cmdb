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

	acMeta "configcenter/src/ac/meta"
	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

// FindModuleHostRelation find host with module by module id
func (s *Service) FindModuleHostRelation(ctx *rest.Contexts) {
	req := ctx.Request
	defErr := ctx.Kit.CCError

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, err: %v, bizID: %s, rid: %s", err,
			req.PathParameter("bk_biz_id"), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}
	if bizID == 0 {
		ctx.RespAutoError(defErr.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	// authorize
	if resp, authorized := s.authHostUnderBiz(ctx.Kit, bizID); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	body := new(meta.FindModuleHostRelationParameter)
	if err := ctx.DecodeInto(body); err != nil {
		ctx.RespAutoError(err)
		return
	}
	rawErr := body.Validate()
	if rawErr.ErrCode != 0 {
		blog.ErrorJSON("validate request body err: %s, body: %s, rid: %s", rawErr.ToCCError(defErr).Error(), *body,
			ctx.Kit.Rid)
		ctx.RespAutoError(rawErr.ToCCError(defErr))
		return
	}

	// get host info
	distinctHostCond := &meta.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{bizID},
		ModuleIDArr:      body.ModuleIDS,
	}
	hostFields := append(body.HostFields, common.BKHostIDField)
	searchHostCond := &meta.QueryCondition{
		Fields: hostFields,
		Page:   body.Page,
	}
	hostRes, err := s.findDistinctHostInfo(ctx, distinctHostCond, searchHostCond)
	if err != nil {
		blog.Errorf("find distinct host info failed, err: %v, distinctHostCond: %+v, searchHostCond: %+v, rid: %s",
			err, *distinctHostCond, *searchHostCond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	hostLen := len(hostRes.Data.Info)
	if hostLen == 0 {
		ctx.RespEntity(meta.FindModuleHostRelationResult{
			Count:    hostRes.Data.Count,
			Relation: []meta.ModuleHostRelation{},
		})
		return
	}
	hostIDArr := make([]int64, hostLen)
	for index, host := range hostRes.Data.Info {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.ErrorJSON("host id not integer, host: %s, rid: %s", host, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		hostIDArr[index] = hostID
	}

	// get module info
	hostModuleConfig, err := s.Logic.GetHostRelations(ctx.Kit, meta.HostModuleRelationRequest{HostIDArr: hostIDArr,
		Fields: []string{common.BKModuleIDField, common.BKHostIDField}})
	if err != nil {
		blog.Errorf("GetConfigByCond failed, err: %v, hostIDArr: %v, rid: %s", err, hostIDArr, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	hostModuleMap := make(map[int64][]int64, hostLen)
	moduleIDArr := make([]int64, 0)
	for _, relation := range hostModuleConfig {
		hostModuleMap[relation.HostID] = append(hostModuleMap[relation.HostID], relation.ModuleID)
		moduleIDArr = append(moduleIDArr, relation.ModuleID)
	}

	moduleFields := append(body.ModuleFields, common.BKModuleIDField)
	moduleInfoMap, err := s.Logic.GetModuleMapByCond(ctx.Kit, moduleFields, map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{common.BKDBIN: moduleIDArr},
	})
	if err != nil {
		blog.Errorf("get module map failed, err: %v, moduleIDArr: %v, rid: %s", err, moduleIDArr, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// assemble host and module info
	relation := make([]meta.ModuleHostRelation, hostLen)
	for index, host := range hostRes.Data.Info {
		hostID, _ := util.GetInt64ByInterface(host[common.BKHostIDField])
		moduleIDs := hostModuleMap[hostID]
		modules := make([]map[string]interface{}, len(moduleIDs))
		for index, moduleID := range moduleIDs {
			modules[index] = moduleInfoMap[moduleID]
		}
		relation[index] = meta.ModuleHostRelation{
			Host:    host,
			Modules: modules,
		}
	}

	ctx.RespEntity(meta.FindModuleHostRelationResult{
		Count:    hostRes.Data.Count,
		Relation: relation,
	})
}

// FindHostsByServiceTemplates find hosts by service templates
func (s *Service) FindHostsByServiceTemplates(ctx *rest.Contexts) {
	defErr := ctx.Kit.CCError

	option := new(meta.FindHostsBySrvTplOpt)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := option.Validate()
	if rawErr.ErrCode != 0 {
		blog.Errorf("validate failed, err: %v, option: %#v, rid: %s", rawErr, *option, ctx.Kit.Rid)
		ctx.RespAutoError(rawErr.ToCCError(defErr))
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		ccErr := defErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
		ctx.RespAutoError(ccErr)
		return
	}
	if bizID == 0 {
		ccErr := defErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
		ctx.RespAutoError(ccErr)
		return
	}

	// authorize
	if resp, authorized := s.authHostUnderBiz(ctx.Kit, bizID); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	moduleCond := []meta.ConditionItem{
		{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    bizID,
		},
		{
			Field:    common.BKServiceTemplateIDField,
			Operator: common.BKDBIN,
			Value:    option.ServiceTemplateIDs,
		},
	}
	if len(option.ModuleIDs) > 0 {
		moduleCond = append(moduleCond, meta.ConditionItem{
			Field:    common.BKModuleIDField,
			Operator: common.BKDBIN,
			Value:    option.ModuleIDs,
		})
	}
	moduleIDArr, err := s.Logic.GetModuleIDByCond(ctx.Kit, meta.ConditionWithTime{Condition: moduleCond})
	if err != nil {
		blog.Errorf("get module id failed, err: %v, cond:%#v, rid: %s", err, moduleCond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(moduleIDArr) == 0 {
		ctx.RespEntity(meta.SearchHost{})
		return
	}

	distinctHostCond := &meta.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{bizID},
		ModuleIDArr:      moduleIDArr,
	}
	searchHostCond := &meta.QueryCondition{
		Fields:         option.Fields,
		Page:           option.Page,
		DisableCounter: true,
	}

	result, err := s.findDistinctHostInfo(ctx, distinctHostCond, searchHostCond)
	if err != nil {
		blog.Errorf("find distinct host info failed, err: %v, distinctHostCond: %#v, searchHostCond: %#v, rid: %s",
			err, *distinctHostCond, *searchHostCond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result.Data)
}

// findDistinctHostInfo find distinct host info
func (s *Service) findDistinctHostInfo(ctx *rest.Contexts, distinctHostCond *meta.DistinctHostIDByTopoRelationRequest,
	searchHostCond *meta.QueryCondition) (*meta.SearchHostResult, error) {

	allHostIDs, err := s.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(ctx.Kit.Ctx, ctx.Kit.Header,
		distinctHostCond)
	if err != nil {
		blog.Errorf("get hostIDs failed, err: %v, input: %#v, rid: %s", err, distinctHostCond, ctx.Kit.Rid)
		return nil, err
	}
	sort.Sort(util.Int64Slice(allHostIDs))

	// get hostIDs according from page info
	hostCnt := len(allHostIDs)
	startIndex := searchHostCond.Page.Start
	if startIndex >= hostCnt {
		return &meta.SearchHostResult{
			BaseResp: meta.SuccessBaseResp,
			Data: &meta.SearchHost{
				Count: hostCnt,
			},
		}, nil
	}
	endindex := startIndex + searchHostCond.Page.Limit
	if endindex > hostCnt {
		endindex = hostCnt
	}
	hostIDs := allHostIDs[startIndex:endindex]

	cond := meta.QueryCondition{
		Fields: searchHostCond.Fields,
		Page: meta.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr{
			common.BKHostIDField: mapstr.MapStr{
				common.BKDBIN: hostIDs,
			},
		},
	}
	hostInfo, err := s.Logic.SearchHostInfo(ctx.Kit, cond)
	if err != nil {
		blog.Errorf("findDistinctHostInfo failed, SearchHostInfo error: %v, input:%#v, rid: %s", err, cond, ctx.Kit.Rid)
		return nil, err
	}

	return &meta.SearchHostResult{
		BaseResp: meta.SuccessBaseResp,
		Data: &meta.SearchHost{
			Count: hostCnt,
			Info:  hostInfo,
		},
	}, nil
}

// FindHostsBySetTemplates find hosts by set templates
func (s *Service) FindHostsBySetTemplates(ctx *rest.Contexts) {

	defErr := ctx.Kit.CCError

	option := new(meta.FindHostsBySetTplOpt)
	if err := ctx.DecodeInto(option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	rawErr := option.Validate()
	if rawErr.ErrCode != 0 {
		blog.Errorf("validate failed, err: %v, option: %#v, rid: %s", rawErr, *option, ctx.Kit.Rid)
		ctx.RespAutoError(rawErr.ToCCError(defErr))
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		ccErr := defErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
		ctx.RespAutoError(ccErr)
		return
	}
	if bizID == 0 {
		ccErr := defErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
		ctx.RespAutoError(ccErr)
		return
	}

	// authorize
	if resp, authorized := s.authHostUnderBiz(ctx.Kit, bizID); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	setCond := []meta.ConditionItem{
		{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    bizID,
		},
		{
			Field:    common.BKSetTemplateIDField,
			Operator: common.BKDBIN,
			Value:    option.SetTemplateIDs,
		},
	}
	if len(option.SetIDs) > 0 {
		setCond = append(setCond, meta.ConditionItem{
			Field:    common.BKSetIDField,
			Operator: common.BKDBIN,
			Value:    option.SetIDs,
		})
	}

	setIDArr, err := s.Logic.GetSetIDByCond(ctx.Kit, meta.ConditionWithTime{Condition: setCond})
	if err != nil {
		blog.Errorf("get set id by cond(%#v) failed, err: %v, rid: %s", setCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(setIDArr) == 0 {
		ctx.RespEntity(meta.SearchHost{})
		return
	}

	distinctHostCond := &meta.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{bizID},
		SetIDArr:         setIDArr,
	}
	searchHostCond := &meta.QueryCondition{
		Fields: option.Fields,
		Page:   option.Page,
	}

	result, err := s.findDistinctHostInfo(ctx, distinctHostCond, searchHostCond)
	if err != nil {
		blog.Errorf("find distinct host info failed, err: %v, distinctHostCond: %#v, searchHostCond: %#v, rid: %s",
			err, *distinctHostCond, *searchHostCond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result.Data)
}

// FindHostsByTopo find hosts by topo node except for biz
func (s *Service) FindHostsByTopo(ctx *rest.Contexts) {

	option := new(meta.FindHostsByTopoOpt)
	if err := ctx.DecodeInto(option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := option.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("find hosts by topo failed, validate err: %v, option: %#v, rid: %s", rawErr, *option, ctx.Kit.Rid)
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil || bizID <= 0 {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	// authorize
	if resp, authorized := s.authHostUnderBiz(ctx.Kit, bizID); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	// generate search condition,
	// if node is not a set or a module, we need to traverse its child topo to the set level to get hosts by relation
	distinctHostCond := &meta.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{bizID},
	}

	if option.ObjID == common.BKInnerObjIDSet {
		distinctHostCond.SetIDArr = []int64{option.InstID}
	} else if option.ObjID == common.BKInnerObjIDModule {
		distinctHostCond.ModuleIDArr = []int64{option.InstID}
	} else {
		setIDArr, err := s.Logic.GetSetIDsByTopo(ctx.Kit, option.ObjID, []int64{option.InstID})
		if err != nil {
			blog.Errorf("find hosts by topo failed, get set ID by topo err: %v, objID: %s, instID: %d, rid: %s",
				err, option.ObjID, option.InstID, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		if len(setIDArr) == 0 {
			ctx.RespEntity(meta.SearchHost{})
			return
		}

		distinctHostCond.SetIDArr = setIDArr
	}

	searchHostCond := &meta.QueryCondition{
		Fields: option.Fields,
		Page:   option.Page,
	}

	result, err := s.findDistinctHostInfo(ctx, distinctHostCond, searchHostCond)
	if err != nil {
		blog.Errorf("find hosts by topo failed, cond: %#v, search cond: %#v, er: %v, rid: %s", err, *distinctHostCond,
			*searchHostCond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result.Data)
}

// ListResourcePoolHosts list hosts of resource pool
func (s *Service) ListResourcePoolHosts(ctx *rest.Contexts) {
	header := ctx.Request.Request.Header
	rid := ctx.Kit.Rid
	defErr := ctx.Kit.CCError

	parameter := meta.ListHostsParameter{}
	if err := ctx.DecodeInto(&parameter); nil != err {
		ctx.RespAutoError(err)
		return
	}

	filter := &meta.QueryCondition{
		Fields: []string{common.BKAppIDField, common.BKAppNameField},
		Page: meta.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKDefaultField: common.DefaultAppFlag,
		},
	}
	appResult, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, header, common.BKInnerObjIDApp,
		filter)
	if err != nil {
		blog.Errorf("get default biz failed, filter: %+v, err: %v, rid:%s", filter, err, rid)
		ccErr := defErr.Error(common.CCErrCommHTTPDoRequestFailed)
		ctx.RespAutoError(ccErr)
		return
	}

	if appResult.Count == 0 {
		blog.Errorf("ListResourcePoolHosts failed, get default app failed, not found, rid: %s", rid)
		ccErr := defErr.Error(common.CCErrCommBizNotFoundError)
		ctx.RespAutoError(ccErr)
		return
	}

	// only use biz with same supplier account if query returns multiple biz
	bizData := appResult.Info[0]
	bizCount := 0
	for _, biz := range appResult.Info {
		supplier, _ := biz.String(common.BkSupplierAccount)
		if supplier == util.GetOwnerID(header) {
			bizCount++
			bizData = biz
		}
	}
	if bizCount > 1 {
		blog.Errorf("ListResourcePoolHosts failed, get multiple default app, result: %+v, rid: %s", appResult, rid)
		ccErr := defErr.Error(common.CCErrCommGetMultipleObject)
		ctx.RespAutoError(ccErr)
		return
	}

	// get biz ID
	bizID, err := util.GetInt64ByInterface(bizData[common.BKAppIDField])
	if err != nil {
		blog.Errorf("parse app data failed, biz: %s, err: %v, rid: %s", bizData, err, rid)
		ccErr := defErr.Error(common.CCErrCommParseDataFailed)
		ctx.RespAutoError(ccErr)
		return
	}

	// do host search
	hostResult, ccErr := s.listBizHosts(ctx, bizID, parameter)
	if ccErr != nil {
		blog.Errorf("listBizHosts failed, bizID: %s, parameter: %s, err: %v, rid:%s", bizID, parameter, ccErr, rid)
		ctx.RespAutoError(ccErr)
		return
	}
	ctx.RespEntity(hostResult)
}

// ListBizHosts list host under business specified by path parameter
func (s *Service) ListBizHosts(ctx *rest.Contexts) {

	rid := ctx.Kit.Rid
	defErr := ctx.Kit.CCError
	req := ctx.Request
	parameter := meta.ListHostsParameter{}
	if err := ctx.DecodeInto(&parameter); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if key, err := parameter.Validate(); err != nil {
		blog.ErrorJSON("ListBizHosts failed, Validate failed,parameter:%s, err: %s, rid:%s", parameter, err,
			ctx.Kit.Rid)
		ccErr := defErr.CCErrorf(common.CCErrCommParamsInvalid, key)
		ctx.RespAutoError(ccErr)
		return
	}
	bizID, err := strconv.ParseInt(req.PathParameter("appid"), 10, 64)
	if err != nil {
		ccErr := defErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
		ctx.RespAutoError(ccErr)
		return
	}
	if bizID == 0 {
		ccErr := defErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
		ctx.RespAutoError(ccErr)
		return
	}

	// authorize
	if resp, authorized := s.authHostUnderBiz(ctx.Kit, bizID); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	hostResult, ccErr := s.listBizHosts(ctx, bizID, parameter)
	if ccErr != nil {
		blog.ErrorJSON("ListBizHosts failed, listBizHosts failed, bizID: %s, parameter: %s, err: %s, rid:%s", bizID,
			parameter, ccErr.Error(), rid)
		ctx.RespAutoError(ccErr)
		return
	}
	ctx.RespEntity(hostResult)
}

func (s *Service) listBizHosts(ctx *rest.Contexts, bizID int64, parameter meta.ListHostsParameter) (
	result *meta.ListHostResult, ccErr errors.CCErrorCoder) {
	header := ctx.Kit.Header
	rid := ctx.Kit.Rid
	defErr := ctx.Kit.CCError

	if parameter.Page.IsIllegal() {
		blog.Errorf("ListBizHosts failed, page limit %d illegal, rid:%s", parameter.Page.Limit, ctx.Kit.Rid)
		return result, defErr.CCErrorf(common.CCErrCommParamsInvalid, "page.limit")
	}

	if len(parameter.SetIDs) != 0 && len(parameter.SetCond) != 0 {
		blog.Errorf("ListBizHosts failed, bk_set_ids and set_cond can't both be set, rid:%s", ctx.Kit.Rid)
		return result, defErr.CCErrorf(common.CCErrCommParamsInvalid, "bk_set_ids and set_cond can't both be set")
	}

	if len(parameter.ModuleIDs) != 0 && len(parameter.ModuleCond) != 0 {
		blog.Errorf("list biz hosts failed, bk_module_ids and module_cond can't both be set, rid: %s", ctx.Kit.Rid)
		return result, defErr.CCErrorf(common.CCErrCommParamsInvalid, "bk_module_ids and module_cond can't both be set")
	}

	setIDList := make([]int64, 0)
	if len(parameter.SetCond) != 0 {
		// set the app id condition
		parameter.SetCond = append(parameter.SetCond, meta.ConditionItem{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    bizID,
		})

		setIDList, ccErr = s.Logic.GetInstIDs(ctx.Kit, common.BKInnerObjIDSet, parameter.SetCond)
		if ccErr != nil {
			return nil, ccErr
		}

		if len(setIDList) == 0 {
			return &meta.ListHostResult{Count: 0, Info: []map[string]interface{}{}}, nil
		}
	}

	if len(parameter.SetIDs) != 0 {
		setIDList = parameter.SetIDs
	}

	moduleIDList := make([]int64, 0)
	if len(parameter.ModuleCond) != 0 {
		// set the app id condition
		parameter.ModuleCond = append(parameter.ModuleCond, meta.ConditionItem{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    bizID,
		})

		moduleIDList, ccErr = s.Logic.GetInstIDs(ctx.Kit, common.BKInnerObjIDModule, parameter.ModuleCond)
		if ccErr != nil {
			return nil, ccErr
		}

		if len(moduleIDList) == 0 {
			return &meta.ListHostResult{Count: 0, Info: []map[string]interface{}{}}, nil
		}
	}

	if len(parameter.ModuleIDs) != 0 {
		moduleIDList = parameter.ModuleIDs
	}

	option := &meta.ListHosts{
		BizID:              bizID,
		SetIDs:             setIDList,
		ModuleIDs:          moduleIDList,
		HostPropertyFilter: parameter.HostPropertyFilter,
		Fields:             parameter.Fields,
		Page:               parameter.Page,
	}
	hostResult, err := s.CoreAPI.CoreService().Host().ListHosts(ctx.Kit.Ctx, header, option)
	if err != nil {
		blog.Errorf("find host failed, err: %s, input:%#v, rid:%s", err.Error(), parameter, rid)
		return result, defErr.CCError(common.CCErrHostGetFail)
	}
	return hostResult, nil
}

// ListHostsWithNoBiz list host for no biz case merely
func (s *Service) ListHostsWithNoBiz(ctx *rest.Contexts) {
	header := ctx.Kit.Header
	rid := ctx.Kit.Rid
	defErr := ctx.Kit.CCError

	parameter := &meta.ListHostsWithNoBizParameter{}
	if err := ctx.DecodeInto(&parameter); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if key, err := parameter.Validate(); err != nil {
		blog.ErrorJSON("ListHostsWithNoBiz failed, decode body failed,parameter:%s, err: %#v, rid:%s", parameter,
			err, ctx.Kit.Rid)
		ccErr := defErr.CCErrorf(common.CCErrCommParamsInvalid, key)
		ctx.RespAutoError(ccErr)
		return
	}

	parameter.Page.Sort = common.BKHostIDField
	option := &meta.ListHosts{
		HostPropertyFilter: parameter.HostPropertyFilter,
		Fields:             parameter.Fields,
		Page:               parameter.Page,
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	host, err := s.CoreAPI.CoreService().Host().ListHosts(ctx.Kit.Ctx, header, option)
	if err != nil {
		blog.Errorf("find host failed, err: %s, input:%#v, rid:%s", err.Error(), parameter, rid)
		ctx.RespAutoError(defErr.Error(common.CCErrHostGetFail))
		return
	}
	ctx.RespEntity(host)

}

// ListBizHostsTopo list hosts under business specified by path parameter with their topology information
func (s *Service) ListBizHostsTopo(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	if bizID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	// authorize
	if resp, authorized := s.authHostUnderBiz(ctx.Kit, bizID); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	parameter := &meta.ListBizHostsTopoParameter{}
	if err := ctx.DecodeInto(&parameter); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := parameter.Validate(ctx.Kit.CCError); err != nil {
		blog.ErrorJSON("list biz host topo but input %s is invalid, err: %s, rid: %s", parameter, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)

	// if set filter or module filter is set, search them first to get ids to filter hosts
	filteredSetIDs, setMap, err := s.parseHostsTopoFilter(ctx.Kit, common.BKInnerObjIDSet, bizID,
		parameter.SetPropertyFilter)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if parameter.SetPropertyFilter != nil && len(filteredSetIDs) == 0 {
		ctx.RespEntityWithCount(0, make([]meta.HostTopo, 0))
		return
	}

	filteredModuleIDs, moduleMap, err := s.parseHostsTopoFilter(ctx.Kit, common.BKInnerObjIDModule, bizID,
		parameter.ModulePropertyFilter)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if parameter.ModulePropertyFilter != nil && len(filteredModuleIDs) == 0 {
		ctx.RespEntityWithCount(0, make([]meta.HostTopo, 0))
		return
	}

	// search all hosts
	option := &meta.ListHosts{
		BizID:              bizID,
		SetIDs:             filteredSetIDs,
		ModuleIDs:          filteredModuleIDs,
		HostPropertyFilter: parameter.HostPropertyFilter,
		Fields:             append(parameter.Fields, common.BKHostIDField),
		Page:               parameter.Page,
	}
	hosts, err := s.CoreAPI.CoreService().Host().ListHosts(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		blog.Errorf("find host failed, err: %s, input:%#v, rid: %s", err.Error(), parameter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrHostGetFail))
		return
	}

	if len(hosts.Info) == 0 {
		ctx.RespEntity(hosts)
		return
	}

	hostTopos, err := s.rearrangeBizHostTopo(ctx.Kit, hosts, bizID, setMap, moduleMap)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(hostTopos)
}

func (s *Service) rearrangeBizHostTopo(kit *rest.Kit, hosts *meta.ListHostResult, bizID int64,
	setMap, moduleMap map[int64]string) (*meta.HostTopoResult, error) {

	// search all hosts' host module relations
	hostIDs := make([]int64, 0)
	for _, host := range hosts.Info {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.ErrorJSON("host: %s bk_host_id field invalid, rid: %s", host, kit.Rid)
			return nil, err
		}
		hostIDs = append(hostIDs, hostID)
	}

	relationCond := meta.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     hostIDs,
		Fields:        []string{common.BKSetIDField, common.BKModuleIDField, common.BKHostIDField},
	}
	relations, err := s.Logic.GetHostRelations(kit, relationCond)
	if err != nil {
		blog.ErrorJSON("read host module relation error: %s, input: %s, rid: %s", err, hosts, kit.Rid)
		return nil, err
	}

	// generate host to set and module relation map
	setIDs := make([]int64, 0)
	moduleIDs := make([]int64, 0)
	relation := make(map[int64]map[int64][]int64)
	for _, r := range relations {
		setIDs = append(setIDs, r.SetID)
		moduleIDs = append(moduleIDs, r.ModuleID)
		setModule, ok := relation[r.HostID]
		if !ok {
			setModule = make(map[int64][]int64)
		}
		setModule[r.SetID] = append(setModule[r.SetID], r.ModuleID)
		relation[r.HostID] = setModule
	}

	// search all module and set info that is not already searched before
	setMap, err = s.getOtherInstInfo(kit, common.BKInnerObjIDSet, setIDs, setMap)
	if err != nil {
		return nil, err
	}

	moduleMap, err = s.getOtherInstInfo(kit, common.BKInnerObjIDModule, moduleIDs, moduleMap)
	if err != nil {
		return nil, err
	}

	// format the output
	hostTopos := &meta.HostTopoResult{
		Count: hosts.Count,
	}
	for _, host := range hosts.Info {
		hostTopo := meta.HostTopo{
			Host: host,
		}
		topos := make([]meta.Topo, 0)
		hostID, _ := util.GetInt64ByInterface(host[common.BKHostIDField])
		if setModule, ok := relation[hostID]; ok {
			for setID, moduleIDs := range setModule {
				topo := meta.Topo{
					SetID:   setID,
					SetName: setMap[setID],
				}
				modules := make([]meta.Module, 0)
				for _, moduleID := range moduleIDs {
					module := meta.Module{
						ModuleID:   moduleID,
						ModuleName: moduleMap[moduleID],
					}
					modules = append(modules, module)
				}
				topo.Module = modules
				topos = append(topos, topo)
			}
		}
		hostTopo.Topo = topos
		hostTopos.Info = append(hostTopos.Info, hostTopo)
	}

	return hostTopos, nil
}

func (s *Service) parseHostsTopoFilter(kit *rest.Kit, objID string, bizID int64, filter *querybuilder.QueryFilter) (
	[]int64, map[int64]string, error) {

	if filter == nil {
		return make([]int64, 0), make(map[int64]string), nil
	}

	cond, key, err := filter.ToMgo()
	if err != nil {
		blog.ErrorJSON("%s filter %s is invalid, err: %s, rid: %s", objID, filter, err, kit.Rid)
		return nil, nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid,
			fmt.Sprintf("%s_property_filter.%s", objID, key))
	}
	cond = mapstr.MapStr{common.BKDBAND: []mapstr.MapStr{{common.BKAppIDField: bizID}, cond}}

	instMap, instIDs, err := s.Logic.GetInstIDNameInfo(kit, objID, cond)
	if err != nil {
		blog.ErrorJSON("get %s by filter(%s) failed, err: %s, rid: %s", cond, err, kit.Rid)
		return nil, nil, err
	}

	return instIDs, instMap, nil
}

func (s *Service) getOtherInstInfo(kit *rest.Kit, objID string, ids []int64, instMap map[int64]string) (
	map[int64]string, error) {

	otherIDs := make([]int64, 0)
	if len(instMap) == 0 {
		otherIDs = ids
	} else {
		for _, id := range ids {
			if _, exists := instMap[id]; !exists {
				otherIDs = append(otherIDs, id)
			}
		}
	}

	if len(otherIDs) == 0 {
		return instMap, nil
	}

	otherIDs = util.IntArrayUnique(otherIDs)
	filter := map[string]interface{}{meta.GetInstIDFieldByObjID(objID): map[string]interface{}{common.BKDBIN: otherIDs}}
	otherMap, _, err := s.Logic.GetInstIDNameInfo(kit, objID, filter)
	if err != nil {
		blog.ErrorJSON("get %s by filter(%s) failed, err: %s, rid: %s", objID, filter, err, kit.Rid)
		return instMap, err
	}

	for key, value := range otherMap {
		instMap[key] = value
	}

	return instMap, nil
}

// ListHostDetailAndTopology obtain host details and corresponding topological relationships.
func (s *Service) ListHostDetailAndTopology(ctx *rest.Contexts) {
	header := ctx.Kit.Header
	rid := ctx.Kit.Rid
	defErr := ctx.Kit.CCError

	options := new(meta.ListHostsDetailAndTopoOption)
	if err := ctx.DecodeInto(&options); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := options.Validate(); rawErr != nil {
		blog.ErrorJSON("list host detail and topology, but validate option failed, option: %s, err: %s, rid:%s",
			options, rawErr, ctx.Kit.Rid)
		ctx.RespAutoError(defErr.CCErrorf(common.CCErrCommParamsInvalid, rawErr.Args))
		return
	}

	// read data from secondary mongodb nodes
	ctx.SetReadPreference(common.SecondaryPreferredMode)

	// search all hosts
	option := &meta.ListHosts{
		HostPropertyFilter: options.HostPropertyFilter,
		Fields:             append(options.Fields, common.BKHostIDField),
		Page:               options.Page,
	}
	hosts, err := s.CoreAPI.CoreService().Host().ListHosts(ctx.Kit.Ctx, header, option)
	if err != nil {
		blog.Errorf("find host failed, err: %s, input:%#v, rid:%s", err.Error(), options, rid)
		ctx.RespAutoError(defErr.Error(common.CCErrHostGetFail))
		return
	}

	if len(hosts.Info) == 0 {
		ctx.RespEntityWithCount(int64(hosts.Count), make([]*meta.HostDetailWithTopo, 0))
		return
	}

	hostTopo, bizList, err := s.Logic.ArrangeHostDetailAndTopology(ctx.Kit, options.WithBiz, hosts.Info)
	if err != nil {
		blog.Errorf("arrange host detail and topology failed, err: %v, rid :%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if options.WithBiz && s.AuthManager.Enabled() {
		if err = s.AuthManager.AuthorizeByInstanceID(ctx.Kit.Ctx, ctx.Kit.Header, acMeta.ViewBusinessResource,
			common.BKInnerObjIDApp, bizList...); err != nil {
			blog.Errorf("authorize failed, bizID: %v, err: %v, rid: %s", bizList, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}
	ctx.RespEntityWithCount(int64(hosts.Count), hostTopo)
	return
}

// CountTopoNodeHosts TODO
func (s *Service) CountTopoNodeHosts(ctx *rest.Contexts) {

	rid := ctx.Kit.Rid
	defErr := ctx.Kit.CCError

	option := meta.CountTopoNodeHostsOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		ctx.RespAutoError(defErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}
	if bizID == 0 {
		ctx.RespAutoError(defErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}
	topoNodeHostCounts, ccErr := s.countTopoNodeHosts(ctx, bizID, option)
	if ccErr != nil {
		blog.ErrorJSON("CountTopoNodeHosts failed, countTopoNodeHosts failed, option: %s, err: %s, rid:%s", option,
			ccErr.Error(), rid)
		ctx.RespAutoError(defErr.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	ctx.RespEntity(topoNodeHostCounts)
}

func (s *Service) countTopoNodeHosts(ctx *rest.Contexts, bizID int64,
	option meta.CountTopoNodeHostsOption) ([]meta.TopoNodeHostCount, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid
	topoRoot, ccErr := s.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(ctx.Kit.Ctx, ctx.Kit.Header, bizID,
		false)
	if ccErr != nil {
		blog.Errorf("search mainline instance topo failed, bizID: %d, err: %v, rid: %s", bizID, ccErr, rid)
		return nil, ccErr
	}
	moduleIDs := make([]int64, 0)
	nodeModuleIDMap := make(map[string]map[int64]bool)
	for _, topoNode := range option.Nodes {
		nodeModuleIDMap[topoNode.String()] = make(map[int64]bool)
		nodes := topoRoot.TraversalFindNode(topoNode.ObjectID, topoNode.InstanceID)
		for _, item := range nodes {
			item.DeepFirstTraversal(func(node *meta.TopoInstanceNode) {
				if node.ObjectID == common.BKInnerObjIDModule {
					moduleIDs = append(moduleIDs, node.InstanceID)
					nodeModuleIDMap[topoNode.String()][node.InstanceID] = true
				}
			})
		}
	}
	relationOption := meta.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   moduleIDs,
		Page: meta.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKModuleIDField, common.BKHostIDField},
	}
	relationResult, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header,
		&relationOption)
	if err != nil {
		blog.Errorf("get host module relation failed, option: %+v, err: %v, rid: %s", relationOption, err, rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	hostCounts := make([]meta.TopoNodeHostCount, 0)
	for _, topoNode := range option.Nodes {
		hostCount := meta.TopoNodeHostCount{
			Node:      topoNode,
			HostCount: 0,
		}
		moduleIDMap, ok := nodeModuleIDMap[topoNode.String()]
		if ok == false {
			hostCounts = append(hostCounts, hostCount)
			continue
		}
		hostIDs := make([]int64, 0)
		for _, item := range relationResult.Info {
			if _, ok := moduleIDMap[item.ModuleID]; ok == true {
				hostIDs = append(hostIDs, item.HostID)
			}
		}
		hostCount.HostCount = len(util.IntArrayUnique(hostIDs))
		hostCounts = append(hostCounts, hostCount)
	}

	return hostCounts, nil
}

// ListServiceTemplateIDsByHost list service template id one by one about hostID
func (s *Service) ListServiceTemplateIDsByHost(ctx *rest.Contexts) {

	input := struct {
		ID []int64 `json:"bk_host_id"`
	}{}
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(input.ID) > 200 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommValExceedMaxFailed, common.BKHostIDField, 200))
		return
	}

	rsp, err := s.Logic.ListServiceTemplateHostIDMap(ctx.Kit, input.ID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(rsp)
}

// ListHostTotalMainlineTopo list host total mainline topo tree
func (s *Service) ListHostTotalMainlineTopo(ctx *rest.Contexts) {

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	if bizID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	// authorize
	if resp, authorized := s.authHostUnderBiz(ctx.Kit, bizID); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	params := meta.FindHostTotalTopo{}
	if err := ctx.DecodeInto(&params); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := params.Validate(ctx.Kit.CCError); err != nil {
		blog.Errorf("validate param failed, param: %v, err: %v, rid: %s", params, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	rsp, err := s.Logic.ListHostTotalMainlineTopo(ctx.Kit, bizID, params)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(int64(len(rsp)), rsp)
}

// authHostUnderBiz 空间级权限版本中，find_module_host_relation、find_host_by_service_template、find_host_by_set_template、
// list_biz_hosts、list_biz_hosts_topo、find_host_by_topo、list_host_total_mainline_topo这几个上esb接口, 可以通过配置变量，
// 决定是否鉴业务访问权限
func (s *Service) authHostUnderBiz(kit *rest.Kit, bizID int64) (*meta.BaseResp, bool) {
	if !s.AuthManager.Enabled() {
		return nil, true
	}

	config := "authServer.skipViewBizAuth"
	skipAuth, err := cc.Bool(config)
	if err != nil {
		blog.Errorf("get config %s failed, err: %v, rid: %s", config, err, kit.Rid)
	}

	if skipAuth {
		return nil, true
	}

	authRes := acMeta.ResourceAttribute{Basic: acMeta.Basic{Type: acMeta.HostInstance, Action: acMeta.Find},
		BusinessID: bizID}

	return s.AuthManager.Authorize(kit, authRes)
}
