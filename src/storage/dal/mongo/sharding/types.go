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

import "configcenter/src/storage/dal/mongo/local"

// ShardOpts is db sharding options
type ShardOpts interface {
	Tenant() string
	IsIgnoreTenant() bool
}

// ShardOptions is db sharding options
type ShardOptions struct {
	// tenantID is the tenant id, used for sharding by tenant
	tenantID string
	// ignoreTenant is used for platform operations that do not use tenant for sharding
	ignoreTenant bool
}

// NewShardOpts create a new db sharding options
func NewShardOpts() *ShardOptions {
	return new(ShardOptions)
}

// WithTenant set tenant id
func (s *ShardOptions) WithTenant(tenantID string) *ShardOptions {
	s.tenantID = tenantID
	return s
}

// WithIgnoreTenant set ignore tenant id = true
func (s *ShardOptions) WithIgnoreTenant() *ShardOptions {
	s.ignoreTenant = true
	return s
}

// Tenant returns tenant id
func (s *ShardOptions) Tenant() string {
	return s.tenantID
}

// IsIgnoreTenant returns whether tenant id is ignored
func (s *ShardOptions) IsIgnoreTenant() bool {
	return s.ignoreTenant
}

// ShardingDBConf is sharding mongodb config
type ShardingDBConf struct {
	ID           string                     `bson:"_id"`
	MasterDB     string                     `bson:"master_db"`
	ForNewTenant string                     `bson:"for_new_tenant"`
	SlaveDB      map[string]local.MongoConf `bson:"slave_db"`
}