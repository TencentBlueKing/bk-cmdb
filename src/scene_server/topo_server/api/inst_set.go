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
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initSet)
}

func (cli *topoAPI) initSet() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/set/{app_id}", HandlerFunc: cli.CreateSet})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/set/{app_id}/{set_id}", HandlerFunc: cli.DeleteSet})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/set/{app_id}/{set_id}", HandlerFunc: cli.UpdateSet})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/set/search/{owner_id}/{app_id}", HandlerFunc: cli.SearchSet})

}

// CreateSet create a new set
func (cli *topoAPI) CreateSet(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {

	return nil, nil
}

// DeleteSet delete the set
func (cli *topoAPI) DeleteSet(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {

	return nil, nil
}

// UpdateSet update the set
func (cli *topoAPI) UpdateSet(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchSet search the set
func (cli *topoAPI) SearchSet(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}
