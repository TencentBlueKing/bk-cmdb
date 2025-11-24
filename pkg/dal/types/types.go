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

	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
	"github.com/TencentBlueKing/bk-cmdb/pkg/structs"
)

// ListOption defines options to list resources.
type ListOption struct {
	Fields []string
	Filter filter.RuleFactory
	Page   *BasePage
	Option *ListOption
}

// DynamicListDetails defines the results of listing resources.
type DynamicListDetails struct {
	Count   uint64         `json:"count,omitempty"`
	Details *structs.Slice `json:"details,omitempty"`
}

// ListDetails defines the results of listing resources.
type ListDetails[T any] struct {
	Count   uint64 `json:"count,omitempty"`
	Details []T    `json:"details,omitempty"`
}

// Validate list option.
func (opt ListOption) Validate(eo *filter.ExprOption, po *PageOption) error {
	err := opt.validatePage(po)
	if err != nil {
		return err
	}

	if eo == nil {
		return errors.New("filter expr option is required")
	}

	if err = opt.validateFilter(eo); err != nil {
		return err
	}

	if err := opt.ValidateFields(eo); err != nil {
		return err
	}

	return nil
}

func (opt ListOption) validateFilter(eo *filter.ExprOption) error {
	if opt.Filter == nil {
		return errors.New("filter expr is required")
	}

	if err := opt.Filter.Validate(eo); err != nil {
		return err
	}

	return nil
}

func (opt ListOption) validatePage(po *PageOption) error {
	if opt.Page == nil {
		return errors.New("page is required")
	}

	if po == nil {
		return errors.New("page option is required")
	}
	if err := opt.Page.Validate(po); err != nil {
		return err
	}
	return nil
}

// ValidateFields check fields correctness, if fields is empty or rule fields is nil, skip check.
func (opt ListOption) ValidateFields(eo *filter.ExprOption) error {
	if eo.RuleFields == nil {
		// skip check
		return nil
	}
	if len(opt.Fields) == 0 {
		return nil
	}

	for _, field := range opt.Fields {
		if _, ok := eo.RuleFields[field]; !ok {
			return fmt.Errorf("unknown field: %s", field)
		}
	}

	return nil

}
