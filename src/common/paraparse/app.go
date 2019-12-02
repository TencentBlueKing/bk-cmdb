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

package params

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// common search struct
type SearchParams struct {
	Condition map[string]interface{} `json:"condition"`
	Page      map[string]interface{} `json:"page,omitempty"`
	Fields    []string               `json:"fields,omitempty"`
}

func ParseCommonParams(input []metadata.ConditionItem, output map[string]interface{}) error {
	for _, i := range input {
		switch i.Operator {
		case common.BKDBEQ:
			if reflect.TypeOf(i.Value).Kind() == reflect.String {
				output[i.Field] = SpecialCharChange(i.Value.(string))
			} else {
				output[i.Field] = i.Value
			}
		case common.BKDBLIKE:
			regex := make(map[string]interface{})
			regex[common.BKDBLIKE] = i.Value
			output[i.Field] = regex

		case common.BKDBMULTIPLELike:
			multi, ok := i.Value.([]interface{})
			if !ok {
				return fmt.Errorf("operator %s only support for string array", common.BKDBMULTIPLELike)
			}
			fields := make([]interface{}, 0)
			for _, m := range multi {
				mstr, ok := m.(string)
				if !ok {
					return fmt.Errorf("operator %s only support for string array", common.BKDBMULTIPLELike)
				}
				fields = append(fields, mapstr.MapStr{i.Field: mapstr.MapStr{common.BKDBLIKE: mstr}})
			}
			if len(fields) != 0 {
				// only when the fields is none empty, then the fields is valid.
				// a or operator can not have a empty value in mongodb.
				output[common.BKDBOR] = fields
			}
		default:
			d := make(map[string]interface{})
			if i.Value == nil {
				d[i.Operator] = i.Value
			} else if reflect.TypeOf(i.Value).Kind() == reflect.String {
				d[i.Operator] = SpecialCharChange(i.Value.(string))
			} else {
				d[i.Operator] = i.Value
			}
			output[i.Field] = d
		}
	}
	return nil
}

func SpecialCharChange(targetStr string) string {

	re := regexp.MustCompile("[.()\\\\|\\[\\]\\-\\*{}\\^\\$\\?]")
	delItems := re.FindAllString(targetStr, -1)
	tmp := map[string]struct{}{}
	for _, target := range delItems {
		if _, ok := tmp[target]; ok {
			continue
		}
		tmp[target] = struct{}{}
		targetStr = strings.Replace(targetStr, target, fmt.Sprintf(`\%s`, target), -1)
	}

	return targetStr
}
