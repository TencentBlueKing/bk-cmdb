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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

// CreateObjectAtt create object's attribute
func (cli *Service) CreateObjectAtt(req *restful.Request, resp *restful.Response) {

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
		blog.Errorf("read http request failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	obj := &meta.Attribute{}

	if err = json.Unmarshal([]byte(value), obj); nil != err {
		blog.Errorf("fail to unmarshal json, error information is %s", err.Error())
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
	id, err := db.NextSequence(ctx, common.BKTableNameObjAttDes)
	if err != nil {
		blog.Errorf("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	obj.ID = int64(id)
	obj.OwnerID = ownerID
	err = db.Table(common.BKTableNameObjAttDes).Insert(ctx, obj)
	if nil != err {
		blog.Errorf("create objectatt failed, error:%s", err.Error())
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
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	pathParameters := req.PathParameters()
	id, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Errorf("failed to parse id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	// delete object from storage
	condition := map[string]interface{}{"id": id}
	if 0 == id {

		js, err := simplejson.NewFromReader(req.Request.Body)
		if err != nil {
			blog.Errorf("failed to read body,error info is %s", err.Error())
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

	// check whether propertys could delete
	propertys := []meta.Attribute{}
	cntErr := db.Table(common.BKTableNameObjAttDes).Find(condition).All(ctx, &propertys)
	if nil != cntErr {
		blog.Errorf("failed to select object by condition(%+v), error is %d", condition, cntErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return

	}
	if len(propertys) <= 0 {
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
		return
	}

	uniques, err := cli.searchObjectUnique(ctx, db, ownerID, propertys[0].ObjectID)
	if nil != err {
		blog.Errorf("failed to search object unique error: %s, params: %v %v", err, ownerID, propertys[0].ObjectID)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
		return
	}

	usedKeyID := map[int64]bool{}
	for _, unique := range uniques {
		for _, key := range unique.Keys {
			if key.Kind == meta.UinqueKeyKindProperty {
				usedKeyID[int64(key.ID)] = true
			}
		}
	}

	for index := range propertys {
		if usedKeyID[propertys[index].ID] {
			blog.Errorf("property %s has bee used by it's unique constrains, not allow delete", propertys[0].PropertyID)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Errorf(common.CCErrTopoObjectPropertyUsedByUnique, propertys[index].PropertyID)})
			return
		}
	}

	// delete propertys from db
	delErr := db.Table(common.BKTableNameObjAttDes).Delete(ctx, condition)
	if nil != delErr {
		blog.Errorf("failed to delete, error info is %s", delErr.Error())
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
		blog.Errorf("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	js.Set(common.LastTimeField, util.GetCurrentTimeStr())

	// decode json string
	data, jsErr := js.Map()
	if nil != jsErr {
		blog.Errorf("unmarshal json failed, error information is %s", jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	condition := map[string]interface{}{"id": id}
	condition = util.SetModOwner(condition, ownerID)
	// update object into storage
	updateErr := db.Table(common.BKTableNameObjAttDes).Update(ctx, condition, data)
	if nil != updateErr {
		blog.Errorf("fail update object by condition, error information is %s", updateErr.Error())
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
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	pathParameters := req.PathParameters()
	id, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Errorf("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	condition := map[string]interface{}{"id": id}
	condition = util.SetQueryOwner(condition, ownerID)
	// select from storage
	result := make([]meta.Attribute, 0)
	if selErr := db.Table(common.BKTableNameObjAttDes).Find(condition).All(ctx, &result); nil != selErr && !db.IsNotFoundError(selErr) {
		blog.Errorf("find object by selector failed, error:%s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	// translate language
	for index := range result {
		result[index].PropertyName = cli.TranslatePropertyName(defLang, &result[index])
		result[index].Placeholder = cli.TranslatePlaceholder(defLang, &result[index])
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
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// decode json object
	js, err := simplejson.NewFromReader(req.Request.Body)
	if err != nil {
		blog.Errorf("read request body failed, error information is %s", err.Error())
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

	if selErr := db.Table(common.BKTableNameObjAttDes).Find(selector).Start(uint64(page.Start)).Limit(uint64(page.Limit)).Sort(page.Sort).All(ctx, &results); nil != selErr && !db.IsNotFoundError(selErr) {
		blog.Errorf("find object by selector failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, selErr.Error())})
		return
	}

	// translate language
	for index := range results {
		results[index].PropertyName = cli.TranslatePropertyName(defLang, &results[index])
		results[index].Placeholder = cli.TranslatePlaceholder(defLang, &results[index])
		if results[index].PropertyType == common.FieldTypeEnum {
			results[index].Option = cli.TranslateEnumName(defLang, &results[index], results[index].Option)
		}
	}
	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})
}
