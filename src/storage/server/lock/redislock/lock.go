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

package redislock

import (
	"encoding/json"
	"fmt"

	redis "gopkg.in/redis.v5"

	"configcenter/src/common/blog"
	"configcenter/src/storage/server/lock/types"
)

// RedisLock  redis lock private data
type RedisLock struct {
	storage        *redis.Client
	prefix         string
	noticeChan     chan *notice
	errRetryCnt    int
	triggerTime    int // compensation trigger time unit sencod
	lockCompare    CompareFunc
	preLcokCompare CompareFunc
}

// NewLock create lock
func NewLock(client *redis.Client, prefix string, retry, compensationTime int, lockCompare, preLcokCompare CompareFunc) *RedisLock {
	rl := &RedisLock{
		storage:        client,
		prefix:         prefix,
		noticeChan:     make(chan *notice, noticeMaxCount),
		errRetryCnt:    retry,
		triggerTime:    compensationTime,
		lockCompare:    lockCompare,
		preLcokCompare: preLcokCompare,
	}
	go rl.compensation()
	return rl
}

// PreLock  lock resources, exclusive mode
func (rl *RedisLock) PreLock(meta *types.Lock) (locked bool, err error) {
	if nil == meta {
		return false, fmt.Errorf("params not allowed null")
	}
	lockKey := getFmtRedisKey(rl.prefix, meta.LockName, true, false)
	lockRelationKey := getFmtRedisKey(rl.prefix, meta.TxnID, true, true)

	locked, _, err = rl.lock(lockKey, lockRelationKey, meta, rl.preLcokCompare)

	return locked, err
}

// PreUnlock  unlock resources, exclusive mode
func (rl *RedisLock) PreUnlock(meta *types.Lock) error {
	if nil == meta {
		return fmt.Errorf("params not allowed null")
	}
	lockKey := getFmtRedisKey(rl.prefix, meta.LockName, true, false)
	lockRelationKey := getFmtRedisKey(rl.prefix, meta.TxnID, true, true)

	return rl.unlock(lockKey, lockRelationKey, meta, rl.preLcokCompare)
}

// Lock  lock resources
func (rl *RedisLock) Lock(meta *types.Lock) (*types.LockResult, error) {
	if nil == meta {
		return nil, fmt.Errorf("params not allowed null")
	}
	lockKey := getFmtRedisKey(rl.prefix, meta.LockName, false, false)
	lockRelationKey := getFmtRedisKey(rl.prefix, meta.TxnID, false, true)

	locked, subTxnID, err := rl.lock(lockKey, lockRelationKey, meta, rl.lockCompare)
	if nil != err {
		return nil, err
	}

	result := new(types.LockResult)
	result.Locked = locked
	result.SubTxnID = meta.SubTxnID
	result.LockSubTxnID = subTxnID
	return result, nil
}

// Unlock  unlock resources
func (rl *RedisLock) Unlock(meta *types.Lock) error {
	if nil == meta {
		return fmt.Errorf("params not allowed null")
	}
	lockKey := getFmtRedisKey(rl.prefix, meta.LockName, false, false)
	lockRelationKey := getFmtRedisKey(rl.prefix, meta.TxnID, false, true)
	return rl.unlock(lockKey, lockRelationKey, meta, rl.lockCompare)
}

// Unlock  unlock resources
func (rl *RedisLock) UnlockAll(id string) error {
	lockRelationKey := getFmtRedisKey(rl.prefix, id, true, true)
	err := rl.unlockall(lockRelationKey, rl.PreUnlock)
	if nil != err {
		return err
	}
	lockRelationKey = getFmtRedisKey(rl.prefix, id, false, true)
	err = rl.unlockall(lockRelationKey, rl.Unlock)
	if nil != err {
		return err
	}

	return nil
}

func (rl *RedisLock) unlockall(lockRelationKey string, unlock func(meta *types.Lock) error) error {
	keys, err := rl.storage.HGetAll(lockRelationKey).Result()
	if nil != err {
		if redis.Nil != err {
			return err
		}
	}
	for _, val := range keys {
		meta := new(types.Lock)
		err := json.Unmarshal([]byte(val), meta)
		if nil != err {
			return err
		}
		err = unlock(meta)
		if nil != err {
			if types.LockNotFound != err && types.LockPermissionDenied != err {
				return err
			}
		}
	}
	err = rl.storage.Del(lockRelationKey).Err()
	if nil != err {
		blog.Errorf("unlockall redis delete key %s error, error:%s", lockRelationKey, err.Error())
	}
	return nil
}
