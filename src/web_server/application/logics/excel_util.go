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

	"github.com/rentiansheng/xlsx"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/util"
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
	return []string{"create_time"}
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
func setExcelRowDataByIndex(rowMap map[string]interface{}, sheet *xlsx.Sheet, rowIndex int, fields map[string]Property) {
	for id, val := range rowMap {
		property, ok := fields[id]
		if false == ok {
			continue
		}
		cell := sheet.Cell(rowIndex, property.ExcelColIndex)
		//cell.NumFmt = "@"

		switch property.PropertyType {
		case common.FieldTypeMultiAsst:
			arrVal, ok := val.([]interface{})
			if true == ok {
				vals := getAssociatePrimaryKey(arrVal, property.AsstObjPrimaryProperty)
				cell.SetString(strings.Join(vals, "\n"))
				style := cell.GetStyle()
				style.Alignment.WrapText = true
			}

		case common.FieldTypeSingleAsst:
			arrVal, ok := val.([]interface{})
			if true == ok {
				vals := getAssociatePrimaryKey(arrVal, property.AsstObjPrimaryProperty)
				cell.SetString(strings.Join(vals, "\n"))
				style := cell.GetStyle()
				style.Alignment.WrapText = true
			}

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

}

func getDataFromByExcelRow(row *xlsx.Row, rowIndex int, fields map[string]Property, defFields common.KvMap, nameIndexMap map[int]string, defLang lang.DefaultCCLanguageIf) (host map[string]interface{}, errMsg []string) {
	host = make(map[string]interface{})
	//errMsg := make([]string, 0)
	for celIDnex, cell := range row.Cells {
		fieldName, ok := nameIndexMap[celIDnex]
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
			cellValue, err := cell.Int64()
			if nil != err {
				errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, (celIDnex+1))) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (celIDnex + 1))
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
				errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", errMsg, fieldName, (celIDnex+1))) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (celIDnex + 1))
				blog.Errorf("%d row %s column get content error:%s", rowIndex+1, fieldName, err.Error())
				continue
			}
			host[fieldName] = cellValue
		default:
			errMsg = append(errMsg, defLang.Languagef("web_excel_row_handle_error", fieldName, (celIDnex+1))) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (celIDnex + 1))
			blog.Error("unknown the type, %v,   %v", reflect.TypeOf(cell), cell.Type())
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
		if true == skip {
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

		asstPrimaryKey := ""
		if 0 < len(field.AsstObjPrimaryProperty) {
			var primaryKeys []string
			for _, f := range field.AsstObjPrimaryProperty {
				primaryKeys = append(primaryKeys, f.Name)
			}
			asstPrimaryKey = fmt.Sprintf("(%s)", strings.Join(primaryKeys, common.ExcelAsstPrimaryKeySplitChar))
		}

		cellType := sheet.Cell(1, index)
		cellType.Value = fieldTypeName + asstPrimaryKey
		cellType.SetStyle(styleCell)

		cellEnName := sheet.Cell(2, index)
		cellEnName.Value = field.ID
		cellEnName.SetStyle(styleCell)

		switch field.PropertyType {
		case common.FieldTypeInt:
			sheet.Col(index).SetType(xlsx.CellTypeNumeric)
		case common.FieldTypeEnum:
			option := field.Option
			optionArr, ok := option.([]interface{})

			if ok {
				enumVals := getEnumNames(optionArr)

				if len(enumVals) < common.ExcelDataValidationListLen {
					dd := xlsx.NewXlsxCellDataValidation(true, true, true)
					dd.SetDropList(enumVals)
					sheet.Col(index).SetDataValidationWithStart(dd, common.HostAddMethodExcelIndexOffset+1)

				}
			}
			sheet.Col(index).SetType(xlsx.CellTypeString)

		case common.FieldTypeBool:
			dd := xlsx.NewXlsxCellDataValidation(true, true, true)
			dd.SetDropList([]string{fieldTypeBoolTrue, fieldTypeBoolFalse})
			sheet.Col(index).SetDataValidationWithStart(dd, common.HostAddMethodExcelIndexOffset+1)
			sheet.Col(index).SetType(xlsx.CellTypeString)
		default:
			sheet.Col(index).SetType(xlsx.CellTypeString)
		}

	}

}
