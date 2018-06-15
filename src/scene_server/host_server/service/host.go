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
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"
	"configcenter/src/source_controller/api/metadata"
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
	condition["condition"] = NewOperation().WithDefaultField(1).WithOwnerID(ownerID).Data()
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
		ctnt := metadata.Content{
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
	}

	attribute, err := s.GetHostAttributes(ownerID, pheader)
	if err != nil {
		blog.Errorf("get host attribute fields failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostDetailFail)})
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
		appID, err := s.GetDefaultAppIDWithSupplier(strconv.FormatInt(hostList.SupplierID, 10), pheader)
		if err != nil {
			blog.Errorf("add host, but get default appid failed, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CC_Err_Comm_APP_QUERY_FAIL)})
			return
		}
	}

	cond := NewOperation().WithDefaultField(int64(common.DefaultResSetFlag)).WithModuleName(common.DefaultResSetName).WithAppID(appID).Data()
	moduleID, err := s.GetResoulePoolModuleID(pheader, cond)
	if err != nil {
		blog.Errorf("add host, but get module id failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

}

func (s *Service) AddHostFromAgent(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetHostFavourites(req *restful.Request, resp *restful.Response) {

}

func (s *Service) AddHostFavourite(req *restful.Request, resp *restful.Response) {

}

func (s *Service) UpdateHostFavouriteByID(req *restful.Request, resp *restful.Response) {

}

func (s *Service) DeleteHostFavouriteByID(req *restful.Request, resp *restful.Response) {

}

func (s *Service) IncrHostFavouritesCount(req *restful.Request, resp *restful.Response) {

}

func (s *Service) AddHistory(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetHistorys(req *restful.Request, resp *restful.Response) {

}

func (s *Service) AddHostMutiltAppModuleRelation(req *restful.Request, resp *restful.Response) {

}

func (s *Service) HostModuleRelation(req *restful.Request, resp *restful.Response) {

}

func (s *Service) MoveHost2EmptyModule(req *restful.Request, resp *restful.Response) {

}

func (s *Service) MoveHost2FaultModule(req *restful.Request, resp *restful.Response) {

}

func (s *Service) MoveHostToResourcePool(req *restful.Request, resp *restful.Response) {

}

func (s *Service) AssignHostToApp(req *restful.Request, resp *restful.Response) {

}

func (s *Service) AssignHostToAppModule(req *restful.Request, resp *restful.Response) {

}

func (s *Service) SaveUserCustom(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetUserCustom(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetDefaultCustom(req *restful.Request, resp *restful.Response) {

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

func (s *Service) HostSearch(req *restful.Request, resp *restful.Response) {

}

func (s *Service) HostSearchWithAsstDetail(req *restful.Request, resp *restful.Response) {

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
