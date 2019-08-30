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

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
)

/*
 * model classification which used for manage models as group
 */

func (am *AuthManager) CollectClassificationByBusinessIDs(ctx context.Context, header http.Header, businessID int64) ([]metadata.Classification, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	condCheckModel := mongo.NewCondition()
	if businessID != 0 {
		_, metaCond := condCheckModel.Embed(metadata.BKMetadata)
		_, labelCond := metaCond.Embed(metadata.BKLabel)
		labelCond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: strconv.FormatInt(businessID, 10)})
	}
	cond := condCheckModel.ToMapStr()
	if businessID == 0 {
		cond.Merge(metadata.BizLabelNotExist)
	}
	query := &metadata.QueryCondition{
		Condition: cond,
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameObjClassifiction, query)
	if err != nil {
		blog.Errorf("get module:%+v by businessID:%d failed, err: %+v, rid: %s", businessID, err, rid)
		return nil, fmt.Errorf("get module by businessID:%d failed, err: %+v", businessID, err)
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

func (am *AuthManager) collectClassificationsByClassificationIDs(ctx context.Context, header http.Header, classificationIDs ...string) ([]metadata.Classification, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	classificationIDs = util.StrArrayUnique(classificationIDs)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKClassificationIDField).In(classificationIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameObjClassifiction, &cond)
	if err != nil {
		blog.V(3).Infof("get classification by id failed, err: %+v, rid: %s", err, rid)
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
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.IntArrayUnique(ids)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKFieldID).In(ids).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameObjClassifiction, &cond)
	if err != nil {
		blog.V(3).Infof("get classification by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get classification by id failed, err: %+v", err)
	}
	classifications := make([]metadata.Classification, 0)
	for _, cls := range result.Data.Info {
		classification := metadata.Classification{}
		_, err = classification.Parse(cls)
		if err != nil {
			blog.Errorf("collectClassificationsByRawIDs %+v failed, parse classification %+v failed, err: %+v, rid: %s", ids, cls, err, rid)
			return nil, fmt.Errorf("parse classification from db data failed, err: %+v", err)
		}
		classifications = append(classifications, classification)
	}
	return classifications, nil
}

func (am *AuthManager) extractBusinessIDFromClassifications(ctx context.Context, classifications ...metadata.Classification) (int64, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	var businessID int64
	var err error
	businessIDs := make([]int64, 0)
	for _, classification := range classifications {
		businessID, err = extractBusinessID(classification.Metadata.Label)
		// we should ignore metadata.LabelBusinessID field not found error
		if err != nil {
			return 0, fmt.Errorf("parse biz id from classification: %+v failed, err: %+v", classification, err)
		}
		businessIDs = append(businessIDs, businessID)
	}
	businessIDs = util.IntArrayUnique(businessIDs)
	if len(businessIDs) > 1 {
		blog.Errorf("extractBusinessIDFromClassifications failed, get multiple business from classifications, business: %+v, rid: %s", businessIDs, rid)
		return 0, fmt.Errorf("get multiple business from classifictions, business: %+v", businessIDs)
	}
	return businessID, nil
}

func (am *AuthManager) MakeResourcesByClassifications(header http.Header, action meta.Action, businessID int64, classifications ...metadata.Classification) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, classification := range classifications {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ModelClassification,
				Name:       classification.ClassificationName,
				InstanceID: classification.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

// func (am *AuthManager) AuthorizeByClassification(ctx context.Context, header http.Header, action meta.Action, classifications ...metadata.Classification) error {
// 	rid := util.ExtractRequestIDFromContext(ctx)
//
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	if am.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
// 		blog.V(4).Infof("skip authorization for reading, classifications: %+v, rid: %s", classifications, rid)
// 		return nil
// 	}
//
// 	// extract business id
// 	bizID, err := am.extractBusinessIDFromClassifications(ctx, classifications...)
// 	if err != nil {
// 		return fmt.Errorf("authorize classifications failed, extract business id from classification failed, err: %+v", err)
// 	}
//
// 	// make auth resources
// 	resources := am.MakeResourcesByClassifications(header, action, bizID, classifications...)
//
// 	return am.authorize(ctx, header, bizID, resources...)
// }

func (am *AuthManager) UpdateRegisteredClassification(ctx context.Context, header http.Header, classifications ...metadata.Classification) error {
	if am.Enabled() == false {
		return nil
	}

	if len(classifications) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromClassifications(ctx, classifications...)
	if err != nil {
		return fmt.Errorf("authorize classifications failed, extract business id from classification failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByClassifications(header, meta.EmptyAction, bizID, classifications...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredClassificationByID(ctx context.Context, header http.Header, classificationIDs ...string) error {
	if am.Enabled() == false {
		return nil
	}

	if len(classificationIDs) == 0 {
		return nil
	}

	classifications, err := am.collectClassificationsByClassificationIDs(ctx, header, classificationIDs...)
	if err != nil {
		return fmt.Errorf("update registered classifications failed, get classfication by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredClassification(ctx, header, classifications...)
}

func (am *AuthManager) UpdateRegisteredClassificationByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	classifications, err := am.collectClassificationsByRawIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered classifications failed, get classfication by raw id failed, err: %+v", err)
	}
	return am.UpdateRegisteredClassification(ctx, header, classifications...)
}

func (am *AuthManager) RegisterClassification(ctx context.Context, header http.Header, classifications ...metadata.Classification) error {
	if am.Enabled() == false {
		return nil
	}

	if len(classifications) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromClassifications(ctx, classifications...)
	if err != nil {
		return fmt.Errorf("register classifications failed, extract business id from classification failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByClassifications(header, meta.EmptyAction, bizID, classifications...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterClassification(ctx context.Context, header http.Header, classifications ...metadata.Classification) error {
	if am.Enabled() == false {
		return nil
	}

	if len(classifications) == 0 {
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromClassifications(ctx, classifications...)
	if err != nil {
		return fmt.Errorf("deregister classifications failed, extract business id from classification failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByClassifications(header, meta.EmptyAction, bizID, classifications...)

	return am.Authorize.DeregisterResource(ctx, resources...)
}
