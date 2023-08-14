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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

// Excel equivalent to an excel
type Excel struct {
	sync.RWMutex
	filePath        string
	file            *excelize.File
	writers         map[string]*excelize.StreamWriter
	delDefaultSheet bool
}

type OperatorFunc func(excel *Excel) error

// NewExcel create an excel
func NewExcel(opts ...OperatorFunc) (*Excel, error) {
	excel := &Excel{
		writers: make(map[string]*excelize.StreamWriter),
	}
	for _, opt := range opts {
		if err := opt(excel); err != nil {
			return nil, err
		}
	}

	return excel, nil
}

const excelSuffix = ".xlsx"

// FilePath set Excel file path
func FilePath(filePath string) OperatorFunc {
	return func(excel *Excel) error {
		excel.Lock()
		defer excel.Unlock()
		if !strings.HasSuffix(filePath, excelSuffix) {
			filePath = filePath + excelSuffix
		}

		excel.filePath = filePath
		return nil
	}
}

func isFileExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// OpenOrCreate open if it exists, and create a file if it does not exist
func OpenOrCreate() OperatorFunc {
	return func(excel *Excel) error {
		excel.Lock()
		defer excel.Unlock()
		if excel.filePath == "" {
			return errors.New("excel filePath can not be empty")
		}

		dirPath := filepath.Dir(excel.filePath)
		if _, err := os.Stat(dirPath); err != nil {
			if err := os.MkdirAll(dirPath, os.ModeDir|os.ModePerm); err != nil {
				return err
			}
		}

		if !isFileExist(excel.filePath) {
			excel.file = excelize.NewFile()
			return nil
		}

		var err error
		excel.file, err = excelize.OpenFile(excel.filePath)
		if err != nil {
			return err
		}

		return nil
	}
}

// DelDefaultSheet delete default sheet
func DelDefaultSheet() OperatorFunc {
	return func(excel *Excel) error {
		excel.Lock()
		defer excel.Unlock()

		excel.delDefaultSheet = true

		return nil
	}
}

// CreateSheet create a new sheet
func (excel *Excel) CreateSheet(sheet string) error {
	excel.Lock()
	defer excel.Unlock()
	if _, err := excel.file.NewSheet(sheet); err != nil {
		return err
	}

	return nil
}

// DeleteSheet delete sheet
func (excel *Excel) DeleteSheet(sheet string) error {
	excel.Lock()
	defer excel.Unlock()

	return excel.deleteSheet(sheet)
}

func (excel *Excel) deleteSheet(sheet string) error {
	return excel.file.DeleteSheet(sheet)
}

// SetAllColsWidth set the same width for all columns
func (excel *Excel) SetAllColsWidth(sheet string, width float64) error {
	return excel.SetColWidth(sheet, excelize.MinColumns, excelize.MaxColumns, width)
}

// SetColWidth provides a function to set the width of a single column or
// multiple columns.
func (excel *Excel) SetColWidth(sheet string, startCol, endCol int, width float64) error {
	excel.Lock()
	defer excel.Unlock()

	var err error
	if excel.writers[sheet] == nil {
		excel.writers[sheet], err = excel.file.NewStreamWriter(sheet)
		if err != nil {
			return err
		}
	}

	if err := excel.writers[sheet].SetColWidth(startCol, endCol, width); err != nil {
		return err
	}

	return nil
}

const (
	rowStartIdx = 1
	colStartIdx = 1
)

// StreamingWrite streaming write to file
func (excel *Excel) StreamingWrite(sheet string, startIdx int, data [][]Cell) error {
	excel.Lock()
	defer excel.Unlock()
	if excel.file == nil {
		return fmt.Errorf("excel file has not been created yet")
	}

	var err error
	if excel.writers[sheet] == nil {
		excel.writers[sheet], err = excel.file.NewStreamWriter(sheet)
		if err != nil {
			return err
		}
	}

	startIdx++
	if startIdx < rowStartIdx {
		return fmt.Errorf("row start index is invalid, val: %d", startIdx)
	}

	for i := 0; i < len(data); i++ {
		firstCell, err := excelize.CoordinatesToCellName(colStartIdx, startIdx+i)
		if err != nil {
			return err
		}

		cells := make([]interface{}, len(data[i]))
		for idx, val := range data[i] {
			cells[idx] = val.transfer()
		}

		if err := excel.writers[sheet].SetRow(firstCell, cells); err != nil {
			return err
		}
	}

	return nil
}

// StreamingRead streaming read from file
func (excel *Excel) StreamingRead(sheet string) ([][]interface{}, error) {
	excel.RLock()
	defer excel.RUnlock()
	if excel.file == nil {
		return nil, fmt.Errorf("excel file has not been created yet")
	}

	rows, err := excel.file.Rows(sheet)
	if err != nil {
		return nil, err
	}

	var result [][]interface{}
	for rows.Next() {
		var rowData []interface{}
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		for _, cell := range cols {
			rowData = append(rowData, cell)
		}

		result = append(result, rowData)
	}

	// close the stream
	if err = rows.Close(); err != nil {
		return nil, err
	}

	return result, nil
}

// MergeCell provides a function to merge cells by a given range reference for the StreamWriter
func (excel *Excel) MergeCell(sheet string, hCell, vCell string) error {
	excel.Lock()
	defer excel.Unlock()

	return excel.mergeCell(sheet, hCell, vCell)
}

func (excel *Excel) mergeCell(sheet string, hCell, vCell string) error {
	if excel.file == nil {
		return fmt.Errorf("excel file has not been created yet")
	}

	var err error
	if excel.writers[sheet] == nil {
		excel.writers[sheet], err = excel.file.NewStreamWriter(sheet)
		if err != nil {
			return err
		}
	}

	return excel.writers[sheet].MergeCell(hCell, vCell)
}

// Flush file
func (excel *Excel) Flush(sheet string) error {
	excel.Lock()
	defer excel.Unlock()
	if excel.file == nil {
		return fmt.Errorf("excel file has not been created yet")
	}

	if excel.writers[sheet] == nil {
		return nil
	}

	if err := excel.writers[sheet].Flush(); err != nil {
		return err
	}
	delete(excel.writers, sheet)

	return nil
}

// Save file
func (excel *Excel) Save() error {
	excel.Lock()
	defer excel.Unlock()

	return excel.save()
}

func (excel *Excel) save() error {
	if excel.file == nil {
		return fmt.Errorf("excel file has not been created yet")
	}

	return excel.file.SaveAs(excel.filePath)
}

// Close Excel file
func (excel *Excel) Close() error {
	excel.Lock()
	defer excel.Unlock()

	if excel.file == nil {
		return fmt.Errorf("excel file has not been created yet")
	}

	if excel.delDefaultSheet {
		if err := excel.deleteSheet(defaultSheet); err != nil {
			return err
		}

		if err := excel.save(); err != nil {
			return err
		}
	}

	if err := excel.file.Close(); err != nil {
		return err
	}

	return nil
}

const defaultSheet = "Sheet1"

// Clean delete temporary file
func (excel *Excel) Clean() error {
	excel.Lock()
	defer excel.Unlock()

	if err := os.Remove(excel.filePath); err != nil {
		return err
	}

	return nil
}

// AddValidation add validation
func (excel *Excel) AddValidation(sheet string, param *ValidationParam) error {
	excel.Lock()
	defer excel.Unlock()

	validation, err := newValidation(param)
	if err != nil {
		return err
	}

	if err := excel.file.AddDataValidation(sheet, validation); err != nil {
		return err
	}

	if err := excel.save(); err != nil {
		return err
	}

	return nil
}

// NewStyle new style
func (excel *Excel) NewStyle(style *Style) (int, error) {
	excelStyle, err := style.convert()
	if err != nil {
		return 0, err
	}

	return excel.file.NewStyle(excelStyle)
}

const singleCellLen = 1

// MergeSameRowCell merge same row cell
func (excel *Excel) MergeSameRowCell(sheet string, colIdx, rowIdx, length int) error {

	if length == singleCellLen {
		return nil
	}

	hCell, err := GetCellIdx(colIdx, rowIdx)
	if err != nil {
		return err
	}

	vCell, err := GetCellIdx(colIdx+length-1, rowIdx)
	if err != nil {
		return err
	}

	excel.Lock()
	defer excel.Unlock()

	if err := excel.mergeCell(sheet, hCell, vCell); err != nil {
		return err
	}

	return nil
}

// MergeSameColCell merge same column cell
func (excel *Excel) MergeSameColCell(sheet string, colIdx, rowIdx, height int) error {
	if height == singleCellLen {
		return nil
	}

	hCell, err := GetCellIdx(colIdx, rowIdx)
	if err != nil {
		return err
	}

	vCell, err := GetCellIdx(colIdx, rowIdx+height-1)
	if err != nil {
		return err
	}

	excel.Lock()
	defer excel.Unlock()

	if err := excel.mergeCell(sheet, hCell, vCell); err != nil {
		return err
	}

	return nil
}
