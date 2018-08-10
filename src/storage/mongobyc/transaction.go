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

package mongobyc

// #include "mongo.h"
import "C"

// Transaction transaction operation methods
type Transaction interface {
	StartTransaction() error
	AbortTransaction() error
	CommitTransaction() error
	Collection(collName string) CollectionInterface
}

func newSessionTransaction(clientSession *session) Transaction {
	return &transaction{
		clientSession: clientSession,
	}
}

type transaction struct {
	txnOpts        *C.mongoc_transaction_opt_t
	clientSession  *session
	collectionMaps map[collectionName]CollectionInterface
}

func (t *transaction) StartTransaction() error {
	var err C.bson_error_t
	if !C.mongoc_client_session_start_transaction(t.clientSession.innerSession, t.txnOpts, &err) {
		return TransformError(err)
	}
	return nil
}

func (t *transaction) AbortTransaction() error {
	var err C.bson_error_t
	if !C.mongoc_client_session_abort_transaction(t.clientSession.innerSession, &err) {
		return TransformError(err)
	}
	return nil
}

func (t *transaction) CommitTransaction() error {
	var err C.bson_error_t
	var reply C.bson_t
	defer C.bson_destroy(&reply)

	if !C.mongoc_client_session_commit_transaction(t.clientSession.innerSession, &reply, &err) {
		return TransformError(err)
	}

	return nil
}

func (t *transaction) Collection(collName string) CollectionInterface {
	target, ok := t.collectionMaps[collectionName(collName)]
	if !ok {
		target = newCollectionWithSession(t.clientSession.mongocli, collName, t.clientSession.innerSession)
		t.collectionMaps[collectionName(collName)] = target
	}
	return target
}

func (t *transaction) Close() error {
	for _, coll := range t.collectionMaps {
		switch target := coll.(type) {
		case *collection:
			if err := target.Close(); nil != err {
				return err
			}
		}
	}
	t.collectionMaps = map[collectionName]CollectionInterface{}
	return nil
}
