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
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initAssociation)
}

func (cli *topoAPI) initAssociation() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/model/mainline", HandlerFunc: cli.CreateMainLineObject})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/model/mainline/owners/{owner_id}/objectids/{obj_id}", HandlerFunc: cli.DeleteMainLineObject})
	cli.actions = append(cli.actions, action{Method: http.MethodGet, Path: "/model/{owner_id}", HandlerFunc: cli.SearchMainLineOBjectTopo})
	cli.actions = append(cli.actions, action{Method: http.MethodGet, Path: "/model/{owner_id}/{cls_id}/{obj_id}", HandlerFunc: cli.SearchObjectByClassificationID})
	cli.actions = append(cli.actions, action{Method: http.MethodGet, Path: "/inst/{owner_id}/{app_id}", HandlerFunc: cli.SearchBusinessTopo})
	cli.actions = append(cli.actions, action{Method: http.MethodGet, Path: "/inst/child/{owner_id}/{obj_id}/{app_id}/{inst_id}", HandlerFunc: cli.SearchMainLineChildInstTopo})
}

// CreateMainLineObject create a new object in the main line topo
func (cli *topoAPI) CreateMainLineObject(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	return nil, nil
}

// DeleteMainLineObject delete a object int the main line topo
func (cli *topoAPI) DeleteMainLineObject(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchMainLineOBjectTopo search the main line topo
func (cli *topoAPI) SearchMainLineOBjectTopo(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchObjectByClassificationID search the object by classification ID
func (cli *topoAPI) SearchObjectByClassificationID(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchBusinessTopo search the business topo
func (cli *topoAPI) SearchBusinessTopo(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchMainLineChildInstTopo search the child inst topo by a inst
func (cli *topoAPI) SearchMainLineChildInstTopo(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}
