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
	"errors"
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
	rid := util.ExtractRequestIDFromContext(ctx)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKObjIDField).Eq(objectID).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, objectID, &cond)
	if err != nil {
		blog.V(3).Infof("get instances by model id %s failed, err: %+v, rid: %s", objectID, err, rid)
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
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	instanceIDs = util.StrArrayUnique(instanceIDs)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKInstIDField).In(instanceIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, objectID, &cond)
	if err != nil {
		blog.V(5).Infof("collectInstancesByInstanceIDs failed, get instances by id failed, id: %+v, err: %+v, rid: %s", instanceIDs, err, rid)
		return nil, fmt.Errorf("get instances by id failed, id: %+v, err: %+v", instanceIDs, err)
	}
	instances := make([]InstanceSimplify, 0)
	for _, cls := range result.Data.Info {
		instance := InstanceSimplify{}
		_, err = instance.Parse(cls)
		if err != nil {
			blog.V(5).Infof("collectInstancesByInstanceIDs failed, parse instance from db data failed, instance: %+v, err: %+v, rid: %s", cls, err, rid)
			return nil, fmt.Errorf("parse instance from db data failed, err: %+v", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (am *AuthManager) collectInstancesByRawIDs(ctx context.Context, header http.Header, modelID string, ids ...int64) ([]InstanceSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.IntArrayUnique(ids)
	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.GetInstIDField(modelID)).In(ids).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, modelID, &cond)
	if err != nil {
		blog.V(3).Infof("get instance by id failed, err: %+v, rid: %s", err, rid)
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

func (am *AuthManager) extractBusinessIDFromInstances(ctx context.Context, instances ...InstanceSimplify) (map[int64]int64, error) {
	businessIDMap := make(map[int64]int64)
	if len(instances) == 0 {
		return businessIDMap, fmt.Errorf("empty instances")
	}
	for _, instance := range instances {
		businessIDMap[instance.InstanceID] = instance.BizID
	}
	return businessIDMap, nil
}

// collectObjectsByInstances collect all instances's related model, group by map
// it support cross multiple business and objects
func (am *AuthManager) collectObjectsByInstances(ctx context.Context, header http.Header, instances ...InstanceSimplify) (map[int64]metadata.Object, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// construct parameters for querying models
	businessObjectsMap := map[int64][]string{}
	for _, instance := range instances {
		businessObjectsMap[instance.BizID] = append(businessObjectsMap[instance.BizID], instance.ObjectID)
	}

	// get models group by business
	bizIDObjID2ObjMap := map[int64]map[string]metadata.Object{}
	for businessID, objectIDs := range businessObjectsMap {
		objectIDs = util.StrArrayUnique(objectIDs)
		objects, err := am.collectObjectsByObjectIDs(ctx, header, businessID, objectIDs...)
		if err != nil {
			blog.Errorf("extractObjectIDFromInstances failed, get models by businessID and object id failed, bizID: %+v, objectIDs: %+v, err: %+v, rid: %s", businessID, objectIDs, err, rid)
			return nil, fmt.Errorf("get models by objectIDs and business id failed")
		}
		if len(objectIDs) != len(objects) {
			blog.Errorf("extractObjectIDFromInstances failed, get models by object id failed, input len %d and output len %d not equal, input: %+v, output: %+v, businessID: %d", len(objectIDs), len(objects), objectIDs, objects, businessID)
			return nil, fmt.Errorf("unexpect error, some models maybe not found")
		}
		if bizIDObjID2ObjMap[businessID] == nil {
			bizIDObjID2ObjMap[businessID] = make(map[string]metadata.Object)
		}
		for _, object := range objects {
			bizIDObjID2ObjMap[businessID][object.ObjectID] = object
		}
	}

	// get instance's model one by one
	instanceIDObjectMap := map[int64]metadata.Object{}
	for _, instance := range instances {
		objectMap, exist := bizIDObjID2ObjMap[instance.BizID]
		if exist == false {
			blog.Errorf("extractObjectIDFromInstances failed, instance's model not found, biz id %d not in bizIDObjID2ObjMap %+v, rid: %s", instance.BizID, bizIDObjID2ObjMap, rid)
			return nil, fmt.Errorf("get model by instance failed, unexpected err, business id:%d related models not found", instance.BizID)
		}

		object, exist := objectMap[instance.ObjectID]
		if exist == false {
			blog.Errorf("extractObjectIDFromInstances failed, instance's model not found, instances: %+v, objectMap: %+v, rid: %s", instance, objectMap, rid)
			return nil, fmt.Errorf("get model by instance failed, not found")
		}
		instanceIDObjectMap[instance.InstanceID] = object
	}

	blog.V(5).Infof("collectObjectsByInstances result: %+v, rid: %s", instanceIDObjectMap, rid)
	return instanceIDObjectMap, nil
}

func (am *AuthManager) MakeResourcesByInstances(ctx context.Context, header http.Header, action meta.Action, instances ...InstanceSimplify) ([]meta.ResourceAttribute, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	if len(instances) == 0 {
		return nil, nil
	}

	businessIDMap, err := am.extractBusinessIDFromInstances(ctx, instances...)
	if err != nil {
		return nil, fmt.Errorf("extract business id from instances failed, err: %+v", err)
	}

	instanceIDObjectMap, err := am.collectObjectsByInstances(ctx, header, instances...)
	if err != nil {
		blog.Errorf("MakeResourcesByInstances failed, collect objects by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("extract object by id failed, err: %+v", err)
	}

	// group instance by object
	objectIDInstancesMap := map[int64][]InstanceSimplify{}
	objectIDMap := map[int64]metadata.Object{}
	for _, instance := range instances {
		object, exist := instanceIDObjectMap[instance.InstanceID]
		if exist == false {
			blog.Errorf("MakeResourcesByInstances failed, unexpected err, instance related object not found, instance: %+v, instanceIDObjectMap: %+v, rid: %s", instance, instanceIDObjectMap, rid)
			return nil, errors.New("unexpected err, get instance related model failed, not found")
		}
		objectIDInstancesMap[object.ID] = append(objectIDInstancesMap[object.ID], instance)
		objectIDMap[object.ID] = object
	}

	mainlineTopo, err := am.clientSet.CoreService().Mainline().SearchMainlineModelTopo(ctx, header, false)
	if err != nil {
		blog.Errorf("list mainline models failed, err: %+v, rid: %s", err, rid)
		return nil, err
	}
	mainlineModels := mainlineTopo.LeftestObjectIDList()

	resultResources := make([]meta.ResourceAttribute, 0)
	for objID, instances := range objectIDInstancesMap {
		object := objectIDMap[objID]

		resourceType := meta.ModelInstance
		if util.InStrArr(mainlineModels, object.ObjectID) {
			resourceType = meta.MainlineInstance
		}

		parentResources, err := am.MakeResourcesByObjects(ctx, header, meta.EmptyAction, object)
		if err != nil {
			blog.Errorf("MakeResourcesByObjects failed, make parent auth resource by objects failed, object: %+v, err: %+v, rid: %s", object, err, rid)
			return nil, fmt.Errorf("make parent auth resource by objects failed, err: %+v", err)
		}
		if len(parentResources) != 1 {
			blog.Errorf("MakeResourcesByInstances failed, make parent auth resource by objects failed, get %d with object %s, rid: %s", len(parentResources), object.ObjectID, rid)
			return nil, fmt.Errorf("make parent auth resource by objects failed, get %d with object %d", len(parentResources), object.ID)
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
					Type:       resourceType,
					Name:       instance.Name,
					InstanceID: instance.InstanceID,
				},
				SupplierAccount: util.GetOwnerID(header),
				BusinessID:      businessIDMap[instance.InstanceID],
				Layers:          layers,
			}

			resources = append(resources, resource)
		}
		resultResources = append(resultResources, resources...)
	}
	return resultResources, nil
}

func (am *AuthManager) AuthorizeByInstanceID(ctx context.Context, header http.Header, action meta.Action, objID string, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	switch objID {
	case common.BKInnerObjIDPlat:
		return am.AuthorizeByPlatIDs(ctx, header, action, ids...)
	case common.BKInnerObjIDHost:
		return am.AuthorizeByHostsIDs(ctx, header, action, ids...)
	case common.BKInnerObjIDModule:
		return am.AuthorizeByModuleID(ctx, header, action, ids...)
	case common.BKInnerObjIDSet:
		return am.AuthorizeBySetID(ctx, header, action, ids...)
	case common.BKInnerObjIDApp:
		return am.AuthorizeByBusinessID(ctx, header, action, ids...)
	}

	instances, err := am.collectInstancesByRawIDs(ctx, header, objID, ids...)
	if err != nil {
		return fmt.Errorf("collect instance of model: %s by id %+v failed, err: %+v", objID, ids, err)
	}
	return am.AuthorizeByInstances(ctx, header, action, instances...)
}

func (am *AuthManager) AuthorizeByInstances(ctx context.Context, header http.Header, action meta.Action, instances ...InstanceSimplify) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if am.Enabled() == false {
		return nil
	}

	if am.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, instances: %+v, rid: %s", instances, rid)
		return nil
	}

	// make auth resources
	resources, err := am.MakeResourcesByInstances(ctx, header, action, instances...)
	if err != nil {
		blog.Errorf("AuthorizeByInstances failed, make resource by instances failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("make resource by instances failed, err: %+v", err)
	}

	return am.batchAuthorize(ctx, header, resources...)
}

func (am *AuthManager) UpdateRegisteredInstances(ctx context.Context, header http.Header, instances ...InstanceSimplify) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if am.Enabled() == false {
		return nil
	}

	// make auth resources
	resources, err := am.MakeResourcesByInstances(ctx, header, meta.EmptyAction, instances...)
	if err != nil {
		blog.Errorf("UpdateRegisteredInstances failed, make resource by instances failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("make resource by instances failed, err: %+v", err)
	}

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return fmt.Errorf("update resource %s/%d to iam failed, err: %v", resource.Type, resource.InstanceID, err)
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredInstanceByID(ctx context.Context, header http.Header, objectID string, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	switch objectID {
	case common.BKInnerObjIDPlat:
		return am.UpdateRegisteredPlatByID(ctx, header, ids...)
	case common.BKInnerObjIDHost:
		return am.UpdateRegisteredHostsByID(ctx, header, ids...)
	case common.BKInnerObjIDModule:
		return am.UpdateRegisteredModuleByID(ctx, header, ids...)
	case common.BKInnerObjIDSet:
		return am.UpdateRegisteredSetByID(ctx, header, ids...)
	case common.BKInnerObjIDApp:
		return am.UpdateRegisteredBusinessByID(ctx, header, ids...)
	}

	instances, err := am.collectInstancesByRawIDs(ctx, header, objectID, ids...)
	if err != nil {
		return fmt.Errorf("update registered instances failed, get instances by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredInstances(ctx, header, instances...)
}

func (am *AuthManager) UpdateRegisteredInstanceByRawID(ctx context.Context, header http.Header, objectID string, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	switch objectID {
	case common.BKInnerObjIDPlat:
		return am.UpdateRegisteredPlatByRawID(ctx, header, ids...)
	case common.BKInnerObjIDHost:
		return am.UpdateRegisteredHostsByID(ctx, header, ids...)
	case common.BKInnerObjIDModule:
		return am.UpdateRegisteredModuleByID(ctx, header, ids...)
	case common.BKInnerObjIDSet:
		return am.UpdateRegisteredSetByID(ctx, header, ids...)
	case common.BKInnerObjIDApp:
		return am.UpdateRegisteredBusinessByID(ctx, header, ids...)
	}

	instances, err := am.collectInstancesByRawIDs(ctx, header, objectID, ids...)
	if err != nil {
		return fmt.Errorf("update registered instances failed, get instances by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredInstances(ctx, header, instances...)
}

func (am *AuthManager) DeregisterInstanceByRawID(ctx context.Context, header http.Header, objectID string, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}
	switch objectID {
	case common.BKInnerObjIDPlat:
		return am.DeregisterPlatByID(ctx, header, ids...)
	case common.BKInnerObjIDHost:
		return am.DeregisterHostsByID(ctx, header, ids...)
	case common.BKInnerObjIDModule:
		return am.DeregisterModuleByID(ctx, header, ids...)
	case common.BKInnerObjIDSet:
		return am.DeregisterSetByID(ctx, header, ids...)
	case common.BKInnerObjIDApp:
		return am.DeregisterBusinessByRawID(ctx, header, ids...)
	}

	instances, err := am.collectInstancesByRawIDs(ctx, header, objectID, ids...)
	if err != nil {
		return fmt.Errorf("deregister instances failed, get instance by id failed, err: %+v", err)
	}
	return am.DeregisterInstances(ctx, header, instances...)
}

func (am *AuthManager) RegisterInstancesByID(ctx context.Context, header http.Header, objectID string, ids ...int64) error {
	if am.Enabled() == false {
		return nil
	}
	switch objectID {
	case common.BKInnerObjIDPlat:
		return am.RegisterPlatByID(ctx, header, ids...)
	case common.BKInnerObjIDHost:
		return am.RegisterHostsByID(ctx, header, ids...)
	case common.BKInnerObjIDModule:
		return am.RegisterModuleByID(ctx, header, ids...)
	case common.BKInnerObjIDSet:
		return am.RegisterSetByID(ctx, header, ids...)
	case common.BKInnerObjIDApp:
		return am.RegisterBusinessesByID(ctx, header, ids...)
	}

	instances, err := am.collectInstancesByRawIDs(ctx, header, objectID, ids...)
	if err != nil {
		return fmt.Errorf("register instances failed, get instance by id failed, err: %+v", err)
	}
	return am.registerInstances(ctx, header, instances...)
}

func (am *AuthManager) registerInstances(ctx context.Context, header http.Header, instances ...InstanceSimplify) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if am.Enabled() == false {
		return nil
	}

	// make auth resources
	resources, err := am.MakeResourcesByInstances(ctx, header, meta.EmptyAction, instances...)
	if err != nil {
		blog.Errorf("RegisterInstances failed, make resource by instances failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("make resource by instances failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterInstances(ctx context.Context, header http.Header, instances ...InstanceSimplify) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if am.Enabled() == false {
		return nil
	}

	// make auth resources
	resources, err := am.MakeResourcesByInstances(ctx, header, meta.EmptyAction, instances...)
	if err != nil {
		blog.Errorf("DeregisterInstances failed, make resource by instances failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("make resource by instances failed, err: %+v", err)
	}

	return am.Authorize.DeregisterResource(ctx, resources...)
}

// AuthorizeInstanceCreateByObjectID authorize create priority by object, plz be note this method only overlay model read/update/delete, without create
// func (am *AuthManager) AuthorizeInstanceCreateByObject(ctx context.Context, header http.Header, action meta.Action, objects ...metadata.Object) error {
// 	rid := util.ExtractRequestIDFromContext(ctx)
//
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	parentResources, err := am.MakeResourcesByObjects(ctx, header, action, objects...)
// 	if err != nil {
// 		blog.V(5).Infof("AuthorizeInstanceCreateByObject failed, make auth resource from objects failed, objects: %+v, err: %+v, rid: %s", objects, err, rid)
// 		return fmt.Errorf("make parent auth resource by models failed, err: %+v", err)
// 	}
//
// 	resources := make([]meta.ResourceAttribute, 0)
// 	for _, parentResource := range parentResources {
// 		layers := parentResource.Layers
// 		layers = append(layers, meta.Item{
// 			Type:       parentResource.Basic.Type,
// 			Action:     parentResource.Basic.Action,
// 			Name:       parentResource.Basic.Name,
// 			InstanceID: parentResource.Basic.InstanceID,
// 		})
// 		resource := meta.ResourceAttribute{
// 			Basic: meta.Basic{
// 				Type:   meta.ModelInstance,
// 				Action: meta.Create,
// 			},
// 			SupplierAccount: parentResource.SupplierAccount,
// 			BusinessID:      parentResource.BusinessID,
// 			Layers:          layers,
// 		}
// 		resources = append(resources, resource)
// 	}
//
// 	return am.batchAuthorize(ctx, header, resources...)
// }
