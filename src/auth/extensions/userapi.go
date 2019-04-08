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
	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"context"
	"fmt"
	"net/http"
)

/*
 * user api
 */

func (am *AuthManager) CollectUserAPIByBusinessID(ctx context.Context, header http.Header, businessID int64) ([]UserAPISimplify, error) {
	cond := metadata.QueryCondition{
		Fields:    []string{common.BKAppIDField, common.BKFieldID, common.BKFieldName},
		Condition: condition.CreateCondition().Field(common.BKAppIDField).Eq(businessID).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameUserAPI, &cond)
	if err != nil {
		blog.V(3).Infof("get user api by business %d failed, err: %+v", businessID, err)
		return nil, fmt.Errorf("get user api by business %d failed, err: %+v", businessID, err)
	}
	userAPIs := make([]UserAPISimplify, 0)
	for _, item := range result.Data.Info {
		userAPI := UserAPISimplify{}
		_, err = userAPI.Parse(item)
		if err != nil {
			return nil, fmt.Errorf("get user api by business %d failed, err: %+v", businessID, err)
		}
		userAPIs = append(userAPIs, userAPI)
	}
	return userAPIs, nil
}

func (am *AuthManager) collectUserAPIByIDs(ctx context.Context, header http.Header, ids ...string) ([]UserAPISimplify, error) {
	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.StrArrayUnique(ids)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKFieldID).In(ids).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameUserAPI, &cond)
	if err != nil {
		blog.Errorf("get user api by id %+v failed, err: %+v", ids, err)
		return nil, fmt.Errorf("get user api by id failed, err: %+v", err)
	}
	userAPIs := make([]UserAPISimplify, 0)
	for _, item := range result.Data.Info {
		userAPI := UserAPISimplify{}
		_, err = userAPI.Parse(item)
		if err != nil {
			blog.Errorf("collectUserAPIByIDs by id %+v failed, parse user api %+v failed, err: %+v ", ids, item, err)
			return nil, fmt.Errorf("parse user api from db data failed, err: %+v", err)
		}
		userAPIs = append(userAPIs, userAPI)
	}
	return userAPIs, nil
}

func (am *AuthManager) MakeResourcesByUserAPIs(header http.Header, action meta.Action, businessID int64, userAPIs ...UserAPISimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, userAPI := range userAPIs {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.DynamicGrouping,
				Name:       userAPI.Name,
				InstanceIDEx: userAPI.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) extractBusinessIDFromUserAPIs(userAPIs ...UserAPISimplify) (int64, error) {
	var businessID int64
	for idx, userAPI := range userAPIs {
		bizID := userAPI.BKAppIDField
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("get multiple business ID from user apis")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) AuthorizeByUserAPIs(ctx context.Context, header http.Header, action meta.Action, userAPIs ...UserAPISimplify) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromUserAPIs(userAPIs...)
	if err != nil {
		return fmt.Errorf("authorize user api failed, extract business id from user api failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByUserAPIs(header, action, bizID, userAPIs...)

	return am.authorize(ctx, header, bizID, resources...)
}

func (am *AuthManager) UpdateRegisteredUserAPIs(ctx context.Context, header http.Header, userAPIs ...UserAPISimplify) error {
	// extract business id
	bizID, err := am.extractBusinessIDFromUserAPIs(userAPIs...)
	if err != nil {
		return fmt.Errorf("authorize user api failed, extract business id from user api failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByUserAPIs(header, meta.EmptyAction, bizID, userAPIs...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredUserAPIByID(ctx context.Context, header http.Header, ids ...string) error {
	userAPIs, err := am.collectUserAPIByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered user apis failed, get user api by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredUserAPIs(ctx, header, userAPIs...)
}

func (am *AuthManager) RegisterUserAPIs(ctx context.Context, header http.Header, userAPIs ...UserAPISimplify) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromUserAPIs(userAPIs...)
	if err != nil {
		return fmt.Errorf("register user api failed, extract business id from user apis failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByUserAPIs(header, meta.EmptyAction, bizID, userAPIs...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) RegisterUserAPIByID(ctx context.Context, header http.Header, ids ...string) error {
	userAPIs, err := am.collectUserAPIByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered user apis failed, get user api by id failed, err: %+v", err)
	}
	return am.RegisterUserAPIs(ctx, header, userAPIs...)
}

func (am *AuthManager) DeregisterUserAPIs(ctx context.Context, header http.Header, userAPIs ...UserAPISimplify) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromUserAPIs(userAPIs...)
	if err != nil {
		return fmt.Errorf("deregister user api failed, extract business id from user api failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByUserAPIs(header, meta.EmptyAction, bizID, userAPIs...)

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterUserAPIByID(ctx context.Context, header http.Header, ids ...string) error {
	userAPIs, err := am.collectUserAPIByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered user apis failed, get user api by id failed, err: %+v", err)
	}
	return am.DeregisterUserAPIs(ctx, header, userAPIs...)
}

