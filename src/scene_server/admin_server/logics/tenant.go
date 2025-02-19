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

package logics

import (
	"fmt"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/cryptor"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/admin_server/app/options"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/sharding"
)

var (
	NewTenantDB   *local.MongoClient
	MongoConf     *local.MongoCliConf
	processConfig options.Config
)

// SetProcessConfig set process config
func SetProcessConfig(c options.Config) {
	processConfig = c
}

// GetNewTenantCli get new tenant db
func GetNewTenantCli(kit *rest.Kit) (*local.Mongo, error) {
	crypto, err := cryptor.NewCrypto(processConfig.Crypto)
	if err != nil {
		return nil, fmt.Errorf("new db index mongo crypto failed, err: %v", err)
	}

	// db 语句的执行时间设置为never timeout
	mongoConf := processConfig.MongoDB
	mongoConf.SocketTimeout = 0
	db, err := sharding.NewShardingMongo(mongoConf.GetMongoConf(), time.Minute, crypto)
	if err != nil {
		return nil, fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	NewTenantDB = db.NewTenantCli()
	MongoConf = db.DBConfig()
	if NewTenantDB == nil {
		blog.Errorf("tenant db not init, rid: %s", kit.Rid)
		return nil, fmt.Errorf("tenant db not init, rid: %s", kit.Rid)
	}

	dbCli, err := local.NewMongo(NewTenantDB, &local.TxnManager{}, MongoConf, &local.MongoOptions{Tenant: kit.TenantID})
	if err != nil {
		return nil, err
	}
	return dbCli, nil
}

// GetNewTenantDB get new tenant db
func GetNewTenantDB() (*local.MongoClient, error) {
	if NewTenantDB == nil {
		blog.Errorf("new tenant db not init")
		return nil, fmt.Errorf("new tenant db not init")
	}
	return NewTenantDB, nil
}
