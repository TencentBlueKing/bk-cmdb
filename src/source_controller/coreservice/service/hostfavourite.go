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
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"

	"github.com/rs/xid"
)

func (s *coreService) AddHostFavourite(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	user := pathParams("user")
	paras := new(meta.FavouriteParms)
	if err := data.MarshalJSONInto(paras); err != nil {
		blog.Errorf("add host favourite, but decode body failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	// check duplicate
	query := map[string]interface{}{
		"user":                user,
		"name":                paras.Name,
		common.BKAppIDField:   paras.BizID,
		common.BKOwnerIDField: params.SupplierAccount,
	}
	rowCount, err := s.db.Table(common.BKTableNameHostFavorite).Find(query).Count(params.Context)
	if err != nil {
		blog.Errorf("query host favorites fail, err: %v, params:%v, rid: %s", err, query, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostFavouriteQueryFail)
	}
	if rowCount > 0 {
		blog.Errorf("AddHostFavourite fail, name duplicated, params:%v, rid: %s", query, params.ReqID)
		return nil, params.Error.CCErrorf(common.CCErrCommDuplicateItem, paras.Name)
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
		OwnerID:     params.SupplierAccount,
		CreateTime:  time.Now().UTC(),
		UpdateTime:  time.Now().UTC(),
		BizID:       paras.BizID,
	}
	err = s.db.Table(common.BKTableNameHostFavorite).Insert(params.Context, fav)
	if err != nil {
		blog.Errorf("create host favorites failed, data:%+v error:%v, rid: %s", fav, err, params.ReqID)
		return nil, params.Error.CCErrorf(common.CCErrHostFavouriteCreateFail)
	}
	return meta.ID{ID: id}, nil
}

func (s *coreService) UpdateHostFavouriteByID(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	id := pathParams("id")
	user := pathParams("user")

	fav := new(meta.FavouriteMeta)
	if err := data.MarshalJSONInto(fav); err != nil {
		blog.Errorf("update host favourite with id[%s], but decode body failed, err: %v, rid: %s", id, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	fav.UpdateTime = time.Now().UTC()

	// check exist
	query := map[string]interface{}{
		"user":                user,
		"id":                  id,
		common.BKOwnerIDField: params.SupplierAccount,
	}
	dbData := make([]meta.FavouriteMeta, 0)
	err := s.db.Table(common.BKTableNameHostFavorite).Find(query).All(params.Context, &dbData)
	if nil != err {
		blog.Errorf("update host favorites with id[%s], but query failed, err: %v, params:%v, rid: %s", id, err, query, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostFavouriteQueryFail)
	}
	if len(dbData) != 1 {
		blog.V(5).Infof("update host favorites with id[%s], but favorites found: %+v, params:%v, rid: %s", id, dbData, query, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostFavouriteUpdateFail)
	}
	hostFavourite := dbData[0]

	// check name duplicate before update name field
	if len(fav.Name) != 0 {
		dupFilter := map[string]interface{}{
			"name":                fav.Name,
			common.BKUser:         user,
			common.BKFieldID:      common.KvMap{common.BKDBNE: id},
			common.BKOwnerIDField: params.SupplierAccount,
			common.BKAppIDField:   fav.BizID,
		}
		rowCount, err := s.db.Table(common.BKTableNameHostFavorite).Find(dupFilter).Count(params.Context)
		if err != nil {
			blog.Errorf("query user api validate name duplicate fail, err: %v, params:%v, rid: %s", err, dupFilter, params.ReqID)
			return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
		}
		if rowCount > 0 {
			blog.Errorf("host user api  name duplicate , params:%v, rid: %s", dupFilter, params.ReqID)
			return nil, params.Error.CCErrorf(common.CCErrCommDuplicateItem, "")
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
	if err = s.db.Table(common.BKTableNameHostFavorite).Update(params.Context, query, fav); err != nil {
		blog.Errorf("update host favorite failed, err: %v, params:%v, rid: %s", err, query, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostFavouriteUpdateFail)
	}

	return nil, nil
}

func (s *coreService) DeleteHostFavouriteByID(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	id := pathParams("id")
	user := pathParams("user")

	query := map[string]interface{}{
		"user":                user,
		"id":                  id,
		common.BKOwnerIDField: params.SupplierAccount,
	}
	rowCount, err := s.db.Table(common.BKTableNameHostFavorite).Find(query).Count(params.Context)
	if nil != err {
		blog.Errorf("delete host favorites with id[%s], but query failed, err: %v, params:%v, rid: %s", id, err, query, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostFavouriteDeleteFail)
	}

	if rowCount != 1 {
		blog.V(5).Infof("delete host favorites with id[%s], but favorites match: %d, params:%v, rid: %s", id, rowCount, query, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostFavouriteDeleteFail)
	}

	if err = s.db.Table(common.BKTableNameHostFavorite).Delete(params.Context, query); err != nil {
		blog.Errorf("delete host favorites with id[%s] failed, err: %v, params:%v, rid: %s", id, err, query, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostFavouriteDeleteFail)
	}

	blog.V(5).Infof("delete host favorites with id[%s] success, info: %v, rid: %s", err, query, params.ReqID)
	return nil, nil
}

func (s *coreService) ListHostFavourites(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	dat := new(meta.ObjQueryInput)
	if err := data.MarshalJSONInto(dat); err != nil {
		blog.Errorf("get host favourite, but decode body failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	// construct get condition
	condition := make(map[string]interface{})
	if nil != dat.Condition {
		condition = dat.Condition.(map[string]interface{})
	}
	condition["user"] = pathParams("user")
	condition[common.BKOwnerIDField] = params.SupplierAccount

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
	count, err := s.db.Table(common.BKTableNameHostFavorite).Find(condition).Count(params.Context)
	if err != nil {
		blog.Errorf("get host favorites failed,input:%+v error:%v, rid: %s", dat, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostFavouriteQueryFail)
	}

	if err = s.db.Table(common.BKTableNameHostFavorite).Find(condition).Fields(fieldArr...).Start(uint64(skip)).Limit(uint64(limit)).Sort(sort).All(params.Context, &resultData); err != nil {
		blog.Errorf("get host favorites failed,input:%+v error:%v, rid: %s", dat, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostFavouriteQueryFail)
	}

	result := meta.FavoriteResult{
		Count: count,
		Info:  resultData,
	}
	return result, nil

}

func (s *coreService) GetHostFavouriteByID(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	ID := pathParams("id")
	user := pathParams("user")
	if "" == ID || "0" == ID {
		blog.Errorf("get host favourite, but id is emtpy, rid: %s", params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommParamsNeedSet)
	}

	query := common.KvMap{
		"user":                user,
		"id":                  ID,
		common.BKOwnerIDField: params.SupplierAccount,
	}
	result := new(meta.FavouriteMeta)
	err := s.db.Table(common.BKTableNameHostFavorite).Find(query).One(params.Context, result)
	if err != nil && !s.db.IsNotFoundError(err) {
		blog.Errorf("get host favourite failed,input: %v error: %v, rid: %s", ID, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrHostFavouriteQueryFail)
	}

	return result, nil
}
