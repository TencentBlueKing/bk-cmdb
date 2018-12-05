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
	"configcenter/src/common/blog"
	"fmt"
	"reflect"

	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
)

// Page the common page definition
type Page struct {
	Start int64    `json:"start"`
	Limit int64    `json:"limit"`
	Sort  []string `json:"sort"`
}

// QueryCondition the common query condition definition
type QueryCondition struct {
	Fields    string        `json:"fields"`
	SplitPage Page          `json:"page"`
	Condition mapstr.MapStr `json:"condition"`
}

// QueryResult common query result
type QueryResult struct {
	Count int64           `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}

type SearchLimit struct {
	Offset int64 `json:"start"`
	Limit  int64 `json:"limit"`
}

type SearchSort struct {
	IsDsc bool   `json:"is_dsc"`
	Field string `json:"field"`
}

type SearchInput struct {
	Fields    []string                   `json:"fields"`
	Limit     *SearchLimit               `json:"limit"`
	SortArr   []SearchSort               `json:"sort"`
	Condition []SearchInputConditionItem `json:"condition"`
}

type SearchInputConditionItem struct {
	Fields   string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

func (s *SearchInput) ToSearchCondition() *DBSearchCondition {
	return s.toMongoSearchCondition()
}

func (s *SearchInput) toMongoSearchCondition() *DBSearchCondition {
	sc := &DBSearchCondition{}
	cond := mapstr.New()
	for _, item := range s.Condition {
		switch item.Fields {
		case condition.BKDBEQ, condition.BKDBGT, condition.BKDBGTE,
			condition.BKDBIN, condition.BKDBNIN, condition.BKDBLIKE,
			condition.BKDBLT, condition.BKDBLTE, condition.BKDBNE, condition.BKDBOR:
			cond.Set(item.Fields, getMongoCond(item.Value))
		default:
			cond.Set(item.Fields, mapstr.MapStr{item.Operator: getMongoCond(item.Value)})

		}
	}
	sc.Condition = cond

	for _, sort := range s.SortArr {
		if sort.IsDsc {
			sc.SortArr = append(sc.SortArr, fmt.Sprintf("%s:-1", sort.Field))
		} else {
			sc.SortArr = append(sc.SortArr, fmt.Sprintf("%s:1", sort.Field))
		}
	}

	sc.Fields = s.Fields
	sc.Limit = s.Limit
	return sc
}

type DBSearchCondition struct {
	Fields    []string      `json:"fields"`
	Limit     *SearchLimit  `json:"sort"`
	SortArr   []string      `json:"sort"`
	Condition mapstr.MapStr `json:"condition"`
}

func getMongoCond(val interface{}) interface{} {

	if nil == val {
		return nil
	}
	var item SearchInputConditionItem
	var ok bool

	valType := reflect.TypeOf(val)

	switch valType.Kind() {
	default:

		return val
	case reflect.Struct:
		item, ok = val.(SearchInputConditionItem)
		if !ok {
			return val
		}

	case reflect.Map:
		tmpMap, err := mapstr.NewFromInterface(val)
		if nil != err {
			blog.Warnf("getMongoCond field error %s", err.Error())
			return val
		}

		field, err := tmpMap.String("field")
		if err != nil {
			blog.Warnf("getMongoCond field error %s", err.Error())
			return val
		}
		op, err := tmpMap.String("operator")
		if err != nil {
			blog.Warnf("getMongoCond operator error %s", err.Error())
			return val
		}
		itemVal, _ := tmpMap.Get("value")
		item.Fields = field
		item.Operator = op
		item.Value = itemVal

	}
	cond := mapstr.New()
	switch item.Fields {
	case condition.BKDBEQ, condition.BKDBGT, condition.BKDBGTE,
		condition.BKDBIN, condition.BKDBNIN, condition.BKDBLIKE,
		condition.BKDBLT, condition.BKDBLTE, condition.BKDBNE, condition.BKDBOR:
		cond.Set(item.Fields, getMongoCond(item.Value))
	default:
		cond.Set(item.Fields, mapstr.MapStr{item.Operator: getMongoCond(item.Value)})

	}
	return cond

}
