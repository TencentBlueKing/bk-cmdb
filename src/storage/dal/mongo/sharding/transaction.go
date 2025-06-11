/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package sharding

import (
	"context"
	"fmt"

	"configcenter/pkg/tenant"
	"configcenter/pkg/tenant/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/mongo/local"
)

// CommitTransaction 提交事务
func (m *ShardingMongoManager) CommitTransaction(ctx context.Context, cap *metadata.TxnCapable) error {
	rid := ctx.Value(common.ContextRequestIDField)

	tenantID := util.GetStrByInterface(ctx.Value(common.ContextRequestTenantField))
	txnTenant, exist := tenant.GetTenant(tenantID)
	if !exist || txnTenant.Status != types.EnabledStatus {
		return fmt.Errorf("transaction tenant id %s is invalid", tenantID)
	}

	sessionInfos, err := m.tm.GetAllSessionInfos(cap.SessionID, txnTenant.Database)
	if err != nil {
		blog.Errorf("get all session infos for session: %s failed, err: %v, rid: %v", cap.SessionID, err, rid)
		return fmt.Errorf("get all session infos failed, err: %v", err)
	}

	for _, sess := range sessionInfos {
		// check if txn number exists, if not, then no db operation with transaction is executed, committing will
		// return an error: "(NoSuchTransaction) Given transaction number 1 does not match any in-progress transactions.
		// The active transaction number is -1.". So we will skip commiting the transaction directly in this situation.
		if sess.TxnNumber == 0 {
			blog.Infof("commit transaction: %s but no transaction to commit, **skip**, rid: %s", sess.SessionID, rid)
			continue
		}

		dbCli, exists := m.dbClientMap[sess.DBID]
		if !exists {
			blog.Errorf("%s session related session %v db id is invalid, rid: %v", cap.SessionID, sess, rid)
			return fmt.Errorf("session db id %s is invalid", sess.DBID)
		}

		reloadSession, err := m.tm.PrepareCommitOrAbort(dbCli.Client(), &sess)
		if err != nil {
			blog.Errorf("commit transaction, but prepare transaction failed, err: %v, rid: %v", err, rid)
			return err
		}
		// reset the transaction state, so that we can commit the transaction after start the
		// transaction immediately.
		if err := local.CmdbPrepareCommitOrAbort(reloadSession); err != nil {
			blog.Errorf("reset the commit transaction state failed, err: %v, rid: %v", err, rid)
			return err
		}

		// we commit the transaction with the session id
		err = reloadSession.CommitTransaction(ctx)
		if err != nil {
			return fmt.Errorf("commit transaction: %s failed, err: %v, rid: %v", sess.SessionID, err, rid)
		}

		if err = m.tm.RemoveTxnNumKey(sess.SessionID); err != nil {
			// this key has ttl, it's ok if we not delete it, cause this key has a ttl.
			blog.Errorf("commit txn, but delete %s txn num key failed, err: %v, rid: %v", sess.SessionID, err, rid)
			// do not return.
		}
	}

	err = m.tm.RemoveSessionKey(cap.SessionID)
	if err != nil {
		// this key has ttl, it's ok if we not delete it, cause this key has a ttl.
		blog.Errorf("commit transaction, but delete txn session: %s key failed, err: %v, rid: %v", cap.SessionID, err,
			rid)
		// do not return.
	}

	return nil
}

// AbortTransaction 取消事务
func (m *ShardingMongoManager) AbortTransaction(ctx context.Context, cap *metadata.TxnCapable) (bool, error) {
	rid := ctx.Value(common.ContextRequestIDField)

	tenantID := util.GetStrByInterface(ctx.Value(common.ContextRequestTenantField))
	txnTenant, exist := tenant.GetTenant(tenantID)
	if !exist || txnTenant.Status != types.EnabledStatus {
		return false, fmt.Errorf("transaction tenant id %s is invalid", tenantID)
	}

	sessionInfos, err := m.tm.GetAllSessionInfos(cap.SessionID, txnTenant.Database)
	if err != nil {
		blog.Errorf("get all session infos for session: %s failed, err: %v, rid: %v", cap.SessionID, err, rid)
		return false, fmt.Errorf("get all session infos failed, err: %v", err)
	}

	for _, sess := range sessionInfos {
		reloadSession, err := m.tm.PrepareCommitOrAbort(m.dbClientMap[sess.DBID].Client(), &sess)
		if err != nil {
			blog.Errorf("abort transaction, but prepare transaction failed, err: %v, rid: %v", err, rid)
			return false, err
		}
		// reset the transaction state, so that we can abort the transaction after start the
		// transaction immediately.
		if err := local.CmdbPrepareCommitOrAbort(reloadSession); err != nil {
			blog.Errorf("reset abort transaction state failed, err: %v, rid: %v", err, rid)
			return false, err
		}

		// we abort the transaction with the session id
		err = reloadSession.AbortTransaction(ctx)
		if err != nil {
			return false, fmt.Errorf("abort transaction: %s failed, err: %v, rid: %v", cap.SessionID, err, rid)
		}

		if err = m.tm.RemoveTxnNumKey(sess.SessionID); err != nil {
			// this key has ttl, it's ok if we not delete it, cause this key has a ttl.
			blog.Errorf("abort txn, but delete %s txn num key failed, err: %v, rid: %v", sess.SessionID, err, rid)
			// do not return.
		}
	}

	err = m.tm.RemoveSessionKey(cap.SessionID)
	if err != nil {
		// this key has ttl, it's ok if we not delete it, cause this key has a ttl.
		blog.Errorf("abort transaction, but delete txn session: %s key failed, err: %v, rid: %v", cap.SessionID, err,
			rid)
		// do not return.
	}

	errorType := m.tm.GetTxnError(cap.SessionID)
	switch errorType {
	// retry when the transaction error type is write conflict, which means the transaction conflicts with another one
	case local.WriteConflictType:
		return true, nil
	}

	return false, nil
}
