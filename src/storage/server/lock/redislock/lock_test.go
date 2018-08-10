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
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	redis "gopkg.in/redis.v5"

	"configcenter/src/storage/server/lock/types"
)

func TestPreLock(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)
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

func TestPrivateCompensationNoticeLockErr(t *testing.T) {
	lock1 := &types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}

	lock2 := &types.Lock{
		TxnID:    "2",
		LockName: "test",
		Timeout:  time.Second,
	}

	s := getRedisInstance()
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)

	// lock alread exists, but the relationship does not exist
	setKey := fmt.Sprintf(lockPreFmtStr, ss.prefix, lock1.LockName)
	str, err := json.Marshal(lock1)
	if nil != err {
		t.Errorf("json marshal error, error:%s", err.Error())
		return
	}
	s.Set(setKey, string(str), 0)
	locked, err := ss.PreLock(lock2)
	if nil != err {
		t.Errorf("lock test lock1 error:%s", err.Error())
		return
	}
	if true == locked {
		t.Errorf("lock must be false, not true")
		return
	}
	ss.notice(setKey, lock1.TxnID, lock1.LockName, noticTypeErrLockCollection)
	time.Sleep(time.Second * 2)
	ok, err := s.Exists(setKey).Result()
	if nil != err {
		t.Error(err.Error())
		return
	}
	if ok {
		t.Error("notice error")
		return
	}
}

func TestPrivateCompensationNoticeSuccess(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)

	lock := &types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}

	// test compensation delete redis emtpy hash key
	setKey := getFmtRedisKey(ss.prefix, lock.TxnID, true, true) //fmt.Sprintf(lockPreCollectionFmtStr, ss.prefix, lock.TxnID)
	locked, err := ss.PreLock(lock)
	if nil != err {
		t.Error(err.Error())
		return
	}
	if false == locked {
		t.Error("lock resource error")
		return
	}
	err = ss.PreUnlock(lock)
	if nil != err {
		t.Errorf("unlock error, error%s", err.Error())
		return
	}
	err = testLockRelationKey(s, setKey, lock)
	if nil != err {
		t.Errorf("testLockRelationKey, error%s", err.Error())
		return
	}
}

func TestPrivateCompensationRelationNotDel(t *testing.T) {

	s := getRedisInstance()
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)

	lock := &types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
	}
	// test compensation,  lock key is delete bu relation key not delete
	setKey := getFmtRedisKey(ss.prefix, lock.TxnID, true, true) //fmt.Sprintf(lockPreCollectionFmtStr, ss.prefix, lock.TxnID)
	ok, err := s.HSet(setKey, lock.LockName, "{}").Result()
	if nil != err {
		t.Errorf("set hash key %s field %s error, error:%s", setKey, lock.LockName, err.Error())
		return
	}
	if false == ok {
		t.Errorf("set hash key %s field %s error, ", setKey, lock.LockName)
		return
	}
	err = s.HDel(setKey, lock.LockName).Err()
	if nil != err {
		t.Errorf("delete hash key %s field %s error, ", setKey, lock.LockName)
		return
	}
	ss.notice(setKey, lock.TxnID, lock.LockName, noticTypeErrUnLockCollection)
	err = testLockRelationKey(s, setKey, lock)
	if nil != err {
		t.Errorf("testLockRelationKey, error%s", err.Error())
		return
	}

}

func TestPrivateCompensationTimeTrigger(t *testing.T) {
	s := getRedisInstance()
	ss := NewLock(s, "cc", 10, 120, LockCompare, LockCompare)

	lock := &types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
		SubTxnID: "1",
	}

	lock1 := &types.Lock{
		TxnID:    "1",
		LockName: "test",
		Timeout:  time.Second,
		SubTxnID: "2",
	}

	str, err := json.Marshal(lock)
	if nil != err {
		t.Errorf("json marshal error, error:%s", err.Error())
		return
	}

	str1, err := json.Marshal(lock1)
	if nil != err {
		t.Errorf("json marshal error, error:%s", err.Error())
		return
	}

	setKey := fmt.Sprintf(lockFmtStr, ss.prefix, lock.LockName)

	// test lock error, lock key create success, but relation create failure
	err = s.Set(setKey, str, 0).Err()
	if nil != err {
		t.Errorf("set key %s  value %s error, error:%s", setKey, str, err.Error())
		return
	}
	err = ss.noticeTimedTrigger()
	if nil != err {
		t.Errorf("noticeTimedTrigger error, error:%s", err.Error())
		return
	}
	ok, err := s.Exists(setKey).Result()
	if nil != err {
		t.Errorf("exist key %s  value %s error, error:%s", setKey, str, err.Error())
		return
	}
	if ok {
		t.Errorf("key %s exist, must be not exist", setKey)
		return

	}

	// test lock relation error, but relation info and lock key  SubTxnID not equal
	err = s.Set(setKey, str, 0).Err()
	if nil != err {
		t.Errorf("set key %s  value %s error, error:%s", setKey, str, err.Error())
		return
	}
	setRelKey := getFmtRedisKey(ss.prefix, lock.TxnID, false, true) //fmt.Sprintf(lockCollectionFmtStr, ss.prefix, lock.TxnID)
	err = s.HSet(setRelKey, lock.LockName, str1).Err()
	if nil != err {
		t.Errorf("set key %s  value %s error, error:%s", setRelKey, str1, err.Error())
		return
	}

	err = ss.noticeTimedTrigger()
	if nil != err {
		t.Errorf("noticeTimedTrigger error, error:%s", err.Error())
		return
	}
	ok, err = s.HExists(setRelKey, lock.LockName).Result()
	if nil != err {
		t.Errorf("exist key %s  field %s error, error:%s", setRelKey, lock.LockName, err.Error())
		return
	}
	if ok {
		t.Errorf("key %s field %sexist, must be not exist", setRelKey, lock.LockName)
		return

	}

	lock1.TxnID = "2"
	str1, err = json.Marshal(lock1)
	if nil != err {
		t.Errorf("json marshal error, error:%s", err.Error())
		return
	}
	// test lock relation error, but relation info and lock key  TxnID not equal
	err = s.Set(setKey, str, 0).Err()
	if nil != err {
		t.Errorf("set key %s  value %s error, error:%s", setKey, str, err.Error())
		return
	}
	setRelKey = getFmtRedisKey(ss.prefix, lock.TxnID, false, true) //fmt.Sprintf(lockCollectionFmtStr, ss.prefix, lock.TxnID)
	err = s.HSet(setRelKey, lock.LockName, str1).Err()
	if nil != err {
		t.Errorf("set key %s  value %s error, error:%s", setRelKey, str1, err.Error())
		return
	}

	err = ss.noticeTimedTrigger()
	if nil != err {
		t.Errorf("noticeTimedTrigger error, error:%s", err.Error())
		return
	}
	ok, err = s.HExists(setRelKey, lock.LockName).Result()
	if nil != err {
		t.Errorf("exist key %s  field %s error, error:%s", setRelKey, lock.LockName, err.Error())
		return
	}
	if ok {
		t.Errorf("key %s field %sexist, must be not exist", setRelKey, lock.LockName)
		return

	}

}

func testLockRelationKey(s *redis.Client, setKey string, lock *types.Lock) error {
	len, err := s.HLen(setKey).Result()
	if nil != err {
		return fmt.Errorf("hlen key %s error, error:%s", setKey, err.Error())

	}
	if 0 != len {
		return fmt.Errorf("hlen key %s must be 0, not %d", setKey, len)
	}
	time.Sleep(2)
	ok, err := s.Exists(setKey).Result()
	if nil != err {
		return fmt.Errorf("hlen key %s error, error:%s", setKey, err.Error())
	}
	if ok == true {
		return fmt.Errorf("key %s exist, must be not exist", setKey)
	}
	return nil
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
