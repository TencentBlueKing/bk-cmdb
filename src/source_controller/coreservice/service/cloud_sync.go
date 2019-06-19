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
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) CreateCloudSyncTask(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := new(meta.CloudTaskList)
	if err := data.MarshalJSONInto(input); nil != err {
		blog.Errorf("create cloud sync task failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data)
		return nil, err
	}

	id, err := s.core.HostOperation().CreateCloudSyncTask(params, input)
	if err != nil {
		blog.Errorf("create cloud sync data: %v error: %v", input, err)
		return nil, err
	}

	return id, nil
}

func (s *coreService) TaskNameCheck(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := make(map[string]interface{})
	if err := data.MarshalJSONInto(input); nil != err {
		blog.Errorf("task name check failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data)
		return nil, err
	}

	condition := common.KvMap{common.CloudSyncTaskName: input[common.CloudSyncTaskName]}
	num, err := s.db.Table(common.BKTableNameCloudTask).Find(condition).Count(params)
	if err != nil {
		blog.Error("get task name [%s] failed, err: %v", input["bk_task_name"], err)
		return nil, err
	}

	return num, nil
}

func (s *coreService) DeleteCloudTask(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	taskID := pathParams("id")
	intTaskID, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		blog.Errorf("string to int64 failed with err: %v", err)
		return nil, err
	}

	opt := common.KvMap{common.CloudSyncTaskID: intTaskID}
	if err := s.db.Table(common.BKTableNameCloudTask).Delete(params, opt); err != nil {
		blog.Errorf("delete cloud sync task failed err: %v", err)
		return nil, err
	}

	return nil, nil
}

func (s *coreService) SearchCloudTask(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := mapstr.MapStr{}
	if err := data.MarshalJSONInto(&input); nil != err {
		blog.Errorf("search cloud sync task failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data)
		return nil, err
	}

	page := meta.ParsePage(input["page"])
	result := make([]meta.CloudTaskInfo, 0)
	err := s.db.Table(common.BKTableNameCloudTask).Find(input).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(params, &result)
	if err != nil {
		blog.Error("get failed, err: %v", err)
		return nil, err
	}

	count, errN := s.db.Table(common.BKTableNameCloudTask).Find(input).Count(params)
	if errN != nil {
		blog.Error("get task name [%s] failed, err: %v", errN)
		return nil, err
	}

	return meta.CloudTaskSearch{
		Count: count,
		Info:  result,
	}, nil
}

func (s *coreService) UpdateCloudTask(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := mapstr.MapStr{}
	if err := data.MarshalJSONInto(&input); nil != err {
		blog.Errorf("update cloud sync task failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), input, params.ReqID)
		return nil, err
	}

	opt := common.KvMap{common.CloudSyncTaskID: data[common.CloudSyncTaskID]}
	err := s.db.Table(common.BKTableNameCloudTask).Update(params, opt, input)
	if nil != err {
		blog.Error("update cloud task failed, error information is %s, params:%v", err.Error(), params)
		return nil, err
	}

	return nil, nil
}

func (s *coreService) CreateConfirm(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := new(meta.ResourceConfirm)
	if err := data.MarshalJSONInto(input); nil != err {
		blog.Errorf("create resource confirm failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), input, params.ReqID)
		return nil, err
	}

	input.CreateTime = time.Now()
	id, err := s.core.HostOperation().CreateResourceConfirm(params, input)
	if err != nil {
		blog.Errorf("create cloud sync data: %v error: %v", input, err)
		return nil, err
	}

	return id, nil
}

func (s *coreService) SearchConfirm(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := mapstr.MapStr{}
	if err := data.MarshalJSONInto(&input); nil != err {
		blog.Errorf("search confirm failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data)
		return nil, err
	}

	page := meta.ParsePage(input["page"])
	result := make([]map[string]interface{}, 0)
	err := s.db.Table(common.BKTableNameCloudResourceConfirm).Find(input).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(params, &result)
	if err != nil {
		blog.Error("search cloud resource confirm %v", err)
		return nil, err
	}

	count, err := s.db.Table(common.BKTableNameCloudResourceConfirm).Find(input).Count(params)
	if err != nil {
		blog.Error("get cloud resource confirm fail, err: %v", err)
		return nil, err
	}

	return meta.FavoriteResult{
		Count: count,
		Info:  result,
	}, nil
}

func (s *coreService) DeleteConfirm(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	resourceID := pathParams("id")
	intResourceID, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		blog.Errorf("string to int64 failed with err: %v", err)
		return nil, err
	}

	opt := common.KvMap{common.CloudSyncResourceConfirmID: intResourceID}
	if err := s.db.Table(common.BKTableNameCloudResourceConfirm).Delete(params, opt); err != nil {
		blog.Errorf("delete failed err: %v", err)
		return nil, err
	}

	return nil, nil
}

func (s *coreService) CreateSyncHistory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := new(meta.CloudHistory)
	if err := data.MarshalJSONInto(&input); nil != err {
		blog.Errorf("create cloud sync history failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), input, params.ReqID)
		return nil, err
	}

	id, err := s.core.HostOperation().CreateCloudSyncHistory(params, input)
	if err != nil {
		blog.Errorf("create cloud history data: %v error: %v", input, err)
		return nil, err
	}

	return id, nil
}

func (s *coreService) SearchSyncHistory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := make(map[string]interface{})
	if err := data.MarshalJSONInto(&input); nil != err {
		blog.Errorf("search sync history failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data)
		return nil, err
	}

	condition := make(map[string]interface{})
	condition["bk_start_time"] = util.ConvParamsTime(input["bk_start_time"])
	condition["bk_task_id"] = input["bk_task_id"]
	page := meta.ParsePage(input["page"])

	result := make([]map[string]interface{}, 0)
	if err := s.db.Table(common.BKTableNameCloudSyncHistory).Find(condition).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(params, &result); err != nil {
		blog.Error("search cloud sync history fail, err: %v, input: %v", err, input)
		return nil, err
	}

	num, err := s.db.Table(common.BKTableNameCloudSyncHistory).Find(condition).Count(params)
	if err != nil {
		blog.Error("get cloud sync history count fail, err: %v", err)
		return nil, err
	}

	return meta.FavoriteResult{
		Count: num,
		Info:  result,
	}, nil
}

func (s *coreService) CreateConfirmHistory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := mapstr.MapStr{}
	if err := data.MarshalJSONInto(&input); nil != err {
		blog.Errorf("create confirm history failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), input, params.ReqID)
		return nil, err
	}

	input[common.CloudSyncConfirmTime] = time.Now()
	id, err := s.core.HostOperation().CreateConfirmHistory(params, input)
	if err != nil {
		blog.Errorf("create cloud history data: %v error: %v", input, err)
		return nil, err
	}

	return id, nil
}

func (s *coreService) SearchConfirmHistory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := make(map[string]interface{})
	if err := data.MarshalJSONInto(&input); nil != err {
		blog.Errorf("search confirm history failed， MarshalJSONInto error, err:%s,input:%v,rid:%s", err.Error(), data)
		return nil, err
	}

	page := meta.ParsePage(input["page"])
	delete(input, "page")
	condition := util.ConvParamsTime(input)

	result := make([]map[string]interface{}, 0)
	err := s.db.Table(common.BKTableNameResourceConfirmHistory).Find(condition).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(params, &result)
	if err != nil {
		blog.Error("search resource confirm history fail, err: %v, input: %v", err, input)
		return nil, err
	}

	num, err := s.db.Table(common.BKTableNameResourceConfirmHistory).Find(condition).Count(params)
	if err != nil {
		blog.Error("get resource confirm count fail, err: %v, input: %v", err, input)
		return nil, err
	}

	return meta.FavoriteResult{
		Count: num,
		Info:  result,
	}, nil
}
