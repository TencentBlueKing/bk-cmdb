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

	"github.com/rentiansheng/xlsx"
)

// GetImportInsts get insts from excel file
func (lgc *Logics) GetImportInsts(f *xlsx.File, objID string, header http.Header, headerRow int, isInst bool, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	fields, err := lgc.GetObjFieldIDs(objID, nil, header)
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
func (lgc *Logics) GetInstData(ownerID, objID, instIDStr string, header http.Header, kvMap mapstr.MapStr) ([]mapstr.MapStr, error) {

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

	result, err := lgc.Engine.CoreAPI.ApiServer().GetInstDetail(context.Background(), header, ownerID, objID, searchCond)
	if nil != err || !result.Result {
		return nil, errors.New(result.ErrMsg)
	}

	if 0 == result.Data.Count {
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
