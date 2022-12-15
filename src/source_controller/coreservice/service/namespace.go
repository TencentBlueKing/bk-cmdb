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
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/kube/orm"
	"configcenter/src/kube/types"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

// CreateNamespace create namespace
func (s *coreService) CreateNamespace(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := new(types.NsCreateOption)
	if err := ctx.DecodeInto(req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	ids, err := mongodb.Client().NextSequences(ctx.Kit.Ctx, types.BKTableNameBaseNamespace, len(req.Data))
	if err != nil {
		blog.Errorf("get namespace ids failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	clusterIDs := make([]int64, 0)
	for _, data := range req.Data {
		clusterIDs = append(clusterIDs, data.ClusterID)
	}

	clusterSpecs, err := s.GetClusterSpec(ctx.Kit, bizID, clusterIDs)
	if err != nil {
		blog.Errorf("get cluster spec failed, bizID: %d, clusterIDs: %v, err: %v, rid: %s", bizID, clusterIDs, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	respData := metadata.RspIDs{
		IDs: make([]int64, len(ids)),
	}
	for idx, data := range req.Data {
		id := int64(ids[idx])
		respData.IDs[idx] = id
		data.ClusterSpec = clusterSpecs[data.ClusterID]
		data.ID = id
		now := time.Now().Unix()
		data.Revision = table.Revision{
			Creator:    ctx.Kit.User,
			Modifier:   ctx.Kit.User,
			CreateTime: now,
			LastTime:   now,
		}
		data.SupplierAccount = ctx.Kit.SupplierAccount

		err = mongodb.Client().Table(types.BKTableNameBaseNamespace).Insert(ctx.Kit.Ctx, &data)
		if err != nil {
			blog.Errorf("add namespace failed, data: %v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
			return
		}
	}

	ctx.RespEntity(respData)
}

// GetClusterSpec get cluster spec
func (s *coreService) GetClusterSpec(kit *rest.Kit, bizID int64, clusterIDs []int64) (map[int64]types.ClusterSpec,
	error) {

	if bizID == 0 {
		blog.Errorf("bizID can not be empty, rid: %s", kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	if len(clusterIDs) == 0 {
		blog.Errorf("clusterIDs can not be empty, rid: %s", kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.BKClusterIDFiled)
	}

	clusterIDs = util.IntArrayUnique(clusterIDs)
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: clusterIDs},
	}
	util.SetModOwner(filter, kit.SupplierAccount)

	field := []string{common.BKFieldID, types.UidField}
	clusters := make([]types.Cluster, 0)

	err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(filter).Fields(field...).All(kit.Ctx, &clusters)
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("find cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(clusterIDs) != len(clusters) {
		blog.Errorf("can not find all cluster, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	specs := make(map[int64]types.ClusterSpec, len(clusters))
	for _, cluster := range clusters {
		specs[cluster.ID] = types.ClusterSpec{
			BizID:      bizID,
			ClusterID:  cluster.ID,
			ClusterUID: *cluster.Uid,
		}
	}

	return specs, nil
}

// UpdateNamespace update namespace
func (s *coreService) UpdateNamespace(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := new(types.NsUpdateOption)
	if err := ctx.DecodeInto(req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// build filter
	filter := mapstr.MapStr{
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: req.IDs},
		common.BKAppIDField: bizID,
	}
	filter = util.SetModOwner(filter, ctx.Kit.SupplierAccount)
	now := time.Now().Unix()
	req.Data.LastTime = now
	req.Data.Modifier = ctx.Kit.User
	// build update data
	opts := orm.NewFieldOptions().AddIgnoredFields(common.BKFieldID, types.ClusterUIDField, common.BKFieldName)
	updateData, err := orm.GetUpdateFieldsWithOption(req.Data, opts)
	if err != nil {
		blog.Errorf("get update data failed, data: %v, err: %v, rid: %s", req.Data, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}

	// update namespace
	err = mongodb.Client().Table(types.BKTableNameBaseNamespace).Update(ctx.Kit.Ctx, filter, updateData)
	if err != nil {
		blog.Errorf("update namespace failed, filter: %v, updateData: %v, err: %v, rid: %s", filter, updateData,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}

	ctx.RespEntity(nil)
}

// DeleteNamespace delete namespace
func (s *coreService) DeleteNamespace(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := new(types.NsDeleteOption)
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
	filter = util.SetModOwner(filter, ctx.Kit.SupplierAccount)
	if err := mongodb.Client().Table(types.BKTableNameBaseNamespace).Delete(ctx.Kit.Ctx, filter); err != nil {
		blog.Errorf("delete namespace failed, filter: %v, err: %v, rid: %s", filter, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
		return
	}

	ctx.RespEntity(nil)
}

// ListNamespace list namespace
func (s *coreService) ListNamespace(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	util.SetQueryOwner(input.Condition, ctx.Kit.SupplierAccount)
	namespaces := make([]types.Namespace, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseNamespace).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &namespaces)
	if err != nil {
		blog.Errorf("search namespace failed, cond: %v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &types.NsDataResp{Data: namespaces}
	ctx.RespEntity(result)
}
