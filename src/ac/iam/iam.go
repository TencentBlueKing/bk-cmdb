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

type Iam struct {
	client *iamClient
}

func NewIam(tls *util.TLSClientConfig, cfg AuthConfig, reg prometheus.Registerer) (*Iam, error) {
	blog.V(5).Infof("new iam with parameters tls: %+v, cfg: %+v", tls, cfg)
	if !auth.EnableAuthorize() {
		return new(Iam), nil
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
		Reg: reg,
	}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	header.Set(iamAppCodeHeader, cfg.AppCode)
	header.Set(iamAppSecretHeader, cfg.AppSecret)

	return &Iam{
		client: &iamClient{
			Config:      cfg,
			client:      rest.NewRESTClient(c, ""),
			basicHeader: header,
		},
	}, nil
}

func (i Iam) RegisterSystem(ctx context.Context, host string) error {
	systemResp, err := i.client.GetSystemInfo(ctx)
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
	for _, resourceType := range systemResp.Data.ResourceTypes {
		existResourceTypeMap[resourceType.ID] = true
		removedResourceTypeMap[resourceType.ID] = struct{}{}
	}
	newResourceTypes := make([]ResourceType, 0)
	for _, resourceType := range GenerateResourceTypes() {
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
	if len(newResourceTypes) > 0 {
		if err = i.client.RegisterResourcesTypes(ctx, newResourceTypes); err != nil {
			blog.ErrorJSON("register resource types failed, error: %s, resource types: %s", err.Error(), newResourceTypes)
			return err
		}
	}

	existInstanceSelectionMap := make(map[InstanceSelectionID]bool)
	removedInstanceSelectionMap := make(map[InstanceSelectionID]struct{})
	for _, instanceSelection := range systemResp.Data.InstanceSelections {
		existInstanceSelectionMap[instanceSelection.ID] = true
		removedInstanceSelectionMap[instanceSelection.ID] = struct{}{}
	}
	newInstanceSelections := make([]InstanceSelection, 0)
	for _, resourceType := range GenerateInstanceSelections() {
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
	if len(newInstanceSelections) > 0 {
		if err = i.client.CreateInstanceSelection(ctx, newInstanceSelections); err != nil {
			blog.ErrorJSON("register instance selections failed, error: %s, resource types: %s", err.Error(), newInstanceSelections)
			return err
		}
	}

	existResourceActionMap := make(map[ActionID]bool)
	removedResourceActionMap := make(map[ActionID]struct{})
	for _, resourceAction := range systemResp.Data.Actions {
		existResourceActionMap[resourceAction.ID] = true
		removedResourceActionMap[resourceAction.ID] = struct{}{}
	}
	newResourceActions := make([]ResourceAction, 0)
	for _, resourceAction := range GenerateActions() {
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
	if len(newResourceActions) > 0 {
		if err = i.client.CreateAction(ctx, newResourceActions); err != nil {
			blog.ErrorJSON("register resource actions failed, error: %s, resource actions: %s", err.Error(), newResourceActions)
			return err
		}
	}

	// remove redundant actions, redundant instance selections and resource types one by one when dependencies are all deleted
	if actionLen := len(removedResourceActionMap); actionLen > 0 {
		removedResourceActionIDs := make([]ActionID, actionLen)
		idx := 0
		for resourceActionID, _ := range removedResourceActionMap {
			removedResourceActionIDs[idx] = resourceActionID
			idx++
		}
		if err = i.client.DeleteAction(ctx, removedResourceActionIDs); err != nil {
			blog.ErrorJSON("delete resource actions failed, error: %s, resource actions: %s", err.Error(), removedResourceActionIDs)
			return err
		}
	}
	if selectionLen := len(removedInstanceSelectionMap); selectionLen > 0 {
		removedInstanceSelectionIDs := make([]InstanceSelectionID, selectionLen)
		idx := 0
		for resourceActionID, _ := range removedInstanceSelectionMap {
			removedInstanceSelectionIDs[idx] = resourceActionID
			idx++
		}
		if err = i.client.DeleteInstanceSelection(ctx, removedInstanceSelectionIDs); err != nil {
			blog.ErrorJSON("delete instance selections failed, error: %s, instance selections: %s", err.Error(), removedInstanceSelectionIDs)
			return err
		}
	}
	if typeLen := len(removedResourceTypeMap); typeLen > 0 {
		removedResourceTypeIDs := make([]TypeID, len(removedResourceTypeMap))
		idx := 0
		for resourceType, _ := range removedResourceTypeMap {
			removedResourceTypeIDs[idx] = resourceType
			idx++
		}
		if err = i.client.DeleteResourcesTypes(ctx, removedResourceTypeIDs); err != nil {
			blog.ErrorJSON("delete resource types failed, error: %s, resource types: %s", err.Error(), removedResourceTypeIDs)
			return err
		}
	}

	// register or update resource action groups
	actionGroups := GenerateActionGroups()
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

	return nil
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

func (a *authorizer) RegisterResourceCreatorAction(ctx context.Context, h http.Header, input metadata.IamInstanceWithCreator) (
	[]metadata.IamCreatorActionPolicy, error) {

	return a.authClientSet.RegisterResourceCreatorAction(ctx, h, input)
}

func (a *authorizer) BatchRegisterResourceCreatorAction(ctx context.Context, h http.Header, input metadata.IamInstancesWithCreator) (
	[]metadata.IamCreatorActionPolicy, error) {

	return a.authClientSet.BatchRegisterResourceCreatorAction(ctx, h, input)
}
