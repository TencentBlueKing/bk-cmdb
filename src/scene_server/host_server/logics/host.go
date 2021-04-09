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
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) GetHostAttributes(kit *rest.Kit, bizMetaOpt mapstr.MapStr) ([]metadata.Attribute, error) {
	searchOp := mapstr.MapStr{
		common.BKObjIDField: common.BKInnerObjIDHost,
	}
	if bizMetaOpt != nil {
		searchOp.Merge(bizMetaOpt)
	}
	query := &metadata.QueryCondition{
		Condition: searchOp,
	}
	result, err := lgc.CoreAPI.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDHost, query)
	if err != nil {
		blog.Errorf("GetHostAttributes http do error, err:%s, input:%+v, rid:%s", err.Error(), query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostAttributes http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func (lgc *Logics) GetHostInstanceDetails(kit *rest.Kit, hostID int64) (map[string]interface{}, string, errors.CCError) {
	// get host details, pre data
	result, err := lgc.CoreAPI.CoreService().Host().GetHostByID(kit.Ctx, kit.Header, hostID)
	if err != nil {
		blog.Errorf("GetHostInstanceDetails http do error, err:%s, input:%+v, rid:%s", err.Error(), hostID, kit.Rid)
		return nil, "", kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostInstanceDetails http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, hostID, kit.Rid)
		return nil, "", kit.CCError.New(result.Code, result.ErrMsg)
	}

	hostInfo := result.Data
	if len(hostInfo) == 0 {
		return nil, "", nil
	}
	ip, ok := hostInfo[common.BKHostInnerIPField].(string)
	if !ok {
		blog.Errorf("GetHostInstanceDetails http response format error,convert bk_biz_id to int error, inst:%#v  input:%#v, rid:%s", hostInfo, hostID, kit.Rid)
		return nil, "", kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostInnerIPField, "string", "not string")

	}
	return hostInfo, ip, nil
}

// GetHostRelations get hosts owned set, module info, where hosts must match condition specify by cond.
func (lgc *Logics) GetHostRelations(kit *rest.Kit, input metadata.HostModuleRelationRequest) ([]metadata.ModuleHost, errors.CCError) {

	result, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("GetConfigByCond http do error, err:%s, input:%+v, rid:%s", err.Error(), input, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetConfigByCond http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, input, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

// EnterIP 将机器导入到指定模块或者空闲模块， 已经存在机器，不操作
func (lgc *Logics) EnterIP(kit *rest.Kit, appID, moduleID int64, ip string, cloudID int64, host map[string]interface{}, isIncrement bool) errors.CCError {

	isExist, err := lgc.IsPlatExist(kit, mapstr.MapStr{common.BKCloudIDField: cloudID})
	if nil != err {
		return err
	}
	if !isExist {
		return kit.CCError.Errorf(common.CCErrTopoCloudNotFound)
	}
	ipArr := strings.Split(ip, ",")
	conds := mapstr.MapStr{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBIN: ipArr,
		},
		common.BKCloudIDField: cloudID,
	}
	hostList, err := lgc.GetHostInfoByConds(kit, conds)
	if nil != err {
		return err
	}

	hostID := int64(0)
	if len(hostList) == 0 {
		//host not exist, add host
		host[common.BKHostInnerIPField] = ip
		host[common.BKCloudIDField] = cloudID
		host["import_from"] = common.HostAddMethodAgent
		defaultFields, hasErr := lgc.getHostFields(kit)
		if nil != hasErr {
			return hasErr
		}
		//补充未填写字段的默认值
		for _, field := range defaultFields {
			_, ok := host[field.PropertyID]
			if !ok {
				if true == util.IsStrProperty(field.PropertyType) {
					host[field.PropertyID] = ""
				} else {
					host[field.PropertyID] = nil
				}
			}
		}

		result, err := lgc.CoreAPI.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDHost, &metadata.CreateModelInstance{Data: host})
		if err != nil {
			blog.Errorf("EnterIP http do error, err:%s, input:%+v, rid:%s", err.Error(), host, kit.Rid)
			return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("EnterIP http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, host, kit.Rid)
			return kit.CCError.New(result.Code, result.ErrMsg)
		}

		// add audit log for create host.
		audit := auditlog.NewHostAudit(lgc.CoreAPI.CoreService())
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, hostID, appID, "", nil)
		if err != nil {
			blog.Errorf("generate audit log failed after create host, hostID: %d, appID: %d, err: %v, rid: %s",
				hostID, appID, err, kit.Rid)
			return err
		}

		// save audit log.
		if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
			blog.Errorf("save audit log failed after create host, hostID: %d, appID: %d,err: %v, rid: %s", hostID,
				appID, err, kit.Rid)
			return err
		}

		hostID = int64(result.Data.Created.ID)
	} else if false == isIncrement {
		// Not an additional relationship model
		return nil
	} else {

		hostID, err = util.GetInt64ByInterface(hostList[0][common.BKHostIDField])
		if err != nil {
			blog.Errorf("EnterIP  get hostID error, err:%s,inst:%+v,input:%+v, rid:%s", err.Error(), hostList[0], host, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", err.Error()) // "查询主机信息失败"
		}

		bl, hasErr := lgc.IsHostExistInApp(kit, appID, hostID)
		if nil != hasErr {
			return hasErr

		}
		if false == bl {
			blog.Errorf("Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, rid:%s", appID, hostID, kit.Rid)
			return kit.CCError.Errorf(common.CCErrHostNotINAPPFail, hostID)
		}

	}

	hmAudit := auditlog.NewHostModuleLog(lgc.CoreAPI.CoreService(), []int64{hostID})
	if err := hmAudit.WithPrevious(kit); err != nil {
		return err
	}

	params := &metadata.HostsModuleRelation{
		ApplicationID: appID,
		HostID:        []int64{hostID},
		ModuleID:      []int64{moduleID},
		IsIncrement:   isIncrement,
	}
	hmResult, ccErr := lgc.CoreAPI.CoreService().Host().TransferToNormalModule(kit.Ctx, kit.Header, params)
	if ccErr != nil {
		blog.Errorf("Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, err:%s, rid:%s", appID, hostID, ccErr.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hmResult.Result {
		blog.Errorf("transfer host to normal module failed, error params:{appID:%d, hostID:%d}, result:%#v, rid:%s", appID, hostID, hmResult, kit.Rid)
		if len(hmResult.Data) > 0 {
			return kit.CCError.New(int(hmResult.Data[0].Code), hmResult.Data[0].Message)
		}
		return kit.CCError.New(hmResult.Code, hmResult.ErrMsg)
	}

	if err := hmAudit.SaveAudit(kit); err != nil {
		return err
	}
	return nil
}

func (lgc *Logics) GetHostInfoByConds(kit *rest.Kit, cond map[string]interface{}) ([]mapstr.MapStr, errors.CCErrorCoder) {
	query := &metadata.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKHostIDField,
	}

	result, err := lgc.CoreAPI.CoreService().Host().GetHosts(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("GetHostInfoByConds GetHosts http do error, err:%s, input:%+v,rid:%s", err.Error(), query, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if err := result.CCError(); err != nil {
		blog.Errorf("GetHostInfoByConds GetHosts http response error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, kit.Rid)
		return nil, err
	}

	return result.Data.Info, nil
}

// SearchHostInfo search host info by QueryCondition
func (lgc *Logics) SearchHostInfo(kit *rest.Kit, cond metadata.QueryCondition) ([]mapstr.MapStr, errors.CCErrorCoder) {
	query := &metadata.QueryInput{
		Condition: cond.Condition,
		Fields:    strings.Join(cond.Fields, ","),
		Start:     cond.Page.Start,
		Limit:     cond.Page.Limit,
		Sort:      cond.Page.Sort,
	}

	result, err := lgc.CoreAPI.CoreService().Host().GetHosts(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("GetHostInfoByConds GetHosts http do error, err:%s, input:%+v,rid:%s", err.Error(), query, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if err := result.CCError(); err != nil {
		blog.Errorf("GetHostInfoByConds GetHosts http response error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, kit.Rid)
		return nil, err
	}

	return result.Data.Info, nil
}

// HostSearch search host by multiple condition
const (
	SplitFlag      = "##"
	TopoSetName    = "TopSetName"
	TopoModuleName = "TopModuleName"
)

// GetHostIDByCond query hostIDs by condition base on cc_ModuleHostConfig
// available condition fields are bk_supplier_account, bk_biz_id, bk_host_id, bk_module_id, bk_set_id
func (lgc *Logics) GetHostIDByCond(kit *rest.Kit, cond metadata.HostModuleRelationRequest) ([]int64, errors.CCError) {

	cond.Fields = []string{common.BKHostIDField}
	result, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, &cond)
	if err != nil {
		blog.Errorf("GetHostIDByCond GetModulesHostConfig http do error, err:%s, input:%+v,rid:%s", err.Error(), cond, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostIDByCond GetModulesHostConfig http response error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, cond, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	hostIDs := make([]int64, 0)
	for _, val := range result.Data.Info {
		hostIDs = append(hostIDs, val.HostID)
	}

	return hostIDs, nil
}

// GetAllHostIDByCond 专用结构， page start 和limit 无效， 获取条件所有满足条件的主机
func (lgc *Logics) GetAllHostIDByCond(kit *rest.Kit, cond metadata.HostModuleRelationRequest) ([]int64, errors.CCError) {
	hostIDs := make([]int64, 0)
	cond.Page.Limit = 2000
	start := 0
	cnt := 0
	cond.Fields = []string{common.BKHostIDField}
	for {
		cond.Page.Start = start
		result, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, &cond)
		if err != nil {
			blog.Errorf("GetHostIDByCond GetModulesHostConfig http do error, err:%s, input:%+v,rid:%s", err.Error(), cond, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("GetHostIDByCond GetModulesHostConfig http response error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, cond, kit.Rid)
			return nil, kit.CCError.New(result.Code, result.ErrMsg)
		}

		for _, val := range result.Data.Info {
			hostIDs = append(hostIDs, val.HostID)
		}
		// 当总数大于现在的总数，使用当前返回值的总是为新的总数值
		if cnt < int(result.Data.Count) {
			// 获取条件的数据总数
			cnt = int(result.Data.Count)
		}
		start += cond.Page.Limit
		if start >= cnt {
			break
		}
	}

	return hostIDs, nil
}

// GetHostModuleRelation  query host and module relation,
// condition key use appID, moduleID,setID,HostID
func (lgc *Logics) GetHostModuleRelation(kit *rest.Kit, cond *metadata.HostModuleRelationRequest) (*metadata.
	HostConfigData, errors.CCErrorCoder) {

	if cond.Empty() {
		return nil, kit.CCError.CCError(common.CCErrCommHTTPBodyEmpty)
	}

	if cond.Page.IsIllegal() {
		return nil, kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded)
	}

	if len(cond.SetIDArr) > 200 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "bk_set_ids", 200)
	}

	if len(cond.ModuleIDArr) > 500 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "bk_module_ids", 500)
	}

	if len(cond.HostIDArr) > 500 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "bk_host_ids", 500)
	}

	result, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("GetHostModuleRelation http do error, err:%s, input:%+v, rid:%s", err.Error(), cond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if retErr := result.CCError(); retErr != nil {
		blog.Errorf("GetHostModuleRelation http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, cond, kit.Rid)
		return nil, retErr
	}

	return &result.Data, nil
}

// TransferHostAcrossBusiness  Transfer host across business, can only transfer between resource set modules
// delete old business  host and module relation
func (lgc *Logics) TransferHostAcrossBusiness(kit *rest.Kit, srcBizID, dstAppID int64, hostID []int64, moduleID int64) errors.CCError {
	// get both biz's resource set's modules
	query := &metadata.QueryCondition{
		Fields: []string{common.BKModuleIDField, common.BKAppIDField},
		Condition: map[string]interface{}{
			common.BKAppIDField:   map[string]interface{}{common.BKDBIN: []int64{srcBizID, dstAppID}},
			common.BKDefaultField: map[string]interface{}{common.BKDBNE: common.NormalModuleFlag},
		},
	}

	moduleRes, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("transfer host across business, get modules failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err := moduleRes.CCError(); err != nil {
		blog.Errorf("transfer host across business, get modules failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return err
	}

	// valid if dest module is dest biz's resource set's module, get src biz module ids
	moduleIDArr := make([]int64, 0)
	isDestModuleValid := false
	for _, module := range moduleRes.Data.Info {
		modID, err := module.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("transfer host across business, get module(%s) id failed, err: %s, rid: %s", module, err.Error(), kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
		}

		bizID, err := module.Int64(common.BKAppIDField)
		if err != nil {
			blog.ErrorJSON("transfer host across business, get module(%s) biz id failed, err: %s, rid: %s", module, err.Error(), kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
		}

		if modID == moduleID && bizID == dstAppID {
			isDestModuleValid = true
		}

		if bizID == srcBizID {
			moduleIDArr = append(moduleIDArr, modID)
		}
	}

	if !isDestModuleValid {
		blog.Errorf("transfer host across business, dest module(%d) does not belong to the resource set of the dest biz, rid: %s", moduleID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	// valid if hosts are in the resource set's modules of the src biz
	notExistHostIDs, err := lgc.notExistAppModuleHost(kit, srcBizID, moduleIDArr, hostID)
	if err != nil {
		blog.Errorf("TransferHostAcrossBusiness IsHostExistInApp err:%s,input:{appID:%d,hostID:%d},rid:%s", err.Error(), srcBizID, hostID, kit.Rid)
		return err
	}

	if len(notExistHostIDs) > 0 {
		notExistHostIP := lgc.convertHostIDToHostIP(kit, notExistHostIDs)
		blog.Errorf("transfer host across business, has host not belong to idle module , host ids: %+v, rid: %s", notExistHostIDs, kit.Rid)
		return kit.CCError.Errorf(common.CCErrHostModuleConfigNotMatch, util.PrettyIPStr(notExistHostIP))
	}

	// do transfer and save audit log
	audit := auditlog.NewHostModuleLog(lgc.CoreAPI.CoreService(), hostID)
	if err := audit.WithPrevious(kit); err != nil {
		blog.Errorf("TransferHostAcrossBusiness, get prev module host config failed, err: %v,hostID:%d,oldbizID:%d,appID:%d, moduleID:%#v,rid:%s", err, hostID, srcBizID, dstAppID, moduleID, kit.Rid)
		return err
	}

	conf := &metadata.TransferHostsCrossBusinessRequest{SrcApplicationID: srcBizID, HostIDArr: hostID, DstApplicationID: dstAppID, DstModuleIDArr: []int64{moduleID}}
	delRet, doErr := lgc.CoreAPI.CoreService().Host().TransferToAnotherBusiness(kit.Ctx, kit.Header, conf)
	if doErr != nil {
		blog.Errorf("TransferHostAcrossBusiness http do error, err:%s, input:%+v, rid:%s", doErr.Error(), conf, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err := delRet.CCError(); err != nil {
		blog.Errorf("TransferHostAcrossBusiness http response error, err code:%d, err msg:%s, input:%#v, rid:%s", delRet.Code, delRet.ErrMsg, conf, kit.Rid)
		return err
	}

	if err := audit.SaveAudit(kit); err != nil {
		blog.Errorf("TransferHostAcrossBusiness, get prev module host config failed, err: %v,hostID:%d,oldbizID:%d,appID:%d, moduleID:%#v,rid:%s", err, hostID, srcBizID, dstAppID, moduleID, kit.Rid)
		return err
	}

	return nil
}

// DeleteHostFromBusiness  delete host from business,
func (lgc *Logics) DeleteHostFromBusiness(kit *rest.Kit, bizID int64, hostIDArr []int64) ([]metadata.ExceptionResult, errors.CCError) {
	// ready audit log of delete host.
	audit := auditlog.NewHostAudit(lgc.CoreAPI.CoreService())
	logContentMap := make(map[int64]*metadata.AuditLog, 0)
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	for _, hostID := range hostIDArr {
		var err error
		logContentMap[hostID], err = audit.GenerateAuditLog(generateAuditParameter, hostID, bizID, "", nil)
		if err != nil {
			blog.Errorf("generate host audit log failed before delete host, hostID: %d, bizID: %d, err: %v, rid: %s", hostID, bizID, err, kit.Rid)
			return nil, err
		}
	}

	// to delete host.
	input := &metadata.DeleteHostRequest{
		ApplicationID: bizID,
		HostIDArr:     hostIDArr,
	}
	result, err := lgc.CoreAPI.CoreService().Host().DeleteHostFromSystem(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("TransferHostAcrossBusiness DeleteHost error, err: %v,hostID:%#v,appID:%d,rid:%s", err, hostIDArr, bizID, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("TransferHostAcrossBusiness DeleteHost failed, err: %v,hostID:%#v,appID:%d,rid:%s", err, hostIDArr, bizID, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	// to save audit log.
	logContents := make([]metadata.AuditLog, len(logContentMap))
	index := 0
	for _, item := range logContentMap {
		logContents[index] = *item
		index++
	}

	if len(logContents) > 0 {
		if err := audit.SaveAuditLog(kit, logContents...); err != nil {
			blog.ErrorJSON("delete host in batch, but add host audit log failed, err: %s, rid: %s",
				err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
		}

	}
	return nil, nil
}

// CloneHostProperty clone host info and host and module relation in same application
func (lgc *Logics) CloneHostProperty(kit *rest.Kit, appID int64, srcHostID int64, dstHostID int64) errors.CCErrorCoder {

	// check if source host and destination host both belong to biz
	relReq := &metadata.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{appID},
		HostIDArr:        []int64{srcHostID, dstHostID},
	}

	relRsp, relErr := lgc.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(kit.Ctx, kit.Header, relReq)
	if relErr != nil {
		blog.ErrorJSON("get host ids in biz failed, err: %s, req: %s, rid: %s", relErr, relReq, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if err := relRsp.CCError(); err != nil {
		blog.ErrorJSON("get host ids in biz failed, err: %s, req: %s, rid: %s", err, relReq, kit.Rid)
		return err
	}

	isSrcHostInBiz, isDstHostInBiz := false, false
	for _, hostID := range relRsp.Data.IDArr {
		if hostID == srcHostID {
			isSrcHostInBiz = true
		}
		if hostID == dstHostID {
			isDstHostInBiz = true
		}
	}

	if !isSrcHostInBiz {
		blog.Errorf("Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, rid:%s", appID, srcHostID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrHostNotINAPPFail, srcHostID)
	}

	if !isDstHostInBiz {
		blog.Errorf("Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, rid:%s", appID, dstHostID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrHostNotINAPPFail, dstHostID)
	}

	attrCond := make(map[string]interface{})
	util.AddModelBizIDCondition(attrCond, appID)
	hostAttributes, attrErr := lgc.GetHostAttributes(kit, attrCond)
	if attrErr != nil {
		blog.Errorf("get host attributes failed, err: %v, biz id: %d, rid: %s", attrErr, appID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoObjectAttributeSelectFailed)
	}

	uniqueReq := metadata.QueryCondition{
		Condition: map[string]interface{}{common.BKObjIDField: common.BKInnerObjIDHost},
	}
	uniqueRsp, uniqueErr := lgc.CoreAPI.CoreService().Model().ReadModelAttrUnique(kit.Ctx, kit.Header, uniqueReq)
	if uniqueErr != nil {
		blog.ErrorJSON("get host ids in biz failed, err: %s, req: %s, rid: %s", relErr, relReq, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	attrIDUniqueMap := make(map[uint64]struct{})
	for _, unique := range uniqueRsp.Data.Info {
		for _, key := range unique.Keys {
			attrIDUniqueMap[key.ID] = struct{}{}
		}
	}

	hostFields := make([]string, 0)
	for _, attr := range hostAttributes {
		if !attr.IsEditable {
			continue
		}
		if _, exists := attrIDUniqueMap[uint64(attr.ID)]; exists {
			continue
		}
		hostFields = append(hostFields, attr.PropertyID)
	}

	if len(hostFields) == 0 {
		blog.Infof("there are no host fields can be cloned, skip, rid: %s", kit.Rid)
		return nil
	}

	cond := metadata.QueryCondition{
		Fields: hostFields,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKHostIDField: srcHostID,
		},
	}

	hostInfoArr, err := lgc.SearchHostInfo(kit, cond)
	if err != nil {
		blog.ErrorJSON("search hosts failed, err: %s, input: %s, rid: %s", err, cond, kit.Rid)
		return err
	}

	if len(hostInfoArr) == 0 {
		blog.Errorf("host not found. hostID:%s, rid:%s", srcHostID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrHostNotFound)
	}
	srcHostInfo := hostInfoArr[0]

	delete(srcHostInfo, common.BKHostIDField)
	delete(srcHostInfo, common.CreateTimeField)

	// generate audit log
	audit := auditlog.NewHostAudit(lgc.CoreAPI.CoreService())
	auditParam := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(srcHostInfo)
	auditLog, auditErr := audit.GenerateAuditLog(auditParam, dstHostID, appID, "", nil)
	if auditErr != nil {
		blog.Errorf("generate audit log failed, err: %v, host id: %d, rid: %s", err, dstHostID, kit.Rid)
		return kit.CCError.CCError(common.CCErrAuditTakeSnapshotFailed)
	}

	input := &metadata.UpdateOption{
		Data: srcHostInfo,
		Condition: mapstr.MapStr{
			common.BKHostIDField: dstHostID,
		},
	}
	result, doErr := lgc.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDHost, input)
	if doErr != nil {
		blog.ErrorJSON("CloneHostProperty UpdateInstance error. err: %s,condition:%s,rid:%s", doErr, input, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if err := result.CCError(); err != nil {
		blog.ErrorJSON("CloneHostProperty UpdateInstance  replay error. err: %s,condition:%s,rid:%s", err, input, kit.Rid)
		return err
	}

	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("save audit log failed, err: %v, host id: %d, rid: %s", err, dstHostID, kit.Rid)
		return kit.CCError.CCError(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

// IPCloudToHost get host id by ip and cloud
func (lgc *Logics) IPCloudToHost(kit *rest.Kit, ip string, cloudID int64) (HostMap mapstr.MapStr, hostID int64, err errors.CCErrorCoder) {
	// FIXME there must be a better ip to hostID solution
	ipArr := strings.Split(ip, ",")
	condition := mapstr.MapStr{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBAll: ipArr,
		},
		common.BKCloudIDField: cloudID,
	}

	hostInfoArr, err := lgc.GetHostInfoByConds(kit, condition)
	if err != nil {
		blog.ErrorJSON("IPCloudToHost GetHostInfoByConds error. err:%s, conditon:%s, rid:%s", err.Error(), condition, kit.Rid)
		return nil, 0, err
	}
	if len(hostInfoArr) == 0 {
		return nil, 0, nil
	}

	hostID, convErr := hostInfoArr[0].Int64(common.BKHostIDField)
	if nil != convErr {
		blog.ErrorJSON("IPCloudToHost bk_host_id field not found hostMap:%s ip:%s, cloudID:%s,rid:%s", hostInfoArr, ip, cloudID, kit.Rid)
		return nil, 0, kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", convErr.Error())
	}

	return hostInfoArr[0], hostID, nil
}

// ArrangeHostDetailAndTopology arrange host's detail and it's topology node's info along with it.
func (lgc *Logics) ArrangeHostDetailAndTopology(kit *rest.Kit, withBiz bool, hosts []map[string]interface{}) (
	[]*metadata.HostDetailWithTopo, error) {

	// get mainline topology rank data, it's the order to arrange the host's topology data.
	rankMap, reverseRankMap, rank, err := lgc.getTopologyRank(kit)
	if err != nil {
		return nil, err
	}

	// search all hosts' host module relations
	hostIDs := make([]int64, 0)
	for _, host := range hosts {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.ErrorJSON("got invalid bk_host_id field in host: %s, rid: %s", host, kit.Rid)
			return nil, err
		}
		hostIDs = append(hostIDs, hostID)
	}
	relationCond := metadata.HostModuleRelationRequest{
		HostIDArr: hostIDs,
	}
	relations, err := lgc.GetHostRelations(kit, relationCond)
	if nil != err {
		blog.ErrorJSON("read host module relation error: %s, input: %s, rid: %s", err, hosts, kit.Rid)
		return nil, err
	}

	bizList := make([]int64, 0)
	moduleList := make([]int64, 0)
	setList := make([]int64, 0)
	hostModule := make(map[int64][]int64)
	for _, one := range relations {
		bizList = append(bizList, one.AppID)
		setList = append(setList, one.SetID)
		moduleList = append(moduleList, one.ModuleID)
		hostModule[one.HostID] = append(hostModule[one.HostID], one.ModuleID)
	}

	// get all the inner object's details info
	bizDetails, setDetails, moduleDetails, err := lgc.getInnerObjectDetails(kit, withBiz, bizList, moduleList, setList)
	if err != nil {
		return nil, err
	}

	// now we get all the custom object's instances with set's parent instance id
	// from low level to the top business level.
	customObjInstMap := make(map[string]map[int64]mapstr.MapStr)
	reversedRank := util.ReverseArrayString(rank)
	var parentsInst []interface{}
loop:
	for _, one := range reversedRank {
		switch one {
		case common.BKInnerObjIDApp:
			break loop
		case common.BKInnerObjIDHost:
			continue
		case common.BKInnerObjIDModule:
			continue
		case common.BKInnerObjIDSet:
			if rankMap[common.BKInnerObjIDSet] == common.BKInnerObjIDApp {
				break loop
			}

			parentsInst = make([]interface{}, 0)
			for _, set := range setDetails {
				if fmt.Sprintf("%v", set[common.BKDefaultField]) != "0" {
					// this is a inner set, do not have custom obj parent instance
					continue
				}
				parentsInst = append(parentsInst, set[common.BKParentIDField])
			}

			continue
		default:
			if len(parentsInst) == 0 {
				// when the set is inner set, which default field value is > 0;
				break loop
			}
			// get custom level instances details with parent instance id list.
			customInst, err := lgc.getCustomObjectInstances(kit, one, parentsInst)
			if err != nil {
				return nil, err
			}

			// reset parent instances
			parentsInst = make([]interface{}, 0)
			// save the custom instances details
			for _, inst := range customInst {
				if _, exists := customObjInstMap[one]; !exists {
					customObjInstMap[one] = make(map[int64]mapstr.MapStr)
				}

				instID, err := util.GetInt64ByInterface(inst[common.BKInstIDField])
				if err != nil {
					blog.Errorf("get inst id from inst: %v failed, err: %v, rid: %s", inst, err, kit.Rid)
					return nil, err
				}

				// save the instances data with object and it's instances
				customObjInstMap[one][instID] = inst
				// update parent instances
				parentsInst = append(parentsInst, inst[common.BKParentIDField])
			}

			if rankMap[one] == common.BKInnerObjIDApp {
				break loop
			}
		}
	}

	// now, we have already get all the data we need, it's time to arrange the data.
	bizMap := make(map[int64]mapstr.MapStr)
	for _, biz := range bizDetails {
		bizID, err := util.GetInt64ByInterface(biz[common.BKAppIDField])
		if err != nil {
			blog.Errorf("get biz id failed, biz: %v, err: %v, rid: %s", biz, err, kit.Rid)
			return nil, err
		}
		bizMap[bizID] = biz
	}
	setMap := make(map[int64]mapstr.MapStr)
	for _, set := range setDetails {
		setID, err := util.GetInt64ByInterface(set[common.BKSetIDField])
		if err != nil {
			blog.Errorf("get set id failed, set: %v, err: %v, rid: %s", set, err, kit.Rid)
			return nil, err
		}
		setMap[setID] = set
	}

	moduleMap := make(map[int64]mapstr.MapStr)
	for _, mod := range moduleDetails {
		modID, err := util.GetInt64ByInterface(mod[common.BKModuleIDField])
		if err != nil {
			blog.Errorf("get module id failed, module: %v, err: %v, rid: %s", mod, err, kit.Rid)
			return nil, err
		}
		moduleMap[modID] = mod
	}

	// reset the rank from biz to module with host one by one.
	rank = util.ReverseArrayString(rank)
	topo := make([]*metadata.HostDetailWithTopo, len(hosts))
	for idx, one := range hosts {
		hostID, err := util.GetInt64ByInterface(one[common.BKHostIDField])
		if err != nil {
			blog.ErrorJSON("got invalid bk_host_id field in host: %s, rid: %s", one, kit.Rid)
			return nil, err
		}

		topo[idx] = &metadata.HostDetailWithTopo{Host: one}
		modules, exists := hostModule[hostID]
		if !exists {
			blog.Errorf("can not find modules host %d belongs to, host: %v, skip, rid: %s", hostID, one, kit.Rid)
			continue
		}

		// one host can only belongs to one business, so the resource in the tree is all belongs to a same business.
		children, err := lgc.arrangeParentTree(kit.Rid, rank, rankMap, reverseRankMap, withBiz, bizMap, setMap,
			moduleMap, modules, customObjInstMap)
		if err != nil {
			return nil, err
		}

		topo[idx].Topo = children

	}

	return topo, nil
}

// arrangeParentTree is to arrange the host's topology tree with the modules it belongs to.
// so all these generated topology tree's nodes must belongs to a same business.
// the tree is constructed with two steps:
// step1: rearrange the object's instances with map[parent_id]child_detail, so that we can know
// which instances have how many children and what exactly its child is.
// step2: arrange the tree from the top object to the bottom module with step 1's map.
func (lgc *Logics) arrangeParentTree(
	rid string,
	rank []string,
	rankMap map[string]string,
	reverseRankMap map[string]string,
	withBiz bool,
	bizsMap, setsMap, modulesMap map[int64]mapstr.MapStr,
	relateModules []int64,
	objInstMap map[string]map[int64]mapstr.MapStr) ([]*metadata.HostTopoNode, error) {

	if len(rank) < 4 {
		// min rank is biz,set,module,host
		return nil, fmt.Errorf("invalid rank, should at least have 4 level, detail: %v", rank)
	}

	// object -> instance_id -> instance_children_list
	parentToChildrenMap := make(map[string]map[int64][]mapstr.MapStr)
	setModules := make(map[int64][]mapstr.MapStr)

	// we separate inner module(default != 0), cause, it's topology is different, it does not have
	// custom object level.
	innerSetModules := make(map[int64][]mapstr.MapStr)
	for _, m := range relateModules {
		mod, exist := modulesMap[m]
		if !exist {
			blog.Errorf("can not find module: %d detail, skip, rid: %s", m, rid)
			continue
		}

		isInner := false
		if fmt.Sprintf("%v", mod[common.BKDefaultField]) != "0" {

			isInner = true
		}

		setID, err := util.GetInt64ByInterface(mod[common.BKSetIDField])
		if err != nil {
			blog.Errorf("get set id from module failed, module: %v, err: %v, rid: %s", m, err, rid)
			return nil, err
		}

		if isInner {
			// this is a inner module, it will be handled with special logic.
			innerSetModules[setID] = append(innerSetModules[setID], mod)
			continue
		}

		setModules[setID] = append(setModules[setID], mod)
	}
	parentToChildrenMap[common.BKInnerObjIDSet] = setModules

	// start from module object
	parent := rankMap[common.BKInnerObjIDModule]

	rootBizInstances := make(map[int64][]mapstr.MapStr)
loop:
	// for loop model from bottom module to business
	for {
		// get next object for prepare
		next := rankMap[parent]

		switch parent {
		case common.BKInnerObjIDSet:

			sets, exists := parentToChildrenMap[common.BKInnerObjIDSet]
			if !exists {
				blog.ErrorJSON("get set list failed, all object detail: %s, err: %s, rid: %s", parentToChildrenMap, rid)
				return nil, fmt.Errorf("can not find all set instances details")
			}

			setParents := make(map[int64][]mapstr.MapStr)
			for setID := range sets {
				set, exists := setsMap[setID]
				if !exists {
					blog.Errorf("can not find set detail with set id: %d, rid: %s", setID, rid)
					return nil, fmt.Errorf("can not find set instance detail with id")
				}
				pid, err := util.GetInt64ByInterface(set[common.BKParentIDField])
				if err != nil {
					blog.Errorf("get set id failed, set: %v, err: %v, rid: %s", set, err, rid)
					return nil, err
				}
				setParents[pid] = append(setParents[pid], set)
			}
			parentToChildrenMap[next] = setParents

		default:
			// get this parent object's all instances found from previous cycle.
			current, exists := parentToChildrenMap[parent]
			if !exists {
				blog.ErrorJSON("can not get %s list instances, all object detail: %s, err: %s, rid: %s", parent,
					parentToChildrenMap, rid)
				return nil, fmt.Errorf("can not find %s instance details", parent)
			}

			customParents := make(map[int64][]mapstr.MapStr)
			for instID := range current {
				pInstances, exists := objInstMap[parent]
				if !exists {
					blog.Errorf("can not get %s instances map, rid: %s", parent, rid)
					return nil, fmt.Errorf("can not find %s instance details", parent)
				}
				one, exists := pInstances[instID]
				if !exists {
					blog.Errorf("can not find set detail with set id: %d, rid: %s", instID, rid)
					return nil, fmt.Errorf("can not find set detail with id: %d", instID)
				}

				// find this one's parent instance id
				pid, err := util.GetInt64ByInterface(one[common.BKParentIDField])
				if err != nil {
					blog.Errorf("get %s inst id failed, inst: %v, err: %v, rid: %s", parent, one, err, rid)
					return nil, err
				}

				// record this instance's parent's children which is itself.
				customParents[pid] = append(customParents[pid], one)

				if next == common.BKInnerObjIDApp {
					if biz, exist := bizsMap[pid]; exist {
						rootBizInstances[pid] = []mapstr.MapStr{biz}
					}
				}
			}

			parentToChildrenMap[next] = customParents
		}

		// update parent level
		parent = next

		// check if we have already hit the topology's root
		if next == common.BKInnerObjIDApp {
			break loop
		}
	}

	// now, we format the topology
	var root string
	nodes := make([]*metadata.HostTopoNode, 0)
	if withBiz {
		// start from biz
		root = common.BKInnerObjIDApp

		for rootID := range rootBizInstances {
			node := &metadata.HostTopoNode{
				Instance: &metadata.NodeInstance{
					Object:   root,
					InstName: bizsMap[rootID][common.BKAppNameField],
					InstID:   rootID,
				},
				Children: getTopologyChildren(rid, reverseRankMap, root, rootID, parentToChildrenMap),
			}
			nodes = append(nodes, node)
		}

	} else {
		// start from the second level
		root = rank[1]
		rootInstances, exist := parentToChildrenMap[common.BKInnerObjIDApp]
		if !exist {
			blog.Errorf("can not find %s object's instances, rid: %s", root, rid)
			return nil, fmt.Errorf("can not find %s object's instances", root)
		}
		// rootInstanceMap = objInstMap[root]
		nameField := common.GetInstNameField(root)
		idField := common.GetInstIDField(root)
		for _, children := range rootInstances {

			childObj := reverseRankMap[root]
			for _, one := range children {
				childID, err := util.GetInt64ByInterface(one[common.GetInstIDField(childObj)])
				if err != nil {
					blog.Errorf("get %s instance id failed, inst: %v, err: %v, rid: %s", childObj, one, err, rid)
					return nil, err
				}
				node := &metadata.HostTopoNode{
					Instance: &metadata.NodeInstance{
						Object:   root,
						InstName: one[nameField],
						InstID:   one[idField],
					},
					Children: getTopologyChildren(rid, reverseRankMap, root, childID, parentToChildrenMap),
				}
				nodes = append(nodes, node)
			}

		}
	}

	if len(innerSetModules) != 0 {
		// these set and modules are in the same business.
		// add inner module and sets to the topology nodes.
		nodes = append(nodes, arrangeInnerModuleTree(rid, withBiz, bizsMap, innerSetModules, setsMap)...)
	}

	return nodes, nil
}

// arrange inner module's topology tree specially.
func arrangeInnerModuleTree(rid string, withBiz bool, bizsMap map[int64]mapstr.MapStr,
	innerSetModules map[int64][]mapstr.MapStr, setMap map[int64]mapstr.MapStr) []*metadata.HostTopoNode {

	nodes := make([]*metadata.HostTopoNode, 0)
	if len(innerSetModules) == 0 {
		return nodes
	}

	var oneSet mapstr.MapStr

	setNodes := make([]*metadata.HostTopoNode, 0)
	for setID, modules := range innerSetModules {
		set, exists := setMap[setID]
		if !exists {
			blog.Errorf("can not find biz id from set: %v, skip, rid: %s", set, rid)
			continue
		}

		oneSet = set
		node := &metadata.HostTopoNode{
			Instance: &metadata.NodeInstance{
				Object:   common.BKInnerObjIDSet,
				InstName: setMap[setID][common.BKSetNameField],
				InstID:   setMap[setID][common.BKSetIDField],
			},
			Children: make([]*metadata.HostTopoNode, 0),
		}

		for _, m := range modules {
			node.Children = append(node.Children, &metadata.HostTopoNode{
				Instance: &metadata.NodeInstance{
					Object:   common.BKInnerObjIDModule,
					InstName: m[common.BKModuleNameField],
					InstID:   m[common.BKModuleIDField],
				},
			})
		}

		setNodes = append(setNodes, node)
	}

	if !withBiz {
		return setNodes
	}

	if len(setNodes) == 0 {
		// seriously? it can not be happen
		blog.Errorf("got 0 set nodes, skip, rid: %s", rid)
		return setNodes
	}

	// these inner module belongs to a same biz
	bizIDField, exists := oneSet[common.BKParentIDField]
	if !exists {
		blog.Errorf("can not find biz id from set: %v, skip, rid: %s", oneSet, rid)
		return nodes
	}

	bizID, err := util.GetInt64ByInterface(bizIDField)
	if err != nil {
		blog.Errorf("can not parse biz id from set: %v, skip, rid: %s", oneSet, rid)
		return nodes
	}

	biz, exists := bizsMap[bizID]
	if !exists {
		blog.Errorf("can not find biz %d, from biz map: %v, skip, rid: %s", bizID, bizsMap, rid)
		return nodes
	}

	nodes = append(nodes, &metadata.HostTopoNode{
		Instance: &metadata.NodeInstance{
			Object:   common.BKInnerObjIDApp,
			InstName: biz[common.BKAppNameField],
			InstID:   bizID,
		},
		Children: setNodes,
	})

	return nodes
}

func getTopologyChildren(rid string, rankMap map[string]string, obj string, instID int64,
	objInstChildrenMap map[string]map[int64][]mapstr.MapStr) []*metadata.HostTopoNode {

	children := make([]*metadata.HostTopoNode, 0)
	instances, exist := objInstChildrenMap[obj][instID]
	if !exist {
		// should not happen
		return children
	}

	next, exist := rankMap[obj]
	if !exist {
		// should not happen, for safety guarantee
		return nil
	}

	// host is the bottom of the topology, do nothing
	if next == common.BKInnerObjIDHost {
		return nil
	}

	idField := common.GetInstIDField(next)
	nameField := common.GetInstNameField(next)

	for _, one := range instances {
		id, err := util.GetInt64ByInterface(one[idField])
		if err != nil {
			// should not happen, because we have already get ids from upper logical
			blog.Errorf("get instance id failed, instance: %v, rid: %s", one, rid)
			return children
		}

		child := &metadata.HostTopoNode{
			Instance: &metadata.NodeInstance{
				Object:   next,
				InstName: one[nameField],
				InstID:   id,
			},
		}

		if next != common.BKInnerObjIDModule {
			child.Children = getTopologyChildren(rid, rankMap, next, id, objInstChildrenMap)
		}

		children = append(children, child)
	}

	return children

}

// get inner object's instance details, as is biz, set, modules from cache
func (lgc *Logics) getInnerObjectDetails(kit *rest.Kit, withBiz bool, bizList, moduleList,
	setList []int64) ([]mapstr.MapStr, []mapstr.MapStr, []mapstr.MapStr, error) {

	wg := sync.WaitGroup{}
	wg.Add(3)
	var hitError error
	var bizDetails []mapstr.MapStr
	var moduleDetails []mapstr.MapStr
	var setDetails []mapstr.MapStr
	if withBiz {
		go func() {
			defer func() { wg.Done() }()
			header := kit.NewHeader()
			opt := &metadata.ListWithIDOption{
				IDs:    bizList,
				Fields: []string{common.BKAppIDField, common.BKAppNameField},
			}
			bizs, err := lgc.CoreAPI.CacheService().Cache().Topology().ListBusiness(kit.Ctx, header, opt)
			if err != nil {
				blog.Errorf("list business from cache failed, err: %v, rid: %s", err, kit.Rid)
				hitError = err
				return
			}

			list := make([]mapstr.MapStr, 0)
			if err := json.Unmarshal([]byte(bizs), &list); err != nil {
				blog.Errorf("unmarshal business from cache failed, detail: %v, err: %v, rid: %s", bizs, err, kit.Rid)
				hitError = err
				return
			}

			bizDetails = list
		}()
	} else {
		bizDetails = make([]mapstr.MapStr, 0)
		wg.Done()
	}

	go func() {
		defer func() { wg.Done() }()
		header := kit.NewHeader()
		opt := &metadata.ListWithIDOption{
			IDs: setList,
			Fields: []string{common.BKSetIDField, common.BKSetNameField, common.BKParentIDField,
				common.BKDefaultField},
		}
		sets, err := lgc.CoreAPI.CacheService().Cache().Topology().ListSets(kit.Ctx, header, opt)
		if err != nil {
			blog.Errorf("list sets from cache failed, err: %v, rid: %s", err, kit.Rid)
			hitError = err
			return
		}

		list := make([]mapstr.MapStr, 0)
		if err := json.Unmarshal([]byte(sets), &list); err != nil {
			blog.Errorf("unmarshal sets from cache failed, detail: %v, err: %v, rid: %s", sets, err, kit.Rid)
			hitError = err
			return
		}

		setDetails = list
	}()

	go func() {
		defer func() { wg.Done() }()
		header := kit.NewHeader()
		opt := &metadata.ListWithIDOption{
			IDs: moduleList,
			Fields: []string{common.BKModuleIDField, common.BKModuleNameField, common.BKSetIDField,
				common.BKDefaultField, common.BKAppIDField},
		}
		modules, err := lgc.CoreAPI.CacheService().Cache().Topology().ListModules(kit.Ctx, header, opt)
		if err != nil {
			blog.Errorf("list modules from cache failed, err: %v, rid: %s", err, kit.Rid)
			hitError = err
			return
		}

		list := make([]mapstr.MapStr, 0)
		if err := json.Unmarshal([]byte(modules), &list); err != nil {
			blog.Errorf("unmarshal modules from cache failed, detail: %v, err: %v, rid: %s", modules, err, kit.Rid)
			hitError = err
			return
		}

		moduleDetails = list
	}()

	wg.Wait()

	if hitError != nil {
		return nil, nil, nil, hitError
	}

	return bizDetails, setDetails, moduleDetails, nil
}

// get mainline topology rank with
func (lgc *Logics) getTopologyRank(kit *rest.Kit) (map[string]string, map[string]string, []string, error) {
	mainlineFilter := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline},
	}
	mainline, err := lgc.CoreAPI.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, mainlineFilter)
	if err != nil {
		blog.Errorf("get mainline association failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, nil, nil, err
	}

	if !mainline.Result {
		blog.Errorf("get mainline association failed, err: %s, rid: %s", mainline.ErrMsg, kit.Rid)
		return nil, nil, nil, errors.New(mainline.Code, mainline.ErrMsg)
	}

	rankMap := make(map[string]string)
	reverseRankMap := make(map[string]string)
	for _, one := range mainline.Data.Info {
		// from host to biz
		// host:module;module:set;set:biz
		rankMap[one.ObjectID] = one.AsstObjID
		reverseRankMap[one.AsstObjID] = one.ObjectID
	}

	rank := make([]string, 0)
	next := "biz"
	rank = append(rank, next)
	for _, relation := range mainline.Data.Info {
		if relation.AsstObjID == next {
			rank = append(rank, relation.ObjectID)
			next = relation.ObjectID
			continue
		} else {
			for _, rel := range mainline.Data.Info {
				if rel.AsstObjID == next {
					rank = append(rank, rel.ObjectID)
					next = rel.ObjectID
					break
				}
			}
		}
	}
	return rankMap, reverseRankMap, rank, nil
}

// get biz custom instances with object
func (lgc *Logics) getCustomObjectInstances(kit *rest.Kit, obj string, instIDs []interface{}) (
	[]mapstr.MapStr, error) {

	opts := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKInstIDField: mapstr.MapStr{common.BKDBIN: instIDs}, common.BKObjIDField: obj},
		Fields:         []string{common.BKInstIDField, common.BKInstNameField, common.BKParentIDField},
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		DisableCounter: true,
	}

	instRes, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, obj, opts)
	if err != nil {
		blog.Errorf("get biz custom object instances failed, options: %s, err: %s, rid: %s", opts, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !instRes.Result {
		blog.Errorf("get biz custom object instances failed, options: %s, err: %s, rid: %s", opts,
			instRes.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(instRes.Code, instRes.ErrMsg)
	}

	return instRes.Data.Info, nil
}
