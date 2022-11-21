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
	"configcenter/src/common/util"
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

	req := types.WlCreateOption{Kind: kind}
	if err := ctx.DecodeInto(&req); err != nil {
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

	nsIDs := make([]int64, 0)
	for _, data := range req.Data {
		nsIDs = append(nsIDs, data.GetWorkloadBase().NamespaceID)
	}
	nsSpecs, err := s.GetNamespaceSpec(ctx.Kit, bizID, nsIDs)
	if err != nil {
		blog.Errorf("get namespace spec message failed, bizID: %s, namespaceIDs: %v, err: %v, rid: %s", bizID, nsIDs,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	respData := metadata.RspIDs{
		IDs: make([]int64, len(ids)),
	}
	for idx, data := range req.Data {
		wlBase := data.GetWorkloadBase()
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
		wlBase.SupplierAccount = ctx.Kit.SupplierAccount
		data.SetWorkloadBase(wlBase)
		err = mongodb.Client().Table(tableName).Insert(ctx.Kit.Ctx, data)
		if err != nil {
			blog.Errorf("add workload failed, table: %s, data: %v, err: %v, rid: %s", tableName, data, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
			return
		}
	}

	ctx.RespEntity(respData)
}

// GetNamespaceInfo get namespace spec
func (s *coreService) GetNamespaceSpec(kit *rest.Kit, bizID int64, nsIDs []int64) (map[int64]types.NamespaceSpec,
	error) {

	if bizID == 0 {
		blog.Errorf("bizID can not be empty, rid: %s", kit.Rid)
		return nil, errors.New("bizID can not be empty")
	}

	if len(nsIDs) == 0 {
		blog.Errorf("namespaceIDs can not be empty, rid: %s", kit.Rid)
		return nil, errors.New("namespaceIDs can not be empty")
	}

	nsIDs = util.IntArrayUnique(nsIDs)
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: nsIDs},
	}
	filter = util.SetQueryOwner(filter, kit.SupplierAccount)

	field := []string{common.BKFieldID, common.BKFieldName, types.BKClusterIDFiled, types.ClusterUIDField}
	namespaces := make([]types.Namespace, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseNamespace).Find(filter).Fields(field...).
		All(kit.Ctx, &namespaces)
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("find namespace failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(nsIDs) != len(namespaces) {
		blog.Errorf("can not find all namespace, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	nsSpecs := make(map[int64]types.NamespaceSpec)
	for _, namespace := range namespaces {
		nsSpecs[namespace.ID] = types.NamespaceSpec{
			ClusterSpec: types.ClusterSpec{
				BizID:      bizID,
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

	req := types.WlUpdateOption{Kind: kind}
	if err := ctx.DecodeInto(&req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := map[string]interface{}{
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: req.IDs},
		common.BKAppIDField: bizID,
	}
	util.SetModOwner(cond, ctx.Kit.SupplierAccount)
	updateData, err := req.Data.BuildUpdateData(ctx.Kit.User)
	if err != nil {
		blog.Errorf("get update data failed, kind: %s, info: %v, err: %v, rid: %s", kind, req.Data, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}

	err = mongodb.Client().Table(table).Update(ctx.Kit.Ctx, cond, updateData)
	if err != nil {
		blog.Errorf("update workload failed, kind: %s, filter: %v, updateData: %v, err: %v, rid: %s", kind, cond,
			updateData, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
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

	req := new(types.WlDeleteOption)
	if err := ctx.DecodeInto(req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	filter := mapstr.MapStr{
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: req.IDs},
		common.BKAppIDField: bizID,
	}
	util.SetModOwner(filter, ctx.Kit.SupplierAccount)
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
	util.SetQueryOwner(input.Condition, ctx.Kit.SupplierAccount)
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
