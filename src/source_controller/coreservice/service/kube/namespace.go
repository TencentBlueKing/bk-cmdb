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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/util/errors"
	"configcenter/src/kube/orm"
	"configcenter/src/kube/types"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

// CreateNamespace create namespace
func (s *service) CreateNamespace(ctx *rest.Contexts) {
	namespaces := make([]types.Namespace, 0)
	if err := ctx.DecodeInto(&namespaces); err != nil {
		ctx.RespAutoError(err)
		return
	}

	clusterIDs := make([]int64, len(namespaces))
	for i, namespace := range namespaces {
		if rawErr := namespace.ValidateCreate(); rawErr.ErrCode != 0 {
			blog.Errorf("namespace %+v is invalid, err: %v, rid: %s", namespace, rawErr, ctx.Kit.Rid)
			ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
			return
		}

		clusterIDs[i] = namespace.ClusterID
	}

	clusterMap, err := s.getClusterMap(ctx.Kit, clusterIDs)
	if err != nil {
		blog.Errorf("get cluster spec failed, clusterIDs: %v, err: %v, rid: %s", clusterIDs, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ids, err := mongodb.Shard(ctx.Kit.SysShardOpts()).NextSequences(ctx.Kit.Ctx, types.BKTableNameBaseNamespace,
		len(namespaces))
	if err != nil {
		blog.Errorf("get namespace ids failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	respData := metadata.RspIDs{IDs: make([]int64, len(ids))}
	sharedRel := make([]types.NsSharedClusterRel, 0)
	for idx, data := range namespaces {
		cluster := clusterMap[data.ClusterID]
		id := int64(ids[idx])
		respData.IDs[idx] = id
		if cluster.Uid != nil {
			data.ClusterUID = *cluster.Uid
		}
		data.ID = id
		now := time.Now().Unix()
		data.Revision = table.Revision{
			Creator:    ctx.Kit.User,
			Modifier:   ctx.Kit.User,
			CreateTime: now,
			LastTime:   now,
		}
		namespaces[idx] = data

		if cluster.BizID == data.BizID {
			continue
		}

		// if cluster and ns biz id is not equal, check if it's shared cluster, add a ns relation for shared cluster
		if cluster.Type == nil || *cluster.Type != types.SharedClusterType {
			blog.Errorf("namespace cluster %d type is not shared cluster, rid: %s", cluster.ID, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.TypeField))
			return
		}

		sharedRel = append(sharedRel, types.NsSharedClusterRel{
			NamespaceID: id,
			ClusterID:   cluster.ID,
			BizID:       data.BizID,
			AsstBizID:   cluster.BizID,
			TenantID:    ctx.Kit.TenantID,
		})
	}

	err = mongodb.Shard(ctx.Kit.ShardOpts()).Table(types.BKTableNameBaseNamespace).Insert(ctx.Kit.Ctx, namespaces)
	if err != nil {
		blog.Errorf("add namespace failed, data: %+v, err: %v, rid: %s", namespaces, err, ctx.Kit.Rid)
		ctx.RespAutoError(errors.ConvDBInsertError(ctx.Kit, err))
		return
	}

	if len(sharedRel) > 0 {
		err = mongodb.Shard(ctx.Kit.ShardOpts()).Table(types.BKTableNameNsSharedClusterRel).Insert(ctx.Kit.Ctx,
			sharedRel)
		if err != nil {
			blog.Errorf("add shared cluster relations failed, rel: %v, err: %v, rid: %s", sharedRel, err,
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
			return
		}
	}

	ctx.RespEntity(respData)
}

// getClusterMap get cluster id to cluster info map
func (s *service) getClusterMap(kit *rest.Kit, clusterIDs []int64) (map[int64]types.Cluster, error) {
	if len(clusterIDs) == 0 {
		blog.Errorf("clusterIDs can not be empty, rid: %s", kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.BKClusterIDFiled)
	}

	clusterIDs = util.IntArrayUnique(clusterIDs)
	filter := map[string]interface{}{
		common.BKFieldID: mapstr.MapStr{common.BKDBIN: clusterIDs},
	}

	field := []string{common.BKFieldID, types.UidField, types.TypeField, common.BKAppIDField}
	clusters := make([]types.Cluster, 0)

	err := mongodb.Shard(kit.ShardOpts()).Table(types.BKTableNameBaseCluster).Find(filter).Fields(field...).
		All(kit.Ctx, &clusters)
	if err != nil && !mongodb.IsNotFoundError(err) {
		blog.Errorf("find cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(clusterIDs) != len(clusters) {
		blog.Errorf("can not find all cluster, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	clusterMap := make(map[int64]types.Cluster, len(clusters))

	for _, cluster := range clusters {
		clusterMap[cluster.ID] = cluster
	}

	return clusterMap, nil
}

// UpdateNamespace update namespace
func (s *service) UpdateNamespace(ctx *rest.Contexts) {
	req := new(types.NsUpdateByIDsOption)
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
		common.BKFieldID: mapstr.MapStr{common.BKDBIN: req.IDs},
	}
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
	err = mongodb.Shard(ctx.Kit.ShardOpts()).Table(types.BKTableNameBaseNamespace).Update(ctx.Kit.Ctx, filter,
		updateData)
	if err != nil {
		blog.Errorf("update namespace failed, filter: %v, updateData: %v, err: %v, rid: %s", filter, updateData,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}

	ctx.RespEntity(nil)
}

// DeleteNamespace delete namespace
func (s *service) DeleteNamespace(ctx *rest.Contexts) {
	req := new(types.NsDeleteByIDsOption)
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
	if err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(types.BKTableNameBaseNamespace).Delete(ctx.Kit.Ctx,
		filter); err != nil {
		blog.Errorf("delete namespace failed, filter: %v, err: %v, rid: %s", filter, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
		return
	}

	// delete all shared cluster relations of the namespaces
	sharedRelCond := mapstr.MapStr{types.BKNamespaceIDField: mapstr.MapStr{common.BKDBIN: req.IDs}}
	err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(types.BKTableNameNsSharedClusterRel).Delete(ctx.Kit.Ctx,
		sharedRelCond)
	if err != nil {
		blog.Errorf("delete shared cluster rel failed, cond: %v, err: %v, rid: %s", sharedRelCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
		return
	}

	ctx.RespEntity(nil)
}

// ListNamespace list namespace
func (s *service) ListNamespace(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	namespaces := make([]types.Namespace, 0)
	err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(types.BKTableNameBaseNamespace).Find(input.Condition).
		Start(uint64(input.Page.Start)).
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
