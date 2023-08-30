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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/web_server/service/excel/core"
)

// TmplOp excel template operator
type TmplOp struct {
	excel        *excel.Excel
	styleCreator *styleCreator
	client       *core.Client
	objID        string
	kit          *rest.Kit
	language     language.CCLanguageIf
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

// FilePath set template operator file path
func FilePath(filePath string) BuildTmplOpFunc {
	return func(tmpl *TmplOp) error {
		var err error
		tmpl.excel, err = excel.NewExcel(excel.FilePath(filePath), excel.OpenOrCreate())
		if err != nil {
			return err
		}

		tmpl.styleCreator, err = newStyleCreator(setExcel(tmpl.excel))
		if err != nil {
			return err
		}

		return nil
	}
}

// Client set template client
func Client(client *core.Client) BuildTmplOpFunc {
	return func(tmpl *TmplOp) error {
		tmpl.client = client
		return nil
	}
}

// ObjID set template operator object id
func ObjID(objID string) BuildTmplOpFunc {
	return func(tmpl *TmplOp) error {
		tmpl.objID = objID
		return nil
	}
}

// Kit set template operator kit
func Kit(kit *rest.Kit) BuildTmplOpFunc {
	return func(tmpl *TmplOp) error {
		tmpl.kit = kit
		return nil
	}
}

// Language set template operator language
func Language(language language.CCLanguageIf) BuildTmplOpFunc {
	return func(tmpl *TmplOp) error {
		tmpl.language = language
		return nil
	}
}

// BuildHeader create an excel with a header
func (t *TmplOp) BuildHeader(colProps ...core.ColProp) error {
	if len(colProps) == 0 {
		var err error
		colProps, err = t.client.GetSortedColProp(t.kit, mapstr.MapStr{common.BKObjIDField: t.objID})
		if err != nil {
			blog.Errorf("get sorted column property failed, err: %v, rid: %s", err, t.kit.Rid)
			return err
		}
	}

	if err := t.productInstHeader(colProps); err != nil {
		blog.Errorf("product excel instance header failed, err: %v, rid: %s", err, t.kit.Rid)
		return err
	}

	if err := t.productAssociationHeader(); err != nil {
		blog.Errorf("product excel instance association failed, err: %v, rid: %s", err, t.kit.Rid)
		return err
	}

	return nil
}

// Close excel
func (t *TmplOp) Close() error {
	if err := t.excel.Flush([]string{t.objID, core.AsstSheet}); err != nil {
		blog.Errorf("flush excel failed, sheet %s, err: %v, rid: %s", t.objID, err, t.kit.Rid)
		return err
	}

	if err := t.excel.Save(); err != nil {
		blog.Errorf("save excel failed, err: %v, rid: %s", err, t.kit.Rid)
		return err
	}

	if err := t.excel.Close(); err != nil {
		blog.Errorf("close excel failed, err: %v, rid: %s", err, t.kit.Rid)
		return err
	}

	return nil
}

// Clean delete temporary file
func (t *TmplOp) Clean() error {
	return t.excel.Clean()
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

func (t *TmplOp) productInstHeader(colProps []core.ColProp) error {
	if err := t.excel.CreateSheet(t.objID); err != nil {
		blog.Errorf("create sheet failed, objID: %s, err: %v, rid: %s", t.objID, err, t.kit.Rid)
		return err
	}

	if err := t.excel.SetAllColsWidth(t.objID, colWidth); err != nil {
		blog.Errorf("set sheet column width failed, objID: %s, err: %v, rid: %s", t.objID, err, t.kit.Rid)
		return err
	}

	header, err := t.handleProperty(colProps)
	if err != nil {
		blog.ErrorJSON("handle excel property failed, property: %s, err: %s, rid: %s", colProps, err, t.kit.Rid)
		return err
	}

	if err := t.excel.StreamingWrite(t.objID, core.NameRowIdx, header); err != nil {
		blog.Errorf("write excel header data to excel failed, header: %v, err: %v, rid: %s", header, err, t.kit.Rid)
		return err
	}

	if err := t.mergeHeaderCell(colProps); err != nil {
		blog.Errorf("merge excel header failed, err: %v, rid: %s", err, t.kit.Rid)
		return err
	}

	return nil
}

func (t *TmplOp) handleProperty(colProps []core.ColProp) ([][]excel.Cell, error) {
	ccLang := t.language.CreateDefaultCCLanguageIf(util.GetLanguage(t.kit.Header))

	firstColStyle, err := t.styleCreator.getStyle(noEditHeader)
	if err != nil {
		blog.Errorf("get style failed, style: %s, err: %v, rid: %s", noEditHeader, err, t.kit.Rid)
		return nil, err
	}

	width, err := core.GetRowWidth(colProps)
	if err != nil {
		blog.Errorf("get row length failed, err: %v, rid: %s", err, t.kit.Rid)
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

	for _, property := range colProps {
		if property.IsRequire {
			property.Name = property.Name + ccLang.Language("web_excel_header_required")
		}

		handleFunc := getHandleColPropFunc(&property)
		headerField, err := handleFunc(t, &property)
		if err != nil {
			blog.ErrorJSON("handle column property failed, sheet: %s, property: %s, err: %s, rid: %s", t.objID,
				property, err, t.kit.Rid)
			return nil, err
		}

		for idx, fields := range headerField {
			for fieldIdx, field := range fields {
				header[idx][property.ExcelColIndex+fieldIdx] = field
			}
		}
	}

	return header, nil
}

func (t *TmplOp) mergeHeaderCell(colProps []core.ColProp) error {
	for _, property := range colProps {
		if property.PropertyType != common.FieldTypeInnerTable {
			err := t.excel.MergeSameColCell(t.objID, property.ExcelColIndex, core.TableNameRowIdx, core.HeaderTableLen)
			if err != nil {
				blog.Errorf("merge same column cell failed, colIdx: %d, rowIdx: %d, height: %d, err: %v, rid: %s",
					property.ExcelColIndex, core.TableNameRowIdx, core.HeaderTableLen, err, t.kit.Rid)
				return err
			}
		}

		if property.PropertyType != common.FieldTypeInnerTable {
			continue
		}

		for _, idx := range rowIndexes {
			if err := t.excel.MergeSameRowCell(t.objID, property.ExcelColIndex, idx, property.Length); err != nil {
				blog.Errorf("merge same row cell failed, colIdx: %d, rowIdx: %d, length: %d, err: %v, rid: %s",
					property.ExcelColIndex, idx, property.Length, err, t.kit.Rid)
				return err
			}
		}
	}

	return nil
}

func (t *TmplOp) productAssociationHeader() error {
	if err := t.excel.CreateSheet(core.AsstSheet); err != nil {
		blog.Errorf("create sheet failed, name: %s, err: %v, rid: %s", core.AsstSheet, err, t.kit.Rid)
		return err
	}

	if err := t.setAsstColWidth(); err != nil {
		return err
	}

	if err := t.writeAsstHeader(); err != nil {
		blog.Errorf("write excel association header failed, err: %v, rid: %s", err, t.kit.Rid)
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
		if err := t.excel.SetColWidth(core.AsstSheet, idx+1, idx+1, width); err != nil {
			blog.Errorf("set sheet width failed, sheet: %s, err: %v, rid: %s", core.AsstSheet, err, t.kit.Rid)
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
	lang := t.language.CreateDefaultCCLanguageIf(util.GetLanguage(t.kit.Header))

	// 设置关联关系sheet表头第一列数据
	firstColStyle, err := t.styleCreator.getStyle(noEditHeader)
	if err != nil {
		blog.Errorf("get style failed, style: %s, err: %v, rid: %s", noEditHeader, err, t.kit.Rid)
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

	if err := t.excel.StreamingWrite(core.AsstSheet, core.AsstStartRowIdx, header); err != nil {
		blog.Errorf("write excel header data to excel failed, header: %v, err: %v, rid: %s", header, err, t.kit.Rid)
		return err
	}

	return nil
}

func (t *TmplOp) setAsstValidation() error {
	params := make([]*excel.ValidationParam, 0)

	// 设置「关联标识」的列表，定义excel选项下拉
	asstList, err := t.client.GetObjAssociation(t.kit, t.objID)
	if err != nil {
		blog.Errorf("get object association failed, err: %v, rid: %s", err, t.kit.Rid)
		return err
	}
	asstNameList := make([]string, len(asstList))
	for idx, asst := range asstList {
		asstNameList[idx] = asst.AssociationName
	}

	lang := t.language.CreateDefaultCCLanguageIf(util.GetLanguage(t.kit.Header))
	refSheet := lang.Language("excel_association_object_id")
	if err := t.excel.CreateSheet(refSheet); err != nil {
		return err
	}
	data := make([][]excel.Cell, len(asstList))
	for idx, asst := range asstList {
		data[idx] = append(data[idx], excel.Cell{Value: asst.AssociationName})
	}
	if err := t.excel.StreamingWrite(refSheet, core.AsstStartRowIdx, data); err != nil {
		return err
	}
	if err := t.excel.Flush([]string{refSheet}); err != nil {
		return err
	}
	if err := t.excel.Save(); err != nil {
		return err
	}

	sqref, err := excel.GetSingleColSqref(core.AsstIDColIdx, core.AsstDataRowIdx+1, excel.GetTotalRows())
	if err != nil {
		blog.Errorf("get single column sqref failed, err: %v, rid: %s", err, t.kit.Rid)
		return err
	}
	params = append(params, &excel.ValidationParam{Type: excel.Ref, Sqref: sqref, Option: refSheet})

	// 设置「操作」的列表，定义excel选项下拉
	sqref, err = excel.GetSingleColSqref(core.AsstOPColIdx, core.AsstDataRowIdx+1, excel.GetTotalRows())
	if err != nil {
		blog.Errorf("get single column sqref failed, err: %v, rid: %s", err, t.kit.Rid)
		return err
	}
	params = append(params, &excel.ValidationParam{Type: excel.Enum, Sqref: sqref, Option: core.AsstOps})

	for _, param := range params {
		if err = t.excel.AddValidation(core.AsstSheet, param); err != nil {
			blog.Errorf("add validation failed, sheet: %s, param: %v, err: %s, rid: %s", core.AsstSheet, param, err,
				t.kit.Rid)

			return err
		}
	}

	return nil
}
