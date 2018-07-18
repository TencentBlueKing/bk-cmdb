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

package transaction

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/txnframe/client"
	"configcenter/src/txnframe/client/types"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

const lockTimeOut = 1 * time.Second

type Txn interface {
	Try(ctx context.Context, txnID types.TxnIDType, preLockName string) Txn
	Prepare(prepare func() (rollback types.RollBackType, before, after interface{}, err error)) Txn
	Commit(commit func() (interface{}, error)) (interface{}, error)
}

func NewTxn(txnClient client.TxnClient) Txn {
	return &txn{txnClient: txnClient}
}

type txn struct {
	// indicate whether the Try operation have already been done.
	tried bool

	// indicate whether the Prepare operation have already been done.
	prepared bool

	// indicate whether the Commit operation have already been done.
	committed bool

	// indicate whether pre lock operation have already been done,
	// or the already locked prelock have already been unlocked.
	preLocked bool

	// preLockName is the name of the lock to do the prelock operation.
	preLockName string

	// indicate whether the operation resources have already been locked.
	docsLocked bool

	ctx context.Context

	// txnID is transaction id of this transaction.
	txnID types.TxnIDType

	// subTxnID is the sub transaction if of txnID, which represent of the operation
	// in Prepare operation
	subTxnID types.TxnIDType

	// rollback indicate the rollback operation in this txn
	rollback types.RollBackType

	// indicate whether this txn operation have already occurred an error.
	// if an error occurred, Prepare and Commit operation should return ASAP.
	err error

	// an client to communicate with transaction framework.
	txnClient client.TxnClient
}

func (t *txn) Try(ctx context.Context, txnID types.TxnIDType, preLockName string) Txn {
	if t.prepared {
		panic("can not do prepare before try.")
	}

	if t.committed {
		panic("can not do commit before try.")
	}

	if t.tried {
		panic("can not do try for more than one time.")
	}
	t.tried = true

	t.ctx = ctx
	t.txnID = txnID
	t.preLockName = preLockName

	// step 1: prepare to lock the resources
	locked, err := t.txnClient.PreLock(&types.PreLockMeta{
		TxnID:    txnID,
		LockName: t.preLockName,
		Timeout:  lockTimeOut,
	})

	if err != nil {
		blog.Errorf("do txn[%s] with lock name[%s] failed, err: %v", txnID, preLockName, err)
		t.err = fmt.Errorf("prelock resource failed, err: %v", err)
		return t
	}

	if !locked {
		t.err = errors.New("prelock resource failed")
		return t
	}
	t.preLocked = true

	return t
}

func (t *txn) Prepare(prepare func() (rollback types.RollBackType, before, after interface{}, err error)) Txn {
	if !t.tried {
		panic("can not do prepare before try.")
	}

	if t.committed {
		panic("can not do prepare after commit.")
	}

	if t.prepared {
		panic("can not do prepare for more than one time.")
	}
	t.prepared = true

	if t.err != nil {
		return t
	}

	// step 2: get the resources for snapshot
	rollback, before, after, err := prepare()
	if err != nil {
		t.err = err
		return t
	}

	fp := client.NewFingerprints()

	if before != nil {
		switch before.(type) {
		case []map[string]interface{}:
			bef := before.([]map[string]interface{})
			for idx := range bef {
				id, ok := t.getResourceID(bef[idx])
				if !ok {
					t.err = errors.New("got empty document _id field")
					return t
				}
				fp.Add(id)
			}

		case map[string]interface{}:
			bef := before.(map[string]interface{})
			id, ok := t.getResourceID(bef)
			if !ok {
				t.err = errors.New("got empty document _id field")
				return t
			}
			fp.Add(id)

		default:
			t.err = fmt.Errorf("unsupported before type: %v", reflect.TypeOf(before).Kind().String())
		}
	}

	// step 3: lock the docs resources
	lockResult, err := t.txnClient.Lock(&types.LockMeta{
		TxnID:        t.txnID,
		Fingerprints: fp,
		Timeout:      lockTimeOut,
	})
	if err != nil {
		t.err = fmt.Errorf("lock the documents failed, err: %v", err)
		return t
	}

	if !lockResult.Locked {
		t.err = fmt.Errorf("get the document's lock failed")
		return t
	}
	t.subTxnID = lockResult.SubTxnID

	t.docsLocked = true

	// step 4: already locked the docs resources, unlock the prelock now.
	if err := t.txnClient.PreUnlock(&types.PreUnlockMeta{
		TxnID:    t.txnID,
		LockName: t.preLockName,
	}); err != nil {
		blog.Errorf("unlock txn[%s] prelock[%s] failed, err: %v", t.txnID, t.preLockName, err)
		t.err = fmt.Errorf("unlock prelock failed, err: %v", err)
		return t
	}
	t.preLocked = false

	// step 5: snapshot the resources to txn frame.
	t.rollback = rollback
	snapshot := types.SubTxnStatus{
		SubTxnID:     t.subTxnID,
		Fingerprints: fp,
		RollbackID:   rollback,
		Before:       before,
		After:        after,
	}
	if err := t.txnClient.Snapshot(&snapshot); err != nil {
		blog.Errorf("operation[%s], but snapshot meta data failed, txnID[%s], subTxnID[%s], err: %v", rollback, t.txnID, lockResult.SubTxnID, err)
		t.err = fmt.Errorf("snapshort meta data failed, err: %v", err)
		return t
	}

	return t
}

func (t *txn) Commit(commit func() (interface{}, error)) (interface{}, error) {
	if !t.tried {
		panic("can not do commit before try.")
	}

	if !t.committed {
		panic("can not do commit before prepare.")
	}

	if t.committed {
		panic("can not do commit more than one time.")
	}
	t.committed = true

	if !t.preLocked {
		if err := t.txnClient.PreUnlock(&types.PreUnlockMeta{
			TxnID:    t.txnID,
			LockName: t.preLockName,
		}); err != nil {
			blog.Errorf("unlock txn[%s] prelock[%s] failed, err: %v", t.txnID, t.preLockName, err)
			t.err = fmt.Errorf("unlock prelock failed, err: %v", err)
			return nil, t.err
		}
	}

	// step 6: do the data operation now
	cResult, err := commit()
	if err != nil {
		blog.Errorf("commit %s txn[%s] sub txn[%s] operation failed, err: %v", t.rollback, t.txnID, t.subTxnID, err)
		return nil, err
	}

	// step 7: record delete success operation
	if err := t.txnClient.SubTxnSuccess(t.txnID, t.subTxnID); err != nil {
		blog.Errorf("commit operation[%s] success, but report sub txn success status failed, txnID[%s], subTxnID[%s], err: %v", t.rollback, t.txnID, t.subTxnID, err)
		return nil, fmt.Errorf("reporet sub txn[%s] success status failed, err: %v", t.subTxnID, err)
	}

	return cResult, nil
}

func (t *txn) getResourceID(ele map[string]interface{}) (hexID string, exist bool) {
	idField, exist := ele["_id"]
	if !exist {
		return "", false
	}

	id, ok := idField.(objectid.ObjectID)
	if !ok {
		return "", false
	}

	return id.Hex(), true
}
