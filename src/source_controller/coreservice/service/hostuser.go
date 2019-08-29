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
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"

	"github.com/rs/xid"
)

func (s *coreService) AddUserConfig(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	addQuery := new(meta.AddConfigQuery)
	if err := data.MarshalJSONInto(addQuery); err != nil {
		blog.Errorf("add user config failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	if len(addQuery.Name) == 0 {
		blog.Error("parameter Name is required, rid: %s", params.ReqID)
		return nil, params.Error.CCErrorf(common.CCErrCommParamsNeedSet, "Name")
	}

	if 0 >= addQuery.AppID {
		blog.Error("add user config, parameter app id is required, rid: %s", params.ReqID)
		return nil, params.Error.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)
	}
	if len(addQuery.CreateUser) == 0 {
		blog.Error("add user config, parameter CreateUser is required, rid: %s", params.ReqID)
		return nil, params.Error.CCErrorf(common.CCErrCommParamsNeedSet, "create_user")
	}

	filter := common.KvMap{
		"name":              addQuery.Name,
		common.BKAppIDField: addQuery.AppID,
	}
	filter = util.SetModOwner(filter, params.SupplierAccount)
	rowCount, err := s.db.Table(common.BKTableNameUserAPI).Find(filter).Count(params.Context)
	if nil != err {
		blog.Errorf("add user config, query user api fail, error information is %s, params:%v, rid: %s", err.Error(), queryParams, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	if 0 != rowCount {
		blog.Errorf("add user config, [%s] user api is exist, rid: %s", addQuery.Name, params.ReqID)
		return nil, params.Error.CCErrorf(common.CCErrCommDuplicateItem, "")
	}

	id := xid.New().String()
	userQuery := meta.UserConfigMeta{
		AppID:      addQuery.AppID,
		Info:       addQuery.Info,
		Name:       addQuery.Name,
		ID:         id,
		CreateTime: time.Now().UTC(),
		CreateUser: addQuery.CreateUser,
		OwnerID:    params.SupplierAccount,
		ModifyUser: addQuery.CreateUser,
		UpdateTime: time.Now().UTC(),
	}

	err = s.db.Table(common.BKTableNameUserAPI).Insert(params.Context, userQuery)
	if err != nil {
		blog.Errorf("add user config, create user query failed, query:%+v err:%v, rid: %s", userQuery, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBInsertFailed)
	}

	return meta.ID{ID: id}, nil
}

func (s *coreService) UpdateUserConfig(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	id := pathParams("id")
	appID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user[%s] config failed, invalid appid[%s], err: %v, rid: %s", id, common.BKAppIDField, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommParamsIsInvalid)
	}

	dat := new(meta.UserConfigMeta)
	if err := data.MarshalJSONInto(dat); err != nil {
		blog.Errorf("update user config failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	filter := common.KvMap{
		"id":                id,
		common.BKAppIDField: appID,
	}
	filter = util.SetModOwner(filter, params.SupplierAccount)
	rowCount, err := s.db.Table(common.BKTableNameUserAPI).Find(filter).Count(params.Context)
	if nil != err {
		blog.Errorf("query user api fail, error information is %s, params:%v, rid: %s", err.Error(), params, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	if 1 != rowCount {
		blog.V(5).Infof("update user api config not permissions or not exists, params:%v, rid: %s", params, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommNotFound)
	}

	if len(dat.Name) != 0 {
		dupParams := common.KvMap{
			"name":              dat.Name,
			common.BKAppIDField: appID,
			common.BKFieldID:    common.KvMap{common.BKDBNE: id},
		}
		dupParams = util.SetModOwner(dupParams, params.SupplierAccount)
		rowCount, getErr := s.db.Table(common.BKTableNameUserAPI).Find(dupParams).Count(params.Context)
		if nil != getErr {
			blog.Errorf("query user api validate name duplicate fail, error information is %s, params:%v, rid: %s", getErr.Error(), dupParams, params.ReqID)
			return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
		}
		if 0 < rowCount {
			blog.V(5).Infof("host user api  name duplicate , params:%v, rid: %s", dupParams, params.ReqID)
			return nil, params.Error.CCErrorf(common.CCErrCommDuplicateItem, "")
		}
	}

	dat.UpdateTime = time.Now().UTC()
	dat.ModifyUser = util.GetUser(params.Header)
	dat.AppID = appID
	dat.OwnerID = params.SupplierAccount
	err = s.db.Table(common.BKTableNameUserAPI).Update(params.Context, filter, dat)
	if nil != err {
		blog.Errorf("update user api fail, error information is %s, params:%v, rid: %s", err.Error(), params, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBUpdateFailed)
	}

	return nil, nil
}

func (s *coreService) DeleteUserConfig(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	id := pathParams("id")
	appID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user[%s] config failed, invalid appid[%s], err: %v, rid: %s", id, common.BKAppIDField, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommParamsIsInvalid)
	}

	filter := common.KvMap{"id": id, common.BKAppIDField: appID}
	filter = util.SetModOwner(filter, params.SupplierAccount)
	rowCount, err := s.db.Table(common.BKTableNameUserAPI).Find(filter).Count(params.Context)
	if nil != err {
		blog.Errorf("query user api fail, error information is %s, params:%v, rid: %s", err.Error(), filter, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	if 1 != rowCount {
		blog.V(5).Infof("host user api not permissions or not exists, params:%v, rid: %s", filter, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommNotFound)
	}

	err = s.db.Table(common.BKTableNameUserAPI).Delete(params.Context, filter)
	if nil != err {
		blog.Errorf("delete user api fail, error information is %s, params:%v, rid: %s", err.Error(), filter, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBDeleteFailed)
	}

	return nil, nil
}

func (s *coreService) GetUserConfig(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	dat := new(meta.ObjQueryInput)
	if err := data.MarshalJSONInto(dat); err != nil {
		blog.Errorf("get user config failed with decode body, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	condition := make(map[string]interface{})
	if nil != dat.Condition {
		condition = dat.Condition.(map[string]interface{})
	}

	appID, err := util.GetInt64ByInterface(condition[common.BKAppIDField])
	if err != nil {
		blog.Errorf("get user config failed, invalid appid[%s], err: %v, rid: %s", common.BKAppIDField, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommParamsIsInvalid)
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
	condition = util.SetModOwner(condition, params.SupplierAccount)
	count, err := s.db.Table(common.BKTableNameUserAPI).Find(condition).Count(params.Context)
	if err != nil {
		blog.Errorf("get user api information failed, err:%v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	result := make([]interface{}, 0)
	err = s.db.Table(common.BKTableNameUserAPI).Find(condition).Fields(fieldArr...).Sort(sort).Start(uint64(start)).Limit(uint64(limit)).All(params.Context, &result)
	if err != nil {
		blog.Errorf("get user api information failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	return meta.UserConfigResult{
		Count: count,
		Info:  result,
	}, nil
}

func (s *coreService) UserConfigDetail(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	id := pathParams("id")
	appID, err := strconv.ParseInt(pathParams(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user[%s] config failed, invalid appid[%s], err: %v, rid: %s", id, common.BKAppIDField, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommParamsIsInvalid)
	}

	filter := common.KvMap{
		"id":                id,
		common.BKAppIDField: appID,
	}
	filter = util.SetModOwner(filter, params.SupplierAccount)
	result := new(meta.UserConfigMeta)
	err = s.db.Table(common.BKTableNameUserAPI).Find(filter).One(params.Context, result)
	if err != nil && !s.db.IsNotFoundError(err) {
		blog.Errorf("get user api information error,input:%v error:%v, rid: %s", id, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	return result, nil

}

func (s *coreService) AddUserCustom(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	ID := xid.New()
	data["id"] = ID.String()
	data["bk_user"] = pathParams("bk_user")
	data = util.SetModOwner(data, params.SupplierAccount)
	err := s.db.Table(common.BKTableNameUserCustom).Insert(params.Context, data)
	if nil != err {
		blog.Errorf("Create  user custom fail, err: %v, params:%v, rid: %s", err, data, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCreateUserCustom)
	}
	return nil, nil
}

func (s *coreService) UpdateUserCustomByID(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	conditons := make(map[string]interface{})
	conditons["id"] = pathParams("id")
	conditons["bk_user"] = pathParams("bk_user")
	conditons = util.SetModOwner(conditons, params.SupplierAccount)
	err := s.db.Table(common.BKTableNameUserCustom).Update(params.Context, conditons, data)
	if nil != err {
		blog.Errorf("update  user custom failed, err: %v, data:%v, rid: %s", err, data, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil, nil
}

func (s *coreService) GetUserCustomByUser(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	conds := make(map[string]interface{})
	conds["bk_user"] = pathParams("bk_user")
	conds = util.SetModOwner(conds, params.SupplierAccount)

	result := make(map[string]interface{})
	err := s.db.Table(common.BKTableNameUserCustom).Find(conds).One(params.Context, &result)
	if nil != err && !s.db.IsNotFoundError(err) {
		blog.Errorf("add  user custom failed, err: %v, params:%v, rid: %s", err, conds, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	return result, nil
}

func (s *coreService) GetDefaultUserCustom(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	conds := make(map[string]interface{})
	conds["is_default"] = 1
	conds = util.SetModOwner(conds, params.SupplierAccount)

	result := make(map[string]interface{})
	err := s.db.Table(common.BKTableNameUserCustom).Find(conds).One(params.Context, &result)
	if nil != err && !s.db.IsNotFoundError(err) {
		blog.Errorf("get default user custom fail, err: %v, params:%v, rid: %s, rid: %s", err, conds, params.ReqID, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	return result, nil
}
