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
		blog.ErrorJSON("[validCreateAccount] account name already exist, bk_account_name: %s", account.AccountName, kit.Rid)
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
			blog.ErrorJSON("[validUpdateAccount] account name already exist, bk_account_name: %s", cloudAccountName, kit.Rid)
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

func (c *cloudOperation) validCreateSyncTask(kit *rest.Kit, task *metadata.CloudSyncTask) errors.CCErrorCoder {
	// accountID check
	cond := mapstr.MapStr{common.BKCloudAccountID: task.AccountID}
	count, err := c.countAccount(kit, cond)
	if nil != err {
		blog.ErrorJSON("[validCreateSyncTask] accountID valid failed, error %s, condition: %s, rid: %s", err, cond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count <= 0 {
		blog.ErrorJSON("[validCreateSyncTask] accountID: %s does not exist, bk_task_name: %s", task.AccountID, task.TaskName, kit.Rid)
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
		blog.ErrorJSON("[validCreateSycTask] task name already exist, bk_task_name: %s", task.TaskName, kit.Rid)
		return kit.CCError.CCError(common.CCErrCloudSyncTaskNameAlreadyExist)
	}

	// vpcID is required
	if len(task.SyncVpcs) > 0 {
		for _, vpc := range task.SyncVpcs {
			if vpc.VpcID == "" {
				blog.ErrorJSON("[validCreateSycTask] vpcID filed is required, bk_task_name: %s", task.TaskName, kit.Rid)
				return kit.CCError.CCError(common.CCErrCloudVpcIDIsRequired)
			}
		}
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

	return nil
}
