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
	"configcenter/src/common/condition"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"

	"github.com/rentiansheng/xlsx"
)

// get date from excel file to import device
func GetImportNetDevices(
	header http.Header, defLang language.DefaultCCLanguageIf, f *xlsx.File) (map[int]map[string]interface{}, []string, error) {
	ctx := util.NewContextFromHTTPHeader(header)

	if 0 == len(f.Sheets) {
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}

	fields := GetNetDevicefield(defLang)

	sheet := f.Sheets[0]
	if nil == sheet {
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}

	return GetExcelData(ctx, sheet, fields, nil, true, 0, defLang)
}

func BuildNetDeviceExcelFromData(ctx context.Context, defLang language.DefaultCCLanguageIf, fields map[string]Property, data []mapstr.MapStr, sheet *xlsx.Sheet) error {
	sortedFields := SortByIsRequired(fields)
	productExcelHeader(ctx, sortedFields, nil, sheet, defLang)

	rowIndex := common.HostAddMethodExcelIndexOffset
	for _, row := range data {
		deviceData := row

		setExcelRowDataByIndex(deviceData, sheet, rowIndex, sortedFields)
		rowIndex++
	}

	return nil
}

// get net device data to export
func (lgc *Logics) GetNetDeviceData(header http.Header, deviceIDStr string) ([]mapstr.MapStr, error) {
	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	deviceIDStrArr := strings.Split(deviceIDStr, ",")
	deviceIDArr := []int64{}

	for _, deviceIDStr := range deviceIDStrArr {
		deviceID, _ := strconv.ParseInt(deviceIDStr, 10, 64)
		deviceIDArr = append(deviceIDArr, deviceID)
	}

	condItem := condition.ConditionItem{
		Field:    common.BKDeviceIDField,
		Operator: common.BKDBIN,
		Value:    deviceIDArr,
	}
	searchCond := condition.CreateCondition()
	searchCond.AddConditionItem(condItem)
	deviceResult, err := lgc.Engine.CoreAPI.ApiServer().SearchNetCollectDevice(context.Background(), header, searchCond)
	if nil != err {
		blog.Errorf("search net device data inst  error:%#v , search condition:%#v, rid: %s", err, searchCond, rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if 0 == deviceResult.Data.Count {
		return deviceResult.Data.Info, errors.New("no device")
	}

	blog.V(5).Infof("[Export Net Device] search return device info:%s, rid: %s", deviceResult, rid)
	return deviceResult.Data.Info, nil
}

// BuildNetDeviceExcelTemplate  return httpcode, error
func BuildNetDeviceExcelTemplate(header http.Header, defLang language.DefaultCCLanguageIf, filename string) error {
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)
	var file *xlsx.File
	file = xlsx.NewFile()

	sheet, err := file.AddSheet(common.BKNetDevice)
	if nil != err {
		blog.Errorf("[Build NetDevice Excel Template] add comment sheet error, sheet name:%s, error:%s, rid: %s", common.BKNetDevice, err.Error(), rid)
		return err
	}

	fields := GetNetDevicefield(defLang)

	blog.V(5).Infof("[Build NetDevice Excel Template] fields count:%d, rid: %s", len(fields), rid)

	productExcelHealer(ctx, fields, nil, sheet, defLang)

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
