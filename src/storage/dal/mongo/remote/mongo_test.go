/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package remote

import (
	"testing"
)

func TestDailRPC(t *testing.T) {
	db, err := New("192.168.100.130:60008")
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatalf("%v", err)
	}
}
