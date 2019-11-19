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
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/rentiansheng/xlsx"
)

func GetImportNetProperty(
	header http.Header, defLang language.DefaultCCLanguageIf, f *xlsx.File) (map[int]map[string]interface{}, []string, error) {
	ctx := util.NewContextFromHTTPHeader(header)

	if 0 == len(f.Sheets) {
		return nil, nil, errors.New(defLang.Language("web_excel_content_empty"))
	}

	fields := GetNetPropertyField(defLang)

	sheet := f.Sheets[0]
	if nil == sheet {
		return nil, nil, errors.New(defLang.Language("web_excel_sheet_not_found"))
	}

	return GetExcelData(ctx, sheet, fields, nil, true, 0, defLang)
}

func BuildNetPropertyExcelFromData(ctx context.Context, defLang language.DefaultCCLanguageIf, fields map[string]Property, data []mapstr.MapStr, sheet *xlsx.Sheet) error {
	sortedFields := SortByIsRequired(fields)
	productExcelHeader(ctx, sortedFields, nil, sheet, defLang)

	rowIndex := common.HostAddMethodExcelIndexOffset
	for _, row := range data {
		propertyData := row

		setExcelRowDataByIndex(propertyData, sheet, rowIndex, sortedFields)
		rowIndex++
	}

	return nil
}

// get net property data to export
func (lgc *Logics) GetNetPropertyData(header http.Header, netPropertyIDStr string) ([]mapstr.MapStr, error) {
	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	netPropertyIDStrArr := strings.Split(netPropertyIDStr, ",")
	netPropertyIDArr := []int64{}

	for _, netPropertyIDStr := range netPropertyIDStrArr {
		netPropertyID, _ := strconv.ParseInt(netPropertyIDStr, 10, 64)
		netPropertyIDArr = append(netPropertyIDArr, netPropertyID)
	}

	condItem := condition.ConditionItem{
		Field:    common.BKNetcollectPropertyIDField,
		Operator: common.BKDBIN,
		Value:    netPropertyIDArr,
	}
	searchCond := condition.CreateCondition()
	searchCond.AddConditionItem(condItem)

	propertyResult, err := lgc.Engine.CoreAPI.ApiServer().SearchNetCollectDevice(context.Background(), header, searchCond)
	if nil != err {
		blog.Errorf("search net property data inst  error:%#v , search condition:%#v, rid: %s", err, searchCond, rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if 0 == propertyResult.Data.Count {
		return propertyResult.Data.Info, errors.New("no device")
	}

	blog.V(5).Infof("[Export Net Device Property] search return device info:%s, rid: %s", propertyResult, rid)
	return propertyResult.Data.Info, nil
}

// BuildNetPropertyExcelTemplate  return httpcode, error
func BuildNetPropertyExcelTemplate(header http.Header, defLang language.DefaultCCLanguageIf, filename string) error {
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)
	var file *xlsx.File
	file = xlsx.NewFile()

	sheet, err := file.AddSheet(common.BKNetProperty)
	if nil != err {
		blog.Errorf("[Build NetProperty Excel Template] add comment sheet error, sheet name:%s, error:%s, rid: %s", common.BKNetProperty, err.Error(), rid)
		return err
	}

	fields := GetNetPropertyField(defLang)

	blog.V(5).Infof("[Build NetProperty Excel Template]  fields count:%d, rid: %s", len(fields), rid)
	sortedFields := SortByIsRequired(fields)
	productExcelHeader(ctx, sortedFields, nil, sheet, defLang)

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
