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
	"encoding/json"
	"testing"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/mongodb/options/findopt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestCRUD(t *testing.T) {

	executeCommand(t, func(dbClient mongodb.CommonClient) {

		// insert one
		coll := dbClient.Collection("cc_tmp")
		val := xid.New().String()
		err := coll.InsertOne(context.TODO(), bson.M{"key": "value_" + val}, nil)
		require.NoError(t, err)

		// insert many
		err = coll.InsertMany(context.TODO(), []interface{}{bson.M{"key": "value-many_" + val}, bson.M{"key": "value-many_" + xid.New().String()}}, nil)
		require.NoError(t, err)

		// find
		cond := mongo.NewCondition()
		cond.Element(&mongo.Regex{Key: "key", Val: "value_"})
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

		// update one
		updateResult, err := coll.UpdateOne(context.TODO(), cond.ToMapStr(), mapstr.MapStr{"key": "value_update"}, nil)
		require.NoError(t, err)
		require.NotNil(t, updateResult)
		t.Logf("update result:%#v", updateResult)

		// update many
		updateResult, err = coll.UpdateMany(context.TODO(), cond.ToMapStr(), mapstr.MapStr{"key": "value_many_update"}, nil)
		require.NoError(t, err)
		require.NotNil(t, updateResult)
		t.Logf("update many result:%#v", updateResult)

		// replace one
		replaceResult, err := coll.ReplaceOne(context.TODO(), cond.ToMapStr(), mapstr.MapStr{"key": "value-replace_update"}, nil)
		require.NoError(t, err)
		require.NotNil(t, replaceResult)
		t.Logf("replace one result:%#v", replaceResult)

		// delete one
		deleteResult, err := coll.DeleteOne(context.TODO(), cond.ToMapStr(), nil)
		require.NoError(t, err)
		require.NotNil(t, deleteResult)
		t.Logf("delete one result:%#v", deleteResult)

		// delete many
		deleteResult, err = coll.DeleteMany(context.TODO(), cond.ToMapStr(), nil)
		require.NoError(t, err)
		require.NotNil(t, deleteResult)
		t.Logf("delete many result:%#v", deleteResult)

	})

}

func TestIndexCRUD(t *testing.T) {
	executeCommand(t, func(dbClient mongodb.CommonClient) {
		coll := dbClient.Collection("cc_index")

		// construct some datas
		val := xid.New().String()
		err := coll.InsertMany(context.TODO(), []interface{}{bson.M{"key": "value-index_" + val}, bson.M{"key": "value-index_" + xid.New().String()}}, nil)
		require.NoError(t, err)

		cond := mongo.NewCondition()
		cond.Element(&mongo.Regex{Key: "key", Val: "value-index"})
		cnt, err := coll.Count(context.TODO(), cond.ToMapStr())
		require.NoError(t, err)
		require.NotEqual(t, uint64(0), cnt)
		t.Logf("cc_index items count:%d", cnt)

		// create a new index
		err = coll.CreateIndex(mongodb.Index{
			Keys: map[string]int32{"key": -1},
			Name: "test_index_name_key",
		})
		require.NoError(t, err)

		// get indexes
		indexs, err := coll.GetIndexes()
		require.NoError(t, err)
		require.NotNil(t, indexs)
		t.Logf("get cc_index indexs:%#v", indexs)

		// delete index
		err = coll.DropIndex("test_index_name_key")
		require.NoError(t, err)

		// get indexes
		indexs, err = coll.GetIndexes()
		require.NoError(t, err)
		require.NotNil(t, indexs)
		t.Logf("get cc_index indexs:%#v", indexs)
	})
}
