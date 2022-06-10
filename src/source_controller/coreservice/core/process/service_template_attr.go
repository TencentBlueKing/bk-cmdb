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

package process

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// validateServiceTemplate validate service template
func (p *processOperation) validateServiceTemplate(kit *rest.Kit, bizID int64,
	serviceTemplateID int64) errors.CCErrorCoder {

	if bizID == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)
	}
	if serviceTemplateID == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKServiceTemplateIDField)
	}

	svcTempFilter := mapstr.MapStr{common.BKAppIDField: bizID, common.BKFieldID: serviceTemplateID}
	svcTempCnt, err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(svcTempFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count service template failed, cond: %+v, err: %v, rid: %s", svcTempFilter, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if svcTempCnt == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}

	return nil
}

// validateServiceTemplateAttrs validate service template attributes
func (p *processOperation) validateServiceTemplateAttrs(kit *rest.Kit, bizID int64, serviceTemplateID int64,
	attrs []metadata.SvcTempAttr) errors.CCErrorCoder {

	// validate service template
	if err := p.validateServiceTemplate(kit, bizID, serviceTemplateID); err != nil {
		return err
	}

	if len(attrs) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "attributes")
	}

	// get module attributes
	attrIDs := make([]int64, 0)
	for _, item := range attrs {
		attrIDs = append(attrIDs, item.AttributeID)
	}

	filter := map[string]interface{}{
		common.BKFieldID:    map[string]interface{}{common.BKDBIN: attrIDs},
		common.BKObjIDField: common.BKInnerObjIDModule,
	}
	util.AddModelBizIDCondition(filter, bizID)

	attributes := make([]metadata.Attribute, 0)
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(filter).All(kit.Ctx, &attributes); err != nil {
		blog.Errorf("get module attribute failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	attrMap := make(map[int64]metadata.Attribute)
	for _, attr := range attributes {
		attrMap[attr.ID] = attr
	}

	// validate attribute values
	for _, attr := range attrs {
		attribute, exists := attrMap[attr.AttributeID]
		if !exists {
			blog.Errorf("module attribute %d not exists, rid: %s", attr.AttributeID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attributes")
		}

		rawError := attribute.Validate(kit.Ctx, attr.PropertyValue, common.BKPropertyValueField)
		if rawError.ErrCode != 0 {
			ccErr := rawError.ToCCError(kit.CCError)
			blog.Errorf("validate attribute value failed, attr: %+v, err: %v, rid: %s", attr, ccErr, kit.Rid)
			return ccErr
		}
	}

	return nil
}

// validateServiceTemplateAttrExist validate service template attribute exist
func (p *processOperation) validateServiceTemplateAttrExist(kit *rest.Kit, bizID int64, serviceTemplateID int64,
	attrIDs []int64) errors.CCErrorCoder {

	filter := map[string]interface{}{
		common.BKAppIDField:             bizID,
		common.BKServiceTemplateIDField: serviceTemplateID,
		common.BKAttributeIDField: map[string]interface{}{
			common.BKDBIN: attrIDs,
		},
	}

	count, err := mongodb.Client().Table(common.BKTableNameServiceTemplateAttr).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count service template attribute failed, filter: %v, err: %v, rid: %s", filter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if count != uint64(len(attrIDs)) {
		blog.Errorf("can not find all service template attributes, rid: %s", kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attributes")
	}

	return nil
}

// UpdateServTempAttr update service template attribute
func (p *processOperation) UpdateServTempAttr(kit *rest.Kit,
	option *metadata.UpdateServTempAttrOption) errors.CCErrorCoder {

	if err := p.validateServiceTemplateAttrs(kit, option.BizID, option.ID, option.Attributes); err != nil {
		return err
	}

	// validate service template attribute is exist
	attrIDs := make([]int64, 0)
	for _, attr := range option.Attributes {
		attrIDs = append(attrIDs, attr.AttributeID)
	}

	if err := p.validateServiceTemplateAttrExist(kit, option.BizID, option.ID, attrIDs); err != nil {
		return err
	}

	// update service template attribute value
	attrFilter := map[string]interface{}{
		common.BKAppIDField:             option.BizID,
		common.BKServiceTemplateIDField: option.ID,
	}
	for _, attribute := range option.Attributes {
		attrFilter[common.BKAttributeIDField] = attribute.AttributeID
		updateData := map[string]interface{}{common.BKPropertyValueField: attribute.PropertyValue}
		err := mongodb.Client().Table(common.BKTableNameServiceTemplateAttr).Update(kit.Ctx, attrFilter, updateData)
		if err != nil {
			blog.Errorf("update service template attribute failed, filter: %s, err: %v, rid: %s", attrFilter, err,
				kit.Rid)
			return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
		}
	}

	return nil
}

// DeleteServiceTemplateAttribute delete service template attribute
func (p *processOperation) DeleteServiceTemplateAttribute(kit *rest.Kit,
	option *metadata.DeleteServTempAttrOption) errors.CCErrorCoder {

	if err := p.validateServiceTemplateAttrExist(kit, option.BizID, option.ID, option.AttributeIDs); err != nil {
		return err
	}

	filter := map[string]interface{}{
		common.BKAppIDField:             option.BizID,
		common.BKServiceTemplateIDField: option.ID,
		common.BKAttributeIDField: map[string]interface{}{
			common.BKDBIN: option.AttributeIDs,
		},
	}

	if err := mongodb.Client().Table(common.BKTableNameServiceTemplateAttr).Delete(kit.Ctx, filter); err != nil {
		blog.Errorf("delete service template attribute failed, filter: %v, err: %v, rid: %s", filter, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

// ListServiceTemplateAttribute list service template attribute
func (p *processOperation) ListServiceTemplateAttribute(kit *rest.Kit, option *metadata.ListServTempAttrOption) (
	*metadata.ServTempAttrData, errors.CCErrorCoder) {

	if err := p.validateServiceTemplate(kit, option.BizID, option.ID); err != nil {
		return nil, err
	}

	filter := map[string]interface{}{
		common.BKAppIDField:             option.BizID,
		common.BKServiceTemplateIDField: option.ID,
	}

	templateAttrs := make([]metadata.ServiceTemplateAttr, 0)
	err := mongodb.Client().Table(common.BKTableNameServiceTemplateAttr).Find(filter).All(kit.Ctx, &templateAttrs)
	if err != nil {
		blog.Errorf("find service template attribute failed, filter: %v, err: %v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &metadata.ServTempAttrData{
		Attributes: templateAttrs,
	}, nil
}
