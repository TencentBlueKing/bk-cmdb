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
	"time"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (ps *ProctrlServer) AddOperateTaskInfo(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := make([]meta.ProcessOperateTask, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("add operate process task failed! decode request body err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if 0 == len(input) {
		resp.WriteEntity(meta.NewSuccessResp(nil))
		return
	}
	insts := make([]interface{}, 0)
	ts := time.Now().UTC()
	for _, item := range input {
		item.OwnerID = util.GetOwnerID(req.Request.Header)
		item.User = util.GetUser(req.Request.Header)
		item.CreateTime = ts
		insts = append(insts, item)
	}
	err := ps.Instance.Table(common.BKTableNameProcOperateTask).Insert(ctx, insts)
	if nil != err {
		blog.Errorf("add  operate process task  to db failed   error:%s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) UpdateOperateTaskInfo(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := new(meta.UpdateParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("update operate process task failed! decode request body err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input.Condition = util.SetModOwner(input.Condition, util.GetOwnerID(req.Request.Header))
	if 0 == len(input.Data) {
		resp.WriteEntity(meta.NewSuccessResp(nil))
		return
	}

	err := ps.Instance.Table(common.BKTableNameProcOperateTask).Update(ctx, input.Condition, input.Data)
	if nil != err {
		blog.Errorf("update  operate process task  to db failed   error:%s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) SearchOperateTaskInfo(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)

	input := new(meta.QueryInput)
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("search operate process task failed! decode request body err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input.Condition = util.SetModOwner(input.Condition, util.GetOwnerID(req.Request.Header))
	cnt, err := ps.Instance.Table(common.BKTableNameProcOperateTask).Find(input.Condition).Count(ctx)
	if err != nil {
		blog.Errorf("search operate process taskfailed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	blog.V(3).Infof("will search operate process task. condition: %v", input)
	data := make([]meta.ProcessOperateTask, 0)
	err = ps.Instance.Table(common.BKTableNameProcOperateTask).Find(input.Condition).Fields(strings.Split(input.Fields, ",")...).
		Sort(input.Sort).Start(uint64(input.Start)).Limit(uint64(input.Limit)).All(ctx, &data)
	if err != nil {
		blog.Errorf("search operate process task failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	ret := meta.ProcessOperateTaskResult{
		BaseResp: meta.SuccessBaseResp,
	}
	ret.Data.Info = data
	ret.Data.Count = int(cnt)
	resp.WriteEntity(ret)

}
