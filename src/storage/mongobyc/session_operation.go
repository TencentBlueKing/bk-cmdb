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

import (
	"fmt"
	"time"
	"unsafe"
)

// SessionOperation session operation methods
type SessionOperation interface {
	WithReadConcernLevel(level string) SessionOperation
	WithWriteConcernMajority(timeout time.Duration) SessionOperation
	Create() Session
}

func newSessionOperation(mongocli *client) SessionOperation {
	return &sessionOperation{
		txnOpts:      C.mongoc_transaction_opts_new(),
		readConcern:  C.mongoc_read_concern_new(),
		writeConcern: C.mongoc_write_concern_new(),
		sessionOpts:  C.mongoc_session_opts_new(),
		mongocli:     mongocli,
	}
}

type sessionOperation struct {
	txnOpts      *C.mongoc_transaction_opt_t
	readConcern  *C.mongoc_read_concern_t
	writeConcern *C.mongoc_write_concern_t
	sessionOpts  *C.mongoc_session_opt_t
	mongocli     *client
	err          error
}

// WithReadConcernLevel set read concern level
func (t *sessionOperation) WithReadConcernLevel(level string) SessionOperation {
	cstrLevel := C.CString(level)
	defer C.free(unsafe.Pointer(cstrLevel))
	if !C.mongoc_read_concern_set_level(t.readConcern, cstrLevel) {
		t.err = fmt.Errorf("faield to set read concern level (%s)", level)
	}

	C.mongoc_transaction_opts_set_read_concern(t.txnOpts, t.readConcern)
	return t
}

// WithWriteConcernMajority set write concern majority
func (t *sessionOperation) WithWriteConcernMajority(timeout time.Duration) SessionOperation {

	millSeconds := timeout.Nanoseconds() / 1e6

	C.mongoc_write_concern_set_wmajority(t.writeConcern, C.int32_t(millSeconds))
	C.mongoc_transaction_opts_set_write_concern(t.txnOpts, t.writeConcern)

	return t
}

// Create create a new transaction instance
func (t *sessionOperation) Create() Session {

	// reset the default transaction options
	C.mongoc_session_opts_set_default_transaction_opts(t.sessionOpts, t.txnOpts)
	return &session{
		mongocli:    t.mongocli,
		sessionOpts: t.sessionOpts,
		txnOpts:     t.txnOpts,
	}
}
