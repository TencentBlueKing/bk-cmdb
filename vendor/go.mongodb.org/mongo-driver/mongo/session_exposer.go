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
	"strconv"

	"go.mongodb.org/mongo-driver/x/bsonx"
)

type SessionExposer struct{}

type SessionInfo struct {
	SessionID    string
	TxnNumber    string
	SessionState string
}

// GetSessionInfo get the session info from the param session
func (s *SessionExposer) GetSessionInfo(session Session) (*SessionInfo, error) {
	i, ok := session.(*sessionImpl)
	if !ok {
		return nil, errors.New("the session is not type *sessionImpl")
	}
	info := &SessionInfo{}
	sessionIDBytes, err := i.clientSession.Server.SessionID.MarshalBSON()
	if err != nil {
		return nil, err
	}
	// use base64 to encode, to prevent the invalid header field value like "\x00"
	info.SessionID = base64.StdEncoding.EncodeToString(sessionIDBytes)
	info.TxnNumber = strconv.FormatInt(i.clientSession.Server.TxnNumber, 10)
	info.SessionState = strconv.Itoa(int(i.clientSession.GetState()))
	return info, nil
}

// SetSessionInfo set the session info into the param session, make different session behave like same one
func (s *SessionExposer) SetSessionInfo(session Session, info *SessionInfo) error {
	i, ok := session.(*sessionImpl)
	if !ok {
		return errors.New("the session is not type *sessionImpl")
	}
	sessionIDBytes, err := base64.StdEncoding.DecodeString(info.SessionID)
	if err != nil {
		return err
	}
	doc := bsonx.Doc{}
	err = doc.UnmarshalBSON(sessionIDBytes)
	if err != nil {
		return err
	}
	i.clientSession.Server.SessionID = doc
	i.clientSession.Server.TxnNumber, _ = strconv.ParseInt(info.TxnNumber, 10, 64)
	stateVal, err := strconv.Atoi(info.SessionState)
	if err != nil {
		return err
	}
	i.clientSession.SetState(uint8(stateVal))
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
}
