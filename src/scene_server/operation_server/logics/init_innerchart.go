package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"context"
)

func (lgc *Logics) InitInnerChart(ctx context.Context) {
	opt := mapstr.MapStr{}
	result, err := lgc.CoreAPI.CoreService().Operation().SearchOperationChart(ctx, lgc.header, opt)
	if err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return
	}

	if result.Data.Count > 0 {
		return
	}

	for _, chart := range InnerCharts {
		_, err := lgc.CoreAPI.CoreService().Operation().CreateOperationChart(ctx, lgc.header, chart)
		if err != nil {
			blog.Errorf("init inner chart fail, err: %v", err)
		}
	}

	//todo 初始化图表位置信息
}

var (
	BizModuleHostChart = metadata.ChartConfig{
		ReportType: common.BizModuleHostChart,
	}

	HostOsChart = metadata.ChartConfig{
		ReportType: common.HostOSChart,
		Name:       "按操作系统类型统计",
		ObjID:      "host",
		Option: metadata.ChartOption{
			ChartType: "pie",
			Field:     "bk_os_type",
		},
	}

	HostBizChart = metadata.ChartConfig{
		ReportType: common.HostBizChart,
		Name:       "按业务统计",
	}

	HostCloudChart = metadata.ChartConfig{
		ReportType: common.HostCloudChart,
		Name:       "按云区域统计",
	}

	HostChangeBizChart = metadata.ChartConfig{
		ReportType: common.HostChangeBizChart,
		Name:       "主机数量变化趋势",
	}

	ModelAndInstCountChart = metadata.ChartConfig{
		ReportType: common.ModelAndInstCount,
	}

	ModelInstChart = metadata.ChartConfig{
		ReportType: common.ModelInstChart,
		Name:       "实例数量统计",
	}

	ModelInstChangeChart = metadata.ChartConfig{
		ReportType: common.ModelInstChangeChart,
		Name:       "实例变更统计",
	}

	InnerCharts = map[string]metadata.ChartConfig{
		common.BizModuleHostChart:   BizModuleHostChart,
		common.HostOSChart:          HostOsChart,
		common.HostBizChart:         HostBizChart,
		common.HostCloudChart:       HostCloudChart,
		common.HostChangeBizChart:   HostChangeBizChart,
		common.ModelInstChart:       ModelInstChart,
		common.ModelInstChangeChart: ModelInstChangeChart,
		common.ModelAndInstCount:    ModelAndInstCountChart,
	}
)
