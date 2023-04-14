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
	"context"
	"encoding/json"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/util"

	"github.com/tealeg/xlsx/v3"
)

// ProductExcelCommentSheet product excel comment sheet
func ProductExcelCommentSheet(ctx context.Context, excel *xlsx.File, defLang lang.DefaultCCLanguageIf) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	sheetName := defLang.Language(common.ExcelCommentSheetCotentLangPrefixKey + "_sheet_name")
	if sheetName == "" {
		sheetName = "comment"
	}

	sheet, err := excel.AddSheet(sheetName)
	if err != nil {
		blog.Errorf("add comment sheet error,sheet name: %s, error: %v, rid: %s", sheetName, err, rid)
		return err
	}
	strJSON := defLang.Language(common.ExcelCommentSheetCotentLangPrefixKey + "_sheet")
	if strJSON == "" {
		blog.Errorf("excel comment sheet content is empty, rid: %s", rid)
		return err
	}
	var jsSheet jsonSheet
	err = json.Unmarshal([]byte(strJSON), &jsSheet)
	if err != nil {
		blog.Errorf("excel comment sheet content not json format , rid: %s", rid)
		return err
	}

	if err := productExcelCommentSheet(jsSheet, sheet, defLang); err != nil {
		return err
	}

	return nil
}

func productExcelCommentSheet(jsSheet jsonSheet, sheet *xlsx.Sheet, defLang lang.DefaultCCLanguageIf) error {
	for idx, col := range jsSheet.Cols {
		if col == nil {
			continue
		}
		sheet.SetColWidth(idx+1, idx+1, *col.Width)
		c := sheet.Col(idx)
		c.Collapsed = col.Collapsed
		c.Hidden = col.Hidden
		if col.Style != nil {
			s := c.GetStyle()
			if s == nil {
				s = xlsx.NewStyle()
			}
			s.Alignment = col.Style.Alignment
			s.ApplyAlignment = col.Style.ApplyAlignment
			s.ApplyBorder = col.Style.ApplyBorder
			s.ApplyFill = col.Style.ApplyFill
			s.ApplyFont = col.Style.ApplyFont
			s.Border = col.Style.Border
			s.Fill = col.Style.Fill
			s.Font = col.Style.Font

			c.SetStyle(s)
		}
	}

	for idx, row := range jsSheet.Rows {
		if row == nil {
			continue
		}
		r, err := sheet.Row(idx)
		if err != nil {
			return err
		}
		if row.Height > 0 {
			r.SetHeight(row.Height)
		}
		r.Hidden = row.Hidden

		for cIdx, c := range row.Cells {
			if c == nil {
				continue
			}
			cell, err := sheet.Cell(idx, cIdx)
			if err != nil {
				return err
			}
			cell.SetFormula(c.Formula)
			cell.Hidden = c.Hidden
			cell.HMerge = c.HMerge
			value := c.Value

			if strings.HasPrefix(c.Value, "_") && !strings.HasPrefix(c.Value, "__") {
				value = defLang.Language(common.ExcelCommentSheetCotentLangPrefixKey + c.Value)
				if "" == value {
					value = c.Value
				}
			}
			cell.SetValue(value)

			cell.VMerge = c.VMerge
			if nil != c.Style {
				s := cell.GetStyle()
				if nil == s {
					s = xlsx.NewStyle()
				}
				s.Alignment = c.Style.Alignment
				s.ApplyAlignment = c.Style.ApplyAlignment
				s.ApplyBorder = c.Style.ApplyBorder
				s.ApplyFill = c.Style.ApplyFill
				s.ApplyFont = c.Style.ApplyFont
				s.Border = c.Style.Border
				s.Fill = c.Style.Fill
				s.Font = c.Style.Font

				cell.SetStyle(s)
			}
		}
	}

	return nil
}

type style struct {
	Border         xlsx.Border
	Fill           xlsx.Fill
	Font           xlsx.Font
	ApplyBorder    bool
	ApplyFill      bool
	ApplyFont      bool
	ApplyAlignment bool
	Alignment      xlsx.Alignment
}

type cell struct {
	Value   string
	Formula string
	Style   *style
	Hidden  bool
	HMerge  int
	VMerge  int
}
type col struct {
	Min       int
	Max       int
	Hidden    *bool
	Width     *float64
	Collapsed *bool
	Style     *style
}

type row struct {
	Cells        []*cell
	Hidden       bool
	Height       float64
	OutlineLevel uint8
	isCustom     bool
}

type jsonSheet struct {
	Cols []*col
	Rows []*row
}
