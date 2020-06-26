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
	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common/http/rest"
)

// Authorize works to check if a user has the authority to operate resources
func (s *AuthService) Authorize(ctx *rest.Contexts) {
	authAttribute := new(meta.AuthAttribute)
	err := ctx.DecodeInto(authAttribute)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	// TODO implement this
	ctx.RespEntity(meta.Decision{Authorized: true})
}

// AuthorizeBath works to check if a user has the authority to operate resources.
func (s *AuthService) AuthorizeBatch(ctx *rest.Contexts) {
	authAttribute := new(meta.AuthAttribute)
	err := ctx.DecodeInto(authAttribute)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	// TODO implement this
	decisions := make([]meta.Decision, len(authAttribute.Resources))
	for key, _ := range authAttribute.Resources {
		decisions[key] = meta.Decision{Authorized: true}
	}
	ctx.RespEntity(decisions)
}

// ListAuthorizedResources returns all specified resources the user has the authority to operate.
func (s *AuthService) ListAuthorizedResources(ctx *rest.Contexts) {
	authAttribute := new(meta.ListAuthorizedResourcesParam)
	err := ctx.DecodeInto(authAttribute)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	// TODO implement this
	ctx.RespEntity([]iam.IamResource{})
}
