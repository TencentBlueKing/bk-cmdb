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

	redis "github.com/go-redis/redis"

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
		content := rl.storage.Get(lockKey)
		if nil == content.Err() {
			locked, SubTxnID, err = rl.compareLock(content.Val(), meta, compare)
		} else if redis.Nil == content.Err() {
			locked, SubTxnID, err = rl.setLock(lockKey, lockcollectionKey, meta)
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
			rl.notice(lockcollectionKey, string(meta.TxnID), meta.LockName, noticTypeErrUnLockCollection)
		}
		rl.notice(lockcollectionKey, meta.TxnID, meta.LockName, noticeTypeUnlockSuccess)
		return nil
	}

	return types.LockPermissionDenied
}

func (rl *RedisLock) setLock(lockKey, lockcollectionKey string, meta *types.Lock) (bool, string, error) {
	meta.Createtime = time.Now().UTC()
	err := rl.storage.SetNX(lockKey, meta, 0).Err()
	if nil == err {
		err := rl.storage.HSet(lockcollectionKey, meta.LockName, meta).Err()
		if nil != err {
			// compensation mechanism  is considered here, notice function clear lockkey
			rl.notice(lockKey, string(meta.TxnID), meta.LockName, noticTypeErrLockCollection)
			return false, "", err
		}
		return true, meta.SubTxnID, nil
	} else {
		// wait next execute, try to lock
		//return false, err
	}
	return false, "", nil
}

func (rl *RedisLock) compareLock(content string, meta *types.Lock, compare compareFunc) (bool, string, error) {
	redisMeta := new(types.Lock)
	err := json.Unmarshal([]byte(content), redisMeta)
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
		rl.noticeTimedTrigger()
	}
}

func (rl *RedisLock) noticeTimedTrigger() error {
	handle := func(prefix string, isHash bool) error {
		var keys []string
		var err error
		var cursor uint64

		keys, cursor, err = rl.storage.Scan(cursor, prefix, redisScanKeyCount).Result()
		if nil != err && redis.Nil != err {
			return err
		}
		for _, key := range keys {
			ret := rl.storage.HLen(key)
			if nil == ret.Err() || redis.Nil == ret.Err() {
				if 0 == ret.Val() {
					rl.storage.Del(key)
				}
			}
		}
		return nil
	}

	prefix := fmt.Sprintf(lockCollectionFmtStr, rl.prefix, "")
	return handle(prefix, true)

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
		ret := rl.storage.Get(n.key)
		if nil == ret {
			meta := new(types.Lock)
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
		}
		if redis.Nil == ret.Err() {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}
