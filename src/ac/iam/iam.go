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
	"time"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"

	"github.com/prometheus/client_golang/prometheus"
)

type Iam struct {
	client *iamClient
}

func NewIam(tls *util.TLSClientConfig, cfg AuthConfig, reg prometheus.Registerer) (*Iam, error) {
	blog.V(5).Infof("new iam with parameters tls: %+v, cfg: %+v", tls, cfg)
	if !auth.IsAuthed() {
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
	if err != nil {
		blog.Errorf("get system info failed, error: %s", err.Error())
		return err
	}
	// if iam cmdb system has not been registered, register system
	if systemResp.Data.BaseInfo.ID == "" {
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

	existResourceTypeMap := make(map[ResourceTypeID]bool)
	removedResourceTypeMap := make(map[ResourceTypeID]struct{})
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

	existResourceActionMap := make(map[ResourceActionID]bool)
	removedResourceActionMap := make(map[ResourceActionID]struct{})
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

	// remove redundant actions, then remove redundant resource types whose related actions are all deleted
	removedResourceActionIDs := make([]ResourceActionID, len(removedResourceActionMap))
	for resourceActionID, _ := range removedResourceActionMap {
		removedResourceActionIDs = append(removedResourceActionIDs, resourceActionID)
	}
	if err = i.client.DeleteAction(ctx, removedResourceActionIDs); err != nil {
		blog.ErrorJSON("delete resource actions failed, error: %s, resource actions: %s", err.Error(), removedResourceActionIDs)
		return err
	}
	removedResourceTypeIDs := make([]ResourceTypeID, len(removedResourceTypeMap))
	for resourceType, _ := range removedResourceTypeMap {
		removedResourceTypeIDs = append(removedResourceTypeIDs, resourceType)
	}
	if err = i.client.DeleteResourcesTypes(ctx, removedResourceTypeIDs); err != nil {
		blog.ErrorJSON("delete resource types failed, error: %s, resource types: %s", err.Error(), removedResourceTypeIDs)
		return err
	}
	return nil
}

var token string
var tokenRefreshTime time.Time

func (i Iam) CheckRequestAuthorization(req *http.Request) (bool, error) {
	name, pwd, ok := req.BasicAuth()
	if !ok || name != SystemIDIAM {
		blog.Errorf("request have no basic authorization")
		return false, nil
	}
	// if cached token is set within a minute, use it to check request authorization
	if token != "" && time.Since(tokenRefreshTime) <= time.Minute && pwd == token {
		return true, nil
	}
	var err error
	token, err = i.client.GetSystemToken(context.Background())
	if err != nil {
		blog.Errorf("check request authorization get system token failed, error: %s", err.Error())
		return false, err
	}
	tokenRefreshTime = time.Now()
	if pwd == token {
		return true, nil
	}
	blog.Errorf("request password not match system token")
	return false, nil
}
