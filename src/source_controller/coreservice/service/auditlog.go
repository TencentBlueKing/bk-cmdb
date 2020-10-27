/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (s *coreService) CreateAuditLog(ctx *rest.Contexts) {
	inputData := new(metadata.CreateAuditLogParam)

	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := s.core.AuditOperation().CreateAuditLog(ctx.Kit, inputData.Data...); nil != err {
		blog.Errorf("CreateAuditLog err:%v, rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrAuditSaveLogFailed))
		return
	}

	ctx.RespEntity(nil)
}

func (s *coreService) SearchAuditLog(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	auditLogs, count, err := s.core.AuditOperation().SearchAuditLog(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("SearchAuditLog err:%v, rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrAuditSelectFailed))
		return
	}

	ctx.RespEntityWithCount(int64(count), auditLogs)
}
