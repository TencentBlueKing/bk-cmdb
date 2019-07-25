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

func (am *AuthManager) collectUniqueByUniqueIDs(ctx context.Context, header http.Header, uniqueIDs ...int64) ([]ModelUniqueSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	uniqueIDs = util.IntArrayUnique(uniqueIDs)

	// get model by objID
	cond := condition.CreateCondition().Field(common.BKFieldID).In(uniqueIDs)
	queryCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	resp, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameObjUnique, queryCond)
	// resp, err := am.clientSet.CoreService().Model().ReadModelAttrUnique(ctx, header, queryCond)
	if err != nil {
		return nil, fmt.Errorf("get model unique by id: %+v failed, err: %+v", uniqueIDs, err)
	}
	if len(resp.Data.Info) == 0 {
		return nil, fmt.Errorf("get model unique by id: %+v failed, not found", uniqueIDs)
	}
	if len(resp.Data.Info) != len(uniqueIDs) {
		return nil, fmt.Errorf("get model unique by id: %+v failed, get %d, expect %d", uniqueIDs, len(resp.Data.Info), len(uniqueIDs))
	}

	uniques := make([]ModelUniqueSimplify, 0)
	for _, item := range resp.Data.Info {
		unique := ModelUniqueSimplify{}
		u, err := unique.Parse(item)
		if err != nil {
			blog.Errorf("collectUniqueByUniqueIDs %+v failed, parse unique %+v failed, err: %+v, rid: %s", uniqueIDs, item, err, rid)
			return nil, fmt.Errorf("parse unique from db data failed, err: %+v", err)
		}
		uniques = append(uniques, *u)
	}
	blog.V(5).Infof("collectUniqueByUniqueIDs result: %+v", uniques)
	return uniques, nil
}

func (am *AuthManager) ExtractBusinessIDFromUniques(uniques ...ModelUniqueSimplify) (int64, error) {
	if len(uniques) == 0 {
		return 0, fmt.Errorf("no object found")
	}

	businessIDs := make([]int64, 0)
	for _, unique := range uniques {
		businessIDs = append(businessIDs, unique.BusinessID)
	}

	businessIDs = util.IntArrayUnique(businessIDs)

	if len(businessIDs) > 1 {
		return 0, fmt.Errorf("uniques belongs to multiple business: [%+v]", businessIDs)
	}

	if len(businessIDs) == 0 {
		return 0, fmt.Errorf("unexpected error, no business found with uniques: %+v", uniques)
	}
	return businessIDs[0], nil
}

func (am *AuthManager) makeResourceByUnique(ctx context.Context, header http.Header, action meta.Action, uniques ...ModelUniqueSimplify) ([]meta.ResourceAttribute, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	businessID, err := am.ExtractBusinessIDFromUniques(uniques...)
	if err != nil {
		blog.Errorf("makeResourceByUnique failed, ExtractBusinessIDFromUniques failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("extract business id from model uniques failed, err: %+v", err)
	}
	objectIDs := make([]string, 0)
	for _, unique := range uniques {
		objectIDs = append(objectIDs, unique.ObjID)
	}
	objectIDs = util.StrArrayUnique(objectIDs)
	if len(objectIDs) > 1 {
		blog.Errorf("makeResourceByUnique failed, model uniques belongs to multiple models: %+v, rid: %s", objectIDs, rid)
		return nil, fmt.Errorf("model uniques belongs to multiple models: %+v", objectIDs)
	}
	if len(objectIDs) == 0 {
		blog.Errorf("makeResourceByUnique failed, model id not found, uniques: %+v", uniques)
		return nil, fmt.Errorf("model id not found")
	}

	objects, err := am.collectObjectsByObjectIDs(ctx, header, businessID, objectIDs...)
	if err != nil {
		blog.Errorf("makeResourceByUnique failed, collectObjectsByObjectIDs failed, objectIDs: %+v, err: %+v, rid: %s", objectIDs, err, rid)
		return nil, fmt.Errorf("collect object id failed, err: %+v", err)
	}

	if len(objects) == 0 {
		blog.Errorf("makeResourceByUnique failed, collectObjectsByObjectIDs no objects found, objectIDs: %+v, err: %+v, rid: %s", objectIDs, err, rid)
		return nil, fmt.Errorf("collect object by id not found")
	}

	parentResources, err := am.MakeResourcesByObjects(ctx, header, meta.EmptyAction, objects...)
	if err != nil {
		blog.Errorf("makeResourceByUnique failed, get parent resource failed, objects: %+v, err: %+v, rid: %s", objects, err, rid)
		return nil, fmt.Errorf("get parent resources failed, objects: %+v, err: %+v", objects, err)
	}
	if len(parentResources) > 1 {
		blog.Errorf("makeResourceByUnique failed, get multiple parent resource, parent resources: %+v, rid: %s", parentResources, rid)
		return nil, fmt.Errorf("get multiple parent resources, parent resources: %+v", parentResources)
	}
	if len(parentResources) == 0 {
		blog.Errorf("makeResourceByUnique failed, get parent resource empty, objects: %+v, rid: %s", objects, rid)
		return nil, fmt.Errorf("get parent resources empty, objects: %+v", objects)
	}

	parentResource := parentResources[0]

	// prepare resource layers for authorization
	resources := make([]meta.ResourceAttribute, 0)
	for _, unique := range uniques {
		// model
		parentLayers := parentResource.Layers
		parentLayers = append(parentLayers, meta.Item{
			Type:       meta.Model,
			Name:       parentResource.Name,
			InstanceID: parentResource.InstanceID,
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
			Layers:          parentLayers,
		}

		resources = append(resources, resource)
	}
	blog.V(5).Infof("makeResourceByUnique result: %+v, rid: %s", resources, rid)
	return resources, nil
}

func (am *AuthManager) RegisterModelUnique(ctx context.Context, header http.Header, uniques ...ModelUniqueSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(uniques) == 0 {
		return nil
	}
	if am.RegisterModelUniqueEnabled == false {
		return nil
	}

	resources, err := am.makeResourceByUnique(ctx, header, meta.EmptyAction, uniques...)
	if err != nil {
		return fmt.Errorf("register model unique failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterModelUnique(ctx context.Context, header http.Header, uniques ...ModelUniqueSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(uniques) == 0 {
		return nil
	}
	if am.RegisterModelUniqueEnabled == false {
		return nil
	}

	resources, err := am.makeResourceByUnique(ctx, header, meta.EmptyAction, uniques...)
	if err != nil {
		return fmt.Errorf("deregister model unique failed, err: %+v", err)
	}

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterModelUniqueByID(ctx context.Context, header http.Header, uniqueIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(uniqueIDs) == 0 {
		return nil
	}
	if am.RegisterModelUniqueEnabled == false {
		return nil
	}

	uniques, err := am.collectUniqueByUniqueIDs(ctx, header, uniqueIDs...)
	if err != nil {
		return fmt.Errorf("update registered model unique failed, get unique by id failed, err: %+v", err)
	}
	return am.DeregisterModelUnique(ctx, header, uniques...)
}

// func (am *AuthManager) AuthorizeModelUnique(ctx context.Context, header http.Header, action meta.Action, uniques ...ModelUniqueSimplify) error {
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	resources, err := am.makeResourceByUnique(ctx, header, action, uniques...)
// 	if err != nil {
// 		return fmt.Errorf("authorize model unique failed, err: %+v", err)
// 	}
//
// 	return am.batchAuthorize(ctx, header, resources...)
// }

func (am *AuthManager) UpdateRegisteredModelUnique(ctx context.Context, header http.Header, uniques ...ModelUniqueSimplify) error {
	if am.Enabled() == false {
		return nil
	}

	if len(uniques) == 0 {
		return nil
	}
	if am.RegisterModelUniqueEnabled == false {
		return nil
	}

	resources, err := am.makeResourceByUnique(ctx, header, meta.EmptyAction, uniques...)
	if err != nil {
		return fmt.Errorf("update registered model unique failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) UpdateRegisteredModelUniqueByID(ctx context.Context, header http.Header, uniqueIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(uniqueIDs) == 0 {
		return nil
	}
	if am.RegisterModelUniqueEnabled == false {
		return nil
	}

	uniques, err := am.collectUniqueByUniqueIDs(ctx, header, uniqueIDs...)
	if err != nil {
		return fmt.Errorf("update registered model unique failed, get unique by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredModelUnique(ctx, header, uniques...)
}

// func (am *AuthManager) AuthorizeByUniqueID(ctx context.Context, header http.Header, action meta.Action, uniqueIDs ...int64) error {
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	uniques, err := am.collectUniqueByUniqueIDs(ctx, header, uniqueIDs...)
// 	if err != nil {
// 		return fmt.Errorf("get uniques by id failed, err: %+v", err)
// 	}
//
// 	businessID, err := am.ExtractBusinessIDFromUniques(uniques...)
// 	if err != nil {
// 		return fmt.Errorf("extract business id by uniques failed, err: %+v", err)
// 	}
//
// 	objectIDs := make([]string, 0)
// 	for _, unique := range uniques {
// 		objectIDs = append(objectIDs, unique.ObjID)
// 	}
//
// 	return am.AuthorizeByObjectID(ctx, header, action, businessID, objectIDs...)
// }

// func (am *AuthManager) AuthorizeByUnique(ctx context.Context, header http.Header, action meta.Action, uniques ...ModelUniqueSimplify) error {
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	bizID, err := am.ExtractBusinessIDFromUniques(uniques...)
// 	if err != nil {
// 		return fmt.Errorf("extract business id from model uniques failed, err: %+v", err)
// 	}
//
// 	objectIDs := make([]string, 0)
// 	for _, unique := range uniques {
// 		objectIDs = append(objectIDs, unique.ObjID)
// 	}
//
// 	if am.RegisterModelUniqueEnabled == false {
// 		objectAction := meta.Update
// 		if action == meta.Find || action == meta.FindMany {
// 			objectAction = action
// 		}
// 		return am.AuthorizeByObjectID(ctx, header, objectAction, bizID, objectIDs...)
// 	}
//
// 	resources, err := am.makeResourceByUnique(ctx, header, action, uniques...)
// 	if err != nil {
// 		return fmt.Errorf("make auth resource from model uniques failed, err: %+v", err)
// 	}
//
// 	return am.authorize(ctx, header, bizID, resources...)
// }

// func (am *AuthManager) AuthorizeModelUniqueResourceCreate(ctx context.Context, header http.Header, businessID int64, objectID string) error {
// 	rid := util.ExtractRequestIDFromContext(ctx)
//
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	objects, err := am.collectObjectsByObjectIDs(ctx, header, businessID, objectID)
// 	if err != nil {
// 		blog.Errorf("AuthorizeModelUniqueResourceCreate failed, get model by id: %s failed, err: %+v, rid: %s", objectID, err, rid)
// 		return fmt.Errorf("get model by id %s failed, err: %+v", objectID, err)
// 	}
//
// 	parentResources, err := am.MakeResourcesByObjects(ctx, header, meta.Update, objects...)
// 	if err != nil {
// 		blog.Errorf("AuthorizeModelUniqueResourceCreate failed, make parent resource by models failed, objects: %+v, err: %+v, rid: %s", objects, err, rid)
// 		return fmt.Errorf("make parent resource from objects failed, err: %+v", err)
// 	}
//
// 	if am.RegisterModelUniqueEnabled == false {
// 		return am.batchAuthorize(ctx, header, parentResources...)
// 	}
//
// 	resources := make([]meta.ResourceAttribute, 0)
// 	for _, parentResource := range parentResources {
// 		layers := parentResource.Layers
// 		layers = append(layers, meta.Item{
// 			Type:       meta.Model,
// 			Action:     parentResource.Action,
// 			Name:       parentResource.Name,
// 			InstanceID: parentResource.InstanceID,
// 		})
// 		resource := meta.ResourceAttribute{
// 			Basic: meta.Basic{
// 				Type:   meta.ModelUnique,
// 				Action: meta.Create,
// 			},
// 			SupplierAccount: parentResource.SupplierAccount,
// 			BusinessID:      businessID,
// 			Layers:          layers,
// 		}
// 		resources = append(resources, resource)
// 	}
//
// 	blog.V(5).Infof("AuthorizeModelUniqueResourceCreate result: %+v, rid: %s", resources, rid)
// 	return am.batchAuthorize(ctx, header, resources...)
// }

func (am *AuthManager) RegisterModuleUniqueByID(ctx context.Context, header http.Header, uniqueIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(uniqueIDs) == 0 {
		return nil
	}
	if am.RegisterModelUniqueEnabled == false {
		return nil
	}

	uniques, err := am.collectUniqueByUniqueIDs(ctx, header, uniqueIDs...)
	if err != nil {
		return fmt.Errorf("update registered model unique failed, get unique by id failed, err: %+v", err)
	}
	return am.RegisterModelUnique(ctx, header, uniques...)
}

// func (am *AuthManager) AuthorizeModelUniqueByID(ctx context.Context, header http.Header, action meta.Action, ids ...int64) error {
// 	rid := util.ExtractRequestIDFromContext(ctx)
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	modelUniques, err := am.collectUniqueByUniqueIDs(ctx, header, ids...)
// 	if err != nil {
// 		blog.Errorf("get model unique by id failed, err: %+v, rid: %s", err, rid)
// 		return fmt.Errorf("get model unique by id failed, err: %+v", err)
// 	}
//
// 	return am.AuthorizeByUnique(ctx, header, action, modelUniques...)
// }
