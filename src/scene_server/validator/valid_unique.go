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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// validCreateUnique  valid create unique
func (valid *ValidMap) validCreateUnique(valData map[string]interface{}) error {
	uniqueresp, err := valid.CoreAPI.ObjectController().Unique().Search(valid.ctx, valid.pheader, valid.objID)
	if nil != err {
		blog.Errorf("[validCreateUnique] search [%s] unique error %v", valid.objID, err)
		return err
	}
	if !uniqueresp.Result {
		blog.Errorf("[validCreateUnique] search [%s] unique error %v", valid.objID, uniqueresp.ErrMsg)
		return valid.errif.New(uniqueresp.Code, uniqueresp.ErrMsg)
	}

	if 0 >= len(uniqueresp.Data) {
		blog.Warnf("[validCreateUnique] there're not unique constraint for %s, return", valid.objID)
		return nil
	}

	for _, unique := range uniqueresp.Data {
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

		result, err := valid.CoreAPI.ObjectController().Instance().SearchObjects(valid.ctx, common.GetObjByType(valid.objID), valid.pheader, &metadata.QueryInput{Condition: cond.ToMapStr()})
		if nil != err {
			blog.Errorf("[validCreateUnique] search [%s] inst error %v", valid.objID, err)
			return err
		}
		if !result.Result {
			blog.Errorf("[validCreateUnique] search [%s] inst error %v", valid.objID, result.ErrMsg)
			return valid.errif.New(result.Code, result.ErrMsg)
		}

		if 0 < result.Data.Count {
			blog.Errorf("[validUpdateUnique] duplicate data condition: %#v, unique keys: %#v, objID %s", cond.ToMapStr(), uniquekeys, valid.objID)

			defLang := valid.Language.CreateDefaultCCLanguageIf(util.GetLanguage(valid.pheader))
			propertyNames := []string{}
			for key := range uniquekeys {
				propertyNames = append(propertyNames, util.FirstNotEmptyString(defLang.Language(valid.objID+"_property_"+key), valid.propertys[key].PropertyName, key))
			}

			return valid.errif.Errorf(common.CCErrCommDuplicateItem, strings.Join(propertyNames, ","))
		}

	}

	return nil
}

func isEmpty(value interface{}) bool {
	return value == nil || value == ""
}

// validUpdateUnique valid update unique
func (valid *ValidMap) validUpdateUnique(valData map[string]interface{}, instID int64) error {

	objID := valid.objID
	mapData, err := valid.getInstDataByID(instID)
	if nil != err {
		blog.Errorf("[validUpdateUnique] search [%s] inst error %v", valid.objID, err)
		return err
	}

	// retrive isonly value
	for key, val := range valData {
		mapData[key] = val
	}

	uniqueresp, err := valid.CoreAPI.ObjectController().Unique().Search(valid.ctx, valid.pheader, valid.objID)
	if nil != err {
		blog.Errorf("[validUpdateUnique] search [%s] unique error %v", valid.objID, err)
		return err
	}
	if !uniqueresp.Result {
		blog.Errorf("[validUpdateUnique] search [%s] unique error %v", valid.objID, uniqueresp.ErrMsg)
		return valid.errif.New(uniqueresp.Code, uniqueresp.ErrMsg)
	}

	if 0 >= len(uniqueresp.Data) {
		blog.Warnf("[validUpdateUnique] there're not unique constraint for %s, return", valid.objID)
		return nil
	}

	for _, unique := range uniqueresp.Data {
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

		result, err := valid.CoreAPI.ObjectController().Instance().SearchObjects(valid.ctx, common.GetObjByType(valid.objID), valid.pheader, &metadata.QueryInput{Condition: cond.ToMapStr()})
		if nil != err {
			blog.Errorf("[validUpdateUnique] search [%s] inst error %v", valid.objID, err)
			return err
		}
		if !result.Result {
			blog.Errorf("[validUpdateUnique] search [%s] inst error %v", valid.objID, result.ErrMsg)
			return valid.errif.New(result.Code, result.ErrMsg)
		}

		if 0 < result.Data.Count {
			blog.Errorf("[validUpdateUnique] duplicate data condition: %#v, origin: %#v, unique keys: %v, objID: %s, instID %v count %d", cond.ToMapStr(), mapData, uniquekeys, valid.objID, instID, result.Data.Count)
			defLang := valid.Language.CreateDefaultCCLanguageIf(util.GetLanguage(valid.pheader))
			propertyNames := []string{}
			for key := range uniquekeys {
				propertyNames = append(propertyNames, util.FirstNotEmptyString(defLang.Language(valid.objID+"_property_"+key), valid.propertys[key].PropertyName, key))
			}

			return valid.errif.Errorf(common.CCErrCommDuplicateItem, strings.Join(propertyNames, " + "))
		}
	}
	return nil
}

// getInstDataByID get inst data by id
func (valid *ValidMap) getInstDataByID(instID int64) (map[string]interface{}, error) {
	objID := valid.objID
	searchCond := make(map[string]interface{})

	searchCond[common.GetInstIDField(objID)] = instID
	if common.GetInstTableName(objID) == common.BKTableNameBaseInst {
		objID = common.BKInnerObjIDObject
		searchCond[common.BKObjIDField] = valid.objID
	}

	blog.V(4).Infof("[getInstDataByID] condition: %#v, objID %s ", searchCond, objID)
	result, err := valid.CoreAPI.ObjectController().Instance().SearchObjects(valid.ctx, objID, valid.pheader, &metadata.QueryInput{Condition: searchCond, Limit: common.BKNoLimit})
	if nil != err {
		return nil, err
	}
	if !result.Result {
		return nil, valid.errif.New(result.Code, result.ErrMsg)
	}
	if len(result.Data.Info) == 0 {
		return nil, valid.errif.Error(common.CCErrCommNotFound)
	}

	if len(result.Data.Info[0]) > 0 {
		return result.Data.Info[0], nil
	}

	return nil, valid.errif.Error(common.CCErrCommNotFound)
}
