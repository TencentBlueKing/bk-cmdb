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
	"io/ioutil"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
	"github.com/rs/xid"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

//CreateUserGroup create group
func (cli *Service) CreateUserGroup(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read json data error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Error("create user group failed, err msg : %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	data, err := js.Map()
	if err != nil {
		blog.Error("create user group failed, err msg : %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	guid := xid.New()
	data[common.BKUserGroupIDField] = guid.String()
	data = util.SetModOwner(data, ownerID)
	err = cli.Instance.Table(common.BKTableNameUserGroup).Insert(context.Background(), data)
	if nil != err {
		blog.Error("create user group error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

//UpdateUserGroup create group
func (cli *Service) UpdateUserGroup(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	pathParams := req.PathParameters()
	groupID := pathParams["group_id"]
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read json data error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Error("update user group failed, err msg : %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	data, err := js.Map()
	if err != nil {
		blog.Error("update user group failed, err msg : %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	cond := make(map[string]interface{})
	cond[common.BKUserGroupIDField] = groupID
	cond = util.SetModOwner(cond, ownerID)
	err = cli.Instance.Table(common.BKTableNameUserGroup).Update(context.Background(), cond, data)
	if nil != err {
		blog.Error("update user group error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

//DeleteUserGroup create group
func (cli *Service) DeleteUserGroup(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	defer req.Request.Body.Close()
	pathParams := req.PathParameters()
	groupID := pathParams["group_id"]
	cond := make(map[string]interface{})
	cond[common.BKUserGroupIDField] = groupID
	cond = util.SetModOwner(cond, ownerID)
	err := cli.Instance.Table(common.BKTableNameUserGroup).Delete(context.Background(), cond)
	if nil != err {
		blog.Error("delete user group error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

//SearchUserGroup create group
func (cli *Service) SearchUserGroup(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read json data error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Error("get user group failed, err msg : %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	cond, err := js.Map()
	if err != nil {
		blog.Error("get user group failed, err msg : %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	cond = util.SetModOwner(cond, ownerID)
	var result []interface{}
	err = cli.Instance.Table(common.BKTableNameUserGroup).Find(cond).All(context.Background(), &result)
	if nil != err {
		blog.Error("get user group error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: result})
}
