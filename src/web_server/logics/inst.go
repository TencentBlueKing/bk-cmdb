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
func (lgc *Logics) GetImportInsts(f *xlsx.File, objID string, header http.Header, headerRow int, isInst bool, defLang lang.DefaultCCLanguageIf, meta metadata.Metadata) (map[int]map[string]interface{}, []string, error) {

	fields, err := lgc.GetObjFieldIDs(objID, nil, nil, header, meta)
	if nil != err {
		return nil, nil, errors.New(defLang.Languagef("web_get_object_field_failure", err.Error()))
	}
	if 0 == len(f.Sheets) {
		blog.Error("the excel file sheets is empty")
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}
	sheet := f.Sheets[0]
	if nil == sheet {
		blog.Error("the excel fiel sheet is nil")
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}
	if isInst {
		return GetExcelData(sheet, fields, common.KvMap{"import_from": common.HostAddMethodExcel}, true, headerRow, defLang)
	} else {
		return GetRawExcelData(sheet, common.KvMap{"import_from": common.HostAddMethodExcel}, headerRow, defLang)
	}
}

//GetInstData get inst data
func (lgc *Logics) GetInstData(ownerID, objID, instIDStr string, header http.Header, kvMap mapstr.MapStr, meta metadata.Metadata) ([]mapstr.MapStr, error) {

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
		common.BKOwnerIDField: ownerID,
		common.BKObjIDField:   objID,
	}
	searchCond["page"] = nil
	searchCond[metadata.BKMetadata] = meta
	result, err := lgc.Engine.CoreAPI.ApiServer().GetInstDetail(context.Background(), header, ownerID, objID, searchCond)
	if nil != err || !result.Result {
		blog.Errorf("get inst detail error:%v , search condition:%#v", err, searchCond)
		return nil, errors.New(result.ErrMsg)
	}

	if 0 == result.Data.Count {
		blog.Errorf("inst inst count is 0 ")
		return nil, errors.New("no inst")
	}

	// read object attributes
	attrCond := mapstr.MapStr{}
	attrCond[common.BKObjIDField] = objID
	attrCond[common.BKOwnerIDField] = ownerID
	attrResult, aErr := lgc.Engine.CoreAPI.ApiServer().GetObjectAttr(context.Background(), header, attrCond)
	if nil != aErr || !attrResult.Result {
		blog.Errorf("get object attr error: %s", aErr.Error())
		return nil, errors.New(result.ErrMsg)
	}
	for _, cell := range attrResult.Data {
		kvMap.Set(cell.PropertyID, cell.PropertyName)

	}
	return result.Data.Info, nil
}

// ImportHosts import host info
func (lgc *Logics) ImportInsts(ctx context.Context, f *xlsx.File, objID string, header http.Header, defLang lang.DefaultCCLanguageIf, meta metadata.Metadata) (resultData mapstr.MapStr, errCode int, err error) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	resultData = mapstr.New()
	insts, errMsg, err := lgc.GetImportInsts(f, objID, header, 0, true, defLang, meta)
	if nil != err {
		blog.Errorf("ImportInsts  get %s inst info from excel error, error:%s logID:%s", objID, err.Error(), util.GetHTTPCCRequestID(header))
		return
	}
	if 0 != len(errMsg) {
		resultData.Set("err", errMsg)
		return resultData, common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, " file empty")
	}
	if 0 == len(insts) {
		return nil, common.CCErrWebFileContentEmpty, defErr.Errorf(common.CCErrWebFileContentEmpty, "")
	}

	params := mapstr.MapStr{}
	params["input_type"] = common.InputTypeExcel
	params["BatchInfo"] = insts
	result, resultErr := lgc.CoreAPI.ApiServer().AddInst(context.Background(), header, util.GetOwnerID(header), objID, params)
	if nil != err {
		blog.Errorf("ImportInsts add inst info  http request  error:%s, rid:%s", resultErr.Error(), util.GetHTTPCCRequestID(header))
		return nil, common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	resultData.Merge(result.Data)
	if !result.Result {
		errCode = result.Code
		err = defErr.New(result.Code, result.ErrMsg)
	}

	if len(f.Sheets) > 2 {
		asstInfoMap := GetAssociationExcelData(f.Sheets[1], common.HostAddMethodExcelAssociationIndexOffset)

		if len(asstInfoMap) > 0 {
			asstInfoMapInput := &metadata.RequestImportAssociation{
				AssociationInfoMap: asstInfoMap,
			}
			asstResult, asstResultErr := lgc.CoreAPI.ApiServer().ImportAssociation(ctx, header, objID, asstInfoMapInput)
			if nil != asstResultErr {
				blog.Errorf("ImportHosts logics http request import %s association error:%s, rid:%s", objID, asstResultErr.Error(), util.GetHTTPCCRequestID(header))
				return nil, common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
			}
			resultData.Set("asst_error", asstResult.Data.ErrMsgMap)
			if result.Result && !asstResult.Result {
				errCode = asstResult.Code
				err = defErr.New(asstResult.Code, asstResult.ErrMsg)
			}
		}
	}

	return

}
