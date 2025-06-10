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

package metadata

import (
	"context"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/transfer-service/app/options"
	"configcenter/src/storage/driver/mongodb"
)

// BluekingBizID is the blueking biz info, resource in blueking biz should not be synced
type bluekingBizInfo struct {
	bizID         int64
	tenantID      string
	lock          sync.RWMutex
	hostModuleMap map[int64]map[int64]struct{}
}

// initBluekingBizInfo get blueking biz id and host ids for source environment
func (m *Metadata) initBluekingBizInfo(ctx context.Context) error {
	kit := rest.NewKit().WithCtx(ctx)

	if m.role != options.SyncRoleSrc {
		return nil
	}
	m.blueking = &bluekingBizInfo{
		tenantID: kit.TenantID,
	}

	bluekingBiz := new(metadata.BizBasicInfo)
	bluekingBizCond := mapstr.MapStr{common.BKAppNameField: common.BKAppName}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseApp).Find(bluekingBizCond).
		Fields(common.BKAppIDField).One(ctx, &bluekingBiz); err != nil {
		blog.Errorf("get blueking biz id by cond(%+v) failed, err: %v", bluekingBizCond, err)
		return err
	}
	m.blueking.bizID = bluekingBiz.BizID

	bkHostRel := make([]metadata.ModuleHost, 0)
	bkHostRelCond := mapstr.MapStr{common.BKAppIDField: m.blueking.bizID}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModuleHostConfig).Find(bkHostRelCond).
		Fields(common.BKHostIDField, common.BKModuleIDField).All(ctx, &bkHostRel); err != nil {
		blog.Errorf("get blueking host relations by cond(%+v) failed, err: %v", bkHostRelCond, err)
		return err
	}

	m.blueking.hostModuleMap = make(map[int64]map[int64]struct{})
	for _, relation := range bkHostRel {
		_, exists := m.blueking.hostModuleMap[relation.HostID]
		if !exists {
			m.blueking.hostModuleMap[relation.HostID] = make(map[int64]struct{})
		}
		m.blueking.hostModuleMap[relation.HostID][relation.ModuleID] = struct{}{}
	}
	return nil
}

func (m *Metadata) needCheckBluekingBiz(kit *rest.Kit) bool {
	return m.blueking != nil && kit.TenantID == m.blueking.tenantID && m.blueking.bizID != 0
}

func (m *Metadata) getBluekingHostIDs(kit *rest.Kit) []int64 {
	if m.blueking == nil || kit.TenantID != m.blueking.tenantID {
		return make([]int64, 0)
	}

	m.blueking.lock.RLock()
	defer m.blueking.lock.RUnlock()

	hostIDs := make([]int64, 0)
	for hostID := range m.blueking.hostModuleMap {
		hostIDs = append(hostIDs, hostID)
	}
	return hostIDs
}

func (m *Metadata) isHostInBluekingBiz(kit *rest.Kit, hostID int64) bool {
	if m.blueking == nil || kit.TenantID != m.blueking.tenantID {
		return false
	}

	m.blueking.lock.RLock()
	defer m.blueking.lock.RUnlock()

	_, exists := m.blueking.hostModuleMap[hostID]
	return exists
}
