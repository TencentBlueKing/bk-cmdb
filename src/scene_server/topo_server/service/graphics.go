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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
	"encoding/json"
)

func (s Service) ParseOriginGraphicsUpdateInput(data []byte) (mapstr.MapStr, error) {
	datas := make([]metadata.TopoGraphics, 0)
	err := json.Unmarshal(data, &datas)
	if nil != err {
		return nil, err
	}
	result := mapstr.New()
	result.Set("origin", datas)
	return result, nil
}
func (s *Service) SelectObjectTopoGraphics(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	return s.Core.GraphicsOperation().SelectObjectTopoGraphics(params, pathParams("scope_type"), pathParams("scope_id"))
}

func (s *Service) UpdateObjectTopoGraphics(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	datas := make([]metadata.TopoGraphics, 0)
	val, exists := data.Get("origin")
	if !exists {
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, "not set anything")
	}

	datas, _ = val.([]metadata.TopoGraphics)

	err := s.Core.GraphicsOperation().UpdateObjectTopoGraphics(params, pathParams("scope_type"), pathParams("scope_id"), datas)
	return nil, err
}

func (s *Service) UpdateObjectTopoGraphicsNew(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	datas := metadata.UpdateTopoGraphicsInput{}
	err := data.MarshalJSONInto(&datas)
	if nil != err {
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, "not set anything")
	}

	err = s.Core.GraphicsOperation().UpdateObjectTopoGraphics(params, pathParams("scope_type"), pathParams("scope_id"), datas.Origin)
	return nil, err
}
