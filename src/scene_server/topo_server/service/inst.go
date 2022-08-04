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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	paraparse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/logics/inst"
)

var whiteList = []string{
	common.BKInnerObjIDHost,
}

// CreateInst create a new inst
func (s *Service) CreateInst(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// forbidden create inner model instance with common api
	if common.IsInnerModel(objID) {
		blog.Errorf("create %s instance with common create api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	exist, err := s.Logics.ObjectOperation().IsObjectExist(ctx.Kit, objID)
	if err != nil {
		blog.Errorf("failed to search the object(%s), err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if !exist {
		blog.Errorf("object(%s) is non-exist, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoModuleSelectFailed))
		return
	}

	setInst := make(mapstr.MapStr)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		setInst, err = s.Logics.InstOperation().CreateInst(ctx.Kit, objID, data)
		if err != nil {
			blog.Errorf("failed to create a new %s, %s, rid: %s", objID, err.Error(), ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(setInst)
}

// CreateManyInstance TODO
func (s *Service) CreateManyInstance(ctx *rest.Contexts) {
	data := new(metadata.CreateManyCommInst)
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	objID := ctx.Request.PathParameter(common.BKObjIDField)

	// forbidden create inner model instance with common api
	if common.IsInnerModel(objID) {
		blog.Errorf("create %s instance with common create api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	isMainline, err := s.Logics.AssociationOperation().IsMainlineObject(ctx.Kit, objID)
	if err != nil {
		blog.Errorf("check if object(%s) is mainline object failed, err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if isMainline {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommForbiddenOperateMainlineInstanceWithCommonAPI))
		return
	}

	var setInst *metadata.CreateManyCommInstResultDetail
	setInst, err = s.Logics.InstOperation().CreateManyInstance(ctx.Kit, objID, data.Details)
	if err != nil {
		blog.Errorf("failed to create %s new instances, err: %s, rid: %s", objID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(setInst)
}

// CreateInstsByImport batch create insts by excel import
func (s *Service) CreateInstsByImport(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	// forbidden create inner model instance with common api
	if common.IsInnerModel(objID) {
		blog.Errorf("create %s instance with common create api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	exist, err := s.Logics.ObjectOperation().IsObjectExist(ctx.Kit, objID)
	if err != nil {
		blog.Errorf("failed to search the object(%s), err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if !exist {
		blog.Errorf("object(%s) is non-exist, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoModuleSelectFailed))
		return
	}

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
	batchInfo := new(metadata.InstBatchInfo)
	if err := ctx.DecodeInto(batchInfo); err != nil {
		blog.Errorf("create instance failed, import object[%s] instance batch, but got invalid BatchInfo:[%v], "+
			"err: %+v, rid: %s", objID, batchInfo, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	setInst, err := s.Logics.InstOperation().CreateInstBatch(ctx.Kit, objID, batchInfo)
	if err != nil {
		blog.Errorf("failed to create new object %s, err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	ctx.RespEntity(setInst)
}

// DeleteInsts TODO
func (s *Service) DeleteInsts(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	data := metadata.OpCondition{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	// forbidden delete inner model instance with common api
	if common.IsInnerModel(objID) {
		blog.Errorf("delete %s instance with common delete api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	// forbidden create mainline instance with common api
	isMainline, err := s.Logics.AssociationOperation().IsMainlineObject(ctx.Kit, objID)
	if err != nil {
		blog.Errorf("check whether model %s to be mainline failed, err: %+v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if isMainline {
		// TODO add custom mainline instance param validation
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.InstOperation().DeleteInstByInstID(ctx.Kit, objID, data.Delete.InstID, true)
		if err != nil {
			blog.Errorf("delete instance failed, err: %v, objID: %s, instIDs: %+v, rid: %s", err, objID,
				data.Delete.InstID, ctx.Kit.Rid)
			return err
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
	if common.IsInnerModel(objID) {
		blog.Errorf("delete %s instance with common delete api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	if ctx.Request.PathParameter("inst_id") == "batch" {
		s.DeleteInsts(ctx)
		return
	}

	instID, err := strconv.ParseInt(ctx.Request.PathParameter("inst_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the inst id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "inst id"))
		return
	}

	// forbidden create mainline instance with common api
	isMainline, err := s.Logics.AssociationOperation().IsMainlineObject(ctx.Kit, objID)
	if err != nil {
		blog.Errorf("check whether model %s to be mainline failed, err: %+v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if isMainline {
		// TODO add custom mainline instance param validation
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.Logics.InstOperation().DeleteInstByInstID(ctx.Kit, objID, []int64{instID}, true); err != nil {
			blog.Errorf("delete instance failed err: %v, objID: %s, instID: %d, rid: %s", err, objID, instID,
				ctx.Kit.Rid)
			return err
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// UpdateInsts batch update insts
func (s *Service) UpdateInsts(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	data := metadata.OpCondition{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// forbidden update inner model instance with common api
	if common.IsInnerModel(objID) && !util.InArray(objID, whiteList) {
		blog.Errorf("update %s instance with common update api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	// check inst_id field to be not empty, is dangerous for empty inst_id field, which will update or delete all instance
	for idx, item := range data.Update {
		if item.InstID == 0 {
			blog.Errorf("%d's update item's field `inst_id` empty, rid: %s", idx, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "inst_id"))
			return
		}
	}
	for idx, instID := range data.Delete.InstID {
		if instID == 0 {
			blog.Errorf("%d's delete item's field `inst_id` empty, rid: %s", idx, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "inst_ids"))
			return
		}
	}

	// forbidden create mainline instance with common api
	isMainline, err := s.Logics.AssociationOperation().IsMainlineObject(ctx.Kit, objID)
	if err != nil {
		blog.Errorf("check whether model %s to be mainline failed, err: %+v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if isMainline {
		// TODO add custom mainline instance param validation
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		instIDField := metadata.GetInstIDFieldByObjID(objID)
		for _, item := range data.Update {
			cond := mapstr.MapStr{instIDField: item.InstID}
			err = s.Logics.InstOperation().UpdateInst(ctx.Kit, cond, item.InstInfo, objID)
			if err != nil {
				blog.Errorf("failed to update the object(%s) inst (%d), the data (%#v), err: %v, rid: %s",
					objID, item.InstID, data, err, ctx.Kit.Rid)
				return err
			}
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
	if common.IsInnerModel(objID) && !util.InArray(objID, whiteList) {
		blog.Errorf("update %s instance with common update api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	if ctx.Request.PathParameter("inst_id") == "batch" {
		s.UpdateInsts(ctx)
		return
	}

	instID, err := strconv.ParseInt(ctx.Request.PathParameter("inst_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the inst id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "inst_id"))
		return
	}

	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// forbidden create mainline instance with common api
	isMainline, err := s.Logics.AssociationOperation().IsMainlineObject(ctx.Kit, objID)
	if err != nil {
		blog.Errorf("check whether model %s to be mainline failed, err: %+v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if isMainline {
		// TODO add custom mainline instance param validation
	}

	cond := mapstr.MapStr{metadata.GetInstIDFieldByObjID(objID): instID}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.InstOperation().UpdateInst(ctx.Kit, cond, data, objID)
		if err != nil {
			blog.Errorf("failed to update the object(%s) inst (%s), the data (%#v), err: %v, rid: %s",
				objID, ctx.Request.PathParameter("inst_id"), data, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// SearchInsts search the insts
func (s *Service) SearchInsts(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// forbidden search inner model instance with common api
	if common.IsInnerModel(objID) {
		blog.Errorf("search %s instance with common search api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	queryCond := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(queryCond); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rsp, err := s.Logics.InstOperation().FindInst(ctx.Kit, objID, queryCond)
	if err != nil {
		blog.Errorf("failed to find the objects(%s), err: %V, rid: %s", ctx.Request.PathParameter("obj_id"), err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(rsp)
}

// SearchObjectInstances searches object instances with the input conditions.
func (s *Service) SearchObjectInstances(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// NOTE: NOT SUPPORT inner model search action in this interface.
	if common.IsInnerModel(objID) {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	// decode input parameter.
	input := new(metadata.CommonSearchFilter)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	input.ObjectID = objID

	// validate input parameter.
	if invalidKey, err := input.Validate(); err != nil {
		blog.Errorf("validate search instances input parameters failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, invalidKey))
		return
	}

	// set read preference.
	ctx.SetReadPreference(common.SecondaryPreferredMode)

	// search object instances.
	result, err := s.Logics.InstOperation().SearchObjectInstances(ctx.Kit, objID, input)
	if err != nil {
		blog.Errorf("search object[%s] instances failed, err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

// CountObjectInstances counts object instances num with the input conditions.
func (s *Service) CountObjectInstances(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// NOTE: NOT SUPPORT inner model count action in this interface.
	if common.IsInnerModel(objID) {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	// decode input parameter.
	input := new(metadata.CommonCountFilter)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	input.ObjectID = objID

	// validate input parameter.
	if invalidKey, err := input.Validate(); err != nil {
		blog.Errorf("validate count instances input parameters failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, invalidKey))
		return
	}

	// set read preference.
	ctx.SetReadPreference(common.SecondaryPreferredMode)

	// count object instances.
	result, err := s.Logics.InstOperation().CountObjectInstances(ctx.Kit, objID, input)
	if err != nil {
		blog.Errorf("count object[%s] instances failed, err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

// SearchInstAndAssociationDetail search the inst with association details
func (s *Service) SearchInstAndAssociationDetail(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// forbidden search inner model instance with common api
	if common.IsInnerModel(objID) {
		blog.Errorf("search %s instance with common search api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	queryCond := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(queryCond); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.Logics.InstOperation().FindInst(ctx.Kit, objID, queryCond)
	if err != nil {
		blog.Errorf("failed to find the objects(%s), err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// SearchInstUniqueFields search the instances' unique fields and including instances' id field
// no need to auth because it only get the unique fields
func (s *Service) SearchInstUniqueFields(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		blog.Errorf("search model unique url parameter id not a number, id: %s, err: %v, rid: %s",
			ctx.Request.PathParameter("id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, "id"))
		return
	}

	// get must check unique to judge if the instance exists
	cond := map[string]interface{}{
		common.BKObjIDField: objID,
		common.BKFieldID:    id,
	}
	uniqueResp, err := s.Engine.CoreAPI.CoreService().Model().ReadModelAttrUnique(ctx.Kit.Ctx, ctx.Kit.Header,
		metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("search model unique failed, cond: %#v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if uniqueResp.Count == 0 {
		blog.Errorf("model %s unique field not found, cond: %#v, rid: %s", objID, cond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrorTopObjectUniqueIndexNotFound, objID, id))
		return
	}

	if uniqueResp.Count != 1 {
		blog.Errorf("model %s unique field count > 1, cond: %#v, rid: %s", objID, cond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrTopoObjectUniqueSearchFailed))
		return
	}

	keyIDs := make([]int64, 0)
	for _, key := range uniqueResp.Info[0].Keys {
		keyIDs = append(keyIDs, int64(key.ID))
	}

	cond = map[string]interface{}{
		common.BKObjIDField: objID,
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: keyIDs,
		},
	}
	attrResp, err := s.Engine.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header,
		objID, &metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("search model attribute failed, cond: %s, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}

	if attrResp.Count <= 0 {
		blog.Errorf("unique model attribute count illegal, cond: %#v, rid: %s", cond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrTopoObjectUniqueSearchFailed))
		return
	}

	instIDKey := metadata.GetInstIDFieldByObjID(objID)
	keys := []string{instIDKey}
	attrIDNameMap := make(map[string]string, len(attrResp.Info))
	for _, attr := range attrResp.Info {
		keys = append(keys, attr.PropertyID)
		attrIDNameMap[attr.PropertyID] = attr.PropertyName
	}

	// construct the query inst condition
	queryCond := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(queryCond); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.Logics.InstOperation().FindInst(ctx.Kit, objID, queryCond)
	if err != nil {
		blog.Errorf("failed to find the object(%s) inst, err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.QueryUniqueFieldsData{InstResult: *result, UniqueAttribute: attrIDNameMap})
}

// SearchInstByObject search the inst of the object
func (s *Service) SearchInstByObject(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// forbidden search inner model instance with common api
	if common.IsInnerModel(objID) {
		blog.Errorf("search %s instance with common search api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	queryCond := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(queryCond); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rsp, err := s.Logics.InstOperation().FindInst(ctx.Kit, objID, queryCond)
	if err != nil {
		blog.Errorf("failed to find the objects(%s), err: %v, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(rsp)
}

// SearchInstByAssociation search inst by the association inst
func (s *Service) SearchInstByAssociation(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	if common.IsInnerModel(objID) {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	data := new(inst.AssociationParams)
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)

	result, err := s.Logics.InstOperation().FindInstByAssociationInst(ctx.Kit, objID, data)
	if nil != err {
		blog.Errorf("failed to find the objects(%s), err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

// SearchInstByInstID search the inst by inst ID
func (s *Service) SearchInstByInstID(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// forbidden search inner model instance with common api
	if common.IsInnerModel(objID) {
		blog.Errorf("search %s instance with common search api forbidden, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI))
		return
	}

	instID, err := strconv.ParseInt(ctx.Request.PathParameter("inst_id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrTopoInstSelectFailed, err.Error()))
		return
	}

	queryCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{metadata.GetInstIDFieldByObjID(objID): instID},
	}

	rsp, err := s.Logics.InstOperation().FindInst(ctx.Kit, objID, queryCond)
	if err != nil {
		blog.Errorf("failed to find the objects(%s), err: %v, rid: %s", ctx.Request.PathParameter("bk_obj_id"), err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(rsp)
}

// SearchInstsNames search instances names
// 只供前端使用，用于在业务的查询主机高级筛选页面根据集群名或者模块名模糊匹配获取相应的实例列表
func (s *Service) SearchInstsNames(ctx *rest.Contexts) {
	defErr := ctx.Kit.CCError
	option := new(metadata.SearchInstsNamesOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := option.Validate()
	if rawErr.ErrCode != 0 {
		blog.ErrorJSON("SearchInstsNames failed, Validate err: %s, option: %s, rid:%s", rawErr.ToCCError(defErr), option, ctx.Kit.Rid)
		ctx.RespAutoError(rawErr.ToCCError(defErr))
		return
	}

	filter := map[string]interface{}{
		common.BKAppIDField: option.BizID,
	}

	switch option.ObjID {
	case common.BKInnerObjIDSet:
		filter[common.BKSetNameField] = map[string]interface{}{
			common.BKDBLIKE: paraparse.SpecialCharChange(option.Name),
			"$options":      "i",
		}
	case common.BKInnerObjIDModule:
		filter[common.BKModuleNameField] = map[string]interface{}{
			common.BKDBLIKE: paraparse.SpecialCharChange(option.Name),
			"$options":      "i",
		}
	default:
		blog.Errorf("SearchInstsNames failed, unsupported obj: %s, rid: %s", option.ObjID, ctx.Kit.Rid)
		ctx.RespAutoError(defErr.CCErrorf(common.CCErrCommParamsInvalid, "bk_obj_id"))
		return
	}

	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.GetInstTableName(option.ObjID, ctx.Kit.SupplierAccount),
		Field:     metadata.GetInstNameFieldName(option.ObjID),
		Filter:    filter,
	}
	names, err := s.Engine.CoreAPI.CoreService().Common().GetDistinctField(ctx.Kit.Ctx, ctx.Kit.Header, distinctOpt)
	if err != nil {
		blog.ErrorJSON("GetDistinctField failed, err: %s, distinctOpt: %s, rid: %s", err, distinctOpt, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(names)
}

// SearchInstChildTopo search the child inst topo for a inst
func (s *Service) SearchInstChildTopo(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	instID, err := strconv.ParseInt(ctx.Request.PathParameter("inst_id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	exist, err := s.Logics.ObjectOperation().IsObjectExist(ctx.Kit, objID)
	if err != nil {
		blog.Errorf("failed to find the objects(%s), err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if !exist {
		blog.Errorf("object %s is non-exist, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField))
		return
	}

	_, instItems, err := s.Logics.InstOperation().FindInstChildTopo(ctx.Kit, objID, instID)
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
	if err != nil {
		blog.Errorf("path parameter inst_id invalid, object: %s, err: %+v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "inst_id"))
		return
	}

	fields := []string{common.BKObjIDField, common.BKObjIconField, common.BKObjNameField}
	obj, err := s.Logics.ObjectOperation().FindSingleObject(ctx.Kit, fields, objID)
	if nil != err {
		blog.Errorf("failed to find the objects(%s), err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	_, instItems, err := s.Logics.InstOperation().FindInstTopo(ctx.Kit, *obj, instID)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(instItems)
}

// SearchInstAssociation TODO
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

	input := metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: objID, common.BKInstIDField: instID},
				{common.BKAsstObjIDField: objID, common.BKAsstInstIDField: instID},
			},
		},
		Page: metadata.BasePage{
			Limit: int(limit),
			Start: int(start),
		},
	}

	if input.IsIllegal() {
		blog.ErrorJSON("parse page illegal, input: %s, rid: %s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded))
		return
	}

	blog.V(5).Infof("input:%#v, rid:%s", input, ctx.Kit.Rid)
	queryCond := &metadata.InstAsstQueryCondition{
		ObjID: objID,
		Cond:  input,
	}

	rsp, err := s.Engine.CoreAPI.CoreService().Association().ReadInstAssociation(ctx.Kit.Ctx, ctx.Kit.Header, queryCond)
	if err != nil {
		blog.Errorf("read instance association failed, err: %v, cond: %#v, rid: %s", err, queryCond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(rsp)
}

// SearchInstAssociationUI TODO
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
	infos, cnt, err := s.Logics.InstAssociationOperation().SearchInstAssociationUIList(ctx.Kit, objID, input)
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

	reqParams := new(metadata.RequestInstAssociationObjectID)
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

	infos, cnt, err := s.Logics.InstAssociationOperation().SearchInstAssociationSingleObjectInstInfo(ctx.Kit,
		reqParams.Condition.AssociationObjectID, input, reqParams.Condition.IsTargetObject)
	if err != nil {
		blog.Errorf("parse page illegal, input: %#v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(map[string]interface{}{
		"info":  infos,
		"count": cnt,
	})
}

// FindInsts find insts by cond
func (s *Service) FindInsts(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	data := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	sets, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, objID, data)
	if err != nil {
		blog.Errorf("failed to get inst, obj id: %s, err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(sets)
}
