/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package driver

import (
	"context"

	"configcenter/src/storage/dal"
	"configcenter/src/storage/mongodb/options/findopt"
	"configcenter/src/storage/types"
)

// Find define a find operation
type Find struct {
	*Collection
	projection types.Document
	filter     dal.Filter
	start      uint64
	limit      uint64
	sort       string
}

// Fields 查询字段
func (f *Find) Fields(fields ...string) dal.Find {
	projection := types.Document{}
	for _, field := range fields {
		projection[field] = true
	}
	f.projection = projection
	return f
}

// Sort 查询排序
func (f *Find) Sort(sort string) dal.Find {
	f.sort = sort
	return f
}

// Start 查询上标
func (f *Find) Start(start uint64) dal.Find {
	f.start = start
	return f
}

// Limit 查询限制
func (f *Find) Limit(limit uint64) dal.Find {
	f.limit = limit
	return f
}

// All 查询多个
func (f *Find) All(ctx context.Context, result interface{}) error {
	opt := findopt.Many{}
	opt.Skip = int64(f.start)
	opt.Limit = int64(f.limit)
	//opt.Fields = mapstr.MapStr(f.projection)

	table := f.getCollection(f.collName)
	err := table.Find(ctx, f.filter, &opt, result)

	return err
}

// One 查询一个
func (f *Find) One(ctx context.Context, result interface{}) error {
	opt := findopt.One{}
	opt.Skip = int64(f.start)
	opt.Limit = int64(f.limit)
	//opt.Fields = mapstr.MapStr(f.projection)

	table := f.getCollection(f.collName)
	err := table.FindOne(ctx, f.filter, &opt, result)

	return err
}

// Count 统计数量(非事务)
func (f *Find) Count(ctx context.Context) (uint64, error) {

	return f.dbProxy.Collection(f.collName).Count(ctx, f.filter)
}
