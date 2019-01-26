/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

func (m *modelAttrUnique) searchModelAttrUnique(ctx core.ContextParams, inputParam metadata.QueryCondition) (results []metadata.ObjectUnique, err error) {
	results = []metadata.ObjectUnique{}
	instHandler := m.dbProxy.Table(common.BKTableNameObjUnique).Find(inputParam.Condition)
	for _, sort := range inputParam.SortArr {
		fileld := sort.Field
		if sort.IsDsc {
			fileld = "-" + fileld
		}
		instHandler = instHandler.Sort(fileld)
	}
	err = instHandler.Start(uint64(inputParam.Limit.Offset)).Limit(uint64(inputParam.Limit.Limit)).All(ctx, &results)

	return results, err
}

func (m *modelAttrUnique) countModelAttrUnique(ctx core.ContextParams, cond mapstr.MapStr) (count uint64, err error) {

	count, err = m.dbProxy.Table(common.BKTableNameObjUnique).Find(cond).Count(ctx)

	return count, err
}

func (m *modelAttrUnique) createModelAttrUnique(ctx core.ContextParams, objID string, inputParam metadata.CreateModelAttrUnique) (uint64, error) {
	for _, key := range inputParam.Data.Keys {
		switch key.Kind {
		case metadata.UniqueKeyKindProperty:
		default:
			blog.Errorf("[CreateObjectUnique] invalid key kind: %s", key.Kind)
			return 0, ctx.Error.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
		}
	}

	if inputParam.Data.MustCheck {
		cond := condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(objID)
		cond.Field("must_check").Eq(true)
		count, err := m.dbProxy.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).Count(ctx)
		if nil != err {
			blog.Errorf("[CreateObjectUnique] check must check error: %#v", err)
			return 0, ctx.Error.Error(common.CCErrObjectDBOpErrno)
		}
		if count > 0 {
			blog.Errorf("[CreateObjectUnique] model could not have multiple must check unique")
			return 0, ctx.Error.Error(common.CCErrTopoObjectUniqueCanNotHasMutiMustCheck)
		}
	}

	err := m.recheckUniqueForExistsInsts(ctx, objID, inputParam.Data.Keys, inputParam.Data.MustCheck)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] recheckUniqueForExistsInsts for %s with %#v error: %#v", objID, inputParam, err)
		return 0, ctx.Error.Errorf(common.CCErrCommDuplicateItem, "")
	}

	id, err := m.dbProxy.NextSequence(ctx, common.BKTableNameObjUnique)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] NextSequence error: %#v", err)
		return 0, ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}

	unique := metadata.ObjectUnique{
		ID:        id,
		ObjID:     objID,
		MustCheck: inputParam.Data.MustCheck,
		Keys:      inputParam.Data.Keys,
		Ispre:     false,
		OwnerID:   ctx.SupplierAccount,
		LastTime:  metadata.Now(),
	}
	_, err = inputParam.Data.Metadata.Label.GetBusinessID()
	if nil == err {
		unique.Metadata = inputParam.Data.Metadata
	}
	err = m.dbProxy.Table(common.BKTableNameObjUnique).Insert(ctx, &unique)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] Insert error: %#v, raw: %#v", err, &unique)
		return 0, ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}

	return id, nil
}

func (m *modelAttrUnique) updateModelAttrUnique(ctx core.ContextParams, objID string, id uint64, data metadata.UpdateModelAttrUnique) error {

	unique := data.Data
	unique.LastTime = metadata.Now()

	for _, key := range unique.Keys {
		switch key.Kind {
		case metadata.UniqueKeyKindProperty:
		default:
			blog.Errorf("[UpdateObjectUnique] invalid key kind: %s", key.Kind)
			return ctx.Error.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
		}
	}

	if unique.MustCheck {
		cond := condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(objID)
		cond.Field("must_check").Eq(true)
		cond.Field("id").NotEq(id)
		count, err := m.dbProxy.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).Count(ctx)
		if nil != err {
			blog.Errorf("[UpdateObjectUnique] check must check  error: %#v", err)
			return ctx.Error.Error(common.CCErrObjectDBOpErrno)
		}
		if count > 0 {
			blog.Errorf("[UpdateObjectUnique] model could not have multiple must check unique")
			return ctx.Error.Error(common.CCErrTopoObjectUniqueCanNotHasMutiMustCheck)
		}
	}

	err := m.recheckUniqueForExistsInsts(ctx, objID, unique.Keys, unique.MustCheck)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] recheckUniqueForExistsInsts for %s with %#v error: %#v", objID, unique, err)
		return ctx.Error.Errorf(common.CCErrCommDuplicateItem, "")
	}

	cond := condition.CreateCondition()
	cond.Field("id").Eq(id)
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ctx.SupplierAccount)

	oldunique := metadata.ObjectUnique{}
	err = m.dbProxy.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).One(ctx, &oldunique)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] find error: %s, raw: %#v", err, cond.ToMapStr())
		return ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}

	if oldunique.Ispre {
		blog.Errorf("[UpdateObjectUnique] could not update preset constrain: %s", err, oldunique)
		return ctx.Error.Error(common.CCErrTopoObjectUniquePresetCouldNotDelOrEdit)
	}

	err = m.dbProxy.Table(common.BKTableNameObjUnique).Update(ctx, cond.ToMapStr(), &unique)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] Update error: %s, raw: %#v", err, &unique)
		return ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}
	return nil
}

func (m *modelAttrUnique) deleteModelAttrUnique(ctx core.ContextParams, objID string, id uint64) error {
	cond := condition.CreateCondition()
	cond.Field("id").Eq(id)
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ctx.SupplierAccount)

	unique := metadata.ObjectUnique{}
	err := m.dbProxy.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).One(ctx, &unique)
	if nil != err {
		blog.Errorf("[DeleteObjectUnique] find error: %s, raw: %#v", err, cond.ToMapStr())
		return ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}

	if unique.Ispre {
		blog.Errorf("[DeleteObjectUnique] could not delete preset constrain: %s", err, unique)
		return ctx.Error.Error(common.CCErrTopoObjectUniquePresetCouldNotDelOrEdit)
	}

	err = m.dbProxy.Table(common.BKTableNameObjUnique).Delete(ctx, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[DeleteObjectUnique] Delete error: %s, raw: %#v", err, cond.ToMapStr())
		return ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}

	return nil
}

func (m *modelAttrUnique) recheckUniqueForExistsInsts(ctx core.ContextParams, objID string, keys []metadata.UniqueKey, mustCheck bool) error {
	propertyIDs := []uint64{}
	for _, key := range keys {
		switch key.Kind {
		case metadata.UniqueKeyKindProperty:
			propertyIDs = append(propertyIDs, key.ID)
		default:
			return ctx.Error.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
		}
	}

	propertys := []metadata.Attribute{}
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ctx.SupplierAccount)
	cond.Field(common.BKFieldID).In(propertyIDs)
	err := m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(ctx, &propertys)
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
	err = m.dbProxy.Table(common.GetInstTableName(objID)).AggregateOne(ctx, pipeline, &result)
	if err != nil && !m.dbProxy.IsNotFoundError(err) {
		blog.ErrorJSON("[ObjectUnique] recheckUniqueForExistsInsts failed %s, pipeline: %s", err, pipeline)
		return err
	}

	if result.Finded > 0 {
		return dal.ErrDuplicated
	}

	return nil
}
