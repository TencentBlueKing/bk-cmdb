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

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateObjectAtt create object's attribute
func (cli *Service) CreateObjectAtt(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	defLang := cli.Core.Language.CreateDefaultCCLanguageIf(language)

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read http request failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	obj := &meta.Attribute{}

	if err = json.Unmarshal([]byte(value), obj); nil != err {
		blog.Error("fail to unmarshal json, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	// save to the storage
	obj.CreateTime = new(time.Time)
	*obj.CreateTime = time.Now()
	obj.LastTime = new(time.Time)
	*obj.LastTime = time.Now()

	if obj.IsPre {
		if obj.PropertyID == common.BKInstNameField {
			obj.PropertyName = util.FirstNotEmptyString(defLang.Language("common_property_"+obj.PropertyID), obj.PropertyName, obj.PropertyID)
		}
	}

	if 0 == len(obj.PropertyGroup) {
		obj.PropertyGroup = "default" // empty value
	}
	if 0 >= obj.PropertyIndex {
		obj.PropertyIndex = -1 // not set any value
	}
	id, err := cli.Instance.GetIncID("cc_ObjAttDes")
	if err != nil && !cli.Instance.IsNotFoundErr(err) {
		blog.Errorf("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	obj.ID = id
	obj.OwnerID = ownerID
	_, err = cli.Instance.Insert("cc_ObjAttDes", obj)
	if nil != err {
		blog.Error("create objectatt failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: obj})

}

// DeleteObjectAttByID delete object's attribute by id
func (cli *Service) DeleteObjectAttByID(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	pathParameters := req.PathParameters()
	appID, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Error("failed to parse id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	// delete object from storage
	condition := map[string]interface{}{"id": appID}
	if 0 == appID {

		js, err := simplejson.NewFromReader(req.Request.Body)
		if err != nil {
			blog.Error("failed to read body,error info is %s", err.Error())
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
	cnt, cntErr := cli.Instance.GetCntByCondition("cc_ObjAttDes", condition)
	if nil != cntErr && !cli.Instance.IsNotFoundErr(cntErr) {
		blog.Error("failed to select object by condition(%+v), error is %d", cntErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return

	}
	if 0 == cnt {

		// success
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
		return
	}
	delErr := cli.Instance.DelByCondition("cc_ObjAttDes", condition)
	if nil != delErr && !cli.Instance.IsNotFoundErr(delErr) {
		blog.Error("failed to delete, error info is %s", delErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

// UpdateObjectAttByID update object's attribute by id
func (cli *Service) UpdateObjectAttByID(req *restful.Request, resp *restful.Response) {

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
	appID, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Error("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	js.Set(common.LastTimeField, util.GetCurrentTimeStr())

	// decode json string
	data, jsErr := js.Map()
	if nil != jsErr {
		blog.Error("unmarshal json failed, error information is %s", jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	condition := map[string]interface{}{"id": appID}
	condition = util.SetModOwner(condition, ownerID)
	// update object into storage
	updateErr := cli.Instance.UpdateByCondition("cc_ObjAttDes", data, condition)
	if nil != updateErr {
		blog.Error("fail update object by condition, error information is %s", updateErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, updateErr.Error())})
		return
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})

}

// SelectObjectAttByID select object's attribute by id
func (cli *Service) SelectObjectAttByID(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	defLang := cli.Core.Language.CreateDefaultCCLanguageIf(language)

	pathParameters := req.PathParameters()
	id, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Error("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	condition := map[string]interface{}{"id": id}
	condition = util.SetQueryOwner(condition, ownerID)
	// select from storage
	result := make([]meta.Attribute, 0)
	if selErr := cli.Instance.GetMutilByCondition("cc_ObjAttDes", nil, condition, &result, "", 0, 0); nil != selErr {
		blog.Error("find object by selector failed, error:%s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	// translate language
	for index := range result {
		result[index].PropertyName = cli.TranslatePropertyName(defLang, &result[index])
		if result[index].PropertyType == common.FieldTypeEnum {
			result[index].Option = cli.TranslateEnumName(defLang, &result[index], result[index].Option)
		}
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: result})
}

// SelectObjectAttWithParams select object's attribute with some params
func (cli *Service) SelectObjectAttWithParams(req *restful.Request, resp *restful.Response) {

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
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	page := meta.BasePage{Limit: common.BKNoLimit}
	if pageJS, ok := js.CheckGet("page"); ok {
		tmpMap, _ := pageJS.Map()
		page = meta.ParsePage(tmpMap)
		js.Del("page")
	}

	results := make([]meta.Attribute, 0)
	// select from storage
	selector, err := js.Map()
	if nil != err {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}
	selector = util.SetQueryOwner(selector, ownerID)

	if selErr := cli.Instance.GetMutilByCondition("cc_ObjAttDes", nil, selector, &results, page.Sort, page.Start, page.Limit); nil != selErr {
		blog.Error("find object by selector failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, selErr.Error())})
		return
	}

	// translate language
	for index := range results {
		results[index].PropertyName = cli.TranslatePropertyName(defLang, &results[index])
		if results[index].PropertyType == common.FieldTypeEnum {
			results[index].Option = cli.TranslateEnumName(defLang, &results[index], results[index].Option)
		}
	}
	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})
}
