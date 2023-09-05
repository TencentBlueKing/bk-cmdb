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
	"fmt"
	"strings"

	"configcenter/pkg/excel"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/service/excel/core"
)

var handleInstFieldFuncMap = make(map[string]handleInstFieldFunc)

var handleSpecialInstFieldFuncMap = make(map[string]handleInstFieldFunc)

func init() {
	handleInstFieldFuncMap[common.FieldTypeInt] = getHandleIntFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeFloat] = getHandleFloatFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeEnum] = getHandleEnumFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeEnumMulti] = getHandleEnumMultiFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeBool] = getHandleBoolFieldFunc()
	handleInstFieldFuncMap[common.FieldTypeInnerTable] = getHandleTableFieldFunc()

	handleSpecialInstFieldFuncMap[common.BKCloudIDField] = getHandleInstCloudAreaFunc()
}

func getHandleInstFieldFunc(property *core.ColProp) handleInstFieldFunc {
	handleFunc, isSpecial := handleSpecialInstFieldFuncMap[property.ID]
	if isSpecial {
		return handleFunc
	}

	handleFunc, ok := handleInstFieldFuncMap[property.PropertyType]
	if !ok {
		handleFunc = getDefaultHandleFieldFunc()
	}

	return handleFunc
}

type handleInstFieldFunc func(e *Exporter, property *core.ColProp, val interface{}) ([][]excel.Cell, error)

const singleCellLen = 1

func getRowWithOneCell() []excel.Cell {
	return []excel.Cell{{}}
}

func getHandleIntFieldFunc() handleInstFieldFunc {
	return func(e *Exporter, property *core.ColProp, val interface{}) ([][]excel.Cell, error) {
		intVal, err := util.GetInt64ByInterface(val)
		if err != nil {
			blog.Errorf("value type is not int, val: %v", val)
			return [][]excel.Cell{getRowWithOneCell()}, nil
		}

		handleFunc := getDefaultHandleFieldFunc()
		return handleFunc(e, property, intVal)
	}
}

func getHandleFloatFieldFunc() handleInstFieldFunc {
	return func(e *Exporter, property *core.ColProp, val interface{}) ([][]excel.Cell, error) {
		floatVal, err := util.GetFloat64ByInterface(val)
		if err != nil {
			blog.Errorf("value type is not float64, val: %v", val)
			return [][]excel.Cell{getRowWithOneCell()}, nil
		}

		handleFunc := getDefaultHandleFieldFunc()
		return handleFunc(e, property, floatVal)
	}
}

func getHandleEnumFieldFunc() handleInstFieldFunc {
	return func(e *Exporter, property *core.ColProp, val interface{}) ([][]excel.Cell, error) {
		option, err := metadata.ParseEnumOption(property.Option)
		if err != nil {
			blog.Errorf("option type is invalid, option: %v", property.Option)
			return [][]excel.Cell{getRowWithOneCell()}, nil
		}

		enumID, ok := val.(string)
		if !ok {
			blog.Errorf("val type is invalid, val: %v", val)
			return [][]excel.Cell{getRowWithOneCell()}, nil
		}

		handleFunc := getDefaultHandleFieldFunc()
		return handleFunc(e, property, getEnumNameByID(enumID, option))
	}
}

func getHandleEnumMultiFieldFunc() handleInstFieldFunc {
	return func(e *Exporter, property *core.ColProp, val interface{}) ([][]excel.Cell, error) {
		if val == nil {
			return [][]excel.Cell{getRowWithOneCell()}, nil
		}

		option, err := metadata.ParseEnumOption(property.Option)
		if err != nil {
			return nil, err
		}

		enumArr, ok := val.([]interface{})
		if !ok {
			return nil, fmt.Errorf("convert enum multiple type value failed, val: %v", val)
		}

		enumMultiName := make([]string, 0)
		for _, enumID := range enumArr {
			id, ok := enumID.(string)
			if !ok {
				return nil, fmt.Errorf("convert enum multiple id [%v] to string failed", enumID)
			}

			name := getEnumNameByID(id, option)
			enumMultiName = append(enumMultiName, name)
		}

		val = strings.Join(enumMultiName, "\n")

		handleFunc := getDefaultHandleFieldFunc()
		return handleFunc(e, property, val)
	}
}

// getEnumNameByID get enum name from option
func getEnumNameByID(id string, option metadata.EnumOption) string {
	for _, item := range option {
		if item.ID != id {
			continue
		}

		return item.Name
	}

	return ""
}

func getHandleBoolFieldFunc() handleInstFieldFunc {
	return func(e *Exporter, property *core.ColProp, val interface{}) ([][]excel.Cell, error) {
		val, ok := val.(bool)
		if !ok {
			blog.Errorf("value type is not boolean, val: %v", val)
			return [][]excel.Cell{getRowWithOneCell()}, nil
		}

		handleFunc := getDefaultHandleFieldFunc()
		return handleFunc(e, property, val)
	}
}

func getHandleTableFieldFunc() handleInstFieldFunc {
	return func(e *Exporter, property *core.ColProp, val interface{}) ([][]excel.Cell, error) {
		table, ok := val.([]mapstr.MapStr)
		if !ok {
			return nil, fmt.Errorf("transfer table struct failed, val: %v", val)
		}
		option, err := metadata.ParseTableAttrOption(property.Option)
		if err != nil {
			return nil, err
		}

		tableRows := make([][]excel.Cell, len(table))
		for idx, data := range table {
			for _, attr := range option.Header {
				handleFunc, ok := handleInstFieldFuncMap[attr.PropertyType]
				if !ok {
					handleFunc = getDefaultHandleFieldFunc()
				}

				colProp := &core.ColProp{ID: attr.PropertyID, Name: attr.PropertyName, PropertyType: attr.PropertyType,
					IsRequire: attr.IsRequired, Option: attr.Option, Group: attr.PropertyGroup}
				rows, err := handleFunc(e, colProp, data[attr.PropertyID])
				if err != nil {
					blog.ErrorJSON("handle instance failed, property: %s, val: %s, err: %s, rid: %s", property, val,
						err, e.GetKit().Rid)
					return nil, err
				}

				if len(rows) != singleCellLen {
					blog.ErrorJSON("instance table field is invalid, property: %s, val: %s, err: %s, rid: %s", property,
						val, err, e.GetKit().Rid)
					return nil, err
				}

				tableRows[idx] = append(tableRows[idx], rows[0]...)
			}
		}

		return tableRows, nil
	}
}

func getDefaultHandleFieldFunc() handleInstFieldFunc {
	return func(e *Exporter, property *core.ColProp, val interface{}) ([][]excel.Cell, error) {
		var styleID int
		if property.NotEditable {
			var err error
			styleID, err = e.styleCreator.getStyle(noEditField)
			if err != nil {
				return nil, err
			}
		}

		result := make([][]excel.Cell, singleCellLen)
		result[0] = append(result[0], excel.Cell{Value: val, StyleID: styleID})

		return result, nil
	}
}

func getHandleInstCloudAreaFunc() handleInstFieldFunc {
	return func(e *Exporter, property *core.ColProp, val interface{}) ([][]excel.Cell, error) {
		cloudArr, err := util.GetMapInterfaceByInterface(val)
		if err != nil {
			return nil, err
		}

		if len(cloudArr) != 1 {
			blog.Errorf("host has many cloud areas, val: %#v, rid: %s", val, e.GetKit().Rid)
			return nil, e.GetKit().CCError.CCError(common.CCErrCommReplyDataFormatError)
		}

		cloudMap, ok := cloudArr[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("transfer cloud area failed, val: %v", cloudArr[0])
		}

		handleFunc := getDefaultHandleFieldFunc()
		return handleFunc(e, property, cloudMap[common.BKInstNameField])
	}
}
