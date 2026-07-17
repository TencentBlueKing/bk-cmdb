/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package y3_15_202607151050

import (
	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal"
)

// updateHostBkCPUArchitectureAttr update the option of host bk_cpu_architecture attribute for each tenant.
func updateHostBkCPUArchitectureAttr(kit *rest.Kit, dal dal.Dal) error {
	return tenant.ExecForAllTenants(func(tenantID string) error {
		kt := kit.NewKit().WithTenant(tenantID)
		db := dal.Shard(kt.ShardOpts())

		attrFilter := map[string]interface{}{
			common.BKObjIDField:      common.BKInnerObjIDHost,
			common.BKPropertyIDField: common.BKCpuArch,
		}
		attr := make(mapstr.MapStr)
		err := db.Table(common.BKTableNameObjAttDes).Find(attrFilter).Fields(common.BKFieldID).One(kt.Ctx, &attr)
		if err != nil {
			blog.Errorf("findOne bk_cpu_architecture attribute for tenant %s failed, err: %v", tenantID, err)
			return err
		}

		updateData := map[string]interface{}{
			common.BKOptionField: []enumVal{
				{ID: "x86", Name: "X86", Type: "text", IsDefault: true},
				{ID: "x86_64", Name: "X86_64", Type: "text"},
				{ID: "arm", Name: "ARM", Type: "text"},
				{ID: "arm64", Name: "ARM64", Type: "text"},
				{ID: "aarch64", Name: "AARCH64", Type: "text"},
				{ID: "powerpc", Name: "POWERPC", Type: "text"},
				{ID: "ppc64", Name: "PPC64", Type: "text"},
				{ID: "ppc", Name: "PPC", Type: "text"},
			},
		}
		updateFilter := map[string]interface{}{
			common.BKFieldID: attr[common.BKFieldID],
		}
		if err = db.Table(common.BKTableNameObjAttDes).Update(kt.Ctx, updateFilter, updateData); err != nil {
			blog.Errorf("update host bk_cpu_architecture attribute failed for tenant %s, err: %v", tenantID, err)
			return err
		}
		return nil
	})
}

// enumVal enum option value of the attribute
type enumVal struct {
	ID        string `bson:"id" json:"id"`
	Name      string `bson:"name" json:"name"`
	Type      string `bson:"type" json:"type"`
	IsDefault bool   `bson:"is_default" json:"is_default"`
}
