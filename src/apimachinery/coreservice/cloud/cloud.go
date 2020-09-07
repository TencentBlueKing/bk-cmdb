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
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

type CloudInterface interface {
	// cloud account
	CreateAccount(ctx context.Context, h http.Header, account *metadata.CloudAccount) (*metadata.CloudAccount, errors.CCErrorCoder)
	SearchAccount(ctx context.Context, h http.Header, option *metadata.SearchCloudOption) (*metadata.MultipleCloudAccount, errors.CCErrorCoder)
	UpdateAccount(ctx context.Context, h http.Header, accountID int64, option map[string]interface{}) errors.CCErrorCoder
	DeleteAccount(ctx context.Context, h http.Header, accountID int64) errors.CCErrorCoder
	SearchAccountConf(ctx context.Context, h http.Header, option *metadata.SearchCloudOption) (*metadata.MultipleCloudAccountConf, errors.CCErrorCoder)

	CreateSyncTask(ctx context.Context, h http.Header, account *metadata.CloudSyncTask) (*metadata.CloudSyncTask, errors.CCErrorCoder)
	SearchSyncTask(ctx context.Context, h http.Header, option *metadata.SearchCloudOption) (*metadata.MultipleCloudSyncTask, errors.CCErrorCoder)
	UpdateSyncTask(ctx context.Context, h http.Header, taskID int64, option map[string]interface{}) errors.CCErrorCoder
	DeleteSyncTask(ctx context.Context, h http.Header, taskID int64) errors.CCErrorCoder
	CreateSyncHistory(ctx context.Context, h http.Header, history *metadata.SyncHistory) (*metadata.SyncHistory, errors.CCErrorCoder)
	SearchSyncHistory(ctx context.Context, h http.Header, option *metadata.SearchSyncHistoryOption) (*metadata.MultipleSyncHistory, errors.CCErrorCoder)
	DeleteDestroyedHostRelated(ctx context.Context, h http.Header, option *metadata.DeleteDestroyedHostRelatedOption) errors.CCErrorCoder
}

func NewCloudInterfaceClient(client rest.ClientInterface) CloudInterface {
	return &cloud{client: client}
}

type cloud struct {
	client rest.ClientInterface
}
