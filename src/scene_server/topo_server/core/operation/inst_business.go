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

package operation

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// BusinessOperationInterface business operation methods
type BusinessOperationInterface interface {
	CreateBusiness(kit *rest.Kit, obj model.Object, data mapstr.MapStr) (inst.Inst, error)
	FindBiz(kit *rest.Kit, cond *metadata.QueryBusinessRequest) (count int, results []mapstr.MapStr, err error)
	GetInternalModule(kit *rest.Kit, bizID int64) (count int, result *metadata.InnterAppTopo, err errors.CCErrorCoder)
	UpdateBusiness(kit *rest.Kit, data mapstr.MapStr, obj model.Object, bizID int64) error
	HasHosts(kit *rest.Kit, bizID int64) (bool, error)
	SetProxy(set SetOperationInterface, module ModuleOperationInterface, inst InstOperationInterface, obj ObjectOperationInterface)
	GenerateAchieveBusinessName(kit *rest.Kit, bizName string) (achieveName string, err error)
	GetBriefTopologyNodeRelation(kit *rest.Kit, opts *metadata.GetBriefBizRelationOptions) ([]*metadata.BriefBizRelations, error)
}

// NewBusinessOperation create a business instance
func NewBusinessOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) BusinessOperationInterface {
	return &business{
		clientSet:   client,
		authManager: authManager,
	}
}

type business struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	inst        InstOperationInterface
	set         SetOperationInterface
	module      ModuleOperationInterface
	obj         ObjectOperationInterface
}

func (b *business) SetProxy(set SetOperationInterface, module ModuleOperationInterface, inst InstOperationInterface, obj ObjectOperationInterface) {
	b.inst = inst
	b.set = set
	b.module = module
	b.obj = obj
}

func (b *business) HasHosts(kit *rest.Kit, bizID int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		Fields:        []string{common.BKHostIDField},
		Page:          metadata.BasePage{Limit: 1},
	}
	rsp, err := b.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, option)
	if nil != err {
		blog.Errorf("[operation-set] failed to request the object controller, error info is %s", err.Error())
		return false, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return 0 != len(rsp.Info), nil
}

func (b *business) CreateBusiness(kit *rest.Kit, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {

	defaultFieldVal, err := data.Int64(common.BKDefaultField)
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is did not set the default field, %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}
	if defaultFieldVal == int64(common.DefaultAppFlag) && kit.SupplierAccount != common.BKDefaultOwnerID {
		// this is a new supplier owner and prepare to create a new business.
		asstQuery := map[string]interface{}{
			common.BKOwnerIDField: common.BKDefaultOwnerID,
		}
		defaultOwnerHeader := util.CloneHeader(kit.Header)
		defaultOwnerHeader.Set(common.BKHTTPOwnerID, common.BKDefaultOwnerID)

		asstRsp, err := b.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, defaultOwnerHeader, &metadata.QueryCondition{Condition: asstQuery})
		if nil != err {
			blog.Errorf("create business failed to get default assoc, error info is %s, rid: %s", err.Error(), kit.Rid)
			return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
		}

		expectAssts := asstRsp.Info
		blog.Infof("copy asst for %s, %+v, rid: %s", kit.SupplierAccount, expectAssts, kit.Rid)

		existAsstRsp, err := b.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: asstQuery})
		if nil != err {
			blog.Errorf("create business failed to get default assoc, error info is %s, rid: %s", err.Error(), kit.Rid)
			return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
		}

		existAssts := existAsstRsp.Info

	expectLoop:
		for _, asst := range expectAssts {
			asst.OwnerID = kit.SupplierAccount
			for _, existAsst := range existAssts {
				if existAsst.ObjectID == asst.ObjectID &&
					existAsst.AsstObjID == asst.AsstObjID &&
					existAsst.AsstKindID == asst.AsstKindID {
					continue expectLoop
				}
			}

			var err error
			if asst.AsstKindID == common.AssociationKindMainline {
				// bk_mainline is a inner association type that can only create in special case,
				// so we separate bk_mainline association type creation with a independent method,
				_, err = b.clientSet.CoreService().Association().CreateMainlineModelAssociation(kit.Ctx, kit.Header,
					&metadata.CreateModelAssociation{Spec: asst})
			} else {
				_, err = b.clientSet.CoreService().Association().CreateModelAssociation(kit.Ctx, kit.Header,
					&metadata.CreateModelAssociation{Spec: asst})
			}
			if nil != err {
				blog.Errorf("create business failed to copy default assoc, error info is %s, rid: %s", err.Error(), kit.Rid)
				return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
			}
		}
	}

	bizInst, err := b.inst.CreateInst(kit, obj, data)
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, err
	}

	bizID, err := bizInst.GetInstID()
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create set
	objSet, err := b.obj.FindSingleObject(kit, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	setData := mapstr.New()
	setData.Set(common.BKAppIDField, bizID)
	setData.Set(common.BKInstParentStr, bizID)
	setData.Set(common.BKSetNameField, common.DefaultResSetName)
	setData.Set(common.BKDefaultField, common.DefaultResSetFlag)

	setInst, err := b.set.CreateSet(kit, objSet, bizID, setData)
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	setID, err := setInst.GetInstID()
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create module
	objModule, err := b.obj.FindSingleObject(kit, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	defaultCategory, err := b.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx, kit.Header)
	if err != nil {
		blog.Errorf("failed to search default category, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.New(common.CCErrProcGetDefaultServiceCategoryFailed, err.Error())
	}

	idleModuleData := mapstr.New()
	idleModuleData.Set(common.BKSetIDField, setID)
	idleModuleData.Set(common.BKInstParentStr, setID)
	idleModuleData.Set(common.BKAppIDField, bizID)
	idleModuleData.Set(common.BKModuleNameField, common.DefaultResModuleName)
	idleModuleData.Set(common.BKDefaultField, common.DefaultResModuleFlag)
	idleModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	idleModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
	idleModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

	_, err = b.module.CreateModule(kit, objModule, bizID, setID, idleModuleData)
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create fault module
	faultModuleData := mapstr.New()
	faultModuleData.Set(common.BKSetIDField, setID)
	faultModuleData.Set(common.BKInstParentStr, setID)
	faultModuleData.Set(common.BKAppIDField, bizID)
	faultModuleData.Set(common.BKModuleNameField, common.DefaultFaultModuleName)
	faultModuleData.Set(common.BKDefaultField, common.DefaultFaultModuleFlag)
	faultModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	faultModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
	faultModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

	_, err = b.module.CreateModule(kit, objModule, bizID, setID, faultModuleData)
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create recycle module
	recycleModuleData := mapstr.New()
	recycleModuleData.Set(common.BKSetIDField, setID)
	recycleModuleData.Set(common.BKInstParentStr, setID)
	recycleModuleData.Set(common.BKAppIDField, bizID)
	recycleModuleData.Set(common.BKModuleNameField, common.DefaultRecycleModuleName)
	recycleModuleData.Set(common.BKDefaultField, common.DefaultRecycleModuleFlag)
	recycleModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	recycleModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
	recycleModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

	_, err = b.module.CreateModule(kit, objModule, bizID, setID, recycleModuleData)
	if nil != err {
		blog.Errorf("create business failed, create recycle module failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return bizInst, kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	return bizInst, nil
}

func (b *business) FindBusiness(kit *rest.Kit, cond *metadata.QueryBusinessRequest) (count int, results []mapstr.MapStr, err error) {

	cond.Condition[common.BKDefaultField] = 0
	query := &metadata.QueryCondition{
		Fields:    cond.Fields,
		Condition: cond.Condition,
		Page:      cond.Page,
	}

	result, err := b.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.ErrorJSON("failed to find business by query condition: %s, err: %s, rid: %s", query, err.Error(), kit.Rid)
		return 0, nil, err
	}

	return result.Count, result.Info, err
}
func (b *business) FindBiz(kit *rest.Kit, cond *metadata.QueryBusinessRequest) (count int, results []mapstr.MapStr, err error) {
	if !cond.Condition.Exists(common.BKDefaultField) {
		cond.Condition[common.BKDefaultField] = 0
	}
	query := &metadata.QueryCondition{
		Fields:    cond.Fields,
		Condition: cond.Condition,
		Page:      cond.Page,
	}

	result, err := b.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.ErrorJSON("failed to find business by query condition: %s, err: %s, rid: %s", query, err.Error(), kit.Rid)
		return 0, nil, err
	}

	return result.Count, result.Info, err
}

var (
	NumRegex = regexp.MustCompile(`^\d+$`)
)

/*
GenerateAchieveBusinessName 生成归档后的业务名称
	- 业务归档的时候，自动重命名为"foo-archived"
	- 归档的时候，如果发现已经存在同名的"foo-archived", 自动在其后+1, 比如 "foo-archived-1", "foo-archived-2"
*/
func (b *business) GenerateAchieveBusinessName(kit *rest.Kit, bizName string) (achieveName string, err error) {
	queryBusinessRequest := &metadata.QueryBusinessRequest{
		Fields: []string{common.BKAppNameField},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKAppNameField: map[string]interface{}{
				common.BKDBLIKE: fmt.Sprintf(`^%s-archived`, regexp.QuoteMeta(bizName)),
			},
		},
	}
	count, data, err := b.FindBusiness(kit, queryBusinessRequest)
	if err != nil {
		return "", err
	}
	if count == 0 {
		return fmt.Sprintf("%s-archived", bizName), nil
	}
	existNums := make([]int64, 0)
	for _, item := range data {
		biz := metadata.BizBasicInfo{}
		if err := mapstruct.Decode2Struct(item, &biz); err != nil {
			blog.Errorf("GenerateBusinessAchieveName failed, Decode2Struct failed, biz: %+v, err: %+v, rid: %s", item, err, kit.Rid)
			return "", kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
		}
		parts := strings.Split(biz.BizName, fmt.Sprintf("%s-archived-", bizName))
		if len(parts) != 2 {
			continue
		}
		numPart := parts[1]
		if !NumRegex.MatchString(numPart) {
			continue
		}
		num, err := util.GetInt64ByInterface(numPart)
		if err != nil {
			blog.Errorf("GenerateBusinessAchieveName failed, GetInt64ByInterface failed, numPart: %s, err: %+v, rid: %s", numPart, err, kit.Rid)
			return "", kit.CCError.CCError(common.CCErrCommParseDataFailed)
		}
		existNums = append(existNums, num)
	}
	// 空数组时默认填充
	existNums = append(existNums, 0)
	maxNum := existNums[0]
	for _, num := range existNums {
		if num > maxNum {
			maxNum = num
		}
	}

	return fmt.Sprintf("%s-archived-%d", bizName, maxNum+1), nil
}

func (b *business) GetInternalModule(kit *rest.Kit,
	bizID int64) (count int, result *metadata.InnterAppTopo, err errors.CCErrorCoder) {
	// get set model
	querySet := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKDefaultField: common.DefaultResSetFlag,
		},
		Fields: []string{common.BKSetIDField, common.BKSetNameField},
	}
	querySet.Page.Limit = 1

	/* setRsp, err := b.inst.FindOriginInst(kit, common.BKInnerObjIDSet, querySet)
	if nil != err {
		return 0, nil, kit.CCError.New(common.CCErrTopoAppSearchFailed, err.Error())
	} */

	setRsp := &metadata.ResponseSetInstance{}
	// 返回数据不包含自定义字段
	if err = b.clientSet.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header,
		common.BKInnerObjIDSet, querySet, setRsp); err != nil {
		return 0, nil, err
	}
	if err := setRsp.CCError(); err != nil {
		blog.ErrorJSON("query set error. filter: %s, result: %s, rid: %s", querySet, setRsp, kit.Rid)
		return 0, nil, err
	}

	// search modules
	queryModule := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKDefaultField: map[string]interface{}{
				common.BKDBNE: 0,
			},
		},
		Fields: []string{common.BKModuleIDField, common.BKModuleNameField, common.BKDefaultField, common.HostApplyEnabledField},
	}
	queryModule.Page.Limit = common.BKNoLimit

	/*
		moduleRsp, err := b.inst.FindOriginInst(kit, common.BKInnerObjIDModule, queryModule)
		if nil != err {
			return 0, nil, kit.CCError.New(common.CCErrTopoAppSearchFailed, err.Error())
		}
	*/

	moduleResp := &metadata.ResponseModuleInstance{}
	// 返回数据不包含自定义字段
	if err = b.clientSet.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header,
		common.BKInnerObjIDModule, queryModule, moduleResp); err != nil {
		return 0, nil, err
	}
	if err := moduleResp.CCError(); err != nil {
		blog.ErrorJSON("query module error. filter: %s, result: %s, rid: %s", queryModule, moduleResp, kit.Rid)
		return 0, nil, err
	}

	// construct result
	result = &metadata.InnterAppTopo{}
	for _, set := range setRsp.Data.Info {
		result.SetID = set.SetID
		result.SetName = set.SetName
		break // should be only one set
	}

	for _, module := range moduleResp.Data.Info {

		result.Module = append(result.Module, metadata.InnerModule{
			ModuleID:         module.ModuleID,
			ModuleName:       module.ModuleName,
			Default:          module.Default,
			HostApplyEnabled: module.HostApplyEnabled,
		})
	}

	return 0, result, nil
}

func (b *business) UpdateBusiness(kit *rest.Kit, data mapstr.MapStr, obj model.Object, bizID int64) error {
	innerCond := condition.CreateCondition()
	innerCond.Field(common.BKAppIDField).Eq(bizID)

	return b.inst.UpdateInst(kit, data, obj, innerCond, bizID)
}

// GetBriefTopologyNodeRelation is used to get directly related business topology node information.
// As is, you can find modules belongs to a set; or you can find the set a module belongs to.
// It has rules as follows:
// 1. if src object is biz, then the destination object can be any mainline object except biz.
// 2. destination object can be biz. otherwise, src and destination object should be the neighbour.
// this api only return business topology relations.
func (b *business) GetBriefTopologyNodeRelation(kit *rest.Kit, opts *metadata.GetBriefBizRelationOptions) (
	[]*metadata.BriefBizRelations, error) {

	// validate the source and destination model is mainline model or not.
	srcDestPriority, err := b.validateMainlineObjectRule(kit, opts.SrcBizObj, opts.DestBizObj)
	if err != nil {
		blog.Errorf("check object is mainline object failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "src_inst_ids or dest_biz_obj")
	}

	filter := make(mapstr.MapStr)
	switch opts.SrcBizObj {
	case common.BKInnerObjIDApp:
		filter[common.BKAppIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
		return b.genBriefTopologyNodeRelation(kit, filter, opts.DestBizObj, common.BKAppIDField,
			common.GetInstIDField(opts.DestBizObj), &opts.Page)

	case common.BKInnerObjIDSet:
		switch opts.DestBizObj {
		case common.BKInnerObjIDApp:
			filter[common.BKSetIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDSet, common.BKSetIDField,
				common.BKAppIDField, &opts.Page)

		case common.BKInnerObjIDModule:
			filter[common.BKSetIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDModule, common.BKSetIDField,
				common.BKModuleIDField, &opts.Page)

		default:
			// search custom level model instance with set ids. which is set's parent id list
			filter[common.BKSetIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDSet, common.BKSetIDField,
				common.BKParentIDField, &opts.Page)
		}

	case common.BKInnerObjIDModule:
		switch opts.DestBizObj {
		case common.BKInnerObjIDApp:
			filter[common.BKModuleIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDModule, common.BKModuleIDField,
				common.BKAppIDField, &opts.Page)

		case common.BKInnerObjIDSet:
			filter[common.BKModuleIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDModule, common.BKModuleIDField,
				common.BKSetIDField, &opts.Page)

		default:
			blog.Errorf("it's not allow to find destination object %s with source module model. rid: %s",
				opts.DestBizObj, kit.Rid)
			return nil, errors.New(common.CCErrCommParamsInvalid, "dest_biz_obj")
		}

	default:
		switch opts.DestBizObj {
		case common.BKInnerObjIDApp:
			filter[common.BKInstIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, opts.SrcBizObj, common.BKInstIDField,
				common.BKAppIDField, &opts.Page)

		case common.BKInnerObjIDSet:
			filter[common.BKParentIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, common.BKInnerObjIDSet, common.BKParentIDField,
				common.BKSetIDField, &opts.Page)

		default:
			if srcDestPriority < 0 {
				// src object is the parent of destination object
				filter[common.BKParentIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
				// src inst id bk_parent_id, dest inst id bk_inst_id, this scene is different with others.
				return b.genBriefTopologyNodeRelation(kit, filter, opts.DestBizObj, common.BKParentIDField,
					common.BKInstIDField, &opts.Page)

			}

			// destination object is the parent of src object
			filter[common.BKInstIDField] = mapstr.MapStr{common.BKDBIN: opts.SrcInstIDs}
			return b.genBriefTopologyNodeRelation(kit, filter, opts.SrcBizObj, common.BKInstIDField,
				common.BKParentIDField, &opts.Page)
		}
	}
}

func (b *business) genBriefTopologyNodeRelation(kit *rest.Kit, filter mapstr.MapStr, destObj, srcInstField,
	destInstField string, page *metadata.BasePage) ([]*metadata.BriefBizRelations, error) {

	// set sort field with destination object instance field as default.
	page.Sort = common.GetInstIDField(destObj)

	input := &metadata.QueryCondition{
		// set all the possible fields.
		Fields: []string{common.BKAppIDField, common.BKParentIDField, common.BKSetIDField, common.BKModuleIDField,
			common.BKInstIDField},
		Page:           *page,
		Condition:      filter,
		DisableCounter: true,
	}
	result, err := b.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, destObj, input)
	if err != nil {
		blog.ErrorJSON("get biz mainline object %s instance with filter: %s failed, err: %v, rid: %s",
			destObj, input, err, kit.Rid)
		return nil, err
	}

	relations := make([]*metadata.BriefBizRelations, 0)
	for _, one := range result.Info {
		relations = append(relations, &metadata.BriefBizRelations{
			Business: one[common.BKAppIDField],
			// source object's instance id, field different with object's type
			SrcInstID: one[srcInstField],
			// destination object's instance id, field different with object's type
			DestInstID: one[destInstField],
		})
	}

	return relations, nil
}

func (b *business) validateMainlineObjectRule(kit *rest.Kit, src, dest string) (int, error) {
	cond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}
	asst, err := b.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return 0, err
	}

	if len(asst.Info) <= 0 {
		return 0, fmt.Errorf("invalid biz mainline object topology")
	}

	next, idx := common.BKInnerObjIDApp, 0
	// save the mainline object with it's index with map
	rankMap := make(map[string]int)
	rankMap[next] = idx
	for _, relation := range asst.Info {
		if relation.AsstObjID == next {
			next = relation.ObjectID
			idx += 1
			rankMap[next] = idx
			continue
		}

		for _, rel := range asst.Info {
			if rel.AsstObjID == next {
				next = rel.ObjectID
				idx += 1
				rankMap[next] = idx
				break
			}
		}

	}

	// src, dest object should all be mainline object.
	srcIdx, destIdx := 0, 0
	srcIdx, exist := rankMap[src]
	if !exist {
		return 0, fmt.Errorf("%s is not mainline object", src)
	}

	destIdx, exist = rankMap[dest]
	if !exist {
		return 0, fmt.Errorf("%s is not mainline object", dest)
	}

	srcDestPriority := srcIdx - destIdx

	if src == common.BKInnerObjIDApp {
		// if src object is biz, then do not care about if the destination object is neighbour or not.
		return srcDestPriority, nil
	}

	// if dest object is not biz, then the src and dest object should be the neighbour.
	// if dest object is biz, we do not check the src or dest is neighbour or not.
	if (dest != common.BKInnerObjIDApp) && (math.Abs(float64(srcDestPriority)) > 1) {
		return 0, fmt.Errorf("src[%s] model and dest[%s] model should be neighbour in the mainline topology", src, dest)
	}

	return srcDestPriority, nil
}
