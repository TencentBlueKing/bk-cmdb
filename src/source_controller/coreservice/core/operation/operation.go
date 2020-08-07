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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.StatisticOperation = (*operationManager)(nil)

type operationManager struct {
	dbProxy dal.RDB
}

type M map[string]interface{}

func New(dbProxy dal.RDB) core.StatisticOperation {
	return &operationManager{
		dbProxy: dbProxy,
	}
}

func (m *operationManager) SearchInstCount(kit *rest.Kit, inputParam map[string]interface{}) (uint64, error) {
	count, err := m.dbProxy.Table(common.BKTableNameBaseInst).Find(inputParam).Count(kit.Ctx)
	if nil != err {
		blog.Errorf("query database error:%s, condition:%v, rid: %v", err.Error(), inputParam, kit.Rid)
		return 0, err
	}

	return count, nil
}

func (m *operationManager) SearchChartData(kit *rest.Kit, inputParam metadata.ChartConfig) (interface{}, error) {
	switch inputParam.ReportType {
	case common.HostCloudChart:
		data, err := m.HostCloudChartData(kit, inputParam)
		if err != nil {
			blog.Error("search host cloud chart data fail, inputParam: %v, err: %v,  rid: %v", inputParam, err, kit.Rid)
			return nil, err
		}
		return data, nil
	case common.HostBizChart:
		data, err := m.HostBizChartData(kit, inputParam)
		if err != nil {
			blog.Error("search biz's host chart data fail, params: %v, err: %v, rid: %v", inputParam, err, kit.Rid)
			return nil, err
		}
		return data, nil
	default:
		data, err := m.CommonModelStatistic(kit, inputParam)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func (m *operationManager) CommonModelStatistic(kit *rest.Kit, inputParam metadata.ChartConfig) (interface{}, error) {
	commonCount := make([]metadata.StringIDCount, 0)
	filterCondition := fmt.Sprintf("$%s", inputParam.Field)

	if inputParam.ObjID == common.BKInnerObjIDHost {
		pipeline := []M{{common.BKDBGroup: M{"_id": filterCondition, "count": M{common.BKDBSum: 1}}}}
		if err := m.dbProxy.Table(common.BKTableNameBaseHost).AggregateAll(kit.Ctx, pipeline, &commonCount); err != nil {
			blog.Errorf("host os type count aggregate fail, chartName: %v, err: %v, rid: %v", inputParam.Name, err, kit.Rid)
			return nil, err
		}
	} else {
		pipeline := []M{{common.BKDBMatch: M{common.BKObjIDField: inputParam.ObjID}}, {common.BKDBGroup: M{"_id": filterCondition, "count": M{common.BKDBSum: 1}}}}
		if err := m.dbProxy.Table(common.BKTableNameBaseInst).AggregateAll(kit.Ctx, pipeline, &commonCount); err != nil {
			blog.Errorf("model's instance count aggregate fail, chartName: %v, ObjID: %v, err: %v, rid: %v", inputParam.Name, inputParam.ObjID, err, kit.Rid)
			return nil, err
		}
	}

	attribute := metadata.Attribute{}
	opt := map[string]interface{}{}
	opt[common.BKObjIDField] = inputParam.ObjID
	opt[common.BKPropertyIDField] = inputParam.Field
	if err := m.dbProxy.Table(common.BKTableNameObjAttDes).Find(opt).One(kit.Ctx, &attribute); err != nil {
		blog.Errorf("model's instance count aggregate fail, chartName: %v, objID: %v, err: %v, rid: %v", inputParam.Name, inputParam.ObjID, err, kit.Rid)
		return nil, err
	}

	instCount := uint64(0)
	cond := M{}
	var countErr error
	if inputParam.ObjID == common.BKInnerObjIDHost {
		instCount, countErr = m.dbProxy.Table(common.BKTableNameBaseHost).Find(cond).Count(kit.Ctx)
		if countErr != nil {
			blog.Errorf("host os type count aggregate fail, chartName: %v, err: %v, rid: %v", inputParam.Name, countErr, kit.Rid)
			return nil, countErr
		}
		if instCount > 0 {
			pipeline := []M{{"$group": M{"_id": filterCondition, "count": M{"$sum": 1}}}}
			if err := m.dbProxy.Table(common.BKTableNameBaseHost).AggregateAll(kit.Ctx, pipeline, &commonCount); err != nil {
				blog.Errorf("host os type count aggregate fail, chartName: %v, err: %v, rid: %v", inputParam.Name, err, kit.Rid)
				return nil, err
			}
		}
	} else {
		instCount, countErr = m.dbProxy.Table(common.BKTableNameBaseInst).Find(cond).Count(kit.Ctx)
		if countErr != nil {
			blog.Errorf("model's instance count aggregate fail, chartName: %v, ObjID: %v, err: %v, rid: %v", inputParam.Name, inputParam.ObjID, countErr, kit.Rid)
			return nil, countErr
		}
		if instCount > 0 {
			pipeline := []M{{"$match": M{"bk_obj_id": inputParam.ObjID}}, {"$group": M{"_id": filterCondition, "count": M{"$sum": 1}}}}
			if err := m.dbProxy.Table(common.BKTableNameBaseInst).AggregateAll(kit.Ctx, pipeline, &commonCount); err != nil {
				blog.Errorf("model's instance count aggregate fail, chartName: %v, ObjID: %v, err: %v, rid: %v", inputParam.Name, inputParam.ObjID, err, kit.Rid)
				return nil, err
			}
		}
	}

	option, err := metadata.ParseEnumOption(kit.Ctx, attribute.Option)
	if err != nil {
		blog.Errorf("count model's instance, parse enum option fail, ObjID: %v, err:%v, rid: %v", inputParam.ObjID, err, kit.Rid)
		return nil, err
	}

	respData := make([]metadata.StringIDCount, 0)
	for _, opt := range option {
		info := metadata.StringIDCount{
			ID:    opt.Name,
			Count: 0,
		}
		for _, count := range commonCount {
			if count.ID == opt.ID {
				info.Count = count.Count
			}
			if opt.Name == common.OptionOther && count.ID == "" {
				info.Count = count.Count
			}
		}
		respData = append(respData, info)
	}

	return respData, nil
}

func (m *operationManager) SearchTimerChartData(kit *rest.Kit, inputParam metadata.ChartConfig) (interface{}, error) {
	condition := map[string]interface{}{}
	condition[common.OperationReportType] = inputParam.ReportType

	switch inputParam.ReportType {
	case common.HostChangeBizChart:
		chartData := make([]metadata.HostChangeChartData, 0)
		if err := m.dbProxy.Table(common.BKTableNameChartData).Find(condition).All(kit.Ctx, &chartData); err != nil {
			blog.Errorf("search chart data fail, chart name: %v err: %v, rid: %v", inputParam.Name, err, kit.Rid)
			return nil, err
		}
		result := make(map[string][]metadata.StringIDCount, 0)
		for _, data := range chartData {
			for _, info := range data.Data {
				if _, ok := result[info.ID]; !ok {
					result[info.ID] = make([]metadata.StringIDCount, 0)
				}
				result[info.ID] = append(result[info.ID], metadata.StringIDCount{
					ID:    data.CreateTime,
					Count: info.Count,
				})
			}
		}
		return result, nil
	case common.ModelInstChart:
		chartData := metadata.ModelInstChartData{}
		if err := m.dbProxy.Table(common.BKTableNameChartData).Find(condition).One(kit.Ctx, &chartData); err != nil {
			blog.Errorf("search chart data fail, chart name: %v err: %v, rid: %v", inputParam.Name, err, kit.Rid)
			return nil, err
		}
		return chartData.Data, nil
	case common.ModelInstChangeChart:
		chartData := metadata.ChartData{}
		if err := m.dbProxy.Table(common.BKTableNameChartData).Find(condition).One(kit.Ctx, &chartData); err != nil {
			blog.Errorf("search chart data fail, chart name: %v err: %v, rid: %v", inputParam.Name, err, kit.Rid)
			return nil, err
		}
		return chartData.Data, nil
	}

	return nil, nil
}
