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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (cli *association) canReset(params types.ContextParams, currentInsts []inst.Inst) error {

	instNames := map[string]struct{}{}
	for _, currInst := range currentInsts {

		currInstParentID, err := currInst.GetParentID()
		if nil != err {
			return err
		}

		// reset the child's parent
		childs, err := currInst.GetMainlineChildInst()
		if nil != err {
			return err
		}

		for _, child := range childs {
			instName, err := child.GetInstName()
			if nil != err {
				return err
			}
			key := fmt.Sprintf("%d_%s", currInstParentID, instName)
			if _, ok := instNames[key]; ok {
				errMsg := params.Err.Error(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat).Error() + " " + instName
				return params.Err.New(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat, errMsg)
			}

			instNames[key] = struct{}{}
		}
	}

	return nil
}

func (cli *association) ResetMainlineInstAssociatoin(params types.ContextParams, current model.Object) error {

	defaultCond := &metadata.QueryInput{}
	cond := condition.CreateCondition()
	if current.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(current.GetID())
	}
	defaultCond.Condition = cond.ToMapStr()

	// fetch all parent inst
	_, currentInsts, err := cli.inst.FindInst(params, current, defaultCond, false)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find current object(%s) inst, error info is %s", current.GetID(), err.Error())
		return err
	}

	if err := cli.canReset(params, currentInsts); nil != err {
		blog.Errorf("[operation-asst] can not be reset, error info is %s", err.Error())
		return err
	}

	// NEED FIX: 下面循环中的continue ，会在处理实例异常的时候跳过当前拓扑的处理，此方式可能会导致某个业务拓扑失败，但是不会影响所有。
	// reset the parent's inst
	for _, currentInst := range currentInsts {
		// delete the current inst
		instID, err := currentInst.GetInstID()
		if nil != err {
			blog.Errorf("[operation-asst] failed to get the inst id from the inst(%#v)", currentInst.ToMapStr())
			continue
		}

		parent, err := currentInst.GetMainlineParentInst()
		if nil != err {
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline parent inst, the current inst(%v), error info is %s", current.GetID(), currentInst.GetValues(), err.Error())
			continue
		}

		// reset the child's parent
		childs, err := currentInst.GetMainlineChildInst()
		if nil != err {
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline child inst, error info is %s", current.GetID(), err.Error())
			continue
		}
		for _, child := range childs {

			// set the child's parent
			if err = child.SetMainlineParentInst(parent); nil != err {
				blog.Errorf("[operation-asst] failed to set the object(%s) mainline child inst, error info is %s", child.GetObject().GetID(), err.Error())
				continue
			}
		}

		// delete the current inst
		cond := condition.CreateCondition()
		cond.Field(currentInst.GetObject().GetInstIDFieldName()).Eq(instID)
		if err = cli.inst.DeleteInst(params, current, cond, false); nil != err {
			blog.Errorf("[operation-asst] failed to delete the current inst(%#v), error info is %s", currentInst.ToMapStr(), err.Error())
			continue
		}
	}

	return nil
}

func (cli *association) SetMainlineInstAssociation(params types.ContextParams, parent, current, child model.Object) error {

	defaultCond := &metadata.QueryInput{}
	cond := condition.CreateCondition()
	if parent.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(parent.GetID())
	}
	defaultCond.Condition = cond.ToMapStr()
	// fetch all parent inst
	_, parentInsts, err := cli.inst.FindInst(params, parent, defaultCond, false)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find parent object(%s) inst, error info is %s", parent.GetID(), err.Error())
		return err
	}

	// reset the parent's inst
	for _, parent := range parentInsts {

		id, err := parent.GetInstID()
		if nil != err {
			blog.Errorf("[operation-asst] failed to find the inst id, error info is %s", err.Error())
			return err
		}

		// create the default inst
		defaultInst := cli.instFactory.CreateInst(params, current)
		defaultInst.SetValue(common.BKOwnerIDField, params.SupplierAccount)
		defaultInst.SetValue(current.GetInstNameFieldName(), current.GetName())
		defaultInst.SetValue(common.BKDefaultField, 0)
		defaultInst.SetValue(common.BKInstParentStr, id)

		// create the inst
		if err = defaultInst.Create(); nil != err {
			blog.Errorf("[operation-asst] failed to create object(%s) default inst, error info is %s", current.GetID(), err.Error())
			return err
		}

		// reset the child's parent
		childs, err := parent.GetMainlineChildInst()
		if nil != err {
			if io.EOF == err {
				continue
			}
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline child inst, error info is %s", parent.GetObject().GetID(), err.Error())
			return err
		}
		for _, child := range childs {

			// set the child's parent
			if err = child.SetMainlineParentInst(defaultInst); nil != err {
				blog.Errorf("[operation-asst] failed to set the object(%s) mainline child inst, error info is %s", child.GetObject().GetID(), err.Error())
				return err
			}
		}

	}

	return nil
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
		cond.Field(common.BKObjIDField).Eq(childObj.GetID())
	} else if childObj.GetID() == common.BKInnerObjIDSet {
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

		tmp := &metadata.TopoInstRst{Child: []*metadata.TopoInstRst{}}
		tmp.InstID = childInstID
		tmp.InstName = childInstName
		tmp.ObjID = childInst.GetObject().GetID()
		tmp.ObjName = childInst.GetObject().GetName()

		childInstMap[parentID] = append(childInstMap[parentID], tmp)
		childTopoInsts = append(childTopoInsts, tmp)
	}

	for _, parentInst := range parentInsts {
		parentInst.Child = append(parentInst.Child, childInstMap[parentInst.InstID]...)
	}

	return cli.fillMainlineChildInst(params, childObj, childTopoInsts)
}

func (cli *association) SearchMainlineAssociationInstTopo(params types.ContextParams, obj model.Object, instID int64) ([]*metadata.TopoInstRst, error) {

	cond := &metadata.QueryInput{}
	cond.Condition = frtypes.MapStr{
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

		tmp := &metadata.TopoInstRst{Child: []*metadata.TopoInstRst{}}
		tmp.InstID = instID
		tmp.InstName = instName
		tmp.ObjID = biz.GetObject().GetID()
		tmp.ObjName = biz.GetObject().GetName()

		results = append(results, tmp)
	}

	if err = cli.fillMainlineChildInst(params, obj, results); err != nil {
		blog.Errorf("[SearchMainlineAssociationInstTopo] fillMainlineChildInst for %+v failed: %v", results, err)
		return nil, err
	}
	return results, nil
}
