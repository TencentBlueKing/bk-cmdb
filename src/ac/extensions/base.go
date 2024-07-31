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

// Package extensions defines extensional utilities for auth
package extensions

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"configcenter/src/ac"
	"configcenter/src/ac/meta"
	"configcenter/src/ac/parser"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// this variable is used to accelerate the way to check if a business is resource pool
// business or not.
var resourcePoolBusinessID int64

// Authorize cc auth resource, returns no permission response(only when not authorized) and if user is authorized
func (am *AuthManager) Authorize(kit *rest.Kit, resources ...meta.ResourceAttribute) (
	*metadata.BaseResp, bool) {

	if !am.Enabled() {
		return nil, true
	}

	blog.V(5).Infof("start authorize resources: %+v, rid: %s", resources, kit.Rid)

	// authorize all auth resources
	user := meta.UserInfo{
		UserName:        kit.User,
		SupplierAccount: kit.SupplierAccount,
	}
	decisions, err := am.Authorizer.AuthorizeBatch(kit.Ctx, kit.Header, user, resources...)
	if err != nil {
		blog.Errorf("authorize failed, resources: %+v, err: %v, rid: %s", resources, err, kit.Rid)
		return &metadata.BaseResp{
			Code:   common.CCErrCommCheckAuthorizeFailed,
			ErrMsg: kit.CCError.Error(common.CCErrCommCheckAuthorizeFailed).Error(),
			Result: false,
		}, false
	}

	authorized := true
	permissionRes := make([]meta.ResourceAttribute, 0)
	for idx, decision := range decisions {
		if decision.Authorized {
			continue
		}

		permissionRes = append(permissionRes, resources[idx])
		authorized = false
	}

	if authorized {
		return nil, true
	}

	// get permissions that user need to apply for this request
	permission, err := am.Authorizer.GetPermissionToApply(kit.Ctx, kit.Header, permissionRes)
	if err != nil {
		blog.Errorf("get permission to apply failed, resources: %+v, err: %v, rid: %s", resources, err, kit.Rid)
		return &metadata.BaseResp{
			Code:   common.CCErrCommCheckAuthorizeFailed,
			ErrMsg: kit.CCError.Error(common.CCErrCommCheckAuthorizeFailed).Error(),
			Result: false,
		}, false
	}

	blog.Errorf("request is not authorized, need permission: %+v, err: %v, rid: %s", permission, err, kit.Rid)

	return &metadata.BaseResp{
		Code:        common.CCNoPermission,
		ErrMsg:      kit.CCError.Error(common.CCErrCommAuthNotHavePermission).Error(),
		Result:      false,
		Permissions: permission,
	}, false
}

// getResourcePoolBusinessID to get bizID of resource pool
// this function is concurrent safe.
func (am *AuthManager) getResourcePoolBusinessID(ctx context.Context, header http.Header) (int64, error) {

	rid := util.ExtractRequestIDFromContext(ctx)
	// this operation is concurrent safe
	if atomic.LoadInt64(&resourcePoolBusinessID) != 0 {
		// resource pool business id is already set, return directly.
		return atomic.LoadInt64(&resourcePoolBusinessID), nil
	}
	// get resource pool business id now.
	query := &metadata.QueryCondition{
		Fields: []string{common.BKAppIDField, common.BkSupplierAccount},
		Condition: map[string]interface{}{
			"default": 1,
		},
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("get biz by query failed, err: %v, rid: %s", err, rid)
		return 0, err
	}

	supplier := httpheader.GetSupplierAccount(header)
	for idx, biz := range result.Info {
		if supplier == biz[common.BkSupplierAccount].(string) {
			if !result.Info[idx].Exists(common.BKAppIDField) {
				// this can not be happen normally.
				return 0, fmt.Errorf("can not find resource pool business id")
			}
			bizID, err := result.Info[idx].Int64(common.BKAppIDField)
			if err != nil {
				return 0, fmt.Errorf("get resource pool biz id failed, err: %v", err)
			}
			// update resource pool business id immediately
			atomic.StoreInt64(&resourcePoolBusinessID, bizID)

			return bizID, nil
		}
	}
	return 0, fmt.Errorf("get resource pool biz id failed, err: %v", err)

}

func (am *AuthManager) batchAuthorize(ctx context.Context, header http.Header,
	resources ...meta.ResourceAttribute) error {
	commonInfo, err := parser.ParseCommonInfo(header)
	if err != nil {
		return fmt.Errorf("authentication failed, parse user info from header failed, err: %+v", err)
	}
	decisions, err := am.Authorizer.AuthorizeBatch(ctx, header, commonInfo.User, resources...)
	if err != nil {
		return fmt.Errorf("authorize failed, err: %+v", err)
	}

	for _, decision := range decisions {
		if !decision.Authorized {
			return ac.NoAuthorizeError
		}
	}

	return nil
}

// Enabled returns if authorization is enabled
func (am *AuthManager) Enabled() bool {
	return auth.EnableAuthorize()
}
