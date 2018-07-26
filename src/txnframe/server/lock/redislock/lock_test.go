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
