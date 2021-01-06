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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"
)

// ImportHost import host
func (s *Service) ImportHost(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromHTTPHeader(c.Request.Header)

	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)
	file, err := c.FormFile("file")
	if nil != err {
		blog.Errorf("ImportHost failed, get file from form data failed, err: %+v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}

	moduleID := int64(0)
	if moduleIDStr := c.PostForm(common.BKModuleIDField); moduleIDStr != "" {
		moduleID, err = strconv.ParseInt(moduleIDStr, 10, 64)
		if err != nil {
			blog.Errorf("ImportHost failed, bk_module_id not integer, err: %+v, bk_module_id: %s,  rid: %s", err, moduleIDStr, rid)
			msg := getReturnStr(common.CCErrCommParamsNeedInt, defErr.CCErrorf(common.CCErrCommParamsNeedInt, common.BKModuleIDField).Error(), nil)
			c.String(http.StatusOK, msg)
			return
		}
	}

	webCommon.SetProxyHeader(c)

	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/import/"
	_, err = os.Stat(dir)
	if nil != err {
		if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("ImportHost failed, save form data to local file failed, mkdir failed, err: %+v, rid: %s", err, rid)
			c.String(http.StatusInternalServerError, fmt.Sprintf("save form data to local file failed, mkdir failed, err: %+v", err))
			return
		}
	}
	filePath := fmt.Sprintf("%s/importhost-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	if err := c.SaveUploadedFile(file, filePath); nil != err {
		blog.Errorf("ImportHost failed, save form data to local file failed, save data as excel failed, err: %+v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}

	// del file
	defer func(filePath string, rid string) {
		if err := os.Remove(filePath); err != nil {
			blog.Errorf("ImportHost, remove temporary file failed, err: %+v, rid: %s", err, rid)
		}
	}(filePath, rid)

	f, err := xlsx.OpenFile(filePath)
	if nil != err {
		blog.Errorf("ImportHost failed, open form data as excel file failed, err: %+v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}
	result := s.Logics.ImportHosts(ctx, f, c.Request.Header, defLang, 0, moduleID)

	c.JSON(http.StatusOK, result)
}

// ExportHost export host
func (s *Service) ExportHost(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)
	webCommon.SetProxyHeader(c)
	header := c.Request.Header
	defLang := s.Language.CreateDefaultCCLanguageIf(util.GetLanguage(header))
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	customFieldsStr := c.PostForm(common.ExportCustomFields)

	hostIDStr := c.PostForm("bk_host_id")
	appIDStr := c.PostForm("bk_biz_id")
	exportCondStr := c.PostForm("export_condition")
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ExportHost failed, bk_biz_id not integer. err: %+v, biz id: %s,  rid: %s", err, appIDStr, rid)
		err := defErr.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
		reply := getReturnStr(err.GetCode(), err.Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}

	objID := common.BKInnerObjIDHost
	filterFields := logics.GetFilterFields(objID)
	customFields := logics.GetCustomFields(filterFields, customFieldsStr)
	fields, err := s.Logics.GetObjFieldIDs(objID, filterFields, customFields, c.Request.Header, appID)
	if nil != err {
		blog.Errorf("ExportHost failed, get host model fields failed, err: %+v, rid: %s", err, rid)
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}

	var hostFields []string
	for _, property := range fields {
		hostFields = append(hostFields, property.ID)
	}

	hostInfo, err := s.Logics.GetHostData(appID, hostIDStr, hostFields, exportCondStr, header, defLang)
	if err != nil {
		blog.Errorf("ExportHost failed, get hosts failed, err: %+v, bk_host_id:%s, export_condition:%s, rid: %s", err, hostIDStr, exportCondStr, rid)
		reply := getReturnStr(common.CCErrWebGetHostFail, defErr.Errorf(common.CCErrWebGetHostFail, err.Error()).Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}
	if len(hostInfo) == 0 {
		blog.Errorf("ExportHost failed, get hosts failed, no host is found, bk_host_id:%s, export_condition:%s, rid: %s", hostIDStr, exportCondStr, rid)
		reply := getReturnStr(common.CCErrWebGetHostFail, defErr.Errorf(common.CCErrWebGetHostFail, "no host is found").Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}

	var file *xlsx.File
	file = xlsx.NewFile()

	usernameMap, propertyList, err := s.getUsernameMapWithPropertyList(c, objID, hostInfo)
	if nil != err {
		blog.Errorf("ExportHost failed, get username map and property list failed, err: %+v, rid: %s", err, rid)
		reply := getReturnStr(common.CCErrWebGetUsernameMapFail, defErr.Errorf(common.CCErrWebGetUsernameMapFail, objID).Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}

	err = s.Logics.BuildHostExcelFromData(context.Background(), objID, fields, nil, hostInfo, file, header, 0, usernameMap, propertyList)
	if nil != err {
		blog.Errorf("ExportHost failed, BuildHostExcelFromData failed, object:%s, err:%+v, rid:%s", objID, err, rid)
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if nil != err {
		if err := os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("ExportHost failed, make local dir to save export file failed, err: %+v, rid: %s", err, rid)
			c.String(http.StatusInternalServerError, fmt.Sprintf("make local dir to save export file failed, err: %+v", err))
			return
		}
	}
	fileName := fmt.Sprintf("%dhost.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)

	logics.ProductExcelCommentSheet(ctx, file, defLang)
	err = file.Save(dirFileName)
	if err != nil {
		blog.Errorf("ExportHost failed, save file failed, err: %+v, rid: %s", err, rid)
		reply := getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}
	logics.AddDownExcelHttpHeader(c, "bk_cmdb_export_host.xlsx")
	c.File(dirFileName)

	if err := os.Remove(dirFileName); err != nil {
		blog.Errorf("ExportHost success, but remove host.xlsx file failed, err: %+v, rid: %s", err, rid)
	}
}

// BuildDownLoadExcelTemplate build download excel template
func (s *Service) BuildDownLoadExcelTemplate(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)

	webCommon.SetProxyHeader(c)
	objID := c.Param(common.BKObjIDField)
	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/template/"
	_, err := os.Stat(dir)
	if nil != err {
		if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("BuildDownLoadExcelTemplate failed, make template dir failed, err: %+v, rid: %s", err, rid)
			c.String(http.StatusInternalServerError, fmt.Sprintf("make template dir failed, err: %+v", err))
			return
		}
	}
	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	modelBizID, err := parseModelBizID(c.PostForm(common.BKAppIDField))
	if err != nil {
		msg := getReturnStr(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}

	file := fmt.Sprintf("%s/%stemplate-%d-%d.xlsx", dir, objID, time.Now().UnixNano(), randNum)
	err = s.Logics.BuildExcelTemplate(ctx, objID, file, c.Request.Header, defLang, modelBizID)
	if nil != err {
		blog.Errorf("BuildDownLoadExcelTemplate failed, build excel template failed, object:%s error:%s, rid: %s", objID, err.Error(), rid)
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}
	if objID == common.BKInnerObjIDHost {
		logics.AddDownExcelHttpHeader(c, "bk_cmdb_import_host.xlsx")
	} else {
		logics.AddDownExcelHttpHeader(c, fmt.Sprintf("bk_cmdb_inst_%s.xlsx", objID))
	}

	// http.ServeFile(c.Writer, c.Request, file)
	c.File(file)
	if err := os.Remove(file); err != nil {
		blog.Errorf("BuildDownLoadExcelTemplate success, but remove template file after response failed, err: %+v, rid: %s", err, rid)
	}
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

func (s *Service) ListenIPOptions(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)
	webCommon.SetProxyHeader(c)
	header := c.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))

	hostIDStr := c.Param("bk_host_id")
	hostID, err := strconv.ParseInt(hostIDStr, 10, 64)
	if err != nil {
		blog.Infof("host id invalid, convert to int failed, hostID: %s, err: %+v, rid: %s", hostID, err, rid)
		result := metadata.ResponseDataMapStr{
			BaseResp: metadata.BaseResp{
				Result: false,
				Code:   common.CCErrCommParamsInvalid,
				ErrMsg: defErr.Errorf(common.CCErrCommParamsInvalid, common.BKHostIDField).Error(),
			},
		}
		c.JSON(http.StatusOK, result)
		return
	}
	option := metadata.ListHostsWithNoBizParameter{
		HostPropertyFilter: &querybuilder.QueryFilter{
			Rule: querybuilder.CombinedRule{
				Condition: querybuilder.ConditionAnd,
				Rules: []querybuilder.Rule{
					querybuilder.AtomRule{
						Field:    common.BKHostIDField,
						Operator: querybuilder.OperatorEqual,
						Value:    hostID,
					},
				},
			},
		},
		Fields: []string{
			common.BKHostIDField,
			common.BKHostNameField,
			common.BKHostInnerIPField,
			common.BKHostOuterIPField,
		},
		Page: metadata.BasePage{
			Start: 0,
			Limit: 1,
		},
	}
	resp, err := s.CoreAPI.ApiServer().ListHostWithoutApp(ctx, c.Request.Header, option)
	if err != nil {
		blog.Errorf("get host by id failed, hostID: %d, err: %+v, rid: %s", hostID, err, rid)
		result := metadata.ResponseDataMapStr{
			BaseResp: metadata.BaseResp{
				Result: false,
				Code:   common.CCErrHostGetFail,
				ErrMsg: defErr.Error(common.CCErrHostGetFail).Error(),
			},
		}
		c.JSON(http.StatusOK, result)
		return
	}
	if resp.Code != 0 || resp.Result == false {
		blog.Errorf("got host by id failed, hostID: %d, response: %+v, rid: %s", hostID, resp, rid)
		c.JSON(http.StatusOK, resp)
		return
	}
	if len(resp.Data.Info) == 0 {
		blog.Errorf("host not found, hostID: %d, rid: %s", hostID, rid)
		result := metadata.ResponseDataMapStr{
			BaseResp: metadata.BaseResp{
				Result: false,
				Code:   common.CCErrCommNotFound,
				ErrMsg: defErr.Error(common.CCErrCommNotFound).Error(),
			},
		}
		c.JSON(http.StatusOK, result)
		return
	}
	type Host struct {
		HostID   int64  `json:"bk_host_id" bson:"bk_host_id"`           // 主机ID(host_id)								数字
		HostName string `json:"bk_host_name" bson:"bk_host_name"`       // 主机名称
		InnerIP  string `json:"bk_host_innerip" bson:"bk_host_innerip"` // 内网IP
		OuterIP  string `json:"bk_host_outerip" bson:"bk_host_outerip"` // 外网IP
	}
	host := Host{}
	raw := resp.Data.Info[0]
	if err := mapstr.DecodeFromMapStr(&host, raw); err != nil {
		msg := fmt.Sprintf("decode response data into host failed, raw: %+v, err: %+v, rid: %s", raw, err, rid)
		blog.Error(msg)
		result := metadata.ResponseDataMapStr{
			BaseResp: metadata.BaseResp{
				Result: false,
				Code:   common.CCErrCommJSONUnmarshalFailed,
				ErrMsg: defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(),
			},
		}
		c.JSON(http.StatusOK, result)
		return
	}

	ipOptions := make([]string, 0)
	ipOptions = append(ipOptions, "127.0.0.1")
	ipOptions = append(ipOptions, "0.0.0.0")
	if len(host.InnerIP) > 0 {
		ipOptions = append(ipOptions, host.InnerIP)
	}
	if len(host.OuterIP) > 0 {
		ipOptions = append(ipOptions, host.OuterIP)
	}
	result := metadata.ResponseDataMapStr{
		BaseResp: metadata.BaseResp{
			Result: true,
			Code:   0,
		},
		Data: map[string]interface{}{
			"options": ipOptions,
		},
	}
	c.JSON(http.StatusOK, result)
	return
}

// UpdateHost Excel update host batch
func (s *Service) UpdateHosts(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromHTTPHeader(c.Request.Header)

	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)
	bizIDStr := c.PostForm("bk_biz_id")
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHosts failed, bk_biz_id not integer. err: %+v, biz id: %s,  rid: %s", err, bizIDStr, rid)
		err := defErr.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
		reply := getReturnStr(err.GetCode(), err.Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}
	file, err := c.FormFile("file")
	if nil != err {
		blog.Errorf("UpdateHost excel import update hosts failed, get file from form data failed, err: %+v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	webCommon.SetProxyHeader(c)

	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/import/"
	_, err = os.Stat(dir)
	if nil != err {
		if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("UpdateHost excel import update hosts, save form data to local file failed, mkdir failed, err: %+v, rid: %s", err, rid)
			c.String(http.StatusInternalServerError, fmt.Sprintf("save form data to local file failed, mkdir failed, err: %+v", err))
			return
		}
	}
	filePath := fmt.Sprintf("%s/importhost-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	if err := c.SaveUploadedFile(file, filePath); nil != err {
		blog.Errorf("UpdateHosts failed, save form data to local file failed, save data as excel failed, err: %+v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	// del file
	defer func(filePath string, rid string) {
		if err := os.Remove(filePath); err != nil {
			blog.Errorf("UpdateHost excel import update hosts, remove temporary file failed, err: %+v, rid: %s", err, rid)
		}
	}(filePath, rid)

	f, err := xlsx.OpenFile(filePath)
	if nil != err {
		blog.Errorf("UpdateHost excel import update hosts failed, open form data as excel file failed, err: %+v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	result := s.Logics.UpdateHosts(ctx, f, c.Request.Header, defLang, bizID)

	c.JSON(http.StatusOK, result)
}
