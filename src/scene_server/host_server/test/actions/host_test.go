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

package actions_test

import (
	"strconv"

	_ "configcenter/src/scene_server/topo_server/actions/inst"   // import inst
	_ "configcenter/src/scene_server/topo_server/actions/object" // import object actions
	_ "configcenter/src/scene_server/topo_server/actions/openapi"
	_ "configcenter/src/scene_server/topo_server/actions/privilege"
	_ "configcenter/src/scene_server/topo_server/logics/object"
	"testing"
)

func TestHost(t *testing.T) {
	testHostAdd(t)
	sids := testHostSearch(t)
	ids := []int{}
	for _, id := range sids {
		iid, _ := strconv.Atoi(id)
		ids = append(ids, iid)
	}
	testGetHostDetailByID(t, sids[0])
	testHostSnapInfo(t, sids[0])

	testMoveHostToResourcePool(t, ids)
	testHostModuleRelation(t, ids)
	testMoveIDleModule(t, ids)
	testMoveFaultModule(t, ids)
	testAssignHostToApp(t, ids)

	testDelete(t, sids...)
}
