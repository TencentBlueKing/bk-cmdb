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

/*
 * instance represent common instances here
 */

func (am *AuthManager) CollectInstancesByModelID(ctx context.Context, header http.Header, objectID string) ([]InstanceSimplify, error) {
	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKObjIDField).Eq(objectID).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDObject, &cond)
	if err != nil {
		blog.V(3).Infof("get instances by model id %s failed, err: %+v", objectID, err)
		return nil, fmt.Errorf("get instances by model id %s failed, err: %+v", objectID, err)
	}
	instances := make([]InstanceSimplify, 0)
	for _, cls := range result.Data.Info {
		instance := InstanceSimplify{}
		_, err = instance.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get instances by object failed, err: %+v", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (am *AuthManager) collectInstancesByInstanceIDs(ctx context.Context, header http.Header, objectID string, instanceIDs ...string) ([]InstanceSimplify, error) {
	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	instanceIDs = util.StrArrayUnique(instanceIDs)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKInstIDField).In(instanceIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, objectID, &cond)
	if err != nil {
		blog.V(5).Infof("collectInstancesByInstanceIDs failed, get instances by id failed, id: %+v, err: %+v", instanceIDs, err)
		return nil, fmt.Errorf("get instances by id failed, id: %+v, err: %+v", instanceIDs, err)
	}
	instances := make([]InstanceSimplify, 0)
	for _, cls := range result.Data.Info {
		instance := InstanceSimplify{}
		_, err = instance.Parse(cls)
		if err != nil {
			blog.V(5).Infof("collectInstancesByInstanceIDs failed, parse instance from db data failed, instance: %+v, err: %+v", cls, err)
			return nil, fmt.Errorf("parse instance from db data failed, err: %+v", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (am *AuthManager) collectInstancesByRawIDs(ctx context.Context, header http.Header, modelID string, ids ...int64) ([]InstanceSimplify, error) {
	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.IntArrayUnique(ids)
	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.GetInstIDField(modelID)).In(ids).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, modelID, &cond)
	if err != nil {
		blog.V(3).Infof("get instance by id failed, err: %+v", err)
		return nil, fmt.Errorf("get instance by id failed, err: %+v", err)
	}
	instances := make([]InstanceSimplify, 0)
	for _, inst := range result.Data.Info {
		instance := InstanceSimplify{}
		_, err = instance.Parse(inst)
		if err != nil {
			return nil, fmt.Errorf("get instance by object failed, err: %+v", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (am *AuthManager) extractBusinessIDFromInstances(instances ...InstanceSimplify) (int64, error) {
	var businessID int64
	for idx, instance := range instances {
		bizID := instance.BizID
		// we should ignore metadata.LabelBusinessID field not found error
		if idx > 0 && bizID != businessID {
			blog.V(5).Infof("extractBusinessIDFromInstances failed, get multiple business ID from instances")
			return 0, fmt.Errorf("get multiple business ID from objects")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) extractObjectIDFromInstances(instances ...InstanceSimplify) (string, error) {
	var objectID string
	for idx, instance := range instances {
		objID := instance.ObjectID
		if idx > 0 && objID != objectID {
			blog.V(5).Infof("extractObjectIDFromInstances failed, get multiple object ID from instances")
			return "", fmt.Errorf("get multiple object ID from instances")
		}
		objectID = objID
	}
	return objectID, nil
}

func (am *AuthManager) MakeResourcesByInstances(ctx context.Context, header http.Header, action meta.Action, businessID int64, instances ...InstanceSimplify) ([]meta.ResourceAttribute, error) {
	objectID, err := am.extractObjectIDFromInstances(instances...)
	if err != nil {
		blog.Errorf("MakeResourcesByInstances failed, extract object id from instances failed, err: %+v", err)
		return nil, fmt.Errorf("extract object id from instances failed, err: %+v", err)
	}

	objects, err := am.collectObjectsByObjectIDs(ctx, header, objectID)
	if err != nil {
		blog.Errorf("MakeResourcesByInstances failed, collect objects by id failed, err: %+v", err)
		return nil, fmt.Errorf("extract object by id failed, err: %+v", err)
	}
	if len(objects) != 1 {
		blog.Errorf("MakeResourcesByInstances failed, get %d objects by id %s", len(objects), objectID)
		return nil, fmt.Errorf("get %d object by id %s", len(objects), objectID)
	}
	object := objects[0]

	parentResources, err := am.MakeResourcesByObjects(ctx, header, meta.EmptyAction, object)
	if err != nil {
		blog.Errorf("MakeResourcesByObjects failed, make parent auth resource by objects failed, object: %+v, err: %+v", object, err)
		return nil, fmt.Errorf("make parent auth resource by objects failed, err: %+v", objectID)
	}
	if len(parentResources) != 1 {
		blog.Errorf("MakeResourcesByInstances failed, make parent auth resource by objects failed, get %d with object %s", len(parentResources), objectID)
		return nil, fmt.Errorf("make parent auth resource by objects failed, get %d with object %s", len(parentResources), objectID)
	}

	parentResource := parentResources[0]

	resources := make([]meta.ResourceAttribute, 0)
	for _, instance := range instances {
		layers := parentResource.Layers
		layers = append(layers, meta.Item{
			Type:       parentResource.Type,
			Action:     parentResource.Action,
			Name:       parentResource.Name,
			InstanceID: parentResource.InstanceID,
		})
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ModelInstance,
				Name:       instance.Name,
				InstanceID: instance.InstanceID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
			Layers:          layers,
		}

		resources = append(resources, resource)
	}
	return resources, nil
}

func (am *AuthManager) AuthorizeByInstances(ctx context.Context, header http.Header, action meta.Action, instances ...InstanceSimplify) error {
	if am.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, instances: %+v", instances)
		return nil
	}

	// extract business id
	bizID, err := am.extractBusinessIDFromInstances(instances...)
	if err != nil {
		blog.V(5).Infof("AuthorizeByInstances failed, extract business ID from instances failed, err: %+v", err)
		return fmt.Errorf("authorize instances failed, extract business id from instance failed, err: %+v", err)
	}

	// make auth resources
	resources, err := am.MakeResourcesByInstances(ctx, header, action, bizID, instances...)
	if err != nil {
		blog.Errorf("AuthorizeByInstances failed, make resource by instances failed, err: %+v", err)
		return fmt.Errorf("make resource by instances failed, err: %+v", err)
	}

	return am.authorize(ctx, header, bizID, resources...)
}

func (am *AuthManager) UpdateRegisteredInstances(ctx context.Context, header http.Header, instances ...InstanceSimplify) error {
	// extract business id
	bizID, err := am.extractBusinessIDFromInstances(instances...)
	if err != nil {
		return fmt.Errorf("deregister instances failed, extract business id from instances failed, err: %+v", err)
	}

	// make auth resources
	resources, err := am.MakeResourcesByInstances(ctx, header, meta.EmptyAction, bizID, instances...)
	if err != nil {
		blog.Errorf("UpdateRegisteredInstances failed, make resource by instances failed, err: %+v", err)
		return fmt.Errorf("make resource by instances failed, err: %+v", err)
	}

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredInstanceByID(ctx context.Context, header http.Header, objectID string, ids ...int64) error {
	instances, err := am.collectInstancesByRawIDs(ctx, header, objectID, ids...)
	if err != nil {
		return fmt.Errorf("update registered instances failed, get instances by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredInstances(ctx, header, instances...)
}

func (am *AuthManager) UpdateRegisteredInstanceByRawID(ctx context.Context, header http.Header, objectID string, ids ...int64) error {
	instances, err := am.collectInstancesByRawIDs(ctx, header, objectID, ids...)
	if err != nil {
		return fmt.Errorf("update registered instances failed, get instances by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredInstances(ctx, header, instances...)
}

func (am *AuthManager) DeregisterInstanceByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	instances, err := am.collectClassificationsByRawIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("deregister instances failed, get instance by id failed, err: %+v", err)
	}
	return am.DeregisterClassification(ctx, header, instances...)
}

func (am *AuthManager) RegisterInstancesByID(ctx context.Context, header http.Header, objectID string, ids ...int64) error {
	instances, err := am.collectInstancesByRawIDs(ctx, header, objectID, ids...)
	if err != nil {
		return fmt.Errorf("register instances failed, get instance by id failed, err: %+v", err)
	}
	return am.RegisterInstances(ctx, header, instances...)
}

func (am *AuthManager) RegisterInstances(ctx context.Context, header http.Header, instances ...InstanceSimplify) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromInstances(instances...)
	if err != nil {
		return fmt.Errorf("register instances failed, extract business id from instances failed, err: %+v", err)
	}

	// make auth resources
	resources, err := am.MakeResourcesByInstances(ctx, header, meta.EmptyAction, bizID, instances...)
	if err != nil {
		blog.Errorf("RegisterInstances failed, make resource by instances failed, err: %+v", err)
		return fmt.Errorf("make resource by instances failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterInstances(ctx context.Context, header http.Header, instances ...InstanceSimplify) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromInstances(instances...)
	if err != nil {
		return fmt.Errorf("deregister instances failed, extract business id from instances failed, err: %+v", err)
	}

	// make auth resources
	resources, err := am.MakeResourcesByInstances(ctx, header, meta.EmptyAction, bizID, instances...)
	if err != nil {
		blog.Errorf("DeregisterInstances failed, make resource by instances failed, err: %+v", err)
		return fmt.Errorf("make resource by instances failed, err: %+v", err)
	}

	return am.Authorize.DeregisterResource(ctx, resources...)
}

// AuthorizeInstanceCreateByObjectID authorize create priority by object, plz be note this method only overlay model read/update/delete, without create
func (am *AuthManager) AuthorizeInstanceCreateByObject(ctx context.Context, header http.Header, action meta.Action, objects ...metadata.Object) error {
	businessID, err := am.ExtractBusinessIDFromObjects(objects...)
	if err != nil {
		blog.V(5).Infof("AuthorizeInstanceCreateByObject failed, extract business id from objects failed, objects: %+v, err: %+v", objects, err)
		return fmt.Errorf("authorize create instance failed, extract business id from models failed, err: %+v", err)
	}

	parentResources, err := am.MakeResourcesByObjects(ctx, header, action, objects...)
	if err != nil {
		blog.V(5).Infof("AuthorizeInstanceCreateByObject failed, make auth resource from objects failed, objects: %+v, err: %+v", objects, err)
		return fmt.Errorf("make parent auth resource by models failed, err: %+v", err)
	}

	resources := make([]meta.ResourceAttribute, 0)
	for _, parentResource := range parentResources {
		layers := parentResource.Layers
		layers = append(layers, meta.Item{
			Type:       parentResource.Basic.Type,
			Action:     parentResource.Basic.Action,
			Name:       parentResource.Basic.Name,
			InstanceID: parentResource.Basic.InstanceID,
		})
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Type:   meta.ModelInstance,
				Action: meta.Create,
			},
			SupplierAccount: parentResource.SupplierAccount,
			BusinessID:      parentResource.BusinessID,
			Layers:          layers,
		}
		resources = append(resources, resource)
	}

	return am.authorize(ctx, header, businessID, resources...)
}
