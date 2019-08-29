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
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"
)

func (s *Service) ImportNetProperty(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)
	header := c.Request.Header
	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)
	webCommon.SetProxyHeader(c)

	// open uploaded file
	file, err, errMsg := openNetPropertyUploadedFile(c, defErr)
	if nil != err {
		blog.Errorf("[Import Net Property] open uploaded file error:%s, rid: %s", err.Error(), rid)
		c.String(http.StatusInternalServerError, string(errMsg))
		return
	}

	netProperty, err, errMsg := getNetPropertiesFromFile(c.Request.Header, defLang, defErr, file)
	if nil != err {
		blog.Errorf("[Import Net Property] failed, error:%s, rid: %s", err.Error(), rid)
		c.String(http.StatusInternalServerError, string(errMsg))
		return
	}

	data := make([]interface{}, 0)
	lineNumbers := make([]int, 0)
	for line, value := range netProperty {
		data = append(data, value)
		lineNumbers = append(lineNumbers, line)
	}
	params := mapstr.MapStr{"Data": data}

	propertyResult, err := s.Engine.CoreAPI.ApiServer().SearchNetDevicePropertyBatch(context.Background(), header, params)
	if nil != err {
		blog.Errorf("search net device property data  batch  error:%#v , search condition:%#v, rid: %s", err, params, rid)
		msg := getReturnStr(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	// rebuild response body
	resultStr, err := rebuildNetPropertyResponseBody(ctx, propertyResult.Data.Info, lineNumbers)
	if nil != err {
		c.String(http.StatusInternalServerError, getReturnStr(common.CCErrWebGetAddNetPropertyResultFail,
			defErr.Errorf(common.CCErrWebGetAddNetPropertyResultFail).Error(), nil))
		return
	}

	c.String(http.StatusOK, resultStr)
}

func (s *Service) ExportNetProperty(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)
	webCommon.SetProxyHeader(c)

	netPropertyIDStr := c.PostForm(common.BKNetcollectPropertyIDField)
	netPropertyInfo, err := s.Logics.GetNetPropertyData(c.Request.Header, netPropertyIDStr)
	if nil != err {
		blog.Errorf("[Export Net Property] get property data error:%s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrWebGetHostFail, defErr.Errorf(common.CCErrWebGetHostFail, err.Error()).Error(), nil)
		c.String(http.StatusBadGateway, msg)
		return
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet(common.BKNetProperty)
	if nil != err {
		blog.Errorf("[Export Net Property] create sheet error:%s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrWebCreateEXCELFail,
			defErr.Errorf(common.CCErrWebCreateEXCELFail, err.Error()).Error(), nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	fields := logics.GetNetPropertyField(defLang)
	logics.AddNetPropertyExtFields(fields, defLang)

	if err = logics.BuildNetPropertyExcelFromData(ctx, defLang, fields, netPropertyInfo, sheet); nil != err {
		blog.Errorf("[Export Net Property] build net property excel data error:%s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrWebCreateEXCELFail,
			defErr.Errorf(common.CCErrWebCreateEXCELFail, err.Error()).Error(), nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	if _, err = os.Stat(dirFileName); nil != err {
		if err = os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm); nil != err {
			blog.Errorf("[Export Net Property] mkdir error:%s, rid: %s", err.Error(), rid)
			msg := getReturnStr(common.CCErrWebCreateEXCELFail,
				defErr.Errorf(common.CCErrWebCreateEXCELFail, err.Error()).Error(), nil)
			c.String(http.StatusInternalServerError, msg)
			return
		}
	}

	fileName := fmt.Sprintf("%dnetproperty.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)

	logics.ProductExcelCommentSheet(ctx, file, defLang)

	if err = file.Save(dirFileName); nil != err {
		blog.Errorf("[Export Net Property] save file error:%s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrWebCreateEXCELFail,
			defErr.Errorf(common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	logics.AddDownExcelHttpHeader(c, "netproperty.xlsx")
	c.File(dirFileName)

	if err = os.Remove(dirFileName); nil != err {
		blog.Errorf("[Export Net Property] remove file error:%s, rid: %s", err.Error(), rid)
	}
}

func (s *Service) BuildDownLoadNetPropertyExcelTemplate(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	webCommon.SetProxyHeader(c)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	dir := webCommon.ResourcePath + "/template/"
	if _, err := os.Stat(dir); nil != err {
		if err = os.MkdirAll(dir, os.ModeDir|os.ModePerm); nil != err {
			blog.Errorf("[Build Net Property Excel Template] mkdir error:%s, rid: %s", err.Error(), rid)
			msg := getReturnStr(common.CCErrCommExcelTemplateFailed,
				defErr.Errorf(common.CCErrCommExcelTemplateFailed, common.BKNetProperty).Error(), nil)
			c.String(http.StatusInternalServerError, msg)
			return
		}
	}

	file := fmt.Sprintf("%s/%stemplate-%d-%d.xlsx", dir, common.BKNetProperty, time.Now().UnixNano(), rand.Uint32())

	if err := logics.BuildNetPropertyExcelTemplate(c.Request.Header, defLang, file); nil != err {
		blog.Errorf("Build Net Property Excel Template, error:%s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrCommExcelTemplateFailed,
			defErr.Errorf(common.CCErrCommExcelTemplateFailed, common.BKNetProperty).Error(),
			nil)
		c.String(http.StatusInternalServerError, msg)
		return
	}

	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("template_%s.xlsx", common.BKNetProperty))

	c.File(file)

	if err := os.Remove(file); nil != err {
		blog.Errorf("[Build Net Property Excel Template] mkdir error:%s, rid: %s", err.Error(), rid)
	}
}

func openNetPropertyUploadedFile(c *gin.Context, defErr errors.DefaultCCErrorIf) (file *xlsx.File, err error, errMsg string) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	fileHeader, err := c.FormFile("file")
	if nil != err {
		errMsg = getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		return nil, err, errMsg
	}

	dir := webCommon.ResourcePath + "/import/"
	if _, err = os.Stat(dir); nil != err {
		if err = os.MkdirAll(dir, os.ModeDir|os.ModePerm); nil != err {
			blog.Errorf("[Import Net Property] mkdir error:%s, rid: %s", err.Error(), rid)
			errMsg = getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
			return nil, err, errMsg
		}
	}

	filePath := fmt.Sprintf("%s/importnetproperty-%d-%d.xlsx", dir, time.Now().UnixNano(), rand.Uint32())
	if err = c.SaveUploadedFile(fileHeader, filePath); nil != err {
		errMsg = getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
		return nil, err, errMsg
	}

	defer func() {
		// del file
		if err := os.Remove(filePath); err != nil {
			blog.Errorf("remove file %s failed, err: %+v, rid: %s", filePath, err, rid)
		}
	}()

	file, err = xlsx.OpenFile(filePath)
	if nil != err {
		errMsg = getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		return nil, err, errMsg
	}

	return file, nil, ""
}

func getNetPropertiesFromFile(header http.Header, defLang lang.DefaultCCLanguageIf, defErr errors.DefaultCCErrorIf, file *xlsx.File) (
	netProperty map[int]map[string]interface{}, err error, errMsg string) {
	rid := util.GetHTTPCCRequestID(header)

	netProperty, errMsgs, err := logics.GetImportNetProperty(header, defLang, file)
	if nil != err {
		blog.Errorf("[Import Net Property] failed, error:%s, rid: %s", err.Error(), rid)
		errMsg = getReturnStr(common.CCErrWebFileContentFail,
			defErr.Errorf(common.CCErrWebFileContentFail, err.Error()).Error(), nil)
		return nil, err, errMsg
	}
	if 0 != len(errMsgs) {
		errMsg = getReturnStr(common.CCErrWebFileContentFail,
			defErr.Errorf(common.CCErrWebFileContentFail, " file empty").Error(),
			common.KvMap{"err": errMsgs})
		return nil, err, errMsg
	}
	if 0 == len(netProperty) {
		errMsg = getReturnStr(common.CCErrWebFileContentEmpty,
			defErr.Errorf(common.CCErrWebFileContentEmpty, "").Error(), nil)
		return nil, err, errMsg
	}

	return netProperty, nil, ""
}

func rebuildNetPropertyResponseBody(ctx context.Context, addPropertyResult []mapstr.MapStr, lineNumbers []int) (string, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	replyBody := new(meta.Response)

	var (
		errRow  []string
		succRow []string
	)
	for i, data := range addPropertyResult {

		result, ok := data["result"].(bool)
		if !ok {
			blog.Errorf("[Import Net Property] data is not bool: %#+v, rid: %s", data["result"], rid)
			return "", fmt.Errorf("convert response body fail")
		}

		switch result {
		case true:
			succRow = append(succRow, strconv.Itoa(lineNumbers[i]))
		case false:
			errMsg, ok := data["error_msg"].(string)
			if !ok {
				blog.Errorf("[Import Net Property] data is not string: %#+v, rid: %s", data["error_msg"], rid)
				return "", fmt.Errorf("convert response body fail")
			}

			errRow = append(errRow, fmt.Sprintf("%d行%s", lineNumbers[i], errMsg))
		}
	}

	retData := make(map[string]interface{})
	if 0 < len(succRow) {
		retData["success"] = succRow
	}
	if 0 < len(errRow) {
		retData["error"] = errRow
	}

	replyBody.Data = retData

	replyByte, err := json.Marshal(replyBody)
	if nil != err {
		blog.Errorf("[Import Net Property] convert rebuilt response body fail, error: %v, rid: %s", err, rid)
		return "", err
	}

	return string(replyByte), nil
}
