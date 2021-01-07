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

package mapstruct

import (
	"encoding/json"

	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"

	"github.com/mitchellh/mapstructure"
)

// Decode2StructWithHook conv map to struct with hook func that can conv time and queryFilter
// 适用场景：用于查询操作，转 查询条件map to struct(重点：转 string to timeDuration 和 map to queryFilter)
func Decode2StructWithHook(m map[string]interface{}, st interface{}) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			metadata.StringToTimeDurationHookFunc(),
			querybuilder.MapToQueryFilterHookFunc(),
		),
		WeaklyTypedInput: true,
		Result:           &st,
	}
	dec, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	if err := dec.Decode(m); err != nil {
		return err
	}
	return nil
}

// Decode2Struct conv map to struct
// 适用场景：仅支持原生转 map to struct
func Decode2Struct(m map[string]interface{}, st interface{}) error {
	if err := mapstructure.Decode(m, &st); err != nil {
		return err
	}
	return nil
}

func Struct2Map(v interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return data, nil
}
