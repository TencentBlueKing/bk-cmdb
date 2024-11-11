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

// Package util TODO
package util

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/driver/mongodb"
)

// DBExecQuery TODO
type DBExecQuery struct {
}

// NewDBExecQuery TODO
func NewDBExecQuery() *DBExecQuery {
	return &DBExecQuery{}
}

// ExecQuery get info from table with condition
func (query DBExecQuery) ExecQuery(kit *rest.Kit, tableName string, fields []string, condMap mapstr.MapStr, result interface{}) error {
	dbFind := mongodb.Client().Table(tableName).Find(condMap)
	if len(fields) > 0 {
		dbFind = dbFind.Fields(fields...)
	}
	err := dbFind.All(kit.Ctx, result)
	if err != nil {
		blog.Errorf("query table[%s] error. condition: %+v, err: %v, rid: %s", tableName, condMap, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}
	return nil
}
