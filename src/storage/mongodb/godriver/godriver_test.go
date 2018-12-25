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
package godriver_test

import (
	"context"
	"encoding/json"
	"testing"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/mongodb/findopt"
	"configcenter/src/storage/mongodb/godriver"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func createConnection() mongodb.CommonClient {
	return godriver.NewClient("mongodb://cc:cc@localhost:27010,localhost:27011,localhost:27012,localhost:27013/cmdb")
}

func executeCommand(t *testing.T, callback func(dbclient mongodb.CommonClient)) {
	dbClient := createConnection()
	require.NoError(t, dbClient.Open())

	callback(dbClient)

	require.NoError(t, dbClient.Close())
}
func TestCRUD(t *testing.T) {

	executeCommand(t, func(dbClient mongodb.CommonClient) {

		// insert
		coll := dbClient.Collection("cc_tmp")
		val := xid.New().String()
		err := coll.InsertOne(context.TODO(), bson.M{"key": "value_" + val}, nil)
		require.NoError(t, err)

		// find
		cond := mongo.NewCondition()
		cond.Element(&mongo.Regex{Key: "key", Val: "value"})
		dataResult := []mapstr.MapStr{}
		err = coll.Find(context.TODO(),
			cond.ToMapStr(),
			&findopt.Many{
				Opts: findopt.Opts{
					Fields: []findopt.FieldItem{
						findopt.FieldItem{
							Name: "_id",
							Hide: true,
						},
						findopt.FieldItem{
							Name: "key",
						},
					},
					Limit: 3,
					Sort: []findopt.SortItem{
						findopt.SortItem{
							Name:       "key",
							Descending: true,
						},
					},
				},
			},
			&dataResult)

		require.NoError(t, err)
		sql, err := cond.ToSQL()
		require.NoError(t, err)
		resultStr, err := json.Marshal(dataResult)
		require.NoError(t, err)
		t.Logf("find data result:%s by the condition:%s", resultStr, sql)
		require.NotEqual(t, 0, len(dataResult))

		// find one
		dataResultOne := mapstr.MapStr{}
		err = coll.FindOne(context.TODO(), cond.ToMapStr(), nil, &dataResultOne)
		require.NoError(t, err)
		resultStr, err = json.Marshal(dataResultOne)
		require.NoError(t, err)
		t.Logf("data result one: %s by the condition:%s", resultStr, sql)

		// find and modify
		dataResultFindUpdate := mapstr.MapStr{}
		update := bsonx.Doc{
			{"$set", bsonx.Document(bsonx.Doc{{"key", bsonx.String("value_test")}})},
		}
		err = coll.FindOneAndModify(context.TODO(), cond.ToMapStr(), update, nil, &dataResultFindUpdate)
		require.NoError(t, err)
		resultStr, err = json.Marshal(dataResultFindUpdate)
		require.NoError(t, err)
		t.Logf("data result find and update:%s", resultStr)

		// delete
		//coll.DeleteOne(context.TODO(),cond.ToMapStr())

	})

}

func TestDatabaseName(t *testing.T) {

	executeCommand(t, func(dbClient mongodb.CommonClient) {
		t.Log("database name:", dbClient.Database().Name())
		require.Equal(t, "cmdb", dbClient.Database().Name())
	})
}

func TestDatabaseHasCollection(t *testing.T) {

	executeCommand(t, func(dbClient mongodb.CommonClient) {
		exists, err := dbClient.Database().HasCollection("cc_tmp")
		require.Equal(t, true, exists)
		require.NoError(t, err)
	})
}

func TestDatabaseDropCollection(t *testing.T) {
	executeCommand(t, func(dbClient mongodb.CommonClient) {
		require.NoError(t, dbClient.Database().DropCollection("cc_tmp"))
	})
}

func TestDatabaseGetCollectionNames(t *testing.T) {
	executeCommand(t, func(dbClient mongodb.CommonClient) {
		collNames, err := dbClient.Database().GetCollectionNames()
		require.NoError(t, err)
		for _, name := range collNames {
			t.Log("colloction:", name)
		}
	})
}
