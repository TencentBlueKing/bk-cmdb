/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package godriver_test

import (
	"configcenter/src/common/mapstr"
	"context"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"configcenter/src/storage/mongobyc"
	"configcenter/src/storage/mongobyc/godriver"

	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/stretchr/testify/require"
)

func createConnection() mongobyc.CommonClient {
	return godriver.NewClient("mongodb://cc:cc@localhost:27010,localhost:27011,localhost:27012,localhost:27013/cmdb")
}

func TestConectMongo(t *testing.T) {

	dbClient := createConnection()
	err := dbClient.Open()
	require.NoError(t, err)
	err = dbClient.Close()
	require.NoError(t, err)

}

func TestInsertOne(t *testing.T) {

	dbClient := createConnection()
	err := dbClient.Open()
	require.NoError(t, err)
	coll := dbClient.Collection("cc_tmp")
	want := bsonx.Elem{Key: "_id", Value: bsonx.ObjectID(objectid.New())}
	doc := bsonx.Doc{want, {"x", bsonx.Int32(1)}}
	err = coll.InsertOne(context.TODO(), doc, nil)
	require.NoError(t, err)
	err = dbClient.Close()
	require.NoError(t, err)
}

func TestFind(t *testing.T) {

	dbClient := createConnection()
	err := dbClient.Open()
	require.NoError(t, err)

	coll := dbClient.Collection("cc_tmp")
	result := []mapstr.MapStr{}
	err = coll.Find(context.TODO(), `{"x":1}`, nil, &result)
	require.NoError(t, err)
	t.Log("result:", result)

	err = dbClient.Close()
	require.NoError(t, err)

}

func TestFindOne(t *testing.T) {

	dbClient := createConnection()
	err := dbClient.Open()
	require.NoError(t, err)

	coll := dbClient.Collection("cc_tmp")
	result := mapstr.MapStr{}
	err = coll.FindOne(context.TODO(), `{"x":1}`, nil, &result)
	require.NoError(t, err)
	t.Log("result:", result)

	err = dbClient.Close()
	require.NoError(t, err)

}
