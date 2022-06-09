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

// validateServiceTemplateAttrs validate service template attributes
func (p *processOperation) validateServiceTemplateAttrs(kit *rest.Kit, bizID int64, serviceTemplateID int64,
	attrs []metadata.SvcTempAttr) errors.CCErrorCoder {

	// validate create options
	if bizID == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)
	}
	if serviceTemplateID == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKServiceTemplateIDField)
	}
	if len(attrs) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "attributes")
	}

	// validate service template
	svcTempFilter := mapstr.MapStr{common.BKAppIDField: bizID, common.BKFieldID: serviceTemplateID}
	svcTempCnt, err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(svcTempFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count service template failed, cond: %+v, err: %v, rid: %s", svcTempFilter, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if svcTempCnt == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
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
	if err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(filter).All(kit.Ctx, &attributes); err != nil {
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
