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

package authorize

import (
	"configcenter/src/auth"
	"configcenter/src/auth/meta"
	"configcenter/src/auth/parser"
	"configcenter/src/common/blog"
	"context"
	"fmt"

	restful "github.com/emicklei/go-restful"
)

// HostAuthorizer manages authorize interface for host server
type HostAuthorizer struct {
	authorizer auth.Authorizer
	register   auth.ResourceHandler
}

// NewHostAuthorizer new authorizer for host server
func NewHostAuthorizer() *HostAuthorizer {
	authorizer := new(HostAuthorizer)
	return authorizer
}

// NewIamAuthorizeData new a meta.Attribute object
func NewIamAuthorizeData(apiVersion string, businessID int64, userInfo meta.UserInfo, resourceType meta.ResourceType, instanceIDs *[]int64, action meta.Action) *meta.Attribute {
	attribute := &meta.Attribute{
		APIVersion: apiVersion,
		BusinessID: businessID,
		User:       userInfo,
	}

	for _, instanceID := range instanceIDs {
		resource := meta.Resource{
			Type: resourceType,
			Name: "",
		}
		Type ResourceType
		Name string
		InstanceID int64
	}
	return attribute
}

// CanDoBusinessAction check permission for operate business
func (ha *HostAuthorizer) CanDoBusinessAction(req *restful.Request, businessID int64, action meta.Action) (decision meta.Decision, err error) {

	commonInfo, err := parser.ParseCommonInfo(req)
	if err != nil {
		return auth.DecisionDeny, fmt.Sprintf("parse common info from request failed, %s", err), err
	}

	iamAuthorizeRequestBody := meta.Attribute{
		APIVersion: commonInfo.APIVersion,
		BusinessID: businessID,
		User:       commonInfo.User,
	}

	resource := meta.Resource{
		Type:       meta.Business,
		Name:       fmt.Sprint("%s[%d]", meta.Business, businessID),
		InstanceID: businessID,
		Action:     action,
	}
	iamAuthorizeRequestBody.Resources = append(iamAuthorizeRequestBody.Resources, resource)

	decision, err = ha.authorizer.Authorize(&iamAuthorizeRequestBody)
	if err != nil {
		message := fmt.Sprintf("auth interface failed, %s", err)
		blog.Errorf(message)
		return auth.DecisionDeny, message, err
	}
	return
}

// CanDoHostAction check whether have permission to view host snapshot
func (ha *HostAuthorizer) CanDoHostAction(req *restful.Request, businessID int64, hostIDs *[]int64, action meta.Action) (decision meta.Decision, err error) {
	commonInfo, err := parser.ParseCommonInfo(req)
	if err != nil {
		return auth.DecisionDeny, fmt.Sprintf("parse common info from request failed, %s", err), err
	}

	iamAuthorizeRequestBody := meta.Attribute{
		APIVersion: commonInfo.APIVersion,
		BusinessID: businessID,
		User:       commonInfo.User,
	}

	resources, err := newResources(req, businessID, hostIDs, meta.Host, action)
	if err != nil {
		return auth.DecisionDeny, fmt.Sprintf("make host transfer authorize resources failed, %s", err), err
	}
	iamAuthorizeRequestBody.Resources = *resources

	authorized, reason, err = ha.authorizer.Authorize(&attribute)
	if err != nil {
		message := fmt.Sprintf("auth interface failed, %s", err)
		blog.Errorf(message)
		return auth.DecisionDeny, message, err
	}
	return
}

// RegisterHosts register host to auth center
func (ha *HostAuthorizer) RegisterHosts(req *restful.Request, businessID int64, hostIDs *[]int64) error {
	var item auth.Item
	// FIXME
	// 1. deal with auth interface failed
	// 2. make it atomic
	// FIXME replace string like host/set with constant
	for _, hostID := range *hostIDs {
		resourceAttribute := auth.ResourceAttribute{
			Object:     "host",
			ObjectName: "host",
		}
		item = auth.Item{
			Object:     "host",
			InstanceID: hostID,
		}
		resourceAttribute.Layers = append(resourceAttribute.Layers, item)

		item = auth.Item{
			Object:     "set",
			InstanceID: businessID,
		}
		resourceAttribute.Layers = append(resourceAttribute.Layers, item)

		requestID, err := ha.register.Register(context.Background(), &resourceAttribute)
		if err == nil {
			blog.Debug("auth register success, requestID=%s, resourceAttribute=%v", requestID, resourceAttribute)
			continue
		}

		message := fmt.Sprintf("auth register failed, requestID=%s, resourceAttribute=%v, error: %s", requestID, resourceAttribute, err)
		blog.Errorf(message)
		return err
	}
	return nil
}

// DeregisterHosts register host to auth center
func (ha *HostAuthorizer) DeregisterHosts(req *restful.Request, businessID int64, hostIDs *[]int64) error {
	var item auth.Item
	// FIXME
	// 1. deal with auth interface failed
	// 2. make it atomic
	// FIXME replace string like host/set with constant
	for _, hostID := range *hostIDs {
		resourceAttribute := auth.ResourceAttribute{
			Object:     "host",
			ObjectName: "host",
		}
		item = auth.Item{
			Object:     "host",
			InstanceID: hostID,
		}
		resourceAttribute.Layers = append(resourceAttribute.Layers, item)

		item = auth.Item{
			Object:     "set",
			InstanceID: businessID,
		}
		resourceAttribute.Layers = append(resourceAttribute.Layers, item)

		requestID, err := ha.register.Deregister(context.Background(), &resourceAttribute)
		if err == nil {
			blog.Debug("auth register success, requestID=%s, resourceAttribute=%v", requestID, resourceAttribute)
			continue
		}

		message := fmt.Sprintf("auth register failed, requestID=%s, resourceAttribute=%v, error: %s", requestID, resourceAttribute, err)
		blog.Errorf(message)
		return err
	}
	return nil
}
