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

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
)

// Deprecated: SearchLimit sub condition
type SearchLimit struct {
	Offset int64 `json:"start" field:"start"`
	Limit  int64 `json:"limit" field:"limit"`
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

// IsIllegal  limit is illegal, if limit = 0; change to default page size
func (qc *QueryCondition) IsIllegal() bool {
	if qc.Limit.Limit == 0 {
		qc.Limit.Limit = common.BKDefaultLimit
	}
	if qc.Limit.Limit > common.BKMaxPageSize && qc.Limit.Limit != common.BKNoLimit {
		return true
	}
	return false
}

// QueryResult common query result
type QueryResult struct {
	Count uint64          `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}

type QueryConditionResult ResponseInstData

// SearchSortParse SearchSort parse interface
type SearchSortParse interface {
	String(sort string) SearchSortParse
	Field(field string, isDesc bool) SearchSortParse
	Set(ssArr []SearchSort) SearchSortParse
	ToMongo() string
	ToSearchSortArr() []SearchSort
}

// searchSortParse SearchSort parse struct
type searchSortParse struct {
	data []SearchSort
}

func NewSearchSortParse() SearchSortParse {
	return &searchSortParse{}
}

//  String convert string sort to cc SearchSort struct array
func (ss *searchSortParse) String(sort string) SearchSortParse {
	if sort == "" {
		return ss
	}
	sortArr := strings.Split(sort, ",")
	for _, sortItem := range sortArr {
		sortItemArr := strings.Split(sortItem, ":")
		ssInst := SearchSort{
			Field: sortItemArr[0],
		}
		if len(sortItemArr) > 1 && strings.TrimSpace(sortItemArr[1]) == "-1" {
			ssInst.IsDsc = true

		}
		ss.data = append(ss.data, ssInst)
	}
	return ss
}

//  Field   cc SearchSort struct array
func (ss *searchSortParse) Field(field string, isDesc bool) SearchSortParse {

	ssInst := SearchSort{
		Field: field,
		IsDsc: isDesc,
	}
	ss.data = append(ss.data, ssInst)
	return ss
}

func (ss *searchSortParse) Set(ssArr []SearchSort) SearchSortParse {
	ss.data = append(ss.data, ssArr...)
	return ss
}

// ToSearchSortArr cc SearchSort struct to mongodb sort filed
func (ss *searchSortParse) ToSearchSortArr() []SearchSort {
	return ss.data
}

// searchSortParse cc SearchSort struct to mongodb sort filed
func (ss *searchSortParse) ToMongo() string {
	var orderByArr []string
	for _, item := range ss.data {
		if item.IsDsc {
			orderByArr = append(orderByArr, item.Field+":-1")
		} else {
			orderByArr = append(orderByArr, item.Field+":1")
		}
	}
	return strings.Join(orderByArr, ",")
}

// IsIllegal  limit is illegal
func (page SearchLimit) IsIllegal() bool {
	if page.Limit > common.BKMaxPageSize && page.Limit != common.BKNoLimit ||
		page.Limit == 0 {
		return true
	}
	return false
}
