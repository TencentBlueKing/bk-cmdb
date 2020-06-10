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

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// CommitTransaction to commit transaction
func (s *coreService) CommitTransaction(ctx *rest.Contexts) {
	cap := new(metadata.TxnCapable)
	if err := ctx.DecodeInto(cap); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "decode transaction request body failed, err: %v", err)
		return
	}

	err := s.db.CommitTransaction(ctx.Kit.Ctx, cap)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommCommitTransactionFailed, "commit transaction: %s failed, err: %v", cap.SessionID, err)
		if err := s.AbortEvent(cap.SessionID); err != nil {
			blog.Errorf("AbortEvent failed, err:%v", err)
		}
		return
	}

	if err := s.CommitEvent(cap.SessionID); err != nil {
		blog.Errorf("CommitEvent failed, err:%v", err)
	}

	ctx.RespEntity(nil)
}

// CommitTransaction to abort transaction
func (s *coreService) AbortTransaction(ctx *rest.Contexts) {
	cap := new(metadata.TxnCapable)
	if err := ctx.DecodeInto(cap); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "decode transaction request body failed, err: %v", err)
		return
	}

	if err := s.AbortEvent(cap.SessionID); err != nil {
		blog.Errorf("AbortEvent failed, err:%v", err)
	}

	err := s.db.AbortTransaction(ctx.Kit.Ctx, cap)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommCommitTransactionFailed, "commit transaction: %s failed, err: %v", cap.SessionID, err)
		return
	}

	ctx.RespEntity(nil)
}

// CommitEvent used when a transaction is committed to make all related events valid
func (s *coreService) CommitEvent(txnID string) error {
	return s.rds.LPush(common.EventCacheEventTxnCommitQueueKey, txnID).Err()
}

// AbortEvent used when a transaction is aborted to make all related events invalid
func (s *coreService) AbortEvent(txnID string) error {
	return s.rds.LPush(common.EventCacheEventTxnAbortQueueKey, txnID).Err()
}