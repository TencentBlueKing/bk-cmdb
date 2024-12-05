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
	"encoding/json"
	"errors"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	errutil "configcenter/src/common/util/errors"
	"configcenter/src/kube/types"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

// CreateWorkload create workload
func (s *service) CreateWorkload(ctx *rest.Contexts) {
	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	tableName, err := kind.Table()
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KindField))
		return
	}

	rawReq := json.RawMessage{}
	if err = ctx.DecodeInto(&rawReq); err != nil {
		ctx.RespAutoError(err)
		return
	}

	workloads, err := types.WlArrayUnmarshalJSON(kind, rawReq)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	for _, workload := range workloads {
		if rawErr := workload.ValidateCreate(); rawErr.ErrCode != 0 {
			blog.Errorf("workload %+v is invalid, err: %v, rid: %s", workload, rawErr, ctx.Kit.Rid)
			ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
			return
		}
	}

	ids, err := mongodb.Shard(ctx.Kit.SysShardOpts()).NextSequences(ctx.Kit.Ctx, tableName, len(workloads))
	if err != nil {
		blog.Errorf("get workload ids failed, table: %s, err: %v, rid: %s", tableName, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	nsIDs := make([]int64, 0)
	for _, data := range workloads {
		nsIDs = append(nsIDs, data.GetWorkloadBase().NamespaceID)
	}
	nsSpecs, err := s.GetNamespaceSpec(ctx.Kit, nsIDs)
	if err != nil {
		blog.Errorf("get namespace spec message failed, namespaceIDs: %v, err: %v, rid: %s", nsIDs, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	respData := metadata.RspIDs{IDs: make([]int64, len(ids))}
	createData := make([]types.WorkloadInterface, len(workloads))
	mismatchNsMap := make(map[int64][]int64)

	for idx, data := range workloads {
		wlBase := data.GetWorkloadBase()

		if wlBase.BizID != nsSpecs[wlBase.NamespaceID].BizID {
			mismatchNsMap[wlBase.BizID] = append(mismatchNsMap[wlBase.BizID], wlBase.NamespaceID)
		}

		wlBase.NamespaceSpec = nsSpecs[wlBase.NamespaceID]
		id := int64(ids[idx])
		wlBase.ID = id
		respData.IDs[idx] = id
		now := time.Now().Unix()
		revision := table.Revision{
			Creator:    ctx.Kit.User,
			Modifier:   ctx.Kit.User,
			CreateTime: now,
			LastTime:   now,
		}
		wlBase.Revision = revision
		wlBase.TenantID = ctx.Kit.TenantID
		data.SetWorkloadBase(wlBase)
		createData[idx] = data
	}

	// checks if workload's namespace is a shared namespace and if its biz id is not the same with the input biz id
	if err = s.core.KubeOperation().CheckPlatBizSharedNs(ctx.Kit, mismatchNsMap); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// create workloads
	err = mongodb.Shard(ctx.Kit.ShardOpts()).Table(tableName).Insert(ctx.Kit.Ctx, createData)
	if err != nil {
		blog.Errorf("add %s workload failed,data: %v, err: %v, rid: %s", tableName, createData, err, ctx.Kit.Rid)
		ctx.RespAutoError(errutil.ConvDBInsertError(ctx.Kit, err))
		return
	}

	ctx.RespEntity(respData)
}

// GetNamespaceSpec get namespace spec
func (s *service) GetNamespaceSpec(kit *rest.Kit, nsIDs []int64) (map[int64]types.NamespaceSpec,
	error) {

	if len(nsIDs) == 0 {
		blog.Errorf("namespaceIDs can not be empty, rid: %s", kit.Rid)
		return nil, errors.New("namespaceIDs can not be empty")
	}

	nsIDs = util.IntArrayUnique(nsIDs)
	filter := map[string]interface{}{
		common.BKFieldID: mapstr.MapStr{common.BKDBIN: nsIDs},
	}

	field := []string{common.BKFieldID, common.BKFieldName, common.BKAppIDField, types.BKClusterIDFiled,
		types.ClusterUIDField}
	namespaces := make([]types.Namespace, 0)
	err := mongodb.Shard(kit.ShardOpts()).Table(types.BKTableNameBaseNamespace).Find(filter).Fields(field...).
		All(kit.Ctx, &namespaces)
	if err != nil && !mongodb.IsNotFoundError(err) {
		blog.Errorf("find namespace failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(nsIDs) != len(namespaces) {
		blog.Errorf("can not find all namespace, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	nsSpecs := make(map[int64]types.NamespaceSpec)
	for _, namespace := range namespaces {
		nsSpecs[namespace.ID] = types.NamespaceSpec{
			ClusterSpec: types.ClusterSpec{
				BizID:      namespace.BizID,
				ClusterID:  namespace.ClusterID,
				ClusterUID: namespace.ClusterUID,
			},
			Namespace:   namespace.Name,
			NamespaceID: namespace.ID,
		}
	}
	return nsSpecs, nil
}

// UpdateWorkload update workload
func (s *service) UpdateWorkload(ctx *rest.Contexts) {
	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	table, err := kind.Table()
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KindField))
		return
	}

	req := types.WlUpdateByIDsOption{Kind: kind}
	if err := ctx.DecodeInto(&req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := map[string]interface{}{
		common.BKFieldID: mapstr.MapStr{common.BKDBIN: req.IDs},
	}
	updateData, err := req.Data.BuildUpdateData(ctx.Kit.User)
	if err != nil {
		blog.Errorf("get update data failed, kind: %s, info: %v, err: %v, rid: %s", kind, req.Data, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}

	err = mongodb.Shard(ctx.Kit.ShardOpts()).Table(table).Update(ctx.Kit.Ctx, cond, updateData)
	if err != nil {
		blog.Errorf("update workload failed, kind: %s, filter: %v, updateData: %v, err: %v, rid: %s", kind, cond,
			updateData, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}

	ctx.RespEntity(nil)
}

// DeleteWorkload delete workload
func (s *service) DeleteWorkload(ctx *rest.Contexts) {
	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	table, err := kind.Table()
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KindField))
		return
	}

	req := new(types.WlDeleteByIDsOption)
	if err := ctx.DecodeInto(req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	filter := mapstr.MapStr{
		common.BKFieldID: mapstr.MapStr{common.BKDBIN: req.IDs},
	}
	if err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(table).Delete(ctx.Kit.Ctx, filter); err != nil {
		blog.Errorf("delete workload failed, filter: %v, err: %v, rid: %s", filter, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
		return
	}

	ctx.RespEntity(nil)
}

// ListWorkload list container
func (s *service) ListWorkload(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	if err := kind.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}

	table, err := kind.Table()
	if err != nil {
		blog.Errorf("workload kind is invalid, kind: %v, rid: %s", kind, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KindField))
		return
	}
	workloads := make([]mapstr.MapStr, 0)
	err = mongodb.Shard(ctx.Kit.ShardOpts()).Table(table).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &workloads)
	if err != nil {
		blog.Errorf("search workload failed, cond: %v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &metadata.QueryResult{
		Info: workloads,
	}
	ctx.RespEntity(result)
}
