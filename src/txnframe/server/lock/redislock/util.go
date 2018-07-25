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
	"configcenter/src/txnframe/server/lock/types"
)

const (
	lockPreCollectionFmtStr = "%s:lock:table:id:%s"
	lockPreFmtStr           = "%s:lock:table:%s"

	lockCollectionFmtStr = "%s:lock:id:%s"
	lockFmtStr           = "%s:lock:%s"

	noticeMaxCount = 500
)

type compareFunc func(val, redisVal *types.Lock) (bool, error)

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

func lockCompare(val, redisMeta *types.Lock) (bool, error) {
	if redisMeta.TxnID == val.TxnID {
		return true, nil
	}
	return false, nil
}
