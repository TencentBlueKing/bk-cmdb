package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
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

func (lgc *Logics) GetInnerChartData(kit *rest.Kit, chartInfo metadata.ChartConfig) (interface{}, error) {
	switch chartInfo.ReportType {
	case "biz_module_host_chart":
		data, err := lgc.GetBizModuleHostCount(kit)
		if err != nil {
			return nil, err
		}
		return data, nil
	case "model_and_inst_count":
		data, err := lgc.GetModelAndInstCount(kit)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		result, err := lgc.CoreAPI.CoreService().Operation().SearchOperationChartData(kit.Ctx, kit.Header, chartInfo.ReportType)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (lgc *Logics) BizHostCount(kit *rest.Kit) (interface{}, error) {
	cond := &metadata.QueryCondition{}
	bizInfo, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("search biz info failed, err: %v", err)
		return nil, err
	}

	opt := mapstr.MapStr{}
	result, err := lgc.CoreAPI.CoreService().Operation().AggregateBizHost(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("search biz's host count fail, err: %v", err)
		return nil, err
	}

	bizHost := mapstr.MapStr{}

	for _, info := range bizInfo.Data.Info {
		for _, data := range result.Data {
			if info[common.BKAppIDField] == data.Id {
				bizName, err := info.String(common.BKAppNameField)
				if err != nil {
					blog.Errorf("interface convert to string fail, err: %v", err)
					continue
				}
				bizHost[bizName] = data.Count
			}
		}
	}

	return bizHost, nil
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

func (lgc *Logics) CommonStatisticFunc(kit *rest.Kit, option metadata.ChartConfig) (interface{}, error) {
	result, err := lgc.CoreAPI.CoreService().Operation().CommonAggregate(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("search data fail, err: %v", err)
		return nil, err
	}

	return result, nil
}
