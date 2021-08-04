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
	"regexp"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	// SearchMainlineAssociationInstTopo search mainline association topo by objID and instID
	SearchMainlineAssociationInstTopo(kit *rest.Kit, objID string, instID int64,
		withStatistics bool, withDefault bool) ([]*metadata.TopoInstRst, errors.CCError)
	// ResetMainlineInstAssociation reset mainline instance association
	ResetMainlineInstAssociation(kit *rest.Kit, currentObjID string) error
	// SetMainlineInstAssociation set mainline instance association by parent object and current object
	SetMainlineInstAssociation(kit *rest.Kit, parentObjID string,
		currObjID string, currObjName string) ([]int64, error)
}

// NewAssociationOperation create a new association operation instance
func NewAssociationOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) AssociationOperationInterface {
	return &association{
		clientSet:   client,
		authManager: authManager,
	}
}

type association struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

// ResetMainlineInstAssociation reset mainline instance association
// while a new mainline object has been created may use this func
func (assoc *association) ResetMainlineInstAssociation(kit *rest.Kit, currentObjID string) error {

	cond := mapstr.New()
	if metadata.IsCommon(currentObjID) {
		cond.Set(common.BKObjIDField, currentObjID)
	}
	// TODO 确认是否需要设置分页
	defaultCond := &metadata.QueryInput{
		Condition: cond,
		Fields:    fmt.Sprintf("%s,%s", common.GetInstIDField(currentObjID), common.BKParentIDField),
	}

	// 获取 current 模型的所有实例
	currentInsts, err := assoc.FindInst(kit, currentObjID, defaultCond)
	if err != nil {
		blog.Errorf("failed to find current object(%s) inst, err: %v, rid: %s", currentObjID, err, kit.Rid)
		return err
	}

	var currentInstIDs []int64
	for _, currentInst := range currentInsts.Info {
		instID, err := currentInst.Int64(common.GetInstIDField(currentObjID))
		if err != nil {
			blog.Errorf("failed to get the inst id from the inst(%#v), rid: %s", currentInst, kit.Rid)
			continue
		}

		currentInstIDs = append(currentInstIDs, instID)
	}

	childObjID, childInst, err := assoc.getMainlineNodeInst(kit, currentObjID, currentInstIDs, true)
	if err != nil {
		blog.Errorf("failed to get the object(%s) mainline child, err: %v, rid: %s",
			currentObjID, err, kit.Rid)
		return err
	}

	currentChildMap := make(map[int64][]int64)
	childNameMap := make(map[string]bool)
	for _, child := range childInst {
		instID, err := child.Int64(common.GetInstIDField(childObjID))
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

		currentChildMap[parentID] = append(currentChildMap[parentID], instID)

		instName, err := child.String(common.GetInstNameField(childObjID))
		if err != nil {
			blog.Errorf("failed to get the inst name from the inst(%#v), rid: %s", child, kit.Rid)
			continue
		}

		// 检查实例删除后，会不会出现重名冲突
		if exist, _ := childNameMap[instName]; exist {
			blog.Errorf("can not be reset, inst name repeated, inst: %s, rid: %s", instName, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrorTopoMultipleObjectInstanceName, instName)
		}

		childNameMap[instName] = true
	}

	// 修改 currentInsts 所有孩子结点的父节点，为 currentInsts 的父节点，并删除 currentInsts
	parentChildMap := make(map[int64][]int64)
	for _, currentInst := range currentInsts.Info {
		instID, err := currentInst.Int64(common.GetInstIDField(currentObjID))
		if err != nil {
			blog.Errorf("failed to get the inst id from the inst(%#v), rid: %s", currentInst, kit.Rid)
			continue
		}

		parentID, err := currentInst.Int64(common.BKParentIDField)
		if err != nil {
			blog.Errorf("failed to get the object(%s) mainline parent id, "+
				"the current inst(%v), err: %v, rid: %s", currentObjID, currentInst, err, kit.Rid)
			continue
		}

		parentChildMap[parentID] = append(parentChildMap[parentID], currentChildMap[instID]...)
	}

	for parentID, childIDs := range parentChildMap {
		if err = assoc.SetMainlineParentInst(kit, childIDs, childObjID, parentID); err != nil {
			blog.Errorf("failed to set the object mainline child inst, err: %v, rid: %s", err, kit.Rid)
			continue
		}
	}

	if err := assoc.deleteMainlineInstWithID(kit, currentObjID, currentInstIDs); err != nil {
		blog.Errorf("failed to delete the current inst(%#v), err: %v, rid: %s",
			currentInstIDs, err, kit.Rid)
		return err
	}

	return nil
}

// SetMainlineInstAssociation set mainline instance association by parent object and current object
func (assoc *association) SetMainlineInstAssociation(kit *rest.Kit, parentObjID, currObjID,
	currObjName string) ([]int64, error) {

	defaultCond := &metadata.QueryInput{}
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
	re, _ := regexp.Compile(`[#/,><|]`)
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
	childObjID, children, err := assoc.getMainlineNodeInst(kit, parentObjID, parentInstIDs, true)
	if err != nil {
		blog.Errorf("failed to get the object(%s) mainline child inst, err: %v, rid: %s",
			parentObjID, err, kit.Rid)
		return nil, err
	}

	for _, child := range children {
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
		if err = assoc.SetMainlineParentInst(kit, childIDs, childObjID, parentID); err != nil {
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
	mainlineAsstRsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: map[string]interface{}{
			common.AssociationKindIDField: common.AssociationKindMainline,
		}})
	if err != nil {
		blog.Errorf("search mainline association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	mainlineObjectChildMap := make(map[string]string)
	isMainline := false
	for _, asst := range mainlineAsstRsp.Data.Info {
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
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = objects.CCError(); err != nil {
		blog.Errorf("search mainline objects(%s) failed, err: %V, rid: %s", objectIDs, err, kit.Rid)
		return nil, err
	}

	for _, object := range objects.Data.Info {
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
		filter := &metadata.QueryInput{Condition: instCond}
		if objectID != objID {
			filter.Sort = common.BKDefaultField
		}

		if objectID == common.BKInnerObjIDSet {
			filter.Fields = fmt.Sprintf("%s,%s,%s,%s,%s,%s",
				common.GetInstIDField(objectID), common.GetInstNameField(objectID),
				common.BKDefaultField, common.BKSetTemplateIDField, common.BKParentIDField, common.BKAppIDField)
		} else if objectID == common.BKInnerObjIDModule {
			filter.Fields = fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s",
				common.GetInstIDField(objectID), common.GetInstNameField(objectID),
				common.BKDefaultField, common.BKServiceTemplateIDField, common.BKSetTemplateIDField,
				common.HostApplyEnabledField, common.BKParentIDField, common.BKAppIDField)
		} else {
			filter.Fields = fmt.Sprintf("%s,%s,%s,%s,%s",
				common.GetInstIDField(objectID), common.GetInstNameField(objectID),
				common.BKDefaultField, common.BKParentIDField, common.BKAppIDField)
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
					blog.Errorf("get instance %#v default failed, err: %v, rid: %s",
						instance, err, kit.Rid)
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
						blog.Errorf("get instance %#v biz id failed, err: %v, rid: %s",
							instance, err, kit.Rid)
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
			blog.Errorf("fill statistics data failed, bizID: %d, err: %v, rid: %s",
				bizID, err, kit.Rid)
			return nil, err
		}
	}
	return results, nil
}

func (assoc *association) getMainlineNodeInst(kit *rest.Kit, objID string, instIDs []int64,
	needChild bool) (string, []mapstr.MapStr, error) {

	objCond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}
	if needChild {
		objCond.Set(common.BKAsstObjIDField, objID)
	} else {
		objCond.Set(common.BKObjIDField, objID)
	}

	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: objCond})
	if err != nil {
		blog.Errorf("search object(%s) mainline association failed, err: %v, rid: %s",
			objID, err, kit.Rid)
		return "", nil, err
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("search object(%s) mainline association failed, err: %v, rid: %s",
			objID, err, kit.Rid)
		return "", nil, err
	}

	if len(rsp.Data.Info) != 1 {
		blog.Errorf("search object mainline association failed, object(%s) isn't mainline object", objID)
		return "", nil, kit.CCError.CCError(common.CCErrTopoMainlineSelectFailed)
	}

	var mainlineObjID string
	if needChild {
		mainlineObjID = rsp.Data.Info[0].ObjectID
	} else {
		mainlineObjID = rsp.Data.Info[0].AsstObjID
	}

	cond := mapstr.New()
	if metadata.IsCommon(mainlineObjID) {
		cond.Set(metadata.ModelFieldObjectID, mainlineObjID)
	} else if mainlineObjID == common.BKInnerObjIDSet {
		cond.Set(common.BKDefaultField, mapstr.MapStr{common.BKDBNE: common.DefaultResSetFlag})
	}

	if needChild {
		cond.Set(common.BKInstParentStr, mapstr.MapStr{common.BKDBIN: instIDs})
	} else {
		cond.Set(common.GetInstIDField(mainlineObjID), mapstr.MapStr{common.BKDBIN: instIDs})
	}

	instRsp, err := assoc.FindInst(kit, mainlineObjID,
		&metadata.QueryInput{
			Condition: cond,
			Fields:    fmt.Sprintf("%s,%s", common.GetInstIDField(mainlineObjID), common.BKParentIDField),
		})
	if err != nil {
		blog.Errorf("search inst by object(%s) failed, err: %v, rid: %s", mainlineObjID, err, kit.Rid)
		return "", nil, err
	}

	return mainlineObjID, instRsp.Info, nil
}

// FindInst search instance by condition
// TODO need to delete after merge
func (assoc *association) FindInst(kit *rest.Kit, objID string,
	cond *metadata.QueryInput) (*metadata.InstResult, error) {

	result := new(metadata.InstResult)
	switch objID {
	case common.BKInnerObjIDHost:
		rsp, err := assoc.clientSet.CoreService().Host().GetHosts(kit.Ctx, kit.Header, cond)
		if err != nil {
			blog.Errorf("get host failed, err: %v, rid: %s", err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if err = rsp.CCError(); err != nil {
			blog.Errorf("search object(%s) inst by the condition(%#v) failed, err: %v, rid: %s",
				objID, cond, err, kit.Rid)
			return nil, err
		}

		result.Count = rsp.Data.Count
		result.Info = rsp.Data.Info
		return result, nil

	default:
		input := &metadata.QueryCondition{Condition: cond.Condition, TimeCondition: cond.TimeCondition}
		input.Page.Start = cond.Start
		input.Page.Limit = cond.Limit
		input.Page.Sort = cond.Sort
		input.Fields = strings.Split(cond.Fields, ",")
		rsp, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, input)
		if err != nil {
			blog.Errorf("search instance failed, err: %v, rid: %s", err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if err = rsp.CCError(); err != nil {
			blog.Errorf("search object(%s) inst by the condition(%#v) failed, err: %v, rid: %s",
				objID, cond, err, kit.Rid)
			return nil, err
		}

		result.Count = rsp.Data.Count
		result.Info = rsp.Data.Info
		return result, nil
	}
}

// TODO should move to another go file after merge
func (assoc *association) SetMainlineParentInst(kit *rest.Kit, childInstID []int64, childObjID string,
	instID int64) error {
	if err := assoc.updateMainlineAssociation(kit, childInstID, childObjID, instID); err != nil {
		blog.Errorf("failed to update the mainline association, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (assoc *association) updateMainlineAssociation(kit *rest.Kit, childID []int64, childObjID string,
	parentID int64) error {

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
	rsp, err := assoc.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, childObjID, &input)
	if err != nil {
		blog.Errorf("failed to request object controller, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to update the association, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// TODO should be deleted after merge, and which call this func use DeleteMainlineInstWithID in inst/inst.go to replace
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

	if err = cnt.CCError(); err != nil {
		blog.Errorf("count association by object(%s) failed, err: %s, rid: %s", objID, err, kit.Rid)
		return err
	}

	if cnt.Data.Count > 0 {
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
	rsp, err := assoc.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, objID, &ops)
	if err != nil {
		blog.Errorf("request to delete instance by condition failed, cond: %#v, err: %v", ops, err)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to delete the object(%s) inst by the condition(%#v), err: %v",
			objID, ops, err)
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

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to create object instance ,err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}

	return rsp.Data.Created.ID, nil
}

// TODO need check this function is here or move to other file
func (assoc *association) fillStatistics(kit *rest.Kit, bizID int64,
	moduleIDs []int64, topoInsts []*metadata.TopoInstRst) errors.CCError {

	// get service instance count
	option := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	serviceInstances, ccErr := assoc.clientSet.CoreService().Process().ListServiceInstance(kit.Ctx, kit.Header, option)
	if ccErr != nil {
		blog.Errorf("list service instances failed, option: %+v, err: %v, rid: %s",
			option, ccErr, kit.Rid)
		return ccErr
	}
	moduleServiceInstanceCount := make(map[int64]int64)
	for _, serviceInstance := range serviceInstances.Info {
		moduleServiceInstanceCount[serviceInstance.ModuleID]++
	}

	// get host count
	listHostOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		Fields:        []string{common.BKAppIDField, common.BKSetIDField, common.BKModuleIDField, common.BKHostIDField},
	}
	hostModules, err := assoc.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, listHostOption)
	if err != nil {
		blog.Errorf("list host modules failed, option: %+v, err: %v, rid: %s",
			listHostOption, err, kit.Rid)
		return err
	}
	// topoObjectID -> topoInstanceID -> []hostIDs
	customLevel := "custom_level"
	hostCount := make(map[string]map[int64][]int64)
	hostCount[common.BKInnerObjIDApp] = make(map[int64][]int64)
	hostCount[common.BKInnerObjIDSet] = make(map[int64][]int64)
	hostCount[common.BKInnerObjIDModule] = make(map[int64][]int64)
	hostCount[customLevel] = make(map[int64][]int64)
	for _, hostModule := range hostModules.Data.Info {
		if _, exist := hostCount[common.BKInnerObjIDModule][hostModule.ModuleID]; exist == false {
			hostCount[common.BKInnerObjIDModule][hostModule.ModuleID] = make([]int64, 0)
		}
		hostCount[common.BKInnerObjIDModule][hostModule.ModuleID] = append(
			hostCount[common.BKInnerObjIDModule][hostModule.ModuleID], hostModule.HostID)

		if _, exist := hostCount[common.BKInnerObjIDSet][hostModule.SetID]; exist == false {
			hostCount[common.BKInnerObjIDSet][hostModule.SetID] = make([]int64, 0)
		}
		hostCount[common.BKInnerObjIDSet][hostModule.SetID] = append(
			hostCount[common.BKInnerObjIDSet][hostModule.SetID], hostModule.HostID)

		if _, exist := hostCount[common.BKInnerObjIDApp][hostModule.AppID]; exist == false {
			hostCount[common.BKInnerObjIDApp][hostModule.AppID] = make([]int64, 0)
		}
		hostCount[common.BKInnerObjIDApp][hostModule.AppID] = append(
			hostCount[common.BKInnerObjIDApp][hostModule.AppID], hostModule.HostID)
	}
	for _, objectID := range []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule} {
		for key := range hostCount[objectID] {
			hostCount[objectID][key] = util.IntArrayUnique(hostCount[objectID][key])
		}
	}

	// get host apply rule count
	listApplyRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: moduleIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostApplyRules, err := assoc.clientSet.CoreService().
		HostApplyRule().ListHostApplyRule(kit.Ctx, kit.Header, bizID, listApplyRuleOption)
	if err != nil {
		blog.Errorf("list host apply rule failed, bizID: %s, option: %#v, err: %v, rid: %s",
			bizID, listApplyRuleOption, err, kit.Rid)
		return err
	}
	moduleRuleCount := make(map[int64]int64)
	for _, item := range hostApplyRules.Info {
		moduleRuleCount[item.ModuleID]++
	}

	exactNodes := []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule}
	// fill hosts
	for _, tir := range topoInsts {
		tir.DeepFirstTraverse(func(node *metadata.TopoInstRst) {
			// calculate service instance count
			subTreeSvcInstCount := int64(0)
			for _, child := range node.Child {
				subTreeSvcInstCount += child.ServiceInstanceCount
			}
			node.ServiceInstanceCount = subTreeSvcInstCount
			if node.ObjID == common.BKInnerObjIDModule {
				if _, exist := moduleServiceInstanceCount[node.InstID]; exist == true {
					node.ServiceInstanceCount = moduleServiceInstanceCount[node.InstID]
				}
				node.HostApplyRuleCount = new(int64)
				*node.HostApplyRuleCount, _ = moduleRuleCount[node.InstID]
			}

			if util.InStrArr(exactNodes, node.ObjID) {
				if _, exist := hostCount[node.ObjID][node.InstID]; exist == true {
					node.HostCount = int64(len(hostCount[node.ObjID][node.InstID]))
				}
				return
			}
			if len(node.Child) == 0 {
				return
			}

			// calculate host count
			subTreeHosts := make([]int64, 0)
			for _, child := range node.Child {
				childHosts := make([]int64, 0)
				if util.InStrArr(exactNodes, child.ObjID) {
					if _, exist := hostCount[child.ObjID][child.InstID]; exist == true {
						childHosts = hostCount[child.ObjID][child.InstID]
					}
				} else {
					if _, exist := hostCount[customLevel][child.InstID]; exist == true {
						childHosts = hostCount[customLevel][child.InstID]
					}
				}
				subTreeHosts = append(subTreeHosts, childHosts...)
			}
			hostCount[customLevel][node.InstID] = util.IntArrayUnique(subTreeHosts)
			node.HostCount = int64(len(hostCount[customLevel][node.InstID]))
		})
	}
	return nil
}
