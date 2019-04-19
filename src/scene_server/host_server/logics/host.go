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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

func (lgc *Logics) GetHostAttributes(ctx context.Context, ownerID string, businessMedatadata *metadata.Metadata) ([]metadata.Header, error) {
	searchOp := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).WithOwnerID(lgc.ownerID).WithAttrComm().MapStr()
	if businessMedatadata != nil {
		searchOp.Set(common.MetadataField, businessMedatadata)
	}
	query := &metadata.QueryCondition{
		Condition: searchOp,
	}
	result, err := lgc.CoreAPI.CoreService().Model().ReadModelAttr(ctx, lgc.header, common.BKInnerObjIDHost, query)
	if err != nil {
		blog.Errorf("GetHostAttributes http do error, err:%s, input:%+v, rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostAttributes http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	headers := make([]metadata.Header, 0)
	for _, p := range result.Data.Info {
		if p.PropertyID == common.BKChildStr {
			continue
		}
		headers = append(headers, metadata.Header{
			PropertyID:   p.PropertyID,
			PropertyName: p.PropertyName,
		})
	}

	return headers, nil
}

func (lgc *Logics) GetHostInstanceDetails(ctx context.Context, ownerID, hostID string) (map[string]interface{}, string, errors.CCError) {
	// get host details, pre data
	result, err := lgc.CoreAPI.HostController().Host().GetHostByID(ctx, hostID, lgc.header)
	if err != nil {
		blog.Errorf("GetHostInstanceDetails http do error, err:%s, input:%+v, rid:%s", err.Error(), hostID, lgc.rid)
		return nil, "", lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostInstanceDetails http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, hostID, lgc.rid)
		return nil, "", lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	hostInfo := result.Data
	ip, ok := hostInfo[common.BKHostInnerIPField].(string)
	if !ok {
		blog.Errorf("GetHostInstanceDetails http response format error,convert bk_biz_id to int error, inst:%#v  input:%#v, rid:%s", hostInfo, hostID, lgc.rid)
		return nil, "", lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDHost, common.BKHostInnerIPField, "string", err.Error())

	}
	return hostInfo, ip, nil
}

func (lgc *Logics) GetConfigByCond(ctx context.Context, cond map[string][]int64) ([]map[string]int64, errors.CCError) {
	configArr := make([]map[string]int64, 0)

	if 0 == len(cond) {
		return configArr, nil
	}

	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, lgc.header, cond)
	if err != nil {
		blog.Errorf("GetConfigByCond http do error, err:%s, input:%+v, rid:%s", err.Error(), cond, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetConfigByCond http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, cond, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	for _, info := range result.Data {
		data := make(map[string]int64)
		data[common.BKAppIDField] = info.AppID
		data[common.BKSetIDField] = info.SetID
		data[common.BKModuleIDField] = info.ModuleID
		data[common.BKHostIDField] = info.HostID
		configArr = append(configArr, data)
	}
	return configArr, nil
}

// EnterIP 将机器导入到制定模块或者空闲机器， 已经存在机器，不操作
func (lgc *Logics) EnterIP(ctx context.Context, ownerID string, appID, moduleID int64, ip string, cloudID int64, host map[string]interface{}, isIncrement bool) errors.CCError {

	isExist, err := lgc.IsPlatExist(ctx, mapstr.MapStr{common.BKCloudIDField: cloudID})
	if nil != err {
		return err
	}
	if !isExist {
		return lgc.ccErr.Errorf(common.CCErrTopoCloudNotFound)
	}
	conds := map[string]interface{}{
		common.BKHostInnerIPField: ip,
		common.BKCloudIDField:     cloudID,
	}
	hostList, err := lgc.GetHostInfoByConds(ctx, conds)
	if nil != err {
		return err
	}

	hostID := int64(0)
	if len(hostList) == 0 {
		//host not exist, add host
		host[common.BKHostInnerIPField] = ip
		host[common.BKCloudIDField] = cloudID
		host["import_from"] = common.HostAddMethodAgent
		defaultFields, hasErr := lgc.getHostFields(ctx, ownerID)
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

		result, err := lgc.CoreAPI.HostController().Host().AddHost(ctx, lgc.header, host)
		if err != nil {
			blog.Errorf("EnterIP http do error, err:%s, input:%+v, rid:%s", err.Error(), host, lgc.rid)
			return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("EnterIP http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, host, lgc.rid)
			return lgc.ccErr.New(result.Code, result.ErrMsg)
		}

		retHost := result.Data.(map[string]interface{})
		hostID, err = util.GetInt64ByInterface(retHost[common.BKHostIDField])
		if err != nil {
			blog.Errorf("EnterIP  get hostID error, err:%s, reply:%+v,input:%+v, rid:%s", err.Error(), retHost, host, lgc.rid)
			return lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", err.Error())
		}
	} else if false == isIncrement {
		//Not an additional relationship model
		return nil
	} else {

		hostID, err = util.GetInt64ByInterface(hostList[0][common.BKHostIDField])
		if err != nil {
			blog.Errorf("EnterIP  get hostID error, err:%s,inst:%+v,input:%+v, rid:%s", err.Error(), hostList[0], host, lgc.rid)
			return lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", err.Error()) // "查询主机信息失败"
		}

		bl, hasErr := lgc.IsHostExistInApp(ctx, appID, hostID)
		if nil != hasErr {
			return hasErr

		}
		if false == bl {
			blog.Errorf("Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, rid:%s", appID, hostID, lgc.rid)
			return lgc.ccErr.Errorf(common.CCErrHostNotINAPPFail, hostID)
		}

	}

	//del host relation from default  module
	conf := &metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		HostID:        hostID,
	}
	result, err := lgc.CoreAPI.HostController().Module().DelDefaultModuleHostConfig(ctx, lgc.header, conf)
	if err != nil {
		blog.Errorf("EnterIP DelDefaultModuleHostConfig http do error, err:%s, input:%+v, rid:%s", err.Error(), conf, lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("EnterIP DelDefaultModuleHostConfig http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, conf, lgc.rid)
		return lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	cfg := &metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		ModuleID:      []int64{moduleID},
		HostID:        hostID,
	}
	result, err = lgc.CoreAPI.HostController().Module().AddModuleHostConfig(ctx, lgc.header, cfg)
	if err != nil {
		blog.Errorf("EnterIP AddModuleHostConfig http do error, err:%s, input:%+v, rid:%s", err.Error(), cfg, lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("EnterIP AddModuleHostConfig http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, cfg, lgc.rid)
		return lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	audit := lgc.NewHostLog(ctx, ownerID)
	if err := audit.WithPrevious(ctx, strconv.FormatInt(hostID, 10), nil); err != nil {
		return err
	}
	content := audit.GetContent(hostID)
	log := common.KvMap{common.BKContentField: content, common.BKOpDescField: "enter ip host", common.BKHostInnerIPField: audit.ip, common.BKOpTypeField: auditoplog.AuditOpTypeAdd, "inst_id": hostID}
	aResult, err := lgc.CoreAPI.AuditController().AddHostLog(ctx, ownerID, strconv.FormatInt(appID, 10), util.GetUser(lgc.header), lgc.header, log)
	if err != nil {
		blog.Errorf("EnterIP AddHostLog http do error, err:%s, rid:%s", err.Error(), lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !aResult.Result {
		blog.Errorf("EnterIP AddHostLog http response error, err code:%d, err msg:%s, rid:%s", result.Code, result.ErrMsg, lgc.rid)
		return lgc.ccErr.New(aResult.Code, aResult.ErrMsg)
	}

	hmAudit := lgc.NewHostModuleLog([]int64{hostID})
	if err := hmAudit.WithPrevious(ctx); err != nil {
		return err
	}
	if err := hmAudit.SaveAudit(ctx, strconv.FormatInt(appID, 10), util.GetUser(lgc.header), "host module change"); err != nil {
		return err
	}
	return nil
}

func (lgc *Logics) GetHostInfoByConds(ctx context.Context, cond map[string]interface{}) ([]mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKHostIDField,
	}

	result, err := lgc.CoreAPI.HostController().Host().GetHosts(ctx, lgc.header, query)
	if err != nil {
		blog.Errorf("GetHostInfoByConds GetHosts http do error, err:%s, input:%+v,rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostInfoByConds GetHosts http response error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

// HostSearch search host by mutiple condition
const (
	SplitFlag      = "##"
	TopoSetName    = "TopSetName"
	TopoModuleName = "TopModuleName"
)

func (lgc *Logics) GetHostIDByCond(ctx context.Context, cond map[string][]int64) ([]int64, errors.CCError) {
	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, lgc.header, cond)
	if err != nil {
		blog.Errorf("GetHostIDByCond GetModulesHostConfig http do error, err:%s, input:%+v,rid:%s", err.Error(), cond, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostIDByCond GetModulesHostConfig http response error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, cond, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	hostIDs := make([]int64, 0)
	for _, val := range result.Data {
		hostIDs = append(hostIDs, val.HostID)
	}

	return hostIDs, nil
}

// DeleteHostBusinessAttributes delete host business private property
func (lgc *Logics) DeleteHostBusinessAttributes(ctx context.Context, hostIDArr []int64, businessMedatadata *metadata.Metadata) error {

	return nil
}

// GetHostModuleRelation  query host and module relation,
// condition key use appID, moduleID,setID,HostID
func (lgc *Logics) GetHostModuleRelation(ctx context.Context, cond map[string][]int64) ([]metadata.ModuleHost, errors.CCError) {

	if 0 == len(cond) {
		return nil, nil
	}

	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, lgc.header, cond)
	if err != nil {
		blog.Errorf("GetHostModuleRelation http do error, err:%s, input:%+v, rid:%s", err.Error(), cond, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostModuleRelation http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, cond, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	return result.Data, nil
}

// TransferHostAcrossBusiness  Transfer host across business,
// delete old business  host and module reltaion
func (lgc *Logics) TransferHostAcrossBusiness(ctx context.Context, srcBizID, dstAppID, hostID int64, moduleID []int64) errors.CCError {

	bl, err := lgc.IsHostExistInApp(ctx, srcBizID, hostID)
	if err != nil {
		blog.Errorf("TransferHostAcrossBusiness IsHostExistInApp err:%s,input:{appID:%d,hostID:%d},rid:%s", err.Error(), srcBizID, hostID, lgc.rid)
		return err
	}
	if !bl {
		blog.Errorf("TransferHostAcrossBusiness Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, rid:%s", srcBizID, hostID, lgc.rid)
		return lgc.ccErr.Errorf(common.CCErrHostNotINAPPFail, hostID)
	}
	audit := lgc.NewHostModuleLog([]int64{hostID})
	if err := audit.WithPrevious(ctx); err != nil {
		blog.Errorf("TransferHostAcrossBusiness, get prev module host config failed, err: %v,hostID:%d,oldbizID:%d,appID:%d, moduleID:%#v,rid:%s", err, hostID, srcBizID, dstAppID, moduleID, lgc.rid)
		return lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}
	delCond := &metadata.ModuleHostConfigParams{ApplicationID: srcBizID, HostID: hostID}
	delRet, err := lgc.CoreAPI.HostController().Module().DelModuleHostConfig(ctx, lgc.header, delCond)
	if err != nil {
		blog.Errorf("TransferHostAcrossBusiness http do error, err:%s, input:%+v, rid:%s", err.Error(), delCond, lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !delRet.Result {
		blog.Errorf("TransferHostAcrossBusiness http response error, err code:%d, err msg:%s, input:%#v, rid:%s", delRet.Code, delRet.ErrMsg, delCond, lgc.rid)
		return lgc.ccErr.New(delRet.Code, delRet.ErrMsg)
	}

	addRelation := &metadata.ModuleHostConfigParams{
		ModuleID:      moduleID,
		ApplicationID: dstAppID,
		HostID:        hostID,
		OwnerID:       lgc.ownerID,
	}
	addRet, err := lgc.CoreAPI.HostController().Module().AddModuleHostConfig(ctx, lgc.header, addRelation)
	if err != nil {
		blog.Errorf("TransferHostAcrossBusiness http do error, err:%s, input:%+v, rid:%s", err.Error(), addRelation, lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !delRet.Result {
		blog.Errorf("TransferHostAcrossBusiness http response error, err code:%d, err msg:%s, input:%#v, rid:%s", addRet.Code, addRet.ErrMsg, addRelation, lgc.rid)
		return lgc.ccErr.New(addRet.Code, addRet.ErrMsg)
	}

	if err := audit.SaveAudit(ctx, strconv.FormatInt(srcBizID, 10), lgc.user, "host to other bussiness module"); err != nil {
		blog.Errorf("TransferHostAcrossBusiness, get prev module host config failed, err: %v,hostID:%d,oldbizID:%d,appID:%d, moduleID:%#v,rid:%s", err, hostID, srcBizID, dstAppID, moduleID, lgc.rid)
		return lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")

	}

	return nil
}

// DeleteHostFromBusiness  delete host from business,
func (lgc *Logics) DeleteHostFromBusiness(ctx context.Context, bizID int64, hostIDArr []int64) ([]metadata.ExceptionResult, errors.CCError) {
	var exceptionArr []metadata.ExceptionResult

	// can delete host
	var newHostIDArr []int64
	for _, hostID := range hostIDArr {
		bl, err := lgc.IsHostExistInApp(ctx, bizID, hostID)
		if err != nil {
			blog.Errorf("DeleteHostFromBusiness IsHostExistInApp err:%s,input:{appID:%d,hostID:%#v},rid:%s", err.Error(), bizID, hostIDArr, lgc.rid)
			errCode, _ := err.(errors.CCErrorCoder)
			exceptionArr = append(exceptionArr, metadata.ExceptionResult{OriginIndex: hostID, Code: int64(errCode.GetCode()), Message: err.Error()})
			continue
		}
		if !bl {
			blog.Warnf("DeleteHostFromBusiness Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, rid:%s", bizID, hostID, lgc.rid)
			continue
		}
		newHostIDArr = append(newHostIDArr, hostID)
	}
	if len(newHostIDArr) == 0 {
		if len(exceptionArr) > 0 {
			return exceptionArr, lgc.ccErr.Error(common.CCErrDeleteHostFromBusiness)
		}
		return nil, nil
	}

	audit := lgc.NewHostModuleLog(hostIDArr)
	if err := audit.WithPrevious(ctx); err != nil {
		blog.Errorf("TransferHostAcrossBusiness, get prev module host config failed, err: %v,hostID:%#v,appID:%d,rid:%s", err, hostIDArr, bizID, lgc.rid)
		return exceptionArr, lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKHostIDField).In(newHostIDArr)

	delInput := &metadata.DeleteOption{
		Condition: cond.ToMapStr(),
	}
	delHostRet, err := lgc.CoreAPI.CoreService().Instance().DeleteInstance(ctx, lgc.header, common.BKInnerObjIDHost, delInput)
	if err != nil {
		return exceptionArr, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !delHostRet.Result {
		return exceptionArr, lgc.ccErr.New(delHostRet.Code, delHostRet.ErrMsg)
	}
	// delete failure error,can ignore, has dirty data
	for _, hostID := range newHostIDArr {
		dat := metadata.ModuleHostConfigParams{
			HostID:        hostID,
			ApplicationID: bizID,
		}
		ret, err := lgc.CoreAPI.HostController().Module().DelModuleHostConfig(ctx, lgc.header, &dat)
		if err != nil {
			blog.Warnf("DeleteHostFromBusiness DelModuleHostConfig http do error. err:%s,input:{appID:%d,hostID:%d},rid:%s", err.Error(), bizID, hostID, lgc.rid)
			continue
		}
		if !ret.Result {
			blog.Warnf("DeleteHostFromBusiness Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, rid:%s", bizID, hostID, lgc.rid)
			continue
		}
	}

	if err := audit.SaveAudit(ctx, strconv.FormatInt(bizID, 10), lgc.user, "delete host from business"); err != nil {
		blog.Errorf("DeleteHostFromBusiness, get prev module host config failed, err: %v,appID:%d, hostID:%#v,rid:%s", err, bizID, hostIDArr, lgc.rid)
		return exceptionArr, lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")

	}
	if len(exceptionArr) > 0 {
		return exceptionArr, lgc.ccErr.Error(common.CCErrDeleteHostFromBusiness)
	}
	return nil, nil
}
