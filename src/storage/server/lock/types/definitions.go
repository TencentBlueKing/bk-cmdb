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
 
package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

var (
	LockPermissionDenied = errors.New("Permission denied")
	LockNotFound         = errors.New("lock not found")
)

var (
	LockIDPrefix = "bkcc"
)

// lcok struct
type Lock struct {
	//  id of this transaction
	TxnID string `json:"txnID"`

	// sub  id of this lock
	SubTxnID string `json:"subTxnID"`

	// lock name is used to define the resources that this lock should be locked
	LockName string `json:"lockName"`

	// timeout means that the time of the client can bear to wait for the lock is locked.
	Timeout time.Duration `json:"timeout"`

	Createtime time.Time `json:"createTime"`
}

// LockResult lock check result
type LockResult struct {
	// the sub txn ID of the txn.
	SubTxnID string `json:"subTxnID"`

	// whether the resources has been locked or not
	Locked bool `json:"locked"`

	// first lock resources TxnID
	LockSubTxnID string `json:"lockSubTxnID"`
}

// GetID   lock tag  ID
func GetID(prefix string) string {
	id := xid.New()
	return fmt.Sprintf("%s-%s", prefix, id.String())
}
