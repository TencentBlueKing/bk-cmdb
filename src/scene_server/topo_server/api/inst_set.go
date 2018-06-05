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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	frcommon "configcenter/src/framework/common"
	frtypes "configcenter/src/framework/core/types"

	"configcenter/src/scene_server/topo_server/core/types"
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
func (cli *topoAPI) CreateSet(params types.LogicParams, parthParams, queryParams func(name string) string, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).Field(common.BKObjIDField).Eq(common.BKInnerObjIDSet)

	objItems, err := cli.core.FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return nil, err
	}

	for _, item := range objItems {
		setInst, err := cli.core.CreateInst(params, item, data)
		if nil != err {
			blog.Errorf("failed to create a new set, %s", err.Error())
			return nil, err
		}

		err = setInst.Save()
		if nil != err {
			blog.Errorf("failed to create a new set, %s", err.Error())
			return nil, err
		}
	}

	return nil, nil
}

// DeleteSet delete the set
func (cli *topoAPI) DeleteSet(params types.LogicParams, parthParams, queryParams func(name string) string, data frtypes.MapStr) (frtypes.MapStr, error) {

	//cli.core.DeleteInst

	return nil, nil
}

// UpdateSet update the set
func (cli *topoAPI) UpdateSet(params types.LogicParams, parthParams, queryParams func(name string) string, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}

// SearchSet search the set
func (cli *topoAPI) SearchSet(params types.LogicParams, parthParams, queryParams func(name string) string, data frtypes.MapStr) (frtypes.MapStr, error) {
	return nil, nil
}
