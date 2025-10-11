/*
 * Tencent is pleased to support the open source community by making
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
	"time"
)

// UserInfoReq 个人信息Req
type UserInfoReq struct {
	Username string     `json:"name" req:"query:name"`
	Age      int        `req:"query:age"`
	Games    *[]*string `json:"games" req:"query:games"`
	BirthDay time.Time  `json:"birthday" req:"query:birthday,format:2006-01-02"`
	Ko       []byte     `json:"ko" req:"query:ko"`
}

// UserInfoResp 个人信息
type UserInfoResp struct {
	Username string     `json:"username"`
	Age      int        `json:"age"`
	Games    *[]*string `json:"games"`
	Ko       string     `json:"ko"`
	BirthDay time.Time  `json:"birthday,format:'2006-01-02'"`
}

// UserInfo 用户信息
func (s *service) UserInfo(ctx context.Context, req *UserInfoReq) (*UserInfoResp, error) {
	resp := &UserInfoResp{
		Username: req.Username,
		Age:      req.Age + 10,
		Games:    req.Games,
		Ko:       string(req.Ko),
		BirthDay: req.BirthDay,
	}
	return resp, nil
}
