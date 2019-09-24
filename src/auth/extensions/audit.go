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
	"configcenter/src/auth/parser"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (am *AuthManager) CollectAuditCategoryByBusinessID(ctx context.Context, header http.Header, businessID int64) ([]AuditCategorySimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	query := metadata.QueryInput{
		Condition: condition.CreateCondition().Field(common.BKAppIDField).Eq(businessID).ToMapStr(),
	}
	response, err := am.clientSet.CoreService().Audit().SearchAuditLog(ctx, header, query)
	if nil != err {
		blog.Errorf("collect audit category by business %d failed, get audit log failed, err: %+v, rid: %s", businessID, err, rid)
		return nil, fmt.Errorf("collect audit category by business %d failed, get audit log failed, err: %+v", businessID, err)
	}

	data, err := mapstr.NewFromInterface(response.Data)
	if nil != err {
		blog.Errorf("collect audit category by business %d failed, parse response data failed, data: %+v, error info is %+v, rid: %s", businessID, response.Data, err, rid)
		return nil, fmt.Errorf("collect audit category by business %d failed, parse response data failed, error info is %+v", businessID, err)
	}
	auditLogs, err := data.MapStrArray("info")
	if nil != err {
		blog.Errorf("collect audit category by business %d failed, extract audit log from response data failed, data: %+v, error info is %+v, rid: %s", businessID, response.Data, err, rid)
		return nil, fmt.Errorf("collect audit category by business %d failed, extract audit log from response data failed, error info is %+v", businessID, err)
	}

	categories := make([]AuditCategorySimplify, 0)
	modelIDFound := map[string]bool{}
	for _, item := range auditLogs {
		category := &AuditCategorySimplify{}
		category, err := category.Parse(item)
		if err != nil {
			blog.Errorf("parse audit category simplify failed, category: %+v, err: %+v, rid: %s", category, err, rid)
			continue
		}
		if _, exist := modelIDFound[category.BKOpTargetField]; exist == false {
			categories = append(categories, *category)
			modelIDFound[category.BKOpTargetField] = true
		}
	}

	blog.V(4).Infof("list audit categories by business %d result: %+v, rid: %s", businessID, categories, rid)
	return categories, nil
}

func (am *AuthManager) ExtractBusinessIDFromAuditCategories(categories ...AuditCategorySimplify) (int64, error) {
	var businessID int64
	for idx, category := range categories {
		bizID := category.BKAppIDField
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("get multiple business ID from audit categories")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) MakeResourcesByAuditCategories(ctx context.Context, header http.Header, action meta.Action, businessID int64, categories ...AuditCategorySimplify) ([]meta.ResourceAttribute, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// prepare resource layers for authorization
	resources := make([]meta.ResourceAttribute, 0)
	for _, category := range categories {
		// instance
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:       action,
				Type:         meta.AuditLog,
				Name:         category.BKOpTargetField,
				InstanceIDEx: category.BKOpTargetField,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}

	blog.V(9).Infof("MakeResourcesByAuditCategories: %+v, rid: %s", resources, rid)
	return resources, nil
}

func (am *AuthManager) RegisterAuditCategories(ctx context.Context, header http.Header, categories ...AuditCategorySimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(categories) == 0 {
		return nil
	}

	businessID, err := am.ExtractBusinessIDFromAuditCategories(categories...)
	if err != nil {
		return fmt.Errorf("extract business id from audit categories failed, err: %+v", err)
	}

	resources, err := am.MakeResourcesByAuditCategories(ctx, header, meta.EmptyAction, businessID, categories...)
	if err != nil {
		return fmt.Errorf("make auth resource by audit categories failed, err: %+v", err)
	}

	if err := am.Authorize.RegisterResource(ctx, resources...); err != nil {
		return fmt.Errorf("register audit categories failed, err: %+v", err)
	}
	return nil
}

// MakeAuthorizedAuditListCondition make a query condition, with which user can only search audit log under it.
// ==> [{"bk_biz_id":2,"op_target":{"$in":["module"]}}]
func (am *AuthManager) MakeAuthorizedAuditListCondition(ctx context.Context, header http.Header, businessID int64) (cond []mapstr.MapStr, hasAuthorization bool, err error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	// businessID 0 means audit log priority of special model on any business

	commonInfo, err := parser.ParseCommonInfo(&header)
	if err != nil {
		return nil, false, fmt.Errorf("parse user info from request header failed, %+v, rid: %s", err, rid)
	}

	businessIDs := make([]int64, 0)
	if businessID == 0 {
		ids, err := am.Authorize.GetAnyAuthorizedBusinessList(ctx, commonInfo.User)
		if err != nil {
			blog.Errorf("make condition from authorization failed, get authorized businesses failed, err: %+v, rid: %s", err, rid)
			return nil, false, fmt.Errorf("make condition from authorization failed, get authorized businesses failed, err: %+v", err)
		}
		businessIDs = ids
	}
	businessIDs = append(businessIDs, 0)
	blog.V(5).Infof("audit on business %+v to be check", businessIDs)

	authorizedBusinessModelMap := map[int64][]string{}
	for _, businessID := range businessIDs {
		auditList, err := am.Authorize.GetAuthorizedAuditList(ctx, commonInfo.User, businessID)
		if err != nil {
			blog.Errorf("get authorized audit by business %d failed, err: %+v, rid: %s", businessID, err, rid)
			return nil, false, fmt.Errorf("get authorized audit by business %d failed, err: %+v", businessID, err)
		}
		blog.Infof("get authorized audit by business %d result: %s", businessID, auditList)
		blog.InfoJSON("get authorized audit by business %s result: %s", businessID, auditList)

		modelIDs := make([]string, 0)
		for _, authorizedList := range auditList {
			for _, resourceID := range authorizedList.ResourceIDs {
				if len(resourceID) == 0 {
					continue
				}
				modelID := resourceID[len(resourceID)-1].ResourceID
				id := util.GetStrByInterface(modelID)
				modelIDs = append(modelIDs, id)
			}
		}

		if len(modelIDs) == 0 {
			continue
		}
		authorizedBusinessModelMap[businessID] = modelIDs
	}

	cond = make([]mapstr.MapStr, 0)

	// extract authorization on any business
	if _, ok := authorizedBusinessModelMap[0]; ok == true {
		if len(authorizedBusinessModelMap[0]) > 0 {
			hasAuthorization = true
			item := condition.CreateCondition()
			item.Field(common.BKOpTargetField).In(authorizedBusinessModelMap[0])

			cond = append(cond, item.ToMapStr())
			delete(authorizedBusinessModelMap, 0)
		}
	}

	// extract authorization on special business and object
	for businessID, objectIDs := range authorizedBusinessModelMap {
		hasAuthorization = true
		item := condition.CreateCondition()
		item.Field(common.BKOpTargetField).In(objectIDs)
		item.Field(common.BKAppIDField).Eq(businessID)

		cond = append(cond, item.ToMapStr())
	}

	blog.V(5).Infof("MakeAuthorizedAuditListCondition result: %+v, rid: %s", cond, rid)
	return cond, hasAuthorization, nil
}

func (am *AuthManager) AuthorizeAuditRead(ctx context.Context, header http.Header, businessID int64) error {
	if am.Enabled() == false {
		return nil
	}

	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Action: meta.Find,
			Type:   meta.AuditLog,
		},
		SupplierAccount: util.GetOwnerID(header),
		BusinessID:      businessID,
	}
	return am.authorize(ctx, header, businessID, resource)
}

func (am *AuthManager) GenAuthorizeAuditReadNoPermissionsResponse(ctx context.Context, header http.Header, businessID int64) (*metadata.BaseResp, error) {
	var p metadata.Permission
	p.SystemID = authcenter.SystemIDCMDB
	p.SystemName = authcenter.SystemNameCMDB
	p.ScopeID = strconv.FormatInt(businessID, 10)
	p.ActionID = string(authcenter.Get)
	p.ActionName = authcenter.ActionIDNameMap[authcenter.Get]
	if businessID > 0 {
		p.Resources = [][]metadata.Resource{
			{{
				ResourceType:     string(authcenter.BizAuditLog),
				ResourceTypeName: authcenter.ResourceTypeIDMap[authcenter.BizAuditLog],
			}},
		}
		businesses, err := am.collectBusinessByIDs(ctx, header, businessID)
		if err != nil {
			return nil, err
		}
		if len(businesses) != 1 {
			return nil, errors.New("get business detail failed")
		}
		p.ScopeType = authcenter.ScopeTypeIDBiz
		p.ScopeTypeName = authcenter.ScopeTypeIDBizName
		p.ScopeID = strconv.FormatInt(businessID, 10)
		p.ScopeName = businesses[0].BKAppNameField
	} else {
		p.ScopeType = authcenter.ScopeTypeIDSystem
		p.ScopeTypeName = authcenter.ScopeTypeIDSystemName
		p.Resources = [][]metadata.Resource{
			{{
				ResourceType:     string(authcenter.SysAuditLog),
				ResourceTypeName: authcenter.ResourceTypeIDMap[authcenter.SysAuditLog],
			}},
		}
	}
    p.ResourceType = p.Resources[0][0].ResourceType
    p.ResourceTypeName = p.Resources[0][0].ResourceTypeName
	resp := metadata.NewNoPermissionResp([]metadata.Permission{p})
	return &resp, nil
}
