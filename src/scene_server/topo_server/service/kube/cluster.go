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

package kube

import (
	"errors"

	acmeta "configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// SearchClusters query based on user-specified criteria cluster
func (s *service) SearchClusters(ctx *rest.Contexts) {
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

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeCluster, Action: acmeta.Find},
		BusinessID: searchCond.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	// compatible for shared cluster scenario
	filter, err := s.Logics.KubeOperation().GenSharedClusterListCond(ctx.Kit, searchCond.BizID, searchCond.Filter)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// get the number of clusters
	if searchCond.Page.EnableCount {
		cond := []map[string]interface{}{filter}
		counts, err := s.ClientSet.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
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
	result, err := s.ClientSet.CoreService().Kube().SearchCluster(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("search cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
	ctx.RespEntityWithCount(0, result.Data)
}

func (s *service) getUpdateClustersInfo(kit *rest.Kit, bizID int64, clusterIDs []int64) ([]types.Cluster, error) {
	cond := map[string]interface{}{
		types.BKIDField:     map[string]interface{}{common.BKDBIN: clusterIDs},
		common.BKAppIDField: bizID,
	}

	input := &metadata.QueryCondition{
		Condition: cond,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}

	result, err := s.ClientSet.CoreService().Kube().SearchCluster(kit.Ctx, kit.Header, input)
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
func (s *service) UpdateClusterFields(ctx *rest.Contexts) {
	data := new(types.UpdateClusterOption)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeCluster, Action: acmeta.Update},
		BusinessID: data.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	clusters, err := s.getUpdateClustersInfo(ctx.Kit, data.BizID, data.IDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// can not update cluster type using batch update cluster api
	if data.UpdateClusterByIDsOption.Data.Type != nil {
		data.UpdateClusterByIDsOption.Data.Type = nil
	}

	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.ClientSet.CoreService().Kube().UpdateClusterFields(ctx.Kit.Ctx, ctx.Kit.Header,
			&data.UpdateClusterByIDsOption)
		if err != nil {
			blog.Errorf("update cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}

		// for audit log.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		audit := auditlog.NewKubeAudit(s.ClientSet.CoreService())

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

// UpdateClusterType update cluster type
func (s *service) UpdateClusterType(ctx *rest.Contexts) {
	opt := new(types.UpdateClusterTypeOpt)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeCluster, Action: acmeta.Update},
		BusinessID: opt.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	// get the cluster to be updated and validate cluster type
	clusters, err := s.getUpdateClustersInfo(ctx.Kit, opt.BizID, []int64{opt.ID})
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(clusters) != 1 || clusters[0].Type == nil {
		blog.Errorf("updated cluster is invalid, opt: %+v, rid: %s", opt, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.BKIDField))
		return
	}

	if opt.Type == *clusters[0].Type {
		ctx.RespEntity(nil)
		return
	}

	if err = s.validateUpdateClusterType(ctx.Kit, opt, *clusters[0].Type); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		updateOpt := &types.UpdateClusterByIDsOption{
			IDs:  []int64{opt.ID},
			Data: types.Cluster{Type: &opt.Type},
		}

		err = s.ClientSet.CoreService().Kube().UpdateClusterFields(ctx.Kit.Ctx, ctx.Kit.Header, updateOpt)
		if err != nil {
			blog.Errorf("update cluster type failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}

		// for audit log.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		audit := auditlog.NewKubeAudit(s.ClientSet.CoreService())

		generateAuditParameter.WithUpdateFields(map[string]interface{}{types.TypeField: opt.Type})
		auditLog, err := audit.GenerateClusterAuditLog(generateAuditParameter, clusters)
		if err != nil {
			blog.Errorf("generate audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}

		err = audit.SaveAuditLog(ctx.Kit, auditLog...)
		if err != nil {
			blog.Errorf("save audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrAuditSaveLogFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// validateUpdateClusterType validate cluster type for update
func (s *service) validateUpdateClusterType(kit *rest.Kit, opt *types.UpdateClusterTypeOpt,
	preType types.ClusterType) error {

	switch opt.Type {
	case types.SharedClusterType:
		if preType != types.IndependentClusterType {
			blog.Errorf("previous cluster type %s is invalid, to update: %s, rid: %s", preType, opt.Type, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.TypeField)
		}
	case types.IndependentClusterType:
		if preType != types.SharedClusterType {
			blog.Errorf("previous cluster type %s is invalid, to update: %s, rid: %s", preType, opt.Type, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.TypeField)
		}

		// check if shared cluster has shared resources, if so, it can't be changed to independent
		filter := []map[string]interface{}{{
			types.BKClusterIDFiled: opt.ID,
			types.BKBizIDField:     mapstr.MapStr{common.BKDBNE: opt.BizID},
		}}

		counts, err := s.ClientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
			types.BKTableNameBaseNamespace, filter)
		if err != nil {
			blog.Errorf("count shared namespace failed, cond: %#v, err: %v, rid: %s", filter, err, kit.Rid)
			return err
		}

		if counts[0] > 0 {
			blog.Errorf("cluster has %d shared ns, filter: %+v, rid: %s", filter, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.BKClusterIDFiled)
		}
	default:
		blog.Errorf("cluster type %s to update is invalid, rid: %s", opt.Type, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.TypeField)
	}

	return nil
}

// CreateCluster create a cluster
func (s *service) CreateCluster(ctx *rest.Contexts) {
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

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeCluster, Action: acmeta.Create},
		BusinessID: data.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	var id int64
	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// create cluster
		result, err := s.ClientSet.CoreService().Kube().CreateCluster(ctx.Kit.Ctx, ctx.Kit.Header, data)
		if err != nil {
			blog.Errorf("create cluster failed, data: %#v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
			return err
		}

		// for audit log.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditCreate)
		audit := auditlog.NewKubeAudit(s.ClientSet.CoreService())
		auditLog, err := audit.GenerateClusterAuditLog(generateAuditParameter, []types.Cluster{*result})
		if err != nil {
			blog.Errorf("create cluster, generate audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}

		err = audit.SaveAuditLog(ctx.Kit, auditLog...)
		if err != nil {
			blog.Errorf("create cluster, save audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrAuditSaveLogFailed)
		}

		id = result.ID
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(metadata.RspID{ID: id})
}

// DeleteCluster delete cluster.
func (s *service) DeleteCluster(ctx *rest.Contexts) {
	option := new(types.DeleteClusterOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeCluster, Action: acmeta.Delete},
		BusinessID: option.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		err = s.Logics.KubeOperation().DeleteCluster(ctx.Kit, option.BizID, option)
		if err != nil {
			blog.Errorf("delete cluster failed, biz: %d, option: %+v, err: %v, rid: %s", option.BizID, option, err,
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
