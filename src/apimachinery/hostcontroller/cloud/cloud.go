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

package cloud

import (
	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
	"context"
	"net/http"
)

type CloudInterface interface {
	AddCloudTask(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)
	ResourceConfirm(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)
	TaskNameCheck(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Uint64Response, err error)
	DeleteCloudTask(ctx context.Context, h http.Header, taskID string) (resp *metadata.Response, err error)
	SearchCloudTask(ctx context.Context, h http.Header, data interface{}) (resp *metadata.CloudTaskSearch, err error)
	UpdateCloudTask(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)
	DeleteConfirm(ctx context.Context, h http.Header, ResourceID int64) (resp *metadata.Response, err error)
	SearchConfirm(ctx context.Context, h http.Header, data interface{}) (resp *metadata.FavoriteResult, err error)
	AddSyncHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)
	SearchSyncHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.FavoriteResult, err error)
	AddConfirmHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)
	SearchConfirmHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.FavoriteResult, err error)
}

func NewCloudInterface(client rest.ClientInterface) CloudInterface {
	return &cloud{client: client}
}

type cloud struct {
	client rest.ClientInterface
}
