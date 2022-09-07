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
	"configcenter/src/kube/orm"
	"configcenter/src/kube/types"
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

	req := types.NsCreateReq{}
	if err := ctx.DecodeInto(&req); nil != err {
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

	respData := types.NsCreateRespData{
		IDs: make([]int64, len(ids)),
	}
	for idx, data := range req.Data {
		spec := &types.ClusterSpec{
			BizID:      &bizID,
			ClusterID:  data.ClusterID,
			ClusterUID: data.ClusterUID,
		}
		spec, err := s.GetClusterSpec(ctx.Kit, spec)
		if err != nil {
			blog.Errorf("get cluster spec message failed, data: %v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		id := int64(ids[idx])
		respData.IDs[idx] = id
		data.ClusterSpec = *spec
		data.ID = &id
		now := time.Now().Unix()
		data.CreateTime = &now
		data.UpdateTime = &now
		data.SupplierAccount = &ctx.Kit.SupplierAccount

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
func (s *coreService) GetClusterSpec(kit *rest.Kit, spec *types.ClusterSpec) (*types.ClusterSpec, error) {
	if spec.BizID == nil {
		blog.Errorf("bizID can not be empty, rid: %s", kit.Rid)
		return nil, errors.New("bizID can not be empty")
	}

	if spec.ClusterID == nil && spec.ClusterUID == nil {
		blog.Errorf("clusterID and clusterUID can not be empty at the same time, rid: %s", kit.Rid)
		return nil, errors.New("clusterID and clusterUID can not be empty at the same time")
	}

	filter := map[string]interface{}{
		common.BKAppIDField:   *spec.BizID,
		common.BKOwnerIDField: kit.SupplierAccount,
	}
	if spec.ClusterID != nil {
		filter[common.BKFieldID] = *spec.ClusterID
	}
	if spec.ClusterUID != nil {
		filter[types.UidField] = *spec.ClusterUID
	}

	cluster := types.Cluster{}
	if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(filter).One(kit.Ctx, &cluster); err != nil {
		if mongodb.Client().IsNotFoundError(err) {
			blog.Errorf("can not find cluster, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommNotFound)
		}

		blog.Errorf("find cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	result := types.ClusterSpec{
		BizID:      &cluster.BizID,
		ClusterID:  &cluster.ID,
		ClusterUID: cluster.Uid,
	}
	return &result, nil
}

// UpdateNamespace update namespace
func (s *coreService) UpdateNamespace(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := types.NsUpdateReq{}
	if err := ctx.DecodeInto(&req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// build filter
	for _, data := range req.Data {
		var filter map[string]interface{}
		if len(data.ID) != 0 {
			filter = map[string]interface{}{
				common.BKAppIDField:   bizID,
				common.BKOwnerIDField: ctx.Kit.SupplierAccount,
				common.BKFieldID: map[string]interface{}{
					common.BKDBIN: data.ID,
				},
			}
		}

		if len(data.Unique) != 0 {
			orCond := make([]map[string]interface{}, 0)
			for _, unique := range data.Unique {
				cond := map[string]interface{}{
					types.ClusterUIDField: unique.ClusterUID,
					common.BKFieldName:    unique.Name,
				}
				orCond = append(orCond, cond)
			}
			filter = map[string]interface{}{
				common.BKAppIDField:   bizID,
				common.BKOwnerIDField: ctx.Kit.SupplierAccount,
				common.BKDBOR:         orCond,
			}
		}
		now := time.Now().Unix()
		data.Info.UpdateTime = &now
		// build update data
		opts := orm.NewFieldOptions().AddIgnoredFields(common.BKFieldID, types.ClusterUIDField, common.BKFieldName)
		updateData, err := orm.GetUpdateFieldsWithOption(data.Info, opts)
		if err != nil {
			blog.Errorf("get update data failed, data: %v, err: %v, rid: %s", data.Info, err, ctx.Kit.Rid)
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

	req := &types.NsDeleteReq{}
	if err := ctx.DecodeInto(&req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	for _, data := range req.Data {
		var filter map[string]interface{}
		if data.ID != 0 {
			filter = map[string]interface{}{
				common.BKAppIDField:   bizID,
				common.BKOwnerIDField: ctx.Kit.SupplierAccount,
				common.BKFieldID:      data.ID,
			}
		} else {
			filter = map[string]interface{}{
				common.BKAppIDField:   bizID,
				common.BKOwnerIDField: ctx.Kit.SupplierAccount,
				types.ClusterUIDField: data.ClusterUID,
				common.BKFieldName:    data.Name,
			}
		}

		err = mongodb.Client().Table(types.BKTableNameBaseNamespace).Delete(ctx.Kit.Ctx, filter)
		if err != nil {
			blog.Errorf("delete namespace failed, filter: %v, err: %v, rid: %s", filter, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
			return
		}
	}

	ctx.RespEntity(nil)
}
