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
	"time"
)

// CreateInstAssociation create instance association map
func (cli *Service) CreateInstAssociation(req *restful.Request, resp *restful.Response) {

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

	request := &meta.CreateAssociationInstRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	data := &meta.InstAsst{
		ObjectAsstID: request.ObjectAsstId,
		InstID:       request.InstId,
		AsstInstID:   request.AsstInstId,
		OwnerID:      ownerID,
		CreateTime:   time.Now(),
	}

	// get id
	id, err := cli.Instance.GetIncID(common.BKTableNameInstAsst)
	if err != nil && !cli.Instance.IsNotFoundErr(err) {
		blog.Error("failed to get id , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}
	data.ID = id

	_, err = cli.Instance.Insert(common.BKTableNameInstAsst, data)
	if nil != err {
		blog.Error("search object association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	result := &meta.CreateAssociationInstResult{BaseResp: meta.SuccessBaseResp}
	result.Data.Id = id
	resp.WriteEntity(result)
}

// DeleteInstAssociation delete inst association map
func (cli *Service) DeleteInstAssociation(req *restful.Request, resp *restful.Response) {

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

	request := &meta.DeleteAssociationInstRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	if request.AsstInstID == 0 && request.InstID == 0 {
		errMsg := "invalid instance delparams"
		blog.Error(errMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommFieldNotValid, errMsg)})
		return
	}
	cond := map[string]interface{}{
		"bk_obj_asst_id":  request.ObjectAsstID,
		"bk_inst_id":      request.InstID,
		"bk_asst_inst_id": request.AsstInstID,
	}
	cond = util.SetModOwner(cond, ownerID)

	// check exist
	if cnt, _ := cli.Instance.GetCntByCondition(common.BKTableNameInstAsst, cond); cnt < 1 {
		err := fmt.Errorf("failed to delete inst association, not found")
		blog.Error(err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBUpdateFailed, err.Error())})
		return
	}

	err = cli.Instance.DelByCondition(common.BKTableNameInstAsst, cond)

	if nil != err {
		blog.Error("delete inst association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBDeleteFailed, err.Error())})
		return
	}

	result := &meta.DeleteAssociationInstResult{BaseResp: meta.SuccessBaseResp, Data: "success"}
	resp.WriteEntity(result)
}

// SearchInstAssociations search inst association map
func (cli *Service) SearchInstAssociations(req *restful.Request, resp *restful.Response) {

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

	request := &meta.SearchAssociationInstRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	cond := request.Condition
	cond = util.SetModOwner(cond, ownerID)
	result := []*meta.InstAsst{}
	err = cli.Instance.GetMutilByCondition(common.BKTableNameInstAsst, []string{}, request.Condition, &result, "", 0, 0)
	if nil != err {
		blog.Error("search association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: result})
}
