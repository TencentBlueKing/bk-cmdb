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
	"context"
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	types "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

func (lgc *Logics) GetResourcePoolModuleID(ctx context.Context, condition mapstr.MapStr) (int64, errors.CCError) {
	query := &metadata.QueryCondition{
		Limit:     metadata.SearchLimit{Offset: 0, Limit: 1},
		SortArr:   metadata.NewSearchSortParse().String(common.BKModuleIDField).ToSearchSortArr(),
		Fields:    []string{common.BKModuleIDField},
		Condition: condition,
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("GetResourcePoolModuleID http do error, err:%s,input:%+v,rid:%s", err.Error(), query, lgc.rid)
		return -1, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetResourcePoolModuleID http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return -1, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		blog.Errorf("GetResourcePoolModuleID http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return -1, lgc.ccErr.Error(common.CCErrTopoGetAppFailed)
	}

	return result.Data.Info[0].Int64(common.BKModuleIDField)
}

func (lgc *Logics) GetNormalModuleByModuleID(ctx context.Context, appID, moduleID int64) ([]mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryCondition{
		Limit:     metadata.SearchLimit{Offset: 0, Limit: 1},
		SortArr:   metadata.NewSearchSortParse().String(common.BKModuleIDField).ToSearchSortArr(),
		Fields:    []string{common.BKModuleIDField},
		Condition: hutil.NewOperation().WithAppID(appID).WithModuleID(moduleID).Data(),
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("GetNormalModuleByModuleID http do error, err:%s,input:%#v,rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetNormalModuleByModuleID http response error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func (lgc *Logics) GetModuleIDByCond(ctx context.Context, cond []metadata.ConditionItem) ([]int64, errors.CCError) {
	condc := make(map[string]interface{})
	if err := parse.ParseCommonParams(cond, condc); err != nil {
		blog.Warnf("ParseCommonParams failed, err: %+v, rid: %s", err, lgc.rid)
	}

	query := &metadata.QueryCondition{
		Limit:     metadata.SearchLimit{Offset: 0, Limit: common.BKNoLimit},
		SortArr:   metadata.NewSearchSortParse().String(common.BKModuleIDField).ToSearchSortArr(),
		Fields:    []string{common.BKModuleIDField},
		Condition: mapstr.NewFromMap(condc),
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("GetModuleIDByCond http do error, err:%s,input:%+v,rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetModuleIDByCond http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	moduleIDArr := make([]int64, 0)
	for _, i := range result.Data.Info {
		moduleID, err := i.Int64(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("GetModuleIDByCond convert  module id to int error, err:%s, module:%+v,input:%+v,rid:%s", err.Error(), i, query, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		moduleIDArr = append(moduleIDArr, moduleID)
	}
	return moduleIDArr, nil
}

func (lgc *Logics) GetModuleMapByCond(ctx context.Context, fields []string, cond mapstr.MapStr) (map[int64]types.MapStr, errors.CCError) {

	query := &metadata.QueryCondition{
		Condition: cond,
		Limit:     metadata.SearchLimit{Offset: 0, Limit: common.BKNoLimit},
		SortArr:   metadata.NewSearchSortParse().String(common.BKModuleIDField).ToSearchSortArr(),
		Fields:    fields,
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("GetModuleMapByCond http do error, err:%s,input:%+v,rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetModuleMapByCond http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	moduleMap := make(map[int64]types.MapStr)
	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("GetModuleMapByCond convert  module id to int error, err:%s, module:%+v,input:%+v,rid:%s", err.Error(), info, query, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		moduleMap[id] = info
	}

	return moduleMap, nil
}

//
func (lgc *Logics) MoveHostToResourcePool(ctx context.Context, conf *metadata.DefaultModuleHostConfigParams) ([]metadata.ExceptionResult, error) {

	ownerAppID, err := lgc.GetDefaultAppID(ctx)
	if err != nil {
		blog.Errorf("move host to resource pool, but get default appid failed, err: %v, input:%+v,rid:%s", err, conf, lgc.rid)
		return nil, err
	}
	if 0 == conf.ApplicationID {
		return nil, lgc.ccErr.Error(common.CCErrHostNotResourceFail)
	}
	if ownerAppID == conf.ApplicationID {
		return nil, lgc.ccErr.Errorf(common.CCErrHostBelongResourceFail)
	}
	owenerModuleIDconds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(ownerAppID)
	ownerModuleID, err := lgc.GetResourcePoolModuleID(ctx, owenerModuleIDconds.MapStr())
	if err != nil {
		blog.Errorf("move host to resource pool, but get module id failed, err: %v, input:%+v,param:%+v,rid:%s", err, conf, owenerModuleIDconds.Data(), lgc.rid)
		return nil, err
	}

	conds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(conf.ApplicationID)
	moduleID, err := lgc.GetResourcePoolModuleID(ctx, conds.MapStr())
	if err != nil {
		blog.Errorf("move host to resource pool, but get module id failed, err: %v, input:%+v,param:%+v,rid:%s", err, conf, conds.Data(), lgc.rid)
		return nil, err
	}
	errHostID, err := lgc.notExistAppModuleHost(ctx, conf.ApplicationID, moduleID, conf.HostIDs)
	if err != nil {
		blog.Errorf("move host to resource pool, notExistAppModuleHost error, err: %v, owneAppID: %d, input:%#v, rid:%s", err, ownerAppID, conf, lgc.rid)
		return nil, err
	}
	if len(errHostID) > 0 {
		errHostIP := lgc.convertHostIDToHostIP(ctx, errHostID)
		blog.Errorf("move host to resource pool, notExistAppModuleHost error, has host not belong to idle module , owneAppID: %d, input:%#v, err host inner ip:%#v, rid:%s", ownerAppID, conf, errHostIP, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrHostNotBelongIDLEModuleErr, util.PrettyIPStr(errHostIP))
	}

	param := &metadata.TransferHostsCrossBusinessRequest{
		SrcApplicationID: conf.ApplicationID,
		HostIDArr:        conf.HostIDs,
		DstApplicationID: ownerAppID,
		DstModuleIDArr:   []int64{ownerModuleID},
	}

	audit := lgc.NewHostModuleLog(conf.HostIDs)
	if err := audit.WithPrevious(ctx); err != nil {
		blog.Errorf("move host to resource pool, but get prev module host config failed, err: %v, input:%+v,rid:%s", err, conf, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}
	result, err := lgc.CoreAPI.CoreService().Host().TransferToAnotherBusiness(ctx, lgc.header, param)
	if err != nil {
		blog.Errorf("move host to resource pool, but update host module http do error, err: %v, input:%#v,params:%#v,rid:%v", err, conf, param, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("move host to resource pool, but update host module http response error, err code:%d, err messge:%s, input:%#v,query:%#v,rid:%v", result.Code, result.ErrMsg, conf, param, lgc.rid)
		return result.Data, lgc.ccErr.New(result.Code, result.ErrMsg)

	}

	if err := audit.SaveAudit(ctx, conf.ApplicationID, lgc.user, "move host to resource pool"); err != nil {
		blog.Errorf("move host to resource pool, but save audit log failed, err: %v, input:%+v,rid:%s", err, conf, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}
	businessMetadata := conf.Metadata
	if businessMetadata.Label == nil {
		businessMetadata.Label = make(metadata.Label)
	}
	businessMetadata.Label.SetBusinessID(conf.ApplicationID)
	if err := lgc.DeleteHostBusinessAttributes(ctx, conf.HostIDs, &businessMetadata); err != nil {
		blog.Errorf("move host to resource pool, delete host bussiness private, err: %v, input:%+v,rid:%s", err, conf, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}
	return nil, nil
}

// notExistAppModuleHost get hostID in the module that does not exist
// 获取不在moduleID中的hostID
func (lgc *Logics) notExistAppModuleHost(ctx context.Context, appID, moduleID int64, hostIDArr []int64) ([]int64, error) {
	hostModuleInput := &metadata.HostModuleRelationRequest{
		ApplicationID: appID,
		ModuleIDArr:   []int64{moduleID},
		HostIDArr:     hostIDArr,
	}

	hmResult, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx, lgc.header, hostModuleInput)
	if err != nil {
		blog.Errorf("existAppModule, GetHostModuleRelation http do error, err: %v, input:%+v,rid:%v", err, hostModuleInput, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hmResult.Result {
		blog.Errorf("existAppModule, GetHostModuleRelation http reply error, result: %#v, input:%+v,rid:%v", hmResult, hostModuleInput, lgc.rid)
		return nil, lgc.ccErr.New(hmResult.Code, hmResult.ErrMsg)
	}
	hostIDMap := make(map[int64]bool, 0)
	for _, row := range hmResult.Data.Info {
		hostIDMap[row.HostID] = true
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
func (lgc *Logics) AssignHostToApp(ctx context.Context, conf *metadata.DefaultModuleHostConfigParams) ([]metadata.ExceptionResult, error) {

	cond := hutil.NewOperation().WithAppID(conf.ApplicationID).Data()
	fields := fmt.Sprintf("%s,%s", common.BKOwnerIDField, common.BKAppNameField)
	appInfo, err := lgc.GetAppDetails(ctx, fields, cond)
	if err != nil {
		blog.Errorf("assign host to app failed, err: %v,input:%+v,rid:%s", err, conf, lgc.rid)
		return nil, err
	}
	if 0 == len(appInfo) {
		blog.Errorf("assign host to app error, not foud app appID: %d,input:%+v,rid:%s", conf.ApplicationID, conf, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommNotFound)
	}

	ownerAppID, err := lgc.GetDefaultAppID(ctx)
	if err != nil {
		blog.Errorf("assign host to app, but get default appid failed, err: %v,input:%+v,rid:%s", err, conf, lgc.rid)
		return nil, err
	}
	if 0 == conf.ApplicationID {
		return nil, lgc.ccErr.Errorf(common.CCErrHostGetResourceFail, "not found")
	}
	if ownerAppID == conf.ApplicationID {
		return nil, nil
	}

	conds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(ownerAppID)
	ownerModuleID, err := lgc.GetResourcePoolModuleID(ctx, conds.MapStr())
	if err != nil {
		blog.Errorf("assign host to app, but get module id failed, err: %v,input:%+v,rid:%s", err, conds.MapStr(), lgc.rid)
		return nil, err
	}
	if 0 == ownerModuleID {
		blog.Errorf("assign host to app, but get module id failed, err: %v,input:%+v,rid:%s", err, conds.MapStr(), lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrHostModuleNotExist, common.DefaultResModuleName)
	}
	errHostID, err := lgc.notExistAppModuleHost(ctx, ownerAppID, ownerModuleID, conf.HostIDs)
	if err != nil {
		blog.Errorf("move host to resource pool, notExistAppModuleHost error, err: %v, input:%+v, rid:%s", err, conf, lgc.rid)
		return nil, err
	}
	if len(errHostID) > 0 {
		errHostIP := lgc.convertHostIDToHostIP(ctx, errHostID)
		blog.Errorf("move host to resource pool, notExistAppModuleHost error, has host not belong to idle module , input:%+v, rid:%s", conf, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrHostNotBelongIDLEModuleErr, strings.Join(errHostIP, ","))
	}

	mConds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(conf.ApplicationID)
	moduleID, err := lgc.GetResourcePoolModuleID(ctx, mConds.MapStr())
	if err != nil {
		blog.Errorf("assign host to app, but get module id failed, err: %v,input:%+v,params:%+v,rid:%s", err, conf, mConds.MapStr(), lgc.rid)
		return nil, err
	}
	if moduleID == 0 {
		blog.Errorf("assign host to app, but get module id failed, %s not found,input:%+v,params:%+v,rid:%s", common.DefaultResModuleName, conf, mConds.MapStr(), lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrHostModuleNotExist, common.DefaultResModuleName)
	}

	assignParams := &metadata.TransferHostsCrossBusinessRequest{
		SrcApplicationID: ownerAppID,
		DstApplicationID: conf.ApplicationID,
		HostIDArr:        conf.HostIDs,
		DstModuleIDArr:   []int64{moduleID},
	}

	audit := lgc.NewHostModuleLog(conf.HostIDs)
	if err := audit.WithPrevious(ctx); err != nil {
		blog.Warnf("WithPrevious failed, err: %+v, rid: %s", err, lgc.rid)
	}

	result, err := lgc.CoreAPI.CoreService().Host().TransferToAnotherBusiness(ctx, lgc.header, assignParams) //.AssignHostToApp(ctx, srvData.header, params)
	if err != nil {
		blog.Errorf("assign host to app, but assign to app http do error. err: %v, input:%+v,param:%+v,rid:%s", err, conf, assignParams, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrHostEditRelationPoolFail)
	}
	if !result.Result {
		blog.Errorf("assign host to app, but assign to app http response error. result:%#v, input:%+v, param:%+v, rid:%s", result, conf, assignParams, lgc.rid)
		return result.Data, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	if err := audit.SaveAudit(ctx, conf.ApplicationID, lgc.user, "assign host to app"); err != nil {
		blog.Errorf("assign host to app, but save audit failed, err: %v, rid:%s", err, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}

	return nil, nil
}

// convertHostIDToHostIP  该方法为专用方法。出现任何错误都会被忽略。
// 尝试将主机ID转换为内网IP，如果转换中出现问题返回主机ID。
func (lgc *Logics) convertHostIDToHostIP(ctx context.Context, hostIDArr []int64) []string {

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
		ReadInstance(ctx, lgc.header, common.BKInnerObjIDHost, input)
	if err != nil {
		blog.Warnf("convertHostIDToHostIP http do error. err:%s, input:%#v, rid:%s", err.Error(), input, lgc.rid)
	}
	if !result.Result {
		blog.Warnf("convertHostIDToHostIP http response error. result:%#v, input:%#v, rid:%s", result, input, lgc.rid)
	}
	for _, host := range result.Data.Info {
		hostID, err := host.Int64(common.BKHostIDField)
		if err != nil {
			// can't not foud host id , skip
			blog.Warnf("convertHostIDToHostIP convert host id to int64 error. err:%s, host:%#v, input:%#v, rid:%s", err.Error(), host, input, lgc.rid)
			continue
		}
		innerIP, err := host.String(common.BKHostInnerIPField)
		if err != nil {
			// can't not foud host inner ip , skip
			blog.Warnf("convertHostIDToHostIP convert host inner ip to string error. err:%s, host:%#v, input:%#v, rid:%s", err.Error(), host, input, lgc.rid)
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
