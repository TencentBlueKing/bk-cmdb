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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// DeleteBusiness delete business instances by condition
func (b *business) DeleteBusiness(kit *rest.Kit, bizIDs []int64) error {
	// pre-check
	if err := b.checkCanDelete(kit, bizIDs); err != nil {
		return err
	}

	// clean business and related resources
	for _, bizID := range bizIDs {
		if err := b.cleanBizAndRelatedResources(kit, bizID); err != nil {
			return err
		}
	}

	return nil
}

func (b *business) checkCanDelete(kit *rest.Kit, bizIDs []int64) error {
	// 1. check built-in business
	has, err := b.checkHasBuiltInBiz(kit, bizIDs)
	if err != nil {
		return err
	}
	if has {
		blog.Errorf("forbidden delete built in business, rid: %s", kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoForbiddenDeleteBuiltInBiz)
	}

	// 2. check unarchived business
	has, err = b.checkHasUnarchivedBiz(kit, bizIDs)
	if err != nil {
		return err
	}
	if has {
		blog.Errorf("forbidden delete unarchived business, rid: %s", kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoForbiddenDeleteUnarchivedBiz)
	}

	// 3. check business has hosts
	has, err = b.checkHasHost(kit, bizIDs)
	if err != nil {
		return err
	}
	if has {
		blog.Errorf("forbidden delete business with hosts, rid: %s", kit.Rid)
		return kit.CCError.CCError(common.CCErrTopoHasHost)
	}

	return nil
}

// checkHasBuiltInBiz check if business list has built-in business
func (b *business) checkHasBuiltInBiz(kit *rest.Kit, bizIDs []int64) (bool, error) {
	// get biz count by filter
	filter := []map[string]interface{}{{
		// built in business
		common.BKDefaultField: common.DefaultAppFlag,
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: bizIDs,
		},
	}}

	rst, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameBaseApp,
		filter)
	if err != nil {
		blog.Errorf("get biz count failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return false, err
	}

	if len(rst) != 1 {
		blog.Errorf("get biz count failed, for result len must be 1, filter: %+v, rid: %s", filter, kit.Rid)
		return false, kit.CCError.Error(common.CCErrOperationBizModuleHostAmountFail)
	}

	if rst[0] <= 0 {
		return false, nil
	}

	return true, nil
}

// checkHasUnarchivedBiz check if business list has unarchived business
func (b *business) checkHasUnarchivedBiz(kit *rest.Kit, bizIDs []int64) (bool, error) {
	// get biz count by filter
	filter := []map[string]interface{}{{
		// unarchived business
		common.BKDataStatusField: map[string]interface{}{
			common.BKDBNE: common.DataStatusDisabled,
		},
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: bizIDs,
		},
	}}

	rst, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameBaseApp,
		filter)
	if err != nil {
		blog.Errorf("get biz count failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return false, err
	}

	if len(rst) != 1 {
		blog.Errorf("get biz count failed, for result len must be 1, filter: %+v, rid: %s", filter, kit.Rid)
		return false, kit.CCError.Error(common.CCErrOperationBizModuleHostAmountFail)
	}

	if rst[0] <= 0 {
		return false, nil
	}

	return true, nil
}

// checkHasHost check if business has hosts
func (b *business) checkHasHost(kit *rest.Kit, bizIDs []int64) (bool, error) {
	// get host count by filter
	filter := []map[string]interface{}{{
		// unarchived business
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: bizIDs,
		},
	}}

	rst, err := b.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameModuleHostConfig, filter)
	if err != nil {
		blog.Errorf("get host count failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return false, err
	}

	if len(rst) != 1 {
		blog.Errorf("get host count failed, for result len must be 1, filter: %+v, rid: %s", filter, kit.Rid)
		return false, kit.CCError.Error(common.CCErrOperationBizModuleHostAmountFail)
	}

	if rst[0] <= 0 {
		return false, nil
	}

	return true, nil
}

func (b *business) cleanBizAndRelatedResources(kit *rest.Kit, bizID int64) error {
	// 1. clean host
	if err := b.cleanHost(kit, bizID); err != nil {
		return err
	}

	// 2. clean process
	if err := b.cleanProcess(kit, bizID); err != nil {
		return err
	}

	// 3. clean service instance
	if err := b.cleanServiceInstance(kit, bizID); err != nil {
		return err
	}

	// 4. clean module
	if err := b.cleanModule(kit, bizID); err != nil {
		return err
	}

	// 5. clean set
	if err := b.cleanSet(kit, bizID); err != nil {
		return err
	}

	// 6. clean mainline topo
	if err := b.cleanTopo(kit, bizID); err != nil {
		return err
	}

	// 7. clean module/set template
	if err := b.cleanTemplate(kit, bizID); err != nil {
		return err
	}

	// 8. clean property
	if err := b.cleanProperty(kit, bizID); err != nil {
		return err
	}

	// 9. clean biz
	if err := b.cleanBiz(kit, bizID); err != nil {
		return err
	}

	return nil
}

func (b *business) cleanHost(kit *rest.Kit, bizID int64) error {
	// 1. clean host instance
	// archived business has no host, need not clean host

	// 2. clean dynamic group
	if err := b.cleanDynamicGroup(kit, bizID); err != nil {
		return err
	}

	return nil
}

func (b *business) cleanDynamicGroup(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameDynamicGroup,
		Field:     common.BKFieldID,
		Filter: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get dynamic group ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToString(rst)
	if err != nil {
		blog.Errorf("dynamic group ids to string failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	for _, id := range ids {
		rsp, err := b.clientSet.CoreService().Host().DeleteDynamicGroup(kit.Ctx, strconv.FormatInt(bizID, 10), id,
			kit.Header)
		if err != nil {
			blog.Errorf("delete dynamic group failed, biz: %v, group id: %s, err: %v, rid: %s", bizID, id, err,
				kit.Rid)
			return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("delete dynamic group failed, biz: %v, group id: %s, err: %v, rid: %s", bizID, id,
				rsp.ErrMsg, kit.Rid)
			return rsp.CCError()
		}
	}

	return nil
}

func (b *business) cleanTemplate(kit *rest.Kit, bizID int64) error {
	// 1. clean set template and set service template relation
	if err := b.cleanSetTemplate(kit, bizID); err != nil {
		return err
	}

	// 2. clean process template
	if err := b.cleanProcessTemplate(kit, bizID); err != nil {
		return err
	}

	// 3. clean service template
	if err := b.cleanServiceTemplate(kit, bizID); err != nil {
		return err
	}

	// 4. clean service category
	if err := b.cleanServiceCategory(kit, bizID); err != nil {
		return err
	}

	return nil
}

func (b *business) cleanProcessTemplate(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameProcessTemplate,
		Field:     common.BKFieldID,
		Filter: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get process template ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("process template ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	// delete process template in batch
	idsLen := len(ids)
	const batchSize = common.BKMaxPageSize
	for i := 0; i < idsLen; i += batchSize {
		idsBatch := make([]int64, 0)
		if (i + batchSize) >= idsLen {
			idsBatch = ids[i:idsLen]
		} else {
			idsBatch = ids[i : i+batchSize]
		}

		if err := b.clientSet.CoreService().Process().DeleteProcessTemplateBatch(kit.Ctx, kit.Header, idsBatch); err != nil {
			blog.Errorf("batch delete process template err: %v, rid: %s", err, kit.Rid)
			return err
		}
	}

	return nil
}

func (b *business) cleanServiceTemplate(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameServiceTemplate,
		Field:     common.BKFieldID,
		Filter: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get service template ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("service template ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	for _, id := range ids {
		if err := b.clientSet.CoreService().Process().DeleteServiceTemplate(kit.Ctx, kit.Header, id); err != nil {
			blog.Errorf("failed to delete service template, id: %v, err: %v, rid: %s", id, err, kit.Rid)
			return err
		}
	}
	return nil
}

func (b *business) cleanSetTemplate(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameSetTemplate,
		Field:     common.BKFieldID,
		Filter: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get set template ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("set template ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	// delete set template in batch
	idsLen := len(ids)
	const batchSize = common.BKMaxPageSize
	for i := 0; i < idsLen; i += batchSize {
		idsBatch := make([]int64, 0)
		if (i + batchSize) >= idsLen {
			idsBatch = ids[i:idsLen]
		} else {
			idsBatch = ids[i : i+batchSize]
		}

		opt := metadata.DeleteSetTemplateOption{
			SetTemplateIDs: idsBatch,
		}
		if err := b.clientSet.CoreService().SetTemplate().DeleteSetTemplate(kit.Ctx, kit.Header, bizID, opt); err != nil {
			blog.Errorf("batch delete set template err: %v, rid: %s", err, kit.Rid)
			return err
		}
	}

	return nil
}

func (b *business) cleanServiceCategory(kit *rest.Kit, bizID int64) error {
	opt := metadata.ListServiceCategoriesOption{
		BusinessID:     bizID,
		WithStatistics: false,
	}

	categories, err := b.clientSet.CoreService().Process().ListServiceCategories(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("get service categories failed, opt: %+v, err: %v, rid: %s", opt, err, kit.Rid)
		return err
	}

	// rearrange service categories in the order of child then parent
	parentMap := make(map[int64][]int64, 0)
	ids := make([]int64, 0)
	idMap := make(map[int64]struct{})
	for _, info := range categories.Info {
		category := info.ServiceCategory

		// ignore the shared categories
		if category.BizID != bizID {
			continue
		}

		if category.ParentID == 0 {
			idMap[category.ID] = struct{}{}
			ids = append(ids, category.ID)
			continue
		}
		parentMap[category.ParentID] = append(parentMap[category.ParentID], category.ID)
	}

	for parentID, childIDs := range parentMap {
		for _, id := range childIDs {
			idMap[id] = struct{}{}
		}

		// if parent category is a child of another category, place its children before that
		if _, exists := idMap[parentID]; exists {
			ids = append(childIDs, ids...)
			continue
		}

		ids = append(ids, childIDs...)
		idMap[parentID] = struct{}{}
		ids = append(ids, parentID)
	}

	for _, id := range ids {
		if err := b.clientSet.CoreService().Process().DeleteServiceCategory(kit.Ctx, kit.Header, id); err != nil {
			blog.Errorf("failed to delete service category, id: %v, err: %v, rid: %s", id, err, kit.Rid)
			return err
		}
	}

	return nil
}

func (b *business) cleanProcess(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameProcessInstanceRelation,
		Field:     common.BKProcessIDField,
		Filter: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	ids, errDist := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDist != nil {
		blog.Errorf("get process ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt, errDist,
			kit.Rid)
		return errDist
	}

	if len(ids) == 0 {
		return nil
	}

	// delete process in batch
	idsLen := len(ids)
	const batchSize = common.BKMaxPageSize
	for i := 0; i < idsLen; i += batchSize {
		idsBatch := make([]interface{}, 0)
		if (i + batchSize) >= idsLen {
			idsBatch = ids[i:idsLen]
		} else {
			idsBatch = ids[i : i+batchSize]
		}

		// clean process instance association
		if err := b.cleanInstAsst(kit, common.BKInnerObjIDProc, idsBatch); err != nil {
			return err
		}

		// clean process process instance
		optDelProc := metadata.DeleteOption{
			Condition: mapstr.MapStr{
				common.BKProcessIDField: mapstr.MapStr{
					common.BKDBIN: idsBatch,
				},
			},
		}

		_, err := b.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, common.BKInnerObjIDProc,
			&optDelProc)
		if err != nil {
			blog.Errorf("failed to delete process instance, ids: %v, err: %v, rid: %s", idsBatch, err, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}

		procIDs, err := util.SliceInterfaceToInt64(idsBatch)
		if err != nil {
			blog.Errorf("process ids to int failed, ids: %v, err: %v, rid: %s", idsBatch, err, kit.Rid)
			return err
		}

		// clean process instance relation
		optDelProcInstRel := metadata.DeleteProcessInstanceRelationOption{
			BusinessID: &bizID,
			ProcessIDs: procIDs,
		}

		if err := b.clientSet.CoreService().Process().DeleteProcessInstanceRelation(kit.Ctx, kit.Header,
			optDelProcInstRel); err != nil {
			return err
		}
	}

	return nil
}

func (b *business) cleanServiceInstance(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameServiceInstance,
		Field:     common.BKFieldID,
		Filter: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get service instance ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("service instance ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	// delete service instance in batch
	idsLen := len(ids)
	const batchSize = common.BKMaxPageSize
	for i := 0; i < idsLen; i += batchSize {
		idsBatch := make([]int64, 0)
		if (i + batchSize) >= idsLen {
			idsBatch = ids[i:idsLen]
		} else {
			idsBatch = ids[i : i+batchSize]
		}

		optDel := &metadata.CoreDeleteServiceInstanceOption{
			BizID:              bizID,
			ServiceInstanceIDs: idsBatch,
		}

		if err := b.clientSet.CoreService().Process().DeleteServiceInstance(kit.Ctx, kit.Header, optDel); err != nil {
			blog.Errorf("failed to delete service instance, option: %+v, err: %v, rid: %s", optDel, err, kit.Rid)
			return err
		}
	}

	return nil
}

func (b *business) cleanModule(kit *rest.Kit, bizID int64) error {
	return b.module.DeleteModule(kit, bizID, nil, nil)
}

func (b *business) cleanSet(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameBaseSet,
		Field:     common.BKSetIDField,
		Filter: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get set ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt, errDistinct, kit.Rid)
		return errDistinct
	}

	ids, err := util.SliceInterfaceToInt64(rst)
	if err != nil {
		blog.Errorf("set ids to int failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	// delete set in batch
	idsLen := len(ids)
	const batchSize = common.BKMaxPageSize
	for i := 0; i < idsLen; i += batchSize {
		idsBatch := make([]int64, 0)
		if (i + batchSize) >= idsLen {
			idsBatch = ids[i:idsLen]
		} else {
			idsBatch = ids[i : i+batchSize]
		}

		if err := b.set.DeleteSet(kit, bizID, idsBatch); err != nil {
			return err
		}
	}

	return nil
}

func (b *business) cleanTopo(kit *rest.Kit, bizID int64) error {
	// get mainline association, generate map of object and its child
	cond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.AssociationKindIDField: common.AssociationKindMainline,
		},
	}

	asstRes, err := b.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("get mainline association err: %v, rid: %s", err, kit.Rid)
		return err
	}

	childObjMap := make(map[string]string)
	for _, asst := range asstRes.Info {
		childObjMap[asst.AsstObjID] = asst.ObjectID
	}

	// delete from child of "biz"
	childObj := childObjMap[common.BKInnerObjIDApp]

	// traverse down topo till set, delete topo instances and associations
	for {
		if childObj == "" {
			blog.Errorf("failed to get mainline association, for obj id is empty, rid: %s", kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrorTopoMainlineObjectAssociationNotExist)
		}

		if childObj == common.BKInnerObjIDSet {
			break
		}

		if err := b.cleanTopoInstAndAsst(kit, bizID, childObj); err != nil {
			return err
		}

		childObj = childObjMap[childObj]
	}

	return nil
}

func (b *business) cleanTopoInstAndAsst(kit *rest.Kit, bizID int64, obj string) error {
	tableName := common.GetInstTableName(obj, kit.SupplierAccount)
	idField := common.GetInstIDField(obj)
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: tableName,
		Field:     idField,
		Filter: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	ids, errDist := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDist != nil {
		blog.Errorf("get topo inst ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt, errDist,
			kit.Rid)
		return errDist
	}

	if len(ids) == 0 {
		return nil
	}

	// delete topo instances and associations in batch
	idsLen := len(ids)
	const batchSize = common.BKMaxPageSize
	for i := 0; i < idsLen; i += batchSize {
		idsBatch := make([]interface{}, 0)
		if (i + batchSize) >= idsLen {
			idsBatch = ids[i:idsLen]
		} else {
			idsBatch = ids[i : i+batchSize]
		}

		// clean topo instance associations
		if err := b.cleanInstAsst(kit, obj, idsBatch); err != nil {
			return err
		}

		// clean topo instances
		optDel := metadata.DeleteOption{
			Condition: mapstr.MapStr{
				idField: mapstr.MapStr{
					common.BKDBIN: idsBatch,
				},
			},
		}

		_, err := b.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, obj, &optDel)
		if err != nil {
			blog.Errorf("failed to delete topo instance, ids: %v, err: %v, rid: %s", idsBatch, err, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
	}

	return nil
}

func (b *business) cleanBiz(kit *rest.Kit, bizID int64) error {
	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}

	return b.inst.DeleteInst(kit, common.BKInnerObjIDApp, cond, false)
}

func (b *business) cleanProperty(kit *rest.Kit, bizID int64) error {
	// 1. clean obj attribute description
	if err := b.cleanObjAttDes(kit, bizID); err != nil {
		return err
	}

	// 2. clean property group
	if err := b.cleanPropertyGroup(kit, bizID); err != nil {
		return err
	}

	return nil
}

func (b *business) cleanObjAttDes(kit *rest.Kit, bizID int64) error {
	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameObjAttDes,
		Field:     common.BKObjIDField,
		Filter: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	rst, errDistinct := b.clientSet.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if errDistinct != nil {
		blog.Errorf("get obj attribute ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt,
			errDistinct, kit.Rid)
		return errDistinct
	}

	objIDs, err := util.SliceInterfaceToString(rst)
	if err != nil {
		blog.Errorf("set ids to string failed, ids: %v, err: %v, rid: %s", rst, err, kit.Rid)
		return err
	}

	delOpt := &metadata.DeleteOption{
		Condition: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	for _, objID := range objIDs {
		rsp, err := b.clientSet.CoreService().Model().DeleteModelAttr(kit.Ctx, kit.Header, objID, delOpt)
		if err != nil {
			blog.Errorf("failed to request object controller, err: %v, rid: %s", err, kit.Rid)
			return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("failed to delete object attribute, obj id: %s, option: %+v,  err: %s, rid: %s", objID,
				delOpt, rsp.ErrMsg, kit.Rid)
			return kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}
	}

	return nil
}

func (b *business) cleanPropertyGroup(kit *rest.Kit, bizID int64) error {
	delOpt := metadata.DeleteOption{
		Condition: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
	}

	_, err := b.clientSet.CoreService().Model().DeleteAttributeGroupByCondition(kit.Ctx, kit.Header, delOpt)
	if err != nil {
		blog.Errorf("failed to request object controller, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (b *business) cleanInstAsst(kit *rest.Kit, objID string, instIDs []interface{}) error {
	if len(instIDs) == 0 {
		return nil
	}

	opt := &metadata.InstAsstDeleteOption{
		Opt: metadata.DeleteOption{
			Condition: mapstr.MapStr{
				common.BKDBOR: []mapstr.MapStr{
					{
						common.BKObjIDField: objID,
						common.BKInstIDField: mapstr.MapStr{
							common.BKDBIN: instIDs,
						},
					},
					{
						common.BKObjIDField: objID,
						common.BKAsstInstIDField: mapstr.MapStr{
							common.BKDBIN: instIDs,
						},
					},
				},
			},
		},
		ObjID: objID,
	}

	_, err := b.clientSet.CoreService().Association().DeleteInstAssociation(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("failed to request delete inst association, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	return nil
}
