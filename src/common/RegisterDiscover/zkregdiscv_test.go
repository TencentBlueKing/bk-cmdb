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
 
package RegisterDiscover

import (
	"testing"
	"time"
)

func Test_sortNode(t *testing.T) {

	zkRegDcv := NewZkRegDiscv("127.0.0.1:2181", time.Second*60)
	t.Log("----- test node1 -----")
	nodes1 := []string{"json.info_0000000003", "json.info_0000000002", "json.info_0000000000", "log_replicas"}
	t.Log(nodes1)
	sortNs1 := zkRegDcv.sortNode(nodes1)
	t.Logf("%+v\n", sortNs1)

	t.Log("----- test node2 -----")
	nodes2 := []string{"_c_cd91a29e1a89e5346014994493ff427c-10.49.110.1620000000087",
		"_c_bc551189ace5a4a2af2779028a2b14d3-10.223.38.880000000086",
		"_c_fb6af66c33df0fbec44fdaf4c820e531-10.49.110.1610000000085",
		"_c_4cf8f02920c631803ba7da24200c9b89-10.223.55.880000000084"}

	t.Log(nodes2)
	sortNs2 := zkRegDcv.sortNode(nodes2)
	t.Logf("%+v\n", sortNs2)
}

// go test -v -test.run Test_RegisterAndWatch
func Test_RegisterAndWatch(t *testing.T) {

	zkRegDcv := NewZkRegDiscv("127.0.0.1:2181", time.Second*60)
	zkRegDcv.Start()
	t.Log("----- start test -----")

	zkRegDcv.RegisterAndWatch("/admin/test", []byte("test"))
	// delete node path
	time.Sleep(5 * time.Minute)
	t.Log("----- end test -----")
}
