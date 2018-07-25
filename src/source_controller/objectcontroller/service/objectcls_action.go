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

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateClassification create object's classification
func (cli *Service) CreateClassification(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)

	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	// execute

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read http request body failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}
	obj := &meta.Classification{}
	if err = json.Unmarshal([]byte(value), obj); nil != err {
		blog.Error("fail to unmarshal json, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return

	}

	// save to the storage
	id, err := cli.Instance.GetIncID(common.BKTableNameObjClassifiction)
	if err != nil {
		blog.Error("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	obj.ID = id
	obj.OwnerID = ownerID
	_, err = cli.Instance.Insert(common.BKTableNameObjClassifiction, obj)
	if nil != err {
		blog.Error("create objectcls failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: obj})
}

// DeleteClassification delete object's classification
func (cli *Service) DeleteClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("delete classification")
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	// execute

	defer req.Request.Body.Close()
	pathParameters := req.PathParameters()
	id, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Error("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	condition := map[string]interface{}{"id": id}

	// delete object from storage
	if 0 == id {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
			return
		}

		if err := json.Unmarshal([]byte(value), &condition); nil != err {
			blog.Error("fail to unmarshal json, error information is %s", err.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
			return
		}
	}
	condition = util.SetModOwner(condition, ownerID)
	cnt, cntErr := cli.Instance.GetCntByCondition(common.BKTableNameObjClassifiction, condition)
	if nil != cntErr {
		blog.Error("failed to select object classification by condition(%+v), error is %d", cntErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	if 0 == cnt {
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
		return
	}
	// execute delete command
	if delErr := cli.Instance.DelByCondition(common.BKTableNameObjClassifiction, condition); nil != delErr {
		blog.Error("fail to delete object by id , error: %s", delErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return

	}
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

// UpdateClassification update object's classification information
func (cli *Service) UpdateClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("update classification")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	// execute

	pathParameters := req.PathParameters()
	id, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Error("failed to get id, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	selector := map[string]interface{}{"id": id}

	// decode json string
	data := map[string]interface{}{}
	if jsErr := json.NewDecoder(req.Request.Body).Decode(&data); nil != jsErr {
		blog.Error("unmarshal json failed, error:%v", jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

	selector = util.SetModOwner(selector, ownerID)
	// update object into storage
	if updateErr := cli.Instance.UpdateByCondition(common.BKTableNameObjClassifiction, &data, selector); nil != updateErr {
		blog.Error("fail update object by condition, error:%v", updateErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, updateErr.Error())})
		return
	}
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

// SelectClassifications select object's classification informations
func (cli *Service) SelectClassifications(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	defLang := cli.Core.Language.CreateDefaultCCLanguageIf(language)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	// execute

	// decode json object
	selector := map[string]interface{}{}
	if jserr := json.NewDecoder(req.Request.Body).Decode(&selector); nil != jserr {
		blog.Error("unmarshal failed, error:%v", jserr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jserr.Error())})
		return
	}
	page := meta.ParsePage(selector["page"])
	delete(selector, "page")

	results := make([]meta.Classification, 0)

	selector = util.SetQueryOwner(selector, ownerID)
	// select from storage
	if selerr := cli.Instance.GetMutilByCondition(common.BKTableNameObjClassifiction, nil, selector, &results, page.Sort, page.Start, page.Limit); nil != selerr {
		blog.Error("select data failed, error: %s", selerr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, selerr.Error())})
		return
	}
	// translate language
	for index := range results {
		results[index].ClassificationName = cli.TranslateClassificationName(defLang, &results[index])
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})

}

// SelectClassificationWithObject select objects by classification information
func (cli *Service) SelectClassificationWithObject(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	defLang := cli.Core.Language.CreateDefaultCCLanguageIf(language)

	// execute

	// decode json object
	selector := map[string]interface{}{}
	if jsErr := json.NewDecoder(req.Request.Body).Decode(&selector); nil != jsErr {
		blog.Error("unmarshal failed, error: %s", jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}
	page := meta.ParsePage(selector["page"])
	delete(selector, "page")

	clsResults := make([]meta.ObjClassificationObject, 0)
	selector = util.SetQueryOwner(selector, ownerID)
	// select from storage
	if selerr := cli.Instance.GetMutilByCondition(common.BKTableNameObjClassifiction, nil, selector, &clsResults, page.Sort, page.Start, page.Limit); nil != selerr {
		blog.Error("select data failed, error:%s", selerr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, selerr.Error())})
		return
	}

	// select object by cls
	blog.InfoJSON("select clsresults: %s", clsResults)
	for tmpidx, tmpobj := range clsResults {
		selector := map[string]interface{}{
			"bk_classification_id": tmpobj.ClassificationID,
			common.BKOwnerIDField:  ownerID,
		}
		selector = util.SetQueryOwner(selector, ownerID)
		if selerr := cli.Instance.GetMutilByCondition(common.BKTableNameObjDes, nil, selector, &clsResults[tmpidx].Objects, "", 0, common.BKNoLimit); nil != selerr {
			blog.Error("select data failed, error:%s", selerr.Error())
			continue
		}

		if len(clsResults[tmpidx].Objects) <= 0 {
			clsResults[tmpidx].Objects = []meta.Object{}
		}
	}

	// translate language
	for index := range clsResults {
		clsResults[index].ClassificationName = cli.TranslateClassificationName(defLang, &clsResults[index].Classification)
		for attindex := range clsResults[index].Objects {
			clsResults[index].Objects[attindex].ObjectName = cli.TranslateObjectName(defLang, &clsResults[index].Objects[attindex])
		}
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: clsResults})

}
