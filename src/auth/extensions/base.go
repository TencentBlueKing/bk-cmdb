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
	"strings"

	"configcenter/src/auth/meta"
	"configcenter/src/auth/parser"
	"configcenter/src/common/metadata"
)

// correctBusinessID correct businessID to 0 if default field is 1, as we need to set it to 0 for iam.
// it's an ugly design here, but it's the cheapest way to set business to 0 for all default resources.
func (am *AuthManager) correctBusinessID(ctx context.Context, header http.Header, businessID int64) (int64, error) {

	businesses, err := am.collectBusinessByIDs(ctx, header, businessID)
	if err != nil {
		return 0, fmt.Errorf("get business:%d detailed failed, err: %+v", businessID, err)
	}
	if len(businesses) != 1 {
		return 0, fmt.Errorf("get business:%d failed, not found", businessID)
	}
	business := businesses[0]

	// set resource pool as global
	if business.IsDefault == 1 {
		return 0, nil
	}

	return business.BKAppIDField, nil
}

func (am *AuthManager) authorize(ctx context.Context, header http.Header, businessID int64, resources ...meta.ResourceAttribute) error {
	commonInfo, err := parser.ParseCommonInfo(&header)
	if err != nil {
		return fmt.Errorf("authentication failed, parse user info from header failed, %+v", err)
	}
	authAttribute := &meta.AuthAttribute{
		User:      commonInfo.User,
		Resources: resources,
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

func (am *AuthManager) batchAuthorize(ctx context.Context, header http.Header, resources ...meta.ResourceAttribute) error {
	commonInfo, err := parser.ParseCommonInfo(&header)
	if err != nil {
		return fmt.Errorf("authentication failed, parse user info from header failed, %+v", err)
	}

	decisions, err := am.Authorize.AuthorizeBatch(ctx, commonInfo.User, resources...)
	if err != nil {
		return fmt.Errorf("authorize failed, err: %+v", err)
	}

	messages := make([]string, 0)
	for _, decision := range decisions {
		if decision.Authorized == false {
			messages = append(messages, fmt.Sprintf("authorize failed, reason: %s", decision.Reason))
		}
	}
	if len(messages) > 0 {
		return fmt.Errorf(strings.Join(messages, "\n"))
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

// this functions works to parse business id from metadata's label
// if the business id key is exist in the label, then it will check
// the value and parse form it. otherwise it will return with 0.
func extractBusinessID(m metadata.Label) (int64, error) {
	if _, exist := m[metadata.LabelBusinessID]; exist {
		return m.Int64(metadata.LabelBusinessID)
	}
	return 0, nil
}
