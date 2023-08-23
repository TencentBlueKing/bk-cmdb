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

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/gin-gonic/gin"
)

func (s *Service) initFieldTemplate(ws *gin.Engine) {
	ws.POST("/findmany/field_template", s.ListFieldTemplate)
	ws.POST("/findmany/field_template/object/count", s.CountFieldTemplateObj)
}

// ListFieldTemplate list field template with object condition ** ONLY FOR UI **
func (s *Service) ListFieldTemplate(c *gin.Context) {
	kit := rest.NewKitFromHeader(c.Request.Header, s.CCErr)

	opt := new(metadata.ListFieldTmplWithObjOption)
	if err := json.NewDecoder(c.Request.Body).Decode(opt); err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommHTTPReadBodyFailed, ErrMsg: err.Error()})
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: rawErr.ErrCode, ErrMsg: rawErr.ToCCError(kit.CCError).Error()})
		return
	}

	tmplFilter := opt.TemplateFilter
	if opt.ObjectFilter != nil {
		combinedFilter, isNone, err := s.parseFieldTmplWithObjFilter(kit, tmplFilter, opt.ObjectFilter)
		if err != nil {
			c.JSON(http.StatusOK, metadata.BaseResp{Code: err.GetCode(), ErrMsg: err.Error()})
			return
		}

		if isNone {
			c.JSON(http.StatusOK, metadata.NewSuccessResp(metadata.FieldTemplateInfo{Count: 0,
				Info: make([]metadata.FieldTemplate, 0)}))
			return
		}

		tmplFilter = combinedFilter
	}

	// list field template
	tmplOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: tmplFilter},
		Page:               opt.Page,
		Fields:             opt.Fields,
	}
	tmplRes, err := s.CoreAPI.ApiServer().FieldTemplate().ListFieldTemplate(kit.Ctx, kit.Header, tmplOpt)
	if err != nil {
		blog.Errorf("list field template failed, err: %v, opt: %+v, rid: %s", err, tmplOpt, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: err.GetCode(), ErrMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, metadata.NewSuccessResp(tmplRes))
}

// parseFieldTmplWithObjFilter parse field template with object filter to field template query filter
func (s *Service) parseFieldTmplWithObjFilter(kit *rest.Kit, tmplFilter, objFilter *filter.Expression) (
	*filter.Expression, bool, errors.CCErrorCoder) {

	if objFilter == nil {
		return tmplFilter, false, nil
	}

	// get object ids by object filter
	objCond, err := objFilter.ToMgo()
	if err != nil {
		blog.Errorf("parse object filter failed, err: %v, filter: %v, rid: %s", err, objFilter, kit.Rid)
		return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "object_filter")
	}

	objOpt := &metadata.QueryCondition{
		Fields:    []string{common.BKFieldID},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Condition: objCond,
	}
	objRes, ccErr := s.CoreAPI.ApiServer().ReadModel(kit.Ctx, kit.Header, objOpt)
	if ccErr != nil {
		blog.Errorf("list object ids failed, err: %v, opt: %+v, rid: %s", ccErr, objOpt, kit.Rid)
		return nil, false, ccErr
	}

	if len(objRes.Info) == 0 {
		return nil, true, nil
	}

	objIDs := make([]int64, len(objRes.Info))
	for i, object := range objRes.Info {
		objIDs[i] = object.ID
	}

	// get object related template ids by relation
	relOpt := &metadata.ListObjFieldTmplRelOption{
		ObjectIDs: objIDs,
	}
	relRes, ccErr := s.CoreAPI.ApiServer().FieldTemplate().ListObjFieldTmplRel(kit.Ctx, kit.Header, relOpt)
	if ccErr != nil {
		blog.ErrorJSON("list obj field template relation failed, err: %s, opt: %s, rid: %s", ccErr, relOpt, kit.Rid)
		return nil, false, ccErr
	}

	if len(relRes.Info) == 0 {
		return nil, true, nil
	}

	tmplIDs := make([]int64, len(relRes.Info))
	for i, relation := range relRes.Info {
		tmplIDs[i] = relation.TemplateID
	}
	tmplIDs = util.IntArrayUnique(tmplIDs)

	tmplFilter, err = filtertools.And(tmplFilter, filtertools.GenAtomFilter(common.BKFieldID, filter.In, tmplIDs))
	if err != nil {
		blog.Errorf("add template id filter failed, err: %v, filter: %v, rid: %s", err, tmplFilter, kit.Rid)
		return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "template_filter")
	}

	return tmplFilter, false, nil
}

// CountFieldTemplateObj count field templates related objects ** ONLY FOR UI **
func (s *Service) CountFieldTemplateObj(c *gin.Context) {
	kit := rest.NewKitFromHeader(c.Request.Header, s.CCErr)

	opt := new(metadata.CountFieldTmplResOption)
	if err := json.NewDecoder(c.Request.Body).Decode(opt); err != nil {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: common.CCErrCommHTTPReadBodyFailed, ErrMsg: err.Error()})
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		c.JSON(http.StatusOK, metadata.BaseResp{Code: rawErr.ErrCode, ErrMsg: rawErr.ToCCError(kit.CCError).Error()})
		return
	}

	// get field template and object relations
	objOpt := &metadata.ListObjFieldTmplRelOption{
		TemplateIDs: opt.TemplateIDs,
	}
	relRes, err := s.CoreAPI.ApiServer().FieldTemplate().ListObjFieldTmplRel(kit.Ctx, kit.Header, objOpt)
	if err != nil {
		blog.Errorf("list field template and object relation failed, err: %v, opt: %+v, rid: %s", err, opt, kit.Rid)
		c.JSON(http.StatusOK, metadata.BaseResp{Code: err.GetCode(), ErrMsg: err.Error()})
		return
	}

	countMap := make(map[int64]int, 0)
	for _, relation := range relRes.Info {
		countMap[relation.TemplateID]++
	}

	countInfos := make([]metadata.FieldTmplResCount, len(opt.TemplateIDs))
	for i, templateID := range opt.TemplateIDs {
		countInfos[i] = metadata.FieldTmplResCount{
			TemplateID: templateID,
			Count:      countMap[templateID],
		}
	}

	c.JSON(http.StatusOK, metadata.NewSuccessResp(countInfos))
}
