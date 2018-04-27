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
	"net/http"
	"strings"

	"reflect"

	lang "configcenter/src/common/language"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

var (
	headerRow int = common.HostAddMethodExcelIndexOffset
)

//BuildExcelTemplate  return httpcode, error
func BuildExcelTemplate(url, objID, filename string, header http.Header, defLang lang.DefaultCCLanguageIf) error {
	conds := common.KvMap{common.BKObjIDField: objID, common.BKOwnerIDField: common.BKDefaultOwnerID, "page": common.KvMap{"start": 0, "limit": common.BKNoLimit}}
	result, err := httpRequest(url, conds, header)
	if nil != err {
		return err
	}
	blog.Info("get %s fields  url:%s", objID, url)
	blog.Info("get %s fields return:%s", objID, result)
	js, _ := simplejson.NewJson([]byte(result))
	hostFields, _ := js.Map()
	fields, _ := hostFields["data"].([]interface{})

	var file *xlsx.File
	file = xlsx.NewFile()
	sheet, err := file.AddSheet("host")
	if err != nil {
		blog.Errorf("get %s fields error:", objID, err.Error())
		return err
	}
	ProductExcelHealer(fields, getFilterFields(objID), sheet, defLang)
	err = file.Save(filename)
	if nil != err {
		return err
	}

	return nil
}

//ProductExcelHealer Excel文件头部，
func ProductExcelHealer(fields []interface{}, filter []string, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) {

	rowName := sheet.AddRow()
	rowType := sheet.AddRow()
	rowEnName := sheet.AddRow()
	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})
		fieldName, okName := mapField[common.BKPropertyNameField].(string)

		if !okName {
			fieldName = defLang.Language("web_excel_header_field_error") //"[未发现字段名(错误)]"
		}

		fieldType, _ := mapField[common.BKPropertyTypeField].(string)
		fieldTypeName, skip := getPropertyTypeAliasName(fieldType)
		if true == skip {
			//不需要用户输入的类型continue
			continue
		}
		isRequire := ""
		require, _ := mapField["bk_is_required"].(bool)
		if require {
			isRequire = defLang.Language("web_excel_header_required") //"(必填)"
		}
		enName, _ := mapField[common.BKPropertyIDField].(string)
		if util.Contains(filter, enName) {
			continue
		}
		cellName := rowName.AddCell()
		cellName.Value = fieldName + isRequire

		cellType := rowType.AddCell()
		cellType.Value = fieldTypeName

		cellEnName := rowEnName.AddCell()
		cellEnName.Value = enName
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
	if 0 == len(errCells) {
		return cells, nil
	} else {
		//web_import_field_not_found
		return nil, errors.New(defLang.Languagef("web_import_field_not_found", strings.Join(errCells, ",")))
	}

}

//getObjFieldIDs 获取properyID和properyName对应的值
func getObjFieldIDs(objID, url string, header http.Header) (common.KvMap, error) {
	conds := common.KvMap{common.BKObjIDField: objID, common.BKOwnerIDField: common.BKDefaultOwnerID, "page": common.KvMap{"start": 0, "limit": common.BKNoLimit}}
	result, err := httpRequest(url, conds, header)
	if nil != err {
		return nil, err
	}
	blog.Info("get %s fields  url:%s", objID, url)
	blog.Info("get %s fields return:%s", objID, result)
	js, _ := simplejson.NewJson([]byte(result))
	hostFields, _ := js.Map()
	fields, _ := hostFields["data"].([]interface{})
	ret := common.KvMap{}

	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})

		fieldName, _ := mapField[common.BKPropertyNameField].(string)
		fieldID, _ := mapField[common.BKPropertyIDField].(string)
		ret[fieldID] = fieldName
	}

	return ret, nil
}
