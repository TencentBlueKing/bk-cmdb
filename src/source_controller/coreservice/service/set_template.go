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
		blog.Errorf("GetSetTemplate failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, setTemplateID, err, params.ReqID)
		return nil, err
	}
	return setTemplate, nil
}
