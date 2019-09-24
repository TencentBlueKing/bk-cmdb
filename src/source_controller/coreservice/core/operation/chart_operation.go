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

func (m *operationManager) SearchOperationChart(ctx core.ContextParams, inputParam interface{}) (*metadata.ChartClassification, error) {
	opt := mapstr.MapStr{}
	chartConfig := make([]metadata.ChartConfig, 0)

	if err := m.dbProxy.Table(common.BKTableNameChartConfig).Find(opt).All(ctx, &chartConfig); err != nil {
		blog.Errorf("SearchOperationChart fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, ctx.Error.CCError(1116005)
	}

	chartPosition := make([]metadata.ChartPosition, 0)
	if err := m.dbProxy.Table(common.BKTableNameChartPosition).Find(opt).All(ctx, &chartPosition); err != nil {
		blog.Errorf("SearchOperationChart fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, ctx.Error.CCError(1116005)
	}

	if len(chartConfig) == 0 || len(chartPosition) == 0 {
		return nil, nil
	}

	// 两个for循环是为了确定图表位置
	chartsInfo := &metadata.ChartClassification{}
	for _, id := range chartPosition[0].Position.Host {
		for _, chart := range chartConfig {
			if chart.ConfigID == id {
				chartsInfo.Host = append(chartsInfo.Host, chart)
				break
			}
		}
	}
	for _, id := range chartPosition[0].Position.Inst {
		for _, chart := range chartConfig {
			if chart.ConfigID == id {
				chartsInfo.Inst = append(chartsInfo.Inst, chart)
				break
			}
		}
	}

	for _, chart := range chartConfig {
		if chart.ReportType == common.BizModuleHostChart {
			chartsInfo.Nav = append(chartsInfo.Nav, chart)
		}
		if chart.ReportType == common.ModelAndInstCount {
			chartsInfo.Nav = append(chartsInfo.Nav, chart)
		}
	}

	return chartsInfo, nil
}

func (m *operationManager) CreateOperationChart(ctx core.ContextParams, inputParam metadata.ChartConfig) (uint64, error) {
	configID, err := m.dbProxy.NextSequence(ctx, common.BKTableNameCloudTask)
	if err != nil {
		return 0, err
	}
	inputParam.ConfigID = configID

	if err := m.dbProxy.Table(common.BKTableNameChartConfig).Insert(ctx, inputParam); err != nil {
		blog.Errorf("CreateOperationChart fail, err: %v, rid: %v", err, ctx.ReqID)
		return 0, ctx.Error.CCError(1116002)
	}

	return configID, nil
}

func (m *operationManager) UpdateChartPosition(ctx core.ContextParams, inputParam interface{}) (interface{}, error) {
	opt := mapstr.MapStr{}

	if err := m.dbProxy.Table(common.BKTableNameChartPosition).Delete(ctx, opt); err != nil {
		blog.Errorf("UpdateChartPosition, delete chart position info fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, ctx.Error.CCError(1116008)
	}

	if err := m.dbProxy.Table(common.BKTableNameChartPosition).Insert(ctx, inputParam); err != nil {
		blog.Errorf("UpdateChartPosition, update chart position fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, ctx.Error.CCError(1116008)
	}

	return nil, nil
}

func (m *operationManager) DeleteOperationChart(ctx core.ContextParams, id int64) (interface{}, error) {
	opt := mapstr.MapStr{}
	opt[common.OperationConfigID] = id
	if err := m.dbProxy.Table(common.BKTableNameChartConfig).Delete(ctx, opt); err != nil {
		blog.Errorf("DeleteOperationChart fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, ctx.Error.CCError(1116004)
	}

	return nil, nil
}

func (m *operationManager) UpdateOperationChart(ctx core.ContextParams, inputParam mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}
	opt[common.OperationConfigID] = inputParam[common.OperationConfigID]
	if err := m.dbProxy.Table(common.BKTableNameChartConfig).Update(ctx, opt, inputParam); err != nil {
		blog.Errorf("UpdateOperationChart fail,chartName: %v, id: %v err: %v, rid: %v", opt["name"], inputParam[common.OperationConfigID], err, ctx.ReqID)
		return nil, ctx.Error.CCError(1116006)
	}

	return nil, nil
}
