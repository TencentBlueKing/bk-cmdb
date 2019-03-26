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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

// CreateUserGroupPrivi create group privi
func (cli *Service) CreateUserGroupPrivi(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	pathParams := req.PathParameters()
	groupID := pathParams["group_id"]
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read json data error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}
	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Errorf("insert user group privi failed, err msg : %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	info, err := js.Map()
	if err != nil {
		blog.Errorf("insert user group privi failed, err msg : %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	data := make(map[string]interface{})
	data[common.BKUserGroupIDField] = groupID
	data[common.BKPrivilegeField] = info
	data = util.SetModOwner(data, ownerID)

	cond := make(map[string]interface{})
	cond[common.BKOwnerIDField] = ownerID
	cond[common.BKUserGroupIDField] = groupID
	cond = util.SetModOwner(cond, ownerID)
	cnt, err := db.Table(common.BKTableNameUserGroupPrivilege).Find(cond).Count(ctx)
	if nil != err && !db.IsNotFoundError(err) {
		blog.Errorf("get user group privi error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	if cnt > 0 {
		blog.V(3).Infof("update user group privi: %+v, by condition %+v ", data, cond)
		err = db.Table(common.BKTableNameUserGroupPrivilege).Update(ctx, cond, data)
		if nil != err {
			blog.Errorf("update user group privi error :%v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
			return
		}
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
		return
	}

	blog.V(3).Infof("create user group privi: %+v", data)
	err = db.Table(common.BKTableNameUserGroupPrivilege).Insert(ctx, data)
	if nil != err {
		blog.Errorf("insert user group privi error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

// UpdateUserGroupPrivi update group privi
func (cli *Service) UpdateUserGroupPrivi(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	pathParams := req.PathParameters()
	groupID := pathParams["group_id"]
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read json data error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}
	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Errorf("update user group privi failed, err msg : %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	info, err := js.Map()
	if err != nil {
		blog.Errorf("update user group privi failed, err msg : %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	cond := make(map[string]interface{})
	data := make(map[string]interface{})
	cond[common.BKUserGroupIDField] = groupID
	data[common.BKPrivilegeField] = info
	cond = util.SetModOwner(cond, ownerID)
	blog.V(3).Infof("update user group privi: %+v, by condition %+v ", data, cond)
	err = db.Table(common.BKTableNameUserGroupPrivilege).Update(ctx, cond, data)
	if nil != err {
		blog.Errorf("update user group privi error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

//GetUserGroupPrivi get group privi
func (cli *Service) GetUserGroupPrivi(req *restful.Request, resp *restful.Response) {

	//get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	pathParams := req.PathParameters()
	groupID := pathParams["group_id"]

	cond := make(map[string]interface{})
	cond[common.BKOwnerIDField] = ownerID
	cond[common.BKUserGroupIDField] = groupID
	cond = util.SetModOwner(cond, ownerID)

	blog.V(3).Infof("get user group privi by condition %+v", cond)
	cnt, err := db.Table(common.BKTableNameUserGroupPrivilege).Find(cond).Count(ctx)
	if nil != err && !db.IsNotFoundError(err) {
		blog.Errorf("get user group privi error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	if 0 == cnt { // TODO: 兼容老的逻辑
		data := make(map[string]interface{})
		data[common.BKOwnerIDField] = ownerID
		data[common.BKUserGroupIDField] = groupID
		data[common.BKPrivilegeField] = common.KvMap{}
		blog.V(3).Infof("get user group privi by condition %+v, returns %+v", cond, data)
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: data})
		return
	}

	var result interface{}
	err = db.Table(common.BKTableNameUserGroupPrivilege).Find(cond).One(ctx, &result)
	if nil != err {
		blog.Errorf("get user group privi error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	blog.V(3).Infof("get user group privi by condition %+v, returns %+v", cond, result)
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: result})
}
