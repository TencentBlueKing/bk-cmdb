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
	"strconv"

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * service categories
 */

func (am *AuthManager) CollectServiceCategoryByBusinessIDs(ctx context.Context, header http.Header, businessID int64) ([]metadata.ServiceCategory, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	option := metadata.ListServiceCategoriesOption{
		BusinessID:     businessID,
		WithStatistics: false,
	}
	result, err := am.clientSet.CoreService().Process().ListServiceCategories(ctx, header, option)
	if err != nil {
		blog.Errorf("list service categories by businessID:%d failed, err: %+v, rid: %s", businessID, err, rid)
		return nil, fmt.Errorf("list service categories by businessID:%d failed, err: %+v", businessID, err)
	}

	categories := make([]metadata.ServiceCategory, 0)
	for _, item := range result.Info {
		categories = append(categories, item.ServiceCategory)
	}
	return categories, nil
}

func (am *AuthManager) collectServiceCategoryByIDs(ctx context.Context, header http.Header, ids ...int64) ([]metadata.ServiceCategory, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.IntArrayUnique(ids)
	categories := make([]metadata.ServiceCategory, 0)
	for _, id := range ids {
		category, err := am.clientSet.CoreService().Process().GetServiceCategory(ctx, header, id)
		if err != nil {
			blog.V(3).Infof("get service categories by id failed, id: %d, err: %+v, rid: %s", id, err, rid)
			return nil, fmt.Errorf("list service categories by id failed, err: %+v", err)
		}
		categories = append(categories, *category)
	}

	return categories, nil
}

func (am *AuthManager) extractBusinessIDFromServiceCategory(categories ...metadata.ServiceCategory) (int64, error) {
	var businessID int64
	for idx, category := range categories {
		bizID := category.BizID
		// we should ignore metadata.LabelBusinessID field not found error
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("get multiple business ID from service category")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) MakeResourcesByServiceCategory(header http.Header, action meta.Action, businessID int64, categories ...metadata.ServiceCategory) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, category := range categories {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ProcessServiceCategory,
				Name:       category.Name,
				InstanceID: category.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) MakeResourcesByServiceCategoryIDs(ctx context.Context, header http.Header, action meta.Action, businessID int64, ids ...int64) ([]meta.ResourceAttribute, error) {
	categories, err := am.collectServiceCategoryByIDs(ctx, header, ids...)
	if err != nil {
		return nil, fmt.Errorf("get service categories by id failed, err: %+v", err)
	}
	resources := am.MakeResourcesByServiceCategory(header, action, businessID, categories...)
	return resources, nil
}

func (am *AuthManager) AuthorizeByServiceCategoryID(ctx context.Context, header http.Header, action meta.Action, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	categories, err := am.collectServiceCategoryByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("get service categories by id failed, err: %+v", err)
	}
	return am.AuthorizeByServiceCategory(ctx, header, action, categories...)
}

func (am *AuthManager) GenServiceCategoryNoPermissionResp() *metadata.BaseResp {
	permission := metadata.Permission{
		SystemID:      authcenter.SystemIDCMDB,
		SystemName:    authcenter.SystemNameCMDB,
		ScopeType:     authcenter.ScopeTypeIDSystem,
		ScopeTypeName: authcenter.ScopeTypeIDSystemName,
		ActionID:      string(authcenter.ModelTopologyOperation),
		ActionName:    authcenter.ActionIDNameMap[authcenter.ModelTopologyOperation],
		Resources: [][]metadata.Resource{
			{{
				ResourceType:     string(authcenter.SysSystemBase),
				ResourceTypeName: authcenter.ResourceTypeIDMap[authcenter.SysSystemBase],
			}},
		},
	}
	permission.ResourceType = permission.Resources[0][0].ResourceType
	permission.ResourceTypeName = permission.Resources[0][0].ResourceTypeName

	resp := metadata.NewNoPermissionResp([]metadata.Permission{permission})
	return &resp
}

func (am *AuthManager) AuthorizeByServiceCategory(ctx context.Context, header http.Header, action meta.Action, categories ...metadata.ServiceCategory) error {
	if am.Enabled() == false {
		return nil
	}

	if len(categories) == 0 {
		return nil
	}
	// extract business id
	bizID, err := am.extractBusinessIDFromServiceCategory(categories...)
	if err != nil {
		return fmt.Errorf("authorize service categories failed, extract business id from service categories failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByServiceCategory(header, action, bizID, categories...)

	return am.authorize(ctx, header, bizID, resources...)
}

func (am *AuthManager) UpdateRegisteredServiceCategory(ctx context.Context, header http.Header, categories ...metadata.ServiceCategory) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if am.Enabled() == false {
		return nil
	}

	if len(categories) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromServiceCategory(categories...)
	if err != nil {
		return fmt.Errorf("authorize service categories failed, extract business id from service categories failed, err: %+v, rid: %s", err, rid)
	}

	// make auth resources
	resources := am.MakeResourcesByServiceCategory(header, meta.EmptyAction, bizID, categories...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredServiceCategoryByID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	categories, err := am.collectServiceCategoryByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered service categories failed, get service categories by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredServiceCategory(ctx, header, categories...)
}

func (am *AuthManager) DeregisterServiceCategoryByIDs(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	categories, err := am.collectServiceCategoryByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("deregister service categories failed, get service categories by id failed, err: %+v", err)
	}
	return am.DeregisterServiceCategory(ctx, header, categories...)
}

func (am *AuthManager) RegisterServiceCategory(ctx context.Context, header http.Header, categories ...metadata.ServiceCategory) error {
	if am.Enabled() == false {
		return nil
	}

	if len(categories) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromServiceCategory(categories...)
	if err != nil {
		return fmt.Errorf("register service categories failed, extract business id from service categories failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByServiceCategory(header, meta.EmptyAction, bizID, categories...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) RegisterServiceCategoryByID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	categories, err := am.collectServiceCategoryByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("register service categories failed, get service categories by id failed, err: %+v", err)
	}
	return am.RegisterServiceCategory(ctx, header, categories...)
}

func (am *AuthManager) DeregisterServiceCategory(ctx context.Context, header http.Header, categories ...metadata.ServiceCategory) error {
	if am.Enabled() == false {
		return nil
	}

	if len(categories) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromServiceCategory(categories...)
	if err != nil {
		return fmt.Errorf("deregister service categories failed, extract business id from service categories failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByServiceCategory(header, meta.EmptyAction, bizID, categories...)

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) ListAuthorizedServiceCategoryIDs(ctx context.Context, header http.Header, bizID int64) ([]int64, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	listOption := &meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:   meta.ProcessServiceCategory,
			Action: meta.FindMany,
		},
		SupplierAccount: util.GetOwnerID(header),
		BusinessID:      bizID,
	}
	resources, err := am.Authorize.ListResources(ctx, listOption)
	if err != nil {
		blog.Errorf("list authorized service category from iam failed, err: %+v, rid: %s", err, rid)
		return nil, err
	}
	ids := make([]int64, 0)
	for _, item := range resources {
		for _, resource := range item {
			id, err := strconv.ParseInt(resource.ResourceID, 10, 64)
			if err != nil {
				blog.Errorf("list authorized service category from iam failed, err: %+v, rid: %s", err, rid)
				return nil, fmt.Errorf("parse resource id into int64 failed, err: %+v", err)
			}
			ids = append(ids, id)
		}
	}
	return ids, nil
}
