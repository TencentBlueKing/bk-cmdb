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
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateObject create a common object
func (cli *Service) CreateObject(req *restful.Request, resp *restful.Response) {

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

	obj := &meta.Object{}
	if jsErr := json.Unmarshal([]byte(value), obj); nil != jsErr {
		blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	// save to the storage
	obj.CreateTime = new(time.Time)
	*obj.CreateTime = time.Now()
	obj.LastTime = new(time.Time)
	*obj.LastTime = time.Now()
	obj.OwnerID = ownerID

	// get id
	id, err := cli.Instance.GetIncID(common.BKTableNameObjDes)
	if err != nil && !cli.Instance.IsNotFoundErr(err) {
		blog.Error("failed to get id , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	obj.ID = id

	// save
	_, err = cli.Instance.Insert(common.BKTableNameObjDes, obj)
	if nil == err && !cli.Instance.IsNotFoundErr(err) {
		resp.WriteEntity(meta.CreateObjectResult{BaseResp: meta.SuccessBaseResp, Data: meta.RspID{ID: id}})
		return
	}
	blog.Error("failed to insert the object, error info is %s", err.Error())

	resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})

}

//删除Object
func (cli *Service) DeleteObject(req *restful.Request, resp *restful.Response) {

	blog.Info("delete object")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	pathParameters := req.PathParameters()
	appID, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Error("failed to get params, error info is %s ", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	condition := map[string]interface{}{"id": appID}
	// delete object from storage
	if 0 == appID {
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
	cnt, cntErr := cli.Instance.GetCntByCondition(common.BKTableNameObjDes, condition)
	if nil != cntErr && !cli.Instance.IsNotFoundErr(cntErr) {
		blog.Error("failed to select object by condition(%+v), error is %d", cntErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, cntErr.Error())})
		return
	}
	if 0 == cnt {
		// success
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
		return
	}
	// execute delete command
	if delErr := cli.Instance.DelByCondition(common.BKTableNameObjDes, condition); nil != delErr && !cli.Instance.IsNotFoundErr(delErr) {
		blog.Error("fail to delete object by id , error information is %s", delErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, delErr.Error())})
		return
	}
	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

func (cli *Service) UpdateObject(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	js, err := simplejson.NewFromReader(req.Request.Body)
	if err != nil {
		blog.Error("read request body failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	pathParameters := req.PathParameters()
	appID, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Error("failed to get params, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	// update object into storage
	js.Set(common.LastTimeField, util.GetCurrentTimeStr())

	// decode json string
	data, jsErr := js.Map()
	if nil != jsErr {
		blog.Error("unmarshal json failed, error information is %v", jsErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}
	condition := map[string]interface{}{"id": appID}
	condition = util.SetModOwner(condition, ownerID)
	err = cli.Instance.UpdateByCondition(common.BKTableNameObjDes, data, condition)
	if nil != err && !cli.Instance.IsNotFoundErr(err) {
		blog.Error("fail update object by condition, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, jsErr.Error())})
		return
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})

}

//查询所有主机信息
func (cli *Service) SelectObjects(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	defLang := cli.Core.Language.CreateDefaultCCLanguageIf(language)

	// decode json object
	js, err := simplejson.NewFromReader(req.Request.Body)
	if err != nil {
		blog.Error("read request body failed, error information is %s", err.Error())
		if nil != err {
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
			return
		}
	}

	page := meta.BasePage{Limit: common.BKNoLimit}
	if pageJS, ok := js.CheckGet("page"); ok {
		tmpMap, err := pageJS.Map()
		if nil != err {
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
			return
		}
		page = meta.ParsePage(tmpMap)
		js.Del("page")
	}
	results := make([]meta.Object, 0)

	// select from storage

	selector, err := js.Map()
	if nil != err {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	selector = util.SetQueryOwner(selector, ownerID)
	if selErr := cli.Instance.GetMutilByCondition(common.BKTableNameObjDes, nil, selector, &results, page.Sort, page.Start, page.Limit); nil != selErr {
		blog.Error("select data failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, selErr.Error())})
		return
	}

	// translate language
	for index := range results {
		results[index].ObjectName = cli.TranslateObjectName(defLang, &results[index])
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})

}
