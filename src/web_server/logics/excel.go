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
	"fmt"
	"net/http"
	"strings"

	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"
	"configcenter/src/web_server/app/options"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx/v3"
)

// BuildExcelFromData product excel from data
func (lgc *Logics) BuildExcelFromData(ctx context.Context, objID string, fields map[string]Property, filter []string,
	data []mapstr.MapStr, xlsxFile *xlsx.File, header http.Header, modelBizID int64, usernameMap map[string]string,
	propertyList []string, org []metadata.DepartmentItem, orgPropertyList []string,
	asstObjectUniqueIDMap map[string]int64, selfObjectUniqueID int64, rowCountArr []int) error {

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

	if err := productExcelHeader(ctx, fields, filter, xlsxFile, sheet, ccLang); err != nil {
		return err
	}

	rowIndex := common.AddExcelDataIndexOffset
	instIDArr := make([]int64, 64)

	for idx, rowMap := range data {

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

		rowMap, err = replaceEnumMultiName(rid, rowMap, fields)
		if err != nil {
			blog.Errorf("rebuild enum multi failed, err: %v, rid: %s", err, rid)
			return err
		}

		if err := setExcelRowDataByIndex(rowMap, sheet, rowIndex, fields, rowCountArr[idx]); err != nil {
			return err
		}

		instIDArr = append(instIDArr, instID)
		rowIndex += rowCountArr[idx]

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
// NOCC:golint/fnsize(后续重构处理)
func (lgc *Logics) BuildHostExcelFromData(c *gin.Context, objID string, fields map[string]Property,
	filter []string, xlsxFile *xlsx.File, header http.Header, objNames, objIDs []string,
	input *metadata.ExcelExportHostInput, config *options.Config) error {

	ctx := util.NewContextFromGinContext(c)
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

	if err := productHostExcelHeader(ctx, fields, filter, xlsxFile, sheet, ccLang, objNames, cloudAreaArr); err != nil {
		blog.Errorf("product host excel header failed, err: %v, rid: %s", err, rid)
		return err
	}

	ids := make([]int64, 0)
	hostCount := input.ExportCond.Page.Limit + input.ExportCond.Page.Start
	limit := input.ExportCond.Page.Limit
	index := common.AddExcelDataIndexOffset
	for start := input.ExportCond.Page.Start; start < hostCount; start = start + common.BKMaxExportLimit {
		input.ExportCond.Page.Start = start
		if limit > common.BKMaxExportLimit {
			input.ExportCond.Page.Limit = common.BKMaxExportLimit
			limit = limit - common.BKMaxExportLimit
		} else {
			input.ExportCond.Page.Limit = limit
		}

		hostInfo, err := lgc.handleHostInfo(ctx, header, fields, objIDs, input)
		if err != nil {
			blog.Errorf("search and handle host info failed, err: %v, rid: %s", err, rid)
			return err
		}

		if len(hostInfo) == 0 {
			break
		}

		hostInfo, err = lgc.HandleExportEnumQuoteInst(c, header, hostInfo, objID, fields, rid)
		if err != nil {
			blog.Errorf("handle enum quote inst failed, err: %v, rid: %s", err, rid)
			return err
		}

		usernameMap, propertyList, err := lgc.GetUsernameMapWithPropertyList(c, objID, hostInfo, config)
		if err != nil {
			blog.Errorf("export host failed, get username map and property list failed, err: %+v, rid: %s", err, rid)
			return err
		}

		org, orgPropertyList, err := lgc.GetDepartmentDetail(c, objID, config, hostInfo)
		if err != nil {
			blog.Errorf("get department map and property list failed, err: %+v, rid: %s", err, rid)
			return err
		}

		hostInfo, rowCountArr, err := lgc.BuildDataWithTable(common.BKInnerObjIDHost, hostInfo, fields, header)
		if err != nil {
			blog.ErrorJSON("get data with table field data failed, hostInfo: %s, err: %s, rid: %s", hostInfo, err, rid)
			return err
		}

		handleHostDataParam := &HandleHostDataParam{
			HostData:          hostInfo,
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

		var instIDs []int64
		instIDs, index, err = lgc.buildHostExcelData(handleHostDataParam, rowCountArr, index)
		if err != nil {
			blog.Errorf("build host excel data failed, err: %v, rid: %s", err, rid)
			return err
		}
		ids = append(ids, instIDs...)
	}

	err = lgc.BuildAssociationExcelFromData(ctx, objID, ids, xlsxFile, header, input.AppID, input.AssociationCond,
		input.ObjectUniqueID)
	if err != nil {
		blog.Errorf("build association excel data failed, err: %v, rid: %s", err, rid)
		return err
	}
	return nil
}

// BuildDataWithTable build data with table data
func (lgc *Logics) BuildDataWithTable(objID string, infos []mapstr.MapStr, fields map[string]Property,
	header http.Header) ([]mapstr.MapStr, []int, error) {

	rowCountArr := make([]int, len(infos))
	for i := range rowCountArr {
		rowCountArr[i] = 1
	}
	// 1. 找出表格字段
	tableProperty := make([]Property, 0)
	for _, property := range fields {
		if property.PropertyType == common.FieldTypeInnerTable {
			tableProperty = append(tableProperty, property)
		}
	}
	if len(tableProperty) == 0 {
		return infos, rowCountArr, nil
	}

	ids := make([]int64, 0)
	dataMap := make(map[int64]map[string][]mapstr.MapStr)
	var err error
	for _, info := range infos {
		infoMapStr := info
		if objID == common.BKInnerObjIDHost {
			infoMapStr, err = info.MapStr(common.BKInnerObjIDHost)
			if err != nil {
				return nil, nil, fmt.Errorf("can not find %s data, err: %s", objID, err)
			}
		}

		id, err := infoMapStr.Int64(common.GetInstIDField(objID))
		if err != nil {
			return nil, nil, fmt.Errorf("data is invalid, err: %v", err)
		}
		ids = append(ids, id)
		dataMap[id] = make(map[string][]mapstr.MapStr)
	}

	// 2. 查询数据对应的表格字段的值
	queryOpt := metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: &filter.Expression{
				RuleFactory: &filter.CombinedRule{
					Condition: filter.And,
					Rules: []filter.RuleFactory{
						&filter.AtomRule{
							Field:    common.BKInstIDField,
							Operator: filter.OpFactory(filter.In),
							Value:    ids,
						},
					},
				},
			},
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	for _, property := range tableProperty {
		opt := &metadata.ListQuotedInstOption{ObjID: objID, PropertyID: property.ID, CommonQueryOption: queryOpt}
		instances, err := lgc.Engine.CoreAPI.ApiServer().ModelQuote().ListQuotedInstance(context.Background(), header,
			opt)
		if err != nil {
			return nil, nil, err
		}
		for _, inst := range instances.Info {
			instID, err := inst.Int64(common.BKInstIDField)
			if err != nil {
				return nil, nil, err
			}
			dataMap[instID][property.ID] = append(dataMap[instID][property.ID], inst)
		}
	}

	// 3. 整理返回带表格数据的结果, 以及每条数据需要占用excel多少行
	return getInfoAndRowCount(objID, infos, dataMap, rowCountArr)
}

func getInfoAndRowCount(objID string, infos []mapstr.MapStr, dataMap map[int64]map[string][]mapstr.MapStr,
	rowCountArr []int) ([]mapstr.MapStr, []int, error) {

	var err error
	for idx, info := range infos {
		infoMapStr := info
		if objID == common.BKInnerObjIDHost {
			infoMapStr, err = info.MapStr(common.BKInnerObjIDHost)
			if err != nil {
				return nil, nil, fmt.Errorf("can not find %s data, err: %s", objID, err)
			}
		}
		id, err := infoMapStr.Int64(common.GetInstIDField(objID))
		if err != nil {
			return nil, nil, fmt.Errorf("data is invalid, err: %v", err)
		}
		count := 1
		for propertyID, data := range dataMap[id] {
			infoMapStr[propertyID] = data
			if len(data) > count {
				count = len(data)
			}
		}
		if objID == common.BKInnerObjIDHost {
			info[objID] = infoMapStr
		}
		infos[idx] = info
		rowCountArr[idx] = count
	}

	return infos, rowCountArr, nil
}

// buildHostExcelData 处理主机数据，生成Excel表格数据
func (lgc *Logics) buildHostExcelData(handleHostDataParam *HandleHostDataParam, rowCountArr []int, rowIndex int) (
	[]int64, int, error) {

	instIDArr := make([]int64, 0)
	for idx, hostData := range handleHostDataParam.HostData {
		rowMap, err := mapstr.NewFromInterface(hostData[common.BKInnerObjIDHost])
		if err != nil {
			blog.Errorf("build host excel data failed, hostData: %#v, err: %v, rid: %s", hostData, err,
				handleHostDataParam.Rid)
			return nil, 0, handleHostDataParam.CcErr.CCError(common.CCErrCommReplyDataFormatError)
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
				return nil, 0, handleHostDataParam.CcErr.CCError(common.CCErrCommReplyDataFormatError)
			}

			if len(cloudAreaArr) != 1 {
				blog.Errorf("host has many cloud areas, host: %#v, err: %v, rid: %s", hostData, err,
					handleHostDataParam.Rid)
				return nil, 0, handleHostDataParam.CcErr.CCError(common.CCErrCommReplyDataFormatError)
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
			return nil, 0, handleHostDataParam.CcErr.Errorf(common.CCErrCommInstFieldNotFound, instIDKey,
				handleHostDataParam.ObjID)
		}

		rowMap, err = replaceData(rowMap, handleHostDataParam)
		if err != nil {
			return nil, 0, err
		}

		err = setExcelRowDataByIndex(rowMap, handleHostDataParam.Sheet, rowIndex, handleHostDataParam.Fields,
			rowCountArr[idx])
		if err != nil {
			return nil, 0, err
		}
		instIDArr = append(instIDArr, instID)
		rowIndex += rowCountArr[idx]
	}
	return instIDArr, rowIndex, nil
}

func replaceData(rowMap mapstr.MapStr, handleHostDataParam *HandleHostDataParam) (mapstr.MapStr, error) {
	// 使用中英文用户名重新构造用户列表(用户列表实际为逗号分隔的string型)
	rowMap, err := replaceEnName(handleHostDataParam.Rid, rowMap, handleHostDataParam.UsernameMap,
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

	rowMap, err = replaceEnumMultiName(handleHostDataParam.Rid, rowMap, handleHostDataParam.Fields)
	if err != nil {
		blog.Errorf("rebuild enum multi failed, err: %v, rid: %s", err, handleHostDataParam.Rid)
		return nil, err
	}

	return rowMap, nil
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
// NOCC:golint/fnsize(后续重构处理)
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
		if err := productExcelAssociationHeader(ctx, sheet, defLang, 0, asstList); err != nil {
			return err
		}
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

	if err := productExcelAssociationHeader(ctx, sheet, defLang, len(instAsst), asstList); err != nil {
		return err
	}

	rowIndex := common.HostAddMethodExcelAssociationIndexOffset

	for _, inst := range instAsst {
		cell, err := sheet.Cell(rowIndex, 1)
		if err != nil {
			return err
		}
		cell.SetString(inst.ObjectAsstID)
		cell, err = sheet.Cell(rowIndex, 2)
		if err != nil {
			return err
		}
		cell.SetString("")

		srcInst, dstInst := buildRowInfo(objID, rid, inst, asstData, objInstData)
		if srcInst == nil || dstInst == nil {
			continue
		}

		// TODO: 注意源和目标顺序
		cell, err = sheet.Cell(rowIndex, 3)
		if err != nil {
			return err
		}
		cell.SetString(buildExcelPrimaryKey(srcInst))

		cell, err = sheet.Cell(rowIndex, 4)
		if err != nil {
			return err
		}
		cell.SetString(buildExcelPrimaryKey(dstInst))

		cell, err = sheet.Cell(rowIndex, 3)
		if err != nil {
			return err
		}
		style := cell.GetStyle()
		style.Alignment.WrapText = true

		cell, err = sheet.Cell(rowIndex, 4)
		if err != nil {
			return err
		}
		style = cell.GetStyle()
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
		blog.Errorf("get object association failed, err: %v, rid: %s", err, rid)
		return err
	}
	if err := productExcelAssociationHeader(ctx, asstSheet, defLang, 0, asstList); err != nil {
		blog.Errorf("product excel association header failed, err: %v, rid: %s", err, rid)
		return err
	}

	blog.V(5).Infof("BuildExcelTemplate fields count:%d, rid: %s", fields, rid)
	if objID == common.BKInnerObjIDHost {
		cloudAreaName, _, err := lgc.getCloudArea(ctx, header)
		if err != nil {
			blog.Errorf("build %s excel template failed, err: %v,  rid: %s", objID, err, rid)
			return err
		}
		err = productHostExcelHeader(ctx, fields, filterFields, file, sheet, defLang, nil, cloudAreaName)
		if err != nil {
			blog.Errorf("product host excel header failed, err: %v, rid: %s", err, rid)
			return err
		}
	} else {
		if err := productExcelHeader(ctx, fields, filterFields, file, sheet, defLang); err != nil {
			blog.Errorf("product excel header failed, err: %v, rid: %s", err, rid)
			return err
		}
	}
	if err := ProductExcelCommentSheet(ctx, file, defLang); err != nil {
		return err
	}

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
func GetExcelData(ctx context.Context, preData *ImportExcelPreData, start, end int, defFields common.KvMap,
	defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	hosts := make(map[int]map[string]interface{})
	errMsg := make([]string, 0)

	for i := start; i < end; i++ {
		host, getErr, err := getDataFromExcel(ctx, preData, preData.DataRange[i].Start, preData.DataRange[i].End,
			defFields, defLang)
		if err != nil {
			return nil, nil, err
		}
		if len(getErr) != 0 {
			errMsg = append(errMsg, getErr...)
			continue
		}
		if len(host) != 0 {
			hosts[preData.DataRange[i].Start+1] = host
		}
	}

	if len(errMsg) != 0 {
		return nil, errMsg, nil
	}

	return hosts, nil, nil
}

// GetRawExcelData excel数据，一个kv结构，key行数（excel中的行数），value内容
func GetRawExcelData(ctx context.Context, sheet *xlsx.Sheet, defFields common.KvMap, firstRow int,
	defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {
	nameIndexMap, err := checkExcelHeader(sheet, nil, firstRow, defLang)
	if err != nil {
		return nil, nil, err
	}
	hosts := make(map[int]map[string]interface{})
	index := headerRow
	if firstRow != 0 {
		index = firstRow
	}
	errMsg := make([]string, 0)
	for ; index < sheet.MaxRow; index++ {
		row, err := sheet.Row(index)
		if err != nil {
			return nil, nil, err
		}
		host, getErr := getDataFromByExcelRow(ctx, row, index, nil, defFields, nameIndexMap, 0, row.GetCellCount(),
			defLang)
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
	map[int]metadata.ExcelAssociation, []metadata.RowMsgData, error) {

	index := firstRow

	asstInfoArr := make(map[int]metadata.ExcelAssociation)
	errMsg := make([]metadata.RowMsgData, 0)
	for ; index < sheet.MaxRow; index++ {
		row, err := sheet.Row(index)
		if err != nil {
			return nil, nil, err
		}

		// 获取单元格内容，使用for循环防止直接获取对应单元格数据导致数组越界
		var asstObjID, op, srcInst, dstInst string
		for i := 0; i < row.GetCellCount(); i++ {
			switch i {
			case associationAsstObjIDIndex:
				asstObjID = row.GetCell(i).String()
			case associationOPColIndex:
				op = row.GetCell(i).String()
			case associationSrcInstIndex:
				srcInst = row.GetCell(i).String()
			case associationDstInstIndex:
				dstInst = row.GetCell(i).String()
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

	return asstInfoArr, errMsg, nil
}

// StatisticsAssociation TODO
func StatisticsAssociation(sheet *xlsx.Sheet, firstRow int) ([]string, map[string]metadata.ObjectAsstIDStatisticsInfo,
	error) {

	index := firstRow
	asstInfoMap := make(map[string]metadata.ObjectAsstIDStatisticsInfo, 0)
	// bk_obj_asst_id
	asstNameArr := make([]string, 0)
	for ; index < sheet.MaxRow; index++ {
		var asstInfo metadata.ObjectAsstIDStatisticsInfo
		var ok bool

		row, err := sheet.Row(index)
		if err != nil {
			return nil, nil, err
		}

		count := row.GetCellCount()
		if count <= associationOPColIndex {
			continue
		}

		// bk_obj_asst_id 的值
		asstObjID := row.GetCell(associationAsstObjIDIndex).String()
		asstInfo, ok = asstInfoMap[asstObjID]
		if !ok {
			// 第一次出现这个关联关系的唯一标识
			asstNameArr = append(asstNameArr, asstObjID)
			asstInfo = metadata.ObjectAsstIDStatisticsInfo{}
		}
		asstInfo.Total += 1
		op := row.GetCell(associationOPColIndex).String()
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

	return asstNameArr, asstInfoMap, nil
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
