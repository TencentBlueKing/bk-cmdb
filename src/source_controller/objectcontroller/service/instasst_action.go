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
	"time"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/metadata"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
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
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	request := &meta.CreateAssociationInstRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
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

	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// get id
	id, err := db.NextSequence(ctx, common.BKTableNameInstAsst)
	if err != nil {
		blog.Errorf("failed to get id , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}
	data.ID = int64(id)

	err = db.Table(common.BKTableNameInstAsst).Insert(ctx, data)
	if nil != err {
		blog.Errorf("search object association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	result := &meta.CreateAssociationInstResult{BaseResp: meta.SuccessBaseResp}
	result.Data.ID = data.ID

	ec := eventclient.NewEventContextByReq(req.Request.Header, cli.Cache)
	err = ec.InsertEvent(metadata.EventTypeAssociation, data.ObjectID, metadata.EventActionCreate, data, nil)
	if err != nil {
		blog.Error("create event error:%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Data: result.Data, Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}
	err = ec.InsertEvent(metadata.EventTypeAssociation, data.AsstObjectID, metadata.EventActionCreate, data, nil)
	if err != nil {
		blog.Error("create event error:%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Data: result.Data, Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

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
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	request := &meta.DeleteAssociationInstRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	if request.AsstInstID == 0 && request.InstID == 0 {
		errMsg := "invalid instance delparams"
		blog.Errorf(errMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommFieldNotValid, errMsg)})
		return
	}
	cond := map[string]interface{}{
		"bk_obj_asst_id":  request.ObjectAsstID,
		"bk_inst_id":      request.InstID,
		"bk_asst_inst_id": request.AsstInstID,
	}
	cond = util.SetModOwner(cond, ownerID)

	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// take snapshot
	assts := []meta.InstAsst{}
	err = db.Table(common.BKTableNameInstAsst).Find(cond).All(ctx, &assts)
	if err != nil {
		blog.Errorf("failed to count inst association , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	err = db.Table(common.BKTableNameInstAsst).Delete(ctx, cond)
	if nil != err {
		blog.Errorf("delete inst association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBDeleteFailed, err.Error())})
		return
	}

	for _, asst := range assts {
		ec := eventclient.NewEventContextByReq(req.Request.Header, cli.Cache)
		err = ec.InsertEvent(metadata.EventTypeAssociation, asst.ObjectID, metadata.EventActionCreate, nil, asst)
		if err != nil {
			blog.Error("create event error:%v", err)
		}
		err = ec.InsertEvent(metadata.EventTypeAssociation, asst.AsstObjectID, metadata.EventActionCreate, nil, asst)
		if err != nil {
			blog.Error("create event error:%v", err)
		}
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
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	request := &meta.SearchAssociationInstRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	cond := request.Condition
	cond = util.SetModOwner(cond, ownerID)
	result := []*meta.InstAsst{}

	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	if err := db.Table(common.BKTableNameInstAsst).Find(cond).All(ctx, &result); err != nil {
		blog.Errorf("select data failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: result})
}
