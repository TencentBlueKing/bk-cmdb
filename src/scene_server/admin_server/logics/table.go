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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
)

// CreateTable create table if not exists
func CreateTable(kit *rest.Kit, db local.DB, table string) error {
	if common.IsPlatformTable(table) {
		db = mongodb.Shard(kit.SysShardOpts())
	}
	exists, err := db.HasTable(kit.Ctx, table)
	if err != nil {
		blog.Errorf("check if %s exists failed, err: %v", table, err)
		return err
	}
	if exists {
		return nil
	}

	if err = db.CreateTable(kit.Ctx, table); err != nil {
		blog.Errorf("create %s table failed, err: %v", table, err)
		return err
	}

	return nil
}
