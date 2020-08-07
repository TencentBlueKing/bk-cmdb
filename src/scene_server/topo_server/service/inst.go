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

package service

import (
	"strconv"
	"strings"

	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	paraparse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/operation"
)

var whiteList = []string{
	common.BKInnerObjIDHost,
}

// CreateInst create a new inst
func (s *Service) CreateInst(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}
	data := dataWithMetadata.Data

	// forbidden create inner model instance with common api
	if common.IsInnerModel(objID) == true {
		blog.V(5).Infof("CreateInst failed, create %s instance with common create api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, dataWithMetadata.Metadata)
	if nil != err {
		blog.Errorf("failed to search the inst, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// forbidden create mainline instance with common api
	isMainline, err := obj.IsMainlineObject()
	if err != nil {
		blog.Errorf("CreateInst failed, check whether model %s to be mainline failed, err: %+v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if isMainline == true {
		// TODO add custom mainline instance param validation
	}

	if data.Exists("BatchInfo") {
		/*
			   BatchInfo data format:
			    {
			      "BatchInfo": {
			        "4": { // excel line number
			          "bk_inst_id": 1,
			          "bk_inst_key": "a22",
			          "bk_inst_name": "a11",
			          "bk_version": "121",
			          "import_from": "1"
					}
				  },
			      "input_type": "excel"
			    }
		*/
		batchInfo := new(operation.InstBatchInfo)
		if err := data.MarshalJSONInto(batchInfo); err != nil {
			blog.Errorf("create instance failed, import object[%s] instance batch, but got invalid BatchInfo:[%v], err: %+v, rid: %s", objID, batchInfo, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
			return
		}

		var setInst *operation.BatchResult
		txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
			var err error
			setInst, err = s.Core.InstOperation().CreateInstBatch(ctx.Kit, obj, batchInfo, dataWithMetadata.Metadata)
			if nil != err {
				blog.Errorf("failed to create new object %s, %s, rid: %s", objID, err.Error(), ctx.Kit.Rid)
				return err
			}

			// auth register new created
			if len(setInst.SuccessCreated) != 0 {
				if err := s.AuthManager.RegisterInstancesByID(ctx.Kit.Ctx, ctx.Kit.Header, objID, setInst.SuccessCreated...); err != nil {
					blog.Errorf("create instance success, but register instances to iam failed, instances: %+v, err: %+v, rid: %s", setInst.SuccessCreated, err, ctx.Kit.Rid)
					return ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed)
				}
			}

			// auth update registered instances
			if len(setInst.SuccessUpdated) != 0 {
				if err := s.AuthManager.UpdateRegisteredInstanceByID(ctx.Kit.Ctx, ctx.Kit.Header, objID, setInst.SuccessUpdated...); err != nil {
					blog.Errorf("update registered instances to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
					return ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed)
				}
			}
			return nil
		})

		if txnErr != nil {
			ctx.RespAutoError(txnErr)
			return
		}
		ctx.RespEntity(setInst)
		return
	}

	var setInst inst.Inst
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		setInst, err = s.Core.InstOperation().CreateInst(ctx.Kit, obj, data)
		if nil != err {
			blog.Errorf("failed to create a new %s, %s, rid: %s", objID, err.Error(), ctx.Kit.Rid)
			return err
		}

		instanceID, err := setInst.GetInstID()
		if err != nil {
			blog.Errorf("create instance failed, unexpected error, create instance success, but get id failed, instance: %+v, err: %+v, rid: %s", setInst, err, ctx.Kit.Rid)
			return err
		}

		// auth: register instances to iam
		if err := s.AuthManager.RegisterInstancesByID(ctx.Kit.Ctx, ctx.Kit.Header, objID, instanceID); err != nil {
			blog.Errorf("create instance success, but register instance to iam failed, instance: %d, err: %+v, rid: %s", instanceID, err, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(setInst.ToMapStr())
}

func (s *Service) DeleteInsts(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	data := struct {
		operation.OpCondition `json:",inline"`
		Metadata              *metadata.Metadata `json:"metadata"`
	}{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	deleteCondition := data.OpCondition

	// forbidden delete inner model instance with common api
	if common.IsInnerModel(objID) == true {
		blog.V(5).Infof("DeleteInsts failed, create %s instance with common create api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, data.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// forbidden create mainline instance with common api
	isMainline, err := obj.IsMainlineObject()
	if err != nil {
		blog.Errorf("DeleteInsts failed, check whether model %s to be mainline failed, err: %+v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if isMainline == true {
		// TODO add custom mainline instance param validation
	}

	authInstances := make([]extensions.InstanceSimplify, 0)
	input := &metadata.QueryInput{
		Condition: map[string]interface{}{
			obj.GetInstIDFieldName(): map[string]interface{}{
				common.BKDBIN: deleteCondition.Delete.InstID,
			}}}

	_, insts, err := s.Core.InstOperation().FindInst(ctx.Kit, obj, input, false)
	if nil != err {
		blog.Errorf("DeleteInst failed, find authInstances to be deleted failed, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	for _, inst := range insts {
		instID, _ := inst.GetInstID()
		instName, _ := inst.GetInstName()
		instBizID, _ := inst.GetBizID()
		authInstances = append(authInstances, extensions.InstanceSimplify{
			InstanceID: instID,
			Name:       instName,
			BizID:      instBizID,
			ObjectID:   objID,
		})
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		if err = s.Core.InstOperation().DeleteInstByInstID(ctx.Kit, obj, deleteCondition.Delete.InstID, true); err != nil {
			blog.Errorf("DeleteInst failed, DeleteInstByInstID failed, err: %s, objID: %s, instIDs: %+v, rid: %s", err.Error(), objID, deleteCondition.Delete.InstID, ctx.Kit.Rid)
			return err
		}

		// auth: deregister resources
		if err := s.AuthManager.DeregisterInstances(ctx.Kit.Ctx, ctx.Kit.Header, authInstances...); err != nil {
			blog.Errorf("batch delete instance failed, deregister instance failed, instID: %d, err: %s, rid: %s", deleteCondition.Delete.InstID, err, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// DeleteInst delete the inst
func (s *Service) DeleteInst(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// forbidden delete inner model instance with common api
	if common.IsInnerModel(objID) == true {
		blog.V(5).Infof("CreateInst failed, create %s instance with common create api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	if "batch" == ctx.Request.PathParameter("inst_id") {
		s.DeleteInsts(ctx)
		return
	}

	instID, err := strconv.ParseInt(ctx.Request.PathParameter("inst_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-inst]failed to parse the inst id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "inst id"))
		return
	}

	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, md.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// forbidden create mainline instance with common api
	isMainline, err := obj.IsMainlineObject()
	if err != nil {
		blog.Errorf("DeleteInst failed, check whether model %s to be mainline failed, err: %+v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if isMainline == true {
		// TODO add custom mainline instance param validation
	}

	authInstances := make([]extensions.InstanceSimplify, 0)
	_, insts, err := s.Core.InstOperation().FindInst(ctx.Kit, obj, &metadata.QueryInput{Condition: map[string]interface{}{obj.GetInstIDFieldName(): instID}}, false)
	if nil != err {
		blog.Errorf("DeleteInst failed, find authInstances to be deleted failed, error info is %s, rid: %s", err.Error(), ctx.Kit)
		ctx.RespAutoError(err)
		return
	}
	for _, inst := range insts {
		instName, _ := inst.GetInstName()
		instBizID, _ := inst.GetBizID()
		authInstances = append(authInstances, extensions.InstanceSimplify{
			InstanceID: instID,
			Name:       instName,
			BizID:      instBizID,
			ObjectID:   objID,
		})
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		if err := s.Core.InstOperation().DeleteInstByInstID(ctx.Kit, obj, []int64{instID}, true); err != nil {
			blog.Errorf("DeleteInst failed, DeleteInstByInstID failed, err: %s, objID: %s, instID: %d, rid: %s", err.Error(), objID, instID, ctx.Kit.Rid)
			return err
		}

		// auth: deregister resources
		if err := s.AuthManager.DeregisterInstances(ctx.Kit.Ctx, ctx.Kit.Header, authInstances...); err != nil {
			blog.Errorf("delete instance failed, deregister instance failed, instID: %d, err: %s, rid: %s", instID, err, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (s *Service) UpdateInsts(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	data := struct {
		operation.OpCondition `json:",inline"`
		Metadata              *metadata.Metadata `json:"metadata"`
	}{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	updateCondition := data.OpCondition

	// forbidden create inner model instance with common api
	if common.IsInnerModel(objID) == true && util.InArray(objID, whiteList) == false {
		blog.V(5).Infof("UpdateInsts failed, update %s instance with common create api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	// check inst_id field to be not empty, is dangerous for empty inst_id field, which will update or delete all instance
	for idx, item := range updateCondition.Update {
		if item.InstID == 0 {
			blog.Errorf("update instance failed, %d's update item's field `inst_id` emtpy, rid: %s", idx, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsInvalid))
			return
		}
	}
	for idx, instID := range updateCondition.Delete.InstID {
		if instID == 0 {
			blog.Errorf("update instance failed, %d's delete item's field `inst_id` emtpy, rid: %s", idx, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsInvalid))
			return
		}
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, data.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// forbidden create mainline instance with common api
	isMainline, err := obj.IsMainlineObject()
	if err != nil {
		blog.Errorf("UpdateInsts failed, check whether model %s to be mainline failed, err: %+v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if isMainline == true {
		// TODO add custom mainline instance param validation
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		instanceIDs := make([]int64, 0)
		for _, item := range updateCondition.Update {
			instanceIDs = append(instanceIDs, item.InstID)
			cond := condition.CreateCondition()
			cond.Field(obj.GetInstIDFieldName()).Eq(item.InstID)
			err = s.Core.InstOperation().UpdateInst(ctx.Kit, item.InstInfo, obj, cond, item.InstID, data.Metadata)
			if nil != err {
				blog.Errorf("[api-inst] failed to update the object(%s) inst (%d),the data (%#v), error info is %s, rid: %s", obj.Object().ObjectID, item.InstID, data, err.Error(), ctx.Kit.Rid)
				return err
			}
		}

		// auth: update resources
		if err := s.AuthManager.UpdateRegisteredInstanceByID(ctx.Kit.Ctx, ctx.Kit.Header, objID, instanceIDs...); err != nil {
			blog.Errorf("update inst success, but update register to iam failed, instanceIDs: %+v, err: %+v, rid: %s", instanceIDs, err, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// UpdateInst update the inst
func (s *Service) UpdateInst(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// forbidden update inner model instance with common api
	if common.IsInnerModel(objID) == true && util.InArray(objID, whiteList) == false {
		blog.V(5).Infof("CreateInst failed, create %s instance with common create api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	if "batch" == ctx.Request.PathParameter("inst_id") {
		s.UpdateInsts(ctx)
		return
	}

	instID, err := strconv.ParseInt(ctx.Request.PathParameter("inst_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-inst]failed to parse the inst id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "inst id"))
		return
	}

	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}
	data := dataWithMetadata.Data
	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, dataWithMetadata.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", objID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// forbidden create mainline instance with common api
	isMainline, err := obj.IsMainlineObject()
	if err != nil {
		blog.Errorf("UpdateInsts failed, check whether model %s to be mainline failed, err: %+v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if isMainline == true {
		// TODO add custom mainline instance param validation
	}

	// this is a special logic for mainline object instance.
	// for auth reason, the front's request add metadata for mainline model's instance update.
	// but actually, it's should not add metadata field in the request.
	// so, we need remove it from the data if it's a mainline model instance.
	yes, err := s.Core.AssociationOperation().IsMainlineObject(ctx.Kit, objID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if yes {
		data.Remove("metadata")
	}

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		err = s.Core.InstOperation().UpdateInst(ctx.Kit, data, obj, cond, instID, dataWithMetadata.Metadata)
		if nil != err {
			blog.Errorf("[api-inst] failed to update the object(%s) inst (%s),the data (%#v), error info is %s, rid: %s", obj.Object().ObjectID, ctx.Request.PathParameter("inst_id"), data, err.Error(), ctx.Kit.Rid)
			return err
		}

		// auth: deregister resources
		if err := s.AuthManager.UpdateRegisteredInstanceByID(ctx.Kit.Ctx, ctx.Kit.Header, objID, instID); err != nil {
			blog.Error("update inst failed, authorization failed, instID: %d, err: %+v, rid: %s", instID, err, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// SearchInst search the inst
func (s *Service) SearchInsts(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	data := struct {
		paraparse.SearchParams `json:",inline"`
		Metadata               *metadata.Metadata `json:"metadata"`
	}{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, data.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	//	if nil != ctx.Kit.MetaData {
	//		data.Set(metadata.BKMetadata, *ctx.Kit.MetaData)
	//	}
	// construct the query inst condition
	queryCond := data.SearchParams
	if queryCond.Condition == nil {
		queryCond.Condition = mapstr.New()
	}
	page := metadata.ParsePage(queryCond.Page)
	query := &metadata.QueryInput{}
	query.Condition = queryCond.Condition
	query.Fields = strings.Join(queryCond.Fields, ",")
	query.Limit = page.Limit
	query.Sort = page.Sort
	query.Start = page.Start

	cnt, instItems, err := s.Core.InstOperation().FindInst(ctx.Kit, obj, query, false)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	ctx.RespEntity(result)
}

// SearchInstAndAssociationDetail search the inst with association details
func (s *Service) SearchInstAndAssociationDetail(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	data := struct {
		paraparse.SearchParams `json:",inline"`
		Metadata               *metadata.Metadata `json:"metadata"`
	}{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, data.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// construct the query inst condition
	queryCond := data.SearchParams
	if queryCond.Condition == nil {
		queryCond.Condition = mapstr.New()
	}
	page := metadata.ParsePage(queryCond.Page)
	query := &metadata.QueryInput{}
	query.Condition = queryCond.Condition
	query.Fields = strings.Join(queryCond.Fields, ",")
	query.Limit = page.Limit
	query.Sort = page.Sort
	query.Start = page.Start

	cnt, instItems, err := s.Core.InstOperation().FindInst(ctx.Kit, obj, query, true)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	ctx.RespEntity(result)
}

// SearchInstByObject search the inst of the object
func (s *Service) SearchInstByObject(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	data := struct {
		paraparse.SearchParams `json:",inline"`
		Metadata               *metadata.Metadata `json:"metadata"`
	}{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, data.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	queryCond := data.SearchParams
	if queryCond.Condition == nil {
		queryCond.Condition = mapstr.New()
	}
	page := metadata.ParsePage(queryCond.Page)
	query := &metadata.QueryInput{}
	query.Condition = queryCond.Condition
	query.Fields = strings.Join(queryCond.Fields, ",")
	query.Limit = page.Limit
	query.Sort = page.Sort
	query.Start = page.Start
	cnt, instItems, err := s.Core.InstOperation().FindInst(ctx.Kit, obj, query, false)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	ctx.RespEntity(result)
}

// SearchInstByAssociation search inst by the association inst
func (s *Service) SearchInstByAssociation(ctx *rest.Contexts) {
	data := struct {
		Metadata                    *metadata.Metadata `json:"metadata"`
		operation.AssociationParams `json:",inline"`
	}{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	objID := ctx.Request.PathParameter("bk_obj_id")

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, data.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	cnt, instItems, err := s.Core.InstOperation().FindInstByAssociationInst(ctx.Kit, obj, &data.AssociationParams)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	ctx.RespEntity(result)
}

// SearchInstByInstID search the inst by inst ID
func (s *Service) SearchInstByInstID(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	instID, err := strconv.ParseInt(ctx.Request.PathParameter("inst_id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrTopoInstSelectFailed, err.Error()))
		return
	}

	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, md.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	queryCond := &metadata.QueryInput{}
	queryCond.Condition = cond.ToMapStr()

	cnt, instItems, err := s.Core.InstOperation().FindInst(ctx.Kit, obj, queryCond, false)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	ctx.RespEntity(result)
}

// SearchInstChildTopo search the child inst topo for a inst
func (s *Service) SearchInstChildTopo(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	instID, err := strconv.ParseInt(ctx.Request.PathParameter("inst_id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, md.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", objID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	query := &metadata.QueryInput{}
	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)

	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	_, instItems, err := s.Core.InstOperation().FindInstChildTopo(ctx.Kit, obj, instID, query)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(instItems)
}

// SearchInstTopo search the inst topo
func (s *Service) SearchInstTopo(ctx *rest.Contexts) {

	objID := ctx.Request.PathParameter("bk_obj_id")
	instID, err := strconv.ParseInt(ctx.Request.PathParameter("inst_id"), 10, 64)
	if nil != err {
		blog.Errorf("search inst topo failed, path parameter inst_id invalid, object: %s inst_id: %s, err: %+v, rid: %s", objID, ctx.Request.PathParameter("inst_id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, objID, md.Metadata)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	query := &metadata.QueryInput{}
	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)

	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	_, instItems, err := s.Core.InstOperation().FindInstTopo(ctx.Kit, obj, instID, query)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(instItems)
}

// Deprecated 2019-09-30 废弃接口
func (s *Service) SearchInstAssociation(ctx *rest.Contexts) {

	objID := ctx.Request.PathParameter("bk_obj_id")
	instID, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "id"))
		return
	}
	start, err := strconv.ParseInt(ctx.Request.PathParameter("start"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "start"))
		return
	}
	limit, err := strconv.ParseInt(ctx.Request.PathParameter("limit"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "limit"))
		return
	}

	cond := condition.CreateCondition()
	condOR := cond.NewOR()
	condOR.Item(map[string]interface{}{common.BKObjIDField: objID, common.BKInstIDField: instID})
	condOR.Item(map[string]interface{}{common.BKAsstObjIDField: objID, common.BKAsstInstIDField: instID})
	input := &metadata.QueryCondition{
		Condition: cond.ToMapStr(),
		Page: metadata.BasePage{
			Limit: int(limit),
			Start: int(start),
		},
	}

	if input.IsIllegal() {
		blog.ErrorJSON("parse page illegal, input:%s,rid:%s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded))
		return
	}

	blog.V(5).Infof("input:%#v, rid:%s", input, ctx.Kit.Rid)
	infos, cnt, err := s.Core.AssociationOperation().SearchInstAssociationList(ctx.Kit, input)
	if err != nil {
		blog.ErrorJSON("parse page illegal, input:%s, err:%s, rid:%s", input, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(map[string]interface{}{
		"info":  infos,
		"count": cnt,
		"page":  input.Page,
	})
}

func (s *Service) SearchInstAssociationUI(ctx *rest.Contexts) {

	objID := ctx.Request.PathParameter(common.BKObjIDField)
	instID, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "id"))
		return
	}
	start, err := strconv.ParseInt(ctx.Request.PathParameter("start"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "start"))
		return
	}
	limit, err := strconv.ParseInt(ctx.Request.PathParameter("limit"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "limit"))
		return
	}

	cond := condition.CreateCondition()
	condOR := cond.NewOR()
	condOR.Item(map[string]interface{}{common.BKObjIDField: objID, common.BKInstIDField: instID})
	condOR.Item(map[string]interface{}{common.BKAsstObjIDField: objID, common.BKAsstInstIDField: instID})
	input := &metadata.QueryCondition{
		Condition: cond.ToMapStr(),
		Page: metadata.BasePage{
			Limit: int(limit),
			Start: int(start),
		},
	}

	if input.IsIllegal() {
		blog.ErrorJSON("parse page illegal, input:%s,rid:%s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded))
		return
	}

	blog.V(5).Infof("input:%#v, rid:%s", input, ctx.Kit.Rid)
	infos, cnt, err := s.Core.AssociationOperation().SearchInstAssociationUIList(ctx.Kit, objID, input)
	if err != nil {
		blog.ErrorJSON("parse page illegal, input:%s, err:%s, rid:%s", input, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(map[string]interface{}{
		"data":              infos,
		"association_count": cnt,
		"page":              input.Page,
	})
}

// SearchInstAssociationWithOtherObject  要求根据实例信息（实例的模型ID，实例ID）和模型ID（关联关系中的源，目的模型ID） 返回实例关联或者被关联模型实例得数据。
func (s *Service) SearchInstAssociationWithOtherObject(ctx *rest.Contexts) {

	reqParams := &metadata.RequestInstAssociationObjectID{}
	if err := ctx.DecodeInto(reqParams); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if reqParams.Condition.ObjectID == "" {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField))
		return
	}
	if reqParams.Condition.InstID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKInstIDField))
		return
	}
	if reqParams.Condition.AssociationObjectID == "" {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "association_obj_id"))
		return
	}

	cond := condition.CreateCondition()
	if reqParams.Condition.IsTargetObject {
		// 作为目标模型
		cond.Field(common.BKAsstObjIDField).Eq(reqParams.Condition.ObjectID)
		cond.Field(common.BKAsstInstIDField).Eq(reqParams.Condition.InstID)
		cond.Field(common.BKObjIDField).Eq(reqParams.Condition.AssociationObjectID)
	} else {
		// 作为源模型
		cond.Field(common.BKObjIDField).Eq(reqParams.Condition.ObjectID)
		cond.Field(common.BKInstIDField).Eq(reqParams.Condition.InstID)
		cond.Field(common.BKAsstObjIDField).Eq(reqParams.Condition.AssociationObjectID)
	}

	input := &metadata.QueryCondition{
		Condition: cond.ToMapStr(),
		Page:      reqParams.Page,
	}

	if input.IsIllegal() {
		blog.ErrorJSON("parse page illegal, input:%s,rid:%s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded))
		return
	}

	blog.V(5).Infof("input:%#v, rid:%s", input, ctx.Kit.Rid)
	infos, cnt, err := s.Core.AssociationOperation().SearchInstAssociationSingleObjectInstInfo(ctx.Kit, reqParams.Condition.AssociationObjectID, input)
	if err != nil {
		blog.ErrorJSON("parse page illegal, input:%s, err:%s, rid:%s", input, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(map[string]interface{}{
		"info":  infos,
		"count": cnt,
		"page":  input.Page,
	})
}
