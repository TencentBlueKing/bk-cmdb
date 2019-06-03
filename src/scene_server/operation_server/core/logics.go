package core

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (o *Operation) GetBizModuleHostCount(params ContextParams, data interface{}) (interface{}, error) {
	cond := metadata.QueryCondition{}
	info := make(map[string]interface{})
	target := [3]string{common.BKInnerObjIDApp, common.BKInnerObjIDModule, common.BKInnerObjIDHost}

	for _, obj := range target {
		result, err := o.CoreAPI.CoreService().Instance().ReadInstance(params.Context, params.Header, obj, &cond)
		if err != nil {
			blog.Errorf("search %v amount failed, err: %v", obj, err)
			return nil, params.Error.Error(common.CCErrOperationBizModuleHostAmountFail)
		}
		info[obj] = result.Data.Count
	}

	return info, nil
}

func (o *Operation) GetModelAndInstCount(params ContextParams, data interface{}) (interface{}, error) {
	cond := &metadata.QueryCondition{}
	result, err := o.CoreAPI.CoreService().Model().ReadModel(params.Context, params.Header, cond)
	if err != nil {
		blog.Errorf("search model fail , err: %v", err)
		return nil, err
	}

	info := mapstr.MapStr{}
	info["model"] = result.Data.Count

	opt := make(map[string]interface{})
	resp, err := o.CoreAPI.CoreService().Operation().SearchInstCount(params.Context, params.Header, opt)
	if err != nil {
		blog.Errorf("get instance number fail, err: %v", err)
		return nil, err
	}
	info["inst"] = resp.Data

	return info, nil
}

func (o *Operation) BizHostCount(params ContextParams) (interface{}, error) {
	cond := &metadata.QueryCondition{}
	bizInfo, err := o.CoreAPI.CoreService().Instance().ReadInstance(params.Context, params.Header, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("search biz info failed, err: %v", err)
		return nil, err
	}

	opt := mapstr.MapStr{}
	result, err := o.CoreAPI.CoreService().Operation().AggregateBizHost(params.Context, params.Header, opt)
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

func (o *Operation) CreateInnerChart(params ContextParams, chartInfo *metadata.ChartConfig) (interface{}, error) {
	opt := InnerCharts[chartInfo.ReportType]
	result, err := o.CoreAPI.CoreService().Operation().CreateOperationChart(params.Context, params.Header, opt)
	if err != nil {
		blog.Errorf("search chart info fail, err: %v", err)
		return nil, err
	}

	return result, nil
}

func (o *Operation) CommonStatisticFunc(params ContextParams, option metadata.ChartOption) (interface{}, error) {
	result, err := o.CoreAPI.CoreService().Operation().CommonAggregate(params.Context, params.Header, option)
	if err != nil {
		blog.Errorf("search data fail, err: %v", err)
		return nil, err
	}

	return result, nil
}
