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
	"net/http"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"
	"github.com/rs/xid"
)

var (
	userCustomTableName    string = "cc_UserCustom"
	mgo_on_not_found_error string = "not found"
)

func (s *Service) AddHostFavourite(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := req.PathParameter("user")
	ownerID := util.GetOwnerID(pheader)

	paras := new(meta.FavouriteParms)
	if err := json.NewDecoder(req.Request.Body).Decode(paras); err != nil {
		blog.Errorf("add host favourite, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	query := common.KvMap{"user": user, "name": paras.Name}
	query = util.SetModOwner(query, ownerID)
	rowCount, err := s.Instance.GetCntByCondition(FavouriteCollection, query)
	if err != nil {
		blog.Error("query host favorites fail, err: %v, params:%v", err, query)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
		return
	}

	if 0 != rowCount {
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteCreateFail)})
		return
	}

	id := xid.New().String()
	fav := meta.FavouriteMeta{
		ID:          id,
		Info:        paras.Info,
		Name:        paras.Name,
		Count:       1,
		User:        user,
		IsDefault:   paras.IsDefault,
		QueryParams: paras.QueryParams,
		OwnerID:     ownerID,
		CreateTime:  time.Now().UTC(),
		UpdateTime:  time.Now().UTC(),
	}
	_, err = s.Instance.Insert(FavouriteCollection, fav)
	if err != nil {
		blog.Errorf("create host favorites failed, data:%+v error:%v", fav, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteCreateFail)})
		return
	}
	resp.WriteEntity(meta.IDResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     meta.ID{ID: id},
	})
}

func (s *Service) UpdateHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	id := req.PathParameter("id")
	user := req.PathParameter("user")
	ownerID := util.GetOwnerID(pheader)

	fav := new(meta.FavouriteMeta)
	if err := json.NewDecoder(req.Request.Body).Decode(fav); err != nil {
		blog.Errorf("update host favourite with id[%s], but decode body failed, err: %v", id, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	fav.UpdateTime = time.Now().UTC()

	query := common.KvMap{"user": user, "id": id}
	query = util.SetModOwner(query, ownerID)
	rowCount, err := s.Instance.GetCntByCondition(FavouriteCollection, query)
	if nil != err {
		blog.Error("update host favorites with id[%s], but query failed, err: %v, params:%v", id, err, query)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
		return
	}

	if 1 != rowCount {
		blog.Info("update host favorites with id[%s], but favorites not exists, params:%v", id, query)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteUpdateFail)})
		return
	}

	if len(fav.Name) != 0 {
		dupParams := make(map[string]interface{})
		dupParams["name"] = fav.Name
		dupParams[common.BKUser] = user
		dupParams[common.BKFieldID] = common.KvMap{common.BKDBNE: id}
		dupParams = util.SetModOwner(dupParams, ownerID)
		rowCount, err := s.Instance.GetCntByCondition(FavouriteCollection, dupParams)
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
	err = s.Instance.UpdateByCondition(FavouriteCollection, fav, query)
	if nil != err {
		blog.Error("update host favorites fail, err: %v, params:%v", err, query)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteUpdateFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (s *Service) DeleteHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	id := req.PathParameter("id")
	user := req.PathParameter("user")
	ownerID := util.GetOwnerID(pheader)

	query := common.KvMap{"user": user, "id": id}
	query = util.SetModOwner(query, ownerID)
	rowCount, err := s.Instance.GetCntByCondition(FavouriteCollection, query)
	if nil != err {
		blog.Error("delete host favorites with id[%s], but query failed, err: %v, params:%v", id, err, query)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteDeleteFail)})
		return
	}

	if 1 != rowCount {
		blog.Info("delete host favorites with id[%s], but favorites not exists, params:%v", id, query)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteDeleteFail)})
		return
	}

	err = s.Instance.DelByCondition(FavouriteCollection, query)
	if nil != err {
		blog.Error("delete host favorites with id[%s] failed, err: %v, params:%v", err, query)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteDeleteFail)})
		return
	}
	resp.WriteEntity(meta.NewSuccessResp(nil))
	blog.Infof("delete host favorites with id[%s] success, info: %v", err, query)
}

func (s *Service) GetHostFavourites(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	dat := new(meta.ObjQueryInput)
	if err := json.NewDecoder(req.Request.Body).Decode(dat); err != nil {
		blog.Errorf("get host favourite, but decode body failed, err: %v", err)
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
	condition = util.SetModOwner(condition, ownerID)
	result := make([]map[string]interface{}, 0)
	count, err := s.Instance.GetCntByCondition(FavouriteCollection, condition)
	if err != nil {
		blog.Errorf("get host favorites failed,input:%+v error:%v", dat, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
		return
	}

	err = s.Instance.GetMutilByCondition(FavouriteCollection, fieldArr, condition, &result, sort, skip, limit)
	if err != nil {
		blog.Errorf("get host favorites failed,input:%+v error:%v", dat, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
		return
	}

	resp.WriteEntity(meta.GetHostFavoriteResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     meta.FavoriteResult{Count: count, Info: result},
	})

}

func (s *Service) GetHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ID := req.PathParameter("id")
	user := req.PathParameter("user")
	ownerID := util.GetOwnerID(pheader)

	if "" == ID || "0" == ID {
		blog.Errorf("get host favourite, but id is emtpy")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsNeedSet)})
		return
	}

	query := common.KvMap{"user": user, "id": ID}
	query = util.SetModOwner(query, ownerID)
	result := new(meta.FavouriteMeta)
	err := s.Instance.GetOneByCondition(FavouriteCollection, nil, query, result)
	if err != nil && mgo_on_not_found_error != err.Error() {
		blog.Errorf("get host favourite failed,input: %v error: %v", ID, err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostFavouriteQueryFail)})
		return
	}

	resp.WriteEntity(meta.GetHostFavoriteWithIDResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *result,
	})
}
