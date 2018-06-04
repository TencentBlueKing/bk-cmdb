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
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initInst)
}


func (cli *topoAPI) initInst() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/{owner_id}/{obj_id}", HandlerFunc: cli.CreateInst})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/inst/{owner_id}/{obj_id}/{inst_id}", HandlerFunc: cli.DeleteInst})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/inst/{owner_id}/{obj_id}/{inst_id}", HandlerFunc: cli.UpdateInst})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/search/{owner_id}/{obj_id}", HandlerFunc: cli.SearchInst})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/search/owner/{owner_id}/object/{obj_id}/detail", HandlerFunc: cli.SearchInstAndAssociationDetail})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/search/owner/{owner_id}/object/{obj_id}", HandlerFunc: cli.SearchInstByObject})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/association/search/owner/{owner_id}/object/{obj_id}", HandlerFunc: cli.SearchInstByAssociation})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/search/{owner_id}/{obj_id}/{inst_id}", HandlerFunc: cli.CreateSet})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/search/topo/owner/{owner_id}/object/{object_id}/inst/{inst_id}", HandlerFunc: cli.SearchInstChildTopo})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/association/topo/search/owner/{owner_id}/object/{object_id}/inst/{inst_id}", HandlerFunc: cli.SearchInstTopo})
}

// CreateInst create a new inst
func (cli *topoAPI) CreateInst(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {

	return nil, nil
}

// DeleteInst delete the inst
func (cli *topoAPI) DeleteInst(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// UpdateInst update the inst
func (cli *topoAPI) UpdateInst(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchInst search the inst
func (cli *topoAPI) SearchInst(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchInstAndAssociationDetail search the inst with association details
func (cli *topoAPI) SearchInstAndAssociationDetail(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchInstByObject search the inst of the object
func (cli *topoAPI) SearchInstByObject(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchInstByAssociation search inst by the association inst
func (cli *topoAPI) SearchInstByAssociation(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchInstByInstID search the inst by inst ID
func (cli *topoAPI) SearchInstByInstID(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchInstChildTopo search the child inst topo for a inst
func (cli *topoAPI) SearchInstChildTopo(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchInstTopo search the inst topo
func (cli *topoAPI) SearchInstTopo(params core.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}
