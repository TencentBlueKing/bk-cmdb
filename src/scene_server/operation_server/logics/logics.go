package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"context"
	"time"
)

func (lgc *Logics) GetBizModuleHostCount(kit *rest.Kit) (mapstr.MapStr, error) {
	cond := metadata.QueryCondition{}
	info := mapstr.MapStr{}
	target := [3]string{common.BKInnerObjIDApp, common.BKInnerObjIDModule, common.BKInnerObjIDHost}

	for _, obj := range target {
		result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, obj, &cond)
		if err != nil {
			blog.Errorf("search %v amount failed, err: %v", obj, err)
			return nil, kit.CCError.Error(common.CCErrOperationBizModuleHostAmountFail)
		}
		info[obj] = result.Data.Count
	}

	return info, nil
}

func (lgc *Logics) GetModelAndInstCount(kit *rest.Kit) (mapstr.MapStr, error) {
	cond := &metadata.QueryCondition{}
	result, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("search model fail , err: %v", err)
		return nil, err
	}

	info := mapstr.MapStr{}
	info["model"] = result.Data.Count

	opt := make(map[string]interface{})
	resp, err := lgc.CoreAPI.CoreService().Operation().SearchInstCount(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("get instance number fail, err: %v", err)
		return nil, err
	}
	info["inst"] = resp.Data

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

	_, err := lgc.CoreAPI.CoreService().Operation().TimerFreshData(ctx, lgc.header, opt)
	if err != nil {
		blog.Error("start collect chart data timer fail, err: %v", err)
	}

	timer := time.NewTicker(time.Duration(12) * time.Hour)
	for range timer.C {
		_, err := lgc.CoreAPI.CoreService().Operation().TimerFreshData(ctx, lgc.header, opt)
		if err != nil {
			blog.Error("start collect chart data timer fail, err: %v", err)
		}
	}
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
		result, err := lgc.CoreAPI.CoreService().Operation().SearchOperationChartData(kit.Ctx, kit.Header, chartInfo.ReportType)
		if err != nil {
			blog.Error("search chart data fail, chart name: %v, err: %v", chartInfo.Name, err)
			return nil, err
		}
		return result.Data, nil
	}
}
