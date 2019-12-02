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
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

func (m *modelAttrUnique) searchModelAttrUnique(ctx core.ContextParams, inputParam metadata.QueryCondition) (results []metadata.ObjectUnique, err error) {
	results = []metadata.ObjectUnique{}
	instHandler := m.dbProxy.Table(common.BKTableNameObjUnique).Find(inputParam.Condition)
	for _, sort := range inputParam.SortArr {
		field := sort.Field
		if sort.IsDsc {
			field = "-" + field
		}
		instHandler = instHandler.Sort(field)
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
			blog.Errorf("[CreateObjectUnique] invalid key kind: %s, rid: %s", key.Kind, ctx.ReqID)
			return 0, ctx.Error.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
		}
	}

	if inputParam.Data.MustCheck {
		cond := condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(objID)
		cond.Field("must_check").Eq(true)
		count, err := m.dbProxy.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).Count(ctx)
		if nil != err {
			blog.Errorf("[CreateObjectUnique] check must check error: %#v, rid: %s", err, ctx.ReqID)
			return 0, ctx.Error.Error(common.CCErrObjectDBOpErrno)
		}
		if count > 0 {
			blog.Errorf("[CreateObjectUnique] model could not have multiple must check unique, rid: %s", ctx.ReqID)
			return 0, ctx.Error.Error(common.CCErrTopoObjectUniqueCanNotHasMultipleMustCheck)
		}
	}

	err := m.recheckUniqueForExistsInstances(ctx, objID, inputParam.Data.Keys, inputParam.Data.MustCheck, inputParam.Data.Metadata)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] recheckUniqueForExistsInsts for %s with %#v err: %#v, rid: %s", objID, inputParam, err, ctx.ReqID)
		return 0, ctx.Error.Errorf(common.CCErrCommDuplicateItem, "instance")
	}

	id, err := m.dbProxy.NextSequence(ctx, common.BKTableNameObjUnique)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] NextSequence error: %#v, rid: %s", err, ctx.ReqID)
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
		blog.Errorf("[CreateObjectUnique] Insert error: %#v, raw: %#v, rid: %s", err, &unique, ctx.ReqID)
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
			blog.Errorf("[UpdateObjectUnique] invalid key kind: %s, rid: %s", key.Kind, ctx.ReqID)
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
			blog.Errorf("[UpdateObjectUnique] check must check  error: %#v, rid: %s", err, ctx.ReqID)
			return ctx.Error.Error(common.CCErrObjectDBOpErrno)
		}
		if count > 0 {
			blog.Errorf("[UpdateObjectUnique] model could not have multiple must check unique, rid: %s", ctx.ReqID)
			return ctx.Error.Error(common.CCErrTopoObjectUniqueCanNotHasMultipleMustCheck)
		}
	}

	err := m.recheckUniqueForExistsInstances(ctx, objID, unique.Keys, unique.MustCheck, unique.Metadata)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] recheckUniqueForExistsInsts for %s with %#v error: %#v, rid: %s", objID, unique, err, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommDuplicateItem, "instance")
	}

	cond := condition.CreateCondition()
	cond.Field("id").Eq(id)
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ctx.SupplierAccount)
	if len(unique.Metadata.Label) > 0 {
		cond.Field(metadata.BKMetadata).Eq(unique.Metadata)
	}

	oldUnique := metadata.ObjectUnique{}
	err = m.dbProxy.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).One(ctx, &oldUnique)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] find error: %s, raw: %#v, rid: %s", err, cond.ToMapStr(), ctx.ReqID)
		return ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}

	if oldUnique.Ispre {
		blog.Errorf("[UpdateObjectUnique] could not update preset constrain: %+v %v, rid: %s", oldUnique, err, ctx.ReqID)
		return ctx.Error.Error(common.CCErrTopoObjectUniquePresetCouldNotDelOrEdit)
	}

	err = m.dbProxy.Table(common.BKTableNameObjUnique).Update(ctx, cond.ToMapStr(), &unique)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] Update error: %s, raw: %#v, rid: %s", err, &unique, ctx.ReqID)
		return ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}
	return nil
}

func (m *modelAttrUnique) deleteModelAttrUnique(ctx core.ContextParams, objID string, id uint64, meta metadata.DeleteModelAttrUnique) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKFieldID).Eq(id)
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKOwnerIDField).Eq(ctx.SupplierAccount)

	unique := metadata.ObjectUnique{}
	err := m.dbProxy.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).One(ctx, &unique)
	if nil != err {
		blog.Errorf("[DeleteObjectUnique] find error: %s, raw: %#v, rid: %s", err, cond.ToMapStr(), ctx.ReqID)
		return ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}

	if unique.Ispre {
		blog.Errorf("[DeleteObjectUnique] could not delete preset constrain: %+v, %v, rid: %s", unique, err, ctx.ReqID)
		return ctx.Error.Error(common.CCErrTopoObjectUniquePresetCouldNotDelOrEdit)
	}

	exist, err := m.checkUniqueRequireExist(ctx, objID, []uint64{id})
	if err != nil {
		blog.ErrorJSON("deleteModelAttrUnique check unique require err:%s, cond:%s, rid:%s", err.Error(), cond.ToMapStr(), ctx.ReqID)
		return err
	}
	if !exist {
		blog.ErrorJSON("deleteModelAttrUnique check unique require result. not found other require unique, cond:%s, rid:%s", cond.ToMapStr(), ctx.ReqID)
		return ctx.Error.CCError(common.CCErrTopoObjectUniqueShouldHaveMoreThanOne)
	}

	fCond := cond.ToMapStr()
	if len(meta.Label) > 0 {
		fCond.Merge(metadata.PublicAndBizCondition(meta.Metadata))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	err = m.dbProxy.Table(common.BKTableNameObjUnique).Delete(ctx, fCond)
	if nil != err {
		blog.Errorf("[DeleteObjectUnique] Delete error: %s, raw: %#v, rid: %s", err, fCond, ctx.ReqID)
		return ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}

	return nil
}

// for create or update a model instance unique check usage.
// the must_check is true, must be check exactly, no matter the check filed is empty or not.
// the must_check is false, only when all the filed is not empty, then it's check exactly, otherwise, skip this check.
func (m *modelAttrUnique) recheckUniqueForExistsInstances(ctx core.ContextParams, objID string, keys []metadata.UniqueKey, mustCheck bool, meta metadata.Metadata) error {
	propertyIDs := make([]uint64, 0)
	for _, key := range keys {
		switch key.Kind {
		case metadata.UniqueKeyKindProperty:
			propertyIDs = append(propertyIDs, key.ID)
		default:
			return ctx.Error.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
		}
	}

	properties := make([]metadata.Attribute, 0)
	attCond := condition.CreateCondition()
	attCond.Field(common.BKObjIDField).Eq(objID)
	attCond.Field(common.BKOwnerIDField).Eq(ctx.SupplierAccount)
	attCond.Field(common.BKFieldID).In(propertyIDs)
	fCond := attCond.ToMapStr()
	if len(meta.Label) > 0 {
		fCond.Merge(metadata.PublicAndBizCondition(meta))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	err := m.dbProxy.Table(common.BKTableNameObjAttDes).Find(fCond).All(ctx, &properties)
	if err != nil {
		blog.ErrorJSON("[ObjectUnique] recheckUniqueForExistsInsts find properties for %s failed %s: %s, rid: %s", objID, err, ctx.ReqID)
		return err
	}

	// now, set the pipeline.
	pipeline := make([]interface{}, 0)

	instCond := mapstr.MapStr{}
	if len(meta.Label) > 0 {
		instCond.Merge(metadata.PublicAndBizCondition(meta))
		instCond.Remove(metadata.BKMetadata)
	} else {
		instCond.Merge(metadata.BizLabelNotExist)
	}
	if common.GetObjByType(objID) == common.BKInnerObjIDObject {
		instCond.Set(common.BKObjIDField, objID)
	}

	if len(properties) <= 0 {
		blog.Warnf("[ObjectUnique] recheckUniqueForExistsInsts keys empty for [%s] %+v, rid: %s", objID, keys, ctx.ReqID)
		return nil
	}

	// if a unique is not a "must check", then it has two scenarios:
	// 1. if all the object's instance's this key is all empty, then it's acceptable, which means that
	//    it matched the unique rules.
	// 2. if one of the object's instance's unique key is set, then unique rules must be check. only when all
	//    the unique rules is matched, then it's acceptable.
	if !mustCheck {
		for _, property := range properties {
			basic, err := getBasicDataType(property.PropertyType)
			if err != nil {
				return err
			}
			// exclude fields that are null(not exist) and "ZERO" value.
			exclude := []interface{}{nil, basic}
			instCond.Set(property.PropertyID, mapstr.MapStr{common.BKDBExists: true, common.BKDBNIN: exclude})
		}
	}

	pipeline = append(pipeline, mapstr.MapStr{common.BKDBMatch: instCond})

	group := mapstr.MapStr{}
	for _, property := range properties {
		group.Set(property.PropertyID, "$"+property.PropertyID)
	}
	pipeline = append(pipeline, mapstr.MapStr{
		common.BKDBGroup: mapstr.MapStr{
			"_id":   group,
			"total": mapstr.MapStr{common.BKDBSum: 1},
		},
	})

	pipeline = append(pipeline, mapstr.MapStr{common.BKDBMatch: mapstr.MapStr{
		// "_id":   mapstr.MapStr{common.BKDBNE: nil},
		"total": mapstr.MapStr{common.BKDBGT: 1},
	}})

	pipeline = append(pipeline, mapstr.MapStr{common.BKDBCount: "unique_count"})

	js, _ := json.Marshal(pipeline)
	fmt.Println("pipeline: ", string(js))

	result := struct {
		UniqueCount uint64 `bson:"unique_count"`
	}{}
	err = m.dbProxy.Table(common.GetInstTableName(objID)).AggregateOne(ctx, pipeline, &result)
	if err != nil && !m.dbProxy.IsNotFoundError(err) {
		blog.ErrorJSON("[ObjectUnique] recheckUniqueForExistsInsts failed %s, pipeline: %s, rid: %s", err, pipeline, ctx.ReqID)
		return err
	}

	if result.UniqueCount > 0 {
		return dal.ErrDuplicated
	}

	return nil
}

// checkUniqueRequireExist  check if either is a required unique check
// ignoreUnqiqueIDS 除ignoreUnqiqueIDS之外是否有唯一校验项目
func (m *modelAttrUnique) checkUniqueRequireExist(ctx core.ContextParams, objID string, ignoreUnqiqueIDS []uint64) (bool, error) {
	cond := condition.CreateCondition()
	if len(ignoreUnqiqueIDS) > 0 {
		cond.Field(common.BKFieldID).NotIn(ignoreUnqiqueIDS)
	}
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field("must_check").Eq(true)

	cnt, err := m.dbProxy.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).Count(ctx)
	if nil != err {
		blog.ErrorJSON("[checkUniqueRequireExist] find error: %s, raw: %s, rid: %s", err, cond.ToMapStr(), ctx.ReqID)
		return false, ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}
	if cnt > 0 {
		return true, nil
	}

	return false, nil
}

func getBasicDataType(propertyType string) (interface{}, error) {
	switch propertyType {
	case common.FieldTypeSingleChar:
		return "", nil
	case common.FieldTypeLongChar:
		return "", nil
	case common.FieldTypeInt:
		return 0, nil
	case common.FieldTypeEnum:
		return "", nil
	case common.FieldTypeDate:
		return "", nil
	case common.FieldTypeTime:
		return "", nil
	case common.FieldTypeTimeZone:
		return "", nil
	case common.FieldTypeBool:
		return false, nil
	case common.FieldTypeFloat:
		return 0.0, nil
	case common.FieldTypeUser:
		return "", nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", propertyType)
	}

}
