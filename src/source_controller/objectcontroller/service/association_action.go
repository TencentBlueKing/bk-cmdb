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
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"strconv"
)

func (cli *Service) SearchAssociationType(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	request := &meta.SearchAssociationTypeRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	cond := request.Condition
	cond = util.SetModOwner(cond, ownerID)
	result := []*meta.Association{}
	err = cli.Instance.GetMutilByCondition(common.BKTableNameAsstDes, []string{}, request.Condition, result, request.Sort, request.Start, request.Limit)
	if nil != err {
		blog.Error("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: result})

}

func (cli *Service) CreateAssociationType(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	request := &meta.AssociationType{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	request.OwnerID = ownerID
	id, err := cli.Instance.Insert(common.BKTableNameAsstDes, request)
	if nil != err {
		blog.Error("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	result := &meta.CreateAssociationTypeResult{meta.SuccessBaseResp, struct{ Id int }{Id: id}}
	resp.WriteEntity(result)
}

func (cli *Service) UpdateAssociationType(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	id := req.PathParameter("id")
	asstTypeId, _ := strconv.Atoi(id)
	request := &meta.UpdateAssociationTypeRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	cond := map[string]interface{}{"id": asstTypeId}
	cond = util.SetModOwner(cond, ownerID)
	err = cli.Instance.UpdateByCondition(common.BKTableNameAsstDes, request, cond)

	if nil != err {
		blog.Error("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBUpdateFailed, err.Error())})
		return
	}

	result := &meta.UpdateAssociationTypeResult{meta.SuccessBaseResp, "success"}
	resp.WriteEntity(result)
}

func (cli *Service) DeleteAssociationType(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	id := req.PathParameter("id")
	asstTypeId, _ := strconv.Atoi(id)

	cond := map[string]interface{}{"id": asstTypeId}
	cond = util.SetModOwner(cond, ownerID)
	err := cli.Instance.DelByCondition(common.BKTableNameAsstDes, cond)

	if nil != err {
		blog.Error("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBDeleteFailed, err.Error())})
		return
	}

	result := &meta.UpdateAssociationTypeResult{meta.SuccessBaseResp, "success"}
	resp.WriteEntity(result)
}
