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
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"

	"github.com/emicklei/go-restful"
)

var metaObjectAtt = &objectAttAction{}

// ObjectAction
type objectAttAction struct {
	base.BaseAction
}

func init() {

	// register action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/objectatt/{id}", Params: nil, Handler: metaObjectAtt.SelectObjectAttByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/objectatts", Params: nil, Handler: metaObjectAtt.SelectObjectAttWithParams})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/meta/objectatt/{id}", Params: nil, Handler: metaObjectAtt.DeleteObjectAttByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/meta/objectatt", Params: nil, Handler: metaObjectAtt.CreateObjectAtt})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/meta/objectatt/{id}", Params: nil, Handler: metaObjectAtt.UpdateObjectAttByID})

	// set cc api resource
	metaObjectAtt.CC = api.NewAPIResource()
}

// CreateObjectAtt create object's attribute
func (cli *objectAttAction) CreateObjectAtt(req *restful.Request, resp *restful.Response) {

	blog.Info("create objectatt")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request failed, error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		obj := &metadata.ObjectAttDes{}

		if err = json.Unmarshal([]byte(value), obj); nil != err {
			blog.Error("fail to unmarshal json, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		// save to the storage
		obj.CreateTime = new(time.Time)
		*obj.CreateTime = time.Now()
		obj.LastTime = new(time.Time)
		*obj.LastTime = time.Now()
		obj.OwnerID = ownerID

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
		id, err := cli.CC.InstCli.GetIncID(obj.TableName())
		if err != nil {
			blog.Errorf("failed to get id, error info is %s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		obj.ID = int(id)
		_, err = cli.CC.InstCli.Insert(obj.TableName(), obj)
		if nil != err {
			blog.Error("create objectatt failed, error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}

		return http.StatusOK, []*metadata.ObjectAttDes{obj}, nil
	}, resp)
}

// DeleteObjectAttByID delete object's attribute by id
func (cli *objectAttAction) DeleteObjectAttByID(req *restful.Request, resp *restful.Response) {

	blog.Info("delete objectatt ")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParameters := req.PathParameters()
		var appID int
		if err := cli.GetParams(cli.CC, &pathParameters, "id", &appID, resp); nil != err {
			blog.Error("failed to parse id, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsInvalid, "id")
		}

		// delete object from storage
		condition := map[string]interface{}{"id": appID}
		if 0 == appID {

			js, err := simplejson.NewFromReader(req.Request.Body)
			if err != nil {
				blog.Error("failed to read body,error info is %s", err.Error())
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
			}
			condition, err = js.Map()
			if nil != err {
				blog.Error("fail to unmarshal json, error information is %s", err.Error())
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
			}
		}
		condition = util.SetModOwner(condition, ownerID)
		cnt, cntErr := cli.CC.InstCli.GetCntByCondition(metadata.ObjectAttDes{}.TableName(), condition)
		if nil != cntErr {
			blog.Error("failed to select object by condition(%+v), error is %d", cntErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		if 0 == cnt {
			// success
			return http.StatusOK, nil, nil
		}
		delErr := cli.CC.InstCli.DelByCondition(metadata.ObjectAttDes{}.TableName(), condition)
		if nil != delErr {
			blog.Error("failed to delete, error info is %s", delErr.Error())

		}

		// success
		return http.StatusOK, nil, nil
	}, resp)
}

// UpdateObjectAttByID update object's attribute by id
func (cli *objectAttAction) UpdateObjectAttByID(req *restful.Request, resp *restful.Response) {

	blog.Info("update objectatt")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		js, err := simplejson.NewFromReader(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		pathParameters := req.PathParameters()
		var appID int
		if err := cli.GetParams(cli.CC, &pathParameters, "id", &appID, resp); nil != err {
			blog.Error("failed to get id, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsInvalid, "id")
		}

		js.Set(common.LastTimeField, util.GetCurrentTimeStr())

		// decode json string
		data, jsErr := js.Map()
		if nil != jsErr {
			blog.Error("unmarshal json failed, error information is %s", jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		condition := util.SetModOwner(map[string]interface{}{"id": appID}, ownerID)
		// update object into storage
		updateErr := cli.CC.InstCli.UpdateByCondition(metadata.ObjectAttDes{}.TableName(), data, condition)
		if nil != updateErr {
			blog.Error("fail update object by condition, error information is %s", updateErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}

		// success
		return http.StatusOK, nil, nil
	}, resp)

}

// SelectObjectAttByID select object's attribute by id
func (cli *objectAttAction) SelectObjectAttByID(req *restful.Request, resp *restful.Response) {

	blog.Info("select objectatt by id")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParameters := req.PathParameters()
		var id int
		if err := cli.GetParams(cli.CC, &pathParameters, "id", &id, resp); nil != err {
			blog.Error("failed to get id, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsInvalid, "id")
		}

		condition := util.SetQueryOwner(map[string]interface{}{"id": id}, ownerID)
		// select from storage
		result := make([]metadata.ObjectAttDes, 0)
		if selErr := cli.CC.InstCli.GetMutilByCondition(metadata.ObjectAttDes{}.TableName(), nil, condition, &result, "", 0, 0); nil != selErr {
			blog.Error("find object by selector failed, error:%s", selErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		// translate language
		for index := range result {
			result[index].PropertyName = commondata.TranslatePropertyName(defLang, &result[index])
			if result[index].PropertyType == common.FieldTypeEnum {
				result[index].Option = commondata.TranslateEnumName(defLang, &result[index], result[index].Option)
			}
		}

		// success
		return http.StatusOK, result, nil

	}, resp)
}

// SelectObjectAttWithParams select object's attribute with some params
func (cli *objectAttAction) SelectObjectAttWithParams(req *restful.Request, resp *restful.Response) {

	blog.Info("select objectatts with params")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		// decode json object
		js, err := simplejson.NewFromReader(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		page := metadata.BasePage{Limit: common.BKNoLimit}
		if pageJS, ok := js.CheckGet("page"); ok {
			tmpMap, _ := pageJS.Map()
			page = metadata.ParsePage(tmpMap)
			js.Del("page")
		}

		results := make([]metadata.ObjectAttDes, 0)
		// select from storage
		selector, _ := js.Map()
		selector = util.SetQueryOwner(selector, ownerID)
		blog.Debug("the condition: %+v the page:%+v", selector, page)
		if selErr := cli.CC.InstCli.GetMutilByCondition(metadata.ObjectAttDes{}.TableName(), nil, selector, &results, page.Sort, page.Start, page.Limit); nil != selErr {
			blog.Error("find object by selector failed, error information is %s", selErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		// blog.Debug("the result:%+v", results)
		// translate language
		for index := range results {
			results[index].PropertyName = commondata.TranslatePropertyName(defLang, &results[index])
			if results[index].PropertyType == common.FieldTypeEnum {
				results[index].Option = commondata.TranslateEnumName(defLang, &results[index], results[index].Option)
			}
		}
		// success
		return http.StatusOK, results, nil
	}, resp)
}
