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
	"configcenter/src/common/util"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"
)

var objattgroup = &objectAttGroupAction{}

// objectAttGroupAction
type objectAttGroupAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/meta/objectatt/group/new", Params: nil, Handler: objattgroup.CreatePropertyGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/meta/objectatt/group/update", Params: nil, Handler: objattgroup.UpdatePropertyGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/meta/objectatt/group/groupid/{id}", Params: nil, Handler: objattgroup.DeletePropertyGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/meta/objectatt/group/property", Params: nil, Handler: objattgroup.UpdatePropertyGroupObjectAtt})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/meta/objectatt/group/owner/{owner_id}/object/{object_id}/propertyids/{property_id}/groupids/{group_id}", Params: nil, Handler: objattgroup.DeletePropertyGroupObjectAtt})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/objectatt/group/property/owner/{owner_id}/object/{object_id}", Params: nil, Handler: objattgroup.SelectPropertyGroupByObjectID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/objectatt/group/search", Params: nil, Handler: objattgroup.SelectGroup})

	// create action
	objattgroup.CreateAction()
}

// CreatePropertyGroup to create property group
func (cli *objectAttGroupAction) CreatePropertyGroup(req *restful.Request, resp *restful.Response) {

	blog.Info("create property group")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)

	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		defer req.Request.Body.Close()

		// read body data
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// parse the body data
		propertyGroup := &metadata.PropertyGroup{}
		jsErr := json.Unmarshal(val, propertyGroup)
		if nil != jsErr {
			blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(val), jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		//  save the data
		blog.Debug("store the property group: %+v ", propertyGroup)
		id, err := cli.CC.InstCli.GetIncID(propertyGroup.TableName())
		if err != nil {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectPropertyGroupInsertFailed)
		}
		propertyGroup.ID = int(id)
		propertyGroup.OwnerID = ownerID

		_, err = cli.CC.InstCli.Insert(propertyGroup.TableName(), propertyGroup)
		if nil == err {
			return http.StatusOK, []*metadata.PropertyGroup{propertyGroup}, nil
		}

		blog.Error("failed to insert the property group , error info is %s", err.Error())
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectPropertyGroupInsertFailed)

	}, resp)

}

// UpdatePropertyGroup to update property group
func (cli *objectAttGroupAction) UpdatePropertyGroup(req *restful.Request, resp *restful.Response) {

	blog.Info("update property group")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)

	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read body data
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		propertyGroup := &PropertyGroupCondition{}
		jsErr := json.Unmarshal(val, propertyGroup)
		if nil != jsErr {
			blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(val), jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		propertyGroup.Condition = util.SetModOwner(propertyGroup.Condition, ownerID)
		propertyGroup.Data = util.SetModOwner(propertyGroup.Data, ownerID)
		blog.Debug("property group:%+v", propertyGroup)
		if updateerr := cli.CC.InstCli.UpdateByCondition(common.BKTableNamePropertyGroup, propertyGroup.Data, propertyGroup.Condition); nil != updateerr {
			blog.Error("fail update object by condition, error:%v", updateerr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectPropertyGroupUpdateFailed)
		}

		return http.StatusOK, nil, nil
	}, resp)
}

// SelectGroup search groups
func (cli *objectAttGroupAction) SelectGroup(req *restful.Request, resp *restful.Response) {

	blog.Info("select property group")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// selector := &metadata.PropertyGroup{Page: &metadata.BasePage{Limit: common.BKNoLimit}}
		condition := map[string]interface{}{}
		if jsErr := json.Unmarshal([]byte(value), &condition); nil != jsErr {
			blog.Error("unmarshal failed, error information %s is %s", value, jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		page := metadata.ParsePage(condition["page"])
		delete(condition, "page")

		condition = util.SetQueryOwner(condition, ownerID)
		results := make([]metadata.PropertyGroup, 0)
		if selerr := cli.CC.InstCli.GetMutilByCondition(common.BKTableNamePropertyGroup, nil, condition, &results, page.Sort, page.Start, page.Limit); nil != selerr {
			blog.Error("find object by selector failed, error information is %s", selerr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrObjectPropertyGroupSelectFailed)
		}
		// translate language
		for index := range results {
			results[index].GroupName = commondata.TranslatePropertyGroupName(defLang, &results[index])
		}
		return http.StatusOK, results, nil

	}, resp)

}

// DeletePropertyGroup to update property group
func (cli *objectAttGroupAction) DeletePropertyGroup(req *restful.Request, resp *restful.Response) {
	blog.Info("delete property group")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		id, conErr := strconv.Atoi(req.PathParameter("id"))
		if nil != conErr {
			blog.Error("id(%s) should be int value, error info is %s", req.PathParameter("id"), conErr.Error())
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsNeedInt, "id")
		}

		condition := map[string]interface{}{"id": id}
		cnt, cntErr := cli.CC.InstCli.GetCntByCondition(common.BKTableNamePropertyGroup, condition)
		if nil != cntErr {
			blog.Error("failed to select object group by condition(%+v), error is %d", cntErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectPropertyGroupDeleteFailed)
		}
		if 0 == cnt {
			return http.StatusOK, nil, nil
		}
		condition = util.SetModOwner(condition, ownerID)
		if delErr := cli.CC.InstCli.DelByCondition(common.BKTableNamePropertyGroup, condition); nil != delErr {
			blog.Error("failed to delete property group  by condition, error:%v", delErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectPropertyGroupDeleteFailed)
		}

		return http.StatusOK, nil, nil
	}, resp)
}

// UpdatePropertyGroupObjectAtt to update property group object attribute
func (cli *objectAttGroupAction) UpdatePropertyGroupObjectAtt(req *restful.Request, resp *restful.Response) {
	blog.Info("update property group")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read body data
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		blog.Debug("property group is %s", string(val))

		// decode the data struct
		propertyGroupObjectAttArr := make([]PropertyGroupObjectAtt, 0)
		jsErr := json.Unmarshal(val, &propertyGroupObjectAttArr)
		if nil != jsErr {
			blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(val), jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		blog.Debug("property group object att, %+v", propertyGroupObjectAttArr)

		for _, objAtt := range propertyGroupObjectAttArr {

			// update the object attributes
			objectAttSelector := map[string]interface{}{
				common.BKOwnerIDField:    objAtt.Condition.OwnerID,
				common.BKObjIDField:      objAtt.Condition.ObjectID,
				common.BKPropertyIDField: objAtt.Condition.PropertyID,
			}

			objectAttValue := map[string]interface{}{
				"bk_property_index": objAtt.Data.PropertyIndex,
				"bk_property_group": objAtt.Data.PropertyGroupID,
			}

			objectAttSelector = util.SetModOwner(objectAttSelector, ownerID)
			objectAttValue = util.SetModOwner(objectAttValue, ownerID)
			// update the object attribute
			if updateerr := cli.CC.InstCli.UpdateByCondition(common.BKTableNameObjAttDes, objectAttValue, objectAttSelector); nil != updateerr {
				blog.Error("fail update object by condition, error:%v", updateerr.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectPropertyGroupUpdateFailed)
			}
		}
		return http.StatusOK, nil, nil
	}, resp)
}

// DeletePropertyGroupObjectAtt to delete property group object attribute
func (cli *objectAttGroupAction) DeletePropertyGroupObjectAtt(req *restful.Request, resp *restful.Response) {

	blog.Info("delete property group object attribute")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		// update the object attributes
		objectAttSelector := map[string]interface{}{
			common.BKOwnerIDField:       req.PathParameter("owner_id"),
			common.BKObjIDField:         req.PathParameter("object_id"),
			common.BKPropertyIDField:    req.PathParameter("property_id"),
			common.BKPropertyGroupField: req.PathParameter("group_id"),
		}

		objectAttValue := map[string]interface{}{
			"bk_property_index":         -1,
			common.BKPropertyGroupField: "default",
		}

		cnt, cntErr := cli.CC.InstCli.GetCntByCondition(common.BKTableNameObjAttDes, objectAttSelector)
		if nil != cntErr {
			blog.Error("failed to select objectatt group by condition(%+v), error is %d", cntErr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectPropertyGroupDeleteFailed)
		}
		if 0 == cnt {
			return http.StatusOK, nil, nil
		}
		objectAttSelector = util.SetModOwner(objectAttSelector, ownerID)
		blog.Debug("group property selector %+v, value %+v", objectAttSelector, objectAttValue)
		// update the object attribute
		if updateerr := cli.CC.InstCli.UpdateByCondition(common.BKTableNameObjAttDes, objectAttValue, objectAttSelector); nil != updateerr {
			blog.Error("fail update object by condition, error:%v", updateerr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrObjectPropertyGroupUpdateFailed)
		}

		return http.StatusOK, nil, nil
	}, resp)
}

// SelectPropertyGroupByObjectID to search
func (cli *objectAttGroupAction) SelectPropertyGroupByObjectID(req *restful.Request, resp *restful.Response) {

	blog.Info("select property group object attribute")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {

		groupSelector := map[string]interface{}{}
		if jserr := json.NewDecoder(req.Request.Body).Decode(&groupSelector); nil != jserr {
			blog.Error("unmarshal failed,  is %s", jserr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		// update the object attributes
		groupSelector[common.BKOwnerIDField] = req.PathParameter("owner_id")
		groupSelector[common.BKObjIDField] = req.PathParameter("object_id")

		page := metadata.ParsePage(groupSelector["page"])
		if page.Sort == "" {
			page.Sort = "bk_group_name"
		}
		delete(groupSelector, "page")

		groupSelector = util.SetQueryOwner(groupSelector, ownerID)
		blog.Debug("group property selector %+v", groupSelector)
		results := make([]metadata.PropertyGroup, 0)
		// select the object group
		if selerr := cli.CC.InstCli.GetMutilByCondition(common.BKTableNamePropertyGroup, nil, groupSelector, &results, page.Sort, page.Start, page.Limit); nil != selerr {
			blog.Error("select data failed, error information is %s", selerr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectPropertyGroupSelectFailed)
		}

		// translate language
		for index := range results {
			results[index].GroupName = commondata.TranslatePropertyGroupName(defLang, &results[index])
		}

		return http.StatusOK, results, nil
	}, resp)
}
