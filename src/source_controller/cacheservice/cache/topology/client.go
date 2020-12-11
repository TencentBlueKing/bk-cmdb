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

package topology

import (
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
)

func (t *Topology) GetBizTopology(kit *rest.Kit, biz int64) (*string, error) {
	// read from secondary in mongodb cluster.
	kit.Ctx = util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	topology, err := t.briefBizKey.getTopology(kit.Ctx, biz)
	if err == nil {
		if len(*topology) != 0 {
			// get data from cache success
			return topology, nil
		}
		// get from db directly.
	}

	blog.Errorf("get biz: %d topology from cache failed, get from db now, err: %v, rid: %s", biz, err, kit.Rid)

	// do not get biz topology from cache, get it from db directly.
	topo, err := t.genBusinessTopology(kit.Ctx, biz)
	if err != nil {
		blog.Errorf("generate biz: %d topology from db failed, err: %v, rid: %s", biz, err, kit.Rid)
		return nil, err
	}

	// update it to cache directly.
	if err := t.briefBizKey.updateTopology(kit.Ctx, topo); err != nil {
		blog.Errorf("refresh biz: %d topology cache failed, err: %v, rid: %s", biz, err, kit.Rid)
		// do not return error
	}

	dat, err := json.Marshal(topo)
	if err != nil {
		blog.Errorf("marshal biz topology failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	topoStr := string(dat)
	return &topoStr, nil
}
