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
	"context"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

func (s *coreService) SearchInstCount(ctx *rest.Contexts) {
	opt := make(map[string]interface{})
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	count, err := s.core.StatisticOperation().SearchInstCount(ctx.Kit, opt)
	if err != nil {
		ctx.RespEntityWithError(0, err)
		return
	}
	ctx.RespEntity(count)
}

func (s *coreService) SearchChartData(ctx *rest.Contexts) {
	condition := metadata.ChartConfig{}
	if err := ctx.DecodeInto(&condition); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.StatisticOperation().SearchChartData(ctx.Kit, condition)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *coreService) DeleteOperationChart(ctx *rest.Contexts) {
	id := ctx.Request.PathParameter("id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		blog.Errorf("delete chart fail, string convert to int64 fail, err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if _, err := s.core.StatisticOperation().DeleteOperationChart(ctx.Kit, int64ID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *coreService) CreateOperationChart(ctx *rest.Contexts) {
	chartConfig := metadata.ChartConfig{}
	if err := ctx.DecodeInto(&chartConfig); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ownerID := util.GetOwnerID(ctx.Kit.Header)
	chartConfig.CreateTime.Time = time.Now()
	chartConfig.OwnerID = ownerID
	result, err := s.core.StatisticOperation().CreateOperationChart(ctx.Kit, chartConfig)
	if err != nil {
		blog.Errorf("create chart fail, err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *coreService) SearchChartWithPosition(ctx *rest.Contexts) {
	opt := make(map[string]interface{})
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.StatisticOperation().SearchOperationChart(ctx.Kit, opt)
	if err != nil {
		blog.Errorf("search chart fail, err: %v, option: %v, rid: %v", err, opt, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if result == nil {
		ctx.RespEntityWithCount(0, nil)
		return
	}

	lang := s.Language(ctx.Kit.Header)
	for index, chart := range result.Host {
		result.Host[index].Name = s.TranslateOperationChartName(lang, chart)
	}
	for index, chart := range result.Inst {
		result.Inst[index].Name = s.TranslateOperationChartName(lang, chart)
	}

	count, err := mongodb.Client().Table(common.BKTableNameChartConfig).Find(opt).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("search chart fail, option: %v, err: %v, rid: %v", opt, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(int64(count), result)
}

func (s *coreService) UpdateOperationChart(ctx *rest.Contexts) {
	opt := make(map[string]interface{})
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.StatisticOperation().UpdateOperationChart(ctx.Kit, opt)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *coreService) UpdateChartPosition(ctx *rest.Contexts) {
	opt := metadata.ChartPosition{}
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.StatisticOperation().UpdateChartPosition(ctx.Kit, opt)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *coreService) SearchTimerChartData(ctx *rest.Contexts) {
	opt := metadata.ChartConfig{}
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.StatisticOperation().SearchTimerChartData(ctx.Kit, opt)
	if err != nil {
		blog.Errorf("search operation chart data fail, chartName: %v, err: %v, rid: %v", opt.Name, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return

	}

	ctx.RespEntity(result)
}

func (s *coreService) SearchChartCommon(ctx *rest.Contexts) {
	opt := make(map[string]interface{})
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	chartConfig := make([]metadata.ChartConfig, 0)
	if err := mongodb.Client().Table(common.BKTableNameChartConfig).Find(opt).All(ctx.Kit.Ctx, &chartConfig); err != nil {
		blog.Errorf("search chart config fail, option: %v, err: %v, rid: %v", opt, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	count, err := mongodb.Client().Table(common.BKTableNameChartConfig).Find(opt).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("search chart fail, opt: %v, err: %v, rid: %v", opt, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(chartConfig) > 0 {
		ctx.RespEntityWithCount(int64(count), chartConfig[0])
		return
	}

	ctx.RespEntityWithCount(int64(count), nil)
}

func (s *coreService) TimerFreshData(ctx *rest.Contexts) {
	exist, err := mongodb.Client().HasTable(context.Background(), common.BKTableNameChartData)
	if err != nil {
		blog.Errorf("TimerFreshData, update timer chart data fail, err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespEntity(false)
		return
	}
	if !exist {
		ctx.RespEntity(false)
		return
	}
	ctx.SetReadPreference(common.SecondaryPreferredMode)
	err = s.core.StatisticOperation().TimerFreshData(ctx.Kit)
	if err != nil {
		blog.Errorf("TimerFreshData fail, err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(true)
}

func (s *coreService) SearchCloudMapping(ctx *rest.Contexts) {
	opt := make(map[string]interface{})
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	respData := new(metadata.CloudMapping)
	if err := mongodb.Client().Table(common.BKTableNameChartConfig).Find(opt).All(ctx.Kit.Ctx, respData); err != nil {
		blog.Errorf("search cloud mapping fail, err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(respData)
}
