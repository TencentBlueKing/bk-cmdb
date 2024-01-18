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

package exporter

import (
	"configcenter/pkg/excel"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/service/excel/core"
)

// Exporter operator who export the instance to excel
type Exporter struct {
	*TmplOp
	exportParam ExportParamI
}

type BuildExporterFunc func(e *Exporter) error

// NewExporter create an operator who export excel data
func NewExporter(opts ...BuildExporterFunc) (*Exporter, error) {
	e := new(Exporter)
	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}

// TmplOperator set template operator
func TmplOperator(tmpl *TmplOp) BuildExporterFunc {
	return func(e *Exporter) error {
		e.TmplOp = tmpl
		return nil
	}
}

// ExportParam set export parameter
func ExportParam(param ExportParamI) BuildExporterFunc {
	return func(e *Exporter) error {
		e.exportParam = param
		return nil
	}
}

// Export export data to excel
func (e *Exporter) Export() error {
	cond, err := e.exportParam.GetPropCond()
	if err != nil {
		blog.Errorf("get property condition failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return err
	}
	colProps, err := e.GetClient().GetSortedColProp(e.GetKit(), cond)
	if err != nil {
		blog.Errorf("get sorted column property failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return err
	}

	colProps, err = e.addExtraProp(colProps)
	if err != nil {
		blog.Errorf("add extra property failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return err
	}

	if err := e.BuildHeader(colProps...); err != nil {
		blog.Errorf("build excel template failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return err
	}

	instIDs, err := e.exportInst(colProps)
	if err != nil {
		blog.Errorf("export instance failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return err
	}

	if err := e.exportAsst(instIDs); err != nil {
		blog.Errorf("export association failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return err
	}

	return nil
}

func (e *Exporter) addExtraProp(colProps []core.ColProp) ([]core.ColProp, error) {
	result := make([]core.ColProp, 0)
	idColIdx := common.HostAddMethodExcelDefaultIndex
	defLang := e.GetLang().CreateDefaultCCLanguageIf(util.GetLanguage(e.GetKit().Header))

	if e.GetObjID() == common.BKInnerObjIDHost {
		topoProps, err := e.getTopoProps()
		if err != nil {
			blog.Errorf("get topo properties failed, err: %v, rid: %s", err, e.GetKit().Rid)
			return nil, err
		}

		result = append(result, topoProps...)
		idColIdx += len(topoProps)
	}

	result = append(result, core.GetIDProp(idColIdx, e.GetObjID(), defLang))

	colProps = moveOldPropIdx(len(result), colProps)

	result = append(result, colProps...)

	return result, nil
}

func moveOldPropIdx(step int, colProps []core.ColProp) []core.ColProp {
	for idx := range colProps {
		colProps[idx].ExcelColIndex = colProps[idx].ExcelColIndex + step
	}

	return colProps
}

const (
	topoName   = "web_ext_field_topo"
	bizName    = "biz_property_bk_biz_name"
	moduleName = "web_ext_field_module_name"
	setName    = "web_ext_field_set_name"
)

func (e *Exporter) getTopoProps() ([]core.ColProp, error) {
	defLang := e.GetLang().CreateDefaultCCLanguageIf(util.GetLanguage(e.GetKit().Header))
	topoMsg := make([]core.TopoBriefMsg, 0)

	topoMsg = append(topoMsg, core.TopoBriefMsg{ObjID: core.TopoObjID, Name: defLang.Language(topoName)})
	topoMsg = append(topoMsg, core.TopoBriefMsg{ObjID: common.BKInnerObjIDApp, Name: defLang.Language(bizName)})

	customTopoMsg, err := e.GetClient().GetCustomTopoBriefMsg(e.GetKit())
	if err != nil {
		blog.Errorf("get custom topo name failed, err: %v, rid: %s", err, e.GetKit())
		return nil, err
	}
	topoMsg = append(topoMsg, customTopoMsg...)

	topoMsg = append(topoMsg, core.TopoBriefMsg{ObjID: common.BKInnerObjIDSet, Name: defLang.Language(setName)})
	topoMsg = append(topoMsg, core.TopoBriefMsg{ObjID: common.BKInnerObjIDModule, Name: defLang.Language(moduleName)})

	colIndex := common.HostAddMethodExcelDefaultIndex
	result := make([]core.ColProp, len(topoMsg))

	for idx, msg := range topoMsg {
		result[idx] = core.ColProp{ID: core.IDPrefix + msg.ObjID, Name: msg.Name, ExcelColIndex: colIndex,
			NotEditable: true}
		colIndex++
	}

	return result, nil
}

func (e *Exporter) exportInst(colProps []core.ColProp) ([]int64, error) {
	rowIndex := core.InstRowIdx
	var result, instIDs []int64
	for e.exportParam.HasInstCond() {
		instCond, err := e.exportParam.GetInstCond()
		if err != nil {
			blog.Errorf("get instance condition failed, err: %v, rid: %s", err, e.GetKit().Rid)
			return nil, err
		}

		rowIndex, instIDs, err = e.exportByCond(instCond, colProps, rowIndex)
		if err != nil {
			blog.Errorf("export instance by condition failed, err: %v, rid: %s", err, e.GetKit().Rid)
			return nil, err
		}

		result = append(result, instIDs...)
	}

	return result, nil
}

func (e *Exporter) exportByCond(cond mapstr.MapStr, colProps []core.ColProp, rowIndex int) (int, []int64, error) {
	insts, err := e.getInst(cond)
	if err != nil {
		blog.Errorf("get instance failed, objID: %s, cond: %v, err: %v, rid: %s", e.GetObjID(), cond, err,
			e.GetKit().Rid)
		return 0, nil, err
	}

	if len(insts) == 0 {
		return rowIndex, nil, nil
	}

	insts, instHeights, err := e.enrichInst(insts, colProps)
	if err != nil {
		blog.Errorf("enrich instance field failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return 0, nil, err
	}

	instIDs := make([]int64, 0)
	instIDKey := metadata.GetInstIDFieldByObjID(e.GetObjID())
	for idx, inst := range insts {
		rows, err := e.handleInst(inst, colProps, instHeights[idx])
		if err != nil {
			blog.ErrorJSON("convert an instance to excel rows failed, inst: %s, property: %s, err: %s, rid: %s", inst,
				colProps, err, e.GetKit().Rid)
			return 0, nil, err
		}

		if err := e.GetExcel().StreamingWrite(e.GetObjID(), rowIndex, rows); err != nil {
			blog.ErrorJSON("write data to excel failed, rows: %s, err: %s, rid: %s", rows, err, e.GetKit().Rid)
			return 0, nil, err
		}

		if err := e.mergeCell(colProps, rowIndex, instHeights[idx]); err != nil {
			blog.Errorf("merge instance cell failed, err: %v, rid: %s", err, e.GetKit().Rid)
			return 0, nil, err
		}

		rowIndex += instHeights[idx]

		instID, err := inst.Int64(instIDKey)
		if err != nil {
			blog.Errorf("parse instance(%+v) id(key:%s) failed, err: %v, objID: %s, rid: %s", inst, instIDKey, err,
				e.GetObjID(), e.GetKit().Rid)
		}
		instIDs = append(instIDs, instID)
	}

	return rowIndex, instIDs, nil
}

func (e *Exporter) getInst(cond mapstr.MapStr) ([]mapstr.MapStr, error) {
	if e.GetObjID() == common.BKInnerObjIDHost {
		return e.GetClient().GetHost(e.GetKit(), cond)
	}

	return e.GetClient().GetInst(e.GetKit(), e.GetObjID(), cond)
}

// enrichInst 第一个返回值是返回实例数据，第二个返回值返回的每个实例数据所占用的excel行数
func (e *Exporter) enrichInst(insts []mapstr.MapStr, colProps []core.ColProp) ([]mapstr.MapStr, []int, error) {
	insts, err := e.GetClient().TransEnumQuoteIDToName(e.GetKit(), insts, colProps)
	if err != nil {
		blog.Errorf("handle instance enum quota field failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return nil, nil, err
	}

	ccLang := e.GetLang().CreateDefaultCCLanguageIf(util.GetLanguage(e.GetKit().Header))
	insts, err = e.GetClient().GetInstWithOrgName(e.GetKit(), ccLang, insts, colProps)
	if err != nil {
		blog.Errorf("get instance with organization name field failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return nil, nil, err
	}

	insts, err = e.GetClient().GetInstWithUserFullName(e.GetKit(), ccLang, e.GetObjID(), insts)
	if err != nil {
		blog.Errorf("get instance with full user name field failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return nil, nil, err
	}

	insts, instHeights, err := e.GetClient().GetInstWithTable(e.GetKit(), e.GetObjID(), insts, colProps)
	if err != nil {
		blog.Errorf("get instance with table field failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return nil, nil, err
	}

	return insts, instHeights, nil
}

func (e *Exporter) handleInst(inst mapstr.MapStr, colProps []core.ColProp, height int) ([][]excel.Cell, error) {
	width, err := core.GetRowWidth(colProps)
	if err != nil {
		blog.Errorf("get row length failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return nil, err
	}

	instRows := make([][]excel.Cell, height)
	for i := range instRows {
		instRows[i] = make([]excel.Cell, width)
	}

	for _, property := range colProps {
		if property.NotExport {
			continue
		}

		val, ok := inst[property.ID]
		if !ok {
			continue
		}

		handleFunc := getHandleInstFieldFunc(&property)
		rows, err := handleFunc(e, &property, val)
		if err != nil {
			blog.ErrorJSON("handle instance failed, property: %s, val: %s, err: %s, rid: %s", property, val, err,
				e.GetKit().Rid)
			return nil, err
		}

		for idx, row := range rows {
			for cellIdx, cell := range row {
				instRows[idx][property.ExcelColIndex+cellIdx] = cell
			}
		}

		// 如果单元格为空，并且它的上一行单元格有样式，那么需要把上一行对应的样式赋值过来，否则进行单元格合并时，会缺少样式
		for rowIdx := 1; rowIdx < height; rowIdx++ {
			for colIdx := 0; colIdx < width; colIdx++ {
				if instRows[rowIdx][colIdx].Value == nil && instRows[rowIdx-1][colIdx].StyleID != 0 {
					instRows[rowIdx][colIdx].StyleID = instRows[rowIdx-1][colIdx].StyleID
				}
			}
		}
	}

	return instRows, nil
}

func (e *Exporter) mergeCell(colProps []core.ColProp, rowIndex, height int) error {
	for _, property := range colProps {
		if property.PropertyType == common.FieldTypeInnerTable {
			continue
		}

		err := e.GetExcel().MergeSameColCell(e.GetObjID(), property.ExcelColIndex, rowIndex, height)
		if err != nil {
			blog.Errorf("merge same column cell failed, colIdx: %d, rowIdx: %d, height: %d, err: %v, rid: %s",
				property.ExcelColIndex, core.TableNameRowIdx, core.HeaderTableLen, err, e.GetKit().Rid)
			return err
		}
	}

	return nil
}

func (e *Exporter) exportAsst(instIDs []int64) error {
	// 未设置, 不导出关联关系数据
	if len(e.exportParam.GetAsstObjUniqueIDMap()) == 0 {
		return nil
	}

	asstData, err := e.getInstAsst(instIDs)
	if err != nil {
		blog.Errorf("get instance association failed, instIDs: %v, err: %v, rid: %s", instIDs, err, e.GetKit().Rid)
		return err
	}

	if err := e.GetExcel().StreamingWrite(core.AsstSheet, core.AsstDataRowIdx, asstData); err != nil {
		blog.Errorf("write instance association to excel failed, data: %v, err: %v, rid: %s", asstData, err,
			e.GetKit().Rid)
		return err
	}

	return nil
}

func (e *Exporter) getInstAsst(instIDs []int64) ([][]excel.Cell, error) {
	// 1. 获取需要导出的模型关联关系，以及判断是否有自关联的关联关系
	asstObjUniqueIDMap := e.exportParam.GetAsstObjUniqueIDMap()
	asstObjIDMap := make(map[string]struct{})
	hasSelfAsst := false

	for key := range asstObjUniqueIDMap {
		if key == e.GetObjID() {
			hasSelfAsst = true
			continue
		}

		asstObjIDMap[key] = struct{}{}
	}

	asstList, err := e.GetClient().GetObjAssociation(e.GetKit(), e.GetObjID())
	if err != nil {
		blog.Errorf("get object association failed, err: %v, rid: %s", err, e.GetKit().Rid)
		return nil, err
	}

	asstIDs := make([]string, 0)
	for _, asst := range asstList {
		if hasSelfAsst && asst.ObjectID == asst.AsstObjID {
			asstIDs = append(asstIDs, asst.AssociationName)
			continue
		}

		_, ok := asstObjIDMap[asst.ObjectID]
		if ok {
			asstIDs = append(asstIDs, asst.AssociationName)
			continue
		}

		_, ok = asstObjIDMap[asst.AsstObjID]
		if ok {
			asstIDs = append(asstIDs, asst.AssociationName)
		}
	}

	// 2. 获取实例的关联关系数据
	instAsstArr, err := e.GetClient().GetInstAsst(e.GetKit(), e.GetObjID(), instIDs, asstIDs, hasSelfAsst)
	if err != nil {
		blog.Errorf("get instance association failed, instIDs: %v, asstIDs: %v, err: %v, rid: %d", instIDs, asstIDs,
			err, e.GetKit().Rid)
		return nil, err
	}

	// 3. 获取实例关联关系中，源实例和目标实例的唯一标识信息; 这里当前模型的唯一标识会单独获取
	asstData, err := e.getInstAsstData(instAsstArr)
	if err != nil {
		blog.Errorf("get instance association data failed, instAsstArr: %v, err: %v, rid: %s", instAsstArr, err,
			e.GetKit().Rid)
		return nil, err
	}

	// 4. 构造需要写到excel的关联关系数据
	result := make([][]excel.Cell, len(asstData))
	for idx, data := range asstData {
		row := make([]excel.Cell, core.AsstDstInstColIdx+1)
		row[core.AsstIDColIdx] = excel.Cell{Value: data.asstID}
		row[core.AsstSrcInstColIdx] = excel.Cell{Value: data.srcInst}
		row[core.AsstDstInstColIdx] = excel.Cell{Value: data.destInst}

		result[idx] = row
	}

	return result, nil
}

type instAsstData struct {
	asstID   string
	srcInst  string
	destInst string
}

func (e *Exporter) getInstAsstData(instAsstArr []*metadata.InstAsst) ([]instAsstData, error) {
	// 当前操作对象实例id
	instIDs := make([]int64, 0)
	// 关联的对方的实例id
	asstInstIDMap := make(map[string][]int64)

	for _, instAsst := range instAsstArr {
		if instAsst.ObjectID == e.GetObjID() {
			instIDs = append(instIDs, instAsst.InstID)
			asstInstIDMap[instAsst.AsstObjectID] = append(asstInstIDMap[instAsst.AsstObjectID], instAsst.AsstInstID)
			continue
		}

		instIDs = append(instIDs, instAsst.AsstInstID)
		asstInstIDMap[instAsst.ObjectID] = append(asstInstIDMap[instAsst.ObjectID], instAsst.InstID)
	}

	asstObjUniqueIDMap := e.exportParam.GetAsstObjUniqueIDMap()

	asstInstUniqueKeyMap := make(map[string]map[int64]string)
	for objID, asstInstIDs := range asstInstIDMap {
		instUniqueKeys, err := e.GetClient().GetInstUniqueKeys(e.GetKit(), objID, asstInstIDs,
			asstObjUniqueIDMap[objID])
		if err != nil {
			blog.Errorf("get instance uniques keys failed, objID: %s, err: %v, rid: %s", objID, err, e.GetKit().Rid)
			return nil, err
		}

		asstInstUniqueKeyMap[objID] = instUniqueKeys
	}

	curInstUniqueKey, err := e.GetClient().GetInstUniqueKeys(e.GetKit(), e.GetObjID(), instIDs,
		e.exportParam.GetObjUniqueID())
	if err != nil {
		blog.Errorf("get instance uniques keys failed, objID: %s, err: %v, rid: %s", e.GetObjID(), err, e.GetKit().Rid)
		return nil, err
	}

	result := make([]instAsstData, 0)
	for _, instAsst := range instAsstArr {
		if instAsst.ObjectID == e.GetObjID() {
			srcInst, ok := curInstUniqueKey[instAsst.InstID]
			if !ok {
				blog.Warnf("association is invalid, val: %v, objID: %s, rid: %s", instAsst, e.GetObjID(),
					e.GetKit().Rid)
				continue
			}

			asstInstUniqueKey, ok := asstInstUniqueKeyMap[instAsst.AsstObjectID]
			if !ok {
				blog.Warnf("association is invalid, val: %v, objID: %s, rid: %s", instAsst, instAsst.AsstObjectID,
					e.GetKit().Rid)
				continue
			}

			dstInst, ok := asstInstUniqueKey[instAsst.AsstInstID]
			if !ok {
				blog.Warnf("association is invalid, val: %v, objID: %s, rid: %s", instAsst, instAsst.AsstObjectID,
					e.GetKit().Rid)
				continue
			}

			result = append(result, instAsstData{asstID: instAsst.ObjectAsstID, srcInst: srcInst, destInst: dstInst})
			continue
		}

		asstInstUniqueKey, ok := asstInstUniqueKeyMap[instAsst.ObjectID]
		if !ok {
			blog.Warnf("association is invalid, val: %v, objID: %s, rid: %s", instAsst, instAsst.ObjectID,
				e.GetKit().Rid)
			continue
		}

		srcInst, ok := asstInstUniqueKey[instAsst.InstID]
		if !ok {
			blog.Warnf("association is invalid, val: %v, objID: %s, rid: %s", instAsst, instAsst.ObjectID,
				e.GetKit().Rid)
			continue
		}

		dstInst, ok := curInstUniqueKey[instAsst.AsstInstID]
		if !ok {
			blog.Warnf("association is invalid, val: %v, objID: %s, rid: %s", instAsst, e.GetObjID(), e.GetKit().Rid)
			continue
		}

		result = append(result, instAsstData{asstID: instAsst.ObjectAsstID, srcInst: srcInst, destInst: dstInst})
	}

	return result, nil
}
