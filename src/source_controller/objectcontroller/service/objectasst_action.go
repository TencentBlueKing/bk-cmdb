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
	"strconv"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateObjectAssociation create object association map
func (cli *Service) CreateObjectAssociation(req *restful.Request, resp *restful.Response) {

	blog.Info("create obj-association")
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read http request body failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	obj := &meta.Association{}
	if err = json.Unmarshal([]byte(value), obj); nil != err {
		blog.Error("fail to unmarshal json, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	// save to the storage
	id, err := cli.Instance.GetIncID("cc_ObjAsst")
	if err != nil && !cli.Instance.IsNotFoundErr(err) {
		blog.Error("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	obj.ID = id
	obj.OwnerID = ownerID
	_, err = cli.Instance.Insert("cc_ObjAsst", obj)
	if nil != err {
		blog.Error("create objectasst failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: []*meta.Association{obj}})
}

// DeleteObjectAssociation delete object association map
func (cli *Service) DeleteObjectAssociation(req *restful.Request, resp *restful.Response) {

	blog.Info("delete obj-association")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	pathParameters := req.PathParameters()
	id, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Error("failed to get params, error info is %s ", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	// delete object from storage
	condition := map[string]interface{}{"id": id}
	if 0 == id {
		js, err := simplejson.NewFromReader(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
			return

		}
		condition, err = js.Map()
		if nil != err {
			blog.Error("fail to unmarshal json, error information is %s", err.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
			return
		}
	}

	condition = util.SetModOwner(condition, ownerID)
	cnt, cntErr := cli.Instance.GetCntByCondition("cc_ObjAsst", condition)
	if nil != cntErr {
		blog.Error("failed to select objectasst by condition(%+v), error is %d", cntErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	if 0 == cnt {
		// success
		// success
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
		return
	}

	// execute delete command
	delErr := cli.Instance.DelByCondition("cc_ObjAsst", condition)
	if nil != delErr && !cli.Instance.IsNotFoundErr(err) {
		blog.Error("fail to delete object by id , error information is %s", delErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

// UpdateObjectAssociation update object association map
func (cli *Service) UpdateObjectAssociation(req *restful.Request, resp *restful.Response) {

	blog.Info("update object association")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	js, err := simplejson.NewFromReader(req.Request.Body)
	if err != nil {
		blog.Error("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	pathParameters := req.PathParameters()
	id, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Error("failed to get params, error info is %s ", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	// decode json string
	data, jsErr := js.Map()
	if nil != jsErr {
		blog.Error("unmarshal json failed, error information is %s", jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}
	condititon := map[string]interface{}{"id": id}
	condititon = util.SetModOwner(condititon, ownerID)
	// update object into storage
	if updateErr := cli.Instance.UpdateByCondition("cc_ObjAsst", data, condititon); nil != updateErr {
		blog.Error("fail update object by condition, error information is %s", updateErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, updateErr.Error())})
		return
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})

}

// SelectObjectAssociations search all object association map
func (cli *Service) SelectObjectAssociations(req *restful.Request, resp *restful.Response) {

	blog.Info("search object association")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	// decode json object
	js, err := simplejson.NewFromReader(req.Request.Body)
	if err != nil {
		blog.Error("read request body failed, error information is %s", err.Error())
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
	selector = util.SetModOwner(selector, ownerID)
	// select from storage
	if selErr := cli.Instance.GetMutilByCondition("cc_ObjAsst", nil, selector, &results, page.Sort, page.Start, page.Limit); nil != selErr && !cli.Instance.IsNotFoundErr(selErr) {
		blog.Error("select data failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, selErr.Error())})
		return
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})
}
