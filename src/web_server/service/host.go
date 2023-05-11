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
	"io"
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
	"configcenter/src/thirdparty/hooks/process"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx/v3"
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

	params := c.PostForm("params")
	if params == "" {
		blog.ErrorJSON("ImportHost failed, not found params value, rid: %s", rid)
		msg := getReturnStr(common.CCErrCommParamsNeedSet,
			defErr.CCErrorf(common.CCErrCommParamsNeedSet, "params").Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}
	inputJSON := &metadata.ExcelImportAddHostInput{}
	if err := json.Unmarshal([]byte(params), inputJSON); err != nil {
		blog.ErrorJSON("ImportHost failed, params unmarshal error, err: %s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrCommParamsValueInvalidError,
			defErr.CCErrorf(common.CCErrCommParamsValueInvalidError, "params", err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
		return
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

	f, err := xlsx.OpenFile(filePath, xlsx.UseDiskVCellStore)
	if nil != err {
		blog.Errorf("ImportHost failed, open form data as excel file failed, err: %+v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}
	result := s.Logics.ImportHosts(ctx, f, c.Request.Header, defLang, 0, inputJSON.ModuleID,
		inputJSON.OpType, inputJSON.AssociationCond, inputJSON.ObjectUniqueID)

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

	input := &metadata.ExcelExportHostInput{}
	if err := c.BindJSON(input); err != nil {
		blog.ErrorJSON("Unmarshal input error. input: %s, err: %s, rid: %s", c.Keys, err.Error(), rid)

		ccErr := defErr.CCError(common.CCErrCommJSONUnmarshalFailed)
		result := metadata.ResponseDataMapStr{
			BaseResp: metadata.BaseResp{
				Result: false,
				Code:   ccErr.GetCode(),
				ErrMsg: ccErr.Error(),
			},
		}
		c.JSON(http.StatusOK, result)
		return
	}

	if input.ExportCond.Page.Limit <= 0 || input.ExportCond.Page.Limit > common.BKMaxOnceExportLimit {
		blog.Errorf("host page input is illegal, page: %v, rid: %s", input.ExportCond.Page, rid)
		_, _ = c.Writer.Write([]byte(getReturnStr(common.CCErrWebGetHostFail, defErr.Errorf(
			common.CCErrWebGetHostFail, defLang.Languagef("export_page_limit_err", common.BKMaxOnceExportLimit)).
			Error(), nil)))
		return
	}

	objectName, objIDs, err := s.getCustomObjectInfo(ctx, header)
	if err != nil {
		blog.Errorf("get custom instance name failed, err: %v, rid: %s", err, rid)
		return
	}

	appID := input.AppID
	objID := common.BKInnerObjIDHost
	filterFields := logics.GetFilterFields(objID)
	customFields := logics.GetCustomFields(filterFields, input.CustomFields)
	// customLen+5为生成主机数据的起始列索引, 5=字段说明1列+业务拓扑，业务名，集群，模块4列
	fields, err := s.Logics.GetObjFieldIDs(objID, filterFields, customFields, c.Request.Header, appID, len(objectName)+5)
	if err != nil {
		blog.Errorf("get host model fields failed, err: %v, rid: %s", err, rid)
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed,
			objID).Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}

	file := xlsx.NewFile()
	err = s.Logics.BuildHostExcelFromData(c, objID, fields, nil, file, header, objectName, objIDs, input, s.Config)
	if nil != err {
		blog.Errorf("ExportHost failed, BuildHostExcelFromData failed, object:%s, err:%+v, rid:%s", objID, err,
			rid)
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed,
			objID).Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	if _, err = os.Stat(dirFileName); err != nil && os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm) != nil {
		blog.Errorf("ExportHost failed, make local dir to save export file failed, err: %+v, rid: %s", err, rid)
		c.String(http.StatusInternalServerError, fmt.Sprintf("make local dir to save export file failed, err: %+v", err))
		return
	}

	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fmt.Sprintf("%dhost.xlsx", time.Now().UnixNano()))
	if err := logics.ProductExcelCommentSheet(ctx, file, defLang); err != nil {
		blog.Errorf("export host failed, err: %+v, rid: %s", err, rid)
		_, _ = c.Writer.Write([]byte(getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(
			common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)))
		return
	}

	if err := file.Save(dirFileName); err != nil {
		blog.Errorf("ExportHost failed, save file failed, err: %+v, rid: %s", err, rid)
		_, _ = c.Writer.Write([]byte(getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(
			common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)))
		return
	}
	logics.AddDownExcelHttpHeader(c, "bk_cmdb_export_host.xlsx")
	if err := s.writeFile(c, dirFileName, rid); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("write exported file failed, err: %+v", err))
		return
	}

	defer func(dirFileName string, rid string) {
		if err := os.Remove(dirFileName); err != nil {
			blog.Errorf("export host success, but remove host.xlsx file failed, err: %+v, rid: %s", err, rid)
		}
	}(dirFileName, rid)
}

func (s *Service) writeFile(c *gin.Context, dirFileName string, rid string) error {
	exportedFile, err := os.Open(dirFileName)
	if err != nil {
		blog.Errorf("open exported file failed, file: %s, err: %+v, rid: %s", dirFileName, err, rid)
		return err
	}
	defer exportedFile.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := exportedFile.Read(buf)
		if err != nil && err != io.EOF {
			blog.Errorf("read exported file failed, file: %s, err: %+v, rid: %s", dirFileName, err, rid)
			return err
		}
		if n == 0 {
			break
		}
		if _, err := c.Writer.Write(buf[:n]); err != nil {
			blog.Errorf("write exported file failed, file: %s, err: %+v, rid: %s", dirFileName, err, rid)
			return err
		}
		c.Writer.Flush()
	}

	return nil
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

// Host simplified host struct
type Host struct {
	// HostID 主机ID(host_id)
	HostID int64 `json:"bk_host_id" bson:"bk_host_id"`
	// HostName 主机名称
	HostName string `json:"bk_host_name" bson:"bk_host_name"`
	// InnerIP 内网IP
	InnerIP string `json:"bk_host_innerip" bson:"bk_host_innerip"`
	// OuterIP 外网IP
	OuterIP string `json:"bk_host_outerip" bson:"bk_host_outerip"`
}

// ListenIPOptions TODO
func (s *Service) ListenIPOptions(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)
	webCommon.SetProxyHeader(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(c.Request.Header))

	hostIDStr := c.Param("bk_host_id")
	hostID, err := strconv.ParseInt(hostIDStr, 10, 64)
	if err != nil {
		blog.Infof("host id invalid, convert to int failed, hostID: %s, err: %+v, rid: %s", hostID, err, rid)
		result := metadata.BaseResp{Result: false, Code: common.CCErrCommParamsInvalid,
			ErrMsg: defErr.Errorf(common.CCErrCommParamsInvalid, common.BKHostIDField).Error()}
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
			common.BKHostInnerIPv6Field,
			common.BKHostOuterIPv6Field,
		},
		Page: metadata.BasePage{
			Start: 0,
			Limit: 1,
		},
	}
	resp, err := s.CoreAPI.ApiServer().ListHostWithoutApp(ctx, c.Request.Header, option)
	if err != nil {
		blog.Errorf("get host by id failed, hostID: %d, err: %+v, rid: %s", hostID, err, rid)
		result := metadata.BaseResp{Result: false, Code: common.CCErrHostGetFail,
			ErrMsg: defErr.Error(common.CCErrHostGetFail).Error()}
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
		result := metadata.BaseResp{Result: false, Code: common.CCErrCommNotFound,
			ErrMsg: defErr.Error(common.CCErrCommNotFound).Error()}
		c.JSON(http.StatusOK, result)
		return
	}
	type hostBase struct {
		HostID    int64  `json:"bk_host_id"`
		HostName  string `json:"bk_host_name"`
		InnerIP   string `json:"bk_host_innerip"`
		InnerIPv6 string `json:"bk_host_innerip_v6"`
		OuterIP   string `json:"bk_host_outerip"`
		OuterIPv6 string `json:"bk_host_outerip_v6"`
	}
	host := hostBase{}
	raw := resp.Data.Info[0]
	if err := mapstr.DecodeFromMapStr(&host, raw); err != nil {
		msg := fmt.Sprintf("decode response data into host failed, raw: %+v, err: %+v, rid: %s", raw, err, rid)
		blog.Error(msg)
		result := metadata.BaseResp{Result: false, Code: common.CCErrCommJSONUnmarshalFailed,
			ErrMsg: defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error()}
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

	// add process ipv6 options if needed
	if process.NeedIPv6OptionsHook() {
		ipOptions = append(ipOptions, "::1")
		ipOptions = append(ipOptions, "::")
		if len(host.InnerIPv6) > 0 {
			ipOptions = append(ipOptions, host.InnerIPv6)
		}
		if len(host.OuterIPv6) > 0 {
			ipOptions = append(ipOptions, host.OuterIPv6)
		}
	}

	result := metadata.ResponseDataMapStr{
		BaseResp: metadata.BaseResp{Result: true, Code: 0},
		Data: map[string]interface{}{
			"options": ipOptions,
		},
	}
	c.JSON(http.StatusOK, result)
	return
}

// UpdateHosts Excel update host batch
func (s *Service) UpdateHosts(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromHTTPHeader(c.Request.Header)

	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	params := c.PostForm("params")
	if params == "" {
		blog.ErrorJSON("ImportHost failed, not found params value, rid: %s", rid)
		msg := getReturnStr(common.CCErrCommParamsNeedSet,
			defErr.CCErrorf(common.CCErrCommParamsNeedSet, "params").Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}
	inputJSON := &metadata.ExcelImportUpdateHostInput{}
	if err := json.Unmarshal([]byte(params), inputJSON); err != nil {
		blog.ErrorJSON("ImportHost failed, params unmarshal error, err: %s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrCommParamsValueInvalidError,
			defErr.CCErrorf(common.CCErrCommParamsValueInvalidError, "params", err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}

	file, err := c.FormFile("file")
	if nil != err {
		blog.Errorf("excel import update hosts failed, get file from form data failed, err: %+v, rid: %s", err, rid)
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
			blog.Errorf("save form data to local file failed, mkdir failed, err: %+v, rid: %s", err, rid)
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("save form data to local file failed, mkdir failed, err: %+v", err))
			return
		}
	}
	filePath := fmt.Sprintf("%s/importhost-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	if err := c.SaveUploadedFile(file, filePath); nil != err {
		blog.Errorf("save form data to local file failed, save data as excel failed, err: %+v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail,
			err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	// del file
	defer func(filePath string, rid string) {
		if err := os.Remove(filePath); err != nil {
			blog.Errorf("UpdateHost excel import update hosts, remove temporary file failed, err: %+v, rid: %s", err, rid)
		}
	}(filePath, rid)

	f, err := xlsx.OpenFile(filePath, xlsx.UseDiskVCellStore)
	if nil != err {
		blog.Errorf("excel import update hosts failed, open form data as excel file failed, err: %+v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail,
			err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	result := s.Logics.UpdateHosts(ctx, f, c.Request.Header, defLang, inputJSON.BizID, inputJSON.OpType,
		inputJSON.AssociationCond, inputJSON.ObjectUniqueID)

	c.JSON(http.StatusOK, result)
}

// getCustomObjectInfo get custom instance object info
func (s *Service) getCustomObjectInfo(ctx context.Context, header http.Header) ([]string, []string, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.AssociationKindIDField: common.AssociationKindMainline,
		},
	}
	mainlineAsstRsp, err := s.Engine.CoreAPI.ApiServer().ReadModuleAssociation(context.Background(), header, query)
	if err != nil {
		blog.Errorf("search mainline association failed, err: %v, rid: %s", err, rid)
		return nil, nil, err
	}

	mainlineObjectChildMap := make(map[string]string, 0)
	for _, asst := range mainlineAsstRsp.Data.Info {
		if asst.ObjectID == common.BKInnerObjIDHost {
			continue
		}
		mainlineObjectChildMap[asst.AsstObjID] = asst.ObjectID
	}

	// get all mainline custom object id
	objectIDs := make([]string, 0)
	for objectID := common.BKInnerObjIDApp; len(objectID) != 0; objectID = mainlineObjectChildMap[objectID] {
		if objectID == common.BKInnerObjIDApp || objectID == common.BKInnerObjIDSet ||
			objectID == common.BKInnerObjIDModule {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}
	if len(objectIDs) == 0 {
		return nil, nil, nil
	}

	input := &metadata.QueryCondition{
		Fields: []string{common.BKObjNameField, common.BKObjIDField},
		Condition: mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objectIDs},
		},
	}

	objectName := make([]string, 0)
	objects, err := s.Logics.CoreAPI.ApiServer().ReadModel(ctx, header, input)
	if err != nil {
		blog.Errorf("search mainline obj failed, objIDs: %#v, err: %v, rid: %s", objectIDs, err, rid)
		return objectName, util.ReverseArrayString(objectIDs), err
	}

	objIDNameMap := make(map[string]string, 0)
	for _, val := range objects.Data.Info {
		objIDNameMap[val.ObjectID] = val.ObjectName
	}

	for _, objID := range objectIDs {
		if objName, ok := objIDNameMap[objID]; ok {
			objectName = append(objectName, objName)
		}
	}

	return objectName, util.ReverseArrayString(objectIDs), nil
}
