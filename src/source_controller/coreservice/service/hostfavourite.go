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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	meta "configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"

	"github.com/rs/xid"
)

func (s *coreService) AddHostFavourite(ctx *rest.Contexts) {
	user := ctx.Request.PathParameter("user")
	paras := new(meta.FavouriteParms)
	if err := ctx.DecodeInto(paras); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// check duplicate
	query := map[string]interface{}{
		"user":                user,
		"name":                paras.Name,
		common.BKAppIDField:   paras.BizID,
		common.BKOwnerIDField: ctx.Kit.SupplierAccount,
	}
	rowCount, err := mongodb.Client().Table(common.BKTableNameHostFavorite).Find(query).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("query host favorites fail, err: %v, ctx:%v, rid: %s", err, query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteQueryFail))
		return
	}
	if rowCount > 0 {
		blog.Errorf("AddHostFavourite fail, name duplicated, ctx:%v, rid: %s", query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, paras.Name))
		return
	}

	// insert
	id := xid.New().String()
	fav := meta.FavouriteMeta{
		ID:          id,
		Info:        paras.Info,
		Name:        paras.Name,
		Count:       1,
		User:        user,
		QueryParams: paras.QueryParams,
		OwnerID:     ctx.Kit.SupplierAccount,
		CreateTime:  time.Now().UTC(),
		UpdateTime:  time.Now().UTC(),
		BizID:       paras.BizID,
	}
	err = mongodb.Client().Table(common.BKTableNameHostFavorite).Insert(ctx.Kit.Ctx, fav)
	if err != nil {
		blog.Errorf("create host favorites failed, data:%+v error:%v, rid: %s", fav, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteCreateFail))
		return
	}
	ctx.RespEntity(meta.ID{ID: id})
}

func (s *coreService) UpdateHostFavouriteByID(ctx *rest.Contexts) {
	id := ctx.Request.PathParameter("id")
	user := ctx.Request.PathParameter("user")

	fav := new(meta.FavouriteMeta)
	if err := ctx.DecodeInto(fav); err != nil {
		ctx.RespAutoError(err)
		return
	}
	fav.UpdateTime = time.Now().UTC()

	// check exist
	query := map[string]interface{}{
		"user":                user,
		"id":                  id,
		common.BKOwnerIDField: ctx.Kit.SupplierAccount,
	}
	dbData := make([]meta.FavouriteMeta, 0)
	err := mongodb.Client().Table(common.BKTableNameHostFavorite).Find(query).All(ctx.Kit.Ctx, &dbData)
	if nil != err {
		blog.Errorf("update host favorites with id[%s], but query failed, err: %v, ctx:%v, rid: %s", id, err, query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteQueryFail))
		return
	}
	if len(dbData) != 1 {
		blog.V(5).Infof("update host favorites with id[%s], but favorites found: %+v, ctx:%v, rid: %s", id, dbData, query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteUpdateFail))
		return
	}
	hostFavourite := dbData[0]

	// check name duplicate before update name field
	if len(fav.Name) != 0 {
		dupFilter := map[string]interface{}{
			"name":                fav.Name,
			common.BKUser:         user,
			common.BKFieldID:      common.KvMap{common.BKDBNE: id},
			common.BKOwnerIDField: ctx.Kit.SupplierAccount,
			common.BKAppIDField:   fav.BizID,
		}
		rowCount, err := mongodb.Client().Table(common.BKTableNameHostFavorite).Find(dupFilter).Count(ctx.Kit.Ctx)
		if err != nil {
			blog.Errorf("query user api validate name duplicate fail, err: %v, ctx:%v, rid: %s", err, dupFilter, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}
		if rowCount > 0 {
			blog.Errorf("host user api  name duplicate , ctx:%v, rid: %s", dupFilter, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, ""))
			return
		}
	}

	// 只有部分字段能更新
	hostFavourite.UpdateTime = time.Now().UTC()
	hostFavourite.Count = fav.Count
	if len(fav.Name) > 0 {
		hostFavourite.Name = fav.Name
	}
	if len(fav.Info) > 0 {
		hostFavourite.Info = fav.Info
	}
	if len(fav.QueryParams) > 0 {
		hostFavourite.QueryParams = fav.QueryParams
	}

	// do update host favorite
	if err = mongodb.Client().Table(common.BKTableNameHostFavorite).Update(ctx.Kit.Ctx, query, fav); err != nil {
		blog.Errorf("update host favorite failed, err: %v, ctx:%v, rid: %s", err, query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteUpdateFail))
		return
	}

	ctx.RespEntity(nil)
}

func (s *coreService) DeleteHostFavouriteByID(ctx *rest.Contexts) {
	id := ctx.Request.PathParameter("id")
	user := ctx.Request.PathParameter("user")

	query := map[string]interface{}{
		"user":                user,
		"id":                  id,
		common.BKOwnerIDField: ctx.Kit.SupplierAccount,
	}
	rowCount, err := mongodb.Client().Table(common.BKTableNameHostFavorite).Find(query).Count(ctx.Kit.Ctx)
	if nil != err {
		blog.Errorf("delete host favorites with id[%s], but query failed, err: %v, ctx:%v, rid: %s", id, err, query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteDeleteFail))
		return
	}

	if rowCount != 1 {
		blog.V(5).Infof("delete host favorites with id[%s], but favorites match: %d, ctx:%v, rid: %s", id, rowCount, query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteDeleteFail))
		return
	}

	if err = mongodb.Client().Table(common.BKTableNameHostFavorite).Delete(ctx.Kit.Ctx, query); err != nil {
		blog.Errorf("delete host favorites with id[%s] failed, err: %v, ctx:%v, rid: %s", id, err, query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteDeleteFail))
		return
	}

	blog.V(5).Infof("delete host favorites with id[%s] success, info: %v, rid: %s", err, query, ctx.Kit.Rid)
	ctx.RespEntity(nil)
}

func (s *coreService) ListHostFavourites(ctx *rest.Contexts) {
	dat := new(meta.ObjQueryInput)
	if err := ctx.DecodeInto(dat); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// construct get condition
	condition := make(map[string]interface{})
	if nil != dat.Condition {
		condition = dat.Condition.(map[string]interface{})
	}
	condition["user"] = ctx.Request.PathParameter("user")
	condition[common.BKOwnerIDField] = ctx.Kit.SupplierAccount

	// read fields and page
	fieldArr := []string{"id", "info", "query_params", "name", "is_default", common.CreateTimeField, "count"}
	if "" != dat.Fields {
		fieldArr = strings.Split(dat.Fields, ",")
	}
	skip, limit, sort := dat.Start, dat.Limit, dat.Sort
	if limit == 0 {
		limit = 20
	}
	if len(sort) == 0 {
		sort = common.CreateTimeField
	}

	resultData := make([]map[string]interface{}, 0)
	count, err := mongodb.Client().Table(common.BKTableNameHostFavorite).Find(condition).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("get host favorites failed,input:%+v error:%v, rid: %s", dat, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteQueryFail))
		return
	}

	if err = mongodb.Client().Table(common.BKTableNameHostFavorite).Find(condition).Fields(fieldArr...).Start(uint64(skip)).Limit(uint64(limit)).Sort(sort).All(ctx.Kit.Ctx, &resultData); err != nil {
		blog.Errorf("get host favorites failed,input:%+v error:%v, rid: %s", dat, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteQueryFail))
		return
	}

	result := meta.FavoriteResult{
		Count: count,
		Info:  resultData,
	}
	ctx.RespEntity(result)

}

func (s *coreService) GetHostFavouriteByID(ctx *rest.Contexts) {
	ID := ctx.Request.PathParameter("id")
	user := ctx.Request.PathParameter("user")
	if "" == ID || "0" == ID {
		blog.Errorf("get host favourite, but id is emtpy, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsNeedSet))
		return
	}

	query := common.KvMap{
		"user":                user,
		"id":                  ID,
		common.BKOwnerIDField: ctx.Kit.SupplierAccount,
	}
	result := new(meta.FavouriteMeta)
	err := mongodb.Client().Table(common.BKTableNameHostFavorite).Find(query).One(ctx.Kit.Ctx, result)
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("get host favourite failed,input: %v error: %v, rid: %s", ID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostFavouriteQueryFail))
		return
	}

	ctx.RespEntity(result)
}
