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
	"encoding/json"
	"io/ioutil"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

// GetRolePri get role privilege
func (cli *Service) GetRolePri(req *restful.Request, resp *restful.Response) {
	// get the language
	//language := util.GetActionLanguage(req)
	// get the error factory by the language
	//defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	defer req.Request.Body.Close()
	pathParams := req.PathParameters()
	ownerID := pathParams["bk_supplier_account"]
	objID := pathParams["bk_obj_id"]
	propertyID := pathParams["bk_property_id"]
	cond := make(map[string]interface{})
	cond[common.BKOwnerIDField] = ownerID
	cond[common.BKObjIDField] = objID
	cond[common.BKPropertyIDField] = propertyID
	var result map[string]interface{}
	err := cli.CC.InstCli.GetOneByCondition(common.BKTableNamePrivilege, []string{}, cond, &result)
	if nil != err {
		blog.Error("get role pri field error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
		return

	}
	privilege, ok := result["privilege"]
	if !ok {
		blog.Error("get role pri field error :%v", err)
		info := make(map[string]interface{})
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: info})

	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: privilege})
}

//CreateRolePri create role privilege
func (cli *Service) CreateRolePri(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	pathParams := req.PathParameters()
	ownerID := pathParams["bk_supplier_account"]
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
	_, err = cli.CC.InstCli.Insert(common.BKTableNamePrivilege, input)
	if nil != err {
		blog.Error("create role privilege error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

//UpdateRolePri update role privilege
func (cli *Service) UpdateRolePri(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	pathParams := req.PathParameters()
	ownerID := pathParams["bk_supplier_account"]
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
	err = cli.CC.InstCli.UpdateByCondition(common.BKTableNamePrivilege, input, cond)
	if nil != err {
		blog.Error("update role privilege error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}
