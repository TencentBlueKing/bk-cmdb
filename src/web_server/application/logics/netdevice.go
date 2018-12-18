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
	webCommon "configcenter/src/web_server/common"
)

// get date from excel file to import device
func GetImportNetDevices(
	header http.Header, defLang language.DefaultCCLanguageIf, f *xlsx.File, url string) (map[int]map[string]interface{}, []string, error) {

	if 0 == len(f.Sheets) {
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}

	fields := GetNetDevicefield(defLang)

	sheet := f.Sheets[0]
	if nil == sheet {
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}

	return GetExcelData(sheet, fields, nil, true, 0, defLang)
}

func BuildNetDeviceExcelFromData(defLang language.DefaultCCLanguageIf, fields map[string]Property, data []interface{}, sheet *xlsx.Sheet) error {
	productExcelHealer(fields, nil, sheet, defLang)

	rowIndex := common.HostAddMethodExcelIndexOffset
	for _, row := range data {
		deviceData, ok := row.(map[string]interface{})
		if !ok {
			msg := fmt.Sprintf("[Export Net Device] Build NetDevice excel from data, convert to map[string]interface{} fail, data: %v", row)
			blog.Errorf(msg)
			return errors.New(msg)
		}

		setExcelRowDataByIndex(deviceData, sheet, rowIndex, fields)
		rowIndex++
	}

	return nil
}

// get net device data to export
func GetNetDeviceData(header http.Header, apiAddr, deviceIDStr string) ([]interface{}, error) {
	deviceIDStrArr := strings.Split(deviceIDStr, ",")
	deviceIDArr := []int64{}

	for _, deviceIDStr := range deviceIDStrArr {
		deviceID, _ := strconv.ParseInt(deviceIDStr, 10, 64)
		deviceIDArr = append(deviceIDArr, deviceID)
	}

	deviceCond := map[string]interface{}{
		"field": []string{},
		"condition": []map[string]interface{}{
			map[string]interface{}{
				"field":    common.BKDeviceIDField,
				"operator": common.BKDBIN,
				"value":    deviceIDArr,
			},
		},
	}

	url := apiAddr + fmt.Sprintf("/api/%s/collector/netcollect/device/action/search", webCommon.API_VERSION)
	result, err := httpRequest(url, deviceCond, header)
	if nil != err {
		blog.Errorf("[Export Net Device] http request error:%v", err)
	}

	blog.V(5).Infof("[Export Net Device] search device url:%s", url)
	blog.V(5).Infof("[Export Net Device] search device return:%s", result)

	js, err := simplejson.NewJson([]byte(result))
	if nil != err {
		blog.Errorf("[Export Net Device] convert http reponse string [%s] to json error:%v", result, err)
	}
	deviceDataResult, err := js.Map()
	if nil != err {
		blog.Errorf("[Export Net Device] convert http reponse json [%#+v] to map[string]interface{} error:%v", deviceDataResult, err)
	}

	deviceResult, ok := deviceDataResult["result"].(bool)
	if !ok {
		blog.Errorf("[Export Net Device] http reponse 'result'[%#+v] is bool", deviceDataResult["result"])
	}
	if !deviceResult {
		return nil, errors.New(deviceDataResult["bk_error_msg"].(string))
	}

	deviceData, ok := deviceDataResult["data"].(map[string]interface{})
	if !ok {
		blog.Errorf("[Export Net Device] http reponse 'data'[%#+v] is not map[string]interface{}", deviceDataResult["data"])
	}
	deviceInfo, ok := deviceData["info"].([]interface{})
	if !ok {
		blog.Errorf("[Export Net Device] http reponse 'info'[%#+v] is not []interface{}", deviceData["info"])
	}
	_, ok = deviceData["count"].(json.Number)
	if !ok {
		blog.Errorf("[Export Net Device] http reponse 'count'[%#+v] is not a number", deviceData["count"])
	}
	deviceCount, err := deviceData["count"].(json.Number).Int64()
	if nil != err {
		blog.Errorf("[Export Net Device] http reponse 'count'[%#+v] convert to int64 error:%v", deviceData["count"], err)
	}

	if 0 == deviceCount {
		return deviceInfo, errors.New("no device")
	}

	blog.V(5).Infof("[Export Net Device] search return device info:%s", deviceInfo)
	return deviceInfo, nil
}

//BuildNetDeviceExcelTemplate  return httpcode, error
func BuildNetDeviceExcelTemplate(header http.Header, defLang language.DefaultCCLanguageIf, url, filename string) error {
	var file *xlsx.File
	file = xlsx.NewFile()

	sheet, err := file.AddSheet(common.BKNetDevice)
	if nil != err {
		blog.Errorf("[Build NetDevice Excel Template] add comment sheet error, sheet name:%s, error:%s", common.BKNetDevice, err.Error())
		return err
	}

	fields := GetNetDevicefield(defLang)

	blog.V(5).Infof("[Build NetDevice Excel Template] fields count:%d", len(fields))

	productExcelHealer(fields, nil, sheet, defLang)

	if err = file.Save(filename); nil != err {
		return err
	}

	return nil
}

// get feild to import device or generate template
func GetNetDevicefield(lang language.DefaultCCLanguageIf) map[string]Property {

	return map[string]Property{
		common.BKDeviceNameField: Property{
			Name: lang.Language("import_device_comment_device_name"), ID: common.BKDeviceNameField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 0, IsRequire: true,
		},
		common.BKDeviceModelField: Property{
			Name: lang.Language("import_device_comment_device_model"), ID: common.BKDeviceModelField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 1, IsRequire: true,
		},
		common.BKObjNameField: Property{
			Name: lang.Language("import_device_comment_obj_name"), ID: common.BKObjNameField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 2, IsRequire: true,
		},
		common.BKVendorField: Property{
			Name: lang.Language("import_device_comment_vendor"), ID: common.BKVendorField,
			PropertyType: common.FieldTypeSingleChar, ExcelColIndex: 3, IsRequire: true,
		},
	}
}

// add extra feild to export device
func AddNetDeviceExtFields(originField map[string]Property, lang language.DefaultCCLanguageIf) {

	field := map[string]Property{
		common.BKDeviceIDField: Property{
			Name:         lang.Language("import_device_comment_device_id"),
			ID:           common.BKDeviceIDField,
			PropertyType: common.FieldTypeInt,
		},
		common.BKObjIDField: Property{
			Name:         lang.Language("import_device_comment_obj_id"),
			ID:           common.BKObjIDField,
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
