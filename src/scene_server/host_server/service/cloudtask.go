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
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	opt := new(meta.DeleteCloudTask)
	if err := json.NewDecoder(req.Request.Body).Decode(opt); err != nil {
		blog.Errorf("delete task , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	_, err := s.CoreAPI.HostController().Cloud().DeleteCloudTask(context.Background(), pheader, opt.TaskID)

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

	response, err := s.CoreAPI.HostController().Cloud().SearchCloudTask(context.Background(), pheader, opt)

	if err != nil {
		blog.Errorf("search %v failed, err: %v", opt["bk_task_name"], err)
		resp.WriteEntity(meta.NewSuccessResp(err))
	} else {
		resp.WriteEntity(meta.NewSuccessResp(response.Data))
	}
}

func (s *Service) UpdateCloudTask(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	updateTask := new(meta.CloudTaskList)
	if err := json.NewDecoder(req.Request.Body).Decode(updateTask); err != nil {
		blog.Errorf("update task failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	_, err := s.CoreAPI.HostController().Cloud().UpdateCloudTask(context.Background(), pheader, updateTask)
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

	taskList, errMap := mapstr.NewFromInterface(response.Data[0])
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

	resourceConfirmMap := make(map[string][]mapstr.MapStr)
	if err := json.NewDecoder(req.Request.Body).Decode(&resourceConfirmMap); err != nil {
		blog.Errorf("resource confirm failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	resourceConfirm := resourceConfirmMap["bk_resource"]

	resourceIDs := make([]interface{}, 0)
	AddHostList := make([]mapstr.MapStr, 0)
	updateHostList := make([]mapstr.MapStr, 0)
	for _, hostInfo := range resourceConfirm {
		resourceID := hostInfo["bk_resource_id"]
		resourceIDs = append(resourceIDs, resourceID)

		addConfirm, ok := hostInfo["bk_confirm"].(bool)
		if !ok {
			blog.Errorf("interface convert to bool fail")
			//resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
			return
		}
		if addConfirm {
			AddHostList = append(AddHostList, hostInfo)
		}

		attrConfirm, ok := hostInfo["bk_attr_confirm"].(bool)
		if !ok {
			blog.Errorf("interface convert to bool fail")
			return
		}
		if attrConfirm {
			updateHostList = append(updateHostList, hostInfo)
		}
	}

	if len(AddHostList) > 0 {
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
		resourceId, ok := id.(int64)
		if !ok {
			blog.Errorf("interface convert to int64 failed")
			return
		}
		_, errD := s.Logics.CoreAPI.HostController().Cloud().DeleteConfirm(context.Background(), pheader, resourceId)
		if errD != nil {
			blog.Errorf("delete resource confirm failed with err: %v", errD)
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
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

	secretID := response.Data[0]["bk_secret_id"]
	secretKey := response.Data[0]["bk_secret_key"]
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

	response, err := s.CoreAPI.HostController().Cloud().SearchHistory(context.Background(), pheader)
	if err != nil {
		blog.Errorf("search cloud sync history failed, err: %v", err)
		resp.WriteEntity(meta.NewSuccessResp(err))
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(response))
}
