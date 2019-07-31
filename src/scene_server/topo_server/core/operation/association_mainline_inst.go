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

	// "strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// checkInstNameRepeat 检查如果将 currentInsts 都删除之后，拥有共同父节点的孩子结点会不会出现名字冲突
// 如果有冲突，返回 (false, 冲突实例名, nil)
func (cli *association) checkInstNameRepeat(params types.ContextParams, currentInsts []inst.Inst) (canReset bool, repeatedInstName string, err error) {
	// TODO: 返回值将bool值与出错情况分开 (bool, error)
	instNames := map[string]bool{}
	for _, currInst := range currentInsts {
		currInstParentID, err := currInst.GetParentID()
		if nil != err {
			return false, "", err
		}

		childs, err := currInst.GetMainlineChildInst()
		if nil != err {
			return false, "", err
		}

		for _, child := range childs {
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

func (cli *association) ResetMainlineInstAssociatoin(params types.ContextParams, current model.Object) error {
	rid := params.ReqID

	cObj := current.Object()
	cond := condition.CreateCondition()
	if current.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(cObj.ObjectID)
	}
	defaultCond := &metadata.QueryInput{}
	defaultCond.Condition = cond.ToMapStr()

	// 获取 current 模型的所有实例
	_, currentInsts, err := cli.inst.FindInst(params, current, defaultCond, false)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find current object(%s) inst, err: %+v, rid: %s", cObj.ObjectID, err, rid)
		return err
	}

	// 检查实例删除后，会不会出现重名冲突
	var canReset bool
	var repeatedInstName string
	if canReset, repeatedInstName, err = cli.checkInstNameRepeat(params, currentInsts); nil != err {
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
		if err := cli.inst.DeleteMainlineInstWithID(params, current, instID); nil != err {
			blog.Errorf("[operation-asst] failed to delete the current inst(%#v), err: %+v, rid: %s", currentInst.ToMapStr(), err, rid)
			continue
		}
	}

	return nil
}

func (cli *association) SetMainlineInstAssociation(params types.ContextParams, parent, current, child model.Object) error {

	defaultCond := &metadata.QueryInput{}
	cond := condition.CreateCondition()
	if parent.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(parent.Object().ObjectID)
	}
	defaultCond.Condition = cond.ToMapStr()
	// fetch all parent instances.
	_, parentInsts, err := cli.inst.FindInst(params, parent, defaultCond, false)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find parent object(%s) inst, err: %s", parent.Object().ObjectID, err.Error())
		return err
	}

	expectParent2Childs := map[int64][]inst.Inst{}
	// create current object instance for each parent instance and insert the current instance to
	for _, parent := range parentInsts {

		id, err := parent.GetInstID()
		if nil != err {
			blog.Errorf("[operation-asst] failed to find the inst id, err: %s", err.Error())
			return err
		}

		// we create the current object's instance for each parent instance belongs to the parent object.
		currentInst := cli.instFactory.CreateInst(params, current)
		currentInst.SetValue(current.GetInstNameFieldName(), current.Object().ObjectName)
		currentInst.SetValue(common.BKDefaultField, 0)
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
			blog.Errorf("[operation-asst] failed to create object(%s) default inst, err: %s", current.Object().ObjectID, err.Error())
			return err
		}
		instID, err := currentInst.GetInstID()
		if err != nil {
			blog.Errorf("create mainline instance for obj: %s, but got invalid instance id, err :%v", current.Object().ObjectID, err)
			return err
		}
		err = cli.authManager.RegisterInstancesByID(params.Context, params.Header, current.Object().ObjectID, instID)
		if err != nil {
			blog.Errorf("create mainline instance for object: %s, but register to auth center failed, err: %v", current.Object().ObjectID, err)
			return err
		}

		// reset the child's parent instance's parent id to current instance's id.
		childs, err := parent.GetMainlineChildInst()
		if nil != err {
			if io.EOF == err {
				continue
			}
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline child inst, err: %s", parent.GetObject().Object().ObjectID, err.Error())
			return err
		}

		curInstID, err := currentInst.GetInstID()
		if err != nil {
			blog.Errorf("[operation-asst] failed to get the instID(%#v), err: %s", currentInst.ToMapStr(), err.Error())
			return err
		}

		expectParent2Childs[curInstID] = childs
	}

	for parentID, childs := range expectParent2Childs {
		for _, child := range childs {
			// set the child's parent
			if err = child.SetMainlineParentInst(parentID); nil != err {
				blog.Errorf("[operation-asst] failed to set the object(%s) mainline child inst, err: %s", child.GetObject().Object().ObjectID, err.Error())
				return err
			}
		}
	}

	return nil
}

func (cli *association) SearchMainlineAssociationInstTopo(params types.ContextParams, obj model.Object, instID int64) ([]*metadata.TopoInstRst, error) {

	cond := &metadata.QueryInput{}
	cond.Condition = mapstr.MapStr{
		obj.GetInstIDFieldName(): instID,
	}

	_, bizInsts, err := cli.inst.FindInst(params, obj, cond, false)
	if nil != err {
		blog.Errorf("[SearchMainlineAssociationInstTopo] FindInst for %+v failed: %v", cond, err)
		return nil, err
	}

	results := make([]*metadata.TopoInstRst, 0)
	for _, biz := range bizInsts {
		instID, err := biz.GetInstID()
		if nil != err {
			blog.Errorf("[SearchMainlineAssociationInstTopo] GetInstID for %+v failed: %v", biz, err)
			return nil, err
		}
		instName, err := biz.GetInstName()
		if nil != err {
			blog.Errorf("[SearchMainlineAssociationInstTopo] GetInstName for %+v failed: %v", biz, err)
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

	if err = cli.fillMainlineChildInst(params, obj, results); err != nil {
		blog.Errorf("[SearchMainlineAssociationInstTopo] fillMainlineChildInst for %+v failed: %v", results, err)
		return nil, err
	}
	return results, nil
}

func (cli *association) fillMainlineChildInst(params types.ContextParams, object model.Object, parentInsts []*metadata.TopoInstRst) error {
	childObj, err := object.GetMainlineChildObject()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		blog.Errorf("[fillMainlineChildInst] GetMainlineChildObject for %+v failed: %v", object, err)
		return err
	}

	parentIDs := []int64{}
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

	_, childInsts, err := cli.inst.FindInst(params, childObj, &metadata.QueryInput{Condition: cond.ToMapStr()}, false)
	if err != nil {
		blog.Errorf("[fillMainlineChildInst] FindInst for %+v failed: %v", cond.ToMapStr(), err)
		return err
	}

	// parentID mapping to child topo insts
	childInstMap := map[int64][]*metadata.TopoInstRst{}
	childTopoInsts := []*metadata.TopoInstRst{}
	for _, childInst := range childInsts {
		childInstID, err := childInst.GetInstID()
		if err != nil {
			blog.Errorf("[fillMainlineChildInst] GetInstID for %+v failed: %v", childInst, err)
			return err
		}
		childInstName, err := childInst.GetInstName()
		if nil != err {
			blog.Errorf("[fillMainlineChildInst] GetInstName for %+v failed: %v", childInst, err)
			return err
		}
		parentID, err := childInst.GetParentID()
		if err != nil {
			blog.Errorf("[fillMainlineChildInst] GetParentID for %+v failed: %v", childInst, err)
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

	return cli.fillMainlineChildInst(params, childObj, childTopoInsts)
}
