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

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/ac/parser"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// MakeAuthorizedAuditListCondition make a query condition, with which user can only search audit log under it.
// ==> [{"bk_biz_id":2,"op_target":{"$in":["module"]}}]
func (am *AuthManager) MakeAuthorizedAuditListCondition(ctx context.Context, header http.Header, businessID int64) (cond []map[string]interface{}, hasAuthorization bool, err error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	// businessID 0 means audit log priority of special model on any business

	commonInfo, err := parser.ParseCommonInfo(&header)
	if err != nil {
		return nil, false, fmt.Errorf("parse user info from request header failed, %+v, rid: %s", err, rid)
	}

	businessIDs := make([]int64, 0)
	if businessID == 0 {
		input := meta.ListAuthorizedResourcesParam{
			UserName:     commonInfo.User.UserName,
			BizID:        0,
			ResourceType: meta.Business,
			Action:       meta.ViewBusinessResource,
		}
		authorizedResources, err := am.clientSet.AuthServer().ListAuthorizedResources(ctx, header, input)
		if err != nil {
			blog.Errorf("make condition from authorization failed, get authorized businesses failed, err: %+v, rid: %s", err, rid)
			return nil, false, fmt.Errorf("make condition from authorization failed, get authorized businesses failed, err: %+v", err)
		}

		for _, resourceID := range authorizedResources {
			bizID, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				return nil, false, fmt.Errorf("make condition from authorization failed, parse authorized businesses id(%s) failed, err: %+v", resourceID, err)
			}
			businessIDs = append(businessIDs, bizID)
		}
	}
	businessIDs = append(businessIDs, 0)
	blog.V(5).Infof("audit on business %+v to be check", businessIDs)

	cond = make([]map[string]interface{}, 0)

	asst, err := am.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), header, &metadata.QueryCondition{Condition: map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline}})
	if err != nil || !asst.Result {
		blog.Errorf("[audit] failed to find mainline association, err: %v, resp: %v, rid: %s", err, asst, rid)
		return
	}

	for _, businessID := range businessIDs {
		input := meta.ListAuthorizedResourcesParam{
			UserName:     commonInfo.User.UserName,
			BizID:        businessID,
			ResourceType: meta.AuditLog,
			Action:       meta.Find,
		}
		auditList, err := am.clientSet.AuthServer().ListAuthorizedResources(ctx, header, input)
		if err != nil {
			blog.Errorf("make condition from authorization failed, get authorized businesses failed, err: %+v, rid: %s", err, rid)
			return nil, false, fmt.Errorf("make condition from authorization failed, get authorized businesses failed, err: %+v", err)
		}
		blog.Infof("get authorized audit by business %d result: %s", businessID, auditList)
		blog.InfoJSON("get authorized audit by business %s result: %s", businessID, auditList)

		resourceTypes := make([]metadata.ResourceType, 0)
		auditTypes := make([]metadata.AuditType, 0)
		for _, modelID := range auditList {
			id := util.GetStrByInterface(modelID)
			isMainline := false
			for _, mainline := range asst.Data.Info {
				if mainline.ObjectID == id || mainline.AsstObjID == id {
					isMainline = true
					break
				}
			}
			resourceTypes = append(resourceTypes, metadata.GetResourceTypeByObjID(id, isMainline))
			auditTypes = append(auditTypes, metadata.GetAuditTypeByObjID(id, isMainline))
		}

		if len(resourceTypes) == 0 {
			continue
		}
		hasAuthorization = true
		item := map[string]interface{}{
			common.BKResourceTypeField: map[string]interface{}{
				common.BKDBIN: resourceTypes,
			},
			common.BKAuditTypeField: map[string]interface{}{
				common.BKDBIN: auditTypes,
			},
		}
		if businessID != 0 {
			item[common.BKAppIDField] = businessID
		}
		cond = append(cond, item)
	}

	blog.V(5).Infof("MakeAuthorizedAuditListCondition result: %+v, rid: %s", cond, rid)
	return cond, hasAuthorization, nil
}

func (am *AuthManager) AuthorizeAuditRead(ctx context.Context, header http.Header, businessID int64) error {
	if !am.Enabled() {
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
	instances := make([][]metadata.IamResourceInstance, 0)
	if businessID > 0 {
		businesses, err := am.collectBusinessByIDs(ctx, header, businessID)
		if err != nil {
			return nil, err
		}
		if len(businesses) != 1 {
			return nil, errors.New("get business detail failed")
		}
		instances = append(instances, []metadata.IamResourceInstance{{
			Type: string(iam.Business),
			ID:   strconv.FormatInt(businessID, 10),
			Name: businesses[0].BKAppNameField,
		}})
	}
	instances = append(instances, []metadata.IamResourceInstance{{
		Type: string(iam.SysAuditLog),
		Name: iam.ResourceTypeIDMap[iam.SysAuditLog],
	}})
	permission := &metadata.IamPermission{
		SystemID: iam.SystemIDCMDB,
		Actions: []metadata.IamAction{{
			ID: string(iam.FindAuditLog),
			RelatedResourceTypes: []metadata.IamResourceType{{
				SystemID:  iam.SystemIDCMDB,
				Type:      string(iam.SysAuditLog),
				Instances: instances,
			}},
		}},
	}
	resp := metadata.NewNoPermissionResp(permission)
	return &resp, nil
}
