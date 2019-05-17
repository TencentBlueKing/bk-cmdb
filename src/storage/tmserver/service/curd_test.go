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

package service_test

import (
	"context"
	"testing"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/storage/dal/mongo/remote"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestCURD(t *testing.T) {

	startCoreService(t, "127.0.0.1", 8899)

	dbDal, err := remote.New("127.0.0.1:8899")
	require.NoError(t, err)
	require.NotNil(t, dbDal)
	doc := mapstr.MapStr{
		"field": "value_" + xid.New().String(),
	}

	// insert
	col := dbDal.Table("cc_rpcCRUD")
	err = col.Insert(context.TODO(), doc)
	require.NoError(t, err)

	// find
	cond := mongo.NewCondition()
	cond.Element(&mongo.Regex{Key: "field", Val: "value_"})
	result := mapstr.NewArray()
	err = col.Find(cond.ToMapStr()).All(context.TODO(), &result)
	require.NoError(t, err)
	t.Logf("find results:%v", result)

	// update
	err = col.Update(context.TODO(), cond.ToMapStr(), mapstr.MapStr{
		"field": "value_update_" + xid.New().String(),
	})
	require.NoError(t, err)
	// count
	cnt, err := col.Find(cond.ToMapStr()).Count(context.TODO())
	require.NoError(t, err)
	require.NotEqual(t, uint64(0), cnt)
	t.Logf("count:%d", cnt)
	// delete
	err = col.Delete(context.TODO(), cond.ToMapStr())
	require.NoError(t, err)

	// count
	cnt, err = col.Find(cond.ToMapStr()).Count(context.TODO())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)
	t.Logf("count:%d", cnt)
}
