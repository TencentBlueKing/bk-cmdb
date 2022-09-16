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
	"errors"
	"fmt"
	"strconv"

	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
)

const (
	// PageName TODO
	PageName = "page"
	// PageSort TODO
	PageSort = "sort"
	// PageStart TODO
	PageStart = "start"
	// PageLimit TODO
	PageLimit = "limit"
	// DBFields TODO
	DBFields = "fields"
	// DBQueryCondition TODO
	DBQueryCondition = "condition"
)

// BasePage for paging query
type BasePage struct {
	Sort        string `json:"sort,omitempty" mapstructure:"sort"`
	Limit       int    `json:"limit,omitempty" mapstructure:"limit"`
	Start       int    `json:"start" mapstructure:"start"`
	EnableCount bool   `json:"enable_count,omitempty" mapstructure:"enable_count,omitempty"`
}

// Validate TODO
func (page BasePage) Validate(allowNoLimit bool) (string, error) {
	// 此场景下如果仅仅是获取查询对象的数量，page的其余参数只能是初始化值
	if page.EnableCount {
		if page.Start > 0 || page.Limit > 0 || page.Sort != "" {
			return "page", errors.New("params page can not be set")
		}
		return "", nil
	}

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
func (page BasePage) ValidateLimit(maxLimit int) error {
	if page.Limit == 0 {
		return errors.New("page limit must not be zero")
	}

	if maxLimit > common.BKMaxPageSize {
		return fmt.Errorf("exceed system max page size: %d", common.BKMaxPageSize)
	}

	if page.Limit > maxLimit {
		return fmt.Errorf("exceed business max page size: %d", maxLimit)
	}

	return nil
}

// ValidateWithEnableCount validate if page has only one of enable count and other param, and if limit is set and valid
func (page BasePage) ValidateWithEnableCount(allowNoLimit bool, maxLimit ...int) ccErr.RawErrorInfo {
	if page.EnableCount {
		if page.Start != 0 || page.Limit != 0 || page.Sort != "" {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{"page.enable_count"},
			}
		}
		return ccErr.RawErrorInfo{}
	}

	if page.Limit == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"page.limit"},
		}
	}

	limit := common.BKMaxPageSize
	if len(maxLimit) > 0 {
		limit = maxLimit[0]
	}

	if page.Limit > limit {
		if allowNoLimit || page.Limit != common.BKNoLimit {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommPageLimitIsExceeded,
			}
		}
	}
	return ccErr.RawErrorInfo{}
}

// ParsePage TODO
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

// ToSearchSort TODO
func (page BasePage) ToSearchSort() []SearchSort {
	return NewSearchSortParse().String(page.Sort).ToSearchSortArr()
}
