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

package fieldtmpl

import (
	"reflect"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

var innerFieldMap = map[string]struct{}{
	common.CreatorField:    {},
	common.CreateTimeField: {},
	common.ModifierField:   {},
	common.LastTimeField:   {},
}

// CompareFieldTemplateAttr compare field template attributes with object attributes,
// comparing priority: template id > property id > property type > property name
// @param forUI: defines if all comparison detail is needed for ui display, or only returns basic info for backend sync
func (t *template) CompareFieldTemplateAttr(kit *rest.Kit, opt *metadata.CompareFieldTmplAttrOption, forUI bool) (
	*metadata.CompareFieldTmplAttrsRes, error) {

	// check compare options that is not related to object attribute
	objID, err := t.comparator.getObjIDAndValidate(kit, opt.ObjectID)
	if err != nil {
		return nil, err
	}

	compParams, err := t.comparator.preCheckAttr(kit, objID, opt.Attrs)
	if err != nil {
		return nil, err
	}

	// get object attributes
	objAttrOpt := &metadata.QueryCondition{
		Condition: make(mapstr.MapStr),
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Fields: []string{common.BKFieldID, common.BKTemplateID, common.BKPropertyIDField, common.BKPropertyNameField,
			common.BKPropertyTypeField, common.BKOptionField, common.BKIsMultipleField, metadata.AttributeFieldDefault,
			common.BKIsRequiredField, metadata.AttributeFieldIsEditable, metadata.AttributeFieldUnit,
			metadata.AttributeFieldPlaceHolder},
	}
	util.AddModelBizIDCondition(objAttrOpt.Condition, 0)

	objAttrRes, err := t.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, objAttrOpt)
	if err != nil {
		blog.Errorf("get object attributes failed, err: %v, opt: %+v, rid: %s", err, opt, kit.Rid)
		return nil, err
	}

	// if object has no attributes, add all field template attributes
	if len(objAttrRes.Info) == 0 {
		createRes := make([]metadata.CompareOneFieldTmplAttrRes, len(opt.Attrs))
		for idx, attr := range opt.Attrs {
			createRes[idx] = metadata.CompareOneFieldTmplAttrRes{
				Index:      idx,
				PropertyID: attr.PropertyID,
			}
		}
		return &metadata.CompareFieldTmplAttrsRes{Create: createRes}, nil
	}

	// cross-compare object attributes with template attributes
	if forUI {
		return t.comparator.compareAttrForUI(kit, compParams, objAttrRes.Info)
	}
	return t.comparator.compareAttrForBackend(kit, compParams, objAttrRes.Info)
}

type compAttrParams struct {
	tmplIDMap     map[int64]metadata.FieldTemplateAttr
	tmplPropIDMap map[string]metadata.FieldTemplateAttr
	// tmplNameMap field template attribute name to property id map, used for conflict check
	tmplNameMap map[string]string
	// tmplIndexMap field template attribute property id to original index map, used to trace back to original attr
	tmplIndexMap map[string]int
	// createTmplMap to be created(has no matching/conflict attr) field template attr property id map
	createTmplMap map[string]struct{}
}

// preCheckAttr pre-check and categorize field template attributes before compare with object attributes
func (c *comparator) preCheckAttr(kit *rest.Kit, objID string, attrs []metadata.FieldTemplateAttr) (*compAttrParams,
	error) {

	params := &compAttrParams{
		tmplIDMap:     make(map[int64]metadata.FieldTemplateAttr),
		tmplPropIDMap: make(map[string]metadata.FieldTemplateAttr),
		tmplNameMap:   make(map[string]string),
		tmplIndexMap:  make(map[string]int),
		createTmplMap: make(map[string]struct{}),
	}

	tmplPropertyIDs, tmplPropertyNames := make([]string, 0), make([]string, 0)

	// check if field template attributes collides with inner attributes, and categorize field template attributes
	for index, attr := range attrs {
		if rawErr := attr.Validate(); rawErr.ErrCode != 0 {
			return nil, rawErr.ToCCError(kit.CCError)
		}
		if strings.HasPrefix(attr.PropertyID, "bk") {
			blog.Errorf("template attribute(%s) has 'bk' prefix, rid: %s", attr.PropertyID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attributes")
		}

		if _, exists := innerFieldMap[attr.PropertyID]; exists {
			blog.Errorf("template attribute(%s) collides with inner ones, rid: %s", attr.PropertyID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attributes")
		}

		if _, exists := params.tmplPropIDMap[attr.PropertyID]; exists {
			blog.Errorf("template attribute(%s) duplicate, rid: %s", attr.PropertyID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attributes")
		}

		if _, exists := params.tmplPropIDMap[attr.PropertyName]; exists {
			blog.Errorf("template attribute name(%s) duplicate, rid: %s", attr.PropertyName, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attributes")
		}

		if attr.ID != 0 {
			params.tmplIDMap[attr.ID] = attr
		}

		params.tmplPropIDMap[attr.PropertyID] = attr
		params.tmplNameMap[attr.PropertyName] = attr.PropertyID
		params.tmplIndexMap[attr.PropertyID] = index
		tmplPropertyIDs = append(tmplPropertyIDs, attr.PropertyID)
		tmplPropertyNames = append(tmplPropertyNames, attr.PropertyName)
		params.createTmplMap[attr.PropertyID] = struct{}{}
	}

	// check if field template attributes collides with biz custom attributes
	bizAttrCond := mapstr.MapStr{
		common.BKObjIDField: objID,
		common.BKAppIDField: mapstr.MapStr{common.BKDBGT: 0},
		common.BKDBOR: []mapstr.MapStr{
			{common.BKPropertyIDField: mapstr.MapStr{common.BKDBIN: tmplPropertyIDs}},
			{common.BKPropertyNameField: mapstr.MapStr{common.BKDBIN: tmplPropertyNames}},
		},
	}

	bizAttrCnt, err := c.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameObjAttDes, []map[string]interface{}{bizAttrCond})
	if err != nil {
		return nil, err
	}

	if len(bizAttrCnt) <= 0 || bizAttrCnt[0] > 0 {
		blog.Errorf("template(%+v) collides with biz custom field, cnt: %v, rid: %s", attrs, bizAttrCnt, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "attributes")
	}

	return params, nil
}

func (c *comparator) compareAttrForUI(kit *rest.Kit, params *compAttrParams, attributes []metadata.Attribute) (
	*metadata.CompareFieldTmplAttrsRes, error) {

	res := new(metadata.CompareFieldTmplAttrsRes)

	// compare object attribute with its template
	noTmplAttr := make([]metadata.Attribute, 0)
	for idx := range attributes {
		attr := attributes[idx]
		// compare attribute without template later, because template id has maximum priority in comparison
		if attr.TemplateID == 0 {
			noTmplAttr = append(noTmplAttr, attr)
			continue
		}

		tmplAttr, exists := params.tmplIDMap[attr.TemplateID]
		if !exists {
			// attribute's template is not exist, check if the attribute conflicts with other templates
			if conflictTmplAttr, ex := params.tmplPropIDMap[attr.PropertyID]; ex {
				c.addConflictAttrInfo(kit, params, &attr, conflictTmplAttr.PropertyID, common.BKPropertyIDField, res)
				continue
			}

			if conflictPropertyID, ex := params.tmplNameMap[attr.PropertyName]; ex {
				c.addConflictAttrInfo(kit, params, &attr, conflictPropertyID, common.BKPropertyNameField, res)
				continue
			}

			res.Unchanged = append(res.Unchanged, attr)
			continue
		}

		c.removeMatchingAttrParams(params, &tmplAttr)

		// compare it with its template if not conflict, check if attr's id and type is the same with its template
		if attr.PropertyID != tmplAttr.PropertyID {
			c.addConflictAttrInfo(kit, params, &attr, tmplAttr.PropertyID, common.BKPropertyIDField, res)
			continue
		}

		if attr.PropertyType != tmplAttr.PropertyType {
			c.addConflictAttrInfo(kit, params, &attr, tmplAttr.PropertyID, common.BKPropertyTypeField, res)
			continue
		}

		c.addUpdateAttrInfoForUI(kit, params, &tmplAttr, &attr, res)
	}

	// compare object attributes without template
	for idx := range noTmplAttr {
		attr := noTmplAttr[idx]

		tmplAttr, exists := params.tmplPropIDMap[attr.PropertyID]
		if !exists {
			// attribute is not related to template, check if its name conflicts with all templates
			if conflictPropertyID, ex := params.tmplNameMap[attr.PropertyName]; ex {
				c.addConflictAttrInfo(kit, params, &attr, conflictPropertyID, common.BKPropertyNameField, res)
				continue
			}

			res.Unchanged = append(res.Unchanged, attr)
			continue
		}

		c.removeMatchingAttrParams(params, &tmplAttr)

		// compare it with its template if not conflict, check if attr's property type is the same with its template
		if attr.PropertyType != tmplAttr.PropertyType {
			c.addConflictAttrInfo(kit, params, &attr, tmplAttr.PropertyID, common.BKPropertyTypeField, res)
			continue
		}

		c.addUpdateAttrInfoForUI(kit, params, &tmplAttr, &attr, res)
	}

	// field template attribute with no matching object attributes should be created
	c.addCreateAttrInfo(params, res)

	return res, nil
}

func (c *comparator) compareAttrForBackend(kit *rest.Kit, params *compAttrParams, attributes []metadata.Attribute) (
	*metadata.CompareFieldTmplAttrsRes, error) {

	res := new(metadata.CompareFieldTmplAttrsRes)

	// compare object attribute with its template
	noTmplAttr := make([]metadata.Attribute, 0)
	for idx := range attributes {
		attr := attributes[idx]
		// compare attribute without template later, because template id has maximum priority in comparison
		if attr.TemplateID == 0 {
			noTmplAttr = append(noTmplAttr, attr)
			continue
		}

		tmplAttr, exists := params.tmplIDMap[attr.TemplateID]
		if !exists {
			// attribute's template is not exist, check if the attribute conflicts with other templates
			if conflictTmplAttr, ex := params.tmplPropIDMap[attr.PropertyID]; ex {
				return nil, kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateAttrConflict, attr.ID,
					common.BKPropertyIDField, conflictTmplAttr.PropertyID)
			}

			if conflictPropertyID, ex := params.tmplNameMap[attr.PropertyName]; ex {
				return nil, kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateAttrConflict, attr.ID,
					common.BKPropertyNameField, conflictPropertyID)
			}

			// attribute template is deleted, so we update attribute's template id to zero
			res.Update = append(res.Update, metadata.CompareOneFieldTmplAttrRes{
				Index:      params.tmplIndexMap[tmplAttr.PropertyID],
				PropertyID: tmplAttr.PropertyID,
				Data:       &metadata.Attribute{ID: attr.ID},
				UpdateData: mapstr.MapStr{common.BKTemplateID: 0},
			})
			continue
		}

		c.removeMatchingAttrParams(params, &tmplAttr)

		// compare it with its template if not conflict, check if attr's id and type is the same with its template
		if attr.PropertyID != tmplAttr.PropertyID {
			return nil, kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateAttrConflict, attr.ID,
				common.BKPropertyIDField, tmplAttr.PropertyID)
		}

		if attr.PropertyType != tmplAttr.PropertyType {
			return nil, kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateAttrConflict, attr.ID,
				common.BKPropertyTypeField, tmplAttr.PropertyID)
		}

		err := c.addUpdateAttrInfoForBackend(kit, params, &attr, &tmplAttr, res)
		if err != nil {
			return nil, err
		}
	}

	// compare object attributes without template
	for idx := range noTmplAttr {
		attr := noTmplAttr[idx]

		tmplAttr, exists := params.tmplPropIDMap[attr.PropertyID]
		if !exists {
			// attribute is not related to template, check if its name conflicts with all templates
			if conflictPropertyID, ex := params.tmplNameMap[attr.PropertyName]; ex {
				return nil, kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateAttrConflict, attr.ID,
					common.BKPropertyNameField, conflictPropertyID)
			}

			continue
		}

		c.removeMatchingAttrParams(params, &tmplAttr)

		// compare it with its template if not conflict, check if attr's property type is the same with its template
		if attr.PropertyType != tmplAttr.PropertyType {
			return nil, kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateAttrConflict, attr.ID,
				common.BKPropertyTypeField, tmplAttr.PropertyID)
		}

		err := c.addUpdateAttrInfoForBackend(kit, params, &attr, &tmplAttr, res)
		if err != nil {
			return nil, err
		}
	}

	// field template attribute with no matching object attributes should be created
	c.addCreateAttrInfo(params, res)

	return res, nil
}

// removeMatchingAttrParams remove template attr that has matching obj attr, because it can't be related to another one
func (c *comparator) removeMatchingAttrParams(params *compAttrParams, tmplAttr *metadata.FieldTemplateAttr) {
	delete(params.tmplIDMap, tmplAttr.ID)
	delete(params.tmplPropIDMap, tmplAttr.PropertyID)
	delete(params.createTmplMap, tmplAttr.PropertyID)
}

// addConflictAttrInfo add conflict attribute info into compare result
func (c *comparator) addConflictAttrInfo(kit *rest.Kit, params *compAttrParams, attr *metadata.Attribute,
	conflictTmplPropID string, field string, res *metadata.CompareFieldTmplAttrsRes) {

	delete(params.createTmplMap, conflictTmplPropID)

	res.Conflict = append(res.Conflict, metadata.CompareOneFieldTmplAttrRes{
		Index:      params.tmplIndexMap[conflictTmplPropID],
		PropertyID: conflictTmplPropID,
		Message: kit.CCError.CCErrorf(common.CCErrTopoFieldTemplateAttrConflict, attr.ID, field,
			conflictTmplPropID).Error(),
		Data: attr,
	})
}

// addUpdateAttrInfoForUI add update/unchanged attribute info into compare result after comparison
func (c *comparator) addUpdateAttrInfoForUI(kit *rest.Kit, params *compAttrParams, tmplAttr *metadata.FieldTemplateAttr,
	attr *metadata.Attribute, res *metadata.CompareFieldTmplAttrsRes) {

	updateData, err := c.compareUpdatedAttr(kit, tmplAttr, attr)
	if err != nil {
		c.addConflictAttrInfo(kit, params, attr, tmplAttr.PropertyID, common.BKTemplateID, res)
		return
	}

	// ui do not compare template id
	delete(updateData, common.BKTemplateID)

	if len(updateData) > 0 {
		res.Update = append(res.Update, metadata.CompareOneFieldTmplAttrRes{
			Index:      params.tmplIndexMap[tmplAttr.PropertyID],
			PropertyID: tmplAttr.PropertyID,
			Data:       attr,
			UpdateData: updateData,
		})
		return
	}

	res.Unchanged = append(res.Unchanged, *attr)
}

// addUpdateAttrInfoForBackend add update attribute info into compare result after comparison
func (c *comparator) addUpdateAttrInfoForBackend(kit *rest.Kit, params *compAttrParams, attr *metadata.Attribute,
	tmplAttr *metadata.FieldTemplateAttr, res *metadata.CompareFieldTmplAttrsRes) error {

	updateData, err := c.compareUpdatedAttr(kit, tmplAttr, attr)
	if err != nil {
		return err
	}

	if len(updateData) > 0 {
		res.Update = append(res.Update, metadata.CompareOneFieldTmplAttrRes{
			Index:      params.tmplIndexMap[tmplAttr.PropertyID],
			PropertyID: tmplAttr.PropertyID,
			Data:       attr,
			UpdateData: updateData,
		})
	}

	return nil
}

// compareUpdatedAttr compare if field template attribute's matching attr conflict attribute info into compare result
func (c *comparator) compareUpdatedAttr(kit *rest.Kit, tmplAttr *metadata.FieldTemplateAttr, attr *metadata.Attribute) (
	mapstr.MapStr, error) {

	updateData := make(mapstr.MapStr)

	if tmplAttr.ID != attr.TemplateID {
		if attr.TemplateID != 0 {
			blog.Errorf("template id mismatch, attribute: %+v, template: %+v, rid: %s", attr, tmplAttr, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKTemplateID)
		}

		updateData[common.BKTemplateID] = tmplAttr.ID
	}

	if tmplAttr.PropertyName != attr.PropertyName {
		updateData[common.BKPropertyNameField] = tmplAttr.PropertyName
	}

	if tmplAttr.Unit != attr.Unit {
		updateData[metadata.AttributeFieldUnit] = tmplAttr.Unit
	}

	if tmplAttr.Placeholder.Lock && tmplAttr.Placeholder.Value != attr.Placeholder {
		updateData[metadata.AttributeFieldPlaceHolder] = tmplAttr.Placeholder.Value
	}

	if tmplAttr.Editable.Lock && tmplAttr.Editable.Value != attr.IsEditable {
		updateData[metadata.AttributeFieldIsEditable] = tmplAttr.Editable.Value
	}

	if tmplAttr.Required.Lock && tmplAttr.Required.Value != attr.IsRequired {
		updateData[common.BKIsRequiredField] = tmplAttr.Required.Value
	}

	if !reflect.DeepEqual(tmplAttr.Option, attr.Option) {
		updateData[common.BKOptionField] = tmplAttr.Option
	}

	if !reflect.DeepEqual(tmplAttr.Default, attr.Default) {
		updateData[metadata.AttributeFieldDefault] = tmplAttr.Default
	}

	isMultiple := false
	if attr.IsMultiple != nil {
		isMultiple = *attr.IsMultiple
	}
	if isMultiple != tmplAttr.IsMultiple {
		updateData[common.BKIsMultipleField] = tmplAttr.IsMultiple
	}

	return updateData, nil
}

// addCreateAttrInfo add create attribute template info into compare result
func (c *comparator) addCreateAttrInfo(params *compAttrParams, res *metadata.CompareFieldTmplAttrsRes) {
	for propertyID := range params.createTmplMap {
		res.Create = append(res.Create, metadata.CompareOneFieldTmplAttrRes{
			Index:      params.tmplIndexMap[propertyID],
			PropertyID: propertyID,
		})
	}
}