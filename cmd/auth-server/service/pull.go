/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

// Package service defines auth-server service.
package service

import (
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
)

// PullResource is iam pull resource callback function.
func (s *Service) PullResource(kt *kit.Kit, req *PullResourceReq) (*PullResourceResp, error) {
	return &PullResourceResp{IDs: []string{"test"}}, nil
}

// PullResourceReq is iam pull resource request.
type PullResourceReq struct {
	Type string `json:"type"`
}

// PullResourceResp is iam pull resource response.
type PullResourceResp struct {
	IDs []string `json:"ids"`
}
