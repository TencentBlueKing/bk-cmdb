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
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/rentiansheng/xlsx"
)

// GetImportInsts get insts from excel file
func (lgc *Logics) GetImportInsts(ctx context.Context, f *xlsx.File, objID string, header http.Header, headerRow int,
	isInst bool, defLang lang.DefaultCCLanguageIf, modelBizID int64) (map[int]map[string]interface{}, []string, error) {

	rid := util.ExtractRequestIDFromContext(ctx)

	fields, err := lgc.GetObjFieldIDs(objID, nil, nil, header, modelBizID, common.HostAddMethodExcelDefaultIndex)

	if nil != err {
		return nil, nil, errors.New(defLang.Languagef("web_get_object_field_failure", err.Error()))
	}
	if len(f.Sheets) == 0 {
		blog.Errorf("the excel file sheets is empty, rid: %s", rid)
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}
	sheet := f.Sheets[0]
	if sheet == nil {
		blog.Errorf("import object %s instance, but the excel file sheet is empty, rid: %s", objID, rid)
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}

	departmentMap, err := lgc.getDepartmentMap(ctx, header)
	if err != nil {
		blog.Errorf("get department failed, err: %v, rid: %s", err, rid)
		return nil, nil, err
	}

	if isInst {
		return GetExcelData(ctx, sheet, fields, common.KvMap{"import_from": common.HostAddMethodExcel}, true, headerRow,
			defLang, departmentMap)
	} else {
		return GetRawExcelData(ctx, sheet, common.KvMap{"import_from": common.HostAddMethodExcel}, headerRow, defLang,
			departmentMap)
	}
}

// GetInstData TODO
func (lgc *Logics) GetInstData(objID string, instIDArr []int64, header http.Header) ([]mapstr.MapStr, error) {
	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	searchCond := mapstr.MapStr{}

	searchCond["fields"] = []string{}
	searchCond["condition"] = mapstr.MapStr{
		common.BKInstIDField: mapstr.MapStr{
			common.BKDBIN: instIDArr,
		},
		common.BKObjIDField: objID,
	}
	searchCond["page"] = nil
	result, err := lgc.Engine.CoreAPI.ApiServer().GetInstDetail(context.Background(), header, objID, searchCond)
	if nil != err {
		blog.Errorf("get inst data detail error:%v , search condition:%#v, rid: %s", err, searchCond, rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("get inst data detail error:%v , search condition:%#v, rid: %s", result.ErrMsg, searchCond, rid)
		return nil, defErr.Error(result.Code)
	}

	if result.Data.Count == 0 {
		blog.Errorf("get inst data detail, but got 0 instances , search condition:%#v, rid: %s", searchCond, rid)
		return nil, defErr.Error(common.CCErrAPINoObjectInstancesIsFound)
	}

	return result.Data.Info, nil
}

// ImportInsts import host info
func (lgc *Logics) ImportInsts(ctx context.Context, f *xlsx.File, objID string, header http.Header,
	defLang lang.DefaultCCLanguageIf, modelBizID int64, opType int64,
	asstObjectUniqueIDMap map[string]int64, objectUniqueID int64) (
	resultData mapstr.MapStr, errCode int, err error) {

	rid := util.GetHTTPCCRequestID(header)

	if opType == 1 {
		if _, exist := f.Sheet["association"]; !exist {
			return nil, 0, nil
		}
		info, err := lgc.importStatisticsAssociation(ctx, header, objID, f.Sheet["association"])
		if err != nil {
			blog.Errorf("ImportHosts failed, GetImportHosts error:%s, rid: %s", err.Error(), rid)
			return nil, err.GetCode(), err
		}
		return mapstr.MapStr{"association": info}, 0, nil
	}

	return lgc.importInsts(ctx, f, objID, header, defLang, modelBizID, asstObjectUniqueIDMap, objectUniqueID)
}

// importInsts import insts info
func (lgc *Logics) importInsts(ctx context.Context, f *xlsx.File, objID string, header http.Header,
	defLang lang.DefaultCCLanguageIf, modelBizID int64, asstObjectUniqueIDMap map[string]int64, objectUniqueID int64) (
	resultData mapstr.MapStr, errCode int, err error) {

	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	resultData = mapstr.New()

	insts, errMsg, err := lgc.GetImportInsts(ctx, f, objID, header, 0, true, defLang, modelBizID)
	if err != nil {
		blog.Errorf("get %s inst info from excel error, err: %v, rid: %s", objID, err, rid)
		return resultData, common.CCErrWebFileContentFail, err
	}
	if len(errMsg) != 0 {
		resultData.Set("err", errMsg)
		return resultData, common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail,
			strings.Join(errMsg, ","))
	}

	var resultErr error
	result := &metadata.ResponseDataMapStr{}
	result.BaseResp.Result = true

	if len(insts) != 0 {
		params := mapstr.MapStr{}
		params["input_type"] = common.InputTypeExcel
		params["BatchInfo"] = insts
		params[common.BKAppIDField] = modelBizID
		result, resultErr = lgc.CoreAPI.ApiServer().AddInstByImport(context.Background(), header,
			util.GetOwnerID(header), objID, params)
		if resultErr != nil {
			blog.Errorf("ImportInsts add inst info  http request  err: %v, rid: %s", resultErr, rid)
			return nil, common.CCErrorUnknownOrUnrecognizedError, resultErr
		}
		resultData.Merge(result.Data)
	}

	resp := &metadata.ResponseDataMapStr{Data: mapstr.New()}
	resp = lgc.handleExcelAssociation(ctx, header, f, objID, rid, asstObjectUniqueIDMap, objectUniqueID, defLang, resp)
	resultData.Merge(resp.Data)

	return
}
