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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"net/http"
)

//CloudAddTask create cloud sync task
func (s *Service) AddCloudTask(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	taskList := new(meta.CloudTaskList)
	if err := json.NewDecoder(req.Request.Body).Decode(taskList); err != nil {
		blog.Errorf("add task failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	errString, err := s.Logics.AddCloudTask(taskList, pheader)
	if err != nil {
		blog.Errorf("add task failed with err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCloudSyncCreateFail)})
		return
	}

	retData := make(map[string]interface{})
	if errString != "" {
		retData["info"] = errString
	}
	resp.WriteEntity(meta.NewSuccessResp(retData))
}

func (s *Service) DeleteCloudTask(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header

	taskID := req.PathParameter("taskID")
	blog.Debug("delete taskID: %v", taskID)
	_, err := s.CoreAPI.HostController().Cloud().DeleteCloudTask(context.Background(), pheader, taskID)

	retData := make(map[string]interface{})
	if err != nil {
		retData["errors:"] = err
	}

	resp.WriteEntity(meta.NewSuccessResp(retData))
}

func (s *Service) SearchCloudTask(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("search task , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	//blog.Debug("opt: %v", opt)
	response, err := s.CoreAPI.HostController().Cloud().SearchCloudTask(context.Background(), pheader, opt)
	if err != nil {
		blog.Errorf("search %v failed, err: %v", opt["bk_task_name"], err)
		resp.WriteEntity(meta.NewSuccessResp(err))
	} else {
		resp.WriteEntity(meta.NewSuccessResp(response))
	}
}

func (s *Service) UpdateCloudTask(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	data := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update task failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	_, err := s.CoreAPI.HostController().Cloud().UpdateCloudTask(context.Background(), pheader, data)
	if err != nil {
		blog.Errorf("update task failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) StartCloudSync(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	opt := make(map[string]interface{}, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("update cloud task failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	_, errUpdateTask := s.CoreAPI.HostController().Cloud().UpdateCloudTask(context.Background(), pheader, opt)
	if errUpdateTask != nil {
		blog.Errorf("update task failed with decode body")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}

	blog.Debug("startSync, opt: %v", opt)
	/*
	* 这个函数可以是ESB接口，接收json:{"bk_task_name": "test", "bk_status": true}开启一个同步
	* 也可以是前端发回一个json:{"bk_task_name": "test", "bk_status": true}，控制一个同步任务的开和关
	 */

	isRequired := make([]string, 0)
	status, ok := opt["bk_status"]
	if ok {
		delete(opt, "bk_status")
	} else {
		isRequired = append(isRequired, "bk_status is required.")
	}

	if _, oK := opt["bk_task_name"]; !oK {
		isRequired = append(isRequired, "bk_task_name is required.")
	}

	if len(isRequired) > 0 {
		blog.Errorf("%v", isRequired)
		resp.WriteEntity(meta.NewSuccessResp(isRequired))
	}

	response, err := s.CoreAPI.HostController().Cloud().SearchCloudTask(context.Background(), pheader, opt)
	if err != nil {
		blog.Errorf("search %v failed, err: %v", opt["bk_task_name"], err)
		resp.WriteEntity(meta.NewSuccessResp(err))
	}

	taskList, errMap := mapstr.NewFromInterface(response.Info[0])
	if errMap != nil {
		blog.Errorf("interface convert to Mapstr failed.")
		resp.WriteEntity(meta.NewSuccessResp(errMap))
	}

	taskList["bk_status"] = status

	errSync := s.Logics.CloudTaskSync(taskList, pheader)
	if errSync != nil {
		blog.Errorf("execute CloudTaskSync failed. err: %v", errSync)
		resp.WriteEntity(meta.NewSuccessResp(errSync))
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) ResourceConfirm(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	resourceIDMap := make(map[string][]int64)
	if err := json.NewDecoder(req.Request.Body).Decode(&resourceIDMap); err != nil {
		blog.Errorf("resource confirm failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	resourceIDs := resourceIDMap["bk_resource_id"]
	blog.Debug("resourceIDs: %v", resourceIDs)

	cloudHostInfo := make([]mapstr.MapStr, 0)
	for _, id := range resourceIDs {
		opt := make(map[string]interface{})
		opt["bk_resource_id"] = id
		response, err := s.CoreAPI.HostController().Cloud().SearchConfirm(context.Background(), pheader, opt)
		if err != nil {
			blog.Errorf("get resourceID %v confirm list failed. err: %v", id, err)
			break
		}

		cloudHostInfo = append(cloudHostInfo, response.Info[0])
	}

	AddHostList := make([]mapstr.MapStr, 0)
	updateHostList := make([]mapstr.MapStr, 0)
	for _, hostInfo := range cloudHostInfo {
		addConfirm, ok := hostInfo["bk_confirm"].(bool)
		if !ok {
			blog.Errorf("interface convert to bool fail")
			break
		}
		if addConfirm {
			AddHostList = append(AddHostList, hostInfo)
		}

		attrConfirm, ok := hostInfo["bk_attr_confirm"].(bool)
		if !ok {
			blog.Errorf("interface convert to bool fail")
			break
		}
		if attrConfirm {
			updateHostList = append(updateHostList, hostInfo)
		}
	}

	if len(AddHostList) > 0 {
		blog.Info("new add confirmed")
		err := s.Logics.AddCloudHosts(pheader, AddHostList)
		if err != nil {
			blog.Errorf("add cloud host failed, err: %v , err")
			return
		}

	}

	if len(updateHostList) > 0 {
		err := s.Logics.UpdateCloudHosts(pheader, updateHostList)
		if err != nil {
			blog.Errorf("update cloud hosts failed, err: %v", err)
		}
	}

	//After resource confirmation, delete the items from table cc_CloudResourceSync
	for _, id := range resourceIDs {
		_, errD := s.Logics.CoreAPI.HostController().Cloud().DeleteConfirm(context.Background(), pheader, id)
		if errD != nil {
			blog.Errorf("delete resource confirm failed with err: %v", errD)
			break
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) SearchConfirm(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("resource confirm failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	response, err := s.CoreAPI.HostController().Cloud().SearchConfirm(context.Background(), pheader, opt)
	blog.Debug("search confirm response: %v", response)
	if err != nil {
		blog.Errorf("search %v failed.")
		resp.WriteEntity(meta.NewSuccessResp(err))
	} else {
		resp.WriteEntity(meta.NewSuccessResp(response))
	}
}

func (s *Service) SearchAccount(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("search task , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	response, err := s.CoreAPI.HostController().Cloud().SearchCloudTask(context.Background(), pheader, opt)

	if err != nil {
		blog.Errorf("search %v failed, err: %v", opt["bk_task_name"], err)
		resp.WriteEntity(meta.NewSuccessResp(err))
	}

	secretID := response.Info[0]["bk_secret_id"]
	secretKey := response.Info[0]["bk_secret_key"]
	secretKeyStr, ok := secretKey.(string)
	if !ok {
		blog.Errorf("interface convert to string failed.")
		resp.WriteEntity(meta.NewSuccessResp(ok))
	}
	//decode secretKey
	decodeBytes, errDecode := base64.StdEncoding.DecodeString(secretKeyStr)
	if errDecode != nil {
		blog.Errorf("Base64 decode secretKey failed.")
		resp.WriteEntity(meta.NewSuccessResp(errDecode))
	}
	secretKeyOrigin := string(decodeBytes)

	result := make(map[string]interface{}, 0)
	result["ID"] = secretID
	result["Key"] = secretKeyOrigin

	resp.WriteEntity(meta.NewSuccessResp(result))
}

func (s *Service) CloudSyncHistory(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	taskID := req.PathParameter("taskID")
	blog.Debug("search cloud history taskID: %v", taskID)

	response, err := s.CoreAPI.HostController().Cloud().SearchHistory(context.Background(), pheader, taskID)
	if err != nil {
		blog.Errorf("search cloud sync history failed, err: %v", err)
		resp.WriteEntity(meta.NewSuccessResp(err))
		return
	}

	//blog.Debug("search sync history, response: %v", response)
	resp.WriteEntity(meta.NewSuccessResp(response))
}
