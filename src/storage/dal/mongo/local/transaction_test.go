/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package local

import (
	"configcenter/src/common"
	"configcenter/src/storage/dal"
	"net/http"
	"context"
	"fmt"
	"sync"

	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

func GetClient() (*Mongo, error) {
	uri := "mongodb://cc:cc@localhost:27010,localhost:27011,localhost:27012,localhost:27013/cmdb"
	return NewMgo(uri, time.Minute)
}

func TestLocalCommitTransaction(t *testing.T) {
	// client1 start tranction
	client1, err := GetClient()
	require.NoError(t, err)
	tnx1, err := client1.StartSession()
	require.NoError(t, err)

	ctx1 := context.Background()
	defer func() {
		tnx1.EndSession(ctx1)
	}()
	err = tnx1.StartTransaction(ctx1)
	require.NoError(t, err)

	coll1 := tnx1.Table("cc_tranTest")
	err = coll1.Insert(ctx1, []interface{}{map[string]string{"local-1": "value-many-01"}, map[string]string{"local-2": "value-many-02"}})
	require.NoError(t, err)

	//client1 commit
	err = tnx1.CommitTransaction(ctx1)
	require.NoError(t, err)
}

func TestLocalStartTransactionTwice(t *testing.T) {
	// client1 start tranction
	client1, err := GetClient()
	require.NoError(t, err)
	tnx1, err := client1.StartSession()
	require.NoError(t, err)

	ctx1 := context.Background()
	defer func() {
		tnx1.EndSession(ctx1)
	}()
	err = tnx1.StartTransaction(ctx1)
	require.NoError(t, err)

	err = tnx1.StartTransaction(ctx1)
	require.EqualError(t, err, "transaction already in progress")
}

func TestDistributedStartTransactionTwice(t *testing.T) {
	// client1 start tranction
	client1, err := GetClient()
	require.NoError(t, err)
	tnx1, err := client1.StartSession()
	require.NoError(t, err)
	// get txn info
	tnxInfo, err := tnx1.TxnInfo()
	require.NoError(t, err)
	fmt.Printf("*******info*******:%#v\n", tnxInfo)

	ctx1 := context.Background()
	defer func() {
		tnx1.EndSession(ctx1)
	}()
	err = tnx1.StartTransaction(ctx1)
	require.NoError(t, err)

	//client2 op
	client2, err := GetClient()
	require.NoError(t, err)

	// get txn info
	tnxInfo, err = tnx1.TxnInfo()
	require.NoError(t, err)

	header := tnxInfo.IntoHeader(http.Header{})
	header.Set(common.BKHTTPCCRequestID, "xxxxx")
	ctx2 := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		SessionID: header.Get(common.BKHTTPCCTxnSessionID),
		SessionState: header.Get(common.BKHTTPCCTxnSessionState),
		TxnNumber: header.Get(common.BKHTTPCCTransactionNumber),
	})
	fmt.Printf("ctx2:%#v\n", ctx2)
	err = client2.StartTransaction(ctx2)
	require.NoError(t, err)
	err = client2.StartTransaction(ctx2)
	require.NoError(t, err)
}

func TestEmbed(t *testing.T) {
	// client1 start tranction
	client1, err := GetClient()
	require.NoError(t, err)
	tnx1, err := client1.StartSession()
	require.NoError(t, err)
	// get txn info
	tnxInfo, err := tnx1.TxnInfo()
	require.NoError(t, err)
	fmt.Printf("*******info*******:%#v\n", tnxInfo)

	ctx1 := context.Background()
	defer func() {
		tnx1.EndSession(ctx1)
	}()
	err = tnx1.StartTransaction(ctx1)
	require.NoError(t, err)


	//client2 op
	client2, err := GetClient()
	require.NoError(t, err)

	// get txn info
	tnxInfo, err = tnx1.TxnInfo()
	require.NoError(t, err)

	header := tnxInfo.IntoHeader(http.Header{})
	header.Set(common.BKHTTPCCRequestID, "xxxxx")
	ctx2 := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		SessionID: header.Get(common.BKHTTPCCTxnSessionID),
		TxnNumber: header.Get(common.BKHTTPCCTransactionNumber),
	})

	//ctx2 := tnxInfo.IntoContext(context.Background())



	coll2 := client2.Table("cc_tranTest")
	// client2 insert
	err = coll2.Insert(ctx2, map[string]string{"client1": "value_bbb"})
	require.NoError(t, err)
	fmt.Println("has inserted one 2")

	//tnx3 op
	//coll3 := client1.Table("cc_tranTest")
	//err = coll3.Insert(ctx3, []interface{}{map[string]string{"txn3-1": "txn3-01"}, map[string]string{"txn3-2": "txn3-02"}})
	//require.NoError(t, err)
	//
	//err = tnx3.AbortTransaction(ctx3)
	//require.NoError(t, err)

	//tnx1 op
	//coll1 := tnx1.Table("cc_tranTest")
	//err = coll1.Insert(ctx1, []interface{}{map[string]string{"client2-1": "value-many-01"}, map[string]string{"client2-2": "value-many-02"}})
	//require.NoError(t, err)
	//fmt.Println("has inserted many")

	//// client1 abort
	//err = tnx1.AbortTransaction(ctx1)
	//require.NoError(t, err)

	coll1 := tnx1.Table("cc_tranTest")
	err = coll1.Insert(ctx1, []interface{}{map[string]string{"client2-1": "value-many-01"}, map[string]string{"client2-2": "value-many-02"}})
	require.NoError(t, err)

	//client1 commit
	err = tnx1.CommitTransaction(ctx1)
	require.NoError(t, err)

	//err = client2.CommitTransaction(ctx2)
	//require.NoError(t, err)

	//err = coll2.Insert(ctx2, map[string]string{"ccc": "value_ccc"})
	//require.NoError(t, err)
	//fmt.Println("has inserted one ")


	//coll1 = tnx1.Table("cc_tranTest")
	//err = coll1.Insert(ctx1, []interface{}{map[string]string{"client2-1": "aaaa"}, map[string]string{"client2-2": "bbb"}})
	//require.NoError(t, err)

}

// 测试验证sessionID生成的唯一性，保证并发操作时，事务上下文信息不会相同冲突
func aTestSessionID(t *testing.T) {
	m1 := map[string]bool{}
	m2 := map[string]bool{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i:=0;i<100000;i++ {
			id, _ := uuid.New()
			m1[fmt.Sprintf("%v",[16]byte(id))] = true
		}
	}()

	go func() {
		defer wg.Done()
		for i:=0;i<100000;i++ {
			id, _ := uuid.New()
			m2[fmt.Sprintf("%v",[16]byte(id))] = true
		}
	}()
	wg.Wait()
	i := 0
	for k := range m1 {
		i++
		if _, ok := m2[k]; ok {
			t.Errorf("has same sessionID:%s\n", k)
		}
	}
	fmt.Println("compare count:",i)
}