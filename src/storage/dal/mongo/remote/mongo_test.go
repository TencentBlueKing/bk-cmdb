/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package remote

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"configcenter/src/common"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo"
)

func MBenchmarkRemoteCUD(b *testing.B) {
	config := mongo.Config{
		Transaction: "enable",
	}
	db, err := NewWithDiscover(getServerFunc, config)
	require.NoError(b, err)

	header := http.Header{}
	header.Set(common.BKHTTPCCRequestID, "xxxxx")
	orgctx := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		TxnID:     header.Get(common.BKHTTPCCTransactionID),
	})
	tablename := "tmptest"
	err = db.Table(tablename).Insert(orgctx, map[string]interface{}{"name": "m"})
	require.NoError(b, err)
	defer db.Table(tablename).Delete(orgctx, map[string]interface{}{})

	for i := 0; i < b.N; i++ {
		tx, err := db.Start(orgctx)
		require.NoError(b, err)
		header = tx.TxnInfo().IntoHeader(header)
		ctx := util.GetDBContext(context.Background(), header)
		opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
		require.True(b, ok)
		require.NotEmpty(b, opt.RequestID)
		require.NotEmpty(b, opt.TxnID)

		err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": "a"})
		require.NoError(b, err)

		err = db.Table(tablename).Update(ctx, map[string]interface{}{"name": "a"}, map[string]interface{}{"name": "b"})
		require.NoError(b, err)

		err = db.Table(tablename).Delete(ctx, map[string]interface{}{"name": "b"})
		require.NoError(b, err)

		err = db.Commit(ctx)
		require.NoError(b, err)
	}

	db.Close()
}

func BenchmarkRemoteCUDParallel(b *testing.B) {
	db, err := NewWithDiscover(func() ([]string, error) { return []string{"http://192.168.100.130:60008"}, nil }, mongo.Config{Transaction: "enable"})
	require.NoError(b, err)
	tablename := "tmptest"
	header := http.Header{}
	header.Set(common.BKHTTPCCRequestID, "xxxxx")
	orgctx := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		TxnID:     header.Get(common.BKHTTPCCTransactionID),
	})
	err = db.Table(tablename).Insert(orgctx, map[string]interface{}{"name": "m"})
	require.NoError(b, err)
	defer db.Table(tablename).Delete(orgctx, map[string]interface{}{})

	var i int32
	b.RunParallel(func(pb *testing.PB) {
		header := http.Header{}
		header.Set(common.BKHTTPCCRequestID, "xxxxx")
		orgctx := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
			RequestID: header.Get(common.BKHTTPCCRequestID),
			TxnID:     header.Get(common.BKHTTPCCTransactionID),
		})

		ii := atomic.AddInt32(&i, 1)
		aa := fmt.Sprintf("a-%d", ii)
		bb := fmt.Sprintf("b-%d", ii)

		for pb.Next() {
			tx, err := db.Start(orgctx)
			require.NoError(b, err)
			header = tx.TxnInfo().IntoHeader(header)
			ctx := util.GetDBContext(context.Background(), header)
			opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
			require.True(b, ok)
			require.NotEmpty(b, opt.RequestID)
			require.NotEmpty(b, opt.TxnID)

			err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": aa})
			require.NoError(b, err)

			err = db.Table(tablename).Update(ctx, map[string]interface{}{"name": aa}, map[string]interface{}{"name": bb})
			require.NoError(b, err)

			err = db.Table(tablename).Delete(ctx, map[string]interface{}{"name": bb})
			require.NoError(b, err)

			err = db.Commit(ctx)
			require.NoError(b, err)
		}
	})
	db.Close()
}
