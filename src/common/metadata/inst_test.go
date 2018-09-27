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

import "testing"

func TestSort(t *testing.T) {
	sorter := HostIdentProcessSorter{
		{BindModules: []int64{1, 2}, ProcessID: 1}, {BindModules: []int64{2, 3}, ProcessID: 2},
	}

	sorter.Swap(0, 1)
	if sorter.Len() != 2 {
		t.Fail()
	}
	if sorter.Less(0, 1) {
		t.Fail()
	}

}
