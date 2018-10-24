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

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateObjectAssociation create object association map
func (cli *Service) CreateObjectAssociation(req *restful.Request, resp *restful.Response) {

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

	request := &meta.Association{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	request.OwnerID = ownerID

	// check uniq bk_obj_asst_id
	if request.ObjectAsstID == "" {
		msg := fmt.Sprintf("failed to create object association, bk_obj_asst_id must be set")
		blog.Errorf(msg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, msg)})
		return
	}

	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// check uniq
	cond := map[string]interface{}{"bk_obj_asst_id": request.ObjectAsstID}
	cond = util.SetModOwner(cond, ownerID)

	cnt, err := db.Table(common.BKTableNameObjAsst).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("failed to count object association , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	if cnt > 0 {
		msg := fmt.Sprintf("failed to create object association, bk_obj_asst_id %s exist", request.ObjectAsstID)
		blog.Errorf(msg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, msg)})
		return
	}

	// get id
	id, err := db.NextSequence(ctx, common.BKTableNameObjAsst)
	if err != nil {
		blog.Errorf("failed to get id , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}
	request.ID = int64(id)

	err = db.Table(common.BKTableNameObjAsst).Insert(ctx, request)
	if nil != err {
		blog.Errorf("search object association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBInsertFailed, err.Error())})
		return
	}

	result := &meta.CreateAssociationObjectResult{BaseResp: meta.SuccessBaseResp}
	result.Data.ID = request.ID
	resp.WriteEntity(result)
}

// DeleteObjectAssociation delete object association map
func (cli *Service) DeleteObjectAssociation(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	id := req.PathParameter("id")
	ID, _ := strconv.Atoi(id)

	cond := map[string]interface{}{"id": ID}
	cond = util.SetModOwner(cond, ownerID)

	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// check exist
	cnt, err := db.Table(common.BKTableNameObjAsst).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("failed to count object association , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	if cnt < 1 {
		msg := fmt.Sprintf("failed to delete object association, id %d not found", ID)
		blog.Errorf(msg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, msg)})
		return
	}

	err = db.Table(common.BKTableNameObjAsst).Delete(ctx, cond)
	if nil != err {
		blog.Errorf("delete object association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBDeleteFailed, err.Error())})
		return
	}

	result := &meta.DeleteAssociationObjectResult{BaseResp: meta.SuccessBaseResp, Data: "success"}
	resp.WriteEntity(result)
}

// UpdateObjectAssociation update object association map
func (cli *Service) UpdateObjectAssociation(req *restful.Request, resp *restful.Response) {

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

	id := req.PathParameter("id")
	ID, _ := strconv.Atoi(id)

	request := &meta.UpdateAssociationObjectRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	cond := map[string]interface{}{"id": ID}
	cond = util.SetModOwner(cond, ownerID)

	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// check exist
	cnt, err := db.Table(common.BKTableNameObjAsst).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("failed to count object association , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, err.Error())})
		return
	}

	if cnt < 1 {
		msg := fmt.Sprintf("failed to update object association, id %d not found", ID)
		blog.Errorf(msg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommNotFound, msg)})
		return
	}

	err = db.Table(common.BKTableNameObjAsst).Update(ctx, cond, request)
	if nil != err {
		blog.Errorf("update object association error :%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBUpdateFailed, err.Error())})
		return
	}

	result := &meta.UpdateAssociationObjectResult{BaseResp: meta.SuccessBaseResp, Data: "success"}
	resp.WriteEntity(result)

}

// SelectObjectAssociations search all object association map
func (cli *Service) SelectObjectAssociations(req *restful.Request, resp *restful.Response) {

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

	request := &meta.SearchAssociationObjectRequest{}
	if jsErr := json.Unmarshal([]byte(value), request); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	cond := map[string]interface{}{
		"bk_asst_id":     request.Condition.AsstID,
		"bk_obj_id":      request.Condition.ObjectID,
		"bk_asst_obj_id": request.Condition.AsstObjID,
	}
	cond = util.SetModOwner(cond, ownerID)

	if request.Condition.BothObjectID != "" {
		cond["$or"] = []map[string]interface{}{
			{
				"bk_object_id": request.Condition.ObjectID,
			},
			{
				"bk_asst_object_id": request.Condition.AsstObjID,
			},
		}
	}

	result := []*meta.Association{}

	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	if err := db.Table(common.BKTableNameObjAsst).Find(cond).All(ctx, &result); err != nil {
		blog.Errorf("select data failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommDBSelectFailed, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: result})
}
