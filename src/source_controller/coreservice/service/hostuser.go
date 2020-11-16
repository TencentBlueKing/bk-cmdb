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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"

	"github.com/rs/xid"
)

func (s *coreService) AddUserCustom(ctx *rest.Contexts) {
	ID := xid.New()
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	data["id"] = ID.String()
	data["bk_user"] = ctx.Kit.User
	if strings.Contains(ctx.Kit.User, ".") {
		transformedData := make(map[string]interface{}, len(data))
		for key, value := range data {
			transformedData[strings.Replace(key, ".", "\u002e", -1)] = value
		}
		data = transformedData
	}

	data = util.SetModOwner(data, ctx.Kit.SupplierAccount)
	err := mongodb.Client().Table(common.BKTableNameUserCustom).Insert(ctx.Kit.Ctx, data)
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
	if strings.Contains(ctx.Request.PathParameter("bk_user"), ".") {
		transformedData := make(map[string]interface{}, len(data))
		for key, value := range data {
			transformedData[strings.Replace(key, ".", "\u002e", -1)] = value
		}
		data = transformedData
	}
	conditons = util.SetModOwner(conditons, ctx.Kit.SupplierAccount)
	err := mongodb.Client().Table(common.BKTableNameUserCustom).Update(ctx.Kit.Ctx, conditons, data)
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
	err := mongodb.Client().Table(common.BKTableNameUserCustom).Find(conds).One(ctx.Kit.Ctx, &result)
	if nil != err && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("add  user custom failed, err: %v, ctx:%v, rid: %s", err, conds, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	if strings.Contains(ctx.Request.PathParameter("bk_user"), ".") {
		transformedData := make(map[string]interface{}, len(result))
		for key, value := range result {
			transformedData[strings.Replace(key, "\u002e", ".", -1)] = value
		}
		result = transformedData
	}

	ctx.RespEntity(result)
}

// GetDefaultUserCustom  find user custom set table header for any object
func (s *coreService) GetDefaultUserCustom(ctx *rest.Contexts) {
	conds := make(map[string]interface{})
	conds[common.BKDefaultField] = 1
	conds = util.SetModOwner(conds, ctx.Kit.SupplierAccount)

	result := make(map[string]interface{})
	err := mongodb.Client().Table(common.BKTableNameUserCustom).Find(conds).One(ctx.Kit.Ctx, &result)
	if nil != err && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("get default user custom fail, err: %v, ctx:%v, rid: %s, rid: %s", err, conds, ctx.Kit.Rid, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	ctx.RespEntity(result)
}

// UpdateDefaultUserCustom update user custom set table header for any object
func (s *coreService) UpdateDefaultUserCustom(ctx *rest.Contexts) {
	conditions := make(map[string]interface{})
	conditions[common.BKDefaultField] = 1
	conditions = util.SetModOwner(conditions, ctx.Kit.SupplierAccount)
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	data[common.ModifierField] = ctx.Kit.User
	data[common.LastTimeField] = util.GetCurrentTimePtr()
	err := mongodb.Client().Table(common.BKTableNameUserCustom).Upsert(ctx.Kit.Ctx, conditions, data)
	if nil != err {
		blog.Errorf("update  default custom failed, err: %v, data:%v, rid: %s", err, data, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
		return
	}
	ctx.RespEntity(nil)
}
