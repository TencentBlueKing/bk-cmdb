/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

package base

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// BenchmarkSetID result on length = 500
//
//	                direct:    530.56 ns/op 1.00x
//	               generic:   1193.46 ns/op 2.25x
//	     interface_nocheck:   1535.09 ns/op 2.89x
//	    interface_precheck:   1766.52 ns/op 3.33x
//	  interface_each_check:   1856.79 ns/op 3.50x
//	reflect_field_by_index:  23154.61 ns/op 43.64x
//	 reflect_field_by_name: 196055.81 ns/op 369.53x
func BenchmarkSetID(b *testing.B) {
	type benchCase struct {
		name string
		fn   func([]testModel, []string)
	}

	cases := []benchCase{
		{name: "direct", fn: directSet},
		{name: "generic", fn: genericSet[testModel, *testModel]},
		{name: "interface_nocheck", fn: interfaceSetNoCheck[testModel]},
		{name: "interface_precheck", fn: interfaceSetPrecheckType[testModel]},
		{name: "interface_each_check", fn: interfaceSetCheckTypeEveryTime[testModel]},
		{name: "reflect_field_by_index", fn: reflectSetFieldByIndex[testModel]},
		{name: "reflect_field_by_name", fn: reflectSetFieldByName[testModel]},
	}

	for _, length := range []int{1, 500} {
		b.Run(fmt.Sprintf("SetID-%d", length), func(b *testing.B) {
			var a = make([]testModel, length)
			var ids = make([]string, length)
			for i := range ids {
				ids[i] = fmt.Sprintf("id-%d", i)
			}
			var cost time.Duration
			var results = make([]float64, len(cases))
			for i, c := range cases {
				b.Run(c.name, func(b *testing.B) {
					start := time.Now()
					for range b.N {
						c.fn(a, ids)
					}
					cost = time.Since(start)
					results[i] = float64(cost) / float64(b.N)
				})
			}
			base := results[0]
			for i := range results {
				b.Logf("%22s: %.2f ns/op %.2fx\n", cases[i].name, results[i], results[i]/base)
			}
		})
	}
}

func directSet(s []testModel, ids []string) {
	for i := range s {
		(&s[i]).ID = ids[i]
	}
}

type settable[T any] interface {
	SetID(string)
	*T
}

func genericSet[T any, PT settable[T]](s []T, ids []string) {
	for i := range s {
		PT(&s[i]).SetID(ids[i])
	}
}
func interfaceSetPrecheckType[T any](s []T, ids []string) {
	var t T
	if _, ok := any(&t).(interface{ SetID(string) }); !ok {
		panic("not settable")
	}
	for i := range s {
		any(&s[i]).(interface{ SetID(string) }).SetID(ids[i])
	}
}
func interfaceSetNoCheck[T any](s []T, ids []string) {
	for i := range s {
		any(&s[i]).(interface{ SetID(string) }).SetID(ids[i])
	}
}

func interfaceSetCheckTypeEveryTime[T any](s []T, ids []string) {
	var t T
	var rt = reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		panic("should not be ptr")
	}
	for i := range s {
		if settable, ok := any(&s[i]).(interface{ SetID(string) }); ok {
			settable.SetID(ids[i])
		} else {
			panic("not settable")
		}
	}
}

func reflectSetFieldByName[T any](s []T, ids []string) {
	sv := reflect.ValueOf(s)
	_, ok := sv.Type().Elem().FieldByName("ID")
	if !ok {
		panic("unknown field id")
	}
	for i := range s {
		v := sv.Index(i)
		v.FieldByName("ID").Set(reflect.ValueOf(ids[i]))
	}
}
func reflectSetFieldByIndex[T any](s []T, ids []string) {
	sv := reflect.ValueOf(s)
	st, ok := sv.Type().Elem().FieldByName("ID")
	if !ok {
		panic("unknown field id")
	}
	for i := range s {
		v := sv.Index(i)
		v.FieldByIndex(st.Index).Set(reflect.ValueOf(ids[i]))
	}
}

type testModel struct {
	base     `gorm:"embedded" json:",inline"`
	Name     string    `gorm:"column:name" json:"name,omitempty"`
	Size     int       `gorm:"column:size" json:"size,omitempty"`
	Weight   float64   `gorm:"column:weight" json:"weight,omitempty"`
	Int64s   []int     `gorm:"column:int64s" json:"int64s,omitempty"`
	Strings  []string  `gorm:"column:strings" json:"strings,omitempty"`
	Strings2 *[]string `gorm:"column:strings2" json:"strings2,omitempty"`
}

type base struct {
	ID           string    `json:"id,omitempty"`
	SysCreatedAt time.Time `json:"sys_created_at"`
	SysCreatedBy string    `json:"sys_created_by,omitempty"`
	SysUpdatedAt time.Time `json:"sys_updated_at"`
	SysUpdatedBy string    `json:"sys_updated_by,omitempty"`
}

// SetID ...
func (t *base) SetID(id string) {
	t.ID = id
}
