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
	//"context"
	"errors"
	"fmt"

	//"reflect"
	//"strings"
	//"time"
	//
	//"configcenter/src/common"
	//"configcenter/src/common/blog"
	//"configcenter/src/common/util"
	//"configcenter/src/storage/dal"
	//"configcenter/src/storage/types"
	//
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)


// Errors defined
var (
	ErrSessionInfoNotFound = errors.New("session info not found in storage")

)

type TxnManager struct{}

var redis = map[string][]string{}  //{sessionID: [sessionState, txnNumber]}


func (t *TxnManager) SaveSession(sess mongo.Session) error {
	se := mongo.SessionExposer{}
	info, err := se.GetSessionInfo(sess)
	if err != nil {
		return err
	}
	redis[info.SessionID] = []string{info.SessionState, info.TxnNumber}
	return nil
}


func (t *TxnManager) GetSessionInfoFromStorage(sessionID string) (*mongo.SessionInfo, error){
	v, ok := redis[sessionID]
	if !ok {
		return nil, ErrSessionInfoNotFound
	}
	return &mongo.SessionInfo{SessionID:sessionID, SessionState:v[0], TxnNumber:v[1]}, nil
}


func (t *TxnManager) ConvertToSameSession(sess mongo.Session, sessionID string) error {
	sessInfo, err := t.GetSessionInfoFromStorage(sessionID)
	if err != nil {
		return err
	}
	fmt.Printf("*****ConvertToSameSession***, sessInfo:%#v\n", sessInfo)
	se := &mongo.SessionExposer{}
	return se.SetSessionInfo(sess, sessInfo)
}
