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
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
)

func (c *cloudOperation) validCreateAccount(kit *rest.Kit, account *metadata.CloudAccount) errors.CCErrorCoder {
	// accountType check
	if !util.InStrArr(metadata.SupportCloudVendors, string(account.CloudVendor)) {
		blog.ErrorJSON("[validCreateAccount] not support cloud vendor: %s, rid: %v", account.CloudVendor, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCloudVendorNotSupport)
	}

	// account name unique check
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKCloudAccountName, Val: account.AccountName})
	count, err := c.countAccount(kit, cond.ToMapStr())
	if nil != err {
		blog.ErrorJSON("[validCreateAccount] count account error %v, condition: %#v, rid: %s", err, cond.ToMapStr(), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCloudValidAccountParamFail, "count account failed, err: %v", err)
	}
	if count > 0 {
		blog.ErrorJSON("[validCreateAccount] account name already exist, bk_account_name: %s", account.AccountName, kit.Rid)
		return kit.CCError.CCError(common.CCErrCloudAccountNameAlreadyExist)
	}

	return nil
}

func (c *cloudOperation) validUpdateAccount(kit *rest.Kit, account *metadata.CloudAccount) errors.CCErrorCoder {
	// accountType check
	if !util.InStrArr(metadata.SupportCloudVendors, string(account.CloudVendor)) {
		blog.ErrorJSON("[validUpdateAccount] not support cloud vendor: %s, rid: %v", account.CloudVendor, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCloudVendorNotSupport)
	}

	// account name unique check
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKCloudAccountName, Val: account.AccountName})
	cond.Element(&mongo.Neq{Key: common.BKCloudAccountIDField, Val: account.AccountID})
	count, err := c.countAccount(kit, cond.ToMapStr())
	if nil != err {
		blog.ErrorJSON("[validUpdateAccount] count account error %v, condition: %#v, rid: %s", err, cond.ToMapStr(), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCloudValidAccountParamFail, "count account failed, err: %v", err)
	}
	if count > 0 {
		blog.ErrorJSON("[validUpdateAccount] account name already exist, bk_account_name: %s", account.AccountName, kit.Rid)
		return kit.CCError.CCError(common.CCErrCloudAccountNameAlreadyExist)
	}

	return nil
}
