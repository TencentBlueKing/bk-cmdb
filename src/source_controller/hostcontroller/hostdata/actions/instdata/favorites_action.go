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
    meta "configcenter/src/common/metadata"
    "configcenter/src/common/util"
    "configcenter/src/source_controller/common/commondata"
    "github.com/emicklei/go-restful"
    "github.com/rs/xid"
)

func init() {
    hostFavouriteAction.CreateAction()
    actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/favorites/{user}", Params: nil, Handler: hostFavouriteAction.AddHostFavourite})
    actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/hosts/favorites/{user}/{id}", Params: nil, Handler: hostFavouriteAction.UpdateHostFavouriteByID})
    actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/hosts/favorites/{user}/{id}", Params: nil, Handler: hostFavouriteAction.DeleteHostFavouriteByID})
    actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/favorites/search/{user}", Params: nil, Handler: hostFavouriteAction.GetHostFavourites})
    actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/favorites/search/{user}/{id}", Params: nil, Handler: hostFavouriteAction.GetHostFavouriteByID})
}

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
    params := make(map[string]interface{})
    if err := json.NewDecoder(req.Request.Body).Decode(params); err != nil {
        blog.Errorf("add host favorite , but decode body failed, err: %v", err)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
        return
    }

    queryParams := make(map[string]interface{})
    queryParams["user"] = req.PathParameter("user")
    queryParams["name"] = params["name"]

    rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, queryParams)
    if nil != err {
        blog.Error("query host favorites fail, err: %v, params:%v", err, queryParams)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
        return
    }
    if 0 != rowCount {
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteCreateFail)})
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
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteCreateFail)})
        return
    }

    resp.WriteEntity(meta.IDResult{
        BaseResp: meta.SuccessBaseResp,
        Data:     meta.ID{ID: xidDevice.String()},
    })
    return
}

//UpdateHostFavouriteByID  update host fav
func (cli *hostFavourite) UpdateHostFavouriteByID(req *restful.Request, resp *restful.Response) {
    language := util.GetActionLanguage(req)
    defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
    cc := api.NewAPIResource()

    ID := req.PathParameter("id")
    value, err := ioutil.ReadAll(req.Request.Body)
    if err != nil {
        blog.Errorf("update host favourite failed, err: %v", err)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPReadBodyFailed)})
        return
    }

    data := make(map[string]interface{})
    if err = json.Unmarshal([]byte(value), &data); nil != err {
        blog.Errorf("update host favourite failed, err: %v, msg:%s", err, string(value))
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
        return
    }

    data[common.LastTimeField] = time.Now()

    params := make(map[string]interface{})
    params["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
    params["id"] = ID
    rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, params)
    if nil != err {
        blog.Error("query host favorites fail, err: %v, params:%v", err, params)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
        return
    }

    if 1 != rowCount {
        blog.Info("host favorites not permissions or not exists, params:%v", params)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteUpdateFail)})
        return
    }

    //edit new not duplicate
    newName, ok := data["name"]
    if ok {
        dupParams := make(map[string]interface{})
        dupParams["name"] = newName
        dupParams[common.BKUser] = req.PathParameter("user")
        dupParams[common.BKFieldID] = common.KvMap{common.BKDBNE: ID}
        rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, dupParams)
        if nil != err {
            blog.Error("query user api validate name duplicate fail, err: %v, params:%v", err, dupParams)
            resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
            return
        }
        if 0 < rowCount {
            blog.Errorf("host user api  name duplicate , params:%v", dupParams)
            resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDuplicateItem)})
            return
        }
    }

    err = cc.InstCli.UpdateByCondition(TABLENAME, data, params)
    if nil != err {
        blog.Error("update host favorites fail, err: %v, params:%v", err, params)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteUpdateFail)})
        return
    }

    resp.WriteEntity(meta.NewSuccessResp(nil))
}

//DeleteHostFavouriteByID  delete host fav
func (cli *hostFavourite) DeleteHostFavouriteByID(req *restful.Request, resp *restful.Response) {
    language := util.GetActionLanguage(req)
    defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
    cc := api.NewAPIResource()
    ID := req.PathParameter("id")
    params := make(map[string]interface{})
    params["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
    params["id"] = ID

    rowCount, err := cc.InstCli.GetCntByCondition(TABLENAME, params)
    if nil != err {
        blog.Error("query host favorites fail, err: %v, params:%v", err, params)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
        return
    }
    if 1 != rowCount {
        blog.Info("host favorites not permissions or not exists, params:%v", params)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteDeleteFail)})
        return
    }
    err = cc.InstCli.DelByCondition(TABLENAME, params)
    if nil != err {
        blog.Error("query host favourite fail, err: %v, params:%v", err, params)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteDeleteFail)})
        return
    }
    resp.WriteEntity(meta.NewSuccessResp(nil))
}

//GetHostFavourites get host favorites
func (cli *hostFavourite) GetHostFavourites(req *restful.Request, resp *restful.Response) {
    language := util.GetActionLanguage(req)
    defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

    cc := api.NewAPIResource()
    value, err := ioutil.ReadAll(req.Request.Body)
    if err != nil {
        blog.Errorf("update host favourite failed, err: %v", err)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPReadBodyFailed)})
        return
    }

    var dat commondata.ObjQueryInput
    err = json.Unmarshal([]byte(value), &dat)
    if err != nil {
        blog.Errorf("get host favourite failed, err: %v, msg:%s", err, string(value))
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
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
        blog.Errorf("get host favorites failed,input:%v error:%v", string(value), err)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
        return
    }

    err = cc.InstCli.GetMutilByCondition(TABLENAME, fieldArr, condition, &result, sort, skip, limit)
    if err != nil {
        blog.Errorf("get host favorites failed,input:%v error:%v", string(value), err)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
        return
    }

    resp.WriteEntity(meta.GetHostFavoriteResult{
        BaseResp: meta.SuccessBaseResp,
        Data:     meta.FavoriteResult{Count: count, Info: result},
    })
}

//GetHostFavouriteByID get host favourite detail
func (cli *hostFavourite) GetHostFavouriteByID(req *restful.Request, resp *restful.Response) {
    language := util.GetActionLanguage(req)
    defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

    cc := api.NewAPIResource()
    ID := req.PathParameter("id")

    if "" == ID || "0" == ID {
        blog.Errorf("get host favourite, but id is emtpy")
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsNeedSet)})
        return
    }
    params := make(map[string]interface{})
    params["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
    params["id"] = ID

    result := make(map[string]interface{})
    err := cc.InstCli.GetOneByCondition(TABLENAME, nil, params, &result)
    if err != nil && mgo_on_not_found_error != err.Error() {
        blog.Errorf("get host favourite failed,input: %v error: %v", ID, err)
        resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
        return
    }
    resp.WriteEntity(meta.GetHostFavoriteWithIDResult{
        BaseResp: meta.SuccessBaseResp,
        Data:     result,
    })
}