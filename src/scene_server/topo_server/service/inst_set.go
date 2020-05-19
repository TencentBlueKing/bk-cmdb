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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	parser "configcenter/src/common/paraparse"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/operation"
)

func (s *Service) BatchCreateSet(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("batch create set failed, parse app_id from url failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	type BatchCreateSetRequest struct {
		// shared fields
		Metadata          metadata.Metadata `json:"metadata"`
		BkSupplierAccount string            `json:"bk_supplier_account"`

		Sets []map[string]interface{} `json:"sets"`
	}
	batchBody := BatchCreateSetRequest{}
	if err := ctx.DecodeInto(&batchBody); err != nil {
		ctx.RespAutoError(err)
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDSet, &batchBody.Metadata)
	if nil != err {
		blog.Errorf("batch create set failed, get set model failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	type OneSetCreateResult struct {
		Index    int         `json:"index"`
		Data     interface{} `json:"data"`
		ErrorMsg string      `json:"error_message"`
	}
	batchCreateResult := make([]OneSetCreateResult, 0)
	var firstErr error
	var errMsg string
	for idx, set := range batchBody.Sets {
		if _, ok := set[common.MetadataField]; ok == false {
			set[common.MetadataField] = batchBody.Metadata
		}
		if _, ok := set[common.BkSupplierAccount]; ok == false {
			set[common.BkSupplierAccount] = batchBody.BkSupplierAccount
		}
		set[common.BKAppIDField] = bizID

		result, err := s.createSet(ctx.Kit, bizID, obj, set, &batchBody.Metadata)
		errMsg = ""
		if err != nil {
			errMsg = err.Error()
		}
		if err != nil && firstErr == nil {
			firstErr = err
		}
		if err != nil && blog.V(3) {
			blog.InfoJSON("batch create set at index:%d failed, data: %s, err: %s, rid: %s", idx, set, err.Error(), ctx.Kit.Rid)
		}
		batchCreateResult = append(batchCreateResult, OneSetCreateResult{
			Index:    idx,
			Data:     result,
			ErrorMsg: errMsg,
		})
	}
	ctx.RespEntityWithError(batchCreateResult, firstErr)
}

// CreateSet create a new set
func (s *Service) CreateSet(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDSet, dataWithMetadata.Metadata)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}
	resp, err := s.createSet(ctx.Kit, bizID, obj, dataWithMetadata.Data, dataWithMetadata.Metadata)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

func (s *Service) createSet(kit *rest.Kit, bizID int64, obj model.Object, data mapstr.MapStr, metaData *metadata.Metadata) (interface{}, error) {
	set, err := s.Core.SetOperation().CreateSet(kit, obj, bizID, data, metaData)
	if err != nil {
		return nil, err
	}

	setID, err := set.GetInstID()
	if err != nil {
		blog.Errorf("unexpected error, create set success, but get id field failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if _, err := s.Core.SetTemplateOperation().UpdateSetSyncStatus(kit, setID); err != nil {
		blog.Errorf("createSet success, but UpdateSetSyncStatus failed, setID: %d, err: %+v, rid: %s", setID, err, kit.Rid)
	}

	// auth: register set
	if err := s.AuthManager.RegisterSetByID(kit.Ctx, kit.Header, setID); err != nil {
		blog.Errorf("create set success,but register to iam failed, err:  %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed)
	}
	return set, nil
}

func (s *Service) CheckIsBuiltInSet(kit *rest.Kit, setIDs ...int64) errors.CCErrorCoder {
	// 检查是否时内置集群
	filter := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKSetIDField: map[string]interface{}{
				common.BKDBIN: setIDs,
			},
			common.BKDefaultField: map[string]interface{}{
				common.BKDBNE: common.DefaultFlagDefaultValue,
			},
		},
	}

	rsp, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, filter)
	if nil != err {
		blog.ErrorJSON("CheckIsBuiltInSet failed, ReadInstance failed, option: %s, err: %s, rid: %s", filter, err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if rsp.Result == false || rsp.Code != 0 {
		blog.ErrorJSON("ReadInstance failed, ReadInstance failed, option: %s, response: %s, rid: %s", filter, rsp, kit.Rid)
		return errors.New(rsp.Code, rsp.ErrMsg)
	}
	if rsp.Data.Count > 0 {
		return kit.CCError.CCError(common.CCErrorTopoForbiddenDeleteBuiltInSetModule)
	}
	return nil
}

func (s *Service) DeleteSets(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	data := struct {
		operation.OpCondition `json:",inline"`
		Metadata              *metadata.Metadata `json:"metadata"`
	}{}
	if err = ctx.DecodeInto(&data); nil != err {
		blog.Errorf("[api-set] failed to parse to the operation condition, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsIsInvalid, err.Error()))
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDSet, data.Metadata)

	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	setIDs := data.Delete.InstID
	// 检查是否时内置集群
	if err := s.CheckIsBuiltInSet(ctx.Kit, setIDs...); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// auth: deregister set
	if err := s.AuthManager.DeregisterSetByID(ctx.Kit.Ctx, ctx.Kit.Header, setIDs...); err != nil {
		blog.Errorf("delete sets failed, deregister sets from iam failed, %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed))
		return
	}

	err = s.Core.SetOperation().DeleteSet(ctx.Kit, obj, bizID, data.Delete.InstID, data.Metadata)
	if err != nil {
		ctx.RespAutoError(err)
		return

	}
	ctx.RespEntity(nil)

}

// DeleteSet delete the set
func (s *Service) DeleteSet(ctx *rest.Contexts) {
	if "batch" == ctx.Request.PathParameter("set_id") {
		s.DeleteSets(ctx)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the set id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "set id"))
		return
	}

	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDSet, md.Metadata)

	if nil != err {
		blog.Errorf("delete set failed, failed to search the set, %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// 检查是否时内置集群
	if err := s.CheckIsBuiltInSet(ctx.Kit, setID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// auth: deregister set
	if err := s.AuthManager.DeregisterSetByID(ctx.Kit.Ctx, ctx.Kit.Header, setID); err != nil {
		blog.Errorf("delete set failed, deregister set from iam failed, %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed))
		return
	}

	err = s.Core.SetOperation().DeleteSet(ctx.Kit, obj, bizID, []int64{setID}, md.Metadata)
	if err != nil {
		blog.Errorf("delete sets failed, %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// UpdateSet update the set
func (s *Service) UpdateSet(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the set id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "set id"))
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDSet, dataWithMetadata.Metadata)
	if nil != err {
		blog.Errorf("update set failed,failed to search the set, %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	err = s.Core.SetOperation().UpdateSet(ctx.Kit, dataWithMetadata.Data, obj, bizID, setID, dataWithMetadata.Metadata)
	if err != nil {
		blog.Errorf("update set failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// auth: update register set
	if err := s.AuthManager.UpdateRegisteredSetByID(ctx.Kit.Ctx, ctx.Kit.Header, setID); err != nil {
		blog.Errorf("update set success, but update registered set failed, %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed))
		return
	}
	ctx.RespEntity(nil)
}

// SearchSet search the set
func (s *Service) SearchSet(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-set]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	data := struct {
		parser.SearchParams `json:",inline"`
		Metadata            *metadata.Metadata `json:"metadata"`
	}{}
	if err = ctx.DecodeInto(&data); nil != err {
		blog.Errorf("search set failed, decode parameter condition failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsInvalid))
		return
	}
	paramsCond := data.SearchParams
	if paramsCond.Condition == nil {
		paramsCond.Condition = mapstr.New()
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDSet, data.Metadata)
	if nil != err {
		blog.Errorf("[api-set]failed to search the set, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	paramsCond.Condition[common.BKAppIDField] = bizID

	queryCond := &metadata.QueryInput{}
	queryCond.Condition = paramsCond.Condition
	queryCond.Fields = strings.Join(paramsCond.Fields, ",")
	page := metadata.ParsePage(paramsCond.Page)
	queryCond.Start = page.Start
	queryCond.Sort = page.Sort
	queryCond.Limit = page.Limit

	cnt, instItems, err := s.Core.SetOperation().FindSet(ctx.Kit, obj, queryCond)
	if nil != err {
		blog.Errorf("[api-set] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	ctx.RespEntity(result)
	return

}

// SearchSetBatch search the sets in one biz
func (s *Service) SearchSetBatch(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.SearchInstBatchOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(option.InstIDs) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "inst_ids"))
		return
	}

	if len(option.InstIDs) > common.BKMaxRecordsAtOnce {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrExceedMaxOperationRecordsAtOnce, common.BKMaxRecordsAtOnce))
		return
	}

	// set default value
	if option.Page.Limit == 0 {
		option.Page.Limit = common.BKDefaultLimit
	}
	if option.Page.IsIllegal() {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded))
		return
	}

	if len(option.Fields) == 0 {
		option.Fields = []string{common.BKSetIDField, common.BKSetNameField}
	}

	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
		common.BKSetIDField: mapstr.MapStr{
			common.BKDBIN: option.InstIDs,
		},
	}

	qc := &metadata.QueryCondition{
		Fields:    option.Fields,
		Page:      option.Page,
		Condition: cond,
	}
	instanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDSet, qc)
	if err != nil {
		blog.Errorf("SearchModuleBatch failed, http request failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if instanceResult.Code != 0 {
		blog.ErrorJSON("SearchModuleBatch failed, ReadInstance failed, filter: %s, response: %s, rid: %s", qc, instanceResult, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(instanceResult.Code, instanceResult.ErrMsg))
		return
	}
	ctx.RespEntity(instanceResult.Data)
}
