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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
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
	versionKey := getSynchronize(common.MetadataField, common.MetaDataSynchronizeVersionField)
	flagKey := getSynchronize(common.MetadataField, common.MetaDataSynchronizeFlagField)

	delConditionParse := condition.CreateCondition()
	delConditionParse.Field(versionKey).Lt(c.input.Version)
	delConditionParse.Field(flagKey).Eq(c.input.SynchronizeFlag)
	deleteConditon := delConditionParse.ToMapStr()

	conditionParse := condition.CreateCondition()
	conditionParse.Field(versionKey).Eq(c.input.Version)
	conditionParse.Field(flagKey).Eq(c.input.SynchronizeFlag)
	queryCondition := conditionParse.ToMapStr()

	tableNameArr := common.AllTables
	for _, tableName := range tableNameArr {
		cnt, err := c.dbProxy.Table(tableName).Find(queryCondition).Count(ctx)
		if err != nil {
			blog.Errorf("clearData  find %s table row error, err:%s,rid:%s", tableName, err.Error(), ctx.ReqID)
			continue
		}
		if cnt <= 0 {
			// not current version data. not execute delete row
			continue
		}

		err = c.dbProxy.Table(tableName).Delete(ctx, deleteConditon)
		if err != nil {
			blog.Errorf("clearData  delete %s table row error, err:%s,rid:%s", tableName, err.Error(), ctx.ReqID)
		}
	}
}

func getSynchronize(key ...string) string {
	return strings.Join(key, ".")
}
