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

package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/gin-gonic/gin"
)

func (s *Service) initResourceCount(ws *gin.Engine) {
	ws.POST("/count/service_instance/processes", s.CountProcessesBySvcInst)
	ws.POST("/count/set_template/:set_template_id/service_template/hosts", s.CountHostsBySvcTmpl)
}

// CountProcessesBySvcInst count processes by service instances ** ONLY FOR UI **
func (s *Service) CountProcessesBySvcInst(c *gin.Context) {
	kit := rest.NewKitFromHeader(c.Request.Header, s.CCErr)

	opt := new(metadata.CountByIDsOption)
	if err := json.NewDecoder(c.Request.Body).Decode(opt); err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommHTTPReadBodyFailed, ErrMsg: err.Error()})
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: rawErr.ErrCode, ErrMsg: rawErr.ToCCError(kit.CCError).Error()})
		return
	}

	cntOpt := &metadata.GroupRelResByIDsOption{
		IDs:      opt.IDs,
		IDField:  common.BKServiceInstanceIDField,
		RelField: common.BKProcessIDField,
	}
	idMap, err := s.ApiCli.GroupRelResByIDs(kit.Ctx, kit.Header, metadata.ProcInstRelGroupByRes, cntOpt)
	if err != nil {
		blog.Errorf("count processes by service instances failed, err: %v, opt: %+v, rid: %s", err, cntOpt, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: err.GetCode(), ErrMsg: err.Error()})
		return
	}

	result := make([]metadata.IDCountInfo, len(opt.IDs))
	for i, id := range opt.IDs {
		result[i] = metadata.IDCountInfo{
			ID:    id,
			Count: len(idMap[id]),
		}
	}

	c.JSON(http.StatusOK, metadata.NewSuccessResp(result))
}

// CountHostsBySvcTmpl count hosts by service templates ** ONLY FOR UI **
func (s *Service) CountHostsBySvcTmpl(c *gin.Context) {
	kit := rest.NewKitFromHeader(c.Request.Header, s.CCErr)

	opt := new(metadata.CountByIDsOption)
	if err := json.NewDecoder(c.Request.Body).Decode(opt); err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommHTTPReadBodyFailed, ErrMsg: err.Error()})
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: rawErr.ErrCode, ErrMsg: rawErr.ToCCError(kit.CCError).Error()})
		return
	}

	setTmplIDStr := c.Param(common.BKSetTemplateIDField)
	setTmplID, rawErr := strconv.ParseInt(setTmplIDStr, 10, 64)
	if rawErr != nil {
		blog.Errorf("convert set template id(%s) to int64 failed, err: %v, rid: %s", setTmplIDStr, rawErr, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommParamsIsInvalid, ErrMsg: rawErr.Error()})
		return
	}

	result, err := s.countHostsBySvcTmpl(kit, setTmplID, opt)
	if err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: err.GetCode(), ErrMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, metadata.NewSuccessResp(result))
}

func (s *Service) countHostsBySvcTmpl(kit *rest.Kit, setTmplID int64, opt *metadata.CountByIDsOption) (
	[]metadata.IDCountInfo, errors.CCErrorCoder) {

	moduleOpt := &metadata.GroupRelResByIDsOption{
		IDs:       opt.IDs,
		IDField:   common.BKServiceTemplateIDField,
		RelField:  common.BKModuleIDField,
		ExtraCond: mapstr.MapStr{common.BKSetTemplateIDField: setTmplID},
	}
	moduleMap, err := s.ApiCli.GroupRelResByIDs(kit.Ctx, kit.Header, metadata.ModuleGroupByRes, moduleOpt)
	if err != nil {
		blog.Errorf("group modules by service templates failed, err: %v, opt: %+v, rid: %s", err, moduleOpt, kit.Rid)
		return nil, err
	}

	if len(moduleMap) == 0 {
		result := make([]metadata.IDCountInfo, len(opt.IDs))
		for i, id := range opt.IDs {
			result[i] = metadata.IDCountInfo{ID: id, Count: 0}
		}
		return result, nil
	}

	allModuleIDs := make([]int64, 0)
	moduleTmplMap := make(map[int64]int64)
	for svcTemp, moduleIDVal := range moduleMap {
		moduleIDs, err := util.SliceInterfaceToInt64(moduleIDVal)
		if err != nil {
			blog.Errorf("convert module ids(%+v) to int64 failed, err: %v, rid: %s", moduleIDVal, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKModuleIDField)
		}
		allModuleIDs = append(allModuleIDs, moduleIDs...)
		for _, moduleID := range moduleIDs {
			moduleTmplMap[moduleID] = svcTemp
		}
	}

	relOpt := &metadata.GroupRelResByIDsOption{
		IDField:  common.BKModuleIDField,
		RelField: common.BKHostIDField,
	}
	total := len(allModuleIDs)
	allRelMap := make(map[int64][]interface{})
	for start := 0; start < total; start += common.BKMaxUpdateOrCreatePageSize {
		if total-start >= common.BKMaxUpdateOrCreatePageSize {
			relOpt.IDs = allModuleIDs[start : start+common.BKMaxUpdateOrCreatePageSize]
		} else {
			relOpt.IDs = allModuleIDs[start:total]
		}

		relMap, err := s.ApiCli.GroupRelResByIDs(kit.Ctx, kit.Header, metadata.ModuleHostRelGroupByRes, relOpt)
		if err != nil {
			blog.Errorf("group host relations by modules failed, err: %v, opt: %+v, rid: %s", err, relOpt, kit.Rid)
			return nil, err
		}
		for moduleID, hostIDVal := range relMap {
			allRelMap[moduleID] = append(allRelMap[moduleID], hostIDVal...)
		}
	}

	tmplHostMap := make(map[int64][]int64)
	for moduleID, hostIDVal := range allRelMap {
		hostIDs, err := util.SliceInterfaceToInt64(hostIDVal)
		if err != nil {
			blog.Errorf("convert host ids(%+v) to int64 failed, err: %v, rid: %s", hostIDVal, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
		}

		tmplID := moduleTmplMap[moduleID]
		tmplHostMap[tmplID] = append(tmplHostMap[tmplID], hostIDs...)
	}

	result := make([]metadata.IDCountInfo, len(opt.IDs))
	for i, id := range opt.IDs {
		hostIDs := util.IntArrayUnique(tmplHostMap[id])
		result[i] = metadata.IDCountInfo{
			ID:    id,
			Count: len(hostIDs),
		}
	}
	return result, nil
}
