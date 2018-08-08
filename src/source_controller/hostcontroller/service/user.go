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
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"github.com/emicklei/go-restful"
	"github.com/rs/xid"
)

func (s *Service) AddUserConfig(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	addQuery := new(meta.AddConfigQuery)
	if err := json.NewDecoder(req.Request.Body).Decode(&addQuery); err != nil {
		blog.Errorf("add user config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(addQuery.Name) == 0 {
		blog.Error("parameter Name is required")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "Name")})
		return
	}

	if 0 >= addQuery.AppID {
		blog.Error("add user config, parameter app id is required")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)})
		return
	}
	if len(addQuery.CreateUser) == 0 {
		blog.Error("add user config, parameter CreateUser is required")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "create_user")})
		return
	}

	queryParams := common.KvMap{"name": addQuery.Name, common.BKAppIDField: addQuery.AppID}
	queryParams = util.SetModOwner(queryParams, ownerID)
	rowCount, err := s.Instance.GetCntByCondition(UserQueryCollection, queryParams)
	if nil != err {
		blog.Errorf("add user config, query user api fail, error information is %s, params:%v", err.Error(), queryParams)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	if 0 != rowCount {
		blog.Errorf("add user config, [%s] user api is exist", addQuery.Name)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDuplicateItem)})
		return
	}

	id := xid.New().String()
	userQuery := meta.UserConfigMeta{
		AppID:      addQuery.AppID,
		Info:       addQuery.Info,
		Name:       addQuery.Name,
		ID:         id,
		CreateTime: time.Now().UTC(),
		CreateUser: addQuery.CreateUser,
		OwnerID:    ownerID,
		ModifyUser: addQuery.CreateUser,
		UpdateTime: time.Now().UTC(),
	}

	_, err = s.Instance.Insert(UserQueryCollection, userQuery)
	if err != nil {
		blog.Error("add user config, create user query failed, query:%+v err:%v", userQuery, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}

	resp.WriteEntity(meta.IDResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     meta.ID{ID: id},
	})
}

func (s *Service) UpdateUserConfig(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)
	id := req.PathParameter("id")

	appID, err := strconv.ParseInt(req.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user[%s] config failed, invalid appid[%s], err: %v", id, common.BKAppIDField, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	data := new(meta.UserConfigMeta)
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update user config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	params := common.KvMap{"id": id, common.BKAppIDField: appID}
	params = util.SetModOwner(params, ownerID)
	rowCount, err := s.Instance.GetCntByCondition(UserQueryCollection, params)
	if nil != err {
		blog.Error("query user api fail, error information is %s, params:%v", err.Error(), params)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	if 1 != rowCount {
		blog.Info("update user api config not permissions or not exists, params:%v", params)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommNotFound)})
		return
	}

	if len(data.Name) != 0 {
		dupParams := common.KvMap{"name": data.Name, common.BKAppIDField: appID, common.BKFieldID: common.KvMap{common.BKDBNE: id}}
		dupParams = util.SetModOwner(dupParams, ownerID)
		rowCount, getErr := s.Instance.GetCntByCondition(UserQueryCollection, dupParams)
		if nil != getErr {
			blog.Error("query user api validate name duplicate fail, error information is %s, params:%v", getErr.Error(), dupParams)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}
		if 0 < rowCount {
			blog.Info("host user api  name duplicate , params:%v", dupParams)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDuplicateItem)})
			return
		}
	}

	data.UpdateTime = time.Now().UTC()
	data.ModifyUser = util.GetUser(req.Request.Header)
	data.AppID = appID
	data.OwnerID = ownerID
	err = s.Instance.UpdateByCondition(UserQueryCollection, data, params)
	if nil != err {
		blog.Error("update user api fail, error information is %s, params:%v", err.Error(), params)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) DeleteUserConfig(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	id := req.PathParameter("id")
	appID, err := strconv.ParseInt(req.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user[%s] config failed, invalid appid[%s], err: %v", id, common.BKAppIDField, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	params := common.KvMap{"id": id, common.BKAppIDField: appID}
	params = util.SetModOwner(params, ownerID)
	rowCount, err := s.Instance.GetCntByCondition(UserQueryCollection, params)
	if nil != err {
		blog.Error("query user api fail, error information is %s, params:%v", err.Error(), params)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	if 1 != rowCount {
		blog.Info("host user api not permissions or not exists, params:%v", params)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommNotFound)})
		return
	}

	err = s.Instance.DelByCondition(UserQueryCollection, params)
	if nil != err {
		blog.Error("delete user api fail, error information is %s, params:%v", err.Error(), params)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBDeleteFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) GetUserConfig(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	dat := new(meta.ObjQueryInput)
	if err := json.NewDecoder(req.Request.Body).Decode(dat); err != nil {
		blog.Errorf("get user config failed with decode body, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	condition := make(map[string]interface{})
	if nil != dat.Condition {
		condition = dat.Condition.(map[string]interface{})
	}

	appID, err := util.GetInt64ByInterface(condition[common.BKAppIDField])
	if err != nil {
		blog.Errorf("get user config failed, invalid appid[%s], err: %v", common.BKAppIDField, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	start, limit, sort := dat.Start, dat.Limit, dat.Sort
	var fieldArr []string
	if "" != dat.Fields {
		fieldArr = strings.Split(",", dat.Fields)
	}

	if 0 == limit {
		limit = 20
	}
	if "" == sort {
		sort = common.CreateTimeField
	}

	condition[common.BKAppIDField] = appID
	condition = util.SetModOwner(condition, ownerID)
	count, err := s.Instance.GetCntByCondition(UserQueryCollection, condition)
	if err != nil {
		blog.Error("get user api information failed, err:%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	var result []interface{}
	err = s.Instance.GetMutilByCondition(UserQueryCollection, fieldArr, condition, &result, sort, start, limit)
	if err != nil {
		blog.Error("get user api information failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.GetUserConfigResult{
		BaseResp: meta.SuccessBaseResp,
		Data: meta.UserConfigResult{
			Count: count,
			Info:  result,
		},
	})
}

func (s *Service) UserConfigDetail(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)
	id := req.PathParameter("id")

	appID, err := strconv.ParseInt(req.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user[%s] config failed, invalid appid[%s], err: %v", id, common.BKAppIDField, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	var fieldArr []string
	params := common.KvMap{"id": id, common.BKAppIDField: appID}
	params = util.SetModOwner(params, ownerID)
	result := new(meta.UserConfigMeta)
	err = s.Instance.GetOneByCondition(UserQueryCollection, fieldArr, params, result)
	if err != nil && mgo_on_not_found_error != err.Error() {
		blog.Error("get user api information error,input:%v error:%v", id, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.GetUserConfigDetailResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *result,
	})

}

func (s *Service) AddUserCustom(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(language)
	ownerID := util.GetOwnerID(req.Request.Header)

	data := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("add user custom, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	ID := xid.New()
	data["id"] = ID.String()
	data["bk_user"] = req.PathParameter("bk_user")
	data = util.SetModOwner(data, ownerID)
	_, err := s.Instance.Insert(userCustomTableName, data)
	if nil != err {
		blog.Errorf("Create  user custom fail, err: %v, params:%v", err, data)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCreateUserCustom)})
		return

	}
	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) UpdateUserCustomByID(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(language)
	ownerID := util.GetOwnerID(req.Request.Header)

	conditons := make(map[string]interface{})
	conditons["id"] = req.PathParameter("id")
	conditons["bk_user"] = req.PathParameter("bk_user")
	data := make(map[string]interface{})

	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update user custom by id, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	conditons = util.SetModOwner(conditons, ownerID)
	err := s.Instance.UpdateByCondition(userCustomTableName, data, conditons)
	if nil != err {
		blog.Errorf("update  user custom failed, err: %v, data:%v", err, data)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) GetUserCustomByUser(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(language)

	conds, result := make(map[string]interface{}), make(map[string]interface{})
	conds["bk_user"] = req.PathParameter("bk_user")
	conds = util.SetModOwner(conds, ownerID)

	err := s.Instance.GetOneByCondition(userCustomTableName, nil, conds, &result)
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

func (s *Service) GetDefaultUserCustom(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(language)

	conds, result := make(map[string]interface{}), make(map[string]interface{})
	conds["is_default"] = 1
	conds = util.SetModOwner(conds, ownerID)

	err := s.Instance.GetOneByCondition(userCustomTableName, nil, conds, &result)
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
