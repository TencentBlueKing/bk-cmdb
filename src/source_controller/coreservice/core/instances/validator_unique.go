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
	common.BKHostInnerIPField:   true,
	common.BKHostOuterIPField:   true,
	common.BKOperatorField:      true,
	common.BKBakOperatorField:   true,
	common.BKHostInnerIPv6Field: true,
	common.BKHostOuterIPv6Field: true,
}

type validUniqueOption struct {
	Condition  universalsql.Condition
	UniqueKeys []string
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
		for _, key := range uniqueKeys {
			val := data[key]
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

		uniqueOpts = append(uniqueOpts, validUniqueOption{Condition: cond, UniqueKeys: uniqueKeys})
	}

	return uniqueOpts, nil
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
