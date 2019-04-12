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
	"configcenter/src/common/language"
	meta "configcenter/src/common/metadata"
	webCommon "configcenter/src/web_server/common"
)

func GetImportNetProperty(
	header http.Header, defLang language.DefaultCCLanguageIf, f *xlsx.File, url string) (map[int]map[string]interface{}, []string, error) {

	if 0 == len(f.Sheets) {
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}

	fields := GetNetPropertyField(defLang)

	sheet := f.Sheets[0]
	if nil == sheet {
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}

	return GetExcelData(sheet, fields, nil, true, 0, defLang)
}

func BuildNetPropertyExcelFromData(defLang language.DefaultCCLanguageIf, fields map[string]Property, data []interface{}, sheet *xlsx.Sheet) error {
	productExcelHealer(fields, nil, sheet, defLang)

	rowIndex := common.HostAddMethodExcelIndexOffset
	for _, row := range data {
		propertyData, ok := row.(map[string]interface{})
		if !ok {
			msg := fmt.Sprintf("[Export Net Property] Build NetProperty excel from data, convert to map[string]interface{} fail, data: %v", row)
			blog.Errorf(msg)
			return errors.New(msg)
		}

		setExcelRowDataByIndex(propertyData, sheet, rowIndex, fields)
		rowIndex++
	}

	return nil
}

// get net property data to export
func GetNetPropertyData(header http.Header, apiAddr, netPropertyIDStr string) ([]interface{}, error) {
	netPropertyIDStrArr := strings.Split(netPropertyIDStr, ",")
	netPropertyIDArr := []int64{}

	for _, netPropertyIDStr := range netPropertyIDStrArr {
		netPropertyID, _ := strconv.ParseInt(netPropertyIDStr, 10, 64)
		netPropertyIDArr = append(netPropertyIDArr, netPropertyID)
	}

	netPropertyCond := map[string]interface{}{
		"field": []string{},
		"condition": []map[string]interface{}{
			map[string]interface{}{
				"field":    common.BKNetcollectPropertyIDField,
				"operator": common.BKDBIN,
				"value":    netPropertyIDArr,
			},
		},
	}

	url := apiAddr + fmt.Sprintf("/api/%s/collector/netcollect/property/action/search", webCommon.API_VERSION)
	result, err := httpRequest(url, netPropertyCond, header)
	if nil != err {
		blog.Errorf("[Export Net Property] http request error:%v", err)
	}

	blog.V(5).Infof("[Export Net Property] search netProperty url:%s", url)
	blog.V(5).Infof("[Export Net Property] search netProperty return:%s", result)

	js, err := simplejson.NewJson([]byte(result))
	if nil != err {
		blog.Errorf("[Export Net Property] convert http reponse string [%s] to json error:%v", result, err)
	}
	netPropertyDataResult, err := js.Map()
	if nil != err {
		blog.Errorf("[Export Net Property] convert http reponse json [%#+v] to map[string]interface{} error:%v", netPropertyDataResult, err)
	}

	netPropertyResult, ok := netPropertyDataResult["result"].(bool)
	if !ok {
		blog.Errorf("[Export Net Property] http reponse 'result'[%#+v] is bool", netPropertyDataResult["result"])
	}
	if !netPropertyResult {
		return nil, errors.New(netPropertyDataResult["bk_error_msg"].(string))
	}

	netPropertyData, ok := netPropertyDataResult["data"].(map[string]interface{})
	if !ok {
		blog.Errorf("[Export Net Property] http reponse 'data'[%#+v] is not map[string]interface{}", netPropertyDataResult["data"])
	}
	netPropertyInfo, ok := netPropertyData["info"].([]interface{})
	if !ok {
		blog.Errorf("[Export Net Property] http reponse 'info'[%#+v] is not []interface{}", netPropertyData["info"])
	}
	_, ok = netPropertyData["count"].(json.Number)
	if !ok {
		blog.Errorf("[Export Net Property] http reponse 'count'[%#+v] is not a number", netPropertyData["count"])
	}
	netPropertyCount, err := netPropertyData["count"].(json.Number).Int64()
	if nil != err {
		blog.Errorf("[Export Net Property] http reponse 'count'[%#+v] convert to int64 error:%v", netPropertyData["count"], err)
	}

	if 0 == netPropertyCount {
		return netPropertyInfo, errors.New("no netProperty")
	}

	blog.V(5).Infof("[Export Net Property] search return netProperty info:%s", netPropertyInfo)
	return netPropertyInfo, nil
}

//BuildNetPropertyExcelTemplate  return httpcode, error
func BuildNetPropertyExcelTemplate(header http.Header, defLang language.DefaultCCLanguageIf, url, filename string) error {
	var file *xlsx.File
	file = xlsx.NewFile()

	sheet, err := file.AddSheet(common.BKNetProperty)
	if nil != err {
		blog.Errorf("[Build NetProperty Excel Template] add comment sheet error, sheet name:%s, error:%s", common.BKNetProperty, err.Error())
		return err
	}

	fields := GetNetPropertyField(defLang)

	blog.V(5).Infof("[Build NetProperty Excel Template]  fields count:%d", len(fields))

	productExcelHealer(fields, nil, sheet, defLang)

	if err = file.Save(filename); nil != err {
		return err
	}

	return nil
}

// get feild to import property or generate template
func GetNetPropertyField(lang language.DefaultCCLanguageIf) map[string]Property {

	return map[string]Property{
		common.BKPropertyNameField: Property{
			Name: lang.Language("import_property_comment_property_name"), ID: common.BKPropertyNameField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 0, IsRequire: true,
		},
		common.BKDeviceNameField: Property{
			Name: lang.Language("import_property_comment_device_name"), ID: common.BKDeviceNameField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 1, IsRequire: true,
		},
		common.BKOIDField: Property{
			Name: lang.Language("import_property_comment_oid"), ID: common.BKOIDField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 2, IsRequire: true,
		},
		common.BKPeriodField: Property{
			Name: lang.Language("import_property_comment_period"), ID: common.BKPeriodField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 3, IsRequire: false,
		},
		common.BKActionField: Property{
			Name: lang.Language("import_property_comment_action"), ID: common.BKActionField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 4, IsRequire: false,
		},
	}
}

// add extra feild to export property
func AddNetPropertyExtFields(originField map[string]Property, lang language.DefaultCCLanguageIf) {

	field := map[string]Property{
		common.BKNetcollectPropertyIDField: Property{
			Name:         lang.Language("import_property_comment_net_property_id"),
			ID:           common.BKNetcollectPropertyIDField,
			PropertyType: common.FieldTypeInt,
		},
		common.BKPropertyIDField: Property{
			Name:         lang.Language("import_property_comment_property_id"),
			ID:           common.BKPropertyIDField,
			PropertyType: common.FieldTypeSingleChar,
		},
		common.BKDeviceIDField: Property{
			Name:         lang.Language("import_property_comment_device_id"),
			ID:           common.BKDeviceIDField,
			PropertyType: common.FieldTypeInt,
		},
		common.BKObjIDField: Property{
			Name:         lang.Language("import_property_comment_object_id"),
			ID:           common.BKObjIDField,
			PropertyType: common.FieldTypeSingleChar,
		},
		common.BKObjNameField: Property{
			Name:         lang.Language("import_property_comment_object_name"),
			ID:           common.BKObjNameField,
			PropertyType: common.FieldTypeSingleChar,
		},
		meta.AttributeFieldUnit: Property{
			Name:         lang.Language("import_property_comment_unit"),
			ID:           meta.AttributeFieldUnit,
			PropertyType: common.FieldTypeSingleChar,
		},
	}

	originFieldLen := len(originField)

	for key, value := range field {
		value.ExcelColIndex = originFieldLen
		originFieldLen++
		originField[key] = value
	}
}
