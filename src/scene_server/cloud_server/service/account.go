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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
	ccom "configcenter/src/scene_server/cloud_server/common"
)

// 云账户连通测试
func (s *Service) VerifyConnectivity(ctx *rest.Contexts) {
	account := new(metadata.CloudAccountVerify)
	if err := ctx.DecodeInto(account); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var pass bool
	var err error
	conf := ccom.AccountConf{account.CloudVendor, account.SecretID, account.SecretKey}
	pass, err = s.Logics.AccountVerify(conf)
	if err != nil {
		blog.ErrorJSON("cloud account verify failed, cloudvendor:%s, err :%v, rid: %s", account.CloudVendor, err, ctx.Kit.Rid)
	}

	rData := mapstr.MapStr{
		"connected": true,
		"error_msg": "",
	}
	if pass == false {
		rData["connected"] = false
		rData["error_msg"] = err.Error()
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

	res, err := s.CoreAPI.CoreService().Cloud().CreateAccount(ctx.Kit.Ctx, ctx.Kit.Header, account)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// add auditLog
	auditLog := s.Logics.NewAccountAuditLog(ctx.Kit, ctx.Kit.SupplierAccount)
	if err := auditLog.WithCurrent(ctx.Kit, res.AccountID); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if err := auditLog.SaveAuditLog(ctx.Kit, metadata.AuditCreate); err != nil {
		ctx.RespAutoError(err)
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

	// set default sort
	if option.Page.Sort == "" {
		option.Page.Sort = "-" + common.CreateTimeField
	}

	if option.Page.IsIllegal() {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded))
		return
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

	// auditLog preData
	auditLog := s.Logics.NewAccountAuditLog(ctx.Kit, ctx.Kit.SupplierAccount)
	if err := auditLog.WithPrevious(ctx.Kit, accountID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	option := map[string]interface{}{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}
	err = s.CoreAPI.CoreService().Cloud().UpdateAccount(ctx.Kit.Ctx, ctx.Kit.Header, accountID, option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// add auditLog
	if err := auditLog.WithCurrent(ctx.Kit, accountID); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if err := auditLog.SaveAuditLog(ctx.Kit, metadata.AuditUpdate); err != nil {
		ctx.RespAutoError(err)
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

	// add auditLog
	auditLog := s.Logics.NewAccountAuditLog(ctx.Kit, ctx.Kit.SupplierAccount)
	if err := auditLog.WithPrevious(ctx.Kit, accountID); err != nil {
		ctx.RespAutoError(err)
		return
	}
	err = s.CoreAPI.CoreService().Cloud().DeleteAccount(ctx.Kit.Ctx, ctx.Kit.Header, accountID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := auditLog.SaveAuditLog(ctx.Kit, metadata.AuditDelete); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}
