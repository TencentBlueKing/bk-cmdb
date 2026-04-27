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

package y3_14_202604271012

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// updateHostBkCPUArchitectureAttr add bk_cpu_architecture attribute for host
func updateHostBkCPUArchitectureAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	attrFilter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: "bk_cpu_architecture",
		common.BkSupplierAccount: conf.OwnerID,
	}

	cnt, err := db.Table(common.BKTableNameObjAttDes).Find(attrFilter).Count(ctx)
	if err != nil {
		blog.Errorf("count bk_cpu_architecture attribute failed, err: %v", err)
		return err
	}

	if cnt <= 0 {
		return fmt.Errorf("must have bk_cpu_architecture attribute")
	}

	updateData := bson.M{
		"option": []EnumVal{
			{ID: "x86", Name: "X86", Type: "text", IsDefault: true},
			{ID: "x86_64", Name: "X86_64", Type: "text"},
			{ID: "arm", Name: "ARM", Type: "text"},
			{ID: "aarch64", Name: "AARCH64", Type: "text"},
			{ID: "powerpc", Name: "POWERPC", Type: "text"},
			{ID: "ppc64", Name: "PPC64", Type: "text"},
			{ID: "ppc", Name: "PPC", Type: "text"},
			{ID: "riscv64", Name: "RISCV64", Type: "text"},
		},
	}
	m, err := db.Table(common.BKTableNameObjAttDes).UpdateMany(ctx, attrFilter, updateData)
	if err != nil {
		blog.Errorf("update bk_cpu_architecture attribute failed, err: %v", err)
		return err
	}
	blog.Errorf("update bk_cpu_architecture attribute %v", m)

	return nil
}

// EnumVal TODO
type EnumVal struct {
	ID        string `bson:"id"           json:"id"`
	Name      string `bson:"name"         json:"name"`
	Type      string `bson:"type"         json:"type"`
	IsDefault bool   `bson:"is_default"   json:"is_default"`
}
