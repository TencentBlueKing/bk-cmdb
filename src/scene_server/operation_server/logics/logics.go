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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"github.com/robfig/cron"
)

func (lgc *Logics) GetBizHostCount(kit *rest.Kit) ([]metadata.IDStringCountInt64, error) {
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
			blog.Errorf("search %v count failed, err: %v, rid: %v", obj, err, kit.Rid)
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
	condition := mapstr.MapStr{
		"ispre":       false,
		"bk_ispaused": false,
	}
	cond.Condition = condition
	result, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("GetModelAndInstCount fail, search model fail , err: %v, rid: %v", err, kit.Rid)
		return nil, err
	}

	info := make([]metadata.IDStringCountInt64, 0)
	info = append(info, metadata.IDStringCountInt64{
		Id:    "model",
		Count: result.Data.Count, // 去除内置的模型(主机、集群等)
	})

	opt := make(map[string]interface{})
	resp, err := lgc.CoreAPI.CoreService().Operation().SearchInstCount(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("GetModelAndInstCount fail, get instance count fail, err: %v, rid: %v", err, kit.Rid)
		return nil, err
	}
	info = append(info, metadata.IDStringCountInt64{
		Id:    "inst",
		Count: int64(resp.Data),
	})

	return info, nil
}

func (lgc *Logics) CreateInnerChart(kit *rest.Kit, chartInfo *metadata.ChartConfig) (uint64, error) {
	opt, ok := metadata.InnerChartsMap[chartInfo.ReportType]
	if !ok {
		return 0, kit.CCError.Error(common.CCErrOperationNewAddStatisticFail)
	}

	opt.Width = chartInfo.Width
	opt.XAxisCount = chartInfo.XAxisCount

	result, err := lgc.CoreAPI.CoreService().Operation().CreateOperationChart(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("create operation chart fail, err: %v, rid: %v", err, kit.Rid)
		return 0, err
	}

	return result.Data, nil
}

func (lgc *Logics) InnerChartData(kit *rest.Kit, chartInfo metadata.ChartConfig) (interface{}, error) {
	switch chartInfo.ReportType {
	case common.BizModuleHostChart:
		data, err := lgc.GetBizHostCount(kit)
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
		result, err := lgc.CoreAPI.CoreService().Operation().SearchTimerChartData(kit.Ctx, kit.Header, chartInfo)
		if err != nil {
			blog.Error("search chart data fail, chart name: %v, err: %v, rid: %v", chartInfo.Name, err, kit.Rid)
			return nil, err
		}
		return result.Data, nil
	}
}

func (lgc *Logics) TimerFreshData(ctx context.Context) {
	lgc.CheckTableExist(ctx)

	c := cron.New()
	spec := lgc.timerSpec // 从配置文件读取的时间
	_, err := c.AddFunc(spec, func() {
		blog.V(3).Infof("begin statistic chart data, time: %v", time.Now())
		// 主服务器跑定时
		opt := mapstr.MapStr{}
		isMaster := lgc.Engine.ServiceManageInterface.IsMaster()
		if isMaster {
			if _, err := lgc.CoreAPI.CoreService().Operation().TimerFreshData(ctx, lgc.header, opt); err != nil {
				blog.Error("statistic chart data fail, err: %v", err)
			}
		}
	})

	if err != nil {
		blog.Errorf("new cron failed, please contact developer, err: %v", err)
		return
	}
	c.Start()

	select {
	case <-ctx.Done():
		return
	}
}

// CheckTableExist 检测cc_chartData集合是否存在
func (lgc *Logics) CheckTableExist(ctx context.Context) {
	opt := mapstr.MapStr{}
	for {
		resp, err := lgc.CoreAPI.CoreService().Operation().TimerFreshData(ctx, lgc.header, opt)
		if err != nil {
			blog.Error("statistic chart data fail, err: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		if resp.Data {
			blog.V(3).Info("collection cc_ChartData inited")
			break
		}
		time.Sleep(10 * time.Second)
		blog.V(3).Info("waiting collection cc_ChartData init")
	}
	return
}
