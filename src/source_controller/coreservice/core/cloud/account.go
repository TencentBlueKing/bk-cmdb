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
	"time"
)

func (c *cloudOperation) CreateAccount(kit *rest.Kit, account *metadata.CloudAccount) (*metadata.CloudAccount, errors.CCErrorCoder) {
	ts := time.Now()
	account.OwnerID = kit.SupplierAccount
	account.CreateTime = ts
	account.LastTime = ts

	err := c.dbProxy.Table(common.BKTableNameCloudAccount).Insert(kit.Ctx, account)
	if err != nil {
		blog.ErrorJSON("CreateAccount failed, db insert failed, accountName: %s, err: %v, rid: %s", account.AccountName, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}
	return account, nil
}

func (c *cloudOperation) SearchAccount(kit *rest.Kit, option *metadata.SearchCloudAccountOption) (*metadata.MultipleCloudAccount, errors.CCErrorCoder) {
	return &metadata.MultipleCloudAccount{}, nil
}

func (c *cloudOperation) UpdateAccount(kit *rest.Kit, accountID int64, account *metadata.CloudAccount) (*metadata.CloudAccount, errors.CCErrorCoder) {
	return &metadata.CloudAccount{}, nil
}

func (c *cloudOperation) DeleteAccount(kit *rest.Kit, accountID int64) errors.CCErrorCoder {
	return nil
}
