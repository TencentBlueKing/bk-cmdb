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

package local

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/redis.v5"
)

// Errors defined
var (
	ErrSessionInfoNotFound = errors.New("session info not found in storage")
	ErrRedisNotInited      = errors.New("redis of TxnManager is not inited")
)

// a transaction manager
type TxnManager struct {
    // enable traction functionality or not.
    enableTransaction bool
	cache *redis.Client
	// transaction timeout time.
	timeout time.Duration
}

var redisCache = map[string][]string{} 
var Sep = "-_-"
var SessPre = "sessinfo_"

// InitTxnManager is to init txn manager, set the redis storage
func (t *TxnManager) InitTxnManager(r *redis.Client) error {
	t.cache = r
	return nil
}

// SaveSessionMock is to save session in a mock storage
func (t *TxnManager) SaveSessionMock(sess mongo.Session) error {
	se := mongo.SessionExposer{}
	info, err := se.GetSessionInfo(sess)
	if err != nil {
		return err
	}
	redisCache[info.SessionID] = []string{info.SessionState, info.TxnNumber}
	return nil
}

// GetSessionInfoFromStorageMock is to get session info from a mock storage
func (t *TxnManager) GetSessionInfoFromStorageMock(sessionID string) (*mongo.SessionInfo, error) {
	v, ok := redisCache[sessionID]
	if !ok {
		return nil, ErrSessionInfoNotFound
	}
	return &mongo.SessionInfo{SessionID: sessionID, SessionState: v[0], TxnNumber: v[1]}, nil
}

// SaveSession is to save session in storage
func (t *TxnManager) SaveSession(sess mongo.Session) error {
	if t.cache == nil {
		return ErrRedisNotInited
	}
	se := mongo.SessionExposer{}
	info, err := se.GetSessionInfo(sess)
	if err != nil {
		return err
	}
	val := info.SessionState + Sep + info.TxnNumber
	return t.cache.Set(SessPre+info.SessionID, val, t.timeout).Err()
}

// SaveSession is to save session in storage
func (t *TxnManager) DeleteSession(sess mongo.Session) error {
	if t.cache == nil {
		return ErrRedisNotInited
	}
	se := mongo.SessionExposer{}
	info, err := se.GetSessionInfo(sess)
	if err != nil {
		return err
	}
	return t.cache.Del(SessPre+info.SessionID).Err()
}

// GetSessionInfoFromStorage is to get session info from storage
func (t *TxnManager) GetSessionInfoFromStorage(sessionID string) (*mongo.SessionInfo, error) {
	if t.cache == nil {
		return nil, ErrRedisNotInited
	}
	v, err := t.cache.Get(SessPre + sessionID).Result()
	if err != nil {
		return nil, err
	}
	if v == "" {
		return nil, ErrSessionInfoNotFound
	}
	items := strings.Split(v, Sep)
	if len(items) != 2 {
		return nil, errors.New(fmt.Sprintf("the session info format in redis is wrong, value:%s", v))
	}
	return &mongo.SessionInfo{SessionID: sessionID, SessionState: items[0], TxnNumber: items[1]}, nil
}

// ConvertToSameSession is to convert a different session to a same session by setting the sessInfo
func (t *TxnManager) ConvertToSameSession(sess mongo.Session, sessionID string) error {
	sessInfo, err := t.GetSessionInfoFromStorage(sessionID)
	if err != nil {
		return err
	}

	se := &mongo.SessionExposer{}
	return se.SetSessionInfo(sess, sessInfo)
}
