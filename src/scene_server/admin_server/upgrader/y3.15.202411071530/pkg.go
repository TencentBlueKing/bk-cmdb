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

package y3_15_202411071530

import (
	"fmt"

	"configcenter/pkg/tenant/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/scene_server/admin_server/upgrader/y3.15.202411071530/data"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/driver/mongodb"
)

func init() {
	upgrader.RegisterUpgrade("y3.15.202411071530", upgrade)
}

func upgrade(kit *rest.Kit, db dal.Dal) error {

	dbCli, dbUUID, err := logics.GetNewTenantCli(kit, db)
	if err != nil {
		blog.Errorf("get new tenant db failed, rid: %s", kit.Rid)
		return fmt.Errorf("get new tenant db failed")
	}
	if err := initTableIndex(kit, dbCli, tableIndexMap); err != nil {
		blog.Errorf("init table index failed: err: %v", err)
		return err
	}

	if err := data.InitData(kit, dbCli); err != nil {
		blog.Errorf("add init data failed, err: %v", err)
		return err
	}

	cond := map[string]interface{}{
		common.TenantID: kit.TenantID,
	}
	count, err := mongodb.Dal().Shard(kit.SysShardOpts()).Table(common.BKTableNameTenant).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("get tenant count failed, err: %v", err)
		return err
	}
	if count > 0 {
		return nil
	}

	// add tenant, system or default
	err = mongodb.Dal().Shard(kit.SysShardOpts()).Table(common.BKTableNameTenant).Insert(kit.Ctx, types.Tenant{
		TenantID: kit.TenantID,
		Status:   types.EnabledStatus,
		Database: dbUUID,
	})
	if err != nil {
		blog.Errorf("add tenant failed for migrate data, tenantID: %s, err: %v", kit.TenantID, err)
		return err
	}
	return nil
}
