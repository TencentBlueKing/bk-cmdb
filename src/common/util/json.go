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
 
package util

import (
	"fmt"
	"reflect"
)

// MapMatch return whether src is partial match to tar,
// means src is smaller than tar
func MapMatch(src, tar interface{}) bool {
	if tar == nil || src == nil ||
		reflect.TypeOf(src).Kind() != reflect.Map ||
		reflect.TypeOf(tar).Kind() != reflect.Map {
		return false
	}

	tarMap := reflect.ValueOf(tar)
	srcMap := reflect.ValueOf(src)

	for _, k := range srcMap.MapKeys() {
		tv := tarMap.MapIndex(k)
		sv := srcMap.MapIndex(k)
		if !tv.IsValid() || !sv.IsValid() {
			return false
		}
		if reflect.DeepEqual(tv, sv) {
			continue
		}
		if sv.IsNil() || tv.IsNil() {
			return false
		}

		tf, tok := tv.Interface().(unix)
		sf, sok := tv.Interface().(unix)
		if tok && sok {
			dt := tf.Unix() - sf.Unix()
			if dt <= -1 || dt >= 1 {
				return false
			}
			continue
		}
		if fmt.Sprint(sv) == fmt.Sprint(tv) {
			continue
		}
		return false
	}
	return true
}

type unix interface {
	Unix() int64
}
