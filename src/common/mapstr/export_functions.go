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
	"errors"
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

// NewArrayFromInterface create a new array from interface
func NewArrayFromInterface(datas []map[string]interface{}) []MapStr {
	results := []MapStr{}
	for _, item := range datas {
		results = append(results, item)
	}
	return results
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
func NewFromInterface(data interface{}) (MapStr, error) {

	switch tmp := data.(type) {
	default:
		return nil, fmt.Errorf("no support the kind(%s)", reflect.TypeOf(data).Kind())
	case nil:
		return MapStr{}, nil
	case MapStr:
		return tmp, nil
	case string:
		result := New()
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

/*NewFromStruct convert the  struct into MapStr , the struct must be taged with 'tagName' .

  eg:
  type targetStruct struct{
  Name string `field:"testName"`
  }
  will be converted the follow map
  {"testName":""}
*/
func NewFromStruct(targetStruct interface{}, tagName string) MapStr {
	return SetValueToMapStrByTagsWithTagName(targetStruct, tagName)
}

// ConvertArrayMapStrInto convert a MapStr array into a new slice instance
func ConvertArrayMapStrInto(datas []MapStr, output interface{}) error {

	resultv := reflect.ValueOf(output)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		return errors.New("result argument must be a slice address")
	}
	slicev := resultv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	idx := 0
	for _, dataItem := range datas {
		if slicev.Len() == idx {
			elemp := reflect.New(elemt)
			if err := dataItem.MarshalJSONInto(elemp.Interface()); nil != err {
				panic(err)
			}
			slicev = reflect.Append(slicev, elemp.Elem())
			slicev = slicev.Slice(0, slicev.Cap())
			idx++
			continue
		}

		if err := dataItem.MarshalJSONInto(slicev.Index(idx).Addr().Interface()); nil != err {
			return err
		}
		idx++
	}
	resultv.Elem().Set(slicev.Slice(0, idx))

	return nil
}
