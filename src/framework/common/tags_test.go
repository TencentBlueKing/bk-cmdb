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

package common_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/types"
	"testing"
)

type testObj struct {
	Filed1 string      `field:"field_one"`
	Filed2 bool        `field:"field_two"`
	Filed3 int         `field:"field_three"`
	Filed4 int64       `field:"field_four"`
	Data   interface{} `field:"field_five"`
}

func TestGetTags(t *testing.T) {
	obj := &testObj{}
	tags := common.GetTags(obj)
	t.Logf("tags:%v", tags)
}

func TestSetValueByTags(t *testing.T) {
	obj := &testObj{}
	data := types.MapStr{
		"field_one":   "test_one_value",
		"field_two":   true,
		"field_three": 3,
		"field_four":  4,
		"field_five": map[string]interface{}{
			"filed_t": 0,
		},
	}
	common.SetValueToStructByTags(obj, data)
	t.Logf("tags:%v", obj)
}

func TestSetValueToMapStrByTags(t *testing.T) {
	obj := &testObj{
		Filed1: "field_1",
		Filed2: false,
		Filed3: 3,
		Filed4: 4,
	}

	data := common.SetValueToMapStrByTags(obj)
	t.Logf("tags:%#v", data)
}
