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

package discovery

import (
	"fmt"
	"testing"

	"configcenter/src/common/types"
)

func TestGetServers(t *testing.T) {
	svr := server{
		name: "demo",
		path: "/",
		servers: []*types.ServerInfo{
			{
				IP:         "127.0.0.1",
				Port:       8081,
				RegisterIP: "127.0.0.1",
				Scheme:     "http",
				UUID:       "1",
			},
			{
				IP:         "127.0.0.2",
				Port:       8082,
				RegisterIP: "127.0.0.2",
				Scheme:     "http",
				UUID:       "2",
			},
			{
				IP:         "127.0.0.3",
				Port:       8083,
				RegisterIP: "127.0.0.3",
				Scheme:     "http",
				UUID:       "3",
			},
			{
				IP:         "127.0.0.4",
				Port:       8084,
				RegisterIP: "127.0.0.4",
				Scheme:     "http",
				UUID:       "4",
			},
		},
		discoverChan: nil,
		serversChan:  nil,
	}

	sum := make(map[string]int)
	for i := 0; i < 10000; i++ {
		s, _ := svr.GetServers()
		sum[s[0]] += 1
	}
	fmt.Println(sum)

}
