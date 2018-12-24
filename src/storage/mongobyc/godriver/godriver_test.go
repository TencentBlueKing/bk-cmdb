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
	"context"
	"testing"

	"configcenter/src/common/mapstr"
	"configcenter/src/storage/mongobyc"
	"configcenter/src/storage/mongobyc/godriver"

	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/stretchr/testify/require"
)

func createConnection() mongobyc.CommonClient {
	return godriver.NewClient("mongodb://cc:cc@localhost:27010,localhost:27011,localhost:27012,localhost:27013/cmdb")
}

func executeCommand(t *testing.T, callback func(dbclient mongobyc.CommonClient)) {
	dbClient := createConnection()
	require.NoError(t, dbClient.Open())

	callback(dbClient)

	require.NoError(t, dbClient.Close())
}

func TestConectMongo(t *testing.T) {

	dbClient := createConnection()
	err := dbClient.Open()
	require.NoError(t, err)
	err = dbClient.Close()
	require.NoError(t, err)

}

func TestInsertOne(t *testing.T) {

	executeCommand(t, func(dbClient mongobyc.CommonClient) {
		coll := dbClient.Collection("cc_tmp")
		want := bsonx.Elem{Key: "_id", Value: bsonx.ObjectID(objectid.New())}
		doc := bsonx.Doc{want, {"x", bsonx.Int32(1)}}
		err := coll.InsertOne(context.TODO(), doc, nil)
		require.NoError(t, err)
	})

}

func TestFind(t *testing.T) {

	executeCommand(t, func(dbClient mongobyc.CommonClient) {

		coll := dbClient.Collection("cc_tmp")
		result := []mapstr.MapStr{}
		err := coll.Find(context.TODO(), `{"x":1}`, nil, &result)
		require.NoError(t, err)
		t.Log("result:", result)
	})

}

func TestFindOne(t *testing.T) {

	executeCommand(t, func(dbClient mongobyc.CommonClient) {
		coll := dbClient.Collection("cc_tmp")
		result := mapstr.MapStr{}
		err := coll.FindOne(context.TODO(), `{"x":1}`, nil, &result)
		require.NoError(t, err)
		t.Log("result:", result)
	})

}

func TestDatabaseName(t *testing.T) {

	executeCommand(t, func(dbClient mongobyc.CommonClient) {
		t.Log("database name:", dbClient.Database().Name())
		require.Equal(t, "cmdb", dbClient.Database().Name())
	})
}

func TestDatabaseHasCollection(t *testing.T) {

	executeCommand(t, func(dbClient mongobyc.CommonClient) {
		exists, err := dbClient.Database().HasCollection("cc_tmp")
		require.Equal(t, true, exists)
		require.NoError(t, err)
	})
}

func TestDatabaseDropCollection(t *testing.T) {
	executeCommand(t, func(dbClient mongobyc.CommonClient) {
		require.NoError(t, dbClient.Database().DropCollection("cc_tmp"))
	})
}

func TestDatabaseGetCollectionNames(t *testing.T) {
	executeCommand(t, func(dbClient mongobyc.CommonClient) {
		collNames, err := dbClient.Database().GetCollectionNames()
		require.NoError(t, err)
		for _, name := range collNames {
			t.Log("colloction:", name)
		}
	})
}
