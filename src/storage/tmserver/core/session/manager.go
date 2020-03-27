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

package session

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/xid"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/mongodb/options/deleteopt"
	"configcenter/src/storage/mongodb/options/findopt"
	"configcenter/src/storage/tmserver/app/options"
	"configcenter/src/storage/types"
)

type Manager struct {
	enable       bool
	processor    string
	txnLifeLimit time.Duration // second
	cache        map[string]*Session
	session      *Session
	db           mongodb.Client

	eventChan   chan *types.Transaction
	subscribers map[chan<- *types.Transaction]bool

	ctx          context.Context
	sessionMutex sync.Mutex
	pubsubMutex  sync.Mutex
}

func New(ctx context.Context, opt options.TransactionConfig, db mongodb.Client, listen string) (*Manager, error) {
	tm := &Manager{
		enable:       opt.IsTransactionEnable(),
		processor:    listen,
		txnLifeLimit: time.Second * time.Duration(float64(opt.GetTransactionLifetimeSecond())*1.5),
		cache:        map[string]*Session{},
		db:           db,

		eventChan:   make(chan *types.Transaction, 2048),
		subscribers: map[chan<- *types.Transaction]bool{},

		ctx: ctx,
	}
	dbRawSession := tm.db.Session().Create()
	err := dbRawSession.Open()
	if err != nil {
		return nil, err
	}

	tm.session = &Session{
		Session: dbRawSession,
	}

	return tm, nil
}

func (tm *Manager) Subscribe(ch chan<- *types.Transaction) {
	tm.pubsubMutex.Lock()
	tm.subscribers[ch] = true
	tm.pubsubMutex.Unlock()
}

func (tm *Manager) UnSubscribe(ch chan<- *types.Transaction) {
	tm.pubsubMutex.Lock()
	delete(tm.subscribers, ch)
	tm.pubsubMutex.Unlock()
}

func (tm *Manager) Publish() {
	for event := range tm.eventChan {
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

func (tm *Manager) Run() error {
	if tm.enable {
		go tm.reconcileCache()
		go tm.reconcilePersistence()
		go tm.Publish()
	}
	<-tm.ctx.Done()
	return nil
}

func (tm *Manager) reconcileCache() {
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

func (tm *Manager) reconcilePersistence() {
	const Limit int64 = 100
	ticker := time.NewTicker(tm.txnLifeLimit * 2)
	for {
		select {
		case <-tm.ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			blog.Infof("reconciling persistence")
			txns := []types.Transaction{}
			opt := findopt.Many{
				Opts: findopt.Opts{
					Limit: Limit,
					Skip:  0,
				},
			}

			tranCond := mongo.NewCondition()
			tranCond.Element(&mongo.Eq{Key: "status", Val: types.TxStatusOnProgress})
			for {
				err := tm.db.Collection(common.BKTableNameTransaction).Find(tm.ctx, tranCond.ToMapStr(), &opt, &txns)
				if err != nil {
					blog.Errorf("reconcile persistence faile: %v, we will retry %v later", err, tm.txnLifeLimit*2)
					break
				}

				if len(txns) <= 0 {
					break
				}

				for _, txn := range txns {
					if time.Since(txn.LastTime) > tm.txnLifeLimit {

						updateCond := mongo.NewCondition()
						updateCond.Element(&mongo.Eq{Key: common.BKTxnIDField, Val: txn.TxnID})
						update := types.Document{
							"status":             types.TxStatusException,
							common.LastTimeField: time.Now(),
						}
						_, err := tm.db.Collection(common.BKTableNameTransaction).UpdateOne(tm.ctx, updateCond.ToMapStr(), update, nil)
						if nil != err {
							// the reconcile will handle this error, so we will not return this error
							blog.Errorf("save transaction [%s] status to %v faile: %s", txn.TxnID, types.TxStatusException, err.Error())
						}
						ntxn := txn
						tm.eventChan <- &ntxn
					}
				}
				txns = txns[:0]
			}

			removeCond := mongo.NewCondition()
			removeCond.Element(&mongo.Gt{Key: common.LastTimeField, Val: time.Now().Add(time.Hour * 24 * 2)})
			if _, err := tm.db.Collection(common.BKTableNameTransaction).DeleteMany(tm.ctx, removeCond.ToMapStr(), &deleteopt.Many{}); err != nil {
				blog.Errorf("delete outdate transaction faile: %s", err.Error())
			}

			blog.Infof("reconcile persistence finish")
		}
	}
}

func (tm *Manager) GetSession(txnID string) *Session {
	tm.sessionMutex.Lock()
	defer tm.sessionMutex.Unlock()

	if tm.enable && txnID != "" {
		return tm.cache[txnID]
	} else {
		return tm.session
	}
	return nil
}

func (tm *Manager) storeSession(txnID string, session *Session) {
	tm.sessionMutex.Lock()
	tm.cache[txnID] = session
	tm.sessionMutex.Unlock()
}

func (tm *Manager) removeSession(txnID string) {
	tm.sessionMutex.Lock()
	delete(tm.cache, txnID)
	tm.sessionMutex.Unlock()
}

func (tm *Manager) CreateTransaction(requestID string) (*Session, error) {
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
	// start transaction return txnID
	txn.TxnID = tm.newTxnID()

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

	inst := &Session{
		Txninst: &txn,
		Session: session,
	}

	tm.storeSession(txn.TxnID, inst)

	return inst, nil
}

func (tm *Manager) newTxnID() string {
	return tm.processor + "-" + xid.New().String()
}

func (tm *Manager) Commit(txnID string) error {
	if !tm.enable || txnID == "" {
		// not start transaction, return
		return nil
	}
	session := tm.GetSession(txnID)
	if session == nil {
		return errors.New("session not found")
	}
	txnerr := session.CommitTransaction()
	defer func() {
		session.Close()
		tm.removeSession(txnID)
	}()
	if nil != txnerr {
		session.Txninst.Status = types.TxStatusException
	} else {
		session.Txninst.Status = types.TxStatusCommitted
	}
	tm.eventChan <- session.Txninst

	tranCond := mongo.NewCondition()
	tranCond.Element(&mongo.Eq{Key: common.BKTxnIDField, Val: txnID})
	filter := tranCond.ToMapStr()
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

func (tm *Manager) Abort(txnID string) error {
	if !tm.enable || txnID == "" {
		// not start transaction, return
		return nil
	}
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
	tranCond := mongo.NewCondition()
	tranCond.Element(&mongo.Eq{Key: common.BKTxnIDField, Val: txnID})
	filter := tranCond.ToMapStr()
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
