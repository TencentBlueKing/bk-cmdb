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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * Dynamic Group
 */

func (am *AuthManager) CollectDynamicGroupByBusinessID(ctx context.Context, header http.Header, businessID int64) ([]DynamicGroupSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKAppIDField).Eq(businessID).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameUserAPI, &cond)
	if err != nil {
		blog.V(3).Infof("get user api by business %d failed, err: %+v, rid: %s", businessID, err, rid)
		return nil, fmt.Errorf("get user api by business %d failed, err: %+v", businessID, err)
	}
	dynamicGroups := make([]DynamicGroupSimplify, 0)
	for _, item := range result.Data.Info {
		dynamicGroup := DynamicGroupSimplify{}
		_, err = dynamicGroup.Parse(item)
		if err != nil {
			return nil, fmt.Errorf("get user api by business %d failed, err: %+v", businessID, err)
		}
		dynamicGroups = append(dynamicGroups, dynamicGroup)
	}
	return dynamicGroups, nil
}

func (am *AuthManager) collectDynamicGroupByIDs(ctx context.Context, header http.Header, ids ...string) ([]DynamicGroupSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.StrArrayUnique(ids)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKFieldID).In(ids).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameUserAPI, &cond)
	if err != nil {
		blog.Errorf("get user api by id %+v failed, err: %+v, rid: %s", ids, err, rid)
		return nil, fmt.Errorf("get user api by id failed, err: %+v", err)
	}
	dynamicGroups := make([]DynamicGroupSimplify, 0)
	for _, item := range result.Data.Info {
		dynamicGroup := DynamicGroupSimplify{}
		_, err = dynamicGroup.Parse(item)
		if err != nil {
			blog.Errorf("collectDynamicGroupByIDs by id %+v failed, parse user api %+v failed, err: %+v, rid: %s", ids, item, err, rid)
			return nil, fmt.Errorf("parse user api from db data failed, err: %+v", err)
		}
		dynamicGroups = append(dynamicGroups, dynamicGroup)
	}
	return dynamicGroups, nil
}

func (am *AuthManager) MakeResourcesByDynamicGroups(header http.Header, action meta.Action, businessID int64, dynamicGroups ...DynamicGroupSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, dynamicGroup := range dynamicGroups {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:       action,
				Type:         meta.DynamicGrouping,
				Name:         dynamicGroup.Name,
				InstanceIDEx: dynamicGroup.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) extractBusinessIDFromDynamicGroups(dynamicGroups ...DynamicGroupSimplify) (int64, error) {
	var businessID int64
	for idx, dynamicGroup := range dynamicGroups {
		bizID := dynamicGroup.BKAppIDField
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("get multiple business ID from user apis")
		}
		businessID = bizID
	}
	return businessID, nil
}

// func (am *AuthManager) AuthorizeByDynamicGroups(ctx context.Context, header http.Header, action meta.Action, dynamicGroups ...DynamicGroupSimplify) error {
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	// extract business id
// 	bizID, err := am.extractBusinessIDFromDynamicGroups(dynamicGroups...)
// 	if err != nil {
// 		return fmt.Errorf("authorize user api failed, extract business id from user api failed, err: %+v", err)
// 	}
//
// 	// make auth resources
// 	resources := am.MakeResourcesByDynamicGroups(header, action, bizID, dynamicGroups...)
//
// 	return am.authorize(ctx, header, bizID, resources...)
// }

func (am *AuthManager) UpdateRegisteredDynamicGroups(ctx context.Context, header http.Header, dynamicGroups ...DynamicGroupSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(dynamicGroups) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromDynamicGroups(dynamicGroups...)
	if err != nil {
		return fmt.Errorf("authorize user api failed, extract business id from user api failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByDynamicGroups(header, meta.EmptyAction, bizID, dynamicGroups...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredDynamicGroupByID(ctx context.Context, header http.Header, ids ...string) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	dynamicGroups, err := am.collectDynamicGroupByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered dynamic group failed, get dynamic group by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredDynamicGroups(ctx, header, dynamicGroups...)
}

func (am *AuthManager) RegisterDynamicGroups(ctx context.Context, header http.Header, dynamicGroups ...DynamicGroupSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(dynamicGroups) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromDynamicGroups(dynamicGroups...)
	if err != nil {
		return fmt.Errorf("register dynamic group failed, extract business id from dynamic group failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByDynamicGroups(header, meta.EmptyAction, bizID, dynamicGroups...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) RegisterDynamicGroupByID(ctx context.Context, header http.Header, ids ...string) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}
	dynamicGroups, err := am.collectDynamicGroupByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered dynamic group failed, get dynamic group by id failed, err: %+v", err)
	}
	return am.RegisterDynamicGroups(ctx, header, dynamicGroups...)
}

func (am *AuthManager) DeregisterDynamicGroupByID(ctx context.Context, header http.Header, configMeta metadata.UserConfigMeta) error {
	if am.Enabled() == false {
		return nil
	}

	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Action:       meta.EmptyAction,
			Type:         meta.DynamicGrouping,
			Name:         configMeta.Name,
			InstanceIDEx: configMeta.ID,
		},
		SupplierAccount: util.GetOwnerID(header),
		BusinessID:      configMeta.AppID,
	}

	return am.Authorize.DeregisterResource(ctx, resource)
}
