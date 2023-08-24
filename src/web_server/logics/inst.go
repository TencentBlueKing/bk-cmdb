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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/tealeg/xlsx/v3"
)

// GetImportAttr get import object attribute from excel file
func (lgc *Logics) GetImportAttr(ctx context.Context, f *xlsx.File, objID string, headerRow int,
	defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	rid := util.ExtractRequestIDFromContext(ctx)
	if len(f.Sheets) == 0 {
		blog.Errorf("the excel file sheets is empty, rid: %s", rid)
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}
	sheet := f.Sheets[0]
	if sheet == nil {
		blog.Errorf("import object %s instance, but the excel file sheet is empty, rid: %s", objID, rid)
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}

	return GetRawExcelData(ctx, sheet, common.KvMap{"import_from": common.HostAddMethodExcel}, headerRow, defLang)
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
	resultData = mapstr.New()

	preData, err := lgc.getImportExcelPreData(objID, header, f, defLang, modelBizID)
	if err != nil {
		blog.Errorf("get import pre data failed, err: %v, rid: %s", err, rid)
		return nil, common.CCErrWebFileContentFail, err
	}

	var successMsgs []string
	var errMsgs []string
	for i := 0; i < len(preData.DataRange); i++ {
		rowNum := preData.DataRange[i].Start + 1
		inst, errMsg, err := GetExcelData(ctx, preData, i, i+1, common.KvMap{"import_from": common.HostAddMethodExcel},
			defLang)
		if err != nil {
			blog.Errorf("get %s inst info from excel error, err: %v, rid: %s", objID, err, rid)
			errMsgs = append(errMsgs, defLang.Languagef("import_data_fail", rowNum, err.Error()))
			continue
		}
		if len(errMsg) != 0 {
			errMsgs = append(errMsgs, defLang.Languagef("import_data_fail", rowNum, errMsg[0]))
			continue
		}

		inst, err = lgc.handleImportEnumQuoteInst(ctx, header, inst, preData.Fields, rid)
		if err != nil {
			blog.Errorf("handle enum quote inst failed, err: %v, rid: %s", err, rid)
			errMsgs = append(errMsgs, defLang.Languagef("import_data_fail", rowNum, err.Error()))
			continue
		}

		if len(inst) == 0 {
			continue
		}
		params := mapstr.MapStr{}
		params["input_type"] = common.InputTypeExcel
		params["BatchInfo"] = inst
		params[common.BKAppIDField] = modelBizID
		//result, err := lgc.CoreAPI.ApiServer().AddInstByImport(context.Background(), header,
		//	util.GetOwnerID(header), objID, params)
		//if err != nil {
		//	blog.Errorf("ImportInsts add inst info  http request  err: %v, rid: %s", err, rid)
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
	resultData["success"] = successMsgs
	resultData["error"] = errMsgs
	if len(errMsgs) != 0 {
		return
	}

	resp := &metadata.ResponseDataMapStr{Data: mapstr.New()}
	resp = lgc.handleExcelAssociation(ctx, header, f, objID, rid, asstObjectUniqueIDMap, objectUniqueID, defLang, resp)
	resultData.Merge(resp.Data)

	return
}

// handleImportEnumQuoteInst search inst detail and return a id-bk_inst_name map
func (lgc *Logics) handleImportEnumQuoteInst(c context.Context, h http.Header, data map[int]map[string]interface{},
	fields map[string]Property, rid string) (map[int]map[string]interface{}, error) {

	enumQuoteFields := make(map[string]Property, 0)
	for id, property := range fields {
		if property.PropertyType == common.FieldTypeEnumQuote {
			enumQuoteFields[id] = property
		}
	}
	if len(enumQuoteFields) == 0 {
		return data, nil
	}

	for _, rowMap := range data {
		for id, property := range enumQuoteFields {
			enumQuoteNameList, exist := rowMap[id]
			if !exist || enumQuoteNameList == nil {
				continue
			}

			quoteObjID, err := GetEnumQuoteObjID(property.Option, rid)
			if err != nil {
				blog.Errorf("get enum quote option obj id failed, err: %s, rid: %s", err, rid)
				return nil, err
			}

			enumQuoteIDs, err := lgc.getEnumQuoteIds(c, h, quoteObjID, rid, enumQuoteNameList)
			if err != nil {
				blog.Errorf("get enum quote id list failed, err: %v, rid: %s", err, rid)
				return nil, err
			}
			rowMap[id] = enumQuoteIDs
		}
	}

	return data, nil
}

// GetEnumQuoteObjID get enum quote field option bk_obj_id and bk_inst_id value
func GetEnumQuoteObjID(option interface{}, rid string) (string, error) {
	var quoteObjID string
	if option == nil {
		return quoteObjID, fmt.Errorf("enum quote option is nil")
	}
	arrOption, ok := option.([]interface{})
	if !ok {
		blog.Errorf("option %v not enum quote option, rid: %s", option, rid)
		return quoteObjID, fmt.Errorf("enum quote option is unvalid")
	}

	for _, o := range arrOption {
		mapOption, ok := o.(map[string]interface{})
		if !ok || mapOption == nil {
			blog.Errorf("option %v not enum quote option, enum quote option item must bk_obj_id, rid: %s", option,
				rid)
			return quoteObjID, fmt.Errorf("convert option map[string]interface{} failed")
		}
		objIDVal, objIDOk := mapOption["bk_obj_id"]
		if !objIDOk || objIDVal == "" {
			blog.Errorf("enum quote option bk_obj_id can't be empty, rid: %s", option, rid)
			return quoteObjID, fmt.Errorf("enum quote option bk_obj_id can't be empty")
		}
		objID, ok := objIDVal.(string)
		if !ok {
			blog.Errorf("objIDVal %v not string, rid: %s", objIDVal, rid)
			return quoteObjID, fmt.Errorf("enum quote option bk_obj_id is not string")
		}

		if quoteObjID == "" {
			quoteObjID = objID
		} else if quoteObjID != objID {
			return quoteObjID, fmt.Errorf("enum quote objID not unique, objID: %s", objID)
		}
	}

	return quoteObjID, nil
}

// getEnumQuoteIds search inst detail and return a inst id list
func (lgc *Logics) getEnumQuoteIds(c context.Context, h http.Header, objID, rid string,
	enumQuoteNameList interface{}) ([]int64, error) {

	input := &metadata.QueryCondition{
		Fields: []string{common.GetInstIDField(objID)},
		Condition: mapstr.MapStr{
			common.GetInstNameField(objID): mapstr.MapStr{common.BKDBIN: enumQuoteNameList},
		},
		DisableCounter: true,
	}
	resp, err := lgc.Engine.CoreAPI.ApiServer().ReadInstance(c, h, objID, input)
	if err != nil {
		blog.Errorf("get quote inst name list failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	enumQuoteIDs := make([]int64, 0)
	for _, info := range resp.Data.Info {
		enumQuoteID, err := info.Int64(common.GetInstIDField(objID))
		if err != nil {
			blog.Errorf("get enum quote id failed, err: %v, rid: %s", err, rid)
			continue
		}
		enumQuoteIDs = append(enumQuoteIDs, enumQuoteID)
	}

	return enumQuoteIDs, nil
}
