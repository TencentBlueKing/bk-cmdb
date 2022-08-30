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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// SearchClusters 根据用户指定的条件查询cluster
func (s *Service) SearchClusters(ctx *rest.Contexts) {

	searchCond := new(types.QueryClusterReq)
	if err := ctx.DecodeInto(searchCond); err != nil {
		blog.Errorf("failed to parse the params, error: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	if cErr := searchCond.Validate(); cErr.ErrCode != 0 {
		blog.Errorf("validate request failed, err: %v, rid: %s", cErr, ctx.Kit.Rid)
		ctx.RespAutoError(cErr.ToCCError(ctx.Kit.CCError))
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	filter := mapstr.New()
	if searchCond.Filter != nil {
		cond, errKey, rawErr := searchCond.Filter.ToMgo()
		if rawErr != nil {
			blog.Errorf("parse biz filter(%#v) failed, err: %v, rid: %s", searchCond.Filter, rawErr, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey))
			return
		}
		filter = cond
	}
	// 无论条件中是否有bk_biz_id、supplier_account,这里统一替换成url中的bk_biz_id 和kit中的supplier_account
	filter[types.BKBizIDField] = bizID
	filter[types.BKSupplierAccountField] = ctx.Kit.SupplierAccount

	// count biz in cluster enable count is set
	if searchCond.Page.EnableCount {
		filter := []map[string]interface{}{filter}

		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBaseCluster, filter)
		if err != nil {
			blog.Errorf("count biz failed, err: %v, cond: %#v, rid: %s", err, filter, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithCount(counts[0], make([]mapstr.MapStr, 0))
		return
	}

	query := &metadata.QueryCondition{
		Condition: filter,
		Page:      searchCond.Page,
		Fields:    searchCond.Fields,
	}
	result, err := s.Logics.ContainerOperation().SearchCluster(ctx.Kit, query)
	if err != nil {
		blog.Errorf("search cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
	ctx.RespEntityWithCount(0, result.Data)
}

// UpdateClusterFields 更新集群字段
func (s *Service) UpdateClusterFields(ctx *rest.Contexts) {
	// 还是用之前的结构体作为参数，json解析后看下不空的字段，根据tag获取是否是可编辑的字段，如果是可编辑的字段那么可以更新，如果是不可编辑字
	// 段那么直接报错
	data := new(types.UpdateClusterOption)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		err = s.Logics.ContainerOperation().UpdateClusterFields(ctx.Kit, data, bizID, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("create cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
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

// CreateCluster create a container cluster
func (s *Service) CreateCluster(ctx *rest.Contexts) {
	data := new(types.ClusterBaseFields)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.ValidateCreate(); err != nil {
		blog.Errorf("validate create container cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	var id int64
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		id, err = s.Logics.ContainerOperation().CreateCluster(ctx.Kit, data, bizID, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("create cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(id)
}

// DeleteCluster delete cluster.
func (s *Service) DeleteCluster(ctx *rest.Contexts) {
	option := new(types.DeleteClusterOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	//1、重新整理一下请求

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		err = s.Logics.ContainerOperation().DeleteCluster(ctx.Kit, bizID, option, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("delete cluster failed, biz: %d, option: %+v, err: %v, rid: %s", bizID, option, err, ctx.Kit.Rid)
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
