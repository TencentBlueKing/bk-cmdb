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
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) SearchInstCount(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}

	count, err := s.core.StatisticOperation().SearchInstCount(params, opt)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *coreService) SearchChartDataCommon(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	condition := metadata.ChartConfig{}
	if err := data.MarshalJSONInto(&condition); err != nil {
		blog.Errorf("search chart data fail, marshal chart config fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	result, err := s.core.StatisticOperation().SearchChartDataCommon(params, condition)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteOperationChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	id := pathParams("id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		blog.Errorf("delete chart fail, string convert to int64 fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}
	if _, err := s.core.StatisticOperation().DeleteOperationChart(params, int64ID); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *coreService) CreateOperationChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	chartConfig := metadata.ChartConfig{}
	if err := data.MarshalJSONInto(&chartConfig); err != nil {
		blog.Errorf("create chart fail, marshal chart config fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	ownerID := util.GetOwnerID(params.Header)
	chartConfig.CreateTime.Time = time.Now()
	chartConfig.OwnerID = ownerID
	result, err := s.core.StatisticOperation().CreateOperationChart(params, chartConfig)
	if err != nil {
		blog.Errorf("create chart fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	return result, nil
}

func (s *coreService) SearchChartWithPosition(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}

	result, err := s.core.StatisticOperation().SearchOperationChart(params, opt)
	if err != nil {
		blog.Errorf("search chart fail, err: %v, option: %v, rid: %v", err, opt, params.ReqID)
		return nil, err
	}

	if result == nil {
		return struct {
			Count uint64      `json:"count"`
			Info  interface{} `json:"info"`
		}{
			Count: 0,
			Info:  result,
		}, err
	}

	for index, chart := range result.Host {
		result.Host[index].Name = s.TranslateOperationChartName(params.Lang, chart)
	}
	for index, chart := range result.Inst {
		result.Inst[index].Name = s.TranslateOperationChartName(params.Lang, chart)
	}

	count, err := s.db.Table(common.BKTableNameChartConfig).Find(opt).Count(params.Context)
	if err != nil {
		blog.Errorf("search chart fail, option: %v, err: %v, rid: %v", opt, err, params.ReqID)
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
		blog.Errorf("update chart fail, marshal chart config fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	result, err := s.core.StatisticOperation().UpdateOperationChart(params, opt)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *coreService) UpdateChartPosition(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := metadata.ChartPosition{}
	if err := data.MarshalJSONInto(&opt); err != nil {
		blog.Errorf("update chart position fail, marshal chart position fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}
	result, err := s.core.StatisticOperation().UpdateChartPosition(params, opt)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *coreService) SearchTimerChartData(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := metadata.ChartConfig{}
	if err := data.MarshalJSONInto(&opt); err != nil {
		blog.Errorf("search chart data fail, marshal chart config fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	result, err := s.core.StatisticOperation().SearchTimerChartData(params, opt)
	if err != nil {
		blog.Errorf("search operation chart data fail, chartName: %v, err: %v, rid: %v", opt.Name, err, params.ReqID)
		return nil, err

	}

	return result, nil
}

func (s *coreService) SearchChartCommon(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}
	if err := data.MarshalJSONInto(&opt); err != nil {
		blog.Errorf(" search chart fail, marshal chart config fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	chartConfig := make([]metadata.ChartConfig, 0)
	if err := s.db.Table(common.BKTableNameChartConfig).Find(opt).All(params.Context, &chartConfig); err != nil {
		blog.Errorf("search chart config fail, option: %v, err: %v, rid: %v", opt, err, params.ReqID)
		return nil, err
	}

	count, err := s.db.Table(common.BKTableNameChartConfig).Find(opt).Count(params.Context)
	if err != nil {
		blog.Errorf("search chart fail, opt: %v, err: %v, rid: %v", opt, err, params.ReqID)
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
	exist, err := s.db.HasTable(params, common.BKTableNameChartData)
	if err != nil {
		blog.Errorf("TimerFreshData, update timer chart data fail, err: %v, rid: %v", err, params.ReqID)
		return false, nil
	}
	if !exist {
		return false, nil
	}

	err = s.core.StatisticOperation().TimerFreshData(params)
	if err != nil {
		blog.Errorf("TimerFreshData fail, err: %v, rid: %v", err, params.ReqID)
	}

	return true, nil
}

func (s *coreService) SearchCloudMapping(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := mapstr.MapStr{}

	respData := new(metadata.CloudMapping)
	if err := s.db.Table(common.BKTableNameChartConfig).Find(opt).All(params.Context, respData); err != nil {
		blog.Errorf("search cloud mapping fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	return respData, nil
}
