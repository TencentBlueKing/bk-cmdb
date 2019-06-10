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

package condition

import (
	"encoding/json"
	"testing"

	"configcenter/src/common/mapstr"
)

func TestCondition(t *testing.T) {

	cond := CreateCondition()
	cond.Field("test_field").Eq(1024).Field("test_field2").In([]int{0, 1, 2, 3}).Field("test").Lt(3)

	conditionItem := ConditionItem{Field: "test_field3", Operator: "$lt", Value: 123}
	if err := cond.AddConditionItem(conditionItem); nil != err {
		t.Errorf("AddContionItem error")
		t.Fail()
	}

	if !cond.IsFieldExist("test_field") {
		t.Errorf("IsFieldExist error")
		t.Fail()
	}

	cond.SetPage(mapstr.New())

	cond.SetLimit(1)

	if cond.GetLimit() != 1 {
		t.Fail()
	}

	cond.SetFields([]string{})
	cond.GetFields()
	cond.SetStart(0)
	if cond.GetStart() != 0 {
		t.Fail()
	}

	cond.SetSort("test_field")
	if cond.GetSort() != "test_field" {
		t.Fail()
	}

	result := cond.ToMapStr()
	rst, _ := result.ToJSON()

	t.Logf("the result:%+v", string(rst))

	newCond := CreateCondition()
	err := newCond.Parse(result)
	if nil != err {
		t.Logf("failed to parse condition, error info is %s", err.Error())
		t.Fail()
		return
	}

	rstT, _ := newCond.ToMapStr().ToJSON()
	t.Logf("the parse result:%+v", string(rstT))

}

func TestORCondition(t *testing.T) {

	cond := CreateCondition()
	cond.Field("test_field").Eq(1024).Field("test_field2").In([]int{0, 1, 2, 3}).Field("test").Lt(3)
	cond.SetPage(mapstr.New())
	cond.SetLimit(1)

	if cond.GetLimit() != 1 {
		t.Fail()
	}

	cond.SetFields([]string{})
	cond.GetFields()
	cond.SetStart(0)
	if cond.GetStart() != 0 {
		t.Fail()
	}

	cond.SetSort("test_field")
	if cond.GetSort() != "test_field" {
		t.Fail()
	}

	or := cond.NewOR()
	or.Item(mapstr.MapStr{"a": "b"})
	or.Item(mapstr.MapStr{"b": "c"})
	or.Array([]interface{}{mapstr.MapStr{"c": "b"}, mapstr.MapStr{"d": "b"}})
	or.MapStrArr([]mapstr.MapStr{mapstr.MapStr{"e": "b"}, mapstr.MapStr{"f": "b"}})

	output := `{"$or":[{"a":"b"},{"b":"c"},{"c":"b"},{"d":"b"},{"e":"b"},{"f":"b"}],"test":{"$lt":3},"test_field":1024,"test_field2":{"$in":[0,1,2,3]}}`

	byteOutput, err := json.Marshal(cond.ToMapStr())
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	if string(byteOutput) != output {
		t.Errorf("expected %s not %s", output, string(byteOutput))
		return
	}

}

func TestParseConditionWithMetaData(t *testing.T) {
	data := `{"aa":"a1","bb":"b1","metadata":{"label":{"bk_biz_id":"123"}}}`
	mData := mapstr.MapStr{}
	json.Unmarshal([]byte(data), &mData)
	cond := CreateCondition()
	cond.Parse(mData)
	t.Logf("parse cond from data %v", cond.ToMapStr())
	if !mData.Exists("metadata") {
		t.Fail()
	}

}
