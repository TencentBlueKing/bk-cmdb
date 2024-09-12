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

// Package synchronize defines multiple cmdb synchronize service
package synchronize

import (
	"configcenter/pkg/synchronize/types"
	acmeta "configcenter/src/ac/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
)

// CreateSyncData create sync data
func (s *service) CreateSyncData(cts *rest.Contexts) {
	option := new(types.CreateSyncDataOption)
	if err := cts.DecodeInto(option); err != nil {
		cts.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		cts.RespAutoError(err.ToCCError(cts.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.SynchronizeData, Action: acmeta.Create}}
	if resp, authorized := s.AuthManager.Authorize(cts.Kit, authRes); !authorized {
		cts.RespNoAuth(resp)
		return
	}

	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(cts.Kit.Ctx, cts.Kit.Header, func() error {
		var err error
		err = s.ClientSet.CoreService().Synchronize().CreateSyncData(cts.Kit.Ctx, cts.Kit.Header, option)
		if err != nil {
			blog.Errorf("create sync data failed, err: %v, option: %+v, rid: %s", err, *option, cts.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		cts.RespAutoError(txnErr)
		return
	}

	cts.RespEntity(nil)
}
