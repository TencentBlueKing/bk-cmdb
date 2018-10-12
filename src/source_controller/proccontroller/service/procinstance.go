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
	"net/http"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (ps *ProctrlServer) CreateProcInstanceModel(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	reqParam := make([]meta.ProcInstanceModel, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("create process instance model failed, decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.Infof("will create process instance model: %+v", reqParam)

	blog.V(3).Infof("will create process instance model: %+v", reqParam)
	if 0 == len(reqParam) {
		resp.WriteEntity(meta.NewSuccessResp(nil))
		return
	}

	addInst := make([]interface{}, 0)
	for _, item := range reqParam {
		item.OwnerID = util.GetOwnerID(req.Request.Header)
		addInst = append(addInst, item)
	}

	if err := ps.Instance.Table(common.BKTableNameProcInstanceModel).Insert(ctx, addInst); err != nil {
		blog.Errorf("create process instance model failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateInstanceModel)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) GetProcInstanceModel(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	reqParam := new(meta.QueryInput)
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("get process instance model failed, decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.Infof("will get process instance model. condition: %v", reqParam)

	cnt, err := ps.Instance.Table(common.BKTableNameProcInstanceModel).Find(reqParam.Condition).Count(ctx)
	if err != nil {
		blog.Errorf("get process instance model failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetInstanceModel)})
		return
	}
	blog.V(3).Infof("will get process instance model. condition: %v", reqParam)
	data := make([]meta.ProcInstanceModel, 0)
	err = ps.Instance.Table(common.BKTableNameProcInstanceModel).Find(reqParam.Condition).Fields(strings.Split(reqParam.Fields, ",")...).
		Sort(reqParam.Sort).Start(uint64(reqParam.Start)).Limit(uint64(reqParam.Limit)).All(ctx, &data)
	if err != nil {
		blog.Errorf("get process instance model failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetInstanceModel)})
		return
	}
	ret := meta.ProcInstModelResult{
		BaseResp: meta.SuccessBaseResp,
	}
	ret.Data.Info = data
	ret.Data.Count = int(cnt)
	resp.WriteEntity(ret)
}

func (ps *ProctrlServer) DeleteProcInstanceModel(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	reqParam := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("delete process instance model failed, decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	reqParam = util.SetModOwner(reqParam, util.GetOwnerID(req.Request.Header))

	blog.Infof("will delete process instance model. condition: %+v", reqParam)

	blog.V(3).Infof("will delete process instance model. condition: %+v", reqParam)
	if err := ps.Instance.Table(common.BKTableNameProcInstanceModel).Delete(ctx, reqParam); err != nil {
		blog.Errorf("delete process instance model failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteInstanceModel)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) RegisterProcInstaceDetail(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := new(meta.GseProcRequest)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("register  process instance detail failed, decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	for _, gseHost := range input.Hosts {
		conds := common.KvMap{common.BKAppIDField: input.AppID, common.BKProcessIDField: input.ProcID, common.BKModuleIDField: input.ModuleID, common.BKHostIDField: gseHost.HostID}
		cnt, err := ps.Instance.Table(common.BKTableNameProcInstaceDetail).Find(conds).Count(ctx)
		if nil != err {
			blog.Errorf("register  process instance detail get info error: %s", err.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}
		detail := new(meta.ProcInstanceDetail)
		detail.OwnerID = util.GetOwnerID(req.Request.Header)
		detail.AppID = input.AppID
		detail.Meta = input.Meta
		detail.ProcID = input.ProcID
		detail.ModuleID = input.ModuleID
		detail.HostID = gseHost.HostID
		detail.Spec = input.Spec
		detail.Hosts = append(detail.Hosts, gseHost)
		detail.Status = meta.ProcInstanceDetailStatusRegisterSucc //1 register gse sucess, 2 register error need retry 3 unregister error need retry
		if 0 == cnt {
			err = ps.Instance.Table(common.BKTableNameProcInstaceDetail).Insert(ctx, detail)
		} else {
			err = ps.Instance.Table(common.BKTableNameProcInstaceDetail).Update(ctx, conds, detail)
		}
		if nil != err {
			blog.Errorf("register  process instance detail save info error: %s", err.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBUpdateFailed)})
			return
		}
	}
	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) ModifyRegisterProcInstanceDetail(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := new(meta.ModifyProcInstanceDetail)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("modify register  process instance detail failed, decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	input.Conds = util.SetModOwner(input.Conds, util.GetOwnerID(req.Request.Header))

	err := ps.Instance.Table(common.BKTableNameProcInstaceDetail).Update(ctx, input.Conds, input.Data)
	if nil != err {
		blog.Errorf("update register  process instance detail  info error: %s, input:%s", err.Error(), input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}
	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) GetProcInstanceDetail(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := new(meta.QueryInput)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("get process instance detail failed, decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	input.Condition = util.SetModOwner(input.Condition, util.GetOwnerID(req.Request.Header))
	cnt, err := ps.Instance.Table(common.BKTableNameProcInstaceDetail).Find(input.Condition).Count(ctx)
	if err != nil {
		blog.Errorf("get process instance detail failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	blog.V(3).Infof("will get process instance detail. condition: %v", input)
	data := make([]meta.ProcInstanceDetail, 0)
	err = ps.Instance.Table(common.BKTableNameProcInstaceDetail).Find(input.Condition).Fields(strings.Split(input.Fields, ",")...).
		Sort(input.Sort).Start(uint64(input.Start)).Limit(uint64(input.Limit)).All(ctx, &data)
	if err != nil {
		blog.Errorf("get process instance detail failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	ret := meta.ProcInstanceDetailResult{
		BaseResp: meta.SuccessBaseResp,
	}
	ret.Data.Info = data
	ret.Data.Count = int(cnt)
	resp.WriteEntity(ret)
}

func (ps *ProctrlServer) DeleteRegisterProcInstanceDetail(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := make(map[string]interface{}, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("modify register  process instance detail failed, decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	input = util.SetModOwner(input, util.GetOwnerID(req.Request.Header))
	err := ps.Instance.Table(common.BKTableNameProcInstaceDetail).Delete(ctx, input)
	if nil != err {
		blog.Errorf("update register  process instance detail  info error: %s, input:%s", err.Error(), input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBDeleteFailed)})
		return
	}
	resp.WriteEntity(meta.NewSuccessResp(nil))
}
