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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	ccredis "configcenter/src/storage/dal/redis"
	daltypes "configcenter/src/storage/types"

	"gopkg.in/redis.v5"
)

func (th *TxnHandler) Run() (err error) {
	th.shouldClose.UnSet()
	defer func() {
		th.shouldClose.Set()
		syserror := recover()
		if syserror != nil {
			err = fmt.Errorf("system error: %v", syserror)
		}
		if err != nil {
			blog.Infof("event inst handle process stoped by %v", err)
			blog.Errorf("%s", debug.Stack())
		}
		th.wg.Wait()
	}()

	blog.Info("transaction handle process started")
	th.wg.Add(1)
	go th.fetchTimeout()
	th.wg.Add(1)
	go th.watchTransaction()
outer:
	for txnID := range th.committed {
		blog.V(4).Infof("transaction %v committed", txnID)
		for {
			err = th.cache.RPopLPush(common.EventCacheEventTxnQueuePrefix+txnID, common.EventCacheEventQueueKey).Err()
			if ccredis.IsNilErr(err) {
				break
			}
			if err != nil {
				blog.Warnf("move committed event to event queue failed: %v, we will try again later", err)
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

func (th *TxnHandler) watchTransaction() {
	defer th.wg.Done()
	if th.rc == nil {
		return
	}

	stream, err := th.rc.CallStream(daltypes.CommandWatchTransactionOperation, nil)
	if err != nil {
		blog.Errorf("WatchTransaction faile %v", err)
		return
	}
	defer stream.Close()
	txn := daltypes.Transaction{}
	var recvErr error
	for recvErr = stream.Recv(&txn); recvErr != nil || th.shouldClose.IsSet(); recvErr = stream.Recv(&txn) {
		go th.handleTxn(txn)
	}
	if recvErr != nil {
		blog.Errorf("watch stream stoped with error: %v", recvErr)
	}
}

func (th *TxnHandler) fetchTimeout() {
	defer th.wg.Done()
	ticker := util.NewTicker(time.Second * 60)
	opt := redis.ZRangeBy{
		Min: "-inf",
	}
	ticker.Tick()
	for now := range ticker.C {
		if th.shouldClose.IsSet() {
			ticker.Stop()
			break
		}
		txnIDs := []string{}
		opt.Max = strconv.FormatInt(now.UTC().Unix(), 10)

		if err := th.cache.ZRangeByScore(common.EventCacheEventTxnSet, opt).ScanSlice(&txnIDs); err != nil {
			blog.Warnf("fetch timeout txnID from redis failed: %v, we will try again later", err)
			continue
		}

		txns := []daltypes.Transaction{} //Transaction
		if err := th.db.Table(common.BKTableNameTransaction).Find(dal.NewFilterBuilder().In(common.BKTxnIDField, txnIDs)).All(th.ctx, &txns); err != nil {
			blog.Warnf("fetch transaction from mongo failed: %v, we will try again later", err)
			continue
		}
		blog.V(4).Infof("fetch transaction by score %v, resturns %v, txns: %v", opt.Max, txnIDs, txns)
		go th.handleTxn(txns...)
	}
}

func (th *TxnHandler) handleTxn(txns ...daltypes.Transaction) {
	droped := []string{}
	for _, txn := range txns {
		switch txn.Status {
		case daltypes.TxStatusOnProgress:
			continue
		case daltypes.TxStatusCommitted:
			th.committed <- txn.TxnID
		case daltypes.TxStatusAborted, daltypes.TxStatusException:
			droped = append(droped, txn.TxnID)
		default:
			blog.Warnf("unknow transaction status %#v", txn.Status)
		}
	}

	th.dropTransaction(droped)
}

func (th *TxnHandler) dropTransaction(txnIDs []string) {
	if len(txnIDs) <= 0 {
		return
	}
	blog.V(4).Infof("transaction %v should drop", txnIDs)
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
