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
	"math/rand"
	"net/http"
	"runtime/debug"
	"time"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/errors"
	"configcenter/src/storage/dal/mongo/local"
)

// Transaction TODO
type Transaction interface {
	// CommitTransaction is to commit the transaction.
	CommitTransaction(ctx context.Context, h http.Header) error
	// AbortTransaction is to abort the transaction.
	AbortTransaction(ctx context.Context, h http.Header) (bool, error)
	// autoRun can only be used at local.
	autoRun(ctx context.Context, h http.Header, run func() error) (bool, error)
}

// NewTransaction TODO
func (t *txn) NewTransaction(h http.Header, opts ...metadata.TxnOption) (Transaction, error) {
	cap, err := local.GenTxnCableAndSetHeader(h, opts...)
	if err != nil {
		return nil, err
	}

	transaction := &transaction{
		sessionID: cap.SessionID,
		timeout:   cap.Timeout,
		client:    t.client,
	}
	return transaction, nil
}

// AutoRunTxn auto run with transaction
func (t *txn) AutoRunTxn(ctx context.Context, h http.Header, run func() error, opts ...metadata.TxnOption) error {
	// to avoid nested txn
	if h.Get(common.TransactionIdHeader) != "" {
		return run()
	}

	startTime := time.Now()

	txn, err := t.NewTransaction(h, opts...)
	if err != nil {
		return ccErr.New(common.CCErrCommStartTransactionFailed, err.Error())
	}

	retry, err := txn.autoRun(ctx, h, run)
	if err != nil {
		return err
	}

	if !retry {
		return nil
	}

	// if need to retry, retry for at most 3 times, each wait for a longer time than the previous one
	rid := util.GetHTTPCCRequestID(h)
	appCode := h.Get(common.BKHTTPRequestAppCode)
	for retryCount := 1; retryCount <= 3; retryCount++ {
		// if the previous operation time exceeds half of the http timeout(25s), do not retry to avoid timeout
		if time.Now().Sub(startTime).Seconds() > 10 {
			return ccErr.New(common.CCErrCommCommitTransactionFailed, "retry conflict transaction timeout")
		}

		blog.Warnf("retry transaction, retry count: %d, app code: %s, rid: %s", retryCount, appCode, rid)
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+300) * time.Duration(retryCount))

		// start a new transaction to retry
		txn, err := t.NewTransaction(h, opts...)
		if err != nil {
			return ccErr.New(common.CCErrCommStartTransactionFailed, err.Error())
		}

		retry, err := txn.autoRun(ctx, h, run)
		if err != nil {
			return err
		}

		if !retry {
			return nil
		}

		// do next retry
	}

	blog.Warnf("retry transaction exceeds maximum count, **skip**, app code: %s, rid: %s", appCode, rid)

	return nil
}

type transaction struct {
	// locked is a flag to indicate whether one of CommitTransaction, AbortTransaction or
	// AutoRun is called. only one of them can be called with in a independent transaction
	// instance. so that we can avoid a transaction to commit or abort for multiple times.
	locked bool

	// enabled the transaction or not.
	enableTxn bool
	sessionID string
	// txnNumber ttl in redis
	timeout time.Duration
	client  rest.ClientInterface
}

// CommitTransaction TODO
// call core service to commit transaction
func (t *transaction) CommitTransaction(ctx context.Context, h http.Header) error {
	if t.locked {
		panic("invalid transaction usage.")
	}
	t.locked = true

	subPath := "/update/transaction/commit"
	body := metadata.TxnCapable{
		Timeout:   t.timeout,
		SessionID: t.sessionID,
	}
	resp := new(metadata.BaseResp)
	err := t.client.Post().
		WithContext(ctx).
		Body(body).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return err
	}

	if !resp.Result {
		return ccErr.New(resp.Code, resp.ErrMsg)
	}

	return nil
}

// AbortTransaction call core service to abort transaction
func (t *transaction) AbortTransaction(ctx context.Context, h http.Header) (bool, error) {
	if t.locked {
		panic("invalid transaction usage.")
	}
	t.locked = true

	subPath := "/update/transaction/abort"
	body := metadata.TxnCapable{
		Timeout:   t.timeout,
		SessionID: t.sessionID,
	}
	resp := new(metadata.AbortTransactionResponse)
	err := t.client.Post().
		WithContext(ctx).
		Body(body).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return false, err
	}

	if err := resp.CCError(); err != nil {
		return false, err
	}

	return resp.AbortTransactionResult.Retry, nil
}

func (t *transaction) autoRun(ctx context.Context, h http.Header, run func() error) (retry bool, err error) {
	rid := util.GetHTTPCCRequestID(h)

	defer func() {
		// if panic ,abort the transaction to avoid WriteConflict when the uncommitted data was processed in another transaction
		if panicErr := recover(); panicErr != nil {
			blog.ErrorfDepthf(3, "run transaction,but server panic, err: %v, rid: %s, debug strace:%s", panicErr, rid, debug.Stack())
			err = ccErr.New(common.CCErrCommInternalServerError, common.GetIdentification()+" Internal Server Error")

			var abortErr error
			retry, abortErr = t.AbortTransaction(ctx, h)
			if abortErr != nil {
				blog.ErrorfDepthf(3, "abort the transaction failed, err: %v, rid: %s", abortErr, rid)
				return
			}
			blog.V(4).InfoDepthf(3, "abort the transaction success. rid: %s", rid)
		}
	}()

	if run == nil {
		return false, errors.New("run function can not be nil")
	}

	if t.locked {
		panic("invalid transaction usage.")
	}

	runErr := run()
	if runErr != nil {
		blog.ErrorfDepthf(2, "run transaction, but run() function failed, err: %v, rid: %s", runErr, rid)

		// check if context is cancelled because of timeout or else, if so, do not need to abort the transaction.
		select {
		case <-ctx.Done():
			blog.Errorf("run transaction, but context is cancelled, err: %v, rid: %s", ctx.Err(), rid)
			return false, runErr
		default:
		}

		// run() logic failed, then abort the transaction.
		retry, err := t.AbortTransaction(ctx, h)
		if err != nil {
			blog.ErrorfDepthf(2, "abort the transaction failed, err: %v, rid: %s", err, rid)
			return false, err
		}
		blog.V(4).InfoDepthf(2, "abort the transaction success. rid: %s", rid)
		if retry {
			return true, nil
		}
		// return the run() original err
		return false, runErr
	}

	// run() logic success, then commit the transaction.
	err = t.CommitTransaction(ctx, h)
	if err != nil {
		blog.ErrorfDepthf(2, "commit the transaction failed, err: %v, rid: %s", err, rid)
		return false, err
	}
	blog.V(4).InfoDepthf(2, "commit the transaction success. rid: %s", rid)

	// roll back the locked flag to true to avoid call transaction function again.
	t.locked = true

	return false, nil
}
