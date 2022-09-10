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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

// CreateWorkload create workload
func (s *coreService) CreateWorkload(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	tableName, err := kind.Table()
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KindField))
		return
	}

	req := types.WlCreateReq{Kind: kind}
	if err := ctx.DecodeInto(&req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	ids, err := mongodb.Client().NextSequences(ctx.Kit.Ctx, tableName, len(req.Data))
	if err != nil {
		blog.Errorf("get workload ids failed, table: %s, err: %v, rid: %s", tableName, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	respData := types.WlCreateRespData{
		IDs: make([]int64, len(ids)),
	}
	for idx, data := range req.Data {
		spec := data.GetNamespaceSpec()
		spec.BizID = &bizID
		nsSpec, err := s.GetNamespaceSpec(ctx.Kit, &spec)
		if err != nil {
			blog.Errorf("get namespace spec message failed, data: %v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		id := int64(ids[idx])
		data.SetNamespaceSpec(*nsSpec)
		respData.IDs[idx] = id
		data.SetID(id)
		now := time.Now().Unix()
		data.SetCreateMsg(
			table.Revision{
				Creator:    ctx.Kit.User,
				Modifier:   ctx.Kit.User,
				CreateTime: now,
				LastTime:   now,
			},
		)
		data.SetSupplierAccount(ctx.Kit.SupplierAccount)
		err = mongodb.Client().Table(tableName).Insert(ctx.Kit.Ctx, data)
		if err != nil {
			blog.Errorf("add workload failed, table: %s, data: %v, err: %v, rid: %s", tableName, data, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
			return
		}
	}

	ctx.RespEntity(respData)
}

// GetNamespaceSpec get namespace spec
func (s *coreService) GetNamespaceSpec(kit *rest.Kit, spec *types.NamespaceSpec) (*types.NamespaceSpec, error) {
	if spec.BizID == nil {
		blog.Errorf("bizID can not be empty, rid: %s", kit.Rid)
		return nil, errors.New("bizID can not be empty")
	}

	if spec.NamespaceID == nil && (spec.ClusterUID == nil || spec.Namespace == nil) {
		blog.Errorf("namespaceID and namespaceUID can not be empty at the same time, rid: %s", kit.Rid)
		return nil, errors.New("namespaceID and namespaceUID can not be empty at the same time")
	}

	filter := map[string]interface{}{
		common.BKAppIDField:   *spec.BizID,
		common.BKOwnerIDField: kit.SupplierAccount,
	}
	if spec.NamespaceID != nil {
		filter[common.BKFieldID] = *spec.NamespaceID
	}
	if spec.ClusterUID != nil && spec.Namespace != nil {
		filter[types.ClusterUIDField] = *spec.ClusterUID
		filter[common.BKFieldName] = *spec.Namespace
	}

	ns := types.Namespace{}
	if err := mongodb.Client().Table(types.BKTableNameBaseNamespace).Find(filter).One(kit.Ctx, &ns); err != nil {
		if mongodb.Client().IsNotFoundError(err) {
			blog.Errorf("can not find namespace, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommNotFound)
		}

		blog.Errorf("find namespace failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	result := types.NamespaceSpec{
		ClusterSpec: types.ClusterSpec{
			BizID:      ns.BizID,
			ClusterID:  ns.ClusterID,
			ClusterUID: ns.ClusterUID,
		},
		Namespace:   ns.Name,
		NamespaceID: ns.ID,
	}
	return &result, nil
}

// UpdateWorkload update workload
func (s *coreService) UpdateWorkload(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	table, err := kind.Table()
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KindField))
		return
	}

	req := types.WlUpdateReq{Kind: kind}
	if err := ctx.DecodeInto(&req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	for _, data := range req.Data {
		filter := data.BuildUpdateFilter(bizID, ctx.Kit.SupplierAccount)
		updateData, err := data.BuildUpdateData(ctx.Kit.User)
		if err != nil {
			blog.Errorf("get update data failed, kind: %s, data: %v, err: %v, rid: %s", kind, data, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
			return
		}

		err = mongodb.Client().Table(table).Update(ctx.Kit.Ctx, filter, updateData)
		if err != nil {
			blog.Errorf("update workload failed, kind: %s, filter: %v, updateData: %v, err: %v, rid: %s", kind, filter,
				updateData, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
			return
		}
	}

	ctx.RespEntity(nil)
}

// DeleteWorkload delete workload
func (s *coreService) DeleteWorkload(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	table, err := kind.Table()
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KindField))
		return
	}

	req := types.WlDeleteReq{}
	if err := ctx.DecodeInto(&req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	filter, err := req.BuildCond(bizID, ctx.Kit.SupplierAccount)
	if err != nil {
		blog.Errorf("delete workload failed, bizID: %s, data: %v, err: %v, rid: %s", bizID, req, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if err := mongodb.Client().Table(table).Delete(ctx.Kit.Ctx, filter); err != nil {
		blog.Errorf("delete workload failed, filter: %v, err: %v, rid: %s", filter, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
		return
	}

	ctx.RespEntity(nil)
}

// ListWorkload list container
func (s *coreService) ListWorkload(ctx *rest.Contexts) {
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
	err = mongodb.Client().Table(table).Find(input.Condition).Start(uint64(input.Page.Start)).
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
