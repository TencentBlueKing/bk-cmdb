/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package importer

import (
	"fmt"

	"configcenter/pkg/excel"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/service/excel/core"
	"configcenter/src/web_server/service/excel/operator"
)

// Importer operator who import the instance into excel
type Importer struct {
	*operator.BaseOp
	param ImportParamI
}

type BuildImporterFunc func(importer *Importer) error

// NewImporter create a excel instance importer
func NewImporter(opts ...BuildImporterFunc) (*Importer, error) {
	importer := new(Importer)
	for _, opt := range opts {
		if err := opt(importer); err != nil {
			return nil, err
		}
	}

	return importer, nil
}

// BaseOperator set base operator
func BaseOperator(op *operator.BaseOp) BuildImporterFunc {
	return func(importer *Importer) error {
		importer.BaseOp = op
		return nil
	}
}

// Param set importer parameter
func Param(param ImportParamI) BuildImporterFunc {
	return func(importer *Importer) error {
		importer.param = param
		return nil
	}
}

// Clean close importer file io and remove excel file
func (i *Importer) Clean() error {
	if err := i.GetExcel().Close(); err != nil {
		blog.Errorf("close excel failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return err
	}

	if err := i.GetExcel().Clean(); err != nil {
		blog.Errorf("remove excel file failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return err
	}

	return nil
}

const (
	onceImportLimit = 100
)

// Handle handle import request
func (i *Importer) Handle() (mapstr.MapStr, error) {
	// 获取association sheet中关联的模型以及关联关系的条数等信息
	if i.param.GetOpType() == getAsstFlag {
		result, err := i.getAsstInfo()
		if err != nil {
			blog.Errorf("get instance association info failed, err: %v, rid: %s", err, i.GetKit().Rid)
			return nil, err
		}

		return result, nil
	}

	// 从excel获取实例数据，进行导入
	result, hasErrMsg, err := i.importInst()
	if err != nil {
		blog.Errorf("import instance failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	if hasErrMsg || len(i.param.GetAsstObjUniqueIDMap()) == 0 {
		return result, nil
	}

	// 从excel获取关联关系数据，进行导入
	return i.importAssociation()
}

func (i *Importer) getAsstInfo() (mapstr.MapStr, error) {
	exist, err := i.GetExcel().IsSheetExist(core.AsstSheet)
	if err != nil {
		return nil, err
	}

	if !exist {
		return mapstr.MapStr{"association": mapstr.New()}, nil
	}

	asstInfo, err := i.getAsstFromExcel()
	if err != nil {
		blog.Errorf("get association info from excel failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	if asstInfo == nil || len(asstInfo.asstIDs) == 0 {
		return mapstr.MapStr{"association": mapstr.New()}, nil
	}

	associations, err := i.GetClient().FindAsstByAsstID(i.GetKit(), i.GetObjID(), asstInfo.asstIDs)
	if err != nil {
		blog.Errorf("find model association by bk_obj_asst_id failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	statisticalInfos := make(map[string]metadata.ObjectAsstIDStatisticsInfo, 0)
	for _, asst := range associations {
		objID := asst.AsstObjID
		// 只统计关联的对象
		if asst.ObjectID != i.GetObjID() {
			objID = asst.ObjectID
		}

		statisticalInfo, ok := statisticalInfos[objID]
		if !ok {
			statisticalInfo = metadata.ObjectAsstIDStatisticsInfo{}
		}

		data := asstInfo.statisticalMap[asst.AssociationName]

		statisticalInfo.Create += data.Create
		statisticalInfo.Delete += data.Delete
		statisticalInfo.Total += data.Total

		statisticalInfos[objID] = statisticalInfo
	}

	return mapstr.MapStr{"association": statisticalInfos}, nil
}

type excelAsstInfo struct {
	// asstIDs 关联标识id数组
	asstIDs []string
	// statisticalMap key为关联标识id，value为关联标识使用的统计结果结构体
	statisticalMap map[string]metadata.ObjectAsstIDStatisticsInfo
	// asstInfoMap 需要导入的excel关联关系数据
	asstInfoMap map[int]metadata.ExcelAssociation
	// errMsg 关联关系sheet不合法的数据信息
	errMsg []metadata.RowMsgData
}

func (i *Importer) getAsstFromExcel() (*excelAsstInfo, error) {
	reader, err := i.GetExcel().NewReader(core.AsstSheet)
	if err != nil {
		blog.Errorf("create excel reader failed, sheet: %s, err: %v, rid: %s", core.AsstSheet, err, i.GetKit().Rid)
		return nil, err
	}
	lang := i.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(i.GetKit().Header))

	statisticalMap := make(map[string]metadata.ObjectAsstIDStatisticsInfo, 0)
	asstIDs := make([]string, 0)
	asstInfoMap := make(map[int]metadata.ExcelAssociation)
	errMsg := make([]metadata.RowMsgData, 0)

	for reader.Next() {
		if reader.GetCurIdx() < core.AsstDataRowIdx {
			continue
		}

		row, err := reader.CurRow()
		if err != nil {
			blog.Errorf("get reader current row data failed, err: %v, rid: %s", err, i.GetKit().Rid)
			return nil, err
		}

		if len(row) < core.AsstDstInstColIdx+1 {
			msg := lang.Languagef("web_excel_row_handle_error", core.AsstSheet, reader.GetCurIdx()+1)
			errMsg = append(errMsg, metadata.RowMsgData{Row: reader.GetCurIdx(), Msg: msg})
			continue
		}

		asstID := row[core.AsstIDColIdx]
		op := row[core.AsstOPColIdx]
		srcInst := row[core.AsstSrcInstColIdx]
		dstInst := row[core.AsstDstInstColIdx]

		idx := reader.GetCurIdx() + 1

		if asstID == "" || op == "" || srcInst == "" || dstInst == "" {
			msg := lang.Languagef("web_excel_row_handle_error", core.AsstSheet, idx)
			errMsg = append(errMsg, metadata.RowMsgData{Row: idx, Msg: msg})
			continue
		}

		statisticalInfo, ok := statisticalMap[asstID]
		if !ok {
			asstIDs = append(asstIDs, asstID)
			statisticalInfo = metadata.ObjectAsstIDStatisticsInfo{}
		}

		operate := core.GetAsstOpFlag(core.AsstOp(op))
		switch operate {
		case metadata.ExcelAssociationOperateDelete:
			statisticalInfo.Delete += 1
		case metadata.ExcelAssociationOperateAdd:
			statisticalInfo.Create += 1
		}

		statisticalMap[asstID] = statisticalInfo

		asstInfoMap[idx] = metadata.ExcelAssociation{
			ObjectAsstID: asstID,
			Operate:      operate,
			SrcPrimary:   srcInst,
			DstPrimary:   dstInst,
		}
	}

	if err := reader.Close(); err != nil {
		blog.Errorf("close reader failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	result := &excelAsstInfo{asstIDs: asstIDs, statisticalMap: statisticalMap, asstInfoMap: asstInfoMap, errMsg: errMsg}
	return result, nil
}

func (i *Importer) preCheck(excelMsg *ExcelMsg) ([]string, error) {
	reader, err := i.GetExcel().NewReader(i.GetObjID())
	if err != nil {
		blog.Errorf("create excel io reader failed, sheet: %s, err: %v, rid: %s", i.GetObjID(), err, i.GetKit().Rid)
		return nil, err
	}
	var instCount int
	for reader.Next() {
		instCount++
	}

	// 实例数需要减去excel表头占用的行数
	instCount -= core.InstHeaderLen

	// 如果存在合并多行作为一个实例，那么需要将这些多出来的行数减掉
	if excelMsg != nil && len(excelMsg.mergeRowRes) != 0 {
		for start, end := range excelMsg.mergeRowRes {
			instCount -= end - start
		}
	}

	lang := i.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(i.GetKit().Header))
	if instCount > common.ExcelImportMaxRow {
		return []string{lang.Languagef("web_excel_import_too_much", common.ExcelImportMaxRow)}, nil
	}

	exist, err := i.isAsstExist()
	if err != nil {
		blog.Errorf("check if there is associated data failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	if !exist && instCount <= 0 {
		return []string{lang.Language("web_excel_not_data")}, nil
	}

	return nil, nil
}

func (i *Importer) isAsstExist() (bool, error) {
	exist, err := i.GetExcel().IsSheetExist(core.AsstSheet)
	if err != nil {
		return false, err
	}

	if !exist {
		return false, nil
	}

	reader, err := i.GetExcel().NewReader(core.AsstSheet)
	if err != nil {
		blog.Errorf("create excel reader failed, sheet: %s, err: %v, rid: %s", core.AsstSheet, err, i.GetKit().Rid)
		return false, err
	}

	exist = false
	for reader.Next() {
		if reader.GetCurIdx() < core.AsstDataRowIdx {
			continue
		}

		return true, nil
	}

	return exist, nil
}

func (i *Importer) importInst() (mapstr.MapStr, bool, error) {
	reader, err := i.GetExcel().NewReader(i.GetObjID())
	if err != nil {
		blog.Errorf("create excel io reader failed, sheet: %s, err: %v, rid: %s", i.GetObjID(), err, i.GetKit().Rid)
		return nil, false, err
	}
	excelMsg, err := i.getExcelMsg(reader)
	if err != nil {
		blog.Errorf("get object excel message failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, false, err
	}

	result := mapstr.New()
	var errMsg []string
	errMsg, err = i.preCheck(excelMsg)
	if err != nil {
		return nil, false, err
	}
	if len(errMsg) != 0 {
		result["error"] = errMsg
		return result, true, nil
	}

	lang := i.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(i.GetKit().Header))
	var successMsg []int64
	insts := make(map[int]map[string]interface{})

	hasDoneNext := false
	for hasDoneNext || reader.Next() {
		hasDoneNext = false
		// skip excel header
		if reader.GetCurIdx() < core.InstRowIdx-1 {
			continue
		}

		idx := reader.GetCurIdx() + 1
		inst, err := i.getNextInst(reader, excelMsg)
		if err != nil {
			blog.Errorf("get next instance from excel failed, err: %v, rid: %s", err, i.GetKit().Rid)
			errMsg = append(errMsg, lang.Languagef("import_data_fail", idx, err.Error()))
			continue
		}
		if inst != nil {
			insts[idx] = inst
		}
		if len(insts) < onceImportLimit && reader.Next() {
			hasDoneNext = true
			continue
		}

		var errRes []string
		insts, errRes = i.doSpecialOp(insts)
		errMsg = append(errMsg, errRes...)
		if len(insts) == 0 {
			continue
		}

		req, err := i.param.BuildParam(insts)
		if err != nil {
			blog.Errorf("get import instances parameter failed, err: %v, rid: %s", err, i.GetKit().Rid)
			return nil, false, err
		}
		importParam := &core.ImportedParam{Language: i.GetLang(), ObjID: i.GetObjID(), Instances: insts,
			Req: req, HandleType: i.param.GetHandleType()}
		successRes, errRes := i.GetClient().HandleImportedInst(i.GetKit(), importParam)
		if len(successRes) != 0 {
			successMsg = append(successMsg, successRes...)
		}
		errMsg = append(errMsg, errRes...)

		insts = make(map[int]map[string]interface{})
	}

	result["success"] = successMsg
	result["error"] = errMsg
	if err := reader.Close(); err != nil {
		blog.Errorf("close read excel io failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, false, err
	}

	return result, len(errMsg) > 0, nil
}

func (i *Importer) getExcelMsg(reader *excel.Reader) (*ExcelMsg, error) {
	propertyMap, err := i.getPropertyMap(reader)
	if err != nil {
		blog.Errorf("get property failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	mergeRes, err := i.getMergeRowRes()
	if err != nil {
		blog.Errorf("get excel merge row resource, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	return &ExcelMsg{propertyMap: propertyMap, mergeRowRes: mergeRes}, nil
}

func (i *Importer) getPropertyMap(reader *excel.Reader) (map[int]PropWithTable, error) {

	cond := mapstr.MapStr{
		common.BKObjIDField: i.GetObjID(),
		common.BKAppIDField: i.param.GetBizID(),
	}
	colProps, err := i.GetClient().GetObjColProp(i.GetKit(), cond)
	if err != nil {
		blog.Errorf("get property failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}
	handleType := i.param.GetHandleType()
	if handleType == core.UpdateHost || handleType == core.AddInst {
		lang := i.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(i.GetKit().Header))
		colProps = append(colProps, core.GetIDProp(core.PropDefaultColIdx, i.GetObjID(), lang))
	}

	propMap := make(map[string]core.ColProp)
	for _, prop := range colProps {
		propMap[prop.ID] = prop
	}

	for reader.GetCurIdx() != core.IDRowIdx && reader.Next() {
		continue
	}
	idRow, err := reader.CurRow()
	if err != nil {
		blog.Errorf("read data from excel failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	for reader.GetCurIdx() != core.TableIDRowIdx && reader.Next() {
		continue
	}
	tableIDRow, err := reader.CurRow()
	if err != nil {
		blog.Errorf("read data from excel failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	return i.buildPropWithTable(propMap, idRow, tableIDRow)
}

func (i *Importer) buildPropWithTable(propMap map[string]core.ColProp, idRow []string, tableIDRow []string) (
	map[int]PropWithTable, error) {

	result := make(map[int]PropWithTable)
	for idx, propID := range idRow {
		prop, ok := propMap[propID]
		if !ok {
			blog.Warnf("can not find property, id: %s, kit: %s", propID, i.GetKit().Rid)
			continue
		}
		prop.ExcelColIndex = idx

		if prop.PropertyType != common.FieldTypeInnerTable {
			result[idx] = PropWithTable{ColProp: prop}
			continue
		}

		option, err := metadata.ParseTableAttrOption(prop.Option)
		if err != nil {
			return nil, err
		}
		tableHeaderMap := make(map[string]core.ColProp, len(option.Header))
		for _, attr := range option.Header {
			subProp := core.ColProp{ID: attr.PropertyID, Name: attr.PropertyName, PropertyType: attr.PropertyType,
				Option: attr.Option, IsRequire: attr.IsRequired, Length: core.PropertyNormalLen}

			tableHeaderMap[attr.PropertyID] = subProp
		}

		subProperties := make(map[int]PropWithTable, len(option.Header))
		for subIdx := idx; subIdx < idx+len(option.Header) && subIdx < len(tableIDRow); subIdx++ {
			tablePropID := tableIDRow[subIdx]

			subProp, ok := tableHeaderMap[tablePropID]
			if !ok {
				blog.Errorf("can not find table sub property, id: %s, table sub id: %s, kit: %s", propID, tablePropID,
					i.GetKit().Rid)
				continue
			}

			subProp.ExcelColIndex = subIdx
			subProperties[subIdx] = PropWithTable{ColProp: subProp}
		}

		prop.Length = len(subProperties)
		result[idx] = PropWithTable{ColProp: prop, subProperties: subProperties}
	}

	return result, nil
}

// getMergeRowRes 获取对同一列的行进行合并的开始和结束范围
func (i *Importer) getMergeRowRes() (map[int]int, error) {
	msgs, err := i.GetExcel().GetMergeCellMsg(i.GetObjID())
	if err != nil {
		return nil, err
	}

	result := make(map[int]int)

	for _, msg := range msgs {
		startCol, startRow, err := excel.CellNameToCoordinates(msg.GetStartAxis())
		if err != nil {
			return nil, err
		}

		endCol, endRow, err := excel.CellNameToCoordinates(msg.GetEndAxis())
		if err != nil {
			return nil, err
		}

		// 目前只会存在同一行合并或者同一列合并的情况，如果有其他情况，那就是不合法的excel
		if startCol != endCol && startRow != endRow {
			return nil, fmt.Errorf("excel is invalid, start: %s, end: %s", msg.GetStartAxis(), msg.GetEndAxis())
		}

		// 这个是对于同一行的列的合并，跳过
		if startCol != endCol {
			continue
		}

		startRowVal := startRow - 1
		endRowVal := endRow - 1

		// excel表头中的表格字段合并不需要记录占用的行数
		if startRowVal == core.TableNameRowIdx {
			continue
		}

		val, ok := result[startRowVal]
		if !ok {
			result[startRowVal] = endRowVal
			continue
		}

		// 如果已经存在同一列的合并结果，但是范围不一样，那就是不合法的excel
		if val != endRowVal {
			return nil, fmt.Errorf("excel is invalid, row start: %d, end: %d, another end: %d", startRow, endRow, val+1)
		}
	}

	return result, nil
}

func (i *Importer) getNextInst(reader *excel.Reader, excelMsg *ExcelMsg) (map[string]interface{}, error) {
	rows := make([][]string, 0)
	row, err := reader.CurRow()
	if err != nil {
		blog.Errorf("read data from excel failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}
	rows = append(rows, row)

	endRow, ok := excelMsg.mergeRowRes[reader.GetCurIdx()]
	if ok {
		for endRow != reader.GetCurIdx() && reader.Next() {
			row, err := reader.CurRow()
			if err != nil {
				blog.Errorf("read data from excel failed, err: %v, rid: %s", err, i.GetKit().Rid)
				return nil, err
			}
			rows = append(rows, row)
		}
	}

	inst := make(map[string]interface{})
	hasInst := false
	for idx, val := range rows[0] {
		if val == "" {
			continue
		}
		hasInst = true

		prop, ok := excelMsg.propertyMap[idx]
		if !ok {
			continue
		}

		handleFunc := getHandleInstFieldFunc(&prop)

		value, err := handleFunc(i, &prop, rows)
		if err != nil {
			blog.ErrorJSON("handle instance failed, property: %s, data: %s, err: %s, rid: %s", prop, rows, err,
				i.GetKit().Rid)
		}
		inst[prop.ID] = value
	}

	if !hasInst {
		return nil, nil
	}

	return inst, nil
}

func (i *Importer) doSpecialOp(insts map[int]map[string]interface{}) (map[int]map[string]interface{}, []string) {
	if i.GetObjID() != common.BKInnerObjIDHost || len(insts) == 0 {
		return insts, nil
	}

	var errMsg []string
	hosts, errRes := i.transferCloudArea(insts)
	errMsg = append(errMsg, errRes...)

	if len(hosts) == 0 {
		return nil, errMsg
	}

	handleType := i.param.GetHandleType()

	switch handleType {
	case core.AddHost:
		hosts, errRes = i.checkAddedHost(hosts)
		errMsg = append(errMsg, errRes...)

	case core.UpdateHost:
		hosts, errRes = i.checkUpdatedHost(hosts)
		errMsg = append(errMsg, errRes...)
	}

	return hosts, errMsg
}

func (i *Importer) transferCloudArea(hosts map[int]map[string]interface{}) (map[int]map[string]interface{}, []string) {
	if i.GetObjID() != common.BKInnerObjIDHost {
		return hosts, nil
	}

	var errMsg []string
	legalHost := make(map[int]map[string]interface{})
	lang := i.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(i.GetKit().Header))

	cloudNames := []string{common.DefaultCloudName}
	for _, host := range hosts {
		if name, ok := host[common.BKCloudIDField]; ok {
			cloudNames = append(cloudNames, util.GetStrByInterface(name))
		}
	}
	_, cloudMap, err := i.GetClient().GetCloudArea(i.GetKit(), util.StrArrayUnique(cloudNames)...)
	if err != nil {
		blog.Errorf("get host cloud area failed, err: %v, rid: %s", err, i.GetKit().Rid)
		for idx := range hosts {
			errMsg = append(errMsg, lang.Languagef("import_data_fail", idx, err.Error()))
		}
		return nil, errMsg
	}

	for _, index := range util.SortedMapIntKeys(hosts) {
		host := hosts[index]
		if host == nil {
			continue
		}

		cloudStr := common.DefaultCloudName
		if _, ok := host[common.BKCloudIDField]; ok {
			cloudStr = util.GetStrByInterface(hosts[index][common.BKCloudIDField])
		}

		if _, ok := cloudMap[cloudStr]; !ok {
			blog.Errorf("cloud area name %s of line %d doesn't exist, rid: %s", cloudStr, index, i.GetKit().Rid)
			msg := lang.Languagef("import_host_cloudID_not_exist", index, host[common.BKHostInnerIPField], cloudStr)
			errMsg = append(errMsg, msg)
			continue
		}

		host[common.BKCloudIDField] = cloudMap[cloudStr]
		legalHost[index] = host
	}

	return legalHost, errMsg
}

func (i *Importer) checkAddedHost(hosts map[int]map[string]interface{}) (map[int]map[string]interface{}, []string) {
	if i.GetObjID() != common.BKInnerObjIDHost {
		return hosts, nil
	}

	var errMsg []string
	legalHost := make(map[int]map[string]interface{})
	lang := i.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(i.GetKit().Header))

	res, err := i.GetClient().GetSameIPRes(i.GetKit(), hosts)
	if err != nil {
		blog.Errorf("get host same ip resource failed, err: %v, rid: %s", err, i.GetKit().Rid)
		for idx := range hosts {
			errMsg = append(errMsg, lang.Languagef("import_data_fail", idx, err.Error()))
		}
		return nil, errMsg
	}

	for _, index := range util.SortedMapIntKeys(hosts) {
		host := hosts[index]
		if host == nil {
			continue
		}

		if _, ok := host[common.BKHostIDField]; ok {
			errMsg = append(errMsg, lang.Languagef("import_host_no_need_hostID", index))
			continue
		}

		if _, ok := host[common.BKAgentIDField]; ok {
			errMsg = append(errMsg, lang.Languagef("import_host_no_need_agentID", index))
			continue
		}

		cloud, ok := host[common.BKCloudIDField]
		if !ok {
			errMsg = append(errMsg, lang.Languagef("import_host_not_provide_cloudID", index))
			continue
		}

		addressType, ok := host[common.BKAddressingField].(string)
		if !ok {
			addressType = common.BKAddressingStatic
		}

		if addressType != common.BKAddressingStatic && addressType != common.BKAddressingDynamic {
			errMsg = append(errMsg, lang.Languagef("import_host_illegal_addressing", index))
			continue
		}

		// in dynamic scenarios, there is no need to do duplication check of ip address.
		if addressType == common.BKAddressingDynamic {
			legalHost[index] = host
			continue
		}

		innerIPv4, v4Ok := host[common.BKHostInnerIPField].(string)
		innerIPv6, v6Ok := host[common.BKHostInnerIPv6Field].(string)
		if (!v4Ok || innerIPv4 == "") && (!v6Ok || innerIPv6 == "") {
			errMsg = append(errMsg, lang.Languagef("host_import_innerip_v4_v6_empty", index))
			continue
		}

		// check if the host ipv4 exist in db
		key := core.HostCloudKey(innerIPv4, cloud)
		if _, exist := res.V4Map[key]; exist {
			errMsg = append(errMsg, lang.Languagef("host_import_innerip_v4_fail", index))
			continue
		}

		// check if the host ipv6 exist in db
		keyV6 := core.HostCloudKey(innerIPv6, cloud)
		if _, exist := res.V6Map[keyV6]; exist {
			errMsg = append(errMsg, lang.Languagef("host_import_innerip_v6_fail", index))
			continue
		}

		legalHost[index] = host
	}

	return legalHost, errMsg
}

func (i *Importer) checkUpdatedHost(hosts map[int]map[string]interface{}) (map[int]map[string]interface{}, []string) {
	if i.GetObjID() != common.BKInnerObjIDHost {
		return hosts, nil
	}

	var errMsg []string
	lang := i.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(i.GetKit().Header))
	res, err := i.getCheckHostRes(hosts)
	if err != nil {
		blog.Errorf("get check updated host resource failed, err: %v, rid: %s", err, i.GetKit().Rid)
		for idx := range hosts {
			errMsg = append(errMsg, lang.Languagef("import_update_data_fail", idx, err.Error()))
		}
		return nil, errMsg
	}

	legalHost := make(map[int]map[string]interface{})
	for _, index := range util.SortedMapIntKeys(hosts) {
		host := hosts[index]
		if host == nil {
			continue
		}

		hostID, ok := host[common.BKHostIDField]
		if !ok {
			blog.Errorf("bk_host_id field doesn't exist, innerIpv4: %v, innerIpv6: %v, rid: %v",
				host[common.BKHostInnerIPField], host[common.BKHostInnerIPv6Field], i.GetKit().Rid)
			errMsg = append(errMsg, lang.Languagef("import_update_host_miss_hostID", index))
			continue
		}
		hostIDVal, err := util.GetInt64ByInterface(hostID)
		if err != nil {
			errMsg = append(errMsg, lang.Languagef("import_update_host_hostID_not_int", index))
			continue
		}

		// check if the host exist in db
		ip := res.existingHosts[hostIDVal].Ip
		ipv6 := res.existingHosts[hostIDVal].Ipv6
		agentID := res.existingHosts[hostIDVal].AgentID
		if ip == "" && ipv6 == "" && agentID == "" {
			errMsg = append(errMsg, lang.Languagef("import_host_no_exist_error", index, hostIDVal))
			continue
		}

		// check if the host innerIP and hostID is consistent
		excelIP := util.GetStrByInterface(host[common.BKHostInnerIPField])
		if ip != excelIP {
			errMsg = append(errMsg, lang.Languagef("import_host_ip_not_consistent", index, excelIP,
				hostIDVal, ip))
			continue
		}

		// check if the host innerIPv6 and hostID is consistent
		excelIPv6 := util.GetStrByInterface(host[common.BKHostInnerIPv6Field])
		if ipv6 != excelIPv6 {
			errMsg = append(errMsg, lang.Languagef("import_host_ipv6_not_consistent", index, excelIPv6,
				hostIDVal, ipv6))
			continue
		}

		// check if the host agentID and hostID is consistent
		excelAgentID := util.GetStrByInterface(host[common.BKAgentIDField])
		if agentID != excelAgentID {
			errMsg = append(errMsg, lang.Languagef("import_host_agentID_not_consistent", index, excelAgentID,
				hostIDVal, agentID))
			continue
		}

		// check if the hostID and bizID is consistent
		if res.hostBizMap[hostIDVal] != res.bizID {
			errMsg = append(errMsg, lang.Languagef("import_hostID_bizID_not_consistent", index, excelIP, excelIPv6))
			continue
		}

		legalHost[index] = host
	}

	return legalHost, errMsg
}

type checkHostRes struct {
	existingHosts map[int64]core.SimpleHost
	bizID         int64
	hostBizMap    map[int64]int64
}

func (i *Importer) getCheckHostRes(hosts map[int]map[string]interface{}) (*checkHostRes, error) {
	existingHosts, err := i.GetClient().GetExistingHost(i.GetKit(), hosts)
	if err != nil {
		blog.Errorf("get existing hosts failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	bizID := i.param.GetBizID()
	if bizID == 0 {
		// get resource pool biz ID
		bizID, err = i.GetClient().GetDefaultBizID(i.GetKit())
		if err != nil {
			blog.Errorf("get resource pool biz ID failed, err: %v, rid:%s", err, i.GetKit().Rid)
			return nil, err
		}
	}

	hostBizMap, err := i.GetClient().GetHostBizRelations(i.GetKit(), hosts)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	return &checkHostRes{existingHosts: existingHosts, bizID: bizID, hostBizMap: hostBizMap}, nil
}

func (i *Importer) importAssociation() (mapstr.MapStr, error) {
	exist, err := i.GetExcel().IsSheetExist(core.AsstSheet)
	if err != nil {
		return nil, err
	}

	result := mapstr.New()
	if !exist {
		return result, nil
	}

	asstInfo, err := i.getAsstFromExcel()
	if err != nil {
		blog.Errorf("get association info from excel failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}
	asstObjUniqueIDMap := i.param.GetAsstObjUniqueIDMap()

	if asstInfo == nil || asstObjUniqueIDMap == nil {
		return result, nil
	}

	// 将不需要导入的关联关系数据过滤出来
	associations, err := i.GetClient().FindAsstByAsstID(i.GetKit(), i.GetObjID(), asstInfo.asstIDs)
	if err != nil {
		blog.Errorf("find model association by bk_obj_asst_id failed, err: %v, rid: %s", err, i.GetKit().Rid)
		return nil, err
	}

	skipAsstID := make(map[string]struct{})
	for _, asst := range associations {
		_, hasAsstObjID := asstObjUniqueIDMap[asst.AsstObjID]
		_, hasObjID := asstObjUniqueIDMap[asst.ObjectID]

		// 如果有一个为true, 表示该类型的关联关系是需要导入的
		if hasAsstObjID || hasObjID {
			continue
		}

		skipAsstID[asst.AssociationName] = struct{}{}
	}

	importedAsst := make(map[int]metadata.ExcelAssociation)
	for idx, asst := range asstInfo.asstInfoMap {
		if _, ok := skipAsstID[asst.ObjectAsstID]; ok {
			continue
		}

		importedAsst[idx] = asst
	}
	if len(importedAsst) == 0 {
		return result, nil
	}

	// 导入指定的关联关系
	input := &metadata.RequestImportAssociation{
		AssociationInfoMap:    importedAsst,
		AsstObjectUniqueIDMap: asstObjUniqueIDMap,
		ObjectUniqueID:        i.param.GetObjUniqueID(),
	}

	asstResp, err := i.GetClient().ImportAssociation(i.GetKit(), i.GetObjID(), input)
	if err != nil {
		blog.Errorf("import association failed, input: %v, err: %v, rid: %s", input, err, i.GetKit().Rid)
		return nil, err
	}

	if len(asstResp.ErrMsgMap) != 0 {
		result["error"] = asstResp.ErrMsgMap
	}

	return result, nil
}
