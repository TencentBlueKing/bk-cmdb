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

	frtypes "configcenter/src/framework/core/types"
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
func (cli *topoAPI) CreateObjectAttribute(params types.LogicParams, parthParams, queryParams func(name string) string, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchObjectAttribute search the object attributes
func (cli *topoAPI) SearchObjectAttribute(params types.LogicParams, parthParams, queryParams func(name string) string, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// UpdateObjectAttribute update the object attribute
func (cli *topoAPI) UpdateObjectAttribute(params types.LogicParams, parthParams, queryParams func(name string) string, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// DeleteObjectAttribute delete the object attribute
func (cli *topoAPI) DeleteObjectAttribute(params types.LogicParams, parthParams, queryParams func(name string) string, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}
