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
	userAliasPattern = `\([a-zA-Z0-9\@\p{Han} .,_-]*\)`
)

var (
	headerRow       = common.HostAddMethodExcelIndexOffset
	userAliasRegexp = regexp.MustCompile(userAliasPattern)
)

// getFilterFields 不需要展示字段
func getFilterFields(objID string) []string {
	switch objID {
	case common.BKInnerObjIDHost:
		return []string{"create_time", "import_from", "bk_agent_status", "bk_agent_version", "bk_set_name", "bk_module_name", "bk_biz_name"}
	default:
		return []string{"create_time"}
	}
	//return []string{"create_time"}
}

func getCustomFields(filterFields []string, customFieldsStr string) []string {
	customFields := strings.Split(customFieldsStr, ",")
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
	// valid excel three row is instance property fields,
	// excel three row  values  exceeding 1/2 does not appear in the field array,
	// indicating that the third line of the excel template was deleted
	if len(errCells) > len(sheet.Rows[headerRow-1].Cells)/2 && true == isCheckHeader {
		blog.Errorf("err:%s, no found fields %s, rid:%s", defLang.Language("web_import_field_not_found"), strings.Join(errCells, ","), rid)
		return ret, errors.New(defLang.Language("web_import_field_not_found"))
	}
	return ret, nil

}

// setExcelRowDataByIndex insert  map[string]interface{}  to excel row by index,
// mapHeaderIndex:Correspondence between head and field
// fields each field description,  field type, isrequire, validate role
func setExcelRowDataByIndex(rowMap mapstr.MapStr, sheet *xlsx.Sheet, rowIndex int, fields map[string]Property) []PropertyPrimaryVal {

	primaryKeyArr := make([]PropertyPrimaryVal, 0)

	// 非模型字段导出是没有field中没有ID 字段，因为导入的时候，第二行是作为Property
	for id, property := range fields {
		val, ok := rowMap[id]
		if false == ok {
			continue
		}
		if property.NotExport {
			if property.IsOnly {
				primaryKeyArr = append(primaryKeyArr, PropertyPrimaryVal{
					ID:     property.ID,
					Name:   property.Name,
					StrVal: getPrimaryKey(val),
				})
			}
			continue
		}

		cell := sheet.Cell(rowIndex, property.ExcelColIndex)
		//cell.NumFmt = "@"

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

		if property.IsOnly {
			primaryKeyArr = append(primaryKeyArr, PropertyPrimaryVal{
				ID:     property.ID,
				Name:   property.Name,
				StrVal: cell.String(),
			})
		}

	}

	return primaryKeyArr

}

func getDataFromByExcelRow(ctx context.Context, row *xlsx.Row, rowIndex int, fields map[string]Property, defFields common.KvMap, nameIndexMap map[int]string, defLang lang.DefaultCCLanguageIf) (host map[string]interface{}, errMsg []string) {
	rid := util.ExtractRequestIDFromContext(ctx)
	host = make(map[string]interface{})
	// errMsg := make([]string, 0)
	for cellIndex, cell := range row.Cells {
		fieldName, ok := nameIndexMap[cellIndex]
		if false == ok {
			continue
		}
		if "" == strings.Trim(fieldName, "") {
			continue
		}
		if "" == cell.Value {
			continue
		}

		switch cell.Type() {
		case xlsx.CellTypeString:
			host[fieldName] = strings.TrimSpace(cell.String())
		case xlsx.CellTypeStringFormula:
			host[fieldName] = strings.TrimSpace(cell.String())
		case xlsx.CellTypeNumeric:
			cellValue, err := cell.Float()
			if nil != err {
				errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, (cellIndex+1))) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (cellIndex + 1))
				blog.Errorf("%d row %s column get content error:%s, rid: %s", rowIndex+1, fieldName, err.Error(), rid)
				continue
			}
			host[fieldName] = cellValue
		case xlsx.CellTypeBool:
			cellValue := cell.Bool()
			host[fieldName] = cellValue
		case xlsx.CellTypeDate:
			cellValue, err := cell.GetTime(true)
			if nil != err {
				errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", errMsg, fieldName, (cellIndex+1))) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (cellIndex + 1))
				blog.Errorf("%d row %s column get content error:%s, rid: %s", rowIndex+1, fieldName, err.Error(), rid)
				continue
			}
			host[fieldName] = cellValue
		default:
			errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, (cellIndex+1))) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (cellIndex + 1))
			blog.Errorf("unknown the type, %v,   %v, rid: %s", reflect.TypeOf(cell), cell.Type(), rid)
			continue
		}

		field, ok := fields[fieldName]
		if !ok {
			blog.Errorf("%d row %s field not found , rid: %s", rowIndex+1, fieldName, rid)
			continue
		}
		switch field.PropertyType {
		case common.FieldTypeBool:
			switch host[fieldName].(type) {
			case bool:
			default:
				bl, err := strconv.ParseBool(cell.Value)
				if nil == err {
					host[fieldName] = bl
				}
			}
		case common.FieldTypeEnum:
			option, optionOk := field.Option.([]interface{})

			if optionOk {
				host[fieldName] = getEnumIDByName(cell.Value, option)
			}
		case common.FieldTypeInt:
			intVal, err := util.GetInt64ByInterface(host[fieldName])
			// convertor int not err , set field value to correct type
			if nil == err {
				host[fieldName] = intVal
			} else {
				blog.Debug("get excel cell value error, field:%s, value:%s, error:%s, rid: %s", fieldName, host[fieldName], err.Error(), rid)
			}
		case common.FieldTypeFloat:
			floatVal, err := util.GetFloat64ByInterface(host[fieldName])
			if nil == err {
				host[fieldName] = floatVal
			} else {
				blog.Debug("get excel cell value error, field:%s, value:%s, error:%s, rid: %s", fieldName, host[fieldName], err.Error(), rid)
			}
		case common.FieldTypeOrganization:
			org := util.GetStrByInterface(host[fieldName])
			if len(org) >= 2 && strings.HasPrefix(org, "[") && strings.HasSuffix(org, "]") {
				if strings.TrimSpace(org[1:len(org)-1]) == "" {
					host[fieldName] = []int64{}
					break
				}
				orgItems := strings.Split(org[1:len(org)-1], ",")
				orgSlice := make([]int64, len(orgItems))
				var err error
				for i, v := range orgItems {
					orgSlice[i], err = strconv.ParseInt(strings.TrimSpace(v), 10, 64)
					if err != nil {
						blog.Debug("get excel cell value error, field:%s, value:%s, error:%s, rid: %s", fieldName, host[fieldName], "not a valid organization type", rid)
						break
					}
				}
				if err == nil {
					host[fieldName] = orgSlice
				}
			} else {
				blog.Debug("get excel cell value error, field:%s, value:%s, error:%s, rid: %s", fieldName, host[fieldName], "not a valid organization type", rid)
			}
		case common.FieldTypeUser:
			// convert userNames,  eg: " admin(admin),xiaoming(小明 ),leo(li hong),  " => "admin,xiaoming,leo"
			userNames := util.GetStrByInterface(host[fieldName])
			userNames = userAliasRegexp.ReplaceAllString(userNames, "")
			userNames = strings.Trim(strings.Trim(userNames, " "), ",")
			host[fieldName] = userNames
		default:
			if util.IsStrProperty(field.PropertyType) {
				host[fieldName] = strings.TrimSpace(cell.Value)
			}

		}

	}
	if 0 != len(errMsg) {
		return nil, errMsg
	}
	if 0 == len(host) {
		return host, nil
	}
	for k, v := range defFields {
		host[k] = v
	}

	return host, nil

}

// ProductExcelHeader Excel文件头部，
func productExcelHeader(ctx context.Context, fields map[string]Property, filter []string, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) {
	rid := util.ExtractRequestIDFromContext(ctx)
	styleCell := getHeaderCellGeneralStyle()
	//橙棕色
	cellStyle := getCellStyle(common.ExcelFirstColumnCellColor, common.ExcelHeaderFirstRowFontColor)
	//粉色
	colStyle := getCellStyle(common.ExcelHeaderFirstColumnColor, common.ExcelHeaderFirstRowFontColor)

	sheet.Col(0).Width = 18
	//字典中的值为国际化之后的"业务拓扑"和"业务名"，用来做判断，命中即变化相应的cell颜色。
	bizTopoMap := map[string]int{
		defLang.Language("web_ext_field_topo"):       1,
		defLang.Language("biz_property_bk_biz_name"): 1,
	}
	firstColFields := []string{common.ExcelFirstColumnFieldName, common.ExcelFirstColumnFieldType, common.ExcelFirstColumnFieldID, common.ExcelFirstColumnInstData}
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

	for _, field := range fields {
		index := field.ExcelColIndex
		sheet.Col(index).Width = 18
		fieldTypeName, skip := getPropertyTypeAliasName(field.PropertyType, defLang)
		if true == skip || field.NotExport {
			// 不需要用户输入的类型continue
			sheet.Col(index).Hidden = true
			continue
		}
		isRequire := ""

		if field.IsRequire {
			// "(必填)"
			isRequire = defLang.Language("web_excel_header_required")
		}
		if util.Contains(filter, field.ID) {
			continue
		}
		cellName := sheet.Cell(0, index)
		cellName.Value = field.Name + isRequire
		cellName.SetStyle(getHeaderFirstRowCellStyle(field.IsRequire))

		cellType := sheet.Cell(1, index)
		cellType.Value = fieldTypeName
		cellType.SetStyle(styleCell)

		cellEnName := sheet.Cell(2, index)
		cellEnName.Value = field.ID
		cellEnName.SetStyle(styleCell)

		switch field.PropertyType {
		case common.FieldTypeInt:
			sheet.Col(index).SetType(xlsx.CellTypeNumeric)
		case common.FieldTypeFloat:
			sheet.Col(index).SetType(xlsx.CellTypeNumeric)
		case common.FieldTypeEnum:
			option := field.Option
			optionArr, ok := option.([]interface{})

			if ok {
				enumVals := getEnumNames(optionArr)
				dd := xlsx.NewXlsxCellDataValidation(true, true, true)
				if err := dd.SetDropList(enumVals); err != nil {
					blog.Errorf("SetDropList failed, err: %+v, rid: %s", err, rid)
				}
				sheet.Col(index).SetDataValidationWithStart(dd, common.HostAddMethodExcelIndexOffset)

			}
			sheet.Col(index).SetType(xlsx.CellTypeString)

		case common.FieldTypeBool:
			dd := xlsx.NewXlsxCellDataValidation(true, true, true)
			if err := dd.SetDropList([]string{fieldTypeBoolTrue, fieldTypeBoolFalse}); err != nil {
				blog.Errorf("SetDropList failed, err: %+v, rid: %s", err, rid)
			}
			sheet.Col(index).SetDataValidationWithStart(dd, common.HostAddMethodExcelIndexOffset)
			sheet.Col(index).SetType(xlsx.CellTypeString)

		default:
			if _, ok := bizTopoMap[field.Name]; ok {
				cellName := sheet.Cell(0, index)
				cellName.Value = field.Name + isRequire
				cellName.SetStyle(cellStyle)
				setExcelCellIgnored(sheet, cellStyle, 1, index)
				setExcelCellIgnored(sheet, cellStyle, 2, index)

				// 给业务拓扑和业务列剩下的空格设置颜色
				for i := 3; i < 1000; i++ {
					cellName := sheet.Cell(i, index)
					cellName.SetStyle(colStyle)
				}
				sheet.Col(index).SetType(xlsx.CellTypeString)
			}

			if field.ID == common.BKCloudIDField {
				setExcelCellIgnored(sheet, styleCell, 1, index)
				setExcelCellIgnored(sheet, styleCell, 2, index)
			}

			sheet.Col(index).SetType(xlsx.CellTypeString)
		}
	}
}

// ProductExcelHeader Excel文件头部，
func productExcelAssociationHeader(ctx context.Context, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf, instNum int, asstList []*metadata.Association) {
	rid := util.ExtractRequestIDFromContext(ctx)

	//第一列(指标说明，橙色)
	cellStyle := getCellStyle(common.ExcelFirstColumnCellColor, common.ExcelHeaderFirstRowFontColor)
	//第一列(其余格，粉色)
	colStyle := getCellStyle(common.ExcelHeaderFirstColumnColor, common.ExcelHeaderFirstRowFontColor)
	//【2-5】列【二】排，(背景色，蓝色)
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

	cellAsstID := sheet.Cell(0, assciationAsstObjIDIndex)
	cellAsstID.SetString(defLang.Language("excel_association_object_id"))
	cellAsstID.SetStyle(getHeaderFirstRowCellStyle(false))
	choiceCell := xlsx.NewXlsxCellDataValidation(true, true, true)
	//确定关联标识的列表，定义excel选项下拉栏。此处需要查cc_ObjAsst表。
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

	cellSrcID := sheet.Cell(0, assciationSrcInstIndex)
	cellSrcID.SetString(defLang.Language("excel_association_src_inst"))
	style := getHeaderFirstRowCellStyle(false)
	style.Alignment.WrapText = true
	cellSrcID.SetStyle(style)

	cellDstID := sheet.Cell(0, assciationDstInstIndex)
	cellDstID.SetString(defLang.Language("excel_association_dst_inst"))
	style = getHeaderFirstRowCellStyle(false)
	style.Alignment.WrapText = true
	cellDstID.SetStyle(style)

	cell := sheet.Cell(1, assciationAsstObjIDIndex)
	cell.SetString(defLang.Language("excel_example_association"))
	cell.SetStyle(backStyle)
	cell = sheet.Cell(1, associationOPColIndex)
	cell.SetString(defLang.Language("excel_example_op"))
	cell.SetStyle(backStyle)
	cell = sheet.Cell(1, assciationSrcInstIndex)
	cell.SetString(defLang.Language("excel_example_association_src_inst"))
	cell.SetStyle(backStyle)
	cell = sheet.Cell(1, assciationDstInstIndex)
	cell.SetString(defLang.Language("excel_example_association_dst_inst"))
	cell.SetStyle(backStyle)
}

const (
	associationOPColIndex    = 2
	assciationAsstObjIDIndex = 1
	assciationSrcInstIndex   = 3
	assciationDstInstIndex   = 4

	associationOPAdd = "add"
	//associationOPUpdate = "update"
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
