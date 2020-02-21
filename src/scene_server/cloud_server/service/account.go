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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
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
	switch account.CloudVendor {
	case metadata.AWS:
		pass, err = s.Logics.AwsAccountVerify(account.SecretID, account.SecretKey)
		if err != nil {
			blog.ErrorJSON("aws cloud account verify failed, err :%v, rid: %s", err, ctx.Kit.Rid)
		}
	case metadata.TencentCloud:
		pass, err = s.Logics.TecentCloudVerify(account.SecretID, account.SecretKey)
		if err != nil {
			blog.ErrorJSON("tencent cloud account verify failed, err :%v, rid: %s", err, ctx.Kit.Rid)
		}
	default:
		ctx.RespErrorCodeOnly(common.CCErrCloudVendorNotSupport, "VerifyConnectivity failed, not support cloud vendor, rid: %v", ctx.Kit.Rid)
		return
	}

	rData := mapstr.MapStr{
		"connected": true,
		"error_msg": "",
	}
	if pass == false {
		rData["connected"] = false
		rData["error_msg"] = err.Error()
	}

	ctx.RespEntity(rData)
}

// 新建云账户
func (s *Service) CreateAccount(ctx *rest.Contexts) {
	account := new(metadata.CloudAccount)
	if err := ctx.DecodeInto(account); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// account name unique check
	// todo

	// accountType check
	if !util.InStrArr(metadata.SupportCloudVendors, string(account.CloudVendor)) {
		ctx.RespErrorCodeOnly(common.CCErrCloudVendorNotSupport, "CreateAccount failed, not support cloud vendor: %s, rid: %v", account.CloudVendor, ctx.Kit.Rid)
		return
	}

	res, err := s.CoreAPI.CoreService().Cloud().CreateAccount(ctx.Kit.Ctx, ctx.Kit.Header, account)
	if err != nil {
		blog.ErrorJSON("CreateAccount failed, core service CreateAccount failed, account info: %v, err: %s, rid: %s", account, err.Error(), ctx.Kit.Rid)
		return
	}

	ctx.RespEntity(res)
}

// 查询云账户
func (s *Service) SearchAccount(ctx *rest.Contexts) {
	ctx.RespEntity("SearchAccount")
}

// 更新云账户
func (s *Service) UpdateAccount(ctx *rest.Contexts) {
	ctx.RespEntity("UpdateAccount")
}

// 删除云账户
func (s *Service) DeleteAccount(ctx *rest.Contexts) {
	ctx.RespEntity("DeleteAccount")
}
