/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"reflect"
	"strconv"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateSetTemplate TODO
func (s *Service) CreateSetTemplate(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.CreateSetTemplateOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(option.ServiceTemplateIDs) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "service_template_ids"))
		return
	}

	var setTemplate metadata.SetTemplate
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		setTemplate, err = s.Engine.CoreAPI.CoreService().SetTemplate().CreateSetTemplate(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option)
		if err != nil {
			blog.Errorf("CreateSetTemplate failed, core service create failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
			return err
		}

		// register set template resource creator action to iam
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.BizSetTemplate),
				ID:      strconv.FormatInt(setTemplate.ID, 10),
				Name:    setTemplate.Name,
				Creator: ctx.Kit.User,
			}
			_, err = s.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created set template to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(setTemplate)
}

// CreateSetTemplateAllInfo create set template all info, including attributes and service template relations
func (s *Service) CreateSetTemplateAllInfo(ctx *rest.Contexts) {
	option := new(metadata.CreateSetTempAllInfoOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := option.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	var templateID int64
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// create set template
		setTempOpt := metadata.CreateSetTemplateOption{
			Name:               option.Name,
			ServiceTemplateIDs: option.ServiceTemplateIDs,
		}

		setTemplate, err := s.Engine.CoreAPI.CoreService().SetTemplate().CreateSetTemplate(ctx.Kit.Ctx, ctx.Kit.Header,
			option.BizID, setTempOpt)
		if err != nil {
			blog.Errorf("create set template failed, bizID: %d, option: %+v, err: %+v, rid: %s", option.BizID,
				setTempOpt, err, ctx.Kit.Rid)
			return err
		}

		templateID = setTemplate.ID

		// create set template attributes
		if len(option.Attributes) > 0 {
			attrOpt := &metadata.CreateSetTempAttrsOption{
				BizID:      option.BizID,
				ID:         templateID,
				Attributes: option.Attributes,
			}

			if _, err = s.Engine.CoreAPI.CoreService().SetTemplate().CreateSetTemplateAttribute(ctx.Kit.Ctx, ctx.Kit.Header,
				attrOpt); err != nil {
				blog.Errorf("create set template attrs(%+v) failed, err: %v, rid: %s", attrOpt, err, ctx.Kit.Rid)
				return err
			}
		}

		// register set template resource creator action to iam
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.BizSetTemplate),
				ID:      strconv.FormatInt(setTemplate.ID, 10),
				Name:    setTemplate.Name,
				Creator: ctx.Kit.User,
			}
			_, err := s.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register set template(%+v) to iam failed, err: %v, rid: %s", iamInstance, err, ctx.Kit.Rid)
				return err
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(metadata.RspID{ID: templateID})
}

// UpdateSetTemplate TODO
func (s *Service) UpdateSetTemplate(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	option := metadata.UpdateSetTemplateOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var setTemplate metadata.SetTemplate
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		setTemplate, err = s.Engine.CoreAPI.CoreService().SetTemplate().UpdateSetTemplate(ctx.Kit.Ctx, ctx.Kit.Header, bizID, setTemplateID, option)
		if err != nil {
			blog.Errorf("UpdateSetTemplate failed, do core service update failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(setTemplate)
}

// UpdateSetTemplateAllInfo update set template all info, including attributes and service template relations
func (s *Service) UpdateSetTemplateAllInfo(ctx *rest.Contexts) {
	option := new(metadata.UpdateSetTempAllInfoOption)
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	allInfo, err := s.getSetTemplateAllInfo(ctx.Kit, option.ID, option.BizID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// compare if service template relations changed
		svcTempMap := make(map[int64]struct{})
		for _, svcTempID := range allInfo.ServiceTemplateIDs {
			svcTempMap[svcTempID] = struct{}{}
		}

		isSvcTempDiff := false
		for _, svcTempID := range option.ServiceTemplateIDs {
			if _, exists := svcTempMap[svcTempID]; !exists {
				isSvcTempDiff = true
				break
			}
			delete(svcTempMap, svcTempID)
		}

		if len(svcTempMap) != 0 {
			isSvcTempDiff = true
		}

		// update set template name and service template relations if there is a difference
		if option.Name != allInfo.Name || isSvcTempDiff {
			svcTempOpt := metadata.UpdateSetTemplateOption{
				Name:               option.Name,
				ServiceTemplateIDs: option.ServiceTemplateIDs,
			}

			if _, err := s.Engine.CoreAPI.CoreService().SetTemplate().UpdateSetTemplate(ctx.Kit.Ctx, ctx.Kit.Header,
				option.BizID, option.ID, svcTempOpt); err != nil {
				blog.Errorf("update set template %d failed, opt: %+v, err: %+v, rid: %s", option.ID, svcTempOpt, err,
					ctx.Kit.Rid)
				return err
			}
		}

		// update service template attributes
		err = s.updateSetTempAllAttrs(ctx.Kit, allInfo.ID, allInfo.BizID, allInfo.Attributes, option.Attributes)
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

// updateSetTempAllAttrs update set template attributes, add new attributes and delete redundant attributes
func (s *Service) updateSetTempAllAttrs(kit *rest.Kit, id, bizID int64, prevAttrs []metadata.SetTemplateAttr,
	updateAttrs []metadata.SetTempAttr) errors.CCErrorCoder {

	attrMap := make(map[int64]interface{})
	for _, attribute := range prevAttrs {
		attrMap[attribute.AttributeID] = attribute
	}

	// cross compare previous attributes and update attributes to find need add/update/delete attributes
	var addedAttrs, updatedAttrs []metadata.SetTempAttr
	for _, attribute := range updateAttrs {
		value, exists := attrMap[attribute.AttributeID]
		if !exists {
			addedAttrs = append(addedAttrs, metadata.SetTempAttr{
				AttributeID:   attribute.AttributeID,
				PropertyValue: attribute.PropertyValue,
			})
			continue
		}

		delete(attrMap, attribute.AttributeID)
		if !reflect.DeepEqual(value, attribute.PropertyValue) {
			updatedAttrs = append(updatedAttrs, metadata.SetTempAttr{
				AttributeID:   attribute.AttributeID,
				PropertyValue: attribute.PropertyValue,
			})
		}
	}

	// add set template attributes
	if len(addedAttrs) > 0 {
		addOpt := &metadata.CreateSetTempAttrsOption{
			BizID:      bizID,
			ID:         id,
			Attributes: addedAttrs,
		}

		_, err := s.Engine.CoreAPI.CoreService().SetTemplate().CreateSetTemplateAttribute(kit.Ctx, kit.Header, addOpt)
		if err != nil {
			blog.Errorf("add set template attrs failed, opt: %+v, err: %v, rid: %s", addOpt, err, kit.Rid)
			return err
		}
	}

	// update set template attributes
	if len(updatedAttrs) > 0 {
		updateOpt := &metadata.UpdateSetTempAttrOption{
			BizID:      bizID,
			ID:         id,
			Attributes: updatedAttrs,
		}

		err := s.Engine.CoreAPI.CoreService().SetTemplate().UpdateSetTemplateAttribute(kit.Ctx, kit.Header, updateOpt)
		if err != nil {
			blog.Errorf("update set template attrs failed, opt: %+v, err: %v, rid: %s", updateOpt, err, kit.Rid)
			return err
		}
	}

	// delete set template attributes
	if len(attrMap) > 0 {
		deletedAttrIDs := make([]int64, 0)
		for attrID := range attrMap {
			deletedAttrIDs = append(deletedAttrIDs, attrID)
		}

		deleteOpt := &metadata.DeleteSetTempAttrOption{
			BizID:        bizID,
			ID:           id,
			AttributeIDs: deletedAttrIDs,
		}
		err := s.Engine.CoreAPI.CoreService().SetTemplate().DeleteSetTemplateAttribute(kit.Ctx, kit.Header, deleteOpt)
		if err != nil {
			blog.Errorf("delete set template attrs failed, opt: %+v, err: %v, rid: %s", deleteOpt, err, kit.Rid)
			return err
		}
	}

	return nil
}

// DeleteSetTemplate TODO
func (s *Service) DeleteSetTemplate(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.DeleteSetTemplateOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.Engine.CoreAPI.CoreService().SetTemplate().DeleteSetTemplate(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option); err != nil {
			blog.Errorf("DeleteSetTemplate failed, do core service update failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
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

// GetSetTemplate TODO
func (s *Service) GetSetTemplate(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	setTemplate, err := s.Engine.CoreAPI.CoreService().SetTemplate().GetSetTemplate(ctx.Kit.Ctx, ctx.Kit.Header, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("GetSetTemplate failed, do core service get failed, bizID: %d, setTemplateID: %d, err: %+v, rid: %s", bizID, setTemplateID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(setTemplate)
}

// GetSetTemplateAllInfo get set template all info, including attributes and service template relations
func (s *Service) GetSetTemplateAllInfo(ctx *rest.Contexts) {
	option := new(metadata.GetSetTempAllInfoOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := option.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	allInfo, err := s.getSetTemplateAllInfo(ctx.Kit, option.ID, option.BizID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(allInfo)
}

func (s *Service) getSetTemplateAllInfo(kit *rest.Kit, id, bizID int64) (*metadata.SetTempAllInfo,
	errors.CCErrorCoder) {

	// get set template
	setTemplate, err := s.Engine.CoreAPI.CoreService().SetTemplate().GetSetTemplate(kit.Ctx, kit.Header, bizID, id)
	if err != nil {
		blog.Errorf("get set template failed, id: %d, bizID: %d, err: %v, rid: %s", id, bizID, err, kit.Rid)
		return nil, err
	}

	// get set template attributes
	attrOpt := &metadata.ListSetTempAttrOption{
		BizID: bizID,
		ID:    id,
	}
	attrs, err := s.Engine.CoreAPI.CoreService().SetTemplate().ListSetTemplateAttribute(kit.Ctx, kit.Header, attrOpt)
	if err != nil {
		blog.Errorf("get set template %d attributes failed, err: %v, rid: %s", id, err, kit.Rid)
		return nil, err
	}

	// get service template relations
	relations, err := s.Engine.CoreAPI.CoreService().SetTemplate().ListSetServiceTemplateRelations(kit.Ctx, kit.Header,
		bizID, id)
	if err != nil {
		blog.Errorf("get set template %d service template relations failed, err: %v, rid: %s", id, err, kit.Rid)
		return nil, err
	}

	serviceTemplateIDs := make([]int64, len(relations))
	for idx, relation := range relations {
		serviceTemplateIDs[idx] = relation.ServiceTemplateID
	}

	return &metadata.SetTempAllInfo{
		ID:                 setTemplate.ID,
		BizID:              setTemplate.BizID,
		Name:               setTemplate.Name,
		ServiceTemplateIDs: serviceTemplateIDs,
		Attributes:         attrs.Attributes,
	}, nil
}

// ListSetTemplate TODO
func (s *Service) ListSetTemplate(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.ListSetTemplateOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}
	// set default value
	if option.Page.Limit == 0 {
		option.Page.Limit = common.BKDefaultLimit
	}

	setTemplate, err := s.Engine.CoreAPI.CoreService().SetTemplate().ListSetTemplate(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option)
	if err != nil {
		blog.Errorf("ListSetTemplate failed, do core service ListSetTemplate failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(setTemplate)
}

// ListSetTemplateWeb TODO
func (s *Service) ListSetTemplateWeb(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	listOption := metadata.ListSetTemplateOption{}
	if err := ctx.DecodeInto(&listOption); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if listOption.Page.Limit == 0 {
		listOption.Page.Limit = common.BKNoLimit
	}

	listResult, err := s.Engine.CoreAPI.CoreService().SetTemplate().ListSetTemplate(ctx.Kit.Ctx, ctx.Kit.Header, bizID, listOption)
	if err != nil {
		blog.Errorf("ListSetTemplate failed, do core service ListSetTemplate failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, listOption, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if listResult == nil {
		ctx.RespEntity(nil)
		return
	}

	// count template instances
	setTemplateIDs := make([]int64, 0)
	for _, item := range listResult.Info {
		setTemplateIDs = append(setTemplateIDs, item.ID)
	}
	option := metadata.CountSetTplInstOption{
		SetTemplateIDs: setTemplateIDs,
	}
	setTplInstCount, err := s.Engine.CoreAPI.CoreService().SetTemplate().CountSetTplInstances(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option)
	if err != nil {
		blog.Errorf("ListSetTemplateWeb failed, CountSetTplInstances failed, bizID: %d, option: %+v, err: %s, rid: %s", bizID, option, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	result := metadata.MultipleSetTemplateWithStatisticsResult{
		Count: listResult.Count,
	}
	for _, setTemplate := range listResult.Info {
		setInstanceCount, exist := setTplInstCount[setTemplate.ID]
		if exist == false {
			setInstanceCount = 0
		}
		result.Info = append(result.Info, metadata.SetTemplateWithStatistics{
			SetInstanceCount: setInstanceCount,
			SetTemplate:      setTemplate,
		})
	}
	ctx.RespEntity(result)
}

// ListSetTplRelatedSvcTpl TODO
func (s *Service) ListSetTplRelatedSvcTpl(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	serviceTemplates, err := s.Engine.CoreAPI.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(ctx.Kit.Ctx, ctx.Kit.Header, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("ListSetTemplateRelatedServiceTemplate failed, ListSetTplRelatedSvcTpl failed, bizID: %d, setTemplateID: %+v, err: %+v, rid: %s", bizID, setTemplateID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(serviceTemplates)
}

// ListSetTplRelatedSvcTplWithStatistics search set template and service template related by statistics
func (s *Service) ListSetTplRelatedSvcTplWithStatistics(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	serviceTemplates, err := s.Engine.CoreAPI.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(ctx.Kit.Ctx,
		ctx.Kit.Header, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("ListSetTplRelatedSvcTplWithStatistics failed, do core service list failed, bizID: %d, "+
			"setTemplateID: %+v, err: %+v, rid: %s", bizID, setTemplateID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	serviceTemplateIDs := make([]int64, 0)
	for _, item := range serviceTemplates {
		serviceTemplateIDs = append(serviceTemplateIDs, item.ID)
	}
	moduleFilter := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKServiceTemplateIDField: map[string]interface{}{
				common.BKDBIN: serviceTemplateIDs,
			},
			common.BKSetTemplateIDField: setTemplateID,
		},
	}

	moduleResult := new(metadata.ResponseModuleInstance)
	if err := s.Engine.CoreAPI.CoreService().Instance().ReadInstanceStruct(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDModule, moduleFilter, moduleResult); err != nil {
		blog.ErrorJSON("ListSetTplRelatedSvcTplWithStatistics failed, ReadInstance of module http failed, "+
			"option: %s, err: %s, rid: %s", moduleFilter, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if ccErr := moduleResult.CCError(); ccErr != nil {
		blog.Errorf("ListSetTplRelatedSvcTplWithStatistics failed, ReadInstance of module failed, filter: %s, "+
			"result: %s, rid: %s", moduleFilter, moduleResult, ctx.Kit.Rid)
		ctx.RespAutoError(ccErr)
		return
	}
	moduleIDs := make([]int64, 0)
	svcTpl2Modules := make(map[int64][]metadata.ModuleInst)
	// map[module]service_template_id
	moduleIDSvcTplID := make(map[int64]int64, 0)
	for _, module := range moduleResult.Data.Info {
		moduleIDSvcTplID[module.ModuleID] = module.ServiceTemplateID
		svcTpl2Modules[module.ServiceTemplateID] = append(svcTpl2Modules[module.ServiceTemplateID], module)
		moduleIDs = append(moduleIDs, module.ModuleID)
	}

	// host module relations
	relationOption := metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   moduleIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKModuleIDField, common.BKHostIDField},
	}
	relationResult, err := s.Engine.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header,
		&relationOption)
	if err != nil {
		blog.Errorf("ListSetTplRelatedSvcTplWithStatistics failed, GetHostModuleRelation http failed, option: %s, "+
			"err: %s, rid: %s", relationOption, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}

	// module hosts
	svcTplIDHostIDs := make(map[int64][]int64)
	for _, item := range relationResult.Info {
		if svcTplID, ok := moduleIDSvcTplID[item.ModuleID]; ok {
			svcTplIDHostIDs[svcTplID] = append(svcTplIDHostIDs[svcTplID], item.HostID)
		}
	}

	result := make([]metadata.ServiceTemplateWithModuleInfo, 0)
	for _, svcTpl := range serviceTemplates {
		info := metadata.ServiceTemplateWithModuleInfo{
			ServiceTemplate: svcTpl,
		}
		modules, ok := svcTpl2Modules[svcTpl.ID]
		if ok == false {
			result = append(result, info)
			continue
		}
		info.Modules = modules
		info.HostCount = len(util.IntArrayUnique(svcTplIDHostIDs[svcTpl.ID]))
		result = append(result, info)
	}

	ctx.RespEntity(result)
}

// ListSetTplRelatedSets get SetTemplate related sets
func (s *Service) ListSetTplRelatedSets(kit *rest.Kit, bizID int64, setTemplateID int64,
	option metadata.ListSetByTemplateOption) (*metadata.InstDataInfo, error) {
	filter := map[string]interface{}{
		common.BKAppIDField:         bizID,
		common.BKSetTemplateIDField: setTemplateID,
	}
	if option.SetIDs != nil {
		filter[common.BKSetIDField] = map[string]interface{}{
			common.BKDBIN: option.SetIDs,
		}
	}
	qc := &metadata.QueryCondition{
		Page:      option.Page,
		Condition: filter,
	}
	return s.Engine.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, qc)
}

// ListSetTplRelatedSetsWeb get SetTemplate related sets, just for web
func (s *Service) ListSetTplRelatedSetsWeb(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}
	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	option := metadata.ListSetByTemplateOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	response, err := s.ListSetTplRelatedSets(ctx.Kit, bizID, setTemplateID, option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	setInstanceResult := response

	topoTree, err := s.Engine.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(ctx.Kit.Ctx,
		ctx.Kit.Header, bizID, false)
	if err != nil {
		blog.Errorf("ListSetTplRelatedSetsWeb failed, bizID: %d, err: %s, rid: %s", bizID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	setIDs := make([]int64, 0)
	for index := range setInstanceResult.Info {
		set := metadata.SetInst{}
		if err := mapstr.DecodeFromMapStr(&set, setInstanceResult.Info[index]); err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
			return
		}
		setIDs = append(setIDs, set.SetID)

		setPath := topoTree.TraversalFindNode(common.BKInnerObjIDSet, set.SetID)
		topoPath := make([]metadata.TopoInstanceNodeSimplify, 0)
		for _, pathNode := range setPath {
			nodeSimplify := metadata.TopoInstanceNodeSimplify{
				ObjectID:     pathNode.ObjectID,
				InstanceID:   pathNode.InstanceID,
				InstanceName: pathNode.InstanceName,
			}
			topoPath = append(topoPath, nodeSimplify)
		}
		setInstanceResult.Info[index]["topo_path"] = topoPath
	}

	// fill with host count
	filter := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		SetIDArr:      setIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKSetIDField, common.BKHostIDField},
	}
	relations, err := s.Engine.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header, filter)
	if err != nil {
		blog.ErrorJSON("SearchMainlineInstanceTopo failed, GetHostModuleRelation failed, filter: %s, err: %s, "+
			"rid: %s", filter, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	set2Hosts := make(map[int64][]int64)
	for _, relation := range relations.Info {
		if _, ok := set2Hosts[relation.SetID]; ok == false {
			set2Hosts[relation.SetID] = make([]int64, 0)
		}
		set2Hosts[relation.SetID] = append(set2Hosts[relation.SetID], relation.HostID)
	}
	for setID := range set2Hosts {
		set2Hosts[setID] = util.IntArrayUnique(set2Hosts[setID])
	}

	for index := range setInstanceResult.Info {
		set := metadata.SetInst{}
		if err := mapstr.DecodeFromMapStr(&set, setInstanceResult.Info[index]); err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
			return
		}
		hostCount := 0
		if _, ok := set2Hosts[set.SetID]; ok == true {
			hostCount = len(set2Hosts[set.SetID])
		}
		setInstanceResult.Info[index]["host_count"] = hostCount
	}

	ctx.RespEntity(setInstanceResult)
}

// SetWithHostFlag 获取集群中要删除的服务模板实例化的模块是否有主机
func (s *Service) SetWithHostFlag(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}
	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	op := metadata.SetWithHostFlagOption{}
	if err := ctx.DecodeInto(&op); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 这里需要判断SetIDs的合法性
	if rawErr := op.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	setModules, err := s.Logics.SetTemplateOperation().SetWithDeleteModulesRelation(ctx.Kit, bizID, setTemplateID, op)
	if err != nil {
		blog.Errorf("get modules failed, bizID: %d, setTemplateID: %d, option: %+v, err: %v, rid: %s", bizID,
			setTemplateID, op, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	moduleIDs := make([]int64, 0)
	for _, modules := range setModules {
		moduleIDs = append(moduleIDs, modules...)
	}

	result := make([]metadata.SetWithHostFlagResult, 0)

	if len(moduleIDs) == 0 {
		for _, setID := range op.SetIDs {
			result = append(result, metadata.SetWithHostFlagResult{
				ID:      setID,
				HasHost: false,
			})
		}
		blog.Warnf("no modules founded, bizID: %d, setTemplateID: %d, option: %+v, rid: %s", bizID, setTemplateID,
			op, ctx.Kit.Rid)
		ctx.RespEntity(result)
		return
	}

	relationOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   moduleIDs,
		Page:          metadata.BasePage{Limit: common.BKNoLimit},
		Fields:        []string{common.BKSetIDField, common.BKModuleIDField},
	}

	relationResult, err := s.Engine.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx,
		ctx.Kit.Header, relationOption)
	if err != nil {
		blog.Errorf("get host module relation failed, option: %+v, err: %v, rid: %s", relationOption, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	setMap := make(map[int64]struct{})
	for _, set := range relationResult.Info {
		setMap[set.SetID] = struct{}{}
	}

	for _, setID := range op.SetIDs {
		if _, ok := setMap[setID]; ok {
			result = append(result, metadata.SetWithHostFlagResult{
				ID:      setID,
				HasHost: true})
			continue
		}
		result = append(result, metadata.SetWithHostFlagResult{
			ID:      setID,
			HasHost: false})
	}
	ctx.RespEntity(result)
}

// DiffSetTplWithInst search different between set template and set inst
func (s *Service) DiffSetTplWithInst(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	option := metadata.DiffSetTplWithInstOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if option.SetID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetIDField))
		return
	}

	serviceTemplates, err := s.Engine.CoreAPI.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(ctx.Kit.Ctx,
		ctx.Kit.Header, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("list service templates failed, bizID: %d, setTemplateID: %d, err: %v, rid: %s", bizID,
			setTemplateID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	setDiff, err := s.Logics.SetTemplateOperation().DiffSetTplWithInst(ctx.Kit, bizID, setTemplateID, option,
		serviceTemplates)
	if err != nil {
		blog.Errorf("DiffSetTplWithInst failed, operation failed, bizID: %d, setTemplateID: %d, option: %+v, err: %s,"+
			" rid: %s", bizID, setTemplateID, option, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	moduleIDs := make([]int64, 0)
	for _, moduleDiff := range setDiff.ModuleDiffs {
		if moduleDiff.ModuleID == 0 {
			continue
		}
		moduleIDs = append(moduleIDs, moduleDiff.ModuleID)
	}

	result := metadata.SetTplDiffResult{
		Difference:      setDiff,
		ModuleHostCount: make(map[int64]int64),
	}

	if len(moduleIDs) > 0 {
		relationOption := &metadata.HostModuleRelationRequest{
			ApplicationID: bizID,
			ModuleIDArr:   moduleIDs,
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
			Fields: []string{common.BKModuleIDField},
		}
		relationResult, err := s.Engine.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx,
			ctx.Kit.Header, relationOption)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		moduleHostsCount := make(map[int64]int64)
		for _, item := range relationResult.Info {
			if _, exist := moduleHostsCount[item.ModuleID]; exist == false {
				moduleHostsCount[item.ModuleID] = 0
			}
			moduleHostsCount[item.ModuleID] += 1
		}
		for _, moduleID := range moduleIDs {
			if _, exist := moduleHostsCount[moduleID]; exist == false {
				moduleHostsCount[moduleID] = 0
			}
		}
		result.ModuleHostCount = moduleHostsCount
	}

	ctx.RespEntity(result)
}

// SyncSetTplToInst  sync set template to set inst
func (s *Service) SyncSetTplToInst(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	option := metadata.SyncSetTplToInstOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(option.SetIDs) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "bk_set_ids"))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.Logics.SetTemplateOperation().SyncSetTplToInst(ctx.Kit, bizID, setTemplateID, option); err != nil {
			blog.Errorf("SyncSetTplToInst failed, operation failed, bizID: %d, setTemplateID: %d, "+
				"option: %+v err: %s, rid: %s", bizID, setTemplateID, option, err.Error(), ctx.Kit.Rid)
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

// GetSetSyncDetails search set sync details
func (s *Service) GetSetSyncDetails(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	option := metadata.SetSyncStatusOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if option.SetIDs == nil {
		filter := &metadata.QueryCondition{
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
			Condition: mapstr.MapStr(map[string]interface{}{
				common.BKAppIDField:         bizID,
				common.BKSetTemplateIDField: setTemplateID,
			}),
		}
		setInstanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKInnerObjIDSet, filter)
		if err != nil {
			blog.Errorf("GetSetSyncStatus failed, get template related set failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		setIDs := make([]int64, 0)
		for _, inst := range setInstanceResult.Info {
			setID, err := inst.Int64(common.BKSetIDField)
			if err != nil {
				blog.Errorf("GetSetSyncStatus failed, get template related set failed, err: %+v, rid: %s", err,
					ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParseDBFailed))
				return
			}
			setIDs = append(setIDs, setID)
		}
		option.SetIDs = setIDs
	}

	taskCond := metadata.ListAPITaskDetail{
		InstID: option.SetIDs,
		Fields: []string{
			common.BKStatusField,
			common.MetaDataSynchronizeFlagField,
			common.BKInstIDField,
			"detail.status",
			"detail.data.module_diff.bk_module_name",
			"detail.response.baseresp.errmsg",
		},
	}

	taskDetail, err := s.Logics.SetTemplateOperation().GetLatestSyncTaskDetail(ctx.Kit, taskCond)
	if err != nil {
		blog.Errorf("get the latest task detail failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(taskDetail)
	return
}

// ListSetTemplateSyncHistory TODO
func (s *Service) ListSetTemplateSyncHistory(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := new(metadata.ListSetTemplateSyncStatusOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	option.BizID = bizID

	result, err := s.Logics.SetTemplateOperation().ListSetTemplateSyncHistory(ctx.Kit, option)
	if err != nil {
		blog.Errorf("list set template sync history failed, option: %#v, err: %v, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
	return
}

// ListSetTemplateSyncStatus TODO
func (s *Service) ListSetTemplateSyncStatus(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := new(metadata.ListSetTemplateSyncStatusOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	option.BizID = bizID

	result, err := s.Logics.SetTemplateOperation().ListSetTemplateSyncStatus(ctx.Kit, option)
	if err != nil {
		blog.Errorf("list set template sync status failed, option: %#v, err: %v, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
	return
}

// CheckSetInstUpdateToDateStatus TODO
func (s *Service) CheckSetInstUpdateToDateStatus(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	result, err := s.Logics.SetTemplateOperation().CheckSetInstUpdateToDateStatus(ctx.Kit, bizID, setTemplateID)
	if err != nil {
		blog.ErrorJSON("CheckSetInstUpdateToDateStatus failed, call core implement failed, bizID: %d, setTemplateID: %d, err: %s, rid: %s", bizID, setTemplateID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
	return
}

// BatchCheckSetInstUpdateToDateStatus TODO
func (s *Service) BatchCheckSetInstUpdateToDateStatus(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.BatchCheckSetInstUpdateToDateStatusOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	batchResult := make([]metadata.SetTemplateUpdateToDateStatus, 0)
	for _, setTemplateID := range option.SetTemplateIDs {
		oneResult, err := s.Logics.SetTemplateOperation().CheckSetInstUpdateToDateStatus(ctx.Kit, bizID, setTemplateID)
		if err != nil {
			blog.ErrorJSON("BatchCheckSetInstUpdateToDateStatus failed, CheckSetInstUpdateToDateStatus failed, bizID: %d, setTemplateID: %d, err: %s, rid: %s", bizID, setTemplateID, err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		batchResult = append(batchResult, *oneResult)
	}
	ctx.RespEntity(batchResult)
}

// UpdateSetTemplateAttribute update set template attribute
func (s *Service) UpdateSetTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.UpdateSetTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.Engine.CoreAPI.CoreService().SetTemplate().UpdateSetTemplateAttribute(ctx.Kit.Ctx, ctx.Kit.Header,
			option); err != nil {
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

// DeleteSetTemplateAttribute delete set template attribute
func (s *Service) DeleteSetTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.DeleteSetTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.Engine.CoreAPI.CoreService().SetTemplate().DeleteSetTemplateAttribute(ctx.Kit.Ctx, ctx.Kit.Header,
			option); err != nil {
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

// ListSetTemplateAttribute list set template attribute
func (s *Service) ListSetTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.ListSetTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	data, err := s.Engine.CoreAPI.CoreService().SetTemplate().ListSetTemplateAttribute(ctx.Kit.Ctx, ctx.Kit.Header,
		option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(data)
}
