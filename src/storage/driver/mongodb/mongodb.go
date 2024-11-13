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

// Package mongodb is cmdb mongodb driver
package mongodb

import (
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/cryptor"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/sharding"
	dbType "configcenter/src/storage/dal/types"
)

type mongoClient interface {
	GetDal(prefix string) dal.Dal
	SetDal(prefix string, db dal.Dal)
	RemoveDal(prefix string)
	Healthz() []metric.HealthItem
}

var mongoInst mongoClient

func init() {
	mongoInst = &mongodb{
		dalMap: make(map[string]dal.Dal),
	}
}

type mongodb struct {
	dalLock sync.RWMutex
	dalMap  map[string]dal.Dal
}

// GetDal get dal by prefix
func (m *mongodb) GetDal(prefix string) dal.Dal {
	m.dalLock.RLock()
	defer m.dalLock.RUnlock()
	return m.dalMap[prefix]
}

// SetDal set dal with prefix
func (m *mongodb) SetDal(prefix string, db dal.Dal) {
	m.dalLock.Lock()
	defer m.dalLock.Unlock()
	m.dalMap[prefix] = db
}

// RemoveDal remove dal by prefix
func (m *mongodb) RemoveDal(prefix string) {
	m.dalLock.Lock()
	defer m.dalLock.Unlock()
	delete(m.dalMap, prefix)
}

// Healthz check db health status
func (m *mongodb) Healthz() []metric.HealthItem {
	m.dalLock.RLock()
	defer m.dalLock.RUnlock()

	items := make([]metric.HealthItem, 0)
	for prefix, db := range m.dalMap {
		if db == nil {
			items = append(items, metric.HealthItem{
				IsHealthy: false,
				Message:   prefix + " db not initialized",
				Name:      types.CCFunctionalityMongo,
			})
			continue
		}

		if err := db.Ping(); err != nil {
			items = append(items, metric.HealthItem{
				IsHealthy: false,
				Message:   prefix + " db connect error. err: " + err.Error(),
				Name:      types.CCFunctionalityMongo,
			})
			continue
		}

		items = append(items, metric.HealthItem{
			IsHealthy: true,
			Name:      types.CCFunctionalityMongo,
		})
	}
	return items
}

// Dal get mongodb sharding client
func Dal(prefix ...string) dal.Dal {
	var pre string
	if len(prefix) > 0 {
		pre = prefix[0]
	}
	return mongoInst.GetDal(pre)
}

// Shard get sharded mongodb client
func Shard(opt sharding.ShardOpts) local.DB {
	return Dal().Shard(opt)
}

// SetShardingCli set mongodb sharding client with prefix
func SetShardingCli(prefix string, config *mongo.Config, cryptoConf *cryptor.Config) error {
	crypto, err := cryptor.NewCrypto(cryptoConf)
	if err != nil {
		blog.Errorf("new %s mongo crypto failed, err: %v", prefix, err)
		return errors.NewCCError(common.CCErrCommResourceInitFailed, "init mongo crypto failed")
	}

	shardingDB, err := sharding.NewShardingMongo(config.GetMongoConf(), time.Minute, crypto)
	if err != nil {
		blog.Errorf("new %s sharding mongo client failed, err: %v", prefix, err)
		return errors.NewCCError(common.CCErrCommResourceInitFailed, "init sharding mongo client failed")
	}
	mongoInst.SetDal(prefix, shardingDB)
	return nil
}

// SetDisableDBShardingCli set mongodb client that disables db sharding with prefix
func SetDisableDBShardingCli(prefix string, config *mongo.Config) error {
	shardingDB, err := sharding.NewDisableDBShardingMongo(config.GetMongoConf(), time.Minute)
	if err != nil {
		blog.Errorf("new %s disable db sharding mongo client failed, err: %v", prefix, err)
		return errors.NewCCError(common.CCErrCommResourceInitFailed, "init disable db sharding mongo client failed")
	}
	mongoInst.SetDal(prefix, shardingDB)
	return nil
}

// 暂时不支持热更新，所以没有加锁
// TODO remove this after all mongodb clients use sharding
var (
	dbMap = make(map[string]dal.DB)

	// 在并发的情况下，这里存在panic的问题
	lastInitErr   errors.CCErrorCoder
	lastConfigErr errors.CCErrorCoder
)

// Client  get default error
func Client(prefix ...string) dal.DB {
	var pre string
	if len(prefix) > 0 {
		pre = prefix[0]
	}
	return dbMap[pre]
}

// Table 获取操作db table的对象
func Table(name string) dbType.Table {
	return Client().Table(name)
}

// InitClient init mongodb client
func InitClient(prefix string, config *mongo.Config) errors.CCErrorCoder {
	lastInitErr = nil
	var dbErr error
	dbMap[prefix], dbErr = local.NewMgo(config.GetMongoConf(), time.Minute)
	if dbErr != nil {
		blog.Errorf("failed to connect the mongo server, error info is %s", dbErr.Error())
		lastInitErr = errors.NewCCError(common.CCErrCommResourceInitFailed,
			"'"+prefix+".mongodb' initialization failed")
		return lastInitErr
	}
	return nil
}

// Healthz check db health status
func Healthz() []metric.HealthItem {
	items := make([]metric.HealthItem, 0)

	for prefix, db := range dbMap {
		if db == nil {
			items = append(items, metric.HealthItem{
				IsHealthy: false,
				Message:   prefix + " db not initialized",
				Name:      types.CCFunctionalityMongo,
			})
			continue
		}

		if err := db.Ping(); err != nil {
			items = append(items, metric.HealthItem{
				IsHealthy: false,
				Message:   prefix + " db connect error. err: " + err.Error(),
				Name:      types.CCFunctionalityMongo,
			})
			continue
		}

		items = append(items, metric.HealthItem{
			IsHealthy: true,
			Name:      types.CCFunctionalityMongo,
		})
	}

	return append(items, mongoInst.Healthz()...)
}

// IsDuplicatedError check duplicated error
func IsDuplicatedError(err error) bool {
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
	return err == dbType.ErrDuplicated
}

// IsNotFoundError check the not found error
func IsNotFoundError(err error) bool {
	return err == dbType.ErrDocumentNotFound
}

// GetDuplicateKey get duplicate key from error, if the error is not a duplicate error, returns the raw error message
// mongodb raw error format example:
// ...{E11000 duplicate key error collection: cmdb.cc_ObjectBase_... index: bkcc_unique_... dup key:
// { bk_inst_name: \"xxx\" }}]},...
func GetDuplicateKey(err error) string {
	if err == nil {
		return ""
	}

	errString := err.Error()
	if !strings.Contains(errString, "E11000 duplicate") {
		return errString
	}

	start := strings.Index(errString, "dup key: ")
	if start == -1 {
		return errString
	}
	start += len("dup key: ") + 1

	end := strings.LastIndex(errString, "}]")
	if end == -1 || end < start {
		return errString
	}

	return errString[start:end]
}

// GetDuplicateValue get duplicate Value from error, if it is not a duplicate error, returns the raw error message
// mongodb raw error format example:
// Index build failed: ... E11000 duplicate key error collection: cmdb.cc_ObjectBase_0_pub_...:  dup key:
// dup key: { field: "xxxx" }
func GetDuplicateValue(field string, err error) string {
	if field == "" {
		return ""
	}
	if err == nil {
		return ""
	}

	errString := err.Error()
	if !strings.Contains(errString, "E11000 duplicate") {
		return errString
	}

	start := strings.Index(errString, "dup key: ")
	if start == -1 {
		return errString
	}
	start += len("dup key: { " + field + ": ")

	end := strings.LastIndex(errString, " }")
	if end == -1 || end < start {
		return errString
	}

	return errString[start:end]
}
