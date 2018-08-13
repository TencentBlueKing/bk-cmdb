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

package mongo_test

import (
	"configcenter/src/storage/dal"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"configcenter/src/storage/dal/mongo"
)

func TestRPCClient(t *testing.T) {
	cli, err := mongo.NewRPC("127.0.0.1:50010")
	require.NoError(t, err)

	ctx := context.Background()
	tablename := "testtable"

	// inset one
	err = cli.Table(tablename).Insert(ctx, map[string]interface{}{
		"name": "name1",
	})
	require.NoError(t, err, "insert one")

	// insert multi
	err = cli.Table(tablename).Insert(ctx, []map[string]interface{}{
		{
			"name": "name2",
		},
		{
			"name": "name3",
		},
	})
	require.NoError(t, err, "insert multi")

	// find all
	findall := []map[string]interface{}{}
	err = cli.Table(tablename).Find(nil).All(ctx, &findall)
	require.NoError(t, err, "find all")
	require.True(t, len(findall) > 0)

	// update
	err = cli.Table(tablename).Update(ctx, map[string]interface{}{"name": "name1"}, map[string]interface{}{"name": "name4"})
	require.NoError(t, err, "update")

	// find one
	findone := map[string]interface{}{}
	err = cli.Table(tablename).Find(map[string]interface{}{"name": "name4"}).One(ctx, &findone)
	require.NoError(t, err)
	require.True(t, findone["name"] == "name4")

	// delete filter
	err = cli.Table(tablename).Delete(ctx, map[string]interface{}{"name": "name4"})
	require.NoError(t, err)
	findone = map[string]interface{}{}
	err = cli.Table(tablename).Find(map[string]interface{}{"name": "name4"}).One(ctx, &findone)
	require.EqualError(t, err, dal.ErrDocumentNotFound.Error())
	require.True(t, findone["name"] == nil)

	// delete all
	err = cli.Table(tablename).Delete(ctx, nil)
	require.NoError(t, err)
	findall = []map[string]interface{}{}
	err = cli.Table(tablename).Find(nil).All(ctx, &findall)
	require.NoError(t, err)
	require.True(t, len(findall) <= 0)
}
