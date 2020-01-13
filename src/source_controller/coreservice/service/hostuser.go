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
	"configcenter/src/common/http/rest"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"github.com/rs/xid"
)

func (s *coreService) AddUserConfig(ctx *rest.Contexts) {
	addQuery := new(meta.AddConfigQuery)
	if err := ctx.DecodeInto(addQuery); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(addQuery.Name) == 0 {
		blog.Error("parameter Name is required, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "Name"))
		return
	}

	if 0 >= addQuery.AppID {
		blog.Error("add user config, parameter app id is required, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAppIDField))
		return
	}
	if len(addQuery.CreateUser) == 0 {
		blog.Error("add user config, parameter CreateUser is required, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "create_user"))
		return
	}

	filter := common.KvMap{
		"name":              addQuery.Name,
		common.BKAppIDField: addQuery.AppID,
	}
	filter = util.SetModOwner(filter, ctx.Kit.SupplierAccount)
	rowCount, err := s.db.Table(common.BKTableNameUserAPI).Find(filter).Count(ctx.Kit.Ctx)
	if nil != err {
		blog.Errorf("add user config, query user api fail, error information is %s, filter: %v, rid: %s", err.Error(), filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if 0 != rowCount {
		blog.Errorf("add user config, [%s] user api is exist, rid: %s", addQuery.Name, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, ""))
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
		OwnerID:    ctx.Kit.SupplierAccount,
		ModifyUser: addQuery.CreateUser,
		UpdateTime: time.Now().UTC(),
	}

	err = s.db.Table(common.BKTableNameUserAPI).Insert(ctx.Kit.Ctx, userQuery)
	if err != nil {
		blog.Errorf("add user config, create user query failed, query:%+v err:%v, rid: %s", userQuery, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
		return
	}
	ctx.RespEntity(meta.ID{ID: id})
}

func (s *coreService) UpdateUserConfig(ctx *rest.Contexts) {
	id := ctx.Request.PathParameter("id")
	appID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user[%s] config failed, invalid appid[%s], err: %v, rid: %s", id, common.BKAppIDField, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	dat := new(meta.UserConfigMeta)
	if err := ctx.DecodeInto(dat); err != nil {
		ctx.RespAutoError(err)
		return
	}

	filter := common.KvMap{
		"id":                id,
		common.BKAppIDField: appID,
	}
	filter = util.SetModOwner(filter, ctx.Kit.SupplierAccount)
	rowCount, err := s.db.Table(common.BKTableNameUserAPI).Find(filter).Count(ctx.Kit.Ctx)
	if nil != err {
		blog.Errorf("query user api fail, error information is %s, ctx:%v, rid: %s", err.Error(), ctx, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if 1 != rowCount {
		blog.V(5).Infof("update user api config not permissions or not exists, ctx:%v, rid: %s", ctx, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommNotFound))
		return
	}

	if len(dat.Name) != 0 {
		dupParams := common.KvMap{
			"name":              dat.Name,
			common.BKAppIDField: appID,
			common.BKFieldID:    common.KvMap{common.BKDBNE: id},
		}
		dupParams = util.SetModOwner(dupParams, ctx.Kit.SupplierAccount)
		rowCount, getErr := s.db.Table(common.BKTableNameUserAPI).Find(dupParams).Count(ctx.Kit.Ctx)
		if nil != getErr {
			blog.Errorf("query user api validate name duplicate fail, error information is %s, ctx:%v, rid: %s", getErr.Error(), dupParams, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}
		if 0 < rowCount {
			blog.V(5).Infof("host user api  name duplicate , ctx:%v, rid: %s", dupParams, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, ""))
			return
		}
	}

	dat.UpdateTime = time.Now().UTC()
	dat.ModifyUser = util.GetUser(ctx.Kit.Header)
	dat.AppID = appID
	dat.OwnerID = ctx.Kit.SupplierAccount
	err = s.db.Table(common.BKTableNameUserAPI).Update(ctx.Kit.Ctx, filter, dat)
	if nil != err {
		blog.Errorf("update user api fail, error information is %s, ctx:%v, rid: %s", err.Error(), ctx, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}

	ctx.RespEntity(nil)
}

func (s *coreService) DeleteUserConfig(ctx *rest.Contexts) {
	id := ctx.Request.PathParameter("id")
	appID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user[%s] config failed, invalid appid[%s], err: %v, rid: %s", id, common.BKAppIDField, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	filter := common.KvMap{"id": id, common.BKAppIDField: appID}
	filter = util.SetModOwner(filter, ctx.Kit.SupplierAccount)
	rowCount, err := s.db.Table(common.BKTableNameUserAPI).Find(filter).Count(ctx.Kit.Ctx)
	if nil != err {
		blog.Errorf("query user api fail, error information is %s, ctx:%v, rid: %s", err.Error(), filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if 1 != rowCount {
		blog.V(5).Infof("host user api not permissions or not exists, ctx:%v, rid: %s", filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommNotFound))
		return
	}

	err = s.db.Table(common.BKTableNameUserAPI).Delete(ctx.Kit.Ctx, filter)
	if nil != err {
		blog.Errorf("delete user api fail, error information is %s, ctx:%v, rid: %s", err.Error(), filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
		return
	}

	ctx.RespEntity(nil)
}

func (s *coreService) GetUserConfig(ctx *rest.Contexts) {
	dat := new(meta.ObjQueryInput)
	if err := ctx.DecodeInto(dat); err != nil {
		blog.Errorf("get user config failed with decode body, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	condition := make(map[string]interface{})
	if nil != dat.Condition {
		condition = dat.Condition.(map[string]interface{})
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

	condition = util.SetModOwner(condition, ctx.Kit.SupplierAccount)
	count, err := s.db.Table(common.BKTableNameUserAPI).Find(condition).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("get user api information failed, err:%v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	result := make([]interface{}, 0)
	err = s.db.Table(common.BKTableNameUserAPI).Find(condition).Fields(fieldArr...).Sort(sort).Start(uint64(start)).Limit(uint64(limit)).All(ctx.Kit.Ctx, &result)
	if err != nil {
		blog.Errorf("get user api information failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	ctx.RespEntity(meta.UserConfigResult{
		Count: count,
		Info:  result,
	})
}

func (s *coreService) UserConfigDetail(ctx *rest.Contexts) {
	id := ctx.Request.PathParameter("id")
	appID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user[%s] config failed, invalid appid[%s], err: %v, rid: %s", id, common.BKAppIDField, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	filter := common.KvMap{
		"id":                id,
		common.BKAppIDField: appID,
	}
	filter = util.SetModOwner(filter, ctx.Kit.SupplierAccount)
	result := new(meta.UserConfigMeta)
	err = s.db.Table(common.BKTableNameUserAPI).Find(filter).One(ctx.Kit.Ctx, result)
	if err != nil && !s.db.IsNotFoundError(err) {
		blog.Errorf("get user api information error,input:%v error:%v, rid: %s", id, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	ctx.RespEntity(result)

}

func (s *coreService) AddUserCustom(ctx *rest.Contexts) {
	ID := xid.New()
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	data["id"] = ID.String()
	data["bk_user"] = ctx.Kit.User
	data = util.SetModOwner(data, ctx.Kit.SupplierAccount)
	err := s.db.Table(common.BKTableNameUserCustom).Insert(ctx.Kit.Ctx, data)
	if nil != err {
		blog.Errorf("Create  user custom fail, err: %v, ctx:%v, rid: %s", err, data, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCreateUserCustom))
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) UpdateUserCustomByID(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	conditons := make(map[string]interface{})
	conditons["id"] = ctx.Request.PathParameter("id")
	conditons["bk_user"] = ctx.Request.PathParameter("bk_user")
	conditons = util.SetModOwner(conditons, ctx.Kit.SupplierAccount)
	err := s.db.Table(common.BKTableNameUserCustom).Update(ctx.Kit.Ctx, conditons, data)
	if nil != err {
		blog.Errorf("update  user custom failed, err: %v, data:%v, rid: %s", err, data, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) GetUserCustomByUser(ctx *rest.Contexts) {
	conds := make(map[string]interface{})
	conds["bk_user"] = ctx.Kit.User
	conds = util.SetModOwner(conds, ctx.Kit.SupplierAccount)

	result := make(map[string]interface{})
	err := s.db.Table(common.BKTableNameUserCustom).Find(conds).One(ctx.Kit.Ctx, &result)
	if nil != err && !s.db.IsNotFoundError(err) {
		blog.Errorf("add  user custom failed, err: %v, ctx:%v, rid: %s", err, conds, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	ctx.RespEntity(result)
}

// GetDefaultUserCustom  find user custom set table heaher for any object
func (s *coreService) GetDefaultUserCustom(ctx *rest.Contexts) {
	conds := make(map[string]interface{})
	conds[common.BKDefaultField] = 1
	conds = util.SetModOwner(conds, ctx.Kit.SupplierAccount)

	result := make(map[string]interface{})
	err := s.db.Table(common.BKTableNameUserCustom).Find(conds).One(ctx.Kit.Ctx, &result)
	if nil != err && !s.db.IsNotFoundError(err) {
		blog.Errorf("get default user custom fail, err: %v, ctx:%v, rid: %s, rid: %s", err, conds, ctx.Kit.Rid, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	ctx.RespEntity(result)
}

// UpdatDefaultUserCustom update user custom set table header for any object
func (s *coreService) UpdatDefaultUserCustom(ctx *rest.Contexts) {
	conditons := make(map[string]interface{})
	conditons[common.BKDefaultField] = 1
	conditons = util.SetModOwner(conditons, ctx.Kit.SupplierAccount)
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	data[common.ModifierField] = ctx.Kit.User
	data[common.LastTimeField] = util.GetCurrentTimePtr()
	err := s.db.Table(common.BKTableNameUserCustom).Upsert(ctx.Kit.Ctx, conditons, data)
	if nil != err {
		blog.Errorf("update  default custom failed, err: %v, data:%v, rid: %s", err, data, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}
	ctx.RespEntity(nil)
}
