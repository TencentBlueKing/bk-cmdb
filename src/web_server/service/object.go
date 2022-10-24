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
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"

	"github.com/alexmullins/zip"
	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"
)

var sortFields = []string{
	"bk_property_id",
	"bk_property_name",
	"bk_property_type",
	"bk_property_group_name",
	"option",
	"unit",
	"description",
	"placeholder",
	"editable",
	"isrequired",
	"isreadonly",
	"isonly",
}

// ImportObject import object attribute
func (s *Service) ImportObject(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	webCommon.SetProxyHeader(c)
	objID := c.Param(common.BKObjIDField)
	ctx := util.NewContextFromGinContext(c)

	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	file, err := c.FormFile("file")
	if err != nil {
		msg := getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	modelBizID, err := parseModelBizID(c.PostForm(common.BKAppIDField))
	if err != nil {
		msg := getReturnStr(common.CCErrCommJSONUnmarshalFailed,
			defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/import/"
	if _, err = os.Stat(dir); err != nil {
		blog.Warnf("os.Stat failed, filename: %s, err: %+v, rid: %s", dir, err, rid)
		if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("os.MkdirAll failed, filename: %s, err: %+v, rid: %s", dir, err, rid)
		}
	}
	filePath := fmt.Sprintf("%s/importinsts-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	if err = c.SaveUploadedFile(file, filePath); err != nil {
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail,
			err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	defer func() {
		if err := os.Remove(filePath); err != nil {
			blog.Errorf("os.Remove failed, filename: %s, err: %+v, rid: %s", filePath, err, rid)
		}
	}()
	f, err := xlsx.OpenFile(filePath)
	if err != nil {
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail,
			err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	attrItems, errMsg, err := s.Logics.GetImportInsts(ctx, f, objID, c.Request.Header, 3, false, defLang, modelBizID)
	if len(attrItems) == 0 {
		var msg string
		if err != nil {
			msg = getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail,
				err.Error()).Error(), nil)
		} else {
			msg = getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail,
				"").Error(), nil)
		}
		c.String(http.StatusOK, string(msg))
		return
	}
	if len(errMsg) != 0 {
		msg := getReturnStr(common.CCErrWebFileContentFail, defErr.Errorf(common.CCErrWebFileContentFail,
			strings.Join(errMsg, ",")).Error(), common.KvMap{"err": errMsg})
		c.String(http.StatusOK, string(msg))
		return
	}

	logics.ConvAttrOption(attrItems)

	params := map[string]interface{}{objID: map[string]interface{}{"attr": attrItems}}

	result, err := s.CoreAPI.ApiServer().AddObjectBatch(ctx, c.Request.Header, params)
	if err != nil {
		msg := getReturnStr(common.CCErrCommHTTPDoRequestFailed, defErr.Errorf(common.CCErrCommHTTPDoRequestFailed,
			"").Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	c.JSON(http.StatusOK, result)
}

func setExcelSubTitle(row *xlsx.Row) *xlsx.Row {
	for _, key := range sortFields {
		cell := row.AddCell()
		cell.Value = key
	}
	return row
}

func setExcelTitle(ctx context.Context, row *xlsx.Row, defLang lang.DefaultCCLanguageIf) *xlsx.Row {
	rid := util.ExtractRequestIDFromContext(ctx)

	fields := logics.GetPropertyFieldDesc(defLang)
	for _, key := range sortFields {
		cell := row.AddCell()
		cell.Value = fields[key]
		blog.V(5).Infof("key:%s value:%v, rid: %s", key, fields[key], rid)
	}
	return row
}

func setExcelTitleType(ctx context.Context, row *xlsx.Row, defLang lang.DefaultCCLanguageIf) *xlsx.Row {
	rid := util.ExtractRequestIDFromContext(ctx)
	fieldType := logics.GetPropertyFieldType(defLang)
	for _, key := range sortFields {
		cell := row.AddCell()
		cell.Value = fieldType[key]
		blog.V(5).Infof("key:%s value:%v, rid: %s", key, fieldType[key], rid)
	}
	return row
}

func setExcelRow(ctx context.Context, row *xlsx.Row, item interface{}) *xlsx.Row {
	rid := util.ExtractRequestIDFromContext(ctx)

	itemMap, ok := item.(map[string]interface{})
	if !ok {
		blog.V(5).Infof("failed to convert to map, rid: %s", rid)
		return row
	}

	// key is the object filed, value is the object filed value
	for _, key := range sortFields {

		cell := row.AddCell()
		// cell.SetValue([]string{"v1", "v2"})
		keyVal, ok := itemMap[key]
		if !ok {
			blog.Warnf("not fount the key(%s), skip it, rid: %s", key, rid)
			continue
		}
		blog.V(5).Infof("key:%s value:%v, rid: %s", key, keyVal, rid)
		if nil == keyVal {
			cell.SetString("")
			continue
		}
		switch t := keyVal.(type) {
		case bool:
			cell.SetBool(t)
		case string:
			if "\"\"" == t {
				cell.SetValue("")
			} else {
				cell.SetValue(t)
			}
		default:
			switch key {
			case common.BKOptionField:

				bOptions, err := json.Marshal(t)
				if nil != err {
					blog.Errorf("option format error:%v, rid: %s", t, rid)
					cell.SetValue("error info:" + err.Error())
				} else {
					cell.SetString(string(bOptions))
				}

			default:
				if nil != keyVal {
					cell.SetValue(t)
				}
			}
		}
	}

	return row
}

// ExportObjectBody TODO
type ExportObjectBody struct {
	BizID int64 `json:"bk_biz_id"`
}

// ExportObject export object
func (s *Service) ExportObject(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)

	webCommon.SetProxyHeader(c)

	objID := c.Param(common.BKObjIDField)

	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	requestBody := ExportObjectBody{}
	err := c.BindJSON(&requestBody)
	if err != nil {
		blog.Error("export model failed, parse request body to json failed, err: %v, rid: %s", err, rid)
		msg := fmt.Sprintf("invalid body, parse json failed, err: %+v", err)
		c.String(http.StatusBadRequest, msg)
		return
	}

	// get the all attribute of the object
	arrItems, err := s.Logics.GetObjectData(objID, c.Request.Header, requestBody.BizID)
	if nil != err {
		blog.Error("export model, but get object data failed, err: %v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebGetObjectFail, defErr.Errorf(common.CCErrWebGetObjectFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}

	// construct the excel file
	var file *xlsx.File
	var sheet *xlsx.Sheet

	file = xlsx.NewFile()
	sheet, err = file.AddSheet(objID)
	if err != nil {
		blog.Errorf("ExportObject failed, AddSheet failed, err: %s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(common.CCErrWebCreateEXCELFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}

	// set the title
	setExcelTitle(ctx, sheet.AddRow(), defLang)
	setExcelTitleType(ctx, sheet.AddRow(), defLang)
	setExcelSubTitle(sheet.AddRow())

	// add the value
	for _, item := range arrItems {

		innerRow := item.(map[string]interface{})
		blog.V(5).Infof("object attribute data :%+v, rid: %s", innerRow, rid)

		// set row value
		setExcelRow(ctx, sheet.AddRow(), innerRow)

	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if nil != err {
		blog.Warnf("os.Stat failed, will retry with os.MkdirAll, filename: %s, err: %+v, rid: %s", dirFileName, err, rid)
		if err := os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("os.MkdirAll failed, filename: %s, err: %+v, rid: %s", dirFileName, err, rid)
		}
	}
	fileName := fmt.Sprintf("%d_%s.xlsx", time.Now().UnixNano(), objID)
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)
	err = file.Save(dirFileName)
	if err != nil {
		blog.Errorf("ExportInst save file error:%s, rid: %s", err.Error(), rid)
		fmt.Printf(err.Error())
	}
	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("bk_cmdb_model_%s.xlsx", objID))
	c.File(dirFileName)

	if err := os.Remove(dirFileName); err != nil {
		blog.Errorf("os.Remove failed, filename: %s, err: %+v, rid: %s", dirFileName, err, rid)
	}

}

// SearchBusiness TODO
func (s *Service) SearchBusiness(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)
	webCommon.SetProxyHeader(c)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	query := new(params.SearchParams)
	err := c.BindJSON(&query)
	if err != nil {
		blog.Errorf("search business, but unmarshal body to json failed, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result:      false,
			Code:        common.CCErrCommJSONUnmarshalFailed,
			ErrMsg:      defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(),
			Permissions: nil,
		})
		return
	}

	// change the string query to regexp, only for frontend usage.
	for k, v := range query.Condition {
		field, ok := v.(string)
		if ok {
			query.Condition[k] = mapstr.MapStr{
				common.BKDBLIKE: params.SpecialCharChange(field),
				// insensitive with the character case.
				common.BKDBOPTIONS: "i",
			}
		}
	}
	ownerID := c.Request.Header.Get(common.BKHTTPOwnerID)
	biz, err := s.Engine.CoreAPI.ApiServer().SearchBiz(ctx, ownerID, c.Request.Header, query)
	if err != nil {
		blog.Error("search business, but request to api failed, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result:      false,
			Code:        common.CCErrCommHTTPDoRequestFailed,
			ErrMsg:      defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(),
			Permissions: nil,
		})
		return
	}

	if !biz.Result {
		if biz.Code == common.CCNoPermission {
			c.JSON(http.StatusOK, biz)
			return
		} else {
			c.JSON(http.StatusBadRequest, biz)
			return
		}
	}

	c.JSON(http.StatusOK, biz)
	return
}

// GetObjectInstanceCount TODO
func (s *Service) GetObjectInstanceCount(c *gin.Context) {
	header := c.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	ctx := util.NewContextFromGinContext(c)
	webCommon.SetProxyHeader(c)
	cond := &metadata.ObjectCountParams{}

	err := c.BindJSON(&cond)
	if err != nil {
		blog.Errorf("unmarshal body to json failed, err: %s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrCommJSONUnmarshalFailed, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(msg))
		return
	}

	resp, err := s.Logics.GetObjectCount(ctx, header, cond)
	if err != nil {
		blog.Errorf("get object count failed, err: %s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrCommHTTPDoRequestFailed, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(msg))
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

// BatchExportObject batch export object into yaml
func (s *Service) BatchExportObject(c *gin.Context) {
	header := c.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	ctx := util.NewContextFromGinContext(c)
	webCommon.SetProxyHeader(c)

	cond := new(metadata.BatchExportObject)
	err := c.BindJSON(cond)
	if err != nil {
		blog.Errorf("unmarshal body to json failed, err: %s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrCommJSONUnmarshalFailed, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(msg))
		return
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if err != nil {
		blog.Warnf("os.Stat failed, will retry with os.MkdirAll, filename: %s, err: %v, rid: %s", dirFileName, err, rid)
		if err := os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("os.MkdirAll failed, filename: %s, err: %v, rid: %s", dirFileName, err, rid)
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("save form data to local file failed, mkdir failed, err: %v", err))
			return
		}
	}

	if cond.FileName == "" {
		cond.FileName = fmt.Sprintf("batch_export_object_%d", time.Now().UnixNano())
	}

	fileDir := fmt.Sprintf("%s/%s_%d.zip", dirFileName, cond.FileName, time.Now().UnixNano())
	fzip, err := os.Create(fileDir)
	if err != nil {
		blog.Errorf("create zip file failed, err: %v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrCommHTTPDoRequestFailed, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(msg))
		return
	}

	defer func() {
		if err := os.Remove(fileDir); err != nil {
			blog.Errorf("os.Remove failed, filename: %s, err: %v, rid: %s", fileDir, err, rid)
		}
	}()

	zipw := zip.NewWriter(fzip)

	objRsp, err := s.Engine.CoreAPI.ApiServer().SearchObjectWithTotalInfo(ctx, header, cond)
	if err != nil {
		blog.Errorf("search object info to build yaml failed, cond: %v, err: %v, rid: %s", cond, err, rid)
		msg := getReturnStr(common.CCErrCommHTTPDoRequestFailed, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(msg))
		return
	}

	for objID, item := range objRsp.Object {
		yamlData, err := s.Logics.BuildExportYaml(header, cond.Expiration, item, "object")
		if err != nil {
			blog.Errorf("build yaml data failed, err: %v, rid: %s", err, rid)
			msg := getReturnStr(common.CCErrWebBuildZipFail, err.Error(), nil)
			c.Writer.Write([]byte(msg))
			return
		}

		fileName := fmt.Sprintf("%s_%d.yaml", objID, time.Now().Unix())
		s.Logics.BuildZipFile(header, zipw, fileName, cond.Password, yamlData)
	}

	if len(objRsp.Asst) != 0 {
		yamlData, err := s.Logics.BuildExportYaml(header, cond.Expiration, objRsp.Asst, "asst_kind")
		if err != nil {
			blog.Errorf("build yaml data failed, err: %v, rid: %s", err, rid)
			msg := getReturnStr(common.CCErrWebBuildZipFail, err.Error(), nil)
			c.Writer.Write([]byte(msg))
			return
		}

		fileName := fmt.Sprintf("asst_kind_%d.yaml", time.Now().Unix())
		s.Logics.BuildZipFile(header, zipw, fileName, cond.Password, yamlData)
	}

	zipw.Close()
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", cond.FileName))
	c.Writer.Header().Set("Content-Type", "application/octet-stream;charset=UTF-8")
	c.File(fileDir)
}

// BatchImportObjectAnalysis batch analysis object and asstkind yaml
func (s *Service) BatchImportObjectAnalysis(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	webCommon.SetProxyHeader(c)

	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	params := c.PostForm("params")
	cond := metadata.ZipFileAnalysis{}
	if len(params) != 0 {
		if err := json.Unmarshal([]byte(params), &cond); err != nil {
			blog.Errorf("params unmarshal error, err: %v, rid: %s", err, rid)
			msg := getReturnStr(common.CCErrCommParamsValueInvalidError,
				defErr.CCErrorf(common.CCErrCommParamsValueInvalidError, "params", err.Error()).Error(), nil)
			c.String(http.StatusOK, msg)
			return
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		blog.Errorf("get file from web form failed, err: %v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}

	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/import/"
	if _, err = os.Stat(dir); err != nil {
		blog.Warnf("os.Stat failed, filename: %s, err: %v, rid: %s", dir, err, rid)
		if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("os.MkdirAll failed, filename: %s, err: %v, rid: %s", dir, err, rid)
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("save form data to local file failed, mkdir failed, err: %v", err))
			return
		}
	}
	filePath := fmt.Sprintf("%s/batch_import_object-%d-%d.zip", dir, time.Now().UnixNano(), randNum)
	if err = c.SaveUploadedFile(file, filePath); err != nil {
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail,
			err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}
	defer func() {
		if err := os.Remove(filePath); err != nil {
			blog.Errorf("os.Remove failed, filename: %s, err: %v, rid: %s", filePath, err, rid)
		}
	}()

	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		blog.Errorf("open zip reader failed, err: %v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail,
			err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}

	defer zipReader.Close()
	result := &metadata.AnalysisResult{}
	for _, item := range zipReader.File {

		if item.FileInfo().IsDir() {
			continue
		}

		// file name start with '.' means hidden file, ignore
		if strings.HasPrefix(item.FileInfo().Name(), ".") {
			continue
		}

		errCode, err := s.Logics.GetDataFromZipFile(c.Request.Header, item, cond.Password, result)
		if err != nil {
			blog.Errorf("get data from zip file failed, err: %v, rid: %s", err, rid)
			msg := getReturnStr(errCode, err.Error(), nil)
			c.String(http.StatusOK, msg)
			return
		}
	}

	result.Result = true
	c.JSON(http.StatusOK, result)
}

// BatchImportObject batch import object
func (s *Service) BatchImportObject(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	webCommon.SetProxyHeader(c)
	ctx := util.NewContextFromGinContext(c)

	cond := new(metadata.BatchImportObject)
	err := c.BindJSON(cond)
	if err != nil {
		blog.Errorf("unmarshal body to json failed, err: %v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrCommJSONUnmarshalFailed, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(msg))
		return
	}

	objInfo := metadata.ImportObjects{Objects: cond.Object, Asst: cond.Asst}
	if _, err := s.Engine.CoreAPI.ApiServer().CreateManyObject(ctx, c.Request.Header, objInfo); err != nil {
		blog.Errorf("create many object failed, err: %v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrTopoModuleCreateFailed, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(msg))
		return
	}

	c.JSON(http.StatusOK, metadata.Response{BaseResp: metadata.BaseResp{Result: true, ErrMsg: "success"}})
}
