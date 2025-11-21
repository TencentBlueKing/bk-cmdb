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

package base

import (
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/types"
	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
)

// Config defines dao operation config
type Config struct {
	PageOption          *types.PageOption
	ExprOptionFunctions []filter.ExprOptionFunc
}

// GetConfig return page option and expr config modified by options
func GetConfig(ruleFields map[string]filter.FieldType, opts []Option) (*filter.ExprOption, *types.PageOption) {
	c := Config{
		// provide a default option to avoid nil pointer check
		PageOption: types.NewDefaultPageOption(),
	}
	for _, opt := range opts {
		opt(&c)
	}
	filterOpt := filter.NewExprOption(filter.RuleFields(ruleFields))
	for _, expOpt := range c.ExprOptionFunctions {
		expOpt(filterOpt)
	}
	return filterOpt, c.PageOption
}

// Option configure page and expr option, both page option and expr option are guaranteed not nil
type Option func(*Config)
