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

// GetConfigByCond get hosts owened set, module info, where hosts must match condition specify by cond.
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
			blog.Errorf("Host does not belong to the current application; error, params:{appid:%d, hostid:%s}, rid:%s", appID, hostID, lgc.rid)
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

// HostSearch search host by multiple condition
const (
	SplitFlag      = "##"
	TopoSetName    = "TopSetName"
	TopoModuleName = "TopModuleName"
)

// GetHostIDByCond query hostIDs by condition base on cc_ModuleHostConfig
// available condition fields are bk_supplier_account, bk_biz_id, bk_host_id, bk_module_id, bk_set_id
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
