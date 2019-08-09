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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) ListHosts(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	option := meta.ListHosts{}
	if err := data.MarshalJSONInto(&option); err != nil {
		blog.Errorf("get host module id failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	hosts, err := s.core.HostOperation().ListHosts(params, option)
	if err != nil {
		blog.Errorf("ListHostByTopoNode failed, call host operation failed, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	return hosts, nil
}
