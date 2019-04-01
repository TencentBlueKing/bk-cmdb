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
	"context"
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"

	"github.com/emicklei/go-restful"
)

// CreateObjectUnique create object's unique
func (cli *Service) CreateObjectUnique(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	objID := req.PathParameter(common.BKObjIDField)
	var dat metadata.CreateUniqueRequest
	if body, err := util.DecodeJSON(req.Request.Body, &dat); err != nil {
		blog.Errorf("[CreateObjectUnique] DecodeJSON error: %v, %s", err, body)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	for _, key := range dat.Keys {
		switch key.Kind {
		case metadata.UinqueKeyKindProperty:
		default:
			blog.Errorf("[CreateObjectUnique] invalid key kind: %s", key.Kind)
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)})
			return
		}
	}

	if dat.MustCheck {
		cond := condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(objID)
		cond.Field("must_check").Eq(true)
		count, err := db.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).Count(ctx)
		if nil != err {
			blog.Errorf("[CreateObjectUnique] check must check error: %v", err)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
			return
		}
		if count > 0 {
			blog.Errorf("[CreateObjectUnique] model could not have multiple must check unique")
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoObjectUniqueCanNotHasMutiMustCheck)})
			return
		}
	}

	err := recheckUniqueForExistsInsts(ctx, db, ownerID, objID, dat.Keys, dat.MustCheck)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] recheckUniqueForExistsInsts for %s with %v error: %v", objID, dat, err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDuplicateItem)})
		return
	}

	id, err := db.NextSequence(ctx, common.BKTableNameObjUnique)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] NextSequence error: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
		return
	}

	unique := metadata.ObjectUnique{
		ID:        id,
		ObjID:     objID,
		MustCheck: dat.MustCheck,
		Keys:      dat.Keys,
		Ispre:     false,
		OwnerID:   ownerID,
		LastTime:  metadata.Now(),
	}

	err = db.Table(common.BKTableNameObjUnique).Insert(ctx, &unique)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] Insert error: %v, raw: %#v", err, &unique)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
		return
	}

	resp.WriteEntity(metadata.CreateUniqueResult{BaseResp: metadata.SuccessBaseResp, Data: metadata.RspID{ID: int64(id)}})
}

// UpdateObjectUnique update object's unique
func (cli *Service) UpdateObjectUnique(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	objID := req.PathParameter(common.BKObjIDField)
	id, err := strconv.ParseUint(req.PathParameter("id"), 10, 64)
	if err != nil {
		blog.Errorf("[UpdateObjectUnique] path param error: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "id")})
		return
	}

	var unique metadata.UpdateUniqueRequest
	if body, err := util.DecodeJSON(req.Request.Body, &unique); err != nil {
		blog.Errorf("[UpdateObjectUnique] DecodeJSON error: %v, %s", err, body)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	unique.LastTime = metadata.Now()

	for _, key := range unique.Keys {
		switch key.Kind {
		case metadata.UinqueKeyKindProperty:
		default:
			blog.Errorf("[UpdateObjectUnique] invalid key kind: %s", key.Kind)
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)})
			return
		}
	}

	if unique.MustCheck {
		cond := condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(objID)
		cond.Field("must_check").Eq(true)
		cond.Field("id").NotEq(id)
		count, err := db.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).Count(ctx)
		if nil != err {
			blog.Errorf("[UpdateObjectUnique] check must check  error: %v", err)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
			return
		}
		if count > 0 {
			blog.Errorf("[UpdateObjectUnique] model could not have multiple must check unique")
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoObjectUniqueCanNotHasMutiMustCheck)})
			return
		}
	}

	err = recheckUniqueForExistsInsts(ctx, db, ownerID, objID, unique.Keys, unique.MustCheck)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] recheckUniqueForExistsInsts for %s with %v error: %v", objID, unique, err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDuplicateItem)})
		return
	}

	cond := condition.CreateCondition()
	cond.Field("id").Eq(id)
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ownerID)

	oldunique := metadata.ObjectUnique{}
	err = db.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).One(ctx, &oldunique)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] find error: %s, raw: %#v", err, cond.ToMapStr())
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
		return
	}

	if oldunique.Ispre {
		blog.Errorf("[UpdateObjectUnique] could not update preset constrain: %+v, %v", oldunique, err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoObjectUniquePresetCouldNotDelOrEdit)})
		return
	}

	err = db.Table(common.BKTableNameObjUnique).Update(ctx, cond.ToMapStr(), &unique)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] Update error: %s, raw: %#v", err, &unique)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
		return
	}

	resp.WriteEntity(metadata.UpdateUniqueResult{BaseResp: metadata.SuccessBaseResp})
}

// DeleteObjectUnique delte object's unique
func (cli *Service) DeleteObjectUnique(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	objID := req.PathParameter(common.BKObjIDField)
	id, err := strconv.ParseUint(req.PathParameter("id"), 10, 64)
	if err != nil {
		blog.Errorf("[DeleteObjectUnique] path param [id] error: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "id")})
		return
	}

	cond := condition.CreateCondition()
	cond.Field("id").Eq(id)
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ownerID)

	unique := metadata.ObjectUnique{}
	err = db.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).One(ctx, &unique)
	if nil != err {
		blog.Errorf("[DeleteObjectUnique] find error: %s, raw: %#v", err, cond.ToMapStr())
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
		return
	}

	if unique.Ispre {
		blog.Errorf("[DeleteObjectUnique] could not delete preset constrain:%+v, %s", unique, err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoObjectUniquePresetCouldNotDelOrEdit)})
		return
	}

	err = db.Table(common.BKTableNameObjUnique).Delete(ctx, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[DeleteObjectUnique] Delete error: %s, raw: %#v", err, cond.ToMapStr())
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
		return
	}

	resp.WriteEntity(metadata.DeleteUniqueResult{BaseResp: metadata.SuccessBaseResp})
}

// SearchObjectUnique delte object's unique
func (cli *Service) SearchObjectUnique(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)
	ctx := util.GetDBContext(context.Background(), req.Request.Header)
	db := cli.Instance.Clone()

	objID := req.PathParameter(common.BKObjIDField)

	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ownerID)

	uniques, err := cli.searchObjectUnique(ctx, db, ownerID, objID)
	if nil != err {
		blog.Errorf("[SearchObjectUnique] Search error: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrObjectDBOpErrno)})
		return
	}

	resp.WriteEntity(metadata.SearchUniqueResult{BaseResp: metadata.SuccessBaseResp, Data: uniques})
}

func (cli *Service) searchObjectUnique(ctx context.Context, db dal.RDB, ownerID, objID string) ([]metadata.ObjectUnique, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ownerID)

	uniques := []metadata.ObjectUnique{}
	err := db.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).All(ctx, &uniques)
	return uniques, err
}

func recheckUniqueForExistsInsts(ctx context.Context, db dal.RDB, ownerID, objID string, keys []metadata.UinqueKey, mustCheck bool) error {
	propertyIDs := []uint64{}
	for _, key := range keys {
		switch key.Kind {
		case metadata.UinqueKeyKindProperty:
			propertyIDs = append(propertyIDs, key.ID)
		default:
			return fmt.Errorf("invalid key kind: %s", key.Kind)
		}
	}

	propertys := []metadata.Attribute{}
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ownerID)
	cond.Field(common.BKFieldID).In(propertyIDs)
	err := db.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(ctx, &propertys)
	if err != nil {
		blog.ErrorJSON("[ObjectUnique] recheckUniqueForExistsInsts find propertys for %s failed %s: %s", objID, err)
		return err
	}

	keynames := []string{}
	for _, property := range propertys {
		keynames = append(keynames, property.PropertyID)
	}
	if len(keynames) <= 0 {
		blog.Warnf("[ObjectUnique] recheckUniqueForExistsInsts keys empty for [%s] %+v", objID, keys)
		return nil
	}

	pipeline := []interface{}{}

	instcond := mapstr.MapStr{
		common.BKObjIDField: objID,
	}
	if common.GetObjByType(objID) == common.BKInnerObjIDObject {
		instcond.Set(common.BKObjIDField, objID)
	}

	if !mustCheck {
		matchs := []mapstr.MapStr{}
		for _, key := range keynames {
			matchs = append(matchs, mapstr.MapStr{key: mapstr.MapStr{common.BKDBNE: nil}})
		}
		if len(matchs) > 0 {
			instcond.Set(common.BKDBOR, matchs)
		}
	}

	pipeline = append(pipeline, mapstr.MapStr{common.BKDBMatch: instcond})

	group := mapstr.MapStr{}
	for _, key := range keynames {
		group.Set(key, "$"+key)
	}
	pipeline = append(pipeline, mapstr.MapStr{
		common.BKDBGroup: mapstr.MapStr{
			"_id":   group,
			"total": mapstr.MapStr{common.BKDBSum: 1},
		},
	})

	pipeline = append(pipeline, mapstr.MapStr{common.BKDBMatch: mapstr.MapStr{
		"_id":   mapstr.MapStr{common.BKDBNE: nil},
		"total": mapstr.MapStr{common.BKDBGT: 1},
	}})

	pipeline = append(pipeline, mapstr.MapStr{common.BKDBCount: "finded"})

	result := struct {
		Finded uint64 `bson:"finded"`
	}{}
	err = db.Table(common.GetInstTableName(objID)).AggregateOne(ctx, pipeline, &result)
	if err != nil && !db.IsNotFoundError(err) {
		blog.ErrorJSON("[ObjectUnique] recheckUniqueForExistsInsts failed %s, pipeline: %s", err, pipeline)
		return err
	}

	if result.Finded > 0 {
		return dal.ErrDuplicated
	}

	return nil
}
