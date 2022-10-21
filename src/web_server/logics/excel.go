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
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"
)

// BuildExcelFromData product excel from data
func (lgc *Logics) BuildExcelFromData(ctx context.Context, objID string, fields map[string]Property, filter []string,
	data []mapstr.MapStr, xlsxFile *xlsx.File, header http.Header, modelBizID int64, usernameMap map[string]string,
	propertyList []string, org []metadata.DepartmentItem, orgPropertyList []string,
	asstObjectUniqueIDMap map[string]int64, selfObjectUniqueID int64) error {

	rid := util.GetHTTPCCRequestID(header)

	ccLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))
	sheet, err := xlsxFile.AddSheet("inst")
	if err != nil {
		blog.Errorf("setExcelRowDataByIndex add excel sheet error, err:%s, rid:%s", err.Error(), rid)
		return err

	}
	// index=1 表格数据的起始索引，excel表格数据第一列为字段说明，第二列为数据列
	addSystemField(fields, common.BKInnerObjIDObject, ccLang, 1)

	if len(filter) == 0 {
		filter = getFilterFields(objID)
	} else {
		filter = append(filter, getFilterFields(objID)...)
	}

	productExcelHeader(ctx, fields, filter, xlsxFile, sheet, ccLang)
	// indexID := getFieldsIDIndexMap(fields)

	rowIndex := common.HostAddMethodExcelIndexOffset
	instIDArr := make([]int64, 64)

	for _, rowMap := range data {

		instIDKey := metadata.GetInstIDFieldByObjID(objID)
		instID, err := rowMap.Int64(instIDKey)
		if err != nil {
			blog.Errorf("parse inst(%+v) id(key:%s) failed, err: %v, objID: %s, rid: %s", rowMap, instIDKey, err,
				objID, rid)
		}
		// 使用中英文用户名重新构造用户列表(用户列表实际为逗号分隔的string型)
		rowMap, err = replaceEnName(rid, rowMap, usernameMap, propertyList, ccLang)
		if err != nil {
			blog.Errorf("rebuild user list failed, err: %v, rid: %s", err, rid)
			return err
		}

		rowMap, err = replaceDepartmentFullName(rid, rowMap, org, orgPropertyList, ccLang)
		if err != nil {
			blog.Errorf("rebuild organization list failed, err: %v, rid: %s", err, rid)
			return err
		}

		setExcelRowDataByIndex(rowMap, sheet, rowIndex, fields)

		instIDArr = append(instIDArr, instID)
		rowIndex++

	}

	err = lgc.BuildAssociationExcelFromData(ctx, objID, instIDArr, xlsxFile, header, modelBizID,
		asstObjectUniqueIDMap, selfObjectUniqueID)
	if err != nil {
		return err
	}
	return nil
}

// BuildHostExcelFromData product excel from data
// selfObjectUniqueID 当前导出对象使用的唯一校验id
func (lgc *Logics) BuildHostExcelFromData(ctx context.Context, objID string, fields map[string]Property,
	filter []string, data []mapstr.MapStr, xlsxFile *xlsx.File, header http.Header, modelBizID int64,
	usernameMap map[string]string, propertyList []string, objNames, objIDs []string, org []metadata.DepartmentItem,
	orgPropertyList []string, asstObjectUniqueIDMap map[string]int64, selfObjectUniqueID int64) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	ccLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))
	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	sheet, err := xlsxFile.AddSheet("host")
	if err != nil {
		blog.Errorf("add excel sheet failed, err: %v, rid: %s", err, rid)
		return err
	}

	extFieldsTopoID := "cc_ext_field_topo"
	extFieldsBizID := "cc_ext_biz"
	extFieldsModuleID := "cc_ext_module"
	extFieldsSetID := "cc_ext_set"
	extFieldKey := make([]string, 0)
	extFieldKey = append(extFieldKey, extFieldsTopoID, extFieldsBizID)

	extFields := map[string]string{
		extFieldsTopoID:   ccLang.Language("web_ext_field_topo"),
		extFieldsBizID:    ccLang.Language("biz_property_bk_biz_name"),
		extFieldsModuleID: ccLang.Language("web_ext_field_module_name"),
		extFieldsSetID:    ccLang.Language("web_ext_field_set_name"),
	}
	// 生成key,用于赋值遍历主机数据进行赋值
	for _, objID := range objIDs {
		extFieldKey = append(extFieldKey, "cc_ext_"+objID)
	}

	// 2 自定义层级名称在extFieldKey切片中起始位置为2，0,1索引为业务拓扑和业务名
	for idx, objName := range objNames {
		extFields[extFieldKey[idx+2]] = objName
	}

	extFieldKey = append(extFieldKey, extFieldsSetID, extFieldsModuleID)
	fields = addExtFields(fields, extFields, extFieldKey)
	// len(objNames)+5=tip + biztopo + biz + set + moudle + customLen, the former indexes is used by these columns
	addSystemField(fields, common.BKInnerObjIDHost, ccLang, len(objNames)+5)

	cloudAreaArr, _, err := lgc.getCloudArea(ctx, header)
	if err != nil {
		blog.Errorf("build host excel data failed, err: %v, rid: %s", err, rid)
		return err
	}

	productHostExcelHeader(ctx, fields, filter, xlsxFile, sheet, ccLang, objNames, cloudAreaArr)

	handleHostDataParam := &HandleHostDataParam{
		HostData:          data,
		ExtFieldsTopoID:   extFieldsTopoID,
		ExtFieldsBizID:    extFieldsBizID,
		ExtFieldsModuleID: extFieldsModuleID,
		ExtFieldsSetID:    extFieldsSetID,
		CcErr:             ccErr,
		ExtFieldKey:       extFieldKey,
		UsernameMap:       usernameMap,
		PropertyList:      propertyList,
		Organization:      org,
		OrgPropertyList:   orgPropertyList,
		CcLang:            ccLang,
		Sheet:             sheet,
		Rid:               rid,
		ObjID:             objID,
		ObjIDs:            objIDs,
		Fields:            fields,
	}

	instIDs, err := lgc.buildHostExcelData(handleHostDataParam)
	if err != nil {
		blog.Errorf("build host excel data failed, err: %v, rid: %s", err, rid)
		return err
	}

	err = lgc.BuildAssociationExcelFromData(ctx, objID, instIDs, xlsxFile, header, modelBizID,
		asstObjectUniqueIDMap, selfObjectUniqueID)
	if err != nil {
		blog.Errorf("build association excel data failed, err: %v, rid: %s", err, rid)
		return err
	}
	return nil
}

// buildHostExcelData 处理主机数据，生成Excel表格数据
func (lgc *Logics) buildHostExcelData(handleHostDataParam *HandleHostDataParam) ([]int64, error) {
	instIDArr := make([]int64, 0)
	rowIndex := common.HostAddMethodExcelIndexOffset
	for _, hostData := range handleHostDataParam.HostData {
		rowMap, err := mapstr.NewFromInterface(hostData[common.BKInnerObjIDHost])
		if err != nil {
			blog.Errorf("build host excel data failed, hostData: %#v, err: %v, rid: %s", hostData, err,
				handleHostDataParam.Rid)
			return nil, handleHostDataParam.CcErr.CCError(common.CCErrCommReplyDataFormatError)
		}

		// handle custom extFieldKey,前两个元素为业务拓扑、业务，后两个元素为集群、模块，中间的为自定义层级列
		for idx, field := range handleHostDataParam.ExtFieldKey[2 : len(handleHostDataParam.ExtFieldKey)-2] {
			rowMap[field] = hostData[handleHostDataParam.ObjIDs[idx]]
		}
		rowMap[handleHostDataParam.ExtFieldsSetID] = hostData["sets"]
		rowMap[handleHostDataParam.ExtFieldsModuleID] = hostData["modules"]

		if _, exist := handleHostDataParam.Fields[common.BKCloudIDField]; exist {
			cloudAreaArr, err := rowMap.MapStrArray(common.BKCloudIDField)
			if err != nil {
				blog.Errorf("get cloud id failed, host: %#v, err: %v, rid: %s", hostData, err, handleHostDataParam.Rid)
				return nil, handleHostDataParam.CcErr.CCError(common.CCErrCommReplyDataFormatError)
			}

			if len(cloudAreaArr) != 1 {
				blog.Errorf("host has many cloud areas, host: %#v, err: %v, rid: %s", hostData, err,
					handleHostDataParam.Rid)
				return nil, handleHostDataParam.CcErr.CCError(common.CCErrCommReplyDataFormatError)
			}

			cloudArea := cloudAreaArr[0][common.BKInstNameField]
			rowMap.Set(common.BKCloudIDField, cloudArea)
		}

		moduleMap, ok := hostData[common.BKInnerObjIDModule].([]interface{})
		if ok {
			topos := util.GetStrValsFromArrMapInterfaceByKey(moduleMap, "TopModuleName")
			if len(topos) > 0 {
				idx := strings.Index(topos[0], logics.SplitFlag)
				if idx > 0 {
					rowMap[handleHostDataParam.ExtFieldsBizID] = topos[0][:idx]
				}

				toposNobiz := make([]string, 0)
				for _, topo := range topos {
					idx := strings.Index(topo, logics.SplitFlag)
					if idx > 0 && len(topo) >= idx+len(logics.SplitFlag) {
						toposNobiz = append(toposNobiz, topo[idx+len(logics.SplitFlag):])
					}
				}
				rowMap[handleHostDataParam.ExtFieldsTopoID] = strings.Join(toposNobiz, ", ")
			}
		}

		instIDKey := metadata.GetInstIDFieldByObjID(handleHostDataParam.ObjID)
		instID, err := rowMap.Int64(instIDKey)
		if err != nil {
			blog.Errorf("get inst id failed, inst: %#v, err: %v, rid: %s", rowMap, err, handleHostDataParam.Rid)
			return nil, handleHostDataParam.CcErr.Errorf(common.CCErrCommInstFieldNotFound, instIDKey,
				handleHostDataParam.ObjID)
		}

		// 使用中英文用户名重新构造用户列表(用户列表实际为逗号分隔的string型)
		rowMap, err = replaceEnName(handleHostDataParam.Rid, rowMap, handleHostDataParam.UsernameMap,
			handleHostDataParam.PropertyList, handleHostDataParam.CcLang)
		if err != nil {
			blog.Errorf("rebuild user list field, err: %v, rid: %s", err, handleHostDataParam.Rid)
			return nil, err
		}

		rowMap, err = replaceDepartmentFullName(handleHostDataParam.Rid, rowMap, handleHostDataParam.Organization,
			handleHostDataParam.OrgPropertyList, handleHostDataParam.CcLang)
		if err != nil {
			blog.Errorf("rebuild organization list failed, err: %v, rid: %s", err, handleHostDataParam.Rid)
			return nil, err
		}

		setExcelRowDataByIndex(rowMap, handleHostDataParam.Sheet, rowIndex, handleHostDataParam.Fields)
		instIDArr = append(instIDArr, instID)
		rowIndex++
	}
	return instIDArr, nil
}

// getObjectAssociation 获取模型关联关系
func (lgc *Logics) getObjectAssociation(ctx context.Context, header http.Header, objID string, modelBizID int64) (
	[]*metadata.Association, errors.CCErrorCoder) {
	rid := util.ExtractRequestIDFromContext(ctx)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	cond := &metadata.SearchAssociationObjectRequest{
		Condition: map[string]interface{}{
			condition.BKDBOR: []mapstr.MapStr{
				{
					common.BKObjIDField: objID,
				},
				{
					common.BKAsstObjIDField: objID,
				},
			},
		},
	}
	// 确定关联标识的列表，定义excel选项下拉栏。此处需要查cc_ObjAsst表。
	resp, err := lgc.CoreAPI.ApiServer().SearchObjectAssociation(ctx, header, cond)
	if err != nil {
		blog.ErrorJSON("get object association list failed, err: %v, rid: %s", err, rid)
		return nil, defErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if err := resp.CCError(); err != nil {
		blog.ErrorJSON("get object association list failed, err: %v, rid: %s", resp.ErrMsg, rid)
		return nil, err
	}

	return resp.Data, nil
}

// BuildAssociationExcelFromData build association excel
func (lgc *Logics) BuildAssociationExcelFromData(ctx context.Context, objID string, instIDArr []int64,
	xlsxFile *xlsx.File, header http.Header, modelBizID int64, asstObjectUniqueIDMap map[string]int64,
	selfObjectUniqueID int64) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	defLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))

	sheet, err := xlsxFile.AddSheet("association")
	if err != nil {
		blog.Errorf("add excel association sheet failed, err: %v, rid:%s", err, rid)
		return err
	}

	asstList, err := lgc.getObjectAssociation(ctx, header, objID, modelBizID)
	if err != nil {
		return err
	}

	// 未设置, 不导出关联关系数据
	if len(asstObjectUniqueIDMap) == 0 {
		productExcelAssociationHeader(ctx, sheet, defLang, 0, asstList)
		return nil
	}

	needAsst := make(map[string]struct{})
	for key := range asstObjectUniqueIDMap {
		needAsst[key] = struct{}{}
	}

	// 2021年06月01日， 判断时候需要导出自关联，这个处理就是打补丁，不友好
	hasSelfAssociation := false
	if _, ok := asstObjectUniqueIDMap[objID]; ok {
		hasSelfAssociation = true
	}

	var asstIDArr []string
	// 单独处理自关联模型, 将自关联模型找到后直接删除该字段，之后的模型关联就不会到受到影响导致所有的关联关系都导出
	for _, asst := range asstList {
		if asst.ObjectID == asst.AsstObjID && hasSelfAssociation {
			asstIDArr = append(asstIDArr, asst.AssociationName)
			delete(needAsst, asst.ObjectID)
			break
		}
	}

	for _, asst := range asstList {

		_, ok := needAsst[asst.ObjectID]
		if ok {
			asstIDArr = append(asstIDArr, asst.AssociationName)
			continue
		}
		_, ok = needAsst[asst.AsstObjID]
		if ok {
			asstIDArr = append(asstIDArr, asst.AssociationName)
			continue
		}
	}

	instAsst, err := lgc.fetchAssociationData(ctx, header, objID, instIDArr, modelBizID, asstIDArr, hasSelfAssociation)
	if err != nil {
		return err
	}
	asstData, objInstData, err := lgc.getAssociationData(ctx, header, objID, instAsst, modelBizID,
		asstObjectUniqueIDMap, selfObjectUniqueID)
	if err != nil {
		return err
	}

	productExcelAssociationHeader(ctx, sheet, defLang, len(instAsst), asstList)

	rowIndex := common.HostAddMethodExcelAssociationIndexOffset

	for _, inst := range instAsst {
		sheet.Cell(rowIndex, 1).SetString(inst.ObjectAsstID)
		sheet.Cell(rowIndex, 2).SetString("")

		srcInst, dstInst := buildRowInfo(objID, rid, inst, asstData, objInstData)
		if srcInst == nil || dstInst == nil {
			continue
		}

		// TODO: 注意源和目标顺序
		sheet.Cell(rowIndex, 3).SetString(buildExcelPrimaryKey(srcInst))
		sheet.Cell(rowIndex, 4).SetString(buildExcelPrimaryKey(dstInst))
		style := sheet.Cell(rowIndex, 3).GetStyle()
		style.Alignment.WrapText = true
		style = sheet.Cell(rowIndex, 4).GetStyle()
		style.Alignment.WrapText = true
		rowIndex++
	}

	return nil
}

func buildRowInfo(objID, rid string, inst *metadata.InstAsst, asstData map[string]map[int64][]PropertyPrimaryVal,
	objInstData map[int64][]PropertyPrimaryVal) (srcInst []PropertyPrimaryVal, dstInst []PropertyPrimaryVal) {
	if inst == nil {
		return
	}
	var ok bool
	if inst.ObjectID == objID {
		srcInst, ok = objInstData[inst.InstID]
		if !ok {
			blog.WarnJSON("association inst:%s, not inst id :%d, objID:%s, rid:%s",
				inst, inst.InstID, objID, rid)
			return
		}
		dstInst, ok = asstData[inst.AsstObjectID][inst.AsstInstID]
		if !ok {
			blog.WarnJSON("association inst:%s, not inst id :%d, objID:%s, rid:%s",
				inst, inst.InstID, inst.AsstObjectID, rid)
			return
		}
	} else {
		srcInst, ok = asstData[inst.ObjectID][inst.InstID]
		if !ok {
			blog.WarnJSON("association inst:%s, not inst id :%d, objID:%s, rid:%s",
				inst, inst.InstID, objID, rid)
			return
		}
		dstInst, ok = objInstData[inst.AsstInstID]
		if !ok {
			blog.WarnJSON("association inst:%s, not inst id :%d, objID:%s, rid:%s",
				inst, inst.InstID, inst.AsstObjectID, rid)
			return
		}
	}
	return
}

func buildExcelPrimaryKey(propertyArr []PropertyPrimaryVal) string {
	var contentArr []string
	for _, property := range propertyArr {
		contentArr = append(contentArr, buildExcelPrimaryStr(property))
	}
	return strings.Join(contentArr, common.ExcelAsstPrimaryKeySplitChar)
}

func buildExcelPrimaryStr(property PropertyPrimaryVal) string {
	return property.Name + common.ExcelAsstPrimaryKeyJoinChar + property.StrVal
}

// BuildExcelTemplate  return httpcode, error
func (lgc *Logics) BuildExcelTemplate(ctx context.Context, objID, filename string, header http.Header,
	defLang lang.DefaultCCLanguageIf, modelBizID int64) error {

	rid := util.GetHTTPCCRequestID(header)
	filterFields := getFilterFields(objID)

	fields, err := lgc.GetObjFieldIDs(objID, filterFields, nil, header, modelBizID,
		common.HostAddMethodExcelDefaultIndex)

	if err != nil {
		blog.Errorf("get %s fields error:%s, rid: %s", objID, err.Error(), rid)
		return err
	}

	var file *xlsx.File
	file = xlsx.NewFile()
	sheet, err := file.AddSheet(objID)
	if err != nil {
		blog.Errorf("get %s fields error: %v, rid: %s", objID, err, rid)
		return err
	}

	asstSheet, err := file.AddSheet("association")
	if err != nil {
		blog.Errorf("setExcelRowDataByIndex add excel  association sheet error. err:%s, rid:%s", err.Error(), rid)
		return err
	}
	asstList, err := lgc.getObjectAssociation(ctx, header, objID, modelBizID)
	if err != nil {
		return err
	}
	productExcelAssociationHeader(ctx, asstSheet, defLang, 0, asstList)

	blog.V(5).Infof("BuildExcelTemplate fields count:%d, rid: %s", fields, rid)
	if objID == common.BKInnerObjIDHost {
		cloudAreaName, _, err := lgc.getCloudArea(ctx, header)
		if err != nil {
			blog.Errorf("build %s excel template failed, err: %v,  rid: %s", objID, err, rid)
			return err
		}
		productHostExcelHeader(ctx, fields, filterFields, file, sheet, defLang, nil, cloudAreaName)
	} else {
		productExcelHeader(ctx, fields, filterFields, file, sheet, defLang)
	}
	ProductExcelCommentSheet(ctx, file, defLang)

	if err = file.Save(filename); nil != err {
		blog.Errorf("save file failed, filename: %s, err: %+v, rid: %s", filename, err, rid)
		return err
	}

	return nil
}

// AddDownExcelHttpHeader TODO
func AddDownExcelHttpHeader(c *gin.Context, name string) {
	if strings.HasSuffix(name, ".xls") {
		c.Header("Content-Type", "application/vnd.ms-excel")
	} else {
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	}
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Disposition", "attachment; filename="+name) // 文件名
	c.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

// GetExcelData excel数据，一个kv结构，key行数（excel中的行数），value内容
func GetExcelData(ctx context.Context, sheet *xlsx.Sheet, fields map[string]Property, defFields common.KvMap,
	isCheckHeader bool, firstRow int, defLang lang.DefaultCCLanguageIf, department map[int64]metadata.DepartmentItem) (
	map[int]map[string]interface{}, []string, error) {

	var err error
	nameIndexMap, err := checkExcelHeader(ctx, sheet, fields, isCheckHeader, defLang)
	if err != nil {
		return nil, nil, err
	}
	hosts := make(map[int]map[string]interface{})
	index := headerRow
	if firstRow != 0 {
		index = firstRow
	}
	errMsg := make([]string, 0)
	rowCnt := len(sheet.Rows)
	for ; index < rowCnt; index++ {
		row := sheet.Rows[index]
		host, getErr := getDataFromByExcelRow(ctx, row, index, fields, defFields, nameIndexMap, defLang, department)
		if len(getErr) != 0 {
			errMsg = append(errMsg, getErr...)
			continue
		}
		if len(host) != 0 {
			hosts[index+1] = host
		}
	}
	if len(errMsg) != 0 {
		return nil, errMsg, nil
	}

	return hosts, nil, nil

}

// GetRawExcelData excel数据，一个kv结构，key行数（excel中的行数），value内容
func GetRawExcelData(ctx context.Context, sheet *xlsx.Sheet, defFields common.KvMap, firstRow int,
	defLang lang.DefaultCCLanguageIf, department map[int64]metadata.DepartmentItem) (map[int]map[string]interface{},
	[]string, error) {

	var err error
	nameIndexMap, err := checkExcelHeader(ctx, sheet, nil, false, defLang)
	if nil != err {
		return nil, nil, err
	}
	hosts := make(map[int]map[string]interface{})
	index := headerRow
	if firstRow != 0 {
		index = firstRow
	}
	errMsg := make([]string, 0)
	rowCnt := len(sheet.Rows)
	for ; index < rowCnt; index++ {
		row := sheet.Rows[index]
		host, getErr := getDataFromByExcelRow(ctx, row, index, nil, defFields, nameIndexMap, defLang, department)
		if getErr != nil {
			errMsg = append(errMsg, getErr...)
			continue
		}
		if len(host) == 0 {
			hosts[index+1] = nil
		} else {
			hosts[index+1] = host
		}
	}
	if len(errMsg) != 0 {
		return nil, errMsg, nil
	}

	return hosts, nil, nil

}

// GetAssociationExcelData read sheet of association data from excel
func GetAssociationExcelData(sheet *xlsx.Sheet, firstRow int, defLang lang.DefaultCCLanguageIf) (
	map[int]metadata.ExcelAssociation, []metadata.RowMsgData) {

	rowCnt := len(sheet.Rows)
	index := firstRow

	asstInfoArr := make(map[int]metadata.ExcelAssociation)
	errMsg := make([]metadata.RowMsgData, 0)
	for ; index < rowCnt; index++ {
		row := sheet.Rows[index]

		// 获取单元格内容，使用for循环防止直接获取对应单元格数据导致数组越界
		var asstObjID, op, srcInst, dstInst string
		for index, item := range row.Cells {
			switch index {
			case associationAsstObjIDIndex:
				asstObjID = item.String()
			case associationOPColIndex:
				op = item.String()
			case associationSrcInstIndex:
				srcInst = item.String()
			case associationDstInstIndex:
				dstInst = item.String()
			}
		}

		if op == "" {
			continue
		}

		if asstObjID == "" || srcInst == "" || dstInst == "" {
			err := defLang.Languagef("web_excel_row_handle_error", sheet.Name, (index + 1))
			errMsg = append(errMsg, metadata.RowMsgData{Row: index, Msg: err})
			continue
		}

		asstInfoArr[index] = metadata.ExcelAssociation{
			ObjectAsstID: asstObjID,
			Operate:      getAssociationExcelOperateFlag(op),
			SrcPrimary:   srcInst,
			DstPrimary:   dstInst,
		}
	}

	return asstInfoArr, errMsg
}

// StatisticsAssociation TODO
func StatisticsAssociation(sheet *xlsx.Sheet, firstRow int) ([]string, map[string]metadata.ObjectAsstIDStatisticsInfo) {

	rowCnt := len(sheet.Rows)
	index := firstRow
	asstInfoMap := make(map[string]metadata.ObjectAsstIDStatisticsInfo, 0)
	// bk_obj_asst_id
	asstNameArr := make([]string, 0)
	for ; index < rowCnt; index++ {
		var asstInfo metadata.ObjectAsstIDStatisticsInfo
		var ok bool

		row := sheet.Rows[index]

		if len(row.Cells) <= associationOPColIndex {
			continue
		}
		// bk_obj_asst_id 的值
		asstObjID := row.Cells[associationAsstObjIDIndex].String()
		asstInfo, ok = asstInfoMap[asstObjID]
		if !ok {
			// 第一次出现这个关联关系的唯一标识
			asstNameArr = append(asstNameArr, asstObjID)
			asstInfo = metadata.ObjectAsstIDStatisticsInfo{}
		}
		asstInfo.Total += 1
		op := row.Cells[associationOPColIndex].String()
		if op == "" {
			continue
		} else {
			operate := getAssociationExcelOperateFlag(op)
			switch operate {
			case metadata.ExcelAssociationOperateDelete:
				asstInfo.Delete += 1
			case metadata.ExcelAssociationOperateAdd:
				asstInfo.Create += 1
			}
		}
		asstInfoMap[asstObjID] = asstInfo
	}

	return asstNameArr, asstInfoMap
}

// GetFilterFields 不需要展示字段
func GetFilterFields(objID string) []string {
	return getFilterFields(objID)
}

// GetCustomFields 用户展示字段export时优先排序
func GetCustomFields(filterFields []string, customFields []string) []string {
	return getCustomFields(filterFields, customFields)
}

func getAssociationExcelOperateFlag(op string) metadata.ExcelAssociationOperate {
	opFlag := metadata.ExcelAssociationOperateError
	switch op {
	case associationOPAdd:
		opFlag = metadata.ExcelAssociationOperateAdd
	case associationOPDelete:
		opFlag = metadata.ExcelAssociationOperateDelete
	}

	return opFlag
}
