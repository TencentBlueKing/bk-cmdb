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

package operation

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.StatisticOperation = (*operationManager)(nil)

type operationManager struct {
	dbProxy dal.RDB
}

type M bson.M

// New create a new instance manager instance
func New(dbProxy dal.RDB) core.StatisticOperation {
	return &operationManager{
		dbProxy: dbProxy,
	}
}

func (m *operationManager) SearchInstCount(ctx core.ContextParams, inputParam mapstr.MapStr) (uint64, error) {
	count, err := m.dbProxy.Table(common.BKTableNameBaseInst).Find(inputParam).Count(ctx)
	if nil != err {
		blog.Errorf("query database error:%s, condition:%v", err.Error(), inputParam)
		return 0, err
	}

	return count, nil
}

func (m *operationManager) SearchBizHost(ctx core.ContextParams) ([]metadata.IntIDCount, error) {

	bizHostCount := make([]metadata.IntIDCount, 0)

	pipeline := []M{{"$group": M{"_id": "$bk_biz_id", "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameModuleHostConfig).AggregateAll(ctx, pipeline, &bizHostCount); err != nil {
		blog.Errorf("biz' host count aggregate fail, err: %v", err)
		return nil, err
	}

	return bizHostCount, nil
}

func (m *operationManager) CommonAggregate(ctx core.ContextParams, inputParam metadata.ChartConfig) (interface{}, error) {
	commonCount := make([]metadata.StringIDCount, 0)
	filterCondition := fmt.Sprintf("$%v", inputParam.Option.Field)

	switch inputParam.ObjID {
	case common.BKInnerObjIDHost:
		pipeline := []M{{"$group": M{"_id": filterCondition, "count": M{"$sum": 1}}}}
		if err := m.dbProxy.Table(common.BKTableNameBaseHost).AggregateAll(ctx, pipeline, &commonCount); err != nil {
			blog.Errorf("model's instance count aggregate fail, err: %v", err)
			return nil, err
		}
	default:
		pipeline := []M{{"$match": M{"bk_obj_id": inputParam.ObjID}}, {"$group": M{"_id": filterCondition, "count": M{"$sum": 1}}}}
		if err := m.dbProxy.Table(common.BKTableNameBaseInst).AggregateAll(ctx, pipeline, &commonCount); err != nil {
			blog.Errorf("model's instance count aggregate fail, err: %v", err)
			return nil, err
		}
	}

	return commonCount, nil
}
