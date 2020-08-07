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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccredis "configcenter/src/storage/dal/redis"
)

/*
TxnHandler commit all event at db transaction commits, elsewhere clear cached events if db transaction abort
*/

func (th *TxnHandler) Run() {
	blog.Info("txn handle process started")
	go th.handleCommit()
	go th.handleAbort()
}

// handleCommit to handle the txn which is committed
func (th *TxnHandler) handleCommit() {
	for {
		// eventStrs format is []string{key, event}
		eventStrs := th.cache.BRPop(time.Second*60, common.EventCacheEventTxnCommitQueueKey).Val()
		if len(eventStrs) == 0 || eventStrs[1] == nilStr || len(eventStrs[1]) == 0 {
			continue
		}

		txnID := eventStrs[1]
		for {
			err := th.cache.RPopLPush(common.EventCacheEventTxnQueuePrefix+txnID, common.EventCacheEventQueueKey).Err()
			if ccredis.IsNilErr(err) {
				break
			}
			if err != nil {
				blog.Errorf("handleCommit RPopLPush failed, txnID:%s, err: %v", txnID, err)
				continue
			}
		}

		blog.V(4).Infof("handleCommit success for txnID:%s", txnID)
	}
}

// handleAbort to handle the txn which is aborted
func (th *TxnHandler) handleAbort() {
	for {
		// eventStrs format is []string{key, event}
		eventStrs := th.cache.BRPop(time.Second*60, common.EventCacheEventTxnAbortQueueKey).Val()
		if len(eventStrs) == 0 || eventStrs[1] == nilStr || len(eventStrs[1]) == 0 {
			continue
		}

		txnID := eventStrs[1]
		if err := th.cache.Del(common.EventCacheEventTxnQueuePrefix + txnID).Err(); err != nil {
			blog.Errorf("handleAbort Del failed, txnID:%s, err:%v", txnID, err)
			continue
		}

		blog.V(4).Infof("handleAbort success for txnID:%s", txnID)
	}
}
