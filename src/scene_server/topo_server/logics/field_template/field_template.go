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

// Package fieldtmpl defines field template logics.
package fieldtmpl

import (
	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/logics/model"
)

// FieldTemplateOperation field template operation methods
type FieldTemplateOperation interface {
	CreateFieldTemplate(kit *rest.Kit, opt *metadata.CreateFieldTmplOption) (*metadata.RspID, error)
	CompareFieldTemplateAttr(kit *rest.Kit, opt *metadata.CompareFieldTmplAttrOption, forUI bool) (
		*metadata.CompareFieldTmplAttrsRes, *metadata.ListFieldTmpltSyncStatusResult, error)
	CompareFieldTemplateUnique(kit *rest.Kit, opt *metadata.CompareFieldTmplUniqueOption, forUI bool) (
		*metadata.CompareFieldTmplUniquesRes, *metadata.ListFieldTmpltSyncStatusResult, error)
	ListFieldTemplateSyncStatus(kit *rest.Kit, option *metadata.ListFieldTmpltSyncStatusOption) (
		[]metadata.ListFieldTmpltSyncStatusResult, error)
	DeleteFieldTemplate(kit *rest.Kit, id int64) error
	DeleteFieldTemplateAttr(kit *rest.Kit, templateID int64, attrIDs []int64, needAuditLog bool) error
	DeleteFieldTemplateUnique(kit *rest.Kit, templateID int64, uniques []int64, needAuditLog bool) error
	UpdateFieldTemplateInfo(kit *rest.Kit, template *metadata.FieldTemplate) error
}

// NewFieldTemplateOperation create a new field template operation instance
func NewFieldTemplateOperation(client apimachinery.ClientSetInterface,
	asst model.AssociationOperationInterface) FieldTemplateOperation {

	return &template{
		clientSet:  client,
		asst:       asst,
		comparator: &comparator{clientSet: client, asst: asst},
	}
}

type template struct {
	clientSet  apimachinery.ClientSetInterface
	asst       model.AssociationOperationInterface
	comparator *comparator
}

// CreateFieldTemplate create field template(contains field template brief information, attributes and uniques)
func (f *template) CreateFieldTemplate(kit *rest.Kit, opt *metadata.CreateFieldTmplOption) (
	*metadata.RspID, error) {

	res, ccErr := f.clientSet.CoreService().FieldTemplate().CreateFieldTemplate(kit.Ctx, kit.Header, &opt.FieldTemplate)
	if ccErr != nil {
		blog.Errorf("create field template failed, err: %v, data: %v, rid: %s", ccErr, opt, kit.Rid)
		return nil, ccErr
	}

	for idx := range opt.Attributes {
		opt.Attributes[idx].TemplateID = res.ID
		opt.Attributes[idx].PropertyIndex = int64(idx)
	}
	attrIDs, ccErr := f.clientSet.CoreService().FieldTemplate().CreateFieldTemplateAttrs(kit.Ctx, kit.Header, res.ID,
		opt.Attributes)
	if ccErr != nil {
		blog.Errorf("create field template attributes failed, err: %v, data: %v, rid: %s", ccErr, opt.Attributes,
			kit.Rid)
		return nil, ccErr
	}

	propertyIDToIDMap := make(map[string]int64)
	for idx, attr := range opt.Attributes {
		propertyIDToIDMap[attr.PropertyID] = attrIDs.IDs[idx]
	}

	if len(opt.Uniques) == 0 {
		return res, nil
	}

	uniques := make([]metadata.FieldTemplateUnique, len(opt.Uniques))
	for idx, uniqueOpt := range opt.Uniques {
		uniqueOpt.TemplateID = res.ID
		unique, ccErr := uniqueOpt.Convert(propertyIDToIDMap)
		if ccErr.ErrCode != 0 {
			return nil, ccErr.ToCCError(kit.CCError)
		}

		uniques[idx] = *unique
	}

	uniqueIDs, ccErr := f.clientSet.CoreService().FieldTemplate().CreateFieldTemplateUniques(kit.Ctx, kit.Header,
		res.ID, uniques)
	if ccErr != nil {
		blog.Errorf("create field template uniques failed, err: %v, data: %v, rid: %s", ccErr, uniques, kit.Rid)
		return nil, ccErr
	}

	// generate and save audit log
	audit := auditlog.NewFieldTmplAuditLog(f.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLogs := make([]metadata.AuditLog, 0)

	tmplLog, err := audit.GenerateFieldTmplAuditLog(generateAuditParameter, res.ID, nil)
	if err != nil {
		blog.Errorf("generate field template audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	auditLogs = append(auditLogs, *tmplLog)

	attrLogs, err := audit.GenerateFieldTmplAttrAuditLog(generateAuditParameter, attrIDs.IDs, nil)
	if err != nil {
		blog.Errorf("generate field template attribute audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	auditLogs = append(auditLogs, attrLogs...)

	uniqueLogs, err := audit.GenerateFieldTmplUniqueAuditLog(generateAuditParameter, uniqueIDs.IDs, nil)
	if err != nil {
		blog.Errorf("generate field template attribute audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	auditLogs = append(auditLogs, uniqueLogs...)

	if err := audit.SaveAuditLog(kit, auditLogs...); err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return res, nil
}

// DeleteFieldTemplate delete field template
func (f *template) DeleteFieldTemplate(kit *rest.Kit, id int64) error {
	query := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: filtertools.GenAtomFilter(common.BKFieldID, filter.Equal, id),
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}

	res, ccErr := f.clientSet.CoreService().FieldTemplate().ListFieldTemplate(kit.Ctx, kit.Header, query)
	if ccErr != nil {
		blog.Errorf("find field template failed, opt: %+v, err: %v, rid: %s", query, ccErr, kit.Rid)
		return ccErr
	}

	if len(res.Info) > 1 {
		blog.Errorf("multiple field templates found, opt: %+v, rid: %s", query, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommGetMultipleObject)
	}

	if len(res.Info) == 0 {
		return nil
	}

	cond := metadata.DeleteOption{Condition: mapstr.MapStr{common.BKFieldID: id}}
	ccErr = f.clientSet.CoreService().FieldTemplate().DeleteFieldTemplate(kit.Ctx, kit.Header, &cond)
	if ccErr != nil {
		blog.Errorf("delete field template failed, cond: %v, err: %v, rid: %s", cond, ccErr, kit.Rid)
		return ccErr
	}

	// generate and save audit log
	audit := auditlog.NewFieldTmplAuditLog(f.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)

	tmplLog, err := audit.GenerateFieldTmplAuditLog(generateAuditParameter, id, &res.Info[0])
	if err != nil {
		blog.Errorf("generate field template audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err := audit.SaveAuditLog(kit, *tmplLog); err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// DeleteFieldTemplateAttr delete field template attribute
func (f *template) DeleteFieldTemplateAttr(kit *rest.Kit, templateID int64, attrIDs []int64, needAuditLog bool) error {
	expr := filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, templateID)
	if len(attrIDs) != 0 {
		var err error
		expr, err = filtertools.And(expr, filtertools.GenAtomFilter(common.BKFieldID, filter.In, attrIDs))
		if err != nil {
			blog.Errorf("build field template attribute filter failed, data: %v, err: %v, rid: %s", attrIDs, err,
				kit.Rid)
			return err
		}
	}
	query := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: expr,
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	res, ccErr := f.clientSet.CoreService().FieldTemplate().ListFieldTemplateAttr(kit.Ctx, kit.Header, query)
	if ccErr != nil {
		blog.Errorf("find field template attribute failed, opt: %+v, err: %v, rid: %s", query, ccErr, kit.Rid)
		return ccErr
	}

	if len(res.Info) == 0 {
		return nil
	}

	var deleteOpt *metadata.DeleteOption
	if len(attrIDs) != 0 {
		deleteOpt = &metadata.DeleteOption{
			Condition: mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: attrIDs}},
		}
	}

	ccErr = f.clientSet.CoreService().FieldTemplate().DeleteFieldTemplateAttr(kit.Ctx, kit.Header, templateID,
		deleteOpt)
	if ccErr != nil {
		blog.Errorf("delete field template attributes failed, template id: %d, cond: %v, err: %v, rid: %s", templateID,
			deleteOpt, ccErr, kit.Rid)
		return ccErr
	}

	if !needAuditLog {
		return nil
	}

	// generate and save audit log
	audit := auditlog.NewFieldTmplAuditLog(f.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)

	auditLogs, err := audit.GenerateFieldTmplAttrAuditLog(generateAuditParameter, attrIDs, res.Info)
	if err != nil {
		blog.Errorf("generate field template audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err := audit.SaveAuditLog(kit, auditLogs...); err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// DeleteFieldTemplateUnique delete field template unique
func (f *template) DeleteFieldTemplateUnique(kit *rest.Kit, templateID int64, uniqueIDs []int64,
	needAuditLog bool) error {

	expr := filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, templateID)
	if len(uniqueIDs) != 0 {
		var err error
		expr, err = filtertools.And(expr, filtertools.GenAtomFilter(common.BKFieldID, filter.In, uniqueIDs))
		if err != nil {
			blog.Errorf("build field template unique filter failed, data: %v, err: %v, rid: %s", uniqueIDs, err,
				kit.Rid)
			return err
		}
	}
	query := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: expr,
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	res, ccErr := f.clientSet.CoreService().FieldTemplate().ListFieldTemplateUnique(kit.Ctx, kit.Header, query)
	if ccErr != nil {
		blog.Errorf("find field template unique failed, opt: %+v, err: %v, rid: %s", query, ccErr, kit.Rid)
		return ccErr
	}

	if len(res.Info) == 0 {
		return nil
	}

	var deleteOpt *metadata.DeleteOption
	if len(uniqueIDs) != 0 {
		deleteOpt = &metadata.DeleteOption{
			Condition: mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: uniqueIDs}},
		}
	}

	ccErr = f.clientSet.CoreService().FieldTemplate().DeleteFieldTemplateUnique(kit.Ctx, kit.Header, templateID,
		deleteOpt)
	if ccErr != nil {
		blog.Errorf("delete field template uniques failed, template id: %d, cond: %v, err: %v, rid: %s", templateID,
			deleteOpt, ccErr, kit.Rid)
		return ccErr
	}

	if !needAuditLog {
		return nil
	}

	// generate and save audit log
	audit := auditlog.NewFieldTmplAuditLog(f.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)

	auditLogs, err := audit.GenerateFieldTmplUniqueAuditLog(generateAuditParameter, uniqueIDs, res.Info)
	if err != nil {
		blog.Errorf("generate field template audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err := audit.SaveAuditLog(kit, auditLogs...); err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// UpdateFieldTemplateInfo update field template brief information
func (f *template) UpdateFieldTemplateInfo(kit *rest.Kit, template *metadata.FieldTemplate) error {
	if template == nil {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "template")
	}

	ccErr := f.clientSet.CoreService().FieldTemplate().UpdateFieldTemplate(kit.Ctx, kit.Header, template)
	if ccErr != nil {
		blog.Errorf("update field template info failed, data: %v, err: %v, rid: %s", template, ccErr, kit.Rid)
		return ccErr
	}

	audit := auditlog.NewFieldTmplAuditLog(f.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate)
	tmplLog, err := audit.GenerateFieldTmplAuditLog(generateAuditParameter, template.ID, nil)
	if err != nil {
		blog.Errorf("generate field template audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err := audit.SaveAuditLog(kit, *tmplLog); err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}
