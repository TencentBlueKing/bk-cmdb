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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/txnframe/mongobyc"
	"configcenter/src/txnframe/types"
	"context"
	"errors"
	"github.com/rs/xid"
	"time"
)

type TxnManager struct {
	cache map[string]*Session
	db    mongobyc.Client
	ctx   context.Context
}

type Session struct {
	Txninst *types.Tansaction
	mongobyc.Session
}

func New(ctx context.Context, db mongobyc.Client) *TxnManager {
	return &TxnManager{
		db:    db,
		cache: map[string]*Session{},
		ctx:   ctx,
	}
}

func (tm *TxnManager) Run() error {
	select {}
	return nil
}

func (tm *TxnManager) GetSession(txnID string) *Session {
	return tm.cache[txnID]
}

func (tm *TxnManager) CreateTransaction(requestID string, processor string) (*Session, error) {
	session := tm.db.Session().Create()
	err := session.Open()
	if nil != err {
		return nil, err
	}
	err = session.StartTransaction()
	if nil != err {
		return nil, err
	}

	now := time.Now()
	txn := types.Tansaction{
		RequestID:  requestID,
		Processor:  processor,
		TxnID:      xid.New().String(),
		Status:     types.TxStatusOnProgress,
		CreateTime: &now,
		LastTime:   &now,
	}

	err = tm.db.Collection(common.BKTableNameTransaction).InsertOne(tm.ctx, txn, nil)
	if err != nil {
		return nil, err
	}

	inst := &Session{
		Txninst: &txn,
		Session: session,
	}
	tm.cache[txn.TxnID] = inst

	return inst, nil
}

func (tm *TxnManager) Commit(txnID string) error {
	session := tm.GetSession(txnID)
	if session == nil {
		return errors.New("session not found")
	}
	err := session.CommitTransaction()
	if nil != err {
		session.Txninst.Status = types.TxStatusException
	} else {
		session.Txninst.Status = types.TxStatusCommited
	}
	session.Close()

	filter := types.NewFilterBuilder().Eq(common.BKTxnIDField, txnID).Build()
	update := types.Document{
		"status":             session.Txninst.Status,
		common.LastTimeField: time.Now(),
	}

	_, err = tm.db.Collection(common.BKTableNameTransaction).UpdateOne(tm.ctx, filter, update, nil)
	if nil != err {
		blog.Errorf("save transaction [%s] status to %#v faile: %s", txnID, session.Txninst.Status, err.Error())
	}
	return nil
}

func (tm *TxnManager) Abort(txnID string) error {
	session := tm.GetSession(txnID)
	if session == nil {
		return errors.New("session not found")
	}
	err := session.AbortTransaction()
	if nil != err {
		session.Txninst.Status = types.TxStatusException
	} else {
		session.Txninst.Status = types.TxStatusAborted
	}
	session.Close()

	filter := types.NewFilterBuilder().Eq(common.BKTxnIDField, txnID).Build()
	update := types.Document{
		"status":             session.Txninst.Status,
		common.LastTimeField: time.Now(),
	}

	_, err = tm.db.Collection(common.BKTableNameTransaction).UpdateOne(tm.ctx, filter, update, nil)
	if nil != err {
		blog.Errorf("save transaction [%s] status to %#v faile: %s", txnID, session.Txninst.Status, err.Error())
	}
	return nil
}
