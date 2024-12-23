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

// Package general defines the general resource cache client
package general

import (
	"context"
	"net/http"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/pkg/cache/general"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
)

// Interface is the general resource cache client interface
type Interface interface {
	CreateFullSyncCond(ctx context.Context, h http.Header, opt *fullsynccond.CreateFullSyncCondOpt) (int64,
		errors.CCErrorCoder)
	UpdateFullSyncCond(ctx context.Context, h http.Header, opt *fullsynccond.UpdateFullSyncCondOpt) errors.CCErrorCoder
	DeleteFullSyncCond(ctx context.Context, h http.Header, opt *fullsynccond.DeleteFullSyncCondOpt) errors.CCErrorCoder
	ListFullSyncCond(ctx context.Context, h http.Header, opt *fullsynccond.ListFullSyncCondOpt) (
		*fullsynccond.ListFullSyncCondRes, errors.CCErrorCoder)
	ListCacheByFullSyncCond(ctx context.Context, h http.Header, opt *fullsynccond.ListCacheByFullSyncCondOpt) (
		*general.ListGeneralCacheRes, errors.CCErrorCoder)
	ListGeneralCacheByIDs(ctx context.Context, h http.Header, opt *general.ListDetailByIDsOpt) (
		*general.ListGeneralCacheRes, errors.CCErrorCoder)
	ListGeneralCacheByUniqueKey(ctx context.Context, h http.Header, opt *general.ListDetailByUniqueKeyOpt) (
		*general.ListGeneralCacheRes, errors.CCErrorCoder)
}

// NewCacheClient new general resource cache client
func NewCacheClient(client rest.ClientInterface) Interface {
	return &cache{client: client}
}

type cache struct {
	client rest.ClientInterface
}
