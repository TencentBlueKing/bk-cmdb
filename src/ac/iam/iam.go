/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

// Package iam TODO
package iam

import (
	"context"
	"net/http"

	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/ac/meta"
	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/authserver"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	apigwcli "configcenter/src/common/resource/apigw"
	"configcenter/src/scene_server/auth_server/sdk/types"
	"configcenter/src/thirdparty/apigw/iam"
)

// IAM TODO
type IAM struct {
	Client iam.ClientI
}

// NewIAM new iam client
func NewIAM() (*IAM, error) {
	return &IAM{
		Client: apigwcli.Client().Iam(),
	}, nil
}

type authorizer struct {
	authClientSet authserver.AuthServerClientInterface
}

// NewAuthorizer new authorizer
func NewAuthorizer(clientSet apimachinery.ClientSetInterface) *authorizer {
	return &authorizer{authClientSet: clientSet.AuthServer()}
}

// AuthorizeBatch batch authorization will not pass if one of them does not have permission
func (a *authorizer) AuthorizeBatch(ctx context.Context, h http.Header, user meta.UserInfo,
	resources ...meta.ResourceAttribute) ([]types.Decision, error) {
	return a.authorizeBatch(ctx, h, true, user, resources...)
}

// AuthorizeAnyBatch batch authorization will pass if one of them has permission
func (a *authorizer) AuthorizeAnyBatch(ctx context.Context, h http.Header, user meta.UserInfo,
	resources ...meta.ResourceAttribute) ([]types.Decision, error) {
	return a.authorizeBatch(ctx, h, false, user, resources...)
}

func (a *authorizer) authorizeBatch(ctx context.Context, h http.Header, exact bool, user meta.UserInfo,
	resources ...meta.ResourceAttribute) ([]types.Decision, error) {

	rid := httpheader.GetRid(h)

	opts, decisions, authCntMap, err := parseAttributesToBatchOptions(rid, user, resources...)
	if err != nil {
		return nil, err
	}

	// all resources are skipped
	if opts == nil {
		return decisions, nil
	}

	if blog.V(5) {
		blog.InfoJSON("auth options: %s, rid: %s", opts, rid)
	}

	var authDecisions []types.Decision
	if exact {
		authDecisions, err = a.authClientSet.AuthorizeBatch(ctx, h, opts)
		if err != nil {
			blog.Errorf("authorize batch failed, err: %s, ops: %s, resources: %s, rid: %s", err, opts, resources, rid)
			return nil, err
		}
	} else {
		authDecisions, err = a.authClientSet.AuthorizeAnyBatch(ctx, h, opts)
		if err != nil {
			blog.Errorf("authorize any batch failed, err: %s, ops: %s, resources: %s, rid: %s", err, opts, resources,
				rid)
			return nil, err
		}

	}

	authIndex := 0
	for i := range decisions {
		// skip resources' decisions are already set as authorized
		if decisions[i].Authorized {
			continue
		}

		cnt := authCntMap[i]
		if cnt == 0 {
			continue
		}

		decisions[i].Authorized = true
		for j := 0; j < cnt; j++ {
			if !authDecisions[authIndex+j].Authorized {
				decisions[i].Authorized = false
				break
			}
		}
		authIndex += cnt
	}

	return decisions, nil
}

func parseAttributesToBatchOptions(rid string, user meta.UserInfo, resources ...meta.ResourceAttribute) (
	*iam.AuthBatchOptions, []types.Decision, map[int]int, error) {

	if !auth.EnableAuthorize() {
		decisions := make([]types.Decision, len(resources))
		for i := range decisions {
			decisions[i].Authorized = true
		}
		return nil, decisions, nil, nil
	}

	authBatchArr := make([]*iam.AuthBatch, 0)
	decisions := make([]types.Decision, len(resources))
	authCntMap := make(map[int]int)
	for index, resource := range resources {

		// this resource should be skipped, do not need to verify in auth center.
		if resource.Action == meta.SkipAction {
			decisions[index].Authorized = true
			blog.V(5).Infof("skip authorization for resource: %+v, rid: %s", resource, rid)
			continue
		}

		authOpts, err := AdaptAuthOptions(&resource)
		if err != nil {
			blog.Errorf("adaptor cmdb resource to iam failed, err: %s, rid: %s", err, rid)
			return nil, nil, nil, err
		}

		if len(authOpts) == 1 && authOpts[0].Action == iamtypes.Skip {
			decisions[index].Authorized = true
			blog.V(5).Infof("skip authorization for resource: %+v, rid: %s", resource, rid)
			continue
		}

		for _, opt := range authOpts {
			authBatchArr = append(authBatchArr, &iam.AuthBatch{
				Action:    iam.Action{ID: string(opt.Action)},
				Resources: opt.Resources,
			})
		}
		authCntMap[index] = len(authOpts)
	}

	// all resources are skipped
	if len(authBatchArr) == 0 {
		return nil, decisions, authCntMap, nil
	}

	ops := &iam.AuthBatchOptions{
		System: iamtypes.SystemIDCMDB,
		Subject: iam.Subject{
			Type: "user",
			ID:   user.UserName,
		},
		Batch: authBatchArr,
	}
	return ops, decisions, authCntMap, nil
}

// ListAuthorizedResources 获取用户有的资源id权限列表
func (a *authorizer) ListAuthorizedResources(ctx context.Context, h http.Header,
	input meta.ListAuthorizedResourcesParam) (*types.AuthorizeList, error) {
	return a.authClientSet.ListAuthorizedResources(ctx, h, input)
}

// GetNoAuthSkipUrl get no auth skip url
func (a *authorizer) GetNoAuthSkipUrl(ctx context.Context, h http.Header,
	input *metadata.IamPermission) (string, error) {
	return a.authClientSet.GetNoAuthSkipUrl(ctx, h, input)
}

// GetPermissionToApply get permission to apply
func (a *authorizer) GetPermissionToApply(ctx context.Context, h http.Header,
	input []meta.ResourceAttribute) (*metadata.IamPermission, error) {
	return a.authClientSet.GetPermissionToApply(ctx, h, input)
}

// RegisterResourceCreatorAction register resourceCreator Action
func (a *authorizer) RegisterResourceCreatorAction(ctx context.Context, h http.Header,
	input metadata.IamInstanceWithCreator) (
	[]metadata.IamCreatorActionPolicy, error) {

	return a.authClientSet.RegisterResourceCreatorAction(ctx, h, input)
}

// BatchRegisterResourceCreatorAction batch register resourceCreator action
func (a *authorizer) BatchRegisterResourceCreatorAction(ctx context.Context, h http.Header,
	input metadata.IamInstancesWithCreator) (
	[]metadata.IamCreatorActionPolicy, error) {

	return a.authClientSet.BatchRegisterResourceCreatorAction(ctx, h, input)
}
