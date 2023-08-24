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

package exporter

import (
	"errors"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// ExportParamI export excel instance parameter interface
type ExportParamI interface {
	// GetPropCond get condition for query property
	GetPropCond() (mapstr.MapStr, error)

	// HasInstCond has condition for query instance
	HasInstCond() bool

	// GetInstCond get condition for query instance
	GetInstCond() (mapstr.MapStr, error)

	// Validate validate parameter
	Validate(kit *rest.Kit, lang language.DefaultCCLanguageIf) error
}

// InstParam export instance parameter
type InstParam struct {
	ObjID string `json:"bk_obj_id"`

	// CustomFields 导出的实例字段
	CustomFields []string `json:"export_custom_fields"`

	// InstIDArr 指定需要导出的实例ID
	InstIDArr []int64 `json:"bk_inst_ids"`

	// AssociationCond 用来限定导出关联关系，map[bk_obj_id]object_unique_id 2021年05月17日
	AssociationCond map[string]int64 `json:"association_condition"`

	// ObjectUniqueID 用来限定当前操作对象导出数据的时候，需要使用的唯一校验关系，
	// 自关联的时候，规定左边对象使用到的唯一索引
	ObjectUniqueID int64 `json:"object_unique_id"`

	cursor *cursor
}

// GetPropCond get condition for query property
func (e *InstParam) GetPropCond() (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKObjIDField: e.ObjID,
		metadata.PageName:   mapstr.MapStr{metadata.PageStart: 0, metadata.PageLimit: common.BKNoLimit},
	}

	if len(e.CustomFields) > 0 {
		cond[common.BKPropertyIDField] = map[string]interface{}{common.BKDBIN: e.CustomFields}
	}

	return cond, nil
}

// HasInstCond has condition for query instance
func (e *InstParam) HasInstCond() bool {
	if e.cursor == nil {
		e.cursor = &cursor{}
	}

	return e.cursor.hasNext()
}

// GetInstCond get condition for query instance
func (e *InstParam) GetInstCond() (mapstr.MapStr, error) {
	fields := make([]string, 0)
	if len(e.CustomFields) > 0 {
		fields = append(fields, e.CustomFields...)
		fields = append(fields, common.BKInstIDField)
	}

	if len(e.InstIDArr) > common.BKInstMaxExportLimit {
		return nil, fmt.Errorf("inst id exceed max len: %d", common.BKInstMaxExportLimit)
	}

	e.cursor.setEnd()

	return mapstr.MapStr{
		metadata.DBQueryCondition: mapstr.MapStr{
			common.BKInstIDField: mapstr.MapStr{common.BKDBIN: e.InstIDArr},
			common.BKObjIDField:  e.ObjID,
		},
		metadata.DBFields: fields,
	}, nil
}

// Validate validate parameter
func (e *InstParam) Validate(kit *rest.Kit, lang language.DefaultCCLanguageIf) error {
	if len(e.InstIDArr) > common.BKInstMaxExportLimit {
		return fmt.Errorf("bk_inst_ids exceed max length: %d", common.BKInstMaxExportLimit)
	}

	return nil
}

// HostParam export host parameter
type HostParam struct {
	// 导出的主机字段
	CustomFields []string `json:"export_custom_fields"`

	// 指定需要导出的主机ID, 设置本参数后， ExportCond限定条件无效
	HostIDArr []int64 `json:"bk_host_ids"`

	// 需要导出主机业务id
	AppID int64 `json:"bk_biz_id"`

	// 导出主机查询参数,就是search host 主机参数
	ExportCond metadata.HostCommonSearch `json:"export_condition"`

	// 用来限定导出关联关系，map[bk_obj_id]object_unique_id 2021年05月17日
	AssociationCond map[string]int64 `json:"association_condition"`

	// 用来限定当前操作对象导出数据的时候，需要使用的唯一校验关系，
	// 自关联的时候，规定左边对象使用到的唯一索引
	ObjectUniqueID int64 `json:"object_unique_id"`

	cursor *cursor
}

// GetPropCond get condition for query property
func (e *HostParam) GetPropCond() (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKObjIDField: common.BKInnerObjIDHost,
		metadata.PageName:   mapstr.MapStr{metadata.PageStart: 0, metadata.PageLimit: common.BKNoLimit},
		common.BKAppIDField: e.AppID,
	}

	if len(e.CustomFields) > 0 {
		e.CustomFields = append(e.CustomFields, common.BKHostInnerIPField, common.BKCloudIDField)
		cond[common.BKPropertyIDField] = map[string]interface{}{common.BKDBIN: e.CustomFields}
	}

	return cond, nil
}

// GetInstCond get condition for query instance
func (e *HostParam) GetInstCond() (mapstr.MapStr, error) {
	if e.cursor == nil {
		return nil, errors.New("need to call the HasInstCond method first")
	}

	fields := make([]string, 0)
	if len(e.CustomFields) > 0 {
		fields = append(fields, e.CustomFields...)
		fields = append(fields, common.BKHostIDField)
	}

	result := mapstr.MapStr{common.BKAppIDField: e.AppID}
	if len(e.HostIDArr) != 0 {
		if len(e.HostIDArr) > common.BKInstMaxExportLimit {
			return nil, fmt.Errorf("host ids exceed max length: %d", common.BKInstMaxExportLimit)
		}
		e.cursor.setEnd()

		cond := make([]interface{}, 0)

		// add host condition
		hostCond := make([]interface{}, 0)
		hostCond = append(hostCond, mapstr.MapStr{common.Field: common.BKHostIDField, common.Operator: common.BKDBIN,
			common.Value: e.HostIDArr})
		cond = append(cond, mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDHost, condition.DBFields: fields,
			common.Condition: hostCond})

		// add topo condition
		objIDs := []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule}

		for _, objID := range objIDs {
			topoFields := []string{common.GetInstIDField(objID), common.GetInstNameField(objID), common.TopoModuleName}
			topoCond := mapstr.MapStr{common.BKObjIDField: objID, condition.DBFields: topoFields}
			cond = append(cond, topoCond)
		}

		result[metadata.DBQueryCondition] = cond
		return result, nil
	}

	for idx, hostCond := range e.ExportCond.Condition {
		if hostCond.ObjectID == common.BKInnerObjIDHost {
			e.ExportCond.Condition[idx].Fields = fields
		}
	}

	result[metadata.DBQueryCondition] = e.ExportCond.Condition
	result[metadata.PageName] = e.cursor.getPage()
	e.cursor.next()

	return result, nil
}

// HasInstCond has condition for query instance
func (e *HostParam) HasInstCond() bool {
	if e.cursor == nil {
		limit := e.ExportCond.Page.Limit
		if common.BKMaxExportLimit < e.ExportCond.Page.Limit {
			limit = common.BKMaxExportLimit
		}

		e.cursor = &cursor{
			start:  e.ExportCond.Page.Start,
			limit:  limit,
			maxIdx: e.ExportCond.Page.Start + e.ExportCond.Page.Limit - 1,
		}
	}

	return e.cursor.hasNext()
}

// Validate validate parameter
func (e *HostParam) Validate(kit *rest.Kit, lang language.DefaultCCLanguageIf) error {
	if e.ExportCond.Page.Limit <= 0 || e.ExportCond.Page.Limit > common.BKMaxOnceExportLimit {
		return fmt.Errorf(lang.Languagef("export_page_limit_err", common.BKMaxOnceExportLimit))
	}

	return nil
}

type cursor struct {
	start  int
	limit  int
	maxIdx int
}

func (c *cursor) getPage() *metadata.BasePage {
	if c.start > c.maxIdx {
		return nil
	}

	if c.start+c.limit > c.maxIdx {
		c.limit = c.maxIdx - c.start + 1
	}

	return &metadata.BasePage{
		Start: c.start,
		Limit: c.limit,
	}
}

func (c *cursor) hasNext() bool {
	return c.maxIdx >= c.start
}

func (c *cursor) next() {
	c.start += c.limit
}

func (c *cursor) setEnd() {
	c.start = c.maxIdx + 1
}
