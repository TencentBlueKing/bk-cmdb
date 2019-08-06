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

package logics

import (
	"context"

	"github.com/robfig/cron"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	charts "configcenter/src/scene_server/admin_server/upgrader/x19.05.22.01"
)

func (lgc *Logics) GetBizModuleHostCount(kit *rest.Kit) ([]metadata.IDStringCountInt64, error) {
	cond := metadata.QueryCondition{}
	data := make([]metadata.IDStringCountInt64, 0)
	target := [2]string{common.BKInnerObjIDApp, common.BKInnerObjIDHost}

	for _, obj := range target {
		if obj == common.BKInnerObjIDApp {
			cond = metadata.QueryCondition{
				Condition: mapstr.MapStr{"bk_data_status": mapstr.MapStr{"$ne": "disabled"}},
			}
		}
		result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, obj, &cond)
		if err != nil {
			blog.Errorf("search %v amount failed, err: %v", obj, err)
			return nil, kit.CCError.Error(common.CCErrOperationBizModuleHostAmountFail)
		}
		info := metadata.IDStringCountInt64{}
		if obj == common.BKInnerObjIDApp {
			info.Id = obj
			info.Count = int64(result.Data.Count) - 1
		} else {
			info.Id = obj
			info.Count = int64(result.Data.Count)
		}
		data = append(data, info)
	}

	return data, nil
}

func (lgc *Logics) GetModelAndInstCount(kit *rest.Kit) ([]metadata.IDStringCountInt64, error) {
	cond := &metadata.QueryCondition{}
	result, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("search model fail , err: %v", err)
		return nil, err
	}

	info := make([]metadata.IDStringCountInt64, 0)
	info = append(info, metadata.IDStringCountInt64{
		Id:    "model",
		Count: result.Data.Count,
	})

	opt := make(map[string]interface{})
	resp, err := lgc.CoreAPI.CoreService().Operation().SearchInstCount(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("get instance number fail, err: %v", err)
		return nil, err
	}
	info = append(info, metadata.IDStringCountInt64{
		Id:    "inst",
		Count: int64(resp.Data),
	})

	return info, nil
}

func (lgc *Logics) CreateInnerChart(kit *rest.Kit, chartInfo *metadata.ChartConfig) (uint64, error) {
	opt, ok := charts.InnerChartsMap[chartInfo.ReportType]
	if !ok {
		return 0, kit.CCError.Error(common.CCErrOperationNewAddStatisticFail)
	}

	opt.Width = chartInfo.Width
	opt.XAxisCount = chartInfo.XAxisCount

	result, err := lgc.CoreAPI.CoreService().Operation().CreateOperationChart(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("search chart info fail, err: %v", err)
		return 0, err
	}

	return result.Data, nil
}

func (lgc *Logics) TimerFreshData(ctx context.Context) {
	opt := mapstr.MapStr{}

	if _, err := lgc.CoreAPI.CoreService().Operation().TimerFreshData(ctx, lgc.header, opt); err != nil {
		blog.Error("start collect chart data timer fail, err: %v", err)
		return
	}

	c := cron.New()
	spec := "0 0 2 * * ?" // 每天凌晨两点，更新定时统计图表数据
	err := c.AddFunc(spec, func() {
		// 主服务器跑定时
		isMaster := lgc.Engine.ServiceManageInterface.IsMaster()
		if isMaster {
			if _, err := lgc.CoreAPI.CoreService().Operation().TimerFreshData(ctx, lgc.header, opt); err != nil {
				blog.Error("start collect chart data timer fail, err: %v", err)
			}
		}
	})

	if err != nil {
		blog.Error("start collect chart data timer fail, err: %v", err)
	}
	c.Start()

	select {}
}

func (lgc *Logics) InnerChartData(kit *rest.Kit, chartInfo metadata.ChartConfig) (interface{}, error) {
	switch chartInfo.ReportType {
	case common.BizModuleHostChart:
		data, err := lgc.GetBizModuleHostCount(kit)
		if err != nil {
			return nil, err
		}
		return data, nil
	case common.ModelAndInstCount:
		data, err := lgc.GetModelAndInstCount(kit)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		result, err := lgc.CoreAPI.CoreService().Operation().SearchOperationChartData(kit.Ctx, kit.Header, chartInfo)
		if err != nil {
			blog.Error("search chart data fail, chart name: %v, err: %v", chartInfo.Name, err)
			return nil, err
		}
		return result.Data, nil
	}
}
