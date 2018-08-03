/*
 * Tencent is pleased to support the open source community by making čé˛¸ available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"

	"configcenter/src/txnframe/server/mongoc"
)

type keyval struct {
	Key string `json:"key"`
}

func main() {

	mongoc.InitMongoc()
	defer mongoc.CleanupMongoc()

	mongo := mongoc.NewClient("mongodb://test.mongoc:27017", "db_name_uri")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	result, err := mongo.Collection("uri_test").InsertOne(context.Background(), `{"key":"uri"}`)
	if nil != err {
		fmt.Println("failed to insert:", err)
		return
	}

	fmt.Println("inserted cnt:", result.Count)
	findResults := []keyval{}
	err = mongo.Collection("uri_test").Find(context.Background(), `{"key":"uri"}`, &findResults)
	if nil != err {
		fmt.Println("failed to find:", err.Error())
		return
	}

	fmt.Printf("find result:%#v\n", findResults)

}
