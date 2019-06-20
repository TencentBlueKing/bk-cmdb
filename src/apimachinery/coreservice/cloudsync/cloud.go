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

package cloudsync

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type CloudSyncClientInterface interface {
	CreateCloudSyncTask(ctx context.Context, h http.Header, input interface{}) (resp *metadata.Uint64DataResponse, err error)
	DeleteCloudSyncTask(ctx context.Context, h http.Header, id int64) (resp *metadata.Response, err error)
	UpdateCloudSyncTask(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)
	SearchCloudSyncTask(ctx context.Context, h http.Header, data interface{}) (resp *metadata.CloudTaskSearch, err error)
	CreateConfirm(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Uint64DataResponse, err error)
	DeleteConfirm(ctx context.Context, h http.Header, id int64) (resp *metadata.Response, err error)
	SearchConfirm(ctx context.Context, h http.Header, data interface{}) (resp *metadata.FavoriteResult, err error)
	CreateSyncHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Uint64Response, err error)
	SearchSyncHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.FavoriteResult, err error)
	CreateConfirmHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)
	SearchConfirmHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.FavoriteResult, err error)
	CheckTaskNameUnique(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Uint64Response, err error)
}

func NewCloudSyncClientInterface(client rest.ClientInterface) CloudSyncClientInterface {
	return &cloud{client: client}
}

type cloud struct {
	client rest.ClientInterface
}
