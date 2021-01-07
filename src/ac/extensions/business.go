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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * business related auth interface
 */

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

func (am *AuthManager) AuthorizeByBusiness(ctx context.Context, header http.Header, action meta.Action, businesses ...BusinessSimplify) error {
	if !am.Enabled() {
		return nil
	}

	// make auth resources
	resources := am.MakeResourcesByBusiness(header, action, businesses...)

	return am.batchAuthorize(ctx, header, resources...)
}

func (am *AuthManager) AuthorizeByBusinessID(ctx context.Context, header http.Header, action meta.Action, businessIDs ...int64) error {
	if !am.Enabled() {
		return nil
	}

	businesses, err := am.collectBusinessByIDs(ctx, header, businessIDs...)
	if err != nil {
		return fmt.Errorf("authorize businesses failed, get business by id failed, err: %+v", err)
	}

	return am.AuthorizeByBusiness(ctx, header, action, businesses...)
}

func (am *AuthManager) GenBusinessAuditNoPermissionResp(ctx context.Context, header http.Header, businessID int64) (*metadata.BaseResp, error) {
	businesses, err := am.collectBusinessByIDs(ctx, header, businessID)
	if err != nil {
		return nil, err
	}
	if len(businesses) != 1 {
		return nil, errors.New("get business detail failed")
	}
	permission := &metadata.IamPermission{
		SystemID: iam.SystemIDCMDB,
		Actions: []metadata.IamAction{{
			ID: string(iam.FindAuditLog),
			RelatedResourceTypes: []metadata.IamResourceType{{
				SystemID: iam.SystemIDCMDB,
				Type:     string(iam.SysAuditLog),
				Instances: [][]metadata.IamResourceInstance{{{
					Type: string(iam.Business),
					ID:   strconv.FormatInt(businessID, 10),
				}, {
					Type: string(iam.SysAuditLog),
				}}},
			}},
		}},
	}
	resp := metadata.NewNoPermissionResp(permission)
	return &resp, nil
}
