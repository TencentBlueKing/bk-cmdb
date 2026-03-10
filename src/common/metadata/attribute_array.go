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

package metadata

import (
	"configcenter/src/common/util"
	"encoding/json"
	"fmt"
	"math"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
)

// ArrayOption len cap ,option is basic type 's option
type ArrayOption[T any] struct {
	Len    int `bson:"len" json:"len"`
	Cap    int `bson:"cap" json:"cap"`
	Option T   `bson:"option" json:"option"`
}

// Valid ArrayOption
func (a *ArrayOption[T]) Valid() error {
	if a.Len < 0 || a.Len > a.Cap {
		return fmt.Errorf("invalid array option,len:%d cap:%d", a.Len, a.Cap)
	}
	return nil
}

// ParseArrayOption len cap ,option is basic type 's option
func ParseArrayOption[T any](option any, handle func(v any) (T, error)) (ArrayOption[T], error) {
	if option == nil {
		return ArrayOption[T]{Len: math.MaxInt, Cap: math.MaxInt}, nil
	}

	var result ArrayOption[T]

	optMap := map[string]interface{}{
		"len": math.MaxInt,
		"cap": math.MaxInt,
	}
	switch value := option.(type) {
	case ArrayOption[T]:
		return value, nil
	case bson.M:
		optMap = value
	case map[string]interface{}:
		optMap = value
	default:
		marshal, err := json.Marshal(option)
		if err != nil {
			return result, fmt.Errorf("invalid array option,type:%v,value:%v,err:%w",
				option, option, err)
		}

		lenItem := gjson.GetBytes(marshal, "len")
		capItem := gjson.GetBytes(marshal, "cap")
		if !lenItem.Exists() || !capItem.Exists() {
			return result, fmt.Errorf("invalid array option,type:%v,value:%v,err: not exist len or cap", option, option)
		}
		optMap["len"] = lenItem.Int()
		optMap["cap"] = capItem.Int()
		optMap["option"] = gjson.GetBytes(marshal, "option").Value()
	}

	lenn, lenOk := optMap["len"]
	capp, capOk := optMap["cap"]
	if !lenOk || !capOk {
		return result, fmt.Errorf("invalid array option,type:%v,value:%v,err: not exist len or cap", option, option)
	}

	lenOpt, err := util.GetIntByInterface(lenn)
	if err != nil {
		return result, err
	}
	result.Len = lenOpt
	capOpt, err := util.GetIntByInterface(capp)
	if err != nil {
		return result, err
	}
	result.Cap = capOpt

	var defaultOption T
	result.Option = defaultOption
	if handle != nil {
		t, err := handle(optMap["option"])
		if err != nil {
			return ArrayOption[T]{}, err
		}
		result.Option = t
	}
	return result, result.Valid()
}
