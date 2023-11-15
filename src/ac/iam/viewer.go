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
	"context"
	"net/http"
	"reflect"

	"configcenter/src/apimachinery"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	commonlgc "configcenter/src/common/logics"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/redis"
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
func (v *viewer) CreateView(ctx context.Context, header http.Header, objects []metadata.Object, redisCli redis.Client,
	rid string) error {

	if !auth.EnableAuthorize() {
		return nil
	}

	locker, err := tryLockRegister(redisCli, rid)
	if err != nil {
		return err
	}
	defer locker.Unlock()

	systemResp, err := v.iam.Client.GetSystemInfo(ctx, []SystemQueryField{})
	if err != nil {
		blog.Errorf("get iam system info failed, err: %v, rid: %s", err, rid)
		return err
	}

	// register order: 1.ResourceType 2.InstanceSelection 3.Action 4.ActionGroup
	if err = v.registerModelResourceTypes(ctx, systemResp.Data.ResourceTypes, objects); err != nil {
		return err
	}

	if err = v.registerModelInstanceSelections(ctx, systemResp.Data.InstanceSelections, objects); err != nil {
		return err
	}

	if err = v.registerModelActions(ctx, systemResp.Data.Actions, objects); err != nil {
		return err
	}

	if err = v.updateModelActionGroups(ctx, header, systemResp.Data.ActionGroups); err != nil {
		return err
	}

	return nil
}

// DeleteView delete iam view for objects
func (v *viewer) DeleteView(ctx context.Context, header http.Header, objects []metadata.Object, redisCli redis.Client,
	rid string) error {

	if !auth.EnableAuthorize() {
		return nil
	}

	locker, err := tryLockRegister(redisCli, rid)
	if err != nil {
		return err
	}
	defer locker.Unlock()

	systemResp, err := v.iam.Client.GetSystemInfo(ctx, []SystemQueryField{})
	if err != nil {
		blog.Errorf("get iam system info failed, err: %v, rid: %s", err, rid)
		return err
	}

	// unregister order: 1.Action 2.InstanceSelection 3.ResourceType 4.ActionGroup
	if err = v.unregisterModelActions(ctx, systemResp.Data.Actions, objects); err != nil {
		return err
	}

	if err = v.unregisterModelInstanceSelections(ctx, systemResp.Data.InstanceSelections, objects); err != nil {
		return err
	}

	if err = v.unregisterModelResourceTypes(ctx, systemResp.Data.ResourceTypes, objects); err != nil {
		return err
	}

	if err = v.updateModelActionGroups(ctx, header, systemResp.Data.ActionGroups); err != nil {
		return err
	}

	return nil
}

// UpdateView update iam view for objects
func (v *viewer) UpdateView(ctx context.Context, header http.Header, objects []metadata.Object, redisCli redis.Client,
	rid string) error {

	if !auth.EnableAuthorize() {
		return nil
	}

	locker, err := tryLockRegister(redisCli, rid)
	if err != nil {
		return err
	}
	defer locker.Unlock()

	systemResp, err := v.iam.Client.GetSystemInfo(ctx, []SystemQueryField{})
	if err != nil {
		blog.Errorf("get iam system info failed, err: %v, rid: %s", err, rid)
		return err
	}

	// update order: 1.ResourceType 2.InstanceSelection 3.Action 4.ActionGroup
	if err = v.updateModelResourceTypes(ctx, objects); err != nil {
		return err
	}

	if err = v.updateModelInstanceSelections(ctx, objects); err != nil {
		return err
	}

	if err = v.updateModelActions(ctx, objects); err != nil {
		return err
	}

	if err = v.updateModelActionGroups(ctx, header, systemResp.Data.ActionGroups); err != nil {
		return err
	}

	return nil
}

// registerModelResourceTypes register resource types for models
func (v *viewer) registerModelResourceTypes(ctx context.Context, existTypes []ResourceType,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	existMap := make(map[TypeID]struct{})
	for _, resourceType := range existTypes {
		existMap[resourceType.ID] = struct{}{}
	}

	resourceTypes := genDynamicResourceTypes(objects)
	newResTypes := make([]ResourceType, 0)
	for _, resourceType := range resourceTypes {
		_, exists := existMap[resourceType.ID]
		if !exists {
			newResTypes = append(newResTypes, resourceType)
			continue
		}

		if err := v.iam.Client.UpdateResourcesType(ctx, resourceType); err != nil {
			blog.Errorf("update resource type failed, err: %v, resourceType: %+v， rid: %s", err, resourceType, rid)
			return err
		}

	}

	if err := v.iam.Client.RegisterResourcesTypes(ctx, newResTypes); err != nil {
		blog.Errorf("register resource types failed, err: %v, objects: %+v, resourceTypes: %+v， rid: %s", err, objects,
			newResTypes, rid)
		return err
	}

	return nil
}

// unregisterModelResourceTypes unregister resourceTypes for models
func (v *viewer) unregisterModelResourceTypes(ctx context.Context, existTypes []ResourceType,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	existMap := make(map[TypeID]struct{})
	for _, resourceType := range existTypes {
		existMap[resourceType.ID] = struct{}{}
	}

	typeIDs := make([]TypeID, 0)
	resourceTypes := genDynamicResourceTypes(objects)
	for _, resourceType := range resourceTypes {
		_, exists := existMap[resourceType.ID]
		if exists {
			typeIDs = append(typeIDs, resourceType.ID)
		}
	}

	if err := v.iam.Client.DeleteResourcesTypes(ctx, typeIDs); err != nil {
		blog.Errorf("unregister resourceTypes failed, err: %v, objects: %+v, resourceTypes: %+v, rid: %s",
			err, objects, resourceTypes, rid)
		return err
	}

	return nil
}

// updateModelResourceTypes update resource types for models
func (v *viewer) updateModelResourceTypes(ctx context.Context, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	resourceTypes := genDynamicResourceTypes(objects)
	for _, resourceType := range resourceTypes {
		if err := v.iam.Client.UpdateResourcesType(ctx, resourceType); err != nil {
			blog.Errorf("update resourceType failed, err: %v, objects: %+v, resourceType: %+v，rid: %s", err, objects,
				resourceType, rid)
			return err
		}
	}

	return nil
}

// registerModelInstanceSelections register instanceSelections for models
func (v *viewer) registerModelInstanceSelections(ctx context.Context, existsSelections []InstanceSelection,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	existMap := make(map[InstanceSelectionID]struct{})
	for _, selection := range existsSelections {
		existMap[selection.ID] = struct{}{}
	}

	instanceSelections := genDynamicInstanceSelections(objects)
	newSelections := make([]InstanceSelection, 0)
	for _, selection := range instanceSelections {
		_, exists := existMap[selection.ID]
		if !exists {
			newSelections = append(newSelections, selection)
			continue
		}

		if err := v.iam.Client.UpdateInstanceSelection(ctx, selection); err != nil {
			blog.Errorf("update instanceSelection %+v failed, err: %v, rid: %s", selection, err, rid)
			return err
		}
	}

	if err := v.iam.Client.RegisterInstanceSelections(ctx, newSelections); err != nil {
		blog.Errorf("register instanceSelections failed, err: %v, objects: %+v, instanceSelections: %+v, rid: %s",
			err, objects, instanceSelections, rid)
		return err
	}

	return nil
}

// unregisterModelInstanceSelections unregister instanceSelections for models
func (v *viewer) unregisterModelInstanceSelections(ctx context.Context, existsSelections []InstanceSelection,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	existMap := make(map[InstanceSelectionID]struct{})
	for _, selection := range existsSelections {
		existMap[selection.ID] = struct{}{}
	}

	instanceSelectionIDs := make([]InstanceSelectionID, 0)
	instanceSelections := genDynamicInstanceSelections(objects)
	for _, instanceSelection := range instanceSelections {
		_, exists := existMap[instanceSelection.ID]
		if exists {
			instanceSelectionIDs = append(instanceSelectionIDs, instanceSelection.ID)
			continue
		}
	}

	if err := v.iam.Client.DeleteInstanceSelections(ctx, instanceSelectionIDs); err != nil {
		blog.Errorf("unregister instanceSelections failed, err: %v, objects: %+v, instanceSelections: %+v, rid: %s",
			err, objects, instanceSelections, rid)
		return err
	}

	return nil
}

// updateModelInstanceSelections update instanceSelections for models
func (v *viewer) updateModelInstanceSelections(ctx context.Context, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	instanceSelections := genDynamicInstanceSelections(objects)
	for _, instanceSelection := range instanceSelections {
		if err := v.iam.Client.UpdateInstanceSelection(ctx, instanceSelection); err != nil {
			blog.Errorf("update instanceSelections failed, err: %v, objects: %+v, instanceSelection: %+v, rid: %s",
				err, objects, instanceSelection, rid)
			return err
		}
	}

	return nil
}

// registerModelActions register actions for models
func (v *viewer) registerModelActions(ctx context.Context, existAction []ResourceAction,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	existMap := make(map[ActionID]struct{})
	for _, action := range existAction {
		existMap[action.ID] = struct{}{}
	}

	actions := genDynamicActions(objects)
	newActions := make([]ResourceAction, 0)
	for _, action := range actions {
		_, exists := existMap[action.ID]
		if !exists {
			newActions = append(newActions, action)
			continue
		}

		if err := v.iam.Client.UpdateAction(ctx, action); err != nil {
			blog.Errorf("update action %+v failed, err: %v, rid: %s", action, err, rid)
			return err
		}
	}

	if err := v.iam.Client.RegisterActions(ctx, newActions); err != nil {
		blog.Errorf("register actions %+v failed, err: %v, objects: %+v, rid: %s", err, objects, newActions, rid)
		return err
	}

	return nil
}

// unregisterModelActions unregister actions for models
func (v *viewer) unregisterModelActions(ctx context.Context, existAction []ResourceAction,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	existMap := make(map[ActionID]struct{})
	for _, action := range existAction {
		existMap[action.ID] = struct{}{}
	}

	actionIDs := make([]ActionID, 0)
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
		if err := v.iam.Client.DeleteActionPolicies(ctx, actionID); err != nil {
			blog.Errorf("delete action %s policies failed, err: %v, rid: %s", actionID, err, rid)
			return err
		}
	}

	if err := v.iam.Client.DeleteActions(ctx, actionIDs); err != nil {
		blog.Errorf("unregister actions failed, err: %v, objects: %+v, actionIDs: %s, rid: %s", err, objects, actionIDs,
			rid)
		return err
	}

	return nil
}

// updateModelActions update actions for models
func (v *viewer) updateModelActions(ctx context.Context, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	actions := genDynamicActions(objects)
	for _, action := range actions {
		if err := v.iam.Client.UpdateAction(ctx, action); err != nil {
			blog.Errorf("update action failed, err: %v, objects: %+v, action: %+v, rid: %s", err, objects, action, rid)
			return err
		}
	}

	return nil
}

// updateModelActionGroups update actionGroups for models
// for now, the update api can only support full update, not incremental update
func (v *viewer) updateModelActionGroups(ctx context.Context, header http.Header, existGroups []ActionGroup) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	objects, err := commonlgc.GetCustomObjects(ctx, header, v.client)
	if err != nil {
		blog.Errorf("get custom objects failed, err: %v, rid: %s", err, rid)
		return err
	}
	actionGroups := GenerateActionGroups(objects)

	if len(existGroups) == 0 {
		if err = v.iam.Client.RegisterActionGroups(ctx, actionGroups); err != nil {
			blog.Errorf("register action groups(%s) failed, err: %v, rid: %s", actionGroups, err, rid)
			return err
		}
		return nil
	}

	if reflect.DeepEqual(existGroups, actionGroups) {
		return nil
	}

	if err = v.iam.Client.UpdateActionGroups(ctx, actionGroups); err != nil {
		blog.Errorf("update actionGroups failed, error: %v, actionGroups: %+v, rid: %s", err, actionGroups, rid)
		return err
	}

	return nil
}
