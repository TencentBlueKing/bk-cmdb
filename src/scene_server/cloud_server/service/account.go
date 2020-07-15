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
	"reflect"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
)

// 云账户连通测试
func (s *Service) VerifyConnectivity(ctx *rest.Contexts) {
	account := new(metadata.CloudAccountVerify)
	if err := ctx.DecodeInto(account); err != nil {
		ctx.RespAutoError(err)
		return
	}

	conf := metadata.CloudAccountConf{VendorName: account.CloudVendor, SecretID: account.SecretID, SecretKey: account.SecretKey}
	err := s.Logics.AccountVerify(conf)
	if err != nil {
		blog.ErrorJSON("cloud account verify failed, cloudvendor:%s, err :%v, rid: %s", account.CloudVendor, err, ctx.Kit.Rid)
		errStr := err.Error()
		if strings.Contains(strings.ToLower(errStr), "authfailure") {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCloudAccoutIDSecretWrong))
			return
		} else if strings.Contains(strings.ToLower(errStr), "timeout") {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCloudHttpRequestTimeout))
			return
		} else {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCloudVendorInterfaceCalledFailed))
			return
		}
	}

	ctx.RespEntity(nil)
}

// 新建云账户
func (s *Service) CreateAccount(ctx *rest.Contexts) {
	account := new(metadata.CloudAccount)
	if err := ctx.DecodeInto(account); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := account.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// 加密云账户密钥
	if s.cryptor != nil {
		secretKey, err := s.cryptor.Encrypt(account.SecretKey)
		if err != nil {
			blog.Errorf("CreateAccount failed, Encrypt err: %s, rid: %d", err, ctx.Kit.Rid)
			ctx.RespWithError(err, common.CCErrCloudAccountCreateFail, "CreateAccount Encrypt err")
			return
		}
		account.SecretKey = secretKey
	}

	var res *metadata.CloudAccount
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		res, err = s.CoreAPI.CoreService().Cloud().CreateAccount(ctx.Kit.Ctx, ctx.Kit.Header, account)
		if err != nil {
			blog.Errorf("CreateAccount failed, CreateAccount err:%s, rid:%s", err, ctx.Kit.Rid)
			return err
		}

		// add auditLog
		auditLog := s.Logics.NewAccountAuditLog(ctx.Kit, ctx.Kit.SupplierAccount)
		if err := auditLog.WithCurrent(ctx.Kit, res.AccountID); err != nil {
			blog.Errorf("CreateAccount failed, NewAccountAuditLog err:%s, rid:%s", err, ctx.Kit.Rid)
			return err
		}
		if err := auditLog.SaveAuditLog(ctx.Kit, metadata.AuditCreate); err != nil {
			blog.Errorf("CreateAccount failed, SaveAuditLog err:%s, rid:%s", err, ctx.Kit.Rid)
			return err
		}

		return nil
	})
	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(res)
}

// 查询云账户
func (s *Service) SearchAccount(ctx *rest.Contexts) {
	option := metadata.SearchCloudOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// set default limit
	if option.Page.Limit == 0 {
		option.Page.Limit = common.BKDefaultLimit
	}
	if option.Page.IsIllegal() {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded))
		return
	}

	// set default sort
	if option.Page.Sort == "" {
		option.Page.Sort = "-" + common.CreateTimeField
	}

	// if not exact search, change the string query to regexp
	if option.Exact != true {
		for k, v := range option.Condition {
			if reflect.TypeOf(v).Kind() == reflect.String {
				field := v.(string)
				option.Condition[k] = mapstr.MapStr{
					common.BKDBLIKE: params.SpecialCharChange(field),
					"$options":      "i",
				}
			}
		}
	}

	res, err := s.CoreAPI.CoreService().Cloud().SearchAccount(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

// 更新云账户
func (s *Service) UpdateAccount(ctx *rest.Contexts) {
	//get accountID
	accountIDStr := ctx.Request.PathParameter(common.BKCloudAccountID)
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudAccountID))
		return
	}

	option := map[string]interface{}{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		// auditLog preData
		auditLog := s.Logics.NewAccountAuditLog(ctx.Kit, ctx.Kit.SupplierAccount)
		if err := auditLog.WithPrevious(ctx.Kit, accountID); err != nil {
			blog.Errorf("UpdateAccount failed, NewAccountAuditLog err:%s, accountID:%d, option:%#v, rid:%s", err, accountID, option, ctx.Kit.Rid)
			return err
		}

		err = s.CoreAPI.CoreService().Cloud().UpdateAccount(ctx.Kit.Ctx, ctx.Kit.Header, accountID, option)
		if err != nil {
			blog.Errorf("UpdateAccount failed, UpdateAccount err:%s, accountID:%d, option:%#v, rid:%s", err, accountID, option, ctx.Kit.Rid)
			return err
		}

		// add auditLog
		if err := auditLog.WithCurrent(ctx.Kit, accountID); err != nil {
			blog.Errorf("UpdateAccount failed, WithCurrent err:%s, accountID:%d, option:%#v, rid:%s", err, accountID, option, ctx.Kit.Rid)
			return err
		}
		if err := auditLog.SaveAuditLog(ctx.Kit, metadata.AuditUpdate); err != nil {
			blog.Errorf("UpdateAccount failed, SaveAuditLog err:%s, accountID:%d, option:%#v, rid:%s", err, accountID, option, ctx.Kit.Rid)
			return err
		}

		return nil
	})
	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

// 删除云账户
func (s *Service) DeleteAccount(ctx *rest.Contexts) {
	//get accountID
	accountIDStr := ctx.Request.PathParameter(common.BKCloudAccountID)
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudAccountID))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		// add auditLog
		auditLog := s.Logics.NewAccountAuditLog(ctx.Kit, ctx.Kit.SupplierAccount)
		if err := auditLog.WithPrevious(ctx.Kit, accountID); err != nil {
			blog.Errorf("DeleteAccount failed, WithPrevious err:%s, accountID:%d, rid:%s", err, accountID, ctx.Kit.Rid)
			return err
		}
		err = s.CoreAPI.CoreService().Cloud().DeleteAccount(ctx.Kit.Ctx, ctx.Kit.Header, accountID)
		if err != nil {
			blog.Errorf("DeleteAccount failed, DeleteAccount err:%s, accountID:%d, rid:%s", err, accountID, ctx.Kit.Rid)
			return err
		}

		if err := auditLog.SaveAuditLog(ctx.Kit, metadata.AuditDelete); err != nil {
			blog.Errorf("DeleteAccount failed, SaveAuditLog err:%s, accountID:%d, rid:%s", err, accountID, ctx.Kit.Rid)
			return err
		}

		return nil
	})
	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}
