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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) CreateAuditLog(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := struct {
		Data []metadata.SaveAuditLogParams `json:"data"`
	}{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	return nil, s.core.AuditOperation().CreateAuditLog(params, inputData.Data...)
}

func (s *coreService) SearchAuditLog(ctx core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	inputData := metadata.QueryInput{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}
	auditlogs, count, err := s.core.AuditOperation().SearchAuditLog(ctx, inputData)
	return struct {
		Count uint64                  `json:"count"`
		Info  []metadata.OperationLog `json:"info"`
	}{
		Count: count,
		Info:  auditlogs,
	}, err
}
