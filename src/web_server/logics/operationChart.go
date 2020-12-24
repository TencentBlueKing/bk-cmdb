package logics

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"
)

func (lgc *Logics) GetOperationChartData(ctx context.Context, header http.Header, configID uint64) (interface{}, error) {
	cond := mapstr.MapStr{
		"config_id": configID,
	}

	resp, err := lgc.Engine.CoreAPI.OperationServer().SearchChartData(ctx, header, cond)
	if nil != err {
		return nil, fmt.Errorf("get operation chart data failed, err: %v", err)
	}

	if !resp.Result {
		return nil, fmt.Errorf("get operation chart data result failed, err: %s", resp.ErrMsg)
	}

	return resp.Data, nil
}

func (lgc *Logics) GetOperationChart(ctx context.Context, header http.Header, configID uint64) (map[string][]metadata.ChartConfig, error) {
	cond := mapstr.MapStr{
		"config_id": configID,
	}

	resp, err := lgc.Engine.CoreAPI.CoreService().Operation().SearchOperationCharts(ctx, header, cond)
	if nil != err {
		return nil, fmt.Errorf("get operation chart failed, err: %v", err)
	}

	if !resp.Result {
		return nil, fmt.Errorf("get operation chart result failed, err: %s", resp.ErrMsg)
	}

	return resp.Data.Info, nil
}

func (lgc *Logics) CreateExcelFile(ctx context.Context, header http.Header, chartType string, data interface{}) (*xlsx.File, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	rowIndex := 1

	xlsxFile := xlsx.NewFile()

	sheet, err := xlsxFile.AddSheet(chartType)
	if err != nil {
		blog.Errorf("CreateExcelFile add excel sheet error, err:%s, rid:%s", err.Error(), rid)
		return nil, err
	}

	if chartType == common.HostChangeBizChart {
		longestDateTarget := ""
		rowNum := 0
		columnNum := 0

		dataPointer, ok := data.(*map[string][]metadata.StringIDCount)
		if !ok {
			err = fmt.Errorf("CreateExcelFile format *map[string][]metadata.StringIDCount error, rid:%s", rid)
			blog.Errorf("CreateExcelFile format *map[string][]metadata.StringIDCount error, rid:%s", rid)
			return nil, err
		}
		formatData := *dataPointer

		// figure out which data to use as the header row
		num := 0
		for bkBizName, value := range formatData {
			if num < len(value) {
				num = len(value)
				longestDateTarget = bkBizName
			}
		}

		// make header row, and record the columnNum with stringIDCount.ID
		dateColumn := map[string]int{}
		cell := sheet.Cell(rowNum, columnNum)
		cell.SetString(`bk_biz_name\data`)

		for i, stringIDCount := range formatData[longestDateTarget] {
			columnNum = i + 1
			cell = sheet.Cell(rowNum, columnNum)
			cell.SetString(stringIDCount.ID)
			dateColumn[stringIDCount.ID] = columnNum
		}

		// Fill in the sheet
		rowNum = 1
		columnNum = 0
		for bkBizName, value := range formatData {
			cell = sheet.Cell(rowNum, 0)
			cell.SetString(bkBizName)
			for _, stringIDCount := range value {
				columnNum, ok := dateColumn[stringIDCount.ID]
				if ok {
					cell = sheet.Cell(rowNum, columnNum)
					cell.SetInt64(stringIDCount.Count)
				}
			}
			rowNum += 1
		}

	} else if chartType == common.ModelInstChangeChart {
		formatData, ok := data.(*metadata.ModelInstChange)
		if !ok {
			err = fmt.Errorf("CreateExcelFile format *metadata.ModelInstChange error, rid:%s", rid)
			blog.Errorf("CreateExcelFile format *metadata.ModelInstChange error, rid:%s", rid)
			return nil, err
		}
		cell := sheet.Cell(0, 0)
		cell.SetString("bk_biz_name")
		cell = sheet.Cell(0, 1)
		cell.SetString("create")
		cell = sheet.Cell(0, 2)
		cell.SetString("update")
		cell = sheet.Cell(0, 3)
		cell.SetString("delete")
		for bkObjName, value := range *formatData {
			cell = sheet.Cell(rowIndex, 0)
			cell.SetString(bkObjName)
			cell = sheet.Cell(rowIndex, 1)
			cell.SetInt64(value.Create)
			cell = sheet.Cell(rowIndex, 2)
			cell.SetInt64(value.Update)
			cell = sheet.Cell(rowIndex, 3)
			cell.SetInt64(value.Delete)
			rowIndex += 1
		}
	} else {
		formatData, ok := data.(*[]metadata.StringIDCount)
		if !ok {
			err = fmt.Errorf("CreateExcelFile format *[]metadata.StringIDCount error, rid:%s", rid)
			blog.Errorf("CreateExcelFile format *[]metadata.StringIDCount error, rid:%s", rid)
			return nil, err
		}
		cell := sheet.Cell(0, 0)
		cell.SetString("bk_obj_id")
		cell = sheet.Cell(0, 1)
		cell.SetString("count")
		for _, inst := range *formatData {
			cell = sheet.Cell(rowIndex, 0)
			cell.SetString(inst.ID)
			cell = sheet.Cell(rowIndex, 1)
			cell.SetInt64(inst.Count)
			rowIndex += 1
		}
	}

	return xlsxFile, nil
}

func SendChartExcel(c *gin.Context, xlsxFile *xlsx.File) {
	rid := util.ExtractRequestIDFromContext(c)
	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err := os.Stat(dirFileName)
	if nil != err {
		if err := os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("establish excel file failed, make local dir to save export file failed, err: %+v, rid: %s", err, rid)
			return
		}
	}
	fileName := fmt.Sprintf("%dchart.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)

	err = xlsxFile.Save(dirFileName)
	if err != nil {
		blog.Errorf("ExportChart failed, save file failed, err: %+v, rid: %s", err, rid)
		_, _ = c.Writer.Write([]byte(fmt.Sprintf("ExportChart failed, save file failed, err: %+v, rid: %s", err, rid)))
		return
	}
	AddDownExcelHttpHeader(c, "bk_cmdb_export_chart.xlsx")
	c.File(dirFileName)

	if err := os.Remove(dirFileName); err != nil {
		blog.Errorf("ExportHost success, but remove chart.xlsx file failed, err: %+v, rid: %s", err, rid)
	}
}
