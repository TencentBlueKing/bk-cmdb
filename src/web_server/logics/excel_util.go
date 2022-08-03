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
	"reflect"
	"regexp"
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

const (
	userBracketsPattern         = `\([a-zA-Z0-9\@\p{Han} .,_-]*\)`
	organizationBracketsPattern = `\[(\d+)\]([^\s]+)`
)

var (
	headerRow          = common.HostAddMethodExcelIndexOffset
	userBracketsRegexp = regexp.MustCompile(userBracketsPattern)
	orgBracketsRegexp  = regexp.MustCompile(organizationBracketsPattern)
)

// getFilterFields 不需要展示字段
func getFilterFields(objID string) []string {
	switch objID {
	case common.BKInnerObjIDHost:
		return []string{"bk_agent_status", "bk_agent_version", "bk_set_name", "bk_module_name", "bk_biz_name"}
	default:
		return []string{"create_time"}
	}
}

func getCustomFields(filterFields []string, customFields []string) []string {
	customFieldsList := make([]string, 0)

	for _, fieldID := range customFields {
		if util.InStrArr(filterFields, fieldID) || "" == fieldID {
			continue
		}
		customFieldsList = append(customFieldsList, fieldID)
	}
	return customFieldsList
}

// checkExcelHeader check whether invalid fields exists in header and return headers
func checkExcelHeader(ctx context.Context, sheet *xlsx.Sheet, fields map[string]Property, isCheckHeader bool, defLang lang.DefaultCCLanguageIf) (map[int]string, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// rowLen := len(sheet.Rows[headerRow-1].Cells)
	var errCells []string
	ret := make(map[int]string)
	if headerRow > len(sheet.Rows) {
		return ret, errors.New(defLang.Language("web_excel_not_data"))
	}
	if headerRow+common.ExcelImportMaxRow < len(sheet.Rows) {
		return ret, errors.New(defLang.Languagef("web_excel_import_too_much", common.ExcelImportMaxRow))
	}
	for index, name := range sheet.Rows[headerRow-1].Cells {
		strName := name.Value
		// skip the ignored cell field
		if strName == common.ExcelCellIgnoreValue {
			continue
		}
		field, ok := fields[strName]
		if true == ok {
			field.ExcelColIndex = index
			fields[strName] = field
		} else {
			errCells = append(errCells, strName)
		}
		ret[index] = strName
	}

	if len(sheet.Rows[headerRow-1].Cells) < 2 && true == isCheckHeader {
		blog.Errorf("err:%s, no found fields %s, rid:%s", defLang.Language("web_import_field_not_found"), strings.Join(errCells, ","), rid)
		return ret, errors.New(defLang.Language("web_import_field_not_found"))
	}
	return ret, nil

}

// setExcelRowDataByIndex insert  map[string]interface{}  to excel row by index,
// mapHeaderIndex:Correspondence between head and field
// fields each field description,  field type, isrequire, validate role
func setExcelRowDataByIndex(rowMap mapstr.MapStr, sheet *xlsx.Sheet, rowIndex int, fields map[string]Property) {

	// 非模型字段导出是没有field中没有ID 字段，因为导入的时候，第二行是作为Property
	for id, property := range fields {
		val, ok := rowMap[id]
		if false == ok {
			continue
		}
		if property.NotExport {
			continue
		}

		cell := sheet.Cell(rowIndex, property.ExcelColIndex)
		// cell.NumFmt = "@"

		switch property.PropertyType {
		case common.FieldTypeEnum:
			var cellVal string
			arrVal, ok := property.Option.([]interface{})
			strEnumID, enumIDOk := val.(string)
			if true == ok && true == enumIDOk {
				cellVal = getEnumNameByID(strEnumID, arrVal)
				cell.SetString(cellVal)
			}

		case common.FieldTypeBool:
			bl, ok := val.(bool)
			if ok {
				if bl {
					cell.SetValue(fieldTypeBoolTrue)
				} else {
					cell.SetValue(fieldTypeBoolFalse)
				}

			}

		case common.FieldTypeInt:
			intVal, err := util.GetInt64ByInterface(val)
			if nil == err {
				cell.SetInt64(intVal)
			}

		case common.FieldTypeFloat:
			floatVal, err := util.GetFloat64ByInterface(val)
			if nil == err {
				cell.SetFloat(floatVal)
			}

		default:
			switch val.(type) {
			case string:
				strVal := val.(string)
				if "" != strVal {
					cell.SetString(val.(string))
				}

			default:
				cell.SetValue(val)
			}
		}

	}

	return

}

func getDataFromByExcelRow(ctx context.Context, row *xlsx.Row, rowIndex int, fields map[string]Property,
	defFields common.KvMap, nameIndexMap map[int]string, defLang lang.DefaultCCLanguageIf,
	department map[int64]metadata.DepartmentItem) (map[string]interface{}, []string) {

	rid := util.ExtractRequestIDFromContext(ctx)
	result := make(map[string]interface{})
	errMsg := make([]string, 0)
	for cellIndex, cell := range row.Cells {
		fieldName, ok := nameIndexMap[cellIndex]
		if !ok || strings.Trim(fieldName, "") == "" || cell.Value == "" {
			continue
		}

		switch cell.Type() {
		case xlsx.CellTypeString:
			result[fieldName] = strings.TrimSpace(cell.String())
		case xlsx.CellTypeStringFormula:
			result[fieldName] = strings.TrimSpace(cell.String())
		case xlsx.CellTypeNumeric:
			cellValue, err := cell.Float()
			if err != nil {
				errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, (rowIndex+1)))
				blog.Errorf("%d row %s column get content err: %v, rid: %s", rowIndex+1, fieldName, err, rid)
				continue
			}
			result[fieldName] = cellValue
		case xlsx.CellTypeBool:
			cellValue := cell.Bool()
			result[fieldName] = cellValue
		case xlsx.CellTypeDate:
			cellValue, err := cell.GetTime(true)
			if err != nil {
				errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", errMsg, fieldName,
					(rowIndex+1)))
				blog.Errorf("%d row %s column get content error:%s, rid: %s", rowIndex+1, fieldName, err, rid)
				continue
			}
			result[fieldName] = cellValue
		default:
			errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, (rowIndex+1)))
			blog.Errorf("unknown the type, %v,   %v, rid: %s", reflect.TypeOf(cell), cell.Type(), rid)
			continue
		}

		field, ok := fields[fieldName]
		if !ok {
			blog.Errorf("%d row %s field not found , rid: %s", rowIndex+1, fieldName, rid)
			continue
		}

		result, errMsg = buildAttrByPropertyType(rid, fieldName, cell.Value, rowIndex, field, result, department,
			defLang, errMsg)
	}
	if len(errMsg) != 0 {
		return nil, errMsg
	}
	if len(result) == 0 {
		return result, nil
	}
	for k, v := range defFields {
		result[k] = v
	}

	return result, nil

}

func buildAttrByPropertyType(rid, fieldName, cellValue string, rowIndex int, field Property,
	result map[string]interface{}, department map[int64]metadata.DepartmentItem, defLang lang.DefaultCCLanguageIf,
	errMsg []string) (map[string]interface{}, []string) {

	switch field.PropertyType {
	case common.FieldTypeBool:
		switch result[fieldName].(type) {
		case bool:
		default:
			if bl, err := strconv.ParseBool(cellValue); err == nil {
				result[fieldName] = bl
			}
		}
	case common.FieldTypeEnum:
		if option, optionOk := field.Option.([]interface{}); optionOk {
			result[fieldName] = getEnumIDByName(cellValue, option)
		}
	case common.FieldTypeInt:
		// convertor int not err, set field value to correct type
		if intVal, err := util.GetInt64ByInterface(result[fieldName]); err != nil {
			blog.Errorf("get excel cell value error, field:%s, value:%s, err: %v, rid: %s", fieldName,
				result[fieldName], err, rid)
		} else {
			result[fieldName] = intVal
		}
	case common.FieldTypeFloat:
		if floatVal, err := util.GetFloat64ByInterface(result[fieldName]); err == nil {
			result[fieldName] = floatVal
		} else {
			blog.Errorf("get excel cell value failed, field:%s, value:%s, err:%v, rid: %s", fieldName,
				result[fieldName], err, rid)
		}
	case common.FieldTypeOrganization:
		result, errMsg = checkOrgnization(result, department, rowIndex, defLang, errMsg, fieldName, rid)
	case common.FieldTypeUser:
		// convert userNames,  eg: " admin(admin),xiaoming(小明 ),leo(li hong),  " => "admin,xiaoming,leo"
		userNames := util.GetStrByInterface(result[fieldName])
		userNames = userBracketsRegexp.ReplaceAllString(userNames, "")
		userNames = strings.Trim(strings.Trim(userNames, " "), ",")
		result[fieldName] = userNames
	default:
		if util.IsStrProperty(field.PropertyType) {
			result[fieldName] = strings.TrimSpace(cellValue)
		}
	}

	return result, errMsg
}

func checkOrgnization(result map[string]interface{}, department map[int64]metadata.DepartmentItem, rowIndex int,
	defLang lang.DefaultCCLanguageIf, errMsg []string, fieldName, rid string) (map[string]interface{}, []string) {

	if len(department) == 0 {
		blog.Debug("no department in paas, rid: %s", rid)
		errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, rowIndex+1)+
			defLang.Languagef("nonexistent_org"))
		return result, errMsg
	}
	// convert Organization,  eg: "[1]总公司,[2]分公司" => "1,2"
	orgStr := util.GetStrByInterface(result[fieldName])
	if len(orgStr) <= 0 {
		blog.Debug("get excel cell value failed, field:%s, value:%s, err:%v, rid: %s", fieldName,
			result[fieldName], "not a valid organization type", rid)
		errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, rowIndex+1)+
			defLang.Languagef("organization_type_invalid"))
		return result, errMsg
	}
	orgItems := strings.Split(orgStr, ",")
	org := make([]int64, len(orgItems))
	for i, v := range orgItems {
		var err error
		orgID := orgBracketsRegexp.FindStringSubmatch(v)
		if len(orgID) != 3 {
			blog.Errorf("regular matching is empty, please enter the correct content, field: %s, value: %s, rid: %s",
				fieldName, result[fieldName], rid)
			errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, rowIndex+1)+
				defLang.Languagef("organization_type_invalid"))
			break
		}

		if org[i], err = strconv.ParseInt(orgID[1], 10, 64); err != nil {
			blog.Debug("get excel cell value error, field: %s, value: %s, err: %v, rid: %s", fieldName,
				result[fieldName], "not a valid organization type", rid)
			errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, rowIndex+1)+
				defLang.Languagef("organization_type_invalid"))
			break
		}

		dp, exist := department[org[i]]
		if !exist {
			blog.Debug("get excel cell value error, field:%s, value:%s, err:%v, rid: %s", fieldName,
				result[fieldName], "organization does not exist", rid)
			errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, rowIndex+1)+
				defLang.Languagef("nonexistent_org"))
			break
		}

		if dp.Name != orgID[2] && dp.FullName != orgID[2] {
			blog.Debug("get excel cell value error, field:%s, value:%s, err:%v, rid: %s", fieldName,
				result[fieldName], "organization name or full_name does not match", rid)
			errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, rowIndex+1)+
				defLang.Languagef("organization_type_invalid"))
			break
		}
	}
	result[fieldName] = org
	return result, errMsg
}

// productExcelHeader Excel文件头部，
func productExcelHeader(ctx context.Context, fields map[string]Property, filter []string, xlsxFile *xlsx.File,
	sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) {
	rid := util.ExtractRequestIDFromContext(ctx)
	styleCell := getHeaderCellGeneralStyle()
	// 橙棕色
	cellStyle := getCellStyle(common.ExcelFirstColumnCellColor, common.ExcelHeaderFirstRowFontColor)
	// 粉色
	colStyle := getCellStyle(common.ExcelHeaderFirstColumnColor, common.ExcelHeaderFirstRowFontColor)
	sheet.Col(0).Width = 18

	firstColFields := []string{common.ExcelFirstColumnFieldName, common.ExcelFirstColumnFieldType,
		common.ExcelFirstColumnFieldID, common.ExcelFirstColumnInstData}
	for index, field := range firstColFields {
		cellName := sheet.Cell(index, 0)
		fieldName := defLang.Language(field)
		cellName.Value = fieldName
		cellName.SetStyle(cellStyle)
	}

	// 给第一列剩下的空格设置颜色
	for i := 3; i < 1003; i++ {
		cellName := sheet.Cell(i, 0)
		cellName.SetStyle(colStyle)
	}

	handleFieldParam := &HandleFieldParam{
		Rid:       rid,
		StyleCell: styleCell,
		Sheet:     sheet,
		File:      xlsxFile,
		Filter:    filter,
		DefLang:   defLang,
		CellStyle: cellStyle,
		ColStyle:  colStyle,
	}

	for _, field := range fields {
		handleField(field, handleFieldParam)
	}
}

// productHostExcelHeader Excel文件头部，
func productHostExcelHeader(ctx context.Context, fields map[string]Property, filter []string, xlsxFile *xlsx.File,
	sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf, objName, cloudAreaName []string) {

	rid := util.ExtractRequestIDFromContext(ctx)
	styleCell := getHeaderCellGeneralStyle()
	// 橙棕色
	cellStyle := getCellStyle(common.ExcelFirstColumnCellColor, common.ExcelHeaderFirstRowFontColor)
	// 粉色
	colStyle := getCellStyle(common.ExcelHeaderFirstColumnColor, common.ExcelHeaderFirstRowFontColor)
	sheet.Col(0).Width = 18
	// 字典中的值为国际化之后的"业务拓扑"和"业务名"，"集群"，”模块“，用来做判断，命中即变化相应的cell颜色。
	bizTopoMap := map[string]int{
		defLang.Language("web_ext_field_topo"):        1,
		defLang.Language("biz_property_bk_biz_name"):  1,
		defLang.Language("web_ext_field_module_name"): 1,
		defLang.Language("web_ext_field_set_name"):    1,
	}
	for _, name := range objName {
		bizTopoMap[name] = 1
	}

	firstColFields := []string{common.ExcelFirstColumnFieldName, common.ExcelFirstColumnFieldType,
		common.ExcelFirstColumnFieldID, common.ExcelFirstColumnInstData}

	for index, field := range firstColFields {
		cellName := sheet.Cell(index, 0)
		fieldName := defLang.Language(field)
		cellName.Value = fieldName
		cellName.SetStyle(cellStyle)
	}

	// 给第一列剩下的空格设置颜色
	for i := 3; i < 1000; i++ {
		cellName := sheet.Cell(i, 0)
		cellName.SetStyle(colStyle)
	}

	handleFieldParam := &HandleFieldParam{
		Rid:       rid,
		StyleCell: styleCell,
		File:      xlsxFile,
		Sheet:     sheet,
		Filter:    filter,
		DefLang:   defLang,
		CellStyle: cellStyle,
		ColStyle:  colStyle,
	}

	for _, field := range fields {
		handleHostField(field, handleFieldParam, cloudAreaName, bizTopoMap)
	}
}

func handleHostField(field Property, handleFieldParam *HandleFieldParam, cloudAreaName []string,
	bizTopoMap map[string]int) {

	isRequire := ""
	if field.IsRequire {
		isRequire = handleFieldParam.DefLang.Language("web_excel_header_required")
	}
	// 主机部分特殊逻辑 针对云区域与topo进行处理
	if field.ID == common.BKCloudIDField {
		// 设置属性的id
		handleFieldParam.Sheet.Cell(2, field.ExcelColIndex).Value = field.ID
		handleFieldParam.Sheet.Col(field.ExcelColIndex).Width = 18
		handleFieldParam.Sheet.Cell(0, field.ExcelColIndex).Value = field.Name + isRequire
		// 设置单元格颜色
		handleFieldParam.Sheet.Cell(0, field.ExcelColIndex).SetStyle(getHeaderFirstRowCellStyle(field.IsRequire))
		setExcelCellIgnored(handleFieldParam.Sheet, handleFieldParam.StyleCell, 1, field.ExcelColIndex)
		handleFieldParam.Sheet.Cell(2, field.ExcelColIndex).SetStyle(handleFieldParam.StyleCell)

		// 设置云区域的下拉选项
		if len(cloudAreaName) != 0 {
			enumSheet, err := handleFieldParam.File.AddSheet(field.Name)
			if err != nil {
				blog.Errorf("add enum sheet failed, err: %s, rid: %s", err, handleFieldParam.Rid)
				return
			}
			for _, enum := range cloudAreaName {
				enumSheet.AddRow().AddCell().SetString(enum)
			}
			dd := xlsx.NewXlsxCellDataValidation(true, true, true)
			if err := dd.SetInFileList(field.Name, 0, 0, 0, len(cloudAreaName)-1); err != nil {
				blog.Errorf("SetDropList failed, err: %+v, rid: %s", err, handleFieldParam.Rid)
			}
			handleFieldParam.Sheet.Col(field.ExcelColIndex).SetDataValidationWithStart(dd,
				common.HostAddMethodExcelIndexOffset)
			handleFieldParam.Sheet.Col(field.ExcelColIndex).SetType(xlsx.CellTypeString)
		}
		return
	}

	if _, ok := bizTopoMap[field.Name]; ok {
		handleFieldParam.Sheet.Col(field.ExcelColIndex).Width = 18

		handleFieldParam.Sheet.Cell(0, field.ExcelColIndex).Value = field.Name + isRequire
		handleFieldParam.Sheet.Cell(0, field.ExcelColIndex).SetStyle(handleFieldParam.CellStyle)
		setExcelCellIgnored(handleFieldParam.Sheet, handleFieldParam.CellStyle, 1, field.ExcelColIndex)
		setExcelCellIgnored(handleFieldParam.Sheet, handleFieldParam.CellStyle, 2, field.ExcelColIndex)

		// 给业务拓扑和业务列剩下的空格设置颜色
		for i := 3; i < 1003; i++ {
			handleFieldParam.Sheet.Cell(i, field.ExcelColIndex).SetStyle(handleFieldParam.ColStyle)
		}
		handleFieldParam.Sheet.Col(field.ExcelColIndex).SetType(xlsx.CellTypeString)
		return
	}

	// 处理其他的通用属性逻辑
	handleField(field, handleFieldParam)
}

func handleField(field Property, handleFieldParam *HandleFieldParam) {
	index := field.ExcelColIndex
	handleFieldParam.Sheet.Col(index).Width = 18
	fieldTypeName, skip := getPropertyTypeAliasName(field.PropertyType, handleFieldParam.DefLang)
	if skip || field.NotExport {
		// 不需要用户输入的类型continue
		handleFieldParam.Sheet.Col(index).Hidden = true
		return
	}
	isRequire := ""

	if field.IsRequire {
		// "(必填)"
		isRequire = handleFieldParam.DefLang.Language("web_excel_header_required")
	}
	if util.Contains(handleFieldParam.Filter, field.ID) {
		return
	}
	cellName := handleFieldParam.Sheet.Cell(0, index)
	cellName.Value = field.Name + isRequire
	cellName.SetStyle(getHeaderFirstRowCellStyle(field.IsRequire))

	cellType := handleFieldParam.Sheet.Cell(1, index)
	cellType.Value = fieldTypeName
	cellType.SetStyle(handleFieldParam.StyleCell)

	cellEnName := handleFieldParam.Sheet.Cell(2, index)
	cellEnName.Value = field.ID
	cellEnName.SetStyle(handleFieldParam.StyleCell)

	switch field.PropertyType {
	case common.FieldTypeInt:
		handleFieldParam.Sheet.Col(index).SetType(xlsx.CellTypeNumeric)
	case common.FieldTypeFloat:
		handleFieldParam.Sheet.Col(index).SetType(xlsx.CellTypeNumeric)
	case common.FieldTypeEnum:
		optionArr, ok := field.Option.([]interface{})

		if ok {

			enumSheet, err := handleFieldParam.File.AddSheet(field.Name)
			if err != nil {
				blog.Errorf("add enum sheet failed, err: %s, rid: %s", err, handleFieldParam.Rid)
			}

			for _, enum := range getEnumNames(optionArr) {
				enumSheet.AddRow().AddCell().SetString(enum)
			}
			dd := xlsx.NewXlsxCellDataValidation(true, true, true)
			if err := dd.SetInFileList(field.Name, 0, 0, 0, len(optionArr)-1); err != nil {
				blog.Errorf("SetDropList failed, err: %+v, rid: %s", err, handleFieldParam.Rid)
			}
			handleFieldParam.Sheet.Col(index).SetDataValidationWithStart(dd, common.HostAddMethodExcelIndexOffset)

		}
		handleFieldParam.Sheet.Col(index).SetType(xlsx.CellTypeString)

	case common.FieldTypeBool:
		dd := xlsx.NewXlsxCellDataValidation(true, true, true)
		if err := dd.SetDropList([]string{fieldTypeBoolTrue, fieldTypeBoolFalse}); err != nil {
			blog.Errorf("set drop list failed, err: %v, rid: %s", err, handleFieldParam.Rid)
		}
		handleFieldParam.Sheet.Col(index).SetDataValidationWithStart(dd, common.HostAddMethodExcelIndexOffset)
		handleFieldParam.Sheet.Col(index).SetType(xlsx.CellTypeString)

	default:
		handleFieldParam.Sheet.Col(index).SetType(xlsx.CellTypeString)
	}
}

// productExcelAssociationHeader TODO
// ProductExcelHeader Excel文件头部，
func productExcelAssociationHeader(ctx context.Context, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf,
	instNum int, asstList []*metadata.Association) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// 第一列(指标说明，橙色)
	cellStyle := getCellStyle(common.ExcelFirstColumnCellColor, common.ExcelHeaderFirstRowFontColor)
	// 第一列(其余格，粉色)
	colStyle := getCellStyle(common.ExcelHeaderFirstColumnColor, common.ExcelHeaderFirstRowFontColor)
	// 【2-5】列【二】排，(背景色，蓝色)
	backStyle := getCellStyle(common.ExcelHeaderOtherRowColor, common.ExcelHeaderFirstRowFontColor)

	sheet.Col(0).Width = 18
	sheet.Col(1).Width = 30
	firstColFields := []string{
		common.ExcelFirstColumnAssociationAttribute,
		common.ExcelFirstColumnFieldDescription,
		common.ExcelFirstColumnInstData,
	}
	for index, field := range firstColFields {
		cellName := sheet.Cell(index, 0)
		cellName.SetString(defLang.Language(field))
		cellName.SetStyle(cellStyle)
	}

	// 给第一列除前两行外的格子设置颜色(粉色)
	for i := 2; i < instNum+2; i++ {
		cellName := sheet.Cell(i, 0)
		cellName.SetStyle(colStyle)
	}
	sheet.Col(3).Width = 60
	sheet.Col(4).Width = 60

	cellAsstID := sheet.Cell(0, associationAsstObjIDIndex)
	cellAsstID.SetString(defLang.Language("excel_association_object_id"))
	cellAsstID.SetStyle(getHeaderFirstRowCellStyle(false))
	choiceCell := xlsx.NewXlsxCellDataValidation(true, true, true)
	// 确定关联标识的列表，定义excel选项下拉栏。此处需要查cc_ObjAsst表。
	pureAsstList := []string{}
	for _, asst := range asstList {
		pureAsstList = append(pureAsstList, asst.AssociationName)
	}
	pureAsstList = util.RemoveDuplicatesAndEmpty(pureAsstList)
	if err := choiceCell.SetDropList(pureAsstList); err != nil {
		blog.Errorf("SetDropList failed, err: %+v, rid: %s", err, rid)
	}
	sheet.Col(1).SetDataValidationWithStart(choiceCell, associationOPColIndex)

	cellOpID := sheet.Cell(0, associationOPColIndex)
	cellOpID.SetString(defLang.Language("excel_association_op"))
	cellOpID.SetStyle(getHeaderFirstRowCellStyle(false))
	dd := xlsx.NewXlsxCellDataValidation(true, true, true)
	if err := dd.SetDropList([]string{associationOPAdd, associationOPDelete}); err != nil {
		blog.Errorf("SetDropList failed, err: %+v, rid: %s", err, rid)
	}
	sheet.Col(2).SetDataValidationWithStart(dd, associationOPColIndex)

	cellSrcID := sheet.Cell(0, associationSrcInstIndex)
	cellSrcID.SetString(defLang.Language("excel_association_src_inst"))
	style := getHeaderFirstRowCellStyle(false)
	style.Alignment.WrapText = true
	cellSrcID.SetStyle(style)

	cellDstID := sheet.Cell(0, associationDstInstIndex)
	cellDstID.SetString(defLang.Language("excel_association_dst_inst"))
	style = getHeaderFirstRowCellStyle(false)
	style.Alignment.WrapText = true
	cellDstID.SetStyle(style)

	cell := sheet.Cell(1, associationAsstObjIDIndex)
	cell.SetString(defLang.Language("excel_example_association"))
	cell.SetStyle(backStyle)
	cell = sheet.Cell(1, associationOPColIndex)
	cell.SetString(defLang.Language("excel_example_op"))
	cell.SetStyle(backStyle)
	cell = sheet.Cell(1, associationSrcInstIndex)
	cell.SetString(defLang.Language("excel_example_association_src_inst"))
	cell.SetStyle(backStyle)
	cell = sheet.Cell(1, associationDstInstIndex)
	cell.SetString(defLang.Language("excel_example_association_dst_inst"))
	cell.SetStyle(backStyle)
}

const (
	associationOPColIndex     = 2
	associationAsstObjIDIndex = 1
	associationSrcInstIndex   = 3
	associationDstInstIndex   = 4

	associationOPAdd = "add"
	// associationOPUpdate = "update"
	associationOPDelete = "delete"
)

func getPrimaryKey(val interface{}) string {
	switch realVal := val.(type) {
	case []interface{}:
		if len(realVal) == 0 {
			return ""
		}
		valMap, ok := realVal[0].(map[string]interface{})
		if !ok {
			return ""
		}
		if valMap == nil {
			return ""
		}
		iVal := valMap[common.BKInstIDField]
		if iVal == nil {
			return ""
		}
		return fmt.Sprintf("%v", iVal)
	default:
		if realVal == nil {
			return ""
		}
		return fmt.Sprintf("%v", val)

	}
}
