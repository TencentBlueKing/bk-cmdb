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
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/lock"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/apigw/iam"
)

// tryLockRegister try lock register iam operation to make sure only one task runs at the same time, retry 3 times.
func tryLockRegister(redisCli redis.Client, rid string) (lock.Locker, error) {
	for i := 0; i < 3; i++ {
		locker := lock.NewLocker(redisCli)
		locked, err := locker.Lock(iamtypes.RegisterIamLock, 2*time.Minute)
		if err != nil {
			blog.Errorf("get register iam lock failed, err: %v, rid: %s", err, rid)
			time.Sleep(5 * time.Second)
			continue
		}

		if locked {
			return locker, nil
		}

		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("there's another register iam task runing, please retry later")
}

// RegisterIamOptions defines options to register iam
type RegisterIamOptions struct {
	Host    string
	Objects map[string][]metadata.Object
}

/**
IAM V4 register order:
  1. register/update System
  2. clear the authorizations of the unused Roles and delete them
  3. unbind the Actions that are no longer used by the remaining Roles, an Action can only be deleted after it is
     unbound from all Roles
  4. delete unused Actions
  5. update ResourceType (use a temporary name when name conflicts)
  6. add ResourceType
  7. update Action (name only; ResourceTypeID change is delete + create)
  8. add Action
  9. delete unused ResourceType
  10. update/register Role and bind the new Actions to them
*/

// Register cc auth resources to iam
func (i IAM) Register(ctx context.Context, h http.Header, redisCli redis.Client, opt *RegisterIamOptions,
	rid string) error {

	if !auth.EnableAuthorize() {
		return nil
	}

	locker, err := tryLockRegister(redisCli, rid)
	if err != nil {
		return err
	}
	defer locker.Unlock()

	if err = i.registerSystem(ctx, h, opt.Host, rid); err != nil {
		return err
	}

	roleCmp, err := i.crossCompareRoles(ctx, h, opt.Objects, rid)
	if err != nil {
		return err
	}

	actionCmp, err := i.crossCompareResActions(ctx, h, opt.Objects, rid)
	if err != nil {
		return err
	}

	resTypeCmp, err := i.crossCompareResTypes(ctx, h, opt.Objects, rid)
	if err != nil {
		return err
	}

	if err = i.removeRoles(ctx, h, roleCmp.removedIDs, rid); err != nil {
		return err
	}

	if err = i.unbindRoleActions(ctx, h, roleCmp.deleteActions, rid); err != nil {
		return err
	}

	if err = i.removeResActions(ctx, h, actionCmp.removedIDs, rid); err != nil {
		return err
	}

	if err = i.updateResTypes(ctx, h, resTypeCmp.updateTypes, rid); err != nil {
		return err
	}

	if err = i.registerResTypes(ctx, h, resTypeCmp.newTypes, rid); err != nil {
		return err
	}

	if err = i.updateResActions(ctx, h, actionCmp.updateActions, rid); err != nil {
		return err
	}

	if err = i.registerResActions(ctx, h, actionCmp.newActions, rid); err != nil {
		return err
	}

	if err = i.removeResTypes(ctx, h, resTypeCmp.removedIDs, rid); err != nil {
		return err
	}

	if err = i.updateRoles(ctx, h, roleCmp.updateRoles, rid); err != nil {
		return err
	}

	if err = i.registerRoles(ctx, h, roleCmp.newRoles, rid); err != nil {
		return err
	}

	if err = i.bindRoleActions(ctx, h, roleCmp.addActions, rid); err != nil {
		return err
	}

	return nil
}

// registerSystem register or update cc system info to IAM
func (i IAM) registerSystem(ctx context.Context, header http.Header, host string, rid string) error {
	systemInfo, err := i.Client.GetSystem(ctx, header)
	if err != nil && !iam.IsSystemNotExistErr(err) {
		blog.Errorf("get system info failed, err: %v, rid: %s", err, rid)
		return err
	}

	sys := &iam.System{
		ID:            iamtypes.SystemIDCMDB,
		Name:          iamtypes.SystemNameCMDB,
		NameEn:        iamtypes.SystemNameCMDBEn,
		Description:   iamtypes.SystemNameCMDB,
		DescriptionEn: iamtypes.SystemNameCMDBEn,
		Managers:      []string{},
		Clients:       []string{iamtypes.SystemIDCMDB},
		CallbackURL:   host,
	}

	if iam.IsSystemNotExistErr(err) {
		if _, err = i.Client.RegisterSystem(ctx, header, sys); err != nil {
			blog.Errorf("register system(%+v) failed, err: %v, rid: %s", sys, err, rid)
			return err
		}
		blog.V(5).Infof("register new system %+v succeed, rid: %s", sys, rid)
		return nil
	}

	if !reflect.DeepEqual(systemInfo, sys) {
		if err = i.Client.UpdateSystem(ctx, header, sys); err != nil {
			blog.Errorf("update system(%+v) failed, err: %v, rid: %s", sys, err, rid)
			return err
		}
		blog.V(5).Infof("update system to %+v succeed, rid: %s", sys, rid)
	}

	return nil
}

// resTypeCompareResult is the cross-compare result of IAM resource types
type resTypeCompareResult struct {
	newTypes    []iam.ResourceType
	updateTypes []iam.ResourceType
	removedIDs  []iamtypes.TypeID
}

// crossCompareResTypes list registered resource types and compare with expected ones, return create/update/delete
// items. Name conflicts are resolved by a temporary name.
func (i IAM) crossCompareResTypes(ctx context.Context, header http.Header, tenantObjects map[string][]metadata.Object,
	rid string) (*resTypeCompareResult, error) {

	registeredResourceTypes, err := i.listAllResTypes(ctx, header, rid)
	if err != nil {
		return nil, err
	}

	registeredResTypeMap := make(map[iamtypes.TypeID]iam.ResourceType)
	for _, resourceType := range registeredResourceTypes {
		registeredResTypeMap[resourceType.ID] = resourceType
	}

	// record the name and resource type id mapping to get the resource types whose name conflicts
	resNameMap := make(map[string]iamtypes.TypeID)
	resNameEnMap := make(map[string]iamtypes.TypeID)
	updateResPrevNameMap := make(map[iamtypes.TypeID]iamName)
	newResTypes := make([]iam.ResourceType, 0)
	updateResTypes := make([]iam.ResourceType, 0)
	for _, resourceType := range GenerateResourceTypes(tenantObjects) {
		resNameMap[resourceType.Name] = resourceType.ID
		resNameEnMap[resourceType.NameEn] = resourceType.ID

		// if current resource type is not registered, register it, otherwise, update it if it is changed
		registeredResType, exists := registeredResTypeMap[resourceType.ID]
		if exists {
			// registered resource type exists in current resource types, should not be removed
			delete(registeredResTypeMap, resourceType.ID)

			if i.isResTypeEqual(registeredResType, resourceType) {
				continue
			}

			updateResPrevNameMap[resourceType.ID] = iamName{
				Name:   registeredResType.Name,
				NameEn: registeredResType.NameEn,
			}
			updateResTypes = append(updateResTypes, resourceType)
			continue
		}
		newResTypes = append(newResTypes, resourceType)
	}

	// if to update resource type previous name conflict with a valid one, change its name to an intermediate one first
	conflictResTypes := getNameConflictItems(updateResTypes, updateResPrevNameMap, resNameMap, resNameEnMap,
		func(data iam.ResourceType) iamtypes.TypeID { return data.ID },
		func(data *iam.ResourceType, name string) { data.Name = name },
		func(data *iam.ResourceType, name string) { data.NameEn = name },
	)

	// remove the resource types that are not exist in new resource types
	removedResTypeIDs := make([]iamtypes.TypeID, 0, len(registeredResTypeMap))
	for resTypeID, resType := range registeredResTypeMap {
		removedResTypeIDs = append(removedResTypeIDs, resTypeID)

		// if to remove resource type name conflicts with a valid one, change its name to an intermediate one first
		isConflict := false
		if _, exists := resNameMap[resType.Name]; exists {
			resType.Name += "_"
			isConflict = true
		}
		if _, exists := resNameEnMap[resType.NameEn]; exists {
			resType.NameEn += "_"
			isConflict = true
		}
		if isConflict {
			conflictResTypes = append(conflictResTypes, resType)
		}
	}

	return &resTypeCompareResult{
		newTypes:    newResTypes,
		updateTypes: append(conflictResTypes, updateResTypes...),
		removedIDs:  removedResTypeIDs,
	}, nil
}

// listAllResTypes list all registered resource types by page
func (i IAM) listAllResTypes(ctx context.Context, header http.Header, rid string) ([]iam.ResourceType, error) {
	resTypes := make([]iam.ResourceType, 0)

	for page := int64(1); ; page++ {
		data, err := i.Client.ListResourceTypes(ctx, header, page, iam.MaxListPageSize)
		if err != nil {
			blog.Errorf("list iam resource types failed, err: %v, page: %d, rid: %s", err, page, rid)
			return nil, err
		}

		resTypes = append(resTypes, data.Results...)
		if len(data.Results) == 0 {
			return resTypes, nil
		}
	}
}

// isResTypeEqual check if registered resource type is equal to expected one
func (i IAM) isResTypeEqual(registeredResType, resType iam.ResourceType) bool {
	if registeredResType.Name != resType.Name || registeredResType.NameEn != resType.NameEn {
		return false
	}
	if len(registeredResType.Ancestors) != len(resType.Ancestors) {
		return false
	}
	for idx, ancestor := range registeredResType.Ancestors {
		if ancestor != resType.Ancestors[idx] {
			return false
		}
	}
	return true
}

// iamName records both chinese and english names of an IAM resource
type iamName struct {
	Name   string
	NameEn string
}

// getNameConflictItems returns copies of update items whose previous names conflict with a currently valid name.
// Conflicting names are replaced with a temporary suffix via updateName so that IAM unique-name constraint is not
// violated during the subsequent actual update.
func getNameConflictItems[T any, ID comparable](items []T, prevNameMap map[ID]iamName, nameMap,
	nameEnMap map[string]ID, getID func(T) ID, updateName, updateEnName func(data *T, name string)) []T {

	conflictItems := make([]T, 0)
	for _, item := range items {
		id := getID(item)
		prevName := prevNameMap[id]
		isConflict := false

		if prevName.Name != "" && nameMap[prevName.Name] != id {
			isConflict = true
			updateName(&item, prevName.Name+"_")
		}
		if prevName.NameEn != "" && nameEnMap[prevName.NameEn] != id {
			isConflict = true
			updateEnName(&item, prevName.NameEn+"_")
		}
		if isConflict {
			conflictItems = append(conflictItems, item)
		}
	}

	return conflictItems
}

// registerResTypes register resource types to IAM
func (i IAM) registerResTypes(ctx context.Context, header http.Header, resTypes []iam.ResourceType, rid string) error {
	if len(resTypes) == 0 {
		return nil
	}

	if _, err := i.Client.RegisterResourcesTypes(ctx, header, resTypes); err != nil {
		blog.Errorf("register resource types(%+v) failed, err: %v, rid: %s", resTypes, err, rid)
		return err
	}

	return nil
}

// updateResTypes update the resource types in IAM, only name and ancestors can be updated
func (i IAM) updateResTypes(ctx context.Context, header http.Header, resTypes []iam.ResourceType, rid string) error {
	for _, resType := range resTypes {
		ancestors := resType.Ancestors
		if ancestors == nil {
			ancestors = make([]iamtypes.TypeID, 0)
		}

		req := &iam.UpdateResourceTypeReq{Name: resType.Name, NameEn: resType.NameEn, Ancestors: ancestors}
		if err := i.Client.UpdateResourcesType(ctx, header, resType.ID, req); err != nil {
			blog.Errorf("update resource type(%+v) failed, err: %v, rid: %s", resType, err, rid)
			return err
		}
	}

	return nil
}

// removeResTypes delete the resource types that are no longer used, IAM only supports deleting one by one
func (i IAM) removeResTypes(ctx context.Context, header http.Header, resTypeIDs []iamtypes.TypeID, rid string) error {
	for _, resTypeID := range resTypeIDs {
		if err := i.Client.DeleteResourcesType(ctx, header, resTypeID); err != nil {
			blog.Errorf("delete resource type(%s) failed, err: %v, rid: %s", resTypeID, err, rid)
			return err
		}
	}

	return nil
}

// resActionCompareResult is the cross-compare result of IAM actions
type resActionCompareResult struct {
	newActions    []iam.ResourceAction
	updateActions []iam.ResourceAction
	removedIDs    []iamtypes.ActionID
}

// crossCompareResActions list registered actions and compare with expected ones. Name change uses Update,
// ResourceTypeID change uses delete + create.
func (i IAM) crossCompareResActions(ctx context.Context, header http.Header,
	tenantObjects map[string][]metadata.Object, rid string) (*resActionCompareResult, error) {

	registeredActions, err := i.listAllResActions(ctx, header, rid)
	if err != nil {
		return nil, err
	}

	registeredResActionMap := make(map[iamtypes.ActionID]iam.ResourceAction)
	for _, resourceAction := range registeredActions {
		registeredResActionMap[resourceAction.ID] = resourceAction
	}

	// record the name and resource action id mapping to get the actions whose name conflicts
	actionNameMap := make(map[string]iamtypes.ActionID)
	actionNameEnMap := make(map[string]iamtypes.ActionID)
	updateActionPrevNameMap := make(map[iamtypes.ActionID]iamName)

	newResActions := make([]iam.ResourceAction, 0)
	updateResActions := make([]iam.ResourceAction, 0)

	for _, resourceAction := range GenerateActions(tenantObjects) {
		actionNameMap[resourceAction.Name] = resourceAction.ID
		actionNameEnMap[resourceAction.NameEn] = resourceAction.ID

		// if current resource action is not registered, register it, otherwise, update it if it is changed
		registeredAction, exists := registeredResActionMap[resourceAction.ID]
		if exists {
			// ResourceTypeID changed: delete the old action and create a new one
			if registeredAction.ResourceTypeID != resourceAction.ResourceTypeID {
				newResActions = append(newResActions, resourceAction)
				continue
			}

			// registered resource action exist in current resource actions, should not be removed
			delete(registeredResActionMap, resourceAction.ID)

			if registeredAction.Name == resourceAction.Name && registeredAction.NameEn == resourceAction.NameEn {
				continue
			}

			updateActionPrevNameMap[resourceAction.ID] = iamName{
				Name:   registeredAction.Name,
				NameEn: registeredAction.NameEn,
			}
			updateResActions = append(updateResActions, resourceAction)
			continue
		}
		newResActions = append(newResActions, resourceAction)
	}

	// if to update action previous name conflict with a valid one, change its name to an intermediate one first
	conflictActions := getNameConflictItems(updateResActions, updateActionPrevNameMap, actionNameMap, actionNameEnMap,
		func(data iam.ResourceAction) iamtypes.ActionID { return data.ID },
		func(data *iam.ResourceAction, name string) { data.Name = name },
		func(data *iam.ResourceAction, name string) { data.NameEn = name },
	)

	removedResActionIDs := make([]iamtypes.ActionID, 0, len(registeredResActionMap))
	for resourceActionID := range registeredResActionMap {
		removedResActionIDs = append(removedResActionIDs, resourceActionID)
	}

	return &resActionCompareResult{
		newActions:    newResActions,
		updateActions: append(conflictActions, updateResActions...),
		removedIDs:    removedResActionIDs,
	}, nil
}

// listAllResActions list all registered actions by page
func (i IAM) listAllResActions(ctx context.Context, header http.Header, rid string) ([]iam.ResourceAction, error) {
	actions := make([]iam.ResourceAction, 0)

	for page := int64(1); ; page++ {
		data, err := i.Client.ListActions(ctx, header, page, iam.MaxListPageSize)
		if err != nil {
			blog.Errorf("list iam actions failed, err: %v, page: %d, rid: %s", err, page, rid)
			return nil, err
		}

		actions = append(actions, data.Results...)
		if len(data.Results) == 0 {
			return actions, nil
		}
	}
}

// registerResActions register actions to IAM
func (i IAM) registerResActions(ctx context.Context, header http.Header, actions []iam.ResourceAction,
	rid string) error {

	if len(actions) == 0 {
		return nil
	}

	if _, err := i.Client.RegisterActions(ctx, header, actions); err != nil {
		blog.Errorf("register actions(%+v) failed, err: %v, rid: %s", actions, err, rid)
		return err
	}

	return nil
}

// updateResActions update the actions in IAM, only the action name can be updated
func (i IAM) updateResActions(ctx context.Context, header http.Header, actions []iam.ResourceAction,
	rid string) error {

	for _, action := range actions {
		req := &iam.UpdateActionReq{Name: action.Name, NameEn: action.NameEn}
		if err := i.Client.UpdateAction(ctx, header, action.ID, req); err != nil {
			blog.Errorf("update action(%+v) failed, err: %v, rid: %s", action, err, rid)
			return err
		}
	}

	return nil
}

// removeResActions delete the actions that are no longer used, an action can only be deleted after it is unbound
// from all the roles. IAM only supports deleting one by one
func (i IAM) removeResActions(ctx context.Context, header http.Header, actionIDs []iamtypes.ActionID,
	rid string) error {

	for _, actionID := range actionIDs {
		if err := i.Client.DeleteAction(ctx, header, actionID); err != nil {
			blog.Errorf("delete action(%s) failed, err: %v, rid: %s", actionID, err, rid)
			return err
		}
	}

	return nil
}

// roleCompareResult is the cross-compare result of IAM roles
type roleCompareResult struct {
	newRoles      []iam.Role
	updateRoles   []iam.Role
	addActions    map[iamtypes.RoleID][]iam.RoleAction
	deleteActions map[iamtypes.RoleID][]iamtypes.ActionID
	removedIDs    []iamtypes.RoleID
}

// crossCompareRoles list registered roles and compare with expected ones, return add/update/action-diff/delete items.
// Name conflicts are resolved by a temporary name.
func (i IAM) crossCompareRoles(ctx context.Context, header http.Header, tenantObjects map[string][]metadata.Object,
	rid string) (*roleCompareResult, error) {

	registeredRoles, err := i.listAllRoles(ctx, header, rid)
	if err != nil {
		return nil, err
	}

	registeredRoleMap := make(map[iamtypes.RoleID]iam.Role)
	for _, role := range registeredRoles {
		registeredRoleMap[role.ID] = role
	}

	// record the name and role id mapping to get the roles whose name conflicts
	roleNameMap := make(map[string]iamtypes.RoleID)
	roleNameEnMap := make(map[string]iamtypes.RoleID)
	updateRolePrevNameMap := make(map[iamtypes.RoleID]iamName)

	addRoleActions := make(map[iamtypes.RoleID][]iam.RoleAction)
	deleteRoleActions := make(map[iamtypes.RoleID][]iamtypes.ActionID)
	newRoles := make([]iam.Role, 0)
	updateRoles := make([]iam.Role, 0)

	for _, role := range GenerateRoles(tenantObjects) {
		roleNameMap[role.Name] = role.ID
		roleNameEnMap[role.NameEn] = role.ID

		// if current role is not registered, register it, otherwise, update it if it is changed
		registered, exists := registeredRoleMap[role.ID]
		if !exists {
			newRoles = append(newRoles, role)
			continue
		}

		// registered role exists in current roles, should not be removed
		delete(registeredRoleMap, role.ID)

		if registered.Name != role.Name || registered.NameEn != role.NameEn ||
			registered.Description != role.Description || registered.DescriptionEn != role.DescriptionEn {
			updateRolePrevNameMap[role.ID] = iamName{
				Name:   registered.Name,
				NameEn: registered.NameEn,
			}
			updateRoles = append(updateRoles, role)
		}

		toAdd, toDelete := diffRoleActions(registered.Actions, role.Actions)
		if len(toAdd) > 0 {
			addRoleActions[role.ID] = toAdd
		}
		if len(toDelete) > 0 {
			deleteRoleActions[role.ID] = toDelete
		}
	}

	// if to update role previous name conflict with a valid one, change its name to an intermediate one first
	conflictRoles := getNameConflictItems(updateRoles, updateRolePrevNameMap, roleNameMap, roleNameEnMap,
		func(data iam.Role) iamtypes.RoleID { return data.ID },
		func(data *iam.Role, name string) { data.Name = name },
		func(data *iam.Role, name string) { data.NameEn = name },
	)

	// remove the roles that are not exist in new roles
	removedRoleIDs := make([]iamtypes.RoleID, 0, len(registeredRoleMap))
	for roleID := range registeredRoleMap {
		removedRoleIDs = append(removedRoleIDs, roleID)
	}

	return &roleCompareResult{
		newRoles:      newRoles,
		updateRoles:   append(conflictRoles, updateRoles...),
		addActions:    addRoleActions,
		deleteActions: deleteRoleActions,
		removedIDs:    removedRoleIDs,
	}, nil
}

// listAllRoles list all registered roles by page
// TODO only list system defined roles after iam support user-defined roles, user cannot update system defined roles.
func (i IAM) listAllRoles(ctx context.Context, header http.Header, rid string) ([]iam.Role, error) {
	roles := make([]iam.Role, 0)

	for page := int64(1); ; page++ {
		data, err := i.Client.ListRoles(ctx, header, page, iam.MaxListPageSize)
		if err != nil {
			blog.Errorf("list iam roles failed, err: %v, page: %d, rid: %s", err, page, rid)
			return nil, err
		}

		roles = append(roles, data.Results...)
		if len(data.Results) == 0 {
			return roles, nil
		}
	}
}

// diffRoleActions compare role action increments, IAM only accepts action ids when unbinding actions from a role
func diffRoleActions(registered, expected []iam.RoleAction) (toAdd []iam.RoleAction, toDelete []iamtypes.ActionID) {
	regMap := make(map[iam.RoleAction]struct{}, len(registered))
	for _, a := range registered {
		regMap[a] = struct{}{}
	}
	expMap := make(map[iam.RoleAction]struct{}, len(expected))
	for _, a := range expected {
		expMap[a] = struct{}{}
	}

	for a := range expMap {
		if _, ok := regMap[a]; !ok {
			toAdd = append(toAdd, a)
		}
	}
	for a := range regMap {
		if _, ok := expMap[a]; !ok {
			toDelete = append(toDelete, a.ID)
		}
	}
	return toAdd, toDelete
}

// registerRoles register roles to IAM
func (i IAM) registerRoles(ctx context.Context, header http.Header, roles []iam.Role, rid string) error {
	if len(roles) == 0 {
		return nil
	}

	if _, err := i.Client.RegisterRoles(ctx, header, roles); err != nil {
		blog.Errorf("register roles(%+v) failed, err: %v, rid: %s", roles, err, rid)
		return err
	}

	return nil
}

// updateRoles update the roles in IAM, only name and description can be updated, the bound actions are updated by
// bindRoleActions and unbindRoleActions
func (i IAM) updateRoles(ctx context.Context, header http.Header, roles []iam.Role, rid string) error {
	for _, role := range roles {
		req := &iam.UpdateRoleReq{
			Name:          role.Name,
			NameEn:        role.NameEn,
			Description:   role.Description,
			DescriptionEn: role.DescriptionEn,
		}
		if err := i.Client.UpdateRole(ctx, header, role.ID, req); err != nil {
			blog.Errorf("update role(%+v) failed, err: %v, rid: %s", role, err, rid)
			return err
		}
	}

	return nil
}

// removeRoles clear the authorizations related to the roles and delete them
func (i IAM) removeRoles(ctx context.Context, header http.Header, roleIDs []iamtypes.RoleID, rid string) error {
	for _, roleID := range roleIDs {
		if err := i.clearRolePolicies(ctx, header, roleID, rid); err != nil {
			return err
		}

		if err := i.Client.DeleteRole(ctx, header, roleID); err != nil {
			blog.Errorf("delete role(%s) failed, err: %v, rid: %s", roleID, err, rid)
			return err
		}
	}

	return nil
}

// clearRolePolicies clear the authorizations related to the role.
// TODO: IAM does not provide the api to clear the role authorizations, need to confirm with IAM how to handle it.
func (i IAM) clearRolePolicies(ctx context.Context, header http.Header, roleID iamtypes.RoleID, rid string) error {
	blog.Infof("clear role(%s) policies is not supported by iam yet, skip it, rid: %s", roleID, rid)
	return nil
}

// bindRoleActions bind the actions to the roles
func (i IAM) bindRoleActions(ctx context.Context, header http.Header, roleActions map[iamtypes.RoleID][]iam.RoleAction,
	rid string) error {

	for roleID, actions := range roleActions {
		if _, err := i.Client.AddRoleActions(ctx, header, roleID, actions); err != nil {
			blog.Errorf("bind role(%s) actions(%+v) failed, err: %v, rid: %s", roleID, actions, err, rid)
			return err
		}
	}

	return nil
}

// unbindRoleActions unbind the actions that are no longer used from the roles
func (i IAM) unbindRoleActions(ctx context.Context, header http.Header,
	roleActionIDs map[iamtypes.RoleID][]iamtypes.ActionID, rid string) error {

	for roleID, actionIDs := range roleActionIDs {
		if err := i.Client.DeleteRoleActions(ctx, header, roleID, actionIDs); err != nil {
			blog.Errorf("unbind role(%s) actions(%+v) failed, err: %v, rid: %s", roleID, actionIDs, err, rid)
			return err
		}
	}

	return nil
}

// SyncIAMSysInstances sync dynamic model resources (ResourceType/Action/Role) between CMDB and IAM
func (i IAM) SyncIAMSysInstances(ctx context.Context, header http.Header, redisCli redis.Client,
	tenantObjects map[string][]metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	for tenantID, objects := range tenantObjects {
		for _, object := range objects {
			if object.ID == 0 || len(object.ObjectID) == 0 || len(object.ObjectName) == 0 {
				blog.Errorf("sync iam system instances but object(%#v) is invalid, tenantID: %s, rid: %s",
					object, tenantID, rid)
				return errors.New("sync iam instances, but object is invalid")
			}
		}
	}

	locker, err := tryLockRegister(redisCli, rid)
	if err != nil {
		return err
	}
	defer locker.Unlock()

	addedRoles, deletedRoleIDs, err := i.compareRoles(ctx, header, tenantObjects, rid)
	if err != nil {
		return err
	}
	addedResourceTypes, deletedResourceTypes, err := i.compareResourceTypes(ctx, header, tenantObjects, rid)
	if err != nil {
		return err
	}
	addedActions, deletedActions, err := i.compareActions(ctx, header, tenantObjects, rid)
	if err != nil {
		return err
	}

	// delete order: Role -> Action -> ResourceType
	if err = i.removeRoles(ctx, header, deletedRoleIDs, rid); err != nil {
		return err
	}
	if err = i.removeResActions(ctx, header, deletedActions, rid); err != nil {
		return err
	}
	if err = i.removeResTypes(ctx, header, deletedResourceTypes, rid); err != nil {
		return err
	}

	// add order: ResourceType -> Action -> Role
	if err = i.registerResTypes(ctx, header, addedResourceTypes, rid); err != nil {
		return err
	}
	if err = i.registerResActions(ctx, header, addedActions, rid); err != nil {
		return err
	}
	if err = i.registerRoles(ctx, header, addedRoles, rid); err != nil {
		return err
	}

	return nil
}

// compareRoles list registered roles and compare with dynamic roles
func (i IAM) compareRoles(ctx context.Context, header http.Header, tenantObjects map[string][]metadata.Object,
	rid string) (added []iam.Role, deleted []iamtypes.RoleID, err error) {

	iamRoles, err := i.listAllRoles(ctx, header, rid)
	if err != nil {
		blog.Errorf("list iam roles failed, err: %v, rid: %s", err, rid)
		return nil, nil, err
	}

	iamRoleMap := map[iamtypes.RoleID]struct{}{}
	for _, role := range iamRoles {
		if isIAMSysInstanceRole(role.ID) {
			iamRoleMap[role.ID] = struct{}{}
		}
	}
	for _, role := range genDynamicInstRoles(tenantObjects) {
		if _, ok := iamRoleMap[role.ID]; !ok {
			added = append(added, role)
		} else {
			delete(iamRoleMap, role.ID)
		}
	}
	for roleID := range iamRoleMap {
		deleted = append(deleted, roleID)
	}
	return added, deleted, nil
}

// compareActions list registered actions and compare with dynamic model actions
func (i IAM) compareActions(ctx context.Context, header http.Header, tenantObjects map[string][]metadata.Object,
	rid string) (addedActions []iam.ResourceAction, deletedActionIDs []iamtypes.ActionID, err error) {

	iamActions, err := i.listAllResActions(ctx, header, rid)
	if err != nil {
		blog.Errorf("list iam actions failed, err: %v, rid: %s", err, rid)
		return nil, nil, err
	}

	iamActionMap := map[iamtypes.ActionID]struct{}{}
	for _, action := range iamActions {
		if isIAMSysInstanceAction(action.ID) {
			iamActionMap[action.ID] = struct{}{}
		}
	}

	for _, action := range genDynamicActions(tenantObjects) {
		if _, ok := iamActionMap[action.ID]; !ok {
			addedActions = append(addedActions, action)
		} else {
			delete(iamActionMap, action.ID)
		}
	}

	for actionID := range iamActionMap {
		deletedActionIDs = append(deletedActionIDs, actionID)
	}

	return addedActions, deletedActionIDs, nil
}

// compareResourceTypes list registered resource types and compare with dynamic model resource types
func (i IAM) compareResourceTypes(ctx context.Context, header http.Header,
	tenantObjects map[string][]metadata.Object, rid string) (addedResourceTypes []iam.ResourceType,
	deletedTypeIDs []iamtypes.TypeID, err error) {

	iamResourceTypes, err := i.listAllResTypes(ctx, header, rid)
	if err != nil {
		blog.Errorf("list iam resource types failed, err: %v, rid: %s", err, rid)
		return nil, nil, err
	}

	iamResourceTypeMap := map[iamtypes.TypeID]struct{}{}
	for _, resourceType := range iamResourceTypes {
		if IsIAMSysInstance(resourceType.ID) {
			iamResourceTypeMap[resourceType.ID] = struct{}{}
		}
	}

	for _, resourceType := range genDynamicResourceTypes(tenantObjects) {
		if _, ok := iamResourceTypeMap[resourceType.ID]; !ok {
			addedResourceTypes = append(addedResourceTypes, resourceType)
		} else {
			delete(iamResourceTypeMap, resourceType.ID)
		}
	}

	for typeID := range iamResourceTypeMap {
		deletedTypeIDs = append(deletedTypeIDs, typeID)
	}

	return addedResourceTypes, deletedTypeIDs, nil
}
