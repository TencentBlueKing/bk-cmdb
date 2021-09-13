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

package inst

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// BusinessOperationInterface business operation methods
type BusinessOperationInterface interface {
	// CreateBusiness create business
	CreateBusiness(kit *rest.Kit, data mapstr.MapStr) (mapstr.MapStr, error)
	// FindBiz find biz
	FindBiz(kit *rest.Kit, cond *metadata.QueryCondition) (count int,
		results []mapstr.MapStr, err error)
	// HasHosts check if this business still has hosts.
	HasHosts(kit *rest.Kit, bizID int64) (bool, error)
	// GenerateAchieveBusinessName 生成归档后的业务名称
	// - 业务归档的时候，自动重命名为"foo-archived"
	// - 归档的时候，如果发现已经存在同名的"foo-archived", 自动在其后+1, 比如 "foo-archived-1", "foo-archived-2"
	GenerateAchieveBusinessName(kit *rest.Kit, bizName string) (achieveName string, err error)
	// GetBriefTopologyNodeRelation is used to get directly related business topology node information.
	// As is, you can find modules belongs to a set; or you can find the set a module belongs to.
	// It has rules as follows:
	// 1. if src object is biz, then the destination object can be any mainline object except biz.
	// 2. destination object can be biz. otherwise, src and destination object should be the neighbour.
	// this api only return business topology relations.
	GetBriefTopologyNodeRelation(kit *rest.Kit, opts *metadata.GetBriefBizRelationOptions) ([]*metadata.
		BriefBizRelations, error)
	// SetProxy SetProxy
	SetProxy(inst InstOperationInterface, module ModuleOperationInterface, set SetOperationInterface)
}

// NewBusinessOperation create a business instance
func NewBusinessOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) BusinessOperationInterface {
	return &business{
		clientSet:   client,
		authManager: authManager,
	}
}

type business struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	inst        InstOperationInterface
	module      ModuleOperationInterface
	set         SetOperationInterface
}

// SetProxy SetProxy
func (b *business) SetProxy(inst InstOperationInterface, module ModuleOperationInterface, set SetOperationInterface) {
	b.inst = inst
	b.module = module
	b.set = set
}

// CreateBusiness create business
func (b *business) CreateBusiness(kit *rest.Kit, data mapstr.MapStr) (mapstr.MapStr, error) {

	// this is a new supplier owner and prepare to create a new business.
	if err := b.createAssociationByNewSupplier(kit, data); err != nil {
		blog.Errorf("create association for new default biz of different supplier account, err: %v, data: %#v, "+
			"rid: %s", err, data, kit.Rid)
		return nil, err
	}

	bizInst, err := b.inst.CreateInst(kit, common.BKInnerObjIDApp, data)
	if err != nil {
		blog.Errorf("create business failed, err: %v, data: %#v, rid: %s", err, data, kit.Rid)
		return nil, err
	}
	bizID, err := bizInst.Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("create business failed, err: %v, data: %#v, rid: %s", err, data, kit.Rid)
		return nil, err
	}

	// create set
	setData := mapstr.MapStr{
		common.BKAppIDField:    bizID,
		common.BKInstParentStr: bizID,
		common.BKSetNameField:  common.DefaultResSetName,
		common.BKDefaultField:  common.DefaultResSetFlag,
	}
	setInst, err := b.set.CreateSet(kit, bizID, setData)
	if err != nil {
		blog.Errorf("create set failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	setID, err := setInst.Int64(common.BKSetIDField)
	if err != nil {
		blog.Errorf("create set failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	// create module
	defaultCategory, err := b.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx, kit.Header)
	if err != nil {
		blog.Errorf("failed to search default category, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	moduleData := mapstr.MapStr{
		common.BKSetIDField:             setID,
		common.BKInstParentStr:          setID,
		common.BKAppIDField:             bizID,
		common.BKModuleNameField:        common.DefaultResModuleName,
		common.BKDefaultField:           common.DefaultResModuleFlag,
		common.BKServiceTemplateIDField: common.ServiceTemplateIDNotSet,
		common.BKSetTemplateIDField:     common.SetTemplateIDNotSet,
		common.BKServiceCategoryIDField: defaultCategory.ID,
	}

	if _, err = b.module.CreateModule(kit, bizID, setID, moduleData); err != nil {
		blog.Errorf("create module failed, err: %v, rid: %s", err, kit.Rid)
		return data, err
	}

	// create fault module
	moduleData.Set(common.BKModuleNameField, common.DefaultFaultModuleName)
	moduleData.Set(common.BKDefaultField, common.DefaultFaultModuleFlag)
	if _, err = b.module.CreateModule(kit, bizID, setID, moduleData); err != nil {
		blog.Errorf("create fault module failed, err: %v, rid: %s", err, kit.Rid)
		return data, err
	}

	// create recycle module
	moduleData.Set(common.BKModuleNameField, common.DefaultRecycleModuleName)
	moduleData.Set(common.BKDefaultField, common.DefaultRecycleModuleFlag)
	if _, err = b.module.CreateModule(kit, bizID, setID, moduleData); err != nil {
		blog.Errorf("create recycle module failed, err: %v, rid: %s", err, kit.Rid)
		return data, err
	}

	return bizInst, nil
}

// FindBiz FindBiz
func (b *business) FindBiz(kit *rest.Kit, cond *metadata.QueryCondition) (count int,
	results []mapstr.MapStr, err error) {
	if !cond.Condition.Exists(common.BKDefaultField) {
		cond.Condition[common.BKDefaultField] = 0
	}

	result, err := b.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("find business by query failed, condition: %s, err: %v, rid: %s", cond, err, kit.Rid)
		return 0, nil, err
	}

	return result.Count, result.Info, err
}

// HasHosts check if this business still has hosts.
func (b *business) HasHosts(kit *rest.Kit, bizID int64) (bool, error) {

	option := []map[string]interface{}{{
		common.BKAppIDField: bizID,
	}}

	rsp, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameModuleHostConfig, option)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, kit.Rid)
		return false, err
	}

	return rsp[0] != 0, nil
}

var (
	numRegex = regexp.MustCompile(`^\d+$`)
)

// GenerateAchieveBusinessName 生成归档后的业务名称
// - 业务归档的时候，自动重命名为"foo-archived"
// - 归档的时候，如果发现已经存在同名的"foo-archived", 自动在其后+1, 比如 "foo-archived-1", "foo-archived-2"
func (b *business) GenerateAchieveBusinessName(kit *rest.Kit, bizName string) (achieveName string, err error) {

	queryBusinessRequest := &metadata.QueryCondition{
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
	count, data, err := b.FindBiz(kit, queryBusinessRequest)
	if err != nil {
		return "", err
	}
	if count == 0 {
		return fmt.Sprintf("%s-archived", bizName), nil
	}
	existNums := make([]int64, 0)
	for _, item := range data {
		parts := strings.Split(util.GetStrByInterface(item[common.BKAppNameField]), fmt.Sprintf("%s-archived-",
			bizName))
		if len(parts) != 2 {
			continue
		}
		numPart := parts[1]
		if !numRegex.MatchString(numPart) {
			continue
		}
		num, err := util.GetInt64ByInterface(numPart)
		if err != nil {
			blog.Errorf("GetInt64ByInterface failed, numPart: %s, err: %v, rid: %s", numPart, err, kit.Rid)
			return "", err
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
		return nil, err
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
		return 0, kit.CCError.New(common.CCErrorTopoMainlineObjectAssociationNotExist, src)
	}

	destIdx, exist = rankMap[dest]
	if !exist {
		return 0, kit.CCError.New(common.CCErrorTopoMainlineObjectAssociationNotExist, dest)
	}

	srcDestPriority := srcIdx - destIdx

	if src == common.BKInnerObjIDApp {
		// if src object is biz, then do not care about if the destination object is neighbour or not.
		return srcDestPriority, nil
	}

	// if dest object is not biz, then the src and dest object should be the neighbour.
	// if dest object is biz, we do not check the src or dest is neighbour or not.
	if (dest != common.BKInnerObjIDApp) && (math.Abs(float64(srcDestPriority)) > 1) {
		return 0, kit.CCError.New(common.CCErrCommTopoModuleNotFoundError, src)
	}

	return srcDestPriority, nil
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

func (b *business) createAssociationByNewSupplier(kit *rest.Kit, data mapstr.MapStr) error {
	defaultFieldVal, err := data.Int64(common.BKDefaultField)
	if err != nil {
		blog.Errorf("failed to create business, error info is did not set the default field,err: %v, rid: %s",
			err, kit.Rid)
		return err
	}
	if defaultFieldVal != int64(common.DefaultAppFlag) || kit.SupplierAccount == common.BKDefaultOwnerID {
		return nil
	}

	asstQuery := map[string]interface{}{
		common.BKOwnerIDField: common.BKDefaultOwnerID,
	}
	defaultOwnerHeader := util.CloneHeader(kit.Header)
	defaultOwnerHeader.Set(common.BKHTTPOwnerID, common.BKDefaultOwnerID)

	asstRsp, err := b.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, defaultOwnerHeader,
		&metadata.QueryCondition{Condition: asstQuery})
	if err != nil {
		blog.Errorf("create business failed to get default assoc, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	expectAssts := asstRsp.Info
	blog.Infof("copy asst for %s, %+v, rid: %s", kit.SupplierAccount, expectAssts, kit.Rid)

	existAsstRsp, err := b.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: asstQuery})
	if err != nil {
		blog.Errorf("create business failed to get default assoc, err: %v, rid: %v", err, kit.Rid)
		return err
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
		if err != nil {
			blog.Errorf("create business failed to copy default assoc, err: %v, rid: %s", err, kit.Rid)
			return err
		}
	}

	return nil
}
