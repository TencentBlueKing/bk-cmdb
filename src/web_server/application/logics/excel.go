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
	//indexID := getFieldsIDIndexMap(fields)

	rowIndex := common.HostAddMethodExcelIndexOffset

	for _, row := range data {
		rowMap, ok := row.(map[string]interface{})

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

// BuildHostExcelFromData product excel from data
func BuildHostExcelFromData(objID string, fields map[string]Property, filter []string, data []interface{}, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) error {
	extFieldsTopoID := "cc_ext_field_topo"

	extFields := map[string]string{
		extFieldsTopoID: defLang.Language("web_ext_field_topo"),
	}
	fields = addExtFields(fields, extFields)

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

		rowMap, ok := hostData[common.BKInnerObjIDHost].(map[string]interface{})
		if false == ok {
			msg := fmt.Sprintf("data format error:%v", row)
			blog.Errorf(msg)
			return errors.New(msg)
		}
		moduleMap, ok := hostData[common.BKInnerObjIDModule].([]interface{})
		if ok {
			topo := util.GetStrValsFromArrMapInterfaceByKey(moduleMap, "TopModuleName")
			rowMap[extFieldsTopoID] = strings.Join(topo, "\n")
		}

		setExcelRowDataByIndex(rowMap, sheet, rowIndex, fields)
		rowIndex++

	}

	return nil
}

//BuildExcelTemplate  return httpcode, error
func BuildExcelTemplate(url, objID, filename string, header http.Header, defLang lang.DefaultCCLanguageIf) error {
	filterFields := getFilterFields(objID)
	fields, err := GetObjFieldIDs(objID, url, filterFields, header)
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
	productExcelHealer(fields, filterFields, sheet, defLang)
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
func GetExcelData(sheet *xlsx.Sheet, fields map[string]Property, defFields common.KvMap, isCheckHeader bool, firstRow int, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, error) {

	var err error
	nameIndexMap, err := checkExcelHealer(sheet, fields, isCheckHeader, defLang)
	if nil != err {

		return nil, err
	}
	hosts := make(map[int]map[string]interface{})
	index := headerRow
	if 0 != firstRow {
		index = firstRow
	}
	rowCnt := len(sheet.Rows)
	for ; index < rowCnt; index++ {
		row := sheet.Rows[index]
		host, getErr := getDataFromByExcelRow(row, index, fields, defFields, nameIndexMap, defLang)
		if nil != getErr {
			getErr = fmt.Errorf("%s;%s", getErr.Error())
			continue
		}
		if 0 == len(host) {
			hosts[index+1] = nil
		} else {
			hosts[index+1] = host
		}
	}
	if nil != err {

		return nil, err
	}

	return hosts, nil

}

//GetFilterFields 不需要展示字段
func GetFilterFields(objID string) []string {
	return getFilterFields(objID)
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
		}
		ret[index] = strName
	}
	if 0 != len(errCells) {
		//web_import_field_not_found
		return ret, errors.New(defLang.Languagef("web_import_field_not_found", strings.Join(errCells, ",")))
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

func getDataFromByExcelRow(row *xlsx.Row, rowIndex int, fields map[string]Property, defFields common.KvMap, nameIndexMap map[int]string, defLang lang.DefaultCCLanguageIf) (map[string]interface{}, error) {
	host := make(map[string]interface{})
	errMsg := ""
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

			cellValue, err := cell.Int()
			if nil != err {
				errMsg = defLang.Languagef("web_excel_row_handle_error", errMsg, fieldName, (celIDnex + 1)) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (celIDnex + 1))
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
				errMsg = defLang.Languagef("web_excel_row_handle_error", errMsg, fieldName, (celIDnex + 1)) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (celIDnex + 1))
				blog.Errorf("%d row %s column get content error:%s", rowIndex+1, fieldName, err.Error())
				continue
			}
			host[fieldName] = cellValue
		default:
			errMsg = defLang.Languagef("web_excel_row_handle_error", errMsg, fieldName, (celIDnex + 1)) //fmt.Sprintf("%s第%d行%d列无法处理内容;", errMsg, (index + 1), (celIDnex + 1))
			blog.Error("unknown the type, %v,   %v", reflect.TypeOf(cell), cell.Type())
			continue
		}
		field, ok := fields[fieldName]

		if true == ok {
			switch field.PropertyType {
			case common.FieldTypeBool:

				switch cell.Value {
				case fieldTypeBoolFalse:
					host[fieldName] = false
				case fieldTypeBoolTrue:
					host[fieldName] = true
				}

			case common.FieldTypeEnum:
				option, optionOk := field.Option.([]interface{})

				if optionOk {
					host[fieldName] = getEnumIDByName(cell.Value, option)
				}

			}
		}
	}
	if "" != errMsg {
		return nil, errors.New(errMsg)
	}
	if 0 == len(host) {
		return host, nil
	}
	for k, v := range defFields {
		host[k] = v
	}

	return host, nil

}
