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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"
)

func (ps *ProctrlServer) CreateConfigTemp(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetLanguage(req.Request.Header)
	// get the error factory by the language
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("create config template failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.Infof("create process config template: %v", input)
	ec := eventclient.NewEventContextByReq(req.Request.Header, ps.CacheDI)
	if _, err := ps.DbInstance.Insert(common.BKTableNameProcConf, input); err != nil {
		blog.Errorf("create config template failed when insert confit template into db. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateProcConf)})
		return
	}

	// recode events
	if err := ec.InsertEvent(meta.EventTypeRelation, "procconf", meta.EventActionCreate, input, nil); err != nil {
		blog.Warnf("insert config template create event failed. err: %v", err)
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) DeleteConfigTemp(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetLanguage(req.Request.Header)
	// get the error factory by the language
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("delete config template failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	// get original data
	var oriData interface{}
	if err := ps.DbInstance.GetOneByCondition(common.BKTableNameProcConf, []string{}, input, &oriData); err != nil {
		blog.Warnf("get original config template data failed. err: %v", err)
	}

	// delete process config template
	blog.Infof("will delete config template, param: %v", input)
	if err := ps.DbInstance.DelByCondition(common.BKTableNameProcConf, input); err != nil {
		blog.Errorf("delete config template failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteProcConf)})
		return
	}

	// record events
	ec := eventclient.NewEventContextByReq(req.Request.Header, ps.CacheDI)
	if err := ec.InsertEvent(meta.EventTypeRelation, "procconf", meta.EventActionDelete, oriData, nil); err != nil {
		blog.Warnf("intert config template delete event failed. err: %v", err)
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) UpdateConfigTemp(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetLanguage(req.Request.Header)
	// get the error factory by the language
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("update config template failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	confTempID := input[common.BKConfTempIdField]
	condition := make(map[string]interface{})
	condition[common.BKConfTempIdField] = confTempID

	// get original data before update in order to save event
	var oriData interface{}
	if err := ps.DbInstance.GetOneByCondition(common.BKTableNameProcConf, []string{}, condition, &oriData); err != nil {
		blog.Warnf("get original config template data failed. err: %v", err)
	}

	if err := ps.DbInstance.UpdateByCondition(common.BKTableNameProcConf, input, condition); err != nil {
		blog.Errorf("update config template failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcUpdateProcConf)})
		return
	}

	// record events
	ec := eventclient.NewEventContextByReq(req.Request.Header, ps.CacheDI)
	if err := ec.InsertEvent(meta.EventTypeRelation, "procconf", meta.EventActionDelete, oriData, input); err != nil {
		blog.Warnf("intert config template delete event failed. err: %v", err)
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) QueryConfigTemp(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetLanguage(req.Request.Header)
	// get the error factory by the language
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("query config template failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.Infof("will query config template. param: %v", input)
	var result interface{}
	if err := ps.DbInstance.GetMutilByCondition(common.BKTableNameProcConf, []string{}, input, &result, "", 0, 0); err != nil {
		blog.Errorf("query config templdate failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetProcConf)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}
