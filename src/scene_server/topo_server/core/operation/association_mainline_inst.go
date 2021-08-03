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
	"regexp"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// checkInstNameRepeat 检查如果将 currentInsts 都删除之后，拥有共同父节点的孩子结点会不会出现名字冲突
// 如果有冲突，返回 (false, 冲突实例名, nil)
func (assoc *association) checkInstNameRepeat(kit *rest.Kit, currentInsts []inst.Inst) (canReset bool, repeatedInstName string, err error) {
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

func (assoc *association) ResetMainlineInstAssociation(kit *rest.Kit, current model.Object) error {
	rid := kit.Rid

	cObj := current.Object()
	cond := condition.CreateCondition()
	if current.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(cObj.ObjectID)
	}
	defaultCond := &metadata.QueryInput{}
	defaultCond.Condition = cond.ToMapStr()

	// 获取 current 模型的所有实例
	_, currentInsts, err := assoc.inst.FindInst(kit, current, defaultCond, false)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find current object(%s) inst, err: %+v, rid: %s", cObj.ObjectID, err, rid)
		return err
	}

	// 检查实例删除后，会不会出现重名冲突
	var canReset bool
	var repeatedInstName string
	if canReset, repeatedInstName, err = assoc.checkInstNameRepeat(kit, currentInsts); nil != err {
		blog.Errorf("[operation-asst] can not be reset, err: %+v, rid: %s", err, rid)
		return err
	}
	if canReset == false {
		blog.Errorf("[operation-asst] can not be reset, inst name repeated, inst: %s, rid: %s", repeatedInstName, rid)
		errMsg := kit.CCError.Error(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat).Error() + " " + repeatedInstName
		return kit.CCError.New(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat, errMsg)
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
		if err := assoc.inst.DeleteMainlineInstWithID(kit, current, instID); nil != err {
			blog.Errorf("[operation-asst] failed to delete the current inst(%#v), err: %+v, rid: %s", currentInst.ToMapStr(), err, rid)
			continue
		}
	}

	return nil
}

func (assoc *association) SetMainlineInstAssociation(kit *rest.Kit, parent, current, child model.Object) ([]int64, error) {
	defaultCond := &metadata.QueryInput{}
	cond := condition.CreateCondition()
	if parent.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(parent.Object().ObjectID)
	}
	defaultCond.Condition = cond.ToMapStr()
	// fetch all parent instances.
	_, parentInsts, err := assoc.inst.FindInst(kit, parent, defaultCond, false)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find parent object(%s) inst, err: %s, rid: %s", parent.Object().ObjectID, err.Error(), kit.Rid)
		return nil, err
	}

	createdInstIDs := make([]int64, len(parentInsts))
	expectParent2Children := map[int64][]inst.Inst{}
	// filters out special character for mainline instances
	re, _ := regexp.Compile(`[#/,><|]`)
	instanceName := re.ReplaceAllString(current.Object().ObjectName, "")
	// create current object instance for each parent instance and insert the current instance to
	for _, parent := range parentInsts {
		id, err := parent.GetInstID()
		if nil != err {
			blog.Errorf("[operation-asst] failed to find the inst id, err: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}

		// we create the current object's instance for each parent instance belongs to the parent object.
		currentInst := assoc.instFactory.CreateInst(kit, current)
		currentInst.SetValue(current.GetInstNameFieldName(), instanceName)
		// set current instance's parent id to parent instance's id, so that they can be chained.
		currentInst.SetValue(common.BKInstParentStr, id)
		object := parent.GetObject()
		if object.GetObjectID() == common.BKInnerObjIDApp {
			currentInst.SetValue(common.BKAppIDField, id)
		} else {
			if bizID, ok := parent.GetValues().Get(common.BKAppIDField); ok {
				currentInst.SetValue(common.BKAppIDField, bizID)
			}
		}

		// create the instance now.
		if err = currentInst.Create(); nil != err {
			blog.Errorf("[operation-asst] failed to create object(%s) default inst, err: %s, rid: %s", current.Object().ObjectID, err.Error(), kit.Rid)
			return nil, err
		}
		instID, err := currentInst.GetInstID()
		if err != nil {
			blog.Errorf("create mainline instance for obj: %s, but got invalid instance id, err :%v, rid: %s", current.Object().ObjectID, err, kit.Rid)
			return nil, err
		}
		createdInstIDs = append(createdInstIDs, instID)

		// reset the child's parent instance's parent id to current instance's id.
		children, err := parent.GetMainlineChildInst()
		if nil != err {
			if io.EOF == err {
				continue
			}
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline child inst, err: %s, rid: %s", parent.GetObject().Object().ObjectID, err.Error(), kit.Rid)
			return nil, err
		}

		expectParent2Children[instID] = children
	}

	for parentID, children := range expectParent2Children {
		for _, child := range children {
			// set the child's parent
			if err = child.SetMainlineParentInst(parentID); nil != err {
				blog.Errorf("[operation-asst] failed to set the object(%s) mainline child inst, err: %s, rid: %s", child.GetObject().Object().ObjectID, err.Error(), kit.Rid)
				return nil, err
			}
		}
	}
	return createdInstIDs, nil
}

func (assoc *association) SearchMainlineAssociationInstTopo(kit *rest.Kit, objID string, instID int64,
	withStatistics bool, withDefault bool) ([]*metadata.TopoInstRst, errors.CCError) {
	// read mainline object association and construct child relation map excluding host
	mainlineAsstRsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline}})
	if nil != err {
		blog.Errorf("search mainline association failed, error: %s, rid: %s", err.Error(), kit.Rid)
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
	objects, err := assoc.obj.FindObject(kit, condition.CreateCondition().Field(common.BKObjIDField).In(objectIDs))
	if nil != err {
		blog.ErrorJSON("search mainline objects(%s) failed, error: %s, rid: %s", objectIDs, err.Error(), kit.Rid)
		return nil, err
	}
	for _, object := range objects {
		objectNameMap[object.GetObjectID()] = object.Object().ObjectName
	}

	// traverse and fill instance topology data
	results := make([]*metadata.TopoInstRst, 0)
	var parents []*metadata.TopoInstRst
	instCond := map[string]interface{}{
		metadata.GetInstIDFieldByObjID(objID): instID,
	}
	var bizID int64
	moduleIDs := make([]int64, 0)
	for objectID := objID; len(objectID) != 0; objectID = mainlineObjectChildMap[objectID] {
		filter := &metadata.QueryInput{Condition: instCond}
		if objectID != objID {
			filter.Sort = common.BKDefaultField
		}
		instanceRsp, err := assoc.inst.FindOriginInst(kit, objectID, filter)
		if err != nil {
			blog.ErrorJSON("search inst failed, err: %s, cond:%s, rid: %s", err, instCond, kit.Rid)
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
			instID, err := instance.Int64(metadata.GetInstIDFieldByObjID(objectID))
			if nil != err {
				blog.ErrorJSON("get instance %s id failed, err: %s, rid: %s", instance, err, kit.Rid)
				return nil, err
			}
			instIDs = append(instIDs, instID)
			instName, err := instance.String(metadata.GetInstNameFieldName(objectID))
			if nil != err {
				blog.ErrorJSON("get instance %s name failed, err: %s, rid: %s", instance, err, kit.Rid)
				return nil, err
			}
			defaultValue := 0
			if defaultFieldValue, exist := instance[common.BKDefaultField]; exist {
				defaultValue, err = util.GetIntByInterface(defaultFieldValue)
				if err != nil {
					blog.ErrorJSON("get instance %s default failed, err: %s, rid: %s", instance, err, kit.Rid)
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
						blog.ErrorJSON("get instance %s biz id failed, err: %s, rid: %s", instance, err, kit.Rid)
						return nil, err
					}
				}
			}
			if objectID == objID {
				results = append(results, topoInst)
			} else {
				parentID, err := instance.Int64(common.BKParentIDField)
				if err != nil {
					blog.ErrorJSON("get instance %s parent id failed, err: %s, rid: %s", instance, err, kit.Rid)
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
			blog.Errorf("[SearchMainlineAssociationInstTopo] fill statistics data failed, bizID: %d, err: %v, rid: %s", bizID, err, kit.Rid)
			return nil, err
		}
	}
	return results, nil
}

// TopoNodeHostAndSerInstCount get topo node host and service instance count
func (assoc *association) TopoNodeHostAndSerInstCount(kit *rest.Kit, instID int64,
	input *metadata.HostAndSerInstCountOption) ([]*metadata.TopoNodeHostAndSerInstCount, errors.CCError) {
	var bizID int64
	setIDs := make([]int64, 0)
	moduleIDs := make([]int64, 0)
	customLevels := make(map[int64]string, 0)
	for _, obj := range input.Condition {
		switch obj.ObjID {
		case common.BKInnerObjIDSet:
			setIDs = append(setIDs, obj.InstID)
		case common.BKInnerObjIDModule:
			moduleIDs = append(moduleIDs, obj.InstID)
		case common.BKInnerObjIDApp:
			bizID = obj.InstID
		default:
			customLevels[obj.InstID] = obj.ObjID
		}
	}

	results := make([]*metadata.TopoNodeHostAndSerInstCount, 0)
	// handle module host number and service instance count
	if len(moduleIDs) > 0 {
		res, err := assoc.getHostSvcInstCountByModuleIDs(kit, moduleIDs)
		if err != nil {
			blog.Errorf("get module host and service instance count failed, err: %v, moduleIDs: %s, rid: %s", err,
				moduleIDs, kit.Rid)
			return nil, err
		}
		results = append(results, res...)
	}

	// handle set host number and service instance count
	if len(setIDs) > 0 {
		res, err := assoc.getHostSvcInstCountBySetIDs(kit, setIDs)
		if err != nil {
			blog.Errorf("get set host and service instance count failed, err: %v, setIDs: %s, rid: %s", err,
				setIDs, kit.Rid)
			return nil, err
		}
		results = append(results, res...)
	}

	// handle biz host and service instance count
	if bizID > 0 {
		res, err := assoc.getHostSvcInstCountByBizID(kit, bizID)
		if err != nil {
			blog.Errorf("get biz host and service instance count failed, err: %v, bizID: %s, rid: %s", err, bizID,
				kit.Rid)
			return nil, err
		}
		results = append(results, res)
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

// getHostSvcInstCountByModuleIDs get host and service instace count by module ids
func (assoc *association) getHostSvcInstCountByModuleIDs(kit *rest.Kit,
	moduleIDs []int64) ([]*metadata.TopoNodeHostAndSerInstCount, error) {
	moduleArr := make([][]int64, 0)
	for _, moduleID := range moduleIDs {
		moduleArr = append(moduleArr, []int64{moduleID})
	}

	svcInstCounts, e := assoc.getServiceInstCount(kit, common.BKModuleIDField, moduleArr)
	if e != nil {
		blog.Errorf("get service instance count failed, err: %v, objID: %s, instIDs: %s, rid: %s", e,
			common.BKModuleIDField, moduleIDs, kit.Rid)
		return nil, e
	}

	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 5)
	results := make([]*metadata.TopoNodeHostAndSerInstCount, len(moduleIDs))
	for idx, moduleID := range moduleIDs {
		pipeline <- true
		wg.Add(1)
		go func(idx int, moduleID int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()
			moduleArr := []int64{moduleID}
			hostCount, err := assoc.getDistinctHostCount(kit, common.BKModuleIDField, moduleArr)
			if err != nil {
				blog.Errorf("get distinct host count failed, err: %v, objID: %s, instID: %s, rid: %s", err,
					common.BKModuleIDField, moduleID, kit.Rid)
				firstErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)

				return
			}

			topoNodeHostCount := &metadata.TopoNodeHostAndSerInstCount{
				ObjID:     common.BKInnerObjIDModule,
				InstID:    moduleID,
				HostCount: hostCount,
			}

			lock.Lock()
			results[idx] = topoNodeHostCount
			lock.Unlock()
		}(idx, moduleID)
	}
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	for idx, count := range svcInstCounts {
		results[idx].ServiceInstanceCount = count
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
				blog.Errorf("find hosts by topo failed, get set ID by topo err: %v, objID: %s, instID: %d, "+
					"rid:%s", err, objID, instID, kit.Rid)
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

// getHostSvcInstCountByBizID get host and service instace count by biz id
func (assoc *association) getHostSvcInstCountByBizID(kit *rest.Kit,
	bizID int64) (*metadata.TopoNodeHostAndSerInstCount, error) {

	bizArr := make([][]int64, 0)
	bizArr = append(bizArr, []int64{bizID})
	hostCount, err := assoc.getDistinctHostCount(kit, common.BKAppIDField, []int64{bizID})
	if err != nil {
		blog.Errorf("get distinct host count failed, err: %v, objID: %s, instIDs: %, rid: %s", err,
			common.BKAppIDField, bizArr, kit.Rid)
		return nil, err
	}

	svcInstCount, e := assoc.getServiceInstCount(kit, common.BKAppIDField, bizArr)
	if e != nil {
		blog.Errorf("get service instance count failed, err: %v, objID: %s, instID: %s, rid: %s", e,
			common.BKAppIDField, bizID, kit.Rid)
		return nil, e
	}

	bizHostSvcInstCount := &metadata.TopoNodeHostAndSerInstCount{
		ObjID:                common.BKInnerObjIDApp,
		InstID:               bizID,
		HostCount:            hostCount,
		ServiceInstanceCount: svcInstCount[0],
	}

	return bizHostSvcInstCount, nil
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

	if !asstRes.Result {
		blog.Errorf("get set IDs by topo failed, get mainline association err: %s, rid: %s", asstRes.ErrMsg,
			kit.Rid)
		return nil, asstRes.CCError()
	}

	childObjMap := make(map[string]string)
	for _, asst := range asstRes.Data.Info {
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
		query := &metadata.QueryCondition{
			Condition: map[string]interface{}{common.BKParentIDField: map[string]interface{}{common.BKDBIN: instIDs}},
			Fields:    []string{idField},
			Page:      metadata.BasePage{Limit: common.BKNoLimit},
		}

		instRes, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, childObj, query)
		if err != nil {
			blog.Errorf("get set IDs by topo failed, read instance err: %s, objID: %s, instIDs: %+v, rid: %s",
				err.Error(), childObj, instIDs, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !instRes.Result {
			blog.Errorf("get set IDs by topo failed, read instance err: %s, objID: %s, instIDs: %+v, rid: %s",
				instRes.ErrMsg, childObj, instIDs, kit.Rid)
			return nil, kit.CCError.New(instRes.Code, instRes.ErrMsg)
		}

		if len(instRes.Data.Info) == 0 {
			return []int64{}, nil
		}

		instIDs = make([]int64, len(instRes.Data.Info))
		for index, insts := range instRes.Data.Info {
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

	setRelModuleMap := make(map[int64][]int64, 0)
	for _, mapStr := range resp.Data.Info {
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

func (assoc *association) fillStatistics(kit *rest.Kit, bizID int64, moduleIDs []int64, topoInsts []*metadata.TopoInstRst) errors.CCError {
	// get service instance count
	option := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	serviceInstances, err := assoc.clientSet.CoreService().Process().ListServiceInstance(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("fillStatistics failed, list service instances failed, option: %+v, err: %s, rid: %s", option, err.Error(), kit.Rid)
		return err
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
	hostModules, e := assoc.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, listHostOption)
	if e != nil {
		blog.Errorf("fillStatistics failed, list host modules failed, option: %+v, err: %s, rid: %s", listHostOption, e.Error(), kit.Rid)
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

	// get host apply rule count
	listApplyRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: moduleIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostApplyRules, err := assoc.clientSet.CoreService().HostApplyRule().ListHostApplyRule(kit.Ctx, kit.Header, bizID, listApplyRuleOption)
	if err != nil {
		blog.ErrorJSON("fillStatistics failed, ListHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, listApplyRuleOption, err, kit.Rid)
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
