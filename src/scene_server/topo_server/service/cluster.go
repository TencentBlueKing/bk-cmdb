/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// SearchClusters query based on user-specified criteria cluster
func (s *Service) SearchClusters(ctx *rest.Contexts) {

	searchCond := new(types.QueryClusterOption)
	if err := ctx.DecodeInto(searchCond); err != nil {
		blog.Errorf("failed to parse the params, error: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	if cErr := searchCond.Validate(); cErr.ErrCode != 0 {
		blog.Errorf("validate request failed, err: %v, rid: %s", cErr, ctx.Kit.Rid)
		ctx.RespAutoError(cErr.ToCCError(ctx.Kit.CCError))
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	filter := mapstr.New()
	if searchCond.Filter != nil {
		cond, errKey, rawErr := searchCond.Filter.ToMgo()
		if rawErr != nil {
			blog.Errorf("parse biz filter(%#v) failed, err: %v, rid: %s", searchCond.Filter, rawErr, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey))
			return
		}
		filter = cond
	}
	// 无论条件中是否有bk_biz_id、supplier_account,这里统一替换成url中的bk_biz_id 和kit中的supplier_account
	filter[types.BKBizIDField] = bizID
	filter[types.BKSupplierAccountField] = ctx.Kit.SupplierAccount

	// count biz in cluster enable count is set
	if searchCond.Page.EnableCount {
		filter := []map[string]interface{}{filter}
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBaseCluster, filter)
		if err != nil {
			blog.Errorf("count biz failed, err: %v, cond: %#v, rid: %s", err, filter, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithCount(counts[0], make([]mapstr.MapStr, 0))
		return
	}

	query := &metadata.QueryCondition{
		Condition: filter,
		Page:      searchCond.Page,
		Fields:    searchCond.Fields,
	}
	result, err := s.Engine.CoreAPI.CoreService().Kube().SearchCluster(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("search cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
	ctx.RespEntityWithCount(0, result.Data)
}

// UpdateClusterFields update cluster fields
func (s *Service) UpdateClusterFields(ctx *rest.Contexts) {

	data := new(types.UpdateClusterOption)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if err := data.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	clusterIDs := make([]int64, 0)
	for _, cluster := range data.Clusters {
		clusterIDs = append(clusterIDs, cluster.ID)
	}

	cond := map[string]interface{}{
		types.BKIDField:       map[string]interface{}{common.BKDBIN: clusterIDs},
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: ctx.Kit.SupplierAccount,
	}

	input := &metadata.QueryCondition{
		Condition: cond,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}

	result, err := s.Engine.CoreAPI.CoreService().Kube().SearchCluster(ctx.Kit.Ctx, ctx.Kit.Header, input)
	if err != nil {
		blog.Errorf("search cluster failed, input: %#v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(result.Data) != len(clusterIDs) {
		blog.Errorf("the number of instances obtained is inconsistent with the param, bizID: %d, ids: %#v, err: %v, "+
			"rid: %s", bizID, clusterIDs, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrTopoInstUpdateFailed))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Engine.CoreAPI.CoreService().Kube().UpdateClusterFields(ctx.Kit.Ctx, ctx.Kit.Header,
			ctx.Kit.SupplierAccount, bizID, data)
		if err != nil {
			blog.Errorf("create cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		// for audit log.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		audit := auditlog.NewKubeAudit(s.Engine.CoreAPI.CoreService())

		for id, cluster := range result.Data {
			updateFields, err := mapstr.Struct2Map(data.Clusters[id].Data)
			if err != nil {
				blog.Errorf("update fields convert failed, data: %+v, err: %v, rid: %s", data.Clusters[id].Data,
					err, ctx.Kit.Rid)
				return err
			}

			generateAuditParameter.WithUpdateFields(updateFields)
			auditLog, err := audit.GenerateClusterAuditLog(generateAuditParameter, []types.Cluster{cluster})
			if err != nil {
				blog.Errorf("create cluster, generate audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}

			err = audit.SaveAuditLog(ctx.Kit, auditLog...)
			if err != nil {
				blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return ctx.Kit.CCError.Error(common.CCErrAuditSaveLogFailed)
			}
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// CreateCluster create a container cluster
func (s *Service) CreateCluster(ctx *rest.Contexts) {
	data := new(types.Cluster)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.CreateValidate(); err != nil {
		blog.Errorf("validate create container cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	var id int64
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		id, err = s.Logics.ContainerOperation().CreateCluster(ctx.Kit, data, bizID, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("create cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(id)
}

// DeleteCluster delete cluster.
func (s *Service) DeleteCluster(ctx *rest.Contexts) {
	option := new(types.DeleteClusterOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		err = s.Logics.ContainerOperation().DeleteCluster(ctx.Kit, bizID, option, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("delete cluster failed, biz: %d, option: %+v, err: %v, rid: %s", bizID, option, err,
				ctx.Kit.Rid)
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
