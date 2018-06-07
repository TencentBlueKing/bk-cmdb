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

package api

import (
	"net/http"

	// "configcenter/src/common"
	// "configcenter/src/common/blog"
	frcommon "configcenter/src/common"
	frtypes "configcenter/src/common/types"
	"configcenter/src/scene_server/topo_server/core/types"
)

func init() {
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initObjectAttribute)
}

func (cli *topoAPI) initObjectAttribute() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/objectatt", HandlerFunc: cli.CreateObjectAttribute})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/objectatt/search", HandlerFunc: cli.SearchObjectAttribute})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/objectatt/{id}", HandlerFunc: cli.UpdateObjectAttribute})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/objectatt/{id}", HandlerFunc: cli.DeleteObjectAttribute})
}

// CreateObjectAttribute create a new object attribute
func (cli *topoAPI) CreateObjectAttribute(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	attr, err := cli.core.CreateObjectAttribute(params, data)
	if nil != err {
		return nil, err
	}

	return attr.ToMapStr()
}

// SearchObjectAttribute search the object attributes
func (cli *topoAPI) SearchObjectAttribute(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()

	// TODO: data => cond

	attrs, err := cli.core.FindObjectAttribute(params, cond)

	if nil != err {
		return nil, err
	}

	result := frtypes.MapStr{}
	items := make([]frtypes.MapStr, 0)
	for _, item := range attrs {

		obj, err := item.ToMapStr()
		if nil != err {
			return nil, err
		}
		items = append(items, obj)
	}

	result.Set("data", items)

	return result, nil
}

// UpdateObjectAttribute update the object attribute
func (cli *topoAPI) UpdateObjectAttribute(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()

	cond.Field("id")

	err := cli.core.UpdateObjectAttribute(params, data, cond)

	return nil, err
}

// DeleteObjectAttribute delete the object attribute
func (cli *topoAPI) DeleteObjectAttribute(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()

	cond.Field("id")

	err := cli.core.DeleteObjectAttribute(params, cond)

	return nil, err
}
