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

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// DeleteProc2Module delete proc module config
func (ps *ProctrlServer) DeleteProc2Module(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("delete process2module failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	// retrieve original data
	var originals []interface{}
	if err := ps.DbInstance.GetMutilByCondition(common.BKTableNameProcModule, []string{}, input, &originals, "", 0, 0); err != nil {
		blog.Warnf("retrieve original error:%v", err)
	}

	// delete proc module config
	blog.Infof("delete proc module config %v", input)
	if err := ps.DbInstance.DelByCondition(common.BKTableNameProcModule, input); err != nil {
		blog.Errorf("delete proc module config error: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteProc2Module)})
		return
	}

	//send  event
	if len(originals) > 0 {
		ec := eventclient.NewEventContextByReq(req.Request.Header, ps.CacheDI)
		for _, i := range originals {
			if err := ec.InsertEvent(meta.EventTypeRelation, "processmodule", meta.EventActionDelete, nil, i); err != nil {
				blog.Warnf("create event error:%s", err.Error())
			}
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) CreateProc2Module(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetLanguage(req.Request.Header)
	// get the error factory by the language
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	input := make([]interface{}, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("create process2module failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.Infof("create proc module config: %v ", input)
	ec := eventclient.NewEventContextByReq(req.Request.Header, ps.CacheDI)
	for _, i := range input {
		if _, err := ps.DbInstance.Insert(common.BKTableNameProcModule, i); err != nil {
			blog.Errorf("create proc module config error:%v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateProc2Module)})
			return
		}
		//  record events
		if err := ec.InsertEvent(meta.EventTypeRelation, "processmodule", meta.EventActionCreate, i, nil); err != nil {
			blog.Errorf("create event error: %v", err)
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) GetProc2Module(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetLanguage(req.Request.Header)
	// get the error factory by the language
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("get process2module failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.Infof("get proc module config condition: %v ", input)
	result := make([]interface{}, 0)
	if err := ps.DbInstance.GetMutilByCondition(common.BKTableNameProcModule, nil, input, &result, "", 0, 0); err != nil {
		blog.Errorf("get process2module config failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcSelectProc2Module)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}
