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
	"fmt"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	acmeta "configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/kube/types"
)

// FindPodPath find pod path
func (s *service) FindPodPath(ctx *rest.Contexts) {
	req := new(types.PodPathOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubePod, Action: acmeta.Find},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	// compatible for shared cluster scenario
	podIDCond := filtertools.GenAtomFilter(common.BKFieldID, filter.In, req.PodIDs)
	cond, err := s.Logics.KubeOperation().GenSharedNsListCond(ctx.Kit, types.KubePod, req.BizID, podIDCond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	fields := []string{common.BKFieldID, common.BKAppIDField, types.BKClusterIDFiled, types.BKNamespaceIDField,
		types.NamespaceField, types.RefField}
	query := &metadata.QueryCondition{
		Condition: cond,
		Fields:    fields,
	}
	resp, err := s.ClientSet.CoreService().Kube().ListPod(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("find pod failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(resp.Info) == 0 {
		ctx.RespEntity(types.PodPathData{Info: []types.PodPath{}})
		return
	}

	paths, rawErr := s.buildPodPaths(ctx.Kit, req.BizID, resp.Info)
	if err != nil {
		ctx.RespAutoError(rawErr)
		return
	}

	ctx.RespEntity(types.PodPathData{
		Info: paths,
	})
}

func (s *service) buildPodPaths(kit *rest.Kit, bizID int64, pods []types.Pod) ([]types.PodPath, error) {
	paths := make([]types.PodPath, 0)
	clusterIDs := make([]int64, 0)
	allBizIDs := make([]int64, 0)

	for _, pod := range pods {
		id := pod.ID

		if pod.ClusterID == 0 {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.BKClusterIDFiled,
				kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.BKClusterIDFiled)
		}
		clusterID := pod.ClusterID
		clusterIDs = append(clusterIDs, clusterID)

		if pod.NamespaceID == 0 {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.BKNamespaceIDField,
				kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.BKNamespaceIDField)
		}
		namespaceID := pod.NamespaceID

		if pod.Namespace == "" {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.NamespaceField, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.NamespaceField)
		}
		namespace := pod.Namespace

		if pod.Ref == nil {
			blog.Errorf("pod %v ref is not set, rid: %s", pod, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.RefField)
		}
		if pod.Ref.Kind == "" {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, rid: %s", types.RefKindField, pod, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.RefKindField)
		}
		if pod.Ref.Name == "" {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, rid: %s", types.RefNameField, pod, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.RefNameField)
		}
		if pod.Ref.ID == 0 {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, rid: %s", types.RefIDField, pod, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.RefIDField)
		}
		ref := pod.Ref

		path := types.PodPath{
			BizID:        pod.BizID,
			ClusterID:    clusterID,
			NamespaceID:  namespaceID,
			Namespace:    namespace,
			Kind:         ref.Kind,
			WorkloadID:   ref.ID,
			WorkloadName: ref.Name,
			PodID:        id,
		}
		paths = append(paths, path)
		allBizIDs = append(allBizIDs, pod.BizID)
	}

	return s.combinePodPath(kit, bizID, allBizIDs, clusterIDs, paths)
}

func (s *service) combinePodPath(kit *rest.Kit, bizID int64, allBizIDs, clusterIDs []int64, paths []types.PodPath) (
	[]types.PodPath, error) {

	// get cluster info, including biz id for shared cluster scenario
	clusterCond := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: clusterIDs}},
		Fields:         []string{common.BKFieldID, common.BKFieldName},
		DisableCounter: true,
	}

	clusterRes, ccErr := s.ClientSet.CoreService().Kube().SearchCluster(kit.Ctx, kit.Header, clusterCond)
	if ccErr != nil {
		blog.Errorf("search cluster failed, cond: %v, err: %v, rid: %s", clusterCond, ccErr, kit.Rid)
		return nil, ccErr
	}

	clusterMap := make(map[int64]types.Cluster)
	for _, cluster := range clusterRes.Data {
		clusterMap[cluster.ID] = cluster
	}

	// get biz names
	bizIDWithName, err := s.getBizIDWithName(kit, allBizIDs)
	if err != nil {
		blog.Errorf("get bizID with name failed, bizID: %+v, err: %v, rid: %s", allBizIDs, err, kit.Rid)
		return nil, err
	}

	// combine cluster and biz info for pod paths
	sharedClusterPaths := make([]types.PodPath, 0)
	for idx, path := range paths {
		cluster, exists := clusterMap[path.ClusterID]
		if !exists {
			blog.Errorf("cluster %d not exists, rid: %s", path.ClusterID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.BKClusterIDFiled)
		}

		if cluster.Name != nil {
			path.ClusterName = *cluster.Name
		}

		if path.BizID != bizID {
			path.BizName = bizIDWithName[path.BizID]
			sharedClusterPaths = append(sharedClusterPaths, path)
		}

		path.BizID = bizID
		path.BizName = bizIDWithName[bizID]
		paths[idx] = path
	}

	return append(paths, sharedClusterPaths...), nil
}

// ListPod list pod
func (s *service) ListPod(ctx *rest.Contexts) {
	req := new(types.PodQueryOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubePod, Action: acmeta.Find},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	// compatible for shared cluster scenario
	cond, err := s.Logics.KubeOperation().GenSharedNsListCond(ctx.Kit, types.KubePod, req.BizID, req.Filter)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if req.Page.EnableCount {
		counts, err := s.ClientSet.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBasePod, []map[string]interface{}{cond})
		if err != nil {
			blog.Errorf("count pod failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithCount(counts[0], make([]mapstr.MapStr, 0))
		return
	}

	if req.Page.Sort == "" {
		req.Page.Sort = common.BKFieldID
	}

	query := &metadata.QueryCondition{
		Condition: cond,
		Page:      req.Page,
		Fields:    req.Fields,
	}
	resp, err := s.ClientSet.CoreService().Kube().ListPod(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("find pod failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(0, resp.Info)
}

// BatchCreatePod batch create pods.
func (s *service) BatchCreatePod(ctx *rest.Contexts) {
	data := new(types.CreatePodsOption)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.Validate(); err.ErrCode != 0 {
		blog.Errorf("batch create pods param verification failed, data: %+v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := make([]acmeta.ResourceAttribute, len(data.Data))
	for i, data := range data.Data {
		authRes[i] = acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubePod, Action: acmeta.Create},
			BusinessID: data.BizID}
	}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes...); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	var ids []int64
	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ids, err = s.Logics.KubeOperation().BatchCreatePod(ctx.Kit, data)
		if err != nil {
			blog.Errorf("create pods failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(metadata.RspIDs{IDs: ids})
}

// DeletePods delete pods and their containers
func (s *service) DeletePods(ctx *rest.Contexts) {
	opt := new(types.DeletePodsOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := make([]acmeta.ResourceAttribute, len(opt.Data))
	for i, data := range opt.Data {
		authRes[i] = acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubePod, Action: acmeta.Delete},
			BusinessID: data.BizID}
	}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes...); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	ids, pods, err := s.checkDelPodData(ctx.Kit, opt)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// if all pods are already deleted, return
	if len(pods) == 0 {
		ctx.RespEntity(nil)
		return
	}

	// generate audit logs
	audit := auditlog.NewKubeAudit(s.ClientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditDelete)
	auditLogs, err := audit.GeneratePodAuditLog(generateAuditParameter, pods)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// delete pods
	delOpt := &types.DeletePodsByIDsOption{PodIDs: ids}

	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.ClientSet.CoreService().Kube().DeletePods(ctx.Kit.Ctx, ctx.Kit.Header, delOpt)
		if err != nil {
			blog.Errorf("delete pods failed, opt: %+v, del opt: %+v, err: %v, rid: %s", opt, delOpt, err, ctx.Kit.Rid)
			return err
		}

		// save audit logs
		if err = audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
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

func (s *service) checkDelPodData(kit *rest.Kit, opt *types.DeletePodsOption) ([]int64, []types.Pod, error) {
	// get to delete pods and containers in them
	ids := make([]int64, 0)
	idBizMap := make(map[int64]int64)

	for _, delData := range opt.Data {
		ids = append(ids, delData.PodIDs...)
		for _, id := range delData.PodIDs {
			idBizMap[id] = delData.BizID
		}
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{types.BKIDField: mapstr.MapStr{common.BKDBIN: ids}},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}

	podResp, err := s.ClientSet.CoreService().Kube().ListPod(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("find pod failed, cond: %v, err: %v, rid: %s", query, err, kit.Rid)
		return nil, nil, err
	}

	// if all pods are already deleted, return
	if len(podResp.Info) == 0 {
		return nil, nil, nil
	}

	if err := s.checkDeletePodSharedNs(kit, podResp.Info, idBizMap); err != nil {
		return nil, nil, err
	}
	return ids, podResp.Info, nil
}

// checkDeletePodSharedNs checks if pod's ns is a shared ns and if its biz id is not the same with the input biz id
func (s *service) checkDeletePodSharedNs(kit *rest.Kit, pods []types.Pod, idBizMap map[int64]int64) error {
	mismatchNsMap := make(map[int64][]int64)
	for _, pod := range pods {
		biz := idBizMap[pod.ID]
		if pod.BizID != biz {
			mismatchNsMap[biz] = append(mismatchNsMap[biz], pod.NamespaceID)
		}
	}

	if err := s.Logics.KubeOperation().CheckPlatBizSharedNs(kit, mismatchNsMap); err != nil {
		return err
	}
	return nil
}

// ListContainer list container
func (s *service) ListContainer(ctx *rest.Contexts) {
	req := new(types.ContainerQueryOption)
	err := ctx.DecodeInto(req)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeContainer, Action: acmeta.Find},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	if err = s.checkContainerPod(ctx.Kit, req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	andCond, err := filtertools.And(filtertools.GenAtomFilter(types.BKPodIDField, filter.Equal, req.PodID), req.Filter)
	if err != nil {
		blog.Errorf("generate container cond with pod failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	cond, err := andCond.ToMgo()
	if err != nil {
		blog.Errorf("parse container cond(%#v) failed, err: %v, rid: %s", andCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if req.Page.EnableCount {
		counts, err := s.ClientSet.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBaseContainer, []map[string]interface{}{cond})
		if err != nil {
			blog.Errorf("count container failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithCount(counts[0], make([]mapstr.MapStr, 0))
		return
	}

	if req.Page.Sort == "" {
		req.Page.Sort = common.BKFieldID
	}

	query := &metadata.QueryCondition{
		Condition: cond,
		Page:      req.Page,
		Fields:    req.Fields,
	}

	resp, err := s.ClientSet.CoreService().Kube().ListContainer(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("find container failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(0, resp.Info)
}

// checkContainerPod check if pod exists, compatible for shared cluster scenario
func (s *service) checkContainerPod(kit *rest.Kit, req *types.ContainerQueryOption) error {
	podQuery := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKFieldID: req.PodID},
		Page:      metadata.BasePage{Limit: 1},
		Fields:    []string{types.BKNamespaceIDField, common.BKAppIDField},
	}
	podResp, err := s.ClientSet.CoreService().Kube().ListPod(kit.Ctx, kit.Header, podQuery)
	if err != nil {
		blog.Errorf("get pod by id %d failed, err: %v, rid: %s", req.PodID, err, kit.Rid)
		return err
	}

	if len(podResp.Info) != 1 {
		blog.Errorf("get no pod by id %d, rid: %s", req.PodID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.BKPodIDField)
	}

	if podResp.Info[0].BizID != req.BizID {
		mismatchNsMap := map[int64][]int64{req.BizID: {podResp.Info[0].NamespaceID}}
		if err := s.Logics.KubeOperation().CheckPlatBizSharedNs(kit, mismatchNsMap); err != nil {
			return err
		}
	}
	return nil
}

// ListContainerByTopo list container by topo
func (s *service) ListContainerByTopo(ctx *rest.Contexts) {
	req := new(types.GetContainerByTopoOption)
	err := ctx.DecodeInto(req)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeContainer, Action: acmeta.Find},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)

	if req.Page.Sort == "" {
		req.Page.Sort = common.BKFieldID
	}

	podCond, err := req.GetPodCond()
	if err != nil {
		blog.Errorf("get pod condition failed, req: %+v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	containerCond, err := req.GetContainerCond()
	if err != nil {
		blog.Errorf("get container condition failed, req: %+v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	fields := make([]string, len(req.ContainerFields))
	copy(fields, req.ContainerFields)
	fields = append(req.ContainerFields, []string{common.BKFieldID, types.BKPodIDField}...)

	query := &types.GetContainerByPodOption{
		PodCond:       podCond,
		ContainerCond: containerCond,
		Page:          req.Page,
		Fields:        fields,
	}
	containerResp, err := s.ClientSet.CoreService().Kube().ListContainerByPod(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("find container failed, cond: %+v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if req.Page.EnableCount || len(containerResp.Info) == 0 {
		ctx.RespEntityWithCount(containerResp.Count, make([]mapstr.MapStr, 0))
		return
	}

	result, err := s.getContainerWithTopo(ctx.Kit, req.ContainerFields, req.PodFields, containerResp.Info)
	if err != nil {
		blog.Errorf("get container with topo failed, req: %+v,  err: %v, rid: %s", req, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(0, result)
}

func (s *service) getContainerWithTopo(kit *rest.Kit, containerFields []string, podFields []string,
	containers []mapstr.MapStr) ([]types.ContainerWithTopo, error) {

	// get pod detail and topo message
	podIDs := make([]int64, 0)
	containerPodMap := make(map[int64]int64)
	for _, container := range containers {
		podID, err := container.Int64(types.BKPodIDField)
		if err != nil {
			blog.Errorf("get container pod id failed, container: %+v, err: %v, rid: %s", container, err, kit.Rid)
			return nil, err
		}
		podIDs = append(podIDs, podID)

		id, err := container.Int64(common.BKFieldID)
		if err != nil {
			blog.Errorf("get container id failed, container: %+v, err: %v, rid: %s", container, err, kit.Rid)
			return nil, err
		}

		containerPodMap[id] = podID
	}
	podIDs = util.IntArrayUnique(podIDs)

	fields := make([]string, len(podFields))
	copy(fields, podFields)
	fields = append(fields, []string{types.BKBizIDField, types.BKClusterIDFiled, types.BKNamespaceIDField,
		types.RefField, common.BKHostIDField}...)
	query := &metadata.QueryCondition{
		Condition: map[string]interface{}{types.BKIDField: map[string]interface{}{common.BKDBIN: podIDs}},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Fields:    fields,
	}
	podResp, err := s.ClientSet.CoreService().Kube().ListPod(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("find pod failed, cond: %+v, err: %v, rid: %s", query, err, kit.Rid)
		return nil, err
	}
	if len(podResp.Info) != len(podIDs) {
		blog.Errorf("can not find all pod from container, podIDs: %+v, pod count: %d, rid: %s", podIDs,
			len(podResp.Info), kit.Rid)
		return nil, fmt.Errorf("can not find all pod from container, podIDs : %+v", podIDs)
	}

	podMap := make(map[int64]mapstr.MapStr)
	podTopoMap := make(map[int64]types.Topo)

	for _, pod := range podResp.Info {
		podTopoMap[pod.ID] = types.Topo{BizID: pod.BizID, ClusterID: pod.ClusterID, NamespaceID: pod.NamespaceID,
			WorkloadID: pod.Ref.ID, WorkloadType: pod.Ref.Kind, HostID: pod.HostID}

		podMap[pod.ID] = make(mapstr.MapStr)
		podKV := mapstr.NewFromStruct(pod, "json")

		for _, field := range podFields {
			podMap[pod.ID][field] = podKV[field]
		}
	}

	// construct the returned result
	result := make([]types.ContainerWithTopo, len(containers))
	for idx, container := range containers {
		containerKV := make(mapstr.MapStr)
		for _, field := range containerFields {
			containerKV[field] = container[field]
		}

		id, err := container.Int64(common.BKFieldID)
		if err != nil {
			blog.Errorf("get container id failed, container: %+v, err: %v, rid: %s", container, err, kit.Rid)
			return nil, err
		}
		podID := containerPodMap[id]

		result[idx] = types.ContainerWithTopo{
			Container: containerKV,
			Pod:       podMap[podID],
			Topo:      podTopoMap[podID],
		}
	}

	return result, nil
}
