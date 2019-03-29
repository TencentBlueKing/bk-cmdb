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

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

// CreateObject create a common object
func (cli *Service) CreateObject(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request body failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	obj := &meta.Object{}
	if jsErr := json.Unmarshal([]byte(value), obj); nil != jsErr {
		blog.Errorf("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}

	// save to the storage
	obj.CreateTime = new(time.Time)
	*obj.CreateTime = time.Now()
	obj.LastTime = new(time.Time)
	*obj.LastTime = time.Now()
	obj.OwnerID = ownerID

	// get id
	id, err := db.NextSequence(ctx, common.BKTableNameObjDes) // .GetIncID(common.BKTableNameObjDes)
	if err != nil {
		blog.Errorf("failed to get id , error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}
	obj.ID = int64(id)

	// save
	err = db.Table(common.BKTableNameObjDes).Insert(ctx, obj)
	if nil == err {
		resp.WriteEntity(meta.CreateObjectResult{BaseResp: meta.SuccessBaseResp, Data: meta.RspID{ID: int64(id)}})
		return
	}
	blog.Errorf("failed to insert the object, error info is %s", err.Error())

	resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})

}

// DeleteObject 删除Object
func (cli *Service) DeleteObject(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	pathParameters := req.PathParameters()
	appID, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Errorf("failed to get params, error info is %s ", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	condition := map[string]interface{}{"id": appID}
	// delete object from storage
	if 0 == appID {
		js, err := simplejson.NewFromReader(req.Request.Body)
		if err != nil {
			blog.Errorf("read http request body failed, error:%s", err.Error())
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
	cnt, cntErr := db.Table(common.BKTableNameObjDes).Find(condition).Count(ctx)
	if nil != cntErr {
		blog.Errorf("failed to select object by condition(%+v), error is %v", condition, cntErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, cntErr.Error())})
		return
	}
	if 0 == cnt {
		// success
		resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
		return
	}
	// execute delete command
	if delErr := db.Table(common.BKTableNameObjDes).Delete(ctx, condition); nil != delErr {
		blog.Errorf("fail to delete object by id , error information is %s", delErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, delErr.Error())})
		return
	}
	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})
}

func (cli *Service) UpdateObject(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	js, err := simplejson.NewFromReader(req.Request.Body)
	if err != nil {
		blog.Errorf("read request body failed, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
		return
	}

	pathParameters := req.PathParameters()
	appID, err := strconv.ParseInt(pathParameters["id"], 10, 64)
	if nil != err {
		blog.Errorf("failed to get params, error info is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	// update object into storage
	js.Set(common.LastTimeField, util.GetCurrentTimeStr())

	// decode json string
	data, jsErr := js.Map()
	if nil != jsErr {
		blog.Errorf("unmarshal json failed, error information is %v", jsErr)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, jsErr.Error())})
		return
	}
	condition := map[string]interface{}{"id": appID}
	condition = util.SetModOwner(condition, ownerID)
	err = db.Table(common.BKTableNameObjDes).Update(ctx, condition, data)
	if nil != err {
		blog.Errorf("fail update object by condition, error information is %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, jsErr.Error())})
		return
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp})

}

// SelectObjects 查询所有主机信息
func (cli *Service) SelectObjects(req *restful.Request, resp *restful.Response) {

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
		if nil != err {
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrCommHTTPReadBodyFailed, err.Error())})
			return
		}
	}

	page := meta.BasePage{Limit: common.BKNoLimit}
	if pageJS, ok := js.CheckGet("page"); ok {
		tmpMap, err := pageJS.Map()
		if nil != err {
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
			return
		}
		page = meta.ParsePage(tmpMap)
		js.Del("page")
	}
	results := make([]meta.Object, 0)

	// select from storage

	selector, err := js.Map()
	if nil != err {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, err.Error())})
		return
	}

	selector = util.SetQueryOwner(selector, ownerID)
	if selErr := db.Table(common.BKTableNameObjDes).Find(selector).Limit(uint64(page.Limit)).Start(uint64(page.Start)).Sort(page.Sort).All(ctx, &results); nil != selErr && !db.IsNotFoundError(selErr) {
		blog.Errorf("select data failed, error information is %s", selErr.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.New(common.CCErrObjectDBOpErrno, selErr.Error())})
		return
	}

	// translate language
	for index := range results {
		results[index].ObjectName = cli.TranslateObjectName(defLang, &results[index])
	}

	// success
	resp.WriteEntity(meta.Response{BaseResp: meta.SuccessBaseResp, Data: results})

}
