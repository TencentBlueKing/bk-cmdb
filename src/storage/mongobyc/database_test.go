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
	"fmt"
	"testing"
)

func TestCreateEmptyCollection(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	err := mongo.Database().CreateEmptyCollection("uri_empty")
	if nil != err {
		fmt.Println("failed to create empty collection:", err)
		return
	}

}

func TestHasCollection(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	ok, err := mongo.Database().HasCollection("uri_empty")
	if nil != err {
		fmt.Println("failed to check:", err)
		return
	}
	if ok {
		fmt.Println("exists")
	}

}

func TestGetCollectionNames(t *testing.T) {
	InitMongoc()
	defer CleanupMongoc()

	mongo := NewClient("mongodb://test.mongoc:27017/cmdb?replicaSet=repseturi")
	if err := mongo.Open(); nil != err {
		fmt.Println("failed  open:", err)
		return
	}
	defer mongo.Close()

	names, err := mongo.Database().GetCollectionNames()
	if nil != err {
		fmt.Println("failed to check:", err)
		return
	}

	fmt.Println("names:", names)

}
