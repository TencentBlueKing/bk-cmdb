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

// Package common defines the common resource cache client
package common

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

// Interface is the common resource cache client interface
type Interface interface {
	ListWithKey(ctx context.Context, h http.Header, typ string, opt *metadata.ListCommonCacheWithKeyOpt) (string, error)
}

// NewCacheClient new common resource cache client
func NewCacheClient(client rest.ClientInterface) Interface {
	return &cache{client: client}
}

type cache struct {
	client rest.ClientInterface
}
