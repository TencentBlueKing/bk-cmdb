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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
)

func (m *modelAttrUnique) searchModelAttrUnique(kit *rest.Kit, inputParam metadata.QueryCondition) (results []metadata.ObjectUnique, err error) {
	results = []metadata.ObjectUnique{}
	instHandler := mongodb.Client().Table(common.BKTableNameObjUnique).Find(inputParam.Condition)
	err = instHandler.Start(uint64(inputParam.Page.Start)).Limit(uint64(inputParam.Page.Limit)).Sort(inputParam.Page.Sort).All(kit.Ctx, &results)

	return results, err
}

func (m *modelAttrUnique) countModelAttrUnique(kit *rest.Kit, cond mapstr.MapStr) (count uint64, err error) {

	count, err = mongodb.Client().Table(common.BKTableNameObjUnique).Find(cond).Count(kit.Ctx)

	return count, err
}

func (m *modelAttrUnique) createModelAttrUnique(kit *rest.Kit, objID string, inputParam metadata.CreateModelAttrUnique) (uint64, error) {
	for _, key := range inputParam.Data.Keys {
		switch key.Kind {
		case metadata.UniqueKeyKindProperty:
		default:
			blog.Errorf("[CreateObjectUnique] invalid key kind: %s, rid: %s", key.Kind, kit.Rid)
			return 0, kit.CCError.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
		}
	}

	if inputParam.Data.MustCheck {
		cond := condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(objID)
		cond.Field("must_check").Eq(true)
		count, err := mongodb.Client().Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).Count(kit.Ctx)
		if nil != err {
			blog.Errorf("[CreateObjectUnique] check must check error: %#v, rid: %s", err, kit.Rid)
			return 0, kit.CCError.Error(common.CCErrObjectDBOpErrno)
		}
		if count > 0 {
			blog.Errorf("[CreateObjectUnique] model could not have multiple must check unique, rid: %s", kit.Rid)
			return 0, kit.CCError.Error(common.CCErrTopoObjectUniqueCanNotHasMultipleMustCheck)
		}
	}

	exist, err := m.checkUniqueRuleExist(kit, objID, 0, inputParam.Data.Keys)
	if err != nil {
		blog.Errorf("[CreateObjectUnique] checkUniqueRuleExist error: %#v, rid: %s", err, kit.Rid)
		return 0, err
	}
	if exist {
		blog.Errorf("[CreateObjectUnique] same unique check rule has been exist: %#v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.Error(common.CCERrrCoreServiceSameUniqueCheckRuleExist)
	}

	properties, err := m.getUniqueProperties(kit, objID, inputParam.Data.Keys, inputParam.Data.MustCheck)
	if nil != err {
		blog.ErrorJSON("[CreateObjectUnique] getUniqueProperties for %s with %s err: %s, rid: %s", objID, inputParam, err, kit.Rid)
		return 0, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "keys")
	}

	err = m.recheckUniqueForExistsInstances(kit, objID, properties, inputParam.Data.MustCheck)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] recheckUniqueForExistsInsts for %s with %#v err: %#v, rid: %s", objID, inputParam, err, kit.Rid)
		return 0, kit.CCError.Errorf(common.CCErrCommDuplicateItem, "instance")
	}

	id, err := mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameObjUnique)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] NextSequence error: %#v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrObjectDBOpErrno)
	}

	unique := metadata.ObjectUnique{
		ID:        id,
		ObjID:     objID,
		MustCheck: inputParam.Data.MustCheck,
		Keys:      inputParam.Data.Keys,
		Ispre:     false,
		OwnerID:   kit.SupplierAccount,
		LastTime:  metadata.Now(),
	}
	err = mongodb.Client().Table(common.BKTableNameObjUnique).Insert(kit.Ctx, &unique)
	if nil != err {
		blog.Errorf("[CreateObjectUnique] Insert error: %#v, raw: %#v, rid: %s", err, &unique, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrObjectDBOpErrno)
	}

	return id, nil
}

func (m *modelAttrUnique) updateModelAttrUnique(kit *rest.Kit, objID string, id uint64, data metadata.UpdateModelAttrUnique) error {

	unique := data.Data
	unique.LastTime = metadata.Now()

	for _, key := range unique.Keys {
		switch key.Kind {
		case metadata.UniqueKeyKindProperty:
		default:
			blog.Errorf("[UpdateObjectUnique] invalid key kind: %s, rid: %s", key.Kind, kit.Rid)
			return kit.CCError.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
		}
	}

	if unique.MustCheck {
		cond := condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(objID)
		cond.Field("must_check").Eq(true)
		cond.Field("id").NotEq(id)
		count, err := mongodb.Client().Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).Count(kit.Ctx)
		if nil != err {
			blog.Errorf("[UpdateObjectUnique] check must check  error: %#v, rid: %s", err, kit.Rid)
			return kit.CCError.Error(common.CCErrObjectDBOpErrno)
		}
		if count > 0 {
			blog.Errorf("[UpdateObjectUnique] model could not have multiple must check unique, rid: %s", kit.Rid)
			return kit.CCError.Error(common.CCErrTopoObjectUniqueCanNotHasMultipleMustCheck)
		}
	}

	exist, err := m.checkUniqueRuleExist(kit, objID, id, unique.Keys)
	if err != nil {
		blog.Errorf("[UpdateObjectUnique] checkUniqueRuleExist error: %#v, rid: %s", err, kit.Rid)
		return err
	}
	if exist {
		blog.Errorf("[UpdateObjectUnique] same unique check rule has been exist: %#v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCERrrCoreServiceSameUniqueCheckRuleExist)
	}

	properties, err := m.getUniqueProperties(kit, objID, unique.Keys, unique.MustCheck)
	if nil != err {
		blog.ErrorJSON("[CreateObjectUnique] getUniqueProperties for %s with %s err: %s, rid: %s", objID, unique, err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "keys")
	}

	err = m.recheckUniqueForExistsInstances(kit, objID, properties, unique.MustCheck)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] recheckUniqueForExistsInsts for %s with %#v error: %#v, rid: %s", objID, unique, err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommDuplicateItem, "instance")
	}

	cond := condition.CreateCondition()
	cond.Field("id").Eq(id)
	cond.Field(common.BKObjIDField).Eq(objID)
	condMap := util.SetModOwner(cond.ToMapStr(), kit.SupplierAccount)

	oldUnique := metadata.ObjectUnique{}
	err = mongodb.Client().Table(common.BKTableNameObjUnique).Find(condMap).One(kit.Ctx, &oldUnique)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] find error: %s, raw: %#v, rid: %s", err, cond.ToMapStr(), kit.Rid)
		return kit.CCError.Error(common.CCErrObjectDBOpErrno)
	}

	if oldUnique.Ispre {
		blog.Errorf("[UpdateObjectUnique] could not update preset constrain: %+v %v, rid: %s", oldUnique, err, kit.Rid)
		return kit.CCError.Error(common.CCErrTopoObjectUniquePresetCouldNotDelOrEdit)
	}

	err = mongodb.Client().Table(common.BKTableNameObjUnique).Update(kit.Ctx, cond.ToMapStr(), &unique)
	if nil != err {
		blog.Errorf("[UpdateObjectUnique] Update error: %s, raw: %#v, rid: %s", err, &unique, kit.Rid)
		return kit.CCError.Error(common.CCErrObjectDBOpErrno)
	}
	return nil
}

func (m *modelAttrUnique) deleteModelAttrUnique(kit *rest.Kit, objID string, id uint64) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKFieldID).Eq(id)
	cond.Field(common.BKObjIDField).Eq(objID)
	condMap := util.SetModOwner(cond.ToMapStr(), kit.SupplierAccount)

	unique := metadata.ObjectUnique{}
	err := mongodb.Client().Table(common.BKTableNameObjUnique).Find(condMap).One(kit.Ctx, &unique)
	if nil != err {
		blog.Errorf("[DeleteObjectUnique] find error: %s, raw: %#v, rid: %s", err, cond.ToMapStr(), kit.Rid)
		return kit.CCError.Error(common.CCErrObjectDBOpErrno)
	}

	if unique.Ispre {
		blog.Errorf("[DeleteObjectUnique] could not delete preset constrain: %+v, %v, rid: %s", unique, err, kit.Rid)
		return kit.CCError.Error(common.CCErrTopoObjectUniquePresetCouldNotDelOrEdit)
	}

	exist, err := m.checkUniqueRequireExist(kit, objID, []uint64{id})
	if err != nil {
		blog.ErrorJSON("deleteModelAttrUnique check unique require err:%s, cond:%s, rid:%s", err.Error(), cond.ToMapStr(), kit.Rid)
		return err
	}
	if !exist {
		blog.ErrorJSON("deleteModelAttrUnique check unique require result. not found other require unique, cond:%s, rid:%s", cond.ToMapStr(), kit.Rid)
		return kit.CCError.CCError(common.CCErrTopoObjectUniqueShouldHaveMoreThanOne)
	}

	fCond := cond.ToMapStr()
	err = mongodb.Client().Table(common.BKTableNameObjUnique).Delete(kit.Ctx, fCond)
	if nil != err {
		blog.Errorf("[DeleteObjectUnique] Delete error: %s, raw: %#v, rid: %s", err, fCond, kit.Rid)
		return kit.CCError.Error(common.CCErrObjectDBOpErrno)
	}

	return nil
}

// get properties via keys
func (m *modelAttrUnique) getUniqueProperties(kit *rest.Kit, objID string, keys []metadata.UniqueKey, mustCheck bool) ([]metadata.Attribute, error) {
	propertyIDs := make([]int64, 0)
	for _, key := range keys {
		propertyIDs = append(propertyIDs, int64(key.ID))
	}
	propertyIDs = util.IntArrayUnique(propertyIDs)

	properties := make([]metadata.Attribute, 0)
	attCond := condition.CreateCondition()
	attCond.Field(common.BKObjIDField).Eq(objID)
	attCond.Field(common.BKFieldID).In(propertyIDs)
	fCond := attCond.ToMapStr()
	fCond = util.SetQueryOwner(fCond, kit.SupplierAccount)

	err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(fCond).All(kit.Ctx, &properties)
	if err != nil {
		blog.ErrorJSON("[ObjectUnique] getUniqueProperties find properties for %s failed %s: %s, rid: %s", objID, err, kit.Rid)
		return nil, err
	}

	if len(properties) <= 0 {
		blog.ErrorJSON("[ObjectUnique] getUniqueProperties keys empty for [%s] %+s, rid: %s", objID, keys, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "keys")
	}
	if len(properties) != len(propertyIDs) {
		blog.ErrorJSON("[ObjectUnique] getUniqueProperties keys have non-existent attribute for [%s] %+s, rid: %s", objID, keys, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "keys")
	}

	return properties, nil
}

// for create or update a model instance unique check usage.
// the must_check is true, must be check exactly, no matter the check filed is empty or not.
// the must_check is false, only when all the filed is not empty, then it's check exactly, otherwise, skip this check.
func (m *modelAttrUnique) recheckUniqueForExistsInstances(kit *rest.Kit, objID string, properties []metadata.Attribute, mustCheck bool) error {
	// now, set the pipeline.
	pipeline := make([]interface{}, 0)

	instCond := mapstr.MapStr{}
	if common.GetObjByType(objID) == common.BKInnerObjIDObject {
		instCond.Set(common.BKObjIDField, objID)
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
	err := mongodb.Client().Table(common.GetInstTableName(objID)).AggregateOne(kit.Ctx, pipeline, &result)
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		blog.ErrorJSON("[ObjectUnique] recheckUniqueForExistsInsts failed %s, pipeline: %s, rid: %s", err, pipeline, kit.Rid)
		return err
	}

	if result.UniqueCount > 0 {
		return types.ErrDuplicated
	}

	return nil
}

// checkUniqueRequireExist  check if either is a required unique check
// ignoreUniqueIDS 除ignoreUniqueIDS之外是否有唯一校验项目
func (m *modelAttrUnique) checkUniqueRequireExist(kit *rest.Kit, objID string, ignoreUnqiqueIDS []uint64) (bool, error) {
	cond := condition.CreateCondition()
	if len(ignoreUnqiqueIDS) > 0 {
		cond.Field(common.BKFieldID).NotIn(ignoreUnqiqueIDS)
	}
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field("must_check").Eq(true)

	cnt, err := mongodb.Client().Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).Count(kit.Ctx)
	if nil != err {
		blog.ErrorJSON("[checkUniqueRequireExist] find error: %s, raw: %s, rid: %s", err, cond.ToMapStr(), kit.Rid)
		return false, kit.CCError.Error(common.CCErrObjectDBOpErrno)
	}
	if cnt > 0 {
		return true, nil
	}

	return false, nil
}

// checkUniqueRuleExist check if same unique rule has already existed
// if ruleID is 0,then it's create operation, otherwise it's update operation
func (m *modelAttrUnique) checkUniqueRuleExist(kit *rest.Kit, objID string, ruleID uint64, keys []metadata.UniqueKey) (bool, error) {
	// get all exist uniques
	uniqueCond := condition.CreateCondition()
	uniqueCond.Field(common.BKObjIDField).Eq(objID)
	cond := util.SetQueryOwner(uniqueCond.ToMapStr(), kit.SupplierAccount)
	existUniques := make([]metadata.ObjectUnique, 0)
	err := mongodb.Client().Table(common.BKTableNameObjUnique).Find(cond).All(kit.Ctx, &existUniques)
	if err != nil {
		return false, kit.CCError.Error(common.CCErrObjectDBOpErrno)
	}

	// compare to see if the input keys has already existed
	keysMap := make(map[uint64]bool)
	for _, key := range keys {
		keysMap[key.ID] = true
	}
	for _, u := range existUniques {
		if len(keysMap) == len(u.Keys) {
			cnt := 0
			for _, key := range u.Keys {
				if keysMap[key.ID] {
					cnt++
				}
			}
			if cnt == len(keysMap) && ruleID != u.ID {
				return true, nil
			}
		}
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
	case common.FieldTypeList:
		return nil, nil
	case common.FieldTypeOrganization:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", propertyType)
	}

}
