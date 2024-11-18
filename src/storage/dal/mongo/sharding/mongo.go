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
	"time"

	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/cryptor"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"

	"github.com/google/uuid"
)

// ShardingMongoManager is the sharding db manager for mongo
type ShardingMongoManager struct {
	// masterCli is the client for master mongodb, master mongodb stores the platform data and some tenant data
	masterCli *local.MongoClient
	// newTenantCli is the client for mongodb that new tenant data will be stored into
	newTenantCli *local.MongoClient
	// tenantCli is the tenant id to mongodb client map
	tenantCli map[string]*local.MongoClient
	// dbClientMap is the db uuid to mongodb client map
	dbClientMap map[string]*local.MongoClient
	// tm is the transaction manager
	tm *local.ShardingTxnManager
	// conf is the mongo client config
	conf *local.MongoCliConf
}

// NewShardingMongo returns new sharding db manager for mongo
func NewShardingMongo(config local.MongoConf, timeout time.Duration, crypto cryptor.Cryptor) (ShardingDB, error) {
	// connect master mongodb
	masterCli, err := local.NewMongoClient(true, "", config, timeout)
	if err != nil {
		return nil, fmt.Errorf("new master mongo client failed, err: %v", err)
	}

	sharding := &ShardingMongoManager{
		masterCli: masterCli,
		tenantCli: make(map[string]*local.MongoClient),
		tm:        new(local.ShardingTxnManager),
		conf:      &local.MongoCliConf{DisableInsert: config.DisableInsert},
	}

	masterMongo, err := local.NewMongo(masterCli, new(local.TxnManager), sharding.conf,
		&local.MongoOptions{IgnoreTenant: true})
	if err != nil {
		return nil, fmt.Errorf("new master mongo db client failed, err: %v", err)
	}

	ctx := context.Background()
	sharding.conf.IDGenStep, err = masterMongo.InitIDGenerator(ctx)
	if err != nil {
		return nil, err
	}

	// get sharding db config
	shardingConf, err := getShardingDBConfig(ctx, masterMongo)
	if err != nil {
		return nil, err
	}

	// fill mongo client info
	sharding.masterCli.SetUUID(shardingConf.MasterDB)

	sharding.dbClientMap = map[string]*local.MongoClient{shardingConf.MasterDB: sharding.masterCli}
	for slaveUUID, mongoConf := range shardingConf.SlaveDB {
		// decrypt slave mongodb uri
		mongoConf.URI, err = crypto.Decrypt(mongoConf.URI)
		if err != nil {
			return nil, fmt.Errorf("decrypt %s slave mongo uri failed, err: %v", slaveUUID, err)
		}

		client, err := local.NewMongoClient(false, slaveUUID, mongoConf, timeout)
		if err != nil {
			return nil, fmt.Errorf("new %s slave mongo client failed, err: %v", slaveUUID, err)
		}
		sharding.dbClientMap[slaveUUID] = client
	}

	newTenantCli, exists := sharding.dbClientMap[shardingConf.ForNewTenant]
	if !exists {
		return nil, fmt.Errorf("add new tenant db %s config not found", shardingConf.ForNewTenant)
	}
	sharding.newTenantCli = newTenantCli

	err = tenant.Init(&tenant.Options{DB: sharding.IgnoreTenant()})
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

	return sharding, nil
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
			ID:           common.ShardingDBConfID,
			MasterDB:     newUUID,
			ForNewTenant: newUUID,
			SlaveDB:      make(map[string]local.MongoConf),
		}
		if err = c.Table(common.BKTableNameSystem).Insert(ctx, conf); err != nil {
			return nil, fmt.Errorf("insert new sharding db config failed, err: %v", err)
		}
		return conf, nil
	}
	blog.Infof("slave mongo config: %+v", conf.SlaveDB)

	return conf, nil
}

func (m *ShardingMongoManager) refreshTenantDBMap() error {
	tenantDBMap := make(map[string]string)
	for _, relation := range tenant.GetAllTenants() {
		tenantDBMap[relation.TenantID] = relation.Database
	}

	tenantCli := make(map[string]*local.MongoClient)
	for tenant, db := range tenantDBMap {
		// TODO add default tenant db client for compatible, remove this later
		if tenant == common.BKDefaultOwnerID {
			tenantCli[tenant] = m.masterCli
			continue
		}

		client, exists := m.dbClientMap[db]
		if !exists {
			return fmt.Errorf("tenant %s related db %s config not found", tenant, db)
		}
		tenantCli[tenant] = client
	}

	m.tenantCli = tenantCli
	return nil
}

// Shard returns the sharded db client
func (m *ShardingMongoManager) Shard(opt ShardOpts) local.DB {
	if opt.IsIgnoreTenant() {
		return m.IgnoreTenant()
	}
	return m.Tenant(opt.Tenant())
}

// Tenant returns the db client for tenant
func (m *ShardingMongoManager) Tenant(tenant string) local.DB {
	if tenant == "" {
		return local.NewErrDB(errors.New("tenant is not set"))
	}

	client, exists := m.tenantCli[tenant]
	if !exists {
		return local.NewErrDB(fmt.Errorf("tenant %s not exists", tenant))
	}

	if client.Disabled() {
		return local.NewErrDB(fmt.Errorf("db client %s is disabled", client.UUID()))
	}

	txnManager, err := m.tm.Tenant(false, tenant)
	if err != nil {
		return local.NewErrDB(err)
	}

	db, err := local.NewMongo(client, txnManager, m.conf, &local.MongoOptions{Tenant: tenant})
	if err != nil {
		return local.NewErrDB(err)
	}
	return db
}

// IgnoreTenant returns the master db client that do not use tenant
func (m *ShardingMongoManager) IgnoreTenant() local.DB {
	txnManager, err := m.tm.Tenant(true, "")
	if err != nil {
		return local.NewErrDB(err)
	}

	db, err := local.NewMongo(m.masterCli, txnManager, m.conf, &local.MongoOptions{IgnoreTenant: true})
	if err != nil {
		return local.NewErrDB(err)
	}
	return db
}

// InitTxnManager TxnID management of initial transaction
func (m *ShardingMongoManager) InitTxnManager(r redis.Client) error {
	return m.tm.InitTxnManager(r)
}

// Ping all sharding db clients
func (m *ShardingMongoManager) Ping() error {
	for uuid, client := range m.dbClientMap {
		err := client.Client().Ping(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("ping db %s failed, err: %v", uuid, err)
		}
	}
	return nil
}

// DisableDBShardingMongo is the disabled db sharding mongo db manager, right now only watch db sharding is disabled
type DisableDBShardingMongo struct {
	client *local.MongoClient
	tm     *local.TxnManager
	conf   *local.MongoCliConf
}

// NewDisableDBShardingMongo returns new disabled db sharding mongo db manager
func NewDisableDBShardingMongo(config local.MongoConf, timeout time.Duration) (ShardingDB, error) {
	client, err := local.NewMongoClient(true, "", config, timeout)
	if err != nil {
		return nil, fmt.Errorf("new mongo client failed, err: %v", err)
	}

	db := &DisableDBShardingMongo{
		client: client,
		tm:     new(local.TxnManager),
		conf:   &local.MongoCliConf{DisableInsert: config.DisableInsert},
	}

	masterMongo, err := local.NewMongo(client, new(local.TxnManager), db.conf, &local.MongoOptions{IgnoreTenant: true})
	if err != nil {
		return nil, fmt.Errorf("new master mongo db client failed, err: %v", err)
	}
	db.conf.IDGenStep, err = masterMongo.InitIDGenerator(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Shard returns the sharded db client
func (m *DisableDBShardingMongo) Shard(opt ShardOpts) local.DB {
	if opt.IsIgnoreTenant() {
		return m.IgnoreTenant()
	}
	return m.Tenant(opt.Tenant())
}

// Tenant returns the db client for tenant
func (m *DisableDBShardingMongo) Tenant(tenant string) local.DB {
	if tenant == "" {
		return local.NewErrDB(errors.New("tenant is not set"))
	}

	db, err := local.NewMongo(m.client, m.tm, m.conf, &local.MongoOptions{Tenant: tenant})
	if err != nil {
		return local.NewErrDB(err)
	}
	return db
}

// IgnoreTenant returns the master db client that do not use tenant
func (m *DisableDBShardingMongo) IgnoreTenant() local.DB {
	db, err := local.NewMongo(m.client, m.tm, m.conf, &local.MongoOptions{IgnoreTenant: true})
	if err != nil {
		return local.NewErrDB(err)
	}
	return db
}

// InitTxnManager TxnID management of initial transaction
func (m *DisableDBShardingMongo) InitTxnManager(r redis.Client) error {
	return m.tm.InitTxnManager(r)
}

// Ping db client
func (m *DisableDBShardingMongo) Ping() error {
	return m.client.Client().Ping(context.Background(), nil)
}
