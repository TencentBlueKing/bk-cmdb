package logics

import (
	"context"

	"github.com/robfig/cron"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) GetBizModuleHostCount(kit *rest.Kit) ([]metadata.IDStringCountInt64, error) {
	cond := metadata.QueryCondition{}
	data := make([]metadata.IDStringCountInt64, 0)
	target := [3]string{common.BKInnerObjIDApp, common.BKInnerObjIDModule, common.BKInnerObjIDHost}

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

func (lgc *Logics) CreateInnerChart(kit *rest.Kit, chartInfo *metadata.ChartConfig) (interface{}, error) {
	opt, ok := InnerCharts[chartInfo.ReportType]
	if !ok {
		return nil, kit.CCError.Error(common.CCErrOperationNewAddStatisticFail)
	}
	blog.Debug("opt: %v", opt)
	result, err := lgc.CoreAPI.CoreService().Operation().CreateOperationChart(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("search chart info fail, err: %v", err)
		return nil, err
	}

	return result.Data, nil
}

func (lgc *Logics) TimerFreshData(ctx context.Context) {
	opt := mapstr.MapStr{}

	// 主服务器跑定时
	if isMaster := lgc.Engine.ServiceManageInterface.IsMaster(); !isMaster {
		return
	}

	if _, err := lgc.CoreAPI.CoreService().Operation().TimerFreshData(ctx, lgc.header, opt); err != nil {
		blog.Error("start collect chart data timer fail, err: %v", err)
		return
	}

	c := cron.New()
	spec := "0 0 2 * * ?" // 每天凌晨两点，更新定时统计图表数据
	err := c.AddFunc(spec, func() {
		if _, err := lgc.CoreAPI.CoreService().Operation().TimerFreshData(ctx, lgc.header, opt); err != nil {
			blog.Error("start collect chart data timer fail, err: %v", err)
			return
		}
	})

	if err != nil {
		blog.Error("start collect chart data timer fail, err: %v", err)
		return
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
