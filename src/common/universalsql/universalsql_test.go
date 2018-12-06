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
package universalsql_test

import (
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUniservalsqlStatements(t *testing.T) {

	mtable := mongo.New()

	resultSql, err := mtable.Create().Fields(universalsql.Field{Key: "test", Val: "helloworld"}).ToSQL()
	require.NoError(t, err)
	t.Logf("create statement sql:%v", resultSql)

	mcond := mongo.NewCondition()
	_, embed := mcond.Element(&mongo.Eq{Key: "test", Val: "testval"}).And(&mongo.Lt{Key: "testand", Val: "testandval"}).Embed("testembedname")
	embed.Or(&mongo.Gt{Key: "testgt", Val: "testgtval"})
	embed.Element(&mongo.Lt{Key: "testeq", Val: "testeqval"})
	resultSql, err = mtable.Where().Conditions(mcond).ToSQL()
	require.NoError(t, err)
	t.Logf("delete statement sql:%v", resultSql)

}
