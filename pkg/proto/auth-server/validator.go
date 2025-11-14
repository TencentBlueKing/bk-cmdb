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

// Package authpb is the grpc generated code and related logics for auth server.
package authpb

import (
	"context"
	"errors"
)

// Validate validate authorize request.
func (r *AuthorizeReq) Validate(ctx context.Context) error {
	if len(r.Resources) == 0 {
		return errors.New("resources are not set")
	}

	for _, resource := range r.Resources {
		if resource.Basic == nil {
			return errors.New("resource basic is not set")
		}
	}

	return nil
}

// Validate validate list auth resource request.
func (r *ListAuthResReq) Validate(ctx context.Context) error {
	if r.ResourceType == "" {
		return errors.New("resource type is not set")
	}

	if r.Action == "" {
		return errors.New("action is not set")
	}
	return nil
}
