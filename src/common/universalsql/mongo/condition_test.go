/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mongo_test

import (
	"configcenter/src/common/universalsql/mongo"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConditionFromMapStr(t *testing.T) {

	target := mongo.NewCondition()
	target.Element(&mongo.Eq{Key: "testelementeq", Val: "testeqval"})
	target.And(&mongo.Lt{Key: "testandlt", Val: "testandltval"})
	target.Or(&mongo.Lt{Key: "testorlt", Val: "testorltval"})
	target.Element(&mongo.In{Key: "testelementin", Val: []string{"testelementin"}})
	_, embed := target.Embed("testembedname")
	embed.Or(&mongo.Gt{Key: "testembedgt", Val: "testembedgtval"})
	embed.And(&mongo.Gt{Key: "testembedgt", Val: "testembedgtval"})
	embed.Element(&mongo.Lt{Key: "testembedeq", Val: "testembedeqval"})
	embed.Element(&mongo.Lt{Key: "testembedeq2", Val: "testembedeqval2"})
	_, subembed := embed.Embed("subembed")
	subembed.Element(&mongo.Eq{Key: "subembedkey", Val: "subembedkeyval"})

	sql, _ := target.ToSQL()
	t.Logf("target sql:%s", sql)

	recoverSql, err := mongo.NewConditionFromMapStr(target.ToMapStr())
	require.NoError(t, err)
	sql, _ = recoverSql.ToSQL()
	t.Logf("recover sql:%s", sql)
}
