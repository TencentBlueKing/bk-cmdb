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
	"fmt"
	"sync"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// KubeOperationInterface container cluster operation methods
type KubeOperationInterface interface {
	CreateCluster(kit *rest.Kit, data *types.Cluster, bizID int64) (int64, error)
	DeleteCluster(kit *rest.Kit, bizID int64, option *types.DeleteClusterOption) error
	BatchDeleteNode(kit *rest.Kit, bizID int64, option *types.BatchDeleteNodeOption) error
	BatchCreateNode(kit *rest.Kit, data *types.CreateNodesOption, bizID int64) ([]int64, error)
	BatchCreatePod(kit *rest.Kit, data *types.CreatePodsOption) ([]int64, error)
}

// NewClusterOperation create a business instance
func NewClusterOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) KubeOperationInterface {
	return &kube{
		clientSet:   client,
		authManager: authManager,
	}
}

type kube struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	cluster     KubeOperationInterface
}

func (b *kube) getDeleteNodeInfo(kit *rest.Kit, ids []int64, bizID int64) ([]types.Node, error) {

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKFieldID: mapstr.MapStr{
				common.BKDBIN: ids,
			},
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	result, err := b.clientSet.CoreService().Kube().SearchNode(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("search node failed, filter: %+v, err: %v, rid: %s", query, err, kit.Rid)
		return nil, err
	}

	bizMap := make(map[int64]struct{})
	for _, node := range result.Data {
		bizMap[node.BizID] = struct{}{}
	}

	if len(bizMap) > 1 {
		blog.Errorf("node ids exist in different businesses, filter: %+v, rid: %s", query, kit.Rid)
		return nil, errors.New("node ids exist in different businesses")
	}

	if _, ok := bizMap[bizID]; !ok {
		blog.Errorf("node ids not in biz %d, filter: %+v, rid: %s", bizID, query, kit.Rid)
		return nil, fmt.Errorf("node ids not in biz %d", bizID)
	}

	return result.Data, nil
}

// BatchDeleteNode batch delete node.
func (b *kube) BatchDeleteNode(kit *rest.Kit, bizID int64, option *types.BatchDeleteNodeOption) error {

	nodes, err := b.getDeleteNodeInfo(kit, option.IDs, bizID)
	if err != nil {
		return err
	}
	// 1、check whether these nodes exist, they must all exist before they can be deleted,
	// otherwise an error will be returned.
	podCond := []map[string]interface{}{
		{
			types.BKBizIDField:  bizID,
			types.BKNodeIDField: map[string]interface{}{common.BKDBIN: option.IDs},
		},
	}
	// 2、check if there is a pod on the node.
	counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameBasePod, podCond)
	if err != nil {
		blog.Errorf("count nodes failed, cond: %#v, err: %v, rid: %s", podCond, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed)
	}

	if counts[0] > 0 {
		blog.Errorf("count nodes failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errors.New("no pods can exist under the node"))
	}

	// 3、batch delete nodes
	if err := b.clientSet.CoreService().Kube().BatchDeleteNode(kit.Ctx, kit.Header, bizID, option); err != nil {
		blog.Errorf("delete node failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return err
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	audit := auditlog.NewKubeAudit(b.clientSet.CoreService())
	auditLog, err := audit.GenerateNodeAuditLog(generateAuditParameter, nodes)
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

func (b *kube) getDeleteClusterInfo(kit *rest.Kit, ids []int64, bizID int64) ([]types.Cluster, error) {

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKFieldID: mapstr.MapStr{
				common.BKDBIN: ids,
			},
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	result, err := b.clientSet.CoreService().Kube().SearchCluster(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("search cluster failed, filter: %+v, err: %v, rid: %s", query, err, kit.Rid)
		return nil, err
	}

	bizMap := make(map[int64]struct{})
	for _, node := range result.Data {
		bizMap[node.BizID] = struct{}{}
	}

	if len(bizMap) != 1 {
		blog.Errorf("cluster ids exist in different businesses, filter: %+v, rid: %s", query, kit.Rid)
		return nil, errors.New("node ids exist in different businesses")
	}

	if _, ok := bizMap[bizID]; !ok {
		blog.Errorf("cluster ids not in biz %d, filter: %+v, rid: %s", bizID, query, kit.Rid)
		return nil, fmt.Errorf("cluster ids not in biz %d", bizID)
	}

	return result.Data, nil
}

// DeleteCluster delete cluster.
func (b *kube) DeleteCluster(kit *rest.Kit, bizID int64, option *types.DeleteClusterOption) error {

	clusters, err := b.getDeleteClusterInfo(kit, option.IDs, bizID)
	if err != nil {
		return err
	}
	// whether the associated resources under the cluster have been deleted. such as namespace, node, deployment, pod.
	exist, cErr := b.isExistKubeResourceUnderCluster(kit, option, bizID)
	if cErr != nil {
		blog.Errorf("failed to obtain resources under the cluster, bizID: %d, cluster IDs: %+v, err: %v, rid: %s",
			bizID, option.IDs, cErr, kit.Rid)
		return cErr
	}
	if exist {
		blog.Errorf("the associated resources under the deleted cluster haven't been deleted, clusterIDs: %+v, rid: %s",
			option.IDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid)
	}

	if err := b.clientSet.CoreService().Kube().DeleteCluster(kit.Ctx, kit.Header, bizID, option); err != nil {
		blog.Errorf("delete cluster failed, cluster IDs: %#v, err: %v, rid: %s", option.IDs, err, kit.Rid)
		return err
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	audit := auditlog.NewKubeAudit(b.clientSet.CoreService())
	auditLog, err := audit.GenerateClusterAuditLog(generateAuditParameter, clusters)
	if err != nil {
		blog.Errorf("generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return nil
}

func (b *kube) isExistKubeResourceUnderCluster(kit *rest.Kit, option *types.DeleteClusterOption, bizID int64) (
	bool, error) {

	if len(option.IDs) == 0 {
		return false, errors.New("ids must be set")
	}

	var (
		wg       sync.WaitGroup
		firstErr ccErr.CCErrorCoder
	)

	workLoads := types.GetWorkLoadTables()
	tables := []string{types.BKTableNameBaseNamespace, types.BKTableNameBaseNode, types.BKTableNameBasePod}
	tables = append(tables, workLoads...)

	filter := make([]map[string]interface{}, 0)
	filter = []map[string]interface{}{
		{
			types.BKClusterIDFiled: map[string]interface{}{common.BKDBIN: option.IDs},
			common.BKAppIDField:    bizID,
		},
	}

	for _, table := range tables {
		wg.Add(1)
		go func(table string, bizID int64) {
			defer func() {
				wg.Done()
			}()

			counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, table, filter)
			if err != nil {
				blog.Errorf("count resource failed, cond: %#v, err: %v, rid: %s", filter, err, kit.Rid)
				firstErr = err
				return
			}
			if counts[0] > 0 {
				blog.Errorf("there are resources under the cluster that cannot be deleted, bizID %d, filter: %+v, "+
					"table: %s, rid: %s", bizID, table, kit.Rid)
				firstErr = kit.CCError.CCErrorf(common.CCErrCommParamsInvalid)
				return
			}

		}(table, bizID)
	}
	wg.Wait()
	if firstErr != nil {
		return false, firstErr
	}

	return false, nil
}

// BatchCreatePod batch create pod.
func (b *kube) BatchCreatePod(kit *rest.Kit, data *types.CreatePodsOption) ([]int64, error) {

	// create pods and containers.
	result, err := b.clientSet.CoreService().Kube().BatchCreatePod(kit.Ctx, kit.Header, data)
	if err != nil {
		blog.Errorf("create pod failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewKubeAudit(b.clientSet.CoreService())

	podIDs := make([]int64, 0)
	for _, pod := range result {
		podIDs = append(podIDs, pod.ID)
	}
	auditLog, err := audit.GeneratePodAuditLog(generateAuditParameter, result)
	if err != nil {
		blog.Errorf("create cluster, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return podIDs, nil
}

// BatchCreateNode batch create node.
func (b *kube) BatchCreateNode(kit *rest.Kit, data *types.CreateNodesOption, bizID int64) ([]int64, error) {

	conds := make([]map[string]interface{}, 0)
	for _, node := range data.Nodes {
		conds = append(conds, map[string]interface{}{
			common.BKFieldName:     *node.Name,
			types.BKClusterIDFiled: node.ClusterID,
		})
	}

	cond := map[string]interface{}{
		common.BKDBOR:       conds,
		common.BKAppIDField: bizID,
	}

	counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameBaseNode, []map[string]interface{}{cond})
	if err != nil {
		blog.Errorf("count cluster failed, cond: %#v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoInstCreateFailed)
	}

	if counts[0] > 0 {
		blog.Errorf("duplicate node name exists, num: %d, err: %v, rid: %s", counts[0], err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "duplicate node name has been created")
	}

	result, err := b.clientSet.CoreService().Kube().BatchCreateNode(kit.Ctx, kit.Header, bizID, data)
	if err != nil {
		blog.Errorf("create nodes failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}
	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewKubeAudit(b.clientSet.CoreService())
	auditLog, err := audit.GenerateNodeAuditLog(generateAuditParameter, result.Info)

	if err != nil {
		blog.Errorf(" creat nodes, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create nodes, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrAuditSaveLogFailed)
	}

	ids := make([]int64, 0)
	for _, node := range result.Info {
		ids = append(ids, node.ID)
	}
	return ids, nil
}

// CreateCluster create container cluster
func (b *kube) CreateCluster(kit *rest.Kit, data *types.Cluster, bizID int64) (int64, error) {

	if data.Name == nil || data.Uid == nil {
		return 0, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "cluster name or uid must be set")
	}

	filterName := map[string]interface{}{
		common.BKFieldName:  *data.Name,
		common.BKAppIDField: bizID,
	}

	filterUid := map[string]interface{}{
		common.BKFieldName:  *data.Uid,
		common.BKAppIDField: bizID,
	}

	cond := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{filterName, filterUid},
	}

	counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameBaseCluster, []map[string]interface{}{cond})
	if err != nil {
		blog.Errorf("count cluster failed, cond: %#v, err: %v, rid: %s", cond, err, kit.Rid)
		return 0, kit.CCError.CCError(common.CCErrTopoInstCreateFailed)
	}
	if counts[0] > 0 {
		blog.Errorf("cluster name or uid has been created, num: %d, err: %v, rid: %s", counts[0], err, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "cluster name or uid has been created")
	}

	result, err := b.clientSet.CoreService().Kube().CreateCluster(kit.Ctx, kit.Header, bizID, data)
	if err != nil {
		blog.Errorf("create cluster failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
		return 0, err
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewKubeAudit(b.clientSet.CoreService())
	auditLog, err := audit.GenerateClusterAuditLog(generateAuditParameter, []types.Cluster{*result.Info})
	if err != nil {
		blog.Errorf("create cluster, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create cluster, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return result.Info.ID, nil
}
