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

// Package iam TODO
package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/ac/meta"
	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/authserver"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/lock"
	"configcenter/src/common/metadata"
	apigwcli "configcenter/src/common/resource/apigw"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/sdk/types"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/apigw/iam"
)

// IAM TODO
type IAM struct {
	Client iam.ClientI
}

// NewIAM new iam client
func NewIAM() (*IAM, error) {
	return &IAM{
		Client: apigwcli.Client().Iam(),
	}, nil
}

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
1. 资源间的依赖关系为 Action 依赖 InstanceSelection 依赖 ResourceType，对资源的增删改操作需要按照这个依赖顺序调整
2. ActionGroup、ResCreatorAction、CommonAction 依赖于 Action，这些资源的增删操作始终放在最后
3. 因为资源的名称在系统中是唯一的，所以可能遇到循环依赖的情况（如两个资源分别更新成对方的名字），此时需要引入一个中间变量进行二次更新

综上，具体操作顺序如下：
  1. 注册cc系统信息
  2. 删除Action。该操作无依赖
  3. 更新ResourceType，先更新名字冲突的(包括需要删除的)为中间值，再更新其它的。该操作无依赖
  4. 新增ResourceType。该操作依赖于上一步中同名的ResourceType均已更新
  5. 更新InstanceSelection，先更新名字冲突的(包括需要删除的)为中间值，再更新其它的。该操作依赖于上一步中的ResourceType均已新增
  6. 新增InstanceSelection。该操作依赖于上一步中同名的InstanceSelection均已更新+第4步中的ResourceType均已新增
  7. 更新ResourceAction，先更新名字冲突的为中间值，再更新其它的。该操作依赖于第2步中同名Action已删除+上一步中InstanceSelection已新增
  8. 新增ResourceAction。该操作依赖于上一步中同名的ResourceAction均已更新+第6步中的InstanceSelection均已新增
  9. 删除InstanceSelection。该操作依赖于第2步和第7步中的原本依赖了这些InstanceSelection的Action均已删除和更新
 10. 删除ResourceType。该操作依赖于第5步和第9步中的原本依赖了这些ResourceType的InstanceSelection均已删除和更新
 11. 注册ActionGroup、ResCreatorAction、CommonAction信息
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

	registeredInfo, err := i.registerSystem(ctx, h, opt.Host)
	if err != nil {
		return err
	}

	newResTypes, updateResTypes, removedResTypeIDs := i.crossCompareResTypes(registeredInfo.ResourceTypes,
		opt.Objects)
	newInstSelections, updateInstSelections, removedInstSelectionIDs := i.crossCompareInstSelections(
		registeredInfo.InstanceSelections, opt.Objects)
	newResActions, updateResActions, removedResActionIDs := i.crossCompareResActions(registeredInfo.Actions,
		opt.Objects)

	if err = i.removeResActions(ctx, h, removedResActionIDs, rid); err != nil {
		return err
	}

	for _, resourceType := range updateResTypes {
		if err = i.Client.UpdateResourcesType(ctx, h, resourceType); err != nil {
			blog.Errorf("update resource type(%v) failed, err: %v, rid: %s", resourceType, err, rid)
			return err
		}
	}

	if err = i.Client.RegisterResourcesTypes(ctx, h, newResTypes); err != nil {
		blog.Errorf("register resource types(%v) failed, err: %v, rid: %s", newResTypes, err, rid)
		return err
	}

	for _, instanceSelection := range updateInstSelections {
		if err = i.Client.UpdateInstanceSelection(ctx, h, instanceSelection); err != nil {
			blog.Errorf("update instance selection(%v) failed, err: %v, rid: %s", instanceSelection, err, rid)
			return err
		}
	}

	if err = i.Client.RegisterInstanceSelections(ctx, h, newInstSelections); err != nil {
		blog.Errorf("register instance selections(%v) failed, err: %v, rid: %s", newInstSelections, err, rid)
		return err
	}

	for _, resourceAction := range updateResActions {
		if err = i.Client.UpdateAction(ctx, h, resourceAction); err != nil {
			blog.Errorf("update resource action(%v) failed, err: %v, rid: %s", resourceAction, err, rid)
			return err
		}
	}

	if err = i.Client.RegisterActions(ctx, h, newResActions); err != nil {
		blog.Errorf("register resource actions(%v) failed, err: %v, rid: %s", newResActions, err, rid)
		return err
	}

	if err = i.Client.DeleteInstanceSelections(ctx, h, removedInstSelectionIDs); err != nil {
		blog.Errorf("delete instance selections(%v) failed, err: %v, rid: %s", removedInstSelectionIDs, err, rid)
		return err
	}

	if err = i.Client.DeleteResourcesTypes(ctx, h, removedResTypeIDs); err != nil {
		blog.Errorf("delete resource types(%v) failed, err: %v, rid: %s", removedResTypeIDs, err, rid)
		return err
	}

	if err := i.registerActionGroups(ctx, h, registeredInfo, opt.Objects, rid); err != nil {
		blog.Errorf("register action groups(%v) failed, err: %v, rid: %s", registeredInfo, err, rid)
		return err
	}

	if err := i.registerResCreatorActions(ctx, h, registeredInfo, rid); err != nil {
		blog.Errorf("register resCreator actions(%v) failed, tenantID: %s, err: %v, rid: %s", registeredInfo, err, rid)
		return err
	}

	if err := i.registerCommonActions(ctx, h, registeredInfo, rid); err != nil {
		blog.Errorf("register common actions(%v) failed, tenantID: %s, err: %v, rid: %s", registeredInfo, err, rid)
		return err
	}

	return nil
}

// registerSystem register cc system to iam
func (i IAM) registerSystem(ctx context.Context, header http.Header, host string) (*iam.RegisteredSystemInfo, error) {

	systemInfo, err := i.Client.GetSystemInfo(ctx, header, []iamtypes.SystemQueryField{})
	if err != nil && err != iam.ErrNotFound {
		blog.Errorf("get system info failed, err: %v", err)
		return nil, err
	}

	// if iam cmdb system has not been registered, register system
	if err == iam.ErrNotFound {
		sys := iam.System{
			ID:          iamtypes.SystemIDCMDB,
			Name:        iamtypes.SystemNameCMDB,
			EnglishName: iamtypes.SystemNameCMDBEn,
			Clients:     iamtypes.SystemIDCMDB,
			ProviderConfig: &iam.SysConfig{
				Host: host,
				Auth: "basic",
			},
		}

		if err = i.Client.RegisterSystem(ctx, header, sys); err != nil {
			blog.Errorf("register system(%s) failed, err: %v", sys, err)
			return nil, err
		}

		blog.V(5).Infof("register new system %+v succeed", sys)
		return new(iam.RegisteredSystemInfo), nil
	}

	providerConfig := systemInfo.BaseInfo.ProviderConfig
	if providerConfig == nil || providerConfig.Host != host {
		// if iam registered cmdb system has no ProviderConfig
		// or registered host config is different with current host config, update system host config
		if err = i.Client.UpdateSystemConfig(ctx, header, &iam.SysConfig{Host: host}); err != nil {
			blog.Errorf("update system host %s config failed, err: %v", host, err)
			return nil, err
		}
		if providerConfig == nil {
			blog.V(5).Infof("update system host to %s succeed", host)
		} else {
			blog.V(5).Infof("update system host %s to %s succeed", providerConfig.Host, host)
		}
	}

	return systemInfo, nil
}

// iamName record iam name and english name to find if name conflicts
type iamName struct {
	Name   string
	NameEn string
}

// crossCompareResTypes cross compare resource types to get need create/update/delete ones
func (i IAM) crossCompareResTypes(registeredResourceTypes []iam.ResourceType,
	tenantObjects map[string][]metadata.Object) ([]iam.ResourceType, []iam.ResourceType, []iamtypes.TypeID) {

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
		// if current resource type is not registered, register it, otherwise, update it if its version is changed
		registeredResType, exists := registeredResTypeMap[resourceType.ID]
		if exists {
			// registered resource type exists in current resource types, should not be removed
			delete(registeredResTypeMap, resourceType.ID)
			if i.compareResType(registeredResType, resourceType) {
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
	conflictResTypes := make([]iam.ResourceType, 0)
	for _, updateResType := range updateResTypes {
		prevName := updateResPrevNameMap[updateResType.ID]
		isConflict := false
		if resNameMap[prevName.Name] != updateResType.ID {
			isConflict = true
			updateResType.Name = prevName.Name + "_"
		}
		if resNameEnMap[prevName.NameEn] != updateResType.ID {
			isConflict = true
			updateResType.NameEn = prevName.NameEn + "_"
		}
		if isConflict {
			conflictResTypes = append(conflictResTypes, updateResType)
		}
	}
	// remove the resource types that are not exist in new resource types
	removedResTypeIDs := make([]iamtypes.TypeID, len(registeredResTypeMap))
	idx := 0
	for resTypeID, resType := range registeredResTypeMap {
		removedResTypeIDs[idx] = resTypeID
		idx++
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
			if resType.Version == 0 {
				resType.Version = 1
			}
			conflictResTypes = append(conflictResTypes, resType)
		}
	}
	return newResTypes, append(conflictResTypes, updateResTypes...), removedResTypeIDs
}

// compareResType compare if registered resource type that iam returns is the same with the new resource type
func (i IAM) compareResType(registeredResType, resType iam.ResourceType) bool {
	if registeredResType.ID != resType.ID ||
		registeredResType.Name != resType.Name ||
		registeredResType.NameEn != resType.NameEn ||
		registeredResType.Description != resType.Description ||
		registeredResType.DescriptionEn != resType.DescriptionEn ||
		registeredResType.Version < resType.Version ||
		registeredResType.ProviderConfig.Path != resType.ProviderConfig.Path {
		return false
	}

	if len(registeredResType.Parents) != len(resType.Parents) {
		return false
	}
	for idx, parent := range registeredResType.Parents {
		resTypeParent := resType.Parents[idx]
		if parent.ResourceID != resTypeParent.ResourceID || parent.SystemID != resTypeParent.SystemID {
			return false
		}
	}

	return true
}

// crossCompareInstSelections cross compare instance selections to get need create/update/delete ones
func (i IAM) crossCompareInstSelections(registeredInstanceSelections []iam.InstanceSelection,
	tenantObjects map[string][]metadata.Object) ([]iam.InstanceSelection, []iam.InstanceSelection,
	[]iamtypes.InstanceSelectionID) {

	registeredInstSelectionMap := make(map[iamtypes.InstanceSelectionID]iam.InstanceSelection)
	for _, instanceSelection := range registeredInstanceSelections {
		registeredInstSelectionMap[instanceSelection.ID] = instanceSelection
	}

	// record the name and instance selection id mapping to get the instance selections whose name conflicts
	selectionNameMap := make(map[string]iamtypes.InstanceSelectionID)
	selectionNameEnMap := make(map[string]iamtypes.InstanceSelectionID)
	updateSelectionPrevNameMap := make(map[iamtypes.InstanceSelectionID]iamName)

	newInstSelections := make([]iam.InstanceSelection, 0)
	updateInstSelections := make([]iam.InstanceSelection, 0)

	for _, instanceSelection := range GenerateInstanceSelections(tenantObjects) {
		selectionNameMap[instanceSelection.Name] = instanceSelection.ID
		selectionNameEnMap[instanceSelection.NameEn] = instanceSelection.ID

		selection, exists := registeredInstSelectionMap[instanceSelection.ID]

		// if current instance selection is not registered, register it, otherwise, update it if it is changed
		if exists {
			// registered instance selection exists in current instance selections, should not be removed
			delete(registeredInstSelectionMap, instanceSelection.ID)

			if reflect.DeepEqual(selection, instanceSelection) {
				continue
			}

			updateSelectionPrevNameMap[instanceSelection.ID] = iamName{
				Name:   selection.Name,
				NameEn: selection.NameEn,
			}
			updateInstSelections = append(updateInstSelections, instanceSelection)
			continue
		}
		newInstSelections = append(newInstSelections, instanceSelection)
	}
	// if to update selection previous name conflict with a valid one, change its name to an intermediate one first
	conflictSelections := make([]iam.InstanceSelection, 0)
	for _, updateSelection := range updateInstSelections {
		prevName := updateSelectionPrevNameMap[updateSelection.ID]
		isConflict := false
		if selectionNameMap[prevName.Name] != updateSelection.ID {
			updateSelection.Name = prevName.Name + "_"
			isConflict = true
		}
		if selectionNameEnMap[prevName.NameEn] != updateSelection.ID {
			updateSelection.NameEn = prevName.NameEn + "_"
			isConflict = true
		}
		if isConflict {
			conflictSelections = append(conflictSelections, updateSelection)
		}
	}
	// remove the resource types that are not exist in new resource types
	removedInstSelectionIDs := make([]iamtypes.InstanceSelectionID, len(registeredInstSelectionMap))
	idx := 0
	for selectionID, selection := range registeredInstSelectionMap {
		removedInstSelectionIDs[idx] = selectionID
		idx++
		// if to remove selection name conflicts with a valid one, change its name to an intermediate one first
		isConflict := false
		if _, exists := selectionNameMap[selection.Name]; exists {
			selection.Name += "_"
			isConflict = true
		}
		if _, exists := selectionNameEnMap[selection.NameEn]; exists {
			selection.NameEn += "_"
			isConflict = true
		}
		if isConflict {
			conflictSelections = append(conflictSelections, selection)
		}
	}
	return newInstSelections, append(conflictSelections, updateInstSelections...), removedInstSelectionIDs
}

// crossCompareResActions cross compare resource actions to get need create/update/delete ones
func (i IAM) crossCompareResActions(registeredActions []iam.ResourceAction,
	tenantObjects map[string][]metadata.Object) ([]iam.ResourceAction, []iam.ResourceAction, []iamtypes.ActionID) {

	registeredResActionMap := make(map[iamtypes.ActionID]iam.ResourceAction)
	for _, resourceAction := range registeredActions {
		registeredResActionMap[resourceAction.ID] = resourceAction
	}

	// record the name and resource action id mapping to get the instance selections whose name conflicts
	actionNameMap := make(map[string]iamtypes.ActionID)
	actionNameEnMap := make(map[string]iamtypes.ActionID)
	updateActionPrevNameMap := make(map[iamtypes.ActionID]iamName)

	newResActions := make([]iam.ResourceAction, 0)
	updateResActions := make([]iam.ResourceAction, 0)

	for _, resourceAction := range GenerateActions(tenantObjects) {
		actionNameMap[resourceAction.Name] = resourceAction.ID
		actionNameEnMap[resourceAction.NameEn] = resourceAction.ID

		// if current resource action is not registered, register it, otherwise, update it if its version is changed
		action, exists := registeredResActionMap[resourceAction.ID]
		if exists {
			// registered resource action exist in current resource actions, should not be removed
			delete(registeredResActionMap, resourceAction.ID)

			if i.compareResAction(action, resourceAction) {
				continue
			}

			updateActionPrevNameMap[action.ID] = iamName{
				Name:   action.Name,
				NameEn: action.NameEn,
			}
			updateResActions = append(updateResActions, resourceAction)
			continue
		}
		newResActions = append(newResActions, resourceAction)
	}

	// if to update action previous name conflict with a valid one, change its name to an intermediate one first
	conflictActions := make([]iam.ResourceAction, 0)
	for _, updateAction := range updateResActions {
		prevName := updateActionPrevNameMap[updateAction.ID]
		isConflict := false

		if actionNameMap[prevName.Name] != updateAction.ID {
			updateAction.Name = prevName.Name + "_"
			isConflict = true
		}

		if actionNameEnMap[prevName.NameEn] != updateAction.ID {
			updateAction.NameEn = prevName.NameEn + "_"
			isConflict = true
		}

		if isConflict {
			conflictActions = append(conflictActions, updateAction)
		}
	}

	removedResActionIDs := make([]iamtypes.ActionID, len(registeredResActionMap))
	idx := 0
	for resourceActionID := range registeredResActionMap {
		removedResActionIDs[idx] = resourceActionID
		idx++
	}

	return newResActions, append(conflictActions, updateResActions...), removedResActionIDs
}

// compareResAction compare if registered resource action that iam returns is the same with the new resource action
func (i IAM) compareResAction(registeredAction, action iam.ResourceAction) bool {
	if registeredAction.ID != action.ID ||
		registeredAction.Name != action.Name ||
		registeredAction.NameEn != action.NameEn ||
		registeredAction.Type != action.Type ||
		registeredAction.Version < action.Version ||
		registeredAction.Hidden != action.Hidden {
		return false
	}

	if len(registeredAction.RelatedResourceTypes) != len(action.RelatedResourceTypes) {
		return false
	}

	for idx, registeredResType := range registeredAction.RelatedResourceTypes {
		resType := action.RelatedResourceTypes[idx]

		isEqual := i.compareResActionType(resType, registeredResType)
		if !isEqual {
			return false
		}

		// TODO since iam returns no related selections & we use matching type & selection, skip this comparison
	}

	if len(registeredAction.RelatedActions) != len(action.RelatedActions) {
		return false
	}

	for idx, actionID := range registeredAction.RelatedActions {
		if actionID != action.RelatedActions[idx] {
			return false
		}
	}

	return true
}

// compareResActionType compare if registered and new resource action's related type are the same
func (i IAM) compareResActionType(resType iam.RelateResourceType, registeredResType iam.RelateResourceType) bool {
	// iam default selection mode is "instance"
	if resType.SelectionMode == "" {
		resType.SelectionMode = iamtypes.ModeInstance
	}

	if registeredResType.ID != resType.ID || registeredResType.SelectionMode != resType.SelectionMode {
		return false
	}

	if registeredResType.Scope == nil && resType.Scope == nil {
		return true
	}

	if registeredResType.Scope == nil && resType.Scope != nil ||
		registeredResType.Scope != nil && resType.Scope == nil {
		return false
	}

	if registeredResType.Scope.Op != resType.Scope.Op {
		return false
	}

	if len(registeredResType.Scope.Content) != len(resType.Scope.Content) {
		return false
	}

	for index, registeredContent := range registeredResType.Scope.Content {
		content := resType.Scope.Content[index]
		if registeredContent.Op != content.Op || registeredContent.Value != content.Value ||
			registeredContent.Field != content.Field {
			return false
		}
	}
	return true
}

// removeResActions remove resource actions and related policies
func (i IAM) removeResActions(ctx context.Context, h http.Header, actionIDs []iamtypes.ActionID, rid string) error {
	if len(actionIDs) == 0 {
		return nil
	}

	// before deleting action, the dependent action policies must be deleted
	for _, resourceActionID := range actionIDs {
		if err := i.Client.DeleteActionPolicies(ctx, h, resourceActionID); err != nil {
			blog.Errorf("delete action %s policies failed, err: %v, rid: %s", resourceActionID, err, rid)
			return err
		}
	}

	if err := i.Client.DeleteActions(ctx, h, actionIDs); err != nil {
		blog.Errorf("delete resource actions(%+v) failed, err: %v, rid: %s", actionIDs, err, rid)
		return err
	}

	return nil
}

// registerActionGroups register or update resource action groups
func (i IAM) registerActionGroups(ctx context.Context, h http.Header, registeredInfo *iam.RegisteredSystemInfo,
	tenantObjects map[string][]metadata.Object, rid string) error {

	actionGroups := GenerateActionGroups(tenantObjects)

	if len(registeredInfo.ActionGroups) == 0 {
		if err := i.Client.RegisterActionGroups(ctx, h, actionGroups); err != nil {
			blog.Errorf("register action groups(%s) failed, err: %s, rid: %s", actionGroups, err, rid)
			return err
		}
		return nil
	}

	if reflect.DeepEqual(registeredInfo.ActionGroups, actionGroups) {
		return nil
	}

	if err := i.Client.UpdateActionGroups(ctx, h, actionGroups); err != nil {
		blog.Errorf("update action groups(%s) failed, err: %s, rid: %s", actionGroups, err, rid)
		return err
	}
	return nil
}

// registerResCreatorActions register or update resource creator actions
func (i IAM) registerResCreatorActions(ctx context.Context, h http.Header, registeredInfo *iam.RegisteredSystemInfo,
	rid string) error {
	rcActions := GenerateResourceCreatorActions()

	if len(registeredInfo.ResourceCreatorActions.Config) == 0 {
		if err := i.Client.RegisterResourceCreatorActions(ctx, h, rcActions); err != nil {
			blog.Errorf("register resource creator actions(%s) failed, err: %s, rid: %s", rcActions, err, rid)
			return err
		}
		return nil
	}

	if reflect.DeepEqual(registeredInfo.ResourceCreatorActions, rcActions) {
		return nil
	}

	if err := i.Client.UpdateResourceCreatorActions(ctx, h, rcActions); err != nil {
		blog.Errorf("update resource creator actions(%s) failed, err: %s, rid: %s", rcActions, err, rid)
		return err
	}
	return nil
}

// registerCommonActions register or update common actions
func (i IAM) registerCommonActions(ctx context.Context, h http.Header, registeredInfo *iam.RegisteredSystemInfo,
	rid string) error {
	commonActions := GenerateCommonActions()

	if len(registeredInfo.CommonActions) == 0 {
		if err := i.Client.RegisterCommonActions(ctx, h, commonActions); err != nil {
			blog.Errorf("register common actions(%s) failed, err: %s, rid: %s", commonActions, err, rid)
			return err
		}
		return nil
	}

	if reflect.DeepEqual(registeredInfo.CommonActions, commonActions) {
		return nil
	}

	if err := i.Client.UpdateCommonActions(ctx, h, commonActions); err != nil {
		blog.Errorf("update common actions(%s) failed, err: %s, rid: %s", commonActions, err, rid)
		return err
	}
	return nil
}

// SyncIAMSysInstances sync system instances between CMDB and IAM
// it check the difference of system instances resource between CMDB and IAM
// if they have difference, sync and make them same
func (i IAM) SyncIAMSysInstances(ctx context.Context, header http.Header, redisCli redis.Client,
	tenantObjects map[string][]metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	// validate the objects
	for tenantID, objects := range tenantObjects {
		for _, object := range objects {
			if object.ID == 0 || len(object.ObjectID) == 0 || len(object.ObjectName) == 0 {
				blog.Errorf("sync iam system instances but object(%#v) is invalid, tenantID: %s, rid: %s",
					tenantID, object, rid)
				return errors.New("sync iam instances, but object is invalid")
			}
		}
	}

	locker, err := tryLockRegister(redisCli, rid)
	if err != nil {
		return err
	}
	defer locker.Unlock()
	fields := []iamtypes.SystemQueryField{iamtypes.FieldResourceTypes, iamtypes.FieldActions,
		iamtypes.FieldActionGroups, iamtypes.FieldInstanceSelections}
	iamInfo, err := i.Client.GetSystemInfo(ctx, header, fields)
	if err != nil {
		blog.Errorf("sync iam sysInstances failed, get system info error: %s, fields: %s, rid: %s", err, fields, rid)
		return err
	}

	// get the cmdb resources
	cmdbActions := genDynamicActions(tenantObjects)
	cmdbInstanceSelections := genDynamicInstanceSelections(tenantObjects)
	cmdbResourceTypes := genDynamicResourceTypes(tenantObjects)

	// compare resources between cmdb and iam
	addedActions, deletedActions := compareActions(cmdbActions, iamInfo.Actions)
	addedInstanceSelections, deletedInstanceSelections := compareInstanceSelections(cmdbInstanceSelections,
		iamInfo.InstanceSelections)
	addedResourceTypes, deletedResourceTypes := compareResourceTypes(cmdbResourceTypes, iamInfo.ResourceTypes)
	// 因为资源间的依赖关系，删除和更新的顺序为 1.Action 2.InstanceSelection 3.IamResourceType
	// 因为资源间的依赖关系，新建的顺序则反过来为 1.IamResourceType 2.InstanceSelection 3.Action
	// ActionGroup依赖于Action，该资源的增删操作始终放在最后
	// 先删除资源，再新增资源，因为实例视图的名称在系统中是唯一的，如果不先删，同样名称的实例视图将创建失败

	err = i.deleteIamResources(ctx, header, deletedActions, deletedInstanceSelections, deletedResourceTypes, rid)
	if err != nil {
		return err
	}

	err = i.registerResource(ctx, header, addedResourceTypes, addedInstanceSelections, addedActions)
	if err != nil {
		return err
	}

	// update action_groups in iam, the action groups contains only the existed actions in iam
	if len(addedActions) > 0 || len(deletedActions) > 0 {
		actionMap := map[iamtypes.ActionID]struct{}{}
		for _, action := range iamInfo.Actions {
			if !isIAMSysInstanceAction(action.ID) {
				actionMap[action.ID] = struct{}{}
			}
		}
		for _, action := range cmdbActions {
			actionMap[action.ID] = struct{}{}
		}
		cmdbActionGroups := GenerateActionGroups(tenantObjects)
		actualActionGroups := getActionGroupWithExistAction(cmdbActionGroups, actionMap)

		// if all exist actions in iam needs no action group(which happens when first initializing), **skip**
		if len(actualActionGroups) > 0 {
			blog.Infof("begin update actionGroups, rid: %s", rid)
			if err = i.Client.UpdateActionGroups(ctx, header, actualActionGroups); err != nil {
				blog.Errorf("update action groups failed, actionGroups: %v,  err: %v, rid: %s", actualActionGroups, err,
					rid)
				return err
			}
		}
	}

	return nil
}

func (i IAM) registerResource(ctx context.Context, header http.Header, addedResourceTypes []iam.ResourceType,
	addedInstanceSelections []iam.InstanceSelection, addedActions []iam.ResourceAction) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	// add cmdb ResourceTypes in iam
	if len(addedResourceTypes) > 0 {
		blog.Infof("begin add resourceTypes, count:%d, detail:%v, rid: %s", len(addedResourceTypes), addedResourceTypes,
			rid)
		if err := i.Client.RegisterResourcesTypes(ctx, header, addedResourceTypes); err != nil {
			blog.Errorf("sync iam sysInstances failed, add resourceType error: %s, resourceType: %s, rid: %s",
				err, addedResourceTypes, rid)
			return err
		}
	}

	// add cmdb InstanceSelections in iam
	if len(addedInstanceSelections) > 0 {
		blog.Infof("begin add instanceSelections, count:%d, detail:%v, rid: %s",
			len(addedInstanceSelections), addedInstanceSelections, rid)
		if err := i.Client.RegisterInstanceSelections(ctx, header, addedInstanceSelections); err != nil {
			blog.Errorf("sync iam sysInstances failed, add instanceSelections error: %s, instanceSelections: %s, "+
				"rid: %s", err, addedInstanceSelections, rid)
			return err
		}
	}

	// add cmdb actions in iam
	if len(addedActions) > 0 {
		blog.Infof("begin add actions, count:%d, detail:%v, rid: %s", len(addedActions), addedActions, rid)
		if err := i.Client.RegisterActions(ctx, header, addedActions); err != nil {
			blog.ErrorJSON("sync iam sysInstances failed, add IAM actions failed, error: %s, actions: %s, rid: %s",
				err, addedActions, rid)
			return err
		}
	}
	return nil
}

func (i IAM) registerIamResource(kit *rest.Kit, addedResourceTypes []iam.ResourceType,
	addedInstanceSelections []iam.InstanceSelection, addedActions []iam.ResourceAction) error {

	// add cmdb ResourceTypes in iam
	if len(addedResourceTypes) > 0 {
		blog.Infof("begin add resourceTypes, count:%d, detail:%v, rid: %s", len(addedResourceTypes), addedResourceTypes,
			kit.Rid)
		if err := i.Client.RegisterResourcesTypes(kit.Ctx, kit.Header, addedResourceTypes); err != nil {
			blog.Errorf("sync iam sysInstances failed, add resourceType, resourceType: %s, error: %v, rid: %s",
				addedResourceTypes, err, kit.Rid)
			return err
		}
	}

	// add cmdb InstanceSelections in iam
	if len(addedInstanceSelections) > 0 {
		blog.Infof("begin add instanceSelections, count:%d, detail:%v, rid: %s",
			len(addedInstanceSelections), addedInstanceSelections, kit.Rid)
		if err := i.Client.RegisterInstanceSelections(kit.Ctx, kit.Header, addedInstanceSelections); err != nil {
			blog.Errorf("sync iam sysInstances failed, add instanceSelections error: %s, instanceSelections: %s, "+
				"rid: %s", err, addedInstanceSelections, kit.Rid)
			return err
		}
	}

	// add cmdb actions in iam
	if len(addedActions) > 0 {
		blog.Infof("begin add actions, count: %d, detail: %v, rid: %s", len(addedActions),
			addedActions, kit.Rid)
		if err := i.Client.RegisterActions(kit.Ctx, kit.Header, addedActions); err != nil {
			blog.Errorf("sync iam sysInstances failed, add IAM actions failed, actions: %s, error: %v, rid: %s",
				addedActions, err, kit.Rid)
			return err
		}
	}
	return nil
}

func (i IAM) deleteIamResources(ctx context.Context, header http.Header, deletedActions []iamtypes.ActionID,
	deletedInstanceSelections []iamtypes.InstanceSelectionID, deletedResourceTypes []iamtypes.TypeID,
	rid string) error {

	// delete unnecessary actions in iam
	if len(deletedActions) > 0 {
		blog.Infof("begin delete actions, count:%d, detail:%v, rid: %s", len(deletedActions), deletedActions, rid)

		// before deleting action, the dependent action policies must be deleted
		for _, actionID := range deletedActions {
			if err := i.Client.DeleteActionPolicies(ctx, header, actionID); err != nil {
				blog.Errorf("delete iam action %s policies failed, err: %v, rid: %s", actionID, err, rid)
				return err
			}
		}

		if err := i.Client.DeleteActions(ctx, header, deletedActions); err != nil {
			blog.Errorf("delete IAM actions failed, err: %s, actions: %s, rid: %s", err, deletedActions, rid)
			return err
		}
	}

	// delete unnecessary InstanceSelections in iam
	if len(deletedInstanceSelections) > 0 {
		blog.Infof("begin delete instanceSelections, count: %d, detail: %v, rid: %s", len(deletedInstanceSelections),
			deletedInstanceSelections, rid)
		if err := i.Client.DeleteInstanceSelections(ctx, header, deletedInstanceSelections); err != nil {
			blog.Errorf("delete instanceSelections failed, err: %s, instanceSelections: %s, rid: %s", err,
				deletedInstanceSelections, rid)
			return err
		}
	}

	// delete unnecessary ResourceTypes in iam
	if len(deletedResourceTypes) > 0 {
		blog.Infof("begin delete resourceTypes, count: %d, detail: %v, rid: %s", len(deletedResourceTypes),
			deletedResourceTypes, rid)
		if err := i.Client.DeleteResourcesTypes(ctx, header, deletedResourceTypes); err != nil {
			blog.Errorf("delete resourceType failed, err: %s, resourceType: %s, rid: %s", err, deletedResourceTypes,
				rid)
			return err
		}
	}
	return nil
}

// getActionGroupWithExistAction get action groups that has actions that exists in iam
func getActionGroupWithExistAction(cmdbActionGroups []iam.ActionGroup,
	actionMap map[iamtypes.ActionID]struct{}) []iam.ActionGroup {

	actualActionGroups := make([]iam.ActionGroup, 0)

	for _, actionGroup := range cmdbActionGroups {
		actualActionGroup := iam.ActionGroup{
			Name:      actionGroup.Name,
			NameEn:    actionGroup.NameEn,
			SubGroups: make([]iam.ActionGroup, 0),
			Actions:   make([]iam.ActionWithID, 0),
		}

		for _, action := range actionGroup.Actions {
			if _, exists := actionMap[action.ID]; exists {
				actualActionGroup.Actions = append(actualActionGroup.Actions, action)
			}
		}

		actualActionGroup.SubGroups = getActionGroupWithExistAction(actionGroup.SubGroups, actionMap)

		if len(actualActionGroup.SubGroups) > 0 || len(actualActionGroup.Actions) > 0 {
			actualActionGroups = append(actualActionGroups, actualActionGroup)
		}
	}

	return actualActionGroups
}

// DeleteCMDBResource delete unnecessary CMDB resource from IAM
// it will  delete the resource if it exists on IAM
func (i IAM) DeleteCMDBResource(ctx context.Context, param *iamtypes.DeleteCMDBResourceParam,
	tenantObjects map[string][]metadata.Object) error {

	rid := util.ExtractRequestIDFromContext(ctx)
	header := headerutil.GenDefaultHeader()
	httpheader.SetRid(header, rid)

	fields := []iamtypes.SystemQueryField{iamtypes.FieldResourceTypes, iamtypes.FieldActions,
		iamtypes.FieldActionGroups, iamtypes.FieldInstanceSelections}
	iamInfo, err := i.Client.GetSystemInfo(ctx, header, fields)
	if err != nil {
		blog.Errorf("sync iam sysInstances failed, get system info error: %s, fields: %s, rid: %s",
			err, fields, rid)
		return err
	}
	// get deleted actions
	deletedActions := getDeletedActions(param.ActionIDs, iamInfo.Actions)
	deletedInstanceSelections := getDeletedInstanceSelections(param.InstanceSelectionIDs,
		iamInfo.InstanceSelections)
	deletedResourceTypes := getDeletedResourceTypes(param.TypeIDs, iamInfo.ResourceTypes)

	// 因为资源间的依赖关系，删除的顺序为 1.Action 2.InstanceSelection 3.IamResourceType
	// ActionGroup依赖于Action，该资源的增删操作始终放在最后
	// delete unnecessary actions in iam
	if len(deletedActions) > 0 {
		// before deleting action, the dependent action policies must be deleted
		for _, actionID := range deletedActions {
			if err = i.Client.DeleteActionPolicies(ctx, header, actionID); err != nil {
				blog.Errorf("delete cmdb resource failed, delete action %s policies err: %s, rid: %s",
					actionID, err, rid)
				return err
			}
		}
		blog.Infof("begin delete actions, count:%d, detail:%v, rid: %s", len(deletedActions), deletedActions, rid)
		if err := i.Client.DeleteActions(ctx, header, deletedActions); err != nil {
			blog.Errorf("delete cmdb resource failed, delete IAM actions error: %s, actions: %s, rid: %s",
				err, deletedActions, rid)
			return err
		}
	}
	// delete unnecessary InstanceSelections in iam
	if len(deletedInstanceSelections) > 0 {
		blog.Infof("begin delete instanceSelections, count:%d, detail:%v, rid: %s",
			len(deletedInstanceSelections), deletedInstanceSelections, rid)
		if err := i.Client.DeleteInstanceSelections(ctx, header, deletedInstanceSelections); err != nil {
			blog.Errorf("delete cmdb resource failed, delete instanceSelections error: %s, instanceSelections: %s,"+
				"rid: %s", err, deletedInstanceSelections, rid)
			return err
		}
	}
	// delete unnecessary ResourceTypes in iam
	if len(deletedResourceTypes) > 0 {
		blog.Infof("begin delete resourceTypes, count:%d, detail:%v, rid: %s",
			len(deletedResourceTypes), deletedResourceTypes, rid)
		if err := i.Client.DeleteResourcesTypes(ctx, header, deletedResourceTypes); err != nil {
			blog.Errorf("delete cmdb resource failed, delete resourceType error: %s, resourceType: %s, "+
				"rid: %s", err, deletedResourceTypes, rid)
			return err
		}
	}
	// update action_groups in iam
	if len(deletedActions) > 0 {
		actionMap := map[iamtypes.ActionID]struct{}{}
		for _, action := range iamInfo.Actions {
			actionMap[action.ID] = struct{}{}
		}
		for _, action := range deletedActions {
			delete(actionMap, action)
		}
		cmdbActionGroups := GenerateActionGroups(tenantObjects)
		actualActionGroups := getActionGroupWithExistAction(cmdbActionGroups, actionMap)
		if len(actualActionGroups) > 0 {
			blog.Infof("begin update action groups")
			if err := i.Client.UpdateActionGroups(ctx, header, actualActionGroups); err != nil {
				blog.Errorf("update action groups(%+v) after delete cmdb resource from iam failed, err: %v, rid: %s",
					actualActionGroups, err, rid)
				return err
			}
		}
	}
	return nil
}

// RegisterToIAM register to iam
func (i IAM) RegisterToIAM(ctx context.Context, header http.Header, host string) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	_, err := i.Client.GetSystemInfo(ctx, header, []iamtypes.SystemQueryField{})
	if err == nil {
		return nil
	}

	if err != iam.ErrNotFound {
		blog.Errorf("get system info failed, error: %v, rid: %s", err, rid)
		return err
	}

	// if iam cmdb system has not been registered, register system
	sys := iam.System{
		ID:          iamtypes.SystemIDCMDB,
		Name:        iamtypes.SystemNameCMDB,
		EnglishName: iamtypes.SystemNameCMDBEn,
		Clients:     iamtypes.SystemIDCMDB,
		ProviderConfig: &iam.SysConfig{
			Host: host,
			Auth: "basic",
		},
	}
	if err = i.Client.RegisterSystem(ctx, header, sys); err != nil {
		blog.Errorf("register system %s failed, error: %v, rid: %s", sys, err, rid)
		return err
	}
	blog.V(5).Infof("register new system %+v succeed", sys)
	return nil
}

// IsRegisteredToIAM checks if cmdb is registered to iam or not
func (i IAM) IsRegisteredToIAM(ctx context.Context) (bool, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	header := headerutil.GenDefaultHeader()
	_, err := i.Client.GetSystemInfo(ctx, header, []iamtypes.SystemQueryField{})

	if err == nil {
		return true, nil
	}

	if err != iam.ErrNotFound {
		blog.Errorf("get system info failed, error: %v, rid: %s", err, rid)
		return false, err
	}

	return false, nil

}

// getDeletedActions get deleted actions
func getDeletedActions(cmdbActionIDs []iamtypes.ActionID, iamActions []iam.ResourceAction) []iamtypes.ActionID {
	deletedActions := make([]iamtypes.ActionID, 0)
	iamActionMap := map[iamtypes.ActionID]struct{}{}
	for _, action := range iamActions {
		iamActionMap[action.ID] = struct{}{}
	}
	for _, actionID := range cmdbActionIDs {
		if _, ok := iamActionMap[actionID]; ok {
			deletedActions = append(deletedActions, actionID)
		}
	}

	return deletedActions
}

// getDeletedInstanceSelections get deleted instance selections
func getDeletedInstanceSelections(cmdbInstanceSelectionIDs []iamtypes.InstanceSelectionID,
	iamInstanceSelections []iam.InstanceSelection) []iamtypes.InstanceSelectionID {
	deletedInstanceSelections := make([]iamtypes.InstanceSelectionID, 0)
	iamInstanceSelectionMap := map[iamtypes.InstanceSelectionID]struct{}{}
	for _, instanceSelection := range iamInstanceSelections {
		iamInstanceSelectionMap[instanceSelection.ID] = struct{}{}
	}
	for _, instanceSelectionID := range cmdbInstanceSelectionIDs {
		if _, ok := iamInstanceSelectionMap[instanceSelectionID]; ok {
			deletedInstanceSelections = append(deletedInstanceSelections, instanceSelectionID)
		}
	}

	return deletedInstanceSelections
}

// getDeletedResourceTypes get deleted resource types
func getDeletedResourceTypes(cmdbTypeIDs []iamtypes.TypeID, iamResourceTypes []iam.ResourceType) []iamtypes.TypeID {
	deletedResourceTypes := make([]iamtypes.TypeID, 0)
	iamResourceTypeMap := map[iamtypes.TypeID]struct{}{}
	for _, resourceType := range iamResourceTypes {
		iamResourceTypeMap[resourceType.ID] = struct{}{}
	}
	for _, typeID := range cmdbTypeIDs {
		if _, ok := iamResourceTypeMap[typeID]; ok {
			deletedResourceTypes = append(deletedResourceTypes, typeID)
		}
	}

	return deletedResourceTypes
}

// compareActions compare actions between cmdb and iam
func compareActions(cmdbActions []iam.ResourceAction, iamActions []iam.ResourceAction) (
	addedActions []iam.ResourceAction, deletedActionIDs []iamtypes.ActionID) {
	iamActionMap := map[iamtypes.ActionID]struct{}{}

	for _, action := range iamActions {
		if isIAMSysInstanceAction(action.ID) {
			iamActionMap[action.ID] = struct{}{}
		}
	}

	for _, action := range cmdbActions {
		if _, ok := iamActionMap[action.ID]; !ok {
			addedActions = append(addedActions, action)
		} else {
			delete(iamActionMap, action.ID)
		}
	}

	for actionID := range iamActionMap {
		deletedActionIDs = append(deletedActionIDs, actionID)
	}

	return addedActions, deletedActionIDs
}

// compareInstanceSelections compare instanceSelections between cmdb and iam
func compareInstanceSelections(cmdbInstanceSelections []iam.InstanceSelection,
	iamInstanceSelections []iam.InstanceSelection) (addInstanceSelection []iam.InstanceSelection,
	deletedInstanceSelectionIDs []iamtypes.InstanceSelectionID) {
	iamInstanceSelectionMap := map[iamtypes.InstanceSelectionID]struct{}{}

	for _, instanceSelection := range iamInstanceSelections {
		if isIAMSysInstanceSelection(instanceSelection.ID) {
			iamInstanceSelectionMap[instanceSelection.ID] = struct{}{}
		}
	}

	for _, instanceSelection := range cmdbInstanceSelections {
		if _, ok := iamInstanceSelectionMap[instanceSelection.ID]; !ok {
			addInstanceSelection = append(addInstanceSelection, instanceSelection)
		} else {
			delete(iamInstanceSelectionMap, instanceSelection.ID)
		}
	}

	for instanceSelectionID := range iamInstanceSelectionMap {
		deletedInstanceSelectionIDs = append(deletedInstanceSelectionIDs, instanceSelectionID)
	}

	return addInstanceSelection, deletedInstanceSelectionIDs
}

// compareResourceTypes compare resourceTypes between cmdb and iam
func compareResourceTypes(cmdbResourceTypes []iam.ResourceType, iamResourceTypes []iam.ResourceType) (
	addedResourceTypes []iam.ResourceType, deletedTypeIDs []iamtypes.TypeID) {
	iamResourceTypeMap := map[iamtypes.TypeID]struct{}{}

	for _, resourceType := range iamResourceTypes {
		if IsIAMSysInstance(resourceType.ID) {
			iamResourceTypeMap[resourceType.ID] = struct{}{}
		}
	}

	for _, resourceType := range cmdbResourceTypes {
		if _, ok := iamResourceTypeMap[resourceType.ID]; !ok {
			addedResourceTypes = append(addedResourceTypes, resourceType)
		} else {
			delete(iamResourceTypeMap, resourceType.ID)
		}
	}

	for typeID := range iamResourceTypeMap {
		deletedTypeIDs = append(deletedTypeIDs, typeID)
	}

	return addedResourceTypes, deletedTypeIDs
}

type authorizer struct {
	authClientSet authserver.AuthServerClientInterface
}

// NewAuthorizer new authorizer
func NewAuthorizer(clientSet apimachinery.ClientSetInterface) *authorizer {
	return &authorizer{authClientSet: clientSet.AuthServer()}
}

// AuthorizeBatch batch authorization will not pass if one of them does not have permission
func (a *authorizer) AuthorizeBatch(ctx context.Context, h http.Header, user meta.UserInfo,
	resources ...meta.ResourceAttribute) ([]types.Decision, error) {
	return a.authorizeBatch(ctx, h, true, user, resources...)
}

// AuthorizeAnyBatch batch authorization will pass if one of them has permission
func (a *authorizer) AuthorizeAnyBatch(ctx context.Context, h http.Header, user meta.UserInfo,
	resources ...meta.ResourceAttribute) ([]types.Decision, error) {
	return a.authorizeBatch(ctx, h, false, user, resources...)
}

func (a *authorizer) authorizeBatch(ctx context.Context, h http.Header, exact bool, user meta.UserInfo,
	resources ...meta.ResourceAttribute) ([]types.Decision, error) {

	rid := httpheader.GetRid(h)

	opts, decisions, err := parseAttributesToBatchOptions(rid, user, resources...)
	if err != nil {
		return nil, err
	}

	// all resources are skipped
	if opts == nil {
		return decisions, nil
	}

	if blog.V(5) {
		blog.InfoJSON("auth options: %s, rid: %s", opts, rid)
	}

	var authDecisions []types.Decision
	if exact {
		authDecisions, err = a.authClientSet.AuthorizeBatch(ctx, h, opts)
		if err != nil {
			blog.Errorf("authorize batch failed, err: %s, ops: %s, resources: %s, rid: %s", err, opts, resources, rid)
			return nil, err
		}
	} else {
		authDecisions, err = a.authClientSet.AuthorizeAnyBatch(ctx, h, opts)
		if err != nil {
			blog.Errorf("authorize any batch failed, err: %s, ops: %s, resources: %s, rid: %s", err, opts, resources,
				rid)
			return nil, err
		}

	}

	index := 0
	for _, decision := range authDecisions {
		// skip resources' decisions are already set as authorized
		for decisions[index].Authorized {
			index++
		}
		decisions[index].Authorized = decision.Authorized
		index++
	}

	return decisions, nil
}

func parseAttributesToBatchOptions(rid string, user meta.UserInfo,
	resources ...meta.ResourceAttribute) (*iam.AuthBatchOptions, []types.Decision, error) {
	if !auth.EnableAuthorize() {
		decisions := make([]types.Decision, len(resources))
		for i := range decisions {
			decisions[i].Authorized = true
		}
		return nil, decisions, nil
	}

	authBatchArr := make([]*iam.AuthBatch, 0)
	decisions := make([]types.Decision, len(resources))
	for index, resource := range resources {

		// this resource should be skipped, do not need to verify in auth center.
		if resource.Action == meta.SkipAction {
			decisions[index].Authorized = true
			blog.V(5).Infof("skip authorization for resource: %+v, rid: %s", resource, rid)
			continue
		}

		action, resources, err := AdaptAuthOptions(&resource)
		if err != nil {
			blog.Errorf("adaptor cmdb resource to iam failed, err: %s, rid: %s", err, rid)
			return nil, nil, err
		}

		// this resource should be skipped, do not need to verify in auth center.
		if action == iamtypes.Skip {
			decisions[index].Authorized = true
			blog.V(5).Infof("skip authorization for resource: %+v, rid: %s", resource, rid)
			continue
		}

		authBatchArr = append(authBatchArr, &iam.AuthBatch{
			Action:    iam.Action{ID: string(action)},
			Resources: resources,
		})
	}

	// all resources are skipped
	if len(authBatchArr) == 0 {
		return nil, decisions, nil
	}

	ops := &iam.AuthBatchOptions{
		System: iamtypes.SystemIDCMDB,
		Subject: iam.Subject{
			Type: "user",
			ID:   user.UserName,
		},
		Batch: authBatchArr,
	}
	return ops, decisions, nil
}

// ListAuthorizedResources 获取用户有的资源id权限列表
func (a *authorizer) ListAuthorizedResources(ctx context.Context, h http.Header,
	input meta.ListAuthorizedResourcesParam) (*types.AuthorizeList, error) {
	return a.authClientSet.ListAuthorizedResources(ctx, h, input)
}

// GetNoAuthSkipUrl get no auth skip url
func (a *authorizer) GetNoAuthSkipUrl(ctx context.Context, h http.Header,
	input *metadata.IamPermission) (string, error) {
	return a.authClientSet.GetNoAuthSkipUrl(ctx, h, input)
}

// GetPermissionToApply get permission to apply
func (a *authorizer) GetPermissionToApply(ctx context.Context, h http.Header,
	input []meta.ResourceAttribute) (*metadata.IamPermission, error) {
	return a.authClientSet.GetPermissionToApply(ctx, h, input)
}

// RegisterResourceCreatorAction register resourceCreator Action
func (a *authorizer) RegisterResourceCreatorAction(ctx context.Context, h http.Header,
	input metadata.IamInstanceWithCreator) (
	[]metadata.IamCreatorActionPolicy, error) {

	return a.authClientSet.RegisterResourceCreatorAction(ctx, h, input)
}

// BatchRegisterResourceCreatorAction batch register resourceCreator action
func (a *authorizer) BatchRegisterResourceCreatorAction(ctx context.Context, h http.Header,
	input metadata.IamInstancesWithCreator) (
	[]metadata.IamCreatorActionPolicy, error) {

	return a.authClientSet.BatchRegisterResourceCreatorAction(ctx, h, input)
}
