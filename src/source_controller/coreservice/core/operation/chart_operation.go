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

	chartPosition := metadata.ChartPosition{}
	if err := m.dbProxy.Table(common.BKTableNameChartPosition).Find(opt).All(ctx, &chartPosition); err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return nil, err
	}

	// 需要改一下chartConfig结构，objID放外面
	for _, chart := range chartConfig {
		matched := false
		for _, info := range chartPosition.Position["host"] {
			if chart.ConfigID == info.ConfigId {
				chart.ChartPosition = info
				matched = true
				continue
			}
		}

		if matched {
			continue
		}

		for _, info := range chartPosition.Position["inst"] {
			if chart.ConfigID == info.ConfigId {
				chart.ChartPosition = info
				continue
			}
		}
	}

	return chartConfig, nil
}

func (m *operationManager) CreateOperationChart(ctx core.ContextParams, inputParam metadata.ChartConfig) (interface{}, error) {
	objID, err := m.dbProxy.NextSequence(ctx, common.BKTableNameCloudTask)
	if err != nil {
		return nil, err
	}
	inputParam.ConfigID = objID

	if err := m.dbProxy.Table(common.BKTableNameChartConfig).Insert(ctx, inputParam); err != nil {
		blog.Errorf("create chart fail, err: %v", err)
		return nil, err
	}

	blog.Debug("objID: %v", objID)
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

func (m *operationManager) DeleteOperationChart(ctx core.ContextParams, id uint64) (interface{}, error) {
	opt := mapstr.MapStr{}
	opt["config_id"] = id
	blog.Debug("opt: %v", opt)
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
