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
	"configcenter/src/common/http/rest"
)

// 云账户连通测试
func (s *Service) VerifyConnectivity(ctx *rest.Contexts) {
	ctx.RespEntity("VerifyConnectivity")
}

// 新建云账户
func (s *Service) CreateAccount(ctx *rest.Contexts) {
	ctx.RespEntity("CreateAccount")
}

// 查询云账户
func (s *Service) SearchAccount(ctx *rest.Contexts) {
	ctx.RespEntity("SearchAccount")
}

// 更新云账户
func (s *Service) UpdateAccount(ctx *rest.Contexts) {
	ctx.RespEntity("UpdateAccount")
}

// 删除云账户
func (s *Service) DeleteAccount(ctx *rest.Contexts) {
	ctx.RespEntity("DeleteAccount")
}
