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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"github.com/emicklei/go-restful"
	"github.com/rs/xid"
)

func init() {
	userCostomAction.CreateAction()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/usercustom/{bk_user}", Params: nil, Handler: userCostomAction.AddUserCustom})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/usercustom/{bk_user}/{id}", Params: nil, Handler: userCostomAction.UpdateUserCustomByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/usercustom/user/search/{bk_user}", Params: nil, Handler: userCostomAction.GetUserCustomByUser})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/usercustom/default/search/{bk_user}", Params: nil, Handler: userCostomAction.GetDefaultUserCustom})
}

var (
	userCustomTableName    string = "cc_UserCustom"
	mgo_on_not_found_error string = "not found"
)

var userCostomAction = &userCustomAction{}

type userCustomAction struct {
	base.BaseAction
}

//保存用户自定义的配置
//表字段User用户名。OwnerID 供应商,IsDefault是否默认配置，HostQueryColumn,HostDisplayColumn主机的自定义配置，

//AddUserCustom add user custom config
func (cli *userCustomAction) AddUserCustom(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	data := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("add user custom, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	ID := xid.New()
	data["id"] = ID.String()
	data["bk_user"] = req.PathParameter("bk_user")
	cc := api.NewAPIResource()
	_, err := cc.InstCli.Insert(userCustomTableName, data)
	if nil != err {
		blog.Errorf("Create  user custom fail, err: %v, params:%v", err, data)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCreateUserCustom)})
		return

	}
	resp.WriteEntity(meta.NewSuccessResp(nil))
}

//UpdateUserCustomByID update user custom
func (cli *userCustomAction) UpdateUserCustomByID(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	conditons := make(map[string]interface{})
	conditons["id"] = req.PathParameter("id")
	conditons["bk_user"] = req.PathParameter("bk_user")
	data := make(map[string]interface{})

	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update user custom by id, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	cc := api.NewAPIResource()
	err := cc.InstCli.UpdateByCondition(userCustomTableName, data, conditons)
	if nil != err {
		blog.Errorf("update  user custom failed, err: %v, data:%v", err, data)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

//GetUserCustomByUser get user custom
func (cli *userCustomAction) GetUserCustomByUser(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	conds, result := make(map[string]interface{}), make(map[string]interface{})
	conds["bk_user"] = req.PathParameter("bk_user")

	cc := api.NewAPIResource()
	err := cc.InstCli.GetOneByCondition(userCustomTableName, nil, conds, &result)
	if nil != err && mgo_on_not_found_error != err.Error() {
		blog.Error("add  user custom failed, err: %v, params:%v", err, conds)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.GetUserCustomResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}

//GetDefaultUserCustom get default user custom
func (cli *userCustomAction) GetDefaultUserCustom(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cc := api.NewAPIResource()

	conds, result := make(map[string]interface{}), make(map[string]interface{})
	conds["is_default"] = 1

	err := cc.InstCli.GetOneByCondition(userCustomTableName, nil, conds, &result)
	if nil != err {
		blog.Error("get default user custom fail, err: %v, params:%v", err, conds)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.GetUserCustomResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}
