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
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccErrrors "configcenter/src/common/errors"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"

	"github.com/tealeg/xlsx/v3"
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

		if len(hostIDArr) > common.BKInstMaxExportLimit {
			return nil, errors.New(defLang.Languagef("host_id_len_err", common.BKInstMaxExportLimit))
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
		if exportCond.Page.Limit <= 0 || exportCond.Page.Limit > common.BKInstMaxExportLimit {
			return nil, errors.New(defLang.Languagef("export_page_limit_err", common.BKInstMaxExportLimit))
		}
		sHostCond["ip"] = exportCond.Ipv4Ip
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
	sHostCond["page"] = exportCond.Page
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

// GetImportHosts get import hosts, return inst array data, errmsg collection, error
func (lgc *Logics) GetImportHosts(header http.Header, defLang lang.DefaultCCLanguageIf, preData *ImportExcelPreData,
	start, end int) (map[int]map[string]interface{}, []string, error) {

	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)

	_, cloudMap, err := lgc.getCloudArea(ctx, header)
	if err != nil {
		blog.Errorf("get cloud area id name map failed, err: %v, rid: %s", err, rid)
		return nil, nil, err
	}

	hostsInfo, errMsg, err := GetExcelData(ctx, preData, start, end,
		common.KvMap{"import_from": common.HostAddMethodExcel}, defLang)
	if err != nil {
		blog.Errorf("get host excel data failed, err: %v, rid: %s", err, rid)
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
			blog.Errorf("check cloud area data failed, cloud area name %s of line %d doesn't exist, rid: %s",
				cloudStr, index, rid)
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

		asstMap, assoErrMsg, err := GetAssociationExcelData(sheet, common.HostAddMethodExcelAssociationIndexOffset,
			defLang)
		if err != nil {
			blog.Errorf("get association excel data failed, err: %v, rid: %s", err, rid)
			resp.Code = common.CCErrCommHTTPDoRequestFailed
			resp.ErrMsg = lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(h)).Errorf(common.
				CCErrCommHTTPDoRequestFailed).Error()
			return resp
		}

		asstMap, err = lgc.handleAsstInfoMap(ctx, h, objID, asstMap, asstObjectUniqueIDMap, rid)
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

		resp.BaseResp = asstResult.BaseResp

		if len(asstResult.Data.ErrMsgMap) > 0 {
			assoErrMsg = append(assoErrMsg, asstResult.Data.ErrMsgMap...)
			resp.Data.Set("error", assoErrMsg)
		}
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

	preData, err := lgc.getImportExcelPreData(common.BKInnerObjIDHost, header, f, defLang, modelBizID)
	if err != nil {
		blog.Errorf("get import hosts failed, err: %v, rid: %s", err, rid)
		resp.Code = common.CCErrWebFileContentFail
		resp.ErrMsg = defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error()
		return resp
	}

	var successMsgs []string
	var errMsgs []string
	for i := 0; i < len(preData.DataRange); i++ {
		rowNum := preData.DataRange[i].Start + 1
		host, errMsg, err := lgc.GetImportHosts(header, defLang, preData, i, i+1)
		if err != nil {
			blog.Errorf("get import hosts failed, err: %v, rid: %s", err, rid)
			errMsgs = append(errMsgs, defLang.Languagef("import_data_fail", rowNum, err.Error()))
			continue
		}
		if len(errMsg) > 0 {
			errMsgs = append(errMsgs, defLang.Languagef("import_data_fail", rowNum, errMsg[0]))
			continue
		}

		host, err = lgc.handleImportEnumQuoteInst(ctx, header, host, preData.Fields, rid)
		if err != nil {
			blog.Errorf("handle enum quote inst failed, err: %v, rid: %s", err, rid)
			errMsgs = append(errMsgs, defLang.Languagef("import_data_fail", rowNum, err.Error()))
			continue
		}

		errMsg, err = lgc.CheckHostsAdded(ctx, header, host)
		if err != nil {
			blog.Errorf("check host added failed, err: %v, rid: %s", err, rid)
			errMsgs = append(errMsgs, defLang.Languagef("import_data_fail", rowNum, err.Error()))
			continue
		}
		if len(errMsg) > 0 {
			errMsgs = append(errMsgs, defLang.Languagef("import_data_fail", rowNum, errMsg[0]))
			continue
		}

		if len(host) == 0 {
			continue
		}

		//params := map[string]interface{}{
		//	"host_info":            host,
		//	"input_type":           common.InputTypeExcel,
		//	common.BKModuleIDField: moduleID,
		//}
		//result, err := lgc.CoreAPI.ApiServer().AddHostByExcel(context.Background(), header, params)
		//if err != nil {
		//	errMsgs = append(errMsgs, defLang.Languagef("import_data_fail", rowNum, err.Error()))
		//	continue
		//}
		//
		//errData, exist := result.Data.Get("error")
		//if exist && errData != nil {
		//	errMsgs = append(errMsgs, defLang.Languagef("import_data_fail", rowNum, errData))
		//	continue
		//}

		successMsgs = append(successMsgs, strconv.Itoa(rowNum))
	}

	resp.Data["success"] = successMsgs
	resp.Data["error"] = errMsgs
	if len(f.Sheets) <= 2 || len(asstObjectUniqueIDMap) == 0 || len(errMsgs) != 0 {
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
	asstNameArr, asstInfoMap, err := StatisticsAssociation(sheet, common.HostAddMethodExcelAssociationIndexOffset)
	if err != nil {
		blog.ErrorJSON("ger statistics association failed, err: %s, objID: %s, input: %s, rid: %s",
			err.Error(), objID, sheet, rid)
		return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
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
// NOCC:golint/fnsize(后续重构处理)
func (lgc *Logics) UpdateHosts(ctx context.Context, f *xlsx.File, header http.Header, defLang lang.DefaultCCLanguageIf,
	modelBizID, opType int64, asstObjectUniqueIDMap map[string]int64,
	objectUniqueID int64) *metadata.ResponseDataMapStr {

	rid := util.ExtractRequestIDFromContext(ctx)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

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

	preData, err := lgc.getImportExcelPreData(common.BKInnerObjIDHost, header, f, defLang, modelBizID)
	if err != nil {
		return returnByErrCode(defErr, common.CCErrWebFileContentFail, nil, err.Error())
	}

	result := returnByErrCode(defErr, 0, mapstr.New(), "")
	var successMsgs []string
	var errMsgs []string
	for i := 0; i < len(preData.DataRange); i++ {
		rowNum := preData.DataRange[i].Start + 1
		host, errMsg, err := lgc.GetImportHosts(header, defLang, preData, i, i+1)
		if err != nil {
			blog.Errorf("get import hosts from excel err, error:%s, rid: %s", err.Error(), rid)
			errMsgs = append(errMsgs, defLang.Languagef("import_update_data_fail", rowNum, err.Error()))
			continue
		}
		if len(errMsg) != 0 {
			errMsgs = append(errMsgs, defLang.Languagef("import_update_data_fail", rowNum, errMsgs[0]))
			continue
		}

		host, err = lgc.handleImportEnumQuoteInst(ctx, header, host, preData.Fields, rid)
		if err != nil {
			blog.Errorf("handle enum quote inst failed, err: %v, rid: %s", err, rid)
			errMsgs = append(errMsgs, defLang.Languagef("import_update_data_fail", rowNum, err.Error()))
			continue
		}

		errMsg, err = lgc.CheckHostsUpdated(ctx, header, host, modelBizID)
		if err != nil {
			blog.Errorf("check hosts updated failed, err: %v, rid: %s", err, rid)
			errMsgs = append(errMsgs, defLang.Languagef("import_update_data_fail", rowNum, err.Error()))
			continue
		}
		if len(errMsg) > 0 {
			errMsgs = append(errMsgs, defLang.Languagef("import_update_data_fail", rowNum, errMsg[0]))
			continue
		}

		if len(host) == 0 {
			continue
		}

		//params := map[string]interface{}{
		//	"host_info":  host,
		//	"input_type": common.InputTypeExcel,
		//}
		//result, err := lgc.CoreAPI.ApiServer().UpdateHost(context.Background(), header, params)
		//if err != nil {
		//	blog.Errorf("update host http request  error: %v, rid: %s", err, util.GetHTTPCCRequestID(header))
		//	return returnByErrCode(defErr, common.CCErrCommHTTPDoRequestFailed, nil, "")
		//}
		//
		//errData, exist := result.Data.Get("error")
		//if exist && errData != nil {
		//	errMsgs = append(errMsgs, defLang.Languagef("import_update_data_fail", rowNum, errData))
		//	continue
		//}

		successMsgs = append(successMsgs, strconv.Itoa(rowNum))
	}

	result.Data["success"] = successMsgs
	result.Data["error"] = errMsgs
	if (len(asstObjectUniqueIDMap) == 0 && objectUniqueID == 0) || len(errMsgs) != 0 {
		return result
	}

	if _, exist := f.Sheet["association"]; !exist {
		return result
	}
	row, err := f.Sheet["association"].Row(common.HostAddMethodExcelAssociationIndexOffset)
	if err != nil {
		blog.Errorf("get sheet association failed, err: %v, rid: %s", err, util.GetHTTPCCRequestID(header))
		return returnByErrCode(defErr, common.CCErrCommHTTPDoRequestFailed, nil, err.Error())
	}
	if row.GetCellCount() < 2 {
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
func (lgc *Logics) GetHostBizRelations(ctx context.Context, header http.Header,
	hostInfos map[int]map[string]interface{}, bizID int64) (map[int64]int64, error) {
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
			blog.Errorf(" GetHostBizRelations failed, GetHostModuleRelation err:%v, params: %#v, rid: %s", err, params,
				rid)
			return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !resp.Result {
			blog.Errorf(" GetHostBizRelations failed, GetHostModuleRelation resp:%#v, params: %#v, rid: %s", resp,
				params, rid)
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
		blog.Errorf(" getExistHostsByHostIDs failed, ListHostWithoutApp resp:%#v, option: %d, rid: %s", resp, option,
			rid)
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

// handleHostInfo handle host info to export host
func (lgc *Logics) handleHostInfo(ctx context.Context, header http.Header, fields map[string]Property, objIDs []string,
	input *metadata.ExcelExportHostInput) ([]mapstr.MapStr, error) {

	rid := util.GetHTTPCCRequestID(header)
	defLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))

	hostFields := make([]string, 0)
	for _, property := range fields {
		hostFields = append(hostFields, property.ID)
	}

	hostInfo, err := lgc.GetHostData(input.AppID, input.HostIDArr, hostFields, input.ExportCond, header, defLang)
	if err != nil {
		blog.Errorf("get host info failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	if len(hostInfo) == 0 {
		blog.Errorf("not find host, host id: %v, cond: %#v, rid: %s", input.HostIDArr, input.ExportCond, rid)
		return nil, nil
	}

	if err := lgc.handleModule(hostInfo, rid); err != nil {
		blog.Errorf("add module name to host failed, err: %v, rid: %s", err, rid)
		return nil, err
	}
	setDIs, hostSetMap, err := lgc.handleSet(hostInfo, rid)
	if err != nil {
		blog.Errorf("add set name to host failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	if len(objIDs) > 0 {
		setParentIDs, setCustomMap, err := lgc.getSetParentID(ctx, header, setDIs, rid)
		if err != nil {
			blog.Errorf("get set parent id and host set rel map failed, err: %v, rid: %s", err, rid)
			return nil, err
		}

		err = lgc.handleCustomData(ctx, header, hostInfo, objIDs, rid, setParentIDs, setCustomMap, hostSetMap)
		if err != nil {
			blog.Errorf("get custom parent id and host custom rel map failed, err: %v, rid: %s", err, rid)
			return nil, err
		}
	}

	return hostInfo, nil
}

// handleModule 处理module数据
func (lgc *Logics) handleModule(hostInfo []mapstr.MapStr, rid string) error {
	// 统计host与module关系
	for _, data := range hostInfo {
		moduleMap, exist := data[common.BKInnerObjIDModule].([]interface{})
		if !exist {
			blog.Errorf("get module map data from host data failed, not exist, data: %#v, rid: %s", data, rid)
			return fmt.Errorf("from host data get module map, not exist, rid: %s", rid)
		}

		moduleNameMap := make(map[string]int)
		for idx, module := range moduleMap {
			rowMap, err := mapstr.NewFromInterface(module)
			if err != nil {
				blog.Errorf("get module data from host data failed, err: %v, rid: %s", err, rid)
				return err
			}

			moduleName, err := rowMap.String(common.BKModuleNameField)
			if err != nil {
				blog.Errorf("get module name from host data failed, err: %v, rid: %s", err, rid)
				return fmt.Errorf("from host data get module name, not exist, rid: %s", rid)
			}
			moduleNameMap[moduleName] = idx
		}

		var moduleStr string
		for moduleName := range moduleNameMap {
			if moduleStr == "" {
				moduleStr = moduleName
			} else {
				moduleStr += "," + moduleName
			}
		}
		data.Set("modules", moduleStr)
	}

	return nil
}

// handleSet 处理set数据
func (lgc *Logics) handleSet(hostInfo []mapstr.MapStr, rid string) ([]int64, map[int64][]int64, error) {
	// 统计host与set关系
	hostSetMap := make(map[int64][]int64, 0)
	setIDs := make([]int64, 0)
	header := util.BuildHeader(common.CCSystemOperatorUserName, common.BKDefaultOwnerID)
	res, err := lgc.CoreAPI.CoreService().System().SearchPlatformSetting(context.Background(), header)
	if err != nil {
		return nil, nil, err
	}
	conf := res.Data

	for _, data := range hostInfo {
		setMap, exist := data[common.BKInnerObjIDSet].([]interface{})
		if !exist {
			blog.Errorf("get set map data from host data, not exist, data: %#v, rid: %s", data, rid)
			return nil, nil, fmt.Errorf("from host data get set map, not exist, rid: %s", rid)
		}

		rowMap, err := mapstr.NewFromInterface(data[common.BKInnerObjIDHost])
		if err != nil {
			blog.Errorf("get host map data failed, hostData: %#v, err: %v, rid: %s", data, err, rid)
			return nil, nil, err
		}

		hostID, err := rowMap.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("get host id failed, host id: %s, err: %v, rid: %s", hostID, err, rid)
			return nil, nil, nil
		}

		setNameMap := make(map[string]int)
		setSubIDs := make([]int64, 0)
		for idx, set := range setMap {
			rowMap, err := mapstr.NewFromInterface(set)
			if err != nil {
				blog.Errorf("get set data from host data failed, err: %v, rid: %s", err, rid)
				return nil, nil, err
			}

			setName, err := rowMap.String(common.BKSetNameField)
			if err != nil {
				blog.Errorf("get set name from host data failed, err: %v, rid: %s", err, rid)
				return nil, nil, fmt.Errorf("from host data get set name, not exist, rid: %s", rid)
			}
			setNameMap[setName] = idx

			setID, err := rowMap.Int64(common.BKSetIDField)
			if err != nil {
				blog.Errorf("get set id from host data failed, err: %v, rid: %s", err, rid)
				return nil, nil, err
			}

			if setName != string(conf.BuiltInSetName) {
				setIDs = append(setIDs, setID)
				setSubIDs = append(setSubIDs, setID)
			}
		}

		hostSetMap[hostID] = setSubIDs

		var setStr string
		for setName := range setNameMap {
			if setStr == "" {
				setStr = setName
			} else {
				setStr += "," + setName
			}
		}
		data.Set("sets", setStr)
	}

	return setIDs, hostSetMap, nil
}

// getSetParentID get set parent id and set custom rel map
func (lgc *Logics) getSetParentID(ctx context.Context, header http.Header, setIDs []int64, rid string) ([]int64,
	map[int64]int64, error) {
	// 获取set数据，统计set parent id
	querySet := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKSetIDField: mapstr.MapStr{
				common.BKDBIN: setIDs,
			},
		},
		Fields: []string{common.BKSetIDField, common.BKInstParentStr, common.BKSetNameField},
	}

	sets, err := lgc.Engine.CoreAPI.ApiServer().ReadInstance(ctx, header, common.BKInnerObjIDSet, querySet)
	if err != nil {
		blog.Errorf("get set data failed, cond: %#v, err: %v,rid:%s", querySet, err, rid)
		return nil, nil, err
	}
	if !sets.Result {
		blog.Errorf("get sets failed, err code: %d, err msg: %s, rid: %s", sets.Code, sets.ErrMsg, rid)
		return nil, nil, fmt.Errorf("get sets failed, err msg: %s", sets.ErrMsg)
	}

	setParentIDs := make([]int64, 0)
	setCustomMap := make(map[int64]int64, 0)
	for _, set := range sets.Data.Info {
		parentID, err := set.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("get set parent id failed, err: %v, rid: %s", err, rid)
			return nil, nil, err
		}
		setParentIDs = append(setParentIDs, parentID)

		setID, err := set.Int64(common.BKSetIDField)
		if err != nil {
			blog.Errorf("get set id failed, err: %v, rid: %s", err, rid)
			return nil, nil, err
		}
		setCustomMap[setID] = parentID
	}

	return setParentIDs, setCustomMap, nil
}

// handleCustomData 处理自定义成层级数据
func (lgc *Logics) handleCustomData(ctx context.Context, header http.Header, hostInfo []mapstr.MapStr, objIDs []string,
	rid string, parentIDs []int64, setCustomMap map[int64]int64, hostSetMap map[int64][]int64) error {
	instIdParentIDMap := make(map[int64]int64, 0)
	instIdNameMap := make(map[int64]string, 0)
	var err error
	for _, objID := range objIDs {
		parentIDs, instIdParentIDMap, instIdNameMap, err = lgc.getCustomData(ctx, header, parentIDs, objID, rid)
		if err != nil {
			blog.Errorf("get custom data failed, cond: %#v, err: %v, rid: %s", parentIDs, err, rid)
			return err
		}

		hostCustomMap := make(map[int64][]int64, 0)
		hostCustomNameMap := make(map[int64]string, 0)
		for hostID, setIDs := range hostSetMap {
			customNameMap := make(map[string]int, 0)
			for idx, setID := range setIDs {
				customNameMap[instIdNameMap[setCustomMap[setID]]] = idx
				hostCustomMap[hostID] = append(hostCustomMap[hostID], setCustomMap[setID])
			}

			customStr := ""
			for customName := range customNameMap {
				if customStr == "" {
					customStr = customName
				} else {
					customStr += "," + customName
				}
			}

			hostCustomNameMap[hostID] = customStr
		}

		for _, data := range hostInfo {
			rowMap, err := mapstr.NewFromInterface(data[common.BKInnerObjIDHost])
			if err != nil {
				blog.Errorf("get host map data failed, hostData: %#v, err: %v, rid: %s", data, err, rid)
				return err
			}

			hostID, err := rowMap.Int64(common.BKHostIDField)
			if err != nil {
				blog.Errorf("get host id failed, host id: %s, err: %v, rid: %s", hostID, err, rid)
				return err
			}

			data[objID] = hostCustomNameMap[hostID]
		}

		setCustomMap = instIdParentIDMap
		hostSetMap = hostCustomMap
	}

	return nil
}

func (lgc *Logics) getCustomData(ctx context.Context, header http.Header, instIDs []int64, objID, rid string) ([]int64,
	map[int64]int64, map[int64]string, error) {
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKInstIDField: mapstr.MapStr{
				common.BKDBIN: instIDs,
			},
		},
		Fields: []string{common.BKInstIDField, common.BKInstNameField, common.BKInstParentStr},
	}

	insts, err := lgc.Engine.CoreAPI.ApiServer().ReadInstance(ctx, header, objID, query)
	if err != nil {
		blog.Errorf("get custom level inst data failed, query cond: %#v, err: %v, rid: %s", query, err, rid)
		return nil, nil, nil, err
	}

	parentIDs := make([]int64, 0)
	instIdParentIdMap := make(map[int64]int64, 0)
	instIdNameMap := make(map[int64]string, 0)
	for _, inst := range insts.Data.Info {
		parentID, err := inst.Int64(common.BKParentIDField)
		if err != nil {
			blog.Errorf("get inst parent id failed, err: %v, rid: %s", err, rid)
			return nil, nil, nil, err
		}
		parentIDs = append(parentIDs, parentID)

		instID, err := inst.Int64(common.BKInstIDField)
		if err != nil {
			blog.Errorf("get inst id failed, err: %v, rid: %s", err, rid)
			return nil, nil, nil, err
		}
		instIdParentIdMap[instID] = parentID

		instName, err := inst.String(common.BKInstNameField)
		if err != nil {
			blog.Errorf("get inst name failed, err: %v, rid: %s", err, rid)
			return nil, nil, nil, err
		}
		instIdNameMap[instID] = instName
	}

	return parentIDs, instIdParentIdMap, instIdNameMap, nil
}
