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

package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

func TestMockInsert(t *testing.T) {
	var err error
	db := NewMock()
	tablename := "test"
	ctx := context.Background()

	var errmsg = "this is an error"
	var mockErr = fmt.Errorf(errmsg)

	err = db.Mock(MockResult{Err: mockErr}).Table(tablename).Insert(ctx, map[string]interface{}{"name": "name"})
	require.NoError(t, err)
	err = db.Table(tablename).Insert(ctx, map[string]interface{}{"name": "name"})
	require.EqualError(t, err, errmsg)
}

func TestMockUpdate(t *testing.T) {
	var err error
	db := NewMock()
	tablename := "test"
	ctx := context.Background()

	var errmsg = "this is an error"
	var mockErr = fmt.Errorf(errmsg)

	err = db.Mock(MockResult{Err: mockErr}).Table(tablename).Update(ctx, map[string]interface{}{"name": "name"}, map[string]interface{}{"name": "name"})
	require.NoError(t, err)
	err = db.Table(tablename).Update(ctx, map[string]interface{}{"name": "name"}, map[string]interface{}{"name": "name"})
	require.EqualError(t, err, errmsg)
}

func TestMockDelete(t *testing.T) {
	var err error
	db := NewMock()
	tablename := "test"
	ctx := context.Background()

	var errmsg = "this is an error"
	var mockErr = fmt.Errorf(errmsg)

	err = db.Mock(MockResult{Err: mockErr}).Table(tablename).Delete(ctx, map[string]interface{}{"name": "name"})
	require.NoError(t, err)
	err = db.Table(tablename).Delete(ctx, map[string]interface{}{"name": "name"})
	require.EqualError(t, err, errmsg)
}

func TestMockAll(t *testing.T) {
	var err error
	db := NewMock()
	tablename := "test"
	ctx := context.Background()

	var errmsg = "this is an error"
	var mockErr = fmt.Errorf(errmsg)
	var mockResult = []map[string]interface{}{
		{
			"name": "name",
		},
	}
	var actualResual = []map[string]interface{}{}

	err = db.Mock(MockResult{Err: mockErr}).Table(tablename).Find(map[string]interface{}{"name": "name"}).Fields().Limit(0).Start(0).All(ctx, &mockResult)
	require.NoError(t, err)
	err = db.Table(tablename).Find(map[string]interface{}{"name": "name"}).Fields().Limit(0).Start(0).All(ctx, &actualResual)
	require.EqualError(t, err, errmsg)

	mockout, err := bson.Marshal(mockResult)
	require.NoError(t, err)
	actualout, err := bson.Marshal(actualResual)
	require.NoError(t, err)

	require.Equal(t, mockout, actualout)
}
func TestMockOne(t *testing.T) {
	var err error
	db := NewMock()
	tablename := "test"
	ctx := context.Background()

	var errmsg = "this is an error"
	var mockErr = fmt.Errorf(errmsg)
	var mockResult = map[string]interface{}{
		"name": "name",
	}
	var actualResual = map[string]interface{}{}

	err = db.Mock(MockResult{Err: mockErr}).Table(tablename).Find(map[string]interface{}{"name": "name"}).Fields().Limit(0).Start(0).One(ctx, &mockResult)
	require.NoError(t, err)
	err = db.Table(tablename).Find(map[string]interface{}{"name": "name"}).Fields().Limit(0).Start(0).One(ctx, &actualResual)
	require.EqualError(t, err, errmsg)

	mockout, err := bson.Marshal(mockResult)
	require.NoError(t, err)
	actualout, err := bson.Marshal(actualResual)
	require.NoError(t, err)

	require.Equal(t, mockout, actualout)
}
