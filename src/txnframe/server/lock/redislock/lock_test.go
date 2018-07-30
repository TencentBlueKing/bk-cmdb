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
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	redis "gopkg.in/redis.v5"

	"configcenter/src/txnframe/server/lock/types"
)

func TestPreLock(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	locked, err := ss.PreLock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == locked {
		t.Errorf("lock is false")
	}
}

func TestPreLockGetMulti(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	locked, err := ss.PreLock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == locked {
		t.Errorf("lock is false")
	}

	locked, err = ss.PreLock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == locked {
		t.Errorf("repeat lock error, expect is just true not false")
	}
}

func TestPreLockErr(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	locked, err := ss.PreLock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == locked {
		t.Errorf("lock is false")
	}

	lock = types.Lock{
		TxnID:    "2",
		LockName: "test",
		Timeout:  time.Second * 2,
	}
	locked, err = ss.PreLock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if true == locked {
		t.Errorf("repeat lock error, expect is just false not true")
	}
}

func TestPreUnlock(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	locked, err := ss.PreLock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == locked {
		t.Errorf("lock is false")
		return
	}
	err = ss.PreUnlock(&lock)
	if nil != err {
		t.Errorf("preunlock  error %s", err.Error())
		return
	}

	locked, err = ss.PreLock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == locked {
		t.Errorf("lock is false")
		return
	}
}

func TestPreUnlockErr(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	locked, err := ss.PreLock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == locked {
		t.Errorf("lock is false")
		return
	}
	lock = types.Lock{
		TxnID:    "2",
		LockName: "test",
		Timeout:  time.Second,
	}
	err = ss.PreUnlock(&lock)
	if nil == err {
		t.Errorf("preunlock  error, should be not permission")
		return
	}
}

func TestLock(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	lockInfo, err := ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lockInfo.Locked {
		t.Errorf("lock is false")
		return
	}
	if lockInfo.SubTxnID != lockInfo.LockSubTxnID {
		t.Errorf("lock id and sub id not equal")
		return
	}
}

func TestLockGetMulti(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	lcokInfo, err := ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lcokInfo.Locked {
		t.Errorf("lock is false")
	}

	lcokInfo, err = ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lcokInfo.Locked {
		t.Errorf("repeat lock error, expect is just true not false")
		return
	}
}

func TestLockGetMultiSameMasterID(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	lockInfo, err := ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lockInfo.Locked {
		t.Errorf("lock is false")
	}

	lock.SubTxnID = ""
	lockInfo, err = ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lockInfo.Locked {
		t.Errorf("repeat lock error, expect is just true not false")
		return
	}

	if lockInfo.SubTxnID == lockInfo.LockSubTxnID {
		t.Errorf("lock id and sub id expect  equal")
		return
	}
}

func TestLockErr(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	lockInfo, err := ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lockInfo.Locked {
		t.Errorf("lock is false")
	}

	lock = types.Lock{
		TxnID:    "2",
		LockName: "test",
		Timeout:  time.Second * 2,
	}
	lockInfo, err = ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if true == lockInfo.Locked {
		t.Errorf("repeat lock error, expect is just false not true")
	}

	if lockInfo.SubTxnID == lockInfo.LockSubTxnID {
		t.Errorf("lock id and sub id expect not equal")
		return
	}
}

func TestUnlock(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	lockInfo, err := ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lockInfo.Locked {
		t.Errorf("lock is false")
		return
	}
	err = ss.Unlock(&lock)
	if nil != err {
		t.Errorf("unlock  error %s", err.Error())
		return
	}

	lockInfo, err = ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lockInfo.Locked {
		t.Errorf("lock is false")
		return
	}
}

func TestUnlockErr(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	lockInfo, err := ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lockInfo.Locked {
		t.Errorf("lock is false")
		return
	}
	lock = types.Lock{
		TxnID:    "2",
		LockName: "test",
		Timeout:  time.Second,
	}
	err = ss.Unlock(&lock)
	if nil == err {
		t.Errorf("unlock  error, should be not permission")
		return
	}

	lock = types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	err = ss.Unlock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lockInfo.Locked {
		t.Errorf("lock is false")
		return
	}

	lockInfo, err = ss.Lock(&lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return

	}
	if false == lockInfo.Locked {
		t.Errorf("lock is false")
		return
	}

}

func TestUnLockAll(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc")
	lock := &types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}

	lock.SubTxnID = ""
	lockInfo, err := ss.Lock(lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return
	}
	if false == lockInfo.Locked {
		t.Errorf("lock is false")
		return
	}
	if nil != err {
		t.Errorf(err.Error())
		return
	}

	lock.SubTxnID = ""
	locked, err := ss.PreLock(lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return
	}
	if false == locked {
		t.Errorf("pre lock is false")
		return
	}
	if nil != err {
		t.Errorf(err.Error())
		return
	}

	err = ss.UnlockAll(lock.TxnID)
	if nil != err {
		t.Errorf(err.Error())
		return
	}

	lock.SubTxnID = ""
	lockInfo, err = ss.Lock(lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return
	}
	if false == lockInfo.Locked {
		t.Errorf("lock is false")
		return
	}
	if nil != err {
		t.Errorf(err.Error())
		return
	}
	lock.SubTxnID = ""
	locked, err = ss.PreLock(lock)
	if nil != err {
		t.Errorf("%s", err.Error())
		return
	}
	if false == locked {
		t.Errorf("pre lock is false")
		return
	}
	if nil != err {
		t.Errorf(err.Error())
		return
	}

}

func TestPrivateCompensationNotice(t *testing.T) {
	lock1 := &types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}

	lock2 := &types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}

	s := getRedisInstance()
	ss := NewLock(s, "cc")

	// hte analog lock alread exists, but the relationship does not exist
	setKey := fmt.Sprintf(lockPreFmtStr, ss.prefix, lock1.LockName)
	s.Set(setKey, "{}", 0)
	locked, err := ss.PreLock(lock1)
	if nil != err {
		t.Errorf("lock test lock1 error:%s", err.Error())
		return
	}
	if true == locked {
		t.Errorf("lock must be false, not true")
		return
	}

}

func getRedisInstance() *redis.Client {
	mockRedis, err := miniredis.Run()
	/*if nil != err {
		panic(err )
	}*/
	storage := redis.NewClient(
		&redis.Options{
			Addr: mockRedis.Addr(),
		})

	err = storage.Ping().Err()
	if err != nil {
		panic(err)
	}

	return storage
}
