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

package distribution

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	ccredis "configcenter/src/storage/dal/redis"
	daltypes "configcenter/src/storage/types"
)

func (th *TxnHandler) Run() (err error) {
	defer func() {
		syserror := recover()
		if syserror != nil {
			err = fmt.Errorf("system error: %v", syserror)
		}
		if err != nil {
			blog.Info("event inst handle process stoped by %v", err)
			blog.Errorf("%s", debug.Stack())
		}
	}()

	blog.Info("transaction handle process started")
	go th.fetchTimeout()
outer:
	for txnID := range th.commited {
		for {
			err = th.cache.RPopLPush(common.EventCacheEventTxnQueuePrefix+txnID, common.EventCacheEventQueueKey).Err()
			if ccredis.IsNilErr(err) {
				break
			}
			if err != nil {
				blog.Warnf("move commited event to event queue failed: %v, we will try again later", err)
				continue outer
			}
		}
		if err = th.cache.Del(common.EventCacheEventTxnQueuePrefix + txnID).Err(); err != nil {
			blog.Warnf("remove [%s] transaction queue failed: %v, we will try again later", txnID, err)
			continue
		}
		if err = th.cache.ZRem(common.EventCacheEventTxnSet, txnID).Err(); err != nil {
			blog.Warnf("remove [%s] from transaction set failed: %v, we will try again later", txnID, err)
		}
	}
	return nil
}

func (th *TxnHandler) fetchTimeout() {
	tick := time.Tick(time.Second * (60))
	opt := redis.ZRangeBy{
		Min: "-inf",
	}
	for now := range tick {
		txnIDs := []string{}
		opt.Max = strconv.Itoa(now.UTC().Second())

		if err := th.cache.ZRangeByScore(common.EventCacheEventTxnSet, opt).ScanSlice(&txnIDs); err != nil {
			blog.Warnf("fetch timeout txnID from redis failed: %v, we will try again later", err)
			continue
		}

		txns := []daltypes.Transaction{} //Transaction
		if err := th.db.Table(common.BKTableNameTransaction).Find(dal.NewFilterBuilder().In(common.BKTxnIDField, txnIDs)).All(th.ctx, &txns); err != nil {
			blog.Warnf("fetch transaction from mongo failed: %v, we will try again later", err)
			continue
		}
		droped := []string{}
		for _, txn := range txns {
			switch txn.Status {
			case daltypes.TxStatusOnProgress:
				continue
			case daltypes.TxStatusCommited:
				th.commited <- txn.TxnID
			case daltypes.TxStatusAborted, daltypes.TxStatusException:
				droped = append(droped, txn.TxnID)
			default:
				blog.Warnf("unknow transaction status %s", txn.Status)
			}
		}
		go th.dropTransaction(droped)
	}
}

func (th *TxnHandler) dropTransaction(txnIDs []string) {
	dropKeys := make([]string, len(txnIDs))
	dropTxnIDs := make([]interface{}, len(txnIDs))
	for index, txnID := range txnIDs {
		dropTxnIDs[index] = txnID
		dropKeys[index] = common.EventCacheEventTxnQueuePrefix + txnID
	}
	if err := th.cache.Del(dropKeys...).Err(); err != nil {
		blog.Warnf("drop transaction queue [%v] failed: %v, we will try again later", dropKeys, err)
		return
	}
	if err := th.cache.ZRem(common.EventCacheEventTxnSet, dropTxnIDs...).Err(); err != nil {
		blog.Warnf("drop [%v] from transaction set failed: %v, we will try again later", dropTxnIDs, err)
	}
}
