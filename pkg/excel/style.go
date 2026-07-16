/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package excel

import (
	"github.com/xuri/excelize/v2"
)

// Style excel style
type Style struct {
	Fill   *Fill
	Border []Border
	Font   *Font
	NumFmt int
}

func (s *Style) convert() (*excelize.Style, error) {
	style := new(excelize.Style)
	style.NumFmt = s.NumFmt
	if s.Fill != nil {
		fill, err := s.Fill.convert()
		if err != nil {
			return nil, err
		}
		style.Fill = fill
	}

	if s.Font != nil {
		font, err := s.Font.convert()
		if err != nil {
			return nil, err
		}
		style.Font = font
	}

	if s.Border != nil {
		for _, border := range s.Border {
			excelBorder, err := border.convert()
			if err != nil {
				return nil, err
			}
			style.Border = append(style.Border, excelBorder)
		}
	}

	return style, nil
}

// Alignment directly maps the alignment settings of the cells.
type Alignment struct {
	Horizontal      string
	Indent          int
	JustifyLastLine bool
	ReadingOrder    uint64
	RelativeIndent  int
	ShrinkToFit     bool
	TextRotation    int
	Vertical        string
	WrapText        bool
}

type borderType string

const (
	Left   borderType = "left"
	Right  borderType = "right"
	Top    borderType = "top"
	Bottom borderType = "bottom"
)

// Border directly maps the border settings of the cells.
type Border struct {
	Type  borderType
	Color string
	Style int
}

func (b *Border) convert() (excelize.Border, error) {
	return excelize.Border{
		Type:  string(b.Type),
		Color: b.Color,
		Style: b.Style,
	}, nil
}

// Font directly maps the font settings of the fonts.
type Font struct {
	Bold         bool
	Italic       bool
	Underline    string
	Family       string
	Size         float64
	Strike       bool
	Color        string
	ColorIndexed int
	ColorTheme   *int
	ColorTint    float64
	VertAlign    string
}

func (f *Font) convert() (*excelize.Font, error) {
	return &excelize.Font{
		Bold:         f.Bold,
		Italic:       f.Italic,
		Underline:    f.Underline,
		Family:       f.Family,
		Size:         f.Size,
		Strike:       f.Strike,
		Color:        f.Color,
		ColorIndexed: f.ColorIndexed,
		ColorTheme:   f.ColorTheme,
		ColorTint:    f.ColorTint,
		VertAlign:    f.VertAlign,
	}, nil
}

type FillType string

const (
	// Pattern pattern fill type
	Pattern FillType = "pattern"
	// Gradient gradient fill type
	Gradient FillType = "gradient"
)

// Fill directly maps the fill settings of the cells.
type Fill struct {
	Type    FillType
	Pattern int
	Color   []string
	Shading int
}

func (f *Fill) convert() (excelize.Fill, error) {
	return excelize.Fill{
		Type:    string(f.Type),
		Pattern: f.Pattern,
		Color:   f.Color,
		Shading: f.Shading,
	}, nil
}

// Protection directly maps the protection settings of the cells.
type Protection struct {
	Hidden bool
	Locked bool
}

func (p *Protection) convert() (*excelize.Protection, error) {
	return &excelize.Protection{
		Hidden: p.Hidden,
		Locked: p.Locked,
	}, nil
}
