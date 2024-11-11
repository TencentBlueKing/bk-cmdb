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

package kube

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
	"configcenter/src/storage/driver/mongodb"
)

// ListNsSharedClusterRel list namespace and shared cluster relations.
func (s *service) ListNsSharedClusterRel(cts *rest.Contexts) {
	opt := new(metadata.QueryCondition)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if key, err := opt.Page.Validate(true); err != nil {
		blog.Errorf("opt(%+v) is invalid, rid: %v", opt, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, key))
		return
	}

	filter := opt.Condition

	if opt.Page.EnableCount {
		count, err := mongodb.Client().Table(types.BKTableNameNsSharedClusterRel).Find(filter).Count(cts.Kit.Ctx)
		if err != nil {
			blog.Errorf("count ns shared cluster rel failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
			cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}

		cts.RespEntity(types.NsSharedClusterRelData{Count: count})
		return
	}

	limit := 0
	if opt.Page.Limit != common.BKNoLimit {
		limit = opt.Page.Limit
	}

	relations := make([]types.NsSharedClusterRel, 0)
	err := mongodb.Client().Table(types.BKTableNameNsSharedClusterRel).Find(filter).Start(uint64(opt.Page.Start)).
		Limit(uint64(limit)).Sort(opt.Page.Sort).Fields(opt.Fields...).All(cts.Kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("list ns shared cluster rel failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	cts.RespEntity(types.NsSharedClusterRelData{Info: relations})
}
