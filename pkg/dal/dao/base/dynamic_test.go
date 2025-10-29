/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package base

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	idgenerator "github.com/TencentBlueKing/bk-cmdb/pkg/dal/dao/id-generator"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/orm"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/table"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/types"
	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
	"github.com/TencentBlueKing/bk-cmdb/pkg/structs"
	"github.com/TencentBlueKing/bk-cmdb/pkg/tests"
)

const dynamicTableName = "test_dynamic"

func TestDynamic(t *testing.T) {
	db, err := tests.GetRealDB(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	kt := tests.GetKit(t)
	ormInst, err := orm.New(kt, db)
	if !assert.NoError(t, err, "init Orm error") {
		return
	}
	idgen := idgenerator.New(db)

	if err := autoDeleteTable(ormInst, dynamicTableName); err != nil {
		t.Fatalf("fail to delete table before test: %v", err)
		return
	}

	listTest0Opt := &types.ListOption{
		Fields: nil,
		Filter: filter.RuleEqual("name", "test-0"),
		Page:   types.NewDefaultPage(),
	}

	// 构造动态结构体
	testDynamicBuilder, err := structs.UpsertBuilderByFields(dynamicTableName, []structs.Field{
		{
			Name:      "Base",
			Type:      table.BaseModelName,
			IsSlice:   false,
			Tags:      map[string]string{"gorm": "embedded", "json": ",inline"},
			Anonymous: true,
			Validator: nil,
		},
		{
			Name:      "Name",
			Type:      structs.StringType,
			IsSlice:   false,
			Tags:      map[string]string{"gorm": "column:name;size:64;not null", "json": "name"},
			Anonymous: false,
			Validator: nil,
		},
	})

	if !assert.Nil(t, err, "upsert builder error") {
		return
	}
	// 建表
	err = ormInst.DBContext(kt).
		Table(dynamicTableName).
		Migrator().
		AutoMigrate(testDynamicBuilder.New().Value())
	if !assert.Nil(t, err, "auto migrate error") {
		return
	}

	t.Cleanup(func() {
		_ = autoDeleteTable(ormInst, dynamicTableName)
	})

	// 注册表结构
	err = table.RegisterBuilder(testDynamicBuilder)
	if !assert.Nil(t, err, "register builder to table error") {
		return
	}

	// 初始化动态DO
	dynamic := NewDynamicConstructor(ormInst, idgen)

	// 设置model
	dynamicWithModel, err := dynamic.Dynamic(dynamicTableName)
	if !assert.Nil(t, err, "with model error") {
		return
	}

	// 测试创建
	// 构造Slice数据
	sliceData := testDynamicBuilder.NewSlice(3, 3)
	for i := range sliceData.Cap() {
		st := testDynamicBuilder.New()
		err = st.Set("Name", fmt.Sprintf("test-%d", i))
		if !assert.Nilf(t, err, "set name error idx %d", i) {
			return
		}
		err = sliceData.SetStruct(i, st)
		if !assert.Nilf(t, err, "set struct error idx %d", i) {
			return
		}
	}
	ids, err := dynamicWithModel.BatchCreate(kt, sliceData)
	if !assert.Nil(t, err, "batch create error") {
		return
	}
	if !assert.Len(t, ids, sliceData.Len()) {
		return
	}

	// 测试查询
	listRet, err := dynamicWithModel.List(kt, listTest0Opt)
	if !assert.Nil(t, err, "list error") {
		return
	}
	if !assert.Equal(t, 1, listRet.Details.Len(), "list by id length mismatch") {
		return
	}
	m0, err := listRet.Details.GetStruct(0)
	assert.Nil(t, err, "get struct error at index 0")
	idVal, err := m0.Get("ID")
	assert.Nil(t, err, "get id error")
	if len(idVal.String()) == 0 {
		t.Errorf("id not set successful, got %s, want non-empty", idVal.String())
		return
	}

	// 测试更新
	dataTest0, err := listRet.Details.GetStruct(0)
	if !assert.Nil(t, err, "get struct error") {
		return
	}
	dataTest0.Set("Name", "test-0-updated")
	// update
	updated, err := dynamicWithModel.Update(kt, filter.RuleEqual("name", "test-0"), dataTest0)
	if !assert.Nil(t, err, "update error") {
		return
	}
	if !assert.Equal(t, updated, int(1), "update count mismatch") {
		return
	}

	// 测试删除
	deleted, err := dynamicWithModel.Delete(kt, filter.RuleEqual("name", "test-0-updated"))
	if !assert.Nil(t, err, "delete error") {
		return
	}
	if !assert.Equal(t, deleted, int(1), "delete count mismatch") {
		return
	}

	// 测试删除后查询
	listAllTestOpt := &types.ListOption{
		Fields: nil,
		Filter: filter.RuleCis("name", "test"),
		Page:   types.NewDefaultPage(),
	}
	listRet, err = dynamicWithModel.List(kt, listAllTestOpt)
	if !assert.Nil(t, err, "list error") {
		return
	}
	// should be 2 after delete 1
	if !assert.Equal(t, listRet.Details.Len(), 2, "list by id length mismatch") {
		return
	}
	names := make([]string, 0)
	for i := range listRet.Details.Len() {
		data, err := listRet.Details.GetStruct(i)
		if !assert.Nil(t, err, "get struct error") {
			return
		}
		name, err := data.Get("Name")
		if !assert.Nil(t, err, "get name error") {
			return
		}
		names = append(names, name.String())
	}
	if !assert.ElementsMatch(t, names, []string{"test-1", "test-2"}, "names mismatch") {
		return
	}

}

func autoDeleteTable(ormInst orm.Interface, table string) error {
	var dummy struct{}
	if ormInst.DB().Table(table).Migrator().HasTable(dummy) {
		err := ormInst.DB().Table(table).
			Migrator().
			DropTable(dummy)
		if err != nil {
			return fmt.Errorf("fail to delete table %s: %v", table, err)
		}
	}

	return nil
}
