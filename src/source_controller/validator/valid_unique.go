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

package validator

import (
	"context"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// validCreateUnique  valid create unique
func (valid *Validator) validCreateUnique(valData map[string]interface{}) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(valid.objID)
	cond.Field(common.BKOwnerIDField).Eq(valid.ownerID)

	uniques := []metadata.ObjectUnique{}
	err := valid.db.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).All(context.Background(), &uniques)
	if nil != err {
		blog.Errorf("[validCreateUnique] search [%s] unique error %v", valid.objID, err)
		return err
	}

	for _, unique := range uniques {
		// retrieve unique value
		uniquekeys := map[string]bool{}
		for _, key := range unique.Keys {
			switch key.Kind {
			case metadata.UinqueKeyKindProperty:
				property, ok := valid.idToProperty[int64(key.ID)]
				if !ok {
					blog.Errorf("[validCreateUnique] find [%s] property [%d] error not found", valid.objID, key.ID)
					return valid.errif.Errorf(common.CCErrTopoObjectPropertyNotFound, key.ID)
				}
				uniquekeys[property.PropertyID] = true
			default:
				blog.Errorf("[validCreateUnique] find [%s] property [%d] unique kind invalid [%d]", valid.objID, key.ID, key.Kind)
				return valid.errif.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
			}
		}

		cond := condition.CreateCondition()

		anyEmpty := false
		for key := range uniquekeys {
			val, ok := valData[key]
			if !ok || isEmpty(val) {
				anyEmpty = true
			}
			cond.Field(key).Eq(val)
		}

		if anyEmpty && !unique.MustCheck {
			continue
		}

		// only search data not in disable status
		cond.Field(common.BKDataStatusField).NotEq(common.DataStatusDisabled)
		if common.GetObjByType(valid.objID) == common.BKInnerObjIDObject {
			cond.Field(common.BKObjIDField).Eq(valid.objID)
		}

		tName := common.GetInstTableName(valid.objID)
		cnt, err := valid.db.Table(tName).Find(cond.ToMapStr()).Count(valid.ctx)
		if nil != err {
			blog.Errorf("[validCreateUnique] search [%s] inst error %v", valid.objID, err.Error())
			return err
		}

		if 0 < cnt {
			blog.Errorf("[validCreateUnique] duplicate data condition: %#v, unique keys: %#v, objID %s", cond.ToMapStr(), uniquekeys, valid.objID)
			propertyNames := []string{}
			for key := range uniquekeys {
				propertyNames = append(propertyNames, util.FirstNotEmptyString(valid.defLang.Language(valid.objID+"_property_"+key), valid.propertys[key].PropertyName, key))
			}

			return valid.errif.Errorf(common.CCErrCommDuplicateItem, strings.Join(propertyNames, " + "))
		}

	}

	return nil
}

func isEmpty(value interface{}) bool {
	return value == nil || value == ""
}

// validUpdateUnique valid update unique
func (valid *Validator) validUpdateUnique(valData map[string]interface{}, originData map[string]interface{}) error {

	objID := valid.objID
	mapData := make(map[string]interface{})
	for key, val := range originData {
		mapData[key] = val
	}
	instID := mapData[common.GetInstIDField(objID)]
	// retrive isonly value
	for key, val := range valData {
		mapData[key] = val
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(valid.objID)
	cond.Field(common.BKOwnerIDField).Eq(valid.ownerID)

	uniques := []metadata.ObjectUnique{}
	err := valid.db.Table(common.BKTableNameObjUnique).Find(cond.ToMapStr()).All(context.Background(), &uniques)
	if nil != err {
		blog.Errorf("[validCreateUnique] search [%s] unique error %v", valid.objID, err)
		return err
	}

	if 0 >= len(uniques) {
		blog.Warnf("[validUpdateUnique] there're not unique constraint for %s, return", valid.objID)
		return nil
	}

	for _, unique := range uniques {
		// retrieve unique value
		uniquekeys := map[string]bool{}
		for _, key := range unique.Keys {
			switch key.Kind {
			case metadata.UinqueKeyKindProperty:
				property, ok := valid.idToProperty[int64(key.ID)]
				if !ok {
					blog.Errorf("[validUpdateUnique] find [%s] property [%d] error: not found", valid.objID, key.ID)
					return valid.errif.Errorf(common.CCErrTopoObjectPropertyNotFound, property.ID)
				}
				uniquekeys[property.PropertyID] = true
			default:
				blog.Errorf("[validUpdateUnique] find [%s] property [%d] unique kind invalid [%d]", valid.objID, key.ID, key.Kind)
				return valid.errif.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
			}
		}

		cond := condition.CreateCondition()
		anyEmpty := false
		for key := range uniquekeys {
			val, ok := mapData[key]
			if !ok || isEmpty(val) {
				anyEmpty = true
			}
			cond.Field(key).Eq(val)
		}

		if anyEmpty && !unique.MustCheck {
			continue
		}

		// only search data not in diable status
		cond.Field(common.BKDataStatusField).NotEq(common.DataStatusDisabled)
		if common.GetObjByType(valid.objID) == common.BKInnerObjIDObject {
			cond.Field(common.BKObjIDField).Eq(valid.objID)
		}
		cond.Field(common.GetInstIDField(objID)).NotEq(instID)

		tName := common.GetInstTableName(objID)
		cnt, err := valid.db.Table(tName).Find(cond.ToMapStr()).Count(valid.ctx)
		if nil != err {
			blog.Errorf("[validUpdateUnique] search [%s] inst error %v", valid.objID, err.Error())
			return err
		}

		if 0 < cnt {
			blog.ErrorJSON("[validUpdateUnique] duplicate data table: %s,  condition: %s, origin: %s, unique keys: %s, objID: %s, instID %s count %d", tName, cond.ToMapStr(), mapData, uniquekeys, valid.objID, instID, cnt)
			propertyNames := []string{}
			for key := range uniquekeys {
				propertyNames = append(propertyNames, util.FirstNotEmptyString(valid.defLang.Language(valid.objID+"_property_"+key), valid.propertys[key].PropertyName, key))
			}

			return valid.errif.Errorf(common.CCErrCommDuplicateItem, strings.Join(propertyNames, " + "))
		}
	}
	return nil
}
