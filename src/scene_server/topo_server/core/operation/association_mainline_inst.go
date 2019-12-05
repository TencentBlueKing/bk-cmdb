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
	"io"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// checkInstNameRepeat 检查如果将 currentInsts 都删除之后，拥有共同父节点的孩子结点会不会出现名字冲突
// 如果有冲突，返回 (false, 冲突实例名, nil)
func (assoc *association) checkInstNameRepeat(params types.ContextParams, currentInsts []inst.Inst) (canReset bool, repeatedInstName string, err error) {
	// TODO: 返回值将bool值与出错情况分开 (bool, error)
	instNames := map[string]bool{}
	for _, currInst := range currentInsts {
		currInstParentID, err := currInst.GetParentID()
		if nil != err {
			return false, "", err
		}

		children, err := currInst.GetMainlineChildInst()
		if nil != err {
			return false, "", err
		}

		for _, child := range children {
			instName, err := child.GetInstName()
			if nil != err {
				return false, "", err
			}
			key := fmt.Sprintf("%d_%s", currInstParentID, instName)
			if _, ok := instNames[key]; ok {
				return false, instName, nil
			}

			instNames[key] = true
		}
	}

	return true, "", nil
}

func (assoc *association) ResetMainlineInstAssociation(params types.ContextParams, current model.Object) error {
	rid := params.ReqID

	cObj := current.Object()
	cond := condition.CreateCondition()
	if current.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(cObj.ObjectID)
	}
	defaultCond := &metadata.QueryInput{}
	defaultCond.Condition = cond.ToMapStr()

	// 获取 current 模型的所有实例
	_, currentInsts, err := assoc.inst.FindInst(params, current, defaultCond, false)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find current object(%s) inst, err: %+v, rid: %s", cObj.ObjectID, err, rid)
		return err
	}

	// 检查实例删除后，会不会出现重名冲突
	var canReset bool
	var repeatedInstName string
	if canReset, repeatedInstName, err = assoc.checkInstNameRepeat(params, currentInsts); nil != err {
		blog.Errorf("[operation-asst] can not be reset, err: %+v, rid: %s", err, rid)
		return err
	}
	if canReset == false {
		blog.Errorf("[operation-asst] can not be reset, inst name repeated, inst: %s, rid: %s", repeatedInstName, rid)
		errMsg := params.Err.Error(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat).Error() + " " + repeatedInstName
		return params.Err.New(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat, errMsg)
	}

	// NEED FIX: 下面循环中的continue ，会在处理实例异常的时候跳过当前拓扑的处理，此方式可能会导致某个业务拓扑失败，但是不会影响所有。
	// 修改 currentInsts 所有孩子结点的父节点，为 currentInsts 的父节点，并删除 currentInsts
	for _, currentInst := range currentInsts {
		instID, err := currentInst.GetInstID()
		if nil != err {
			blog.Errorf("[operation-asst] failed to get the inst id from the inst(%#v), rid: %s", currentInst.ToMapStr(), rid)
			continue
		}

		parentID, err := currentInst.GetParentID()
		if nil != err {
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline parent id, the current inst(%v), err: %+v, rid: %s", cObj.ObjectID, currentInst.GetValues(), err, rid)
			continue
		}

		// reset the child's parent
		children, err := currentInst.GetMainlineChildInst()
		if nil != err {
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline child inst, err: %+v, rid: %s", cObj.ObjectID, err, rid)
			continue
		}
		for _, child := range children {
			// set the child's parent
			if err = child.SetMainlineParentInst(parentID); nil != err {
				blog.Errorf("[operation-asst] failed to set the object(%s) mainline child inst, err: %+v, rid: %s", child.GetObject().Object().ObjectID, err, rid)
				continue
			}
		}

		// delete the current inst
		if err := assoc.inst.DeleteMainlineInstWithID(params, current, instID); nil != err {
			blog.Errorf("[operation-asst] failed to delete the current inst(%#v), err: %+v, rid: %s", currentInst.ToMapStr(), err, rid)
			continue
		}
	}

	return nil
}

func (assoc *association) SetMainlineInstAssociation(params types.ContextParams, parent, current, child model.Object) error {

	defaultCond := &metadata.QueryInput{}
	cond := condition.CreateCondition()
	if parent.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(parent.Object().ObjectID)
	}
	defaultCond.Condition = cond.ToMapStr()
	// fetch all parent instances.
	_, parentInsts, err := assoc.inst.FindInst(params, parent, defaultCond, false)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find parent object(%s) inst, err: %s, rid: %s", parent.Object().ObjectID, err.Error(), params.ReqID)
		return err
	}

	expectParent2Children := map[int64][]inst.Inst{}
	// create current object instance for each parent instance and insert the current instance to
	for _, parent := range parentInsts {

		id, err := parent.GetInstID()
		if nil != err {
			blog.Errorf("[operation-asst] failed to find the inst id, err: %s, rid: %s", err.Error(), params.ReqID)
			return err
		}

		// we create the current object's instance for each parent instance belongs to the parent object.
		currentInst := assoc.instFactory.CreateInst(params, current)
		currentInst.SetValue(current.GetInstNameFieldName(), current.Object().ObjectName)
		currentInst.SetValue(common.BKDefaultField, common.DefaultFlagDefaultValue)
		// set current instance's parent id to parent instance's id, so that they can be chained.
		currentInst.SetValue(common.BKInstParentStr, id)
		object := parent.GetObject()
		if object.GetObjectID() == common.BKInnerObjIDApp {
			metaInfo := metadata.NewMetaDataFromBusinessID(strconv.FormatInt(id, 10))
			currentInst.SetValue(metadata.BKMetadata, metaInfo)
		} else {
			currentInst.SetValue(metadata.BKMetadata, parent.GetValues()[metadata.BKMetadata])
		}

		// create the instance now.
		if err = currentInst.Create(); nil != err {
			blog.Errorf("[operation-asst] failed to create object(%s) default inst, err: %s, rid: %s", current.Object().ObjectID, err.Error(), params.ReqID)
			return err
		}
		instID, err := currentInst.GetInstID()
		if err != nil {
			blog.Errorf("create mainline instance for obj: %s, but got invalid instance id, err :%v, rid: %s", current.Object().ObjectID, err, params.ReqID)
			return err
		}
		err = assoc.authManager.RegisterInstancesByID(params.Context, params.Header, current.Object().ObjectID, instID)
		if err != nil {
			blog.Errorf("create mainline instance for object: %s, but register to auth center failed, instID: %d, err: %v, rid: %s", current.Object().ObjectID, instID, err, params.ReqID)
			return err
		}

		// reset the child's parent instance's parent id to current instance's id.
		children, err := parent.GetMainlineChildInst()
		if nil != err {
			if io.EOF == err {
				continue
			}
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline child inst, err: %s, rid: %s", parent.GetObject().Object().ObjectID, err.Error(), params.ReqID)
			return err
		}

		curInstID, err := currentInst.GetInstID()
		if err != nil {
			blog.Errorf("[operation-asst] failed to get the instID(%#v), err: %s, rid: %s", currentInst.ToMapStr(), err.Error(), params.ReqID)
			return err
		}

		expectParent2Children[curInstID] = children
	}

	for parentID, children := range expectParent2Children {
		for _, child := range children {
			// set the child's parent
			if err = child.SetMainlineParentInst(parentID); nil != err {
				blog.Errorf("[operation-asst] failed to set the object(%s) mainline child inst, err: %s, rid: %s", child.GetObject().Object().ObjectID, err.Error(), params.ReqID)
				return err
			}
		}
	}

	return nil
}

func (assoc *association) SearchMainlineAssociationInstTopo(params types.ContextParams, obj model.Object, instID int64, withStatistics bool) ([]*metadata.TopoInstRst, error) {

	cond := &metadata.QueryInput{}
	cond.Condition = mapstr.MapStr{
		obj.GetInstIDFieldName(): instID,
	}

	_, bizInsts, err := assoc.inst.FindInst(params, obj, cond, false)
	if nil != err {
		blog.Errorf("[SearchMainlineAssociationInstTopo] FindInst for %+v failed: %v, rid: %s", cond, err, params.ReqID)
		return nil, err
	}

	results := make([]*metadata.TopoInstRst, 0)
	for _, biz := range bizInsts {
		instID, err := biz.GetInstID()
		if nil != err {
			blog.Errorf("[SearchMainlineAssociationInstTopo] GetInstID for %+v failed: %v, rid: %s", biz, err, params.ReqID)
			return nil, err
		}
		instName, err := biz.GetInstName()
		if nil != err {
			blog.Errorf("[SearchMainlineAssociationInstTopo] GetInstName for %+v failed: %v, rid: %s", biz, err, params.ReqID)
			return nil, err
		}

		object := biz.GetObject().Object()
		tmp := &metadata.TopoInstRst{Child: []*metadata.TopoInstRst{}}
		tmp.InstID = instID
		tmp.InstName = instName
		tmp.ObjID = object.ObjectID
		tmp.ObjName = object.ObjectName

		results = append(results, tmp)
	}

	if err = assoc.fillMainlineChildInst(params, obj, results); err != nil {
		blog.Errorf("[SearchMainlineAssociationInstTopo] fillMainlineChildInst for %+v failed: %v, rid: %s", results, err, params.ReqID)
		return nil, err
	}
	if withStatistics && len(bizInsts) > 0 {
		instance := bizInsts[0]
		bizID, err := instance.GetBizID()
		if err != nil {
			blog.Errorf("[SearchMainlineAssociationInstTopo] parse biz id failed, inst: %+v, err: %v, rid: %s", instance, err, params.ReqID)
			return nil, params.Err.CCError(common.CCErrCommParseBizIDFromMetadataInDBFailed)
		}
		if err := assoc.fillStatistics(params, bizID, results); err != nil {
			blog.Errorf("[SearchMainlineAssociationInstTopo] fill statistics data failed, bizID: %d, err: %v, rid: %s", bizID, err, params.ReqID)
			return nil, err
		}
	}
	return results, nil
}

func (assoc *association) fillStatistics(params types.ContextParams, bizID int64, parentInsts []*metadata.TopoInstRst) error {
	// fill service instance count
	option := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	serviceInstances, err := assoc.clientSet.CoreService().Process().ListServiceInstance(params.Context, params.Header, option)
	if err != nil {
		blog.Errorf("fillStatistics failed, list service instances failed, option: %+v, err: %s, rid: %s", option, err.Error(), params.ReqID)
		return err
	}
	moduleServiceInstanceCount := map[int64]int64{}
	for _, serviceInstance := range serviceInstances.Info {
		if _, exist := moduleServiceInstanceCount[serviceInstance.ModuleID]; exist == false {
			moduleServiceInstanceCount[serviceInstance.ModuleID] = 0
		}
		moduleServiceInstanceCount[serviceInstance.ModuleID]++
	}
	listHostOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
	}
	hostModules, e := assoc.clientSet.CoreService().Host().GetHostModuleRelation(params.Context, params.Header, listHostOption)
	if e != nil {
		blog.Errorf("fillStatistics failed, list host modules failed, option: %+v, err: %s, rid: %s", listHostOption, e.Error(), params.ReqID)
		return e
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
		hostCount[common.BKInnerObjIDModule][hostModule.ModuleID] = append(hostCount[common.BKInnerObjIDModule][hostModule.ModuleID], hostModule.HostID)

		if _, exist := hostCount[common.BKInnerObjIDSet][hostModule.SetID]; exist == false {
			hostCount[common.BKInnerObjIDSet][hostModule.SetID] = make([]int64, 0)
		}
		hostCount[common.BKInnerObjIDSet][hostModule.SetID] = append(hostCount[common.BKInnerObjIDSet][hostModule.SetID], hostModule.HostID)

		if _, exist := hostCount[common.BKInnerObjIDApp][hostModule.AppID]; exist == false {
			hostCount[common.BKInnerObjIDApp][hostModule.AppID] = make([]int64, 0)
		}
		hostCount[common.BKInnerObjIDApp][hostModule.AppID] = append(hostCount[common.BKInnerObjIDApp][hostModule.AppID], hostModule.HostID)
	}
	for _, objectID := range []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule} {
		for key := range hostCount[objectID] {
			hostCount[objectID][key] = util.IntArrayUnique(hostCount[objectID][key])
		}
	}

	// module bound service_templateID
	moduleFilter := map[string]interface{}{
		common.BKAppIDField: bizID,
	}
	moduleQueryCondition := &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr(moduleFilter),
	}
	modules, e := assoc.clientSet.CoreService().Instance().ReadInstance(params.Context, params.Header, common.BKInnerObjIDModule, moduleQueryCondition)
	if e != nil {
		blog.Errorf("fillStatistics failed, list modules failed, option: %+v, err: %s, rid: %s", listHostOption, e.Error(), params.ReqID)
		return e
	}
	moduleServiceTemplateIDMap := make(map[int64]int64)
	moduleSetTemplateIDMap := make(map[int64]int64)
	setSetTemplateIDMap := make(map[int64]int64)
	hostApplyEnabledMap := make(map[int64]bool)
	moduleIDs := make([]int64, 0)
	for _, module := range modules.Data.Info {
		moduleStruct := &metadata.ModuleInst{}
		if err := module.ToStructByTag(moduleStruct, "field"); err != nil {
			blog.Errorf("fillStatistics failed, parse module data to struct failed, module: %+v, err: %s, rid: %s", module, e.Error(), params.ReqID)
			return err
		}
		moduleIDs = append(moduleIDs, moduleStruct.ModuleID)
		moduleServiceTemplateIDMap[moduleStruct.ModuleID] = moduleStruct.ServiceTemplateID
		moduleSetTemplateIDMap[moduleStruct.ModuleID] = moduleStruct.SetTemplateID
		setSetTemplateIDMap[moduleStruct.ParentID] = moduleStruct.SetTemplateID
		hostApplyEnabledMap[moduleStruct.ModuleID] = moduleStruct.HostApplyEnabled
	}

	exactNodes := []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule}
	// fill service instances
	for _, tir := range parentInsts {
		tir.DeepFirstTraverse(func(node *metadata.TopoInstRst) {
			if len(node.Child) > 0 {
				// calculate service instance count
				subTreeSvcInstCount := int64(0)
				for _, child := range node.Child {
					subTreeSvcInstCount += child.ServiceInstanceCount
				}
				node.ServiceInstanceCount = subTreeSvcInstCount
			}
			if node.ObjID == common.BKInnerObjIDModule {
				if _, exist := moduleServiceInstanceCount[node.InstID]; exist == true {
					node.ServiceInstanceCount = moduleServiceInstanceCount[node.InstID]
				}
				if id, exist := moduleServiceTemplateIDMap[node.InstID]; exist == true {
					node.ServiceTemplateID = id
				}
				if id, exist := moduleSetTemplateIDMap[node.InstID]; exist == true {
					node.SetTemplateID = id
				}
				node.HostApplyEnabled = new(bool)
				if enabled, exist := hostApplyEnabledMap[node.InstID]; exist == true {
					*node.HostApplyEnabled = enabled
				}
			}
			if node.ObjID == common.BKInnerObjIDSet {
				if id, exist := setSetTemplateIDMap[node.InstID]; exist == true {
					node.SetTemplateID = id
				}
			}
		})
	}
	// fill hosts
	for _, tir := range parentInsts {
		tir.DeepFirstTraverse(func(node *metadata.TopoInstRst) {
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

	// fill host apply rules
	listApplyRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: moduleIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostApplyRules, err := assoc.clientSet.CoreService().HostApplyRule().ListHostApplyRule(params.Context, params.Header, bizID, listApplyRuleOption)
	if err != nil {
		blog.Errorf("fillStatistics failed, ListHostApplyRule failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, listApplyRuleOption, err, params.ReqID)
		return err
	}
	moduleRuleCount := make(map[int64]int64)
	for _, item := range hostApplyRules.Info {
		if _, exist := moduleRuleCount[item.ModuleID]; exist == false {
			moduleRuleCount[item.ModuleID] = 0
		}
		moduleRuleCount[item.ModuleID] += 1
	}
	for _, tir := range parentInsts {
		tir.DeepFirstTraverse(func(node *metadata.TopoInstRst) {
			if node.ObjID == common.BKInnerObjIDModule {
				node.HostApplyRuleCount = new(int64)
				*node.HostApplyRuleCount, _ = moduleRuleCount[node.InstID]
			}
		})
	}

	return nil
}

func (assoc *association) fillMainlineChildInst(params types.ContextParams, object model.Object, parentInsts []*metadata.TopoInstRst) error {
	childObj, err := object.GetMainlineChildObject()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		blog.Errorf("[fillMainlineChildInst] GetMainlineChildObject for %+v failed: %v, rid: %s", object, err, params.ReqID)
		return err
	}

	parentIDs := make([]int64, 0)
	for index := range parentInsts {
		parentIDs = append(parentIDs, parentInsts[index].InstID)
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKInstParentStr).In(parentIDs)
	if childObj.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(childObj.Object().ObjectID)
	} else if childObj.Object().ObjectID == common.BKInnerObjIDSet {
		cond.Field(common.BKDefaultField).NotEq(common.DefaultResSetFlag)
	}

	_, childInsts, err := assoc.inst.FindInst(params, childObj, &metadata.QueryInput{Condition: cond.ToMapStr()}, false)
	if err != nil {
		blog.Errorf("[fillMainlineChildInst] FindInst for %+v failed: %v, rid: %s", cond.ToMapStr(), err, params.ReqID)
		return err
	}

	// parentID mapping to child topo instances
	childInstMap := map[int64][]*metadata.TopoInstRst{}
	childTopoInsts := make([]*metadata.TopoInstRst, 0)
	for _, childInst := range childInsts {
		childInstID, err := childInst.GetInstID()
		if err != nil {
			blog.Errorf("[fillMainlineChildInst] GetInstID for %+v failed: %v, rid: %s", childInst, err, params.ReqID)
			return err
		}
		childInstName, err := childInst.GetInstName()
		if nil != err {
			blog.Errorf("[fillMainlineChildInst] GetInstName for %+v failed: %v, rid: %s", childInst, err, params.ReqID)
			return err
		}
		parentID, err := childInst.GetParentID()
		if err != nil {
			blog.Errorf("[fillMainlineChildInst] GetParentID for %+v failed: %v, rid: %s", childInst, err, params.ReqID)
			return err
		}

		object := childInst.GetObject().Object()
		tmp := &metadata.TopoInstRst{Child: []*metadata.TopoInstRst{}}
		tmp.InstID = childInstID
		tmp.InstName = childInstName
		tmp.ObjID = object.ObjectID
		tmp.ObjName = object.ObjectName

		childInstMap[parentID] = append(childInstMap[parentID], tmp)
		childTopoInsts = append(childTopoInsts, tmp)
	}

	for _, parentInst := range parentInsts {
		parentInst.Child = append(parentInst.Child, childInstMap[parentInst.InstID]...)
	}

	return assoc.fillMainlineChildInst(params, childObj, childTopoInsts)
}
