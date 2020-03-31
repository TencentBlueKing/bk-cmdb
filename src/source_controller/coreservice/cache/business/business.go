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

package business

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/reflector"
	"configcenter/src/storage/stream/types"
	"github.com/tidwall/gjson"
	"gopkg.in/redis.v5"
)

func NewBusinessCache(event reflector.Interface, rds *redis.Client) *BusinessCache {
	return &BusinessCache{
		client: rds, 
		event:event,
		hkeyName: common.BKCacheKeyV3Prefix + "bizlist",
		expireKeyName: "expireTime",
		lockKeyName: common.BKCacheKeyV3Prefix + "bizlistlock",
	}
	
}

type BusinessCache struct {
	client *redis.Client
	event  reflector.Interface
	hkeyName string
	expireKeyName string
	lockKeyName string
	db dal.DB
}

func (bc *BusinessCache) Run() error {
	
	
	return nil
}

func (bc *BusinessCache) onUpsert(e *types.Event) {
	blog.V(4).Infof("received biz upsert event, oid: %s, operate: %s, doc: %s", e.Oid, e.OperationType, e.DocBytes)
	
	bizInfo := gjson.GetManyBytes(e.DocBytes, "bk_biz_id", "bk_biz_name")
	bizID := bizInfo[0].Int()
	bizName := bizInfo[1].String()
	if bizID == 0 {
		blog.Errorf("received biz upsert event, invalid biz id: %d, oid: %s", bizID, e.Oid)
		return
	}
	
	if len(bizName) == 0 {
		blog.Errorf("received biz upsert event, invalid biz name: %s, oid: %s", bizName, e.Oid)
		return
	}
	bc.upsertBusinessCache(bizID, bizName, e.DocBytes)
}

func (bc *BusinessCache) onDelete(e *types.Event) {
	blog.V(4).Infof("received *unexpected* delete biz event, oid: %s, operate: %s, doc: %s", e.Oid, e.OperationType, e.DocBytes)
}

func (bc *BusinessCache) onListDone() {
	// set the expire key for later use
	if err := bc.client.HSet(bc.hkeyName, bc.expireKeyName, time.Now().Unix()).Err(); err != nil {
		blog.Errorf("list biz cache done, but set bizlist expire key %s failed, err: %v", bc.expireKeyName, err)
		return
	}
}

func (bc *BusinessCache) upsertBusinessCache(id int64, name string, bizInfo []byte) {
	keyName := bc.genKey(id, name)
	
	// set the business
	if err := bc.client.HSet(bc.hkeyName, keyName, bizInfo).Err(); err != nil {
		blog.Errorf("upsert biz id: %d, name: %s cache failed, err: %v", id, name, err)
		return
	}
	
	
	// check if the business is already exist or not
	key, exist := bc.isExist(id)
	if exist {
		if keyName != key {
			// business name has changed, delete and reset
			if err := bc.client.HDel(bc.hkeyName, key).Err(); err != nil {
				blog.Errorf("delete invalid biz cache, key: %s failed, err: %v", key, err)
			}
		}
	}
}

func (bc *BusinessCache) isExist(bizID int64) (string, bool) {
	// get all keys which contains the biz id.
	keys, err := bc.client.HKeys(bc.hkeyName).Result()
	if err != nil {
		blog.Errorf("hget bizlist keys %s falied. err: %v", bc.hkeyName, err)
		return "", false
	}
	for _, key := range keys {
		if key == bc.expireKeyName {
			// skip the expire key
			continue
		}
		id, _ ,err := bc.parseKey(key)
		if err != nil {
			// invalid key, delete immediately
			if bc.client.HDel(bc.hkeyName, key).Err() != nil {
				blog.Errorf("delete invalid biz hash %s key: %s failed,", bc.hkeyName, key)
			}
			return "", false
		}
		if id == bizID {
			return key, true
		}
	}
	
	return "", false
}

func (bc *BusinessCache) getKeys() (keys []string, err error) {
	// get all keys which contains the biz id.
	keys, err := bc.client.HKeys(bc.hkeyName).Result()
	if err != nil {
		blog.Errorf("hget bizlist keys %s falied. err: %v", bc.hkeyName, err)
		return nil, err
	}
	for _, key := range keys {
		if key == bc.expireKeyName {
			// skip the expire key
			continue
		}
		id, _ ,err := bc.parseKey(key)
		if err != nil {
			// invalid key, delete immediately
			if bc.client.HDel(bc.hkeyName, key).Err() != nil {
				blog.Errorf("delete invalid biz hash %s key: %s failed,", bc.hkeyName, key)
			}
			return "", false
		}
		if id == bizID {
			return key, true
		}
	}

	return "", false
}


func (bc *BusinessCache) parseKey(key string)(int64, string, error) {
	index := strings.Index(key,":")
	if index == -1 {
		return 0, "", errors.New("invalid key")
	}
	
	bizID, err := strconv.ParseInt(key[:index], 10, 64)
	if err != nil {
		return 0, "", errors.New("key with invalid biz id")
	}
	
	bizName := key[index+1:]
	if len(bizName) == 0 {
		return 0, "", errors.New("invalid key with empty biz name")
	}
	
	return bizID, bizName, nil
}

func (bc *BusinessCache) genKey(bizID int64, name string) string {
	return strconv.FormatInt(bizID, 10) + ":" + name
}

func (bc *BusinessCache) tryUpdate() {
	
}