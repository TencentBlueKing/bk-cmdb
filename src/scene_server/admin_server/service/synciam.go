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
		// 每次同步生成新的kit
		kit := newKit()
		err := s.SyncIAMSysInstanceResources(kit)
		if err != nil {
			blog.Errorf("sync action with IAM failed, err:%v", err)
		}
		time.Sleep(s.Config.Iam.Interval)
	}
}

// SyncIAMSysInstanceResources check iam dynamic-model resources with CC models.
// In most cases, this func will unregister IAM model resources which are discard.
func (s *Service) SyncIAMSysInstanceResources(kit *rest.Kit) error {
	// Direct call IAM, get infos from iam.
	sysResp, err := s.iam.Client.GetSystemDynamicInfo(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("Synchronize actions with IAM failed, get resource actions from IAM failed, error: %s, resource actions: %s, rid: %s", err.Error(), sysResp, kit.Rid)
		return err
	}

	// 这是有顺序的1.ActionGroup 2.Action 3.InstanceSelection 4.ResourceType

	// 需要先拿到当前已存在的模型, 再与IAM返回结果进行对比
	models, err := s.GetCustomObjects(kit.Header)
	if err != nil {
		blog.Errorf("Synchronize actions with IAM failed, get custom models failed, err: %s, rid:%s", err.Error(),
			kit.Rid)
		return err
	}

	cmdbActionGroupList := iam.GenModelInstanceManageActionGroups(models)
	cmdbActionList := iam.GenModelInstanceActions(models)
	cmdbInstanceSelectionList := iam.GenDynamicInstanceSelectionWithModel(models)
	cmdbResourceTypeList := iam.GenDynamicResourceTypeWithModel(models)

	// 筛选需要删除的对象, 其中ActionGroup是全量同步, 不需要此步骤
	deleteActionList := checkActionList(cmdbActionList, sysResp.Data.Actions)
	deleteInstanceSelectionList := checkInstanceSelectionList(cmdbInstanceSelectionList, sysResp.Data.InstanceSelections)
	deleteResourceTypeList := checkResourceTypeList(cmdbResourceTypeList, sysResp.Data.ResourceTypes)

	// update action_groups in iam
	if err := s.iam.Client.UpdateActionGroups(kit.Ctx, cmdbActionGroupList); err != nil {
		blog.ErrorJSON("Synchronize actions with IAM failed, delete IAM actions failed, error: %s, resource actions: %s, rid:%s", err.Error(), deleteActionList, kit.Rid)
		return err
	}

	// delete unnecessary actions in iam
	if err := s.iam.Client.DeleteAction(kit.Ctx, deleteActionList); err != nil {
		blog.ErrorJSON("Synchronize actions with IAM failed, delete IAM actions failed, error: %s, resource actions: %s, rid:%s", err.Error(), deleteActionList, kit.Rid)
		return err
	}

	// delete unnecessary InstanceSelections in iam
	if err := s.iam.Client.DeleteInstanceSelection(kit.Ctx, deleteInstanceSelectionList); err != nil {
		blog.ErrorJSON("Synchronize actions with IAM failed, delete IAM actions failed, error: %s, resource actions: %s, rid:%s", err.Error(), deleteActionList, kit.Rid)
		return err
	}

	// delete unnecessary ResourceTypes in iam
	if err := s.iam.Client.DeleteResourcesTypes(kit.Ctx, deleteResourceTypeList); err != nil {
		blog.ErrorJSON("Synchronize actions with IAM failed, delete IAM actions failed, error: %s, resource actions: %s, rid:%s", err.Error(), deleteActionList, kit.Rid)
		return err
	}
	return nil
}

func checkActionList(cmdbActionList []iam.ResourceAction, iamActionList []iam.ResourceAction) []iam.ActionID {
	// 由整体的cmdbAction列表转换为cmdbAction集合
	cmdbActionMap := map[iam.ActionID]struct{}{}
	for _, act := range cmdbActionList {
		cmdbActionMap[act.ID] = struct{}{}
	}

	// 对比出IAM中多余的动作
	deleteActionList := []iam.ActionID{}
	for _, act := range iamActionList {
		if _, exists := cmdbActionMap[act.ID]; exists {
			continue
		}
		deleteActionList = append(deleteActionList, act.ID)
	}
	return deleteActionList
}

func checkInstanceSelectionList(cmdbInstanceSelectionList []iam.InstanceSelection, iamInstanceSelectionList []iam.InstanceSelection) []iam.InstanceSelectionID {
	cmdbInstanceSelectionMap := map[iam.InstanceSelectionID]struct{}{}
	for _, instSelection := range cmdbInstanceSelectionList {
		cmdbInstanceSelectionMap[instSelection.ID] = struct{}{}
	}

	deleteInstanceSelectionList := []iam.InstanceSelectionID{}
	for _, instSelection := range iamInstanceSelectionList {
		if _, exists := cmdbInstanceSelectionMap[instSelection.ID]; exists {
			continue
		}
		deleteInstanceSelectionList = append(deleteInstanceSelectionList, instSelection.ID)
	}
	return deleteInstanceSelectionList
}

func checkResourceTypeList(cmdbResourceTypeList []iam.ResourceType, iamResourceTypeList []iam.ResourceType) []iam.TypeID {
	cmdbResourceTypeMap := map[iam.TypeID]struct{}{}
	for _, resourceType := range cmdbResourceTypeList {
		cmdbResourceTypeMap[resourceType.ID] = struct{}{}
	}

	deleteResourceTypeList := []iam.TypeID{}
	for _, resourceType := range iamResourceTypeList {
		if _, exists := cmdbResourceTypeMap[resourceType.ID]; exists {
			continue
		}
		deleteResourceTypeList = append(deleteResourceTypeList, resourceType.ID)
	}
	return deleteResourceTypeList
}