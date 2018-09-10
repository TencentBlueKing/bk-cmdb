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
	"testing"
)

func TestField(t *testing.T) {
	f := &field{
		condition: CreateCondition(),
	}

	f.Eq("f_eq")
	f.NotEq("f_neq")
	f.Like("f_like")
	f.In("f_in")
	f.NotIn("f_notin")
	f.Lt("f_lt")
	f.Lte("f_lte")
	f.Gt("f_gt")
	f.Gte("f_gte")
	_, err := f.ToMapStr().ToJSON()
	if err != nil {
		t.Fail()
	}

}
