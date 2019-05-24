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
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

type validator struct {
	errif         errors.DefaultCCErrorIf
	propertys     map[string]metadata.Attribute
	idToProperty  map[int64]metadata.Attribute
	propertyslice []metadata.Attribute
	require       map[string]bool
	requirefields []string
	dependent     OperationDependences
	objID         string
}

// Init init
func NewValidator(ctx core.ContextParams, dependent OperationDependences, objID string, bizID int64) (*validator, error) {
	valid := &validator{}
	valid.propertys = make(map[string]metadata.Attribute)
	valid.idToProperty = make(map[int64]metadata.Attribute)
	valid.propertyslice = make([]metadata.Attribute, 0)
	valid.require = make(map[string]bool)
	valid.requirefields = make([]string, 0)
	valid.errif = ctx.Error
	result, err := dependent.SelectObjectAttWithParams(ctx, objID, bizID)
	if nil != err {
		return valid, err
	}
	for _, attr := range result {
		if attr.PropertyID == common.BKChildStr || attr.PropertyID == common.BKParentStr {
			continue
		}
		valid.propertys[attr.PropertyID] = attr
		valid.idToProperty[attr.ID] = attr
		valid.propertyslice = append(valid.propertyslice, attr)
		if attr.IsRequired {
			valid.require[attr.PropertyID] = true
			valid.requirefields = append(valid.requirefields, attr.PropertyID)
		}
	}
	valid.objID = objID
	valid.dependent = dependent
	return valid, nil
}
