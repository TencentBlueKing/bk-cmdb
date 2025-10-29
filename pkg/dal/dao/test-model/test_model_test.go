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

package testmodel

import (
	"context"
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	idgenerator "github.com/TencentBlueKing/bk-cmdb/pkg/dal/dao/id-generator"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/orm"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/table"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/types"
	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
	"github.com/TencentBlueKing/bk-cmdb/pkg/tests"
)

var predefinedModels = []table.TestModel{
	{
		Name:    "a",
		Size:    1,
		Weight:  2.0,
		Int64s:  []int64{1, 2, 3},
		Strings: []string{"a", "b", "c"},
	},
	{
		Name:    "b",
		Size:    3,
		Weight:  4.0,
		Int64s:  []int64{5, 6, 7},
		Strings: []string{"d", "e", "f"},
	},
}

func TestGenericDaoCRUD(t *testing.T) {
	testDO, err := prepareTestModelDo(t)
	if err != nil {
		t.Errorf("fail to prepare test model do, err: %v", err)
		return
	}
	kt := tests.GetKit(t)
	ids, err := testDO.BatchCreate(kt, predefinedModels)
	if err != nil {
		t.Errorf("fail to create err: %v", err)
		return
	}
	if len(ids) != len(predefinedModels) {
		t.Errorf("generated id length mismatch")
		return
	}
	t.Logf("success to create, ids: %v", ids)

	t.Run("get name by id", func(t *testing.T) {
		kt := tests.GetKit(t)
		name, err := testDO.GetNameByID(kt, ids[1])
		assert.Nil(t, err, "fail to get name by id")
		assert.Equal(t, name, predefinedModels[1].Name)
		t.Logf("success to get name by id: %s", name)
	})

	listOpt := &types.ListOption{
		Fields: nil,
		Filter: filter.AllExpression(),
		Page:   types.NewDefaultPage(),
	}

	t.Run("list all", func(t *testing.T) {
		kt := tests.GetKit(t)
		result, err := testDO.List(kt, listOpt)
		assert.Nil(t, err, "fail to list")
		assert.Equal(t, result.Details, predefinedModels)
	})

	t.Run("list id only", func(t *testing.T) {
		// list with fields
		listOpt.Fields = []string{"id"}
		kt := tests.GetKit(t)
		result, err := testDO.List(kt, listOpt)
		assert.Nil(t, err, "fail to list with fields")
		idOnlyModels := lo.Map(predefinedModels, func(item table.TestModel, index int) table.TestModel {
			return table.TestModel{Base: table.Base{ID: item.ID}}
		})
		assert.Equal(t, result.Details, idOnlyModels)
	})
}

func TestGenericDaoTxn(t *testing.T) {
	testDO, err := prepareTestModelDo(t)
	if err != nil {
		t.Errorf("fail to prepare test model do, err: %v", err)
		return
	}

	txnCreateData := []table.TestModel{
		{Name: "txn-1"}, {Name: "txn-2"},
	}
	kt := tests.GetKit(t)
	err = testDO.AutoTxn(kt, func(txn orm.Interface) error {

		tx := testDO.WithTx(txn)
		newIDs, err := tx.BatchCreate(kt, txnCreateData)
		if err != nil {
			t.Logf("fail to create, ids: %v", newIDs)
			return err
		}
		listOpt := &types.ListOption{
			Fields: nil,
			Filter: filter.ExpressionAnd(filter.RuleIn("id", lo.ToAnySlice(newIDs))),
			Page:   types.NewCountPage(),
		}
		ret, err := tx.List(kt, listOpt)
		if err != nil {
			t.Logf("fail to count, ids: %v", newIDs)
			return err
		}
		if ret.Count != uint64(len(newIDs)) {
			t.Logf("count mismatch, ids: %v", newIDs)
			return fmt.Errorf("count mismatch")
		}
		return abortTransactionErr
	})

	assert.ErrorIs(t, err, abortTransactionErr)

	// check if the data is still there
	listOpt := &types.ListOption{
		Fields: nil,
		Filter: filter.ExpressionAnd(filter.RuleIn("name", []any{"txn-1"})),
		Page:   types.NewCountPage(),
	}
	ret, err := testDO.List(kt, listOpt)
	assert.Nil(t, err, "fail to count outside transaction")
	assert.Equal(t, ret.Count, uint64(0), "data still exists outside transaction")
}

func prepareTestModelDo(t *testing.T) (Interface, error) {
	db, err := tests.GetRealDB(t)
	ctx := context.Background()
	if err != nil {
		t.Errorf("fail to get test db, err: %v", err)
		return nil, err
	}
	ormInst, err := orm.New(ctx, db)
	if err != nil {
		t.Errorf("fail to init orm for test, err: %v", err)
		return nil, err
	}
	idGen := idgenerator.New(db)
	if db.Migrator().HasTable(table.TestModel{}) {
		err := db.Migrator().DropTable(table.TestModel{})
		if err != nil {
			t.Errorf("fail to drop table err: %v", err)
			return nil, err
		}
	}
	err = db.Migrator().AutoMigrate(table.TestModel{})
	if err != nil {
		t.Errorf("fail to auto migrate err: %v", err)
		return nil, err
	}

	testDO := NewDao(ormInst, idGen)
	return testDO, nil
}

var abortTransactionErr = fmt.Errorf("abort transaction for test")
