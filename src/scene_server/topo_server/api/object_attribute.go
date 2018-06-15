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
	"fmt"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
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
func (cli *topoAPI) CreateObjectAttribute(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("CreateObjectAttribute")
	attr, err := cli.core.AttributeOperation().CreateObjectAttribute(params, data)
	if nil != err {
		return nil, err
	}

	return attr.ToMapStr()
}

// SearchObjectAttribute search the object attributes
func (cli *topoAPI) SearchObjectAttribute(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchObjectAttribute")
	cond := condition.CreateCondition()
	data.Remove(metadata.PageName)
	if err := cond.Parse(data); nil != err {
		blog.Errorf("failed to parset the data into condition, error info is %s", err.Error())
		return nil, err
	}

	attrs, err := cli.core.AttributeOperation().FindObjectAttribute(params, cond)
	if nil != err {
		blog.Errorf("failed to parse the data into condition, error info is %s", err.Error())
		return nil, err
	}

	result := frtypes.MapStr{}
	result.Set("data", attrs)
	return result, nil
}

// UpdateObjectAttribute update the object attribute
func (cli *topoAPI) UpdateObjectAttribute(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("UpdateObjectAttribute")
	cond := condition.CreateCondition()
	cond.Field("id").Eq(queryParams("id"))
	err := cli.core.AttributeOperation().UpdateObjectAttribute(params, data, cond)

	return nil, err
}

// DeleteObjectAttribute delete the object attribute
func (cli *topoAPI) DeleteObjectAttribute(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("DeleteObjectAttribute")
	cond := condition.CreateCondition()
	cond.Field("id").Eq(queryParams("id"))
	err := cli.core.AttributeOperation().DeleteObjectAttribute(params, cond)

	return nil, err
}
