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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccErrrors "configcenter/src/common/errors"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"

	"github.com/rentiansheng/xlsx"
)

// GetHostData get host data from excel
func (lgc *Logics) GetHostData(appID int64, hostIDArr []int64, hostFields []string,
	exportCond metadata.HostCommonSearch,
	header http.Header, defLang lang.DefaultCCLanguageIf) ([]mapstr.MapStr, error) {
	rid := util.GetHTTPCCRequestID(header)
	// defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	if len(hostIDArr) == 0 && len(exportCond.Condition) == 0 {
		return nil, errors.New(defLang.Language("both_hostid_exportcond_empty"))
	}

	hostInfo := make([]mapstr.MapStr, 0)
	sHostCond := make(map[string]interface{})
	sHostCond["ip"] = make(map[string]interface{})
	sHostCond["condition"] = make([]interface{}, 0)
	sHostCond["page"] = make(map[string]interface{})
	sHostCond[common.BKAppIDField] = appID

	// hostIDStr has the higher priority
	if len(hostIDArr) != 0 {

		if len(hostIDArr) > common.BKMaxExportLimit {
			return nil, errors.New(defLang.Languagef("host_id_len_err", common.BKMaxExportLimit))
		}
		condArr := make([]interface{}, 0)

		// host condition
		condition := make(map[string]interface{})
		hostCondArr := make([]interface{}, 0)
		hostCond := make(map[string]interface{})
		hostCond["field"] = common.BKHostIDField
		hostCond["operator"] = common.BKDBIN
		hostCond["value"] = hostIDArr
		hostCondArr = append(hostCondArr, hostCond)
		condition[common.BKObjIDField] = common.BKInnerObjIDHost
		condition["fields"] = make([]string, 0)
		if len(hostFields) > 0 {
			condition["fields"] = hostFields
		}
		condition["condition"] = hostCondArr
		condArr = append(condArr, condition)

		// biz condition
		condition = make(map[string]interface{})
		condition[common.BKObjIDField] = common.BKInnerObjIDApp
		condition["fields"] = make([]interface{}, 0)
		condition["condition"] = make([]interface{}, 0)
		condArr = append(condArr, condition)

		// set condition
		condition = make(map[string]interface{})
		condition[common.BKObjIDField] = common.BKInnerObjIDSet
		condition["fields"] = make([]interface{}, 0)
		condition["condition"] = make([]interface{}, 0)
		condArr = append(condArr, condition)

		// module condition
		condition = make(map[string]interface{})
		condition[common.BKObjIDField] = common.BKInnerObjIDModule
		condition["fields"] = make([]interface{}, 0)
		condition["condition"] = make([]interface{}, 0)
		condArr = append(condArr, condition)

		sHostCond["condition"] = condArr
	} else {
		if exportCond.Page.Limit <= 0 || exportCond.Page.Limit > common.BKMaxExportLimit {
			return nil, errors.New(defLang.Languagef("export_page_limit_err", common.BKMaxExportLimit))
		}
		sHostCond["ip"] = exportCond.Ip
		sHostCond["page"] = exportCond.Page

		// set host fields
		if len(hostFields) > 0 {
			for idx, cond := range exportCond.Condition {
				if cond.ObjectID == common.BKInnerObjIDHost {
					exportCond.Condition[idx].Fields = hostFields
				}
			}
		}
		sHostCond["condition"] = exportCond.Condition
	}

	result, err := lgc.Engine.CoreAPI.ApiServer().GetHostData(context.Background(), header, sHostCond)
	if nil != err {
		blog.Errorf("GetHostData failed, search condition: %+v, err: %+v, rid: %s", sHostCond, err, rid)
		return hostInfo, err
	}

	if !result.Result {
		blog.Errorf("GetHostData failed, search condition: %+v, result: %+v, rid: %s", sHostCond, result, rid)
		return nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

// GetImportHosts get import hosts
// return inst array data, errmsg collection, error
func (lgc *Logics) GetImportHosts(f *xlsx.File, header http.Header, defLang lang.DefaultCCLanguageIf,
	modelBizID int64) (map[int]map[string]interface{}, []string, error) {

	ctx := util.NewContextFromHTTPHeader(header)
	if len(f.Sheets) == 0 {
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}

	fields, err := lgc.GetObjFieldIDs(common.BKInnerObjIDHost, nil, nil, header, modelBizID,
		common.HostAddMethodExcelDefaultIndex)

	if nil != err {
		return nil, nil, errors.New(defLang.Languagef("web_get_object_field_failure", err.Error()))
	}

	sheet := f.Sheets[0]
	if nil == sheet {
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}

	departmentMap, err := lgc.getDepartmentMap(ctx, header)
	if err != nil {
		blog.Errorf("get department failed, err: %v, rid: %s", err, util.GetHTTPCCRequestID(header))
		return nil, nil, err
	}

	_, cloudMap, err := lgc.getCloudArea(ctx, header)
	if err != nil {
		blog.Errorf("get cloud area id name map failed, err: %v, rid: %s", err, util.GetHTTPCCRequestID(header))
		return nil, nil, err
	}
	hostsInfo, errMsg, err := GetExcelData(ctx, sheet, fields, common.KvMap{"import_from": common.HostAddMethodExcel},
		true, 0, defLang, departmentMap)
	if err != nil {
		blog.Errorf("get host excel data failed, err: %v, rid: %s", err, util.GetHTTPCCRequestID(header))
		return nil, errMsg, err
	}

	for index := range hostsInfo {
		cloudStr := common.DefaultCloudName
		if _, ok := hostsInfo[index][common.BKCloudIDField]; ok {
			cloudStr = util.GetStrByInterface(hostsInfo[index][common.BKCloudIDField])
		}

		addressType := common.BKAddressingStatic
		if _, ok := hostsInfo[index][common.BKAddressingField]; ok {
			addressType = util.GetStrByInterface(hostsInfo[index][common.BKAddressingField])
		}

		if _, ok := cloudMap[cloudStr]; !ok {
			blog.Errorf("check cloud area data failed, cloud area name %s of line %d doesn't exist, rid: %s", cloudStr,
				index, util.GetHTTPCCRequestID(header))
			errMsg = append(errMsg, defLang.Languagef("import_host_cloudID_not_exist", index,
				hostsInfo[index][common.BKHostInnerIPField], cloudStr))
			return nil, errMsg, nil
		}

		hostsInfo[index][common.BKCloudIDField] = cloudMap[cloudStr]
		hostsInfo[index][common.BKAddressingField] = addressType
	}

	return hostsInfo, errMsg, nil
}

// ImportHosts import host info
func (lgc *Logics) ImportHosts(ctx context.Context, f *xlsx.File, header http.Header,
	defLang lang.DefaultCCLanguageIf, modelBizID int64, moduleID int64, opType int64,
	asstObjectUniqueIDMap map[string]int64, objectUniqueID int64) *metadata.ResponseDataMapStr {

	rid := util.GetHTTPCCRequestID(header)

	if opType == 1 {
		if _, exist := f.Sheet["association"]; !exist {
			return &metadata.ResponseDataMapStr{}
		}
		info, err := lgc.importStatisticsAssociation(ctx, header, common.BKInnerObjIDHost, f.Sheet["association"])
		if err != nil {
			blog.Errorf("ImportHosts failed, GetImportHosts error:%s, rid: %s", err.Error(), rid)
			return &metadata.ResponseDataMapStr{
				BaseResp: metadata.BaseResp{
					Result: false,
					Code:   err.GetCode(),
					ErrMsg: err.Error(),
				},
				Data: nil,
			}
		}
		return &metadata.ResponseDataMapStr{
			BaseResp: metadata.BaseResp{Result: true},
			Data:     mapstr.MapStr{"association": info},
		}

	}
	return lgc.importHosts(ctx, f, header, defLang, modelBizID, modelBizID, asstObjectUniqueIDMap, objectUniqueID)
}

func (lgc *Logics) handleAsstInfoMap(ctx context.Context, header http.Header, objID string,
	asstInfoMap map[int]metadata.ExcelAssociation, asstObjectUniqueIDMap map[string]int64,
	rid string) (map[int]metadata.ExcelAssociation, error) {

	var associationFlag []string
	for _, info := range asstInfoMap {
		associationFlag = append(associationFlag, info.ObjectAsstID)
	}
	resp, err := lgc.CoreAPI.ApiServer().FindAssociationByObjectAssociationID(ctx, header, objID,
		metadata.FindAssociationByObjectAssociationIDRequest{ObjAsstIDArr: associationFlag})
	if err != nil {
		blog.Errorf("find association by object asstID failed, err: %v, rid: %s", err, rid)
		return nil, err
	}
	tempAsstInfo := make(map[string]int64, 0)
	for _, asstInfo := range resp.Data {
		_, asstObjID := asstObjectUniqueIDMap[asstInfo.AsstObjID]
		_, objID := asstObjectUniqueIDMap[asstInfo.ObjectID]
		if asstObjID || objID {
			continue
		}
		tempAsstInfo[asstInfo.AssociationName] = asstInfo.ID
	}

	tempAsstMap := asstInfoMap
	for index, asst := range asstInfoMap {
		if _, ok := tempAsstInfo[asst.ObjectAsstID]; ok {
			delete(tempAsstMap, index)
		}
	}

	return tempAsstMap, nil
}

func (lgc *Logics) handleExcelAssociation(ctx context.Context, h http.Header, f *xlsx.File, objID string, rid string,
	asstObjectUniqueIDMap map[string]int64, objectUniqueID int64, defLang lang.DefaultCCLanguageIf,
	resp *metadata.ResponseDataMapStr) *metadata.ResponseDataMapStr {
	// if sheet name is 'association', the sheet is association data to be import
	for _, sheet := range f.Sheets {
		if sheet.Name != "association" {
			continue
		}

		asstMap, assoErrMsg := GetAssociationExcelData(sheet, common.HostAddMethodExcelAssociationIndexOffset, defLang)
		asstMap, err := lgc.handleAsstInfoMap(ctx, h, objID, asstMap, asstObjectUniqueIDMap, rid)
		if err != nil {
			blog.Errorf("handle asst info map failed, err: %v, rid: %s", err, rid)
			return resp
		}
		if len(asstMap) == 0 {
			blog.Errorf("not found association data need add, rid: %s", rid)
			return resp
		}

		asstInfoMapInput := &metadata.RequestImportAssociation{
			AssociationInfoMap:    asstMap,
			AsstObjectUniqueIDMap: asstObjectUniqueIDMap,
			ObjectUniqueID:        objectUniqueID,
		}
		asstResult, asstResultErr := lgc.CoreAPI.ApiServer().ImportAssociation(ctx, h, objID, asstInfoMapInput)
		if asstResultErr != nil {
			blog.Errorf("import host association failed, err: %v, rid: %s", asstResultErr, rid)
			resp.Code = common.CCErrCommHTTPDoRequestFailed
			resp.ErrMsg = lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(h)).Errorf(common.
				CCErrCommHTTPDoRequestFailed).Error()
			return resp
		}

		assoErrMsg = append(assoErrMsg, asstResult.Data.ErrMsgMap...)
		if resp.Result && !asstResult.Result {
			resp.BaseResp = asstResult.BaseResp
		}

		resp.Data.Set("asst_error", assoErrMsg)
		return resp
	}

	return resp
}

func (lgc *Logics) importHosts(ctx context.Context, f *xlsx.File, header http.Header, defLang lang.DefaultCCLanguageIf,
	modelBizID int64, moduleID int64, asstObjectUniqueIDMap map[string]int64,
	objectUniqueID int64) *metadata.ResponseDataMapStr {

	rid := util.ExtractRequestIDFromContext(ctx)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	resp := &metadata.ResponseDataMapStr{Data: mapstr.New()}

	hosts, errMsg, err := lgc.GetImportHosts(f, header, defLang, modelBizID)
	if err != nil {
		blog.Errorf("get import hosts failed, err: %v, rid: %s", err, rid)
		resp.Code = common.CCErrWebFileContentFail
		resp.ErrMsg = defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error()
		return resp
	}
	if len(errMsg) > 0 {
		resp.Code = common.CCErrWebFileContentFail
		resp.ErrMsg = defErr.Errorf(common.CCErrWebFileContentFail, "").Error()
		resp.Data.Set("error", errMsg)
		return resp
	}
	errMsg, err = lgc.CheckHostsAdded(ctx, header, hosts)
	if err != nil {
		blog.Errorf("check host added failed, err: %v, rid: %s", err, rid)
		resp.Code = common.CCErrWebHostCheckFail
		resp.ErrMsg = defErr.Errorf(common.CCErrWebHostCheckFail, err.Error()).Error()
		return resp
	}
	if len(errMsg) > 0 {
		resp.Code = common.CCErrWebHostCheckFail
		resp.ErrMsg = defErr.Errorf(common.CCErrWebHostCheckFail, "").Error()
		resp.Data.Set("error", errMsg)
		return resp
	}

	if len(hosts) > 0 {
		params := map[string]interface{}{
			"host_info":            hosts,
			"input_type":           common.InputTypeExcel,
			common.BKModuleIDField: moduleID,
		}
		result, resultErr := lgc.CoreAPI.ApiServer().AddHostByExcel(context.Background(), header, params)
		if resultErr != nil {
			blog.Errorf("add host info failed, err: %v, rid: %s", resultErr, rid)
			resp.Code = common.CCErrCommHTTPDoRequestFailed
			resp.ErrMsg = defErr.Errorf(common.CCErrCommHTTPDoRequestFailed).Error()
			return resp
		}

		resp = result
	}

	if len(f.Sheets) <= 2 || len(asstObjectUniqueIDMap) == 0 {
		resp.Result = true
		return resp
	}

	return lgc.handleExcelAssociation(ctx, header, f, common.BKInnerObjIDHost, rid, asstObjectUniqueIDMap,
		objectUniqueID, defLang, resp)
}

// importStatisticsAssociation TODO
// Statistics
func (lgc *Logics) importStatisticsAssociation(ctx context.Context, header http.Header, objID string,
	sheet *xlsx.Sheet) (map[string]metadata.ObjectAsstIDStatisticsInfo, ccErrrors.CCErrorCoder) {

	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	rid := util.ExtractRequestIDFromContext(ctx)

	// if len(f.Sheets) >= 2, the second sheet is association data to be import
	asstNameArr, asstInfoMap := StatisticsAssociation(sheet, common.HostAddMethodExcelAssociationIndexOffset)
	if len(asstInfoMap) == 0 {
		return nil, nil
	}

	input := metadata.FindAssociationByObjectAssociationIDRequest{
		ObjAsstIDArr: asstNameArr,
	}
	resp, err := lgc.CoreAPI.ApiServer().FindAssociationByObjectAssociationID(ctx, header, objID, input)
	if err != nil {
		blog.ErrorJSON("find model association by bk_obj_asst_id http do error. err: %s, objID: %s, input: %s, rid: %s",
			err.Error(), objID, input, rid)
		return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if ccErr := resp.CCError(); ccErr != nil {
		blog.ErrorJSON("find model association by bk_obj_asst_id http reply error. reply: %s, objID: %s, input: %s,"+
			" rid: %s", resp, objID, input, rid)
		return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	objIDStatisMap := make(map[string]metadata.ObjectAsstIDStatisticsInfo, 0)
	for _, row := range resp.Data {
		// bk_obj_asst_id
		excelAsstNameStatis := asstInfoMap[row.AssociationName]

		statisObjID := row.AsstObjID
		// 只统计关联对象
		if row.ObjectID != objID {
			statisObjID = row.ObjectID
		}

		objIDStatis, ok := objIDStatisMap[statisObjID]
		if !ok {
			objIDStatis = metadata.ObjectAsstIDStatisticsInfo{}
		}
		objIDStatis.Create += excelAsstNameStatis.Create
		objIDStatis.Delete += excelAsstNameStatis.Delete
		objIDStatis.Total += excelAsstNameStatis.Total
		objIDStatisMap[statisObjID] = objIDStatis
	}

	return objIDStatisMap, nil

}

// UpdateHosts update excel import hosts
func (lgc *Logics) UpdateHosts(ctx context.Context, f *xlsx.File, header http.Header, defLang lang.DefaultCCLanguageIf,
	modelBizID, opType int64, asstObjectUniqueIDMap map[string]int64, objectUniqueID int64) *metadata.ResponseDataMapStr {

	rid := util.ExtractRequestIDFromContext(ctx)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	hosts, errMsg, err := lgc.GetImportHosts(f, header, defLang, modelBizID)
	if err != nil {
		blog.Errorf("ImportHost get import hosts from excel err, error:%s, rid: %s", err.Error(), rid)
		return returnByErrCode(defErr, common.CCErrWebFileContentFail, nil, err.Error())
	}
	if len(errMsg) != 0 {
		return returnByErrCode(defErr, common.CCErrWebFileContentFail, mapstr.MapStr{"error": errMsg},
			strings.Join(errMsg, ","))
	}

	if opType == 1 {

		if _, exist := f.Sheet["association"]; !exist {
			return returnByErrCode(defErr, 0, mapstr.MapStr{"association": mapstr.New()}, "")
		}

		statisAsstInfo, err := lgc.importStatisticsAssociation(ctx, header, common.BKInnerObjIDHost,
			f.Sheet["association"])
		if err != nil {
			return returnByErrCode(defErr, err.GetCode(), nil, "")
		}

		return returnByErrCode(defErr, 0, mapstr.MapStr{"association": statisAsstInfo}, "")
	}

	errMsg, err = lgc.CheckHostsUpdated(ctx, header, hosts, modelBizID)
	if err != nil {
		blog.Errorf("ImportHosts failed,  CheckHostsAdded error:%s, rid: %s", err.Error(), rid)
		return returnByErrCode(defErr, common.CCErrWebHostCheckFail, nil, err.Error())
	}
	if len(errMsg) > 0 {
		return returnByErrCode(defErr, common.CCErrWebHostCheckFail, mapstr.MapStr{"error": errMsg}, "")
	}

	var resultErr error
	result := returnByErrCode(defErr, 0, mapstr.New(), "")
	if len(hosts) != 0 {
		params := map[string]interface{}{
			"host_info":  hosts,
			"input_type": common.InputTypeExcel,
		}
		result, resultErr = lgc.CoreAPI.ApiServer().UpdateHost(context.Background(), header, params)
		if resultErr != nil {
			blog.Errorf("UpdateHosts update host http request  error:%s, rid:%s", resultErr.Error(),
				util.GetHTTPCCRequestID(header))
			return returnByErrCode(defErr, common.CCErrCommHTTPDoRequestFailed, nil, "")
		}
	}

	if len(asstObjectUniqueIDMap) == 0 && objectUniqueID == 0 {
		return result
	}

	if _, exist := f.Sheet["association"]; !exist {
		return result
	}
	if len(f.Sheet["association"].Rows[common.HostAddMethodExcelAssociationIndexOffset].Cells) < 2 {
		return result
	}

	return lgc.handleExcelAssociation(ctx, header, f, common.BKInnerObjIDHost, rid, asstObjectUniqueIDMap,
		objectUniqueID, defLang, result)
}

func returnByErrCode(defErr ccErrrors.DefaultCCErrorIf, errCode int, data mapstr.MapStr,
	errStr string) *metadata.ResponseDataMapStr {

	result := true
	var errMsg string
	if errCode != 0 {
		result = false

		if len(errStr) != 0 {
			errMsg = defErr.CCErrorf(errCode, errStr).Error()
		} else {
			errMsg = defErr.CCError(errCode).Error()
		}
	}

	return &metadata.ResponseDataMapStr{
		BaseResp: metadata.BaseResp{
			Result: result,
			Code:   errCode,
			ErrMsg: errMsg,
		},
		Data: data,
	}
}

// getIpField get ipv4 and ipv6 address, ipv4 and ipv6 address cannot be null at the same time.
func getIpField(host map[string]interface{}) (string, string, string) {

	innerIP, v4Ok := host[common.BKHostInnerIPField].(string)
	innerIPv6, v6Ok := host[common.BKHostInnerIPv6Field].(string)
	if (!v4Ok || innerIP == "") && (!v6Ok || innerIPv6 == "") {
		return "host_import_innerip_v4_v6_empty", "", ""
	}

	// 存在ipv6地址的场景下需要将录入的ipv6地址转化为完整的ipv6地址
	if v6Ok && innerIPv6 != "" {
		ipv6, err := common.ConvertIPv6ToStandardFormat(innerIPv6)
		if err != nil {
			return "host_import_innerip_v6_transfer_fail", "", ""
		}
		innerIPv6 = ipv6
	}

	return "", innerIP, innerIPv6
}

// CheckHostsAdded check the hosts to be added
func (lgc *Logics) CheckHostsAdded(ctx context.Context, header http.Header, hostInfos map[int]map[string]interface{}) (
	errMsg []string, err error) {

	rid := util.ExtractRequestIDFromContext(ctx)
	ccLang := lgc.Engine.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))

	existentV4Hosts, existentV6Hosts, err := lgc.getExistHostsByInnerIPs(ctx, header, hostInfos)
	if err != nil {
		blog.Errorf("CheckHostsAdded failed, getExistHostsByInnerIPs err:%s, rid:%s", err.Error(), rid)
		return nil, err
	}

	for _, index := range util.SortedMapIntKeys(hostInfos) {
		host := hostInfos[index]
		if host == nil {
			continue
		}

		if _, ok := host[common.BKHostIDField]; ok {
			errMsg = append(errMsg, ccLang.Languagef("import_host_no_need_hostID", index))
			continue
		}
		if _, ok := host[common.BKAgentIDField]; ok {
			errMsg = append(errMsg, ccLang.Languagef("import_host_no_need_agentID", index))
			continue
		}
		cloud, ok := host[common.BKCloudIDField]
		if !ok {
			errMsg = append(errMsg, ccLang.Languagef("import_host_not_provide_cloudID", index))
			continue
		}

		addressType, ok := host[common.BKAddressingField].(string)
		if !ok {
			errMsg = append(errMsg, ccLang.Languagef("import_host_not_provide_addressing", index))
			continue
		}

		if addressType != common.BKAddressingStatic && addressType != common.BKAddressingDynamic {
			errMsg = append(errMsg, ccLang.Languagef("import_host_illegal_addressing", index))
			continue
		}

		// in dynamic scenarios, there is no need to do duplication check of ip address.
		if addressType == common.BKAddressingDynamic {
			continue
		}

		errStr, innerIPv4, innerIPv6 := getIpField(host)
		if errStr != "" {
			errMsg = append(errMsg, ccLang.Languagef(errStr, index))
			continue
		}

		// check if the host ipv4 exist in db
		key := generateHostCloudKey(innerIPv4, cloud)
		if _, exist := existentV4Hosts[key]; exist {
			errMsg = append(errMsg, ccLang.Languagef("host_import_innerip_v4_fail", index))
			continue
		}

		// check if the host ipv6 exist in db
		keyV6 := generateHostCloudKey(innerIPv6, cloud)
		if _, exist := existentV6Hosts[keyV6]; exist {
			errMsg = append(errMsg, ccLang.Languagef("host_import_innerip_v6_fail", index))
			continue
		}
	}

	return errMsg, nil
}

// CheckHostsUpdated check the hosts to be updated
func (lgc *Logics) CheckHostsUpdated(ctx context.Context, header http.Header, hostInfos map[int]map[string]interface{},
	modelBizID int64) (errMsg []string, err error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	ccLang := lgc.Engine.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))

	existentHost, err := lgc.getExistHostsByHostIDs(ctx, header, hostInfos)
	if err != nil {
		blog.Errorf("CheckHostsUpdated failed, getExistHostsByHostIDs err:%s, rid:%s", err.Error(), rid)
		return nil, err
	}

	if modelBizID == 0 {
		// get resource pool biz ID
		modelBizID, err = lgc.GetDefaultBizID(ctx, header)
		if err != nil {
			blog.Errorf("CheckHostsUpdated failed, GetDefaultBizID err:%s, rid:%s", err.Error(), rid)
			return nil, err
		}
	}

	hostBizMap, err := lgc.GetHostBizRelations(ctx, header, hostInfos, modelBizID)
	if err != nil {
		blog.Errorf("CheckHostsUpdated failed, GetHostBizRelations err:%s, rid:%s", err.Error(), rid)
		return nil, err
	}

	for _, index := range util.SortedMapIntKeys(hostInfos) {
		host := hostInfos[index]
		if host == nil {
			continue
		}

		hostID, ok := host[common.BKHostIDField]
		if !ok {
			blog.Errorf("bk_host_id field doesn't exist, innerIpv4: %v, innerIpv6: %v, rid: %v",
				host[common.BKHostInnerIPField], host[common.BKHostInnerIPv6Field], rid)
			errMsg = append(errMsg, ccLang.Languagef("import_update_host_miss_hostID", index))
			continue
		}
		hostIDVal, err := util.GetInt64ByInterface(hostID)
		if err != nil {
			errMsg = append(errMsg, ccLang.Languagef("import_update_host_hostID_not_int", index))
			continue
		}

		// check if the host exist in db
		ip := existentHost[hostIDVal].Ip
		ipv6 := existentHost[hostIDVal].Ipv6
		agentID := existentHost[hostIDVal].AgentID
		if ip == "" && ipv6 == "" && agentID == "" {
			errMsg = append(errMsg, ccLang.Languagef("import_host_no_exist_error", index, hostIDVal))
			continue
		}

		// check if the host innerIP and hostID is consistent
		excelIP := util.GetStrByInterface(host[common.BKHostInnerIPField])
		if ip != excelIP {
			errMsg = append(errMsg, ccLang.Languagef("import_host_ip_not_consistent", index, excelIP,
				hostIDVal, ip))
			continue
		}

		// check if the host innerIPv6 and hostID is consistent
		excelIPv6 := util.GetStrByInterface(host[common.BKHostInnerIPv6Field])
		if ipv6 != excelIPv6 {
			errMsg = append(errMsg, ccLang.Languagef("import_host_ipv6_not_consistent", index, excelIPv6,
				hostIDVal, ipv6))
			continue
		}

		// check if the host agentID and hostID is consistent
		excelAgentID := util.GetStrByInterface(host[common.BKAgentIDField])
		if agentID != excelAgentID {
			errMsg = append(errMsg, ccLang.Languagef("import_host_agentID_not_consistent", index, excelAgentID,
				hostIDVal, agentID))
			continue
		}

		// check if the hostID and bizID is consistent
		if hostBizMap[hostIDVal] != modelBizID {
			errMsg = append(errMsg, ccLang.Languagef("import_hostID_bizID_not_consistent", index, excelIP, excelIPv6))
			continue
		}
	}

	return errMsg, nil
}

// GetHostBizRelations get host and biz relations
func (lgc *Logics) GetHostBizRelations(ctx context.Context, header http.Header, hostInfos map[int]map[string]interface{}, bizID int64) (map[int64]int64, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	var hostIDs []int64
	for _, host := range hostInfos {
		if hostID, ok := host[common.BKHostIDField]; ok {
			if hostIDVal, err := util.GetInt64ByInterface(hostID); err == nil {
				hostIDs = append(hostIDs, hostIDVal)
			}
		}
	}

	hostLen := len(hostIDs)
	if hostLen == 0 {
		return make(map[int64]int64), nil
	}

	hostBizMap := make(map[int64]int64)
	// the length of GetHostModuleRelation's param bk_host_id can't bigger than 500
	for idx := 0; idx < hostLen; idx += 500 {
		endIdx := idx + 500
		if endIdx > hostLen {
			endIdx = hostLen
		}
		hosts := hostIDs[idx:endIdx]

		params := mapstr.MapStr{
			common.BKAppIDField:  bizID,
			common.BKHostIDField: hosts,
		}
		resp, err := lgc.CoreAPI.ApiServer().GetHostModuleRelation(ctx, header, params)
		if err != nil {
			blog.Errorf(" GetHostBizRelations failed, GetHostModuleRelation err:%v, params: %#v, rid: %s", err, params, rid)
			return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !resp.Result {
			blog.Errorf(" GetHostBizRelations failed, GetHostModuleRelation resp:%#v, params: %#v, rid: %s", resp, params, rid)
			return nil, resp.CCError()
		}

		for _, relation := range resp.Data {
			hostBizMap[relation.HostID] = relation.AppID
		}
	}

	return hostBizMap, nil
}

// GetDefaultBizID get resource pool biz ID
func (lgc *Logics) GetDefaultBizID(ctx context.Context, header http.Header) (int64, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	ownerID := util.GetOwnerID(header)

	resp, err := lgc.CoreAPI.ApiServer().SearchDefaultApp(context.Background(), header, ownerID)
	if err != nil {
		blog.Errorf(" GetDefaultBizID failed, SearchDefaultApp err:%v, ownerID: %d, rid: %s", err, ownerID, rid)
		return 0, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result {
		blog.Errorf(" GetDefaultBizID failed, SearchDefaultApp resp:%#v, ownerID: %d, rid: %s", resp, ownerID, rid)
		return 0, resp.CCError()
	}

	if len(resp.Data.Info) == 0 {
		return 0, defErr.CCError(common.CCErrHostNotResourceFail)
	}
	bizID, err := resp.Data.Info[0].Int64(common.BKAppIDField)
	if err != nil {
		return 0, defErr.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp,
			common.BKAppIDField, "int64", err.Error())
	}

	return bizID, nil
}

// generateHostCloudKey generate a cloudKey for host that is unique among clouds by appending the cloudID.
func generateHostCloudKey(ip, cloudID interface{}) string {
	return fmt.Sprintf("%v-%v", ip, cloudID)
}

// getExistHostsByInnerIPs get hosts that already in db(same bk_host_innerip host)
func (lgc *Logics) getExistHostsByInnerIPs(ctx context.Context, header http.Header,
	hostInfos map[int]map[string]interface{}) (map[string]struct{}, map[string]struct{}, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	// step1. extract all innerIP from hostInfos
	ipArr, ipv6Arr := make([]string, 0), make([]string, 0)
	for _, host := range hostInfos {
		innerIP, ok := host[common.BKHostInnerIPField].(string)
		if ok && innerIP != "" {
			ipArr = append(ipArr, innerIP)
		}
		innerIPv6, ok := host[common.BKHostInnerIPv6Field].(string)
		if ok && innerIPv6 != "" {
			ipv6Arr = append(ipv6Arr, innerIPv6)
		}
	}
	if len(ipArr) == 0 && len(ipv6Arr) == 0 {
		return make(map[string]struct{}), make(map[string]struct{}), nil
	}

	// step2. query host info by innerIPs
	innerIPs, innerIPv6s := make([]string, 0), make([]string, 0)
	for _, innerIP := range ipArr {
		innerIPArr := strings.Split(innerIP, ",")
		innerIPs = append(innerIPs, innerIPArr...)
	}
	for _, innerIPv6 := range ipv6Arr {
		innerIPv6Arr := strings.Split(innerIPv6, ",")
		innerIPv6s = append(innerIPv6s, innerIPv6Arr...)
	}

	rules := make([]querybuilder.Rule, 0)
	if len(innerIPs) > 0 {
		rules = append(rules, querybuilder.AtomRule{
			Field:    common.BKHostInnerIPField,
			Operator: querybuilder.OperatorIn,
			Value:    innerIPs,
		})
	}
	if len(innerIPv6s) > 0 {
		rules = append(rules, querybuilder.AtomRule{
			Field:    common.BKHostInnerIPv6Field,
			Operator: querybuilder.OperatorIn,
			Value:    innerIPv6s,
		})
	}

	option := metadata.ListHostsWithNoBizParameter{
		HostPropertyFilter: &querybuilder.QueryFilter{
			Rule: querybuilder.CombinedRule{
				Condition: querybuilder.ConditionOr,
				Rules:     rules,
			},
		},
		Fields: []string{common.BKHostIDField, common.BKHostInnerIPField,
			common.BKHostInnerIPv6Field, common.BKCloudIDField},
	}
	resp, err := lgc.CoreAPI.ApiServer().ListHostWithoutApp(ctx, header, option)
	if err != nil {
		blog.Errorf("list host without app failed, option: %+v, err: %v, rid: %s", option, err, rid)
		return nil, nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result {
		blog.Errorf("ListHostWithoutApp resp:%#v, option: %d, rid: %s", resp, option, rid)
		return nil, nil, resp.CCError()
	}

	// step3. arrange data as a map, cloudKey: hostID
	hostV4Map, hostV6Map := make(map[string]struct{}, 0), make(map[string]struct{}, 0)

	for _, host := range resp.Data.Info {
		ipv4 := util.GetStrByInterface(host[common.BKHostInnerIPField])
		if ipv4 != "" {
			keyV4 := generateHostCloudKey(ipv4, host[common.BKCloudIDField])
			hostV4Map[keyV4] = struct{}{}
		}

		ipv6 := util.GetStrByInterface(host[common.BKHostInnerIPv6Field])
		if ipv6 != "" {
			keyV6 := generateHostCloudKey(ipv6, host[common.BKCloudIDField])
			hostV6Map[keyV6] = struct{}{}
		}
	}
	return hostV4Map, hostV6Map, nil
}

type excelSimpleHost struct {
	Ip      string
	Ipv6    string
	AgentID string
}

// getExistHostsByHostIDs get hosts that already in db(same bk_host_id host)
// return: map[hostID]innerIP
func (lgc *Logics) getExistHostsByHostIDs(ctx context.Context, header http.Header,
	hostInfos map[int]map[string]interface{}) (map[int64]excelSimpleHost, error) {

	rid := util.ExtractRequestIDFromContext(ctx)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	// step1. extract all innerIP from hostInfos
	var hostIDs []int64
	for _, host := range hostInfos {
		if hostID, ok := host[common.BKHostIDField]; ok {
			if hostIDVal, err := util.GetInt64ByInterface(hostID); err == nil {
				hostIDs = append(hostIDs, hostIDVal)
			}
		}
	}

	if len(hostIDs) == 0 {
		return make(map[int64]excelSimpleHost), nil
	}

	// step2. query host info by hostIDs
	rules := []querybuilder.Rule{
		querybuilder.AtomRule{
			Field:    common.BKHostIDField,
			Operator: querybuilder.OperatorIn,
			Value:    hostIDs,
		},
	}

	option := metadata.ListHostsWithNoBizParameter{
		HostPropertyFilter: &querybuilder.QueryFilter{
			Rule: querybuilder.CombinedRule{
				Condition: querybuilder.ConditionOr,
				Rules:     rules,
			},
		},
		Fields: []string{
			common.BKHostIDField,
			common.BKHostInnerIPField,
			common.BKHostInnerIPv6Field,
			common.BKAgentIDField,
		},
	}
	resp, err := lgc.CoreAPI.ApiServer().ListHostWithoutApp(ctx, header, option)
	if err != nil {
		blog.Errorf(" getExistHostsByHostIDs failed, ListHostWithoutApp err:%v, option: %d, rid: %s", err, option, rid)
		return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result {
		blog.Errorf(" getExistHostsByHostIDs failed, ListHostWithoutApp resp:%#v, option: %d, rid: %s", resp, option, rid)
		return nil, resp.CCError()
	}

	// step3. arrange data as a map, cloudKey: hostID
	hostMap := make(map[int64]excelSimpleHost, 0)
	for _, host := range resp.Data.Info {
		if hostID, ok := host[common.BKHostIDField]; ok {
			if hostIDVal, err := util.GetInt64ByInterface(hostID); err == nil {
				hostMap[hostIDVal] = excelSimpleHost{
					Ip:      util.GetStrByInterface(host[common.BKHostInnerIPField]),
					Ipv6:    util.GetStrByInterface(host[common.BKHostInnerIPv6Field]),
					AgentID: util.GetStrByInterface(host[common.BKAgentIDField]),
				}

			}
		}
	}

	return hostMap, nil
}

// getCloudArea search total cloud area id and name
// return an array of cloud name and a name-id map
func (lgc *Logics) getCloudArea(ctx context.Context, header http.Header) ([]string, map[string]int64, error) {

	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	cloudArea := make([]mapstr.MapStr, 0)
	start := 0
	for {
		input := metadata.CloudAreaSearchParam{
			SearchCloudOption: metadata.SearchCloudOption{
				Fields: []string{common.BKCloudIDField, common.BKCloudNameField},
				Page:   metadata.BasePage{Start: start, Limit: common.BKMaxPageSize},
			},
			SyncTaskIDs: false,
		}
		rsp, err := lgc.Engine.CoreAPI.ApiServer().SearchCloudArea(ctx, header, input)
		if err != nil {
			blog.Errorf("search cloud area failed, err: %v, rid: %s", err, rid)
			return nil, nil, err
		}

		cloudArea = append(cloudArea, rsp.Info...)
		if len(rsp.Info) < common.BKMaxPageSize {
			break
		}

		start += common.BKMaxPageSize
	}

	if len(cloudArea) == 0 {
		blog.Errorf("search cloud area failed, return empty, rid: %s", rid)
		return nil, nil, defErr.CCError(common.CCErrTopoCloudNotFound)
	}

	cloudAreaArr := make([]string, 0)
	cloudAreaMap := make(map[string]int64)
	for _, item := range cloudArea {
		areaName, err := item.String(common.BKCloudNameField)
		if err != nil {
			blog.Errorf("get type of string cloud name failed, err: %v, rid: %s", err, rid)
			return nil, nil, err
		}
		cloudAreaArr = append(cloudAreaArr, areaName)

		areaID, err := item.Int64(common.BKCloudIDField)
		if err != nil {
			blog.Errorf("get type of int64 cloud id failed, err: %v, rid: %s", err, rid)
			return nil, nil, err
		}
		// cloud area name is unique
		cloudAreaMap[areaName] = areaID
	}

	return cloudAreaArr, cloudAreaMap, nil
}
