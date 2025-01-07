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

package core

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// GetSortedColProp get sort column property
func (d *Client) GetSortedColProp(kit *rest.Kit, cond mapstr.MapStr) ([]ColProp, error) {
	colProps, err := d.GetObjColProp(kit, cond)
	if err != nil {
		blog.Errorf("get object column property failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	objID, err := cond.String(common.BKObjIDField)
	if err != nil {
		blog.Errorf("get objID from condition failed, instCond: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	filterPropID := getFilterPropID(objID)
	filterPropType := getFilterPropType()
	for idx := range colProps {
		if util.InStrArr(filterPropType, colProps[idx].PropertyType) {
			colProps = append(colProps[:idx], colProps[idx+1:]...)
			continue
		}

		if util.InStrArr(filterPropID, colProps[idx].ID) {
			colProps[idx].NotExport = true
		}
	}

	bizID, err := cond.Int64(common.BKAppIDField)
	if err != nil {
		bizID = 0
	}
	groups, err := d.getObjGroup(kit, objID, bizID)
	if err != nil {
		blog.Errorf("get object attribute group failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	colProps, err = sortColProp(colProps, groups)
	if err != nil {
		blog.Errorf("sort column property failed, column property: %v, attribute group: %v, err: %v, rid: %s",
			colProps, groups, err, kit.Rid)
		return nil, err
	}

	return colProps, nil
}

// getFilterPropID 不需要展示字段id
func getFilterPropID(objID string) []string {
	switch objID {
	case common.BKInnerObjIDHost:
		return []string{common.BKSetNameField, common.BKModuleNameField, common.BKAppNameField}
	default:
		return []string{common.CreateTimeField}
	}
}

// getFilterPropType 不需要展示字段类型
func getFilterPropType() []string {
	return []string{common.FieldTypeIDRule}
}

// GetObjColProp get object column property
func (d *Client) GetObjColProp(kit *rest.Kit, cond mapstr.MapStr) ([]ColProp, error) {
	attrs, err := d.ApiClient.ModelQuote().GetObjectAttrWithTable(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("get object fields failed, condition: %v, err: %v ,rid: %s", cond, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	result := make([]ColProp, len(attrs))
	for idx, attr := range attrs {
		colProp := ColProp{ID: attr.PropertyID, Name: attr.PropertyName, PropertyType: attr.PropertyType,
			IsRequire: attr.IsRequired, Option: attr.Option, Group: attr.PropertyGroup, RefSheet: attr.PropertyName,
			Length: PropertyNormalLen,
		}

		result[idx] = colProp
	}

	return result, nil
}

// getObjGroup get object property group
func (d *Client) getObjGroup(kit *rest.Kit, objID string, bizID int64) ([]metadata.AttributeGroup, error) {
	cond := mapstr.MapStr{
		common.BKObjIDField: objID,
		common.BKAppIDField: bizID,
		metadata.PageName: mapstr.MapStr{metadata.PageStart: 0, metadata.PageLimit: common.BKNoLimit,
			metadata.PageSort: common.BKPropertyGroupIndexField},
	}

	result, err := d.ApiClient.GetObjectGroup(kit.Ctx, kit.Header, objID, cond)
	if err != nil {
		blog.Errorf("get %s fields group failed, err:%+v, rid: %s", objID, err, kit.Rid)
		return nil, fmt.Errorf("get attribute group failed, err: %+v", err)
	}

	if !result.Result {
		blog.Errorf("get %s fields group result failed. code: %d, message: %s, rid: %s", objID, result.Code,
			result.ErrMsg, kit.Rid)

		return nil, fmt.Errorf("get attribute group result false, result: %+v", result)
	}

	return result.Data, nil
}

// GetObjectData get object data
func (d *Client) GetObjectData(kit *rest.Kit, objID string) ([]interface{}, error) {
	cond := &metadata.ExportObjectCondition{ObjIDs: []string{objID}}
	result, err := d.ApiClient.GetObjectData(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("ger object data failed, cond: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	if result.CCError() != nil {
		blog.Errorf("ger object data failed, cond: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, result.CCError()
	}

	attrs := result.Data[objID].Attr
	filterPropType := getFilterPropType()
	for idx, attr := range attrs {
		attrMap, ok := attr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("object attribute is invalid, val: %v", attr)
		}

		if util.InStrArr(filterPropType, util.GetStrByInterface(attrMap[common.BKPropertyTypeField])) {
			attrs = append(attrs[:idx], attrs[idx+1:]...)
		}
	}

	return attrs, nil
}

// AddObjectBatch batch add object
func (d *Client) AddObjectBatch(kit *rest.Kit, param map[string]interface{}) (*metadata.Response, error) {
	result, err := d.ApiClient.AddObjectBatch(kit.Ctx, kit.Header, param)
	if err != nil {
		return nil, err
	}

	return result, nil
}
