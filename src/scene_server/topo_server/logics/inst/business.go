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
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/logics/model"
)

// BusinessOperationInterface business operation methods
type BusinessOperationInterface interface {
	// CreateBusiness create business
	CreateBusiness(kit *rest.Kit, data mapstr.MapStr) (mapstr.MapStr, error)
	// FindBiz find biz
	FindBiz(kit *rest.Kit, cond *metadata.QueryBusinessRequest, disableCounter bool) (count int,
		results []mapstr.MapStr, err error)
	// UpdateBusiness update business
	UpdateBusiness(kit *rest.Kit, data mapstr.MapStr, obj metadata.Object, bizID int64) error
	// HasHosts check if this business still has hosts.
	HasHosts(kit *rest.Kit, bizID int64) (bool, error)
	//	GenerateAchieveBusinessName 生成归档后的业务名称
	//	- 业务归档的时候，自动重命名为"foo-archived"
	//	- 归档的时候，如果发现已经存在同名的"foo-archived", 自动在其后+1, 比如 "foo-archived-1", "foo-archived-2"
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
	SetProxy(obj model.ObjectOperationInterface)
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
	obj         model.ObjectOperationInterface
}

var (
	numRegex = regexp.MustCompile(`^\d+$`)
)

// SetProxy SetProxy
func (b *business) SetProxy(obj model.ObjectOperationInterface) {
	b.obj = obj
}

// CreateBusiness create business
func (b *business) CreateBusiness(kit *rest.Kit, data mapstr.MapStr) (mapstr.MapStr, error) {
	// this is a new supplier owner and prepare to create a new business.
	if err := b.createAssociationByNewSupplier(kit, data); err != nil {
		blog.Errorf("create business failed, err: %v, data: %#v, rid: %s", err, data, kit.Rid)
		return nil, err
	}
	// TODO 调用inst中createInst
	bizInst, err := b.createInst(kit, metadata.Object{ObjectID: common.BKInnerObjIDApp}, data)
	if err != nil {
		blog.Errorf("create business failed, err: %v, data: %#v, rid: %s", err, data, kit.Rid)
		return nil, err
	}
	bizID := int64(bizInst.Created.ID)

	// create set
	objSet, err := b.obj.FindSingleObject(kit, nil, common.BKInnerObjIDSet)
	if err != nil {
		blog.Errorf("failed to search the set, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	setData := mapstr.MapStr{
		common.BKAppIDField:    bizID,
		common.BKInstParentStr: bizID,
		common.BKSetNameField:  common.DefaultResSetName,
		common.BKDefaultField:  common.DefaultResSetFlag,
	}
	// TODO 调用set中CreateSet
	setInst, err := b.createSet(kit, *objSet, setData)
	if err != nil {
		blog.Errorf("create business failed, err: %s, rid: %s", err, kit.Rid)
		return nil, err
	}
	setID := int64(setInst.Created.ID)

	// create module
	objModule, err := b.obj.FindSingleObject(kit, nil, common.BKInnerObjIDModule)
	if err != nil {
		blog.Errorf("failed to search the set, %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	defaultCategory, err := b.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx, kit.Header)
	if err != nil {
		blog.Errorf("failed to search default category, err: %+v, rid: %s", err, kit.Rid)
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
	// TODO 调用module中CreateModule
	if _, err = b.createModule(kit, *objModule, bizID, setID, moduleData); err != nil {
		blog.Errorf("create business failed, err: %v, rid: %s", err, kit.Rid)
		return data, err
	}
	// create fault module
	moduleData.Set(common.BKModuleNameField, common.DefaultFaultModuleName)
	moduleData.Set(common.BKDefaultField, common.DefaultFaultModuleFlag)
	// TODO 调用module中CreateModule
	if _, err = b.createModule(kit, *objModule, bizID, setID, moduleData); err != nil {
		blog.Errorf("create business failed, err: %v, rid: %s", err, kit.Rid)
		return data, err
	}
	// create recycle module
	moduleData.Set(common.BKModuleNameField, common.DefaultRecycleModuleName)
	moduleData.Set(common.BKDefaultField, common.DefaultRecycleModuleFlag)
	// TODO 调用module中CreateModule
	if _, err = b.createModule(kit, *objModule, bizID, setID, moduleData); err != nil {
		blog.Errorf("create recycle module failed, err: %v, rid: %s", err, kit.Rid)
		return data, err
	}
	return data, nil
}

// FindBiz FindBiz
func (b *business) FindBiz(kit *rest.Kit, cond *metadata.QueryBusinessRequest, disableCounter bool) (count int,
	results []mapstr.MapStr, err error) {
	if !cond.Condition.Exists(common.BKDefaultField) {
		cond.Condition[common.BKDefaultField] = 0
	}
	query := &metadata.QueryCondition{
		Fields:         cond.Fields,
		Condition:      cond.Condition,
		Page:           cond.Page,
		DisableCounter: disableCounter,
	}

	result, err := b.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("find business by query failed, condition: %s, err: %v, rid: %s", query, err, kit.Rid)
		return 0, nil, err
	}

	return result.Count, result.Info, err
}

// UpdateBusiness update business
func (b *business) UpdateBusiness(kit *rest.Kit, data mapstr.MapStr, obj metadata.Object, bizID int64) error {
	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}

	// TODO 这里调用 inst.go 中UpdateInst ,将在合并后修改
	return b.updateInst(kit, cond, data, obj.ObjectID)
}

// HasHosts check if this business still has hosts.
func (b *business) HasHosts(kit *rest.Kit, bizID int64) (bool, error) {
	option := new(metadata.HostModuleRelationRequest)
	option.ApplicationID = bizID
	option.Fields = []string{common.BKHostIDField}
	option.Page = metadata.BasePage{Limit: 1}

	rsp, err := b.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, kit.Rid)
		return false, err
	}

	return len(rsp.Info) != 0, nil
}

//	GenerateAchieveBusinessName 生成归档后的业务名称
//	- 业务归档的时候，自动重命名为"foo-archived"
//	- 归档的时候，如果发现已经存在同名的"foo-archived", 自动在其后+1, 比如 "foo-archived-1", "foo-archived-2"
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
	count, data, err := b.FindBiz(kit, queryBusinessRequest, false)
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
		return 0, fmt.Errorf("src[%s] model and dest[%s] model should be neighbour in the mainline topology",
			src, dest)
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

// TODO 以下为临时函数
func (b *business) findSingleObject(kit *rest.Kit, objectID string) (*metadata.Object, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objectID)

	objs, err := b.findObject(kit, cond)
	if err != nil {
		blog.Errorf("get model failed, failed to get model by supplier account(%s) objects(%s), err: %s, "+
			"rid: %s", kit.SupplierAccount, objectID, err.Error(), kit.Rid)
		return nil, err
	}

	if len(objs) == 0 {
		blog.Errorf("get model failed, get model by supplier account(%s) objects(%s) not found, result: %+v, "+
			"rid: %s", kit.SupplierAccount, objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed, kit.CCError.Error(common.CCErrCommNotFound).
			Error())
	}
	if len(objs) > 1 {
		blog.Errorf("get model failed, get model by supplier account(%s) objects(%s) get multiple, result:"+
			" %+v, rid: %s", kit.SupplierAccount, objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed,
			kit.CCError.Error(common.CCErrCommGetMultipleObject).Error())
	}

	objects := make([]metadata.Object, 0)
	for _, obj := range objs {
		objects = append(objects, obj)
	}

	for _, item := range objs {
		return &item, nil
	}
	return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed,
		kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, objectID).Error())
}

func (b *business) updateInst(kit *rest.Kit, cond, data mapstr.MapStr, objID string) error {
	// not allowed to update these fields, need to use specialized function
	data.Remove(common.BKParentIDField)
	data.Remove(common.BKAppIDField)

	inputParams := metadata.UpdateOption{
		Data:      data,
		Condition: cond,
	}

	// generate audit log of instance.
	audit := auditlog.NewInstanceAudit(b.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(data)
	auditLog, ccErr := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, cond)
	if ccErr != nil {
		blog.Errorf(" update inst, generate audit log failed, err: %v, rid: %s", ccErr, kit.Rid)
		return ccErr
	}

	// to update.
	_, err := b.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, objID, &inputParams)
	if err != nil {
		blog.Errorf("update instance failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	// save audit log.
	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return nil
}

func (b *business) createSet(kit *rest.Kit, obj metadata.Object, data mapstr.MapStr) (*metadata.CreateOneDataResult,
	error) {
	if obj.ObjectID == common.BKInnerObjIDPlat {
		data.Set(common.BkSupplierAccount, kit.SupplierAccount)
	}

	data.Set(common.BKObjIDField, obj.ObjectID)

	instCond := &metadata.CreateModelInstance{Data: data}
	rsp, err := b.clientSet.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, obj.ObjectID, instCond)
	if err != nil {
		blog.Errorf("failed to create object instance, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if rsp.Created.ID == 0 {
		blog.Errorf("failed to create object instance, return nothing, rid: %s", kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoInstCreateFailed)
	}

	data.Set(obj.GetInstIDFieldName(), rsp.Created.ID)
	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewInstanceAudit(b.clientSet.CoreService())
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.GetObjectID(), []mapstr.MapStr{data})
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return rsp, nil
}

func (b *business) createInst(kit *rest.Kit, obj metadata.Object, data mapstr.MapStr) (*metadata.CreateOneDataResult,
	error) {

	if obj.ObjectID == common.BKInnerObjIDPlat {
		data.Set(common.BkSupplierAccount, kit.SupplierAccount)
	}

	data.Set(common.BKObjIDField, obj.ObjectID)

	instCond := &metadata.CreateModelInstance{Data: data}
	rsp, err := b.clientSet.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, obj.ObjectID, instCond)
	if err != nil {
		blog.Errorf("failed to create object instance, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if rsp.Created.ID == 0 {
		blog.Errorf("failed to create object instance, return nothing, rid: %s", kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoInstCreateFailed)
	}

	data.Set(obj.GetInstIDFieldName(), rsp.Created.ID)
	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewInstanceAudit(b.clientSet.CoreService())
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.GetObjectID(), []mapstr.MapStr{data})
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return rsp, nil
}

func (b *business) findObject(kit *rest.Kit, cond condition.Condition) ([]metadata.Object, error) {
	fCond := cond.ToMapStr()

	rsp, err := b.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: fCond})
	if err != nil {
		blog.Errorf("find object failed, cond: %+v, err: %s, rid: %s", fCond, err.Error(),
			kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return rsp.Info, nil
}

// TODO 这个函数后续会调用module CreateModule，这里先这么写
func (b *business) createModule(kit *rest.Kit, obj metadata.Object, bizID, setID int64,
	data mapstr.MapStr) (*metadata.CreateOneDataResult, error) {

	data.Set(common.BKSetIDField, setID)
	data.Set(common.BKAppIDField, bizID)
	if !data.Exists(common.BKDefaultField) {
		data.Set(common.BKDefaultField, common.DefaultFlagDefaultValue)
	}
	_, err := data.Int64(common.BKDefaultField)
	if err != nil {
		blog.Errorf("parse default field into int failed, data: %+v, rid: %s", data, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKDefaultField)
		return nil, err
	}

	// validate service category id and service template id
	// 如果服务分类没有设置，则从服务模版中获取，如果服务模版也没有设置，则参数错误
	// 有效参数参数形式:
	// 1. serviceCategoryID > 0  && serviceTemplateID == 0
	// 2. serviceCategoryID unset && serviceTemplateID > 0
	// 3. serviceCategoryID > 0 && serviceTemplateID > 0 && serviceTemplate.ServiceCategoryID == serviceCategoryID
	// 4. serviceCategoryID unset && serviceTemplateID unset, then module create with default category
	var serviceCategoryID int64
	serviceCategoryIDIf, serviceCategoryExist := data.Get(common.BKServiceCategoryIDField)
	if serviceCategoryExist == true {
		scID, err := util.GetInt64ByInterface(serviceCategoryIDIf)
		if err != nil {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
		serviceCategoryID = scID
	}

	var serviceTemplateID int64
	serviceTemplateIDIf, serviceTemplateFieldExist := data.Get(common.BKServiceTemplateIDField)
	if serviceTemplateFieldExist == true {
		serviceTemplateID, err = util.GetInt64ByInterface(serviceTemplateIDIf)
		if err != nil {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
	}

	if serviceCategoryID == 0 && serviceTemplateID == 0 {
		// set default service template id
		defaultServiceCategory, err := b.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx,
			kit.Header)
		if err != nil {
			blog.Errorf("create module failed, GetDefaultServiceCategory failed, err: %s, rid: %s", err.Error(),
				kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrProcGetDefaultServiceCategoryFailed)
		}
		serviceCategoryID = defaultServiceCategory.ID
	} else if serviceTemplateID != common.ServiceTemplateIDNotSet {
		// 校验 serviceCategoryID 与 serviceTemplateID 对应
		templateIDs := []int64{serviceTemplateID}
		option := metadata.ListServiceTemplateOption{
			BusinessID:         bizID,
			ServiceTemplateIDs: templateIDs,
		}
		stResult, err := b.clientSet.CoreService().Process().ListServiceTemplates(kit.Ctx, kit.Header, &option)
		if err != nil {
			return nil, err
		}
		if len(stResult.Info) == 0 {
			blog.ErrorJSON("create module failed, service template not found, filter: %s, rid: %s", option,
				kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		if serviceCategoryExist == true && serviceCategoryID != stResult.Info[0].ServiceCategoryID {
			return nil, kit.CCError.Error(common.CCErrProcServiceTemplateAndCategoryNotCoincide)
		}
		serviceCategoryID = stResult.Info[0].ServiceCategoryID
	} else {
		// 检查 service category id 是否有效
		serviceCategory, err := b.clientSet.CoreService().Process().GetServiceCategory(kit.Ctx, kit.Header,
			serviceCategoryID)
		if err != nil {
			return nil, err
		}
		if serviceCategory.BizID != 0 && serviceCategory.BizID != bizID {
			blog.V(3).Info("create module failed, service category and module belong to two business, "+
				"categoryBizID: %d, bizID: %d, rid: %s", serviceCategory.BizID, bizID, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
	}
	data.Set(common.BKServiceCategoryIDField, serviceCategoryID)
	data.Set(common.BKServiceTemplateIDField, serviceTemplateID)
	data.Set(common.HostApplyEnabledField, false)

	// set default set template
	_, exist := data[common.BKSetTemplateIDField]
	if exist == false {
		data[common.BKSetTemplateIDField] = common.SetTemplateIDNotSet
	}

	// convert bk_parent_id to int
	parentIDIf, ok := data[common.BKParentIDField]
	if ok == true {
		parentID, err := util.GetInt64ByInterface(parentIDIf)
		if err != nil {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKParentIDField)
		}
		data[common.BKParentIDField] = parentID
	}

	data.Remove(common.MetadataField)
	inst, err := b.createInst(kit, obj, data)
	if err != nil {
		return inst, err
	}

	return inst, nil
}
