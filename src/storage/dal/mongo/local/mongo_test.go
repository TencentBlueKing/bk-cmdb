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
	"net/http"
	"os"
	"testing"
	"time"

	"configcenter/src/common"
	"configcenter/src/storage/dal"

	"github.com/stretchr/testify/require"
)

func BenchmarkLocalCUD(b *testing.B) {
	db, err := NewMgo("127.0.0.1:27010", time.Second*5)
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
	if cnt != 2 {
		t.Errorf("update error.")
		return
	}
}
