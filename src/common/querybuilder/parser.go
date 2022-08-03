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

package querybuilder

import (
	"encoding/json"
	"fmt"
	"reflect"

	"configcenter/src/common/mapstr"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
)

// RuleGroup TODO
type RuleGroup struct {
	Condition Condition                `json:"condition" field:"condition"`
	Rules     []map[string]interface{} `json:"rules" field:"rules"`
}

// ParseRule TODO
func ParseRule(data map[string]interface{}) (queryFilter Rule, errKey string, err error) {
	if data == nil {
		return nil, "", nil
	}
	if _, ok := data["condition"]; ok {
		ruleGroupData := &RuleGroup{}
		// shouldn't use mapstr here as it doesn't support nest struct
		// TODO: replace it with more efficient way
		if err := mapstr.DecodeFromMapStr(ruleGroupData, data); err != nil {
			return nil, "", fmt.Errorf("decode as combined rule struct failed, err: %+v", err)
		}
		combinedRule := CombinedRule{
			Condition: ruleGroupData.Condition,
			Rules:     make([]Rule, 0),
		}
		for idx, item := range ruleGroupData.Rules {
			qf, errKey, err := ParseRule(item)
			if err != nil {
				return nil, fmt.Sprintf("rules[%d].%s", idx, errKey), err
			}
			if qf != nil {
				combinedRule.Rules = append(combinedRule.Rules, qf)
			}
		}
		queryFilter = combinedRule
	} else if _, ok := data["operator"]; ok {
		rule := AtomRule{}
		if err := mapstr.DecodeFromMapStr(&rule, data); err != nil {
			return nil, "", fmt.Errorf("decode to rule struct failed, err: %+v", err)
		}
		queryFilter = rule
	} else {
		return nil, "", fmt.Errorf("no query filter found")
	}
	return queryFilter, "", nil
}

// ParseRuleFromBytes TODO
func ParseRuleFromBytes(bs []byte) (queryFilter Rule, errKey string, err error) {
	data := make(map[string]interface{})
	if err := json.Unmarshal(bs, &data); err != nil {
		return nil, "", err
	}
	return ParseRule(data)
}

// QueryFilter is aimed at export as a struct member
type QueryFilter struct {
	Rule `json:",inline"`
}

// Validate validates query filter conditions.
func (qf *QueryFilter) Validate(option *RuleOption) (string, error) {
	if qf.Rule == nil {
		return "", nil
	}

	if _, ok := qf.Rule.(CombinedRule); !ok {
		return "", fmt.Errorf("query filter must be combined rules")
	}

	return qf.Rule.Validate(option)
}

// MarshalJSON TODO
func (qf *QueryFilter) MarshalJSON() ([]byte, error) {
	if qf.Rule != nil {
		return json.Marshal(qf.Rule)
	}
	return make([]byte, 0), nil
}

// UnmarshalJSON TODO
func (qf *QueryFilter) UnmarshalJSON(raw []byte) error {
	rule, errKey, err := ParseRuleFromBytes(raw)
	if err != nil {
		return fmt.Errorf("UnmarshalJSON failed, key: %s, err: %+v", errKey, err)
	}
	qf.Rule = rule
	return nil
}

// MarshalBSON marshal query filter into bson value
func (qf *QueryFilter) MarshalBSON() ([]byte, error) {
	if qf.Rule != nil {
		return bson.Marshal(qf.Rule)
	}
	return make([]byte, 0), nil
}

// UnmarshalBSON unmarshal query filter from bson value by first parse bson into map and then parse map into filter
func (qf *QueryFilter) UnmarshalBSON(raw []byte) error {
	data := make(map[string]interface{})
	if err := bson.Unmarshal(raw, &data); err != nil {
		return err
	}

	rule, errKey, err := ParseRule(data)
	if err != nil {
		return fmt.Errorf("parse rule failed, key: %s, err: %v", errKey, err)
	}
	qf.Rule = rule
	return nil
}

// MapToQueryFilterHookFunc TODO
func MapToQueryFilterHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(QueryFilter{}) {
			return data, nil
		}
		if f.Kind() != reflect.Map {
			return data, nil
		}
		dataMap, ok := data.(map[string]interface{})
		if ok == false {
			return data, nil
		}
		rule, errKey, err := ParseRule(dataMap)
		if err != nil {
			return nil, fmt.Errorf("key: %s, err: %s", errKey, err.Error())
		}
		filter := QueryFilter{
			Rule: rule,
		}
		return filter, nil
	}
}
