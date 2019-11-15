/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package local

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"

	"github.com/stretchr/testify/require"
)

func BenchmarkLocalCUD(b *testing.B) {

	db, err := NewMgo("mongodb://cc:cc@localhost:27011,localhost:27012,localhost:27013,localhost:27014/cmdb", time.Second*5)
	require.NoError(b, err)

	header := http.Header{}
	header.Set(common.BKHTTPCCRequestID, "xxxxx")
	ctx := context.WithValue(context.Background(), common.CCContextKeyJoinOption, dal.JoinOption{
		RequestID: header.Get(common.BKHTTPCCRequestID),
		TxnID:     header.Get(common.BKHTTPCCTransactionID),
	})
	tablename := "tmptest"
	err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": "m"})
	require.NoError(b, err)
	defer db.Table(tablename).Delete(ctx, map[string]interface{}{})

	for i := 0; i < b.N; i++ {

		err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": "a"})
		require.NoError(b, err)

		err = db.Table(tablename).Update(ctx, map[string]interface{}{"name": "a"}, map[string]interface{}{"name": "b"})
		require.NoError(b, err)

		err = db.Table(tablename).Delete(ctx, map[string]interface{}{"name": "b"})
		require.NoError(b, err)

	}
}

func dbCleint(t *testing.T) *Mongo {
	uri := os.Getenv("MONGOURI")
	db, err := NewMgo(uri, time.Second*5)
	require.NoError(t, err)
	err = db.Ping()
	require.NoError(t, err)
	return db
}

func TestTableOperate(t *testing.T) {

	ctx := context.Background()
	db := dbCleint(t)
	tableName := "tmp_test_table_operate"

	exist, err := db.HasTable(ctx, tableName)
	require.NoError(t, err)
	if !exist {
		err := db.CreateTable(ctx, tableName)
		require.NoError(t, err)
	}
	exist, err = db.HasTable(ctx, tableName)
	require.NoError(t, err)
	if !exist {
		t.Errorf("table %s not exist", tableName)
		return
	}

	err = db.DropTable(ctx, tableName)
	require.NoError(t, err)

	exist, err = db.HasTable(ctx, tableName)
	require.NoError(t, err)
	if exist {
		t.Errorf("drop table %s, table already exist", tableName)
	}
}

func TestIndex(t *testing.T) {
	ctx := context.Background()
	tableName := "tmptest_index"

	db := dbCleint(t)
	table := db.Table(tableName)

	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	err = db.CreateTable(ctx, tableName)
	require.NoError(t, err)

	// 检查默认索引
	indexes, err := table.Indexes(ctx)
	require.NoError(t, err)
	if len(indexes) != 1 {
		t.Errorf("table %s index not equal one, indexes:%#v", tableName, indexes)
		return
	} else if indexes[0].Name != "_id_" {
		t.Errorf("table %s index name not equal _id_", tableName)
		return
	}

	// 创建索引
	createIndexes := map[string]dal.Index{
		"test_one": dal.Index{
			Name: "test_one",
			Keys: map[string]int32{"a": 1, "b": 1},
		},
		"test_backgroud": dal.Index{
			Name:       "test_backgroud",
			Keys:       map[string]int32{"aa": 1, "bb": -1},
			Background: true,
		},
		"test_unique": dal.Index{
			Name:   "test_unique",
			Keys:   map[string]int32{"aa": 1, "bb": 1},
			Unique: true,
		},
	}

	indexNameMap := make(map[string]string, 0)
	for indexName, index := range createIndexes {
		err := table.CreateIndex(ctx, index)
		require.NoError(t, err)
		indexNameMap[indexName] = indexName
	}
	// 检查新加的索引索引
	indexes, err = table.Indexes(ctx)
	require.NoError(t, err)

	for _, dbIndex := range indexes {
		delete(indexNameMap, dbIndex.Name)
		if index, ok := createIndexes[dbIndex.Name]; ok {
			require.Equal(t, index, dbIndex)
		}
	}
	if len(indexNameMap) > 0 {
		t.Errorf("table %s index %v not found", tableName, indexNameMap)
	}

	deleteIndexName := "test_unique"
	err = table.DropIndex(ctx, deleteIndexName)
	require.NoError(t, err)

	// 检查新加的索引索引
	indexes, err = table.Indexes(ctx)
	require.NoError(t, err)

	for _, dbIndex := range indexes {
		if dbIndex.Name == deleteIndexName {
			t.Errorf("table %s index name %s have drop. but already exist", tableName, deleteIndexName)
		}
	}

}

func TestInsertAndFind(t *testing.T) {
	ctx := context.Background()
	tableName := "tmptest_insert_find"

	db := dbCleint(t)

	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	resultOne := make(map[string]interface{}, 0)
	err = table.Find(map[string]string{"test_xxx_xxx": "1"}).One(ctx, &resultOne)
	if err != dal.ErrDocumentNotFound {
		require.NoError(t, err)
	}
	if len(resultOne) != 0 {
		t.Errorf("find one not data. but return data")
	}

	resultMany := make([]map[string]interface{}, 0)
	err = table.Find(map[string]string{"test_xxx_xxx": "1"}).All(ctx, &resultMany)
	require.NoError(t, err)
	if len(resultMany) != 0 {
		t.Errorf("find many not data. but return data")
	}

	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"a1": "a1",
		},
		map[string]interface{}{
			"a2": "a2",
		},
	}
	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)

	err = table.Insert(ctx, map[string]string{"b1": "b2"})
	require.NoError(t, err)

	// 查询单个数据
	resultOne = make(map[string]interface{}, 0)
	err = table.Find(map[string]string{"a1": "a1"}).One(ctx, &resultOne)
	require.NoError(t, err)

	if len(resultOne) == 0 {
		t.Errorf("find one not data.")
	}

	// 查询多个数据返回一条
	resultMany = make([]map[string]interface{}, 0)
	err = table.Find(map[string]string{"a2": "a1"}).All(ctx, &resultMany)
	require.NoError(t, err)
	if len(resultOne) == 0 {
		t.Errorf("find many not data.")
	}

	// 查询多个数据返回多条
	resultMany = make([]map[string]interface{}, 0)
	err = table.Find(nil).All(ctx, &resultMany)
	require.NoError(t, err)
	if len(resultOne) == 0 {
		t.Errorf("find many not data.")
	}
}

func TestFindOpt(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_find_option"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"a1":   "a1",
			"ext":  "ext",
			"sort": "1",
		},
		map[string]interface{}{
			"a2":   "a2",
			"ext":  "ext",
			"sort": "2",
		},
		map[string]interface{}{
			"a3":   "a4",
			"ext":  "ext",
			"sort": "3", // 排序会校验这个值，不要修改及在后面新加
		},
	}
	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)

	filter := map[string]string{"ext": "ext"}
	resultMany := make([]map[string]string, 0)
	err = table.Find(filter).All(ctx, &resultMany)
	require.NoError(t, err)
	if len(resultMany) != len(insertDataMany) {
		t.Errorf("find db data. row error")
		return
	}

	filter = map[string]string{"ext": "ext"}
	resultMany = make([]map[string]string, 0)
	err = table.Find(filter).Start(1).All(ctx, &resultMany)
	require.NoError(t, err)
	if len(resultMany) != (len(insertDataMany) - 1) {
		t.Errorf("find db skip one data. row error")
		return
	}

	filter = map[string]string{"ext": "ext"}
	resultMany = make([]map[string]string, 0)
	err = table.Find(filter).Start(uint64(len(insertDataMany))).All(ctx, &resultMany)
	require.NoError(t, err)
	if len(resultMany) != 0 {
		t.Errorf("find db skip %d data. row error", len(resultMany))
		return
	}

	filter = map[string]string{"ext": "ext"}
	resultMany = make([]map[string]string, 0)
	err = table.Find(filter).Limit(1).All(ctx, &resultMany)
	require.NoError(t, err)
	if len(resultMany) != 1 {
		t.Errorf("find db limit one data. row error")
		return
	}

	filter = map[string]string{"ext": "ext"}
	resultMany = make([]map[string]string, 0)
	err = table.Find(filter).Fields("sort").Sort("sort:1").All(ctx, &resultMany)
	require.NoError(t, err)
	for idx, row := range resultMany {
		// 是否只有一个字段
		if len(row) != 1 {
			t.Errorf("db find field, no effect")
			return
		}
		val, ok := row["sort"]
		if !ok {
			t.Errorf("db find field, fields not found")
			return
		}
		if idx == 0 && val != "1" {
			t.Errorf("db find sort, no effect")
			return
		}
	}

	filter = map[string]string{"ext": "ext"}
	resultMany = make([]map[string]string, 0)
	err = table.Find(filter).Fields("sort", "ext").Sort("sort:-1").All(ctx, &resultMany)
	require.NoError(t, err)
	for idx, row := range resultMany {
		if len(row) != 2 {
			t.Errorf("db find field, no effect")
			return
		}
		val, ok := row["sort"]
		if !ok {
			t.Errorf("db find field, fields not found")
			return
		}
		if idx == 0 && val != "3" {
			t.Errorf("db find sort, no effect")
			return
		}
	}

}

func TestFindOneOpt(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_find_option"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	insertDataMany := []map[string]string{
		map[string]string{
			"a1":   "a1",
			"ext":  "ext",
			"sort": "1",
		},
		map[string]string{
			"a2":   "a2",
			"ext":  "ext",
			"sort": "2",
		},
		map[string]string{
			"a3":   "a4",
			"ext":  "ext",
			"sort": "3", // 排序会校验这个值，不要修改及在后面新加
		},
	}
	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)

	filter := map[string]string{"ext": "ext"}
	resultOne := make(map[string]string, 0)
	err = table.Find(filter).One(ctx, &resultOne)
	require.NoError(t, err)

	filter = map[string]string{"ext": "ext"}
	resultOne = make(map[string]string, 0)
	err = table.Find(filter).Start(1).One(ctx, &resultOne)
	require.NoError(t, err)
	require.Equal(t, insertDataMany[1], resultOne)

	filter = map[string]string{"ext": "ext"}
	resultOne = make(map[string]string, 0)
	err = table.Find(filter).Start(uint64(len(insertDataMany))).One(ctx, &resultOne)
	if err != dal.ErrDocumentNotFound {
		require.NoError(t, err)
	}
	if len(resultOne) > 0 {
		t.Errorf("find one skip %d, not effect, data:%s", len(insertDataMany), resultOne)
		return
	}

	filter = map[string]string{"ext": "ext"}
	resultOne = make(map[string]string, 0)
	err = table.Find(filter).Fields("sort").Sort("sort:1").One(ctx, &resultOne)
	require.NoError(t, err)
	if len(resultOne) != 1 {
		t.Errorf("db find field, no effect")
		return
	}
	val, ok := resultOne["sort"]
	if !ok {
		t.Errorf("db find field, fields not found")
		return
	}
	if val != "1" {
		t.Errorf("db find sort, no effect")
		return
	}

	filter = map[string]string{"ext": "ext"}
	resultOne = make(map[string]string, 0)
	err = table.Find(filter).Fields("sort", "ext").Sort("sort:-1").One(ctx, &resultOne)
	require.NoError(t, err)
	if len(resultOne) != 2 {
		t.Errorf("db find field, no effect")
		return
	}
	val, ok = resultOne["sort"]
	if !ok {
		t.Errorf("db find field, fields not found")
		return
	}
	if val != "3" {
		t.Errorf("db find sort, no effect")
		return
	}

}

func TestCount(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_find_count"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)
	table := db.Table(tableName)

	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"a1":   "a1",
			"ext":  "ext",
			"sort": "1",
		},
		map[string]interface{}{
			"a2":   "a2",
			"ext":  "ext",
			"sort": "2",
		},
		map[string]interface{}{
			"a3":   "a4",
			"ext":  "ext",
			"sort": "3", // 排序会校验这个值，不要修改及在后面新加
		},
	}
	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)

	cnt, err := table.Find(nil).Count(ctx)
	require.NoError(t, err)
	if cnt != uint64(len(insertDataMany)) {
		t.Errorf("db count result error. not equal %d", cnt)
		return
	}
	filter := map[string]string{"ext": "ext"}
	cnt, err = table.Find(filter).Count(ctx)
	require.NoError(t, err)
	if cnt != uint64(len(insertDataMany)) {
		t.Errorf("db count result error. not equal %d", cnt)
		return
	}

	filter = map[string]string{"ext": "ext", "a1": "a1"}
	cnt, err = table.Find(filter).Count(ctx)
	require.NoError(t, err)
	if cnt != 1 {
		t.Errorf("db count result error. not equal 1")
		return
	}

	filter = map[string]string{"ext": "ext", "a1": "a1", "not foudnd row": "xxx_.xxx"}
	cnt, err = table.Find(filter).Count(ctx)
	require.NoError(t, err)
	if cnt != 0 {
		t.Errorf("db count result error. not equal 0")
		return
	}

}

func TestUpdate(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_update"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"a1":             "a1",
			"ext":            "ext",
			"sort":           "1",
			"change_version": "1",
		},
		map[string]interface{}{
			"a2":   "a2",
			"ext":  "ext",
			"sort": "2",
		},
	}
	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)

	filter := map[string]string{"ext": "ext"}
	update := map[string]string{"change_version": "2"}
	err = table.Update(ctx, filter, update)
	require.NoError(t, err)

	filter = map[string]string{"change_version": "2"}
	cnt, err := table.Find(filter).Count(ctx)
	require.NoError(t, err)
	if cnt != 2 {
		t.Errorf("update error.")
		return
	}

	update = map[string]string{"change_version": "2"}
	err = table.Update(ctx, nil, update)
	require.NoError(t, err)

}

func TestUpdateMulti(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_update_multi"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	type RowStruct struct {
		A     string  `bson:"a"`
		Ext   string  `bson:"ext"`
		Sort  string  `bson:"sort"`
		Inc   int64   `bson:"inc"`
		Unset *string `bson:"unset"`
	}
	unsetVal := "test_val"
	insertData := RowStruct{
		A:     "a",
		Ext:   "ext",
		Sort:  "2",
		Inc:   1,
		Unset: &unsetVal,
	}
	err = table.Insert(ctx, insertData)
	require.NoError(t, err)

	resultData := RowStruct{
		A:     "a_update_multi_model",
		Ext:   "ext",
		Sort:  "2",
		Inc:   2,
		Unset: nil,
	}

	filter := map[string]string{"ext": "ext"}
	update := []dal.ModeUpdate{
		dal.ModeUpdate{Op: "set", Doc: map[string]string{"a": "a_update_multi_model"}},
		dal.ModeUpdate{Op: "unset", Doc: map[string]string{"unset": ""}},
		dal.ModeUpdate{Op: "inc", Doc: map[string]interface{}{"inc": 1}},
	}
	err = table.UpdateMultiModel(ctx, filter, update...)
	require.NoError(t, err)

	resultOne := RowStruct{}
	err = table.Find(nil).One(ctx, &resultOne)
	require.NoError(t, err)
	require.Equal(t, resultData, resultOne)

}

func TestUpsert(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_upsert"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)
	// 更新或者新加接口
	notDataFilter := map[string]string{
		"not_data_test_upsert": "__cc_cc_upsert",
	}
	upsertAddData := map[string]string{
		"not_data_test_upsert_add_data": "not_data_test_upsert",
	}
	// add
	err = table.Upsert(ctx, notDataFilter, upsertAddData)
	require.NoError(t, err)

	resultOne := make(map[string]string, 0)
	err = table.Find(upsertAddData).One(ctx, &resultOne)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"not_data_test_upsert":          "__cc_cc_upsert",
		"not_data_test_upsert_add_data": "not_data_test_upsert",
	}, resultOne)

	dataFilter := map[string]string{
		"not_data_test_upsert_add_data": "not_data_test_upsert",
	}
	upsertUpdateData := map[string]string{
		"not_data_test_upsert_update_data": "update_data",
	}
	// change
	err = table.Upsert(ctx, dataFilter, upsertUpdateData)
	require.NoError(t, err)

	resultOneUpsert := make(map[string]string, 0)
	err = table.Find(dataFilter).One(ctx, &resultOneUpsert)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"not_data_test_upsert":             "__cc_cc_upsert",
		"not_data_test_upsert_add_data":    "not_data_test_upsert",
		"not_data_test_upsert_update_data": "update_data",
	}, resultOneUpsert)

}

func TestDelete(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_delete"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"a1": "a1",
		},
		map[string]interface{}{
			"a2": "a2",
		},
	}
	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)

	filter := map[string]string{"a1": "a1"}
	err = table.Delete(ctx, filter)
	require.NoError(t, err)

	resultMany := make([]map[string]string, 0)
	err = table.Find(filter).All(ctx, &resultMany)
	require.NoError(t, err)

	if len(resultMany) != 0 {
		t.Errorf("delete db row error. ")
		return
	}
}

func TestColumn(t *testing.T) {
	ctx := context.Background()
	tableName := "tmptest_Column"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)
	require.NoError(t, err)

	table := db.Table(tableName)

	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"a1":         "a1",
			"delete_col": "delete_col",
			"rename_col": "rename_col1",
			"add_col":    "add_col_exist",
		},
		map[string]interface{}{
			"a2": "a2",
		},
	}
	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)

	err = table.AddColumn(ctx, "add_col", "add_col_test")
	require.NoError(t, err)

	cnt, err := table.Find(map[string]string{"add_col": "add_col_exist"}).Count(ctx)
	require.NoError(t, err)
	if cnt != 1 {
		t.Errorf("add_col exist column data change")
		return

	}

	cnt, err = table.Find(map[string]string{"add_col": "add_col_test"}).Count(ctx)
	require.NoError(t, err)
	if cnt != 1 {
		t.Errorf("add_col exist failure")
		return
	}

	err = table.RenameColumn(ctx, "rename_col", "rename_col_after")
	require.NoError(t, err)

	cnt, err = table.Find(map[string]string{"rename_col_after": "rename_col1"}).Count(ctx)
	require.NoError(t, err)
	if cnt != 1 {
		t.Errorf("RenameColumn rename colume not found")
		return
	}
	cnt, err = table.Find(map[string]string{"rename_col": "rename_col1"}).Count(ctx)
	require.NoError(t, err)
	if cnt != 0 {
		t.Errorf("RenameColumn name  column already exist")
		return
	}

	err = table.DropColumn(ctx, "delete_col")
	require.NoError(t, err)

	cnt, err = table.Find(map[string]string{"delete_col": "delete_col"}).Count(ctx)
	require.NoError(t, err)
	if cnt != 0 {
		t.Errorf("DropColumn error. name  column already exist")
		return
	}
}

func TestAggregate(t *testing.T) {
	ctx := context.Background()
	tableName := "tmptest_Aggregate"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	err = db.CreateTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	err = table.Insert(ctx, map[string]string{"aa": "aa"})
	require.NoError(t, err)

	aggregateCond := []interface{}{
		map[string]interface{}{
			"$group": map[string]interface{}{
				"_id": "$aa",
				"num": map[string]interface{}{"$sum": 1},
			},
		},
	}
	type aggregateRowStruct struct {
		ID  string `bson:"_id"`
		Num int64  `bson:"num"`
	}
	resultOne := &aggregateRowStruct{}
	err = table.AggregateOne(ctx, aggregateCond, resultOne)
	require.NoError(t, err)
	require.Equal(t, aggregateRowStruct{
		ID:  "aa",
		Num: 1,
	}, *resultOne)

	resultAll := make([]aggregateRowStruct, 0)
	err = table.AggregateAll(ctx, aggregateCond, &resultAll)
	require.NoError(t, err)
	if len(resultAll) == 0 {
		t.Errorf("AggregateOne error")
		return
	}
	require.Equal(t, aggregateRowStruct{
		ID:  "aa",
		Num: 1,
	}, resultAll[0])

}

func TestUpdateModifyCount(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_update_modify_count"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"a1":             "a1",
			"ext":            "ext",
			"sort":           "1",
			"change_version": "1",
		},
		map[string]interface{}{
			"a2":   "a2",
			"ext":  "ext",
			"sort": "2",
		},
	}
	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)

	// update one
	filter := map[string]string{"ext": "ext"}
	update := map[string]string{"change_version": "2"}
	var modifyCount int64 = 0
	modifyCount, err = table.UpdateModifyCount(ctx, filter, update)
	require.NoError(t, err)
	require.NotEqual(t, 1, modifyCount)

	filter = map[string]string{"change_version": "2"}
	cnt, err := table.Find(filter).Count(ctx)
	require.NoError(t, err)
	if cnt != 2 {
		t.Errorf("update error.")
		return
	}

	//  update may
	filterMay := map[string]string{}
	update = map[string]string{"change_modify_count_many": "4"}
	modifyCount, err = table.UpdateModifyCount(ctx, filterMay, update)
	require.NoError(t, err)
	require.NotEqual(t, 0, modifyCount)
	filter = map[string]string{"change_modify_count_many": "4"}
	cnt, err = table.Find(filter).Count(ctx)
	require.NoError(t, err)
	if cnt != 2 {
		t.Errorf("update error.")
		return
	}

	// not row update
	filterNotFound := map[string]string{"ext_not_found": "ext"}
	update = map[string]string{"change_modify_count_not_found": "4"}
	modifyCount, err = table.UpdateModifyCount(ctx, filterNotFound, update)
	require.NoError(t, err)
	require.NotEqual(t, 0, modifyCount)
	filter = map[string]string{"change_modify_count_not_found": "4"}
	cnt, err = table.Find(filter).Count(ctx)
	require.NoError(t, err)
	if cnt != 0 {
		t.Errorf("update error.")
		return
	}
}

func TestConvInterface(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_decode_interface"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	type SubStruct struct {
		Int int    `bson:"sub_int" json:"sub_int"`
		Str string `bson:"sub_aa" json:"sub_aa"`
	}

	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"str":         "str",
			"bool":        true,
			"int":         int(1),
			"int8":        int8(8),
			"int16":       int16(8),
			"int32":       int32(8),
			"int64":       int64(8),
			"uint":        uint(8),
			"uint8":       uint8(8),
			"uint16":      uint16(8),
			"uint32":      uint32(8),
			"uint64":      uint64(8),
			"float32":     float32(8),
			"float64":     float64(8),
			"str_ptr":     "str",
			"bool_ptr":    true,
			"int_ptr":     int(1),
			"int8_ptr":    int8(8),
			"int16_ptr":   int16(8),
			"int32_ptr":   int32(8),
			"int64_ptr":   int64(8),
			"uint_ptr":    uint(8),
			"uint8_ptr":   uint8(8),
			"uint16_ptr":  uint16(8),
			"uint32_ptr":  uint32(8),
			"uint64_ptr":  uint64(8),
			"float32_ptr": float32(8),
			"float64_ptr": float64(8),
			"str_arr":     []string{"1", "2"},
			"int_arr":     []int{1, 2},

			"struct": map[string]interface{}{
				"sub_int": int(1),
				"sub_aa":  "struct",
			},
			"struct_ptr": map[string]interface{}{
				"sub_int": int(11),
				"sub_aa":  "ptr",
			},
			"tag_test": "11",
			//"sub_int":  (8888888),
			//"sub_aa":   "inline",
			"struct_arr": []SubStruct{

				SubStruct{
					Int: 1,
					Str: "sub_arr_str",
				},
			},
		},
		map[string]interface{}{
			"str":         "str",
			"bool":        true,
			"int":         int(12),
			"int8":        int8(82),
			"int16":       int16(82),
			"int32":       int32(82),
			"int64":       int64(82),
			"uint":        uint(82),
			"uint8":       uint8(82),
			"uint16":      uint16(82),
			"uint32":      uint32(82),
			"uint64":      uint64(82),
			"float32":     float32(82),
			"float64":     float64(82),
			"str_arr_ptr": []string{"1", "2"},
			"int_arr_ptr": []int{1, 2},

			"struct": map[string]interface{}{
				"sub_int": int(1),
				"sub_aa":  "struct",
			},
			"struct_ptr": map[string]interface{}{
				"sub_int": int(11),
				"sub_aa":  "struct ptr",
			},
			"tag_test": "11",
		},
	}

	// inline , inline ptr
	type resultStruct struct {
		Str        string      `bson:"str" json:"str"`
		Bool       bool        `bson:"bool" json:"bool"`
		Int        int         `bson:"int" json:"int"`
		Int8       int8        `bson:"int8" json:"int8"`
		Int16      int16       `bson:"int16" json:"int16"`
		Int32      int32       `bson:"int32" json:"int32"`
		Int64      int64       `bson:"int64" json:"int64"`
		Uint       uint        `bson:"uint" json:"uint"`
		Uint8      uint8       `bson:"uint8" json:"uint8"`
		Uint16     uint16      `bson:"uint16" json:"uint16"`
		Uint32     uint32      `bson:"uint32" json:"uint32"`
		Uint64     uint64      `bson:"uint64" json:"uint64"`
		Float32    float32     `bson:"float32" json:"float32"`
		Float64    float64     `bson:"float64" json:"float64"`
		StrPtr     *string     `bson:"str_ptr" json:"str_ptr"`
		BoolPtr    *bool       `bson:"bool_ptr" json:"bool_ptr"`
		IntPtr     *int        `bson:"int_ptr" json:"int_ptr"`
		Int8Ptr    *int8       `bson:"int8_ptr" json:"int8_ptr"`
		Int16Ptr   *int16      `bson:"int16_ptr" json:"int16_ptr"`
		Int32Ptr   *int32      `bson:"int32_ptr" json:"int32_ptr"`
		Int64Ptr   *int64      `bson:"int64_ptr" json:"int64_ptr"`
		UintPtr    *uint       `bson:"uint_ptr" json:"uint_ptr"`
		Uint8Ptr   *uint8      `bson:"uint8_ptr" json:"uint8_ptr"`
		Uint16Ptr  *uint16     `bson:"uint16_ptr" json:"uint16_ptr"`
		Uint32Ptr  *uint32     `bson:"uint32_ptr" json:"uint32_ptr"`
		Uint64Ptr  *uint64     `bson:"uint64_ptr" json:"uint64_ptr"`
		Float32Ptr *float32    `bson:"float32_ptr" json:"float32_ptr"`
		Float64Ptr *float64    `bson:"float64_ptr" json:"float64_ptr"`
		StrArr     []string    `bson:"str_arr" json:"str_arr"`
		IntArr     []int       `bson:"int_arr" json:"int_arr"`
		Struct     SubStruct   `bson:"struct" json:"struct"`
		TagTest    interface{} `bson:"tag_test" json:"tag_test"`
		StructPtr  *SubStruct  `bson:"struct_ptr" json:"struct_ptr"`
		//*SubStruct
		StructArr []SubStruct `bson:"struct_arr" json:"struct_arr"`
	}

	inserJSONByte, err := json.Marshal(insertDataMany)
	require.NoError(t, err)
	insertJSONUnmarshal := make([]resultStruct, 0)
	err = json.Unmarshal(inserJSONByte, &insertJSONUnmarshal)
	require.NoError(t, err)
	dbJSONUnmarshal := make([]resultStruct, 0)

	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)
	resultStructMany := make([]resultStruct, 0)
	err = table.Find(nil).Limit(2).All(ctx, &resultStructMany)
	require.NoError(t, err)
	require.Equal(t, insertDataMany[0]["bool"], resultStructMany[0].Bool)
	require.Equal(t, insertDataMany[0]["bool_ptr"], *resultStructMany[0].BoolPtr)
	require.Equal(t, insertDataMany[0]["int"], resultStructMany[0].Int)
	require.Equal(t, insertDataMany[0]["int_ptr"], *resultStructMany[0].IntPtr)
	require.Equal(t, insertDataMany[0]["int16"], resultStructMany[0].Int16)
	require.Equal(t, insertDataMany[0]["int16_ptr"], *resultStructMany[0].Int16Ptr)
	require.Equal(t, insertDataMany[0]["int32"], resultStructMany[0].Int32)
	require.Equal(t, insertDataMany[0]["int32_ptr"], *resultStructMany[0].Int32Ptr)
	require.Equal(t, insertDataMany[0]["int64"], resultStructMany[0].Int64)
	require.Equal(t, insertDataMany[0]["int64_ptr"], *resultStructMany[0].Int64Ptr)
	require.Equal(t, insertDataMany[0]["uint"], resultStructMany[0].Uint)
	require.Equal(t, insertDataMany[0]["uint_ptr"], *resultStructMany[0].UintPtr)
	require.Equal(t, insertDataMany[0]["uint32"], resultStructMany[0].Uint32)
	require.Equal(t, insertDataMany[0]["uint32_ptr"], *resultStructMany[0].Uint32Ptr)
	require.Equal(t, insertDataMany[0]["uint64"], resultStructMany[0].Uint64)
	require.Equal(t, insertDataMany[0]["uint64_ptr"], *resultStructMany[0].Uint64Ptr)
	require.Empty(t, resultStructMany[1].Uint64Ptr)

	require.Equal(t, insertDataMany[0]["struct"], map[string]interface{}{
		"sub_int": resultStructMany[0].Struct.Int,
		"sub_aa":  resultStructMany[0].Struct.Str,
	})
	require.Equal(t, insertDataMany[0]["struct_ptr"], map[string]interface{}{
		"sub_int": resultStructMany[0].StructPtr.Int,
		"sub_aa":  resultStructMany[0].StructPtr.Str,
	})
	require.Equal(t, insertDataMany[0]["str_arr"], resultStructMany[0].StrArr)
	require.Equal(t, insertDataMany[0]["int_arr"], resultStructMany[0].IntArr)
	require.Equal(t, insertDataMany[0]["struct_arr"], resultStructMany[0].StructArr)

	dbJSON, err := json.Marshal(resultStructMany)
	require.NoError(t, err)
	dbJSONUnmarshal = make([]resultStruct, 0)
	err = json.Unmarshal(dbJSON, &dbJSONUnmarshal)
	require.NoError(t, err)
	blog.ErrorJSON("%s  %s", insertJSONUnmarshal, resultStructMany)
	require.Equal(t, insertJSONUnmarshal, dbJSONUnmarshal)

	// test interface
	resultInterfaceMany := make([]interface{}, 0)
	err = table.Find(nil).Limit(2).All(ctx, &resultInterfaceMany)
	require.NoError(t, err)
	dbRow := resultInterfaceMany[0].(map[string]interface{})
	require.Equal(t, insertDataMany[0]["bool"], dbRow["bool"])
	require.Equal(t, insertDataMany[0]["bool_ptr"], dbRow["bool_ptr"])
	require.Equal(t, insertDataMany[0]["sub_aa"], dbRow["sub_aa"])

	dbJSON, err = json.Marshal(resultInterfaceMany)
	require.NoError(t, err)

	dbJSONUnmarshal = make([]resultStruct, 0)
	err = json.Unmarshal(dbJSON, &dbJSONUnmarshal)
	require.NoError(t, err)
	require.Equal(t, insertJSONUnmarshal, dbJSONUnmarshal)

	// test map
	resultMapMany := make([]map[string]interface{}, 0)
	err = table.Find(nil).Limit(2).All(ctx, &resultMapMany)
	require.NoError(t, err)

	dbJSON, err = json.Marshal(resultMapMany)
	require.NoError(t, err)

	dbJSONUnmarshal = make([]resultStruct, 0)
	err = json.Unmarshal(dbJSON, &dbJSONUnmarshal)
	require.NoError(t, err)
	require.Equal(t, insertJSONUnmarshal, dbJSONUnmarshal)
	/*
		err = table.Insert(ctx, insertDataMany)
		require.NoError(t, err)
		resultMany := make([]interface{}, 0)
		err = table.Find(nil).Limit(1).All(ctx, &resultMany)
		require.NoError(t, err)
		blog.InfoJSON("%s  \n %s", insertDataMany, resultMany)

		resultStruct := make([]tmp, 0)
		err = table.Find(nil).Limit(1).All(ctx, &resultStruct)
		require.NoError(t, err)
		blog.InfoJSON("%s  \n %s", insertDataMany, resultStruct)*/
}

func TestConvInterfaceMap(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_decode_interface_map"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	insertDataMany := []map[string]string{
		map[string]string{
			"str":  "str",
			"str2": "str2",
		},
		map[string]string{
			"str":  "str22",
			"str2": "str2222",
		},
	}

	type strMapStruct struct {
		Str  string `bson:"str"`
		Str2 string `bson:"str2"`
	}

	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)
	resultMapMany := make([]map[string]string, 0)
	err = table.Find(nil).Limit(2).All(ctx, &resultMapMany)
	require.NoError(t, err)
	require.Equal(t, insertDataMany, resultMapMany)

	resultStructMany := make([]strMapStruct, 0)
	err = table.Find(nil).Limit(2).All(ctx, &resultStructMany)
	require.NoError(t, err)
	require.Equal(t, insertDataMany[0], map[string]string{
		"str":  resultStructMany[0].Str,
		"str2": resultStructMany[0].Str2,
	})

	require.Equal(t, insertDataMany[1], map[string]string{
		"str":  resultStructMany[1].Str,
		"str2": resultStructMany[1].Str2,
	})

}

func TestConvInterfaceStructInline(t *testing.T) {

	ctx := context.Background()
	tableName := "tmptest_decode_interface_map"

	db := dbCleint(t)
	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)

	table := db.Table(tableName)

	insertDataMany := []map[string]string{
		map[string]string{
			"str":  "str",
			"str2": "str2",
		},
		map[string]string{
			"str":  "str22",
			"str2": "str2222",
		},
	}

	type StrMapStruct struct {
		Str  string `bson:"str"`
		Str2 string `bson:"str2"`
	}

	type resultStruct struct {
		StrMapStruct
	}

	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)
	resultStructMany := make([]resultStruct, 0)
	err = table.Find(nil).Limit(2).All(ctx, &resultStructMany)
	require.NoError(t, err)
	require.Equal(t, insertDataMany[0], map[string]string{
		"str":  resultStructMany[0].Str,
		"str2": resultStructMany[0].Str2,
	})

	require.Equal(t, insertDataMany[1], map[string]string{
		"str":  resultStructMany[1].Str,
		"str2": resultStructMany[1].Str2,
	})

}

func TestConvInterfaceStructTagInline(t *testing.T) {

	ctx := context.Background()
	db := dbCleint(t)

	tableName := "tmptest_decode_interface_struct_tag"

	// 清理数据
	err := db.DropTable(ctx, tableName)
	require.NoError(t, err)
	table := db.Table(tableName)

	insertDataMany := []map[string]string{
		map[string]string{
			"str":  "str",
			"str2": "str2",
		},
		map[string]string{
			"str":  "str22",
			"str2": "str2222",
		},
	}

	type StrMapStruct struct {
		Str  string `bson:"str"`
		Str2 string `bson:"str2"`
	}

	type resultStruct struct {
		StrMapStruct
	}

	err = table.Insert(ctx, insertDataMany)
	require.NoError(t, err)
	resultStructMany := make([]resultStruct, 0)
	err = table.Find(nil).Limit(2).All(ctx, &resultStructMany)
	require.NoError(t, err)
	require.Equal(t, insertDataMany[0], map[string]string{
		"str":  resultStructMany[0].Str,
		"str2": resultStructMany[0].Str2,
	})

	require.Equal(t, insertDataMany[1], map[string]string{
		"str":  resultStructMany[1].Str,
		"str2": resultStructMany[1].Str2,
	})

}

func TestNextSequence(t *testing.T) {

	ctx := context.Background()

	db := dbCleint(t)

	// 清理数据
	err := db.DropTable(ctx, "cc_idgenerator")
	require.NoError(t, err)

	id, err := db.NextSequence(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, uint64(1), id)
	id, err = db.NextSequence(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, uint64(2), id)
}
