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
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

//func GetClient() (*Mongo, error) {
//	uri := "mongodb://cc:cc@localhost:27010,localhost:27011,localhost:27012,localhost:27013/cmdb"
//	return NewMgo(uri, time.Minute)
//}

func TestDeliverSession(t *testing.T) {
	// client1 start tranction
	client1, err := GetClient()
	require.NoError(t, err)
	tnx1, err := client1.StartSession()
	require.NoError(t, err)
	ctx1 := context.Background()
	defer func() {
		tnx1.EndSession(ctx1)
	}()
	err = tnx1.StartTransaction()
	require.NoError(t, err)

	// client2 op
	client2, err := GetClient()
	require.NoError(t, err)

	// get txn info
	tnxInfo, err := tnx1.TxnInfo()
	require.NoError(t, err)
	header := tnxInfo.IntoHeader(http.Header{})
	header.Set(common.BKHTTPCCRequestID, "xxxxx")
	ctx2 := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		SessionID: header.Get(common.BKHTTPCCTxnSessionID),
		TxnNumber: header.Get(common.BKHTTPCCTransactionNumber),
	})

	coll2 := client2.Table("cc_tranTest")
	// client2 insert
	err = coll2.Insert(ctx2, map[string]string{"client1": "value_bbb"})
	require.NoError(t, err)
	//fmt.Println("has inserted one 2")

	//tnx1 op
	coll1 := tnx1.Table("cc_tranTest")
	err = coll1.Insert(ctx1, []interface{}{map[string]string{"client2-1": "value-many-01"}, map[string]string{"client2-2": "value-many-02"}})
	require.NoError(t, err)
	//fmt.Println("has inserted many")
	// client1 commit
	//err = tnx1.AbortTransaction(ctx1)
	//require.NoError(t, err)
	err = tnx1.CommitTransaction(ctx1)
	require.NoError(t, err)

	/********** ??? why after sesssion1 commit or abort, the op can still be successful, while it will be fail if sesssion2 commit or abort ************/
	//err = coll2.Insert(ctx2, map[string]string{"ccc": "value_ccc"})
	//require.NoError(t, err)
	//fmt.Println("has inserted one ")
}

//// TestCreateApp test create app which will be successful
//func TestCreateApp(t *testing.T) {
//	ClearAllBizSet(t)
//
//	s, err := GetSession()
//	require.NoError(t, err)
//	err = s.StartTransaction()
//	require.NoError(t, err)
//	err = ActivateTransaction(s)
//	require.NoError(t, err)
//
//	se := &mongo.SessionExposer{}
//	info, err := se.GetSessionInfo(s)
//	require.NoError(t, err)
//
//	err = CreateBiz(info)
//	require.NoError(t, err)
//
//	err = CreateSet(info)
//	require.NoError(t, err)
//
//	s.CommitTransaction(context.TODO())
//
//	sNoTxn, err := GetSession()
//	require.NoError(t, err)
//
//	rs, err := SearchAllBiz(sNoTxn)
//	require.NoError(t, err)
//	require.Equal(t, 1, len(rs))
//
//	rs, err = SearchAllSet(sNoTxn)
//	require.NoError(t, err)
//	require.Equal(t, 1, len(rs))
//
//}
//
//// TestCreateAppFailed test create app which will be failed
//func TestCreateAppFailed(t *testing.T) {
//	ClearAllBizSet(t)
//
//	s, err := GetSession()
//	require.NoError(t, err)
//	err = s.StartTransaction()
//	require.NoError(t, err)
//	err = ActivateTransaction(s)
//	require.NoError(t, err)
//
//	se := &mongo.SessionExposer{}
//	info, err := se.GetSessionInfo(s)
//	require.NoError(t, err)
//
//	err = CreateBiz(info)
//	require.NoError(t, err)
//
//	err = CreateSetFailed(info)
//	require.Error(t, err)
//
//	sNoTxn, err := GetSession()
//	require.NoError(t, err)
//
//	rs, err := SearchAllBiz(sNoTxn)
//	require.NoError(t, err)
//	require.Equal(t, 0, len(rs))
//
//	rs, err = SearchAllSet(sNoTxn)
//	require.NoError(t, err)
//	require.Equal(t, 0, len(rs))
//}
//
//// TestCreateAppFailed test create app which will be aborted
//func TestCreateAppAbort(t *testing.T) {
//	ClearAllBizSet(t)
//
//	s, err := GetSession()
//	require.NoError(t, err)
//	err = s.StartTransaction()
//	require.NoError(t, err)
//	err = ActivateTransaction(s)
//	require.NoError(t, err)
//
//	se := &mongo.SessionExposer{}
//	info, err := se.GetSessionInfo(s)
//	require.NoError(t, err)
//
//	err = CreateBiz(info)
//	require.NoError(t, err)
//
//	err = CreateSet(info)
//	require.NoError(t, err)
//
//	s.AbortTransaction(context.TODO())
//
//	sNoTxn, err := GetSession()
//	require.NoError(t, err)
//
//	rs, err := SearchAllBiz(sNoTxn)
//	require.NoError(t, err)
//	require.Equal(t, 0, len(rs))
//
//	rs, err = SearchAllSet(sNoTxn)
//	require.NoError(t, err)
//	require.Equal(t, 0, len(rs))
//}
//
//// TestTxnIsolation test the isolation of transaction
//func TestTxnIsolation(t *testing.T) {
//	ClearAllBizSet(t)
//
//	s, err := GetSession()
//	require.NoError(t, err)
//	err = s.StartTransaction()
//	require.NoError(t, err)
//	err = ActivateTransaction(s)
//	require.NoError(t, err)
//
//	se := &mongo.SessionExposer{}
//	info, err := se.GetSessionInfo(s)
//	require.NoError(t, err)
//
//	err = CreateBizNoTxn()
//	require.NoError(t, err)
//
//	sNoTxn, err := GetSession()
//	require.NoError(t, err)
//
//	rs, err := SearchAllBiz(sNoTxn)
//	require.NoError(t, err)
//	require.Equal(t, 1, len(rs))
//
//	rs, err = SearchAllBiz(s)
//	require.NoError(t, err)
//	require.Equal(t, 0, len(rs))
//
//	err = CreateBiz(info)
//	require.NoError(t, err)
//
//	rs, err = SearchAllBiz(s)
//	require.NoError(t, err)
//	require.Equal(t, 1, len(rs))
//
//	rs, err = SearchAllBiz(sNoTxn)
//	require.NoError(t, err)
//	require.Equal(t, 1, len(rs))
//
//	s.CommitTransaction(context.TODO())
//
//	rs, err = SearchAllBiz(sNoTxn)
//	require.NoError(t, err)
//	require.Equal(t, 2, len(rs))
//}
//
//func ClearAllBizSet(t *testing.T) {
//	s, err := GetSession()
//	require.NoError(t, err)
//
//	err = DeleteAllBiz(s)
//	require.NoError(t, err)
//	rs, err := SearchAllBiz(s)
//	require.NoError(t, err)
//	require.Equal(t, 0, len(rs))
//
//	err = DeleteAllSet(s)
//	require.NoError(t, err)
//	rs, err = SearchAllSet(s)
//	require.NoError(t, err)
//	require.Equal(t, 0, len(rs))
//}
//
//func CreateBiz(info *mongo.SessionInfo) error {
//	s, err := GetSession()
//	if err != nil {
//		return err
//	}
//	err = s.StartTransaction()
//	if err != nil {
//		return err
//	}
//	se := &mongo.SessionExposer{}
//	se.SetSessionInfo(s, info)
//	bizColl := s.Client().Database("cmdb").Collection("app")
//	biz := map[string]string{"app": "appname"}
//	if err := bizColl.Insert(context.Background(), biz); err != nil {
//		return err
//	}
//	return nil
//}
//
//func CreateBizNoTxn() error {
//	s, err := GetSession()
//	if err != nil {
//		return err
//	}
//	bizColl := s.Client().Database("cmdb").Collection("app")
//	biz := map[string]string{"app": "appNoTxn"}
//	if err := bizColl.Insert(context.Background(), biz); err != nil {
//		return err
//	}
//	return nil
//}
//
//func CreateSet(info *mongo.SessionInfo) error {
//	s, err := GetSession()
//	if err != nil {
//		return err
//	}
//	err = s.StartTransaction()
//	if err != nil {
//		return err
//	}
//	se := &mongo.SessionExposer{}
//	se.SetSessionInfo(s, info)
//	setColl := s.Client().Database("cmdb").Collection("set")
//	set := map[string]string{"set": "setname"}
//	if err := setColl.Insert(context.Background(), set); err != nil {
//		return err
//	}
//	return nil
//}
//
//func CreateSetFailed(info *mongo.SessionInfo) error {
//	s, err := GetSession()
//	if err != nil {
//		return err
//	}
//	err = s.StartTransaction()
//	if err != nil {
//		return err
//	}
//	se := &mongo.SessionExposer{}
//	se.SetSessionInfo(s, info)
//	setColl := s.Client().Database("cmdb").Collection("set")
//	set := map[string]string{"set": "failed"}
//	if err := setColl.Insert(context.Background(), set); err != nil {
//		return err
//	}
//	return errors.New("create set failed")
//}
//
//func SearchAllBiz(s mongo.Session) ([]mapstr.MapStr, error) {
//	coll := s.Client().Database("cmdb").Collection("app")
//	resultMany = make([]map[string]string, 0)
//	err := coll.Find(nil).
//	return resultOne, err
//}
//
//func SearchAllSet(s mongo.Session) ([]mapstr.MapStr, error) {
//	coll := s.Client().Database("cmdb").Collection("set")
//	resultMany = make([]map[string]string, 0)
//	err := coll.Find(context.TODO(), mapstr.MapStr{}, nil, &resultOne)
//	return resultOne, err
//}
//
//func DeleteAllBiz(s mongo.Session) error {
//	coll := s.Client().Database("cmdb").Collection("app")
//	return coll.Delete(context.TODO(), mapstr.MapStr{})
//}
//
//func DeleteAllSet(s mongo.Session) error {
//	coll := s.Client().Database("cmdb").Collection("set")
//	return coll.Delete(context.TODO(), mapstr.MapStr{})
//}
//
//func GetSession() (mongo.Session, error) {
//	client, err := GetClient()
//	if err != nil {
//		return nil, err
//	}
//	return client.dbc.StartSession()
//}
