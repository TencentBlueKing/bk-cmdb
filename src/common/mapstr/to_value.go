/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mapstr

import (
	"encoding/json"
)

func toBool(val interface{}) bool {
	if val, ok := val.(bool); ok {
		return val
	}
	return false
}

func toInt(val interface{}) int {
	switch t := val.(type) {
	default:
		return 0
	case float32:
		return int(t)
	case float64:
		return int(t)
	case int:
		return t
	case int16:
		return int(t)
	case int32:
		return int(t)
	case int64:
		return int(t)
	case int8:
		return int(t)
	case json.Number:
		data, _ := t.Int64()
		return int(data)
	}
}

func toUint(val interface{}) uint {

	switch t := val.(type) {
	default:
		return 0
	case float32:
		return uint(t)
	case float64:
		return uint(t)
	case uint:
		return t
	case uint16:
		return uint(t)
	case uint32:
		return uint(t)
	case uint64:
		return uint(t)
	case uint8:
		return uint(t)
	case json.Number:
		data, _ := t.Float64()
		return uint(data)
	}
}

func toFloat(tagVal interface{}) float64 {
	switch t := tagVal.(type) {
	default:
		return float64(0)
	case float32:
		return float64(t)
	case float64:
		return float64(t)
	case int:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case int8:
		return float64(t)
	case uint:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case uint8:
		return float64(t)
	case json.Number:
		data, _ := t.Float64()
		return data
	}
}
