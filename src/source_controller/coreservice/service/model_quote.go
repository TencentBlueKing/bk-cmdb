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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

// ListModelQuoteRelation list model quote relationships.
func (s *coreService) ListModelQuoteRelation(cts *rest.Contexts) {
	req := new(metadata.CommonQueryOption)
	if err := cts.DecodeInto(req); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	filter, err := req.ToMgo()
	if err != nil {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	if req.Page.EnableCount {
		count, err := mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Find(filter).Count(cts.Kit.Ctx)
		if err != nil {
			blog.Errorf("count model quote relations failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
			cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}

		cts.RespEntity(metadata.ListModelQuoteRelRes{Count: count})
		return
	}

	relations := make([]metadata.ModelQuoteRelation, 0)
	err = mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Find(filter).Start(uint64(req.Page.Start)).
		Limit(uint64(req.Page.Limit)).Fields(req.Fields...).All(cts.Kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("list model quote relations failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	cts.RespEntity(metadata.ListModelQuoteRelRes{Info: relations})
}

// CreateModelQuoteRelation create model quote relationships.
func (s *coreService) CreateModelQuoteRelation(cts *rest.Contexts) {
	relations := make([]metadata.ModelQuoteRelation, 0)
	if err := cts.DecodeInto(&relations); err != nil {
		cts.RespAutoError(err)
		return
	}

	if len(relations) == 0 || len(relations) > common.BKMaxLimitSize {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrArrayLengthWrong, "relations", common.BKMaxLimitSize))
		return
	}

	for idx := range relations {
		relations[idx].TenantID = cts.Kit.TenantID
	}

	err := mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Insert(cts.Kit.Ctx, relations)
	if err != nil {
		blog.Errorf("create model quote relations failed, err: %v, data: %+v, rid: %v", err, relations, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
		return
	}

	cts.RespEntity(nil)
}

// DeleteModelQuoteRelation delete model quote relationships.
func (s *coreService) DeleteModelQuoteRelation(cts *rest.Contexts) {
	req := new(metadata.CommonFilterOption)
	if err := cts.DecodeInto(req); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	filter, err := req.ToMgo()
	if err != nil {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	err = mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Delete(cts.Kit.Ctx, filter)
	if err != nil {
		blog.Errorf("delete model quote relations failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBDeleteFailed))
		return
	}

	cts.RespEntity(nil)
}
