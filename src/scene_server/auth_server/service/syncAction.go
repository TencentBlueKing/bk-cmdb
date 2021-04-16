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

package service

import (
	"context"
	"net/http"
	"time"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/app/options"
)

type SyncServer struct {
	Config  *options.Config
	Service *AuthService
}

// 根据配置文件, 每隔固定时间就同步下IAM的action列表, 将IAM中多余的action删除。
func (s *AuthService) LoopSyncActionWithIAM(ctx context.Context, config *options.Config) {
	kit := s.NewKit(config)
	timer := time.NewTimer(config.Auth.Interval)
	for true {
		select {
		// 计时器信号
		case <-timer.C:
			err := s.SyncIAMModelResources(*kit)
			if err != nil {
				blog.Errorf("sync action with IAM failed, err:%v", err)
			}
		// authServer退出信号
		case <-ctx.Done():
			blog.Infof("auth server will exit!")
			return
		}
	}
}

// NewHeader
func (s *AuthService) NewHeader(config *options.Config) http.Header {
	header := make(http.Header)
	header.Add(common.BKHTTPOwnerID, common.BKSuperOwnerID)
	header.Add(common.BKHTTPHeaderUser, common.BKAuthUser)
	header.Add(common.BKHTTPLanguage, "cn")
	header.Add(common.BKHTTPCCRequestID, util.GenerateRID())
	header.Add("Content-Type", "application/json")

	header.Add(iam.IamAppCodeHeader, config.Auth.AppCode)
	header.Add(iam.IamAppSecretHeader, config.Auth.AppSecret)
	return header
}

// NewKit
func (s *AuthService) NewKit(config *options.Config) *rest.Kit {
	header := s.NewHeader(config)

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

// SyncIAMModelResources check iam dynamic-model resources with CC models.
// In most cases, this func will unregister IAM model resources which are discard.
func (s *AuthService) SyncIAMModelResources(kit rest.Kit) error {

	// Direct call IAM, get infos from iam.
	sysResp, err := s.acIam.Client.GetSystemDynamicInfo(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("Synchronize actions with IAM failed, get resource actions from IAM failed, error: %s, resource actions: %s, rid: %s", err.Error(), sysResp, kit.Rid)
		return err
	}

	// 这是有顺序的1.ActionGroup 2.Action 3.InstanceSelection 4.ResourceType
	if staticActionGroupList == nil {
		staticActionGroupList = iam.GenerateStaticActionGroups()
	}
	if staticActionList == nil {
		staticActionList = iam.GenerateStaticActions()
	}
	if staticInstanceSelectionList == nil {
		staticInstanceSelectionList = iam.GenerateStaticInstanceSelections()
	}
	if staticResourceTypeList == nil {
		staticResourceTypeList = iam.GenerateStaticResourceTypes()
	}

	// 需要先拿到当前已存在的模型, 再与IAM返回结果进行对比
	models, err := s.lgc.CollectObjectsNotPre(&kit)
	if err != nil {
		blog.Errorf("Synchronize actions with IAM failed, collect notPre-models failed, err: %s, rid:%s", err.Error(), kit.Rid)
		return err
	}

	cmdbActionGroupList := append(staticActionGroupList, iam.GenModelInstanceManageActionGroups(models)...)
	cmdbActionList := append(staticActionList, iam.GenModelInstanceActions(models)...)
	cmdbInstanceSelectionList := append(staticInstanceSelectionList, iam.GenDynamicInstanceSelectionWithModel(models)...)
	cmdbResourceTypeList := append(staticResourceTypeList, iam.GenDynamicResourceTypeWithModel(models)...)

	// 筛选需要删除的对象, 其中ActionGroup是全量同步, 不需要此步骤
	deleteActionList := checkActionList(cmdbActionList, sysResp.Data.Actions)
	deleteInstanceSelectionList := checkInstanceSelectionList(cmdbInstanceSelectionList, sysResp.Data.InstanceSelections)
	deleteResourceTypeList := checkResourceTypeList(cmdbResourceTypeList, sysResp.Data.ResourceTypes)

	// 1.Direct call IAM, update action_groups in iam
	if err := s.acIam.Client.UpdateActionGroups(kit.Ctx, cmdbActionGroupList); err != nil {
		blog.ErrorJSON("Synchronize actions with IAM failed, delete IAM actions failed, error: %s, resource actions: %s, rid:%s", err.Error(), deleteActionList, kit.Rid)
		return err
	}
	// 2.Direct call IAM, delete certain actions in iam
	if err := s.acIam.Client.DeleteAction(kit.Ctx, deleteActionList); err != nil {
		blog.ErrorJSON("Synchronize actions with IAM failed, delete IAM actions failed, error: %s, resource actions: %s, rid:%s", err.Error(), deleteActionList, kit.Rid)
		return err
	}
	// 3.Direct call IAM, delete certain InstanceSelections in iam
	if err := s.acIam.Client.DeleteInstanceSelection(kit.Ctx, deleteInstanceSelectionList); err != nil {
		blog.ErrorJSON("Synchronize actions with IAM failed, delete IAM actions failed, error: %s, resource actions: %s, rid:%s", err.Error(), deleteActionList, kit.Rid)
		return err
	}
	// 4.Direct call IAM, delete certain ResourceTypes in iam
	if err := s.acIam.Client.DeleteResourcesTypes(kit.Ctx, deleteResourceTypeList); err != nil {
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
