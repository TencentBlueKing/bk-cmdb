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

package tools

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/kube/types"
	"configcenter/src/storage/driver/mongodb"
)

// GenKubeSharedNsCond generate shared namespace condition by biz id
func GenKubeSharedNsCond(kit *rest.Kit, bizID int64, nsIDField string) (mapstr.MapStr, error) {
	kit.Ctx = util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	sharedCond := mapstr.MapStr{types.BKAsstBizIDField: bizID}

	relations := make([]types.NsSharedClusterRel, 0)
	err := mongodb.Shard(kit.ShardOpts()).Table(types.BKTableNameNsSharedClusterRel).Find(sharedCond).
		Fields(types.BKNamespaceIDField).All(kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("list kube shared namespace rel failed, err: %v, cond: %+v, rid: %v", err, sharedCond, kit.Rid)
		return nil, err
	}

	if len(relations) == 0 {
		return mapstr.MapStr{types.BKBizIDField: bizID}, nil
	}

	sharedNsIDs := make([]int64, 0)
	for _, relation := range relations {
		sharedNsIDs = append(sharedNsIDs, relation.NamespaceID)
	}

	return mapstr.MapStr{
		common.BKDBOR: []mapstr.MapStr{
			{types.BKBizIDField: bizID},
			{nsIDField: mapstr.MapStr{common.BKDBIN: sharedNsIDs}},
		},
	}, nil
}
