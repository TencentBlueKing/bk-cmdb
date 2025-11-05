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
	"errors"
	"sync"
	"time"

	"github.com/TencentBlueKing/bk-cmdb/pkg/auth/meta"
	"github.com/TencentBlueKing/bk-cmdb/pkg/i18n"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest"
)

// UserInfoReq 个人信息Req
type UserInfoReq struct {
	Username string     `req:"name,in:query" validate:"required"`
	Age      int        `req:"age,in:query" validate:"required"`
	Games    *[]*string `json:"games" req:"-" validate:"required"`
	BirthDay time.Time  `req:"birthday,in:query,format:2006-01-02" validate:"required"`
	Ko       []byte     `json:"-" req:"ko,in:query" validate:"required"`
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
func (s *Service) UserInfo(kt *kit.Kit, req *UserInfoReq) (*UserInfoResp, error) {
	log.Info(kt, "handle UserInfo")

	var wg sync.WaitGroup

	wg.Go(func() { doBiz(kt) })
	wg.Go(func() { doBiz(kt) })
	wg.Wait()

	// authorize, NOTE: this is only a demo
	authRes := meta.ResourceAttribute{
		Basic: &meta.Basic{Type: "skip"},
	}
	decisions, err := s.authorizer.Authorize(kt, authRes)
	if err != nil || !decisions[0].Authorized {
		log.Error(kt, "authorize failed", log.E(err), "decisions", decisions)
		return nil, errors.New("authorize failed")
	}

	resp := &UserInfoResp{
		Username: i18n.Sys(kt, req.Username),
		Age:      req.Age + 10,
		Games:    req.Games,
		Ko:       string(req.Ko),
		BirthDay: req.BirthDay,
	}
	return resp, nil
}

func doBiz(kt *kit.Kit) {
	kt, span := kt.StartSpan("")
	defer span.End()

	log.Info(kt, "do Biz")
}

// ListAuthorizedUsers 获取有权限的用户
func (s *Service) ListAuthorizedUsers(kt *kit.Kit, _ *rest.EmptyReq) (*[]UserInfoResp, error) {
	log.Info(kt, "handle ListAuthorizedUsers")

	authOpt := &meta.ListAuthResOptions{ResourceType: "user", Action: "list"}
	resInfo, err := s.authorizer.ListAuthorizedResources(kt, authOpt)
	if err != nil {
		log.Error(kt, "list authorized resources failed", log.E(err))
		return nil, err
	}

	users := make([]UserInfoResp, len(resInfo.IDs))
	for i, id := range resInfo.IDs {
		users[i] = UserInfoResp{
			Username: id,
		}
	}
	return &users, nil
}
