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

package instances

import (
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

// validCreateUnique  valid create inst data unique
func (valid *validator) validCreateUnique(ctx core.ContextParams, instanceData mapstr.MapStr, instMedataData metadata.Metadata, instanceManager *instanceManager) error {
	uniqueAttr, err := valid.dependent.SearchUnique(ctx, valid.objID)
	if nil != err {
		blog.Errorf("[validCreateUnique] search [%s] unique error %v, rid: %s", valid.objID, err, ctx.ReqID)
		return err
	}

	if 0 >= len(uniqueAttr) {
		blog.Warnf("[validCreateUnique] there're not unique constraint for %s, return, rid: %s", valid.objID, ctx.ReqID)
		return nil
	}

	for _, unique := range uniqueAttr {
		// retrieve unique value
		uniqueKeys := make([]string, 0)
		for _, key := range unique.Keys {
			switch key.Kind {
			case metadata.UniqueKeyKindProperty:
				property, ok := valid.idToProperty[int64(key.ID)]
				if !ok {
					blog.Errorf("[validCreateUnique] find [%s] property [%d] error %v, rid: %s", valid.objID, key.ID, ctx.ReqID)
					return valid.errif.Errorf(common.CCErrTopoObjectPropertyNotFound, key.ID)
				}
				uniqueKeys = append(uniqueKeys, property.PropertyID)
			default:
				blog.Errorf("[validCreateUnique] find [%s] property [%d] unique kind invalid [%d], rid: %s", valid.objID, key.ID, key.Kind, ctx.ReqID)
				return valid.errif.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
			}
		}

		cond := mongo.NewCondition()

		anyEmpty := false
		for _, key := range uniqueKeys {
			val, ok := instanceData[key]
			if !ok || isEmpty(val) {
				anyEmpty = true
			}
			cond.Element(&mongo.Eq{Key: key, Val: val})
		}

		if anyEmpty && !unique.MustCheck {
			continue
		}

		// only search data not in disable status
		cond.Element(&mongo.Neq{Key: common.BKDataStatusField, Val: common.DataStatusDisabled})
		if common.GetObjByType(valid.objID) == common.BKInnerObjIDObject {
			cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: valid.objID})
		}

		isExist, bizID := instMedataData.Label.Get(common.BKAppIDField)
		if isExist {
			_, metaCond := cond.Embed(metadata.BKMetadata)
			_, labelCond := metaCond.Embed(metadata.BKLabel)
			labelCond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bizID})
		}

		result, err := instanceManager.countInstance(ctx, valid.objID, cond.ToMapStr())
		if nil != err {
			blog.Errorf("[validCreateUnique] count [%s] inst error %v, condition: %#v, rid: %s", valid.objID, err, cond.ToMapStr(), ctx.ReqID)
			return err
		}

		if 0 < result {
			blog.Errorf("[validCreateUnique] duplicate data condition: %#v, unique keys: %#v, objID %s, rid: %s", cond.ToMapStr(), uniqueKeys, valid.objID, ctx.ReqID)
			propertyNames := make([]string, 0)
			for _, key := range uniqueKeys {
				propertyNames = append(propertyNames, util.FirstNotEmptyString(ctx.Lang.Language(valid.objID+"_property_"+key), valid.propertys[key].PropertyName, key))
			}

			return valid.errif.Errorf(common.CCErrCommDuplicateItem, strings.Join(propertyNames, ","))
		}

	}

	return nil
}

// validUpdateUnique valid update unique
func (valid *validator) validUpdateUnique(ctx core.ContextParams, updateData mapstr.MapStr, instMedataData metadata.Metadata, instID uint64, instanceManager *instanceManager) error {
	uniqueAttr, err := valid.dependent.SearchUnique(ctx, valid.objID)
	if nil != err {
		blog.Errorf("[validUpdateUnique] search [%s] unique error %v, rid: %s", valid.objID, err, ctx.ReqID)
		return err
	}

	if 0 >= len(uniqueAttr) {
		blog.Warnf("[validUpdateUnique] there're not unique constraint for %s, return, rid: %s", valid.objID, ctx.ReqID)
		return nil
	}

	for _, unique := range uniqueAttr {
		// retrieve unique value
		uniqueKeys := make([]string, 0)
		for _, key := range unique.Keys {
			switch key.Kind {
			case metadata.UniqueKeyKindProperty:
				property, ok := valid.idToProperty[int64(key.ID)]
				if !ok {
					blog.Errorf("[validUpdateUnique] find [%s] property [%d] error %v, rid: %s", valid.objID, key.ID, ctx.ReqID)
					return valid.errif.Errorf(common.CCErrTopoObjectPropertyNotFound, property.ID)
				}
				uniqueKeys = append(uniqueKeys, property.PropertyID)
			default:
				blog.Errorf("[validUpdateUnique] find [%s] property [%d] unique kind invalid [%d], rid: %s", valid.objID, key.ID, key.Kind, ctx.ReqID)
				return valid.errif.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
			}
		}

		cond := mongo.NewCondition()
		anyEmpty := false
		for _, key := range uniqueKeys {
			val, ok := updateData[key]
			if !ok || isEmpty(val) {
				anyEmpty = true
			}
			cond.Element(&mongo.Eq{Key: key, Val: val})
		}

		if anyEmpty && !unique.MustCheck {
			continue
		}

		// only search data not in disable status
		cond.Element(&mongo.Neq{Key: common.BKDataStatusField, Val: common.DataStatusDisabled})
		if common.GetObjByType(valid.objID) == common.BKInnerObjIDObject {
			cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: valid.objID})
		}
		cond.Element(&mongo.Neq{Key: common.GetInstIDField(valid.objID), Val: instID})
		isExist, bizID := instMedataData.Label.Get(common.BKAppIDField)
		if isExist {
			_, metaCond := cond.Embed(metadata.BKMetadata)
			_, lableCond := metaCond.Embed(metadata.BKLabel)
			lableCond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bizID})
		}

		result, err := instanceManager.countInstance(ctx, valid.objID, cond.ToMapStr())
		if nil != err {
			blog.Errorf("[validUpdateUnique] count [%s] inst error %v, condition: %#v, rid: %s", valid.objID, err, cond.ToMapStr(), ctx.ReqID)
			return err
		}

		if 0 < result {
			blog.Errorf("[validUpdateUnique] duplicate data condition: %#v, unique keys: %#v, objID %s, rid: %s", cond.ToMapStr(), uniqueKeys, valid.objID, ctx.ReqID)
			propertyNames := make([]string, 0)
			for _, key := range uniqueKeys {
				propertyNames = append(propertyNames, util.FirstNotEmptyString(ctx.Lang.Language(valid.objID+"_property_"+key), valid.propertys[key].PropertyName, key))
			}

			return valid.errif.Errorf(common.CCErrCommDuplicateItem, strings.Join(propertyNames, ","))
		}
	}
	return nil
}
