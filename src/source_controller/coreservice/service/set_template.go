/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common/mapstruct"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) CreateSetTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.CreateSetTemplateOption{}
	if err := mapstr.DecodeFromMapStr(&option, data); err != nil {
		blog.Errorf("CreateSetTemplate failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.SetTemplateOperation().CreateSetTemplate(params, bizID, option)
	if err != nil {
		blog.Errorf("CreateSetTemplate failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) UpdateSetTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	setTemplateIDStr := pathParams(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
	}

	option := metadata.UpdateSetTemplateOption{}
	if err := mapstr.DecodeFromMapStr(&option, data); err != nil {
		blog.Errorf("UpdateSetTemplate failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.SetTemplateOperation().UpdateSetTemplate(params, setTemplateID, option)
	if err != nil {
		blog.Errorf("UpdateSetTemplate failed, setTemplateID: %d, option: %+v, err: %+v, rid: %s", setTemplateID, option, err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) DeleteSetTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.DeleteSetTemplateOption{}
	if err := mapstr.DecodeFromMapStr(&option, data); err != nil {
		blog.Errorf("DeleteSetTemplate failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if err := s.core.SetTemplateOperation().DeleteSetTemplate(params, bizID, option); err != nil {
		blog.Errorf("UpdateSetTemplate failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return nil, nil
}

func (s *coreService) GetSetTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	setTemplateIDStr := pathParams(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
	}

	setTemplate, err := s.core.SetTemplateOperation().GetSetTemplate(params, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("GetSetTemplate failed, bizID: %d, setTemplateID: %d, err: %+v, rid: %s", bizID, setTemplateID, err, params.ReqID)
		return nil, err
	}
	return setTemplate, nil
}

func (s *coreService) ListSetTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.ListSetTemplateOption{}
	if err := mapstr.DecodeFromMapStr(&option, data); err != nil {
		blog.Errorf("ListSetTemplate failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	setTemplateResult, err := s.core.SetTemplateOperation().ListSetTemplate(params, bizID, option)
	if err != nil {
		blog.Errorf("ListSetTemplate failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return setTemplateResult, nil
}

func (s *coreService) CountSetTplInstances(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.CountSetTplInstOption{}
	if err := mapstr.DecodeFromMapStr(&option, data); err != nil {
		blog.Errorf("CountSetTplInstances failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	filter := map[string]interface{}{
		common.BKSetTemplateIDField: map[string]interface{}{
			common.BKDBIN: option.SetTemplateIDs,
		},
		common.BKAppIDField: bizID,
	}
	pipeline := []map[string]interface{}{
		{common.BKDBMatch: filter},
		{common.BKDBGroup: map[string]interface{}{
			"_id":                 "$" + common.BKSetTemplateIDField,
			"set_instances_count": map[string]interface{}{common.BKDBSum: 1}},
		},
	}
	result := make([]metadata.CountSetTplInstItem, 0)
	if err := s.db.Table(common.BKTableNameBaseSet).AggregateAll(params.Context, pipeline, &result); err != nil {
		if s.db.IsNotFoundError(err) == true {
			result = make([]metadata.CountSetTplInstItem, 0)
		} else {
			blog.Errorf("CountSetTplInstances failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
			return result, params.Error.Error(common.CCErrCommDBSelectFailed)
		}
	}

	return result, nil
}

func (s *coreService) ListSetServiceTemplateRelations(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	setTemplateIDStr := pathParams(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
	}

	relations, err := s.core.SetTemplateOperation().ListSetServiceTemplateRelations(params, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("ListSetServiceTemplateRelations failed, bizID: %d, setTemplateID: %+v, err: %+v, rid: %s", bizID, setTemplateID, err, params.ReqID)
		return nil, err
	}
	return relations, nil
}

func (s *coreService) ListSetTplRelatedSvcTpl(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	setTemplateIDStr := pathParams(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
	}

	serviceTemplates, err := s.core.SetTemplateOperation().ListSetTplRelatedSvcTpl(params, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("ListSetTplRelatedSvcTpl failed, bizID: %d, setTemplateID: %d, err: %s, rid: %s", bizID, setTemplateID, err.Error(), params.ReqID)
		return nil, err
	}
	return serviceTemplates, nil
}

func (s *coreService) UpdateSetTemplateSyncStatus(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	setIDStr := pathParams(common.BKSetIDField)
	setID, err := strconv.ParseInt(setIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetIDField)
	}

	option := metadata.SetTemplateSyncStatus{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("UpdateSetTemplateSyncStatus failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if err := s.core.SetTemplateOperation().UpdateSetTemplateSyncStatus(params, setID, option); err != nil {
		blog.Errorf("UpdateSetTemplateSyncStatus failed, setID: %d, option: %+v, err: %+v, rid: %s", setID, option, err, params.ReqID)
		return nil, err
	}
	return nil, nil
}

func (s *coreService) ListSetTemplateSyncStatus(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.ListSetTemplateSyncStatusOption{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("ListSetTemplateSyncStatus failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}
	option.BizID = bizID

	result, err := s.core.SetTemplateOperation().ListSetTemplateSyncStatus(params, option)
	if err != nil {
		blog.Errorf("ListSetTemplateSyncStatus failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) ListSetTemplateSyncHistory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.ListSetTemplateSyncStatusOption{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("ListSetTemplateSyncHistory failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}
	option.BizID = bizID

	result, err := s.core.SetTemplateOperation().ListSetTemplateSyncHistory(params, option)
	if err != nil {
		blog.Errorf("ListSetTemplateSyncHistory failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) DeleteSetTemplateSyncStatus(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.DeleteSetTemplateSyncStatusOption{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("DeleteSetTemplateSyncStatus failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}
	option.BizID = bizID

	ccErr := s.core.SetTemplateOperation().DeleteSetTemplateSyncStatus(params, option)
	if ccErr != nil {
		blog.Errorf("DeleteSetTemplateSyncStatus failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, ccErr, params.ReqID)
		return nil, err
	}
	return nil, nil
}
