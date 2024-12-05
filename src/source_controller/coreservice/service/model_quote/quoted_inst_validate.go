/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package modelquote

import (
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

func getQuoteAttributes(kit *rest.Kit, objID string) (string, map[string]metadata.Attribute, error) {
	quoteRelCond := mapstr.MapStr{common.BKDestModelField: objID}
	quoteRelation := new(metadata.ModelQuoteRelation)

	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModelQuoteRelation).Find(quoteRelCond).
		Fields(common.BKSrcModelField, common.BKPropertyIDField).One(kit.Ctx, &quoteRelation)
	if err != nil {
		blog.Errorf("get quoted relations failed, err: %v, dest object: %s, rid: %s", err, objID, kit.Rid)
		return "", nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	attrCond := mapstr.MapStr{
		common.BKObjIDField:      quoteRelation.SrcModel,
		common.BKPropertyIDField: quoteRelation.PropertyID,
	}
	attribute := new(metadata.Attribute)

	err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAttDes).Find(attrCond).One(kit.Ctx, attribute)
	if err != nil {
		blog.Errorf("get quote attribute failed, err: %v, cond: %+v, rid: %s", err, attrCond, kit.Rid)
		return "", nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	option, err := metadata.ParseTableAttrOption(attribute.Option)
	if err != nil {
		blog.Errorf("parse table attribute option failed, err: %v, attr: %+v, rid: %s", err, attribute, kit.Rid)
		return "", nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKOptionField)
	}

	attrMap := make(map[string]metadata.Attribute)
	for _, attr := range option.Header {
		attrMap[attr.PropertyID] = attr
	}

	return quoteRelation.SrcModel, attrMap, nil
}

func validateCreateQuotedInstances(kit *rest.Kit, objID string, instances []mapstr.MapStr) error {
	// get source model info and quote attributes
	srcObj, attrMap, err := getQuoteAttributes(kit, objID)
	if err != nil {
		return err
	}

	// validate instances
	srcInstIDs := make([]int64, 0)
	for _, instance := range instances {
		srcInstID, err := validateCreateQuotedInst(kit, objID, instance, attrMap)
		if err != nil {
			return err
		}

		if srcInstID != 0 {
			srcInstIDs = append(srcInstIDs, srcInstID)
		}
	}

	// validate source instance ids
	if len(srcInstIDs) == 0 {
		return nil
	}

	srcInstIDs = util.IntArrayUnique(srcInstIDs)
	srcTable := common.GetInstTableName(srcObj, kit.TenantID)
	srcCond := mapstr.MapStr{common.GetInstIDField(srcObj): mapstr.MapStr{common.BKDBIN: srcInstIDs}}

	cnt, err := mongodb.Shard(kit.ShardOpts()).Table(srcTable).Find(srcCond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count source inst failed, err: %v, obj: %s, cond: %+v, rid: %s", err, srcObj, srcCond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if int(cnt) != len(srcInstIDs) {
		blog.Errorf("not all source instances exists, count: %d, ids: %+v, rid: %s", cnt, srcInstIDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKInstIDField)
	}

	return nil
}

func validateCreateQuotedInst(kit *rest.Kit, objID string, instance mapstr.MapStr,
	attrMap map[string]metadata.Attribute) (int64, error) {

	var srcInstID int64

	for key, val := range instance {

		if key == common.BKInstIDField {
			instIDVal, err := util.GetInt64ByInterface(val)
			if err != nil {
				blog.Errorf("source instance id is not int64, err: %v, val: %v, rid: %s", err, val, kit.Rid)
				return 0, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKInstIDField)
			}

			srcInstID = instIDVal
			continue
		}

		attr, ok := attrMap[key]
		if !ok {
			delete(instance, key)
			continue
		}

		if str, ok := val.(string); ok {
			val = strings.TrimSpace(str)
			instance[key] = val
		}

		rawErr := attr.Validate(kit.Ctx, val, key)
		if rawErr.ErrCode != 0 {
			err := rawErr.ToCCError(kit.CCError)
			blog.Errorf("validate %s inst failed, err: %v, key: %s, val: %v, rid: %s", objID, err, key, val, kit.Rid)
			return 0, err
		}
	}

	for attrID, attribute := range attrMap {
		if _, exists := instance[attrID]; exists {
			continue
		}

		if attribute.IsRequired {
			blog.Errorf("%s inst attr %s is not set, rid: %s", objID, attrID, kit.Rid)
			return 0, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, attrID)
		}

		instance[attrID] = getLostFieldDefaultValue(kit, attribute)
	}

	if _, exists := instance[common.BKInstIDField]; !exists {
		instance[common.BKInstIDField] = 0
	}

	return srcInstID, nil
}

// getLostFieldDefaultValue fill lost field with zero value, right now quoted attribute default value is only for ui
func getLostFieldDefaultValue(kit *rest.Kit, attr metadata.Attribute) interface{} {
	switch attr.PropertyType {
	case common.FieldTypeSingleChar, common.FieldTypeLongChar:
		return ""
	case common.FieldTypeEnumMulti:
		return make([]interface{}, 0)
	case common.FieldTypeInt, common.FieldTypeFloat:
		return nil
	case common.FieldTypeBool:
		return false
	}

	return nil
}

func validateUpdateQuotedInst(kit *rest.Kit, objID string, instance mapstr.MapStr) error {
	_, attrMap, err := getQuoteAttributes(kit, objID)
	if err != nil {
		return err
	}

	for key, val := range instance {
		if key == common.BKInstIDField || key == common.BKFieldID {
			delete(instance, key)
			continue
		}

		attr, ok := attrMap[key]
		if !ok || !attr.IsEditable {
			delete(instance, key)
			continue
		}

		if str, ok := val.(string); ok {
			val = strings.TrimSpace(str)
			instance[key] = val
		}

		rawErr := attr.Validate(kit.Ctx, val, key)
		if rawErr.ErrCode != 0 {
			err := rawErr.ToCCError(kit.CCError)
			blog.Errorf("validate %s inst failed, err: %v, key: %s, val: %v, rid: %s", objID, err, key, val, kit.Rid)
			return err
		}
	}

	return nil
}
