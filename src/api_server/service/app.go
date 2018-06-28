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

func (s *Service) GetAppList(req *restful.Request, resp *restful.Response) {
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

func (s *Service) GetAppList(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetAppList start !")

	//set empty to get all fields
	param := map[string]interface{}{
		"condition": nil,
	}

	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))
	ownerID, user := util.GetOwnerIDAndUser(req.Request.Header)

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	paramJson, _ := json.Marshal(param)
	//url := fmt.Sprintf("%s/topo/v1/app/search/"+common.BKDefaultOwnerID, cli.CC.TopoAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))

	result, err := s.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, req.Request.Header, &query)
	if err != nil {
		blog.Error("GetAppList url:%s, params:%s error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForAppList(rspV3)
	if err != nil {
		blog.Error("convert app res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) GetAppByID(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetAppByUIN(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetUserRoleApp(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetAppSetModuleTreeByAppId(req *restful.Request, resp *restful.Response) {

}

func (s *Service) AddApp(req *restful.Request, resp *restful.Response) {

}

func (s *Service) DeleteApp(req *restful.Request, resp *restful.Response) {

}

func (s *Service) EditApp(req *restful.Request, resp *restful.Response) {

}

func (s *Service) GetHostAppByCompanyId(req *restful.Request, resp *restful.Response) {

}
