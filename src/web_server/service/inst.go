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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"

	"github.com/gin-gonic/gin"
	"github.com/rentiansheng/xlsx"
)

type excelExportInstInput struct {
	// 导出的实例字段
	CustomFields []string `json:"export_custom_fields"`
	// 指定需要导出的实例ID, 设置本参数后，
	InstIDArr []int64 `json:"bk_inst_ids"`
	// Deprecated 兼容，历史原因
	AppID int64 `json:"bk_biz_id"`

	// 用来限定导出关联关系，map[bk_obj_id]object_unique_id 2021年05月17日
	AssociationCond map[string]int64 `json:"association_condition"`

	// 用来限定当前操作对象导出数据的时候，需要使用的唯一校验关系，
	// 自关联的时候，规定左边对象使用到的唯一索引
	ObjectUniqueID int64 `json:"object_unique_id"`
}

type excelImportInstInput struct {
	BizID  int64 `json:"bk_biz_id"`
	OpType int64 `json:"op"`
	// 用来限定导出关联关系，map[bk_obj_id]object_unique_id 2021年05月17日
	AssociationCond map[string]int64 `json:"association_condition"`
	// 用来限定当前操作对象导出数据的时候，需要使用的唯一校验关系，
	// 自关联的时候，规定左边对象使用到的唯一索引
	ObjectUniqueID int64 `json:"object_unique_id"`
}

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
	params := c.PostForm("params")
	if params == "" {
		blog.ErrorJSON("ImportHost failed, not found params value, rid: %s", rid)
		msg := getReturnStr(common.CCErrCommParamsNeedSet,
			defErr.CCErrorf(common.CCErrCommParamsNeedSet, "params").Error(), nil)
		c.String(http.StatusOK, msg)
		return
	}
	inputJSON := &excelImportInstInput{}
	if err := json.Unmarshal([]byte(params), inputJSON); err != nil {
		blog.ErrorJSON("ImportHost failed, params unmarshal error, err: %s, rid: %s", err.Error(), rid)
		msg := getReturnStr(common.CCErrCommParamsValueInvalidError,
			defErr.CCErrorf(common.CCErrCommParamsValueInvalidError, "params", err.Error()).Error(), nil)
		c.String(http.StatusOK, msg)
	}

	randNum := rand.Uint32()
	dir := webCommon.ResourcePath + "/import/"
	_, err = os.Stat(dir)
	if err != nil {
		if err != nil {
			blog.Warnf("os.Stat failed, filename: %s, will retry with os.MkdirAll, err: %+v, rid: %s", dir, err, rid)
		}
		if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("os.MkdirAll failed, filename: %s, err: %+v, rid: %s", dir, err, rid)
		}
	}
	filePath := fmt.Sprintf("%s/importinsts-%d-%d.xlsx", dir, time.Now().UnixNano(), randNum)
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
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

	data, errCode, err := s.Logics.ImportInsts(context.Background(), f, objID, c.Request.Header, defLang,
		inputJSON.BizID, inputJSON.OpType, inputJSON.AssociationCond, inputJSON.ObjectUniqueID)

	if err != nil {
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

	input := &excelExportInstInput{}
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

	// ownerID := c.Param(common.BKOwnerIDField)
	objID := c.Param(common.BKObjIDField)

	modelBizID := input.AppID

	instInfo, err := s.Logics.GetInstData(objID, input.InstIDArr, pheader)
	if err != nil {
		msg := getReturnStr(common.CCErrWebGetObjectFail, defErr.Errorf(common.CCErrWebGetObjectFail,
			err.Error()).Error(), nil)
		blog.ErrorJSON("get inst data error. err: %s, inst id: %s, rid: %s", err.Error(), input.InstIDArr, rid)
		c.String(http.StatusForbidden, msg)
		return
	}

	customFields := logics.GetCustomFields(nil, input.CustomFields)
	fields, err := s.Logics.GetObjFieldIDs(objID, nil, customFields, pheader, modelBizID,
		common.HostAddMethodExcelDefaultIndex)
	if err != nil {
		blog.Errorf("get object:%s attribute field failed, err: %v, rid: %s", objID, err, rid)
		_, _ = c.Writer.Write([]byte(getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(
			common.CCErrCommExcelTemplateFailed, objID).Error(), nil)))
		return
	}

	usernameMap, propertyList, err := s.getUsernameMapWithPropertyList(c, objID, instInfo)
	if err != nil {
		blog.Errorf("ExportInst failed, get username map and property list failed, err: %+v, rid: %s", err, rid)
		_, _ = c.Writer.Write([]byte(getReturnStr(common.CCErrWebGetUsernameMapFail, defErr.Errorf(
			common.CCErrWebGetUsernameMapFail, objID).Error(), nil)))
	}

	org, orgPropertyList, err := s.getDepartment(c, objID)
	if err != nil {
		blog.Errorf("get department map and property list failed, err: %+v, rid: %s", err, rid)
		_, _ = c.Writer.Write([]byte(getReturnStr(common.CCErrWebGetDepartmentMapFail, defErr.Errorf(
			common.CCErrWebGetDepartmentMapFail, err.Error()).Error(), nil)))
	}

	file := xlsx.NewFile()
	err = s.Logics.BuildExcelFromData(ctx, objID, fields, nil, instInfo, file, pheader, modelBizID, usernameMap,
		propertyList, org, orgPropertyList, input.AssociationCond, input.ObjectUniqueID)
	if nil != err {
		blog.Errorf("ExportHost object:%s error:%s, rid: %s", objID, err.Error(), rid)
		_, _ = c.Writer.Write([]byte(getReturnStr(common.CCErrCommExcelTemplateFailed, defErr.Errorf(
			common.CCErrCommExcelTemplateFailed, objID).Error(), nil)))
		return
	}

	dirFileName := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	if _, err = os.Stat(dirFileName); err != nil {
		blog.Warnf("os.Stat failed, filename: %s, will retry with os.MkdirAll, err: %+v, rid: %s", dirFileName, err, rid)
		if err := os.MkdirAll(dirFileName, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("os.MkdirAll failed, filename: %s, err: %+v, rid: %s", dirFileName, err, rid)
		}
	}

	dirFileName = fmt.Sprintf("%s/%s", dirFileName, fmt.Sprintf("%dinst.xlsx", time.Now().UnixNano()))
	logics.ProductExcelCommentSheet(ctx, file, defLang)
	if err = file.Save(dirFileName); err != nil {
		blog.Errorf("ExportInst save file error:%s, rid: %s", err.Error(), rid)
		_, _ = c.Writer.Write([]byte(getReturnStr(common.CCErrWebCreateEXCELFail, defErr.Errorf(
			common.CCErrCommExcelTemplateFailed, err.Error()).Error(), nil)))
		return
	}
	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("bk_cmdb_export_inst_%s.xlsx", objID))
	c.File(dirFileName)
	if err := os.Remove(dirFileName); err != nil {
		blog.Errorf("remove file %s failed, err: %+v, rid: %s", dirFileName, err, rid)
	}
}
