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

	frtypes "configcenter/src/common/types"
	"configcenter/src/scene_server/topo_server/core/types"
)

func init() {
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initObjectGroup)
}

func (cli *topoAPI) initObjectGroup() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/objectatt/group/new", HandlerFunc: cli.CreateObjectGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/objectatt/group/update", HandlerFunc: cli.UpdateObjectGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/objectatt/group/groupid/{id}", HandlerFunc: cli.DeleteObjectGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/objectatt/group/property", HandlerFunc: cli.UpdateObjectAttributeGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/objectatt/group/owner/{owner_id}/object/{object_id}/propertyids/{property_id}/groupids/{group_id}", HandlerFunc: cli.DeleteObjectAttributeGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/objectatt/group/property/owner/{owner_id}/object/{object_id}", HandlerFunc: cli.SearchGroupByObject})
}

// CreateObjectGroup create a new object group
func (cli *topoAPI) CreateObjectGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// UpdateObjectGroup update the object group information
func (cli *topoAPI) UpdateObjectGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// DeleteObjectGroup delete the object group
func (cli *topoAPI) DeleteObjectGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// UpdateObjectAttributeGroup update the object attribute belongs to group information
func (cli *topoAPI) UpdateObjectAttributeGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// DeleteObjectAttributeGroup delete the object attribute belongs to group information
func (cli *topoAPI) DeleteObjectAttributeGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchGroupByObject search the groups by the object
func (cli *topoAPI) SearchGroupByObject(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}
