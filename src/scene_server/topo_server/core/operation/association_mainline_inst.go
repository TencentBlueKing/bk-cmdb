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
	"io"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (cli *association) ResetMainlineInstAssociatoin(params types.LogicParams, current model.Object) error {

	defaultCond := &metadata.QueryInput{}

	// fetch all parent inst
	_, currentInsts, err := cli.inst.FindInst(params, current, defaultCond)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find current object(%s) inst, error info is %s", current.GetID(), err.Error())
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
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline parent inst, error info is %s", current.GetID(), err.Error())
			return err
		}

		// reset the child's parent
		child, err := currentInst.GetMainlineChildInst()
		if nil != err {
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline child inst, error info is %s", current.GetID(), err.Error())
			return err
		}
		blog.Infof("the child: %s", child.GetObject().GetID())

		// set the child's parent
		if err = child.SetMainlineParentInst(parent); nil != err {
			blog.Errorf("[operation-asst] failed to set the object(%s) mainline child inst, error info is %s", child.GetObject().GetID(), err.Error())
			return err
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

func (cli *association) SetMainlineInstAssociation(params types.LogicParams, parent, current, child model.Object) error {

	defaultCond := &metadata.QueryInput{}

	// fetch all parent inst
	_, parentInsts, err := cli.inst.FindInst(params, parent, defaultCond)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find parent object(%s) inst, error info is %s", parent.GetID(), err.Error())
		return err
	}

	// reset the parent's inst
	for _, parent := range parentInsts {

		// create the default inst
		defaultInst := cli.instFactory.CreateInst(params, current)
		defaultInst.SetValue(common.BKOwnerIDField, params.Header.OwnerID)
		defaultInst.SetValue(current.GetInstNameFieldName(), current.GetName())
		defaultInst.SetValue(common.BKDefaultField, 0)

		// create the inst
		if err = defaultInst.Create(); nil != err {
			blog.Errorf("[operation-asst] failed to create object(%s) default inst, error info is %s", current.GetID(), err.Error())
			return err
		}

		// reset the child's parent
		child, err := parent.GetMainlineChildInst()
		if nil != err {
			if io.EOF == err {
				continue
			}
			blog.Errorf("[operation-asst] failed to get the object(%s) mainline child inst, error info is %s", parent.GetObject().GetID(), err.Error())
			return err
		}
		blog.Infof("the child: %s", child.GetObject().GetID())

		// set the child's parent
		if err = child.SetMainlineParentInst(defaultInst); nil != err {
			blog.Errorf("[operation-asst] failed to set the object(%s) mainline child inst, error info is %s", child.GetObject().GetID(), err.Error())
			return err
		}

	}

	return nil
}

func (cli *association) SearchMainlineAssociationInstTopo(params types.LogicParams, bizID int64) ([]*metadata.TopoInstRst, error) {

	bizObj, err := cli.obj.FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		return nil, err
	}

	cond := &metadata.QueryInput{}
	cond.Condition = frtypes.MapStr{
		bizObj.GetInstIDFieldName(): bizID,
	}
	_, bizInsts, err := cli.inst.FindInst(params, bizObj, cond)

	if nil != err {
		return nil, err
	}

	results := make([]*metadata.TopoInstRst, 0)

	for _, biz := range bizInsts {

		var lastRst *metadata.TopoInstRst
	exist_for:
		for {
			instID, err := biz.GetInstID()
			if nil != err {
				return nil, err
			}
			instName, err := biz.GetInstName()
			if nil != err {
				return nil, err
			}

			tmp := metadata.TopoInstRst{}
			tmp.InstID = instID
			tmp.InstName = instName
			tmp.ObjID = biz.GetObject().GetID()
			tmp.ObjName = biz.GetObject().GetName()

			tmpChild, err := biz.GetMainlineChildInst()
			if nil != err {
				if io.EOF == err {
					break exist_for
				}
				return nil, err
			}
			biz = tmpChild
			if nil != lastRst {
				lastRst.Child = append(lastRst.Child, tmp)
			} else {
				lastRst = &tmp
			}

		}

		results = append(results, lastRst)

	}
	return results, nil
}
