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

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"
	hutil "configcenter/src/scene_server/host_server/util"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/metadata"
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
	condition = hutil.NewOperation().WithDefaultField(int64(common.DefaultAppFlag)).WithOwnerID(ownerID).Data()
	query := meta.QueryInput{Condition: condition}
	query.Limit = 1
	result, err := s.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, req.Request.Header, &query)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("delete host in batch, but search instance failed, err: %v, result err: %v", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPReadBodyFailed)})
		return
	}
	if err != nil || (err != nil && result.Result == false) {
		blog.Errorf("delete host batch , but unmarshal result failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(result.Data.Info) == 0 {
		blog.Error("delete host batch, but can not found it's instance.")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostNotFound)})
		return
	}

	appID, err := result.Data.Info[0].Int64(common.BKAppIDField)
	if err != nil {
		blog.Error("delete host batch, but got invalid app id, err: %v, appinfo:%v", err, result.Data.Info[0])
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
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
		logger := s.Logics.NewHostLog(req.Request.Header, ownerID)
		if err := logger.WithPrevious(strconv.FormatInt(hostID, 10), hostFields); err != nil {
			blog.Errorf("delet host batch, but get pre host data failed, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
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

		delOp := hutil.NewOperation().WithInstID(hostID).WithObjID(common.BKInnerObjIDHost).Data()
		delResult, err := s.CoreAPI.ObjectController().Instance().DelObject(context.Background(), common.BKTableNameInstAsst, req.Request.Header, delOp)
		if err != nil || (err == nil && !delResult.Result) {
			blog.Errorf("delete host in batch, but delete object failed, err: %v, result err: %v", err, delResult.ErrMsg)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostDeleteFail)})
			return
		}

		if err := logger.WithCurrent(strconv.FormatInt(hostID, 10)); err != nil {
			blog.Errorf("delet host batch, but get current host data failed, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
			return
		}

		logConents = append(logConents, *logger.AuditLog(hostID))
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

	addHostLogs := common.KvMap{common.BKContentField: logConents, common.BKOpDescField: "delete host", common.BKOpTypeField: auditoplog.AuditOpTypeDel}
	auditResult, err := s.CoreAPI.AuditController().AddHostLogs(context.Background(), ownerID, strconv.FormatInt(appID, 10), user, req.Request.Header, addHostLogs)
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

	details, _, err := s.GetHostInstanceDetails(pheader, ownerID, hostID)
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
		appID, err = s.GetDefaultAppIDWithSupplier(pheader)
		if err != nil {
			blog.Errorf("add host, but get default appid failed, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CC_Err_Comm_APP_QUERY_FAIL)})
			return
		}
	}

	cond := hutil.NewOperation().WithModuleName(common.DefaultResModuleName).WithAppID(appID).Data()
	cond[common.BKDefaultField] = common.DefaultResModuleFlag
	moduleID, err := s.GetResoulePoolModuleID(pheader, cond)
	if err != nil {
		blog.Errorf("add host, but get module id failed, err: %s", err.Error())
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
		return
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

	opt := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(appID)
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
		return
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

func (s *Service) SearchHost(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	body := new(meta.HostCommonSearch)
	if err := json.NewDecoder(req.Request.Body).Decode(body); err != nil {
		blog.Errorf("search host failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	host, err := s.Logics.SearchHost(pheader, body, false)
	if err != nil {
		blog.Errorf("search host failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	resp.WriteEntity(meta.SearchHostResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *host,
	})
}

func (s *Service) SearchHostWithAsstDetail(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	body := new(meta.HostCommonSearch)
	if err := json.NewDecoder(req.Request.Body).Decode(body); err != nil {
		blog.Errorf("search host failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	host, err := s.Logics.SearchHost(pheader, body, true)
	if err != nil {
		blog.Errorf("search host failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	resp.WriteEntity(meta.SearchHostResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *host,
	})
}

func (s *Service) UpdateHostBatch(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)
	data := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update host batch failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	hostIDStr, ok := data[common.BKHostIDField].(string)
	if !ok {
		blog.Errorf("update host batch failed, but without host id field.")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFeildValidFail)})
		return

	}
	delete(data, common.BKHostIDField)
	ownerID := util.GetOwnerID(req.Request.Header)
	valid := validator.NewValidMap(ownerID, common.BKInnerObjIDHost, req.Request.Header, s.Engine)

	hostFields, err := s.Logics.GetHostAttributes(ownerID, pheader)
	if err != nil {
		blog.Errorf("update host batch, but get host attribute for audit failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrHostDetailFail)})
		return
	}
	audit := s.Logics.NewHostLog(pheader, common.BKDefaultOwnerID)

	logPreConents := make(map[int64]auditoplog.AuditLogExt, 0)
	hostIDs := make([]int64, 0)
	for _, id := range strings.Split(hostIDStr, ",") {
		hostID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			blog.Errorf("update host batch, but got invalid host id[%s], err: %v", id, err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsInvalid)})
			return
		}
		hostIDs = append(hostIDs, hostID)

		err = valid.ValidMap(data, common.ValidUpdate, hostID)
		if nil != err {
			blog.Errorf("update host batch, but invalid host failed, id[%s], err: %v", id, err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFeildValidFail)})
			return

		}

		if err := audit.WithPrevious(id, hostFields); err != nil {
			blog.Errorf("update host batch, but get host[%s] pre data for audit failed, err: %v", id, err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrHostDetailFail)})
			return
		}
		logPreConents[hostID] = *audit.AuditLog(hostID)
	}

	opt := common.KvMap{"condition": common.KvMap{common.BKHostIDField: common.KvMap{common.BKDBIN: hostIDs}}, "data": data}
	result, err := s.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDHost, pheader, opt)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("update host batch failed, ids[%v], err: %v, %v", hostIDs, err, result.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrHostUpdateFail)})
		return
	}

	hostModuleConfig, err := s.Logics.GetConfigByCond(pheader, map[string][]int64{common.BKHostIDField: hostIDs})
	if err != nil {
		blog.Errorf("update host batch failed, ids[%v], err: %v", hostIDs, err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrHostUpdateFail)})
		return
	}

	appID := "0"
	if len(hostModuleConfig) != 0 {
		appID = strconv.FormatInt(hostModuleConfig[0][common.BKAppIDField], 10)
	}

	logLastConents := make([]auditoplog.AuditLogExt, 0)
	for _, hostID := range hostIDs {
		inst := logics.NewImportInstance(common.BKDefaultOwnerID, pheader, s.Engine)
		err := inst.UpdateInstAssociation(pheader, hostID, common.BKDefaultOwnerID, common.BKInnerObjIDHost, data)
		if err != nil {
			blog.Errorf("update host batch, but update inst association failed, id[%v], err: %v", hostID, err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrHostUpdateFail)})
			return
		}

		if err := audit.WithPrevious(strconv.FormatInt(hostID, 10), hostFields); err != nil {
			blog.Errorf("update host batch, but get host[%s] pre data for audit failed, err: %v", hostID, err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrHostDetailFail)})
			return
		}
		logContent := audit.Content
		logContent.CurData = logContent.PreData
		preLogContent, ok := logPreConents[hostID]
		if ok {
			content, ok := preLogContent.Content.(*metadata.Content)
			if ok {
				logContent.PreData = content.PreData
			}
		}

		logLastConents = append(logLastConents, auditoplog.AuditLogExt{ID: hostID, Content: logContent, ExtKey: preLogContent.ExtKey})
	}
	log := common.KvMap{common.BKContentField: logLastConents, common.BKOpDescField: "update host", common.BKOpTypeField: auditoplog.AuditOpTypeModify}
	aResult, err := s.CoreAPI.AuditController().AddHostLogs(context.Background(), common.BKDefaultOwnerID, appID, user, pheader, log)
	if err != nil || (err == nil && !aResult.Result) {
		blog.Errorf("update host batch, but add host[%s] audit failed, err: %v, %v", hostIDs, err, aResult.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrHostDetailFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}
