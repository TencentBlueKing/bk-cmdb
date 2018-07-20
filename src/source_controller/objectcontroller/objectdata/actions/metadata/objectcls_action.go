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

package metadata

import (
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
)

var objcls = &objectClassificationAction{}

// ObjectAction
type objectClassificationAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/object/classification/{owner_id}/objects", Params: nil, Handler: objcls.SelectClassificationWithObject})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/object/classification/search", Params: nil, Handler: objcls.SelectClassifications})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/meta/object/classification/{id}", Params: nil, Handler: objcls.DeleteClassification})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/meta/object/classification", Params: nil, Handler: objcls.CreateClassification})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/meta/object/classification/{id}", Params: nil, Handler: objcls.UpdateClassification})

	// set cc api resource
	objcls.CC = api.NewAPIResource()
}

// CreateClassification create object's classification
func (cli *objectClassificationAction) CreateClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("create classification")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		obj := &metadata.ObjClassification{}
		if err = json.Unmarshal([]byte(value), obj); nil != err {
			blog.Error("fail to unmarshal json, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		obj.OwnerID = ownerID

		// save to the storage
		id, err := cli.CC.InstCli.GetIncID(obj.TableName())
		if err != nil {
			blog.Error("failed to get id, error info is %s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		obj.ID = int(id)
		_, err = cli.CC.InstCli.Insert(obj.TableName(), obj)
		if nil != err {
			blog.Error("create objectcls failed, error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, []*metadata.ObjClassification{obj}, nil
	}, resp)
}

// DeleteClassification delete object's classification
func (cli *objectClassificationAction) DeleteClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("delete classification")
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		defer req.Request.Body.Close()
		pathParameters := req.PathParameters()
		var id int
		if err := cli.GetParams(cli.CC, &pathParameters, "id", &id, resp); nil != err {
			blog.Error("failed to get id, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsInvalid, "id")
		}

		condition := map[string]interface{}{"id": id}

		// delete object from storage
		if 0 == id {
			value, err := ioutil.ReadAll(req.Request.Body)
			if err != nil {
				blog.Error("read http request body failed, error:%s", err.Error())
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
			}

			if err := json.Unmarshal([]byte(value), &condition); nil != err {
				blog.Error("fail to unmarshal json, error information is %s", err.Error())
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
			}
		}
		condition = util.SetModOwner(condition, ownerID)
		cnt, cntErr := cli.CC.InstCli.GetCntByCondition(common.BKTableNameObjClassifiction, condition)
		if nil != cntErr {
			blog.Error("failed to select object classification by condition(%+v), error is %d", cntErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		if 0 == cnt {
			return http.StatusOK, nil, nil
		}
		// execute delete command
		if delErr := cli.CC.InstCli.DelByCondition(common.BKTableNameObjClassifiction, condition); nil != delErr {
			blog.Error("fail to delete object by id , error: %s", delErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

// UpdateClassification update object's classification information
func (cli *objectClassificationAction) UpdateClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("update classification")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParameters := req.PathParameters()
		var id int
		if err := cli.GetParams(cli.CC, &pathParameters, "id", &id, resp); nil != err {
			blog.Error("failed to get id, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsInvalid, "id")
		}

		selector := map[string]interface{}{"id": id}

		// decode json string
		data := map[string]interface{}{}
		if jsErr := json.NewDecoder(req.Request.Body).Decode(&data); nil != jsErr {
			blog.Error("unmarshal json failed, error:%v", jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		selector = util.SetModOwner(selector, ownerID)
		data = util.SetModOwner(data, ownerID)
		// update object into storage
		if updateErr := cli.CC.InstCli.UpdateByCondition(common.BKTableNameObjClassifiction, &data, selector); nil != updateErr {
			blog.Error("fail update object by condition, error:%v", updateErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

// SelectClassifications select object's classification informations
func (cli *objectClassificationAction) SelectClassifications(req *restful.Request, resp *restful.Response) {

	blog.Info("select classification")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)

	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {
		// decode json object
		selector := map[string]interface{}{}
		if jserr := json.NewDecoder(req.Request.Body).Decode(&selector); nil != jserr {
			blog.Error("unmarshal failed, error:%v", jserr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		page := metadata.ParsePage(selector["page"])
		delete(selector, "page")

		results := make([]metadata.ObjClassification, 0)

		selector = util.SetQueryOwner(selector, ownerID)
		// select from storage
		if selerr := cli.CC.InstCli.GetMutilByCondition(common.BKTableNameObjClassifiction, nil, selector, &results, page.Sort, page.Start, page.Limit); nil != selerr {
			blog.Error("select data failed, error: %s", selerr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		// translate language
		for index := range results {
			results[index].ClassificationName = commondata.TranslateClassificationName(defLang, &results[index])
		}
		return http.StatusOK, results, nil
	}, resp)
}

// SelectClassificationWithObject select objects by classification information
func (cli *objectClassificationAction) SelectClassificationWithObject(req *restful.Request, resp *restful.Response) {

	blog.Info("select classification with object")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// decode json object
		selector := map[string]interface{}{}
		if jsErr := json.NewDecoder(req.Request.Body).Decode(&selector); nil != jsErr {
			blog.Error("unmarshal failed, error: %s", jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		page := metadata.ParsePage(selector["page"])
		delete(selector, "page")
		selector = util.SetQueryOwner(selector, ownerID)

		clsResults := make([]metadata.ObjClassificationObject, 0)
		// select from storage
		if selerr := cli.CC.InstCli.GetMutilByCondition(common.BKTableNameObjClassifiction, nil, selector, &clsResults, page.Sort, page.Start, page.Limit); nil != selerr {
			blog.Error("select data failed, error:%s", selerr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}

		// select object by cls
		blog.InfoJSON("select clsresults: %s", clsResults)
		for tmpidx, tmpobj := range clsResults {
			selector := map[string]interface{}{
				"bk_classification_id": tmpobj.ClassificationID,
			}
			selector = util.SetQueryOwner(selector, ownerID)
			if selerr := cli.CC.InstCli.GetMutilByCondition(common.BKTableNameObjDes, nil, selector, &clsResults[tmpidx].Objects, "", 0, common.BKNoLimit); nil != selerr {
				blog.Error("select data failed, error:%s", selerr.Error())
				continue
			}

			if len(clsResults[tmpidx].Objects) <= 0 {
				clsResults[tmpidx].Objects = []metadata.ObjectDes{}
			}
		}

		// translate language
		for index := range clsResults {
			clsResults[index].ClassificationName = commondata.TranslateClassificationName(defLang, &clsResults[index].ObjClassification)
			for attindex := range clsResults[index].Objects {
				clsResults[index].Objects[attindex].ObjectName = commondata.TranslateObjectName(defLang, &clsResults[index].Objects[attindex])
			}
		}

		return http.StatusOK, clsResults, nil
	}, resp)
}
