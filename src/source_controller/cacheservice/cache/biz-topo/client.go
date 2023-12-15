/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package biztopo

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/key"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/topo"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
)

// GetBizTopo get business topology cache info by topology type
func (t *Topo) GetBizTopo(kit *rest.Kit, typ string, opt *types.GetBizTopoOption) (*string, error) {
	// read from secondary in mongodb cluster.
	kit.Ctx = util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	topoType := types.TopoType(typ)
	topoKey, exists := key.TopoKeyMap[topoType]
	if !exists {
		blog.Errorf("biz topo type %s is invalid, rid: %s", topoType, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "type")
	}

	if opt == nil || opt.BizID <= 0 {
		blog.Errorf("get biz topo option %+v is invalid, rid: %s", opt, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "opt")
	}

	topology, err := topoKey.GetTopology(kit.Ctx, opt.BizID)
	if err == nil {
		if len(*topology) != 0 {
			// get data from cache succeed
			return topology, nil
		}
		// get from db directly.
	}

	blog.Errorf("get biz: %d %s topology from cache failed, get from db now, err: %v, rid: %s", opt.BizID, topoType,
		err, kit.Rid)

	// do not get biz topology from cache, get it from db directly.
	bizTopo, err := topo.GenBizTopo(kit.Ctx, opt.BizID, topoType, false, kit.Rid)
	if err != nil {
		blog.Errorf("generate biz: %d %s topology from db failed, err: %v, rid: %s", opt.BizID, topoType, err, kit.Rid)
		return nil, err
	}

	// update it to cache directly.
	topology, err = topoKey.UpdateTopology(kit.Ctx, bizTopo)
	if err != nil {
		blog.Errorf("update biz: %d %s topology cache failed, err: %v, rid: %s", opt.BizID, topoType, err, kit.Rid)
		// do not return error
	}

	return topology, nil
}
