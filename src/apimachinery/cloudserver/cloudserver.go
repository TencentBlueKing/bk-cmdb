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

package cloudserver

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/metadata"
)

type CloudServerClientInterface interface {
	// cloud account
	CreateAccount(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	SearchAccount(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.SearchResp, err error)
	UpdateAccount(ctx context.Context, h http.Header, accountID int64, data map[string]interface{}) (resp *metadata.Response, err error)
	DeleteAccount(ctx context.Context, h http.Header, accountID int64) (resp *metadata.Response, err error)

	CreateSyncTask(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	SearchSyncTask(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.SearchResp, err error)
	UpdateSyncTask(ctx context.Context, h http.Header, taskID int64, data map[string]interface{}) (resp *metadata.Response, err error)
	DeleteSyncTask(ctx context.Context, h http.Header, taskID int64) (resp *metadata.Response, err error)
	SearchSyncHistory(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.SearchResp, err error)
	SearchSyncRegion(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.SearchResp, err error)
}

func NewCloudServerClientInterface(c *util.Capability, version string) CloudServerClientInterface {
	base := fmt.Sprintf("/cloud/%s", version)

	return &cloudserver{
		client: rest.NewRESTClient(c, base),
	}
}

type cloudserver struct {
	client rest.ClientInterface
}
