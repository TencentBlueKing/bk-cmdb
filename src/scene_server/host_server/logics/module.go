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

package logics

import (
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	types "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

// GetResourcePoolModuleID get module id,module name.
func (lgc *Logics) GetResourcePoolModuleID(kit *rest.Kit, condition mapstr.MapStr) (int64, string, errors.CCError) {
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKModuleIDField, common.BkSupplierAccount, common.BKModuleNameField},
		Condition: condition,
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		query)
	if err != nil {
		blog.Errorf("GetResourcePoolModuleID http do error, err:%s,input:%+v,rid:%s", err.Error(), query, kit.Rid)
		return -1, "", kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(result.Info) == 0 {
		blog.Errorf("GetResourcePoolModuleID http response error,input:%+v,rid:%s", query, kit.Rid)
		return -1, "", kit.CCError.Error(common.CCErrTopoGetAppFailed)
	}

	supplier := kit.SupplierAccount
	for idx, mod := range result.Info {
		if supplier == mod[common.BkSupplierAccount].(string) {
			moduleId, err := result.Info[idx].Int64(common.BKModuleIDField)
			if err != nil {
				return moduleId, "", err
			}

			moduleName := ""
			if name, ok := result.Info[idx][common.BKModuleNameField].(string); ok {
				moduleName = name
				return moduleId, moduleName, nil
			}
			return -1, "", kit.CCError.Error(common.CCErrTopoGetModuleFailed)
		}
	}

	blog.Errorf("can not get resource pool module id rid:%s", kit.Rid)
	return -1, "", kit.CCError.Error(common.CCErrTopoGetAppFailed)
}

// GetNormalModuleByModuleID TODO
func (lgc *Logics) GetNormalModuleByModuleID(kit *rest.Kit, appID, moduleID int64) ([]mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryCondition{
		Page:      metadata.BasePage{Start: 0, Limit: 1, Sort: common.BKModuleIDField},
		Fields:    []string{common.BKModuleIDField},
		Condition: hutil.NewOperation().WithAppID(appID).WithModuleID(moduleID).Data(),
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		query)
	if err != nil {
		return nil, kit.CCError.Errorf(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	return result.Info, nil
}

// GetModuleIDByCond TODO
func (lgc *Logics) GetModuleIDByCond(kit *rest.Kit, cond metadata.ConditionWithTime) (
	[]int64, errors.CCError) {

	condc := make(map[string]interface{})
	if err := parse.ParseCommonParams(cond.Condition, condc); err != nil {
		blog.Errorf("ParseCommonParams failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommUtilHandleFail, "parse condition", err.Error())
	}

	query := &metadata.QueryCondition{
		Page:          metadata.BasePage{Start: 0, Limit: common.BKNoLimit, Sort: common.BKModuleIDField},
		Fields:        []string{common.BKModuleIDField},
		Condition:     mapstr.NewFromMap(condc),
		TimeCondition: cond.TimeCondition,
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		query)
	if err != nil {
		blog.Errorf("GetModuleIDByCond http do error, err:%s,input:%+v,rid:%s", err.Error(), query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	moduleIDArr := make([]int64, 0)
	for _, i := range result.Info {
		moduleID, err := i.Int64(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("GetModuleIDByCond convert  module id to int error, err:%s, module:%+v,input:%+v,rid:%s",
				err.Error(), i, query, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDModule,
				common.BKModuleIDField, "int", err.Error())
		}
		moduleIDArr = append(moduleIDArr, moduleID)
	}
	return moduleIDArr, nil
}

// GetModuleMapByCond search object map by condition
func (lgc *Logics) GetModuleMapByCond(kit *rest.Kit, fields []string, cond mapstr.MapStr) (map[int64]types.MapStr,
	errors.CCError) {

	query := &metadata.QueryCondition{
		Condition: cond,
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit, Sort: common.BKModuleIDField},
		Fields:    fields,
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		query)
	if err != nil {
		blog.Errorf("GetModuleMapByCond http do error, err:%s,input:%+v,rid:%s", err.Error(), query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	moduleMap := make(map[int64]types.MapStr)
	for _, info := range result.Info {
		id, err := info.Int64(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("GetModuleMapByCond convert  module id to int error, err:%s, module:%+v,input:%+v,rid:%s",
				err.Error(), info, query, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDModule,
				common.BKModuleIDField, "int", err.Error())
		}
		moduleMap[id] = info
	}

	return moduleMap, nil
}

// GetModuleIDAndIsInternal TODO
func (lgc *Logics) GetModuleIDAndIsInternal(kit *rest.Kit, bizID, moduleID int64) (int64, bool, error) {
	if moduleID == 0 {
		cond := map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKDefaultField: common.DefaultResModuleFlag,
		}
		moduleID, _, err := lgc.GetResourcePoolModuleID(kit, cond)
		if err != nil {
			blog.Errorf("GetModuleIDAndIsInternal get default moduleID failed, err: %s, bizID: %d, rid: %s", err.Error(), bizID, kit.Rid)
			return 0, false, err
		}
		return moduleID, true, nil
	} else {
		cond := map[string]interface{}{
			common.BKAppIDField:    bizID,
			common.BKModuleIDField: moduleID,
		}
		moduleMap, err := lgc.GetModuleMapByCond(kit, []string{common.BKDefaultField, common.BKModuleIDField}, cond)
		if err != nil {
			blog.Errorf("GetModuleIDAndIsInternal get module info failed, err: %s, bizID: %d, moduleID: %d, rid: %s", err.Error(), bizID, moduleID, kit.Rid)
			return 0, false, err
		}
		module, ok := moduleMap[moduleID]
		if !ok {
			blog.Errorf("GetModuleIDAndIsInternal module not exist, bizID: %d, moduleID: %d, rid: %s", bizID, moduleID, kit.Rid)
			return 0, false, kit.CCError.CCErrorf(common.CCErrHostModuleNotBelongBusinessErr, moduleID, bizID)
		}
		return moduleID, module[common.BKDefaultField] != common.NormalModuleFlag, nil
	}
}

// MoveHostToResourcePool transfer hosts to a resource pool
func (lgc *Logics) MoveHostToResourcePool(kit *rest.Kit, conf *metadata.DefaultModuleHostConfigParams) ([]metadata.ExceptionResult, error) {

	ownerAppID, err := lgc.GetDefaultAppID(kit)
	if err != nil {
		blog.Errorf("move host to resource pool, but get default appid failed, err: %v, input:%+v,rid:%s", err, conf, kit.Rid)
		return nil, err
	}
	if 0 == conf.ApplicationID {
		return nil, kit.CCError.Error(common.CCErrHostNotResourceFail)
	}
	if ownerAppID == conf.ApplicationID {
		return nil, kit.CCError.Errorf(common.CCErrHostBelongResourceFail)
	}

	ownerModuleIDCond := map[string]interface{}{
		common.BKAppIDField: ownerAppID,
	}

	// if directory id is specified, transfer to it, if not, transfer host to the default directory
	if conf.ModuleID == 0 {
		ownerModuleIDCond[common.BKDefaultField] = common.DefaultResModuleFlag
	} else {
		ownerModuleIDCond[common.BKModuleIDField] = conf.ModuleID
	}

	ownerModuleID, _, err := lgc.GetResourcePoolModuleID(kit, ownerModuleIDCond)
	if err != nil {
		blog.Errorf("move host to resource pool, but get module id failed, err: %v, input:%+v,param:%+v,rid:%s", err, conf, ownerModuleIDCond, kit.Rid)
		return nil, err
	}

	cond := metadata.ConditionWithTime{
		Condition: []metadata.ConditionItem{
			{Field: common.BKDefaultField, Operator: common.BKDBNE, Value: common.NormalModuleFlag},
			{Field: common.BKAppIDField, Operator: common.BKDBEQ, Value: conf.ApplicationID},
		},
	}
	moduleIDs, err := lgc.GetModuleIDByCond(kit, cond)
	if err != nil {
		blog.Errorf("move host to resource pool, but get module ids failed, input: %+v, param: %+v, err: %v, rid: %s",
			conf, cond, err, kit.Rid)
		return nil, err
	}

	if len(moduleIDs) == 0 {
		blog.Errorf("move host to resource pool, module not found, input: %+v, param: %+v, rid: %s",
			conf, cond, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrHostModuleNotExist, "idle pool")
	}

	errHostID, err := lgc.notExistAppModuleHost(kit, []int64{conf.ApplicationID}, moduleIDs, conf.HostIDs)
	if err != nil {
		blog.Errorf("move host to resource pool, notExistAppModuleHost error, err: %v, owneAppID: %d, input:%#v, rid:%s", err, ownerAppID, conf, kit.Rid)
		return nil, err
	}
	if len(errHostID) > 0 {
		errHostIP := lgc.convertHostIDToHostIP(kit, errHostID)
		blog.Errorf("move host to resource pool, notExistAppModuleHost error, has host not belong to idle module , owneAppID: %d, input:%#v, err host inner ip:%#v, rid:%s", ownerAppID, conf, errHostIP, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrHostModuleConfigNotMatch, util.PrettyIPStr(errHostIP))
	}

	param := &metadata.TransferHostsCrossBusinessRequest{
		SrcApplicationIDs: []int64{conf.ApplicationID},
		HostIDArr:         conf.HostIDs,
		DstApplicationID:  ownerAppID,
		DstModuleIDArr:    []int64{ownerModuleID},
	}

	audit := auditlog.NewHostModuleLog(lgc.CoreAPI.CoreService(), conf.HostIDs)
	if err := audit.WithPrevious(kit); err != nil {
		blog.Errorf("move host to resource pool, but get prev module host config failed, err: %v, input:%+v,rid:%s", err, conf, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}
	result, err := lgc.CoreAPI.CoreService().Host().TransferToAnotherBusiness(kit.Ctx, kit.Header, param)
	if err != nil {
		blog.Errorf("move host to resource pool, but update host module http do error, err: %v, input:%#v,params:%#v,rid:%v", err, conf, param, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("move host to resource pool, but update host module http response error, err code:%d, err messge:%s, input:%#v,query:%#v,rid:%v", result.Code, result.ErrMsg, conf, param, kit.Rid)
		return result.Data, kit.CCError.New(result.Code, result.ErrMsg)

	}

	if err := audit.SaveAudit(kit); err != nil {
		blog.Errorf("move host to resource pool, but save audit log failed, err: %v, input:%+v,rid:%s", err, conf, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}

	return nil, nil
}

// notExistAppModuleHost get hostIDs those don't exist in the modules
func (lgc *Logics) notExistAppModuleHost(kit *rest.Kit, appIDs []int64, moduleIDs []int64, hostIDArr []int64) ([]int64,
	errors.CCErrorCoder) {

	hostModuleInput := &metadata.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: appIDs,
		ModuleIDArr:      moduleIDs,
		HostIDArr:        hostIDArr,
	}

	hmResult, err := lgc.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(kit.Ctx, kit.Header, hostModuleInput)
	if err != nil {
		blog.ErrorJSON("get host ID by topology failed, err: %s, input:%s,rid:%s", err, hostModuleInput, kit.Rid)
		return nil, err
	}

	hostIDMap := make(map[int64]bool, 0)
	for _, id := range hmResult {
		hostIDMap[id] = true
	}
	var errHostIDArr []int64
	for _, hostID := range hostIDArr {
		if _, ok := hostIDMap[hostID]; !ok {
			errHostIDArr = append(errHostIDArr, hostID)
		}
	}

	return errHostIDArr, nil
}

// AssignHostToApp transfer resource host to  idle module
func (lgc *Logics) AssignHostToApp(kit *rest.Kit, conf *metadata.DefaultModuleHostConfigParams) ([]metadata.ExceptionResult, error) {

	cond := hutil.NewOperation().WithAppID(conf.ApplicationID).Data()
	fields := fmt.Sprintf("%s,%s", common.BKOwnerIDField, common.BKAppNameField)
	appInfo, err := lgc.GetAppDetails(kit, fields, cond)
	if err != nil {
		blog.Errorf("assign host to app failed, err: %v,input:%+v,rid:%s", err, conf, kit.Rid)
		return nil, err
	}
	if 0 == len(appInfo) {
		blog.Errorf("assign host to app error, not foud app appID: %d,input:%+v,rid:%s", conf.ApplicationID, conf, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommNotFound)
	}

	ownerAppID, err := lgc.GetDefaultAppID(kit)
	if err != nil {
		blog.Errorf("assign host to app, but get default appid failed, err: %v,input:%+v,rid:%s", err, conf, kit.Rid)
		return nil, err
	}
	if 0 == conf.ApplicationID {
		return nil, kit.CCError.Errorf(common.CCErrHostGetResourceFail, "not found")
	}
	if ownerAppID == conf.ApplicationID {
		return nil, nil
	}
	moduleCond := []metadata.ConditionItem{
		{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    ownerAppID,
		},
		{
			Field:    common.BKDefaultField,
			Operator: common.BKDBIN,
			Value:    []int{common.DefaultResModuleFlag, common.DefaultResSelfDefinedModuleFlag},
		},
	}

	resourceModuleIDs, err := lgc.GetModuleIDByCond(kit, metadata.ConditionWithTime{Condition: moduleCond})
	if err != nil {
		blog.Errorf("assign host to app failed, GetModuleIDByCond err: %v, moduleCond:%+v, rid:%s", err, moduleCond, kit.Rid)
		return nil, err
	}

	errHostID, err := lgc.notExistAppModuleHost(kit, []int64{ownerAppID}, resourceModuleIDs, conf.HostIDs)
	if err != nil {
		blog.Errorf("assign host to app, notExistAppModuleHost error, err: %v, input:%+v, rid:%s", err, conf, kit.Rid)
		return nil, err
	}
	if len(errHostID) > 0 {
		errHostIP := lgc.convertHostIDToHostIP(kit, errHostID)
		blog.Errorf("assign host to app, notExistAppModuleHost error, has host not belong to resource directory module , input:%+v, rid:%s", conf, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrHostModuleConfigNotMatch, strings.Join(errHostIP, ","))
	}

	mConds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithAppID(conf.ApplicationID)
	moduleID, moduleName, err := lgc.GetResourcePoolModuleID(kit, mConds.MapStr())
	if err != nil {
		blog.Errorf("assign host to app, but get module id failed, err: %v,input:%+v,params:%+v,rid:%s", err, conf, mConds.MapStr(), kit.Rid)
		return nil, err
	}
	if moduleID == 0 {
		blog.Errorf("assign host to app, but get module id failed, %s not found,input: %+v,params:% +v,rid: %s",
			moduleName, conf, mConds.MapStr(), kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrHostModuleNotExist, moduleName)
	}

	assignParams := &metadata.TransferHostsCrossBusinessRequest{
		SrcApplicationIDs: []int64{ownerAppID},
		DstApplicationID:  conf.ApplicationID,
		HostIDArr:         conf.HostIDs,
		DstModuleIDArr:    []int64{moduleID},
	}

	audit := auditlog.NewHostModuleLog(lgc.CoreAPI.CoreService(), conf.HostIDs)
	if err := audit.WithPrevious(kit); err != nil {
		blog.Warnf("WithPrevious failed, err: %+v, rid: %s", err, kit.Rid)
	}

	result, err := lgc.CoreAPI.CoreService().Host().TransferToAnotherBusiness(kit.Ctx, kit.Header, assignParams)
	if err != nil {
		blog.Errorf("assign host to app, but assign to app http do error. err: %v, input:%+v,param:%+v,rid:%s", err, conf, assignParams, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrHostEditRelationPoolFail)
	}
	if !result.Result {
		blog.Errorf("assign host to app, but assign to app http response error. result:%#v, input:%+v, param:%+v, rid:%s", result, conf, assignParams, kit.Rid)
		return result.Data, kit.CCError.New(result.Code, result.ErrMsg)
	}

	if err := audit.SaveAudit(kit); err != nil {
		blog.Errorf("assign host to app, but save audit failed, err: %v, rid:%s", err, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}

	return nil, nil
}

// convertHostIDToHostIP  该方法为专用方法。出现任何错误都会被忽略。
// 尝试将主机ID转换为内网IP，如果转换中出现问题返回主机ID。
func (lgc *Logics) convertHostIDToHostIP(kit *rest.Kit, hostIDArr []int64) []string {

	if len(hostIDArr) == 0 {
		return nil
	}
	cond := condition.CreateCondition()
	cond.Field(common.BKHostIDField).In(hostIDArr)
	input := &metadata.QueryCondition{
		Condition: cond.ToMapStr(),
		Fields:    []string{common.BKHostIDField, common.BKHostInnerIPField},
	}

	// 找不到主机ID对应的IP， 返回主机ID
	hostIDIPMap := make(map[int64]string, 0)
	for _, hostID := range hostIDArr {
		hostIDIPMap[hostID] = strconv.FormatInt(hostID, 10)
	}

	result, err := lgc.CoreAPI.
		CoreService().
		Instance().
		ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Warnf("convertHostIDToHostIP http do error. err:%s, input:%#v, rid:%s", err.Error(), input, kit.Rid)
	}

	for _, host := range result.Info {
		hostID, err := host.Int64(common.BKHostIDField)
		if err != nil {
			// can't not foud host id , skip
			blog.Warnf("convertHostIDToHostIP convert host id to int64 error. err:%s, host:%#v, input:%#v, rid:%s",
				err.Error(), host, input, kit.Rid)
			continue
		}
		innerIP, err := host.String(common.BKHostInnerIPField)
		if err != nil {
			// can't not foud host inner ip , skip
			blog.Warnf("convertHostIDToHostIP convert host inner ip to string error. err:%s, host:%#v, input:%#v, "+
				"rid:%s", err.Error(), host, input, kit.Rid)
			continue
		}
		hostIDIPMap[hostID] = innerIP
	}
	var ips []string
	for _, ip := range hostIDIPMap {
		ips = append(ips, ip)
	}

	return ips
}

// ValidateHostInModule validate if host are all in modules specified by cond
func (lgc *Logics) ValidateHostInModule(kit *rest.Kit, hostIDs []int64, moduleCond mapstr.MapStr) errors.CCErrorCoder {

	// get module ids by condition
	moduleReq := &metadata.QueryCondition{
		Condition:      moduleCond,
		Fields:         []string{common.BKModuleIDField},
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		DisableCounter: true,
	}

	moduleRes := new(metadata.ResponseModuleInstance)
	err := lgc.CoreAPI.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		moduleReq, &moduleRes)
	if err != nil {
		blog.Errorf("get modules failed, input: %#v, err: %v, rid: %s", moduleReq, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if err = moduleRes.CCError(); err != nil {
		blog.Errorf("get modules failed, input: %#v, err: %v, rid: %s", moduleReq, err, kit.Rid)
		return err
	}

	moduleIDs := make([]int64, len(moduleRes.Data.Info))
	for index, module := range moduleRes.Data.Info {
		moduleIDs[index] = module.ModuleID
	}

	// check if all hosts belongs to the modules
	relationCond := mapstr.MapStr{
		common.BKHostIDField:   mapstr.MapStr{common.BKDBIN: hostIDs},
		common.BKModuleIDField: mapstr.MapStr{common.BKDBNIN: moduleIDs},
	}

	counts, err := lgc.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameModuleHostConfig, []map[string]interface{}{relationCond})
	if err != nil {
		blog.Error("get illegal relation count failed, cond: %+v, err: %v, rid: %s", relationCond, err, kit.Rid)
		return err
	}

	if len(counts) != 1 || int(counts[0]) > 0 {
		blog.Error("not all hosts(%+v) belongs to the modules(%+v), input: %+v, rid: %s", hostIDs, moduleIDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
	}

	return nil
}
