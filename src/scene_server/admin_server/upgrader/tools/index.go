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

package tools

import (
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	idx "configcenter/src/common/index"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

// CreateIndexes create table index
func CreateIndexes(kit *rest.Kit, db dal.Dal, table string, indexes []types.Index) error {
	existIndexes, err := db.Shard(kit.ShardOpts()).Table(table).Indexes(kit.Ctx)
	if err != nil {
		blog.Errorf("get %s table exist index failed, err: %v", table, err)
		return err
	}

	existIndexMap := make(map[string]struct{})
	for _, index := range existIndexes {
		existIndexMap[index.Name] = struct{}{}
	}

	insertIndexes := make([]types.Index, 0)
	for _, index := range indexes {
		if _, exist := existIndexMap[index.Name]; exist {
			continue
		}

		isExist := false
		for _, existIndex := range existIndexes {
			if idx.IndexEqual(existIndex, index) {
				isExist = true
				break
			}
		}
		if isExist {
			continue
		}
		insertIndexes = append(insertIndexes, index)
	}

	if len(insertIndexes) == 0 {
		blog.Infof("table %s index is up to date", table)
		return nil
	}
	if err = db.Shard(kit.ShardOpts()).Table(table).BatchCreateIndexes(kit.Ctx, insertIndexes); err != nil {
		blog.Errorf("create %s table index %+v failed, err: %v", table, insertIndexes, err)
		return err
	}

	return nil
}
