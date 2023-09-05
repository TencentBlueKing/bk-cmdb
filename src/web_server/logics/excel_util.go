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
	attrvalid "configcenter/src/common/valid/attribute"

	"github.com/tealeg/xlsx/v3"
)

const (
	userBracketsPattern         = `\([a-zA-Z0-9\@\p{Han} .,_-]*\)`
	organizationBracketsPattern = `\[(\d+)\]([^\s]+)`
)

var (
	headerRow          = common.AddExcelDataIndexOffset
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

func (lgc *Logics) getImportExcelPreData(objID string, header http.Header, f *xlsx.File,
	defLang lang.DefaultCCLanguageIf, bizID int64) (*ImportExcelPreData, error) {

	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	if len(f.Sheets) == 0 {
		return nil, defErr.Errorf(common.CCErrWebFileContentFail, defLang.Language("web_excel_content_empty"))
	}

	sheet := f.Sheets[0]
	if sheet == nil {
		return nil, defErr.Errorf(common.CCErrWebFileContentFail, defLang.Language("web_excel_sheet_not_found"))
	}

	// 获取模型字段信息
	fields, err := lgc.GetObjFieldIDs(objID, nil, nil, header, bizID, common.HostAddMethodExcelDefaultIndex)
	if err != nil {
		blog.Errorf("get object fields failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	// 检查excel表头，得到模型每一个字段从excel哪一列开始，并返回列号与字段id的映射map
	nameIndexMap, err := checkExcelHeader(sheet, fields, tableNameFieldIndex, defLang)
	if err != nil {
		blog.Errorf("check excel header failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	// 获取每一条数据，得到它从哪一行开始，哪一行结束
	dataRange, err := lgc.getCount(sheet, common.AddExcelDataIndexOffset, fields, defLang)
	if err != nil {
		blog.Errorf("get data range failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	// 获取表格字段以及表格字段里的属性在excel中的列的位置，字段的所占列的启始位置，以及将表格字段里的属性构造出property
	tableMap, err := lgc.getTableMap(fields, sheet)
	if err != nil {
		blog.Errorf("get table map failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	return &ImportExcelPreData{
		Fields:       fields,
		NameIndexMap: nameIndexMap,
		DataRange:    dataRange,
		TableMap:     tableMap,
		Sheet:        sheet,
	}, nil
}

// checkExcelHeader check whether invalid fields exists in header and return headers
func checkExcelHeader(sheet *xlsx.Sheet, fields map[string]Property, index int, defLang lang.DefaultCCLanguageIf) (
	map[int]string, error) {

	ret := make(map[int]string)
	if index > sheet.MaxRow {
		return ret, errors.New(defLang.Language("web_excel_not_data"))
	}

	row, err := sheet.Row(index - 1)
	if err != nil {
		return nil, err
	}
	for i := 0; i < row.GetCellCount(); i++ {
		cell := row.GetCell(i)
		strName := cell.String()
		// skip the ignored cell field
		if strName == common.ExcelCellIgnoreValue {
			continue
		}
		field, ok := fields[strName]
		if ok {
			field.ExcelColIndex = i
			fields[strName] = field
		}
		ret[i] = strName
	}
	return ret, nil

}

func (lgc *Logics) getTableMap(fields map[string]Property, sheet *xlsx.Sheet) (map[string]HeaderTable, error) {
	result := make(map[string]HeaderTable)
	for propertyID, field := range fields {
		if field.PropertyType != common.FieldTypeInnerTable {
			continue
		}

		tablePropertyIDIndex := make(map[string]int)
		for i := field.ExcelColIndex; i < field.ExcelColIndex+field.Length; i++ {
			cell, err := sheet.Cell(tableIDFieldIndex, i)
			if err != nil {
				return nil, err
			}
			tablePropertyIDIndex[cell.String()] = i
		}

		option, err := metadata.ParseTableAttrOption(field.Option)
		if err != nil {
			return nil, err
		}

		tableField := make(map[string]Property)
		nameIndexMap := make(map[int]string)
		for _, attr := range option.Header {
			tableField[attr.PropertyID] = Property{
				ID:            attr.PropertyID,
				Name:          attr.PropertyName,
				PropertyType:  attr.PropertyType,
				IsRequire:     attr.IsRequired,
				IsPre:         attr.IsPre,
				Option:        attr.Option,
				Group:         attr.PropertyGroup,
				ExcelColIndex: tablePropertyIDIndex[attr.PropertyID],
				Length:        1,
			}
			nameIndexMap[tablePropertyIDIndex[attr.PropertyID]] = attr.PropertyID
		}

		result[propertyID] = HeaderTable{
			Start:        field.ExcelColIndex,
			End:          field.ExcelColIndex + field.Length,
			Field:        tableField,
			NameIndexMap: nameIndexMap,
		}
	}

	return result, nil
}

func (lgc *Logics) getCount(sheet *xlsx.Sheet, index int, fields map[string]Property,
	defLang lang.DefaultCCLanguageIf) ([]ExcelDataRange, error) {

	indexMap := make(map[int]Property)
	for _, field := range fields {
		indexMap[field.ExcelColIndex] = field
	}

	var result []ExcelDataRange
	for index < sheet.MaxRow {
		data := ExcelDataRange{
			Start: index,
		}
		row, err := sheet.Row(index)
		if err != nil {
			return nil, err
		}

		for i := 0; i < row.GetCellCount(); i++ {
			property, ok := indexMap[i]
			if !ok {
				// 如果轮训的到表格字段的中间，那么这个时候会找不到，需要进行下一次查找
				continue
			}

			if property.PropertyType == common.FieldTypeInnerTable {
				continue
			}
			index += row.GetCell(i).VMerge
			break
		}
		index++
		data.End = index
		result = append(result, data)
	}

	if len(result) > common.ExcelImportMaxRow {
		return nil, errors.New(defLang.Languagef("web_excel_import_too_much", common.ExcelImportMaxRow))
	}

	return result, nil
}

// setExcelRowDataByIndex insert  map[string]interface{}  to excel row by index,
// mapHeaderIndex:Correspondence between head and field
// fields each field description,  field type, isrequire, validate role
func setExcelRowDataByIndex(rowMap mapstr.MapStr, sheet *xlsx.Sheet, rowIndex int, fields map[string]Property,
	rowCount int) error {

	style := getCellStyle(common.ExcelHeaderFirstColumnColor, common.ExcelHeaderFirstRowFontColor)
	for id, property := range fields {
		if property.NotExport {
			continue
		}

		if property.NotEditable {
			for i := rowIndex; i < rowIndex+rowCount; i++ {
				cell, err := sheet.Cell(i, property.ExcelColIndex)
				if err != nil {
					return err
				}
				cell.SetStyle(style)
			}
		}

		cell, err := sheet.Cell(rowIndex, property.ExcelColIndex)
		if err != nil {
			return err
		}
		if property.PropertyType != common.FieldTypeInnerTable && rowCount > 1 {
			cell.Merge(0, rowCount-1)
		}

		val, ok := rowMap[id]
		if !ok {
			continue
		}
		switch property.PropertyType {
		case common.FieldTypeEnum:
			setEnumCellValue(cell, val, property.Option)
		case common.FieldTypeBool:
			setBoolCellValue(cell, val)
		case common.FieldTypeInt:
			setIntCellValue(cell, val)
		case common.FieldTypeFloat:
			setFloatCellValue(cell, val)
		case common.FieldTypeInnerTable:
			err := setTableCellValue(val, sheet, rowIndex, property)
			if err != nil {
				return err
			}
		default:
			setDefaultCellValue(cell, val)
		}
	}
	return nil
}

func setEnumCellValue(cell *xlsx.Cell, val interface{}, option interface{}) {
	arrVal, ok := option.([]interface{})
	if !ok {
		blog.Errorf("option type is invalid, option: %v", option)
		return
	}
	strEnumID, enumIDOk := val.(string)
	if !enumIDOk {
		blog.Errorf("val type is invalid, val: %v", val)
		return
	}

	cellVal := getEnumNameByID(strEnumID, arrVal)
	cell.SetString(cellVal)
}

func setBoolCellValue(cell *xlsx.Cell, val interface{}) {
	bl, ok := val.(bool)
	if !ok {
		blog.Errorf("value type is not boolean, val: %s", val)
		return
	}
	if bl {
		cell.SetValue(fieldTypeBoolTrue)
		return
	}

	cell.SetValue(fieldTypeBoolFalse)
}

func setIntCellValue(cell *xlsx.Cell, val interface{}) {
	intVal, err := util.GetInt64ByInterface(val)
	if err != nil {
		blog.Errorf("val type is not int64, val: %v", val)
		return
	}
	cell.SetInt64(intVal)
}

func setFloatCellValue(cell *xlsx.Cell, val interface{}) {
	floatVal, err := util.GetFloat64ByInterface(val)
	if err != nil {
		blog.Errorf("val type is not float64, val: %v", val)
		return
	}
	cell.SetFloat(floatVal)
}

func setTableCellValue(val interface{}, sheet *xlsx.Sheet, rowIndex int, property Property) error {
	table, ok := val.([]mapstr.MapStr)
	if !ok {
		return fmt.Errorf("transfer table struct failed, val: %v", val)
	}
	option, err := metadata.ParseTableAttrOption(property.Option)
	if err != nil {
		return err
	}
	for tableIdx, data := range table {
		for idx, attr := range option.Header {
			val, exist := data[attr.PropertyID]
			cell, err := sheet.Cell(rowIndex+tableIdx, property.ExcelColIndex+idx)
			if err != nil {
				return err
			}
			switch attr.PropertyType {
			case common.FieldTypeBool:
				setBoolCellValue(cell, val)
			case common.FieldTypeInt:
				setIntCellValue(cell, val)
			case common.FieldTypeFloat:
				setFloatCellValue(cell, val)
			case common.FieldTypeEnumMulti:
				if !exist || val == nil {
					continue
				}
				items, ok := attr.Option.([]interface{})
				if !ok {
					return fmt.Errorf("enum multiple option param is invalid, option: %v", property.Option)
				}
				enumArr, ok := val.([]interface{})
				if !ok {
					return fmt.Errorf("convert enum multiple type value failed, val: %v", val)
				}
				enumMultiName := make([]string, 0)
				for _, enumID := range enumArr {
					id, ok := enumID.(string)
					if !ok {
						return fmt.Errorf("convert enum multiple id [%v] to string failed", enumID)
					}
					name := getEnumNameByID(id, items)
					enumMultiName = append(enumMultiName, name)
				}
				val = strings.Join(enumMultiName, "\n")
				setDefaultCellValue(cell, val)
			default:
				setDefaultCellValue(cell, val)
			}
		}
	}
	return nil
}

func setDefaultCellValue(cell *xlsx.Cell, val interface{}) {
	switch val.(type) {
	case string:
		strVal := val.(string)
		if strVal != "" {
			cell.SetString(val.(string))
		}
	default:
		cell.SetValue(val)
	}
}

func getDataFromExcel(ctx context.Context, preData *ImportExcelPreData, start, end int, defFields common.KvMap,
	defLang lang.DefaultCCLanguageIf) (map[string]interface{}, []string, error) {

	row, err := preData.Sheet.Row(start)
	if err != nil {
		return nil, nil, err
	}

	data, getErr := getDataFromByExcelRow(ctx, row, start, preData.Fields, defFields, preData.NameIndexMap, 1,
		row.GetCellCount(), defLang)
	var errMsg []string
	data, errMsg, err = buildDataWithTable(ctx, data, preData.Sheet, start, end, preData.TableMap, defLang)
	if err != nil {
		return nil, nil, err
	}
	getErr = append(getErr, errMsg...)
	return data, getErr, nil
}

func buildDataWithTable(ctx context.Context, data map[string]interface{}, sheet *xlsx.Sheet, start int, end int,
	tableMap map[string]HeaderTable, defLang lang.DefaultCCLanguageIf) (map[string]interface{}, []string, error) {

	var errMeg []string
	tableData := make(map[string][]map[string]interface{})
	for start < end {
		row, err := sheet.Row(start)
		if err != nil {
			return nil, nil, err
		}

		for propertyID, table := range tableMap {
			tableVal, getErr := getDataFromByExcelRow(ctx, row, start, table.Field, nil, table.NameIndexMap,
				table.Start, table.End, defLang)
			if len(tableVal) == 0 {
				continue
			}
			tableData[propertyID] = append(tableData[propertyID], tableVal)
			errMeg = append(errMeg, getErr...)
		}

		start++
	}

	for propertyID, table := range tableData {
		if len(table) == 0 {
			continue
		}
		data[propertyID] = table
	}

	return data, errMeg, nil
}

func getDataFromByExcelRow(ctx context.Context, row *xlsx.Row, rowIndex int, fields map[string]Property,
	defFields common.KvMap, nameIndexMap map[int]string, start, end int, defLang lang.DefaultCCLanguageIf) (
	map[string]interface{}, []string) {

	rid := util.ExtractRequestIDFromContext(ctx)
	result := make(map[string]interface{})
	errMsg := make([]string, 0)

	for i := start; i < end; i++ {
		cell := row.GetCell(i)
		fieldName, ok := nameIndexMap[i]
		if !ok || strings.Trim(fieldName, "") == "" || cell.String() == "" {
			continue
		}

		var hasField bool
		var field Property
		// 获取实例数据时，不为nil
		if fields != nil {
			field, hasField = fields[fieldName]
			// 如果这个字段是表格类型，那么不在这个函数中去处理获取
			if hasField && field.PropertyType == common.FieldTypeInnerTable {
				continue
			}
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

		if !hasField {
			continue
		}
		result, errMsg = buildAttrByPropertyType(rid, fieldName, cell.String(), rowIndex, field, result, defLang,
			errMsg)
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
	result map[string]interface{}, defLang lang.DefaultCCLanguageIf, errMsg []string) (map[string]interface{},
	[]string) {

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
	case common.FieldTypeEnumMulti:
		option, optionOK := field.Option.([]interface{})
		if !optionOK {
			break
		}
		cellValueList := strings.Split(cellValue, "\n")
		cellIDList := make([]string, 0)
		for _, cellName := range cellValueList {
			cellID := getEnumIDByName(cellName, option)
			cellIDList = append(cellIDList, cellID)
		}
		result[fieldName] = cellIDList
	case common.FieldTypeEnumQuote:
		result[fieldName] = strings.Split(cellValue, "\n")
	case common.FieldTypeInt:
		// convertor int not err, set field value to correct type
		if intVal, err := util.GetInt64ByInterface(result[fieldName]); err != nil {
			blog.Errorf("get excel cell value error, field: %s, value: %s, err: %v, rid: %s", fieldName,
				result[fieldName], err, rid)
		} else {
			result[fieldName] = intVal
		}
	case common.FieldTypeFloat:
		if floatVal, err := util.GetFloat64ByInterface(result[fieldName]); err == nil {
			result[fieldName] = floatVal
		} else {
			blog.Errorf("get excel cell value failed, field: %s, value: %s, err: %v, rid: %s", fieldName,
				result[fieldName], err, rid)
		}
	case common.FieldTypeOrganization:
		errMsg = parseOrganizationID(rid, fieldName, rowIndex, result, defLang, errMsg)
		if len(errMsg) != 0 {
			return nil, errMsg
		}
	case common.FieldTypeUser:
		userNames := util.GetStrByInterface(result[fieldName])
		userNames = userBracketsRegexp.ReplaceAllString(userNames, "")
		userNames = strings.Trim(strings.Trim(userNames, " "), ",")
		result[fieldName] = userNames
	default:
		if attrvalid.IsStrProperty(field.PropertyType) {
			result[fieldName] = strings.TrimSpace(cellValue)
		}
	}

	return result, errMsg
}

// parseOrganizationID parse organization id from excel
func parseOrganizationID(rid, fieldName string, rowIndex int, result map[string]interface{},
	defLang lang.DefaultCCLanguageIf, errMsg []string) []string {
	orgStr := util.GetStrByInterface(result[fieldName])
	if len(orgStr) <= 0 {
		errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, rowIndex+1)+
			defLang.Languagef("organization_type_invalid"))
		return errMsg
	}

	orgItems := strings.Split(orgStr, ",")
	org := make([]int64, len(orgItems))
	for i, v := range orgItems {
		var err error
		orgID := orgBracketsRegexp.FindStringSubmatch(v)
		if len(orgID) != 3 {
			blog.Errorf("regular matching is empty, please enter the correct content, field: %s, value: %s, "+
				"rid: %s", fieldName, result[fieldName], rid)
			errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, rowIndex+1)+
				defLang.Languagef("organization_type_invalid"))
			return errMsg
		}

		if org[i], err = strconv.ParseInt(orgID[1], 10, 64); err != nil {
			blog.Debug("get excel cell value error, field: %s, value: %s, err: %v, rid: %s", fieldName,
				result[fieldName], "not a valid organization type", rid)
			errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, rowIndex+1)+
				defLang.Languagef("organization_type_invalid"))
			return errMsg
		}
	}
	result[fieldName] = org
	return nil
}

// productExcelHeader Excel文件头部，
func productExcelHeader(ctx context.Context, fields map[string]Property, filter []string, xlsxFile *xlsx.File,
	sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	styleCell := getHeaderCellGeneralStyle()
	// 橙棕色
	cellStyle := getCellStyle(common.ExcelFirstColumnCellColor, common.ExcelHeaderFirstRowFontColor)
	// 粉色
	colStyle := getCellStyle(common.ExcelHeaderFirstColumnColor, common.ExcelHeaderFirstRowFontColor)
	width := float64(24)
	sheet.SetColWidth(1, 1, width)

	firstColFields := []string{common.ExcelFirstColumnFieldName, common.ExcelFirstColumnFieldType,
		common.ExcelFirstColumnFieldID, common.ExcelFirstColumnTableFieldName, common.ExcelFirstColumnTableFieldType,
		common.ExcelFirstColumnTableFieldID, common.ExcelFirstColumnInstData}

	for index, field := range firstColFields {
		cellName, err := sheet.Cell(index, 0)
		if err != nil {
			return err
		}
		fieldName := defLang.Language(field)
		cellName.Value = fieldName
		cellName.SetStyle(cellStyle)
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
		if err := handleField(field, handleFieldParam); err != nil {
			return err
		}
	}

	return nil
}

// productHostExcelHeader Excel文件头部，
func productHostExcelHeader(ctx context.Context, fields map[string]Property, filter []string, xlsxFile *xlsx.File,
	sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf, objName, cloudAreaName []string) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	styleCell := getHeaderCellGeneralStyle()
	// 橙棕色
	cellStyle := getCellStyle(common.ExcelFirstColumnCellColor, common.ExcelHeaderFirstRowFontColor)
	// 粉色
	colStyle := getCellStyle(common.ExcelHeaderFirstColumnColor, common.ExcelHeaderFirstRowFontColor)
	width := float64(24)
	sheet.SetColWidth(1, 1, width)
	// 字典中的值为国际化之后的"业务拓扑"和"业务名"，"集群"，”模块“，用来做判断，命中即变化相应的cell颜色。
	bizTopoMap := map[string]struct{}{
		defLang.Language("web_ext_field_topo"):        {},
		defLang.Language("biz_property_bk_biz_name"):  {},
		defLang.Language("web_ext_field_module_name"): {},
		defLang.Language("web_ext_field_set_name"):    {},
	}
	for _, name := range objName {
		bizTopoMap[name] = struct{}{}
	}

	firstColFields := []string{common.ExcelFirstColumnFieldName, common.ExcelFirstColumnFieldType,
		common.ExcelFirstColumnFieldID, common.ExcelFirstColumnTableFieldName, common.ExcelFirstColumnTableFieldType,
		common.ExcelFirstColumnTableFieldID, common.ExcelFirstColumnInstData}

	for index, field := range firstColFields {
		cellName, err := sheet.Cell(index, 0)
		if err != nil {
			return err
		}
		fieldName := defLang.Language(field)
		cellName.Value = fieldName
		cellName.SetStyle(cellStyle)
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
		if err := handleHostField(field, handleFieldParam, cloudAreaName, bizTopoMap); err != nil {
			return err
		}
	}

	return nil
}

func handleHostField(field Property, handleFieldParam *HandleFieldParam, cloudAreaName []string,
	bizTopoMap map[string]struct{}) error {

	isRequire := ""
	if field.IsRequire {
		isRequire = handleFieldParam.DefLang.Language("web_excel_header_required")
	}

	width := float64(24)
	// 主机部分特殊逻辑 针对云区域与topo进行处理
	if field.ID == common.BKCloudIDField {
		return handHostCloudIDField(handleFieldParam, field, width, cloudAreaName, isRequire)
	}

	if _, ok := bizTopoMap[field.Name]; ok {
		handleFieldParam.Sheet.SetColWidth(field.ExcelColIndex+1, field.ExcelColIndex+1, width)
		cellInRow0, err := handleFieldParam.Sheet.Cell(0, field.ExcelColIndex)
		if err != nil {
			return err
		}
		cellInRow0.Value = field.Name + isRequire
		cellInRow0.SetStyle(handleFieldParam.CellStyle)
		err = setExcelCellIgnored(handleFieldParam.Sheet, handleFieldParam.CellStyle, 1, field.ExcelColIndex)
		if err != nil {
			return err
		}
		err = setExcelCellIgnored(handleFieldParam.Sheet, handleFieldParam.CellStyle, 2, field.ExcelColIndex)
		if err != nil {
			return err
		}

		err = setTableRowMerge(handleFieldParam.Sheet, handleFieldParam.CellStyle, field.ExcelColIndex)
		if err != nil {
			return err
		}
		return nil
	}

	// 处理其他的通用属性逻辑
	if err := handleField(field, handleFieldParam); err != nil {
		return err
	}

	return nil
}

func handHostCloudIDField(handleFieldParam *HandleFieldParam, field Property, width float64, cloudAreaName []string,
	isRequire string) error {

	cellInRow0, err := handleFieldParam.Sheet.Cell(0, field.ExcelColIndex)
	if err != nil {
		return err
	}
	cellInRow0.Value = field.Name + isRequire
	// 设置单元格颜色
	cellInRow0.SetStyle(getHeaderFirstRowCellStyle(field.IsRequire))

	cellInRow2, err := handleFieldParam.Sheet.Cell(2, field.ExcelColIndex)
	if err != nil {
		return err
	}
	// 设置属性的id
	cellInRow2.Value = field.ID
	handleFieldParam.Sheet.SetColWidth(field.ExcelColIndex+1, field.ExcelColIndex+1, width)
	cellInRow2.SetStyle(handleFieldParam.StyleCell)

	err = setExcelCellIgnored(handleFieldParam.Sheet, handleFieldParam.StyleCell, 1, field.ExcelColIndex)
	if err != nil {
		return err
	}

	err = setTableRowMerge(handleFieldParam.Sheet, handleFieldParam.StyleCell, field.ExcelColIndex)
	if err != nil {
		return err
	}

	// 设置云区域的下拉选项
	if len(cloudAreaName) != 0 {
		enumSheet, err := handleFieldParam.File.AddSheet(field.Name)
		if err != nil {
			blog.Errorf("add enum sheet failed, err: %v, rid: %s", err, handleFieldParam.Rid)
			return err
		}
		for _, enum := range cloudAreaName {
			enumSheet.AddRow().AddCell().SetString(enum)
		}
		dd := xlsx.NewDataValidation(common.AddExcelDataIndexOffset, field.ExcelColIndex,
			xlsx.Excel2006MaxRowIndex, field.ExcelColIndex, true)
		if err := dd.SetInFileList(field.Name, 0, 0, 0, len(cloudAreaName)-1); err != nil {
			blog.Errorf("SetDropList failed, err: %+v, rid: %s", err, handleFieldParam.Rid)
		}
		dd.ShowInputMessage = true
		dd.ShowErrorMessage = true
		handleFieldParam.Sheet.AddDataValidation(dd)
		handleFieldParam.Sheet.SetType(field.ExcelColIndex+1, field.ExcelColIndex+1, xlsx.CellTypeString)
	}

	return nil
}

const (
	nameFieldIndex        = 0
	typeFieldIndex        = 1
	idFieldIndex          = 2
	tableNameFieldIndex   = 3
	tableTypeFieldIndex   = 4
	tableIDFieldIndex     = 5
	tableColMergeRowCount = 2
)

func handleField(field Property, handleFieldParam *HandleFieldParam) error {
	index := field.ExcelColIndex
	width := float64(24)
	handleFieldParam.Sheet.SetColWidth(index+1, index+field.Length, width)
	fieldTypeName, skip := getPropertyTypeAliasName(field.PropertyType, handleFieldParam.DefLang)
	if skip || field.NotExport {
		hidden := true
		handleFieldParam.Sheet.Col(index).Hidden = &hidden
		return nil
	}
	isRequire := ""

	if field.IsRequire {
		isRequire = handleFieldParam.DefLang.Language("web_excel_header_required")
	}
	if util.Contains(handleFieldParam.Filter, field.ID) {
		return nil
	}

	cellName, err := handleFieldParam.Sheet.Cell(nameFieldIndex, index)
	if err != nil {
		return err
	}
	cellName.Value = field.Name + isRequire
	style := getHeaderFirstRowCellStyle(field.IsRequire)
	cellName.SetStyle(style)
	for i := nameFieldIndex; i <= idFieldIndex; i++ {
		if err := setColMerge(handleFieldParam.Sheet, style, i, index, field.Length); err != nil {
			return err
		}
	}

	cellType, err := handleFieldParam.Sheet.Cell(typeFieldIndex, index)
	if err != nil {
		return err
	}
	cellType.Value = fieldTypeName
	cellType.SetStyle(handleFieldParam.StyleCell)
	for i := nameFieldIndex; i <= idFieldIndex; i++ {
		if err := setColMerge(handleFieldParam.Sheet, handleFieldParam.StyleCell, i, index, field.Length); err != nil {
			return err
		}
	}

	cellEnName, err := handleFieldParam.Sheet.Cell(idFieldIndex, index)
	if err != nil {
		return err
	}
	cellEnName.Value = field.ID
	cellEnName.SetStyle(handleFieldParam.StyleCell)
	for i := nameFieldIndex; i <= idFieldIndex; i++ {
		if err := setColMerge(handleFieldParam.Sheet, handleFieldParam.StyleCell, i, index, field.Length); err != nil {
			return err
		}
	}

	if field.PropertyType != common.FieldTypeInnerTable {
		if err := setTableRowMerge(handleFieldParam.Sheet, handleFieldParam.StyleCell, index); err != nil {
			return err
		}
	}

	tableColor := getCellStyle(common.ExcelTableHeaderColor, common.ExcelHeaderFirstRowFontColor)
	switch field.PropertyType {
	case common.FieldTypeInt:
		handleFieldTypeInt(handleFieldParam, index)
	case common.FieldTypeFloat:
		handleFieldTypeFloat(handleFieldParam, index)
	case common.FieldTypeEnum:
		handleFieldTypeEnum(handleFieldParam, index, &field)
	case common.FieldTypeEnumMulti:
		optionArr, ok := field.Option.([]interface{})
		if ok {
			handleFieldTypeEnumMulti(handleFieldParam, optionArr, field.Name)
		}
	case common.FieldTypeBool:
		handleFieldTypeBool(handleFieldParam, index)
	case common.FieldTypeInnerTable:
		if err := handleFieldTypeTable(handleFieldParam, index, &field, tableColor); err != nil {
			return err
		}

	default:
		handleFieldTypeDefault(handleFieldParam, index)
	}

	return nil
}

func handleFieldTypeInt(handleFieldParam *HandleFieldParam, index int) {
	handleFieldParam.Sheet.SetType(index+1, index+1, xlsx.CellTypeNumeric)
}

func handleFieldTypeFloat(handleFieldParam *HandleFieldParam, index int) {
	handleFieldParam.Sheet.SetType(index+1, index+1, xlsx.CellTypeNumeric)
}

func handleFieldTypeEnum(handleFieldParam *HandleFieldParam, index int, field *Property) {
	optionArr, ok := field.Option.([]interface{})
	if ok {
		enumSheet, err := handleFieldParam.File.AddSheet(field.Name)
		if err != nil {
			blog.Errorf("add enum sheet failed, err: %s, rid: %s", err, handleFieldParam.Rid)
		}

		for _, enum := range getEnumNames(optionArr) {
			enumSheet.AddRow().AddCell().SetString(enum)
		}
		dd := xlsx.NewDataValidation(common.AddExcelDataIndexOffset, index, xlsx.Excel2006MaxRowIndex, index,
			true)
		if err := dd.SetInFileList(field.Name, 0, 0, 0, len(optionArr)-1); err != nil {
			blog.Errorf("SetDropList failed, err: %+v, rid: %s", err, handleFieldParam.Rid)
		}
		dd.ShowInputMessage = true
		dd.ShowErrorMessage = true
		handleFieldParam.Sheet.AddDataValidation(dd)
	}
	handleFieldParam.Sheet.SetType(index+1, index+1, xlsx.CellTypeString)
}

func handleFieldTypeEnumMulti(handleFieldParam *HandleFieldParam, optionArr []interface{}, name string) {
	enumSheet, err := handleFieldParam.File.AddSheet(name)
	if err != nil {
		blog.Errorf("add enum sheet failed, err: %v, rid: %s", err, handleFieldParam.Rid)
	}

	for _, enum := range getEnumNames(optionArr) {
		enumSheet.AddRow().AddCell().SetString(enum)
	}
}

func handleFieldTypeBool(handleFieldParam *HandleFieldParam, index int) {
	dd := xlsx.NewDataValidation(common.AddExcelDataIndexOffset, index, xlsx.Excel2006MaxRowIndex, index,
		true)
	if err := dd.SetDropList([]string{fieldTypeBoolTrue, fieldTypeBoolFalse}); err != nil {
		blog.Errorf("set drop list failed, err: %v, rid: %s", err, handleFieldParam.Rid)
	}
	dd.ShowInputMessage = true
	dd.ShowErrorMessage = true
	handleFieldParam.Sheet.AddDataValidation(dd)
	handleFieldParam.Sheet.SetType(index+1, index+1, xlsx.CellTypeString)
}

func handleFieldTypeDefault(handleFieldParam *HandleFieldParam, index int) {
	handleFieldParam.Sheet.SetType(index+1, index+1, xlsx.CellTypeString)
}

func handleFieldTypeTable(handleFieldParam *HandleFieldParam, index int, field *Property, color *xlsx.Style) error {
	option, err := metadata.ParseTableAttrOption(field.Option)
	if err != nil {
		return err
	}
	for idx, attr := range option.Header {
		tableNameField, err := handleFieldParam.Sheet.Cell(tableNameFieldIndex, index+idx)
		if err != nil {
			return err
		}
		tableNameField.Value = attr.PropertyName
		tableNameField.SetStyle(color)

		tableTypeField, err := handleFieldParam.Sheet.Cell(tableTypeFieldIndex, index+idx)
		if err != nil {
			return err
		}
		tableTypeField.Value = handleFieldParam.DefLang.Language("field_type_" + attr.PropertyType)
		tableTypeField.SetStyle(color)

		tableIDField, err := handleFieldParam.Sheet.Cell(tableIDFieldIndex, index+idx)
		if err != nil {
			return err
		}
		tableIDField.Value = attr.PropertyID
		tableIDField.SetStyle(color)

		switch attr.PropertyType {
		case common.FieldTypeInt:
			handleFieldTypeInt(handleFieldParam, index+idx)
		case common.FieldTypeFloat:
			handleFieldTypeFloat(handleFieldParam, index+idx)
		case common.FieldTypeEnumMulti:
			optionArr, ok := attr.Option.([]interface{})
			if ok {
				handleFieldTypeEnumMulti(handleFieldParam, optionArr, field.Name+"##"+attr.PropertyName)
			}
		case common.FieldTypeBool:
			handleFieldTypeBool(handleFieldParam, index+idx)
		default:
			handleFieldTypeDefault(handleFieldParam, index+idx)
		}
	}

	return nil
}

// productExcelAssociationHeader Excel文件头部
// NOCC:golint/fnsize(后续重构处理)
func productExcelAssociationHeader(ctx context.Context, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf,
	instNum int, asstList []*metadata.Association) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	// 第一列(指标说明，橙色)
	cellStyle := getCellStyle(common.ExcelFirstColumnCellColor, common.ExcelHeaderFirstRowFontColor)
	// 第一列(其余格，粉色)
	colStyle := getCellStyle(common.ExcelHeaderFirstColumnColor, common.ExcelHeaderFirstRowFontColor)
	// 【2-5】列【二】排，(背景色，蓝色)
	backStyle := getCellStyle(common.ExcelHeaderOtherRowColor, common.ExcelHeaderFirstRowFontColor)

	firstColWidth := float64(24)
	sheet.SetColWidth(1, 1, firstColWidth)
	secondColWidth := float64(36)
	sheet.SetColWidth(2, 2, secondColWidth)
	firstColFields := []string{
		common.ExcelFirstColumnAssociationAttribute,
		common.ExcelFirstColumnFieldDescription,
		common.ExcelFirstColumnInstData,
	}
	for index, field := range firstColFields {
		cellName, err := sheet.Cell(index, 0)
		if err != nil {
			return err
		}
		cellName.SetString(defLang.Language(field))
		cellName.SetStyle(cellStyle)
	}

	// 给第一列除前两行外的格子设置颜色(粉色)
	for i := 2; i < instNum+2; i++ {
		cellName, err := sheet.Cell(i, 0)
		if err != nil {
			return err
		}
		cellName.SetStyle(colStyle)
	}
	thirdWidth := float64(12)
	sheet.SetColWidth(3, 3, thirdWidth)
	fourthAndFifthWidth := float64(80)
	sheet.SetColWidth(4, 5, fourthAndFifthWidth)

	cellAsstID, err := sheet.Cell(0, associationAsstObjIDIndex)
	if err != nil {
		return err
	}
	cellAsstID.SetString(defLang.Language("excel_association_object_id"))
	cellAsstID.SetStyle(getHeaderFirstRowCellStyle(false))
	choiceCell := xlsx.NewDataValidation(associationOPColIndex, 1, xlsx.Excel2006MaxRowIndex, 1, true)
	// 确定关联标识的列表，定义excel选项下拉栏。此处需要查cc_ObjAsst表。
	pureAsstList := []string{}
	for _, asst := range asstList {
		pureAsstList = append(pureAsstList, asst.AssociationName)
	}
	pureAsstList = util.RemoveDuplicatesAndEmpty(pureAsstList)
	if err := choiceCell.SetDropList(pureAsstList); err != nil {
		blog.Errorf("SetDropList failed, err: %+v, rid: %s", err, rid)
	}
	choiceCell.ShowInputMessage = true
	choiceCell.ShowErrorMessage = true
	sheet.AddDataValidation(choiceCell)

	cellOpID, err := sheet.Cell(0, associationOPColIndex)
	if err != nil {
		return err
	}
	cellOpID.SetString(defLang.Language("excel_association_op"))
	cellOpID.SetStyle(getHeaderFirstRowCellStyle(false))
	dd := xlsx.NewDataValidation(associationOPColIndex, 2, xlsx.Excel2006MaxRowIndex, 2, true)
	if err := dd.SetDropList([]string{associationOPAdd, associationOPDelete}); err != nil {
		blog.Errorf("SetDropList failed, err: %+v, rid: %s", err, rid)
	}
	dd.ShowInputMessage = true
	dd.ShowErrorMessage = true
	sheet.AddDataValidation(dd)

	cellSrcID, err := sheet.Cell(0, associationSrcInstIndex)
	if err != nil {
		return err
	}
	cellSrcID.SetString(defLang.Language("excel_association_src_inst"))
	style := getHeaderFirstRowCellStyle(false)
	style.Alignment.WrapText = true
	cellSrcID.SetStyle(style)

	cellDstID, err := sheet.Cell(0, associationDstInstIndex)
	if err != nil {
		return err
	}
	cellDstID.SetString(defLang.Language("excel_association_dst_inst"))
	style = getHeaderFirstRowCellStyle(false)
	style.Alignment.WrapText = true
	cellDstID.SetStyle(style)

	cell, err := sheet.Cell(1, associationAsstObjIDIndex)
	if err != nil {
		return err
	}
	cell.SetString(defLang.Language("excel_example_association"))
	cell.SetStyle(backStyle)
	cell, err = sheet.Cell(1, associationOPColIndex)
	if err != nil {
		return err
	}
	cell.SetString(defLang.Language("excel_example_op"))
	cell.SetStyle(backStyle)
	cell, err = sheet.Cell(1, associationSrcInstIndex)
	if err != nil {
		return err
	}
	cell.SetString(defLang.Language("excel_example_association_src_inst"))
	cell.SetStyle(backStyle)
	cell, err = sheet.Cell(1, associationDstInstIndex)
	if err != nil {
		return err
	}
	cell.SetString(defLang.Language("excel_example_association_dst_inst"))
	cell.SetStyle(backStyle)

	return nil
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

func setTableRowMerge(sheet *xlsx.Sheet, style *xlsx.Style, index int) error {
	cellTableType, err := sheet.Cell(tableTypeFieldIndex, index)
	if err != nil {
		return err
	}
	cellTableType.SetStyle(style)

	cellTableID, err := sheet.Cell(tableIDFieldIndex, index)
	if err != nil {
		return err
	}
	cellTableID.SetStyle(style)

	cellTableName, err := sheet.Cell(tableNameFieldIndex, index)
	if err != nil {
		return err
	}
	cellTableName.SetStyle(style)
	cellTableName.Merge(0, tableColMergeRowCount)
	return nil
}

func setColMerge(sheet *xlsx.Sheet, style *xlsx.Style, rowIdx, colIdx, length int) error {
	if length == 1 {
		return nil
	}
	for i := 1; i < length; i++ {
		cell, err := sheet.Cell(rowIdx, colIdx+i)
		if err != nil {
			return err
		}
		cell.SetStyle(style)
	}

	cell, err := sheet.Cell(rowIdx, colIdx)
	if err != nil {
		return err
	}
	cell.Merge(length-1, 0)
	return nil
}
