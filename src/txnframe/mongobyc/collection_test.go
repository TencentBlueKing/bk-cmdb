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

	"configcenter/src/common/mapstr"
	"configcenter/src/txnframe/mongobyc/findopt"
)

type keyval struct {
	Key string `json:"key"`
}

func TestInsertOne(t *testing.T) {

	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	err := mongo.Collection("uri_test").InsertOne(context.Background(), `{"key":"urid"}`, nil)
	if nil != err {
		fmt.Println("failed to insert:", err)
		return
	}

}

func TestInsertMany(t *testing.T) {

	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	err := mongo.Collection("uri_test").InsertMany(context.Background(), []interface{}{`{"key":"uri"}`, `{"key":"uri2"}`}, nil)
	if nil != err {
		fmt.Println("failed to insert:", err)
		return
	}

}

func TestUpdateOne(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	err := mongo.Collection("uri_test").InsertOne(context.Background(), `{"key":"uri"}`, nil)
	if nil != err {
		fmt.Println("failed to insert:", err)
		return
	}

	rst, err := mongo.Collection("uri_test").UpdateOne(context.Background(), `{"key":{"$eq":"uri"}}`, `{"$set":{"key":"urid"}}`, nil)
	if nil != err {
		fmt.Println("failed to update:", err.Error())
		return
	}

	t.Logf("rst:%#v", rst)

}

func TestReplaceOne(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	rst, err := mongo.Collection("uri_test_replace").ReplaceOne(context.Background(), `{"key":"urid"}`, `{"key":"uri"}`, nil)
	if nil != err {
		fmt.Println("failed to replace:", err.Error())
		return
	}

	t.Logf("rst:%#v", rst)

}

func TestUpdateMany(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	err := mongo.Collection("uri_test_rename").InsertOne(context.Background(), `{"key":"uri"}`, nil)
	if nil != err {
		fmt.Println("failed to insert:", err)
		return
	}

	rst, err := mongo.Collection("uri_test_rename").UpdateMany(context.Background(), `{"key":{"$eq":"uri"}}`, `{"$rename":{"key":"uridmany"}}`, nil)
	if nil != err {
		fmt.Println("failed to update:", err.Error())
		return
	}

	t.Logf("rst:%#v", rst)

}

func TestDeleteOne(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	err := mongo.Collection("uri_test").InsertOne(context.Background(), `{"key":"uri"}`, nil)
	if nil != err {
		fmt.Println("failed to insert:", err)
		return
	}

	rst, err := mongo.Collection("uri_test").DeleteOne(context.Background(), `{"key":{"$eq":"urid"}}`, nil)
	if nil != err {
		fmt.Println("failed to update:", err.Error())
		return
	}

	t.Logf("rst:%#v", rst)

}

func TestDeleteMany(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	rst, err := mongo.Collection("uri_test").DeleteMany(context.Background(), `{"key":{"$eq":"urid"}}`, nil)
	if nil != err {
		fmt.Println("failed to update:", err.Error())
		return
	}

	t.Logf("rst:%#v", rst)

}

func TestDeleteCollection(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	if err := mongo.Collection("uri_test").Drop(context.Background()); nil != err {
		fmt.Println("failed to drop collection:", err.Error())
		return
	}
}

func TestCount(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	err := mongo.Collection("uri_test").InsertOne(context.Background(), `{"key":"uri"}`, nil)
	if nil != err {
		fmt.Println("failed to insert:", err)
		return
	}

	cnt, err := mongo.Collection("uri_test").Count(context.Background(), `{"key":"uri"}`)
	if nil != err {
		fmt.Println("failed to count:", err.Error())
		return
	}
	fmt.Println("cnt:", cnt)

}

func TestCreateIndex(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	err := mongo.Collection("uri_test").InsertOne(context.Background(), `{"key":"urid"}`, nil)
	if nil != err {
		fmt.Println("failed to insert:", err)
		return
	}

	err = mongo.Collection("uri_test").CreateIndex(Index{
		Keys: mapstr.MapStr{
			"key": 1,
		},
		Name: "key",
	})

	if nil != err {
		fmt.Println("failed to create index:", err)
		return
	}

}

func TestDeleteIndex(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	err := mongo.Collection("uri_test").DropIndex("key")

	if nil != err {
		fmt.Println("failed to delete index:", err)
		return
	}

}

func TestGetIndex(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	rst, err := mongo.Collection("uri_test").GetIndexes()
	if nil != err {
		fmt.Println("failed to create index:", err)
		return
	}

	t.Logf("indexes rst:%#v", rst)
}

func TestFindCollection(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	err := mongo.Collection("uri_test").InsertMany(context.Background(), []interface{}{`{"keyd":"urid3"}`, `{"key_modify":"uri2"}`}, nil)
	if nil != err {
		fmt.Println("failed to find insert:", err)
		return
	}

	results := []map[string]interface{}{}
	err = mongo.Collection("uri_test").Find(context.Background(), map[string]interface{}{
		"keyd": map[string]interface{}{
			"$regex":   "urid3",
			"$options": "",
		},
	}, &findopt.Many{
		Opts: findopt.Opts{
			Limit: 0,
			Skip:  0,
			Sort:  "",
			Fields: mapstr.MapStr{
				"_id": 0,
			},
		},
	}, &results)
	if nil != err {
		fmt.Println("failed to find:", err)
		return
	}

	fmt.Println("result:", results)
}

func TestFindAndModifyCollection(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	results := map[string]interface{}{}
	err := mongo.Collection("uri_test").FindAndModify(context.Background(), `{"key":"uri"}`, `{"$set":{"key_modify":"modifyd"}}`, &findopt.FindAndModify{
		Opts: findopt.Opts{
			Fields: mapstr.MapStr{
				"key":        1,
				"key_modify": 1,
			},
		},
		Upsert: true,
		New:    true,
	}, &results)

	if nil != err {
		fmt.Println("failed to find:", err)
		return
	}

	fmt.Println("result:", results)
}
