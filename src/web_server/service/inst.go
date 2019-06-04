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
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"
)

// ImportInst import inst
func (s *Service) ImportInst(c *gin.Context) {
	logics.SetProxyHeader(c)
	objID := c.Param(common.BKObjIDField)
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	file, err := c.FormFile("file")
	if nil != err {
		msg := getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/import/"
	_, err = os.Stat(dir)
	if nil != err {
		os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	}
	filePath := fmt.Sprintf("%s/importinsts-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	err = c.SaveUploadedFile(file, filePath)
	if nil != err {
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	defer os.Remove(filePath) //delete file
	f, err := xlsx.OpenFile(filePath)
	if nil != err {
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	data, errCode, err := s.Logics.ImportInsts(context.Background(), f, objID, c.Request.Header, defLang)

	if nil != err {
		msg := getReturnStr(errCode, err.Error(), data)
		c.String(http.StatusOK, string(msg))
		return
	}

	c.String(http.StatusOK, getReturnStr(0, "", data))
}

// ExportInst export inst
func (s *Service) ExportInst(c *gin.Context) {
	logics.SetProxyHeader(c)
	language := logics.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)
	pheader := c.Request.Header

	ownerID := c.Param(common.BKOwnerIDField)
	objID := c.Param(common.BKObjIDField)
	instIDStr := c.PostForm(common.BKInstIDField)

	kvMap := mapstr.MapStr{}
	instInfo, err := s.Logics.GetInstData(ownerID, objID, instIDStr, pheader, kvMap)

	if err != nil {
		blog.Error(err.Error())
		msg := getReturnStr(common.CCErrWebGetObjectFail, defErr.Errorf(common.CCErrWebGetObjectFail, err.Error()).Error(), nil)

		c.String(http.StatusInternalServerError, msg, nil)
		return
	}

	var file *xlsx.File

	file = xlsx.NewFile()

	fields, err := s.Logics.GetObjFieldIDs(objID, nil, pheader)
	if err != nil {
		blog.Errorf("ExportInst object:%s error:%s", objID, err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, err.Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}
	err = s.Logics.BuildExcelFromData(context.Background(), objID, fields, nil, instInfo, file, pheader)
	if nil != err {
		blog.Errorf("ExportInst object:%s error:%s", objID, err.Error())
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}
	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if nil != err {
		os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm)
	}
	fileName := fmt.Sprintf("%dinst.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)
	//fileName := fmt.Sprintf("tmp/%s_inst.xls", time.Now().UnixNano())
	logics.ProductExcelCommentSheet(file, defLang)
	err = file.Save(dirFileName)
	if err != nil {
		blog.Errorf("ExportInst save file error:%s", err.Error())
		if err != nil {
			blog.Errorf("ExportInst save file error:%s", err.Error())
			reply := getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)
			c.Writer.Write([]byte(reply))
			return
		}
	}
	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("bk_cmdb_export_inst_%s.xlsx", objID))
	c.File(dirFileName)

	os.Remove(dirFileName)
}
