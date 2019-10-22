// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"errors"

	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

type SessionExposer struct{}

type SessionInfo struct {
	SessionID bsonx.Doc
	TxnNumber int64
}

// GetSessionInfo get the session info from the param session
func (s *SessionExposer) GetSessionInfo(session Session) (*SessionInfo, error) {
	i, ok := session.(*sessionImpl)
	if !ok {
		return nil, errors.New("the session is not type *sessionImpl")
	}
	info := &SessionInfo{}
	info.SessionID = i.Client.Server.SessionID
	info.TxnNumber = i.Client.Server.TxnNumber
	return info, nil
}

// SetSessionInfo set the session info into the param session and update the session state
func (s *SessionExposer) SetSessionInfo(session Session, info *SessionInfo) error {
	i, ok := session.(*sessionImpl)
	if !ok {
		return errors.New("the session is not type *sessionImpl")
	}
	i.Client.Server.SessionID = info.SessionID
	i.Client.Server.TxnNumber = info.TxnNumber
	// update the session state to InProgress
	i.Client.ApplyCommand()
	return nil
}

// EndSession ends the session, just return the session to pool, not AbortTransaction
func (s *SessionExposer) EndSession(session Session) error {
	i, ok := session.(*sessionImpl)
	if !ok {
		return errors.New("the session is not type *sessionImpl")
	}
	i.Client.EndSession()
	return nil
}
