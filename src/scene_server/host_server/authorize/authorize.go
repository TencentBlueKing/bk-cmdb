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

	restful "github.com/emicklei/go-restful"

	"configcenter/src/auth"
	"configcenter/src/auth/meta"
	"configcenter/src/auth/parser"
	"configcenter/src/common/blog"
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
func NewIamAuthorizeData(commonInfo *meta.CommonInfo, businessID int64,
	resourceType meta.ResourceType, instanceIDs *[]int64, action meta.Action) *meta.AuthAttribute {
	iamAuthorizeRequestBody := &meta.AuthAttribute{
		APIVersion: commonInfo.APIVersion,
		BusinessID: businessID,
		User:       commonInfo.User,
	}

	for _, instanceID := range *instanceIDs {
		resource := NewResoruceAttribute(commonInfo, businessID, resourceType, instanceID, action)
		iamAuthorizeRequestBody.Resources = append(iamAuthorizeRequestBody.Resources, *resource)
	}
	return iamAuthorizeRequestBody
}

// NewResoruceAttribute new a resource attribute
func NewResoruceAttribute(commonInfo *meta.CommonInfo, businessID int64,
	resourceType meta.ResourceType, instanceID int64, action meta.Action) *meta.ResourceAttribute {
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

// CanDoBusinessAction check permission for operate business
func (ha *HostAuthorizer) CanDoBusinessAction(req *restful.Request, businessID int64, action meta.Action) (decision meta.Decision, err error) {

	commonInfo, err := parser.ParseCommonInfo(req)
	if err != nil {
		decision := meta.Decision{
			Authorized: false,
			Reason:     fmt.Sprintf("parse common info from request failed, %s", err),
		}
		return decision, err
	}
	iamAuthorizeRequestBody := NewIamAuthorizeData(commonInfo, businessID, meta.Business, &[]int64{businessID}, action)

	decision, err = ha.authorizer.Authorize(context.Background(), iamAuthorizeRequestBody)
	if err != nil {
		message := fmt.Sprintf("auth interface failed, %s", err)
		blog.Errorf(message)

		decision = meta.Decision{
			Authorized: false,
			Reason:     message,
		}
		return decision, err
	}
	return
}

// CanDoHostAction check whether have permission to view host snapshot
func (ha *HostAuthorizer) CanDoHostAction(req *restful.Request, businessID int64, hostIDs *[]int64, action meta.Action) (decision meta.Decision, err error) {
	commonInfo, err := parser.ParseCommonInfo(req)
	if err != nil {
		decision := meta.Decision{
			Authorized: false,
			Reason:     fmt.Sprintf("parse common info from request failed, %s", err),
		}
		return decision, nil
	}

	iamAuthorizeRequestBody := NewIamAuthorizeData(commonInfo, businessID, meta.Host, hostIDs, action)

	decision, err = ha.authorizer.Authorize(context.Background(), iamAuthorizeRequestBody)
	if err != nil {
		message := fmt.Sprintf("auth interface failed, %s", err)
		blog.Errorf(message)

		decision = meta.Decision{
			Authorized: false,
			Reason:     message,
		}
		return decision, err
	}
	return
}

// RegisterHosts register host to auth center
func (ha *HostAuthorizer) RegisterHosts(req *restful.Request, businessID int64, hostIDs *[]int64) error {
	// TODO make it atomic
	commonInfo, err := parser.ParseCommonInfo(req)
	if err != nil {
		return fmt.Errorf("parse common info from request failed, %s", err)
	}

	/*
		// TODO what action should i use for register resource
		// batch register hosts
		iamAuthorizeRequestBody := NewIamAuthorizeData(commonInfo, businessID, meta.Host, hostIDs, "")
		err = ha.register.Register(context.Background(), iamAuthorizeRequestBody)
		if err != nil {
			blog.Errorf("auth register hosts failed, iamAuthorizeRequestBody=%v, error: %s", iamAuthorizeRequestBody, err)
		}
	*/

	// register host one by one
	for _, hostID := range *hostIDs {
		resource := NewResoruceAttribute(commonInfo, businessID, meta.Host, hostID, meta.EmptyAction)
		err := ha.register.RegisterResource(context.Background(), *resource)
		if err == nil {
			blog.Debug("auth register success, resourceAttribute=%v", resource)
			continue
		}

		blog.Errorf("auth register failed, resourceAttribute=%v, error: %s", resource, err)
		return err
	}
	return nil
}

// DeregisterHosts register host to auth center
func (ha *HostAuthorizer) DeregisterHosts(req *restful.Request, businessID int64, hostIDs *[]int64) error {
	commonInfo, err := parser.ParseCommonInfo(req)
	if err != nil {
		return fmt.Errorf("parse common info from request failed, %s", err)
	}

	for _, hostID := range *hostIDs {
		resource := NewResoruceAttribute(commonInfo, businessID, meta.Host, hostID, meta.EmptyAction)
		err := ha.register.RegisterResource(context.Background(), *resource)
		if err == nil {
			blog.Debug("auth register success, resourceAttribute=%v", resource)
			continue
		}

		blog.Errorf("auth register failed, resourceAttribute=%v, error: %s", resource, err)
		return err
	}
	return nil
}
