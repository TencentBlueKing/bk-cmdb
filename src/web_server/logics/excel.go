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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"
)

// BuildExcelFromData product excel from data
func (lgc *Logics) BuildExcelFromData(ctx context.Context, objID string, fields map[string]Property, filter []string, data []mapstr.MapStr, xlsxFile *xlsx.File, header http.Header, modelBizID int64, usernameMap map[string]string, propertyList []string) error {
	rid := util.GetHTTPCCRequestID(header)

	ccLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))
	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	sheet, err := xlsxFile.AddSheet("inst")
	if err != nil {
		blog.Errorf("setExcelRowDataByIndex add excel sheet error, err:%s, rid:%s", err.Error(), rid)
		return err

	}
	addSystemField(fields, common.BKInnerObjIDObject, ccLang, 1)

	if 0 == len(filter) {
		filter = getFilterFields(objID)
	} else {
		filter = append(filter, getFilterFields(objID)...)
	}

	instPrimaryKeyValMap := make(map[int64][]PropertyPrimaryVal)
	productExcelHeader(ctx, fields, filter, sheet, ccLang)
	// indexID := getFieldsIDIndexMap(fields)

	rowIndex := common.HostAddMethodExcelIndexOffset

	for _, rowMap := range data {

		instIDKey := metadata.GetInstIDFieldByObjID(objID)
		instID, err := rowMap.Int64(instIDKey)
		if err != nil {
			blog.Errorf("setExcelRowDataByIndex inst:%+v, not inst id key:%s, objID:%s, rid:%s", rowMap, instIDKey, objID, rid)
			return ccErr.Errorf(common.CCErrCommInstFieldNotFound, "instIDKey", objID)
		}
		// 使用中英文用户名重新构造用户列表(用户列表实际为逗号分隔的string型)
		rowMap, err = replaceEnName(rid, rowMap, usernameMap, propertyList, ccLang)
		if err != nil {
			blog.Errorf("rebuild user list field, rid: %s", rid)
			return err
		}

		primaryKeyArr := setExcelRowDataByIndex(rowMap, sheet, rowIndex, fields)

		instPrimaryKeyValMap[instID] = primaryKeyArr
		rowIndex++

	}

	err = lgc.BuildAssociationExcelFromData(ctx, objID, instPrimaryKeyValMap, xlsxFile, header, modelBizID)
	if err != nil {
		return err
	}
	return nil
}

// BuildHostExcelFromData product excel from data
func (lgc *Logics) BuildHostExcelFromData(ctx context.Context, objID string, fields map[string]Property,
	filter []string, data []mapstr.MapStr, xlsxFile *xlsx.File, header http.Header, modelBizID int64,
	usernameMap map[string]string, propertyList []string, customLen int, objName []string) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	ccLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))
	ccErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	sheet, err := xlsxFile.AddSheet("host")
	if err != nil {
		blog.Errorf("BuildHostExcelFromData add excel sheet error, err:%s, rid:%s", err.Error(), rid)
		return err
	}

	extFieldKey := make([]string, 0)
	extFieldsTopoID := "cc_ext_field_topo"
	extFieldsBizID := "cc_ext_biz"
	extFieldsModuleID := "cc_ext_module"
	extFieldsSetID := "cc_ext_set"
	extFieldKey = append(extFieldKey, extFieldsTopoID, extFieldsBizID)
	extFields := map[string]string{
		extFieldsTopoID:   ccLang.Language("web_ext_field_topo"),
		extFieldsBizID:    ccLang.Language("biz_property_bk_biz_name"),
		extFieldsModuleID: ccLang.Language("bk_module_name"),
		extFieldsSetID:    ccLang.Language("bk_set_name"),
	}
	extFieldsCustomID1 := "cc_ext_custom1"
	extFieldsCustomID2 := "cc_ext_custom2"
	extFieldsCustomID3 := "cc_ext_custom3"
	switch customLen {
	case 1:
		extFields[extFieldsCustomID1] = objName[0]
		extFieldKey = append(extFieldKey, extFieldsCustomID1, extFieldsSetID, extFieldsModuleID)
	case 2:
		extFields[extFieldsCustomID1] = objName[1]
		extFields[extFieldsCustomID2] = objName[0]
		extFieldKey = append(extFieldKey, extFieldsCustomID2, extFieldsCustomID1, extFieldsSetID, extFieldsModuleID)
	case 3:
		extFields[extFieldsCustomID1] = objName[2]
		extFields[extFieldsCustomID2] = objName[1]
		extFields[extFieldsCustomID3] = objName[0]
		extFieldKey = append(extFieldKey, extFieldsCustomID3, extFieldsCustomID2, extFieldsCustomID1, extFieldsSetID,
			extFieldsModuleID)
	default:
		extFieldKey = append(extFieldKey, extFieldsSetID, extFieldsModuleID)
	}

	fields = addExtFields(fields, extFields, extFieldKey)
	addSystemField(fields, common.BKInnerObjIDHost, ccLang, customLen+5)

	productHostExcelHeader(ctx, fields, filter, sheet, ccLang, customLen, objName)

	instPrimaryKeyValMap := make(map[int64][]PropertyPrimaryVal)
	// indexID := getFieldsIDIndexMap(fields)
	rowIndex := common.HostAddMethodExcelIndexOffset
	for _, hostData := range data {

		rowMap, err := mapstr.NewFromInterface(hostData[common.BKInnerObjIDHost])
		if err != nil {
			blog.ErrorJSON("BuildHostExcelFromData failed, hostData: %s, err: %s, rid: %s", hostData, err.Error(), rid)
			return ccErr.CCError(common.CCErrCommReplyDataFormatError)
		}

		if _, exist := fields[common.BKCloudIDField]; exist {
			cloudAreaArr, err := rowMap.MapStrArray(common.BKCloudIDField)
			if err != nil {
				blog.ErrorJSON("build host excel failed, cloud area not array, host: %s, err: %s, rid: %s", hostData, err, rid)
				return ccErr.CCError(common.CCErrCommReplyDataFormatError)
			}

			if len(cloudAreaArr) != 1 {
				blog.ErrorJSON("build host excel failed, host has many cloud areas, host: %s, err: %s, rid: %s", hostData, err, rid)
				return ccErr.CCError(common.CCErrCommReplyDataFormatError)
			}

			cloudArea := fmt.Sprintf("%v[%v]", cloudAreaArr[0][common.BKInstNameField], cloudAreaArr[0][common.BKInstIDField])
			rowMap.Set(common.BKCloudIDField, cloudArea)
		}

		moduleMap, ok := hostData[common.BKInnerObjIDModule].([]interface{})
		if ok {
			topos := util.GetStrValsFromArrMapInterfaceByKey(moduleMap, "TopModuleName")
			if len(topos) > 0 {
				idx := strings.Index(topos[0], logics.SplitFlag)
				if idx > 0 {
					rowMap[extFieldsBizID] = topos[0][:idx]
				}

				toposNobiz := make([]string, 0)
				for _, topo := range topos {
					idx := strings.Index(topo, logics.SplitFlag)
					if idx > 0 && len(topo) >= idx+len(logics.SplitFlag) {
						toposNobiz = append(toposNobiz, topo[idx+len(logics.SplitFlag):])
					}
				}
				rowMap[extFieldsTopoID] = strings.Join(toposNobiz, ", ")
			}
		}

		result, err := lgc.getTopoMainlineInstRoot(ctx, header, modelBizID, moduleMap)
		if err != nil {
			blog.Errorf("get topo mainline instance root failed, err: %s, rid: %s", err.Error(), rid)
		}

		var moduleStr, setStr, customStr1, customStr2, customStr3 string
		for _, res := range result.Nodes {
			length := len(res.Path) - 3
			switch length {
			case 1:
				if customStr1 == "" {
					customStr1 = res.Path[2].InstanceName
				} else {
					ok := util.Contains(strings.Split(customStr1, ","), res.Path[2].InstanceName)
					if !ok {
						customStr1 += "," + res.Path[2].InstanceName
					}
				}
			case 2:
				if customStr1 == "" {
					customStr1 = res.Path[2].InstanceName
				} else {
					ok := util.Contains(strings.Split(customStr1, ","), res.Path[2].InstanceName)
					if !ok {
						customStr1 += "," + res.Path[2].InstanceName
					}
				}
				if customStr2 == "" {
					customStr2 = res.Path[3].InstanceName
				} else {
					ok := util.Contains(strings.Split(customStr2, ","), res.Path[3].InstanceName)
					if !ok {
						customStr2 += "," + res.Path[3].InstanceName
					}
				}
			case 3:
				if customStr1 == "" {
					customStr1 = res.Path[2].InstanceName
				} else {
					ok := util.Contains(strings.Split(customStr1, ","), res.Path[2].InstanceName)
					if !ok {
						customStr1 += "," + res.Path[2].InstanceName
					}
				}

				if customStr2 == "" {
					customStr2 = res.Path[3].InstanceName
				} else {
					ok := util.Contains(strings.Split(customStr2, ","), res.Path[3].InstanceName)
					if !ok {
						customStr2 += "," + res.Path[3].InstanceName
					}
				}

				if customStr3 == "" {
					customStr3 = res.Path[4].InstanceName
				} else {
					ok := util.Contains(strings.Split(customStr3, ","), res.Path[4].InstanceName)
					if !ok {
						customStr3 += "," + res.Path[4].InstanceName
					}
				}
			}

			if moduleStr == "" {
				moduleStr = res.Path[0].InstanceName
			} else {
				ok := util.Contains(strings.Split(moduleStr, ","), res.Path[0].InstanceName)
				if !ok {
					moduleStr += "," + res.Path[0].InstanceName
				}
			}

			if setStr == "" {
				setStr = res.Path[1].InstanceName
			} else {
				ok := util.Contains(strings.Split(setStr, ","), res.Path[1].InstanceName)
				if !ok {
					setStr += "," + res.Path[1].InstanceName
				}
			}
		}

		rowMap[extFieldsModuleID] = moduleStr
		rowMap[extFieldsSetID] = setStr
		switch customLen {
		case 1:
			rowMap[extFieldsCustomID1] = customStr1
		case 2:
			rowMap[extFieldsCustomID1] = customStr1
			rowMap[extFieldsCustomID2] = customStr2
		case 3:
			rowMap[extFieldsCustomID1] = customStr1
			rowMap[extFieldsCustomID2] = customStr2
			rowMap[extFieldsCustomID3] = customStr3
		}

		instIDKey := metadata.GetInstIDFieldByObjID(objID)
		instID, err := rowMap.Int64(instIDKey)
		if err != nil {
			blog.Errorf("setExcelRowDataByIndex inst:%+v, not inst id key:%s, objID:%s, rid:%s", rowMap, instIDKey, objID, rid)
			return ccErr.Errorf(common.CCErrCommInstFieldNotFound, instIDKey, objID)
		}

		// 使用中英文用户名重新构造用户列表(用户列表实际为逗号分隔的string型)
		rowMap, err = replaceEnName(rid, rowMap, usernameMap, propertyList, ccLang)
		if err != nil {
			blog.Errorf("rebuild user list field, rid: %s", rid)
			return err
		}

		primaryKeyArr := setExcelRowDataByIndex(rowMap, sheet, rowIndex, fields)
		instPrimaryKeyValMap[instID] = primaryKeyArr
		rowIndex++
	}

	err = lgc.BuildAssociationExcelFromData(ctx, objID, instPrimaryKeyValMap, xlsxFile, header, modelBizID)
	if err != nil {
		return err
	}
	return nil
}

// getTopoMainlineInstRoot get topo mainline inst root
func (lgc *Logics) getTopoMainlineInstRoot(ctx context.Context, header http.Header, modelBizID int64,
	moduleMap []interface{}) (*metadata.TopoPathResult, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	nodes := make([]metadata.TopoNode, 0)
	for _, row := range moduleMap {
		mapRow, ok := row.(map[string]interface{})
		if ok {
			moduleID, err := util.GetIntByInterface(mapRow[common.BKModuleIDField])
			if err != nil {
				return nil, err
			}
			node := metadata.TopoNode{
				ObjectID:   common.BKInnerObjIDModule,
				InstanceID: int64(moduleID),
			}
			nodes = append(nodes, node)
		}
	}
	input := metadata.FindTopoPathRequest{
		Nodes: nodes,
	}

	topoRoot, err := lgc.Engine.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(ctx, header,
		modelBizID, false)
	if err != nil {
		blog.Errorf("search mainline instance topo path failed, bizID:%d, err:%s, rid:%s", modelBizID,
			err.Error(), rid)
		return nil, err
	}
	result := &metadata.TopoPathResult{}
	for _, node := range input.Nodes {
		topoPath := topoRoot.TraversalFindNode(node.ObjectID, node.InstanceID)
		path := make([]*metadata.TopoInstanceNodeSimplify, 0)
		for _, item := range topoPath {
			simplify := item.ToSimplify()
			path = append(path, simplify)
		}
		nodeTopoPath := metadata.NodeTopoPath{
			BizID: modelBizID,
			Node:  node,
			Path:  path,
		}
		result.Nodes = append(result.Nodes, nodeTopoPath)
	}
	return result, err
}

// GetCustomCntAndInstName get custom level count and instance name
func (lgc *Logics) GetCustomCntAndInstName(ctx context.Context, header http.Header) (int, []string, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	mainlineAsstRsp, err := lgc.CoreAPI.CoreService().Association().ReadModelAssociation(ctx, header,
		&metadata.QueryCondition{Condition: map[string]interface{}{common.AssociationKindIDField: common.
			AssociationKindMainline}})
	if nil != err {
		blog.Errorf("search mainline association failed, error: %s, rid: %s", err.Error(), rid)
		return 0, []string{}, err
	}

	mainlineObjectChildMap := make(map[string]string, 0)
	objectName := make([]string, 0)
	customLen := 0
	isMainline := false
	for _, asst := range mainlineAsstRsp.Data.Info {
		if asst.ObjectID == common.BKInnerObjIDHost {
			continue
		}
		mainlineObjectChildMap[asst.AsstObjID] = asst.ObjectID
		if asst.AsstObjID == common.BKInnerObjIDApp {
			isMainline = true
		}
	}
	if !isMainline {
		return customLen, objectName, nil
	}

	// get all mainline object name map
	objectIDs := make([]string, 0)
	for objectID := common.BKInnerObjIDApp; len(objectID) != 0; objectID = mainlineObjectChildMap[objectID] {
		objectIDs = append(objectIDs, objectID)
	}
	cond := make([]string, 0)
	for _, obj := range objectIDs {
		if obj == common.BKInnerObjIDApp || obj == common.BKInnerObjIDSet || obj == common.BKInnerObjIDModule {
			continue
		}
		cond = append(cond, obj)
	}

	input := &metadata.QueryCondition{
		Fields: []string{common.BKObjNameField, common.BKAsstObjIDField},
		Condition: mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{common.BKDBIN: cond},
		},
	}

	objects, err := lgc.CoreAPI.CoreService().Model().ReadModel(ctx, header, input)
	if nil != err {
		blog.ErrorJSON("search mainline objects(%s) failed, error: %s, rid: %s", objectIDs, err.Error(), rid)
		return customLen, objectName, err
	}

	for _, objID := range objectIDs {
		for _, val := range objects.Data.Info {
			if val.Spec.ObjectID == objID {
				objectName = append(objectName, val.Spec.ObjectName)
			}
		}
	}

	customLen = int(mainlineAsstRsp.Data.Count) - 3

	return customLen, objectName, nil
}

func (lgc *Logics) BuildAssociationExcelFromData(ctx context.Context, objID string, instPrimaryInfo map[int64][]PropertyPrimaryVal, xlsxFile *xlsx.File, header http.Header, modelBizID int64) error {
	defLang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))
	rid := util.ExtractRequestIDFromContext(ctx)
	var instIDArr []int64
	for instID := range instPrimaryInfo {
		instIDArr = append(instIDArr, instID)
	}
	instAsst, err := lgc.fetchAssociationData(ctx, header, objID, instIDArr, modelBizID)
	if err != nil {
		return err
	}
	asstData, err := lgc.getAssociationData(ctx, header, objID, instAsst, modelBizID)
	if err != nil {
		return err
	}

	sheet, err := xlsxFile.AddSheet("assocation")
	if err != nil {
		blog.Errorf("setExcelRowDataByIndex add excel  assocation sheet error. err:%s, rid:%s", err.Error(), rid)
		return err
	}

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
	//确定关联标识的列表，定义excel选项下拉栏。此处需要查cc_ObjAsst表。
	resp, err := lgc.CoreAPI.TopoServer().Association().SearchObject(ctx, header, cond)
	if err != nil {
		blog.ErrorJSON("get object association list failed, err: %v, rid: %s", err, rid)
		return err
	}
	if err := resp.CCError(); err != nil {
		blog.ErrorJSON("get object association list failed, err: %v, rid: %s", resp.ErrMsg, rid)
		return err
	}
	asstList := resp.Data
	productExcelAssociationHeader(ctx, sheet, defLang, len(instAsst), asstList)

	rowIndex := common.HostAddMethodExcelAssociationIndexOffset

	for _, inst := range instAsst {
		sheet.Cell(rowIndex, 1).SetString(inst.ObjectAsstID)
		sheet.Cell(rowIndex, 2).SetString("")
		srcInst, ok := asstData[inst.ObjectID][inst.InstID]
		if !ok {
			blog.Warnf("BuildAssociationExcelFromData association inst:%+v, not inst id :%d, objID:%s, rid:%s", inst, inst.InstID, objID, rid)
			// return lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).Errorf(common.CCErrCommInstDataNil, fmt.Sprintf("%s %d", objID, inst.InstID))
			continue
		}
		dstInst, ok := asstData[inst.AsstObjectID][inst.AsstInstID]
		if !ok {
			blog.Warnf("BuildAssociationExcelFromData association inst:%+v, not inst id :%d, objID:%s, rid:%s", inst, inst.InstID, inst.AsstObjectID, rid)
			continue
		}
		sheet.Cell(rowIndex, 3).SetString(buildEexcelPrimaryKey(srcInst))
		sheet.Cell(rowIndex, 4).SetString(buildEexcelPrimaryKey(dstInst))
		style := sheet.Cell(rowIndex, 3).GetStyle()
		style.Alignment.WrapText = true
		style = sheet.Cell(rowIndex, 4).GetStyle()
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

// BuildExcelTemplate  return httpcode, error
func (lgc *Logics) BuildExcelTemplate(ctx context.Context, objID, filename string, header http.Header, defLang lang.DefaultCCLanguageIf, modelBizID int64) error {
	rid := util.GetHTTPCCRequestID(header)
	filterFields := getFilterFields(objID)
	// host excel template doesn't need export field bk_cloud_id
	if objID == common.BKInnerObjIDHost {
		filterFields = append(filterFields, common.BKCloudIDField)
	}
	fields, err := lgc.GetObjFieldIDs(objID, filterFields, nil, header, modelBizID, common.HostAddMethodExcelDefaultIndex)
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
	blog.V(5).Infof("BuildExcelTemplate fields count:%d, rid: %s", fields, rid)
	productExcelHeader(ctx, fields, filterFields, sheet, defLang)
	ProductExcelCommentSheet(ctx, file, defLang)

	if err = file.Save(filename); nil != err {
		blog.Errorf("save file failed, filename: %s, err: %+v, rid: %s", filename, err, rid)
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
	c.Header("Content-Disposition", "attachment; filename="+name) // 文件名
	c.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

// GetExcelData excel数据，一个kv结构，key行数（excel中的行数），value内容
func GetExcelData(ctx context.Context, sheet *xlsx.Sheet, fields map[string]Property, defFields common.KvMap, isCheckHeader bool, firstRow int, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	var err error
	nameIndexMap, err := checkExcelHeader(ctx, sheet, fields, isCheckHeader, defLang)
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
		host, getErr := getDataFromByExcelRow(ctx, row, index, fields, defFields, nameIndexMap, defLang)
		if 0 != len(getErr) {
			errMsg = append(errMsg, getErr...)
			continue
		}
		if 0 != len(host) {
			hosts[index+1] = host
		}
	}
	if 0 != len(errMsg) {
		return nil, errMsg, nil
	}

	return hosts, nil, nil

}

// GetExcelData excel数据，一个kv结构，key行数（excel中的行数），value内容
func GetRawExcelData(ctx context.Context, sheet *xlsx.Sheet, defFields common.KvMap, firstRow int, defLang lang.DefaultCCLanguageIf) (map[int]map[string]interface{}, []string, error) {

	var err error
	nameIndexMap, err := checkExcelHeader(ctx, sheet, nil, false, defLang)
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
		host, getErr := getDataFromByExcelRow(ctx, row, index, nil, defFields, nameIndexMap, defLang)
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

func GetAssociationExcelData(sheet *xlsx.Sheet, firstRow int) map[int]metadata.ExcelAssociation {

	rowCnt := len(sheet.Rows)
	index := firstRow

	asstInfoArr := make(map[int]metadata.ExcelAssociation, 0)
	for ; index < rowCnt; index++ {
		row := sheet.Rows[index]
		op := row.Cells[associationOPColIndex].String()
		if op == "" {
			continue
		}

		asstObjID := row.Cells[assciationAsstObjIDIndex].String()
		srcInst := row.Cells[assciationSrcInstIndex].String()
		dstInst := row.Cells[assciationDstInstIndex].String()
		asstInfoArr[index] = metadata.ExcelAssociation{
			ObjectAsstID: asstObjID,
			Operate:      getAssociationExcelOperateFlag(op),
			SrcPrimary:   srcInst,
			DstPrimary:   dstInst,
		}
	}

	return asstInfoArr
}

// GetFilterFields 不需要展示字段
func GetFilterFields(objID string) []string {
	return getFilterFields(objID)
}

// GetCustomFields 用户展示字段export时优先排序
func GetCustomFields(filterFields []string, customFieldsStr string) []string {
	return getCustomFields(filterFields, customFieldsStr)
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
