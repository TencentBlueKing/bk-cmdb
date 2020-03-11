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
		pass, err = s.Logics.AwsAccountVerify(ctx.Kit, account.SecretID, account.SecretKey)
		if err != nil {
			blog.ErrorJSON("aws cloud account verify failed, err :%v, rid: %s", err, ctx.Kit.Rid)
		}
	case metadata.TencentCloud:
		pass, err = s.Logics.TecentCloudVerify(ctx.Kit, account.SecretID, account.SecretKey)
		if err != nil {
			blog.ErrorJSON("tencent cloud account verify failed, err :%v, rid: %s", err, ctx.Kit.Rid)
		}
	default:
		ctx.RespErrorCodeOnly(common.CCErrCloudVendorNotSupport, "VerifyConnectivity failed, not support cloud vendor: %v", account.CloudVendor)
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

	res, errCoder := s.CoreAPI.CoreService().Cloud().CreateAccount(ctx.Kit.Ctx, ctx.Kit.Header, account)
	if errCoder != nil {
		ctx.RespAutoError(errCoder)
		return
	}

	// audiLog should remove SecreteID,SecreteKey
	res.SecreteID = ""
	res.SecreteKey = ""
	curData := metadata.CloudAccountOpData{
		AccountName: res.AccountName,
		CloudVendor: res.CloudVendor,
		Description: res.Description,
	}
	auditLog := metadata.AuditLog{
		AuditType:    metadata.CloudResourceType,
		ResourceType: metadata.CloudAccountRes,
		Action:       metadata.AuditCreate,
		OperationDetail: &metadata.CloudAccountOpDetail{
			AccountID:  res.AccountID,
			CurData:    curData,
			Properties: metadata.CloudAccountAudiLogProperty,
		},
	}

	result, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(ctx.Kit.Ctx, ctx.Kit.Header, auditLog)
	if err != nil {
		blog.Errorf("CreateAccount success but create audit log failed, http failed, err:%v", err)
		return
	}
	if !result.Result {
		blog.Errorf("CreateAccount success but create audit log failed, err code:%d, err msg:%s", result.Code, result.ErrMsg)
		return
	}

	ctx.RespEntity(res)
}

// 查询云账户
func (s *Service) SearchAccount(ctx *rest.Contexts) {
	option := metadata.SearchCloudAccountOption{}
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
	accountIDStr := ctx.Request.PathParameter(common.BKCloudAccountIDField)
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudAccountIDField))
		return
	}

	// update cloud account audiLog preData
	cond := metadata.SearchCloudAccountOption{
		Condition: mapstr.MapStr{common.BKCloudAccountIDField: accountID},
	}
	res, err := s.CoreAPI.CoreService().Cloud().SearchAccount(ctx.Kit.Ctx, ctx.Kit.Header, &cond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if len(res.Info) <= 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCloudAccountIDNoExistFail))
		return
	}
	// should remove secreteID, secreteKey
	res.Info[0].SecreteKey = ""
	res.Info[0].SecreteID = ""
	preData := metadata.CloudAccountOpData{
		AccountName: res.Info[0].AccountName,
		CloudVendor: res.Info[0].CloudVendor,
		Description: res.Info[0].Description,
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

	// audiLog curData
	rsp, err := s.CoreAPI.CoreService().Cloud().SearchAccount(ctx.Kit.Ctx, ctx.Kit.Header, &cond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if len(res.Info) <= 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCloudAccountIDNoExistFail))
		return
	}
	// should remove secreteID, secreteKey
	rsp.Info[0].SecreteKey = ""
	rsp.Info[0].SecreteID = ""
	curData := metadata.CloudAccountOpData{
		AccountName: rsp.Info[0].AccountName,
		CloudVendor: rsp.Info[0].CloudVendor,
		Description: rsp.Info[0].Description,
	}
	// audiLog
	auditLog := metadata.AuditLog{
		AuditType:    metadata.CloudResourceType,
		ResourceType: metadata.CloudAccountRes,
		Action:       metadata.AuditUpdate,
		OperationDetail: &metadata.CloudAccountOpDetail{
			AccountID:  rsp.Info[0].AccountID,
			PreData:    preData,
			CurData:    curData,
			Properties: metadata.CloudAccountAudiLogProperty,
		},
	}
	result, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(ctx.Kit.Ctx, ctx.Kit.Header, auditLog)
	if err != nil {
		blog.Errorf("UpdateAccount success but create audit log failed, http failed, accountID: %d, err:%v", accountID, err)
		return
	}
	if !result.Result {
		blog.Errorf("UpdateAccount success but create audit log failed, accountID: %d, err code:%d, err msg:%s", accountID, result.Code, result.ErrMsg)
		return
	}

	ctx.RespEntity(nil)
}

// 删除云账户
func (s *Service) DeleteAccount(ctx *rest.Contexts) {
	//get accountID
	accountIDStr := ctx.Request.PathParameter(common.BKCloudAccountIDField)
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudAccountIDField))
		return
	}

	// delete account audiLog preData
	cond := metadata.SearchCloudAccountOption{
		Condition: mapstr.MapStr{common.BKCloudAccountIDField: accountID},
	}
	res, err := s.CoreAPI.CoreService().Cloud().SearchAccount(ctx.Kit.Ctx, ctx.Kit.Header, &cond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if len(res.Info) <= 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCloudAccountIDNoExistFail))
		return
	}
	preData := metadata.CloudAccountOpData{
		AccountName: res.Info[0].AccountName,
		CloudVendor: res.Info[0].CloudVendor,
		Description: res.Info[0].Description,
	}

	err = s.CoreAPI.CoreService().Cloud().DeleteAccount(ctx.Kit.Ctx, ctx.Kit.Header, accountID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// audiLog
	auditLog := metadata.AuditLog{
		AuditType:    metadata.CloudResourceType,
		ResourceType: metadata.CloudAccountRes,
		Action:       metadata.AuditDelete,
		OperationDetail: &metadata.CloudAccountOpDetail{
			AccountID:  res.Info[0].AccountID,
			PreData:    preData,
			Properties: metadata.CloudAccountAudiLogProperty,
		},
	}
	result, err := s.CoreAPI.CoreService().Audit().SaveAuditLog(ctx.Kit.Ctx, ctx.Kit.Header, auditLog)
	if err != nil {
		blog.Errorf("DeleteAccount success but create audit log failed, http failed, accountID: %d, err:%v", accountID, err)
		return
	}
	if !result.Result {
		blog.Errorf("DeleteAccount success but create audit log failed, accountID: %d, err code:%d, err msg:%s", accountID, result.Code, result.ErrMsg)
		return
	}

	ctx.RespEntity(nil)
}
