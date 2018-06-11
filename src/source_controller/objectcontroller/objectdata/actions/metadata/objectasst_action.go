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
	"encoding/json"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
)

var objasst = &objectAssociationAction{}

// ObjectAction
type objectAssociationAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/objectassts", Params: nil, Handler: objasst.SelectObjectAssociations})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/meta/objectasst/{id}", Params: nil, Handler: objasst.DeleteObjectAssociation})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/meta/objectasst", Params: nil, Handler: objasst.CreateObjectAssociation})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/meta/objectasst/{id}", Params: nil, Handler: objasst.UpdateObjectAssociation})

	// set cc api resource
	objasst.CC = api.NewAPIResource()
}

// CreateObjectAssociation create object association map
func (cli *objectAssociationAction) CreateObjectAssociation(req *restful.Request, resp *restful.Response) {

	blog.Info("create obj-association")
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		obj := &metadata.ObjectAsst{}
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
			blog.Error("create objectasst failed, error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}

		return http.StatusOK, []*metadata.ObjectAsst{obj}, nil
	}, resp)
}

// DeleteObjectAssociation delete object association map
func (cli *objectAssociationAction) DeleteObjectAssociation(req *restful.Request, resp *restful.Response) {

	blog.Info("delete obj-association")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParameters := req.PathParameters()
		var id int
		if err := cli.GetParams(cli.CC, &pathParameters, "id", &id, resp); nil != err {
			blog.Error("failed to get params, error info is %s ", err.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsInvalid, "id")
		}

		// delete object from storage
		condition := map[string]interface{}{"id": id}
		if 0 == id {
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

		cnt, cntErr := cli.CC.InstCli.GetCntByCondition(metadata.ObjectAsst{}.TableName(), condition)
		if nil != cntErr {
			blog.Error("failed to select objectasst by condition(%+v), error is %d", cntErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		if 0 == cnt {
			// success
			return http.StatusOK, nil, nil
		}
		// execute delete command
		delErr := cli.CC.InstCli.DelByCondition(metadata.ObjectAsst{}.TableName(), condition)
		if nil != delErr {
			blog.Error("fail to delete object by id , error information is %s", delErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		// success
		return http.StatusOK, nil, nil
	}, resp)
}

// UpdateObjectAssociation update object association map
func (cli *objectAssociationAction) UpdateObjectAssociation(req *restful.Request, resp *restful.Response) {

	blog.Info("update object association")

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
		var id int
		if err := cli.GetParams(cli.CC, &pathParameters, "id", &id, resp); nil != err {
			blog.Error("failed to get params, error info is %s ", err.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsInvalid, "id")
		}

		// decode json string
		data, jsErr := js.Map()
		if nil != jsErr {
			blog.Error("unmarshal json failed, error information is %s", jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		blog.Debug("update:%+v", data)
		// update object into storage
		condition := util.SetModOwner(map[string]interface{}{"id": id}, ownerID)
		if updateErr := cli.CC.InstCli.UpdateByCondition(metadata.ObjectAsst{}.TableName(), data, condition); nil != updateErr {
			blog.Error("fail update object by condition, error information is %s", updateErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}

		// success
		return http.StatusOK, nil, nil

	}, resp)

}

// SelectObjectAssociations search all object association map
func (cli *objectAssociationAction) SelectObjectAssociations(req *restful.Request, resp *restful.Response) {

	blog.Info("search object association")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

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

		results := make([]metadata.ObjectAsst, 0)
		selector, _ := js.Map()
		selector = util.SetModOwner(selector, ownerID)
		// select from storage
		if selErr := cli.CC.InstCli.GetMutilByCondition(metadata.ObjectAsst{}.TableName(), nil, selector, &results, page.Sort, page.Start, page.Limit); nil != selErr {
			blog.Error("select data failed, error information is %s", selErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}

		// success
		return http.StatusOK, results, nil
	}, resp)
}
