package service

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"

	"github.com/gin-gonic/gin"
)

func (s *Service) ExportOperationChart(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)
	webCommon.SetProxyHeader(c)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)
	header := c.Request.Header

	value, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error:%s", err.Error())
		return
	}
	configID, err := parseConfigID(value)
	if err != nil {
		blog.Errorf("[export operation chart]unmarshal configID failed, error:%s, rid:%s", err.Error(), rid)
		msg := getReturnStr(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}
	chartInfo, err := s.Logics.GetOperationChart(ctx, header, configID)
	if err != nil {
		blog.Errorf("[export operation chart]get operation chart failed, error:%s, rid:%s", err.Error(), rid)
		msg := getReturnStr(common.CCErrWebGetObjectFail, defErr.Errorf(common.CCErrWebGetObjectFail, err.Error()).Error(), nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}
	if chartInfo == nil {
		blog.Errorf("[export operation chart]got none chartInfo, rid:%s", rid)
		c.String(http.StatusInternalServerError, "")
		return
	}
	var chartType string

	for _, i := range []string{"inst", "host", "nav"} {
		if chartInfo[i] != nil {
			chartType = chartInfo[i][0].ReportType
			break
		}
	}

	instInfo, err := s.Logics.GetOperationChartData(ctx, header, configID)
	if err != nil {
		blog.Errorf("[export operation chart]get operation chart data failed, error:%s, rid:%s", err.Error(), rid)
		msg := getReturnStr(common.CCErrWebGetObjectFail, defErr.Errorf(common.CCErrWebGetObjectFail, err.Error()).Error(), nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	chartData, err := parseOperationChartData(chartType, instInfo)
	if err != nil {
		blog.Errorf("[export operation chart]parse operation chart instance info failed, error:%s, rid:%s", err.Error(), rid)
		msg := getReturnStr(common.CCErrCommJSONUnmarshalFailed, err.Error(), nil)
		c.String(http.StatusInternalServerError, msg)
	}

	xlsxFile, err := s.Logics.CreateExcelFile(ctx, header, chartType, chartData)
	if err != nil {
		blog.Errorf("[export operation chart]create excel file failed, error:%s, rid:%s", err.Error(), rid)
		msg := fmt.Sprintf("create excel file failed, error:%v", err)
		c.String(http.StatusInternalServerError, msg)
	}

	logics.SendChartExcel(c, xlsxFile)
}
