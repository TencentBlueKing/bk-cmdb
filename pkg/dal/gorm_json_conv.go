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

package dal

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm/clause"

	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
	"github.com/TencentBlueKing/bk-cmdb/pkg/util"
)

// jsonRuleToClauseExpr convert JSON rule to gorm clause.
// Ref:
// postgresql  https://www.postgresql.org/docs/current/functions-json.html
// mysql: https://dev.mysql.com/doc/refman/8.0/en/json-search-functions.html
func jsonRuleToClauseExpr(rule *filter.AtomRule) (clause.Expression, error) {
	parts := strings.Split(rule.Field, filter.JSONFieldSeparator)
	col := parts[0]
	keys := parts[1:]
	switch filter.OpType(rule.Op) {
	case filter.JSONEqual:
		return datatypes.JSONQuery(col).Equals(rule.Value, keys...), nil
	case filter.JSONNotEqual:
		eqExpr := datatypes.JSONQuery(col).Equals(rule.Value, keys...)
		return clause.NotConditions{Exprs: []clause.Expression{eqExpr}}, nil
	case filter.JSONContains:
		return datatypes.JSONArrayQuery(col).Contains(rule.Value, keys...), nil
	case filter.JSONHasKey:
		return buildJSONHasKey(col, keys, rule.Value)
	case filter.JSONNotHasKey:
		containsExpr, err := buildJSONHasKey(col, keys, rule.Value)
		if err != nil {
			return nil, fmt.Errorf("fail to build not json has key: %w", err)
		}
		return clause.NotConditions{Exprs: []clause.Expression{containsExpr}}, nil
	default:
		return nil, fmt.Errorf("json op is not supported: %s", rule.Op)
	}
}

func buildJSONHasKey(col string, keys []string, value any) (clause.Expression, error) {
	path, ok := util.GetString(value)
	if !ok {
		return nil, errors.New("json has key operator's path value should be an string")
	}

	if len(path) == 0 {
		return nil, errors.New("json has key operator's path value can not be a empty string")
	}

	mergedKeys := append(slices.Clip(keys), path)
	// HasKey方法会转换为pg的`?` 操作符，最后应为可选的含.分隔的string形式待查找 Path；其他的是数组，用于拼接列的JSON路径
	// 结果为： `col->'keys[0]'->'keys[1]' ? 'path'`
	return datatypes.JSONQuery(col).HasKey(mergedKeys...), nil
}
