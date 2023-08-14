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
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/logics"
	"configcenter/src/web_server/service/excel/operator"

	"github.com/gin-gonic/gin"
)

// BuildTemplate build excel download template
func (s *service) BuildTemplate(c *gin.Context) {
	objID := c.Param(common.BKObjIDField)
	kit := rest.NewKitFromHeader(c.Request.Header, s.engine.CCErr)

	// 1. 创建excel模版
	dir := fmt.Sprintf("%s/template", webCommon.ResourcePath)
	randNum := rand.Uint32()
	filePath := fmt.Sprintf("%s/%stemplate-%d-%d.xlsx", dir, objID, time.Now().UnixNano(), randNum)

	tmplOp, err := operator.NewTmplOp(operator.FilePath(filePath), operator.Dao(s.dao), operator.ObjID(objID),
		operator.Kit(kit), operator.Language(s.engine.Language))
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
		logics.AddDownExcelHttpHeader(c, "bk_cmdb_import_host.xlsx")
	} else {
		logics.AddDownExcelHttpHeader(c, fmt.Sprintf("bk_cmdb_inst_%s.xlsx", objID))
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
	kit := rest.NewKitFromHeader(c.Request.Header, s.engine.CCErr)
	input := &operator.ExportInstParam{}
	if err := c.BindJSON(input); err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommJSONUnmarshalFailed, ErrMsg: err.Error()})
		return
	}

	objID := c.Param(common.BKObjIDField)
	input.ObjID = objID

	// 1. 初始化导出excel对象
	dir := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	filePath := fmt.Sprintf("%s/%s", dir, fmt.Sprintf("%dinst.xlsx", time.Now().UnixNano()))

	tmplOp, err := operator.NewTmplOp(operator.FilePath(filePath), operator.Dao(s.dao), operator.ObjID(objID),
		operator.Kit(kit), operator.Language(s.engine.Language))
	if err != nil {
		blog.Errorf("create excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}

	exporter, err := operator.NewExporter(operator.SetTmplOp(tmplOp), operator.SetExportParam(input))

	// 2. 导出实例数据到excel
	if err := exporter.Export(); err != nil {
		blog.Errorf("export instance data failed, err: %v, rid: %s", err, kit.Rid)
		c.String(http.StatusInternalServerError, fmt.Errorf("export instance data failed, err: %+v", err).Error())
		return
	}

	if err := exporter.Close(); err != nil {
		blog.Errorf("close excel io failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}

	// 3. 将excel模版文件返回，并删除临时文件
	logics.AddDownExcelHttpHeader(c, fmt.Sprintf("bk_cmdb_export_inst_%s.xlsx", objID))
	c.File(filePath)

	if err := exporter.Clean(); err != nil {
		blog.Errorf("clean excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}
}

// ExportHost export host
func (s *service) ExportHost(c *gin.Context) {
	kit := rest.NewKitFromHeader(c.Request.Header, s.engine.CCErr)
	defLang := s.engine.Language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header))

	input := &operator.ExportHostParam{}
	if err := c.BindJSON(input); err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommJSONUnmarshalFailed, ErrMsg: err.Error()})
		return
	}

	if input.ExportCond.Page.Limit <= 0 || input.ExportCond.Page.Limit > common.BKMaxOnceExportLimit {
		blog.Errorf("host page input is illegal, page: %v, rid: %s", input.ExportCond.Page, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrWebGetHostFail,
			ErrMsg: kit.CCError.Errorf(common.CCErrWebGetHostFail,
				defLang.Languagef("export_page_limit_err", common.BKMaxOnceExportLimit)).Error()})
		return
	}

	// 1. 初始化导出excel对象
	dir := fmt.Sprintf("%s/export", webCommon.ResourcePath)
	filePath := fmt.Sprintf("%s/%s", dir, fmt.Sprintf("%dinst.xlsx", time.Now().UnixNano()))
	objID := common.BKInnerObjIDHost
	tmplOp, err := operator.NewTmplOp(operator.FilePath(filePath), operator.Dao(s.dao), operator.ObjID(objID),
		operator.Kit(kit), operator.Language(s.engine.Language))
	if err != nil {
		blog.Errorf("create excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}

	exporter, err := operator.NewExporter(operator.SetTmplOp(tmplOp), operator.SetExportParam(input))

	// 2. 导出实例数据到excel
	if err := exporter.Export(); err != nil {
		blog.Errorf("export instance data failed, err: %v, rid: %s", err, kit.Rid)
		c.String(http.StatusInternalServerError, fmt.Errorf("export instance data failed, err: %+v", err).Error())
		return
	}

	if err := exporter.Close(); err != nil {
		blog.Errorf("close excel io failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}

	// 3. 将excel模版文件返回，并删除临时文件
	logics.AddDownExcelHttpHeader(c, "bk_cmdb_export_host.xlsx")
	c.File(filePath)

	if err := exporter.Clean(); err != nil {
		blog.Errorf("clean excel template failed, err: %v, rid: %s", err, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommExcelTemplateFailed, ErrMsg: err.Error()})
		return
	}
}
