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

package model

import (
	"encoding/json"
	"fmt"
	"strconv"

	"configcenter/pkg/excel"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/service/excel/core"
)

// Operator model operator
type Operator struct {
	excel    *excel.Excel
	client   *core.Client
	objID    string
	kit      *rest.Kit
	language language.CCLanguageIf
}

type BuildModelOpFunc func(modelOp *Operator) error

// NewOp create a excel model operator
func NewOp(opts ...BuildModelOpFunc) (*Operator, error) {
	modelOp := new(Operator)
	for _, opt := range opts {
		if err := opt(modelOp); err != nil {
			return nil, err
		}
	}

	return modelOp, nil
}

// FilePath set model file path
func FilePath(filePath string) BuildModelOpFunc {
	return func(modelOp *Operator) error {
		var err error
		modelOp.excel, err = excel.NewExcel(excel.FilePath(filePath), excel.OpenOrCreate(), excel.DelDefaultSheet())
		if err != nil {
			return err
		}

		return nil
	}
}

// Client set model client
func Client(dao *core.Client) BuildModelOpFunc {
	return func(modelOp *Operator) error {
		modelOp.client = dao
		return nil
	}
}

// ObjID set model object id
func ObjID(objID string) BuildModelOpFunc {
	return func(modelOp *Operator) error {
		modelOp.objID = objID
		return nil
	}
}

// Kit set model kit
func Kit(kit *rest.Kit) BuildModelOpFunc {
	return func(modelOp *Operator) error {
		modelOp.kit = kit
		return nil
	}
}

// Language set model language
func Language(language language.CCLanguageIf) BuildModelOpFunc {
	return func(modelOp *Operator) error {
		modelOp.language = language
		return nil
	}
}

// Close excel
func (op *Operator) Close() error {
	if err := op.excel.Flush(op.objID); err != nil {
		blog.Errorf("flush excel failed, sheet %s, err: %v, rid: %s", op.objID, err, op.kit.Rid)
		return err
	}

	if err := op.excel.Save(); err != nil {
		blog.Errorf("save excel failed, err: %v, rid: %s", err, op.kit.Rid)
		return err
	}

	if err := op.excel.Close(); err != nil {
		blog.Errorf("close excel failed, err: %v, rid: %s", err, op.kit.Rid)
		return err
	}

	return nil
}

// Clean delete temporary file
func (op *Operator) Clean() error {
	return op.excel.Clean()
}

// Export export data to excel
func (op *Operator) Export() error {
	if err := op.excel.CreateSheet(op.objID); err != nil {
		blog.Errorf("create sheet failed, objID: %s, err: %v, rid: %s", op.objID, err, op.kit.Rid)
		return err
	}

	if err := op.setExcelTitle(); err != nil {
		blog.Errorf("set excel title failed, err: %v, rid: %s", err, op.kit.Rid)
		return err
	}

	if err := op.setExcelData(); err != nil {
		blog.Errorf("set excel data failed, err: %v, rid: %s", err, op.kit.Rid)
		return err
	}

	return nil
}

const (
	headerLen = 3
	descIdx   = 0
	typeIdx   = 1
	idIdx     = 2
	dataIdx   = 3
)

func (op *Operator) setExcelTitle() error {
	lang := op.language.CreateDefaultCCLanguageIf(util.GetLanguage(op.kit.Header))
	fields := getSortFields(lang)

	header := make([][]excel.Cell, headerLen)
	for i := range header {
		header[i] = make([]excel.Cell, len(fields))
	}

	for idx, field := range fields {
		header[descIdx][idx].Value = field.desc
		header[typeIdx][idx].Value = field.fType
		header[idIdx][idx].Value = field.id
	}

	if err := op.excel.StreamingWrite(op.objID, descIdx, header); err != nil {
		blog.Errorf("write excel header data to excel failed, header: %v, err: %v, rid: %s", header, err, op.kit.Rid)
		return err
	}

	return nil
}

type fieldBriefMsg struct {
	id    string
	desc  string
	fType string
}

func getSortFields(lang language.DefaultCCLanguageIf) []fieldBriefMsg {
	textType := lang.Language("val_type_text")
	boolType := lang.Language("val_type_bool")

	var fields = []fieldBriefMsg{
		{id: common.BKPropertyIDField, desc: lang.Language("web_en_name_required"), fType: textType},
		{id: common.BKPropertyNameField, desc: lang.Language("web_bk_alias_name_required"), fType: textType},
		{id: common.BKPropertyTypeField, desc: lang.Language("web_bk_data_type_required"), fType: textType},
		{id: common.BKPropertyGroupField, desc: lang.Language("property_group"), fType: textType},
		{id: common.BKOptionField, desc: lang.Language("property_option"), fType: textType},
		{id: common.BKUnitField, desc: lang.Language("unit"), fType: textType},
		{id: common.BKDescriptionField, desc: lang.Language("desc"), fType: textType},
		{id: common.BKPlaceholderField, desc: lang.Language("placeholder"), fType: textType},
		{id: common.BKEditableField, desc: lang.Language("is_editable"), fType: boolType},
		{id: common.BKIsRequiredField, desc: lang.Language("property_is_required"), fType: boolType},
		{id: common.BKIsreadonlyField, desc: lang.Language("property_is_readonly"), fType: boolType},
		{id: common.BKIsOnlyField, desc: lang.Language("property_is_only"), fType: boolType},
		{id: common.BKIsMultipleField, desc: lang.Language("property_is_multiple"), fType: boolType},
		{id: common.BKDefaultFiled, desc: lang.Language("property_default"), fType: textType},
	}

	return fields
}

func (op *Operator) setExcelData() error {
	attrs, err := op.client.GetObjectData(op.kit, op.objID)
	if err != nil {
		blog.Errorf("get excel object data failed, err: %v, rid: %s", err, op.kit.Rid)
		return err
	}

	lang := op.language.CreateDefaultCCLanguageIf(util.GetLanguage(op.kit.Header))
	fields := getSortFields(lang)
	data := make([][]excel.Cell, len(attrs))
	for i := range data {
		data[i] = make([]excel.Cell, len(fields))
	}

	for rowIdx, attr := range attrs {
		row, ok := attr.(map[string]interface{})
		if !ok {
			return fmt.Errorf("object attribute is invalid, val: %v", attr)
		}

		delete(row, common.BKTemplateID)

		for idx, field := range fields {
			id := field.id
			cell := row[id]

			if id == common.BKOptionField || id == common.BKDefaultFiled {
				if cell == nil {
					data[rowIdx][idx].Value = ""
					continue
				}

				value, err := json.Marshal(cell)
				if err != nil {
					blog.Errorf("value is invalid, field: %s, val: %v, err: %v, rid: %s", id, cell, err, op.kit.Rid)
					data[rowIdx][idx].Value = "error info:" + err.Error()
					continue
				}

				data[rowIdx][idx].Value = string(value)
				continue
			}

			data[rowIdx][idx].Value = cell
		}
	}

	if err := op.excel.StreamingWrite(op.objID, dataIdx, data); err != nil {
		blog.Errorf("write excel data to excel failed, header: %v, err: %v, rid: %s", data, err, op.kit.Rid)
		return err
	}

	return nil
}

// Import object attributes
func (op *Operator) Import() (*metadata.Response, error) {
	attrs, err := op.getImportAttr()
	if err != nil {
		blog.Errorf("get imported attribute failed, err: %v, rid: %s", err, op.kit.Rid)
		return nil, err
	}

	attrs = convAttr(attrs)

	param := map[string]interface{}{op.objID: map[string]interface{}{"attr": attrs}}
	result, err := op.client.AddObjectBatch(op.kit, param)
	if err != nil {
		blog.ErrorJSON("add object attribute failed, attrs: %s, err: %s, rid: %s", attrs, err, op.kit.Rid)
		return nil, err
	}

	return result, nil
}

func (op *Operator) getImportAttr() (map[int]map[string]interface{}, error) {
	reader, err := op.excel.NewReader(op.objID)
	if err != nil {
		blog.Errorf("create excel io reader failed, sheet: %s, err: %v, rid: %s", op.objID, err, op.kit.Rid)
		return nil, err
	}

	// get the column location where the id is located
	idMap := make(map[int]string)
	for reader.Next() {
		if reader.GetCurIdx() < idIdx {
			continue
		}

		row, err := reader.CurRow()
		if err != nil {
			blog.Errorf("get reader current row data failed, err: %v, rid: %s", err, op.kit.Rid)
			return nil, err
		}

		for idx, val := range row {
			idMap[idx] = val
		}
		break
	}

	lang := op.language.CreateDefaultCCLanguageIf(util.GetLanguage(op.kit.Header))
	fields := getSortFields(lang)
	idTypeMap := make(map[string]string)
	for _, field := range fields {
		idTypeMap[field.id] = field.fType
	}
	boolType := lang.Language("val_type_bool")

	// get attribute
	attrs := make(map[int]map[string]interface{})
	for reader.Next() {
		if reader.GetCurIdx() < dataIdx {
			continue
		}

		row, err := reader.CurRow()
		if err != nil {
			blog.Errorf("get reader current row data failed, err: %v, rid: %s", err, op.kit.Rid)
			return nil, err
		}

		attr := make(map[string]interface{})
		for idx, val := range row {
			id, ok := idMap[idx]
			if !ok {
				continue
			}

			if idTypeMap[id] != boolType {
				attr[id] = val
				continue
			}

			boolVal, err := strconv.ParseBool(val)
			if err != nil {
				blog.Errorf("convert string to bool type failed, val: %v, err: %v, rid: %s", val, err, op.kit.Rid)
				return nil, err
			}
			attr[id] = boolVal
		}

		attrs[reader.GetCurIdx()+1] = attr
	}

	return attrs, nil
}

// convAttr convert attribute in excel to cmdb attributes
func convAttr(attrItems map[int]map[string]interface{}) map[int]map[string]interface{} {
	for index, attr := range attrItems {
		fieldType, ok := attr[common.BKPropertyTypeField].(string)
		if !ok {
			continue
		}

		val, ok := attr[common.BKOptionField].(string)
		if ok && val == "\"\"" {
			attrItems[index][common.BKOptionField] = ""
		}

		switch fieldType {
		case common.FieldTypeEnum, common.FieldTypeList, common.FieldTypeEnumMulti, common.FieldTypeEnumQuote,
			common.FieldTypeInnerTable:
			var iOption interface{}
			attrItems[index] = unmarshalAttrStrVal(attrItems[index], common.BKOptionField, iOption)
		case common.FieldTypeInt:
			iOption := make(map[string]interface{})
			attrItems[index] = unmarshalAttrStrVal(attrItems[index], common.BKOptionField, iOption)

			var iDefault int64
			attrItems[index] = unmarshalAttrStrVal(attrItems[index], common.BKDefaultFiled, iDefault)
		case common.FieldTypeFloat:
			iOption := make(map[string]interface{})
			attrItems[index] = unmarshalAttrStrVal(attrItems[index], common.BKOptionField, iOption)

			var iDefault float64
			attrItems[index] = unmarshalAttrStrVal(attrItems[index], common.BKDefaultFiled, iDefault)
		case common.FieldTypeOrganization:
			iDefault := make([]interface{}, 0)
			attrItems[index] = unmarshalAttrStrVal(attrItems[index], common.BKDefaultFiled, iDefault)
		case common.FieldTypeBool:
			var iOption bool
			attrItems[index] = unmarshalAttrStrVal(attrItems[index], common.BKOptionField, iOption)
		}
	}

	return attrItems
}

func unmarshalAttrStrVal(attr map[string]interface{}, field string, value interface{}) map[string]interface{} {
	val, ok := attr[field].(string)
	if !ok {
		return attr
	}

	if val == "\"\"" {
		attr[field] = value
		return attr
	}

	err := json.Unmarshal([]byte(val), &value)
	if err != nil {
		return attr
	}

	attr[field] = value
	return attr
}
