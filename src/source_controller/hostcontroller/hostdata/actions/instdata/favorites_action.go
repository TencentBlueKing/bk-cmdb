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
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/rs/xid"
)

var (
	// TABLENAME cc_HostFavourite table name
	TABLENAME = "cc_HostFavourite"
)

var hostFavouriteAction = &hostFavourite{}

type hostFavourite struct {
	base.BaseAction
}

//AddHostFavourite add host favorites
func (cli *hostFavourite) AddHostFavourite(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := make(map[string]interface{}) //favouriteParms{}
		if err = json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %s, msg:%s", err.Error(), string(value))
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		queryParams := make(map[string]interface{})
		queryParams["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
		queryParams["name"] = params["name"]
		queryParams = util.SetModOwner(queryParams, ownerID)

		rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, queryParams)
		if nil != err {
			blog.Error("query host favorites fail, error information is %s, params:%v", err.Error(), queryParams)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteQueryFail)
		}
		if 0 != rowCount {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteCreateFail)
		}
		//mogo 需要使用生产的id
		xidDevice := xid.New()
		params["id"] = xidDevice.String()
		params["count"] = 1
		params[common.CreateTimeField] = time.Now()
		params["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
		params = util.SetModOwner(params, ownerID)
		_, err = cc.InstCli.Insert(TABLENAME, params)

		if err != nil {
			blog.Error("create host favorites type:data:%v error:%v", params, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteCreateFail)
		}

		info := make(map[string]interface{})
		info["id"] = xidDevice.String()

		return http.StatusOK, info, nil
	}, resp)

}

//UpdateHostFavouriteByID  update host fav
func (cli *hostFavourite) UpdateHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()

		ID := req.PathParameter("id")

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		data := make(map[string]interface{})
		if err = json.Unmarshal([]byte(value), &data); nil != err {
			blog.Error("fail to unmarshal json, error information is %s, msg:%s", err.Error(), string(value))
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		data[common.LastTimeField] = time.Now()

		params := make(map[string]interface{})
		params["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
		params["id"] = ID
		params = util.SetModOwner(params, ownerID)
		rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, params)
		if nil != err {
			blog.Error("query host favorites fail, error information is %s, params:%v", err.Error(), params)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteQueryFail)
		}
		if 1 != rowCount {
			blog.Info("host favorites not permissions or not exists, params:%v", params)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteUpdateFail)
		}

		//edit new not duplicate
		newName, ok := data["name"]
		if ok {
			dupParams := make(map[string]interface{})
			dupParams["name"] = newName
			dupParams[common.BKUser] = req.PathParameter("user")
			dupParams[common.BKFieldID] = common.KvMap{common.BKDBNE: ID}
			dupParams = util.SetModOwner(dupParams, ownerID)
			rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, dupParams)
			if nil != err {
				blog.Error("query user api validate name duplicatie fail, error information is %s, params:%v", err.Error(), dupParams)
				return http.StatusBadGateway, nil, defErr.Error(common.CCErrCommDBSelectFailed)
			}
			if 0 < rowCount {
				blog.Info("host user api  name duplicatie , params:%v", dupParams)
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommDuplicateItem)
			}
		}

		err = cc.InstCli.UpdateByCondition(TABLENAME, data, params)
		if nil != err {
			blog.Error("updata host favorites fail, error information is %s, params:%v", err.Error(), params)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteUpdateFail)
		}
		return http.StatusOK, nil, nil
	}, resp)

}

//DeleteHostFavouriteByID  delete host fav
func (cli *hostFavourite) DeleteHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		cc := api.NewAPIResource()
		ID := req.PathParameter("id")
		params := make(map[string]interface{})
		params["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
		params["id"] = ID
		params = util.SetModOwner(params, ownerID)

		rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, params)
		if nil != err {
			blog.Error("query host favorites fail, error information is %s, params:%v", err.Error(), params)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteQueryFail)
		}
		if 1 != rowCount {
			blog.Info("host favorites not permissions or not exists, params:%v", params)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteDeleteFail)
		}
		err = cc.InstCli.DelByCondition(TABLENAME, params)
		if nil != err {
			blog.Error("query host favourite fail, error information is %s, params:%v", err.Error(), params)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteDeleteFail)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

//GetHostFavourites get host favorites
func (cli *hostFavourite) GetHostFavourites(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()

		var dat commondata.ObjQueryInput
		value, err := ioutil.ReadAll(req.Request.Body)

		//if no params use default
		if nil == err && nil == value {
			value = []byte("{}")
		}
		err = json.Unmarshal([]byte(value), &dat)
		if err != nil {
			blog.Error("fail to unmarshal json, error information is,input:%v error:%v", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		fields := dat.Fields

		condition := make(map[string]interface{})
		if nil != dat.Condition {
			condition = dat.Condition.(map[string]interface{})
		}

		skip := dat.Start
		limit := dat.Limit
		sort := dat.Sort

		fieldArr := []string{"id", "info", "query_params", "name", "is_default", common.CreateTimeField, "count"}
		if "" != fields {
			fieldArr = strings.Split(fields, ",")
		}

		if 0 == limit {
			limit = 20
		}
		if "" == sort {
			sort = common.CreateTimeField
		}

		condition["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
		condition = util.SetModOwner(condition, ownerID)
		result := make([]interface{}, 0)
		count, err := cc.InstCli.GetCntByCondition(TABLENAME, condition)
		if err != nil {
			blog.Error("get host favorites infomation error,input:%v error:%v", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteQueryFail)
		}
		err = cc.InstCli.GetMutilByCondition(TABLENAME, fieldArr, condition, &result, sort, skip, limit)
		if err != nil {
			blog.Error("get host favorites infomation error,input:%v error:%v", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteQueryFail)
		}
		info := make(map[string]interface{})
		info["count"] = count
		info["info"] = result
		return http.StatusOK, info, nil
	}, resp)
	return
}

//GetHostFavouriteByID get host favourite detail
func (cli *hostFavourite) GetHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		cc := api.NewAPIResource()
		ID := req.PathParameter("id")

		if "" == ID || "0" == ID {
			blog.Error("get host favourite id  emtpy")
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, "id")
		}
		params := make(map[string]interface{})
		params["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
		params["id"] = ID
		params = util.SetModOwner(params, ownerID)
		result := make(map[string]interface{})
		err := cc.InstCli.GetOneByCondition(TABLENAME, nil, params, &result)
		if err != nil && mgo_on_not_found_error != err.Error() {
			blog.Error("get host favourite infomation error,input:%v error:%v", ID, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavouriteQueryFail)
		}

		return http.StatusOK, result, nil
	}, resp)

}

func init() {
	hostFavouriteAction.CreateAction()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/favorites/{user}", Params: nil, Handler: hostFavouriteAction.AddHostFavourite})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/hosts/favorites/{user}/{id}", Params: nil, Handler: hostFavouriteAction.UpdateHostFavouriteByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/hosts/favorites/{user}/{id}", Params: nil, Handler: hostFavouriteAction.DeleteHostFavouriteByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/favorites/search/{user}", Params: nil, Handler: hostFavouriteAction.GetHostFavourites})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/favorites/search/{user}/{id}", Params: nil, Handler: hostFavouriteAction.GetHostFavouriteByID})

}
