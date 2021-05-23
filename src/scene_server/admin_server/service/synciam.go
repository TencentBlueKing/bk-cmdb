/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"net/http"
	"time"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
)

const (
	// 同步周期最小值
	SyncIAMPeriodMin = 5 * time.Minute
	// 同步周期默认值
	SyncIAMPeriodDefault = 30 * time.Minute
)

// 同步周期
var SyncIAMPeriod time.Duration

// newHeader 创建IAM同步需要的header
func newHeader() http.Header {
	header := make(http.Header)
	header.Add(common.BKHTTPOwnerID, common.BKSuperOwnerID)
	header.Add(common.BKHTTPHeaderUser, common.BKIAMSyncUser)
	header.Add(common.BKHTTPLanguage, "cn")
	header.Add(common.BKHTTPCCRequestID, util.GenerateRID())
	header.Add("Content-Type", "application/json")
	return header
}

// newKit 创建新的Kit
func newKit() *rest.Kit {
	header := newHeader()
	if header.Get(common.BKHTTPCCRequestID) == "" {
		header.Set(common.BKHTTPCCRequestID, util.GenerateRID())
	}
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)
	user := util.GetUser(header)
	supplierAccount := util.GetOwnerID(header)
	defaultCCError := util.GetDefaultCCError(header)

	return &rest.Kit{
		Rid:             rid,
		Header:          header,
		Ctx:             ctx,
		CCError:         defaultCCError,
		User:            user,
		SupplierAccount: supplierAccount,
	}
}

// SyncIAM sync the sys instances resource between CMDB and IAM
func (s *Service) SyncIAM() {
	for {
		// new kit with a different rid, header, ctx
		kit := newKit()
		err := s.SyncIAMSysInstances(kit)
		if err != nil {
			blog.Errorf("sync action with IAM failed, err:%v", err)
		}
		time.Sleep(SyncIAMPeriod)
	}
}

// SyncIAMSysInstances check the difference of system instances resource between IAM and CMDB
// In most cases, this func will unregister IAM system instances resources which are discard.
func (s *Service) SyncIAMSysInstances(kit *rest.Kit) error {
	fields := []iam.SystemQueryField{iam.FieldResourceTypes,
		iam.FieldActions, iam.FieldActionGroups, iam.FieldInstanceSelections}
	iamResp, err := s.iam.Client.GetSystemInfo(kit.Ctx, fields)
	if err != nil {
		blog.ErrorJSON("syc iam sysInstances failed, get system info error: %s, fields: %s, rid: %s",
			err.Error(), fields, kit.Rid)
		return err
	}

	// get all custom models in cmdb
	models, err := s.GetCustomObjects(kit.Header)
	if err != nil {
		blog.Errorf("syc iam sysInstances failed, get custom models err: %s, rid:%s", err.Error(), kit.Rid)
		return err
	}

	// get the cmdb resources
	cmdbActionGroups := iam.GenModelInstanceManageActionGroups(models)
	cmdbActions := iam.GenModelInstanceActions(models)
	cmdbInstanceSelections := iam.GenDynamicInstanceSelections(models)
	cmdbResourceTypes := iam.GenDynamicResourceTypes(models)

	// compare resources between cmdb and iam
	addActions, deleteActions := compareActions(cmdbActions, iamResp.Data.Actions)
	addInstanceSelections, deleteInstanceSelections := compareInstanceSelections(cmdbInstanceSelections,
		iamResp.Data.InstanceSelections)
	addResourceTypes, deleteResourceTypes := compareResourceTypes(cmdbResourceTypes, iamResp.Data.ResourceTypes)

	// 因为资源间的依赖关系，删除和更新的顺序为 1.Action 2.InstanceSelection 3.ResourceType
	// 因为资源间的依赖关系，新建的顺序则反过来为 1.ResourceType 2.InstanceSelection 3.Action
	// ActionGroup依赖于Action，增删操作始终放在最后
	// 先删除资源，再新增资源，因为实例视图的名称在系统中是唯一的，如果不先删，同样名称的实例视图将创建失败

	// delete unnecessary actions in iam
	if len(deleteActions) > 0 {
		if err := s.iam.Client.DeleteActions(kit.Ctx, deleteActions); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, delete IAM actions failed, error: %s, "+
				"actions: %s, rid:%s", err.Error(), deleteActions, kit.Rid)
			return err
		}
	}

	// delete unnecessary InstanceSelections in iam
	if len(deleteInstanceSelections) > 0 {
		if err := s.iam.Client.DeleteInstanceSelections(kit.Ctx, deleteInstanceSelections); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, delete instanceSelections error: %s, "+
				"instanceSelections: %s, rid:%s", err.Error(), deleteInstanceSelections, kit.Rid)
			return err
		}
	}

	// delete unnecessary ResourceTypes in iam
	if len(deleteResourceTypes) > 0 {
		if err := s.iam.Client.DeleteResourcesTypes(kit.Ctx, deleteResourceTypes); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, delete resourceType error: %s, "+
				"resourceType: %s, rid:%s", err.Error(), deleteResourceTypes, kit.Rid)
			return err
		}
	}

	// add cmdb actions in iam
	if len(addActions) > 0 {
		if err := s.iam.Client.RegisterActions(kit.Ctx, addActions); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, add IAM actions failed, error: %s, "+
				"actions: %s, rid:%s", err.Error(), addActions, kit.Rid)
			return err
		}
	}

	// add cmdb InstanceSelections in iam
	if len(addInstanceSelections) > 0 {
		if err := s.iam.Client.RegisterInstanceSelections(kit.Ctx, addInstanceSelections); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, add instanceSelections error: %s, "+
				"instanceSelections: %s, rid:%s", err.Error(), addInstanceSelections, kit.Rid)
			return err
		}
	}

	// add cmdb ResourceTypes in iam
	if len(addResourceTypes) > 0 {
		if err := s.iam.Client.RegisterResourcesTypes(kit.Ctx, addResourceTypes); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, add resourceType error: %s, "+
				"resourceType: %s, rid:%s", err.Error(), addResourceTypes, kit.Rid)
			return err
		}
	}

	// update action_groups in iam
	if len(addActions) > 0 || len(deleteActions) > 0 {
		if err := s.iam.Client.UpdateActionGroups(kit.Ctx, cmdbActionGroups); err != nil {
			blog.ErrorJSON("syc iam sysInstances failed, update actionGroups error: %s, "+
				"actionGroups: %s, rid:%s", err.Error(), cmdbActionGroups, kit.Rid)
			return err
		}
	}

	return nil
}

func compareActions(cmdbActions []iam.ResourceAction, iamActions []iam.ResourceAction) (
	addActions []iam.ResourceAction, deleteActionIDs []iam.ActionID) {
	cmdbActionMap := map[iam.ActionID]struct{}{}
	iamActionMap := map[iam.ActionID]struct{}{}

	for _, act := range iamActions {
		iamActionMap[act.ID] = struct{}{}
	}

	for _, act := range cmdbActions {
		cmdbActionMap[act.ID] = struct{}{}
		if _, ok := iamActionMap[act.ID]; !ok {
			addActions = append(addActions, act)
		}
	}

	for _, act := range iamActions {
		iamActionMap[act.ID] = struct{}{}
		if _, ok := cmdbActionMap[act.ID]; !ok {
			deleteActionIDs = append(deleteActionIDs, act.ID)
		}
	}

	return addActions, deleteActionIDs
}

func compareInstanceSelections(cmdbInstanceSelections []iam.InstanceSelection,
	iamInstanceSelections []iam.InstanceSelection) (addInstanceSelection []iam.InstanceSelection,
	deleteInstanceSelectionIDs []iam.InstanceSelectionID) {
	cmdbInstanceSelectionMap := map[iam.InstanceSelectionID]struct{}{}
	iamInstanceSelectionMap := map[iam.InstanceSelectionID]struct{}{}

	for _, instanceSelection := range iamInstanceSelections {
		iamInstanceSelectionMap[instanceSelection.ID] = struct{}{}
	}

	for _, instanceSelection := range cmdbInstanceSelections {
		cmdbInstanceSelectionMap[instanceSelection.ID] = struct{}{}
		if _, ok := iamInstanceSelectionMap[instanceSelection.ID]; !ok {
			addInstanceSelection = append(addInstanceSelection, instanceSelection)
		}
	}

	for _, instanceSelection := range iamInstanceSelections {
		iamInstanceSelectionMap[instanceSelection.ID] = struct{}{}
		if _, ok := cmdbInstanceSelectionMap[instanceSelection.ID]; !ok {
			deleteInstanceSelectionIDs = append(deleteInstanceSelectionIDs, instanceSelection.ID)
		}
	}

	return addInstanceSelection, deleteInstanceSelectionIDs
}

func compareResourceTypes(cmdbResourceTypes []iam.ResourceType, iamResourceTypes []iam.ResourceType) (
	addResourceTypes []iam.ResourceType, deleteTypeIDs []iam.TypeID) {
	cmdbResourceTypeMap := map[iam.TypeID]struct{}{}
	iamResourceTypeMap := map[iam.TypeID]struct{}{}

	for _, act := range iamResourceTypes {
		iamResourceTypeMap[act.ID] = struct{}{}
	}

	for _, act := range cmdbResourceTypes {
		cmdbResourceTypeMap[act.ID] = struct{}{}
		if _, ok := iamResourceTypeMap[act.ID]; !ok {
			addResourceTypes = append(addResourceTypes, act)
		}
	}

	for _, act := range iamResourceTypes {
		iamResourceTypeMap[act.ID] = struct{}{}
		if _, ok := cmdbResourceTypeMap[act.ID]; !ok {
			deleteTypeIDs = append(deleteTypeIDs, act.ID)
		}
	}

	return addResourceTypes, deleteTypeIDs
}
