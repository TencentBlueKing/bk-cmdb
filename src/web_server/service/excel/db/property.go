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

package db

import (
	"fmt"

	"configcenter/pkg/excel"
	"configcenter/src/common"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
)

// ColProp excel column property
type ColProp struct {
	ID            string      // 字段标识
	Name          string      // 字段名称
	PropertyType  string      // 字段类型
	Option        interface{} // 枚举下拉列表
	IsRequire     bool        // 是否必填
	Group         string      // 字段分组
	ExcelColIndex int         // 列所处excel位置
	Length        int         // 属性需要占用多少列excel
	NotExport     bool        // 是否导出
	RefSheet      string      // 如果需要引用另一个sheet的值进行校验，那个sheet的名称
	NotEditable   bool        // 是否不可编辑
}

// GetInstWidth get instance width
func GetInstWidth(colProp []ColProp) (int, error) {
	if colProp == nil || len(colProp) == 0 {
		return 0, fmt.Errorf("column properties is invalid, val: %v", colProp)
	}

	lastProperty := colProp[len(colProp)-1]
	return lastProperty.ExcelColIndex + lastProperty.Length, nil
}

// sortColProp sort column property
func sortColProp(colProps []ColProp, groups []metadata.AttributeGroup) ([]ColProp, error) {
	props := make([]ColProp, 0)
	requiredPropMap := make(map[string][]ColProp)
	noRequiredPropMap := make(map[string][]ColProp)

	// 构造必填字段和非必填字段所在分组的map
	for _, property := range colProps {
		if property.IsRequire {
			requiredPropMap[property.Group] = append(requiredPropMap[property.Group], property)
			continue
		}
		noRequiredPropMap[property.Group] = append(noRequiredPropMap[property.Group], property)
	}

	// 第二步，根据字段分组，对必填字段排序
	index := common.HostAddMethodExcelDefaultIndex
	requiredProps, index, err := setColPropIndexAndLen(groups, requiredPropMap, index)
	if err != nil {
		return nil, err
	}
	props = append(props, requiredProps...)

	// 第三步，根据字段分组，用必填字段使用的index，继续对非必填字段进行排序
	noRequiredProps, index, err := setColPropIndexAndLen(groups, noRequiredPropMap, index)
	if err != nil {
		return nil, err
	}

	props = append(props, noRequiredProps...)

	return props, nil
}

const propertyNormalLen = 1

func setColPropIndexAndLen(groups []metadata.AttributeGroup, fieldsGroupMap map[string][]ColProp, index int) (
	[]ColProp, int, error) {

	result := make([]ColProp, 0)
	for _, group := range groups {
		fields, ok := fieldsGroupMap[group.GroupID]
		if !ok {
			continue
		}

		for _, field := range fields {
			field.ExcelColIndex = index
			if field.PropertyType == common.FieldTypeInnerTable {
				option, err := metadata.ParseTableAttrOption(field.Option)
				if err != nil {
					return nil, 0, err
				}
				index += len(option.Header)
				field.Length = len(option.Header)
				result = append(result, field)
				continue
			}
			field.Length = propertyNormalLen
			result = append(result, field)
			index++
		}
	}

	return result, index, nil
}

// GetSingleColSqref get single column sqref
func GetSingleColSqref(col int) (string, error) {
	sqref, err := excel.GetSingleColSqref(col, InstRowIdx+1, excel.GetTotalRows())
	if err != nil {
		return "", err
	}

	return sqref, nil
}

// GetTypeAliasName get type alias name
func GetTypeAliasName(ccLang language.DefaultCCLanguageIf, propertyType string) string {
	name := ccLang.Language("field_type_" + propertyType)
	if name == "" {
		return propertyType
	}

	return name
}

const instIDGroup = "default"

// GetIDProp get instance id property
func GetIDProp(colIndex int, objID string, defLang language.DefaultCCLanguageIf) ColProp {
	idProperty := ColProp{
		PropertyType:  common.FieldTypeInt,
		Group:         instIDGroup,
		ExcelColIndex: colIndex,
		Length:        propertyNormalLen,
	}

	switch objID {
	case common.BKInnerObjIDHost:
		idProperty.ID = common.BKHostIDField
		idProperty.Name = defLang.Languagef("host_property_bk_host_id")
	default:
		idProperty.ID = common.BKInstIDField
		idProperty.Name = defLang.Languagef("common_property_bk_inst_id")
	}

	return idProperty
}
