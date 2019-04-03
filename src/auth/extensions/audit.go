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

func (am *AuthManager) CollectAuditCategoryByBusinessID(ctx context.Context, header http.Header, businessID int64) ([]AuditCategorySimplify, error) {
	query := &metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKAppIDField).Eq(businessID).ToMapStr(),
	}
	response, err := am.clientSet.CoreService().Instance().ReadInstance(context.Background(), header, common.BKTableNameOperationLog, query)
	if err != nil {
		blog.Errorf("get models by business %d failed, err: %+v", businessID, err)
		return nil, fmt.Errorf("get models by business %d failed, err: %+v", businessID, err)
	}

	categories := make([]AuditCategorySimplify, 0)
	modelIDs := make([]string, 0)
	for _, item := range response.Data.Info {
		category := &AuditCategorySimplify{}
		category, err :=  category.Parse(item)
		if err != nil {
			blog.Errorf("parse audit category simplify failed, category: %+v, err: %+v", category, err)
			continue
		}
		modelIDs = append(modelIDs, category.BKOpTargetField)
		categories = append(categories, *category)
	}
	modelIDs = util.StrArrayUnique(modelIDs)
	objects, err := am.collectObjectsByObjectIDs(ctx, header, modelIDs...)
	if err != nil {
		blog.Errorf("collectObjectsByObjectIDs failed, model: %+v, err: %+v", modelIDs, err)
		return nil, fmt.Errorf("get audit category related models failed, err: %+v", err)
	}
	objectIDMap := map[string]int64{}
	for _, object := range objects {
		objectIDMap[object.ObjectID] = object.ID
	}
	
	// invalid categories will be filter out
	validCategories := make([]AuditCategorySimplify, 0)
	for _, category := range categories {
		modelID, existed := objectIDMap[category.BKOpTargetField]
		if existed == true {
			category.ModelID = modelID
			validCategories = append(validCategories, category)
		} else {
			blog.Errorf("unexpect audit op_target: %s", category.BKOpTargetField)
		}
	}

	blog.V(4).Infof("list audit categories by business %d result: %+v", businessID, validCategories)
	return validCategories, nil
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
	// step2 prepare resource layers for authorization
	resources := make([]meta.ResourceAttribute, 0)
	for _, category := range categories {
		// instance
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Model,
				Name:       category.BKOpTargetField,
				InstanceID: category.ModelID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}

	blog.V(9).Infof("MakeResourcesByAuditCategories: %+v", resources)
	return resources, nil
}

func (am *AuthManager) RegisterAuditCategories(ctx context.Context, header http.Header, categories ...AuditCategorySimplify) error {
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
func (am *AuthManager) MakeAuthorizedAuditListCondition(ctx context.Context, user meta.UserInfo) (map[string]interface{}, error) {
	return nil, nil
}
