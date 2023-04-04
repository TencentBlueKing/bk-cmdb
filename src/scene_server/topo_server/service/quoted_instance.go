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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
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

		// TODO generate and save audit logs

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

	updateOpt := &metadata.CommonUpdateOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: filtertools.GenAtomFilter(common.BKFieldID, filter.Equal, opt.IDs),
		},
		Data: opt.Data,
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(cts.Kit.Ctx, cts.Kit.Header, func() error {
		// TODO generate audit logs

		// update quoted instances
		if err = s.Engine.CoreAPI.CoreService().ModelQuote().BatchUpdateQuotedInstance(cts.Kit.Ctx, cts.Kit.Header,
			objID, updateOpt); err != nil {
			blog.Errorf("update quoted instances failed, err: %v, rid: %s", err, cts.Kit.Rid)
			return err
		}

		// TODO save audit logs

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

	deleteOpt := &metadata.CommonFilterOption{
		Filter: filtertools.GenAtomFilter(common.BKFieldID, filter.Equal, opt.IDs),
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(cts.Kit.Ctx, cts.Kit.Header, func() error {
		// TODO generate audit logs

		// delete quoted instances
		if err := s.Engine.CoreAPI.CoreService().ModelQuote().BatchDeleteQuotedInstance(cts.Kit.Ctx, cts.Kit.Header,
			objID, deleteOpt); err != nil {
			blog.Errorf("delete quoted instances failed, err: %v, rid: %s", err, cts.Kit.Rid)
			return err
		}

		// TODO save audit logs

		return nil
	})

	if txnErr != nil {
		cts.RespAutoError(txnErr)
		return
	}
	cts.RespEntity(nil)
}
