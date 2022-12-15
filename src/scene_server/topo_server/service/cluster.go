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
	"errors"
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
		cond, rawErr := searchCond.Filter.ToMgo()
		if rawErr != nil {
			blog.Errorf("parse cluster filter(%#v) failed, err: %v, rid: %s", searchCond.Filter, rawErr, ctx.Kit.Rid)
			ctx.RespAutoError(rawErr)
			return
		}
		filter = cond
	}

	filter[types.BKBizIDField] = bizID

	// get the number of clusters
	if searchCond.Page.EnableCount {
		cond := []map[string]interface{}{filter}
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBaseCluster, cond)
		if err != nil {
			blog.Errorf("count cluster failed, cond: %#v, err: %v, rid: %s", filter, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithCount(counts[0], make([]types.Cluster, 0))
		return
	}

	query := &metadata.QueryCondition{
		Condition:      filter,
		Page:           searchCond.Page,
		Fields:         searchCond.Fields,
		DisableCounter: true,
	}
	result, err := s.Engine.CoreAPI.CoreService().Kube().SearchCluster(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("search cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
	ctx.RespEntityWithCount(0, result.Data)
}

func (s *Service) getUpdateClustersInfo(kit *rest.Kit, bizID int64, clusterIDs []int64) ([]types.Cluster, error) {

	cond := map[string]interface{}{
		types.BKIDField:     map[string]interface{}{common.BKDBIN: clusterIDs},
		common.BKAppIDField: bizID,
	}

	input := &metadata.QueryCondition{
		Condition: cond,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}

	result, err := s.Engine.CoreAPI.CoreService().Kube().SearchCluster(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("search cluster failed, input: %#v, err: %v, rid: %s", input, err, kit.Rid)
		return nil, err
	}

	if len(result.Data) != len(clusterIDs) {
		blog.Errorf("the number of cluster obtained is inconsistent with the param, bizID: %d, ids: %#v, err: %v, "+
			"rid: %s", bizID, clusterIDs, err, kit.Rid)
		return nil, errors.New("the clusterIDs must all be under the given business")
	}
	return result.Data, nil
}

// UpdateClusterFields update cluster fields
func (s *Service) UpdateClusterFields(ctx *rest.Contexts) {

	data := new(types.UpdateClusterOption)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	clusters, err := s.getUpdateClustersInfo(ctx.Kit, bizID, data.IDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Engine.CoreAPI.CoreService().Kube().UpdateClusterFields(ctx.Kit.Ctx, ctx.Kit.Header, bizID, data)
		if err != nil {
			blog.Errorf("create cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		// for audit log.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		audit := auditlog.NewKubeAudit(s.Engine.CoreAPI.CoreService())

		for _, cluster := range clusters {
			updateFields, err := mapstr.Struct2Map(data.Data)
			if err != nil {
				blog.Errorf("update fields convert failed, data: %+v, err: %v, rid: %s", data.Data,
					err, ctx.Kit.Rid)
				return err
			}

			generateAuditParameter.WithUpdateFields(updateFields)
			auditLog, err := audit.GenerateClusterAuditLog(generateAuditParameter, []types.Cluster{cluster})
			if err != nil {
				blog.Errorf("generate audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}

			err = audit.SaveAuditLog(ctx.Kit, auditLog...)
			if err != nil {
				blog.Errorf("save audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
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

// CreateCluster create a cluster
func (s *Service) CreateCluster(ctx *rest.Contexts) {
	data := new(types.Cluster)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if err := data.ValidateCreate(); err.ErrCode != 0 {
		blog.Errorf("validate create kube cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	var id int64
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		id, err = s.Logics.KubeOperation().CreateCluster(ctx.Kit, data, bizID)
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

	ctx.RespEntity(metadata.RspID{ID: id})
}

// DeleteCluster delete cluster.
func (s *Service) DeleteCluster(ctx *rest.Contexts) {

	option := new(types.DeleteClusterOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
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
		err = s.Logics.KubeOperation().DeleteCluster(ctx.Kit, bizID, option)
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
