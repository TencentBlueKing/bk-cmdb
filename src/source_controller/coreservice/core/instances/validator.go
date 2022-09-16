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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
)

type validator struct {
	errIf         errors.DefaultCCErrorIf
	properties    map[string]metadata.Attribute
	idToProperty  map[int64]metadata.Attribute
	propertySlice []metadata.Attribute
	require       map[string]bool
	requireFields []string
	uniqueAttrs   []metadata.ObjectUnique
	dependent     OperationDependences
	objID         string
	language      language.CCLanguageIf
}

// NewValidator TODO
// Init init
func NewValidator(kit *rest.Kit, dependent OperationDependences, objID string, bizID int64, language language.CCLanguageIf) (*validator, error) {
	valid := &validator{}
	valid.properties = make(map[string]metadata.Attribute)
	valid.idToProperty = make(map[int64]metadata.Attribute)
	valid.propertySlice = make([]metadata.Attribute, 0)
	valid.require = make(map[string]bool)
	valid.requireFields = make([]string, 0)
	valid.uniqueAttrs = make([]metadata.ObjectUnique, 0)
	valid.objID = objID
	valid.errIf = kit.CCError
	valid.dependent = dependent
	valid.language = language

	result, err := dependent.SelectObjectAttWithParams(kit, objID, []int64{bizID})
	if nil != err {
		return nil, err
	}
	for _, attr := range result {
		valid.properties[attr.PropertyID] = attr
		valid.idToProperty[attr.ID] = attr
		valid.propertySlice = append(valid.propertySlice, attr)
		if attr.IsRequired {
			valid.require[attr.PropertyID] = true
			valid.requireFields = append(valid.requireFields, attr.PropertyID)
		}
	}

	uniqueAttrs, err := valid.dependent.SearchUnique(kit, valid.objID)
	if nil != err {
		return nil, err
	}

	// process model do not have the unique rules, so we ignore it's attribute's unique check
	if len(uniqueAttrs) == 0 && objID != common.BKProcessObjectName {
		blog.Errorf("[validUpdateUnique] there're not unique constraint for %s, return, rid: %s", valid.objID, kit.Rid)
		uniqueAttrs = make([]metadata.ObjectUnique, 0)
	}
	valid.uniqueAttrs = uniqueAttrs

	return valid, nil
}

// NewValidators init validators, returns mapping of biz id to corresponding validator
func NewValidators(kit *rest.Kit, dependent OperationDependences, objID string, bizIDs []int64,
	language language.CCLanguageIf) (map[int64]*validator, error) {

	uniqueAttrs, err := dependent.SearchUnique(kit, objID)
	if err != nil {
		return nil, err
	}

	// process model do not have the unique rules, so we ignore it's attribute's unique check
	if len(uniqueAttrs) == 0 && objID != common.BKProcessObjectName {
		blog.Errorf("there's no unique constraint for %s, rid: %s", objID, kit.Rid)
		uniqueAttrs = make([]metadata.ObjectUnique, 0)
	}

	attributes, err := dependent.SelectObjectAttributes(kit, objID, bizIDs)
	if err != nil {
		return nil, err
	}

	attrMap := make(map[int64][]metadata.Attribute)
	for _, attr := range attributes {
		attrMap[attr.BizID] = append(attrMap[attr.BizID], attr)
	}

	bizValidatorMap := make(map[int64]*validator)
	for _, bizID := range bizIDs {
		if bizValidatorMap[bizID] != nil {
			continue
		}

		validator := &validator{
			properties:    make(map[string]metadata.Attribute),
			idToProperty:  make(map[int64]metadata.Attribute),
			propertySlice: make([]metadata.Attribute, 0),
			require:       make(map[string]bool),
			requireFields: make([]string, 0),
			uniqueAttrs:   uniqueAttrs,
			objID:         objID,
			errIf:         kit.CCError,
			dependent:     dependent,
			language:      language,
		}

		// the instances in biz has both biz attributes and global attributes that has no biz id
		attrArr := attrMap[0]
		if bizID != 0 {
			attrArr = append(attrArr, attrMap[bizID]...)
		}
		for _, attr := range attrArr {
			validator.properties[attr.PropertyID] = attr
			validator.idToProperty[attr.ID] = attr
			validator.propertySlice = append(validator.propertySlice, attr)
			if attr.IsRequired {
				validator.require[attr.PropertyID] = true
				validator.requireFields = append(validator.requireFields, attr.PropertyID)
			}
		}

		bizValidatorMap[bizID] = validator
	}

	return bizValidatorMap, nil
}
