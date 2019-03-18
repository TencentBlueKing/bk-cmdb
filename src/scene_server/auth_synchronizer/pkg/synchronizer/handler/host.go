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

package handler

import (
    "context"
    "io"

    "configcenter/src/apimachinery"
    "configcenter/src/common"
    "configcenter/src/common/blog"
    "configcenter/src/common/condition"
    "configcenter/src/common/metadata"
    "configcenter/src/scene_server/auth_synchronizer/pkg/synchronizer/meta"
    "configcenter/src/scene_server/auth_synchronizer/pkg/utils"
    hostutil "configcenter/src/scene_server/host_server/util"
    "configcenter/src/scene_server/topo_server/core/model"
    "configcenter/src/scene_server/topo_server/core/types"
)

// HandleHostSync do sync host of one business
func (ih *IAMHandler) HandleHostSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(meta.BusinessSimplify)

	// step1 get host by business from core service
	cond := hostutil.NewOperation().WithAppID(businessSimplify.BKAppIDField).Data()
	query := &metadata.QueryCondition{
	    Fields: []string{common.BKAppIDField, common.BKModuleIDField, common.BKSetIDField, common.BKHostIDField},
		Condition: cond,
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}

	header := utils.NewAPIHeaderByBusiness(&businessSimplify)
	// we only got common fields(bk_host_id, bk_module_id, bk_set_id, bk_biz_id) here,
	// got additional level of mainline must use some other operations.
	result, err := ih.CoreAPI.CoreService().Instance().ReadInstance(
		context.TODO(), *header, common.BKTableNameModuleHostConfig, query)

	if err != nil {
		blog.Errorf("list hosts by business failed: %v", err.Error())
	}
	blog.Infof("list hosts by business result: %+v", result)

	// step2 generate host layers
	// reference http interface: `api/v3/find/topomodelmainline`
	
	// step3 generate host resource id
	// step4 get host by business from iam
	// step5 diff step2 and step4 result
	// step6 register host not exist in iam
	// step7 deregister and register hosts that layers has changed
	// step8 deregister resource id that not in cmdb
	return nil
}

func SearchMainlineAssociationTopo(supplierAccount string, targetObj model.Object) ([]*metadata.MainlineObjectTopo, error) {

	foundObjIDMap := make(map[string]bool)
	results := make([]*metadata.MainlineObjectTopo, 0)
	for {
		tObject := targetObj.Object()

		resultsLen := len(results)
		tmpRst := &metadata.MainlineObjectTopo{}
		tmpRst.ObjID = tObject.ObjectID
		tmpRst.ObjName = tObject.ObjectName
		tmpRst.OwnerID = supplierAccount

		parentObj, err := targetObj.GetMainlineParentObject()
		if nil == err {
			tmpRst.PreObjID = parentObj.Object().ObjectID
			tmpRst.PreObjName = parentObj.Object().ObjectName
		} else if nil != err && io.EOF != err {
			return nil, err
		}

		childObj, err := targetObj.GetMainlineChildObject()
		if nil == err {
			tmpRst.NextObj = childObj.Object().ObjectID
			tmpRst.NextName = childObj.Object().ObjectName
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

func (o *object) FindSingleObject(params types.ContextParams, objectID string) (model.Object, error) {

    cond := condition.CreateCondition()
    cond.Field(common.BKObjIDField).Eq(objectID)

    objs, err := o.FindObject(params, cond)
    if nil != err {
        blog.Errorf("[api-inst] failed to find the supplier account(%s) objects(%s), err: %s", params.SupplierAccount, objectID, err.Error())
        return nil, err
    }
    for _, item := range objs {
        return item, nil
    }
    return nil, params.Err.New(common.CCErrTopoObjectSelectFailed, params.Err.Errorf(common.CCErrCommParamsIsInvalid, objectID).Error())
}

func (o *object) FindObject(params types.ContextParams, cond condition.Condition) ([]model.Object, error) {
    fCond := cond.ToMapStr()
    if nil != params.MetaData {
        fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
        fCond.Remove(metadata.BKMetadata)
    } else {
        fCond.Merge(metadata.BizLabelNotExist)
    }
    rsp, err := o.clientSet.CoreService().Model().ReadModel(context.Background(), params.Header, &metadata.QueryCondition{Condition: fCond})
    if nil != err {
        blog.Errorf("[operation-obj] failed to request the object controller, err: %s", err.Error())
        return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
    }

    if !rsp.Result {
        blog.Errorf("[operation-obj] failed to search the objects by the condition(%#v) , error info is %s", fCond, rsp.ErrMsg)
        return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
    }

    models := []metadata.Object{}
    for index := range rsp.Data.Info {
        models = append(models, rsp.Data.Info[index].Spec)
    }
    return model.CreateObject(params, o.clientSet, models), nil
}

// CreateObject create  objects
func CreateObject(params types.ContextParams, clientSet apimachinery.ClientSetInterface, objItems []metadata.Object) []Object {
    results := make([]Object, 0)
    for _, obj := range objItems {

        results = append(results, &object{
            obj:       obj,
            params:    params,
            clientSet: clientSet,
        })
    }

    return results
}
