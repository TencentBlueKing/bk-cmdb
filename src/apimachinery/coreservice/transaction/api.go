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
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

func NewTxn(client rest.ClientInterface) Interface {
	return &txn{client: client}
}

type txn struct {
	client rest.ClientInterface
}

// Transaction interface
type Interface interface {
	// StartTransaction 开启新事务
	NewTransaction(h http.Header, opts ...metadata.TxnOption) (Transaction, error)
	// AutoRun is a transaction wrapper. it will automatically commit or abort the
	// transaction depend on the f(), if f() returns with an error, then abort the
	// transaction, otherwise, it will commit the transaction.
	AutoRunTxn(ctx context.Context, h http.Header, run func() error, opts ...metadata.TxnOption) error
}
