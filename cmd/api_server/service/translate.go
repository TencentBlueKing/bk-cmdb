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

package service

import (
	"context"

	ccError "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
)

// TranslateInfoReq translate info request
type TranslateInfoReq struct {
	Context string `req:"context,in:query"`
}

// TranslateInfoResp translate info response
type TranslateInfoResp struct {
	Context string `json:"context"`
}

// Translate ...
func (s *service) Translate(ctx context.Context, req *TranslateInfoReq) (*TranslateInfoResp, *ccError.RespError) {
	resp := &TranslateInfoResp{
		Context: s.T(ctx, req.Context),
	}
	return resp, nil
}
