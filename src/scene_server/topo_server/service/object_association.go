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

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateObjectAssociation create a new object association
func (s *Service) CreateObjectAssociation(ctx *rest.Contexts) {
	assoc := &metadata.Association{}
	if err := ctx.DecodeInto(assoc); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	var association *metadata.Association
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		association, err = s.Logics.AssociationOperation().CreateCommonAssociation(ctx.Kit, assoc)
		if nil != err {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(association)
}

// SearchObjectAssociation search object association by object id
func (s *Service) SearchObjectAssociation(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if data.Exists("condition") {
		// ATTENTION:
		// compatible with new query structures
		// the new condition format:
		// { "condition":{}}
		cond, err := data.MapStr("condition")
		if err != nil {
			blog.Errorf("get the condition failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsIsInvalid, err.Error()))
			return
		}

		if len(cond) == 0 {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
			return
		}

		needAuth := s.AuthManager.Enabled()
		if needAuth {
			condFields := []string{common.BKObjIDField, common.BKAsstObjIDField}
			for _, field := range condFields {
				if val, exist := cond.Get(field); exist {
					authResp, authorized, err := s.AuthManager.HasFindModelAuthUseObjID(ctx.Kit,
						[]string{util.GetStrByInterface(val)})
					if err != nil {
						ctx.RespAutoError(err)
						return
					}
					if !authorized {
						ctx.RespNoAuth(authResp)
						return
					}
					needAuth = false
					break
				}
			}
		}
		s.searchObjAssociationWithCond(ctx, cond, needAuth)
		return
	}

	objID, err := data.String(metadata.AssociationFieldObjectID)
	if err != nil {
		blog.Errorf("search object association, but get object id failed from: %v, err: %v, rid: %s",
			data, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	if len(objID) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	// authorize
	authResp, authorized, err := s.AuthManager.HasFindModelAuthUseObjID(ctx.Kit, []string{objID})
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	input := &metadata.QueryCondition{Condition: mapstr.MapStr{
		common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objID},
	}}
	resp, err := s.Engine.CoreAPI.CoreService().Association().ReadModelAssociation(ctx.Kit.Ctx, ctx.Kit.Header, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp.Info)
}

func (s *Service) searchObjAssociationWithCond(ctx *rest.Contexts, cond mapstr.MapStr, needAuth bool) {
	input := &metadata.QueryCondition{Condition: cond}
	resp, err := s.Engine.CoreAPI.CoreService().Association().ReadModelAssociation(ctx.Kit.Ctx, ctx.Kit.Header,
		input)
	if err != nil {
		blog.Errorf("search object association with cond[%v] failed, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if !needAuth {
		ctx.RespEntity(resp.Info)
		return
	}

	authInput := meta.ListAuthorizedResourcesParam{
		UserName:     ctx.Kit.User,
		ResourceType: meta.Model,
		Action:       meta.Find,
	}
	authorizedRes, err := s.AuthManager.Authorizer.ListAuthorizedResources(ctx.Kit.Ctx, ctx.Kit.Header, authInput)
	if err != nil {
		blog.Errorf("list authorized resources failed, user: %s, err: %v, rid: %s", ctx.Kit.User, err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrorTopoGetAuthorizedBusinessListFailed, "")
		return
	}

	if authorizedRes.IsAny {
		ctx.RespEntity(resp.Info)
		return
	}

	ids := make([]int64, 0)
	result := make([]metadata.Association, 0)
	for _, resourceID := range authorizedRes.Ids {
		id, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			blog.Errorf("get authorized object id failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		ctx.RespEntity(result)
		return
	}

	queryCond := mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: ids}}
	query := &metadata.QueryCondition{Condition: queryCond, DisableCounter: true}
	modelResp, err := s.Engine.CoreAPI.CoreService().Model().ReadModel(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	authMap := make(map[string]struct{})
	for _, model := range modelResp.Info {
		authMap[model.ObjectID] = struct{}{}
	}

	for _, association := range resp.Info {
		_, haveObjIDAuth := authMap[association.ObjectID]
		_, haveAsstObjIDAuth := authMap[association.AsstObjID]
		if haveObjIDAuth || haveAsstObjIDAuth {
			result = append(result, association)
		}
	}

	ctx.RespEntity(result)
	return
}

// DeleteObjectAssociation delete object association
func (s *Service) DeleteObjectAssociation(ctx *rest.Contexts) {

	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		blog.Errorf("delete object association failed, got a invalid object association id[%v], err: %v, rid: %s",
			ctx.Request.PathParameter("id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoInvalidObjectAssociationID))
		return
	}

	if id <= 0 {
		blog.Errorf("delete object association failed, got a invalid objAsst id[%d], rid: %s", id, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoInvalidObjectAssociationID))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.AssociationOperation().DeleteAssociationWithPreCheck(ctx.Kit, id)
		if err != nil {
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

// UpdateObjectAssociation update object association
func (s *Service) UpdateObjectAssociation(ctx *rest.Contexts) {
	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		blog.Errorf("update object association, but got invalid id[%v], err: %v, rid: %s",
			ctx.Request.PathParameter("id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommParamsIsInvalid))
		return
	}

	data := mapstr.New()
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.AssociationOperation().UpdateObjectAssociation(ctx.Kit, data, id)
		if err != nil {
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

// ImportInstanceAssociation import instance  association
func (s *Service) ImportInstanceAssociation(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")
	request := new(metadata.RequestImportAssociation)
	if err := ctx.DecodeInto(request); err != nil {
		blog.Errorf("ImportInstanceAssociation, json unmarshal error, objID: %s, err: %v, rid: %s",
			objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	var ret metadata.ResponeImportAssociationData
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Logics.ImportAssociationOperation().ImportInstAssociation(ctx.Kit, s.Language, objID,
			request.AssociationInfoMap, request.AsstObjectUniqueIDMap, request.ObjectUniqueID)
		if err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(ret, txnErr)
		return
	}
	ctx.RespEntity(ret)
}

// SearchModuleAssociation search model association
func (s *Service) SearchModuleAssociation(ctx *rest.Contexts) {
	data := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Engine.CoreAPI.CoreService().Association().ReadModelAssociation(ctx.Kit.Ctx, ctx.Kit.Header, data)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	unique := make(map[string]struct{}, 0)
	objIDs := make([]string, 0)
	for _, model := range resp.Info {
		if _, ok := unique[model.ObjectID]; !ok {
			unique[model.ObjectID] = struct{}{}
			objIDs = append(objIDs, model.ObjectID)
		}
	}
	// authorize
	authResp, authorized, err := s.AuthManager.HasFindModelAuthUseObjID(ctx.Kit, objIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	ctx.RespEntity(resp)
}

// FindAssociationByObjectAssociationID 根据关联关系bk_obj_asst_id 获取关联信息
// 专用方法，提供给关联关系导入使用
func (s *Service) FindAssociationByObjectAssociationID(ctx *rest.Contexts) {

	objID := ctx.Request.PathParameter(common.BKObjIDField)
	request := new(metadata.FindAssociationByObjectAssociationIDRequest)
	if err := ctx.DecodeInto(request); err != nil {
		blog.Errorf("FindObjectByObjectAssociationID, json unmarshal error, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	// authorize
	authResp, authorized, err := s.AuthManager.HasFindModelAuthUseObjID(ctx.Kit, []string{objID})
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	var association []metadata.Association
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		association, err = s.Logics.ImportAssociationOperation().FindAssociationByObjectAssociationID(ctx.Kit, objID,
			request.ObjAsstIDArr)
		if err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(association)
}
