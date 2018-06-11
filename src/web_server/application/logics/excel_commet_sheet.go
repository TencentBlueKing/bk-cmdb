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
	"encoding/json"
	"strings"

	//simplejson "github.com/bitly/go-simplejson"
	"github.com/rentiansheng/xlsx"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
)

// ProductExcelHealer Excel comment sheet，
func ProductExcelCommentSheet(excel *xlsx.File, defLang lang.DefaultCCLanguageIf) {
	sheetName := defLang.Language(common.ExcelCommentSheetCotentLangPrefixKey + "_sheet_name")
	if "" == sheetName {
		sheetName = "comment"
	}

	sheet, err := excel.AddSheet(sheetName)
	if nil != err {
		blog.Errorf("add comment sheet error,sheet name:%s, error:%s ", sheetName, err.Error())
		return
	}
	strJSON := defLang.Language(common.ExcelCommentSheetCotentLangPrefixKey + "_sheet")
	if "" == strJSON {
		blog.Errorf("excel comment sheet content is empty")
		return
	}
	var jsSheet jsonSheet
	err = json.Unmarshal([]byte(strJSON), &jsSheet)
	if nil != err {
		blog.Errorf("excel comment sheet content not json format ")
		return
	}

	for idx, col := range jsSheet.Cols {
		if nil == col {
			continue
		}
		c := sheet.Col(idx)
		c.Collapsed = c.Collapsed
		c.Hidden = c.Hidden
		c.Max = col.Max
		c.Min = col.Min
		c.Width = col.Width
		if nil != col.Style {
			s := c.GetStyle()
			if nil == s {
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
		if nil == row {
			continue
		}
		r := sheet.Row(idx)
		if 0 < row.Height {
			r.SetHeight(row.Height)

		}
		r.Hidden = row.Hidden

		for cIdx, c := range row.Cells {
			if nil == c {
				continue
			}
			cell := sheet.Cell(idx, cIdx)
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
	Hidden    bool
	Width     float64
	Collapsed bool
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
