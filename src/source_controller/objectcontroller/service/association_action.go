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
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

// SearchAssociationType Search Association Type
func (cli *Service) SearchAssociationType(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	defLang := cli.Core.Language.CreateDefaultCCLanguageIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	request := &meta.SearchAssociationTypeRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	cond := request.Condition
	cond = util.SetQueryOwner(cond, ownerID)
	result := []*meta.AssociationKind{}

	if selErr := db.Table(common.BKTableNameAsstDes).Find(cond).Limit(uint64(request.Limit)).Start(uint64(request.Start)).Sort(request.Sort).All(ctx, &result); nil != selErr && !db.IsNotFoundError(selErr) {
		blog.Errorf("select data failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, selErr.Error())})
		return
	}
	cnt, err := db.Table(common.BKTableNameAsstDes).Find(cond).Count(ctx)
	if nil != err {
		blog.Errorf("select data failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	for index := range result {
		cli.TranslateAssociationKind(defLang, result[index])
	}

	ret := &meta.SearchAssociationTypeResult{BaseResp: meta.SuccessBaseResp}
	ret.Data.Count = int(cnt)
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
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	request := &meta.AssociationKind{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	// check uniq bk_asst_id
	cond := map[string]interface{}{"bk_asst_id": request.AssociationKindID}
	cond = util.SetModOwner(cond, ownerID)

	cnt, err := db.Table(common.BKTableNameAsstDes).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("failed to count association , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	if cnt > 0 {
		msg := fmt.Sprintf("failed to create association, bk_asst_id %s exist", request.AssociationKindID)
		blog.Errorf(msg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, msg)})
		return
	}

	// get id
	id, err := db.NextSequence(ctx, common.BKTableNameAsstDes)
	if err != nil {
		blog.Errorf("failed to get id , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	request.ID = int64(id)
	request.OwnerID = ownerID

	err = db.Table(common.BKTableNameAsstDes).Insert(ctx, request)
	if nil != err {
		blog.Errorf("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	result := &meta.CreateAssociationTypeResult{BaseResp: meta.SuccessBaseResp}
	result.Data.ID = int64(id)
	resp.WriteEntity(result)
}

// UpdateAssociationType Update Association Type
func (cli *Service) UpdateAssociationType(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	id := req.PathParameter("id")
	asstTypeID, _ := strconv.Atoi(id)
	request := &meta.UpdateAssociationTypeRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	cond := map[string]interface{}{"id": asstTypeID}
	cond = util.SetModOwner(cond, ownerID)
	cnt, err := db.Table(common.BKTableNameAsstDes).Find(cond).Count(ctx)

	if err != nil {
		blog.Errorf("failed to count association , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	if cnt < 1 {
		msg := fmt.Sprintf("failed to update association, id %d not found", asstTypeID)
		blog.Errorf(msg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, msg)})
		return
	}

	err = db.Table(common.BKTableNameAsstDes).Update(ctx, cond, request)
	if nil != err {
		blog.Errorf("search association error :%v", err)
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

	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	id := req.PathParameter("id")
	asstTypeID, _ := strconv.Atoi(id)

	cond := map[string]interface{}{"id": asstTypeID}
	cond = util.SetModOwner(cond, ownerID)

	cnt, err := db.Table(common.BKTableNameAsstDes).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf(err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	if cnt < 1 {
		msg := fmt.Sprintf("failed to delete association, id %d not found", asstTypeID)
		blog.Errorf(msg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, msg)})
		return
	}

	err = db.Table(common.BKTableNameAsstDes).Delete(ctx, cond)
	if nil != err {
		blog.Errorf("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBDeleteFailed, err.Error())})
		return
	}

	result := &meta.UpdateAssociationTypeResult{BaseResp: meta.SuccessBaseResp, Data: "success"}
	resp.WriteEntity(result)
}
