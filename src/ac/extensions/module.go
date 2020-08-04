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

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * module instance
 */

func (am *AuthManager) collectModuleByModuleIDs(ctx context.Context, header http.Header, moduleIDs ...int64) ([]ModuleSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	moduleIDs = util.IntArrayUnique(moduleIDs)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKModuleIDField).In(moduleIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDModule, &cond)
	if err != nil {
		blog.V(3).Infof("get modules by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get modules by id failed, err: %+v", err)
	}
	modules := make([]ModuleSimplify, 0)
	for _, cls := range result.Data.Info {
		module := ModuleSimplify{}
		_, err = module.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get modules by object failed, err: %+v", err)
		}
		modules = append(modules, module)
	}
	return modules, nil
}

func (am *AuthManager) extractBusinessIDFromModules(modules ...ModuleSimplify) (int64, error) {
	var businessID int64
	for idx, module := range modules {
		bizID := module.BKAppIDField
		// we should ignore metadata.LabelBusinessID field not found error
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("authorization failed, get multiple business ID from modules")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) MakeResourcesByModule(header http.Header, action meta.Action, businessID int64, modules ...ModuleSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, module := range modules {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ModelModule,
				Name:       module.BKModuleNameField,
				InstanceID: module.BKModuleIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) AuthorizeByModuleID(ctx context.Context, header http.Header, action meta.Action, ids ...int64) error {
	if !am.Enabled() {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}
	if !am.RegisterModuleEnabled {
		return nil
	}

	modules, err := am.collectModuleByModuleIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered modules failed, get modules by id failed, err: %+v", err)
	}
	return am.AuthorizeByModule(ctx, header, action, modules...)
}

func (am *AuthManager) GenModuleSetNoPermissionResp() *metadata.BaseResp {
	permission := &metadata.IamPermission{
		SystemID: iam.SystemIDCMDB,
		Actions: []metadata.IamAction{{
			ID:                   string(iam.EditBusinessLayer),
			RelatedResourceTypes: nil,
		}},
	}
	resp := metadata.NewNoPermissionResp(permission)
	return &resp
}

func (am *AuthManager) AuthorizeByModule(ctx context.Context, header http.Header, action meta.Action, modules ...ModuleSimplify) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}

	if !am.RegisterModuleEnabled {
		return nil
	}

	if len(modules) == 0 {
		return nil
	}
	if am.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, modules: %+v, rid: %s", modules, rid)
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromModules(modules...)
	if err != nil {
		return fmt.Errorf("authorize modules failed, extract business id from modules failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByModule(header, action, bizID, modules...)

	return am.batchAuthorize(ctx, header, resources...)
}
