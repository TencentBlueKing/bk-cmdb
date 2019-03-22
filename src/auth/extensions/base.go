/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package extensions

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/auth/meta"
	"configcenter/src/auth/parser"
)

func (am *AuthManager) authorize(ctx context.Context, header http.Header, businessID int64, resources ...meta.ResourceAttribute) error {
	commonInfo, err := parser.ParseCommonInfo(&header)
	if err != nil {
		return fmt.Errorf("authentication failed, parse user info from header failed, %+v", err)
	}
	authAttribute := &meta.AuthAttribute{
		User:       commonInfo.User,
		BusinessID: businessID,
		Resources:  resources,
	}

	decision, err := am.Authorize.Authorize(ctx, authAttribute)
	if err != nil {
		return fmt.Errorf("authorize failed, err: %+v", err)
	}
	if decision.Authorized == false {
		return fmt.Errorf("authorize failed, reason: %s", decision.Reason)
	}

	return nil
}

func (am *AuthManager) updateResources(ctx context.Context, resources ...meta.ResourceAttribute) error {
	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return nil
		}
	}
	return nil
}
