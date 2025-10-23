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

package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
)

type jsonTestModel struct {
	ID           int64                       `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
	Name         string                      `gorm:"primaryKey;column:name" json:"name"`
	JSONCol      datatypes.JSON              `gorm:"column:json_col" json:"json_col"`
	JSONArrayCol datatypes.JSONSlice[string] `gorm:"column:json_array_col" json:"json_array_col"`
}

func (jsonTestModel) TableName() string {
	return "json_test"
}

//nolint:funlen
func Test_atomJSONRuleToClauseExpr(t *testing.T) {
	type args struct {
		rule *filter.AtomRule
	}
	tests := []struct {
		name         string
		args         args
		want         clause.Expression
		wantConvErr  string
		wantQueryErr string
		wantSQL      string
		wantVars     []any
		shouldFound  []jsonTestModel
	}{
		{
			name: "wholeColumnEqualEmptyPath",
			args: args{
				rule: filter.RuleJSONEqual("json_col", `{"a": "b"}`),
			},
			want:         datatypes.JSONQuery("json_col").Equals(`{"a": "b"}`, []string{}...),
			wantQueryErr: "42601",
			// field with no path will be ignored, should use eq op
			wantSQL:     ``,
			wantVars:    nil,
			shouldFound: []jsonTestModel{},
		},

		{
			name: "simplePathEqual",
			args: args{
				rule: filter.RuleJSONEqual("json_col.key_a", `val_a`),
			},
			want:        datatypes.JSONQuery("json_col").Equals("val_a", "key_a"),
			wantSQL:     `json_extract_path_text("json_col"::json,$1) = $2`,
			wantVars:    []any{"key_a", "val_a"},
			shouldFound: []jsonTestModel{singleLevelJson},
		},
		{
			name: "simplePathNotEqual",
			args: args{
				rule: filter.RuleJSONNotEqual("json_col.key_a", `val_b`),
			},
			want:     clause.NotConditions{Exprs: []clause.Expression{datatypes.JSONQuery("json_col").Equals("val_b", "key_a")}},
			wantSQL:  `NOT json_extract_path_text("json_col"::json,$1) = $2`,
			wantVars: []any{"key_a", "val_b"},
			// will found single level json
			shouldFound: []jsonTestModel{singleLevelJson},
		},
		{
			name: "simplePathNotEqual2",
			args: args{
				rule: filter.RuleJSONNotEqual("json_col.key_a", `val_a`),
			},
			want: clause.NotConditions{Exprs: []clause.Expression{datatypes.JSONQuery("json_col").Equals("val_a", "key_a")}},

			wantSQL:  `NOT json_extract_path_text("json_col"::json,$1) = $2`,
			wantVars: []any{"key_a", "val_a"},
			// will found nothing, because another row does not have key_a
			shouldFound: []jsonTestModel{},
		},

		{
			name: "JSONHasKey",
			args: args{
				rule: filter.RuleJSONHasKey("json_col.l1.l2.l3", "l4"),
			},
			want: datatypes.JSONQuery("json_col").HasKey("l1", "l2", "l3", "l4"),

			wantSQL:     `"json_col"::jsonb -> $1 -> $2 -> $3 ? $4`,
			wantVars:    []any{"l1", "l2", "l3", "l4"},
			shouldFound: []jsonTestModel{multiLevelJson},
		},
		{
			name: "ArrayContains",
			args: args{
				rule: filter.RuleJSONContains("json_array_col", "1"),
			},
			want: datatypes.JSONArrayQuery("json_array_col").Contains("1", []string{}...),

			wantSQL:     `"json_array_col" ? $1`,
			wantVars:    []any{"1"},
			shouldFound: []jsonTestModel{singleLevelJson},
		},
	}

	g := prepareTestDB(t, jsonTestModel{}, jsonTestModels)
	assert.NotNil(t, g, "failed to open gorm db")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonRuleToClauseExpr(tt.args.rule)
			if tt.wantConvErr != "" {
				if !assert.ErrorContains(t, err, tt.wantConvErr, "convert to clause failed got unexpected error") {
					return
				}
			} else {
				assert.Nil(t, err, "convert to clause failed")
				return
			}
			if !assert.Equal(t, tt.want, got) {
				return
			}
			assert.Equal(t, tt.want, got)
			s := &gorm.Statement{DB: g.Session(&gorm.Session{DryRun: true})}
			got.Build(s)
			sql := s.SQL.String()
			t.Log("SQL:", sql)
			assert.Equal(t, tt.wantSQL, sql, "sql mismatch")
			assert.Equal(t, tt.wantVars, s.Vars, "bind vars mismatch")
			// try query
			found := make([]jsonTestModel, 0)
			err = g.Model(jsonTestModel{}).Clauses(got).Find(&found).Error
			if (err != nil) != (tt.wantConvErr != "") {
				assert.ErrorContainsf(t, err, tt.wantQueryErr, "query failed: %v", err)
			}
			assert.Equalf(t, tt.shouldFound, found, "found data mismatch")
		})
	}
}

var jsonTestModels = []jsonTestModel{singleLevelJson, multiLevelJson}
var singleLevelJson = jsonTestModel{
	ID:           1,
	Name:         "single-level",
	JSONCol:      datatypes.JSON(`{"key_a": "val_a"}`),
	JSONArrayCol: datatypes.JSONSlice[string]{"1", "2", "3"},
}
var multiLevelJson = jsonTestModel{
	ID:           2,
	Name:         "multi-level",
	JSONCol:      datatypes.JSON(`{"l1": {"l2": {"l3": {"l4": 1}}}}`),
	JSONArrayCol: datatypes.JSONSlice[string]{"4", "5", "6"},
}
