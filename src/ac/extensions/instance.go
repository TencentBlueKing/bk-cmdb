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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * instance represent common instances here
 */

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
		if !exist {
			blog.Errorf("extractObjectIDFromInstances failed, instance's model not found, biz id %d not in bizIDObjID2ObjMap %+v, rid: %s", instance.BizID, bizIDObjID2ObjMap, rid)
			return nil, fmt.Errorf("get model by instance failed, unexpected err, business id:%d related models not found", instance.BizID)
		}

		object, exist := objectMap[instance.ObjectID]
		if !exist {
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
		if !exist {
			blog.Errorf("MakeResourcesByInstances failed, unexpected err, instance related object not found, instance: %+v, instanceIDObjectMap: %+v, rid: %s", instance, instanceIDObjectMap, rid)
			return nil, errors.New("unexpected err, get instance related model failed, not found")
		}
		objectIDInstancesMap[object.ID] = append(objectIDInstancesMap[object.ID], instance)
		objectIDMap[object.ID] = object
	}

	mainlineAsst, err := am.clientSet.CoreService().Association().ReadModelAssociation(ctx, header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}})
	if err != nil {
		blog.Errorf("list mainline models failed, err: %+v, rid: %s", err, rid)
		return nil, err
	}

	resultResources := make([]meta.ResourceAttribute, 0)
	for objID, instances := range objectIDInstancesMap {
		object := objectIDMap[objID]
		resourceType := iam.GenCMDBDynamicResType(object.ObjectID, object.ID)
		for _, mainline := range mainlineAsst.Data.Info {
			if object.ObjectID == mainline.ObjectID {
				resourceType = meta.MainlineInstance
			}
		}

		resources := make([]meta.ResourceAttribute, 0)
		for _, instance := range instances {
			resource := meta.ResourceAttribute{
				Basic: meta.Basic{
					Action:     action,
					Type:       resourceType,
					Name:       instance.Name,
					InstanceID: instance.InstanceID,
				},
				SupplierAccount: util.GetOwnerID(header),
				BusinessID:      businessIDMap[instance.InstanceID],
			}

			resources = append(resources, resource)
		}
		resultResources = append(resultResources, resources...)
	}
	return resultResources, nil
}

func (am *AuthManager) AuthorizeByInstanceID(ctx context.Context, header http.Header, action meta.Action, objID string, ids ...int64) error {
	if !am.Enabled() {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	switch objID {
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

	if !am.Enabled() {
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

// RegisterModelResourceTypes register resource types for models
func (am *AuthManager) RegisterModelResourceTypes(ctx context.Context, iamClient iam.IAMClientInterface,
	objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}

	resourceTypes := iam.GenDynamicResourceTypes(objects)
	if err := iamClient.RegisterResourcesTypes(ctx, resourceTypes); err != nil {
		blog.ErrorJSON("register resourceTypes failed, error: %s, objects: %s, resourceTypes: %s, rid: %s",
			err.Error(), objects, resourceTypes, rid)
		return err
	}

	return nil
}

// UnregisterModelResourceTypes unregister resourceTypes for models
func (am *AuthManager) UnregisterModelResourceTypes(ctx context.Context, iamClient iam.IAMClientInterface,
	objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}

	typeIDs := []iam.TypeID{}
	resourceTypes := iam.GenDynamicResourceTypes(objects)
	for _, resourceType := range resourceTypes {
		typeIDs = append(typeIDs, resourceType.ID)
	}
	if err := iamClient.DeleteResourcesTypes(ctx, typeIDs); err != nil {
		blog.ErrorJSON("unregister resourceTypes failed, error: %s, objects: %s, resourceTypes: %s, rid: %s",
			err.Error(), objects, resourceTypes, rid)
		return err
	}

	return nil
}

// RegisterModelInstanceSelections register instanceSelections for models
func (am *AuthManager) RegisterModelInstanceSelections(ctx context.Context, iamClient iam.IAMClientInterface,
	objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}

	instanceSelections := iam.GenDynamicInstanceSelections(objects)
	if err := iamClient.RegisterInstanceSelections(ctx, instanceSelections); err != nil {
		blog.ErrorJSON("register instanceSelections failed, error: %s, objects: %s, instanceSelections: %s, rid: %s",
			err.Error(), objects,
			instanceSelections, rid)
		return err
	}

	return nil
}

// UnregisterModelInstanceSelections unregister instanceSelections for models
func (am *AuthManager) UnregisterModelInstanceSelections(ctx context.Context, iamClient iam.IAMClientInterface,
	objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}

	instanceSelectionIDs := []iam.InstanceSelectionID{}
	instanceSelections := iam.GenDynamicInstanceSelections(objects)
	for _, instanceSelection := range instanceSelections {
		instanceSelectionIDs = append(instanceSelectionIDs, instanceSelection.ID)
	}
	if err := iamClient.DeleteInstanceSelections(ctx, instanceSelectionIDs); err != nil {
		blog.ErrorJSON("unregister instanceSelections failed, error: %s, objects: %s, instanceSelections: %s, rid: %s",
			err.Error(), objects,
			instanceSelections, rid)
		return err
	}

	return nil
}

// RegisterModelActions register actions for models
func (am *AuthManager) RegisterModelActions(ctx context.Context, iamClient iam.IAMClientInterface,
	objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}

	actions := iam.GenDynamicActions(objects)
	if err := iamClient.RegisterActions(ctx, actions); err != nil {
		blog.ErrorJSON("register actions failed, error: %s, objects: %s, actions: %s, rid: %s", err.Error(), objects,
			actions, rid)
		return err
	}

	return nil
}

// UnregisterModelActions unregister actions for models
func (am *AuthManager) UnregisterModelActions(ctx context.Context, iamClient iam.IAMClientInterface,
	objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}

	actionIDs := []iam.ActionID{}
	for _, obj := range objects {
		actionIDs = append(actionIDs, iam.GenDynamicActionIDs(obj)...)
	}
	if err := iamClient.DeleteActions(ctx, actionIDs); err != nil {
		blog.ErrorJSON("unregister actions failed, error: %s, objects: %s, actionIDs: %s, rid: %s", err.Error(),
			objects, actionIDs, rid)
		return err
	}

	return nil
}

// UpdateModelActionGroups update actionGroups for models
func (am *AuthManager) UpdateModelActionGroups(ctx context.Context, header http.Header,
	iamClient iam.IAMClientInterface) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}

	// IAM目前不支持增量更新，仅支持完全覆盖更新，即更新所有actionGroup
	actionGroups := iam.GenerateStaticActionGroups()
	objects, err := am.GetCustomObjects(ctx, header)
	if err != nil {
		blog.Errorf("get custom objects failed, err: %s, rid:%s", err.Error(), rid)
		return err
	}
	actionGroups = append(actionGroups, iam.GenModelInstanceManageActionGroups(objects)...)

	if err := iamClient.UpdateActionGroups(ctx, actionGroups); err != nil {
		blog.ErrorJSON("update actionGroups failed, error: %s, actionGroups: %s, rid: %s", err.Error(),
			actionGroups, rid)
		return err
	}

	return nil
}

// GetCustomObjects get objects which are custom
func (am *AuthManager) GetCustomObjects(ctx context.Context, header http.Header) ([]metadata.Object, error) {
	resp, err := am.clientSet.CoreService().Model().ReadModel(ctx, header, &metadata.QueryCondition{
		Fields: []string{common.BKObjIDField, common.BKObjNameField, common.BKFieldID},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKIsPre: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get custom objects failed, err: %+v", err)
	}
	if len(resp.Data.Info) == 0 {
		blog.Info("get custom objects failed, no custom objects were found")
	}

	objects := make([]metadata.Object, 0)
	for _, item := range resp.Data.Info {
		objects = append(objects, item.Spec)
	}

	return objects, nil
}
