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
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *Service) AddCloudTask(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("add cloud sync task failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input[common.CreateTimeField] = time.Now()
	input = util.SetModOwner(input, ownerID)

	err := s.Logics.CreateCloudTask(ctx, input)
	if err != nil {
		blog.Errorf("create cloud sync data: %v error: %v", input, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudCreateSyncTaskFail)})
		return
	}

	result := make(map[string]interface{})
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}

func (s *Service) TaskNameCheck(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("task name check failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	condition := common.KvMap{"bk_task_name": input["bk_task_name"]}
	condition = util.SetModOwner(condition, ownerID)
	num, err := s.Instance.Table(common.BKTableNameCloudTask).Find(condition).Count(ctx)
	if err != nil {
		blog.Error("get task name [%s] failed, err: %v", input["bk_task_name"], err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     num,
	})
}

func (s *Service) DeleteCloudTask(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	taskID := req.PathParameter("taskID")
	intTaskID, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		blog.Errorf("string to int64 failed with err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	params := common.KvMap{"bk_task_id": intTaskID}

	if err := s.Instance.Table(common.BKTableNameCloudTask).Delete(ctx, params); err != nil {
		blog.Errorf("delete failed err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBDeleteFailed)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     "success",
	})
}

func (s *Service) SearchCloudTask(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	opt := mapstr.MapStr{}
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("add cloud sync task failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	page := mapstr.MapStr{}
	result := make([]meta.CloudTaskInfo, 0)
	var num uint64
	if opt["page"] != nil {
		pageM, err := mapstr.NewFromInterface(opt["page"])
		delete(opt, "page")
		if err != nil {
			blog.Error("interface convert to mapstr failed")
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
			return
		}
		page = pageM

		sort, errS := page.String("sort")
		if errS != nil {
			blog.Error("interface convert to string failed")
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
			return
		}
		limit, errL := page.Int64("limit")
		if errL != nil {
			blog.Error("interface convert to int64 failed")
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
			return
		}
		start, errStart := page.Int64("start")
		if errStart != nil {
			blog.Error("interface convert to string failed")
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
			return
		}

		errR := s.Instance.Table(common.BKTableNameCloudTask).Find(opt).Sort(sort).Start(uint64(start)).Limit(uint64(limit)).All(ctx, &result)
		if errR != nil {
			blog.Error("get failed, err: %v", errR)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}

		number, errN := s.Instance.Table(common.BKTableNameCloudTask).Find(opt).Count(ctx)
		if errN != nil {
			blog.Error("get task name [%s] failed, err: %v", errN)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}
		num = number
	} else {
		if err := s.Instance.Table(common.BKTableNameCloudTask).Find(opt).All(ctx, &result); err != nil {
			blog.Error("get failed, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}

		number, err := s.Instance.Table(common.BKTableNameCloudTask).Find(opt).Count(ctx)
		if err != nil {
			blog.Error("get task name [%s] failed, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}
		num = number
	}

	resp.WriteEntity(meta.CloudTaskSearch{
		Count: num,
		Info:  result,
	})
}

func (s *Service) UpdateCloudTask(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ctx := util.GetDBContext(context.Background(), pheader)

	data := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update cloud task failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	params := common.KvMap{"bk_task_id": data["bk_task_id"]}
	err := s.Instance.Table(common.BKTableNameCloudTask).Update(ctx, params, data)
	if nil != err {
		blog.Error("update cloud task failed, error information is %s, params:%v", err.Error(), params)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) ResourceConfirm(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("add cloud sync task failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input[common.CreateTimeField] = time.Now()
	input = util.SetModOwner(input, ownerID)

	err := s.Logics.CreateResourceConfirm(ctx, input)
	if err != nil {
		blog.Errorf("create cloud sync data: %v error: %v", input, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudConfirmCreateFail)})
		return
	}

	result := make(map[string]interface{})
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})

}

func (s *Service) SearchConfirm(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("search confirm failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	result := make([]map[string]interface{}, 0)
	var number uint64
	if opt["page"] != nil {
		page, err := mapstr.NewFromInterface(opt["page"])
		if err != nil {
			blog.Errorf("interface convert to mapstr fail, error: %v", err)
		}
		delete(opt, "page")

		sort, errS := page.String("sort")
		if errS != nil {
			blog.Error("interface convert to string failed")
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
			return
		}
		limit, errL := page.Int64("limit")
		if errL != nil {
			blog.Error("interface convert to string failed")
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
			return
		}
		start, errStart := page.Int64("start")
		if errStart != nil {
			blog.Error("interface convert to string failed")
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
			return
		}

		errR := s.Instance.Table(common.BKTableNameCloudResourceConfirm).Find(opt).Sort(sort).Start(uint64(start)).Limit(uint64(limit)).All(ctx, &result)
		if errR != nil {
			blog.Error("get failed, err: %v", errR)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}

		num, errN := s.Instance.Table(common.BKTableNameCloudResourceConfirm).Find(opt).Count(ctx)
		if errN != nil {
			blog.Error("get task name [%s] failed, err: %v", errN)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}
		number = num
	} else {
		if err := s.Instance.Table(common.BKTableNameCloudResourceConfirm).Find(opt).All(ctx, &result); err != nil {
			blog.Error("get failed, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}

		num, err := s.Instance.Table(common.BKTableNameCloudResourceConfirm).Find(opt).Count(ctx)
		if err != nil {
			blog.Error("get task name [%s] failed, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}
		number = num
	}

	resp.WriteEntity(meta.FavoriteResult{
		Count: number,
		Info:  result,
	})
}

func (s *Service) DeleteConfirm(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	resourceID := req.PathParameter("resourceID")
	intResourceID, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		blog.Errorf("string to int64 failed with err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	params := common.KvMap{"bk_resource_id": intResourceID}
	if err := s.Instance.Table(common.BKTableNameCloudResourceConfirm).Delete(ctx, params); err != nil {
		blog.Errorf("delete failed err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBDeleteFailed)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     "success",
	})
}

func (s *Service) AddSyncHistory(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("add cloud sync task failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err := s.Logics.CreateCloudHistory(ctx, input)
	if err != nil {
		blog.Errorf("create cloud history data: %v error: %v", input, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudHistoryCreateFail)})
		return
	}

	result := make(map[string]interface{})
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}

func (s *Service) SearchSyncHistory(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("add cloud sync task failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	condition := make(map[string]interface{})
	condition["bk_start_time"] = util.ConvParamsTime(opt["bk_start_time"])
	condition["bk_task_id"] = opt["bk_task_id"]
	page, err := mapstr.NewFromInterface(opt["page"])
	if err != nil {
		blog.Error("interface convert to mapstr failed")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
		return
	}

	sort, errS := page.String("sort")
	if errS != nil {
		blog.Error("interface convert to string failed")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
		return
	}
	limit, errL := page.Int64("limit")
	if errL != nil {
		blog.Error("interface convert to string failed")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
		return
	}
	start, errStart := page.Int64("start")
	if errStart != nil {
		blog.Error("interface convert to string failed")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
		return
	}

	result := make([]map[string]interface{}, 0)
	if err := s.Instance.Table(common.BKTableNameCloudSyncHistory).Find(condition).Sort(sort).Start(uint64(start)).Limit(uint64(limit)).All(ctx, &result); err != nil {
		blog.Error("get failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	num, errN := s.Instance.Table(common.BKTableNameCloudSyncHistory).Find(condition).Count(ctx)
	if errN != nil {
		blog.Error("get task name [%s] failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.FavoriteResult{
		Count: num,
		Info:  result,
	})
}

func (s *Service) AddConfirmHistory(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("add confirm history failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input["confirm_time"] = time.Now()
	err := s.Logics.CreateConfirmHistory(ctx, input)
	if err != nil {
		blog.Errorf("create cloud history data: %v error: %v", input, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudConfirmHistoryAddFail)})
		return
	}

	result := make(map[string]interface{})
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}

func (s *Service) SearchConfirmHistory(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("search confirm history failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	page, err := mapstr.NewFromInterface(opt["page"])
	delete(opt, "page")
	condition := util.ConvParamsTime(opt)
	if err != nil {
		blog.Error("interface convert to mapstr failed")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
		return
	}

	sort, errS := page.String("sort")
	if errS != nil {
		blog.Error("interface convert to string failed")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
		return
	}
	limit, errL := page.Int64("limit")
	if errL != nil {
		blog.Error("interface convert to string failed")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
		return
	}
	start, errStart := page.Int64("start")
	if errStart != nil {
		blog.Error("interface convert to string failed")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncHistorySearchFail)})
		return
	}

	result := make([]map[string]interface{}, 0)
	errR := s.Instance.Table(common.BKTableNameResourceConfirmHistory).Find(condition).Sort(sort).Start(uint64(start)).Limit(uint64(limit)).All(ctx, &result)
	if errR != nil {
		blog.Error("get failed, err: %v", errR)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	num, errN := s.Instance.Table(common.BKTableNameResourceConfirmHistory).Find(condition).Count(ctx)
	if errN != nil {
		blog.Error("get task name [%s] failed, err: %v", errN)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.FavoriteResult{
		Count: num,
		Info:  result,
	})
}
