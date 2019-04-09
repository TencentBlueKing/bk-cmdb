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
	"configcenter/src/common/mapstr"
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

func (am *AuthManager) CollectHostByBusinessID(ctx context.Context, header http.Header, businessID int64) ([]HostSimplify, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(businessID)
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKHostInnerIPField, common.BKHostIDField},
		Condition: cond.ToMapStr(),
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	hosts, err := am.clientSet.CoreService().Instance().ReadInstance(context.Background(), header, common.BKTableNameModuleHostConfig, query)
	if err != nil {
		blog.Errorf("get host:%+v by businessID:%d failed, err: %+v", businessID, err)
		return nil, fmt.Errorf("get host by businessID:%d failed, err: %+v", businessID, err)
	}

	// extract hostID
	hostIDs := make([]int64, 0)
	for _, host := range hosts.Data.Info {
		hostIDVal, exist := host[common.BKHostIDField]
		if exist == false {
			continue
		}
		hostID, err := util.GetInt64ByInterface(hostIDVal)
		if err != nil {
			blog.V(2).Infof("synchronize task skip host:%+v, as parse hostID field failed, err: %+v", host, err)
			continue
		}
		hostIDs = append(hostIDs, hostID)
	}

	blog.V(4).Infof("list hosts by business:%d result: %+v", businessID, hostIDs)
	return am.collectHostByHostIDs(ctx, header, hostIDs...)
}

func (am *AuthManager) constructHostFromSearchResult(ctx context.Context, header http.Header, rawData []mapstr.MapStr) ([]HostSimplify, error) {
	hostIDs := make([]int64, 0)
	hosts := make([]HostSimplify, 0)
	for _, cls := range rawData {
		host := HostSimplify{}
		if _, err := host.Parse(cls); err != nil {
			return nil, fmt.Errorf("get hosts by object failed, err: %+v", err)
		}
		hosts = append(hosts, host)
		hostIDs = append(hostIDs, host.BKHostIDField)
	}

	// inject business,set,module info to HostSimplify
	hostModuleCondition := condition.CreateCondition()
	hostModuleCondition.Field(common.BKHostIDField).In(hostIDs)
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKAppIDField, common.BKModuleIDField, common.BKSetIDField, common.BKHostIDField},
		Condition: hostModuleCondition.ToMapStr(),
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	hostModuleResult, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameModuleHostConfig, query)
	if err != nil {
		err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDs, err)
		return nil, err
	}
	blog.V(5).Infof("get host module config: %+v", hostModuleResult.Data.Info)
	if len(rawData) == 0 {
		err = fmt.Errorf("get host:%+v layer failed, get host module config by host id not found, maybe hostID invalid", hostIDs)
		return nil, err
	}
	hostModuleMap := map[int64]HostSimplify{}
	for _, cls := range hostModuleResult.Data.Info {
		host := HostSimplify{}
		_, err = host.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get hosts by object failed, err: %+v", err)
		}
		hostModuleMap[host.BKHostIDField] = host
	}
	blog.V(9).Infof("hostModuleMap: %+v", hostModuleMap)
	for idx, host := range hosts {
		hostModule, exist := hostModuleMap[host.BKHostIDField]
		if exist == false {
			return nil, fmt.Errorf("hostID:%+d doesn't exist in any module", host.BKHostIDField)
		}
		hosts[idx].BKAppIDField = hostModule.BKAppIDField
		hosts[idx].BKSetIDField = hostModule.BKSetIDField
		hosts[idx].BKModuleIDField = hostModule.BKModuleIDField
	}
	blog.V(9).Infof("hosts: %+v", hosts)
	return hosts, nil
}

func (am *AuthManager) collectHostByHostIDs(ctx context.Context, header http.Header, hostIDs ...int64) ([]HostSimplify, error) {
	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	hostIDs = util.IntArrayUnique(hostIDs)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKHostIDField).In(hostIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDHost, &cond)
	if err != nil {
		blog.V(3).Infof("get hosts by id failed, err: %+v", err)
		return nil, fmt.Errorf("get hosts by id failed, err: %+v", err)
	}

	return am.constructHostFromSearchResult(ctx, header, result.Data.Info)
}

func (am *AuthManager) extractBusinessIDFromHosts(ctx context.Context, header http.Header, hosts ...HostSimplify) (int64, error) {
	var businessID int64
	for idx, host := range hosts {
		bizID := host.BKAppIDField
		// we should ignore metadata.LabelBusinessID field not found error
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("authorization failed, get multiple business ID from hosts")
		}
		businessID = bizID
	}
	return am.correctBusinessID(ctx, header, businessID)
}

func (am *AuthManager) MakeResourcesByHosts(header http.Header, action meta.Action, businessID int64, hosts ...HostSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, host := range hosts {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.HostInstance,
				Name:       host.BKHostInnerIPField,
				InstanceID: host.BKHostIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	blog.V(9).Infof("host resources for iam: %+v", resources)
	return resources
}

func (am *AuthManager) makeHostsResourcesGroupByBusiness(header http.Header, action meta.Action, hosts ...HostSimplify) map[int64][]meta.ResourceAttribute {
	result := map[int64][]meta.ResourceAttribute{}
	for _, host := range hosts {
		bizID := host.BKAppIDField
		_, exist := result[bizID]
		if exist == false {
			result[bizID] = make([]meta.ResourceAttribute, 0)
		}
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.HostInstance,
				Name:       host.BKHostInnerIPField,
				InstanceID: host.BKHostIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      bizID,
		}

		result[bizID] = append(result[bizID], resource)
	}
	blog.V(9).Infof("host resources group by business for iam: %+v", result)
	return result
}

func (am *AuthManager) AuthorizeHostsCrossMultipleBusiness(ctx context.Context, header http.Header, action meta.Action, hosts ...HostSimplify) error {

	// make auth resources
	bizResourcesMap := am.makeHostsResourcesGroupByBusiness(header, action, hosts...)

	for bizID, resources := range bizResourcesMap {
		err := am.authorize(ctx, header, bizID, resources...)
		if err != nil {
			return fmt.Errorf("authorize resources with businessID %d failed, err: %+v", bizID, err)
		}
	}
	return nil
}

func (am *AuthManager) AuthorizeByHosts(ctx context.Context, header http.Header, action meta.Action, hosts ...HostSimplify) error {
	if len(hosts) == 0 {
		return nil
	}
	// extract business id
	bizID, err := am.extractBusinessIDFromHosts(ctx, header, hosts...)
	if err != nil {
		return fmt.Errorf("authorize hosts failed, extract business id from hosts failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByHosts(header, action, bizID, hosts...)
	return am.authorize(ctx, header, bizID, resources...)
}

func (am *AuthManager) AuthorizeByHostsIDs(ctx context.Context, header http.Header, action meta.Action, hostIDs ...int64) error {
	if am.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, hosts: %+v", hostIDs)
		return nil
	}

	if len(hostIDs) == 0 {
		return nil
	}
	hosts, err := am.collectHostByHostIDs(ctx, header, hostIDs...)
	if err != nil {
		return fmt.Errorf("authorize hosts failed, get hosts by id failed, err: %+v", err)
	}
	return am.AuthorizeByHosts(ctx, header, action, hosts...)
}

func (am *AuthManager) DryRunAuthorizeByHostsIDs(ctx context.Context, header http.Header, action meta.Action, hostIDs ...int64) ([]meta.ResourceAttribute, error) {
	if len(hostIDs) == 0 {
		return nil, nil
	}

	hosts, err := am.collectHostByHostIDs(ctx, header, hostIDs...)
	if err != nil {
		return nil, fmt.Errorf("authorize hosts failed, get hosts by id failed, err: %+v", err)
	}
	// extract business id
	bizID, err := am.extractBusinessIDFromHosts(ctx, header, hosts...)
	if err != nil {
		return nil, fmt.Errorf("authorize hosts failed, extract business id from hosts failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByHosts(header, action, bizID, hosts...)

	realResources, err := am.Authorize.DryRunRegisterResource(context.Background(), resources...)
	if err != nil {
		blog.Errorf("dry dun register failed, err: %+v", err)
		return nil, err
	}
	blog.InfoJSON("realResources is: %s", realResources)
	return resources, nil
}

func (am *AuthManager) AuthorizeCreateHost(ctx context.Context, header http.Header, bizID int64) error {
	return am.AuthorizeResourceCreate(ctx, header, bizID, meta.HostInstance)
}

func (am *AuthManager) UpdateRegisteredHosts(ctx context.Context, header http.Header, hosts ...HostSimplify) error {
	if len(hosts) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromHosts(ctx, header, hosts...)
	if err != nil {
		return fmt.Errorf("authorize hosts failed, extract business id from hosts failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByHosts(header, meta.EmptyAction, bizID, hosts...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredHostsByID(ctx context.Context, header http.Header, hostIDs ...int64) error {
	if len(hostIDs) == 0 {
		return nil
	}

	hosts, err := am.collectHostByHostIDs(ctx, header, hostIDs...)
	if err != nil {
		return fmt.Errorf("update registered hosts failed, get hosts by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredHosts(ctx, header, hosts...)
}

func (am *AuthManager) DeregisterHostsByID(ctx context.Context, header http.Header, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}

	hosts, err := am.collectHostByHostIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("deregister hosts failed, get hosts by id failed, err: %+v", err)
	}
	return am.DeregisterHosts(ctx, header, hosts...)
}

func (am *AuthManager) RegisterHosts(ctx context.Context, header http.Header, hosts ...HostSimplify) error {
	if len(hosts) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromHosts(ctx, header, hosts...)
	if err != nil {
		return fmt.Errorf("register hosts failed, extract business id from hosts failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByHosts(header, meta.EmptyAction, bizID, hosts...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) RegisterHostsByID(ctx context.Context, header http.Header, hostIDs ...int64) error {
	if len(hostIDs) == 0 {
		return nil
	}

	hosts, err := am.collectHostByHostIDs(ctx, header, hostIDs...)
	if err != nil {
		return fmt.Errorf("register host failed, get hosts by id failed, err: %+v", err)
	}
	return am.RegisterHosts(ctx, header, hosts...)
}

func (am *AuthManager) DeregisterHosts(ctx context.Context, header http.Header, hosts ...HostSimplify) error {
	if len(hosts) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromHosts(ctx, header, hosts...)
	if err != nil {
		return fmt.Errorf("deregister hosts failed, extract business id from hosts failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByHosts(header, meta.EmptyAction, bizID, hosts...)

	return am.Authorize.DeregisterResource(ctx, resources...)
}
