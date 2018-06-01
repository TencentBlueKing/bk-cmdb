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

package mgoclient

import (
	"configcenter/src/storage"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUint(t *testing.T) {
	db, err := NewMgoCli("127.0.0.1", "27017", "user", "pwd", "", "cmdb")
	require.NoError(t, err)
	err = db.Open()
	require.NoError(t, err)

	tablename := "testcollection"
	require.NotNil(t, db.GetSession())
	require.Error(t, db.ExecSql(""))
	db.HasFields(tablename, "field")

	assert.Equal(t, storage.DI_MONGO, db.GetType())

	testTable(t, db)
	testSingle(t, db)
	testMulti(t, db)
	testIncID(t, db)
	db.DropTable(tablename)

	db.Close()
}

func TestStruct(t *testing.T) {
	db, err := NewMgoCli("127.0.0.1", "27017", "user", "pwd", "", "cmdb")
	require.NoError(t, err)
	err = db.Open()
	require.NoError(t, err)

	tablename := "testcollection2"
	// doc := collection{
	// 	Name:    "testname1",
	// 	Comment: "testcomment1",
	// }
	// db.Insert(tablename, doc)
	doc := collection{
		Comment: "testcomment2",
	}
	doc.Name = "testname2"
	db.Insert(tablename, doc)

	doc2 := collection{
		Comment: "testcomment5",
	}
	doc2.Name = "testname5"
	condiction := collection{}
	condiction.Name = "testname2"
	db.UpdateByCondition(tablename, doc2, condiction)

	result := collection{}
	db.GetOneByCondition(tablename, nil, doc2, &result)
	assert.Equal(t, doc2, result)

	count, err := db.GetCntByCondition(tablename, doc2)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	results := []collection{}
	db.GetMutilByCondition(tablename, nil, doc2, &results, "", 0, 0)
	fmt.Println(results)

	// db.DelByCondition(tablename, doc2)

	// db.DropTable(tablename)

}

type Namee struct {
	Name string `bson:"Name,omitempty"`
}
type collection struct {
	Namee   `bson:",inline"`
	Comment string `bson:"Comment,omitempty"`
}

func testIncID(t *testing.T, db *MgoCli) {
	tablename := "testcollection"

	id, err := db.GetIncID(tablename)
	require.NoError(t, err)
	require.Equal(t, int64(1), id)

	id, err = db.GetIncID(tablename)
	require.NoError(t, err)
	require.Equal(t, int64(2), id)

	err = db.DelByCondition("cc_idgenerator", map[string]interface{}{"_id": tablename})
	require.NoError(t, err)
}

func testSingle(t *testing.T, db *MgoCli) {
	tablename := "testcollection"
	originData := map[string]interface{}{"ID": int64(1), "Name": "name", "comment": "comment"}
	_, err := db.Insert(tablename, originData)
	require.NoError(t, err)

	condiction := map[string]interface{}{"ID": int64(1)}
	result := map[string]interface{}{}
	err = db.GetOneByCondition(tablename, []string{"ID", "Name", "comment"}, condiction, &result)
	require.NoError(t, err)
	require.Equal(t, originData, result)

	err = db.UpdateByCondition(tablename, map[string]interface{}{"ID": int64(1), "Name": "newname", "comment": "comment"}, condiction)
	require.NoError(t, err)
	result = map[string]interface{}{}
	err = db.GetOneByCondition(tablename, []string{"ID", "Name", "time"}, condiction, &result)
	require.NoError(t, err)
	require.Equal(t, "newname", result["Name"])

	require.Error(t, db.UpdateByCondition(tablename, nil, nil))

	err = db.DelByCondition(tablename, condiction)
	require.NoError(t, err)
}

func testMulti(t *testing.T, db *MgoCli) {
	tablename := "testcollection"

	originData := []interface{}{
		map[string]interface{}{"ID": 2, "Name": "name2", "comment": "comment00"},
		map[string]interface{}{"ID": 3, "Name": "name3", "comment": "comment00"},
	}
	err := db.InsertMuti(tablename, originData)
	require.NoError(t, err)

	condiction := map[string]interface{}{"comment": "comment00"}
	result := []map[string]interface{}{}
	err = db.GetMutilByCondition(tablename, []string{"ID", "Name", "comment"},
		condiction, &result, "", 0, 0)
	require.NoError(t, err)

	require.Len(t, result, 2)
	require.ElementsMatch(t, originData, result)

	count, err := db.GetCntByCondition(tablename, condiction)
	require.NoError(t, err)
	require.Equal(t, 2, count)

	err = db.DelByCondition(tablename, condiction)
	require.NoError(t, err)
}

func testTable(t *testing.T, db *MgoCli) {
	tablename := "testcollection"
	defer db.DropTable(tablename)
	require.NoError(t, db.CreateTable(tablename))

	hasTable, err := db.HasTable(tablename)
	require.NoError(t, err)
	require.True(t, hasTable)

	index1 := storage.GetMongoIndex("test", []string{"bbb"}, false, false)
	index2 := storage.GetMongoIndex("test1", []string{"t"}, false, false)
	err = db.Index(tablename, index1)
	require.NoError(t, err)
	err = db.Index(tablename, index2)
	require.NoError(t, err)

	column1 := storage.GetMongoColumn("t22", 1)
	column2 := storage.GetMongoColumn("t3", "1")
	column3 := storage.GetMongoColumn("t1", 11)

	err = db.AddColumn(tablename, column1)
	require.NoError(t, err)

	err = db.AddColumn(tablename, column2)
	require.NoError(t, err)

	err = db.AddColumn(tablename, column3)
	require.NoError(t, err)

	err = db.DropColumn(tablename, "test")
	require.NoError(t, err)
	err = db.ModifyColumn(tablename, "test1", "tt121")
	require.NoError(t, err)
}

func TestMongoTime(t *testing.T) {
	return
	db, err := NewMgoCli("127.0.0.1", "27017", "user", "pwd", "", "cmdb")
	err = db.Open()
	if nil != err {
		t.Errorf("%s", err)
	}
	data := make(map[string]interface{})
	data["Time"] = time.Now()
	startTime, err := time.Parse("2006-01-02 15:04:05", "2017-12-25 12:02:00")
	fmt.Println(startTime)
	db.Insert("test", data)
}

func TestSearchTime(t *testing.T) {
	return
	db, err := NewMgoCli("127.0.0.1", "27017", "user", "pwd", "", "cmdb")
	err = db.Open()
	if nil != err {
		t.Errorf("%s", err)
	}
	var result interface{}
	startTime, err := time.Parse("2006-01-02 15:04:05", "2017-12-25 12:50:51")
	startTime = time.Unix(startTime.Unix()-8*3600, 0)
	cond := make(map[string]interface{})
	fmt.Println(startTime)
	cond["$gt"] = startTime
	condition := make(map[string]interface{})
	condition["Time"] = cond
	db.GetOneByCondition("test", nil, condition, &result)
	fmt.Println(result)
}
