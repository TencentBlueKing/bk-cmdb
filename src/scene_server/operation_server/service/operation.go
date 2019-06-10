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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (o *OperationServer) CreateStatisticChart(ctx *rest.Contexts) {
	chartInfo := new(metadata.ChartConfig)
	if err := ctx.DecodeInto(chartInfo); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 自定义报表
	if chartInfo.ReportType == common.OperationCustom {
		result, err := o.Engine.CoreAPI.CoreService().Operation().CreateOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, chartInfo)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrOperationNewAddStatisticFail, "new add statistic fail, err: %v", err)
			return
		}

		blog.Debug("count: %v", result.Data)
		ctx.RespEntity(result.Data)
		return
	}

	blog.Debug("create inner chart")
	// 内置报表
	srvData := o.newSrvComm(ctx.Kit.Header)
	resp, err := srvData.lgc.CreateInnerChart(ctx.Kit, chartInfo)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationNewAddStatisticFail, "new add statistic fail, err: %v", err)
		return
	}

	ctx.RespEntity(resp)
}

func (o *OperationServer) DeleteStatisticChart(ctx *rest.Contexts) {
	id := ctx.Request.PathParameter("id")
	_, err := o.Engine.CoreAPI.CoreService().Operation().DeleteOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, id)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationDeleteStatisticFail, "search chart info fail, err: %v, id: %v", err)
		return
	}

	ctx.RespEntity(nil)
}

func (o *OperationServer) SearchStatisticCharts(ctx *rest.Contexts) {

	opt := make(map[string]interface{})
	result, err := o.Engine.CoreAPI.CoreService().Operation().SearchOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationSearchStatisticsFail, "search chart info fail, err: %v", err)
		return
	}

	ctx.RespEntity(result)
}

func (o *OperationServer) UpdateStatisticChart(ctx *rest.Contexts) {
	chartInfo := new(metadata.ChartConfig)
	if err := ctx.DecodeInto(chartInfo); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := o.Engine.CoreAPI.CoreService().Operation().UpdateOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, chartInfo)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationSearchStatisticsFail, "update statistic info fail, err: %v", err)
		return
	}

	ctx.RespEntity(result)
}

func (o *OperationServer) SearchChartData(ctx *rest.Contexts) {
	// todo 判断模型是否存在，不存在返回模型不存在

	innerChart := []string{
		"host_change_biz_chart", "model_inst_chart", "model_inst_change_chart",
		"biz_module_host_chart", "model_and_inst_count",
	}

	inputData := new(metadata.ChartConfig)
	if err := ctx.DecodeInto(inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}

	srvData := o.newSrvComm(ctx.Kit.Header)
	if !util.InStrArr(innerChart, inputData.ReportType) {
		result, err := srvData.lgc.CommonStatisticFunc(ctx.Kit, inputData.Option)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, err: %v, chart name: %v", err, inputData.Name)
			return
		}
		ctx.RespEntity(result)
		return
	}

	result, err := srvData.lgc.GetInnerChartData(ctx.Kit, inputData)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, err: %v, chart name: %v", err, inputData.Name)
		return
	}

	ctx.RespEntity(result)
}

func (o *OperationServer) UpdateChartPosition(ctx *rest.Contexts) {
	result, err := o.CoreAPI.CoreService().Operation().UpdateOperationChartPosition(ctx.Kit.Ctx, ctx.Kit.Header, ctx.Request.Request.Body)
	if err != nil {
		blog.Errorf("update chart position fail, err: %v", err)
		return
	}

	ctx.RespEntity(result)
}
