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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/rentiansheng/xlsx"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	webCommon "configcenter/src/web_server/common"
)

//GetImportInsts get insts from excel file
func GetImportInsts(f *xlsx.File, objID, url string, header http.Header, headerRow int, isInst bool, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	fields, err := GetObjFieldIDs(objID, url, nil, header)
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
func GetInstData(ownerID, objID, instIDStr, apiAddr string, header http.Header, kvMap map[string]string) ([]interface{}, error) {

	instInfo := make([]interface{}, 0)
	sInstCond := make(map[string]interface{})
	instIDArr := strings.Split(instIDStr, ",")

	iInstIDArr := make([]int, 0)
	for _, j := range instIDArr {
		instID, _ := strconv.Atoi(j)
		iInstIDArr = append(iInstIDArr, instID)
	}

	// construct the search condition

	sInstCond["fields"] = []string{}
	sInstCond["condition"] = map[string]interface{}{
		common.BKInstIDField: map[string]interface{}{
			"$in": iInstIDArr,
		},
		common.BKOwnerIDField: ownerID,
		common.BKObjIDField:   objID,
	}
	sInstCond["page"] = nil

	// read insts
	url := apiAddr + fmt.Sprintf("/api/%s/inst/search/owner/%s/object/%s/detail", webCommon.API_VERSION, ownerID, objID)
	result, _ := httpRequest(url, sInstCond, header)
	blog.Info("search inst  url:%s", url)
	blog.Info("search inst  return:%s", result)
	js, _ := simplejson.NewJson([]byte(result))
	instData, _ := js.Map()
	instResult := instData["result"].(bool)
	if !instResult {
		return nil, errors.New(instData["bk_error_msg"].(string))
	}

	instDataArr := instData["data"].(map[string]interface{})
	instInfo = instDataArr["info"].([]interface{})
	instCnt, _ := instDataArr["count"].(json.Number).Int64()
	if !instResult || 0 == instCnt {
		return instInfo, errors.New("no inst")
	}

	// read object attributes
	url = apiAddr + fmt.Sprintf("/api/%s/object/attr/search", webCommon.API_VERSION)
	attrCond := make(map[string]interface{})
	attrCond[common.BKObjIDField] = objID
	attrCond[common.BKOwnerIDField] = ownerID
	result, _ = httpRequest(url, attrCond, header)
	blog.Info("get inst attr  url:%s", url)
	blog.Info("get inst attr return:%s", result)
	js, _ = simplejson.NewJson([]byte(result))
	instAttr, _ := js.Map()
	attrData := instAttr["data"].([]interface{})
	for _, j := range attrData {
		cell := j.(map[string]interface{})
		key := cell[common.BKPropertyIDField].(string)
		value, ok := cell[common.BKPropertyNameField].(string)
		if ok {
			kvMap[key] = value
		} else {
			kvMap[key] = ""
		}

	}
	return instInfo, nil
}
