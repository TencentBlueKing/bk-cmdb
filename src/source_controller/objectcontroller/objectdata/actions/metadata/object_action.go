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
	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
	"time"
)

var objaction = &objectAction{}

// ObjectAction
type objectAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/objects", Params: nil, Handler: objaction.SelectObjects})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/meta/object/{id}", Params: nil, Handler: objaction.DeleteObject})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/meta/object", Params: nil, Handler: objaction.CreateObject})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/meta/object/{id}", Params: nil, Handler: objaction.UpdateObject})

	// set cc api resource
	objaction.CC = api.NewAPIResource()
}

// CreateObject create a common object
func (cli *objectAction) CreateObject(req *restful.Request, resp *restful.Response) {

	blog.Info("create object")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		obj := &metadata.ObjectDes{}
		if jsErr := json.Unmarshal([]byte(value), obj); nil != jsErr {
			blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		// save to the storage
		obj.CreateTime = new(time.Time)
		*obj.CreateTime = time.Now()
		obj.LastTime = new(time.Time)
		*obj.LastTime = time.Now()
		obj.OwnerID = ownerID

		// get id
		id, err := cli.CC.InstCli.GetIncID(obj.TableName())
		if err != nil {
			blog.Error("failed to get id , error info is %s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		obj.ID = int(id)

		// save
		_, err = cli.CC.InstCli.Insert(obj.TableName(), obj)
		if nil == err {
			return http.StatusOK, []*metadata.ObjectDes{obj}, nil
		}
		blog.Error("failed to insert the object, error info is %s", err.Error())
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
	}, resp)
}

//删除Object
func (cli *objectAction) DeleteObject(req *restful.Request, resp *restful.Response) {

	blog.Info("delete object")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParameters := req.PathParameters()
		var appID int
		if err := cli.GetParams(cli.CC, &pathParameters, "id", &appID, resp); nil != err {
			blog.Error("failed to get params, error info is %s ", err.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsInvalid, "id")
		}

		condition := map[string]interface{}{"id": appID}
		// delete object from storage
		if 0 == appID {
			js, err := simplejson.NewFromReader(req.Request.Body)
			if err != nil {
				blog.Error("read http request body failed, error:%s", err.Error())
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
			}

			condition, err = js.Map()
			if nil != err {
				blog.Error("fail to unmarshal json, error information is %s", err.Error())
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
			}

		}
		util.SetModOwner(condition, ownerID)
		cnt, cntErr := cli.CC.InstCli.GetCntByCondition(metadata.ObjectDes{}.TableName(), condition)
		if nil != cntErr {
			blog.Error("failed to select object by condition(%+v), error is %d", cntErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		if 0 == cnt {
			// success
			return http.StatusOK, nil, nil
		}
		// execute delete command
		if delErr := cli.CC.InstCli.DelByCondition(metadata.ObjectDes{}.TableName(), condition); nil != delErr {
			blog.Error("fail to delete object by id , error information is %s", delErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		// success
		return http.StatusOK, nil, nil
	}, resp)
}

//编辑Object信息
func (cli *objectAction) UpdateObject(req *restful.Request, resp *restful.Response) {

	blog.Info("update object")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		js, err := simplejson.NewFromReader(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		pathParameters := req.PathParameters()
		var appID int
		if err = cli.GetParams(cli.CC, &pathParameters, "id", &appID, resp); nil != err {
			blog.Error("failed to get params, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsInvalid, "id")
		}

		// update object into storage
		js.Set(common.LastTimeField, util.GetCurrentTimeStr())

		// decode json string
		data, jsErr := js.Map()
		if nil != jsErr {
			blog.Error("unmarshal json failed, error information is %v", jsErr)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		condition := util.SetModOwner(map[string]interface{}{"id": appID}, ownerID)
		err = cli.CC.InstCli.UpdateByCondition(metadata.ObjectDes{}.TableName(), data, condition)
		if nil != err {
			blog.Error("fail update object by condition, error information is %s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}

		// success
		return http.StatusOK, nil, nil
	}, resp)

}

//查询所有主机信息
func (cli *objectAction) SelectObjects(req *restful.Request, resp *restful.Response) {

	blog.Info("select object")

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
		/*
			selector := &metadata.ObjectDes{Page: &metadata.BasePage{Limit: common.BKNoLimit}}
			if jserr := json.Unmarshal([]byte(value), selector); nil != jserr {
				blog.Error("unmarshal failed, error information is %s", jserr.Error())
				return http.StatusBadRequest, nil, defErr.Error(common.CC_Err_Comm_JSONUnmarshalFailed)
			}
		*/
		page := metadata.BasePage{Limit: common.BKNoLimit}
		if pageJS, ok := js.CheckGet("page"); ok {
			tmpMap, _ := pageJS.Map()
			page = metadata.ParsePage(tmpMap)
			js.Del("page")
		}
		results := make([]metadata.ObjectDes, 0)

		// select from storage
		//blog.Debug("selector:%+v", selector)
		selector, _ := js.Map()
		selector = util.SetQueryOwner(selector, ownerID)
		if selErr := cli.CC.InstCli.GetMutilByCondition(metadata.ObjectDes{}.TableName(), nil, selector, &results, page.Sort, page.Start, page.Limit); nil != selErr {
			blog.Error("select data failed, error information is %s", selErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}

		// translate language
		for index := range results {
			results[index].ObjectName = commondata.TranslateObjectName(defLang, &results[index])
		}

		// success
		return http.StatusOK, results, nil

	}, resp)
}
