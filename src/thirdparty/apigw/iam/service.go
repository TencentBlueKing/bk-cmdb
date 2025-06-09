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

	RegisterSystem(ctx context.Context, header http.Header, sys System) error
	GetSystemInfo(ctx context.Context, header http.Header, fields []types.SystemQueryField) (*RegisteredSystemInfo,
		error)
	UpdateSystemConfig(ctx context.Context, header http.Header, config *SysConfig) error

	// RegisterResourcesTypes register resource types in IAM
	RegisterResourcesTypes(ctx context.Context, header http.Header, resTypes []ResourceType) error
	// UpdateResourcesType update resource type in IAM
	UpdateResourcesType(ctx context.Context, header http.Header, resType ResourceType) error

	// RegisterActions register actions in IAM
	RegisterActions(ctx context.Context, header http.Header, actions []ResourceAction) error
	// UpdateAction update action in IAM
	UpdateAction(ctx context.Context, header http.Header, action ResourceAction) error
	// DeleteActions delete actions in IAM
	DeleteActions(ctx context.Context, header http.Header, actionIDs []types.ActionID) error

	// RegisterActionGroups register action groups in IAM
	RegisterActionGroups(ctx context.Context, header http.Header, actionGroups []ActionGroup) error
	// UpdateActionGroups update action groups in IAM
	UpdateActionGroups(ctx context.Context, header http.Header, actionGroups []ActionGroup) error

	// RegisterInstanceSelections register instance selections in IAM
	RegisterInstanceSelections(ctx context.Context, header http.Header, instanceSelections []InstanceSelection) error
	// UpdateInstanceSelection update instance selection in IAM
	UpdateInstanceSelection(ctx context.Context, header http.Header, instanceSelection InstanceSelection) error
	// DeleteInstanceSelections delete instance selections in IAM
	DeleteInstanceSelections(ctx context.Context, header http.Header,
		instanceSelectionIDs []types.InstanceSelectionID) error

	// RegisterResourceCreatorActions regitser resource creator actions in IAM
	RegisterResourceCreatorActions(ctx context.Context, header http.Header,
		resourceCreatorActions ResourceCreatorActions) error
	// UpdateResourceCreatorActions update resource creator actions in IAM
	UpdateResourceCreatorActions(ctx context.Context, header http.Header,
		resourceCreatorActions ResourceCreatorActions) error
	// RegisterCommonActions register common actions in IAM
	RegisterCommonActions(ctx context.Context, header http.Header, commonActions []CommonAction) error
	// UpdateCommonActions update common actions in IAM
	UpdateCommonActions(ctx context.Context, header http.Header, commonActions []CommonAction) error
	// DeleteActionPolicies delete action policies in IAM
	DeleteActionPolicies(ctx context.Context, header http.Header, actionID types.ActionID) error
	// ListPolicies list action policies in IAM
	ListPolicies(ctx context.Context, header http.Header, params *ListPoliciesParams) (*ListPoliciesData, error)
	DeleteResourcesTypes(ctx context.Context, header http.Header, resTypeIDs []types.TypeID) error

	ListUserPolicies(ctx context.Context, header http.Header, opts *ListPolicyOptions) ([]*ActionPolicy, error)
	GetSystemToken(ctx context.Context, header http.Header) (string, error)
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
