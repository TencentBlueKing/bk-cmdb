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

	uniques := make([]ModelUniqueSimplify, 0)
	for _, item := range resp.Data.Info {
		unique := ModelUniqueSimplify{}
		u, err := unique.Parse(item)
		if err != nil {
			blog.Errorf("collectUniqueByUniqueIDs %+v failed, parse unique %+v failed, err: %+v ", uniqueIDs, item, err)
			return nil, fmt.Errorf("parse unique from db data failed, err: %+v", err)
		}
		uniques = append(uniques, *u)
	}
	blog.V(5).Infof("collectUniqueByUniqueIDs result: %+v", uniques)
	return uniques, nil
}

func (am *AuthManager) makeResourceByUnique(ctx context.Context, header http.Header, action meta.Action, uniques ...ModelUniqueSimplify) ([]meta.ResourceAttribute, error) {
	objectIDs := make([]string, 0)
	for _, unique := range uniques {
		objectIDs = append(objectIDs, unique.ObjID)
	}
	objectIDs = util.StrArrayUnique(objectIDs)
	if len(objectIDs) > 1 {
		blog.Errorf("makeResourceByUnique failed, model uniques belongs to multiple models: %+v", objectIDs)
		return nil, fmt.Errorf("model uniques belongs to multiple models: %+v", objectIDs)
	}
	if len(objectIDs) == 0 {
		blog.Errorf("makeResourceByUnique failed, model id not found, uniques: %+v", uniques)
		return nil, fmt.Errorf("model id not found")
	}

	objects, err := am.collectObjectsByObjectIDs(ctx, header, objectIDs...)
	if err != nil {
		blog.Errorf("makeResourceByUnique failed, collectObjectsByObjectIDs failed, objectIDs: %+v, err: %+v", objectIDs, err)
		return nil, fmt.Errorf("collect object id failed, err: %+v", err)
	}
	
	if len(objects) == 0 {
		blog.Errorf("makeResourceByUnique failed, collectObjectsByObjectIDs no objects found, objectIDs: %+v, err: %+v", objectIDs, err)
		return nil, fmt.Errorf("collect object by id not found")
	}

	businessID, err := am.ExtractBusinessIDFromObjects(objects...)
	if err != nil {
		blog.Errorf("makeResourceByUnique failed, extract business id failed, uniques: %+v, err: %+v", uniques, err)
		return nil, fmt.Errorf("extract business id failed, err: %+v", err)
	}
	
	parentResources, err := am.MakeResourcesByObjects(ctx, header, meta.EmptyAction, objects...)
	if err != nil {
		blog.Errorf("makeResourceByUnique failed, get parent resource failed, objects: %+v, err: %+v", objects, err)
		return nil, fmt.Errorf("get parent resources failed, objects: %+v, err: %+v", objects, err)
	}
	if len(parentResources) > 1 {
		blog.Errorf("makeResourceByUnique failed, get multiple parent resource, parent resources: %+v", parentResources)
		return nil, fmt.Errorf("get multiple parent resources, parent resources: %+v", parentResources)
	}
	if len(parentResources) == 0 {
		blog.Errorf("makeResourceByUnique failed, get parent resource empty, objects: %+v", objects)
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
	blog.V(5).Infof("makeResourceByUnique result: %+v", resources)
	return resources, nil
}

func (am *AuthManager) RegisterModelUnique(ctx context.Context, header http.Header, uniques ...ModelUniqueSimplify) error {
	resources, err := am.makeResourceByUnique(ctx, header, meta.EmptyAction, uniques...)
	if err != nil {
		return fmt.Errorf("register model unique failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterModelUnique(ctx context.Context, header http.Header, uniques ...ModelUniqueSimplify) error {
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

func (am *AuthManager) AuthorizeModelUnique(ctx context.Context, header http.Header, action meta.Action, uniques ...ModelUniqueSimplify) error {
	resources, err := am.makeResourceByUnique(ctx, header, action, uniques...)
	if err != nil {
		return fmt.Errorf("authorize model unique failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) UpdateRegisteredModelUnique(ctx context.Context, header http.Header, uniques ...ModelUniqueSimplify) error {
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

func (am *AuthManager) AuthorizeByUnique(ctx context.Context, header http.Header, action meta.Action, uniques ...ModelUniqueSimplify) error {
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
	objects, err := am.collectObjectsByObjectIDs(ctx, header, objectID)
	if err != nil {
		blog.Errorf("AuthorizeModelUniqueResourceCreate failed, get model by id: %s failed, err: %+v", objectID, err)
		return fmt.Errorf("get model by id %s failed, err: %+v", objectID, err)
	}
	
	parentResources, err := am.MakeResourcesByObjects(ctx, header, meta.Update, objects...)
	if err != nil {
		blog.Errorf("AuthorizeModelUniqueResourceCreate failed, make parent resource by models failed, objects: %+v, err: %+v", objects, err)
		return fmt.Errorf("make parent resource from objects failed, err: %+v", err)
	}
	
	if am.RegisterModelUniqueEnabled == false {
		return am.batchAuthorize(ctx, header, parentResources...)
	}
	
	resources := make([]meta.ResourceAttribute, 0)
	for _, parentResource := range parentResources {
		layers := parentResource.Layers
		layers = append(layers, meta.Item{
			Type:       meta.Model,
			Action:     parentResource.Action,
			Name:       parentResource.Name,
			InstanceID: parentResource.InstanceID,
		})
		resource := meta.ResourceAttribute{
			Basic:           meta.Basic{
				Type:       meta.ModelUnique,
				Action:     meta.Create,
			},
			SupplierAccount: parentResource.SupplierAccount,
			BusinessID:      parentResource.BusinessID,
			Layers:          layers,
		}
		resources = append(resources, resource)
	}
	
	blog.V(5).Infof("AuthorizeModelUniqueResourceCreate result: %+v", resources)
	return am.batchAuthorize(ctx, header, resources...)
}

func (am *AuthManager) RegisterModuleUniqueByID(ctx context.Context, header http.Header, uniqueIDs ...int64) error {
	if am.RegisterModelUniqueEnabled == false {
		return nil
	}
	
	uniques, err := am.collectUniqueByUniqueIDs(ctx, header, uniqueIDs...)
	if err != nil {
		return fmt.Errorf("update registered model unique failed, get unique by id failed, err: %+v", err)
	}
	return am.RegisterModelUnique(ctx, header, uniques...)
}

func (am *AuthManager) AuthorizeModelUniqueByID(ctx context.Context, header http.Header, action meta.Action, ids ...int64) error {
	modelUniques, err := am.collectUniqueByUniqueIDs(ctx, header, ids...)
	if err != nil {
		blog.Errorf("get model unique by id failed, err: %+v", err)
		return fmt.Errorf("get model unique by id failed, err: %+v", err)
	}
	
	return am.AuthorizeByUnique(ctx, header, action, modelUniques...)
}
