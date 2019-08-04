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
	"fmt"

	"configcenter/src/common/mapstr"
)

type RuleGroupData struct {
	Condition Condition                `json:"condition" field:"condition"`
	Rules     []map[string]interface{} `json:"rules" field:"rules"`
}

func ParseRule(data map[string]interface{}) (queryFilter Rule, errKey string, err error) {
	if _, ok := data["condition"]; ok == true {
		ruleGroupData := &RuleGroupData{}
		// shouldn't use mapstr here as it doesn't support nest struct
		// TODO: replace it with more efficient way
		if err := mapstr.DecodeFromMapStr(ruleGroupData, data); err != nil {
			return nil, "", fmt.Errorf("decode to rule group struct failed, err: %+v", err)
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
			combinedRule.Rules = append(combinedRule.Rules, qf)
		}
		queryFilter = combinedRule
	} else if _, ok := data["operator"]; ok == true {
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
