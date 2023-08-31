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
)

type styleType string

const (
	// requiredField 必填的表头单元格类型
	requiredField styleType = "requiredField"
	// noEditHeader 不可编辑的表头单元格类型
	noEditHeader styleType = "noEditHeader"
	// firstRow 第一行数据的单元格类型
	firstRow styleType = "firstRow"
	// tableHeader 表格表头的单元格类型
	tableHeader styleType = "tableHeader"
	// generalHeader 表头正常单元格的类型
	generalHeader styleType = "generalHeader"
	// noEditField 不可编辑的单元格类型
	noEditField styleType = "noEditField"
	// example 例子数据的单元格类型
	example styleType = "example"

	requiredFieldColor = "#FF0000"
	noEditHeaderColor  = "fabf8f"
	noEditFieldColor   = "fee9da"
	firstRowColor      = "92d050"
	tableHeaderColor   = "d1e0b6"
	generalHeaderColor = "c6efce"
	borderColor        = "d4d4d4"
	exampleColor       = "c6efce"
)

var generalBorder = []excel.Border{
	{Type: excel.Left, Color: borderColor, Style: 1}, {Type: excel.Right, Color: borderColor, Style: 1},
	{Type: excel.Top, Color: borderColor, Style: 1}, {Type: excel.Bottom, Color: borderColor, Style: 1},
}

var createStyleFuncMap = make(map[styleType]createStyleFunc)

func init() {
	createStyleFuncMap[requiredField] = getRequiredFieldStyleFunc()
	createStyleFuncMap[noEditHeader] = getNoEditHeaderStyleFunc()
	createStyleFuncMap[noEditField] = getNoEditFieldStyleFunc()
	createStyleFuncMap[firstRow] = getFirstRowStyleFunc()
	createStyleFuncMap[generalHeader] = getGeneralHeaderStyleFunc()
	createStyleFuncMap[tableHeader] = getTableHeaderStyleFunc()
	createStyleFuncMap[example] = getExampleStyleFunc()
}

type createStyleFunc func(s *styleCreator) (int, error)

func getNoEditHeaderStyleFunc() createStyleFunc {
	return func(s *styleCreator) (int, error) {
		style := &excel.Style{Fill: &excel.Fill{Type: "pattern", Color: []string{noEditHeaderColor}, Pattern: 1},
			Border: generalBorder}

		result, err := s.excel.NewStyle(style)
		if err != nil {
			return 0, err
		}

		return result, nil
	}
}

func getNoEditFieldStyleFunc() createStyleFunc {
	return func(s *styleCreator) (int, error) {
		style := &excel.Style{Fill: &excel.Fill{Type: excel.Pattern, Color: []string{noEditFieldColor}, Pattern: 1},
			Border: generalBorder}

		result, err := s.excel.NewStyle(style)
		if err != nil {
			return 0, err
		}

		return result, nil
	}
}

func getFirstRowStyleFunc() createStyleFunc {
	return func(s *styleCreator) (int, error) {
		style := &excel.Style{Fill: &excel.Fill{Type: excel.Pattern, Color: []string{firstRowColor}, Pattern: 1},
			Border: generalBorder}

		result, err := s.excel.NewStyle(style)
		if err != nil {
			return 0, err
		}

		return result, nil
	}
}

func getGeneralHeaderStyleFunc() createStyleFunc {
	return func(s *styleCreator) (int, error) {
		style := &excel.Style{Fill: &excel.Fill{Type: excel.Pattern, Color: []string{generalHeaderColor}, Pattern: 1},
			Border: generalBorder}

		result, err := s.excel.NewStyle(style)
		if err != nil {
			return 0, err
		}

		return result, nil
	}
}

func getTableHeaderStyleFunc() createStyleFunc {
	return func(s *styleCreator) (int, error) {
		style := &excel.Style{Fill: &excel.Fill{Type: excel.Pattern, Color: []string{tableHeaderColor}, Pattern: 1},
			Border: generalBorder}

		result, err := s.excel.NewStyle(style)
		if err != nil {
			return 0, err
		}

		return result, nil
	}
}

func getExampleStyleFunc() createStyleFunc {
	return func(s *styleCreator) (int, error) {
		style := &excel.Style{Fill: &excel.Fill{Type: excel.Pattern, Color: []string{exampleColor}, Pattern: 1},
			Border: generalBorder}

		result, err := s.excel.NewStyle(style)
		if err != nil {
			return 0, err
		}

		return result, nil
	}
}

func getRequiredFieldStyleFunc() createStyleFunc {
	return func(s *styleCreator) (int, error) {
		style := &excel.Style{Fill: &excel.Fill{Type: excel.Pattern, Color: []string{firstRowColor}, Pattern: 1},
			Border: generalBorder, Font: &excel.Font{Color: requiredFieldColor}}

		result, err := s.excel.NewStyle(style)
		if err != nil {
			return 0, err
		}

		return result, nil
	}
}

type styleCreator struct {
	excel    *excel.Excel
	styleMap map[styleType]int
}

type styleOperatorFunc func(style *styleCreator) error

func newStyleCreator(opts ...styleOperatorFunc) (*styleCreator, error) {
	style := &styleCreator{
		styleMap: make(map[styleType]int),
	}
	for _, opt := range opts {
		if err := opt(style); err != nil {
			return nil, err
		}
	}

	return style, nil
}

func setExcel(excel *excel.Excel) styleOperatorFunc {
	return func(style *styleCreator) error {
		style.excel = excel
		return nil
	}
}

func (s *styleCreator) getStyle(style styleType) (int, error) {
	result, ok := s.styleMap[style]
	if !ok {
		styleFunc := createStyleFuncMap[style]
		var err error
		result, err = styleFunc(s)
		if err != nil {
			return 0, err
		}
	}

	return result, nil
}
