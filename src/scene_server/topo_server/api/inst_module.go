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
	"configcenter/src/scene_server/topo_server/core"
)

func init() {
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initModule)
}

func (cli *topoAPI) initModule() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/module/{app_id}/{set_id}", HandlerFunc: cli.CreateModule})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/module/{app_id}/{set_id}/{module_id}", HandlerFunc: cli.DeleteModule})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/module/{app_id}/{set_id}/{module_id}", HandlerFunc: cli.UpdateModule})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/module/search/{owner_id}/{app_id}/{set_id}", HandlerFunc: cli.SearchModule})

}

// CreateModule create a new module
func (cli *topoAPI) CreateModule(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {

	return nil, nil
}

// DeleteModule delete the module
func (cli *topoAPI) DeleteModule(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {

	return nil, nil
}

// UpdateModule update the module
func (cli *topoAPI) UpdateModule(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchModule search the modules
func (cli *topoAPI) SearchModule(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}
