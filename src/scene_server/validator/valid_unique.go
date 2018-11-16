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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
)

// validCreateUnique  valid create unique
func (valid *ValidMap) validCreateUnique(valData map[string]interface{}) error {
	uniqueresp, err := valid.CoreAPI.ObjectController().Unique().Search(valid.ctx, valid.pheader, valid.objID)
	if nil != err {
		return err
	}
	if !uniqueresp.Result {
		return valid.errif.New(uniqueresp.Code, uniqueresp.ErrMsg)
	}

	if 0 >= len(uniqueresp.Data) {
		blog.V(3).Infof("is only array is zero for %s", valid.objID)
		return nil
	}

	for _, unique := range uniqueresp.Data {
		// retrive unique value
		uniquekeys := map[string]bool{}
		for _, key := range unique.Keys {
			switch key.Kind {
			case metadata.UinqueKeyKindProperty:
				property, ok := valid.idToProperty[int64(key.ID)]
				if !ok {
					return valid.errif.Errorf(common.CCErrTopoObjectPropertyNotFound, property.ID)
				}
				uniquekeys[property.PropertyID] = true
			default:
				return valid.errif.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
			}
		}

		cond := condition.CreateCondition()
		allEmpty := true
		for key, val := range valData {
			if uniquekeys[key] {
				cond.Field(key).Eq(val)
				if isEmpty(val) {
					allEmpty = false
				}
			}
		}

		if allEmpty && !unique.MustCheck {
			return nil
		}

		// only search data not in diable status
		cond.Field(common.BKDataStatusField).NotEq(common.DataStatusDisabled)
		if common.GetObjByType(valid.objID) == common.BKInnerObjIDObject {
			cond.Field(common.BKObjIDField).Eq(valid.objID)
		}

		result, err := valid.CoreAPI.ObjectController().Instance().SearchObjects(valid.ctx, common.GetObjByType(valid.objID), valid.pheader, &metadata.QueryInput{Condition: cond.ToMapStr()})
		if nil != err {
			return err
		}
		if !result.Result {
			return valid.errif.New(result.Code, result.ErrMsg)
		}

		if 0 < result.Data.Count {
			blog.Errorf("[validUpdateUnique] duplicate data condition: %#v, unique keys: %#v, objID %s", cond.ToMapStr(), uniquekeys, valid.objID)
			return valid.errif.Error(common.CCErrCommDuplicateItem)
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
	searchCond := make(map[string]interface{})
	mapData, err := valid.getInstDataByID(instID)
	if nil != err {
		return err
	}

	// retrive isonly value
	for key, val := range valData {
		mapData[key] = val
	}

	uniqueresp, err := valid.CoreAPI.ObjectController().Unique().Search(valid.ctx, valid.pheader, valid.objID)
	if nil != err {
		return err
	}
	if !uniqueresp.Result {
		return valid.errif.New(uniqueresp.Code, uniqueresp.ErrMsg)
	}

	if 0 >= len(uniqueresp.Data) {
		blog.V(3).Infof("is only array is zero for %s", valid.objID)
		return nil
	}

	for _, unique := range uniqueresp.Data {
		// retrive unique value
		uniquekeys := map[string]bool{}
		for _, key := range unique.Keys {
			switch key.Kind {
			case metadata.UinqueKeyKindProperty:
				property, ok := valid.idToProperty[int64(key.ID)]
				if !ok {
					return valid.errif.Errorf(common.CCErrTopoObjectPropertyNotFound, property.ID)
				}
				uniquekeys[property.PropertyID] = true
			default:
				return valid.errif.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
			}
		}

		cond := condition.CreateCondition()
		allEmpty := true
		for key, val := range mapData {
			if uniquekeys[key] {
				cond.Field(key).Eq(val)
				if isEmpty(val) {
					allEmpty = false
				}
			}
		}

		if allEmpty && !unique.MustCheck {
			return nil
		}

		// only search data not in diable status
		cond.Field(common.BKDataStatusField).NotEq(common.DataStatusDisabled)
		if common.GetObjByType(valid.objID) == common.BKInnerObjIDObject {
			cond.Field(common.BKObjIDField).Eq(valid.objID)
		}
		cond.Field(common.GetInstIDField(objID)).NotEq(instID)

		result, err := valid.CoreAPI.ObjectController().Instance().SearchObjects(valid.ctx, common.GetObjByType(valid.objID), valid.pheader, &metadata.QueryInput{Condition: searchCond})
		if nil != err {
			return err
		}
		if !result.Result {
			return valid.errif.Error(result.Code)
		}

		if 0 < result.Data.Count {
			blog.Errorf("[validUpdateUnique] duplicate data condition: %#v, origin: %#v, unique keys: %v, objID: %s, instID %v count %d", searchCond, mapData, uniquekeys, valid.objID, instID, result.Data.Count)
			return valid.errif.Error(common.CCErrCommDuplicateItem)
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
		return nil, valid.errif.Error(result.Code)
	}
	if len(result.Data.Info) == 0 {
		return nil, nil
	}

	if len(result.Data.Info[0]) > 0 {
		return result.Data.Info[0], nil
	}

	return nil, valid.errif.Error(common.CCErrCommNotFound)
}
