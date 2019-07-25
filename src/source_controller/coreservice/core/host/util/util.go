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

package util

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type DBExecQuery struct {
	DbProxy dal.RDB
}

func NewDBExecQuery(dbProxy dal.RDB) *DBExecQuery {
	return &DBExecQuery{
		DbProxy: dbProxy,
	}
}

// dbExecQuery get info from table with condition
func (query DBExecQuery) ExecQuery(ctx core.ContextParams, tableName string, fields []string, condMap mapstr.MapStr, result interface{}) error {
	newCondMap := util.SetQueryOwner(condMap, ctx.SupplierAccount)
	dbFind := query.DbProxy.Table(tableName).Find(newCondMap)
	if len(fields) > 0 {
		dbFind = dbFind.Fields(fields...)
	}
	err := dbFind.All(ctx, result)
	if err != nil {
		blog.ErrorJSON("ExecQuery query table[%s] error. condition: %s, err:%s, rid:%s", tableName, newCondMap, err.Error(), ctx.ReqID)
		return ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	return nil
}
