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
	"configcenter/src/common/util"
	"configcenter/src/web_server/application/logics"
	webCommon "configcenter/src/web_server/common"
)

func init() {
	wactions.RegisterNewAction(wactions.Action{common.HTTPCreate, "/netdevice/import", nil, ImportNetDevice})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectPost, "/netdevice/export", nil, ExportNetDevice})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectGet, "/netcollect/importtemplate/netdevice", nil, BuildDownLoadNetDeviceExcelTemplate})
}

func ImportNetDevice(c *gin.Context) {
	cc := api.NewAPIResource()
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	file, err := c.FormFile("file")
	if nil != err {
		blog.Errorf("Import Net Device get file error:%s", err.Error())
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
	filePath := fmt.Sprintf("%s/importnetdevice-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
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
	netdevice, errMsg, err := logics.GetImportNetDevices(f, apiSite, c.Request.Header, defLang) //TODO

	if nil != err {
		blog.Errorf("ImportNetDevice logID:%s, error:%s", util.GetHTTPCCRequestID(c.Request.Header), err.Error())
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 != len(errMsg) {
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, " file empty").Error(), common.KvMap{"err": errMsg})
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 == len(netdevice) {
		msg := getReturnStr(common.CCErrWebFileContentEmpty, defErr.Errorf(common.CCErrWebFileContentEmpty, "").Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	url := apiSite + fmt.Sprintf("/api/%s/netcollect/device/action/create", webCommon.API_VERSION)
	params := make(map[string]interface{})
	params["host_info"] = netdevice //TODO
	params["bk_supplier_id"] = common.BKDefaultSupplierID
	params["input_type"] = common.InputTypeExcel

	blog.Infof("add host url: %v", url)
	blog.Infof("add host content: %v", params)
	reply, err := deviceHttpRequest(url, params, c.Request.Header) //TODO
	blog.Infof("add host result: %v", reply)

	if nil != err {
		c.String(http.StatusOK, err.Error())
	} else {
		c.String(http.StatusOK, reply)
	}

}

func ExportNetDevice(c *gin.Context) {
	cc := api.NewAPIResource()
	appIDStr := c.PostForm("bk_biz_id")
	hostIDStr := c.PostForm("bk_host_id")

	logics.SetProxyHeader(c)

	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	hostInfo, err := logics.GetHostData(appIDStr, hostIDStr, apiSite, c.Request.Header) //TODO
	if err != nil {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebGetHostFail, defErr.Errorf(common.CCErrWebGetHostFail, err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg, nil)
		return
	}
	var file *xlsx.File
	var sheet *xlsx.Sheet

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("netdevice")

	objID := common.BKInnerObjIDHost
	fields, err := logics.GetObjFieldIDs(objID, apiSite, logics.GetFilterFields(objID), c.Request.Header) //TODO
	if nil != err {
		blog.Errorf("ExportNetDevice get %s field error:%s error:%s", objID, err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}
	err = logics.BuildHostExcelFromData(objID, fields, nil, hostInfo, sheet, defLang) //TODO
	if nil != err {
		blog.Errorf("ExportNetDevice object:%s error:%s", objID, err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if nil != err {
		os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm)
	}
	fileName := fmt.Sprintf("%dnetdevice.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)

	logics.ProductExcelCommentSheet(file, defLang)
	if err = file.Save(dirFileName); err != nil {
		blog.Error("ExportNetDevice save file error:%s", err.Error())
		reply := getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}
	logics.AddDownExcelHttpHeader(c, "netdevice.xlsx")
	c.File(dirFileName)

	os.Remove(dirFileName)
}

func BuildDownLoadNetDeviceExcelTemplate(c *gin.Context) {
	logics.SetProxyHeader(c)
	objID := c.Param(common.BKObjIDField)
	cc := api.NewAPIResource()
	apiSite := cc.APIAddr()
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

	if err = logics.BuildNetDeviceExcelTemplate(apiSite, objID, file, c.Request.Header, defLang); nil != err {
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

func deviceHttpRequest(url string, body interface{}, header http.Header) (string, error) {
	params, _ := json.Marshal(body)
	blog.Info("input:%s", string(params))
	httpClient := httpclient.NewHttpClient()
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetHeader("Accept", "application/json")

	reply, err := httpClient.POST(url, header, params)

	return string(reply), err
}

// getReturnStr get return string
func getDeviceReturnStr(code int, message string, data interface{}) string {
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
