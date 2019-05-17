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
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
)

// CloudAddTask create cloud sync task
func (s *Service) AddCloudTask(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	taskList := new(meta.CloudTaskList)
	if err := json.NewDecoder(req.Request.Body).Decode(taskList); err != nil {
		blog.Errorf("add task failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	taskList.User = srvData.user

	if err := srvData.lgc.AddCloudTask(srvData.ctx, taskList); err != nil {
		blog.Errorf("add task failed with err: %v, rid: %s", err.Error(), srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCloudSyncCreateFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) DeleteCloudTask(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	taskID := req.PathParameter("taskID")
	_, err := s.CoreAPI.HostController().Cloud().DeleteCloudTask(srvData.ctx, srvData.header, taskID)

	retData := make(map[string]interface{})
	if err != nil {
		retData["errors:"] = err
	}

	resp.WriteEntity(meta.NewSuccessResp(retData))
}

func (s *Service) SearchCloudTask(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("search task , but decode body failed, err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	response, err := s.CoreAPI.HostController().Cloud().SearchCloudTask(srvData.ctx, srvData.header, opt)
	if err != nil {
		blog.Errorf("search %v failed, err: %v, rid: %s", opt["bk_task_name"], err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCloudGetTaskFail)})
	}

	resp.WriteEntity(meta.NewSuccessResp(response))
}

func (s *Service) UpdateCloudTask(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	data := make(mapstr.MapStr, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update task failed with decode body err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	// TaskName Uniqueness check
	response, err := s.CoreAPI.HostController().Cloud().TaskNameCheck(srvData.ctx, srvData.header, data)
	if err != nil {
		blog.Debug("task name unique check fail, error: %v, rid: %s", err, srvData.rid)
		return
	}

	if response.Data != 0.0 {
		blog.Errorf("update task failed, task name already exits. rid: %s", srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCloudTaskNameAlreadyExist)})
		return
	}

	status, err := data.Bool("bk_status")
	if err != nil {
		blog.Errorf("interface convert to bool fail, err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}

	if status {
		// 开启同步状态下，update；先关闭同步，更新数据后，再开启同步
		data["bk_status"] = false
		if _, err := s.CoreAPI.HostController().Cloud().UpdateCloudTask(srvData.ctx, srvData.header, data); err != nil {
			blog.Errorf("update task failed with decode body err: %v, rid: %s", err, srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommDBUpdateFailed)})
			return
		}

		if err := srvData.lgc.FrontEndSyncSwitch(srvData.ctx, data); err != nil {
			blog.Errorf("stop cloud sync fail, err: %v, rid: %s", err, srvData.rid)
			return
		}

		data["bk_status"] = true
		if _, err := s.CoreAPI.HostController().Cloud().UpdateCloudTask(srvData.ctx, srvData.header, data); err != nil {
			blog.Errorf("update task failed with decode body err: %v, rid: %s", err, srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommDBUpdateFailed)})
			return
		}

		if err := srvData.lgc.FrontEndSyncSwitch(srvData.ctx, data); err != nil {
			blog.Errorf("stop cloud sync fail, err: %v, rid: %s", err, srvData.rid)
			return
		}
	} else {
		if _, err := s.CoreAPI.HostController().Cloud().UpdateCloudTask(srvData.ctx, srvData.header, data); err != nil {
			blog.Errorf("update task failed with decode body err: %v, rid: %s", err, srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommDBUpdateFailed)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) StartCloudSync(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	opt := make(map[string]interface{}, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("update cloud task failed, err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if _, err := s.CoreAPI.HostController().Cloud().UpdateCloudTask(srvData.ctx, srvData.header, opt); err != nil {
		blog.Errorf("update task failed with decode body, %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}

	isRequired := make([]string, 0)
	if _, ok := opt["bk_status"]; ok {
		delete(opt, "bk_status")
	} else {
		isRequired = append(isRequired, "bk_status is required.")
	}

	if _, ok := opt["bk_task_name"]; !ok {
		isRequired = append(isRequired, "bk_task_name is required.")
	}

	if len(isRequired) > 0 {
		blog.Errorf("%v", isRequired)
		resp.WriteEntity(meta.NewSuccessResp(isRequired))
		return
	}

	delete(opt, "bk_task_name")

	if err := srvData.lgc.FrontEndSyncSwitch(srvData.ctx, opt); err != nil {
		blog.Errorf("start cloud sync fail, err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCloudSyncStartFail)})
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) ResourceConfirm(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	resourceIDMap := make(map[string][]int64)
	if err := json.NewDecoder(req.Request.Body).Decode(&resourceIDMap); err != nil {
		blog.Errorf("resource confirm failed with decode body err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	resourceIDs := resourceIDMap["bk_resource_id"]
	cloudHostInfo := make([]mapstr.MapStr, 0)
	for _, id := range resourceIDs {
		opt := make(map[string]interface{})
		opt["bk_resource_id"] = id
		response, err := s.CoreAPI.HostController().Cloud().SearchConfirm(srvData.ctx, srvData.header, opt)
		if err != nil {
			blog.Errorf("get resourceID %v confirm list failed. err: %v, rid: %s", id, err, srvData.rid)
			continue
		}
		if response.Count > 0 {
			cloudHostInfo = append(cloudHostInfo, response.Info[0])
		}
	}

	AddHostList := make([]mapstr.MapStr, 0)
	updateHostList := make([]mapstr.MapStr, 0)
	for _, hostInfo := range cloudHostInfo {
		addConfirm, ok := hostInfo["bk_confirm"].(bool)
		if !ok {
			blog.Errorf("interface convert to bool fail")
			continue
		}
		if addConfirm {
			AddHostList = append(AddHostList, hostInfo)
		}

		attrConfirm, ok := hostInfo["bk_attr_confirm"].(bool)
		if !ok {
			blog.Errorf("interface convert to bool fail")
			continue
		}
		if attrConfirm {
			updateHostList = append(updateHostList, hostInfo)
		}
	}

	if len(AddHostList) > 0 {
		blog.Info("new add confirmed")
		err := srvData.lgc.AddCloudHosts(srvData.ctx, AddHostList)
		if err != nil {
			blog.Errorf("add cloud host failed, err: %v, rid: %s", err, srvData.rid)
			srvData.ccErr.Error(1110080)
		}
	}

	if len(updateHostList) > 0 {
		err := srvData.lgc.UpdateCloudHosts(srvData.ctx, updateHostList)
		if err != nil {
			blog.Errorf("update cloud hosts failed, err: %v, rid: %s", err, srvData.rid)
			srvData.ccErr.Error(1110003)
		}
	}

	// After resource confirmation, delete the items from table cc_CloudResourceSync
	for _, id := range resourceIDs {
		_, errD := srvData.lgc.CoreAPI.HostController().Cloud().DeleteConfirm(srvData.ctx, srvData.header, id)
		if errD != nil {
			blog.Errorf("delete resource confirm failed with err: %v, rid: %s", errD, srvData.rid)
			continue
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) SearchConfirm(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("resource confirm failed with decode body err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	response, err := s.CoreAPI.HostController().Cloud().SearchConfirm(srvData.ctx, srvData.header, opt)
	if err != nil {
		blog.Errorf("search confirm instance failed, rid: %s", srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCloudGetConfirmFail)})
	}

	resp.WriteEntity(meta.NewSuccessResp(response))
}

func (s *Service) SearchAccount(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("search task , but decode body failed, err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	response, err := s.CoreAPI.HostController().Cloud().SearchCloudTask(srvData.ctx, srvData.header, opt)

	if err != nil {
		blog.Errorf("search %v failed, err: %v, rid: %s", opt["bk_task_name"], err, srvData.rid)
		resp.WriteEntity(meta.NewSuccessResp(err))
	}

	secretID := response.Info[0].SecretID
	secretKey := response.Info[0].SecretKey

	// decode secretKey
	decodeBytes, errDecode := base64.StdEncoding.DecodeString(secretKey)
	if errDecode != nil {
		blog.Errorf("Base64 decode secretKey failed. rid: %s", srvData.rid)
		resp.WriteEntity(meta.NewSuccessResp(errDecode))
	}
	secretKeyOrigin := string(decodeBytes)

	result := make(map[string]interface{}, 0)
	result["ID"] = secretID
	result["Key"] = secretKeyOrigin

	resp.WriteEntity(meta.NewSuccessResp(response))
}

func (s *Service) CloudSyncHistory(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("search task , but decode body failed, err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	response, err := s.CoreAPI.HostController().Cloud().SearchSyncHistory(srvData.ctx, srvData.header, opt)
	if err != nil {
		blog.Errorf("search cloud sync history failed, err: %v, rid: %s", err, srvData.rid)
		resp.WriteEntity(meta.NewSuccessResp(err))
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(response))
}

func (s *Service) AddConfirmHistory(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	resourceIDMap := make(map[string][]int64)
	if err := json.NewDecoder(req.Request.Body).Decode(&resourceIDMap); err != nil {
		blog.Errorf("resource confirm failed with decode body err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	resourceIDs := resourceIDMap["bk_resource_id"]

	for _, id := range resourceIDs {
		condition := make(map[string]interface{})
		condition["bk_resource_id"] = id
		response, err := s.CoreAPI.HostController().Cloud().SearchConfirm(srvData.ctx, srvData.header, condition)
		if err != nil {
			blog.Errorf("search confirm instance failed, rid: %s", srvData.rid)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCloudGetConfirmFail)})
		}

		if response.Count > 0 {
			opt := response.Info[0]
			if _, err := s.CoreAPI.HostController().Cloud().AddConfirmHistory(srvData.ctx, srvData.header, opt); err != nil {
				blog.Errorf("add confirm history failed, rid: %s", srvData.rid)
				resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCloudAddConfirmHistoryFail)})
			}
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) SearchConfirmHistory(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	opt := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&opt); err != nil {
		blog.Errorf("resource confirm failed with decode body err: %v, rid: %s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	response, err := s.CoreAPI.HostController().Cloud().SearchConfirmHistory(context.Background(), srvData.header, opt)
	if err != nil {
		blog.Errorf("search confirm history failed, rid: %s", srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCloudGetConfirmHistoryFail)})
	}

	resp.WriteEntity(meta.NewSuccessResp(response))
}
