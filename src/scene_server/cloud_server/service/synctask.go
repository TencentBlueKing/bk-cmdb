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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
)

func (s *Service) SearchVpc(ctx *rest.Contexts) {
	vpcOpt := new(metadata.SearchVpcOption)
	if err := ctx.DecodeInto(vpcOpt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	accountIDStr := ctx.Request.PathParameter(common.BKCloudAccountID)
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudAccountID))
		return
	}

	accountConf, err := s.Logics.GetAccountConf(accountID)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	result, err := s.Logics.GetVpcHostCnt(vpcOpt.Region, *accountConf)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCloudVendorInterfaceCalledFailed))
		return
	}

	ctx.RespEntity(result)
}

func (s *Service) CreateSyncTask(ctx *rest.Contexts) {
	task := new(metadata.CloudSyncTask)
	if err := ctx.DecodeInto(task); err != nil {
		ctx.RespAutoError(err)
		return
	}

	res, err := s.CoreAPI.CoreService().Cloud().CreateSyncTask(ctx.Kit.Ctx, ctx.Kit.Header, task)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// add auditLog
	auditLog := s.Logics.NewSyncTaskAuditLog(ctx.Kit, ctx.Kit.SupplierAccount)
	if err := auditLog.WithCurrent(ctx.Kit, res.TaskID); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if err := auditLog.SaveAuditLog(ctx.Kit, metadata.AuditCreate); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

func (s *Service) SearchSyncTask(ctx *rest.Contexts) {
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

	res, err := s.CoreAPI.CoreService().Cloud().SearchSyncTask(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

func (s *Service) UpdateSyncTask(ctx *rest.Contexts) {
	taskIDStr := ctx.Request.PathParameter(common.BKCloudSyncTaskID)
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudSyncTaskID))
		return
	}

	option := map[string]interface{}{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// add auditLog preData
	auditLog := s.Logics.NewSyncTaskAuditLog(ctx.Kit, ctx.Kit.SupplierAccount)
	if err := auditLog.WithPrevious(ctx.Kit, taskID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err = s.CoreAPI.CoreService().Cloud().UpdateSyncTask(ctx.Kit.Ctx, ctx.Kit.Header, taskID, option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// add auditLog
	if err := auditLog.WithCurrent(ctx.Kit, taskID); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if err := auditLog.SaveAuditLog(ctx.Kit, metadata.AuditUpdate); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *Service) DeleteSyncTask(ctx *rest.Contexts) {
	taskIDStr := ctx.Request.PathParameter(common.BKCloudSyncTaskID)
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudSyncTaskID))
		return
	}

	// add auditLog preData
	auditLog := s.Logics.NewSyncTaskAuditLog(ctx.Kit, ctx.Kit.SupplierAccount)
	if err := auditLog.WithPrevious(ctx.Kit, taskID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err = s.CoreAPI.CoreService().Cloud().DeleteSyncTask(ctx.Kit.Ctx, ctx.Kit.Header, taskID)
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

func (s *Service) SearchSyncHistory(ctx *rest.Contexts) {
	option := metadata.SearchSyncHistoryOption{}
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

	res, err := s.CoreAPI.CoreService().Cloud().SearchSyncHistory(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

func (s *Service) SearchSyncRegion(ctx *rest.Contexts) {
	option := metadata.SearchSyncRegionOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	accountConf, err := s.Logics.GetAccountConf(option.AccountID)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	result, err := s.Logics.GetRegionsInfo(option.WithHostCount, *accountConf)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCloudVendorInterfaceCalledFailed))
		return
	}

	ctx.RespEntity(result)
}
