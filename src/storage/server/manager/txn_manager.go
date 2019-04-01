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
	"context"
	"errors"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/mongobyc"
	"configcenter/src/storage/server/app/options"
	"configcenter/src/storage/types"

	"github.com/rs/xid"
)

type TxnManager struct {
	enable       bool
	processor    string
	txnLifeLimit time.Duration // second
	cache        map[string]*Session
	db           mongobyc.Client

	eventChan   chan *types.Transaction
	subscribers map[chan<- *types.Transaction]bool

	ctx          context.Context
	sessionMutex sync.Mutex
	pubsubMutex  sync.Mutex
}

type Session struct {
	Txninst *types.Transaction
	mongobyc.Session
}

func New(ctx context.Context, opt options.TransactionConfig, db mongobyc.Client, listen string) *TxnManager {
	tm := &TxnManager{
		enable:       opt.ShouldEnable(),
		processor:    listen,
		txnLifeLimit: time.Second * time.Duration(float64(opt.GetTransactionLifetimeSecond())*1.5),
		cache:        map[string]*Session{},
		db:           db,

		eventChan:   make(chan *types.Transaction, 2048),
		subscribers: map[chan<- *types.Transaction]bool{},

		ctx: ctx,
	}
	return tm
}

func (tm *TxnManager) Subscribe(ch chan<- *types.Transaction) {
	tm.pubsubMutex.Lock()
	tm.subscribers[ch] = true
	tm.pubsubMutex.Unlock()
}

func (tm *TxnManager) UnSubscribe(ch chan<- *types.Transaction) {
	tm.pubsubMutex.Lock()
	delete(tm.subscribers, ch)
	tm.pubsubMutex.Unlock()
}

func (tm *TxnManager) Publish() {
	for {
		event := <-tm.eventChan
		tm.pubsubMutex.Lock()
		for subscriber := range tm.subscribers {
			select {
			case subscriber <- event:
			case <-time.After(time.Second):
			}
		}
		tm.pubsubMutex.Unlock()
	}
}

func (tm *TxnManager) Run() error {
	if tm.enable {
		go tm.reconcileCache()
		go tm.reconcilePersistence()
		go tm.Publish()
	}
	<-tm.ctx.Done()
	return nil
}

func (tm *TxnManager) reconcileCache() {
	ticker := time.NewTicker(tm.txnLifeLimit)
	for {
		select {
		case <-tm.ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			tm.sessionMutex.Lock()
			for _, session := range tm.cache {
				if time.Since(session.Txninst.LastTime) > tm.txnLifeLimit {
					// ignore the abort error, cause the session will not be used again
					go tm.Abort(session.Txninst.TxnID)
				}
			}
			tm.sessionMutex.Unlock()
		}
	}
}

func (tm *TxnManager) reconcilePersistence() {
	ticker := time.NewTicker(tm.txnLifeLimit * 2)
	for {
		select {
		case <-tm.ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			txns := []types.Transaction{}
			err := tm.db.Collection(common.BKTableNameTransaction).Find(tm.ctx, nil, nil, &txns)
			if err != nil {
				blog.Errorf("reconcile persistence faile: %v, we will retry %v later", err, tm.txnLifeLimit)
				continue
			}

			for _, txn := range txns {
				if time.Since(txn.LastTime) > tm.txnLifeLimit {
					filter := dal.NewFilterBuilder().Eq(common.BKTxnIDField, txn.TxnID).Build()
					update := types.Document{
						"status":             types.TxStatusException,
						common.LastTimeField: time.Now(),
					}
					_, err := tm.db.Collection(common.BKTableNameTransaction).UpdateOne(tm.ctx, filter, update, nil)
					if nil != err {
						// the reconcile will handle this error, so we will not return this error
						blog.Errorf("save transaction [%s] status to %#v faile: %s", txn.TxnID, txn.Status, err.Error())
					}
					ntxn := txn
					tm.eventChan <- &ntxn
				}
			}
		}
	}
}

func (tm *TxnManager) GetSession(txnID string) *Session {
	tm.sessionMutex.Lock()
	session := tm.cache[txnID]
	tm.sessionMutex.Unlock()
	return session
}

func (tm *TxnManager) storeSession(txnID string, session *Session) {
	tm.sessionMutex.Lock()
	tm.cache[txnID] = session
	tm.sessionMutex.Unlock()
}

func (tm *TxnManager) removeSession(txnID string) {
	tm.sessionMutex.Lock()
	delete(tm.cache, txnID)
	tm.sessionMutex.Unlock()
}

func (tm *TxnManager) CreateTransaction(requestID string) (*Session, error) {
	txn := types.Transaction{
		RequestID:  requestID,
		Processor:  tm.processor,
		Status:     types.TxStatusOnProgress,
		CreateTime: time.Now(),
		LastTime:   time.Now(),
	}

	if !tm.enable {
		return &Session{
			Txninst: &txn,
		}, nil
	}
	session := tm.db.Session().Create()
	err := session.Open()
	defer func() {
		if err != nil {
			session.Close()
		}
	}()
	if nil != err {
		return nil, err
	}
	err = session.StartTransaction()
	if nil != err {
		return nil, err
	}

	err = tm.db.Collection(common.BKTableNameTransaction).InsertOne(tm.ctx, txn, nil)
	if err != nil {
		// we should return this error,
		// cause the transaction life cycle will not under txn manager's controll
		return nil, err
	}

	txn.TxnID = tm.newTxnID()
	inst := &Session{
		Txninst: &txn,
		Session: session,
	}

	tm.storeSession(txn.TxnID, inst)

	return inst, nil
}

func (tm *TxnManager) newTxnID() string {
	return tm.processor + "-" + xid.New().String()
}

func (tm *TxnManager) Commit(txnID string) error {
	session := tm.GetSession(txnID)
	if session == nil {
		return errors.New("session not found")
	}
	txnerr := session.CommitTransaction()
	defer func() {
		if txnerr == nil {
			session.Close()
			tm.removeSession(txnID)
		}
	}()
	if nil != txnerr {
		session.Txninst.Status = types.TxStatusException
	} else {
		session.Txninst.Status = types.TxStatusCommitted
	}
	tm.eventChan <- session.Txninst

	filter := dal.NewFilterBuilder().Eq(common.BKTxnIDField, txnID).Build()
	update := types.Document{
		"status":             session.Txninst.Status,
		common.LastTimeField: time.Now(),
	}
	_, err := tm.db.Collection(common.BKTableNameTransaction).UpdateOne(tm.ctx, filter, update, nil)
	if nil != err {
		// the reconcile will handle this error, so we will not return this error
		blog.Errorf("save transaction [%s] status to %#v faile: %s", txnID, session.Txninst.Status, err.Error())
	}
	return nil
}

func (tm *TxnManager) Abort(txnID string) error {
	session := tm.GetSession(txnID)
	if session == nil {
		return errors.New("session not found")
	}
	txnerr := session.AbortTransaction()
	defer func() {
		session.Close()
		tm.removeSession(txnID)
	}()
	if nil != txnerr {
		session.Txninst.Status = types.TxStatusException
	} else {
		session.Txninst.Status = types.TxStatusAborted
	}
	tm.eventChan <- session.Txninst

	filter := dal.NewFilterBuilder().Eq(common.BKTxnIDField, txnID).Build()
	update := types.Document{
		"status":             session.Txninst.Status,
		common.LastTimeField: time.Now(),
	}

	_, err := tm.db.Collection(common.BKTableNameTransaction).UpdateOne(tm.ctx, filter, update, nil)
	if nil != err {
		// the reconcile will handle this error, so we will not return this error
		blog.Errorf("save transaction [%s] status to %#v faile: %s", txnID, session.Txninst.Status, err.Error())
	}
	return nil
}
