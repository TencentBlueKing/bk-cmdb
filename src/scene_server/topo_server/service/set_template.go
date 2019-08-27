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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (s *Service) CreateSetTemplate(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (output interface{}, retErr error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.Atoi(bizIDStr)
	if err != nil {
		return nil, params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}
	setTemplate := metadata.SetTemplate{}
	if err := data.MarshalJSONInto(&setTemplate); err != nil {
		return nil, params.Err.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	setTemplate.BizID = int64(bizID)
	_ = setTemplate

	return setTemplate, nil
}
