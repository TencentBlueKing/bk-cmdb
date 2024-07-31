/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package filter

import (
	"fmt"
	"strings"
	"time"

	"configcenter/src/common/util"

	"github.com/tidwall/gjson"
)

// MatchedData is the data to be matched with the expression's rule
type MatchedData interface {
	GetValue(field string) (interface{}, error)
}

// JsonString is the json string type
type JsonString string

// GetValue get value by field from json string
func (j JsonString) GetValue(field string) (interface{}, error) {
	if j == "" {
		return nil, fmt.Errorf("json data is empty")
	}
	val := gjson.Get(string(j), field)
	switch val.Type {
	case gjson.Null:
		return nil, nil
	case gjson.Number:
		if strings.Contains(val.Raw, ".") {
			return val.Float(), nil
		}
		return val.Int(), nil
	default:
		return val.Value(), nil
	}
}

// MapStr is the map[string]interface{} type
type MapStr map[string]interface{}

// GetValue get value by field from map
func (m MapStr) GetValue(field string) (interface{}, error) {
	if m == nil {
		return nil, fmt.Errorf("map data is nil")
	}
	return m[field], nil
}

func parseNumericValues(value1, value2 interface{}) (float64, float64, error) {
	val1, err := util.GetFloat64ByInterface(value1)
	if err != nil {
		return 0, 0, fmt.Errorf("parse input value(%+v) failed, err: %v", value1, err)
	}

	val2, err := util.GetFloat64ByInterface(value2)
	if err != nil {
		return 0, 0, fmt.Errorf("parse rule value(%+v) failed, err: %v", value2, err)
	}

	return val1, val2, nil
}

func parseTimeValues(value1, value2 interface{}) (time.Time, time.Time, error) {
	val1, err := util.ConvToTime(value1)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse input value(%+v) failed, err: %v", value1, err)
	}

	val2, err := util.ConvToTime(value2)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse rule value(%+v) failed, err: %v", value2, err)
	}

	return val1, val2, nil
}

func parseStringValues(value1, value2 interface{}) (string, string, error) {
	val1, ok := value1.(string)
	if !ok {
		return "", "", fmt.Errorf("input value(%+v) is not string type", value1)
	}

	val2, ok := value2.(string)
	if !ok {
		return "", "", fmt.Errorf("rule value(%+v) is not string type", value2)
	}

	return val1, val2, nil
}
