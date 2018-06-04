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
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initObjectClassification)
}

func (cli *topoAPI) initObjectClassification() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/object/classification", HandlerFunc: cli.CreateClassification})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/object/classification/{owner_id}/objects", HandlerFunc: cli.SearchClassificationWithObjects})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/object/classifications", HandlerFunc: cli.SearchClassification})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/object/classification/{id}", HandlerFunc: cli.UpdateClassification})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/object/classification/{id}", HandlerFunc: cli.DeleteClassification})
}

// CreateClassification create a new object classification
func (cli *topoAPI) CreateClassification(params types.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchClassificationWithObjects search the classification with objects
func (cli *topoAPI) SearchClassificationWithObjects(params types.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchClassification search the classifications
func (cli *topoAPI) SearchClassification(params types.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// UpdateClassification update the object classification
func (cli *topoAPI) UpdateClassification(params types.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// DeleteClassification delete the object classification
func (cli *topoAPI) DeleteClassification(params types.LogicParams, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}
