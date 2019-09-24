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
	"strconv"

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * business related auth interface
 */

// CollectAllBusiness get all business
func (am *AuthManager) CollectAllBusiness(ctx context.Context, header http.Header) ([]BusinessSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	cond := metadata.QueryCondition{}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDApp, &cond)
	if err != nil {
		blog.Errorf("list business failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	// step1 get business from logics service
	businessList := make([]BusinessSimplify, 0)
	for _, business := range result.Data.Info {
		businessSimplify := BusinessSimplify{}
		_, err := businessSimplify.Parse(business)
		if err != nil {
			blog.Errorf("parse businesses %+v simplify information failed, err: %+v, rid: %s", business, err, rid)
			continue
		}
		if businessSimplify.IsDefault > 0 {
			continue
		}

		businessList = append(businessList, businessSimplify)
	}
	return businessList, nil
}

func (am *AuthManager) collectBusinessByIDs(ctx context.Context, header http.Header, businessIDs ...int64) ([]BusinessSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	businessIDs = util.IntArrayUnique(businessIDs)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKAppIDField).In(businessIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDApp, &cond)
	if err != nil {
		blog.V(3).Infof("get businesses by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get businesses by id failed, err: %+v", err)
	}
	blog.V(5).Infof("get businesses by id result: %+v", result)
	instances := make([]BusinessSimplify, 0)
	for _, cls := range result.Data.Info {
		instance := BusinessSimplify{}
		_, err = instance.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("parse business from db data failed, err: %+v", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (am *AuthManager) MakeResourcesByBusiness(header http.Header, action meta.Action, businesses ...BusinessSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, business := range businesses {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Business,
				Name:       business.BKAppNameField,
				InstanceID: business.BKAppIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) extractBusinessIDFromBusinesses(businesses ...BusinessSimplify) (int64, error) {
	var bizID int64
	for idx, business := range businesses {
		if idx != 0 && business.BKAppIDField != bizID {
			return 0, fmt.Errorf("get multiple business id[%d:%d] from businesses", bizID, business.BKAppIDField)
		}
		bizID = business.BKAppIDField
	}
	return bizID, nil
}

func (am *AuthManager) AuthorizeByBusiness(ctx context.Context, header http.Header, action meta.Action, businesses ...BusinessSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromBusinesses(businesses...)
	if err != nil {
		return fmt.Errorf("authorize instances failed, extract business id from instance failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByBusiness(header, action, businesses...)

	return am.authorize(ctx, header, bizID, resources...)
}

func (am *AuthManager) AuthorizeByBusinessID(ctx context.Context, header http.Header, action meta.Action, businessIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	businesses, err := am.collectBusinessByIDs(ctx, header, businessIDs...)
	if err != nil {
		return fmt.Errorf("authorize businesses failed, get business by id failed, err: %+v", err)
	}

	return am.AuthorizeByBusiness(ctx, header, action, businesses...)
}

func (am *AuthManager) GenBusinessAuditNoPermissionResp(ctx context.Context, header http.Header, businessID int64) (*metadata.BaseResp, error) {
	var p metadata.Permission
	p.SystemID = authcenter.SystemIDCMDB
	p.SystemName = authcenter.SystemNameCMDB
	p.ScopeType = authcenter.ScopeTypeIDSystem
	p.ScopeTypeName = authcenter.ScopeTypeIDSystemName
	p.ScopeID = strconv.FormatInt(businessID, 10)
	p.ActionID = string(authcenter.Get)
	p.ActionName = authcenter.ActionIDNameMap[authcenter.Get]

	p.Resources = [][]metadata.Resource{
		{{
			ResourceType:     string(authcenter.SysBusinessInstance),
			ResourceTypeName: authcenter.ResourceTypeIDMap[authcenter.SysBusinessInstance],
		}},
	}

	businesses, err := am.collectBusinessByIDs(ctx, header, businessID)
	if err != nil {
		return nil, err
	}
	if len(businesses) != 1 {
		return nil, errors.New("get business detail failed")
	}
	p.ScopeName = businesses[0].BKAppNameField
    p.ResourceType = p.Resources[0][0].ResourceType
    p.ResourceTypeName = p.Resources[0][0].ResourceTypeName
	resp := metadata.NewNoPermissionResp([]metadata.Permission{p})
	return &resp, nil
}

func (am *AuthManager) UpdateRegisteredBusiness(ctx context.Context, header http.Header, businesses ...BusinessSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(businesses) == 0 {
		return nil
	}

	// make auth resources
	resources := am.MakeResourcesByBusiness(header, meta.EmptyAction, businesses...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredBusinessByID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	businesses, err := am.collectBusinessByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered businesses failed, get businesses by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredBusiness(ctx, header, businesses...)
}

func (am *AuthManager) UpdateRegisteredBusinessByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	businesses, err := am.collectBusinessByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered businesses failed, get businesses by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredBusiness(ctx, header, businesses...)
}

func (am *AuthManager) DeregisterBusinessByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	businesses, err := am.collectBusinessByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("deregister businesses failed, get businesses by id failed, err: %+v", err)
	}
	return am.DeregisterBusinesses(ctx, header, businesses...)
}

func (am *AuthManager) RegisterBusinesses(ctx context.Context, header http.Header, businesses ...BusinessSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(businesses) == 0 {
		return nil
	}

	// make auth resources
	resources := am.MakeResourcesByBusiness(header, meta.EmptyAction, businesses...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) RegisterBusinessesByID(ctx context.Context, header http.Header, businessIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(businessIDs) == 0 {
		return nil
	}

	businesses, err := am.collectBusinessByIDs(ctx, header, businessIDs...)
	if err != nil {
		return fmt.Errorf("get businesses by id failed, err: %+v", err)
	}
	return am.RegisterBusinesses(ctx, header, businesses...)
}

func (am *AuthManager) DeregisterBusinesses(ctx context.Context, header http.Header, businesses ...BusinessSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(businesses) == 0 {
		return nil
	}

	// make auth resources
	resources := am.MakeResourcesByBusiness(header, meta.EmptyAction, businesses...)

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterBusinessesByID(ctx context.Context, header http.Header, businessIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(businessIDs) == 0 {
		return nil
	}

	businesses, err := am.collectBusinessByIDs(ctx, header, businessIDs...)
	if err != nil {
		return fmt.Errorf("get businesses by id failed, err: %+v", err)
	}
	return am.DeregisterBusinesses(ctx, header, businesses...)
}

func (am *AuthManager) GenBusinessNoPermissionResp(ctx context.Context, header http.Header, businessID int64) (*metadata.BaseResp, error) {
	var p metadata.Permission
	p.SystemID = authcenter.SystemIDCMDB
	p.SystemName = authcenter.SystemNameCMDB
	p.ScopeType = authcenter.ScopeTypeIDSystem
	p.ScopeTypeName = authcenter.ScopeTypeIDSystemName
	p.ActionID = string(authcenter.Get)
	p.ActionName = authcenter.ActionIDNameMap[authcenter.Get]
    p.ResourceType = string(authcenter.SysBusinessInstance)
	p.ResourceTypeName = authcenter.ResourceTypeIDMap[authcenter.SysBusinessInstance]
	p.Resources = [][]metadata.Resource{
		{{
			ResourceType:     string(authcenter.SysBusinessInstance),
			ResourceTypeName: authcenter.ResourceTypeIDMap[authcenter.SysBusinessInstance],
		}},
	}

	businesses, err := am.collectBusinessByIDs(ctx, header, businessID)
	if err != nil {
		return nil, err
	}
	if len(businesses) != 1 {
		return nil, errors.New("get business detail failed")
	}
	p.ScopeName = businesses[0].BKAppNameField
	resp := metadata.NewNoPermissionResp([]metadata.Permission{p})
	return &resp, nil
}

// func (am *AuthManager) AuthorizeBusinessesByID(ctx context.Context, header http.Header, action meta.Action, businessIDs ...int64) error {
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	businesses, err := am.collectBusinessByIDs(ctx, header, businessIDs...)
// 	if err != nil {
// 		return fmt.Errorf("get businesses by id failed, err: %+v", err)
// 	}
// 	return am.AuthorizeByBusiness(ctx, header, action, businesses...)
// }
