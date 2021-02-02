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
	"configcenter/src/storage/driver/mongodb"
)

var _ core.StatisticOperation = (*operationManager)(nil)

type operationManager struct {
}

type M map[string]interface{}

func New() core.StatisticOperation {
	return &operationManager{}
}

func (m *operationManager) SearchInstCount(kit *rest.Kit, inputParam map[string]interface{}) (uint64, error) {
	count, err := mongodb.Client().Table(common.BKTableNameBaseInst).Find(inputParam).Count(kit.Ctx)
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
	// get enum options by model's field
	attribute := metadata.Attribute{}
	opt := map[string]interface{}{}
	opt[common.BKObjIDField] = inputParam.ObjID
	opt[common.BKPropertyIDField] = inputParam.Field
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(opt).One(kit.Ctx, &attribute); err != nil {
		blog.Errorf("model's instance count aggregate failed, chartName: %v, objID: %v, err: %v, rid: %v", inputParam.Name, inputParam.ObjID, err, kit.Rid)
		return nil, err
	}

	option, err := metadata.ParseEnumOption(kit.Ctx, attribute.Option)
	if err != nil {
		blog.Errorf("count model's instance, parse enum option fail, ObjID: %v, err:%v, rid: %v", inputParam.ObjID, err, kit.Rid)
		return nil, err
	}

	// get model instances' count group by its field
	// eg: get host count group by bk_os_type
	groupCountArr := make([]metadata.StringIDCount, 0)
	groupField := fmt.Sprintf("$%s", inputParam.Field)
	instCount := uint64(0)
	cond := M{}
	var countErr error
	if inputParam.ObjID == common.BKInnerObjIDHost {
		instCount, countErr = mongodb.Client().Table(common.BKTableNameBaseHost).Find(cond).Count(kit.Ctx)
		if countErr != nil {
			blog.Errorf("model's instance count aggregate failed, chartName: %v, err: %v, rid: %v", inputParam.Name, countErr, kit.Rid)
			return nil, countErr
		}
		if instCount > 0 {
			pipeline := []M{
				{common.BKDBMatch: M{common.BKDBAND: []M{
					{inputParam.Field: M{common.BKDBExists: true}},
					{inputParam.Field: M{common.BKDBNE: nil}},
				}}},
				{common.BKDBGroup: M{"_id": groupField, "count": M{common.BKDBSum: 1}}},
			}
			if err := mongodb.Client().Table(common.BKTableNameBaseHost).AggregateAll(kit.Ctx, pipeline, &groupCountArr); err != nil {
				blog.Errorf("model's instance count aggregate failed, chartName: %v, err: %v, rid: %v", inputParam.Name, err, kit.Rid)
				return nil, err
			}
		}
	} else {
		instCount, countErr = mongodb.Client().Table(common.BKTableNameBaseInst).Find(cond).Count(kit.Ctx)
		if countErr != nil {
			blog.Errorf("model's instance count aggregate fail, chartName: %v, ObjID: %v, err: %v, rid: %v", inputParam.Name, inputParam.ObjID, countErr, kit.Rid)
			return nil, countErr
		}
		if instCount > 0 {
			pipeline := []M{
				{common.BKDBMatch: M{common.BKDBAND: []M{
					{inputParam.Field: M{common.BKDBExists: true}},
					{inputParam.Field: M{common.BKDBNE: nil}},
					{common.BKDBMatch: M{common.BKObjIDField: inputParam.ObjID}},
				}}},
				{common.BKDBGroup: M{"_id": groupField, "count": M{common.BKDBSum: 1}}},
			}
			if err := mongodb.Client().Table(common.BKTableNameBaseInst).AggregateAll(kit.Ctx, pipeline, &groupCountArr); err != nil {
				blog.Errorf("model's instance count aggregate failed, chartName: %v, ObjID: %v, err: %v, rid: %v", inputParam.Name, inputParam.ObjID, err, kit.Rid)
				return nil, err
			}
		}
	}

	if len(groupCountArr) == 0 {
		return []metadata.StringIDCount{}, nil
	}

	groupCountMap := make(map[string]int64)
	for _, groupCount := range groupCountArr {
		groupCountMap[groupCount.ID] = groupCount.Count
	}

	respData := make([]metadata.StringIDCount, 0)
	for _, opt := range option {
		if opt.Name == common.OptionOther {
			respData = append(respData, metadata.StringIDCount{
				ID:    opt.Name,
				Count: groupCountMap[""],
			})
			continue
		}
		respData = append(respData, metadata.StringIDCount{
			ID:    opt.Name,
			Count: groupCountMap[opt.ID],
		})
	}

	return respData, nil
}

func (m *operationManager) SearchTimerChartData(kit *rest.Kit, inputParam metadata.ChartConfig) (interface{}, error) {
	condition := map[string]interface{}{}
	condition[common.OperationReportType] = inputParam.ReportType

	switch inputParam.ReportType {
	case common.HostChangeBizChart:
		chartData := make([]metadata.HostChangeChartData, 0)
		if err := mongodb.Client().Table(common.BKTableNameChartData).Find(condition).All(kit.Ctx, &chartData); err != nil {
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
		if err := mongodb.Client().Table(common.BKTableNameChartData).Find(condition).One(kit.Ctx, &chartData); err != nil {
			blog.Errorf("search chart data fail, chart name: %v err: %v, rid: %v", inputParam.Name, err, kit.Rid)
			return nil, err
		}
		return chartData.Data, nil
	case common.ModelInstChangeChart:
		chartData := metadata.ChartData{}
		if err := mongodb.Client().Table(common.BKTableNameChartData).Find(condition).One(kit.Ctx, &chartData); err != nil {
			blog.Errorf("search chart data fail, chart name: %v err: %v, rid: %v", inputParam.Name, err, kit.Rid)
			return nil, err
		}
		return chartData.Data, nil
	}

	return nil, nil
}
