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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

var updateIgnoreKeys = []string{
	common.BKOwnerIDField,
	common.BKDefaultField,
	common.BKInstParentStr,
	common.BKAppIDField,
	common.BKDataStatusField,
	common.BKSupplierIDField,
	common.BKInstIDField,
}

var createIgnoreKeys = []string{
	common.BKOwnerIDField,
	common.BKDefaultField,
	common.BKInstParentStr,
	common.BKAppIDField,
	common.BKSupplierIDField,
	common.BKInstIDField,
	common.BKDataStatusField,
}

func (m *instanceManager) validCreateInstanceData(ctx core.ContextParams, objID string, instanceData mapstr.MapStr) error {
	valid, err := NewValidator(ctx, m.dependent, objID)
	if nil != err {
		blog.Errorf("init validator failed %s", err.Error())
		return err
	}
	FillLostedFieldValue(instanceData, valid.propertyslice, valid.requirefields)
	for _, key := range valid.requirefields {
		if _, ok := instanceData[key]; !ok {
			blog.Errorf("params in need, valid %s, data: %+v", objID, instanceData)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
		}
	}
	var instMedataData metadata.Metadata
	instMedataData.Label = make(metadata.Label)
	for key, val := range instanceData {
		if metadata.BKMetadata == key {
			bizID := metadata.GetBusinessIDFromMeta(val)
			if "" != bizID {
				instMedataData.Label.Set(metadata.LabelBusinessID, metadata.GetBusinessIDFromMeta(val))
			}
			continue
		}
		if util.InStrArr(createIgnoreKeys, key) {
			// ignore the key field
			continue
		}
		property, ok := valid.propertys[key]
		if !ok {
			blog.Errorf("params is not valid, the key is %s", key)
			return valid.errif.Errorf(common.CCErrCommParamsIsInvalid, key)
		}
		fieldType := property.PropertyType
		switch fieldType {
		case common.FieldTypeSingleChar:
			err = valid.validChar(val, key)
		case common.FieldTypeLongChar:
			err = valid.validLongChar(val, key)
		case common.FieldTypeInt:
			err = valid.validInt(val, key)
		case common.FieldTypeEnum:
			err = valid.validEnum(val, key)
		case common.FieldTypeDate:
			err = valid.validDate(val, key)
		case common.FieldTypeTime:
			err = valid.validTime(val, key)
		case common.FieldTypeTimeZone:
			err = valid.validTimeZone(val, key)
		case common.FieldTypeBool:
			err = valid.validBool(val, key)
		case common.FieldTypeForeignKey:
			err = valid.validForeignKey(val, key)
		default:
			continue
		}
		if nil != err {
			return err
		}
	}
	return valid.validCreateUnique(ctx, instanceData, instMedataData, m)
}

func (m *instanceManager) validUpdateInstanceData(ctx core.ContextParams, objID string, instanceData mapstr.MapStr, instMetaData metadata.Metadata, instID uint64) error {
	valid, err := NewValidator(ctx, m.dependent, objID)
	if nil != err {
		blog.Errorf("init validator faile %s", err.Error())
		return err
	}

	for key, val := range instanceData {

		if util.InStrArr(updateIgnoreKeys, key) {
			// ignore the key field
			continue
		}

		property, ok := valid.propertys[key]
		if !ok {
			blog.Errorf("params is not valid, the key is %s", key)
			return valid.errif.Errorf(common.CCErrCommParamsIsInvalid, key)
		}
		fieldType := property.PropertyType
		switch fieldType {
		case common.FieldTypeSingleChar:
			err = valid.validChar(val, key)
		case common.FieldTypeLongChar:
			err = valid.validLongChar(val, key)
		case common.FieldTypeInt:
			err = valid.validInt(val, key)
		case common.FieldTypeEnum:
			err = valid.validEnum(val, key)
		case common.FieldTypeDate:
			err = valid.validDate(val, key)
		case common.FieldTypeTime:
			err = valid.validTime(val, key)
		case common.FieldTypeTimeZone:
			err = valid.validTimeZone(val, key)
		case common.FieldTypeBool:
			err = valid.validBool(val, key)
		case common.FieldTypeForeignKey:
			err = valid.validForeignKey(val, key)
		default:
			continue
		}
		if nil != err {
			return err
		}
	}
	return valid.validUpdateUnique(ctx, instanceData, instMetaData, instID, m)
}
