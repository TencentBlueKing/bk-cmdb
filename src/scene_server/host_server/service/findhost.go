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
	"sort"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
)

// FindModuleHostRelation find host with module by module id
func (s *Service) FindModuleHostRelation(ctx *rest.Contexts) {
	req := ctx.Request
	defErr := ctx.Kit.CCError

	bizID, err := util.GetInt64ByInterface(req.PathParameter("bk_biz_id"))
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", req.PathParameter("bk_biz_id"), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}
	if bizID == 0 {
		ctx.RespAutoError(defErr.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	body := new(meta.FindModuleHostRelationParameter)
	if err := ctx.DecodeInto(body); err != nil {
		ctx.RespAutoError(err)
		return
	}
	rawErr := body.Validate()
	if rawErr.ErrCode != 0 {
		blog.ErrorJSON("validate request body err: %s, body: %s, rid: %s", rawErr.ToCCError(defErr).Error(), *body, ctx.Kit.Rid)
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
		blog.Errorf("findDistinctHostInfo failed, err: %s, distinctHostCond: %s, searchHostCond: %s, rid:%s", err.Error(), *distinctHostCond, *searchHostCond, ctx.Kit.Rid)
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
	hostModuleConfig, err := s.Logic.GetConfigByCond(ctx.Kit, meta.HostModuleRelationRequest{HostIDArr: hostIDArr,
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
		blog.Errorf("GetModuleMapByCond failed, err: %s, moduleIDArr: %v, rid:%s", err.Error(), moduleIDArr, ctx.Kit.Rid)
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
		blog.Errorf("FindHostsByServiceTemplates failed, Validate err: %v, option:%#v, rid:%s", rawErr.ToCCError(defErr).Error(), *option, ctx.Kit.Rid)
		ctx.RespAutoError(rawErr.ToCCError(defErr))
		return
	}

	bizID, err := util.GetInt64ByInterface(ctx.Request.PathParameter("bk_biz_id"))
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
	moduleIDArr, err := s.Logic.GetModuleIDByCond(ctx.Kit, moduleCond)
	if err != nil {
		blog.Errorf("FindHostsByServiceTemplates failed, GetModuleIDByCond err:%s, cond:%#v, rid:%s", err.Error(), moduleCond, ctx.Kit.Rid)
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
		Fields: option.Fields,
		Page:   option.Page,
	}

	result, err := s.findDistinctHostInfo(ctx, distinctHostCond, searchHostCond)
	if err != nil {
		blog.Errorf("FindHostsByServiceTemplates failed, findDistinctHostInfo err: %s, distinctHostCond:%#v, searchHostCond:%#v, rid:%s", err.Error(), *distinctHostCond, *searchHostCond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result.Data)
}

// findDistinctHostInfo find distinct host info
func (s *Service) findDistinctHostInfo(ctx *rest.Contexts, distinctHostCond *meta.DistinctHostIDByTopoRelationRequest, searchHostCond *meta.QueryCondition) (*meta.SearchHostResult, error) {

	defErr := ctx.Kit.CCError

	hmResult, err := s.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(ctx.Kit.Ctx, ctx.Kit.Header, distinctHostCond)
	if err != nil {
		blog.Errorf("findDistinctHostInfo failed, GetDistinctHostIDByTopology error: %v, input:%#v, rid: %s", err, *hmResult, ctx.Kit.Rid)
		return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !hmResult.Result {
		blog.Errorf("findDistinctHostInfo failed, GetDistinctHostIDByTopology error: %v, input:%#v, rid: %s", hmResult.ErrMsg, *hmResult, ctx.Kit.Rid)
		return nil, hmResult.CCError()
	}

	allHostIDs := hmResult.Data.IDArr
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
		blog.Errorf("FindHostsBySetTemplates failed, Validate err: %v, option:%#v, rid:%s", rawErr.ToCCError(defErr).Error(), *option, ctx.Kit.Rid)
		ctx.RespAutoError(rawErr.ToCCError(defErr))
		return
	}

	bizID, err := util.GetInt64ByInterface(ctx.Request.PathParameter("bk_biz_id"))
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

	setIDArr, err := s.Logic.GetSetIDByCond(ctx.Kit, setCond)
	if err != nil {
		blog.Errorf("FindHostsBySetTemplates failed, GetSetIDByCond err:%s, cond:%#v, rid:%s", err.Error(), setCond, ctx.Kit.Rid)
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
		blog.Errorf("FindHostsBySetTemplates failed, findDistinctHostInfo err: %s, distinctHostCond:%#v, searchHostCond:%#v, rid:%s", err.Error(), *distinctHostCond, *searchHostCond, ctx.Kit.Rid)
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

	// generate search condition, if node is not a set or a module, we need to traverse its child topo to the set level to get hosts by relation
	distinctHostCond := &meta.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{bizID},
	}

	if option.ObjID == common.BKInnerObjIDSet {
		distinctHostCond.SetIDArr = []int64{option.InstID}
	} else if option.ObjID == common.BKInnerObjIDModule {
		distinctHostCond.ModuleIDArr = []int64{option.InstID}
	} else {
		setIDArr, err := s.Logic.GetSetIDsByTopo(ctx.Kit, option.ObjID, option.InstID)
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
	appResult, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, header, common.BKInnerObjIDApp, filter)
	if err != nil {
		blog.Errorf("ListResourcePoolHosts failed, ReadInstance of default app failed, filter: %+v, err: %#v, rid:%s", filter, err, rid)
		ccErr := defErr.Error(common.CCErrCommHTTPDoRequestFailed)
		ctx.RespAutoError(ccErr)
		return
	}
	if ccErr := appResult.CCError(); ccErr != nil {
		blog.ErrorJSON("ListResourcePoolHosts failed, ReadInstance of default app failed, filter: %s, result: %s, rid:%s", filter, appResult, rid)
		ccErr := defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		ctx.RespAutoError(ccErr)
		return
	}
	if appResult.Data.Count == 0 {
		blog.Errorf("ListResourcePoolHosts failed, get default app failed, not found, rid: %s", rid)
		ccErr := defErr.Error(common.CCErrCommBizNotFoundError)
		ctx.RespAutoError(ccErr)
		return
	}

	// only use biz with same supplier account if query returns multiple biz
	bizData := appResult.Data.Info[0]
	bizCount := 0
	for _, biz := range appResult.Data.Info {
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
		blog.ErrorJSON("ListResourcePoolHosts failed, parse app data failed, bizData: %s, err: %s, rid: %s", bizData, err.Error(), rid)
		ccErr := defErr.Error(common.CCErrCommParseDataFailed)
		ctx.RespAutoError(ccErr)
		return
	}

	// do host search
	hostResult, ccErr := s.listBizHosts(ctx, bizID, parameter)
	if ccErr != nil {
		blog.ErrorJSON("ListResourcePoolHosts failed, listBizHosts failed, bizID: %s, parameter: %s, err: %s, rid:%s", bizID, parameter, ccErr.Error(), rid)
		ctx.RespAutoError(ccErr)
		return
	}
	ctx.RespEntity(hostResult)
}

// ListHosts list host under business specified by path parameter
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
		blog.ErrorJSON("ListBizHosts failed, Validate failed,parameter:%s, err: %s, rid:%s", parameter, err, ctx.Kit.Rid)
		ccErr := defErr.CCErrorf(common.CCErrCommParamsInvalid, key)
		ctx.RespAutoError(ccErr)
		return
	}
	bizID, err := util.GetInt64ByInterface(req.PathParameter("appid"))
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
	ctx.SetReadPreference(common.SecondaryPreferredMode)
	hostResult, ccErr := s.listBizHosts(ctx, bizID, parameter)
	if ccErr != nil {
		blog.ErrorJSON("ListBizHosts failed, listBizHosts failed, bizID: %s, parameter: %s, err: %s, rid:%s", bizID, parameter, ccErr.Error(), rid)
		ctx.RespAutoError(ccErr)
		return
	}
	ctx.RespEntity(hostResult)
}

func (s *Service) listBizHosts(ctx *rest.Contexts, bizID int64, parameter meta.ListHostsParameter) (result *meta.ListHostResult, ccErr errors.CCErrorCoder) {
	header := ctx.Kit.Header
	rid := ctx.Kit.Rid
	defErr := ctx.Kit.CCError

	if parameter.Page.IsIllegal() {
		blog.Errorf("ListBizHosts failed, page limit %d illegal, rid:%s", parameter.Page.Limit, ctx.Kit.Rid)
		return result, defErr.CCErrorf(common.CCErrCommParamsInvalid, "page.limit")
	}

	if parameter.SetIDs != nil && len(parameter.SetIDs) != 0 && parameter.SetCond != nil && len(parameter.SetCond) != 0 {
		blog.Errorf("ListBizHosts failed, bk_set_ids and set_cond can't both be set, rid:%s", ctx.Kit.Rid)
		return result, defErr.CCErrorf(common.CCErrCommParamsInvalid, "bk_set_ids and set_cond can't both be set")
	}

	setIDList := make([]int64, 0)
	if parameter.SetCond != nil {
		setCond := make(map[string]interface{})
		if err := parse.ParseCommonParams(parameter.SetCond, setCond); err != nil {
			blog.Errorf("parse set cond failed, err: %v, rid: %s", err, rid)
			return nil, errors.New(common.CCErrCommParamsInvalid, "set_cond")
		}

		// set the app id condition
		setCond[common.BKAppIDField] = bizID
		query := meta.QueryCondition{
			Fields:    []string{common.BKSetIDField},
			Condition: setCond,
		}

		setList, setErr := s.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, header, common.BKInnerObjIDSet, &query)
		if setErr != nil {
			blog.Errorf("get set with cond: %v failed, err: %v, rid: %s", setCond, setErr, rid)
			return nil, errors.New(common.CCErrCommParamsInvalid, "set_cond")
		}

		if !setList.Result {
			blog.Errorf("get set with cond: %v failed, err: %v, rid: %s", setCond, setErr, rid)
			return nil, errors.New(setList.Code, setList.ErrMsg)
		}

		if len(setList.Data.Info) == 0 {
			return &meta.ListHostResult{
				Count: 0,
				Info:  []map[string]interface{}{},
			}, nil
		}

		for _, set := range setList.Data.Info {
			id, err := util.GetInt64ByInterface(set[common.BKSetIDField])
			if err != nil {
				blog.Errorf("get set id: %v failed, err: %v, rid: %s", set[common.BKSetIDField], err, rid)
				return nil, errors.New(common.CCErrCommParamsInvalid, "bk_set_id")
			}

			if id == 0 {
				continue
			}

			setIDList = append(setIDList, id)
		}
	}

	if len(parameter.SetIDs) != 0 {
		setIDList = append(setIDList, parameter.SetIDs...)
	}

	option := &meta.ListHosts{
		BizID:              bizID,
		SetIDs:             setIDList,
		ModuleIDs:          parameter.ModuleIDs,
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
		blog.ErrorJSON("ListHostsWithNoBiz failed, decode body failed,parameter:%s, err: %#v, rid:%s", parameter, err, ctx.Kit.Rid)
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
	header := ctx.Kit.Header
	rid := ctx.Kit.Rid
	defErr := ctx.Kit.CCError

	parameter := &meta.ListHostsWithNoBizParameter{}
	if err := ctx.DecodeInto(&parameter); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if key, err := parameter.Validate(); err != nil {
		blog.ErrorJSON("ListHostByTopoNode failed, Validate failed,parameter:%s, err: %s, rid:%s", parameter, err, ctx.Kit.Rid)
		ccErr := defErr.CCErrorf(common.CCErrCommParamsInvalid, key)
		ctx.RespAutoError(ccErr)
		return
	}
	bizID, err := util.GetInt64ByInterface(ctx.Request.PathParameter("bk_biz_id"))
	if err != nil {
		ctx.RespAutoError(defErr.Errorf(common.CCErrCommParamsInvalid, "bk_app_id"))
		return
	}
	if bizID == 0 {
		ctx.RespAutoError(defErr.Errorf(common.CCErrCommParamsInvalid, "bk_app_id"))
		return
	}

	if parameter.Page.IsIllegal() {
		blog.Errorf("ListHostByTopoNode failed, page limit %d illegal, rid:%s", parameter.Page.Limit, ctx.Kit.Rid)
		ctx.RespAutoError(defErr.CCErrorf(common.CCErrCommParamsInvalid, "page.limit"))
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	// search all hosts
	option := &meta.ListHosts{
		BizID:              bizID,
		HostPropertyFilter: parameter.HostPropertyFilter,
		Fields:             append(parameter.Fields, common.BKHostIDField),
		Page:               parameter.Page,
	}
	hosts, err := s.CoreAPI.CoreService().Host().ListHosts(ctx.Kit.Ctx, header, option)
	if err != nil {
		blog.Errorf("find host failed, err: %s, input:%#v, rid:%s", err.Error(), parameter, rid)
		ctx.RespAutoError(defErr.Error(common.CCErrHostGetFail))
		return
	}

	if len(hosts.Info) == 0 {
		ctx.RespEntity(hosts)
		return
	}

	// search all hosts' host module relations
	hostIDs := make([]int64, 0)
	for _, host := range hosts.Info {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.ErrorJSON("host: %s bk_host_id field invalid, rid:%s", host, rid)
			ctx.RespAutoError(err)
			return
		}
		hostIDs = append(hostIDs, hostID)
	}
	relationCond := meta.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     hostIDs,
		Fields:        []string{common.BKSetIDField, common.BKModuleIDField, common.BKHostIDField},
	}
	relations, err := s.Logic.GetConfigByCond(ctx.Kit, relationCond)
	if nil != err {
		blog.ErrorJSON("read host module relation error: %s, input: %s, rid: %s", err, hosts, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// search all module and set info
	setIDs := make([]int64, 0)
	moduleIDs := make([]int64, 0)
	relation := make(map[int64]map[int64][]int64)
	for _, r := range relations {
		setIDs = append(setIDs, r.SetID)
		moduleIDs = append(moduleIDs, r.ModuleID)
		if setModule, ok := relation[r.HostID]; ok {
			setModule[r.SetID] = append(setModule[r.SetID], r.ModuleID)
			relation[r.HostID] = setModule
		} else {
			setModule := make(map[int64][]int64)
			setModule[r.SetID] = append(setModule[r.SetID], r.ModuleID)
			relation[r.HostID] = setModule
		}
	}
	setIDs = util.IntArrayUnique(setIDs)
	moduleIDs = util.IntArrayUnique(moduleIDs)

	cond := condition.CreateCondition()
	cond.Field(common.BKSetIDField).In(setIDs)
	query := &meta.QueryCondition{
		Fields:    []string{common.BKSetIDField, common.BKSetNameField},
		Condition: cond.ToMapStr(),
	}
	sets, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.ErrorJSON("get set by condition: %s failed, err: %s, rid: %s", cond.ToMapStr(), err, rid)
		ctx.RespAutoError(err)
		return
	}
	setMap := make(map[int64]string)
	for _, set := range sets.Data.Info {
		setID, err := set.Int64(common.BKSetIDField)
		if err != nil {
			blog.ErrorJSON("set %s id invalid, error: %s, rid: %s", set, err, rid)
			ctx.RespAutoError(err)
			return
		}
		setName, err := set.String(common.BKSetNameField)
		if err != nil {
			blog.ErrorJSON("set %s name invalid, error: %s, rid: %s", set, err, rid)
			ctx.RespAutoError(err)
			return
		}
		setMap[setID] = setName
	}

	cond = condition.CreateCondition()
	cond.Field(common.BKModuleIDField).In(moduleIDs)
	query = &meta.QueryCondition{
		Fields:    []string{common.BKModuleIDField, common.BKModuleNameField},
		Condition: cond.ToMapStr(),
	}
	modules, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.ErrorJSON("get module by condition: %s failed, err: %s, rid: %s", cond.ToMapStr(), err, rid)
		ctx.RespAutoError(err)
		return
	}
	moduleMap := make(map[int64]string)
	for _, module := range modules.Data.Info {
		moduleID, err := module.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("module %s id invalid, error: %s, rid: %s", module, err, rid)
			ctx.RespAutoError(err)
			return
		}
		moduleName, err := module.String(common.BKModuleNameField)
		if err != nil {
			blog.ErrorJSON("module %s name invalid, error: %s, rid: %s", module, err, rid)
			ctx.RespAutoError(err)
			return
		}
		moduleMap[moduleID] = moduleName
	}

	// format the output
	hostTopos := meta.HostTopoResult{
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
	ctx.RespEntity(hostTopos)

}

func (s *Service) CountTopoNodeHosts(ctx *rest.Contexts) {

	rid := ctx.Kit.Rid
	defErr := ctx.Kit.CCError

	option := meta.CountTopoNodeHostsOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := util.GetInt64ByInterface(ctx.Request.PathParameter(common.BKAppIDField))
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
		blog.ErrorJSON("CountTopoNodeHosts failed, countTopoNodeHosts failed, option: %s, err: %s, rid:%s", option, ccErr.Error(), rid)
		ctx.RespAutoError(defErr.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	ctx.RespEntity(topoNodeHostCounts)
}

func (s *Service) countTopoNodeHosts(ctx *rest.Contexts, bizID int64, option meta.CountTopoNodeHostsOption) ([]meta.TopoNodeHostCount, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid
	topoRoot, ccErr := s.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(ctx.Kit.Ctx, ctx.Kit.Header, bizID, false)
	if ccErr != nil {
		blog.Errorf("countTopoNodeHosts failed, SearchMainlineInstanceTopo failed, bizID: %d, err: %s, rid: %s", bizID, ccErr.Error(), rid)
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
	relationResult, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header, &relationOption)
	if err != nil {
		blog.Errorf("countTopoNodeHosts failed, GetHostModuleRelation failed, option: %+v, err: %s, rid: %s", relationOption, err.Error(), rid)
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
		for _, item := range relationResult.Data.Info {
			if _, ok := moduleIDMap[item.ModuleID]; ok == true {
				hostIDs = append(hostIDs, item.HostID)
			}
		}
		hostCount.HostCount = len(util.IntArrayUnique(hostIDs))
		hostCounts = append(hostCounts, hostCount)
	}

	return hostCounts, nil
}
