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
	"encoding/json"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
	"strconv"
)

var objcls = &objClsAction{}

// objClsAction
type objClsAction struct {
	base.BaseAction
	mgr manager.Manager
}

func init() {

	// register action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/object/classification", Params: nil, Handler: objcls.CreateClassification})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/object/classification/{owner_id}/objects", Params: nil, Handler: objcls.SelectClassificationWithObjects})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/object/classifications", Params: nil, Handler: objcls.SelectClassificationWithParams})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/object/classification/{id}", Params: nil, Handler: objcls.UpdateClassification})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/object/classification/{id}", Params: nil, Handler: objcls.DeleteClassification})

	// create cc object
	objcls.CreateAction()
	// set httpclient
	manager.SetManager(objcls)

}

// SetManager implement the manager's Hooker interface
func (cli *objClsAction) SetManager(mgr manager.Manager) error {
	cli.mgr = mgr
	return nil
}

func (cli *objClsAction) CreateClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("create object classification information")

	// get the language
	language := util.GetActionLanguage(req)
	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read body
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// deal data
		result, ctrErr := cli.mgr.CreateObjectClass(val, defErr)
		if nil == ctrErr {
			return http.StatusOK, map[string]int{"id": result}, nil
		}
		blog.Error("create object classification failed, error information is %s", ctrErr.Error())
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectClassificationCreateFailed)

	}, resp)
}

func (cli *objClsAction) SelectClassificationWithParams(req *restful.Request, resp *restful.Response) {

	blog.Info("select object classification information with params")
	// get the language
	language := util.GetActionLanguage(req)
	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read body
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// deal data
		result, ctrErr := cli.mgr.SelectObjectClass(val, defErr)
		if nil == ctrErr {
			return http.StatusOK, result, nil
		}
		blog.Error("search object classification failed, error information is %s", ctrErr.Error())
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectClassificationSelectFailed)

	}, resp)
}

func (cli *objClsAction) SelectClassificationWithObjects(req *restful.Request, resp *restful.Response) {

	blog.Info("select object classification with objects")

	// get the language
	language := util.GetActionLanguage(req)
	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read parameters
		ownerID := req.PathParameter("owner_id")

		// read body
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// deal data
		result, ctrErr := cli.mgr.SelectObjectClassWithObjects(ownerID, val, defErr)
		for idx, item := range result {
			result[idx].AsstObjects = make(map[string][]interface{})
			for _, subItem := range item.Objects {
				search := map[string]interface{}{}
				search[common.BKOwnerIDField] = subItem.OwnerID
				search[common.BKObjIDField] = subItem.ObjectID
				asstRst, asstErr := cli.mgr.SelectObjectAsst(search, defErr)
				if nil != asstErr {
					blog.Error("failed to search class object , error info is %s", asstErr.Error())
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectClassificationSelectFailed)
				}
				for _, asstItem := range asstRst {
					if asstItem.ObjectAttID == common.BKChildStr {
						continue
					}
					condition := map[string]interface{}{}
					condition[common.BKOwnerIDField] = asstItem.OwnerID
					condition[common.BKObjIDField] = asstItem.AsstObjID
					conditionStr, _ := json.Marshal(condition)
					objs, objErr := cli.mgr.SelectObject(conditionStr, defErr)
					if nil != objErr {
						blog.Error("failed to find the asst , error info is %s", objErr.Error())
						return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectClassificationSelectFailed)
					}

					for _, obj := range objs {
						result[idx].AsstObjects[subItem.ObjectID] = append(result[idx].AsstObjects[subItem.ObjectID], obj)
					}
				}
			}
			//cli.mgr.SelectObjectAsst
		}
		if nil == ctrErr {
			return http.StatusOK, result, nil
		}
		blog.Error("search object classification failed, error information is %s", ctrErr.Error())
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectClassificationSelectFailed)

	}, resp)

}

func (cli *objClsAction) UpdateClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("update objectatt information")

	// get the language
	language := util.GetActionLanguage(req)
	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// get params
		id, idErr := strconv.Atoi(req.PathParameter("id"))
		if nil != idErr {
			blog.Error("id(%s) is invalid, error info is %s", req.PathParameter("id"), idErr.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsNeedInt, "id")
		}

		// deal data
		ctrErr := cli.mgr.UpdateObjectClass(id, val, defErr)
		if nil == ctrErr {
			return http.StatusOK, nil, nil
		}

		blog.Error("update objectattr failed, error information is %s", ctrErr.Error())
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectClassificationUpdateFailed)

	}, resp)
}

func (cli *objClsAction) DeleteClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("delete object classification information")

	// get the language
	language := util.GetActionLanguage(req)
	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read body
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// get params
		id, idErr := strconv.Atoi(req.PathParameter("id"))
		if nil != idErr {
			blog.Error("id(%s) is invalid, error info is %s", req.PathParameter("id"), idErr.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsNeedInt, "id")
		}

		// search classifications
		condition := map[string]int{
			"id": id,
		}
		conditionStr, _ := json.Marshal(condition)
		rstItems, rstErr := cli.mgr.SelectObjectClass(conditionStr, defErr)
		if nil != rstErr {
			blog.Error("failed to search classification by id(%d), error info is %s", id, rstErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectClassificationDeleteFailed)
		}
		if 0 == len(rstItems) {
			blog.Warn("there is not  any classification, id(%d)", id)
			return http.StatusOK, nil, nil
		}

		// check the object with classification
		for _, rst := range rstItems {
			objCondition := map[string]interface{}{
				"bk_classification_id": rst.ClassificationID,
			}
			objConditionStr, _ := json.Marshal(objCondition)
			objItems, objErr := cli.mgr.SelectObject(objConditionStr, defErr)
			if nil != objErr {
				blog.Error("failed to search object with classificationid(%s), error info is %s", rst.ClassificationID, objErr.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectClassificationDeleteFailed)
			}

			if 0 != len(objItems) {
				blog.Info("there is some object assgined classification(%s)", rst.ClassificationID)
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrTopoObjectClassificationHasObject)
			}
		}

		// deal data
		ctrErr := cli.mgr.DeleteObjectClass(id, val, defErr)
		if nil == ctrErr {
			return http.StatusOK, nil, nil
		}
		blog.Error("update objectattr failed, error information is %s", ctrErr.Error())
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectClassificationDeleteFailed)
	}, resp)
}
