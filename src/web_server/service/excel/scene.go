/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package excel

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/service/excel/core"
	"configcenter/src/web_server/service/excel/operator"
	"configcenter/src/web_server/service/excel/operator/inst/exporter"
	"configcenter/src/web_server/service/excel/operator/inst/importer"
	"configcenter/src/web_server/service/excel/operator/model"

	"github.com/gin-gonic/gin"
)

// FileType 文件类型
type FileType string

const (
	// FileTypeXlsx 文件格式xlsx
	FileTypeXlsx FileType = "xlsx"
	// FileTypeXls 文件格式xls
	FileTypeXls FileType = "xls"
	// FileTypeZip 文件格式zip
	FileTypeZip FileType = "zip"
	// FileTypeYaml 文件格式yaml
	FileTypeYaml FileType = "yaml"
)

// ImportType 导入类型
type ImportType string

const (
	// ImportTypeInst 导入类型为导入实例
	ImportTypeInst ImportType = "importInst"
	// ImportTypeObject 导入类型为导入模型
	ImportTypeObject ImportType = "importObject"
	// ImportTypeObjectYaml 导入类型为导入模型yaml文件
	ImportTypeObjectYaml ImportType = "importObjectYaml"
	// ImportTypeObjectAttr 导入类型为导入模型属性字段
	ImportTypeObjectAttr ImportType = "importObjectAttr"
)

// ImportTypeMap 导入类型与文件类型对应关系map
var ImportTypeMap = map[ImportType][]FileType{
	ImportTypeInst:       {FileTypeXlsx, FileTypeXls},
	ImportTypeObjectAttr: {FileTypeXlsx, FileTypeXls},
	ImportTypeObject:     {FileTypeZip},
	ImportTypeObjectYaml: {FileTypeYaml},
}

// BuildTemplate build excel download template
func (s *service) BuildTemplate(c *gin.Context) {
	objID := c.Param(common.BKObjIDField)
	kit := rest.NewKitFromHeader(c.Request.Header, s.engine.CCErr)

	// 1. 创建excel模版
	dir := fmt.Sprintf("%s/template", webCommon.ResourcePath)
	randNum := rand.Uint32()
	filePath := fmt.Sprintf("%s/%stemplate-%d-%d.xlsx", dir, objID, time.Now().UnixNano(), randNum)

	client := &core.Client{ApiClient: s.apiCli}
	baseOp, err := operator.NewBaseOp(operator.FilePath(filePath), operator.Client(client), operator.ObjID(objID),
		operator.Kit(kit), operator.Language(s.engine.Language))
	if err != nil {
		blog.Errorf("create excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}

	tmplOp, err := exporter.NewTmplOp(exporter.BaseOperator(baseOp))
	if err != nil {
		blog.Errorf("create excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}

	if err := tmplOp.BuildHeader(); err != nil {
		blog.Errorf("create excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}

	if err := tmplOp.Close(); err != nil {
		blog.Errorf("close excel io failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}

	// 2. 将excel模版文件返回，并删除临时文件
	if objID == common.BKInnerObjIDHost {
		addDownExcelHttpHeader(c, "bk_cmdb_import_host.xlsx")
	} else {
		addDownExcelHttpHeader(c, fmt.Sprintf("bk_cmdb_inst_%s.xlsx", objID))
	}

	c.File(filePath)

	if err := tmplOp.Clean(); err != nil {
		blog.Errorf("clean excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}
}

// ExportInst export instance
func (s *service) ExportInst(c *gin.Context) {
	objID := c.Param(common.BKObjIDField)
	s.exportInstFunc(c, objID)
}

// ExportHost export host
func (s *service) ExportHost(c *gin.Context) {
	s.exportInstFunc(c, common.BKInnerObjIDHost)
}

// ExportBiz export biz
func (s *service) ExportBiz(c *gin.Context) {
	s.exportInstFunc(c, common.BKInnerObjIDApp)
}

// ExportProject export project
func (s *service) ExportProject(c *gin.Context) {
	s.exportInstFunc(c, common.BKInnerObjIDProject)
}

func (s *service) exportInstFunc(c *gin.Context, objID string) {
	kit := rest.NewKitFromHeader(c.Request.Header, s.engine.CCErr)
	input := exporter.GetExportParamInterface(objID)

	if err := c.BindJSON(input); err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommJSONUnmarshalFailed, ErrMsg: err.Error()})
		return
	}
	lang := s.engine.Language.CreateDefaultCCLanguageIf(httpheader.GetLanguage(kit.Header))
	if err := input.Validate(kit, lang); err != nil {
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebGetObjectFail, err.Error()))
		return
	}

	// 1. 初始化导出excel对象
	dir := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	filePath := fmt.Sprintf("%s/%s", dir, fmt.Sprintf("%dinst.xlsx", time.Now().UnixNano()))

	client := &core.Client{ApiClient: s.apiCli, GinCtx: c}
	baseOp, err := operator.NewBaseOp(operator.FilePath(filePath), operator.Client(client), operator.ObjID(objID),
		operator.Kit(kit), operator.Language(s.engine.Language))
	if err != nil {
		blog.Errorf("create excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrCommExcelTemplateFailed, err.Error()))
		return
	}

	tmplOp, err := exporter.NewTmplOp(exporter.BaseOperator(baseOp))
	if err != nil {
		blog.Errorf("create excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrCommExcelTemplateFailed, err.Error()))
		return
	}

	operator, err := exporter.NewExporter(exporter.TmplOperator(tmplOp), exporter.ExportParam(input))
	if err != nil {
		blog.Errorf("create excel exporter failed, err: %v, rid: %s", err, kit.Rid)
		c.String(http.StatusInternalServerError, fmt.Errorf("create exporter failed, err: %+v", err).Error())
		return
	}

	// 2. 导出实例数据到excel
	if err := operator.Export(); err != nil {
		blog.Errorf("export instance data failed, err: %v, rid: %s", err, kit.Rid)
		c.String(http.StatusInternalServerError, fmt.Errorf("export instance data failed, err: %+v", err).Error())
		return
	}

	if err := operator.Close(); err != nil {
		blog.Errorf("close excel io failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebGetObjectFail, err.Error()))
		return
	}

	// 3. 将excel文件返回，并删除临时文件
	switch objID {
	case common.BKInnerObjIDHost:
		addDownExcelHttpHeader(c, "bk_cmdb_export_host.xlsx")
	case common.BKInnerObjIDApp:
		addDownExcelHttpHeader(c, "bk_cmdb_export_biz.xlsx")
	case common.BKInnerObjIDProject:
		addDownExcelHttpHeader(c, "bk_cmdb_export_project.xlsx")
	default:
		addDownExcelHttpHeader(c, fmt.Sprintf("bk_cmdb_export_inst_%s.xlsx", objID))
	}

	c.File(filePath)

	if err := operator.Clean(); err != nil {
		blog.Errorf("clean excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebGetObjectFail, err.Error()))
		return
	}
}

const param = "params"

// AddInst add instance
func (s *service) AddInst(c *gin.Context) {
	objID := c.Param(common.BKObjIDField)
	s.importInstFunc(c, objID, core.AddInst)
}

// AddHost add host
func (s *service) AddHost(c *gin.Context) {
	s.importInstFunc(c, common.BKInnerObjIDHost, core.AddHost)
}

// UpdateHost update host
func (s *service) UpdateHost(c *gin.Context) {
	s.importInstFunc(c, common.BKInnerObjIDHost, core.UpdateHost)
}

// importInstFunc import instance function
func (s *service) importInstFunc(c *gin.Context, objID string, handleType core.HandleType) {
	kit := rest.NewKitFromHeader(c.Request.Header, s.engine.CCErr)
	params := c.PostForm(param)
	if params == "" {
		blog.Errorf("not found params value, rid: %s", kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrCommParamsNeedSet, param))
		return
	}

	var input importer.ImportParamI
	switch handleType {
	case core.AddHost:
		input = &importer.AddHostParam{}
	case core.UpdateHost:
		input = &importer.UpdateHostParam{}
	case core.AddInst:
		input = &importer.InstParam{}
	}

	if err := json.Unmarshal([]byte(params), input); err != nil {
		blog.Errorf("params unmarshal error, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrCommParamsValueInvalidError, params, err.Error()))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		blog.Errorf("get file failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebFileNoFound))
		return
	}
	if err := VerifyFileType(ImportTypeInst, file.Filename, kit.Rid); err != nil {
		blog.Errorf("file type verify failed, err: %v, fileName: %s, rid: %s", err, file.Filename, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrInvalidFileTypeFail, err.Error()))
		return
	}

	dir := webCommon.ResourcePath + "/import/"
	if _, err = os.Stat(dir); err != nil {
		blog.Warnf("os.Stat failed, filename: %s, will retry with os.MkdirAll, err: %v, rid: %s", dir, err, kit.Rid)
		if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("os.MkdirAll failed, filename: %s, err: %v, rid: %s", dir, err, kit.Rid)
		}
	}

	filePath := fmt.Sprintf("%s/importinsts-%d-%d.xlsx", dir, time.Now().UnixNano(), rand.Uint32())
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebFileSaveFail))
		return
	}

	client := &core.Client{ApiClient: s.apiCli}
	baseOp, err := operator.NewBaseOp(operator.FilePath(filePath), operator.Client(client), operator.ObjID(objID),
		operator.Kit(kit), operator.Language(s.engine.Language))
	if err != nil {
		blog.Errorf("create importer failed, err: %v, rid: %s", err, kit.Rid)
		c.String(http.StatusInternalServerError, fmt.Errorf("create importer failed, err: %+v", err).Error())
		return
	}

	op, err := importer.NewImporter(importer.BaseOperator(baseOp), importer.Param(input))
	if err != nil {
		blog.Errorf("create importer failed, err: %v, rid: %s", err, kit.Rid)
		c.String(http.StatusInternalServerError, fmt.Errorf("create importer failed, err: %+v", err).Error())
		return
	}

	result, err := op.Handle()
	if err != nil {
		blog.Errorf("handle excel import request failed, err: %v, rid: %s", err, kit.Rid)
		c.String(http.StatusInternalServerError, fmt.Errorf("handle import request failed, err: %+v", err).Error())
		return
	}
	if err := op.Clean(); err != nil {
		blog.Errorf("clean importer resource failed, err: %v, rid: %s", err, kit.Rid)
		c.String(http.StatusInternalServerError, fmt.Errorf("clean importer resource failed, err: %+v", err).Error())
		return
	}

	c.JSON(http.StatusOK, metadata.NewSuccessResp(result))
}

// ExportObject export object
func (s *service) ExportObject(c *gin.Context) {
	kit := rest.NewKitFromHeader(c.Request.Header, s.engine.CCErr)
	objID := c.Param(common.BKObjIDField)

	dir := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	filePath := fmt.Sprintf("%s/%d_%s.xlsx", dir, time.Now().UnixNano(), objID)

	client := &core.Client{ApiClient: s.apiCli}
	baseOp, err := operator.NewBaseOp(operator.FilePath(filePath), operator.Client(client), operator.ObjID(objID),
		operator.Kit(kit), operator.Language(s.engine.Language))
	if err != nil {
		blog.Errorf("create model operator failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebOpenFileFail, err.Error()))
		return
	}

	modelOp, err := model.NewOp(model.BaseOperator(baseOp))
	if err != nil {
		blog.Errorf("create model operator failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebOpenFileFail, err.Error()))
		return
	}

	// 导出模型数据到excel
	if err := modelOp.Export(); err != nil {
		blog.Errorf("export model data failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebCreateEXCELFail, err.Error()))
		return
	}

	if err := modelOp.Close(); err != nil {
		blog.Errorf("close excel io failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrWebFileContentFail, ErrMsg: err.Error()})
		return
	}

	// 将excel文件返回，并删除临时文件
	addDownExcelHttpHeader(c, fmt.Sprintf("bk_cmdb_model_%s.xlsx", objID))
	c.File(filePath)

	if err := modelOp.Clean(); err != nil {
		blog.Errorf("clean model operator resource failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrWebFileContentFail, ErrMsg: err.Error()})
		return
	}
}

// ImportObject import object attribute
func (s *service) ImportObject(c *gin.Context) {
	kit := rest.NewKitFromHeader(c.Request.Header, s.engine.CCErr)
	objID := c.Param(common.BKObjIDField)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebFileNoFound))
		return
	}

	if err := VerifyFileType(ImportTypeObjectAttr, file.Filename, kit.Rid); err != nil {
		blog.Errorf("file type verify failed, err: %v, fileName: %s, rid: %s", err, file.Filename, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrInvalidFileTypeFail, err.Error()))
		return
	}

	dir := webCommon.ResourcePath + "/import/"
	if _, err = os.Stat(dir); err != nil {
		blog.Warnf("os.Stat failed, filename: %s, will retry with os.MkdirAll, err: %v, rid: %s", dir, err, kit.Rid)
		if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			blog.Errorf("os.MkdirAll failed, filename: %s, err: %v, rid: %s", dir, err, kit.Rid)
		}
	}

	filePath := fmt.Sprintf("%s/importinsts-%d-%d.xlsx", dir, time.Now().UnixNano(), rand.Uint32())
	if err = c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebFileSaveFail))
		return
	}

	client := &core.Client{ApiClient: s.apiCli}
	baseOp, err := operator.NewBaseOp(operator.FilePath(filePath), operator.Client(client), operator.ObjID(objID),
		operator.Kit(kit), operator.Language(s.engine.Language))
	if err != nil {
		blog.Errorf("create importer failed, err: %v, rid: %s", err, kit.Rid)
		c.String(http.StatusInternalServerError, fmt.Errorf("create importer failed, err: %+v", err).Error())
		return
	}

	modelOp, err := model.NewOp(model.BaseOperator(baseOp))
	if err != nil {
		blog.Errorf("create model operator failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebOpenFileFail, err.Error()))
		return
	}

	result, err := modelOp.Import()
	if err != nil {
		blog.Errorf("export model data failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, getErrResp(kit, common.CCErrWebFileContentFail, err.Error()))
		return
	}

	if err := modelOp.Close(); err != nil {
		blog.Errorf("close excel io failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrWebFileContentFail, ErrMsg: err.Error()})
		return
	}

	if err := modelOp.Clean(); err != nil {
		blog.Errorf("close excel io failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrWebFileContentFail, ErrMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func getErrResp(kit *rest.Kit, code int, params ...string) metadata.BaseResp {
	return metadata.BaseResp{
		Code:   code,
		ErrMsg: kit.CCError.CCErrorf(code, params).Error(),
	}
}

func addDownExcelHttpHeader(c *gin.Context, name string) {
	if strings.HasSuffix(name, ".xls") {
		c.Header("Content-Type", "application/vnd.ms-excel")
	} else {
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	}
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Disposition", "attachment; filename="+name) // 文件名
	c.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

// VerifyFileType 校验导入的文件格式
func VerifyFileType(importType ImportType, fileName, rid string) error {
	if fileName == "" {
		blog.Errorf("verify file type failed, file name is empty, rid: %s", rid)
		return errors.NewCCError(common.CCErrInvalidFileTypeFail, "file name is empty")
	}

	lastIndex := strings.LastIndex(fileName, ".")
	if lastIndex == -1 {
		blog.Errorf("verify file type failed, file type is empty, rid: %s", rid)
		return errors.NewCCError(common.CCErrInvalidFileTypeFail, "file type is empty")
	}

	fileType := fileName[lastIndex+1:]
	fileTypeArr, exist := ImportTypeMap[importType]
	if !exist {
		blog.Errorf("verify file type failed, unknown import type %s, rid: %s", importType, rid)
		return errors.NewCCError(common.CCErrInvalidFileTypeFail, "unknown import type")
	}

	for _, ft := range fileTypeArr {
		if ft == FileType(fileType) {
			return nil
		}
	}
	blog.Errorf("verify file type failed, unknown import type %s, rid: %s", fileType, rid)
	return errors.NewCCError(common.CCErrInvalidFileTypeFail, "file type is illegality")
}
