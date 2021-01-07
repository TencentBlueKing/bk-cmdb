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
	"fmt"
	"reflect"
)

// New create a new MapStr instance
func New() MapStr {
	return MapStr{}
}

// NewArray create MapStr array
func NewArray() []MapStr {
	return []MapStr{}
}

// NewArrayFromMapStr create a new array from mapstr array
func NewArrayFromMapStr(datas []MapStr) []MapStr {
	results := []MapStr{}
	for _, item := range datas {
		results = append(results, item)
	}
	return results
}

// NewFromInterface create a mapstr instance from the interface
// Support Input Type: []byte, string, base-type map, struct.
// If the input value type is []byte or string, then the value must be a valid json.
// Like: map[string]int will be converted into MapStr
// Like: struct { TestStr string TestInt int } will be converted into  MapStr{"TestStr":"", "TestInt":0}
func NewFromInterface(data interface{}) (MapStr, error) {

	switch tmp := data.(type) {
	default:
		return convertInterfaceIntoMapStrByReflection(data, "")
	case nil:
		return MapStr{}, nil
	case MapStr:
		return tmp, nil
	case []byte:
		result := New()
		if 0 == len(tmp) {
			return result, nil
		}
		err := json.Unmarshal(tmp, &result)
		return result, err
	case string:
		result := New()
		if 0 == len(tmp) {
			return result, nil
		}
		err := json.Unmarshal([]byte(tmp), &result)
		return result, err
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

// NewFromMap create a new MapStr from map[string]interface{} type
func NewFromMap(data map[string]interface{}) MapStr {
	return MapStr(data)
}

// NewFromStruct convert the  struct into MapStr , the struct must be taged with 'tagName' .
//
//  eg:
//  type targetStruct struct{
//       Name string `field:"testName"`
//  }
//  will be converted the follow map
//  {"testName":""}
//
func NewFromStruct(targetStruct interface{}, tagName string) MapStr {
	return SetValueToMapStrByTagsWithTagName(targetStruct, tagName)
}

// NewArrayFromInterface create a new array from interface
func NewArrayFromInterface(datas []map[string]interface{}) []MapStr {
	results := []MapStr{}
	for _, item := range datas {
		results = append(results, item)
	}
	return results
}

// SetValueToMapStrByTags  convert a struct to MapStr by tags default tag name is field
func SetValueToMapStrByTags(source interface{}) MapStr {
	return SetValueToMapStrByTagsWithTagName(source, "field")
}

// SetValueToMapStrByTagsWithTagName convert a struct to MapStr by tags
func SetValueToMapStrByTagsWithTagName(source interface{}, tagName string) MapStr {
	values := MapStr{}

	if nil == source {
		return values
	}

	targetType := getTypeElem(reflect.TypeOf(source))
	targetValue := getValueElem(reflect.ValueOf(source))

	setMapStrByStruct(targetType, targetValue, values, tagName)

	return values
}

// SetValueToStructByTags set the struct object field value by tags, default tag name is field
func SetValueToStructByTags(target interface{}, values MapStr) error {
	return SetValueToStructByTagsWithTagName(target, values, "field")
}

// SetValueToStructByTagsWithTagName set the struct object field value by tags
func SetValueToStructByTagsWithTagName(target interface{}, values MapStr, tagName string) error {

	targetType := reflect.TypeOf(target)
	targetValue := reflect.ValueOf(target)

	return setStructByMapStr(targetType, targetValue, values, tagName)
}

func convertInterfaceIntoMapStrByReflection(target interface{}, tagName string) (MapStr, error) {

	value := reflect.ValueOf(target)
	switch value.Kind() {
	case reflect.Map:
		return dealMap(value, tagName)
	case reflect.Struct:
		return dealStruct(value.Type(), value, tagName)
	}

	return nil, fmt.Errorf("no support the kind(%s)", value.Kind())
}
