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

package service

import (
	"configcenter/pkg/tenant"
	"configcenter/pkg/tenant/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/driver/mongodb"
)

// GetAllTenants get all tenants from db
func (s *coreService) GetAllTenants(ctx *rest.Contexts) {
	result := tenant.GetAllTenants()
	if len(result) == 0 {
		blog.Errorf("tenant is empty, rid: %s", ctx.Kit.Rid)
	}
	ctx.RespEntity(result)
}

// RefreshAllTenants refresh all tenants from db
func (s *coreService) RefreshAllTenants(ctx *rest.Contexts) {
	tenants := make([]types.Tenant, 0)
	err := mongodb.Shard(ctx.Kit.SysShardOpts()).Table(common.BKTableNameTenant).Find(mapstr.MapStr{}).All(ctx.Kit.Ctx,
		&tenants)
	if err != nil {
		blog.Errorf("find all tenants failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	tenant.SetTenant(tenants)
	// refresh tenant db map
	shardingMongoManager, ok := mongodb.Dal().(*sharding.ShardingMongoManager)
	if !ok {
		blog.Errorf("convert to ShardingMongoManager failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if err = shardingMongoManager.RefreshTenantDBMap(); err != nil {
		blog.Errorf("refresh tenant db map failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(tenants)
}
