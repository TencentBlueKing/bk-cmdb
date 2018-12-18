/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"strings"

	"configcenter/src/common/mapstr"
)

// SearchLimit sub condition
type SearchLimit struct {
	Offset int64 `json:"start"`
	Limit  int64 `json:"limit"`
}

// SearchSort sub condition
type SearchSort struct {
	IsDsc bool   `json:"is_dsc"`
	Field string `json:"field"`
}

// QueryCondition the common query condition definition
type QueryCondition struct {
	Fields    []string      `json:"fields"`
	Limit     SearchLimit   `json:"limit"`
	SortArr   []SearchSort  `json:"sort"`
	Condition mapstr.MapStr `json:"condition"`
}

// QueryResult common query result
type QueryResult struct {
	Count uint64          `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}

type QueryConditionResult ResponseInstData

// SearchSortParse SearchSort parse interface
type SearchSortParse interface {
	String(sort string) []SearchSort
	ToMongo(ssArr []SearchSort) string
}

// searchSortParse SearchSort parse struct
type searchSortParse struct {
}

func NewSearchSortParse() SearchSortParse {
	return &searchSortParse{}
}

//  str convert string srot to cc SearchSort struct array
func (ss *searchSortParse) String(sort string) []SearchSort {
	if sort == "" {
		return nil
	}
	sortArr := strings.Split(sort, ",")
	var ssArr []SearchSort
	for _, sortItem := range sortArr {
		sortItemArr := strings.Split(sortItem, ":")
		ssInst := SearchSort{
			Field: sortItemArr[0],
		}
		if len(sortItemArr) > 1 && strings.TrimSpace(sortItemArr[1]) == "-1" {
			ssInst.IsDsc = true

		}
		ssArr = append(ssArr, ssInst)
	}
	return ssArr
}

// searchSortParse cc SearchSort struct to mongodb sort filed
func (ss *searchSortParse) ToMongo(ssArr []SearchSort) string {
	var orderByArr []string
	for _, item := range ssArr {
		if item.IsDsc {
			orderByArr = append(orderByArr, item.Field+":-1")
		} else {
			orderByArr = append(orderByArr, item.Field+":1")
		}
	}
	return strings.Join(orderByArr, ",")
}
