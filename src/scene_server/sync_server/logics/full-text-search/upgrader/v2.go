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

package upgrader

import (
	"context"
)

func init() {
	RegisterUpgrader(2, upgraderInst.upgradeV2)
}

var v2Indexes = []string{
	"bk_cmdb.bk_biz_set_obj_20250317",
	"bk_cmdb.biz_20250317",
	"bk_cmdb.set_20250317",
	"bk_cmdb.module_20250317",
	"bk_cmdb.host_20250317",
	"bk_cmdb.model_20250317",
	"bk_cmdb.object_instance_20250317",
}

// upgradeV2 update tenant id field name from meta_bk_supplier_account to meta_tenant_id
// and update es id format from {_id}:{type suffix} to {id}
// since the data structure is very different from monstache structure, we use full sync instead of reindex to upgrade
func (u *upgrader) upgradeV2(ctx context.Context, rid string) (*UpgraderFuncResult, error) {
	return &UpgraderFuncResult{Indexes: v2Indexes, NeedSyncAll: true}, nil
}
