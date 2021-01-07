/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cloud

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (c *cloudOperation) validCreateAccount(kit *rest.Kit, account *metadata.CloudAccount) errors.CCErrorCoder {
	// cloud vendor check
	if !util.InStrArr(metadata.SupportedCloudVendors, string(account.CloudVendor)) {
		blog.ErrorJSON("[validCreateAccount] not support cloud vendor: %s, rid: %s", account.CloudVendor, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCloudVendorNotSupport)
	}

	// account name unique check
	cond := mapstr.MapStr{common.BKCloudAccountName: account.AccountName}
	count, err := c.countAccount(kit, cond)
	if nil != err {
		blog.ErrorJSON("[validCreateAccount] count account error %s, condition: %s, rid: %s", err, cond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.ErrorJSON("[validCreateAccount] account name already exist, bk_account_name: %s, rid: %s", account.AccountName, kit.Rid)
		return kit.CCError.CCError(common.CCErrCloudAccountNameAlreadyExist)
	}

	// SecretID check, one SecretID can only have one account
	option := &metadata.SearchCloudOption{Condition: mapstr.MapStr{common.BKSecretID: account.SecretID}}
	multiAccount, err := c.SearchAccount(kit, option)
	if nil != err {
		blog.ErrorJSON("[validCreateAccount] SearchAccount error %s, option: %s, rid: %s", err, option, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if len(multiAccount.Info) > 0 {
		blog.ErrorJSON("[validCreateAccount] same SecretID %s has already exist in cloud account:%s, rid: %s", account.SecretID, multiAccount.Info[0].AccountName, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCloudAccountSecretIDAlreadyExist, multiAccount.Info[0].AccountName)
	}

	return nil
}

func (c *cloudOperation) validUpdateAccount(kit *rest.Kit, accountID int64, option mapstr.MapStr) errors.CCErrorCoder {
	// accountID exist check
	if err := c.validAccountExist(kit, accountID); err != nil {
		return err
	}
	// cloud vendor check
	if option.Exists(common.BKCloudVendor) {
		cloudVendor, err := option.String(common.BKCloudVendor)
		if err != nil {
			blog.ErrorJSON("[validUpdateAccount] not invalid cloud vendor, option: %s, rid: %s", option, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCloudValidAccountParamFail, common.BKCloudVendor)
		}
		if !util.InStrArr(metadata.SupportedCloudVendors, cloudVendor) {
			blog.ErrorJSON("[validUpdateAccount] not support cloud vendor: %s, rid: %s", cloudVendor, kit.Rid)
			return kit.CCError.CCError(common.CCErrCloudVendorNotSupport)
		}
	}

	// account name unique check
	if option.Exists(common.BKCloudAccountName) {
		cloudAccountName, err := option.String(common.BKCloudAccountName)
		if err != nil {
			blog.ErrorJSON("[validUpdateAccount] not invalid cloud vendor, option: %s, rid: %s", option, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCloudValidAccountParamFail, cloudAccountName)
		}
		cond := mapstr.MapStr{common.BKCloudAccountName: cloudAccountName,
			common.BKCloudAccountID: map[string]interface{}{common.BKDBNE: accountID}}
		count, err := c.countAccount(kit, cond)
		if nil != err {
			blog.ErrorJSON("[validUpdateAccount] count account error %s, condition: %s, rid: %s", err, cond, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		if count > 0 {
			blog.ErrorJSON("[validUpdateAccount] account name already exist, bk_account_name: %s, rid: %s", cloudAccountName, kit.Rid)
			return kit.CCError.CCError(common.CCErrCloudAccountNameAlreadyExist)
		}
	}

	return nil
}

func (c *cloudOperation) validDeleteAccount(kit *rest.Kit, accountID int64) errors.CCErrorCoder {
	// accountID exist check
	if err := c.validAccountExist(kit, accountID); err != nil {
		return err
	}

	// accountID has cloud sync task check
	if err := c.validAccountHasTask(kit, accountID); err != nil {
		return err
	}
	return nil
}

func (c *cloudOperation) validAccountExist(kit *rest.Kit, accountID int64) errors.CCErrorCoder {
	cond := mapstr.MapStr{common.BKCloudAccountID: accountID}
	count, err := c.countAccount(kit, cond)
	if nil != err {
		blog.ErrorJSON("[validAccountExist] count account error %s, condition: %s, rid: %s", err, cond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count <= 0 {
		blog.ErrorJSON("[validAccountExist] no account exist, bk_account_id: %s", accountID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCloudAccountIDNoExistFail, accountID)
	}

	return nil
}

// validAccountHasTask valid whether the account has any cloud sync task
func (c *cloudOperation) validAccountHasTask(kit *rest.Kit, accountID int64) errors.CCErrorCoder {
	accountTaskcntMap, err := c.getAcccountTaskcntMap(kit)
	if err != nil {
		blog.ErrorJSON("getAcccountTaskcntMap error %s, rid: %s", err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if accountTaskcntMap[accountID] > 0 {
		return kit.CCError.CCError(common.CCErrCloudAccountDeletedFailedForSyncTask)
	}

	return nil
}

func (c *cloudOperation) validCreateSyncTask(kit *rest.Kit, task *metadata.CloudSyncTask) errors.CCErrorCoder {
	// accountID check
	cond := mapstr.MapStr{common.BKCloudAccountID: task.AccountID}
	count, err := c.countAccount(kit, cond)
	if nil != err {
		blog.ErrorJSON("[validCreateSyncTask] accountID valid failed, error %s, condition: %s, rid: %s", err, cond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count <= 0 {
		blog.ErrorJSON("[validCreateSyncTask] accountID: %s does not exist, bk_task_name: %s, rid: %s", task.AccountID, task.TaskName, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCloudValidSyncTaskParamFail, common.BKCloudAccountID)
	}

	// task name unique check
	cond = mapstr.MapStr{common.BKCloudSyncTaskName: task.TaskName}
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)
	taskCount, err := c.dbProxy.Table(common.BKTableNameCloudSyncTask).Find(cond).Count(kit.Ctx)
	if nil != err {
		blog.ErrorJSON("[validCreateSycTask] count task error %s, condition: %s, rid: %s", err, cond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if taskCount > 0 {
		blog.ErrorJSON("[validCreateSycTask] task name already exist, bk_task_name: %s, rid: %s", task.TaskName, kit.Rid)
		return kit.CCError.CCError(common.CCErrCloudSyncTaskNameAlreadyExist)
	}

	if err := c.validSyncVpcInfo(kit, task.SyncVpcs); err != nil {
		blog.ErrorJSON("validUpdateSyncTask failed, error %s, syncVpcs:%s, rid: %s", err, task.SyncVpcs, kit.Rid)
		return err
	}

	// account task count check, one account can only have one task
	option := &metadata.SearchCloudOption{Condition: mapstr.MapStr{common.BKCloudAccountID: task.AccountID}}
	multiTask, err := c.SearchSyncTask(kit, option)
	if nil != err {
		blog.ErrorJSON("[validCreateSycTask] SearchSyncTask error %s, option: %s, rid: %s", err, option, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if len(multiTask.Info) > 0 {
		blog.ErrorJSON("[validCreateSycTask] this cloud account has had cloud sync task %s, rid: %s", multiTask.Info[0].TaskName, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCloudTaskAlreadyExistInAccount, multiTask.Info[0].TaskName)
	}

	return nil
}

func (c *cloudOperation) validUpdateSyncTask(kit *rest.Kit, taskID int64, option mapstr.MapStr) errors.CCErrorCoder {
	// accountID check
	if option.Exists(common.BKCloudAccountID) {
		cond := mapstr.MapStr{common.BKCloudAccountID: option[common.BKCloudAccountID]}
		count, err := c.countAccount(kit, cond)
		if nil != err {
			blog.ErrorJSON("[validCreateSyncTask] accountID valid failed, error %s, condition: %s, rid: %s", err, cond, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		if count <= 0 {
			blog.ErrorJSON("[validCreateSyncTask] accountID: %s does not exist, rid: %s", option[common.BKCloudAccountID], kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCloudValidSyncTaskParamFail, common.BKCloudAccountID)
		}
	}

	// task name unique check
	if option.Exists(common.BKCloudSyncTaskName) {
		taskName, err := option.String(common.BKCloudSyncTaskName)
		if err != nil {
			blog.ErrorJSON("[validUpdateSyncTask] not invalid task name, option: %s, rid: %s", option, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCloudValidSyncTaskParamFail, taskName)
		}
		cond := mapstr.MapStr{common.BKCloudSyncTaskName: taskName,
			common.BKCloudSyncTaskID: map[string]interface{}{common.BKDBNE: taskID}}
		count, err := c.countTask(kit, cond)
		if nil != err {
			blog.ErrorJSON("[validUpdateSyncTask] count task name failed error %s, condition: %s, rid: %s", err, cond, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		if count > 0 {
			blog.ErrorJSON("[validUpdateSyncTask] task name already exist, bk_account_name: %s", taskName, kit.Rid)
			return kit.CCError.CCError(common.CCErrCloudSyncTaskNameAlreadyExist)
		}
	}

	if vpcInfo, ok := option.Get(common.BKCloudSyncVpcs); ok {
		bs, err := json.Marshal(vpcInfo)
		if err != nil {
			blog.ErrorJSON("validUpdateSyncTask failed, error %s, vpcInfo:%s, rid: %s", err, vpcInfo, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommJSONMarshalFailed)
		}
		syncVpcs := make([]metadata.VpcSyncInfo, 0)
		err = json.Unmarshal(bs, &syncVpcs)
		if err != nil {
			blog.ErrorJSON("validUpdateSyncTask failed, error %s, vpcInfo:%s, rid: %s", err, vpcInfo, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
		}

		if err := c.validSyncVpcInfo(kit, syncVpcs); err != nil {
			blog.ErrorJSON("validUpdateSyncTask failed, error %s, vpcInfo:%s, rid: %s", err, vpcInfo, kit.Rid)
			return err
		}

	}

	return nil
}

// Valid sync vpc info
func (c *cloudOperation) validSyncVpcInfo(kit *rest.Kit, syncVpcs []metadata.VpcSyncInfo) errors.CCErrorCoder {
	if len(syncVpcs) == 0 {
		return nil
	}

	// vpcID is required
	for _, vpc := range syncVpcs {
		if vpc.VpcID == "" {
			blog.ErrorJSON("validUpdateSyncTask failed, rid: %s", kit.Rid)
			return kit.CCError.CCError(common.CCErrCloudVpcIDIsRequired)
		}
	}

	// resource dir must be exist
	if err := c.validResourceDirExist(kit, syncVpcs); err != nil {
		blog.ErrorJSON("validUpdateSyncTask failed, err:%s, rid: %s", err, kit.Rid)
		return err
	}

	// cloudID must be exist
	if err := c.validCloudIDExist(kit, syncVpcs); err != nil {
		blog.ErrorJSON("validCloudIDExist failed, err:%s, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// Valid resource dir which must be exist
func (c *cloudOperation) validResourceDirExist(kit *rest.Kit, syncVpcs []metadata.VpcSyncInfo) errors.CCErrorCoder {
	syncDirs := make(map[int64]bool)
	for _, syncInfo := range syncVpcs {
		syncDirs[syncInfo.SyncDir] = true
	}
	if len(syncDirs) == 0 {
		blog.ErrorJSON("validResourceDirExist failed, no sync dir is chosen, rid: %s", kit.Rid)
		return kit.CCError.CCError(common.CCErrCloudSyncDirNoChosen)
	}
	cond := mapstr.MapStr{}
	cond[common.BKDBOR] = []mapstr.MapStr{{common.BKDefaultField: 1}, {common.BKDefaultField: 4}}
	result := make([]struct {
		DirID int64 `json:"bk_module_id" bson:"bk_module_id"`
	}, 0)
	err := c.dbProxy.Table(common.BKTableNameBaseModule).Find(cond).Fields(common.BKModuleIDField).All(kit.Ctx, &result)
	if err != nil {
		blog.ErrorJSON("validResourceDirExist failed, err: %s, cond:%s, rid: %s", err.Error(), cond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	moduleIDs := make(map[int64]bool)
	for _, dir := range result {
		moduleIDs[dir.DirID] = true
	}
	for dir := range syncDirs {
		if _, ok := moduleIDs[dir]; !ok {
			blog.ErrorJSON("validResourceDirExist failed, syncDir %d not in moduleIDs, cond:%s, rid: %s", dir, cond, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCloudSyncDirNoExist, dir)
		}
	}

	return nil
}

// Valid cloudID which must be exist
func (c *cloudOperation) validCloudIDExist(kit *rest.Kit, syncVpcs []metadata.VpcSyncInfo) errors.CCErrorCoder {
	cloudIDs := make(map[int64]bool)
	for _, syncInfo := range syncVpcs {
		if syncInfo.CloudID == 0 {
			blog.ErrorJSON("validCloudIDExist failed, can't be default cloud area, rid: %s", kit.Rid)
			return kit.CCError.CCError(common.CCErrDefaultCloudIDProvided)
		}
		cloudIDs[syncInfo.CloudID] = true
	}
	if len(cloudIDs) == 0 {
		blog.ErrorJSON("validCloudIDExist failed, no cloudID is provided, rid: %s", kit.Rid)
		return kit.CCError.CCError(common.CCErrCloudIDNoProvided)
	}

	cloudIDArr := make([]int64, 0)
	for id := range cloudIDs {
		cloudIDArr = append(cloudIDArr, id)
	}
	cond := mapstr.MapStr{common.BKCloudIDField: map[string]interface{}{
		common.BKDBIN: cloudIDArr,
	}}
	result := make([]struct {
		CloudID int64 `json:"bk_cloud_id" bson:"bk_cloud_id"`
	}, 0)
	err := c.dbProxy.Table(common.BKTableNameBasePlat).Find(cond).Fields(common.BKCloudIDField).All(kit.Ctx, &result)
	if err != nil {
		blog.ErrorJSON("validCloudIDExist failed, err: %s, cond:%s, rid: %s", err.Error(), cond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	allIDs := make(map[int64]bool)
	for _, r := range result {
		allIDs[r.CloudID] = true
	}
	for id := range cloudIDs {
		if _, ok := allIDs[id]; !ok {
			blog.ErrorJSON("validCloudIDExist failed, cloudID %d is not exist, cond:%s, rid: %s", id, cond, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCloudIDNoExist, id)
		}
	}

	return nil
}
