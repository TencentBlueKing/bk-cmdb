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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	restful "github.com/emicklei/go-restful"
)

// CreatePropertyGroup to create property group
func (cli *Service) CreatePropertyGroup(req *restful.Request, resp *restful.Response) {

	blog.Info("create property group")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// execute

	// read body data
	val, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	// parse the body data
	propertyGroup := &meta.Group{}
	jsErr := json.Unmarshal(val, propertyGroup)
	if nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(val), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	//  save the data
	id, err := db.NextSequence(ctx, common.BKTableNamePropertyGroup)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectPropertyGroupInsertFailed, err.Error())})
		return
	}

	propertyGroup.ID = int64(id)
	propertyGroup.OwnerID = ownerID
	err = db.Table(common.BKTableNamePropertyGroup).Insert(ctx, propertyGroup)
	if nil == err {
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: propertyGroup})
		return
	}

	blog.Errorf("failed to insert the property group , error info is %s", err.Error())
	resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectPropertyGroupInsertFailed, err.Error())})

}

// UpdatePropertyGroup to update property group
func (cli *Service) UpdatePropertyGroup(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// read body data
	val, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}
	propertyGroup := &meta.PropertyGroupCondition{}
	jsErr := json.Unmarshal(val, propertyGroup)
	if nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(val), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	propertyGroup.Condition = util.SetModOwner(propertyGroup.Condition, ownerID)
	if updateErr := db.Table(common.BKTableNamePropertyGroup).Update(ctx, propertyGroup.Condition, propertyGroup.Data); nil != updateErr {
		blog.Errorf("fail update object by condition, error:%v", updateErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectPropertyGroupUpdateFailed, updateErr.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})

}

// SelectGroup search groups
func (cli *Service) SelectGroup(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	defLang := cli.Core.Language.CreateDefaultCCLanguageIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// execute
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read request body failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	// selector := &metadata.PropertyGroup{Page: &metadata.BasePage{Limit: common.BKNoLimit}}
	condition := map[string]interface{}{}
	if jsErr := json.Unmarshal([]byte(value), &condition); nil != jsErr {
		blog.Errorf("unmarshal failed, error information %s is %s", value, jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	page := meta.ParsePage(condition["page"])
	delete(condition, "page")
	condition = util.SetQueryOwner(condition, ownerID)

	results := make([]meta.Group, 0)
	if selErr := db.Table(common.BKTableNamePropertyGroup).Find(condition).Limit(uint64(page.Limit)).Start(uint64(page.Start)).Sort(page.Sort).All(ctx, &results); nil != selErr && db.IsNotFoundError(selErr) {
		blog.Errorf("find object by selector failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectPropertyGroupSelectFailed, selErr.Error())})
		return
	}

	// translate language
	for index := range results {
		results[index].GroupName = cli.TranslatePropertyGroupName(defLang, &results[index])
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})

}

// DeletePropertyGroup to update property group
func (cli *Service) DeletePropertyGroup(req *restful.Request, resp *restful.Response) {
	blog.Info("delete property group")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// execute

	id, conErr := strconv.Atoi(req.PathParameter("id"))
	if nil != conErr {
		blog.Errorf("id(%s) should be int value, error info is %s", req.PathParameter("id"), conErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsNeedInt, conErr.Error())})
		return
	}

	condition := map[string]interface{}{"id": id}
	condition = util.SetModOwner(condition, ownerID)
	cnt, cntErr := db.Table(common.BKTableNamePropertyGroup).Find(condition).Count(ctx)
	if nil != cntErr {
		blog.Errorf("failed to select object group by condition(%+v), error is %d", condition, cntErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectPropertyGroupDeleteFailed, cntErr.Error())})
		return
	}
	if 0 == cnt {
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
		return
	}
	if delErr := db.Table(common.BKTableNamePropertyGroup).Delete(ctx, condition); nil != delErr {
		blog.Errorf("failed to delete property group  by condition, error:%v", delErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectPropertyGroupDeleteFailed, delErr.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})

}

// UpdatePropertyGroupObjectAtt to update property group object attribute
func (cli *Service) UpdatePropertyGroupObjectAtt(req *restful.Request, resp *restful.Response) {
	blog.Info("update property group")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	// read body data
	val, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	// decode the data struct
	propertyGroupObjectAttArr := make([]meta.PropertyGroupObjectAtt, 0)
	jsErr := json.Unmarshal(val, &propertyGroupObjectAttArr)
	if nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(val), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return
	}

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
		// update the object attribute
		if updateErr := db.Table(common.BKTableNameObjAttDes).Update(ctx, objectAttSelector, objectAttValue); nil != updateErr {
			blog.Errorf("fail update object by condition, error:%v", updateErr.Error())
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectPropertyGroupUpdateFailed, updateErr.Error())})
			return
		}
	}
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

// DeletePropertyGroupObjectAtt to delete property group object attribute
func (cli *Service) DeletePropertyGroupObjectAtt(req *restful.Request, resp *restful.Response) {

	blog.Info("delete property group object attribute")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

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
	objectAttSelector = util.SetModOwner(objectAttSelector, ownerID)

	cnt, cntErr := db.Table(common.BKTableNameObjAttDes).Find(objectAttSelector).Count(ctx)
	if nil != cntErr {
		blog.Errorf("failed to select objectatt group by condition(%+v), error is %d", objectAttSelector, cntErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectPropertyGroupDeleteFailed, cntErr.Error())})
		return
	}
	if 0 == cnt {
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
		return
	}

	// update the object attribute
	if updateErr := db.Table(common.BKTableNameObjAttDes).Update(ctx, objectAttSelector, objectAttValue); nil != updateErr {
		blog.Errorf("fail update object by condition, error:%v", updateErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectPropertyGroupUpdateFailed, updateErr.Error())})
		return
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

// SelectPropertyGroupByObjectID to search
func (cli *Service) SelectPropertyGroupByObjectID(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	defLang := cli.Core.Language.CreateDefaultCCLanguageIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	groupSelector := map[string]interface{}{}
	if jsErr := json.NewDecoder(req.Request.Body).Decode(&groupSelector); nil != jsErr {
		blog.Errorf("unmarshal failed,  is %s", jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	// update the object attributes
	groupSelector[common.BKOwnerIDField] = req.PathParameter("owner_id")
	groupSelector[common.BKObjIDField] = req.PathParameter("object_id")

	page := meta.ParsePage(groupSelector["page"])
	if page.Sort == "" {
		page.Sort = "bk_group_name"
	}
	delete(groupSelector, "page")
	groupSelector = util.SetQueryOwner(groupSelector, ownerID)

	blog.V(3).Infof("group property selector %+v", groupSelector)
	results := make([]meta.Group, 0)
	// select the object group
	if selErr := db.Table(common.BKTableNamePropertyGroup).Find(groupSelector).Limit(uint64(page.Limit)).Start(uint64(page.Start)).Sort(page.Sort).All(ctx, &results); nil != selErr {
		blog.Errorf("select data failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectPropertyGroupSelectFailed, selErr.Error())})
		return
	}

	// translate language
	for index := range results {
		results[index].GroupName = cli.TranslatePropertyGroupName(defLang, &results[index])
	}

	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})
}
