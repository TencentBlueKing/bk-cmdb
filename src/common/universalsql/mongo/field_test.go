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

package mongo

import (
	"testing"
)

func TestComparisionField(t *testing.T) {
	sql, _ := Field("name").Eq("sam").Neq("uri  ").In([]string{"jim", "berg"}).ToSQL()
	t.Logf("%s", sql)

	sql, _ = Field("age").Lt(100).Gte(10).In([]int{22, 35}).Nin([]int{44, 54, 64}).Regex("jim").ToSQL()
	t.Logf("%s", sql)
}

func TestElementField(t *testing.T) {
	sql, _ := Field("school").Exists(true).ToSQL()
	t.Logf("%s", sql)
}

func TestArrayField(t *testing.T) {
	sql, _ := Field("school").Size(2).All([]string{"qinghua", "beida"}).ToSQL()
	t.Logf("%s", sql)
}
