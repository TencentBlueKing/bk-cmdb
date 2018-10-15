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

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"fmt"
	"strconv"
)

// SearchAssociationType Search Association Type
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
	result := []*meta.AssociationType{}
	err = cli.Instance.GetMutilByCondition(common.BKTableNameAsstDes, []string{}, request.Condition, &result, request.Sort, request.Start, request.Limit)
	if nil != err {
		blog.Error("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	ret := &meta.SearchAssociationTypeResult{BaseResp: meta.SuccessBaseResp}
	ret.Data.Count = len(result)
	ret.Data.Info = result
	resp.WriteEntity(ret)
}

// CreateAssociationType Create Association Type
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

	// check uniq bk_asst_id
	cond := map[string]interface{}{"bk_asst_id": request.AsstID}
	cond = util.SetModOwner(cond, ownerID)
	cnt, err := cli.Instance.GetCntByCondition(common.BKTableNameAsstDes, cond)
	if err != nil {
		blog.Error("failed to count association , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	if cnt > 1 {
		err = fmt.Errorf("failed to create association, bk_asst_id %s exist", request.AsstID)
		blog.Error(err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	// get id
	id, err := cli.Instance.GetIncID(common.BKTableNameAsstDes)
	if err != nil && !cli.Instance.IsNotFoundErr(err) {
		blog.Error("failed to get id , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}
	request.ID = id
	request.OwnerID = ownerID

	_, err = cli.Instance.Insert(common.BKTableNameAsstDes, request)
	if nil != err {
		blog.Error("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	result := &meta.CreateAssociationTypeResult{BaseResp: meta.SuccessBaseResp}
	result.Data.Id = id
	resp.WriteEntity(result)
}

// UpdateAssociationType Update Association Type
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
	asstTypeID, _ := strconv.Atoi(id)
	request := &meta.UpdateAssociationTypeRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	cond := map[string]interface{}{"id": asstTypeID}
	cond = util.SetModOwner(cond, ownerID)
	if cnt, _ := cli.Instance.GetCntByCondition(common.BKTableNameAsstDes, cond); cnt < 1 {
		err = fmt.Errorf("failed to update association, id %d not found", asstTypeID)
		blog.Error(err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	err = cli.Instance.UpdateByCondition(common.BKTableNameAsstDes, request, cond)

	if nil != err {
		blog.Error("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBUpdateFailed, err.Error())})
		return
	}

	result := &meta.UpdateAssociationTypeResult{BaseResp: meta.SuccessBaseResp, Data: "success"}
	resp.WriteEntity(result)
}

// DeleteAssociationType Delete Association Type
func (cli *Service) DeleteAssociationType(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	id := req.PathParameter("id")
	asstTypeID, _ := strconv.Atoi(id)

	cond := map[string]interface{}{"id": asstTypeID}
	cond = util.SetModOwner(cond, ownerID)
	if cnt, _ := cli.Instance.GetCntByCondition(common.BKTableNameAsstDes, cond); cnt < 1 {
		err := fmt.Errorf("failed to delete association, id %d not found", asstTypeID)
		blog.Error(err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	err := cli.Instance.DelByCondition(common.BKTableNameAsstDes, cond)

	if nil != err {
		blog.Error("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBDeleteFailed, err.Error())})
		return
	}

	result := &meta.UpdateAssociationTypeResult{BaseResp: meta.SuccessBaseResp, Data: "success"}
	resp.WriteEntity(result)
}
