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
	"strconv"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// ListFieldTemplate list field templates.
// NOTICE: this api only returns basic info of field template, do not need to authorize
func (s *service) ListFieldTemplate(cts *rest.Contexts) {
	opt := new(metadata.CommonQueryOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// list field templates
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplate(cts.Kit.Ctx, cts.Kit.Header, opt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, req: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}

// FindFieldTemplateByID find the field combination template involved by ID.
func (s *service) FindFieldTemplateByID(ctx *rest.Contexts) {

	idStr := ctx.Request.PathParameter(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the path params id(%s), err: %v, rid: %s", idStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKFieldID))
		return
	}

	if id == 0 {
		blog.Errorf("template id cannot be zero, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKFieldID))
		return
	}

	query := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: filtertools.GenAtomFilter(common.BKFieldID, filter.Equal, id),
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}

	// list field template by id
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplate(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("get field template failed, req: %+v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(res.Info) > 1 {
		blog.Errorf("multiple field templates found, req: %+v, rid: %s", query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommGetMultipleObject))
		return
	}

	if len(res.Info) == 0 {
		ctx.RespEntity(nil)
		return
	}
	ctx.RespEntity(res.Info[0])
}

// canObjsBindFieldTemplate currently only supports "host" in the mainline model
func (s *service) canObjsBindFieldTemplate(kit *rest.Kit, ids []int64) ccErr.CCErrorCoder {

	query := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKFieldID: mapstr.MapStr{
				common.BKDBIN: ids,
			},
		},
		Fields:         []string{common.BKObjIDField},
		DisableCounter: true,
	}

	obj, err := s.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("failed to find objects by query(%#v), err: %v, rid: %s", query, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommNotFound, err.Error())
	}

	if len(obj.Info) == 0 {
		blog.Errorf("failed to find objIDs by query(%#v),  rid: %s", query, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommNotFound)
	}
	if len(obj.Info) != len(ids) {
		blog.Errorf("obj ids are invalid, input: %+v, rid: %s", query, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "object_ids")
	}

	objIDs := make([]string, 0)
	for _, info := range obj.Info {
		objIDs = append(objIDs, info.ObjectID)
	}

	cond := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
		common.BKDBAND: []mapstr.MapStr{
			{
				common.BKObjIDField: mapstr.MapStr{
					common.BKDBNE: common.BKInnerObjIDHost,
				},
			},
			{
				common.BKObjIDField: mapstr.MapStr{
					common.BKDBIN: objIDs,
				},
			},
		},
	}

	counts, ccErr := s.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameObjAsst, []map[string]interface{}{cond})
	if ccErr != nil {
		blog.Errorf("get mainline obj count failed, cond: %+v, err: %v, rid: %s", cond, ccErr, kit.Rid)
		return ccErr
	}

	if len(counts) != 1 || int(counts[0]) > 0 {
		blog.Errorf("obj ids cannot be a non-host mainline model, ids: %+v, rid: %s", ids, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "object_ids")
	}

	return nil
}

// FieldTemplateBindObject field template binding model
func (s *service) FieldTemplateBindObject(ctx *rest.Contexts) {

	opt := new(metadata.FieldTemplateBindObjOpt)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	objIDs := make([]int64, 0)
	objIDs = util.IntArrayUnique(opt.ObjectIDs)

	if authResp, authorized := s.authorizeObjsBindFieldTemplate(ctx.Kit, opt.ID, objIDs); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	if err := s.canObjsBindFieldTemplate(ctx.Kit, objIDs); err != nil {
		ctx.RespAutoError(err)
		return
	}

	option := &metadata.FieldTemplateBindObjOpt{
		ID:        opt.ID,
		ObjectIDs: objIDs,
	}
	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.clientSet.CoreService().FieldTemplate().FieldTemplateBindObject(ctx.Kit.Ctx, ctx.Kit.Header, option)
		if err != nil {
			blog.Errorf("field template bind model failed, err: %v , rid: %s", err, ctx.Kit.Rid)
			return err
		}

		audit := auditlog.NewObjectAuditLog(s.clientSet.CoreService())
		parameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)

		auditLogs, ccErr := audit.GenerateAuditLogForBindingFieldTemplate(parameter, objIDs, opt.ID)
		if ccErr != nil {
			blog.Errorf("generate audit log failed before update object, objName: %s, err: %v, rid: %s",
				opt.ObjectIDs, ccErr, ctx.Kit.Rid)
			return ccErr
		}
		// save audit log.
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, cond: %+v, err: %v, rid: %s", opt, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (s *service) authorizeObjsBindFieldTemplate(kit *rest.Kit, templateID int64, objIDs []int64) (
	*metadata.BaseResp, bool) {

	resource := make([]meta.ResourceAttribute, 0)
	for _, id := range objIDs {
		resource = append(resource, meta.ResourceAttribute{
			Basic: meta.Basic{
				Type:       meta.Model,
				Action:     meta.Update,
				InstanceID: id,
			},
		})
	}

	resource = append(resource, meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:       meta.FieldTemplate,
			Action:     meta.Find,
			InstanceID: templateID},
	})

	return s.auth.Authorize(kit, resource...)
}

// FieldTemplateUnbindObject field template binding model
func (s *service) FieldTemplateUnbindObject(ctx *rest.Contexts) {
	opt := new(metadata.FieldTemplateUnbindObjOpt)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	if authResp, authorized := s.auth.Authorize(ctx.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.Model, Action: meta.Update, InstanceID: opt.ObjectID}}); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	if err := s.canObjsBindFieldTemplate(ctx.Kit, []int64{opt.ObjectID}); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.clientSet.CoreService().FieldTemplate().FieldTemplateUnbindObject(ctx.Kit.Ctx, ctx.Kit.Header, opt)
		if err != nil {
			blog.Errorf("field template unbind model failed, cond: %+v, err: %v , rid: %s", opt, err, ctx.Kit.Rid)
			return err
		}
		audit := auditlog.NewObjectAuditLog(s.clientSet.CoreService())

		parameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		auditLogs, ccErr := audit.GenerateAuditLogForBindingFieldTemplate(parameter, []int64{opt.ObjectID}, opt.ID)
		if ccErr != nil {
			blog.Errorf("generate audit log failed , object id: %d, err: %v, rid: %s", opt.ObjectID, err, ctx.Kit.Rid)
			return ccErr
		}
		// save audit log.
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, cond: %+v, err: %v, rid: %s", opt, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// CreateFieldTemplate create field template(contains field template brief information, attributes and uniques)
func (s *service) CreateFieldTemplate(ctx *rest.Contexts) {
	opt := new(metadata.CreateFieldTmplOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	if authResp, authorized := s.auth.Authorize(ctx.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.FieldTemplate, Action: meta.Create}}); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	var res *metadata.RspID
	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		res, err = s.logics.FieldTemplateOperation().CreateFieldTemplate(ctx.Kit, opt)
		if err != nil {
			blog.Errorf("create field template failed, opt: %v, err: %v, rid: %s", opt, err, ctx.Kit.Rid)
			return err
		}

		// register business resource creator action to iam
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.FieldGroupingTemplate),
				ID:      strconv.FormatInt(res.ID, 10),
				Name:    opt.Name,
				Creator: ctx.Kit.User,
			}
			_, err = s.auth.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created field template to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(res)
}

// DeleteFieldTemplate delete field template(contains field template brief information, attributes and uniques)
func (s *service) DeleteFieldTemplate(ctx *rest.Contexts) {
	opt := new(metadata.DeleteFieldTmplOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	if authResp, authorized := s.auth.Authorize(ctx.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.FieldTemplate, Action: meta.Delete, InstanceID: opt.ID}}); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.logics.FieldTemplateOperation().DeleteFieldTemplateUnique(ctx.Kit, opt.ID, nil, false)
		if err != nil {
			return err
		}

		err = s.logics.FieldTemplateOperation().DeleteFieldTemplateAttr(ctx.Kit, opt.ID, nil, false)
		if err != nil {
			return err
		}

		err = s.logics.FieldTemplateOperation().DeleteFieldTemplate(ctx.Kit, opt.ID)
		if err != nil {
			return err
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

// CloneFieldTemplate clone field template(contains field template attributes and uniques)
func (s *service) CloneFieldTemplate(ctx *rest.Contexts) {
	opt := new(metadata.CloneFieldTmplOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	if err := s.judgeFieldTmplIsExist(ctx.Kit, opt.ID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	resources := []meta.ResourceAttribute{
		{Basic: meta.Basic{Type: meta.FieldTemplate, Action: meta.Create}},
		{Basic: meta.Basic{Type: meta.FieldTemplate, Action: meta.Find, InstanceID: opt.ID}},
	}
	if authResp, authorized := s.auth.Authorize(ctx.Kit, resources...); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	createOpt, err := s.buildCreateOpt(ctx.Kit, opt)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	var res *metadata.RspID
	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		res, err = s.logics.FieldTemplateOperation().CreateFieldTemplate(ctx.Kit, createOpt)
		if err != nil {
			blog.Errorf("create field template failed, opt: %v, err: %v, rid: %s", createOpt, err, ctx.Kit.Rid)
			return err
		}

		// register business resource creator action to iam
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.FieldGroupingTemplate),
				ID:      strconv.FormatInt(res.ID, 10),
				Name:    opt.Name,
				Creator: ctx.Kit.User,
			}
			_, err = s.auth.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created field template to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(res)
}

func (s *service) judgeFieldTmplIsExist(kit *rest.Kit, id int64) error {
	tmplFilter := filtertools.GenAtomFilter(common.BKFieldID, filter.Equal, id)
	tmplOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: tmplFilter},
		Page:               metadata.BasePage{EnableCount: true},
	}

	tmplInfo, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplate(kit.Ctx, kit.Header, tmplOpt)
	if err != nil {
		blog.Errorf("find field template failed, err: %v, opt: %+v, rid: %s", err, tmplOpt, kit.Rid)
		return err
	}

	if tmplInfo.Count != 1 {
		blog.Errorf("field template id is invalid, id: %d, rid: %s", id, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKFieldID)
	}

	return nil
}

func (s *service) buildCreateOpt(kit *rest.Kit, cloneOpt *metadata.CloneFieldTmplOption) (
	*metadata.CreateFieldTmplOption, error) {

	tmplOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, cloneOpt.ID),
		},
		Page: metadata.BasePage{Limit: common.BKNoLimit},
	}

	tmplAttrs, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplateAttr(kit.Ctx, kit.Header, tmplOpt)
	if err != nil {
		blog.Errorf("find field template attribute failed, opt: %+v, err: %v, rid: %s", tmplOpt, err, kit.Rid)
		return nil, err
	}

	idToPropertyIDMap := make(map[int64]string)
	for _, attr := range tmplAttrs.Info {
		idToPropertyIDMap[attr.ID] = attr.PropertyID
	}

	tmplUniques, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplateUnique(kit.Ctx, kit.Header, tmplOpt)
	if err != nil {
		blog.Errorf("find field template unique failed, opt: %+v, err: %v, rid: %s", tmplOpt, err, kit.Rid)
		return nil, err
	}

	uniques := make([]metadata.FieldTmplUniqueOption, len(tmplUniques.Info))
	for idx, unique := range tmplUniques.Info {
		createUnique, err := unique.Convert(idToPropertyIDMap)
		if err.ErrCode != 0 {
			return nil, err.ToCCError(kit.CCError)
		}
		uniques[idx] = *createUnique
	}

	createOpt := new(metadata.CreateFieldTmplOption)
	createOpt.FieldTemplate = cloneOpt.FieldTemplate
	createOpt.Attributes = tmplAttrs.Info
	createOpt.Uniques = uniques

	return createOpt, nil
}

// UpdateFieldTemplate update field template(contains field template brief information, attributes and uniques)
func (s *service) UpdateFieldTemplate(ctx *rest.Contexts) {
	opt := new(metadata.UpdateFieldTmplOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	if authResp, authorized := s.auth.Authorize(ctx.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.FieldTemplate, Action: meta.Update, InstanceID: opt.ID}}); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.logics.FieldTemplateOperation().UpdateFieldTemplateInfo(ctx.Kit, &opt.FieldTemplate); err != nil {
			blog.Errorf("update field template info failed, data: %v, err: %v, rid: %s", opt.FieldTemplate, err,
				ctx.Kit.Rid)
			return err
		}

		// because deleting attribute requires deleting its unique, we need to delete the unique first.
		if err := s.deleteFieldTmplUnique(ctx.Kit, opt.ID, opt.Uniques); err != nil {
			blog.Errorf("delete field template unique failed, template id: %d, cond: %v, err: %v, rid: %s", opt.ID,
				opt.Uniques, err, ctx.Kit.Rid)
			return err
		}

		propertyIDToIDMap, err := s.updateFieldTmplAttr(ctx.Kit, opt.ID, opt.Attributes)
		if err != nil {
			blog.Errorf("update field template attribute failed, data: %v, err: %v, rid: %s", opt.Attributes, err,
				ctx.Kit.Rid)
			return err
		}

		if err := s.updateFieldTmplUnique(ctx.Kit, opt.ID, propertyIDToIDMap, opt.Uniques); err != nil {
			blog.Errorf("update field template unique failed, data: %v, err: %v, rid: %s", opt.Uniques, err,
				ctx.Kit.Rid)
			return err
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)

}

func (s *service) deleteFieldTmplUnique(kit *rest.Kit, templateID int64,
	uniques []metadata.FieldTmplUniqueOption) error {

	dbIDMap, err := s.getFieldTmplUniqueIDs(kit, templateID)
	if err != nil {
		blog.Errorf("get field template unique ids failed, template id: %d, err: %v, rid: %s", templateID, err, kit.Rid)
		return err
	}

	if len(dbIDMap) == 0 {
		return nil
	}

	for _, unique := range uniques {
		if unique.ID == 0 {
			continue
		}

		delete(dbIDMap, unique.ID)
	}

	if len(dbIDMap) == 0 {
		return nil
	}

	deleteIDs := make([]int64, 0)
	for id := range dbIDMap {
		deleteIDs = append(deleteIDs, id)
	}

	err = s.logics.FieldTemplateOperation().DeleteFieldTemplateUnique(kit, templateID, deleteIDs, true)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) getFieldTmplUniqueIDs(kit *rest.Kit, templateID int64) (map[int64]struct{}, error) {
	uniqueFilter := filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, templateID)
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: uniqueFilter},
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
		Fields:             []string{common.BKFieldID},
	}

	uniques, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplateUnique(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field template uniques failed, opt: %+v, err: %v, rid: %s", uniqueFilter, err, kit.Rid)
		return nil, err
	}

	result := make(map[int64]struct{})
	for _, unique := range uniques.Info {
		result[unique.ID] = struct{}{}
	}

	return result, nil
}

// updateFieldTmplAttr contains update, create and delete field template attribute
func (s *service) updateFieldTmplAttr(kit *rest.Kit, templateID int64, attrs []metadata.FieldTemplateAttr) (
	map[string]int64, error) {

	attrOp, err := s.getFieldTmplAttrOperation(kit, templateID, attrs)
	if err != nil {
		return nil, err
	}

	if len(attrOp.deleteAttrIDs) != 0 {
		err = s.logics.FieldTemplateOperation().DeleteFieldTemplateAttr(kit, templateID, attrOp.deleteAttrIDs, true)
		if err != nil {
			return nil, err
		}
	}

	propertyIDToIDMap := make(map[string]int64)

	if len(attrOp.createAttrs) == 0 && len(attrOp.updateAttrs) == 0 {
		return propertyIDToIDMap, nil
	}

	audit := auditlog.NewFieldTmplAuditLog(s.clientSet.CoreService())
	auditLogs := make([]metadata.AuditLog, 0)
	if len(attrOp.createAttrs) != 0 {
		resp, ccErr := s.clientSet.CoreService().FieldTemplate().CreateFieldTemplateAttrs(kit.Ctx, kit.Header,
			templateID, attrOp.createAttrs)
		if ccErr != nil {
			blog.Errorf("create field template attribute failed, data: %v, err: %v, rid: %s", attrOp.createAttrs, ccErr,
				kit.Rid)
			return nil, ccErr
		}
		for idx, attr := range attrOp.createAttrs {
			propertyIDToIDMap[attr.PropertyID] = resp.IDs[idx]
		}

		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
		createLogs, err := audit.GenerateFieldTmplAttrAuditLog(generateAuditParameter, resp.IDs, nil)
		if err != nil {
			blog.Errorf("generate field template attribute audit log failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		auditLogs = append(auditLogs, createLogs...)
	}

	if len(attrOp.updateAttrs) != 0 {
		err = s.clientSet.CoreService().FieldTemplate().UpdateFieldTemplateAttrs(kit.Ctx, kit.Header, templateID,
			attrOp.updateAttrs)
		if err != nil {
			blog.Errorf("update field template attributes failed, template id: %d, data: %v, err: %v, rid: %s",
				templateID, attrOp.updateAttrs, err, kit.Rid)
			return nil, err
		}
		for _, attr := range attrOp.updateAttrs {
			propertyIDToIDMap[attr.PropertyID] = attr.ID
		}

		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate)
		updateLogs, err := audit.GenerateFieldTmplAttrAuditLog(generateAuditParameter, nil, attrOp.updateAttrs)
		if err != nil {
			blog.Errorf("generate field template attribute audit log failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		auditLogs = append(auditLogs, updateLogs...)
	}

	if err := audit.SaveAuditLog(kit, auditLogs...); err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return propertyIDToIDMap, nil
}

type attrOperation struct {
	createAttrs   []metadata.FieldTemplateAttr
	updateAttrs   []metadata.FieldTemplateAttr
	deleteAttrIDs []int64
}

func (s *service) getFieldTmplAttrOperation(kit *rest.Kit, templateID int64, attrs []metadata.FieldTemplateAttr) (
	op *attrOperation, err error) {

	dbIDMap, err := s.getFieldTmplAttrIDs(kit, templateID)
	if err != nil {
		blog.Errorf("get field template attribute ids failed, template id: %d, err: %v, rid: %s", templateID, err,
			kit.Rid)
		return nil, err
	}

	updateAttrs := make([]metadata.FieldTemplateAttr, 0)
	createAttrs := make([]metadata.FieldTemplateAttr, 0)

	for idx, attr := range attrs {
		attr.TemplateID = templateID
		attr.PropertyIndex = int64(idx)

		if attr.ID == 0 {
			createAttrs = append(createAttrs, attr)
			continue
		}

		updateAttrs = append(updateAttrs, attr)
		delete(dbIDMap, attr.ID)
	}

	deleteIDs := make([]int64, 0)
	for id := range dbIDMap {
		deleteIDs = append(deleteIDs, id)
	}

	return &attrOperation{
		createAttrs:   createAttrs,
		updateAttrs:   updateAttrs,
		deleteAttrIDs: deleteIDs,
	}, nil
}

func (s *service) getFieldTmplAttrIDs(kit *rest.Kit, templateID int64) (map[int64]struct{}, error) {
	attrFilter := filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, templateID)
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: attrFilter},
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
		Fields:             []string{common.BKFieldID},
	}

	attrs, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplateAttr(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field template attribute failed, opt: %+v, err: %v, rid: %s", attrFilter, err, kit.Rid)
		return nil, err
	}

	result := make(map[int64]struct{})
	for _, attr := range attrs.Info {
		result[attr.ID] = struct{}{}
	}

	return result, nil
}

// updateFieldTmplUnique contains update and create field template unique
func (s *service) updateFieldTmplUnique(kit *rest.Kit, templateID int64, propertyIDToIDMap map[string]int64,
	uniques []metadata.FieldTmplUniqueOption) error {

	if len(uniques) == 0 {
		return nil
	}

	updateUniques := make([]metadata.FieldTemplateUnique, 0)
	createUniques := make([]metadata.FieldTemplateUnique, 0)

	for _, uniqueOpt := range uniques {
		unique, err := uniqueOpt.Convert(propertyIDToIDMap)
		if err.ErrCode != 0 {
			return err.ToCCError(kit.CCError)
		}
		unique.TemplateID = templateID

		if unique.ID == 0 {
			createUniques = append(createUniques, *unique)
			continue
		}

		updateUniques = append(updateUniques, *unique)
	}

	auditLogs := make([]metadata.AuditLog, 0)
	audit := auditlog.NewFieldTmplAuditLog(s.clientSet.CoreService())
	if len(updateUniques) != 0 {
		ccErr := s.clientSet.CoreService().FieldTemplate().UpdateFieldTemplateUniques(kit.Ctx, kit.Header, templateID,
			updateUniques)
		if ccErr != nil {
			blog.Errorf("update field template uniques failed, template id: %d, data: %v, err: %v, rid: %s", templateID,
				updateUniques, ccErr, kit.Rid)
			return ccErr
		}

		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate)
		updateLogs, err := audit.GenerateFieldTmplUniqueAuditLog(generateAuditParameter, nil, updateUniques)
		if err != nil {
			blog.Errorf("generate field template unique audit log failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		auditLogs = append(auditLogs, updateLogs...)
	}

	if len(createUniques) != 0 {
		resp, ccErr := s.clientSet.CoreService().FieldTemplate().CreateFieldTemplateUniques(kit.Ctx, kit.Header,
			templateID, createUniques)
		if ccErr != nil {
			blog.Errorf("create field template uniques failed, data: %v, err: %v, rid: %s", createUniques, ccErr,
				kit.Rid)
			return ccErr
		}

		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
		createLogs, err := audit.GenerateFieldTmplUniqueAuditLog(generateAuditParameter, resp.IDs, nil)
		if err != nil {
			blog.Errorf("generate field template unique audit log failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		auditLogs = append(auditLogs, createLogs...)
	}

	if len(auditLogs) == 0 {
		return nil
	}

	if err := audit.SaveAuditLog(kit, auditLogs...); err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// UpdateFieldTemplateInfo update field template brief information
func (s *service) UpdateFieldTemplateInfo(ctx *rest.Contexts) {
	opt := new(metadata.FieldTemplate)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if opt.ID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKFieldID))
		return
	}

	if err := opt.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	if authResp, authorized := s.auth.Authorize(ctx.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.FieldTemplate, Action: meta.Update, InstanceID: opt.ID}}); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.logics.FieldTemplateOperation().UpdateFieldTemplateInfo(ctx.Kit, opt); err != nil {
			blog.Errorf("update field template info failed, data: %v, err: %v, rid: %s", opt, err, ctx.Kit.Rid)
			return err
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

// ListFieldTmplByUniqueTmplIDForUI query the ID and name of the corresponding field
// template according to the TemplateID on the unique verification of the model
func (s *service) ListFieldTmplByUniqueTmplIDForUI(cts *rest.Contexts) {
	opt := new(metadata.ListTmplSimpleByUniqueOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// list field templates brief info
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTmplSimplyByUniqueTemplateID(cts.Kit.Ctx,
		cts.Kit.Header, opt)
	if err != nil {
		blog.Errorf("list field templates failed, req: %+v, err: %v, rid: %s", opt, err, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}

// ListFieldTmplByObjectTmplIDForUI query the ID and name of the corresponding field
// template according to the TemplateID on the model attribute
func (s *service) ListFieldTmplByObjectTmplIDForUI(cts *rest.Contexts) {
	opt := new(metadata.ListTmplSimpleByAttrOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// list field templates brief info.
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTmplSimplyByAttrTemplateID(cts.Kit.Ctx,
		cts.Kit.Header, opt)
	if err != nil {
		blog.Errorf("list field templates brief info failed, req: %+v, err: %v, rid: %s", opt, err, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}
