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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"
)

// BuildExcelFromData product excel from data
func (lgc *Logics) BuildExcelFromData(ctx context.Context, objID string, fields map[string]Property, filter []string, data []mapstr.MapStr, xlsxFile *xlsx.File, header http.Header, meta metadata.Metadata) error {

	ccLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))
	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	sheet, err := xlsxFile.AddSheet("inst")
	if err != nil {
		blog.Errorf("setExcelRowDataByIndex add excel sheet error, rid:%s", err.Error(), util.GetHTTPCCRequestID(header))
		return err

	}
	addSystemField(fields, common.BKInnerObjIDObject, ccLang)

	if 0 == len(filter) {
		filter = getFilterFields(objID)
	} else {
		filter = append(filter, getFilterFields(objID)...)
	}

	instPrimaryKeyValMap := make(map[int64][]PropertyPrimaryVal)
	productExcelHealer(fields, filter, sheet, ccLang)
	//indexID := getFieldsIDIndexMap(fields)

	rowIndex := common.HostAddMethodExcelIndexOffset

	for _, rowMap := range data {

		instIDKey := metadata.GetInstIDFieldByObjID(objID)
		instID, err := rowMap.Int64(instIDKey)
		if err != nil {
			blog.Errorf("setExcelRowDataByIndex inst:%+v, not inst id key:%s, objID:%s, rid:%s", rowMap, instIDKey, objID, util.GetHTTPCCRequestID(header))
			return ccErr.Errorf(common.CCErrCommInstFieldNotFound, "instIDKey", objID)
		}

		primaryKeyArr := setExcelRowDataByIndex(rowMap, sheet, rowIndex, fields)
		instPrimaryKeyValMap[instID] = primaryKeyArr
		rowIndex++

	}

	err = lgc.BuildAssociationExcelFromData(ctx, objID, instPrimaryKeyValMap, xlsxFile, header, meta)
	if err != nil {
		return err
	}
	return nil
}

// BuildHostExcelFromData product excel from data
func (lgc *Logics) BuildHostExcelFromData(ctx context.Context, objID string, fields map[string]Property, filter []string, data []mapstr.MapStr, xlsxFile *xlsx.File, header http.Header, meta metadata.Metadata) error {
	ccLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))
	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	sheet, err := xlsxFile.AddSheet("host")
	if err != nil {
		blog.Errorf("BuildHostExcelFromData add excel sheet error, rid:%s", err.Error(), util.GetHTTPCCRequestID(header))
		return err
	}
	extFieldsTopoID := "cc_ext_field_topo"
	extFields := map[string]string{
		extFieldsTopoID: ccLang.Language("web_ext_field_topo"),
	}
	fields = addExtFields(fields, extFields)
	addSystemField(fields, common.BKInnerObjIDHost, ccLang)

	productExcelHealer(fields, filter, sheet, ccLang)

	instPrimaryKeyValMap := make(map[int64][]PropertyPrimaryVal)
	//indexID := getFieldsIDIndexMap(fields)
	rowIndex := common.HostAddMethodExcelIndexOffset
	for _, hostData := range data {

		rowMap, err := mapstr.NewFromInterface(hostData[common.BKInnerObjIDHost])
		if err != nil {
			msg := fmt.Sprintf("data format error:%v", hostData)
			blog.Errorf(msg)
			return errors.New(msg)
		}
		moduleMap, ok := hostData[common.BKInnerObjIDModule].([]interface{})
		if ok {
			topo := util.GetStrValsFromArrMapInterfaceByKey(moduleMap, "TopModuleName")
			rowMap[extFieldsTopoID] = strings.Join(topo, "\n")
		}

		instIDKey := metadata.GetInstIDFieldByObjID(objID)
		instID, err := rowMap.Int64(instIDKey)
		if err != nil {
			blog.Errorf("setExcelRowDataByIndex inst:%+v, not inst id key:%s, objID:%s, rid:%s", rowMap, instIDKey, objID, util.GetHTTPCCRequestID(header))
			return ccErr.Errorf(common.CCErrCommInstFieldNotFound, "instIDKey", objID)
		}
		primaryKeyArr := setExcelRowDataByIndex(rowMap, sheet, rowIndex, fields)
		instPrimaryKeyValMap[instID] = primaryKeyArr
		rowIndex++

	}

	err = lgc.BuildAssociationExcelFromData(ctx, objID, instPrimaryKeyValMap, xlsxFile, header, meta)
	if err != nil {
		return err
	}
	return nil
}

func (lgc *Logics) BuildAssociationExcelFromData(ctx context.Context, objID string, instPrimaryInfo map[int64][]PropertyPrimaryVal, xlsxFile *xlsx.File, header http.Header, meta metadata.Metadata) error {
	var instIDArr []int64
	for instID, _ := range instPrimaryInfo {
		instIDArr = append(instIDArr, instID)
	}
	instAsst, err := lgc.fetchAssocationData(ctx, header, objID, instIDArr)
	if err != nil {
		return err
	}
	asstData, err := lgc.getAssociationData(ctx, header, objID, instAsst, meta)
	if err != nil {
		return err
	}

	sheet, err := xlsxFile.AddSheet("assocation")
	if err != nil {
		blog.Errorf("setExcelRowDataByIndex add excel  assocation sheet error, rid:%s", err.Error(), util.GetHTTPCCRequestID(header))
		return err
	}
	productExcelAssociationHealer(sheet, lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header)))

	rowIndex := common.HostAddMethodExcelAssociationIndexOffset
	for _, inst := range instAsst {
		sheet.Cell(rowIndex, 0).SetString(inst.ObjectAsstID)
		sheet.Cell(rowIndex, 1).SetString("")
		srcInst, ok := instPrimaryInfo[inst.InstID]
		if !ok {
			blog.Warnf("BuildAssociationExcelFromData association inst:%+v, not inst id :%d, objID:%s, rid:%s", inst, inst.InstID, objID, util.GetHTTPCCRequestID(header))
			//return lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).Errorf(common.CCErrCommInstDataNil, fmt.Sprintf("%s %d", objID, inst.InstID))
			continue
		}
		dstInst, ok := asstData[inst.AsstObjectID][inst.AsstInstID]
		if !ok {
			blog.Warnf("BuildAssociationExcelFromData association inst:%+v, not inst id :%d, objID:%s, rid:%s", inst, inst.InstID, inst.AsstObjectID, util.GetHTTPCCRequestID(header))
			continue
		}
		sheet.Cell(rowIndex, 2).SetString(buildEexcelPrimaryKey(srcInst))
		sheet.Cell(rowIndex, 3).SetString(buildEexcelPrimaryKey(dstInst))
		style := sheet.Cell(rowIndex, 2).GetStyle()
		style.Alignment.WrapText = true
		style = sheet.Cell(rowIndex, 3).GetStyle()
		style.Alignment.WrapText = true
		rowIndex++
	}

	return nil

}

func buildEexcelPrimaryKey(propertyArr []PropertyPrimaryVal) string {
	var contentArr []string
	for _, property := range propertyArr {
		contentArr = append(contentArr, buildExcelPrimaryStr(property))
	}
	return strings.Join(contentArr, common.ExcelAsstPrimaryKeySplitChar)
}

func buildExcelPrimaryStr(property PropertyPrimaryVal) string {
	return property.Name + common.ExcelAsstPrimaryKeyJoinChar + property.StrVal
}

//BuildExcelTemplate  return httpcode, error
func (lgc *Logics) BuildExcelTemplate(objID, filename string, header http.Header, defLang lang.DefaultCCLanguageIf, meta metadata.Metadata) error {
	filterFields := getFilterFields(objID)
	fields, err := lgc.GetObjFieldIDs(objID, filterFields, header, meta)
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
func GetExcelData(sheet *xlsx.Sheet, fields map[string]Property, defFields common.KvMap, isCheckHeader bool, firstRow int, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	var err error
	nameIndexMap, err := checkExcelHealer(sheet, fields, isCheckHeader, defLang)
	if nil != err {
		return nil, nil, err
	}
	hosts := make(map[int]map[string]interface{})
	index := headerRow
	if 0 != firstRow {
		index = firstRow
	}
	errMsg := make([]string, 0)
	rowCnt := len(sheet.Rows)
	for ; index < rowCnt; index++ {
		row := sheet.Rows[index]
		host, getErr := getDataFromByExcelRow(row, index, fields, defFields, nameIndexMap, defLang)
		if 0 != len(getErr) {
			errMsg = append(errMsg, getErr...)
			continue
		}
		if 0 == len(host) {
			hosts[index+1] = nil
		} else {
			hosts[index+1] = host
		}
	}
	if 0 != len(errMsg) {
		return nil, errMsg, nil
	}

	return hosts, nil, nil

}

// GetExcelData excel数据，一个kv结构，key行数（excel中的行数），value内容
func GetRawExcelData(sheet *xlsx.Sheet, defFields common.KvMap, firstRow int, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	var err error
	nameIndexMap, err := checkExcelHealer(sheet, nil, false, defLang)
	if nil != err {
		return nil, nil, err
	}
	hosts := make(map[int]map[string]interface{})
	index := headerRow
	if 0 != firstRow {
		index = firstRow
	}
	errMsg := make([]string, 0)
	rowCnt := len(sheet.Rows)
	for ; index < rowCnt; index++ {
		row := sheet.Rows[index]
		host, getErr := getDataFromByExcelRow(row, index, nil, defFields, nameIndexMap, defLang)
		if nil != getErr {
			errMsg = append(errMsg, getErr...)
			continue
		}
		if 0 == len(host) {
			hosts[index+1] = nil
		} else {
			hosts[index+1] = host
		}
	}
	if 0 != len(errMsg) {
		return nil, errMsg, nil
	}

	return hosts, nil, nil

}

func GetAssociationExcelData(sheet *xlsx.Sheet, firstRow int) map[int]metadata.ExcelAssocation {

	rowCnt := len(sheet.Rows)
	index := firstRow

	asstInfoArr := make(map[int]metadata.ExcelAssocation, 0)
	for ; index < rowCnt; index++ {
		row := sheet.Rows[index]
		op := row.Cells[associationOPColIndex].String()
		if op == "" {
			continue
		}

		asstObjID := row.Cells[assciationAsstObjIDIndex].String()
		srcInst := row.Cells[assciationSrcInstIndex].String()
		dstInst := row.Cells[assciationDstInstIndex].String()
		asstInfoArr[index] = metadata.ExcelAssocation{
			ObjectAsstID: asstObjID,
			Operate:      getAssociationExcelOperateFlag(op),
			SrcPrimary:   srcInst,
			DstPrimary:   dstInst,
		}
	}

	return asstInfoArr
}

//GetFilterFields 不需要展示字段
func GetFilterFields(objID string) []string {
	return getFilterFields(objID)
}

func getAssociationExcelOperateFlag(op string) metadata.ExcelAssocationOperate {
	opFlag := metadata.ExcelAssocationOperateError
	switch op {
	case associationOPAdd:
		opFlag = metadata.ExcelAssocationOperateAdd
	case associationOPDelete:
		opFlag = metadata.ExcelAssocationOperateDelete
	}

	return opFlag
}
