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

package container

import (
	"errors"
	"sync"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// ClusterOperationInterface container cluster operation methods
type ClusterOperationInterface interface {
	CreateCluster(kit *rest.Kit, data *types.ClusterBaseFields, bizID int64, bkSupplierAccount string) (int64, error)
	DeleteCluster(kit *rest.Kit, bizID int64, option *types.DeleteClusterOption, bkSupplierAccount string) error
	BatchDeleteNode(kit *rest.Kit, bizID int64, option *types.BatchDeleteNodeOption, bkSupplierAccount string) error
	BatchCreateNode(kit *rest.Kit, data *types.CreateNodesOption, bizID int64, supplierAccount string) ([]int64, error)
	BatchCreatePod(kit *rest.Kit, data *types.CreatePodsOption, bizID int64) ([]int64, error)
	SetProxy(inst ClusterOperationInterface)
}

// NewClusterOperation create a business instance
func NewClusterOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) ClusterOperationInterface {
	return &kube{
		clientSet:   client,
		authManager: authManager,
	}
}

type kube struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	cluster     ClusterOperationInterface
}

// SetProxy  SetProxy
func (b *kube) SetProxy(cluster ClusterOperationInterface) {
	b.cluster = cluster
}

// BatchDeleteNode  batch delete nodes.
func (b *kube) BatchDeleteNode(kit *rest.Kit, bizID int64, option *types.BatchDeleteNodeOption,
	supplierAccount string) error {

	// 1、check whether these nodes exist, they must all exist before they can be deleted,
	// otherwise an error will be returned.
	podCond := make([]map[string]interface{}, 0)
	num := 0
	var nodeCount int
	nodes := make([]types.Node, 0)
	for _, node := range option.Nodes {
		cond := map[string]interface{}{
			types.BKBizIDField:     bizID,
			types.BKClusterIDFiled: node.ClusterID,
			types.BKIDField:        map[string]interface{}{common.BKDBIN: node.IDs},
			common.BKOwnerIDField:  supplierAccount,
		}
		podCond = append(podCond, map[string]interface{}{
			types.BKBizIDField:     bizID,
			types.BKClusterIDFiled: node.ClusterID,
			types.BKNodeIDField:    map[string]interface{}{common.BKDBIN: node.IDs},
			common.BKOwnerIDField:  supplierAccount,
		})
		num += len(node.IDs)

		query := &metadata.QueryCondition{
			Condition: cond,
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
		}
		result, err := b.clientSet.CoreService().Container().SearchNode(kit.Ctx, kit.Header, query)
		if err != nil {
			blog.Errorf("search node failed, filter: %+v, err: %v, rid: %s", cond, err, kit.Rid)
			return err
		}
		nodeCount += len(result.Data)
		nodes = append(nodes, result.Data...)
	}

	if nodeCount != num {
		blog.Errorf("count nodes failed, option: %#v, err: %v, rid: %s", option, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed)
	}

	// 2、check if there is a pod on the node.
	counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameBasePod, podCond)
	if err != nil {
		blog.Errorf("count nodes failed, cond: %#v, err: %v, rid: %s", podCond, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed)
	}

	var podNum int64
	for _, count := range counts {
		podNum += count
	}

	if podNum > 0 {
		blog.Errorf("count nodes failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errors.New("no pods can exist under the node"))
	}

	// 3、batch delete nodes
	if err := b.clientSet.CoreService().Container().BatchDeleteNode(kit.Ctx, kit.Header, bizID, option); err != nil {
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

// DeleteCluster delete cluster instance.
func (b *kube) DeleteCluster(kit *rest.Kit, bizID int64, option *types.DeleteClusterOption,
	supplierAccount string) error {

	cond := make(map[string]interface{})

	if len(option.IDs) > 0 {
		cond = map[string]interface{}{

			types.BKIDField: map[string]interface{}{
				common.BKDBIN: option.IDs,
			},
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: supplierAccount,
		}
	}

	input := &metadata.QueryCondition{
		Condition: cond,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	result, err := b.clientSet.CoreService().Container().SearchCluster(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("search cluster failed, input: %#v, err: %v, rid: %s", input, err, kit.Rid)
		return err
	}

	if len(option.IDs) > 0 && len(result.Data) != len(option.IDs) {
		blog.Errorf("the number of instances obtained is inconsistent with the param, bizID: %d, ids: %#v, err: %v,"+
			" rid: %s", option.IDs, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed, "id")
	}

	// whether the associated resources under the cluster have been deleted.
	// such as namespace, node, deployment, pod.
	exist, cErr := b.isExistKubeResource(kit, option, bizID, supplierAccount)
	if cErr != nil {
		blog.Errorf("failed to obtain resources under the cluster, bizID: %d, option: %+v, err: %v, rid: %s",
			bizID, option, cErr, kit.Rid)
		return cErr
	}
	if exist {
		blog.Errorf("the associated resources under the deleted cluster have not been deleted, ids: %+v, rid: %s",
			option.IDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed, "there are resources that have not been deleted")
	}

	if err := b.clientSet.CoreService().Container().DeleteCluster(kit.Ctx, kit.Header, bizID, option); err != nil {
		blog.Errorf("delete cluster failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return err
	}
	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	audit := auditlog.NewKubeAudit(b.clientSet.CoreService())
	auditLog, err := audit.GenerateClusterAuditLog(generateAuditParameter, result.Data)
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

func (b *kube) isExistKubeResource(kit *rest.Kit, option *types.DeleteClusterOption, bizID int64,
	supplierAccount string) (bool, error) {

	var wg sync.WaitGroup
	var firstErr ccErr.CCErrorCoder

	workLoads := types.GetWorkLoadTables()
	tables := []string{
		types.BKTableNameBaseNamespace,
		types.BKTableNameBaseNode,
		types.BKTableNameBasePod,
	}
	count := make([]int64, 0)
	tables = append(tables, workLoads...)
	for _, table := range tables {
		wg.Add(1)
		go func(table string, bizID int64, supplierAccount string) {
			defer func() {
				// one search goroutine done.
				wg.Done()
			}()

			filter := make([]map[string]interface{}, 0)
			if len(option.IDs) > 0 {
				filter = []map[string]interface{}{
					{
						types.BKClusterIDFiled: map[string]interface{}{
							common.BKDBIN: option.IDs,
						},
						common.BKAppIDField:   bizID,
						common.BKOwnerIDField: supplierAccount,
					},
				}
			}
			counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
				table, filter)
			if err != nil {
				blog.Errorf("count cluster failed, cond: %#v, err: %v, rid: %s", filter, err, kit.Rid)
				firstErr = err
				return
			}
			count = append(count, counts[0])

		}(table, bizID, supplierAccount)
	}
	wg.Wait()
	if firstErr != nil {
		return false, firstErr
	}
	for _, num := range count {
		if num > 0 {
			return true, nil
		}
	}
	return false, nil
}

// BatchCreatePod batch create pod.
func (b *kube) BatchCreatePod(kit *rest.Kit, data *types.CreatePodsOption, bizID int64) (
	[]int64, error) {

	filters := make([]map[string]interface{}, 0)

	for _, pod := range data.Pods {
		filter := map[string]interface{}{
			common.BKOwnerIDField: kit.SupplierAccount,
			types.BKBizIDField:    bizID,
		}
		filter[types.BKClusterIDFiled] = *pod.Spec.ClusterID
		filter[types.BKNamespaceIDField] = *pod.Spec.NamespaceID
		filter[types.BKNodeIDField] = *pod.Spec.NodeID
		filter[types.KubeNameField] = *pod.Name
		filter[types.RefKindField] = *pod.Spec.WorkloadKind
		filter[types.RefIDField] = *pod.Spec.WorkloadID
		filters = append(filters, filter)
	}

	counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameBasePod, filters)
	if err != nil {
		blog.Errorf("count cluster failed, cond: %#v, err: %v, rid: %s", filters, err, kit.Rid)
		return nil, err
	}

	var podNum int64
	for _, count := range counts {
		podNum += count
	}
	if podNum > 0 {
		blog.Errorf("some pod already exists and the creation fails, rid: %s", kit.Rid)
		return nil, errors.New("some pod already exists and the creation fails")
	}

	//2、create pods and containers.
	result, err := b.clientSet.CoreService().Container().BatchCreatePod(kit.Ctx, kit.Header, bizID, data)
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
func (b *kube) BatchCreateNode(kit *rest.Kit, data *types.CreateNodesOption, bizID int64, supplierAccount string) (
	[]int64, error) {

	names := make([]string, 0)

	for _, node := range data.Nodes {
		names = append(names, *node.Name)
	}

	cond := map[string]interface{}{
		common.BKFieldName: map[string]interface{}{
			common.BKDBIN: names,
		},
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: supplierAccount,
	}

	counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameBaseNode, []map[string]interface{}{cond})
	if err != nil {
		blog.Errorf("count cluster failed, cond: %#v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrTopoInstCreateFailed, "cluster name or uid has been created")
	}
	if counts[0] > 0 {
		blog.Errorf("duplicate node name exists, num: %d, err: %v, rid: %s", counts[0], err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrTopoInstCreateFailed, "duplicate node name has been created")
	}

	result, err := b.clientSet.CoreService().Container().BatchCreateNode(kit.Ctx, kit.Header, bizID, data)
	if err != nil {
		blog.Errorf("create nodes failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}
	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewKubeAudit(b.clientSet.CoreService())
	auditLog, err := audit.GenerateNodeAuditLog(generateAuditParameter, result.Info)

	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrAuditSaveLogFailed)
	}
	ids := make([]int64, 0)
	for _, node := range result.Info {
		ids = append(ids, node.ID)
	}
	return ids, nil
}

// CreateCluster create container  cluster
func (b *kube) CreateCluster(kit *rest.Kit, data *types.ClusterBaseFields, bizID int64, supplierAccount string) (
	int64, error) {

	cond := map[string]interface{}{common.BKDBOR: []map[string]interface{}{
		{
			common.BKFieldName:    *data.Name,
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: supplierAccount,
		},
		{
			common.BKFieldName:    *data.Uid,
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: supplierAccount,
		},
	},
	}
	// kit.SupplierAccount = bkSupplierAccount
	counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameBaseCluster, []map[string]interface{}{cond})
	if err != nil {
		blog.Errorf("count cluster failed, cond: %#v, err: %v, rid: %s", cond, err, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrTopoInstCreateFailed, "cluster name or uid has been created")
	}
	if counts[0] > 0 {
		blog.Errorf("cluster name or uid has been created, num: %d, err: %v, rid: %s", counts[0], err, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrTopoInstCreateFailed, "cluster name or uid has been created")
	}

	result, err := b.clientSet.CoreService().Container().CreateCluster(kit.Ctx, kit.Header, bizID, data)
	if err != nil {
		blog.Errorf("create business failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
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
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return result.Info.ID, nil
}
