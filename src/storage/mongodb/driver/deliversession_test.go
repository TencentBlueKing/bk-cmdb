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

package driver_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"configcenter/src/common/mapstr"
	"configcenter/src/storage/mongodb"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/require"
)

// TestCreateApp test create app which will be successful
func TestCreateApp(t *testing.T) {
	ClearAllBizSet(t)

	s, err := GetSession()
	require.NoError(t, err)
	err = s.StartTransaction()
	require.NoError(t, err)
	err = ActivateTransaction(s)
	require.NoError(t, err)

	se := &mongo.SessionExposer{}
	info, err := se.GetSessionInfo(s.GetInnerSession())
	require.NoError(t, err)

	err = CreateBiz(info)
	require.NoError(t, err)

	err = CreateSet(info)
	require.NoError(t, err)

	s.CommitTransaction()

	sNoTxn, err := GetSession()
	require.NoError(t, err)

	rs, err := SearchAllBiz(sNoTxn)
	require.NoError(t, err)
	require.Equal(t, 1, len(rs))

	rs, err = SearchAllSet(sNoTxn)
	require.NoError(t, err)
	require.Equal(t, 1, len(rs))

}

// TestCreateAppFailed test create app which will be failed
func TestCreateAppFailed(t *testing.T) {
	ClearAllBizSet(t)

	s, err := GetSession()
	require.NoError(t, err)
	err = s.StartTransaction()
	require.NoError(t, err)
	err = ActivateTransaction(s)
	require.NoError(t, err)

	se := &mongo.SessionExposer{}
	info, err := se.GetSessionInfo(s.GetInnerSession())
	require.NoError(t, err)

	err = CreateBiz(info)
	require.NoError(t, err)

	err = CreateSetFailed(info)
	require.Error(t, err)

	sNoTxn, err := GetSession()
	require.NoError(t, err)

	rs, err := SearchAllBiz(sNoTxn)
	require.NoError(t, err)
	require.Equal(t, 0, len(rs))

	rs, err = SearchAllSet(sNoTxn)
	require.NoError(t, err)
	require.Equal(t, 0, len(rs))
}

// TestCreateAppFailed test create app which will be aborted
func TestCreateAppAbort(t *testing.T) {
	ClearAllBizSet(t)

	s, err := GetSession()
	require.NoError(t, err)
	err = s.StartTransaction()
	require.NoError(t, err)
	err = ActivateTransaction(s)
	require.NoError(t, err)

	se := &mongo.SessionExposer{}
	info, err := se.GetSessionInfo(s.GetInnerSession())
	require.NoError(t, err)

	err = CreateBiz(info)
	require.NoError(t, err)

	err = CreateSet(info)
	require.NoError(t, err)

	s.AbortTransaction()

	sNoTxn, err := GetSession()
	require.NoError(t, err)

	rs, err := SearchAllBiz(sNoTxn)
	require.NoError(t, err)
	require.Equal(t, 0, len(rs))

	rs, err = SearchAllSet(sNoTxn)
	require.NoError(t, err)
	require.Equal(t, 0, len(rs))
}

// TestTxnIsolation test the isolation of transaction
func TestTxnIsolation(t *testing.T) {
	ClearAllBizSet(t)

	s, err := GetSession()
	require.NoError(t, err)
	err = s.StartTransaction()
	require.NoError(t, err)
	err = ActivateTransaction(s)
	require.NoError(t, err)

	se := &mongo.SessionExposer{}
	info, err := se.GetSessionInfo(s.GetInnerSession())
	require.NoError(t, err)

	err = CreateBizNoTxn()
	require.NoError(t, err)

	sNoTxn, err := GetSession()
	require.NoError(t, err)

	rs, err := SearchAllBiz(sNoTxn)
	require.NoError(t, err)
	require.Equal(t, 1, len(rs))

	rs, err = SearchAllBiz(s)
	require.NoError(t, err)
	require.Equal(t, 0, len(rs))

	err = CreateBiz(info)
	require.NoError(t, err)

	rs, err = SearchAllBiz(s)
	require.NoError(t, err)
	require.Equal(t, 1, len(rs))

	rs, err = SearchAllBiz(sNoTxn)
	require.NoError(t, err)
	require.Equal(t, 1, len(rs))

	s.CommitTransaction()

	rs, err = SearchAllBiz(sNoTxn)
	require.NoError(t, err)
	require.Equal(t, 2, len(rs))
}

func ClearAllBizSet(t *testing.T) {
	s, err := GetSession()
	require.NoError(t, err)

	_, err = DeleteAllBiz(s)
	require.NoError(t, err)
	rs, err := SearchAllBiz(s)
	require.NoError(t, err)
	require.Equal(t, 0, len(rs))

	_, err = DeleteAllSet(s)
	require.NoError(t, err)
	rs, err = SearchAllSet(s)
	require.NoError(t, err)
	require.Equal(t, 0, len(rs))
}

func CreateBiz(info *mongo.SessionInfo) error {
	s, err := GetSession()
	if err != nil {
		return err
	}
	err = s.StartTransaction()
	if err != nil {
		return err
	}
	se := &mongo.SessionExposer{}
	se.SetSessionInfo(s.GetInnerSession(), info)
	bizColl := s.Collection("app")
	biz := map[string]string{"app": "appname"}
	if err := bizColl.InsertOne(context.Background(), biz, nil); err != nil {
		return err
	}
	return nil
}

func CreateBizNoTxn() error {
	s, err := GetSession()
	if err != nil {
		return err
	}
	bizColl := s.Collection("app")
	biz := map[string]string{"app": "appNoTxn"}
	if err := bizColl.InsertOne(context.Background(), biz, nil); err != nil {
		return err
	}
	return nil
}

func CreateSet(info *mongo.SessionInfo) error {
	s, err := GetSession()
	if err != nil {
		return err
	}
	err = s.StartTransaction()
	if err != nil {
		return err
	}
	se := &mongo.SessionExposer{}
	se.SetSessionInfo(s.GetInnerSession(), info)
	setColl := s.Collection("set")
	set := map[string]string{"set": "setname"}
	if err := setColl.InsertOne(context.Background(), set, nil); err != nil {
		return err
	}
	return nil
}

func CreateSetFailed(info *mongo.SessionInfo) error {
	s, err := GetSession()
	if err != nil {
		return err
	}
	err = s.StartTransaction()
	if err != nil {
		return err
	}
	se := &mongo.SessionExposer{}
	se.SetSessionInfo(s.GetInnerSession(), info)
	setColl := s.Collection("set")
	set := map[string]string{"set": "failed"}
	if err := setColl.InsertOne(context.Background(), set, nil); err != nil {
		return err
	}
	return errors.New("create set failed")
}

func SearchAllBiz(s mongodb.Session) ([]mapstr.MapStr, error) {
	coll := s.Collection("app")
	resultOne := []mapstr.MapStr{}
	err := coll.Find(context.TODO(), mapstr.MapStr{}, nil, &resultOne)
	return resultOne, err
}

func SearchAllSet(s mongodb.Session) ([]mapstr.MapStr, error) {
	coll := s.Collection("set")
	resultOne := []mapstr.MapStr{}
	err := coll.Find(context.TODO(), mapstr.MapStr{}, nil, &resultOne)
	return resultOne, err
}

func DeleteAllBiz(s mongodb.Session) (*mongodb.DeleteResult, error) {
	coll := s.Collection("app")
	return coll.DeleteMany(context.TODO(), mapstr.MapStr{}, nil)
}

func DeleteAllSet(s mongodb.Session) (*mongodb.DeleteResult, error) {
	coll := s.Collection("set")
	return coll.DeleteMany(context.TODO(), mapstr.MapStr{}, nil)
}

func GetSession() (mongodb.Session, error) {
	client := createConnection()
	err := client.Open()
	if err != nil {
		return nil, err
	}
	session := client.Session().Create()
	err = session.Open()
	if err != nil {
		return nil, err
	}
	return session, nil
}

// ActivateTransaction activate transaction number in mongo server
func ActivateTransaction(s mongodb.Session) error {
	// find
	coll := s.Collection("111")
	resultOne := []mapstr.MapStr{}
	err := coll.Find(context.TODO(), mapstr.MapStr{"aaa": "bbb"}, nil, &resultOne)
	return err
}

func aTestDeliverSession(t *testing.T) {

	var err error

	/*************** session1 start tranction and op ************/
	client1 := createConnection()
	err = client1.Open()
	require.NoError(t, err)
	session1 := client1.Session().Create()
	err = session1.Open()
	require.NoError(t, err)
	defer func() {
		session1.Close()
	}()
	err = session1.StartTransaction()
	require.NoError(t, err)

	// get seesion info
	se := &mongo.SessionExposer{}
	info, err := se.GetSessionInfo(session1.GetInnerSession())
	require.NoError(t, err)
	fmt.Printf("info:%#v", info)

	coll1 := session1.Collection("cc_tranTest")

	// insert one
	err = coll1.InsertOne(context.TODO(), bson.M{"key": "value_aaa"}, nil)
	require.NoError(t, err)
	fmt.Println("has inserted one 1")

	/****************** session2 op *******************/
	client2 := createConnection()
	err = client2.Open()
	require.NoError(t, err)
	session2 := client2.Session().Create()
	err = session2.Open()
	require.NoError(t, err)
	err = session2.StartTransaction()
	require.NoError(t, err)

	// update session by using info
	err = se.SetSessionInfo(session2.GetInnerSession(), info)
	require.NoError(t, err)
	coll2 := session2.Collection("cc_tranTest")

	// insert one
	err = coll2.InsertOne(context.TODO(), bson.M{"key": "value_aaa"}, nil)
	require.NoError(t, err)
	fmt.Println("has inserted one 2")

	se.EndSession(session2.GetInnerSession())

	/******session1 op again and then commit or abort**********/
	// insert many
	err = coll1.InsertMany(context.TODO(), []interface{}{bson.M{"key": "value-many-01"}, bson.M{"key": "value-many-02"}}, nil)
	require.NoError(t, err)
	fmt.Println("has inserted many")

	err = session1.AbortTransaction()
	require.NoError(t, err)
	// err = session1.CommitTransaction()
	// require.NoError(t, err)
	// err = session2.AbortTransaction()
	// require.NoError(t, err)
	// err = session2.CommitTransaction()
	// require.NoError(t, err)

	/********** ??? why after sesssion1 commit or abort, the op can still be successful, while it will be fail if sesssion2 commit or abort ************/
	// err = coll1.InsertOne(context.TODO(), bson.M{"key": "value_bbb"}, nil)
	// require.NoError(t, err)
	// fmt.Println("has inserted one ")
}
