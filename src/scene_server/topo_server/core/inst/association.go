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
	"context"
	"io"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (cli *inst) updateMainlineAssociation(child Inst, parentID int64) error {

	childID, err := child.GetInstID()
	if nil != err {
		return err
	}

	object := child.GetObject().Object()

	cond := condition.CreateCondition()
	cond.Field(object.GetInstIDFieldName()).Eq(int(childID))
	if object.IsCommon() {
		cond.Field(metadata.ModelFieldObjectID).Eq(object.ObjectID)
	}

	input := metadata.UpdateOption{
		Data: mapstr.MapStr{
			common.BKInstParentStr: parentID,
		},
		Condition: cond.ToMapStr(),
	}
	rsp, err := cli.clientSet.CoreService().Instance().UpdateInstance(context.Background(), cli.kit.Header, object.ObjectID, &input)
	if nil != err {
		blog.Errorf("[inst-inst] failed to request object controller, error info %s, rid: %s", err.Error(), cli.kit.Rid)
		return cli.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[inst-inst] failed to update the association, err: %s, rid: %s", rsp.ErrMsg, cli.kit.Rid)
		return cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (cli *inst) searchInstAssociation(cond condition.Condition) ([]metadata.InstAsst, error) {

	rsp, err := cli.clientSet.CoreService().Association().ReadInstAssociation(context.Background(), cli.kit.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[inst-inst] failed to request the object controller , err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return nil, cli.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[inst-inst] failed to search the inst association, err: %s, rid: %s", rsp.ErrMsg, cli.kit.Rid)
		return nil, cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil

}

func (cli *inst) deleteInstAssociation(instID, asstInstID int64, objID, asstObjID string) error {

	cond := condition.CreateCondition()

	cond.Field(common.BKInstIDField).Eq(instID)
	cond.Field(common.BKAsstInstIDField).Eq(asstInstID)
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKAsstObjIDField).Eq(asstObjID)

	rsp, err := cli.clientSet.CoreService().Association().DeleteInstAssociation(context.Background(), cli.kit.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[inst-inst] failed to request the object controller , err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return cli.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[inst-inst] failed to delete the inst association, err: %s, rid: %s", rsp.ErrMsg, cli.kit.Rid)
		return cli.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return nil

}

func (cli *inst) GetMainlineParentInst() (Inst, error) {

	parentObj, err := cli.target.GetMainlineParentObject()
	if nil != err {
		return nil, err
	}

	parentID, err := cli.GetParentID()
	if nil != err {
		blog.Errorf("[inst-inst] failed to get the inst id, err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return nil, err
	}

	cond := condition.CreateCondition()
	if parentObj.IsCommon() {
		cond.Field(metadata.ModelFieldObjectID).Eq(parentObj.Object().ObjectID)
	}
	cond.Field(parentObj.GetInstIDFieldName()).Eq(parentID)

	rspItems, err := cli.searchInsts(parentObj, cond)
	if nil != err {
		blog.Errorf("[inst-inst] failed to request the object controller , err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return nil, cli.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	for _, item := range rspItems {
		return item, nil // only one mainline parent
	}

	return nil, io.EOF
}
func (cli *inst) GetMainlineChildInst() ([]Inst, error) {

	childObj, err := cli.target.GetMainlineChildObject()
	if nil != err {
		if err == io.EOF {
			return []Inst{}, nil
		}
		blog.Errorf("[inst-inst]failed to get the object(%s)'s child object, err: %s, rid: %s", cli.target.Object().ObjectID, err.Error(), cli.kit.Rid)
		return nil, err
	}

	currInstID, err := cli.GetInstID()
	if nil != err {
		blog.Errorf("[inst-inst] failed to get the inst id, err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return nil, err
	}

	cObj := childObj.Object()
	cond := condition.CreateCondition()
	if childObj.IsCommon() {
		cond.Field(metadata.ModelFieldObjectID).Eq(cObj.ObjectID)
	} else if cObj.ObjectID == common.BKInnerObjIDSet {
		cond.Field(common.BKDefaultField).NotEq(common.DefaultResSetFlag)
	}
	cond.Field(common.BKInstParentStr).Eq(currInstID)
	return cli.searchInsts(childObj, cond)
}
func (cli *inst) GetParentObjectWithInsts() ([]*ObjectWithInsts, error) {

	result := make([]*ObjectWithInsts, 0)
	objPairs, err := cli.target.GetParentObject()
	if nil != err {
		blog.Errorf("[inst-inst] failed to get the object(%s)'s parent, err: %s, rid: %s", cli.target.Object().ObjectID, err.Error(), cli.kit.Rid)
		return result, err
	}

	currInstID, err := cli.GetInstID()
	if nil != err {
		blog.Errorf("[inst-inst] failed to get the inst id, err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return result, err
	}

	for _, objPair := range objPairs {

		rstObj := &ObjectWithInsts{Object: objPair.Object}
		cond := condition.CreateCondition()
		cond.Field(common.BKAsstInstIDField).Eq(currInstID)
		cond.Field(common.BKObjIDField).Eq(objPair.Object.Object().ObjectID)
		cond.Field(common.BKAsstObjIDField).Eq(cli.target.Object().ObjectID)
		cond.Field(common.AssociationObjAsstIDField).Eq(objPair.Association.AssociationName)

		asstItems, err := cli.searchInstAssociation(cond)
		if nil != err {
			blog.Errorf("[inst-inst] failed to search the inst association, the err: %s, rid: %s", err.Error(), cli.kit.Rid)
			return result, err
		}

		// found no noe inst association with this object and association info.
		// which means that, this object association has not been instantiated.
		if len(asstItems) == 0 {
			continue
		}

		relation := make(map[int64]int64)
		parentInstIDS := []int64{}
		for _, item := range asstItems {

			parentInstID := item.InstID
			assoID := item.ID
			relation[parentInstID] = assoID
			parentInstIDS = append(parentInstIDS, parentInstID)
		}

		innerCond := condition.CreateCondition()
		innerCond.Field(objPair.Object.GetInstIDFieldName()).In(parentInstIDS)
		if objPair.Object.IsCommon() {
			innerCond.Field(metadata.ModelFieldObjectID).Eq(objPair.Object.Object().ObjectID)
		}

		rspItems, err := cli.searchInsts(objPair.Object, innerCond)
		if nil != err {
			blog.Errorf("[inst-inst] failed to search the insts by the condition(%#v), err: %s, rid: %s", innerCond, err.Error(), cli.kit.Rid)
			return result, err
		}

		for _, item := range rspItems {
			id, err := item.GetInstID()
			if err != nil {
				blog.Errorf("[inst-inst] failed to parse the instance id , err: %s, rid: %s", err.Error(), cli.kit.Rid)
				return result, err
			}
			item.SetAssoID(relation[id])
		}

		rstObj.Insts = rspItems
		result = append(result, rstObj)

	}

	return result, nil
}

func (cli *inst) GetChildObjectWithInsts() ([]*ObjectWithInsts, error) {

	result := make([]*ObjectWithInsts, 0)

	objPairs, err := cli.target.GetChildObject()
	if nil != err {
		blog.Errorf("[inst-inst] failed to get the object(%s)'s child, err: %s, rid: %s", cli.target.Object().ObjectID, err.Error(), cli.kit.Rid)
		return result, err
	}

	currInstID, err := cli.GetInstID()
	if nil != err {
		blog.Errorf("[inst-inst] failed to get the inst id, err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return result, err
	}

	for _, objPair := range objPairs {

		rstObj := &ObjectWithInsts{Object: objPair.Object}
		cond := condition.CreateCondition()
		cond.Field(common.BKInstIDField).Eq(currInstID)
		cond.Field(common.BKObjIDField).Eq(cli.target.Object().ObjectID)
		cond.Field(common.BKAsstObjIDField).Eq(objPair.Object.Object().ObjectID)
		cond.Field(common.AssociationObjAsstIDField).Eq(objPair.Association.AssociationName)

		asstItems, err := cli.searchInstAssociation(cond)
		if nil != err {
			blog.Errorf("[inst-inst] failed to search the inst association,  the err: %s, rid: %s", err.Error(), cli.kit.Rid)
			return result, err
		}

		// found no one inst association with this object and association info.
		// which means that, this object association has not been instantiated.
		if len(asstItems) == 0 {
			continue
		}

		relations := make(map[int64]int64, 0)

		childInstIDS := make([]int64, 0)
		for _, item := range asstItems {
			childInstID := item.AsstInstID
			assoID := item.ID
			childInstIDS = append(childInstIDS, childInstID)
			relations[childInstID] = assoID
		}

		innerCond := condition.CreateCondition()
		innerCond.Field(objPair.Object.GetInstIDFieldName()).In(childInstIDS)
		if objPair.Object.IsCommon() {
			innerCond.Field(metadata.ModelFieldObjectID).Eq(objPair.Object.Object().ObjectID)
		}

		rspItems, err := cli.searchInsts(objPair.Object, innerCond)
		if nil != err {
			blog.Errorf("[inst-inst] failed to search the insts by the condition(%#v), err: %s, rid: %s", innerCond, err.Error(), cli.kit.Rid)
			return result, err
		}

		for _, item := range rspItems {
			id, err := item.GetInstID()
			if err != nil {
				blog.Errorf("[inst-inst] failed to parse the association id , err: %s, rid: %s", err.Error(), cli.kit.Rid)
				return result, err
			}

			item.SetAssoID(relations[id])
		}

		rstObj.Insts = rspItems
		result = append(result, rstObj)
	}

	return result, nil
}

func (cli *inst) SetMainlineParentInst(instID int64) error {
	if err := cli.updateMainlineAssociation(cli, instID); nil != err {
		blog.Errorf("[inst-inst] failed to update the mainline association, err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return err
	}

	return nil
}
func (cli *inst) SetMainlineChildInst(targetInst Inst) error {

	instID, err := targetInst.GetInstID()
	if err != nil {
		return err
	}

	childInsts, err := cli.GetMainlineChildInst()
	if nil != err {
		blog.Errorf("[inst-inst] failed to get the child inst, err:  %s, rid: %s", err.Error(), cli.kit.Rid)
		return err
	}
	for _, childInst := range childInsts {
		if err = cli.updateMainlineAssociation(childInst, instID); nil != err {
			blog.Errorf("[inst-inst] failed to set the mainline child inst, err: %s, rid: %s", err.Error(), cli.kit.Rid)
			return err
		}
	}

	id, err := cli.GetInstID()
	if err != nil {
		return err
	}

	if err = cli.updateMainlineAssociation(targetInst, id); nil != err {
		blog.Errorf("[inst-inst] failed to update the mainline association, err: %s, rid: %s", err.Error(), cli.kit.Rid)
		return err
	}

	return nil
}
