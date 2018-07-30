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

func (rl *RedisLock) lock(lockKey, lockcollectionKey string, meta *types.Lock, compare CompareFunc) (locked bool, SubTxnID string, err error) {

	sleepTime := time.Millisecond * 100
	total := int(meta.Timeout/time.Second)*(int(time.Second)/int(sleepTime)) + 1
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
		time.Sleep(sleepTime)
	}

	return false, SubTxnID, nil
}

func (rl *RedisLock) unlock(lockKey, lockcollectionKey string, meta *types.Lock, compare CompareFunc) error {
	content := rl.storage.Get(lockKey)
	if nil != content.Err() {
		if redis.Nil == content.Err() {
			return types.LockNotFound
		}
		blog.Errorf("unlock redis get key %s error, error:%s", lockKey, content.Err().Error())
		return content.Err()
	}

	redisMeta := new(types.Lock)
	err := json.Unmarshal([]byte(content.Val()), redisMeta)
	if nil != err {
		blog.Errorf("unlock redis  key %s  content  json unmarshal error, content:%s error:%s", lockKey, content.Val(), content.Err().Error())
		return err
	}
	hasLocked, err := compare(meta, redisMeta)
	if nil != err {
		blog.Errorf("unlock compare lock error, error:%s, lock:%s, redis info:%v", err.Error(), meta, redisMeta)
		return err
	}
	if hasLocked {
		err := rl.storage.Del(lockKey).Err()
		if nil != err {
			blog.Errorf("unlock  delete redis key %s error, error:%s", lockKey, err.Error())
			return err
		}
		err = rl.storage.HDel(lockcollectionKey, meta.LockName).Err()
		if nil != err {
			rl.notice(lockcollectionKey, meta.TxnID, meta.LockName, noticTypeErrUnLockCollection)
		} else {
			rl.notice(lockcollectionKey, meta.TxnID, meta.LockName, noticeTypeUnlockSuccess)
		}
		return nil
	}

	return types.LockPermissionDenied
}

func (rl *RedisLock) setLock(lockKey, lockcollectionKey string, meta *types.Lock) (bool, string, error) {
	meta.Createtime = time.Now().UTC()
	strVal, err := json.Marshal(meta)
	if nil != err {
		blog.Errorf("setLock redis  key %s  content  json Marshal error, content:%s error:%s", lockKey, meta, err.Error())
		return false, "", err
	}
	locked, err := rl.storage.SetNX(lockKey, string(strVal), 0).Result()
	if nil != err {
		blog.Errorf("setLock redis  key %s  SetNX error, error:%s", lockKey, err.Error())
		return false, "", err

	}
	if false == locked {
		return false, "", nil
	}

	err = rl.storage.HSet(lockcollectionKey, meta.LockName, string(strVal)).Err()
	if nil != err {
		// compensation mechanism  is considered here, notice function clear lockkey
		rl.notice(lockKey, meta.TxnID, meta.LockName, noticTypeErrLockCollection)
		blog.Errorf("setLock redis HSet key %s fields %s error, error:%s", lockKey, meta.LockName, err.Error())
		return false, "", err
	}
	return true, meta.SubTxnID, nil
}

func (rl *RedisLock) compareLock(lockKey string, meta *types.Lock, compare CompareFunc) (bool, string, error) {
	content, err := rl.storage.Get(lockKey).Result()
	if nil != err {
		if redis.Nil == err {
			return false, "", nil
		}
		blog.Errorf("setLock redis  get key %s  content error,error:%s", lockKey, err.Error())
		return false, "", err
	}
	redisMeta := new(types.Lock)
	err = json.Unmarshal([]byte(content), redisMeta)
	if nil != err {
		blog.Errorf("setLock redis key %s  content json unmarshal error,content:%s error:%s", lockKey, content, err.Error())
		return false, "", err
	}
	bl, diffErr := compare(meta, redisMeta)
	if nil != diffErr {
		blog.Errorf("setLock compare key %s  content error, lock:%v, redis info:%v:%s error:%s", lockKey, content, err.Error())
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
	timer := time.NewTicker(time.Second * time.Duration(rl.triggerTime))

	for {
		select {
		case n := <-rl.noticeChan:
			switch n.noticType {
			case noticTypeErrLockCollection:
				//lock successfully, but record lock and sub id relationship failed
				go rl.noticTypeErrLockCollection(n)
			case noticTypeErrUnLockCollection:
				//unlock successfully, but delete lock and sub id relationship failed
				go rl.noticTypeErrUnLockCollection(n)
			case noticeTypeUnlockSuccess:
				// unlock sucessfully, clear lock and sub id relationship emtpy key
				go rl.noticeTypeUnlockSuccess(n)
			}
		case <-timer.C:
			go func() {
				err := rl.noticeTimedTrigger()
				if nil != err {
					blog.Errorf("compensation error:%s", err.Error())
				}
			}()
		}
	}

}

func (rl *RedisLock) noticeTimedTrigger() error {
	prefix := getFmtRedisKey(rl.prefix, "", true, false)
	relationPrefix := getFmtRedisKey(rl.prefix, "", true, true)
	rl.compensationLock(prefix, relationPrefix)
	rl.compensationRelation(relationPrefix, true)

	prefix = getFmtRedisKey(rl.prefix, "", false, false)
	relationPrefix = getFmtRedisKey(rl.prefix, "", false, true)
	rl.compensationLock(prefix, relationPrefix)
	rl.compensationRelation(relationPrefix, false)

	return nil
}

func (rl *RedisLock) compensationLock(prefix, relationPrefix string) {
	var keys []string
	var err error
	var cursor uint64

	keys, cursor, err = rl.storage.Scan(cursor, fmt.Sprintf("%s%s", prefix, "*"), redisScanKeyCount).Result()
	if nil != err && redis.Nil != err {
		blog.Errorf("compensationLock redis scan key %s error,error %s", prefix, err.Error())
	}

	for _, key := range keys {
		lockInfo, isExist, err := getRedisLockInfoByKey(rl.storage, key)
		if nil != err {
			blog.Errorf("compensationLock key %s, error:%s", err.Error())
			continue
		}
		if false == isExist {
			continue
		}
		relKey := fmt.Sprintf("%s%s", relationPrefix, lockInfo.TxnID)
		_, isExist, err = getRedisRelationInfoBy(rl.storage, relKey, lockInfo.LockName)
		if nil != err {
			blog.Errorf("compensationLock key %s field %s, error: %s", relKey, lockInfo.LockName, err.Error())
			continue
		}
		if false == isExist {
			err := rl.storage.Del(key).Err()
			if nil != err {
				blog.Errorf("compensationLock delete redis key %s error, error:%s", key, err.Error())
			}
		}

	}
}

func (rl *RedisLock) compensationRelation(prefix string, isPre bool) {
	var keys []string
	var err error
	var cursor uint64

	keys, cursor, err = rl.storage.Scan(cursor, fmt.Sprintf("%s%s", prefix, "*"), redisScanKeyCount).Result()
	if nil != err && redis.Nil != err {
		blog.Errorf("compensationRelation redis scan error %s", err.Error())
		return
	}
	for _, key := range keys {
		ret := rl.storage.HLen(key)
		if nil == ret.Err() || redis.Nil == ret.Err() {
			if 0 == ret.Val() {
				err = rl.storage.Del(key).Err()
				if nil != err {
					blog.Errorf("compensationRelation delete redis key %s error %s", key, err.Error())
				}
			} else {
				err := rl.compensationRelationHashFields(key, isPre)
				if nil != err {
					blog.Errorf("compensationRelation  key %s error %s", key, err.Error())
				}
			}
		}
	}
}

func (rl *RedisLock) noticeTypeUnlockSuccess(n *notice) {
	for i := 0; i < rl.errRetryCnt; i++ {
		ret := rl.storage.HLen(n.key)
		if redis.Nil == ret.Err() {
			break
		} else if nil == ret.Err() {
			if 0 == ret.Val() {
				err := rl.storage.Del(n.key).Err()
				if nil != err {
					blog.Errorf("redis delete key %s error, error:%s", n.key, err.Error())
					continue
				}
			}
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (rl *RedisLock) noticTypeErrUnLockCollection(n *notice) {
	for i := 0; i < rl.errRetryCnt; i++ {
		err := rl.storage.HDel(n.key, n.lockName).Err()
		if nil == err || redis.Nil == err {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (rl *RedisLock) noticTypeErrLockCollection(n *notice) {
	for i := 0; i < rl.errRetryCnt; i++ {
		val, err := rl.storage.Get(n.key).Result()
		if redis.Nil == err {
			break
		} else if nil == err {
			meta := new(types.Lock)
			err = json.Unmarshal([]byte(val), meta)
			if nil != err {
				blog.Errorf("redis key %s content %s not json", n.key, val)
				break
			}
			if meta.TxnID != n.tid {
				break
			}
			err = rl.storage.Del(n.key).Err()
			if nil != err {
				blog.Errorf("redis delete key %s error, error:%s", n.key, err.Error())
			}
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func (rl *RedisLock) compensationRelationHashFields(key string, isPre bool) error {
	var fields []string
	var err error
	var cursor uint64

	fields, cursor, err = rl.storage.HScan(key, cursor, "*", redisScanKeyCount).Result() //Scan(cursor, prefix, redisScanKeyCount).Result()
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
		lockRealtionKey := getFmtRedisKey(rl.prefix, "", isPre, true)

		lockInfo, isExist, err := getRedisLockInfoByKey(rl.storage, lockRealtionKey)
		if nil != err {
			blog.Errorf("compensationLock %s", err.Error())
			continue
		}
		if false == isExist {
			err = rl.storage.HDel(key, field).Err()
		} else if lockInfo.TxnID != relLockInfo.TxnID || lockInfo.SubTxnID != relLockInfo.SubTxnID {
			err = rl.storage.HDel(key, field).Err()
		}
		if nil != err {
			blog.Errorf("compensationLock  delete key %s error, error:%s", key, err.Error())
			continue
		}

	}

	return nil
}
