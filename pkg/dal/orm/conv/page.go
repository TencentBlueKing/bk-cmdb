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

package conv

import (
	"github.com/samber/lo"
	"gorm.io/gorm/clause"

	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/types"
)

// Page convert page to clause
func Page(page *types.BasePage, option *types.PageOption) []clause.Expression {
	var exps []clause.Expression
	if option != nil && !option.DisabledSort {
		exps = append(exps, OrderBy(page))
	}
	if option != nil && !option.EnableUnlimitedLimit {
		exps = append(exps, Limit(page))
	}
	return exps
}

// Limit convert limit offset
func Limit(page *types.BasePage) clause.Expression {
	limit := clause.Limit{Offset: int(page.Start)}
	if page.Limit > 0 {
		limit.Limit = lo.ToPtr(int(page.Limit))
	}
	return limit
}

// OrderBy convert sort with direction
func OrderBy(page *types.BasePage) clause.Expression {
	// default sort by id
	if len(page.Sort) == 0 {
		return clause.OrderBy{Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: "id"},
			Desc:   false,
		}}}
	}
	orderBy := clause.OrderBy{}
	for i := range page.Sort {
		column := clause.OrderByColumn{Column: clause.Column{Name: page.Sort[i].Field}}
		if page.Sort[i].Order == types.Descending {
			column.Desc = true
		}
		orderBy.Columns = append(orderBy.Columns, column)
	}
	return orderBy
}
