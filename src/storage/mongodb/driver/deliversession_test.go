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
	"fmt"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/require"
)

func TestDeliverSession(t *testing.T) {

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
	coll1 := session1.Collection("cc_tranTest")

	// get seesion info
	se := &mongo.SessionExposer{}
	info, err := se.GetSessionInfo(session1.GetInnerSession())
	require.NoError(t, err)
	fmt.Printf("info:%#v", info)

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

	// err = session1.AbortTransaction()
	// require.NoError(t, err)
	err = session1.CommitTransaction()
	require.NoError(t, err)
	// err = session2.AbortTransaction()
	// require.NoError(t, err)
	// err = session2.CommitTransaction()
	// require.NoError(t, err)

	/********** ??? why after sesssion1 commit or abort, the op can still be successful, while it will be fail if sesssion2 commit or abort ************/
	// err = coll1.InsertOne(context.TODO(), bson.M{"key": "value_bbb"}, nil)
	// require.NoError(t, err)
	// fmt.Println("has inserted one ")
}
