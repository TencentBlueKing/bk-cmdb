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
package datasynchronize

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type clearDataInterface interface {
	clearData(ctx core.ContextParams)
}

type clearData struct {
	dbProxy dal.RDB
	input   *metadata.SynchronizeClearDataParameter
}

func NewClearData(dbProxy dal.RDB, input *metadata.SynchronizeClearDataParameter) clearDataInterface {
	return &clearData{
		dbProxy: dbProxy,
		input:   input,
	}
}

func (c *clearData) clearData(ctx core.ContextParams) {

	versionKey := util.BuildMongoSyncItemField(common.MetaDataSynchronizeVersionField)
	flagKey := util.BuildMongoSyncItemField(common.MetaDataSynchronizeFlagField)

	delConditionParse := condition.CreateCondition()
	delConditionParse.Field(versionKey).Lt(c.input.Version)
	delConditionParse.Field(flagKey).Eq(c.input.SynchronizeFlag)
	deleteConditon := delConditionParse.ToMapStr()

	blog.V(5).Infof(" clearData condition:%#v, rid:%s", deleteConditon, ctx.ReqID)
	tableNameArr := common.AllTables
	for _, tableName := range tableNameArr {
		cnt, err := c.dbProxy.Table(tableName).Find(deleteConditon).Count(ctx)
		if err != nil {
			blog.Warnf("clearData  find %s table row error, err:%s, condition:%#v, rid:%s", tableName, err.Error(), deleteConditon, ctx.ReqID)
			continue
		}
		if cnt <= 0 {
			// not current version data. not execute delete row
			continue
		}

		err = c.dbProxy.Table(tableName).Delete(ctx, deleteConditon)
		if err != nil {
			blog.Errorf("clearData  delete %s table row error, err:%s, condition:%#v, rid:%s", tableName, err.Error(), deleteConditon, ctx.ReqID)
		}
	}
}
