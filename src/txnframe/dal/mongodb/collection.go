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

package mongodb

import (
	"context"
	"errors"
	"fmt"

	"configcenter/src/txnframe/client"
	"configcenter/src/txnframe/client/types"
	flt "configcenter/src/txnframe/dal/filter"
	"configcenter/src/txnframe/dal/transaction"
)

type CollectionInterface interface {
	Name() string
	Indexes() IndexView
	Count(ctx context.Context, filter *flt.Filter) (int64, error)
	DeleteMany(ctx context.Context, txnID types.TxnIDType, filter *flt.Filter) (*DeleteResult, error)
	DeleteOne(ctx context.Context, txnID types.TxnIDType, filter *flt.Filter) (*DeleteResult, error)
	Drop(ctx context.Context) error
	Find(ctx context.Context, filter *flt.Filter) (Cursor, error)
	FindOne(ctx context.Context, filter *flt.Filter) *DocumentResult
	InsertMany(ctx context.Context, txnID types.TxnIDType, documents []interface{}) (*InsertManyResult, error)
	InsertOne(ctx context.Context, txnID types.TxnIDType, document interface{}) (*InsertOneResult, error)
	UpdateMany(ctx context.Context, txnID types.TxnIDType, filter *flt.Filter, update interface{}) (*UpdateResult, error)
	UpdateOne(ctx context.Context, txnID types.TxnIDType, filter *flt.Filter, update interface{}) (*UpdateResult, error)
}

type Collection struct {
	MgoClient   MongoCollectionClient
	PreLockPath string
	TxnClient   client.TxnClient
}

func (coll *Collection) Count(ctx context.Context, filter *flt.Filter) (int64, error) {
	return coll.MgoClient.Count(ctx, filter.ToDoc())
}

func (coll *Collection) DeleteMany(ctx context.Context, txnID types.TxnIDType, filter *flt.Filter) (*DeleteResult, error) {
	if len(txnID) == 0 {
		return nil, errors.New("empty transaction id")
	}

	result, err := transaction.NewTxn(coll.TxnClient).
		Try(ctx, txnID, coll.PreLockPath).
		Prepare(func() (rollback types.RollBackType, before, after interface{}, err error) {
			cur, err := coll.MgoClient.Find(ctx, filter.ToDoc())
			if err != nil {
				return types.Unknown, nil, nil, fmt.Errorf("txn, query filter failed, err: %v", err)
			}

			defer cur.Close(ctx)
			bef := make([]map[string]interface{}, 0)
			for cur.Next(ctx) {
				ele := make(map[string]interface{})
				if err := cur.Decode(&ele); err != nil {
					return types.Unknown, nil, nil, fmt.Errorf("txn, query filter failed, err: %v", err)
				}

				bef = append(bef, ele)
			}
			return types.DeleteMany, bef, nil, nil
		}).
		Commit(func() (interface{}, error) {
			rtn, err := coll.MgoClient.DeleteMany(ctx, filter.ToDoc())
			return (*DeleteResult)(rtn), err
		})

	if err != nil {
		return nil, err
	}

	return result.(*DeleteResult), nil
}

func (coll *Collection) DeleteOne(ctx context.Context, txnID types.TxnIDType, filter *flt.Filter) (*DeleteResult, error) {
	if len(txnID) == 0 {
		return nil, errors.New("empty transaction id")
	}

	result, err := transaction.NewTxn(coll.TxnClient).
		Try(ctx, txnID, coll.PreLockPath).
		Prepare(func() (rollback types.RollBackType, before, after interface{}, err error) {
			fResult := coll.MgoClient.FindOne(ctx, filter.ToDoc())
			ele := make(map[string]interface{})
			if err := fResult.Decode(&ele); err != nil {
				return types.Unknown, nil, nil, fmt.Errorf("txn, query find one failed, err: %v", err)
			}

			return types.DeleteOne, ele, nil, nil
		}).
		Commit(func() (interface{}, error) {
			rtn, err := coll.MgoClient.DeleteOne(ctx, filter.ToDoc())
			return (*DeleteResult)(rtn), err
		})

	if err != nil {
		return nil, err
	}

	return result.(*DeleteResult), nil
}

// attention: drop collection does not support transaction function.
func (coll *Collection) Drop(ctx context.Context) error {
	return coll.MgoClient.Drop(ctx)
}

func (coll *Collection) Find(ctx context.Context, filter *flt.Filter) (Cursor, error) {
	rtn, err := coll.MgoClient.Find(ctx, filter.ToDoc())
	return Cursor(rtn), err
}

func (coll *Collection) FindOne(ctx context.Context, filter *flt.Filter) *DocumentResult {
	return (*DocumentResult)(coll.MgoClient.FindOne(ctx, filter.ToDoc()))
}

func (coll *Collection) Indexes() IndexView {
	return IndexView(coll.MgoClient.Indexes())
}

func (coll *Collection) InsertMany(ctx context.Context, txnID types.TxnIDType, documents []interface{}) (*InsertManyResult, error) {
	if len(txnID) == 0 {
		return nil, errors.New("empty transaction id")
	}

	result, err := transaction.NewTxn(coll.TxnClient).
		Try(ctx, txnID, coll.PreLockPath).
		Prepare(func() (rollback types.RollBackType, before, after interface{}, err error) {
			return types.InsertMany, nil, documents, nil
		}).
		Commit(func() (interface{}, error) {
			rtn, err := coll.MgoClient.InsertMany(ctx, documents)
			return (*InsertManyResult)(rtn), err
		})

	if err != nil {
		return nil, err
	}

	return result.(*InsertManyResult), nil
}

func (coll *Collection) InsertOne(ctx context.Context, txnID types.TxnIDType, document interface{}) (*InsertOneResult, error) {
	if len(txnID) == 0 {
		return nil, errors.New("empty transaction id")
	}

	result, err := transaction.NewTxn(coll.TxnClient).
		Try(ctx, txnID, coll.PreLockPath).
		Prepare(func() (rollback types.RollBackType, before, after interface{}, err error) {
			return types.InsertOne, nil, document, nil
		}).
		Commit(func() (interface{}, error) {
			rtn, err := coll.MgoClient.InsertOne(ctx, document)
			return (*InsertOneResult)(rtn), err
		})

	if err != nil {
		return nil, err
	}

	return result.(*InsertOneResult), nil
}

func (coll *Collection) Name() string {
	return coll.MgoClient.Name()
}

func (coll *Collection) UpdateMany(ctx context.Context, txnID types.TxnIDType, filter *flt.Filter, update interface{}) (*UpdateResult, error) {
	if len(txnID) == 0 {
		return nil, errors.New("empty transaction id")
	}

	result, err := transaction.NewTxn(coll.TxnClient).
		Try(ctx, txnID, coll.PreLockPath).
		Prepare(func() (rollback types.RollBackType, before, after interface{}, err error) {
			cur, err := coll.MgoClient.Find(ctx, filter.ToDoc())
			if err != nil {
				return types.Unknown, nil, nil, fmt.Errorf("txn, query filter failed, err: %v", err)
			}

			defer cur.Close(ctx)
			bef := make([]map[string]interface{}, 0)
			for cur.Next(ctx) {
				ele := make(map[string]interface{})
				if err := cur.Decode(&ele); err != nil {
					return types.Unknown, nil, nil, fmt.Errorf("txn, query filter failed, err: %v", err)
				}

				bef = append(bef, ele)
			}
			return types.UpdateMany, bef, update, nil
		}).
		Commit(func() (interface{}, error) {
			rtn, err := coll.MgoClient.UpdateMany(ctx, filter.ToDoc(), update)
			return (*UpdateResult)(rtn), err
		})

	if err != nil {
		return nil, err
	}

	return result.(*UpdateResult), nil
}

func (coll *Collection) UpdateOne(ctx context.Context, txnID types.TxnIDType, filter *flt.Filter, update interface{}) (*UpdateResult, error) {
	if len(txnID) == 0 {
		return nil, errors.New("empty transaction id")
	}

	result, err := transaction.NewTxn(coll.TxnClient).
		Try(ctx, txnID, coll.PreLockPath).
		Prepare(func() (rollback types.RollBackType, before, after interface{}, err error) {
			fResult := coll.MgoClient.FindOne(ctx, filter.ToDoc())

			ele := make(map[string]interface{})
			if err := fResult.Decode(&ele); err != nil {
				return types.Unknown, nil, nil, fmt.Errorf("txn, query find one failed, err: %v", err)
			}

			return types.UpdateOne, ele, nil, nil
		}).
		Commit(func() (interface{}, error) {
			rtn, err := coll.MgoClient.UpdateOne(ctx, filter.ToDoc(), update)
			return (*UpdateResult)(rtn), err
		})

	if err != nil {
		return nil, err
	}

	return result.(*UpdateResult), nil
}
