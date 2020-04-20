// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"encoding/base64"
	"errors"

	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

type SessionInfo struct {
	TxnNubmer int64
	SessionID string
}

// GetSessionInfo get the transaction id and txnNumber from a session.
func CmdbGetSessionInfo(sess Session) (string, int64, error) {
	i, ok := sess.(*sessionImpl)
	if !ok {
		return "", 0, errors.New("the session is not type *sessionImpl")
	}
	_, sessID := i.clientSession.Server.SessionID.Lookup("id").Binary()
	return base64.StdEncoding.EncodeToString(sessID[:]),
		i.clientSession.Server.TxnNumber, nil
}

// CmdbReloadSession is used to reset a created session's session id, so that we can 
// put all the business operation 
func CmdbReloadSession(sess Session, info *SessionInfo) error {
	i, ok := sess.(*sessionImpl)
	if !ok {
		panic("the session is not type *sessionImpl")
	}
	sessionIDBytes, err := base64.StdEncoding.DecodeString(info.SessionID)
	if err != nil {
		return err
	}
	idDoc := bsonx.Doc{{"id", bsonx.Binary(session.UUIDSubtype, sessionIDBytes[:])}}
	i.clientSession.Server.SessionID = idDoc
	i.clientSession.SessionID=idDoc
	// i.didCommitAfterStart=false
	if info.TxnNubmer > 1 {
		// when the txnNumber is large than 1, it means that it's not the first transaction in 
		// this session, we do not need to create a new transaction with this txnNumber and mongodb does
		// not allow this, so we need to change the session status from Starting to InProgressing.
		// set state to InProgressing in a same session id, then we can use the same
		// transaction number as a transaction in a single transaction session.
		// otherwise a error like this will be occured as follows:
		// (NoSuchTransaction) Given transaction number 2 does not match any in-progress transactions. The active transaction number is 1
		i.clientSession.SetState(2)
	}
	return nil
}
// CmdbPrepareCommitOrAbort set state to InProgress, so that we can commit with other
// operation directly. otherwise mongodriver will do a false commit
func CmdbPrepareCommitOrAbort(sess Session) {
	i, ok := sess.(*sessionImpl)
	if !ok {
		panic("the session is not type *sessionImpl")
	}

	i.clientSession.SetState(2)
	i.didCommitAfterStart=false
}

// CmdbContextWithSession set the session into context if context includes session info
func CmdbContextWithSession(ctx context.Context, sess Session) SessionContext {
	return contextWithSession(ctx, sess)
}

// CmdbReleaseSession is almost same with session.EndSession(), the difference is 
// that CmdbReleaseSession do not abrot the transaction, and just release the net connection
// it panic if it's not a valid sessionImpl
// Note: do not use this, because our transaction plan do not allows to reuse a session server.
// if we do this, the session server will be reused, and the txnNumber will be increased, we 
// do not allow it happen.
func CmdbReleaseSession(ctx context.Context, sess Session)  {
	i, ok := sess.(*sessionImpl)
	if !ok {
		panic("the session is not type *sessionImpl")
	}
	i.clientSession.EndSession()
}
