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

import "time"

type PreLockMeta struct {
	// transaction id of this transaction
	TxnID TxnIDType `json:"txnID"`

	// lock name is used to define the resources that this lock should be locked
	LockName string `json:"lockName"`

	// timeout means that the time of the client can bear to wait for the lock is locked.
	Timeout time.Duration `json:"timeout"`
}

type PreUnlockMeta struct {
	// transaction id of this transaction
	TxnID TxnIDType `json:"txnID"`

	// lock name is used to define the resources that this lock should be locked.
	// same with PreLockMeta's LockName
	LockName string `json:"lockName"`
}

type LockMeta struct {
	// transaction id of this transaction
	TxnID TxnIDType `json:"txnID"`

	// fingerprints is used to define the resources that this lock should be locked
	// lock server will uses these fingerprints to lock the resources that it describes.
	Fingerprints FingerprintsType `json:"fingerprints"`

	// timeout means that the time of the client can bear to wait for the lock is locked.
	Timeout time.Duration `json:"timeout"`
}

type LockResult struct {
	// the sub txn ID of the txn.
	SubTxnID TxnIDType `json:"subTxnID"`

	// whether the resources has been locked or not
	Locked bool `json:"locked"`

	// when the resources is locked by the other sub transactions in the same transaction,
	// then this lock can be shared.
	// the user has the right to determine whether the lock should be shared or not.
	CanShare bool `json:"canShare"`
}

type UnlockMeta struct {
	// transaction id of this transaction
	TxnID TxnIDType `json:"txnID"`

	// the sub txn ID of the txn.
	SubTxnID TxnIDType `json:"subTxnID"`
}
