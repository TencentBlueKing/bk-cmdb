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
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (a *association) DeleteMainlineAssociation(params types.ContextParams, objID string) error {

	targetObj, err := a.obj.FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find the target object(%s), error info is %s", objID, err.Error())
		return err
	}

	parentObj, err := targetObj.GetMainlineParentObject()
	if nil != err {
		blog.Errorf("[operation-asst] failed to find the object(%s)'s parent, error info is %s", objID, err.Error())
		return err
	}

	// update associations
	childObj, err := targetObj.GetMainlineChildObject()
	if nil != err && io.EOF != err {
		blog.Errorf("[operation-asst] failed to find the object(%s)'s child, error info is %s", objID, err.Error())
		return err
	}

	if err = a.ResetMainlineInstAssociatoin(params, targetObj); nil != err && io.EOF != err {
		blog.Errorf("[operation-asst] failed to delete the object(%s)'s insts, error info %s", objID, err.Error())
		return err
	}

	if nil != childObj {
		// FIX: 正常情况下 childObj 不可以能为 nil，只有在拓扑异常的时候才会出现
		if err = childObj.SetMainlineParentObject(parentObj.GetID()); nil != err && io.EOF != err {
			blog.Errorf("[operation-asst] failed to update the association, error info is %s", err.Error())
			return err
		}

	}
	// delete objects
	if err = a.obj.DeleteObject(params, targetObj.GetRecordID(), nil, false); nil != err && io.EOF != err {
		blog.Errorf("[operation-asst] failed to delete the object(%s), error info is %s", targetObj.GetID(), err.Error())
		return err
	}

	// delete this object related association.
	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldObjectID).Eq(targetObj.GetID())
	cond.Field(common.BKOwnerIDField).Eq(targetObj.GetSupplierAccount())
	if err = a.DeleteAssociation(params, cond); nil != err {
		blog.Errorf("[operation-asst] failed to delete the association, error info is %s", err.Error())
		return err
	}

	return nil
}

func (a *association) SearchMainlineAssociationTopo(params types.ContextParams, targetObj model.Object) ([]*metadata.MainlineObjectTopo, error) {

	foundObjIDMap := make(map[string]bool)
	results := make([]*metadata.MainlineObjectTopo, 0)
	for {
		resultsLen := len(results)

		tmpRst := &metadata.MainlineObjectTopo{}
		tmpRst.ObjID = targetObj.GetID()
		tmpRst.ObjName = targetObj.GetName()
		tmpRst.OwnerID = params.SupplierAccount

		parentObj, err := targetObj.GetMainlineParentObject()
		if nil == err {
			tmpRst.PreObjID = parentObj.GetID()
			tmpRst.PreObjName = parentObj.GetName()
		} else if nil != err && io.EOF != err {
			return nil, err
		}

		childObj, err := targetObj.GetMainlineChildObject()
		if nil == err {
			tmpRst.NextObj = childObj.GetID()
			tmpRst.NextName = childObj.GetName()
		} else if nil != err {
			if io.EOF != err {
				return nil, err
			}
			if _, ok := foundObjIDMap[tmpRst.ObjID]; !ok {
				results = append(results, tmpRst)
				foundObjIDMap[tmpRst.ObjID] = true
			}
			return results, nil
		}

		if _, ok := foundObjIDMap[tmpRst.ObjID]; !ok {
			results = append(results, tmpRst)
			foundObjIDMap[tmpRst.ObjID] = true
		}
		targetObj = childObj

		// detect infinite loop by checking whether there are new added objects in current loop.
		if resultsLen == len(results) {
			// merely return found objects here to avoid infinite loop.
			// returned results here maybe parts of all mainline objects.
			// better to prevent loop from taking shape seriously, at adding or editing association.
			return results, nil
		}
	}

}

func (a *association) CreateMainlineAssociation(params types.ContextParams, data *metadata.Association) (model.Object, error) {
	// find the mainline module's head, which is biz.
	bizObj, err := a.obj.FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("[operation-asst] failed to check the mainline topo level, error info is %s", err.Error())
		return nil, err
	}

	if data.AsstObjID == "" {
		blog.Errorf("[operation-asst] bk_asst_obj_id empty,rid:%s", util.GetHTTPCCRequestID(params.Header))
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedSet, common.BKAsstObjIDField)
	}

	if data.ClassificationID == "" {
		blog.Errorf("[operation-asst] bk_classification_id empty,rid:%s", util.GetHTTPCCRequestID(params.Header))
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedSet, common.BKClassificationIDField)
	}

	items, err := a.SearchMainlineAssociationTopo(params, bizObj)
	if nil != err {
		blog.Errorf("[operation-asst] failed to check the mainline topo level, error info is %s", err.Error())
		return nil, err
	}

	if len(items) >= params.MaxTopoLevel {
		blog.Errorf("[operation-asst] the mainline topo leve is %d, the max limit is %d", len(items), params.MaxTopoLevel)
		return nil, params.Err.Error(common.CCErrTopoBizTopoLevelOverLimit)
	}

	// find the mainline parent object
	parentObj, err := a.obj.FindSingleObject(params, data.AsstObjID)
	switch t := err.(type) {
	case nil:
	default:
		blog.Errorf("[operation-asst] failed to find the single object(%s), error info is %s", data.ObjectID, t.Error())
		return nil, t
	case errors.CCErrorCoder:
		if t.GetCode() == common.CCErrTopoObjectSelectFailed {
			blog.Errorf("[operation-asst] failed to find the single object(%s), error info is %s", data.ObjectID, t.Error())
			return nil, t
		}
	}

	// find the mainline child object for the parent
	childObj, err := parentObj.GetMainlineChildObject()
	if nil != err {
		blog.Errorf("[operation-asst] failed to find the child object for the object(%s), error info is %s", parentObj.GetID(), err.Error())
		return nil, err
	}

	// check and create the association mainline object
	if err = a.obj.IsValidObject(params, data.ObjectID); nil == err {
		blog.Errorf("[operation-asst] the object(%s) is duplicate", data.ObjectID)
		return nil, params.Err.Errorf(common.CCErrCommDuplicateItem, data.ObjectID)
	}

	objData := mapstr.MapStr{
		common.BKObjIDField:            data.ObjectID,
		common.BKObjNameField:          data.ObjectName,
		common.BKObjIconField:          data.ObjectIcon,
		common.BKClassificationIDField: data.ClassificationID,
	}
	currentObj, err := a.obj.CreateObject(params, true, objData)
	if err != nil {
		return nil, err
	}

	// update the mainline topo inst association
	if err = a.SetMainlineInstAssociation(params, parentObj, currentObj, childObj); nil != err {
		blog.Errorf("[operation-asst] failed set the mainline inst association, error info is %s", err.Error())
		return nil, err
	}

	if err = currentObj.CreateMainlineObjectAssociation(parentObj.GetID()); err != nil {
		blog.Errorf("[operation-asst] create mainline object[%s] association related to object[%s] failed, err: %v",
			currentObj.GetID(), parentObj.GetID(), err)
		return nil, err
	}

	if err = childObj.UpdateMainlineObjectAssociationTo(parentObj.GetID(), currentObj.GetID()); err != nil {
		blog.Errorf("[operation-asst] update mainline current object's[%s] child object[%s] association to current failed, err: %v",
			currentObj.GetID(), childObj.GetID(), err)
		return nil, err
	}

	return currentObj, nil
}
