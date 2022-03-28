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
	"sync"

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
	mainlineSpecialCharacterRegexp = regexp.MustCompile(`[#/,><|]`)
)

type buildTopoInstRst struct {
	result    []*metadata.TopoInstRst
	moduleIDs []int64
	bizID     int64
}

// NewBuildTopoInstRst return a buildTopoInstRst struct instance
func NewBuildTopoInstRst() *buildTopoInstRst {
	return &buildTopoInstRst{
		result:    make([]*metadata.TopoInstRst, 0),
		moduleIDs: make([]int64, 0),
		bizID:     0,
	}
}

// ResetMainlineInstAssociation reset mainline instance association
// while a mainline object deleted may use this func
func (assoc *association) ResetMainlineInstAssociation(kit *rest.Kit, currentObjID, childObjID string) error {

	defaultCond := &metadata.QueryCondition{Condition: mapstr.New()}
	if metadata.IsCommon(currentObjID) {
		defaultCond.Condition.Set(common.BKObjIDField, currentObjID)
	}

	// 获取 current 模型的所有实例
	currentInsts, err := assoc.inst.FindInst(kit, currentObjID, defaultCond)
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
			blog.Errorf("get inst id(%d) in current insts failed, err: %v, rid: %s", instID, err, kit.Rid)
			return err
		}

		instIDs = append(instIDs, instID)

		instParentID, err := currInst.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("get bk_parent_id(%d) in current insts failed, err: %v, rid: %s", instParentID, err, kit.Rid)
			return err
		}

		instParentMap[instID] = instParentID
	}

	children, err := assoc.getMainlineChildInst(kit, currentObjID, childObjID, instIDs)
	if err != nil {
		return err
	}

	// 检查实例删除后，会不会出现重名冲突
	canReset, repeatedInstName, err := assoc.checkInstNameRepeat(kit, instParentMap, childObjID, children)
	if err != nil {
		blog.Errorf("can not be reset, err: %+v, rid: %s", err, kit.Rid)
		return err
	}

	if !canReset {
		blog.Errorf("can not be reset, inst name repeated, inst: %s, rid: %s", repeatedInstName, kit.Rid)
		return kit.CCError.CCError(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat)
	}

	// 修改 currentInsts 所有孩子结点的父节点，为 currentInsts 的父节点，并删除 currentInsts
	parentChildMap := map[int64][]int64{}
	for _, child := range children {
		childInstID, err := child.Int64(common.GetInstIDField(childObjID))
		if err != nil {
			blog.Errorf("get inst id in current insts failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		childParentID, err := child.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("get bk_parent_id in current insts failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		parentChildMap[instParentMap[childParentID]] = append(parentChildMap[instParentMap[childParentID]], childInstID)
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

	defaultCond := &metadata.QueryCondition{Condition: mapstr.New()}
	if metadata.IsCommon(parentObjID) {
		defaultCond.Condition.Set(common.BKObjIDField, parentObjID)
	}

	// fetch all parent instances.
	parentInsts, err := assoc.inst.FindInst(kit, parentObjID, defaultCond)
	if err != nil {
		blog.Errorf("failed to find parent object(%s) inst, err: %v, rid: %s", parentObjID, err, kit.Rid)
		return nil, err
	}

	createdInstIDs := make([]int64, len(parentInsts.Info))
	newParentCurrentMap := make(map[int64]int64)
	// filters out special character for mainline instances
	instanceName := mainlineSpecialCharacterRegexp.ReplaceAllString(currObjName, "")
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
		instRsp, err := assoc.inst.CreateInst(kit, currObjID, currentInst)
		if err != nil {
			blog.Errorf("failed to create object(%s) default inst, err: %v, rid: %s", currObjID, err, kit.Rid)
			return nil, err
		}
		currentInstID, err := instRsp.Int64(common.GetInstIDField(currObjID))
		if err != nil {
			blog.Errorf("get current object(%s) inst id failed, err: %v, rid: %s", currObjID, err, kit.Rid)
			return nil, err
		}

		newParentCurrentMap[id] = currentInstID
		createdInstIDs = append(createdInstIDs, currentInstID)
	}

	// reset the child's parent instance's parent id to current instance's id.
	childInst, err := assoc.getMainlineChildInst(kit, parentObjID, childObjID, parentInstIDs)
	if err != nil {
		blog.Errorf("failed to get the object(%s) mainline child inst, err: %v, rid: %s", parentObjID, err, kit.Rid)
		return nil, err
	}

	expectParent2Children := make(map[int64][]int64)
	for _, child := range childInst {
		childID, err := child.Int64(common.GetInstIDField(childObjID))
		if err != nil {
			blog.Errorf("failed to get the inst id from the inst(%#v), rid: %s", child, kit.Rid)
			continue
		}

		parentID, err := child.Int64(common.BKParentIDField)
		if err != nil {
			blog.Errorf("get parent id failed, current inst(%v), err: %v, rid: %s", child, err, kit.Rid)
			continue
		}
		newParentID := newParentCurrentMap[parentID]
		expectParent2Children[newParentID] = append(expectParent2Children[newParentID], childID)
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

	queryCond = &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objectIDs}},
		Fields:    []string{common.BKObjIDField, common.BKObjNameField},
	}
	objects, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("search mainline objects(%s) failed, err: %V, rid: %s", objectIDs, err, kit.Rid)
		return nil, err
	}

	objectNameMap := make(map[string]string)
	for _, object := range objects.Info {
		objectNameMap[object.GetObjectID()] = object.ObjectName
	}

	results, err := assoc.buildTopoInstRst(kit, instID, objID, objectNameMap, mainlineObjectChildMap, withDefault,
		withStatistics)
	if err != nil {
		blog.Errorf("build topo inst result failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	if withStatistics && len(results.result) > 0 {
		if err := assoc.fillStatistics(kit, results.bizID, results.moduleIDs, results.result); err != nil {
			blog.Errorf("fill statistics data failed, bizID: %d, err: %v, rid: %s", results.bizID, err, kit.Rid)
			return nil, err
		}
	}

	return results.result, nil
}

func (assoc *association) buildTopoInstRst(kit *rest.Kit, instID int64, objID string, objNameMap,
	mlObjChildMap map[string]string, withDefault, withStatistics bool) (*buildTopoInstRst, error) {

	results := NewBuildTopoInstRst()
	parents := make([]*metadata.TopoInstRst, 0)
	instCond := map[string]interface{}{common.GetInstIDField(objID): instID}

	for objectID := objID; len(objectID) != 0; objectID = mlObjChildMap[objectID] {

		instanceRsp, err := assoc.searchMainlineObjInstForTopo(kit, objectID, instCond, objectID != objID)
		if err != nil {
			blog.Errorf("search topo mainline inst failed, objID: %s, err: %v, rid: %s", objectID, err, kit.Rid)
			return nil, err
		}

		if len(instanceRsp) == 0 {
			return results, nil
		}

		instIDs := make([]int64, 0)
		instances := make([]*metadata.TopoInstRst, 0)
		childInstMap := make(map[int64][]*metadata.TopoInstRst)
		childDefaultSetMap := make(map[int64][]*metadata.TopoInstRst)
		for _, instance := range instanceRsp {
			topoInst, err := assoc.makeTopoInstRst(kit, objectID, objNameMap[objectID], instance)
			if err != nil {
				return nil, err
			}

			instIDs = append(instIDs, topoInst.InstID)

			if withStatistics {
				if objectID == common.BKInnerObjIDSet {
					topoInst.SetTemplateID, _ = instance.Int64(common.BKSetTemplateIDField)
				}
				if objectID == common.BKInnerObjIDModule {
					topoInst.ServiceTemplateID, _ = instance.Int64(common.BKServiceTemplateIDField)
					topoInst.SetTemplateID, _ = instance.Int64(common.BKSetTemplateIDField)
					enabled, _ := instance.Bool(common.HostApplyEnabledField)
					topoInst.HostApplyEnabled = &enabled
					results.moduleIDs = append(results.moduleIDs, topoInst.InstID)
				}
				if results.bizID == 0 {
					results.bizID, err = instance.Int64(common.BKAppIDField)
					if err != nil {
						blog.Errorf("get instance %#v biz id failed, err: %v, rid: %s", instance, err, kit.Rid)
						return nil, err
					}
				}
			}
			if objectID == objID {
				results.result = append(results.result, topoInst)
			} else {
				parentID, err := instance.Int64(common.BKParentIDField)
				if err != nil {
					blog.Errorf("get instance %#v parent id failed, err: %v, rid: %s", instance, err, kit.Rid)
					return nil, err
				}
				if objectID == common.BKInnerObjIDSet && topoInst.Default == common.DefaultResSetFlag {
					childDefaultSetMap[parentID] = append(childDefaultSetMap[parentID], topoInst)
				} else {
					childInstMap[parentID] = append(childInstMap[parentID], topoInst)
				}
			}
			instances = append(instances, topoInst)
		}
		for _, parentInst := range parents {
			parentInst.Child = append(parentInst.Child, childInstMap[parentInst.InstID]...)
		}

		if objectID == common.BKInnerObjIDSet && objID == common.BKInnerObjIDApp {
			for _, parentInst := range results.result {
				parentInst.Child = append(parentInst.Child, childDefaultSetMap[parentInst.InstID]...)
			}
		}
		instCond = make(map[string]interface{})
		if mlObjChildMap[objectID] == common.BKInnerObjIDSet {
			if withDefault {
				instIDs = append(instIDs, results.bizID)
			} else {
				instCond[common.BKDefaultField] = map[string]interface{}{common.BKDBNE: common.DefaultResSetFlag}
			}
		}
		parents = instances
		instCond[common.BKInstParentStr] = map[string]interface{}{common.BKDBIN: instIDs}
	}

	return results, nil
}

func (assoc *association) searchMainlineObjInstForTopo(kit *rest.Kit, objID string, instCond map[string]interface{},
	isSortDefault bool) ([]mapstr.MapStr, error) {

	filter := &metadata.QueryCondition{Condition: instCond}
	if isSortDefault {
		filter.Page.Sort = common.BKDefaultField
	}

	filter.Fields = []string{common.GetInstIDField(objID), common.GetInstNameField(objID),
		common.BKDefaultField, common.BKParentIDField, common.BKAppIDField}
	if objID == common.BKInnerObjIDSet {
		filter.Fields = append(filter.Fields, common.BKSetTemplateIDField)
	} else if objID == common.BKInnerObjIDModule {
		filter.Fields = append(filter.Fields, []string{common.BKServiceTemplateIDField, common.BKSetTemplateIDField,
			common.HostApplyEnabledField}...)
	}

	instanceRsp, err := assoc.inst.FindInst(kit, objID, filter)
	if err != nil {
		blog.Errorf("search inst failed, err: %s, cond:%s, rid: %s", err, instCond, kit.Rid)
		return nil, err
	}

	return instanceRsp.Info, nil
}

func (assoc *association) makeTopoInstRst(kit *rest.Kit, objID, objName string, inst mapstr.MapStr) (
	*metadata.TopoInstRst, error) {

	instID, err := inst.Int64(common.GetInstIDField(objID))
	if err != nil {
		blog.Errorf("get instance %#v id failed, err: %v, rid: %s", inst, err, kit.Rid)
		return nil, err
	}

	instName, err := inst.String(common.GetInstNameField(objID))
	if err != nil {
		blog.Errorf("get instance %#v name failed, err: %v, rid: %s", inst, err, kit.Rid)
		return nil, err
	}
	defaultValue := 0
	if defaultFieldValue, exist := inst[common.BKDefaultField]; exist {
		defaultValue, err = util.GetIntByInterface(defaultFieldValue)
		if err != nil {
			blog.Errorf("get instance %#v default failed, err: %v, rid: %s", inst, err, kit.Rid)
			return nil, err
		}
	}
	topoInst := &metadata.TopoInstRst{
		TopoInst: metadata.TopoInst{
			InstID:   instID,
			InstName: instName,
			ObjID:    objID,
			ObjName:  objName,
			Default:  defaultValue,
		},
		Child: []*metadata.TopoInstRst{},
	}

	return topoInst, nil
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
		Fields: []string{
			common.GetInstIDField(childObjID),
			common.GetInstNameField(childObjID),
			common.BKParentIDField,
		},
	}
	instRsp, err := assoc.inst.FindInst(kit, childObjID, instCond)
	if err != nil {
		blog.Errorf("search inst by object(%s) failed, err: %v, rid: %s", childObjID, err, kit.Rid)
		return nil, err
	}

	return instRsp.Info, nil
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

	if err := assoc.CheckAssociations(kit, objID, instID); err != nil {
		blog.Errorf("check object(%s) insts(%v) associations failed, err: %v, rid: %s", objID, instID, err, kit.Rid)
		return err
	}

	// delete this instance now.
	delCond := mapstr.MapStr{common.GetInstIDField(objID): mapstr.MapStr{common.BKDBIN: instID}}
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
	ops := metadata.DeleteOption{Condition: delCond}
	if _, err = assoc.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, objID, &ops); err != nil {
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
				*node.HostApplyRuleCount = moduleRuleCount[node.InstID]
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
func (assoc *association) checkInstNameRepeat(kit *rest.Kit, instParentMap map[int64]int64, childObjID string,
	children []mapstr.MapStr) (canReset bool, repeatedInstName string, err error) {

	parentChildName := map[int64]map[string]struct{}{}
	for _, child := range children {
		instName, err := child.String(common.GetInstNameField(childObjID))
		if err != nil {
			blog.Errorf("get child name in child insts failed, err: %v, rid: %s", err, kit.Rid)
			return false, "", err
		}

		childParentID, err := child.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("get child parent id in child insts failed, err: %v, rid: %s", err, kit.Rid)
			return false, "", err
		}

		childNameMap, exist := parentChildName[instParentMap[childParentID]]
		if !exist {
			parentChildName[instParentMap[childParentID]] = map[string]struct{}{instName: {}}
			continue
		}

		if _, exist := childNameMap[instName]; exist {
			return false, instName, nil
		}

		parentChildName[instParentMap[childParentID]][instName] = struct{}{}

	}

	return true, "", nil
}

// TopoNodeHostAndSerInstCount get topo node host and service instance count
func (assoc *association) TopoNodeHostAndSerInstCount(kit *rest.Kit, input *metadata.HostAndSerInstCountOption) (
	[]*metadata.TopoNodeHostAndSerInstCount, errors.CCError) {

	bizIDs := make([]int64, 0)
	setIDs := make([]int64, 0)
	moduleIDs := make([]int64, 0)
	customLevels := make(map[int64]string)
	for _, obj := range input.Condition {
		switch obj.ObjID {
		case common.BKInnerObjIDSet:
			setIDs = append(setIDs, obj.InstID)
		case common.BKInnerObjIDModule:
			moduleIDs = append(moduleIDs, obj.InstID)
		case common.BKInnerObjIDApp:
			bizIDs = append(bizIDs, obj.InstID)
		default:
			customLevels[obj.InstID] = obj.ObjID
		}
	}

	results := make([]*metadata.TopoNodeHostAndSerInstCount, 0)
	// handle module host number and service instance count
	if len(moduleIDs) > 0 {
		res, err := assoc.getHostSvcInstCountByInstIDsDirectly(kit, common.BKInnerObjIDModule, moduleIDs)
		if err != nil {
			blog.Errorf("get module host and service instance count failed, err: %v, moduleIDs: %v, rid: %s", err,
				moduleIDs, kit.Rid)
			return nil, err
		}
		results = append(results, res...)
	}

	// handle set host number and service instance count
	if len(setIDs) > 0 {
		res, err := assoc.getHostSvcInstCountBySetIDs(kit, setIDs)
		if err != nil {
			blog.Errorf("get set host and service instance count failed, err: %v, setIDs: %v, rid: %s", err,
				setIDs, kit.Rid)
			return nil, err
		}
		results = append(results, res...)
	}

	// handle biz host and service instance count
	if len(bizIDs) > 0 {
		res, err := assoc.getHostSvcInstCountByInstIDsDirectly(kit, common.BKInnerObjIDApp, bizIDs)
		if err != nil {
			blog.Errorf("get biz host and service instance count failed, err: %v, bizID: %v, rid: %s", err, bizIDs,
				kit.Rid)
			return nil, err
		}
		results = append(results, res...)
	}

	// handle custom level host count
	if len(customLevels) > 0 {
		res, err := assoc.getCustomLevHostSvcInstCount(kit, customLevels)
		if err != nil {
			blog.Errorf("get custom level host and service instance count failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		results = append(results, res...)
	}

	return results, nil
}

// getHostSvcInstCountByInstIDsDirectly get host and service instance count by biz or module ids directly
func (assoc *association) getHostSvcInstCountByInstIDsDirectly(kit *rest.Kit, objID string, instIDs []int64) (
	[]*metadata.TopoNodeHostAndSerInstCount, error) {

	instArr := make([][]int64, 0)
	for _, instID := range instIDs {
		instArr = append(instArr, []int64{instID})
	}

	instIDField := common.GetInstIDField(objID)
	svcInstCounts, e := assoc.getServiceInstCount(kit, instIDField, instArr)
	if e != nil {
		blog.Errorf("count service instance failed, err: %v, obj: %s, instIDs: %v, rid: %s", e, objID, instIDs, kit.Rid)
		return nil, e
	}

	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 5)
	results := make([]*metadata.TopoNodeHostAndSerInstCount, len(instIDs))
	for idx, instID := range instIDs {
		pipeline <- true
		wg.Add(1)
		go func(idx int, instID int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()
			instArr := []int64{instID}
			hostCount, err := assoc.getDistinctHostCount(kit, instIDField, instArr)
			if err != nil {
				blog.Errorf("get distinct host count failed, err: %v, objID: %s, instID: %s, rid: %s", err,
					instIDField, instID, kit.Rid)
				firstErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)

				return
			}

			topoNodeHostCount := &metadata.TopoNodeHostAndSerInstCount{
				ObjID:                objID,
				InstID:               instID,
				HostCount:            hostCount,
				ServiceInstanceCount: svcInstCounts[idx],
			}

			lock.Lock()
			results[idx] = topoNodeHostCount
			lock.Unlock()
		}(idx, instID)
	}
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return results, nil
}

// getHostSvcInstCountBySetIDs get host and service instace count by set ids
func (assoc *association) getHostSvcInstCountBySetIDs(kit *rest.Kit,
	setIDs []int64) ([]*metadata.TopoNodeHostAndSerInstCount, error) {
	setRelModuleMap, e := assoc.getSetRelationModule(kit, setIDs)
	if e != nil {
		blog.Errorf("get set module rel map failed, err: %s, rid: %s", e.Error(), kit.Rid)
		return nil, e
	}

	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr errors.CCErrorCoder
	moduleIDs := make([][]int64, 0)
	pipeline := make(chan bool, 5)
	results := make([]*metadata.TopoNodeHostAndSerInstCount, len(setIDs))
	for idx, setID := range setIDs {
		moduleIDs = append(moduleIDs, setRelModuleMap[setID])
		pipeline <- true
		wg.Add(1)
		go func(idx int, setID int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()
			setArr := []int64{setID}
			hostCount, err := assoc.getDistinctHostCount(kit, common.BKSetIDField, setArr)
			if err != nil {
				blog.Errorf("get distinct host count failed, err: %v, objID: %s, instID: %s, rid: %s", err,
					common.BKSetIDField, setID, kit.Rid)
				firstErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)

				return
			}

			topoNodeHostCount := &metadata.TopoNodeHostAndSerInstCount{
				ObjID:     common.BKInnerObjIDSet,
				InstID:    setID,
				HostCount: hostCount,
			}
			lock.Lock()
			results[idx] = topoNodeHostCount
			lock.Unlock()
		}(idx, setID)
	}
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	svcInstCounts, e := assoc.getServiceInstCount(kit, common.BKModuleIDField, moduleIDs)
	if e != nil {
		blog.Errorf("get service instance count failed, err: %v, objID: %s, instID: %s, rid: %s", e,
			common.BKSetIDField, moduleIDs, kit.Rid)
		return nil, e
	}

	for idx, count := range svcInstCounts {
		results[idx].ServiceInstanceCount = count
	}

	return results, nil
}

// getCustomLevHostSvcInstCount get coustom level host and service instace
func (assoc *association) getCustomLevHostSvcInstCount(kit *rest.Kit,
	customLevels map[int64]string) ([]*metadata.TopoNodeHostAndSerInstCount, error) {

	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 5)
	results := make([]*metadata.TopoNodeHostAndSerInstCount, 0)

	for instID, objID := range customLevels {
		pipeline <- true
		wg.Add(1)
		go func(instID int64, objID string) {
			defer func() {
				wg.Done()
				<-pipeline
			}()
			setIDArr, err := assoc.getSetIDsByTopo(kit, objID, []int64{instID})
			if err != nil {
				blog.Errorf("get set ID by topo err: %v, objID: %s, instID: %d, rid:%s", err, objID, instID, kit.Rid)
				firstErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)
				return
			}

			if len(setIDArr) == 0 {
				topoNodeCount := &metadata.TopoNodeHostAndSerInstCount{
					ObjID:                objID,
					InstID:               instID,
					HostCount:            0,
					ServiceInstanceCount: 0,
				}

				lock.Lock()
				results = append(results, topoNodeCount)
				lock.Unlock()
				return
			}

			// get host count by set ids
			hostCount, err := assoc.getDistinctHostCount(kit, common.BKSetIDField, setIDArr)
			if err != nil {
				blog.Errorf("get distinct host count failed, err: %v, objID: %s, instIDs: %, rid: %s", err,
					common.BKSetIDField, setIDArr, kit.Rid)
				firstErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)
				return
			}

			// get service instance count by set ids
			setRelModuleMap, e := assoc.getSetRelationModule(kit, setIDArr)
			if e != nil {
				blog.Errorf("get set module rel map failed, err: %s, rid: %s", e.Error(), kit.Rid)
				firstErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)

				return
			}
			moduleIDs := make([]int64, 0)
			for _, moduleSlice := range setRelModuleMap {
				moduleIDs = append(moduleIDs, moduleSlice...)
			}

			svcInstCount := int64(0)
			if len(moduleIDs) > 0 {
				moduleIDs = util.IntArrayUnique(moduleIDs)

				cond := make([][]int64, 0)
				cond = append(cond, moduleIDs)
				svcInstCounts, e := assoc.getServiceInstCount(kit, common.BKModuleIDField, cond)
				if e != nil {
					blog.Errorf("get service instance count failed, err: %v, objID: %s, instIDs: %s, rid: %s", e,
						common.BKSetIDField, moduleIDs, kit.Rid)
					firstErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)
					return
				}

				svcInstCount = svcInstCounts[0]
			}

			topoNodeCount := &metadata.TopoNodeHostAndSerInstCount{
				ObjID:                objID,
				InstID:               instID,
				HostCount:            hostCount,
				ServiceInstanceCount: svcInstCount,
			}

			lock.Lock()
			results = append(results, topoNodeCount)
			lock.Unlock()
		}(instID, objID)
	}
	wg.Wait()

	return results, firstErr
}

// getServiceInstCount get toponode service instance count
func (assoc *association) getServiceInstCount(kit *rest.Kit, objID string, instIDs [][]int64) ([]int64, error) {
	filters := make([]map[string]interface{}, 0)
	results := make([]int64, len(instIDs))
	zeroIndexMap := make(map[int]struct{})

	for idx, instID := range instIDs {
		if len(instID) == 0 {
			zeroIndexMap[idx] = struct{}{}
			continue
		}
		filters = append(filters, mapstr.MapStr{objID: mapstr.MapStr{common.BKDBIN: instID}})
	}

	svcInstCount, err := assoc.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameServiceInstance, filters)
	if err != nil {
		blog.Errorf("find service instance count failed, err: %v, objID: %s, instIDs: %, rid: %s", err,
			objID, instIDs, kit.Rid)
		return svcInstCount, err
	}

	if len(zeroIndexMap) == 0 {
		return svcInstCount, nil
	}

	index := 0
	for i := 0; i < len(instIDs); i++ {
		if _, exists := zeroIndexMap[i]; exists {
			results[i] = 0
			continue
		}
		results[i] = svcInstCount[index]
		index++
	}
	return results, nil
}

// getDistinctHostCount get distinct host count
func (assoc *association) getDistinctHostCount(kit *rest.Kit, objID string, instIDs []int64) (int64, error) {
	opt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameModuleHostConfig,
		Field:     common.BKHostIDField,
		Filter:    mapstr.MapStr{objID: mapstr.MapStr{common.BKDBIN: instIDs}},
	}

	count, err := assoc.clientSet.CoreService().Common().GetDistinctCount(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("find distinct host count failed, err: %v, objID: %s, instIDs: %, rid: %s", err,
			objID, instIDs, kit.Rid)
		return count, err
	}

	return count, nil
}

// getSetIDsByTopo get set IDs by custom layer node
func (assoc *association) getSetIDsByTopo(kit *rest.Kit, objID string, instIDs []int64) ([]int64, error) {

	if objID == common.BKInnerObjIDApp || objID == common.BKInnerObjIDSet || objID == common.BKInnerObjIDModule {
		blog.Errorf("get set IDs by topo failed, obj(%s) is a inner object, rid: %s", objID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjIDField)
	}

	// get mainline association, generate map of object and its child
	asstRes, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{
			Condition: map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline}})
	if err != nil {
		blog.Errorf("get set IDs by topo failed, get mainline association err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	childObjMap := make(map[string]string)
	for _, asst := range asstRes.Info {
		childObjMap[asst.AsstObjID] = asst.ObjectID
	}

	childObj := childObjMap[objID]
	if childObj == "" {
		blog.Errorf("get set IDs by topo failed, obj(%s) is not a mainline object, rid: %s", objID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjIDField)
	}

	// traverse down topo till set, get set ids
	for {
		idField := common.GetInstIDField(childObj)
		instCond := make(map[string]interface{})
		instCond[common.BKParentIDField] = map[string]interface{}{
			common.BKDBIN: instIDs,
		}
		// exclude default sets
		if childObj == common.BKInnerObjIDSet {
			instCond[common.BKDefaultField] = map[string]interface{}{
				common.BKDBNE: common.DefaultResSetFlag,
			}
		}
		query := &metadata.QueryCondition{
			Condition: instCond,
			Fields:    []string{idField},
			Page:      metadata.BasePage{Limit: common.BKNoLimit},
		}

		instRes, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, childObj, query)
		if err != nil {
			blog.Errorf("get set IDs by topo failed, read instance err: %s, objID: %s, instIDs: %+v, rid: %s",
				err.Error(), childObj, instIDs, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if len(instRes.Info) == 0 {
			return []int64{}, nil
		}

		instIDs = make([]int64, len(instRes.Info))
		for index, insts := range instRes.Info {
			id, err := insts.Int64(idField)
			if err != nil {
				blog.Errorf("get set IDs by topo failed, parse inst id err: %s, inst: %#v, rid: %s", err.Error(),
					insts, kit.Rid)
				return nil, err
			}
			instIDs[index] = id
		}

		if childObj == common.BKInnerObjIDSet {
			break
		}
		childObj = childObjMap[childObj]
	}

	return instIDs, nil
}

// getSetRelationModule get set relation module by set ids
func (assoc *association) getSetRelationModule(kit *rest.Kit, setIDs []int64) (map[int64][]int64, error) {
	queryCond := &metadata.QueryCondition{
		Fields:         []string{common.BKSetIDField, common.BKModuleIDField},
		Condition:      mapstr.MapStr{common.BKSetIDField: mapstr.MapStr{common.BKDBIN: setIDs}},
		DisableCounter: true,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}

	resp, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		queryCond)
	if err != nil {
		blog.Errorf("get instance data failed, error info is %s , rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	setRelModuleMap := make(map[int64][]int64)
	for _, mapStr := range resp.Info {
		setID, err := mapStr.Int64(common.BKSetIDField)
		if err != nil {
			blog.Errorf("failed to parse the interface to int64, error info is %s , rid: %s",
				err.Error(), kit.Rid)
			return nil, err
		}

		moduleID, err := mapStr.Int64(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("failed to parse the interface to int64, error info is %s , rid: %s",
				err.Error(), kit.Rid)
			return nil, err
		}

		setRelModuleMap[setID] = append(setRelModuleMap[setID], moduleID)
	}

	return setRelModuleMap, nil
}
