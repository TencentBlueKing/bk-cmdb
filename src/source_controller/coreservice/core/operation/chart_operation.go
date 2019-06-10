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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *operationManager) SearchOperationChart(ctx core.ContextParams, inputParam interface{}) (interface{}, error) {
	opt := mapstr.MapStr{}
	chartConfig := make([]metadata.ChartConfig, 0)

	if err := m.dbProxy.Table(common.BKTableNameChartConfig).Find(opt).All(ctx, &chartConfig); err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return nil, err
	}

	chartPosition := make([]metadata.ChartPosition, 0)
	if err := m.dbProxy.Table(common.BKTableNameChartPosition).Find(opt).All(ctx, &chartPosition); err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return nil, err
	}

	if len(chartConfig) == 0 || len(chartPosition) == 0 {
		return nil, nil
	}

	chartsInfo := make(map[string][]interface{})
	for _, info := range chartPosition[0].Position["host"] {
		for _, chart := range chartConfig {
			if chart.ConfigID == info.ConfigId {
				chartsInfo["host"] = append(chartsInfo["host"], chart)
			}
		}
	}

	for _, info := range chartPosition[0].Position["inst"] {
		for _, chart := range chartConfig {
			if chart.ConfigID == info.ConfigId {
				chartsInfo["inst"] = append(chartsInfo["inst"], chart)
			}
		}
	}

	return chartsInfo, nil
}

func (m *operationManager) CreateOperationChart(ctx core.ContextParams, inputParam metadata.ChartConfig) (uint64, error) {
	objID, err := m.dbProxy.NextSequence(ctx, common.BKTableNameCloudTask)
	if err != nil {
		return 0, err
	}
	inputParam.ConfigID = objID

	if err := m.dbProxy.Table(common.BKTableNameChartConfig).Insert(ctx, inputParam); err != nil {
		blog.Errorf("create chart fail, err: %v", err)
		return 0, err
	}

	return objID, nil
}

func (m *operationManager) UpdateChartPosition(ctx core.ContextParams, inputParam interface{}) (interface{}, error) {
	opt := mapstr.MapStr{}

	if err := m.dbProxy.Table(common.BKTableNameChartPosition).Delete(ctx, opt); err != nil {
		blog.Errorf("delete chart position info fail, err: %v", err)
		return nil, err
	}

	if err := m.dbProxy.Table(common.BKTableNameChartPosition).Insert(ctx, inputParam); err != nil {
		blog.Errorf("update chart position fail, err: %v", err)
		return nil, err
	}

	return nil, nil
}

func (m *operationManager) DeleteOperationChart(ctx core.ContextParams, inputParam mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}
	condition := mapstr.MapStr{}
	condition["$in"] = inputParam["id"]
	opt["config_id"] = condition
	if err := m.dbProxy.Table(common.BKTableNameChartConfig).Delete(ctx, opt); err != nil {
		blog.Errorf("create chart fail, err: %v", err)
		return nil, err
	}

	return nil, nil
}

func (m *operationManager) UpdateOperationChart(ctx core.ContextParams, inputParam metadata.ChartConfig) (interface{}, error) {
	opt := mapstr.MapStr{}
	opt["config_id"] = inputParam.ConfigID
	blog.Debug("input: %v", inputParam)
	if err := m.dbProxy.Table(common.BKTableNameChartConfig).Update(ctx, opt, inputParam); err != nil {
		blog.Errorf("update chart config fail,id: %v err: %v", inputParam.ConfigID, err)
		return nil, err
	}

	return nil, nil
}
