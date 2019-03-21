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
 * model classification which used for manage models as group
 */

func (am *AuthManager) collectClassificationsByClassificationIDs(ctx context.Context, header http.Header, classificationIDs ...string) ([]metadata.Classification, error) {

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKClassificationIDField).In(classificationIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameObjClassifiction, &cond)
	if err != nil {
		blog.V(3).Infof("get classification by id failed, err: %+v", err)
		return nil, fmt.Errorf("get classification by id failed, err: %+v", err)
	}
	classifications := make([]metadata.Classification, 0)
	for _, cls := range result.Data.Info {
		classification := metadata.Classification{}
		_, err = classification.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get classication by object failed, err: %+v", err)
		}
		classifications = append(classifications, classification)
	}
	return classifications, nil
}

func (am *AuthManager) collectClassificationsByRawIDs(ctx context.Context, header http.Header, ids ...int64) ([]metadata.Classification, error) {

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKFieldID).In(ids).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameObjClassifiction, &cond)
	if err != nil {
		blog.V(3).Infof("get classification by id failed, err: %+v", err)
		return nil, fmt.Errorf("get classification by id failed, err: %+v", err)
	}
	classifications := make([]metadata.Classification, 0)
	for _, cls := range result.Data.Info {
		classification := metadata.Classification{}
		_, err = classification.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get classication by object failed, err: %+v", err)
		}
		classifications = append(classifications, classification)
	}
	return classifications, nil
}

func (am *AuthManager) extractBusinessIDFromClassifications(classifications ...metadata.Classification) (int64, error) {
	var businessID int64
	for idx, classification := range classifications {
		bizID, err := classification.Metadata.Label.Int64(metadata.LabelBusinessID)
		// we should ignore metadata.LabelBusinessID field not found error
		if err != nil && err != metadata.LabelKeyNotExistError {
			return 0, fmt.Errorf("parse biz id from classification: %+v failed, err: %+v", classification, err)
		}
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("authorization failed, get multiple business ID from objects")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) makeResourcesByClassifications(header http.Header, action meta.Action, businessID int64, classifications ...metadata.Classification) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, classification := range classifications {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Model,
				Name:       classification.ClassificationID,
				InstanceID: classification.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) AuthorizeByClassification(ctx context.Context, header http.Header, action meta.Action, classifications ...metadata.Classification) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromClassifications(classifications...)
	if err != nil {
		return fmt.Errorf("authorize classifications failed, extract business id from classification failed, err: %+v", err)
	}

	// make auth resources
	resources := am.makeResourcesByClassifications(header, action, bizID, classifications...)

	return am.authorize(ctx, header, bizID, resources...)
}

func (am *AuthManager) UpdateRegisteredClassification(ctx context.Context, header http.Header, classifications ...metadata.Classification) error {
	// extract business id
	bizID, err := am.extractBusinessIDFromClassifications(classifications...)
	if err != nil {
		return fmt.Errorf("authorize classifications failed, extract business id from classification failed, err: %+v", err)
	}

	// make auth resources
	resources := am.makeResourcesByClassifications(header, meta.EmptyAction, bizID, classifications...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredClassificationByID(ctx context.Context, header http.Header, classificationIDs ...string) error {
	classifications, err := am.collectClassificationsByClassificationIDs(ctx, header, classificationIDs...)
	if err != nil {
		return fmt.Errorf("update registered classifications failed, get classfication by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredClassification(ctx, header, classifications...)
}

func (am *AuthManager) UpdateRegisteredClassificationByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	classifications, err := am.collectClassificationsByRawIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered classifications failed, get classfication by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredClassification(ctx, header, classifications...)
}

func (am *AuthManager) DeregisterClassificationByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	classifications, err := am.collectClassificationsByRawIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("deregister classifications failed, get classfication by id failed, err: %+v", err)
	}
	return am.DeregisterClassification(ctx, header, classifications...)
}

func (am *AuthManager) RegisterClassification(ctx context.Context, header http.Header, classifications ...metadata.Classification) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromClassifications(classifications...)
	if err != nil {
		return fmt.Errorf("register classifications failed, extract business id from classification failed, err: %+v", err)
	}

	// make auth resources
	resources := am.makeResourcesByClassifications(header, meta.EmptyAction, bizID, classifications...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterClassification(ctx context.Context, header http.Header, classifications ...metadata.Classification) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromClassifications(classifications...)
	if err != nil {
		return fmt.Errorf("deregister classifications failed, extract business id from classification failed, err: %+v", err)
	}

	// make auth resources
	resources := am.makeResourcesByClassifications(header, meta.EmptyAction, bizID, classifications...)

	return am.Authorize.DeregisterResource(ctx, resources...)
}
