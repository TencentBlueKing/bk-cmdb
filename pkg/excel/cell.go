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

package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// Cell can be used directly in StreamWriter.SetRow to specify a style and a value.
type Cell struct {
	StyleID int
	Value   interface{}
}

func (c *Cell) transfer() *excelize.Cell {
	return &excelize.Cell{
		Value:   c.Value,
		StyleID: c.StyleID,
	}
}

// GetCellIdx get cell index
func GetCellIdx(col int, row int) (string, error) {
	// 由于第三方库的行和列不是从0开始，所以这里加上开始数，使调用者可以按照从0开始进行计数
	return excelize.CoordinatesToCellName(col+colStartIdx, row+rowStartIdx)
}

// GetSingleColSqref get single column sqref
// Example: GetSingleColSqref(0, 1, 2) // return A1:A2
func GetSingleColSqref(col, startRow, endRow int) (string, error) {
	colNum, err := ColumnNumberToName(col)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%d:%s%d", colNum, startRow, colNum, endRow), nil
}

// ColumnNumberToName provides a function to convert the integer to Excel
// sheet column title.
// Example:	ColumnNumberToName(0) // returns "A", nil
func ColumnNumberToName(col int) (string, error) {
	// 由于第三方库的列是从1开始，所以这里进行了+1操作，使调用者可以按照从0开始进行计数
	return excelize.ColumnNumberToName(col + 1)
}

// GetTotalRows get total rows
func GetTotalRows() int {
	return excelize.TotalRows
}

// CellNameToCoordinates converts alphanumeric cell name to [X, Y] coordinates
// or returns an error.
//
// Example:
//
//	CellNameToCoordinates("A1") // returns 1, 1, nil
//	CellNameToCoordinates("Z3") // returns 26, 3, nil
func CellNameToCoordinates(cell string) (int, int, error) {
	return excelize.CellNameToCoordinates(cell)
}

// CellMergeMsg cell merge message
type CellMergeMsg struct {
	start string
	end   string
}

// GetStartAxis returns the top left cell reference of merged range, for
// example: "C2".
func (c *CellMergeMsg) GetStartAxis() string {
	return c.start
}

// GetEndAxis returns the bottom right cell reference of merged range, for
// example: "D4".
func (c *CellMergeMsg) GetEndAxis() string {
	return c.end
}
