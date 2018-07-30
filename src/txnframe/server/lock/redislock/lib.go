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
	"time"

	redis "gopkg.in/redis.v5"

	"configcenter/src/common/blog"
	"configcenter/src/txnframe/server/lock/types"
)

func (rl *RedisLock) tryLockSubTxnID(meta *types.Lock) {
	if "" == meta.SubTxnID {
		meta.SubTxnID = types.GetID(types.LockIDPrefix)
	}
}

func (rl *RedisLock) lock(lockKey, lockcollectionKey string, meta *types.Lock, compare compareFunc) (locked bool, SubTxnID string, err error) {
	total := int(meta.Timeout / time.Second)
	rl.tryLockSubTxnID(meta)
	for idx := 0; idx < total; idx++ {
		locked, SubTxnID, err = rl.setLock(lockKey, lockcollectionKey, meta)
		if nil == err && false == locked {
			locked, SubTxnID, err = rl.compareLock(lockKey, meta, compare)
		}

		// has error or get locked return
		if nil != err || true == locked {
			return locked, SubTxnID, err
		}
		time.Sleep(time.Millisecond * 100)
	}

	return false, SubTxnID, nil
}

func (rl *RedisLock) unlock(lockKey, lockcollectionKey string, meta *types.Lock, compare compareFunc) error {
	content := rl.storage.Get(lockKey)
	if nil != content.Err() {
		if redis.Nil == content.Err() {
			return types.LockNotFound
		}
		return content.Err()
	}

	redisMeta := new(types.Lock)
	err := json.Unmarshal([]byte(content.Val()), redisMeta)
	if nil != err {
		return err
	}
	hasLocked, err := compare(meta, redisMeta)
	if nil != err {
		return err
	}
	if hasLocked {
		err := rl.storage.Del(lockKey).Err()
		if nil != err {
			return err
		}
		err = rl.storage.HDel(lockcollectionKey, meta.LockName).Err()
		if nil != err {
			rl.notice(lockcollectionKey, meta.TxnID, meta.LockName, noticTypeErrUnLockCollection)
		}
		rl.notice(lockcollectionKey, meta.TxnID, meta.LockName, noticeTypeUnlockSuccess)
		return nil
	}

	return types.LockPermissionDenied
}

func (rl *RedisLock) setLock(lockKey, lockcollectionKey string, meta *types.Lock) (bool, string, error) {
	meta.Createtime = time.Now().UTC()
	strVal, err := json.Marshal(meta)
	if nil != err {
		return false, "", err
	}
	locked, err := rl.storage.SetNX(lockKey, string(strVal), 0).Result()
	if nil != err {
		return false, "", err

	}
	if false == locked {
		return false, "", nil
	}

	err = rl.storage.HSet(lockcollectionKey, meta.LockName, string(strVal)).Err()
	if nil != err {
		// compensation mechanism  is considered here, notice function clear lockkey
		rl.notice(lockKey, string(meta.TxnID), meta.LockName, noticTypeErrLockCollection)
		return false, "", err
	}
	return true, meta.SubTxnID, nil
}

func (rl *RedisLock) compareLock(lockKey string, meta *types.Lock, compare compareFunc) (bool, string, error) {
	content, err := rl.storage.Get(lockKey).Result()
	if nil != err {
		if redis.Nil == err {
			return false, "", nil
		}
		return false, "", err
	}
	redisMeta := new(types.Lock)
	err = json.Unmarshal([]byte(content), redisMeta)
	if nil != err {
		return false, "", err
	}
	bl, diffErr := compare(meta, redisMeta)
	if nil != diffErr {
		return false, "", err
	}

	return bl, redisMeta.SubTxnID, nil
}

func (rl *RedisLock) notice(key, tid, lockName string, t noticType) {
	if len(rl.noticeChan) < noticeMaxCount {
		go func() {
			rl.noticeChan <- &notice{key: key, tid: tid, lockName: lockName, noticType: t}
		}()
	}
}

func (rl *RedisLock) compensation() {
	timer := time.NewTicker(time.Minute * 5)

	select {
	case n := <-rl.noticeChan:
		switch n.noticType {
		case noticTypeErrLockCollection:
			//lock successfully, but record lock and sub id relationship failed
			rl.noticTypeErrLockCollection(n)
		case noticTypeErrUnLockCollection:
			//unlock successfully, but delete lock and sub id relationship failed
			rl.noticTypeErrUnLockCollection(n)
		case noticeTypeUnlockSuccess:
			// unlock sucessfully, clear lock and sub id relationship emtpy key
			rl.noticeTypeUnlockSuccess(n)
		}
	case <-timer.C:
		err := rl.noticeTimedTrigger()
		if nil != err {
			blog.Errorf("compensation error:%s", err.Error())
		}
	}
}

func (rl *RedisLock) noticeTimedTrigger() error {
	prefix := fmt.Sprintf(lockPreFmtStr, rl.prefix, "*")
	relationPrefix := fmt.Sprintf(lockPreCollectionFmtStr, rl.prefix, "*")
	rl.compensationLock(prefix, relationPrefix)
	rl.compensationRelation(relationPrefix, true)

	prefix = fmt.Sprintf(lockFmtStr, rl.prefix, "*")
	relationPrefix = fmt.Sprintf(lockCollectionFmtStr, rl.prefix, "*")
	rl.compensationLock(prefix, relationPrefix)
	rl.compensationRelation(relationPrefix, false)

	return nil
}

func (rl *RedisLock) compensationLock(prefix, relationPrefix string) {
	var keys []string
	var err error
	var cursor uint64

	keys, cursor, err = rl.storage.Scan(cursor, prefix, redisScanKeyCount).Result()
	if nil != err && redis.Nil != err {
		blog.Errorf("compensationLock redis scan error %s", err.Error())
	}

	for _, key := range keys {
		lockInfo, isExist, err := getRedisLockInfoByKey(rl.storage, key)
		if nil != err {
			blog.Errorf("compensationLock %s", err.Error())
			continue
		}
		if false == isExist {
			continue
		}
		_, isExist, err = getRedisLockInfoByKey(rl.storage, fmt.Sprintf("%s%s", relationPrefix, lockInfo.LockName))
		if nil != err {
			blog.Errorf("compensationLock %s", err.Error())
			continue
		}
		if false == isExist {
			rl.storage.Del(key)
		}

	}
}

func (rl *RedisLock) compensationRelation(prefix string, isPre bool) {
	var keys []string
	var err error
	var cursor uint64

	keys, cursor, err = rl.storage.Scan(cursor, prefix, redisScanKeyCount).Result()
	if nil != err && redis.Nil != err {
		blog.Errorf("compensationRelation redis scan error %s", err.Error())
	}
	for _, key := range keys {
		ret := rl.storage.HLen(key)
		if nil == ret.Err() || redis.Nil == ret.Err() {
			if 0 == ret.Val() {
				rl.storage.Del(key)
			} else {
				err := rl.compensationRelationHashFields(key, isPre)
				if nil != err {
					blog.Errorf("compensationRelationHashFields  error %s", err.Error())
				}
			}
		}
	}
}

func (rl *RedisLock) noticeTypeUnlockSuccess(n *notice) {
	for {
		ret := rl.storage.HLen(n.key)
		if nil == ret.Err() || redis.Nil == ret.Err() {
			if 0 == ret.Val() {
				rl.storage.Del(n.key)
			}
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (rl *RedisLock) noticTypeErrUnLockCollection(n *notice) {
	for {
		err := rl.storage.HDel(n.key, n.lockName).Err()
		if nil == err || redis.Nil == err {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (rl *RedisLock) noticTypeErrLockCollection(n *notice) {
	for {
		val, err := rl.storage.Get(n.key).Result()
		if nil == err {
			meta := new(types.Lock)
			err := json.Unmarshal([]byte(val), meta)
			if nil != err {
				blog.Errorf("redis key %s content %s not json", n.key, val)
				break
			}
			if meta.TxnID == n.tid {
				err := rl.storage.Del(n.key).Err()
				if nil != err && redis.Nil != err {
					continue
				} else {
					break
				}
			} else {
				break
			}
		} else if redis.Nil == err {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (rl *RedisLock) compensationRelationHashFields(key string, isPre bool) error {
	var fields []string
	var err error
	var cursor uint64

	fields, cursor, err = rl.storage.HScan(key, cursor, "", redisScanKeyCount).Result() //Scan(cursor, prefix, redisScanKeyCount).Result()
	if nil != err && redis.Nil != err {
		return fmt.Errorf("compensationRelation redis scan error %s", err.Error())

	}
	for _, field := range fields {
		relLockInfo, isExist, err := getRedisRelationInfoBy(rl.storage, key, field)
		if nil != err {
			blog.Errorf("compensationLock %s", err.Error())
			continue
		}
		if false == isExist {
			continue
		}
		lockKey := ""
		if isPre {
			lockKey = fmt.Sprintf("%s%s", lockPreCollectionFmtStr, relLockInfo.LockName)
		} else {
			lockKey = fmt.Sprintf("%s%s", lockCollectionFmtStr, relLockInfo.LockName)

		}
		lockInfo, isExist, err := getRedisLockInfoByKey(rl.storage, lockKey)
		if nil != err {
			blog.Errorf("compensationLock %s", err.Error())
			continue
		}
		if false == isExist {
			rl.storage.HDel(key, field)
		} else if lockInfo.TxnID != relLockInfo.TxnID || lockInfo.SubTxnID != relLockInfo.SubTxnID {
			rl.storage.HDel(key, field)
		}

	}

	return nil
}
