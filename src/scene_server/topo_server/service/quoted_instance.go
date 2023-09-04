/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// BatchCreateQuotedInstance batch create quoted instances.
func (s *Service) BatchCreateQuotedInstance(cts *rest.Contexts) {
	opt := new(metadata.BatchCreateQuotedInstOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err.ErrCode != 0 {
		cts.RespAutoError(err.ToCCError(cts.Kit.CCError))
		return
	}

	// authorize, ** right now use source instance create or update action to authorize, change it when confirmed **
	instIDs := make([]int64, 0)
	for _, data := range opt.Data {
		instIDVal, exists := data[common.BKInstIDField]
		if !exists {
			continue
		}

		instID, err := util.GetInt64ByInterface(instIDVal)
		if err != nil {
			cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKInstIDField))
			return
		}

		if instID == 0 {
			continue
		}
		instIDs = append(instIDs, instID)
	}

	if len(instIDs) > 0 {
		instIDs = util.IntArrayUnique(instIDs)

		uAuthErr := s.AuthManager.AuthorizeByInstanceID(cts.Kit.Ctx, cts.Kit.Header, meta.Update, opt.ObjID, instIDs...)
		cAuthErr := s.AuthManager.AuthorizeByInstanceID(cts.Kit.Ctx, cts.Kit.Header, meta.Create, opt.ObjID, instIDs...)
		if uAuthErr != nil && cAuthErr != nil {
			blog.Errorf("authorize failed, create err: %v, update err: %v, rid: %s", cAuthErr, uAuthErr, cts.Kit.Rid)
			cts.RespAutoError(uAuthErr)
			return
		}
	}

	// get quoted object id
	objID, err := s.Logics.ModelQuoteOperation().GetQuotedObjID(cts.Kit, opt.ObjID, opt.PropertyID)
	if err != nil {
		cts.RespAutoError(err)
		return
	}

	res := new(metadata.BatchCreateResult)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(cts.Kit.Ctx, cts.Kit.Header, func() error {
		// create quoted instances
		ids, err := s.Engine.CoreAPI.CoreService().ModelQuote().BatchCreateQuotedInstance(cts.Kit.Ctx, cts.Kit.Header,
			objID, opt.Data)
		if err != nil {
			blog.Errorf("create quoted instances failed, err: %v, rid: %s", err, cts.Kit.Rid)
			return err
		}

		// generate and save audit logs
		for i := range opt.Data {
			opt.Data[i][common.BKFieldID] = ids[i]
		}

		audit := auditlog.NewQuotedInstAuditLog(s.Engine.CoreAPI.CoreService())
		genAuditParams := auditlog.NewGenerateAuditCommonParameter(cts.Kit, metadata.AuditCreate)
		auditLogs, ccErr := audit.GenerateAuditLog(genAuditParams, objID, opt.ObjID, opt.PropertyID, opt.Data)
		if ccErr != nil {
			return ccErr
		}

		err = audit.SaveAuditLog(cts.Kit, auditLogs...)
		if err != nil {
			return cts.Kit.CCError.Error(common.CCErrAuditSaveLogFailed)
		}

		res.IDs = ids
		return nil
	})

	if txnErr != nil {
		cts.RespAutoError(txnErr)
		return
	}
	cts.RespEntity(res)
}

// ListQuotedInstance list quoted instances.
func (s *Service) ListQuotedInstance(cts *rest.Contexts) {
	opt := new(metadata.ListQuotedInstOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// skip find authorize, ** add it when confirmed **

	// get quoted object id
	objID, err := s.Logics.ModelQuoteOperation().GetQuotedObjID(cts.Kit, opt.ObjID, opt.PropertyID)
	if err != nil {
		cts.RespAutoError(err)
		return
	}

	// list quoted instances
	res, err := s.Engine.CoreAPI.CoreService().ModelQuote().ListQuotedInstance(cts.Kit.Ctx, cts.Kit.Header, objID,
		&opt.CommonQueryOption)
	if err != nil {
		blog.Errorf("list quoted instances failed, err: %v, req: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}

// BatchUpdateQuotedInstance batch update quoted instances.
func (s *Service) BatchUpdateQuotedInstance(cts *rest.Contexts) {
	opt := new(metadata.BatchUpdateQuotedInstOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err.ErrCode != 0 {
		cts.RespAutoError(err.ToCCError(cts.Kit.CCError))
		return
	}

	// get quoted object id
	objID, err := s.Logics.ModelQuoteOperation().GetQuotedObjID(cts.Kit, opt.ObjID, opt.PropertyID)
	if err != nil {
		cts.RespAutoError(err)
		return
	}

	filterOpt := metadata.CommonFilterOption{
		Filter: filtertools.GenAtomFilter(common.BKFieldID, filter.In, opt.IDs),
	}

	// get quoted instance info for audit and authorization
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: filterOpt,
		Page:               metadata.BasePage{Limit: common.BKMaxPageSize},
	}
	listRes, err := s.Engine.CoreAPI.CoreService().ModelQuote().ListQuotedInstance(cts.Kit.Ctx, cts.Kit.Header, objID,
		listOpt)
	if err != nil {
		blog.Errorf("list quoted instance failed, err: %v, opt: %+v, rid: %s", err, listOpt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	// authorize, ** right now use source instance update action to authorize, change it when confirmed **
	if err = s.authorizeQuotedInstance(cts.Kit, objID, listRes.Info); err != nil {
		return
	}

	// generate audit logs
	audit := auditlog.NewQuotedInstAuditLog(s.Engine.CoreAPI.CoreService())
	genAuditParams := auditlog.NewGenerateAuditCommonParameter(cts.Kit, metadata.AuditUpdate).WithUpdateFields(opt.Data)
	auditLogs, ccErr := audit.GenerateAuditLog(genAuditParams, objID, opt.ObjID, opt.PropertyID, listRes.Info)
	if ccErr != nil {
		cts.RespAutoError(ccErr)
		return
	}

	updateOpt := &metadata.CommonUpdateOption{
		CommonFilterOption: filterOpt,
		Data:               opt.Data,
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(cts.Kit.Ctx, cts.Kit.Header, func() error {
		// update quoted instances
		if err = s.Engine.CoreAPI.CoreService().ModelQuote().BatchUpdateQuotedInstance(cts.Kit.Ctx, cts.Kit.Header,
			objID, updateOpt); err != nil {
			blog.Errorf("update quoted instances failed, err: %v, rid: %s", err, cts.Kit.Rid)
			return err
		}

		// save audit logs
		err = audit.SaveAuditLog(cts.Kit, auditLogs...)
		if err != nil {
			return cts.Kit.CCError.Error(common.CCErrAuditSaveLogFailed)
		}

		return nil
	})

	if txnErr != nil {
		cts.RespAutoError(txnErr)
		return
	}
	cts.RespEntity(nil)
}

// BatchDeleteQuotedInstance batch delete quoted instances.
func (s *Service) BatchDeleteQuotedInstance(cts *rest.Contexts) {
	opt := new(metadata.BatchDeleteQuotedInstOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err.ErrCode != 0 {
		cts.RespAutoError(err.ToCCError(cts.Kit.CCError))
		return
	}

	// get quoted object id
	objID, err := s.Logics.ModelQuoteOperation().GetQuotedObjID(cts.Kit, opt.ObjID, opt.PropertyID)
	if err != nil {
		cts.RespAutoError(err)
		return
	}

	filterOpt := metadata.CommonFilterOption{
		Filter: filtertools.GenAtomFilter(common.BKFieldID, filter.In, opt.IDs),
	}

	// get quoted instance info for audit and authorization
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: filterOpt,
		Page:               metadata.BasePage{Limit: common.BKMaxPageSize},
	}
	listRes, err := s.Engine.CoreAPI.CoreService().ModelQuote().ListQuotedInstance(cts.Kit.Ctx, cts.Kit.Header, objID,
		listOpt)
	if err != nil {
		blog.Errorf("list quoted instance failed, err: %v, opt: %+v, rid: %s", err, listOpt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	// authorize, ** right now use source instance update action to authorize, change it when confirmed **
	if err = s.authorizeQuotedInstance(cts.Kit, objID, listRes.Info); err != nil {
		return
	}

	// generate audit logs
	audit := auditlog.NewQuotedInstAuditLog(s.Engine.CoreAPI.CoreService())
	genAuditParams := auditlog.NewGenerateAuditCommonParameter(cts.Kit, metadata.AuditDelete)
	auditLogs, ccErr := audit.GenerateAuditLog(genAuditParams, objID, opt.ObjID, opt.PropertyID, listRes.Info)
	if ccErr != nil {
		cts.RespAutoError(ccErr)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(cts.Kit.Ctx, cts.Kit.Header, func() error {
		// delete quoted instances
		if err := s.Engine.CoreAPI.CoreService().ModelQuote().BatchDeleteQuotedInstance(cts.Kit.Ctx, cts.Kit.Header,
			objID, &filterOpt); err != nil {
			blog.Errorf("delete quoted instances failed, err: %v, rid: %s", err, cts.Kit.Rid)
			return err
		}

		// save audit logs
		err = audit.SaveAuditLog(cts.Kit, auditLogs...)
		if err != nil {
			return cts.Kit.CCError.Error(common.CCErrAuditSaveLogFailed)
		}

		return nil
	})

	if txnErr != nil {
		cts.RespAutoError(txnErr)
		return
	}
	cts.RespEntity(nil)
}

func (s *Service) authorizeQuotedInstance(kit *rest.Kit, objID string, data []mapstr.MapStr) error {
	instIDs := make([]int64, 0)
	for _, info := range data {
		instIDVal, exists := info[common.BKInstIDField]
		if !exists {
			continue
		}

		instID, err := util.GetInt64ByInterface(instIDVal)
		if err != nil {
			blog.Errorf("parse inst id failed, err: %v, id: %+v, rid: %s", err, instIDVal, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKInstIDField)
		}
		instIDs = append(instIDs, instID)
	}

	instIDs = util.IntArrayUnique(instIDs)

	err := s.AuthManager.AuthorizeByInstanceID(kit.Ctx, kit.Header, meta.Update, objID, instIDs...)
	if err != nil {
		blog.Errorf("authorize failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}
