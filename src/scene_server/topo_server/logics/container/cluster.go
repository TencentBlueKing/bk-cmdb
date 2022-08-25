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
	BatchDeleteNode(kit *rest.Kit, bizID int64, option *types.ArrangeDeleteNodeOption, bkSupplierAccount string) error
	SearchCluster(kit *rest.Kit, input *metadata.QueryCondition) (*types.ResponseCluster, error)
	BatchCreateNode(kit *rest.Kit, data *types.CreateNodesReq, bizID int64, bkSupplierAccount string) ([]int64, error)
	SearchNode(kit *rest.Kit, input *metadata.QueryCondition) (*types.ResponseNode, error)
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
func (b *kube) BatchDeleteNode(kit *rest.Kit, bizID int64, option *types.ArrangeDeleteNodeOption,
	bkSupplierAccount string) error {

	// 1、检查是否存在这些node，必须都存在才能删除否则返回报错
	cond := make([]map[string]interface{}, 0)
	if option.Flag {
		for clusterUid, names := range option.Option {
			cond = append(cond, map[string]interface{}{
				types.ClusterUIDField: clusterUid,
				types.NodeField:       map[string]interface{}{common.BKDBIN: names},
			})
		}
	} else {
		for clusterID, ids := range option.Option {
			cond = append(cond, map[string]interface{}{
				types.BKClusterIDFiled: clusterID,
				types.BKNodeIDField:    map[string]interface{}{common.BKDBIN: ids},
			})
		}
	}

	// 查找是否有pod
	counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameBaseCluster, cond)
	if err != nil {
		blog.Errorf("count nodes failed, cond: %#v, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed)
	}

	var podNum int64
	for _, count := range counts {
		podNum += count
	}

	if podNum > 0 {
		blog.Errorf("count nodes failed, cond: %#v, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errors.New("no pods can exist under the node"))
	}

	// 3、进行批量删除node
	if err := b.clientSet.CoreService().Container().BatchDeleteNode(kit.Ctx, kit.Header, bizID, option); err != nil {
		blog.Errorf("delete cluster failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return err
	}
	return nil
}

// DeleteCluster delete cluster instance.
func (b *kube) DeleteCluster(kit *rest.Kit, bizID int64, option *types.DeleteClusterOption,
	bkSupplierAccount string) error {

	cond := make([]map[string]interface{}, 0)
	if len(option.Uids) > 0 {
		cond = []map[string]interface{}{
			{
				types.ClusterUIDField: map[string]interface{}{
					common.BKDBIN: option.Uids,
				},
				common.BKAppIDField:   bizID,
				common.BKOwnerIDField: bkSupplierAccount,
			},
		}
	}

	if len(option.IDs) > 0 {
		cond = []map[string]interface{}{
			{
				types.BKIDField: map[string]interface{}{
					common.BKDBIN: option.IDs,
				},
				common.BKAppIDField:   bizID,
				common.BKOwnerIDField: bkSupplierAccount,
			},
		}
	}

	counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		types.BKTableNameBaseCluster, cond)
	if err != nil {
		blog.Errorf("count cluster failed, cond: %#v, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed)
	}

	if len(option.Uids) > 0 && int(counts[0]) != len(option.Uids) {
		blog.Errorf("the number of instances obtained is inconsistent with the param, bizID: %d, uid: %#v, "+
			"err: %v, rid: %s", bizID, option.Uids, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed, "uid")
	}

	if len(option.IDs) > 0 && int(counts[0]) != len(option.IDs) {
		blog.Errorf("the number of instances obtained is inconsistent with the param, bizID: %d, ids: %#v, err: %v,"+
			" rid: %s", option.IDs, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed, "id")
	}

	// whether the associated resources under the cluster have been deleted.
	// Such as namespace, node, deployment, pod, container.
	exsit, cErr := b.isExsitKubeResource(kit, option, bizID, bkSupplierAccount)
	if cErr != nil {
		blog.Errorf("failed to obtain resources under the cluster, bizID: %d, option: %+v, err: %v, rid: %s",
			bizID, option, cErr, kit.Rid)
		return cErr
	}
	if exsit {
		blog.Errorf("the associated resources under the deleted cluster have not been deleted, ids: %+v, uids: %+v, "+
			"rid: %s", option.IDs, option.Uids, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrTopoInstDeleteFailed, "there are resources that have not been deleted")
	}

	if err := b.clientSet.CoreService().Container().DeleteCluster(kit.Ctx, kit.Header, bizID, option); err != nil {
		blog.Errorf("delete cluster failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return err
	}
	return nil
}

func (b *kube) isExsitKubeResource(kit *rest.Kit, option *types.DeleteClusterOption, bizID int64,
	bkSupplierAccount string) (bool, error) {

	var wg sync.WaitGroup
	var firstErr ccErr.CCErrorCoder

	workLoads := types.GetWorkLoadTables()
	tables := []string{
		types.BKTableNameBaseNamespace,
		types.BKTableNameBaseNode,
		types.BKTableNameBasePod,
		types.BKTableNameBaseContainer,
	}
	count := make([]int64, 0)
	tables = append(tables, workLoads...)
	for _, table := range tables {
		wg.Add(1)
		go func(table string, bizID int64, bkSupplierAccount string) {
			defer func() {
				// one search gcoroutine done.
				wg.Done()
			}()

			filter := make([]map[string]interface{}, 0)
			if len(option.Uids) > 0 {
				filter = []map[string]interface{}{
					{
						types.ClusterUIDField: option.Uids,
						common.BKAppIDField:   bizID,
						common.BKOwnerIDField: bkSupplierAccount,
					},
				}
			}
			if len(option.IDs) > 0 {
				filter = []map[string]interface{}{
					{
						types.ClusterUIDField: option.IDs,
						common.BKAppIDField:   bizID,
						common.BKOwnerIDField: bkSupplierAccount,
					},
				}
			}
			kit.SupplierAccount = bkSupplierAccount
			counts, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
				table, filter)
			if err != nil {
				blog.Errorf("count cluster failed, cond: %#v, err: %v, rid: %s", filter, err, kit.Rid)
				firstErr = err
				return
			}
			if counts[0] > 0 {
				blog.Errorf("count cluster failed, cond: %#v, err: %v, rid: %s", filter, err, kit.Rid)
				firstErr = err
				return
			}
			count = append(count, counts[0])

		}(table, bizID, bkSupplierAccount)
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

// BatchCreateNode batch create node.
func (b *kube) BatchCreateNode(kit *rest.Kit, data *types.CreateNodesReq, bizID int64, bkSupplierAccount string) (
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
		common.BKOwnerIDField: bkSupplierAccount,
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
	return result, nil
}

// CreateCluster create container  cluster
func (b *kube) CreateCluster(kit *rest.Kit, data *types.ClusterBaseFields, bizID int64, bkSupplierAccount string) (
	int64, error) {

	cond := map[string]interface{}{common.BKDBOR: []map[string]interface{}{
		{
			common.BKFieldName:    *data.Name,
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: bkSupplierAccount,
		},
		{
			common.BKFieldName:    *data.Uid,
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: bkSupplierAccount,
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

	return result.ID, nil
}

// SearchCluster search clusters
func (b *kube) SearchCluster(kit *rest.Kit, input *metadata.QueryCondition) (*types.ResponseCluster, error) {

	result, err := b.clientSet.CoreService().Container().SearchCluster(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("search cluster failed, input: %#v, err: %v, rid: %s", input, err, kit.Rid)
		return nil, err
	}

	return result, nil
}

// SearchNode search node by condition
func (b *kube) SearchNode(kit *rest.Kit, input *metadata.QueryCondition) (*types.ResponseNode, error) {

	result, err := b.clientSet.CoreService().Container().SearchNode(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("search cluster failed, input: %#v, err: %v, rid: %s", input, err, kit.Rid)
		return nil, err
	}

	return result, nil
}
