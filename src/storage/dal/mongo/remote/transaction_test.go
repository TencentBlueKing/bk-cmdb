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
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"configcenter/src/common"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

func getTMServer() []string {
	ips := strings.TrimSpace(os.Getenv("tmserver_ips"))
	if ips == "" {
		panic("env tmserver_ips require")
	}

	return strings.Split(ips, ",")
}

func getServerFunc() ([]string, error) {
	return getTMServer(), nil
}

func TestTransactionQuery(t *testing.T) {
	NewWithDiscover(getServerFunc)
}

func TestTransaction(t *testing.T) {

	db, err := NewWithDiscover(getServerFunc)
	require.NoError(t, err)

	header := http.Header{}
	header.Set(common.BKHTTPCCRequestID, "xxxxx")
	orgctx := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		TxnID:     header.Get(common.BKHTTPCCTransactionID),
	})
	tablename := "tmptest"
	db.Table(tablename).Delete(orgctx, map[string]interface{}{})
	err = db.Table(tablename).Insert(orgctx, map[string]interface{}{"name": "m"})
	require.NoError(t, err)
	defer db.Table(tablename).Delete(orgctx, map[string]interface{}{})

	tx, err := db.Start(orgctx)
	require.NoError(t, err)
	header = tx.TxnInfo().IntoHeader(header)
	ctx := util.GetDBContext(context.Background(), header)

	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	require.True(t, ok)
	require.NotEmpty(t, opt.RequestID)
	require.NotEmpty(t, opt.TxnID)

	err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": "a"})
	require.NoError(t, err)

	result := []map[string]interface{}{}
	err = db.Table(tablename).Find(map[string]interface{}{"name": "a"}).All(ctx, &result)
	require.NoError(t, err)
	require.Equal(t, 1, len(result))
	require.Equal(t, "a", result[0]["name"])

	count, err := db.Table(tablename).Find(map[string]interface{}{"name": "a"}).Count(ctx)
	require.NoError(t, err)
	require.True(t, count == 1)

	err = db.Table(tablename).Update(ctx, map[string]interface{}{"name": "a"}, map[string]interface{}{"name": "b"})
	require.NoError(t, err)

	result = []map[string]interface{}{}
	err = db.Table(tablename).Find(map[string]interface{}{"name": "b"}).All(ctx, &result)
	require.NoError(t, err)
	require.Equal(t, 1, len(result))
	require.Equal(t, "b", result[0]["name"])

	err = db.Table(tablename).Delete(ctx, map[string]interface{}{"name": "b"})
	require.NoError(t, err)

	result = []map[string]interface{}{}
	err = db.Table(tablename).Find(map[string]interface{}{"name": "b"}).All(ctx, &result)
	require.NoError(t, err)
	require.Equal(t, 0, len(result))

	err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": "c"})
	require.NoError(t, err)

	result = []map[string]interface{}{}
	err = db.Table(tablename).Find(map[string]interface{}{"name": "c"}).All(ctx, &result)
	require.NoError(t, err)
	require.Equal(t, 1, len(result))

	err = tx.Commit(ctx)
	require.NoError(t, err)

	result = []map[string]interface{}{}
	err = db.Table(tablename).Find(map[string]interface{}{"name": "c"}).All(orgctx, &result)
	require.NoError(t, err)
	require.Equal(t, 1, len(result))
	require.Equal(t, "c", result[0]["name"])

	tx, err = db.Start(orgctx)
	require.NoError(t, err)
	header = tx.TxnInfo().IntoHeader(header)
	ctx = util.GetDBContext(context.Background(), header)
	opt, ok = ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	require.True(t, ok)
	require.NotEmpty(t, opt.RequestID)
	require.NotEmpty(t, opt.TxnID)
	err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": "d"})
	require.NoError(t, err)
	err = tx.Abort(ctx)
	require.NoError(t, err)

	result = []map[string]interface{}{}
	err = db.Table(tablename).Find(map[string]interface{}{"name": "d"}).All(orgctx, &result)
	require.NoError(t, err)
	require.Equal(t, 0, len(result))

	count, err = db.Table(tablename).Find(map[string]interface{}{}).Count(orgctx)
	require.NoError(t, err)
	require.EqualValues(t, 2, count)

}

func TestInsertCommit(t *testing.T) {
	db, err := NewWithDiscover(getServerFunc)
	require.NoError(t, err)

	header := http.Header{}
	header.Set(common.BKHTTPCCRequestID, "xxxxx")
	orgctx := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		TxnID:     header.Get(common.BKHTTPCCTransactionID),
	})
	tablename := "tmptest"

	tx, err := db.Start(orgctx)
	require.NoError(t, err)
	header = tx.TxnInfo().IntoHeader(header)
	ctx := util.GetDBContext(context.Background(), header)
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	require.True(t, ok)
	require.NotEmpty(t, opt.RequestID)
	require.NotEmpty(t, opt.TxnID)

	err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": "tmptest_insert_commit"})
	require.NoError(t, err)

	err = db.Commit(ctx)
	require.NoError(t, err)

	db.Close()

}

func TestInsertAbort(t *testing.T) {
	db, err := NewWithDiscover(getServerFunc)
	require.NoError(t, err)

	header := http.Header{}
	header.Set(common.BKHTTPCCRequestID, "xxxxx")
	orgctx := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		TxnID:     header.Get(common.BKHTTPCCTransactionID),
	})
	tablename := "tmptest"

	tx, err := db.Start(orgctx)
	require.NoError(t, err)
	header = tx.TxnInfo().IntoHeader(header)
	ctx := util.GetDBContext(context.Background(), header)
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	require.True(t, ok)
	require.NotEmpty(t, opt.RequestID)
	require.NotEmpty(t, opt.TxnID)

	row := map[string]string{"name": "tmptest_insert_abort"}
	err = db.Table(tablename).Insert(ctx, row)
	require.NoError(t, err)

	newRow := make(map[string]string, 0)
	err = db.Table(tablename).Find(row).One(ctx, &newRow)
	require.NoError(t, err)
	if len(newRow) == 0 {
		t.Errorf("not found data")
		return
	}
	if newRow["name"] != row["name"] {
		t.Errorf(" %#v != %#v", newRow, row)
		return
	}

	err = db.Abort(ctx)
	require.NoError(t, err)

	newRow = make(map[string]string, 0)
	err = db.Table(tablename).Find(row).One(context.Background(), &newRow)
	if !db.IsNotFoundError(err) {
		t.Errorf(err.Error())
		return
	}

	if len(newRow) != 0 {
		t.Errorf("abort txn failure")
		return
	}

	db.Close()

}
