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
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/wactions"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/types"
	"configcenter/src/web_server/application/logics"
	webCommon "configcenter/src/web_server/common"
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
	cc := api.NewAPIResource()
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	file, err := c.FormFile("file")
	if nil != err {
		blog.Errorf("Import Host get file error:%s", err.Error())
		msg := getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
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
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	defer os.Remove(filePath) //del file
	f, err := xlsx.OpenFile(filePath)
	if nil != err {
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	hosts, err := logics.GetImportHosts(f, apiSite, c.Request.Header, defLang)

	if nil != err {
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 == len(hosts) {
		msg := getReturnStr(common.CCErrWebFileContentEmpty, defErr.Errorf(common.CCErrWebFileContentEmpty, "").Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	url := apiSite + fmt.Sprintf("/api/%s/hosts/add", webCommon.API_VERSION)
	params := make(map[string]interface{})
	params["host_info"] = hosts
	params["bk_supplier_id"] = common.BKDefaultSupplierID
	params["input_type"] = common.InputTypeExcel

	blog.Infof("add host url: %v", url)
	blog.Infof("add host content: %v", params)
	reply, err := httpRequest(url, params, c.Request.Header)
	blog.Infof("add host result: %v", reply)

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

	logics.SetProxyHeader(c)

	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	hostInfo, err := logics.GetHostData(appIDStr, hostIDStr, apiSite, c.Request.Header)
	if err != nil {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebGetHostFail, defErr.Errorf(common.CCErrWebGetHostFail, err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg, nil)
		return
	}
	var file *xlsx.File
	var sheet *xlsx.Sheet

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("host")

	objID := common.BKInnerObjIDHost
	fields, err := logics.GetObjFieldIDs(objID, apiSite, logics.GetFilterFields(objID), c.Request.Header)
	if nil != err {
		blog.Errorf("ExportHost get %s field error:%s error:%s", objID, err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}
	err = logics.BuildHostExcelFromData(objID, fields, nil, hostInfo, sheet, defLang)
	if nil != err {
		blog.Errorf("ExportHost object:%s error:%s", objID, err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if nil != err {
		os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm)
	}
	fileName := fmt.Sprintf("%dhost.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)

	logics.ProductExcelCommentSheet(file, defLang)
	err = file.Save(dirFileName)
	if err != nil {
		blog.Error("ExportHost save file error:%s", err.Error())
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
	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/template/"
	_, err := os.Stat(dir)
	if nil != err {
		os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	}
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	file := fmt.Sprintf("%s/%stemplate-%d-%d.xlsx", dir, objID, time.Now().UnixNano(), randNum)
	err = logics.BuildExcelTemplate(apiSite, objID, file, c.Request.Header, defLang)
	if nil != err {
		blog.Errorf("BuildDownLoadExcelTemplate object:%s error:%s", objID, err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
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
