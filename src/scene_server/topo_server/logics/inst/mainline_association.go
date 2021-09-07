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
	"regexp"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

var (
	re = regexp.MustCompile(`[#/,><|]`)
)

// ResetMainlineInstAssociation reset mainline instance association
// while a mainline object deleted may use this func
func (assoc *association) ResetMainlineInstAssociation(kit *rest.Kit, currentObjID, childObjID string) error {

	cond := mapstr.New()
	if metadata.IsCommon(currentObjID) {
		cond.Set(common.BKObjIDField, currentObjID)
	}
	// TODO 确认是否需要设置分页
	defaultCond := &metadata.QueryCondition{Condition: cond}

	// 获取 current 模型的所有实例
	currentInsts, err := assoc.FindInst(kit, currentObjID, defaultCond)
	if err != nil {
		blog.Errorf("failed to find current object(%s) inst, err: %v, rid: %s", currentObjID, err, kit.Rid)
		return err
	}

	if len(currentInsts.Info) == 0 {
		return nil
	}

	instIDs := make([]int64, len(currentInsts.Info))
	instParentMap := map[int64]int64{}
	for _, currInst := range currentInsts.Info {
		instID, err := currInst.Int64(common.GetInstIDField(currentObjID))
		if err != nil {
			blog.Errorf("get inst id in current insts failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		instIDs = append(instIDs, instID)

		currInstParentID, err := currInst.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("get bk_parent_id in current insts failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		instParentMap[instID] = currInstParentID
	}

	children, err := assoc.getMainlineChildInst(kit, currentObjID, childObjID, instIDs)
	if err != nil {
		return err
	}

	// 检查实例删除后，会不会出现重名冲突
	canReset, repeatedInstName, err := assoc.checkInstNameRepeat(kit, instParentMap, children)
	if err != nil {
		blog.Errorf("can not be reset, err: %+v, rid: %s", err, kit.Rid)
		return err
	}

	if canReset == false {
		blog.Errorf("can not be reset, inst name repeated, inst: %s, rid: %s", repeatedInstName, kit.Rid)
		errMsg := kit.CCError.CCError(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat).Error() +
			" " + repeatedInstName
		return kit.CCError.New(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat, errMsg)
	}

	// NEED FIX: 下面循环中的continue ，会在处理实例异常的时候跳过当前拓扑的处理，此方式可能会导致某个业务拓扑失败，但是不会影响所有。
	// 修改 currentInsts 所有孩子结点的父节点，为 currentInsts 的父节点，并删除 currentInsts

	childIDMap := make(map[int64][]int64)
	for _, child := range children {
		childInstID, err := child.Int64(common.GetInstIDField(currentObjID))
		if err != nil {
			blog.Errorf("get inst id in current insts failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		childParentID, err := child.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("get bk_parent_id in current insts failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		childIDMap[childParentID] = append(childIDMap[childParentID], childInstID)
	}

	parentChildMap := map[int64][]int64{}
	for inst, parent := range instParentMap {
		parentChildMap[parent] = append(parentChildMap[parent], childIDMap[inst]...)
	}

	// set the child's parent
	for parent, child := range parentChildMap {
		if len(child) == 0 {
			continue
		}

		if err = assoc.setMainlineParentInst(kit, child, childObjID, parent); err != nil {
			blog.Errorf("failed to set the object mainline child inst, parent: %d, child: %v, err: %v, rid: %s",
				parent, child, err, kit.Rid)
			return err
		}
	}

	// delete the current inst
	if err := assoc.deleteMainlineInstWithID(kit, currentObjID, instIDs); err != nil {
		blog.Errorf("failed to delete the current inst, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// SetMainlineInstAssociation set mainline instance association by parent object and current object
func (assoc *association) SetMainlineInstAssociation(kit *rest.Kit, parentObjID, childObjID, currObjID,
	currObjName string) ([]int64, error) {

	defaultCond := &metadata.QueryCondition{}
	cond := mapstr.New()
	if metadata.IsCommon(parentObjID) {
		cond.Set(common.BKObjIDField, parentObjID)
	}
	defaultCond.Condition = cond
	// fetch all parent instances.
	// TODO replace to FindInst in inst/inst.go after merge
	parentInsts, err := assoc.FindInst(kit, parentObjID, defaultCond)
	if err != nil {
		blog.Errorf("failed to find parent object(%s) inst, err: %v, rid: %s", parentObjID, err, kit.Rid)
		return nil, err
	}

	createdInstIDs := make([]int64, len(parentInsts.Info))
	expectParent2Children := make(map[int64][]int64)
	// filters out special character for mainline instances
	instanceName := re.ReplaceAllString(currObjName, "")
	var parentInstIDs []int64
	// create current object instance for each parent instance and insert the current instance to
	for _, parentInst := range parentInsts.Info {
		id, err := parentInst.Int64(common.GetInstIDField(parentObjID))
		if err != nil {
			blog.Errorf("failed to find the inst id, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		parentInstIDs = append(parentInstIDs, id)

		// we create the current object's instance for each parent instance belongs to the parent object.
		currentInst := mapstr.MapStr{common.BKObjIDField: currObjID}
		currentInst.Set(common.GetInstNameField(currObjID), instanceName)
		// set current instance's parent id to parent instance's id, so that they can be chained.
		currentInst.Set(common.BKInstParentStr, id)

		if parentObjID == common.BKInnerObjIDApp {
			currentInst.Set(common.BKAppIDField, id)
		} else {
			if bizID, ok := parentInst.Get(common.BKAppIDField); ok {
				currentInst.Set(common.BKAppIDField, bizID)
			}
		}

		// create the instance now.
		currentInstID, err := assoc.createInst(kit, currObjID, currentInst)
		if err != nil {
			blog.Errorf("failed to create object(%s) default inst, err: %v, rid: %s", currObjID, err, kit.Rid)
			return nil, err
		}

		createdInstIDs = append(createdInstIDs, int64(currentInstID))
	}

	// reset the child's parent instance's parent id to current instance's id.
	childInst, err := assoc.getMainlineChildInst(kit, parentObjID, childObjID, parentInstIDs)
	if err != nil {
		blog.Errorf("failed to get the object(%s) mainline child inst, err: %v, rid: %s",
			parentObjID, err, kit.Rid)
		return nil, err
	}

	for _, child := range childInst {
		childID, err := child.Int64(common.GetInstIDField(childObjID))
		if err != nil {
			blog.Errorf("failed to get the inst id from the inst(%#v), rid: %s", child, kit.Rid)
			continue
		}

		parentID, err := child.Int64(common.BKParentIDField)
		if err != nil {
			blog.Errorf("failed to get the object(%s) mainline parent id, "+
				"the current inst(%v), err: %v, rid: %s", childObjID, child, err, kit.Rid)
			continue
		}
		expectParent2Children[parentID] = append(expectParent2Children[parentID], childID)
	}

	for parentID, childIDs := range expectParent2Children {
		if err = assoc.setMainlineParentInst(kit, childIDs, childObjID, parentID); err != nil {
			blog.Errorf("failed to set the object mainline child inst, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
	}
	return createdInstIDs, nil
}

// SearchMainlineAssociationInstTopo search mainline association topo by objID and instID
func (assoc *association) SearchMainlineAssociationInstTopo(kit *rest.Kit, objID string, instID int64,
	withStatistics bool, withDefault bool) ([]*metadata.TopoInstRst, errors.CCError) {

	// read mainline object association and construct child relation map excluding host
	queryCond := &metadata.QueryCondition{
		Condition: map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline},
		Fields:    []string{common.BKObjIDField, common.BKAsstObjIDField},
	}
	mlAsstRsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("search mainline association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	mainlineObjectChildMap := make(map[string]string)
	isMainline := false
	for _, asst := range mlAsstRsp.Info {
		if asst.ObjectID == common.BKInnerObjIDHost {
			continue
		}
		mainlineObjectChildMap[asst.AsstObjID] = asst.ObjectID
		if asst.AsstObjID == objID {
			isMainline = true
		}
	}
	if !isMainline {
		return nil, nil
	}

	// get all mainline object name map
	objectIDs := make([]string, 0)
	for objectID := objID; len(objectID) != 0; objectID = mainlineObjectChildMap[objectID] {
		objectIDs = append(objectIDs, objectID)
	}

	objectNameMap := make(map[string]string)
	objects, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objectIDs}},
		})
	if err != nil {
		blog.Errorf("search mainline objects(%s) failed, err: %V, rid: %s", objectIDs, err, kit.Rid)
		return nil, err
	}

	for _, object := range objects.Info {
		objectNameMap[object.GetObjectID()] = object.ObjectName
	}

	// traverse and fill instance topology data
	results := make([]*metadata.TopoInstRst, 0)
	var parents []*metadata.TopoInstRst
	instCond := map[string]interface{}{
		common.GetInstIDField(objID): instID,
	}
	var bizID int64
	moduleIDs := make([]int64, 0)
	for objectID := objID; len(objectID) != 0; objectID = mainlineObjectChildMap[objectID] {
		filter := &metadata.QueryCondition{Condition: instCond}
		if objectID != objID {
			filter.Page.Sort = common.BKDefaultField
		}

		filter.Fields = []string{common.GetInstIDField(objectID), common.GetInstNameField(objectID),
			common.BKDefaultField, common.BKParentIDField, common.BKAppIDField}
		if objectID == common.BKInnerObjIDSet {
			filter.Fields = append(filter.Fields, common.BKSetTemplateIDField)
		} else if objectID == common.BKInnerObjIDModule {
			filter.Fields = append(filter.Fields, []string{common.BKServiceTemplateIDField, common.BKSetTemplateIDField,
				common.HostApplyEnabledField}...)
		}

		instanceRsp, err := assoc.FindInst(kit, objectID, filter)
		if err != nil {
			blog.Errorf("search inst failed, err: %s, cond:%s, rid: %s", err, instCond, kit.Rid)
			return nil, err
		}
		// already reached the deepest level, stop the loop
		if len(instanceRsp.Info) == 0 {
			break
		}
		instIDs := make([]int64, 0)
		objectName := objectNameMap[objectID]
		instances := make([]*metadata.TopoInstRst, 0)
		// map parentID to its children, not including default set
		childInstMap := make(map[int64][]*metadata.TopoInstRst)
		// map parentID to its default set children, default sets are children of biz
		childDefaultSetMap := make(map[int64][]*metadata.TopoInstRst)
		for _, instance := range instanceRsp.Info {
			instID, err := instance.Int64(common.GetInstIDField(objectID))
			if err != nil {
				blog.Errorf("get instance %#v id failed, err: %v, rid: %s", instance, err, kit.Rid)
				return nil, err
			}
			instIDs = append(instIDs, instID)
			instName, err := instance.String(common.GetInstNameField(objectID))
			if err != nil {
				blog.Errorf("get instance %#v name failed, err: %v, rid: %s", instance, err, kit.Rid)
				return nil, err
			}
			defaultValue := 0
			if defaultFieldValue, exist := instance[common.BKDefaultField]; exist {
				defaultValue, err = util.GetIntByInterface(defaultFieldValue)
				if err != nil {
					blog.Errorf("get instance %#v default failed, err: %v, rid: %s", instance, err, kit.Rid)
					return nil, err
				}
			}
			topoInst := &metadata.TopoInstRst{
				TopoInst: metadata.TopoInst{
					InstID:   instID,
					InstName: instName,
					ObjID:    objectID,
					ObjName:  objectName,
					Default:  defaultValue,
				},
				Child: []*metadata.TopoInstRst{},
			}
			if withStatistics {
				if objectID == common.BKInnerObjIDSet {
					topoInst.SetTemplateID, _ = instance.Int64(common.BKSetTemplateIDField)
				}
				if objectID == common.BKInnerObjIDModule {
					topoInst.ServiceTemplateID, _ = instance.Int64(common.BKServiceTemplateIDField)
					topoInst.SetTemplateID, _ = instance.Int64(common.BKSetTemplateIDField)
					enabled, _ := instance.Bool(common.HostApplyEnabledField)
					topoInst.HostApplyEnabled = &enabled
					moduleIDs = append(moduleIDs, instID)
				}
				if bizID == 0 {
					bizID, err = instance.Int64(common.BKAppIDField)
					if err != nil {
						blog.Errorf("get instance %#v biz id failed, err: %v, rid: %s", instance, err, kit.Rid)
						return nil, err
					}
				}
			}
			if objectID == objID {
				results = append(results, topoInst)
			} else {
				parentID, err := instance.Int64(common.BKParentIDField)
				if err != nil {
					blog.Errorf("get instance %#v parent id failed, err: %v, rid: %s", instance, err, kit.Rid)
					return nil, err
				}
				if objectID == common.BKInnerObjIDSet && defaultValue == common.DefaultResSetFlag {
					childDefaultSetMap[parentID] = append(childDefaultSetMap[parentID], topoInst)
				} else {
					childInstMap[parentID] = append(childInstMap[parentID], topoInst)
				}
			}
			instances = append(instances, topoInst)
		}
		// set children for parents, default sets are children of biz
		for _, parentInst := range parents {
			parentInst.Child = append(parentInst.Child, childInstMap[parentInst.InstID]...)
		}
		if objectID == common.BKInnerObjIDSet && objID == common.BKInnerObjIDApp {
			for _, parentInst := range results {
				parentInst.Child = append(parentInst.Child, childDefaultSetMap[parentInst.InstID]...)
			}
		}
		// set current instances as parents and generate condition for next level
		instCond = make(map[string]interface{})
		if mainlineObjectChildMap[objectID] == common.BKInnerObjIDSet {
			if withDefault {
				// default sets are children of biz, so need to add biz into parent condition
				instIDs = append(instIDs, bizID)
			} else {
				instCond[common.BKDefaultField] = map[string]interface{}{
					common.BKDBNE: common.DefaultResSetFlag,
				}
			}
		}
		parents = instances
		instCond[common.BKInstParentStr] = map[string]interface{}{
			common.BKDBIN: instIDs,
		}
	}

	if withStatistics && len(results) > 0 {
		if err := assoc.fillStatistics(kit, bizID, moduleIDs, results); err != nil {
			blog.Errorf("fill statistics data failed, bizID: %d, err: %v, rid: %s", bizID, err, kit.Rid)
			return nil, err
		}
	}
	return results, nil
}

func (assoc *association) getMainlineChildInst(kit *rest.Kit, objID, childObjID string, instIDs []int64) (
	[]mapstr.MapStr, error) {

	cond := mapstr.MapStr{common.BKInstParentStr: mapstr.MapStr{common.BKDBIN: instIDs}}
	if metadata.IsCommon(childObjID) {
		cond.Set(metadata.ModelFieldObjectID, childObjID)
	} else if childObjID == common.BKInnerObjIDSet {
		cond.Set(common.BKDefaultField, mapstr.MapStr{common.BKDBNE: common.DefaultResSetFlag})
	}

	instCond := &metadata.QueryCondition{
		Condition: cond,
		Fields:    []string{common.GetInstIDField(childObjID), common.BKParentIDField},
	}
	instRsp, err := assoc.FindInst(kit, childObjID, instCond)
	if err != nil {
		blog.Errorf("search inst by object(%s) failed, err: %v, rid: %s", childObjID, err, kit.Rid)
		return nil, err
	}

	return instRsp.Info, nil
}

// FindInst search instance by condition
// TODO need to delete after merge
func (assoc *association) FindInst(kit *rest.Kit, objID string, cond *metadata.QueryCondition) (*metadata.InstResult,
	error) {

	switch objID {
	case common.BKInnerObjIDHost:
		input := &metadata.QueryInput{
			Condition:     cond.Condition,
			Fields:        strings.Join(cond.Fields, ","),
			TimeCondition: cond.TimeCondition,
			Start:         cond.Page.Start,
			Limit:         cond.Page.Limit,
			Sort:          cond.Page.Sort,
		}
		rsp, err := assoc.clientSet.CoreService().Host().GetHosts(kit.Ctx, kit.Header, input)
		if err != nil {
			blog.Errorf("search object(%s) inst by the input(%#v) failed, err: %v, rid: %s", objID, input, err, kit.Rid)
			return nil, err
		}

		return &metadata.InstResult{Count: rsp.Count, Info: rsp.Info}, nil

	default:
		rsp, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, cond)
		if err != nil {
			blog.Errorf("search object(%s) inst by the cond(%#v) failed, err: %v, rid: %s", objID, cond, err, kit.Rid)
			return nil, err
		}

		return &metadata.InstResult{Count: rsp.Count, Info: rsp.Info}, nil
	}
}

func (assoc *association) setMainlineParentInst(kit *rest.Kit, childID []int64, childObjID string,
	parentID int64) error {

	if len(childID) == 0 {
		return nil
	}

	cond := mapstr.MapStr{common.GetInstIDField(childObjID): mapstr.MapStr{common.BKDBIN: childID}}
	if metadata.IsCommon(childObjID) {
		cond.Set(metadata.ModelFieldObjectID, childObjID)
	}

	input := metadata.UpdateOption{
		Data: mapstr.MapStr{
			common.BKInstParentStr: parentID,
		},
		Condition: cond,
	}
	_, err := assoc.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, childObjID, &input)
	if err != nil {
		blog.Errorf("failed to update the association, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (assoc *association) deleteMainlineInstWithID(kit *rest.Kit, objID string, instID []int64) error {

	// if this instance has been bind to a instance by the association, then this instance should not be deleted.
	cnt, err := assoc.clientSet.CoreService().Association().CountInstanceAssociations(kit.Ctx, kit.Header, objID,
		&metadata.Condition{
			Condition: mapstr.MapStr{common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: objID, common.BKInstIDField: mapstr.MapStr{common.BKDBIN: instID}},
				{common.BKAsstObjIDField: objID, common.BKAsstInstIDField: mapstr.MapStr{common.BKDBIN: instID}},
			}},
		})
	if err != nil {
		blog.Errorf("count association by object(%s) failed, err: %s, rid: %s", objID, err, kit.Rid)
		return err
	}

	if cnt.Count > 0 {
		return kit.CCError.CCError(common.CCErrorInstHasAsst)
	}

	// delete this instance now.
	delCond := mapstr.MapStr{common.GetInstIDField(objID): instID}
	if metadata.IsCommon(objID) {
		delCond.Set(common.BKObjIDField, objID)
	}

	// generate audit log.
	audit := auditlog.NewInstanceAudit(assoc.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, delCond)
	if err != nil {
		blog.Errorf(" delete inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// to delete.
	ops := metadata.DeleteOption{
		Condition: delCond,
	}
	_, err = assoc.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, objID, &ops)
	if err != nil {
		blog.Errorf("failed to delete the object(%s) inst by the condition(%#v), err: %v", objID, ops, err)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
		blog.Errorf("delete inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

// TODO should be deleted after merge
func (assoc *association) createInst(kit *rest.Kit, objID string, data mapstr.MapStr) (uint64, error) {
	cond := &metadata.CreateModelInstance{
		Data: data,
	}
	rsp, err := assoc.clientSet.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, objID, cond)
	if err != nil {
		blog.Errorf("failed to create object instance, err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}

	return rsp.Created.ID, nil
}

// TODO need check this function is here or move to other file
func (assoc *association) fillStatistics(kit *rest.Kit, bizID int64, moduleIDs []int64,
	topoInsts []*metadata.TopoInstRst) errors.CCError {
	// get host apply rule count
	listApplyRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: moduleIDs,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}
	hostApplyRules, err := assoc.clientSet.CoreService().HostApplyRule().ListHostApplyRule(kit.Ctx, kit.Header,
		bizID, listApplyRuleOption)
	if err != nil {
		blog.Errorf("fillStatistics failed, list host apply rule failed, bizID: %s, option: %s, err: %s, rid: %s",
			bizID, listApplyRuleOption, err, kit.Rid)
		return err
	}
	moduleRuleCount := make(map[int64]int64)
	for _, item := range hostApplyRules.Info {
		moduleRuleCount[item.ModuleID]++
	}

	// fill hosts
	for _, tir := range topoInsts {
		tir.DeepFirstTraverse(func(node *metadata.TopoInstRst) {
			if node.ObjID == common.BKInnerObjIDModule {
				node.HostApplyRuleCount = new(int64)
				*node.HostApplyRuleCount, _ = moduleRuleCount[node.InstID]
			}

			if len(node.Child) == 0 {
				return
			}
		})
	}
	return nil
}

// checkInstNameRepeat 检查如果将 currentInsts 都删除之后，拥有共同父节点的孩子结点会不会出现名字冲突
// 如果有冲突，返回 (false, 冲突实例名, nil)
func (assoc *association) checkInstNameRepeat(kit *rest.Kit, instParentMap map[int64]int64, children []mapstr.MapStr) (
	canReset bool, repeatedInstName string, err error) {

	parentChildName := map[int64]map[string]struct{}{}
	for _, child := range children {
		instName, err := child.String(common.BKInstNameField)
		if err != nil {
			return false, "", err
		}

		childInstID, err := child.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("get child parent id in child insts failed, err: %v, rid: %s", err, kit.Rid)
			return false, "", err
		}

		for instID, parentID := range instParentMap {
			if instID != childInstID {
				continue
			}

			childName, exist := parentChildName[parentID]
			if !exist {
				childName = make(map[string]struct{})
				parentChildName[parentID] = childName
			}

			if _, exist := childName[instName]; exist {
				return false, instName, nil
			}

			childName[instName] = struct{}{}
		}
	}

	return true, "", nil
}
