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

	"configcenter/src/ac"
	"configcenter/src/ac/meta"
	"configcenter/src/ac/parser"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
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
		Fields: []string{common.BKAppIDField, common.BkSupplierAccount},
		Condition: map[string]interface{}{
			"default": 1,
		},
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDApp, query)
	if err != nil {
		return 0, err
	}

	if !result.Result {
		return 0, errors.New(result.ErrMsg)
	}

	supplier := util.GetOwnerID(header)
	for idx, biz := range result.Data.Info {
		if supplier == biz[common.BkSupplierAccount].(string) {
			if !result.Data.Info[idx].Exists(common.BKAppIDField) {
				// this can not be happen normally.
				return 0, fmt.Errorf("can not find resource pool business id")
			}
			bizID, err := result.Data.Info[idx].Int64(common.BKAppIDField)
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

func (am *AuthManager) batchAuthorize(ctx context.Context, header http.Header, resources ...meta.ResourceAttribute) error {
	commonInfo, err := parser.ParseCommonInfo(&header)
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

func (am *AuthManager) Enabled() bool {
	return auth.EnableAuthorize()
}
