/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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
	"strconv"

	"configcenter/pkg/excel"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/mapstr"
	"configcenter/src/web_server/service/excel/core"
	"configcenter/src/web_server/service/excel/operator"
)

// TmplOp excel template operator
type TmplOp struct {
	*operator.BaseOp
	styleCreator *styleCreator
}

type BuildTmplOpFunc func(tmpl *TmplOp) error

// NewTmplOp create a excel template operator
func NewTmplOp(opts ...BuildTmplOpFunc) (*TmplOp, error) {
	tmpl := new(TmplOp)
	for _, opt := range opts {
		if err := opt(tmpl); err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

// BaseOperator set base operator
func BaseOperator(op *operator.BaseOp) BuildTmplOpFunc {
	return func(tmpl *TmplOp) error {
		tmpl.BaseOp = op

		var err error
		tmpl.styleCreator, err = newStyleCreator(setExcel(tmpl.GetExcel()))
		if err != nil {
			return err
		}
		return nil
	}
}

// BuildHeader create an excel with a header
func (t *TmplOp) BuildHeader(colProps ...core.ColProp) error {
	if len(colProps) == 0 {
		var err error
		colProps, err = t.GetClient().GetSortedColProp(t.GetKit(), mapstr.MapStr{common.BKObjIDField: t.GetObjID()})
		if err != nil {
			blog.Errorf("get sorted column property failed, err: %v, rid: %s", err, t.GetKit().Rid)
			return err
		}
	}

	if err := t.productSheet(); err != nil {
		blog.Errorf("product sheet failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return err
	}

	if err := t.writeInstHeader(colProps); err != nil {
		blog.Errorf("write excel instance header failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return err
	}

	if err := t.writeAssociationHeader(); err != nil {
		blog.Errorf("write excel instance association failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return err
	}

	return nil
}

func (t *TmplOp) productSheet() error {
	if err := t.GetExcel().CreateSheet(t.GetObjID()); err != nil {
		blog.Errorf("create sheet failed, objID: %s, err: %v, rid: %s", t.GetObjID(), err, t.GetKit().Rid)
		return err
	}

	if err := t.GetExcel().CreateSheet(core.AsstSheet); err != nil {
		blog.Errorf("create sheet failed, name: %s, err: %v, rid: %s", core.AsstSheet, err, t.GetKit().Rid)
		return err
	}

	return nil
}

// Close excel
func (t *TmplOp) Close() error {
	if err := t.GetExcel().Flush([]string{t.GetObjID(), core.AsstSheet}); err != nil {
		blog.Errorf("flush excel failed, sheet %s, err: %v, rid: %s", t.GetObjID(), err, t.GetKit().Rid)
		return err
	}

	if err := t.GetExcel().Save(); err != nil {
		blog.Errorf("save excel failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return err
	}

	if err := t.GetExcel().Close(); err != nil {
		blog.Errorf("close excel failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return err
	}

	return nil
}

// Clean delete temporary file
func (t *TmplOp) Clean() error {
	return t.GetExcel().Clean()
}

var (
	// InstFirstColFields excel实例sheet第0列0-5格的cell值
	InstFirstColFields = []string{common.ExcelFirstColumnFieldName, common.ExcelFirstColumnFieldType,
		common.ExcelFirstColumnFieldID, common.ExcelFirstColumnTableFieldName, common.ExcelFirstColumnTableFieldType,
		common.ExcelFirstColumnTableFieldID}

	// rowIndexes 表头中，非表格相关的字段所在行号
	rowIndexes = []int{core.NameRowIdx, core.TypeRowIdx, core.IDRowIdx}

	// tableRowIndexes 表头中，表格相关的字段所在行号
	tableRowIndexes = []int{core.TableNameRowIdx, core.TableTypeRowIdx, core.TableIDRowIdx}
)

const (
	colWidth = 24
)

func (t *TmplOp) writeInstHeader(colProps []core.ColProp) error {
	if err := t.GetExcel().SetAllColsWidth(t.GetObjID(), colWidth); err != nil {
		blog.Errorf("set sheet column width failed, objID: %s, err: %v, rid: %s", t.GetObjID(), err, t.GetKit().Rid)
		return err
	}

	header, err := t.handleProperty(colProps)
	if err != nil {
		blog.ErrorJSON("handle excel property failed, property: %s, err: %s, rid: %s", colProps, err, t.GetKit().Rid)
		return err
	}

	if err := t.GetExcel().StreamingWrite(t.GetObjID(), core.NameRowIdx, header); err != nil {
		blog.Errorf("write excel header data to excel failed, header: %v, err: %v, rid: %s", header, err,
			t.GetKit().Rid)
		return err
	}

	if err := t.mergeHeaderCell(colProps); err != nil {
		blog.Errorf("merge excel header failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return err
	}

	return nil
}

func (t *TmplOp) handleProperty(colProps []core.ColProp) ([][]excel.Cell, error) {
	ccLang := t.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(t.GetKit().Header))

	firstColStyle, err := t.styleCreator.getStyle(noEditHeader)
	if err != nil {
		blog.Errorf("get style failed, style: %s, err: %v, rid: %s", noEditHeader, err, t.GetKit().Rid)
		return nil, err
	}

	width, err := core.GetRowWidth(colProps)
	if err != nil {
		blog.Errorf("get row length failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return nil, err
	}
	header := make([][]excel.Cell, core.InstHeaderLen)
	for i := range header {
		header[i] = make([]excel.Cell, width)
	}

	for idx, field := range InstFirstColFields {
		fieldName := ccLang.Language(field)
		header[idx][0] = excel.Cell{Value: fieldName, StyleID: firstColStyle}
	}

	requiredStyle, err := t.styleCreator.getStyle(requiredField)
	if err != nil {
		blog.Errorf("get style failed, style: %s, err: %v, rid: %s", requiredField, err, t.GetKit().Rid)
		return nil, err
	}

	enumIdx := 0
	data := make([][]excel.Cell, len(colProps))
	for _, property := range colProps {
		if property.IsRequire {
			property.Name = property.Name + ccLang.Language("web_excel_header_required")
		}

		if property.PropertyType == common.FieldTypeEnumMulti || property.PropertyType == common.FieldTypeEnum {
			sheetName := "enum_" + strconv.Itoa(enumIdx)
			data[enumIdx] = []excel.Cell{{Value: sheetName + ":" + property.RefSheet}}
			property.RefSheet = sheetName
			enumIdx++
		}

		handleFunc := getHandleColPropFunc(&property)
		headerField, err := handleFunc(t, &property)
		if err != nil {
			blog.ErrorJSON("handle column property failed, sheet: %s, property: %s, err: %s, rid: %s", t.GetObjID(),
				property, err, t.GetKit().Rid)
			return nil, err
		}

		for idx, fields := range headerField {
			for fieldIdx, field := range fields {
				if property.IsRequire && idx == core.NameRowIdx {
					field.StyleID = requiredStyle
				}
				header[idx][property.ExcelColIndex+fieldIdx] = field
			}
		}
	}

	// The number of characters in an excel sheet cannot exceed 31; Since the enumeration field name may exceed this
	// length, a custom sheet name is defined and a new sheet is created to map the enumeration field name to the
	// sheet name.
	enumSheetMap := "枚举字段名映射"
	if enumIdx > 0 {
		if err := createSheetWithData(t, enumSheetMap, core.NameRowIdx, data); err != nil {
			return nil, err
		}
	}

	return header, nil
}

func (t *TmplOp) mergeHeaderCell(colProps []core.ColProp) error {
	for _, property := range colProps {
		if property.PropertyType != common.FieldTypeInnerTable {
			err := t.GetExcel().MergeSameColCell(t.GetObjID(), property.ExcelColIndex, core.TableNameRowIdx,
				core.HeaderTableLen)
			if err != nil {
				blog.Errorf("merge same column cell failed, colIdx: %d, rowIdx: %d, height: %d, err: %v, rid: %s",
					property.ExcelColIndex, core.TableNameRowIdx, core.HeaderTableLen, err, t.GetKit().Rid)
				return err
			}
			continue
		}

		for _, idx := range rowIndexes {
			if err := t.GetExcel().MergeSameRowCell(t.GetObjID(), property.ExcelColIndex, idx,
				property.Length); err != nil {
				blog.Errorf("merge same row cell failed, colIdx: %d, rowIdx: %d, length: %d, err: %v, rid: %s",
					property.ExcelColIndex, idx, property.Length, err, t.GetKit().Rid)
				return err
			}
		}
	}

	return nil
}

func (t *TmplOp) writeAssociationHeader() error {
	if err := t.setAsstColWidth(); err != nil {
		return err
	}

	if err := t.writeAsstHeader(); err != nil {
		blog.Errorf("write excel association header failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return err
	}

	if err := t.setAsstValidation(); err != nil {
		return err
	}

	return nil
}

const (
	asstFirstColWidth  = 24
	asstSecondColWidth = 36
	asstThirdColWidth  = 12
	asstFourthColWidth = 80
	asstFifthColWidth  = 80
)

func (t *TmplOp) setAsstColWidth() error {
	colWidths := []float64{
		asstFirstColWidth, asstSecondColWidth, asstThirdColWidth, asstFourthColWidth, asstFifthColWidth,
	}
	for idx, width := range colWidths {
		if err := t.GetExcel().SetColWidth(core.AsstSheet, idx+1, idx+1, width); err != nil {
			blog.Errorf("set sheet width failed, sheet: %s, err: %v, rid: %s", core.AsstSheet, err, t.GetKit().Rid)
			return err
		}
	}

	return nil
}

var firstAsstColFields = []string{
	common.ExcelFirstColumnAssociationAttribute,
	common.ExcelFirstColumnFieldDescription,
}

func (t *TmplOp) writeAsstHeader() error {
	header := make([][]excel.Cell, core.AsstExampleRowIdx+1)
	for idx := range header {
		header[idx] = make([]excel.Cell, 0)
	}
	lang := t.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(t.GetKit().Header))

	// 设置关联关系sheet表头第一列数据
	firstColStyle, err := t.styleCreator.getStyle(noEditHeader)
	if err != nil {
		blog.Errorf("get style failed, style: %s, err: %v, rid: %s", noEditHeader, err, t.GetKit().Rid)
		return err
	}
	for idx, field := range firstAsstColFields {
		header[idx] = append(header[idx], excel.Cell{StyleID: firstColStyle, Value: lang.Language(field)})
	}

	// 设置关联关系sheet表头第一行数据(除第一列的单元格)
	firstRowStyle, err := t.styleCreator.getStyle(firstRow)
	if err != nil {
		return err
	}
	firstRowFields := []string{lang.Language("excel_association_object_id"), lang.Language("excel_association_op"),
		lang.Language("excel_association_src_inst"), lang.Language("excel_association_dst_inst")}
	for idx := range firstRowFields {
		header[core.AsstStartRowIdx] = append(header[core.AsstStartRowIdx],
			excel.Cell{StyleID: firstRowStyle, Value: firstRowFields[idx]})
	}

	// 设置关联关系sheet表头第二行数据(除第一列的单元格)
	exampleStyle, err := t.styleCreator.getStyle(example)
	if err != nil {
		return err
	}
	exampleFields := []string{lang.Language("excel_example_association"), lang.Language("excel_example_op"),
		lang.Language("excel_example_association_src_inst"), lang.Language("excel_example_association_dst_inst")}
	for idx := range exampleFields {
		header[core.AsstExampleRowIdx] = append(header[core.AsstExampleRowIdx],
			excel.Cell{StyleID: exampleStyle, Value: exampleFields[idx]})
	}

	if err := t.GetExcel().StreamingWrite(core.AsstSheet, core.AsstStartRowIdx, header); err != nil {
		blog.Errorf("write excel header data to excel failed, header: %v, err: %v, rid: %s", header, err,
			t.GetKit().Rid)
		return err
	}

	return nil
}

func (t *TmplOp) setAsstValidation() error {
	params := make([]*excel.ValidationParam, 0)

	// 设置「关联标识」的列表，定义excel选项下拉
	asstList, err := t.GetClient().GetObjAssociation(t.GetKit(), t.GetObjID())
	if err != nil {
		blog.Errorf("get object association failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return err
	}
	asstNameList := make([]string, len(asstList))
	for idx, asst := range asstList {
		asstNameList[idx] = asst.AssociationName
	}

	lang := t.GetLang().CreateDefaultCCLanguageIf(httpheader.GetLanguage(t.GetKit().Header))
	refSheet := lang.Language("excel_association_object_id")
	if err := t.GetExcel().CreateSheet(refSheet); err != nil {
		return err
	}
	data := make([][]excel.Cell, len(asstList))
	for idx, asst := range asstList {
		data[idx] = append(data[idx], excel.Cell{Value: asst.AssociationName})
	}
	if err := t.GetExcel().StreamingWrite(refSheet, core.AsstStartRowIdx, data); err != nil {
		return err
	}
	if err := t.GetExcel().Flush([]string{refSheet}); err != nil {
		return err
	}
	if err := t.GetExcel().Save(); err != nil {
		return err
	}

	sqref, err := excel.GetSingleColSqref(core.AsstIDColIdx, core.AsstDataRowIdx+1, excel.GetTotalRows())
	if err != nil {
		blog.Errorf("get single column sqref failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return err
	}
	params = append(params, &excel.ValidationParam{Type: excel.Ref, Sqref: sqref, Option: refSheet})

	// 设置「操作」的列表，定义excel选项下拉
	sqref, err = excel.GetSingleColSqref(core.AsstOPColIdx, core.AsstDataRowIdx+1, excel.GetTotalRows())
	if err != nil {
		blog.Errorf("get single column sqref failed, err: %v, rid: %s", err, t.GetKit().Rid)
		return err
	}
	params = append(params, &excel.ValidationParam{Type: excel.Enum, Sqref: sqref, Option: core.AsstOps})

	for _, param := range params {
		if err = t.GetExcel().AddValidation(core.AsstSheet, param); err != nil {
			blog.Errorf("add validation failed, sheet: %s, param: %v, err: %s, rid: %s", core.AsstSheet, param, err,
				t.GetKit().Rid)

			return err
		}
	}

	return nil
}
