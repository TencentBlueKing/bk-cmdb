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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util/errors"
	"configcenter/src/kube/orm"
	"configcenter/src/kube/types"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

// SearchClusters search clusters
func (s *service) SearchClusters(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	clusters := make([]types.Cluster, 0)

	err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &clusters)
	if err != nil {
		blog.Errorf("search cluster failed, cond: %+v, err: %v, rid: %s", input.Condition, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	result := &types.ResponseCluster{Data: clusters}

	ctx.RespEntity(result)
}

// BatchUpdateCluster update cluster.
func (s *service) BatchUpdateCluster(ctx *rest.Contexts) {
	input := new(types.UpdateClusterByIDsOption)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := input.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	filter := map[string]interface{}{
		types.BKIDField: map[string]interface{}{
			common.BKDBIN: input.IDs,
		},
	}

	opts := orm.NewFieldOptions().AddIgnoredFields(types.IgnoredUpdateClusterFields...)
	updateData, err := orm.GetUpdateFieldsWithOption(input.Data, opts)
	if err != nil {
		blog.Errorf("get update data failed, data: %v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommDBUpdateFailed))
		return
	}

	err = mongodb.Client().Table(types.BKTableNameBaseCluster).Update(ctx.Kit.Ctx, filter, updateData)
	if err != nil {
		blog.Errorf("update cluster failed, filter: %v, updateData: %v, err: %v, rid: %s", filter, updateData,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommDBUpdateFailed))
		return
	}
	ctx.RespEntity(nil)
}

// CreateCluster create kube cluster.
func (s *service) CreateCluster(ctx *rest.Contexts) {
	cluster := new(types.Cluster)
	if err := ctx.DecodeInto(cluster); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := cluster.ValidateCreate(); err.ErrCode != 0 {
		blog.Errorf("create cluster failed, data: %+v, err: %+v, rid: %s", cluster, err, ctx.Kit.Rid)
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	id, err := mongodb.Client().NextSequence(ctx.Kit.Ctx, types.BKTableNameBaseCluster)
	if err != nil {
		blog.Errorf("create cluster failed, generate id failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed))
		return
	}
	cluster.ID = int64(id)

	now := time.Now().Unix()
	cluster.Revision = table.Revision{
		Creator:    ctx.Kit.User,
		Modifier:   ctx.Kit.User,
		CreateTime: now,
		LastTime:   now,
	}
	cluster.SupplierAccount = ctx.Kit.SupplierAccount

	err = mongodb.Client().Table(types.BKTableNameBaseCluster).Insert(ctx.Kit.Ctx, cluster)
	if err != nil {
		blog.Errorf("create cluster failed, db insert failed, doc: %+v, err: %+v, rid: %s", cluster, err, ctx.Kit.Rid)
		ctx.RespAutoError(errors.ConvDBInsertError(ctx.Kit, mongodb.Client(), err))
		return
	}

	ctx.RespEntity(cluster)
}

// BatchDeleteCluster delete clusters.
func (s *service) BatchDeleteCluster(ctx *rest.Contexts) {
	option := new(types.DeleteClusterByIDsOption)
	if err := ctx.DecodeInto(option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	filter := map[string]interface{}{
		common.BKOwnerIDField: ctx.Kit.SupplierAccount,
		types.BKIDField: map[string]interface{}{
			common.BKDBIN: option.IDs,
		},
	}

	if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Delete(ctx.Kit.Ctx, filter); err != nil {
		blog.Errorf("delete cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
