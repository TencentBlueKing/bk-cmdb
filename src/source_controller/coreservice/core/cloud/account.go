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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (c *cloudOperation) CreateAccount(kit *rest.Kit, account *metadata.CloudAccount) (*metadata.CloudAccount, errors.CCErrorCoder) {
	if err := c.validCreateAccount(kit, account); nil != err {
		blog.Errorf("CreateAccount failed, valid error: %+v, rid: %s", err, kit.Rid)
		return nil, err
	}

	id, err := c.dbProxy.NextSequence(kit.Ctx, common.BKTableNameCloudAccount)
	if nil != err {
		blog.Errorf("CreateAccount failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	account.AccountID = int64(id)
	ts := time.Now()
	account.OwnerID = kit.SupplierAccount
	account.Creator = kit.User
	account.LastEditor = kit.User
	account.CreateTime = ts
	account.LastTime = ts

	err = c.dbProxy.Table(common.BKTableNameCloudAccount).Insert(kit.Ctx, account)
	if err != nil {
		blog.ErrorJSON("CreateAccount failed, db insert failed, accountName: %s, err: %s, rid: %s", account.AccountName, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
	}
	// 不返回bk_secret_key的值
	account.SecretKey = ""
	return account, nil
}

func (c *cloudOperation) SearchAccount(kit *rest.Kit, option *metadata.SearchCloudOption) (*metadata.MultipleCloudAccount, errors.CCErrorCoder) {
	accounts := []metadata.CloudAccount{}
	option.Condition = util.SetQueryOwner(option.Condition, kit.SupplierAccount)
	err := c.dbProxy.Table(common.BKTableNameCloudAccount).Find(option.Condition).Fields(option.Fields...).
		Start(uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).Sort(option.Page.Sort).All(kit.Ctx, &accounts)
	if err != nil {
		blog.ErrorJSON("search cloud account failed, option: %s, err: %s, rid: %s", option, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	// 不返回bk_secret_key的值
	for i, _ := range accounts {
		accounts[i].SecretKey = ""
	}

	// 账户总个数
	count, err := c.countAccount(kit, option.Condition)
	if err != nil {
		blog.ErrorJSON("SearchAccount countAccount error %s, option: %s, rid: %s", err, option.Condition, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	// 获取账户是否能被删除的信息
	accountTaskcntMap, err := c.getAcccountTaskcntMap(kit)
	if err != nil {
		blog.ErrorJSON("getAcccountTaskcntMap error %s, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	results := make([]metadata.CloudAccountWithExtraInfo, 0)
	for i, _ := range accounts {
		canDeleteAccount := false
		// 账户下没有同步任务，则可以删除
		if accountTaskcntMap[accounts[i].AccountID] == 0 {
			canDeleteAccount = true
		}
		results = append(results, metadata.CloudAccountWithExtraInfo{CloudAccount: accounts[i], CanDeleteAccount: canDeleteAccount})
	}

	return &metadata.MultipleCloudAccount{Count: int64(count), Info: results}, nil
}

func (c *cloudOperation) UpdateAccount(kit *rest.Kit, accountID int64, option mapstr.MapStr) errors.CCErrorCoder {
	if err := c.validUpdateAccount(kit, accountID, option); nil != err {
		blog.Errorf("UpdateAccount failed, valid error: %+v, rid: %s", err, kit.Rid)
		return err
	}
	filter := map[string]interface{}{common.BKCloudAccountID: accountID}
	filter = util.SetModOwner(filter, kit.SupplierAccount)
	option.Set(common.BKLastEditor, kit.User)
	option.Set(common.LastTimeField, time.Now())
	// 确保不会更新云厂商类型、云账户id、开发商id
	option.Remove(common.BKCloudVendor)
	option.Remove(common.BKCloudIDField)
	option.Remove(common.BKOwnerIDField)
	if e := c.dbProxy.Table(common.BKTableNameCloudAccount).Update(kit.Ctx, filter, option); e != nil {
		blog.Errorf("UpdateAccount failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameCloudAccount, filter, e, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

func (c *cloudOperation) DeleteAccount(kit *rest.Kit, accountID int64) errors.CCErrorCoder {
	if err := c.validDeleteAccount(kit, accountID); nil != err {
		blog.Errorf("DeleteAccount failed, valid error: %+v, rid: %s", err, kit.Rid)
		return err
	}

	filter := map[string]interface{}{common.BKCloudAccountID: accountID}
	filter = util.SetModOwner(filter, kit.SupplierAccount)
	if e := c.dbProxy.Table(common.BKTableNameCloudAccount).Delete(kit.Ctx, filter); e != nil {
		blog.Errorf("DeleteAccount failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameCloudAccount, filter, e, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
	}
	return nil
}

// 查询云厂商账户配置
func (c *cloudOperation) SearchAccountConf(kit *rest.Kit, option *metadata.SearchCloudOption) (*metadata.MultipleCloudAccountConf, errors.CCErrorCoder) {
	accountconfs := []metadata.CloudAccountConf{}
	option.Condition = util.SetQueryOwner(option.Condition, kit.SupplierAccount)
	err := c.dbProxy.Table(common.BKTableNameCloudAccount).Find(option.Condition).Fields(option.Fields...).
		Start(uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).Sort(option.Page.Sort).All(kit.Ctx, &accountconfs)
	if err != nil {
		blog.ErrorJSON("SearchAccount failed, db insert failed, option: %s, err: %s, rid: %s", option, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	// 账户总个数
	count, err := c.countAccount(kit, option.Condition)
	if err != nil {
		blog.ErrorJSON("SearchAccount countAccount error %s, option: %s, rid: %s", err, option.Condition, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return &metadata.MultipleCloudAccountConf{Count: int64(count), Info: accountconfs}, nil
}

// 获取账户总个数
func (c *cloudOperation) countAccount(kit *rest.Kit, cond mapstr.MapStr) (uint64, error) {
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)
	count, err := c.dbProxy.Table(common.BKTableNameCloudAccount).Find(cond).Count(kit.Ctx)
	return count, err

}

// 获取各个账户的同步任务数
func (c *cloudOperation) getAcccountTaskcntMap(kit *rest.Kit) (map[int64]int64, error) {
	accountTaskcntMap := make(map[int64]int64)
	aggregateCond := []interface{}{
		map[string]interface{}{
			"$group": map[string]interface{}{
				"_id": "$" + common.BKCloudAccountID,
				"num": map[string]interface{}{"$sum": 1},
			},
		},
	}
	type accountTaskcnt struct {
		AccountID int64 `bson:"_id"`
		Num       int64 `bson:"num"`
	}
	resultAll := make([]accountTaskcnt, 0)
	if err := c.dbProxy.Table(common.BKTableNameCloudSyncTask).AggregateAll(kit.Ctx, aggregateCond, &resultAll); err != nil {
		return nil, err
	}

	for _, r := range resultAll {
		accountTaskcntMap[r.AccountID] = r.Num
	}

	return accountTaskcntMap, nil
}
