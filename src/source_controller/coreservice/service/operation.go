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

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"time"
)

func (s *coreService) SearchInstCount(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}

	count, err := s.core.StatisticOperation().SearchInstCount(params, opt)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *coreService) CommonAggregate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	condition := metadata.ChartConfig{}
	if err := data.MarshalJSONInto(&condition); err != nil {
		blog.Errorf("marshal chart config fail, err: %v", err)
		return nil, err
	}

	result, err := s.core.StatisticOperation().CommonAggregate(params, condition)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteOperationChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}
	if err := data.MarshalJSONInto(&opt); err != nil {
		blog.Errorf("marshal request data fail, err: %v", err)
		return nil, err
	}
	if _, err := s.core.StatisticOperation().DeleteOperationChart(params, opt); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *coreService) CreateOperationChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	chartConfig := metadata.ChartConfig{}
	if err := data.MarshalJSONInto(&chartConfig); err != nil {
		blog.Errorf("marshal chart config fail, err: %v", err)
		return nil, err
	}

	ownerID := util.GetOwnerID(params.Header)
	chartConfig.CreateTime = time.Now()
	chartConfig.OwnerID = ownerID
	result, err := s.core.StatisticOperation().CreateOperationChart(params, chartConfig)
	if err != nil {
		blog.Errorf("save chart config fail, err: %v", err)
		return nil, err
	}
	blog.Debug("result: %v", result)
	return result, nil
}

func (s *coreService) SearchOperationChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}

	result, err := s.core.StatisticOperation().SearchOperationChart(params, opt)
	if err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return nil, err
	}

	count, err := s.db.Table(common.BKTableNameChartConfig).Find(opt).Count(params.Context)
	if err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return nil, err
	}

	return struct {
		Count uint64      `json:"count"`
		Info  interface{} `json:"info"`
	}{
		Count: count,
		Info:  result,
	}, err
}

func (s *coreService) UpdateOperationChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}
	if err := data.MarshalJSONInto(&opt); err != nil {
		blog.Errorf("marshal chart config fail, err: %v", err)
		return nil, err
	}

	result, err := s.core.StatisticOperation().UpdateOperationChart(params, opt)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *coreService) UpdateOperationChartPosition(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	result, err := s.core.StatisticOperation().UpdateChartPosition(params, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *coreService) SearchOperationChartData(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := metadata.ChartConfig{}
	if err := data.MarshalJSONInto(&opt); err != nil {
		blog.Errorf("marshal chart config fail, err: %v", err)
		return nil, err
	}

	result, err := s.core.StatisticOperation().SearchOperationChartData(params, opt)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *coreService) SearchChartCommon(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}
	if err := data.MarshalJSONInto(&opt); err != nil {
		blog.Errorf("marshal chart config fail, err: %v", err)
		return nil, err
	}

	chartConfig := make([]metadata.ChartConfig, 0)
	if err := s.db.Table(common.BKTableNameChartConfig).Find(opt).All(params.Context, &chartConfig); err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return nil, err
	}

	count, err := s.db.Table(common.BKTableNameChartConfig).Find(opt).Count(params.Context)
	if err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return nil, err
	}

	if len(chartConfig) > 0 {
		return struct {
			Count uint64      `json:"count"`
			Info  interface{} `json:"info"`
		}{
			Count: count,
			Info:  chartConfig[0],
		}, err
	}

	return struct {
		Count uint64      `json:"count"`
		Info  interface{} `json:"info"`
	}{
		Count: count,
		Info:  nil,
	}, err
}

func (s *coreService) TimerFreshData(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	s.core.StatisticOperation().TimerFreshData(params)

	return nil, nil
}

func (s *coreService) SearchCloudMapping(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}

	respData := new(metadata.CloudMapping)
	if err := s.db.Table(common.BKTableNameChartConfig).Find(opt).All(params.Context, respData); err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return nil, err
	}

	return respData, nil
}
