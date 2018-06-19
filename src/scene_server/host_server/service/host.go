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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"
	"github.com/emicklei/go-restful"
)

type AppResult struct {
	Result  bool        `json:result`
	Code    int         `json:code`
	Message interface{} `json:message`
	Data    DataInfo    `json:data`
}

type DataInfo struct {
	Count int                      `json:count`
	Info  []map[string]interface{} `json:info`
}

func (s *Service) DeleteHostBatch(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))
	ownerID, user := util.GetOwnerIDAndUser(req.Request.Header)

	opt := new(meta.DeleteHostBatchOpt)
	if err := json.NewDecoder(req.Request.Body).Decode(opt); err != nil {
		blog.Errorf("delete host batch , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	hostIDArr := strings.Split(opt.HostID, ",")
	var iHostIDArr []int64
	for _, i := range hostIDArr {
		iHostID, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			blog.Errorf("delete host batch, but got invalid host id, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsInvalid)})
			return
		}
		iHostIDArr = append(iHostIDArr, iHostID)
	}

	condition := make(map[string]interface{})
	condition["condition"] = NewOperation().WithSupplierID(1).WithOwnerID(ownerID).Data()
	query := meta.QueryInput{Condition: condition}
	result, err := s.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, req.Request.Header, &query)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("delete host in batch, but search instance failed, err: %v, result err: %v", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPReadBodyFailed)})
	}

	js, _ := json.Marshal(result)
	appData := new(AppResult)
	err = json.Unmarshal(js, appData)
	if err != nil {
		blog.Errorf("delete host batch , but unmarshal result failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
	}

	if len(appData.Data.Info) == 0 {
		blog.Error("delete host batch, but can not found it's instance.")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostNotFound)})
	}

	id, exist := appData.Data.Info[0][common.BKAppIDField]
	if !exist {
		blog.Errorf("search host result, but can not find app id.")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CC_Err_Comm_APP_Field_VALID_FAIL)})
		return
	}

	appID, err := util.GetInt64ByInterface(id)
	if err != nil {
		blog.Error("delete host batch, but got invalid app id, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CC_Err_Comm_APP_Field_VALID_FAIL)})
		return
	}

	hostFields, err := s.GetHostAttributes(ownerID, req.Request.Header)
	if err != nil {
		blog.Errorf("delete host batch failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostDeleteFail)})
		return
	}

	var logConents []auditoplog.AuditLogExt
	for _, hostID := range iHostIDArr {
		preHostLog, ip, err := s.GetHostInstanceDetails(ownerID, strconv.FormatInt(hostID, 10), req.Request.Header)
		if err != nil {
			blog.Errorf("get pre host snap failed, err: %v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrHostDeleteFail)})
			return
		}

		opt := meta.ModuleHostConfigParams{
			HostID:        hostID,
			ApplicationID: appID,
		}
		result, err := s.CoreAPI.HostController().Module().DelModuleHostConfig(context.Background(), req.Request.Header, &opt)
		if err != nil || (err == nil && !result.Result) {
			blog.Errorf("delete host in batch, but delete module failed, err: %v, result err: %v", err, result.ErrMsg)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostDeleteFail)})
			return
		}

		delOp := NewOperation().WithInstID(hostID).WithObjID(common.BKInnerObjIDHost).Data()
		delResult, err := s.CoreAPI.ObjectController().Instance().DelObject(context.Background(), common.BKTableNameInstAsst, req.Request.Header, delOp)
		if err != nil || (err == nil && !delResult.Result) {
			blog.Errorf("delete host in batch, but delete object failed, err: %v, result err: %v", err, delResult.ErrMsg)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostDeleteFail)})
			return
		}
		ctnt := meta.Content{
			PreData: preHostLog,
			CurData: nil,
			Headers: hostFields,
		}
		logConents = append(logConents, auditoplog.AuditLogExt{ID: hostID, Content: ctnt, ExtKey: ip})
	}

	hostCond := make(map[string]interface{})
	condInput := make(map[string]interface{})
	hostCond[common.BKDBIN] = iHostIDArr
	condInput[common.BKHostIDField] = hostCond
	delResult, err := s.CoreAPI.ObjectController().Instance().DelObject(context.Background(), common.BKInnerObjIDHost, req.Request.Header, condInput)
	if err != nil || (err == nil && !delResult.Result) {
		blog.Errorf("delete host in batch, but delete host failed, err: %v, result err: %v", err, delResult.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostDeleteFail)})
		return
	}

	auditResult, err := s.CoreAPI.AuditController().AddHostLogs(context.Background(), ownerID, strconv.FormatInt(appID, 10), user, req.Request.Header, logConents)
	if err != nil || (err == nil && !auditResult.Result) {
		blog.Errorf("delete host in batch, but add host audit log failed, err: %v, result err: %v", err, auditResult.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrAuditSaveLogFaile)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) GetHostInstanceProperties(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	hostID := req.PathParameter("bk_host_id")
	ownerID := util.GetOwnerID(pheader)

	details, _, err := s.GetHostInstanceDetails(ownerID, hostID, pheader)
	if err != nil {
		blog.Errorf("get host defails failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostDetailFail)})
		return
	}

	attribute, err := s.GetHostAttributes(ownerID, pheader)
	if err != nil {
		blog.Errorf("get host attribute fields failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostDetailFail)})
		return
	}

	result := make([]meta.HostInstanceProperties, 0)
	for _, attr := range attribute {
		if attr.PropertyID == common.BKChildStr {
			continue
		}
		result = append(result, meta.HostInstanceProperties{
			PropertyID:    attr.PropertyID,
			PropertyName:  attr.PropertyName,
			PropertyValue: details[attr.PropertyID],
		})
	}

	resp.WriteEntity(meta.HostInstancePropertiesResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}

func (s *Service) HostSnapInfo(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	hostID := req.PathParameter(common.BKHostIDField)
	result, err := s.CoreAPI.HostController().Host().GetHostSnap(context.Background(), hostID, pheader)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("get host snap info failed, err: %v, result err: %v", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetSnapshot)})
		return
	}

	snap, err := logics.ParseHostSnap(result.Data.Data)
	if err != nil {
		blog.Errorf("get host snap info, but parse snap info failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetSnapshot)})
		return
	}

	resp.WriteEntity(meta.HostSnapResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     snap,
	})
}

func (s *Service) AddHost(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	hostList := new(meta.HostList)
	if err := json.NewDecoder(req.Request.Body).Decode(hostList); err != nil {
		blog.Errorf("add host failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if hostList.HostInfo == nil {
		blog.Errorf("add host, but host info is nil.")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	appID := hostList.ApplicationID
	if appID == 0 {
		// get default app id
		var err error
		appID, err = s.GetDefaultAppIDWithSupplier(hostList.SupplierID, pheader)
		if err != nil {
			blog.Errorf("add host, but get default appid failed, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CC_Err_Comm_APP_QUERY_FAIL)})
			return
		}
	}

	cond := NewOperation().WithSupplierID(int64(common.DefaultResSetFlag)).WithModuleName(common.DefaultResSetName).WithAppID(appID).Data()
	moduleID, err := s.GetResoulePoolModuleID(pheader, cond)
	if err != nil {
		blog.Errorf("add host, but get module id failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

	succ, updateErrRow, errRow, err := s.Logics.AddHost(appID, moduleID, common.BKDefaultOwnerID, pheader, hostList.HostInfo, hostList.InputType)
	if err != nil {
		blog.Errorf("add host failed, succ: %v, update: %v, err: %v, %v", succ, updateErrRow, err, errRow)

		retData := make(map[string]interface{})
		retData["success"] = succ
		retData["error"] = errRow
		retData["update_error"] = updateErrRow
		resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{false, common.CCErrHostCreateFail, defErr.Error(common.CCErrHostCreateFail).Error()},
			Data:     retData,
		})
	}

	resp.WriteEntity(meta.NewSuccessResp(succ))
}

func (s *Service) AddHostFromAgent(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := common.BKDefaultOwnerID

	agents := make(mapstr.MapStr)
	if err := json.NewDecoder(req.Request.Body).Decode(&agents); err != nil {
		blog.Errorf("add host from agent failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(agents) == 0 {
		blog.Errorf("add host from agent, but got 0 agents from body.")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "HostInfo")})
		return
	}

	appID, err := s.GetDefaultAppID(ownerID, pheader)
	if 0 == appID || nil != err {
		blog.Errorf("add host from agent, but got invalid appid, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrAddHostToModule, err.Error())})
		return
	}

	opt := NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(appID)
	moduleID, err := s.GetResoulePoolModuleID(pheader, opt.Data())
	if err != nil {
		blog.Errorf("add host from agent , but get module id failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

	agents["import_from"] = common.HostAddMethodAgent
	addHost := make(map[int64]map[string]interface{})
	addHost[1] = agents

	succ, updateErrRow, errRow, err := s.Logics.AddHost(appID, moduleID, common.BKDefaultOwnerID, pheader, addHost, "")
	if err != nil {
		blog.Errorf("add host failed, succ: %v, update: %v, err: %v, %v", succ, updateErrRow, err, errRow)

		retData := make(map[string]interface{})
		retData["success"] = succ
		retData["error"] = errRow
		retData["update_error"] = updateErrRow
		resp.WriteEntity(meta.Response{
			BaseResp: meta.BaseResp{false, common.CCErrHostCreateFail, defErr.Error(common.CCErrHostCreateFail).Error()},
			Data:     retData,
		})
	}

	resp.WriteEntity(meta.NewSuccessResp(succ))

}

func (s *Service) AddHistory(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	data := make(mapstr.MapStr)
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("add host from agent failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	content, ok := data["content"].(string)
	if !ok || "" == content {
		blog.Error("add history, but content is empty. data: %v", data)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return

	}
	params := make(map[string]interface{}, 1)
	params["content"] = content

	result, err := s.CoreAPI.HostController().History().AddHistory(context.Background(), user, pheader, &meta.HistoryContent{Content: content})
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("add history failed, err: %v, result err: %v", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostHisCreateFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(result.Data))
}

func (s *Service) GetHistorys(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	start := req.PathParameter("start")
	limit := req.PathParameter("limit")

	result, err := s.CoreAPI.HostController().History().GetHistorys(context.Background(), user, start, limit, pheader)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("get history failed, err: %v, result err: %v", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostHisGetFail)})
		return
	}

	resp.WriteEntity(meta.GetHistoryResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     result.Data,
	})

}

func (s *Service) GetAgentStatus(req *restful.Request, resp *restful.Response) {

}

func (s *Service) UpdateHost(req *restful.Request, resp *restful.Response) {

}

func (s *Service) UpdateHostByAppID(req *restful.Request, resp *restful.Response) {

}

func (s *Service) HostSearchByIP(req *restful.Request, resp *restful.Response) {

}

func (s *Service) HostSearchByConds(req *restful.Request, resp *restful.Response) {

}

func (s *Service) HostSearchByModuleID(req *restful.Request, resp *restful.Response) {

}

func (s *Service) HostSearchBySetID(req *restful.Request, resp *restful.Response) {

}

func (s *Service) HostSearchByAppID(req *restful.Request, resp *restful.Response) {

}

func (s *Service) HostSearchByProperty(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetIPAndProxyByCompany(req *restful.Request, resp *restful.Response) {

}

func (s *Service) UpdateCustomProperty(req *restful.Request, resp *restful.Response) {

}

func (s *Service) CloneHostProperty(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetHostAppByCompanyId(req *restful.Request, resp *restful.Response) {

}

func (s *Service) DelHostInApp(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetGitServerIp(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetPlat(req *restful.Request, resp *restful.Response) {

}

func (s *Service) CreatePlat(req *restful.Request, resp *restful.Response) {

}

func (s *Service) DelPlat(req *restful.Request, resp *restful.Response) {

}

func (s *Service) SearchHost(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	body := new(meta.HostCommonSearch)
	if err := json.NewDecoder(req.Request.Body).Decode(body); err != nil {
		blog.Errorf("search host failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

}

func (s *Service) SearchHostWithAsstDetail(req *restful.Request, resp *restful.Response) {

}

func (s *Service) UpdateHostBatch(req *restful.Request, resp *restful.Response) {

}

func (s *Service) Add(req *restful.Request, resp *restful.Response) {

}

func (s *Service) Update(req *restful.Request, resp *restful.Response) {

}

func (s *Service) Delete(req *restful.Request, resp *restful.Response) {

}

func (s *Service) Get(req *restful.Request, resp *restful.Response) {

}

func (s *Service) Detail(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetUserAPIData(req *restful.Request, resp *restful.Response) {

}
