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
		blog.Errorf("create cloud sync task failed， MarshalJSONInto error, err:%s, input:%v, rid: %v", err.Error(), data, params.ReqID)
		return nil, err
	}

	id, err := s.core.HostOperation().CreateCloudSyncTask(params, input)
	if err != nil {
		blog.Errorf("create cloud sync fail input: %v, error: %v, rid: %v", input, err, params.ReqID)
		return nil, err
	}

	return id, nil
}

func (s *coreService) CheckTaskNameUnique(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	condition := common.KvMap{common.CloudSyncTaskName: data[common.CloudSyncTaskName]}
	num, err := s.db.Table(common.BKTableNameCloudTask).Find(condition).Count(params)
	if err != nil {
		blog.Errorf("CheckTaskNameUnique fail, get task name [%s] failed, err: %v, rid: %v", data["bk_task_name"], err, params.ReqID)
		return nil, err
	}

	return num, nil
}

func (s *coreService) DeleteCloudSyncTask(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	taskID := pathParams("taskID")
	intTaskID, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		blog.Errorf("DeleteCloudSyncTask fail, taskID string to int64 failed with err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	opt := common.KvMap{common.CloudSyncTaskID: intTaskID}
	if err := s.db.Table(common.BKTableNameCloudTask).Delete(params, opt); err != nil {
		blog.Errorf("DeleteCloudSyncTask failed err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	return nil, nil
}

func (s *coreService) SearchCloudSyncTask(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	page := meta.ParsePage(data["page"])
	result := make([]meta.CloudTaskInfo, 0)
	err := s.db.Table(common.BKTableNameCloudTask).Find(data).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(params, &result)
	if err != nil {
		blog.Error("SearchCloudSyncTask failed err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	count, err := s.db.Table(common.BKTableNameCloudTask).Find(data).Count(params)
	if err != nil {
		blog.Error("SearchCloudSyncTask failed, get task count fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	return meta.CloudTaskSearch{
		Count: count,
		Info:  result,
	}, nil
}

func (s *coreService) UpdateCloudSyncTask(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	opt := common.KvMap{common.CloudSyncTaskID: data[common.CloudSyncTaskID]}
	err := s.db.Table(common.BKTableNameCloudTask).Update(params, opt, data)
	if nil != err {
		blog.Error("UpdateCloudSyncTask fail, error information is %s, params:%v, rid: %v", err.Error(), data, params.ReqID)
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
		blog.Errorf("CreateConfirm fail, input: %v error: %v, rid: %s", input, err, params.ReqID)
		return nil, err
	}

	return id, nil
}

func (s *coreService) SearchConfirm(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	page := meta.ParsePage(data["page"])
	result := make([]map[string]interface{}, 0)
	err := s.db.Table(common.BKTableNameCloudResourceConfirm).Find(data).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(params, &result)
	if err != nil {
		blog.Error("search cloud resource confirm  fail, search condition: %v err: %v, rid: %v", data, err, params.ReqID)
		return nil, err
	}

	count, err := s.db.Table(common.BKTableNameCloudResourceConfirm).Find(data).Count(params)
	if err != nil {
		blog.Error("search cloud resource confirm count fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	return meta.FavoriteResult{
		Count: count,
		Info:  result,
	}, nil
}

func (s *coreService) DeleteConfirm(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	resourceID := pathParams("taskID")
	intResourceID, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		blog.Errorf("DeleteConfirm fail, taskID string to int64 failed with err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	opt := common.KvMap{common.CloudSyncResourceConfirmID: intResourceID}
	if err := s.db.Table(common.BKTableNameCloudResourceConfirm).Delete(params, opt); err != nil {
		blog.Errorf("DeleteConfirm failed err: %v, rid: %v", err, params.ReqID)
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
		blog.Errorf("create cloud history fail, input: %v error: %v, rid: %v", input, err, params.ReqID)
		return nil, err
	}

	return id, nil
}

func (s *coreService) SearchSyncHistory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	condition := make(map[string]interface{})
	condition["bk_start_time"] = util.ConvParamsTime(data["bk_start_time"])
	condition["bk_task_id"] = data["bk_task_id"]
	page := meta.ParsePage(data["page"])

	result := make([]map[string]interface{}, 0)
	if err := s.db.Table(common.BKTableNameCloudSyncHistory).Find(condition).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(params, &result); err != nil {
		blog.Error("search cloud sync history fail, err: %v, input: %v, rid: %v", err, data, params.ReqID)
		return nil, err
	}

	count, err := s.db.Table(common.BKTableNameCloudSyncHistory).Find(condition).Count(params)
	if err != nil {
		blog.Error("search cloud sync history count fail, err: %v, rid: %v", err, params.ReqID)
		return nil, err
	}

	return meta.FavoriteResult{
		Count: count,
		Info:  result,
	}, nil
}

func (s *coreService) CreateConfirmHistory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	data[common.CloudSyncConfirmTime] = time.Now()
	id, err := s.core.HostOperation().CreateConfirmHistory(params, data)
	if err != nil {
		blog.Errorf("CreateConfirmHistory fail, input: %v error: %v, rid: %v", data, err, params.ReqID)
		return nil, err
	}

	return id, nil
}

func (s *coreService) SearchConfirmHistory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	page := meta.ParsePage(data["page"])
	delete(data, "page")
	condition := util.ConvParamsTime(data)

	result := make([]map[string]interface{}, 0)
	err := s.db.Table(common.BKTableNameResourceConfirmHistory).Find(condition).Sort(page.Sort).Start(uint64(page.Start)).Limit(uint64(page.Limit)).All(params, &result)
	if err != nil {
		blog.Error("search resource confirm history fail, err: %v, input: %v, rid: %v", err, data, params.ReqID)
		return nil, err
	}

	num, err := s.db.Table(common.BKTableNameResourceConfirmHistory).Find(condition).Count(params)
	if err != nil {
		blog.Error("SearchConfirmHistory count fail, err: %v, input: %v, rid: %v", err, data, params.ReqID)
		return nil, err
	}

	return meta.FavoriteResult{
		Count: num,
		Info:  result,
	}, nil
}
