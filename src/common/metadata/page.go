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

package metadata

import (
	"fmt"
	"strconv"

	"configcenter/src/common"
)

const (
	PageName         = "page"
	PageSort         = "sort"
	PageStart        = "start"
	DBFields         = "fields"
	DBQueryCondition = "condition"
)

// BasePage for paging query
type BasePage struct {
	Sort  string `json:"sort,omitempty" mapstructure:"sort"`
	Limit int    `json:"limit,omitempty" mapstructure:"limit"`
	Start int    `json:"start" mapstructure:"start"`
}

func (page BasePage) Validate(allowNoLimit bool) (string, error) {
	if page.Limit > common.BKMaxPageSize {
		if page.Limit != common.BKNoLimit || allowNoLimit != true {
			return "limit", fmt.Errorf("exceed max page size: %d", common.BKMaxPageSize)
		}
	}
	return "", nil
}

// IsIllegal  limit is illegal
func (page BasePage) IsIllegal() bool {
	if page.Limit > common.BKMaxPageSize && page.Limit != common.BKNoLimit ||
		page.Limit <= 0 {
		return true
	}
	return false
}

// ValidateLimit validates target page limit.
func (page BasePage) ValidateLimit(maxLimit, defaultLimit int) error {
	if page.Limit == 0 {
		page.Limit = defaultLimit
	}

	if page.Limit == 0 {
		return fmt.Errorf("page limit must not be zero")
	}

	if maxLimit > common.BKMaxPageSize {
		return fmt.Errorf("exceed system max page size: %d", common.BKMaxPageSize)
	}

	if page.Limit > maxLimit {
		return fmt.Errorf("exceed business max page size: %d", maxLimit)
	}

	return nil
}

func ParsePage(origin interface{}) BasePage {
	if origin == nil {
		return BasePage{Limit: common.BKNoLimit}
	}
	page, ok := origin.(map[string]interface{})
	if !ok {
		return BasePage{Limit: common.BKNoLimit}
	}
	result := BasePage{}
	if sort, ok := page["sort"]; ok && sort != nil {
		result.Sort = fmt.Sprint(sort)
	}
	if start, ok := page["start"]; ok {
		result.Start, _ = strconv.Atoi(fmt.Sprint(start))
	}
	if limit, ok := page["limit"]; ok {
		result.Limit, _ = strconv.Atoi(fmt.Sprint(limit))
		if result.Limit <= 0 {
			result.Limit = common.BKNoLimit
		}
	}
	return result
}

func (page BasePage) ToSearchSort() []SearchSort {
	return NewSearchSortParse().String(page.Sort).ToSearchSortArr()
}
