/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package iam

import (
	"context"
	"net/http"

	"configcenter/src/ac/iam/types"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/thirdparty/apigw/apigwutil"
	"configcenter/src/thirdparty/apigw/apigwutil/user"
)

// ClientI is the iam api gateway client
type ClientI interface {
	GetNoAuthSkipUrl(ctx context.Context, header http.Header, p metadata.IamPermission) (string, error)
	RegisterResourceCreatorAction(ctx context.Context, header http.Header, instance metadata.IamInstanceWithCreator) (
		[]metadata.IamCreatorActionPolicy, error)
	BatchRegisterResourceCreatorAction(ctx context.Context, header http.Header,
		instance metadata.IamInstancesWithCreator) ([]metadata.IamCreatorActionPolicy, error)
	BatchOperateInstanceAuth(ctx context.Context, header http.Header, req *metadata.IamBatchOperateInstanceAuthReq) (
		[]metadata.IamBatchOperateInstanceAuthRes, error)

	// RegisterSystem register cmdb system in IAM, returns the registered system id
	RegisterSystem(ctx context.Context, header http.Header, sys *System) (string, error)
	// GetSystem get cmdb system info from IAM
	GetSystem(ctx context.Context, header http.Header) (*System, error)
	// UpdateSystem update cmdb system info in IAM
	UpdateSystem(ctx context.Context, header http.Header, sys *System) error
	// GetSystemToken get cmdb system auth token from IAM
	GetSystemToken(ctx context.Context, header http.Header) (string, error)

	// ListResourceTypes list resource types by page
	ListResourceTypes(ctx context.Context, header http.Header, page, pageSize int64) (*ListResourceTypesData, error)
	// RegisterResourcesTypes register resource types in IAM, returns the registered resource type ids
	RegisterResourcesTypes(ctx context.Context, header http.Header, resTypes []ResourceType) ([]string, error)
	// UpdateResourcesType update resource type in IAM
	UpdateResourcesType(ctx context.Context, header http.Header, resTypeID types.TypeID,
		req *UpdateResourceTypeReq) error
	// DeleteResourcesType delete resource type in IAM
	DeleteResourcesType(ctx context.Context, header http.Header, resTypeID types.TypeID) error

	// ListActions list actions by page
	ListActions(ctx context.Context, header http.Header, page, pageSize int64) (*ListActionsData, error)
	// RegisterActions register actions in IAM, returns the registered action ids
	RegisterActions(ctx context.Context, header http.Header, actions []ResourceAction) ([]string, error)
	// UpdateAction update action in IAM
	UpdateAction(ctx context.Context, header http.Header, actionID types.ActionID, req *UpdateActionReq) error
	// DeleteAction delete action in IAM
	DeleteAction(ctx context.Context, header http.Header, actionID types.ActionID) error

	// ListRoles list roles by page
	ListRoles(ctx context.Context, header http.Header, page, pageSize int64) (*ListRolesData, error)
	// RegisterRoles register roles in IAM, returns the registered role ids
	RegisterRoles(ctx context.Context, header http.Header, roles []Role) ([]string, error)
	// UpdateRole update role in IAM
	UpdateRole(ctx context.Context, header http.Header, roleID types.RoleID, req *UpdateRoleReq) error
	// DeleteRole delete role in IAM
	DeleteRole(ctx context.Context, header http.Header, roleID types.RoleID) error
	// AddRoleActions bind actions to a role, returns the added action ids
	AddRoleActions(ctx context.Context, header http.Header, roleID types.RoleID, actions []RoleAction) ([]string, error)
	// DeleteRoleActions unbind actions from a role
	DeleteRoleActions(ctx context.Context, header http.Header, roleID types.RoleID, actionIDs []types.ActionID) error

	// DeleteActionPolicies delete action policies in IAM
	DeleteActionPolicies(ctx context.Context, header http.Header, actionID types.ActionID) error
	// ListPolicies list action policies in IAM
	ListPolicies(ctx context.Context, header http.Header, params *ListPoliciesParams) (*ListPoliciesData, error)

	ListUserPolicies(ctx context.Context, header http.Header, opts *ListPolicyOptions) ([]*ActionPolicy, error)
	GetUserPolicy(ctx context.Context, header http.Header, opt *GetPolicyOption) (*operator.Policy, error)
}

type iam struct {
	service *apigwutil.ApiGWSrv
	userCli user.VirtualUserClientI
}

// NewClient create gse api gateway client
func NewClient(options *apigwutil.ApiGWOptions, userCli user.VirtualUserClientI) (ClientI, error) {
	service, err := apigwutil.NewApiGW(options, "apiGW.bkIamApiGatewayUrl")
	if err != nil {
		return nil, err
	}

	return &iam{
		service: service,
		userCli: userCli,
	}, nil
}
