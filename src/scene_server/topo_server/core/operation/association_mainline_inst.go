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

	// reset the parent's inst
	for _, currentInst := range currentInsts {
		// delete the current inst
		instID, err := currentInst.GetInstID()
		if nil != err {
			blog.Errorf("[operation-asst] failed to get the inst id from the inst(%#v)", currentInst.ToMapStr())
			return err
		}

		parent, err := currentInst.GetMainlineParentInst()
		if nil != err {
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline parent inst, the current inst(%v), error info is %s", current.GetID(), currentInst.GetValues(), err.Error())
			return err
		}

		// reset the child's parent
		childs, err := currentInst.GetMainlineChildInst()
		if nil != err {
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline child inst, error info is %s", current.GetID(), err.Error())
			return err
		}
		for _, child := range childs {
			blog.Infof("the child: %s", child.GetObject().GetID())

			// set the child's parent
			if err = child.SetMainlineParentInst(parent); nil != err {
				blog.Errorf("[operation-asst] failed to set the object(%s) mainline child inst, error info is %s", child.GetObject().GetID(), err.Error())
				return err
			}
		}

		// delete the current inst
		cond := condition.CreateCondition()
		cond.Field(currentInst.GetObject().GetInstIDFieldName()).Eq(instID)
		if err = cli.inst.DeleteInst(params, current, cond); nil != err {
			blog.Errorf("[operation-asst] failed to delete the current inst(%#v), error info is %s", currentInst.ToMapStr(), err.Error())
			return err
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
func (cli *association) constructTopo(params types.ContextParams, targetInst inst.Inst) ([]metadata.TopoInstRst, error) {

	childs, err := targetInst.GetMainlineChildInst()
	if nil != err {
		return nil, err
	}

	results := []metadata.TopoInstRst{}

	if 0 == len(childs) {
		return []metadata.TopoInstRst{}, nil
	}

	for _, child := range childs {

		rst := metadata.TopoInstRst{
			Child: []metadata.TopoInstRst{},
		}

		id, err := child.GetInstID()
		if nil != err {
			return nil, err
		}

		name, err := child.GetInstName()
		if nil != err {
			return nil, err
		}

		rst.InstID = id
		rst.InstName = name
		rst.ObjID = child.GetObject().GetID()
		rst.ObjName = child.GetObject().GetName()

		childRst, err := cli.constructTopo(params, child)
		if nil != err {
			return nil, err
		}

		rst.Child = append(rst.Child, childRst...)
		results = append(results, rst)
	}

	return results, nil
}
func (cli *association) SearchMainlineAssociationInstTopo(params types.ContextParams, obj model.Object, instID int64) ([]*metadata.TopoInstRst, error) {

	cond := &metadata.QueryInput{}
	cond.Condition = frtypes.MapStr{
		obj.GetInstIDFieldName(): instID,
	}

	_, bizInsts, err := cli.inst.FindInst(params, obj, cond, false)
	if nil != err {
		return nil, err
	}

	results := make([]*metadata.TopoInstRst, 0)

	for _, biz := range bizInsts {

		instID, err := biz.GetInstID()
		if nil != err {
			return nil, err
		}
		instName, err := biz.GetInstName()
		if nil != err {
			return nil, err
		}

		tmp := &metadata.TopoInstRst{Child: []metadata.TopoInstRst{}}
		tmp.InstID = instID
		tmp.InstName = instName
		tmp.ObjID = biz.GetObject().GetID()
		tmp.ObjName = biz.GetObject().GetName()

		rst, err := cli.constructTopo(params, biz)
		if nil != err {
			return nil, err
		}

		tmp.Child = append(tmp.Child, rst...)
		results = append(results, tmp)

	}
	return results, nil
}
