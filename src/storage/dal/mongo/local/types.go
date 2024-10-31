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

// MongoConf is mongodb config
type MongoConf struct {
	Name     string `bson:"name"`
	Disabled bool   `bson:"disabled"`

	MaxOpenConns  uint64 `bson:"max_open_conns"`
	MaxIdleConns  uint64 `bson:"max_idle_conns"`
	URI           string `bson:"uri"`
	RsName        string `bson:"rs_name"`
	SocketTimeout int    `bson:"socket_timeout"`

	DisableInsert bool
}

// ShardingDBConf is sharding mongodb config
type ShardingDBConf struct {
	ID             string               `bson:"_id"`
	MasterDB       string               `bson:"master_db"`
	AddNewTenantDB string               `bson:"add_new_tenant_db"`
	SlaveDB        map[string]MongoConf `bson:"slave_db"`
}

// TenantDBRelation tenant to db uuid relation
type TenantDBRelation struct {
	TenantID string `bson:"tenant_id"`
	Database string `bson:"database"`
}
