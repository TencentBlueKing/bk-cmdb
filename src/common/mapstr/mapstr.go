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

package mapstr

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// MapStr the common event data definition
type MapStr map[string]interface{}

// New create a new MapStr instance
func New() MapStr {
	return MapStr{}
}

// NewArrayFromInterface create a new array from interface
func NewArrayFromInterface(datas []map[string]interface{}) []MapStr {
	results := []MapStr{}
	for _, item := range datas {
		results = append(results, item)
	}
	return results
}

// NewFromInterface create a mapstr instance from the interface
func NewFromInterface(data interface{}) (MapStr, error) {

	switch tmp := data.(type) {
	default:
		return nil, fmt.Errorf("not support the kind(%s)", reflect.TypeOf(data).Kind())
	case nil:
		return MapStr{}, nil
	case *map[string]interface{}:
		return MapStr(*tmp), nil
	case map[string]string:
		result := New()
		for key, val := range tmp {
			result.Set(key, val)
		}
		return result, nil
	case map[string]interface{}:
		return MapStr(tmp), nil
	}
}

// Merge merge second into self,if the key is the same then the new value replaces the old value.
func (cli MapStr) Merge(second MapStr) {
	for key, val := range second {
		cli[key] = val
	}
}

// MarshalJSONInto convert to the input value
func (cli MapStr) MarshalJSONInto(target interface{}) error {

	data, err := cli.ToJSON()
	if nil != err {
		return err
	}

	return json.Unmarshal(data, target)
}

// ToJSON convert to json string
func (cli MapStr) ToJSON() ([]byte, error) {
	js, err := json.Marshal(cli)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// Get return the origin value by the key
func (cli MapStr) Get(key string) (val interface{}, exists bool) {

	val, exists = cli[key]
	return val, exists
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
func (cli MapStr) Bool(key string) (bool, error) {
	switch t := cli[key].(type) {
	case nil:
		return false, fmt.Errorf("the key (%s) is invalid", key)
	default:
		return false, fmt.Errorf("the key (%s) is invalid", key)
	case bool:
		return t, nil
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
		return strconv.ParseInt(t, 10, 64)
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
func (cli MapStr) String(key string) (string, error) {
	switch t := cli[key].(type) {
	case nil:
		return "", nil
	default:
		return fmt.Sprintf("%v", t), nil
	case map[string]interface{}, []interface{}:
		rest, err := json.Marshal(t)
		if nil != err {
			return "", err
		}
		return string(rest), nil
	case json.Number:
		return t.String(), nil
	case string:
		return t, nil
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
			return nil, fmt.Errorf("the data is not a valid type,%s", val.Kind().String())
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
