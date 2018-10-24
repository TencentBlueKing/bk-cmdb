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

package service

import (
	"context"
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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"
)

var (
	CODE_SUCESS            = 0
	CODE_ERROR_UPLOAD_FILE = 100
	CODE_ERROR_OPEN_FILE   = 101
)

// ImportHost import host
func (s *Service) ImportHost(c *gin.Context) {

	language := logics.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)
	pheader := c.Request.Header
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

	hosts, errMsg, err := s.Logics.GetImportHosts(f, c.Request.Header, defLang)

	if nil != err {
		blog.Errorf("ImportHost logID:%s, error:%s", util.GetHTTPCCRequestID(c.Request.Header), err.Error())
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 != len(errMsg) {
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail, " file empty").Error(), common.KvMap{"err": errMsg})
		c.String(http.StatusOK, string(msg))
		return
	}
	if 0 == len(hosts) {
		msg := getReturnStr(common.CCErrWebFileContentEmpty, defErr.Errorf(common.CCErrWebFileContentEmpty, "").Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	params := mapstr.MapStr{}
	params["host_info"] = hosts
	params["bk_supplier_id"] = common.BKDefaultSupplierID
	params["input_type"] = common.InputTypeExcel

	result, err := s.CoreAPI.ApiServer().AddHost(context.Background(), pheader, params)

	if nil != err {
		msg := getReturnStr(common.CCErrCommHTTPDoRequestFailed, defErr.Errorf(common.CCErrCommHTTPDoRequestFailed, "").Error(), nil)
		c.String(http.StatusOK, string(msg))
	}

	c.JSON(http.StatusOK, result)

}

// ExportHost export host
func (s *Service) ExportHost(c *gin.Context) {

	appIDStr := c.PostForm("bk_biz_id")
	hostIDStr := c.PostForm("bk_host_id")

	logics.SetProxyHeader(c)
	pheader := c.Request.Header
	defLang := s.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	hostInfo, err := s.Logics.GetHostData(appIDStr, hostIDStr, pheader)
	if err != nil {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebGetHostFail, defErr.Errorf(common.CCErrWebGetHostFail, err.Error()).Error(), nil)
		c.String(http.StatusInternalServerError, msg, nil)
		return
	}
	var file *xlsx.File
	var sheet *xlsx.Sheet

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("host")

	objID := common.BKInnerObjIDHost
	fields, err := s.Logics.GetObjFieldIDs(objID, logics.GetFilterFields(objID), c.Request.Header)
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
		reply := getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}
	logics.AddDownExcelHttpHeader(c, "host.xlsx")
	c.File(dirFileName)

	os.Remove(dirFileName)

}

//BuildDownLoadExcelTemplate build download excel template
func (s *Service) BuildDownLoadExcelTemplate(c *gin.Context) {
	logics.SetProxyHeader(c)
	objID := c.Param(common.BKObjIDField)
	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/template/"
	_, err := os.Stat(dir)
	if nil != err {
		os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	}
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	file := fmt.Sprintf("%s/%stemplate-%d-%d.xlsx", dir, objID, time.Now().UnixNano(), randNum)
	err = s.Logics.BuildExcelTemplate(objID, file, c.Request.Header, defLang)
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
