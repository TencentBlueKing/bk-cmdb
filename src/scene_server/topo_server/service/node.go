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
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// FindNodePathForHost find node path for host
func (s *Service) FindNodePathForHost(ctx *rest.Contexts) {
	req := types.HostPathReq{}
	if err := ctx.DecodeInto(&req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	relation, err := s.getHostNodeRelation(ctx.Kit, req.HostIDs)
	if err != nil {
		blog.Errorf("get host and node relation failed, ids: %v, err: %v, rid: %s", req.HostIDs, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizIDWithName, err := s.getBizIDWithName(ctx.Kit, relation.BizIDs)
	if err != nil {
		blog.Errorf("get bizID with name failed, bizIDs: %v, err: %v, rid: %s", relation.BizIDs, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	hostsPath := make([]types.HostNodePath, len(req.HostIDs))
	for outerIdx, hostID := range req.HostIDs {
		nodes := relation.HostWithNode[hostID]
		paths := make([]types.NodePath, 0)
		uniqueMap := make(map[string]struct{})
		for _, node := range nodes {
			clusterID := relation.NodeIDWithClusterID[node.ID]
			bizID := relation.NodeIDWithBizID[node.ID]

			unique := strconv.FormatInt(bizID, 10) + ":" + strconv.FormatInt(clusterID, 10)
			if _, ok := uniqueMap[unique]; ok {
				continue
			}
			uniqueMap[unique] = struct{}{}

			path := types.NodePath{
				BizID:       bizID,
				BizName:     bizIDWithName[bizID],
				ClusterID:   clusterID,
				ClusterName: relation.ClusterIDWithName[clusterID],
			}
			paths = append(paths, path)
		}

		hostsPath[outerIdx] = types.HostNodePath{
			HostID: hostID,
			Paths:  paths,
		}
	}

	ctx.RespEntity(types.HostPathData{
		Info: hostsPath,
	})
}

func (s *Service) getHostNodeRelation(kit *rest.Kit, hostIDs []int64) (*types.HostNodeRelation, error) {
	cond := mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: hostIDs}}
	fields := []string{
		common.BKFieldID, common.BKAppIDField, types.BKClusterIDFiled, common.BKHostIDField,
	}
	query := &metadata.QueryCondition{
		Condition: cond,
		Fields:    fields,
	}

	resp, ccErr := s.Engine.CoreAPI.CoreService().Kube().SearchNode(kit.Ctx, kit.Header, query)
	if ccErr != nil {
		blog.Errorf("find node failed, cond: %v, err: %v, rid: %s", query, ccErr, kit.Rid)
		return nil, ccErr
	}

	bizIDs := make([]int64, 0)
	hostWithNode := make(map[int64][]types.Node)
	nodeIDWithBizID := make(map[int64]int64)
	nodeIDWithClusterID := make(map[int64]int64)
	clusterIDs := make([]int64, 0)
	for _, node := range resp.Data {
		bizIDs = append(bizIDs, node.BizID)
		nodeIDWithBizID[node.ID] = node.BizID
		hostWithNode[node.HostID] = append(hostWithNode[node.HostID], node)
		nodeIDWithClusterID[node.ID] = node.ClusterID
		clusterIDs = append(clusterIDs, node.ClusterID)
	}

	clusterIDWithName, err := s.getClusterIDWithName(kit, clusterIDs)
	if err != nil {
		blog.Errorf("get cluster id with name failed, clusterIDs: %v, err: %v, rid: %s", clusterIDs, err, kit.Rid)
		return nil, err
	}

	return &types.HostNodeRelation{
		BizIDs:              bizIDs,
		HostWithNode:        hostWithNode,
		NodeIDWithBizID:     nodeIDWithBizID,
		NodeIDWithClusterID: nodeIDWithClusterID,
		ClusterIDWithName:   clusterIDWithName,
	}, nil
}

func (s *Service) getClusterIDWithName(kit *rest.Kit, clusterIDs []int64) (map[int64]string, error) {
	cond := mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: clusterIDs}}
	fields := []string{common.BKFieldID, common.BKFieldName}
	query := &metadata.QueryCondition{
		Condition: cond,
		Fields:    fields,
	}

	resp, ccErr := s.Engine.CoreAPI.CoreService().Kube().SearchCluster(kit.Ctx, kit.Header, query)
	if ccErr != nil {
		blog.Errorf("find node failed, cond: %v, err: %v, rid: %s", query, ccErr, kit.Rid)
		return nil, ccErr
	}

	idWithName := make(map[int64]string)
	for _, cluster := range resp.Data {
		if cluster.Name == nil {
			blog.Errorf("get node attribute failed, attr: %s, cluster: %v, rid: %s", common.BKFieldName, cluster,
				kit.Rid)
			return nil, fmt.Errorf("get node attribute failed, attr: %s", common.BKFieldName)
		}

		idWithName[cluster.ID] = *cluster.Name
	}

	return idWithName, nil
}

func (s *Service) getBizIDWithName(kit *rest.Kit, bizIDs []int64) (map[int64]string, error) {
	query := &metadata.QueryCondition{
		Fields: []string{
			common.BKAppIDField,
			common.BKAppNameField,
		},
		Condition: mapstr.MapStr{
			common.BKDataStatusField: mapstr.MapStr{common.BKDBNE: bizIDs},
		},
		DisableCounter: true,
	}
	_, instItems, err := s.Logics.BusinessOperation().FindBiz(kit, query)
	if err != nil {
		blog.Errorf("find business failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	bizIDWithName := make(map[int64]string, len(instItems))
	for _, biz := range instItems {
		bizID, err := biz.Int64(common.BKAppIDField)
		if err != nil {
			blog.Errorf("the biz is invalid, data: %v, err: %v, rid: %s", biz, err, kit.Rid)
			return nil, err
		}

		name, err := biz.String(common.BKAppNameField)
		if err != nil {
			blog.Errorf("the biz is invalid, data: %v, err: %v, rid: %s", biz, err, kit.Rid)
			return nil, err
		}

		bizIDWithName[bizID] = name
	}

	return bizIDWithName, nil
}

// BatchDeleteNode delete nodes.
func (s *Service) BatchDeleteNode(ctx *rest.Contexts) {
	option := new(types.BatchDeleteNodeOption)
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
		err = s.Logics.KubeOperation().BatchDeleteNode(ctx.Kit, bizID, option)
		if err != nil {
			blog.Errorf("delete node failed, biz: %d, option: %+v, err: %v, rid: %s", bizID, option, err, ctx.Kit.Rid)
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

// BatchCreateNode batch create nodes.
func (s *Service) BatchCreateNode(ctx *rest.Contexts) {

	data := new(types.CreateNodesOption)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.ValidateCreate(); err != nil {
		blog.Errorf("batch create nodes param verification failed, data: %+v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	var ids []int64
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ids, err = s.Logics.KubeOperation().BatchCreateNode(ctx.Kit, data, bizID)
		if err != nil {
			blog.Errorf("create node failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(ids)
}

// SearchNodes query nodes based on user-specified criteria
func (s *Service) SearchNodes(ctx *rest.Contexts) {

	searchCond := new(types.QueryNodeOption)
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
			blog.Errorf("parse biz failed, filter: %+v, err: %v, rid: %s", searchCond.Filter, rawErr, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey))
			return
		}
		filter = cond
	}

	// regardless of whether there is bk_biz_id or supplier_account in the condition,
	// it is uniformly replaced with bk_biz_id in url and supplier_account in kit.
	filter[types.BKBizIDField] = bizID
	filter[types.BKSupplierAccountField] = ctx.Kit.SupplierAccount
	if searchCond.HostID != 0 {
		filter[common.BKHostIDField] = searchCond.HostID
	}

	if searchCond.ClusterID != 0 {
		filter[types.BKClusterIDFiled] = searchCond.ClusterID
	}

	// count biz in cluster enable count is set
	if searchCond.Page.EnableCount {
		filter := []map[string]interface{}{filter}
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBaseNode, filter)
		if err != nil {
			blog.Errorf("count node failed, cond: %+v, err: %v, rid: %s", filter, err, ctx.Kit.Rid)
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
	result, err := s.Engine.CoreAPI.CoreService().Kube().SearchNode(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("search node failed, filter: %+v, err: %v, rid: %s", filter, err, ctx.Kit.Rid)
		return
	}
	ctx.RespEntityWithCount(0, result.Data)
}

// UpdateNodeFields update the node field.
func (s *Service) UpdateNodeFields(ctx *rest.Contexts) {

	data := new(types.UpdateNodeOption)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if err := data.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}

	nodeIDs, nodeIDUpdateFields := make([]int64, 0), make(map[int64]types.Node)

	for _, node := range data.Nodes {
		nodeIDs = append(nodeIDs, node.NodeIDs...)
		for _, id := range node.NodeIDs {
			nodeIDUpdateFields[id] = node.Data
		}
	}

	// duplicate nodeIDs are not allowed
	nodeIDMap := make(map[int64]struct{})
	for _, nodeID := range nodeIDs {
		if _, ok := nodeIDMap[nodeID]; ok {
			blog.Errorf("duplicate value for nodeID parameter, nodeID: %d, rid: %s", nodeID, ctx.Kit.Rid)
			ctx.RespAutoError(errors.New("duplicate value for nodeID parameter"))
			return
		}
		nodeIDMap[nodeID] = struct{}{}
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	query := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			types.BKIDField:       map[string]interface{}{common.BKDBIN: nodeIDs},
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: ctx.Kit.SupplierAccount},
		Page: metadata.BasePage{Limit: common.BKNoLimit},
	}
	result, err := s.Engine.CoreAPI.CoreService().Kube().SearchNode(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("search node failed, filter: %+v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.Engine.CoreAPI.CoreService().Kube().UpdateNodeFields(ctx.Kit.Ctx, ctx.Kit.Header,
			ctx.Kit.SupplierAccount, bizID, data); err != nil {
			blog.Errorf("create cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}

		for _, node := range result.Data {
			generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
			updateFields, err := mapstr.Struct2Map(nodeIDUpdateFields[node.ID])
			if err != nil {
				blog.Errorf("update fields convert failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}

			generateAuditParameter.WithUpdateFields(updateFields)
			audit := auditlog.NewKubeAudit(s.Engine.CoreAPI.CoreService())
			auditLog, err := audit.GenerateNodeAuditLog(generateAuditParameter, result.Data)
			if err != nil {
				blog.Errorf("update node failed, generate audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}

			err = audit.SaveAuditLog(ctx.Kit, auditLog...)
			if err != nil {
				blog.Errorf("update node failed, save audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return ctx.Kit.CCError.CCErrorf(common.CCErrAuditSaveLogFailed)
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
