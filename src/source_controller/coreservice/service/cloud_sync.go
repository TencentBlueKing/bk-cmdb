/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *coreService) CreateCloudSyncTask(ctx *rest.Contexts) {
	input := new(meta.CloudTaskList)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	id, err := s.core.HostOperation().CreateCloudSyncTask(ctx.Kit, input)
	if err != nil {
		blog.Errorf("create cloud sync fail input: %v, error: %v, rid: %v", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(id)
}

func (s *coreService) CheckTaskNameUnique(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	condition := common.KvMap{common.CloudSyncTaskName: data[common.CloudSyncTaskName]}
	num, err := s.db.Table(common.BKTableNameCloudTask).Find(condition).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("CheckTaskNameUnique fail, get task name [%s] failed, err: %v, rid: %v", data["bk_task_name"], err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(num)
}

func (s *coreService) DeleteCloudSyncTask(ctx *rest.Contexts) {
	taskID := ctx.Request.PathParameter("taskID")
	intTaskID, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		blog.Errorf("DeleteCloudSyncTask fail, taskID string to int64 failed with err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	opt := common.KvMap{common.CloudSyncTaskID: intTaskID}
	if err := s.db.Table(common.BKTableNameCloudTask).Delete(ctx.Kit.Ctx, opt); err != nil {
		blog.Errorf("DeleteCloudSyncTask failed err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) SearchCloudSyncTask(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	page := meta.ParsePage(data["page"])
	result := make([]meta.CloudTaskInfo, 0)
	err := s.db.Table(common.BKTableNameCloudTask).Find(data).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(ctx.Kit.Ctx, &result)
	if err != nil {
		blog.Error("SearchCloudSyncTask failed err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	count, err := s.db.Table(common.BKTableNameCloudTask).Find(data).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Error("SearchCloudSyncTask failed, get task count fail, err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(meta.CloudTaskSearch{
		Count: count,
		Info:  result,
	})
}

func (s *coreService) UpdateCloudSyncTask(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	opt := common.KvMap{common.CloudSyncTaskID: data[common.CloudSyncTaskID]}
	err := s.db.Table(common.BKTableNameCloudTask).Update(ctx.Kit.Ctx, opt, data)
	if nil != err {
		blog.Error("UpdateCloudSyncTask fail, error information is %s, ctx:%v, rid: %v", err.Error(), data, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) CreateConfirm(ctx *rest.Contexts) {
	input := new(meta.ResourceConfirm)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	input.CreateTime = time.Now()
	id, err := s.core.HostOperation().CreateResourceConfirm(ctx.Kit, input)
	if err != nil {
		blog.Errorf("CreateConfirm fail, input: %v error: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(id)
}

func (s *coreService) SearchConfirm(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	page := meta.ParsePage(data["page"])
	result := make([]map[string]interface{}, 0)
	err := s.db.Table(common.BKTableNameCloudResourceConfirm).Find(data).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(ctx.Kit.Ctx, &result)
	if err != nil {
		blog.Error("search cloud resource confirm  fail, search condition: %v err: %v, rid: %v", data, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	count, err := s.db.Table(common.BKTableNameCloudResourceConfirm).Find(data).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Error("search cloud resource confirm count fail, err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(meta.FavoriteResult{
		Count: count,
		Info:  result,
	})
}

func (s *coreService) DeleteConfirm(ctx *rest.Contexts) {
	resourceID := ctx.Request.PathParameter("taskID")
	intResourceID, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		blog.Errorf("DeleteConfirm fail, taskID string to int64 failed with err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	opt := common.KvMap{common.CloudSyncResourceConfirmID: intResourceID}
	if err := s.db.Table(common.BKTableNameCloudResourceConfirm).Delete(ctx.Kit.Ctx, opt); err != nil {
		blog.Errorf("DeleteConfirm failed err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *coreService) CreateSyncHistory(ctx *rest.Contexts) {
	input := new(meta.CloudHistory)
	if err := ctx.DecodeInto(&input); nil != err {
		blog.Errorf("create cloud sync history failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	id, err := s.core.HostOperation().CreateCloudSyncHistory(ctx.Kit, input)
	if err != nil {
		blog.Errorf("create cloud history fail, input: %v error: %v, rid: %v", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(id)
}

func (s *coreService) SearchSyncHistory(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	condition := make(map[string]interface{})
	condition["bk_start_time"] = util.ConvParamsTime(data["bk_start_time"])
	condition["bk_task_id"] = data["bk_task_id"]
	page := meta.ParsePage(data["page"])

	result := make([]map[string]interface{}, 0)
	if err := s.db.Table(common.BKTableNameCloudSyncHistory).Find(condition).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(ctx.Kit.Ctx, &result); err != nil {
		blog.Error("search cloud sync history fail, err: %v, input: %v, rid: %v", err, data, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	count, err := s.db.Table(common.BKTableNameCloudSyncHistory).Find(condition).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Error("search cloud sync history count fail, err: %v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(meta.FavoriteResult{
		Count: count,
		Info:  result,
	})
}

func (s *coreService) CreateConfirmHistory(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	data[common.CloudSyncConfirmTime] = time.Now()
	id, err := s.core.HostOperation().CreateConfirmHistory(ctx.Kit, data)
	if err != nil {
		blog.Errorf("CreateConfirmHistory fail, input: %v error: %v, rid: %v", data, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(id)
}

func (s *coreService) SearchConfirmHistory(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	page := meta.ParsePage(data["page"])
	delete(data, "page")
	condition := util.ConvParamsTime(data)

	result := make([]map[string]interface{}, 0)
	err := s.db.Table(common.BKTableNameResourceConfirmHistory).Find(condition).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(ctx.Kit.Ctx, &result)
	if err != nil {
		blog.Error("search resource confirm history fail, err: %v, input: %v, rid: %v", err, data, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	num, err := s.db.Table(common.BKTableNameResourceConfirmHistory).Find(condition).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Error("SearchConfirmHistory count fail, err: %v, input: %v, rid: %v", err, data, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(meta.FavoriteResult{
		Count: num,
		Info:  result,
	})
}
