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
 
package controllers

import (
	"configcenter/src/common"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/wactions"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/types"
	"configcenter/src/web_server/application/logics"
	"encoding/json"
	"fmt"
	_ "io"
	"math/rand"

	webCommon "configcenter/src/web_server/common"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

var (
	CODE_SUCESS            = 0
	CODE_ERROR_UPLOAD_FILE = 100
	CODE_ERROR_OPEN_FILE   = 101
)

func init() {
	wactions.RegisterNewAction(wactions.Action{common.HTTPCreate, "/hosts/import", nil, ImportHost})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectPost, "/hosts/export", nil, ExportHost})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectGet, "/importtemplate/:bk_obj_id", nil, BuildDownLoadExcelTemplate})
}

// ImportHost import host
func ImportHost(c *gin.Context) {
	file, err := c.FormFile("file")
	if nil != err {
		msg := getReturnStr(CODE_ERROR_UPLOAD_FILE, "未找到上传文件", nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	logics.SetProxyHeader(c)

	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/import/"
	_, err = os.Stat(dir)
	if nil != err {
		os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	}
	filePath := fmt.Sprintf("%s/importhost-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	err = c.SaveUploadedFile(file, filePath)
	if nil != err {
		msg := getReturnStr(CODE_ERROR_UPLOAD_FILE, fmt.Sprintf("保存文件失败;error:%s", err.Error()), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	defer os.Remove(filePath) //del file
	f, err := xlsx.OpenFile(filePath)
	if nil != err {
		msg := getReturnStr(CODE_ERROR_OPEN_FILE, fmt.Sprintf("保存上传文件;error:%s", err.Error()), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	cc := api.NewAPIResource()
	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	hosts, err := logics.GetImportHosts(f, apiSite, c.Request.Header)
	if 0 == len(hosts) {
		msg := getReturnStr(CODE_ERROR_OPEN_FILE, "文件内容不能为空", nil)
		if nil != err {
			msg = getReturnStr(CODE_ERROR_OPEN_FILE, "文件内容不能为空;"+err.Error(), nil)
		}
		c.String(http.StatusOK, string(msg))
		return
	}
	if nil != err {
		msg := getReturnStr(CODE_ERROR_OPEN_FILE, err.Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	url := apiSite + fmt.Sprintf("/api/%s/hosts/add", webCommon.API_VERSION)
	params := make(map[string]interface{})
	params["host_info"] = hosts
	params["bk_supplier_id"] = common.BKDefaultSupplierID

	blog.Info("add host url: %v", url)
	blog.Info("add host content: %v", params)
	reply, err := httpRequest(url, params, c.Request.Header)
	blog.Info("add host result: %v", reply)

	if nil != err {
		c.String(http.StatusOK, err.Error())
	} else {
		c.String(http.StatusOK, reply)
	}

}

// ExportHost export host
func ExportHost(c *gin.Context) {
	cc := api.NewAPIResource()
	appIDStr := c.PostForm("bk_biz_id")
	hostIDStr := c.PostForm("bk_host_id")
	kvMap := make(map[string]string)

	logics.SetProxyHeader(c)

	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	hostInfo, err := logics.GetHostData(appIDStr, hostIDStr, apiSite, c.Request.Header, kvMap)
	if err != nil {
		blog.Error(err.Error())
		c.String(http.StatusBadGateway, "获取主机数据失败, %s", err.Error())
		return
	}
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("host")
	if err != nil {
		blog.Error(err.Error())
		c.String(http.StatusBadGateway, "创建EXCEL文件失败，%s", err.Error())
		return
	}
	row = sheet.AddRow()
	kArray := make([]string, 0)

	for i, k := range kvMap {
		cell = row.AddCell()
		cell.Value = k
		kArray = append(kArray, i)
	}
	kLength := len(kArray)
	for _, j := range hostInfo {
		hostData := j.(map[string]interface{})
		hostcell := hostData["host"].(map[string]interface{})
		row = sheet.AddRow()
		for i := 0; i != kLength; i++ {
			cell = row.AddCell()
			kName := kArray[i]

			n, ok := hostcell[kName]
			if ok {
				if nil == n {
					cell.Value = ""
					continue
				}
				objtype := reflect.TypeOf(n)
				switch objtype.Kind() {
				case reflect.String:
					cell.Value = reflect.ValueOf(n).String()
				default:
					cell.Value = ""
				}
			} else {
				cell.Value = ""
			}
		}
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if nil != err {
		os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm)
	}
	fileName := fmt.Sprintf("%dhost.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)

	err = file.Save(dirFileName)
	if err != nil {
		blog.Error("ExportHost save file error:%s", err.Error())
		fmt.Printf(err.Error())
	}
	logics.AddDownExcelHttpHeader(c, "host.xlsx")
	c.File(dirFileName)

	os.Remove(dirFileName)

}

//BuildDownLoadExcelTemplate build download excel template
func BuildDownLoadExcelTemplate(c *gin.Context) {
	logics.SetProxyHeader(c)
	objID := c.Param(common.BKObjIDField)
	cc := api.NewAPIResource()
	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	url := apiSite + fmt.Sprintf("/api/%s/object/attr/search", webCommon.API_VERSION)
	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/template/"
	_, err := os.Stat(dir)
	if nil != err {
		os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	}
	file := fmt.Sprintf("%s/%stemplate-%d-%d.xlsx", dir, objID, time.Now().UnixNano(), randNum)
	err = logics.BuildExcelTemplate(url, objID, file, c.Request.Header)
	if nil != err {
		blog.Errorf("BuildDownLoadExcelTemplate object:%s error:%s", objID, err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, fmt.Sprintf("未找到下载模板%s", objID), nil)
		c.Writer.Write([]byte(reply))
		return
	}

	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("template_%s.xlsx", objID))

	//http.ServeFile(c.Writer, c.Request, file)
	c.File(file)
	os.Remove(file)
	return
}

//httpRequest do http request
func httpRequest(url string, body interface{}, header http.Header) (string, error) {
	params, _ := json.Marshal(body)
	blog.Info("input:%s", string(params))
	httpClient := httpclient.NewHttpClient()
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetHeader("Accept", "application/json")

	reply, err := httpClient.POST(url, header, params)

	return string(reply), err
}

// getReturnStr get return string
func getReturnStr(code int, message string, data interface{}) string {
	ret := make(map[string]interface{})
	ret["bk_error_code"] = code
	if 0 == code {
		ret["result"] = true
	} else {
		ret["result"] = false
	}
	ret["bk_error_msg"] = message
	ret["data"] = data
	msg, _ := json.Marshal(ret)

	return string(msg)

}
