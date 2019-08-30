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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateClassification create a new object classification
func (s *Service) CreateClassification(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	cls, err := s.Core.ClassificationOperation().CreateClassification(params, data)
	if nil != err {
		return nil, err
	}

	return cls.ToMapStr()
}

// SearchClassificationWithObjects search the classification with objects
func (s *Service) SearchClassificationWithObjects(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	if data.Exists(metadata.PageName) {

		page, err := data.MapStr(metadata.PageName)
		if nil != err {
			blog.Errorf("failed to get the page , error info is %s, rid: %s", err.Error(), params.ReqID)
			return nil, err
		}

		if err = cond.SetPage(page); nil != err {
			blog.Errorf("failed to parse the page, error info is %s, rid: %s", err.Error(), params.ReqID)
			return nil, err
		}

		data.Remove(metadata.PageName)
	}

	if err := cond.Parse(data); nil != err {
		blog.Errorf("failed to parse the condition, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	return s.Core.ClassificationOperation().FindClassificationWithObjects(params, cond)
}

// SearchClassification search the classifications
func (s *Service) SearchClassification(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	if data.Exists(metadata.PageName) {

		page, err := data.MapStr(metadata.PageName)
		if nil != err {
			blog.Errorf("failed to get the page , error info is %s, rid: %s", err.Error(), params.ReqID)
			return nil, err
		}

		if err = cond.SetPage(page); nil != err {
			blog.Errorf("failed to parse the page, error info is %s, rid: %s", err.Error(), params.ReqID)
			return nil, err
		}

		data.Remove(metadata.PageName)
	}
	if err := cond.Parse(data); err != nil {
		blog.Errorf("parse condition from data failed, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	return s.Core.ClassificationOperation().FindClassification(params, cond)
}

// UpdateClassification update the object classification
func (s *Service) UpdateClassification(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	paramPath := mapstr.MapStr{}
	paramPath.Set("id", pathParams("id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-cls] failed to parse the path params id(%s), error info is %s , rid: %s", pathParams("id"), err.Error(), params.ReqID)
		return nil, err
	}
	data.Remove(metadata.BKMetadata)

	err = s.Core.ClassificationOperation().UpdateClassification(params, data, id, cond)
	if err != nil {
		return nil, err
	}

	return nil, err
}

// DeleteClassification delete the object classification
func (s *Service) DeleteClassification(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	paramPath := mapstr.MapStr{}
	paramPath.Set("id", pathParams("id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-cls] failed to parse the path params id(%s), error info is %s , rid: %s", pathParams("id"), err.Error(), params.ReqID)
		return nil, err
	}

	// data.Remove(metadata.BKMetadata)
	err = s.Core.ClassificationOperation().DeleteClassification(params, id, data, cond)
	return nil, err
}
