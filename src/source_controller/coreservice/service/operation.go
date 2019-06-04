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
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"strconv"
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
	condition := metadata.ChartOption{}
	if err := data.MarshalJSONInto(condition); err != nil {
		return nil, err
	}

	result, err := s.core.StatisticOperation().CommonAggregate(params, condition)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteOperationChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	id, err := strconv.ParseUint(pathParams("id"), 10, 64)
	if err != nil {
		blog.Errorf("string convert to u")
		return nil, err
	}
	if _, err := s.core.StatisticOperation().DeleteOperationChart(params, id); err != nil {
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
	blog.Debug("chartInfo---1： %v", chartConfig)
	result, err := s.core.StatisticOperation().CreateOperationChart(params, chartConfig)
	if err != nil {
		blog.Errorf("save chart config fail, err: %v", err)
		return nil, err
	}
	return result, nil
}

func (s *coreService) SearchOperationChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	result, err := s.core.StatisticOperation().SearchOperationChart(params, data)
	if err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return nil, err
	}

	blog.Debug("result: %v", result)
	return result, err
}

func (s *coreService) UpdateOperationChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	chartConfig := metadata.ChartConfig{}
	if err := data.MarshalJSONInto(chartConfig); err != nil {
		blog.Errorf("marshal chart config fail, err: %v", err)
		return nil, err
	}

	result, err := s.core.StatisticOperation().UpdateOperationChart(params, chartConfig)
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
