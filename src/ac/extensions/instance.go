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

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * instance represent common instances here
 */

func (a *AuthManager) collectInstancesByRawIDs(ctx context.Context, header http.Header, modelID string, ids ...int64) (
	[]InstanceSimplify, error) {

	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.IntArrayUnique(ids)
	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.GetInstIDField(modelID)).In(ids).ToMapStr(),
	}
	result, err := a.clientSet.CoreService().Instance().ReadInstance(ctx, header, modelID, &cond)
	if err != nil {
		blog.V(3).Infof("get instance by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get instance by id failed, err: %+v", err)
	}
	instances := make([]InstanceSimplify, 0)
	for _, inst := range result.Info {
		instance := InstanceSimplify{}
		_, err = instance.Parse(inst)
		if err != nil {
			return nil, fmt.Errorf("get instance by object failed, err: %+v", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (a *AuthManager) extractBusinessIDFromInstances(ctx context.Context,
	instances ...InstanceSimplify) (map[int64]int64, error) {
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
func (a *AuthManager) collectObjectsByInstances(ctx context.Context, header http.Header,
	instances ...InstanceSimplify) (map[int64]metadata.Object, error) {
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
		objects, err := a.collectObjectsByObjectIDs(ctx, header, businessID, objectIDs...)
		if err != nil {
			blog.Errorf("extractObjectIDFromInstances failed, get models by businessID and object id failed, bizID: %+v, objectIDs: %+v, err: %+v, rid: %s",
				businessID, objectIDs, err, rid)
			return nil, fmt.Errorf("get models by objectIDs and business id failed")
		}
		if len(objectIDs) != len(objects) {
			blog.Errorf("extractObjectIDFromInstances failed, get models by object id failed, input len %d and output len %d not equal, input: %+v, output: %+v, businessID: %d",
				len(objectIDs), len(objects), objectIDs, objects, businessID)
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
		if !exist {
			blog.Errorf("extractObjectIDFromInstances failed, instance's model not found, biz id %d not in bizIDObjID2ObjMap %+v, rid: %s",
				instance.BizID, bizIDObjID2ObjMap, rid)
			return nil, fmt.Errorf("get model by instance failed, unexpected err, business id:%d related models not found",
				instance.BizID)
		}

		object, exist := objectMap[instance.ObjectID]
		if !exist {
			blog.Errorf("extractObjectIDFromInstances failed, instance's model not found, instances: %+v, objectMap: %+v, rid: %s",
				instance, objectMap, rid)
			return nil, fmt.Errorf("get model by instance failed, not found")
		}
		instanceIDObjectMap[instance.InstanceID] = object
	}

	blog.V(5).Infof("collectObjectsByInstances result: %+v, rid: %s", instanceIDObjectMap, rid)
	return instanceIDObjectMap, nil
}

// MakeResourcesByInstances make the resources by simplify instances and action
func (a *AuthManager) MakeResourcesByInstances(ctx context.Context, header http.Header, action meta.Action,
	instances ...InstanceSimplify) ([]meta.ResourceAttribute, error) {

	rid := util.ExtractRequestIDFromContext(ctx)

	if len(instances) == 0 {
		return nil, nil
	}

	businessIDMap, err := a.extractBusinessIDFromInstances(ctx, instances...)
	if err != nil {
		return nil, fmt.Errorf("extract business id from instances failed, err: %+v", err)
	}

	instanceIDObjectMap, err := a.collectObjectsByInstances(ctx, header, instances...)
	if err != nil {
		blog.Errorf("MakeResourcesByInstances failed, collect objects by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("extract object by id failed, err: %+v", err)
	}

	// group instance by object
	objectIDInstancesMap := map[int64][]InstanceSimplify{}
	objectIDMap := map[int64]metadata.Object{}
	for _, instance := range instances {
		object, exist := instanceIDObjectMap[instance.InstanceID]
		if !exist {
			blog.Errorf("MakeResourcesByInstances failed, unexpected err, instance related object not found, "+
				"instance: %+v, instanceIDObjectMap: %+v, rid: %s", instance, instanceIDObjectMap, rid)
			return nil, errors.New("unexpected err, get instance related model failed, not found")
		}
		objectIDInstancesMap[object.ID] = append(objectIDInstancesMap[object.ID], instance)
		objectIDMap[object.ID] = object
	}

	mainlineAsst, err := a.clientSet.CoreService().Association().ReadModelAssociation(ctx, header,
		&metadata.QueryCondition{
			Condition: mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline},
		})
	if err != nil {
		blog.Errorf("list mainline models failed, err: %+v, rid: %s", err, rid)
		return nil, err
	}

	resultResources := make([]meta.ResourceAttribute, 0)
	for objID, instances := range objectIDInstancesMap {
		object := objectIDMap[objID]

		resourceType := iam.GenCMDBDynamicResType(object.ID)
		for _, mainline := range mainlineAsst.Info {
			if object.ObjectID == mainline.ObjectID {
				resourceType = meta.MainlineInstance
			}
		}

		if resourceType == meta.MainlineInstance {
			parentResources, err := a.MakeResourcesByObjects(ctx, header, meta.EmptyAction, object)
			if err != nil {
				blog.Errorf("MakeResourcesByObjects failed, make parent auth resource by objects failed, "+
					"object: %+v, err: %+v, rid: %s", object, err, rid)
				return nil, fmt.Errorf("make parent auth resource by objects failed, err: %+v", err)
			}
			if len(parentResources) != 1 {
				blog.Errorf("MakeResourcesByInstances failed, make parent auth resource by objects failed, "+
					"get %d with object %s, rid: %s", len(parentResources), object.ObjectID, rid)
				return nil, fmt.Errorf("make parent auth resource by objects failed, get %d with object %d",
					len(parentResources), object.ID)
			}

			parentResource := parentResources[0]
			layers := parentResource.Layers
			layers = append(layers, meta.Item{
				Type:       parentResource.Type,
				Action:     parentResource.Action,
				Name:       parentResource.Name,
				InstanceID: parentResource.InstanceID,
			})
			resources := make([]meta.ResourceAttribute, 0)
			for _, instance := range instances {

				resource := meta.ResourceAttribute{
					Basic: meta.Basic{
						Action:     action,
						Type:       resourceType,
						Name:       instance.Name,
						InstanceID: instance.InstanceID,
					},
					TenantID:   httpheader.GetTenantID(header),
					BusinessID: businessIDMap[instance.InstanceID],
					Layers:     layers,
				}

				resources = append(resources, resource)
			}
			resultResources = append(resultResources, resources...)
		} else {
			resources := make([]meta.ResourceAttribute, 0)
			for _, instance := range instances {
				resource := meta.ResourceAttribute{
					Basic: meta.Basic{
						Action:     action,
						Type:       resourceType,
						Name:       instance.Name,
						InstanceID: instance.InstanceID,
					},
					TenantID:   httpheader.GetTenantID(header),
					BusinessID: businessIDMap[instance.InstanceID],
				}

				resources = append(resources, resource)
			}
			resultResources = append(resultResources, resources...)
		}
	}
	return resultResources, nil
}

// AuthorizeByInstIDs authorize by instance ids, returns auth response if not authorized
func (a *AuthManager) AuthorizeByInstIDs(kit *rest.Kit, action meta.Action, objID string, ids ...int64) (
	*metadata.BaseResp, bool, error) {

	if !a.Enabled() {
		return nil, true, nil
	}

	if len(ids) == 0 {
		return nil, true, nil
	}

	if a.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, obj: %s, ids: %+v, kit.Rid: %s", objID, ids, kit.Rid)
		return nil, true, nil
	}

	var resources []meta.ResourceAttribute
	var resErr error
	switch objID {
	case common.BKInnerObjIDHost:
		hosts, err := a.collectHostByHostIDs(kit.Ctx, kit.Header, ids...)
		if err != nil {
			return nil, true, fmt.Errorf("get hosts by ids: %+v for authorization failed, err: %v", ids, err)
		}
		resources, resErr = a.MakeResourcesByHosts(kit.Ctx, kit.Header, action, hosts...)
	case common.BKInnerObjIDModule:
		if !a.RegisterModuleEnabled {
			return nil, true, nil
		}
		modules, err := a.collectModuleByModuleIDs(kit.Ctx, kit.Header, ids...)
		if err != nil {
			return nil, true, fmt.Errorf("get modules by ids: %+v for authorization failed, err: %v", ids, err)
		}
		bizID, err := a.extractBusinessIDFromModules(modules...)
		if err != nil {
			return nil, true, fmt.Errorf("extract business id from modules failed, err: %v", err)
		}
		resources = a.MakeResourcesByModule(kit.Header, action, bizID, modules...)
	case common.BKInnerObjIDSet:
		if !a.RegisterSetEnabled {
			return nil, true, nil
		}
		sets, err := a.collectSetBySetIDs(kit.Ctx, kit.Header, ids...)
		if err != nil {
			return nil, true, fmt.Errorf("collect set by id failed, err: %v", err)
		}
		bizID, err := a.extractBusinessIDFromSets(sets...)
		if err != nil {
			return nil, true, fmt.Errorf("authorize sets failed, extract business id from sets failed, err: %v", err)
		}
		resources = a.MakeResourcesBySet(kit.Header, action, bizID, sets...)
	case common.BKInnerObjIDApp:
		resources, resErr = a.genBizAuthRes(kit, action, ids)
	case common.BKInnerObjIDBizSet:
		bizSets, err := a.collectBizSetByIDs(kit.Ctx, kit.Header, kit.Rid, ids...)
		if err != nil {
			return nil, true, fmt.Errorf("get biz set by ids: %+v failed, err: %v", ids, err)
		}
		resources = a.makeResourcesByBizSet(kit.Header, action, bizSets...)
	default:
		instances, err := a.collectInstancesByRawIDs(kit.Ctx, kit.Header, objID, ids...)
		if err != nil {
			return nil, true, fmt.Errorf("get %s instance by ids: %+v for auth failed, err: %v", objID, ids, err)
		}
		resources, resErr = a.MakeResourcesByInstances(kit.Ctx, kit.Header, action, instances...)
	}

	if resErr != nil {
		blog.Errorf("make resource by %s instances(%+v) failed, err: %+v, kit.Rid: %s", objID, ids, resErr, kit.Rid)
		return nil, true, fmt.Errorf("make resource by instances failed, err: %v", resErr)
	}

	authResp, authorized := a.Authorize(kit, resources...)
	return authResp, authorized, nil
}

func (a *AuthManager) genBizAuthRes(kit *rest.Kit, action meta.Action, ids []int64) ([]meta.ResourceAttribute, error) {
	businesses, err := a.collectBusinessByIDs(kit.Ctx, kit.Header, ids...)
	if err != nil {
		return nil, fmt.Errorf("authorize businesses failed, get business by id failed, err: %v", err)
	}
	resourcePoolBusinessID, err := a.getResourcePoolBusinessID(kit.Ctx, kit.Header)
	if err != nil {
		return nil, err
	}

	bizArr := make([]BusinessSimplify, 0)
	if action == meta.ViewBusinessResource {
		for _, biz := range businesses {
			if biz.BKAppIDField == resourcePoolBusinessID {
				continue
			}
			bizArr = append(bizArr, biz)
		}
	} else {
		bizArr = businesses
	}

	// make auth resources
	resources := a.MakeResourcesByBusiness(kit.Header, action, bizArr...)
	return resources, nil
}

// AuthorizeByInstanceID TODO
func (a *AuthManager) AuthorizeByInstanceID(ctx context.Context, header http.Header, action meta.Action, objID string,
	ids ...int64) error {
	if !a.Enabled() {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	switch objID {
	case common.BKInnerObjIDHost:
		return a.AuthorizeByHostsIDs(ctx, header, action, ids...)
	case common.BKInnerObjIDModule:
		return a.AuthorizeByModuleID(ctx, header, action, ids...)
	case common.BKInnerObjIDSet:
		return a.AuthorizeBySetID(ctx, header, action, ids...)
	case common.BKInnerObjIDApp:
		return a.AuthorizeByBusinessID(ctx, header, action, ids...)
	case common.BKInnerObjIDBizSet:
		return a.AuthorizeByBizSetID(ctx, header, action, ids...)
	}

	instances, err := a.collectInstancesByRawIDs(ctx, header, objID, ids...)
	if err != nil {
		return fmt.Errorf("collect instance of model: %s by id %+v failed, err: %+v", objID, ids, err)
	}
	return a.AuthorizeByInstances(ctx, header, action, instances...)
}

// AuthorizeByInstances TODO
func (a *AuthManager) AuthorizeByInstances(ctx context.Context, header http.Header, action meta.Action,
	instances ...InstanceSimplify) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !a.Enabled() {
		return nil
	}

	if a.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, instances: %+v, rid: %s", instances, rid)
		return nil
	}

	// make auth resources
	resources, err := a.MakeResourcesByInstances(ctx, header, action, instances...)
	if err != nil {
		blog.Errorf("AuthorizeByInstances failed, make resource by instances failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("make resource by instances failed, err: %+v", err)
	}

	return a.batchAuthorize(ctx, header, resources...)
}
