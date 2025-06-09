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

	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/apimachinery"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	commonlgc "configcenter/src/common/logics"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
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
func (v *viewer) CreateView(ctx context.Context, h http.Header, objects []metadata.Object, redisCli redis.Client,
	rid string) error {

	if !auth.EnableAuthorize() {
		return nil
	}

	locker, err := tryLockRegister(redisCli, rid)
	if err != nil {
		return err
	}
	defer locker.Unlock()

	systemInfo, err := v.iam.Client.GetSystemInfo(ctx, h, []iamtypes.SystemQueryField{})
	if err != nil {
		blog.Errorf("get iam system info failed, err: %v, rid: %s", err, rid)
		return err
	}

	// register order: 1.ResourceType 2.InstanceSelection 3.Action 4.ActionGroup
	if err = v.registerModelResourceTypes(ctx, h, systemInfo.ResourceTypes, objects); err != nil {
		return err
	}

	if err = v.registerModelInstanceSelections(ctx, h, systemInfo.InstanceSelections, objects); err != nil {
		return err
	}

	if err = v.registerModelActions(ctx, h, systemInfo.Actions, objects); err != nil {
		return err
	}

	if err = v.updateModelActionGroups(ctx, h, systemInfo.ActionGroups); err != nil {
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

	systemInfo, err := v.iam.Client.GetSystemInfo(ctx, header, []iamtypes.SystemQueryField{})
	if err != nil {
		blog.Errorf("get iam system info failed, err: %v, rid: %s", err, rid)
		return err
	}

	// unregister order: 1.Action 2.InstanceSelection 3.ResourceType 4.ActionGroup
	if err = v.unregisterModelActions(ctx, header, systemInfo.Actions, objects); err != nil {
		return err
	}

	if err = v.unregisterModelInstanceSelections(ctx, header, systemInfo.InstanceSelections,
		objects); err != nil {
		return err
	}

	if err = v.unregisterModelResourceTypes(ctx, header, systemInfo.ResourceTypes, objects); err != nil {
		return err
	}

	if err = v.updateModelActionGroups(ctx, header, systemInfo.ActionGroups); err != nil {
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

	systemInfo, err := v.iam.Client.GetSystemInfo(ctx, header, []iamtypes.SystemQueryField{})
	if err != nil {
		blog.Errorf("get iam system info failed, err: %v, rid: %s", err, rid)
		return err
	}

	// update order: 1.ResourceType 2.InstanceSelection 3.Action 4.ActionGroup
	if err = v.updateModelResourceTypes(ctx, header, objects); err != nil {
		return err
	}

	if err = v.updateModelInstanceSelections(ctx, header, objects); err != nil {
		return err
	}

	if err = v.updateModelActions(ctx, header, objects); err != nil {
		return err
	}

	if err = v.updateModelActionGroups(ctx, header, systemInfo.ActionGroups); err != nil {
		return err
	}

	return nil
}

// registerModelResourceTypes register resource types for models
func (v *viewer) registerModelResourceTypes(ctx context.Context, h http.Header, existTypes []iam.ResourceType,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	existMap := make(map[iamtypes.TypeID]struct{})
	for _, resourceType := range existTypes {
		existMap[resourceType.ID] = struct{}{}
	}

	tenantID := httpheader.GetTenantID(h)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	resourceTypes := genDynamicResourceTypes(tenantObjects)
	newResTypes := make([]iam.ResourceType, 0)
	for _, resourceType := range resourceTypes {
		_, exists := existMap[resourceType.ID]
		if !exists {
			newResTypes = append(newResTypes, resourceType)
			continue
		}

		if err := v.iam.Client.UpdateResourcesType(ctx, h, resourceType); err != nil {
			blog.Errorf("update resource type failed, err: %v, resourceType: %+v， rid: %s", err, resourceType, rid)
			return err
		}

	}

	if err := v.iam.Client.RegisterResourcesTypes(ctx, h, newResTypes); err != nil {
		blog.Errorf("register resource types failed, err: %v, objects: %+v, resourceTypes: %+v， rid: %s", err,
			tenantObjects, newResTypes, rid)
		return err
	}

	return nil
}

// unregisterModelResourceTypes unregister resourceTypes for models
func (v *viewer) unregisterModelResourceTypes(ctx context.Context, header http.Header, existTypes []iam.ResourceType,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	existMap := make(map[iamtypes.TypeID]struct{})
	for _, resourceType := range existTypes {
		existMap[resourceType.ID] = struct{}{}
	}

	typeIDs := make([]iamtypes.TypeID, 0)
	tenantID := httpheader.GetTenantID(header)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	resourceTypes := genDynamicResourceTypes(tenantObjects)
	for _, resourceType := range resourceTypes {
		_, exists := existMap[resourceType.ID]
		if exists {
			typeIDs = append(typeIDs, resourceType.ID)
		}
	}

	if err := v.iam.Client.DeleteResourcesTypes(ctx, header, typeIDs); err != nil {
		blog.Errorf("unregister resourceTypes failed, err: %v, objects: %+v, resourceTypes: %+v, rid: %s",
			err, tenantObjects, resourceTypes, rid)
		return err
	}

	return nil
}

// updateModelResourceTypes update resource types for models
func (v *viewer) updateModelResourceTypes(ctx context.Context, header http.Header, objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	tenantID := httpheader.GetTenantID(header)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	resourceTypes := genDynamicResourceTypes(tenantObjects)
	for _, resourceType := range resourceTypes {
		if err := v.iam.Client.UpdateResourcesType(ctx, header, resourceType); err != nil {
			blog.Errorf("update resourceType failed, err: %v, objects: %+v, resourceType: %+v，rid: %s", err,
				tenantObjects, resourceType, rid)
			return err
		}
	}

	return nil
}

// registerModelInstanceSelections register instanceSelections for models
func (v *viewer) registerModelInstanceSelections(ctx context.Context, h http.Header,
	existsSelections []iam.InstanceSelection, objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	existMap := make(map[iamtypes.InstanceSelectionID]struct{})
	for _, selection := range existsSelections {
		existMap[selection.ID] = struct{}{}
	}

	tenantID := httpheader.GetTenantID(h)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	instanceSelections := genDynamicInstanceSelections(tenantObjects)
	newSelections := make([]iam.InstanceSelection, 0)
	for _, selection := range instanceSelections {
		_, exists := existMap[selection.ID]
		if !exists {
			newSelections = append(newSelections, selection)
			continue
		}

		if err := v.iam.Client.UpdateInstanceSelection(ctx, h, selection); err != nil {
			blog.Errorf("update instanceSelection %+v failed, err: %v, rid: %s", selection, err, rid)
			return err
		}
	}

	if err := v.iam.Client.RegisterInstanceSelections(ctx, h, newSelections); err != nil {
		blog.Errorf("register instanceSelections failed, err: %v, objects: %+v, instanceSelections: %+v, rid: %s",
			err, tenantObjects, instanceSelections, rid)
		return err
	}

	return nil
}

// unregisterModelInstanceSelections unregister instanceSelections for models
func (v *viewer) unregisterModelInstanceSelections(ctx context.Context, header http.Header,
	existsSelections []iam.InstanceSelection, objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	existMap := make(map[iamtypes.InstanceSelectionID]struct{})
	for _, selection := range existsSelections {
		existMap[selection.ID] = struct{}{}
	}

	instanceSelectionIDs := make([]iamtypes.InstanceSelectionID, 0)
	tenantID := httpheader.GetTenantID(header)
	tenantObjectMap := map[string][]metadata.Object{
		tenantID: objects,
	}
	instanceSelections := genDynamicInstanceSelections(tenantObjectMap)
	for _, instanceSelection := range instanceSelections {
		_, exists := existMap[instanceSelection.ID]
		if exists {
			instanceSelectionIDs = append(instanceSelectionIDs, instanceSelection.ID)
			continue
		}
	}

	if err := v.iam.Client.DeleteInstanceSelections(ctx, header, instanceSelectionIDs); err != nil {
		blog.Errorf("unregister instanceSelections failed, err: %v, objects: %+v, instanceSelections: %+v, rid: %s",
			err, tenantObjectMap, instanceSelections, rid)
		return err
	}

	return nil
}

// updateModelInstanceSelections update instanceSelections for models
func (v *viewer) updateModelInstanceSelections(ctx context.Context, header http.Header,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	tenantID := httpheader.GetTenantID(header)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	instanceSelections := genDynamicInstanceSelections(tenantObjects)
	for _, instanceSelection := range instanceSelections {
		if err := v.iam.Client.UpdateInstanceSelection(ctx, header, instanceSelection); err != nil {
			blog.Errorf("update instanceSelections failed, err: %v, objects: %+v, instanceSelection: %+v, rid: %s",
				err, tenantObjects, instanceSelection, rid)
			return err
		}
	}

	return nil
}

// registerModelActions register actions for models
func (v *viewer) registerModelActions(ctx context.Context, h http.Header, existAction []iam.ResourceAction,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	existMap := make(map[iamtypes.ActionID]struct{})
	for _, action := range existAction {
		existMap[action.ID] = struct{}{}
	}

	tenantID := httpheader.GetTenantID(h)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	actions := genDynamicActions(tenantObjects)
	newActions := make([]iam.ResourceAction, 0)
	for _, action := range actions {
		_, exists := existMap[action.ID]
		if !exists {
			newActions = append(newActions, action)
			continue
		}

		if err := v.iam.Client.UpdateAction(ctx, h, action); err != nil {
			blog.Errorf("update action %+v failed, err: %v, rid: %s", action, err, rid)
			return err
		}
	}

	if err := v.iam.Client.RegisterActions(ctx, h, newActions); err != nil {
		blog.Errorf("register actions %+v failed, err: %v, objects: %+v, rid: %s", err, tenantObjects, newActions, rid)
		return err
	}

	return nil
}

// unregisterModelActions unregister actions for models
func (v *viewer) unregisterModelActions(ctx context.Context, header http.Header, existAction []iam.ResourceAction,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)
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
		if err := v.iam.Client.DeleteActionPolicies(ctx, header, actionID); err != nil {
			blog.Errorf("delete action %s policies failed, err: %v, rid: %s", actionID, err, rid)
			return err
		}
	}

	if err := v.iam.Client.DeleteActions(ctx, header, actionIDs); err != nil {
		blog.Errorf("unregister actions failed, err: %v, objects: %+v, actionIDs: %s, rid: %s", err, objects,
			actionIDs, rid)
		return err
	}

	return nil
}

// updateModelActions update actions for models
func (v *viewer) updateModelActions(ctx context.Context, header http.Header, objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	tenantID := httpheader.GetTenantID(header)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	actions := genDynamicActions(tenantObjects)
	for _, action := range actions {
		if err := v.iam.Client.UpdateAction(ctx, header, action); err != nil {
			blog.Errorf("update action failed, err: %v, objects: %+v, action: %+v, rid: %s", err, tenantObjects, action,
				rid)
			return err
		}
	}

	return nil
}

// updateModelActionGroups update actionGroups for models
// for now, the update api can only support full update, not incremental update
func (v *viewer) updateModelActionGroups(ctx context.Context, h http.Header, existGroups []iam.ActionGroup) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	objects, err := commonlgc.GetCustomObjects(ctx, h, v.client)
	if err != nil {
		blog.Errorf("get custom objects failed, err: %v, rid: %s", err, rid)
		return err
	}

	tenantID := httpheader.GetTenantID(h)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	actionGroups := GenerateActionGroups(tenantObjects)

	if len(existGroups) == 0 {
		if err = v.iam.Client.RegisterActionGroups(ctx, h, actionGroups); err != nil {
			blog.Errorf("register action groups(%s) failed, err: %v, rid: %s", actionGroups, err, rid)
			return err
		}
		return nil
	}

	if reflect.DeepEqual(existGroups, actionGroups) {
		return nil
	}

	if err = v.iam.Client.UpdateActionGroups(ctx, h, actionGroups); err != nil {
		blog.Errorf("update actionGroups failed, error: %v, actionGroups: %+v, rid: %s", err, actionGroups, rid)
		return err
	}

	return nil
}
