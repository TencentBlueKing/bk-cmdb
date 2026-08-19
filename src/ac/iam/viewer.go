/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package iam

import (
	"context"
	"net/http"

	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/apimachinery"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
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

	if err = v.registerModelResourceTypes(ctx, h, objects); err != nil {
		return err
	}

	if err = v.registerModelActions(ctx, h, objects); err != nil {
		return err
	}

	if err = v.registerModelRoles(ctx, h, objects); err != nil {
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

	if err = v.unregisterModelRoles(ctx, header, objects); err != nil {
		return err
	}

	if err = v.unregisterModelActions(ctx, header, objects); err != nil {
		return err
	}

	if err = v.unregisterModelResourceTypes(ctx, header, objects); err != nil {
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

	if err = v.updateModelResourceTypes(ctx, header, objects); err != nil {
		return err
	}

	if err = v.updateModelActions(ctx, header, objects); err != nil {
		return err
	}

	if err = v.updateModelRoles(ctx, header, objects); err != nil {
		return err
	}

	return nil
}

// registerModelResourceTypes register resource types for models
func (v *viewer) registerModelResourceTypes(ctx context.Context, h http.Header, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	existTypes, err := v.iam.listAllResTypes(ctx, h, rid)
	if err != nil {
		blog.Errorf("list iam resource types failed, err: %v, rid: %s", err, rid)
		return err
	}

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
	updateResTypes := make([]iam.ResourceType, 0)
	for _, resourceType := range resourceTypes {
		_, exists := existMap[resourceType.ID]
		if !exists {
			newResTypes = append(newResTypes, resourceType)
			continue
		}
		updateResTypes = append(updateResTypes, resourceType)
	}

	if err = v.iam.updateResTypes(ctx, h, updateResTypes, rid); err != nil {
		return err
	}

	return v.iam.registerResTypes(ctx, h, newResTypes, rid)
}

// unregisterModelResourceTypes unregister resourceTypes for models
func (v *viewer) unregisterModelResourceTypes(ctx context.Context, header http.Header,
	objects []metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	existTypes, err := v.iam.listAllResTypes(ctx, header, rid)
	if err != nil {
		blog.Errorf("list iam resource types failed, err: %v, rid: %s", err, rid)
		return err
	}

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

	return v.iam.removeResTypes(ctx, header, typeIDs, rid)
}

// updateModelResourceTypes update model resource types
func (v *viewer) updateModelResourceTypes(ctx context.Context, header http.Header, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	tenantID := httpheader.GetTenantID(header)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	return v.iam.updateResTypes(ctx, header, genDynamicResourceTypes(tenantObjects), rid)
}

// registerModelActions register actions for models
func (v *viewer) registerModelActions(ctx context.Context, h http.Header, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	existActions, err := v.iam.listAllResActions(ctx, h, rid)
	if err != nil {
		blog.Errorf("list iam actions failed, err: %v, rid: %s", err, rid)
		return err
	}

	existMap := make(map[iamtypes.ActionID]struct{})
	for _, action := range existActions {
		existMap[action.ID] = struct{}{}
	}

	tenantID := httpheader.GetTenantID(h)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	actions := genDynamicActions(tenantObjects)
	newActions := make([]iam.ResourceAction, 0)
	updateActions := make([]iam.ResourceAction, 0)
	for _, action := range actions {
		if _, exists := existMap[action.ID]; !exists {
			newActions = append(newActions, action)
			continue
		}
		updateActions = append(updateActions, action)
	}

	if err := v.iam.updateResActions(ctx, h, updateActions, rid); err != nil {
		return err
	}

	return v.iam.registerResActions(ctx, h, newActions, rid)
}

// unregisterModelActions unregister model actions
func (v *viewer) unregisterModelActions(ctx context.Context, header http.Header, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	existActions, err := v.iam.listAllResActions(ctx, header, rid)
	if err != nil {
		blog.Errorf("list iam actions failed, err: %v, rid: %s", err, rid)
		return err
	}

	existMap := make(map[iamtypes.ActionID]struct{})
	for _, action := range existActions {
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

	return v.iam.removeResActions(ctx, header, actionIDs, rid)
}

// updateModelActions update model actions (name only)
func (v *viewer) updateModelActions(ctx context.Context, header http.Header, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	tenantID := httpheader.GetTenantID(header)
	tenantObjects := map[string][]metadata.Object{
		tenantID: objects,
	}
	return v.iam.updateResActions(ctx, header, genDynamicActions(tenantObjects), rid)
}

// registerModelRoles register model manager and viewer roles
func (v *viewer) registerModelRoles(ctx context.Context, h http.Header, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	existRoles, err := v.iam.listAllRoles(ctx, h, rid)
	if err != nil {
		blog.Errorf("list iam roles failed, err: %v, rid: %s", err, rid)
		return err
	}

	existMap := make(map[iamtypes.RoleID]iam.Role)
	for _, role := range existRoles {
		existMap[role.ID] = role
	}

	tenantID := httpheader.GetTenantID(h)
	tenantObjects := map[string][]metadata.Object{tenantID: objects}
	roles := genDynamicInstRoles(tenantObjects)

	newRoles := make([]iam.Role, 0)
	updateRoles := make([]iam.Role, 0)
	addActions := make(map[iamtypes.RoleID][]iam.RoleAction)
	deleteActions := make(map[iamtypes.RoleID][]iamtypes.ActionID)
	for _, role := range roles {
		registered, exists := existMap[role.ID]
		if !exists {
			newRoles = append(newRoles, role)
			continue
		}
		updateRoles = append(updateRoles, role)

		toAdd, toDelete := diffRoleActions(registered.Actions, role.Actions)
		if len(toAdd) > 0 {
			addActions[role.ID] = toAdd
		}
		if len(toDelete) > 0 {
			deleteActions[role.ID] = toDelete
		}
	}

	if err = v.iam.unbindRoleActions(ctx, h, deleteActions, rid); err != nil {
		return err
	}

	if err = v.iam.updateRoles(ctx, h, updateRoles, rid); err != nil {
		return err
	}

	if err = v.iam.registerRoles(ctx, h, newRoles, rid); err != nil {
		return err
	}

	return v.iam.bindRoleActions(ctx, h, addActions, rid)
}

// unregisterModelRoles unregister model manager and viewer roles
func (v *viewer) unregisterModelRoles(ctx context.Context, header http.Header, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	existRoles, err := v.iam.listAllRoles(ctx, header, rid)
	if err != nil {
		blog.Errorf("list iam roles failed, err: %v, rid: %s", err, rid)
		return err
	}

	existMap := make(map[iamtypes.RoleID]struct{})
	for _, role := range existRoles {
		existMap[role.ID] = struct{}{}
	}

	roleIDs := make([]iamtypes.RoleID, 0)
	for _, obj := range objects {
		for _, roleID := range []iamtypes.RoleID{
			GenIAMDynamicRoleID(obj.ID, "manager"),
			GenIAMDynamicRoleID(obj.ID, "viewer"),
		} {
			if _, exists := existMap[roleID]; exists {
				roleIDs = append(roleIDs, roleID)
			}
		}
	}

	return v.iam.removeRoles(ctx, header, roleIDs, rid)
}

// updateModelRoles update model manager and viewer roles
func (v *viewer) updateModelRoles(ctx context.Context, header http.Header, objects []metadata.Object) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	existRoles, err := v.iam.listAllRoles(ctx, header, rid)
	if err != nil {
		blog.Errorf("list iam roles failed, err: %v, rid: %s", err, rid)
		return err
	}

	existMap := make(map[iamtypes.RoleID]iam.Role)
	for _, role := range existRoles {
		existMap[role.ID] = role
	}

	tenantID := httpheader.GetTenantID(header)
	tenantObjects := map[string][]metadata.Object{tenantID: objects}

	updateRoles := make([]iam.Role, 0)
	addActions := make(map[iamtypes.RoleID][]iam.RoleAction)
	deleteActions := make(map[iamtypes.RoleID][]iamtypes.ActionID)
	for _, role := range genDynamicInstRoles(tenantObjects) {
		registered, exists := existMap[role.ID]
		if !exists {
			continue
		}
		updateRoles = append(updateRoles, role)

		toAdd, toDelete := diffRoleActions(registered.Actions, role.Actions)
		if len(toAdd) > 0 {
			addActions[role.ID] = toAdd
		}
		if len(toDelete) > 0 {
			deleteActions[role.ID] = toDelete
		}
	}

	if err = v.iam.unbindRoleActions(ctx, header, deleteActions, rid); err != nil {
		return err
	}

	if err = v.iam.updateRoles(ctx, header, updateRoles, rid); err != nil {
		return err
	}

	return v.iam.bindRoleActions(ctx, header, addActions, rid)
}
