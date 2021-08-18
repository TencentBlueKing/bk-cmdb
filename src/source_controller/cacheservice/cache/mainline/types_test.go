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

package mainline

import (
	"testing"
)

func TestObserver(t *testing.T) {
	observer := watchObserver{
		observer: make(map[string]chan struct{}),
	}

	// add a custom object
	notifier := make(chan struct{})
	custom := "country"
	observer.add("country", notifier)

	// test a custom is exist or not
	if !observer.exist(custom) {
		t.Fatal("test country custom level exist failed, not exist")
	}

	// test get custom object
	if observer.getAllObjects()[0] != custom {
		t.Fatal("test country custom level get failed, not equal")
	}

	// test delete custom object
	observer.delete(custom)
	if observer.exist(custom) {
		t.Fatal("test country custom level delete failed, still exist")
	}

}
