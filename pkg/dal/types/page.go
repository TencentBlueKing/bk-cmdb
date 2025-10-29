/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

// Package types ...
package types

import (
	"errors"
	"fmt"
)

// PageWithoutSort page without sort.
type PageWithoutSort struct {
	// Count describe if this query only return the total request
	// count of the resources.
	// If true, then the request will only return the total count
	// without the resource's detail infos. and start, limit must
	// be 0.
	Count bool `json:"count"`
	// Start is the start position of the queried resource's page.
	// Note:
	// 1. Start only works when the Count = false.
	// 2. Start's minimum value is 0, not 1.
	// 3. if PageOption.EnableUnlimitedLimit = true, then Start = 0
	//   and Limit = 0 means query all the resources at once.
	Start uint32 `json:"start"`
	// Limit is the total returned resources at once query.
	// Limit only works when the Count = false.
	Limit uint `json:"limit"`
}

const (
	// DefaultMaxPageLimit is the default value of the max page limitation.
	DefaultMaxPageLimit = uint64(500)
	// AggregationQueryMaxPageLimit 聚合查询最大数量限制
	AggregationQueryMaxPageLimit = uint64(50)
)

// NewDefaultPage define default base page.
func NewDefaultPage() *BasePage {
	return &BasePage{
		Count: false,
		Start: 0,
		Limit: DefaultMaxPageLimit,
	}
}

// NewCountPage define count page.
func NewCountPage() *BasePage {
	return &BasePage{
		Count: true,
	}
}

// NewDefaultPageOption is the default BasePage's option.
func NewDefaultPageOption() *PageOption {
	return &PageOption{
		EnableUnlimitedLimit: false,
		MaxLimit:             DefaultMaxPageLimit,
		DisabledSort:         false,
	}
}

// PageOption defines the options to validate the
// BasePage's configuration.
type PageOption struct {
	// EnableUnlimitedLimit allows user to query resources with unlimited
	// limitation. if true, then the 'Limit' option will not be checked.
	EnableUnlimitedLimit bool `json:"enable_unlimited_limit"`
	// MaxLimit defines max limit value of a page.
	MaxLimit uint64 `json:"max_limit"`
	// DisableSort defines the sort field is not allowed to be defined by the user.
	// then system defined sort field is used.
	// Note: this option does not work when use the page to generate SQL expression,
	// which means call the method of BasePage's SQLExprAndValue().
	DisabledSort bool `json:"disabled_sort"`
}

// Order is the direction when do sort operation.
type Order string

const (
	// Ascending sort data with ascending direction
	// this is the default sort direction.
	Ascending Order = "ASC"
	// Descending sort data with descending direction
	Descending Order = "DESC"
)

// Validate the sort direction is valid or not
func (sd Order) Validate() error {
	if len(sd) == 0 {
		return nil
	}

	switch sd {
	case Ascending:
	case Descending:
	default:
		return fmt.Errorf("unsupported sort direction: %s", sd)
	}

	return nil
}

// Order returns the sort direction, if not set, use
// ascending as the default direction.
func (sd Order) Order() Order {
	switch sd {
	case Ascending:
		return Ascending
	case Descending:
		return Descending
	default:
		// set Ascending as the default sort direction.
		return Ascending
	}
}

// BasePage define the basic page limitation to query resources.
type BasePage struct {
	// Count describe if this query only return the total request
	// count of the resources.
	// If true, then the request will only return the total count
	// without the resource's detail infos. and start, limit must
	// be 0.
	Count bool `json:"count"`
	// Start is the start position of the queried resource's page.
	// Note:
	// 1. Start only works when the Count = false.
	// 2. Start's minimum value is 0, not 1.
	// 3. if PageOption.EnableUnlimitedLimit = true, then Start = 0
	//   and Limit = 0 means query all the resources at once.
	Start uint64 `json:"start"`
	// Limit is the total returned resources at once query.
	// Limit only works when the Count = false.
	Limit uint64 `json:"limit"`
	// Sort defines use which field to sort the queried resources.
	// only 'one' field is supported to do sort.
	// Sort only works when the Count = false.
	Sort string `json:"sort"`
	// Order is the direction when do sort operation.
	// it works only when the Sort is set.
	Order Order `json:"order"`
}

// Validate the base page's options.
// if the page option is not set, use the default configuration.
func (bp BasePage) Validate(opt ...*PageOption) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("invalid page option: %w", err)
		}
	}()

	if len(opt) >= 2 {
		return errors.New("at most one page options is allows")
	}

	if bp.Count {
		if bp.Start > 0 {
			return errors.New("count is enabled, page.start should be 0")
		}

		if bp.Limit > 0 {
			return errors.New("count is enabled, page.limit should be 0")
		}

		if len(bp.Sort) > 0 {
			return errors.New("count is enabled, page.sort should be null")
		}

		if len(bp.Order) > 0 {
			return errors.New("count is enabled, page.order should be empty")
		}

		return nil
	}

	maxLimit := DefaultMaxPageLimit
	enableUnlimited := false
	if len(opt) != 0 {
		// option is configured, validate it
		one := opt[0]
		if one.MaxLimit > 0 {
			maxLimit = one.MaxLimit
		}

		enableUnlimited = one.EnableUnlimitedLimit

		if one.DisabledSort {
			if len(bp.Sort) > 0 {
				return errors.New("page.sort is not allowed")
			}

			if len(bp.Order) > 0 {
				return errors.New("invalid page.order, page.order is not allowed")
			}
		}
	}

	if !enableUnlimited {
		// if the user is not allowed to query with unlimited limit, then
		// 1. limit should >=1
		// 2. validate whether the limit is larger than the max limit value
		if bp.Limit < 1 {
			return errors.New("page.limit value should >= 1")
		}

		if bp.Limit > maxLimit {
			return fmt.Errorf("invalid page.limit max value: %d", maxLimit)
		}
	}

	// if direction is set, then validate it.
	if len(bp.Order) != 0 {
		if err := bp.Order.Validate(); err != nil {
			return err
		}
	}

	return nil
}
