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
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/web_server/application/logics"
	webCommon "configcenter/src/web_server/common"
)

func init() {
	wactions.RegisterNewAction(wactions.Action{common.HTTPCreate, "/netdevice/import", nil, ImportNetDevice})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectPost, "/netdevice/export", nil, ExportNetDevice})
	wactions.RegisterNewAction(wactions.Action{
		common.HTTPSelectGet, "/netcollect/importtemplate/netdevice", nil, BuildDownLoadNetDeviceExcelTemplate})
}

func ImportNetDevice(c *gin.Context) {
	cc := api.NewAPIResource()
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	fileHeader, err := c.FormFile("file")
	if nil != err {
		blog.Errorf("Import Net Device get file error:%s", err.Error())
		msg := getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	logics.SetProxyHeader(c)

	dir := webCommon.ResourcePath + "/import/"
	if _, err = os.Stat(dir); nil != err {
		os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	}

	randNum := rand.Uint32()
	filePath := fmt.Sprintf("%s/importnetdevice-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	if err = c.SaveUploadedFile(fileHeader, filePath); nil != err {
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	defer os.Remove(filePath) //del file

	file, err := xlsx.OpenFile(filePath)
	if nil != err {
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	netdevice, errMsg, err := logics.GetImportNetDevices(c.Request.Header, defLang, file, apiSite)
	if nil != err {
		blog.Errorf("ImportNetDevice logID:%s, error:%s", util.GetHTTPCCRequestID(c.Request.Header), err.Error())
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 != len(errMsg) {
		msg := getReturnStr(common.CCErrWebFileContentFail,
			defErr.Errorf(common.CCErrWebFileContentFail, " file empty").Error(),
			common.KvMap{"err": errMsg})
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 == len(netdevice) {
		msg := getReturnStr(common.CCErrWebFileContentEmpty, defErr.Errorf(common.CCErrWebFileContentEmpty, "").Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	url := apiSite + fmt.Sprintf("/api/%s/netcollect/device/action/create", webCommon.API_VERSION)
	blog.Infof("add device url: %v", url)

	params := make([]interface{}, 0)
	for _, value := range netdevice {
		params = append(params, value)
	}
	blog.Infof("add device content: %v", params)

	reply, err := httpRequest(url, params, c.Request.Header)
	blog.Infof("add device result: %v", reply)

	if nil != err {
		c.String(http.StatusOK, err.Error())
		return
	}

	c.String(http.StatusOK, reply)
}

func ExportNetDevice(c *gin.Context) {
	logics.SetProxyHeader(c)

	cc := api.NewAPIResource()
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	deviceIDstr := c.PostForm(common.BKDeviceIDField)
	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)

	deviceInfo, err := logics.GetNetDeviceData(c.Request.Header, apiSite, deviceIDstr)
	if nil != err {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebGetHostFail, defErr.Errorf(common.CCErrWebGetHostFail, err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg, nil)
		return
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet(common.BKNetDevice)
	if nil != err {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebCreateEXCELFail,
			defErr.Errorf(common.CCErrWebCreateEXCELFail,
				err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg, nil)
		return
	}

	fields := logics.GetNetDevicefield(defLang)
	logics.AddNetDeviceExtFields(&fields, defLang)

	if err = logics.BuildNetDeviceExcelFromData(defLang, fields, deviceInfo, sheet); nil != err {
		blog.Errorf("ExportNetDevice object:%s error:%s", err.Error())
		reply := getReturnStr(
			common.CCErrCommExcelTemplateFailed,
			defErr.Errorf(common.CCErrCommExcelTemplateFailed, common.BKNetDevice).Error(),
			nil)
		c.Writer.Write([]byte(reply))
		return
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	if _, err = os.Stat(dirFileName); nil != err {
		os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm)
	}

	fileName := fmt.Sprintf("%dnetdevice.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)

	logics.ProductExcelCommentSheet(file, defLang)

	if err = file.Save(dirFileName); nil != err {
		blog.Error("ExportNetDevice save file error:%s", err.Error())
		reply := getReturnStr(common.CCErrWebCreateEXCELFail,
			defErr.Errorf(common.CCErrCommExcelTemplateFailed, err.Error()).Error(),
			nil)
		c.Writer.Write([]byte(reply))
		return
	}

	logics.AddDownExcelHttpHeader(c, "netdevice.xlsx")
	c.File(dirFileName)

	os.Remove(dirFileName)
}

func BuildDownLoadNetDeviceExcelTemplate(c *gin.Context) {
	logics.SetProxyHeader(c)
	cc := api.NewAPIResource()

	dir := webCommon.ResourcePath + "/template/"
	if _, err := os.Stat(dir); nil != err {
		os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	}

	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	randNum := rand.Uint32()
	file := fmt.Sprintf("%s/%stemplate-%d-%d.xlsx", dir, common.BKNetDevice, time.Now().UnixNano(), randNum)

	apiSite := cc.APIAddr()
	if err := logics.BuildNetDeviceExcelTemplate(c.Request.Header, defLang, apiSite, file); nil != err {
		blog.Errorf("Build NetDevice Excel Template, error:%s", err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed,
			defErr.Errorf(common.CCErrCommExcelTemplateFailed, common.BKNetDevice).Error(),
			nil)
		c.Writer.Write([]byte(reply))
		return
	}

	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("template_%s.xlsx", common.BKNetDevice))

	c.File(file)
	os.Remove(file)
	return
}
