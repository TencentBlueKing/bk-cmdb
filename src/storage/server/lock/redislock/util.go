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

const (
	lockPreRelationFmtStr = "%s:lock:pre:id:%s"
	lockPreFmtStr         = "%s:lock:pre:name:%s"

	lockRelationFmtStr = "%s:lock:detail:id:%s"
	lockFmtStr         = "%s:lock:detail:name:%s"

	noticeMaxCount = 500

	redisScanKeyCount = 256
)

type CompareFunc func(val, redisVal *types.Lock) (bool, error)

type noticType int

const (
	_ noticType = iota
	noticTypeErrLockCollection
	noticTypeErrUnLockCollection

	noticeTypeUnlockSuccess
)

type notice struct {
	key       string
	tid       string
	lockName  string
	noticType noticType
}

func LockCompare(val, redisMeta *types.Lock) (bool, error) {
	if redisMeta.TxnID == val.TxnID {
		return true, nil
	}
	return false, nil
}

func getRedisRelationInfoBy(storage *redis.Client, key, field string) (*types.Lock, bool, error) {
	return getLockInfoByKey(storage.HGet(key, field).Result())
}

func getRedisLockInfoByKey(storage *redis.Client, key string) (*types.Lock, bool, error) {
	return getLockInfoByKey(storage.Get(key).Result())
}

func getLockInfoByKey(str string, err error) (*types.Lock, bool, error) {
	if nil == err {
		lockInfo := new(types.Lock)
		err := json.Unmarshal([]byte(str), lockInfo)
		if nil != err {
			err := fmt.Errorf("json unmarshal error, reply:%s, error:%s", str, err.Error())
			blog.Error(err.Error())
			return nil, false, err
		} else {
			return lockInfo, true, nil
		}
	} else if redis.Nil == err {
		return nil, false, nil
	}

	return nil, false, err
}

func getFmtRedisKey(prefix, name string, isPre, isRelation bool) string {
	if isPre {
		if isRelation {
			return fmt.Sprintf(lockPreRelationFmtStr, prefix, name)
		} else {
			return fmt.Sprintf(lockPreFmtStr, prefix, name)
		}
	} else {
		if isRelation {
			return fmt.Sprintf(lockRelationFmtStr, prefix, name)
		} else {
			return fmt.Sprintf(lockFmtStr, prefix, name)
		}
	}

	return ""
}
