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
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
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
		instIDs = append(instIDs, instID)
	}

	instIDs = util.IntArrayUnique(instIDs)

	uAuthErr := s.AuthManager.AuthorizeByInstanceID(cts.Kit.Ctx, cts.Kit.Header, meta.Update, opt.ObjID, instIDs...)
	cAuthErr := s.AuthManager.AuthorizeByInstanceID(cts.Kit.Ctx, cts.Kit.Header, meta.Create, opt.ObjID, instIDs...)
	if uAuthErr != nil && cAuthErr != nil {
		blog.Errorf("authorize failed, create err: %v, update err: %v, rid: %s", cAuthErr, uAuthErr, cts.Kit.Rid)
		cts.RespAutoError(uAuthErr)
		return
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
