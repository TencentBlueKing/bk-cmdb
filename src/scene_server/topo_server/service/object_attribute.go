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
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateObjectAttribute create a new object attribute
func (s *topoService) CreateObjectAttribute(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	attr, err := s.core.AttributeOperation().CreateObjectAttribute(params, data)
	if nil != err {
		return nil, err
	}

	return attr.ToMapStr()
}

// SearchObjectAttribute search the object attributes
func (s *topoService) SearchObjectAttribute(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	data.Remove(metadata.PageName)
	if err := cond.Parse(data); nil != err {
		blog.Errorf("failed to parset the data into condition, error info is %s", err.Error())
		return nil, err
	}
	cond.Field(metadata.AttributeFieldIsSystem).Eq(false)
	return s.core.AttributeOperation().FindObjectAttributeWithDetail(params, cond)
}

// UpdateObjectAttribute update the object attribute
func (s *topoService) UpdateObjectAttribute(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	paramPath := frtypes.MapStr{}
	paramPath.Set("id", pathParams("id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-att] failed to parse the path params id(%s), error info is %s ", pathParams("id"), err.Error())
		return nil, err
	}

	err = s.core.AttributeOperation().UpdateObjectAttribute(params, data, id, cond)

	return nil, err
}

// DeleteObjectAttribute delete the object attribute
func (s *topoService) DeleteObjectAttribute(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()

	paramPath := frtypes.MapStr{}
	paramPath.Set("id", pathParams("id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-att] failed to parse the path params id(%s), error info is %s ", pathParams("id"), err.Error())
		return nil, err
	}

	err = s.core.AttributeOperation().DeleteObjectAttribute(params, id, cond)

	return nil, err
}
