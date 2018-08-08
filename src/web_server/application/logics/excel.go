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
	"net/http"
	"strings"

	//simplejson "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/util"
)

// BuildExcelFromData product excel from data
func BuildExcelFromData(objID string, fields map[string]Property, filter []string, data []interface{}, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) error {
	addSystemField(fields, common.BKINnerObjIDObject, defLang)
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
	addSystemField(fields, common.BKInnerObjIDHost, defLang)

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
	blog.V(5).Infof("BuildExcelTemplate fields count:%d", fields)
	productExcelHealer(fields, filterFields, sheet, defLang)
	ProductExcelCommentSheet(file, defLang)

	err = file.Save(filename)
	if nil != err {
		return err
	}

	return nil
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

//GetExcelData excel数据，一个kv结构，key行数（excel中的行数），value内容
func GetRawExcelData(sheet *xlsx.Sheet, defFields common.KvMap, firstRow int, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, error) {

	var err error
	nameIndexMap, err := checkExcelHealer(sheet, nil, false, defLang)
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
		host, getErr := getDataFromByExcelRow(row, index, nil, defFields, nameIndexMap, defLang)
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
