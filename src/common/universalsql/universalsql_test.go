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
	"testing"

	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/universalsql"
)

func TestUniservalsqlCondition(t *testing.T) {

	sql := universalsql.CreateMongoSQLHelper()


	// {"$and":[{"field_name":{"$eq":"field_name_value"}}]}
	sql.Where().And(&universalsql.Eq{Key: "field_name", Val: "field_name_val"}).ToSQL()


	// {"field_name":"$in":["field_name_val"],"$or":[{"field_name1":{"$neq":"field_name1_val"}}]}
	sql.Where().Element(&universalsql.In{Key: "field_name", Val: []string{"field_name_val"}}).Or(&universalsql.Neq{Key:"field_name1", Val:"field_name1_val"}).ToSQL()

}
