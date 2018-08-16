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
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// GetRolePri get role privilege
func (cli *Service) GetRolePri(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	defer req.Request.Body.Close()
	pathParams := req.PathParameters()
	objID := pathParams["bk_obj_id"]
	propertyID := pathParams["bk_property_id"]
	cond := make(map[string]interface{})
	cond[common.BKObjIDField] = objID
	cond[common.BKPropertyIDField] = propertyID
	var result map[string]interface{}
	cond = util.SetModOwner(cond, ownerID)

	cnt, err := db.Table(common.BKTableNamePrivilege).Find(cond).Count(ctx)
	if nil != err {
		blog.Error("get user group privi error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	if 0 == cnt { // TODO:
		blog.V(3).Infof("failed to find the cnt")
		info := make(map[string]interface{})
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: info})
		return
	}

	err = db.Table(common.BKTableNamePrivilege).Find(cond).All(ctx, &result)
	if nil != err {
		blog.Error("get role pri field error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
		return

	}
	privilege, ok := result["privilege"]
	if !ok {
		blog.Errorf("not privilege, the origin data is %#v", result)
		info := make(map[string]interface{})
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: info})
		return

	}
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: privilege})
}

//CreateRolePri create role privilege
func (cli *Service) CreateRolePri(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	pathParams := req.PathParameters()
	objID := pathParams["bk_obj_id"]
	propertyID := pathParams["bk_property_id"]
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read json data error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}
	var roleJSON []string
	err = json.Unmarshal([]byte(value), &roleJSON)
	if err != nil {
		blog.Error("read json data error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	input := make(map[string]interface{})
	input[common.BKOwnerIDField] = ownerID
	input[common.BKObjIDField] = objID
	input[common.BKPropertyIDField] = propertyID
	input[common.BKPrivilegeField] = roleJSON
	input = util.SetModOwner(input, ownerID)

	err = db.Table(common.BKTableNamePrivilege).Insert(ctx, input)
	if nil != err {
		blog.Error("create role privilege error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

//UpdateRolePri update role privilege
func (cli *Service) UpdateRolePri(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	pathParams := req.PathParameters()
	objID := pathParams["bk_obj_id"]
	propertyID := pathParams["bk_property_id"]
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read json data error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}
	var roleJSON []string
	err = json.Unmarshal([]byte(value), &roleJSON)
	if err != nil {
		blog.Error("read json data error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	input := make(map[string]interface{})
	cond := make(map[string]interface{})
	cond[common.BKOwnerIDField] = ownerID
	cond[common.BKObjIDField] = objID
	cond[common.BKPropertyIDField] = propertyID
	input[common.BKPrivilegeField] = roleJSON
	cond = util.SetModOwner(cond, ownerID)

	err = db.Table(common.BKTableNamePrivilege).Update(ctx, cond, input)
	if nil != err {
		blog.Error("update role privilege error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}
