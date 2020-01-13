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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (s *coreService) CreateAuditLog(ctx *rest.Contexts) {
	inputData := struct {
		Data []metadata.SaveAuditLogParams `json:"data"`
	}{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	if err := s.core.AuditOperation().CreateAuditLog(ctx.Kit, inputData.Data...); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) SearchAuditLog(ctx *rest.Contexts) {
	inputData := metadata.QueryInput{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	auditLogs, count, err := s.core.AuditOperation().SearchAuditLog(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(struct {
		Count uint64                  `json:"count"`
		Info  []metadata.OperationLog `json:"info"`
	}{
		Count: count,
		Info:  auditLogs,
	})
}
