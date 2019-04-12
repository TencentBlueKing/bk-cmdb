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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

// CreateMainlineObjectAssociation used for create association of type bk_mainline, as it can only create by special method,
// for example add a level to business model
func (cli *Service) CreateMainlineObjectAssociation(req *restful.Request, resp *restful.Response) {
	// bk_maineline is a inner association type that can only create in special case,
	// so we separate bk_mainline association type creation with a independent method,
	enableMainlineAssociationType := true
	cli.createObjectAssociation(req, resp, enableMainlineAssociationType)
}

// CreateObjectAssociation create object association map
func (cli *Service) CreateObjectAssociation(req *restful.Request, resp *restful.Response) {
	enableMainlineAssociationType := false
	cli.createObjectAssociation(req, resp, enableMainlineAssociationType)
}

func (cli *Service) createObjectAssociation(req *restful.Request, resp *restful.Response, enableMainlineAssociationType bool) {
	// enableMainlineAssociationType used for distinguish two creation mode
	// when enableMainlineAssociationType enabled, only bk_mainline type could be create
	// when enableMainlineAssociationType disabled, all type except bk_mainline could be create

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	obj := &meta.Association{}
	if err = json.Unmarshal([]byte(value), obj); nil != err {
		blog.Errorf("fail to unmarshal json, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	if enableMainlineAssociationType == false {
		// AsstKindID shouldn't be use bk_mainline
		if obj.AsstKindID == common.AssociationKindMainline {
			blog.Errorf("use inner association type: %v is forbidden", common.AssociationKindMainline)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrorTopoAssociationKindMainlineUnavailable, obj.AsstKindID)})
			return
		}
	} else {
		// AsstKindID could only be bk_mainline
		if obj.AsstKindID != common.AssociationKindMainline {
			blog.Errorf("use CreateMainlineObjectAssociation method but bk_asst_id is: %s", obj.AsstKindID)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrorTopoAssociationKindInconsistent, obj.AsstKindID)})
			return
		}
	}

	// save to the storage
	id, err := db.NextSequence(ctx, common.BKTableNameObjAsst)
	if err != nil {
		blog.Errorf("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	obj.ID = int64(id)
	obj.OwnerID = ownerID
	err = db.Table(common.BKTableNameObjAsst).Insert(ctx, obj)
	if nil != err {
		blog.Errorf("create objectasst failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: []*meta.Association{obj}})
}

// DeleteObjectAssociation delete object association map
func (cli *Service) DeleteObjectAssociation(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	pathParameters := req.PathParameters()
	id, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Errorf("failed to get params, error info is %s ", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	// delete object from storage
	condition := map[string]interface{}{"id": id}
	if 0 == id {
		js, err := simplejson.NewFromReader(req.Request.Body)
		if err != nil {
			blog.Errorf("read http request body failed, error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
			return

		}
		condition, err = js.Map()
		if nil != err {
			blog.Errorf("fail to unmarshal json, error information is %s", err.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
			return
		}
	}

	condition = util.SetModOwner(condition, ownerID)
	cnt, cntErr := db.Table(common.BKTableNameObjAsst).Find(condition).Count(ctx)
	if nil != cntErr {
		blog.Errorf("failed to select objectasst by condition(%+v), error is %d", condition, cntErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	if 0 == cnt {
		// success
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
		return
	}

	// execute delete command
	delErr := db.Table(common.BKTableNameObjAsst).Delete(ctx, condition)
	if nil != delErr {
		blog.Errorf("fail to delete object by id , error information is %s", delErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

// UpdateObjectAssociation update object association map
func (cli *Service) UpdateObjectAssociation(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	js, err := simplejson.NewFromReader(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	pathParameters := req.PathParameters()
	id, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Errorf("failed to get params, error info is %s ", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	// decode json string
	data, jsErr := js.Map()
	if nil != jsErr {
		blog.Errorf("unmarshal json failed, error information is %s", jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	// only field in white list could be update
	// bk_asst_obj_id is allowed for add business model level
	validFields := []string{"bk_obj_asst_name", "bk_asst_obj_id"}
	validData := map[string]interface{}{}
	filterOutFields := []string{}
	for key, val := range data {
		if isValidField := util.Contains(validFields, key); isValidField == false {
			filterOutFields = append(filterOutFields, key)
			continue
		}
		validData[key] = val
	}

	if len(filterOutFields) > 0 {
		blog.Warnf("update object association got invalid fields: %v", filterOutFields)
	}
	condititon := map[string]interface{}{"id": id}
	condititon = util.SetModOwner(condititon, ownerID)
	// update object into storage
	if updateErr := db.Table(common.BKTableNameObjAsst).Update(ctx, condititon, validData); nil != updateErr {
		blog.Errorf("fail update object by condition, error information is: %v", updateErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, updateErr.Error())})
		return
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})

}

// SelectObjectAssociations search all object association map
func (cli *Service) SelectObjectAssociations(req *restful.Request, resp *restful.Response) {

	// TODO: 输入参数有变化

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// decode json object
	js, err := simplejson.NewFromReader(req.Request.Body)
	if err != nil {
		blog.Errorf("read request body failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	page := metadata.BasePage{Limit: common.BKNoLimit}
	if pageJS, ok := js.CheckGet("page"); ok {
		tmpMap, _ := pageJS.Map()
		page = meta.BasePage{}
		tmp, err := mapstr.NewFromInterface(tmpMap)
		if nil != err {
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
			return
		}
		if err := tmp.MarshalJSONInto(&tmp); nil != err {
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
			return
		}
		js.Del("page")
	}

	results := make([]meta.Association, 0)
	selector, _ := js.Map()

	// compatibility the new format. eg: {"condition":{"key":"val"}}
	if _, existsCondition := selector["condition"]; existsCondition {
		selector, _ = selector["condition"].(map[string]interface{})
	}
	selector = util.SetModOwner(selector, ownerID)
	// select from storage
	if selErr := db.Table(common.BKTableNameObjAsst).Find(selector).Limit(uint64(page.Limit)).Start(uint64(page.Start)).Sort(page.Sort).All(ctx, &results); nil != selErr && !db.IsNotFoundError(selErr) {
		blog.Errorf("select data failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, selErr.Error())})
		return
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})
}
