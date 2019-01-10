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
	"errors"
	"reflect"
)

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
