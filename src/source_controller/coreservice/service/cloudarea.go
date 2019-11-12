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
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) UpdateHostCloudAreaField(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := metadata.UpdateHostCloudAreaFieldOption{}
	if err := data.MarshalJSONInto(&input); nil != err {
		blog.Errorf("UpdateHostCloudAreaField failed, err:%s, input:%v, rid: %v", data, err.Error(), params.ReqID)
		return nil, err
	}

	err := s.core.HostOperation().UpdateHostCloudAreaField(params, input)
	if err != nil {
		blog.Errorf("UpdateHostCloudAreaField failed, call core operation failed, input: %+v, err: %v, rid: %v", input, err, params.ReqID)
		return nil, err
	}

	return nil, nil
}
