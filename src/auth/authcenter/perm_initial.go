/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package authcenter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/auth/meta"
	"configcenter/src/common/blog"
)

func (ac *AuthCenter) Init(ctx context.Context, configs meta.InitConfig) error {
	if err := ac.initAuthResources(ctx, configs); err != nil {
		return fmt.Errorf("initial auth resources failed, err: %v", err)
	}

	if err := ac.initDefaultUserRoleWithAuth(ctx); err != nil {
		return fmt.Errorf("initail default user's role with auth failed: %v", err)
	}

	// TODO: remove this operation, when it's auth is roll back.
	// remove biz instance resource
	if err := ac.authClient.DeleteResources(ctx, http.Header{}, "biz", BizInstance); err != nil {
		blog.Errorf("delete biz instance resource from iam failed, err: %v", err)
	}
	// TODO: remove this operation, when it's auth is roll back
	// remove biz audit log resource
	if err := ac.authClient.DeleteResources(ctx, http.Header{}, "biz", BizAuditLog); err != nil {
		blog.Errorf("delete biz audit log resource from iam failed, err: %v", err)
	}

	blog.Info("initial auth center success.")
	return nil
}

func (ac *AuthCenter) initAuthResources(ctx context.Context, configs meta.InitConfig) error {
	header := http.Header{}
	if err := ac.authClient.RegistSystem(ctx, header, expectSystem); err != nil && err != ErrDuplicated {
		return err
	}

	if err := ac.authClient.UpdateSystem(ctx, header, System{SystemID: expectSystem.SystemID, SystemName: expectSystem.SystemName}); err != nil {
		return err
	}

	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDSystem, expectSystemResourceType); err != nil {
		return err
	}
	if err := ac.authClient.UpsertResourceTypeBatch(ctx, header, SystemIDCMDB, ScopeTypeIDBiz, expectBizResourceType); err != nil {
		return err
	}

	// init model classification
	clsName2ID := map[string]int64{}
	for _, cls := range configs.Classifications {
		bizID, _ := cls.Metadata.Label.GetBusinessID()
		clsName2ID[fmt.Sprintf("%d:%s", bizID, cls.ClassificationID)] = cls.ID

		scopeType := ScopeTypeIDSystem
		ScopeID := SystemIDCMDB
		groupType := SysModelGroup

		if bizID > 0 {
			scopeType = ScopeTypeIDBiz
			ScopeID = strconv.FormatInt(bizID, 10)
			groupType = BizModelGroup
		}

		info := RegisterInfo{
			CreatorID:   "system",
			CreatorType: "user",
			Resources: []ResourceEntity{
				{
					ResourceType: groupType,
					ResourceID: []RscTypeAndID{
						{ResourceType: groupType, ResourceID: strconv.FormatInt(cls.ID, 10)},
					},
					ResourceName: cls.ClassificationName,
					ScopeInfo: ScopeInfo{
						ScopeType: scopeType,
						ScopeID:   ScopeID,
					},
				},
			},
		}

		if err := ac.authClient.registerResource(ctx, header, &info); err != nil && err != ErrDuplicated {
			return err
		}
	}

	// init model description
	for _, model := range configs.Models {
		bizID, _ := model.Metadata.Label.GetBusinessID()

		scopeType := ScopeTypeIDSystem
		ScopeID := SystemIDCMDB
		modelType := SysModel

		if bizID > 0 {
			scopeType = ScopeTypeIDBiz
			ScopeID = strconv.FormatInt(bizID, 10)
			modelType = BizModel
		}

		info := RegisterInfo{
			CreatorID:   "system",
			CreatorType: "user",
			Resources: []ResourceEntity{
				{
					ResourceType: modelType,
					ResourceID: []RscTypeAndID{
						{ResourceType: modelType, ResourceID: strconv.FormatInt(model.ID, 10)},
					},
					ResourceName: model.ObjectName,
					ScopeInfo: ScopeInfo{
						ScopeType: scopeType,
						ScopeID:   ScopeID,
					},
				},
			},
		}

		if err := ac.authClient.registerResource(ctx, header, &info); err != nil && err != ErrDuplicated {
			return err
		}
	}

	// init business inst
	for _, biz := range configs.Bizs {
		bkbiz := RegisterInfo{
			CreatorID:   "system",
			CreatorType: "user",
			Resources: []ResourceEntity{
				{
					ResourceType: SysBusinessInstance,
					ResourceID: []RscTypeAndID{
						{ResourceType: SysBusinessInstance, ResourceID: strconv.FormatInt(biz.BizID, 10)},
					},
					ResourceName: biz.BizName,
					ScopeInfo: ScopeInfo{
						ScopeType: "system",
						ScopeID:   SystemIDCMDB,
					},
				},
			},
		}

		if err := ac.authClient.registerResource(ctx, header, &bkbiz); err != nil && err != ErrDuplicated {
			return err
		}
	}

	// init association kind
	for _, kind := range configs.AssociationKinds {
		info := RegisterInfo{
			CreatorID:   "system",
			CreatorType: "user",
			Resources: []ResourceEntity{
				{
					ResourceType: SysAssociationType,
					ResourceID: []RscTypeAndID{
						{ResourceType: SysAssociationType, ResourceID: strconv.FormatInt(kind.ID, 10)},
					},
					ResourceName: kind.AssociationKindID,
					ScopeInfo: ScopeInfo{
						ScopeType: ScopeTypeIDSystem,
						ScopeID:   SystemIDCMDB,
					},
				},
			},
		}

		if err := ac.authClient.registerResource(ctx, header, &info); err != nil && err != ErrDuplicated {
			return err
		}
	}
	return nil
}

func (ac *AuthCenter) initDefaultUserRoleWithAuth(ctx context.Context) error {
	normalActions := []RoleAction{
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizCustomQuery,
			ActionID:       Get,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizInstance,
			ActionID:       Get,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizAuditLog,
			ActionID:       Get,
		},
	}
	rolesWithAuth := []RoleWithAuthResources{
		bizOperatorRoleAuth,
		{
			RoleTemplateName: "产品",
			TemplateID:       "product_manager",
			Desc:             "Product Manager",
			ResourceActions:  normalActions,
		},
		{
			RoleTemplateName: "开发",
			TemplateID:       "product_developer",
			Desc:             "Product Developer",
			ResourceActions:  normalActions,
		},
		{
			RoleTemplateName: "测试",
			TemplateID:       "product_tester",
			Desc:             "Product Tester",
			ResourceActions:  normalActions,
		},
		{
			RoleTemplateName: "职能化",
			TemplateID:       "product_operator",
			Desc:             "Product Operator",
			ResourceActions:  normalActions,
		},
	}
	for _, role := range rolesWithAuth {
		id, err := ac.authClient.RegisterUserRole(ctx, http.Header{}, role)
		if err != nil {
			return err
		}
		blog.Infof("register auth with role: %s, id: %d", role.RoleTemplateName, id)
	}
	return nil
}

var bizOperatorRoleAuth = RoleWithAuthResources{
	RoleTemplateName: "运维",
	TemplateID:       "business_maintainer",
	Desc:             "a business's maintainer",
	ResourceActions: []RoleAction{
		// business host instance related
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizHostInstance,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizHostInstance,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizHostInstance,
			ActionID:       Delete,
		},
		// dynamic group related
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizCustomQuery,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizCustomQuery,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizCustomQuery,
			ActionID:       Delete,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizCustomQuery,
			ActionID:       Get,
		},
		// biz topology related
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizTopology,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizTopology,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizTopology,
			ActionID:       Delete,
		},
		// process related
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessInstance,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessInstance,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessInstance,
			ActionID:       Delete,
		},
		// model classification related
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizModelGroup,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizModelGroup,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizModelGroup,
			ActionID:       Delete,
		},
		// model related
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizModel,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizModel,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizModel,
			ActionID:       Delete,
		},
		// model instance
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizInstance,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizInstance,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizInstance,
			ActionID:       Delete,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizInstance,
			ActionID:       Get,
		},
		// audit related
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizAuditLog,
			ActionID:       Get,
		},
		// service template related.
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessServiceTemplate,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessServiceTemplate,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessServiceTemplate,
			ActionID:       Delete,
		},
		// service category related.
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessServiceCategory,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessServiceCategory,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessServiceCategory,
			ActionID:       Delete,
		},
		// service instance related.
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessServiceInstance,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessServiceInstance,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizProcessServiceInstance,
			ActionID:       Delete,
		},
		// global resource related.
		{
			ScopeTypeID:    ScopeTypeIDSystem,
			ResourceTypeID: SysBusinessInstance,
			ActionID:       Get,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizSetTemplate,
			ActionID:       Create,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizSetTemplate,
			ActionID:       Edit,
		},
		{
			ScopeTypeID:    ScopeTypeIDBiz,
			ResourceTypeID: BizSetTemplate,
			ActionID:       Delete,
		},
	},
}
