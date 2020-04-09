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

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

func (lgc *Logics) GetHostAttributes(ctx context.Context, ownerID string, businessMedatadata *metadata.Metadata) ([]metadata.Header, error) {
	searchOp := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).MapStr()
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
	result, err := lgc.CoreAPI.CoreService().Host().GetHostByID(ctx, lgc.header, hostID)
	if err != nil {
		blog.Errorf("GetHostInstanceDetails http do error, err:%s, input:%+v, rid:%s", err.Error(), hostID, lgc.rid)
		return nil, "", lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostInstanceDetails http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, hostID, lgc.rid)
		return nil, "", lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	hostInfo := result.Data
	if len(hostInfo) == 0 {
		return nil, "", nil
	}
	ip, ok := hostInfo[common.BKHostInnerIPField].(string)
	if !ok {
		blog.Errorf("GetHostInstanceDetails http response format error,convert bk_biz_id to int error, inst:%#v  input:%#v, rid:%s", hostInfo, hostID, lgc.rid)
		return nil, "", lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostInnerIPField, "string", "not string")

	}
	return hostInfo, ip, nil
}

// GetConfigByCond get hosts owened set, module info, where hosts must match condition specify by cond.
func (lgc *Logics) GetConfigByCond(ctx context.Context, input metadata.HostModuleRelationRequest) ([]metadata.ModuleHost, errors.CCError) {

	result, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx, lgc.header, &input)
	if err != nil {
		blog.Errorf("GetConfigByCond http do error, err:%s, input:%+v, rid:%s", err.Error(), input, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetConfigByCond http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, input, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

// EnterIP 将机器导入到指定模块或者空闲模块， 已经存在机器，不操作
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

		result, err := lgc.CoreAPI.CoreService().Instance().CreateInstance(ctx, lgc.header, common.BKInnerObjIDHost, &metadata.CreateModelInstance{Data: host}) //HostController().Host().AddHost(ctx, lgc.header, host)
		if err != nil {
			blog.Errorf("EnterIP http do error, err:%s, input:%+v, rid:%s", err.Error(), host, lgc.rid)
			return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("EnterIP http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, host, lgc.rid)
			return lgc.ccErr.New(result.Code, result.ErrMsg)
		}
		// add create host log
		audit := lgc.NewHostLog(ctx, ownerID)
		if err := audit.WithPrevious(ctx, strconv.FormatInt(hostID, 10), nil); err != nil {
			return err
		}
		content := audit.GetContent(hostID)
		log := metadata.SaveAuditLogParams{
			ID:      hostID,
			Model:   common.BKInnerObjIDHost,
			Content: content,
			OpDesc:  "enter ip host",
			OpType:  auditoplog.AuditOpTypeAdd,
			ExtKey:  audit.ip,
			BizID:   appID,
		}

		aResult, err := lgc.CoreAPI.CoreService().Audit().SaveAuditLog(context.Background(), lgc.header, log)
		if err != nil {
			blog.Errorf("EnterIP AddHostLog http do error, err:%s, rid:%s", err.Error(), lgc.rid)
			return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !aResult.Result {
			blog.Errorf("EnterIP AddHostLog http response error, err code:%d, err msg:%s, rid:%s", result.Code, result.ErrMsg, lgc.rid)
			return lgc.ccErr.New(aResult.Code, aResult.ErrMsg)
		}

		hostID = int64(result.Data.Created.ID)
	} else if false == isIncrement {
		// Not an additional relationship model
		return nil
	} else {

		hostID, err = util.GetInt64ByInterface(hostList[0][common.BKHostIDField])
		if err != nil {
			blog.Errorf("EnterIP  get hostID error, err:%s,inst:%+v,input:%+v, rid:%s", err.Error(), hostList[0], host, lgc.rid)
			return lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", err.Error()) // "查询主机信息失败"
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

	params := &metadata.HostsModuleRelation{
		ApplicationID: appID,
		HostID:        []int64{hostID},
		ModuleID:      []int64{moduleID},
		IsIncrement:   isIncrement,
	}
	hmResult, ccErr := lgc.CoreAPI.CoreService().Host().TransferToNormalModule(ctx, lgc.header, params)
	if ccErr != nil {
		blog.Errorf("Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, err:%s, rid:%s", appID, hostID, err.Error(), lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hmResult.Result {
		blog.Errorf("transfer host to normal module failed, error params:{appID:%d, hostID:%d}, result:%#v, rid:%s", appID, hostID, hmResult, lgc.rid)
		if len(hmResult.Data) > 0 {
			return lgc.ccErr.New(int(hmResult.Data[0].Code), hmResult.Data[0].Message)
		}
		return lgc.ccErr.New(hmResult.Code, hmResult.ErrMsg)
	}

	audit := lgc.NewHostLog(ctx, ownerID)
	if err := audit.WithPrevious(ctx, strconv.FormatInt(hostID, 10), nil); err != nil {
		return err
	}
	hmAudit := lgc.NewHostModuleLog([]int64{hostID})
	if err := hmAudit.WithPrevious(ctx); err != nil {
		return err
	}
	if err := hmAudit.SaveAudit(ctx, appID, util.GetUser(lgc.header), "host module change"); err != nil {
		return err
	}
	return nil
}

func (lgc *Logics) GetHostInfoByConds(ctx context.Context, cond map[string]interface{}) ([]mapstr.MapStr, errors.CCErrorCoder) {
	query := &metadata.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKHostIDField,
	}

	result, err := lgc.CoreAPI.CoreService().Host().GetHosts(ctx, lgc.header, query)
	if err != nil {
		blog.Errorf("GetHostInfoByConds GetHosts http do error, err:%s, input:%+v,rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if err := result.CCError(); err != nil {
		blog.Errorf("GetHostInfoByConds GetHosts http response error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
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
func (lgc *Logics) GetHostIDByCond(ctx context.Context, cond metadata.HostModuleRelationRequest) ([]int64, errors.CCError) {

	result, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx, lgc.header, &cond)
	if err != nil {
		blog.Errorf("GetHostIDByCond GetModulesHostConfig http do error, err:%s, input:%+v,rid:%s", err.Error(), cond, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHostIDByCond GetModulesHostConfig http response error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, cond, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	hostIDs := make([]int64, 0)
	for _, val := range result.Data.Info {
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
func (lgc *Logics) GetHostModuleRelation(ctx context.Context, cond metadata.HostModuleRelationRequest) (*metadata.HostConfigData, errors.CCErrorCoder) {

	if cond.Empty() {
		return nil, lgc.ccErr.CCError(common.CCErrCommHTTPBodyEmpty)
	}

	if cond.Page.IsIllegal() {
		return nil, lgc.ccErr.CCError(common.CCErrCommPageLimitIsExceeded)
	}

	result, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx, lgc.header, &cond)
	if err != nil {
		blog.Errorf("GetHostModuleRelation http do error, err:%s, input:%+v, rid:%s", err.Error(), cond, lgc.rid)
		return nil, lgc.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if retErr := result.CCError(); retErr != nil {
		blog.Errorf("GetHostModuleRelation http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, cond, lgc.rid)
		return nil, retErr
	}

	return &result.Data, nil
}

// TransferHostAcrossBusiness  Transfer host across business,
// delete old business  host and module relation
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
	// auth: check host authorization
	if err := lgc.AuthManager.AuthorizeByHostsIDs(ctx, lgc.header, meta.MoveHostToAnotherBizModule, hostID); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostID, err, lgc.rid)
		return lgc.ccErr.Errorf(common.CCErrCommAuthorizeFailed)
	}
	// auth: deregister
	if err := lgc.AuthManager.DeregisterHostsByID(ctx, lgc.header, hostID); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v, rid: %s", hostID, err, lgc.rid)
		return lgc.ccErr.Errorf(common.CCErrCommUnRegistResourceToIAMFailed)
	}
	conf := &metadata.TransferHostsCrossBusinessRequest{SrcApplicationID: srcBizID, HostIDArr: []int64{hostID}, DstApplicationID: dstAppID, DstModuleIDArr: moduleID}
	delRet, doErr := lgc.CoreAPI.CoreService().Host().TransferToAnotherBusiness(ctx, lgc.header, conf)
	if err != nil {
		blog.Errorf("TransferHostAcrossBusiness http do error, err:%s, input:%+v, rid:%s", doErr.Error(), conf, lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !delRet.Result {
		blog.Errorf("TransferHostAcrossBusiness http response error, err code:%d, err msg:%s, input:%#v, rid:%s", delRet.Code, delRet.ErrMsg, conf, lgc.rid)
		return lgc.ccErr.New(delRet.Code, delRet.ErrMsg)
	}

	if err := audit.SaveAudit(ctx, srcBizID, lgc.user, "host to other bussiness module"); err != nil {
		blog.Errorf("TransferHostAcrossBusiness, get prev module host config failed, err: %v,hostID:%d,oldbizID:%d,appID:%d, moduleID:%#v,rid:%s", err, hostID, srcBizID, dstAppID, moduleID, lgc.rid)
		return lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")

	}

	// auth: register host
	if err := lgc.AuthManager.RegisterHostsByID(ctx, lgc.header, hostID); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid: %s", hostID, err, lgc.rid)
		return lgc.ccErr.Errorf(common.CCErrCommRegistResourceToIAMFailed)
	}
	return nil
}

// DeleteHostFromBusiness  delete host from business,
func (lgc *Logics) DeleteHostFromBusiness(ctx context.Context, bizID int64, hostIDArr []int64) ([]metadata.ExceptionResult, errors.CCError) {

	// auth: check host authorization
	if err := lgc.AuthManager.AuthorizeByHostsIDs(ctx, lgc.header, meta.MoveHostFromModuleToResPool, hostIDArr...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", hostIDArr, err, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrCommAuthorizeFailed)
	}
	// auth: deregister
	if err := lgc.AuthManager.DeregisterHostsByID(ctx, lgc.header, hostIDArr...); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v, rid: %s", hostIDArr, err, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrCommUnRegistResourceToIAMFailed)
	}
	audit := lgc.NewHostModuleLog(hostIDArr)
	if err := audit.WithPrevious(ctx); err != nil {
		blog.Errorf("TransferHostAcrossBusiness, get prev module host config failed, err: %v,hostID:%#v,appID:%d,rid:%s", err, hostIDArr, bizID, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}

	input := &metadata.DeleteHostRequest{
		ApplicationID: bizID,
		HostIDArr:     hostIDArr,
	}
	result, err := lgc.CoreAPI.CoreService().Host().DeleteHostFromSystem(ctx, lgc.header, input)
	if err != nil {
		blog.Errorf("TransferHostAcrossBusiness DeleteHost error, err: %v,hostID:%#v,appID:%d,rid:%s", err, hostIDArr, bizID, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("TransferHostAcrossBusiness DeleteHost failed, err: %v,hostID:%#v,appID:%d,rid:%s", err, hostIDArr, bizID, lgc.rid)
		return result.Data, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	if err := audit.SaveAudit(ctx, bizID, lgc.user, "delete host from business"); err != nil {
		blog.Errorf("DeleteHostFromBusiness, get prev module host config failed, err: %v,appID:%d, hostID:%#v,rid:%s", err, bizID, hostIDArr, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")
	}

	// auth: register host
	if err := lgc.AuthManager.RegisterHostsByID(ctx, lgc.header, hostIDArr...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid: %s", hostIDArr, err, lgc.rid)
		return nil, lgc.ccErr.Errorf(common.CCErrCommRegistResourceToIAMFailed)
	}
	return nil, nil
}

// CloneHostProperty clone host info and host and module relation in same application
func (lgc *Logics) CloneHostProperty(ctx context.Context, appID int64, srcHostID int64, dstHostID int64) errors.CCErrorCoder {

	// source host belong app
	ok, err := lgc.IsHostExistInApp(ctx, appID, srcHostID)
	if err != nil {
		blog.Errorf("IsHostExistInApp error. err:%s, params:{appID:%d, hostID:%d}, rid:%s", err.Error(), srcHostID, lgc.rid)
		return err
	}
	if !ok {
		blog.Errorf("Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, rid:%s", appID, srcHostID, lgc.rid)
		return lgc.ccErr.CCErrorf(common.CCErrHostNotINAPPFail, srcHostID)
	}

	// destination host belong app
	ok, err = lgc.IsHostExistInApp(ctx, appID, dstHostID)
	if err != nil {
		blog.Errorf("IsHostExistInApp error. err:%s, params:{appID:%d, hostID:%d}, rid:%s", err.Error(), dstHostID, lgc.rid)
		return err
	}
	if !ok {
		blog.Errorf("Host does not belong to the current application; error, params:{appID:%d, hostID:%d}, rid:%s", appID, dstHostID, lgc.rid)
		return lgc.ccErr.CCErrorf(common.CCErrHostNotINAPPFail, dstHostID)
	}

	hostInfoArr, err := lgc.GetHostInfoByConds(ctx, map[string]interface{}{common.BKHostIDField: srcHostID})
	if err != nil {
		return err
	}
	if len(hostInfoArr) == 0 {
		blog.Errorf("host not found. hostID:%s, rid:%s", srcHostID, lgc.rid)
		return lgc.ccErr.CCErrorf(common.CCErrHostNotFound)
	}
	srcHostInfo := hostInfoArr[0]

	delete(srcHostInfo, common.BKHostIDField)
	delete(srcHostInfo, common.CreateTimeField)
	delete(srcHostInfo, common.BKHostInnerIPField)
	delete(srcHostInfo, common.BKHostOuterIPField)
	delete(srcHostInfo, common.BKAssetIDField)
	delete(srcHostInfo, common.BKSNField)
	delete(srcHostInfo, common.BKImportFrom)

	// get source host and module relation
	hostModuleRelationCond := metadata.HostModuleRelationRequest{
		ApplicationID: appID,
		HostIDArr:     []int64{srcHostID},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
			Start: 0,
		},
	}
	relationArr, err := lgc.GetHostModuleRelation(ctx, hostModuleRelationCond)
	if err != nil {
		return err
	}
	var moduleIDArr []int64
	for _, relation := range relationArr.Info {
		moduleIDArr = append(moduleIDArr, relation.ModuleID)
	}

	exist, err := lgc.ExistInnerModule(ctx, moduleIDArr)
	if err != nil {
		return err
	}
	if exist {
		if len(moduleIDArr) != 1 {
			return lgc.ccErr.CCErrorf(common.CCErrHostModuleIDNotFoundORHasMultipleInnerModuleIDFailed)
		}
		dstModuleHostRelation := &metadata.TransferHostToInnerModule{
			ApplicationID: appID,
			HostID:        []int64{dstHostID},
			ModuleID:      moduleIDArr[0],
		}
		relationRet, doErr := lgc.CoreAPI.CoreService().Host().TransferToInnerModule(ctx, lgc.header, dstModuleHostRelation)
		if doErr != nil {
			blog.ErrorJSON("CloneHostProperty UpdateInstance error. err: %s,condition:%s,rid:%s", doErr, relationRet, lgc.rid)
			return lgc.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if err := relationRet.CCError(); err != nil {
			return err
		}
	} else {
		// destination host new module relation
		dstModuleHostRelation := &metadata.HostsModuleRelation{
			ApplicationID: appID,
			HostID:        []int64{dstHostID},
			ModuleID:      moduleIDArr,
			IsIncrement:   false,
		}
		relationRet, doErr := lgc.CoreAPI.CoreService().Host().TransferToNormalModule(ctx, lgc.header, dstModuleHostRelation)
		if doErr != nil {
			blog.ErrorJSON("CloneHostProperty UpdateInstance error. err: %s,condition:%s,rid:%s", doErr, relationRet, lgc.rid)
			return lgc.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if err := relationRet.CCError(); err != nil {
			return err
		}
	}

	input := &metadata.UpdateOption{
		Data: srcHostInfo,
		Condition: mapstr.MapStr{
			common.BKHostIDField: dstHostID,
		},
	}
	result, doErr := lgc.CoreAPI.CoreService().Instance().UpdateInstance(ctx, lgc.header, common.BKInnerObjIDHost, input)
	if doErr != nil {
		blog.ErrorJSON("CloneHostProperty UpdateInstance error. err: %s,condition:%s,rid:%s", doErr, input, lgc.rid)
		return lgc.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if err := result.CCError(); err != nil {
		blog.ErrorJSON("CloneHostProperty UpdateInstance  replay error. err: %s,condition:%s,rid:%s", err, input, lgc.rid)
		return err
	}

	return nil
}

// IPCloudToHost get host id by ip and cloud
func (lgc *Logics) IPCloudToHost(ctx context.Context, ip string, cloudID int64) (HostMap mapstr.MapStr, hostID int64, err errors.CCErrorCoder) {
	// FIXME there must be a better ip to hostID solution
	condition := mapstr.MapStr{
		common.BKHostInnerIPField: ip,
		common.BKCloudIDField:     cloudID,
	}

	hostInfoArr, err := lgc.GetHostInfoByConds(ctx, condition)
	if err != nil {
		blog.ErrorJSON("IPCloudToHost GetHostInfoByConds error. err:%s, conditon:%s, rid:%s", err.Error(), condition, lgc.rid)
		return nil, 0, err
	}
	if len(hostInfoArr) == 0 {
		return nil, 0, nil
	}

	hostID, convErr := hostInfoArr[0].Int64(common.BKHostIDField)
	if nil != convErr {
		blog.ErrorJSON("IPCloudToHost bk_host_id field not found hostMap:%s ip:%s, cloudID:%s,rid:%s", hostInfoArr, ip, cloudID, lgc.rid)
		return nil, 0, lgc.ccErr.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", convErr.Error())
	}

	return hostInfoArr[0], hostID, nil
}
