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

package y3_9_202410101500

import (
	"context"
	"reflect"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

var (
	newIndexMap = map[string]types.Index{
		common.BKTableNameBaseHost: {
			Name: "bkcc_unique_bkHostInnerIP_bkCloudID",
			Keys: map[string]int32{
				common.BKCloudIDField: 1, common.BKHostInnerIPField: 1,
			},
			Background: true,
			Unique:     true,
		},
		common.BKTableNameBaseInst: {
			Name: "bkcc_unique_bkObjectID_bkInstName",
			Keys: map[string]int32{
				common.BKObjIDField: 1, common.BKInstNameField: 1,
			},
			Background: true,
			Unique:     true,
		},
	}
)

// updateTableIndex update cc_HostBase and cc_ObjectBase table index
func updateTableIndex(ctx context.Context, db dal.RDB) error {
	for table, newIndex := range newIndexMap {
		existIndexes, err := db.Table(table).Indexes(ctx)
		if err != nil {
			return err
		}

		shouldCreate := true

		for _, idx := range existIndexes {
			if idx.Name == newIndex.Name {
				shouldCreate = false
				break
			}
			// 判断索引是否与目标索引键相同
			if !reflect.DeepEqual(idx.Keys, newIndex.Keys) {
				continue
			}
			// 判断索引是否是唯一索引
			if idx.Unique {
				shouldCreate = false
				break
			}
			// 如果索引键匹配但非唯一，删除该索引
			if err := db.Table(table).DropIndex(ctx, idx.Name); err != nil {
				blog.Errorf("drop table index failed, table: %s, index: %s, err: %v", table, idx.Name, err)
				return err
			}
		}

		if !shouldCreate {
			continue
		}
		if err := db.Table(table).CreateIndex(ctx, newIndex); err != nil {
			blog.Errorf("create table index failed, table: %s, index: %s, err: %v", table, newIndex.Name, err)
			return err
		}
	}

	return nil
}
