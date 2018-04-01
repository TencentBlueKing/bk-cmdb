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
 
package object

import (
	"configcenter/src/common"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/topo_service/manager"
	"io/ioutil"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"
)

var objattgroup = &objectAttGroupAction{}

// objectAttGroupAction 模型属性分组的Action定义
type objectAttGroupAction struct {
	base.BaseAction
	mgr manager.Manager
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/objectatt/group/new", Params: nil, Handler: objattgroup.CreatePropertyGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/objectatt/group/update", Params: nil, Handler: objattgroup.UpdatePropertyGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/objectatt/group/groupid/{id}", Params: nil, Handler: objattgroup.DeletePropertyGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/objectatt/group/property", Params: nil, Handler: objattgroup.UpdatePropertyGroupObjectAtt})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/objectatt/group/owner/{owner_id}/object/{object_id}/propertyids/{property_id}/groupids/{group_id}", Params: nil, Handler: objattgroup.DeletePropertyGroupObjectAtt})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/objectatt/group/property/owner/{owner_id}/object/{object_id}", Params: nil, Handler: objattgroup.SelectPropertyGroupByObjectID})
	// create action
	objattgroup.CreateAction()
	// set httpclient
	manager.SetManager(objattgroup)
}

// SetManager implement the manager's Hooker interface
func (cli *objectAttGroupAction) SetManager(mgr manager.Manager) error {
	cli.mgr = mgr
	return nil
}

// CreatePropertyGroup to create property group
func (cli *objectAttGroupAction) CreatePropertyGroup(req *restful.Request, resp *restful.Response) {
	blog.Info("object attribute group, create")
	// get the language
	language := util.GetActionLanguage(req)

	// get the default error by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %v", err)
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// deal request
		id, err := cli.mgr.CreateObjectGroup(val, defErr)
		if nil != err {
			blog.Error("failed to create object group , error info is %s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectGroupCreateFailed)
		}

		return http.StatusOK, map[string]int{"id": id}, nil
	}, resp)
}

// UpdatePropertyGroup to update property group
func (cli *objectAttGroupAction) UpdatePropertyGroup(req *restful.Request, resp *restful.Response) {
	blog.Info("object attribute group, update")
	// get the language
	language := util.GetActionLanguage(req)

	// get the default error by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %v", err)
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// deal request
		err = cli.mgr.UpdateObjectGroup(val, defErr)
		if nil != err {
			blog.Error("failed to update object group, error info is %s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectGroupUpdateFailed)
		}
		return http.StatusOK, nil, nil

	}, resp)
}

// DeletePropertyGroup to update property group
func (cli *objectAttGroupAction) DeletePropertyGroup(req *restful.Request, resp *restful.Response) {
	blog.Info("object attribute group, delete")
	// get the language
	language := util.GetActionLanguage(req)

	// get the default error by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		id, conErr := strconv.Atoi(req.PathParameter("id"))
		if nil != conErr {
			blog.Error("failed to convert 'id(%s)' to int, error info is %s ", id, conErr.Error())
			return http.StatusBadRequest, "", defErr.Errorf(common.CCErrCommParamsNeedInt, "id")
		}

		// deal request
		err := cli.mgr.DeleteObjectGroup(id, defErr)
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectGroupDeleteFailed)
		}
		return http.StatusOK, nil, nil

	}, resp)
}

// UpdatePropertyGroupObjectAtt to update property group object attribute
func (cli *objectAttGroupAction) UpdatePropertyGroupObjectAtt(req *restful.Request, resp *restful.Response) {

	blog.Info("object attribute group property, update")
	// get the language
	language := util.GetActionLanguage(req)

	// get the default error by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %v", err)
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// deal request
		err = cli.mgr.UpdateObjectGroupProperty(val, defErr)
		if nil != err {
			blog.Error("failed to update object group property, error info is %s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectGroupUpdateFailed)
		}
		return http.StatusOK, nil, nil

	}, resp)
}

// DeletePropertyGroupObjectAtt to delete property group object attribute
func (cli *objectAttGroupAction) DeletePropertyGroupObjectAtt(req *restful.Request, resp *restful.Response) {

	blog.Info("object attribute group property, delete")
	// get the language
	language := util.GetActionLanguage(req)

	// get the default error by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {

		ownerID := req.PathParameter("owner_id")
		objectID := req.PathParameter("object_id")
		propertyID := req.PathParameter("property_id")
		groupID := req.PathParameter("group_id")

		// deal request
		err := cli.mgr.DeleteObjectGroupProperty(ownerID, objectID, propertyID, groupID, defErr)
		if nil != err {
			blog.Error("failed to delete object group property, error info is %s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectGroupDeleteFailed)
		}

		return http.StatusOK, nil, nil

	}, resp)
}

// SelectPropertyGroupByObjectID to select property group by objectid
func (cli *objectAttGroupAction) SelectPropertyGroupByObjectID(req *restful.Request, resp *restful.Response) {

	blog.Info("object attribute group property, select")
	// get the language
	language := util.GetActionLanguage(req)

	// get the default error by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		ownerID := req.PathParameter("owner_id")
		objectID := req.PathParameter("object_id")

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {

			blog.Error("read request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// deal request
		rst, err := cli.mgr.SelectPropertyGroupByObjectID(ownerID, objectID, value, defErr)
		if nil != err {
			blog.Error("failed to delete object group property, error info is %s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectGroupSelectFailed)
		}
		return http.StatusOK, rst, nil
	}, resp)
}
