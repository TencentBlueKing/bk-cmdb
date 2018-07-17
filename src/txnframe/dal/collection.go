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

package dal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/txnframe/client"
	"configcenter/src/txnframe/client/lock"
	"configcenter/src/txnframe/client/types"
	flt "configcenter/src/txnframe/dal/filter"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type CollectionInterface interface {
	Name() string
	Indexes() mongo.IndexView
	Drop(ctx context.Context) error
	Count(ctx context.Context, filter *flt.Filter) (int64, error)
	DeleteMany(ctx context.Context, txnID string, filter *flt.Filter) (*mongo.DeleteResult, error)
	DeleteOne(ctx context.Context, txnID string, filter *flt.Filter) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter *flt.Filter) (mongo.Cursor, error)
	FindOne(ctx context.Context, filter *flt.Filter) *mongo.DocumentResult
	InsertMany(ctx context.Context, txnID string, documents []interface{}) (*mongo.InsertManyResult, error)
	InsertOne(ctx context.Context, txnID string, document interface{}) (*mongo.InsertOneResult, error)
	UpdateMany(ctx context.Context, txnID string, filter *flt.Filter, update interface{}) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, txnID string, filter *flt.Filter, update interface{}) (*mongo.UpdateResult, error)
}

const lockTimeOut = 1 * time.Second

type Collection struct {
	mgoClient   MongoCollectionClient
	preLockPath string
	lock        lock.LockInterface
	txnClient   client.TxnClient
}

func (coll *Collection) Count(ctx context.Context, filter *flt.Filter) (int64, error) {
	return coll.mgoClient.Count(ctx, filter.ToDoc())
}

func (coll *Collection) DeleteMany(ctx context.Context, txnID string, filter *flt.Filter) (*mongo.DeleteResult, error) {
	if len(txnID) == 0 {
		return nil, errors.New("empty transaction id")
	}

	// step 1: prepare to lock the resources
	locked, err := coll.txnClient.PreLock(&types.PreLockMeta{
		TxnID:    txnID,
		LockName: coll.preLockPath,
		Timeout:  lockTimeOut,
	})

	if err != nil {
		return nil, fmt.Errorf("prelock resource failed, err: %v", err)
	}

	if !locked {
		return nil, errors.New("prelock resource failed")
	}

	fp := client.NewFingerprints()

	// step 2: find the resources for snapshot
	cur, err := coll.mgoClient.Find(ctx, filter.ToDoc())
	if err != nil {
		return nil, fmt.Errorf("txn, query filter failed, err: %v", err)
	}

	defer cur.Close(ctx)
	before := make([]map[string]interface{}, 0)
	for cur.Next(ctx) {
		ele := make(map[string]interface{})
		if err := cur.Decode(&ele); err != nil {
			return nil, fmt.Errorf("txn, query filter failed, err: %v", err)
		}

		id, ok := getResourceID(ele)
		if !ok {
			return nil, fmt.Errorf("get empty document id failed")
		}
		fp.Add(id)
		before = append(before, ele)
	}

	// step 3: lock the resources
	lockResult, err := coll.txnClient.Lock(&types.LockMeta{
		TxnID:        txnID,
		Fingerprints: fp,
		Timeout:      lockTimeOut,
	})

	if err != nil {
		if lerr := coll.txnClient.PreUnlock(&types.PreUnlockMeta{
			TxnID:    txnID,
			LockName: coll.preLockPath,
		}); lerr != nil {
			blog.Errorf("unlock txn[%s] prelock[%s] failed, err: %v", txnID, coll.preLockPath, lerr)
		}
		return nil, fmt.Errorf("lock the documents failed, err: %v", err)
	}

	if !lockResult.Locked {
		return nil, fmt.Errorf("get the documents's lock failed")
	}

	if err := coll.txnClient.PreUnlock(&types.PreUnlockMeta{
		TxnID:    txnID,
		LockName: coll.preLockPath,
	}); err != nil {
		blog.Errorf("unlock txn[%s] prelock[%s] failed, err: %v", txnID, coll.preLockPath, err)
		return nil, fmt.Errorf("unlock prelock failed, err: %v", err)
	}

	// step 4: snapshot the resources to txn frame.
	snapshot := types.SubTxnStatus{
		SubTxnID:     lockResult.SubTxnID,
		Fingerprints: fp,
		RollbackID:   types.DeleteMany,
		Before:       before,
		After:        nil,
	}
	if err := coll.txnClient.Snapshot(&snapshot); err != nil {
		blog.Errorf("delete many operation, but snapshot meta data failed, txnID[%s], subTxnID[%s], err: %v", txnID, lockResult.SubTxnID, err)
		return nil, fmt.Errorf("snapshort meta data failed, err: %v", err)
	}

	// step 5: delete the resources
	delResult, err := coll.mgoClient.DeleteMany(ctx, filter)
	if err != nil {
		blog.Errorf("delete many failed, txnID[%s], subTxnID[%s], err: %v", txnID, lockResult.SubTxnID, err)
		return nil, fmt.Errorf("delete data failed, err: %v", err)
	}

	// step 6: record delete success operation
	if err := coll.txnClient.SubTxnSuccess(txnID, lockResult.SubTxnID); err != nil {
		blog.Errorf("delete many success, but report sub txn success status failed, txnID[%s], subTxnID[%s], err: %v", txnID, lockResult.SubTxnID, err)
		return nil, fmt.Errorf("reporet sub txn[%s] success status failed, err: %v", lockResult.SubTxnID, err)
	}

	return delResult, nil
}

func (coll *Collection) DeleteOne(ctx context.Context, txnID string, filter *flt.Filter) (*mongo.DeleteResult, error) {

}

// attention: drop collection does not support transaction function.
func (coll *Collection) Drop(ctx context.Context) error {
	return coll.mgoClient.Drop(ctx)
}

func (coll *Collection) Find(ctx context.Context, filter *flt.Filter) (mongo.Cursor, error) {
	return coll.mgoClient.Find(ctx, filter.ToDoc())
}

func (coll *Collection) FindOne(ctx context.Context, filter *flt.Filter) *mongo.DocumentResult {
	return coll.mgoClient.FindOne(ctx, filter.ToDoc())
}

func (coll *Collection) Indexes() mongo.IndexView {
	return coll.mgoClient.Indexes()
}

func (coll *Collection) InsertMany(ctx context.Context, txnID string, documents []interface{}) (*mongo.InsertManyResult, error) {
	return
}

func (coll *Collection) InsertOne(ctx context.Context, txnID string, document interface{}) (*mongo.InsertOneResult, error) {

}

func (coll *Collection) Name() string {
	return coll.mgoClient.Name()
}

func (coll *Collection) UpdateMany(ctx context.Context, txnID string, filter *flt.Filter, update interface{}) (*mongo.UpdateResult, error) {

}

func (coll *Collection) UpdateOne(ctx context.Context, txnID string, filter *flt.Filter, update interface{}) (*mongo.UpdateResult, error) {

}

func getResourceID(ele map[string]interface{}) (hexID string, exist bool) {
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
