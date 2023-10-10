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
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/service/excel/core"
	"configcenter/src/web_server/service/excel/operator"
)

// Operator model operator
type Operator struct {
	*operator.BaseOp
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

// BaseOperator set base operator
func BaseOperator(op *operator.BaseOp) BuildModelOpFunc {
	return func(modelOp *Operator) error {
		modelOp.BaseOp = op
		return nil
	}
}

// Close excel
func (op *Operator) Close() error {
	if err := op.GetExcel().Flush([]string{op.GetObjID(), core.AsstSheet}); err != nil {
		blog.Errorf("flush excel failed, sheet %s, err: %v, rid: %s", op.GetObjID(), err, op.GetKit().Rid)
		return err
	}

	if err := op.GetExcel().Save(); err != nil {
		blog.Errorf("save excel failed, err: %v, rid: %s", err, op.GetKit().Rid)
		return err
	}

	if err := op.GetExcel().Close(); err != nil {
		blog.Errorf("close excel failed, err: %v, rid: %s", err, op.GetKit().Rid)
		return err
	}

	return nil
}

// Clean delete temporary file
func (op *Operator) Clean() error {
	return op.GetExcel().Clean()
}

// Export export data to excel
func (op *Operator) Export() error {
	if err := op.GetExcel().CreateSheet(op.GetObjID()); err != nil {
		blog.Errorf("create sheet failed, objID: %s, err: %v, rid: %s", op.GetObjID(), err, op.GetKit().Rid)
		return err
	}

	if err := op.setExcelTitle(); err != nil {
		blog.Errorf("set excel title failed, err: %v, rid: %s", err, op.GetKit().Rid)
		return err
	}

	if err := op.setExcelData(); err != nil {
		blog.Errorf("set excel data failed, err: %v, rid: %s", err, op.GetKit().Rid)
		return err
	}

	return nil
}

const (
	// headerLen excel模型数据表头长度
	headerLen = 3
	// descIdx excel模型字段名称所在行位置
	descIdx = 0
	// typeIdx excel模型类型所在行位置
	typeIdx = 1
	// idIdx excel模型id所在行位置
	idIdx = 2
	// dataIdx excel模型数据开始所在行位置
	dataIdx = 3
)

func (op *Operator) setExcelTitle() error {
	lang := op.GetLang().CreateDefaultCCLanguageIf(util.GetLanguage(op.GetKit().Header))
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

	if err := op.GetExcel().StreamingWrite(op.GetObjID(), descIdx, header); err != nil {
		blog.Errorf("write excel header data to excel failed, header: %v, err: %v, rid: %s", header, err,
			op.GetKit().Rid)
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
		{id: metadata.AttributeFieldPropertyID, desc: lang.Language("web_en_name_required"), fType: textType},
		{id: metadata.AttributeFieldPropertyName, desc: lang.Language("web_bk_alias_name_required"), fType: textType},
		{id: metadata.AttributeFieldPropertyType, desc: lang.Language("web_bk_data_type_required"), fType: textType},
		{id: metadata.AttributeFieldPropertyGroup, desc: lang.Language("property_group"), fType: textType},
		{id: metadata.AttributeFieldOption, desc: lang.Language("property_option"), fType: textType},
		{id: metadata.AttributeFieldUnit, desc: lang.Language("unit"), fType: textType},
		{id: metadata.AttributeFieldDescription, desc: lang.Language("desc"), fType: textType},
		{id: metadata.AttributeFieldPlaceHolder, desc: lang.Language("placeholder"), fType: textType},
		{id: metadata.AttributeFieldIsEditable, desc: lang.Language("is_editable"), fType: boolType},
		{id: metadata.AttributeFieldIsRequired, desc: lang.Language("property_is_required"), fType: boolType},
		{id: metadata.AttributeFieldIsReadOnly, desc: lang.Language("property_is_readonly"), fType: boolType},
		{id: metadata.AttributeFieldIsOnly, desc: lang.Language("property_is_only"), fType: boolType},
		{id: metadata.AttributeFieldIsMultiple, desc: lang.Language("property_is_multiple"), fType: boolType},
		{id: metadata.AttributeFieldDefault, desc: lang.Language("property_default"), fType: textType},
	}

	return fields
}

func (op *Operator) setExcelData() error {
	attrs, err := op.GetClient().GetObjectData(op.GetKit(), op.GetObjID())
	if err != nil {
		blog.Errorf("get excel object data failed, err: %v, rid: %s", err, op.GetKit().Rid)
		return err
	}

	lang := op.GetLang().CreateDefaultCCLanguageIf(util.GetLanguage(op.GetKit().Header))
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
					blog.Errorf("value is invalid, field: %s, val: %v, err: %v, rid: %s", id, cell, err,
						op.GetKit().Rid)
					data[rowIdx][idx].Value = "error info:" + err.Error()
					continue
				}

				data[rowIdx][idx].Value = string(value)
				continue
			}

			data[rowIdx][idx].Value = cell
		}
	}

	if err := op.GetExcel().StreamingWrite(op.GetObjID(), dataIdx, data); err != nil {
		blog.Errorf("write excel data to excel failed, header: %v, err: %v, rid: %s", data, err, op.GetKit().Rid)
		return err
	}

	return nil
}

// Import object attributes
func (op *Operator) Import() (*metadata.Response, error) {
	attrs, err := op.getImportAttr()
	if err != nil {
		blog.Errorf("get imported attribute failed, err: %v, rid: %s", err, op.GetKit().Rid)
		return nil, err
	}

	attrs = convAttr(attrs)

	param := map[string]interface{}{op.GetObjID(): map[string]interface{}{"attr": attrs}}
	result, err := op.GetClient().AddObjectBatch(op.GetKit(), param)
	if err != nil {
		blog.ErrorJSON("add object attribute failed, attrs: %s, err: %s, rid: %s", attrs, err, op.GetKit().Rid)
		return nil, err
	}

	return result, nil
}

func (op *Operator) getImportAttr() (map[int]map[string]interface{}, error) {
	reader, err := op.GetExcel().NewReader(op.GetObjID())
	if err != nil {
		blog.Errorf("create excel io reader failed, sheet: %s, err: %v, rid: %s", op.GetObjID(), err, op.GetKit().Rid)
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
			blog.Errorf("get reader current row data failed, err: %v, rid: %s", err, op.GetKit().Rid)
			return nil, err
		}

		for idx, val := range row {
			idMap[idx] = val
		}
		break
	}

	lang := op.GetLang().CreateDefaultCCLanguageIf(util.GetLanguage(op.GetKit().Header))
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
			blog.Errorf("get reader current row data failed, err: %v, rid: %s", err, op.GetKit().Rid)
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
				blog.Errorf("convert string to bool type failed, val: %v, err: %v, rid: %s", val, err, op.GetKit().Rid)
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
