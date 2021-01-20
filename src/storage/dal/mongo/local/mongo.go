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

package local

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/dal/types"
	dtype "configcenter/src/storage/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Mongo struct {
	dbc    *mongo.Client
	dbname string
	sess   mongo.Session
	tm     *TxnManager
}

var _ dal.DB = new(Mongo)

type MongoConf struct {
	TimeoutSeconds int
	MaxOpenConns   uint64
	MaxIdleConns   uint64
	URI            string
	RsName         string
	SocketTimeout  int
}

// NewMgo returns new RDB
func NewMgo(config MongoConf, timeout time.Duration) (*Mongo, error) {
	connStr, err := connstring.Parse(config.URI)
	if nil != err {
		return nil, err
	}
	if config.RsName == "" {
		return nil, fmt.Errorf("mongodb rsName not set")
	}
	socketTimeout := time.Second * time.Duration(config.SocketTimeout)
	// do not change this, our transaction plan need it to false.
	// it's related with the transaction number(eg txnNumber) in a transaction session.
	disableWriteRetry := false
	conOpt := options.ClientOptions{
		MaxPoolSize:    &config.MaxOpenConns,
		MinPoolSize:    &config.MaxIdleConns,
		ConnectTimeout: &timeout,
		SocketTimeout:  &socketTimeout,
		ReplicaSet:     &config.RsName,
		RetryWrites:    &disableWriteRetry,
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(config.URI), &conOpt)
	if nil != err {
		return nil, err
	}

	if err := client.Connect(context.TODO()); nil != err {
		return nil, err
	}

	// TODO: add this check later, this command needs authorize to get version.
	// if err := checkMongodbVersion(connStr.Database, client); err != nil {
	// 	return nil, err
	// }

	// initialize mongodb related metrics
	initMongoMetric()

	return &Mongo{
		dbc:    client,
		dbname: connStr.Database,
		tm:     &TxnManager{},
	}, nil
}

// from now on, mongodb version must >= 4.2.0
func checkMongodbVersion(db string, client *mongo.Client) error {
	serverStatus, err := client.Database(db).RunCommand(
		context.Background(),
		bsonx.Doc{{"serverStatus", bsonx.Int32(1)}},
	).DecodeBytes()
	if err != nil {
		return err
	}

	version, err := serverStatus.LookupErr("version")
	if err != nil {
		return err
	}

	fields := strings.Split(version.StringValue(), ".")
	if len(fields) != 3 {
		return fmt.Errorf("got invalid mongodb version: %v", version.StringValue())
	}
	// version must be >= v4.2.0
	major, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("parse mongodb version %s major failed, err: %v", version.StringValue(), err)
	}
	if major < 4 {
		return errors.New("mongodb version must be >= v4.2.0")
	}

	minor, err := strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("parse mongodb version %s minor failed, err: %v", version.StringValue(), err)
	}
	if minor < 2 {
		return errors.New("mongodb version must be >= v4.2.0")
	}
	return nil
}

// InitTxnManager TxnID management of initial transaction
func (c *Mongo) InitTxnManager(r redis.Client) error {
	return c.tm.InitTxnManager(r)
}

// Close replica client
func (c *Mongo) Close() error {
	c.dbc.Disconnect(context.TODO())
	return nil
}

// Ping replica client
func (c *Mongo) Ping() error {
	return c.dbc.Ping(context.TODO(), nil)
}

// IsDuplicatedError check duplicated error
func (c *Mongo) IsDuplicatedError(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "The existing index") {
			return true
		}
		if strings.Contains(err.Error(), "There's already an index with name") {
			return true
		}
		if strings.Contains(err.Error(), "E11000 duplicate") {
			return true
		}
		if strings.Contains(err.Error(), "IndexOptionsConflict") {
			return true
		}
		if strings.Contains(err.Error(), "already exists with a different name") {
			return true
		}
		if strings.Contains(err.Error(), "already exists with different options") {
			return true
		}
	}
	return err == types.ErrDuplicated
}

// IsNotFoundError check the not found error
func (c *Mongo) IsNotFoundError(err error) bool {
	return err == types.ErrDocumentNotFound
}

// Table collection operation
func (c *Mongo) Table(collName string) types.Table {
	col := Collection{}
	col.collName = collName
	col.Mongo = c
	return &col
}

// get db client
func (c *Mongo) GetDBClient() *mongo.Client {
	return c.dbc
}

// get db name
func (c *Mongo) GetDBName() string {
	return c.dbname
}

// Collection implement client.Collection interface
type Collection struct {
	collName string // 集合名
	*Mongo
}

// Find 查询多个并反序列化到 Result
func (c *Collection) Find(filter types.Filter, opts ...types.FindOpts) types.Find {
	find := &Find{
		Collection: c,
		filter:     filter,
		projection: make(map[string]int),
	}

	if len(opts) == 0 {
		find.projection["_id"] = 0
		return find
	}

	if !opts[0].WithObjectID {
		find.projection["_id"] = 0
		return find
	}
	return find
}

// Find define a find operation
type Find struct {
	*Collection

	projection map[string]int
	filter     types.Filter
	start      int64
	limit      int64
	sort       map[string]interface{}
}

// Fields 查询字段
func (f *Find) Fields(fields ...string) types.Find {
	for _, field := range fields {
		if len(field) <= 0 {
			continue
		}
		f.projection[field] = 1
	}
	return f
}

// Sort 查询排序
func (f *Find) Sort(sort string) types.Find {
	if sort != "" {
		sortArr := strings.Split(sort, ",")
		f.sort = make(map[string]interface{}, 0)
		for _, sortItem := range sortArr {
			sortItemArr := strings.Split(sortItem, ":")
			sortKey := strings.TrimLeft(sortItemArr[0], "+-")
			if len(sortItemArr) == 2 {
				sortDescFlag := strings.TrimSpace(sortItemArr[1])
				if sortDescFlag == "-1" {
					f.sort[sortKey] = -1
				} else {
					f.sort[sortKey] = 1
				}
			} else {
				if strings.HasPrefix(sortItemArr[0], "-") {
					f.sort[sortKey] = -1
				} else {
					f.sort[sortKey] = 1
				}
			}
		}
	}

	return f
}

// Start 查询上标
func (f *Find) Start(start uint64) types.Find {
	// change to int64,后续改成int64
	dbStart := int64(start)
	f.start = dbStart
	return f
}

// Limit 查询限制
func (f *Find) Limit(limit uint64) types.Find {
	// change to int64,后续改成int64
	dbLimit := int64(limit)
	f.limit = dbLimit
	return f
}

var hostSpecialFieldMap = map[string]bool{
	common.BKHostInnerIPField: true,
	common.BKHostOuterIPField: true,
	common.BKOperatorField:    true,
	common.BKBakOperatorField: true,
}

// All 查询多个
func (f *Find) All(ctx context.Context, result interface{}) error {
	mtc.collectOperCount(f.collName, findOper)

	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	defer func() {
		mtc.collectOperDuration(f.collName, findOper, time.Since(start))
	}()

	err := validHostType(f.collName, f.projection, result, rid)
	if err != nil {
		return err
	}

	findOpts := &options.FindOptions{}
	if len(f.projection) != 0 {
		findOpts.Projection = f.projection
	}
	if f.start != 0 {
		findOpts.SetSkip(f.start)
	}
	if f.limit != 0 {
		findOpts.SetLimit(f.limit)
	}
	if len(f.sort) != 0 {
		findOpts.SetSort(f.sort)
	}
	// 查询条件为空时候，mongodb 不返回数据
	if f.filter == nil {
		f.filter = bson.M{}
	}

	opt := getCollectionOption(ctx)

	return f.tm.AutoRunWithTxn(ctx, f.dbc, func(ctx context.Context) error {
		cursor, err := f.dbc.Database(f.dbname).Collection(f.collName, opt).Find(ctx, f.filter, findOpts)
		if err != nil {
			mtc.collectErrorCount(f.collName, findOper)
			return err
		}
		return cursor.All(ctx, result)
	})
}

// One 查询一个
func (f *Find) One(ctx context.Context, result interface{}) error {
	mtc.collectOperCount(f.collName, findOper)

	start := time.Now()
	rid := ctx.Value(common.ContextRequestIDField)
	defer func() {
		mtc.collectOperDuration(f.collName, findOper, time.Since(start))
	}()

	err := validHostType(f.collName, f.projection, result, rid)
	if err != nil {
		return err
	}

	findOpts := &options.FindOptions{}
	if len(f.projection) != 0 {
		findOpts.Projection = f.projection
	}
	if f.start != 0 {
		findOpts.SetSkip(f.start)
	}
	if f.limit != 0 {
		findOpts.SetLimit(1)
	}
	if len(f.sort) != 0 {
		findOpts.SetSort(f.sort)
	}
	// 查询条件为空时候，mongodb panic
	if f.filter == nil {
		f.filter = bson.M{}
	}

	opt := getCollectionOption(ctx)
	return f.tm.AutoRunWithTxn(ctx, f.dbc, func(ctx context.Context) error {
		cursor, err := f.dbc.Database(f.dbname).Collection(f.collName, opt).Find(ctx, f.filter, findOpts)
		if err != nil {
			mtc.collectErrorCount(f.collName, findOper)
			return err
		}

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			return cursor.Decode(result)
		}
		return types.ErrDocumentNotFound
	})

}

// Count 统计数量(非事务)
func (f *Find) Count(ctx context.Context) (uint64, error) {
	mtc.collectOperCount(f.collName, countOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(f.collName, countOper, time.Since(start))
	}()

	if f.filter == nil {
		f.filter = bson.M{}
	}

	opt := getCollectionOption(ctx)

	sessCtx, _, useTxn, err := f.tm.GetTxnContext(ctx, f.dbc)
	if err != nil {
		return 0, err
	}
	if !useTxn {
		// not use transaction.
		cnt, err := f.dbc.Database(f.dbname).Collection(f.collName, opt).CountDocuments(ctx, f.filter)
		if err != nil {
			mtc.collectErrorCount(f.collName, countOper)
			return 0, err
		}

		return uint64(cnt), err
	} else {
		// use transaction
		cnt, err := f.dbc.Database(f.dbname).Collection(f.collName, opt).CountDocuments(sessCtx, f.filter)
		// do not release th session, otherwise, the session will be returned to the
		// session pool and will be reused. then mongodb driver will increase the transaction number
		// automatically and do read/write retry if policy is set.
		// mongo.CmdbReleaseSession(ctx, session)
		if err != nil {
			mtc.collectErrorCount(f.collName, countOper)
			return 0, err
		}
		return uint64(cnt), nil
	}
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *Collection) Insert(ctx context.Context, docs interface{}) error {
	mtc.collectOperCount(c.collName, insertOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, insertOper, time.Since(start))
	}()

	rows := util.ConverToInterfaceSlice(docs)

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
	mtc.collectOperCount(c.collName, deleteOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, deleteOper, time.Since(start))
	}()

	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		if err := c.tryArchiveDeletedDoc(ctx, filter); err != nil {
			mtc.collectErrorCount(c.collName, deleteOper)
			return err
		}
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).DeleteMany(ctx, filter)
		if err != nil {
			mtc.collectErrorCount(c.collName, deleteOper)
			return err
		}

		return nil
	})

}

func (c *Collection) tryArchiveDeletedDoc(ctx context.Context, filter types.Filter) error {
	switch c.collName {
	case common.BKTableNameModuleHostConfig:
	case common.BKTableNameBaseHost:
	case common.BKTableNameBaseApp:
	case common.BKTableNameBaseSet:
	case common.BKTableNameBaseModule:
	case common.BKTableNameSetTemplate:
	case common.BKTableNameBaseInst:
	case common.BKTableNameBaseProcess:
	case common.BKTableNameProcessInstanceRelation:
	default:
		// do not archive the delete docs
		return nil
	}

	docs := make([]bsonx.Doc, 0)
	cursor, err := c.dbc.Database(c.dbname).Collection(c.collName).Find(ctx, filter, nil)
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
		archives[idx] = metadata.DeleteArchive{
			Oid:    doc.Lookup("_id").ObjectID().Hex(),
			Detail: doc.Delete("_id"),
			Coll:   c.collName,
		}
	}

	_, err = c.dbc.Database(c.dbname).Collection(common.BKTableNameDelArchive).InsertMany(ctx, archives)
	return err
}

// NextSequence 获取新序列号(非事务)
func (c *Mongo) NextSequence(ctx context.Context, sequenceName string) (uint64, error) {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	defer func() {
		blog.V(4).InfoDepthf(2, "mongo next-sequence cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	}()

	// 直接使用新的context，确保不会用到事务,不会因为context含有session而使用分布式事务，防止产生相同的序列号
	ctx = context.Background()

	coll := c.dbc.Database(c.dbname).Collection("cc_idgenerator")

	Update := bson.M{
		"$inc":         bson.M{"SequenceID": int64(1)},
		"$setOnInsert": bson.M{"create_time": time.Now()},
		"$set":         bson.M{"last_time": time.Now()},
	}
	filter := bson.M{"_id": sequenceName}
	upsert := true
	returnChange := options.After
	opt := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnChange,
	}

	doc := Idgen{}
	err := coll.FindOneAndUpdate(ctx, filter, Update, opt).Decode(&doc)
	if err != nil {
		return 0, err
	}
	return doc.SequenceID, err
}

// NextSequences 批量获取新序列号(非事务)
func (c *Mongo) NextSequences(ctx context.Context, sequenceName string, num int) ([]uint64, error) {
	if num == 0 {
		return make([]uint64, 0), nil
	}

	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	defer func() {
		blog.V(4).InfoDepthf(2, "mongo next-sequences cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	}()

	// 直接使用新的context，确保不会用到事务,不会因为context含有session而使用分布式事务，防止产生相同的序列号
	ctx = context.Background()

	coll := c.dbc.Database(c.dbname).Collection("cc_idgenerator")

	Update := bson.M{
		"$inc":         bson.M{"SequenceID": num},
		"$setOnInsert": bson.M{"create_time": time.Now()},
		"$set":         bson.M{"last_time": time.Now()},
	}
	filter := bson.M{"_id": sequenceName}
	upsert := true
	returnChange := options.After
	opt := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnChange,
	}

	doc := Idgen{}
	err := coll.FindOneAndUpdate(ctx, filter, Update, opt).Decode(&doc)
	if err != nil {
		return nil, err
	}

	sequences := make([]uint64, num)
	for i := 0; i < num; i++ {
		sequences[i] = uint64(i-num) + doc.SequenceID + 1
	}

	return sequences, err
}

type Idgen struct {
	ID         string `bson:"_id"`
	SequenceID uint64 `bson:"SequenceID"`
}

// HasTable 判断是否存在集合
func (c *Mongo) HasTable(ctx context.Context, collName string) (bool, error) {
	cursor, err := c.dbc.Database(c.dbname).ListCollections(ctx, bson.M{"name": collName, "type": "collection"})
	if err != nil {
		return false, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		return true, nil
	}

	return false, nil
}

// DropTable 移除集合
func (c *Mongo) DropTable(ctx context.Context, collName string) error {
	return c.dbc.Database(c.dbname).Collection(collName).Drop(ctx)
}

// CreateTable 创建集合 TODO test
func (c *Mongo) CreateTable(ctx context.Context, collName string) error {
	return c.dbc.Database(c.dbname).RunCommand(ctx, map[string]interface{}{"create": collName}).Err()
}

// CreateIndex 创建索引
func (c *Collection) CreateIndex(ctx context.Context, index types.Index) error {
	createIndexOpt := &options.IndexOptions{
		Background: &index.Background,
		Unique:     &index.Unique,
	}
	if index.Name != "" {
		createIndexOpt.Name = &index.Name
	}

	if index.ExpireAfterSeconds != 0 {
		createIndexOpt.SetExpireAfterSeconds(index.ExpireAfterSeconds)
	}

	createIndexInfo := mongo.IndexModel{
		Keys:    index.Keys,
		Options: createIndexOpt,
	}

	indexView := c.dbc.Database(c.dbname).Collection(c.collName).Indexes()
	_, err := indexView.CreateOne(ctx, createIndexInfo)
	if err != nil {
		// ignore the duplicated index error
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
	}

	return err
}

// DropIndex remove index by name
func (c *Collection) DropIndex(ctx context.Context, indexName string) error {
	indexView := c.dbc.Database(c.dbname).Collection(c.collName).Indexes()
	_, err := indexView.DropOne(ctx, indexName)
	return err
}

// Indexes get all indexes for the collection
func (c *Collection) Indexes(ctx context.Context) ([]types.Index, error) {
	indexView := c.dbc.Database(c.dbname).Collection(c.collName).Indexes()
	cursor, err := indexView.List(ctx)
	if nil != err {
		return nil, err
	}
	defer cursor.Close(ctx)
	var indexs []types.Index
	for cursor.Next(ctx) {
		idxResult := types.Index{}
		cursor.Decode(&idxResult)
		indexs = append(indexs, idxResult)
	}

	return indexs, nil
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
func (c *Collection) RenameColumn(ctx context.Context, oldName, newColumn string) error {
	mtc.collectOperCount(c.collName, columnOper)

	start := time.Now()
	defer func() {
		mtc.collectOperDuration(c.collName, columnOper, time.Since(start))
	}()

	datac := dtype.Document{"$rename": dtype.Document{oldName: newColumn}}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, dtype.Document{}, datac)
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
func (c *Collection) AggregateAll(ctx context.Context, pipeline interface{}, result interface{}) error {
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
		return decodeCusorIntoSlice(ctx, cursor, result)
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

// Distinct Finds the distinct values for a specified field across a single collection or view and returns the results in an
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

func decodeCusorIntoSlice(ctx context.Context, cursor *mongo.Cursor, result interface{}) error {
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

// validHostType valid if host query uses specified type that transforms ip & operator array to string
func validHostType(collection string, projection map[string]int, result interface{}, rid interface{}) error {
	if result == nil {
		blog.Errorf("host query result is nil, rid: %s", rid)
		return fmt.Errorf("host query result type invalid")
	}

	if collection != common.BKTableNameBaseHost {
		return nil
	}

	// check if specified fields include special fields
	if len(projection) != 0 {
		needCheck := false
		for field := range projection {
			if hostSpecialFieldMap[field] {
				needCheck = true
				break
			}
		}
		if !needCheck {
			return nil
		}
	}

	resType := reflect.TypeOf(result)
	if resType.Kind() != reflect.Ptr {
		blog.Errorf("host query result type(%v) not pointer type, rid: %v", resType, rid)
		return fmt.Errorf("host query result type invalid")
	}
	// if result is *map[string]interface{} type, it must be *metadata.HostMapStr type
	if resType.ConvertibleTo(reflect.TypeOf(&map[string]interface{}{})) {
		if resType != reflect.TypeOf(&metadata.HostMapStr{}) {
			blog.Errorf("host query result type(%v) not match *metadata.HostMapStr type, rid: %v", resType, rid)
			return fmt.Errorf("host query result type invalid")
		}
		return nil
	}

	resElem := resType.Elem()
	switch resElem.Kind() {
	case reflect.Struct:
		// if result is *struct type, the special field in it must be metadata.StringArrayToString type
		numField := resElem.NumField()
		validType := reflect.TypeOf(metadata.StringArrayToString(""))
		for i := 0; i < numField; i++ {
			field := resElem.Field(i)
			bsonTag := field.Tag.Get("bson")
			if bsonTag == "" {
				blog.Errorf("host query result field(%s) has empty bson tag, rid: %v", field.Name, rid)
				return fmt.Errorf("host query result type invalid")
			}
			if hostSpecialFieldMap[bsonTag] && field.Type != validType {
				blog.Errorf("host query result field type(%v) not match *metadata.StringArrayToString type", field.Type)
				return fmt.Errorf("host query result type invalid")
			}
		}
	case reflect.Slice:
		// check if slice item is valid type, map or struct validation is similar as before
		elem := resElem.Elem()
		if elem.ConvertibleTo(reflect.TypeOf(map[string]interface{}{})) {
			if elem != reflect.TypeOf(metadata.HostMapStr{}) {
				blog.Errorf("host query result type(%v) not match *[]metadata.HostMapStr type", resType)
				return fmt.Errorf("host query result type invalid")
			}
			return nil
		}

		if elem.Kind() != reflect.Struct {
			blog.Errorf("host query result type(%v) not struct pointer type or map type", resType)
			return fmt.Errorf("host query result type invalid")
		}
		numField := elem.NumField()
		validType := reflect.TypeOf(metadata.StringArrayToString(""))
		for i := 0; i < numField; i++ {
			field := elem.Field(i)
			bsonTag := field.Tag.Get("bson")
			if bsonTag == "" {
				blog.Errorf("host query result field(%s) has empty bson tag, rid: %v", field.Name, rid)
				return fmt.Errorf("host query result type invalid")
			}
			if hostSpecialFieldMap[bsonTag] && field.Type != validType {
				blog.Errorf("host query result field type(%v) not match *metadata.StringArrayToString type", field.Type)
				return fmt.Errorf("host query result type invalid")
			}
		}
	default:
		blog.Errorf("host query result type(%v) not pointer of map, struct or slice, rid: %v", resType, rid)
		return fmt.Errorf("host query result type invalid")
	}
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
