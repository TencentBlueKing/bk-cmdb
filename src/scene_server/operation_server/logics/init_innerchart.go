package logics

import (
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
		ReportType: "biz_module_host_chart",
	}

	HostOsChart = metadata.ChartConfig{
		ReportType: "host_os_chart",
		Name:       "按操作系统类型统计",
		ObjID:      "host",
		Option: metadata.ChartOption{
			ChartType: "pie",
			Field:     "bk_os_type",
		},
	}

	HostBizChart = metadata.ChartConfig{
		ReportType: "host_biz_chart",
		Name:       "按业务统计",
	}

	HostCloudChart = metadata.ChartConfig{
		ReportType: "host_cloud_chart",
		Name:       "按云区域统计",
	}

	HostChangeBizChart = metadata.ChartConfig{
		ReportType: "host_change_biz_chart",
		Name:       "主机数量变化趋势",
	}

	ModelAndInstCountChart = metadata.ChartConfig{
		ReportType: "model_and_inst_count",
	}

	ModelInstChart = metadata.ChartConfig{
		ReportType: "model_inst_chart",
		Name:       "实例数量统计",
	}

	ModelInstChangeChart = metadata.ChartConfig{
		ReportType: "model_inst_change_chart",
		Name:       "实例变更统计",
	}

	InnerCharts = map[string]metadata.ChartConfig{
		"biz_module_host_chart":   BizModuleHostChart,
		"host_os_chart":           HostOsChart,
		"host_biz_chart":          HostBizChart,
		"host_cloud_chart":        HostCloudChart,
		"host_change_biz_chart":   HostChangeBizChart,
		"model_inst_chart":        ModelInstChart,
		"model_inst_change_chart": ModelInstChangeChart,
		"model_and_inst_count":    ModelAndInstCountChart,
	}
)
