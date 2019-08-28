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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (s *Service) CreateSetTemplate(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (output interface{}, retErr error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.CreateSetTemplateOption{}
	if err := data.MarshalJSONInto(&option); err != nil {
		return nil, params.Err.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	setTemplate, err := s.Engine.CoreAPI.CoreService().SetTemplate().CreateSetTemplate(params.Context, params.Header, bizID, option)
	if err != nil {
		blog.Errorf("CreateSetTemplate failed, core service create failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return setTemplate, nil
}

func (s *Service) UpdateSetTemplate(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (output interface{}, retErr error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	setTemplateIDStr := pathParams(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		return nil, params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
	}

	option := metadata.UpdateSetTemplateOption{}
	if err := data.MarshalJSONInto(&option); err != nil {
		return nil, params.Err.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	setTemplate, err := s.Engine.CoreAPI.CoreService().SetTemplate().UpdateSetTemplate(params.Context, params.Header, bizID, setTemplateID, option)
	if err != nil {
		blog.Errorf("UpdateSetTemplate failed, do core service update failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return setTemplate, nil
}

func (s *Service) DeleteSetTemplate(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (output interface{}, retErr error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.DeleteSetTemplateOption{}
	if err := data.MarshalJSONInto(&option); err != nil {
		return nil, params.Err.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	if err := s.Engine.CoreAPI.CoreService().SetTemplate().DeleteSetTemplate(params.Context, params.Header, bizID, option); err != nil {
		blog.Errorf("DeleteSetTemplate failed, do core service update failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return nil, nil
}

func (s *Service) GetSetTemplate(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (output interface{}, retErr error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	setTemplateIDStr := pathParams(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		return nil, params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
	}

	setTemplate, err := s.Engine.CoreAPI.CoreService().SetTemplate().GetSetTemplate(params.Context, params.Header, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("GetSetTemplate failed, do core service get failed, bizID: %d, setTemplateID: %d, err: %+v, rid: %s", bizID, setTemplateID, err, params.ReqID)
		return nil, err
	}
	return setTemplate, nil
}

func (s *Service) ListSetTemplate(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (output interface{}, retErr error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.ListSetTemplateOption{}
	if err := data.MarshalJSONInto(&option); err != nil {
		return nil, params.Err.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	setTemplate, err := s.Engine.CoreAPI.CoreService().SetTemplate().ListSetTemplate(params.Context, params.Header, bizID, option)
	if err != nil {
		blog.Errorf("ListSetTemplate failed, do core service list failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return setTemplate, nil
}

func (s *Service) ListSetTemplateWeb(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (output interface{}, retErr error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	setTemplates, err := s.ListSetTemplate(params, pathParams, queryParams, data)
	if err != nil {
		return nil, err
	}
	listResult := setTemplates.(*metadata.MultipleSetTemplateResult)
	if listResult == nil {
		return nil, nil
	}

	// count template instances
	setTemplateIDs := make([]int64, 0)
	for _, item := range listResult.Info {
		setTemplateIDs = append(setTemplateIDs, item.ID)
	}
	option := metadata.CountSetTplInstOption{
		SetTemplateIDs: setTemplateIDs,
	}
	setTplInstCount, err := s.Engine.CoreAPI.CoreService().SetTemplate().CountSetTplInstances(params.Context, params.Header, bizID, option)
	if err != nil {
		blog.Errorf("ListSetTemplateWeb failed, CountSetTplInstances failed, bizID: %d, option: %+v, err: %s, rid: %s", bizID, option, err.Error(), params.ReqID)
		return nil, err
	}
	result := metadata.MultipleSetTemplateWithStatisticsResult{
		Count: listResult.Count,
	}
	for _, setTemplate := range listResult.Info {
		setInstanceCount, exist := setTplInstCount[setTemplate.ID]
		if exist == false {
			setInstanceCount = 0
		}
		result.Info = append(result.Info, metadata.SetTemplateWithStatistics{
			SetInstanceCount: setInstanceCount,
			SetTemplate:      setTemplate,
		})
	}
	return result, nil
}
