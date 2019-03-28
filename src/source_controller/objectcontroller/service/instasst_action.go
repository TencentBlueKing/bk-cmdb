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
	"configcenter/src/common/mapstr"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/metadata"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

// CreateInstAssociation create instance association map
func (cli *Service) CreateInstAssociation(req *restful.Request, resp *restful.Response) {

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

	request := &meta.CreateAssociationInstRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	// find object id
	objCond := map[string]interface{}{
		meta.AssociationFieldAsstID:          request.ObjectAsstID,
		meta.AssociationFieldSupplierAccount: ownerID,
	}

	objResult := &meta.Association{}
	err = db.Table(common.BKTableNameObjAsst).Find(objCond).One(ctx, &objResult)
	if nil != err {
		blog.Errorf("not found object association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, request.ObjectAsstID)})
		return
	}

	// bk_mainline shouldn't be use
	if objResult.AsstKindID == common.AssociationKindMainline {
		blog.Errorf("use inner association type: %v is forbidden", common.AssociationKindMainline)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrorTopoAssociationKindMainlineUnavailable, request.ObjectAsstID)})
		return
	}

	// get insert id
	id, err := db.NextSequence(ctx, common.BKTableNameInstAsst)
	if err != nil {
		blog.Errorf("failed to get id , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	data := &meta.InstAsst{
		ID:                int64(id),
		ObjectID:          objResult.ObjectID,
		AsstObjectID:      objResult.AsstObjID,
		ObjectAsstID:      request.ObjectAsstID,
		InstID:            request.InstID,
		AsstInstID:        request.AsstInstID,
		AssociationKindID: objResult.AsstKindID,
		OwnerID:           ownerID,
	}

	err = db.Table(common.BKTableNameInstAsst).Insert(ctx, data)
	if nil != err {
		blog.Errorf("search object association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	result := &meta.CreateAssociationInstResult{BaseResp: meta.SuccessBaseResp}
	result.Data.ID = data.ID

	srcevent := eventclient.NewEventWithHeader(req.Request.Header)
	srcevent.EventType = metadata.EventTypeAssociation
	srcevent.ObjType = data.ObjectID
	srcevent.Action = metadata.EventActionCreate
	srcevent.Data = []metadata.EventData{
		{
			CurData: data,
		},
	}
	err = cli.EventC.Push(ctx, srcevent)
	if err != nil {
		blog.Errorf("create event error:%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Data: result.Data, Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	destevent := eventclient.NewEventWithHeader(req.Request.Header)
	destevent.EventType = metadata.EventTypeAssociation
	destevent.ObjType = data.AsstObjectID
	destevent.Action = metadata.EventActionCreate
	destevent.Data = []metadata.EventData{
		{
			CurData: data,
		},
	}
	err = cli.EventC.Push(ctx, destevent)
	if err != nil {
		blog.Errorf("create event error:%v", err)
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

	id, err := strconv.ParseInt(req.PathParameter("association_id"), 10, 64)
	if err != nil {
		blog.Errorf("delete inst association, but failed to parse association_id, err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
	}

	cond := mapstr.MapStr{
		common.BKOwnerIDField: ownerID,
		"id":                  id,
	}

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
		srcevent := eventclient.NewEventWithHeader(req.Request.Header)
		srcevent.EventType = metadata.EventTypeAssociation
		srcevent.ObjType = asst.ObjectID
		srcevent.Action = metadata.EventActionDelete
		srcevent.Data = []metadata.EventData{
			{
				PreData: asst,
			},
		}
		err = cli.EventC.Push(ctx, srcevent)
		if err != nil {
			blog.Errorf("create event error:%v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
			return
		}

		destevent := eventclient.NewEventWithHeader(req.Request.Header)
		destevent.EventType = metadata.EventTypeAssociation
		destevent.ObjType = asst.AsstObjectID
		destevent.Action = metadata.EventActionDelete
		destevent.Data = []metadata.EventData{
			{
				PreData: asst,
			},
		}
		err = cli.EventC.Push(ctx, destevent)
		if err != nil {
			blog.Errorf("create event error:%v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
			return
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
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

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

	result := []*meta.InstAsst{}
	cond := util.SetModOwner(request.Condition.ToMapInterface(), ownerID)
	if err := db.Table(common.BKTableNameInstAsst).Find(cond).All(ctx, &result); err != nil {
		blog.Errorf("select data failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: result})
}
