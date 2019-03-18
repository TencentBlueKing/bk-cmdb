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
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/util"
	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/meta"
	"configcenter/src/auth/parser"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
)

// HostAuthorizer manages authorize interface for host server
type HostAuthorizer struct {
	authorizer auth.Authorizer
	register   auth.ResourceHandler
}

// NewHostAuthorizer new authorizer for host server
func NewHostAuthorizer(tls *util.TLSClientConfig, optionConfig authcenter.AuthConfig) (*HostAuthorizer, error) {
	config := authcenter.AuthConfig{
		Address:   optionConfig.Address,
		AppCode:   optionConfig.AppCode,
		AppSecret: optionConfig.AppSecret,
		SystemID:  authcenter.SystemIDCMDB,
		Enable:    optionConfig.Enable,
	}
	authAuthorizer, err := auth.NewAuthorize(tls, config)
	if err != nil {
		blog.Errorf("new host authorizer failed, err: %+v", err)
		return nil, fmt.Errorf("new host authorizer failed, err: %+v", err)
	}
	authRegister, err := auth.NewAuthorize(tls, config)
	if err != nil {
		blog.Errorf("new host authorizer failed, err: %+v", err)
		return nil, fmt.Errorf("new host authorizer failed, err: %+v", err)
	}
	authorizer := &HostAuthorizer{
		authorizer: authAuthorizer,
		register:   authRegister,
	}
	return authorizer, nil
}

// NewIamAuthorizeData new a meta.Attribute object
func NewIamAuthorizeData(commonInfo *meta.CommonInfo, businessID int64,
	resourceType meta.ResourceType, instanceIDs *[]int64, action meta.Action) *meta.AuthAttribute {

	iamAuthorizeRequestBody := &meta.AuthAttribute{
		BusinessID: businessID,
		User:       commonInfo.User,
		Resources:  make([]meta.ResourceAttribute, 0),
	}

	for _, instanceID := range *instanceIDs {
		resource := NewResourceAttribute(commonInfo, businessID, resourceType, instanceID, action)
		iamAuthorizeRequestBody.Resources = append(iamAuthorizeRequestBody.Resources, *resource)
	}
	return iamAuthorizeRequestBody
}

// NewIamAuthorizeData new a meta.Attribute object
func NewBatchResourceAttributeWithLayers(commonInfo *meta.CommonInfo, businessID int64,
	resourceType meta.ResourceType, layers [][]meta.Item, action meta.Action) *meta.AuthAttribute {

	iamAuthorizeRequestBody := &meta.AuthAttribute{
		BusinessID: businessID,
		User:       commonInfo.User,
		Resources:  make([]meta.ResourceAttribute, 0),
	}

	for _, layer := range layers {
		resource := NewResourceAttributeWithLayers(commonInfo, businessID, resourceType, layer, action)
		iamAuthorizeRequestBody.Resources = append(iamAuthorizeRequestBody.Resources, *resource)
	}
	return iamAuthorizeRequestBody
}

// NewResourceAttribute new a resource attribute
func NewResourceAttributeWithLayers(commonInfo *meta.CommonInfo, businessID int64, resourceType meta.ResourceType, layer []meta.Item, action meta.Action) *meta.ResourceAttribute {
	resource := &meta.ResourceAttribute{
		Layers:          layer,
		BusinessID:      businessID,
		SupplierAccount: commonInfo.User.SupplierAccount,
	}

	if len(layer) > 0 {
		instance := layer[len(layer)-1]
		basic := meta.Basic{
			Type:       resourceType,
			Action:     action,
			Name:       instance.Name,
			InstanceID: instance.InstanceID,
		}
		resource.Basic = basic
	}
	return resource
}

// NewResourceAttribute new a resource attribute
func NewResourceAttribute(commonInfo *meta.CommonInfo, businessID int64, resourceType meta.ResourceType, instanceID int64, action meta.Action) *meta.ResourceAttribute {
	resource := &meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       resourceType,
			Name:       fmt.Sprintf("%s[%d]", resourceType, instanceID),
			InstanceID: instanceID,
			Action:     action,
		},
		BusinessID:      businessID,
		SupplierAccount: commonInfo.User.SupplierAccount,
	}
	return resource
}

// canDoResourceAction check permission for operate business
func (ha *HostAuthorizer) canDoResourceAction(requestHeader *http.Header, resourceType meta.ResourceType,
	businessID int64, instanceIDs *[]int64, action meta.Action) (decision meta.Decision, err error) {

	commonInfo, err := parser.ParseCommonInfo(requestHeader)
	if err != nil {
		decision := meta.Decision{
			Authorized: false,
			Reason:     fmt.Sprintf("parse common info from request failed, %s", err),
		}
		return decision, err
	}
	iamAuthorizeRequestBody := NewIamAuthorizeData(commonInfo, businessID, resourceType, &[]int64{businessID}, action)

	decision, err = ha.authorizer.Authorize(context.Background(), iamAuthorizeRequestBody)
	if err != nil {
		blog.Errorf("auth interface failed, %s", err)

		decision = meta.Decision{
			Authorized: false,
			Reason:     fmt.Sprintf("auth interface failed, %s", err),
		}
		return decision, fmt.Errorf("auth interface failed, %s", err)
	}
	return
}

// canDoResourceAction check permission for operate business
func (ha *HostAuthorizer) CanDoResourceActionWithLayers(requestHeader *http.Header, resourceType meta.ResourceType,
	businessID int64, layers [][]meta.Item, action meta.Action) (decision meta.Decision, err error) {

	commonInfo, err := parser.ParseCommonInfo(requestHeader)
	if err != nil {
		decision := meta.Decision{
			Authorized: false,
			Reason:     fmt.Sprintf("parse common info from request failed, %s", err),
		}
		return decision, err
	}
	iamAuthorizeRequestBody := NewBatchResourceAttributeWithLayers(commonInfo, businessID, resourceType, layers, action)

	decision, err = ha.authorizer.Authorize(context.Background(), iamAuthorizeRequestBody)
	if err != nil {
		blog.Errorf("auth interface failed, %+v", err)

		decision = meta.Decision{
			Authorized: false,
			Reason:     fmt.Sprintf("auth interface failed, %+v", err),
		}
		return decision, fmt.Errorf("auth interface failed, %+v", err)
	}
	return
}

// CanDoBusinessAction check permission for operate business
func (ha *HostAuthorizer) CanDoBusinessAction(requestHeader *http.Header, businessID int64,
	action meta.Action) (decision meta.Decision, err error) {

	return ha.canDoResourceAction(requestHeader, meta.Business, businessID, &[]int64{businessID}, action)
}

// CanDoModuleAction check permission for operate business
func (ha *HostAuthorizer) CanDoModuleAction(requestHeader *http.Header, businessID int64,
	moduleID int64, action meta.Action) (decision meta.Decision, err error) {

	return ha.canDoResourceAction(requestHeader, meta.ModelModule, businessID, &[]int64{moduleID}, action)
}

// RegisterHosts register host to auth center
func (ha *HostAuthorizer) RegisterHosts(requestHeader *http.Header, businessID int64, hostIDs *[]int64) error {
	return ha.registerResource(requestHeader, meta.HostInstance, businessID, hostIDs)
}

// registerResource register resource of resourceType type to auth center
func (ha *HostAuthorizer) registerResource(requestHeader *http.Header, resourceType meta.ResourceType,
	businessID int64, instanceIDs *[]int64) error {

	// TODO make it atomic
	commonInfo, err := parser.ParseCommonInfo(requestHeader)
	if err != nil {
		return fmt.Errorf("parse common info from request failed, %s", err)
	}

	resources := make([]meta.ResourceAttribute, 0)
	for _, instanceID := range *instanceIDs {
		resource := NewResourceAttribute(commonInfo, businessID, resourceType, instanceID, meta.EmptyAction)
		resources = append(resources, *resource)
	}

	if err := ha.register.RegisterResource(context.TODO(), resources...); err != nil {
		blog.Errorf("auth register failed, resourceAttribute: %+v, error: %s", resources, err)
		return fmt.Errorf("auth register failed, resourceAttribute: %+v, error: %s", resources, err)
	}
	return nil
}

// RegisterResourceWithLayers register resource of resourceType type to auth center
func (ha *HostAuthorizer) RegisterResourceWithLayers(requestHeader *http.Header, resourceType meta.ResourceType,
	businessID int64, layers *[][]meta.Item) error {

	// TODO make it atomic
	commonInfo, err := parser.ParseCommonInfo(requestHeader)
	if err != nil {
		return fmt.Errorf("parse common info from request failed, %s", err)
	}

	resources := make([]meta.ResourceAttribute, 0)
	for _, layer := range *layers {
		resource := NewResourceAttributeWithLayers(commonInfo, businessID, resourceType, layer, meta.EmptyAction)
		resources = append(resources, *resource)
	}

	resourcesData, err := json.Marshal(resources)
	if err != nil {
		blog.Errorf("auth register failed, resourceAttribute: %+v, error: %s", resources, err)
		return fmt.Errorf("auth register failed, resourceAttribute: %+v, error: %s", resources, err)
	}
	blog.Infof("auth register data: %s", resourcesData)
	if err := ha.register.RegisterResource(context.TODO(), resources...); err != nil {
		blog.Errorf("auth register failed, resourceAttribute: %+v, error: %s", resources, err)
		return fmt.Errorf("auth register failed, resourceAttribute: %+v, error: %s", resources, err)
	}
	return nil
}

// deregisterResource register resource of resourceType type to auth center
func (ha *HostAuthorizer) deregisterResource(requestHeader *http.Header, resourceType meta.ResourceType,
	businessID int64, instanceIDs *[]int64) error {

	// TODO make it atomic
	commonInfo, err := parser.ParseCommonInfo(requestHeader)
	if err != nil {
		return fmt.Errorf("parse common info from request failed, %s", err)
	}

	resources := make([]meta.ResourceAttribute, 0)
	for _, instanceID := range *instanceIDs {
		resource := NewResourceAttribute(commonInfo, businessID, resourceType, instanceID, meta.EmptyAction)
		resources = append(resources, *resource)
	}

	if err := ha.register.DeregisterResource(context.TODO(), resources...); err != nil {
		blog.Errorf("auth deregister failed, resourceAttribute: %+v, error: %s", resources, err)
		return fmt.Errorf("auth deregister failed, resourceAttribute: %+v, error: %s", resources, err)
	}
	return nil
}

// DeregisterHosts register host to auth center
func (ha *HostAuthorizer) DeregisterHosts(requestHeader *http.Header, businessID int64, hostIDs *[]int64) error {
	return ha.deregisterResource(requestHeader, meta.Host, businessID, hostIDs)
}

// DeregisterHosts register host to auth center
func (ha *HostAuthorizer) DeregisterResourceWithLayers(requestHeader *http.Header, resourceType meta.ResourceType,
	businessID int64, layers *[][]meta.Item) error {

	commonInfo, err := parser.ParseCommonInfo(requestHeader)
	if err != nil {
		return fmt.Errorf("parse common info from request failed, %+v", err)
	}

	resources := make([]meta.ResourceAttribute, 0)
	for _, layer := range *layers {
		resource := NewResourceAttributeWithLayers(commonInfo, businessID, resourceType, layer, meta.EmptyAction)
		resources = append(resources, *resource)
	}

	resourcesData, err := json.Marshal(resources)
	if err != nil {
		return fmt.Errorf("auth register failed, resourceAttribute: %+v, error: %s", resources, err)
	}
	blog.Infof("auth register data: %s", resourcesData)
	if err := ha.register.DeregisterResource(context.TODO(), resources...); err != nil {
		return fmt.Errorf("auth deregister failed, resourceAttribute: %+v, error: %s", resources, err)
	}
	return nil
}
