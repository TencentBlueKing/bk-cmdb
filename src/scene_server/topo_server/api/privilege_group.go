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

	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/types"
)

func init() {
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initObject)
}

func (cli *topoAPI) initPrivilegeGroup() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/privilege/group/{bk_supplier_account}", HandlerFunc: cli.CreateUserGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/privilege/group/{bk_supplier_account}/{group_id}", HandlerFunc: cli.DeleteUserGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/privilege/group/{bk_supplier_account}/{group_id}", HandlerFunc: cli.UpdateUserGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/privilege/group/{bk_supplier_account}/search", HandlerFunc: cli.SearchUserGroup})
}

// CreateUserGroup create user goup
func (cli *topoAPI) CreateUserGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchObjectBatch")
	return nil, nil
}

// DeleteUserGroup delete user goup
func (cli *topoAPI) DeleteUserGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchObjectBatch")
	return nil, nil
}

// UpdateUserGroup update user goup
func (cli *topoAPI) UpdateUserGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchObjectBatch")
	return nil, nil
}

// SearchUserGroup search user goup
func (cli *topoAPI) SearchUserGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchObjectBatch")
	return nil, nil
}
