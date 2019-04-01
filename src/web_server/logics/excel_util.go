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
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"

	"github.com/rentiansheng/xlsx"
)

var (
	headerRow = common.HostAddMethodExcelIndexOffset
)

//getFilterFields 不需要展示字段
func getFilterFields(objID string) []string {
	switch objID {
	case common.BKInnerObjIDHost:
		return []string{"create_time", "import_from", "bk_cloud_id", "bk_agent_status", "bk_agent_version"}
	default:
		return []string{"create_time"}
	}
	//return []string{"create_time"}
}

// checkExcelHealer check whether invalid fields exists in header and return headers
func checkExcelHealer(sheet *xlsx.Sheet, fields map[string]Property, isCheckHeader bool, defLang lang.DefaultCCLanguageIf) (map[int]string, error) {

	//rowLen := len(sheet.Rows[headerRow-1].Cells)
	var errCells []string
	ret := make(map[int]string)
	if headerRow > len(sheet.Rows) {
		return ret, errors.New(defLang.Language("web_excel_not_data"))
	}
	for index, name := range sheet.Rows[headerRow-1].Cells {
		strName := name.Value
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
		//web_import_field_not_found
		blog.Errorf(defLang.Languagef("web_import_field_not_found", strings.Join(errCells, ",")))
		return ret, errors.New(defLang.Languagef("web_import_field_not_found", errCells[0]+"..."))
	}
	return ret, nil

}

// setExcelRowDataByIndex insert  map[string]interface{}  to excel row by index,
// mapHeaderIndex:Correspondence between head and field
// fields each field description,  field type, isrequire, validate role
func setExcelRowDataByIndex(rowMap mapstr.MapStr, sheet *xlsx.Sheet, rowIndex int, fields map[string]Property) []PropertyPrimaryVal {

	primaryKeyArr := make([]PropertyPrimaryVal, 0)

	for id, val := range rowMap {
		property, ok := fields[id]
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
			if true == ok || true == enumIDOk {
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

func getDataFromByExcelRow(row *xlsx.Row, rowIndex int, fields map[string]Property, defFields common.KvMap, nameIndexMap map[int]string, defLang lang.DefaultCCLanguageIf) (host map[string]interface{}, errMsg []string) {
	host = make(map[string]interface{})
	//errMsg := make([]string, 0)
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
			host[fieldName] = cell.String()
		case xlsx.CellTypeStringFormula:
			host[fieldName] = cell.String()
		case xlsx.CellTypeNumeric:
			cellValue, err := cell.Float()
			if nil != err {
				errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, (cellIndex+1))) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (cellIndex + 1))
				blog.Errorf("%d row %s column get content error:%s", rowIndex+1, fieldName, err.Error())
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
				blog.Errorf("%d row %s column get content error:%s", rowIndex+1, fieldName, err.Error())
				continue
			}
			host[fieldName] = cellValue
		default:
			errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, (cellIndex+1))) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (cellIndex + 1))
			blog.Errorf("unknown the type, %v,   %v", reflect.TypeOf(cell), cell.Type())
			continue
		}

		field, ok := fields[fieldName]
		if !ok {
			blog.Errorf("%d row %s field not found ", rowIndex+1, fieldName)
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
			//convertor int not err , set field value to correct type
			if nil == err {
				host[fieldName] = intVal
			} else {
				blog.Debug("get excel cell value error, field:%s, value:%s, error:%s", fieldName, host[fieldName], err.Error())
			}
		case common.FieldTypeFloat:
			floatVal, err := util.GetFloat64ByInterface(host[fieldName])
			if nil == err {
				host[fieldName] = floatVal
			} else {
				blog.Debug("get excel cell value error, field:%s, value:%s, error:%s", fieldName, host[fieldName], err.Error())
			}
		default:
			if util.IsStrProperty(field.PropertyType) {
				host[fieldName] = cell.Value
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

// ProductExcelHealer Excel文件头部，
func productExcelHealer(fields map[string]Property, filter []string, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) {

	type excelHeader struct {
		first  string
		second string
		third  string
	}

	styleCell := getHeaderCellGeneralStyle()

	for _, field := range fields {
		index := field.ExcelColIndex
		sheet.Col(index).Width = 18
		fieldTypeName, skip := getPropertyTypeAliasName(field.PropertyType, defLang)
		if true == skip || field.NotExport {
			//不需要用户输入的类型continue
			continue
		}
		isRequire := ""

		if field.IsRequire {
			isRequire = defLang.Language("web_excel_header_required") //"(必填)"
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
				dd.SetDropList(enumVals)
				sheet.Col(index).SetDataValidationWithStart(dd, common.HostAddMethodExcelIndexOffset)

			}
			sheet.Col(index).SetType(xlsx.CellTypeString)

		case common.FieldTypeBool:
			dd := xlsx.NewXlsxCellDataValidation(true, true, true)
			dd.SetDropList([]string{fieldTypeBoolTrue, fieldTypeBoolFalse})
			sheet.Col(index).SetDataValidationWithStart(dd, common.HostAddMethodExcelIndexOffset)
			sheet.Col(index).SetType(xlsx.CellTypeString)
		default:
			sheet.Col(index).SetType(xlsx.CellTypeString)
		}

	}

}

// ProductExcelHealer Excel文件头部，
func productExcelAssociationHealer(sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) {

	cellAsstID := sheet.Cell(0, assciationAsstObjIDIndex)
	cellAsstID.SetString(defLang.Language("excel_association_object_id"))
	cellAsstID.SetStyle(getHeaderFirstRowCellStyle(false))

	cellOpID := sheet.Cell(0, associationOPColIndex)
	cellOpID.SetString(defLang.Language("excel_association_op"))
	cellOpID.SetStyle(getHeaderFirstRowCellStyle(false))
	dd := xlsx.NewXlsxCellDataValidation(true, true, true)
	dd.SetDropList([]string{associationOPAdd, associationOPDelete})
	sheet.Col(associationOPColIndex).SetDataValidationWithStart(dd, 1)

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
	sheet.Col(2).Width = 60
	sheet.Col(3).Width = 60
}

const (
	associationOPColIndex    = 1
	assciationAsstObjIDIndex = 0
	assciationSrcInstIndex   = 2
	assciationDstInstIndex   = 3

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
