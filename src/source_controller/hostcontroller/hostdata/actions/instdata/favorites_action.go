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

package instdata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	. "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/common/commondata"

	"github.com/emicklei/go-restful"
	"github.com/rs/xid"
)

var (
	TABLENAME string = "cc_HostFavourite"
)

var hostFavouriteAction = &hostFavourite{}

type hostFavourite struct {
	base.BaseAction
}

//AddHostFavourite add host favorites
func (cli *hostFavourite) AddHostFavourite(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cc := api.NewAPIResource()
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("add host favourite failed, err: %v", err)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommHTTPReadBodyFailed).Error()})
		return
	}

	params := make(map[string]interface{})
	if err = json.Unmarshal([]byte(value), &params); nil != err {
		blog.Errorf("add host favourite failed, err: %v", err)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error()})
		return
	}

	queryParams := make(map[string]interface{})
	queryParams["user"] = req.PathParameter("user")
	queryParams["name"] = params["name"]

	rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, queryParams)
	if nil != err {
		blog.Error("query host favorites fail, err: %v, params:%v", err, queryParams)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteQueryFail).Error()})
		return
	}
	if 0 != rowCount {
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteCreateFail).Error()})
		return
	}
	//mogo 需要使用生产的id
	xidDevice := xid.New()
	params["id"] = xidDevice.String()
	params["count"] = 1
	params[common.CreateTimeField] = time.Now()
	params["user"] = req.PathParameter("user")
	_, err = cc.InstCli.Insert(TABLENAME, params)
	if err != nil {
		blog.Errorf("create host favorites failed, data:%v error:%v", params, err)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteCreateFail).Error()})
		return
	}

	resp.WriteAsJson(HostFavorite{
		BaseResp: BaseResp{true, http.StatusOK, ""},
		Data:     ID{ID: xidDevice.String()},
	})
	return
}

//UpdateHostFavouriteByID  update host fav
func (cli *hostFavourite) UpdateHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cc := api.NewAPIResource()

	id := req.PathParameter("id")
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("update host favourite failed, err: %v", err)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommHTTPReadBodyFailed).Error()})
		return
	}

	data := make(map[string]interface{})
	if err = json.Unmarshal([]byte(value), &data); nil != err {
		blog.Errorf("update host favourite failed, err: %v, msg:%s", err, string(value))
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error()})
		return
	}

	data[common.LastTimeField] = time.Now()

	params := make(map[string]interface{})
	params["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
	params["id"] = id
	rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, params)
	if nil != err {
		blog.Error("query host favorites fail, err: %v, params:%v", err, params)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteQueryFail).Error()})
		return
	}

	if 1 != rowCount {
		blog.Info("host favorites not permissions or not exists, params:%v", params)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteUpdateFail).Error()})
		return
	}

	//edit new not duplicate
	newName, ok := data["name"]
	if ok {
		dupParams := make(map[string]interface{})
		dupParams["name"] = newName
		dupParams[common.BKUser] = req.PathParameter("user")
		dupParams[common.BKFieldID] = common.KvMap{common.BKDBNE: id}
		rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, dupParams)
		if nil != err {
			blog.Error("query user api validate name duplicatie fail, err: %v, params:%v", err, dupParams)
			resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommDBSelectFailed).Error()})
			return
		}
		if 0 < rowCount {
			blog.Info("host user api  name duplicatie , params:%v", dupParams)
			resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommDuplicateItem).Error()})
			return
		}
	}

	err = cc.InstCli.UpdateByCondition(TABLENAME, data, params)
	if nil != err {
		blog.Error("update host favorites fail, err: %v, params:%v", err, params)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteUpdateFail).Error()})
		return
	}
	resp.WriteAsJson(BaseResp{Result: true, Code: http.StatusOK})
}

//DeleteHostFavouriteByID  delete host fav
func (cli *hostFavourite) DeleteHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cc := api.NewAPIResource()
	id := req.PathParameter("id")
	params := make(map[string]interface{})
	params["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
	params["id"] = id

	rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, params)
	if nil != err {
		blog.Error("query host favorites fail, err: %v, params:%v", err, params)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteQueryFail).Error()})
		return
	}
	if 1 != rowCount {
		blog.Info("host favorites not permissions or not exists, params:%v", params)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteDeleteFail).Error()})
		return
	}
	err = cc.InstCli.DelByCondition(TABLENAME, params)
	if nil != err {
		blog.Error("query host favourite fail, err: %v, params:%v", err, params)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteDeleteFail).Error()})
		return
	}
	resp.WriteAsJson(BaseResp{Result: true, Code: http.StatusOK})
}

//GetHostFavourites get host favorites
func (cli *hostFavourite) GetHostFavourites(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cc := api.NewAPIResource()
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("update host favourite failed, err: %v", err)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommHTTPReadBodyFailed).Error()})
		return
	}

	var dat commondata.ObjQueryInput
	err = json.Unmarshal([]byte(value), &dat)
	if err != nil {
		blog.Errorf("get host favourite failed, err: %v, msg:%s", err, string(value))
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error()})
		return
	}

	condition := make(map[string]interface{})
	if nil != dat.Condition {
		condition = dat.Condition.(map[string]interface{})
	}

	fieldArr := []string{"id", "info", "query_params", "name", "is_default", common.CreateTimeField, "count"}
	if "" != dat.Fields {
		fieldArr = strings.Split(dat.Fields, ",")
	}

	skip, limit, sort := dat.Start, dat.Limit, dat.Sort
	if 0 == limit {
		limit = 20
	}

	if "" == sort {
		sort = common.CreateTimeField
	}

	condition["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
	result := make([]interface{}, 0)
	count, err := cc.InstCli.GetCntByCondition(TABLENAME, condition)
	if err != nil {
		blog.Error("get host favorites failed,input:%v error:%v", string(value), err)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteQueryFail).Error()})
		return
	}

	err = cc.InstCli.GetMutilByCondition(TABLENAME, fieldArr, condition, &result, sort, skip, limit)
	if err != nil {
		blog.Error("get host favorites failed,input:%v error:%v", string(value), err)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteQueryFail).Error()})
		return
	}

	info := make(map[string]interface{})
	info["count"] = count
	info["info"] = result
	resp.WriteAsJson(Response{
		BaseResp: BaseResp{true, http.StatusOK, ""},
		Data:     info,
	})
}

//GetHostFavouriteByID get host favourite detail
func (cli *hostFavourite) GetHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cc := api.NewAPIResource()
	id := req.PathParameter("id")

	if "" == id || "0" == id {
		blog.Error("get host favourite id  emtpy")
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommParamsNeedSet).Error()})
		return
	}
	params := make(map[string]interface{})
	params["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
	params["id"] = id

	result := make(map[string]interface{})
	err := cc.InstCli.GetOneByCondition(TABLENAME, nil, params, &result)
	if err != nil && mgo_on_not_found_error != err.Error() {
		blog.Error("get host favourite failed,input: %v error: %v", id, err)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrHostFavouriteQueryFail).Error()})
		return
	}
	resp.WriteAsJson(Response{
		BaseResp: BaseResp{true, http.StatusOK, ""},
		Data:     result,
	})
}

func init() {
	hostFavouriteAction.CreateAction()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/favorites/{user}", Params: nil, Handler: hostFavouriteAction.AddHostFavourite})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/hosts/favorites/{user}/{id}", Params: nil, Handler: hostFavouriteAction.UpdateHostFavouriteByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/hosts/favorites/{user}/{id}", Params: nil, Handler: hostFavouriteAction.DeleteHostFavouriteByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/favorites/search/{user}", Params: nil, Handler: hostFavouriteAction.GetHostFavourites})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/favorites/search/{user}/{id}", Params: nil, Handler: hostFavouriteAction.GetHostFavouriteByID})

}
