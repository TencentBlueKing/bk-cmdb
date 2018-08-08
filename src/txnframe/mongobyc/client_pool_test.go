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

package mongobyc

import (
	"context"
	"fmt"
	"testing"
)

func TestPoolInsertOne(t *testing.T) {

	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClientPool("mongodb://test.mongoc:27017", "db_name_uri")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	poolCli := mongo.Pop()
	defer mongo.Push(poolCli)

	err := poolCli.Collection("uri_test_pool").InsertOne(context.Background(), `{"key":"uri"}`, nil)
	if nil != err {
		fmt.Println("failed to insert:", err)
		return
	}

}

func TestPoolTransaction(t *testing.T) {

	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClientPool("mongodb://test.mongoc:27017", "db_name_uri")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	poolCli := mongo.Pop()
	defer mongo.Push(poolCli)

	cliSession := poolCli.Session().Create()
	if err := cliSession.Open(); nil != err {
		t.Errorf("failed to  start session: %s", err.Error())
		return
	}
	defer cliSession.Close()

	if err := cliSession.StartTransaction(); nil != err {
		t.Errorf("failed to  start txn: %s", err.Error())
		return
	}

	txnCol := cliSession.Collection("txn_uri_pool")
	txnCol2 := cliSession.Collection("txn_uri_pool2")

	if err := txnCol.InsertOne(context.Background(), `{"txn":"txn_uri_vald3"}`, nil); nil != err {
		t.Errorf("err:%s", err.Error())
		return
	}
	if err := txnCol2.InsertOne(context.Background(), `{"txn":"txn_uri_val3"}`, nil); nil != err {
		t.Errorf("err:%s", err.Error())
		return
	}
	if err := cliSession.CommitTransaction(); nil != err {
		t.Errorf("failed to  commit coll: %s", err.Error())
		return
	}
	t.Logf("finish")
}
