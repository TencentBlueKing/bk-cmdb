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
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/kube/orm"
	"configcenter/src/kube/types"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

// SearchClusters search clusters
func (s *coreService) SearchClusters(ctx *rest.Contexts) {

	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	clusters := make([]types.Cluster, 0)
	util.SetQueryOwner(input.Condition, ctx.Kit.SupplierAccount)

	err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &clusters)
	if err != nil {
		blog.Errorf("search cluster failed, cond: %+v, err: %v, rid: %s", input.Condition, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &types.ResponseCluster{Data: clusters}

	ctx.RespEntity(result)
}

// SearchNsClusterRelation search ns clusters relation in the shared cluster scenario.
func (s *coreService) SearchNsClusterRelation(ctx *rest.Contexts) {

	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	relations := make([]types.NsClusterRelation, 0)
	util.SetQueryOwner(input.Condition, ctx.Kit.SupplierAccount)

	err := mongodb.Client().Table(types.BKTableNsClusterRelation).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("search ns and cluster relation failed, cond: %+v, err: %v, rid: %s", input.Condition,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &types.ResponseNsClusterRelation{Data: relations}

	ctx.RespEntity(result)
}

// validateUpdateClusterType for shared cluster scenarios, if you need to change the type
// field from non-share to share then it is necessary to judge whether there is an associated
// relationship table between ns and node. If there is a relationship table, it cannot be deleted.
func validateUpdateClusterType(kit *rest.Kit, input *types.UpdateClusterOption, bizID int64,
	filter map[string]interface{}) error {
	if input.Data.Type == nil || *input.Data.Type == types.ClusterShareTypeField {
		return nil
	}

	clusters := make([]types.Cluster, 0)
	util.SetQueryOwner(filter, kit.SupplierAccount)

	if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(filter).
		Fields(types.TypeField, types.BKIDField).All(kit.Ctx, &clusters); err != nil {
		blog.Errorf("search cluster failed, cond: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return err
	}
	// get the clusterID of the shared cluster
	ids := make([]int64, 0)
	for _, cluster := range clusters {
		if *cluster.Type == types.ClusterShareTypeField {
			ids = append(ids, cluster.ID)
		}
	}
	if len(ids) == 0 {
		return nil
	}

	// if it is a shared cluster field, you need to check whether
	// the node relation table and ns relation table have information
	// about the cluster, and if so, it cannot be updated.
	nsFilter := map[string]interface{}{
		types.BKClusterIDField: map[string]interface{}{
			common.BKDBIN: ids,
		},
		common.BKAppIDField: bizID,
	}
	count, err := mongodb.Client().Table(types.BKTableNsClusterRelation).Find(nsFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("query ns relation failed, filter: %+v, err: %+v, rid: %s", nsFilter, err, kit.Rid)
		return err
	}
	if count > 0 {
		blog.Errorf("existing ns relationship tables cannot be deleted, filter: %+v, rid: %s", nsFilter, kit.Rid)
		return kit.CCError.Error(common.CCErrCoreServiceNsRelaitionExist)
	}

	nodeFilter := map[string]interface{}{
		types.BKClusterIDField: map[string]interface{}{
			common.BKDBIN: ids,
		},
		common.BKAppIDField: bizID,
	}
	count, err = mongodb.Client().Table(types.BKTableNodeClusterRelation).Find(nodeFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("query node relation failed, filter: %+v, err: %+v, rid: %s", nodeFilter, err, kit.Rid)
		return err
	}
	if count > 0 {
		blog.Errorf("existing node relationship tables cannot be deleted, filter: %+v, rid: %s", nodeFilter, kit.Rid)
		return kit.CCError.Error(common.CCErrCoreServiceNodeRelaitionExist)
	}

	return nil
}

// BatchUpdateCluster update cluster.
func (s *coreService) BatchUpdateCluster(ctx *rest.Contexts) {

	input := new(types.UpdateClusterOption)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, err: %v, rid: %s",
			ctx.Request.PathParameter("bk_biz_id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	filter := map[string]interface{}{
		types.BKBizIDField: bizID,
		types.BKIDField: map[string]interface{}{
			common.BKDBIN: input.IDs,
		},
	}

	util.SetModOwner(filter, ctx.Kit.SupplierAccount)

	opts := orm.NewFieldOptions().AddIgnoredFields(types.IgnoredUpdateClusterFields...)
	updateData, err := orm.GetUpdateFieldsWithOption(input.Data, opts)
	if err != nil {
		blog.Errorf("get update data failed, data: %v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommDBUpdateFailed))
		return
	}

	if err := validateUpdateClusterType(ctx.Kit, input, bizID, filter); err != nil {
		ctx.RespAutoError(err)
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
func (s *coreService) CreateCluster(ctx *rest.Contexts) {

	inputData := new(types.Cluster)
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	bizStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizStr, 10, 64)
	if err != nil {
		blog.Error("url param bk_biz_id not integer, bizID: %s, rid: %s", bizStr, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	cluster, err := createCluster(ctx.Kit, bizID, inputData)

	ctx.RespEntityWithError(cluster, err)
}

func checkClusterInfoDuplicatedOrNot(kit *rest.Kit, bizID int64, data *types.Cluster) ccErr.CCErrorCoder {

	filterName := map[string]interface{}{
		common.BKFieldName:  *data.Name,
		common.BKAppIDField: bizID,
	}
	util.SetModOwner(filterName, kit.SupplierAccount)

	filterUid := map[string]interface{}{
		common.BKFieldName:  *data.Uid,
		common.BKAppIDField: bizID,
	}
	util.SetModOwner(filterUid, kit.SupplierAccount)

	filter := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			filterName, filterUid,
		},
	}

	count, err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("create cluster failed, name or uid duplicated, name: %s, uid: %s, rid: %s", data.Name,
			data.Uid, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, "name or uid")
	}
	return nil
}

// createCluster create cluster instance.
func createCluster(kit *rest.Kit, bizID int64, data *types.Cluster) (*types.Cluster,
	ccErr.CCErrorCoder) {

	// it is necessary to judge whether there is duplicate data here,
	// to prevent subsequent calls to coreservice directly and lack of verification.
	if err := data.ValidateCreate(); err.ErrCode != 0 {
		blog.Errorf("create cluster failed, data: %+v, err: %+v, rid: %s", data, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	if err := checkClusterInfoDuplicatedOrNot(kit, bizID, data); err != nil {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error())
	}

	// generate id field
	id, err := mongodb.Client().NextSequence(kit.Ctx, types.BKTableNameBaseCluster)
	if err != nil {
		blog.Errorf("create cluster failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	now := time.Now().Unix()
	cluster := &types.Cluster{
		ID:               int64(id),
		BizID:            bizID,
		SupplierAccount:  kit.SupplierAccount,
		Name:             data.Name,
		SchedulingEngine: data.SchedulingEngine,
		Uid:              data.Uid,
		Xid:              data.Xid,
		Version:          data.Version,
		NetworkType:      data.NetworkType,
		Region:           data.Region,
		Vpc:              data.Vpc,
		Environment:      data.Environment,
		NetWork:          data.NetWork,
		Type:             data.Type,
		Revision: table.Revision{
			CreateTime: now,
			LastTime:   now,
			Creator:    kit.User,
			Modifier:   kit.User,
		},
	}

	if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Insert(kit.Ctx, cluster); err != nil {
		blog.Errorf("create cluster failed, db insert failed, doc: %+v, err: %+v, rid: %s", cluster, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed, "cluster")
	}
	return cluster, nil
}

// BatchDeleteCluster delete clusters.
func (s *coreService) BatchDeleteCluster(ctx *rest.Contexts) {

	option := new(types.DeleteClusterOption)
	if err := ctx.DecodeInto(option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	filter := make(map[string]interface{}, 0)
	if len(option.IDs) > 0 {
		filter = map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: ctx.Kit.SupplierAccount,
			types.BKIDField: map[string]interface{}{
				common.BKDBIN: option.IDs,
			},
		}
	}
	if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Delete(ctx.Kit.Ctx, filter); err != nil {
		blog.Errorf("delete cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
