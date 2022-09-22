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

package metadata

import "time"

// TxnOption TODO
type TxnOption struct {
	// transaction timeout time
	// min value: 5 * time.Second
	// default: 5min
	Timeout time.Duration
}

// TxnCapable TODO
type TxnCapable struct {
	Timeout   time.Duration `json:"timeout"`
	SessionID string        `json:"session_id"`
}

// AbortTransactionResult abort transaction result
type AbortTransactionResult struct {
	// Retry defines if the transaction needs to retry, the following are the scenario that needs to retry:
	// 1. the write operation in the transaction conflicts with another transaction,
	// then do retry in the scene layer with server times depends on conditions.
	Retry bool `json:"retry"`
}

// AbortTransactionResponse abort transaction response
type AbortTransactionResponse struct {
	BaseResp               `json:",inline"`
	AbortTransactionResult `json:"data"`
}
