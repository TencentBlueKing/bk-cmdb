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
	wactions.RegisterNewAction(wactions.Action{common.HTTPCreate, "/netproperty/import", nil, ImportNetProperty})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectPost, "/netproperty/export", nil, ExportNetProperty})
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectGet, "/netcollect/importtemplate/netproperty", nil, BuildDownLoadNetPropertyExcelTemplate})
}

func ImportNetProperty(c *gin.Context) {
	cc := api.NewAPIResource()
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	fileHeader, err := c.FormFile("file")
	if nil != err {
		blog.Errorf("Import Net Property get file error:%s", err.Error())
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
	filePath := fmt.Sprintf("%s/importnetproperty-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
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
	netproperty, errMsg, err := logics.GetImportNetProperty(c.Request.Header, defLang, file, apiSite) //TODO
	if nil != err {
		blog.Errorf("ImportNetProperty logID:%s, error:%s", util.GetHTTPCCRequestID(c.Request.Header), err.Error())
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 != len(errMsg) {
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, " file empty").Error(), common.KvMap{"err": errMsg})
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 == len(netproperty) {
		msg := getReturnStr(common.CCErrWebFileContentEmpty, defErr.Errorf(common.CCErrWebFileContentEmpty, "").Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	url := apiSite + fmt.Sprintf("/api/%s/netcollect/property/action/create", webCommon.API_VERSION)
	blog.Infof("add net property url: %v", url)

	params := make([]interface{}, 0)
	for _, value := range netproperty {
		params = append(params, value)
	}
	blog.Infof("add net property content: %v", params)

	reply, err := httpRequest(url, params, c.Request.Header)
	blog.Infof("add net property result: %v", reply)

	if nil != err {
		c.String(http.StatusOK, err.Error())
		return
	}

	c.String(http.StatusOK, reply)
}

func ExportNetProperty(c *gin.Context) {
	logics.SetProxyHeader(c)

	cc := api.NewAPIResource()
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := cc.Lang.CreateDefaultCCLanguageIf(language)
	defErr := cc.Error.CreateDefaultCCErrorIf(language)

	netPropertyIDstr := c.PostForm(common.BKNetcollectPropertyIDlField)
	apiSite, _ := cc.AddrSrv.GetServer(types.CC_MODULE_APISERVER)

	netPropertyInfo, err := logics.GetNetPropertyData(c.Request.Header, apiSite, netPropertyIDstr)
	if nil != err {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebGetHostFail, defErr.Errorf(common.CCErrWebGetHostFail, err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg, nil)
		return
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet(common.BKNetProperty)
	if nil != err {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebCreateEXCELFail,
			defErr.Errorf(common.CCErrWebCreateEXCELFail,
				err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg, nil)
		return
	}

	fields := logics.GetNetPropertyField(defLang)
	logics.AddNetPropertyExtFields(&fields, defLang)

	if err = logics.BuildNetPropertyExcelFromData(defLang, fields, netPropertyInfo, sheet); nil != err {
		blog.Errorf("ExportNetProperty object:%s error:%s", err.Error())
		reply := getReturnStr(
			common.CCErrCommExcelTemplateFailed,
			defErr.Errorf(common.CCErrCommExcelTemplateFailed, common.BKNetProperty).Error(),
			nil)
		c.Writer.Write([]byte(reply))
		return
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	if _, err = os.Stat(dirFileName); nil != err {
		os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm)
	}

	fileName := fmt.Sprintf("%dnetproperty.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)

	logics.ProductExcelCommentSheet(file, defLang)

	if err = file.Save(dirFileName); nil != err {
		blog.Error("ExportNetProperty save file error:%s", err.Error())
		reply := getReturnStr(common.CCErrWebCreateEXCELFail,
			defErr.Errorf(common.CCErrCommExcelTemplateFailed, err.Error()).Error(),
			nil)
		c.Writer.Write([]byte(reply))
		return
	}

	logics.AddDownExcelHttpHeader(c, "netproperty.xlsx")
	c.File(dirFileName)

	os.Remove(dirFileName)

}

func BuildDownLoadNetPropertyExcelTemplate(c *gin.Context) {
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
	file := fmt.Sprintf("%s/%stemplate-%d-%d.xlsx", dir, common.BKNetProperty, time.Now().UnixNano(), randNum)

	apiSite := cc.APIAddr()
	if err := logics.BuildNetPropertyExcelTemplate(c.Request.Header, defLang, apiSite, file); nil != err {
		blog.Errorf("Build NetProperty Excel Template, error:%s", err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed,
			defErr.Errorf(common.CCErrCommExcelTemplateFailed, common.BKNetProperty).Error(),
			nil)
		c.Writer.Write([]byte(reply))
		return
	}

	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("template_%s.xlsx", common.BKNetProperty))

	c.File(file)
	os.Remove(file)
	return
}
