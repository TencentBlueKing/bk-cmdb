// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
)

type SessionExposer struct{}

type SessionInfo struct {
	SessionID string
	TxnNumber string
}

// GetSessionInfo get the session info from the param session
func (s *SessionExposer) GetSessionInfo(session Session) (*SessionInfo, error) {
	i, ok := session.(*sessionImpl)
	if !ok {
		return nil, errors.New("the session is not type *sessionImpl")
	}
	info := &SessionInfo{}
	sessionIDBytes, _ := i.clientSession.Server.SessionID.MarshalBSON()
	info.SessionID = string(sessionIDBytes)
	info.TxnNumber = strconv.FormatInt(i.clientSession.Server.TxnNumber, 10)
	return info, nil
}

// SetSessionInfo set the session info into the param session and update the session state
func (s *SessionExposer) SetSessionInfo(session Session, info *SessionInfo) error {
	i, ok := session.(*sessionImpl)
	if !ok {
		return errors.New("the session is not type *sessionImpl")
	}
	doc := bsonx.Doc{}
	err := doc.UnmarshalBSON([]byte(info.SessionID))
	if err != nil {
		return err
	}
	i.clientSession.Server.SessionID = doc
	i.clientSession.Server.TxnNumber, _ = strconv.ParseInt(info.TxnNumber, 10, 64)
	// update the session state to InProgress
	i.clientSession.ApplyCommand(description.Server{})
	return nil
}

// EndSession ends the session, just return the session to pool, not AbortTransaction
func (s *SessionExposer) EndSession(session Session) error {
	i, ok := session.(*sessionImpl)
	if !ok {
		return errors.New("the session is not type *sessionImpl")
	}
	i.clientSession.EndSession()
	return nil
}

// SetContextSession set the session into context if context includes session info
func (s *SessionExposer) ContextWithSession(ctx context.Context, sess Session) context.Context {

	return contextWithSession(ctx, sess)

	//// set txn
	//opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	//if ok {
	//	msg.RequestID = opt.RequestID
	//	msg.TxnID = opt.TxnID
	//}
	//if c.TxnID != "" {
	//	msg.TxnID = c.TxnID
	//}
}