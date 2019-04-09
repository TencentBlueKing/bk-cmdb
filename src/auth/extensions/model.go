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

func (am *AuthManager) CollectObjectsByBusinessID(ctx context.Context, header http.Header, businessID int64) ([]metadata.Object, error) {
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
	models, err := am.clientSet.CoreService().Model().ReadModel(context.Background(), header, query)
	if err != nil {
		blog.Errorf("get models by business %d failed, err: %+v", businessID, err)
		return nil, fmt.Errorf("get models by business %d failed, err: %+v", businessID, err)
	}

	objects := make([]metadata.Object, 0)
	for _, model := range models.Data.Info {
		objects = append(objects, model.Spec)
	}

	blog.V(4).Infof("list model by business %d result: %+v", businessID, objects)
	return objects, nil
}

func (am *AuthManager) collectObjectsByObjectIDs(ctx context.Context, header http.Header, objIDs ...string) ([]metadata.Object, error) {
	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	objIDs = util.StrArrayUnique(objIDs)

	// get model by objID
	cond := condition.CreateCondition().Field(common.BKObjIDField).In(objIDs)
	queryCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	resp, err := am.clientSet.CoreService().Model().ReadModel(ctx, header, queryCond)
	if err != nil {
		return nil, fmt.Errorf("get model by id: %+v failed, err: %+v", objIDs, err)
	}
	if len(resp.Data.Info) == 0 {
		return nil, fmt.Errorf("get model by id: %+v failed, not found", objIDs)
	}
	if len(resp.Data.Info) != len(objIDs) {
		return nil, fmt.Errorf("get model by id: %+v failed, get multiple model", objIDs)
	}

	objects := make([]metadata.Object, 0)
	for _, item := range resp.Data.Info {
		objects = append(objects, item.Spec)
	}

	return objects, nil
}

func (am *AuthManager) collectObjectsByRawIDs(ctx context.Context, header http.Header, ids ...int64) ([]metadata.Object, error) {
	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.IntArrayUnique(ids)

	// get model by objID
	cond := condition.CreateCondition().Field(common.BKFieldID).In(ids)
	queryCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	resp, err := am.clientSet.CoreService().Model().ReadModel(ctx, header, queryCond)
	if err != nil {
		return nil, fmt.Errorf("get model by id: %+v failed, err: %+v", ids, err)
	}
	if len(resp.Data.Info) == 0 {
		return nil, fmt.Errorf("get model by id: %+v failed, not found", ids)
	}
	if len(resp.Data.Info) != len(ids) {
		return nil, fmt.Errorf("get model by id: %+v failed, result count %d not equal to expect %d", ids, len(resp.Data.Info), len(ids))
	}

	objects := make([]metadata.Object, 0)
	for _, item := range resp.Data.Info {
		objects = append(objects, item.Spec)
	}

	return objects, nil
}

func (am *AuthManager) ExtractBusinessIDFromObject(object metadata.Object) (int64, error) {
	bizID, err := object.Metadata.Label.Int64(metadata.LabelBusinessID)
	// we should ignore metadata.LabelBusinessID field not found error
	if err != nil && err != metadata.LabelKeyNotExistError {
		return 0, fmt.Errorf("parse biz id from model: %+v failed, err: %+v", object, err)
	}
	return bizID, nil
}

func (am *AuthManager) ExtractBusinessIDFromObjects(objects ...metadata.Object) (int64, error) {
	var businessID int64
	for idx, object := range objects {
		bizID, err := am.ExtractBusinessIDFromObject(object)
		if err != nil {
			return 0, fmt.Errorf("parse business id from model failed, model: %+v, err: %+v", object, err)
		}
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("authorization failed, get multiple business ID from objects")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) MakeResourcesByObjects(ctx context.Context, header http.Header, action meta.Action, objects ...metadata.Object) ([]meta.ResourceAttribute, error) {
	// prepare resource layers for authorization
	resources := make([]meta.ResourceAttribute, 0)
	for _, object := range objects {
		businessID, err := am.ExtractBusinessIDFromObject(object)
		if err != nil {
			blog.V(3).Infof("parse business id from object failed, err: %+v", err)
			return nil, fmt.Errorf("parse business id from object failed, err: %+v", err)
		}

		// instance
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Model,
				Name:       object.ObjectName,
				InstanceID: object.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// AuthorizeByObjectID authorize model by id
func (am *AuthManager) AuthorizeByObjectID(ctx context.Context, header http.Header, action meta.Action, objIDs ...string) error {
	if len(objIDs) == 0 {
		return nil
	}
	if am.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, models: %+v", objIDs)
		return nil
	}

	objects, err := am.collectObjectsByObjectIDs(ctx, header, objIDs...)
	if err != nil {
		return fmt.Errorf("get model by id failed, err: %+v", err)
	}

	return am.AuthorizeByObject(ctx, header, action, objects...)
}

// AuthorizeObject authorize by object, plz be note this method only overlay model read/update/delete, without create
func (am *AuthManager) AuthorizeByObject(ctx context.Context, header http.Header, action meta.Action, objects ...metadata.Object) error {
	if am.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany || action == meta.ModelTopologyView) {
		blog.V(4).Infof("skip authorization for reading, models: %+v", objects)
		return nil
	}

	// make resources from objects
	resources, err := am.MakeResourcesByObjects(ctx, header, action, objects...)
	if err != nil {
		return fmt.Errorf("make auth resource by models failed, err: %+v", err)
	}

	return am.batchAuthorize(ctx, header, resources...)
}

// AuthorizeObject authorize by object, plz be note this method only overlay model read/update/delete, without create
func (am *AuthManager) AuthorizeResourceCreateByObject(ctx context.Context, header http.Header, action meta.Action, objects ...metadata.Object) error {
	// step1: extract business ID from object, business ID from all objects must be identical to one value
	businessID, err := am.ExtractBusinessIDFromObjects(objects...)
	if err != nil {
		return fmt.Errorf("authrize create instance failed, extract business id from models failed, err: %+v", err)
	}

	resources, err := am.MakeResourcesByObjects(ctx, header, action, objects...)
	if err != nil {
		return fmt.Errorf("make auth resource by models failed, err: %+v", err)
	}

	return am.authorize(ctx, header, businessID, resources...)
}

func (am *AuthManager) AuthorizeResourceCreate(ctx context.Context, header http.Header, businessID int64, resourceType meta.ResourceType) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:   resourceType,
			Action: meta.Create,
		},
		SupplierAccount: util.GetOwnerID(header),
		BusinessID:      businessID,
	}

	return am.authorize(ctx, header, businessID, resource)
}

func (am *AuthManager) RegisterObject(ctx context.Context, header http.Header, objects ...metadata.Object) error {
	if len(objects) == 0 {
		return nil
	}
	resources, err := am.MakeResourcesByObjects(ctx, header, meta.EmptyAction, objects...)
	if err != nil {
		return fmt.Errorf("make auth resource by models failed, err: %+v", err)
	}

	if err := am.Authorize.RegisterResource(ctx, resources...); err != nil {
		return fmt.Errorf("deregister models failed, err: %+v", err)
	}
	return nil
}

func (am *AuthManager) UpdateRegisteredObjects(ctx context.Context, header http.Header, objects ...metadata.Object) error {
	if len(objects) == 0 {
		return nil
	}
	resources, err := am.MakeResourcesByObjects(ctx, header, meta.EmptyAction, objects...)
	if err != nil {
		return fmt.Errorf("make auth resource by models failed, err: %+v", err)
	}

	if err := am.updateResources(ctx, resources...); err != nil {
		return fmt.Errorf("deregister models failed, err: %+v", err)
	}
	return nil
}
func (am *AuthManager) UpdateRegisteredObjectsByRawIDs(ctx context.Context, header http.Header, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}
	ids = util.IntArrayUnique(ids)

	objects, err := am.collectObjectsByRawIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("get model by id failed, id: %+v, err: %+v", ids, err)
	}

	return am.UpdateRegisteredObjects(ctx, header, objects...)
}
func (am *AuthManager) DeregisterObject(ctx context.Context, header http.Header, objects ...metadata.Object) error {
	if len(objects) == 0 {
		return nil
	}
	resources, err := am.MakeResourcesByObjects(ctx, header, meta.EmptyAction, objects...)
	if err != nil {
		return fmt.Errorf("make auth resource by models failed, err: %+v", err)
	}

	if err := am.Authorize.DeregisterResource(ctx, resources...); err != nil {
		return fmt.Errorf("deregister models failed, err: %+v", err)
	}
	return nil
}

func (am *AuthManager) RegisterMainlineObject(ctx context.Context, header http.Header, objects ...metadata.Object) error {
	return am.RegisterObject(ctx, header, objects...)
}

func (am *AuthManager) DeregisterMainlineModelByObjectID(ctx context.Context, header http.Header, objectIDs ...string) error {
	if len(objectIDs) == 0 {
		return nil
	}
	objects, err := am.collectObjectsByObjectIDs(ctx, header, objectIDs...)
	if err != nil {
		return fmt.Errorf("deregister mainline model failed, get model by id failed, err: %+v", err)
	}
	return am.DeregisterObject(ctx, header, objects...)
}
