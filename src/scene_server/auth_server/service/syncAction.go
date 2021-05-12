/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"context"
	"net/http"
	"time"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/app/options"
)

type SyncServer struct {
	Config  *options.Config
	Service *AuthService
}

// 根据配置文件, 每隔固定时间就同步下IAM的action列表, 将IAM中多余的action删除。
func (s *AuthService) LoopSyncActionWithIAM(ctx context.Context, config *options.Config) {
	kit := s.NewKit(config)
	timer := time.NewTimer(config.Auth.Interval)
	for true {
		select {
		// 计时器信号
		case <-timer.C:
			err := s.SyncIAMModelResources(*kit)
			if err != nil {
				blog.Errorf("sync action with IAM failed, err:%v", err)
			}
		// authServer退出信号
		case <-ctx.Done():
			blog.Infof("auth server will exit!")
			return
		}
	}
}

// NewHeader
func (s *AuthService) NewHeader(config *options.Config) http.Header {
	header := make(http.Header)
	header.Add(common.BKHTTPOwnerID, common.BKSuperOwnerID)
	header.Add(common.BKHTTPHeaderUser, common.BKAuthUser)
	header.Add(common.BKHTTPLanguage, "cn")
	header.Add(common.BKHTTPCCRequestID, util.GenerateRID())
	header.Add("Content-Type", "application/json")

	header.Add(iam.IamAppCodeHeader, config.Auth.AppCode)
	header.Add(iam.IamAppSecretHeader, config.Auth.AppSecret)
	return header
}

// NewKit
func (s *AuthService) NewKit(config *options.Config) *rest.Kit {
	header := s.NewHeader(config)

	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.GetHTTPCCRequestID(header)
	user := util.GetUser(header)
	supplierAccount := util.GetOwnerID(header)
	defaultCCError := util.GetDefaultCCError(header)

	return &rest.Kit{
		Rid:             rid,
		Header:          header,
		Ctx:             ctx,
		CCError:         defaultCCError,
		User:            user,
		SupplierAccount: supplierAccount,
	}
}
