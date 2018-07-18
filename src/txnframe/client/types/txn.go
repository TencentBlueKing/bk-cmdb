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

// fingerprints is used to define the resources that this lock should be locked
// lock server will uses these fingerprints to lock the resources that it describes.
type FingerprintsType []string

func (f *FingerprintsType) Add(finger string) {
	*f = append(*f, finger)
}

//
type RollBackType string

const (
	InsertOne  RollBackType = "InsertOne"
	InsertMany RollBackType = "InsertMany"
	UpdateOne  RollBackType = "UpdateOne"
	UpdateMany RollBackType = "UpdateMany"
	DeleteOne  RollBackType = "DeleteOne"
	DeleteMany RollBackType = "DeleteMany"
	Drop       RollBackType = "Drop"
)

type TxnIDType string

type TxnMeta struct {
	// TxnID is the transaction id of this transaction
	TxnID TxnIDType `json:"txnID"`

	// when did the transaction launched, which is a unix nano time value
	CreatedAt int64 `json:"createdAt"`

	// status contains all the sub transaction status
	Status []SubTxnStatus `json:"status"`
}

type SubTxnStatus struct {
	// SubTxnID is the sub transaction ID of this transaction.
	// which is a sub operation of a transaction.
	SubTxnID TxnIDType `json:"subTxnID"`

	// sub transaction's fingerprints for locking the resources.
	Fingerprints FingerprintsType `json:"fingerprints"`

	// rollback records the transaction roll back api function relationships
	RollbackID RollBackType `json:"rollbackID"`

	// Before describes the data before the transaction is done,
	// which is also a snapshot for rollback usage.
	Before interface{} `json:"before"`

	// After describes the data after the transaction is done,
	After interface{} `json:"after"`
}
