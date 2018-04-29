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
	"configcenter/src/common"
	"github.com/tealeg/xlsx"
)

const (
	fieldTypeBoolTrue  = "true"
	fieldTypeBoolFalse = "false"
)

// getFieldsIDIndexMap get field property index
func getFieldsIDIndexMap(fields map[string]Property) map[string]int {
	index := 0
	IDNameMap := make(map[string]int)
	for id := range fields {
		IDNameMap[id] = index
		index++
	}
	return IDNameMap
}

// getAssociateName  get getAssociate object name
func getAssociateNames(a []interface{}) []string {
	vals := []string{}

	for _, valRow := range a {
		mapVal, ok := valRow.(map[string]interface{})
		if ok {
			nameval, ok := mapVal[common.BKInstNameField].(string)
			if true == ok {
				vals = append(vals, nameval)
			}
		}
	}

	return vals
}

// getEnumNameByID get enum name from option
func getEnumNameByID(id string, items []interface{}) string {
	var name string
	for _, valRow := range items {
		mapVal, ok := valRow.(map[string]interface{})
		if ok {
			enumID, ok := mapVal["id"].(string)
			if true == ok {
				if enumID == id {
					name = mapVal["name"].(string)
				}
			}
		}
	}

	return name
}

// getEnumNames get enum name from option
func getEnumNames(items []interface{}) []string {
	var names []string
	for _, valRow := range items {
		mapVal, ok := valRow.(map[string]interface{})
		if ok {

			name, ok := mapVal["name"].(string)
			if ok {
				names = append(names, name)
			}

		}
	}

	return names
}

// getHeaderCellGeneralStyle get excel header general style by C6EFCE,000000
func getHeaderCellGeneralStyle() *xlsx.Style {
	return getCellStyle(common.ExcelHeaderOtherRowColor, common.ExcelHeaderOtherRowFontColor)
}

// getHeaderFirstRowCellStyle
func getHeaderFirstRowCellStyle(isRequire bool) *xlsx.Style {
	if isRequire {
		return getCellStyle(common.ExcelHeaderFirstRowColor, common.ExcelHeaderFirstRowRequireFontColor)
	}

	return getCellStyle(common.ExcelHeaderFirstRowColor, common.ExcelHeaderFirstRowFontColor)
}

// getCellStyle get cell style from fgColor and fontcolor
func getCellStyle(fgColor, fontColor string) *xlsx.Style {
	style := xlsx.NewStyle()
	style.Fill = *xlsx.DefaultFill()
	style.Font = *xlsx.DefaultFont()
	style.ApplyFill = true
	style.ApplyFont = true
	style.ApplyBorder = true

	style.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	style.Border.BottomColor = common.ExcelCellDefaultBorderColor
	style.Border.TopColor = common.ExcelCellDefaultBorderColor
	style.Border.LeftColor = common.ExcelCellDefaultBorderColor
	style.Border.RightColor = common.ExcelCellDefaultBorderColor

	style.Fill.FgColor = fgColor
	style.Fill.PatternType = "solid"

	style.Font.Color = fontColor

	return style
}
