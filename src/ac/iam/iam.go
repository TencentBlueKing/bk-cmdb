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

package iam

import (
	"context"
	"net/http"

	"configcenter/src/ac"
	"configcenter/src/ac/meta"
	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/authserver"
	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/sdk/types"

	"github.com/prometheus/client_golang/prometheus"
)

type IAM struct {
	client iamClientInterface
}

func NewIAM(tls *util.TLSClientConfig, cfg AuthConfig, reg prometheus.Registerer) (*IAM, error) {
	blog.V(5).Infof("new iam with parameters tls: %+v, cfg: %+v", tls, cfg)
	if !auth.EnableAuthorize() {
		return new(IAM), nil
	}
	client, err := util.NewClient(tls)
	if err != nil {
		return nil, err
	}
	c := &util.Capability{
		Client: client,
		Discover: &iamDiscovery{
			servers: cfg.Address,
		},
		Throttle: flowctrl.NewRateLimiter(5000, 5000),
		Mock: util.MockInfo{
			Mocked: false,
		},
		MetricOpts: util.MetricOption{Register: reg},
	}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	header.Set(iamAppCodeHeader, cfg.AppCode)
	header.Set(iamAppSecretHeader, cfg.AppSecret)

	return &IAM{
		client: NewIAMClient(&IAMClientCfg{
			Config:      cfg,
			Client:      rest.NewRESTClient(c, ""),
			BasicHeader: header,
		}),
	}, nil
}

func (i *IAM) RegisterSystem(ctx context.Context, host string, objects []metadata.Object) error {
	systemResp, err := i.client.GetSystemInfo(ctx, []SystemQueryField{})
	if err != nil && err != ErrNotFound {
		blog.Errorf("get system info failed, error: %s", err.Error())
		return err
	}
	// if iam cmdb system has not been registered, register system
	if err == ErrNotFound {
		sys := System{
			ID:          SystemIDCMDB,
			Name:        SystemNameCMDB,
			EnglishName: SystemNameCMDBEn,
			Clients:     SystemIDCMDB,
			ProviderConfig: &SysConfig{
				Host: host,
				Auth: "basic",
			},
		}
		if err = i.client.RegisterSystem(ctx, sys); err != nil {
			blog.ErrorJSON("register system %s failed, error: %s", sys, err.Error())
			return err
		}
		blog.V(5).Infof("register new system %+v succeed", sys)
	} else if systemResp.Data.BaseInfo.ProviderConfig == nil || systemResp.Data.BaseInfo.ProviderConfig.Host != host {
		// if iam registered cmdb system has no ProviderConfig or registered host config is different with current host config, update system host config
		if err = i.client.UpdateSystemConfig(ctx, &SysConfig{Host: host}); err != nil {
			blog.Errorf("update system host %s config failed, error: %s", host, err.Error())
			return err
		}
		if systemResp.Data.BaseInfo.ProviderConfig == nil {
			blog.V(5).Infof("update system host to %s succeed", systemResp.Data.BaseInfo.ProviderConfig.Host, host)
		} else {
			blog.V(5).Infof("update system host %s to %s succeed", systemResp.Data.BaseInfo.ProviderConfig.Host, host)
		}
	}

	existResourceTypeMap := make(map[TypeID]bool)
	removedResourceTypeMap := make(map[TypeID]struct{})
	newResourceTypes := make([]ResourceType, 0)
	for _, resourceType := range systemResp.Data.ResourceTypes {
		existResourceTypeMap[resourceType.ID] = true
		removedResourceTypeMap[resourceType.ID] = struct{}{}
	}
	for _, resourceType := range GenerateResourceTypes(objects) {
		// registered resource type exist in current resource types, should not be removed
		delete(removedResourceTypeMap, resourceType.ID)
		// if current resource type is registered, update it, or else register it
		if existResourceTypeMap[resourceType.ID] {
			if err = i.client.UpdateResourcesTypes(ctx, resourceType); err != nil {
				blog.ErrorJSON("update resource type %s failed, error: %s, input resource type: %s", resourceType.ID, err.Error(), resourceType)
				return err
			}
		} else {
			newResourceTypes = append(newResourceTypes, resourceType)
		}
	}

	existInstanceSelectionMap := make(map[InstanceSelectionID]bool)
	removedInstanceSelectionMap := make(map[InstanceSelectionID]struct{})
	newInstanceSelections := make([]InstanceSelection, 0)
	for _, instanceSelection := range systemResp.Data.InstanceSelections {
		existInstanceSelectionMap[instanceSelection.ID] = true
		removedInstanceSelectionMap[instanceSelection.ID] = struct{}{}
	}

	for _, resourceType := range GenerateInstanceSelections(objects) {
		// registered instance selection exist in current instance selections, should not be removed
		delete(removedInstanceSelectionMap, resourceType.ID)
		// if current instance selection is registered, update it, or else register it
		if existInstanceSelectionMap[resourceType.ID] {
			if err = i.client.UpdateInstanceSelection(ctx, resourceType); err != nil {
				blog.ErrorJSON("update instance selection %s failed, error: %s, input resource type: %s", resourceType.ID, err.Error(), resourceType)
				return err
			}
		} else {
			newInstanceSelections = append(newInstanceSelections, resourceType)
		}
	}

	existResourceActionMap := make(map[ActionID]bool)
	removedResourceActionMap := make(map[ActionID]struct{})
	newResourceActions := make([]ResourceAction, 0)
	for _, resourceAction := range systemResp.Data.Actions {
		existResourceActionMap[resourceAction.ID] = true
		removedResourceActionMap[resourceAction.ID] = struct{}{}
	}
	for _, resourceAction := range GenerateActions(objects) {
		// registered resource action exist in current resource actions, should not be removed
		delete(removedResourceActionMap, resourceAction.ID)
		// if current resource action is registered, update it, or else register it
		if existResourceActionMap[resourceAction.ID] {
			if err = i.client.UpdateAction(ctx, resourceAction); err != nil {
				blog.ErrorJSON("update resource action %s failed, error: %s, input resource action: %s", resourceAction.ID, err.Error(), resourceAction)
				return err
			}
		} else {
			newResourceActions = append(newResourceActions, resourceAction)
		}
	}

	// 因为资源间的依赖关系，删除和更新的顺序为 1.Action 2.InstanceSelection 3.ResourceType
	// 因为资源间的依赖关系，新建的顺序则反过来为 1.ResourceType 2.InstanceSelection 3.Action
	// ActionGroup依赖于Action，该资源的增删操作始终放在最后
	// 先删除资源，再新增资源，因为实例视图的名称在系统中是唯一的，如果不先删，同样名称的实例视图将创建失败
	// remove redundant actions, redundant instance selections and resource types one by one when dependencies are all deleted
	if len(removedResourceActionMap) > 0 {
		removedResourceActionIDs := make([]ActionID, len(removedResourceActionMap))
		idx := 0
		for resourceActionID := range removedResourceActionMap {
			removedResourceActionIDs[idx] = resourceActionID
			idx++
		}
		if err = i.client.DeleteActions(ctx, removedResourceActionIDs); err != nil {
			blog.ErrorJSON("delete resource actions failed, error: %s, resource actions: %s", err.Error(), removedResourceActionIDs)
			return err
		}
	}

	if len(removedInstanceSelectionMap) > 0 {
		removedInstanceSelectionIDs := make([]InstanceSelectionID, len(removedInstanceSelectionMap))
		idx := 0
		for resourceActionID := range removedInstanceSelectionMap {
			removedInstanceSelectionIDs[idx] = resourceActionID
			idx++
		}
		if err = i.client.DeleteInstanceSelections(ctx, removedInstanceSelectionIDs); err != nil {
			blog.ErrorJSON("delete instance selections failed, error: %s, instance selections: %s", err.Error(), removedInstanceSelectionIDs)
			return err
		}
	}

	if len(removedResourceTypeMap) > 0 {
		removedResourceTypeIDs := make([]TypeID, len(removedResourceTypeMap))
		idx := 0
		for resourceType := range removedResourceTypeMap {
			removedResourceTypeIDs[idx] = resourceType
			idx++
		}
		if err = i.client.DeleteResourcesTypes(ctx, removedResourceTypeIDs); err != nil {
			blog.ErrorJSON("delete resource types failed, error: %s, resource types: %s", err.Error(), removedResourceTypeIDs)
			return err
		}
	}

	if len(newResourceTypes) > 0 {
		if err = i.client.RegisterResourcesTypes(ctx, newResourceTypes); err != nil {
			blog.ErrorJSON("register resource types failed, error: %s, resource types: %s", err.Error(), newResourceTypes)
			return err
		}
	}

	if len(newInstanceSelections) > 0 {
		if err = i.client.RegisterInstanceSelections(ctx, newInstanceSelections); err != nil {
			blog.ErrorJSON("register instance selections failed, error: %s, resource types: %s", err.Error(), newInstanceSelections)
			return err
		}
	}

	if len(newResourceActions) > 0 {
		if err = i.client.RegisterActions(ctx, newResourceActions); err != nil {
			blog.ErrorJSON("register resource actions failed, error: %s, resource actions: %s", err.Error(), newResourceActions)
			return err
		}
	}

	// register or update resource action groups
	actionGroups := GenerateActionGroups(objects)
	if len(systemResp.Data.ActionGroups) == 0 {
		if err = i.client.RegisterActionGroups(ctx, actionGroups); err != nil {
			blog.ErrorJSON("register action groups failed, error: %s, action groups: %s", err.Error(), actionGroups)
			return err
		}
	} else {
		if err = i.client.UpdateActionGroups(ctx, actionGroups); err != nil {
			blog.ErrorJSON("update action groups failed, error: %s, action groups: %s", err.Error(), actionGroups)
			return err
		}
	}

	// register or update resource creator actions
	resourceCreatorActions := GenerateResourceCreatorActions()
	if len(systemResp.Data.ResourceCreatorActions.Config) == 0 {
		if err = i.client.RegisterResourceCreatorActions(ctx, resourceCreatorActions); err != nil {
			blog.ErrorJSON("register resource creator actions failed, error: %s, resource creator actions: %s", err.Error(), resourceCreatorActions)
			return err
		}
	} else {
		if err = i.client.UpdateResourceCreatorActions(ctx, resourceCreatorActions); err != nil {
			blog.ErrorJSON("update resource creator actions failed, error: %s, resource creator actions: %s", err.Error(), resourceCreatorActions)
			return err
		}
	}

	//register or update common actions
	commonActions := GenerateCommonActions()
	if len(systemResp.Data.CommonActions) == 0 {
		if err = i.client.RegisterCommonActions(ctx, commonActions); err != nil {
			blog.ErrorJSON("register common actions failed, error: %s, common actions: %s", err.Error(), commonActions)
			return err
		}
	} else {
		if err = i.client.UpdateCommonActions(ctx, commonActions); err != nil {
			blog.ErrorJSON("update common actions failed, error: %s, common actions: %s", err.Error(), commonActions)
			return err
		}
	}

	return nil
}

// SyncIAMSysInstances sync system instances between CMDB and IAM
// it check the difference of system instances resource between CMDB and IAM
// if they have difference, sync and make them same
func (i *IAM) SyncIAMSysInstances(ctx context.Context, objects []metadata.Object) error {

	fields := []SystemQueryField{FieldResourceTypes,
		FieldActions, FieldActionGroups, FieldInstanceSelections}
	iamResp, err := i.client.GetSystemInfo(ctx, fields)
	if err != nil {
		blog.ErrorJSON("syc iam sysInstances failed, get system info error: %s, fields: %s",
			err.Error(), fields)
		return err
	}

	// get the cmdb resources
	cmdbActions := GenModelInstanceActions(objects)
	cmdbInstanceSelections := GenDynamicInstanceSelections(objects)
	cmdbResourceTypes := GenDynamicResourceTypes(objects)
	cmdbActionGroups := GenerateActionGroups(objects)

	// compare resources between cmdb and iam
	addActions, deleteActions := compareActions(cmdbActions, iamResp.Data.Actions)
	addInstanceSelections, deleteInstanceSelections := compareInstanceSelections(cmdbInstanceSelections,
		iamResp.Data.InstanceSelections)
	addResourceTypes, deleteResourceTypes := compareResourceTypes(cmdbResourceTypes, iamResp.Data.ResourceTypes)

	// 因为资源间的依赖关系，删除和更新的顺序为 1.Action 2.InstanceSelection 3.ResourceType
	// 因为资源间的依赖关系，新建的顺序则反过来为 1.ResourceType 2.InstanceSelection 3.Action
	// ActionGroup依赖于Action，该资源的增删操作始终放在最后
	// 先删除资源，再新增资源，因为实例视图的名称在系统中是唯一的，如果不先删，同样名称的实例视图将创建失败

	// delete unnecessary actions in iam
	if len(deleteActions) > 0 {
		blog.Infof("begin delete actions, count:%d, detail:%s", len(deleteActions),
			deleteActions)
		if err := i.client.DeleteActions(ctx, deleteActions); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, delete IAM actions failed, error: %s, actions: %s",
				err.Error(), deleteActions)
			return err
		}
	}

	// delete unnecessary InstanceSelections in iam
	if len(deleteInstanceSelections) > 0 {
		blog.Infof("begin delete instanceSelections, count:%d, detail:%s", len(deleteInstanceSelections),
			deleteInstanceSelections)
		if err := i.client.DeleteInstanceSelections(ctx, deleteInstanceSelections); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, delete instanceSelections error: %s, instanceSelections: %s",
				err.Error(), deleteInstanceSelections)
			return err
		}
	}

	// delete unnecessary ResourceTypes in iam
	if len(deleteResourceTypes) > 0 {
		blog.Infof("begin delete resourceTypes, count:%d, detail:%s", len(deleteResourceTypes),
			deleteResourceTypes)
		if err := i.client.DeleteResourcesTypes(ctx, deleteResourceTypes); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, delete resourceType error: %s, resourceType: %s",
				err.Error(), deleteResourceTypes)
			return err
		}
	}

	// add cmdb ResourceTypes in iam
	if len(addResourceTypes) > 0 {
		blog.Infof("begin add resourceTypes, count:%d, detail:%s", len(addResourceTypes),
			addResourceTypes)
		if err := i.client.RegisterResourcesTypes(ctx, addResourceTypes); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, add resourceType error: %s, resourceType: %s",
				err.Error(), addResourceTypes)
			return err
		}
	}

	// add cmdb InstanceSelections in iam
	if len(addInstanceSelections) > 0 {
		blog.Infof("begin add instanceSelections, count:%d, detail:%s", len(addInstanceSelections),
			addInstanceSelections)
		if err := i.client.RegisterInstanceSelections(ctx, addInstanceSelections); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, add instanceSelections error: %s, instanceSelections: %s",
				err.Error(), addInstanceSelections)
			return err
		}
	}

	// add cmdb actions in iam
	if len(addActions) > 0 {
		blog.Infof("begin add actions, count:%d, detail:%s", len(addActions), addActions)
		if err := i.client.RegisterActions(ctx, addActions); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, add IAM actions failed, error: %s, actions: %s",
				err.Error(), addActions)
			return err
		}
	}

	// update action_groups in iam
	if len(addActions) > 0 || len(deleteActions) > 0 {
		blog.Infof("begin update actionGroups")
		if err := i.client.UpdateActionGroups(ctx, cmdbActionGroups); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, update actionGroups error: %s, actionGroups: %s",
				err.Error(), cmdbActionGroups)
			return err
		}
	}

	return nil
}

// compareActions compare actions between cmdb and iam
func compareActions(cmdbActions []ResourceAction, iamActions []ResourceAction) (
	addActions []ResourceAction, deleteActionIDs []ActionID) {
	cmdbActionMap := map[ActionID]struct{}{}
	iamActionMap := map[ActionID]struct{}{}

	for _, action := range iamActions {
		if IsOldIAMSysInstanceActionID(action.ID) || IsIAMSysInstanceAction(action.ID) {
			iamActionMap[action.ID] = struct{}{}
		}
	}

	for _, action := range cmdbActions {
		cmdbActionMap[action.ID] = struct{}{}
		if _, ok := iamActionMap[action.ID]; !ok {
			addActions = append(addActions, action)
		}
	}

	for actionID := range iamActionMap {
		if _, ok := cmdbActionMap[actionID]; !ok {
			deleteActionIDs = append(deleteActionIDs, actionID)
		}
	}

	return addActions, deleteActionIDs
}

// compareInstanceSelections compare instanceSelections between cmdb and iam
func compareInstanceSelections(cmdbInstanceSelections []InstanceSelection,
	iamInstanceSelections []InstanceSelection) (addInstanceSelection []InstanceSelection,
	deleteInstanceSelectionIDs []InstanceSelectionID) {
	cmdbInstanceSelectionMap := map[InstanceSelectionID]struct{}{}
	iamInstanceSelectionMap := map[InstanceSelectionID]struct{}{}

	for _, instanceSelection := range iamInstanceSelections {
		if IsIAMSysInstanceSelection(instanceSelection.ID) || IsIAMSysInstanceSelection(instanceSelection.ID) {
			iamInstanceSelectionMap[instanceSelection.ID] = struct{}{}
		}
	}

	for _, instanceSelection := range cmdbInstanceSelections {
		cmdbInstanceSelectionMap[instanceSelection.ID] = struct{}{}
		if _, ok := iamInstanceSelectionMap[instanceSelection.ID]; !ok {
			addInstanceSelection = append(addInstanceSelection, instanceSelection)
		}
	}

	for instanceSelectionID := range iamInstanceSelectionMap {
		if _, ok := cmdbInstanceSelectionMap[instanceSelectionID]; !ok {
			deleteInstanceSelectionIDs = append(deleteInstanceSelectionIDs, instanceSelectionID)
		}
	}

	return addInstanceSelection, deleteInstanceSelectionIDs
}

// compareResourceTypes compare resourceTypes between cmdb and iam
func compareResourceTypes(cmdbResourceTypes []ResourceType, iamResourceTypes []ResourceType) (
	addResourceTypes []ResourceType, deleteTypeIDs []TypeID) {
	cmdbResourceTypeMap := map[TypeID]struct{}{}
	iamResourceTypeMap := map[TypeID]struct{}{}

	for _, resourceType := range iamResourceTypes {
		if IsOldIAMSysInstanceTypeID(resourceType.ID) || IsIAMSysInstance(resourceType.ID) {
			iamResourceTypeMap[resourceType.ID] = struct{}{}
		}
	}

	for _, resourceType := range cmdbResourceTypes {
		cmdbResourceTypeMap[resourceType.ID] = struct{}{}
		if _, ok := iamResourceTypeMap[resourceType.ID]; !ok {
			addResourceTypes = append(addResourceTypes, resourceType)
		}
	}

	for typeID := range iamResourceTypeMap {
		if _, ok := cmdbResourceTypeMap[typeID]; !ok {
			deleteTypeIDs = append(deleteTypeIDs, typeID)
		}
	}

	return addResourceTypes, deleteTypeIDs
}

type authorizer struct {
	authClientSet authserver.AuthServerClientInterface
}

func NewAuthorizer(clientSet apimachinery.ClientSetInterface) ac.AuthorizeInterface {
	return &authorizer{authClientSet: clientSet.AuthServer()}
}

func (a *authorizer) AuthorizeBatch(ctx context.Context, h http.Header, user meta.UserInfo,
	resources ...meta.ResourceAttribute) ([]types.Decision, error) {
	return a.authorizeBatch(ctx, h, true, user, resources...)
}

func (a *authorizer) AuthorizeAnyBatch(ctx context.Context, h http.Header, user meta.UserInfo,
	resources ...meta.ResourceAttribute) ([]types.Decision, error) {
	return a.authorizeBatch(ctx, h, false, user, resources...)
}

func (a *authorizer) authorizeBatch(ctx context.Context, h http.Header, exact bool, user meta.UserInfo,
	resources ...meta.ResourceAttribute) ([]types.Decision, error) {

	rid := commonutil.GetHTTPCCRequestID(h)

	opts, decisions, err := parseAttributesToBatchOptions(rid, user, resources...)
	if err != nil {
		return nil, err
	}

	// all resources are skipped
	if opts == nil {
		return decisions, nil
	}

	var authDecisions []types.Decision
	if exact {
		authDecisions, err = a.authClientSet.AuthorizeBatch(ctx, h, opts)
		if err != nil {
			blog.ErrorJSON("authorize batch failed, err: %s, ops: %s, resources: %s, rid: %s", err, opts, resources, rid)
			return nil, err
		}
	} else {
		authDecisions, err = a.authClientSet.AuthorizeAnyBatch(ctx, h, opts)
		blog.ErrorJSON("authorize any batch failed, err: %s, ops: %s, resources: %s, rid: %s", err, opts, resources, rid)
		if err != nil {
			return nil, err
		}

	}

	index := 0
	for _, decision := range authDecisions {
		// skip resources' decisions are already set as authorized
		for decisions[index].Authorized {
			index++
		}
		decisions[index].Authorized = decision.Authorized
		index++
	}

	return decisions, nil
}

func parseAttributesToBatchOptions(rid string, user meta.UserInfo, resources ...meta.ResourceAttribute) (*types.AuthBatchOptions, []types.Decision, error) {
	if !auth.EnableAuthorize() {
		decisions := make([]types.Decision, len(resources))
		for i := range decisions {
			decisions[i].Authorized = true
		}
		return nil, decisions, nil
	}

	authBatchArr := make([]*types.AuthBatch, 0)
	decisions := make([]types.Decision, len(resources))
	for index, resource := range resources {

		// this resource should be skipped, do not need to verify in auth center.
		if resource.Action == meta.SkipAction {
			decisions[index].Authorized = true
			blog.V(5).Infof("skip authorization for resource: %+v, rid: %s", resource, rid)
			continue
		}

		action, resources, err := AdaptAuthOptions(&resource)
		if err != nil {
			blog.Errorf("adaptor cmdb resource to iam failed, err: %s, rid: %s", err, rid)
			return nil, nil, err
		}

		// this resource should be skipped, do not need to verify in auth center.
		if action == Skip {
			decisions[index].Authorized = true
			blog.V(5).Infof("skip authorization for resource: %+v, rid: %s", resource, rid)
			continue
		}

		authBatchArr = append(authBatchArr, &types.AuthBatch{
			Action:    types.Action{ID: string(action)},
			Resources: resources,
		})
	}

	// all resources are skipped
	if len(authBatchArr) == 0 {
		return nil, decisions, nil
	}

	ops := &types.AuthBatchOptions{
		System: SystemIDCMDB,
		Subject: types.Subject{
			Type: "user",
			ID:   user.UserName,
		},
		Batch: authBatchArr,
	}
	return ops, decisions, nil
}

func (a *authorizer) ListAuthorizedResources(ctx context.Context, h http.Header, input meta.ListAuthorizedResourcesParam) ([]string, error) {
	return a.authClientSet.ListAuthorizedResources(ctx, h, input)
}

func (a *authorizer) GetNoAuthSkipUrl(ctx context.Context, h http.Header, input *metadata.IamPermission) (string, error) {
	return a.authClientSet.GetNoAuthSkipUrl(ctx, h, input)
}

func (a *authorizer) GetPermissionToApply(ctx context.Context, h http.Header, input []meta.ResourceAttribute) (*metadata.IamPermission, error) {
	return a.authClientSet.GetPermissionToApply(ctx, h, input)
}

func (a *authorizer) RegisterResourceCreatorAction(ctx context.Context, h http.Header, input metadata.IamInstanceWithCreator) (
	[]metadata.IamCreatorActionPolicy, error) {

	return a.authClientSet.RegisterResourceCreatorAction(ctx, h, input)
}

func (a *authorizer) BatchRegisterResourceCreatorAction(ctx context.Context, h http.Header, input metadata.IamInstancesWithCreator) (
	[]metadata.IamCreatorActionPolicy, error) {

	return a.authClientSet.BatchRegisterResourceCreatorAction(ctx, h, input)
}
