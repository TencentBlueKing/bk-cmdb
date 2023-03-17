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
	"strconv"
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
func (lgc *Logics) GetImportInsts(ctx context.Context, f *xlsx.File, objID string, req *http.Request, headerRow int,
	isInst bool, defLang lang.DefaultCCLanguageIf, modelBizID int64) (map[int]map[string]interface{}, []string, error) {

	rid := util.ExtractRequestIDFromContext(ctx)

	fields, err := lgc.GetObjFieldIDs(objID, nil, nil, req.Header, modelBizID,
		common.HostAddMethodExcelDefaultIndex)

	if nil != err {
		return nil, nil, errors.New(defLang.Languagef("web_get_object_field_failure", err.Error()))
	}
	if len(f.Sheets) == 0 {
		blog.Errorf("the excel file sheets is empty, rid: %s", rid)
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}
	sheet := f.Sheets[0]
	if nil == sheet {
		blog.Errorf("import object %s instance, but the excel file sheet is empty, rid: %s", objID, rid)
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}

	if isInst {
		return GetExcelData(ctx, sheet, fields, common.KvMap{"import_from": common.HostAddMethodExcel}, true,
		headerRow, defLang)
	} else {
		return GetRawExcelData(ctx, sheet, common.KvMap{"import_from": common.HostAddMethodExcel}, headerRow, defLang)
	}
}

func (lgc *Logics) GetInstData(ownerID, objID, instIDStr string, header http.Header, kvMap mapstr.MapStr) ([]mapstr.MapStr, error) {
	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	instIDArr := strings.Split(instIDStr, ",")
	searchCond := mapstr.MapStr{}

	iInstIDArr := make([]int, 0)
	for _, j := range instIDArr {
		instID, _ := strconv.Atoi(j)
		iInstIDArr = append(iInstIDArr, instID)
	}

	searchCond["fields"] = []string{}
	searchCond["condition"] = mapstr.MapStr{
		common.BKInstIDField: mapstr.MapStr{
			common.BKDBIN: iInstIDArr,
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

	if 0 == result.Data.Count {
		blog.Errorf("get inst data detail, but got 0 instances , search condition:%#v, rid: %s", searchCond, rid)
		return nil, defErr.Error(common.CCErrAPINoObjectInstancesIsFound)
	}

	// read object attributes
	attrCond := mapstr.MapStr{}
	attrCond[common.BKObjIDField] = objID
	attrCond[common.BKOwnerIDField] = ownerID
	attrResult, aErr := lgc.Engine.CoreAPI.ApiServer().GetObjectAttr(context.Background(), header, attrCond)
	if nil != aErr {
		blog.Errorf("get object: %s instance, but get object attr error: %v, rid: %s", objID, aErr, rid)
		return nil, defErr.Error(common.CCErrTopoObjectAttributeSelectFailed)
	}

	if !attrResult.Result {
		blog.Errorf("get object: %s instance, but get object attr error: %s, rid: %s", objID, attrResult.Code, rid)
		return nil, defErr.Error(common.CCErrTopoObjectAttributeSelectFailed)
	}

	for _, cell := range attrResult.Data {
		kvMap.Set(cell.PropertyID, cell.PropertyName)
	}

	return result.Data.Info, nil
}

// ImportHosts import host info
func (lgc *Logics) ImportInsts(ctx context.Context, f *xlsx.File, objID string, req *http.Request,
	defLang lang.DefaultCCLanguageIf, modelBizID int64) (resultData mapstr.MapStr, errCode int, err error) {

	rid := util.GetHTTPCCRequestID(req.Header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Header))
	resultData = mapstr.New()

	insts, errMsg, err := lgc.GetImportInsts(ctx, f, objID, req, 0, true, defLang, modelBizID)
	if err != nil {
		blog.Errorf("get %s inst info from excel error, err: %v, rid: %s", objID, err, rid)
		return
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
		result, resultErr = lgc.CoreAPI.ApiServer().AddInst(context.Background(), req.Header,
			util.GetOwnerID(req.Header), objID, params)
		if err != nil {
			blog.Errorf("add inst info failed, err: %v, rid: %s", resultErr, rid)
			return nil, common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		resultData.Merge(result.Data)
		if !result.Result {
			errCode = result.Code
			err = defErr.New(result.Code, result.ErrMsg)
		}

	}

	for _, sheet := range f.Sheets {
		if sheet.Name != "association" {
			continue
		}

		asstInfoMap, errMsg := getAssociationExcelData(sheet, common.HostAddMethodExcelAssociationIndexOffset,
			defLang)

		if len(asstInfoMap) > 0 {
			asstInfoMapInput := &metadata.RequestImportAssociation{
				AssociationInfoMap: asstInfoMap,
			}
			asstResult, asstResultErr := lgc.CoreAPI.ApiServer().ImportAssociation(ctx, req.Header, objID,
				asstInfoMapInput)
			if asstResultErr != nil {
				blog.Errorf("import %s association failed, err: %v, rid:%s", objID, asstResultErr, rid)
				return nil, common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
			}
			errMsg = append(errMsg, asstResult.Data.ErrMsgMap...)
			if errCode == 0 && !asstResult.Result {
				errCode = asstResult.Code
				err = defErr.New(asstResult.Code, asstResult.ErrMsg)
			}
		}

		resultData.Set("asst_error", errMsg)
	}

	return
}
