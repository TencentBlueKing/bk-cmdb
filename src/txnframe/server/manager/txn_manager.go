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

package manager

import (
	"configcenter/src/txnframe/mongobyc"
	"configcenter/src/txnframe/types"
)

type TxnManager struct {
	cache map[string]*Session
	db    mongobyc.Client
}

type Session struct {
	*types.Tansaction
	mongo mongobyc.Session
}

func (s *Session) Txn() mongobyc.Transaction {
	return s.mongo
}

func New() *TxnManager {
	return new(TxnManager)
}

func (tm *TxnManager) Start() error {
	return nil
}

func (tm *TxnManager) Store(txn *types.Tansaction, db mongobyc.Session) {
	tm.cache[txn.TxnID] = &Session{
		Tansaction: txn,
		mongo:      db,
	}
}

func (tm *TxnManager) GetSession(txnID string) *Session {
	return tm.cache[txnID]
}

func (tm *TxnManager) CreateTransaction() *Session {
	session := tm.db.Session().Create()
	session.Open()
	session.CreateTransaction()

	return tm.cache[txnID]
}
