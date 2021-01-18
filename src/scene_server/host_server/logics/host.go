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
	"strings"

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

// GetConfigByCond get hosts owned set, module info, where hosts must match condition specify by cond.
func (lgc *Logics) GetConfigByCond(kit *rest.Kit, input metadata.HostModuleRelationRequest) ([]metadata.ModuleHost, errors.CCError) {

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

// DeleteHostBusinessAttributes delete host business private property
func (lgc *Logics) DeleteHostBusinessAttributes(kit *rest.Kit, hostIDArr []int64, bizID int64) error {

	return nil
}

// GetHostModuleRelation  query host and module relation,
// condition key use appID, moduleID,setID,HostID
func (lgc *Logics) GetHostModuleRelation(kit *rest.Kit, cond metadata.HostModuleRelationRequest) (*metadata.HostConfigData, errors.CCErrorCoder) {

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

	result, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, &cond)
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
	util.AddModelBizIDConditon(attrCond, appID)
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
