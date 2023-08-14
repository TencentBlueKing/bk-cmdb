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

package operator

import (
	"configcenter/pkg/excel"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/web_server/service/excel/db"
)

// Exporter operator who export excel data
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

// SetTmplOp set template operator
func SetTmplOp(tmpl *TmplOp) BuildExporterFunc {
	return func(e *Exporter) error {
		e.TmplOp = tmpl

		return nil
	}
}

// SetExportParam set export parameter
func SetExportParam(param ExportParamI) BuildExporterFunc {
	return func(e *Exporter) error {
		e.exportParam = param
		return nil
	}
}

// Export export data to excel
func (e *Exporter) Export() error {
	cond, err := e.exportParam.GetQueryPropCond()
	if err != nil {
		blog.Errorf("get property condition failed, err: %v, rid: %s", err, e.kit.Rid)
		return err
	}
	colProps, err := e.dao.GetSortedColProp(e.kit, cond)
	if err != nil {
		blog.Errorf("get sorted column property failed, err: %v, rid: %s", err, e.kit.Rid)
		return err
	}

	colProps, err = e.addExtraProp(colProps)
	if err != nil {
		blog.Errorf("add extra property failed, err: %v, rid: %s", err, e.kit.Rid)
		return err
	}

	if err := e.BuildHeader(colProps...); err != nil {
		blog.Errorf("build excel template failed, err: %v, rid: %s", err, e.kit.Rid)
		return err
	}

	rowIndex := db.InstRowIdx
	for e.exportParam.HasQueryInstCond() {
		instCond, err := e.exportParam.GetQueryInstCond()
		if err != nil {
			blog.Errorf("get instance condition failed, err: %v, rid: %s", err, e.kit.Rid)
			return err
		}

		rowIndex, err = e.exportByCond(instCond, colProps, rowIndex)
		if err != nil {
			blog.Errorf("export instance by condition failed, err: %v, rid: %s", err, e.kit.Rid)
			return err
		}
	}

	return nil
}

func (e *Exporter) addExtraProp(colProps []db.ColProp) ([]db.ColProp, error) {
	result := make([]db.ColProp, 0)
	idColIdx := common.HostAddMethodExcelDefaultIndex
	defLang := e.language.CreateDefaultCCLanguageIf(util.GetLanguage(e.kit.Header))

	if e.objID == common.BKInnerObjIDHost {
		topoProps, err := e.getTopoProps()
		if err != nil {
			blog.Errorf("get topo properties failed, err: %v, rid: %s", err, e.kit.Rid)
			return nil, err
		}

		result = append(result, topoProps...)
		idColIdx += len(topoProps)
	}

	result = append(result, db.GetIDProp(idColIdx, e.objID, defLang))

	colProps = moveOldPropIdx(len(result), colProps)

	result = append(result, colProps...)

	return result, nil
}

func moveOldPropIdx(step int, colProps []db.ColProp) []db.ColProp {
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

func (e *Exporter) getTopoProps() ([]db.ColProp, error) {
	defLang := e.language.CreateDefaultCCLanguageIf(util.GetLanguage(e.kit.Header))
	topoMsg := make([]db.TopoBriefMsg, 0)

	topoMsg = append(topoMsg, db.TopoBriefMsg{ObjID: db.TopoObjID, Name: defLang.Language(topoName)})
	topoMsg = append(topoMsg, db.TopoBriefMsg{ObjID: common.BKInnerObjIDApp, Name: defLang.Language(bizName)})

	customTopoMsg, err := e.dao.GetCustomTopoBriefMsg(e.kit)
	if err != nil {
		blog.Errorf("get custom topo name failed, err: %v, rid: %s", err, e.kit)
		return nil, err
	}
	topoMsg = append(topoMsg, customTopoMsg...)

	topoMsg = append(topoMsg, db.TopoBriefMsg{ObjID: common.BKInnerObjIDSet, Name: defLang.Language(setName)})
	topoMsg = append(topoMsg, db.TopoBriefMsg{ObjID: common.BKInnerObjIDModule, Name: defLang.Language(moduleName)})

	colIndex := common.HostAddMethodExcelDefaultIndex
	result := make([]db.ColProp, len(topoMsg))

	for idx, msg := range topoMsg {
		result[idx] = db.ColProp{ID: db.IDPrefix + msg.ObjID, Name: msg.Name, ExcelColIndex: colIndex,
			NotEditable: true}
		colIndex++
	}

	return result, nil
}

func (e *Exporter) exportByCond(cond mapstr.MapStr, colProps []db.ColProp, rowIndex int) (int, error) {
	insts, err := e.getInst(cond)
	if err != nil {
		blog.Errorf("get instance failed, objID: %s, cond: %v, err: %v, rid: %s", e.objID, cond, err, e.kit.Rid)
		return 0, err
	}

	insts, instHeights, err := e.enrichInst(insts, colProps)
	if err != nil {
		blog.Errorf("enrich instance field failed, err: %v, rid: %s", err, e.kit.Rid)
		return 0, err
	}

	for idx, inst := range insts {
		rows, err := e.handleInst(inst, colProps, instHeights[idx])
		if err != nil {
			blog.ErrorJSON("convert an instance to excel rows failed, inst: %s, property: %s, err: %s, rid: %s", inst,
				colProps, err, e.kit.Rid)
			return 0, err
		}

		if err := e.excel.StreamingWrite(e.objID, rowIndex, rows); err != nil {
			blog.ErrorJSON("write data to excel failed, rows: %s, err: %s, rid: %s", rows, err, e.kit.Rid)
			return 0, err
		}

		if err := e.mergeCell(colProps, rowIndex, instHeights[idx]); err != nil {
			blog.Errorf("merge instance cell failed, err: %v, rid: %s", err, e.kit.Rid)
			return 0, err
		}

		rowIndex += instHeights[idx]
	}

	return rowIndex, nil
}

func (e *Exporter) getInst(cond mapstr.MapStr) ([]mapstr.MapStr, error) {
	if e.objID == common.BKInnerObjIDHost {
		return e.dao.GetHost(e.kit, cond)
	}

	return e.dao.GetInst(e.kit, e.objID, cond)
}

// enrichInst 第一个返回值是返回实例数据，第二个返回值返回的每个实例数据所占用的excel行数
func (e *Exporter) enrichInst(insts []mapstr.MapStr, colProps []db.ColProp) ([]mapstr.MapStr, []int, error) {
	insts, err := e.dao.HandleEnumQuoteInst(e.kit, insts, colProps)
	if err != nil {
		blog.Errorf("handle instance enum quota field failed, err: %v, rid: %s", err, e.kit.Rid)
		return nil, nil, err
	}

	// todo 处理实例用户字段(调用GetUsernameMapWithPropertyList获取用户，但是需要解决getUsernameFromEsb的问题)

	ccLang := e.language.CreateDefaultCCLanguageIf(util.GetLanguage(e.kit.Header))
	insts, err = e.dao.GetInstWithOrgName(e.kit, ccLang, e.objID, insts, colProps)
	if err != nil {
		blog.Errorf("get instance with organization name field failed, err: %v, rid: %s", err, e.kit.Rid)
		return nil, nil, err
	}

	insts, instHeights, err := e.dao.GetInstWithTable(e.kit, e.objID, insts, colProps)
	if err != nil {
		blog.Errorf("get instance with table field failed, err: %v, rid: %s", err, e.kit.Rid)
		return nil, nil, err
	}

	return insts, instHeights, nil
}

func (e *Exporter) handleInst(inst mapstr.MapStr, colProps []db.ColProp, height int) ([][]excel.Cell, error) {
	width, err := db.GetInstWidth(colProps)
	if err != nil {
		blog.Errorf("get row length failed, err: %v, rid: %s", err, e.kit.Rid)
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

		handleFunc, isSpecial := handleSpecialInstFieldFuncMap[property.ID]
		if !isSpecial {
			var ok bool
			handleFunc, ok = handleInstFieldFuncMap[property.PropertyType]
			if !ok {
				handleFunc = getDefaultHandleFieldFunc()
			}
		}

		rows, err := handleFunc(e, &property, val)
		if err != nil {
			blog.ErrorJSON("handle instance failed, property: %s, val: %s, err: %s, rid: %s", property, val, err,
				e.kit.Rid)
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

func (e *Exporter) mergeCell(colProps []db.ColProp, rowIndex, height int) error {
	for _, property := range colProps {
		if property.PropertyType == common.FieldTypeInnerTable {
			continue
		}

		err := e.excel.MergeSameColCell(e.objID, property.ExcelColIndex, rowIndex, height)
		if err != nil {
			blog.Errorf("merge same column cell failed, colIdx: %d, rowIdx: %d, height: %d, err: %v, rid: %s",
				property.ExcelColIndex, db.TableNameRowIdx, db.HeaderTableLen, err, e.kit.Rid)
			return err
		}
	}

	return nil
}
