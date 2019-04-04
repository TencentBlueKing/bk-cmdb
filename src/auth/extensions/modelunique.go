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
	"configcenter/src/common/util"
)

func (am *AuthManager) collectUniqueByUniqueIDs(ctx context.Context, header http.Header, uniqueIDs ...int64) ([]metadata.ObjectUnique, error) {
	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	uniqueIDs = util.IntArrayUnique(uniqueIDs)

	// get model by objID
	cond := condition.CreateCondition().Field(common.BKFieldID).In(uniqueIDs)
	queryCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	resp, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameObjUnique, queryCond)
	if err != nil {
		return nil, fmt.Errorf("get model unique by id: %+v failed, err: %+v", uniqueIDs, err)
	}
	if len(resp.Data.Info) == 0 {
		return nil, fmt.Errorf("get model unique by id: %+v failed, not found", uniqueIDs)
	}
	if len(resp.Data.Info) != len(uniqueIDs) {
		return nil, fmt.Errorf("get model unique by id: %+v failed, get %d, expect %d", uniqueIDs, len(resp.Data.Info), len(uniqueIDs))
	}

	uniques := make([]metadata.ObjectUnique, 0)
	for _, item := range resp.Data.Info {
		unique := metadata.ObjectUnique{}
		_, err := unique.Parse(item)
		if err != nil {
			blog.Errorf("collectUniqueByUniqueIDs %+v failed, parse unique %+v failed, err: %+v ", uniqueIDs, item, err)
			return nil, fmt.Errorf("parse unique from db data failed, err: %+v", err)
		}
		uniques = append(uniques, unique)
	}
	return uniques, nil
}

func (am *AuthManager) makeResourceByUnique(ctx context.Context, header http.Header, action meta.Action, uniques ...metadata.ObjectUnique) ([]meta.ResourceAttribute, error) {
	objectIDs := make([]string, 0)
	for _, unique := range uniques {
		objectIDs = append(objectIDs, unique.ObjID)
	}

	objects, err := am.collectObjectsByObjectIDs(ctx, header, objectIDs...)
	if err != nil {
		return nil, fmt.Errorf("register model unique failed, get related models failed, err: %+v", err)
	}
	objectMap := map[string]metadata.Object{}
	for _, object := range objects {
		objectMap[object.ObjectID] = object
	}

	businessID, err := am.ExtractBusinessIDFromObjects(objects...)
	if err != nil {
		return nil, fmt.Errorf("make auth resource for model unique failed, err: %+v", err)
	}

	classificationIDs := make([]string, 0)
	for _, object := range objects {
		classificationIDs = append(classificationIDs, object.ObjCls)
	}
	classifications, err := am.collectClassificationsByClassificationIDs(ctx, header, classificationIDs...)
	if err != nil {
		return nil, fmt.Errorf("register model unique failed, get related models failed, err: %+v", err)
	}
	classificationMap := map[string]metadata.Classification{}
	for _, classification := range classifications {
		classificationMap[classification.ClassificationID] = classification
	}

	// step2 prepare resource layers for authorization
	resources := make([]meta.ResourceAttribute, 0)
	for _, unique := range uniques {

		object := objectMap[unique.ObjID]

		// check obj's group id in map
		if _, exist := classificationMap[object.ObjCls]; exist == false {
			blog.V(3).Infof("authorization failed, get classification by object failed, err: bk_classification_id not exist")
			return nil, fmt.Errorf("authorization failed, get classification by object failed, err: bk_classification_id not exist")
		}

		parentLayers := meta.Layers{}
		// model group
		parentLayers = append(parentLayers, meta.Item{
			Type:       meta.Model,
			Name:       classificationMap[object.ObjCls].ClassificationID,
			InstanceID: classificationMap[object.ObjCls].ID,
		})

		// model
		parentLayers = append(parentLayers, meta.Item{
			Type:       meta.Model,
			Name:       object.ObjectID,
			InstanceID: object.ID,
		})

		// unique
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ModelUnique,
				Name:       strconv.FormatUint(unique.ID, 10),
				InstanceID: int64(unique.ID),
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return nil, nil
}

func (am *AuthManager) RegisterModelUnique(ctx context.Context, header http.Header, uniques ...metadata.ObjectUnique) error {
	resources, err := am.makeResourceByUnique(ctx, header, meta.EmptyAction, uniques...)
	if err != nil {
		return fmt.Errorf("register model unique failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterModelUnique(ctx context.Context, header http.Header, uniques ...metadata.ObjectUnique) error {
	resources, err := am.makeResourceByUnique(ctx, header, meta.EmptyAction, uniques...)
	if err != nil {
		return fmt.Errorf("deregister model unique failed, err: %+v", err)
	}

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterModelUniqueByID(ctx context.Context, header http.Header, uniqueIDs ...int64) error {
	uniques, err := am.collectUniqueByUniqueIDs(ctx, header, uniqueIDs...)
	if err != nil {
		return fmt.Errorf("update registered model unique failed, get unique by id failed, err: %+v", err)
	}
	return am.DeregisterModelUnique(ctx, header, uniques...)
}

func (am *AuthManager) AuthorizeModelUnique(ctx context.Context, header http.Header, action meta.Action, uniques ...metadata.ObjectUnique) error {
	resources, err := am.makeResourceByUnique(ctx, header, action, uniques...)
	if err != nil {
		return fmt.Errorf("authorize model unique failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) UpdateRegisteredModelUnique(ctx context.Context, header http.Header, uniques ...metadata.ObjectUnique) error {
	resources, err := am.makeResourceByUnique(ctx, header, meta.EmptyAction, uniques...)
	if err != nil {
		return fmt.Errorf("update registered model unique failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) UpdateRegisteredModelUniqueByID(ctx context.Context, header http.Header, uniqueIDs ...int64) error {
	uniques, err := am.collectUniqueByUniqueIDs(ctx, header, uniqueIDs...)
	if err != nil {
		return fmt.Errorf("update registered model unique failed, get unique by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredModelUnique(ctx, header, uniques...)
}

func (am *AuthManager) AuthorizeByUniqueID(ctx context.Context, header http.Header, action meta.Action, uniqueIDs ...int64) error {
	uniques, err := am.collectUniqueByUniqueIDs(ctx, header, uniqueIDs...)
	if err != nil {
		return fmt.Errorf("get uniques by id failed, err: %+v", err)
	}

	objectIDs := make([]string, 0)
	for _, unique := range uniques {
		objectIDs = append(objectIDs, unique.ObjID)
	}

	return am.AuthorizeByObjectID(ctx, header, action, objectIDs...)
}

func (am *AuthManager) AuthorizeByUnique(ctx context.Context, header http.Header, action meta.Action, uniques ...metadata.ObjectUnique) error {
	objectIDs := make([]string, 0)
	for _, unique := range uniques {
		objectIDs = append(objectIDs, unique.ObjID)
	}

	if am.RegisterModelUniqueEnabled == false {
		return am.AuthorizeByObjectID(ctx, header, meta.Update, objectIDs...)
	}

	objects, err := am.collectObjectsByObjectIDs(ctx, header, objectIDs...)
	if err != nil {
		return fmt.Errorf("get model by id failed, err: %+v", err)
	}

	bizID, err := am.ExtractBusinessIDFromObjects(objects...)
	if err != nil {
		return fmt.Errorf("extract business id from model failed, err: %+v", err)
	}

	resources, err := am.makeResourceByUnique(ctx, header, action, uniques...)
	if err != nil {
		return fmt.Errorf("make auth resource from model uniques failed, err: %+v", err)
	}

	return am.authorize(ctx, header, bizID, resources...)
}

func (am *AuthManager) AuthorizeModelUniqueResourceCreate(ctx context.Context, header http.Header, objectID string) error {
	return am.AuthorizeByObjectID(ctx, header, meta.Update, objectID)
}

func (am *AuthManager) RegisterModuleUniqueByID(ctx context.Context, header http.Header, uniqueIDs ...int64) error {
	uniques, err := am.collectUniqueByUniqueIDs(ctx, header, uniqueIDs...)
	if err != nil {
		return fmt.Errorf("update registered model unique failed, get unique by id failed, err: %+v", err)
	}
	return am.RegisterModelUnique(ctx, header, uniques...)
}
