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

package types

import (
	"configcenter/src/framework/core/log"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Get return the origin value by the key
func (cli MapStr) Get(key string) (val interface{}, exists bool) {

	val, exists = cli[key]
	return val, exists
}

// Merge merge second into self,if the key is the same then the new value replaces the old value.
func (cli MapStr) Merge(second MapStr) {
	for key, val := range second {
		cli[key] = val
	}
}

// ToJSON convert to json string
func (cli MapStr) ToJSON() []byte {
	js, err := json.Marshal(cli)
	if err != nil {
		log.Errorf("to json error %v", err)
	}
	return js
}

// Set set a new value for the key, the old value will be replaced
func (cli MapStr) Set(key string, value interface{}) {
	cli[key] = value
}

// Reset  reset the mapstr into the init state
func (cli MapStr) Reset() {
	for key := range cli {
		delete(cli, key)
	}
}

// Bool get the value as bool
func (cli MapStr) Bool(key string) bool {
	switch t := cli[key].(type) {
	case nil:
		return false
	default:
		return false
	case bool:
		return t
	}
}

// Int64 return the value by the key
func (cli MapStr) Int64(key string) (int64, error) {

	switch t := cli[key].(type) {
	default:
		return 0, errors.New("invalid num")
	case nil:

		return 0, errors.New("invalid key(" + key + "), not found value")
	case int:
		return int64(t), nil
	case int16:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case int64:
		return t, nil
	case float32:
		return int64(t), nil
	case float64:
		return int64(t), nil
	case json.Number:
		num, err := t.Int64()
		return int64(num), err
	case string:
		tv, err := strconv.Atoi(t)
		if nil != err {
			return 0, err
		}
		return int64(tv), nil
	}
}

// Int return the value by the key
func (cli MapStr) Int(key string) (int, error) {

	switch t := cli[key].(type) {
	default:
		return 0, errors.New("invalid num")
	case nil:

		return 0, errors.New("invalid key(" + key + "), not found value")
	case int:
		return t, nil
	case int16:
		return int(t), nil
	case int32:
		return int(t), nil
	case int64:
		return int(t), nil
	case float32:
		return int(t), nil
	case float64:
		return int(t), nil
	case json.Number:
		num, err := t.Int64()
		return int(num), err
	case string:
		return strconv.Atoi(t)
	}
}

// Float get the value as float64
func (cli MapStr) Float(key string) (float64, error) {
	switch t := cli[key].(type) {
	default:
		return 0, errors.New("invalid num")
	case nil:
		return 0, errors.New("invalid key, not found value")
	case int:
		return float64(t), nil
	case int16:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float32:
		return float64(t), nil
	case float64:
		return t, nil
	case json.Number:
		return t.Float64()
	}
}

// String get the value as string
func (cli MapStr) String(key string) string {
	switch t := cli[key].(type) {
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", t)
	case map[string]interface{}, []interface{}:
		rest, _ := json.Marshal(t)
		return string(rest)
	case json.Number:
		return t.String()
	case string:
		return t
	}
}

// Time get the value as time.Time
func (cli MapStr) Time(key string) (*time.Time, error) {
	switch t := cli[key].(type) {
	default:
		return nil, errors.New("invalid time value")
	case nil:
		return nil, errors.New("invalid key")
	case time.Time:
		return &t, nil
	case *time.Time:
		return t, nil
	case string:
		if tm, tmErr := time.Parse(time.RFC1123, t); nil == tmErr {
			return &tm, nil
		}

		if tm, tmErr := time.Parse(time.RFC1123Z, t); nil == tmErr {
			return &tm, nil
		}

		if tm, tmErr := time.Parse(time.RFC3339, t); nil == tmErr {
			return &tm, nil
		}

		if tm, tmErr := time.Parse(time.RFC3339Nano, t); nil == tmErr {
			return &tm, nil
		}

		if tm, tmErr := time.Parse(time.RFC822, t); nil == tmErr {
			return &tm, nil
		}

		if tm, tmErr := time.Parse(time.RFC822Z, t); nil == tmErr {
			return &tm, nil
		}

		if tm, tmErr := time.Parse(time.RFC850, t); nil == tmErr {
			return &tm, nil
		}

		return nil, errors.New("can not parse the datetime")
	}
}

// MapStr get the MapStr object
func (cli MapStr) MapStr(key string) (MapStr, error) {

	switch t := cli[key].(type) {
	default:
		return nil, errors.New("the data is not a map[string]interface{} type")
	case nil:
		if _, ok := cli[key]; ok {
			return MapStr{}, nil
		}
		return nil, errors.New("the key is invalid")
	case map[string]interface{}:
		return MapStr(t), nil
	}

}

// MapStrArray get the MapStr object array
func (cli MapStr) MapStrArray(key string) ([]MapStr, error) {

	switch t := cli[key].(type) {
	default:
		val := reflect.ValueOf(cli[key])
		switch val.Kind() {
		default:

			return []MapStr{
				MapStr{
					key: val.Interface(),
				},
			}, nil

		case reflect.Slice:
			tmpval, ok := val.Interface().([]MapStr)
			if ok {
				return tmpval, nil
			}

			return nil, fmt.Errorf("the data is not a valid type,%s", val.Kind().String())
		}

	case nil:

		return nil, fmt.Errorf("the key(%s) is invalid", key)
	case []map[string]interface{}:
		items := make([]MapStr, 0)
		for _, item := range t {
			items = append(items, item)
		}
		return items, nil
	case []interface{}:
		items := make([]MapStr, 0)
		for _, item := range t {
			switch childType := item.(type) {
			case map[string]interface{}:
				items = append(items, childType)
			}
		}
		return items, nil
	}

}

// ForEach for each the every item
func (cli MapStr) ForEach(callItem func(key string, val interface{})) {

	for key, val := range cli {
		callItem(key, val)
	}

}

// Remove delete the item by the key and return the deleted one
func (cli MapStr) Remove(key string) interface{} {

	if val, ok := cli[key]; ok {
		delete(cli, key)
		return val
	}

	return nil
}

// Exists check the key exists
func (cli MapStr) Exists(key string) bool {
	_, ok := cli[key]
	return ok
}

// IsEmpty check the empty status
func (cli MapStr) IsEmpty() bool {
	return len(cli) == 0
}

// Different the current value is different from the content of the given data
func (cli MapStr) Different(target MapStr) (more MapStr, less MapStr, changes MapStr) {

	// init
	more = make(MapStr)
	less = make(MapStr)
	changes = make(MapStr)

	// check more
	cli.ForEach(func(key string, val interface{}) {
		if targetVal, ok := target[key]; ok {

			if !reflect.DeepEqual(val, targetVal) {
				changes[key] = val
			}
			return
		}

		more.Set(key, val)
	})

	// check less
	target.ForEach(func(key string, val interface{}) {
		if !cli.Exists(key) {
			less[key] = val
		}

	})

	return more, less, changes
}
