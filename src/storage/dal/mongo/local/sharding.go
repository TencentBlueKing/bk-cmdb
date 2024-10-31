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
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/cryptor"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
)

// ShardingMongoManager is the sharding db manager for mongo
type ShardingMongoManager struct {
	// masterCli is the client for master mongodb, master mongodb stores the platform data and some tenant data
	masterCli *mongoClient
	// newTenantCli is the client for mongodb that new tenant data will be stored into
	newTenantCli *mongoClient
	// tenantCli is the tenant id to mongodb client map
	tenantCli map[string]*mongoClient
	// dbClientMap is the db uuid to mongodb client map
	dbClientMap map[string]*mongoClient
	// tm is the transaction manager
	tm *ShardingTxnManager
	// conf is the mongo client config
	conf *mongoCliConf
}

// NewShardingMongo returns new sharding db manager for mongo
func NewShardingMongo(config MongoConf, timeout time.Duration, crypto cryptor.Cryptor) (dal.ShardingDB, error) {
	// connect master mongodb
	masterCli, err := newMongoClient(true, "", config, timeout)
	if err != nil {
		return nil, fmt.Errorf("new master mongo client failed, err: %v", err)
	}

	sharding := &ShardingMongoManager{
		masterCli: masterCli,
		tenantCli: make(map[string]*mongoClient),
		tm:        new(ShardingTxnManager),
		conf:      &mongoCliConf{disableInsert: config.DisableInsert},
	}

	masterMongo := &Mongo{ignoreTenant: true, mongoClient: masterCli}

	ctx := context.Background()
	sharding.conf.idGenStep, err = masterMongo.initIDGenerator(ctx)
	if err != nil {
		return nil, err
	}

	// get sharding db config
	shardingConf, err := masterMongo.getShardingDBConfig(ctx)
	if err != nil {
		return nil, err
	}

	// fill mongo client info
	sharding.masterCli.uuid = shardingConf.MasterDB

	sharding.dbClientMap = map[string]*mongoClient{shardingConf.MasterDB: sharding.masterCli}
	for slaveUUID, mongoConf := range shardingConf.SlaveDB {
		// decrypt slave mongodb uri
		mongoConf.URI, err = crypto.Decrypt(mongoConf.URI)
		if err != nil {
			return nil, fmt.Errorf("decrypt %s slave mongo uri failed, err: %v", slaveUUID, err)
		}

		client, err := newMongoClient(false, slaveUUID, mongoConf, timeout)
		if err != nil {
			return nil, fmt.Errorf("new %s slave mongo client failed, err: %v", slaveUUID, err)
		}
		sharding.dbClientMap[slaveUUID] = client
	}

	newTenantCli, exists := sharding.dbClientMap[shardingConf.AddNewTenantDB]
	if !exists {
		return nil, fmt.Errorf("add new tenant db %s config not found", shardingConf.AddNewTenantDB)
	}
	sharding.newTenantCli = newTenantCli

	if err = sharding.refreshTenantDBMap(masterMongo); err != nil {
		return nil, err
	}

	// loop refresh tenant to db relation
	go func() {
		for {
			time.Sleep(time.Minute)
			if err = sharding.refreshTenantDBMap(masterMongo); err != nil {
				blog.Errorf("refresh tenant to db relation failed, err: %v", err)
				continue
			}
		}
	}()

	return sharding, nil
}

func (m *ShardingMongoManager) refreshTenantDBMap(client dal.DB) error {
	// get tenant to db map
	tenantDBMap := make(map[string]string)
	relations := make([]TenantDBRelation, 0)
	err := client.Table(common.BKTableNameTenantDBRelation).Find(nil).All(context.Background(), &relations)
	if err != nil {
		return fmt.Errorf("get tenant db relations failed, err: %v", err)
	}
	for _, relation := range relations {
		tenantDBMap[relation.TenantID] = relation.Database
	}

	tenantCli := make(map[string]*mongoClient)
	for tenant, db := range tenantDBMap {
		client, exists := m.dbClientMap[db]
		if !exists {
			return fmt.Errorf("tenant %s related db %s config not found", tenant, db)
		}
		tenantCli[tenant] = client
	}

	m.tenantCli = tenantCli
	return nil
}

// Tenant returns the db client for tenant
func (m *ShardingMongoManager) Tenant(tenant string) dal.DB {
	if tenant == "" {
		return &errDB{err: errors.New("tenant is not set")}
	}

	client, exists := m.tenantCli[tenant]
	if !exists {
		return &errDB{err: fmt.Errorf("tenant %s not exists", tenant)}
	}

	if client.disabled {
		return &errDB{err: fmt.Errorf("db client %s is disabled", client.uuid)}
	}

	txnManager, err := m.tm.Tenant(false, tenant)
	if err != nil {
		return &errDB{err: err}
	}

	return &Mongo{
		tenant:         tenant,
		mongoClient:    client,
		tm:             txnManager,
		conf:           m.conf,
		enableSharding: true,
	}
}

// IgnoreTenant returns the master db client that do not use tenant
func (m *ShardingMongoManager) IgnoreTenant() dal.DB {
	txnManager, err := m.tm.Tenant(true, "")
	if err != nil {
		return &errDB{err: err}
	}

	return &Mongo{
		ignoreTenant:   true,
		mongoClient:    m.masterCli,
		tm:             txnManager,
		conf:           m.conf,
		enableSharding: true,
	}
}

// InitTxnManager TxnID management of initial transaction
func (m *ShardingMongoManager) InitTxnManager(r redis.Client) error {
	return m.tm.InitTxnManager(r)
}
