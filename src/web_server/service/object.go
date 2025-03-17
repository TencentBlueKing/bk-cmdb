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
	"encoding/json"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user"
	"configcenter/src/web_server/service/excel"

	"github.com/alexmullins/zip"
	"github.com/gin-gonic/gin"
)

// GetObjectInstanceCount TODO
func (s *Service) GetObjectInstanceCount(c *gin.Context) {
	header := c.Request.Header
	rid := httpheader.GetRid(header)
	ctx := util.NewContextFromGinContext(c)
	if err := webCommon.SetProxyHeader(c); err != nil {
		user := user.NewUser(*s.Config, s.Engine, s.CacheCli)
		url := user.GetLoginUrl(c)
		c.Redirect(302, url)
	}
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
	rid := httpheader.GetRid(header)
	if err := webCommon.SetProxyHeader(c); err != nil {
		user := user.NewUser(*s.Config, s.Engine, s.CacheCli)
		url := user.GetLoginUrl(c)
		c.Redirect(302, url)
	}

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
	s.exportObj(c, cond, dirFileName)
}

func (s *Service) exportObj(c *gin.Context, cond *metadata.BatchExportObject, dirFileName string) {
	header := c.Request.Header
	rid := httpheader.GetRid(header)
	ctx := util.NewContextFromGinContext(c)

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

	objRsp, err := s.ApiCli.SearchObjectWithTotalInfo(ctx, header, cond)
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
	rid := httpheader.GetRid(c.Request.Header)
	if err := webCommon.SetProxyHeader(c); err != nil {
		user := user.NewUser(*s.Config, s.Engine, s.CacheCli)
		url := user.GetLoginUrl(c)
		c.Redirect(302, url)
	}

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
	if err := excel.VerifyFileType(excel.ImportTypeObject, file.Filename, rid); err != nil {
		blog.Errorf("ImportObject failed, file type verify failed, err: %v, fileName: %s, rid: %s",
			err, file.Filename, rid)
		msg := getReturnStr(common.CCErrInvalidFileTypeFail, err.Error(), nil)
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

	result, err := s.getAnalysisResult(c, filePath, cond, file)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, result)
}

func (s *Service) getAnalysisResult(c *gin.Context, filePath string,
	cond metadata.ZipFileAnalysis, file *multipart.FileHeader) (*metadata.AnalysisResult, error) {
	rid := httpheader.GetRid(c.Request.Header)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		blog.Errorf("open zip reader failed, err: %v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail,
			err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
		return nil, err
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
		if err := excel.VerifyFileType(excel.ImportTypeObjectYaml, item.FileInfo().Name(), rid); err != nil {
			blog.Errorf("ImportObjectYaml failed, file type verify failed, err: %v, fileName: %s, rid: %s",
				err, file.Filename, rid)
			msg := getReturnStr(common.CCErrInvalidFileTypeFail, err.Error(), nil)
			c.String(http.StatusOK, msg)
			return nil, err
		}

		errCode, err := s.Logics.GetDataFromZipFile(c.Request.Header, item, cond.Password, result)
		if err != nil {
			blog.Errorf("get data from zip file failed, err: %v, rid: %s", err, rid)
			msg := getReturnStr(errCode, err.Error(), nil)
			c.String(http.StatusOK, msg)
			return nil, err
		}
	}
	result.Result = true
	return result, nil
}

// BatchImportObject batch import object
func (s *Service) BatchImportObject(c *gin.Context) {
	rid := httpheader.GetRid(c.Request.Header)
	if err := webCommon.SetProxyHeader(c); err != nil {
		user := user.NewUser(*s.Config, s.Engine, s.CacheCli)
		url := user.GetLoginUrl(c)
		c.Redirect(302, url)
	}
	ctx := util.NewContextFromGinContext(c)

	cond := new(metadata.BatchImportObject)
	err := c.BindJSON(cond)
	if err != nil {
		blog.Errorf("unmarshal body to json failed, err: %v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrCommJSONUnmarshalFailed, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(msg))
		return
	}

	if len(cond.Object) == 0 {
		c.JSON(http.StatusOK, metadata.Response{BaseResp: metadata.BaseResp{Result: true, ErrMsg: "success"}})
		return
	}

	objInfo := metadata.ImportObjects{Objects: cond.Object, Asst: cond.Asst}
	if _, err := s.ApiCli.CreateManyObject(ctx, c.Request.Header, objInfo); err != nil {
		blog.Errorf("create many object failed, err: %v, rid: %s", err, rid)
		msg := getReturnStr(common.CCErrTopoModuleCreateFailed, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(msg))
		return
	}

	c.JSON(http.StatusOK, metadata.Response{BaseResp: metadata.BaseResp{Result: true, ErrMsg: "success"}})
}
