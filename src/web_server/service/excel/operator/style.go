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
)

type styleType string

const (
	noEditHeader  styleType = "noEditHeader"
	fieldName     styleType = "fieldName"
	tableHeader   styleType = "tableHeader"
	generalHeader styleType = "generalHeader"
	noEditField   styleType = "noEditField"

	noEditHeaderColor  = "fabf8f"
	noEditFieldColor   = "fee9da"
	fieldNameColor     = "92d050"
	tableHeaderColor   = "d1e0b6"
	generalHeaderColor = "c6efce"
	borderColor        = "d4d4d4"
)

var generalBorder = []excel.Border{
	{Type: excel.Left, Color: borderColor, Style: 1}, {Type: excel.Right, Color: borderColor, Style: 1},
	{Type: excel.Top, Color: borderColor, Style: 1}, {Type: excel.Bottom, Color: borderColor, Style: 1},
}

var createStyleFuncMap = make(map[styleType]createStyleFunc)

func init() {
	createStyleFuncMap[noEditHeader] = getNoEditHeaderStyleFunc()
	createStyleFuncMap[noEditField] = getNoEditFieldStyleFunc()
	createStyleFuncMap[fieldName] = getFieldNameStyleFunc()
	createStyleFuncMap[generalHeader] = getGeneralHeaderStyleFunc()
	createStyleFuncMap[tableHeader] = getTableHeaderStyleFunc()
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

func getFieldNameStyleFunc() createStyleFunc {
	return func(s *styleCreator) (int, error) {
		style := &excel.Style{Fill: &excel.Fill{Type: excel.Pattern, Color: []string{fieldNameColor}, Pattern: 1},
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
