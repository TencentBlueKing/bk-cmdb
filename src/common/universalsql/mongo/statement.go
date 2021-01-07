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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/universalsql"
	"encoding/json"
)

//var _ universalsql.WhereStatement = (*statement)(nil)
var _ universalsql.Result = (*statementResult)(nil)

type statementResult struct {
	conds []universalsql.Condition
}

type statement struct {
}

func (d *statementResult) ToSQL() (string, error) {
	sql, err := json.Marshal(d.ToMapStr())
	return string(sql), err
}

func (d *statementResult) ToMapStr() mapstr.MapStr {

	condResult := mapstr.New()
	for _, cond := range d.conds {
		condResult.Merge(cond.ToMapStr())
	}
	return condResult
}

func (d *statement) Conditions(conds ...universalsql.Condition) universalsql.Result {

	result := &statementResult{}
	for _, cond := range conds {
		result.conds = append(result.conds, cond)
	}
	return result
}
