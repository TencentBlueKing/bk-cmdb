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

package local

import (
	"context"
	"net/http"
	"testing"
	"time"

	"configcenter/src/common"
	"configcenter/src/storage/dal"

	"github.com/stretchr/testify/require"
)

func BenchmarkLocalCUD(b *testing.B) {
	db, err := NewMgo("192.168.100.130:27010", time.Second*5)
	require.NoError(b, err)

	header := http.Header{}
	header.Set(common.BKHTTPCCRequestID, "xxxxx")
	ctx := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		TxnID:     header.Get(common.BKHTTPCCTransactionID),
	})
	tablename := "tmptest"
	err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": "m"})
	require.NoError(b, err)
	defer db.Table(tablename).Delete(ctx, map[string]interface{}{})

	for i := 0; i < b.N; i++ {

		err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": "a"})
		require.NoError(b, err)

		err = db.Table(tablename).Update(ctx, map[string]interface{}{"name": "a"}, map[string]interface{}{"name": "b"})
		require.NoError(b, err)

		err = db.Table(tablename).Delete(ctx, map[string]interface{}{"name": "b"})
		require.NoError(b, err)

	}
}
