/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package iam

import (
	"reflect"

	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/apimachinery"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	commonlgc "configcenter/src/common/logics"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/apigw/iam"
)

type viewer struct {
	client apimachinery.ClientSetInterface
	iam    *IAM
}

// NewViewer TODO
func NewViewer(client apimachinery.ClientSetInterface, iam *IAM) *viewer {
	return &viewer{
		client: client,
		iam:    iam,
	}
}

// CreateView create iam view for objects
func (v *viewer) CreateView(kit *rest.Kit, objects []metadata.Object, redisCli redis.Client,
	rid string) error {

	if !auth.EnableAuthorize() {
		return nil
	}

	locker, err := tryLockRegister(redisCli, rid)
	if err != nil {
		return err
	}
	defer locker.Unlock()

	systemInfo, err := v.iam.Client.GetSystemInfo(kit.Ctx, kit.Header, []iamtypes.SystemQueryField{})
	if err != nil {
		blog.Errorf("get iam system info failed, err: %v, rid: %s", err, rid)
		return err
	}

	// register order: 1.ResourceType 2.InstanceSelection 3.Action 4.ActionGroup
	if err = v.registerModelResourceTypes(kit, systemInfo.ResourceTypes, objects); err != nil {
		return err
	}

	if err = v.registerModelInstanceSelections(kit, systemInfo.InstanceSelections, objects); err != nil {
		return err
	}

	if err = v.registerModelActions(kit, systemInfo.Actions, objects); err != nil {
		return err
	}

	if err = v.updateModelActionGroups(kit, systemInfo.ActionGroups); err != nil {
		return err
	}

	return nil
}

// DeleteView delete iam view for objects
func (v *viewer) DeleteView(kit *rest.Kit, objects []metadata.Object, redisCli redis.Client,
	rid string) error {

	if !auth.EnableAuthorize() {
		return nil
	}

	locker, err := tryLockRegister(redisCli, rid)
	if err != nil {
		return err
	}
	defer locker.Unlock()

	systemInfo, err := v.iam.Client.GetSystemInfo(kit.Ctx, kit.Header, []iamtypes.SystemQueryField{})
	if err != nil {
		blog.Errorf("get iam system info failed, err: %v, rid: %s", err, rid)
		return err
	}

	// unregister order: 1.Action 2.InstanceSelection 3.ResourceType 4.ActionGroup
	if err = v.unregisterModelActions(kit, systemInfo.Actions, objects); err != nil {
		return err
	}

	if err = v.unregisterModelInstanceSelections(kit, systemInfo.InstanceSelections, objects); err != nil {
		return err
	}

	if err = v.unregisterModelResourceTypes(kit, systemInfo.ResourceTypes, objects); err != nil {
		return err
	}

	if err = v.updateModelActionGroups(kit, systemInfo.ActionGroups); err != nil {
		return err
	}

	return nil
}

// UpdateView update iam view for objects
func (v *viewer) UpdateView(kit *rest.Kit, objects []metadata.Object, redisCli redis.Client, rid string) error {

	if !auth.EnableAuthorize() {
		return nil
	}

	locker, err := tryLockRegister(redisCli, rid)
	if err != nil {
		return err
	}
	defer locker.Unlock()

	systemInfo, err := v.iam.Client.GetSystemInfo(kit.Ctx, kit.Header, []iamtypes.SystemQueryField{})
	if err != nil {
		blog.Errorf("get iam system info failed, err: %v, rid: %s", err, rid)
		return err
	}

	// update order: 1.ResourceType 2.InstanceSelection 3.Action 4.ActionGroup
	if err = v.updateModelResourceTypes(kit, objects); err != nil {
		return err
	}

	if err = v.updateModelInstanceSelections(kit, objects); err != nil {
		return err
	}

	if err = v.updateModelActions(kit, objects); err != nil {
		return err
	}

	if err = v.updateModelActionGroups(kit, systemInfo.ActionGroups); err != nil {
		return err
	}

	return nil
}

// registerModelResourceTypes register resource types for models
func (v *viewer) registerModelResourceTypes(kit *rest.Kit, existTypes []iam.ResourceType,
	objects []metadata.Object) error {

	rid := kit.Rid

	existMap := make(map[iamtypes.TypeID]struct{})
	for _, resourceType := range existTypes {
		existMap[resourceType.ID] = struct{}{}
	}

	resourceTypes := genDynamicResourceTypes(objects)
	newResTypes := make([]iam.ResourceType, 0)
	for _, resourceType := range resourceTypes {
		_, exists := existMap[resourceType.ID]
		if !exists {
			newResTypes = append(newResTypes, resourceType)
			continue
		}

		if err := v.iam.Client.UpdateResourcesType(kit.Ctx, kit.Header, resourceType); err != nil {
			blog.Errorf("update resource type failed, err: %v, resourceType: %+v， rid: %s", err, resourceType, rid)
			return err
		}

	}

	if err := v.iam.Client.RegisterResourcesTypes(kit.Ctx, kit.Header, newResTypes); err != nil {
		blog.Errorf("register resource types failed, err: %v, objects: %+v, resourceTypes: %+v， rid: %s", err, objects,
			newResTypes, rid)
		return err
	}

	return nil
}

// unregisterModelResourceTypes unregister resourceTypes for models
func (v *viewer) unregisterModelResourceTypes(kit *rest.Kit, existTypes []iam.ResourceType,
	objects []metadata.Object) error {

	rid := kit.Rid

	existMap := make(map[iamtypes.TypeID]struct{})
	for _, resourceType := range existTypes {
		existMap[resourceType.ID] = struct{}{}
	}

	typeIDs := make([]iamtypes.TypeID, 0)
	resourceTypes := genDynamicResourceTypes(objects)
	for _, resourceType := range resourceTypes {
		_, exists := existMap[resourceType.ID]
		if exists {
			typeIDs = append(typeIDs, resourceType.ID)
		}
	}

	if err := v.iam.Client.DeleteResourcesTypes(kit.Ctx, kit.Header, typeIDs); err != nil {
		blog.Errorf("unregister resourceTypes failed, err: %v, objects: %+v, resourceTypes: %+v, rid: %s",
			err, objects, resourceTypes, rid)
		return err
	}

	return nil
}

// updateModelResourceTypes update resource types for models
func (v *viewer) updateModelResourceTypes(kit *rest.Kit, objects []metadata.Object) error {
	rid := kit.Rid
	resourceTypes := genDynamicResourceTypes(objects)
	for _, resourceType := range resourceTypes {
		if err := v.iam.Client.UpdateResourcesType(kit.Ctx, kit.Header, resourceType); err != nil {
			blog.Errorf("update resourceType failed, err: %v, objects: %+v, resourceType: %+v，rid: %s", err, objects,
				resourceType, rid)
			return err
		}
	}

	return nil
}

// registerModelInstanceSelections register instanceSelections for models
func (v *viewer) registerModelInstanceSelections(kit *rest.Kit, existsSelections []iam.InstanceSelection,
	objects []metadata.Object) error {

	rid := kit.Rid

	existMap := make(map[iamtypes.InstanceSelectionID]struct{})
	for _, selection := range existsSelections {
		existMap[selection.ID] = struct{}{}
	}

	instanceSelections := genDynamicInstanceSelections(objects)
	newSelections := make([]iam.InstanceSelection, 0)
	for _, selection := range instanceSelections {
		_, exists := existMap[selection.ID]
		if !exists {
			newSelections = append(newSelections, selection)
			continue
		}

		if err := v.iam.Client.UpdateInstanceSelection(kit.Ctx, kit.Header, selection); err != nil {
			blog.Errorf("update instanceSelection %+v failed, err: %v, rid: %s", selection, err, rid)
			return err
		}
	}

	if err := v.iam.Client.RegisterInstanceSelections(kit.Ctx, kit.Header, newSelections); err != nil {
		blog.Errorf("register instanceSelections failed, err: %v, objects: %+v, instanceSelections: %+v, rid: %s",
			err, objects, instanceSelections, rid)
		return err
	}

	return nil
}

// unregisterModelInstanceSelections unregister instanceSelections for models
func (v *viewer) unregisterModelInstanceSelections(kit *rest.Kit, existsSelections []iam.InstanceSelection,
	objects []metadata.Object) error {

	rid := kit.Rid

	existMap := make(map[iamtypes.InstanceSelectionID]struct{})
	for _, selection := range existsSelections {
		existMap[selection.ID] = struct{}{}
	}

	instanceSelectionIDs := make([]iamtypes.InstanceSelectionID, 0)
	instanceSelections := genDynamicInstanceSelections(objects)
	for _, instanceSelection := range instanceSelections {
		_, exists := existMap[instanceSelection.ID]
		if exists {
			instanceSelectionIDs = append(instanceSelectionIDs, instanceSelection.ID)
			continue
		}
	}

	if err := v.iam.Client.DeleteInstanceSelections(kit.Ctx, kit.Header, instanceSelectionIDs); err != nil {
		blog.Errorf("unregister instanceSelections failed, err: %v, objects: %+v, instanceSelections: %+v, rid: %s",
			err, objects, instanceSelections, rid)
		return err
	}

	return nil
}

// updateModelInstanceSelections update instanceSelections for models
func (v *viewer) updateModelInstanceSelections(kit *rest.Kit, objects []metadata.Object) error {
	rid := kit.Rid
	instanceSelections := genDynamicInstanceSelections(objects)
	for _, instanceSelection := range instanceSelections {
		if err := v.iam.Client.UpdateInstanceSelection(kit.Ctx, kit.Header, instanceSelection); err != nil {
			blog.Errorf("update instanceSelections failed, err: %v, objects: %+v, instanceSelection: %+v, rid: %s",
				err, objects, instanceSelection, rid)
			return err
		}
	}

	return nil
}

// registerModelActions register actions for models
func (v *viewer) registerModelActions(kit *rest.Kit, existAction []iam.ResourceAction,
	objects []metadata.Object) error {

	rid := kit.Rid
	existMap := make(map[iamtypes.ActionID]struct{})
	for _, action := range existAction {
		existMap[action.ID] = struct{}{}
	}

	actions := genDynamicActions(objects)
	newActions := make([]iam.ResourceAction, 0)
	for _, action := range actions {
		_, exists := existMap[action.ID]
		if !exists {
			newActions = append(newActions, action)
			continue
		}

		if err := v.iam.Client.UpdateAction(kit.Ctx, kit.Header, action); err != nil {
			blog.Errorf("update action %+v failed, err: %v, rid: %s", action, err, rid)
			return err
		}
	}

	if err := v.iam.Client.RegisterActions(kit.Ctx, kit.Header, newActions); err != nil {
		blog.Errorf("register actions %+v failed, err: %v, objects: %+v, rid: %s", err, objects, newActions, rid)
		return err
	}

	return nil
}

// unregisterModelActions unregister actions for models
func (v *viewer) unregisterModelActions(kit *rest.Kit, existAction []iam.ResourceAction,
	objects []metadata.Object) error {

	rid := kit.Rid

	existMap := make(map[iamtypes.ActionID]struct{})
	for _, action := range existAction {
		existMap[action.ID] = struct{}{}
	}

	actionIDs := make([]iamtypes.ActionID, 0)
	for _, obj := range objects {
		ids := genDynamicActionIDs(obj)
		for _, id := range ids {
			_, exists := existMap[id]
			if exists {
				actionIDs = append(actionIDs, id)
				continue
			}
		}
	}

	// before deleting action, the dependent action policies must be deleted
	for _, actionID := range actionIDs {
		if err := v.iam.Client.DeleteActionPolicies(kit.Ctx, kit.Header, actionID); err != nil {
			blog.Errorf("delete action %s policies failed, err: %v, rid: %s", actionID, err, rid)
			return err
		}
	}

	if err := v.iam.Client.DeleteActions(kit.Ctx, kit.Header, actionIDs); err != nil {
		blog.Errorf("unregister actions failed, err: %v, objects: %+v, actionIDs: %s, rid: %s", err, objects, actionIDs,
			rid)
		return err
	}

	return nil
}

// updateModelActions update actions for models
func (v *viewer) updateModelActions(kit *rest.Kit, objects []metadata.Object) error {
	rid := kit.Rid
	actions := genDynamicActions(objects)
	for _, action := range actions {
		if err := v.iam.Client.UpdateAction(kit.Ctx, kit.Header, action); err != nil {
			blog.Errorf("update action failed, err: %v, objects: %+v, action: %+v, rid: %s", err, objects, action, rid)
			return err
		}
	}

	return nil
}

// updateModelActionGroups update actionGroups for models
// for now, the update api can only support full update, not incremental update
func (v *viewer) updateModelActionGroups(kit *rest.Kit, existGroups []iam.ActionGroup) error {
	rid := kit.Rid
	objects, err := commonlgc.GetCustomObjects(kit, v.client)
	if err != nil {
		blog.Errorf("get custom objects failed, err: %v, rid: %s", err, rid)
		return err
	}
	actionGroups := GenerateActionGroups(objects)

	if len(existGroups) == 0 {
		if err = v.iam.Client.RegisterActionGroups(kit.Ctx, kit.Header, actionGroups); err != nil {
			blog.Errorf("register action groups(%s) failed, err: %v, rid: %s", actionGroups, err, rid)
			return err
		}
		return nil
	}

	if reflect.DeepEqual(existGroups, actionGroups) {
		return nil
	}

	if err = v.iam.Client.UpdateActionGroups(kit.Ctx, kit.Header, actionGroups); err != nil {
		blog.Errorf("update actionGroups failed, error: %v, actionGroups: %+v, rid: %s", err, actionGroups, rid)
		return err
	}

	return nil
}
