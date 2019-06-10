/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
<<<<<<< HEAD
=======
	"configcenter/src/common/http/rest"
>>>>>>> c7685d399... fix: operation crud bugs
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/operation_server/core"
)

func (s *Service) CreateStatisticChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	chartInfo := new(metadata.ChartConfig)
	if err := data.MarshalJSONInto(chartInfo); err != nil {
		blog.Errorf("create statistical chart fail, err: %v", err)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	// 自定义报表
	if chartInfo.ReportType == common.OperationCustom {
		result, err := s.Engine.CoreAPI.CoreService().Operation().CreateOperationChart(params.Context, params.Header, chartInfo)
		if err != nil {
			blog.Errorf("new add statistic fail, err: %v", err)
			return nil, params.Error.Error(common.CCErrOperationNewAddStatisticFail)
		}

		return result, nil
	}

	// 内置报表
	resp, err := s.Core.CreateInnerChart(params, chartInfo)
	if err != nil {
		blog.Errorf("new add statistic fail, err: %v", err)
		return nil, params.Error.Error(common.CCErrOperationNewAddStatisticFail)
	}

	return resp, nil
}

<<<<<<< HEAD
func (s *Service) DeleteStatisticChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	result, err := s.Engine.CoreAPI.CoreService().Operation().DeleteOperationChart(params.Context, params.Header, pathParams("id"))
=======
func (o *OperationServer) DeleteStatisticChart(ctx *rest.Contexts) {
	opt := mapstr.MapStr{}
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	_, err := o.Engine.CoreAPI.CoreService().Operation().DeleteOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, opt)
>>>>>>> c7685d399... fix: operation crud bugs
	if err != nil {
		blog.Errorf("search chart info fail, err: %v, id: %v", err, pathParams)
		return nil, err
	}

	return result, nil
}

<<<<<<< HEAD
func (s *Service) SearchStatisticCharts(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	opt := make(map[string]interface{})
	result, err := s.Engine.CoreAPI.CoreService().Operation().SearchOperationChart(params.Context, params.Header, opt)
=======
func (o *OperationServer) SearchStatisticChart(ctx *rest.Contexts) {
	opt := make(map[string]interface{})

	result, err := o.Engine.CoreAPI.CoreService().Operation().SearchOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, opt)
>>>>>>> c7685d399... fix: operation crud bugs
	if err != nil {
		blog.Errorf("search operation field info fail, err: %v", err)
		return nil, params.Error.Error(common.CCErrOperationSearchStatisticsFail)
	}

<<<<<<< HEAD
	return result, nil
=======
	blog.Debug("result: %v", result)
	ctx.RespEntity(result.Data)
>>>>>>> c7685d399... fix: operation crud bugs
}

func (s *Service) UpdateStatisticChart(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	chartInfo := new(metadata.ChartConfig)
	if err := data.MarshalJSONInto(chartInfo); err != nil {
		blog.Errorf("create statistical chart fail, err: %v", err)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.Engine.CoreAPI.CoreService().Operation().UpdateOperationChart(params.Context, params.Header, chartInfo)
	if err != nil {
		blog.Errorf("update statistic info fail, err: %v", err)
		return nil, params.Error.Error(common.CCErrOperationUpdateStatisticsFail)
	}

<<<<<<< HEAD
	return result, nil
}

func (s *Service) SearchChartData(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	innerChart := []string{"host_change_biz_chart", "model_inst_chart", "model_inst_change_chart"}

	inputData := metadata.ChartConfig{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	if !util.InStrArr(innerChart, inputData.ReportType) {
		result, err := s.Core.CommonStatisticFunc(params, inputData.Option)
=======
	ctx.RespEntity(result.Data)
}

func (o *OperationServer) SearchChartData(ctx *rest.Contexts) {
	inputData := new(metadata.ChartConfig)
	if err := ctx.DecodeInto(inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 判断模型是否存在，不存在返回nil
	cond := make(map[string]interface{}, 0)
	cond[common.BKObjIDField] = inputData.ObjID
	query := metadata.QueryCondition{Condition: cond}
	models, err := o.CoreAPI.CoreService().Model().ReadModel(ctx.Kit.Ctx, ctx.Kit.Header, &query)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, err: %v, chart name: %v", err)
		return
	}
	if models.Data.Count <= 0 {
		ctx.RespEntity(nil)
		return
	}

	innerChart := []string{
		"host_change_biz_chart", "model_inst_chart", "model_inst_change_chart",
		"biz_module_host_chart", "model_and_inst_count",
	}

	srvData := o.newSrvComm(ctx.Kit.Header)
	if !util.InStrArr(innerChart, inputData.ReportType) {
		result, err := srvData.lgc.CommonStatisticFunc(ctx.Kit, inputData)
>>>>>>> c7685d399... fix: operation crud bugs
		if err != nil {
			blog.Errorf("search chart data fail, err: %v, chart name: %v", err, inputData.Name)
			return nil, err
		}
		return result, nil
	}

<<<<<<< HEAD
	result, err := s.Engine.CoreAPI.CoreService().Operation().SearchOperationChartData(params.Context, params.Header, inputData.ReportType)
=======
	result, err := srvData.lgc.GetInnerChartData(ctx.Kit, inputData)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, err: %v, chart name: %v", err, inputData.Name)
		return
	}

	ctx.RespEntity(result)
}

func (o *OperationServer) UpdateChartPosition(ctx *rest.Contexts) {
	opt := mapstr.MapStr{}
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := o.CoreAPI.CoreService().Operation().UpdateOperationChartPosition(ctx.Kit.Ctx, ctx.Kit.Header, opt)
>>>>>>> c7685d399... fix: operation crud bugs
	if err != nil {
		blog.Errorf("search chart data fail, err: %v, chart name: %v", err, inputData.Name)
		return nil, err
	}

	return result, nil
}
