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
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"

	"configcenter/src/auth"
	"configcenter/src/auth/meta"
	"configcenter/src/auth/parser"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

// this variable is used to accelerate the way to check if a business is resource pool
// business or not.
var resourcePoolBusinessID int64

// this function is concurrent safe.
func (am *AuthManager) getResourcePoolBusinessID(ctx context.Context, header http.Header) (int64, error) {

	// this operation is concurrent safe
	if atomic.LoadInt64(&resourcePoolBusinessID) != 0 {
		// resource pool business id is already set, return directly.
		return atomic.LoadInt64(&resourcePoolBusinessID), nil
	}
	// get resource pool business id now.
	query := &metadata.QueryCondition{
		Fields: []string{common.BKAppIDField},
		Condition: map[string]interface{}{
			"bk_biz_name": common.DefaultAppName,
			"default":     1,
		},
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDApp, query)
	if err != nil {
		return 0, err
	}

	if !result.Result {
		return 0, errors.New(result.ErrMsg)
	}

	if len(result.Data.Info) != 1 {
		// normally, this can not be happen.
		return 0, errors.New("get resource pool business id, but got multiple or not found")
	}

	// set resource pool as global
	if !result.Data.Info[0].Exists(common.BKAppIDField) {
		// this can not be happen normally.
		return 0, fmt.Errorf("can not find resource pool business id")
	}
	bizID, err := result.Data.Info[0].Int64(common.BKAppIDField)
	if err != nil {
		return 0, fmt.Errorf("get resource pool biz id failed, err: %v", err)
	}
	// update resource pool business id immediately
	atomic.StoreInt64(&resourcePoolBusinessID, bizID)

	return bizID, nil

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
		return auth.NoAuthorizeError
	}

	return nil
}

func (am *AuthManager) batchAuthorize(ctx context.Context, header http.Header, resources ...meta.ResourceAttribute) error {
	commonInfo, err := parser.ParseCommonInfo(&header)
	if err != nil {
		return fmt.Errorf("authentication failed, parse user info from header failed, err: %+v", err)
	}
	decisions, err := am.Authorize.AuthorizeBatch(ctx, commonInfo.User, resources...)
	if err != nil {
		return fmt.Errorf("authorize failed, err: %+v", err)
	}

	for _, decision := range decisions {
		if decision.Authorized == false {
			return auth.NoAuthorizeError
		}
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

func (am *AuthManager) Enabled() bool {
	if am == nil {
		return false
	}
	if am.Authorize == nil {
		return false
	}
	return am.Authorize.Enabled()
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

func (am *AuthManager) DeregisterResource(ctx context.Context, rs ...meta.ResourceAttribute) error {
	if am.Enabled() == false {
		return nil
	}
	return am.Authorize.DeregisterResource(ctx, rs...)
}
