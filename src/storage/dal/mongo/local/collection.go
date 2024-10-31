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

package local

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/util/table"
	"configcenter/src/storage/dal/types"
	dtype "configcenter/src/storage/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Collection implement client.Collection interface
type Collection struct {
	collName string // 集合名
	*Mongo
}

// Find 查询多个并反序列化到 Result
func (c *Collection) Find(filter types.Filter, opts ...*types.FindOpts) types.Find {
	find := &Find{
		Collection: c,
		filter:     filter,
		projection: make(map[string]int),
	}

	find.Option(opts...)

	return find
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *Collection) Insert(ctx context.Context, docs interface{}) error {
	mtc.collectOperCount(c.collName, insertOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, insertOper, time.Since(start))
	}()

	rows := util.ConvertToInterfaceSlice(docs)

	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).InsertMany(ctx, rows)
		if err != nil {
			mtc.collectErrorCount(c.collName, insertOper)
			return err
		}

		return nil
	})
}

// Update 更新数据
func (c *Collection) Update(ctx context.Context, filter types.Filter, doc interface{}) error {
	mtc.collectOperCount(c.collName, updateOper)
	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, updateOper, time.Since(start))
	}()

	if filter == nil {
		filter = bson.M{}
	}

	data := bson.M{"$set": doc}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, data)
		if err != nil {
			mtc.collectErrorCount(c.collName, updateOper)
			return err
		}
		return nil
	})
}

// UpdateMany TODO
// Update 更新数据, 返回修改成功的条数
func (c *Collection) UpdateMany(ctx context.Context, filter types.Filter, doc interface{}) (uint64, error) {
	mtc.collectOperCount(c.collName, updateOper)
	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, updateOper, time.Since(start))
	}()

	if filter == nil {
		filter = bson.M{}
	}

	data := bson.M{"$set": doc}
	var modifiedCount uint64
	err := c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		updateRet, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, data)
		if err != nil {
			mtc.collectErrorCount(c.collName, updateOper)
			return err
		}
		modifiedCount = uint64(updateRet.ModifiedCount)
		return nil
	})
	return modifiedCount, err
}

// Upsert 数据存在更新数据，否则新加数据。
// 注意：该接口非原子操作，可能存在插入多条相同数据的风险。
func (c *Collection) Upsert(ctx context.Context, filter types.Filter, doc interface{}) error {
	mtc.collectOperCount(c.collName, upsertOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, upsertOper, time.Since(start))
	}()

	// set upsert option
	doUpsert := true
	replaceOpt := &options.UpdateOptions{
		Upsert: &doUpsert,
	}
	data := bson.M{"$set": doc}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateOne(ctx, filter, data, replaceOpt)
		if err != nil {
			mtc.collectErrorCount(c.collName, upsertOper)
			return err
		}
		return nil
	})

}

// UpdateMultiModel 根据不同的操作符去更新数据
func (c *Collection) UpdateMultiModel(ctx context.Context, filter types.Filter, updateModel ...types.ModeUpdate) error {
	mtc.collectOperCount(c.collName, updateOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, updateOper, time.Since(start))
	}()

	data := bson.M{}
	for _, item := range updateModel {
		if _, ok := data[item.Op]; ok {
			return errors.New(item.Op + " appear multiple times")
		}
		data["$"+item.Op] = item.Doc
	}

	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, data)
		if err != nil {
			mtc.collectErrorCount(c.collName, updateOper)
			return err
		}
		return nil
	})

}

// Delete 删除数据
func (c *Collection) Delete(ctx context.Context, filter types.Filter) error {
	_, err := c.DeleteMany(ctx, filter)
	return err
}

// DeleteMany TODO
// Delete 删除数据， 返回删除的行数
func (c *Collection) DeleteMany(ctx context.Context, filter types.Filter) (uint64, error) {
	mtc.collectOperCount(c.collName, deleteOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, deleteOper, time.Since(start))
	}()

	var deleteCount uint64
	err := c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		if err := c.tryArchiveDeletedDoc(ctx, filter); err != nil {
			mtc.collectErrorCount(c.collName, deleteOper)
			return err
		}
		deleteRet, err := c.dbc.Database(c.dbname).Collection(c.collName).DeleteMany(ctx, filter)
		if err != nil {
			mtc.collectErrorCount(c.collName, deleteOper)
			return err
		}

		deleteCount = uint64(deleteRet.DeletedCount)
		return nil
	})

	return deleteCount, err
}

func (c *Collection) tryArchiveDeletedDoc(ctx context.Context, filter types.Filter) error {
	delArchiveTable, exists := table.GetDelArchiveTable(c.collName)
	if !exists {
		// do not archive the delete docs
		return nil
	}

	// only archive the specified fields for delete docs
	var findOpts *options.FindOptions
	fields := table.GetDelArchiveFields(c.collName)
	if len(fields) > 0 {
		projection := map[string]int{"_id": 1}
		for _, field := range fields {
			projection[field] = 1
		}
		findOpts = &options.FindOptions{Projection: projection}
	}

	docs := make([]bson.D, 0)
	cursor, err := c.dbc.Database(c.dbname).Collection(c.collName).Find(ctx, filter, findOpts)
	if err != nil {
		return err
	}

	if err := cursor.All(ctx, &docs); err != nil {
		return err
	}

	if len(docs) == 0 {
		return nil
	}

	archives := make([]interface{}, len(docs))
	for idx, doc := range docs {
		detail := make(bson.D, 0)
		var oid string
		for _, e := range doc {
			if e.Key == "_id" {
				rawOid, ok := e.Value.(primitive.ObjectID)
				if !ok {
					return errors.New("invalid object id")
				}
				oid = rawOid.Hex()
				continue
			}
			detail = append(detail, e)
		}
		archives[idx] = metadata.DeleteArchive{
			Oid:    oid,
			Detail: detail,
			Time:   time.Now(),
			Coll:   c.collName,
		}
	}

	_, err = c.dbc.Database(c.dbname).Collection(delArchiveTable).InsertMany(ctx, archives)
	return err
}

// BatchCreateIndexes 批量创建索引
func (c *Collection) BatchCreateIndexes(ctx context.Context, indexes []types.Index) error {
	mtc.collectOperCount(c.collName, indexCreateOper)

	createIndexInfos := make([]mongo.IndexModel, len(indexes))
	for idx, index := range indexes {
		createIndexInfo, err := buildIndex(index)
		if err != nil {
			return err
		}

		createIndexInfos[idx] = createIndexInfo
	}

	indexView := c.dbc.Database(c.dbname).Collection(c.collName).Indexes()
	_, err := indexView.CreateMany(ctx, createIndexInfos)
	if err != nil {
		mtc.collectErrorCount(c.collName, indexCreateOper)
		// ignore the following case
		// 1.the new index is exactly the same as the existing one
		// 2.the new index has same keys with the existing one, but its name is different
		if strings.Contains(err.Error(), "all indexes already exist") ||
			strings.Contains(err.Error(), "already exists with a different name") {
			return nil
		}
	}

	return err
}

// CreateIndex 创建索引
func (c *Collection) CreateIndex(ctx context.Context, index types.Index) error {
	mtc.collectOperCount(c.collName, indexCreateOper)

	createIndexInfo, err := buildIndex(index)
	if err != nil {
		return err
	}

	indexView := c.dbc.Database(c.dbname).Collection(c.collName).Indexes()
	_, err = indexView.CreateOne(ctx, createIndexInfo)
	if err != nil {
		mtc.collectErrorCount(c.collName, indexCreateOper)
		// ignore the following case
		// 1.the new index is exactly the same as the existing one
		// 2.the new index has same keys with the existing one, but its name is different
		if strings.Contains(err.Error(), "all indexes already exist") ||
			strings.Contains(err.Error(), "already exists with a different name") {
			return nil
		}
	}

	return err
}

func buildIndex(index types.Index) (mongo.IndexModel, error) {
	createIndexOpt := &options.IndexOptions{
		Background:              &index.Background,
		Unique:                  &index.Unique,
		PartialFilterExpression: index.PartialFilterExpression,
	}
	if index.Name != "" {
		createIndexOpt.Name = &index.Name
	}

	if index.ExpireAfterSeconds != 0 {
		createIndexOpt.SetExpireAfterSeconds(index.ExpireAfterSeconds)
	}

	keys := index.Keys
	for idx, key := range keys {
		val, err := util.GetInt32ByInterface(key.Value)
		if err != nil {
			return mongo.IndexModel{}, err
		}
		key.Value = val
		keys[idx] = key
	}

	return mongo.IndexModel{
		Keys:    keys,
		Options: createIndexOpt,
	}, nil
}

// DropIndex remove index by name
func (c *Collection) DropIndex(ctx context.Context, indexName string) error {
	mtc.collectOperCount(c.collName, indexDropOper)
	indexView := c.dbc.Database(c.dbname).Collection(c.collName).Indexes()
	_, err := indexView.DropOne(ctx, indexName)
	if err != nil {
		if strings.Contains(err.Error(), "IndexNotFound") {
			return nil
		}
		mtc.collectErrorCount(c.collName, indexDropOper)
		return err
	}
	return nil
}

// Indexes get all indexes for the collection
func (c *Collection) Indexes(ctx context.Context) ([]types.Index, error) {
	indexView := c.dbc.Database(c.dbname).Collection(c.collName).Indexes()
	cursor, err := indexView.List(ctx)
	if nil != err {
		return nil, err
	}
	defer cursor.Close(ctx)
	var indexes []types.Index
	for cursor.Next(ctx) {
		idxResult := types.Index{}
		cursor.Decode(&idxResult)
		indexes = append(indexes, idxResult)
	}

	return indexes, nil
}

// AddColumn add a new column for the collection
func (c *Collection) AddColumn(ctx context.Context, column string, value interface{}) error {
	mtc.collectOperCount(c.collName, columnOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, columnOper, time.Since(start))
	}()

	selector := dtype.Document{column: dtype.Document{"$exists": false}}
	datac := dtype.Document{"$set": dtype.Document{column: value}}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, selector, datac)
		if err != nil {
			mtc.collectErrorCount(c.collName, columnOper)
			return err
		}
		return nil
	})
}

// RenameColumn rename a column for the collection
func (c *Collection) RenameColumn(ctx context.Context, filter types.Filter, oldName, newColumn string) error {
	mtc.collectOperCount(c.collName, columnOper)
	if filter == nil {
		filter = dtype.Document{}
	}

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, columnOper, time.Since(start))
	}()

	datac := dtype.Document{"$rename": dtype.Document{oldName: newColumn}}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, datac)
		if err != nil {
			mtc.collectErrorCount(c.collName, columnOper)
			return err
		}

		return nil
	})
}

// DropColumn remove a column by the name
func (c *Collection) DropColumn(ctx context.Context, field string) error {
	mtc.collectOperCount(c.collName, columnOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, columnOper, time.Since(start))
	}()

	datac := dtype.Document{"$unset": dtype.Document{field: ""}}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, dtype.Document{}, datac)
		if err != nil {
			mtc.collectErrorCount(c.collName, columnOper)
			return err
		}

		return nil
	})
}

// DropColumns remove many columns by the name
func (c *Collection) DropColumns(ctx context.Context, filter types.Filter, fields []string) error {
	mtc.collectOperCount(c.collName, columnOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, columnOper, time.Since(start))
	}()

	unsetFields := make(map[string]interface{})
	for _, field := range fields {
		unsetFields[field] = ""
	}

	datac := dtype.Document{"$unset": unsetFields}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, datac)
		if err != nil {
			mtc.collectErrorCount(c.collName, columnOper)
			return err
		}

		return nil
	})
}

// DropDocsColumn remove a column by the name for doc use filter
func (c *Collection) DropDocsColumn(ctx context.Context, field string, filter types.Filter) error {
	mtc.collectOperCount(c.collName, columnOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, columnOper, time.Since(start))
	}()

	// 查询条件为空时候，mongodb 不返回数据
	if filter == nil {
		filter = bson.M{}
	}

	datac := dtype.Document{"$unset": dtype.Document{field: ""}}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, datac)
		if err != nil {
			mtc.collectErrorCount(c.collName, columnOper)
			return err
		}

		return nil
	})
}

// AggregateAll aggregate all operation
func (c *Collection) AggregateAll(ctx context.Context, pipeline interface{}, result interface{},
	opts ...*types.AggregateOpts) error {

	mtc.collectOperCount(c.collName, aggregateOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, aggregateOper, time.Since(start))
	}()

	var aggregateOption *options.AggregateOptions
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.AllowDiskUse != nil {
			aggregateOption = &options.AggregateOptions{AllowDiskUse: opt.AllowDiskUse}
		}
	}

	opt := getCollectionOption(ctx)

	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		cursor, err := c.dbc.Database(c.dbname).Collection(c.collName, opt).Aggregate(ctx, pipeline, aggregateOption)
		if err != nil {
			mtc.collectErrorCount(c.collName, aggregateOper)
			return err
		}
		defer cursor.Close(ctx)
		return decodeCursorIntoSlice(ctx, cursor, result)
	})

}

// AggregateOne aggregate one operation
func (c *Collection) AggregateOne(ctx context.Context, pipeline interface{}, result interface{}) error {
	mtc.collectOperCount(c.collName, aggregateOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, aggregateOper, time.Since(start))
	}()

	opt := getCollectionOption(ctx)

	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		cursor, err := c.dbc.Database(c.dbname).Collection(c.collName, opt).Aggregate(ctx, pipeline)
		if err != nil {
			mtc.collectErrorCount(c.collName, aggregateOper)
			return err
		}

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			return cursor.Decode(result)
		}
		return types.ErrDocumentNotFound
	})

}

// Distinct Finds the distinct values for a specified field in a single collection or view and returns the result array
// field the field for which to return distinct values.
// filter query that specifies the documents from which to retrieve the distinct values.
func (c *Collection) Distinct(ctx context.Context, field string, filter types.Filter) ([]interface{}, error) {
	mtc.collectOperCount(c.collName, distinctOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, distinctOper, time.Since(start))
	}()

	if filter == nil {
		filter = bson.M{}
	}

	opt := getCollectionOption(ctx)
	var results []interface{} = nil
	err := c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		var err error
		results, err = c.dbc.Database(c.dbname).Collection(c.collName, opt).Distinct(ctx, field, filter)
		if err != nil {
			mtc.collectErrorCount(c.collName, distinctOper)
			return err
		}

		return nil
	})
	return results, err
}

func decodeCursorIntoSlice(ctx context.Context, cursor *mongo.Cursor, result interface{}) error {
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		return errors.New("result argument must be a slice address")
	}

	elemt := resultv.Elem().Type().Elem()
	slice := reflect.MakeSlice(resultv.Elem().Type(), 0, 10)
	for cursor.Next(ctx) {
		elemp := reflect.New(elemt)
		if err := cursor.Decode(elemp.Interface()); nil != err {
			return err
		}
		slice = reflect.Append(slice, elemp.Elem())
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	resultv.Elem().Set(slice)
	return nil
}

const (
	// reference doc:
	// https://docs.mongodb.com/manual/core/read-preference-staleness/#replica-set-read-preference-max-staleness
	// this is the minimum value of maxStalenessSeconds allowed.
	// specifying a smaller maxStalenessSeconds value will raise an error. Clients estimate secondaries’ staleness
	// by periodically checking the latest write date of each replica set member. Since these checks are infrequent,
	// the staleness estimate is coarse. Thus, clients cannot enforce a maxStalenessSeconds value of less than
	// 90 seconds.
	maxStalenessSeconds = 90 * time.Second
)

func getCollectionOption(ctx context.Context) *options.CollectionOptions {
	var opt *options.CollectionOptions
	switch util.GetDBReadPreference(ctx) {

	case common.NilMode:

	case common.PrimaryMode:
		opt = &options.CollectionOptions{
			ReadPreference: readpref.Primary(),
		}
	case common.PrimaryPreferredMode:
		opt = &options.CollectionOptions{
			ReadPreference: readpref.PrimaryPreferred(readpref.WithMaxStaleness(maxStalenessSeconds)),
		}
	case common.SecondaryMode:
		opt = &options.CollectionOptions{
			ReadPreference: readpref.Secondary(readpref.WithMaxStaleness(maxStalenessSeconds)),
		}
	case common.SecondaryPreferredMode:
		opt = &options.CollectionOptions{
			ReadPreference: readpref.SecondaryPreferred(readpref.WithMaxStaleness(maxStalenessSeconds)),
		}
	case common.NearestMode:
		opt = &options.CollectionOptions{
			ReadPreference: readpref.Nearest(readpref.WithMaxStaleness(maxStalenessSeconds)),
		}
	}

	return opt
}
