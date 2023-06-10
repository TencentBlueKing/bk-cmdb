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

package auditlog

import (
	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

type fieldTmplAuditLog struct {
	audit
}

// NewFieldTmplAuditLog new field template audit log
func NewFieldTmplAuditLog(clientSet coreservice.CoreServiceClientInterface) *fieldTmplAuditLog {
	return &fieldTmplAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}

// GenerateFieldTmplAuditLog generate audit of field template, if data is nil, will auto get current data by id.
func (h *fieldTmplAuditLog) GenerateFieldTmplAuditLog(parameter *generateAuditCommonParameter, id int64,
	tmpl *metadata.FieldTemplate) (*metadata.AuditLog, error) {

	kit := parameter.kit

	if tmpl == nil {
		query := &metadata.CommonQueryOption{
			CommonFilterOption: metadata.CommonFilterOption{
				Filter: filtertools.GenAtomFilter(common.BKFieldID, filter.Equal, id),
			},
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
		}
		rsp, err := h.clientSet.FieldTemplate().ListFieldTemplate(kit.Ctx, kit.Header, query)
		if err != nil {
			blog.Errorf("failed to read field template, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		if len(rsp.Info) <= 0 {
			blog.Errorf("failed to read field template, id: %d, rid: %s", id, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID)
		}

		tmpl = &rsp.Info[0]
	}

	return &metadata.AuditLog{
		AuditType:       metadata.FieldTemplateType,
		ResourceType:    metadata.FieldTemplateRes,
		Action:          parameter.action,
		ResourceID:      tmpl.ID,
		ResourceName:    tmpl.Name,
		OperateFrom:     parameter.operateFrom,
		OperationDetail: &metadata.GenericOpDetail{Data: tmpl},
	}, nil
}

// GenerateFieldTmplAttrAuditLog generate audit of field template attributes,
// if data is nil, will auto get current data by ids.
func (h *fieldTmplAuditLog) GenerateFieldTmplAttrAuditLog(parameter *generateAuditCommonParameter, ids []int64,
	attrs []metadata.FieldTemplateAttr) ([]metadata.AuditLog, error) {

	kit := parameter.kit

	if len(attrs) == 0 {
		query := &metadata.CommonQueryOption{
			CommonFilterOption: metadata.CommonFilterOption{
				Filter: filtertools.GenAtomFilter(common.BKFieldID, filter.In, ids),
			},
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
		}
		rsp, err := h.clientSet.FieldTemplate().ListFieldTemplateAttr(kit.Ctx, kit.Header, query)
		if err != nil {
			blog.Errorf("failed to read field template attributes, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		if len(rsp.Info) <= 0 {
			blog.Errorf("failed to read field template attributes, ids: %v, rid: %s", ids, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "ids")
		}

		attrs = rsp.Info
	}

	audits := make([]metadata.AuditLog, len(attrs))
	for idx := range attrs {
		attrAudit := metadata.AuditLog{
			AuditType:       metadata.FieldTemplateType,
			ResourceType:    metadata.FieldTemplateAttrRes,
			Action:          parameter.action,
			ResourceID:      attrs[idx].ID,
			ResourceName:    attrs[idx].PropertyName,
			OperateFrom:     parameter.operateFrom,
			OperationDetail: &metadata.GenericOpDetail{Data: attrs[idx]},
		}

		audits[idx] = attrAudit
	}

	return audits, nil
}

// GenerateFieldTmplUniqueAuditLog generate audit of field template uniques,
// if data is nil, will auto get current data by ids.
func (h *fieldTmplAuditLog) GenerateFieldTmplUniqueAuditLog(parameter *generateAuditCommonParameter, ids []int64,
	uniques []metadata.FieldTemplateUnique) ([]metadata.AuditLog, error) {

	kit := parameter.kit

	if len(uniques) == 0 {
		query := &metadata.CommonQueryOption{
			CommonFilterOption: metadata.CommonFilterOption{
				Filter: filtertools.GenAtomFilter(common.BKFieldID, filter.In, ids),
			},
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
		}
		rsp, err := h.clientSet.FieldTemplate().ListFieldTemplateUnique(kit.Ctx, kit.Header, query)
		if err != nil {
			blog.Errorf("failed to read field template uniques, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		if len(rsp.Info) <= 0 {
			blog.Errorf("failed to read field template uniques, ids: %v, rid: %s", ids, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "ids")
		}

		uniques = rsp.Info
	}

	audits := make([]metadata.AuditLog, len(uniques))
	for idx := range uniques {
		uniqueAudit := metadata.AuditLog{
			AuditType:       metadata.FieldTemplateType,
			ResourceType:    metadata.FieldTemplateUniqueRes,
			Action:          parameter.action,
			ResourceID:      uniques[idx].ID,
			OperateFrom:     parameter.operateFrom,
			OperationDetail: &metadata.GenericOpDetail{Data: uniques[idx]},
		}

		audits[idx] = uniqueAudit
	}

	return audits, nil
}
