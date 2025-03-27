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

// Package local is the local mongodb dao implementation
package local

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	idgen "configcenter/pkg/id-gen"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/dal/types"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// Mongo is the mongo client for one tenant
type Mongo struct {
	// tenant is the tenant id
	tenant string
	// ignoreTenant is used for platform operations that do not need to specify tenant
	ignoreTenant bool
	cli          *MongoClient
	tm           *TxnManager
	conf         *MongoCliConf

	// TODO: remove this when all db clients use ShardingDB
	enableSharding bool
}

var _ DB = new(Mongo)

// NewMgo returns new RDB
func NewMgo(config MongoConf, timeout time.Duration) (*Mongo, error) {
	client, err := NewMongoClient(true, "", &config, timeout)
	if err != nil {
		return nil, err
	}

	mgo := &Mongo{
		cli: client,
		tm:  new(TxnManager),
		conf: &MongoCliConf{
			DisableInsert: config.DisableInsert,
		},
	}

	ctx := context.Background()
	mgo.conf.IDGenStep, err = mgo.InitIDGenerator(ctx)
	if err != nil {
		return nil, err
	}

	return mgo, nil
}

// MongoOptions is the mongo options
type MongoOptions struct {
	Tenant       string
	IgnoreTenant bool
}

// NewMongo new DB
func NewMongo(cli *MongoClient, tm *TxnManager, conf *MongoCliConf, opt ...*MongoOptions) (*Mongo, error) {
	if cli == nil || tm == nil || conf == nil {
		return nil, errors.New("not all mongo client info is set")
	}

	m := &Mongo{
		cli:            cli,
		tm:             tm,
		conf:           conf,
		enableSharding: true,
	}

	if len(opt) > 0 {
		m.tenant = opt[0].Tenant
		m.ignoreTenant = opt[0].IgnoreTenant
	}

	return m, nil
}

// MongoClient is the mongodb client
type MongoClient struct {
	dbc      *mongo.Client
	dbname   string
	uuid     string
	disabled bool
}

// NewMongoClient new mongodb client
func NewMongoClient(isMaster bool, uuid string, config *MongoConf, timeout time.Duration) (*MongoClient, error) {
	connStr, err := connstring.Parse(config.URI)
	if err != nil {
		return nil, err
	}
	if config.RsName == "" {
		return nil, fmt.Errorf("mongodb rsName not set")
	}
	if !isMaster && uuid == "" {
		return nil, fmt.Errorf("mongodb client uuid not set")
	}

	socketTimeout := time.Second * time.Duration(config.SocketTimeout)
	maxConnIdleTime := 25 * time.Minute
	appName := common.GetIdentification()
	// do not change this, our transaction plan need it to false.
	// it's related with the transaction number(eg txnNumber) in a transaction session.
	disableWriteRetry := false
	conOpt := options.ClientOptions{
		MaxPoolSize:     &config.MaxOpenConns,
		MinPoolSize:     &config.MaxIdleConns,
		ConnectTimeout:  &timeout,
		SocketTimeout:   &socketTimeout,
		ReplicaSet:      &config.RsName,
		RetryWrites:     &disableWriteRetry,
		MaxConnIdleTime: &maxConnIdleTime,
		AppName:         &appName,
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.URI), &conOpt)
	if nil != err {
		return nil, err
	}

	// TODO: add this check later, this command needs authorize to get version.
	// if err := checkMongodbVersion(connStr.Database, client); err != nil {
	// 	return nil, err
	// }

	// initialize mongodb related metrics
	initMongoMetric()

	return &MongoClient{
		dbc:      client,
		dbname:   connStr.Database,
		uuid:     uuid,
		disabled: config.Disabled,
	}, nil
}

// Client returns mongodb client
func (m *MongoClient) Client() *mongo.Client {
	return m.dbc
}

// DBName returns mongodb database name
func (m *MongoClient) DBName() string {
	return m.dbname
}

// Database returns mongodb database client
func (m *MongoClient) Database() *mongo.Database {
	return m.dbc.Database(m.dbname)
}

// UUID returns uuid of the mongodb client
func (m *MongoClient) UUID() string {
	return m.uuid
}

// SetUUID set mongodb client uuid
func (m *MongoClient) SetUUID(uuid string) {
	m.uuid = uuid
}

// Disabled returns whether the mongodb client is disabled
func (m *MongoClient) Disabled() bool {
	return m.disabled
}

// SetDisabled set mongodb client disabled status
func (m *MongoClient) SetDisabled(disabled bool) {
	m.disabled = disabled
}

// convColl convert collection name
func (c *Mongo) convColl(collection string) (string, error) {
	if common.IsPlatformTableWithTenant(collection) {
		return collection, nil
	}

	if c.ignoreTenant {
		if !common.IsPlatformTable(collection) {
			return "", fmt.Errorf("tenant table %s has no tenant", collection)
		}
		return collection, nil
	}

	if common.IsPlatformTable(collection) {
		return "", fmt.Errorf("platform table %s do not need tenant", collection)
	}
	return common.GenTenantTableName(c.tenant, collection), nil
}

// InitIDGenerator init id generator by config admin, returns id generator step
func (c *Mongo) InitIDGenerator(ctx context.Context) (int, error) {
	cond := map[string]interface{}{"_id": common.ConfigAdminID}

	confData := make(map[string]string)
	err := c.Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(ctx, &confData)
	if err != nil {
		// watch database & low version has no config admin field, so we use default step 1
		if c.IsNotFoundError(err) {
			return 1, nil
		}
		return 0, fmt.Errorf("init id generator but get config admin failed, err: %v", err)
	}

	// id generator config is not set, use default step 1
	confStr := confData[common.ConfigAdminValueField]
	if !gjson.Get(confStr, "id_generator").Exists() {
		return 1, nil
	}

	conf := new(metadata.PlatformSettingConfig)
	if err = json.Unmarshal([]byte(confStr), conf); err != nil {
		return 0, fmt.Errorf("unmarshal config admin failed, err: %v, config: %s", err, confStr)
	}

	idGenConf := conf.IDGenerator
	if err = idGenConf.Validate(); err != nil {
		return 0, fmt.Errorf("config admin id gen config is invalid, err: %v, config: %+v", err, idGenConf)
	}

	if len(idGenConf.InitID) == 0 {
		return idGenConf.Step, nil
	}

	// update id generator sequence id by config admin
	for typ, id := range idGenConf.InitID {
		if err = c.updateIDGenSeqID(ctx, typ, id); err != nil {
			return 0, err
		}
	}

	// delete config admin id generator init id config to avoid updating again
	conf.IDGenerator.InitID = nil
	updateVal, err := json.Marshal(conf)
	if err != nil {
		return 0, fmt.Errorf("marshal config admin failed, err: %v, config: %+v", err, conf)
	}

	data := map[string]interface{}{
		common.ConfigAdminValueField: string(updateVal),
		common.LastTimeField:         time.Now(),
	}

	if err = c.Table(common.BKTableNameSystem).Update(ctx, cond, data); err != nil {
		return 0, fmt.Errorf("update config admin failed, err: %v, data: %+v", err, data)
	}

	return idGenConf.Step, nil
}

func (c *Mongo) updateIDGenSeqID(ctx context.Context, typ idgen.IDGenType, id uint64) error {
	sequenceName, exists := idgen.GetIDGenSequenceName(typ)
	if !exists {
		return fmt.Errorf("id generator type %s is invalid", typ)
	}

	// add id generator if not exists
	cnt, err := c.Table(common.BKTableNameIDgenerator).Find(map[string]interface{}{"_id": sequenceName}).Count(ctx)
	if err != nil {
		return fmt.Errorf("check if %s id generator exists failed, err: %v, ", sequenceName, err)
	}

	if cnt == 0 {
		insertData := map[string]interface{}{
			"_id":        sequenceName,
			"SequenceID": id,
		}
		err = c.Table(common.BKTableNameIDgenerator).Insert(ctx, insertData)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("insert id generator failed, err: %v, data: %+v", err, insertData)
		}
	}

	// update id generator sequence id if it is greater than current sequence id
	updateCond := map[string]interface{}{
		"_id":        sequenceName,
		"SequenceID": map[string]interface{}{common.BKDBLT: id},
	}

	data := map[string]interface{}{"SequenceID": id}
	if err = c.Table(common.BKTableNameIDgenerator).Update(ctx, updateCond, data); err != nil {
		return fmt.Errorf("update id generator failed, err: %v, cond: %+v, data: %+v", err, updateCond, data)
	}
	return nil
}

// checkMongodbVersion TODO
// from now on, mongodb version must >= 7.0
func checkMongodbVersion(db string, client *mongo.Client) error {
	serverStatus, err := client.Database(db).RunCommand(
		context.Background(),
		bson.D{{"serverStatus", 1}},
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
	// version must be >= v7.0
	major, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("parse mongodb version %s major failed, err: %v", version.StringValue(), err)
	}
	if major < 7 {
		return errors.New("mongodb version must be >= v7.0")
	}

	return nil
}

// InitTxnManager TxnID management of initial transaction
// TODO remove this
func (c *Mongo) InitTxnManager(r redis.Client) error {
	return c.tm.InitTxnManager(r)
}

// Close replica client
func (c *Mongo) Close() error {
	c.cli.Client().Disconnect(context.TODO())
	return nil
}

// Ping replica client
func (c *Mongo) Ping() error {
	return c.cli.Client().Ping(context.TODO(), nil)
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
		if strings.Contains(err.Error(), "all indexes already exist") {
			return true
		}
		if strings.Contains(err.Error(), "already exists with a different name") {
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
	if c.enableSharding {
		var err error
		collName, err = c.convColl(collName)
		if err != nil {
			return &errColl{err: err}
		}
	}

	col := Collection{}
	col.collName = collName
	col.Mongo = c
	return &col
}

// GetDBClient TODO
// get db client
func (c *Mongo) GetDBClient() *mongo.Client {
	return c.cli.Client()
}

// GetDBName TODO
// get db name
func (c *Mongo) GetDBName() string {
	return c.cli.DBName()
}
func (c *Mongo) redirectTable(tableName string) string {
	if common.IsObjectInstShardingTable(tableName) {
		tableName = common.BKTableNameBaseInst
	} else if common.IsObjectInstAsstShardingTable(tableName) {
		tableName = common.BKTableNameInstAsst
	}
	return tableName
}

// NextSequence 获取新序列号(非事务)
func (c *Mongo) NextSequence(ctx context.Context, sequenceName string) (uint64, error) {
	if c.enableSharding && !c.ignoreTenant {
		return 0, errors.New("next sequence do not need tenant")
	}

	sequenceName = c.redirectTable(sequenceName)

	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	defer func() {
		blog.V(4).InfoDepthf(2, "mongo next-sequence cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	}()

	// 直接使用新的context，确保不会用到事务,不会因为context含有session而使用分布式事务，防止产生相同的序列号
	ctx = context.Background()

	coll := c.cli.Database().Collection(common.BKTableNameIDgenerator)

	Update := bson.M{
		"$inc":         bson.M{"SequenceID": c.conf.IDGenStep},
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
	if c.enableSharding && !c.ignoreTenant {
		return nil, errors.New("next sequence do not need tenant")
	}

	if num == 0 {
		return make([]uint64, 0), nil
	}
	sequenceName = c.redirectTable(sequenceName)

	if c.conf.DisableInsert && idgen.IsIDGenSeqName(sequenceName) {
		return nil, errors.New("insertion is disabled")
	}

	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	defer func() {
		blog.V(4).InfoDepthf(2, "mongo next-sequences cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	}()

	// 直接使用新的context，确保不会用到事务,不会因为context含有session而使用分布式事务，防止产生相同的序列号
	ctx = context.Background()

	coll := c.cli.Database().Collection(common.BKTableNameIDgenerator)

	Update := bson.M{
		"$inc":         bson.M{"SequenceID": num * c.conf.IDGenStep},
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
		sequences[i] = uint64((i-num+1)*c.conf.IDGenStep) + doc.SequenceID
	}

	return sequences, err
}

// Idgen TODO
type Idgen struct {
	ID         string `bson:"_id"`
	SequenceID uint64 `bson:"SequenceID"`
}

// HasTable 判断是否存在集合
func (c *Mongo) HasTable(ctx context.Context, collName string) (bool, error) {
	if c.enableSharding {
		var err error
		collName, err = c.convColl(collName)
		if err != nil {
			return false, err
		}
	}

	cursor, err := c.cli.Database().ListCollections(ctx, bson.M{"name": collName, "type": "collection"})
	if err != nil {
		return false, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		return true, nil
	}

	return false, nil
}

// ListTables 获取所有的表名
func (c *Mongo) ListTables(ctx context.Context) ([]string, error) {
	filter := bson.M{"type": "collection"}
	if c.enableSharding {
		if c.ignoreTenant {
			filter["name"] = bson.M{common.BKDBIN: common.PlatformTables()}
		} else {
			filter["name"] = bson.M{common.BKDBLIKE: fmt.Sprintf("^%s_", c.tenant)}
		}
	}
	return c.cli.Database().ListCollectionNames(ctx, filter)
}

// DropTable 移除集合
func (c *Mongo) DropTable(ctx context.Context, collName string) error {
	if c.enableSharding {
		var err error
		collName, err = c.convColl(collName)
		if err != nil {
			return err
		}
	}
	return c.cli.Database().Collection(collName).Drop(ctx)
}

// CreateTable 创建集合 TODO test
func (c *Mongo) CreateTable(ctx context.Context, collName string) error {
	if c.enableSharding {
		var err error
		collName, err = c.convColl(collName)
		if err != nil {
			return err
		}
	}
	return c.cli.Database().RunCommand(ctx, map[string]interface{}{"create": collName}).Err()
}

// RenameTable 更新集合名称
func (c *Mongo) RenameTable(ctx context.Context, prevName, currName string) error {
	if c.enableSharding {
		var err error
		prevName, err = c.convColl(prevName)
		if err != nil {
			return err
		}
		currName, err = c.convColl(currName)
		if err != nil {
			return err
		}
	}

	cmd := bson.D{
		{"renameCollection", c.cli.DBName() + "." + prevName},
		{"to", c.cli.DBName() + "." + currName},
	}
	return c.cli.Client().Database("admin").RunCommand(ctx, cmd).Err()
}
