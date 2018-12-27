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
	"errors"

	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

// Collection implement client.Collection interface
type Collection struct {
	collName string // 集合名
	*Client
}

// AggregateOne TODO: need to implement
func (c *Collection) AggregateOne(ctx context.Context, pipeline interface{}, result interface{}) error {
	return nil
}

// AggregateAll TODO: need to implement
func (c *Collection) AggregateAll(ctx context.Context, pipeline interface{}, result interface{}) error {
	return nil
}

// Indexes 查询索引
func (c *Collection) Indexes(ctx context.Context) ([]dal.Index, error) {
	return nil, errors.New("not implement")
}

// Find 查询多个并反序列化到 Result
func (c *Collection) Find(filter dal.Filter) dal.Find {
	return &Find{Collection: c, filter: filter}
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *Collection) Insert(ctx context.Context, docs interface{}) error {

	table := c.getCollection(c.collName)
	err := table.InsertMany(ctx, util.ConverToInterfaceSlice(docs), nil)
	return err
}

// Update 更新数据
func (c *Collection) Update(ctx context.Context, filter dal.Filter, doc interface{}) error {
	table := c.getCollection(c.collName)
	_, err := table.UpdateMany(ctx, filter, doc, nil)
	return err
}

// Delete 删除数据
func (c *Collection) Delete(ctx context.Context, filter dal.Filter) error {
	table := c.getCollection(c.collName)
	_, err := table.DeleteMany(ctx, filter, nil)
	return err
}
