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

package local

import (
	"context"
	"encoding/base64"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// SessionInfo session information for mongo distributed transactions
type SessionInfo struct {
	TxnNubmer int64
	SessionID string
}

// CmdbReloadSession is used to reset a created session's session id, so that we can
// put all the business operation
func CmdbReloadSession(sess mongo.Session, info *SessionInfo) error {
	xsess, ok := sess.(mongo.XSession)
	if !ok {
		return errors.New("the session is not type XSession")
	}
	clientSession := xsess.ClientSession()

	sessionIDBytes, err := base64.StdEncoding.DecodeString(info.SessionID)
	if err != nil {
		return err
	}
	idx, idDoc := bsoncore.AppendDocumentStart(nil)
	idDoc = bsoncore.AppendBinaryElement(idDoc, "id", session.UUIDSubtype, sessionIDBytes[:])
	idDoc, _ = bsoncore.AppendDocumentEnd(idDoc, idx)

	clientSession.Server.SessionID = idDoc
	clientSession.SessionID = idDoc
	// i.didCommitAfterStart=false
	if info.TxnNubmer > 1 {
		// when the txnNumber is large than 1, it means that it's not the first transaction in
		// this session, we do not need to create a new transaction with this txnNumber and mongodb does
		// not allow this, so we need to change the session status from Starting to InProgressing.
		// set state to InProgressing in a same session id, then we can use the same
		// transaction number as a transaction in a single transaction session.
		// otherwise a error like this will be occured as follows:
		// (NoSuchTransaction) Given transaction number 2 does not match any in-progress transactions.
		// The active transaction number is 1
		clientSession.TransactionState = session.InProgress
	}
	return nil
}

// CmdbPrepareCommitOrAbort set state to InProgress, so that we can commit with other
// operation directly. otherwise mongodriver will do a false commit
func CmdbPrepareCommitOrAbort(sess mongo.Session) error {
	xsess, ok := sess.(mongo.XSession)
	if !ok {
		return errors.New("the session is not type XSession")
	}
	clientSession := xsess.ClientSession()

	clientSession.TransactionState = session.InProgress

	return nil
}

// CmdbContextWithSession set the session into context if context includes session info
func CmdbContextWithSession(ctx context.Context, sess mongo.Session) mongo.SessionContext {
	return mongo.NewSessionContext(ctx, sess)
}
