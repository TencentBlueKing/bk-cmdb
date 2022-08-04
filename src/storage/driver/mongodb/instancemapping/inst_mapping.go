/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package instancemapping

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

/***
deprecated
不建议使用， 改功能是在通用模型的实例数据分表后， 只有实例id，没有bk_obj_id的时候使用，负责将实例id 转为bk_obj_id
使用前必须先初始化 cc 的mongodb driver（configcenter/src/storage/driver/mongodb）

mongodb 集群 配置：
节点数量:3
规格:2核；4GB；100GB(SSD)

索引表
document 只包含bk_obj_id, bk_inst_id 两个字段
索引， bk_inst_id 单个字段索引
数据量在1亿个document

插入行数：20w（单任务处理，没有并发）
耗时： 3m1.183887155s， 每次数据库耗时1ms


查询次数：100w	 每次查询ID个数：1	耗时：17m51.628125676s	QPS: 933
***/

var (
	tableName = "cc_ObjectBaseMapping"
)

// GetInstanceMapping TODO
// deprecated 不建议使用，新加的要求用户必须传bk_obj_id, 改功能是在实例数据分表后， 只有实例id，没有bk_obj_id的时候使用，
//     负责将实例id 转为bk_obj_id,
func GetInstanceMapping(ids []int64) (map[int64]metadata.ObjectMapping, error) {
	if len(ids) > 200 {
		return nil, fmt.Errorf("id array count must lt 200")
	}
	filter := map[string]interface{}{
		common.BKInstIDField: map[string]interface{}{
			common.BKDBIN: ids,
		},
	}
	rows := make([]metadata.ObjectMapping, 0)
	// 看不到事务中未提交的数据
	if err := mongodb.Table(tableName).Find(filter).All(context.Background(), &rows); err != nil {
		return nil, err
	}

	mapping := make(map[int64]metadata.ObjectMapping, 0)
	for _, row := range rows {
		mapping[row.ID] = row
	}

	return mapping, nil

}

// GetInstanceObjectMapping TODO
func GetInstanceObjectMapping(ids []int64) ([]metadata.ObjectMapping, error) {
	mapping := make([]metadata.ObjectMapping, 0)
	total := len(ids)
	if total == 0 {
		return mapping, nil
	}

	const step = 500
	for start := 0; start < total; start += step {
		var paged []int64
		if total-start >= step {
			paged = ids[start : start+step]
		} else {
			paged = ids[start:total]
		}

		filter := map[string]interface{}{
			common.BKInstIDField: map[string]interface{}{
				common.BKDBIN: paged,
			},
		}
		rows := make([]metadata.ObjectMapping, 0)
		// 看不到事务中未提交的数据
		if err := mongodb.Table(tableName).Find(filter).All(context.Background(), &rows); err != nil {
			return nil, err
		}

		mapping = append(mapping, rows...)
	}

	return mapping, nil
}

// Create 新加实例id与模型id的对应关系就， ctx 是为了保证事务， doc 为数组的时候表示插入多条数据
func Create(ctx context.Context, doc interface{}) error {
	return mongodb.Table(tableName).Insert(ctx, doc)
}

// Delete TODO
//  移除实例id与模型id的对应关系就， ctx 是为了保证事务
func Delete(ctx context.Context, ids []int64) error {
	filter := map[string]interface{}{
		common.BKInstIDField: map[string]interface{}{
			common.BKDBIN: ids,
		},
	}
	return mongodb.Table(tableName).Delete(ctx, filter)
}
