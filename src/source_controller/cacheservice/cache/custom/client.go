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

package custom

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/custom/cache"
	"configcenter/src/source_controller/cacheservice/cache/custom/types"
)

// ListPodLabelKey list pod label keys cache info
func (c *Cache) ListPodLabelKey(kit *rest.Kit, opt *types.ListPodLabelKeyOption) ([]string, error) {
	// read from secondary in mongodb cluster.
	kit.Ctx = util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	if opt == nil || opt.BizID <= 0 {
		blog.Errorf("list pod label key option %+v is invalid, rid: %s", opt, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "opt")
	}

	return c.cacheSet.Label.GetKeys(kit.Ctx, opt.BizID, kit.Rid)
}

// ListPodLabelValue list pod label values cache info
func (c *Cache) ListPodLabelValue(kit *rest.Kit, opt *types.ListPodLabelValueOption) ([]string, error) {
	// read from secondary in mongodb cluster.
	kit.Ctx = util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	if opt == nil || opt.BizID <= 0 || opt.Key == "" {
		blog.Errorf("list pod label value option %+v is invalid, rid: %s", opt, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "opt")
	}

	return c.cacheSet.Label.GetValues(kit.Ctx, opt.BizID, opt.Key, kit.Rid)
}

// RefreshPodLabel refresh biz pod label key and value cache
func (c *Cache) RefreshPodLabel(kit *rest.Kit, opt *types.RefreshPodLabelOption) error {
	// read from secondary in mongodb cluster.
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	if opt == nil || opt.BizID <= 0 {
		blog.Errorf("refresh pod label option %+v is invalid, rid: %s", opt, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "opt")
	}

	refreshOpt := &cache.RefreshPodLabelOpt{BizID: opt.BizID}

	go func() {
		blog.Infof("start refresh biz: %d pod label cache, rid: %s", opt.BizID, kit.Rid)
		_, err := c.cacheSet.Label.RefreshPodLabel(ctx, refreshOpt, kit.Rid)
		if err != nil {
			blog.Errorf("refresh biz: %d pod label cache failed, err: %v, rid: %s", opt.BizID, err, kit.Rid)
			return
		}
		blog.Infof("refresh biz: %d pod label cache successfully, rid: %s", opt.BizID, kit.Rid)
	}()

	return nil
}
