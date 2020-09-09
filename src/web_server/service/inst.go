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
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"
)

// ImportInst import inst
func (s *Service) ImportInst(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	webCommon.SetProxyHeader(c)
	objID := c.Param(common.BKObjIDField)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	file, err := c.FormFile("file")
	if nil != err {
		msg := getReturnStr(common.CCErrWebFileNoFound, defErr.Error(common.CCErrWebFileNoFound).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	modelBizID, err := parseModelBizID(c.PostForm(common.BKAppIDField))
	if err != nil {
		msg := getReturnStr(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/import/"
	_, err = os.Stat(dir)
	if nil != err {
		if err != nil {
			blog.Warnf("os.Stat failed, filename: %s, will retry with os.MkdirAll, err: %+v, rid: %s", dir, err, rid)
		}
		if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("os.MkdirAll failed, filename: %s, err: %+v, rid: %s", dir, err, rid)
		}
	}
	filePath := fmt.Sprintf("%s/importinsts-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	err = c.SaveUploadedFile(file, filePath)
	if nil != err {
		msg := getReturnStr(common.CCErrWebFileSaveFail, defErr.Errorf(common.CCErrWebFileSaveFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}
	defer func() {
		if err := os.Remove(filePath); err != nil {
			blog.Errorf("os.Remove failed, filename: %s, err: %+v, rid: %s", filePath, err, rid)
		}
	}()
	f, err := xlsx.OpenFile(filePath)
	if nil != err {
		msg := getReturnStr(common.CCErrWebOpenFileFail, defErr.Errorf(common.CCErrWebOpenFileFail, err.Error()).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	data, errCode, err := s.Logics.ImportInsts(context.Background(), f, objID, c.Request.Header, defLang, modelBizID)

	if nil != err {
		msg := getReturnStr(errCode, err.Error(), data)
		c.String(http.StatusOK, string(msg))
		return
	}

	c.String(http.StatusOK, getReturnStr(0, "", data))
}

// ExportInst export inst
func (s *Service) ExportInst(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)
	webCommon.SetProxyHeader(c)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defLang := s.Language.CreateDefaultCCLanguageIf(language)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)
	pheader := c.Request.Header

	ownerID := c.Param(common.BKOwnerIDField)
	objID := c.Param(common.BKObjIDField)
	instIDStr := c.PostForm(common.BKInstIDField)
	customFieldsStr := c.PostForm(common.ExportCustomFields)

	modelBizID, err := parseModelBizID(c.PostForm(common.BKAppIDField))
	if err != nil {
		msg := getReturnStr(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), nil)
		c.String(http.StatusOK, string(msg))
		return
	}

	kvMap := mapstr.MapStr{}
	instInfo, err := s.Logics.GetInstData(ownerID, objID, instIDStr, pheader, kvMap)
	if err != nil {
		msg := getReturnStr(common.CCErrWebGetObjectFail, defErr.Errorf(common.CCErrWebGetObjectFail, err.Error()).Error(), nil)
		fmt.Println("return msg: ", msg)
		c.String(http.StatusForbidden, msg)
		return
	}

	var file *xlsx.File

	file = xlsx.NewFile()

	customFields := logics.GetCustomFields(nil, customFieldsStr)
	fields, err := s.Logics.GetObjFieldIDs(objID, nil, customFields, pheader, modelBizID)
	if err != nil {
		blog.Errorf("export object instance, but get object:%s attribute field failed, err: %v, rid: %s", objID, err, rid)
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}

	err = s.Logics.BuildExcelFromData(ctx, objID, fields, nil, instInfo, file, pheader, modelBizID)
	if nil != err {
		blog.Errorf("ExportHost object:%s error:%s, rid: %s", objID, err.Error(), rid)
		reply := getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(common.CCErrCommExcelTemplateFailed, objID).Error(), nil)
		c.Writer.Write([]byte(reply))
		return
	}
	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	_, err = os.Stat(dirFileName)
	if nil != err {
		blog.Warnf("os.Stat failed, filename: %s, will retry with os.MkdirAll, err: %+v, rid: %s", dirFileName, err, rid)
		if err := os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("os.MkdirAll failed, filename: %s, err: %+v, rid: %s", dirFileName, err, rid)
		}
	}
	fileName := fmt.Sprintf("%dinst.xlsx", time.Now().UnixNano())
	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fileName)
	logics.ProductExcelCommentSheet(ctx, file, defLang)
	err = file.Save(dirFileName)
	if err != nil {
		blog.Errorf("ExportInst save file error:%s, rid: %s", err.Error(), rid)
		if err != nil {
			blog.Errorf("ExportInst save file error:%s, rid: %s", err.Error(), rid)
			reply := getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)
			c.Writer.Write([]byte(reply))
			return
		}
	}
	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("bk_cmdb_export_inst_%s.xlsx", objID))
	c.File(dirFileName)
	if err := os.Remove(dirFileName); err != nil {
		blog.Errorf("remove file %s failed, err: %+v, rid: %s", dirFileName, err, rid)
	}
}
