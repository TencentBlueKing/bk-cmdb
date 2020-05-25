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

package redis

import (
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/types"
	dalRedis "configcenter/src/storage/dal/redis"

	redis "gopkg.in/redis.v5"
)

/*
 暂时不支持，多个mongodb实例连接， 暂时不值热更新，所以没有加锁
*/

var (
	defaultPrefix string = "redis"
	// redis is default
	cacheMap = make(map[string]*redis.Client, 0)
	lock     = &sync.RWMutex{}

	// 在并发的情况下，这里存在panic的问题
	lastInitErr   errors.CCErrorCoder
	lastConfigErr errors.CCErrorCoder

	Nil = redis.Nil
)

// Client  get default error
func Client() *redis.Client {
	lock.RLock()
	defer lock.RUnlock()
	if db, ok := cacheMap[defaultPrefix]; ok {
		return db
	}
	return nil
}

// ClientInstance  获取指定的redis
func ClientInstance(prefix string) *redis.Client {
	lock.RLock()
	defer lock.RUnlock()
	if db, ok := cacheMap[prefix]; ok {
		return db
	}
	return nil
}

func ParseConfig(prefix string, configMap map[string]string) (*dalRedis.Config, errors.CCErrorCoder) {
	lastConfigErr = nil
	config := dalRedis.ParseConfigFromKV(prefix, configMap)
	if config.Address == "" {
		lastConfigErr = errors.NewCCError(common.CCErrCommConfMissItem, "Configuration file missing ["+prefix+".Address] configuration item")
		return nil, lastConfigErr
	}
	if config.Password == "" {
		lastConfigErr = errors.NewCCError(common.CCErrCommConfMissItem, "Configuration file missing ["+prefix+".pwd] configuration item")
		return nil, lastConfigErr
	}

	return &config, nil
}

func InitClient(prefix string, config *dalRedis.Config) errors.CCErrorCoder {
	lock.Lock()
	defer lock.Unlock()
	if cacheMap[prefix] != nil {
		// 不支持热更新
		blog.V(5).Infof("duplicate open redis. prefix:%s, host:%s", prefix, config.Address)
		return nil
	}
	lastInitErr = nil
	db, dbErr := dalRedis.NewFromConfig(*config)
	if dbErr != nil {
		blog.Errorf("failed to connect the redis server, error info is %s", dbErr.Error())
		lastInitErr = errors.NewCCError(common.CCErrCommResourceInitFailed, "'"+prefix+" redis' initialization failed")
		return lastInitErr
	}
	if len(cacheMap) == 0 {
		cacheMap[defaultPrefix] = db
	}
	cacheMap[prefix] = db
	return nil
}

func Validate() errors.CCErrorCoder {
	return nil
}

func UpdateConfig(prefix string, config dalRedis.Config) {
	// 不支持热更行
	return
}

func Healthz() (items []metric.HealthItem) {
	lock.RLock()
	defer lock.RUnlock()

	for prefix, db := range cacheMap {
		item := &metric.HealthItem{
			IsHealthy: true,
			Name:      types.CCFunctionalityRedis + " " + prefix,
		}
		items = append(items, *item)
		if db == nil {
			item.IsHealthy = false
			item.Message = "[" + prefix + "] not initialized"
			continue
		}
		if err := db.Ping().Err(); err != nil {
			item.IsHealthy = false
			item.Message = "[" + prefix + "] connect error. err: " + err.Error()
			continue
		}
	}
	if len(items) == 0 {
		items = append(items, metric.HealthItem{
			IsHealthy: false,
			Name:      "not found intance",
		})
	}
	return
}
