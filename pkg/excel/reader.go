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
	"sync"

	"github.com/xuri/excelize/v2"
)

// Reader read data in the way of io stream
type Reader struct {
	sync.RWMutex
	rows   *excelize.Rows
	curIdx int
}

// Next will return true if the next row element is found
func (r *Reader) Next() bool {
	r.RLock()
	defer r.RUnlock()

	r.curIdx++
	return r.rows.Next()
}

// CurRow return the current row's column values. This fetches the worksheet
// data as a stream, returns each cell in a row as is, and will not skip empty
// rows in the tail of the worksheet.
func (r *Reader) CurRow() ([]string, error) {
	r.Lock()
	defer r.Unlock()

	row, err := r.rows.Columns()
	if err != nil {
		return nil, err
	}

	return row, nil
}

// Close reader
func (r *Reader) Close() error {
	r.Lock()
	defer r.Unlock()

	return r.rows.Close()
}

// GetCurIdx get current index
func (r *Reader) GetCurIdx() int {
	r.RLock()
	defer r.RUnlock()

	return r.curIdx
}
