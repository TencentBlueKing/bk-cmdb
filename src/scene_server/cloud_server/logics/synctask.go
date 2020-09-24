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

package logics

import (
	"fmt"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
)

func (lgc *Logics) SearchVpc(kit *rest.Kit, accountID int64, vpcOpt *metadata.SearchVpcOption) (*metadata.VpcHostCntResult, error) {
	accountConf, err := lgc.GetCloudAccountConf(kit, accountID)
	if err != nil {
		blog.Errorf("SearchVpc failed, rid:%s, accountID:%d, vpcOpt:%+v, err:%+v", kit.Rid, accountID, vpcOpt, err)
		return nil, kit.CCError.CCError(common.CCErrCloudVpcGetFail)
	}

	result, err := lgc.GetVpcHostCntInOneRegion(kit, *accountConf, vpcOpt.Region)
	if err != nil {
		blog.Errorf("SearchVpc failed, rid:%s, accountConf:%+v, vpcOpt:%+v, err:%+v", kit.Rid, accountConf, vpcOpt, err)
		return nil, kit.CCError.CCError(common.CCErrCloudVpcGetFail)
	}

	if len(result.Info) == 0 {
		return result, nil
	}

	vpcIDs := make([]string, 0)
	for _, info := range result.Info {
		vpcIDs = append(vpcIDs, info.VpcID)
	}

	vpcCloud, err := lgc.GetVpcCloudArea(kit, vpcIDs)
	if err != nil {
		blog.Errorf("SearchVpc failed, rid:%s, accountConf:%+v, vpcOpt:%+v, err:%+v", kit.Rid, accountConf, vpcOpt, err)
		return nil, kit.CCError.CCError(common.CCErrCloudVpcGetFail)
	}

	for i, info := range result.Info {
		if cloudID, ok := vpcCloud[info.VpcID]; ok {
			result.Info[i].CloudID = cloudID
		} else {
			result.Info[i].CloudID = -1
		}
	}

	return result, nil
}

func (lgc *Logics) GetVpcCloudArea(kit *rest.Kit, vpcIDs []string) (map[string]int64, error) {
	query := &metadata.QueryCondition{
		Fields: []string{common.BKCloudIDField, common.BKVpcID},
		Condition: mapstr.MapStr{common.BKVpcID: map[string]interface{}{
			common.BKDBIN: vpcIDs,
		}},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDPlat, query)
	if err != nil {
		blog.Errorf("GetVpcCloudIDs fail, rid:%s, err:%s, query:%+v", kit.Rid, err.Error(), *query)
		return nil, err
	}
	if !result.Result {
		blog.Errorf("GetVpcCloudIDs fail, rid:%s, err:%s, query:%+v", kit.Rid, result.ErrMsg, *query)
		return nil, fmt.Errorf("%s", result.ErrMsg)
	}

	ret := make(map[string]int64)
	for _, info := range result.Data.Info {
		cloudID, err := info.Int64(common.BKCloudIDField)
		if err != nil {
			blog.Errorf("GetVpcCloudIDs fail, rid:%s, err:%s, info:%+v", kit.Rid, err.Error(), result.Data.Info)
			return nil, err
		}
		vpcID, err := info.String(common.BKVpcID)
		if err != nil {
			blog.Errorf("GetVpcCloudIDs fail, rid:%s, err:%s, info:%+v", kit.Rid, err.Error(), result.Data.Info)
			return nil, err
		}
		ret[vpcID] = cloudID
	}

	return ret, nil
}

func (lgc *Logics) CreateSyncTask(kit *rest.Kit, task *metadata.CloudSyncTask) (*metadata.CloudSyncTask, error) {
	result, err := lgc.CoreAPI.CoreService().Cloud().CreateSyncTask(kit.Ctx, kit.Header, task)
	if err != nil {
		blog.Errorf("CreateSyncTask failed, rid:%s, task:%+v, err:%+v", kit.Rid, task, err)
		return nil, err
	}

	// generate audit log.
	audit := auditlog.NewSyncTaskAuditLog(lgc.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, auditErr := audit.GenerateAuditLog(generateAuditParameter, result.TaskID, result)
	if auditErr != nil {
		blog.Errorf("generate audit log failed after create sync task, taskID: %d, err: %v, rid: %s",
			result.TaskID, auditErr, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("save audit log failed after create sync task, taskID: %d, err: %v, rid: %s",
			result.TaskID, err, kit.Rid)
		return nil, err
	}

	return result, nil
}

func (lgc *Logics) SearchSyncTask(kit *rest.Kit, option *metadata.SearchSyncTaskOption) (*metadata.MultipleCloudSyncTask, error) {
	// set default limit
	if option.Page.Limit == 0 {
		option.Page.Limit = common.BKDefaultLimit
	}
	if option.Page.IsIllegal() {
		blog.Errorf("SearchSyncTask failed, Page is IsIllegal, rid:%s, page:%+v", kit.Rid, option.Page)
		return nil, kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded)
	}

	// set default sort
	if option.Page.Sort == "" {
		option.Page.Sort = "-" + common.CreateTimeField
	}

	// if fuzzy search, change the string query to regexp
	if option.IsFuzzy == true {
		for k, v := range option.Condition {
			field, ok := v.(string)
			if ok {
				option.Condition[k] = mapstr.MapStr{
					common.BKDBLIKE: params.SpecialCharChange(field),
					"$options":      "i",
				}
			}
		}
	}

	if auth.EnableAuthorize() {
		list, err := lgc.ListAuthorizedResources(kit, meta.CloudResourceTask, meta.Find)
		if err != nil {
			blog.Errorf("SearchSyncTask failed, rid:%s, option:%+v, ListAuthorizedResources err:%+v", kit.Rid, option, err)
			return nil, err
		}

		if option.Condition == nil {
			option.Condition = make(map[string]interface{})
		}

		option.Condition = map[string]interface{}{
			common.BKDBAND: []map[string]interface{}{
				option.Condition,
				{
					common.BKCloudSyncTaskID: map[string]interface{}{
						common.BKDBIN: list,
					},
				},
			},
		}
	}

	result, err := lgc.CoreAPI.CoreService().Cloud().SearchSyncTask(kit.Ctx, kit.Header, &option.SearchCloudOption)
	if err != nil {
		blog.Errorf("SearchSyncTask failed, rid:%s, option:%+v, err:%+v", kit.Rid, option, err)
		return nil, err
	}

	// 是否实时获取云厂商vpc下最新的主机数
	if option.LastestHostCount {
		if err := lgc.updateVpcHostCount(kit, result); err != nil {
			blog.Errorf("SearchSyncTask failed, rid:%s, option:%+v, err:%+v", kit.Rid, option, err)
			return nil, err
		}
	}

	return result, nil
}

// 更新vpc对应的主机数
func (lgc *Logics) updateVpcHostCount(kit *rest.Kit, multiTask *metadata.MultipleCloudSyncTask) error {
	// 待更新主机数的vpc对象，获取其地址
	vpcToUpdate := make(map[string]*metadata.VpcSyncInfo)
	accountOption := make(map[int64]*metadata.SearchVpcHostCntOption)

	for i := range multiTask.Info {
		task := multiTask.Info[i]
		for j := range task.SyncVpcs {
			syncVpc := &task.SyncVpcs[j]
			vpcToUpdate[syncVpc.VpcID] = syncVpc
			if accountOption[task.AccountID] == nil {
				accountOption[task.AccountID] = new(metadata.SearchVpcHostCntOption)
			}
			accountOption[task.AccountID].RegionVpcs = append(accountOption[task.AccountID].RegionVpcs, metadata.RegionVpc{
				Region: syncVpc.Region,
				VpcID:  syncVpc.VpcID,
			})
		}
	}

	// 获取所有vpc对应的主机数
	allVpcHostCnt := make(map[string]int64)
	for accountID, option := range accountOption {
		// 获取所有的vpc对应的主机数
		accountConf, err := lgc.GetCloudAccountConf(kit, accountID)
		if err != nil {
			blog.Errorf("updateVpcHostCount failed, rid:%s, accountID:%d, err:%+v", kit.Rid, accountID, err)
			return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		vpcHostCnt, err := lgc.GetVpcHostCnt(kit, *accountConf, *option)
		if err != nil {
			blog.Errorf("updateVpcHostCount failed, rid:%s, accountID:%d, err:%+v", kit.Rid, accountID, err)
			return kit.CCError.CCError(common.CCErrCloudVpcGetFail)
		}
		for vpcID, hostCnt := range vpcHostCnt {
			allVpcHostCnt[vpcID] = hostCnt
		}
	}

	// 更新返回结果里vpc对应的主机数
	for vpcID := range vpcToUpdate {
		vpcToUpdate[vpcID].VpcHostCount = allVpcHostCnt[vpcID]
	}

	return nil
}

func (lgc *Logics) UpdateSyncTask(kit *rest.Kit, taskID int64, option map[string]interface{}) error {

	// generate audit log.
	audit := auditlog.NewSyncTaskAuditLog(lgc.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(option)
	auditLog, auditErr := audit.GenerateAuditLog(generateAuditParameter, taskID, nil)
	if auditErr != nil {
		blog.Errorf("generate audit log failed before update sync task, taskID: %d, err: %v, rid: %s", taskID, auditErr, kit.Rid)
		return auditErr
	}

	// to update.
	err := lgc.CoreAPI.CoreService().Cloud().UpdateSyncTask(kit.Ctx, kit.Header, taskID, option)
	if err != nil {
		blog.Errorf("UpdateSyncTask failed, rid:%s, taskID:%d, option:%+v, err:%+v", kit.Rid, taskID, option, err)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("save audit log failed after update sync task, taskID: %d, err: %v, rid: %s", taskID, err, kit.Rid)
		return err
	}
	return nil
}

func (lgc *Logics) DeleteSyncTask(kit *rest.Kit, taskID int64) error {

	// generate audit log.
	audit := auditlog.NewSyncTaskAuditLog(lgc.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, auditErr := audit.GenerateAuditLog(generateAuditParameter, taskID, nil)
	if auditErr != nil {
		blog.Errorf("generate audit log failed before delete sync task, taskID: %d, err: %v, rid: %s", taskID,
			auditErr, kit.Rid)
		return auditErr
	}

	// to delete.
	err := lgc.CoreAPI.CoreService().Cloud().DeleteSyncTask(kit.Ctx, kit.Header, taskID)
	if err != nil {
		blog.Errorf("DeleteSyncTask failed, rid:%s, taskID:%d, err:%+v", kit.Rid, taskID, err)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("save audit log failed after delete sync task, taskID: %d, err: %v, rid: %s", taskID, err, kit.Rid)
		return err
	}

	return nil
}

func (lgc *Logics) CreateSyncHistory(kit *rest.Kit, history *metadata.SyncHistory) (*metadata.SyncHistory, error) {
	result, err := lgc.CoreAPI.CoreService().Cloud().CreateSyncHistory(kit.Ctx, kit.Header, history)
	if err != nil {
		blog.Errorf("CreateSyncHistory failed, rid:%s, history:%+v, err:%+v", kit.Rid, history, err)
		return nil, err
	}

	return result, nil
}

func (lgc *Logics) SearchSyncHistory(kit *rest.Kit, option *metadata.SearchSyncHistoryOption) (*metadata.MultipleSyncHistory, error) {
	// set default limit
	if option.Page.Limit == 0 {
		option.Page.Limit = common.BKDefaultLimit
	}
	if option.Page.IsIllegal() {
		blog.Errorf("SearchSyncHistory failed, Page is IsIllegal, rid:%s, page:%+v", kit.Rid, option.Page)
		return nil, kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded)
	}

	// set default sort
	if option.Page.Sort == "" {
		option.Page.Sort = "-" + common.CreateTimeField
	}

	// if fuzzy search, change the string query to regexp
	if option.IsFuzzy == true {
		for k, v := range option.Condition {
			field, ok := v.(string)
			if ok {
				option.Condition[k] = mapstr.MapStr{
					common.BKDBLIKE: params.SpecialCharChange(field),
					"$options":      "i",
				}
			}
		}
	}

	result, err := lgc.CoreAPI.CoreService().Cloud().SearchSyncHistory(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("SearchSyncHistory failed, rid:%s, option:%+v, err:%+v", kit.Rid, option, err)
		return nil, err
	}

	return result, nil
}

func (lgc *Logics) SearchSyncRegion(kit *rest.Kit, option *metadata.SearchSyncRegionOption) ([]metadata.SyncRegion, error) {
	accountConf, err := lgc.GetCloudAccountConf(kit, option.AccountID)
	if err != nil {
		blog.Errorf("SearchSyncRegion failed, rid:%s, option:%+v, err:%+v", kit.Rid, option, err)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	result, err := lgc.GetRegionsInfo(kit, *accountConf, option.WithHostCount)
	if err != nil {
		blog.Errorf("SearchSyncRegion failed, rid:%s, option:%+v, err:%+v", kit.Rid, option, err)
		return nil, kit.CCError.CCError(common.CCErrCloudRegionGetFail)
	}

	return result, nil
}
