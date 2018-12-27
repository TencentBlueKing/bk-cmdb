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

package driver

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"configcenter/src/common"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/mongodb/options/findopt"
	"configcenter/src/storage/types"
)

// ErrSessionMissing session missing
var ErrSessionMissing = errors.New("session missing")

// Client implement client.DALRDB interface
type Client struct {
	dbProxy mongodb.Client
	session mongodb.Session
}

var _ dal.DB = (*Client)(nil)

var initMongoc sync.Once

// NewClient returns new RDB
func NewClient(uri string) (*Client, error) {
	/*
		pool := mongodb.NewClientPool(uri)
		err := pool.Open()
		if err != nil {
			return nil, err
		}
		return &Client{
			pool: pool,
		}, nil
	*/

	return nil, nil
}

// Close replica client
func (c *Client) Close() error {
	return nil
}

// Ping replica client
func (c *Client) Ping() error {
	return nil
}

// Clone return the new client
func (c *Client) Clone() dal.DB {
	return nil
}

// IsDuplicatedError interface compatibility
func (c *Client) IsDuplicatedError(err error) bool {
	return false
}

// IsNotFoundError interface compatibility
func (c *Client) IsNotFoundError(err error) bool {
	return false
}

// Table collection operation
func (c *Client) Table(collName string) dal.Table {
	col := Collection{}
	col.collName = collName
	col.Client = c
	return &col
}

// NextSequence 获取新序列号(非事务)
func (c *Client) NextSequence(ctx context.Context, sequenceName string) (uint64, error) {
	data := types.Document{
		"$inc": types.Document{"SequenceID": 1},
	}
	filter := types.Document{
		"_id": sequenceName,
	}

	opt := findopt.FindAndModify{}
	opt.Upsert = true
	opt.New = true

	results := types.Documents{}

	err := c.dbProxy.Collection(common.BKTableNameIDgenerator).FindOneAndModify(ctx, filter, data, &opt, &results)
	if nil != err {
		return 0, err
	}

	if len(results) <= 0 {
		return 0, dal.ErrDocumentNotFound
	}

	return strconv.ParseUint(fmt.Sprint(results[0]["SequenceID"]), 10, 64)
}

// StartTransaction 开启新事务
func (c *Client) StartTransaction(ctx context.Context) (dal.DB, error) {
	/*
		txc := c.pool.Pop()
		c.txc = txc
		session := txc.Session().Create()
		if err := session.Open(); err != nil {
			session.Close()
			return nil, err
		}
		c.session = session
		err := session.StartTransaction()
		if err != nil {
			session.Close()
		}
		return c, nil
	*/
	return nil, nil
}

// Commit 提交事务
func (c *Client) Commit(context.Context) error {
	/*
		if c.session == nil {
			return ErrSessionMissing
		}
		commitErr := c.session.CommitTransaction()
		if commitErr == nil {
			closeErr := c.session.Close()
			c.pool.Push(c.txc)
			if closeErr != nil {
				blog.Warnf("[mongoc dal] session close faile: %v", closeErr)
			}
			c.session = nil
			return nil
		}
		return commitErr
	*/
	return nil
}

// Abort 取消事务
func (c *Client) Abort(context.Context) error {
	/*
		if c.session == nil {
			return ErrSessionMissing
		}
		abortErr := c.session.AbortTransaction()
		closeErr := c.session.Close()
		c.pool.Push(c.txc)
		if closeErr != nil {
			blog.Warnf("[mongoc dal] session close faile: %v", closeErr)
		}
		c.session = nil
		return abortErr
	*/
	return nil
}

// TxnInfo 当前事务信息，用于事务发起者往下传递
func (c *Client) TxnInfo() *types.Transaction {
	return &types.Transaction{}
}

// HasTable 判断是否存在集合
func (c *Client) HasTable(collName string) (bool, error) {

	exists, err := c.dbProxy.Database().HasCollection(collName)
	return exists, err
}

// DropTable 移除集合
func (c *Client) DropTable(collName string) error {
	err := c.dbProxy.Database().DropCollection(collName)
	return err
}

// CreateTable 创建集合
func (c *Client) CreateTable(collName string) error {
	err := c.dbProxy.Database().CreateEmptyCollection(collName)
	return err
}

// CreateIndex 创建索引
func (c *Collection) CreateIndex(ctx context.Context, index dal.Index) error {
	i := mongodb.Index{
		Keys:       index.Keys,
		Name:       index.Name,
		Unique:     index.Unique,
		Background: index.Background,
	}

	err := c.dbProxy.Collection(c.collName).CreateIndex(i)
	return err
}

// DropIndex 移除索引
func (c *Collection) DropIndex(ctx context.Context, indexName string) error {
	err := c.dbProxy.Collection(c.collName).DropIndex(indexName)
	return err
}

// AddColumn 添加字段
func (c *Collection) AddColumn(ctx context.Context, column string, value interface{}) error {

	selector := types.Document{column: types.Document{"$exists": false}}
	datac := types.Document{"$set": types.Document{column: value}}

	_, err := c.dbProxy.Collection(c.collName).UpdateMany(ctx, selector, datac, nil)
	return err
}

// RenameColumn 重命名字段
func (c *Collection) RenameColumn(ctx context.Context, oldName, newColumn string) error {

	datac := types.Document{"$rename": types.Document{oldName: newColumn}}
	_, err := c.dbProxy.Collection(c.collName).UpdateMany(ctx, nil, datac, nil)
	return err
}

// DropColumn 移除字段
func (c *Collection) DropColumn(ctx context.Context, field string) error {

	datac := types.Document{"$unset": types.Document{field: "1"}}
	_, err := c.dbProxy.Collection(c.collName).UpdateMany(ctx, nil, datac, nil)
	return err
}

func (c *Client) getCollection(collName string) mongodb.CollectionInterface {
	var table mongodb.CollectionInterface
	if c.session == nil {
		table = c.dbProxy.Collection(collName)
	} else {
		table = c.session.Collection(collName)
	}
	return table
}
