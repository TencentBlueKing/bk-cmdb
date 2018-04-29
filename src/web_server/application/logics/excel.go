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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	lang "configcenter/src/common/language"
	//simplejson "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

var (
	headerRow = common.HostAddMethodExcelIndexOffset
)

// BuildExcelFromData product excel from data
func BuildExcelFromData(objID string, fields map[string]Property, filter []string, data []interface{}, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) error {
	if 0 == len(filter) {
		filter = getFilterFields(objID)
	} else {
		filter = append(filter, getFilterFields(objID)...)
	}

	productExcelHealer(fields, filter, sheet, defLang)
	indexID := getFieldsIDIndexMap(fields)

	var xlsRow *xlsx.Row
	rowIndex := common.HostAddMethodExcelIndexOffset

	for _, row := range data {
		hostData, ok := row.(map[string]interface{})
		if false == ok {
			msg := fmt.Sprintf("data format error:%v", row)
			blog.Errorf(msg)
			return errors.New(msg)
		}

		rowMap, ok := hostData["host"].(map[string]interface{})
		if false == ok {
			msg := fmt.Sprintf("data format error:%v", row)
			blog.Errorf(msg)
			return errors.New(msg)
		}
		isEmpty := true
		for id, val := range rowMap {
			// row unequal nil, pre row not used
			if nil == xlsRow {
				//row = sheet.AddRow()
			}
			index, ok := indexID[id]
			if false == ok {
				continue
			}
			isEmpty = true
			sheet.Cell(rowIndex, index).SetValue(val)

		}
		if false == isEmpty {
			rowIndex += 1
		}
	}

	return nil
}

// BuildExcelFromData product excel from data
func BuildHostExcelFromData(objID string, fields map[string]Property, filter []string, data []interface{}, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) error {

	productExcelHealer(fields, filter, sheet, defLang)
	//indexID := getFieldsIDIndexMap(fields)
	rowIndex := common.HostAddMethodExcelIndexOffset
	for _, row := range data {
		hostData, ok := row.(map[string]interface{})
		if false == ok {
			msg := fmt.Sprintf("data format error:%v", row)
			blog.Errorf(msg)
			return errors.New(msg)
		}

		rowMap, ok := hostData["host"].(map[string]interface{})
		if false == ok {
			msg := fmt.Sprintf("data format error:%v", row)
			blog.Errorf(msg)
			return errors.New(msg)
		}

		setExcelRowDataByIndex(rowMap, sheet, rowIndex, fields)
		rowIndex++

	}

	return nil
}

//BuildExcelTemplate  return httpcode, error
func BuildExcelTemplate(url, objID, filename string, header http.Header, defLang lang.DefaultCCLanguageIf) error {
	fields, err := GetObjFieldIDs(objID, url, header)
	if err != nil {
		blog.Errorf("get %s fields error:%s", objID, err.Error())
		return err
	}

	var file *xlsx.File
	file = xlsx.NewFile()
	sheet, err := file.AddSheet("host")
	if err != nil {
		blog.Errorf("get %s fields error:", objID, err.Error())
		return err
	}
	productExcelHealer(fields, getFilterFields(objID), sheet, defLang)
	err = file.Save(filename)
	if nil != err {
		return err
	}

	return nil
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

		cellType := sheet.Cell(1, index)
		cellType.Value = fieldTypeName
		cellType.SetStyle(styleCell)

		cellEnName := sheet.Cell(2, index)
		cellEnName.Value = field.ID
		cellEnName.SetStyle(styleCell)

		switch field.PropertyType {
		case common.FiledTypeInt:
			sheet.Col(index).SetType(xlsx.CellTypeNumeric)
		case common.FiledTypeEnum:
			option := field.Option
			optionArr, ok := option.([]interface{})

			if ok {
				enumVals := getEnumNames(optionArr)
				dd := xlsx.NewXlsxCellDataValidation(true, true, true)
				dd.SetDropList(enumVals)
				sheet.Col(index).SetDataValidationWithStart(dd, common.HostAddMethodExcelIndexOffset+1)
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

func AddDownExcelHttpHeader(c *gin.Context, name string) {
	if strings.HasSuffix(name, ".xls") {
		c.Header("Content-Type", "application/vnd.ms-excel")
	} else {
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	}
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Disposition", "attachment; filename="+name) //文件名
	c.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

//GetExcelData excel数据，一个kv结构，key行数（excel中的行数），value内容
func GetExcelData(sheet *xlsx.Sheet, fields, defFields common.KvMap, isCheckHeader bool, firstRow int, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, error) {

	cols, err := checkExcelHealer(sheet, fields, isCheckHeader, defLang)
	if nil != err {
		return nil, err
	}
	var errMsg string
	hosts := make(map[int]map[string]interface{})
	index := headerRow
	if 0 != firstRow {
		index = firstRow
	}
	rowCnt := len(sheet.Rows)
	maxCellLen := len(cols) //每行处理最大字段个数
	for ; index < rowCnt; index++ {
		isEmpty := true
		host := make(map[string]interface{})
		row := sheet.Rows[index]
		for celIDnex, cell := range row.Cells {
			if celIDnex >= maxCellLen { //当前行数字段，比对象属性多时候，忽略
				break
			}

			switch cell.Type() {
			case xlsx.CellTypeString:
				if "" == cell.String() {
					continue
				}
				isEmpty = false
				host[cols[celIDnex]] = cell.String()
			case xlsx.CellTypeStringFormula:
				if "" == cell.String() {
					continue
				}
				isEmpty = false
				host[cols[celIDnex]] = cell.String()
			case xlsx.CellTypeNumeric:

				cellValue, err := cell.Int()
				if nil != err {
					blog.Errorf("%d row %s column get content error:%s", index+1, cols[celIDnex], err.Error())
					continue
				}
				if 0 == cellValue {
					continue
				}
				isEmpty = false
				host[cols[celIDnex]] = cellValue
			case xlsx.CellTypeBool:
				cellValue := cell.Bool()
				isEmpty = false
				host[cols[celIDnex]] = cellValue
			case xlsx.CellTypeDate:
				cellValue, err := cell.GetTime(true)
				if nil != err {
					blog.Errorf("%d row %s column get content error:%s", index+1, cols[celIDnex], err.Error())
					continue
				}
				isEmpty = false
				host[cols[celIDnex]] = cellValue
			default:
				errMsg = defLang.Languagef("web_excel_row_handle_error", errMsg, (index + 1), (celIDnex + 1)) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (celIDnex + 1))
				blog.Error("unknown the type, %v,   %v", reflect.TypeOf(cell), cell.Type())
			}
		}
		for k, v := range defFields {
			host[k] = v
		}

		//内容不为空，加入返回数据中
		if false == isEmpty {
			hosts[index+1] = host
		} else {
			hosts[index+1] = nil
		}
	}
	if "" != errMsg {
		return nil, errors.New(errMsg)
	}

	return hosts, nil

}

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
func checkExcelHealer(sheet *xlsx.Sheet, fields common.KvMap, isCheckHeader bool, defLang lang.DefaultCCLanguageIf) ([]string, error) {

	//rowLen := len(sheet.Rows[headerRow-1].Cells)
	var cells []string
	var errCells []string
	if headerRow > len(sheet.Rows) {
		return nil, errors.New(defLang.Language("web_excel_not_data"))
	}
	nameCells := sheet.Rows[0].Cells
	for index, name := range sheet.Rows[headerRow-1].Cells {
		strName := name.Value
		cells = append(cells, strName)
		//是否坚持头部的字段存
		if !isCheckHeader {
			continue
		}
		_, ok := fields[strName]

		if !ok || "" != strName {

			cnName := nameCells[index].Value
			if ok {
				errCells = append(errCells, cnName)
			} else {
				errCells = append(errCells, strName)
			}

		}
	}
	if 0 != len(errCells) {
		//web_import_field_not_found
		return nil, errors.New(defLang.Languagef("web_import_field_not_found", strings.Join(errCells, ",")))
	}
	return cells, nil

}

// setExcelRowDataByIndex insert  map[string]interface{}  to excel row by index,
// mapHeaderIndex:Correspondence between head and field
// fields each field description,  field type, isrequire, validate role
func setExcelRowDataByIndex(rowMap map[string]interface{}, sheet *xlsx.Sheet, rowIndex int, fields map[string]Property) {
	for id, val := range rowMap {
		proptery, ok := fields[id]
		if false == ok {
			continue
		}
		cell := sheet.Cell(rowIndex, proptery.ExcelColIndex)
		//cell.NumFmt = "@"

		switch proptery.PropertyType {
		case common.FieldTypeMultiAsst:
			arrVal, ok := val.([]interface{})
			if true == ok {
				vals := getAssociateNames(arrVal)
				cell.SetString(strings.Join(vals, "\n"))
			}

		case common.FiledTypeSingleAsst:
			arrVal, ok := val.([]interface{})
			if true == ok {
				vals := getAssociateNames(arrVal)
				cell.SetString(strings.Join(vals, "\n"))
			}

		case common.FiledTypeEnum:
			var cellVal string
			arrVal, ok := proptery.Option.([]interface{})
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

		case common.FiledTypeInt:
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
