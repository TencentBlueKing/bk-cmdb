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

package extensions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// GetHostLayers get resource layers id by hostID(layers is a data structure for call iam)
func GetHostLayers(coreService coreservice.CoreServiceClientInterface, requestHeader *http.Header, hostIDArr *[]int64) (
	bkBizID int64, batchLayers [][]meta.Item, err error) {
	batchLayers = make([][]meta.Item, 0)

	cond := condition.CreateCondition()
	cond.Field(common.BKHostIDField).In(*hostIDArr)
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKAppIDField, common.BKModuleIDField, common.BKSetIDField, common.BKHostIDField},
		Condition: cond.ToMapStr(),
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	result, err := coreService.Instance().ReadInstance(
		context.Background(), *requestHeader, common.BKTableNameModuleHostConfig, query)
	if err != nil {
		err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
		return
	}
	blog.V(5).Infof("get host module config: %+v", result.Data.Info)
	if len(result.Data.Info) == 0 {
		err = fmt.Errorf("get host:%+v layer failed, get host module config by host id not found, maybe hostID invalid",
			hostIDArr)
		return
	}
	bkBizID, err = util.GetInt64ByInterface(result.Data.Info[0][common.BKAppIDField])
	if err != nil {
		err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
		return
	}

	bizTopoTreeRoot, err := coreService.Mainline().SearchMainlineInstanceTopo(
		context.Background(), *requestHeader, bkBizID, true)
	if err != nil {
		err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
		return
	}

	bizTopoTreeRootJSON, err := json.Marshal(bizTopoTreeRoot)
	if err != nil {
		err = fmt.Errorf("json encode bizTopoTreeRootJSON failed: %+v", err)
		return
	}
	blog.V(5).Infof("bizTopoTreeRoot: %s", bizTopoTreeRootJSON)

	dataInfo, err := json.Marshal(result.Data.Info)
	if err != nil {
		err = fmt.Errorf("json encode dataInfo failed: %+v", err)
		return
	}
	blog.V(5).Infof("dataInfo: %s", dataInfo)

	hostIDs := make([]int64, 0)
	for _, item := range result.Data.Info {
		hostID, e := util.GetInt64ByInterface(item[common.BKHostIDField])
		if e != nil {
			err = fmt.Errorf("extract hostID from host info failed, host: %+v, err: %+v", item, e)
			return
		}
		hostIDs = append(hostIDs, hostID)
	}

	hostIDInnerIPMap, err := getInnerIPByHostIDs(coreService, *requestHeader, &hostIDs)
	if err != nil {
		err = fmt.Errorf("get host:%+v InnerIP failed, err: %+v", hostIDs, err)
		return
	}

	for _, item := range result.Data.Info {
		bizID, err := util.GetInt64ByInterface(item[common.BKAppIDField])
		if err != nil {
			err = fmt.Errorf("get host:%+v layer failed, get bk_app_id field failed, err: %+v", item, err)
		}
		if bizID != bkBizID {
			continue
		}
		moduleID, err := util.GetInt64ByInterface(item[common.BKModuleIDField])
		if err != nil {
			err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
		}
		path := bizTopoTreeRoot.TraversalFindModule(moduleID)
		blog.V(9).Infof("traversal find module: %d result: %+v", moduleID, path)

		hostID, err := util.GetInt64ByInterface(item[common.BKHostIDField])
		if err != nil {
			err = fmt.Errorf("get host:%+v layer failed, err: %+v", item, err)
		}

		// prepare host layer
		var innerIP string
		var exist bool
		innerIP, exist = hostIDInnerIPMap[hostID]
		if exist == false {
			innerIP = fmt.Sprintf("host:%d", hostID)
		}
		hostLayer := meta.Item{
			Type:       meta.HostInstance,
			Name:       innerIP,
			InstanceID: hostID,
		}

		// layers from topo instance tree
		layers := make([]meta.Item, 0)
		for i := len(path) - 1; i >= 0; i-- {
			node := path[i]
			item := meta.Item{
				Name:       node.Name(),
				InstanceID: node.InstanceID,
				Type:       meta.GetResourceTypeByObjectType(node.ObjectID),
			}
			layers = append(layers, item)
		}
		layers = append(layers, hostLayer)
		blog.V(9).Infof("layers from traversal find module:%d result: %+v", moduleID, layers)
		batchLayers = append(batchLayers, layers)
	}
	batchLayersStr, err := json.Marshal(batchLayers)
	if err != nil {
		blog.Errorf("json encode GetHostLayers failed, err: %+v", err)
		err = fmt.Errorf("json encode GetHostLayers failed, err: %+v", err)
		return
	}
	blog.V(5).Infof("batchLayersStr: %s", batchLayersStr)

	return
}

func getInnerIPByHostIDs(coreService coreservice.CoreServiceClientInterface, rHeader http.Header, hostIDArr *[]int64) (hostIDInnerIPMap map[int64]string, err error) {
	hostIDInnerIPMap = map[int64]string{}

	cond := condition.CreateCondition()
	cond.Field(common.BKHostIDField).In(*hostIDArr)
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKHostInnerIPField, common.BKHostIDField},
		Condition: cond.ToMapStr(),
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	hosts, err := coreService.Instance().ReadInstance(
		context.Background(), rHeader, common.BKInnerObjIDHost, query)
	if err != nil {
		err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, err)
		return
	}
	for _, host := range hosts.Data.Info {
		hostID, e := util.GetInt64ByInterface(host[common.BKHostIDField])
		if e != nil {
			err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDArr, e)
			return
		}
		innerIP := util.GetStrByInterface(host[common.BKHostInnerIPField])
		hostIDInnerIPMap[hostID] = innerIP
	}
	return
}
