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

package mapstr_test

import (
	. "configcenter/src/common/mapstr"
	"fmt"
	"testing"
	"time"
)

func TestMapStrInto(t *testing.T) {
	type testData struct {
		Ignor int
		Data  string `field:"data"`
		Test  int    `field:"test"`
	}
	target := New()
	target.Set("test", 245)
	target.Set("data", "test_data")

	tmp := &testData{}
	target.MarshalJSONInto(tmp)
	//t.Logf("the test tmp %#v", tmp)

	maps := NewArrayFromInterface([]map[string]interface{}{
		{"k": "value"}, {"i": 0},
	})
	target1 := maps[0]

	target2, err := NewFromInterface(map[string]interface{}{
		"k": "v", "i": 1, "j": 2, "time": time.Now(), "map": map[string]interface{}{}, "bool": true,
	})
	if err != nil {
		t.Fail()
	}
	target1.Different(target2)
	target1.Merge(target2)

	target1.Set("set_key", "set_value")
	_, ok := target1.Get("set_key")
	if !ok {
		t.Fail()
	}

	if i, _ := target1.Int64("i"); i != 1 {
		t.Fail()
	}

	if b, _ := target1.Bool("bool"); !b {
		t.Fail()
	}

	if i, _ := target1.Float("i"); i != 1 {
		t.Fail()
	}

	if s, _ := target1.String("k"); s != "v" {
		t.Fail()
	}

	if _, err := target1.Time("time"); err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if _, err := target1.MapStr("map"); err != nil {
		t.Fail()
	}

	target1.Set("maps", maps)
	if _, err := target1.MapStrArray("maps"); err != nil {
		t.Fail()
	}

	if target1.ForEach(func(key string, val interface{}) error {
		return nil
	}) != nil {
		t.Fail()
	}

	target1.Remove("maps")
	if target1.Exists("maps") {
		t.Fail()
	}

	if target1.IsEmpty() {
		t.Fail()
	}

	target1.Reset()

}

func TestMapStrToMapstr(t *testing.T) {

	testData := MapStr{"aa": "bb"}

	_, err := NewFromInterface(testData)
	if err != nil {
		t.Error(err.Error())
		return
	}

	testData2 := MapStr{"aa": []MapStr{
		MapStr{"aa": "bb"},
	}}
	_, err = testData2.MapStrArray("aa")
	if err != nil {
		t.Error(err.Error())
		return
	}
}
