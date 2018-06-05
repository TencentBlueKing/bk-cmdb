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
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initBusiness)
}

func (cli *topoAPI) initBusiness() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/app/{owner_id}", HandlerFunc: cli.CreateBusiness})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/app/{owner_id}/{app_id}", HandlerFunc: cli.DeleteBusiness})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/app/{owner_id}/{app_id}", HandlerFunc: cli.UpdateBusiness})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/app/status/{flag}/{owner_id}/{app_id}", HandlerFunc: cli.UpdateBusinessStatus})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/app/search/{owner_id}", HandlerFunc: cli.SearchBusiness})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/app/default/{owner_id}/search", HandlerFunc: cli.SearchDefaultBusiness})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/app/default/{owner_id}", HandlerFunc: cli.CreateDefaultBusiness})
}

// CreateBusiness create a new business
func (cli *topoAPI) CreateBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	return nil, nil
}

// DeleteBusiness delete the business
func (cli *topoAPI) DeleteBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// UpdateBusiness update the business
func (cli *topoAPI) UpdateBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// UpdateBusinessStatus update the business status
func (cli *topoAPI) UpdateBusinessStatus(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchBusiness search the business by condition
func (cli *topoAPI) SearchBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchDefaultBusiness search the business by condition
func (cli *topoAPI) SearchDefaultBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// CreateDefaultBusiness create the default business
func (cli *topoAPI) CreateDefaultBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}
