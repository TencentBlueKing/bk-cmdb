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

package sharding

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"configcenter/pkg/tenant"
	"configcenter/pkg/tenant/logics"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/cryptor"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"

	"github.com/google/uuid"
)

// ShardingMongoManager is the sharding db manager for mongo
type ShardingMongoManager struct {
	*shardingMongoClient
	// conf is the mongo client config
	conf *local.MongoCliConf
}

// NewShardingMongo returns new sharding db manager for mongo
func NewShardingMongo(config local.MongoConf, timeout time.Duration, crypto cryptor.Cryptor) (ShardingDB, error) {
	clientInfo, masterMongo, err := newShardingMongoClient(config, timeout, crypto)
	if err != nil {
		return nil, err
	}

	sharding := &ShardingMongoManager{
		shardingMongoClient: clientInfo,
		conf:                &local.MongoCliConf{DisableInsert: config.DisableInsert},
	}

	sharding.conf.IDGenStep, err = masterMongo.InitIDGenerator(context.Background())
	if err != nil {
		return nil, err
	}

	err = logics.Init(&logics.Options{DB: sharding.IgnoreTenant()})
	if err != nil {
		return nil, err
	}

	if err = sharding.refreshTenantDBMap(); err != nil {
		return nil, err
	}

	// loop refresh tenant to db relation
	go func() {
		for {
			time.Sleep(time.Minute)
			if err = sharding.refreshTenantDBMap(); err != nil {
				blog.Errorf("refresh tenant to db relation failed, err: %v", err)
				continue
			}
		}
	}()

	tenantChan := tenant.NewTenantEventChan(fmt.Sprintf("sharding_db_%s", clientInfo.masterCli.UUID()))
	go func() {
		for e := range tenantChan {
			switch e.EventType {
			case tenant.Create:
				client, exists := sharding.dbClientMap[e.Tenant.Database]
				if !exists {
					blog.Errorf("tenant %s related db %s config not found", e.Tenant.TenantID, e.Tenant.Database)
					continue
				}
				sharding.tenantCli.set(e.Tenant.TenantID, client)
			case tenant.Delete:
				sharding.tenantCli.delete(e.Tenant.TenantID)
			}
		}
	}()

	return sharding, nil
}

// shardingMongoClient is the common structure that stores all sharding db mongo client info
type shardingMongoClient struct {
	// masterCli is the client for master mongodb, master mongodb stores the platform data and some tenant data
	masterCli *local.MongoClient
	// newDataCli is the client for mongodb that new data without specified db will be stored into
	newDataCli *local.MongoClient
	// tenantCli is the tenant id to mongodb client map
	tenantCli *tenantMongoCliMap
	// dbClientMap is the db uuid to mongodb client map
	dbClientMap map[string]*local.MongoClient
	// tm is the transaction manager
	tm *local.ShardingTxnManager
}

type tenantMongoCliMap struct {
	tenantCli map[string]*local.MongoClient
	sync.RWMutex
}

func (m *tenantMongoCliMap) get(tenantID string) (*local.MongoClient, bool) {
	m.RLock()
	cli, exists := m.tenantCli[tenantID]
	m.RUnlock()
	return cli, exists
}

func (m *tenantMongoCliMap) set(tenantID string, cli *local.MongoClient) {
	m.Lock()
	m.tenantCli[tenantID] = cli
	m.Unlock()
}

func (m *tenantMongoCliMap) delete(tenantID string) {
	m.Lock()
	delete(m.tenantCli, tenantID)
	m.Unlock()
}

func newShardingMongoClient(config local.MongoConf, timeout time.Duration, crypto cryptor.Cryptor) (
	*shardingMongoClient, *local.Mongo, error) {

	// connect master mongodb
	masterCli, err := local.NewMongoClient(true, "", &config, timeout)
	if err != nil {
		return nil, nil, fmt.Errorf("new master mongo client failed, err: %v", err)
	}

	clientInfo := &shardingMongoClient{
		masterCli:  masterCli,
		newDataCli: nil,
		tenantCli: &tenantMongoCliMap{
			tenantCli: make(map[string]*local.MongoClient),
		},
		dbClientMap: nil,
		tm:          new(local.ShardingTxnManager),
	}

	masterMongo, err := local.NewMongo(masterCli, new(local.TxnManager), &local.MongoCliConf{IDGenStep: 1},
		&local.MongoOptions{IgnoreTenant: true})
	if err != nil {
		return nil, nil, fmt.Errorf("new master mongo db client failed, err: %v", err)
	}

	// get sharding db config
	shardingConf, err := getShardingDBConfig(context.Background(), masterMongo)
	if err != nil {
		return nil, nil, err
	}

	// fill mongo client info
	clientInfo.masterCli.SetUUID(shardingConf.MasterDB)
	clientInfo.dbClientMap = map[string]*local.MongoClient{shardingConf.MasterDB: clientInfo.masterCli}
	for slaveUUID, mongoConf := range shardingConf.SlaveDB {
		// decrypt slave mongodb uri
		mongoConf.URI, err = crypto.Decrypt(mongoConf.URI)
		if err != nil {
			return nil, nil, fmt.Errorf("decrypt %s slave mongo uri failed, err: %v", slaveUUID, err)
		}

		client, err := local.NewMongoClient(false, slaveUUID, &mongoConf, timeout)
		if err != nil {
			return nil, nil, fmt.Errorf("new %s slave mongo client failed, err: %v", slaveUUID, err)
		}
		clientInfo.dbClientMap[slaveUUID] = client
	}

	newDataCli, exists := clientInfo.dbClientMap[shardingConf.ForNewData]
	if !exists {
		return nil, nil, fmt.Errorf("add new tenant db %s config not found", shardingConf.ForNewData)
	}
	clientInfo.newDataCli = newDataCli

	return clientInfo, masterMongo, nil
}

// newTenantDB new db client for tenant
func (c *shardingMongoClient) newTenantDB(tenantID string, conf *local.MongoCliConf) local.DB {
	if tenantID == "" {
		return local.NewErrDB(errors.New("tenant is not set"))
	}

	client, exists := c.tenantCli.get(tenantID)
	if !exists {
		return local.NewErrDB(fmt.Errorf("tenant %s not exists", tenantID))
	}

	if client.Disabled() {
		return local.NewErrDB(fmt.Errorf("db client %s is disabled", client.UUID()))
	}

	txnManager, err := c.tm.DB(client.UUID())
	if err != nil {
		return local.NewErrDB(err)
	}

	db, err := local.NewMongo(client, txnManager, conf, &local.MongoOptions{Tenant: tenantID})
	if err != nil {
		return local.NewErrDB(err)
	}
	return db
}

// newIgnoreTenantDB new master db client that do not use tenant
func (c *shardingMongoClient) newIgnoreTenantDB(conf *local.MongoCliConf) local.DB {
	txnManager, err := c.tm.DB(c.masterCli.UUID())
	if err != nil {
		return local.NewErrDB(err)
	}

	db, err := local.NewMongo(c.masterCli, txnManager, conf, &local.MongoOptions{IgnoreTenant: true})
	if err != nil {
		return local.NewErrDB(err)
	}
	return db
}

// ping all sharding db clients
func (c *shardingMongoClient) ping() error {
	for uuid, client := range c.dbClientMap {
		err := client.Client().Ping(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("ping db %s failed, err: %v", uuid, err)
		}
	}
	return nil
}

// execForAllDB execute handler for all db clients
func (c *shardingMongoClient) execForAllDB(handler func(db local.DB) error, conf *local.MongoCliConf) error {
	for uuid, client := range c.dbClientMap {
		txnManager, err := c.tm.DB(client.UUID())
		if err != nil {
			return fmt.Errorf("get txn manager failed, err: %v", err)
		}

		db, err := local.NewMongo(client, txnManager, conf, &local.MongoOptions{IgnoreTenant: true})
		if err != nil {
			return fmt.Errorf("generate %s db client failed, err: %v", uuid, err)
		}

		if err = handler(db); err != nil {
			return fmt.Errorf("execute for db %s failed, err: %v", uuid, err)
		}
	}
	return nil
}

// getShardingDBConfig get sharding db config
func getShardingDBConfig(ctx context.Context, c *local.Mongo) (*ShardingDBConf, error) {
	cond := map[string]any{common.MongoMetaID: common.ShardingDBConfID}
	conf := new(ShardingDBConf)
	err := c.Table(common.BKTableNameSystem).Find(cond).One(ctx, &conf)
	if err != nil {
		if !c.IsNotFoundError(err) {
			return nil, fmt.Errorf("get sharding db config failed, err: %v", err)
		}

		// generate new sharding db config and save it if not exists, new tenant will be added to master db by default
		newUUID := uuid.NewString()
		conf = &ShardingDBConf{
			ID:         common.ShardingDBConfID,
			MasterDB:   newUUID,
			ForNewData: newUUID,
			SlaveDB:    make(map[string]local.MongoConf),
		}
		if err = c.Table(common.BKTableNameSystem).Insert(ctx, conf); err != nil {
			return nil, fmt.Errorf("insert new sharding db config failed, err: %v", err)
		}
		return conf, nil
	}
	blog.Infof("slave mongo config: %+v", conf.SlaveDB)

	return conf, nil
}

// refreshTenantDBMap refresh tenant to db relation
func (m *ShardingMongoManager) refreshTenantDBMap() error {
	tenantDBMap := make(map[string]string)
	for _, relation := range tenant.GetAllTenants() {
		tenantDBMap[relation.TenantID] = relation.Database
	}

	tenantCli := make(map[string]*local.MongoClient)
	for tenant, db := range tenantDBMap {
		client, exists := m.dbClientMap[db]
		if !exists {
			return fmt.Errorf("tenant %s related db %s config not found", tenant, db)
		}
		tenantCli[tenant] = client
	}

	m.tenantCli.tenantCli = tenantCli
	return nil
}

// Shard returns the sharded db client
func (m *ShardingMongoManager) Shard(opt ShardOpts) local.DB {
	if opt.IsIgnoreTenant() {
		return m.IgnoreTenant()
	}
	return m.Tenant(opt.Tenant())
}

// NewTenantCli returns the new tenant db client
func (m *ShardingMongoManager) NewTenantCli(tenant string) (local.DB, string, error) {
	client := m.newDataCli
	txnManager, err := m.tm.DB(client.UUID())
	if err != nil {
		return nil, "", err
	}

	db, err := local.NewMongo(client, txnManager, m.conf, &local.MongoOptions{Tenant: tenant})
	if err != nil {
		return nil, "", err
	}
	return db, m.newDataCli.UUID(), nil
}

// Tenant returns the db client for tenant
func (m *ShardingMongoManager) Tenant(tenant string) local.DB {
	return m.shardingMongoClient.newTenantDB(tenant, m.conf)
}

// IgnoreTenant returns the master db client that do not use tenant
func (m *ShardingMongoManager) IgnoreTenant() local.DB {
	return m.shardingMongoClient.newIgnoreTenantDB(m.conf)
}

// InitTxnManager TxnID management of initial transaction
func (m *ShardingMongoManager) InitTxnManager(r redis.Client) error {
	return m.tm.InitTxnManager(r)
}

// Ping all sharding db clients
func (m *ShardingMongoManager) Ping() error {
	return m.shardingMongoClient.ping()
}

// ExecForAllDB execute handler for all db clients
func (m *ShardingMongoManager) ExecForAllDB(handler func(db local.DB) error) error {
	return m.shardingMongoClient.execForAllDB(handler, m.conf)
}

// WatchMongo is the watch mongo db manager
type WatchMongo struct {
	*shardingMongoClient
	// dbWatchDBMap is the db uuid to watch db uuid map
	dbWatchDBMap map[string]string
}

// NewWatchMongo returns new watch mongo db manager
func NewWatchMongo(config local.MongoConf, timeout time.Duration, crypto cryptor.Cryptor) (ShardingDB, error) {
	clientInfo, masterMongo, err := newShardingMongoClient(config, timeout, crypto)
	if err != nil {
		return nil, err
	}

	sharding := &WatchMongo{
		shardingMongoClient: clientInfo,
		dbWatchDBMap:        make(map[string]string),
	}

	// generate db uuid to watch db uuid map
	relations := make([]WatchDBRelation, 0)
	err = masterMongo.Table(common.BKTableNameWatchDBRelation).Find(nil).All(context.Background(), &relations)
	if err != nil {
		return nil, fmt.Errorf("get db and watch db relation failed, err: %v", err)
	}

	for _, relation := range relations {
		sharding.dbWatchDBMap[relation.DB] = relation.WatchDB
	}

	// refresh tenant to db relation
	if err = sharding.refreshTenantDBMap(); err != nil {
		return nil, err
	}

	tenantChan := tenant.NewTenantEventChan(fmt.Sprintf("watch_sharding_db_%s", clientInfo.masterCli.UUID()))
	go func() {
		for e := range tenantChan {
			switch e.EventType {
			case tenant.Create:
				watchDBUUID, exists := sharding.dbWatchDBMap[e.Tenant.Database]
				if !exists {
					blog.Errorf("tenant %s db %s watch db config not found, use default watch db: %s",
						e.Tenant.TenantID, e.Tenant.Database, clientInfo.newDataCli.UUID())

					sharding.tenantCli.set(e.Tenant.TenantID, clientInfo.newDataCli)
					continue
				}
				client, exists := sharding.dbClientMap[watchDBUUID]
				if !exists {
					blog.Errorf("tenant %s related watch db %s config not found", e.Tenant.TenantID, watchDBUUID)
					continue
				}
				sharding.tenantCli.set(e.Tenant.TenantID, client)
			case tenant.Delete:
				sharding.tenantCli.delete(e.Tenant.TenantID)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Minute)
			if err = sharding.refreshTenantDBMap(); err != nil {
				blog.Errorf("refresh tenant to db relation failed, err: %v", err)
				continue
			}
		}
	}()

	return sharding, nil
}

func (m *WatchMongo) refreshTenantDBMap() error {
	tenantDBMap := make(map[string]string)
	for _, relation := range tenant.GetAllTenants() {
		watchDBUUID, exists := m.dbWatchDBMap[relation.Database]
		if exists {
			tenantDBMap[relation.TenantID] = watchDBUUID
		} else {
			blog.Warnf("tenant %s related db %s watch db not found, use default watch db %s", relation.TenantID,
				relation.Database, m.newDataCli.UUID())
			tenantDBMap[relation.TenantID] = m.newDataCli.UUID()
		}
	}

	tenantCli := make(map[string]*local.MongoClient)
	for tenant, db := range tenantDBMap {
		client, exists := m.dbClientMap[db]
		if !exists {
			return fmt.Errorf("tenant %s related db %s config not found", tenant, db)
		}
		tenantCli[tenant] = client
	}

	m.tenantCli.tenantCli = tenantCli
	return nil
}

// Shard returns the sharded db client
func (m *WatchMongo) Shard(opt ShardOpts) local.DB {
	if opt.IsIgnoreTenant() {
		return m.IgnoreTenant()
	}
	return m.Tenant(opt.Tenant())
}

// Tenant returns the db client for tenant
func (m *WatchMongo) Tenant(tenant string) local.DB {
	return m.shardingMongoClient.newTenantDB(tenant, &local.MongoCliConf{IDGenStep: 1})
}

// IgnoreTenant returns the master db client that do not use tenant
func (m *WatchMongo) IgnoreTenant() local.DB {
	return m.shardingMongoClient.newIgnoreTenantDB(&local.MongoCliConf{IDGenStep: 1})
}

// Ping all sharding db clients
func (m *WatchMongo) Ping() error {
	return m.shardingMongoClient.ping()
}

// ExecForAllDB execute handler for all db clients
func (m *WatchMongo) ExecForAllDB(handler func(db local.DB) error) error {
	return m.shardingMongoClient.execForAllDB(handler, &local.MongoCliConf{IDGenStep: 1})
}

// InitTxnManager TxnID management of initial transaction
func (m *WatchMongo) InitTxnManager(_ redis.Client) error {
	return fmt.Errorf("watch db do not support transaction")
}

// CommitTransaction commit transaction
func (m *WatchMongo) CommitTransaction(ctx context.Context, cap *metadata.TxnCapable) error {
	return fmt.Errorf("watch db do not support transaction")
}

// AbortTransaction abort transaction
func (m *WatchMongo) AbortTransaction(context.Context, *metadata.TxnCapable) (bool, error) {
	return false, fmt.Errorf("watch db do not support transaction")
}

// NewTenantCli returns the new tenant db client
func (m *WatchMongo) NewTenantCli(tenant string) (local.DB, string, error) {
	db, err := local.NewMongo(m.newDataCli, new(local.TxnManager), &local.MongoCliConf{IDGenStep: 1},
		&local.MongoOptions{Tenant: tenant})
	if err != nil {
		return nil, "", err
	}
	return db, m.newDataCli.UUID(), nil
}
