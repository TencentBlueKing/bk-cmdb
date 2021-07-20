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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
)

var hostSpecialFieldMap = map[string]bool{
	common.BKHostInnerIPField: true,
	common.BKHostOuterIPField: true,
	common.BKOperatorField:    true,
	common.BKBakOperatorField: true,
}

type validUniqueOption struct {
	Condition  universalsql.Condition
	UniqueKeys []string
}

// validCreateUnique valid create inst data unique
func (valid *validator) validCreateUnique(kit *rest.Kit, instanceData mapstr.MapStr, instanceManager *instanceManager) error {
	uniqueOpts, err := valid.getValidUniqueOptions(kit, instanceData, instanceManager)
	if err != nil {
		blog.Errorf("[validCreateUnique] getValidUniqueOptions error %v, data: %#v, rid: %s", err, instanceData, kit.Rid)
		return err
	}

	for _, opt := range uniqueOpts {
		cond := opt.Condition
		// only search data not in disable status
		cond.Element(&mongo.Neq{Key: common.BKDataStatusField, Val: common.DataStatusDisabled})
		if common.GetObjByType(valid.objID) == common.BKInnerObjIDObject {
			cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: valid.objID})
		}

		if err := valid.validUniqueByCond(kit, instanceManager, cond.ToMapStr(), opt.UniqueKeys); err != nil {
			return err
		}
	}

	return nil
}

// validUpdateUnique valid update unique
func (valid *validator) validUpdateUnique(kit *rest.Kit, updateData mapstr.MapStr, instanceData mapstr.MapStr, instID int64, instanceManager *instanceManager) error {
	// we need the complete updated instance data, override the db's original instance with updataData
	// considering the following scene
	// when updating a module's name, the whole update data is just {"bk_module_name":"new_name"}
	// as we know, the module's unique key has 3 fields: bk_biz_id, bk_set_id and bk_module_name
	// we can't just validate the "bk_module_name", but validate "bk_biz_id - bk_set_id - bk_module_name" as a whole
	// we need know all of the three fields value so that we can validate if the module's name is duplicate
	for k, v := range updateData {
		instanceData[k] = v
	}
	uniqueOpts, err := valid.getValidUniqueOptions(kit, instanceData, instanceManager)
	if err != nil {
		blog.Errorf("[validCreateUnique] getValidUniqueOptions error %v, data: %#v, rid: %s", err, instanceData, kit.Rid)
		return err
	}

	for _, opt := range uniqueOpts {
		needCheck := false
		// only check the unique field which need update
		for _, key := range opt.UniqueKeys {
			if _, ok := updateData[key]; ok {
				needCheck = true
			}
		}
		if !needCheck {
			continue
		}

		cond := opt.Condition
		// only search data not in disable status
		cond.Element(&mongo.Neq{Key: common.BKDataStatusField, Val: common.DataStatusDisabled})
		if common.GetObjByType(valid.objID) == common.BKInnerObjIDObject {
			cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: valid.objID})
		}
		cond.Element(&mongo.Neq{Key: common.GetInstIDField(valid.objID), Val: instID})

		if err := valid.validUniqueByCond(kit, instanceManager, cond.ToMapStr(), opt.UniqueKeys); err != nil {
			return err
		}
	}

	return nil
}

// getValidUniqueOptions get unique option used for validating the instances
func (valid *validator) getValidUniqueOptions(kit *rest.Kit, data mapstr.MapStr, instanceManager *instanceManager) ([]validUniqueOption, error) {
	uniqueOpts := make([]validUniqueOption, 0)

	for _, unique := range valid.uniqueAttrs {
		// retrieve unique value
		uniqueKeys := make([]string, 0)
		for _, key := range unique.Keys {
			switch key.Kind {
			case metadata.UniqueKeyKindProperty:
				property, ok := valid.idToProperty[int64(key.ID)]
				if !ok {
					blog.Errorf("find [%s] property [%d] error %v, rid: %s", valid.objID, key.ID, kit.Rid)
					return nil, valid.errIf.Errorf(common.CCErrTopoObjectPropertyNotFound, property.ID)
				}
				uniqueKeys = append(uniqueKeys, property.PropertyID)
			default:
				blog.Errorf("find [%s] property [%d] unique kind invalid [%d], rid: %s", valid.objID, key.ID, key.Kind, kit.Rid)
				return nil, valid.errIf.Errorf(common.CCErrTopoObjectUniqueKeyKindInvalid, key.Kind)
			}
		}

		cond := mongo.NewCondition()
		anyEmpty := false
		for _, key := range uniqueKeys {
			val, ok := data[key]
			if !ok || isEmpty(val) {
				anyEmpty = true
			}
			if valid.objID == common.BKInnerObjIDHost && hostSpecialFieldMap[key] {
				valStr, _ := val.(string)
				valArr := strings.Split(valStr, ",")
				cond.Element(&mongo.KV{
					Key: key,
					Val: map[string]interface{}{
						common.BKDBIN: valArr,
					},
				})
			} else {
				cond.Element(&mongo.Eq{Key: key, Val: val})
			}
		}

		if anyEmpty && !unique.MustCheck {
			continue
		}
		uniqueOpts = append(uniqueOpts, validUniqueOption{Condition: cond, UniqueKeys: uniqueKeys})
	}

	return uniqueOpts, nil
}

// validUniqueByCond valid instance unique by condition
func (valid *validator) validUniqueByCond(kit *rest.Kit, instanceManager *instanceManager, cond mapstr.MapStr, uniqueKeys []string) error {
	result, err := instanceManager.countInstance(kit, valid.objID, cond)
	if nil != err {
		blog.Errorf("[validUniqueByCond] count [%s] inst error %v, condition: %#v, rid: %s", valid.objID, err, cond, kit.Rid)
		return err
	}

	if result > 0 {
		blog.Errorf("[validUniqueByCond] duplicate data condition: %#v, unique keys: %#v, objID %s, rid: %s", cond, uniqueKeys, valid.objID, kit.Rid)
		propertyNames := make([]string, 0)
		lang := util.GetLanguage(kit.Header)
		language := valid.language.CreateDefaultCCLanguageIf(lang)
		for _, key := range uniqueKeys {
			propertyNames = append(propertyNames, util.FirstNotEmptyString(language.Language(valid.objID+"_property_"+key), valid.properties[key].PropertyName, key))
		}
		return valid.errIf.Errorf(common.CCErrCommDuplicateItem, strings.Join(propertyNames, ","))
	}

	return nil
}

// hasUniqueFields judge if the update data has any unique fields
func (valid *validator) hasUniqueFields(updateData mapstr.MapStr, uniqueOpts []validUniqueOption) (has bool, uniqueFields []string) {
	for _, opt := range uniqueOpts {
		if len(opt.UniqueKeys) == 0 {
			continue
		}

		keyCnt := 0
		// opt.UniqueKeys may have many keys as a union unique key
		// eg: for module, opt.UniqueKeys is [bk_biz_id, bk_set_id, bk_module_name]
		for _, key := range opt.UniqueKeys {
			if _, ok := updateData[key]; !ok {
				continue
			}
			keyCnt++
		}
		if keyCnt == len(opt.UniqueKeys) {
			return true, opt.UniqueKeys
		}
	}

	return false, make([]string, 0)
}

// validUpdateUniqFieldInMulti validate if it update unique field in multiple records
func (valid *validator) validUpdateUniqFieldInMulti(kit *rest.Kit, updateData mapstr.MapStr, instanceManager *instanceManager) error {
	uniqueOpts, err := valid.getValidUniqueOptions(kit, updateData, instanceManager)
	if err != nil {
		blog.Errorf("validUpdateUniqFieldInMulti failed, getValidUniqueOptions error %v, updateData: %#v, rid: %s", err, updateData, kit.Rid)
		return err
	}

	hasUniqueField, uniqueFields := valid.hasUniqueFields(updateData, uniqueOpts)
	if hasUniqueField == false {
		return nil
	}

	propertyNames := make([]string, 0)
	lang := util.GetLanguage(kit.Header)
	language := valid.language.CreateDefaultCCLanguageIf(lang)
	for _, key := range uniqueFields {
		propertyNames = append(propertyNames, util.FirstNotEmptyString(language.Language(valid.objID+"_property_"+key), valid.properties[key].PropertyName, key))
	}
	return valid.errIf.Errorf(common.CCErrCommDuplicateItem, strings.Join(propertyNames, ","))
}
