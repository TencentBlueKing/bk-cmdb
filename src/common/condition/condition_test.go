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

package condition_test

import (
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"testing"
)

func TestCondition(t *testing.T) {

	cond := condition.CreateCondition()
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

	result := cond.ToMapStr()
	rst, _ := result.ToJSON()

	t.Logf("the result:%+v", string(rst))

	newCond := condition.CreateCondition()
	err := newCond.Parse(result)
	if nil != err {
		t.Logf("failed to parse condition, error info is %s", err.Error())
		return
	}

	rstT, _ := newCond.ToMapStr().ToJSON()
	t.Logf("the parse result:%+v", string(rstT))

}
