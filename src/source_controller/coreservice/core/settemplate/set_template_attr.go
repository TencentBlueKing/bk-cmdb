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

package settemplate

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// validateSetTemplate validate set template
func (p *setTemplateOperation) validateSetTemplate(kit *rest.Kit, bizID int64,
	setTemplateID int64) errors.CCErrorCoder {

	if bizID == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)
	}
	if setTemplateID == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKSetTemplateIDField)
	}

	filter := mapstr.MapStr{common.BKAppIDField: bizID, common.BKFieldID: setTemplateID}
	cnt, err := mongodb.Client().Table(common.BKTableNameSetTemplate).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count set template failed, cond: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if cnt == 0 {
		blog.Errorf("set template is invalid, cond: %+v, cnt: %d, err: %v, rid: %s", filter, cnt, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
	}

	return nil
}

// validateSetTemplateAttrs validate set template attributes
func (p *setTemplateOperation) validateSetTemplateAttrs(kit *rest.Kit, bizID int64, setTemplateID int64,
	attrs []metadata.SetTempAttr) errors.CCErrorCoder {

	// validate set template
	if err := p.validateSetTemplate(kit, bizID, setTemplateID); err != nil {
		return err
	}

	if len(attrs) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "attributes")
	}

	// get set attributes
	attrIDs := make([]int64, 0)
	for _, item := range attrs {
		attrIDs = append(attrIDs, item.AttributeID)
	}

	filter := map[string]interface{}{
		common.BKFieldID:                  map[string]interface{}{common.BKDBIN: attrIDs},
		common.BKObjIDField:               common.BKInnerObjIDSet,
		common.BKPropertyIDField:          map[string]interface{}{common.BKDBNIN: []string{common.BKSetNameField}},
		metadata.AttributeFieldIsEditable: true,
	}
	util.AddModelBizIDCondition(filter, bizID)

	attributes := make([]metadata.Attribute, 0)
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(filter).All(kit.Ctx, &attributes); err != nil {
		blog.Errorf("get set attribute failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	attrMap := make(map[int64]metadata.Attribute)
	for _, attr := range attributes {
		attrMap[attr.ID] = attr
	}

	// validate attribute values
	for index, attr := range attrs {
		attribute, exists := attrMap[attr.AttributeID]
		if !exists {
			blog.Errorf("set attribute %d not exists, rid: %s", attr.AttributeID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attributes")
		}

		rawError := attribute.Validate(kit.Ctx, attr.PropertyValue, common.BKPropertyValueField)
		if rawError.ErrCode != 0 {
			ccErr := rawError.ToCCError(kit.CCError)
			blog.Errorf("validate attribute value failed, attr: %+v, err: %v, rid: %s", attr, ccErr, kit.Rid)
			return ccErr
		}

		// transfer to time value, because the field of the set instance is a time type, it needs to be the same.
		if attribute.PropertyType != common.FieldTypeTime {
			continue
		}

		val, err := util.ConvToTime(attr.PropertyValue)
		if err != nil {
			blog.Errorf("can not transfer property value to time type, attr: %+v, val: %v, rid: %s", attr,
				attr.PropertyValue, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attributes")
		}

		attrs[index].PropertyValue = val
	}

	return nil
}

// validateSetTemplateAttrExist validate set template attribute exist
func (p *setTemplateOperation) validateSetTemplateAttrExist(kit *rest.Kit, bizID int64, setTemplateID int64,
	attrIDs []int64) errors.CCErrorCoder {

	if len(attrIDs) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "attributes")
	}

	filter := map[string]interface{}{
		common.BKAppIDField:         bizID,
		common.BKSetTemplateIDField: setTemplateID,
		common.BKAttributeIDField: map[string]interface{}{
			common.BKDBIN: attrIDs,
		},
	}

	count, err := mongodb.Client().Table(common.BKTableNameSetTemplateAttr).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count set template attribute failed, filter: %v, err: %v, rid: %s", filter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if count != uint64(len(attrIDs)) {
		blog.Errorf("can't find all set template attributes, filter: %+v, count: %d, rid: %s", filter, count, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attributes")
	}

	return nil
}

// CreateSetTempAttr create set template attribute
func (p *setTemplateOperation) CreateSetTempAttr(kit *rest.Kit, option *metadata.CreateSetTempAttrsOption) (
	[]uint64, errors.CCErrorCoder) {

	if err := p.validateSetTemplateAttrs(kit, option.BizID, option.ID, option.Attributes); err != nil {
		return nil, err
	}

	ids, err := mongodb.Client().NextSequences(kit.Ctx, common.BKTableNameSetTemplateAttr, len(option.Attributes))
	if err != nil {
		blog.Errorf("get set template attribute ids failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	now := time.Now()

	setTempAttrs := make([]metadata.SetTemplateAttr, len(option.Attributes))
	for idx, attr := range option.Attributes {
		setTempAttrs[idx] = metadata.SetTemplateAttr{
			ID:              int64(ids[idx]),
			BizID:           option.BizID,
			SetTemplateID:   option.ID,
			AttributeID:     attr.AttributeID,
			PropertyValue:   attr.PropertyValue,
			Creator:         kit.User,
			Modifier:        kit.User,
			CreateTime:      now,
			LastTime:        now,
			SupplierAccount: kit.SupplierAccount,
		}
	}

	if err := mongodb.Client().Table(common.BKTableNameSetTemplateAttr).Insert(kit.Ctx, setTempAttrs); err != nil {
		blog.Errorf("create set template attributes(%+v) failed, err: %v, rid: %s", setTempAttrs, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	return ids, nil
}

// UpdateSetTempAttr update set template attribute
func (p *setTemplateOperation) UpdateSetTempAttr(kit *rest.Kit,
	option *metadata.UpdateSetTempAttrOption) errors.CCErrorCoder {

	if err := p.validateSetTemplateAttrs(kit, option.BizID, option.ID, option.Attributes); err != nil {
		return err
	}

	// validate set template attribute is exist
	attrIDs := make([]int64, 0)
	for _, attr := range option.Attributes {
		attrIDs = append(attrIDs, attr.AttributeID)
	}

	if err := p.validateSetTemplateAttrExist(kit, option.BizID, option.ID, attrIDs); err != nil {
		return err
	}

	// update set template attribute value
	attrFilter := map[string]interface{}{
		common.BKAppIDField:         option.BizID,
		common.BKSetTemplateIDField: option.ID,
	}
	for _, attribute := range option.Attributes {
		attrFilter[common.BKAttributeIDField] = attribute.AttributeID
		updateData := map[string]interface{}{common.BKPropertyValueField: attribute.PropertyValue}
		err := mongodb.Client().Table(common.BKTableNameSetTemplateAttr).Update(kit.Ctx, attrFilter, updateData)
		if err != nil {
			blog.Errorf("update set template attribute failed, filter: %v, err: %v, rid: %s", attrFilter, err,
				kit.Rid)
			return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
		}
	}

	return nil
}

// DeleteSetTemplateAttribute delete set template attribute
func (p *setTemplateOperation) DeleteSetTemplateAttribute(kit *rest.Kit,
	option *metadata.DeleteSetTempAttrOption) errors.CCErrorCoder {

	if err := p.validateSetTemplateAttrExist(kit, option.BizID, option.ID, option.AttributeIDs); err != nil {
		return err
	}

	filter := map[string]interface{}{
		common.BKAppIDField:         option.BizID,
		common.BKSetTemplateIDField: option.ID,
		common.BKAttributeIDField: map[string]interface{}{
			common.BKDBIN: option.AttributeIDs,
		},
	}

	if err := mongodb.Client().Table(common.BKTableNameSetTemplateAttr).Delete(kit.Ctx, filter); err != nil {
		blog.Errorf("delete set template attribute failed, filter: %v, err: %v, rid: %s", filter, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

// ListSetTemplateAttribute list set template attribute
func (p *setTemplateOperation) ListSetTemplateAttribute(kit *rest.Kit, option *metadata.ListSetTempAttrOption) (
	*metadata.SetTempAttrData, errors.CCErrorCoder) {

	if err := p.validateSetTemplate(kit, option.BizID, option.ID); err != nil {
		return nil, err
	}

	filter := map[string]interface{}{
		common.BKAppIDField:         option.BizID,
		common.BKSetTemplateIDField: option.ID,
	}

	templateAttrs := make([]metadata.SetTemplateAttr, 0)
	err := mongodb.Client().Table(common.BKTableNameSetTemplateAttr).Find(filter).Fields(option.Fields...).
		All(kit.Ctx, &templateAttrs)
	if err != nil {
		blog.Errorf("find set template attribute failed, filter: %v, err: %v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &metadata.SetTempAttrData{
		Attributes: templateAttrs,
	}, nil
}
