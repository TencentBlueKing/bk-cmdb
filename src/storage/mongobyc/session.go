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

// Session mongodb session operation methods
type Session interface {
	OpenCloser
	Transaction
}

type session struct {
	*transaction
	mongocli     *client
	innerSession *C.mongoc_client_session_t
	sessionOpts  *C.mongoc_session_opt_t
	txnOpts      *C.mongoc_transaction_opt_t
}

func (s *session) Open() error {

	var err C.bson_error_t
	s.innerSession = C.mongoc_client_start_session(s.mongocli.innerClient, s.sessionOpts, &err)
	if nil == s.innerSession {
		return TransformError(err)
	}
	s.transaction = &transaction{
		txnOpts:        s.txnOpts,
		clientSession:  s,
		collectionMaps: map[collectionName]CollectionInterface{},
	}
	return nil
}

func (s *session) Close() error {
	if nil != s.innerSession {
		C.mongoc_client_session_destroy(s.innerSession)
		s.innerSession = nil
	}

	if nil != s.txnOpts {
		C.mongoc_transaction_opts_destroy(s.txnOpts)
		s.txnOpts = nil
	}

	if nil != s.sessionOpts {
		C.mongoc_session_opts_destroy(s.sessionOpts)
		s.sessionOpts = nil
	}

	s.transaction.Close()
	return nil
}
