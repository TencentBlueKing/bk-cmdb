package core

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (o *Operation) InitInnerChart(ctx context.Context) {
	header := make(http.Header, 0)
	if "" == util.GetOwnerID(header) {
		header.Set(common.BKHTTPOwnerID, common.BKSuperOwnerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	}

	opt := mapstr.MapStr{}
	result, err := o.CoreAPI.CoreService().Operation().SearchOperationChart(ctx, header, opt)
	if err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return
	}

	if result.Data.Count > 0 {
		return
	}

	for _, chart := range InnerCharts {
		_, err := o.CoreAPI.CoreService().Operation().CreateOperationChart(ctx, header, chart)
		if err != nil {
			blog.Errorf("init inner chart fail, err: %v", err)
		}
	}
}

var (
	BizModuleHostChart = metadata.ChartConfig{
		ReportType: "biz_module_host_chart",
	}

	HostOsChart = metadata.ChartConfig{
		ReportType: "host_os_chart",
		Name:       "按操作系统类型统计",
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
