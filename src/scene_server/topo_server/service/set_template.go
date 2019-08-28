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
